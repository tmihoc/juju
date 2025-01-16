// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/resource/downloader (interfaces: DownloadClient)
//
// Generated by this command:
//
//	mockgen -typed -package downloader_test -destination downloadclient_mock_test.go github.com/juju/juju/internal/resource/downloader DownloadClient
//

// Package downloader_test is a generated GoMock package.
package downloader_test

import (
	context "context"
	url "net/url"
	reflect "reflect"

	charmhub "github.com/juju/juju/internal/charmhub"
	gomock "go.uber.org/mock/gomock"
)

// MockDownloadClient is a mock of DownloadClient interface.
type MockDownloadClient struct {
	ctrl     *gomock.Controller
	recorder *MockDownloadClientMockRecorder
}

// MockDownloadClientMockRecorder is the mock recorder for MockDownloadClient.
type MockDownloadClientMockRecorder struct {
	mock *MockDownloadClient
}

// NewMockDownloadClient creates a new mock instance.
func NewMockDownloadClient(ctrl *gomock.Controller) *MockDownloadClient {
	mock := &MockDownloadClient{ctrl: ctrl}
	mock.recorder = &MockDownloadClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDownloadClient) EXPECT() *MockDownloadClientMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockDownloadClient) Download(arg0 context.Context, arg1 *url.URL, arg2 string, arg3 ...charmhub.DownloadOption) (*charmhub.Digest, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Download", varargs...)
	ret0, _ := ret[0].(*charmhub.Digest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockDownloadClientMockRecorder) Download(arg0, arg1, arg2 any, arg3 ...any) *MockDownloadClientDownloadCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2}, arg3...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockDownloadClient)(nil).Download), varargs...)
	return &MockDownloadClientDownloadCall{Call: call}
}

// MockDownloadClientDownloadCall wrap *gomock.Call
type MockDownloadClientDownloadCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockDownloadClientDownloadCall) Return(arg0 *charmhub.Digest, arg1 error) *MockDownloadClientDownloadCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockDownloadClientDownloadCall) Do(f func(context.Context, *url.URL, string, ...charmhub.DownloadOption) (*charmhub.Digest, error)) *MockDownloadClientDownloadCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockDownloadClientDownloadCall) DoAndReturn(f func(context.Context, *url.URL, string, ...charmhub.DownloadOption) (*charmhub.Digest, error)) *MockDownloadClientDownloadCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
