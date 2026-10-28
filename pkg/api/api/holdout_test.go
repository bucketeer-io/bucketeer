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

package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestLoadLiveHoldoutGroup(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	live := &featureproto.HoldoutGroup{
		Id:     "hg",
		Status: featureproto.HoldoutGroup_ANALYSIS,
		Features: []*featureproto.HoldoutGroupFeature{{
			FeatureId:          "flag",
			HoldoutVariationId: "hold",
		}},
	}

	t.Run("memory hit", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env").Return(live, nil)
		got := loadLiveHoldoutGroup(context.Background(), memory, nil, "env", nil, zap.NewNop())
		assert.Equal(t, "hg", got.Id)
	})
	t.Run("redis hit is copied into memory", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env").Return(live, nil)
		memory.EXPECT().Put(live, "env").Return(nil)
		got := loadLiveHoldoutGroup(context.Background(), memory, redis, "env", nil, zap.NewNop())
		assert.Equal(t, live, got)
	})
	t.Run("empty snapshot stays a snapshot", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		empty := &featureproto.HoldoutGroup{}
		memory.EXPECT().Get("env").Return(empty, nil)
		got := loadLiveHoldoutGroup(context.Background(), memory, nil, "env", nil, zap.NewNop())
		assert.Empty(t, got.GetId())
		assert.Equal(t, int64(0), got.GetUpdatedAt())
	})
	t.Run("clearance snapshot keeps the stop revision", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		cleared := &featureproto.HoldoutGroup{UpdatedAt: 1700000000}
		memory.EXPECT().Get("env").Return(cleared, nil)
		assert.Equal(t, cleared, loadLiveHoldoutGroup(context.Background(), memory, nil, "env", nil, zap.NewNop()))
	})
	t.Run("both misses without a fallback are unavailable", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env-miss").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env-miss").Return(nil, cache.ErrNotFound)
		assert.Equal(t,
			evaluation.UnavailableHoldoutSnapshot(),
			loadLiveHoldoutGroup(context.Background(), memory, redis, "env-miss", nil, zap.NewNop()),
		)
	})
	t.Run("a read error without a fallback is unavailable", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env").Return(nil, errors.New("redis"))
		assert.Equal(t, evaluation.UnavailableHoldoutSnapshot(), loadLiveHoldoutGroup(context.Background(), memory, redis, "env", nil, zap.NewNop()))
	})
	t.Run("a redis miss falls back to the feature service and fills memory", func(t *testing.T) {
		// Before the cacher's first run, or after a flush, the key is missing
		// for every environment. That must not degrade evaluation.
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env-fallback").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env-fallback").Return(nil, cache.ErrNotFound)
		memory.EXPECT().Put(live, "env-fallback").Return(nil)
		calls := 0
		fallback := func(_ context.Context, epoch uint64) (*featureproto.HoldoutGroup, error) {
			calls++
			assert.Equal(t, loadHoldoutFillEpoch("env-fallback"), epoch, "the fallback receives the epoch the fill is committed under")
			return live, nil
		}
		got := loadLiveHoldoutGroup(context.Background(), memory, redis, "env-fallback", fallback, zap.NewNop())
		assert.Equal(t, live, got)
		assert.Equal(t, 1, calls)
	})
	t.Run("a redis miss with no group falls back to a clearance, not unavailable", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env-clear").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env-clear").Return(nil, cache.ErrNotFound)
		cleared := evaluation.ClearanceSnapshot(0)
		memory.EXPECT().Put(cleared, "env-clear").Return(nil)
		fallback := func(context.Context, uint64) (*featureproto.HoldoutGroup, error) { return cleared, nil }
		got := loadLiveHoldoutGroup(context.Background(), memory, redis, "env-clear", fallback, zap.NewNop())
		assert.Equal(t, cleared, got)
		assert.False(t, evaluation.HoldoutStale(holdoutGroups(got), 1), "an environment without a group must not force updates")
	})
	t.Run("a redis error also uses the fallback", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env-err").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env-err").Return(nil, errors.New("redis"))
		memory.EXPECT().Put(live, "env-err").Return(nil)
		fallback := func(context.Context, uint64) (*featureproto.HoldoutGroup, error) { return live, nil }
		assert.Equal(t, live, loadLiveHoldoutGroup(context.Background(), memory, redis, "env-err", fallback, zap.NewNop()))
	})
	t.Run("a failing fallback is unavailable", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env-down").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env-down").Return(nil, cache.ErrNotFound)
		fallback := func(context.Context, uint64) (*featureproto.HoldoutGroup, error) { return nil, errors.New("grpc") }
		assert.Equal(t,
			evaluation.UnavailableHoldoutSnapshot(),
			loadLiveHoldoutGroup(context.Background(), memory, redis, "env-down", fallback, zap.NewNop()),
		)
	})
}

func TestHoldoutServingSnapshot(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	live := &featureproto.HoldoutGroup{Id: "hg", Status: featureproto.HoldoutGroup_RUNNING, UpdatedAt: 30}
	isLiveQuery := func(req *featureproto.ListHoldoutGroupsRequest) bool {
		return len(req.Statuses) == 2 &&
			req.Statuses[0] == featureproto.HoldoutGroup_RUNNING &&
			req.Statuses[1] == featureproto.HoldoutGroup_ANALYSIS &&
			req.Archived != nil && !req.Archived.Value &&
			req.PageSize == 1 &&
			req.OrderBy == featureproto.ListHoldoutGroupsRequest_UPDATED_AT &&
			req.OrderDirection == featureproto.ListHoldoutGroupsRequest_DESC
	}
	isStoppedQuery := func(req *featureproto.ListHoldoutGroupsRequest) bool {
		return len(req.Statuses) == 2 &&
			req.Statuses[0] == featureproto.HoldoutGroup_STOPPED &&
			req.Statuses[1] == featureproto.HoldoutGroup_FORCE_STOPPED &&
			req.Archived == nil && // archived stopped groups still carry the stop revision
			req.PageSize == 1 &&
			req.OrderBy == featureproto.ListHoldoutGroupsRequest_UPDATED_AT &&
			req.OrderDirection == featureproto.ListHoldoutGroupsRequest_DESC
	}

	t.Run("serving group", func(t *testing.T) {
		client := featureclientmock.NewMockClient(controller)
		client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, req *featureproto.ListHoldoutGroupsRequest, _ ...grpc.CallOption) (*featureproto.ListHoldoutGroupsResponse, error) {
				require.True(t, isLiveQuery(req))
				return &featureproto.ListHoldoutGroupsResponse{HoldoutGroups: []*featureproto.HoldoutGroup{live}}, nil
			})
		got, err := HoldoutServingSnapshot(context.Background(), client, "env")
		require.NoError(t, err)
		assert.Equal(t, live, got)
	})
	t.Run("no serving group keeps the newest stop revision", func(t *testing.T) {
		client := featureclientmock.NewMockClient(controller)
		gomock.InOrder(
			client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, req *featureproto.ListHoldoutGroupsRequest, _ ...grpc.CallOption) (*featureproto.ListHoldoutGroupsResponse, error) {
					require.True(t, isLiveQuery(req))
					return &featureproto.ListHoldoutGroupsResponse{}, nil
				}),
			client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, req *featureproto.ListHoldoutGroupsRequest, _ ...grpc.CallOption) (*featureproto.ListHoldoutGroupsResponse, error) {
					require.True(t, isStoppedQuery(req))
					return &featureproto.ListHoldoutGroupsResponse{HoldoutGroups: []*featureproto.HoldoutGroup{
						{Id: "old", Status: featureproto.HoldoutGroup_STOPPED, UpdatedAt: 20},
					}}, nil
				}),
		)
		got, err := HoldoutServingSnapshot(context.Background(), client, "env")
		require.NoError(t, err)
		assert.Equal(t, evaluation.ClearanceSnapshot(20), got)
	})
	t.Run("never had a group is the zero clearance", func(t *testing.T) {
		client := featureclientmock.NewMockClient(controller)
		client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).
			Return(&featureproto.ListHoldoutGroupsResponse{}, nil).Times(2)
		got, err := HoldoutServingSnapshot(context.Background(), client, "env")
		require.NoError(t, err)
		assert.Equal(t, evaluation.ClearanceSnapshot(0), got)
		assert.Empty(t, evaluation.EvaluationHoldout(got).GetId())
	})
	t.Run("service error is returned", func(t *testing.T) {
		client := featureclientmock.NewMockClient(controller)
		client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).Return(nil, errors.New("grpc"))
		_, err := HoldoutServingSnapshot(context.Background(), client, "env")
		assert.Error(t, err)
	})
}

func TestHoldoutSnapshotFromService(t *testing.T) {
	t.Parallel()
	assert.Nil(t, holdoutSnapshotFromService(&singleflight.Group{}, nil, "env"), "no client disables the fallback")

	live := &featureproto.HoldoutGroup{Id: "hg", Status: featureproto.HoldoutGroup_RUNNING}
	response := &featureproto.ListHoldoutGroupsResponse{HoldoutGroups: []*featureproto.HoldoutGroup{live}}
	const wait = 5 * time.Second

	t.Run("concurrent callers in one epoch all receive the snapshot", func(t *testing.T) {
		controller := gomock.NewController(t)
		client := featureclientmock.NewMockClient(controller)
		// Whether a caller joins the in-flight read or starts the next one
		// depends on scheduling, so only the results are asserted; the
		// coalescing itself is singleflight's contract.
		release := make(chan struct{})
		client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).
			DoAndReturn(func(context.Context, *featureproto.ListHoldoutGroupsRequest, ...grpc.CallOption) (*featureproto.ListHoldoutGroupsResponse, error) {
				<-release
				return response, nil
			}).MinTimes(1).MaxTimes(2)
		fallback := holdoutSnapshotFromService(&singleflight.Group{}, client, "env")
		results := make(chan *featureproto.HoldoutGroup, 2)
		for i := 0; i < 2; i++ {
			go func() {
				got, err := fallback(context.Background(), 1)
				assert.NoError(t, err)
				results <- got
			}()
		}
		close(release)
		for i := 0; i < 2; i++ {
			select {
			case got := <-results:
				assert.Equal(t, live, got)
			case <-time.After(wait):
				t.Fatal("a caller never returned")
			}
		}
	})

	t.Run("a later epoch does not join an older read", func(t *testing.T) {
		controller := gomock.NewController(t)
		client := featureclientmock.NewMockClient(controller)
		entered := make(chan struct{}, 2)
		release := make(chan struct{})
		client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).
			DoAndReturn(func(context.Context, *featureproto.ListHoldoutGroupsRequest, ...grpc.CallOption) (*featureproto.ListHoldoutGroupsResponse, error) {
				entered <- struct{}{}
				<-release
				return response, nil
			}).Times(2)
		fallback := holdoutSnapshotFromService(&singleflight.Group{}, client, "env")
		results := make(chan error, 2)
		for _, epoch := range []uint64{1, 2} {
			go func() {
				_, err := fallback(context.Background(), epoch)
				results <- err
			}()
		}
		// Both epochs must reach the service while the first read is still
		// blocked: an invalidation between them means the older snapshot is
		// not safe to share.
		for i := 0; i < 2; i++ {
			select {
			case <-entered:
			case <-time.After(wait):
				t.Fatalf("only %d read(s) reached the feature service", i)
			}
		}
		close(release)
		for i := 0; i < 2; i++ {
			select {
			case err := <-results:
				assert.NoError(t, err)
			case <-time.After(wait):
				t.Fatal("a caller never returned")
			}
		}
	})

	t.Run("a canceled caller stops waiting without aborting the read", func(t *testing.T) {
		controller := gomock.NewController(t)
		client := featureclientmock.NewMockClient(controller)
		entered := make(chan struct{}, 1)
		release := make(chan struct{})
		readCtxErr := make(chan error, 1)
		client.EXPECT().ListHoldoutGroups(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ *featureproto.ListHoldoutGroupsRequest, _ ...grpc.CallOption) (*featureproto.ListHoldoutGroupsResponse, error) {
				entered <- struct{}{}
				<-release
				readCtxErr <- ctx.Err()
				return response, nil
			}).Times(1)
		fallback := holdoutSnapshotFromService(&singleflight.Group{}, client, "env")
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() {
			_, err := fallback(ctx, 1)
			done <- err
		}()
		select {
		case <-entered:
		case <-time.After(wait):
			t.Fatal("the caller never reached the feature service")
		}
		cancel()
		select {
		case err := <-done:
			// A closed stream or canceled request returns at once, flagged as
			// a caller cancellation rather than a service failure.
			assert.ErrorIs(t, err, errCallerCanceled)
		case <-time.After(wait):
			t.Fatal("the canceled caller kept waiting for the read")
		}
		// The shared read was not canceled by the caller leaving: it runs to
		// completion with a live context.
		close(release)
		select {
		case err := <-readCtxErr:
			assert.NoError(t, err)
		case <-time.After(wait):
			t.Fatal("the shared read never completed")
		}
	})
}
