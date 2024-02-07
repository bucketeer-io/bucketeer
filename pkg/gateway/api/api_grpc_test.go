// Copyright 2024 The Bucketeer Authors.
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
	"errors"
	"fmt"
	"testing"
	"time"

	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestWithAPIKeyMemoryCacheTTL(t *testing.T) {
	t.Parallel()
	dur := time.Second
	f := WithAPIKeyMemoryCacheTTL(dur)
	opt := &options{}
	f(opt)
	assert.Equal(t, dur, opt.apiKeyMemoryCacheTTL)
}

func TestWithAPIKeyMemoryCacheEvictionInterval(t *testing.T) {
	t.Parallel()
	dur := time.Second
	f := WithAPIKeyMemoryCacheEvictionInterval(dur)
	opt := &options{}
	f(opt)
	assert.Equal(t, dur, opt.apiKeyMemoryCacheEvictionInterval)
}

func TestWithMetrics(t *testing.T) {
	t.Parallel()
	metrics := metrics.NewMetrics(
		9999,
		"/metrics",
	)
	reg := metrics.DefaultRegisterer()
	f := WithMetrics(reg)
	opt := &options{}
	f(opt)
	assert.Equal(t, reg, opt.metrics)
}

func TestWithLogger(t *testing.T) {
	t.Parallel()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	f := WithLogger(logger)
	opt := &options{}
	f(opt)
	assert.Equal(t, logger, opt.logger)
}

func TestNewGrpcGatewayService(t *testing.T) {
	t.Parallel()
	g := NewGrpcGatewayService(nil, nil, nil, nil, nil, nil, nil)
	assert.IsType(t, &grpcGatewayService{}, g)
}

func TestGrpcExtractAPIKeyID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testcases := []struct {
		ctx    context.Context
		key    string
		failed bool
	}{
		{
			ctx:    metadata.NewIncomingContext(context.TODO(), metadata.MD{}),
			key:    "",
			failed: true,
		},
		{
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{},
			}),
			key:    "",
			failed: true,
		},
		{
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{""},
			}),
			key:    "",
			failed: true,
		},
		{
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			}),
			key:    "test-key",
			failed: false,
		},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("index %d", i)
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		key, err := gs.extractAPIKeyID(tc.ctx)
		assert.Equal(t, tc.key, key, des)
		assert.Equal(t, tc.failed, err != nil, des)
	}
}

func TestGrpcGetEnvironmentAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		ctx         context.Context
		expected    *accountproto.EnvironmentAPIKey
		expectedErr error
	}{
		{
			desc: "exists in cache",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey:      &accountproto.APIKey{Id: "id-0"},
					}, nil)
			},
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			}),
			expected: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey:      &accountproto.APIKey{Id: "id-0"},
			},
			expectedErr: nil,
		},
		{
			desc: "ErrInvalidAPIKey",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.NotFound, "test"))
			},
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			}),
			expected:    nil,
			expectedErr: ErrInvalidAPIKey,
		},
		{
			desc: "ErrInternal",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Unknown, "test"))
			},
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			}),
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					&accountproto.GetAPIKeyBySearchingAllEnvironmentsResponse{EnvironmentApiKey: &accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey:      &accountproto.APIKey{Id: "id-0"},
					}}, nil)
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Put(gomock.Any()).Return(nil)
			},
			ctx: metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			}),
			expected: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey:      &accountproto.APIKey{Id: "id-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		id, err := gs.extractAPIKeyID(p.ctx)
		assert.NoError(t, err)
		actual, err := gs.getEnvironmentAPIKey(p.ctx, id)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcGetEnvironmentAPIKeyFromCache(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*cachev3mock.MockEnvironmentAPIKeyCache)
		expected    *accountproto.EnvironmentAPIKey
		expectedErr error
	}{
		{
			desc: "no error",
			setup: func(mtf *cachev3mock.MockEnvironmentAPIKeyCache) {
				mtf.EXPECT().Get(gomock.Any()).Return(&accountproto.EnvironmentAPIKey{}, nil)
			},
			expected:    &accountproto.EnvironmentAPIKey{},
			expectedErr: nil,
		},
		{
			desc: "error",
			setup: func(mtf *cachev3mock.MockEnvironmentAPIKeyCache) {
				mtf.EXPECT().Get(gomock.Any()).Return(nil, cache.ErrNotFound)
			},
			expected:    nil,
			expectedErr: cache.ErrNotFound,
		},
	}
	for _, p := range patterns {
		mock := cachev3mock.NewMockEnvironmentAPIKeyCache(mockController)
		p.setup(mock)
		actual, err := getEnvironmentAPIKeyFromCache(context.Background(), "id", mock, "caller", "layer")
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcCheckEnvironmentAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		inputEnvAPIKey *accountproto.EnvironmentAPIKey
		inputRole      accountproto.APIKey_Role
		expected       error
	}{
		{
			desc: "ErrBadRole",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SERVICE,
					Disabled: false,
				},
			},
			inputRole: accountproto.APIKey_SDK,
			expected:  ErrBadRole,
		},
		{
			desc: "ErrDisabledAPIKey: environment disabled",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SDK,
					Disabled: false,
				},
				EnvironmentDisabled: true,
			},
			inputRole: accountproto.APIKey_SDK,
			expected:  ErrDisabledAPIKey,
		},
		{
			desc: "ErrDisabledAPIKey: api key disabled",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SDK,
					Disabled: true,
				},
				EnvironmentDisabled: false,
			},
			inputRole: accountproto.APIKey_SDK,
			expected:  ErrDisabledAPIKey,
		},
		{
			desc: "no error",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SDK,
					Disabled: false,
				},
			},
			inputRole: accountproto.APIKey_SDK,
			expected:  nil,
		},
	}
	for _, p := range patterns {
		actual := checkEnvironmentAPIKey(p.inputEnvAPIKey, p.inputRole)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func TestGrpcValidateTrackRequest(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc     string
		input    *gwproto.TrackRequest
		expected error
	}{
		{
			desc:     "error: missing api key",
			input:    &gwproto.TrackRequest{},
			expected: ErrMissingAPIKey,
		},
		{
			desc:     "error: user ID is requried",
			input:    &gwproto.TrackRequest{Apikey: "api-key"},
			expected: ErrUserIDRequired,
		},
		{
			desc:     "error: goal ID is required",
			input:    &gwproto.TrackRequest{Apikey: "api-key", Userid: "user-id"},
			expected: ErrGoalIDRequired,
		},
		{
			desc: "error: tag is required",
			input: &gwproto.TrackRequest{
				Apikey: "api-key",
				Userid: "user-id",
				Goalid: "goal-id",
			},
			expected: ErrTagRequired,
		},
		{
			desc: "error: invalid timestamp",
			input: &gwproto.TrackRequest{
				Apikey: "api-key",
				Userid: "user-id",
				Goalid: "goal-id",
				Tag:    "tag",
			},
			expected: ErrInvalidTimestamp,
		},
		{
			desc: "success",
			input: &gwproto.TrackRequest{
				Apikey:    "api-key",
				Userid:    "user-id",
				Goalid:    "goal-id",
				Tag:       "tag",
				Timestamp: time.Now().Unix(),
			},
			expected: nil,
		},
	}
	gs := newGrpcGatewayServiceWithMock(t, mockController)
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := gs.validateTrackRequest(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGrpcValidateGetEvaluationsRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    *gwproto.GetEvaluationsRequest
		expected error
	}{
		{
			desc:     "user is empty",
			input:    &gwproto.GetEvaluationsRequest{Tag: "test"},
			expected: ErrUserRequired,
		},
		{
			desc:     "user ID is empty",
			input:    &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{}},
			expected: ErrUserIDRequired,
		},
		{
			desc:  "pass",
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "id"}},
		},
	}
	gs := grpcGatewayService{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := gs.validateGetEvaluationsRequest(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGrpcValidateGetEvaluationRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    *gwproto.GetEvaluationRequest
		expected error
	}{
		{
			desc:     "tag is empty",
			input:    &gwproto.GetEvaluationRequest{},
			expected: ErrTagRequired,
		},
		{
			desc:     "user is empty",
			input:    &gwproto.GetEvaluationRequest{Tag: "test"},
			expected: ErrUserRequired,
		},
		{
			desc:     "user ID is empty",
			input:    &gwproto.GetEvaluationRequest{Tag: "test", User: &userproto.User{}},
			expected: ErrUserIDRequired,
		},
		{
			desc:     "feature ID is empty",
			input:    &gwproto.GetEvaluationRequest{Tag: "test", User: &userproto.User{Id: "id"}},
			expected: ErrFeatureIDRequired,
		},
		{
			desc:  "pass",
			input: &gwproto.GetEvaluationRequest{Tag: "test", User: &userproto.User{Id: "id"}, FeatureId: "id"},
		},
	}
	gs := grpcGatewayService{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := gs.validateGetEvaluationRequest(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGrpcGetFeaturesFromCache(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		setup         func(*cachev3mock.MockFeaturesCache)
		environmentId string
		expected      *featureproto.Features
		expectedErr   error
	}{
		{
			desc: "no error",
			setup: func(mtf *cachev3mock.MockFeaturesCache) {
				mtf.EXPECT().Get(gomock.Any()).Return(&featureproto.Features{}, nil)
			},
			environmentId: "ns0",
			expected:      &featureproto.Features{},
			expectedErr:   nil,
		},
		{
			desc: "error",
			setup: func(mtf *cachev3mock.MockFeaturesCache) {
				mtf.EXPECT().Get(gomock.Any()).Return(nil, cache.ErrNotFound)
			},
			environmentId: "ns0",
			expected:      nil,
			expectedErr:   cache.ErrNotFound,
		},
	}
	for _, p := range patterns {
		mtfc := cachev3mock.NewMockFeaturesCache(mockController)
		p.setup(mtfc)
		gs := grpcGatewayService{featuresCache: mtfc}
		actual, err := gs.getFeaturesFromCache(context.Background(), p.environmentId)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcGetFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	now := time.Now()
	twentyNineDaysAgo := now.Add(-29 * 24 * time.Hour)
	thirtyOneDaysAgo := now.Add(-31 * 24 * time.Hour)

	patterns := []struct {
		desc          string
		setup         func(*grpcGatewayService)
		environmentId string
		expected      []*featureproto.Feature
		expectedErr   error
	}{
		{
			desc: "exists in redis",
			setup: func(gs *grpcGatewayService) {
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{{}},
					}, nil)
			},
			environmentId: "ns0",
			expectedErr:   nil,
			expected:      []*featureproto.Feature{{}},
		},
		{
			desc: "listFeatures fails",
			setup: func(gs *grpcGatewayService) {
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("test"))
			},
			environmentId: "ns0",
			expected:      nil,
			expectedErr:   ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					&featureproto.ListFeaturesResponse{Features: []*featureproto.Feature{
						{
							Id:      "id-0",
							Enabled: true,
						},
					}}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
			},
			environmentId: "ns0",
			expected: []*featureproto.Feature{
				{
					Id:      "id-0",
					Enabled: true,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: including off-variation features",
			setup: func(gs *grpcGatewayService) {
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					&featureproto.ListFeaturesResponse{Features: []*featureproto.Feature{
						{
							Id:      "id-0",
							Enabled: true,
						},
						{
							Id:           "id-1",
							Enabled:      true,
							OffVariation: "",
						},
						{
							Id:           "id-2",
							Enabled:      false,
							OffVariation: "var-2",
						},
						{
							Id:           "id-3",
							Enabled:      false,
							OffVariation: "",
						},
					}}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
			},
			environmentId: "ns0",
			expected: []*featureproto.Feature{
				{
					Id:      "id-0",
					Enabled: true,
				},
				{
					Id:           "id-1",
					Enabled:      true,
					OffVariation: "",
				},
				{
					Id:           "id-2",
					Enabled:      false,
					OffVariation: "var-2",
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: including archived features",
			setup: func(gs *grpcGatewayService) {
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					&featureproto.ListFeaturesResponse{Features: []*featureproto.Feature{
						{
							Id:       "id-0",
							Enabled:  true,
							Archived: false,
						},
						{
							Id:        "id-1",
							Enabled:   true,
							Archived:  true,
							UpdatedAt: twentyNineDaysAgo.Unix(),
						},
						{
							Id:        "id-2",
							Enabled:   true,
							Archived:  true,
							UpdatedAt: thirtyOneDaysAgo.Unix(),
						},
					}}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
			},
			environmentId: "ns0",
			expected: []*featureproto.Feature{
				{
					Id:       "id-0",
					Enabled:  true,
					Archived: false,
				},
				{
					Id:        "id-1",
					Enabled:   true,
					Archived:  true,
					UpdatedAt: twentyNineDaysAgo.Unix(),
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			actual, err := gs.getFeatures(context.Background(), p.environmentId)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcTrack(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.TrackRequest
		expected    *gwproto.TrackResponse
		expectedErr error
	}{
		{
			desc: "error: invalid api key",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.NotFound, "error: apy key not found"))
			},
			input: &gwproto.TrackRequest{
				Apikey:    "api-key",
				Userid:    "user-id",
				Goalid:    "goal-id",
				Tag:       "tag",
				Timestamp: time.Now().Unix(),
			},
			expected:    nil,
			expectedErr: ErrInvalidAPIKey,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.goalPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.TrackRequest{
				Apikey:    "api-key",
				Userid:    "user-id",
				Goalid:    "goal-id",
				Tag:       "tag",
				Timestamp: time.Now().Unix(),
			},
			expected:    &gwproto.TrackResponse{},
			expectedErr: nil,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		actual, err := gs.Track(ctx, p.input)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcGetEvaluationsContextCanceled(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		cancel      bool
		expected    *gwproto.GetEvaluationsResponse
		expectedErr error
	}{
		{
			desc:        "error: context canceled",
			cancel:      true,
			expected:    nil,
			expectedErr: ErrContextCanceled,
		},
		{
			desc:        "error: missing API key",
			cancel:      false,
			expected:    nil,
			expectedErr: ErrMissingAPIKey,
		},
	}
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		ctx, cancel := context.WithCancel(context.Background())
		if p.cancel {
			cancel()
		} else {
			defer cancel()
		}
		actual, err := gs.GetEvaluations(ctx, &gwproto.GetEvaluationsRequest{})
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcGetEvaluationsValidation(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetEvaluationsRequest
		expected    *gwproto.GetEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "missing tag",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-1",
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input:       &gwproto.GetEvaluationsRequest{User: &userproto.User{Id: "id-0"}},
			expected:    nil,
			expectedErr: ErrTagRequired,
		},
		{
			desc: "missing user id",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
			},
			input:       &gwproto.GetEvaluationsRequest{Tag: "test"},
			expected:    nil,
			expectedErr: ErrUserRequired,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "id-0"}},
			expected: &gwproto.GetEvaluationsResponse{
				State:             featureproto.UserEvaluations_FULL,
				Evaluations:       emptyUserEvaluations(t),
				UserEvaluationsId: "no_evaluations",
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
			actual, err := gs.GetEvaluations(ctx, p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetEvaluationsZeroFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetEvaluationsRequest
		expected    *gwproto.GetEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "zero feature",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "id-0"}},
			expected: &gwproto.GetEvaluationsResponse{
				State:             featureproto.UserEvaluations_FULL,
				Evaluations:       emptyUserEvaluations(t),
				UserEvaluationsId: "no_evaluations",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
			"authorization": []string{"test-key"},
		})
		actual, err := gs.GetEvaluations(ctx, p.input)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expected.State, actual.State, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcGetEvaluationsUserEvaluationsID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	vID1 := newUUID(t)
	vID2 := newUUID(t)
	vID3 := newUUID(t)
	vID4 := newUUID(t)
	vID5 := newUUID(t)
	vID6 := newUUID(t)

	features := []*featureproto.Feature{
		{
			Id: newUUID(t),
			Variations: []*featureproto.Variation{
				{
					Id:    vID1,
					Name:  "variation name true",
					Value: "true",
				},
				{
					Id:    newUUID(t),
					Name:  "variation name false",
					Value: "false",
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: vID1,
				},
			},
			Tags: []string{"android"},
		},
		{
			Id: newUUID(t),
			Variations: []*featureproto.Variation{
				{
					Id:    newUUID(t),
					Name:  "variation name true",
					Value: "true",
				},
				{
					Id:    vID2,
					Name:  "variation name false",
					Value: "false",
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: vID2,
				},
			},
			Tags: []string{"android"},
		},
	}

	features2 := []*featureproto.Feature{
		{
			Id: newUUID(t),
			Variations: []*featureproto.Variation{
				{
					Id:    vID3,
					Name:  "variation name true",
					Value: "true",
				},
				{
					Id:    newUUID(t),
					Name:  "variation name false",
					Value: "false",
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: vID3,
				},
			},
			Tags: []string{"ios"},
		},
		{
			Id: newUUID(t),
			Variations: []*featureproto.Variation{
				{
					Id:    newUUID(t),
					Name:  "variation name true",
					Value: "true",
				},
				{
					Id:    vID4,
					Name:  "variation name false",
					Value: "false",
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: vID4,
				},
			},
			Tags: []string{"ios"},
		},
	}

	features3 := []*featureproto.Feature{
		{
			Id: newUUID(t),
			Variations: []*featureproto.Variation{
				{
					Id:    vID5,
					Value: "true",
				},
				{
					Id:    newUUID(t),
					Value: "false",
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: vID5,
				},
			},
			Tags: []string{"web"},
		},
		{
			Id: newUUID(t),
			Variations: []*featureproto.Variation{
				{
					Id:    newUUID(t),
					Value: "true",
				},
				{
					Id:    vID6,
					Value: "false",
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: vID6,
				},
			},
			Tags: []string{"web"},
		},
	}
	multiFeatures := append(features, features2...)
	multiFeatures = append(multiFeatures, features3...)
	androidFeatures := features
	userID := "user-id-0"
	userMetadata := map[string]string{"b": "value-b", "c": "value-c", "a": "value-a", "d": "value-d"}
	ueidFromAndroidFeatures := featuredomain.UserEvaluationsID(userID, nil, androidFeatures)
	ueidWithDataFromAndroidFeatures := featuredomain.UserEvaluationsID(userID, userMetadata, androidFeatures)

	patterns := []struct {
		desc                      string
		setup                     func(*grpcGatewayService)
		input                     *gwproto.GetEvaluationsRequest
		expected                  *gwproto.GetEvaluationsResponse
		expectedErr               error
		expectedEvaluationsAssert func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			desc: "user evaluations id not set",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag: "android",
				User: &userproto.User{
					Id:   userID,
					Data: userMetadata,
				},
			},
			expected: &gwproto.GetEvaluationsResponse{
				State:             featureproto.UserEvaluations_FULL,
				UserEvaluationsId: ueidWithDataFromAndroidFeatures,
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user evaluations id is the same",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag: "android",
				User: &userproto.User{
					Id:   userID,
					Data: userMetadata,
				},
				UserEvaluationsId: ueidWithDataFromAndroidFeatures,
			},
			expected: &gwproto.GetEvaluationsResponse{
				State:             featureproto.UserEvaluations_FULL,
				UserEvaluationsId: ueidWithDataFromAndroidFeatures,
				Evaluations:       &featureproto.UserEvaluations{},
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user evaluations id is different",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag: "android",
				User: &userproto.User{
					Id:   userID,
					Data: userMetadata,
				},
				UserEvaluationsId: "evaluation-id",
			},
			expected: &gwproto.GetEvaluationsResponse{
				State:             featureproto.UserEvaluations_FULL,
				UserEvaluationsId: ueidWithDataFromAndroidFeatures,
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user_with_no_metadata_and_the_id_is_same",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag:               "android",
				User:              &userproto.User{Id: userID},
				UserEvaluationsId: ueidFromAndroidFeatures,
			},
			expected: &gwproto.GetEvaluationsResponse{
				State:             featureproto.UserEvaluations_FULL,
				UserEvaluationsId: ueidFromAndroidFeatures,
				Evaluations:       &featureproto.UserEvaluations{},
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user_with_no_metadata_and_the_id_is_different",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag:               "android",
				User:              &userproto.User{Id: userID},
				UserEvaluationsId: "evaluation-id",
			},
			expected: &gwproto.GetEvaluationsResponse{
				State:             featureproto.UserEvaluations_FULL,
				UserEvaluationsId: ueidFromAndroidFeatures,
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetEvaluations(ctx, p.input)
			assert.Equal(t, p.expected.State, actual.State, "%s", p.desc)
			assert.Equal(t, p.expected.UserEvaluationsId, actual.UserEvaluationsId, "%s", p.desc)
			p.expectedEvaluationsAssert(t, actual.Evaluations, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetEvaluationsNoSegmentList(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	vID1 := newUUID(t)
	vID2 := newUUID(t)
	vID3 := newUUID(t)
	vID4 := newUUID(t)

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetEvaluationsRequest
		expected    *gwproto.GetEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "state: full",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-a",
								Variations: []*featureproto.Variation{
									{
										Id:    vID1,
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    newUUID(t),
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID1,
									},
								},
								Tags: []string{"android"},
							},
							{
								Id: "feature-b",
								Variations: []*featureproto.Variation{
									{
										Id:    newUUID(t),
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    vID2,
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID2,
									},
								},
								Tags: []string{"android"},
							},
							{
								Id: "feature-c",
								Variations: []*featureproto.Variation{
									{
										Id:    vID3,
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    newUUID(t),
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID3,
									},
								},
								Tags: []string{"ios"},
							},
							{
								Id: "feature-d",
								Variations: []*featureproto.Variation{
									{
										Id:    newUUID(t),
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    vID4,
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: vID4,
									},
								},
								Tags: []string{"ios"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "ios", User: &userproto.User{Id: "id-0"}},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							VariationId: vID3,
						},
						{
							VariationId: vID4,
						},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
			"authorization": []string{"test-key"},
		})
		actual, err := gs.GetEvaluations(ctx, p.input)
		ev := p.expected.Evaluations.Evaluations
		av := actual.Evaluations.Evaluations
		assert.Equal(t, len(ev), len(av), "%s", p.desc)
		assert.Equal(t, p.expected.State, actual.State, "%s", p.desc)
		assert.Equal(t, ev[0].VariationId, av[0].VariationId, "%s", p.desc)
		assert.Equal(t, ev[1].VariationId, av[1].VariationId, "%s", p.desc)
		assert.NotEmpty(t, actual.UserEvaluationsId, "%s", p.desc)
		assert.ElementsMatch(t, p.expected.Evaluations.ArchivedFeatureIds, actual.Evaluations.ArchivedFeatureIds, p.desc)
		assert.Equal(t, p.expected.Evaluations.ForceUpdate, actual.Evaluations.ForceUpdate, p.desc)
		require.NoError(t, err)
	}
}

func TestGrpcGetEvaluationsEvaluateFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetEvaluationsRequest
		expected    *gwproto.GetEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "errInternal",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: "variation-b",
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("random error"))
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegmentUsers(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			input:       &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "id-0"}},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "state: full, evaluate features list segment from cache",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{

						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: "variation-b",
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-a",
									},
								},
								Tags: []string{"test"},
							},
						},
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&featureproto.SegmentUsers{
						SegmentId: "segment-id",
						Users: []*featureproto.SegmentUser{
							{
								SegmentId: "segment-id",
								UserId:    "user-id-1",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
							{
								SegmentId: "segment-id",
								UserId:    "user-id-2",
								State:     featureproto.SegmentUser_INCLUDED,
								Deleted:   false,
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "user-id-1"}},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "state: full, evaluate features list segment from storage",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: "variation-b",
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
						},
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("random error"))
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegmentUsers(gomock.Any(), gomock.Any()).Return(
					&featureproto.ListSegmentUsersResponse{}, nil)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "user-id-1"}},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "state: full, evaluate features",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "user-id-1"}},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: the cache includes archived features but the evaluation doesn't target them",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
							{
								Id:      "feature-id-2",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-d",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-d",
									},
								},
								Archived: true,
								Tags:     []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "user-id-1"}},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
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
			actual, err := gs.GetEvaluations(ctx, p.input)
			if err != nil {
				assert.Equal(t, p.expected, actual, p.desc)
				assert.Equal(t, p.expectedErr, err, p.desc)
			} else {
				assert.Equal(t, len(actual.Evaluations.Evaluations), 1, p.desc)
				assert.Equal(t, p.expected.State, actual.State, p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].Id, featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"), p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].UserId, "user-id-1", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].FeatureId, "feature-id-1", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].FeatureVersion, int32(2), p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].VariationId, "variation-b", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].VariationName, "variation name false", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].VariationValue, "false", p.desc)
				assert.Empty(t, p.expected.Evaluations.Evaluations[0].Variation.Description, p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].Reason, actual.Evaluations.Evaluations[0].Reason, p.desc)
				assert.ElementsMatch(t, p.expected.Evaluations.ArchivedFeatureIds, actual.Evaluations.ArchivedFeatureIds, p.desc)
				assert.Equal(t, p.expected.Evaluations.ForceUpdate, actual.Evaluations.ForceUpdate, p.desc)
				assert.NotEmpty(t, actual.UserEvaluationsId, p.desc)
				require.NoError(t, err)
			}
		})
	}
}

func TestGrpcGetEvaluationsByEvaluatedAt(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetEvaluationsRequest
		expected    *gwproto.GetEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "success: evaluate only flags that have benn updated since the last evaluation 10min ago",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags:      []string{"test"},
								UpdatedAt: now.Add(-5 * time.Minute).Unix(),
							},
							{
								Id:      "feature-id-2",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-d",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-d",
									},
								},
								Archived:  true,
								Tags:      []string{"test"},
								UpdatedAt: now.Add(-5 * time.Minute).Unix(),
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag:               "test",
				User:              &userproto.User{Id: "user-id-1"},
				UserEvaluationsId: "user-evaluations-id-1",
				UserEvaluationCondition: &gwproto.GetEvaluationsRequest_UserEvaluationCondition{
					EvaluatedAt:           now.Add(-10 * time.Minute).Unix(),
					UserAttributesUpdated: false,
				},
			},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
					ArchivedFeatureIds: []string{"feature-id-2"},
					ForceUpdate:        false,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: forceUpdate=true because UserEvaluationsId is empty",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags:      []string{"test"},
								UpdatedAt: now.Add(-5 * time.Minute).Unix(),
							},
							{
								Id:      "feature-id-2",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-d",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-d",
									},
								},
								Archived:  true,
								Tags:      []string{"test"},
								UpdatedAt: now.Add(-5 * time.Minute).Unix(),
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag:               "test",
				User:              &userproto.User{Id: "user-id-1"},
				UserEvaluationsId: "",
				UserEvaluationCondition: &gwproto.GetEvaluationsRequest_UserEvaluationCondition{
					EvaluatedAt:           now.Add(-10 * time.Minute).Unix(),
					UserAttributesUpdated: false,
				},
			},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
					ArchivedFeatureIds: []string{"feature-id-2"},
					ForceUpdate:        true,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: forceUpdate=true because UserEvaluationCondition.EvaluatedAt is 0",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
							{
								Id:      "feature-id-2",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-d",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-d",
									},
								},
								Archived: true,
								Tags:     []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag:               "test",
				User:              &userproto.User{Id: "user-id-1"},
				UserEvaluationsId: "user-evaluations-id-1",
				UserEvaluationCondition: &gwproto.GetEvaluationsRequest_UserEvaluationCondition{
					EvaluatedAt:           0,
					UserAttributesUpdated: false,
				},
			},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
					ArchivedFeatureIds: []string{},
					ForceUpdate:        true,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: forceUpdate=true because UserEvaluationCondition.EvaluatedAt is 30days ago",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
							{
								Id:      "feature-id-2",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-d",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-d",
									},
								},
								Archived: true,
								Tags:     []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag:               "test",
				User:              &userproto.User{Id: "user-id-1"},
				UserEvaluationsId: "user-evaluations-id-1",
				UserEvaluationCondition: &gwproto.GetEvaluationsRequest_UserEvaluationCondition{
					EvaluatedAt:           now.Add(-31 * 24 * time.Hour).Unix(),
					UserAttributesUpdated: false,
				},
			},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"),
							UserId:         "user-id-1",
							FeatureId:      "feature-id-1",
							FeatureVersion: int32(2),
							VariationId:    "variation-b",
							VariationName:  "variation name false",
							VariationValue: "false",
							Variation: &featureproto.Variation{
								Id:          "variation-b",
								Name:        "variation name false",
								Value:       "false",
								Description: "",
							},
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
					ArchivedFeatureIds: []string{},
					ForceUpdate:        true,
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
			actual, err := gs.GetEvaluations(ctx, p.input)
			if err != nil {
				assert.Equal(t, p.expected, actual, p.desc)
				assert.Equal(t, p.expectedErr, err, p.desc)
			} else {
				assert.Equal(t, len(actual.Evaluations.Evaluations), 1, p.desc)
				assert.Equal(t, p.expected.State, actual.State, p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].Id, featuredomain.EvaluationID("feature-id-1", int32(2), "user-id-1"), p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].UserId, "user-id-1", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].FeatureId, "feature-id-1", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].FeatureVersion, int32(2), p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].VariationId, "variation-b", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].VariationName, "variation name false", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].VariationValue, "false", p.desc)
				assert.Empty(t, p.expected.Evaluations.Evaluations[0].Variation.Description, p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].Reason, actual.Evaluations.Evaluations[0].Reason, p.desc)
				assert.ElementsMatch(t, p.expected.Evaluations.ArchivedFeatureIds, actual.Evaluations.ArchivedFeatureIds, p.desc)
				assert.Equal(t, p.expected.Evaluations.ForceUpdate, actual.Evaluations.ForceUpdate, p.desc)
				assert.NotEmpty(t, actual.UserEvaluationsId, p.desc)
				require.NoError(t, err)
			}
		})
	}
}

func TestGrpcGetEvaluation(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetEvaluationRequest
		expected    *featureproto.Evaluation
		expectedErr error
	}{
		{
			desc: "errFeatureNotFound",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-id-1",
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
							{
								Id: "feature-id-2",
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-d",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-d",
									},
								},
								Tags: []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input:       &gwproto.GetEvaluationRequest{Tag: "test", User: &userproto.User{Id: "id-0"}, FeatureId: "feature-id-3"},
			expectedErr: ErrFeatureNotFound,
		},
		{
			desc: "errInternal",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id: "feature-id-1",
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
							{
								Id: "feature-id-2",
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-d",
										Name:  "variation name false",
										Value: "false",
									},
								},
								Rules: []*featureproto.Rule{
									{
										Id: "rule-1",
										Strategy: &featureproto.Strategy{
											Type: featureproto.Strategy_FIXED,
											FixedStrategy: &featureproto.FixedStrategy{
												Variation: "variation-b",
											},
										},
										Clauses: []*featureproto.Clause{
											{
												Id:        "clause-1",
												Attribute: "name",
												Operator:  featureproto.Clause_SEGMENT,
												Values: []string{
													"id-0",
												},
											},
										},
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-d",
									},
								},
								Tags: []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("random error"))
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegmentUsers(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			input:       &gwproto.GetEvaluationRequest{Tag: "test", User: &userproto.User{Id: "id-0"}, FeatureId: "feature-id-2"},
			expectedErr: ErrInternal,
		},
		{
			desc: "return evaluation",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{
							{
								Id:      "feature-id-1",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
							{
								Id:      "feature-id-2",
								Version: int32(2),
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Name:  "variation name true",
										Value: "true",
									},
									{
										Id:    "variation-b",
										Name:  "variation name false",
										Value: "false",
									},
								},
								DefaultStrategy: &featureproto.Strategy{
									Type: featureproto.Strategy_FIXED,
									FixedStrategy: &featureproto.FixedStrategy{
										Variation: "variation-b",
									},
								},
								Tags: []string{"test"},
							},
						},
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.GetEvaluationRequest{Tag: "test", User: &userproto.User{Id: "user-id-2"}, FeatureId: "feature-id-2"},
			expected: &featureproto.Evaluation{
				Id:             featuredomain.EvaluationID("feature-id-2", int32(2), "user-id-2"),
				UserId:         "user-id-2",
				FeatureId:      "feature-id-2",
				FeatureVersion: int32(2),
				VariationId:    "variation-b",
				VariationName:  "variation name false",
				VariationValue: "false",
				Variation: &featureproto.Variation{
					Id:          "variation-b",
					Name:        "variation name false",
					Value:       "false",
					Description: "",
				},
				Reason: &featureproto.Reason{
					Type: featureproto.Reason_DEFAULT,
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
			actual, err := gs.GetEvaluation(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.Equal(t, p.expected.Id, actual.Evaluation.Id)
				assert.Equal(t, p.expected.UserId, actual.Evaluation.UserId)
				assert.Equal(t, p.expected.FeatureId, actual.Evaluation.FeatureId)
				assert.Equal(t, p.expected.FeatureVersion, actual.Evaluation.FeatureVersion)
				assert.Equal(t, p.expected.VariationId, actual.Evaluation.VariationId)
				assert.Equal(t, p.expected.VariationName, actual.Evaluation.VariationName)
				assert.Equal(t, p.expected.VariationValue, actual.Evaluation.VariationValue)
				assert.Empty(t, actual.Evaluation.Variation.Description)
			}
		})
	}
}

func TestGrpcRegisterEventsContextCanceled(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		cancel      bool
		expected    *gwproto.RegisterEventsResponse
		expectedErr error
	}{
		{
			desc:        "error: context canceled",
			cancel:      true,
			expected:    nil,
			expectedErr: ErrContextCanceled,
		},
		{
			desc:        "error: missing API key",
			cancel:      false,
			expected:    nil,
			expectedErr: ErrMissingAPIKey,
		},
	}
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		ctx, cancel := context.WithCancel(context.Background())
		if p.cancel {
			cancel()
		} else {
			defer cancel()
		}
		actual, err := gs.RegisterEvents(ctx, &gwproto.RegisterEventsRequest{})
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrcpRegisterEvents(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{Timestamp: time.Now().Unix()})
	if err != nil {
		t.Fatal("could not serialize goal event")
	}
	bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{Timestamp: time.Now().Unix()})
	if err != nil {
		t.Fatal("could not serialize evaluation event")
	}
	bInvalidEvent, err := proto.Marshal(&any.Any{})
	if err != nil {
		t.Fatal("could not serialize experiment event")
	}
	bMetricsEvent, err := proto.Marshal(&eventproto.MetricsEvent{Timestamp: time.Now().Unix()})
	if err != nil {
		t.Fatal("could not serialize metrics event")
	}
	uuid0 := newUUID(t)
	uuid1 := newUUID(t)
	uuid2 := newUUID(t)

	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.RegisterEventsRequest
		expected    *gwproto.RegisterEventsResponse
		expectedErr error
	}{
		{
			desc: "error: zero event",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
			},
			input:       &gwproto.RegisterEventsRequest{},
			expectedErr: ErrMissingEvents,
		},
		{
			desc: "error: ErrMissingEventID",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
			},
			input: &gwproto.RegisterEventsRequest{
				Events: []*eventproto.Event{
					{
						Id: "",
					},
				},
			},
			expectedErr: ErrMissingEventID,
		},
		{
			desc: "error: invalid message type",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.goalPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
				gs.evaluationPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.RegisterEventsRequest{
				Events: []*eventproto.Event{
					{
						Id: uuid0,
						Event: &any.Any{
							TypeUrl: "github.com/golang/protobuf/ptypes/any",
							Value:   bInvalidEvent,
						},
					},
				},
			},
			expected: &gwproto.RegisterEventsResponse{
				Errors: map[string]*gwproto.RegisterEventsResponse_Error{
					uuid0: {
						Retriable: false,
						Message:   "Invalid message type",
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK,
							Disabled: false,
						},
					}, nil)
				gs.goalPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
				gs.evaluationPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: &gwproto.RegisterEventsRequest{
				Events: []*eventproto.Event{
					{
						Id: uuid0,
						Event: &any.Any{
							TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
							Value:   bGoalEvent,
						},
					},
					{
						Id: uuid1,
						Event: &any.Any{
							TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
							Value:   bEvaluationEvent,
						},
					},
					{
						Id: uuid2,
						Event: &any.Any{
							TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.MetricsEvent",
							Value:   bMetricsEvent,
						},
					},
				},
			},
			expected:    &gwproto.RegisterEventsResponse{Errors: make(map[string]*gwproto.RegisterEventsResponse_Error)},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
			"authorization": []string{"test-key"},
		})
		actual, err := gs.RegisterEvents(ctx, p.input)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcContainsInvalidTimestampError(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc     string
		errs     map[string]*gwproto.RegisterEventsResponse_Error
		expected bool
	}{
		{
			desc: "error: invalid timestamp",
			errs: map[string]*gwproto.RegisterEventsResponse_Error{
				"id-test": {
					Retriable: false,
					Message:   errInvalidTimestamp.Error(),
				},
			},
			expected: true,
		},
		{
			desc: "error: different error",
			errs: map[string]*gwproto.RegisterEventsResponse_Error{
				"id-test": {
					Retriable: true,
					Message:   errUnmarshalFailed.Error(),
				},
			},
			expected: false,
		},
		{
			desc:     "error: empty",
			errs:     make(map[string]*gwproto.RegisterEventsResponse_Error),
			expected: false,
		},
	}
	for _, p := range patterns {
		gs := newGrpcGatewayServiceWithMock(t, mockController)
		actual := gs.containsInvalidTimestampError(p.errs)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func newGrpcGatewayServiceWithMock(t *testing.T, mockController *gomock.Controller) *grpcGatewayService {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &grpcGatewayService{
		featureClient:          featureclientmock.NewMockClient(mockController),
		accountClient:          accountclientmock.NewMockClient(mockController),
		goalPublisher:          publishermock.NewMockPublisher(mockController),
		userPublisher:          publishermock.NewMockPublisher(mockController),
		evaluationPublisher:    publishermock.NewMockPublisher(mockController),
		featuresCache:          cachev3mock.NewMockFeaturesCache(mockController),
		segmentUsersCache:      cachev3mock.NewMockSegmentUsersCache(mockController),
		environmentAPIKeyCache: cachev3mock.NewMockEnvironmentAPIKeyCache(mockController),
		opts:                   &defaultOptions,
		logger:                 logger,
	}
}

func TestGetTargetFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	multipleFs := []*featureproto.Feature{
		{
			Id: "fid3",
		},
		{
			Id: "fid2",
		},
		{
			Id: "fid10",
		},
		{
			Id: "fid",
		},
	}
	multiplePreFs := []*featureproto.Feature{
		{
			Id: "fid3",
			Prerequisites: []*featureproto.Prerequisite{
				{
					FeatureId: "fid10",
				},
			},
		},
		{
			Id: "fid2",
		},
		{
			Id: "fid10",
		},
		{
			Id: "fid",
			Prerequisites: []*featureproto.Prerequisite{
				{
					FeatureId: "fid3",
				},
			},
		},
	}
	patterns := []struct {
		desc        string
		fs          []*featureproto.Feature
		id          string
		expected    []*featureproto.Feature
		expectedErr error
	}{
		{
			desc:        "err: not found feature",
			id:          "not_found",
			fs:          multipleFs,
			expected:    nil,
			expectedErr: ErrFeatureNotFound,
		},
		{
			desc: "success: not configure prerequisite",
			id:   "fid",
			fs:   multipleFs,
			expected: []*featureproto.Feature{
				multipleFs[3],
			},
			expectedErr: nil,
		},
		{
			desc: "success: configure prerequisite",
			id:   "fid",
			fs:   multiplePreFs,
			expected: []*featureproto.Feature{
				multiplePreFs[0],
				multiplePreFs[2],
				multiplePreFs[3],
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			actual, err := gs.getTargetFeatures(p.fs, p.id)
			assert.ElementsMatch(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func emptyUserEvaluations(t *testing.T) *featureproto.UserEvaluations {
	t.Helper()
	return &featureproto.UserEvaluations{
		Id:                 "no_evaluations",
		Evaluations:        []*featureproto.Evaluation{},
		CreatedAt:          time.Now().Unix(),
		ArchivedFeatureIds: []string{},
		ForceUpdate:        false,
	}
}
