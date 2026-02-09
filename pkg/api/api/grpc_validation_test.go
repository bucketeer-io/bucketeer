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
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	"github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	oldestTimestampDuration   = 744 * time.Hour // 31 days - aligns with 30-day DB retention + 1 day buffer
	furthestTimestampDuration = 1 * time.Hour   // 1 hour - handles legitimate clock skew while preventing malicious timestamps
)

func TestNewEventValidator(t *testing.T) {
	t.Parallel()
	bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{})
	if err != nil {
		t.Fatal("could not serialize evaluation event")
	}
	bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{})
	if err != nil {
		t.Fatal("could not serialize goal event")
	}
	bMetricsEvent, err := proto.Marshal(&eventproto.MetricsEvent{})
	if err != nil {
		t.Fatal("could not serialize metrics event")
	}
	patterns := []struct {
		desc     string
		input    *eventproto.Event
		expected eventValidator
	}{
		{
			desc: "evaluationValidator",
			input: &eventproto.Event{
				Id: newUUID(t),
				Event: &any.Any{
					TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
					Value:   bEvaluationEvent,
				},
			},
			expected: &eventEvaluationValidator{},
		},
		{
			desc: "GoalValidator",
			input: &eventproto.Event{
				Id: newUUID(t),
				Event: &any.Any{
					TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.GoalEvent",
					Value:   bGoalEvent,
				},
			},
			expected: &eventGoalValidator{},
		},
		{
			desc: "MetricsEvent",
			input: &eventproto.Event{
				Id: newUUID(t),
				Event: &any.Any{
					TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.MetricsEvent",
					Value:   bMetricsEvent,
				},
			},
			expected: &eventMetricsValidator{},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			logger, _ := log.NewLogger()
			actual := newEventValidator(p.input, oldestTimestampDuration, furthestTimestampDuration, logger)
			assert.IsType(t, p.expected, actual)
		})
	}
}

func TestValidateTimestamp(t *testing.T) {
	testcases := []struct {
		desc      string
		timestamp int64
		expected  bool
	}{
		{
			desc:      "success",
			timestamp: time.Now().Unix(),
			expected:  true,
		},
		{
			desc:      "success: within 30-day retention window",
			timestamp: time.Now().Add(-30 * 24 * time.Hour).Unix(),
			expected:  true,
		},
		{
			desc:      "fail: invalid past time - older than 30-day retention",
			timestamp: time.Now().Add(-31 * 24 * time.Hour).Unix(),
			expected:  false,
		},
		{
			desc:      "fail: invalid future time",
			timestamp: time.Now().Add(2 * time.Hour).Unix(),
			expected:  false,
		},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("index %d", i)
		res := validateTimestamp(tc.timestamp, oldestTimestampDuration, furthestTimestampDuration)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestGrpcValidateGoalEvent(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
	}{
		{
			desc: "unmarshal fails",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		{
			desc: "invalid uuid",
			inputFunc: func() *eventproto.Event {
				bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
					Timestamp: time.Now().Unix(),
					GoalId:    "goal-id",
					User: &user.User{
						Id: "user-id",
					},
					Evaluations: []*feature.Evaluation{
						{
							Id: "evaluation-id",
						},
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		{
			desc: "empty goal_id",
			inputFunc: func() *eventproto.Event {
				bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
					Timestamp: time.Now().Unix(),
					GoalId:    "",
					User: &user.User{
						Id: "user-id",
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expected:    codeEmptyField,
			expectedErr: errEmptyGoalID,
		},
		{
			desc: "empty user_id",
			inputFunc: func() *eventproto.Event {
				bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
					Timestamp: time.Now().Unix(),
					GoalId:    "goal-id",
					User: &user.User{
						Id: "",
					},
					UserId: "",
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expected:    codeEmptyField,
			expectedErr: errEmptyUserID,
		},
		{
			desc: "invalid timestamp",
			inputFunc: func() *eventproto.Event {
				bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
					Timestamp: int64(999999999999999),
					GoalId:    "goal-id",
					User: &user.User{
						Id: "user-id",
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expected:    codeInvalidTimestamp,
			expectedErr: errInvalidTimestamp,
		},
		{
			desc: "success",
			inputFunc: func() *eventproto.Event {
				bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
					Timestamp: time.Now().Unix(),
					GoalId:    "goal-id",
					User: &user.User{
						Id: "user-id",
					},
					Evaluations: []*feature.Evaluation{
						{
							Id: "evaluation-id",
						},
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			logger, _ := log.NewLogger()
			v := &eventGoalValidator{
				event:                     p.inputFunc(),
				logger:                    logger,
				oldestTimestampDuration:   oldestTimestampDuration,
				furthestTimestampDuration: furthestTimestampDuration,
			}
			actual, err := v.validate(context.Background())
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGrpcValidateEvaluationEvent(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
	}{
		{
			desc: "unmarshal fails",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		{
			desc: "invalid uuid",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:      time.Now().Unix(),
					FeatureId:      "feature-id",
					FeatureVersion: 1,
					VariationId:    "variation-id",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		{
			desc: "empty feature_id",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "",
					VariationId: "variation-id",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    codeEmptyField,
			expectedErr: errEmptyFeatureID,
		},
		{
			desc: "empty variation_id",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    codeEmptyField,
			expectedErr: errEmptyVariationID,
		},
		{
			desc: "empty variation_id with ERROR_NO_EVALUATIONS reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_ERROR_NO_EVALUATIONS,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_FLAG_NOT_FOUND reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_ERROR_FLAG_NOT_FOUND,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_WRONG_TYPE reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_ERROR_WRONG_TYPE,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_USER_ID_NOT_SPECIFIED reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_ERROR_USER_ID_NOT_SPECIFIED,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_EXCEPTION reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_ERROR_EXCEPTION,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_CACHE_NOT_FOUND reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_ERROR_CACHE_NOT_FOUND,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with CLIENT reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_CLIENT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty user_id",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "variation-id",
					User: &user.User{
						Id: "",
					},
					UserId: "",
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    codeEmptyField,
			expectedErr: errEmptyUserID,
		},
		{
			desc: "nil user",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "variation-id",
					User:        nil,
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    codeEmptyField,
			expectedErr: errEmptyUserID,
		},
		{
			desc: "nil reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   time.Now().Unix(),
					FeatureId:   "feature-id",
					VariationId: "variation-id",
					User: &user.User{
						Id: "user-id",
					},
					Reason: nil,
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    codeEmptyField,
			expectedErr: errNilReason,
		},
		{
			desc: "invalid timestamp",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   int64(999999999999999),
					FeatureId:   "feature-id",
					VariationId: "variation-id",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expected:    codeInvalidTimestamp,
			expectedErr: errInvalidTimestamp,
		},
		{
			desc: "success",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:      time.Now().Unix(),
					FeatureId:      "feature-id",
					FeatureVersion: 1,
					VariationId:    "variation-id",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			logger, _ := log.NewLogger()
			v := &eventEvaluationValidator{
				event:                     p.inputFunc(),
				logger:                    logger,
				oldestTimestampDuration:   oldestTimestampDuration,
				furthestTimestampDuration: furthestTimestampDuration,
			}
			actual, err := v.validate(context.Background())
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

// TestGrpcValidateEvaluationEventAllowsEmptyVariationIdForAllProtoErrorTypes fails when a new error
// reason type is added to the proto (CLIENT or ERROR_* naming) but the grpc_validation.go's isErrorReason
// logic hasn't been updated. This forces us to update the validation when adding new error types.
func TestGrpcValidateEvaluationEventAllowsEmptyVariationIdForAllProtoErrorTypes(t *testing.T) {
	t.Parallel()
	for value, name := range feature.Reason_Type_name {
		isErrorType := name == "CLIENT" || strings.HasPrefix(name, "ERROR_")
		if !isErrorType {
			continue
		}
		reasonType := feature.Reason_Type(value)
		bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
			Timestamp:   time.Now().Unix(),
			FeatureId:   "feature-id",
			VariationId: "",
			User: &user.User{
				Id: "user-id",
			},
			Reason: &feature.Reason{
				Type: reasonType,
			},
		})
		if err != nil {
			t.Fatalf("could not serialize evaluation event: %v", err)
		}
		event := &eventproto.Event{
			Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
			Event: &any.Any{
				TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
				Value:   bEvaluationEvent,
			},
		}
		logger, _ := log.NewLogger()
		v := &eventEvaluationValidator{
			event:                     event,
			logger:                    logger,
			oldestTimestampDuration:   oldestTimestampDuration,
			furthestTimestampDuration: furthestTimestampDuration,
		}
		code, err := v.validate(context.Background())
		assert.Empty(t, code, "Reason %s (value=%d): validation should pass for error types with empty variation_id. "+
			"Update isErrorReason in grpc_validation.go to include this type.", name, value)
		assert.NoError(t, err)
	}
}

func TestIsEvaluationEventErrorReason(t *testing.T) {
	t.Parallel()
	assert.False(t, isEvaluationEventErrorReason(nil), "nil reason (e.g. from old SDKs) must return false")
	assert.False(t, isEvaluationEventErrorReason(&feature.Reason{Type: feature.Reason_DEFAULT}))
	assert.False(t, isEvaluationEventErrorReason(&feature.Reason{Type: feature.Reason_RULE}))
	assert.True(t, isEvaluationEventErrorReason(&feature.Reason{Type: feature.Reason_CLIENT}))
	assert.True(t, isEvaluationEventErrorReason(&feature.Reason{Type: feature.Reason_ERROR_FLAG_NOT_FOUND}))
	assert.True(t, isEvaluationEventErrorReason(&feature.Reason{Type: feature.Reason_ERROR_CACHE_NOT_FOUND}))
}

func TestEventEvaluationValidatorStoresLastUnmarshaledEvent(t *testing.T) {
	t.Parallel()
	bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
		Timestamp:   time.Now().Unix(),
		FeatureId:   "feature-id",
		VariationId: "",
		User: &user.User{
			Id: "user-id",
		},
		Reason: &feature.Reason{
			Type: feature.Reason_ERROR_FLAG_NOT_FOUND,
		},
		Tag:        "tag1",
		SdkVersion: "1.0.0",
		SourceId:   eventproto.SourceId_GO_SERVER,
	})
	if err != nil {
		t.Fatalf("could not serialize evaluation event: %v", err)
	}
	event := &eventproto.Event{
		Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
		Event: &any.Any{
			TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.EvaluationEvent",
			Value:   bEvaluationEvent,
		},
	}
	logger, _ := log.NewLogger()
	v := &eventEvaluationValidator{
		event:                     event,
		logger:                    logger,
		oldestTimestampDuration:   oldestTimestampDuration,
		furthestTimestampDuration: furthestTimestampDuration,
	}
	code, err := v.validate(context.Background())
	assert.Empty(t, code)
	assert.NoError(t, err)
	assert.NotNil(t, v.lastUnmarshaledEvent, "lastUnmarshaledEvent should be set after successful validation")
	assert.Equal(t, "feature-id", v.lastUnmarshaledEvent.FeatureId)
	assert.Equal(t, feature.Reason_ERROR_FLAG_NOT_FOUND, v.lastUnmarshaledEvent.Reason.Type)
	assert.Equal(t, "tag1", v.lastUnmarshaledEvent.Tag)
	assert.Equal(t, "1.0.0", v.lastUnmarshaledEvent.SdkVersion)
	assert.Equal(t, eventproto.SourceId_GO_SERVER, v.lastUnmarshaledEvent.SourceId)
}

func TestGrpcValidateMetrics(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
	}{
		{
			desc: "unmarshal fails",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		{
			desc: "invalid uuid",
			inputFunc: func() *eventproto.Event {
				b, err := proto.Marshal(&eventproto.MetricsEvent{
					Timestamp: time.Now().Unix(),
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.MetricsEvent",
						Value:   b,
					},
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		{
			desc: "invalid timestamp",
			inputFunc: func() *eventproto.Event {
				b, err := proto.Marshal(&eventproto.MetricsEvent{
					Timestamp: int64(999999999999999),
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.MetricsEvent",
						Value:   b,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "success",
			inputFunc: func() *eventproto.Event {
				b, err := proto.Marshal(&eventproto.MetricsEvent{
					Timestamp: time.Now().Unix(),
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/v2/proto/event/client/bucketeer.event.client.MetricsEvent",
						Value:   b,
					},
				}
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			logger, _ := log.NewLogger()
			v := &eventMetricsValidator{
				event:                     p.inputFunc(),
				oldestTimestampDuration:   oldestTimestampDuration,
				furthestTimestampDuration: furthestTimestampDuration,
				logger:                    logger,
			}
			actual, err := v.validate(context.Background())
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
