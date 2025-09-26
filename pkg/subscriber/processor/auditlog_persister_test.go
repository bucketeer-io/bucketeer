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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func TestExtractAuditLogs(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	editor := &eventproto.Editor{Email: "test@example.com"}
	event0, err := domainevent.NewEvent(
		editor,
		eventproto.Event_FEATURE,
		"fId-0",
		eventproto.Event_FEATURE_CREATED,
		&eventproto.FeatureCreatedEvent{Id: "fId-0"},
		"ns0",
		"{\"id\": \"curr\"}",
		"{\"id\": \"prev\"}",
	)
	assert.NoError(t, err)
	event1, err := domainevent.NewEvent(
		editor,
		eventproto.Event_FEATURE,
		"fId-1",
		eventproto.Event_FEATURE_CREATED,
		&eventproto.FeatureCreatedEvent{Id: "fId-1"},
		"ns0",
		"{\"id\": \"curr\"}",
		"{\"id\": \"prev\"}",
	)
	assert.NoError(t, err)
	adninEvent0, err := domainevent.NewAdminEvent(
		editor,
		eventproto.Event_FEATURE,
		"fId-2",
		eventproto.Event_FEATURE_CREATED,
		&eventproto.FeatureCreatedEvent{Id: "fId-2"},
		"{\"id\": \"curr\"}",
		"{\"id\": \"prev\"}",
	)
	assert.NoError(t, err)
	chunk := createChunk(t, []*domain.Event{event0, event1, adninEvent0})

	p := newPersister(t, mockController)
	auditLogs, adminAuditLogs, messages, adminMessages := p.extractAuditLogs(chunk)
	for i, al := range auditLogs {
		msg, ok := chunk[al.Id]
		assert.True(t, ok)
		assert.Equal(t, msg.ID, al.Id)
		assert.Equal(t, messages[i].ID, al.Id)
	}
	for i, al := range adminAuditLogs {
		msg, ok := chunk[al.Id]
		assert.True(t, ok)
		assert.Equal(t, msg.ID, al.Id)
		assert.Equal(t, adminMessages[i].ID, al.Id)
	}
}

func newPersister(t *testing.T, mockController *gomock.Controller) *auditLogPersister {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &auditLogPersister{
		logger: logger.Named("persister"),
	}
}

func createChunk(t *testing.T, events []*domain.Event) map[string]*puller.Message {
	t.Helper()
	chunk := make(map[string]*puller.Message)
	for _, e := range events {
		data, err := proto.Marshal(e)
		require.NoError(t, err)
		chunk[e.Id] = &puller.Message{
			ID:         e.Id,
			Data:       data,
			Attributes: map[string]string{"id": e.Id},
			Ack:        func() {},
			Nack:       func() {},
		}
	}
	return chunk
}
