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
	"sync"
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

// Expiry removal must only delete the exact entry that was observed to be
// expired. Deleting by key alone would race with a Put/PutIfNewer that stored
// a fresh entry in between, removing the new value (or only its generation
// marker, which would let a later stale write through).
func TestInMemoryCacheRemoveIfExpired(t *testing.T) {
	t.Parallel()
	now := time.Now()
	expired := &entry{value: "old", expiration: now.Add(-time.Second)}
	fresh := &entry{value: "new", expiration: now.Add(time.Hour)}
	noExpiry := &entry{value: "forever"}

	t.Run("not expired: untouched, returns false", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		c.entries.Store("k", fresh)
		assert.False(t, c.removeIfExpired("k", fresh, now))
		v, loaded := c.entries.Load("k")
		require.True(t, loaded)
		assert.Same(t, fresh, v)
	})
	t.Run("no expiry: untouched, returns false", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		c.entries.Store("k", noExpiry)
		assert.False(t, c.removeIfExpired("k", noExpiry, now))
		_, loaded := c.entries.Load("k")
		assert.True(t, loaded)
	})
	t.Run("exactly at expiration is not yet expired", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		e := &entry{value: "v", expiration: now}
		c.entries.Store("k", e)
		assert.False(t, c.removeIfExpired("k", e, now))
		_, loaded := c.entries.Load("k")
		assert.True(t, loaded)
	})
	t.Run("expired and still stored: deleted, returns true", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		c.entries.Store("k", expired)
		assert.True(t, c.removeIfExpired("k", expired, now))
		_, loaded := c.entries.Load("k")
		assert.False(t, loaded)
	})
	t.Run("expired but replaced by a fresh entry: fresh entry survives", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		// Simulates: reader loads `expired`, then a writer stores `fresh`, then
		// the reader attempts the expiry delete with its stale observation.
		c.entries.Store("k", expired)
		c.entries.Store("k", fresh)
		assert.True(t, c.removeIfExpired("k", expired, now))
		v, loaded := c.entries.Load("k")
		require.True(t, loaded)
		assert.Same(t, fresh, v)
	})
	t.Run("expired and already gone: no-op, returns true", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		assert.True(t, c.removeIfExpired("k", expired, now))
		_, loaded := c.entries.Load("k")
		assert.False(t, loaded)
	})
}

// Get must report an expired entry as missing even when the expiry delete was
// a no-op because a concurrent writer already replaced the entry; the caller
// retries and sees the new value on the next call.
func TestInMemoryCacheGetStaleObservationDoesNotDeleteFreshEntry(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache()
	expired := &entry{value: "old", expiration: time.Now().Add(-time.Second)}
	c.entries.Store("k", expired)
	// Writer wins the race before Get's delete.
	require.NoError(t, c.Put("k", "new", time.Hour))
	assert.True(t, c.removeIfExpired("k", expired, time.Now()))
	v, err := c.Get("k")
	require.NoError(t, err)
	assert.Equal(t, "new", v)
}

// The background sweep must not remove a value or marker that PutIfNewer
// rewrote after the sweep observed the previous (expired) pair.
func TestInMemoryCacheEvictExpiredKeepsFreshConditionalPair(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache()
	oldValue := &entry{value: []byte("old"), expiration: time.Now().Add(-time.Second)}
	oldMarker := &entry{value: "marker-old", expiration: time.Now().Add(-time.Second)}
	c.entries.Store("k", oldValue)
	c.entries.Store("g", oldMarker)

	// Sweep observes the expired pair, then PutIfNewer replaces both before the
	// sweep's deletes run.
	accepted, err := c.PutIfNewer("k", "g", []byte("new"), 10, time.Hour)
	require.NoError(t, err)
	require.True(t, accepted)
	assert.True(t, c.removeIfExpired("k", oldValue, time.Now()))
	assert.True(t, c.removeIfExpired("g", oldMarker, time.Now()))

	v, err := c.Get("k")
	require.NoError(t, err)
	assert.Equal(t, []byte("new"), v)
	_, err = c.Get("g")
	require.NoError(t, err)
	// Marker intact → an older generation is still rejected.
	accepted, err = c.PutIfNewer("k", "g", []byte("stale"), 5, time.Hour)
	require.NoError(t, err)
	assert.False(t, accepted)
}

func TestInMemoryCacheGetRejectsNonEntryValue(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache()
	c.entries.Store("k", "raw")
	_, err := c.Get("k")
	assert.Equal(t, cache.ErrInvalidType, err)
}

// evictExpired itself must route through the identity-checked delete.
func TestInMemoryCacheEvictExpiredOnlyRemovesExpiredEntries(t *testing.T) {
	t.Parallel()
	c := NewInMemoryCache()
	now := time.Now()
	c.entries.Store("expired", &entry{value: "a", expiration: now.Add(-time.Second)})
	c.entries.Store("fresh", &entry{value: "b", expiration: now.Add(time.Hour)})
	c.entries.Store("forever", &entry{value: "c"})
	c.entries.Store("not-an-entry", "raw")

	c.evictExpired(now)

	_, loaded := c.entries.Load("expired")
	assert.False(t, loaded)
	for _, k := range []string{"fresh", "forever", "not-an-entry"} {
		_, loaded := c.entries.Load(k)
		assert.True(t, loaded, k)
	}
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

func TestInMemoryCacheDeleteWithGeneration(t *testing.T) {
	t.Parallel()

	t.Run("removes value and marker", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		defer c.Destroy()
		accepted, err := c.PutIfNewer("k", "g", []byte("v"), 10, 0)
		require.NoError(t, err)
		require.True(t, accepted)
		require.NoError(t, c.DeleteWithGeneration("k", "g"))
		_, err = c.Get("k")
		assert.Equal(t, cache.ErrNotFound, err)
		_, err = c.Get("g")
		assert.Equal(t, cache.ErrNotFound, err)
	})
	t.Run("missing keys is not an error", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		defer c.Destroy()
		assert.NoError(t, c.DeleteWithGeneration("k", "g"))
	})
	t.Run("does not touch other keys", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		defer c.Destroy()
		require.NoError(t, c.Put("other", []byte("x"), 0))
		require.NoError(t, c.DeleteWithGeneration("k", "g"))
		v, err := c.Get("other")
		require.NoError(t, err)
		assert.Equal(t, []byte("x"), v)
	})
	// A put racing with the delete must end up either fully present (value and
	// marker) or fully absent; a value without a marker would let a later stale
	// write through.
	t.Run("concurrent put is never left unmarked", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		defer c.Destroy()
		for i := 0; i < 200; i++ {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				_, err := c.PutIfNewer("k", "g", []byte("new"), 200, 0)
				assert.NoError(t, err)
			}()
			go func() {
				defer wg.Done()
				assert.NoError(t, c.DeleteWithGeneration("k", "g"))
			}()
			wg.Wait()
			_, valueErr := c.Get("k")
			_, markerErr := c.Get("g")
			assert.Equal(t, valueErr == nil, markerErr == nil, "iteration %d: value/marker presence diverged", i)
			if valueErr == nil {
				accepted, err := c.PutIfNewer("k", "g", []byte("old"), 100, 0)
				require.NoError(t, err)
				assert.False(t, accepted, "iteration %d: stale write accepted", i)
			}
			require.NoError(t, c.DeleteWithGeneration("k", "g"))
		}
	})
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
