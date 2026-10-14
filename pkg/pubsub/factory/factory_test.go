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

package factory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/redis"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func TestNewClientRedisStreamMaxLen(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc string
		opts []Option
		// XAdd must receive this limit when publishing.
		expectedXAddMaxLen int64
	}{
		{desc: "custom limit is passed through", opts: []Option{WithStreamMaxLen(99)}, expectedXAddMaxLen: 99},
		{desc: "explicit zero disables trimming", opts: []Option{WithStreamMaxLen(0)}, expectedXAddMaxLen: 0},
		{desc: "explicit negative disables trimming", opts: []Option{WithStreamMaxLen(-1)}, expectedXAddMaxLen: -1},
		// Backend default: the factory only overrides when the option is used.
		{desc: "unset keeps the backend default", opts: nil, expectedXAddMaxLen: 10000},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			rc := redismock.NewMockClient(ctrl)
			rc.EXPECT().
				XAdd(gomock.Any(), gomock.Any(), gomock.Any(), p.expectedXAddMaxLen).
				Return("1-0", nil)

			opts := append([]Option{
				WithPubSubType(RedisStream),
				WithRedisClient(rc),
			}, p.opts...)
			client, err := NewClient(context.Background(), opts...)
			require.NoError(t, err)
			_, ok := client.(*redis.StreamClient)
			require.True(t, ok)

			pub, err := client.CreatePublisher("topic")
			require.NoError(t, err)
			assert.NoError(t, pub.Publish(context.Background(), &domaineventproto.Event{Id: "e1"}))
		})
	}
}
