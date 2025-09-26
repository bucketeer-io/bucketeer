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
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func TestGrpcGatewayService_CreateSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.CreateSegmentResponse
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
			desc: "fails: internal error",
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().CreateSegment(
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
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().CreateSegment(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.CreateSegmentResponse{
					Segment: &featureproto.Segment{
						Id:   "segment-id",
						Name: "segment-name",
					},
				}, nil)
			},
			expected: &gwproto.CreateSegmentResponse{
				Segment: &featureproto.Segment{
					Id:   "segment-id",
					Name: "segment-name",
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
			actual, err := gs.CreateSegment(ctx, &gwproto.CreateSegmentRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_GetSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetSegmentResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: internal error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().GetSegment(
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
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.GetSegmentResponse{
					Segment: &featureproto.Segment{
						Id:   "segment-id",
						Name: "segment-name",
					},
				}, nil)
			},
			expected: &gwproto.GetSegmentResponse{
				Segment: &featureproto.Segment{
					Id:   "segment-id",
					Name: "segment-name",
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
			actual, err := gs.GetSegment(ctx, &gwproto.GetSegmentRequest{Id: "segment-id"})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_ListSegments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListSegmentsResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: true,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: internal error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegments(
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
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegments(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.ListSegmentsResponse{
					Segments: []*featureproto.Segment{
						{Id: "segment-1", Name: "Segment 1"},
						{Id: "segment-2", Name: "Segment 2"},
					},
				}, nil)
			},
			expected: &gwproto.ListSegmentsResponse{
				Segments: []*featureproto.Segment{
					{Id: "segment-1", Name: "Segment 1"},
					{Id: "segment-2", Name: "Segment 2"},
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
			actual, err := gs.ListSegments(ctx, &gwproto.ListSegmentsRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_UpdateSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.UpdateSegmentResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: internal error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().UpdateSegment(
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
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.UpdateSegmentResponse{}, nil)
			},
			expected:    &gwproto.UpdateSegmentResponse{},
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
			actual, err := gs.UpdateSegment(ctx, &gwproto.UpdateSegmentRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_DeleteSegment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.DeleteSegmentResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: internal error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().DeleteSegment(
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
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().DeleteSegment(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.DeleteSegmentResponse{}, nil)
			},
			expected:    &gwproto.DeleteSegmentResponse{},
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
			actual, err := gs.DeleteSegment(ctx, &gwproto.DeleteSegmentRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGatewayService_BulkUploadSegmentUsers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.BulkUploadSegmentUsersResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fails: internal error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().BulkUploadSegmentUsers(
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
						Environment: &environmentproto.EnvironmentV2{Id: "env-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "key-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().BulkUploadSegmentUsers(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.BulkUploadSegmentUsersResponse{}, nil)
			},
			expected:    &gwproto.BulkUploadSegmentUsersResponse{},
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
			actual, err := gs.BulkUploadSegmentUsers(ctx, &gwproto.BulkUploadSegmentUsersRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}
