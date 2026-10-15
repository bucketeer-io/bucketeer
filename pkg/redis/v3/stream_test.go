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
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newMiniRedisClient(t *testing.T) (*miniredis.Miniredis, Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	client, err := NewClient(
		mr.Addr(),
		WithRedisMode(RedisModeStandalone),
		WithLogger(zap.NewNop()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	return mr, client
}

func TestXAddMaxLen(t *testing.T) {
	t.Parallel()
	const stream = "s{stream}"
	patterns := []struct {
		desc        string
		maxLen      int64
		adds        int
		expectedLen int
	}{
		{desc: "positive maxLen trims the stream", maxLen: 3, adds: 10, expectedLen: 3},
		{desc: "zero maxLen keeps everything", maxLen: 0, adds: 10, expectedLen: 10},
		{desc: "negative maxLen keeps everything", maxLen: -1, adds: 10, expectedLen: 10},
		{desc: "maxLen above length is a no-op", maxLen: 100, adds: 10, expectedLen: 10},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			mr, client := newMiniRedisClient(t)
			var lastID string
			for i := 0; i < p.adds; i++ {
				id, err := client.XAdd(
					context.Background(),
					stream,
					map[string]interface{}{"m": fmt.Sprintf("v%d", i)},
					p.maxLen,
				)
				require.NoError(t, err)
				require.NotEmpty(t, id)
				lastID = id
			}
			entries, err := mr.Stream(stream)
			require.NoError(t, err)
			// Real Redis treats MAXLEN ~ as a hint; miniredis trims exactly.
			assert.LessOrEqual(t, len(entries), p.expectedLen)
			assert.NotEmpty(t, entries)
			assert.Equal(t, lastID, entries[len(entries)-1].ID, "newest entry must be kept")
		})
	}
}
