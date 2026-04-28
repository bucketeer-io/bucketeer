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

	accountdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	accstorage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	accountstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2/mock"
	autoopsclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
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

type cacheRefresherMocks struct {
	featureClient          *featureclientmock.MockClient
	experimentClient       *experimentclientmock.MockClient
	autoOpsClient          *autoopsclientmock.MockClient
	accountStorage         *accountstoragemock.MockAccountStorage
	featuresCache          *cachev3mock.MockFeaturesCache
	segmentUsersCache      *cachev3mock.MockSegmentUsersCache
	environmentAPIKeyCache *cachev3mock.MockEnvironmentAPIKeyCache
	experimentsCache       *cachev3mock.MockExperimentsCache
	autoOpsRulesCache      *cachev3mock.MockAutoOpsRulesCache
	invalidationPublisher  *publishermock.MockPublisher
}

func newCacheRefresherWithMocks(t *testing.T) (*cacheRefresher, *cacheRefresherMocks) {
	ctrl := gomock.NewController(t)
	m := &cacheRefresherMocks{
		featureClient:          featureclientmock.NewMockClient(ctrl),
		experimentClient:       experimentclientmock.NewMockClient(ctrl),
		autoOpsClient:          autoopsclientmock.NewMockClient(ctrl),
		accountStorage:         accountstoragemock.NewMockAccountStorage(ctrl),
		featuresCache:          cachev3mock.NewMockFeaturesCache(ctrl),
		segmentUsersCache:      cachev3mock.NewMockSegmentUsersCache(ctrl),
		environmentAPIKeyCache: cachev3mock.NewMockEnvironmentAPIKeyCache(ctrl),
		experimentsCache:       cachev3mock.NewMockExperimentsCache(ctrl),
		autoOpsRulesCache:      cachev3mock.NewMockAutoOpsRulesCache(ctrl),
		invalidationPublisher:  publishermock.NewMockPublisher(ctrl),
	}
	p := NewCacheRefresher(
		m.featureClient,
		m.experimentClient,
		m.autoOpsClient,
		m.accountStorage,
		m.featuresCache,
		m.segmentUsersCache,
		m.environmentAPIKeyCache,
		m.experimentsCache,
		m.autoOpsRulesCache,
		m.invalidationPublisher,
		zap.NewNop(),
	).(*cacheRefresher)
	return p, m
}

func TestCacheRefresherHandleMessage(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc  string
		event *domaineventproto.Event
		setup func(m *cacheRefresherMocks, event *domaineventproto.Event)
	}{
		{
			desc: "feature event refreshes features cache and publishes invalidation",
			event: &domaineventproto.Event{
				Id:            "evt-feature",
				EntityType:    domaineventproto.Event_FEATURE,
				EntityId:      "feature-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_FEATURE_UPDATED,
			},
			setup: func(m *cacheRefresherMocks, _ *domaineventproto.Event) {
				m.featureClient.EXPECT().
					ListFeatures(gomock.Any(), gomock.Any()).
					Return(&featureproto.ListFeaturesResponse{
						Features: []*featureproto.Feature{
							{Id: "feature-id-1"},
						},
						Cursor: "",
					}, nil)
				m.featuresCache.EXPECT().
					Put(gomock.Any(), "env-1").
					Return(nil)
				m.invalidationPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			desc: "segment delete evicts L2 and publishes invalidation without calling ListSegmentUsers/GetSegment",
			event: &domaineventproto.Event{
				Id:            "evt-segment-delete",
				EntityType:    domaineventproto.Event_SEGMENT,
				EntityId:      "segment-id-deleted",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_SEGMENT_DELETED,
			},
			setup: func(m *cacheRefresherMocks, _ *domaineventproto.Event) {
				// Crucially: NO ListSegmentUsers / GetSegment calls — those would
				// return NotFound for a deleted segment and the refresh would be
				// dropped, leaving stale users in L2.
				m.segmentUsersCache.EXPECT().
					Evict("segment-id-deleted", "env-1").
					Return(nil)
				m.invalidationPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			desc: "segment event refreshes segment users cache and publishes invalidation",
			event: &domaineventproto.Event{
				Id:            "evt-segment",
				EntityType:    domaineventproto.Event_SEGMENT,
				EntityId:      "segment-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_SEGMENT_USER_ADDED,
			},
			setup: func(m *cacheRefresherMocks, event *domaineventproto.Event) {
				m.featureClient.EXPECT().
					ListSegmentUsers(gomock.Any(), gomock.Any()).
					Return(&featureproto.ListSegmentUsersResponse{
						Users: []*featureproto.SegmentUser{
							{Id: "user-1", SegmentId: "segment-id-1"},
						},
					}, nil)
				m.featureClient.EXPECT().
					GetSegment(gomock.Any(), gomock.Any()).
					Return(&featureproto.GetSegmentResponse{
						Segment: &featureproto.Segment{
							Id:        "segment-id-1",
							UpdatedAt: 12345,
						},
					}, nil)
				m.segmentUsersCache.EXPECT().
					Put(gomock.Any(), "env-1").
					Return(nil)
				m.invalidationPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			desc: "api key event refreshes environment api key cache and publishes invalidation",
			event: &domaineventproto.Event{
				Id:            "evt-apikey",
				EntityType:    domaineventproto.Event_APIKEY,
				EntityId:      "apikey-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_APIKEY_CHANGED,
				EntityData:    `{"api_key": "secret-123"}`,
			},
			setup: func(m *cacheRefresherMocks, event *domaineventproto.Event) {
				m.accountStorage.EXPECT().
					GetEnvironmentAPIKey(gomock.Any(), "secret-123").
					Return(&accountdomain.EnvironmentAPIKey{
						EnvironmentAPIKey: &accountproto.EnvironmentAPIKey{
							ApiKey: &accountproto.APIKey{Id: "apikey-id-1", ApiKey: "secret-123"},
						},
					}, nil)
				m.environmentAPIKeyCache.EXPECT().
					Put(gomock.Any()).
					Return(nil)
				m.invalidationPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			desc: "api key event refreshes both previous and current secrets",
			event: &domaineventproto.Event{
				Id:                 "evt-apikey-2",
				EntityType:         domaineventproto.Event_APIKEY,
				EntityId:           "apikey-id-1",
				EnvironmentId:      "env-1",
				Type:               domaineventproto.Event_APIKEY_CHANGED,
				EntityData:         `{"api_key": "new-secret"}`,
				PreviousEntityData: `{"api_key": "old-secret"}`,
			},
			setup: func(m *cacheRefresherMocks, event *domaineventproto.Event) {
				m.accountStorage.EXPECT().
					GetEnvironmentAPIKey(gomock.Any(), "old-secret").
					Return(&accountdomain.EnvironmentAPIKey{
						EnvironmentAPIKey: &accountproto.EnvironmentAPIKey{
							ApiKey: &accountproto.APIKey{Id: "apikey-id-1", ApiKey: "old-secret"},
						},
					}, nil)
				m.environmentAPIKeyCache.EXPECT().Put(gomock.Any()).Return(nil)
				m.accountStorage.EXPECT().
					GetEnvironmentAPIKey(gomock.Any(), "new-secret").
					Return(&accountdomain.EnvironmentAPIKey{
						EnvironmentAPIKey: &accountproto.EnvironmentAPIKey{
							ApiKey: &accountproto.APIKey{Id: "apikey-id-1", ApiKey: "new-secret"},
						},
					}, nil)
				m.environmentAPIKeyCache.EXPECT().Put(gomock.Any()).Return(nil)
				m.invalidationPublisher.EXPECT().
					Publish(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			desc: "api key event evicts and skips publish when key is missing in DB",
			event: &domaineventproto.Event{
				Id:            "evt-apikey-missing",
				EntityType:    domaineventproto.Event_APIKEY,
				EntityId:      "apikey-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_APIKEY_CHANGED,
				EntityData:    `{"api_key": "deleted-secret"}`,
			},
			setup: func(m *cacheRefresherMocks, event *domaineventproto.Event) {
				m.accountStorage.EXPECT().
					GetEnvironmentAPIKey(gomock.Any(), "deleted-secret").
					Return(nil, accstorage.ErrAPIKeyNotFound)
				m.environmentAPIKeyCache.EXPECT().Evict("deleted-secret").Return(nil)
				m.invalidationPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			desc: "experiment event refreshes experiments cache and publishes invalidation",
			event: &domaineventproto.Event{
				Id:            "evt-experiment",
				EntityType:    domaineventproto.Event_EXPERIMENT,
				EntityId:      "experiment-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_EXPERIMENT_CREATED,
			},
			setup: func(m *cacheRefresherMocks, _ *domaineventproto.Event) {
				m.experimentClient.EXPECT().
					ListExperiments(gomock.Any(), gomock.Any()).
					Return(&experimentproto.ListExperimentsResponse{
						Experiments: []*experimentproto.Experiment{
							{Id: "experiment-id-1"},
						},
					}, nil)
				m.experimentsCache.EXPECT().
					Put(gomock.Any(), "env-1").
					Return(nil)
				m.invalidationPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			desc: "auto ops rule event refreshes auto ops rules cache and publishes invalidation",
			event: &domaineventproto.Event{
				Id:            "evt-autoops",
				EntityType:    domaineventproto.Event_AUTOOPS_RULE,
				EntityId:      "autoops-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_AUTOOPS_RULE_CREATED,
			},
			setup: func(m *cacheRefresherMocks, _ *domaineventproto.Event) {
				m.autoOpsClient.EXPECT().
					ListAutoOpsRules(gomock.Any(), gomock.Any()).
					Return(&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{Id: "autoops-id-1"},
						},
					}, nil)
				m.autoOpsRulesCache.EXPECT().
					Put(gomock.Any(), "env-1").
					Return(nil)
				m.invalidationPublisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			desc: "unrelated entity type does not touch any cache",
			event: &domaineventproto.Event{
				Id:            "evt-account",
				EntityType:    domaineventproto.Event_ACCOUNT,
				EntityId:      "account-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_ACCOUNT_V2_CREATED,
			},
			setup: func(_ *cacheRefresherMocks, _ *domaineventproto.Event) {},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cr, mocks := newCacheRefresherWithMocks(t)
			p.setup(mocks, p.event)

			data, err := proto.Marshal(p.event)
			assert.NoError(t, err)

			acked := false
			nacked := false
			msg := &puller.Message{
				Data: data,
				Ack:  func() { acked = true },
				Nack: func() { nacked = true },
			}
			cr.handleMessage(context.Background(), msg)
			assert.True(t, acked, "message should be acked on success")
			assert.False(t, nacked, "message should not be nacked on success")
		})
	}
}

func TestCacheRefresherHandleMessageNackOnRepeatableError(t *testing.T) {
	t.Parallel()

	cr, mocks := newCacheRefresherWithMocks(t)

	mocks.featureClient.EXPECT().
		ListFeatures(gomock.Any(), gomock.Any()).
		Return(nil, context.DeadlineExceeded)

	data, err := proto.Marshal(&domaineventproto.Event{
		Id:            "evt-feature",
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
		Ack:  func() { acked = true },
		Nack: func() { nacked = true },
	}
	cr.handleMessage(context.Background(), msg)

	assert.False(t, acked)
	assert.True(t, nacked)
}

func TestCacheRefresherHandleMessageAcksOnPutFailureNonRepeatable(t *testing.T) {
	t.Parallel()

	cr, mocks := newCacheRefresherWithMocks(t)

	mocks.featureClient.EXPECT().
		ListFeatures(gomock.Any(), gomock.Any()).
		Return(&featureproto.ListFeaturesResponse{}, nil)
	mocks.featuresCache.EXPECT().
		Put(gomock.Any(), "env-1").
		Return(errors.New("redis: oom"))
	// non-repeatable error => Ack (drop the event)

	data, err := proto.Marshal(&domaineventproto.Event{
		Id:            "evt-feature",
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
		Ack:  func() { acked = true },
		Nack: func() { nacked = true },
	}
	cr.handleMessage(context.Background(), msg)

	assert.True(t, acked)
	assert.False(t, nacked)
}

// Regression test for the shutdown-mid-publish bug: if L2 was successfully
// refreshed but the publish step then fails because the worker context was
// canceled (graceful shutdown), the message must be NACKed so a healthy pod
// can re-publish the L1 invalidation. Acking would leave api pods serving
// stale L1 entries until TTL expiry.
func TestCacheRefresherHandleMessageNacksWhenPublishFailsWithContextCanceled(t *testing.T) {
	t.Parallel()

	cr, mocks := newCacheRefresherWithMocks(t)

	mocks.featureClient.EXPECT().
		ListFeatures(gomock.Any(), gomock.Any()).
		Return(&featureproto.ListFeaturesResponse{}, nil)
	mocks.featuresCache.EXPECT().
		Put(gomock.Any(), "env-1").
		Return(nil)
	mocks.invalidationPublisher.EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		Return(context.Canceled)

	data, err := proto.Marshal(&domaineventproto.Event{
		Id:            "evt-feature",
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
		Ack:  func() { acked = true },
		Nack: func() { nacked = true },
	}
	cr.handleMessage(context.Background(), msg)

	assert.False(t, acked, "must not ack when publish failed mid-flight")
	assert.True(t, nacked, "must nack so another pod re-publishes the invalidation")
}

func TestCacheRefresherWorkerIndexStability(t *testing.T) {
	t.Parallel()

	// Same env hashes to the same bucket every call.
	idx1 := workerIndex("env-abc", 16)
	idx2 := workerIndex("env-abc", 16)
	assert.Equal(t, idx1, idx2)
	assert.GreaterOrEqual(t, idx1, 0)
	assert.Less(t, idx1, 16)

	// Empty env is mapped, deterministically, to bucket 0.
	assert.Equal(t, 0, workerIndex("", 16))

	// n<=0 is also defensively handled.
	assert.Equal(t, 0, workerIndex("env-abc", 0))
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
			desc:     "context.Canceled is repeatable (shutdown mid-refresh; let another pod retry)",
			err:      context.Canceled,
			expected: true,
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
