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

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/email"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type EmailSenderConfig struct {
	// Add config fields if needed in the future
}

type emailSender struct {
	emailService email.Service
	config       EmailSenderConfig
	logger       *zap.Logger
}

func NewEmailSender(
	config interface{},
	emailService email.Service,
	logger *zap.Logger,
) subscriber.PubSubProcessor {
	var emailSenderConfig EmailSenderConfig
	if config != nil {
		jsonConfigMap, ok := config.(map[string]interface{})
		if ok {
			configBytes, err := json.Marshal(jsonConfigMap)
			if err != nil {
				logger.Warn("emailSender: failed to marshal config, using defaults", zap.Error(err))
			} else {
				if err := json.Unmarshal(configBytes, &emailSenderConfig); err != nil {
					logger.Warn("emailSender: failed to unmarshal config, using defaults", zap.Error(err))
				}
			}
		}
	}

	return &emailSender{
		emailService: emailService,
		config:       emailSenderConfig,
		logger:       logger,
	}
}

func (e *emailSender) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				e.logger.Error("emailSender: message channel closed")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberEmailSender).Inc()
			e.handleMessage(msg)
		case <-ctx.Done():
			e.logger.Debug("subscriber context done, stopped processing messages")
			return nil
		}
	}
}

func (e *emailSender) handleMessage(msg *puller.Message) {
	startTime := time.Now()

	// Check for message ID
	if id := msg.Attributes["id"]; id == "" {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberEmailSender, codes.MissingID.String()).Inc()
		return
	}

	// Unmarshal domain event
	domainEvent, err := e.unmarshalMessage(msg)
	if err != nil {
		e.logger.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.String("msgID", msg.ID),
		)
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberEmailSender, codes.BadMessage.String()).Inc()
		return
	}

	// Filter: Only handle ACCOUNT_V2_CREATED events
	if domainEvent.EntityType != domaineventproto.Event_ACCOUNT ||
		domainEvent.Type != domaineventproto.Event_ACCOUNT_V2_CREATED {
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberEmailSender, codes.OK.String()).Inc()
		return
	}

	// Extract AccountV2CreatedEvent
	var accountCreatedEvent domaineventproto.AccountV2CreatedEvent
	if err := domainEvent.Data.UnmarshalTo(&accountCreatedEvent); err != nil {
		e.logger.Error("Failed to unmarshal AccountV2CreatedEvent",
			zap.String("event id", domainEvent.Id),
			zap.Error(err),
		)
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberEmailSender, codes.BadMessage.String()).Inc()
		return
	}

	// Send welcome email
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.emailService.SendWelcomeEmail(ctx, accountCreatedEvent.Email, accountCreatedEvent.Language); err != nil {
		e.logger.Error("Failed to send welcome email",
			zap.Error(err),
			zap.String("event id", domainEvent.Id),
			zap.String("email", accountCreatedEvent.Email),
		)
		// Don't retry on email failures - acknowledge the message
		msg.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberEmailSender, codes.NonRepeatableError.String()).Inc()
		subscriberHandledHistogram.WithLabelValues(subscriberEmailSender, codes.NonRepeatableError.String()).
			Observe(time.Since(startTime).Seconds())
		return
	}

	e.logger.Info("Successfully sent welcome email",
		zap.String("event id", domainEvent.Id),
		zap.String("email", accountCreatedEvent.Email),
	)

	msg.Ack()
	subscriberHandledCounter.WithLabelValues(subscriberEmailSender, codes.OK.String()).Inc()
	subscriberHandledHistogram.WithLabelValues(subscriberEmailSender, codes.OK.String()).
		Observe(time.Since(startTime).Seconds())
}

func (e *emailSender) unmarshalMessage(msg *puller.Message) (*domaineventproto.Event, error) {
	event := &domaineventproto.Event{}
	err := proto.Unmarshal(msg.Data, event)
	if err != nil {
		e.logger.Error("Failed to unmarshal message", zap.Error(err))
		return nil, err
	}
	return event, nil
}
