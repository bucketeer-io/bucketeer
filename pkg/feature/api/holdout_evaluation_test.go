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
		cacheMock.EXPECT().Put(&featureproto.HoldoutGroup{}, "env").Return(nil)
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
		cacheMock.EXPECT().Put(group, "env").Return(nil)
		service := &FeatureService{
			holdoutGroupCache: cacheMock,
			holdoutStorage:    storage,
			logger:            zap.NewNop(),
		}
		service.refreshHoldoutCache(ctx, "env")
	})
	t.Run("evicts when the read fails", func(t *testing.T) {
		cacheMock := cachev3mock.NewMockHoldoutGroupCache(controller)
		storage := holdoutmock.NewMockHoldoutStorage(controller)
		storage.EXPECT().GetLiveServingGroup(gomock.Any(), "env").Return(nil, errors.New("db"))
		cacheMock.EXPECT().Evict("env").Return(nil)
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
