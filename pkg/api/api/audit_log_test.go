// Copyright 2025 The Bucketeer Authors.
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
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	auditlogclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	auditlogproto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	"github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func TestGrpcGetAuditLog(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetAuditLogResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.auditLogClient.(*auditlogclientmock.MockClient).EXPECT().
					GetAuditLog(gomock.Any(), gomock.Any()).
					Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.auditLogClient.(*auditlogclientmock.MockClient).EXPECT().
					GetAuditLog(gomock.Any(), gomock.Any()).
					Return(&auditlogproto.GetAuditLogResponse{
						AuditLog: &auditlogproto.AuditLog{
							Id:         "audit-log-id",
							Timestamp:  1749548894,
							EntityType: domain.Event_FEATURE,
							Type:       domain.Event_FEATURE_CREATED,
							EntityId:   "entity-id",
							Editor: &domain.Editor{
								Email: "bucketeer@demo.io",
								Name:  "Bucketeer Demo",
							},
							EntityData:         "{\"key\":\"value\"}",
							PreviousEntityData: "{\"key\":\"previous-value\"}",
						},
					}, nil)
			},
			expected: &gwproto.GetAuditLogResponse{
				AuditLog: &auditlogproto.AuditLog{
					Id:         "audit-log-id",
					Timestamp:  1749548894,
					EntityType: domain.Event_FEATURE,
					Type:       domain.Event_FEATURE_CREATED,
					EntityId:   "entity-id",
					Editor: &domain.Editor{
						Email: "bucketeer@demo.io",
						Name:  "Bucketeer Demo",
					},
					EntityData:         "{\"key\":\"value\"}",
					PreviousEntityData: "{\"key\":\"previous-value\"}",
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetAuditLog(ctx, &gwproto.GetAuditLogRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcListAuditLogs(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListAuditLogsResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.auditLogClient.(*auditlogclientmock.MockClient).EXPECT().
					ListAuditLogs(gomock.Any(), gomock.Any()).
					Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.auditLogClient.(*auditlogclientmock.MockClient).EXPECT().
					ListAuditLogs(gomock.Any(), gomock.Any()).
					Return(&auditlogproto.ListAuditLogsResponse{
						AuditLogs: []*auditlogproto.AuditLog{
							{Id: "audit-log-1"},
							{Id: "audit-log-2"},
						},
						Cursor:     "cursor-1",
						TotalCount: 2,
					}, nil)
			},
			expected: &gwproto.ListAuditLogsResponse{
				AuditLogs: []*auditlogproto.AuditLog{
					{Id: "audit-log-1"},
					{Id: "audit-log-2"},
				},
				Cursor:     "cursor-1",
				TotalCount: 2,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.ListAuditLogs(ctx, &gwproto.ListAuditLogsRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcListFeatureHistory(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListFeatureHistoryResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.auditLogClient.(*auditlogclientmock.MockClient).EXPECT().
					ListFeatureHistory(gomock.Any(), gomock.Any()).
					Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.auditLogClient.(*auditlogclientmock.MockClient).EXPECT().
					ListFeatureHistory(gomock.Any(), gomock.Any()).
					Return(&auditlogproto.ListFeatureHistoryResponse{
						AuditLogs: []*auditlogproto.AuditLog{
							{Id: "feature-history-1"},
							{Id: "feature-history-2"},
						},
						Cursor:     "cursor-1",
						TotalCount: 2,
					}, nil)
			},
			expected: &gwproto.ListFeatureHistoryResponse{
				AuditLogs: []*auditlogproto.AuditLog{
					{Id: "feature-history-1"},
					{Id: "feature-history-2"},
				},
				Cursor:     "cursor-1",
				TotalCount: 2,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.ListFeatureHistory(ctx, &gwproto.ListFeatureHistoryRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}
