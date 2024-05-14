//  Copyright 2024 The Bucketeer Authors.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package processor

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	storage "github.com/bucketeer-io/bucketeer/pkg/eventpersisterdwh/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

type eventDWHMap map[string]proto.Message
type environmentEventDWHMap map[string]eventDWHMap

type evaluationEventsDWHPersisterConfig struct {
	FlushInterval int `json:"flushInterval"`
	FlushTimeout  int `json:"flushTimeout"`
	FlushSize     int `json:"flushSize"`
}

type evaluationEventsDWHPersister struct {
	evaluationEventsDWHPersisterConfig evaluationEventsDWHPersisterConfig
	mysqlClient                        mysql.Client
	writer                             Writer
	logger                             *zap.Logger
}

func (e *evaluationEventsDWHPersister) Process(
	ctx context.Context,
	msgChan <-chan *puller.Message,
) error {
	batch := make(map[string]*puller.Message)
	ticker := time.NewTicker(time.Duration(e.evaluationEventsDWHPersisterConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberEvaluationEventDWH).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				// TODO: better log format for msg data
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				previous.Ack()
				e.logger.Warn("Message with duplicate id", zap.String("id", id))
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codes.DuplicateID.String()).Inc()
			}
			batch[id] = msg
			if len(batch) < e.evaluationEventsDWHPersisterConfig.FlushSize {
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
			e.logger.Info("Context is done", zap.Int("batchSize", batchSize))
			if len(batch) > 0 {
				e.send(batch)
				e.logger.Info("All the left messages are processed successfully", zap.Int("batchSize", batchSize))
			}
			return nil
		}
	}
}

func (e *evaluationEventsDWHPersister) Switch(ctx context.Context) (bool, error) {
	experimentStorage := storage.NewExperimentStorage(e.mysqlClient)
	count, err := experimentStorage.CountRunningExperiments(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (e *evaluationEventsDWHPersister) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(e.evaluationEventsDWHPersisterConfig.FlushTimeout)*time.Second,
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
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codes.RepeatableError.String()).Inc()
			} else {
				m.Ack()
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codes.NonRepeatableError.String()).Inc()
			}
			continue
		}
		m.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codes.OK.String()).Inc()
	}
}
func (e *evaluationEventsDWHPersister) extractEvents(messages map[string]*puller.Message) environmentEventDWHMap {
	envEvents := environmentEventDWHMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		e.logger.Error("bad message", zap.Error(err), zap.Any("msg", m))
		subscriberHandledCounter.WithLabelValues(subscriberEvaluationEventDWH, codes.BadMessage.String()).Inc()
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
		if innerEvents, ok := envEvents[event.EnvironmentNamespace]; ok {
			innerEvents[event.Id] = innerEvent.Message
			continue
		}
		envEvents[event.EnvironmentNamespace] = eventDWHMap{event.Id: innerEvent.Message}
	}
	return envEvents
}
