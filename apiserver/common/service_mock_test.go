// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common (interfaces: ModelAgentService,MachineRebootService,EnsureDeadMachineService,WatchableMachineService,UnitStateService,MachineService)
//
// Generated by this command:
//
//	mockgen -typed -package common_test -destination service_mock_test.go github.com/juju/juju/apiserver/common ModelAgentService,MachineRebootService,EnsureDeadMachineService,WatchableMachineService,UnitStateService,MachineService
//

// Package common_test is a generated GoMock package.
package common_test

import (
	context "context"
	reflect "reflect"

	machine "github.com/juju/juju/core/machine"
	watcher "github.com/juju/juju/core/watcher"
	unitstate "github.com/juju/juju/domain/unitstate"
	version "github.com/juju/version/v2"
	gomock "go.uber.org/mock/gomock"
)

// MockModelAgentService is a mock of ModelAgentService interface.
type MockModelAgentService struct {
	ctrl     *gomock.Controller
	recorder *MockModelAgentServiceMockRecorder
}

// MockModelAgentServiceMockRecorder is the mock recorder for MockModelAgentService.
type MockModelAgentServiceMockRecorder struct {
	mock *MockModelAgentService
}

// NewMockModelAgentService creates a new mock instance.
func NewMockModelAgentService(ctrl *gomock.Controller) *MockModelAgentService {
	mock := &MockModelAgentService{ctrl: ctrl}
	mock.recorder = &MockModelAgentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelAgentService) EXPECT() *MockModelAgentServiceMockRecorder {
	return m.recorder
}

// GetModelAgentVersion mocks base method.
func (m *MockModelAgentService) GetModelAgentVersion(arg0 context.Context) (version.Number, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModelAgentVersion", arg0)
	ret0, _ := ret[0].(version.Number)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetModelAgentVersion indicates an expected call of GetModelAgentVersion.
func (mr *MockModelAgentServiceMockRecorder) GetModelAgentVersion(arg0 any) *MockModelAgentServiceGetModelAgentVersionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModelAgentVersion", reflect.TypeOf((*MockModelAgentService)(nil).GetModelAgentVersion), arg0)
	return &MockModelAgentServiceGetModelAgentVersionCall{Call: call}
}

// MockModelAgentServiceGetModelAgentVersionCall wrap *gomock.Call
type MockModelAgentServiceGetModelAgentVersionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelAgentServiceGetModelAgentVersionCall) Return(arg0 version.Number, arg1 error) *MockModelAgentServiceGetModelAgentVersionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelAgentServiceGetModelAgentVersionCall) Do(f func(context.Context) (version.Number, error)) *MockModelAgentServiceGetModelAgentVersionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelAgentServiceGetModelAgentVersionCall) DoAndReturn(f func(context.Context) (version.Number, error)) *MockModelAgentServiceGetModelAgentVersionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockMachineRebootService is a mock of MachineRebootService interface.
type MockMachineRebootService struct {
	ctrl     *gomock.Controller
	recorder *MockMachineRebootServiceMockRecorder
}

// MockMachineRebootServiceMockRecorder is the mock recorder for MockMachineRebootService.
type MockMachineRebootServiceMockRecorder struct {
	mock *MockMachineRebootService
}

// NewMockMachineRebootService creates a new mock instance.
func NewMockMachineRebootService(ctrl *gomock.Controller) *MockMachineRebootService {
	mock := &MockMachineRebootService{ctrl: ctrl}
	mock.recorder = &MockMachineRebootServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachineRebootService) EXPECT() *MockMachineRebootServiceMockRecorder {
	return m.recorder
}

// ClearMachineReboot mocks base method.
func (m *MockMachineRebootService) ClearMachineReboot(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearMachineReboot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearMachineReboot indicates an expected call of ClearMachineReboot.
func (mr *MockMachineRebootServiceMockRecorder) ClearMachineReboot(arg0, arg1 any) *MockMachineRebootServiceClearMachineRebootCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearMachineReboot", reflect.TypeOf((*MockMachineRebootService)(nil).ClearMachineReboot), arg0, arg1)
	return &MockMachineRebootServiceClearMachineRebootCall{Call: call}
}

// MockMachineRebootServiceClearMachineRebootCall wrap *gomock.Call
type MockMachineRebootServiceClearMachineRebootCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineRebootServiceClearMachineRebootCall) Return(arg0 error) *MockMachineRebootServiceClearMachineRebootCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineRebootServiceClearMachineRebootCall) Do(f func(context.Context, string) error) *MockMachineRebootServiceClearMachineRebootCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineRebootServiceClearMachineRebootCall) DoAndReturn(f func(context.Context, string) error) *MockMachineRebootServiceClearMachineRebootCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetMachineUUID mocks base method.
func (m *MockMachineRebootService) GetMachineUUID(arg0 context.Context, arg1 machine.Name) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMachineUUID", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMachineUUID indicates an expected call of GetMachineUUID.
func (mr *MockMachineRebootServiceMockRecorder) GetMachineUUID(arg0, arg1 any) *MockMachineRebootServiceGetMachineUUIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMachineUUID", reflect.TypeOf((*MockMachineRebootService)(nil).GetMachineUUID), arg0, arg1)
	return &MockMachineRebootServiceGetMachineUUIDCall{Call: call}
}

// MockMachineRebootServiceGetMachineUUIDCall wrap *gomock.Call
type MockMachineRebootServiceGetMachineUUIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineRebootServiceGetMachineUUIDCall) Return(arg0 string, arg1 error) *MockMachineRebootServiceGetMachineUUIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineRebootServiceGetMachineUUIDCall) Do(f func(context.Context, machine.Name) (string, error)) *MockMachineRebootServiceGetMachineUUIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineRebootServiceGetMachineUUIDCall) DoAndReturn(f func(context.Context, machine.Name) (string, error)) *MockMachineRebootServiceGetMachineUUIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsMachineRebootRequired mocks base method.
func (m *MockMachineRebootService) IsMachineRebootRequired(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMachineRebootRequired", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMachineRebootRequired indicates an expected call of IsMachineRebootRequired.
func (mr *MockMachineRebootServiceMockRecorder) IsMachineRebootRequired(arg0, arg1 any) *MockMachineRebootServiceIsMachineRebootRequiredCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMachineRebootRequired", reflect.TypeOf((*MockMachineRebootService)(nil).IsMachineRebootRequired), arg0, arg1)
	return &MockMachineRebootServiceIsMachineRebootRequiredCall{Call: call}
}

// MockMachineRebootServiceIsMachineRebootRequiredCall wrap *gomock.Call
type MockMachineRebootServiceIsMachineRebootRequiredCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineRebootServiceIsMachineRebootRequiredCall) Return(arg0 bool, arg1 error) *MockMachineRebootServiceIsMachineRebootRequiredCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineRebootServiceIsMachineRebootRequiredCall) Do(f func(context.Context, string) (bool, error)) *MockMachineRebootServiceIsMachineRebootRequiredCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineRebootServiceIsMachineRebootRequiredCall) DoAndReturn(f func(context.Context, string) (bool, error)) *MockMachineRebootServiceIsMachineRebootRequiredCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RequireMachineReboot mocks base method.
func (m *MockMachineRebootService) RequireMachineReboot(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequireMachineReboot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RequireMachineReboot indicates an expected call of RequireMachineReboot.
func (mr *MockMachineRebootServiceMockRecorder) RequireMachineReboot(arg0, arg1 any) *MockMachineRebootServiceRequireMachineRebootCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequireMachineReboot", reflect.TypeOf((*MockMachineRebootService)(nil).RequireMachineReboot), arg0, arg1)
	return &MockMachineRebootServiceRequireMachineRebootCall{Call: call}
}

// MockMachineRebootServiceRequireMachineRebootCall wrap *gomock.Call
type MockMachineRebootServiceRequireMachineRebootCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineRebootServiceRequireMachineRebootCall) Return(arg0 error) *MockMachineRebootServiceRequireMachineRebootCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineRebootServiceRequireMachineRebootCall) Do(f func(context.Context, string) error) *MockMachineRebootServiceRequireMachineRebootCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineRebootServiceRequireMachineRebootCall) DoAndReturn(f func(context.Context, string) error) *MockMachineRebootServiceRequireMachineRebootCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ShouldRebootOrShutdown mocks base method.
func (m *MockMachineRebootService) ShouldRebootOrShutdown(arg0 context.Context, arg1 string) (machine.RebootAction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldRebootOrShutdown", arg0, arg1)
	ret0, _ := ret[0].(machine.RebootAction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShouldRebootOrShutdown indicates an expected call of ShouldRebootOrShutdown.
func (mr *MockMachineRebootServiceMockRecorder) ShouldRebootOrShutdown(arg0, arg1 any) *MockMachineRebootServiceShouldRebootOrShutdownCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldRebootOrShutdown", reflect.TypeOf((*MockMachineRebootService)(nil).ShouldRebootOrShutdown), arg0, arg1)
	return &MockMachineRebootServiceShouldRebootOrShutdownCall{Call: call}
}

// MockMachineRebootServiceShouldRebootOrShutdownCall wrap *gomock.Call
type MockMachineRebootServiceShouldRebootOrShutdownCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineRebootServiceShouldRebootOrShutdownCall) Return(arg0 machine.RebootAction, arg1 error) *MockMachineRebootServiceShouldRebootOrShutdownCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineRebootServiceShouldRebootOrShutdownCall) Do(f func(context.Context, string) (machine.RebootAction, error)) *MockMachineRebootServiceShouldRebootOrShutdownCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineRebootServiceShouldRebootOrShutdownCall) DoAndReturn(f func(context.Context, string) (machine.RebootAction, error)) *MockMachineRebootServiceShouldRebootOrShutdownCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockEnsureDeadMachineService is a mock of EnsureDeadMachineService interface.
type MockEnsureDeadMachineService struct {
	ctrl     *gomock.Controller
	recorder *MockEnsureDeadMachineServiceMockRecorder
}

// MockEnsureDeadMachineServiceMockRecorder is the mock recorder for MockEnsureDeadMachineService.
type MockEnsureDeadMachineServiceMockRecorder struct {
	mock *MockEnsureDeadMachineService
}

// NewMockEnsureDeadMachineService creates a new mock instance.
func NewMockEnsureDeadMachineService(ctrl *gomock.Controller) *MockEnsureDeadMachineService {
	mock := &MockEnsureDeadMachineService{ctrl: ctrl}
	mock.recorder = &MockEnsureDeadMachineServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEnsureDeadMachineService) EXPECT() *MockEnsureDeadMachineServiceMockRecorder {
	return m.recorder
}

// EnsureDeadMachine mocks base method.
func (m *MockEnsureDeadMachineService) EnsureDeadMachine(arg0 context.Context, arg1 machine.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsureDeadMachine", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureDeadMachine indicates an expected call of EnsureDeadMachine.
func (mr *MockEnsureDeadMachineServiceMockRecorder) EnsureDeadMachine(arg0, arg1 any) *MockEnsureDeadMachineServiceEnsureDeadMachineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureDeadMachine", reflect.TypeOf((*MockEnsureDeadMachineService)(nil).EnsureDeadMachine), arg0, arg1)
	return &MockEnsureDeadMachineServiceEnsureDeadMachineCall{Call: call}
}

// MockEnsureDeadMachineServiceEnsureDeadMachineCall wrap *gomock.Call
type MockEnsureDeadMachineServiceEnsureDeadMachineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockEnsureDeadMachineServiceEnsureDeadMachineCall) Return(arg0 error) *MockEnsureDeadMachineServiceEnsureDeadMachineCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockEnsureDeadMachineServiceEnsureDeadMachineCall) Do(f func(context.Context, machine.Name) error) *MockEnsureDeadMachineServiceEnsureDeadMachineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockEnsureDeadMachineServiceEnsureDeadMachineCall) DoAndReturn(f func(context.Context, machine.Name) error) *MockEnsureDeadMachineServiceEnsureDeadMachineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockWatchableMachineService is a mock of WatchableMachineService interface.
type MockWatchableMachineService struct {
	ctrl     *gomock.Controller
	recorder *MockWatchableMachineServiceMockRecorder
}

// MockWatchableMachineServiceMockRecorder is the mock recorder for MockWatchableMachineService.
type MockWatchableMachineServiceMockRecorder struct {
	mock *MockWatchableMachineService
}

// NewMockWatchableMachineService creates a new mock instance.
func NewMockWatchableMachineService(ctrl *gomock.Controller) *MockWatchableMachineService {
	mock := &MockWatchableMachineService{ctrl: ctrl}
	mock.recorder = &MockWatchableMachineServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatchableMachineService) EXPECT() *MockWatchableMachineServiceMockRecorder {
	return m.recorder
}

// WatchMachineReboot mocks base method.
func (m *MockWatchableMachineService) WatchMachineReboot(arg0 context.Context, arg1 string) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchMachineReboot", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchMachineReboot indicates an expected call of WatchMachineReboot.
func (mr *MockWatchableMachineServiceMockRecorder) WatchMachineReboot(arg0, arg1 any) *MockWatchableMachineServiceWatchMachineRebootCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchMachineReboot", reflect.TypeOf((*MockWatchableMachineService)(nil).WatchMachineReboot), arg0, arg1)
	return &MockWatchableMachineServiceWatchMachineRebootCall{Call: call}
}

// MockWatchableMachineServiceWatchMachineRebootCall wrap *gomock.Call
type MockWatchableMachineServiceWatchMachineRebootCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockWatchableMachineServiceWatchMachineRebootCall) Return(arg0 watcher.Watcher[struct{}], arg1 error) *MockWatchableMachineServiceWatchMachineRebootCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockWatchableMachineServiceWatchMachineRebootCall) Do(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockWatchableMachineServiceWatchMachineRebootCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockWatchableMachineServiceWatchMachineRebootCall) DoAndReturn(f func(context.Context, string) (watcher.Watcher[struct{}], error)) *MockWatchableMachineServiceWatchMachineRebootCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockUnitStateService is a mock of UnitStateService interface.
type MockUnitStateService struct {
	ctrl     *gomock.Controller
	recorder *MockUnitStateServiceMockRecorder
}

// MockUnitStateServiceMockRecorder is the mock recorder for MockUnitStateService.
type MockUnitStateServiceMockRecorder struct {
	mock *MockUnitStateService
}

// NewMockUnitStateService creates a new mock instance.
func NewMockUnitStateService(ctrl *gomock.Controller) *MockUnitStateService {
	mock := &MockUnitStateService{ctrl: ctrl}
	mock.recorder = &MockUnitStateServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnitStateService) EXPECT() *MockUnitStateServiceMockRecorder {
	return m.recorder
}

// GetState mocks base method.
func (m *MockUnitStateService) GetState(arg0 context.Context, arg1 string) (unitstate.RetrievedUnitState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", arg0, arg1)
	ret0, _ := ret[0].(unitstate.RetrievedUnitState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockUnitStateServiceMockRecorder) GetState(arg0, arg1 any) *MockUnitStateServiceGetStateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockUnitStateService)(nil).GetState), arg0, arg1)
	return &MockUnitStateServiceGetStateCall{Call: call}
}

// MockUnitStateServiceGetStateCall wrap *gomock.Call
type MockUnitStateServiceGetStateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUnitStateServiceGetStateCall) Return(arg0 unitstate.RetrievedUnitState, arg1 error) *MockUnitStateServiceGetStateCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUnitStateServiceGetStateCall) Do(f func(context.Context, string) (unitstate.RetrievedUnitState, error)) *MockUnitStateServiceGetStateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUnitStateServiceGetStateCall) DoAndReturn(f func(context.Context, string) (unitstate.RetrievedUnitState, error)) *MockUnitStateServiceGetStateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetUnitUUIDForName mocks base method.
func (m *MockUnitStateService) GetUnitUUIDForName(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnitUUIDForName", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnitUUIDForName indicates an expected call of GetUnitUUIDForName.
func (mr *MockUnitStateServiceMockRecorder) GetUnitUUIDForName(arg0, arg1 any) *MockUnitStateServiceGetUnitUUIDForNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnitUUIDForName", reflect.TypeOf((*MockUnitStateService)(nil).GetUnitUUIDForName), arg0, arg1)
	return &MockUnitStateServiceGetUnitUUIDForNameCall{Call: call}
}

// MockUnitStateServiceGetUnitUUIDForNameCall wrap *gomock.Call
type MockUnitStateServiceGetUnitUUIDForNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUnitStateServiceGetUnitUUIDForNameCall) Return(arg0 string, arg1 error) *MockUnitStateServiceGetUnitUUIDForNameCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUnitStateServiceGetUnitUUIDForNameCall) Do(f func(context.Context, string) (string, error)) *MockUnitStateServiceGetUnitUUIDForNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUnitStateServiceGetUnitUUIDForNameCall) DoAndReturn(f func(context.Context, string) (string, error)) *MockUnitStateServiceGetUnitUUIDForNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetState mocks base method.
func (m *MockUnitStateService) SetState(arg0 context.Context, arg1 unitstate.UnitState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetState", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetState indicates an expected call of SetState.
func (mr *MockUnitStateServiceMockRecorder) SetState(arg0, arg1 any) *MockUnitStateServiceSetStateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetState", reflect.TypeOf((*MockUnitStateService)(nil).SetState), arg0, arg1)
	return &MockUnitStateServiceSetStateCall{Call: call}
}

// MockUnitStateServiceSetStateCall wrap *gomock.Call
type MockUnitStateServiceSetStateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUnitStateServiceSetStateCall) Return(arg0 error) *MockUnitStateServiceSetStateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUnitStateServiceSetStateCall) Do(f func(context.Context, unitstate.UnitState) error) *MockUnitStateServiceSetStateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUnitStateServiceSetStateCall) DoAndReturn(f func(context.Context, unitstate.UnitState) error) *MockUnitStateServiceSetStateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockMachineService is a mock of MachineService interface.
type MockMachineService struct {
	ctrl     *gomock.Controller
	recorder *MockMachineServiceMockRecorder
}

// MockMachineServiceMockRecorder is the mock recorder for MockMachineService.
type MockMachineServiceMockRecorder struct {
	mock *MockMachineService
}

// NewMockMachineService creates a new mock instance.
func NewMockMachineService(ctrl *gomock.Controller) *MockMachineService {
	mock := &MockMachineService{ctrl: ctrl}
	mock.recorder = &MockMachineServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachineService) EXPECT() *MockMachineServiceMockRecorder {
	return m.recorder
}

// EnsureDeadMachine mocks base method.
func (m *MockMachineService) EnsureDeadMachine(arg0 context.Context, arg1 machine.Name) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsureDeadMachine", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureDeadMachine indicates an expected call of EnsureDeadMachine.
func (mr *MockMachineServiceMockRecorder) EnsureDeadMachine(arg0, arg1 any) *MockMachineServiceEnsureDeadMachineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureDeadMachine", reflect.TypeOf((*MockMachineService)(nil).EnsureDeadMachine), arg0, arg1)
	return &MockMachineServiceEnsureDeadMachineCall{Call: call}
}

// MockMachineServiceEnsureDeadMachineCall wrap *gomock.Call
type MockMachineServiceEnsureDeadMachineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineServiceEnsureDeadMachineCall) Return(arg0 error) *MockMachineServiceEnsureDeadMachineCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineServiceEnsureDeadMachineCall) Do(f func(context.Context, machine.Name) error) *MockMachineServiceEnsureDeadMachineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineServiceEnsureDeadMachineCall) DoAndReturn(f func(context.Context, machine.Name) error) *MockMachineServiceEnsureDeadMachineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetMachineUUID mocks base method.
func (m *MockMachineService) GetMachineUUID(arg0 context.Context, arg1 machine.Name) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMachineUUID", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMachineUUID indicates an expected call of GetMachineUUID.
func (mr *MockMachineServiceMockRecorder) GetMachineUUID(arg0, arg1 any) *MockMachineServiceGetMachineUUIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMachineUUID", reflect.TypeOf((*MockMachineService)(nil).GetMachineUUID), arg0, arg1)
	return &MockMachineServiceGetMachineUUIDCall{Call: call}
}

// MockMachineServiceGetMachineUUIDCall wrap *gomock.Call
type MockMachineServiceGetMachineUUIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineServiceGetMachineUUIDCall) Return(arg0 string, arg1 error) *MockMachineServiceGetMachineUUIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineServiceGetMachineUUIDCall) Do(f func(context.Context, machine.Name) (string, error)) *MockMachineServiceGetMachineUUIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineServiceGetMachineUUIDCall) DoAndReturn(f func(context.Context, machine.Name) (string, error)) *MockMachineServiceGetMachineUUIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InstanceID mocks base method.
func (m *MockMachineService) InstanceID(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceID", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceID indicates an expected call of InstanceID.
func (mr *MockMachineServiceMockRecorder) InstanceID(arg0, arg1 any) *MockMachineServiceInstanceIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceID", reflect.TypeOf((*MockMachineService)(nil).InstanceID), arg0, arg1)
	return &MockMachineServiceInstanceIDCall{Call: call}
}

// MockMachineServiceInstanceIDCall wrap *gomock.Call
type MockMachineServiceInstanceIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineServiceInstanceIDCall) Return(arg0 string, arg1 error) *MockMachineServiceInstanceIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineServiceInstanceIDCall) Do(f func(context.Context, string) (string, error)) *MockMachineServiceInstanceIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineServiceInstanceIDCall) DoAndReturn(f func(context.Context, string) (string, error)) *MockMachineServiceInstanceIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// InstanceIDAndName mocks base method.
func (m *MockMachineService) InstanceIDAndName(arg0 context.Context, arg1 string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceIDAndName", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// InstanceIDAndName indicates an expected call of InstanceIDAndName.
func (mr *MockMachineServiceMockRecorder) InstanceIDAndName(arg0, arg1 any) *MockMachineServiceInstanceIDAndNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceIDAndName", reflect.TypeOf((*MockMachineService)(nil).InstanceIDAndName), arg0, arg1)
	return &MockMachineServiceInstanceIDAndNameCall{Call: call}
}

// MockMachineServiceInstanceIDAndNameCall wrap *gomock.Call
type MockMachineServiceInstanceIDAndNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockMachineServiceInstanceIDAndNameCall) Return(arg0, arg1 string, arg2 error) *MockMachineServiceInstanceIDAndNameCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockMachineServiceInstanceIDAndNameCall) Do(f func(context.Context, string) (string, string, error)) *MockMachineServiceInstanceIDAndNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockMachineServiceInstanceIDAndNameCall) DoAndReturn(f func(context.Context, string) (string, string, error)) *MockMachineServiceInstanceIDAndNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
