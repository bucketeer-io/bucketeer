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
	"github.com/golang/protobuf/ptypes"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const totalVariationWeight = int32(100000)

func ExecuteProgressiveRolloutOperation(
	progressiveRollout *domain.ProgressiveRollout,
	feature *ftdomain.Feature,
	scheduleID string,
) (*featureproto.Strategy, error) {
	// Extract control and target variation IDs
	controlVariationID, err := progressiveRollout.GetControlVariationID(feature)
	if err != nil {
		return nil, err
	}
	targetVariationID, err := progressiveRollout.GetTargetVariationID()
	if err != nil {
		return nil, err
	}

	// Get weight for this schedule
	var weight int32
	switch progressiveRollout.Type {
	case autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE:
		c := &autoopsproto.ProgressiveRolloutManualScheduleClause{}
		if err := ptypes.UnmarshalAny(progressiveRollout.Clause, c); err != nil {
			return nil, err
		}
		weight, err = getTargetWeight(c.Schedules, scheduleID)
		if err != nil {
			return nil, err
		}
	case autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE:
		c := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{}
		if err := ptypes.UnmarshalAny(progressiveRollout.Clause, c); err != nil {
			return nil, err
		}
		weight, err = getTargetWeight(c.Schedules, scheduleID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, domain.ErrProgressiveRolloutInvalidType
	}

	return newRolloutStrategy(
		controlVariationID,
		targetVariationID,
		weight,
		feature,
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
	controlVariationID string,
	targetVariationID string,
	weight int32,
	feature *ftdomain.Feature,
) (*featureproto.Strategy, error) {
	variations := getRolloutStrategyVariations(controlVariationID, targetVariationID, weight, feature)
	strategy := &featureproto.Strategy{
		Type: featureproto.Strategy_ROLLOUT,
		RolloutStrategy: &featureproto.RolloutStrategy{
			Variations: variations,
		},
	}
	return strategy, nil
}

func getRolloutStrategyVariations(
	controlVariationID string,
	targetVariationID string,
	weight int32,
	feature *ftdomain.Feature,
) []*featureproto.RolloutStrategy_Variation {
	// Create variations for all feature variations
	// Control and target get their calculated weights, all others get 0
	variations := make([]*featureproto.RolloutStrategy_Variation, 0, len(feature.Variations))

	for _, v := range feature.Variations {
		var varWeight int32
		switch v.Id {
		case targetVariationID:
			varWeight = weight
		case controlVariationID:
			varWeight = totalVariationWeight - weight
		default:
			// All other variations are reset to 0
			varWeight = 0
		}

		variations = append(variations, &featureproto.RolloutStrategy_Variation{
			Variation: v.Id,
			Weight:    varWeight,
		})
	}

	return variations
}
