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

package rediscounter

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ec "github.com/bucketeer-io/bucketeer/pkg/eventcounter/api"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

const (
	redisScanMaxSize  = int64(100)
	redisChunkMaxSize = 100
	day               = 24 * 60 * 60
)

var (
	kinds = []string{
		ec.UserCountPrefix,
		ec.EventCountPrefix,
	}

	errSubmatchStringNotFound = errors.New("batch: submatch string not found")
	errParseInt               = errors.New("batch: failed to parse int from string")
)

type redisCounterDeleter struct {
	envClient envclient.Client
	cache     cache.MultiGetDeleteCountCache
	opts      *jobs.Options
	logger    *zap.Logger
}

func NewRedisCounterDeleter(
	redis cache.MultiGetDeleteCountCache,
	environmentClient envclient.Client,
	opts ...jobs.Option) jobs.Job {

	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &redisCounterDeleter{
		cache:     redis,
		envClient: environmentClient,
		opts:      dopts,
		logger:    dopts.Logger.Named("redis-counter-deleter"),
	}
}

func (r *redisCounterDeleter) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobRedisCounterDeleter, lastErr, time.Since(startTime))
	}()

	ctx, cancel := context.WithTimeout(ctx, r.opts.Timeout)
	defer cancel()

	r.logger.Info("Starting to delete old counters from Redis")
	envs, err := r.listEnvironments(ctx)
	if err != nil {
		return err
	}
	deletedKeys := 0
	for _, env := range envs {
		for _, kind := range kinds {
			keysSize, err := r.deleteKeysByKind(env.Id, kind)
			if err != nil {
				return err
			}
			deletedKeys += keysSize
		}
	}
	r.logger.Info("Finished deleting old counters from Redis",
		zap.Duration("elapsedTime", time.Since(startTime)),
		zap.Int("keysDeletedSize", deletedKeys),
	)
	return nil
}

func (r *redisCounterDeleter) listEnvironments(ctx context.Context) ([]*envproto.EnvironmentV2, error) {
	resp, err := r.envClient.ListEnvironmentsV2(ctx, &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
		Archived: &wrapperspb.BoolValue{Value: false},
	})
	if err != nil {
		r.logger.Error("Failed to list environments", zap.Error(err))
		return nil, err
	}
	return resp.Environments, nil
}

func (r *redisCounterDeleter) deleteKeysByKind(environmentId, kind string) (int, error) {
	keyPrefix := r.newKeyPrefix(environmentId, kind)
	keys, err := r.scan(environmentId, kind, keyPrefix)
	if err != nil {
		r.logger.Error("Failed to scan keys from redis",
			zap.String("environmentId", environmentId),
			zap.String("kind", kind),
			zap.String("keyPrefix", keyPrefix),
		)
		return 0, err
	}
	if len(keys) == 0 {
		r.logger.Info("No keys was found",
			zap.String("environmentId", environmentId),
			zap.String("kind", kind),
		)
		return 0, nil
	}
	filteredKeys, err := r.filterKeysOlderThanThirtyOneDays(environmentId, kind, keys)
	if err != nil {
		return 0, err
	}
	r.logger.Info("Filtered keys older than 31 days",
		zap.String("environmentId", environmentId),
		zap.String("kind", kind),
		zap.Int("filteredKeysSize", len(filteredKeys)),
	)
	// To avoid blocking Redis for too much time while deleting all the keys
	// we split the keys in chunks
	chunks := r.chunkSlice(filteredKeys, redisChunkMaxSize)
	r.logger.Info("Chunked the filtered keys", zap.Int("chunkSize", len(chunks)))
	deletedKeys := 0
	for _, chunk := range chunks {
		if err := r.deleteKeys(chunk); err != nil {
			r.logger.Error("Failed to delete chunk from redis",
				zap.String("environmentId", environmentId),
				zap.String("kind", kind),
				zap.Strings("keys", chunk),
				zap.Int("deletedKeysSizeUntilTheError", deletedKeys),
			)
			// Return the number of deleted keys until the error
			return deletedKeys, err
		}
		deletedKeys += len(chunk)
		r.logger.Info("Chunk deleted successfully", zap.Strings("keys", chunk))
	}
	return deletedKeys, nil
}

func (r *redisCounterDeleter) newKeyPrefix(environmentId, kind string) string {
	keyPrefix := cache.MakeKeyPrefix(kind, environmentId)
	return keyPrefix + "*"
}

func (r *redisCounterDeleter) scan(environmentId, kind, key string) ([]string, error) {
	r.logger.Info("Starting scan keys from Redis",
		zap.String("environmentId", environmentId),
		zap.String("kind", kind),
	)
	startTime := time.Now()
	var cursor uint64
	var k []string
	var err error
	keys := []string{}
	for {
		cursor, k, err = r.cache.Scan(cursor, key, redisScanMaxSize)
		if err != nil {
			break
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	r.logger.Info("Finished scanning keys from Redis",
		zap.String("environmentId", environmentId),
		zap.String("kind", kind),
		zap.Duration("elapsedTime", time.Since(startTime)),
		zap.Int("keysSize", len(keys)),
	)
	return keys, nil
}

func (r *redisCounterDeleter) filterKeysOlderThanThirtyOneDays(
	environmentId, kind string,
	keys []string,
) ([]string, error) {
	filteredKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		// E.g. environment_id:uc:1689835532:feature_id:variation_id
		var regex string
		if environmentId == "" {
			regex = fmt.Sprintf("%s:([0-9]+):", kind)
		} else {
			regex = fmt.Sprintf("%s:%s:([0-9]+):", environmentId, kind)
		}
		re := regexp.MustCompile(regex)
		match := re.FindStringSubmatch(key)
		if len(match) == 0 {
			r.logger.Error("Failed to find submatch string",
				zap.String("environmentId", environmentId),
				zap.String("kind", kind),
				zap.String("key", key),
				zap.String("regex", regex),
			)
			return nil, errSubmatchStringNotFound
		}
		createdAt, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			r.logger.Error("Failed to convert string to int",
				zap.String("created_at", match[1]),
				zap.String("environmentId", environmentId),
				zap.String("kind", kind),
			)
			return nil, errParseInt
		}
		now := time.Now()
		if now.Unix()-createdAt < 30*day {
			continue
		}
		filteredKeys = append(filteredKeys, key)
	}
	return filteredKeys, nil
}

func (r *redisCounterDeleter) chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func (r *redisCounterDeleter) deleteKeys(keys []string) error {
	for _, key := range keys {
		err := r.cache.Delete(key)
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return fmt.Errorf("failed to delete key %s: %w", key, err)
		}
	}
	return nil
}
