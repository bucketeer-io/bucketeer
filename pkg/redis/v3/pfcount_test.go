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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPFCountClusterSingleKey(t *testing.T) {
	t.Parallel()
	_, c := newMiniredisClient(t, ClientTypeCluster)

	_, err := c.PFAdd("env:uc:1:feature:variation", "u1", "u2", "u3")
	require.NoError(t, err)

	count, err := c.PFCount("env:uc:1:feature:variation")
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestPFCountClusterSameSlotKeys(t *testing.T) {
	t.Parallel()
	_, c := newMiniredisClient(t, ClientTypeCluster)

	// Keys sharing a hash tag land in one slot, so multi-key PFCOUNT is a
	// single valid server-side command. Disjoint members keep the expected
	// value identical under union and sum semantics (miniredis sums).
	_, err := c.PFAdd("{tag}a", "u1", "u2")
	require.NoError(t, err)
	_, err = c.PFAdd("{tag}b", "u3")
	require.NoError(t, err)

	count, err := c.PFCount("{tag}a", "{tag}b")
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestPFCountClusterCrossSlotAllKeysMissing(t *testing.T) {
	t.Parallel()
	_, c := newMiniredisClient(t, ClientTypeCluster)

	keys := []string{"env:uc:1:feature:variation", "env:uc:2:feature:variation"}
	require.NotEqual(t, keyHashSlot(keys[0]), keyHashSlot(keys[1]))

	// Cross-slot keys go through the staged-union path; when none of the
	// source keys exist, the union is empty and nothing is staged.
	count, err := c.PFCount(keys...)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestPFCountClusterCrossSlotErrorCleansUpStagingKeys(t *testing.T) {
	t.Parallel()
	// The happy cross-slot path cannot run against miniredis: it stores
	// HLLs as a dedicated type, so MGET cannot fetch their raw values the
	// way real Redis does (where HLLs are plain strings). Instead, store
	// plain strings so the values are fetched and staged, and let the
	// staged PFCOUNT fail (miniredis rejects PFCOUNT on string keys) to
	// verify the error surfaces and the staging keys are cleaned up.
	mr, c := newMiniredisClient(t, ClientTypeCluster)

	keys := []string{"env:uc:1:feature:variation", "env:uc:2:feature:variation"}
	require.NotEqual(t, keyHashSlot(keys[0]), keyHashSlot(keys[1]))
	require.NoError(t, mr.Set(keys[0], "not-an-hll"))
	require.NoError(t, mr.Set(keys[1], "not-an-hll-either"))

	_, err := c.PFCount(keys...)
	require.Error(t, err)

	// Only the two source keys may remain; all staging keys must be gone.
	assert.ElementsMatch(t, keys, mr.Keys())
}

// A wrong-type source (here a list) must fail the cross-slot count, like
// native PFCOUNT raises WRONGTYPE. MGET would silently treat such keys as
// missing; the per-key GETs used by the staged path propagate the error.
func TestPFCountClusterCrossSlotWrongTypeKeyFails(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeCluster)

	keys := []string{"env:uc:1:feature:variation", "env:uc:2:feature:variation"}
	require.NotEqual(t, keyHashSlot(keys[0]), keyHashSlot(keys[1]))
	_, err := mr.Lpush(keys[0], "element")
	require.NoError(t, err)

	_, err = c.PFCount(keys...)
	require.Error(t, err)
	assert.ElementsMatch(t, []string{keys[0]}, mr.Keys(),
		"no staging keys may be written when fetching the sources fails")
}

// An existing empty-string source must be staged and rejected server side,
// like native PFCOUNT rejects a key that is not a valid HLL string, instead
// of being silently skipped.
func TestPFCountClusterCrossSlotEmptyStringKeyFails(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeCluster)

	keys := []string{"env:uc:1:feature:variation", "env:uc:2:feature:variation"}
	require.NotEqual(t, keyHashSlot(keys[0]), keyHashSlot(keys[1]))
	require.NoError(t, mr.Set(keys[0], ""))

	_, err := c.PFCount(keys...)
	require.Error(t, err)
	// Only the source key may remain; the staging keys must be cleaned up.
	assert.ElementsMatch(t, []string{keys[0]}, mr.Keys())
}
