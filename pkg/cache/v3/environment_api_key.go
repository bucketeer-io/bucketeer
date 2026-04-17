// Copyright 2026 The Bucketeer Authors.
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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v3

import (
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

const (
	environmentAPIKeyKind = "environment_apikey"
)

type EnvironmentAPIKeyCache interface {
	Get(string) (*accountproto.EnvironmentAPIKey, error)
	Put(*accountproto.EnvironmentAPIKey) error
	Evict(apiKey string) error
}

type environmentAPIKeyCache struct {
	cache cache.Cache
	ttl   time.Duration
}

func NewEnvironmentAPIKeyCache(c cache.Cache, ttl time.Duration) EnvironmentAPIKeyCache {
	return &environmentAPIKeyCache{cache: c, ttl: ttl}
}

func (c *environmentAPIKeyCache) Get(apiKey string) (*accountproto.EnvironmentAPIKey, error) {
	key := c.key(apiKey)
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	environmentAPIKey := &accountproto.EnvironmentAPIKey{}
	if err := proto.Unmarshal(b, environmentAPIKey); err != nil {
		return nil, err
	}
	return environmentAPIKey, nil
}

func (c *environmentAPIKeyCache) Put(environmentAPIKey *accountproto.EnvironmentAPIKey) error {
	buffer, err := proto.Marshal(environmentAPIKey)
	if err != nil {
		return err
	}
	key := c.key(environmentAPIKey.ApiKey.ApiKey)
	return c.cache.Put(key, buffer, c.ttl)
}

func (c *environmentAPIKeyCache) Evict(apiKey string) error {
	return evictKey(c.cache, c.key(apiKey))
}

func (c *environmentAPIKeyCache) key(apiKey string) string {
	// Pass "" as environmentId to skip the environment prefix in the cache key.
	// At lookup time we only have the API key string from the authorization header
	// and don't yet know which environment it belongs to.
	return cache.MakeKey(environmentAPIKeyKind, apiKey, "")
}
