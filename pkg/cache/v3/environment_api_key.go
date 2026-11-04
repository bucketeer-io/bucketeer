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

	"google.golang.org/protobuf/proto"

	accountdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

const (
	// environmentAPIKeyKind names entries addressed by the digest of the key. It is
	// distinct from legacyEnvironmentAPIKeyKind so that an entry written by this
	// version is never read back by a version that still addresses entries by the
	// key itself: that version would otherwise accept the digest as a key.
	environmentAPIKeyKind = "environment_apikey_hash"
	// legacyEnvironmentAPIKeyKind names the entries written before key names were
	// hashed. They have no expiry, so they are removed whenever the matching entry
	// is written or evicted; the periodic cacher writes every key, so a single run
	// clears them all after a rollout.
	legacyEnvironmentAPIKeyKind = "environment_apikey"
)

// EnvironmentAPIKeyCache is keyed by the API key presented by the client. The key
// is hashed before it is used as a cache key and the cached value only keeps its
// obfuscated form, so neither the key names nor the values hold the key itself.
type EnvironmentAPIKeyCache interface {
	Get(apiKey string) (*accountproto.EnvironmentAPIKey, error)
	Put(apiKey string, environmentAPIKey *accountproto.EnvironmentAPIKey) error
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

func (c *environmentAPIKeyCache) Put(apiKey string, environmentAPIKey *accountproto.EnvironmentAPIKey) error {
	// The caller's message is left untouched: it may be shared between caches.
	cached, ok := proto.Clone(environmentAPIKey).(*accountproto.EnvironmentAPIKey)
	if !ok {
		return cache.ErrInvalidType
	}
	if cached.ApiKey != nil {
		// Only the obfuscated form is cached. It is what the public API relays as
		// the editor of its changes, and it cannot be used to authenticate.
		cached.ApiKey.ApiKey = accountdomain.ObfuscateAPIKey(cached.ApiKey.ApiKey)
	}
	buffer, err := proto.Marshal(cached)
	if err != nil {
		return err
	}
	if err := c.cache.Put(c.key(apiKey), buffer, c.ttl); err != nil {
		return err
	}
	return evictKey(c.cache, c.legacyKey(apiKey))
}

func (c *environmentAPIKeyCache) Evict(apiKey string) error {
	if err := evictKey(c.cache, c.key(apiKey)); err != nil {
		return err
	}
	return evictKey(c.cache, c.legacyKey(apiKey))
}

func (c *environmentAPIKeyCache) key(apiKey string) string {
	// Pass "" as environmentId to skip the environment prefix in the cache key.
	// At lookup time we only have the API key string from the authorization header
	// and don't yet know which environment it belongs to.
	// The key is stored by its digest so the cache never holds the key itself.
	return cache.MakeKey(environmentAPIKeyKind, accountdomain.HashAPIKey(apiKey), "")
}

// legacyKey is the name under which previous versions stored the entry for apiKey.
func (c *environmentAPIKeyCache) legacyKey(apiKey string) string {
	return cache.MakeKey(legacyEnvironmentAPIKeyKind, apiKey, "")
}
