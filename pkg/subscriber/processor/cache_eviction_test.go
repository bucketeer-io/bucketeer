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

package processor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type testNetError struct {
	timeout   bool
	temporary bool
	msg       string
}

func (e *testNetError) Error() string   { return e.msg }
func (e *testNetError) Timeout() bool   { return e.timeout }
func (e *testNetError) Temporary() bool { return e.temporary }

var _ net.Error = (*testNetError)(nil)

func TestCacheEvictionHandleMessage(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc  string
		event *domaineventproto.Event
		setup func(
			fc *cachev3mock.MockFeaturesCache,
			sc *cachev3mock.MockSegmentUsersCache,
			ak *cachev3mock.MockEnvironmentAPIKeyCache,
			ec *cachev3mock.MockExperimentsCache,
			ao *cachev3mock.MockAutoOpsRulesCache,
		)
	}{
		{
			desc: "feature event evicts features cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				EntityId:      "feature-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_FEATURE_UPDATED,
			},
			setup: func(fc *cachev3mock.MockFeaturesCache, _ *cachev3mock.MockSegmentUsersCache, _ *cachev3mock.MockEnvironmentAPIKeyCache, _ *cachev3mock.MockExperimentsCache, _ *cachev3mock.MockAutoOpsRulesCache) {
				fc.EXPECT().Evict("env-1").Return(nil)
			},
		},
		{
			desc: "segment event evicts segment users cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_SEGMENT,
				EntityId:      "segment-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_SEGMENT_CREATED,
			},
			setup: func(_ *cachev3mock.MockFeaturesCache, sc *cachev3mock.MockSegmentUsersCache, _ *cachev3mock.MockEnvironmentAPIKeyCache, _ *cachev3mock.MockExperimentsCache, _ *cachev3mock.MockAutoOpsRulesCache) {
				sc.EXPECT().Evict("segment-id-1", "env-1").Return(nil)
			},
		},
		{
			desc: "api key event evicts environment api key cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_APIKEY,
				EntityId:      "apikey-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_APIKEY_CHANGED,
				EntityData:    `{"api_key": "secret-123"}`,
			},
			setup: func(_ *cachev3mock.MockFeaturesCache, _ *cachev3mock.MockSegmentUsersCache, ak *cachev3mock.MockEnvironmentAPIKeyCache, _ *cachev3mock.MockExperimentsCache, _ *cachev3mock.MockAutoOpsRulesCache) {
				ak.EXPECT().Evict("secret-123").Return(nil)
			},
		},
		{
			desc: "api key event evicts both previous and current secrets",
			event: &domaineventproto.Event{
				EntityType:         domaineventproto.Event_APIKEY,
				EntityId:           "apikey-id-1",
				EnvironmentId:      "env-1",
				Type:               domaineventproto.Event_APIKEY_CHANGED,
				EntityData:         `{"api_key": "new-secret"}`,
				PreviousEntityData: `{"api_key": "old-secret"}`,
			},
			setup: func(_ *cachev3mock.MockFeaturesCache, _ *cachev3mock.MockSegmentUsersCache, ak *cachev3mock.MockEnvironmentAPIKeyCache, _ *cachev3mock.MockExperimentsCache, _ *cachev3mock.MockAutoOpsRulesCache) {
				ak.EXPECT().Evict("old-secret").Return(nil)
				ak.EXPECT().Evict("new-secret").Return(nil)
			},
		},
		{
			desc: "experiment event evicts experiments cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_EXPERIMENT,
				EntityId:      "experiment-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_EXPERIMENT_CREATED,
			},
			setup: func(_ *cachev3mock.MockFeaturesCache, _ *cachev3mock.MockSegmentUsersCache, _ *cachev3mock.MockEnvironmentAPIKeyCache, ec *cachev3mock.MockExperimentsCache, _ *cachev3mock.MockAutoOpsRulesCache) {
				ec.EXPECT().Evict("env-1").Return(nil)
			},
		},
		{
			desc: "auto ops rule event evicts auto ops rules cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_AUTOOPS_RULE,
				EntityId:      "autoops-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_AUTOOPS_RULE_CREATED,
			},
			setup: func(_ *cachev3mock.MockFeaturesCache, _ *cachev3mock.MockSegmentUsersCache, _ *cachev3mock.MockEnvironmentAPIKeyCache, _ *cachev3mock.MockExperimentsCache, ao *cachev3mock.MockAutoOpsRulesCache) {
				ao.EXPECT().Evict("env-1").Return(nil)
			},
		},
		{
			desc: "unrelated entity type does not evict any cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_ACCOUNT,
				EntityId:      "account-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_ACCOUNT_V2_CREATED,
			},
			setup: func(_ *cachev3mock.MockFeaturesCache, _ *cachev3mock.MockSegmentUsersCache, _ *cachev3mock.MockEnvironmentAPIKeyCache, _ *cachev3mock.MockExperimentsCache, _ *cachev3mock.MockAutoOpsRulesCache) {
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			fc := cachev3mock.NewMockFeaturesCache(ctrl)
			sc := cachev3mock.NewMockSegmentUsersCache(ctrl)
			ak := cachev3mock.NewMockEnvironmentAPIKeyCache(ctrl)
			ec := cachev3mock.NewMockExperimentsCache(ctrl)
			ao := cachev3mock.NewMockAutoOpsRulesCache(ctrl)

			p.setup(fc, sc, ak, ec, ao)

			processor := NewCacheEviction(fc, sc, ak, ec, ao, zap.NewNop())
			ce := processor.(*cacheEviction)

			data, err := proto.Marshal(p.event)
			assert.NoError(t, err)

			msg := &puller.Message{
				Data: data,
				Ack:  func() {},
				Nack: func() {},
			}
			ce.handleMessage(msg)
		})
	}
}

func TestCacheEvictionHandleMessageNackOnRepeatableError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	fc := cachev3mock.NewMockFeaturesCache(ctrl)
	sc := cachev3mock.NewMockSegmentUsersCache(ctrl)
	ak := cachev3mock.NewMockEnvironmentAPIKeyCache(ctrl)
	ec := cachev3mock.NewMockExperimentsCache(ctrl)
	ao := cachev3mock.NewMockAutoOpsRulesCache(ctrl)

	fc.EXPECT().Evict("env-1").Return(context.DeadlineExceeded)

	processor := NewCacheEviction(fc, sc, ak, ec, ao, zap.NewNop())
	ce := processor.(*cacheEviction)

	data, err := proto.Marshal(&domaineventproto.Event{
		EntityType:    domaineventproto.Event_FEATURE,
		EntityId:      "feature-id-1",
		EnvironmentId: "env-1",
		Type:          domaineventproto.Event_FEATURE_UPDATED,
	})
	assert.NoError(t, err)

	acked := false
	nacked := false
	msg := &puller.Message{
		Data: data,
		Ack: func() {
			acked = true
		},
		Nack: func() {
			nacked = true
		},
	}

	ce.handleMessage(msg)

	assert.False(t, acked)
	assert.True(t, nacked)
}

func TestIsRepeatable(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		err      error
		expected bool
	}{
		{
			desc:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			desc:     "context.Canceled is not repeatable",
			err:      context.Canceled,
			expected: false,
		},
		{
			desc:     "context.DeadlineExceeded is repeatable",
			err:      context.DeadlineExceeded,
			expected: true,
		},
		{
			desc:     "cache.ErrNotFound is not repeatable",
			err:      cache.ErrNotFound,
			expected: false,
		},
		{
			desc:     "cache.ErrInvalidType is not repeatable",
			err:      cache.ErrInvalidType,
			expected: false,
		},
		{
			desc:     "net.Error with timeout is repeatable",
			err:      &testNetError{timeout: true, msg: "i/o timeout"},
			expected: true,
		},
		{
			desc:     "net.Error without timeout falls through to string check",
			err:      &testNetError{timeout: false, msg: "some net error"},
			expected: false,
		},
		{
			desc:     "net.Error without timeout but with connection reset message is repeatable",
			err:      &testNetError{timeout: false, msg: "read: connection reset by peer"},
			expected: true,
		},
		{
			desc:     "net.Error without timeout but with broken pipe message is repeatable",
			err:      &testNetError{timeout: false, msg: "write: broken pipe"},
			expected: true,
		},
		{
			desc:     "error containing timeout string is repeatable",
			err:      errors.New("redis: command Timeout exceeded"),
			expected: true,
		},
		{
			desc:     "error containing connection reset is repeatable",
			err:      errors.New("read: connection reset by peer"),
			expected: true,
		},
		{
			desc:     "error containing broken pipe is repeatable",
			err:      errors.New("write: broken pipe"),
			expected: true,
		},
		{
			desc:     "error containing connection refused is repeatable",
			err:      errors.New("dial tcp: connection refused"),
			expected: true,
		},
		{
			desc:     "error containing eof is repeatable",
			err:      errors.New("unexpected EOF"),
			expected: true,
		},
		{
			desc:     "wrapped context.DeadlineExceeded is repeatable",
			err:      fmt.Errorf("operation failed: %w", context.DeadlineExceeded),
			expected: true,
		},
		{
			desc:     "wrapped cache.ErrNotFound is not repeatable",
			err:      fmt.Errorf("lookup failed: %w", cache.ErrNotFound),
			expected: false,
		},
		{
			desc:     "generic error is not repeatable",
			err:      errors.New("something unexpected"),
			expected: false,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, isRepeatable(p.err))
		})
	}
}
