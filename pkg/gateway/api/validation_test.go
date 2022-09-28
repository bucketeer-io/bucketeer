// Copyright 2022 The Bucketeer Authors.
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
	bGoalBatchEvent, err := proto.Marshal(&eventproto.GoalBatchEvent{})
	if err != nil {
		t.Fatal("could not serialize goal batch event")
	}
	bMetricsEvent, err := proto.Marshal(&eventproto.MetricsEvent{})
	if err != nil {
		t.Fatal("could not serialize metrics event")
	}
	patterns := map[string]struct {
		input    *eventproto.Event
		expected eventValidator
	}{
		"evaluationValidator": {
			input: &eventproto.Event{
				Id: newUUID(t),
				Event: &any.Any{
					TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.EvaluationEvent",
					Value:   bEvaluationEvent,
				},
			},
			expected: &eventEvaluationValidator{},
		},
		"GoalValidator": {
			input: &eventproto.Event{
				Id: newUUID(t),
				Event: &any.Any{
					TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
					Value:   bGoalEvent,
				},
			},
			expected: &eventGoalValidator{},
		},
		"GoalBatchValidator": {
			input: &eventproto.Event{
				Id: newUUID(t),
				Event: &any.Any{
					TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalBatchEvent",
					Value:   bGoalBatchEvent,
				},
			},
			expected: &eventGoalBatchValidator{},
		},
		"MetricsEvent": {
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
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			logger, _ := log.NewLogger()
			actual := newEventValidator(p.input, oldestTimestampDuration, furthestTimestampDuration, logger)
			assert.IsType(t, p.expected, actual)
		})
	}
}

func TestValidateTimestamp(t *testing.T) {
	testcases := []struct {
		timestamp int64
		expected  bool
	}{
		{
			timestamp: time.Now().Unix(),
			expected:  true,
		},
		{
			timestamp: time.Now().AddDate(0, 0, -2).Unix(),
			expected:  false,
		},
		{
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

func TestValidateGoalEvent(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
	}{
		"invalid uuid": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		"unmarshal fails": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		"invalid timestamp": {
			inputFunc: func() *eventproto.Event {
				bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
					Timestamp: int64(999999999999999),
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
			expected:    codeInvalidTimestamp,
			expectedErr: errInvalidTimestamp,
		},
		"success": {
			inputFunc: func() *eventproto.Event {
				bGoalEvent, err := proto.Marshal(&eventproto.GoalEvent{
					Timestamp: time.Now().Unix(),
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalEvent",
						Value:   bGoalEvent,
					},
				}
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			logger, _ := log.NewLogger()
			v := &eventGoalValidator{
				event:                     p.inputFunc(),
				logger:                    logger,
				oldestTimestampDuration:   24 * time.Hour,
				furthestTimestampDuration: 24 * time.Hour,
			}
			actual, err := v.validate(context.Background())
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateGoalBatchEvent(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
	}{
		"err: invalid uuid": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		"err: unmarshal failed": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		"err: empty user id": {
			inputFunc: func() *eventproto.Event {
				bGoalBatchEvent, err := proto.Marshal(&eventproto.GoalBatchEvent{
					UserId: "",
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalBatchEvent",
						Value:   bGoalBatchEvent,
					},
				}
			},
			expected:    codeEmptyUserID,
			expectedErr: errEmptyUserID,
		},
		"err: empty tag": {
			inputFunc: func() *eventproto.Event {
				bGoalBatchEvent, err := proto.Marshal(&eventproto.GoalBatchEvent{
					UserId: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					UserGoalEventsOverTags: []*eventproto.UserGoalEventsOverTag{
						{
							Tag: "",
						},
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalBatchEvent",
						Value:   bGoalBatchEvent,
					},
				}
			},
			expected:    codeEmptyTag,
			expectedErr: errEmptyTag,
		},
		"success": {
			inputFunc: func() *eventproto.Event {
				bGoalBatchEvent, err := proto.Marshal(&eventproto.GoalBatchEvent{
					UserId: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					UserGoalEventsOverTags: []*eventproto.UserGoalEventsOverTag{
						{
							Tag: "tag",
						},
					},
				})
				if err != nil {
					t.Fatal("could not serialize event")
				}
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
					Event: &any.Any{
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.GoalBatchEvent",
						Value:   bGoalBatchEvent,
					},
				}
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			logger, _ := log.NewLogger()
			v := &eventGoalBatchValidator{
				event:                     p.inputFunc(),
				logger:                    logger,
				oldestTimestampDuration:   24 * time.Hour,
				furthestTimestampDuration: 24 * time.Hour,
			}
			actual, err := v.validate(context.Background())
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateEvaluationEvent(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
	}{
		"invalid uuid": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		"unmarshal fails": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		"invalid timestamp": {
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp: int64(999999999999999),
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
			expected:    codeInvalidTimestamp,
			expectedErr: errInvalidTimestamp,
		},
		"success": {
			inputFunc: func() *eventproto.Event {
				bEvaluationEvent, err := proto.Marshal(&eventproto.EvaluationEvent{
					Timestamp: time.Now().Unix(),
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
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			logger, _ := log.NewLogger()
			v := &eventEvaluationValidator{
				event:                     p.inputFunc(),
				logger:                    logger,
				oldestTimestampDuration:   24 * time.Hour,
				furthestTimestampDuration: 24 * time.Hour,
			}
			actual, err := v.validate(context.Background())
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateMetrics(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
	}{
		"invalid uuid": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e 2fd2 4996 c5c3 194f05444f1f",
				}
			},
			expected:    codeInvalidID,
			expectedErr: errInvalidIDFormat,
		},
		"unmarshal fails": {
			inputFunc: func() *eventproto.Event {
				return &eventproto.Event{
					Id: "0efe416e-2fd2-4996-b5c3-194f05444f1f",
				}
			},
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		"invalid timestamp": {
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.MetricsEvent",
						Value:   b,
					},
				}
			},
			expected:    "",
			expectedErr: nil,
		},
		"success": {
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.MetricsEvent",
						Value:   b,
					},
				}
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
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
