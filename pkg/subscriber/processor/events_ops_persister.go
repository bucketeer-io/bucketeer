// Copyright 2026 The Bucketeer Authors.
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

	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	storage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

type eventsOPSPersisterConfig struct {
	FlushInterval int `json:"flushInterval"`
	FlushTimeout  int `json:"flushTimeout"`
	FlushSize     int `json:"flushSize"`
}

type eventsOPSPersister struct {
	eventsOPSPersisterConfig eventsOPSPersisterConfig
	mysqlClient              mysql.Client
	updater                  Updater
	subscriberType           string
	logger                   *zap.Logger
}

func NewEventsOPSPersister(
	ctx context.Context,
	config interface{},
	mysqlClient mysql.Client,
	redisClient redisv3.Client,
	opsClient autoopsclient.Client,
	ftClient featureclient.Client,
	persisterName string,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	jsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("eventsOPSPersister: invalid config")
		return nil, ErrEventsOPSPersisterInvalidConfig
	}
	configBytes, err := json.Marshal(jsonConfig)
	if err != nil {
		logger.Error("eventsOPSPersister: failed to marshal config", zap.Error(err))
		return nil, err
	}
	var persisterConfig eventsOPSPersisterConfig
	err = json.Unmarshal(configBytes, &persisterConfig)
	if err != nil {
		logger.Error("eventsOPSPersister: failed to unmarshal config", zap.Error(err))
		return nil, err
	}
	e := &eventsOPSPersister{
		eventsOPSPersisterConfig: persisterConfig,
		mysqlClient:              mysqlClient,
		logger:                   logger,
	}
	switch persisterName {
	case EvaluationCountEventOPSPersisterName:
		e.subscriberType = subscriberEvaluationEventOPS
		e.updater = NewEvalUserCountUpdater(
			ctx,
			ftClient,
			opsClient,
			cachev3.NewEventCountCache(cachev3.NewRedisCache(redisClient)),
			cachev3.NewAutoOpsRulesCache(cachev3.NewRedisCache(redisClient)),
			logger,
		)
	case GoalCountEventOPSPersisterName:
		e.subscriberType = subscriberGoalEventOPS
		e.updater = NewGoalUserCountUpdater(
			ctx,
			ftClient,
			opsClient,
			cachev3.NewEventCountCache(cachev3.NewRedisCache(redisClient)),
			cachev3.NewAutoOpsRulesCache(cachev3.NewRedisCache(redisClient)),
			logger,
		)
	}
	return e, nil
}

func (e eventsOPSPersister) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	batch := make(map[string]*puller.Message)
	ticker := time.NewTicker(time.Duration(e.eventsOPSPersisterConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(e.subscriberType).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				// TODO: better log format for msg data
				subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				previous.Ack()
				subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.DuplicateID.String()).Inc()
			}
			batch[id] = msg
			if len(batch) < e.eventsOPSPersisterConfig.FlushSize {
				continue
			}
			e.send(batch)
			batch = make(map[string]*puller.Message)
		case <-ticker.C:
			if len(batch) > 0 {
				e.send(batch)
				batch = make(map[string]*puller.Message)
			}
		case <-ctx.Done():
			batchSize := len(batch)
			e.logger.Debug("Context is done", zap.Int("batchSize", batchSize))
			if len(batch) > 0 {
				e.send(batch)
				e.logger.Debug("All the left messages are processed successfully", zap.Int("batchSize", batchSize))
			}
			return nil
		}
	}
}

func (e eventsOPSPersister) Switch(ctx context.Context) (bool, error) {
	autoOpsRuleStorage := storage.NewAutoOpsRuleStorage(e.mysqlClient)
	autoOpsRuleCount, err := autoOpsRuleStorage.CountOpsEventRate(ctx)
	if err != nil {
		return false, err
	}
	return autoOpsRuleCount > 0, nil
}

func (e eventsOPSPersister) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(e.eventsOPSPersisterConfig.FlushTimeout)*time.Second,
	)
	defer cancel()
	envEvents := e.extractEvents(messages)
	if len(envEvents) == 0 {
		e.logger.Error("All messages were bad")
		return
	}
	fails := e.updater.UpdateUserCounts(ctx, envEvents)
	for id, m := range messages {
		if repeatable, ok := fails[id]; ok {
			if repeatable {
				m.Nack()
				subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.RepeatableError.String()).Inc()
			} else {
				m.Ack()
				subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.NonRepeatableError.String()).Inc()
			}
			continue
		}
		m.Ack()
		subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.OK.String()).Inc()
	}
}

func (e eventsOPSPersister) extractEvents(messages map[string]*puller.Message) environmentEventOPSMap {
	envEvents := environmentEventOPSMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		e.logger.Error("Bad message", zap.Error(err), zap.Any("msg", m))
		subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		var innerEvent ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(event.Event, &innerEvent); err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentId]; ok {
			innerEvents[event.Id] = innerEvent.Message
			continue
		}
		envEvents[event.EnvironmentId] = eventOPSMap{event.Id: innerEvent.Message}
	}
	return envEvents
}
