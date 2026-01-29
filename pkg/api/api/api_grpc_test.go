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
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	auditlogclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/client/mock"
	autoopsclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	coderefclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/coderef/client/mock"
	environmentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	eventcounterclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	notificationclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/notification/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	pushclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/push/client/mock"
	tagclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/client/mock"
	teamclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/team/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
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
	g := NewGrpcGatewayService(context.Background(), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
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
		key, err := gs.extractAPIKey(tc.ctx)
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
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetEnvironmentAPIKey(gomock.Any(), gomock.Any()).Return(
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
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetEnvironmentAPIKey(gomock.Any(), gomock.Any()).Return(
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
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetEnvironmentAPIKey(gomock.Any(), gomock.Any()).Return(
					&accountproto.GetEnvironmentAPIKeyResponse{EnvironmentApiKey: &accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey:      &accountproto.APIKey{Id: "id-0"},
					}}, nil)
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
		id, err := gs.extractAPIKey(p.ctx)
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
		inputRole      []accountproto.APIKey_Role
		expected       error
	}{
		{
			desc: "ErrBadRole",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SDK_SERVER,
					Disabled: false,
				},
			},
			inputRole: []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT},
			expected:  ErrBadRole,
		},
		{
			desc: "ErrDisabledAPIKey: environment disabled",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SDK_CLIENT,
					Disabled: false,
				},
				EnvironmentDisabled: true,
			},
			inputRole: []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT},
			expected:  ErrDisabledAPIKey,
		},
		{
			desc: "ErrDisabledAPIKey: api key disabled",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SDK_CLIENT,
					Disabled: true,
				},
				EnvironmentDisabled: false,
			},
			inputRole: []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT},
			expected:  ErrDisabledAPIKey,
		},
		{
			desc: "no error",
			inputEnvAPIKey: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_SDK_CLIENT,
					Disabled: false,
				},
			},
			inputRole: []accountproto.APIKey_Role{accountproto.APIKey_SDK_CLIENT},
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
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetEnvironmentAPIKey(gomock.Any(), gomock.Any()).Return(
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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

func TestGrpcGetSegmentUsersContextCanceled(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		cancel      bool
		expected    *gwproto.GetSegmentUsersResponse
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
		actual, err := gs.GetSegmentUsers(ctx, &gwproto.GetSegmentUsersRequest{})
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcGetSegmentUsers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	timeNow := time.Now()
	apiKey := "api-key-id"
	envID := "ns0"
	tag := "tag"
	multiSegmentUsers := []*featureproto.SegmentUsers{
		{
			SegmentId: "segment-id-2",
			Users: []*featureproto.SegmentUser{
				{
					SegmentId: "segment-id-2",
					UserId:    "user-id",
				},
			},
			UpdatedAt: timeNow.Add(-30 * time.Minute).Unix(),
		},
		{
			SegmentId: "segment-id-4",
			Users: []*featureproto.SegmentUser{
				{
					SegmentId: "segment-id-4",
					UserId:    "user-id",
				},
			},
			UpdatedAt: timeNow.Add(-1 * time.Hour).Unix(),
		},
		{
			SegmentId: "segment-id-5",
			Users: []*featureproto.SegmentUser{
				{
					SegmentId: "segment-id-5",
					UserId:    "user-id",
				},
			},
			UpdatedAt: timeNow.Add(-1 * time.Hour).Unix(),
		},
	}
	singleFeature := []*featureproto.Feature{
		{
			Id:        "feature-id-1",
			Version:   1,
			Tags:      []string{tag},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
	}
	multiFeatures := []*featureproto.Feature{
		{
			Id:        "feature-id-1",
			Version:   1,
			Tags:      []string{tag},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
		{
			Id:      "feature-id-2",
			Version: 1,
			Tags:    []string{},
			Rules: []*featureproto.Rule{
				{
					Id: "rule-id",
					Clauses: []*featureproto.Clause{
						{
							Id:       "clause-id",
							Operator: featureproto.Clause_SEGMENT,
							Values:   []string{"segment-id-2"},
						},
					},
				},
			},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
		{
			Id:        "feature-id-3",
			Version:   1,
			Tags:      []string{},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
		{
			Id:      "feature-id-4",
			Version: 1,
			Tags:    []string{tag},
			Rules: []*featureproto.Rule{
				{
					Id: "rule-id",
					Clauses: []*featureproto.Clause{
						{
							Id:       "clause-id",
							Operator: featureproto.Clause_SEGMENT,
							Values:   []string{"segment-id-4"},
						},
					},
				},
			},
			UpdatedAt: timeNow.Add(-secondsToReturnAllFlags * time.Second).Unix(),
			Archived:  true,
		},
		{
			Id:      "feature-id-5",
			Version: 1,
			Tags:    []string{tag},
			Rules: []*featureproto.Rule{
				{
					Id: "rule-id",
					Clauses: []*featureproto.Clause{
						{
							Id:       "clause-id",
							Operator: featureproto.Clause_SEGMENT,
							Values:   []string{"segment-id-5"},
						},
					},
				},
			},
			UpdatedAt: timeNow.Add(-10 * time.Minute).Unix(),
			Archived:  false,
		},
	}
	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetSegmentUsersRequest
		expected    *gwproto.GetSegmentUsersResponse
		expectedErr error
	}{
		{
			desc: "err: environment api key not found",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					nil, errors.New("internal error"))
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetEnvironmentAPIKey(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.NotFound, "test"))
			},
			input:       &gwproto.GetSegmentUsersRequest{},
			expected:    nil,
			expectedErr: ErrInvalidAPIKey,
		},
		{
			desc: "err: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			input:       &gwproto.GetSegmentUsersRequest{},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "err: source id is required",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SourceId: eventproto.SourceId_UNKNOWN,
			},
			expected:    nil,
			expectedErr: ErrSourceIDRequired,
		},
		{
			desc: "err: sdk version is required",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SourceId:   eventproto.SourceId_GO_SERVER,
				SdkVersion: "",
			},
			expected:    nil,
			expectedErr: ErrSDKVersionRequired,
		},
		{
			desc: "err: internal error while getting feature flags",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					nil, errors.New("internal error"))
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{},
				RequestedAt: 0,
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "err: internal error while getting segments users",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{multiFeatures[1]},
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-2", envID).Return(
					nil, errors.New("internal error"))
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegmentUsers(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{},
				RequestedAt: 0,
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "err: internal error while getting segment",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{multiFeatures[1]},
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-2", envID).Return(
					nil, errors.New("internal error"))
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegmentUsers(gomock.Any(), gomock.Any()).Return(
					&featureproto.ListSegmentUsersResponse{
						Users: []*featureproto.SegmentUser{
							{
								SegmentId: "segment-id-2",
							},
						},
					}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().GetSegment(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{},
				RequestedAt: 0,
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success: zero feature",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{},
					}, nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{},
				RequestedAt: 0,
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      make([]*featureproto.SegmentUsers, 0),
				DeletedSegmentIds: make([]string, 0),
				RequestedAt:       timeNow.Unix(),
				ForceUpdate:       true,
			},
			expectedErr: nil,
		},
		{
			desc: "success: no segments",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: singleFeature,
					}, nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{},
				RequestedAt: 0,
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      make([]*featureproto.SegmentUsers, 0),
				DeletedSegmentIds: make([]string, 0),
				RequestedAt:       timeNow.Unix(),
				ForceUpdate:       true,
			},
			expectedErr: nil,
		},
		{
			desc: "success: request at older than 30 days",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-2", envID).Return(
					multiSegmentUsers[0], nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-5", envID).Return(
					multiSegmentUsers[2], nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{"segment-id-1", "segment-id-2", "segment-id-3"},
				RequestedAt: timeNow.Add(-31 * 24 * time.Hour).Unix(),
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers: []*featureproto.SegmentUsers{
					multiSegmentUsers[0],
					multiSegmentUsers[2],
				},
				DeletedSegmentIds: make([]string, 0),
				RequestedAt:       timeNow.Unix(),
				ForceUpdate:       true,
			},
			expectedErr: nil,
		},
		{
			desc: "success: return deleted segment ids",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-2", envID).Return(
					multiSegmentUsers[0], nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-5", envID).Return(
					multiSegmentUsers[2], nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{"segment-id-1", "segment-id-2", "segment-id-3"},
				RequestedAt: timeNow.Add(-10 * time.Minute).Unix(),
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      make([]*featureproto.SegmentUsers, 0),
				DeletedSegmentIds: []string{"segment-id-1", "segment-id-3"},
				RequestedAt:       timeNow.Unix(),
				ForceUpdate:       false,
			},
			expectedErr: nil,
		},
		{
			desc: "success: return updated segment users",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-2", envID).Return(
					multiSegmentUsers[0], nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-5", envID).Return(
					multiSegmentUsers[2], nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{"segment-id-1", "segment-id-2", "segment-id-3"},
				RequestedAt: timeNow.Add(-40 * time.Minute).Unix(),
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      []*featureproto.SegmentUsers{multiSegmentUsers[0]},
				DeletedSegmentIds: []string{"segment-id-1", "segment-id-3"},
				RequestedAt:       timeNow.Unix(),
				ForceUpdate:       false,
			},
			expectedErr: nil,
		},
		{
			desc: "success: nothing to update",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envID},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envID).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-2", envID).Return(
					multiSegmentUsers[0], nil)
				gs.segmentUsersCache.(*cachev3mock.MockSegmentUsersCache).EXPECT().Get("segment-id-5", envID).Return(
					multiSegmentUsers[2], nil)
			},
			input: &gwproto.GetSegmentUsersRequest{
				SegmentIds:  []string{"segment-id-2"},
				RequestedAt: timeNow.Add(-20 * time.Minute).Unix(),
				SourceId:    eventproto.SourceId_GO_SERVER,
				SdkVersion:  "v0.0.1",
			},
			expected: &gwproto.GetSegmentUsersResponse{
				SegmentUsers:      make([]*featureproto.SegmentUsers, 0),
				DeletedSegmentIds: make([]string, 0),
				RequestedAt:       timeNow.Unix(),
				ForceUpdate:       false,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{apiKey},
			})
			actual, err := gs.GetSegmentUsers(ctx, p.input)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
			if p.expectedErr != nil {
				assert.Nil(t, actual, "%s", p.desc)
				return
			}
			assert.Equal(t, p.expected.SegmentUsers, actual.SegmentUsers, "%s", p.desc)
			assert.Equal(t, p.expected.DeletedSegmentIds, actual.DeletedSegmentIds, "%s", p.desc)
			assert.GreaterOrEqual(t, actual.RequestedAt, p.expected.RequestedAt, "%s", p.desc)
			assert.Equal(t, p.expected.ForceUpdate, actual.ForceUpdate, "%s", p.desc)
		})
	}
}

func TestGrpcGetFeatureFlagsContextCanceled(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		cancel      bool
		expected    *gwproto.GetFeatureFlagsResponse
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
		actual, err := gs.GetFeatureFlags(ctx, &gwproto.GetFeatureFlagsRequest{})
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGrpcGetFeatureFlags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	timeNow := time.Now()
	apiKey := "api-key-id"
	envNamespace := "ns0"
	tag := "tag"
	singleFeature := []*featureproto.Feature{
		{
			Id:        "feature-id-1",
			Version:   1,
			Tags:      []string{tag},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
	}
	multiFeatures := []*featureproto.Feature{
		{
			Id:        "feature-id-1",
			Version:   1,
			Tags:      []string{tag},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
		{
			Id:        "feature-id-2",
			Version:   1,
			Tags:      []string{},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
		{
			Id:        "feature-id-3",
			Version:   1,
			Tags:      []string{},
			UpdatedAt: timeNow.Add(-20 * time.Minute).Unix(),
		},
		{
			Id:        "feature-id-4",
			Version:   1,
			Tags:      []string{tag},
			UpdatedAt: timeNow.Add(-secondsToReturnAllFlags * time.Second).Unix(),
			Archived:  true,
		},
		{
			Id:        "feature-id-5",
			Version:   1,
			Tags:      []string{tag},
			UpdatedAt: timeNow.Add(-10 * time.Minute).Unix(),
			Archived:  true,
		},
	}
	// Calculate expected feature flags IDs dynamically using UpdatedAt
	multiFeaturesID := evaluation.GenerateFeaturesID(multiFeatures[:3])
	singleFeatureID := evaluation.GenerateFeaturesID(singleFeature)
	patterns := []struct {
		desc        string
		setup       func(*grpcGatewayService)
		input       *gwproto.GetFeatureFlagsRequest
		expected    *gwproto.GetFeatureFlagsResponse
		expectedErr error
	}{
		{
			desc: "err: environment api key not found",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					nil, errors.New("internal error"))
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetEnvironmentAPIKey(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.NotFound, "test"))
			},
			input:       &gwproto.GetFeatureFlagsRequest{Tag: "test", FeatureFlagsId: ""},
			expected:    nil,
			expectedErr: ErrInvalidAPIKey,
		},
		{
			desc: "err: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			input:       &gwproto.GetFeatureFlagsRequest{Tag: "test", FeatureFlagsId: ""},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "err: source id is required",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            "test",
				FeatureFlagsId: "",
				SourceId:       eventproto.SourceId_UNKNOWN,
			},
			expected:    nil,
			expectedErr: ErrSourceIDRequired,
		},
		{
			desc: "err: sdk version is required",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            "test",
				FeatureFlagsId: "",
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "",
			},
			expected:    nil,
			expectedErr: ErrSDKVersionRequired,
		},
		{
			desc: "err: internal error while getting feature flags",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					nil, errors.New("internal error"))
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            "test",
				FeatureFlagsId: "",
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "zero feature",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: []*featureproto.Feature{},
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            "test",
				FeatureFlagsId: "",
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         "",
				Features:               []*featureproto.Feature{},
				ArchivedFeatureFlagIds: make([]string, 0),
				RequestedAt:            timeNow.Unix(),
				ForceUpdate:            false,
			},
			expectedErr: nil,
		},
		{
			desc: "success: with no tag and no feature flags ID",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            "",
				FeatureFlagsId: "",
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId: multiFeaturesID,
				Features: []*featureproto.Feature{
					multiFeatures[0],
					multiFeatures[1],
					multiFeatures[2],
				},
				ArchivedFeatureFlagIds: make([]string, 0),
				RequestedAt:            timeNow.Unix(),
				ForceUpdate:            true,
			},
			expectedErr: nil,
		},
		{
			desc: "success: with no tag and with same feature flags ID",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            "",
				FeatureFlagsId: multiFeaturesID,
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         multiFeaturesID,
				Features:               []*featureproto.Feature{},
				ArchivedFeatureFlagIds: make([]string, 0),
				RequestedAt:            timeNow.Unix(),
				ForceUpdate:            false,
			},
			expectedErr: nil,
		},
		{
			desc: "success: with no tag and with different feature flags ID, and with old requested at",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            "",
				FeatureFlagsId: "random-id",
				RequestedAt:    timeNow.Add(-20 * time.Minute).Unix(),
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId: multiFeaturesID,
				Features: []*featureproto.Feature{
					multiFeatures[0],
					multiFeatures[1],
					multiFeatures[2],
				},
				ArchivedFeatureFlagIds: []string{multiFeatures[4].Id},
				RequestedAt:            timeNow.Unix(),
				ForceUpdate:            false,
			},
			expectedErr: nil,
		},
		{
			desc: "success: with tag and no feature flags ID",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            tag,
				FeatureFlagsId: "",
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         singleFeatureID,
				Features:               singleFeature,
				RequestedAt:            timeNow.Unix(),
				ArchivedFeatureFlagIds: make([]string, 0),
				ForceUpdate:            true,
			},
			expectedErr: nil,
		},
		{
			desc: "success: with tag and same feature flags ID",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: singleFeature,
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            tag,
				FeatureFlagsId: singleFeatureID,
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         singleFeatureID,
				Features:               []*featureproto.Feature{},
				ArchivedFeatureFlagIds: make([]string, 0),
				RequestedAt:            timeNow.Unix(),
				ForceUpdate:            false,
			},
			expectedErr: nil,
		},
		{
			desc: "success: with tag and with different feature flags ID, and with no requested at",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            tag,
				FeatureFlagsId: "random-id",
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         singleFeatureID,
				Features:               singleFeature,
				ArchivedFeatureFlagIds: make([]string, 0),
				RequestedAt:            timeNow.Unix(),
				ForceUpdate:            true,
			},
			expectedErr: nil,
		},
		{
			desc: "success: with tag and with different feature flags ID, and with old requested at",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(apiKey).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: envNamespace},
						ApiKey: &accountproto.APIKey{
							Id:       apiKey,
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(envNamespace).Return(
					&featureproto.Features{
						Features: multiFeatures,
					}, nil)
			},
			input: &gwproto.GetFeatureFlagsRequest{
				Tag:            tag,
				FeatureFlagsId: "random-id",
				RequestedAt:    timeNow.Add(-10 * time.Minute).Unix(),
				SourceId:       eventproto.SourceId_GO_SERVER,
				SdkVersion:     "v0.0.1",
			},
			expected: &gwproto.GetFeatureFlagsResponse{
				FeatureFlagsId:         singleFeatureID,
				Features:               make([]*featureproto.Feature, 0),
				ArchivedFeatureFlagIds: []string{multiFeatures[4].Id},
				RequestedAt:            timeNow.Unix(),
				ForceUpdate:            false,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{apiKey},
			})
			actual, err := gs.GetFeatureFlags(ctx, p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
	ueidFromAndroidFeatures := evaluation.UserEvaluationsID(userID, nil, androidFeatures)
	ueidWithDataFromAndroidFeatures := evaluation.UserEvaluationsID(userID, userMetadata, androidFeatures)

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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListSegmentUsers(gomock.Any(), gomock.Any()).Return(
					&featureproto.ListSegmentUsersResponse{}, nil)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().GetSegment(gomock.Any(), gomock.Any()).Return(
					&featureproto.GetSegmentResponse{Segment: &featureproto.Segment{}}, nil)
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "user-id-1"}},
			expected: &gwproto.GetEvaluationsResponse{
				State: featureproto.UserEvaluations_FULL,
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].Id, evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"), p.desc)
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:             evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"),
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
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].Id, evaluation.EvaluationID("feature-id-1", int32(2), "user-id-1"), p.desc)
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Role:     accountproto.APIKey_SDK_CLIENT,
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
				Id:             evaluation.EvaluationID("feature-id-2", int32(2), "user-id-2"),
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

	bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
		Timestamp: time.Now().Unix(),
		GoalId:    "goal-id-1",
		UserId:    "user-id-1",
		User: &userproto.User{
			Id: "user-id-1",
		},
	})
	if err != nil {
		t.Fatal("could not serialize goal event")
	}
	bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
		Timestamp:   time.Now().Unix(),
		FeatureId:   "feature-id-1",
		VariationId: "variation-id-1",
		User: &userproto.User{
			Id: "user-id-1",
		},
		Reason: &featureproto.Reason{
			Type: featureproto.Reason_DEFAULT,
		},
	})
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
							Id:       "api-key-id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:       "api-key-id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:       "api-key-id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							Id:       "api-key-id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
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
							TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.GoalEvent",
							Value:   bGoalEvent,
						},
					},
					{
						Id: uuid1,
						Event: &any.Any{
							TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
							Value:   bEvaluationEvent,
						},
					},
					{
						Id: uuid2,
						Event: &any.Any{
							TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.MetricsEvent",
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
		featureClient:            featureclientmock.NewMockClient(mockController),
		accountClient:            accountclientmock.NewMockClient(mockController),
		pushClient:               pushclientmock.NewMockClient(mockController),
		codeRefClient:            coderefclientmock.NewMockClient(mockController),
		auditLogClient:           auditlogclientmock.NewMockClient(mockController),
		autoOpsClient:            autoopsclientmock.NewMockClient(mockController),
		goalPublisher:            publishermock.NewMockPublisher(mockController),
		tagClient:                tagclientmock.NewMockClient(mockController),
		teamClient:               teamclientmock.NewMockClient(mockController),
		notificationClient:       notificationclientmock.NewMockClient(mockController),
		experimentClient:         experimentclientmock.NewMockClient(mockController),
		eventCounterClient:       eventcounterclientmock.NewMockClient(mockController),
		environmentClient:        environmentclientmock.NewMockClient(mockController),
		userPublisher:            publishermock.NewMockPublisher(mockController),
		evaluationPublisher:      publishermock.NewMockPublisher(mockController),
		featuresCache:            cachev3mock.NewMockFeaturesCache(mockController),
		segmentUsersCache:        cachev3mock.NewMockSegmentUsersCache(mockController),
		environmentAPIKeyCache:   cachev3mock.NewMockEnvironmentAPIKeyCache(mockController),
		apiKeyLastUsedInfoCacher: sync.Map{},
		opts:                     &defaultOptions,
		logger:                   logger,
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

func TestGrpcDebugEvaluateFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.DebugEvaluateFeaturesResponse
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
					}, nil,
				)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: debug evaluate features error",
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().DebugEvaluateFeatures(
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().DebugEvaluateFeatures(
					gomock.Any(), gomock.Any(),
				).Return(&featureproto.DebugEvaluateFeaturesResponse{
					Evaluations: []*featureproto.Evaluation{
						{
							FeatureId: "feature-id-1",
							UserId:    "user-id-1",
							Reason: &featureproto.Reason{
								Type: featureproto.Reason_DEFAULT,
							},
						},
					},
					ArchivedFeatureIds: []string{"feature-id-2"},
				}, nil)
			},
			expected: &gwproto.DebugEvaluateFeaturesResponse{
				Evaluations: []*featureproto.Evaluation{
					{
						FeatureId: "feature-id-1",
						UserId:    "user-id-1",
						Reason: &featureproto.Reason{
							Type: featureproto.Reason_DEFAULT,
						},
					},
				},
				ArchivedFeatureIds: []string{"feature-id-2"},
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
			actual, err := gs.DebugEvaluateFeatures(ctx, &gwproto.DebugEvaluateFeaturesRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetFeatureResponse
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
			desc: "fail: getFeature error",
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					&feature.GetFeatureResponse{Feature: &featureproto.Feature{Id: "id-0", Enabled: true}},
					nil)
			},
			expected: &gwproto.GetFeatureResponse{
				Feature: &featureproto.Feature{Id: "id-0", Enabled: true},
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
			actual, err := gs.GetFeature(ctx, &gwproto.GetFeatureRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcCreateFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		req         *gwproto.CreateFeatureRequest
		expected    *gwproto.CreateFeatureResponse
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
			desc: "fail: create feature error",
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().CreateFeature(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().CreateFeature(gomock.Any(), gomock.Any()).Return(
					&feature.CreateFeatureResponse{
						Feature: &featureproto.Feature{Id: "id-0", Enabled: true},
					},
					nil)
			},
			expected: &gwproto.CreateFeatureResponse{
				Feature: &featureproto.Feature{Id: "id-0", Enabled: true},
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
			actual, err := gs.CreateFeature(ctx, &gwproto.CreateFeatureRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcUpdateFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.UpdateFeatureResponse
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
			desc: "fail: update feature error",
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().UpdateFeature(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().UpdateFeature(gomock.Any(), gomock.Any()).Return(
					&feature.UpdateFeatureResponse{Feature: &feature.Feature{Id: "fid"}}, nil)
			},
			expected:    &gwproto.UpdateFeatureResponse{Feature: &feature.Feature{Id: "fid"}},
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
			actual, err := gs.UpdateFeature(ctx, &gwproto.UpdateFeatureRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
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

func TestGrpcListFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListFeaturesResponse
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
			desc: "fail: listFeatures error",
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
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
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					&featureproto.ListFeaturesResponse{Features: []*featureproto.Feature{
						{
							Id:      "id-0",
							Enabled: true,
						},
					}}, nil)
			},
			expected: &gwproto.ListFeaturesResponse{Features: []*featureproto.Feature{
				{Id: "id-0", Enabled: true},
			}},
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
			actual, err := gs.ListFeatures(ctx, &gwproto.ListFeaturesRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestObfuscateString(t *testing.T) {
	tests := []struct {
		desc       string // Description comes first
		input      string
		showLength int
		expected   string
	}{
		{
			desc:       "String of 32 characters, showing the first 4 and last 4 characters.",
			input:      "12345678901234567890123456789012",
			showLength: 4,
			expected:   "1234....9012",
		},
		{
			desc:       "String of 10 characters, showing the first 3 and last 3 characters.",
			input:      "abcdefghij",
			showLength: 3,
			expected:   "abc....hij",
		},
		{
			desc:       "String shorter than showLength*2. Should not obfuscate.",
			input:      "abcd",
			showLength: 4,
			expected:   "abcd", // No obfuscation needed
		},
		{
			desc:       "String of exactly twice the showLength (8 characters), showing the first 3 and last 3 characters.",
			input:      "abcdefgh",
			showLength: 3,
			expected:   "abc....fgh",
		},
		{
			desc:       "String of 5 characters, showing 2 characters from both the start and end.",
			input:      "abcde",
			showLength: 2,
			expected:   "ab....de",
		},
		{
			desc:       "String with special characters, ensuring the function works with non-alphanumeric characters.",
			input:      "@bcde!@#f",
			showLength: 3,
			expected:   "@bc....@#f",
		},
		{
			desc:       "String longer than showLength*2, showing the first 2 and last 2 characters.",
			input:      "1234567890abcdef1234567890",
			showLength: 2,
			expected:   "12....90",
		},
		{
			desc:       "String exactly equal to 2*showLength, obfuscate with the middle part replaced by dots.",
			input:      "12345678",
			showLength: 4,
			expected:   "12345678",
		},
		{
			desc:       "Single character string. Should not obfuscate as it's too short.",
			input:      "a",
			showLength: 1,
			expected:   "a",
		},
	}

	// Loop through each test case
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual := obfuscateString(tt.input, tt.showLength)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
