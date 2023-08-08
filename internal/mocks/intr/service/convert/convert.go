// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/convert/convert.go

// Package mock_convert is a generated GoMock package.
package mock_convert

import (
	io "io"
	reflect "reflect"

	types "github.com/jhekau/favicon/internal/core/types"
	convert "github.com/jhekau/favicon/internal/service/convert"
	domain "github.com/jhekau/favicon/pkg/domain"
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
func (m *MockStorageOBJ) Key() domain.StorageKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Key")
	ret0, _ := ret[0].(domain.StorageKey)
	return ret0
}

// Key indicates an expected call of Key.
func (mr *MockStorageOBJMockRecorder) Key() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Key", reflect.TypeOf((*MockStorageOBJ)(nil).Key))
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

// MockConverterT is a mock of ConverterT interface.
type MockConverterT struct {
	ctrl     *gomock.Controller
	recorder *MockConverterTMockRecorder
}

// MockConverterTMockRecorder is the mock recorder for MockConverterT.
type MockConverterTMockRecorder struct {
	mock *MockConverterT
}

// NewMockConverterT creates a new mock instance.
func NewMockConverterT(ctrl *gomock.Controller) *MockConverterT {
	mock := &MockConverterT{ctrl: ctrl}
	mock.recorder = &MockConverterTMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConverterT) EXPECT() *MockConverterTMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockConverterT) Do(source, save convert.StorageOBJ, size_px int, typ types.FileType) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", source, save, size_px, typ)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockConverterTMockRecorder) Do(source, save, size_px, typ interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockConverterT)(nil).Do), source, save, size_px, typ)
}

// MockCheckPreview is a mock of CheckPreview interface.
type MockCheckPreview struct {
	ctrl     *gomock.Controller
	recorder *MockCheckPreviewMockRecorder
}

// MockCheckPreviewMockRecorder is the mock recorder for MockCheckPreview.
type MockCheckPreviewMockRecorder struct {
	mock *MockCheckPreview
}

// NewMockCheckPreview creates a new mock instance.
func NewMockCheckPreview(ctrl *gomock.Controller) *MockCheckPreview {
	mock := &MockCheckPreview{ctrl: ctrl}
	mock.recorder = &MockCheckPreviewMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckPreview) EXPECT() *MockCheckPreviewMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockCheckPreview) Check(typ types.FileType, size_px int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", typ, size_px)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockCheckPreviewMockRecorder) Check(typ, size_px interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockCheckPreview)(nil).Check), typ, size_px)
}

// MockCheckSource is a mock of CheckSource interface.
type MockCheckSource struct {
	ctrl     *gomock.Controller
	recorder *MockCheckSourceMockRecorder
}

// MockCheckSourceMockRecorder is the mock recorder for MockCheckSource.
type MockCheckSourceMockRecorder struct {
	mock *MockCheckSource
}

// NewMockCheckSource creates a new mock instance.
func NewMockCheckSource(ctrl *gomock.Controller) *MockCheckSource {
	mock := &MockCheckSource{ctrl: ctrl}
	mock.recorder = &MockCheckSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckSource) EXPECT() *MockCheckSourceMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockCheckSource) Check(original convert.StorageOBJ, originalSVG bool, thumb_size int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", original, originalSVG, thumb_size)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockCheckSourceMockRecorder) Check(original, originalSVG, thumb_size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockCheckSource)(nil).Check), original, originalSVG, thumb_size)
}
