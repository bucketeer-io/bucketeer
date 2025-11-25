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

package api

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/proto/account"
)

type apikeyLastUsedAt struct {
	apiKeyID      string
	environmentID string
	lastUsedAt    int64
}

func (s *grpcGatewayService) cacheAPIKeyLastUsedAt(
	envAPIKey *account.EnvironmentAPIKey,
	lastUsedAt int64,
) {
	if cache, ok := s.apiKeyLastUsedInfoCacher.Load(envAPIKey.ApiKey.Id); ok {
		lastUsedAtCache := cache.(apikeyLastUsedAt)
		if lastUsedAtCache.lastUsedAt < lastUsedAt {
			s.apiKeyLastUsedInfoCacher.Store(envAPIKey.ApiKey.Id, apikeyLastUsedAt{
				apiKeyID:      envAPIKey.ApiKey.Id,
				lastUsedAt:    lastUsedAt,
				environmentID: envAPIKey.Environment.Id,
			})
		}
		return
	}
	s.apiKeyLastUsedInfoCacher.Store(envAPIKey.ApiKey.Id, apikeyLastUsedAt{
		apiKeyID:      envAPIKey.ApiKey.Id,
		lastUsedAt:    lastUsedAt,
		environmentID: envAPIKey.Environment.Id,
	})
}

func (s *grpcGatewayService) writeAPIKeyLastUsedAtCacheToDatabase(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("writeAPIKeyLastUsedAtCacheToDatabase stopped")
			return
		case <-ticker.C:
			s.logger.Debug("writing API key last used at cache to database")
			s.writeAPIKeyLastUsedAt(context.Background())
		}
	}
}

func (s *grpcGatewayService) writeAPIKeyLastUsedAt(ctx context.Context) {
	updatedAPIKeys := make([]string, 0)
	s.apiKeyLastUsedInfoCacher.Range(func(key, value interface{}) bool {
		apiKeyID := key.(string)
		lastUsedAtInfo := value.(apikeyLastUsedAt)

		_, err := s.accountStorage.UpdateAPIKeyLastUsedAt(
			ctx,
			apiKeyID,
			lastUsedAtInfo.environmentID,
			lastUsedAtInfo.lastUsedAt,
		)
		if err != nil {
			s.logger.Error(
				"Failed to update API Key Last Used At",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("apiKeyID", apiKeyID),
					zap.String("environmentID", lastUsedAtInfo.environmentID),
					zap.Int64("lastUsedAt", lastUsedAtInfo.lastUsedAt),
					zap.Error(err),
				)...,
			)
			return true
		}
		updatedAPIKeys = append(updatedAPIKeys, apiKeyID)
		return true
	})

	// Clear the cache for the updated API keys
	for _, apiKeyID := range updatedAPIKeys {
		s.apiKeyLastUsedInfoCacher.Delete(apiKeyID)
	}
}
