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
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestExecuteProgressiveRolloutOperation(t *testing.T) {
	t.Parallel()

	// Helper to create a feature with variations
	createFeature := func(variationIDs ...string) *ftdomain.Feature {
		variations := make([]*featureproto.Variation, len(variationIDs))
		for i, id := range variationIDs {
			variations[i] = &featureproto.Variation{Id: id}
		}
		return &ftdomain.Feature{
			Feature: &featureproto.Feature{
				Variations: variations,
			},
		}
	}

	patterns := []struct {
		desc             string
		rolloutClause    interface{}
		rolloutType      autoopsproto.ProgressiveRollout_Type
		feature          *ftdomain.Feature
		scheduleID       string
		expectedStrategy *featureproto.Strategy
		expectedErr      error
	}{
		{
			desc: "new format: manual schedule with both control and target IDs",
			rolloutClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				ControlVariationId: "variation-a",
				TargetVariationId:  "variation-b",
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{ScheduleId: "schedule-1", Weight: 30000},
				},
			},
			rolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			feature:     createFeature("variation-a", "variation-b"),
			scheduleID:  "schedule-1",
			expectedStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{Variation: "variation-b", Weight: 30000},
						{Variation: "variation-a", Weight: 70000},
					},
				},
			},
		},
		{
			desc: "old format: manual schedule with only variation_id - infers control",
			rolloutClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				VariationId: "variation-b", // Old field - this is target
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{ScheduleId: "schedule-1", Weight: 20000},
				},
			},
			rolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			feature:     createFeature("variation-a", "variation-b"),
			scheduleID:  "schedule-1",
			expectedStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{Variation: "variation-b", Weight: 20000},
						{Variation: "variation-a", Weight: 80000}, // Inferred control
					},
				},
			},
		},
		{
			desc: "new format: template schedule with both control and target IDs",
			rolloutClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				ControlVariationId: "var-1",
				TargetVariationId:  "var-2",
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{ScheduleId: "sched-1", Weight: 50000},
				},
			},
			rolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			feature:     createFeature("var-1", "var-2"),
			scheduleID:  "sched-1",
			expectedStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{Variation: "var-2", Weight: 50000},
						{Variation: "var-1", Weight: 50000},
					},
				},
			},
		},
		{
			desc: "old format: template schedule with only variation_id - infers control",
			rolloutClause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				VariationId: "var-2",
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{ScheduleId: "sched-1", Weight: 10000},
				},
			},
			rolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			feature:     createFeature("var-1", "var-2"),
			scheduleID:  "sched-1",
			expectedStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{Variation: "var-2", Weight: 10000},
						{Variation: "var-1", Weight: 90000}, // Inferred control
					},
				},
			},
		},
		{
			desc: "full rollout: 100% to target variation",
			rolloutClause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				ControlVariationId: "control",
				TargetVariationId:  "target",
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{ScheduleId: "final-schedule", Weight: 100000},
				},
			},
			rolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			feature:     createFeature("control", "target"),
			scheduleID:  "final-schedule",
			expectedStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
				RolloutStrategy: &featureproto.RolloutStrategy{
					Variations: []*featureproto.RolloutStrategy_Variation{
						{Variation: "target", Weight: 100000},
						{Variation: "control", Weight: 0},
					},
				},
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// Create progressive rollout with the clause
			var clause *anypb.Any
			var err error
			switch c := p.rolloutClause.(type) {
			case *autoopsproto.ProgressiveRolloutManualScheduleClause:
				clause, err = ptypes.MarshalAny(c)
			case *autoopsproto.ProgressiveRolloutTemplateScheduleClause:
				clause, err = ptypes.MarshalAny(c)
			}
			require.NoError(t, err)

			progressiveRollout := &domain.ProgressiveRollout{
				ProgressiveRollout: &autoopsproto.ProgressiveRollout{
					Id:     "test-rollout",
					Type:   p.rolloutType,
					Clause: clause,
				},
			}

			// Execute the operation
			strategy, err := ExecuteProgressiveRolloutOperation(
				progressiveRollout,
				p.feature,
				p.scheduleID,
			)

			// Verify results
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				assert.Equal(t, p.expectedStrategy, strategy)
			}
		})
	}
}

func TestGetRolloutStrategyVariations(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc               string
		controlVariationID string
		targetVariationID  string
		targetWeight       int32
		expected           []*featureproto.RolloutStrategy_Variation
	}{
		{
			desc:               "success: weight is max",
			controlVariationID: "vid-2",
			targetVariationID:  "vid-1",
			targetWeight:       totalVariationWeight,
			expected: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: "vid-1",
					Weight:    totalVariationWeight,
				},
				{
					Variation: "vid-2",
					Weight:    0,
				},
			},
		},
		{
			desc:               "success: weight is not max",
			controlVariationID: "vid-1",
			targetVariationID:  "vid-2",
			targetWeight:       20,
			expected: []*featureproto.RolloutStrategy_Variation{
				{
					Variation: "vid-2",
					Weight:    20,
				},
				{
					Variation: "vid-1",
					Weight:    totalVariationWeight - 20,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getRolloutStrategyVariations(
				p.controlVariationID,
				p.targetVariationID,
				p.targetWeight,
			)
			assert.Equal(t, p.expected, actual)
		})
	}
}
