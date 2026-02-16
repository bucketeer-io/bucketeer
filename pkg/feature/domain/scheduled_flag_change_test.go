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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewScheduledFlagChange(t *testing.T) {
	t.Parallel()

	featureID := "feature-1"
	environmentID := "env-1"
	scheduledAt := time.Now().Add(time.Hour).Unix()
	timezone := "Asia/Tokyo"
	payload := &proto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(true),
	}
	comment := "Enable flag for testing"
	flagVersion := int32(5)
	createdBy := "user@example.com"

	sfc, err := NewScheduledFlagChange(
		featureID,
		environmentID,
		scheduledAt,
		timezone,
		payload,
		comment,
		flagVersion,
		createdBy,
	)

	require.NoError(t, err)
	assert.NotEmpty(t, sfc.Id)
	assert.Equal(t, featureID, sfc.FeatureId)
	assert.Equal(t, environmentID, sfc.EnvironmentId)
	assert.Equal(t, scheduledAt, sfc.ScheduledAt)
	assert.Equal(t, timezone, sfc.Timezone)
	assert.Equal(t, payload, sfc.Payload)
	assert.Equal(t, comment, sfc.Comment)
	assert.Equal(t, proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_PENDING, sfc.Status)
	assert.Equal(t, flagVersion, sfc.FlagVersionAtCreation)
	assert.Equal(t, createdBy, sfc.CreatedBy)
	assert.True(t, sfc.CreatedAt > 0)
	assert.Equal(t, sfc.CreatedAt, sfc.UpdatedAt)
}

func TestNewScheduledFlagChangeDefaultTimezone(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"", // empty timezone
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)

	require.NoError(t, err)
	assert.Equal(t, "UTC", sfc.Timezone)
}

func TestScheduledFlagChangeUpdateSchedule(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	oldUpdatedAt := sfc.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference

	newScheduledAt := time.Now().Add(2 * time.Hour).Unix()
	newTimezone := "America/New_York"
	updatedBy := "admin@example.com"

	sfc.UpdateSchedule(newScheduledAt, newTimezone, updatedBy)

	assert.Equal(t, newScheduledAt, sfc.ScheduledAt)
	assert.Equal(t, newTimezone, sfc.Timezone)
	assert.Equal(t, updatedBy, sfc.UpdatedBy)
	assert.True(t, sfc.UpdatedAt >= oldUpdatedAt)
}

func TestScheduledFlagChangeUpdatePayload(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"original comment",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	newPayload := &proto.ScheduledChangePayload{
		Enabled: wrapperspb.Bool(false),
	}
	newComment := "updated comment"
	updatedBy := "admin@example.com"

	sfc.UpdatePayload(newPayload, newComment, updatedBy)

	assert.Equal(t, newPayload, sfc.Payload)
	assert.Equal(t, newComment, sfc.Comment)
	assert.Equal(t, updatedBy, sfc.UpdatedBy)
}

func TestScheduledFlagChangeCancel(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	updatedBy := "admin@example.com"
	reason := "Flag was archived"

	sfc.Cancel(updatedBy, reason)

	assert.Equal(t, proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CANCELLED, sfc.Status)
	assert.Equal(t, reason, sfc.FailureReason)
	assert.Equal(t, updatedBy, sfc.UpdatedBy)
}

func TestScheduledFlagChangeMarkExecuted(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	sfc.MarkExecuted()

	assert.Equal(t, proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_EXECUTED, sfc.Status)
	assert.True(t, sfc.ExecutedAt > 0)
}

func TestScheduledFlagChangeMarkFailed(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	reason := "Variation not found"

	sfc.MarkFailed(reason)

	assert.Equal(t, proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_FAILED, sfc.Status)
	assert.Equal(t, reason, sfc.FailureReason)
}

func TestScheduledFlagChangeMarkConflict(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	conflicts := []*proto.ScheduledChangeConflict{
		{
			Type:        proto.ScheduledChangeConflict_CONFLICT_TYPE_VERSION_MISMATCH,
			Description: "Flag was modified",
		},
	}

	sfc.MarkConflict(conflicts)

	assert.Equal(t, proto.ScheduledFlagChangeStatus_SCHEDULED_FLAG_CHANGE_STATUS_CONFLICT, sfc.Status)
	assert.Equal(t, conflicts, sfc.Conflicts)
}

func TestScheduledFlagChangeIsPending(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	assert.True(t, sfc.IsPending())

	sfc.MarkExecuted()
	assert.False(t, sfc.IsPending())
}

func TestScheduledFlagChangeIsConflict(t *testing.T) {
	t.Parallel()

	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)

	assert.False(t, sfc.IsConflict())

	sfc.MarkConflict([]*proto.ScheduledChangeConflict{})
	assert.True(t, sfc.IsConflict())
}

func TestScheduledFlagChangeIsDue(t *testing.T) {
	t.Parallel()

	// Not due (scheduled in the future)
	sfc, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)
	assert.False(t, sfc.IsDue())

	// Due (scheduled in the past)
	sfc2, err := NewScheduledFlagChange(
		"feature-1",
		"env-1",
		time.Now().Add(-time.Hour).Unix(),
		"UTC",
		&proto.ScheduledChangePayload{},
		"",
		1,
		"user@example.com",
	)
	require.NoError(t, err)
	assert.True(t, sfc2.IsDue())
}

func TestScheduledFlagChangeDetermineCategory(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		payload  *proto.ScheduledChangePayload
		expected proto.ScheduledChangeCategory
	}{
		{
			desc:     "nil payload",
			payload:  nil,
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_UNSPECIFIED,
		},
		{
			desc:     "empty payload",
			payload:  &proto.ScheduledChangePayload{},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_UNSPECIFIED,
		},
		{
			desc: "targeting: rule changes",
			payload: &proto.ScheduledChangePayload{
				RuleChanges: []*proto.RuleChange{
					{ChangeType: proto.ChangeType_CREATE},
				},
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_TARGETING,
		},
		{
			desc: "targeting: target changes",
			payload: &proto.ScheduledChangePayload{
				TargetChanges: []*proto.TargetChange{
					{ChangeType: proto.ChangeType_CREATE},
				},
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_TARGETING,
		},
		{
			desc: "targeting: default strategy",
			payload: &proto.ScheduledChangePayload{
				DefaultStrategy: &proto.Strategy{},
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_TARGETING,
		},
		{
			desc: "variations: variation changes",
			payload: &proto.ScheduledChangePayload{
				VariationChanges: []*proto.VariationChange{
					{ChangeType: proto.ChangeType_UPDATE},
				},
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_VARIATIONS,
		},
		{
			desc: "variations: off variation",
			payload: &proto.ScheduledChangePayload{
				OffVariation: wrapperspb.String("var-1"),
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_VARIATIONS,
		},
		{
			desc: "settings: enabled",
			payload: &proto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(true),
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_SETTINGS,
		},
		{
			desc: "settings: name",
			payload: &proto.ScheduledChangePayload{
				Name: wrapperspb.String("new name"),
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_SETTINGS,
		},
		{
			desc: "settings: tag changes",
			payload: &proto.ScheduledChangePayload{
				TagChanges: []*proto.TagChange{
					{ChangeType: proto.ChangeType_CREATE, Tag: "tag1"},
				},
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_SETTINGS,
		},
		{
			desc: "mixed: targeting + settings",
			payload: &proto.ScheduledChangePayload{
				RuleChanges: []*proto.RuleChange{
					{ChangeType: proto.ChangeType_CREATE},
				},
				Enabled: wrapperspb.Bool(true),
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_MIXED,
		},
		{
			desc: "mixed: all categories",
			payload: &proto.ScheduledChangePayload{
				RuleChanges: []*proto.RuleChange{
					{ChangeType: proto.ChangeType_CREATE},
				},
				VariationChanges: []*proto.VariationChange{
					{ChangeType: proto.ChangeType_UPDATE},
				},
				Enabled: wrapperspb.Bool(true),
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_MIXED,
		},
		// ResetSamplingSeed tests - it should NOT cause MIXED when combined with targeting/variations
		{
			desc: "settings: reset sampling seed alone",
			payload: &proto.ScheduledChangePayload{
				ResetSamplingSeed: true,
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_SETTINGS,
		},
		{
			desc: "targeting: rule changes + reset sampling seed (should NOT be MIXED)",
			payload: &proto.ScheduledChangePayload{
				RuleChanges: []*proto.RuleChange{
					{ChangeType: proto.ChangeType_CREATE},
				},
				ResetSamplingSeed: true,
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_TARGETING,
		},
		{
			desc: "targeting: default strategy + reset sampling seed (should NOT be MIXED)",
			payload: &proto.ScheduledChangePayload{
				DefaultStrategy:   &proto.Strategy{},
				ResetSamplingSeed: true,
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_TARGETING,
		},
		{
			desc: "variations: variation changes + reset sampling seed (should NOT be MIXED)",
			payload: &proto.ScheduledChangePayload{
				VariationChanges: []*proto.VariationChange{
					{ChangeType: proto.ChangeType_UPDATE},
				},
				ResetSamplingSeed: true,
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_VARIATIONS,
		},
		{
			desc: "variations: off variation + reset sampling seed (should NOT be MIXED)",
			payload: &proto.ScheduledChangePayload{
				OffVariation:      wrapperspb.String("var-1"),
				ResetSamplingSeed: true,
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_VARIATIONS,
		},
		{
			desc: "mixed: targeting + variations + reset sampling seed",
			payload: &proto.ScheduledChangePayload{
				RuleChanges: []*proto.RuleChange{
					{ChangeType: proto.ChangeType_CREATE},
				},
				VariationChanges: []*proto.VariationChange{
					{ChangeType: proto.ChangeType_UPDATE},
				},
				ResetSamplingSeed: true,
			},
			expected: proto.ScheduledChangeCategory_SCHEDULED_CHANGE_CATEGORY_MIXED,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sfc := &ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Payload: p.payload,
				},
			}
			assert.Equal(t, p.expected, sfc.DetermineCategory())
		})
	}
}

func TestGenerateChangeSummaries_EnableFlag(t *testing.T) {
	t.Parallel()

	sfc := &ScheduledFlagChange{
		ScheduledFlagChange: &proto.ScheduledFlagChange{
			Payload: &proto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(true),
			},
		},
	}

	summaries := sfc.GenerateChangeSummaries(nil)

	assert.Len(t, summaries, 1)
	assert.Equal(t, MsgKeyEnableFlag, summaries[0].MessageKey)
	assert.Nil(t, summaries[0].Values)
}

func TestGenerateChangeSummaries_DisableFlag(t *testing.T) {
	t.Parallel()

	sfc := &ScheduledFlagChange{
		ScheduledFlagChange: &proto.ScheduledFlagChange{
			Payload: &proto.ScheduledChangePayload{
				Enabled: wrapperspb.Bool(false),
			},
		},
	}

	summaries := sfc.GenerateChangeSummaries(nil)

	assert.Len(t, summaries, 1)
	assert.Equal(t, MsgKeyDisableFlag, summaries[0].MessageKey)
}

func TestGenerateChangeSummaries_UpdateVariation(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc              string
		flag              *proto.Feature
		variation         *proto.Variation
		expectedCount     int
		expectedFirstKey  string
		expectedSecondKey string
		assertions        func(t *testing.T, summaries []*proto.ChangeSummary)
	}{
		{
			desc: "value only changed",
			flag: &proto.Feature{
				Variations: []*proto.Variation{
					{Id: "v1", Name: "Variation A", Value: "old-value"},
				},
			},
			variation: &proto.Variation{
				Id:    "v1",
				Name:  "Variation A", // Same name
				Value: "new-value",   // Different value
			},
			expectedCount:    1,
			expectedFirstKey: MsgKeyChangeVariationValue,
			assertions: func(t *testing.T, summaries []*proto.ChangeSummary) {
				assert.Equal(t, "Variation A", summaries[0].Values["name"])
				assert.Equal(t, "old-value", summaries[0].Values["oldValue"])
				assert.Equal(t, "new-value", summaries[0].Values["newValue"])
			},
		},
		{
			desc: "name only changed",
			flag: &proto.Feature{
				Variations: []*proto.Variation{
					{Id: "v1", Name: "Old Name", Value: "same-value"},
				},
			},
			variation: &proto.Variation{
				Id:    "v1",
				Name:  "New Name",   // Different name
				Value: "same-value", // Same value
			},
			expectedCount:    1,
			expectedFirstKey: MsgKeyRenameVariation,
			assertions: func(t *testing.T, summaries []*proto.ChangeSummary) {
				assert.Equal(t, "Old Name", summaries[0].Values["oldName"])
				assert.Equal(t, "New Name", summaries[0].Values["newName"])
			},
		},
		{
			desc: "both name and value changed",
			flag: &proto.Feature{
				Variations: []*proto.Variation{
					{Id: "v1", Name: "Old Name", Value: "old-value"},
				},
			},
			variation: &proto.Variation{
				Id:    "v1",
				Name:  "New Name",  // Different name
				Value: "new-value", // Different value
			},
			expectedCount:     2,
			expectedFirstKey:  MsgKeyChangeVariationValue,
			expectedSecondKey: MsgKeyRenameVariation,
			assertions: func(t *testing.T, summaries []*proto.ChangeSummary) {
				// First summary: value change
				assert.Equal(t, "New Name", summaries[0].Values["name"])
				assert.Equal(t, "old-value", summaries[0].Values["oldValue"])
				assert.Equal(t, "new-value", summaries[0].Values["newValue"])
				// Second summary: name change
				assert.Equal(t, "Old Name", summaries[1].Values["oldName"])
				assert.Equal(t, "New Name", summaries[1].Values["newName"])
			},
		},
		{
			desc: "no change",
			flag: &proto.Feature{
				Variations: []*proto.Variation{
					{Id: "v1", Name: "Same Name", Value: "same-value"},
				},
			},
			variation: &proto.Variation{
				Id:    "v1",
				Name:  "Same Name",  // Same
				Value: "same-value", // Same
			},
			expectedCount:    1,
			expectedFirstKey: MsgKeyUpdateVariation,
			assertions: func(t *testing.T, summaries []*proto.ChangeSummary) {
				assert.Equal(t, "Same Name", summaries[0].Values["name"])
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sfc := &ScheduledFlagChange{
				ScheduledFlagChange: &proto.ScheduledFlagChange{
					Payload: &proto.ScheduledChangePayload{
						VariationChanges: []*proto.VariationChange{
							{
								ChangeType: proto.ChangeType_UPDATE,
								Variation:  p.variation,
							},
						},
					},
				},
			}

			summaries := sfc.GenerateChangeSummaries(p.flag)

			assert.Len(t, summaries, p.expectedCount)
			assert.Equal(t, p.expectedFirstKey, summaries[0].MessageKey)
			if p.expectedSecondKey != "" {
				assert.Equal(t, p.expectedSecondKey, summaries[1].MessageKey)
			}
			if p.assertions != nil {
				p.assertions(t, summaries)
			}
		})
	}
}

func TestGenerateChangeSummaries_MultipleChanges(t *testing.T) {
	t.Parallel()

	sfc := &ScheduledFlagChange{
		ScheduledFlagChange: &proto.ScheduledFlagChange{
			Payload: &proto.ScheduledChangePayload{
				Enabled:           wrapperspb.Bool(true),
				Name:              wrapperspb.String("Updated Name"),
				ResetSamplingSeed: true,
				TagChanges: []*proto.TagChange{
					{ChangeType: proto.ChangeType_CREATE, Tag: "tag1"},
				},
			},
		},
	}

	summaries := sfc.GenerateChangeSummaries(nil)

	assert.Len(t, summaries, 4)
	assert.Equal(t, MsgKeyEnableFlag, summaries[0].MessageKey)
	assert.Equal(t, MsgKeyRenameFlag, summaries[1].MessageKey)
	assert.Equal(t, MsgKeyResetSamplingSeed, summaries[2].MessageKey)
	assert.Equal(t, MsgKeyAddTag, summaries[3].MessageKey)
}

func TestGenerateChangeSummaries_NilPayload(t *testing.T) {
	t.Parallel()

	sfc := &ScheduledFlagChange{
		ScheduledFlagChange: &proto.ScheduledFlagChange{
			Payload: nil,
		},
	}

	summaries := sfc.GenerateChangeSummaries(nil)

	assert.Nil(t, summaries)
}

func TestGenerateChangeSummaries_EmptyPayload(t *testing.T) {
	t.Parallel()

	sfc := &ScheduledFlagChange{
		ScheduledFlagChange: &proto.ScheduledFlagChange{
			Payload: &proto.ScheduledChangePayload{},
		},
	}

	summaries := sfc.GenerateChangeSummaries(nil)

	assert.Empty(t, summaries)
}
