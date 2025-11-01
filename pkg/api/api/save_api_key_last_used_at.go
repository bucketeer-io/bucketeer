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

	"github.com/bucketeer-io/bucketeer/v2/proto/account"
)

func (s *grpcGatewayService) cacheAPIKeyLastUsedAt(
	apiKey *account.APIKey,
	lastUsedAt int64,
) {
	if cache, ok := s.apiKeyLastUsedInfoCacher.Load(apiKey.Id); ok {
		lastUsedAtCache := cache.(int64)
		if lastUsedAtCache < lastUsedAt {
			s.apiKeyLastUsedInfoCacher.Store(apiKey.Id, lastUsedAt)
		}
		return
	}
	s.apiKeyLastUsedInfoCacher.Store(apiKey.Id, lastUsedAt)
}

func (s *grpcGatewayService) writeAPIKeyLastUsedAtCacheToDatabase(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.logger.Debug("writing API key last used at cache to database")
			s.writeAPIKeyLastUsedAt(ctx)
		}
	}
}

func (s *grpcGatewayService) writeAPIKeyLastUsedAt(ctx context.Context) {
	s.apiKeyLastUsedInfoCacher.Range(func(key, value interface{}) bool {
		apiKey := key.(string)
		lastUsedAt := value.(int64)

		_, err := s.accountClient.UpdateAPIKeyLastUsedAt(ctx, &account.UpdateAPIKeyLastUsedAtRequest{
			ApiKeyId:   apiKey,
			LastUsedAt: lastUsedAt,
		})
		if err != nil {
			s.logger.Error("failed to update API key last used at", zap.Error(err),
				zap.String("apiKeyId", apiKey),
				zap.Int64("lastUsedAt", lastUsedAt),
			)
			return true
		}
		return true
	})

	s.apiKeyLastUsedInfoCacher.Clear()
}
