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

package apikey_last_used_at_writer

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	accountstotage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/proto/account"
)

type options struct {
	logger *zap.Logger
}

var defaultOptions = options{
	logger: zap.NewNop(),
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type apikeyLastUsedAt struct {
	apiKeyID      string
	environmentID string
	lastUsedAt    int64
}

type APIKeyLastUsedWriter struct {
	APIKeyLastUsedInfoCacher sync.Map
	mysqlClient              mysql.Client
	accountStorage           accountstotage.AccountStorage
	opts                     *options
	logger                   *zap.Logger
}

func NewAPIKeyLastUsedWriter(
	mysqlClient mysql.Client,
	opts ...Option,
) *APIKeyLastUsedWriter {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &APIKeyLastUsedWriter{
		mysqlClient:    mysqlClient,
		accountStorage: accountstotage.NewAccountStorage(mysqlClient),
		opts:           &options,
		logger:         options.logger.Named("api_grpc"),
	}
}

func (s *APIKeyLastUsedWriter) WriteAPIKeyLastUsedAt(ctx context.Context) {
	updatedAPIKeys := make([]string, 0)
	s.APIKeyLastUsedInfoCacher.Range(func(key, value interface{}) bool {
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
		s.APIKeyLastUsedInfoCacher.Delete(apiKeyID)
	}
}

func (s *APIKeyLastUsedWriter) CacheAPIKeyLastUsedAt(
	envAPIKey *account.EnvironmentAPIKey,
	lastUsedAt int64,
) {
	if cache, ok := s.APIKeyLastUsedInfoCacher.Load(envAPIKey.ApiKey.Id); ok {
		lastUsedAtCache := cache.(apikeyLastUsedAt)
		if lastUsedAtCache.lastUsedAt < lastUsedAt {
			s.APIKeyLastUsedInfoCacher.Store(envAPIKey.ApiKey.Id, apikeyLastUsedAt{
				apiKeyID:      envAPIKey.ApiKey.Id,
				lastUsedAt:    lastUsedAt,
				environmentID: envAPIKey.Environment.Id,
			})
		}
		return
	}
	s.APIKeyLastUsedInfoCacher.Store(envAPIKey.ApiKey.Id, apikeyLastUsedAt{
		apiKeyID:      envAPIKey.ApiKey.Id,
		lastUsedAt:    lastUsedAt,
		environmentID: envAPIKey.Environment.Id,
	})
}

func (s *APIKeyLastUsedWriter) WriteAPIKeyLastUsedAtCacheToDatabase(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("WriteAPIKeyLastUsedAtCacheToDatabase stopped")
			return
		case <-ticker.C:
			s.logger.Debug("writing API key last used at cache to database")
			s.WriteAPIKeyLastUsedAt(ctx)
		}
	}
}
