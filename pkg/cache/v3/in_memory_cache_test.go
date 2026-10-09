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

func TestInMemoryCachePutIfNewer(t *testing.T) {
	t.Parallel()
	const (
		key    = "env:features"
		genKey = "{env:features}:gen"
	)
	gen := func(g int64) string {
		s, err := cache.FormatGeneration(g)
		require.NoError(t, err)
		return s
	}
	patterns := []struct {
		desc             string
		setup            func(t *testing.T, c *InMemoryCache)
		value            []byte
		generation       int64
		expiration       time.Duration
		expectedAccepted bool
		expectedErr      error
		expectedValue    interface{}
		expectedMarker   interface{}
	}{
		{
			desc:             "empty cache: accepted",
			setup:            func(*testing.T, *InMemoryCache) {},
			value:            []byte("v1"),
			generation:       10,
			expectedAccepted: true,
			expectedValue:    []byte("v1"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(10), []byte("v1")),
		},
		{
			desc: "newer over older: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("old"), 0))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(10), []byte("old")), 0))
			},
			value:            []byte("new!"),
			generation:       20,
			expectedAccepted: true,
			expectedValue:    []byte("new!"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(20), []byte("new!")),
		},
		{
			desc: "older over newer: rejected",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("new!"), 0))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(20), []byte("new!")), 0))
			},
			value:            []byte("old"),
			generation:       10,
			expectedAccepted: false,
			expectedValue:    []byte("new!"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(20), []byte("new!")),
		},
		{
			desc: "equal generation: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("same"), 0))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(10), []byte("same")), 0))
			},
			value:            []byte("same2"),
			generation:       10,
			expectedAccepted: true,
			expectedValue:    []byte("same2"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(10), []byte("same2")),
		},
		{
			desc: "value missing, stale marker left behind: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(999), []byte("x")), 0))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "value expired, stale marker not expired: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("expired"), time.Nanosecond))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(999), []byte("expired")), 0))
				time.Sleep(2 * time.Millisecond)
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "value present, no marker: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("legacy"), 0))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "newer marker but value overwritten unconditionally (length mismatch): accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("legacy-overwrite"), 0))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(999), []byte("xxxx")), 0))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "newer marker, same-length different content (legacy overwrite): accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("abd"), 0))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(999), []byte("abc")), 0))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "newer marker but value is not bytes: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, "a-string", 0))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(999), []byte("a-string")), 0))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "marker is not a string: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("val"), 0))
				require.NoError(t, c.Put(genKey, 12345, 0))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "malformed marker: accepted",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte("val"), 0))
				require.NoError(t, c.Put(genKey, "garbage", 0))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    []byte("v"),
			expectedMarker:   cache.EncodeGenerationMarker(gen(1), []byte("v")),
		},
		{
			desc: "empty value with newer marker: rejected",
			setup: func(t *testing.T, c *InMemoryCache) {
				require.NoError(t, c.Put(key, []byte{}, 0))
				require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(20), nil), 0))
			},
			value:            []byte("stale"),
			generation:       10,
			expectedAccepted: false,
			expectedValue:    []byte{},
			expectedMarker:   cache.EncodeGenerationMarker(gen(20), nil),
		},
		{
			desc:             "negative generation: error, nothing written",
			setup:            func(*testing.T, *InMemoryCache) {},
			value:            []byte("v"),
			generation:       -1,
			expectedAccepted: false,
			expectedErr:      cache.ErrInvalidGeneration,
			expectedValue:    nil,
			expectedMarker:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			c := NewInMemoryCache()
			defer c.Destroy()
			p.setup(t, c)

			accepted, err := c.PutIfNewer(key, genKey, p.value, p.generation, p.expiration)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedAccepted, accepted)

			gotValue, _ := c.Get(key)
			assert.Equal(t, p.expectedValue, gotValue)
			gotMarker, _ := c.Get(genKey)
			assert.Equal(t, p.expectedMarker, gotMarker)
		})
	}
}

func TestInMemoryCachePutIfNewerAppliesExpiration(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache()
	defer c.Destroy()
	accepted, err := c.PutIfNewer("k", "g", []byte("v"), 1, 10*time.Millisecond)
	require.NoError(t, err)
	require.True(t, accepted)
	for _, k := range []string{"k", "g"} {
		v, ok := c.entries.Load(k)
		require.True(t, ok)
		assert.False(t, v.(*entry).expiration.IsZero(), "%s should carry an expiration", k)
	}
	accepted, err = c.PutIfNewer("k2", "g2", []byte("v"), 1, 0)
	require.NoError(t, err)
	require.True(t, accepted)
	for _, k := range []string{"k2", "g2"} {
		v, ok := c.entries.Load(k)
		require.True(t, ok)
		assert.True(t, v.(*entry).expiration.IsZero(), "%s should not carry an expiration", k)
	}
}

func TestInMemoryCachePutIfNewerConcurrentHighestWins(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache()
	defer c.Destroy()
	const writers = 64
	done := make(chan struct{})
	for i := 1; i <= writers; i++ {
		go func(g int64) {
			defer func() { done <- struct{}{} }()
			_, err := c.PutIfNewer("k", "g", []byte{byte(g)}, g, 0)
			assert.NoError(t, err)
		}(int64(i))
	}
	for i := 0; i < writers; i++ {
		<-done
	}
	v, err := c.Get("k")
	require.NoError(t, err)
	assert.Equal(t, []byte{writers}, v)
	formatted, err := cache.FormatGeneration(writers)
	require.NoError(t, err)
	m, err := c.Get("g")
	require.NoError(t, err)
	assert.Equal(t, cache.EncodeGenerationMarker(formatted, []byte{writers}), m)
}
