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
	"encoding/json"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	ustorage "github.com/bucketeer-io/bucketeer/pkg/subscriber/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	ecproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/service"
)

type userEventPersisterConfig struct {
	FlushSize     int `json:"flushSize"`
	FlushInterval int `json:"flushInterval"`
}

type userEventPersister struct {
	userEventPersisterConfig userEventPersisterConfig
	timeNow                  func() time.Time
	newUUID                  func() (*uuid.UUID, error)
	mysqlClient              mysql.Client
	logger                   *zap.Logger
}

func NewUserEventPersister(
	config interface{},
	mysqlClient mysql.Client,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	userEventPerisiterJsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("UserEventPersister: invalid config")
		return nil, ErrUserEventInvalidConfig
	}
	configBytes, err := json.Marshal(userEventPerisiterJsonConfig)
	if err != nil {
		logger.Error("UserEventPersister: failed to marshal config", zap.Error(err))
		return nil, err
	}
	var userEventPeristerConfig userEventPersisterConfig
	err = json.Unmarshal(configBytes, &userEventPeristerConfig)
	if err != nil {
		logger.Error("UserEventPersister: failed to unmarshal config", zap.Error(err))
		return nil, err
	}
	return &userEventPersister{
		userEventPersisterConfig: userEventPeristerConfig,
		mysqlClient:              mysqlClient,
		timeNow:                  time.Now,
		newUUID:                  uuid.NewUUID,
		logger:                   logger,
	}, nil
}

func (p *userEventPersister) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	chunk := make(map[string]*puller.Message, p.userEventPersisterConfig.FlushSize)
	ticker := time.NewTicker(time.Duration(p.userEventPersisterConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberUserEvent).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				subscriberHandledCounter.WithLabelValues(subscriberUserEvent, codes.MissingID.String()).Inc()
				continue
			}
			if pre, ok := chunk[id]; ok {
				pre.Ack()
				subscriberHandledCounter.WithLabelValues(subscriberUserEvent, codes.DuplicateID.String()).Inc()
			}
			chunk[id] = msg
			if len(chunk) >= p.userEventPersisterConfig.FlushSize {
				p.handleChunk(chunk)
				chunk = make(map[string]*puller.Message, p.userEventPersisterConfig.FlushSize)
			}
		case <-ticker.C:
			if len(chunk) > 0 {
				p.handleChunk(chunk)
				chunk = make(map[string]*puller.Message, p.userEventPersisterConfig.FlushSize)
			}
		case <-ctx.Done():
			chunkSize := len(chunk)
			p.logger.Debug("Context is done", zap.Int("chunkSize", chunkSize))
			if chunkSize > 0 {
				p.handleChunk(chunk)
				p.logger.Debug(
					"All the left messages are processed successfully",
					zap.Int("chunkSize", chunkSize),
				)
			}
			return nil
		}
	}
}

func (p *userEventPersister) handleChunk(chunk map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	events := make(map[string][]*eventproto.UserEvent)
	messages := make(map[string][]*puller.Message)
	for _, msg := range chunk {
		// Extract the user event
		event := p.extractUserEvent(msg)
		if event == nil {
			continue
		}
		// Append events per environment
		listEvents, ok := events[event.EnvironmentId]
		if ok {
			events[event.EnvironmentId] = append(listEvents, event)
		} else {
			events[event.EnvironmentId] = []*eventproto.UserEvent{event}
		}
		// Append PubSub messages per environment
		listMessages, ok := messages[event.EnvironmentId]
		if ok {
			messages[event.EnvironmentId] = append(listMessages, msg)
		} else {
			messages[event.EnvironmentId] = []*puller.Message{msg}
		}
	}
	// Upsert events
	for environmentId, events := range events {
		// Upsert events per environment
		err := p.upsertMAUs(ctx, events, environmentId)
		if err != nil {
			p.nackMessages(messages[environmentId])
		} else {
			p.ackMessages(messages[environmentId])
		}
	}
}

func (p *userEventPersister) extractUserEvent(message *puller.Message) *eventproto.UserEvent {
	event, err := p.unmarshalMessage(message)
	if err != nil {
		message.Nack()
		subscriberHandledCounter.WithLabelValues(subscriberUserEvent, codes.BadMessage.String()).Inc()
		return nil
	}
	if !p.validateEvent(event) {
		message.Nack()
		subscriberHandledCounter.WithLabelValues(subscriberUserEvent, codes.BadMessage.String()).Inc()
		return nil
	}
	return event
}

func (p *userEventPersister) unmarshalMessage(msg *puller.Message) (*eventproto.UserEvent, error) {
	event := &ecproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		return nil, err
	}
	var userEvent eventproto.UserEvent
	if err := ptypes.UnmarshalAny(event.Event, &userEvent); err != nil {
		p.logger.Error("Failed to unmarshal Event -> UserEvent", zap.Error(err), zap.Any("msg", msg))
		return nil, err
	}
	return &userEvent, err
}

func (p *userEventPersister) validateEvent(event *eventproto.UserEvent) bool {
	if event.UserId == "" {
		p.logger.Warn("Message contains an empty User Id", zap.Any("event", event))
		return false
	}
	if event.LastSeen == 0 {
		p.logger.Warn("Message's LastSeen is zero", zap.Any("event", event))
		return false
	}
	return true
}

func (p *userEventPersister) nackMessages(messages []*puller.Message) {
	for _, msg := range messages {
		msg.Nack()
		subscriberHandledCounter.WithLabelValues(subscriberUserEvent, codes.RepeatableError.String()).Inc()
	}
}

func (p *userEventPersister) ackMessages(messages []*puller.Message) {
	for _, msg := range messages {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberUserEvent, codes.OK.String()).Inc()
	}
}

func (p *userEventPersister) upsertMAUs(
	ctx context.Context,
	events []*eventproto.UserEvent,
	environmentId string,
) error {
	s := ustorage.NewMysqlMAUStorage(p.mysqlClient)
	if err := s.UpsertMAUs(ctx, events, environmentId); err != nil {
		p.logger.Error("Failed to upsert user events",
			zap.Error(err),
			zap.String("environmentId", environmentId),
			zap.Int("size", len(events)),
		)
		return err
	}
	return nil
}
