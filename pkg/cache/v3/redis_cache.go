// Copyright 2024 The Bucketeer Authors.
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

package v3

import (
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	redis "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

type redisCache struct {
	client redis.Client
}

func NewRedisCache(client redis.Client) cache.MultiGetDeleteCountCache {
	return &redisCache{
		client: client,
	}
}

func (r *redisCache) Get(key interface{}) (interface{}, error) {
	value, err := r.client.Get(key.(string))
	if err != nil {
		if err == redis.ErrNil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}
	return value, nil
}

func (r *redisCache) Put(key interface{}, value interface{}, expiration time.Duration) error {
	return r.client.Set(key.(string), value, expiration)
}

func (r *redisCache) GetMulti(keys interface{}) ([]interface{}, error) {
	value, err := r.client.GetMulti(keys.([]string))
	switch err {
	case nil:
		return value, nil
	case redis.ErrNil:
		return nil, cache.ErrNotFound
	case redis.ErrInvalidType:
		return nil, cache.ErrInvalidType
	default:
		return nil, err
	}
}

func (r *redisCache) Scan(cursor, key, count interface{}) (uint64, []string, error) {
	c, keys, err := r.client.Scan(cursor.(uint64), key.(string), count.(int64))
	switch err {
	case nil:
		return c, keys, nil
	case redis.ErrNil:
		return 0, nil, cache.ErrNotFound
	default:
		return 0, nil, err
	}
}

func (r *redisCache) Delete(key string) error {
	return r.client.Del(key)
}

func (r *redisCache) Increment(key string) (int64, error) {
	return r.client.Incr(key)
}

func (r *redisCache) PFCount(keys ...string) (int64, error) {
	return r.client.PFCount(keys...)
}

func (r *redisCache) PFMerge(dest string, keys ...string) error {
	return r.client.PFMerge(dest, keys...)
}

func (r *redisCache) PFAdd(key string, els ...string) (int64, error) {
	return r.client.PFAdd(key, els...)
}

func (r *redisCache) Pipeline() redis.PipeClient {
	return r.client.Pipeline()
}

func (r *redisCache) Expire(key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(key, expiration)
}
