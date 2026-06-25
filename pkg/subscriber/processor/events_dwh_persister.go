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
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	ecdwh "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/dwh_database"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	dwhstorage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage"
	operationalstorage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/operationalstorage"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

type eventsDWHPersisterConfig struct {
	// Persister-specific settings
	FlushInterval           int `json:"flushInterval"`
	FlushTimeout            int `json:"flushTimeout"`
	FlushSize               int `json:"flushSize"`
	MaxRetryGoalEventPeriod int `json:"maxRetryGoalEventPeriod,omitempty"`
	RetryGoalEventInterval  int `json:"retryGoalEventInterval,omitempty"`
}

type eventsDWHPersister struct {
	eventsDWHPersisterConfig eventsDWHPersisterConfig
	experimentStorage        operationalstorage.ExperimentStorage
	writer                   Writer
	subscriberType           string
	logger                   *zap.Logger
}

// NewEventsDWHPersister builds a data-warehouse events persister.
//
// All data-warehouse storages (evalEventWriter / goalEventWriter / goalEventStorage) are
// created by the server and injected here — the processor no longer depends on any DWH
// client or dialect. experimentStorage is backed by the operational database (used by Switch()).
func NewEventsDWHPersister(
	ctx context.Context,
	config interface{},
	evalEventWriter dwhstorage.EvalEventWriter,
	goalEventWriter dwhstorage.GoalEventWriter,
	goalEventStorage ecdwh.EventStorage,
	experimentStorage operationalstorage.ExperimentStorage,
	location *time.Location,
	redisClient redisv3.Client,
	persistentRedisClient redisv3.Client,
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
		experimentStorage:        experimentStorage,
		logger:                   logger,
	}
	experimentsCache := cachev3.NewExperimentsCache(cachev3.NewRedisCache(redisClient))

	switch persisterName {
	case EvaluationCountEventDWHPersisterName:
		e.subscriberType = subscriberEvaluationEventDWH
		e.writer = NewEvalEventWriter(logger, exClient, experimentsCache, location, evalEventWriter)

	case GoalCountEventDWHPersisterName:
		e.subscriberType = subscriberGoalEventDWH
		maxRetryPeriod := time.Duration(e.eventsDWHPersisterConfig.MaxRetryGoalEventPeriod) * time.Second
		retryInterval := time.Duration(e.eventsDWHPersisterConfig.RetryGoalEventInterval) * time.Second
		e.writer = NewGoalEventWriter(
			ctx,
			logger,
			exClient,
			ftClient,
			experimentsCache,
			location,
			persistentRedisClient,
			maxRetryPeriod,
			retryInterval,
			goalEventWriter,
			goalEventStorage,
		)
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
	count, err := e.experimentStorage.CountRunningExperiments(ctx)
	if err != nil {
		e.logger.Error("Failed to count running experiments", zap.Error(err))
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
		e.logger.Error("Bad proto message",
			zap.Error(err),
			zap.String("messageID", m.ID),
			zap.ByteString("data", m.Data),
			zap.Any("attributes", m.Attributes),
		)
		subscriberHandledCounter.WithLabelValues(e.subscriberType, codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		// Check if message data is empty
		if len(m.Data) == 0 {
			handleBadMessage(m, fmt.Errorf("message data is empty"))
			continue
		}
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		innerEvent, err := event.Event.UnmarshalNew()
		if err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentId]; ok {
			innerEvents[event.Id] = innerEvent
			continue
		}
		envEvents[event.EnvironmentId] = eventDWHMap{event.Id: innerEvent}
	}
	return envEvents
}
