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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The cluster branch groups keys by hash slot and fetches all slots in one
// pipeline. Running it against a single node still exercises the grouping,
// the pipelined MGETs, and the reassembly in original key order.
func TestGetMultiCluster(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeCluster)

	keys := make([]string, 0, 20)
	want := make([]interface{}, 0, 20)
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("env:ec:%d:feature:variation", i)
		keys = append(keys, key)
		if i%3 == 0 {
			// Leave every third key missing.
			want = append(want, "")
			continue
		}
		val := fmt.Sprintf("value-%d", i)
		require.NoError(t, mr.Set(key, val))
		want = append(want, []byte(val))
	}

	// Sanity check: the keys must span multiple slots so the pipeline
	// actually contains multiple MGET commands.
	slots := map[int]bool{}
	for _, k := range keys {
		slots[keyHashSlot(k)] = true
	}
	require.Greater(t, len(slots), 1)

	values, err := c.GetMulti(keys, true)
	require.NoError(t, err)
	assert.Equal(t, want, values, "values must come back in the original key order")
}

func TestGetMultiClusterNotFound(t *testing.T) {
	t.Parallel()
	mr, c := newMiniredisClient(t, ClientTypeCluster)
	require.NoError(t, mr.Set("env:ec:1:feature:variation", "v"))

	_, err := c.GetMulti([]string{"env:ec:1:feature:variation", "missing"}, false)
	assert.ErrorIs(t, err, ErrNil)
}
