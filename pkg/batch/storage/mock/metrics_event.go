// Code generated by MockGen. DO NOT EDIT.
// Source: metrics_event.go
//
// Generated by this command:
//
//	mockgen -source=metrics_event.go -package=mock -destination=./mock/metrics_event.go
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"

	client "github.com/bucketeer-io/bucketeer/proto/event/client"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// SaveGetEvaluationLatencyMetricsEvent mocks base method.
func (m *MockStorage) SaveGetEvaluationLatencyMetricsEvent(tag, status string, duration time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveGetEvaluationLatencyMetricsEvent", tag, status, duration)
}

// SaveGetEvaluationLatencyMetricsEvent indicates an expected call of SaveGetEvaluationLatencyMetricsEvent.
func (mr *MockStorageMockRecorder) SaveGetEvaluationLatencyMetricsEvent(tag, status, duration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveGetEvaluationLatencyMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveGetEvaluationLatencyMetricsEvent), tag, status, duration)
}

// SaveGetEvaluationSizeMetricsEvent mocks base method.
func (m *MockStorage) SaveGetEvaluationSizeMetricsEvent(tag, status string, sizeByte int32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveGetEvaluationSizeMetricsEvent", tag, status, sizeByte)
}

// SaveGetEvaluationSizeMetricsEvent indicates an expected call of SaveGetEvaluationSizeMetricsEvent.
func (mr *MockStorageMockRecorder) SaveGetEvaluationSizeMetricsEvent(tag, status, sizeByte any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveGetEvaluationSizeMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveGetEvaluationSizeMetricsEvent), tag, status, sizeByte)
}

// SaveInternalErrorCountMetricsEvent mocks base method.
func (m *MockStorage) SaveInternalErrorCountMetricsEvent(tag string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveInternalErrorCountMetricsEvent", tag)
}

// SaveInternalErrorCountMetricsEvent indicates an expected call of SaveInternalErrorCountMetricsEvent.
func (mr *MockStorageMockRecorder) SaveInternalErrorCountMetricsEvent(tag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveInternalErrorCountMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveInternalErrorCountMetricsEvent), tag)
}

// SaveInternalErrorMetricsEvent mocks base method.
func (m *MockStorage) SaveInternalErrorMetricsEvent(tag, sdkVersion string, api client.ApiId) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveInternalErrorMetricsEvent", tag, sdkVersion, api)
}

// SaveInternalErrorMetricsEvent indicates an expected call of SaveInternalErrorMetricsEvent.
func (mr *MockStorageMockRecorder) SaveInternalErrorMetricsEvent(tag, sdkVersion, api any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveInternalErrorMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveInternalErrorMetricsEvent), tag, sdkVersion, api)
}

// SaveInternalSdkErrorMetricsEvent mocks base method.
func (m *MockStorage) SaveInternalSdkErrorMetricsEvent(tag, sdkVersion string, api client.ApiId) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveInternalSdkErrorMetricsEvent", tag, sdkVersion, api)
}

// SaveInternalSdkErrorMetricsEvent indicates an expected call of SaveInternalSdkErrorMetricsEvent.
func (mr *MockStorageMockRecorder) SaveInternalSdkErrorMetricsEvent(tag, sdkVersion, api any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveInternalSdkErrorMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveInternalSdkErrorMetricsEvent), tag, sdkVersion, api)
}

// SaveLatencyMetricsEvent mocks base method.
func (m *MockStorage) SaveLatencyMetricsEvent(tag, status, sdkVersion string, api client.ApiId, duration time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveLatencyMetricsEvent", tag, status, sdkVersion, api, duration)
}

// SaveLatencyMetricsEvent indicates an expected call of SaveLatencyMetricsEvent.
func (mr *MockStorageMockRecorder) SaveLatencyMetricsEvent(tag, status, sdkVersion, api, duration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveLatencyMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveLatencyMetricsEvent), tag, status, sdkVersion, api, duration)
}

// SaveNetworkErrorMetricsEvent mocks base method.
func (m *MockStorage) SaveNetworkErrorMetricsEvent(tag, sdkVersion string, api client.ApiId) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveNetworkErrorMetricsEvent", tag, sdkVersion, api)
}

// SaveNetworkErrorMetricsEvent indicates an expected call of SaveNetworkErrorMetricsEvent.
func (mr *MockStorageMockRecorder) SaveNetworkErrorMetricsEvent(tag, sdkVersion, api any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNetworkErrorMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveNetworkErrorMetricsEvent), tag, sdkVersion, api)
}

// SaveSizeMetricsEvent mocks base method.
func (m *MockStorage) SaveSizeMetricsEvent(tag, status, sdkVersion string, api client.ApiId, sizeByte int32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveSizeMetricsEvent", tag, status, sdkVersion, api, sizeByte)
}

// SaveSizeMetricsEvent indicates an expected call of SaveSizeMetricsEvent.
func (mr *MockStorageMockRecorder) SaveSizeMetricsEvent(tag, status, sdkVersion, api, sizeByte any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSizeMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveSizeMetricsEvent), tag, status, sdkVersion, api, sizeByte)
}

// SaveTimeoutErrorCountMetricsEvent mocks base method.
func (m *MockStorage) SaveTimeoutErrorCountMetricsEvent(tag string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTimeoutErrorCountMetricsEvent", tag)
}

// SaveTimeoutErrorCountMetricsEvent indicates an expected call of SaveTimeoutErrorCountMetricsEvent.
func (mr *MockStorageMockRecorder) SaveTimeoutErrorCountMetricsEvent(tag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTimeoutErrorCountMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveTimeoutErrorCountMetricsEvent), tag)
}

// SaveTimeoutErrorMetricsEvent mocks base method.
func (m *MockStorage) SaveTimeoutErrorMetricsEvent(tag, sdkVersion string, api client.ApiId) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTimeoutErrorMetricsEvent", tag, sdkVersion, api)
}

// SaveTimeoutErrorMetricsEvent indicates an expected call of SaveTimeoutErrorMetricsEvent.
func (mr *MockStorageMockRecorder) SaveTimeoutErrorMetricsEvent(tag, sdkVersion, api any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTimeoutErrorMetricsEvent", reflect.TypeOf((*MockStorage)(nil).SaveTimeoutErrorMetricsEvent), tag, sdkVersion, api)
}