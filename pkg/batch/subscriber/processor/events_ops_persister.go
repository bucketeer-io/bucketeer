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
	"encoding/json"
	autoopsclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	"github.com/bucketeer-io/bucketeer/pkg/batch/subscriber"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"go.uber.org/zap"
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
	persistentRedisClient redisv3.Client,
	opsClient autoopsclient.Client,
	ftClient featureclient.Client,
	persisterName string,
	registerer metrics.Registerer,
	logger *zap.Logger,
) (subscriber.Processor, error) {
	jsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("eventsOPSPersister: invalid config")
		return nil, ErrSegmentInvalidConfig
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
		eventsOPSPersisterConfig: eventsOPSPersisterConfig{},
		mysqlClient:              mysqlClient,
		logger:                   logger,
	}
	switch persisterName {
	case EvaluationCountEventOPSPersisterName:
		e.subscriberType = subscriberEvaluationEventOPS
	case GoalCountEventOPSPersisterName:
		e.subscriberType = subscriberGoalEventOPS
	}
	return e, nil
}

func (e eventsOPSPersister) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	return nil
}

func (e eventsOPSPersister) Switch(ctx context.Context) (bool, error) {
	return false, nil
}
