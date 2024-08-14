// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package controlsocket

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"path"
	"strings"

	jujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/core/permission"
	coreuser "github.com/juju/juju/core/user"
	usertesting "github.com/juju/juju/core/user/testing"
	usererrors "github.com/juju/juju/domain/access/errors"
	"github.com/juju/juju/domain/access/service"
	auth "github.com/juju/juju/internal/auth"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	jujujujutesting "github.com/juju/juju/internal/testing"
)

type workerSuite struct {
	jujutesting.IsolationSuite

	logger        logger.Logger
	accessService *MockAccessService

	controllerModelID permission.ID
	metricsUserName   coreuser.Name
}

var _ = gc.Suite(&workerSuite{})

func (s *workerSuite) SetUpTest(c *gc.C) {
	s.logger = loggertesting.WrapCheckLog(c)
	s.metricsUserName = usertesting.GenNewName(c, "juju-metrics-r0")
	s.controllerModelID = permission.ID{
		ObjectType: permission.Model,
		Key:        jujujujutesting.ModelTag.Id(),
	}
}

type handlerTest struct {
	// Request
	method   string
	endpoint string
	body     string
	// Response
	statusCode int
	response   string // response body
	ignoreBody bool   // if true, test will not read the request body
}

func (s *workerSuite) runHandlerTest(c *gc.C, test handlerTest) {
	tmpDir := c.MkDir()
	socket := path.Join(tmpDir, "test.socket")

	_, err := NewWorker(Config{
		AccessService:       s.accessService,
		Logger:              s.logger,
		SocketName:          socket,
		NewSocketListener:   NewSocketListener,
		ControllerModelUUID: model.UUID(jujujujutesting.ModelTag.Id()),
	})
	c.Assert(err, jc.ErrorIsNil)

	serverURL := "http://localhost:8080"
	req, err := http.NewRequest(
		test.method,
		serverURL+test.endpoint,
		strings.NewReader(test.body),
	)
	c.Assert(err, jc.ErrorIsNil)

	resp, err := client(socket).Do(req)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(resp.StatusCode, gc.Equals, test.statusCode)

	if test.ignoreBody {
		return
	}
	data, err := io.ReadAll(resp.Body)
	c.Assert(err, jc.ErrorIsNil)
	err = resp.Body.Close()
	c.Assert(err, jc.ErrorIsNil)

	// Response should be valid JSON
	c.Check(resp.Header.Get("Content-Type"), gc.Equals, "application/json")
	err = json.Unmarshal(data, &struct{}{})
	c.Assert(err, jc.ErrorIsNil)
	if test.response != "" {
		c.Check(string(data), gc.Matches, test.response)
	}
}

func (s *workerSuite) TestMetricsUsersAddInvalidMethod(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodGet,
		endpoint:   "/metrics-users",
		statusCode: http.StatusMethodNotAllowed,
		ignoreBody: true,
	})
}

func (s *workerSuite) TestMetricsUsersAddMissingBody(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		statusCode: http.StatusBadRequest,
		response:   ".*missing request body.*",
	})
}

func (s *workerSuite) TestMetricsUsersAddInvalidBody(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       "username foo, password bar",
		statusCode: http.StatusBadRequest,
		response:   ".*request body is not valid JSON.*",
	})
}

func (s *workerSuite) TestMetricsUsersAddMissingUsername(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       `{"password":"bar"}`,
		statusCode: http.StatusBadRequest,
		response:   ".*missing username.*",
	})
}

func (s *workerSuite) TestMetricsUsersAddUsernameMissingPrefix(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       `{"username":"foo","password":"bar"}`,
		statusCode: http.StatusBadRequest,
		response:   `.*username .* should have prefix \\\"juju-metrics-\\\".*`,
	})
}

func (s *workerSuite) TestMetricsUsersAddSuccess(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), usertesting.GenNewName(c, userCreator)).Return(coreuser.User{
		UUID: coreuser.UUID("deadbeef"),
	}, nil)
	s.accessService.EXPECT().AddUser(gomock.Any(), service.AddUserArg{
		Name:        s.metricsUserName,
		DisplayName: "juju-metrics-r0",
		Password:    ptr(auth.NewPassword("bar")),
		CreatorUUID: coreuser.UUID("deadbeef"),
		Permission: permission.AccessSpec{
			Target: s.controllerModelID,
			Access: permission.ReadAccess,
		},
	}).Return(coreuser.UUID("foobar"), nil, nil)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       `{"username":"juju-metrics-r0","password":"bar"}`,
		statusCode: http.StatusOK,
		response:   `.*created user \\\"juju-metrics-r0\\\".*`,
	})
}

func (s *workerSuite) TestMetricsUsersAddAlreadyExists(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), usertesting.GenNewName(c, userCreator)).Return(coreuser.User{
		UUID: coreuser.UUID("deadbeef"),
	}, nil)
	s.accessService.EXPECT().AddUser(gomock.Any(), service.AddUserArg{
		Name:        s.metricsUserName,
		DisplayName: "juju-metrics-r0",
		Password:    ptr(auth.NewPassword("bar")),
		CreatorUUID: coreuser.UUID("deadbeef"),
		Permission: permission.AccessSpec{
			Target: s.controllerModelID,
			Access: permission.ReadAccess,
		},
	}).Return(coreuser.UUID("foobar"), nil, usererrors.UserAlreadyExists)
	s.accessService.EXPECT().GetUserByAuth(gomock.Any(), s.metricsUserName, auth.NewPassword("bar")).Return(coreuser.User{
		CreatorName: usertesting.GenNewName(c, "not-you"),
	}, nil)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       `{"username":"juju-metrics-r0","password":"bar"}`,
		statusCode: http.StatusConflict,
		response:   ".*user .* already exists.*",
	})
}

func (s *workerSuite) TestMetricsUsersAddAlreadyExistsButDisabled(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), usertesting.GenNewName(c, userCreator)).Return(coreuser.User{
		UUID: coreuser.UUID("deadbeef"),
	}, nil)
	s.accessService.EXPECT().AddUser(gomock.Any(), service.AddUserArg{
		Name:        s.metricsUserName,
		DisplayName: "juju-metrics-r0",
		Password:    ptr(auth.NewPassword("bar")),
		CreatorUUID: coreuser.UUID("deadbeef"),
		Permission: permission.AccessSpec{
			Target: s.controllerModelID,
			Access: permission.ReadAccess,
		},
	}).Return(coreuser.UUID("foobar"), nil, usererrors.UserAlreadyExists)
	s.accessService.EXPECT().GetUserByAuth(gomock.Any(), s.metricsUserName, auth.NewPassword("bar")).Return(coreuser.User{
		CreatorName: usertesting.GenNewName(c, "not-you"),
		Disabled:    true,
	}, nil)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       `{"username":"juju-metrics-r0","password":"bar"}`,
		statusCode: http.StatusForbidden,
		response:   ".*user .* is disabled.*",
	})
}

func (s *workerSuite) TestMetricsUsersAddAlreadyExistsButWrongPermissions(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), usertesting.GenNewName(c, userCreator)).Return(coreuser.User{
		UUID: coreuser.UUID("deadbeef"),
	}, nil)
	s.accessService.EXPECT().AddUser(gomock.Any(), service.AddUserArg{
		Name:        s.metricsUserName,
		DisplayName: "juju-metrics-r0",
		Password:    ptr(auth.NewPassword("bar")),
		CreatorUUID: coreuser.UUID("deadbeef"),
		Permission: permission.AccessSpec{
			Target: s.controllerModelID,
			Access: permission.ReadAccess,
		},
	}).Return(coreuser.UUID("foobar"), nil, usererrors.UserAlreadyExists)
	s.accessService.EXPECT().GetUserByAuth(gomock.Any(), s.metricsUserName, auth.NewPassword("bar")).Return(coreuser.User{
		CreatorName: usertesting.GenNewName(c, userCreator),
	}, nil)
	s.accessService.EXPECT().ReadUserAccessLevelForTarget(gomock.Any(), s.metricsUserName, s.controllerModelID).Return(
		permission.WriteAccess, nil,
	)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       `{"username":"juju-metrics-r0","password":"bar"}`,
		statusCode: http.StatusNotFound,
		response:   ".*unexpected permission for user .*",
	})
}

func (s *workerSuite) TestMetricsUsersAddIdempotent(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), usertesting.GenNewName(c, userCreator)).Return(coreuser.User{
		UUID: coreuser.UUID("deadbeef"),
	}, nil)
	s.accessService.EXPECT().AddUser(gomock.Any(), service.AddUserArg{
		Name:        s.metricsUserName,
		DisplayName: "juju-metrics-r0",
		Password:    ptr(auth.NewPassword("bar")),
		CreatorUUID: coreuser.UUID("deadbeef"),
		Permission: permission.AccessSpec{
			Target: s.controllerModelID,
			Access: permission.ReadAccess,
		},
	}).Return(coreuser.UUID("foobar"), nil, usererrors.UserAlreadyExists)
	s.accessService.EXPECT().GetUserByAuth(gomock.Any(), s.metricsUserName, auth.NewPassword("bar")).Return(coreuser.User{
		CreatorName: usertesting.GenNewName(c, userCreator),
	}, nil)
	s.accessService.EXPECT().ReadUserAccessLevelForTarget(gomock.Any(), s.metricsUserName, s.controllerModelID).Return(
		permission.ReadAccess, nil,
	)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodPost,
		endpoint:   "/metrics-users",
		body:       `{"username":"juju-metrics-r0","password":"bar"}`,
		statusCode: http.StatusOK, // succeed as a no-op
		response:   `.*created user \\\"juju-metrics-r0\\\".*`,
	})
}

func (s *workerSuite) TestMetricsUsersRemoveInvalidMethod(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodGet,
		endpoint:   "/metrics-users/foo",
		statusCode: http.StatusMethodNotAllowed,
		ignoreBody: true,
	})
}

func (s *workerSuite) TestMetricsUsersRemoveUsernameMissingPrefix(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodDelete,
		endpoint:   "/metrics-users/foo",
		statusCode: http.StatusBadRequest,
		response:   `.*username .* should have prefix \\\"juju-metrics-\\\".*`,
	})
}

func (s *workerSuite) TestMetricsUsersRemoveSuccess(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), s.metricsUserName).Return(coreuser.User{
		UUID:        coreuser.UUID("deadbeef"),
		CreatorName: usertesting.GenNewName(c, userCreator),
	}, nil)
	s.accessService.EXPECT().RemoveUser(gomock.Any(), s.metricsUserName).Return(nil)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodDelete,
		endpoint:   "/metrics-users/juju-metrics-r0",
		statusCode: http.StatusOK,
		response:   `.*deleted user \\\"juju-metrics-r0\\\".*`,
	})
}

func (s *workerSuite) TestMetricsUsersRemoveForbidden(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), s.metricsUserName).Return(coreuser.User{
		UUID:        coreuser.UUID("deadbeef"),
		Name:        s.metricsUserName,
		CreatorName: usertesting.GenNewName(c, "not-you"),
	}, nil)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodDelete,
		endpoint:   "/metrics-users/juju-metrics-r0",
		statusCode: http.StatusForbidden,
		response:   `.*cannot remove user \\\"juju-metrics-r0\\\" created by \\\"not-you\\\".*`,
	})
}

func (s *workerSuite) TestMetricsUsersRemoveNotFound(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.accessService.EXPECT().GetUserByName(gomock.Any(), s.metricsUserName).Return(coreuser.User{
		UUID:        coreuser.UUID("deadbeef"),
		Name:        s.metricsUserName,
		CreatorName: usertesting.GenNewName(c, "not-you"),
	}, usererrors.UserNotFound)

	s.runHandlerTest(c, handlerTest{
		method:     http.MethodDelete,
		endpoint:   "/metrics-users/juju-metrics-r0",
		statusCode: http.StatusOK, // succeed as a no-op
		response:   `.*deleted user \\\"juju-metrics-r0\\\".*`,
	})
}

func (s *workerSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.accessService = NewMockAccessService(ctrl)
	return ctrl
}

// Return an *http.Client with custom transport that allows it to connect to
// the given Unix socket.
func client(socketPath string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (conn net.Conn, err error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}
