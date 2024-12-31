// Code generated by MockGen. DO NOT EDIT.
// Source: eventcounter.go
//
// Generated by this command:
//
//	mockgen -source=eventcounter.go -package=mock -destination=./mock/eventcounter.go
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockEventCounterCache is a mock of EventCounterCache interface.
type MockEventCounterCache struct {
	ctrl     *gomock.Controller
	recorder *MockEventCounterCacheMockRecorder
}

// MockEventCounterCacheMockRecorder is the mock recorder for MockEventCounterCache.
type MockEventCounterCacheMockRecorder struct {
	mock *MockEventCounterCache
}

// NewMockEventCounterCache creates a new mock instance.
func NewMockEventCounterCache(ctrl *gomock.Controller) *MockEventCounterCache {
	mock := &MockEventCounterCache{ctrl: ctrl}
	mock.recorder = &MockEventCounterCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventCounterCache) EXPECT() *MockEventCounterCacheMockRecorder {
	return m.recorder
}

// DeleteKey mocks base method.
func (m *MockEventCounterCache) DeleteKey(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteKey", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteKey indicates an expected call of DeleteKey.
func (mr *MockEventCounterCacheMockRecorder) DeleteKey(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKey", reflect.TypeOf((*MockEventCounterCache)(nil).DeleteKey), key)
}

// GetEventCounts mocks base method.
func (m *MockEventCounterCache) GetEventCounts(keys []string) ([]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventCounts", keys)
	ret0, _ := ret[0].([]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventCounts indicates an expected call of GetEventCounts.
func (mr *MockEventCounterCacheMockRecorder) GetEventCounts(keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventCounts", reflect.TypeOf((*MockEventCounterCache)(nil).GetEventCounts), keys)
}

// GetEventCountsV2 mocks base method.
func (m *MockEventCounterCache) GetEventCountsV2(keys [][]string) ([]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventCountsV2", keys)
	ret0, _ := ret[0].([]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventCountsV2 indicates an expected call of GetEventCountsV2.
func (mr *MockEventCounterCacheMockRecorder) GetEventCountsV2(keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventCountsV2", reflect.TypeOf((*MockEventCounterCache)(nil).GetEventCountsV2), keys)
}

// GetUserCount mocks base method.
func (m *MockEventCounterCache) GetUserCount(key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCount", key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCount indicates an expected call of GetUserCount.
func (mr *MockEventCounterCacheMockRecorder) GetUserCount(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCount", reflect.TypeOf((*MockEventCounterCache)(nil).GetUserCount), key)
}

// GetUserCounts mocks base method.
func (m *MockEventCounterCache) GetUserCounts(keys []string) ([]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCounts", keys)
	ret0, _ := ret[0].([]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCounts indicates an expected call of GetUserCounts.
func (mr *MockEventCounterCacheMockRecorder) GetUserCounts(keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCounts", reflect.TypeOf((*MockEventCounterCache)(nil).GetUserCounts), keys)
}

// GetUserCountsV2 mocks base method.
func (m *MockEventCounterCache) GetUserCountsV2(keys []string) ([]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserCountsV2", keys)
	ret0, _ := ret[0].([]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserCountsV2 indicates an expected call of GetUserCountsV2.
func (mr *MockEventCounterCacheMockRecorder) GetUserCountsV2(keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserCountsV2", reflect.TypeOf((*MockEventCounterCache)(nil).GetUserCountsV2), keys)
}

// MergeMultiKeys mocks base method.
func (m *MockEventCounterCache) MergeMultiKeys(dest string, keys []string, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MergeMultiKeys", dest, keys, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// MergeMultiKeys indicates an expected call of MergeMultiKeys.
func (mr *MockEventCounterCacheMockRecorder) MergeMultiKeys(dest, keys, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MergeMultiKeys", reflect.TypeOf((*MockEventCounterCache)(nil).MergeMultiKeys), dest, keys, expiration)
}

// UpdateUserCount mocks base method.
func (m *MockEventCounterCache) UpdateUserCount(key, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserCount", key, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserCount indicates an expected call of UpdateUserCount.
func (mr *MockEventCounterCacheMockRecorder) UpdateUserCount(key, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserCount", reflect.TypeOf((*MockEventCounterCache)(nil).UpdateUserCount), key, userID)
}
