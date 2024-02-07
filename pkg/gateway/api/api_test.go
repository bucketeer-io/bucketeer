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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/pkg/cache/v3/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const dummyURL = "http://example.com"

func TestNewGatewayService(t *testing.T) {
	t.Parallel()
	g := NewGatewayService(nil, nil, nil, nil, nil, nil, nil, nil)
	assert.IsType(t, &gatewayService{}, g)
}

func TestGetEnvironmentAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*gatewayService)
		auth        string
		expected    *accountproto.EnvironmentAPIKey
		expectedErr error
	}{
		{
			desc: "exists in cache",
			setup: func(gs *gatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey:      &accountproto.APIKey{Id: "id-0"},
					}, nil)
			},
			auth: "test-key",
			expected: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey:      &accountproto.APIKey{Id: "id-0"},
			},
			expectedErr: nil,
		},
		{
			desc: "ErrInvalidAPIKey",
			setup: func(gs *gatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.NotFound, "test"))
			},
			auth:        "test-key",
			expected:    nil,
			expectedErr: errInvalidAPIKey,
		},
		{
			desc: "ErrInternal",
			setup: func(gs *gatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Unknown, "test"))
			},
			auth:        "test-key",
			expected:    nil,
			expectedErr: errInternal,
		},
		{
			desc: "success",
			setup: func(gs *gatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAPIKeyBySearchingAllEnvironments(gomock.Any(), gomock.Any()).Return(
					&accountproto.GetAPIKeyBySearchingAllEnvironmentsResponse{EnvironmentApiKey: &accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey:      &accountproto.APIKey{Id: "id-0"},
					}}, nil)
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Put(gomock.Any()).Return(nil)
			},
			auth: "test-key",
			expected: &accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey:      &accountproto.APIKey{Id: "id-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		req := httptest.NewRequest(
			"POST",
			dummyURL,
			nil,
		)
		req.Header.Add(authorizationKey, p.auth)
		actual, err := gs.findEnvironmentAPIKey(context.Background(), req)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGetEnvironmentAPIKeyFromCache(t *testing.T) {
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

func TestCheckEnvironmentAPIKey(t *testing.T) {
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
			expected:  errBadRole,
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
			expected:  errDisabledAPIKey,
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
			expected:  errDisabledAPIKey,
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
	gs := gatewayService{}
	for _, p := range patterns {
		actual := gs.checkEnvironmentAPIKey(p.inputEnvAPIKey, p.inputRole)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func TestValidateGetEvaluationsRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    *getEvaluationsRequest
		expected error
	}{
		{
			desc:     "tag is empty",
			input:    &getEvaluationsRequest{},
			expected: errTagRequired,
		},
		{
			desc:     "user is empty",
			input:    &getEvaluationsRequest{Tag: "test"},
			expected: errUserRequired,
		},
		{
			desc:     "user ID is empty",
			input:    &getEvaluationsRequest{Tag: "test", User: &userproto.User{}},
			expected: errUserIDRequired,
		},
		{
			desc:  "pass",
			input: &getEvaluationsRequest{Tag: "test", User: &userproto.User{Id: "id"}},
		},
	}
	gs := gatewayService{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := gs.validateGetEvaluationsRequest(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestValidateGetEvaluationRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		input    *getEvaluationRequest
		expected error
	}{
		{
			desc:     "tag is empty",
			input:    &getEvaluationRequest{},
			expected: errTagRequired,
		},
		{
			desc:     "user is empty",
			input:    &getEvaluationRequest{Tag: "test"},
			expected: errUserRequired,
		},
		{
			desc:     "user ID is empty",
			input:    &getEvaluationRequest{Tag: "test", User: &userproto.User{}},
			expected: errUserIDRequired,
		},
		{
			desc:     "feature ID is empty",
			input:    &getEvaluationRequest{Tag: "test", User: &userproto.User{Id: "id"}},
			expected: errFeatureIDRequired,
		},
		{
			desc:  "pass",
			input: &getEvaluationRequest{Tag: "test", User: &userproto.User{Id: "id"}, FeatureID: "id"},
		},
	}
	gs := gatewayService{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := gs.validateGetEvaluationRequest(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGetFeaturesFromCache(t *testing.T) {
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
		gs := gatewayService{featuresCache: mtfc}
		actual, err := gs.getFeaturesFromCache(context.Background(), p.environmentId)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func TestGetFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	now := time.Now()
	twentyNineDaysAgo := now.Add(-29 * 24 * time.Hour)
	thirtyOneDaysAgo := now.Add(-31 * 24 * time.Hour)

	patterns := []struct {
		desc          string
		setup         func(*gatewayService)
		environmentId string
		expected      []*featureproto.Feature
		expectedErr   error
	}{
		{
			desc: "exists in redis",
			setup: func(gs *gatewayService) {
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
			setup: func(gs *gatewayService) {
				gs.featuresCache.(*cachev3mock.MockFeaturesCache).EXPECT().Get(gomock.Any()).Return(
					nil, cache.ErrNotFound)
				gs.featureClient.(*featureclientmock.MockClient).EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("test"))
			},
			environmentId: "ns0",
			expected:      nil,
			expectedErr:   errInternal,
		},
		{
			desc: "success",
			setup: func(gs *gatewayService) {
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
			setup: func(gs *gatewayService) {
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
			setup: func(gs *gatewayService) {
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
			gs := newGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			actual, err := gs.getFeatures(context.Background(), p.environmentId)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGetEvaluationsContextCanceled(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		cancel      bool
		expected    *getEvaluationsResponse
		expectedErr error
	}{
		{
			desc:        "error: context canceled",
			cancel:      true,
			expected:    nil,
			expectedErr: errContextCanceled,
		},
		{
			desc:        "error: missing API key",
			cancel:      false,
			expected:    nil,
			expectedErr: errMissingAPIKey,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		req := httptest.NewRequest(
			"POST",
			dummyURL,
			nil,
		)
		ctx, cancel := context.WithCancel(req.Context())
		if p.cancel {
			cancel()
		} else {
			defer cancel()
		}
		actual := httptest.NewRecorder()
		gs.getEvaluations(actual, req.WithContext(ctx))
		assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
	}
}

func TestGetEvaluationsValidation(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*gatewayService)
		input       *http.Request
		expected    *getEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "missing tag",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected:    nil,
			expectedErr: errTagRequired,
		},
		{
			desc: "missing user id",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag: "test",
					},
				),
			),
			expected:    nil,
			expectedErr: errUserRequired,
		},
		{
			desc: "success",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:  "test",
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected: &getEvaluationsResponse{
				Evaluations:       emptyUserEvaluationsForREST(t),
				UserEvaluationsID: "no_evaluations",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		actual := httptest.NewRecorder()
		p.input.Header.Add(authorizationKey, "test-key")
		gs.getEvaluations(actual, p.input)
		if actual.Code != http.StatusOK {
			assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
			continue
		}
		var respBody getEvaluationsResponse
		decoded := decodeSuccessResponse(t, actual.Body)
		err := json.Unmarshal(decoded, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, p.expected, &respBody, "%s", p.desc)
	}
}

func TestGetEvaluationsZeroFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*gatewayService)
		input       *http.Request
		expected    *getEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "zero feature",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:               "test",
						User:              &userproto.User{Id: "id-0"},
						UserEvaluationsID: "evaluation-id",
					},
				),
			),
			expected: &getEvaluationsResponse{
				Evaluations:       emptyUserEvaluationsForREST(t),
				UserEvaluationsID: "no_evaluations",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		p.setup(gs)
		actual := httptest.NewRecorder()
		p.input.Header.Add(authorizationKey, "test-key")
		gs.getEvaluations(actual, p.input)
		var respBody getEvaluationsResponse
		decoded := decodeSuccessResponse(t, actual.Body)
		err := json.Unmarshal(decoded, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, p.expected, &respBody, "%s", p.desc)
	}
}

func TestGetEvaluationsUserEvaluationsID(t *testing.T) {
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
					Value: "true",
				},
				{
					Id:    vID2,
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
					Value: "true",
				},
				{
					Id:    vID4,
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
	userID := "user-id-0"
	userMetadata := map[string]string{"b": "value-b", "c": "value-c", "a": "value-a", "d": "value-d"}
	ueid := featuredomain.UserEvaluationsID(userID, nil, features)
	ueidWithData := featuredomain.UserEvaluationsID(userID, userMetadata, features)

	patterns := []struct {
		desc                      string
		setup                     func(*gatewayService)
		input                     *http.Request
		expected                  *getEvaluationsResponse
		expectedErr               error
		expectedEvaluationsAssert func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			desc: "user evaluations id not set",
			setup: func(gs *gatewayService) {
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
						Features: features,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag: "test",
						User: &userproto.User{
							Id:   userID,
							Data: userMetadata,
						},
					},
				),
			),
			expected: &getEvaluationsResponse{
				UserEvaluationsID: ueidWithData,
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user evaluations id is the same",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag: "test",
						User: &userproto.User{
							Id:   userID,
							Data: userMetadata,
						},
						UserEvaluationsID: featuredomain.UserEvaluationsID(userID, userMetadata, multiFeatures),
					},
				),
			),
			expected: &getEvaluationsResponse{
				UserEvaluationsID: featuredomain.UserEvaluationsID(userID, userMetadata, multiFeatures),
				Evaluations:       &featureproto.UserEvaluations{},
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user evaluations id is different",
			setup: func(gs *gatewayService) {
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
						Features: features,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag: "test",
						User: &userproto.User{
							Id:   userID,
							Data: userMetadata,
						},
						UserEvaluationsID: "evaluation-id",
					},
				),
			),
			expected: &getEvaluationsResponse{
				UserEvaluationsID: ueidWithData,
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user_with_no_metadata_and_the_id_is_same",
			setup: func(gs *gatewayService) {
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
						Features: features,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:               "test",
						User:              &userproto.User{Id: userID},
						UserEvaluationsID: ueid,
					},
				),
			),

			expected: &getEvaluationsResponse{
				UserEvaluationsID: ueid,
				Evaluations:       &featureproto.UserEvaluations{},
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
		{
			desc: "user_with_no_metadata_and_the_id_is_different",
			setup: func(gs *gatewayService) {
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
						Features: features,
					}, nil)
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:               "test",
						User:              &userproto.User{Id: userID},
						UserEvaluationsID: "evaluation-id",
					},
				),
			),
			expected: &getEvaluationsResponse{
				UserEvaluationsID: ueid,
			},
			expectedErr:               nil,
			expectedEvaluationsAssert: assert.NotNil,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		if p.setup != nil {
			p.setup(gs)
		}
		actual := httptest.NewRecorder()
		p.input.Header.Add(authorizationKey, "test-key")
		gs.getEvaluations(actual, p.input)
		if actual.Code != http.StatusOK {
			assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
			continue
		}
		var respBody getEvaluationsResponse
		decoded := decodeSuccessResponse(t, actual.Body)
		err := json.Unmarshal(decoded, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, p.expected.UserEvaluationsID, respBody.UserEvaluationsID, "%s", p.desc)
		p.expectedEvaluationsAssert(t, respBody.Evaluations, "%s", p.desc)
	}
}

func testGetEvaluationsNoSegmentList(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	vID1 := newUUID(t)
	vID2 := newUUID(t)
	vID3 := newUUID(t)
	vID4 := newUUID(t)

	patterns := []struct {
		desc        string
		setup       func(*gatewayService)
		input       *http.Request
		expected    *getEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "state: full",
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    vID2,
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
										Value: "true",
									},
									{
										Id:    vID4,
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:  "ios",
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected: &getEvaluationsResponse{
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
		gs := newGatewayServiceWithMock(t, mockController)
		if p.setup != nil {
			p.setup(gs)
		}
		actual := httptest.NewRecorder()
		p.input.Header.Add(authorizationKey, "test-key")
		gs.getEvaluations(actual, p.input)
		if actual.Code != http.StatusOK {
			assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
			return
		}
		var respBody getEvaluationsResponse
		decoded := decodeSuccessResponse(t, actual.Body)
		err := json.Unmarshal(decoded, &respBody)
		assert.NoError(t, err)
		ev := p.expected.Evaluations.Evaluations
		av := respBody.Evaluations.Evaluations
		assert.Equal(t, len(ev), len(av), "%s", p.desc)
		assert.Equal(t, ev[0].VariationId, av[0].VariationId, "%s", p.desc)
		assert.Equal(t, ev[1].VariationId, av[1].VariationId, "%s", p.desc)
		assert.NotEmpty(t, respBody.UserEvaluationsID, "%s", p.desc)
	}
}

func TestGetEvaluationsEvaluateFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*gatewayService)
		input       *http.Request
		expected    *getEvaluationsResponse
		expectedErr error
	}{
		{
			desc: "errInternal",
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    "variation-b",
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:  "user",
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected:    nil,
			expectedErr: errInternal,
		},
		{
			desc: "state: full, evaluate features list segment from cache",
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    "variation-b",
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:  "test",
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected: &getEvaluationsResponse{
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							VariationId: "variation-b",
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
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    "variation-b",
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:  "test",
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected: &getEvaluationsResponse{
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							VariationId: "variation-b",
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
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    "variation-b",
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:  "test",
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected: &getEvaluationsResponse{
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							VariationId: "variation-b",
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
			setup: func(gs *gatewayService) {
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
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-a",
										Value: "true",
									},
									{
										Id:    "variation-b",
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
								Id:       "feature-2",
								Archived: true,
								Variations: []*featureproto.Variation{
									{
										Id:    "variation-c",
										Value: "true",
									},
									{
										Id:    "variation-d",
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationsRequest{
						Tag:  "test",
						User: &userproto.User{Id: "id-0"},
					},
				),
			),
			expected: &getEvaluationsResponse{
				Evaluations: &featureproto.UserEvaluations{
					Evaluations: []*featureproto.Evaluation{
						{
							VariationId: "variation-b",
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
			gs := newGatewayServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual := httptest.NewRecorder()
			p.input.Header.Add(authorizationKey, "test-key")
			gs.getEvaluations(actual, p.input)
			if actual.Code != http.StatusOK {
				assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
			} else {
				var respBody getEvaluationsResponse
				decoded := decodeSuccessResponse(t, actual.Body)
				err := json.Unmarshal(decoded, &respBody)
				assert.NoError(t, err)
				assert.Equal(t, len(respBody.Evaluations.Evaluations), 1, "%s", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].VariationId, "variation-b", "%s", p.desc)
				assert.Equal(t, p.expected.Evaluations.Evaluations[0].Reason, respBody.Evaluations.Evaluations[0].Reason, p.desc)
				assert.NotEmpty(t, respBody.UserEvaluationsID, "%s", p.desc)
			}
		})
	}
}

func TestGetEvaluation(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc              string
		setup             func(*gatewayService)
		input             *http.Request
		expectedFeatureID string
		expectedErr       error
	}{
		{
			desc: "errFeatureNotFound",
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    "variation-b",
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
										Value: "true",
									},
									{
										Id:    "variation-d",
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationRequest{
						Tag:       "test",
						User:      &userproto.User{Id: "id-0"},
						FeatureID: "feature-id-3",
					},
				),
			),
			expectedFeatureID: "",
			expectedErr:       errFeatureNotFound,
		},
		{
			desc: "errInternal",
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    "variation-b",
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
										Value: "true",
									},
									{
										Id:    "variation-d",
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationRequest{
						Tag:       "test",
						User:      &userproto.User{Id: "id-0"},
						FeatureID: "feature-id-2",
					},
				),
			),
			expectedFeatureID: "",
			expectedErr:       errInternal,
		},
		{
			desc: "return evaluation",
			setup: func(gs *gatewayService) {
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
										Value: "true",
									},
									{
										Id:    "variation-b",
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
										Id:    "variation-a",
										Value: "true",
									},
									{
										Id:    "variation-b",
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
				gs.userPublisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			},
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					getEvaluationRequest{
						Tag:       "test",
						User:      &userproto.User{Id: "id-0"},
						FeatureID: "feature-id-2",
					},
				),
			),
			expectedFeatureID: "feature-id-2",
			expectedErr:       nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGatewayServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(gs)
			}
			actual := httptest.NewRecorder()
			p.input.Header.Add(authorizationKey, "test-key")
			gs.getEvaluation(actual, p.input)
			if actual.Code != http.StatusOK {
				assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
				return
			}
			var respBody getEvaluationResponse
			decoded := decodeSuccessResponse(t, actual.Body)
			err := json.Unmarshal(decoded, &respBody)
			assert.NoError(t, err)
			assert.Equal(t, p.expectedFeatureID, respBody.Evaluation.FeatureId)
		})
	}
}

func TestRegisterEventsContextCanceled(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		cancel      bool
		expectedErr error
	}{
		{
			desc:        "error: context canceled",
			cancel:      true,
			expectedErr: errContextCanceled,
		},
		{
			desc:        "error: missing API key",
			cancel:      false,
			expectedErr: errMissingAPIKey,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		req := httptest.NewRequest(
			"POST",
			dummyURL,
			nil,
		)
		ctx, cancel := context.WithCancel(req.Context())
		if p.cancel {
			cancel()
		} else {
			defer cancel()
		}
		actual := httptest.NewRecorder()
		gs.registerEvents(
			actual,
			req.WithContext(ctx),
		)
		assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
	}
}

func TestRegisterEvents(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	bGoalEvent, err := protojson.Marshal(&eventproto.GoalEvent{Timestamp: time.Now().Unix()})
	if err != nil {
		t.Fatal("could not serialize goal event")
	}
	bEvaluationEvent, err := protojson.Marshal(&eventproto.EvaluationEvent{Timestamp: time.Now().Unix()})
	if err != nil {
		t.Fatal("could not serialize evaluation event")
	}
	bLatencyEvent, err := json.Marshal(&latencyMetricsEvent{
		Labels:   map[string]string{"tag": "test", "status": "success"},
		Duration: time.Duration(1),
	})
	if err != nil {
		t.Fatal("could not serialize goal event")
	}
	bMetricsEvent, err := json.Marshal(&metricsEvent{
		Timestamp: time.Now().Unix(),
		Event:     json.RawMessage(string(bLatencyEvent)),
		Type:      latencyMetricsEventType,
	})
	if err != nil {
		t.Fatal("could not serialize metrics event")
	}
	uuid0 := newUUID(t)
	uuid1 := newUUID(t)
	uuid2 := newUUID(t)

	patterns := []struct {
		desc        string
		setup       func(*gatewayService)
		input       *http.Request
		expected    *registerEventsResponse
		expectedErr error
	}{
		{
			desc:  "error: invalid http method",
			setup: nil,
			input: httptest.NewRequest(
				"GET",
				dummyURL,
				nil,
			),
			expectedErr: errInvalidHttpMethod,
		},
		{
			desc: "error: body is nil",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				nil,
			),
			expectedErr: errBodyRequired,
		},
		{
			desc: "error: zero event",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					registerEventsRequest{},
				),
			),
			expectedErr: errMissingEvents,
		},
		{
			desc: "error: ErrMissingEventID",
			setup: func(gs *gatewayService) {
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
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					registerEventsRequest{
						Events: []event{
							{
								ID: "",
							},
						},
					},
				),
			),
			expectedErr: errMissingEventID,
		},
		{
			desc: "error: invalid message type",
			setup: func(gs *gatewayService) {
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
				gs.metricsPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					registerEventsRequest{
						Events: []event{
							{
								ID:    uuid0,
								Event: json.RawMessage(string(bGoalEvent)),
								Type:  8,
							},
						},
					},
				),
			),
			expected: &registerEventsResponse{
				Errors: map[string]*registerEventsResponseError{
					uuid0: {
						Retriable: false,
						Message:   errInvalidType.Error(),
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success",
			setup: func(gs *gatewayService) {
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
				gs.metricsPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(
					nil).MaxTimes(1)
			},
			input: httptest.NewRequest(
				"POST",
				dummyURL,
				renderBody(
					t,
					registerEventsRequest{
						Events: []event{
							{
								ID:    uuid0,
								Event: json.RawMessage(bGoalEvent),
								Type:  GoalEventType,
							},
							{
								ID:    uuid1,
								Event: json.RawMessage(bEvaluationEvent),
								Type:  EvaluationEventType,
							},
							{
								ID:    uuid2,
								Event: json.RawMessage(bMetricsEvent),
								Type:  MetricsEventType,
							},
						},
					},
				),
			),
			expected:    &registerEventsResponse{Errors: map[string]*registerEventsResponseError(nil)},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		if p.setup != nil {
			p.setup(gs)
		}
		actual := httptest.NewRecorder()
		p.input.Header.Add("authorization", "test-key")
		gs.registerEvents(actual, p.input)
		if actual.Code != http.StatusOK {
			assert.Equal(t, newErrResponse(t, p.expectedErr), actual.Body.String(), "%s", p.desc)
			continue
		}
		var respBody registerEventsResponse
		decoded := decodeSuccessResponse(t, actual.Body)
		err := json.Unmarshal(decoded, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, p.expected, &respBody, p.desc)
	}
}

func TestContainsInvalidTimestampError(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc     string
		errs     map[string]*registerEventsResponseError
		expected bool
	}{
		{
			desc: "error: invalid timestamp",
			errs: map[string]*registerEventsResponseError{
				"id-test": {
					Retriable: false,
					Message:   errInvalidTimestamp.Error(),
				},
			},
			expected: true,
		},
		{
			desc: "error: unmarshal failed",
			errs: map[string]*registerEventsResponseError{
				"id-test": {
					Retriable: true,
					Message:   errUnmarshalFailed.Error(),
				},
			},
			expected: false,
		},
		{
			desc:     "error: empty",
			errs:     make(map[string]*registerEventsResponseError),
			expected: false,
		},
	}
	for _, p := range patterns {
		gs := newGatewayServiceWithMock(t, mockController)
		actual := gs.containsInvalidTimestampError(p.errs)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}

func TestGetMetricsEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	bLatencyEvent, err := json.Marshal(&latencyMetricsEvent{
		Labels:   map[string]string{"tag": "test", "status": "success"},
		Duration: time.Duration(1),
	})
	if err != nil {
		t.Fatal("could not serialize goal event")
	}
	ctx := context.TODO()
	patterns := []struct {
		desc        string
		input       metricsEvent
		expected    *eventproto.MetricsEvent
		expectedErr error
	}{
		{
			desc: "error: invalid message type",
			input: metricsEvent{
				Timestamp: time.Now().Unix(),
				Event:     json.RawMessage(string(bLatencyEvent)),
				Type:      0,
			},
			expectedErr: errInvalidType,
		},
		{
			desc: "success",
			input: metricsEvent{
				Timestamp: time.Now().Unix(),
				Event:     json.RawMessage(string(bLatencyEvent)),
				Type:      latencyMetricsEventType,
			},
			expected: &eventproto.MetricsEvent{
				Timestamp: time.Now().Unix(),
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGatewayServiceWithMock(t, mockController)
			bMetricsEvent, err := json.Marshal(p.input)
			assert.NoError(t, err)
			ev := event{
				ID:    newUUID(t),
				Event: json.RawMessage(bMetricsEvent),
				Type:  MetricsEventType,
			}
			event, _, err := gs.getMetricsEvent(ctx, ev)
			if err != nil {
				assert.Equal(t, p.expectedErr, err)
				return
			}
			assert.Equal(t, event.Timestamp, p.expected.Timestamp)
			assert.NotNil(t, event.Event)
		})
	}
}

type successResponse struct {
	Data json.RawMessage `json:"data"`
}

func decodeSuccessResponse(t *testing.T, body *bytes.Buffer) json.RawMessage {
	t.Helper()
	var resp successResponse
	err := json.NewDecoder(body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}
	return resp.Data
}

type failureResponse struct {
	Error errorResponse `json:"error"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newErrResponse(t *testing.T, err error) string {
	t.Helper()
	status, _ := convertToErrStatus(err)
	res := &failureResponse{
		Error: errorResponse{
			Code:    status.GetStatusCode(),
			Message: status.GetErrMessage(),
		},
	}
	encoded, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

type errStatus interface {
	GetErrMessage() string
	GetStatusCode() int
}

func convertToErrStatus(err error) (errStatus, bool) {
	s, ok := err.(errStatus)
	if !ok {
		return nil, false
	}
	return s, true
}

func renderBody(t *testing.T, res interface{}) io.Reader {
	t.Helper()
	encoded, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(encoded)
}

func newGatewayServiceWithMock(t *testing.T, mockController *gomock.Controller) *gatewayService {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &gatewayService{
		featureClient:          featureclientmock.NewMockClient(mockController),
		accountClient:          accountclientmock.NewMockClient(mockController),
		goalPublisher:          publishermock.NewMockPublisher(mockController),
		userPublisher:          publishermock.NewMockPublisher(mockController),
		metricsPublisher:       publishermock.NewMockPublisher(mockController),
		evaluationPublisher:    publishermock.NewMockPublisher(mockController),
		featuresCache:          cachev3mock.NewMockFeaturesCache(mockController),
		segmentUsersCache:      cachev3mock.NewMockSegmentUsersCache(mockController),
		environmentAPIKeyCache: cachev3mock.NewMockEnvironmentAPIKeyCache(mockController),
		opts:                   &defaultOptions,
		logger:                 logger,
	}
}

func emptyUserEvaluationsForREST(t *testing.T) *featureproto.UserEvaluations {
	t.Helper()
	return &featureproto.UserEvaluations{
		Id:                 "no_evaluations",
		Evaluations:        []*featureproto.Evaluation{},
		CreatedAt:          time.Now().Unix(),
		ForceUpdate:        false,
		ArchivedFeatureIds: []string{},
	}
}
