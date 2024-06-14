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

	feature "github.com/bucketeer-io/bucketeer/proto/feature"
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

// AddSegmentUser mocks base method.
func (m *MockClient) AddSegmentUser(ctx context.Context, in *feature.AddSegmentUserRequest, opts ...grpc.CallOption) (*feature.AddSegmentUserResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddSegmentUser", varargs...)
	ret0, _ := ret[0].(*feature.AddSegmentUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSegmentUser indicates an expected call of AddSegmentUser.
func (mr *MockClientMockRecorder) AddSegmentUser(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSegmentUser", reflect.TypeOf((*MockClient)(nil).AddSegmentUser), varargs...)
}

// ArchiveFeature mocks base method.
func (m *MockClient) ArchiveFeature(ctx context.Context, in *feature.ArchiveFeatureRequest, opts ...grpc.CallOption) (*feature.ArchiveFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ArchiveFeature", varargs...)
	ret0, _ := ret[0].(*feature.ArchiveFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArchiveFeature indicates an expected call of ArchiveFeature.
func (mr *MockClientMockRecorder) ArchiveFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveFeature", reflect.TypeOf((*MockClient)(nil).ArchiveFeature), varargs...)
}

// BulkDownloadSegmentUsers mocks base method.
func (m *MockClient) BulkDownloadSegmentUsers(ctx context.Context, in *feature.BulkDownloadSegmentUsersRequest, opts ...grpc.CallOption) (*feature.BulkDownloadSegmentUsersResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BulkDownloadSegmentUsers", varargs...)
	ret0, _ := ret[0].(*feature.BulkDownloadSegmentUsersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BulkDownloadSegmentUsers indicates an expected call of BulkDownloadSegmentUsers.
func (mr *MockClientMockRecorder) BulkDownloadSegmentUsers(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkDownloadSegmentUsers", reflect.TypeOf((*MockClient)(nil).BulkDownloadSegmentUsers), varargs...)
}

// BulkUploadSegmentUsers mocks base method.
func (m *MockClient) BulkUploadSegmentUsers(ctx context.Context, in *feature.BulkUploadSegmentUsersRequest, opts ...grpc.CallOption) (*feature.BulkUploadSegmentUsersResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BulkUploadSegmentUsers", varargs...)
	ret0, _ := ret[0].(*feature.BulkUploadSegmentUsersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BulkUploadSegmentUsers indicates an expected call of BulkUploadSegmentUsers.
func (mr *MockClientMockRecorder) BulkUploadSegmentUsers(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkUploadSegmentUsers", reflect.TypeOf((*MockClient)(nil).BulkUploadSegmentUsers), varargs...)
}

// CloneFeature mocks base method.
func (m *MockClient) CloneFeature(ctx context.Context, in *feature.CloneFeatureRequest, opts ...grpc.CallOption) (*feature.CloneFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CloneFeature", varargs...)
	ret0, _ := ret[0].(*feature.CloneFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloneFeature indicates an expected call of CloneFeature.
func (mr *MockClientMockRecorder) CloneFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloneFeature", reflect.TypeOf((*MockClient)(nil).CloneFeature), varargs...)
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

// CreateFeature mocks base method.
func (m *MockClient) CreateFeature(ctx context.Context, in *feature.CreateFeatureRequest, opts ...grpc.CallOption) (*feature.CreateFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateFeature", varargs...)
	ret0, _ := ret[0].(*feature.CreateFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFeature indicates an expected call of CreateFeature.
func (mr *MockClientMockRecorder) CreateFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeature", reflect.TypeOf((*MockClient)(nil).CreateFeature), varargs...)
}

// CreateFlagTrigger mocks base method.
func (m *MockClient) CreateFlagTrigger(ctx context.Context, in *feature.CreateFlagTriggerRequest, opts ...grpc.CallOption) (*feature.CreateFlagTriggerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateFlagTrigger", varargs...)
	ret0, _ := ret[0].(*feature.CreateFlagTriggerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFlagTrigger indicates an expected call of CreateFlagTrigger.
func (mr *MockClientMockRecorder) CreateFlagTrigger(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFlagTrigger", reflect.TypeOf((*MockClient)(nil).CreateFlagTrigger), varargs...)
}

// CreateSegment mocks base method.
func (m *MockClient) CreateSegment(ctx context.Context, in *feature.CreateSegmentRequest, opts ...grpc.CallOption) (*feature.CreateSegmentResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSegment", varargs...)
	ret0, _ := ret[0].(*feature.CreateSegmentResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSegment indicates an expected call of CreateSegment.
func (mr *MockClientMockRecorder) CreateSegment(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSegment", reflect.TypeOf((*MockClient)(nil).CreateSegment), varargs...)
}

// DeleteFeature mocks base method.
func (m *MockClient) DeleteFeature(ctx context.Context, in *feature.DeleteFeatureRequest, opts ...grpc.CallOption) (*feature.DeleteFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteFeature", varargs...)
	ret0, _ := ret[0].(*feature.DeleteFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFeature indicates an expected call of DeleteFeature.
func (mr *MockClientMockRecorder) DeleteFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFeature", reflect.TypeOf((*MockClient)(nil).DeleteFeature), varargs...)
}

// DeleteFlagTrigger mocks base method.
func (m *MockClient) DeleteFlagTrigger(ctx context.Context, in *feature.DeleteFlagTriggerRequest, opts ...grpc.CallOption) (*feature.DeleteFlagTriggerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteFlagTrigger", varargs...)
	ret0, _ := ret[0].(*feature.DeleteFlagTriggerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFlagTrigger indicates an expected call of DeleteFlagTrigger.
func (mr *MockClientMockRecorder) DeleteFlagTrigger(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFlagTrigger", reflect.TypeOf((*MockClient)(nil).DeleteFlagTrigger), varargs...)
}

// DeleteSegment mocks base method.
func (m *MockClient) DeleteSegment(ctx context.Context, in *feature.DeleteSegmentRequest, opts ...grpc.CallOption) (*feature.DeleteSegmentResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSegment", varargs...)
	ret0, _ := ret[0].(*feature.DeleteSegmentResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSegment indicates an expected call of DeleteSegment.
func (mr *MockClientMockRecorder) DeleteSegment(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegment", reflect.TypeOf((*MockClient)(nil).DeleteSegment), varargs...)
}

// DeleteSegmentUser mocks base method.
func (m *MockClient) DeleteSegmentUser(ctx context.Context, in *feature.DeleteSegmentUserRequest, opts ...grpc.CallOption) (*feature.DeleteSegmentUserResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSegmentUser", varargs...)
	ret0, _ := ret[0].(*feature.DeleteSegmentUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSegmentUser indicates an expected call of DeleteSegmentUser.
func (mr *MockClientMockRecorder) DeleteSegmentUser(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegmentUser", reflect.TypeOf((*MockClient)(nil).DeleteSegmentUser), varargs...)
}

// DisableFeature mocks base method.
func (m *MockClient) DisableFeature(ctx context.Context, in *feature.DisableFeatureRequest, opts ...grpc.CallOption) (*feature.DisableFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DisableFeature", varargs...)
	ret0, _ := ret[0].(*feature.DisableFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableFeature indicates an expected call of DisableFeature.
func (mr *MockClientMockRecorder) DisableFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableFeature", reflect.TypeOf((*MockClient)(nil).DisableFeature), varargs...)
}

// DisableFlagTrigger mocks base method.
func (m *MockClient) DisableFlagTrigger(ctx context.Context, in *feature.DisableFlagTriggerRequest, opts ...grpc.CallOption) (*feature.DisableFlagTriggerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DisableFlagTrigger", varargs...)
	ret0, _ := ret[0].(*feature.DisableFlagTriggerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableFlagTrigger indicates an expected call of DisableFlagTrigger.
func (mr *MockClientMockRecorder) DisableFlagTrigger(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableFlagTrigger", reflect.TypeOf((*MockClient)(nil).DisableFlagTrigger), varargs...)
}

// EnableFeature mocks base method.
func (m *MockClient) EnableFeature(ctx context.Context, in *feature.EnableFeatureRequest, opts ...grpc.CallOption) (*feature.EnableFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EnableFeature", varargs...)
	ret0, _ := ret[0].(*feature.EnableFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableFeature indicates an expected call of EnableFeature.
func (mr *MockClientMockRecorder) EnableFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableFeature", reflect.TypeOf((*MockClient)(nil).EnableFeature), varargs...)
}

// EnableFlagTrigger mocks base method.
func (m *MockClient) EnableFlagTrigger(ctx context.Context, in *feature.EnableFlagTriggerRequest, opts ...grpc.CallOption) (*feature.EnableFlagTriggerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EnableFlagTrigger", varargs...)
	ret0, _ := ret[0].(*feature.EnableFlagTriggerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableFlagTrigger indicates an expected call of EnableFlagTrigger.
func (mr *MockClientMockRecorder) EnableFlagTrigger(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableFlagTrigger", reflect.TypeOf((*MockClient)(nil).EnableFlagTrigger), varargs...)
}

// EvaluateFeatures mocks base method.
func (m *MockClient) EvaluateFeatures(ctx context.Context, in *feature.EvaluateFeaturesRequest, opts ...grpc.CallOption) (*feature.EvaluateFeaturesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EvaluateFeatures", varargs...)
	ret0, _ := ret[0].(*feature.EvaluateFeaturesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EvaluateFeatures indicates an expected call of EvaluateFeatures.
func (mr *MockClientMockRecorder) EvaluateFeatures(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EvaluateFeatures", reflect.TypeOf((*MockClient)(nil).EvaluateFeatures), varargs...)
}

// FlagTriggerWebhook mocks base method.
func (m *MockClient) FlagTriggerWebhook(ctx context.Context, in *feature.FlagTriggerWebhookRequest, opts ...grpc.CallOption) (*feature.FlagTriggerWebhookResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FlagTriggerWebhook", varargs...)
	ret0, _ := ret[0].(*feature.FlagTriggerWebhookResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FlagTriggerWebhook indicates an expected call of FlagTriggerWebhook.
func (mr *MockClientMockRecorder) FlagTriggerWebhook(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlagTriggerWebhook", reflect.TypeOf((*MockClient)(nil).FlagTriggerWebhook), varargs...)
}

// GetFeature mocks base method.
func (m *MockClient) GetFeature(ctx context.Context, in *feature.GetFeatureRequest, opts ...grpc.CallOption) (*feature.GetFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFeature", varargs...)
	ret0, _ := ret[0].(*feature.GetFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeature indicates an expected call of GetFeature.
func (mr *MockClientMockRecorder) GetFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeature", reflect.TypeOf((*MockClient)(nil).GetFeature), varargs...)
}

// GetFeatures mocks base method.
func (m *MockClient) GetFeatures(ctx context.Context, in *feature.GetFeaturesRequest, opts ...grpc.CallOption) (*feature.GetFeaturesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFeatures", varargs...)
	ret0, _ := ret[0].(*feature.GetFeaturesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeatures indicates an expected call of GetFeatures.
func (mr *MockClientMockRecorder) GetFeatures(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeatures", reflect.TypeOf((*MockClient)(nil).GetFeatures), varargs...)
}

// GetFlagTrigger mocks base method.
func (m *MockClient) GetFlagTrigger(ctx context.Context, in *feature.GetFlagTriggerRequest, opts ...grpc.CallOption) (*feature.GetFlagTriggerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFlagTrigger", varargs...)
	ret0, _ := ret[0].(*feature.GetFlagTriggerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlagTrigger indicates an expected call of GetFlagTrigger.
func (mr *MockClientMockRecorder) GetFlagTrigger(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlagTrigger", reflect.TypeOf((*MockClient)(nil).GetFlagTrigger), varargs...)
}

// GetSegment mocks base method.
func (m *MockClient) GetSegment(ctx context.Context, in *feature.GetSegmentRequest, opts ...grpc.CallOption) (*feature.GetSegmentResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSegment", varargs...)
	ret0, _ := ret[0].(*feature.GetSegmentResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegment indicates an expected call of GetSegment.
func (mr *MockClientMockRecorder) GetSegment(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegment", reflect.TypeOf((*MockClient)(nil).GetSegment), varargs...)
}

// GetSegmentUser mocks base method.
func (m *MockClient) GetSegmentUser(ctx context.Context, in *feature.GetSegmentUserRequest, opts ...grpc.CallOption) (*feature.GetSegmentUserResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSegmentUser", varargs...)
	ret0, _ := ret[0].(*feature.GetSegmentUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegmentUser indicates an expected call of GetSegmentUser.
func (mr *MockClientMockRecorder) GetSegmentUser(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegmentUser", reflect.TypeOf((*MockClient)(nil).GetSegmentUser), varargs...)
}

// ListEnabledFeatures mocks base method.
func (m *MockClient) ListEnabledFeatures(ctx context.Context, in *feature.ListEnabledFeaturesRequest, opts ...grpc.CallOption) (*feature.ListEnabledFeaturesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListEnabledFeatures", varargs...)
	ret0, _ := ret[0].(*feature.ListEnabledFeaturesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEnabledFeatures indicates an expected call of ListEnabledFeatures.
func (mr *MockClientMockRecorder) ListEnabledFeatures(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEnabledFeatures", reflect.TypeOf((*MockClient)(nil).ListEnabledFeatures), varargs...)
}

// ListFeatures mocks base method.
func (m *MockClient) ListFeatures(ctx context.Context, in *feature.ListFeaturesRequest, opts ...grpc.CallOption) (*feature.ListFeaturesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFeatures", varargs...)
	ret0, _ := ret[0].(*feature.ListFeaturesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeatures indicates an expected call of ListFeatures.
func (mr *MockClientMockRecorder) ListFeatures(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeatures", reflect.TypeOf((*MockClient)(nil).ListFeatures), varargs...)
}

// ListFlagTriggers mocks base method.
func (m *MockClient) ListFlagTriggers(ctx context.Context, in *feature.ListFlagTriggersRequest, opts ...grpc.CallOption) (*feature.ListFlagTriggersResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFlagTriggers", varargs...)
	ret0, _ := ret[0].(*feature.ListFlagTriggersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFlagTriggers indicates an expected call of ListFlagTriggers.
func (mr *MockClientMockRecorder) ListFlagTriggers(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFlagTriggers", reflect.TypeOf((*MockClient)(nil).ListFlagTriggers), varargs...)
}

// ListSegmentUsers mocks base method.
func (m *MockClient) ListSegmentUsers(ctx context.Context, in *feature.ListSegmentUsersRequest, opts ...grpc.CallOption) (*feature.ListSegmentUsersResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListSegmentUsers", varargs...)
	ret0, _ := ret[0].(*feature.ListSegmentUsersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSegmentUsers indicates an expected call of ListSegmentUsers.
func (mr *MockClientMockRecorder) ListSegmentUsers(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSegmentUsers", reflect.TypeOf((*MockClient)(nil).ListSegmentUsers), varargs...)
}

// ListSegments mocks base method.
func (m *MockClient) ListSegments(ctx context.Context, in *feature.ListSegmentsRequest, opts ...grpc.CallOption) (*feature.ListSegmentsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListSegments", varargs...)
	ret0, _ := ret[0].(*feature.ListSegmentsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSegments indicates an expected call of ListSegments.
func (mr *MockClientMockRecorder) ListSegments(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSegments", reflect.TypeOf((*MockClient)(nil).ListSegments), varargs...)
}

// ListTags mocks base method.
func (m *MockClient) ListTags(ctx context.Context, in *feature.ListTagsRequest, opts ...grpc.CallOption) (*feature.ListTagsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListTags", varargs...)
	ret0, _ := ret[0].(*feature.ListTagsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTags indicates an expected call of ListTags.
func (mr *MockClientMockRecorder) ListTags(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTags", reflect.TypeOf((*MockClient)(nil).ListTags), varargs...)
}

// ResetFlagTrigger mocks base method.
func (m *MockClient) ResetFlagTrigger(ctx context.Context, in *feature.ResetFlagTriggerRequest, opts ...grpc.CallOption) (*feature.ResetFlagTriggerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ResetFlagTrigger", varargs...)
	ret0, _ := ret[0].(*feature.ResetFlagTriggerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetFlagTrigger indicates an expected call of ResetFlagTrigger.
func (mr *MockClientMockRecorder) ResetFlagTrigger(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetFlagTrigger", reflect.TypeOf((*MockClient)(nil).ResetFlagTrigger), varargs...)
}

// UnarchiveFeature mocks base method.
func (m *MockClient) UnarchiveFeature(ctx context.Context, in *feature.UnarchiveFeatureRequest, opts ...grpc.CallOption) (*feature.UnarchiveFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UnarchiveFeature", varargs...)
	ret0, _ := ret[0].(*feature.UnarchiveFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnarchiveFeature indicates an expected call of UnarchiveFeature.
func (mr *MockClientMockRecorder) UnarchiveFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnarchiveFeature", reflect.TypeOf((*MockClient)(nil).UnarchiveFeature), varargs...)
}

// UpdateFeature mocks base method.
func (m *MockClient) UpdateFeature(ctx context.Context, in *feature.UpdateFeatureRequest, opts ...grpc.CallOption) (*feature.UpdateFeatureResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateFeature", varargs...)
	ret0, _ := ret[0].(*feature.UpdateFeatureResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFeature indicates an expected call of UpdateFeature.
func (mr *MockClientMockRecorder) UpdateFeature(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeature", reflect.TypeOf((*MockClient)(nil).UpdateFeature), varargs...)
}

// UpdateFeatureDetails mocks base method.
func (m *MockClient) UpdateFeatureDetails(ctx context.Context, in *feature.UpdateFeatureDetailsRequest, opts ...grpc.CallOption) (*feature.UpdateFeatureDetailsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateFeatureDetails", varargs...)
	ret0, _ := ret[0].(*feature.UpdateFeatureDetailsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFeatureDetails indicates an expected call of UpdateFeatureDetails.
func (mr *MockClientMockRecorder) UpdateFeatureDetails(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeatureDetails", reflect.TypeOf((*MockClient)(nil).UpdateFeatureDetails), varargs...)
}

// UpdateFeatureTargeting mocks base method.
func (m *MockClient) UpdateFeatureTargeting(ctx context.Context, in *feature.UpdateFeatureTargetingRequest, opts ...grpc.CallOption) (*feature.UpdateFeatureTargetingResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateFeatureTargeting", varargs...)
	ret0, _ := ret[0].(*feature.UpdateFeatureTargetingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFeatureTargeting indicates an expected call of UpdateFeatureTargeting.
func (mr *MockClientMockRecorder) UpdateFeatureTargeting(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeatureTargeting", reflect.TypeOf((*MockClient)(nil).UpdateFeatureTargeting), varargs...)
}

// UpdateFeatureVariations mocks base method.
func (m *MockClient) UpdateFeatureVariations(ctx context.Context, in *feature.UpdateFeatureVariationsRequest, opts ...grpc.CallOption) (*feature.UpdateFeatureVariationsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateFeatureVariations", varargs...)
	ret0, _ := ret[0].(*feature.UpdateFeatureVariationsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFeatureVariations indicates an expected call of UpdateFeatureVariations.
func (mr *MockClientMockRecorder) UpdateFeatureVariations(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeatureVariations", reflect.TypeOf((*MockClient)(nil).UpdateFeatureVariations), varargs...)
}

// UpdateFlagTrigger mocks base method.
func (m *MockClient) UpdateFlagTrigger(ctx context.Context, in *feature.UpdateFlagTriggerRequest, opts ...grpc.CallOption) (*feature.UpdateFlagTriggerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateFlagTrigger", varargs...)
	ret0, _ := ret[0].(*feature.UpdateFlagTriggerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFlagTrigger indicates an expected call of UpdateFlagTrigger.
func (mr *MockClientMockRecorder) UpdateFlagTrigger(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlagTrigger", reflect.TypeOf((*MockClient)(nil).UpdateFlagTrigger), varargs...)
}

// UpdateSegment mocks base method.
func (m *MockClient) UpdateSegment(ctx context.Context, in *feature.UpdateSegmentRequest, opts ...grpc.CallOption) (*feature.UpdateSegmentResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateSegment", varargs...)
	ret0, _ := ret[0].(*feature.UpdateSegmentResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSegment indicates an expected call of UpdateSegment.
func (mr *MockClientMockRecorder) UpdateSegment(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSegment", reflect.TypeOf((*MockClient)(nil).UpdateSegment), varargs...)
}
