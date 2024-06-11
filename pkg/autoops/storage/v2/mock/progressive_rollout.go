// Code generated by MockGen. DO NOT EDIT.
// Source: progressive_rollout.go
//
// Generated by this command:
//
//	mockgen -source=progressive_rollout.go -package=mock -destination=./mock/progressive_rollout.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	mysql "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	autoops "github.com/bucketeer-io/bucketeer/proto/autoops"
)

// MockProgressiveRolloutStorage is a mock of ProgressiveRolloutStorage interface.
type MockProgressiveRolloutStorage struct {
	ctrl     *gomock.Controller
	recorder *MockProgressiveRolloutStorageMockRecorder
}

// MockProgressiveRolloutStorageMockRecorder is the mock recorder for MockProgressiveRolloutStorage.
type MockProgressiveRolloutStorageMockRecorder struct {
	mock *MockProgressiveRolloutStorage
}

// NewMockProgressiveRolloutStorage creates a new mock instance.
func NewMockProgressiveRolloutStorage(ctrl *gomock.Controller) *MockProgressiveRolloutStorage {
	mock := &MockProgressiveRolloutStorage{ctrl: ctrl}
	mock.recorder = &MockProgressiveRolloutStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProgressiveRolloutStorage) EXPECT() *MockProgressiveRolloutStorageMockRecorder {
	return m.recorder
}

// CreateProgressiveRollout mocks base method.
func (m *MockProgressiveRolloutStorage) CreateProgressiveRollout(ctx context.Context, progressiveRollout *domain.ProgressiveRollout, environmentNamespace string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProgressiveRollout", ctx, progressiveRollout, environmentNamespace)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProgressiveRollout indicates an expected call of CreateProgressiveRollout.
func (mr *MockProgressiveRolloutStorageMockRecorder) CreateProgressiveRollout(ctx, progressiveRollout, environmentNamespace any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProgressiveRollout", reflect.TypeOf((*MockProgressiveRolloutStorage)(nil).CreateProgressiveRollout), ctx, progressiveRollout, environmentNamespace)
}

// DeleteProgressiveRollout mocks base method.
func (m *MockProgressiveRolloutStorage) DeleteProgressiveRollout(ctx context.Context, id, environmentNamespace string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProgressiveRollout", ctx, id, environmentNamespace)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProgressiveRollout indicates an expected call of DeleteProgressiveRollout.
func (mr *MockProgressiveRolloutStorageMockRecorder) DeleteProgressiveRollout(ctx, id, environmentNamespace any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProgressiveRollout", reflect.TypeOf((*MockProgressiveRolloutStorage)(nil).DeleteProgressiveRollout), ctx, id, environmentNamespace)
}

// GetProgressiveRollout mocks base method.
func (m *MockProgressiveRolloutStorage) GetProgressiveRollout(ctx context.Context, id, environmentNamespace string) (*domain.ProgressiveRollout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProgressiveRollout", ctx, id, environmentNamespace)
	ret0, _ := ret[0].(*domain.ProgressiveRollout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProgressiveRollout indicates an expected call of GetProgressiveRollout.
func (mr *MockProgressiveRolloutStorageMockRecorder) GetProgressiveRollout(ctx, id, environmentNamespace any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProgressiveRollout", reflect.TypeOf((*MockProgressiveRolloutStorage)(nil).GetProgressiveRollout), ctx, id, environmentNamespace)
}

// ListProgressiveRollouts mocks base method.
func (m *MockProgressiveRolloutStorage) ListProgressiveRollouts(ctx context.Context, whereParts []mysql.WherePart, orders []*mysql.Order, limit, offset int) ([]*autoops.ProgressiveRollout, int64, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProgressiveRollouts", ctx, whereParts, orders, limit, offset)
	ret0, _ := ret[0].([]*autoops.ProgressiveRollout)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListProgressiveRollouts indicates an expected call of ListProgressiveRollouts.
func (mr *MockProgressiveRolloutStorageMockRecorder) ListProgressiveRollouts(ctx, whereParts, orders, limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProgressiveRollouts", reflect.TypeOf((*MockProgressiveRolloutStorage)(nil).ListProgressiveRollouts), ctx, whereParts, orders, limit, offset)
}

// UpdateProgressiveRollout mocks base method.
func (m *MockProgressiveRolloutStorage) UpdateProgressiveRollout(ctx context.Context, progressiveRollout *domain.ProgressiveRollout, environmentNamespace string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProgressiveRollout", ctx, progressiveRollout, environmentNamespace)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProgressiveRollout indicates an expected call of UpdateProgressiveRollout.
func (mr *MockProgressiveRolloutStorageMockRecorder) UpdateProgressiveRollout(ctx, progressiveRollout, environmentNamespace any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProgressiveRollout", reflect.TypeOf((*MockProgressiveRolloutStorage)(nil).UpdateProgressiveRollout), ctx, progressiveRollout, environmentNamespace)
}
