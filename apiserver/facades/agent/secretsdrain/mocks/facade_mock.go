// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facade (interfaces: Context,Authorizer)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	facade "github.com/juju/juju/apiserver/facade"
	cache "github.com/juju/juju/core/cache"
	database "github.com/juju/juju/core/database"
	leadership "github.com/juju/juju/core/leadership"
	lease "github.com/juju/juju/core/lease"
	multiwatcher "github.com/juju/juju/core/multiwatcher"
	permission "github.com/juju/juju/core/permission"
	state "github.com/juju/juju/state"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockContext is a mock of Context interface.
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext.
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance.
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockContext) Auth() facade.Authorizer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth")
	ret0, _ := ret[0].(facade.Authorizer)
	return ret0
}

// Auth indicates an expected call of Auth.
func (mr *MockContextMockRecorder) Auth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockContext)(nil).Auth))
}

// CachedModel mocks base method.
func (m *MockContext) CachedModel(arg0 string) (*cache.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CachedModel", arg0)
	ret0, _ := ret[0].(*cache.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CachedModel indicates an expected call of CachedModel.
func (mr *MockContextMockRecorder) CachedModel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CachedModel", reflect.TypeOf((*MockContext)(nil).CachedModel), arg0)
}

// Cancel mocks base method.
func (m *MockContext) Cancel() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Cancel indicates an expected call of Cancel.
func (mr *MockContextMockRecorder) Cancel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockContext)(nil).Cancel))
}

// Controller mocks base method.
func (m *MockContext) Controller() *cache.Controller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Controller")
	ret0, _ := ret[0].(*cache.Controller)
	return ret0
}

// Controller indicates an expected call of Controller.
func (mr *MockContextMockRecorder) Controller() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Controller", reflect.TypeOf((*MockContext)(nil).Controller))
}

// ControllerDB mocks base method.
func (m *MockContext) ControllerDB() (database.TrackedDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerDB")
	ret0, _ := ret[0].(database.TrackedDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerDB indicates an expected call of ControllerDB.
func (mr *MockContextMockRecorder) ControllerDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerDB", reflect.TypeOf((*MockContext)(nil).ControllerDB))
}

// Dispose mocks base method.
func (m *MockContext) Dispose() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispose")
}

// Dispose indicates an expected call of Dispose.
func (mr *MockContextMockRecorder) Dispose() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispose", reflect.TypeOf((*MockContext)(nil).Dispose))
}

// HTTPClient mocks base method.
func (m *MockContext) HTTPClient(arg0 facade.HTTPClientPurpose) facade.HTTPClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HTTPClient", arg0)
	ret0, _ := ret[0].(facade.HTTPClient)
	return ret0
}

// HTTPClient indicates an expected call of HTTPClient.
func (mr *MockContextMockRecorder) HTTPClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPClient", reflect.TypeOf((*MockContext)(nil).HTTPClient), arg0)
}

// Hub mocks base method.
func (m *MockContext) Hub() facade.Hub {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hub")
	ret0, _ := ret[0].(facade.Hub)
	return ret0
}

// Hub indicates an expected call of Hub.
func (mr *MockContextMockRecorder) Hub() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hub", reflect.TypeOf((*MockContext)(nil).Hub))
}

// ID mocks base method.
func (m *MockContext) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockContextMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockContext)(nil).ID))
}

// LeadershipChecker mocks base method.
func (m *MockContext) LeadershipChecker() (leadership.Checker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeadershipChecker")
	ret0, _ := ret[0].(leadership.Checker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipChecker indicates an expected call of LeadershipChecker.
func (mr *MockContextMockRecorder) LeadershipChecker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipChecker", reflect.TypeOf((*MockContext)(nil).LeadershipChecker))
}

// LeadershipClaimer mocks base method.
func (m *MockContext) LeadershipClaimer(arg0 string) (leadership.Claimer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeadershipClaimer", arg0)
	ret0, _ := ret[0].(leadership.Claimer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipClaimer indicates an expected call of LeadershipClaimer.
func (mr *MockContextMockRecorder) LeadershipClaimer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipClaimer", reflect.TypeOf((*MockContext)(nil).LeadershipClaimer), arg0)
}

// LeadershipPinner mocks base method.
func (m *MockContext) LeadershipPinner(arg0 string) (leadership.Pinner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeadershipPinner", arg0)
	ret0, _ := ret[0].(leadership.Pinner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipPinner indicates an expected call of LeadershipPinner.
func (mr *MockContextMockRecorder) LeadershipPinner(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipPinner", reflect.TypeOf((*MockContext)(nil).LeadershipPinner), arg0)
}

// LeadershipReader mocks base method.
func (m *MockContext) LeadershipReader(arg0 string) (leadership.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeadershipReader", arg0)
	ret0, _ := ret[0].(leadership.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipReader indicates an expected call of LeadershipReader.
func (mr *MockContextMockRecorder) LeadershipReader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipReader", reflect.TypeOf((*MockContext)(nil).LeadershipReader), arg0)
}

// LeadershipRevoker mocks base method.
func (m *MockContext) LeadershipRevoker(arg0 string) (leadership.Revoker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeadershipRevoker", arg0)
	ret0, _ := ret[0].(leadership.Revoker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipRevoker indicates an expected call of LeadershipRevoker.
func (mr *MockContextMockRecorder) LeadershipRevoker(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipRevoker", reflect.TypeOf((*MockContext)(nil).LeadershipRevoker), arg0)
}

// MultiwatcherFactory mocks base method.
func (m *MockContext) MultiwatcherFactory() multiwatcher.Factory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiwatcherFactory")
	ret0, _ := ret[0].(multiwatcher.Factory)
	return ret0
}

// MultiwatcherFactory indicates an expected call of MultiwatcherFactory.
func (mr *MockContextMockRecorder) MultiwatcherFactory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiwatcherFactory", reflect.TypeOf((*MockContext)(nil).MultiwatcherFactory))
}

// Presence mocks base method.
func (m *MockContext) Presence() facade.Presence {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Presence")
	ret0, _ := ret[0].(facade.Presence)
	return ret0
}

// Presence indicates an expected call of Presence.
func (mr *MockContextMockRecorder) Presence() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Presence", reflect.TypeOf((*MockContext)(nil).Presence))
}

// RequestRecorder mocks base method.
func (m *MockContext) RequestRecorder() facade.RequestRecorder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestRecorder")
	ret0, _ := ret[0].(facade.RequestRecorder)
	return ret0
}

// RequestRecorder indicates an expected call of RequestRecorder.
func (mr *MockContextMockRecorder) RequestRecorder() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestRecorder", reflect.TypeOf((*MockContext)(nil).RequestRecorder))
}

// Resources mocks base method.
func (m *MockContext) Resources() facade.Resources {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resources")
	ret0, _ := ret[0].(facade.Resources)
	return ret0
}

// Resources indicates an expected call of Resources.
func (mr *MockContextMockRecorder) Resources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resources", reflect.TypeOf((*MockContext)(nil).Resources))
}

// SingularClaimer mocks base method.
func (m *MockContext) SingularClaimer() (lease.Claimer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SingularClaimer")
	ret0, _ := ret[0].(lease.Claimer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SingularClaimer indicates an expected call of SingularClaimer.
func (mr *MockContextMockRecorder) SingularClaimer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SingularClaimer", reflect.TypeOf((*MockContext)(nil).SingularClaimer))
}

// State mocks base method.
func (m *MockContext) State() *state.State {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(*state.State)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockContextMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockContext)(nil).State))
}

// StatePool mocks base method.
func (m *MockContext) StatePool() *state.StatePool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StatePool")
	ret0, _ := ret[0].(*state.StatePool)
	return ret0
}

// StatePool indicates an expected call of StatePool.
func (mr *MockContextMockRecorder) StatePool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatePool", reflect.TypeOf((*MockContext)(nil).StatePool))
}

// MockAuthorizer is a mock of Authorizer interface.
type MockAuthorizer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizerMockRecorder
}

// MockAuthorizerMockRecorder is the mock recorder for MockAuthorizer.
type MockAuthorizerMockRecorder struct {
	mock *MockAuthorizer
}

// NewMockAuthorizer creates a new mock instance.
func NewMockAuthorizer(ctrl *gomock.Controller) *MockAuthorizer {
	mock := &MockAuthorizer{ctrl: ctrl}
	mock.recorder = &MockAuthorizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizer) EXPECT() *MockAuthorizerMockRecorder {
	return m.recorder
}

// AuthApplicationAgent mocks base method.
func (m *MockAuthorizer) AuthApplicationAgent() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthApplicationAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthApplicationAgent indicates an expected call of AuthApplicationAgent.
func (mr *MockAuthorizerMockRecorder) AuthApplicationAgent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthApplicationAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthApplicationAgent))
}

// AuthClient mocks base method.
func (m *MockAuthorizer) AuthClient() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthClient")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthClient indicates an expected call of AuthClient.
func (mr *MockAuthorizerMockRecorder) AuthClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthClient", reflect.TypeOf((*MockAuthorizer)(nil).AuthClient))
}

// AuthController mocks base method.
func (m *MockAuthorizer) AuthController() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthController")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthController indicates an expected call of AuthController.
func (mr *MockAuthorizerMockRecorder) AuthController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthController", reflect.TypeOf((*MockAuthorizer)(nil).AuthController))
}

// AuthMachineAgent mocks base method.
func (m *MockAuthorizer) AuthMachineAgent() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthMachineAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthMachineAgent indicates an expected call of AuthMachineAgent.
func (mr *MockAuthorizerMockRecorder) AuthMachineAgent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthMachineAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthMachineAgent))
}

// AuthModelAgent mocks base method.
func (m *MockAuthorizer) AuthModelAgent() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthModelAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthModelAgent indicates an expected call of AuthModelAgent.
func (mr *MockAuthorizerMockRecorder) AuthModelAgent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthModelAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthModelAgent))
}

// AuthOwner mocks base method.
func (m *MockAuthorizer) AuthOwner(arg0 names.Tag) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthOwner", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthOwner indicates an expected call of AuthOwner.
func (mr *MockAuthorizerMockRecorder) AuthOwner(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthOwner", reflect.TypeOf((*MockAuthorizer)(nil).AuthOwner), arg0)
}

// AuthUnitAgent mocks base method.
func (m *MockAuthorizer) AuthUnitAgent() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthUnitAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthUnitAgent indicates an expected call of AuthUnitAgent.
func (mr *MockAuthorizerMockRecorder) AuthUnitAgent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUnitAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthUnitAgent))
}

// ConnectedModel mocks base method.
func (m *MockAuthorizer) ConnectedModel() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectedModel")
	ret0, _ := ret[0].(string)
	return ret0
}

// ConnectedModel indicates an expected call of ConnectedModel.
func (mr *MockAuthorizerMockRecorder) ConnectedModel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectedModel", reflect.TypeOf((*MockAuthorizer)(nil).ConnectedModel))
}

// EntityHasPermission mocks base method.
func (m *MockAuthorizer) EntityHasPermission(arg0 names.Tag, arg1 permission.Access, arg2 names.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EntityHasPermission", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// EntityHasPermission indicates an expected call of EntityHasPermission.
func (mr *MockAuthorizerMockRecorder) EntityHasPermission(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EntityHasPermission", reflect.TypeOf((*MockAuthorizer)(nil).EntityHasPermission), arg0, arg1, arg2)
}

// GetAuthTag mocks base method.
func (m *MockAuthorizer) GetAuthTag() names.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthTag")
	ret0, _ := ret[0].(names.Tag)
	return ret0
}

// GetAuthTag indicates an expected call of GetAuthTag.
func (mr *MockAuthorizerMockRecorder) GetAuthTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthTag", reflect.TypeOf((*MockAuthorizer)(nil).GetAuthTag))
}

// HasPermission mocks base method.
func (m *MockAuthorizer) HasPermission(arg0 permission.Access, arg1 names.Tag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasPermission", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HasPermission indicates an expected call of HasPermission.
func (mr *MockAuthorizerMockRecorder) HasPermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPermission", reflect.TypeOf((*MockAuthorizer)(nil).HasPermission), arg0, arg1)
}
