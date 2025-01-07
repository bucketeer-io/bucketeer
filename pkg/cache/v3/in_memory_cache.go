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

package v3

import (
	"sync"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
)

const (
	defaultEvictionInterval = 0
)

type entry struct {
	value      interface{}
	expiration time.Time
}

type InMemoryCache struct {
	entries          sync.Map
	evictionInterval time.Duration
	doneCh           chan struct{}
}

type InMemoryCacheOption func(*InMemoryCache)

func WithEvictionInterval(evictionInterval time.Duration) InMemoryCacheOption {
	return func(c *InMemoryCache) {
		c.evictionInterval = evictionInterval
	}
}

func NewInMemoryCache(opts ...InMemoryCacheOption) *InMemoryCache {
	c := &InMemoryCache{
		evictionInterval: defaultEvictionInterval,
		doneCh:           make(chan struct{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.evictionInterval > 0 {
		go c.startEvicter(c.evictionInterval)
	}
	return c
}

func (c *InMemoryCache) startEvicter(evictionInterval time.Duration) {
	ticker := time.NewTicker(evictionInterval)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			c.evictExpired(now)
		case <-c.doneCh:
			return
		}
	}
}

func (c *InMemoryCache) evictExpired(t time.Time) {
	c.entries.Range(func(key, value interface{}) bool {
		if e, ok := value.(*entry); ok && e.expiration.Before(t) {
			c.entries.Delete(key)
		}
		return true
	})
}

func (c *InMemoryCache) Get(key interface{}) (interface{}, error) {
	value, ok := c.entries.Load(key)
	if !ok {
		return nil, cache.ErrNotFound
	}
	e, ok := value.(*entry)
	if !ok {
		return nil, cache.ErrInvalidType
	}
	return e.value, nil
}

func (c *InMemoryCache) Put(key, value interface{}, expiration time.Duration) error {
	c.entries.Store(key, &entry{
		value:      value,
		expiration: time.Now().Add(expiration),
	})
	return nil
}

func (c *InMemoryCache) Destroy() {
	close(c.doneCh)
	c.entries.Range(func(key, value interface{}) bool {
		c.entries.Delete(key)
		return true
	})
}
