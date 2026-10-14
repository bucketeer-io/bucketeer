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
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHashTagForKey(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc    string
		key     string
		wantTag string
		wantOK  bool
	}{
		{
			desc:    "plain key becomes its own tag",
			key:     "env-1:pfmerge-key:feature:variation",
			wantTag: "env-1:pfmerge-key:feature:variation",
			wantOK:  true,
		},
		{
			desc:    "existing tag is reused",
			key:     "a{user1000}c",
			wantTag: "user1000",
			wantOK:  true,
		},
		{
			desc:    "tag containing an opening brace is reused",
			key:     "x{y{z}w",
			wantTag: "y{z",
			wantOK:  true,
		},
		{
			desc:   "empty tag cannot be wrapped",
			key:    "{}foo",
			wantOK: false,
		},
		{
			desc:   "empty key cannot be wrapped",
			key:    "",
			wantOK: false,
		},
		{
			desc:   "closing brace without tag cannot be wrapped",
			key:    "a}b",
			wantOK: false,
		},
		{
			desc:   "unclosed brace cannot be wrapped",
			key:    "a{b",
			wantOK: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			tag, ok := hashTagForKey(p.key)
			assert.Equal(t, p.wantOK, ok)
			if !p.wantOK {
				return
			}
			assert.Equal(t, p.wantTag, tag)
			// The whole point of the tag: a staging key built from it must
			// hash to the same cluster slot as the original key.
			stagingKey := "{" + tag + "}staging-suffix"
			assert.Equal(t, keyHashSlot(p.key), keyHashSlot(stagingKey),
				"staging key %q must share the slot of %q", stagingKey, p.key)
		})
	}
}

func newMiniredisClient(t *testing.T, clientType ClientType) (*miniredis.Miniredis, *client) {
	t.Helper()
	mr := miniredis.RunT(t)
	rc := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rc.Close() })
	return mr, &client{
		rc:         rc,
		opts:       defaultOptions(),
		logger:     zap.NewNop(),
		clientType: clientType,
	}
}

func TestPFMergeStandalone(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeStandard)

	_, err := c.PFAdd("src1", "a", "b")
	require.NoError(t, err)
	_, err = c.PFAdd("src2", "b", "c")
	require.NoError(t, err)

	require.NoError(t, c.PFMerge("dest", time.Hour, "src1", "src2"))

	count, err := c.PFCount("dest")
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
	assert.Equal(t, time.Hour, mr.TTL("dest"))
}

func TestPFMergeStandaloneNoSources(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeStandard)

	require.NoError(t, c.PFMerge("dest", time.Hour, "missing1", "missing2"))

	count, err := c.PFCount("dest")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
	assert.True(t, mr.Exists("dest"))
	assert.Equal(t, time.Hour, mr.TTL("dest"))
}

// The cluster branch with no existing source keys only issues same-key
// commands (PFMERGE dest + EXPIRE dest), so it can run against miniredis.
func TestPFMergeClusterNoSources(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeCluster)

	require.NoError(t, c.PFMerge("dest", time.Hour, "missing1", "missing2"))

	count, err := c.PFCount("dest")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
	assert.True(t, mr.Exists("dest"))
	assert.Equal(t, time.Hour, mr.TTL("dest"))
}

// Miniredis stores HLLs as a dedicated type, so PFMERGE rejects the staged
// string keys. That failure path verifies two properties of the cluster
// branch: errors from the merge pipeline are propagated, and the deferred
// cleanup still removes every staging key.
func TestPFMergeClusterMergeErrorCleansUpStagingKeys(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeCluster)

	require.NoError(t, mr.Set("src1", "not-an-hll-1"))
	require.NoError(t, mr.Set("src2", "not-an-hll-2"))

	err := c.PFMerge("dest", time.Hour, "src1", "src2")
	require.Error(t, err)

	for _, key := range mr.Keys() {
		assert.False(t, strings.HasPrefix(key, "{"),
			"staging key %q must be cleaned up after a failed merge", key)
	}
}

// Same failure-path coverage for the brace-containing dest fallback.
func TestPFMergeClusterTmpDestErrorCleansUpStagingKeys(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeCluster)

	require.NoError(t, mr.Set("src1", "not-an-hll-1"))

	err := c.PFMerge("{}dest", time.Hour, "src1")
	require.Error(t, err)

	for _, key := range mr.Keys() {
		assert.False(t, strings.HasPrefix(key, "{") && key != "{}dest",
			"staging key %q must be cleaned up after a failed merge", key)
	}
}
