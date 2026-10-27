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
	"errors"

	"go.uber.org/zap"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// loadLiveHoldoutGroup reads the serving snapshot from memory, then redis.
// A miss is "no live group". The gateway does not query MySQL.
func loadLiveHoldoutGroup(
	memory, redis cachev3.HoldoutGroupCache,
	environmentID string,
	logger *zap.Logger,
) *featureproto.HoldoutGroup {
	if memory != nil {
		group, err := memory.Get(environmentID)
		if err == nil {
			return evaluation.EvaluationHoldout(group)
		}
		if logger != nil && !errors.Is(err, cache.ErrNotFound) {
			logger.Debug("Failed to get the holdout group from memory",
				zap.Error(err),
				zap.String("environmentId", environmentID),
			)
		}
	}
	if redis == nil {
		return nil
	}
	group, err := redis.Get(environmentID)
	if err != nil {
		if logger != nil && !errors.Is(err, cache.ErrNotFound) {
			logger.Debug("Failed to get the holdout group from redis",
				zap.Error(err),
				zap.String("environmentId", environmentID),
			)
		}
		return nil
	}
	if memory != nil {
		if putErr := memory.Put(group, environmentID); putErr != nil && logger != nil {
			logger.Debug("Failed to store the holdout group in memory",
				zap.Error(putErr),
				zap.String("environmentId", environmentID),
			)
		}
	}
	return evaluation.EvaluationHoldout(group)
}

func holdoutGroups(group *featureproto.HoldoutGroup) []*featureproto.HoldoutGroup {
	if group == nil {
		return nil
	}
	return []*featureproto.HoldoutGroup{group}
}

func countHoldoutEvaluation(group *featureproto.HoldoutGroup, sourceID string) {
	if group == nil || len(group.Features) == 0 {
		return
	}
	holdoutEvaluationsCounter.WithLabelValues(sourceID).Inc()
}

func (s *gatewayService) liveHoldoutGroup(environmentID string) *featureproto.HoldoutGroup {
	return loadLiveHoldoutGroup(s.holdoutGroupCache, s.holdoutGroupRedisCache, environmentID, s.logger)
}

func (s *grpcGatewayService) liveHoldoutGroup(environmentID string) *featureproto.HoldoutGroup {
	return loadLiveHoldoutGroup(s.holdoutGroupCache, s.holdoutGroupRedisCache, environmentID, s.logger)
}
