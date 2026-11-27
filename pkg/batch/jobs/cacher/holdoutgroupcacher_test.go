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
//

package cacher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	envclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	holdoutmock "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2/mock"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestHoldoutGroupCacherRun(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	ctx := context.Background()

	t.Run("environment list fails", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("list"))
		err := cacher.Run(ctx)
		assert.EqualError(t, err, "list")
	})

	t.Run("missing group stores an empty snapshot", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(nil, v2hs.ErrHoldoutGroupNotFound)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetClearedHoldoutUpdatedAt(gomock.Any(), "env").
			Return(int64(0), nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				PutIfNewer(&ftproto.HoldoutGroup{}, "env", gomock.Any()).
				Return(true, nil)
		}
		expectMeasuredGroups(cacher, "env")
		assert.NoError(t, cacher.Run(ctx))
	})

	t.Run("stopped group keeps its stop time on the empty snapshot", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(nil, v2hs.ErrHoldoutGroupNotFound)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetClearedHoldoutUpdatedAt(gomock.Any(), "env").
			Return(int64(1700000000), nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				PutIfNewer(&ftproto.HoldoutGroup{UpdatedAt: 1700000000}, "env", gomock.Any()).
				Return(true, nil)
		}
		expectMeasuredGroups(cacher, "env")
		assert.NoError(t, cacher.Run(ctx))
	})

	t.Run("stores the live group", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		group := &ftproto.HoldoutGroup{Id: "hg", EnvironmentId: "env"}
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(&domain.HoldoutGroup{HoldoutGroup: group}, nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				PutIfNewer(group, "env", gomock.Any()).
				Return(true, nil)
		}
		expectMeasuredGroups(cacher, "env")
		assert.NoError(t, cacher.Run(ctx))
	})

	t.Run("storage error is reported after the remaining environments are refreshed", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		group := &ftproto.HoldoutGroup{Id: "hg", EnvironmentId: "env-2"}
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env-1"}, {Id: "env-2"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env-1").
			Return(nil, errors.New("db"))
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env-2").
			Return(&domain.HoldoutGroup{HoldoutGroup: group}, nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				PutIfNewer(group, "env-2", gomock.Any()).
				Return(true, nil)
		}
		// The measured groups of an environment are refreshed even when its
		// serving snapshot could not be read.
		expectMeasuredGroups(cacher, "env-1")
		expectMeasuredGroups(cacher, "env-2")
		err := cacher.Run(ctx)
		assert.EqualError(t, err, "db")
	})

	t.Run("a canceled context stops the run", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env-1"}, {Id: "env-2"}},
			}, nil)
		canceled, cancel := context.WithCancel(ctx)
		cancel()
		assert.ErrorIs(t, cacher.Run(canceled), context.Canceled)
	})

	t.Run("stale write is not an error", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		group := &ftproto.HoldoutGroup{Id: "hg", EnvironmentId: "env"}
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(&domain.HoldoutGroup{HoldoutGroup: group}, nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				PutIfNewer(group, "env", gomock.Any()).
				Return(false, nil)
		}
		expectMeasuredGroups(cacher, "env")
		assert.NoError(t, cacher.Run(ctx))
	})
}

func TestHoldoutGroupCacherRunMeasuredGroups(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	ctx := context.Background()

	expectServingSnapshot := func(cacher *holdoutGroupCacher, environmentID string) {
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), environmentID).
			Return(nil, v2hs.ErrHoldoutGroupNotFound)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetClearedHoldoutUpdatedAt(gomock.Any(), environmentID).
			Return(int64(0), nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				PutIfNewer(&ftproto.HoldoutGroup{}, environmentID, gomock.Any()).
				Return(true, nil)
		}
	}
	expectEnvironments := func(cacher *holdoutGroupCacher, ids ...string) {
		envs := make([]*envproto.EnvironmentV2, 0, len(ids))
		for _, id := range ids {
			envs = append(envs, &envproto.EnvironmentV2{Id: id})
		}
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{Environments: envs}, nil)
	}

	t.Run("stores the measured groups in every cache, in the order of the storage", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		expectEnvironments(cacher, "env")
		expectServingSnapshot(cacher, "env")
		stopped := &ftproto.HoldoutGroup{
			Id:           "stopped",
			ExperimentId: "e1",
			Status:       ftproto.HoldoutGroup_STOPPED,
			Features: []*ftproto.HoldoutGroupFeature{
				{Id: "r1", FeatureId: "flag", AddedAt: 10, RemovedAt: 20},
			},
		}
		running := &ftproto.HoldoutGroup{Id: "running", ExperimentId: "e2", Status: ftproto.HoldoutGroup_RUNNING}
		expectMeasuredGroups(cacher, "env", stopped, running)
		assert.NoError(t, cacher.Run(ctx))
	})

	t.Run("the list reaches back as far as stopped experiments still link events", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		expectEnvironments(cacher, "env")
		expectServingSnapshot(cacher, "env")
		before := time.Now()
		var endedSince int64
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			ListMeasuredGroups(gomock.Any(), "env", gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, since int64) ([]*domain.HoldoutGroup, error) {
				endedSince = since
				return nil, nil
			})
		generations := make([]int64, len(cacher.measuredGroupsCaches))
		for i, cache := range cacher.measuredGroupsCaches {
			cache.(*cachev3mock.MockHoldoutMeasuredGroupsCache).EXPECT().
				PutIfNewer(gomock.Any(), "env", gomock.Any()).
				DoAndReturn(func(_ *ftproto.HoldoutGroups, _ string, generation int64) (bool, error) {
					generations[i] = generation
					return true, nil
				})
		}
		require.NoError(t, cacher.Run(ctx))
		after := time.Now()
		assert.GreaterOrEqual(t, endedSince, before.Add(-experimentdomain.StoppedExperimentEventWindow).Unix())
		assert.LessOrEqual(t, endedSince, after.Add(-experimentdomain.StoppedExperimentEventWindow).Unix())
		// The generation is the time of the read, so that a list read
		// earlier never replaces a list read later.
		for _, generation := range generations {
			assert.GreaterOrEqual(t, generation, before.UnixNano())
			assert.LessOrEqual(t, generation, after.UnixNano())
			assert.Equal(t, generations[0], generation)
		}
	})

	t.Run("a failed read is reported, writes nothing, and does not stop the run", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		expectEnvironments(cacher, "env-1", "env-2")
		expectServingSnapshot(cacher, "env-1")
		expectServingSnapshot(cacher, "env-2")
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			ListMeasuredGroups(gomock.Any(), "env-1", gomock.Any()).
			Return(nil, errors.New("measured"))
		// No PutIfNewer is expected for env-1: gomock fails the test if the
		// list of an environment is replaced after a failed read.
		expectMeasuredGroups(cacher, "env-2")
		assert.EqualError(t, cacher.Run(ctx), "measured")
	})

	t.Run("a failed or refused cache write is not an error", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		expectEnvironments(cacher, "env")
		expectServingSnapshot(cacher, "env")
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			ListMeasuredGroups(gomock.Any(), "env", gomock.Any()).
			Return(nil, nil)
		cacher.measuredGroupsCaches[0].(*cachev3mock.MockHoldoutMeasuredGroupsCache).EXPECT().
			PutIfNewer(gomock.Any(), "env", gomock.Any()).
			Return(false, errors.New("redis"))
		cacher.measuredGroupsCaches[1].(*cachev3mock.MockHoldoutMeasuredGroupsCache).EXPECT().
			PutIfNewer(gomock.Any(), "env", gomock.Any()).
			Return(false, nil)
		assert.NoError(t, cacher.Run(ctx))
	})
}

func TestNewHoldoutGroupCacherBuildsBothCachesPerInstance(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	job := NewHoldoutGroupCacher(
		envclientmock.NewMockClient(controller),
		holdoutmock.NewMockHoldoutStorage(controller),
		[]cache.MultiGetCache{
			cachemock.NewMockMultiGetCache(controller),
			cachemock.NewMockMultiGetCache(controller),
		},
	)
	cacher, ok := job.(*holdoutGroupCacher)
	require.True(t, ok)
	assert.Len(t, cacher.caches, 2)
	assert.Len(t, cacher.measuredGroupsCaches, 2)
}

// expectMeasuredGroups expects one read of the measured groups of the
// environment and the list stored in every cache. Without groups the stored
// list is empty.
func expectMeasuredGroups(cacher *holdoutGroupCacher, environmentID string, groups ...*ftproto.HoldoutGroup) {
	stored := make([]*domain.HoldoutGroup, 0, len(groups))
	cached := &ftproto.HoldoutGroups{HoldoutGroups: make([]*ftproto.HoldoutGroup, 0, len(groups))}
	for _, group := range groups {
		stored = append(stored, &domain.HoldoutGroup{HoldoutGroup: group})
		cached.HoldoutGroups = append(cached.HoldoutGroups, group)
	}
	cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
		ListMeasuredGroups(gomock.Any(), environmentID, gomock.Any()).
		Return(stored, nil)
	for _, cache := range cacher.measuredGroupsCaches {
		cache.(*cachev3mock.MockHoldoutMeasuredGroupsCache).EXPECT().
			PutIfNewer(cached, environmentID, gomock.Any()).
			Return(true, nil)
	}
}

func newHoldoutGroupCacher(t *testing.T, controller *gomock.Controller) *holdoutGroupCacher {
	t.Helper()
	require.NotNil(t, controller)
	return &holdoutGroupCacher{
		environmentClient: envclientmock.NewMockClient(controller),
		storage:           holdoutmock.NewMockHoldoutStorage(controller),
		caches: []cachev3.HoldoutGroupCache{
			cachev3mock.NewMockHoldoutGroupCache(controller),
			cachev3mock.NewMockHoldoutGroupCache(controller),
		},
		measuredGroupsCaches: []cachev3.HoldoutMeasuredGroupsCache{
			cachev3mock.NewMockHoldoutMeasuredGroupsCache(controller),
			cachev3mock.NewMockHoldoutMeasuredGroupsCache(controller),
		},
		logger: zap.NewNop(),
	}
}
