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
	"context"
	"errors"
	"math"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	redis "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
)

const (
	testValueKey = "env-1:features"
	testGenKey   = "{env-1:features}:gen"
)

// newMiniRedisCache runs the real Lua script against an in-process Redis so
// the script's semantics (not a mock of them) are what is tested.
func newMiniRedisCache(t *testing.T) (*miniredis.Miniredis, *redisCache) {
	t.Helper()
	mr := miniredis.RunT(t)
	client, err := redis.NewClient(
		mr.Addr(),
		redis.WithRedisMode(redis.RedisModeStandalone),
		redis.WithLogger(zap.NewNop()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	return mr, NewRedisCache(client).(*redisCache)
}

func formatGen(t *testing.T, g int64) string {
	t.Helper()
	s, err := cache.FormatGeneration(g)
	require.NoError(t, err)
	return s
}

// mk builds the marker a conditional writer would store for value at
// generation g.
func mk(t *testing.T, g int64, value string) string {
	t.Helper()
	return cache.EncodeGenerationMarker(formatGen(t, g), []byte(value))
}

func TestRedisCachePutIfNewer(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc             string
		setup            func(t *testing.T, mr *miniredis.Miniredis)
		value            []byte
		generation       int64
		expiration       time.Duration
		expectedAccepted bool
		expectedValue    string
		expectedMarker   string
		expectedTTL      time.Duration
	}{
		{
			desc:             "empty cache: accepted and both keys written without TTL",
			setup:            func(*testing.T, *miniredis.Miniredis) {},
			value:            []byte("v1"),
			generation:       10,
			expiration:       0,
			expectedAccepted: true,
			expectedValue:    "v1",
			expectedMarker:   mk(t, 10, "v1"),
			expectedTTL:      0,
		},
		{
			desc:             "empty cache with TTL: both keys get the same TTL",
			setup:            func(*testing.T, *miniredis.Miniredis) {},
			value:            []byte("v1"),
			generation:       10,
			expiration:       90 * time.Second,
			expectedAccepted: true,
			expectedValue:    "v1",
			expectedMarker:   mk(t, 10, "v1"),
			expectedTTL:      90 * time.Second,
		},
		{
			desc: "newer generation over older: accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "old"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 10, "old")))
			},
			value:            []byte("new!"),
			generation:       20,
			expectedAccepted: true,
			expectedValue:    "new!",
			expectedMarker:   mk(t, 20, "new!"),
		},
		{
			desc: "older generation over newer: rejected, nothing changes",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "new!"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 20, "new!")))
			},
			value:            []byte("old"),
			generation:       10,
			expectedAccepted: false,
			expectedValue:    "new!",
			expectedMarker:   mk(t, 20, "new!"),
		},
		{
			desc: "equal generation: accepted (idempotent rewrite)",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "same"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 10, "same")))
			},
			value:            []byte("same2"),
			generation:       10,
			expectedAccepted: true,
			expectedValue:    "same2",
			expectedMarker:   mk(t, 10, "same2"),
		},
		{
			desc: "value missing but stale newer marker left behind: accepted (repopulation never blocked)",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testGenKey, mk(t, 999, "xxx")))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "value present but no marker (written by legacy pod): accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "legacy"))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "newer marker but value overwritten by legacy pod (length mismatch): accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "legacy-overwrite"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 999, "xxxx")))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "newer marker, value overwritten by legacy pod with SAME length: accepted (digest mismatch)",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "abd"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 999, "abc")))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "newer marker, length matches but marker digest is for a different value: accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "abc"))
				require.NoError(t, mr.Set(testGenKey, formatGen(t, 999)+":3:"+cache.ValueDigest([]byte("xyz"))))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "marker with upper-case digest is malformed: accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "abc"))
				require.NoError(t, mr.Set(testGenKey, formatGen(t, 999)+":3:"+strings.ToUpper(cache.ValueDigest([]byte("abc")))))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "binary (non-UTF8) value round-trips through digest check: older rejected",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "\x00\xff\x10\x80"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 20, "\x00\xff\x10\x80")))
			},
			value:            []byte("old"),
			generation:       10,
			expectedAccepted: false,
			expectedValue:    "\x00\xff\x10\x80",
			expectedMarker:   mk(t, 20, "\x00\xff\x10\x80"),
		},
		{
			desc: "malformed marker (garbage): accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "val"))
				require.NoError(t, mr.Set(testGenKey, "garbage"))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "malformed marker (wrong generation width): accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "val"))
				require.NoError(t, mr.Set(testGenKey, "999:3:"+cache.ValueDigest([]byte("val"))))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "malformed marker (non-numeric length): accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "val"))
				require.NoError(t, mr.Set(testGenKey, formatGen(t, 999)+":abc:"+cache.ValueDigest([]byte("val"))))
			},
			value:            []byte("v"),
			generation:       1,
			expectedAccepted: true,
			expectedValue:    "v",
			expectedMarker:   mk(t, 1, "v"),
		},
		{
			desc: "empty value cached with newer marker: rejected (EXISTS, not STRLEN, decides presence)",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, ""))
				require.NoError(t, mr.Set(testGenKey, mk(t, 20, "")))
			},
			value:            []byte("stale-with-content"),
			generation:       10,
			expectedAccepted: false,
			expectedValue:    "",
			expectedMarker:   mk(t, 20, ""),
		},
		{
			desc: "empty value accepted when newer",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "content"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 10, "content")))
			},
			value:            []byte{},
			generation:       20,
			expectedAccepted: true,
			expectedValue:    "",
			expectedMarker:   mk(t, 20, ""),
		},
		{
			desc: "unix-nano magnitude generations compare exactly (beyond double precision)",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "new"))
				require.NoError(t, mr.Set(testGenKey,
					mk(t, 1_791_600_000_000_000_001, "new")))
			},
			value:            []byte("old"),
			generation:       1_791_600_000_000_000_000,
			expectedAccepted: false,
			expectedValue:    "new",
			expectedMarker:   mk(t, 1_791_600_000_000_000_001, "new"),
		},
		{
			desc: "one nanosecond newer at unix-nano magnitude: accepted",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "old"))
				require.NoError(t, mr.Set(testGenKey,
					mk(t, 1_791_600_000_000_000_000, "old")))
			},
			value:            []byte("new"),
			generation:       1_791_600_000_000_000_001,
			expectedAccepted: true,
			expectedValue:    "new",
			expectedMarker:   mk(t, 1_791_600_000_000_000_001, "new"),
		},
		{
			desc: "rejected write does not touch an existing TTL",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "new!"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 20, "new!")))
				mr.SetTTL(testValueKey, 30*time.Second)
				mr.SetTTL(testGenKey, 30*time.Second)
			},
			value:            []byte("old"),
			generation:       10,
			expiration:       5 * time.Minute,
			expectedAccepted: false,
			expectedValue:    "new!",
			expectedMarker:   mk(t, 20, "new!"),
			expectedTTL:      30 * time.Second,
		},
		{
			desc: "accepted write with zero expiration clears a previous TTL",
			setup: func(t *testing.T, mr *miniredis.Miniredis) {
				require.NoError(t, mr.Set(testValueKey, "old"))
				require.NoError(t, mr.Set(testGenKey, mk(t, 10, "old")))
				mr.SetTTL(testValueKey, 30*time.Second)
				mr.SetTTL(testGenKey, 30*time.Second)
			},
			value:            []byte("new"),
			generation:       20,
			expiration:       0,
			expectedAccepted: true,
			expectedValue:    "new",
			expectedMarker:   mk(t, 20, "new"),
			expectedTTL:      0,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			mr, rc := newMiniRedisCache(t)
			p.setup(t, mr)

			accepted, err := rc.PutIfNewer(testValueKey, testGenKey, p.value, p.generation, p.expiration)
			require.NoError(t, err)
			assert.Equal(t, p.expectedAccepted, accepted)

			gotValue, err := mr.Get(testValueKey)
			require.NoError(t, err)
			assert.Equal(t, p.expectedValue, gotValue)
			gotMarker, err := mr.Get(testGenKey)
			require.NoError(t, err)
			assert.Equal(t, p.expectedMarker, gotMarker)
			assert.Equal(t, p.expectedTTL, mr.TTL(testValueKey))
			assert.Equal(t, p.expectedTTL, mr.TTL(testGenKey))
		})
	}
}

func TestRedisCachePutIfNewerNegativeGeneration(t *testing.T) {
	t.Parallel()
	mr, rc := newMiniRedisCache(t)
	accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("v"), -1, 0)
	assert.ErrorIs(t, err, cache.ErrInvalidGeneration)
	assert.False(t, accepted)
	assert.False(t, mr.Exists(testValueKey))
	assert.False(t, mr.Exists(testGenKey))
}

func TestRedisCachePutIfNewerMaxGeneration(t *testing.T) {
	t.Parallel()
	mr, rc := newMiniRedisCache(t)
	accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("v"), math.MaxInt64, 0)
	require.NoError(t, err)
	assert.True(t, accepted)
	marker, err := mr.Get(testGenKey)
	require.NoError(t, err)
	assert.Equal(t, mk(t, math.MaxInt64, "v"), marker)
}

// Interleaving: regardless of the order in which an older and a newer snapshot
// reach Redis, the newer one is what remains.
func TestRedisCachePutIfNewerOrderIndependence(t *testing.T) {
	t.Parallel()
	t.Run("new then old", func(t *testing.T) {
		t.Parallel()
		mr, rc := newMiniRedisCache(t)
		accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("new"), 2, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
		accepted, err = rc.PutIfNewer(testValueKey, testGenKey, []byte("old"), 1, 0)
		require.NoError(t, err)
		assert.False(t, accepted)
		v, err := mr.Get(testValueKey)
		require.NoError(t, err)
		assert.Equal(t, "new", v)
	})
	t.Run("old then new", func(t *testing.T) {
		t.Parallel()
		mr, rc := newMiniRedisCache(t)
		accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("old"), 1, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
		accepted, err = rc.PutIfNewer(testValueKey, testGenKey, []byte("new"), 2, 0)
		require.NoError(t, err)
		assert.True(t, accepted)
		v, err := mr.Get(testValueKey)
		require.NoError(t, err)
		assert.Equal(t, "new", v)
	})
	t.Run("concurrent writers: highest generation wins", func(t *testing.T) {
		t.Parallel()
		mr, rc := newMiniRedisCache(t)
		const writers = 32
		var wg sync.WaitGroup
		for i := 1; i <= writers; i++ {
			wg.Add(1)
			go func(g int64) {
				defer wg.Done()
				_, err := rc.PutIfNewer(testValueKey, testGenKey, []byte{byte(g)}, g, 0)
				assert.NoError(t, err)
			}(int64(i))
		}
		wg.Wait()
		v, err := mr.Get(testValueKey)
		require.NoError(t, err)
		assert.Equal(t, string([]byte{writers}), v)
		marker, err := mr.Get(testGenKey)
		require.NoError(t, err)
		assert.Equal(t, mk(t, writers, string([]byte{writers})), marker)
	})
}

// Mixed-version rollout: a pod running the previous build does a plain SET on
// the value key. The next conditional writer must not be pinned out by the
// stale marker, even if its generation is older than the marker's.
func TestRedisCachePutIfNewerLegacyWriterInterleaving(t *testing.T) {
	t.Parallel()
	mr, rc := newMiniRedisCache(t)

	accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("gen-100"), 100, 0)
	require.NoError(t, err)
	require.True(t, accepted)

	// Legacy pod overwrites the value through the unconditional path.
	require.NoError(t, rc.Put(testValueKey, []byte("legacy"), 0))

	// A conditional writer with an older generation than the marker (100) is
	// still accepted because the marker no longer describes the stored value.
	accepted, err = rc.PutIfNewer(testValueKey, testGenKey, []byte("gen-50"), 50, 0)
	require.NoError(t, err)
	assert.True(t, accepted)
	v, err := mr.Get(testValueKey)
	require.NoError(t, err)
	assert.Equal(t, "gen-50", v)

	// From here on monotonicity is restored.
	accepted, err = rc.PutIfNewer(testValueKey, testGenKey, []byte("gen-40"), 40, 0)
	require.NoError(t, err)
	assert.False(t, accepted)
	v, err = mr.Get(testValueKey)
	require.NoError(t, err)
	assert.Equal(t, "gen-50", v)
}

// Legacy Evict deletes only the value key. The orphaned marker must not block
// the next write.
func TestRedisCachePutIfNewerAfterLegacyEvict(t *testing.T) {
	t.Parallel()
	mr, rc := newMiniRedisCache(t)

	accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("gen-100"), 100, 0)
	require.NoError(t, err)
	require.True(t, accepted)
	require.NoError(t, rc.Delete(testValueKey))
	require.True(t, mr.Exists(testGenKey))

	accepted, err = rc.PutIfNewer(testValueKey, testGenKey, []byte("gen-1"), 1, 0)
	require.NoError(t, err)
	assert.True(t, accepted)
	v, err := mr.Get(testValueKey)
	require.NoError(t, err)
	assert.Equal(t, "gen-1", v)
}

func TestRedisCachePutIfNewerEvalErrors(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		cmd         *goredis.Cmd
		expectedErr string
	}{
		{
			desc:        "transport error",
			cmd:         goredis.NewCmdResult(nil, errors.New("connection refused")),
			expectedErr: "connection refused",
		},
		{
			desc:        "non-integer reply",
			cmd:         goredis.NewCmdResult("not-an-int", nil),
			expectedErr: "invalid syntax",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			client := redismock.NewMockClient(ctrl)
			client.EXPECT().
				Eval(gomock.Any(), putIfNewerScript, []string{testValueKey, testGenKey},
					[]byte("v"), formatGen(t, 7), int64(1000), mk(t, 7, "v")).
				Return(p.cmd)
			rc := NewRedisCache(client).(*redisCache)
			accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("v"), 7, time.Second)
			assert.False(t, accepted)
			require.Error(t, err)
			assert.Contains(t, err.Error(), p.expectedErr)
		})
	}
}

func TestRedisCachePutIfNewerPassesArgs(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	client := redismock.NewMockClient(ctrl)
	client.EXPECT().
		Eval(
			gomock.AssignableToTypeOf(context.Background()),
			putIfNewerScript,
			[]string{testValueKey, testGenKey},
			[]byte("payload"),
			formatGen(t, 123),
			int64(1500),
			mk(t, 123, "payload"),
		).
		Return(goredis.NewCmdResult(int64(1), nil))
	rc := NewRedisCache(client).(*redisCache)
	accepted, err := rc.PutIfNewer(testValueKey, testGenKey, []byte("payload"), 123, 1500*time.Millisecond)
	require.NoError(t, err)
	assert.True(t, accepted)
}

// The script template must have its placeholder substituted, otherwise the
// Lua pattern would never match a well-formed marker.
func TestPutIfNewerScriptRendered(t *testing.T) {
	t.Parallel()
	assert.NotContains(t, putIfNewerScript, "__GENERATION_DIGITS__")
	assert.Contains(t, putIfNewerScript, "string.rep('%d', 20)")
	assert.Contains(t, putIfNewerScript, "redis.sha1hex")
}
