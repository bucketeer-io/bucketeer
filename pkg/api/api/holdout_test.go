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
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
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
		memory.EXPECT().Get("env").Return(nil, cache.ErrNotFound)
		redis.EXPECT().Get("env").Return(nil, cache.ErrNotFound)
		assert.Nil(t, loadLiveHoldoutGroup(memory, redis, "env", zap.NewNop()))
	})
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
