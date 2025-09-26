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

package processor

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	ecproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/service"
)

func TestValidateEvent(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		input    *eventproto.UserEvent
		expected bool
	}{
		{
			input: &eventproto.UserEvent{
				UserId:   "hoge",
				LastSeen: 3456789,
			},
			expected: true,
		},
		{
			input:    &eventproto.UserEvent{},
			expected: false,
		},
		{
			input: &eventproto.UserEvent{
				UserId:   "",
				LastSeen: 3456789,
			},
			expected: false,
		},
		{
			input: &eventproto.UserEvent{
				UserId:   "hoge",
				LastSeen: 0,
			},
			expected: false,
		},
	}
	logger, _ := log.NewLogger()
	pst := userEventPersister{logger: logger}
	for _, p := range patterns {
		actual := pst.validateEvent(p.input)
		assert.Equal(t, p.expected, actual)
	}
}

func TestUpsert(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	patterns := []struct {
		desc, environmentId string
		setup               func(persister *userEventPersister)
		input               []*eventproto.UserEvent
		expected            error
	}{
		{
			desc:          "upsert mau error",
			environmentId: "env1",
			setup: func(p *userEventPersister) {
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("internal"))
			},
			input: []*eventproto.UserEvent{
				{
					EnvironmentId: "env1",
					UserId:        "id-1",
					LastSeen:      3,
				},
			},
			expected: errors.New("internal"),
		},
		{
			desc:          "upsert success",
			environmentId: "env1",
			setup: func(p *userEventPersister) {
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: []*eventproto.UserEvent{
				{
					EnvironmentId: "env1",
					UserId:        "id-1",
					LastSeen:      3,
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pst := newPersisterWithMock(t, mockController, now, uuid)
			if p.setup != nil {
				p.setup(pst)
			}
			err := pst.upsertMAUs(ctx, p.input, p.environmentId)
			assert.Equal(t, p.expected, err)
		})
	}
}

func newPersisterWithMock(
	t *testing.T,
	mockController *gomock.Controller,
	now time.Time,
	id *uuid.UUID,
) *userEventPersister {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &userEventPersister{
		mysqlClient: mysqlmock.NewMockClient(mockController),
		timeNow:     func() time.Time { return now },
		newUUID:     func() (*uuid.UUID, error) { return id, nil },
		logger:      logger,
	}
}

func TestHandleChunk(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	now := time.Now()
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)
	aMap := generatePullerMessages(t, 2, "env-1")
	bMap := generatePullerMessages(t, 3, "env-2")

	patterns := []struct {
		desc, environmentId string
		setup               func(persister *userEventPersister)
		input               map[string]*puller.Message
		expected            error
	}{
		{
			desc:          "upsert mau error",
			environmentId: "env1",
			setup: func(p *userEventPersister) {
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("internal")).Times(2)
			},
			input:    mergeMap(t, aMap, bMap),
			expected: errors.New("internal"),
		},
		{
			desc:          "upsert success",
			environmentId: "env1",
			setup: func(p *userEventPersister) {
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil).Times(2)
			},
			input:    mergeMap(t, aMap, bMap),
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pst := newPersisterWithMock(t, mockController, now, uuid)
			if p.setup != nil {
				p.setup(pst)
			}
			pst.handleChunk(p.input)
		})
	}
}

func generatePullerMessages(
	t *testing.T,
	size int,
	environmentId string,
) map[string]*puller.Message {
	t.Helper()
	messages := make(map[string]*puller.Message)
	for i := 0; i < size; i++ {
		userEvent := generateUserEvent(t, environmentId)
		msg := generatePullerMessage(t, userEvent)
		messages[msg.ID] = msg
	}
	return messages
}

func generatePullerMessage(t *testing.T, userEvent *eventproto.UserEvent) *puller.Message {
	t.Helper()
	ue, err := ptypes.MarshalAny(userEvent)
	if err != nil {
		t.Fatalf("Failed to marshal any user event: %v", err)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}
	ev := &ecproto.Event{
		Id:            id.String(),
		Event:         ue,
		EnvironmentId: userEvent.EnvironmentId,
	}
	data, err := proto.Marshal(ev)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	msg := &puller.Message{
		ID:   fmt.Sprintf("message-id-%s", id.String()),
		Data: data,
		Nack: func() {},
		Ack:  func() {},
	}
	return msg
}

func generateUserEvent(t *testing.T, environmentId string) *eventproto.UserEvent {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}
	return &eventproto.UserEvent{
		EnvironmentId: environmentId,
		UserId:        fmt.Sprintf("user-id-%s", id),
		LastSeen:      time.Now().Unix(),
	}
}

func mergeMap(t *testing.T, a, b map[string]*puller.Message) map[string]*puller.Message {
	t.Helper()
	for k, v := range b {
		a[k] = v
	}
	return a
}
