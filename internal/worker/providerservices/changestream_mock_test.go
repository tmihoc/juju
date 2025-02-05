// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/changestream (interfaces: WatchableDBGetter)
//
// Generated by this command:
//
//	mockgen -typed -package providerservices -destination changestream_mock_test.go github.com/juju/juju/core/changestream WatchableDBGetter
//

// Package providerservices is a generated GoMock package.
package providerservices

import (
	reflect "reflect"

	changestream "github.com/juju/juju/core/changestream"
	gomock "go.uber.org/mock/gomock"
)

// MockWatchableDBGetter is a mock of WatchableDBGetter interface.
type MockWatchableDBGetter struct {
	ctrl     *gomock.Controller
	recorder *MockWatchableDBGetterMockRecorder
}

// MockWatchableDBGetterMockRecorder is the mock recorder for MockWatchableDBGetter.
type MockWatchableDBGetterMockRecorder struct {
	mock *MockWatchableDBGetter
}

// NewMockWatchableDBGetter creates a new mock instance.
func NewMockWatchableDBGetter(ctrl *gomock.Controller) *MockWatchableDBGetter {
	mock := &MockWatchableDBGetter{ctrl: ctrl}
	mock.recorder = &MockWatchableDBGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatchableDBGetter) EXPECT() *MockWatchableDBGetterMockRecorder {
	return m.recorder
}

// GetWatchableDB mocks base method.
func (m *MockWatchableDBGetter) GetWatchableDB(arg0 string) (changestream.WatchableDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWatchableDB", arg0)
	ret0, _ := ret[0].(changestream.WatchableDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWatchableDB indicates an expected call of GetWatchableDB.
func (mr *MockWatchableDBGetterMockRecorder) GetWatchableDB(arg0 any) *MockWatchableDBGetterGetWatchableDBCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWatchableDB", reflect.TypeOf((*MockWatchableDBGetter)(nil).GetWatchableDB), arg0)
	return &MockWatchableDBGetterGetWatchableDBCall{Call: call}
}

// MockWatchableDBGetterGetWatchableDBCall wrap *gomock.Call
type MockWatchableDBGetterGetWatchableDBCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockWatchableDBGetterGetWatchableDBCall) Return(arg0 changestream.WatchableDB, arg1 error) *MockWatchableDBGetterGetWatchableDBCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockWatchableDBGetterGetWatchableDBCall) Do(f func(string) (changestream.WatchableDB, error)) *MockWatchableDBGetterGetWatchableDBCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockWatchableDBGetterGetWatchableDBCall) DoAndReturn(f func(string) (changestream.WatchableDB, error)) *MockWatchableDBGetterGetWatchableDBCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
