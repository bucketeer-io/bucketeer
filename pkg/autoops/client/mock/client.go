// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"

	autoops "github.com/bucketeer-io/bucketeer/proto/autoops"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// CreateAutoOpsRule mocks base method.
func (m *MockClient) CreateAutoOpsRule(ctx context.Context, in *autoops.CreateAutoOpsRuleRequest, opts ...grpc.CallOption) (*autoops.CreateAutoOpsRuleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAutoOpsRule", varargs...)
	ret0, _ := ret[0].(*autoops.CreateAutoOpsRuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAutoOpsRule indicates an expected call of CreateAutoOpsRule.
func (mr *MockClientMockRecorder) CreateAutoOpsRule(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAutoOpsRule", reflect.TypeOf((*MockClient)(nil).CreateAutoOpsRule), varargs...)
}

// CreateProgressiveRollout mocks base method.
func (m *MockClient) CreateProgressiveRollout(ctx context.Context, in *autoops.CreateProgressiveRolloutRequest, opts ...grpc.CallOption) (*autoops.CreateProgressiveRolloutResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateProgressiveRollout", varargs...)
	ret0, _ := ret[0].(*autoops.CreateProgressiveRolloutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProgressiveRollout indicates an expected call of CreateProgressiveRollout.
func (mr *MockClientMockRecorder) CreateProgressiveRollout(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProgressiveRollout", reflect.TypeOf((*MockClient)(nil).CreateProgressiveRollout), varargs...)
}

// CreateWebhook mocks base method.
func (m *MockClient) CreateWebhook(ctx context.Context, in *autoops.CreateWebhookRequest, opts ...grpc.CallOption) (*autoops.CreateWebhookResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateWebhook", varargs...)
	ret0, _ := ret[0].(*autoops.CreateWebhookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebhook indicates an expected call of CreateWebhook.
func (mr *MockClientMockRecorder) CreateWebhook(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebhook", reflect.TypeOf((*MockClient)(nil).CreateWebhook), varargs...)
}

// DeleteAutoOpsRule mocks base method.
func (m *MockClient) DeleteAutoOpsRule(ctx context.Context, in *autoops.DeleteAutoOpsRuleRequest, opts ...grpc.CallOption) (*autoops.DeleteAutoOpsRuleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAutoOpsRule", varargs...)
	ret0, _ := ret[0].(*autoops.DeleteAutoOpsRuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAutoOpsRule indicates an expected call of DeleteAutoOpsRule.
func (mr *MockClientMockRecorder) DeleteAutoOpsRule(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAutoOpsRule", reflect.TypeOf((*MockClient)(nil).DeleteAutoOpsRule), varargs...)
}

// DeleteProgressiveRollout mocks base method.
func (m *MockClient) DeleteProgressiveRollout(ctx context.Context, in *autoops.DeleteProgressiveRolloutRequest, opts ...grpc.CallOption) (*autoops.DeleteProgressiveRolloutResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteProgressiveRollout", varargs...)
	ret0, _ := ret[0].(*autoops.DeleteProgressiveRolloutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteProgressiveRollout indicates an expected call of DeleteProgressiveRollout.
func (mr *MockClientMockRecorder) DeleteProgressiveRollout(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProgressiveRollout", reflect.TypeOf((*MockClient)(nil).DeleteProgressiveRollout), varargs...)
}

// DeleteWebhook mocks base method.
func (m *MockClient) DeleteWebhook(ctx context.Context, in *autoops.DeleteWebhookRequest, opts ...grpc.CallOption) (*autoops.DeleteWebhookResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteWebhook", varargs...)
	ret0, _ := ret[0].(*autoops.DeleteWebhookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteWebhook indicates an expected call of DeleteWebhook.
func (mr *MockClientMockRecorder) DeleteWebhook(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebhook", reflect.TypeOf((*MockClient)(nil).DeleteWebhook), varargs...)
}

// ExecuteAutoOps mocks base method.
func (m *MockClient) ExecuteAutoOps(ctx context.Context, in *autoops.ExecuteAutoOpsRequest, opts ...grpc.CallOption) (*autoops.ExecuteAutoOpsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecuteAutoOps", varargs...)
	ret0, _ := ret[0].(*autoops.ExecuteAutoOpsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteAutoOps indicates an expected call of ExecuteAutoOps.
func (mr *MockClientMockRecorder) ExecuteAutoOps(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteAutoOps", reflect.TypeOf((*MockClient)(nil).ExecuteAutoOps), varargs...)
}

// ExecuteProgressiveRollout mocks base method.
func (m *MockClient) ExecuteProgressiveRollout(ctx context.Context, in *autoops.ExecuteProgressiveRolloutRequest, opts ...grpc.CallOption) (*autoops.ExecuteProgressiveRolloutResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecuteProgressiveRollout", varargs...)
	ret0, _ := ret[0].(*autoops.ExecuteProgressiveRolloutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteProgressiveRollout indicates an expected call of ExecuteProgressiveRollout.
func (mr *MockClientMockRecorder) ExecuteProgressiveRollout(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteProgressiveRollout", reflect.TypeOf((*MockClient)(nil).ExecuteProgressiveRollout), varargs...)
}

// GetAutoOpsRule mocks base method.
func (m *MockClient) GetAutoOpsRule(ctx context.Context, in *autoops.GetAutoOpsRuleRequest, opts ...grpc.CallOption) (*autoops.GetAutoOpsRuleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAutoOpsRule", varargs...)
	ret0, _ := ret[0].(*autoops.GetAutoOpsRuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAutoOpsRule indicates an expected call of GetAutoOpsRule.
func (mr *MockClientMockRecorder) GetAutoOpsRule(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAutoOpsRule", reflect.TypeOf((*MockClient)(nil).GetAutoOpsRule), varargs...)
}

// GetProgressiveRollout mocks base method.
func (m *MockClient) GetProgressiveRollout(ctx context.Context, in *autoops.GetProgressiveRolloutRequest, opts ...grpc.CallOption) (*autoops.GetProgressiveRolloutResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProgressiveRollout", varargs...)
	ret0, _ := ret[0].(*autoops.GetProgressiveRolloutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProgressiveRollout indicates an expected call of GetProgressiveRollout.
func (mr *MockClientMockRecorder) GetProgressiveRollout(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProgressiveRollout", reflect.TypeOf((*MockClient)(nil).GetProgressiveRollout), varargs...)
}

// GetWebhook mocks base method.
func (m *MockClient) GetWebhook(ctx context.Context, in *autoops.GetWebhookRequest, opts ...grpc.CallOption) (*autoops.GetWebhookResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetWebhook", varargs...)
	ret0, _ := ret[0].(*autoops.GetWebhookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWebhook indicates an expected call of GetWebhook.
func (mr *MockClientMockRecorder) GetWebhook(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWebhook", reflect.TypeOf((*MockClient)(nil).GetWebhook), varargs...)
}

// ListAutoOpsRules mocks base method.
func (m *MockClient) ListAutoOpsRules(ctx context.Context, in *autoops.ListAutoOpsRulesRequest, opts ...grpc.CallOption) (*autoops.ListAutoOpsRulesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListAutoOpsRules", varargs...)
	ret0, _ := ret[0].(*autoops.ListAutoOpsRulesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAutoOpsRules indicates an expected call of ListAutoOpsRules.
func (mr *MockClientMockRecorder) ListAutoOpsRules(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAutoOpsRules", reflect.TypeOf((*MockClient)(nil).ListAutoOpsRules), varargs...)
}

// ListOpsCounts mocks base method.
func (m *MockClient) ListOpsCounts(ctx context.Context, in *autoops.ListOpsCountsRequest, opts ...grpc.CallOption) (*autoops.ListOpsCountsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListOpsCounts", varargs...)
	ret0, _ := ret[0].(*autoops.ListOpsCountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOpsCounts indicates an expected call of ListOpsCounts.
func (mr *MockClientMockRecorder) ListOpsCounts(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOpsCounts", reflect.TypeOf((*MockClient)(nil).ListOpsCounts), varargs...)
}

// ListProgressiveRollouts mocks base method.
func (m *MockClient) ListProgressiveRollouts(ctx context.Context, in *autoops.ListProgressiveRolloutsRequest, opts ...grpc.CallOption) (*autoops.ListProgressiveRolloutsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListProgressiveRollouts", varargs...)
	ret0, _ := ret[0].(*autoops.ListProgressiveRolloutsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProgressiveRollouts indicates an expected call of ListProgressiveRollouts.
func (mr *MockClientMockRecorder) ListProgressiveRollouts(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProgressiveRollouts", reflect.TypeOf((*MockClient)(nil).ListProgressiveRollouts), varargs...)
}

// ListWebhooks mocks base method.
func (m *MockClient) ListWebhooks(ctx context.Context, in *autoops.ListWebhooksRequest, opts ...grpc.CallOption) (*autoops.ListWebhooksResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListWebhooks", varargs...)
	ret0, _ := ret[0].(*autoops.ListWebhooksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWebhooks indicates an expected call of ListWebhooks.
func (mr *MockClientMockRecorder) ListWebhooks(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWebhooks", reflect.TypeOf((*MockClient)(nil).ListWebhooks), varargs...)
}

// UpdateAutoOpsRule mocks base method.
func (m *MockClient) UpdateAutoOpsRule(ctx context.Context, in *autoops.UpdateAutoOpsRuleRequest, opts ...grpc.CallOption) (*autoops.UpdateAutoOpsRuleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateAutoOpsRule", varargs...)
	ret0, _ := ret[0].(*autoops.UpdateAutoOpsRuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAutoOpsRule indicates an expected call of UpdateAutoOpsRule.
func (mr *MockClientMockRecorder) UpdateAutoOpsRule(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAutoOpsRule", reflect.TypeOf((*MockClient)(nil).UpdateAutoOpsRule), varargs...)
}

// UpdateWebhook mocks base method.
func (m *MockClient) UpdateWebhook(ctx context.Context, in *autoops.UpdateWebhookRequest, opts ...grpc.CallOption) (*autoops.UpdateWebhookResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateWebhook", varargs...)
	ret0, _ := ret[0].(*autoops.UpdateWebhookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWebhook indicates an expected call of UpdateWebhook.
func (mr *MockClientMockRecorder) UpdateWebhook(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWebhook", reflect.TypeOf((*MockClient)(nil).UpdateWebhook), varargs...)
}
