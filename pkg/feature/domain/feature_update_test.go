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

package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	"github.com/bucketeer-io/bucketeer/proto/common"
	"github.com/bucketeer-io/bucketeer/proto/feature"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestUpdateNoTimestampChangeWithSameValues(t *testing.T) {
	t.Parallel()

	// Generate valid UUIDs for variations
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)

	// Create a feature with specific timestamp
	fixedTimestamp := time.Now().Unix() - 3600 // 1 hour ago
	original := &Feature{
		Feature: &feature.Feature{
			Id:          "test-feature",
			Name:        "Test Feature",
			Description: "Test Description",
			Enabled:     true,
			Archived:    false,
			Version:     5,
			UpdatedAt:   fixedTimestamp,
			Variations: []*feature.Variation{
				{Id: v1.String(), Name: "v1", Value: "true"},
				{Id: v2.String(), Name: "v2", Value: "false"},
			},
			Targets: []*feature.Target{
				{Variation: v1.String(), Users: []string{}},
				{Variation: v2.String(), Users: []string{}},
			},
			DefaultStrategy: &feature.Strategy{
				Type: feature.Strategy_FIXED,
				FixedStrategy: &feature.FixedStrategy{
					Variation: v1.String(),
				},
			},
			OffVariation: v2.String(),
		},
	}

	testCases := []struct {
		name     string
		enabled  *wrapperspb.BoolValue
		archived *wrapperspb.BoolValue
		desc     string
	}{
		{
			name:     "enabled same value true",
			enabled:  wrapperspb.Bool(true),  // same as original
			archived: wrapperspb.Bool(false), // same as original
			desc:     "Should not update timestamp when enabled=true (same as original)",
		},
		{
			name:     "enabled same value false",
			enabled:  wrapperspb.Bool(false), // different from original
			archived: wrapperspb.Bool(false), // same as original
			desc:     "Should update timestamp when enabled=false (different from original)",
		},
		{
			name:     "archived same value false",
			enabled:  wrapperspb.Bool(true),  // same as original
			archived: wrapperspb.Bool(false), // same as original
			desc:     "Should not update timestamp when archived=false (same as original)",
		},
		{
			name:     "archived same value true",
			enabled:  wrapperspb.Bool(true), // same as original
			archived: wrapperspb.Bool(true), // different from original
			desc:     "Should update timestamp when archived=true (different from original)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh copy of original for each test
			testFeature := &Feature{
				Feature: &feature.Feature{
					Id:          original.Id,
					Name:        original.Name,
					Description: original.Description,
					Enabled:     original.Enabled,
					Archived:    original.Archived,
					Version:     original.Version,
					UpdatedAt:   original.UpdatedAt,
					Variations:  original.Variations,
					Targets:     original.Targets,
					DefaultStrategy: &feature.Strategy{
						Type: feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{
							Variation: original.DefaultStrategy.FixedStrategy.Variation,
						},
					},
					OffVariation: original.OffVariation,
				},
			}

			// Update with the test values
			updated, err := testFeature.Update(
				nil,         // name
				nil,         // description
				nil,         // tags
				tc.enabled,  // enabled
				tc.archived, // archived
				nil,         // defaultStrategy
				nil,         // offVariation
				false,       // resetSamplingSeed
				nil,         // prerequisiteChanges
				nil,         // targetChanges
				nil,         // ruleChanges
				nil,         // variationChanges
				nil,         // tagChanges
			)

			require.NoError(t, err)

			// Determine if values are changing
			enabledChanging := tc.enabled != nil && tc.enabled.Value != original.Enabled
			archivedChanging := tc.archived != nil && tc.archived.Value != original.Archived
			shouldHaveChanges := enabledChanging || archivedChanging

			if shouldHaveChanges {
				// Should have changes - version should increment and timestamp should update
				assert.Equal(t, original.Version+1, updated.Version, tc.desc)
				assert.NotEqual(t, original.UpdatedAt, updated.UpdatedAt, tc.desc)
				assert.True(t, updated.UpdatedAt > original.UpdatedAt, tc.desc)

				// Values should be updated
				if tc.enabled != nil {
					assert.Equal(t, tc.enabled.Value, updated.Enabled, tc.desc)
				}
				if tc.archived != nil {
					assert.Equal(t, tc.archived.Value, updated.Archived, tc.desc)
				}
			} else {
				// Should not have changes - version and timestamp should remain the same
				assert.Equal(t, original.Version, updated.Version, tc.desc)
				assert.Equal(t, original.UpdatedAt, updated.UpdatedAt, tc.desc)

				// Values should remain the same
				assert.Equal(t, original.Enabled, updated.Enabled, tc.desc)
				assert.Equal(t, original.Archived, updated.Archived, tc.desc)
			}
		})
	}
}

func TestUpdateWithIdenticalDefaultStrategy(t *testing.T) {
	t.Parallel()

	// Generate valid UUIDs for variations
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)

	// Create a feature with rollout strategy
	fixedTimestamp := time.Now().Unix() - 3600 // 1 hour ago
	original := &Feature{
		Feature: &feature.Feature{
			Id:        "test-feature",
			Name:      "Test Feature",
			Enabled:   true,
			Version:   10,
			UpdatedAt: fixedTimestamp,
			Variations: []*feature.Variation{
				{Id: v1.String(), Name: "v1", Value: "20"},
				{Id: v2.String(), Name: "v2", Value: "30"},
			},
			Targets: []*feature.Target{
				{Variation: v1.String(), Users: []string{}},
				{Variation: v2.String(), Users: []string{}},
			},
			DefaultStrategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{
						{Variation: v1.String(), Weight: 40000},
						{Variation: v2.String(), Weight: 60000},
					},
				},
			},
			OffVariation: v2.String(),
		},
	}

	// Update with identical default strategy
	identicalStrategy := &feature.Strategy{
		Type: feature.Strategy_ROLLOUT,
		RolloutStrategy: &feature.RolloutStrategy{
			Variations: []*feature.RolloutStrategy_Variation{
				{Variation: v1.String(), Weight: 40000},
				{Variation: v2.String(), Weight: 60000},
			},
		},
	}

	updated, err := original.Update(
		nil,                            // name
		nil,                            // description
		nil,                            // tags
		wrapperspb.Bool(true),          // enabled (same as original)
		nil,                            // archived
		identicalStrategy,              // defaultStrategy (identical to original)
		wrapperspb.String(v2.String()), // offVariation (same as original)
		false,                          // resetSamplingSeed
		nil,                            // prerequisiteChanges
		nil,                            // targetChanges
		nil,                            // ruleChanges
		nil,                            // variationChanges
		nil,                            // tagChanges
	)

	require.NoError(t, err)

	// Should not have changes - version and timestamp should remain the same
	assert.Equal(t, original.Version, updated.Version, "Version should not increment when no actual changes")
	assert.Equal(t, original.UpdatedAt, updated.UpdatedAt, "Timestamp should not change when no actual changes")

	// Values should remain the same
	assert.Equal(t, original.Enabled, updated.Enabled)
	assert.Equal(t, original.DefaultStrategy.RolloutStrategy.Variations, updated.DefaultStrategy.RolloutStrategy.Variations)
	assert.Equal(t, original.OffVariation, updated.OffVariation)
}

func TestUpdateEnable(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		expectedFunc func() *Feature
	}{
		{
			desc: "enable when already enabled - no-op",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = true
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = true
				return f
			},
		},
		{
			desc: "enable when disabled - should enable",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = false
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = true
				return f
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateEnable()
			require.NoError(t, err)
			expected := p.expectedFunc()
			assert.Equal(t, expected.Enabled, actual.Enabled)
		})
	}
}

func TestUpdateDisable(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		expectedFunc func() *Feature
	}{
		{
			desc: "disable when already disabled - no-op",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = false
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = false
				return f
			},
		},
		{
			desc: "disable when enabled - should disable",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = true
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Enabled = false
				return f
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateDisable()
			require.NoError(t, err)
			expected := p.expectedFunc()
			assert.Equal(t, expected.Enabled, actual.Enabled)
		})
	}
}

func TestUpdateArchive(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		expectedFunc func() *Feature
	}{
		{
			desc: "archive when already archived - no-op",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = true
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = true
				return f
			},
		},
		{
			desc: "archive when not archived - should archive",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = false
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = true
				return f
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateArchive()
			require.NoError(t, err)
			expected := p.expectedFunc()
			assert.Equal(t, expected.Archived, actual.Archived)
		})
	}
}

func TestUpdateUnarchive(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		expectedFunc func() *Feature
	}{
		{
			desc: "unarchive when already unarchived - no-op",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = false
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = false
				return f
			},
		},
		{
			desc: "unarchive when archived - should unarchive",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = true
				return f
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Archived = false
				return f
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateUnarchive()
			require.NoError(t, err)
			expected := p.expectedFunc()
			assert.Equal(t, expected.Archived, actual.Archived)
		})
	}
}

func TestUpdateAddVariation(t *testing.T) {
	newV, err := uuid.NewUUID()
	require.NoError(t, err)

	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		id           string
		value        string
		name         string
		description  string
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - add new variation",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			id:          newV.String(),
			value:       "new-value",
			name:        "new-name",
			description: "new-description",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Variations = append(f.Variations, &feature.Variation{
					Id:          newV.String(),
					Value:       "new-value",
					Name:        "new-name",
					Description: "new-description",
				})
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "error - empty id",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			id:          "",
			value:       "new-value",
			name:        "new-name",
			description: "new-description",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errVariationIDRequired,
		},
		{
			desc: "error - empty value",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			id:          newV.String(),
			value:       "",
			name:        "new-name",
			description: "new-description",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errVariationValueRequired,
		},
		{
			desc: "error - empty name",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			id:          newV.String(),
			value:       "new-value",
			name:        "",
			description: "new-description",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errVariationNameRequired,
		},
		{
			desc: "error - duplicate variation id",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				return f
			},
			id:          "variation-A", // Using existing variation ID from makeFeature
			value:       "new-value",
			name:        "new-name",
			description: "new-description",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errVariationValueUnique,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateAddVariation(p.id, p.value, p.name, p.description)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				expected := p.expectedFunc()
				assert.Equal(t, len(expected.Variations), len(actual.Variations))
				if len(actual.Variations) > 0 {
					lastVar := actual.Variations[len(actual.Variations)-1]
					assert.Equal(t, p.id, lastVar.Id)
					assert.Equal(t, p.value, lastVar.Value)
					assert.Equal(t, p.name, lastVar.Name)
					assert.Equal(t, p.description, lastVar.Description)
				}
			}
		})
	}
}

func TestUpdateChangeVariation(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		variation    *feature.Variation
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - change variation",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			variation: &feature.Variation{
				Id:          "variation-A", // Using existing variation ID from makeFeature
				Value:       "updated-value",
				Name:        "updated-name",
				Description: "updated-description",
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Variations[0].Value = "updated-value"
				f.Variations[0].Name = "updated-name"
				f.Variations[0].Description = "updated-description"
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "no change - same variation",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				return f
			},
			variation: &feature.Variation{
				Id:          "variation-A",  // Using existing variation ID from makeFeature
				Value:       "A",            // Original value from makeFeature
				Name:        "Variation A",  // Original name from makeFeature
				Description: "Thing does A", // Original description from makeFeature
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: nil,
		},
		{
			desc: "error - nil variation",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			variation: nil,
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errVariationRequired,
		},
		{
			desc: "error - empty name",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			variation: &feature.Variation{
				Id:          "variation-A", // Using existing variation ID from makeFeature
				Value:       "updated-value",
				Name:        "",
				Description: "updated-description",
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errVariationNameRequired,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateChangeVariation(p.variation)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				if p.variation != nil {
					idx, findErr := actual.findVariationIndex(p.variation.Id)
					require.NoError(t, findErr)
					if p.desc == "success - change variation" {
						assert.Equal(t, p.variation.Value, actual.Variations[idx].Value)
						assert.Equal(t, p.variation.Name, actual.Variations[idx].Name)
						assert.Equal(t, p.variation.Description, actual.Variations[idx].Description)
					}
				}
			}
		})
	}
}

func TestUpdateAddTargetUsers(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		target       *feature.Target
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - add new users",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			target: &feature.Target{
				Variation: "variation-A", // Using existing variation ID from makeFeature
				Users:     []string{"new-user1", "new-user2"},
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Targets[0].Users = []string{"user1", "new-user1", "new-user2"} // makeFeature creates user1 initially
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "no-op - add existing users",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				return f // makeFeature already has user1 in variation-A
			},
			target: &feature.Target{
				Variation: "variation-A",     // Using existing variation ID from makeFeature
				Users:     []string{"user1"}, // user1 already exists in makeFeature
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				return f // No change expected
			},
			expectedErr: nil,
		},
		{
			desc: "error - nil users",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			target: &feature.Target{
				Variation: "variation-A", // Using existing variation ID from makeFeature
				Users:     nil,
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errTargetUsersRequired,
		},
		{
			desc: "error - empty user",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			target: &feature.Target{
				Variation: "variation-A", // Using existing variation ID from makeFeature
				Users:     []string{""},
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errTargetUserRequired,
		},
		{
			desc: "error - target not found",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			target: &feature.Target{
				Variation: "non-existent",
				Users:     []string{"user1"},
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errTargetNotFound,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateAddTargetUsers(p.target)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				idx, findErr := actual.findTarget(p.target.Variation)
				require.NoError(t, findErr)
				expected := p.expectedFunc()
				expectedIdx, expectedFindErr := expected.findTarget(p.target.Variation)
				require.NoError(t, expectedFindErr)
				assert.ElementsMatch(t, expected.Targets[expectedIdx].Users, actual.Targets[idx].Users)
			}
		})
	}
}

func TestUpdateRemoveTargetUsers(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		target       *feature.Target
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - remove existing users",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Targets[0].Users = []string{"user1", "user2", "user3"}
				return f
			},
			target: &feature.Target{
				Variation: "variation-A", // Using existing variation ID from makeFeature
				Users:     []string{"user1", "user3"},
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Targets[0].Users = []string{"user2"}
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "no-op - remove non-existent users",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Targets[0].Users = []string{"user1"}
				return f
			},
			target: &feature.Target{
				Variation: "variation-A", // Using existing variation ID from makeFeature
				Users:     []string{"user2"},
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Targets[0].Users = []string{"user1"}
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "error - nil users",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			target: &feature.Target{
				Variation: "variation-A", // Using existing variation ID from makeFeature
				Users:     nil,
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errTargetUsersRequired,
		},
		{
			desc: "error - empty user",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			target: &feature.Target{
				Variation: "variation-A", // Using existing variation ID from makeFeature
				Users:     []string{""},
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errTargetUserRequired,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateRemoveTargetUsers(p.target)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				idx, findErr := actual.findTarget(p.target.Variation)
				require.NoError(t, findErr)
				expected := p.expectedFunc()
				expectedIdx, expectedFindErr := expected.findTarget(p.target.Variation)
				require.NoError(t, expectedFindErr)
				assert.ElementsMatch(t, expected.Targets[expectedIdx].Users, actual.Targets[idx].Users)
			}
		})
	}
}

func TestUpdateAddRule(t *testing.T) {
	ruleID, err := uuid.NewUUID()
	require.NoError(t, err)
	clauseID, err := uuid.NewUUID()
	require.NoError(t, err)

	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		rule         *feature.Rule
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - add new rule",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			rule: &feature.Rule{
				Id: ruleID.String(),
				Strategy: &feature.Strategy{
					Type:          feature.Strategy_FIXED,
					FixedStrategy: &feature.FixedStrategy{Variation: "variation-A"}, // Using existing variation ID from makeFeature
				},
				Clauses: []*feature.Clause{
					{
						Id:        clauseID.String(),
						Attribute: "attr",
						Operator:  feature.Clause_EQUALS,
						Values:    []string{"val"},
					},
				},
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Rules = append(f.Rules, &feature.Rule{
					Id: ruleID.String(),
					Strategy: &feature.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{Variation: "variation-A"}, // Using existing variation ID from makeFeature
					},
					Clauses: []*feature.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"val"},
						},
					},
				})
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "error - nil rule",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			rule: nil,
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errRuleRequired,
		},
		{
			desc: "error - nil strategy",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			rule: &feature.Rule{
				Id:       ruleID.String(),
				Strategy: nil,
				Clauses: []*feature.Clause{
					{
						Id:        clauseID.String(),
						Attribute: "attr",
						Operator:  feature.Clause_EQUALS,
						Values:    []string{"val"},
					},
				},
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errStrategyRequired,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateAddRule(p.rule)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				expected := p.expectedFunc()
				assert.Equal(t, len(expected.Rules), len(actual.Rules))
				if len(actual.Rules) > 0 && p.rule != nil {
					lastRule := actual.Rules[len(actual.Rules)-1]
					assert.Equal(t, p.rule.Id, lastRule.Id)
				}
			}
		})
	}
}

func TestUpdateChangeRule(t *testing.T) {
	ruleID, err := uuid.NewUUID()
	require.NoError(t, err)
	clauseID, err := uuid.NewUUID()
	require.NoError(t, err)

	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		rule         *feature.Rule
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - change rule",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Rules = append(f.Rules, &feature.Rule{
					Id: ruleID.String(),
					Strategy: &feature.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{Variation: "variation-A"}, // Using existing variation ID from makeFeature
					},
					Clauses: []*feature.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"val"},
						},
					},
				})
				return f
			},
			rule: &feature.Rule{
				Id: ruleID.String(),
				Strategy: &feature.Strategy{
					Type:          feature.Strategy_FIXED,
					FixedStrategy: &feature.FixedStrategy{Variation: "variation-A"}, // Using existing variation ID from makeFeature
				},
				Clauses: []*feature.Clause{
					{
						Id:        clauseID.String(),
						Attribute: "attr",
						Operator:  feature.Clause_EQUALS,
						Values:    []string{"updated-val"},
					},
				},
			},
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Rules = append(f.Rules, &feature.Rule{
					Id: ruleID.String(),
					Strategy: &feature.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{Variation: "variation-A"}, // Using existing variation ID from makeFeature
					},
					Clauses: []*feature.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"updated-val"},
						},
					},
				})
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "error - nil rule",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			rule: nil,
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errRuleRequired,
		},
		{
			desc: "error - rule not found",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			rule: &feature.Rule{
				Id: "non-existent",
				Strategy: &feature.Strategy{
					Type:          feature.Strategy_FIXED,
					FixedStrategy: &feature.FixedStrategy{Variation: "variation-A"}, // Using existing variation ID from makeFeature
				},
				Clauses: []*feature.Clause{
					{
						Id:        clauseID.String(),
						Attribute: "attr",
						Operator:  feature.Clause_EQUALS,
						Values:    []string{"val"},
					},
				},
			},
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errRuleNotFound,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateChangeRule(p.rule)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				if p.rule != nil {
					idx, findErr := actual.findRuleIndex(p.rule.Id)
					require.NoError(t, findErr)
					assert.Equal(t, p.rule.Id, actual.Rules[idx].Id)
				}
			}
		})
	}
}

func TestUpdateRemoveRule(t *testing.T) {
	ruleID, err := uuid.NewUUID()
	require.NoError(t, err)
	clauseID, err := uuid.NewUUID()
	require.NoError(t, err)

	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		ruleID       string
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - remove existing rule",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Rules = append(f.Rules, &feature.Rule{
					Id: ruleID.String(),
					Strategy: &feature.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{Variation: "variation-A"}, // Using existing variation ID from makeFeature
					},
					Clauses: []*feature.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"val"},
						},
					},
				})
				return f
			},
			ruleID: ruleID.String(),
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: nil,
		},
		{
			desc: "error - rule not found",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			ruleID: "non-existent",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errRuleNotFound,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateRemoveRule(p.ruleID)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				expected := p.expectedFunc()
				assert.Equal(t, len(expected.Rules), len(actual.Rules))
			}
		})
	}
}

func TestUpdateAddPrerequisite(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		featureID    string
		variationID  string
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - add new prerequisite",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			featureID:   "feature-1",
			variationID: "variation-1",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
					FeatureId:   "feature-1",
					VariationId: "variation-1",
				})
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "error - empty feature ID",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			featureID:   "",
			variationID: "variation-1",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errFeatureIDRequired,
		},
		{
			desc: "error - empty variation ID",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			featureID:   "feature-1",
			variationID: "",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errVariationIDRequired,
		},
		{
			desc: "error - prerequisite already exists",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
					FeatureId:   "feature-1",
					VariationId: "variation-1",
				})
				return f
			},
			featureID:   "feature-1",
			variationID: "variation-2",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errPrerequisiteAlreadyExists,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateAddPrerequisite(p.featureID, p.variationID)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				expected := p.expectedFunc()
				assert.Equal(t, len(expected.Prerequisites), len(actual.Prerequisites))
				if len(actual.Prerequisites) > 0 {
					lastPrereq := actual.Prerequisites[len(actual.Prerequisites)-1]
					assert.Equal(t, p.featureID, lastPrereq.FeatureId)
					assert.Equal(t, p.variationID, lastPrereq.VariationId)
				}
			}
		})
	}
}

func TestUpdateChangePrerequisiteVariation(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		featureID    string
		variationID  string
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - change prerequisite variation",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
					FeatureId:   "feature-1",
					VariationId: "variation-1",
				})
				return f
			},
			featureID:   "feature-1",
			variationID: "variation-2",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
					FeatureId:   "feature-1",
					VariationId: "variation-2",
				})
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "no-op - same prerequisite variation",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
					FeatureId:   "feature-1",
					VariationId: "variation-1",
				})
				return f
			},
			featureID:   "feature-1",
			variationID: "variation-1",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
					FeatureId:   "feature-1",
					VariationId: "variation-1",
				})
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "error - prerequisite not found",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			featureID:   "non-existent",
			variationID: "variation-1",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errPrerequisiteNotFound,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateChangePrerequisiteVariation(p.featureID, p.variationID)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				if p.expectedErr == nil {
					idx, findErr := actual.findPrerequisiteIndex(p.featureID)
					require.NoError(t, findErr)
					assert.Equal(t, p.variationID, actual.Prerequisites[idx].VariationId)
				}
			}
		})
	}
}

func TestUpdateRemovePrerequisite(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		featureID    string
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - remove existing prerequisite",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
					FeatureId:   "feature-1",
					VariationId: "variation-1",
				})
				return f
			},
			featureID: "feature-1",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: nil,
		},
		{
			desc: "error - prerequisite not found",
			inputFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			featureID: "non-existent",
			expectedFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedErr: errPrerequisiteNotFound,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateRemovePrerequisite(p.featureID)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				expected := p.expectedFunc()
				assert.Equal(t, len(expected.Prerequisites), len(actual.Prerequisites))
			}
		})
	}
}

func TestUpdateAddTag(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		tag          string
		expectedFunc func() *Feature
	}{
		{
			desc: "success - add new tag",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1"}
				return f
			},
			tag: "tag2",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1", "tag2"}
				return f
			},
		},
		{
			desc: "no-op - add existing tag",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1"}
				return f
			},
			tag: "tag1",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1"}
				return f
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateAddTag(p.tag)
			require.NoError(t, err)
			expected := p.expectedFunc()
			assert.ElementsMatch(t, expected.Tags, actual.Tags)
		})
	}
}

func TestUpdateRemoveTag(t *testing.T) {
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		tag          string
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc: "success - remove existing tag",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1", "tag2"}
				return f
			},
			tag: "tag1",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag2"}
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "error - remove non-existent tag",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1"}
				return f
			},
			tag: "tag2",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1"}
				return f
			},
			expectedErr: errors.New("feature: tag not found"),
		},
		{
			desc: "success - remove last tag",
			inputFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{"tag1"}
				return f
			},
			tag: "tag1",
			expectedFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.Tags = []string{}
				return f
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.inputFunc()
			err := actual.updateRemoveTag(p.tag)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				require.NoError(t, err)
				expected := p.expectedFunc()
				assert.ElementsMatch(t, expected.Tags, actual.Tags)
			}
		})
	}
}

func TestUpdateCompleteNoChangesScenario(t *testing.T) {
	t.Parallel()

	// Create a comprehensive feature with all possible fields set
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)
	ruleID, err := uuid.NewUUID()
	require.NoError(t, err)
	clauseID, err := uuid.NewUUID()
	require.NoError(t, err)

	fixedTimestamp := time.Now().Unix() - 3600 // 1 hour ago
	originalVersion := int32(10)

	originalFeature := &Feature{
		Feature: &feature.Feature{
			Id:            "comprehensive-test-feature",
			Name:          "Comprehensive Test Feature",
			Description:   "A feature with all possible fields set",
			Tags:          []string{"tag1", "tag2", "tag3"},
			Enabled:       true,
			Archived:      false,
			Version:       originalVersion,
			UpdatedAt:     fixedTimestamp,
			VariationType: feature.Feature_BOOLEAN,
			Variations: []*feature.Variation{
				{Id: v1.String(), Name: "v1", Value: "true", Description: "First variation"},
				{Id: v2.String(), Name: "v2", Value: "false", Description: "Second variation"},
			},
			Targets: []*feature.Target{
				{Variation: v1.String(), Users: []string{"user1", "user2"}},
				{Variation: v2.String(), Users: []string{"user3"}},
			},
			Rules: []*feature.Rule{
				{
					Id: ruleID.String(),
					Strategy: &feature.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{Variation: v1.String()},
					},
					Clauses: []*feature.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "user_type",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"premium"},
						},
					},
				},
			},
			Prerequisites: []*feature.Prerequisite{
				{FeatureId: "prerequisite-feature-1", VariationId: v1.String()},
			},
			DefaultStrategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{
						{Variation: v1.String(), Weight: 30000},
						{Variation: v2.String(), Weight: 70000},
					},
				},
			},
			OffVariation: v2.String(),
			SamplingSeed: "test-sampling-seed",
		},
	}

	// Call Update with IDENTICAL values to all current fields
	updated, err := originalFeature.Update(
		// Basic fields - all identical to original
		wrapperspb.String("Comprehensive Test Feature"),                   // name - same
		wrapperspb.String("A feature with all possible fields set"),       // description - same
		&common.StringListValue{Values: []string{"tag1", "tag2", "tag3"}}, // tags - same order
		wrapperspb.Bool(true),  // enabled - same
		wrapperspb.Bool(false), // archived - same
		&feature.Strategy{ // defaultStrategy - identical
			Type: feature.Strategy_ROLLOUT,
			RolloutStrategy: &feature.RolloutStrategy{
				Variations: []*feature.RolloutStrategy_Variation{
					{Variation: v1.String(), Weight: 30000},
					{Variation: v2.String(), Weight: 70000},
				},
			},
		},
		wrapperspb.String(v2.String()), // offVariation - same
		false,                          // resetSamplingSeed - no reset

		// Granular changes - all empty (no changes)
		nil, // prerequisiteChanges - no changes
		nil, // targetChanges - no changes
		nil, // ruleChanges - no changes
		nil, // variationChanges - no changes
		nil, // tagChanges - no changes
	)

	require.NoError(t, err)

	// CRITICAL ASSERTIONS: Version and timestamp should NOT change
	assert.Equal(t, originalVersion, updated.Version, "Version should NOT increment when no actual changes occur")
	assert.Equal(t, fixedTimestamp, updated.UpdatedAt, "UpdatedAt should NOT change when no actual changes occur")

	// Verify all field values remain exactly the same
	assert.Equal(t, originalFeature.Id, updated.Id)
	assert.Equal(t, originalFeature.Name, updated.Name)
	assert.Equal(t, originalFeature.Description, updated.Description)
	assert.ElementsMatch(t, originalFeature.Tags, updated.Tags)
	assert.Equal(t, originalFeature.Enabled, updated.Enabled)
	assert.Equal(t, originalFeature.Archived, updated.Archived)
	assert.Equal(t, originalFeature.OffVariation, updated.OffVariation)
	assert.Equal(t, originalFeature.SamplingSeed, updated.SamplingSeed)

	// Deep comparison of complex fields
	assert.Equal(t, originalFeature.Variations, updated.Variations)
	assert.Equal(t, originalFeature.Targets, updated.Targets)
	assert.Equal(t, originalFeature.Rules, updated.Rules)
	assert.Equal(t, originalFeature.Prerequisites, updated.Prerequisites)

	// Compare default strategy
	assert.True(t, proto.Equal(originalFeature.DefaultStrategy, updated.DefaultStrategy),
		"Default strategy should remain identical")
}

// TestUpdateWithActualChangesIncrementsVersionAndTimestamp verifies the opposite scenario
func TestUpdateWithActualChangesIncrementsVersionAndTimestamp(t *testing.T) {
	t.Parallel()

	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)

	fixedTimestamp := time.Now().Unix() - 3600 // 1 hour ago
	originalVersion := int32(5)

	originalFeature := &Feature{
		Feature: &feature.Feature{
			Id:          "test-feature",
			Name:        "Original Name",
			Description: "Original Description",
			Enabled:     false,
			Archived:    false,
			Version:     originalVersion,
			UpdatedAt:   fixedTimestamp,
			Variations: []*feature.Variation{
				{Id: v1.String(), Name: "v1", Value: "true"},
				{Id: v2.String(), Name: "v2", Value: "false"},
			},
			Targets: []*feature.Target{
				{Variation: v1.String(), Users: []string{}},
				{Variation: v2.String(), Users: []string{}},
			},
		},
	}

	// Make an actual change (different name)
	updated, err := originalFeature.Update(
		wrapperspb.String("Updated Name"), // CHANGED - different from original
		nil, nil, nil, nil, nil, nil, false, nil, nil, nil, nil, nil,
	)

	require.NoError(t, err)

	// CRITICAL ASSERTIONS: Version and timestamp SHOULD change
	assert.Equal(t, originalVersion+1, updated.Version, "Version should increment when actual changes occur")
	assert.NotEqual(t, fixedTimestamp, updated.UpdatedAt, "UpdatedAt should change when actual changes occur")
	assert.True(t, updated.UpdatedAt > fixedTimestamp, "UpdatedAt should be more recent")

	// Verify the change took effect
	assert.Equal(t, "Updated Name", updated.Name)
}

func TestUpdateRemoveVariationComprehensiveCleanup(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setupFunc   func() *Feature
		variationID string
		expectedErr error
	}{
		{
			desc: "success - comprehensive cleanup",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				expected := "variation-D"
				f.AddVariation(expected, "D", "Variation D", "Thing does D")

				// Set up rollout strategy in default strategy with the variation (weight=0 so it can be removed)
				f.ChangeDefaultStrategy(&ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{
								Variation: "variation-A",
								Weight:    100000,
							},
							{
								Variation: expected,
								Weight:    0, // Weight 0 means not "in use" so can be removed
							},
						},
					},
				})

				// Add a rule with rollout strategy containing the variation
				rule := &ftproto.Rule{
					Id: "test-rule-rollout",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_ROLLOUT,
						RolloutStrategy: &ftproto.RolloutStrategy{
							Variations: []*ftproto.RolloutStrategy_Variation{
								{
									Variation: "variation-B",
									Weight:    50000,
								},
								{
									Variation: expected,
									Weight:    0, // Weight 0 means not "in use" so can be removed
								},
							},
						},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        "clause-1",
							Attribute: "user_id",
							Operator:  ftproto.Clause_EQUALS,
							Values:    []string{"user-1"},
						},
					},
				}
				f.AddRule(rule)
				return f
			},
			variationID: "variation-D",
			expectedErr: nil,
		},
		{
			desc: "error - variation in use",
			setupFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			variationID: "variation-C", // Has users in target
			expectedErr: errVariationInUse,
		},
		{
			desc: "error - minimum variation constraint",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Remove one variation to get down to exactly 2 variations
				f.Targets[2].Users = []string{}        // Remove users from variation-C
				f.updateRemoveVariation("variation-C") // This should succeed, leaving 2 variations
				return f
			},
			variationID: "variation-A", // Try to remove when only 2 variations remain
			expectedErr: errVariationsMustHaveAtLeastTwoVariations,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f := p.setupFunc()
			err := f.updateRemoveVariation(p.variationID)
			assert.Equal(t, p.expectedErr, err)

			if err == nil {
				// Verify complete cleanup for successfully removed variation
				if _, err := f.updateFindVariationIndex(p.variationID); err == nil {
					t.Fatalf("Variation not deleted from Variations. Actual: %v", f.Variations)
				}
				if _, err := f.updateFindTarget(p.variationID); err == nil {
					t.Fatalf("Target not deleted. Actual: %v", f.Targets)
				}

				// Verify variation removed from default strategy rollout
				for _, v := range f.DefaultStrategy.RolloutStrategy.Variations {
					if v.Variation == p.variationID {
						t.Fatalf("Variation not removed from default strategy. Actual: %v", f.DefaultStrategy.RolloutStrategy.Variations)
					}
				}

				// Verify variation removed from rule rollout strategy
				for _, r := range f.Rules {
					if r.Id == "test-rule-rollout" && r.Strategy.Type == ftproto.Strategy_ROLLOUT {
						for _, v := range r.Strategy.RolloutStrategy.Variations {
							if v.Variation == p.variationID {
								t.Fatalf("Variation not removed from rule strategy. Actual: %v", r.Strategy.RolloutStrategy.Variations)
							}
						}
					}
				}
			}
		})
	}
}

func TestUpdateRemoveVariationMultipleRulesCleanup(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setupFunc   func() *Feature
		variationID string
		expectedErr error
	}{
		{
			desc: "success - multiple rules cleanup",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				expected := "variation-D"
				f.AddVariation(expected, "D", "Variation D", "Thing does D")

				// Clear existing rules and add new ones with rollout strategies
				f.Rules = []*ftproto.Rule{}

				// Add multiple rules with rollout strategies containing the variation
				rule1 := &ftproto.Rule{
					Id: "test-rule-1",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_ROLLOUT,
						RolloutStrategy: &ftproto.RolloutStrategy{
							Variations: []*ftproto.RolloutStrategy_Variation{
								{
									Variation: "variation-A",
									Weight:    50000,
								},
								{
									Variation: expected,
									Weight:    0, // Weight 0 means not "in use"
								},
							},
						},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        "0efe416e-2fd2-4996-b5c3-194f05444f1f",
							Attribute: "user_id",
							Operator:  ftproto.Clause_EQUALS,
							Values:    []string{"user-1"},
						},
					},
				}
				f.AddRule(rule1)

				rule2 := &ftproto.Rule{
					Id: "test-rule-2",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_ROLLOUT,
						RolloutStrategy: &ftproto.RolloutStrategy{
							Variations: []*ftproto.RolloutStrategy_Variation{
								{
									Variation: "variation-B",
									Weight:    70000,
								},
								{
									Variation: expected,
									Weight:    0, // Weight 0 means not "in use"
								},
							},
						},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        "1efe416e-2fd2-4996-b5c3-194f05444f1f",
							Attribute: "user_type",
							Operator:  ftproto.Clause_EQUALS,
							Values:    []string{"premium"},
						},
					},
				}
				f.AddRule(rule2)
				return f
			},
			variationID: "variation-D",
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f := p.setupFunc()

			// Verify both rules have the variation before removal
			rule1HasVariation := false
			rule2HasVariation := false
			for _, r := range f.Rules {
				if r.Id == "test-rule-1" && r.Strategy.Type == ftproto.Strategy_ROLLOUT {
					for _, v := range r.Strategy.RolloutStrategy.Variations {
						if v.Variation == p.variationID {
							rule1HasVariation = true
							break
						}
					}
				}
				if r.Id == "test-rule-2" && r.Strategy.Type == ftproto.Strategy_ROLLOUT {
					for _, v := range r.Strategy.RolloutStrategy.Variations {
						if v.Variation == p.variationID {
							rule2HasVariation = true
							break
						}
					}
				}
			}
			if !rule1HasVariation {
				t.Fatalf("Rule 1 should have variation before removal. Looking for: %s", p.variationID)
			}
			if !rule2HasVariation {
				t.Fatalf("Rule 2 should have variation before removal. Looking for: %s", p.variationID)
			}

			err := f.updateRemoveVariation(p.variationID)
			assert.Equal(t, p.expectedErr, err)

			if err == nil {
				// Verify variation is removed from BOTH rules (this would fail with the early return bug)
				rule1HasVariation = false
				rule2HasVariation = false
				for _, r := range f.Rules {
					if r.Id == "test-rule-1" && r.Strategy.Type == ftproto.Strategy_ROLLOUT {
						for _, v := range r.Strategy.RolloutStrategy.Variations {
							if v.Variation == p.variationID {
								rule1HasVariation = true
								break
							}
						}
					}
					if r.Id == "test-rule-2" && r.Strategy.Type == ftproto.Strategy_ROLLOUT {
						for _, v := range r.Strategy.RolloutStrategy.Variations {
							if v.Variation == p.variationID {
								rule2HasVariation = true
								break
							}
						}
					}
				}
				if rule1HasVariation {
					t.Fatalf("Rule 1 should NOT have variation after removal")
				}
				if rule2HasVariation {
					t.Fatalf("Rule 2 should NOT have variation after removal")
				}
			}
		})
	}
}
