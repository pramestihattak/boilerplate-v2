// Code generated by MockGen. DO NOT EDIT.
// Source: ./backend/pkg/jwt/jwt.go
//
// Generated by this command:
//
//	mockgen -source ./backend/pkg/jwt/jwt.go -destination ./backend/pkg/jwt/mock/jwt_mock.go
//
// Package mock_jwt is a generated GoMock package.
package mock_jwt

import (
	jwt "boilerplate-v2/pkg/jwt"
	postgres "boilerplate-v2/storage/postgres"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockJWTInterface is a mock of JWTInterface interface.
type MockJWTInterface struct {
	ctrl     *gomock.Controller
	recorder *MockJWTInterfaceMockRecorder
}

// MockJWTInterfaceMockRecorder is the mock recorder for MockJWTInterface.
type MockJWTInterfaceMockRecorder struct {
	mock *MockJWTInterface
}

// NewMockJWTInterface creates a new mock instance.
func NewMockJWTInterface(ctrl *gomock.Controller) *MockJWTInterface {
	mock := &MockJWTInterface{ctrl: ctrl}
	mock.recorder = &MockJWTInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTInterface) EXPECT() *MockJWTInterfaceMockRecorder {
	return m.recorder
}

// GetClaims mocks base method.
func (m *MockJWTInterface) GetClaims(token string) (*jwt.Auth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClaims", token)
	ret0, _ := ret[0].(*jwt.Auth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClaims indicates an expected call of GetClaims.
func (mr *MockJWTInterfaceMockRecorder) GetClaims(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClaims", reflect.TypeOf((*MockJWTInterface)(nil).GetClaims), token)
}

// IsValidToken mocks base method.
func (m *MockJWTInterface) IsValidToken(token string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValidToken", token)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValidToken indicates an expected call of IsValidToken.
func (mr *MockJWTInterfaceMockRecorder) IsValidToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidToken", reflect.TypeOf((*MockJWTInterface)(nil).IsValidToken), token)
}

// Sign mocks base method.
func (m *MockJWTInterface) Sign(data *postgres.LoginOutput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockJWTInterfaceMockRecorder) Sign(data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockJWTInterface)(nil).Sign), data)
}

// MockJWTReader is a mock of JWTReader interface.
type MockJWTReader struct {
	ctrl     *gomock.Controller
	recorder *MockJWTReaderMockRecorder
}

// MockJWTReaderMockRecorder is the mock recorder for MockJWTReader.
type MockJWTReaderMockRecorder struct {
	mock *MockJWTReader
}

// NewMockJWTReader creates a new mock instance.
func NewMockJWTReader(ctrl *gomock.Controller) *MockJWTReader {
	mock := &MockJWTReader{ctrl: ctrl}
	mock.recorder = &MockJWTReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTReader) EXPECT() *MockJWTReaderMockRecorder {
	return m.recorder
}

// GetClaims mocks base method.
func (m *MockJWTReader) GetClaims(token string) (*jwt.Auth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClaims", token)
	ret0, _ := ret[0].(*jwt.Auth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClaims indicates an expected call of GetClaims.
func (mr *MockJWTReaderMockRecorder) GetClaims(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClaims", reflect.TypeOf((*MockJWTReader)(nil).GetClaims), token)
}

// IsValidToken mocks base method.
func (m *MockJWTReader) IsValidToken(token string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValidToken", token)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValidToken indicates an expected call of IsValidToken.
func (mr *MockJWTReaderMockRecorder) IsValidToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidToken", reflect.TypeOf((*MockJWTReader)(nil).IsValidToken), token)
}

// MockJWTWriter is a mock of JWTWriter interface.
type MockJWTWriter struct {
	ctrl     *gomock.Controller
	recorder *MockJWTWriterMockRecorder
}

// MockJWTWriterMockRecorder is the mock recorder for MockJWTWriter.
type MockJWTWriterMockRecorder struct {
	mock *MockJWTWriter
}

// NewMockJWTWriter creates a new mock instance.
func NewMockJWTWriter(ctrl *gomock.Controller) *MockJWTWriter {
	mock := &MockJWTWriter{ctrl: ctrl}
	mock.recorder = &MockJWTWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTWriter) EXPECT() *MockJWTWriterMockRecorder {
	return m.recorder
}

// Sign mocks base method.
func (m *MockJWTWriter) Sign(data *postgres.LoginOutput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockJWTWriterMockRecorder) Sign(data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockJWTWriter)(nil).Sign), data)
}