// Code generated by MockGen. DO NOT EDIT.
// Source: experiment_result.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/bucketeer-io/bucketeer/pkg/eventcounter/domain"
)

// MockExperimentResultStorage is a mock of ExperimentResultStorage interface.
type MockExperimentResultStorage struct {
	ctrl     *gomock.Controller
	recorder *MockExperimentResultStorageMockRecorder
}

// MockExperimentResultStorageMockRecorder is the mock recorder for MockExperimentResultStorage.
type MockExperimentResultStorageMockRecorder struct {
	mock *MockExperimentResultStorage
}

// NewMockExperimentResultStorage creates a new mock instance.
func NewMockExperimentResultStorage(ctrl *gomock.Controller) *MockExperimentResultStorage {
	mock := &MockExperimentResultStorage{ctrl: ctrl}
	mock.recorder = &MockExperimentResultStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExperimentResultStorage) EXPECT() *MockExperimentResultStorageMockRecorder {
	return m.recorder
}

// GetExperimentResult mocks base method.
func (m *MockExperimentResultStorage) GetExperimentResult(ctx context.Context, id, environmentNamespace string) (*domain.ExperimentResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExperimentResult", ctx, id, environmentNamespace)
	ret0, _ := ret[0].(*domain.ExperimentResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExperimentResult indicates an expected call of GetExperimentResult.
func (mr *MockExperimentResultStorageMockRecorder) GetExperimentResult(ctx, id, environmentNamespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExperimentResult", reflect.TypeOf((*MockExperimentResultStorage)(nil).GetExperimentResult), ctx, id, environmentNamespace)
}
