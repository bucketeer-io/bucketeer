// Code generated by MockGen. DO NOT EDIT.
// Source: auto_ops_rule.go
//
// Generated by this command:
//
//	mockgen -source=auto_ops_rule.go -package=mock -destination=./mock/auto_ops_rule.go
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

// MockAutoOpsRuleStorage is a mock of AutoOpsRuleStorage interface.
type MockAutoOpsRuleStorage struct {
	ctrl     *gomock.Controller
	recorder *MockAutoOpsRuleStorageMockRecorder
}

// MockAutoOpsRuleStorageMockRecorder is the mock recorder for MockAutoOpsRuleStorage.
type MockAutoOpsRuleStorageMockRecorder struct {
	mock *MockAutoOpsRuleStorage
}

// NewMockAutoOpsRuleStorage creates a new mock instance.
func NewMockAutoOpsRuleStorage(ctrl *gomock.Controller) *MockAutoOpsRuleStorage {
	mock := &MockAutoOpsRuleStorage{ctrl: ctrl}
	mock.recorder = &MockAutoOpsRuleStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAutoOpsRuleStorage) EXPECT() *MockAutoOpsRuleStorageMockRecorder {
	return m.recorder
}

// CreateAutoOpsRule mocks base method.
func (m *MockAutoOpsRuleStorage) CreateAutoOpsRule(ctx context.Context, e *domain.AutoOpsRule, environmentId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAutoOpsRule", ctx, e, environmentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAutoOpsRule indicates an expected call of CreateAutoOpsRule.
func (mr *MockAutoOpsRuleStorageMockRecorder) CreateAutoOpsRule(ctx, e, environmentId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAutoOpsRule", reflect.TypeOf((*MockAutoOpsRuleStorage)(nil).CreateAutoOpsRule), ctx, e, environmentId)
}

// GetAutoOpsRule mocks base method.
func (m *MockAutoOpsRuleStorage) GetAutoOpsRule(ctx context.Context, id, environmentId string) (*domain.AutoOpsRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAutoOpsRule", ctx, id, environmentId)
	ret0, _ := ret[0].(*domain.AutoOpsRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAutoOpsRule indicates an expected call of GetAutoOpsRule.
func (mr *MockAutoOpsRuleStorageMockRecorder) GetAutoOpsRule(ctx, id, environmentId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAutoOpsRule", reflect.TypeOf((*MockAutoOpsRuleStorage)(nil).GetAutoOpsRule), ctx, id, environmentId)
}

// ListAutoOpsRules mocks base method.
func (m *MockAutoOpsRuleStorage) ListAutoOpsRules(ctx context.Context, whereParts []mysql.WherePart, orders []*mysql.Order, limit, offset int) ([]*autoops.AutoOpsRule, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAutoOpsRules", ctx, whereParts, orders, limit, offset)
	ret0, _ := ret[0].([]*autoops.AutoOpsRule)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListAutoOpsRules indicates an expected call of ListAutoOpsRules.
func (mr *MockAutoOpsRuleStorageMockRecorder) ListAutoOpsRules(ctx, whereParts, orders, limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAutoOpsRules", reflect.TypeOf((*MockAutoOpsRuleStorage)(nil).ListAutoOpsRules), ctx, whereParts, orders, limit, offset)
}

// UpdateAutoOpsRule mocks base method.
func (m *MockAutoOpsRuleStorage) UpdateAutoOpsRule(ctx context.Context, e *domain.AutoOpsRule, environmentId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAutoOpsRule", ctx, e, environmentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAutoOpsRule indicates an expected call of UpdateAutoOpsRule.
func (mr *MockAutoOpsRuleStorageMockRecorder) UpdateAutoOpsRule(ctx, e, environmentId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAutoOpsRule", reflect.TypeOf((*MockAutoOpsRuleStorage)(nil).UpdateAutoOpsRule), ctx, e, environmentId)
}
