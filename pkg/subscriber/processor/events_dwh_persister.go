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

	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	storage "github.com/bucketeer-io/bucketeer/pkg/subscriber/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

type eventsDWHPersisterConfig struct {
	FlushInterval     int    `json:"flushInterval"`
	FlushTimeout      int    `json:"flushTimeout"`
	FlushSize         int    `json:"flushSize"`
	Project           string `json:"project"`
	BigQueryDataSet   string `json:"bigQueryDataSet"`
	BigQueryBatchSize int    `json:"bigQueryBatchSize"`
	Timezone          string `json:"timezone"`
}

type eventsDWHPersister struct {
	eventsDWHPersisterConfig eventsDWHPersisterConfig
	mysqlClient              mysql.Client
	writer                   Writer
	subscriberType           string
	logger                   *zap.Logger
}

func NewEventsDWHPersister(
	ctx context.Context,
	config interface{},
	mysqlClient mysql.Client,
	redisClient redisv3.Client,
	exClient experimentclient.Client,
	ftClient featureclient.Client,
	persisterName string,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	jsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("eventsDWHPersister: invalid config")
		return nil, ErrEventsDWHPersisterInvalidConfig
	}
	configBytes, err := json.Marshal(jsonConfig)
	if err != nil {
		logger.Error("eventsDWHPersister: failed to marshal config", zap.Error(err))
		return nil, err
	}
	var persisterConfig eventsDWHPersisterConfig
	err = json.Unmarshal(configBytes, &persisterConfig)
	if err != nil {
		logger.Error("eventsDWHPersister: failed to unmarshal config", zap.Error(err))
		return nil, err
	}
	e := &eventsDWHPersister{
		eventsDWHPersisterConfig: persisterConfig,
		mysqlClient:              mysqlClient,
		logger:                   logger,
	}
	experimentsCache := cachev3.NewExperimentsCache(cachev3.NewRedisCache(redisClient))
	location, err := locale.GetLocation(e.eventsDWHPersisterConfig.Timezone)
	if err != nil {
		return nil, err
	}
	switch persisterName {
	case EvaluationCountEventDWHPersisterName:
		e.subscriberType = subscriberEvaluationEventDWH
		evalEventWriter, err := NewEvalEventWriter(
			ctx,
			logger,
			exClient,
			experimentsCache,
			e.eventsDWHPersisterConfig.Project,
			e.eventsDWHPersisterConfig.BigQueryDataSet,
			e.eventsDWHPersisterConfig.BigQueryBatchSize,
			location,
		)
		if err != nil {
			return nil, err
		}
		e.writer = evalEventWriter
	case GoalCountEventDWHPersisterName:
		e.subscriberType = subscriberGoalEventDWH
		goalEventWriter, err := NewGoalEventWriter(
			ctx,
			logger,
			exClient,
			ftClient,
			experimentsCache,
			e.eventsDWHPersisterConfig.Project,
			e.eventsDWHPersisterConfig.BigQueryDataSet,
			e.eventsDWHPersisterConfig.BigQueryBatchSize,
			location,
		)
		if err != nil {
			return nil, err
		}
		e.writer = goalEventWriter
	}
	return e, nil
}

func (e *eventsDWHPersister) Process(
	ctx context.Context,
	msgChan <-chan *puller.Message,
) error {
	batch := make(map[string]*puller.Message)
	ticker := time.NewTicker(time.Duration(e.eventsDWHPersisterConfig.FlushInterval) * time.Second)
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
			if len(batch) < e.eventsDWHPersisterConfig.FlushSize {
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
				e.logger.Debug(
					"All the left messages are processed successfully",
					zap.Int("batchSize", batchSize),
				)
			}
			return nil
		}
	}
}

func (e *eventsDWHPersister) Switch(ctx context.Context) (bool, error) {
	experimentStorage := storage.NewExperimentStorage(e.mysqlClient)
	count, err := experimentStorage.CountRunningExperiments(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *eventsDWHPersister) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(e.eventsDWHPersisterConfig.FlushTimeout)*time.Second,
	)
	defer cancel()
	envEvents := e.extractEvents(messages)
	if len(envEvents) == 0 {
		e.logger.Error("all messages were bad")
		return
	}
	fails := e.writer.Write(ctx, envEvents)
	for id, m := range messages {
		if repeatable, ok := fails[id]; ok {
			if repeatable {
				m.Nack()
				subscriberHandledCounter.WithLabelValues(
					e.subscriberType,
					codes.RepeatableError.String(),
				).Inc()
			} else {
				m.Ack()
				subscriberHandledCounter.WithLabelValues(
					e.subscriberType,
					codes.NonRepeatableError.String(),
				).Inc()
			}
			continue
		}
		m.Ack()
		subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.OK.String()).Inc()
	}
}
func (e *eventsDWHPersister) extractEvents(messages map[string]*puller.Message) environmentEventDWHMap {
	envEvents := environmentEventDWHMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		e.logger.Error("bad message", zap.Error(err), zap.Any("msg", m))
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
		envEvents[event.EnvironmentId] = eventDWHMap{event.Id: innerEvent.Message}
	}
	return envEvents
}
