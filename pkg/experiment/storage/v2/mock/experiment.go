// Code generated by MockGen. DO NOT EDIT.
// Source: experiment.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	mysql "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	experiment "github.com/bucketeer-io/bucketeer/proto/experiment"
)

// MockExperimentStorage is a mock of ExperimentStorage interface.
type MockExperimentStorage struct {
	ctrl     *gomock.Controller
	recorder *MockExperimentStorageMockRecorder
}

// MockExperimentStorageMockRecorder is the mock recorder for MockExperimentStorage.
type MockExperimentStorageMockRecorder struct {
	mock *MockExperimentStorage
}

// NewMockExperimentStorage creates a new mock instance.
func NewMockExperimentStorage(ctrl *gomock.Controller) *MockExperimentStorage {
	mock := &MockExperimentStorage{ctrl: ctrl}
	mock.recorder = &MockExperimentStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExperimentStorage) EXPECT() *MockExperimentStorageMockRecorder {
	return m.recorder
}

// CreateExperiment mocks base method.
func (m *MockExperimentStorage) CreateExperiment(ctx context.Context, e *domain.Experiment, environmentNamespace string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateExperiment", ctx, e, environmentNamespace)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateExperiment indicates an expected call of CreateExperiment.
func (mr *MockExperimentStorageMockRecorder) CreateExperiment(ctx, e, environmentNamespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExperiment", reflect.TypeOf((*MockExperimentStorage)(nil).CreateExperiment), ctx, e, environmentNamespace)
}

// GetExperiment mocks base method.
func (m *MockExperimentStorage) GetExperiment(ctx context.Context, id, environmentNamespace string) (*domain.Experiment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExperiment", ctx, id, environmentNamespace)
	ret0, _ := ret[0].(*domain.Experiment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExperiment indicates an expected call of GetExperiment.
func (mr *MockExperimentStorageMockRecorder) GetExperiment(ctx, id, environmentNamespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExperiment", reflect.TypeOf((*MockExperimentStorage)(nil).GetExperiment), ctx, id, environmentNamespace)
}

// ListExperiments mocks base method.
func (m *MockExperimentStorage) ListExperiments(ctx context.Context, whereParts []mysql.WherePart, orders []*mysql.Order, limit, offset int) ([]*experiment.Experiment, int, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListExperiments", ctx, whereParts, orders, limit, offset)
	ret0, _ := ret[0].([]*experiment.Experiment)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(int64)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListExperiments indicates an expected call of ListExperiments.
func (mr *MockExperimentStorageMockRecorder) ListExperiments(ctx, whereParts, orders, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListExperiments", reflect.TypeOf((*MockExperimentStorage)(nil).ListExperiments), ctx, whereParts, orders, limit, offset)
}

// UpdateExperiment mocks base method.
func (m *MockExperimentStorage) UpdateExperiment(ctx context.Context, e *domain.Experiment, environmentNamespace string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateExperiment", ctx, e, environmentNamespace)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExperiment indicates an expected call of UpdateExperiment.
func (mr *MockExperimentStorageMockRecorder) UpdateExperiment(ctx, e, environmentNamespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExperiment", reflect.TypeOf((*MockExperimentStorage)(nil).UpdateExperiment), ctx, e, environmentNamespace)
}
