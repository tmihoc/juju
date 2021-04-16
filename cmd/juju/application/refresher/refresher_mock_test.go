// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/juju/application/refresher (interfaces: RefresherFactory,Refresher,CharmResolver,CharmRepository)

// Package refresher is a generated GoMock package.
package refresher

import (
	gomock "github.com/golang/mock/gomock"
	v8 "github.com/juju/charm/v8"
	charm "github.com/juju/juju/api/common/charm"
	reflect "reflect"
)

// MockRefresherFactory is a mock of RefresherFactory interface
type MockRefresherFactory struct {
	ctrl     *gomock.Controller
	recorder *MockRefresherFactoryMockRecorder
}

// MockRefresherFactoryMockRecorder is the mock recorder for MockRefresherFactory
type MockRefresherFactoryMockRecorder struct {
	mock *MockRefresherFactory
}

// NewMockRefresherFactory creates a new mock instance
func NewMockRefresherFactory(ctrl *gomock.Controller) *MockRefresherFactory {
	mock := &MockRefresherFactory{ctrl: ctrl}
	mock.recorder = &MockRefresherFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRefresherFactory) EXPECT() *MockRefresherFactoryMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockRefresherFactory) Run(arg0 RefresherConfig) (*CharmID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(*CharmID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run
func (mr *MockRefresherFactoryMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockRefresherFactory)(nil).Run), arg0)
}

// MockRefresher is a mock of Refresher interface
type MockRefresher struct {
	ctrl     *gomock.Controller
	recorder *MockRefresherMockRecorder
}

// MockRefresherMockRecorder is the mock recorder for MockRefresher
type MockRefresherMockRecorder struct {
	mock *MockRefresher
}

// NewMockRefresher creates a new mock instance
func NewMockRefresher(ctrl *gomock.Controller) *MockRefresher {
	mock := &MockRefresher{ctrl: ctrl}
	mock.recorder = &MockRefresherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRefresher) EXPECT() *MockRefresherMockRecorder {
	return m.recorder
}

// Allowed mocks base method
func (m *MockRefresher) Allowed(arg0 RefresherConfig) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Allowed", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Allowed indicates an expected call of Allowed
func (mr *MockRefresherMockRecorder) Allowed(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Allowed", reflect.TypeOf((*MockRefresher)(nil).Allowed), arg0)
}

// Refresh mocks base method
func (m *MockRefresher) Refresh() (*CharmID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(*CharmID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh
func (mr *MockRefresherMockRecorder) Refresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockRefresher)(nil).Refresh))
}

// String mocks base method
func (m *MockRefresher) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockRefresherMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockRefresher)(nil).String))
}

// MockCharmResolver is a mock of CharmResolver interface
type MockCharmResolver struct {
	ctrl     *gomock.Controller
	recorder *MockCharmResolverMockRecorder
}

// MockCharmResolverMockRecorder is the mock recorder for MockCharmResolver
type MockCharmResolverMockRecorder struct {
	mock *MockCharmResolver
}

// NewMockCharmResolver creates a new mock instance
func NewMockCharmResolver(ctrl *gomock.Controller) *MockCharmResolver {
	mock := &MockCharmResolver{ctrl: ctrl}
	mock.recorder = &MockCharmResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharmResolver) EXPECT() *MockCharmResolverMockRecorder {
	return m.recorder
}

// ResolveCharm mocks base method
func (m *MockCharmResolver) ResolveCharm(arg0 *v8.URL, arg1 charm.Origin) (*v8.URL, charm.Origin, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveCharm", arg0, arg1)
	ret0, _ := ret[0].(*v8.URL)
	ret1, _ := ret[1].(charm.Origin)
	ret2, _ := ret[2].([]string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ResolveCharm indicates an expected call of ResolveCharm
func (mr *MockCharmResolverMockRecorder) ResolveCharm(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveCharm", reflect.TypeOf((*MockCharmResolver)(nil).ResolveCharm), arg0, arg1)
}

// MockCharmRepository is a mock of CharmRepository interface
type MockCharmRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCharmRepositoryMockRecorder
}

// MockCharmRepositoryMockRecorder is the mock recorder for MockCharmRepository
type MockCharmRepositoryMockRecorder struct {
	mock *MockCharmRepository
}

// NewMockCharmRepository creates a new mock instance
func NewMockCharmRepository(ctrl *gomock.Controller) *MockCharmRepository {
	mock := &MockCharmRepository{ctrl: ctrl}
	mock.recorder = &MockCharmRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharmRepository) EXPECT() *MockCharmRepositoryMockRecorder {
	return m.recorder
}

// NewCharmAtPathForceSeries mocks base method
func (m *MockCharmRepository) NewCharmAtPathForceSeries(arg0, arg1 string, arg2 bool) (v8.Charm, *v8.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCharmAtPathForceSeries", arg0, arg1, arg2)
	ret0, _ := ret[0].(v8.Charm)
	ret1, _ := ret[1].(*v8.URL)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// NewCharmAtPathForceSeries indicates an expected call of NewCharmAtPathForceSeries
func (mr *MockCharmRepositoryMockRecorder) NewCharmAtPathForceSeries(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCharmAtPathForceSeries", reflect.TypeOf((*MockCharmRepository)(nil).NewCharmAtPathForceSeries), arg0, arg1, arg2)
}
