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
	"errors"
	"net"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	domaineventdomain "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

var errCacheEvictionBadMessage = errors.New("cache eviction bad message")

type cacheEviction struct {
	featuresCache          cachev3.FeaturesCache
	segmentUsersCache      cachev3.SegmentUsersCache
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache
	experimentsCache       cachev3.ExperimentsCache
	autoOpsRulesCache      cachev3.AutoOpsRulesCache
	logger                 *zap.Logger
}

func NewCacheEviction(
	featuresCache cachev3.FeaturesCache,
	segmentUsersCache cachev3.SegmentUsersCache,
	environmentAPIKeyCache cachev3.EnvironmentAPIKeyCache,
	experimentsCache cachev3.ExperimentsCache,
	autoOpsRulesCache cachev3.AutoOpsRulesCache,
	logger *zap.Logger,
) subscriber.PubSubProcessor {
	return &cacheEviction{
		featuresCache:          featuresCache,
		segmentUsersCache:      segmentUsersCache,
		environmentAPIKeyCache: environmentAPIKeyCache,
		experimentsCache:       experimentsCache,
		autoOpsRulesCache:      autoOpsRulesCache,
		logger:                 logger,
	}
}

func (c *cacheEviction) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				c.logger.Error("cacheEviction: message channel closed")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberCacheEviction).Inc()
			c.handleMessage(msg)
		case <-ctx.Done():
			c.logger.Debug("cacheEviction: context done, stopped processing messages")
			return nil
		}
	}
}

func (c *cacheEviction) handleMessage(msg *puller.Message) {
	event := &domaineventproto.Event{}
	if err := proto.Unmarshal(msg.Data, event); err != nil {
		c.logger.Error("Failed to unmarshal domain event",
			zap.Error(err),
			zap.String("msgID", msg.ID),
		)
		subscriberHandledCounter.WithLabelValues(subscriberCacheEviction, codes.BadMessage.String()).Inc()
		msg.Ack()
		return
	}
	if err := c.evict(event); err != nil {
		if errors.Is(err, errCacheEvictionBadMessage) {
			subscriberHandledCounter.WithLabelValues(subscriberCacheEviction, codes.BadMessage.String()).Inc()
			msg.Ack()
			return
		}
		if isRepeatable(err) {
			c.logger.Warn("Failed to evict cache with repeatable error",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("entityId", event.EntityId),
				zap.String("entityType", event.EntityType.String()),
				zap.String("type", event.Type.String()),
			)
			subscriberHandledCounter.WithLabelValues(subscriberCacheEviction, codes.RepeatableError.String()).Inc()
			msg.Nack()
			return
		}
		c.logger.Error("Failed to evict cache with non-repeatable error",
			zap.Error(err),
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("entityType", event.EntityType.String()),
			zap.String("type", event.Type.String()),
		)
		subscriberHandledCounter.WithLabelValues(subscriberCacheEviction, codes.NonRepeatableError.String()).Inc()
		msg.Ack()
		return
	}
	subscriberHandledCounter.WithLabelValues(subscriberCacheEviction, codes.OK.String()).Inc()
	msg.Ack()
}

func (c *cacheEviction) evict(event *domaineventproto.Event) error {
	switch event.EntityType {
	case domaineventproto.Event_FEATURE:
		if err := c.featuresCache.Evict(event.EnvironmentId); err != nil {
			return err
		}
		c.logger.Debug("Evicted features redis cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	case domaineventproto.Event_SEGMENT:
		if err := c.segmentUsersCache.Evict(event.EntityId, event.EnvironmentId); err != nil {
			return err
		}
		c.logger.Debug("Evicted segment users redis cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("segmentId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	case domaineventproto.Event_APIKEY:
		secrets, err := domaineventdomain.ExtractAPIKeySecrets(event)
		if err != nil {
			if len(secrets) > 0 {
				c.logger.Warn(
					"Partially failed to extract api_key from entity data; evicting available secrets",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
					zap.String("entityId", event.EntityId),
					zap.String("type", event.Type.String()),
				)
			} else {
				c.logger.Error(
					"Failed to extract api_key from entity data",
					zap.Error(err),
					zap.String("environmentId", event.EnvironmentId),
					zap.String("entityId", event.EntityId),
					zap.String("type", event.Type.String()),
				)
				return errCacheEvictionBadMessage
			}
		}
		if len(secrets) == 0 {
			return nil
		}
		for _, s := range secrets {
			if err := c.environmentAPIKeyCache.Evict(s); err != nil {
				return err
			}
		}
		c.logger.Debug("Evicted environment API key redis cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	case domaineventproto.Event_EXPERIMENT:
		if err := c.experimentsCache.Evict(event.EnvironmentId); err != nil {
			return err
		}
		c.logger.Debug("Evicted experiments redis cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	case domaineventproto.Event_AUTOOPS_RULE:
		if err := c.autoOpsRulesCache.Evict(event.EnvironmentId); err != nil {
			return err
		}
		c.logger.Debug("Evicted auto ops rules redis cache",
			zap.String("environmentId", event.EnvironmentId),
			zap.String("entityId", event.EntityId),
			zap.String("type", event.Type.String()),
		)
	default:
		return nil
	}
	return nil
}

func isRepeatable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	if errors.Is(err, cache.ErrNotFound) || errors.Is(err, cache.ErrInvalidType) {
		return false
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "eof") {
		return true
	}
	return false
}
