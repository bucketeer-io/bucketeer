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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestBuildFeatureContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc     string
		feature  *featureproto.Feature
		contains []string
		excludes []string
	}{
		{
			desc: "full feature with all fields",
			feature: &featureproto.Feature{
				Id:            "flag-123",
				Name:          "Dark Mode",
				Description:   "Enable dark mode for users",
				Enabled:       true,
				VariationType: featureproto.Feature_BOOLEAN,
				Variations: []*featureproto.Variation{
					{
						Id:          "var-on",
						Value:       "true",
						Name:        "ON",
						Description: "Dark mode enabled",
					},
					{
						Id:          "var-off",
						Value:       "false",
						Name:        "OFF",
						Description: "Dark mode disabled",
					},
				},
				Tags: []string{"ui", "frontend", "theme"},
				Rules: []*featureproto.Rule{
					{
						Id: "rule-1",
						Strategy: &featureproto.Strategy{
							Type: featureproto.Strategy_FIXED,
						},
						Clauses: []*featureproto.Clause{
							{
								Id:        "clause-1",
								Attribute: "email",
								Operator:  featureproto.Clause_ENDS_WITH,
								Values:    []string{"@example.com"},
							},
						},
					},
				},
				DefaultStrategy: &featureproto.Strategy{
					Type: featureproto.Strategy_ROLLOUT,
				},
				Prerequisites: []*featureproto.Prerequisite{
					{
						FeatureId:   "base-flag",
						VariationId: "var-1",
					},
				},
				Targets: []*featureproto.Target{
					{
						Variation: "var-on",
						Users:     []string{"user-1", "user-2"},
					},
				},
				OffVariation: "var-off",
			},
			contains: []string{
				`Name: "Dark Mode"`,
				`Description: "Enable dark mode for users"`,
				"Enabled: true",
				"Variation Type: BOOLEAN",
				"ON",
				"Dark mode enabled",
				"OFF",
				"Dark mode disabled",
				"Tags: ui, frontend, theme",
				"Targeting Rules: 1 rule(s)",
				"Strategy: FIXED",
				"Conditions: 1",
				"Operator: ENDS_WITH",
				"Default Strategy: ROLLOUT",
				`Depends on flag: "base-flag"`,
			},
			excludes: []string{
				// Privacy: variation values must NOT appear
				"\"true\"",
				"\"false\"",
				// Privacy: attribute values must NOT appear
				"@example.com",
				// Privacy: attribute names must NOT appear
				"email",
				// Privacy: user IDs must NOT appear
				"user-1",
				"user-2",
				// Privacy: clause values must NOT appear
				"example.com",
			},
		},
		{
			desc: "minimal feature with no optional fields",
			feature: &featureproto.Feature{
				Id:            "flag-minimal",
				Name:          "Simple Flag",
				Enabled:       false,
				VariationType: featureproto.Feature_STRING,
			},
			contains: []string{
				`Name: "Simple Flag"`,
				"Enabled: false",
				"Variation Type: STRING",
			},
			excludes: []string{
				"Description:",
				"Tags:",
				"Targeting Rules:",
				"Default Strategy:",
				"Prerequisites:",
			},
		},
		{
			desc: "variation with empty name falls back to ID",
			feature: &featureproto.Feature{
				Id:            "flag-no-name",
				Name:          "No Name Vars",
				VariationType: featureproto.Feature_NUMBER,
				Variations: []*featureproto.Variation{
					{
						Id:    "var-abc",
						Value: "42",
						Name:  "",
					},
				},
			},
			contains: []string{
				"var-abc",
			},
			excludes: []string{
				"42",
			},
		},
		{
			desc: "multiple rules with rollout strategy",
			feature: &featureproto.Feature{
				Id:            "flag-multi-rule",
				Name:          "Multi Rule Flag",
				VariationType: featureproto.Feature_BOOLEAN,
				Rules: []*featureproto.Rule{
					{
						Id: "rule-1",
						Strategy: &featureproto.Strategy{
							Type: featureproto.Strategy_FIXED,
						},
						Clauses: []*featureproto.Clause{
							{
								Operator: featureproto.Clause_EQUALS,
								Values:   []string{"secret-value"},
							},
							{
								Operator: featureproto.Clause_IN,
								Values:   []string{"val-a", "val-b"},
							},
						},
					},
					{
						Id: "rule-2",
						Strategy: &featureproto.Strategy{
							Type: featureproto.Strategy_ROLLOUT,
						},
						Clauses: []*featureproto.Clause{
							{
								Operator: featureproto.Clause_SEGMENT,
								Values:   []string{"segment-id"},
							},
						},
					},
				},
			},
			contains: []string{
				"Targeting Rules: 2 rule(s)",
				"Rule 1:",
				"Rule 2:",
				"Strategy: FIXED",
				"Strategy: ROLLOUT",
				"Conditions: 2",
				"Conditions: 1",
				"Operator: EQUALS",
				"Operator: IN",
				"Operator: SEGMENT",
			},
			excludes: []string{
				"secret-value",
				"val-a",
				"val-b",
				"segment-id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			result := buildFeatureContext(tt.feature)
			for _, want := range tt.contains {
				assert.Contains(t, result, want, "expected to contain: %s", want)
			}
			for _, notWant := range tt.excludes {
				assert.NotContains(t, result, notWant, "must NOT contain (privacy): %s", notWant)
			}
		})
	}
}

func TestBuildFeatureContext_VariationValueNeverLeaks(t *testing.T) {
	t.Parallel()
	// Exhaustive check: variation Value field must never appear in output
	secretValue := "super-secret-config-json-{\"key\":\"val\"}"
	f := &featureproto.Feature{
		Name:          "Value Leak Test",
		VariationType: featureproto.Feature_JSON,
		Variations: []*featureproto.Variation{
			{
				Id:    "v1",
				Value: secretValue,
				Name:  "Config A",
			},
		},
	}
	result := buildFeatureContext(f)
	assert.NotContains(t, result, secretValue)
	assert.NotContains(t, result, "super-secret")
	assert.Contains(t, result, "Config A")
}

func TestBuildFeatureContext_ClauseValuesNeverLeak(t *testing.T) {
	t.Parallel()
	f := &featureproto.Feature{
		Name:          "Clause Leak Test",
		VariationType: featureproto.Feature_BOOLEAN,
		Rules: []*featureproto.Rule{
			{
				Strategy: &featureproto.Strategy{Type: featureproto.Strategy_FIXED},
				Clauses: []*featureproto.Clause{
					{
						Attribute: "user.plan",
						Operator:  featureproto.Clause_EQUALS,
						Values:    []string{"enterprise", "pro"},
					},
				},
			},
		},
	}
	result := buildFeatureContext(f)
	assert.NotContains(t, result, "enterprise")
	assert.NotContains(t, result, "pro")
	assert.NotContains(t, result, "user.plan")
	assert.Contains(t, result, "EQUALS")
}

func TestBuildFeatureContext_NilFeature(t *testing.T) {
	t.Parallel()
	result := buildFeatureContext(nil)
	assert.Empty(t, result)
}

func TestBuildFeatureContext_TruncatesLongOutput(t *testing.T) {
	t.Parallel()
	// Create a feature with many rules to exceed maxFeatureContextLength
	rules := make([]*featureproto.Rule, 100)
	for i := range rules {
		clauses := make([]*featureproto.Clause, 5)
		for j := range clauses {
			clauses[j] = &featureproto.Clause{
				Operator: featureproto.Clause_EQUALS,
			}
		}
		rules[i] = &featureproto.Rule{
			Id: fmt.Sprintf("rule-%d", i),
			Strategy: &featureproto.Strategy{
				Type: featureproto.Strategy_ROLLOUT,
			},
			Clauses: clauses,
		}
	}
	f := &featureproto.Feature{
		Name:          "Many Rules Flag",
		VariationType: featureproto.Feature_BOOLEAN,
		Rules:         rules,
	}
	result := buildFeatureContext(f)
	assert.Contains(t, result, "... (truncated)")
	// Rune count of the truncated portion must not exceed the limit
	runes := []rune(result)
	truncSuffix := []rune("\n... (truncated)\n")
	assert.LessOrEqual(t, len(runes), maxFeatureContextLength+len(truncSuffix))
}

func TestBuildFeatureContext_OutputIsNonEmpty(t *testing.T) {
	t.Parallel()
	f := &featureproto.Feature{
		Name: "Test",
	}
	result := buildFeatureContext(f)
	assert.True(t, len(strings.TrimSpace(result)) > 0)
}
