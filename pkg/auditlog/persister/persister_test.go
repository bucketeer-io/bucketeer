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

package persister

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	metricsmock "github.com/bucketeer-io/bucketeer/pkg/metrics/mock"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	pullermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/proto/event/domain"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestNewPersister(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	puller := pullermock.NewMockPuller(mockController)
	mysqlClient := mysqlmock.NewMockClient(mockController)
	registerer := metricsmock.NewMockRegisterer(mockController)
	registerer.EXPECT().MustRegister(gomock.Any()).Return()
	p := NewPersister(
		puller,
		mysqlClient,
		WithMaxMPS(1000),
		WithNumWorkers(1),
		WithFlushSize(100),
		WithFlushInterval(time.Second),
		WithMetrics(registerer),
		WithLogger(zap.NewNop()),
	)
	assert.IsType(t, &Persister{}, p)
}

func TestCheck(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		setup    func(p *Persister)
		expected health.Status
	}{
		{
			setup:    func(p *Persister) { p.cancel() },
			expected: health.Unhealthy,
		},
		{
			setup: func(p *Persister) {
				p.group.Go(func() error { return nil })
				time.Sleep(100 * time.Millisecond) // wait for p.group.FinishedCount() is incremented
			},
			expected: health.Unhealthy,
		},
		{
			setup:    nil,
			expected: health.Healthy,
		},
	}

	for _, pat := range patterns {
		p := newPersister(t, mockController)
		if pat.setup != nil {
			pat.setup(p)
		}
		assert.Equal(t, pat.expected, p.Check(context.Background()))
	}
}

func TestExtractAuditLogs(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	editor := &eventproto.Editor{Email: "test@example.com"}
	event0, err := domainevent.NewEvent(editor, eventproto.Event_FEATURE, "fId-0", eventproto.Event_FEATURE_CREATED, &eventproto.FeatureCreatedEvent{Id: "fId-0"}, "ns0")
	assert.NoError(t, err)
	event1, err := domainevent.NewEvent(editor, eventproto.Event_FEATURE, "fId-1", eventproto.Event_FEATURE_CREATED, &eventproto.FeatureCreatedEvent{Id: "fId-1"}, "ns0")
	assert.NoError(t, err)
	adninEvent0, err := domainevent.NewAdminEvent(editor, eventproto.Event_FEATURE, "fId-2", eventproto.Event_FEATURE_CREATED, &eventproto.FeatureCreatedEvent{Id: "fId-2"})
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

func newPersister(t *testing.T, mockController *gomock.Controller) *Persister {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &Persister{
		puller: pullermock.NewMockRateLimitedPuller(mockController),
		logger: logger.Named("persister"),
		ctx:    ctx,
		cancel: cancel,
		doneCh: make(chan struct{}),
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
