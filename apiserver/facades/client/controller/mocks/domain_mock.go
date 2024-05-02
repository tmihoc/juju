// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/controller (interfaces: ControllerAccessService,ControllerConfigService)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/domain_mock.go github.com/juju/juju/apiserver/facades/client/controller ControllerAccessService,ControllerConfigService
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	controller "github.com/juju/juju/controller"
	model "github.com/juju/juju/core/model"
	permission "github.com/juju/juju/core/permission"
	access "github.com/juju/juju/domain/access"
	gomock "go.uber.org/mock/gomock"
)

// MockControllerAccessService is a mock of ControllerAccessService interface.
type MockControllerAccessService struct {
	ctrl     *gomock.Controller
	recorder *MockControllerAccessServiceMockRecorder
}

// MockControllerAccessServiceMockRecorder is the mock recorder for MockControllerAccessService.
type MockControllerAccessServiceMockRecorder struct {
	mock *MockControllerAccessService
}

// NewMockControllerAccessService creates a new mock instance.
func NewMockControllerAccessService(ctrl *gomock.Controller) *MockControllerAccessService {
	mock := &MockControllerAccessService{ctrl: ctrl}
	mock.recorder = &MockControllerAccessServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControllerAccessService) EXPECT() *MockControllerAccessServiceMockRecorder {
	return m.recorder
}

// LastModelLogin mocks base method.
func (m *MockControllerAccessService) LastModelLogin(arg0 context.Context, arg1 string, arg2 model.UUID) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastModelLogin", arg0, arg1, arg2)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastModelLogin indicates an expected call of LastModelLogin.
func (mr *MockControllerAccessServiceMockRecorder) LastModelLogin(arg0, arg1, arg2 any) *MockControllerAccessServiceLastModelLoginCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastModelLogin", reflect.TypeOf((*MockControllerAccessService)(nil).LastModelLogin), arg0, arg1, arg2)
	return &MockControllerAccessServiceLastModelLoginCall{Call: call}
}

// MockControllerAccessServiceLastModelLoginCall wrap *gomock.Call
type MockControllerAccessServiceLastModelLoginCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAccessServiceLastModelLoginCall) Return(arg0 time.Time, arg1 error) *MockControllerAccessServiceLastModelLoginCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAccessServiceLastModelLoginCall) Do(f func(context.Context, string, model.UUID) (time.Time, error)) *MockControllerAccessServiceLastModelLoginCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAccessServiceLastModelLoginCall) DoAndReturn(f func(context.Context, string, model.UUID) (time.Time, error)) *MockControllerAccessServiceLastModelLoginCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ReadUserAccessLevelForTarget mocks base method.
func (m *MockControllerAccessService) ReadUserAccessLevelForTarget(arg0 context.Context, arg1 string, arg2 permission.ID) (permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadUserAccessLevelForTarget", arg0, arg1, arg2)
	ret0, _ := ret[0].(permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadUserAccessLevelForTarget indicates an expected call of ReadUserAccessLevelForTarget.
func (mr *MockControllerAccessServiceMockRecorder) ReadUserAccessLevelForTarget(arg0, arg1, arg2 any) *MockControllerAccessServiceReadUserAccessLevelForTargetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadUserAccessLevelForTarget", reflect.TypeOf((*MockControllerAccessService)(nil).ReadUserAccessLevelForTarget), arg0, arg1, arg2)
	return &MockControllerAccessServiceReadUserAccessLevelForTargetCall{Call: call}
}

// MockControllerAccessServiceReadUserAccessLevelForTargetCall wrap *gomock.Call
type MockControllerAccessServiceReadUserAccessLevelForTargetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAccessServiceReadUserAccessLevelForTargetCall) Return(arg0 permission.Access, arg1 error) *MockControllerAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAccessServiceReadUserAccessLevelForTargetCall) Do(f func(context.Context, string, permission.ID) (permission.Access, error)) *MockControllerAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAccessServiceReadUserAccessLevelForTargetCall) DoAndReturn(f func(context.Context, string, permission.ID) (permission.Access, error)) *MockControllerAccessServiceReadUserAccessLevelForTargetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdatePermission mocks base method.
func (m *MockControllerAccessService) UpdatePermission(arg0 context.Context, arg1 access.UpdatePermissionArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePermission", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePermission indicates an expected call of UpdatePermission.
func (mr *MockControllerAccessServiceMockRecorder) UpdatePermission(arg0, arg1 any) *MockControllerAccessServiceUpdatePermissionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePermission", reflect.TypeOf((*MockControllerAccessService)(nil).UpdatePermission), arg0, arg1)
	return &MockControllerAccessServiceUpdatePermissionCall{Call: call}
}

// MockControllerAccessServiceUpdatePermissionCall wrap *gomock.Call
type MockControllerAccessServiceUpdatePermissionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerAccessServiceUpdatePermissionCall) Return(arg0 error) *MockControllerAccessServiceUpdatePermissionCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerAccessServiceUpdatePermissionCall) Do(f func(context.Context, access.UpdatePermissionArgs) error) *MockControllerAccessServiceUpdatePermissionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerAccessServiceUpdatePermissionCall) DoAndReturn(f func(context.Context, access.UpdatePermissionArgs) error) *MockControllerAccessServiceUpdatePermissionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockControllerConfigService is a mock of ControllerConfigService interface.
type MockControllerConfigService struct {
	ctrl     *gomock.Controller
	recorder *MockControllerConfigServiceMockRecorder
}

// MockControllerConfigServiceMockRecorder is the mock recorder for MockControllerConfigService.
type MockControllerConfigServiceMockRecorder struct {
	mock *MockControllerConfigService
}

// NewMockControllerConfigService creates a new mock instance.
func NewMockControllerConfigService(ctrl *gomock.Controller) *MockControllerConfigService {
	mock := &MockControllerConfigService{ctrl: ctrl}
	mock.recorder = &MockControllerConfigServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControllerConfigService) EXPECT() *MockControllerConfigServiceMockRecorder {
	return m.recorder
}

// ControllerConfig mocks base method.
func (m *MockControllerConfigService) ControllerConfig(arg0 context.Context) (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig", arg0)
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig.
func (mr *MockControllerConfigServiceMockRecorder) ControllerConfig(arg0 any) *MockControllerConfigServiceControllerConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockControllerConfigService)(nil).ControllerConfig), arg0)
	return &MockControllerConfigServiceControllerConfigCall{Call: call}
}

// MockControllerConfigServiceControllerConfigCall wrap *gomock.Call
type MockControllerConfigServiceControllerConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerConfigServiceControllerConfigCall) Return(arg0 controller.Config, arg1 error) *MockControllerConfigServiceControllerConfigCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerConfigServiceControllerConfigCall) Do(f func(context.Context) (controller.Config, error)) *MockControllerConfigServiceControllerConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerConfigServiceControllerConfigCall) DoAndReturn(f func(context.Context) (controller.Config, error)) *MockControllerConfigServiceControllerConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateControllerConfig mocks base method.
func (m *MockControllerConfigService) UpdateControllerConfig(arg0 context.Context, arg1 controller.Config, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateControllerConfig", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateControllerConfig indicates an expected call of UpdateControllerConfig.
func (mr *MockControllerConfigServiceMockRecorder) UpdateControllerConfig(arg0, arg1, arg2 any) *MockControllerConfigServiceUpdateControllerConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateControllerConfig", reflect.TypeOf((*MockControllerConfigService)(nil).UpdateControllerConfig), arg0, arg1, arg2)
	return &MockControllerConfigServiceUpdateControllerConfigCall{Call: call}
}

// MockControllerConfigServiceUpdateControllerConfigCall wrap *gomock.Call
type MockControllerConfigServiceUpdateControllerConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockControllerConfigServiceUpdateControllerConfigCall) Return(arg0 error) *MockControllerConfigServiceUpdateControllerConfigCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockControllerConfigServiceUpdateControllerConfigCall) Do(f func(context.Context, controller.Config, []string) error) *MockControllerConfigServiceUpdateControllerConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockControllerConfigServiceUpdateControllerConfigCall) DoAndReturn(f func(context.Context, controller.Config, []string) error) *MockControllerConfigServiceUpdateControllerConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
