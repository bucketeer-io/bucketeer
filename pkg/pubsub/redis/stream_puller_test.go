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
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	v3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

const (
	pullerTestTopic   = "topic"
	pullerTestGroup   = "group"
	pullerTestTimeout = 5 * time.Second
	// How long to wait to be confident that nothing else is delivered.
	pullerTestQuiet = 300 * time.Millisecond
)

// startPuller runs Pull in the background and returns a channel of received
// message IDs (the published event's ID, decoded from the payload).
// It blocks until the consumer group exists on partition 0 so the caller can
// publish "after group creation" deterministically.
func startPuller(
	ctx context.Context,
	t *testing.T,
	client v3.Client,
	subscription string,
) (<-chan string, <-chan error) {
	t.Helper()
	received := make(chan string, 100)
	p := NewStreamPuller(client, subscription, pullerTestTopic,
		WithStreamPullerPartitionCount(1),
		WithStreamPullerBlockTime(50*time.Millisecond),
	)
	errCh := make(chan error, 1)
	go func() {
		errCh <- p.Pull(ctx, func(_ context.Context, msg *puller.Message) {
			msg.Ack()
			event := &domaineventproto.Event{}
			if err := proto.Unmarshal(msg.Data, event); err != nil {
				received <- "unmarshal error: " + err.Error()
				return
			}
			received <- event.Id
		})
	}()
	require.Eventually(t, func() bool {
		groups, err := client.XInfoGroups(ctx, testStreamKey(pullerTestTopic, 0))
		if err != nil {
			return false
		}
		for _, g := range groups {
			if g.Name == subscription {
				return true
			}
		}
		return false
	}, pullerTestTimeout, 10*time.Millisecond, "consumer group was not created")
	return received, errCh
}

func collect(t *testing.T, received <-chan string, n int) []string {
	t.Helper()
	ids := make([]string, 0, n)
	deadline := time.After(pullerTestTimeout)
	for len(ids) < n {
		select {
		case id := <-received:
			ids = append(ids, id)
		case <-deadline:
			t.Fatalf("expected %d messages, got %d: %v", n, len(ids), ids)
		}
	}
	return ids
}

func assertQuiet(t *testing.T, received <-chan string) {
	t.Helper()
	select {
	case id := <-received:
		t.Fatalf("unexpected extra message delivered: %s", id)
	case <-time.After(pullerTestQuiet):
	}
}

// A new consumer group must start at the stream tail: messages published
// before the group existed are never delivered, messages published after it
// are. This is the Google Pub/Sub behaviour for a new subscription and is what
// stops a re-pointed processor from replaying the whole retained stream.
func TestStreamPullerNewGroupStartsAtTail(t *testing.T) {
	t.Parallel()
	_, client := newTestRedis(t)
	pub := NewStreamPublisher(client, pullerTestTopic, WithStreamPublisherPartitionCount(1))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 3; i++ {
		require.NoError(t, pub.Publish(ctx, testEvent(fmt.Sprintf("old-%d", i))))
	}

	received, errCh := startPuller(ctx, t, client, pullerTestGroup)

	require.NoError(t, pub.Publish(ctx, testEvent("new-0")))
	require.NoError(t, pub.Publish(ctx, testEvent("new-1")))

	assert.ElementsMatch(t, []string{"new-0", "new-1"}, collect(t, received, 2))
	assertQuiet(t, received)

	cancel()
	select {
	case err := <-errCh:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(pullerTestTimeout):
		t.Fatal("Pull did not return after cancel")
	}
}

// An existing consumer group is left untouched, so an operator can pre-create
// a group at any position (for example "0" to replay, or "$" ahead of a
// deployment) and the puller will honour it.
func TestStreamPullerKeepsExistingGroup(t *testing.T) {
	t.Parallel()
	_, client := newTestRedis(t)
	pub := NewStreamPublisher(client, pullerTestTopic, WithStreamPublisherPartitionCount(1))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 3; i++ {
		require.NoError(t, pub.Publish(ctx, testEvent(fmt.Sprintf("old-%d", i))))
	}
	// Pre-created at "0": the whole history is wanted.
	require.NoError(t, client.XGroupCreateMkStream(testStreamKey(pullerTestTopic, 0), pullerTestGroup, "0"))

	received, _ := startPuller(ctx, t, client, pullerTestGroup)
	require.NoError(t, pub.Publish(ctx, testEvent("new-0")))

	assert.ElementsMatch(t, []string{"old-0", "old-1", "old-2", "new-0"}, collect(t, received, 4))
	assertQuiet(t, received)
}

// Two independent groups on the same stream each start from their own
// creation point, so adding a second consumer never replays into it.
func TestStreamPullerSecondGroupDoesNotReplay(t *testing.T) {
	t.Parallel()
	_, client := newTestRedis(t)
	pub := NewStreamPublisher(client, pullerTestTopic, WithStreamPublisherPartitionCount(1))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	first, _ := startPuller(ctx, t, client, "first")
	for i := 0; i < 3; i++ {
		require.NoError(t, pub.Publish(ctx, testEvent(fmt.Sprintf("a-%d", i))))
	}
	assert.ElementsMatch(t, []string{"a-0", "a-1", "a-2"}, collect(t, first, 3))

	second, _ := startPuller(ctx, t, client, "second")
	require.NoError(t, pub.Publish(ctx, testEvent("b-0")))

	assert.Equal(t, []string{"b-0"}, collect(t, second, 1))
	assertQuiet(t, second)
	assert.Equal(t, []string{"b-0"}, collect(t, first, 1))
}

// The group is created with "$" on every partition, and not re-created when
// it already exists.
func TestStreamPullerCreatesGroupsAtTail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	client := redismock.NewMockClient(ctrl)
	const partitions = 3
	for partition := 0; partition < partitions; partition++ {
		key := testStreamKey(pullerTestTopic, partition)
		if partition == 1 {
			// Existing group on this partition: must not be recreated.
			client.EXPECT().Exists(key).Return(int64(1), nil)
			client.EXPECT().XInfoGroups(gomock.Any(), key).
				Return([]goredis.XInfoGroup{{Name: pullerTestGroup}}, nil)
			continue
		}
		client.EXPECT().Exists(key).Return(int64(0), nil)
		client.EXPECT().XGroupCreateMkStream(key, pullerTestGroup, "$").Return(nil)
	}
	// Stop the read loop at the first read.
	client.EXPECT().XReadGroup(
		gomock.Any(), pullerTestGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil, context.Canceled).AnyTimes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := NewStreamPuller(client, pullerTestGroup, pullerTestTopic,
		WithStreamPullerPartitionCount(partitions),
	)
	err := p.Pull(ctx, func(context.Context, *puller.Message) {})
	assert.ErrorIs(t, err, context.Canceled)
}

var _ publisher.Publisher = (*StreamPublisher)(nil)
