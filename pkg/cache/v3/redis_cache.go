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

package v3

import (
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	redis "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

type redisCache struct {
	client redis.Client
}

func NewRedisCache(client redis.Client) cache.MultiGetDeleteCache {
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

func (r *redisCache) Put(key interface{}, value interface{}) error {
	return r.client.Set(key.(string), value, 0)
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
