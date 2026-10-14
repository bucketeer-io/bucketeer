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
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	v3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

// newTestRedis runs an in-process Redis and returns a real v3 client bound to
// it, so stream commands (XADD MAXLEN, XGROUP CREATE $, XREADGROUP) are
// exercised for real rather than mocked.
func newTestRedis(t *testing.T) (*miniredis.Miniredis, v3.Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	client, err := v3.NewClient(
		mr.Addr(),
		v3.WithRedisMode(v3.RedisModeStandalone),
		v3.WithLogger(zap.NewNop()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	return mr, client
}

// testStreamKey mirrors the partition naming used by publisher and puller.
func testStreamKey(topic string, partition int) string {
	return fmt.Sprintf("%s-%d{stream}", topic, partition)
}

func testEvent(id string) *domaineventproto.Event {
	return &domaineventproto.Event{
		Id:            id,
		EntityType:    domaineventproto.Event_FEATURE,
		EntityId:      "feature-" + id,
		EnvironmentId: "env-1",
		Type:          domaineventproto.Event_FEATURE_UPDATED,
	}
}
