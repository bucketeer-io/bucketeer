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
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	"github.com/bucketeer-io/bucketeer/proto/feature"
	"github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	oldestTimestampDuration   = 24 * time.Hour
	furthestTimestampDuration = 24 * time.Hour
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
					TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
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
					TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
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
					TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.MetricsEvent",
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
			desc:      "fail: invalid past time",
			timestamp: time.Now().AddDate(0, 0, -2).Unix(),
			expected:  false,
		},
		{
			desc:      "fail: invalid future time",
			timestamp: time.Now().AddDate(0, 0, 2).Unix(),
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
	timestamp := time.Now().Unix()
	patterns := []struct {
		desc          string
		inputFunc     func() *eventproto.Event
		expectedEvent *eventproto.GoalEvent
		expected      string
		expectedErr   error
	}{
		{
			desc: "invalid uuid",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
				}
			},
			expectedEvent: nil,
			expected:      codeInvalidID,
			expectedErr:   errInvalidIDFormat,
		},
		{
			desc: "unmarshal fails",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expectedEvent: nil,
			expected:      codeUnmarshalFailed,
			expectedErr:   errUnmarshalFailed,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeEmptyField,
			expectedErr:   errEmptyGoalID,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeEmptyField,
			expectedErr:   errEmptyUserID,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeInvalidTimestamp,
			expectedErr:   errInvalidTimestamp,
		},
		{
			desc: "success",
			inputFunc: func() *eventproto.Event {
				goalEvent := &eventproto.GoalEvent{
					Timestamp: timestamp,
					GoalId:    "goal-id",
					User: &user.User{
						Id: "user-id",
					},
					Evaluations: []*feature.Evaluation{
						{
							Id: "evaluation-id",
						},
					},
				}
				bGoalEvent, err := proto.Marshal(goalEvent)
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
			expectedEvent: &eventproto.GoalEvent{
				Timestamp: timestamp,
				GoalId:    "goal-id",
				User: &user.User{
					Id: "user-id",
				},
				Evaluations: []*feature.Evaluation{
					{
						Id: "evaluation-id",
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			logger, _ := log.NewLogger()
			v := &eventGoalValidator{
				event:                     p.inputFunc(),
				logger:                    logger,
				oldestTimestampDuration:   24 * time.Hour,
				furthestTimestampDuration: 24 * time.Hour,
			}
			ev, actual, err := v.validate(context.Background())
			if p.expectedEvent == nil {
				assert.Nil(t, ev)
			} else {
				goalEv, ok := ev.(*eventproto.GoalEvent)
				assert.True(t, ok, "Failed to type assert to GoalEvent")
				assert.True(t, proto.Equal(p.expectedEvent, goalEv))
			}
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGrpcValidateEvaluationEvent(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	patterns := []struct {
		desc          string
		inputFunc     func() *eventproto.Event
		expectedEvent *eventproto.EvaluationEvent
		expected      string
		expectedErr   error
	}{
		{
			desc: "invalid uuid",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		{
			desc: "unmarshal fails",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expectedEvent: nil,
			expected:      codeUnmarshalFailed,
			expectedErr:   errUnmarshalFailed,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeEmptyField,
			expectedErr:   errEmptyFeatureID,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeEmptyField,
			expectedErr:   errEmptyVariationID,
		},
		{
			desc: "empty variation_id with ERROR_NO_EVALUATIONS reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   timestamp,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: &eventproto.EvaluationEvent{
				Timestamp:   timestamp,
				FeatureId:   "feature-id",
				VariationId: "",
				User: &user.User{
					Id: "user-id",
				},
				Reason: &feature.Reason{
					Type: feature.Reason_ERROR_NO_EVALUATIONS,
				},
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_FLAG_NOT_FOUND reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   timestamp,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: &eventproto.EvaluationEvent{
				Timestamp:   timestamp,
				FeatureId:   "feature-id",
				VariationId: "",
				User: &user.User{
					Id: "user-id",
				},
				Reason: &feature.Reason{
					Type: feature.Reason_ERROR_FLAG_NOT_FOUND,
				},
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_WRONG_TYPE reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   timestamp,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: &eventproto.EvaluationEvent{
				Timestamp:   timestamp,
				FeatureId:   "feature-id",
				VariationId: "",
				User: &user.User{
					Id: "user-id",
				},
				Reason: &feature.Reason{
					Type: feature.Reason_ERROR_WRONG_TYPE,
				},
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_USER_ID_NOT_SPECIFIED reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   timestamp,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: &eventproto.EvaluationEvent{
				Timestamp:   timestamp,
				FeatureId:   "feature-id",
				VariationId: "",
				User: &user.User{
					Id: "user-id",
				},
				Reason: &feature.Reason{
					Type: feature.Reason_ERROR_USER_ID_NOT_SPECIFIED,
				},
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   timestamp,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: &eventproto.EvaluationEvent{
				Timestamp:   timestamp,
				FeatureId:   "feature-id",
				VariationId: "",
				User: &user.User{
					Id: "user-id",
				},
				Reason: &feature.Reason{
					Type: feature.Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED,
				},
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with ERROR_EXCEPTION reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   timestamp,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: &eventproto.EvaluationEvent{
				Timestamp:   timestamp,
				FeatureId:   "feature-id",
				VariationId: "",
				User: &user.User{
					Id: "user-id",
				},
				Reason: &feature.Reason{
					Type: feature.Reason_ERROR_EXCEPTION,
				},
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			desc: "empty variation_id with CLIENT reason",
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp:   timestamp,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: &eventproto.EvaluationEvent{
				Timestamp:   timestamp,
				FeatureId:   "feature-id",
				VariationId: "",
				User: &user.User{
					Id: "user-id",
				},
				Reason: &feature.Reason{
					Type: feature.Reason_CLIENT,
				},
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeEmptyField,
			expectedErr:   errEmptyUserID,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeEmptyField,
			expectedErr:   errEmptyUserID,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeEmptyField,
			expectedErr:   errNilReason,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: nil,
			expected:      codeInvalidTimestamp,
			expectedErr:   errInvalidTimestamp,
		},
		{
			desc: "success",
			inputFunc: func() *eventproto.Event {
				evalEvent := &eventproto.EvaluationEvent{
					Timestamp:      timestamp,
					FeatureId:      "feature-id",
					FeatureVersion: 1,
					VariationId:    "variation-id",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				}
				bEvaluationEvent, err := proto.Marshal(evalEvent)
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
						Value:   bEvaluationEvent,
					},
				}
			},
			expectedEvent: func() *eventproto.EvaluationEvent {
				evalEvent := &eventproto.EvaluationEvent{
					Timestamp:      timestamp,
					FeatureId:      "feature-id",
					FeatureVersion: 1,
					VariationId:    "variation-id",
					User: &user.User{
						Id: "user-id",
					},
					Reason: &feature.Reason{
						Type: feature.Reason_DEFAULT,
					},
				}
				return evalEvent
			}(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			logger, _ := log.NewLogger()
			v := &eventEvaluationValidator{
				event:                     p.inputFunc(),
				logger:                    logger,
				oldestTimestampDuration:   24 * time.Hour,
				furthestTimestampDuration: 24 * time.Hour,
			}
			ev, actual, err := v.validate(context.Background())
			if p.expectedEvent == nil {
				assert.Nil(t, ev)
			} else {
				evalEv, ok := ev.(*eventproto.EvaluationEvent)
				assert.True(t, ok, "Failed to type assert to EvaluationEvent")
				assert.True(t, proto.Equal(p.expectedEvent, evalEv))
			}
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGrpcValidateMetrics(t *testing.T) {
	t.Parallel()
	timestamp := time.Now().Unix()
	patterns := []struct {
		desc          string
		inputFunc     func() *eventproto.Event
		expectedEvent *eventproto.MetricsEvent
		expected      string
		expectedErr   error
	}{
		{
			desc: "invalid uuid",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
				}
			},
			expectedEvent: nil,
			expected:      codeInvalidID,
			expectedErr:   errInvalidIDFormat,
		},
		{
			desc: "unmarshal fails",
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expectedEvent: nil,
			expected:      codeUnmarshalFailed,
			expectedErr:   errUnmarshalFailed,
		},
		{
			desc: "success",
			inputFunc: func() *eventproto.Event {
				metricsEvent := &eventproto.MetricsEvent{
					Timestamp: timestamp,
				}
				b, err := proto.Marshal(metricsEvent)
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.MetricsEvent",
						Value:   b,
					},
				}
			},
			expectedEvent: func() *eventproto.MetricsEvent {
				metricsEvent := &eventproto.MetricsEvent{
					Timestamp: timestamp,
				}
				return metricsEvent
			}(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			logger, _ := log.NewLogger()
			v := &eventMetricsValidator{
				event:                     p.inputFunc(),
				oldestTimestampDuration:   oldestTimestampDuration,
				furthestTimestampDuration: furthestTimestampDuration,
				logger:                    logger,
			}
			ev, actual, err := v.validate(context.Background())
			if p.expectedEvent == nil {
				assert.Nil(t, ev)
			} else {
				metricsEv, ok := ev.(*eventproto.MetricsEvent)
				assert.True(t, ok, "Failed to type assert to MetricsEvent")
				assert.True(t, proto.Equal(p.expectedEvent, metricsEv))
			}
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
