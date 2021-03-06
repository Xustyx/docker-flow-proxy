// +build !integration

package main

import (
	"./proxy"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"testing"
)

type ArgsTestSuite struct {
	suite.Suite
	args Args
}

func (s *ArgsTestSuite) SetupTest() {
	httpListenAndServe = func(addr string, handler http.Handler) error {
		return nil
	}
	osRemove = func(name string) error {
		return nil
	}
}

// NewArgs

func (s ArgsTestSuite) Test_NewArgs_ReturnsNewStruct() {
	a := NewArgs()

	s.IsType(Args{}, a)
}

// Parse

func (s ArgsTestSuite) Test_Parse_ReturnsError_WhenFailure() {
	os.Args = []string{"myProgram", "myCommand", "--this-flag-does-not-exist=something"}

	actual := Args{}.Parse()

	s.Error(actual)
}

// Parse > Reconfigure

func (s ArgsTestSuite) Test_Parse_ParsesReconfigureLongArgsStrings() {
	os.Args = []string{"myProgram", "reconfigure"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"serviceNameFromArgs", "service-name", &reconfigure.ServiceName},
		{"serviceColorFromArgs", "service-color", &reconfigure.ServiceColor},
		{"instanceNameFromArgs", "proxy-instance-name", &reconfigure.InstanceName},
		{"templatesPathFromArgs", "templates-path", &reconfigure.TemplatesPath},
		{"configsPathFromArgs", "configs-path", &reconfigure.ConfigsPath},
		{"consulTemplateFePath", "consul-template-fe-path", &reconfigure.ConsulTemplateFePath},
		{"consulTemplateBePath", "consul-template-be-path", &reconfigure.ConsulTemplateBePath},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("--%s", d.key), d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ParsesReconfigureLongArgsSlices() {
	os.Args = []string{"myProgram", "reconfigure"}
	data := []struct {
		expected []string
		key      string
		value    *[]string
	}{
		{[]string{"path1", "path2"}, "service-path", &reconfigure.ServicePath},
		{[]string{"service-domain"}, "service-domain", &reconfigure.ServiceDomain},
	}

	for _, d := range data {
		for _, v := range d.expected {
			os.Args = append(os.Args, fmt.Sprintf("--%s", d.key), v)
		}

	}

	Args{}.Parse()

	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ParsesReconfigureShortArgsStrings() {
	os.Args = []string{"myProgram", "reconfigure"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"serviceNameFromArgs", "s", &reconfigure.ServiceName},
		{"serviceColorFromArgs", "C", &reconfigure.ServiceColor},
		{"templatesPathFromArgs", "t", &reconfigure.TemplatesPath},
		{"configsPathFromArgs", "c", &reconfigure.ConfigsPath},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("-%s", d.key), d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ParsesReconfigureShortArgsSlices() {
	os.Args = []string{"myProgram", "reconfigure"}
	data := []struct {
		expected []string
		key      string
		value    *[]string
	}{
		{[]string{"p1", "p2"}, "p", &reconfigure.ServicePath},
	}
	for _, d := range data {
		for _, v := range d.expected {
			os.Args = append(os.Args, fmt.Sprintf("-%s", d.key), v)
		}
	}

	Args{}.Parse()

	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ReconfigureHasDefaultValues() {
	os.Args = []string{
		"myProgram", "reconfigure",
		"--service-name", "myService",
		"--service-path", "my/service/path",
	}
	data := []struct {
		expected string
		value    *string
	}{
		{"/cfg/tmpl", &reconfigure.TemplatesPath},
		{"/cfg", &reconfigure.ConfigsPath},
	}
	reconfigure.ConsulAddresses = []string{"myConsulAddress"}
	reconfigure.ServicePath = []string{"p1", "p2"}
	reconfigure.ServiceName = "myServiceName"

	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ReconfigureDefaultsToEnvVars() {
	os.Args = []string{
		"myProgram", "reconfigure",
		"--service-name", "serviceName",
		"--service-path", "servicePath",
	}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"proxyInstanceNameFromEnv", "PROXY_INSTANCE_NAME", &reconfigure.InstanceName},
	}

	for _, d := range data {
		os.Setenv(d.key, d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

// Parse > Remove

func (s ArgsTestSuite) Test_Parse_ParsesRemoveLongArgsStrings() {
	os.Args = []string{"myProgram", "remove"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"serviceNameFromArgs", "service-name", &remove.ServiceName},
		{"templatesPathFromArgs", "templates-path", &remove.TemplatesPath},
		{"configsPathFromArgs", "configs-path", &remove.ConfigsPath},
		{"instanceNameFromArgs", "proxy-instance-name", &remove.InstanceName},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("--%s", d.key), d.expected)
	}
	err := Args{}.Parse()
	s.NoError(err)
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ParsesRemoveShortArgsStrings() {
	os.Args = []string{"myProgram", "remove"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"serviceNameFromArgs", "s", &remove.ServiceName},
		{"templatesPathFromArgs", "t", &remove.TemplatesPath},
		{"configsPathFromArgs", "c", &remove.ConfigsPath},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("-%s", d.key), d.expected)
	}
	err := Args{}.Parse()
	s.NoError(err)
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_RemoveDefaultsToEnvVars() {
	os.Args = []string{"myProgram", "remove"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"proxyInstanceNameFromEnv", "PROXY_INSTANCE_NAME", &remove.InstanceName},
	}

	for _, d := range data {
		os.Setenv(d.key, d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

// Parse > Server

func (s ArgsTestSuite) Test_Parse_ParsesServerLongArgs() {
	os.Args = []string{"myProgram", "server"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"ipFromArgs", "ip", &serverImpl.IP},
		{"portFromArgs", "port", &serverImpl.Port},
		{"modeFromArgs", "mode", &serverImpl.Mode},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("--%s", d.key), d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ParsesServerShortArgs() {
	os.Args = []string{"myProgram", "server"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"ipFromArgs", "i", &serverImpl.IP},
		{"portFromArgs", "p", &serverImpl.Port},
		{"modeFromArgs", "m", &serverImpl.Mode},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("-%s", d.key), d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ServerHasDefaultValues() {
	os.Args = []string{"myProgram", "server"}
	os.Unsetenv("IP")
	os.Unsetenv("PORT")
	data := []struct {
		expected string
		value    *string
	}{
		{"0.0.0.0", &serverImpl.IP},
		{"8080", &serverImpl.Port},
	}

	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ServerDefaultsToEnvVars() {
	os.Args = []string{"myProgram", "server"}
	data := []struct {
		expected string
		key      string
		value    *string
	}{
		{"ipFromEnv", "IP", &serverImpl.IP},
		{"portFromEnv", "PORT", &serverImpl.Port},
		{"modeFromEnv", "MODE", &serverImpl.Mode},
	}

	for _, d := range data {
		os.Setenv(d.key, d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

// Suite

func TestArgsUnitTestSuite(t *testing.T) {
	mockObj := getRegistrarableMock("")
	registryInstanceOrig := registryInstance
	defer func() { registryInstance = registryInstanceOrig }()
	registryInstance = mockObj
	logPrintf = func(format string, v ...interface{}) {}
	proxyOrig := proxy.Instance
	defer func() { proxy.Instance = proxyOrig }()
	proxy.Instance = getProxyMock("")
	suite.Run(t, new(ArgsTestSuite))
}

// Mock

type ArgsMock struct {
	mock.Mock
}

func (m *ArgsMock) Parse(args *Args) error {
	params := m.Called(args)
	return params.Error(0)
}

func getArgsMock() *ArgsMock {
	mockObj := new(ArgsMock)
	mockObj.On("Parse", mock.Anything).Return(nil)
	return mockObj
}

type ProxyMock struct {
	mock.Mock
}

func (m *ProxyMock) RunCmd(extraArgs []string) error {
	params := m.Called(extraArgs)
	return params.Error(0)
}

func (m *ProxyMock) CreateConfigFromTemplates() error {
	params := m.Called()
	return params.Error(0)
}

func (m *ProxyMock) ReadConfig() (string, error) {
	params := m.Called()
	return params.String(0), params.Error(1)
}

func (m *ProxyMock) Reload() error {
	params := m.Called()
	return params.Error(0)
}

func (m *ProxyMock) AddCert(certName string) {
	m.Called(certName)
}

func (m *ProxyMock) GetCerts() map[string]string {
	params := m.Called()
	return params.Get(0).(map[string]string)
}

func getProxyMock(skipMethod string) *ProxyMock {
	mockObj := new(ProxyMock)
	if skipMethod != "RunCmd" {
		mockObj.On("RunCmd", mock.Anything).Return(nil)
	}
	if skipMethod != "CreateConfigFromTemplates" {
		mockObj.On("CreateConfigFromTemplates").Return(nil)
	}
	if skipMethod != "ReadConfig" {
		mockObj.On("ReadConfig").Return("", nil)
	}
	if skipMethod != "Reload" {
		mockObj.On("Reload").Return(nil)
	}
	if skipMethod != "AddCert" {
		mockObj.On("AddCert", mock.Anything).Return(nil)
	}
	if skipMethod != "GetCerts" {
		mockObj.On("GetCerts").Return(map[string]string{})
	}
	return mockObj
}
