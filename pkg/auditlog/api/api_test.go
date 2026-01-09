// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	v2asmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	v2alsmock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2/mock"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func TestNewAuditLogService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	accountClientMock := accountclientmock.NewMockClient(mockController)
	mysqlClient := mysqlmock.NewMockClient(mockController)
	logger := zap.NewNop()
	s := NewAuditLogService(accountClientMock, mysqlClient, WithLogger(logger))
	assert.IsType(t, &auditlogService{}, s)
}

func TestGetAuditLog(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		service     *auditlogService
		context     context.Context
		setup       func(*auditlogService)
		input       *proto.GetAuditLogRequest
		expected    *proto.GetAuditLogResponse
		expectedErr error
	}{
		{
			desc:    "errPermissionDenied",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithToken(t, false),
			setup:   func(s *auditlogService) {},
			input: &proto.GetAuditLogRequest{
				Id: "id-1",
			},
			expected:    nil,
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "err: missing ID",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, true),
			setup:   nil,
			input: &proto.GetAuditLogRequest{
				Id:            "",
				EnvironmentId: "env-1",
			},
			expected:    nil,
			expectedErr: statusMissingID.Err(),
		},
		{
			desc:    "err: ErrInternal",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, true),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().GetAuditLog(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal"))
			},
			input: &proto.GetAuditLogRequest{
				Id:            "id-1",
				EnvironmentId: "env-1",
			},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal")).Err(),
		},
		{
			desc:    "success",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, true),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().GetAuditLog(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&proto.AuditLog{
					Id: "id-1",
					Editor: &domaineventproto.Editor{
						Email: "test@bucketeer.io",
					},
				}, nil)
				s.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().GetAvatarAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{
					{
						Email:       "test@bucketeer.io",
						AvatarImage: []byte{0x1},
					},
				}, nil)
			},
			input: &proto.GetAuditLogRequest{
				Id:            "id-1",
				EnvironmentId: "env-1",
			},
			expected: &proto.GetAuditLogResponse{
				AuditLog: &proto.AuditLog{
					Id: "id-1",
					Editor: &domaineventproto.Editor{
						Email:       "test@bucketeer.io",
						AvatarImage: []byte{0x1},
					},
					LocalizedMessage: domainevent.LocalizedMessage(domaineventproto.Event_UNKNOWN, locale.NewLocalizer(context.Background())),
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := p.service
			if p.setup != nil {
				p.setup(s)
			}
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"en"},
			})

			actual, err := s.GetAuditLog(ctx, p.input)
			if err != nil {
				assert.Equal(t, p.expectedErr, err)
				return
			}
			assert.Equal(t, p.expected.AuditLog, actual.AuditLog)
		})
	}
}

func TestListAuditLogsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		service     *auditlogService
		context     context.Context
		setup       func(*auditlogService)
		input       *proto.ListAuditLogsRequest
		expected    *proto.ListAuditLogsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			service:     newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context:     createContextWithToken(t, true),
			setup:       nil,
			input:       &proto.ListAuditLogsRequest{Cursor: "XXX", EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc:    "err: ErrInternal",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, true),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal"))
			},
			input:       &proto.ListAuditLogsRequest{EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal")).Err(),
		},
		{
			desc:        "errPermissionDenied",
			service:     newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context:     createContextWithToken(t, false),
			setup:       func(s *auditlogService) {},
			input:       &proto.ListAuditLogsRequest{PageSize: 2, Cursor: "", EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "success",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, true),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
				s.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().GetAvatarAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, nil)
			},
			input:       &proto.ListAuditLogsRequest{PageSize: 2, Cursor: "", EnvironmentId: "ns0"},
			expected:    &proto.ListAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: 10},
			expectedErr: nil,
		},
		{
			desc:    "success with default page size when page_size is 0",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, true),
			setup: func(s *auditlogService) {
				// Expect the default page size (200) to be used
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(),
					&mysql.ListOptions{
						Limit:  200, // maxAuditLogPageSize (used as default when page_size is 0)
						Offset: 0,
						Filters: []*mysql.FilterV2{
							{
								Column:   "environment_id",
								Operator: mysql.OperatorEqual,
								Value:    "ns0",
							},
						},
						Orders: []*mysql.Order{
							{
								Column:    "timestamp",
								Direction: mysql.OrderDirectionDesc,
							},
						},
					},
				).Return(createAuditLogs(t), 200, int64(10), nil)
				s.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().GetAvatarAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, nil)
			},
			input:       &proto.ListAuditLogsRequest{PageSize: 0, Cursor: "", EnvironmentId: "ns0"},
			expected:    &proto.ListAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "200", TotalCount: 10},
			expectedErr: nil,
		},
		{
			desc:    "success with Viewer Account",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: createContextWithToken(t, false),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
				s.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().GetAvatarAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, nil)
			},
			input:       &proto.ListAuditLogsRequest{PageSize: 2, Cursor: "", EnvironmentId: "ns0"},
			expected:    &proto.ListAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: 10},
			expectedErr: nil,
		},
		{
			desc:    "success: page size exceeds maximum",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, true),
			setup: func(s *auditlogService) {
				// Expect the maximum page size (200) to be used even though 1000 was requested
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(),
					&mysql.ListOptions{
						Limit:  200, // maxAuditLogPageSize
						Offset: 0,
						Filters: []*mysql.FilterV2{
							{
								Column:   "environment_id",
								Operator: mysql.OperatorEqual,
								Value:    "ns0",
							},
						},
						Orders: []*mysql.Order{
							{
								Column:    "timestamp",
								Direction: mysql.OrderDirectionDesc,
							},
						},
					},
				).Return(createAuditLogs(t), 200, int64(10), nil)
				s.accountStorage.(*v2asmock.MockAccountStorage).EXPECT().GetAvatarAccountsV2(
					gomock.Any(), gomock.Any(),
				).Return([]*accountproto.AccountV2{}, nil)
			},
			input: &proto.ListAuditLogsRequest{
				PageSize:       1000, // Exceeds maximum, should be capped at 200
				EnvironmentId:  "ns0",
				OrderBy:        proto.ListAuditLogsRequest_TIMESTAMP,
				OrderDirection: proto.ListAuditLogsRequest_DESC,
			},
			expected: &proto.ListAuditLogsResponse{
				AuditLogs:  createAuditLogs(t),
				Cursor:     "200", // Capped at maximum page size
				TotalCount: 10,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := p.service
			if p.setup != nil {
				p.setup(s)
			}
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})

			actual, err := s.ListAuditLogs(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestListAdminAuditLogsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*auditlogService)
		input       *proto.ListAdminAuditLogsRequest
		expected    *proto.ListAdminAuditLogsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			setup:       nil,
			input:       &proto.ListAdminAuditLogsRequest{Cursor: "invalid"},
			expected:    nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *auditlogService) {
				s.adminAuditLogStorage.(*v2alsmock.MockAdminAuditLogStorage).EXPECT().ListAdminAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal"))
			},
			input:       &proto.ListAdminAuditLogsRequest{},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal")).Err(),
		},
		{
			desc: "success",
			setup: func(s *auditlogService) {
				s.adminAuditLogStorage.(*v2alsmock.MockAdminAuditLogStorage).EXPECT().ListAdminAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
			},
			input:       &proto.ListAdminAuditLogsRequest{PageSize: 2, Cursor: ""},
			expected:    &proto.ListAdminAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: 10},
			expectedErr: nil,
		},
		{
			desc: "success with default page size when page_size is 0",
			setup: func(s *auditlogService) {
				s.adminAuditLogStorage.(*v2alsmock.MockAdminAuditLogStorage).EXPECT().ListAdminAuditLogs(
					gomock.Any(),
					&mysql.ListOptions{
						Limit:   200, // maxAuditLogPageSize (used as default when page_size is 0)
						Offset:  0,
						Orders:  []*mysql.Order{{Column: "timestamp", Direction: mysql.OrderDirectionDesc}},
						Filters: []*mysql.FilterV2{},
					},
				).Return(createAuditLogs(t), 200, int64(10), nil)
			},
			input:       &proto.ListAdminAuditLogsRequest{PageSize: 0, Cursor: ""},
			expected:    &proto.ListAdminAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "200", TotalCount: 10},
			expectedErr: nil,
		},
		{
			desc: "success: page size exceeds maximum",
			setup: func(s *auditlogService) {
				s.adminAuditLogStorage.(*v2alsmock.MockAdminAuditLogStorage).EXPECT().ListAdminAuditLogs(
					gomock.Any(),
					&mysql.ListOptions{
						Limit:   200, // maxAuditLogPageSize
						Offset:  0,
						Orders:  []*mysql.Order{{Column: "timestamp", Direction: mysql.OrderDirectionDesc}},
						Filters: []*mysql.FilterV2{},
					},
				).Return(createAuditLogs(t), 200, int64(10), nil)
			},
			input:       &proto.ListAdminAuditLogsRequest{PageSize: 1000, Cursor: ""},
			expected:    &proto.ListAdminAuditLogsResponse{AuditLogs: createAuditLogs(t), Cursor: "200", TotalCount: 10},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListAdminAuditLogs(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestListFeatureHistoryMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		service     *auditlogService
		context     context.Context
		setup       func(*auditlogService)
		input       *proto.ListFeatureHistoryRequest
		expected    *proto.ListFeatureHistoryResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			service:     newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context:     createContextWithToken(t, false),
			setup:       nil,
			input:       &proto.ListFeatureHistoryRequest{Cursor: "XXX", EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc:    "err: ErrInternal",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, false),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal"))
			},
			input:       &proto.ListFeatureHistoryRequest{EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.AuditlogPackageName, "internal")).Err(),
		},
		{
			desc:    "errPermissionDenied",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: createContextWithTokenRoleUnassigned(t),
			setup:   func(s *auditlogService) {},
			input: &proto.ListFeatureHistoryRequest{
				FeatureId: "fid-1", PageSize: 2, Cursor: "", EnvironmentId: "ns0",
			},
			expected:    nil,
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "success",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, false),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
			},
			input: &proto.ListFeatureHistoryRequest{
				FeatureId: "fid-1", PageSize: 2, Cursor: "", EnvironmentId: "ns0",
			},
			expected:    &proto.ListFeatureHistoryResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: int64(10)},
			expectedErr: nil,
		},
		{
			desc:    "success with viewer account",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: createContextWithTokenRoleUnassigned(t),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(), gomock.Any(),
				).Return(createAuditLogs(t), 2, int64(10), nil)
			},
			input: &proto.ListFeatureHistoryRequest{
				FeatureId: "fid-1", PageSize: 2, Cursor: "", EnvironmentId: "ns0",
			},
			expected:    &proto.ListFeatureHistoryResponse{AuditLogs: createAuditLogs(t), Cursor: "2", TotalCount: int64(10)},
			expectedErr: nil,
		},
		{
			desc:    "success with default page size when page_size is 0",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, false),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(),
					&mysql.ListOptions{
						Limit:  200, // maxAuditLogPageSize (used as default when page_size is 0)
						Offset: 0,
						Filters: []*mysql.FilterV2{
							{
								Column:   "environment_id",
								Operator: mysql.OperatorEqual,
								Value:    "ns0",
							},
							{
								Column:   "entity_type",
								Operator: mysql.OperatorEqual,
								Value:    int32(domaineventproto.Event_FEATURE),
							},
							{
								Column:   "entity_id",
								Operator: mysql.OperatorEqual,
								Value:    "fid-1",
							},
						},
						Orders: []*mysql.Order{
							{
								Column:    "timestamp",
								Direction: mysql.OrderDirectionDesc,
							},
						},
					},
				).Return(createAuditLogs(t), 200, int64(10), nil)
			},
			input: &proto.ListFeatureHistoryRequest{
				FeatureId: "fid-1", PageSize: 0, Cursor: "", EnvironmentId: "ns0",
			},
			expected:    &proto.ListFeatureHistoryResponse{AuditLogs: createAuditLogs(t), Cursor: "200", TotalCount: int64(10)},
			expectedErr: nil,
		},
		{
			desc:    "success: page size exceeds maximum",
			service: newAuditLogServiceWithGetAccountByEnvironmentMock(t, mockController, accountproto.AccountV2_Role_Organization_OWNER, accountproto.AccountV2_Role_Environment_EDITOR),
			context: createContextWithToken(t, false),
			setup: func(s *auditlogService) {
				s.auditLogStorage.(*v2alsmock.MockAuditLogStorage).EXPECT().ListAuditLogs(
					gomock.Any(),
					&mysql.ListOptions{
						Limit:  200, // maxAuditLogPageSize
						Offset: 0,
						Filters: []*mysql.FilterV2{
							{
								Column:   "environment_id",
								Operator: mysql.OperatorEqual,
								Value:    "ns0",
							},
							{
								Column:   "entity_type",
								Operator: mysql.OperatorEqual,
								Value:    int32(domaineventproto.Event_FEATURE),
							},
							{
								Column:   "entity_id",
								Operator: mysql.OperatorEqual,
								Value:    "fid-1",
							},
						},
						Orders: []*mysql.Order{
							{
								Column:    "timestamp",
								Direction: mysql.OrderDirectionDesc,
							},
						},
					},
				).Return(createAuditLogs(t), 200, int64(10), nil)
			},
			input: &proto.ListFeatureHistoryRequest{
				FeatureId: "fid-1", PageSize: 1000, Cursor: "", EnvironmentId: "ns0",
			},
			expected:    &proto.ListFeatureHistoryResponse{AuditLogs: createAuditLogs(t), Cursor: "200", TotalCount: int64(10)},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := p.service
			if p.setup != nil {
				p.setup(s)
			}
			ctx := p.context
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{
				"accept-language": []string{"ja"},
			})

			actual, err := s.ListFeatureHistory(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newAuditLogServiceWithGetAccountByEnvironmentMock(t *testing.T, mockController *gomock.Controller, ro accountproto.AccountV2_Role_Organization, re accountproto.AccountV2_Role_Environment) *auditlogService {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: ro,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          re,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	return &auditlogService{
		accountClient:        accountClientMock,
		accountStorage:       v2asmock.NewMockAccountStorage(mockController),
		auditLogStorage:      v2alsmock.NewMockAuditLogStorage(mockController),
		adminAuditLogStorage: v2alsmock.NewMockAdminAuditLogStorage(mockController),
		logger:               logger.Named("api"),
	}
}

func createAuditLogs(t *testing.T) []*proto.AuditLog {
	t.Helper()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	msgUnknown := domainevent.LocalizedMessage(domaineventproto.Event_UNKNOWN, localizer)
	return []*proto.AuditLog{
		{Id: "id-0", LocalizedMessage: msgUnknown, Editor: &domaineventproto.Editor{}},
		{Id: "id-1", LocalizedMessage: msgUnknown, Editor: &domaineventproto.Editor{}},
	}
}

func createContextWithToken(t *testing.T, isSystemAdmin bool) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Email:         "test@example.com",
		IsSystemAdmin: isSystemAdmin,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}

func createContextWithTokenRoleUnassigned(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Email: "test@example.com",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}
