// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/controller/caasapplicationprovisioner (interfaces: ControllerConfigService,ModelConfigService,ModelInfoService,ApplicationService)
//
// Generated by this command:
//
//	mockgen -package caasapplicationprovisioner_test -destination service_mock_test.go github.com/juju/juju/apiserver/facades/controller/caasapplicationprovisioner ControllerConfigService,ModelConfigService,ModelInfoService,ApplicationService
//

// Package caasapplicationprovisioner_test is a generated GoMock package.
package caasapplicationprovisioner_test

import (
	context "context"
	reflect "reflect"

	controller "github.com/juju/juju/controller"
	leadership "github.com/juju/juju/core/leadership"
	life "github.com/juju/juju/core/life"
	model "github.com/juju/juju/core/model"
	unit "github.com/juju/juju/core/unit"
	watcher "github.com/juju/juju/core/watcher"
	service "github.com/juju/juju/domain/application/service"
	config "github.com/juju/juju/environs/config"
	gomock "go.uber.org/mock/gomock"
)

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
func (mr *MockControllerConfigServiceMockRecorder) ControllerConfig(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockControllerConfigService)(nil).ControllerConfig), arg0)
}

// WatchControllerConfig mocks base method.
func (m *MockControllerConfigService) WatchControllerConfig() (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchControllerConfig")
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchControllerConfig indicates an expected call of WatchControllerConfig.
func (mr *MockControllerConfigServiceMockRecorder) WatchControllerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchControllerConfig", reflect.TypeOf((*MockControllerConfigService)(nil).WatchControllerConfig))
}

// MockModelConfigService is a mock of ModelConfigService interface.
type MockModelConfigService struct {
	ctrl     *gomock.Controller
	recorder *MockModelConfigServiceMockRecorder
}

// MockModelConfigServiceMockRecorder is the mock recorder for MockModelConfigService.
type MockModelConfigServiceMockRecorder struct {
	mock *MockModelConfigService
}

// NewMockModelConfigService creates a new mock instance.
func NewMockModelConfigService(ctrl *gomock.Controller) *MockModelConfigService {
	mock := &MockModelConfigService{ctrl: ctrl}
	mock.recorder = &MockModelConfigServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelConfigService) EXPECT() *MockModelConfigServiceMockRecorder {
	return m.recorder
}

// ModelConfig mocks base method.
func (m *MockModelConfigService) ModelConfig(arg0 context.Context) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig", arg0)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig.
func (mr *MockModelConfigServiceMockRecorder) ModelConfig(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockModelConfigService)(nil).ModelConfig), arg0)
}

// Watch mocks base method.
func (m *MockModelConfigService) Watch() (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch")
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockModelConfigServiceMockRecorder) Watch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockModelConfigService)(nil).Watch))
}

// MockModelInfoService is a mock of ModelInfoService interface.
type MockModelInfoService struct {
	ctrl     *gomock.Controller
	recorder *MockModelInfoServiceMockRecorder
}

// MockModelInfoServiceMockRecorder is the mock recorder for MockModelInfoService.
type MockModelInfoServiceMockRecorder struct {
	mock *MockModelInfoService
}

// NewMockModelInfoService creates a new mock instance.
func NewMockModelInfoService(ctrl *gomock.Controller) *MockModelInfoService {
	mock := &MockModelInfoService{ctrl: ctrl}
	mock.recorder = &MockModelInfoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelInfoService) EXPECT() *MockModelInfoServiceMockRecorder {
	return m.recorder
}

// GetModelInfo mocks base method.
func (m *MockModelInfoService) GetModelInfo(arg0 context.Context) (model.ReadOnlyModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModelInfo", arg0)
	ret0, _ := ret[0].(model.ReadOnlyModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetModelInfo indicates an expected call of GetModelInfo.
func (mr *MockModelInfoServiceMockRecorder) GetModelInfo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModelInfo", reflect.TypeOf((*MockModelInfoService)(nil).GetModelInfo), arg0)
}

// MockApplicationService is a mock of ApplicationService interface.
type MockApplicationService struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationServiceMockRecorder
}

// MockApplicationServiceMockRecorder is the mock recorder for MockApplicationService.
type MockApplicationServiceMockRecorder struct {
	mock *MockApplicationService
}

// NewMockApplicationService creates a new mock instance.
func NewMockApplicationService(ctrl *gomock.Controller) *MockApplicationService {
	mock := &MockApplicationService{ctrl: ctrl}
	mock.recorder = &MockApplicationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicationService) EXPECT() *MockApplicationServiceMockRecorder {
	return m.recorder
}

// DestroyUnit mocks base method.
func (m *MockApplicationService) DestroyUnit(arg0 context.Context, arg1 unit.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyUnit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyUnit indicates an expected call of DestroyUnit.
func (mr *MockApplicationServiceMockRecorder) DestroyUnit(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyUnit", reflect.TypeOf((*MockApplicationService)(nil).DestroyUnit), arg0, arg1)
}

// GetApplicationLife mocks base method.
func (m *MockApplicationService) GetApplicationLife(arg0 context.Context, arg1 string) (life.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationLife", arg0, arg1)
	ret0, _ := ret[0].(life.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationLife indicates an expected call of GetApplicationLife.
func (mr *MockApplicationServiceMockRecorder) GetApplicationLife(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationLife", reflect.TypeOf((*MockApplicationService)(nil).GetApplicationLife), arg0, arg1)
}

// GetApplicationScale mocks base method.
func (m *MockApplicationService) GetApplicationScale(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationScale", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationScale indicates an expected call of GetApplicationScale.
func (mr *MockApplicationServiceMockRecorder) GetApplicationScale(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationScale", reflect.TypeOf((*MockApplicationService)(nil).GetApplicationScale), arg0, arg1)
}

// GetApplicationScalingState mocks base method.
func (m *MockApplicationService) GetApplicationScalingState(arg0 context.Context, arg1 string) (service.ScalingState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationScalingState", arg0, arg1)
	ret0, _ := ret[0].(service.ScalingState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationScalingState indicates an expected call of GetApplicationScalingState.
func (mr *MockApplicationServiceMockRecorder) GetApplicationScalingState(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationScalingState", reflect.TypeOf((*MockApplicationService)(nil).GetApplicationScalingState), arg0, arg1)
}

// GetUnitLife mocks base method.
func (m *MockApplicationService) GetUnitLife(arg0 context.Context, arg1 unit.Name) (life.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitLife", arg0, arg1)
	ret0, _ := ret[0].(life.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitLife indicates an expected call of GetUnitLife.
func (mr *MockApplicationServiceMockRecorder) GetUnitLife(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitLife", reflect.TypeOf((*MockApplicationService)(nil).GetUnitLife), arg0, arg1)
}

// RemoveUnit mocks base method.
func (m *MockApplicationService) RemoveUnit(arg0 context.Context, arg1 unit.Name, arg2 leadership.Revoker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUnit", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUnit indicates an expected call of RemoveUnit.
func (mr *MockApplicationServiceMockRecorder) RemoveUnit(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUnit", reflect.TypeOf((*MockApplicationService)(nil).RemoveUnit), arg0, arg1, arg2)
}

// SetApplicationScalingState mocks base method.
func (m *MockApplicationService) SetApplicationScalingState(arg0 context.Context, arg1 string, arg2 int, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetApplicationScalingState", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetApplicationScalingState indicates an expected call of SetApplicationScalingState.
func (mr *MockApplicationServiceMockRecorder) SetApplicationScalingState(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetApplicationScalingState", reflect.TypeOf((*MockApplicationService)(nil).SetApplicationScalingState), arg0, arg1, arg2, arg3)
}

// UpdateCAASUnit mocks base method.
func (m *MockApplicationService) UpdateCAASUnit(arg0 context.Context, arg1 unit.Name, arg2 service.UpdateCAASUnitParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCAASUnit", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCAASUnit indicates an expected call of UpdateCAASUnit.
func (mr *MockApplicationServiceMockRecorder) UpdateCAASUnit(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCAASUnit", reflect.TypeOf((*MockApplicationService)(nil).UpdateCAASUnit), arg0, arg1, arg2)
}
