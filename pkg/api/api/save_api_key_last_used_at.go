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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/proto/account"
)

type lastUsedInfo struct {
	envAPIKey  *account.EnvironmentAPIKey
	lastUsedAt int64
}

type apikeyLastUsedInfoCache map[string]lastUsedInfo

type envAPIKeyLastUsedInfoCache map[string]apikeyLastUsedInfoCache

func (s *grpcGatewayService) cacheAPIKeyLastUsedAt(
	envAPIKey *account.EnvironmentAPIKey,
	lastUsedAt int64,
) {
	s.envAPIKeyLastUsedInfoMutex.Lock()
	defer s.envAPIKeyLastUsedInfoMutex.Unlock()

	if cache, ok := s.envAPIKeyLastUsedInfoCacher[envAPIKey.Environment.Id]; ok {
		if info, ok := cache[envAPIKey.ApiKey.Id]; ok {
			if info.lastUsedAt < lastUsedAt {
				info.lastUsedAt = lastUsedAt
			}
			s.envAPIKeyLastUsedInfoCacher[envAPIKey.Environment.Id][envAPIKey.ApiKey.Id] = info
			return
		}
		cache[envAPIKey.ApiKey.Id] = lastUsedInfo{
			envAPIKey:  envAPIKey,
			lastUsedAt: lastUsedAt,
		}
		s.envAPIKeyLastUsedInfoCacher[envAPIKey.Environment.Id] = cache
		return
	}
	cache := apikeyLastUsedInfoCache{}
	cache[envAPIKey.ApiKey.Id] = lastUsedInfo{
		envAPIKey:  envAPIKey,
		lastUsedAt: lastUsedAt,
	}
	s.envAPIKeyLastUsedInfoCacher[envAPIKey.Environment.Id] = cache
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
	s.envAPIKeyLastUsedInfoMutex.Lock()
	defer s.envAPIKeyLastUsedInfoMutex.Unlock()

	for _, cache := range s.envAPIKeyLastUsedInfoCacher {
		for _, info := range cache {
			envAPIKey, err := s.getEnvironmentAPIKey(ctx, info.envAPIKey.ApiKey.Id)
			if err != nil {
				s.logger.Error("failed to get environment API key", zap.Error(err),
					zap.String("apiKeyId", info.envAPIKey.ApiKey.Id),
				)
				continue
			}

			if envAPIKey == nil {
				s.logger.Error("environment API key not found",
					zap.String("apiKeyId", info.envAPIKey.ApiKey.Id),
				)
				continue
			}

			if envAPIKey.ApiKey.LastUsedAt >= info.lastUsedAt {
				continue
			}

			_, err = s.accountClient.UpdateAPIKey(ctx, &account.UpdateAPIKeyRequest{
				EnvironmentId: envAPIKey.Environment.Id,
				Id:            envAPIKey.ApiKey.Id,
				LastUsedAt:    wrapperspb.Int64(info.lastUsedAt),
			})
			if err != nil {
				s.logger.Error("failed to update API key last used at", zap.Error(err),
					zap.String("apiKeyId", info.envAPIKey.ApiKey.Id),
				)
				continue
			}
		}
	}
}
