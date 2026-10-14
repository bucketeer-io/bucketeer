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

package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
)

func TestStreamPublisherDefaultMaxLen(t *testing.T) {
	t.Parallel()
	p := NewStreamPublisher(nil, "topic").(*StreamPublisher)
	assert.Equal(t, defaultStreamMaxLen, p.maxLen)
	assert.Greater(t, defaultStreamMaxLen, int64(0), "trimming must be on by default")
}

func TestWithStreamPublisherMaxLen(t *testing.T) {
	t.Parallel()
	p := NewStreamPublisher(nil, "topic", WithStreamPublisherMaxLen(42)).(*StreamPublisher)
	assert.Equal(t, int64(42), p.maxLen)
}

// Publish must pass the configured limit to XADD on every call; the trim is
// done by Redis as part of the add, not as a separate command.
func TestStreamPublisherPublishPassesMaxLenToXAdd(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc   string
		maxLen int64
	}{
		{desc: "custom limit", maxLen: 123},
		{desc: "zero disables trimming", maxLen: 0},
	}
	for _, pt := range patterns {
		t.Run(pt.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			client := redismock.NewMockClient(ctrl)
			client.EXPECT().
				XAdd(gomock.Any(), testStreamKey("topic", 0), gomock.Any(), pt.maxLen).
				Return("1-0", nil)
			p := NewStreamPublisher(client, "topic",
				WithStreamPublisherPartitionCount(1),
				WithStreamPublisherMaxLen(pt.maxLen),
			)
			require.NoError(t, p.Publish(context.Background(), testEvent("e1")))
		})
	}
}

// Against a real stream, publishing more than maxLen messages must leave the
// stream bounded and keep the newest entries.
func TestStreamPublisherTrimsStream(t *testing.T) {
	t.Parallel()
	const maxLen = 5
	mr, client := newTestRedis(t)
	p := NewStreamPublisher(client, "topic",
		WithStreamPublisherPartitionCount(1),
		WithStreamPublisherMaxLen(maxLen),
	)
	for i := 0; i < 20; i++ {
		require.NoError(t, p.Publish(context.Background(), testEvent(fmt.Sprintf("e%02d", i))))
	}
	entries, err := mr.Stream(testStreamKey("topic", 0))
	require.NoError(t, err)
	// MAXLEN ~ lets Redis keep a few extra entries; miniredis trims exactly.
	assert.LessOrEqual(t, len(entries), maxLen)
	assert.NotEmpty(t, entries)
	// The newest message must survive trimming.
	last := entries[len(entries)-1]
	require.Len(t, last.Values, 2, "one field/value pair per entry")
	assert.Equal(t, "e19", last.Values[0], "field name is the message ID")
}

func TestStreamPublisherMaxLenZeroDoesNotTrim(t *testing.T) {
	t.Parallel()
	mr, client := newTestRedis(t)
	p := NewStreamPublisher(client, "topic",
		WithStreamPublisherPartitionCount(1),
		WithStreamPublisherMaxLen(0),
	)
	for i := 0; i < 20; i++ {
		require.NoError(t, p.Publish(context.Background(), testEvent(fmt.Sprintf("e%02d", i))))
	}
	entries, err := mr.Stream(testStreamKey("topic", 0))
	require.NoError(t, err)
	assert.Len(t, entries, 20)
}

// Trimming is per partition stream, so messages hashed to different
// partitions do not count against each other.
func TestStreamPublisherTrimsPerPartition(t *testing.T) {
	t.Parallel()
	const maxLen = 3
	mr, client := newTestRedis(t)
	sp := NewStreamPublisher(client, "topic",
		WithStreamPublisherPartitionCount(4),
		WithStreamPublisherMaxLen(maxLen),
	)
	p := sp.(*StreamPublisher)
	for i := 0; i < 100; i++ {
		require.NoError(t, sp.Publish(context.Background(), testEvent(fmt.Sprintf("e%03d", i))))
	}
	total := 0
	for partition := 0; partition < 4; partition++ {
		entries, err := mr.Stream(testStreamKey("topic", partition))
		if err != nil {
			continue // partition never written
		}
		assert.LessOrEqual(t, len(entries), maxLen, "partition %d", partition)
		total += len(entries)
	}
	assert.Greater(t, total, maxLen, "100 ids should spread over more than one partition")
	assert.Equal(t, 4, p.partitionCount)
}
