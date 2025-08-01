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
	"testing"

	"github.com/stretchr/testify/assert"

	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestValidateVariationReferences(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc                string
		setupFunc           func() *Feature
		expectedOrphanedIDs []string
	}{
		{
			desc: "no orphaned references",
			setupFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedOrphanedIDs: []string{},
		},
		{
			desc: "orphaned target reference",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Add orphaned target
				f.Targets = append(f.Targets, &ftproto.Target{
					Variation: "orphaned-variation-1",
					Users:     []string{"user1"},
				})
				return f
			},
			expectedOrphanedIDs: []string{"orphaned-variation-1"},
		},
		{
			desc: "orphaned variation in rule rollout strategy",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Add rule with orphaned variation
				rule := &ftproto.Rule{
					Id: "test-rule",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_ROLLOUT,
						RolloutStrategy: &ftproto.RolloutStrategy{
							Variations: []*ftproto.RolloutStrategy_Variation{
								{
									Variation: "variation-A", // Valid
									Weight:    50000,
								},
								{
									Variation: "orphaned-variation-2", // Orphaned
									Weight:    0,
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
				f.Rules = []*ftproto.Rule{rule}
				return f
			},
			expectedOrphanedIDs: []string{"orphaned-variation-2"},
		},
		{
			desc: "orphaned variation in default strategy",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Set default strategy with orphaned variation
				f.DefaultStrategy = &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{
								Variation: "variation-A", // Valid
								Weight:    100000,
							},
							{
								Variation: "orphaned-variation-3", // Orphaned
								Weight:    0,
							},
						},
					},
				}
				return f
			},
			expectedOrphanedIDs: []string{"orphaned-variation-3"},
		},
		{
			desc: "orphaned OffVariation",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.OffVariation = "orphaned-off-variation"
				return f
			},
			expectedOrphanedIDs: []string{"orphaned-off-variation"},
		},
		{
			desc: "multiple orphaned references",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")

				// Orphaned target
				f.Targets = append(f.Targets, &ftproto.Target{
					Variation: "orphaned-1",
					Users:     []string{"user1"},
				})

				// Orphaned in rule
				rule := &ftproto.Rule{
					Id: "test-rule",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_ROLLOUT,
						RolloutStrategy: &ftproto.RolloutStrategy{
							Variations: []*ftproto.RolloutStrategy_Variation{
								{
									Variation: "orphaned-2",
									Weight:    0,
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
				f.Rules = []*ftproto.Rule{rule}

				// Orphaned in default strategy
				f.DefaultStrategy = &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{
								Variation: "orphaned-3",
								Weight:    0,
							},
						},
					},
				}

				return f
			},
			expectedOrphanedIDs: []string{"orphaned-1", "orphaned-2", "orphaned-3"},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f := p.setupFunc()

			orphanedIDs := f.ValidateVariationReferences()

			// Convert to map for easier comparison (order doesn't matter)
			expectedMap := make(map[string]bool)
			for _, id := range p.expectedOrphanedIDs {
				expectedMap[id] = true
			}

			actualMap := make(map[string]bool)
			for _, id := range orphanedIDs {
				actualMap[id] = true
			}

			assert.Equal(t, expectedMap, actualMap, "Orphaned variation IDs don't match")
		})
	}
}

func TestCleanupOrphanedVariationReferences(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc            string
		setupFunc       func() *Feature
		expectedChanged bool
		verifyFunc      func(*testing.T, *Feature)
	}{
		{
			desc: "no cleanup needed",
			setupFunc: func() *Feature {
				return makeFeature("test-feature")
			},
			expectedChanged: false,
			verifyFunc: func(t *testing.T, f *Feature) {
				// Should remain unchanged
				assert.Equal(t, 3, len(f.Variations))
				assert.Equal(t, 3, len(f.Targets))
			},
		},
		{
			desc: "cleanup orphaned target",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Add orphaned target
				f.Targets = append(f.Targets, &ftproto.Target{
					Variation: "orphaned-variation",
					Users:     []string{"user1"},
				})
				return f
			},
			expectedChanged: true,
			verifyFunc: func(t *testing.T, f *Feature) {
				// Orphaned target should be removed
				assert.Equal(t, 3, len(f.Targets))
				for _, target := range f.Targets {
					assert.NotEqual(t, "orphaned-variation", target.Variation)
				}
			},
		},
		{
			desc: "cleanup orphaned variation in rule",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Clear existing rules and add one with orphaned variation
				f.Rules = []*ftproto.Rule{
					{
						Id: "test-rule",
						Strategy: &ftproto.Strategy{
							Type: ftproto.Strategy_ROLLOUT,
							RolloutStrategy: &ftproto.RolloutStrategy{
								Variations: []*ftproto.RolloutStrategy_Variation{
									{
										Variation: "variation-A", // Valid
										Weight:    50000,
									},
									{
										Variation: "orphaned-variation", // Orphaned
										Weight:    0,
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
					},
				}
				return f
			},
			expectedChanged: true,
			verifyFunc: func(t *testing.T, f *Feature) {
				// Orphaned variation should be removed from rule
				assert.Equal(t, 1, len(f.Rules))
				rule := f.Rules[0]
				assert.Equal(t, 1, len(rule.Strategy.RolloutStrategy.Variations))
				assert.Equal(t, "variation-A", rule.Strategy.RolloutStrategy.Variations[0].Variation)
			},
		},
		{
			desc: "cleanup orphaned variation in default strategy",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Set default strategy with orphaned variation
				f.DefaultStrategy = &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{
								Variation: "variation-A", // Valid
								Weight:    100000,
							},
							{
								Variation: "orphaned-variation", // Orphaned
								Weight:    0,
							},
						},
					},
				}
				return f
			},
			expectedChanged: true,
			verifyFunc: func(t *testing.T, f *Feature) {
				// Orphaned variation should be removed from default strategy
				assert.Equal(t, 1, len(f.DefaultStrategy.RolloutStrategy.Variations))
				assert.Equal(t, "variation-A", f.DefaultStrategy.RolloutStrategy.Variations[0].Variation)
			},
		},
		{
			desc: "cleanup orphaned OffVariation",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.OffVariation = "orphaned-off-variation"
				return f
			},
			expectedChanged: true,
			verifyFunc: func(t *testing.T, f *Feature) {
				// OffVariation should be reset to second valid variation
				assert.Equal(t, "variation-B", f.OffVariation)
			},
		},
		{
			desc: "cleanup orphaned OffVariation with only one variation",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Keep only one variation
				f.Variations = f.Variations[:1]
				f.OffVariation = "orphaned-off-variation"
				return f
			},
			expectedChanged: true,
			verifyFunc: func(t *testing.T, f *Feature) {
				// OffVariation should fallback to first variation when only one exists
				assert.Equal(t, "variation-A", f.OffVariation)
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f := p.setupFunc()

			initialTime := f.UpdatedAt
			result := f.CleanupOrphanedVariationReferences()

			assert.Equal(t, p.expectedChanged, result.Changed, "Changed flag doesn't match expected")

			if result.Changed {
				assert.Greater(t, f.UpdatedAt, initialTime, "UpdatedAt should be updated when changes are made")
			} else {
				assert.Equal(t, initialTime, f.UpdatedAt, "UpdatedAt should not change when no changes are made")
			}

			p.verifyFunc(t, f)
		})
	}
}

func TestEnsureVariationsInStrategies(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc           string
		setupFunc      func() *Feature
		expectedResult VariationMigrationResult
	}{
		{
			desc: "no changes needed - all variations already in strategies",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.DefaultStrategy = &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{Variation: "variation-A", Weight: 50000},
							{Variation: "variation-B", Weight: 30000},
							{Variation: "variation-C", Weight: 20000},
						},
					},
				}
				return f
			},
			expectedResult: VariationMigrationResult{
				Changed:           false,
				AddedToRules:      0,
				AddedToDefault:    0,
				AddedVariationIDs: []string{},
			},
		},
		{
			desc: "add missing variations to default strategy",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.DefaultStrategy = &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{Variation: "variation-A", Weight: 100000},
						},
					},
				}
				return f
			},
			expectedResult: VariationMigrationResult{
				Changed:           true,
				AddedToRules:      0,
				AddedToDefault:    2,
				AddedVariationIDs: []string{"variation-A", "variation-B", "variation-C"},
			},
		},
		{
			desc: "no changes for fixed strategies",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.DefaultStrategy = &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: "variation-A",
					},
				}
				return f
			},
			expectedResult: VariationMigrationResult{
				Changed:           false,
				AddedToRules:      0,
				AddedToDefault:    0,
				AddedVariationIDs: []string{},
			},
		},
		{
			desc: "preserve variation order when adding to strategies",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				f.DefaultStrategy = &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{Variation: "variation-B", Weight: 100000},
						},
					},
				}
				return f
			},
			expectedResult: VariationMigrationResult{
				Changed:           true,
				AddedToRules:      0,
				AddedToDefault:    2,
				AddedVariationIDs: []string{"variation-A", "variation-B", "variation-C"},
			},
		},
		{
			desc: "preserve variation order when adding to rules",
			setupFunc: func() *Feature {
				f := makeFeature("test-feature")
				// Clear default rules and add our test rule
				f.Rules = []*ftproto.Rule{
					{
						Id: "rule-1",
						Strategy: &ftproto.Strategy{
							Type: ftproto.Strategy_ROLLOUT,
							RolloutStrategy: &ftproto.RolloutStrategy{
								Variations: []*ftproto.RolloutStrategy_Variation{
									{Variation: "variation-C", Weight: 50000},
									{Variation: "variation-A", Weight: 50000},
								},
							},
						},
					},
				}
				return f
			},
			expectedResult: VariationMigrationResult{
				Changed:           true,
				AddedToRules:      1,
				AddedToDefault:    0,
				AddedVariationIDs: []string{"variation-A", "variation-B", "variation-C"},
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.setupFunc()
			result := actual.EnsureVariationsInStrategies()

			// Verify the result
			assert.Equal(t, p.expectedResult.Changed, result.Changed)
			assert.Equal(t, p.expectedResult.AddedToRules, result.AddedToRules)
			assert.Equal(t, p.expectedResult.AddedToDefault, result.AddedToDefault)

			// Check that the expected variation IDs are present (order may vary)
			assert.Equal(t, len(p.expectedResult.AddedVariationIDs), len(result.AddedVariationIDs))

			// Special test for order preservation
			if p.desc == "preserve variation order when adding to strategies" {
				// Verify that strategy variations match f.Variations order exactly
				// makeFeature creates variations in order: A, B, C
				// Strategy should be reconstructed to match this order: A, B, C
				expectedOrder := []string{"variation-A", "variation-B", "variation-C"}
				actualOrder := make([]string, len(actual.DefaultStrategy.RolloutStrategy.Variations))
				for i, v := range actual.DefaultStrategy.RolloutStrategy.Variations {
					actualOrder[i] = v.Variation
				}
				assert.Equal(t, expectedOrder, actualOrder, "Strategy variations should match f.Variations order exactly")

				// Verify weights are preserved for existing variations
				for _, v := range actual.DefaultStrategy.RolloutStrategy.Variations {
					if v.Variation == "variation-B" {
						assert.Equal(t, int32(100000), v.Weight, "Existing variation weight should be preserved")
					} else {
						assert.Equal(t, int32(0), v.Weight, "New variations should have weight 0")
					}
				}
			}

			// Special test for rule order preservation
			if p.desc == "preserve variation order when adding to rules" {
				// Verify that rule strategy variations match f.Variations order exactly
				// makeFeature creates variations in order: A, B, C
				// Rule strategy should be reconstructed to match this order: A, B, C
				expectedOrder := []string{"variation-A", "variation-B", "variation-C"}
				actualOrder := make([]string, len(actual.Rules[0].Strategy.RolloutStrategy.Variations))
				for i, v := range actual.Rules[0].Strategy.RolloutStrategy.Variations {
					actualOrder[i] = v.Variation
				}
				assert.Equal(t, expectedOrder, actualOrder, "Rule strategy variations should match f.Variations order exactly")

				// Verify weights are preserved for existing variations
				for _, v := range actual.Rules[0].Strategy.RolloutStrategy.Variations {
					if v.Variation == "variation-C" {
						assert.Equal(t, int32(50000), v.Weight, "Existing variation weight should be preserved")
					} else if v.Variation == "variation-A" {
						assert.Equal(t, int32(50000), v.Weight, "Existing variation weight should be preserved")
					} else {
						assert.Equal(t, int32(0), v.Weight, "New variations should have weight 0")
					}
				}
			}
		})
	}
}
