// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/agent (interfaces: Agent,Config)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	agent "github.com/juju/juju/agent"
	api "github.com/juju/juju/api"
	controller "github.com/juju/juju/controller"
	model "github.com/juju/juju/core/model"
	mongo "github.com/juju/juju/internal/mongo"
	names "github.com/juju/names/v4"
	shell "github.com/juju/utils/v3/shell"
	version "github.com/juju/version/v2"
	gomock "go.uber.org/mock/gomock"
)

// MockAgent is a mock of Agent interface.
type MockAgent struct {
	ctrl     *gomock.Controller
	recorder *MockAgentMockRecorder
}

// MockAgentMockRecorder is the mock recorder for MockAgent.
type MockAgentMockRecorder struct {
	mock *MockAgent
}

// NewMockAgent creates a new mock instance.
func NewMockAgent(ctrl *gomock.Controller) *MockAgent {
	mock := &MockAgent{ctrl: ctrl}
	mock.recorder = &MockAgentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgent) EXPECT() *MockAgentMockRecorder {
	return m.recorder
}

// ChangeConfig mocks base method.
func (m *MockAgent) ChangeConfig(arg0 agent.ConfigMutator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeConfig indicates an expected call of ChangeConfig.
func (mr *MockAgentMockRecorder) ChangeConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeConfig", reflect.TypeOf((*MockAgent)(nil).ChangeConfig), arg0)
}

// CurrentConfig mocks base method.
func (m *MockAgent) CurrentConfig() agent.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentConfig")
	ret0, _ := ret[0].(agent.Config)
	return ret0
}

// CurrentConfig indicates an expected call of CurrentConfig.
func (mr *MockAgentMockRecorder) CurrentConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentConfig", reflect.TypeOf((*MockAgent)(nil).CurrentConfig))
}

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// APIAddresses mocks base method.
func (m *MockConfig) APIAddresses() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIAddresses")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// APIAddresses indicates an expected call of APIAddresses.
func (mr *MockConfigMockRecorder) APIAddresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIAddresses", reflect.TypeOf((*MockConfig)(nil).APIAddresses))
}

// APIInfo mocks base method.
func (m *MockConfig) APIInfo() (*api.Info, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIInfo")
	ret0, _ := ret[0].(*api.Info)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// APIInfo indicates an expected call of APIInfo.
func (mr *MockConfigMockRecorder) APIInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIInfo", reflect.TypeOf((*MockConfig)(nil).APIInfo))
}

// AgentLogfileMaxBackups mocks base method.
func (m *MockConfig) AgentLogfileMaxBackups() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentLogfileMaxBackups")
	ret0, _ := ret[0].(int)
	return ret0
}

// AgentLogfileMaxBackups indicates an expected call of AgentLogfileMaxBackups.
func (mr *MockConfigMockRecorder) AgentLogfileMaxBackups() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentLogfileMaxBackups", reflect.TypeOf((*MockConfig)(nil).AgentLogfileMaxBackups))
}

// AgentLogfileMaxSizeMB mocks base method.
func (m *MockConfig) AgentLogfileMaxSizeMB() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentLogfileMaxSizeMB")
	ret0, _ := ret[0].(int)
	return ret0
}

// AgentLogfileMaxSizeMB indicates an expected call of AgentLogfileMaxSizeMB.
func (mr *MockConfigMockRecorder) AgentLogfileMaxSizeMB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentLogfileMaxSizeMB", reflect.TypeOf((*MockConfig)(nil).AgentLogfileMaxSizeMB))
}

// CACert mocks base method.
func (m *MockConfig) CACert() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CACert")
	ret0, _ := ret[0].(string)
	return ret0
}

// CACert indicates an expected call of CACert.
func (mr *MockConfigMockRecorder) CACert() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CACert", reflect.TypeOf((*MockConfig)(nil).CACert))
}

// Controller mocks base method.
func (m *MockConfig) Controller() names.ControllerTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Controller")
	ret0, _ := ret[0].(names.ControllerTag)
	return ret0
}

// Controller indicates an expected call of Controller.
func (mr *MockConfigMockRecorder) Controller() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Controller", reflect.TypeOf((*MockConfig)(nil).Controller))
}

// DataDir mocks base method.
func (m *MockConfig) DataDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DataDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// DataDir indicates an expected call of DataDir.
func (mr *MockConfigMockRecorder) DataDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DataDir", reflect.TypeOf((*MockConfig)(nil).DataDir))
}

// Dir mocks base method.
func (m *MockConfig) Dir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dir")
	ret0, _ := ret[0].(string)
	return ret0
}

// Dir indicates an expected call of Dir.
func (mr *MockConfigMockRecorder) Dir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dir", reflect.TypeOf((*MockConfig)(nil).Dir))
}

// DqlitePort mocks base method.
func (m *MockConfig) DqlitePort() (int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DqlitePort")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// DqlitePort indicates an expected call of DqlitePort.
func (mr *MockConfigMockRecorder) DqlitePort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DqlitePort", reflect.TypeOf((*MockConfig)(nil).DqlitePort))
}

// Jobs mocks base method.
func (m *MockConfig) Jobs() []model.MachineJob {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Jobs")
	ret0, _ := ret[0].([]model.MachineJob)
	return ret0
}

// Jobs indicates an expected call of Jobs.
func (mr *MockConfigMockRecorder) Jobs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Jobs", reflect.TypeOf((*MockConfig)(nil).Jobs))
}

// JujuDBSnapChannel mocks base method.
func (m *MockConfig) JujuDBSnapChannel() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JujuDBSnapChannel")
	ret0, _ := ret[0].(string)
	return ret0
}

// JujuDBSnapChannel indicates an expected call of JujuDBSnapChannel.
func (mr *MockConfigMockRecorder) JujuDBSnapChannel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JujuDBSnapChannel", reflect.TypeOf((*MockConfig)(nil).JujuDBSnapChannel))
}

// LogDir mocks base method.
func (m *MockConfig) LogDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// LogDir indicates an expected call of LogDir.
func (mr *MockConfigMockRecorder) LogDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogDir", reflect.TypeOf((*MockConfig)(nil).LogDir))
}

// LoggingConfig mocks base method.
func (m *MockConfig) LoggingConfig() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoggingConfig")
	ret0, _ := ret[0].(string)
	return ret0
}

// LoggingConfig indicates an expected call of LoggingConfig.
func (mr *MockConfigMockRecorder) LoggingConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoggingConfig", reflect.TypeOf((*MockConfig)(nil).LoggingConfig))
}

// MetricsSpoolDir mocks base method.
func (m *MockConfig) MetricsSpoolDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MetricsSpoolDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// MetricsSpoolDir indicates an expected call of MetricsSpoolDir.
func (mr *MockConfigMockRecorder) MetricsSpoolDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MetricsSpoolDir", reflect.TypeOf((*MockConfig)(nil).MetricsSpoolDir))
}

// Model mocks base method.
func (m *MockConfig) Model() names.ModelTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Model")
	ret0, _ := ret[0].(names.ModelTag)
	return ret0
}

// Model indicates an expected call of Model.
func (mr *MockConfigMockRecorder) Model() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockConfig)(nil).Model))
}

// MongoInfo mocks base method.
func (m *MockConfig) MongoInfo() (*mongo.MongoInfo, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MongoInfo")
	ret0, _ := ret[0].(*mongo.MongoInfo)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// MongoInfo indicates an expected call of MongoInfo.
func (mr *MockConfigMockRecorder) MongoInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MongoInfo", reflect.TypeOf((*MockConfig)(nil).MongoInfo))
}

// MongoMemoryProfile mocks base method.
func (m *MockConfig) MongoMemoryProfile() mongo.MemoryProfile {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MongoMemoryProfile")
	ret0, _ := ret[0].(mongo.MemoryProfile)
	return ret0
}

// MongoMemoryProfile indicates an expected call of MongoMemoryProfile.
func (mr *MockConfigMockRecorder) MongoMemoryProfile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MongoMemoryProfile", reflect.TypeOf((*MockConfig)(nil).MongoMemoryProfile))
}

// Nonce mocks base method.
func (m *MockConfig) Nonce() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nonce")
	ret0, _ := ret[0].(string)
	return ret0
}

// Nonce indicates an expected call of Nonce.
func (mr *MockConfigMockRecorder) Nonce() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nonce", reflect.TypeOf((*MockConfig)(nil).Nonce))
}

// OldPassword mocks base method.
func (m *MockConfig) OldPassword() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OldPassword")
	ret0, _ := ret[0].(string)
	return ret0
}

// OldPassword indicates an expected call of OldPassword.
func (mr *MockConfigMockRecorder) OldPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OldPassword", reflect.TypeOf((*MockConfig)(nil).OldPassword))
}

// OpenTelemetryEnabled mocks base method.
func (m *MockConfig) OpenTelemetryEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenTelemetryEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// OpenTelemetryEnabled indicates an expected call of OpenTelemetryEnabled.
func (mr *MockConfigMockRecorder) OpenTelemetryEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenTelemetryEnabled", reflect.TypeOf((*MockConfig)(nil).OpenTelemetryEnabled))
}

// OpenTelemetryEndpoint mocks base method.
func (m *MockConfig) OpenTelemetryEndpoint() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenTelemetryEndpoint")
	ret0, _ := ret[0].(string)
	return ret0
}

// OpenTelemetryEndpoint indicates an expected call of OpenTelemetryEndpoint.
func (mr *MockConfigMockRecorder) OpenTelemetryEndpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenTelemetryEndpoint", reflect.TypeOf((*MockConfig)(nil).OpenTelemetryEndpoint))
}

// OpenTelemetryInsecure mocks base method.
func (m *MockConfig) OpenTelemetryInsecure() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenTelemetryInsecure")
	ret0, _ := ret[0].(bool)
	return ret0
}

// OpenTelemetryInsecure indicates an expected call of OpenTelemetryInsecure.
func (mr *MockConfigMockRecorder) OpenTelemetryInsecure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenTelemetryInsecure", reflect.TypeOf((*MockConfig)(nil).OpenTelemetryInsecure))
}

// OpenTelemetryStackTraces mocks base method.
func (m *MockConfig) OpenTelemetryStackTraces() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenTelemetryStackTraces")
	ret0, _ := ret[0].(bool)
	return ret0
}

// OpenTelemetryStackTraces indicates an expected call of OpenTelemetryStackTraces.
func (mr *MockConfigMockRecorder) OpenTelemetryStackTraces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenTelemetryStackTraces", reflect.TypeOf((*MockConfig)(nil).OpenTelemetryStackTraces))
}

// QueryTracingEnabled mocks base method.
func (m *MockConfig) QueryTracingEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTracingEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// QueryTracingEnabled indicates an expected call of QueryTracingEnabled.
func (mr *MockConfigMockRecorder) QueryTracingEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTracingEnabled", reflect.TypeOf((*MockConfig)(nil).QueryTracingEnabled))
}

// QueryTracingThreshold mocks base method.
func (m *MockConfig) QueryTracingThreshold() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTracingThreshold")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// QueryTracingThreshold indicates an expected call of QueryTracingThreshold.
func (mr *MockConfigMockRecorder) QueryTracingThreshold() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTracingThreshold", reflect.TypeOf((*MockConfig)(nil).QueryTracingThreshold))
}

// StateServingInfo mocks base method.
func (m *MockConfig) StateServingInfo() (controller.StateServingInfo, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateServingInfo")
	ret0, _ := ret[0].(controller.StateServingInfo)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// StateServingInfo indicates an expected call of StateServingInfo.
func (mr *MockConfigMockRecorder) StateServingInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateServingInfo", reflect.TypeOf((*MockConfig)(nil).StateServingInfo))
}

// SystemIdentityPath mocks base method.
func (m *MockConfig) SystemIdentityPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SystemIdentityPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// SystemIdentityPath indicates an expected call of SystemIdentityPath.
func (mr *MockConfigMockRecorder) SystemIdentityPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SystemIdentityPath", reflect.TypeOf((*MockConfig)(nil).SystemIdentityPath))
}

// Tag mocks base method.
func (m *MockConfig) Tag() names.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names.Tag)
	return ret0
}

// Tag indicates an expected call of Tag.
func (mr *MockConfigMockRecorder) Tag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockConfig)(nil).Tag))
}

// TransientDataDir mocks base method.
func (m *MockConfig) TransientDataDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransientDataDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// TransientDataDir indicates an expected call of TransientDataDir.
func (mr *MockConfigMockRecorder) TransientDataDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransientDataDir", reflect.TypeOf((*MockConfig)(nil).TransientDataDir))
}

// UpgradedToVersion mocks base method.
func (m *MockConfig) UpgradedToVersion() version.Number {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpgradedToVersion")
	ret0, _ := ret[0].(version.Number)
	return ret0
}

// UpgradedToVersion indicates an expected call of UpgradedToVersion.
func (mr *MockConfigMockRecorder) UpgradedToVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradedToVersion", reflect.TypeOf((*MockConfig)(nil).UpgradedToVersion))
}

// Value mocks base method.
func (m *MockConfig) Value(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Value indicates an expected call of Value.
func (mr *MockConfigMockRecorder) Value(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockConfig)(nil).Value), arg0)
}

// WriteCommands mocks base method.
func (m *MockConfig) WriteCommands(arg0 shell.Renderer) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteCommands", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteCommands indicates an expected call of WriteCommands.
func (mr *MockConfigMockRecorder) WriteCommands(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteCommands", reflect.TypeOf((*MockConfig)(nil).WriteCommands), arg0)
}
