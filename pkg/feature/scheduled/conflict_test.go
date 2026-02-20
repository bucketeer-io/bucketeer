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

package scheduled

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestFindDeletedReferencesNeededByPayload(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		earlier  *featureproto.ScheduledChangePayload
		later    *featureproto.ScheduledChangePayload
		expected []string
	}{
		{
			desc:     "both nil",
			earlier:  nil,
			later:    nil,
			expected: nil,
		},
		{
			desc: "no deleted references",
			earlier: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			later: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			expected: nil,
		},
		{
			desc: "earlier deletes variation used in later update",
			earlier: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			later: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			expected: []string{"variation var-1"},
		},
		{
			desc: "earlier deletes variation used in later off_variation",
			earlier: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			later: &featureproto.ScheduledChangePayload{
				OffVariation: wrapperspb.String("var-1"),
			},
			expected: []string{"off_variation var-1"},
		},
		{
			desc: "earlier deletes rule used in later update",
			earlier: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{ChangeType: featureproto.ChangeType_DELETE, Rule: &featureproto.Rule{Id: "rule-1"}},
				},
			},
			later: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Rule: &featureproto.Rule{Id: "rule-1"}},
				},
			},
			expected: []string{"rule rule-1"},
		},
		{
			desc: "earlier deletes variation used in later CREATE: no conflict",
			earlier: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			later: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_CREATE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			expected: nil,
		},
		{
			desc: "earlier deletes variation used in later default strategy fixed",
			earlier: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			later: &featureproto.ScheduledChangePayload{
				DefaultStrategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: "var-1"},
				},
			},
			expected: []string{"default_strategy.variation"},
		},
		{
			desc: "earlier deletes variation used in later target",
			earlier: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			later: &featureproto.ScheduledChangePayload{
				TargetChanges: []*featureproto.TargetChange{
					{ChangeType: featureproto.ChangeType_CREATE, Target: &featureproto.Target{Variation: "var-1", Users: []string{"user-1"}}},
				},
			},
			expected: []string{"target.variation var-1"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			result := findDeletedReferencesNeededByPayload(p.earlier, p.later)
			assert.Equal(t, p.expected, result)
		})
	}
}

func TestValidatePayloadReferences(t *testing.T) {
	t.Parallel()
	defaultFlag := &featureproto.Feature{
		Id:      "feature-id",
		Version: 1,
		Variations: []*featureproto.Variation{
			{Id: "var-1", Name: "A", Value: "true"},
			{Id: "var-2", Name: "B", Value: "false"},
		},
		Rules: []*featureproto.Rule{
			{Id: "rule-1"},
		},
	}

	patterns := []struct {
		desc        string
		flag        *featureproto.Feature
		payload     *featureproto.ScheduledChangePayload
		expectedLen int
		checkField  string
	}{
		{
			desc:        "nil payload",
			flag:        defaultFlag,
			payload:     nil,
			expectedLen: 0,
		},
		{
			desc:        "nil flag",
			flag:        nil,
			payload:     &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
			expectedLen: 0,
		},
		{
			desc: "valid variation update",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-1"}},
				},
			},
			expectedLen: 0,
		},
		{
			desc: "invalid variation update",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-999"}},
				},
			},
			expectedLen: 1,
			checkField:  "variations",
		},
		{
			desc: "valid variation create: no reference check needed",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_CREATE, Variation: &featureproto.Variation{Id: "var-new"}},
				},
			},
			expectedLen: 0,
		},
		{
			desc: "invalid rule delete",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{ChangeType: featureproto.ChangeType_DELETE, Rule: &featureproto.Rule{Id: "rule-999"}},
				},
			},
			expectedLen: 1,
			checkField:  "rules",
		},
		{
			desc: "valid rule delete",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{ChangeType: featureproto.ChangeType_DELETE, Rule: &featureproto.Rule{Id: "rule-1"}},
				},
			},
			expectedLen: 0,
		},
		{
			desc: "invalid off_variation",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				OffVariation: wrapperspb.String("var-nonexistent"),
			},
			expectedLen: 1,
			checkField:  "off_variation",
		},
		{
			desc: "valid off_variation",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				OffVariation: wrapperspb.String("var-1"),
			},
			expectedLen: 0,
		},
		{
			desc: "invalid default strategy fixed",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				DefaultStrategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: "var-gone"},
				},
			},
			expectedLen: 1,
			checkField:  "default_strategy",
		},
		{
			desc: "invalid default strategy rollout",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				DefaultStrategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{Variation: "var-1", Weight: 50000},
							{Variation: "var-gone", Weight: 50000},
						},
					},
				},
			},
			expectedLen: 1,
			checkField:  "default_strategy",
		},
		{
			desc: "invalid target variation reference",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				TargetChanges: []*featureproto.TargetChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Target: &featureproto.Target{
							Variation: "var-gone",
							Users:     []string{"user-1"},
						},
					},
				},
			},
			expectedLen: 1,
			checkField:  "targets",
		},
		{
			desc: "valid target variation reference",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				TargetChanges: []*featureproto.TargetChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Target: &featureproto.Target{
							Variation: "var-1",
							Users:     []string{"user-1"},
						},
					},
				},
			},
			expectedLen: 0,
		},
		{
			desc: "invalid rule strategy fixed variation reference",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-new",
							Strategy: &featureproto.Strategy{
								Type: featureproto.Strategy_FIXED,
								FixedStrategy: &featureproto.FixedStrategy{
									Variation: "var-gone",
								},
							},
						},
					},
				},
			},
			expectedLen: 1,
			checkField:  "rules.strategy",
		},
		{
			desc: "invalid rule strategy rollout variation reference",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-new",
							Strategy: &featureproto.Strategy{
								Type: featureproto.Strategy_ROLLOUT,
								RolloutStrategy: &featureproto.RolloutStrategy{
									Variations: []*featureproto.RolloutStrategy_Variation{
										{Variation: "var-1", Weight: 50000},
										{Variation: "var-gone", Weight: 50000},
									},
								},
							},
						},
					},
				},
			},
			expectedLen: 1,
			checkField:  "rules.strategy",
		},
		{
			desc: "valid rule strategy variation references",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-new",
							Strategy: &featureproto.Strategy{
								Type: featureproto.Strategy_FIXED,
								FixedStrategy: &featureproto.FixedStrategy{
									Variation: "var-1",
								},
							},
						},
					},
				},
			},
			expectedLen: 0,
		},
		{
			desc: "no changes: no conflicts",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(true),
			},
			expectedLen: 0,
		},
		{
			desc: "payload with only ResetSamplingSeed: no references to validate",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				ResetSamplingSeed: true,
			},
			expectedLen: 0,
		},
		{
			desc: "multiple targets: one valid, one invalid",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				TargetChanges: []*featureproto.TargetChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Target:     &featureproto.Target{Variation: "var-1", Users: []string{"user-1"}},
					},
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Target:     &featureproto.Target{Variation: "var-gone", Users: []string{"user-2"}},
					},
				},
			},
			expectedLen: 1,
			checkField:  "targets",
		},
		{
			desc: "rollout strategy: one of many variations deleted",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				DefaultStrategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
					RolloutStrategy: &featureproto.RolloutStrategy{
						Variations: []*featureproto.RolloutStrategy_Variation{
							{Variation: "var-1", Weight: 30000},
							{Variation: "var-2", Weight: 30000},
							{Variation: "var-deleted", Weight: 40000},
						},
					},
				},
			},
			expectedLen: 1,
			checkField:  "default_strategy",
		},
		{
			desc: "multiple invalid references in same payload",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-gone-1"}},
				},
				RuleChanges: []*featureproto.RuleChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Rule: &featureproto.Rule{Id: "rule-gone"}},
				},
				OffVariation: wrapperspb.String("var-gone-2"),
			},
			expectedLen: 3,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			now := time.Now().Unix()
			conflicts := validatePayloadReferences(p.flag, p.payload, now)
			assert.Len(t, conflicts, p.expectedLen)
			if p.expectedLen > 0 && p.checkField != "" {
				assert.Equal(t, p.checkField, conflicts[0].ConflictingField)
				assert.Equal(t, featureproto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE, conflicts[0].Type)
			}
		})
	}
}

func TestConflictDetector_DetectConflictsOnCreate(t *testing.T) {
	t.Parallel()
	defaultFlag := &featureproto.Feature{
		Id:      "feature-id",
		Version: 1,
		Variations: []*featureproto.Variation{
			{Id: "var-1", Name: "A", Value: "true"},
			{Id: "var-2", Name: "B", Value: "false"},
		},
		Rules: []*featureproto.Rule{
			{Id: "rule-1"},
		},
	}

	patterns := []struct {
		desc               string
		flag               *featureproto.Feature
		payload            *featureproto.ScheduledChangePayload
		scheduledAt        int64
		excludeScheduleID  string
		existingSchedules  []*featureproto.ScheduledFlagChange
		expectedLen        int
		expectedTypes      []featureproto.ScheduledChangeConflict_ConflictType
		expectedFieldCheck func(t *testing.T, conflicts []*featureproto.ScheduledChangeConflict)
	}{
		{
			desc: "no conflicts: no existing schedules",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(true),
			},
			scheduledAt:       time.Now().Add(2 * time.Hour).Unix(),
			existingSchedules: []*featureproto.ScheduledFlagChange{},
			expectedLen:       0,
		},
		{
			desc: "no conflict: same field at different times is valid",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(true),
			},
			scheduledAt: time.Now().Add(2 * time.Hour).Unix(),
			existingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:          "sfc-1",
					FeatureId:   "feature-id",
					ScheduledAt: time.Now().Add(time.Hour).Unix(),
					Payload:     &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(false)},
					Status:      featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				},
			},
			expectedLen: 0, // Scheduling the same field is NOT a conflict
		},
		{
			desc: "exclude self: no conflict when own ID is excluded",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(true),
			},
			scheduledAt:       time.Now().Add(2 * time.Hour).Unix(),
			excludeScheduleID: "sfc-1",
			existingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:          "sfc-1",
					FeatureId:   "feature-id",
					ScheduledAt: time.Now().Add(time.Hour).Unix(),
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
						},
					},
					Status: featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				},
			},
			expectedLen: 0,
		},
		{
			desc: "invalid variation reference",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-999", Name: "X", Value: "x"}},
				},
			},
			scheduledAt:       time.Now().Add(2 * time.Hour).Unix(),
			existingSchedules: []*featureproto.ScheduledFlagChange{},
			expectedLen:       1,
			expectedTypes:     []featureproto.ScheduledChangeConflict_ConflictType{featureproto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE},
			expectedFieldCheck: func(t *testing.T, conflicts []*featureproto.ScheduledChangeConflict) {
				t.Helper()
				assert.Contains(t, conflicts[0].Description, "var-999")
				assert.Equal(t, "variations", conflicts[0].ConflictingField)
			},
		},
		{
			desc: "invalid rule reference",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{ChangeType: featureproto.ChangeType_DELETE, Rule: &featureproto.Rule{Id: "rule-999"}},
				},
			},
			scheduledAt:       time.Now().Add(2 * time.Hour).Unix(),
			existingSchedules: []*featureproto.ScheduledFlagChange{},
			expectedLen:       1,
			expectedTypes:     []featureproto.ScheduledChangeConflict_ConflictType{featureproto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE},
			expectedFieldCheck: func(t *testing.T, conflicts []*featureproto.ScheduledChangeConflict) {
				t.Helper()
				assert.Contains(t, conflicts[0].Description, "rule-999")
				assert.Equal(t, "rules", conflicts[0].ConflictingField)
			},
		},
		{
			desc: "dependency missing: earlier schedule deletes variation referenced by new",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-1", Name: "A", Value: "updated"}},
				},
			},
			scheduledAt: time.Now().Add(2 * time.Hour).Unix(),
			existingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:          "sfc-earlier",
					FeatureId:   "feature-id",
					ScheduledAt: time.Now().Add(time.Hour).Unix(),
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
						},
					},
					Status: featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				},
			},
			expectedLen: 1,
			expectedTypes: []featureproto.ScheduledChangeConflict_ConflictType{
				featureproto.ScheduledChangeConflict_CONFLICT_TYPE_DEPENDENCY_MISSING,
			},
			expectedFieldCheck: func(t *testing.T, conflicts []*featureproto.ScheduledChangeConflict) {
				t.Helper()
				assert.Contains(t, conflicts[0].Description, "sfc-earlier")
				assert.Equal(t, "variation var-1", conflicts[0].ConflictingField)
			},
		},
		{
			desc: "valid references: create variation not a conflict",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_CREATE, Variation: &featureproto.Variation{Id: "var-new", Name: "New", Value: "new"}},
				},
			},
			scheduledAt:       time.Now().Add(2 * time.Hour).Unix(),
			existingSchedules: []*featureproto.ScheduledFlagChange{},
			expectedLen:       0,
		},
		{
			desc: "dependency + invalid ref combined",
			flag: defaultFlag,
			payload: &featureproto.ScheduledChangePayload{
				VariationChanges: []*featureproto.VariationChange{
					{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-ghost"}},
				},
				OffVariation: wrapperspb.String("var-1"),
			},
			scheduledAt: time.Now().Add(2 * time.Hour).Unix(),
			existingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:          "sfc-earlier",
					FeatureId:   "feature-id",
					ScheduledAt: time.Now().Add(time.Hour).Unix(),
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-1"}},
						},
					},
					Status: featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				},
			},
			expectedLen: 2, // DEPENDENCY_MISSING (off_variation) + INVALID_REFERENCE (var-ghost)
			expectedTypes: []featureproto.ScheduledChangeConflict_ConflictType{
				featureproto.ScheduledChangeConflict_CONFLICT_TYPE_DEPENDENCY_MISSING,
				featureproto.ScheduledChangeConflict_CONFLICT_TYPE_INVALID_REFERENCE,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewMockScheduledFlagChangeStorage(ctrl)
			detector := NewConflictDetector(storage)

			storage.EXPECT().ListScheduledFlagChanges(gomock.Any(), gomock.Any()).
				Return(p.existingSchedules, len(p.existingSchedules), int64(len(p.existingSchedules)), nil)

			conflicts, err := detector.DetectConflictsOnCreate(
				context.Background(),
				p.flag,
				p.payload,
				p.scheduledAt,
				"ns0",
				p.excludeScheduleID,
			)
			require.NoError(t, err)
			assert.Len(t, conflicts, p.expectedLen)
			for i, expectedType := range p.expectedTypes {
				assert.Equal(t, expectedType, conflicts[i].Type)
			}
			if p.expectedFieldCheck != nil {
				p.expectedFieldCheck(t, conflicts)
			}
		})
	}
}

func TestConflictDetector_DetectConflictsOnFlagChange(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc              string
		flag              *featureproto.Feature
		pendingSchedules  []*featureproto.ScheduledFlagChange
		setupUpdateExpect func(t *testing.T, storage *mock.MockScheduledFlagChangeStorage)
		expectedCount     int
	}{
		{
			desc: "no pending schedules",
			flag: &featureproto.Feature{
				Id:      "feature-id",
				Version: 3,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Name: "A", Value: "true"},
				},
			},
			pendingSchedules:  []*featureproto.ScheduledFlagChange{},
			setupUpdateExpect: func(_ *testing.T, _ *mock.MockScheduledFlagChangeStorage) {},
			expectedCount:     0,
		},
		{
			desc: "schedule version matches: no conflict",
			flag: &featureproto.Feature{
				Id:      "feature-id",
				Version: 2,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Name: "A", Value: "true"},
				},
			},
			pendingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:                    "sfc-1",
					FlagVersionAtCreation: 2,
					Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
					Payload:               &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
				},
			},
			setupUpdateExpect: func(_ *testing.T, _ *mock.MockScheduledFlagChangeStorage) {},
			expectedCount:     0,
		},
		{
			desc: "version mismatch but no stale references: no conflict",
			flag: &featureproto.Feature{
				Id:      "feature-id",
				Version: 5,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Name: "A", Value: "true"},
				},
			},
			pendingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:                    "sfc-1",
					FeatureId:             "feature-id",
					EnvironmentId:         "ns0",
					FlagVersionAtCreation: 1,
					Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
					Payload:               &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
				},
			},
			setupUpdateExpect: func(_ *testing.T, _ *mock.MockScheduledFlagChangeStorage) {},
			expectedCount:     0, // Enable flag doesn't reference any variation/rule, so no conflict
		},
		{
			desc: "version mismatch with stale variation reference: marks conflict",
			flag: &featureproto.Feature{
				Id:      "feature-id",
				Version: 5,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Name: "A", Value: "true"},
				},
			},
			pendingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:                    "sfc-1",
					FeatureId:             "feature-id",
					EnvironmentId:         "ns0",
					FlagVersionAtCreation: 2,
					Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{ChangeType: featureproto.ChangeType_UPDATE, Variation: &featureproto.Variation{Id: "var-deleted"}},
						},
					},
				},
			},
			setupUpdateExpect: func(t *testing.T, storage *mock.MockScheduledFlagChangeStorage) {
				t.Helper()
				storage.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, sfc *domain.ScheduledFlagChange) error {
						assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT, sfc.Status)
						assert.NotEmpty(t, sfc.Conflicts)
						return nil
					})
			},
			expectedCount: 1,
		},
		{
			desc: "multiple schedules: only stale-reference ones get marked",
			flag: &featureproto.Feature{
				Id:      "feature-id",
				Version: 4,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Name: "A", Value: "true"},
				},
			},
			pendingSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:                    "sfc-stale-ref",
					FeatureId:             "feature-id",
					EnvironmentId:         "ns0",
					FlagVersionAtCreation: 1,
					Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{ChangeType: featureproto.ChangeType_DELETE, Variation: &featureproto.Variation{Id: "var-removed"}},
						},
					},
				},
				{
					Id:                    "sfc-stale-but-valid",
					FeatureId:             "feature-id",
					EnvironmentId:         "ns0",
					FlagVersionAtCreation: 1,
					Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
					Payload:               &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(true)},
				},
				{
					Id:                    "sfc-fresh",
					FeatureId:             "feature-id",
					EnvironmentId:         "ns0",
					FlagVersionAtCreation: 4,
					Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
					Payload:               &featureproto.ScheduledChangePayload{Enabled: wrapperspb.Bool(false)},
				},
			},
			setupUpdateExpect: func(_ *testing.T, storage *mock.MockScheduledFlagChangeStorage) {
				// Only sfc-stale-ref gets updated (has invalid var-removed reference)
				// sfc-stale-but-valid has version mismatch but no stale refs, so no update
				storage.EXPECT().UpdateScheduledFlagChange(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedCount: 1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewMockScheduledFlagChangeStorage(ctrl)
			detector := NewConflictDetector(storage)

			storage.EXPECT().ListScheduledFlagChanges(gomock.Any(), gomock.Any()).
				Return(p.pendingSchedules, len(p.pendingSchedules), int64(len(p.pendingSchedules)), nil)
			p.setupUpdateExpect(t, storage)

			count, err := detector.DetectConflictsOnFlagChange(context.Background(), p.flag, "ns0")
			require.NoError(t, err)
			assert.Equal(t, p.expectedCount, count)
		})
	}
}

func TestDetectConflictsOnFlagChange_AutoRecovery(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		flag           *featureproto.Feature
		schedule       *featureproto.ScheduledFlagChange
		setupFeatureFS func(*gomock.Controller) *mock.MockFeatureStorage
		expectUpdate   bool
		expectedStatus featureproto.ScheduledFlagChangeStatus
		expectedCount  int
	}{
		{
			desc: "CONFLICT recovers to PENDING when variation re-added",
			flag: &featureproto.Feature{
				Id:      "feature-1",
				Version: 5,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Value: "true"},
					{Id: "var-2", Value: "false"},
				},
			},
			schedule: &featureproto.ScheduledFlagChange{
				Id:                    "sfc-conflict",
				FeatureId:             "feature-1",
				EnvironmentId:         "env-1",
				Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
				FlagVersionAtCreation: 2,
				Payload: &featureproto.ScheduledChangePayload{
					VariationChanges: []*featureproto.VariationChange{
						{
							ChangeType: featureproto.ChangeType_UPDATE,
							Variation:  &featureproto.Variation{Id: "var-2", Value: "updated"},
						},
					},
				},
			},
			expectUpdate:   true,
			expectedStatus: featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			expectedCount:  1,
		},
		{
			desc: "CONFLICT stays when reference still invalid",
			flag: &featureproto.Feature{
				Id:      "feature-1",
				Version: 5,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Value: "true"},
				},
			},
			schedule: &featureproto.ScheduledFlagChange{
				Id:                    "sfc-conflict",
				FeatureId:             "feature-1",
				EnvironmentId:         "env-1",
				Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
				FlagVersionAtCreation: 2,
				Payload: &featureproto.ScheduledChangePayload{
					VariationChanges: []*featureproto.VariationChange{
						{
							ChangeType: featureproto.ChangeType_UPDATE,
							Variation:  &featureproto.Variation{Id: "var-deleted", Value: "nope"},
						},
					},
				},
			},
			expectUpdate:  false,
			expectedCount: 0,
		},
		{
			desc: "CONFLICT stays when same-flag refs valid but cross-flag prereq still broken",
			flag: &featureproto.Feature{
				Id:      "feature-1",
				Version: 5,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Value: "true"},
				},
			},
			schedule: &featureproto.ScheduledFlagChange{
				Id:                    "sfc-conflict",
				FeatureId:             "feature-1",
				EnvironmentId:         "env-1",
				Status:                featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
				FlagVersionAtCreation: 2,
				Payload: &featureproto.ScheduledChangePayload{
					PrerequisiteChanges: []*featureproto.PrerequisiteChange{
						{
							ChangeType: featureproto.ChangeType_CREATE,
							Prerequisite: &featureproto.Prerequisite{
								FeatureId:   "flag-b",
								VariationId: "var-beta-deleted",
							},
						},
					},
				},
			},
			setupFeatureFS: func(ctrl *gomock.Controller) *mock.MockFeatureStorage {
				fs := mock.NewMockFeatureStorage(ctrl)
				// flag-b exists but var-beta-deleted does not
				fs.EXPECT().GetFeature(
					gomock.Any(), "flag-b", "env-1",
				).Return(&domain.Feature{
					Feature: &featureproto.Feature{
						Id: "flag-b",
						Variations: []*featureproto.Variation{
							{Id: "var-beta-active"},
						},
					},
				}, nil)
				return fs
			},
			expectUpdate:  false,
			expectedCount: 0,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock.NewMockScheduledFlagChangeStorage(ctrl)

			var detector *ConflictDetector
			if p.setupFeatureFS != nil {
				fs := p.setupFeatureFS(ctrl)
				detector = NewConflictDetectorWithFeatureStorage(
					storage, fs, nil,
				)
			} else {
				detector = NewConflictDetector(storage)
			}

			storage.EXPECT().ListScheduledFlagChanges(
				gomock.Any(), gomock.Any(),
			).Return(
				[]*featureproto.ScheduledFlagChange{p.schedule},
				1, int64(1), nil,
			)

			if p.expectUpdate {
				storage.EXPECT().UpdateScheduledFlagChange(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(
					func(_ context.Context, sfc *domain.ScheduledFlagChange) error {
						assert.Equal(t, p.expectedStatus, sfc.Status)
						if p.expectedStatus == featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING {
							assert.Nil(t, sfc.Conflicts)
						}
						return nil
					},
				)
			}

			count, err := detector.DetectConflictsOnFlagChange(
				context.Background(), p.flag, "env-1",
			)
			require.NoError(t, err)
			assert.Equal(t, p.expectedCount, count)
		})
	}
}

func TestDetectCrossFlagConflicts(t *testing.T) {
	t.Parallel()

	flagAFeature := &domain.Feature{
		Feature: &featureproto.Feature{
			Id:         "flag-a",
			Variations: []*featureproto.Variation{{Id: "var-a1"}},
		},
	}

	patterns := []struct {
		desc           string
		schedule       *featureproto.ScheduledFlagChange
		flagBVars      []*featureproto.Variation
		expectUpdate   bool
		expectedStatus featureproto.ScheduledFlagChangeStatus
		expectedCount  int
	}{
		{
			desc: "PENDING marked CONFLICT when prereq variation deleted",
			schedule: &featureproto.ScheduledFlagChange{
				Id:            "sfc-1",
				FeatureId:     "flag-a",
				EnvironmentId: "env-1",
				Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				Payload: &featureproto.ScheduledChangePayload{
					PrerequisiteChanges: []*featureproto.PrerequisiteChange{
						{
							ChangeType: featureproto.ChangeType_CREATE,
							Prerequisite: &featureproto.Prerequisite{
								FeatureId:   "flag-b",
								VariationId: "var-beta-deleted",
							},
						},
					},
				},
			},
			flagBVars:      []*featureproto.Variation{{Id: "var-beta-active"}},
			expectUpdate:   true,
			expectedStatus: featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
			expectedCount:  1,
		},
		{
			desc: "PENDING stays when prereq variation still exists",
			schedule: &featureproto.ScheduledFlagChange{
				Id:            "sfc-1",
				FeatureId:     "flag-a",
				EnvironmentId: "env-1",
				Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
				Payload: &featureproto.ScheduledChangePayload{
					PrerequisiteChanges: []*featureproto.PrerequisiteChange{
						{
							ChangeType: featureproto.ChangeType_CREATE,
							Prerequisite: &featureproto.Prerequisite{
								FeatureId:   "flag-b",
								VariationId: "var-beta",
							},
						},
					},
				},
			},
			flagBVars:     []*featureproto.Variation{{Id: "var-beta"}},
			expectUpdate:  false,
			expectedCount: 0,
		},
		{
			desc: "CONFLICT auto-recovers when prereq variation re-added",
			schedule: &featureproto.ScheduledFlagChange{
				Id:            "sfc-1",
				FeatureId:     "flag-a",
				EnvironmentId: "env-1",
				Status:        featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
				Payload: &featureproto.ScheduledChangePayload{
					PrerequisiteChanges: []*featureproto.PrerequisiteChange{
						{
							ChangeType: featureproto.ChangeType_CREATE,
							Prerequisite: &featureproto.Prerequisite{
								FeatureId:   "flag-b",
								VariationId: "var-beta",
							},
						},
					},
				},
			},
			flagBVars:      []*featureproto.Variation{{Id: "var-beta"}},
			expectUpdate:   true,
			expectedStatus: featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING,
			expectedCount:  1,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sfcStorage := mock.NewMockScheduledFlagChangeStorage(ctrl)
			featureStorage := mock.NewMockFeatureStorage(ctrl)
			detector := NewConflictDetectorWithFeatureStorage(
				sfcStorage, featureStorage, nil,
			)

			sfcStorage.EXPECT().ListScheduledFlagChanges(
				gomock.Any(), gomock.Any(),
			).Return(
				[]*featureproto.ScheduledFlagChange{p.schedule},
				1, int64(1), nil,
			)

			featureStorage.EXPECT().GetFeature(
				gomock.Any(), "flag-a", "env-1",
			).Return(flagAFeature, nil)

			featureStorage.EXPECT().GetFeature(
				gomock.Any(), "flag-b", "env-1",
			).Return(&domain.Feature{
				Feature: &featureproto.Feature{
					Id:         "flag-b",
					Variations: p.flagBVars,
				},
			}, nil)

			if p.expectUpdate {
				sfcStorage.EXPECT().UpdateScheduledFlagChange(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(
					func(_ context.Context, sfc *domain.ScheduledFlagChange) error {
						assert.Equal(t, p.expectedStatus, sfc.Status)
						if p.expectedStatus == featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING {
							assert.Nil(t, sfc.Conflicts)
						} else {
							assert.NotEmpty(t, sfc.Conflicts)
						}
						return nil
					},
				)
			}

			count, err := detector.DetectCrossFlagConflicts(
				context.Background(), "flag-b", "env-1",
			)
			require.NoError(t, err)
			assert.Equal(t, p.expectedCount, count)
		})
	}
}

func TestDetectCrossFlagConflicts_NilFeatureStorage(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sfcStorage := mock.NewMockScheduledFlagChangeStorage(ctrl)
	detector := NewConflictDetector(sfcStorage) // No feature storage

	count, err := detector.DetectCrossFlagConflicts(
		context.Background(), "flag-b", "env-1",
	)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestAttemptRecoveryForFlag(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc              string
		featureID         string
		environmentID     string
		flag              *featureproto.Feature
		conflictSchedules []*featureproto.ScheduledFlagChange
		expectedRecovered int
		expectError       bool
	}{
		{
			desc:          "successfully recovers CONFLICT schedule when references valid",
			featureID:     "feature-1",
			environmentID: "env-1",
			flag: &featureproto.Feature{
				Id:      "feature-1",
				Version: 3,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Value: "true"},
					{Id: "var-2", Value: "false"},
				},
			},
			conflictSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:        "sfc-1",
					FeatureId: "feature-1",
					Status:    featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{
								ChangeType: featureproto.ChangeType_UPDATE,
								Variation:  &featureproto.Variation{Id: "var-2", Value: "updated"},
							},
						},
					},
				},
			},
			expectedRecovered: 1,
			expectError:       false,
		},
		{
			desc:          "no recovery when references still invalid",
			featureID:     "feature-1",
			environmentID: "env-1",
			flag: &featureproto.Feature{
				Id:      "feature-1",
				Version: 3,
				Variations: []*featureproto.Variation{
					{Id: "var-1", Value: "true"},
				},
			},
			conflictSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:        "sfc-1",
					FeatureId: "feature-1",
					Status:    featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{
								ChangeType: featureproto.ChangeType_UPDATE,
								Variation:  &featureproto.Variation{Id: "var-deleted", Value: "nope"},
							},
						},
					},
				},
			},
			expectedRecovered: 0,
			expectError:       false,
		},
		{
			desc:              "no recovery when no conflict schedules",
			featureID:         "feature-1",
			environmentID:     "env-1",
			conflictSchedules: []*featureproto.ScheduledFlagChange{},
			expectedRecovered: 0,
			expectError:       false,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sfcStorage := mock.NewMockScheduledFlagChangeStorage(ctrl)
			featureStorage := mock.NewMockFeatureStorage(ctrl)
			detector := NewConflictDetectorWithFeatureStorage(sfcStorage, featureStorage, nil)

			// Mock listConflictSchedulesByFlag
			sfcStorage.EXPECT().ListScheduledFlagChanges(
				gomock.Any(), gomock.Any(),
			).Return(p.conflictSchedules, len(p.conflictSchedules), int64(len(p.conflictSchedules)), nil)

			if len(p.conflictSchedules) > 0 {
				// Mock GetFeature
				featureStorage.EXPECT().GetFeature(
					gomock.Any(), p.featureID, p.environmentID,
				).Return(&domain.Feature{Feature: p.flag}, nil)

				// Mock UpdateScheduledFlagChange if recovery expected
				if p.expectedRecovered > 0 {
					sfcStorage.EXPECT().UpdateScheduledFlagChange(
						gomock.Any(), gomock.Any(),
					).Times(p.expectedRecovered).DoAndReturn(
						func(_ context.Context, sfc *domain.ScheduledFlagChange) error {
							assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, sfc.Status)
							assert.Nil(t, sfc.Conflicts)
							return nil
						},
					)
				}
			}

			recovered, err := detector.attemptRecoveryForFlag(
				context.Background(), p.featureID, p.environmentID,
			)

			if p.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, p.expectedRecovered, recovered)
			}
		})
	}
}

func TestAttemptRecoveryForEnvironment(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc              string
		environmentID     string
		excludeFeatureID  string
		conflictSchedules []*featureproto.ScheduledFlagChange
		flagsMap          map[string]*featureproto.Feature
		expectedRecovered int
	}{
		{
			desc:             "recovers schedules from multiple flags",
			environmentID:    "env-1",
			excludeFeatureID: "",
			conflictSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:        "sfc-1",
					FeatureId: "feature-1",
					Status:    featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{
								ChangeType: featureproto.ChangeType_UPDATE,
								Variation:  &featureproto.Variation{Id: "var-1", Value: "updated"},
							},
						},
					},
				},
				{
					Id:        "sfc-2",
					FeatureId: "feature-2",
					Status:    featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{
								ChangeType: featureproto.ChangeType_UPDATE,
								Variation:  &featureproto.Variation{Id: "var-a", Value: "updated"},
							},
						},
					},
				},
			},
			flagsMap: map[string]*featureproto.Feature{
				"feature-1": {
					Id:         "feature-1",
					Version:    3,
					Variations: []*featureproto.Variation{{Id: "var-1"}},
				},
				"feature-2": {
					Id:         "feature-2",
					Version:    2,
					Variations: []*featureproto.Variation{{Id: "var-a"}},
				},
			},
			expectedRecovered: 2,
		},
		{
			desc:             "excludes specific flag",
			environmentID:    "env-1",
			excludeFeatureID: "feature-1",
			conflictSchedules: []*featureproto.ScheduledFlagChange{
				{
					Id:        "sfc-2",
					FeatureId: "feature-2",
					Status:    featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
					Payload: &featureproto.ScheduledChangePayload{
						VariationChanges: []*featureproto.VariationChange{
							{
								ChangeType: featureproto.ChangeType_UPDATE,
								Variation:  &featureproto.Variation{Id: "var-a", Value: "updated"},
							},
						},
					},
				},
			},
			flagsMap: map[string]*featureproto.Feature{
				"feature-2": {
					Id:         "feature-2",
					Version:    2,
					Variations: []*featureproto.Variation{{Id: "var-a"}},
				},
			},
			expectedRecovered: 1,
		},
		{
			desc:              "no schedules to recover",
			environmentID:     "env-1",
			excludeFeatureID:  "",
			conflictSchedules: []*featureproto.ScheduledFlagChange{},
			expectedRecovered: 0,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sfcStorage := mock.NewMockScheduledFlagChangeStorage(ctrl)
			featureStorage := mock.NewMockFeatureStorage(ctrl)
			detector := NewConflictDetectorWithFeatureStorage(sfcStorage, featureStorage, nil)

			// Mock listConflictSchedulesInEnvironment
			sfcStorage.EXPECT().ListScheduledFlagChanges(
				gomock.Any(), gomock.Any(),
			).Return(p.conflictSchedules, len(p.conflictSchedules), int64(len(p.conflictSchedules)), nil)

			// Mock GetFeature for each unique flag
			for featureID, flag := range p.flagsMap {
				featureStorage.EXPECT().GetFeature(
					gomock.Any(), featureID, p.environmentID,
				).Return(&domain.Feature{Feature: flag}, nil)
			}

			// Mock UpdateScheduledFlagChange for recoveries
			if p.expectedRecovered > 0 {
				sfcStorage.EXPECT().UpdateScheduledFlagChange(
					gomock.Any(), gomock.Any(),
				).Times(p.expectedRecovered).DoAndReturn(
					func(_ context.Context, sfc *domain.ScheduledFlagChange) error {
						assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, sfc.Status)
						return nil
					},
				)
			}

			recovered, err := detector.attemptRecoveryForEnvironment(
				context.Background(), p.environmentID, p.excludeFeatureID,
			)

			require.NoError(t, err)
			assert.Equal(t, p.expectedRecovered, recovered)
		})
	}
}

func TestAttemptRecoveryForFlag_WithPrerequisites(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sfcStorage := mock.NewMockScheduledFlagChangeStorage(ctrl)
	featureStorage := mock.NewMockFeatureStorage(ctrl)
	detector := NewConflictDetectorWithFeatureStorage(sfcStorage, featureStorage, nil)

	// Schedule with prerequisite that references another flag
	schedule := &featureproto.ScheduledFlagChange{
		Id:        "sfc-1",
		FeatureId: "feature-a",
		Status:    featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT,
		Payload: &featureproto.ScheduledChangePayload{
			PrerequisiteChanges: []*featureproto.PrerequisiteChange{
				{
					ChangeType: featureproto.ChangeType_CREATE,
					Prerequisite: &featureproto.Prerequisite{
						FeatureId:   "feature-b",
						VariationId: "var-beta",
					},
				},
			},
		},
	}

	flagA := &featureproto.Feature{
		Id:         "feature-a",
		Version:    2,
		Variations: []*featureproto.Variation{{Id: "var-a"}},
	}

	flagB := &featureproto.Feature{
		Id:         "feature-b",
		Version:    1,
		Variations: []*featureproto.Variation{{Id: "var-beta"}}, // Prerequisite variation now exists
	}

	// Mock list schedules
	sfcStorage.EXPECT().ListScheduledFlagChanges(
		gomock.Any(), gomock.Any(),
	).Return([]*featureproto.ScheduledFlagChange{schedule}, 1, int64(1), nil)

	// Mock get feature-a
	featureStorage.EXPECT().GetFeature(
		gomock.Any(), "feature-a", "env-1",
	).Return(&domain.Feature{Feature: flagA}, nil)

	// Mock get feature-b (for prerequisite validation)
	featureStorage.EXPECT().GetFeature(
		gomock.Any(), "feature-b", "env-1",
	).Return(&domain.Feature{Feature: flagB}, nil)

	// Mock update (should restore to PENDING since prerequisite is now valid)
	sfcStorage.EXPECT().UpdateScheduledFlagChange(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(
		func(_ context.Context, sfc *domain.ScheduledFlagChange) error {
			assert.Equal(t, featureproto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, sfc.Status)
			return nil
		},
	)

	recovered, err := detector.attemptRecoveryForFlag(
		context.Background(), "feature-a", "env-1",
	)

	require.NoError(t, err)
	assert.Equal(t, 1, recovered)
}

func TestAttemptRecoveryForFlag_NilFeatureStorage(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sfcStorage := mock.NewMockScheduledFlagChangeStorage(ctrl)
	detector := NewConflictDetector(sfcStorage) // No feature storage

	recovered, err := detector.attemptRecoveryForFlag(
		context.Background(), "feature-1", "env-1",
	)

	require.Error(t, err)
	assert.Equal(t, 0, recovered)
	assert.Contains(t, err.Error(), "feature storage not configured")
}

func TestScheduleReferencesFlag_FeatureFlagClauses(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		payload      *featureproto.ScheduledChangePayload
		targetFlagID string
		expected     bool
	}{
		{
			desc: "references flag via FEATURE_FLAG clause",
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-1",
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "flag-b", // References flag-b
									Values:    []string{"var-true"},
								},
							},
						},
					},
				},
			},
			targetFlagID: "flag-b",
			expected:     true,
		},
		{
			desc: "does not reference flag (different flag in clause)",
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-1",
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "flag-c", // References flag-c, not flag-b
									Values:    []string{"var-true"},
								},
							},
						},
					},
				},
			},
			targetFlagID: "flag-b",
			expected:     false,
		},
		{
			desc: "references flag via prerequisite AND FEATURE_FLAG clause",
			payload: &featureproto.ScheduledChangePayload{
				PrerequisiteChanges: []*featureproto.PrerequisiteChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Prerequisite: &featureproto.Prerequisite{
							FeatureId:   "flag-b",
							VariationId: "var-1",
						},
					},
				},
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-1",
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "flag-b",
									Values:    []string{"var-true"},
								},
							},
						},
					},
				},
			},
			targetFlagID: "flag-b",
			expected:     true,
		},
		{
			desc: "no FEATURE_FLAG clauses",
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-1",
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_SEGMENT,
									Attribute: "segment-1",
									Values:    []string{"value"},
								},
							},
						},
					},
				},
			},
			targetFlagID: "flag-b",
			expected:     false,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			result := scheduleReferencesFlag(p.payload, p.targetFlagID)
			assert.Equal(t, p.expected, result)
		})
	}
}

func TestValidateFeatureFlagReferences(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc              string
		payload           *featureproto.ScheduledChangePayload
		environmentID     string
		referencedFlags   map[string]*featureproto.Feature
		expectedConflicts int
	}{
		{
			desc: "valid FEATURE_FLAG clause - flag and variation exist",
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-1",
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "flag-b",
									Values:    []string{"var-true"},
								},
							},
						},
					},
				},
			},
			environmentID: "env-1",
			referencedFlags: map[string]*featureproto.Feature{
				"flag-b": {
					Id:         "flag-b",
					Variations: []*featureproto.Variation{{Id: "var-true"}},
				},
			},
			expectedConflicts: 0,
		},
		{
			desc: "invalid FEATURE_FLAG clause - referenced flag not found",
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-1",
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "flag-deleted",
									Values:    []string{"var-true"},
								},
							},
						},
					},
				},
			},
			environmentID:     "env-1",
			referencedFlags:   map[string]*featureproto.Feature{},
			expectedConflicts: 1,
		},
		{
			desc: "invalid FEATURE_FLAG clause - variation not found",
			payload: &featureproto.ScheduledChangePayload{
				RuleChanges: []*featureproto.RuleChange{
					{
						ChangeType: featureproto.ChangeType_CREATE,
						Rule: &featureproto.Rule{
							Id: "rule-1",
							Clauses: []*featureproto.Clause{
								{
									Operator:  featureproto.Clause_FEATURE_FLAG,
									Attribute: "flag-b",
									Values:    []string{"var-deleted"},
								},
							},
						},
					},
				},
			},
			environmentID: "env-1",
			referencedFlags: map[string]*featureproto.Feature{
				"flag-b": {
					Id:         "flag-b",
					Variations: []*featureproto.Variation{{Id: "var-true"}},
				},
			},
			expectedConflicts: 1,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sfcStorage := mock.NewMockScheduledFlagChangeStorage(ctrl)
			featureStorage := mock.NewMockFeatureStorage(ctrl)
			detector := NewConflictDetectorWithFeatureStorage(sfcStorage, featureStorage, nil)

			// Mock GetFeature calls
			for flagID, flag := range p.referencedFlags {
				featureStorage.EXPECT().GetFeature(
					gomock.Any(), flagID, p.environmentID,
				).Return(&domain.Feature{Feature: flag}, nil)
			}

			// Mock GetFeature failures for missing flags
			for _, rc := range p.payload.RuleChanges {
				if rc != nil && rc.Rule != nil {
					for _, clause := range rc.Rule.Clauses {
						if clause != nil && clause.Operator == featureproto.Clause_FEATURE_FLAG {
							if _, exists := p.referencedFlags[clause.Attribute]; !exists {
								featureStorage.EXPECT().GetFeature(
									gomock.Any(), clause.Attribute, p.environmentID,
								).Return(nil, fmt.Errorf("not found"))
							}
						}
					}
				}
			}

			conflicts := detector.validateFeatureFlagReferences(
				context.Background(), p.payload, p.environmentID, time.Now().Unix(),
			)

			assert.Equal(t, p.expectedConflicts, len(conflicts))
		})
	}
}
