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

package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
)

func TestInMemoryCachePutIfNewer(t *testing.T) {
	t.Parallel()
	const (
		key    = "k"
		genKey = "g"
	)
	gen := func(g int64) string {
		s, err := cache.FormatGeneration(g)
		require.NoError(t, err)
		return s
	}

	t.Run("negative generation", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		accepted, err := c.PutIfNewer(key, genKey, []byte("v"), -1, 0)
		assert.ErrorIs(t, err, cache.ErrInvalidGeneration)
		assert.False(t, accepted)
		_, err = c.Get(key)
		assert.Equal(t, cache.ErrNotFound, err)
	})
	t.Run("empty then older rejected then newer accepted", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		accepted, err := c.PutIfNewer(key, genKey, []byte("g10"), 10, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
		accepted, err = c.PutIfNewer(key, genKey, []byte("g5"), 5, 0)
		require.NoError(t, err)
		assert.False(t, accepted)
		v, err := c.Get(key)
		require.NoError(t, err)
		assert.Equal(t, []byte("g10"), v)
		accepted, err = c.PutIfNewer(key, genKey, []byte("g20"), 20, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
		v, err = c.Get(key)
		require.NoError(t, err)
		assert.Equal(t, []byte("g20"), v)
		m, err := c.Get(genKey)
		require.NoError(t, err)
		assert.Equal(t, cache.EncodeGenerationMarker(gen(20), []byte("g20")), m)
	})
	t.Run("unconditional overwrite invalidates marker", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		accepted, err := c.PutIfNewer(key, genKey, []byte("g100"), 100, 0)
		require.NoError(t, err)
		require.True(t, accepted)
		require.NoError(t, c.Put(key, []byte("legacy-write"), 0))
		accepted, err = c.PutIfNewer(key, genKey, []byte("g1"), 1, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
	})
	t.Run("same-length unconditional overwrite invalidates marker", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		accepted, err := c.PutIfNewer(key, genKey, []byte("abc"), 100, 0)
		require.NoError(t, err)
		require.True(t, accepted)
		require.NoError(t, c.Put(key, []byte("abd"), 0))
		accepted, err = c.PutIfNewer(key, genKey, []byte("g1"), 1, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
	})
	t.Run("non-byte value is treated as untracked", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		require.NoError(t, c.Put(key, "string-value", 0))
		require.NoError(t, c.Put(genKey, cache.EncodeGenerationMarker(gen(100), []byte("string-value")), 0))
		accepted, err := c.PutIfNewer(key, genKey, []byte("g1"), 1, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
	})
	t.Run("non-string marker is treated as missing", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		require.NoError(t, c.Put(key, []byte("v"), 0))
		require.NoError(t, c.Put(genKey, 42, 0))
		accepted, err := c.PutIfNewer(key, genKey, []byte("g1"), 1, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
	})
	t.Run("deleted value with leftover marker is accepted", func(t *testing.T) {
		t.Parallel()
		c := NewInMemoryCache()
		accepted, err := c.PutIfNewer(key, genKey, []byte("g100"), 100, 0)
		require.NoError(t, err)
		require.True(t, accepted)
		require.NoError(t, c.Delete(key))
		accepted, err = c.PutIfNewer(key, genKey, []byte("g1"), 1, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
	})
}
