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
	"errors"
	"time"

	pb "github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/auditlog/domain"
	v2als "github.com/bucketeer-io/bucketeer/pkg/auditlog/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	domainevent "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type auditLogPersisterConfig struct {
	FlushSize     int `json:"flushSize"`
	FlushInterval int `json:"flushInterval"`
	FlushTimeout  int `json:"flushTimeout"`
}

type auditLogPersister struct {
	auditLogPersisterConfig auditLogPersisterConfig
	mysqlAdminStorage       v2als.AdminAuditLogStorage
	mysqlStorage            v2als.AuditLogStorage
	logger                  *zap.Logger
}

func NewAuditLogPersister(
	config interface{},
	mysqlClient mysql.Client,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	auditLogPersisterJsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("AuditLogPersister: invalid config")
		return nil, errAuditLogInvalidConfig
	}
	configBytes, err := json.Marshal(auditLogPersisterJsonConfig)
	if err != nil {
		logger.Error("AuditLogPersister: failed to marshal config", zap.Error(err))
		return nil, err
	}
	var persisterConfig auditLogPersisterConfig
	err = json.Unmarshal(configBytes, &persisterConfig)
	if err != nil {
		logger.Error("AuditLogPersister: failed to unmarshal config", zap.Error(err))
		return nil, err
	}
	return &auditLogPersister{
		auditLogPersisterConfig: persisterConfig,
		mysqlAdminStorage:       v2als.NewAdminAuditLogStorage(mysqlClient),
		mysqlStorage:            v2als.NewAuditLogStorage(mysqlClient),
		logger:                  logger,
	}, nil
}

func (a auditLogPersister) Process(
	ctx context.Context,
	msgChan <-chan *puller.Message,
) error {
	chunk := make(map[string]*puller.Message, a.auditLogPersisterConfig.FlushSize)
	ticker := time.NewTicker(time.Duration(a.auditLogPersisterConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberAuditLog).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				subscriberHandledCounter.WithLabelValues(subscriberAuditLog, codes.MissingID.String()).Inc()
				continue
			}
			if _, ok := chunk[id]; ok {
				subscriberHandledCounter.WithLabelValues(subscriberAuditLog, codes.DuplicateID.String()).Inc()
			}
			chunk[id] = msg
			if len(chunk) >= a.auditLogPersisterConfig.FlushSize {
				a.flushChunk(chunk)
				chunk = make(map[string]*puller.Message, a.auditLogPersisterConfig.FlushSize)
			}
		case <-ticker.C:
			if len(chunk) > 0 {
				a.flushChunk(chunk)
				chunk = make(map[string]*puller.Message, a.auditLogPersisterConfig.FlushSize)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (a auditLogPersister) flushChunk(chunk map[string]*puller.Message) {
	auditlogs, adminAuditLogs, messages, adminMessages := a.extractAuditLogs(chunk)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(a.auditLogPersisterConfig.FlushTimeout)*time.Second,
	)
	defer cancel()
	// Environment audit logs
	a.createAuditLogsMySQL(ctx, auditlogs, messages, a.mysqlStorage.CreateAuditLog)
	// Admin audit logs
	a.createAuditLogsMySQL(ctx, adminAuditLogs, adminMessages, a.mysqlAdminStorage.CreateAdminAuditLog)
}

func (a auditLogPersister) extractAuditLogs(
	chunk map[string]*puller.Message,
) (auditlogs, adminAuditLogs []*domain.AuditLog, messages, adminMessages []*puller.Message) {
	for _, msg := range chunk {
		event := &domainevent.Event{}
		if err := pb.Unmarshal(msg.Data, event); err != nil {
			a.logger.Error("Failed to unmarshal message", zap.Error(err))
			subscriberHandledCounter.WithLabelValues(subscriberAuditLog, codes.BadMessage.String()).Inc()
			msg.Ack()
			continue
		}
		if event.IsAdminEvent {
			adminAuditLogs = append(adminAuditLogs, domain.NewAuditLog(
				event,
				storage.AdminEnvironmentID,
			))
			adminMessages = append(adminMessages, msg)
		} else {
			auditlogs = append(auditlogs, domain.NewAuditLog(event, event.EnvironmentId))
			messages = append(messages, msg)
		}
	}
	return
}

func (a auditLogPersister) createAuditLogsMySQL(
	ctx context.Context,
	auditlogs []*domain.AuditLog,
	messages []*puller.Message,
	createFunc func(ctx context.Context, auditLog *domain.AuditLog) error,
) {
	for i, aud := range auditlogs {
		if err := createFunc(ctx, aud); err != nil {
			if errors.Is(err, v2als.ErrAuditLogAlreadyExists) || errors.Is(err, v2als.ErrAdminAuditLogAlreadyExists) {
				subscriberHandledCounter.WithLabelValues(subscriberAuditLog, codes.NonRepeatableError.String()).Inc()
				messages[i].Ack()
			} else {
				var publicApiEditor string
				if aud.Editor.PublicApiEditor != nil {
					publicApiEditor = aud.Editor.PublicApiEditor.Maintainer
				}
				a.logger.Error("Failed to put audit logs",
					zap.Error(err),
					zap.String("entity_id", aud.EntityId),
					zap.String("editor", aud.Editor.Email),
					zap.String("public_api_editor", publicApiEditor),
					zap.String("type", aud.Type.String()),
					zap.Int("entity_data_size", len(aud.EntityData)),
					zap.Int("previous_entity_data_size", len(aud.PreviousEntityData)),
					zap.Any("entity_type", aud.EntityType),
					zap.String("entity_data", aud.EntityData),
					zap.String("previous_entity_data", aud.PreviousEntityData),
				)
				subscriberHandledCounter.WithLabelValues(subscriberAuditLog, codes.RepeatableError.String()).Inc()
				messages[i].Nack()
			}
			continue
		}
		subscriberHandledCounter.WithLabelValues(subscriberAuditLog, codes.OK.String()).Inc()
		messages[i].Ack()
	}
}
