// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces/storage/storage.obj.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	io "io"
	reflect "reflect"
	time "time"

	storage "github.com/jhekau/favicon/interfaces/storage"
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

// IsExists mocks base method.
func (m *MockStorageOBJ) IsExists() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExists")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExists indicates an expected call of IsExists.
func (mr *MockStorageOBJMockRecorder) IsExists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExists", reflect.TypeOf((*MockStorageOBJ)(nil).IsExists))
}

// Key mocks base method.
func (m *MockStorageOBJ) Key() storage.StorageKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Key")
	ret0, _ := ret[0].(storage.StorageKey)
	return ret0
}

// Key indicates an expected call of Key.
func (mr *MockStorageOBJMockRecorder) Key() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Key", reflect.TypeOf((*MockStorageOBJ)(nil).Key))
}

// ModTime mocks base method.
func (m *MockStorageOBJ) ModTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// ModTime indicates an expected call of ModTime.
func (mr *MockStorageOBJMockRecorder) ModTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModTime", reflect.TypeOf((*MockStorageOBJ)(nil).ModTime))
}

// Reader mocks base method.
func (m *MockStorageOBJ) Reader() (io.ReadSeekCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reader")
	ret0, _ := ret[0].(io.ReadSeekCloser)
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