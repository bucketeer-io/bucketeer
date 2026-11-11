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

package notifier

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	domainproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription/sender"
)

func featureSnapshot(t *testing.T, feature *featureproto.Feature) string {
	t.Helper()
	data, err := json.Marshal(feature)
	require.NoError(t, err)
	return string(data)
}

func TestFeatureChangeSummary(t *testing.T) {
	t.Parallel()
	for _, lang := range []string{locale.En, locale.Ja} {
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			notification := &senderproto.DomainEventNotification{
				EntityId: "flag",
				PreviousEntityData: `{"id":"flag","enabled":false,"tags":["team-a"],
                    "variations":[{"id":"v1","value":"true","name":"On"}],
                    "targets":[{"variation":"v1","users":["old","unchanged"]}]}`,
				EntityData: `{"id":"flag","enabled":true,"tags":["team-a","payments"],
                    "variations":[{"id":"v1","value":"true","name":"On"}],
                    "targets":[{"variation":"v1","users":["9941","unchanged"]}]}`,
			}
			actual := featureChangeSummary(notification, newTestLocalizer(lang))
			assert.Contains(t, actual, "false → true")
			assert.Contains(t, actual, `["team-a"] → ["team-a","payments"]`)
			assert.Contains(t, actual, `"On (true)"`)
			assert.Contains(t, actual, `["9941"]`)
			assert.Contains(t, actual, `["old"]`)
			assert.NotContains(t, actual, "unchanged")
			assert.Len(t, strings.Split(actual, "\n"), 4)
			if lang == locale.Ja {
				assert.Contains(t, actual, "有効")
				assert.Contains(t, actual, "を追加")
				assert.Contains(t, actual, "を削除")
			} else {
				assert.Contains(t, actual, "Enabled")
				assert.Contains(t, actual, "added users")
				assert.Contains(t, actual, "removed users")
			}
		})
	}
}

func TestFeatureChangeFallback(t *testing.T) {
	t.Parallel()
	for _, data := range []string{"", "null", "{}", "[]", "invalid", `{"id":"other"}`} {
		t.Run(data, func(t *testing.T) {
			t.Parallel()
			for _, previous := range []bool{true, false} {
				n := &senderproto.DomainEventNotification{
					EntityId: "flag", EntityType: domainproto.Event_FEATURE,
					Type: domainproto.Event_FEATURE_UPDATED, Editor: &domainproto.Editor{},
					EntityData: `{"id":"flag"}`, PreviousEntityData: `{"id":"flag"}`,
				}
				if previous {
					n.PreviousEntityData = data
				} else {
					n.EntityData = data
				}
				assert.Empty(t, featureChangeSummary(n, newTestLocalizer(locale.En)))
				attachment, err := (&slackNotifier{webURL: "https://example.com"}).createDomainEventAttachment(
					n, newTestLocalizer(locale.En))
				require.NoError(t, err)
				assert.Contains(t, attachment.Text, "Feature flag has been updated")
				assert.Contains(t, attachment.Text, "Entity ID: flag")
			}
		})
	}
}

func TestFeatureChangeConfiguration(t *testing.T) {
	t.Parallel()
	n := &senderproto.DomainEventNotification{
		EntityId: "flag",
		PreviousEntityData: `{"id":"flag","version":1,"updated_at":1,
            "variations":[{"id":"v1","value":"old"}],
            "rules":[{"id":"rule1"}], "default_strategy":{"type":0,"fixed_strategy":{"variation":"v1"}}}`,
		EntityData: `{"id":"flag","version":2,"updated_at":2,
            "variations":[{"id":"v1","value":"new"},{"id":"v2","value":"other"}],
            "rules":[], "default_strategy":{"type":1,"rollout_strategy":{"variations":[{"variation":"v1","weight":100000}]}}}`,
	}
	summary := featureChangeSummary(n, newTestLocalizer(locale.En))
	assert.Contains(t, summary, "Default strategy / Strategy type:")
	assert.Contains(t, summary, "Rule:")
	assert.Contains(t, summary, "Variation (v1):")
	assert.Contains(t, summary, `"old"`)
	assert.Contains(t, summary, `"new"`)
	assert.Contains(t, summary, "Variation (v2): null →")
	assert.NotContains(t, summary, "updated_at")
	assert.NotContains(t, summary, "version")
	n.EntityData, n.PreviousEntityData = n.PreviousEntityData, n.EntityData
	assert.Contains(t, featureChangeSummary(n, newTestLocalizer(locale.En)), `"other","name":"","description":""} → null`)
}

func TestFeatureChangeIgnoresBookkeepingAndTargetOrder(t *testing.T) {
	t.Parallel()
	n := &senderproto.DomainEventNotification{
		EntityId: "flag",
		PreviousEntityData: `{"id":"flag","version":1,"updated_at":1,
            "variations":[{"id":"v1","value":"true"},{"id":"v2","value":"false"}],
            "targets":[{"variation":"v1","users":["a","b"]}]}`,
		EntityData: `{"id":"flag","version":2,"updated_at":2,
            "variations":[{"id":"v2","value":"false"},{"id":"v1","value":"true"}],
            "targets":[{"variation":"v1","users":["b","a","a"]}]}`,
	}
	assert.Empty(t, featureChangeSummary(n, newTestLocalizer(locale.En)))
}

func TestFeatureChangeLimitsAndEscaping(t *testing.T) {
	t.Parallel()
	before := &featureproto.Feature{Id: "flag"}
	after := &featureproto.Feature{Id: "flag", Description: "<!channel>&" + strings.Repeat("界", 1000)}
	for i := 0; i < 15; i++ {
		after.Variations = append(after.Variations, &featureproto.Variation{
			Id: string(rune('a' + i)), Value: strings.Repeat("界", 1000),
		})
	}
	n := &senderproto.DomainEventNotification{
		EntityId: "flag", PreviousEntityData: featureSnapshot(t, before), EntityData: featureSnapshot(t, after),
	}
	for _, lang := range []string{locale.En, locale.Ja} {
		summary := featureChangeSummary(n, newTestLocalizer(lang))
		assert.True(t, utf8.ValidString(summary))
		assert.Len(t, strings.Split(summary, "\n"), maxFeatureChanges+1)
		assert.Less(t, len(summary), 8000)
		assert.Contains(t, summary, "&lt;!channel&gt;&amp;")
		assert.NotContains(t, summary, "<!channel>")
		assert.Contains(t, summary, "…")
		assert.Contains(t, summary, newTestLocalizer(lang).MustLocalize(locale.NotificationFeatureChangesOmitted))
	}
}

func TestFeatureChangeAttachment(t *testing.T) {
	t.Parallel()
	for _, entity := range []domainproto.Event_EntityType{domainproto.Event_FEATURE, domainproto.Event_SEGMENT} {
		n := &senderproto.DomainEventNotification{
			EntityId: "flag", EntityType: entity, Type: domainproto.Event_FEATURE_UPDATED,
			Editor: &domainproto.Editor{Email: "editor@example.com"}, EnvironmentName: "Production", EnvironmentUrlCode: "prod",
			PreviousEntityData: `{"id":"flag","enabled":false}`, EntityData: `{"id":"flag","enabled":true}`,
		}
		attachment, err := (&slackNotifier{webURL: "https://example.com"}).createDomainEventAttachment(n, newTestLocalizer(locale.En))
		require.NoError(t, err)
		assert.Equal(t, "editor@example.com", attachment.AuthorName)
		assert.Contains(t, attachment.Text, "Environment: Production")
		assert.Contains(t, attachment.Text, "Entity ID: flag")
		if entity == domainproto.Event_FEATURE {
			assert.Contains(t, attachment.Text, "Enabled: false → true")
			assert.Contains(t, attachment.Text, "https://example.com/prod/features/flag")
		} else {
			assert.NotContains(t, attachment.Text, "Enabled:")
		}
	}
}

func TestFeatureChangeWithVariationSchema(t *testing.T) {
	t.Parallel()
	before := &featureproto.Feature{
		Id: "flag",
		VariationValueSchema: &featureproto.VariationValueSchema{
			Type: featureproto.VariationValueSchema_REGEX,
			Validator: &featureproto.VariationValueSchema_RegexValidator_{
				RegexValidator: &featureproto.VariationValueSchema_RegexValidator{Pattern: "^old$"},
			},
		},
	}
	// Use the actual domain-event JSON encoding, including the oneof wrapper.
	after := &featureproto.Feature{Id: "flag", Enabled: true, VariationValueSchema: before.VariationValueSchema}
	n := &senderproto.DomainEventNotification{
		EntityId: "flag", PreviousEntityData: featureSnapshot(t, before), EntityData: featureSnapshot(t, after),
	}
	summary := featureChangeSummary(n, newTestLocalizer(locale.En))
	assert.Equal(t, "• Enabled: false → true", summary)
	after.VariationValueSchema = &featureproto.VariationValueSchema{Type: featureproto.VariationValueSchema_ENUM}
	n.EntityData = featureSnapshot(t, after)
	assert.Contains(t, featureChangeSummary(n, newTestLocalizer(locale.En)), "Variation value schema:")
}

func TestFeatureChangeRuleClauseValuesWithRealisticIDs(t *testing.T) {
	t.Parallel()
	const (
		ruleID      = "d48293e5-b591-4e62-a72d-3f61ce371277"
		clauseID    = "959bcbed-9cfb-48b0-9134-44f743ad2172"
		variationID = "3749d3ec-6273-428c-8c41-1d5ea0998911"
	)
	before := &featureproto.Feature{
		Id: "flag",
		Rules: []*featureproto.Rule{{
			Id: ruleID,
			Strategy: &featureproto.Strategy{
				Type:          featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{Variation: variationID},
			},
			Clauses: []*featureproto.Clause{{
				Id: clauseID, Attribute: "country", Operator: featureproto.Clause_EQUALS, Values: []string{"US"},
			}},
		}},
	}
	var after featureproto.Feature
	require.NoError(t, json.Unmarshal([]byte(featureSnapshot(t, before)), &after))
	after.Rules[0].Clauses[0].Values = []string{"JP"}
	// The original whole-rule rendering loses the changed values behind UUIDs.
	assert.Equal(t, featureChangeValue(before.Rules), featureChangeValue(after.Rules))
	n := &senderproto.DomainEventNotification{
		EntityId: "flag", PreviousEntityData: featureSnapshot(t, before), EntityData: featureSnapshot(t, &after),
	}
	for _, lang := range []string{locale.En, locale.Ja} {
		summary := featureChangeSummary(n, newTestLocalizer(lang))
		assert.Contains(t, summary, ruleID)
		assert.Contains(t, summary, clauseID)
		assert.Contains(t, summary, `country: ["US"] → ["JP"]`)
		assert.NotContains(t, summary, variationID)
		assert.Len(t, strings.Split(summary, "\n"), 1)
	}
}

func TestFeatureChangeRulesAndClauses(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, before, after string
		contains            []string
		empty               bool
	}{
		{
			name:   "changed clause fields",
			before: `[{"id":"rule","clauses":[{"id":"clause","attribute":"country","operator":0,"values":["US"]}]}]`,
			after:  `[{"id":"rule","clauses":[{"id":"clause","attribute":"region","operator":1,"values":["JP"]}]}]`,
			contains: []string{`Rule attribute: "country" → "region"`, `Rule operator: "EQUALS" → "IN"`,
				`region: ["US"] → ["JP"]`},
		},
		{
			name:     "rule added",
			before:   `[]`,
			after:    `[{"id":"rule","clauses":[{"id":"clause","attribute":"country","values":["JP"]}]}]`,
			contains: []string{`Rule: null → "rule"`, `Rule condition: null → "clause"`, `country: null → ["JP"]`},
		},
		{
			name:     "clause removed",
			before:   `[{"id":"rule","clauses":[{"id":"clause","attribute":"country","values":["US"]}]}]`,
			after:    `[{"id":"rule"}]`,
			contains: []string{`Rule condition: "clause" → null`, `country: ["US"] → null`},
		},
		{
			name:     "rule order changed",
			before:   `[{"id":"first"},{"id":"second"}]`,
			after:    `[{"id":"second"},{"id":"first"}]`,
			contains: []string{`Rules order (first): 1 → 2`, `Rules order (second): 2 → 1`},
		},
		{
			name:   "clause and condition value order ignored",
			before: `[{"id":"rule","clauses":[{"id":"a","values":["US","JP"]},{"id":"b"}]}]`,
			after:  `[{"id":"rule","clauses":[{"id":"b"},{"id":"a","values":["JP","US"]}]}]`,
			empty:  true,
		},
		{
			name:     "shared condition values omitted",
			before:   `[{"id":"rule","clauses":[{"id":"clause","attribute":"country","values":["shared","US"]}]}]`,
			after:    `[{"id":"rule","clauses":[{"id":"clause","attribute":"country","values":["shared","JP"]}]}]`,
			contains: []string{`country: ["US"] → ["JP"]`},
		},
		{
			name:   "null entries ignored",
			before: `[null,{"id":"rule","clauses":[null]}]`,
			after:  `[{"id":"rule"}]`,
			empty:  true,
		},
		{
			name:     "rule and attribute escaped",
			before:   `[{"id":"<!channel>","clauses":[{"id":"<clause>","attribute":"<country>","values":["US"]}]}]`,
			after:    `[{"id":"<!channel>","clauses":[{"id":"<clause>","attribute":"<country>","values":["JP"]}]}]`,
			contains: []string{`Rule (&lt;!channel&gt;)`, `Rule condition (&lt;clause&gt;)`, `&lt;country&gt;: ["US"] → ["JP"]`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := &senderproto.DomainEventNotification{
				EntityId: "flag", PreviousEntityData: `{"id":"flag","rules":` + tt.before + `}`,
				EntityData: `{"id":"flag","rules":` + tt.after + `}`,
			}
			summary := featureChangeSummary(n, newTestLocalizer(locale.En))
			if tt.empty {
				assert.Empty(t, summary)
			}
			for _, expected := range tt.contains {
				assert.Contains(t, summary, expected)
			}
			assert.NotContains(t, summary, "shared")
		})
	}
}

func TestFeatureChangeRolloutAudienceWithRealisticIDs(t *testing.T) {
	t.Parallel()
	const (
		firstID  = "3749d3ec-6273-428c-8c41-1d5ea0998911"
		secondID = "d48293e5-b591-4e62-a72d-3f61ce371277"
	)
	beforeStrategy := &featureproto.Strategy{
		Type: featureproto.Strategy_ROLLOUT,
		RolloutStrategy: &featureproto.RolloutStrategy{
			Variations: []*featureproto.RolloutStrategy_Variation{
				{Variation: firstID, Weight: 50000}, {Variation: secondID, Weight: 50000},
			},
			Audience: &featureproto.Audience{Percentage: 10, DefaultVariation: firstID},
		},
	}
	var afterStrategy featureproto.Strategy
	data, err := json.Marshal(beforeStrategy)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &afterStrategy))
	afterStrategy.RolloutStrategy.Audience.Percentage = 20
	// Demonstrate that the original truncated strategy snippets hide the change.
	assert.Equal(t, featureChangeValue(beforeStrategy), featureChangeValue(&afterStrategy))
	for _, scope := range []string{"default", "rule"} {
		t.Run(scope, func(t *testing.T) {
			t.Parallel()
			before, after := &featureproto.Feature{Id: "flag"}, &featureproto.Feature{Id: "flag"}
			if scope == "default" {
				before.DefaultStrategy, after.DefaultStrategy = beforeStrategy, &afterStrategy
			} else {
				before.Rules = []*featureproto.Rule{{Id: "rule", Strategy: beforeStrategy}}
				after.Rules = []*featureproto.Rule{{Id: "rule", Strategy: &afterStrategy}}
			}
			n := &senderproto.DomainEventNotification{
				EntityId: "flag", PreviousEntityData: featureSnapshot(t, before), EntityData: featureSnapshot(t, after),
			}
			for _, lang := range []string{locale.En, locale.Ja} {
				l := newTestLocalizer(lang)
				summary := featureChangeSummary(n, l)
				assert.Contains(t, summary, l.MustLocalize(locale.NotificationFeatureAudiencePercentage)+": 10 → 20")
				assert.NotContains(t, summary, firstID)
				assert.NotContains(t, summary, secondID)
				assert.Len(t, strings.Split(summary, "\n"), 1)
			}
		})
	}
}

func TestFeatureChangeStrategyLeaves(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, before, after string
		want                []string
		empty               bool
	}{
		{
			name:   "weights",
			before: `{"type":1,"rollout_strategy":{"variations":[{"variation":"a","weight":10000},{"variation":"b","weight":90000}]}}`,
			after:  `{"type":1,"rollout_strategy":{"variations":[{"variation":"a","weight":20000},{"variation":"b","weight":80000}]}}`,
			want:   []string{"Rollout weight (a): 10000 → 20000", "Rollout weight (b): 90000 → 80000"},
		},
		{
			name:   "audience default",
			before: `{"type":1,"rollout_strategy":{"audience":{"percentage":10,"default_variation":"a"}}}`,
			after:  `{"type":1,"rollout_strategy":{"audience":{"percentage":10,"default_variation":"b"}}}`,
			want:   []string{`Audience default variation: "a" → "b"`},
		},
		{
			name:   "audience added at zero percent",
			before: `{"type":1,"rollout_strategy":{}}`,
			after:  `{"type":1,"rollout_strategy":{"audience":{"percentage":0}}}`,
			want:   []string{"Audience percentage: null → 0"},
		},
		{
			name:   "audience removed",
			before: `{"type":1,"rollout_strategy":{"audience":{"percentage":10,"default_variation":"a"}}}`,
			after:  `{"type":1,"rollout_strategy":{}}`,
			want:   []string{"Audience percentage: 10 → null", `Audience default variation: "a" → ""`},
		},
		{
			name:   "fixed variation",
			before: `{"type":0,"fixed_strategy":{"variation":"a"}}`,
			after:  `{"type":0,"fixed_strategy":{"variation":"b"}}`,
			want:   []string{`Variation: "a" → "b"`},
		},
		{
			name:   "strategy added",
			before: `null`, after: `{"type":0,"fixed_strategy":{"variation":"a"}}`,
			want: []string{`Strategy type: null → "FIXED"`, `Variation: "" → "a"`},
		},
		{
			name:   "strategy removed",
			before: `{"type":1,"rollout_strategy":{"variations":[{"variation":"a","weight":100000}]}}`, after: `null`,
			want: []string{`Strategy type: "ROLLOUT" → null`, `Rollout weight (a): 100000 → null`},
		},
		{
			name:   "rollout variation order ignored",
			before: `{"type":1,"rollout_strategy":{"variations":[{"variation":"a","weight":10000},{"variation":"b","weight":90000}]}}`,
			after:  `{"type":1,"rollout_strategy":{"variations":[{"variation":"b","weight":90000},{"variation":"a","weight":10000}]}}`,
			empty:  true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := &senderproto.DomainEventNotification{
				EntityId:           "flag",
				PreviousEntityData: `{"id":"flag","default_strategy":` + tt.before + `}`,
				EntityData:         `{"id":"flag","default_strategy":` + tt.after + `}`,
			}
			summary := featureChangeSummary(n, newTestLocalizer(locale.En))
			for _, want := range tt.want {
				assert.Contains(t, summary, want)
			}
			if tt.empty {
				assert.Empty(t, summary)
			}
		})
	}
}

func TestFeatureChangeRulePriorityWithAdditionsAndRemovals(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name          string
		before, after []string
		reordered     bool
	}{
		{"delete and reorder", []string{"A", "B", "C"}, []string{"B", "A"}, true},
		{"add and reorder", []string{"A", "B"}, []string{"B", "A", "C"}, true},
		{"replace and reorder", []string{"A", "B", "C"}, []string{"B", "A", "D"}, true},
		{"delete only", []string{"A", "C", "B"}, []string{"A", "B"}, false},
		{"insert only", []string{"A", "B"}, []string{"C", "A", "B"}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			before, after := &featureproto.Feature{Id: "flag"}, &featureproto.Feature{Id: "flag"}
			for _, id := range tt.before {
				before.Rules = append(before.Rules, &featureproto.Rule{Id: id})
			}
			for _, id := range tt.after {
				after.Rules = append(after.Rules, &featureproto.Rule{Id: id})
			}
			n := &senderproto.DomainEventNotification{
				EntityId: "flag", PreviousEntityData: featureSnapshot(t, before), EntityData: featureSnapshot(t, after),
			}
			summary := featureChangeSummary(n, newTestLocalizer(locale.En))
			if tt.reordered {
				assert.Contains(t, summary, "Rules order (A): 1 → 2")
				assert.Contains(t, summary, "Rules order (B): 2 → 1")
			} else {
				assert.NotContains(t, summary, "Rules order")
			}
		})
	}
}
