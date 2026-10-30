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

package v3

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	accountdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

// rawAPIKey has the length of a generated key (32 random bytes, hex encoded).
const rawAPIKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func newEnvironmentAPIKey(apiKey string) *accountproto.EnvironmentAPIKey {
	return &accountproto.EnvironmentAPIKey{
		ApiKey: &accountproto.APIKey{
			Id:           "api-key-id",
			Name:         "api-key-name",
			Role:         accountproto.APIKey_SDK_CLIENT,
			ApiKey:       apiKey,
			ApiKeyPrefix: accountdomain.APIKeyPrefix(apiKey),
			Maintainer:   "test@bucketeer.io",
		},
		Environment: &environmentproto.EnvironmentV2{Id: "env-id"},
	}
}

func TestEnvironmentAPIKeyCacheKey(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		apiKey   string
		expected string
	}{
		{
			desc:     "generated key: the key name holds the digest, not the key",
			apiKey:   rawAPIKey,
			expected: fmt.Sprintf("%s:%s", environmentAPIKeyKind, accountdomain.HashAPIKey(rawAPIKey)),
		},
		{
			desc:     "empty key",
			apiKey:   "",
			expected: fmt.Sprintf("%s:%s", environmentAPIKeyKind, accountdomain.HashAPIKey("")),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			c := &environmentAPIKeyCache{}
			key := c.key(p.apiKey)
			assert.Equal(t, p.expected, key)
			if p.apiKey != "" {
				assert.NotContains(t, key, p.apiKey)
			}
		})
	}
}

func TestEnvironmentAPIKeyCacheGet(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	envAPIKey := newEnvironmentAPIKey(rawAPIKey)
	data, err := proto.Marshal(envAPIKey)
	require.NoError(t, err)
	key := fmt.Sprintf("%s:%s", environmentAPIKeyKind, accountdomain.HashAPIKey(rawAPIKey))

	patterns := []struct {
		desc        string
		setup       func(*environmentAPIKeyCache)
		expected    *accountproto.EnvironmentAPIKey
		expectedErr error
	}{
		{
			desc: "error: not found",
			setup: func(c *environmentAPIKeyCache) {
				c.cache.(*cachemock.MockCache).EXPECT().Get(key).Return(nil, cache.ErrNotFound)
			},
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "error: invalid type",
			setup: func(c *environmentAPIKeyCache) {
				c.cache.(*cachemock.MockCache).EXPECT().Get(key).Return("not bytes", nil)
			},
			expectedErr: cache.ErrInvalidType,
		},
		{
			desc: "success: looked up by the digest of the presented key",
			setup: func(c *environmentAPIKeyCache) {
				c.cache.(*cachemock.MockCache).EXPECT().Get(key).Return(data, nil)
			},
			expected: envAPIKey,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			c := &environmentAPIKeyCache{cache: cachemock.NewMockCache(mockController)}
			p.setup(c)
			got, err := c.Get(rawAPIKey)
			assert.Equal(t, p.expectedErr, err)
			if p.expected == nil {
				assert.Nil(t, got)
				return
			}
			assert.True(t, proto.Equal(p.expected, got))
		})
	}
}

func TestEnvironmentAPIKeyCachePut(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	key := fmt.Sprintf("%s:%s", environmentAPIKeyKind, accountdomain.HashAPIKey(rawAPIKey))
	ttl := time.Minute

	patterns := []struct {
		desc        string
		setup       func(*environmentAPIKeyCache)
		expectedErr error
	}{
		{
			desc: "error: backing cache",
			setup: func(c *environmentAPIKeyCache) {
				c.cache.(*cachemock.MockCache).EXPECT().
					Put(key, gomock.Any(), ttl).
					Return(cache.ErrNotFound)
			},
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "success: stored under the digest without the key itself",
			setup: func(c *environmentAPIKeyCache) {
				c.cache.(*cachemock.MockCache).EXPECT().
					Put(key, gomock.Any(), ttl).
					DoAndReturn(func(_ interface{}, value interface{}, _ time.Duration) error {
						data, ok := value.([]byte)
						require.True(t, ok)
						assert.False(t, bytes.Contains(data, []byte(rawAPIKey)), "the cached value must not hold the key")
						cached := &accountproto.EnvironmentAPIKey{}
						require.NoError(t, proto.Unmarshal(data, cached))
						assert.Equal(t, accountdomain.ObfuscateAPIKey(rawAPIKey), cached.ApiKey.ApiKey)
						assert.Equal(t, accountdomain.APIKeyPrefix(rawAPIKey), cached.ApiKey.ApiKeyPrefix)
						assert.Equal(t, "api-key-id", cached.ApiKey.Id)
						assert.Equal(t, "env-id", cached.Environment.Id)
						return nil
					})
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			c := &environmentAPIKeyCache{cache: cachemock.NewMockCache(mockController), ttl: ttl}
			p.setup(c)
			envAPIKey := newEnvironmentAPIKey(rawAPIKey)
			err := c.Put(rawAPIKey, envAPIKey)
			assert.Equal(t, p.expectedErr, err)
			// The caller's message is shared between caches and must be left intact.
			assert.Equal(t, rawAPIKey, envAPIKey.ApiKey.ApiKey)
		})
	}
}

func TestEnvironmentAPIKeyCachePutWithoutAPIKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	key := fmt.Sprintf("%s:%s", environmentAPIKeyKind, accountdomain.HashAPIKey(rawAPIKey))
	c := &environmentAPIKeyCache{cache: cachemock.NewMockCache(mockController)}
	c.cache.(*cachemock.MockCache).EXPECT().Put(key, gomock.Any(), time.Duration(0)).Return(nil)
	err := c.Put(rawAPIKey, &accountproto.EnvironmentAPIKey{
		Environment: &environmentproto.EnvironmentV2{Id: "env-id"},
	})
	assert.NoError(t, err)
}

func TestEnvironmentAPIKeyCacheInMemoryRoundTrip(t *testing.T) {
	t.Parallel()
	const otherAPIKey = "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"

	backing := NewInMemoryCache(WithEvictionInterval(0))
	c := NewEnvironmentAPIKeyCache(backing, 0)
	envAPIKey := newEnvironmentAPIKey(rawAPIKey)
	require.NoError(t, c.Put(rawAPIKey, envAPIKey))

	// The presented key is matched through its digest.
	got, err := c.Get(rawAPIKey)
	require.NoError(t, err)
	assert.Equal(t, envAPIKey.ApiKey.Id, got.ApiKey.Id)
	assert.Equal(t, envAPIKey.Environment.Id, got.Environment.Id)
	assert.Equal(t, accountdomain.ObfuscateAPIKey(rawAPIKey), got.ApiKey.ApiKey)
	assert.Equal(t, accountdomain.APIKeyPrefix(rawAPIKey), got.ApiKey.ApiKeyPrefix)

	// A different key does not resolve to the cached entry.
	_, err = c.Get(otherAPIKey)
	assert.Equal(t, cache.ErrNotFound, err)

	// Neither key names nor values hold the key itself.
	backing.entries.Range(func(k, v interface{}) bool {
		assert.NotContains(t, k.(string), rawAPIKey)
		e, ok := v.(*entry)
		require.True(t, ok)
		assert.False(t, bytes.Contains(e.value.([]byte), []byte(rawAPIKey)))
		assert.True(t, strings.HasPrefix(k.(string), environmentAPIKeyKind+":"))
		return true
	})

	// Eviction by the presented key removes the entry.
	require.NoError(t, c.Evict(rawAPIKey))
	_, err = c.Get(rawAPIKey)
	assert.Equal(t, cache.ErrNotFound, err)
}
