// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"

	notification "github.com/bucketeer-io/bucketeer/proto/notification"
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

// CreateAdminSubscription mocks base method.
func (m *MockClient) CreateAdminSubscription(ctx context.Context, in *notification.CreateAdminSubscriptionRequest, opts ...grpc.CallOption) (*notification.CreateAdminSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAdminSubscription", varargs...)
	ret0, _ := ret[0].(*notification.CreateAdminSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdminSubscription indicates an expected call of CreateAdminSubscription.
func (mr *MockClientMockRecorder) CreateAdminSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdminSubscription", reflect.TypeOf((*MockClient)(nil).CreateAdminSubscription), varargs...)
}

// CreateSubscription mocks base method.
func (m *MockClient) CreateSubscription(ctx context.Context, in *notification.CreateSubscriptionRequest, opts ...grpc.CallOption) (*notification.CreateSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSubscription", varargs...)
	ret0, _ := ret[0].(*notification.CreateSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubscription indicates an expected call of CreateSubscription.
func (mr *MockClientMockRecorder) CreateSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubscription", reflect.TypeOf((*MockClient)(nil).CreateSubscription), varargs...)
}

// DeleteAdminSubscription mocks base method.
func (m *MockClient) DeleteAdminSubscription(ctx context.Context, in *notification.DeleteAdminSubscriptionRequest, opts ...grpc.CallOption) (*notification.DeleteAdminSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAdminSubscription", varargs...)
	ret0, _ := ret[0].(*notification.DeleteAdminSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAdminSubscription indicates an expected call of DeleteAdminSubscription.
func (mr *MockClientMockRecorder) DeleteAdminSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdminSubscription", reflect.TypeOf((*MockClient)(nil).DeleteAdminSubscription), varargs...)
}

// DeleteSubscription mocks base method.
func (m *MockClient) DeleteSubscription(ctx context.Context, in *notification.DeleteSubscriptionRequest, opts ...grpc.CallOption) (*notification.DeleteSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSubscription", varargs...)
	ret0, _ := ret[0].(*notification.DeleteSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSubscription indicates an expected call of DeleteSubscription.
func (mr *MockClientMockRecorder) DeleteSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscription", reflect.TypeOf((*MockClient)(nil).DeleteSubscription), varargs...)
}

// DisableAdminSubscription mocks base method.
func (m *MockClient) DisableAdminSubscription(ctx context.Context, in *notification.DisableAdminSubscriptionRequest, opts ...grpc.CallOption) (*notification.DisableAdminSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DisableAdminSubscription", varargs...)
	ret0, _ := ret[0].(*notification.DisableAdminSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableAdminSubscription indicates an expected call of DisableAdminSubscription.
func (mr *MockClientMockRecorder) DisableAdminSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableAdminSubscription", reflect.TypeOf((*MockClient)(nil).DisableAdminSubscription), varargs...)
}

// DisableSubscription mocks base method.
func (m *MockClient) DisableSubscription(ctx context.Context, in *notification.DisableSubscriptionRequest, opts ...grpc.CallOption) (*notification.DisableSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DisableSubscription", varargs...)
	ret0, _ := ret[0].(*notification.DisableSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableSubscription indicates an expected call of DisableSubscription.
func (mr *MockClientMockRecorder) DisableSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableSubscription", reflect.TypeOf((*MockClient)(nil).DisableSubscription), varargs...)
}

// EnableAdminSubscription mocks base method.
func (m *MockClient) EnableAdminSubscription(ctx context.Context, in *notification.EnableAdminSubscriptionRequest, opts ...grpc.CallOption) (*notification.EnableAdminSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EnableAdminSubscription", varargs...)
	ret0, _ := ret[0].(*notification.EnableAdminSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableAdminSubscription indicates an expected call of EnableAdminSubscription.
func (mr *MockClientMockRecorder) EnableAdminSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableAdminSubscription", reflect.TypeOf((*MockClient)(nil).EnableAdminSubscription), varargs...)
}

// EnableSubscription mocks base method.
func (m *MockClient) EnableSubscription(ctx context.Context, in *notification.EnableSubscriptionRequest, opts ...grpc.CallOption) (*notification.EnableSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EnableSubscription", varargs...)
	ret0, _ := ret[0].(*notification.EnableSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableSubscription indicates an expected call of EnableSubscription.
func (mr *MockClientMockRecorder) EnableSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableSubscription", reflect.TypeOf((*MockClient)(nil).EnableSubscription), varargs...)
}

// GetAdminSubscription mocks base method.
func (m *MockClient) GetAdminSubscription(ctx context.Context, in *notification.GetAdminSubscriptionRequest, opts ...grpc.CallOption) (*notification.GetAdminSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAdminSubscription", varargs...)
	ret0, _ := ret[0].(*notification.GetAdminSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminSubscription indicates an expected call of GetAdminSubscription.
func (mr *MockClientMockRecorder) GetAdminSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminSubscription", reflect.TypeOf((*MockClient)(nil).GetAdminSubscription), varargs...)
}

// GetSubscription mocks base method.
func (m *MockClient) GetSubscription(ctx context.Context, in *notification.GetSubscriptionRequest, opts ...grpc.CallOption) (*notification.GetSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSubscription", varargs...)
	ret0, _ := ret[0].(*notification.GetSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscription indicates an expected call of GetSubscription.
func (mr *MockClientMockRecorder) GetSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscription", reflect.TypeOf((*MockClient)(nil).GetSubscription), varargs...)
}

// ListAdminSubscriptions mocks base method.
func (m *MockClient) ListAdminSubscriptions(ctx context.Context, in *notification.ListAdminSubscriptionsRequest, opts ...grpc.CallOption) (*notification.ListAdminSubscriptionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListAdminSubscriptions", varargs...)
	ret0, _ := ret[0].(*notification.ListAdminSubscriptionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAdminSubscriptions indicates an expected call of ListAdminSubscriptions.
func (mr *MockClientMockRecorder) ListAdminSubscriptions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAdminSubscriptions", reflect.TypeOf((*MockClient)(nil).ListAdminSubscriptions), varargs...)
}

// ListEnabledAdminSubscriptions mocks base method.
func (m *MockClient) ListEnabledAdminSubscriptions(ctx context.Context, in *notification.ListEnabledAdminSubscriptionsRequest, opts ...grpc.CallOption) (*notification.ListEnabledAdminSubscriptionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListEnabledAdminSubscriptions", varargs...)
	ret0, _ := ret[0].(*notification.ListEnabledAdminSubscriptionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEnabledAdminSubscriptions indicates an expected call of ListEnabledAdminSubscriptions.
func (mr *MockClientMockRecorder) ListEnabledAdminSubscriptions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEnabledAdminSubscriptions", reflect.TypeOf((*MockClient)(nil).ListEnabledAdminSubscriptions), varargs...)
}

// ListEnabledSubscriptions mocks base method.
func (m *MockClient) ListEnabledSubscriptions(ctx context.Context, in *notification.ListEnabledSubscriptionsRequest, opts ...grpc.CallOption) (*notification.ListEnabledSubscriptionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListEnabledSubscriptions", varargs...)
	ret0, _ := ret[0].(*notification.ListEnabledSubscriptionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEnabledSubscriptions indicates an expected call of ListEnabledSubscriptions.
func (mr *MockClientMockRecorder) ListEnabledSubscriptions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEnabledSubscriptions", reflect.TypeOf((*MockClient)(nil).ListEnabledSubscriptions), varargs...)
}

// ListSubscriptions mocks base method.
func (m *MockClient) ListSubscriptions(ctx context.Context, in *notification.ListSubscriptionsRequest, opts ...grpc.CallOption) (*notification.ListSubscriptionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListSubscriptions", varargs...)
	ret0, _ := ret[0].(*notification.ListSubscriptionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSubscriptions indicates an expected call of ListSubscriptions.
func (mr *MockClientMockRecorder) ListSubscriptions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSubscriptions", reflect.TypeOf((*MockClient)(nil).ListSubscriptions), varargs...)
}

// UpdateAdminSubscription mocks base method.
func (m *MockClient) UpdateAdminSubscription(ctx context.Context, in *notification.UpdateAdminSubscriptionRequest, opts ...grpc.CallOption) (*notification.UpdateAdminSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateAdminSubscription", varargs...)
	ret0, _ := ret[0].(*notification.UpdateAdminSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAdminSubscription indicates an expected call of UpdateAdminSubscription.
func (mr *MockClientMockRecorder) UpdateAdminSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdminSubscription", reflect.TypeOf((*MockClient)(nil).UpdateAdminSubscription), varargs...)
}

// UpdateSubscription mocks base method.
func (m *MockClient) UpdateSubscription(ctx context.Context, in *notification.UpdateSubscriptionRequest, opts ...grpc.CallOption) (*notification.UpdateSubscriptionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateSubscription", varargs...)
	ret0, _ := ret[0].(*notification.UpdateSubscriptionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSubscription indicates an expected call of UpdateSubscription.
func (mr *MockClientMockRecorder) UpdateSubscription(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSubscription", reflect.TypeOf((*MockClient)(nil).UpdateSubscription), varargs...)
}
