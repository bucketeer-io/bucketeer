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

package domain

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoiface"

	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewProgressiveRollout(t *testing.T) {
	t.Parallel()
	aor := createProgressiveRollout(t)
	assert.IsType(t, &ProgressiveRollout{}, aor)
	assert.Equal(t, "feature-id", aor.FeatureId)
	assert.Equal(t, autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE, aor.Type)
	assert.NotNil(t, aor.Clause)
	assert.NotZero(t, aor.CreatedAt)
	assert.NotZero(t, aor.UpdatedAt)
}

func createProgressiveRollout(t *testing.T) *ProgressiveRollout {
	aor, err := NewProgressiveRollout(
		"feature-id",
		nil,
		&autoopsproto.ProgressiveRolloutTemplateScheduleClause{
			Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
				{
					ScheduleId: "schedule-id-0",
					ExecuteAt:  time.Now().Unix(),
					Weight:     0,
				},
				{
					ScheduleId: "schedule-id-1",
					ExecuteAt:  time.Now().AddDate(1, 0, 0).Unix(),
					Weight:     20,
				},
				{
					ScheduleId: "schedule-id-2",
					ExecuteAt:  time.Now().AddDate(2, 0, 0).Unix(),
					Weight:     40,
				},
				{
					ScheduleId: "schedule-id-3",
					ExecuteAt:  time.Now().AddDate(3, 0, 0).Unix(),
					Weight:     60,
				},
				{
					ScheduleId: "schedule-id-4",
					ExecuteAt:  time.Now().AddDate(4, 0, 0).Unix(),
					Weight:     80,
				},
				{
					ScheduleId: "schedule-id-5",
					ExecuteAt:  time.Now().AddDate(5, 0, 0).Unix(),
					Weight:     100,
				},
			},
			Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
			Increments:  20,
			VariationId: "vid-1",
		},
	)
	require.NoError(t, err)
	return aor
}

func TestAlreadyTriggered(t *testing.T) {
	patterns := []struct {
		desc                   string
		progressiveRolloutType autoopsproto.ProgressiveRollout_Type
		clause                 protoiface.MessageV1
		targetScheduleID       string
		expected               bool
		expectedErr            error
	}{
		{
			desc:                   "err: template ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: manual ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: ErrProgressiveRolloutInvalidType",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_Type(10),
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutInvalidType,
		},
		{
			desc:                   "success: false",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-2",
			expected:         false,
		},
		{
			desc:                   "success: true",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId:  "sid-2",
						TriggeredAt: time.Now().AddDate(0, -1, 0).Unix(),
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-2",
			expected:         true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createProgressiveRollout(t)
			ac, err := ptypes.MarshalAny(p.clause)
			assert.NoError(t, err)
			s.Clause = ac
			s.Type = p.progressiveRolloutType
			triggered, err := s.AlreadyTriggered(p.targetScheduleID)
			assert.Equal(t, triggered, p.expected)
			assert.Equal(t, err, p.expectedErr)
		})
	}
}

func TestProgressiveRolloutSetTriggeredAt(t *testing.T) {
	patterns := []struct {
		desc                   string
		progressiveRolloutType autoopsproto.ProgressiveRollout_Type
		clause                 protoiface.MessageV1
		targetScheduleID       string
		expectedErr            error
		expectedStatus         autoopsproto.ProgressiveRollout_Status
	}{
		{
			desc:                   "err: template ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: manual ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: ErrProgressiveRolloutInvalidType",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_Type(10),
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutInvalidType,
		},
		{
			desc:                   "success",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-2",
			expectedStatus:   autoopsproto.ProgressiveRollout_RUNNING,
		},
		{
			desc:                   "success last schedule is executed",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-4",
			expectedStatus:   autoopsproto.ProgressiveRollout_FINISHED,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createProgressiveRollout(t)
			ac, err := ptypes.MarshalAny(p.clause)
			assert.NoError(t, err)
			s.Clause = ac
			s.Type = p.progressiveRolloutType
			err = s.SetTriggeredAt(p.targetScheduleID)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				c, err := unmarshalProgressiveRolloutManualClause(s.Clause)
				assert.NoError(t, err)
				s, err := findTargetSchedule(c.Schedules, p.targetScheduleID)
				assert.NoError(t, err)
				assert.NotZero(t, s.TriggeredAt)
			}
			assert.Equal(t, p.expectedStatus, s.Status)
		})
	}
}

func TestAddManualScheduleClause(t *testing.T) {
	patterns := []struct {
		desc   string
		clause *autoopsproto.ProgressiveRolloutManualScheduleClause
	}{
		{
			desc: "success",
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-1",
						Weight:     10,
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pro := createProgressiveRollout(t)
			pro.Clause = nil
			assert.Nil(t, pro.Clause)
			pro.addManualScheduleClause(p.clause)
			assert.NotNil(t, pro.Clause)
		})
	}
}

func TestAddTemplateScheduleClause(t *testing.T) {
	patterns := []struct {
		desc   string
		clause *autoopsproto.ProgressiveRolloutTemplateScheduleClause
	}{
		{
			desc: "success",
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-1",
						Weight:     10,
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pro := createProgressiveRollout(t)
			pro.Clause = nil
			assert.Nil(t, pro.Clause)
			pro.addTemplatelScheduleClause(p.clause)
			assert.NotNil(t, pro.Clause)
		})
	}
}

func TestExtractSchedules(t *testing.T) {
	p := createProgressiveRollout(t)
	actual, err := p.ExtractSchedules()
	assert.NoError(t, err)
	assert.Len(t, actual, 6)
	assert.Equal(t, actual[1].Weight, int32(20))
	assert.Equal(t, actual[5].Weight, int32(100))
}

func TestStop(t *testing.T) {
	patterns := []struct {
		desc     string
		input    autoopsproto.ProgressiveRollout_StoppedBy
		expected error
	}{
		{
			desc:     "err: stopped by is required",
			input:    autoopsproto.ProgressiveRollout_UNKNOWN,
			expected: ErrProgressiveRolloutStoopedByRequired,
		},
		{
			desc:     "success: by user",
			input:    autoopsproto.ProgressiveRollout_USER,
			expected: nil,
		},
		{
			desc:     "success: by schedule",
			input:    autoopsproto.ProgressiveRollout_OPS_SCHEDULE,
			expected: nil,
		},
		{
			desc:     "success: by kill switch",
			input:    autoopsproto.ProgressiveRollout_OPS_KILL_SWITCH,
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pr := createProgressiveRollout(t)
			err := pr.Stop(p.input)
			if err != nil {
				assert.Equal(t, p.expected, err, p.desc)
				assert.Equal(t, autoopsproto.ProgressiveRollout_WAITING, pr.Status, p.desc)
				assert.Equal(t, autoopsproto.ProgressiveRollout_UNKNOWN, pr.StoppedBy, p.desc)
				assert.Zero(t, pr.StoppedAt, p.desc)
				assert.NotZero(t, pr.UpdatedAt, p.desc)
			} else {
				assert.Equal(t, p.expected, err, p.desc)
				assert.Equal(t, autoopsproto.ProgressiveRollout_STOPPED, pr.Status, p.desc)
				assert.Equal(t, p.input, pr.StoppedBy, p.desc)
				assert.NotZero(t, pr.StoppedAt, p.desc)
				assert.True(t, pr.UpdatedAt > time.Now().Add(time.Second*-2).Unix(), p.desc)
			}
		})
	}
}

func TestGetControlVariationID(t *testing.T) {
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
		desc                   string
		progressiveRolloutType autoopsproto.ProgressiveRollout_Type
		clause                 protoiface.MessageV1
		feature                *ftdomain.Feature
		expected               string
		expectedErr            error
	}{
		{
			desc:                   "manual: new format with both control and target IDs",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules:          []*autoopsproto.ProgressiveRolloutSchedule{},
				ControlVariationId: "control-var-id",
				TargetVariationId:  "target-var-id",
			},
			feature:  createFeature("control-var-id", "target-var-id"),
			expected: "control-var-id",
		},
		{
			desc:                   "manual: old format with variation_id only - infer control from 2 variations",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules:   []*autoopsproto.ProgressiveRolloutSchedule{},
				VariationId: "target-var-id", // Old field - this is the target
			},
			feature:  createFeature("control-var-id", "target-var-id"),
			expected: "control-var-id", // Should infer the other variation
		},
		{
			desc:                   "manual: old format with variation_id but feature has 3 variations - cannot infer",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules:   []*autoopsproto.ProgressiveRolloutSchedule{},
				VariationId: "target-var-id",
			},
			feature:     createFeature("var-1", "var-2", "var-3"),
			expected:    "",
			expectedErr: ErrProgressiveRolloutInvalidVariationCount,
		},
		{
			desc:                   "manual: no variation IDs set",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{},
			},
			feature:  createFeature("var-1", "var-2"),
			expected: "",
		},
		{
			desc:                   "template: new format with both control and target IDs",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules:          []*autoopsproto.ProgressiveRolloutSchedule{},
				ControlVariationId: "control-var-id",
				TargetVariationId:  "target-var-id",
			},
			feature:  createFeature("control-var-id", "target-var-id"),
			expected: "control-var-id",
		},
		{
			desc:                   "template: old format with variation_id only - infer control from 2 variations",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules:   []*autoopsproto.ProgressiveRolloutSchedule{},
				VariationId: "target-var-id",
			},
			feature:  createFeature("control-var-id", "target-var-id"),
			expected: "control-var-id",
		},
		{
			desc:                   "template: old format with variation_id but feature has 3 variations - cannot infer",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules:   []*autoopsproto.ProgressiveRolloutSchedule{},
				VariationId: "target-var-id",
			},
			feature:     createFeature("var-1", "var-2", "var-3"),
			expected:    "",
			expectedErr: ErrProgressiveRolloutInvalidVariationCount,
		},
		{
			desc:                   "invalid type",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_Type(99),
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{},
			},
			feature:     createFeature("var-1", "var-2"),
			expected:    "",
			expectedErr: ErrProgressiveRolloutInvalidType,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pr := createProgressiveRollout(t)
			ac, err := ptypes.MarshalAny(p.clause)
			require.NoError(t, err)
			pr.Clause = ac
			pr.Type = p.progressiveRolloutType

			actual, err := pr.GetControlVariationID(p.feature)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestInferControlVariationID(t *testing.T) {
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
		desc              string
		variations        []*featureproto.Variation
		targetVariationID string
		expectedControlID string
		expectedErr       error
	}{
		{
			desc:              "2 variations - returns the other variation",
			variations:        createFeature("variation-a", "variation-b").Variations,
			targetVariationID: "variation-b",
			expectedControlID: "variation-a",
		},
		{
			desc:              "2 variations - returns the other variation (reversed)",
			variations:        createFeature("variation-a", "variation-b").Variations,
			targetVariationID: "variation-a",
			expectedControlID: "variation-b",
		},
		{
			desc:              "3 variations - returns error",
			variations:        createFeature("var-1", "var-2", "var-3").Variations,
			targetVariationID: "var-2",
			expectedControlID: "",
			expectedErr:       ErrProgressiveRolloutInvalidVariationCount,
		},
		{
			desc:              "1 variation - returns error",
			variations:        createFeature("var-1").Variations,
			targetVariationID: "var-1",
			expectedControlID: "",
			expectedErr:       ErrProgressiveRolloutInvalidVariationCount,
		},
		{
			desc:              "0 variations - returns error",
			variations:        createFeature().Variations,
			targetVariationID: "var-1",
			expectedControlID: "",
			expectedErr:       ErrProgressiveRolloutInvalidVariationCount,
		},
		{
			desc:              "4 variations - returns error",
			variations:        createFeature("var-1", "var-2", "var-3", "var-4").Variations,
			targetVariationID: "var-2",
			expectedControlID: "",
			expectedErr:       ErrProgressiveRolloutInvalidVariationCount,
		},
		{
			desc: "2 identical variations (pathological) - returns error",
			variations: []*featureproto.Variation{
				{Id: "same-id"},
				{Id: "same-id"},
			},
			targetVariationID: "same-id",
			expectedControlID: "",
			expectedErr:       ErrProgressiveRolloutControlVariationNotFound,
		},
		{
			desc:              "nil variations - returns error",
			variations:        nil,
			targetVariationID: "var-1",
			expectedControlID: "",
			expectedErr:       ErrProgressiveRolloutInvalidVariationCount,
		},
		{
			desc:              "empty targetVariationID - returns empty",
			variations:        createFeature("var-1", "var-2").Variations,
			targetVariationID: "",
			expectedControlID: "",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := inferControlVariationID(p.variations, p.targetVariationID)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedControlID, actual)
		})
	}
}
