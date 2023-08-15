// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/thumb/thumb.go

// Package mock_thumb is a generated GoMock package.
package mock_thumb

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Mockcache is a mock of cache interface.
type Mockcache struct {
	ctrl     *gomock.Controller
	recorder *MockcacheMockRecorder
}

// MockcacheMockRecorder is the mock recorder for Mockcache.
type MockcacheMockRecorder struct {
	mock *Mockcache
}

// NewMockcache creates a new mock instance.
func NewMockcache(ctrl *gomock.Controller) *Mockcache {
	mock := &Mockcache{ctrl: ctrl}
	mock.recorder = &MockcacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockcache) EXPECT() *MockcacheMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *Mockcache) Delete(key any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", key)
}

// Delete indicates an expected call of Delete.
func (mr *MockcacheMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*Mockcache)(nil).Delete), key)
}

// Load mocks base method.
func (m *Mockcache) Load(key any) (any, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", key)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockcacheMockRecorder) Load(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*Mockcache)(nil).Load), key)
}

// Range mocks base method.
func (m *Mockcache) Range(f func(any, any) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Range", f)
}

// Range indicates an expected call of Range.
func (mr *MockcacheMockRecorder) Range(f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Range", reflect.TypeOf((*Mockcache)(nil).Range), f)
}

// Store mocks base method.
func (m *Mockcache) Store(key, value any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Store", key, value)
}

// Store indicates an expected call of Store.
func (mr *MockcacheMockRecorder) Store(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*Mockcache)(nil).Store), key, value)
}
