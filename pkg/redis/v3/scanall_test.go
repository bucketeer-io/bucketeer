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
	"sort"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestScanAllStandalone verifies that ScanAll enumerates every matching key
// exactly once on the standalone path, across multiple SCAN pages. The
// cluster path (one iteration per master) is covered by
// TestClusterScanAllIntegration.
func TestScanAllStandalone(t *testing.T) {
	t.Parallel()
	mr := miniredis.RunT(t)
	cl, err := NewClient(mr.Addr(), WithRedisMode(RedisModeStandalone), WithLogger(zap.NewNop()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = cl.Close() })

	expected := make([]string, 0, 25)
	for i := 0; i < 25; i++ {
		key := fmt.Sprintf("env:ec:%d:feature:variation", i)
		require.NoError(t, cl.Set(key, "1", time.Minute))
		expected = append(expected, key)
	}
	// Non-matching keys must not be returned.
	require.NoError(t, cl.Set("env:ec1d:0:feature:variation", "1", time.Minute))
	require.NoError(t, cl.Set("other:ec:0:feature:variation", "1", time.Minute))

	// A small page size forces multiple SCAN iterations.
	var keys []string
	err = cl.ScanAll("env:ec:*", 5, func(page []string) error {
		keys = append(keys, page...)
		return nil
	})
	require.NoError(t, err)
	sort.Strings(keys)
	sort.Strings(expected)
	assert.Equal(t, expected, keys)

	// No matches yields no callback invocations, not an error.
	err = cl.ScanAll("missing:*", 5, func(page []string) error {
		t.Errorf("fn must not be called for an empty result, got %v", page)
		return nil
	})
	require.NoError(t, err)

	// An fn error aborts the scan and is propagated.
	wantErr := fmt.Errorf("stop")
	err = cl.ScanAll("env:ec:*", 5, func(page []string) error {
		return wantErr
	})
	assert.ErrorIs(t, err, wantErr)
}

// TestSerializeLatchesFirstError verifies the cluster-mode abort semantics:
// once fn returns an error, parallel master workers that deliver later pages
// must get that same error back without fn being invoked again.
func TestSerializeLatchesFirstError(t *testing.T) {
	t.Parallel()
	wantErr := fmt.Errorf("stop")
	calls := 0
	fn := serialize(func(keys []string) error {
		calls++
		if calls == 2 {
			return wantErr
		}
		return nil
	})

	require.NoError(t, fn([]string{"page-1"}))
	require.ErrorIs(t, fn([]string{"page-2"}), wantErr)
	// Later pages from other masters: same error, fn not invoked again.
	require.ErrorIs(t, fn([]string{"page-3"}), wantErr)
	require.ErrorIs(t, fn([]string{"page-4"}), wantErr)
	assert.Equal(t, 2, calls, "fn must not run after it requested an abort")
}
