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

package cache

import (
	"github.com/bucketeer-io/bucketeer/pkg/redis"
)

type redisCache struct {
	cluster redis.Cluster
}

func NewRedisCache(cluster redis.Cluster) Cache {
	return &redisCache{
		cluster: cluster,
	}
}

func (r *redisCache) Get(key interface{}) (interface{}, error) {
	conn := r.cluster.Get(redis.WithReadOnly())
	defer conn.Close()
	value, err := conn.Do("GET", key)
	if err != nil {
		if err == redis.ErrNil {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return value, nil
}

func (r *redisCache) Put(key interface{}, value interface{}) error {
	conn := r.cluster.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}
