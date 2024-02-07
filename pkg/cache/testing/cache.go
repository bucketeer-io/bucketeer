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

package testing

import (
	"sync"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	redis "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

type inMemoryCache struct {
	data  map[interface{}]interface{}
	mutex sync.Mutex
}

func NewInMemoryCache() cache.MultiGetDeleteCountCache {
	return &inMemoryCache{
		data: make(map[interface{}]interface{}),
	}
}

func (c *inMemoryCache) Get(key interface{}) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if val, ok := c.data[key]; ok {
		return val, nil
	}
	return nil, cache.ErrNotFound
}

func (c *inMemoryCache) Put(key interface{}, value interface{}, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	return nil
}

func (c *inMemoryCache) GetMulti(keys interface{}) ([]interface{}, error) {
	// TODO: implement
	return nil, nil
}

func (c *inMemoryCache) Scan(cursor, key, count interface{}) (uint64, []string, error) {
	// TODO: implement
	return 0, nil, nil
}

func (c *inMemoryCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
	return nil
}

func (c *inMemoryCache) Increment(key string) (int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if val, ok := c.data[key]; ok {
		if intVal, ok := val.(int64); ok {
			intVal += 1
			c.data[key] = intVal
		}
	} else {
		c.data[key] = 1
	}
	return 0, nil
}

func (c *inMemoryCache) PFAdd(key string, els ...string) (int64, error) {
	// TODO: implement
	return 0, nil
}

func (c *inMemoryCache) PFCount(keys ...string) (int64, error) {
	// TODO: implement
	return 0, nil
}

func (c *inMemoryCache) PFMerge(dest string, keys ...string) error {
	// TODO: implement
	return nil
}

func (c *inMemoryCache) Expire(key string, expiration time.Duration) (bool, error) {
	// TODO: implement
	return true, nil
}

func (c *inMemoryCache) Pipeline() redis.PipeClient {
	// TODO: implement
	return nil
}
