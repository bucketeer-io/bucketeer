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

	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	mockfeatureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func TestGrpcGatewayService_GetFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetFlagTriggerResponse
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
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: GetFlagTrigger error",
			ctx:  context.Background(),
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().
					GetFlagTrigger(gomock.Any(), gomock.Any()).
					Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			ctx:  context.Background(),
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().
					GetFlagTrigger(gomock.Any(), gomock.Any()).
					Return(&feature.GetFlagTriggerResponse{
						FlagTrigger: &feature.FlagTrigger{
							Id:            "sub-1",
							EnvironmentId: "env-1",
							FeatureId:     "feature-1",
						},
					}, nil)
			},
			expected: &gwproto.GetFlagTriggerResponse{
				FlagTrigger: &feature.FlagTrigger{
					Id:            "sub-1",
					EnvironmentId: "env-1",
					FeatureId:     "feature-1",
				},
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
			actual, err := gs.GetFlagTrigger(ctx, &gwproto.GetFlagTriggerRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_ListFlagTriggers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListFlagTriggersResponse
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
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: ListFlagTriggers error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().ListFlagTriggers(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
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
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().ListFlagTriggers(
					gomock.Any(), gomock.Any(),
				).Return(&feature.ListFlagTriggersResponse{
					FlagTriggers: []*feature.ListFlagTriggersResponse_FlagTriggerWithUrl{
						{
							FlagTrigger: &feature.FlagTrigger{
								Id:            "sub-1",
								EnvironmentId: "env-1",
								FeatureId:     "feature-1",
							},
							Url: "https://example.com/flag/sub-1",
						},
					},
					Cursor:     "0",
					TotalCount: 1,
				}, nil)
			},
			expected: &gwproto.ListFlagTriggersResponse{
				FlagTriggers: []*feature.ListFlagTriggersResponse_FlagTriggerWithUrl{
					{
						FlagTrigger: &feature.FlagTrigger{
							Id:            "sub-1",
							EnvironmentId: "env-1",
							FeatureId:     "feature-1",
						},
						Url: "https://example.com/flag/sub-1",
					},
				},
				Cursor:     "0",
				TotalCount: 1,
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
			actual, err := gs.ListFlagTriggers(ctx, &gwproto.ListFlagTriggersRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_CreateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.CreateFlagTriggerResponse
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
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: CreateFlagTriggers error",
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
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().CreateFlagTrigger(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			ctx:  context.Background(),
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().CreateFlagTrigger(
					gomock.Any(), gomock.Any(),
				).Return(&feature.CreateFlagTriggerResponse{
					FlagTrigger: &feature.FlagTrigger{
						Id:            "sub-1",
						EnvironmentId: "env-1",
						FeatureId:     "feature-1",
					},
				}, nil)
			},
			expected: &gwproto.CreateFlagTriggerResponse{
				FlagTrigger: &feature.FlagTrigger{
					Id:            "sub-1",
					EnvironmentId: "env-1",
					FeatureId:     "feature-1",
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
			actual, err := gs.CreateFlagTrigger(ctx, &gwproto.CreateFlagTriggerRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_DeleteFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.DeleteFlagTriggerResponse
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
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: DeleteFlagTrigger error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().DeleteFlagTrigger(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
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
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().DeleteFlagTrigger(
					gomock.Any(), gomock.Any(),
				).Return(&feature.DeleteFlagTriggerResponse{}, nil)
			},
			expected:    &gwproto.DeleteFlagTriggerResponse{},
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
			actual, err := gs.DeleteFlagTrigger(ctx, &gwproto.DeleteFlagTriggerRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_UpdateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.UpdateFlagTriggerResponse
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
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: UpdateFlagTrigger error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().UpdateFlagTrigger(
					gomock.Any(), gomock.Any(),
				).Return(nil, ErrInternal)
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
							Role:     accountproto.APIKey_PUBLIC_API_ADMIN,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*mockfeatureclient.MockClient).EXPECT().UpdateFlagTrigger(
					gomock.Any(), gomock.Any(),
				).Return(&feature.UpdateFlagTriggerResponse{}, nil)
			},
			expected:    &gwproto.UpdateFlagTriggerResponse{},
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
			actual, err := gs.UpdateFlagTrigger(ctx, &gwproto.UpdateFlagTriggerRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}
