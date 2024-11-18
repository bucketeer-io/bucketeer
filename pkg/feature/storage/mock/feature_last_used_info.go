// Code generated by MockGen. DO NOT EDIT.
// Source: feature_last_used_info.go
//
// Generated by this command:
//
//	mockgen -source=feature_last_used_info.go -package=mock -destination=./mock/feature_last_used_info.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	storage "github.com/bucketeer-io/bucketeer/pkg/storage"
	feature "github.com/bucketeer-io/bucketeer/proto/feature"
)

// MockFeatureLastUsedStorage is a mock of FeatureLastUsedStorage interface.
type MockFeatureLastUsedStorage struct {
	ctrl     *gomock.Controller
	recorder *MockFeatureLastUsedStorageMockRecorder
}

// MockFeatureLastUsedStorageMockRecorder is the mock recorder for MockFeatureLastUsedStorage.
type MockFeatureLastUsedStorageMockRecorder struct {
	mock *MockFeatureLastUsedStorage
}

// NewMockFeatureLastUsedStorage creates a new mock instance.
func NewMockFeatureLastUsedStorage(ctrl *gomock.Controller) *MockFeatureLastUsedStorage {
	mock := &MockFeatureLastUsedStorage{ctrl: ctrl}
	mock.recorder = &MockFeatureLastUsedStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFeatureLastUsedStorage) EXPECT() *MockFeatureLastUsedStorageMockRecorder {
	return m.recorder
}

// GetFeatureLastUsedInfos mocks base method.
func (m *MockFeatureLastUsedStorage) GetFeatureLastUsedInfos(ctx context.Context, ids []string, environmentId string) ([]*domain.FeatureLastUsedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeatureLastUsedInfos", ctx, ids, environmentId)
	ret0, _ := ret[0].([]*domain.FeatureLastUsedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeatureLastUsedInfos indicates an expected call of GetFeatureLastUsedInfos.
func (mr *MockFeatureLastUsedStorageMockRecorder) GetFeatureLastUsedInfos(ctx, ids, environmentId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeatureLastUsedInfos", reflect.TypeOf((*MockFeatureLastUsedStorage)(nil).GetFeatureLastUsedInfos), ctx, ids, environmentId)
}

// UpsertFeatureLastUsedInfos mocks base method.
func (m *MockFeatureLastUsedStorage) UpsertFeatureLastUsedInfos(ctx context.Context, featureLastUsedInfos []*domain.FeatureLastUsedInfo, environmentId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertFeatureLastUsedInfos", ctx, featureLastUsedInfos, environmentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertFeatureLastUsedInfos indicates an expected call of UpsertFeatureLastUsedInfos.
func (mr *MockFeatureLastUsedStorageMockRecorder) UpsertFeatureLastUsedInfos(ctx, featureLastUsedInfos, environmentId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertFeatureLastUsedInfos", reflect.TypeOf((*MockFeatureLastUsedStorage)(nil).UpsertFeatureLastUsedInfos), ctx, featureLastUsedInfos, environmentId)
}

// MockFeatureLastUsedLister is a mock of FeatureLastUsedLister interface.
type MockFeatureLastUsedLister struct {
	ctrl     *gomock.Controller
	recorder *MockFeatureLastUsedListerMockRecorder
}

// MockFeatureLastUsedListerMockRecorder is the mock recorder for MockFeatureLastUsedLister.
type MockFeatureLastUsedListerMockRecorder struct {
	mock *MockFeatureLastUsedLister
}

// NewMockFeatureLastUsedLister creates a new mock instance.
func NewMockFeatureLastUsedLister(ctrl *gomock.Controller) *MockFeatureLastUsedLister {
	mock := &MockFeatureLastUsedLister{ctrl: ctrl}
	mock.recorder = &MockFeatureLastUsedListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFeatureLastUsedLister) EXPECT() *MockFeatureLastUsedListerMockRecorder {
	return m.recorder
}

// ListFeatureLastUsedInfo mocks base method.
func (m *MockFeatureLastUsedLister) ListFeatureLastUsedInfo(ctx context.Context, pageSize int, cursor, environmentId string, filters ...*storage.Filter) ([]*feature.FeatureLastUsedInfo, string, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, pageSize, cursor, environmentId}
	for _, a := range filters {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFeatureLastUsedInfo", varargs...)
	ret0, _ := ret[0].([]*feature.FeatureLastUsedInfo)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListFeatureLastUsedInfo indicates an expected call of ListFeatureLastUsedInfo.
func (mr *MockFeatureLastUsedListerMockRecorder) ListFeatureLastUsedInfo(ctx, pageSize, cursor, environmentId any, filters ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, pageSize, cursor, environmentId}, filters...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeatureLastUsedInfo", reflect.TypeOf((*MockFeatureLastUsedLister)(nil).ListFeatureLastUsedInfo), varargs...)
}
