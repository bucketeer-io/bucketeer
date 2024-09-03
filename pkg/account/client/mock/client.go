// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source=client.go -package=mock -destination=./mock/client.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"

	account "github.com/bucketeer-io/bucketeer/proto/account"
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

// ChangeAPIKeyName mocks base method.
func (m *MockClient) ChangeAPIKeyName(ctx context.Context, in *account.ChangeAPIKeyNameRequest, opts ...grpc.CallOption) (*account.ChangeAPIKeyNameResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChangeAPIKeyName", varargs...)
	ret0, _ := ret[0].(*account.ChangeAPIKeyNameResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeAPIKeyName indicates an expected call of ChangeAPIKeyName.
func (mr *MockClientMockRecorder) ChangeAPIKeyName(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAPIKeyName", reflect.TypeOf((*MockClient)(nil).ChangeAPIKeyName), varargs...)
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

// CreateAPIKey mocks base method.
func (m *MockClient) CreateAPIKey(ctx context.Context, in *account.CreateAPIKeyRequest, opts ...grpc.CallOption) (*account.CreateAPIKeyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAPIKey", varargs...)
	ret0, _ := ret[0].(*account.CreateAPIKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAPIKey indicates an expected call of CreateAPIKey.
func (mr *MockClientMockRecorder) CreateAPIKey(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAPIKey", reflect.TypeOf((*MockClient)(nil).CreateAPIKey), varargs...)
}

// CreateAccountV2 mocks base method.
func (m *MockClient) CreateAccountV2(ctx context.Context, in *account.CreateAccountV2Request, opts ...grpc.CallOption) (*account.CreateAccountV2Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAccountV2", varargs...)
	ret0, _ := ret[0].(*account.CreateAccountV2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccountV2 indicates an expected call of CreateAccountV2.
func (mr *MockClientMockRecorder) CreateAccountV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccountV2", reflect.TypeOf((*MockClient)(nil).CreateAccountV2), varargs...)
}

// CreateSearchFilter mocks base method.
func (m *MockClient) CreateSearchFilter(ctx context.Context, in *account.CreateSearchFilterRequest, opts ...grpc.CallOption) (*account.CreateSearchFilterResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSearchFilter", varargs...)
	ret0, _ := ret[0].(*account.CreateSearchFilterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSearchFilter indicates an expected call of CreateSearchFilter.
func (mr *MockClientMockRecorder) CreateSearchFilter(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSearchFilter", reflect.TypeOf((*MockClient)(nil).CreateSearchFilter), varargs...)
}

// DeleteAccountV2 mocks base method.
func (m *MockClient) DeleteAccountV2(ctx context.Context, in *account.DeleteAccountV2Request, opts ...grpc.CallOption) (*account.DeleteAccountV2Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAccountV2", varargs...)
	ret0, _ := ret[0].(*account.DeleteAccountV2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAccountV2 indicates an expected call of DeleteAccountV2.
func (mr *MockClientMockRecorder) DeleteAccountV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccountV2", reflect.TypeOf((*MockClient)(nil).DeleteAccountV2), varargs...)
}

// DeleteSearchFilterV2 mocks base method.
func (m *MockClient) DeleteSearchFilterV2(ctx context.Context, in *account.DeleteSearchFilterRequest, opts ...grpc.CallOption) (*account.DeleteSearchFilterResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSearchFilterV2", varargs...)
	ret0, _ := ret[0].(*account.DeleteSearchFilterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSearchFilterV2 indicates an expected call of DeleteSearchFilterV2.
func (mr *MockClientMockRecorder) DeleteSearchFilterV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSearchFilterV2", reflect.TypeOf((*MockClient)(nil).DeleteSearchFilterV2), varargs...)
}

// DisableAPIKey mocks base method.
func (m *MockClient) DisableAPIKey(ctx context.Context, in *account.DisableAPIKeyRequest, opts ...grpc.CallOption) (*account.DisableAPIKeyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DisableAPIKey", varargs...)
	ret0, _ := ret[0].(*account.DisableAPIKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableAPIKey indicates an expected call of DisableAPIKey.
func (mr *MockClientMockRecorder) DisableAPIKey(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableAPIKey", reflect.TypeOf((*MockClient)(nil).DisableAPIKey), varargs...)
}

// DisableAccountV2 mocks base method.
func (m *MockClient) DisableAccountV2(ctx context.Context, in *account.DisableAccountV2Request, opts ...grpc.CallOption) (*account.DisableAccountV2Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DisableAccountV2", varargs...)
	ret0, _ := ret[0].(*account.DisableAccountV2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableAccountV2 indicates an expected call of DisableAccountV2.
func (mr *MockClientMockRecorder) DisableAccountV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableAccountV2", reflect.TypeOf((*MockClient)(nil).DisableAccountV2), varargs...)
}

// EnableAPIKey mocks base method.
func (m *MockClient) EnableAPIKey(ctx context.Context, in *account.EnableAPIKeyRequest, opts ...grpc.CallOption) (*account.EnableAPIKeyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EnableAPIKey", varargs...)
	ret0, _ := ret[0].(*account.EnableAPIKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableAPIKey indicates an expected call of EnableAPIKey.
func (mr *MockClientMockRecorder) EnableAPIKey(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableAPIKey", reflect.TypeOf((*MockClient)(nil).EnableAPIKey), varargs...)
}

// EnableAccountV2 mocks base method.
func (m *MockClient) EnableAccountV2(ctx context.Context, in *account.EnableAccountV2Request, opts ...grpc.CallOption) (*account.EnableAccountV2Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EnableAccountV2", varargs...)
	ret0, _ := ret[0].(*account.EnableAccountV2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableAccountV2 indicates an expected call of EnableAccountV2.
func (mr *MockClientMockRecorder) EnableAccountV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableAccountV2", reflect.TypeOf((*MockClient)(nil).EnableAccountV2), varargs...)
}

// GetAPIKey mocks base method.
func (m *MockClient) GetAPIKey(ctx context.Context, in *account.GetAPIKeyRequest, opts ...grpc.CallOption) (*account.GetAPIKeyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAPIKey", varargs...)
	ret0, _ := ret[0].(*account.GetAPIKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAPIKey indicates an expected call of GetAPIKey.
func (mr *MockClientMockRecorder) GetAPIKey(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIKey", reflect.TypeOf((*MockClient)(nil).GetAPIKey), varargs...)
}

// GetAPIKeyBySearchingAllEnvironments mocks base method.
func (m *MockClient) GetAPIKeyBySearchingAllEnvironments(ctx context.Context, in *account.GetAPIKeyBySearchingAllEnvironmentsRequest, opts ...grpc.CallOption) (*account.GetAPIKeyBySearchingAllEnvironmentsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAPIKeyBySearchingAllEnvironments", varargs...)
	ret0, _ := ret[0].(*account.GetAPIKeyBySearchingAllEnvironmentsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAPIKeyBySearchingAllEnvironments indicates an expected call of GetAPIKeyBySearchingAllEnvironments.
func (mr *MockClientMockRecorder) GetAPIKeyBySearchingAllEnvironments(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAPIKeyBySearchingAllEnvironments", reflect.TypeOf((*MockClient)(nil).GetAPIKeyBySearchingAllEnvironments), varargs...)
}

// GetAccountV2 mocks base method.
func (m *MockClient) GetAccountV2(ctx context.Context, in *account.GetAccountV2Request, opts ...grpc.CallOption) (*account.GetAccountV2Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccountV2", varargs...)
	ret0, _ := ret[0].(*account.GetAccountV2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountV2 indicates an expected call of GetAccountV2.
func (mr *MockClientMockRecorder) GetAccountV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountV2", reflect.TypeOf((*MockClient)(nil).GetAccountV2), varargs...)
}

// GetAccountV2ByEnvironmentID mocks base method.
func (m *MockClient) GetAccountV2ByEnvironmentID(ctx context.Context, in *account.GetAccountV2ByEnvironmentIDRequest, opts ...grpc.CallOption) (*account.GetAccountV2ByEnvironmentIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccountV2ByEnvironmentID", varargs...)
	ret0, _ := ret[0].(*account.GetAccountV2ByEnvironmentIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountV2ByEnvironmentID indicates an expected call of GetAccountV2ByEnvironmentID.
func (mr *MockClientMockRecorder) GetAccountV2ByEnvironmentID(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountV2ByEnvironmentID", reflect.TypeOf((*MockClient)(nil).GetAccountV2ByEnvironmentID), varargs...)
}

// GetMe mocks base method.
func (m *MockClient) GetMe(ctx context.Context, in *account.GetMeRequest, opts ...grpc.CallOption) (*account.GetMeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMe", varargs...)
	ret0, _ := ret[0].(*account.GetMeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMe indicates an expected call of GetMe.
func (mr *MockClientMockRecorder) GetMe(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockClient)(nil).GetMe), varargs...)
}

// GetMyOrganizations mocks base method.
func (m *MockClient) GetMyOrganizations(ctx context.Context, in *account.GetMyOrganizationsRequest, opts ...grpc.CallOption) (*account.GetMyOrganizationsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMyOrganizations", varargs...)
	ret0, _ := ret[0].(*account.GetMyOrganizationsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyOrganizations indicates an expected call of GetMyOrganizations.
func (mr *MockClientMockRecorder) GetMyOrganizations(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyOrganizations", reflect.TypeOf((*MockClient)(nil).GetMyOrganizations), varargs...)
}

// GetMyOrganizationsByEmail mocks base method.
func (m *MockClient) GetMyOrganizationsByEmail(ctx context.Context, in *account.GetMyOrganizationsByEmailRequest, opts ...grpc.CallOption) (*account.GetMyOrganizationsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMyOrganizationsByEmail", varargs...)
	ret0, _ := ret[0].(*account.GetMyOrganizationsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyOrganizationsByEmail indicates an expected call of GetMyOrganizationsByEmail.
func (mr *MockClientMockRecorder) GetMyOrganizationsByEmail(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyOrganizationsByEmail", reflect.TypeOf((*MockClient)(nil).GetMyOrganizationsByEmail), varargs...)
}

// ListAPIKeys mocks base method.
func (m *MockClient) ListAPIKeys(ctx context.Context, in *account.ListAPIKeysRequest, opts ...grpc.CallOption) (*account.ListAPIKeysResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListAPIKeys", varargs...)
	ret0, _ := ret[0].(*account.ListAPIKeysResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAPIKeys indicates an expected call of ListAPIKeys.
func (mr *MockClientMockRecorder) ListAPIKeys(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAPIKeys", reflect.TypeOf((*MockClient)(nil).ListAPIKeys), varargs...)
}

// ListAccountsV2 mocks base method.
func (m *MockClient) ListAccountsV2(ctx context.Context, in *account.ListAccountsV2Request, opts ...grpc.CallOption) (*account.ListAccountsV2Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListAccountsV2", varargs...)
	ret0, _ := ret[0].(*account.ListAccountsV2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAccountsV2 indicates an expected call of ListAccountsV2.
func (mr *MockClientMockRecorder) ListAccountsV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAccountsV2", reflect.TypeOf((*MockClient)(nil).ListAccountsV2), varargs...)
}

// UpdateAccountV2 mocks base method.
func (m *MockClient) UpdateAccountV2(ctx context.Context, in *account.UpdateAccountV2Request, opts ...grpc.CallOption) (*account.UpdateAccountV2Response, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateAccountV2", varargs...)
	ret0, _ := ret[0].(*account.UpdateAccountV2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccountV2 indicates an expected call of UpdateAccountV2.
func (mr *MockClientMockRecorder) UpdateAccountV2(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccountV2", reflect.TypeOf((*MockClient)(nil).UpdateAccountV2), varargs...)
}

// UpdateSearchFilter mocks base method.
func (m *MockClient) UpdateSearchFilter(ctx context.Context, in *account.UpdateSearchFilterRequest, opts ...grpc.CallOption) (*account.UpdateSearchFilterResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateSearchFilter", varargs...)
	ret0, _ := ret[0].(*account.UpdateSearchFilterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSearchFilter indicates an expected call of UpdateSearchFilter.
func (mr *MockClientMockRecorder) UpdateSearchFilter(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSearchFilter", reflect.TypeOf((*MockClient)(nil).UpdateSearchFilter), varargs...)
}
