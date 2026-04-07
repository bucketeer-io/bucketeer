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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
)

func TestInMemoryCacheGet(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func(t *testing.T, c *InMemoryCache)
		key         string
		expected    interface{}
		expectedErr error
	}{
		{
			desc:        "not found: key does not exist",
			setup:       func(t *testing.T, c *InMemoryCache) {},
			key:         "missing-key",
			expectedErr: cache.ErrNotFound,
		},
		{
			desc: "success: entry within TTL",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put("key1", "value1", 10*time.Minute))
			},
			key:      "key1",
			expected: "value1",
		},
		{
			desc: "success: entry with no expiry (TTL=0)",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put("key1", "value1", 0))
			},
			key:      "key1",
			expected: "value1",
		},
		{
			desc: "not found: entry expired",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put("key1", "value1", 50*time.Millisecond))
				time.Sleep(100 * time.Millisecond)
			},
			key:         "key1",
			expectedErr: cache.ErrNotFound,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			c := NewInMemoryCache()
			p.setup(t, c)
			val, err := c.Get(p.key)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, p.expected, val)
		})
	}
}

func TestInMemoryCacheGetDeletesExpiredEntry(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache()
	require.NoError(t, c.Put("key1", "value1", 50*time.Millisecond))

	time.Sleep(100 * time.Millisecond)

	_, err := c.Get("key1")
	assert.Equal(t, cache.ErrNotFound, err)

	_, loaded := c.entries.Load("key1")
	assert.False(t, loaded)
}

func TestInMemoryCacheEvicterCleansExpiredEntries(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache(WithEvictionInterval(50 * time.Millisecond))
	defer c.Destroy()

	require.NoError(t, c.Put("key1", "value1", 50*time.Millisecond))

	assert.Eventually(t, func() bool {
		_, loaded := c.entries.Load("key1")
		return !loaded
	}, 1*time.Second, 50*time.Millisecond)
}

func TestInMemoryCacheEvicterSkipsNoExpiryEntries(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache(WithEvictionInterval(50 * time.Millisecond))
	defer c.Destroy()

	require.NoError(t, c.Put("no-expiry", "value1", 0))
	require.NoError(t, c.Put("with-expiry", "value2", 50*time.Millisecond))

	assert.Eventually(t, func() bool {
		_, loaded := c.entries.Load("with-expiry")
		return !loaded
	}, 1*time.Second, 50*time.Millisecond)

	val, err := c.Get("no-expiry")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestInMemoryCacheDestroy(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache(WithEvictionInterval(50 * time.Millisecond))

	require.NoError(t, c.Put("key1", "value1", 5*time.Minute))
	require.NoError(t, c.Put("key2", "value2", 5*time.Minute))

	c.Destroy()

	_, err := c.Get("key1")
	assert.Equal(t, cache.ErrNotFound, err)
	_, err = c.Get("key2")
	assert.Equal(t, cache.ErrNotFound, err)
}
