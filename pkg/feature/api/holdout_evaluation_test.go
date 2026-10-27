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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	holdoutmock "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestEvaluateFeaturesAppliesCachedHoldout(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	holdoutCache := cachev3mock.NewMockHoldoutGroupCache(controller)
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
	holdoutCache.EXPECT().Get("env").Return(group, nil)
	service := &FeatureService{
		holdoutGroupCache: holdoutCache,
		logger:            zap.NewNop(),
	}
	evaluations, err := service.evaluateFeatures(
		context.Background(),
		[]*featureproto.Feature{disabledFlag()},
		&userproto.User{Id: "user-1"},
		"env",
		"",
	)
	require.NoError(t, err)
	require.Len(t, evaluations.Evaluations, 1)
	assert.Equal(t, featureproto.Reason_HOLDOUT, evaluations.Evaluations[0].Reason.Type)
	assert.Equal(t, "hg", evaluations.Evaluations[0].Reason.HoldoutGroupId)
	assert.Equal(t, "hold", evaluations.Evaluations[0].VariationId)
}

func TestEvaluateFeaturesCacheMissLeavesTheFlagOff(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	holdoutCache := cachev3mock.NewMockHoldoutGroupCache(controller)
	holdoutCache.EXPECT().Get("env").Return(nil, cache.ErrNotFound)
	service := &FeatureService{
		holdoutGroupCache: holdoutCache,
		logger:            zap.NewNop(),
	}
	evaluations, err := service.evaluateFeatures(
		context.Background(),
		[]*featureproto.Feature{disabledFlag()},
		&userproto.User{Id: "user-1"},
		"env",
		"",
	)
	require.NoError(t, err)
	require.Len(t, evaluations.Evaluations, 1)
	assert.Equal(t, featureproto.Reason_OFF_VARIATION, evaluations.Evaluations[0].Reason.Type)
	assert.Equal(t, "off", evaluations.Evaluations[0].VariationId)
}

func TestEvaluateFeaturesClearanceOmitsDisabledFlagWithoutOffVariation(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	holdoutCache := cachev3mock.NewMockHoldoutGroupCache(controller)
	holdoutCache.EXPECT().Get("env").Return(&featureproto.HoldoutGroup{UpdatedAt: 1700000000}, nil)
	service := &FeatureService{
		holdoutGroupCache: holdoutCache,
		logger:            zap.NewNop(),
	}
	evaluations, err := service.evaluateFeatures(
		context.Background(),
		[]*featureproto.Feature{{Id: "flag", Enabled: false}},
		&userproto.User{Id: "user-1"},
		"env",
		"",
	)
	require.NoError(t, err)
	assert.Empty(t, evaluations.Evaluations)
	assert.Equal(t, []string{"flag"}, evaluations.ArchivedFeatureIds)
}

func TestEvaluateFeaturesCacheErrorOmitsDisabledFlagWithoutOffVariation(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	holdoutCache := cachev3mock.NewMockHoldoutGroupCache(controller)
	holdoutCache.EXPECT().Get("env").Return(nil, errors.New("redis"))
	service := &FeatureService{
		holdoutGroupCache: holdoutCache,
		logger:            zap.NewNop(),
	}
	evaluations, err := service.evaluateFeatures(
		context.Background(),
		[]*featureproto.Feature{{Id: "flag", Enabled: false}},
		&userproto.User{Id: "user-1"},
		"env",
		"",
	)
	require.NoError(t, err)
	assert.Empty(t, evaluations.Evaluations)
	assert.Equal(t, []string{"flag"}, evaluations.ArchivedFeatureIds)
}

func TestRefreshHoldoutCache(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	ctx := context.Background()

	t.Run("nil cache is a no-op", func(t *testing.T) {
		service := &FeatureService{logger: zap.NewNop()}
		service.refreshHoldoutCache(ctx, "env")
	})
	t.Run("stores an empty snapshot when nothing is serving", func(t *testing.T) {
		cacheMock := cachev3mock.NewMockHoldoutGroupCache(controller)
		storage := holdoutmock.NewMockHoldoutStorage(controller)
		storage.EXPECT().GetLiveServingGroup(gomock.Any(), "env").Return(nil, v2hs.ErrHoldoutGroupNotFound)
		storage.EXPECT().GetClearedHoldoutUpdatedAt(gomock.Any(), "env").Return(int64(0), nil)
		cacheMock.EXPECT().PutIfNewer(&featureproto.HoldoutGroup{}, "env", gomock.Any()).Return(true, nil)
		service := &FeatureService{
			holdoutGroupCache: cacheMock,
			holdoutStorage:    storage,
			logger:            zap.NewNop(),
		}
		service.refreshHoldoutCache(ctx, "env")
	})
	t.Run("stamps the stop time on the empty snapshot", func(t *testing.T) {
		cacheMock := cachev3mock.NewMockHoldoutGroupCache(controller)
		storage := holdoutmock.NewMockHoldoutStorage(controller)
		storage.EXPECT().GetLiveServingGroup(gomock.Any(), "env").Return(nil, v2hs.ErrHoldoutGroupNotFound)
		storage.EXPECT().GetClearedHoldoutUpdatedAt(gomock.Any(), "env").Return(int64(1700000000), nil)
		cacheMock.EXPECT().
			PutIfNewer(&featureproto.HoldoutGroup{UpdatedAt: 1700000000}, "env", gomock.Any()).
			Return(true, nil)
		service := &FeatureService{
			holdoutGroupCache: cacheMock,
			holdoutStorage:    storage,
			logger:            zap.NewNop(),
		}
		service.refreshHoldoutCache(ctx, "env")
	})
	t.Run("writes a barrier when the clearance read fails", func(t *testing.T) {
		cacheMock := cachev3mock.NewMockHoldoutGroupCache(controller)
		storage := holdoutmock.NewMockHoldoutStorage(controller)
		storage.EXPECT().GetLiveServingGroup(gomock.Any(), "env").Return(nil, v2hs.ErrHoldoutGroupNotFound)
		storage.EXPECT().GetClearedHoldoutUpdatedAt(gomock.Any(), "env").Return(int64(0), errors.New("db"))
		cacheMock.EXPECT().PutIfNewer(gomock.Any(), "env", gomock.Any()).
			DoAndReturn(func(group *featureproto.HoldoutGroup, _ string, _ int64) (bool, error) {
				assert.Empty(t, group.Id)
				assert.Positive(t, group.UpdatedAt)
				return true, nil
			})
		service := &FeatureService{
			holdoutGroupCache: cacheMock,
			holdoutStorage:    storage,
			logger:            zap.NewNop(),
		}
		service.refreshHoldoutCache(ctx, "env")
	})
	t.Run("stores the live group", func(t *testing.T) {
		cacheMock := cachev3mock.NewMockHoldoutGroupCache(controller)
		storage := holdoutmock.NewMockHoldoutStorage(controller)
		group := &featureproto.HoldoutGroup{Id: "hg"}
		storage.EXPECT().GetLiveServingGroup(gomock.Any(), "env").
			Return(&domain.HoldoutGroup{HoldoutGroup: group}, nil)
		cacheMock.EXPECT().PutIfNewer(group, "env", gomock.Any()).Return(true, nil)
		service := &FeatureService{
			holdoutGroupCache: cacheMock,
			holdoutStorage:    storage,
			logger:            zap.NewNop(),
		}
		service.refreshHoldoutCache(ctx, "env")
	})
	t.Run("writes a barrier when the read fails", func(t *testing.T) {
		cacheMock := cachev3mock.NewMockHoldoutGroupCache(controller)
		storage := holdoutmock.NewMockHoldoutStorage(controller)
		storage.EXPECT().GetLiveServingGroup(gomock.Any(), "env").Return(nil, errors.New("db"))
		cacheMock.EXPECT().PutIfNewer(gomock.Any(), "env", gomock.Any()).
			DoAndReturn(func(group *featureproto.HoldoutGroup, _ string, _ int64) (bool, error) {
				assert.Empty(t, group.Id)
				assert.Positive(t, group.UpdatedAt)
				return true, nil
			})
		service := &FeatureService{
			holdoutGroupCache: cacheMock,
			holdoutStorage:    storage,
			logger:            zap.NewNop(),
		}
		service.refreshHoldoutCache(ctx, "env")
	})
}

func disabledFlag() *featureproto.Feature {
	return &featureproto.Feature{
		Id:           "flag",
		Enabled:      false,
		OffVariation: "off",
		Variations: []*featureproto.Variation{
			{Id: "hold", Value: "hold"},
			{Id: "off", Value: "off"},
		},
	}
}
