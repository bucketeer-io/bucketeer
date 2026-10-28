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

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
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

func TestLoadLiveHoldoutGroupDoesNotRestoreAStoppedSnapshot(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	environmentID := "env-fill-race"
	memory := cachev3.NewHoldoutGroupCache(cachev3.NewInMemoryCache(), time.Minute)
	redis := cachev3mock.NewMockHoldoutGroupCache(controller)
	running := &featureproto.HoldoutGroup{
		Id:        "hg",
		Status:    featureproto.HoldoutGroup_RUNNING,
		UpdatedAt: 10,
	}
	cleared := &featureproto.HoldoutGroup{UpdatedAt: 20}
	redis.EXPECT().Get(environmentID).DoAndReturn(func(string) (*featureproto.HoldoutGroup, error) {
		invalidateHoldoutFill(environmentID, func() {
			require.NoError(t, memory.Put(cleared, environmentID))
		})
		return running, nil
	})
	redis.EXPECT().Get(environmentID).Return(cleared, nil)

	got := loadLiveHoldoutGroup(context.Background(), memory, redis, environmentID, nil, zap.NewNop())
	// The in-flight read observed the stop, so it does not answer with the
	// group it already fetched. The retry reads the clearance.
	assert.Empty(t, got.GetId())
	assert.Equal(t, int64(20), got.GetUpdatedAt())

	cached, err := memory.Get(environmentID)
	require.NoError(t, err)
	assert.Empty(t, cached.Id)
	assert.Equal(t, int64(20), cached.UpdatedAt)
}

func TestLoadLiveHoldoutGroupFallbackDoesNotRestoreAStoppedSnapshot(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	environmentID := "env-fallback-race"
	memory := cachev3.NewHoldoutGroupCache(cachev3.NewInMemoryCache(), time.Minute)
	redis := cachev3mock.NewMockHoldoutGroupCache(controller)
	redis.EXPECT().Get(environmentID).Return(nil, cache.ErrNotFound).Times(2)
	running := &featureproto.HoldoutGroup{
		Id:        "hg",
		Status:    featureproto.HoldoutGroup_RUNNING,
		UpdatedAt: 10,
	}
	cleared := evaluation.ClearanceSnapshot(20)
	// The first feature-service read is answered with the pre-stop group, but
	// the group is stopped while it is in flight: the invalidator bumps the
	// epoch and writes the clearance to memory. The read that follows the
	// invalidation sees the stop.
	var epochs []uint64
	fallback := func(_ context.Context, epoch uint64) (*featureproto.HoldoutGroup, error) {
		epochs = append(epochs, epoch)
		if len(epochs) == 1 {
			invalidateHoldoutFill(environmentID, func() {
				require.NoError(t, memory.Put(cleared, environmentID))
			})
			return running, nil
		}
		return cleared, nil
	}

	got := loadLiveHoldoutGroup(context.Background(), memory, redis, environmentID, fallback, zap.NewNop())
	// The stale read is not committed and not served. The retry runs under
	// the new epoch, and the clearance is what the caller gets.
	require.Len(t, epochs, 2)
	assert.Equal(t, epochs[0]+1, epochs[1], "the retry must carry the epoch bumped by the invalidation")
	assert.Empty(t, got.GetId())
	assert.Equal(t, int64(20), got.GetUpdatedAt())

	cached, err := memory.Get(environmentID)
	require.NoError(t, err)
	assert.Empty(t, cached.Id, "the stopped group must not be restored to memory")
	assert.Equal(t, int64(20), cached.UpdatedAt)
}

func TestFeatureFillDoesNotRestorePreHoldoutSnapshot(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	environmentID := "env-feature-fill"
	memory := cachev3.NewFeaturesCache(cachev3.NewInMemoryCache(), time.Minute)
	redis := cachev3mock.NewMockFeaturesCache(controller)
	before := &featureproto.Features{Features: []*featureproto.Feature{
		{Id: "normal", Enabled: true, OffVariation: "on"},
	}}
	after := &featureproto.Features{Features: []*featureproto.Feature{
		{Id: "normal", Enabled: true, OffVariation: "on"},
		{Id: "targeted", Enabled: false},
	}}
	redis.EXPECT().Get(environmentID).DoAndReturn(func(string) (*featureproto.Features, error) {
		invalidateHoldoutFill(environmentID, func() {
			require.NoError(t, memory.Evict(environmentID))
		})
		return before, nil
	})
	redis.EXPECT().Get(environmentID).Return(after, nil)

	got, err := FillFeatureSnapshot(environmentID, func() (*featureproto.Features, error) {
		return redis.Get(environmentID)
	}, func(features *featureproto.Features) {
		require.NoError(t, memory.Put(features, environmentID))
	})
	require.NoError(t, err)
	require.Len(t, got.Features, 2)
	assert.Equal(t, "targeted", got.Features[1].Id)

	cached, err := memory.Get(environmentID)
	require.NoError(t, err)
	require.Len(t, cached.Features, 2)
	assert.Equal(t, "targeted", cached.Features[1].Id)
}

func TestRemovedIDsForTag(t *testing.T) {
	t.Parallel()
	service := &grpcGatewayService{}
	features := []*featureproto.Feature{
		{Id: "android-flag", Tags: []string{"android"}},
		{Id: "ios-flag", Tags: []string{"ios"}},
		{Id: "both", Tags: []string{"android", "ios"}},
	}
	removed := []string{"android-flag", "ios-flag", "both", "missing"}
	assert.Equal(t, []string{"android-flag", "both"}, service.removedIDsForTag(features, removed, "android"))
	assert.Equal(t, removed, service.removedIDsForTag(features, removed, ""))
	assert.Empty(t, service.removedIDsForTag(features, nil, "android"))
}

func TestGatewayEvaluateFeaturesHoldout(t *testing.T) {
	t.Parallel()
	service := &gatewayService{logger: zap.NewNop()}
	group := &featureproto.HoldoutGroup{
		Id:             "hg",
		EnvironmentId:  "env",
		Status:         featureproto.HoldoutGroup_RUNNING,
		AssignmentType: featureproto.HoldoutGroup_HASH,
		SamplingSeed:   "seed",
		Percentage:     100,
		Features: []*featureproto.HoldoutGroupFeature{{
			FeatureId:          "flag",
			HoldoutVariationId: "hold",
		}},
	}
	evaluations, err := service.evaluateFeatures(
		context.Background(),
		&userproto.User{Id: "user-1"},
		[]*featureproto.Feature{{
			Id:           "flag",
			Enabled:      false,
			OffVariation: "off",
			Variations: []*featureproto.Variation{
				{Id: "hold", Value: "hold"},
				{Id: "off", Value: "off"},
			},
		}},
		"env",
		"",
		group,
	)
	require.NoError(t, err)
	require.Len(t, evaluations.Evaluations, 1)
	assert.Equal(t, featureproto.Reason_HOLDOUT, evaluations.Evaluations[0].Reason.Type)
	assert.Equal(t, "hold", evaluations.Evaluations[0].VariationId)

	without, err := service.evaluateFeatures(
		context.Background(),
		&userproto.User{Id: "user-1"},
		[]*featureproto.Feature{{
			Id:           "flag",
			Enabled:      false,
			OffVariation: "off",
			Variations:   []*featureproto.Variation{{Id: "off", Value: "off"}},
		}},
		"env",
		"",
		nil,
	)
	require.NoError(t, err)
	assert.Equal(t, featureproto.Reason_OFF_VARIATION, without.Evaluations[0].Reason.Type)
}

func TestGatewayEvaluateDisabledHoldoutTarget(t *testing.T) {
	t.Parallel()
	service := &gatewayService{logger: zap.NewNop()}
	disabled := &featureproto.Feature{
		Id:      "flag",
		Enabled: false,
		Variations: []*featureproto.Variation{
			{Id: "hold", Value: "hold"},
			{Id: "on", Value: "on"},
		},
		DefaultStrategy: &featureproto.Strategy{
			Type:          featureproto.Strategy_FIXED,
			FixedStrategy: &featureproto.FixedStrategy{Variation: "on"},
		},
	}
	sibling := &featureproto.Feature{
		Id:           "sibling",
		Enabled:      true,
		OffVariation: "off",
		Variations: []*featureproto.Variation{
			{Id: "off", Value: "off"},
			{Id: "on", Value: "on"},
		},
		DefaultStrategy: &featureproto.Strategy{
			Type:          featureproto.Strategy_FIXED,
			FixedStrategy: &featureproto.FixedStrategy{Variation: "on"},
		},
	}
	group := &featureproto.HoldoutGroup{
		Id:             "hg",
		EnvironmentId:  "env",
		Status:         featureproto.HoldoutGroup_RUNNING,
		AssignmentType: featureproto.HoldoutGroup_HASH,
		SamplingSeed:   "seed",
		Percentage:     100,
		Features: []*featureproto.HoldoutGroupFeature{{
			FeatureId:          "flag",
			HoldoutVariationId: "hold",
		}},
	}
	member, err := service.evaluateFeatures(
		context.Background(),
		&userproto.User{Id: "user-1"},
		[]*featureproto.Feature{disabled, sibling},
		"env",
		"",
		group,
	)
	require.NoError(t, err)
	require.Len(t, member.Evaluations, 2)
	byID := map[string]*featureproto.Evaluation{}
	for _, evaluation := range member.Evaluations {
		byID[evaluation.FeatureId] = evaluation
	}
	assert.Equal(t, featureproto.Reason_HOLDOUT, byID["flag"].Reason.Type)
	assert.Equal(t, "hold", byID["flag"].VariationId)
	assert.Equal(t, "sibling", byID["sibling"].FeatureId)

	group.Percentage = 0
	nonMember, err := service.evaluateFeatures(
		context.Background(),
		&userproto.User{Id: "user-1"},
		[]*featureproto.Feature{disabled, sibling},
		"env",
		"",
		group,
	)
	require.NoError(t, err)
	require.Len(t, nonMember.Evaluations, 1)
	assert.Equal(t, "sibling", nonMember.Evaluations[0].FeatureId)
}

func TestListFeaturesKeepsDisabledHoldoutTarget(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	client := featureclientmock.NewMockClient(controller)
	client.EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(&featureproto.ListFeaturesResponse{
		Features: []*featureproto.Feature{
			{Id: "kept", Enabled: false},
			{Id: "dropped", Enabled: false},
			{Id: "normal", Enabled: true, OffVariation: "on"},
		},
	}, nil)
	got, err := ListFeatures(
		context.Background(),
		client,
		"env",
		map[string]struct{}{"kept": {}},
	)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "kept", got[0].Id)
	assert.Equal(t, "normal", got[1].Id)
}

func TestListFeaturesKeepsDependencyOfStoppedTarget(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	client := featureclientmock.NewMockClient(controller)
	client.EXPECT().ListFeatures(gomock.Any(), gomock.Any()).Return(&featureproto.ListFeaturesResponse{
		Features: []*featureproto.Feature{
			{
				Id:           "dependent",
				Enabled:      true,
				OffVariation: "off",
				Prerequisites: []*featureproto.Prerequisite{
					{FeatureId: "flag", VariationId: "on"},
				},
			},
			{Id: "flag", Enabled: false},
			{Id: "dropped", Enabled: false},
		},
	}, nil)
	got, err := ListFeatures(context.Background(), client, "env", nil)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "dependent", got[0].Id)
	assert.Equal(t, "flag", got[1].Id)
}

func TestCountHoldoutEvaluationRequiresTargets(t *testing.T) {
	t.Parallel()
	metric := holdoutEvaluationsCounter.WithLabelValues("holdout-unit")
	before := testutil.ToFloat64(metric)
	countHoldoutEvaluation(nil, "holdout-unit")
	countHoldoutEvaluation(&featureproto.HoldoutGroup{Id: "hg"}, "holdout-unit")
	assert.Equal(t, before, testutil.ToFloat64(metric))
	countHoldoutEvaluation(&featureproto.HoldoutGroup{
		Id:       "hg",
		Features: []*featureproto.HoldoutGroupFeature{{FeatureId: "flag"}},
	}, "holdout-unit")
	assert.Equal(t, before+1, testutil.ToFloat64(metric))
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
