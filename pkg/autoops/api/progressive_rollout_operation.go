// Copyright 2023 The Bucketeer Authors.
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

	"github.com/golang/protobuf/ptypes"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var errVariationNotFound = errors.New("autoops: a variation for a progressive rollout is not found")

const totalVariationWeight = int32(100000)

func ExecuteProgressiveRolloutOperation(
	ctx context.Context,
	progressiveRollout *domain.ProgressiveRollout,
	featureClient featureclient.Client,
	scheduleID, environmentNamespace string,
) error {
	switch progressiveRollout.Type {
	case autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE:
		c := &autoopsproto.ProgressiveRolloutManualScheduleClause{}
		if err := ptypes.UnmarshalAny(progressiveRollout.Clause, c); err != nil {
			return err
		}
		s, err := getTargetSchedule(c.Schedules, scheduleID)
		if err != nil {
			return err
		}
		if err := updateRolloutStrategy(
			ctx,
			s,
			featureClient,
			c.VariationId,
			progressiveRollout.FeatureId,
			environmentNamespace,
		); err != nil {
			return err
		}
	case autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE:
		c := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{}
		if err := ptypes.UnmarshalAny(progressiveRollout.Clause, c); err != nil {
			return err
		}
		s, err := getTargetSchedule(c.Schedules, scheduleID)
		if err != nil {
			return err
		}
		if err := updateRolloutStrategy(
			ctx,
			s,
			featureClient,
			c.VariationId,
			progressiveRollout.FeatureId,
			environmentNamespace,
		); err != nil {
			return err
		}
	default:
		return domain.ErrProgressiveRolloutInvalidType
	}
	return nil
}

func getTargetSchedule(
	schedules []*autoopsproto.ProgressiveRolloutSchedule,
	scheduleID string,
) (*autoopsproto.ProgressiveRolloutSchedule, error) {
	for _, s := range schedules {
		if s.ScheduleId == scheduleID {
			return s, nil
		}
	}
	return nil, domain.ErrProgressiveRolloutScheduleNotFound
}

func updateRolloutStrategy(
	ctx context.Context,
	schedule *autoopsproto.ProgressiveRolloutSchedule,
	featureClient featureclient.Client,
	targetVariationID, featureID, environmentNamespace string,
) error {
	f, err := fetchFeature(ctx, featureClient, featureID, environmentNamespace)
	if err != nil {
		return err
	}
	if err := updateFeatureTargeting(
		ctx,
		schedule,
		featureClient,
		f,
		targetVariationID,
		environmentNamespace,
	); err != nil {
		return err
	}
	return nil
}

func fetchFeature(
	ctx context.Context,
	featureClient featureclient.Client,
	featureID, environmentNamespace string,
) (*featureproto.Feature, error) {
	resp, err := featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentNamespace: environmentNamespace,
		Id:                   featureID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Feature, nil
}

func updateFeatureTargeting(
	ctx context.Context,
	schedule *autoopsproto.ProgressiveRolloutSchedule,
	featureClient featureclient.Client,
	feature *featureproto.Feature,
	targetVariationID, environmentNamespace string,
) error {
	cmds, err := getResetFeatureCmds(feature)
	if err != nil {
		return err
	}
	c, err := getNewRuleCmd(feature, schedule, targetVariationID)
	if err != nil {
		return err
	}
	cmds = append(cmds, c)
	_, err = featureClient.UpdateFeatureTargeting(
		ctx,
		&featureproto.UpdateFeatureTargetingRequest{
			EnvironmentNamespace: environmentNamespace,
			Id:                   feature.Id,
			Commands:             cmds,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func getResetFeatureCmds(
	feature *featureproto.Feature,
) ([]*featureproto.Command, error) {
	// In resetting feature, we don't delete users in feature.Targets and prerequisite.
	cmds := make([]*featureproto.DeleteRuleCommand, 0, len(feature.Rules))
	for _, r := range feature.Rules {
		c := &featureproto.DeleteRuleCommand{
			Id: r.Id,
		}
		cmds = append(cmds, c)
	}
	featureCmds := make([]*featureproto.Command, 0, len(cmds))
	for _, c := range cmds {
		ac, err := ptypes.MarshalAny(c)
		if err != nil {
			return nil, err
		}
		featureCmds = append(featureCmds, &featureproto.Command{
			Command: ac,
		})
	}
	return featureCmds, nil
}

func getNewRuleCmd(
	feature *featureproto.Feature,
	schedule *autoopsproto.ProgressiveRolloutSchedule,
	targetVariationID string,
) (*featureproto.Command, error) {
	variations, err := getRolloutStrategyVariations(feature, schedule, targetVariationID)
	if err != nil {
		return nil, err
	}
	c := &featureproto.ChangeDefaultStrategyCommand{
		Strategy: &featureproto.Strategy{
			Type: featureproto.Strategy_ROLLOUT,
			RolloutStrategy: &featureproto.RolloutStrategy{
				Variations: variations,
			},
		},
	}
	ac, err := ptypes.MarshalAny(c)
	if err != nil {
		return nil, err
	}
	return &featureproto.Command{
		Command: ac,
	}, nil
}

func getRolloutStrategyVariations(
	feature *featureproto.Feature,
	schedule *autoopsproto.ProgressiveRolloutSchedule,
	targetVariationID string,
) ([]*featureproto.RolloutStrategy_Variation, error) {
	nonTargetVariationID, err := findNonTargetVariationID(feature, targetVariationID)
	if err != nil {
		return nil, err
	}
	return []*featureproto.RolloutStrategy_Variation{
		{
			Variation: targetVariationID,
			Weight:    schedule.Weight,
		},
		{
			Variation: nonTargetVariationID,
			Weight:    totalVariationWeight - schedule.Weight,
		},
	}, nil
}

func findNonTargetVariationID(
	feature *featureproto.Feature,
	variationID string,
) (string, error) {
	for _, v := range feature.Variations {
		if v.Id != variationID {
			return v.Id, nil
		}
	}
	return "", errVariationNotFound
}
