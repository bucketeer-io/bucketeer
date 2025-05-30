// Code generated by MockGen. DO NOT EDIT.
// Source: admin_audit_log.go
//
// Generated by this command:
//
//	mockgen -source=admin_audit_log.go -package=mock -destination=./mock/admin_audit_log.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/bucketeer-io/bucketeer/pkg/auditlog/domain"
	mysql "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	auditlog "github.com/bucketeer-io/bucketeer/proto/auditlog"
)

// MockAdminAuditLogStorage is a mock of AdminAuditLogStorage interface.
type MockAdminAuditLogStorage struct {
	ctrl     *gomock.Controller
	recorder *MockAdminAuditLogStorageMockRecorder
}

// MockAdminAuditLogStorageMockRecorder is the mock recorder for MockAdminAuditLogStorage.
type MockAdminAuditLogStorageMockRecorder struct {
	mock *MockAdminAuditLogStorage
}

// NewMockAdminAuditLogStorage creates a new mock instance.
func NewMockAdminAuditLogStorage(ctrl *gomock.Controller) *MockAdminAuditLogStorage {
	mock := &MockAdminAuditLogStorage{ctrl: ctrl}
	mock.recorder = &MockAdminAuditLogStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdminAuditLogStorage) EXPECT() *MockAdminAuditLogStorageMockRecorder {
	return m.recorder
}

// CreateAdminAuditLog mocks base method.
func (m *MockAdminAuditLogStorage) CreateAdminAuditLog(ctx context.Context, auditLog *domain.AuditLog) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdminAuditLog", ctx, auditLog)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAdminAuditLog indicates an expected call of CreateAdminAuditLog.
func (mr *MockAdminAuditLogStorageMockRecorder) CreateAdminAuditLog(ctx, auditLog any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdminAuditLog", reflect.TypeOf((*MockAdminAuditLogStorage)(nil).CreateAdminAuditLog), ctx, auditLog)
}

// CreateAdminAuditLogs mocks base method.
func (m *MockAdminAuditLogStorage) CreateAdminAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdminAuditLogs", ctx, auditLogs)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAdminAuditLogs indicates an expected call of CreateAdminAuditLogs.
func (mr *MockAdminAuditLogStorageMockRecorder) CreateAdminAuditLogs(ctx, auditLogs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdminAuditLogs", reflect.TypeOf((*MockAdminAuditLogStorage)(nil).CreateAdminAuditLogs), ctx, auditLogs)
}

// ListAdminAuditLogs mocks base method.
func (m *MockAdminAuditLogStorage) ListAdminAuditLogs(ctx context.Context, options *mysql.ListOptions) ([]*auditlog.AuditLog, int, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAdminAuditLogs", ctx, options)
	ret0, _ := ret[0].([]*auditlog.AuditLog)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(int64)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListAdminAuditLogs indicates an expected call of ListAdminAuditLogs.
func (mr *MockAdminAuditLogStorageMockRecorder) ListAdminAuditLogs(ctx, options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAdminAuditLogs", reflect.TypeOf((*MockAdminAuditLogStorage)(nil).ListAdminAuditLogs), ctx, options)
}
