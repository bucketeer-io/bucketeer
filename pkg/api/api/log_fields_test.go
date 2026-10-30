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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	evaluationLogUserID = "user-1"
	// evaluationLogUserAttribute is an end-user attribute value that must
	// never be written to the log.
	evaluationLogUserAttribute = "user-1@example.com"
	evaluationLogFeatureID     = "feature-1"
)

func newEvaluationLogUser() *userproto.User {
	return &userproto.User{
		Id:   evaluationLogUserID,
		Data: map[string]string{"email": evaluationLogUserAttribute},
	}
}

func newObservedLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.DebugLevel)
	return zap.New(core), logs
}

// assertEntriesOmit fails when any recorded entry mentions value in its
// message or in any of its fields.
func assertEntriesOmit(t *testing.T, logs *observer.ObservedLogs, value string) {
	t.Helper()
	for _, entry := range logs.All() {
		assert.NotContains(t, entry.Message, value)
		for key, field := range entry.ContextMap() {
			assert.NotContains(t, fmt.Sprint(field), value, "field %q of %q", key, entry.Message)
		}
	}
}

// requireUserIDEntry checks that the request failure was logged once with
// the user id as a plain field.
func requireUserIDEntry(t *testing.T, logs *observer.ObservedLogs, message, fieldName, userID string) {
	t.Helper()
	entries := logs.FilterMessage(message).All()
	require.Len(t, entries, 1, "expected exactly one %q entry", message)
	assert.Equal(t, userID, entries[0].ContextMap()[fieldName])
}

func expectSDKClientAPIKey(cache *cachev3mock.MockEnvironmentAPIKeyCache, disabled bool) {
	cache.EXPECT().Get(gomock.Any()).Return(
		&accountproto.EnvironmentAPIKey{
			Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
			ApiKey: &accountproto.APIKey{
				Id:       "id-0",
				Role:     accountproto.APIKey_SDK_CLIENT,
				Disabled: disabled,
			},
		}, nil)
}

func expectSingleFixedFeature(cache *cachev3mock.MockFeaturesCache) {
	cache.EXPECT().Get(gomock.Any()).Return(
		&featureproto.Features{
			Features: []*featureproto.Feature{
				{
					Id:      evaluationLogFeatureID,
					Version: 1,
					Variations: []*featureproto.Variation{
						{Id: "variation-a", Name: "true", Value: "true"},
						{Id: "variation-b", Name: "false", Value: "false"},
					},
					DefaultStrategy: &featureproto.Strategy{
						Type:          featureproto.Strategy_FIXED,
						FixedStrategy: &featureproto.FixedStrategy{Variation: "variation-b"},
					},
					Tags: []string{"test"},
				},
			},
		}, nil)
}

func TestGrpcGetEvaluationsLogsUserIDOnly(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setup          func(*grpcGatewayService)
		input          *gwproto.GetEvaluationsRequest
		expectedErr    error
		expectedLogMsg string
		expectedUserID string
	}{
		{
			desc: "err: request check failure is logged with the user id only",
			setup: func(gs *grpcGatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), true)
			},
			input:          &gwproto.GetEvaluationsRequest{Tag: "test", User: newEvaluationLogUser()},
			expectedErr:    ErrDisabledAPIKey,
			expectedLogMsg: "Failed to check GetEvaluations request",
			expectedUserID: evaluationLogUserID,
		},
		{
			// Validation only inspects the user, so the failing request has
			// attributes but no id.
			desc: "err: validation failure is logged without the user attributes",
			setup: func(gs *grpcGatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
			},
			input: &gwproto.GetEvaluationsRequest{
				Tag:  "test",
				User: &userproto.User{Data: map[string]string{"email": evaluationLogUserAttribute}},
			},
			expectedErr:    ErrUserIDRequired,
			expectedLogMsg: "Failed to validate GetEvaluations request",
			expectedUserID: "",
		},
		{
			desc: "success: evaluations are returned",
			setup: func(gs *grpcGatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
				expectSingleFixedFeature(gs.featuresCache.(*cachev3mock.MockFeaturesCache))
			},
			input: &gwproto.GetEvaluationsRequest{Tag: "test", User: newEvaluationLogUser()},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			logger, logs := newObservedLogger()
			gs.logger = logger
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetEvaluations(ctx, p.input)
			assertEntriesOmit(t, logs, evaluationLogUserAttribute)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
				assert.Nil(t, actual)
				requireUserIDEntry(t, logs, p.expectedLogMsg, "userId", p.expectedUserID)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Len(t, actual.Evaluations.Evaluations, 1)
			assert.Equal(t, evaluationLogFeatureID, actual.Evaluations.Evaluations[0].FeatureId)
		})
	}
}

func TestGrpcGetEvaluationLogsUserIDOnly(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setup          func(*grpcGatewayService)
		input          *gwproto.GetEvaluationRequest
		expectedErr    error
		expectedLogMsg string
	}{
		{
			desc: "err: request check failure is logged with the user id only",
			setup: func(gs *grpcGatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), true)
			},
			input: &gwproto.GetEvaluationRequest{
				Tag: "test", User: newEvaluationLogUser(), FeatureId: evaluationLogFeatureID,
			},
			expectedErr:    ErrDisabledAPIKey,
			expectedLogMsg: "Failed to check GetEvaluation request",
		},
		{
			desc: "err: validation failure is logged with the user id only",
			setup: func(gs *grpcGatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
			},
			input:          &gwproto.GetEvaluationRequest{Tag: "test", User: newEvaluationLogUser()},
			expectedErr:    ErrFeatureIDRequired,
			expectedLogMsg: "Failed to validate GetEvaluation request",
		},
		{
			desc: "success: the evaluation is returned",
			setup: func(gs *grpcGatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
				expectSingleFixedFeature(gs.featuresCache.(*cachev3mock.MockFeaturesCache))
			},
			input: &gwproto.GetEvaluationRequest{
				Tag: "test", User: newEvaluationLogUser(), FeatureId: evaluationLogFeatureID,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			logger, logs := newObservedLogger()
			gs.logger = logger
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetEvaluation(ctx, p.input)
			assertEntriesOmit(t, logs, evaluationLogUserAttribute)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
				assert.Nil(t, actual)
				requireUserIDEntry(t, logs, p.expectedLogMsg, "userId", evaluationLogUserID)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, actual)
			assert.Equal(t, evaluationLogFeatureID, actual.Evaluation.FeatureId)
			assert.Equal(t, "variation-b", actual.Evaluation.VariationId)
		})
	}
}

func TestGetEvaluationsLogsUserIDOnly(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setup          func(*gatewayService)
		body           getEvaluationsRequest
		expectedErr    error
		expectedLogMsg string
	}{
		{
			desc: "err: request check failure is logged with the user id only",
			setup: func(gs *gatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), true)
			},
			body:           getEvaluationsRequest{Tag: "test", User: newEvaluationLogUser()},
			expectedErr:    errDisabledAPIKey,
			expectedLogMsg: "Failed to check GetEvaluations request",
		},
		{
			desc: "err: validation failure is logged with the user id only",
			setup: func(gs *gatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
			},
			body:           getEvaluationsRequest{User: newEvaluationLogUser()},
			expectedErr:    errTagRequired,
			expectedLogMsg: "Failed to validate GetEvaluations request",
		},
		{
			desc: "success: evaluations are returned",
			setup: func(gs *gatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
				expectSingleFixedFeature(gs.featuresCache.(*cachev3mock.MockFeaturesCache))
			},
			body: getEvaluationsRequest{Tag: "test", User: newEvaluationLogUser()},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			gs := newGatewayServiceWithMock(t, mockController)
			logger, logs := newObservedLogger()
			gs.logger = logger
			p.setup(gs)
			req := httptest.NewRequest(http.MethodPost, dummyURL, renderBody(t, p.body))
			req.Header.Add(authorizationKey, "test-key")
			rec := httptest.NewRecorder()
			gs.getEvaluations(rec, req)
			assertEntriesOmit(t, logs, evaluationLogUserAttribute)
			if p.expectedErr != nil {
				assert.NotEqual(t, http.StatusOK, rec.Code)
				assert.Equal(t, newErrResponse(t, p.expectedErr), rec.Body.String())
				requireUserIDEntry(t, logs, p.expectedLogMsg, "userId", evaluationLogUserID)
				return
			}
			require.Equal(t, http.StatusOK, rec.Code)
			var resp getEvaluationsResponse
			require.NoError(t, json.Unmarshal(decodeSuccessResponse(t, rec.Body), &resp))
			require.NotNil(t, resp.Evaluations)
			require.Len(t, resp.Evaluations.Evaluations, 1)
			assert.Equal(t, evaluationLogFeatureID, resp.Evaluations.Evaluations[0].FeatureId)
		})
	}
}

func TestGetEvaluationLogsUserIDOnly(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		setup          func(*gatewayService)
		body           getEvaluationRequest
		expectedErr    error
		expectedLogMsg string
	}{
		{
			desc: "err: request check failure is logged with the user id only",
			setup: func(gs *gatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), true)
			},
			body: getEvaluationRequest{
				Tag: "test", User: newEvaluationLogUser(), FeatureID: evaluationLogFeatureID,
			},
			expectedErr:    errDisabledAPIKey,
			expectedLogMsg: "Failed to check GetEvaluation request",
		},
		{
			desc: "err: validation failure is logged with the user id only",
			setup: func(gs *gatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
			},
			body:           getEvaluationRequest{Tag: "test", User: newEvaluationLogUser()},
			expectedErr:    errFeatureIDRequired,
			expectedLogMsg: "Failed to validate GetEvaluation request",
		},
		{
			desc: "success: the evaluation is returned",
			setup: func(gs *gatewayService) {
				expectSDKClientAPIKey(gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache), false)
				expectSingleFixedFeature(gs.featuresCache.(*cachev3mock.MockFeaturesCache))
			},
			body: getEvaluationRequest{
				Tag: "test", User: newEvaluationLogUser(), FeatureID: evaluationLogFeatureID,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			gs := newGatewayServiceWithMock(t, mockController)
			logger, logs := newObservedLogger()
			gs.logger = logger
			p.setup(gs)
			req := httptest.NewRequest(http.MethodPost, dummyURL, renderBody(t, p.body))
			req.Header.Add(authorizationKey, "test-key")
			rec := httptest.NewRecorder()
			gs.getEvaluation(rec, req)
			assertEntriesOmit(t, logs, evaluationLogUserAttribute)
			if p.expectedErr != nil {
				assert.NotEqual(t, http.StatusOK, rec.Code)
				assert.Equal(t, newErrResponse(t, p.expectedErr), rec.Body.String())
				requireUserIDEntry(t, logs, p.expectedLogMsg, "userId", evaluationLogUserID)
				return
			}
			require.Equal(t, http.StatusOK, rec.Code)
			var resp getEvaluationResponse
			require.NoError(t, json.Unmarshal(decodeSuccessResponse(t, rec.Body), &resp))
			require.NotNil(t, resp.Evaluation)
			assert.Equal(t, evaluationLogFeatureID, resp.Evaluation.FeatureId)
			assert.Equal(t, "variation-b", resp.Evaluation.VariationId)
		})
	}
}

func TestEventValidatorsLogUserIDOnly(t *testing.T) {
	t.Parallel()
	const validEventID = "0efe416e-2fd2-4996-b5c3-194f05444f1f"
	newEvent := func(t *testing.T, id string, msg proto.Message) *eventproto.Event {
		t.Helper()
		any, err := anypb.New(msg)
		require.NoError(t, err)
		return &eventproto.Event{Id: id, EnvironmentId: "ns0", Event: any}
	}
	goalEvent := func() *eventproto.GoalEvent {
		return &eventproto.GoalEvent{
			Timestamp: time.Now().Unix(),
			GoalId:    "goal-1",
			UserId:    evaluationLogUserID,
			User:      newEvaluationLogUser(),
		}
	}
	evaluationEvent := func() *eventproto.EvaluationEvent {
		return &eventproto.EvaluationEvent{
			Timestamp:   time.Now().Unix(),
			FeatureId:   evaluationLogFeatureID,
			VariationId: "variation-b",
			UserId:      evaluationLogUserID,
			User:        newEvaluationLogUser(),
			Reason:      &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
		}
	}
	patterns := []struct {
		desc           string
		newValidator   func(*testing.T, *zap.Logger) eventValidator
		expectedCode   string
		expectedErr    error
		expectedLogMsg string
	}{
		{
			desc: "err: invalid goal event id is logged with the user id only",
			newValidator: func(t *testing.T, logger *zap.Logger) eventValidator {
				return newEventValidator(
					newEvent(t, "not-a-uuid", goalEvent()),
					oldestTimestampDuration, furthestTimestampDuration, logger,
				)
			},
			expectedCode:   codeInvalidID,
			expectedErr:    errInvalidIDFormat,
			expectedLogMsg: "Failed to validate goal event id format",
		},
		{
			desc: "err: invalid evaluation event id is logged with the user id only",
			newValidator: func(t *testing.T, logger *zap.Logger) eventValidator {
				return newEventValidator(
					newEvent(t, "not-a-uuid", evaluationEvent()),
					oldestTimestampDuration, furthestTimestampDuration, logger,
				)
			},
			expectedCode:   codeInvalidID,
			expectedErr:    errInvalidIDFormat,
			expectedLogMsg: "Failed to validate evaluation event id format",
		},
		{
			desc: "success: valid goal event",
			newValidator: func(t *testing.T, logger *zap.Logger) eventValidator {
				return newEventValidator(
					newEvent(t, validEventID, goalEvent()),
					oldestTimestampDuration, furthestTimestampDuration, logger,
				)
			},
		},
		{
			desc: "success: valid evaluation event",
			newValidator: func(t *testing.T, logger *zap.Logger) eventValidator {
				return newEventValidator(
					newEvent(t, validEventID, evaluationEvent()),
					oldestTimestampDuration, furthestTimestampDuration, logger,
				)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			logger, logs := newObservedLogger()
			code, err := p.newValidator(t, logger).validate(context.Background())
			assertEntriesOmit(t, logs, evaluationLogUserAttribute)
			assert.Equal(t, p.expectedCode, code)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr != nil {
				requireUserIDEntry(t, logs, p.expectedLogMsg, "user_id", evaluationLogUserID)
			}
		})
	}
}
