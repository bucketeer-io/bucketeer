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

package api

import (
	"context"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"go.uber.org/zap"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	domaineventdomain "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type cacheInvalidator struct {
	featuresCache          cachev3.FeaturesCache
	segmentUsersCache      cachev3.SegmentUsersCache
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache
	logger                 *zap.Logger
}

func NewCacheInvalidator(
	featuresCache cachev3.FeaturesCache,
	segmentUsersCache cachev3.SegmentUsersCache,
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache,
	logger *zap.Logger,
) *cacheInvalidator {
	return &cacheInvalidator{
		featuresCache:          featuresCache,
		segmentUsersCache:      segmentUsersCache,
		environmentAPIKeyCache: environmentAPIKeyCache,
		logger:                 logger.Named("cache-invalidator"),
	}
}

func (ci *cacheInvalidator) Run(ctx context.Context, msgChan <-chan *puller.Message) error {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			ci.handleMessage(msg)
			msg.Ack()
		case <-ctx.Done():
			return nil
		}
	}
}

func (ci *cacheInvalidator) handleMessage(msg *puller.Message) {
	event := &domaineventproto.Event{}
	if err := proto.Unmarshal(msg.Data, event); err != nil {
		ci.logger.Warn("Failed to unmarshal domain event", zap.Error(err))
		return
	}
	switch event.EntityType {
	case domaineventproto.Event_FEATURE:
		if err := ci.featuresCache.Evict(event.EnvironmentId); err != nil {
			ci.logger.Warn("Failed to evict features cache",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("entityId", event.EntityId),
				zap.String("type", event.Type.String()),
			)
			return
		}
		cacheInvalidationCounter.WithLabelValues(
			event.EntityType.String(), event.Type.String(), event.EnvironmentId,
		).Inc()
		ci.logger.Debug("Evicted features cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	case domaineventproto.Event_SEGMENT:
		if err := ci.segmentUsersCache.Evict(event.EntityId, event.EnvironmentId); err != nil {
			ci.logger.Warn("Failed to evict segment users cache",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("segmentId", event.EntityId),
				zap.String("type", event.Type.String()),
			)
			return
		}
		cacheInvalidationCounter.WithLabelValues(
			event.EntityType.String(), event.Type.String(), event.EnvironmentId,
		).Inc()
		ci.logger.Debug("Evicted segment users cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("segmentId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	case domaineventproto.Event_APIKEY:
		secrets, err := domaineventdomain.ExtractAPIKeySecrets(event)
		if err != nil {
			if len(secrets) > 0 {
				ci.logger.Warn(
					"Partially failed to extract api_key from entity data; evicting available secrets",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
					zap.String("entityId", event.EntityId),
					zap.String("type", event.Type.String()),
				)
			} else {
				ci.logger.Error(
					"Failed to extract api_key from entity data",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
					zap.String("entityId", event.EntityId),
					zap.String("type", event.Type.String()),
				)
				return
			}
		}
		if len(secrets) == 0 {
			return
		}
		for _, s := range secrets {
			if err := ci.environmentAPIKeyCache.Evict(s); err != nil {
				ci.logger.Warn("Failed to evict environment API key cache",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
					zap.String("entityId", event.EntityId),
					zap.String("type", event.Type.String()),
				)
				return
			}
		}
		cacheInvalidationCounter.WithLabelValues(
			event.EntityType.String(), event.Type.String(), event.EnvironmentId,
		).Inc()
		ci.logger.Debug("Evicted environment API key cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	}
}
