// Copyright 2022 The Bucketeer Authors.
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

package v2

import (
	"fmt"

	goredis "github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	redis "github.com/bucketeer-io/bucketeer/pkg/redis/v2"
)

const (
	defaultKeysMaxSize = 3000000
)

type redisCache struct {
	cluster redis.Cluster
	logger  *zap.Logger
}

func NewRedisCache(cluster redis.Cluster, logger *zap.Logger) cache.Cache {
	return &redisCache{
		cluster: cluster,
		logger:  logger.Named("redis_cache_v2"),
	}
}

func NewRedisCacheLister(cluster redis.Cluster, logger *zap.Logger) cache.Lister {
	return &redisCache{
		cluster: cluster,
		logger:  logger.Named("redis_cache_lister_v2"),
	}
}

func NewRedisCacheDeleter(cluster redis.Cluster, logger *zap.Logger) cache.Deleter {
	return &redisCache{
		cluster: cluster,
		logger:  logger.Named("redis_cache_deleter_v2"),
	}
}

func (r *redisCache) Get(key interface{}) (interface{}, error) {
	value, err := r.cluster.Get(key.(string))
	if err != nil {
		if err == redis.ErrNil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	return value, nil
}

func (r *redisCache) Put(key interface{}, value interface{}) error {
	return r.cluster.Set(key.(string), value, 0)
}

func (r *redisCache) Delete(key string) error {
	return r.cluster.Del(key)
}

// If maxSize is less than or equals to zero, it is regarded as the default.
func (r *redisCache) Keys(pattern string, maxSize int) ([]string, error) {
	keys := []string{}
	if maxSize <= 0 {
		maxSize = defaultKeysMaxSize
	}
	fn := func(client *goredis.Client) error {
		r.logger.Debug(fmt.Sprintf("keys runs for %s", client.String()),
			zap.String("pattern", pattern),
		)
		iter := client.Scan(0, pattern, 0).Iterator()
		for iter.Next() {
			keys = append(keys, iter.Val())
			if len(keys) >= int(maxSize) {
				return nil
			}
		}
		if err := iter.Err(); err != nil {
			return err
		}
		return nil
	}
	err := r.cluster.ForEachMaster(fn)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
