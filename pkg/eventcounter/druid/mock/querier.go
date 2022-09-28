// Code generated by MockGen. DO NOT EDIT.
// Source: querier.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"

	eventcounter "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// QueryCount mocks base method.
func (m *MockQuerier) QueryCount(ctx context.Context, environmentNamespace string, startAt, endAt time.Time, goalID, featureID string, featureVersion int32, reason string, segmnets []string, filters []*eventcounter.Filter) (*eventcounter.Row, []*eventcounter.Row, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryCount", ctx, environmentNamespace, startAt, endAt, goalID, featureID, featureVersion, reason, segmnets, filters)
	ret0, _ := ret[0].(*eventcounter.Row)
	ret1, _ := ret[1].([]*eventcounter.Row)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryCount indicates an expected call of QueryCount.
func (mr *MockQuerierMockRecorder) QueryCount(ctx, environmentNamespace, startAt, endAt, goalID, featureID, featureVersion, reason, segmnets, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryCount", reflect.TypeOf((*MockQuerier)(nil).QueryCount), ctx, environmentNamespace, startAt, endAt, goalID, featureID, featureVersion, reason, segmnets, filters)
}

// QueryEvaluationCount mocks base method.
func (m *MockQuerier) QueryEvaluationCount(ctx context.Context, environmentNamespace string, startAt, endAt time.Time, featureID string, featureVersion int32, reason string, segmnets []string, filters []*eventcounter.Filter) (*eventcounter.Row, []*eventcounter.Row, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryEvaluationCount", ctx, environmentNamespace, startAt, endAt, featureID, featureVersion, reason, segmnets, filters)
	ret0, _ := ret[0].(*eventcounter.Row)
	ret1, _ := ret[1].([]*eventcounter.Row)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryEvaluationCount indicates an expected call of QueryEvaluationCount.
func (mr *MockQuerierMockRecorder) QueryEvaluationCount(ctx, environmentNamespace, startAt, endAt, featureID, featureVersion, reason, segmnets, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryEvaluationCount", reflect.TypeOf((*MockQuerier)(nil).QueryEvaluationCount), ctx, environmentNamespace, startAt, endAt, featureID, featureVersion, reason, segmnets, filters)
}

// QueryEvaluationTimeseriesCount mocks base method.
func (m *MockQuerier) QueryEvaluationTimeseriesCount(ctx context.Context, environmentNamespace string, startAt, endAt time.Time, featureID string, featureVersion int32, variationID string) (map[string]*eventcounter.VariationTimeseries, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryEvaluationTimeseriesCount", ctx, environmentNamespace, startAt, endAt, featureID, featureVersion, variationID)
	ret0, _ := ret[0].(map[string]*eventcounter.VariationTimeseries)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryEvaluationTimeseriesCount indicates an expected call of QueryEvaluationTimeseriesCount.
func (mr *MockQuerierMockRecorder) QueryEvaluationTimeseriesCount(ctx, environmentNamespace, startAt, endAt, featureID, featureVersion, variationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryEvaluationTimeseriesCount", reflect.TypeOf((*MockQuerier)(nil).QueryEvaluationTimeseriesCount), ctx, environmentNamespace, startAt, endAt, featureID, featureVersion, variationID)
}

// QueryGoalCount mocks base method.
func (m *MockQuerier) QueryGoalCount(ctx context.Context, environmentNamespace string, startAt, endAt time.Time, goalID, featureID string, featureVersion int32, reason string, segmnets []string, filters []*eventcounter.Filter) (*eventcounter.Row, []*eventcounter.Row, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryGoalCount", ctx, environmentNamespace, startAt, endAt, goalID, featureID, featureVersion, reason, segmnets, filters)
	ret0, _ := ret[0].(*eventcounter.Row)
	ret1, _ := ret[1].([]*eventcounter.Row)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryGoalCount indicates an expected call of QueryGoalCount.
func (mr *MockQuerierMockRecorder) QueryGoalCount(ctx, environmentNamespace, startAt, endAt, goalID, featureID, featureVersion, reason, segmnets, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryGoalCount", reflect.TypeOf((*MockQuerier)(nil).QueryGoalCount), ctx, environmentNamespace, startAt, endAt, goalID, featureID, featureVersion, reason, segmnets, filters)
}

// QuerySegmentMetadata mocks base method.
func (m *MockQuerier) QuerySegmentMetadata(ctx context.Context, environmentNamespace, dataType string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QuerySegmentMetadata", ctx, environmentNamespace, dataType)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QuerySegmentMetadata indicates an expected call of QuerySegmentMetadata.
func (mr *MockQuerierMockRecorder) QuerySegmentMetadata(ctx, environmentNamespace, dataType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QuerySegmentMetadata", reflect.TypeOf((*MockQuerier)(nil).QuerySegmentMetadata), ctx, environmentNamespace, dataType)
}

// QueryUserCount mocks base method.
func (m *MockQuerier) QueryUserCount(ctx context.Context, environmentNamespace string, startAt, endAt time.Time) (*eventcounter.Row, []*eventcounter.Row, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryUserCount", ctx, environmentNamespace, startAt, endAt)
	ret0, _ := ret[0].(*eventcounter.Row)
	ret1, _ := ret[1].([]*eventcounter.Row)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryUserCount indicates an expected call of QueryUserCount.
func (mr *MockQuerierMockRecorder) QueryUserCount(ctx, environmentNamespace, startAt, endAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryUserCount", reflect.TypeOf((*MockQuerier)(nil).QueryUserCount), ctx, environmentNamespace, startAt, endAt)
}
