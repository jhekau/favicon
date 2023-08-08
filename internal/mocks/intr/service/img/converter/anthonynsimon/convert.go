// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/img/converter/anthonynsimon/convert.go

// Package mock_converter_exec_anthonynsimon is a generated GoMock package.
package mock_converter_exec_anthonynsimon

import (
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorageOBJ is a mock of StorageOBJ interface.
type MockStorageOBJ struct {
	ctrl     *gomock.Controller
	recorder *MockStorageOBJMockRecorder
}

// MockStorageOBJMockRecorder is the mock recorder for MockStorageOBJ.
type MockStorageOBJMockRecorder struct {
	mock *MockStorageOBJ
}

// NewMockStorageOBJ creates a new mock instance.
func NewMockStorageOBJ(ctrl *gomock.Controller) *MockStorageOBJ {
	mock := &MockStorageOBJ{ctrl: ctrl}
	mock.recorder = &MockStorageOBJMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageOBJ) EXPECT() *MockStorageOBJMockRecorder {
	return m.recorder
}

// Reader mocks base method.
func (m *MockStorageOBJ) Reader() (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reader")
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reader indicates an expected call of Reader.
func (mr *MockStorageOBJMockRecorder) Reader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reader", reflect.TypeOf((*MockStorageOBJ)(nil).Reader))
}

// Writer mocks base method.
func (m *MockStorageOBJ) Writer() (io.WriteCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writer")
	ret0, _ := ret[0].(io.WriteCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Writer indicates an expected call of Writer.
func (mr *MockStorageOBJMockRecorder) Writer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writer", reflect.TypeOf((*MockStorageOBJ)(nil).Writer))
}
