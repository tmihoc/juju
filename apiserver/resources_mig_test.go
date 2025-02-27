// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/apiserverhttp"
	"github.com/juju/juju/core/resource"
	"github.com/juju/juju/core/unit"
	domainresource "github.com/juju/juju/domain/resource"
	"github.com/juju/juju/rpc/params"
)

type resourcesUploadSuite struct {
	applicationsServiceGetter *MockApplicationServiceGetter
	applicationsService       *MockApplicationService
	resourceServiceGetter     *MockResourceServiceGetter
	resourceService           *MockResourceService

	mux *apiserverhttp.Mux
	srv *httptest.Server
}

var _ = gc.Suite(&resourcesUploadSuite{})

func (s *resourcesUploadSuite) SetUpTest(c *gc.C) {
	s.mux = apiserverhttp.NewMux()
	s.srv = httptest.NewServer(s.mux)
}

func (s *resourcesUploadSuite) TearDownTest(c *gc.C) {
	s.srv.Close()
}
func (s *resourcesUploadSuite) TestStub(c *gc.C) {
	c.Skip("This suite is missing tests for the following scenarios:\n" +
		"- Sending a POST req requires authorization via unit or application only.\n" +
		"- Rejects an unknown model with http.StatusNotFound.\n" +
		"- Test fails when model not importing.")
}

// TestServeMethodNotSupported ensures that the handler rejects HTTP methods
// other than POST with a 405 Method Not Allowed response.
func (s *resourcesUploadSuite) TestServeMethodNotSupported(c *gc.C) {
	// Arrange
	handler := &resourcesMigrationUploadHandler{}
	unsupportedMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	// This is a bit pathological, but we want to make sure that the handler
	// logic only actions on POST requests.
	for _, method := range unsupportedMethods {
		s.mux.AddHandler(method, migrateResourcesPrefix, handler)
		defer s.mux.RemoveHandler(method, migrateResourcesPrefix)
	}

	url := s.srv.URL + migrateResourcesPrefix
	for _, method := range unsupportedMethods {
		// Act
		request, err := http.NewRequest(method, url, nil)
		c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while building request. method: %s", method))
		response, err := http.DefaultClient.Do(request)
		c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request. method: %s", method))

		// Assert
		c.Check(response.StatusCode, gc.Equals, http.StatusMethodNotAllowed,
			gc.Commentf("(Assert) unexpected status code. method: %s", method))
	}
}

// TestGetUploadTargetBadRequest verifies the handling of bad requests in upload
// target retrieval when query parameters are wrong
func (s *resourcesUploadSuite) TestGetUploadTargetBadRequest(c *gc.C) {
	type testCase struct {
		name       string
		query      url.Values
		errMatches string
	}

	for _, tc := range []testCase{
		{
			name:       "missing name",
			query:      url.Values{},
			errMatches: "missing resource name",
		},
		{
			name:       "missing application and unit",
			query:      url.Values{"name": {"test"}},
			errMatches: "missing application/unit",
		},
		{
			name:       "both application and unit defined",
			query:      url.Values{"name": {"test"}, "application": {"testapp"}, "unit": {"testunit"}},
			errMatches: "application and unit can't be set at the same time",
		},
		{
			name:       "malformed unit name",
			query:      url.Values{"name": {"test"}, "unit": {"testunit"}},
			errMatches: "invalid unit name.*",
		},
	} {
		// Act
		_, err := getUploadTarget(context.Background(), s.applicationsService, tc.query)

		// Assert
		c.Check(err, jc.ErrorIs, errors.BadRequest, gc.Commentf("(Assert) unexpected error. test case: %s", tc.name))
		c.Check(err, gc.ErrorMatches, tc.errMatches, gc.Commentf("(Assert) errors doesn't match. test case: %s",
			tc.name))
	}
}

// TestGetUploadTargetCannotGetUnitUUID verifies that getUploadTarget returns
// an expected error when unit UUID retrieval fails.
func (s *resourcesUploadSuite) TestGetUploadTargetCannotGetUnitUUID(c *gc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	expectedErr := errors.New("cannot get unit uuid")
	s.applicationsService.EXPECT().GetUnitUUID(gomock.Any(), gomock.Any()).Return("", expectedErr)
	query := url.Values{"name": {"test"}, "unit": {"testunit/0"}}

	// Act
	_, err := getUploadTarget(context.Background(), s.applicationsService, query)

	// Assert
	c.Check(err, jc.ErrorIs, expectedErr, gc.Commentf("(Assert) unexpected error."))
}

// TestGetUploadTargetCannotGetApplicationByUnitName verifies that an error
// is returned when  while retrieving an application by unit name
func (s *resourcesUploadSuite) TestGetUploadTargetCannotGetApplicationByUnitName(c *gc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	expectedErr := errors.New("cannot get application by unit name")
	s.applicationsService.EXPECT().GetUnitUUID(gomock.Any(), gomock.Any()).Return("whatever", nil)
	s.applicationsService.EXPECT().GetApplicationIDByUnitName(gomock.Any(), gomock.Any()).Return("", expectedErr)
	query := url.Values{"name": {"test"}, "unit": {"testunit/0"}}

	// Act
	_, err := getUploadTarget(context.Background(), s.applicationsService, query)

	// Assert
	c.Check(err, jc.ErrorIs, expectedErr, gc.Commentf("(Assert) unexpected error."))
}

// TestGetUploadTargetCannotGetApplicationByName verifies behavior
// when GetApplicationIDByName fails while retrieving an application by name
func (s *resourcesUploadSuite) TestGetUploadTargetCannotGetApplicationByName(c *gc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	expectedErr := errors.New("cannot get application by name")
	s.applicationsService.EXPECT().GetApplicationIDByName(gomock.Any(), gomock.Any()).Return("", expectedErr)
	query := url.Values{"name": {"test"}, "application": {"testapplication"}}

	// Act
	_, err := getUploadTarget(context.Background(), s.applicationsService, query)

	// Assert
	c.Check(err, jc.ErrorIs, expectedErr, gc.Commentf("(Assert) unexpected error."))
}

// TestGetUploadTargetByApplication tests retrieving the upload target based
// on application name.
func (s *resourcesUploadSuite) TestGetUploadTargetByApplication(c *gc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	s.applicationsService.EXPECT().GetApplicationIDByName(gomock.Any(), "testapplication").Return("testapp-id", nil)
	query := url.Values{"name": {"test"}, "application": {"testapplication"}}

	// Act
	target, err := getUploadTarget(context.Background(), s.applicationsService, query)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error."))

	// Assert
	c.Check(target, gc.Equals, resourceUploadTarget{
		name:  "test",
		appID: "testapp-id",
	}, gc.Commentf("(Assert) unexpected result."))
}

// TestGetUploadTargetByUnittests retrieving the upload target based
// on unit name.
func (s *resourcesUploadSuite) TestGetUploadTargetByUnit(c *gc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	s.applicationsService.EXPECT().GetUnitUUID(gomock.Any(), unit.Name("testunit/0")).Return("testunit-id", nil)
	s.applicationsService.EXPECT().GetApplicationIDByUnitName(gomock.Any(), unit.Name("testunit/0")).Return("testapp-id", nil)
	query := url.Values{"name": {"test"}, "unit": {"testunit/0"}}

	// Act
	target, err := getUploadTarget(context.Background(), s.applicationsService, query)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error."))

	// Assert
	c.Check(target, gc.Equals, resourceUploadTarget{
		name:     "test",
		appID:    "testapp-id",
		unitUUID: "testunit-id",
	}, gc.Commentf("(Assert) unexpected result."))
}

// TestServeUploadApplicationResourceNotFound verifies the handler's behavior
// when the application resource is not found.
func (s *resourcesUploadSuite) TestServeUploadApplicationResourceNotFound(c *gc.C) {
	// Arrange
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:  "test",
		appID: "testapp-id",
	})
	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), domainresource.GetApplicationResourceIDArgs{
		ApplicationID: "testapp-id",
		Name:          "test",
	}).Return("", errors.NotFound)

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", nil)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	c.Check(response.StatusCode, gc.Equals, http.StatusNotFound,
		gc.Commentf("(Assert) unexpected status code."))
}

// TestServeUploadApplicationStoreResourceError verifies error handling
// when storing a resource during an upload operation.
func (s *resourcesUploadSuite) TestServeUploadApplicationStoreResourceError(c *gc.C) {
	// Arrange
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:  "test",
		appID: "testapp-id",
	})
	query.Add("timestamp", "not-placeholder")
	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), gomock.Any()).Return("res-uuid", nil)
	s.resourceService.EXPECT().StoreResource(gomock.Any(), gomock.Any()).Return(errors.New("cannot store resource"))

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", nil)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	c.Check(response.StatusCode, gc.Equals, http.StatusInternalServerError,
		gc.Commentf("(Assert) unexpected status code."))
}

// TestServeUploadApplicationGetResourceError validates the behavior
// when an error occurs while retrieving a resource.
func (s *resourcesUploadSuite) TestServeUploadApplicationGetResourceError(c *gc.C) {
	// Arrange
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:  "test",
		appID: "testapp-id",
	})
	query.Add("timestamp", "not-placeholder")
	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), gomock.Any()).Return("res-uuid", nil)
	s.resourceService.EXPECT().StoreResource(gomock.Any(), gomock.Any()).Return(nil)
	s.resourceService.EXPECT().GetResource(gomock.Any(), gomock.Any()).Return(domainresource.Resource{}, errors.New(
		"cannot get resource"))

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", nil)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	c.Check(response.StatusCode, gc.Equals, http.StatusInternalServerError,
		gc.Commentf("(Assert) unexpected status code."))
}

// TestServeUploadApplicationWithPlaceholder tests the application's ability to
// handle uploading with placeholders correctly. It verifies that StoreResource is
// not called through not configuring related mock.
func (s *resourcesUploadSuite) TestServeUploadApplicationWithPlaceholder(c *gc.C) {
	// Arrange
	now := time.Now().Truncate(time.Second).UTC()
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:  "test",
		appID: "testapp-id",
	})

	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), gomock.Any()).Return("res-uuid", nil)
	s.resourceService.EXPECT().GetResource(gomock.Any(), gomock.Any()).Return(domainresource.Resource{
		UUID:      "res-uuid",
		Timestamp: now,
	}, nil)

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", http.NoBody)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	var obtained params.ResourceUploadResult
	c.Check(response.StatusCode, gc.Equals, http.StatusOK,
		gc.Commentf("(Assert) unexpected status code."))
	body, err := io.ReadAll(response.Body)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Assert) unexpected error while reading response body"))
	c.Assert(json.Unmarshal(body, &obtained), jc.ErrorIsNil,
		gc.Commentf("(Assert) unexpected error while unmarshalling response"))
	c.Check(obtained, gc.Equals, params.ResourceUploadResult{
		ID:        "res-uuid",
		Timestamp: now,
	})
}

// TestServeUploadApplication tests the HTTP endpoint for uploading application
// resources, ensuring correct handling and storage. It verifies correct values
// are passed as argument for calling the different underlying services.
func (s *resourcesUploadSuite) TestServeUploadApplication(c *gc.C) {
	// Arrange
	now := time.Now().Truncate(time.Second).UTC()
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:  "test",
		appID: "testapp-id",
	})
	query.Add("timestamp", "not-placeholder")

	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), domainresource.GetApplicationResourceIDArgs{
		ApplicationID: "testapp-id",
		Name:          "test",
	}).Return("res-uuid", nil)
	s.resourceService.EXPECT().StoreResource(gomock.Any(), domainresource.StoreResourceArgs{
		ResourceUUID:    "res-uuid",
		Reader:          http.NoBody,
		RetrievedBy:     "testapp-id",
		RetrievedByType: domainresource.Application,
	}).Return(nil)
	s.resourceService.EXPECT().GetResource(gomock.Any(), resource.UUID("res-uuid")).Return(domainresource.Resource{
		UUID:      "res-uuid",
		Timestamp: now,
	}, nil)

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", http.NoBody)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	var obtained params.ResourceUploadResult
	c.Check(response.StatusCode, gc.Equals, http.StatusOK,
		gc.Commentf("(Assert) unexpected status code."))
	body, err := io.ReadAll(response.Body)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Assert) unexpected error while reading response body"))
	c.Assert(json.Unmarshal(body, &obtained), jc.ErrorIsNil,
		gc.Commentf("(Assert) unexpected error while unmarshalling response"))
	c.Check(obtained, gc.Equals, params.ResourceUploadResult{
		ID:        "res-uuid",
		Timestamp: now,
	})
}

// TestServeUploadUnitWithPlaceholder tests the upload functionality for a unit
// with a placeholder in the resource upload service. It is basically the same
// test than the one with application, with one call to SetUnitResource.
func (s *resourcesUploadSuite) TestServeUploadUnitWithPlaceholder(c *gc.C) {
	// Arrange
	now := time.Now().Truncate(time.Second).UTC()
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:     "test",
		appID:    "testapp-id",
		unitUUID: "testunit-id",
	})

	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), gomock.Any()).Return("res-uuid", nil)
	s.resourceService.EXPECT().SetUnitResource(gomock.Any(), resource.UUID("res-uuid"),
		unit.UUID("testunit-id")).Return(nil)
	s.resourceService.EXPECT().GetResource(gomock.Any(), gomock.Any()).Return(domainresource.Resource{
		UUID:      "res-uuid",
		Timestamp: now,
	}, nil)

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", http.NoBody)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	var obtained params.ResourceUploadResult
	c.Check(response.StatusCode, gc.Equals, http.StatusOK,
		gc.Commentf("(Assert) unexpected status code."))
	body, err := io.ReadAll(response.Body)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Assert) unexpected error while reading response body"))
	c.Assert(json.Unmarshal(body, &obtained), jc.ErrorIsNil,
		gc.Commentf("(Assert) unexpected error while unmarshalling response"))
	c.Check(obtained, gc.Equals, params.ResourceUploadResult{
		ID:        "res-uuid",
		Timestamp: now,
	})
}

// TestServeUploadUnit tests the process of uploading a resource unit and
// verifies correct resource handling and response. It is basically the same
// test than the one with application, with one call to SetUnitResource.
func (s *resourcesUploadSuite) TestServeUploadUnit(c *gc.C) {
	// Arrange
	now := time.Now().Truncate(time.Second).UTC()
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:     "test",
		appID:    "testapp-id",
		unitUUID: "testunit-id",
	})
	query.Add("timestamp", "not-placeholder")

	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), domainresource.GetApplicationResourceIDArgs{
		ApplicationID: "testapp-id",
		Name:          "test",
	}).Return("res-uuid", nil)
	s.resourceService.EXPECT().SetUnitResource(gomock.Any(), resource.UUID("res-uuid"),
		unit.UUID("testunit-id")).Return(nil)
	s.resourceService.EXPECT().StoreResource(gomock.Any(), domainresource.StoreResourceArgs{
		ResourceUUID:    "res-uuid",
		Reader:          http.NoBody,
		RetrievedBy:     "testunit-id",
		RetrievedByType: domainresource.Unit,
	}).Return(nil)
	s.resourceService.EXPECT().GetResource(gomock.Any(), resource.UUID("res-uuid")).Return(domainresource.Resource{
		UUID:      "res-uuid",
		Timestamp: now,
	}, nil)

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", http.NoBody)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	var obtained params.ResourceUploadResult
	c.Check(response.StatusCode, gc.Equals, http.StatusOK,
		gc.Commentf("(Assert) unexpected status code."))
	body, err := io.ReadAll(response.Body)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Assert) unexpected error while reading response body"))
	c.Assert(json.Unmarshal(body, &obtained), jc.ErrorIsNil,
		gc.Commentf("(Assert) unexpected error while unmarshalling response"))
	c.Check(obtained, gc.Equals, params.ResourceUploadResult{
		ID:        "res-uuid",
		Timestamp: now,
	})
}

// TestServeUploadUnitErrorSetUnitResource tests the error handling when failing
// to set a unit resource during upload processing.
func (s *resourcesUploadSuite) TestServeUploadUnitErrorSetUnitResource(c *gc.C) {
	// Arrange
	defer s.setupHandler(c).Finish()
	query := s.setupUploadTarget(c, resourceUploadTarget{
		name:     "test",
		appID:    "testapp-id",
		unitUUID: "testunit-id",
	})
	query.Add("timestamp", "not-placeholder")

	s.resourceService.EXPECT().GetApplicationResourceID(gomock.Any(), domainresource.GetApplicationResourceIDArgs{
		ApplicationID: "testapp-id",
		Name:          "test",
	}).Return("res-uuid", nil)
	s.resourceService.EXPECT().SetUnitResource(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("cannot set unit resource"))

	// Act
	response, err := http.Post(s.srv.URL+migrateResourcesPrefix+"?"+query.Encode(), "application/octet-stream", http.NoBody)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Act) unexpected error while executing request"))

	// Assert
	c.Check(response.StatusCode, gc.Equals, http.StatusInternalServerError,
		gc.Commentf("(Assert) unexpected status code."))
}

// setupUploadTarget configures the upload query parameters for a given resource
// upload target and sets up necessary mock expectations.
func (s *resourcesUploadSuite) setupUploadTarget(c *gc.C, target resourceUploadTarget) (query url.Values) {
	query = url.Values{"name": {target.name}}

	// Get app & unit uuids from unit name
	if target.unitUUID != "" {
		query.Add("unit", "testunit/0")
		s.applicationsService.EXPECT().GetUnitUUID(gomock.Any(), gomock.Any()).Return(target.unitUUID, nil)
		s.applicationsService.EXPECT().GetApplicationIDByUnitName(gomock.Any(), gomock.Any()).Return(target.appID, nil)
		return
	}

	// get app uuid from app name
	query.Add("application", "testapplication")
	s.applicationsService.EXPECT().GetApplicationIDByName(gomock.Any(), gomock.Any()).Return(target.appID, nil)
	return query
}

// setupHandler configures the resources migration upload HTTP handler, init
// mocks and registers it to the mux. It provides cleanup logic.
func (s *resourcesUploadSuite) setupHandler(c *gc.C) Finisher {
	finish := s.setupMocks(c).Finish
	s.expectApplicationService()
	s.expectResourceService()

	handler := &resourcesMigrationUploadHandler{
		resourceServiceGetter:    s.resourceServiceGetter,
		applicationServiceGetter: s.applicationsServiceGetter,
	}

	err := s.mux.AddHandler("POST", migrateResourcesPrefix, handler)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("(Arrange) unexpected error while adding handler"))

	return &finisherWrapper{
		finish: func() {
			s.mux.RemoveHandler("POST", migrateResourcesPrefix)
			finish()
		},
	}
}

// expectApplicationService prepare mocks for application service
func (s *resourcesUploadSuite) expectApplicationService() {
	s.applicationsServiceGetter.EXPECT().Application(gomock.Any()).Return(s.applicationsService, nil)
}

// expectResourceService prepare mocks for resource service
func (s *resourcesUploadSuite) expectResourceService() {
	s.resourceServiceGetter.EXPECT().Resource(gomock.Any()).Return(s.resourceService, nil)
}

// setupMocks initializes mock services and returns a gomock.Controller
// for managing mock lifecycle.
func (s *resourcesUploadSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.applicationsServiceGetter = NewMockApplicationServiceGetter(ctrl)
	s.applicationsService = NewMockApplicationService(ctrl)
	s.resourceServiceGetter = NewMockResourceServiceGetter(ctrl)
	s.resourceService = NewMockResourceService(ctrl)

	return ctrl
}

type Finisher interface {
	Finish()
}

type finisherWrapper struct {
	finish func()
}

func (f *finisherWrapper) Finish() {
	f.finish()
}
