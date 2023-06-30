// Code generated by MockGen. DO NOT EDIT.
// Source: mau.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	service "github.com/bucketeer-io/bucketeer/proto/event/service"
)

// MockMAUStorage is a mock of MAUStorage interface.
type MockMAUStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMAUStorageMockRecorder
}

// MockMAUStorageMockRecorder is the mock recorder for MockMAUStorage.
type MockMAUStorageMockRecorder struct {
	mock *MockMAUStorage
}

// NewMockMAUStorage creates a new mock instance.
func NewMockMAUStorage(ctrl *gomock.Controller) *MockMAUStorage {
	mock := &MockMAUStorage{ctrl: ctrl}
	mock.recorder = &MockMAUStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMAUStorage) EXPECT() *MockMAUStorageMockRecorder {
	return m.recorder
}

// UpsertMAU mocks base method.
func (m *MockMAUStorage) UpsertMAU(ctx context.Context, event *service.UserEvent, environmentNamespace string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertMAU", ctx, event, environmentNamespace)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertMAU indicates an expected call of UpsertMAU.
func (mr *MockMAUStorageMockRecorder) UpsertMAU(ctx, event, environmentNamespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMAU", reflect.TypeOf((*MockMAUStorage)(nil).UpsertMAU), ctx, event, environmentNamespace)
}
