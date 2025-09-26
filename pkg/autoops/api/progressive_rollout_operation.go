// Copyright 2025 The Bucketeer Authors.
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
	"github.com/golang/protobuf/ptypes"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var errVariationNotFound = err.NewErrorNotFound(
	err.AutoopsPackageName,
	"a variation for a progressive rollout is not found",
	"variation_id",
)

const totalVariationWeight = int32(100000)

func ExecuteProgressiveRolloutOperation(
	progressiveRollout *domain.ProgressiveRollout,
	feature *ftdomain.Feature,
	scheduleID string,
) (*featureproto.Strategy, error) {
	var variationID string
	var weight int32
	switch progressiveRollout.Type {
	case autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE:
		c := &autoopsproto.ProgressiveRolloutManualScheduleClause{}
		if err := ptypes.UnmarshalAny(progressiveRollout.Clause, c); err != nil {
			return nil, err
		}
		variationID = c.VariationId
		var err error
		weight, err = getTargetWeight(c.Schedules, scheduleID)
		if err != nil {
			return nil, err
		}
	case autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE:
		c := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{}
		if err := ptypes.UnmarshalAny(progressiveRollout.Clause, c); err != nil {
			return nil, err
		}
		variationID = c.VariationId
		var err error
		weight, err = getTargetWeight(c.Schedules, scheduleID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, domain.ErrProgressiveRolloutInvalidType
	}
	return newRolloutStrategy(
		weight,
		feature,
		variationID,
	)
}

func getTargetWeight(
	schedules []*autoopsproto.ProgressiveRolloutSchedule,
	scheduleID string,
) (int32, error) {
	for _, s := range schedules {
		if s.ScheduleId == scheduleID {
			return s.Weight, nil
		}
	}
	return 0, domain.ErrProgressiveRolloutScheduleNotFound
}

func newRolloutStrategy(
	weight int32,
	feature *ftdomain.Feature,
	targetVariationID string,
) (*featureproto.Strategy, error) {
	variations, err := getRolloutStrategyVariations(feature, weight, targetVariationID)
	if err != nil {
		return nil, err
	}
	strategy := &featureproto.Strategy{
		Type: featureproto.Strategy_ROLLOUT,
		RolloutStrategy: &featureproto.RolloutStrategy{
			Variations: variations,
		},
	}
	return strategy, nil
}

func getRolloutStrategyVariations(
	feature *ftdomain.Feature,
	weight int32,
	targetVariationID string,
) ([]*featureproto.RolloutStrategy_Variation, error) {
	nonTargetVariationID, err := findNonTargetVariationID(feature, targetVariationID)
	if err != nil {
		return nil, err
	}
	return []*featureproto.RolloutStrategy_Variation{
		{
			Variation: targetVariationID,
			Weight:    weight,
		},
		{
			Variation: nonTargetVariationID,
			Weight:    totalVariationWeight - weight,
		},
	}, nil
}

func findNonTargetVariationID(
	feature *ftdomain.Feature,
	variationID string,
) (string, error) {
	for _, v := range feature.Variations {
		if v.Id != variationID {
			return v.Id, nil
		}
	}
	return "", errVariationNotFound
}
