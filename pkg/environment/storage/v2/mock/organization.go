// Code generated by MockGen. DO NOT EDIT.
// Source: organization.go
//
// Generated by this command:
//
//	mockgen -source=organization.go -package=mock -destination=./mock/organization.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	domain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	mysql "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	environment "github.com/bucketeer-io/bucketeer/proto/environment"
)

// MockOrganizationStorage is a mock of OrganizationStorage interface.
type MockOrganizationStorage struct {
	ctrl     *gomock.Controller
	recorder *MockOrganizationStorageMockRecorder
}

// MockOrganizationStorageMockRecorder is the mock recorder for MockOrganizationStorage.
type MockOrganizationStorageMockRecorder struct {
	mock *MockOrganizationStorage
}

// NewMockOrganizationStorage creates a new mock instance.
func NewMockOrganizationStorage(ctrl *gomock.Controller) *MockOrganizationStorage {
	mock := &MockOrganizationStorage{ctrl: ctrl}
	mock.recorder = &MockOrganizationStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrganizationStorage) EXPECT() *MockOrganizationStorageMockRecorder {
	return m.recorder
}

// CreateOrganization mocks base method.
func (m *MockOrganizationStorage) CreateOrganization(ctx context.Context, p *domain.Organization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrganization", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrganization indicates an expected call of CreateOrganization.
func (mr *MockOrganizationStorageMockRecorder) CreateOrganization(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrganization", reflect.TypeOf((*MockOrganizationStorage)(nil).CreateOrganization), ctx, p)
}

// GetOrganization mocks base method.
func (m *MockOrganizationStorage) GetOrganization(ctx context.Context, id string) (*domain.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganization", ctx, id)
	ret0, _ := ret[0].(*domain.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganization indicates an expected call of GetOrganization.
func (mr *MockOrganizationStorageMockRecorder) GetOrganization(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganization", reflect.TypeOf((*MockOrganizationStorage)(nil).GetOrganization), ctx, id)
}

// GetSystemAdminOrganization mocks base method.
func (m *MockOrganizationStorage) GetSystemAdminOrganization(ctx context.Context) (*domain.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSystemAdminOrganization", ctx)
	ret0, _ := ret[0].(*domain.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystemAdminOrganization indicates an expected call of GetSystemAdminOrganization.
func (mr *MockOrganizationStorageMockRecorder) GetSystemAdminOrganization(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystemAdminOrganization", reflect.TypeOf((*MockOrganizationStorage)(nil).GetSystemAdminOrganization), ctx)
}

// ListOrganizations mocks base method.
func (m *MockOrganizationStorage) ListOrganizations(ctx context.Context, options *mysql.ListOptions) ([]*environment.Organization, int, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrganizations", ctx, options)
	ret0, _ := ret[0].([]*environment.Organization)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(int64)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListOrganizations indicates an expected call of ListOrganizations.
func (mr *MockOrganizationStorageMockRecorder) ListOrganizations(ctx, options any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrganizations", reflect.TypeOf((*MockOrganizationStorage)(nil).ListOrganizations), ctx, options)
}

// UpdateOrganization mocks base method.
func (m *MockOrganizationStorage) UpdateOrganization(ctx context.Context, p *domain.Organization) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrganization", ctx, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrganization indicates an expected call of UpdateOrganization.
func (mr *MockOrganizationStorageMockRecorder) UpdateOrganization(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrganization", reflect.TypeOf((*MockOrganizationStorage)(nil).UpdateOrganization), ctx, p)
}
