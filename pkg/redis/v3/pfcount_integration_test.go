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
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClusterPFCountIntegration verifies the staged cross-slot PFCount and
// PFMerge paths against a real Redis cluster, which miniredis cannot emulate
// (it stores HLLs as a dedicated type instead of plain strings). Skipped
// unless REDIS_CLUSTER_ADDR points at a cluster node, e.g.:
//
//	docker run -d --rm --name redis-cluster -e IP=0.0.0.0 \
//	  -p 7000-7005:7000-7005 grokzen/redis-cluster:7.0.10
//	REDIS_CLUSTER_ADDR=localhost:7000 go test ./pkg/redis/v3/ \
//	  -run TestClusterPFCountIntegration -v -count=1
func TestClusterPFCountIntegration(t *testing.T) {
	addr := os.Getenv("REDIS_CLUSTER_ADDR")
	if addr == "" {
		t.Skip("REDIS_CLUSTER_ADDR not set; skipping real-cluster integration test")
	}
	cl, err := NewClient(addr, WithRedisMode(RedisModeCluster))
	require.NoError(t, err)
	// Registered first so it runs last: cleanups run in reverse registration
	// order, and the key deletions below need a working client.
	t.Cleanup(func() { _ = cl.Close() })
	c := cl.(*client)
	require.Equal(t, ClientTypeCluster, c.clientType)

	prefix := fmt.Sprintf("pfcount-it:%d", time.Now().UnixNano())
	key1 := prefix + ":key1"
	key2 := prefix + ":key2"
	require.NotEqual(t, keyHashSlot(key1), keyHashSlot(key2),
		"test keys must land in different slots")
	t.Cleanup(func() {
		_ = c.Del(key1)
		_ = c.Del(key2)
	})

	// Overlapping members across slots: the union must deduplicate u2.
	_, err = c.PFAdd(key1, "u1", "u2")
	require.NoError(t, err)
	_, err = c.PFAdd(key2, "u2", "u3")
	require.NoError(t, err)

	count, err := c.PFCount(key1, key2)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count, "cross-slot union must deduplicate, not sum")

	// Missing keys contribute nothing to the union.
	count, err = c.PFCount(key1, key2, prefix+":missing")
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// All-missing cross-slot keys count to zero.
	count, err = c.PFCount(prefix+":none-1", prefix+":none-2")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// The cross-slot PFMerge path shares the staging helpers; verify it
	// produces the same union.
	dest := prefix + ":merged"
	require.NoError(t, c.PFMerge(dest, time.Minute, key1, key2))
	t.Cleanup(func() { _ = c.Del(dest) })
	count, err = c.PFCount(dest)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// A wrong-type source must fail, like native PFCOUNT raises WRONGTYPE.
	wrongType := prefix + ":wrong-type"
	require.NoError(t, c.rc.LPush(context.TODO(), wrongType, "element").Err())
	t.Cleanup(func() { _ = c.Del(wrongType) })
	_, err = c.PFCount(key1, wrongType)
	require.Error(t, err)

	// After successful counts and merges, no staging keys (always built as
	// "{tag}uuid") may remain anywhere in the cluster.
	cc, ok := c.rc.(*goredis.ClusterClient)
	require.True(t, ok)
	var mu sync.Mutex
	leftover := []string{}
	err = cc.ForEachMaster(context.TODO(), func(ctx context.Context, master *goredis.Client) error {
		keys, err := master.Keys(ctx, "{*").Result()
		if err != nil {
			return err
		}
		mu.Lock()
		leftover = append(leftover, keys...)
		mu.Unlock()
		return nil
	})
	require.NoError(t, err)
	assert.Empty(t, leftover, "staging keys must be cleaned up")
}

// TestClusterScanAllIntegration verifies that ScanAll enumerates keys from
// EVERY cluster master: plain Scan on a cluster client reaches only one node
// per call and silently misses the rest, which is exactly what ScanAll exists
// to fix (the daily-count backfill job relies on exhaustive enumeration).
func TestClusterScanAllIntegration(t *testing.T) {
	addr := os.Getenv("REDIS_CLUSTER_ADDR")
	if addr == "" {
		t.Skip("REDIS_CLUSTER_ADDR not set; skipping real-cluster integration test")
	}
	cl, err := NewClient(addr, WithRedisMode(RedisModeCluster))
	require.NoError(t, err)
	t.Cleanup(func() { _ = cl.Close() })
	c := cl.(*client)
	require.Equal(t, ClientTypeCluster, c.clientType)

	prefix := fmt.Sprintf("scanall-it-%d", time.Now().UnixNano())
	// Enough keys to land on every master with overwhelming probability
	// (slots are spread across 3 masters in the grokzen cluster).
	expected := make([]string, 0, 50)
	slots := make(map[int]bool)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("%s:ec:%d:feature:variation", prefix, i)
		require.NoError(t, c.Set(key, "1", time.Minute))
		expected = append(expected, key)
		slots[keyHashSlot(key)] = true
	}
	require.Greater(t, len(slots), 1, "test keys must span multiple slots")
	t.Cleanup(func() {
		for _, key := range expected {
			_ = c.Del(key)
		}
	})

	// The callback must never run concurrently (asserted via the running
	// flag) and must deliver every key from every master.
	var running atomic.Bool
	keys := []string{}
	err = c.ScanAll(prefix+":ec:*", 10, func(page []string) error {
		require.True(t, running.CompareAndSwap(false, true), "fn must not be called concurrently")
		defer running.Store(false)
		keys = append(keys, page...)
		return nil
	})
	require.NoError(t, err)
	assert.ElementsMatch(t, expected, keys,
		"ScanAll must deliver every key from every master")
}
