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
		got := loadLiveHoldoutGroup(memory, nil, "env", zap.NewNop())
		assert.Equal(t, "hg", got.Id)
	})
	t.Run("redis hit is copied into memory", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env").Return(live, nil)
		memory.EXPECT().Put(live, "env").Return(nil)
		got := loadLiveHoldoutGroup(memory, redis, "env", zap.NewNop())
		assert.Equal(t, live, got)
	})
	t.Run("empty snapshot is not a live group", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env").Return(&featureproto.HoldoutGroup{}, nil)
		assert.Nil(t, loadLiveHoldoutGroup(memory, nil, "env", zap.NewNop()))
	})
	t.Run("clearance snapshot keeps the stop revision", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		cleared := &featureproto.HoldoutGroup{UpdatedAt: 1700000000}
		memory.EXPECT().Get("env").Return(cleared, nil)
		assert.Equal(t, cleared, loadLiveHoldoutGroup(memory, nil, "env", zap.NewNop()))
	})
	t.Run("both misses stay off the database", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env-miss").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env-miss").Return(nil, cache.ErrNotFound)
		assert.Equal(t,
			evaluation.UnavailableHoldoutSnapshot(),
			loadLiveHoldoutGroup(memory, redis, "env-miss", zap.NewNop()),
		)
	})
	t.Run("a read error keeps disabled flags omitted", func(t *testing.T) {
		memory := cachev3mock.NewMockHoldoutGroupCache(controller)
		redis := cachev3mock.NewMockHoldoutGroupCache(controller)
		memory.EXPECT().Get("env").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env").Return(nil, errors.New("redis"))
		assert.Equal(t, evaluation.UnavailableHoldoutSnapshot(), loadLiveHoldoutGroup(memory, redis, "env", zap.NewNop()))
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

	got := loadLiveHoldoutGroup(memory, redis, environmentID, zap.NewNop())
	// The in-flight read observed the stop, so it does not answer with the
	// group it already fetched. The retry reads the clearance.
	assert.Empty(t, got.GetId())
	assert.Equal(t, int64(20), got.GetUpdatedAt())

	cached, err := memory.Get(environmentID)
	require.NoError(t, err)
	assert.Empty(t, cached.Id)
	assert.Equal(t, int64(20), cached.UpdatedAt)
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
