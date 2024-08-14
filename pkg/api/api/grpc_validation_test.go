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
	patterns := []struct {
		desc        string
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
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
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		{
			desc: "invalid timestamp",
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
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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

func TestGrpcValidateEvaluationEvent(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
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
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
		},
		{
			desc: "invalid timestamp",
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
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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

func TestGrpcValidateMetrics(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		inputFunc   func() *eventproto.Event
		expected    string
		expectedErr error
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
			expected:    codeUnmarshalFailed,
			expectedErr: errUnmarshalFailed,
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.MetricsEvent",
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
						TypeUrl: "github.com/bucketeer-io/bucketeer/proto/event/client/bucketeer.event.client.MetricsEvent",
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
