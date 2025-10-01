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
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

type APIKeyLastUsedInfoWriterConfig struct {
	FlushSize           int `json:"flushSize"`
	FlushInterval       int `json:"flushInterval"`
	WriteCacheInterval  int `json:"writeCacheInterval"`
	UserAttributeKeyTTL int `json:"userAttributeKeyTtl"`
}

type apikeyLastUsedInfoCache map[string]*domain.APIKeyLastUsedInfo

type envAPIKeyLastUsedInfoCache map[string]apikeyLastUsedInfoCache

type apiKeyUsageEventMap map[string]*eventproto.APIKeyUsageEvent

type envAPIKeyUsageEventMap map[string]apiKeyUsageEventMap

type apikeyLastUsedInfoWriter struct {
	config                   APIKeyLastUsedInfoWriterConfig
	apikeyLastUsedInfoCacher envAPIKeyLastUsedInfoCache

	envLastUsedCacheMutex sync.Mutex

	mysqlClient mysql.Client
	logger      *zap.Logger
}

func NewAPIKeyLastUsedInfoWriter(
	config interface{},
	mysqlClient mysql.Client,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	apikeyLastUsedInfWriterConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("apikeyLastUsedInfoWriter: invalid config")
		return nil, ErrAPIKeyLastUsedInfoWriterInvalidConfig
	}
	configBytes, err := json.Marshal(apikeyLastUsedInfWriterConfig)
	if err != nil {
		logger.Error("apikeyLastUsedInfoWriter: failed to marshal config", zap.Error(err))
		return nil, ErrAPIKeyLastUsedInfoWriterInvalidConfig
	}
	var apikeyLastUsedInfoWriterConfig APIKeyLastUsedInfoWriterConfig
	err = json.Unmarshal(configBytes, &apikeyLastUsedInfoWriterConfig)
	if err != nil {
		logger.Error("apikeyLastUsedInfoWriter: failed to unmarshal config", zap.Error(err))
		return nil, ErrAPIKeyLastUsedInfoWriterInvalidConfig
	}
	w := &apikeyLastUsedInfoWriter{
		config:      apikeyLastUsedInfoWriterConfig,
		mysqlClient: mysqlClient,
		logger:      logger,
	}

	return w, nil
}

func (a *apikeyLastUsedInfoWriter) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	batch := make(map[string]*puller.Message)
	ticket := time.NewTicker(time.Duration(a.config.UserAttributeKeyTTL) * time.Second)
	defer ticket.Stop()

	resetBatch := func() {
		for _, msg := range batch {
			msg.Ack()
			subscriberHandledCounter.WithLabelValues(subscriberAPIKeyLastUsedInfo, codes.OK.String()).Inc()
		}
		batch = make(map[string]*puller.Message)
	}

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				a.logger.Error("Failed to pull message")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberAPIKeyLastUsedInfo).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				a.logger.Error("apikeyLastUsedInfoWriter: id is empty")
				subscriberHandledCounter.WithLabelValues(subscriberAPIKeyLastUsedInfo, codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				subscriberHandledCounter.WithLabelValues(subscriberAPIKeyLastUsedInfo, codes.MissingID.String()).Inc()
				previous.Ack()
			}
			batch[id] = msg
			if len(batch) < a.config.FlushSize {
				continue
			}
			envEvents := a.extractEvents(batch)
			a.cacheAPIKeyLastUsedInfoPerEnv(envEvents)
			resetBatch()
		case <-ticket.C:
			envEvents := a.extractEvents(batch)
			a.cacheAPIKeyLastUsedInfoPerEnv(envEvents)
			resetBatch()
		case <-ctx.Done():
			// Nack the messages to be redelivered
			for _, msg := range batch {
				msg.Nack()
			}
			a.logger.Debug("All the left messages were Nack successfully before shutting down",
				zap.Int("batchSize", len(batch)))
			return nil
		}
	}
}

func (a *apikeyLastUsedInfoWriter) cacheAPIKeyLastUsedInfoPerEnv(envEvents envAPIKeyUsageEventMap) {
	for environmentID, events := range envEvents {
		for _, event := range events {
			a.cacheEnvAPIKeyLastUsedInfo(event, environmentID)
		}
		a.logger.Debug("Update cache API key last used info",
			zap.String("environmentID", environmentID),
			zap.Int("cacheSize", len(a.apikeyLastUsedInfoCacher[environmentID])),
			zap.Int("eventSize", len(events)),
		)
	}
}

func (a *apikeyLastUsedInfoWriter) cacheEnvAPIKeyLastUsedInfo(
	event *eventproto.APIKeyUsageEvent,
	environmentID string,
) {
	a.envLastUsedCacheMutex.Lock()
	defer a.envLastUsedCacheMutex.Unlock()
	if cache, ok := a.apikeyLastUsedInfoCacher[environmentID]; ok {
		if info, ok := cache[event.ApiKeyId]; ok {
			info.UsedAt(event.Timestamp)
			return
		}
		cache[event.ApiKeyId] = domain.NewAPIKeyLastUsedInfo(event.ApiKeyId, event.Timestamp, environmentID)
		return
	}
	cache := apikeyLastUsedInfoCache{}
	cache[event.ApiKeyId] = domain.NewAPIKeyLastUsedInfo(event.ApiKeyId, event.Timestamp, environmentID)
	a.apikeyLastUsedInfoCacher[environmentID] = cache
}

func (a *apikeyLastUsedInfoWriter) extractEvents(messages map[string]*puller.Message) envAPIKeyUsageEventMap {
	envEvents := envAPIKeyUsageEventMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		a.logger.Error("Bad proto message",
			zap.Error(err),
			zap.String("messageID", m.ID),
			zap.ByteString("data", m.Data),
			zap.Any("attributes", m.Attributes),
		)
		subscriberHandledCounter.WithLabelValues(subscriberAPIKeyLastUsedInfo, codes.BadMessage.String()).Inc()
	}
	for _, msg := range messages {
		// check if data is empty
		if len(msg.Data) == 0 {
			handleBadMessage(msg, fmt.Errorf("message data is empty"))
			continue
		}
		event := &eventproto.Event{}
		if err := proto.Unmarshal(msg.Data, event); err != nil {
			handleBadMessage(msg, err)
			continue
		}
		innerEvent := &eventproto.APIKeyUsageEvent{}
		if err := anypb.UnmarshalTo(event.Event, innerEvent, proto.UnmarshalOptions{}); err != nil {
			handleBadMessage(msg, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentId]; ok {
			innerEvents[event.Id] = innerEvent
			continue
		}
		envEvents[event.EnvironmentId] = apiKeyUsageEventMap{event.Id: innerEvent}
	}
	return envEvents
}
