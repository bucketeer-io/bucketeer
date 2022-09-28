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

package transformer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	pullermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/mock"
	ucmock "github.com/bucketeer-io/bucketeer/pkg/user/client/mock"
	clienteventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestNewTransformer(t *testing.T) {
	t.Parallel()
	tf := NewTransformer(nil, nil, nil)
	assert.IsType(t, &transformer{}, tf)
}

func TestCheck(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup    func(tf *transformer)
		expected health.Status
	}{
		{
			setup:    func(tf *transformer) { tf.cancel() },
			expected: health.Unhealthy,
		},
		{
			setup: func(tf *transformer) {
				tf.errgroup.Go(func() error { return nil })
				time.Sleep(100 * time.Millisecond) // wait for p.group.FinishedCount() is incremented
			},
			expected: health.Unhealthy,
		},
		{
			setup:    nil,
			expected: health.Healthy,
		},
	}

	for _, p := range patterns {
		tf := newTransformer(t, mockController)
		if p.setup != nil {
			p.setup(tf)
		}
		assert.Equal(t, p.expected, tf.Check(context.Background()))
	}
}

func TestHandle(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*transformer)
		input       *clienteventproto.GoalBatchEvent
		expectedErr error
	}{
		"error: transform": {
			setup:       nil,
			input:       &clienteventproto.GoalBatchEvent{UserId: "uid-0"},
			expectedErr: nil,
		},
		"internal error": {
			setup: func(t *transformer) {
				t.userClient.(*ucmock.MockClient).EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("internal error"))
			},
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
				UserGoalEventsOverTags: []*clienteventproto.UserGoalEventsOverTag{
					{
						Tag: "t-0",
						UserGoalEvents: []*clienteventproto.UserGoalEvent{
							{Timestamp: 0, GoalId: "gid-0", Value: 0.0},
						},
					},
				},
			},
			expectedErr: errors.New("internal error"),
		},
		"user not found": {
			setup: func(t *transformer) {
				t.userClient.(*ucmock.MockClient).EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "user: not found"))
			},
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
				UserGoalEventsOverTags: []*clienteventproto.UserGoalEventsOverTag{
					{
						Tag: "t-0",
						UserGoalEvents: []*clienteventproto.UserGoalEvent{
							{Timestamp: 0, GoalId: "gid-0", Value: 0.0},
						},
					},
				},
			},
			expectedErr: nil,
		},
		"success": {
			setup: func(t *transformer) {
				t.userClient.(*ucmock.MockClient).EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(
					&userproto.GetUserResponse{User: &userproto.User{
						Id: "uid-0",
						TaggedData: map[string]*userproto.User_Data{
							"t-0": {Value: map[string]string{"key": "value"}},
						},
					}}, nil)
				t.publisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
				UserGoalEventsOverTags: []*clienteventproto.UserGoalEventsOverTag{
					{
						Tag: "t-0",
						UserGoalEvents: []*clienteventproto.UserGoalEvent{
							{Timestamp: 0, GoalId: "gid-0", Value: 0.0},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			tf := newTransformer(t, mockController)
			if p.setup != nil {
				p.setup(tf)
			}
			err := tf.handle(p.input, "n0")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestTransform(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	goalEventWithData := &clienteventproto.GoalEvent{
		SourceId:    clienteventproto.SourceId_GOAL_BATCH,
		Tag:         "t-0",
		Timestamp:   0,
		GoalId:      "gid-0",
		UserId:      "uid-0",
		Value:       0,
		User:        &userproto.User{Id: "uid-0", Data: map[string]string{"key": "value"}},
		Evaluations: nil,
	}
	goalEventWithDataAny, err := ptypes.MarshalAny(goalEventWithData)
	require.NoError(t, err)

	goalEvent := &clienteventproto.GoalEvent{
		SourceId:    clienteventproto.SourceId_GOAL_BATCH,
		Tag:         "t-0",
		Timestamp:   0,
		GoalId:      "gid-0",
		UserId:      "uid-0",
		Value:       0,
		User:        &userproto.User{Id: "uid-0"},
		Evaluations: nil,
	}
	goalEventAny, err := ptypes.MarshalAny(goalEvent)
	require.NoError(t, err)

	patterns := map[string]struct {
		setup       func(*transformer)
		input       *clienteventproto.GoalBatchEvent
		expected    []*clienteventproto.Event
		expectedErr error
	}{
		"no UserGoalEventsOverTags": {
			setup: nil,
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
			},
			expected:    nil,
			expectedErr: nil,
		},
		"fail: getUser": {
			setup: func(t *transformer) {
				t.userClient.(*ucmock.MockClient).EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("internal error"))
			},
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
				UserGoalEventsOverTags: []*clienteventproto.UserGoalEventsOverTag{
					{
						Tag: "t-0",
						UserGoalEvents: []*clienteventproto.UserGoalEvent{
							{Timestamp: 0, GoalId: "gid-0", Value: 0.0},
						},
					},
				},
			},
			expected:    nil,
			expectedErr: errors.New("internal error"),
		},
		"not found: getUser": {
			setup: func(t *transformer) {
				t.userClient.(*ucmock.MockClient).EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(
					nil, status.Error(codes.NotFound, "user: not found"))
			},
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
				UserGoalEventsOverTags: []*clienteventproto.UserGoalEventsOverTag{
					{
						Tag: "t-0",
						UserGoalEvents: []*clienteventproto.UserGoalEvent{
							{Timestamp: 0, GoalId: "gid-0", Value: 0.0},
						},
					},
				},
			},
			expected:    nil,
			expectedErr: nil,
		},
		"tagged data not found": {
			setup: func(t *transformer) {
				t.userClient.(*ucmock.MockClient).EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(
					&userproto.GetUserResponse{User: &userproto.User{
						Id: "uid-0",
						TaggedData: map[string]*userproto.User_Data{
							"t-1": {Value: map[string]string{"key": "value"}},
						},
					}}, nil)
			},
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
				UserGoalEventsOverTags: []*clienteventproto.UserGoalEventsOverTag{
					{
						Tag: "t-0",
						UserGoalEvents: []*clienteventproto.UserGoalEvent{
							{Timestamp: 0, GoalId: "gid-0", Value: 0.0},
						},
					},
				},
			},
			expected: []*clienteventproto.Event{
				{
					EnvironmentNamespace: "n0",
					Event:                goalEventAny,
				},
			},
			expectedErr: nil,
		},
		"success": {
			setup: func(t *transformer) {
				t.userClient.(*ucmock.MockClient).EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(
					&userproto.GetUserResponse{User: &userproto.User{
						Id: "uid-0",
						TaggedData: map[string]*userproto.User_Data{
							"t-0": {Value: map[string]string{"key": "value"}},
						},
					}}, nil)
			},
			input: &clienteventproto.GoalBatchEvent{
				UserId: "uid-0",
				UserGoalEventsOverTags: []*clienteventproto.UserGoalEventsOverTag{
					{
						Tag: "t-0",
						UserGoalEvents: []*clienteventproto.UserGoalEvent{
							{Timestamp: 0, GoalId: "gid-0", Value: 0.0},
						},
					},
				},
			},
			expected: []*clienteventproto.Event{
				{
					EnvironmentNamespace: "n0",
					Event:                goalEventWithDataAny,
				},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			tf := newTransformer(t, mockController)
			if p.setup != nil {
				p.setup(tf)
			}
			actual, err := tf.transform(p.input, "n0")
			if p.expected != nil || actual != nil {
				for i := range p.expected {
					assert.Equal(t, p.expected[i].EnvironmentNamespace, actual[i].EnvironmentNamespace)
					expectedGoalEvent, err := unmarshalGoalEvent(p.expected[0].Event)
					assert.NoError(t, err)
					actualGoalEvent, err := unmarshalGoalEvent(actual[0].Event)
					assert.NoError(t, err)
					assert.Equal(t, expectedGoalEvent, actualGoalEvent)
				}
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func unmarshalGoalEvent(event *any.Any) (*clienteventproto.GoalEvent, error) {
	goalEvent := &clienteventproto.GoalEvent{}
	if err := ptypes.UnmarshalAny(event, goalEvent); err != nil {
		return nil, err
	}
	return goalEvent, nil
}

func newTransformer(t *testing.T, mockController *gomock.Controller) *transformer {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &transformer{
		userClient: ucmock.NewMockClient(mockController),
		puller:     pullermock.NewMockRateLimitedPuller(mockController),
		publisher:  publishermock.NewMockPublisher(mockController),
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		doneCh:     make(chan struct{}),
	}
}
