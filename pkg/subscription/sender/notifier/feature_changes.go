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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription/sender"
)

const (
	maxFeatureChanges          = 10
	maxFeatureChangeValueRunes = 180
)

// Domain events use encoding/json for snapshots. Keep the schema as JSON
// because its protobuf oneof cannot be decoded by encoding/json into an interface.
type featureChangeSnapshot struct {
	featureproto.Feature
	VariationValueSchema interface{} `json:"variation_value_schema"`
}

// Snapshots are optional on older notifications. Invalid or mismatched snapshots
// must never prevent delivery of the original notification.
func featureChangeSummary(n *senderproto.DomainEventNotification, l locale.Localizer) string {
	var before, after featureChangeSnapshot
	if json.Unmarshal([]byte(n.PreviousEntityData), &before) != nil ||
		json.Unmarshal([]byte(n.EntityData), &after) != nil ||
		before.Id == "" || before.Id != after.Id || after.Id != n.EntityId {
		return ""
	}
	lines := make([]string, 0, maxFeatureChanges+1)
	add := func(label string, oldValue, newValue interface{}) {
		if reflect.DeepEqual(oldValue, newValue) {
			return
		}
		lines = append(lines, fmt.Sprintf("• %s: %s → %s", label,
			featureChangeValue(oldValue), featureChangeValue(newValue)))
	}
	fields := []struct {
		label         string
		before, after interface{}
	}{
		{locale.NotificationFeatureEnabled, before.Enabled, after.Enabled},
		{locale.NotificationFeatureArchived, before.Archived, after.Archived},
		{locale.NotificationFeatureDeleted, before.Deleted, after.Deleted},
		{locale.NotificationFeatureName, before.Name, after.Name},
		{locale.NotificationFeatureDescription, before.Description, after.Description},
		{locale.Tag, before.Tags, after.Tags},
		{locale.OffVariation,
			variationLabel(&before.Feature, before.OffVariation),
			variationLabel(&after.Feature, after.OffVariation)},
		{locale.Prerequisite, before.Prerequisites, after.Prerequisites},
		{locale.NotificationFeatureVariationType, before.VariationType.String(), after.VariationType.String()},
		{locale.NotificationFeatureVariationSchema, before.VariationValueSchema, after.VariationValueSchema},
	}
	for _, field := range fields {
		add(l.MustLocalize(field.label), field.before, field.after)
	}
	addFeatureStrategyChanges(l.MustLocalize(locale.DefaultStrategy),
		before.DefaultStrategy, after.DefaultStrategy, l, add)
	addFeatureRuleChanges(before.Rules, after.Rules, l, add)
	// Match variations by stable ID so adding one does not report unrelated
	// variations as changed just because their array positions moved.
	oldVariations, newVariations := variationMap(&before.Feature), variationMap(&after.Feature)
	for _, id := range sortedKeys(oldVariations, newVariations) {
		label := l.MustLocalize(locale.Variation) + " (" + slackEscape(limitChangeValue(id)) + ")"
		add(label, oldVariations[id], newVariations[id])
	}
	oldTargets, newTargets := targetMap(&before.Feature), targetMap(&after.Feature)
	for _, id := range sortedKeys(oldTargets, newTargets) {
		label := variationLabel(&after.Feature, id)
		if newVariations[id] == nil {
			label = variationLabel(&before.Feature, id)
		}
		for _, change := range []struct {
			template string
			users    []string
		}{
			{locale.NotificationFeatureTargetAdded, addedUsers(newTargets[id], oldTargets[id])},
			{locale.NotificationFeatureTargetRemoved, addedUsers(oldTargets[id], newTargets[id])},
		} {
			if len(change.users) > 0 {
				lines = append(lines, "• "+l.MustLocalizeWithTemplate(change.template,
					featureChangeValue(label), featureChangeValue(change.users)))
			}
		}
	}
	if len(lines) > maxFeatureChanges {
		lines = append(lines[:maxFeatureChanges], l.MustLocalize(locale.NotificationFeatureChangesOmitted))
	}
	return strings.Join(lines, "\n")
}

// Render strategy leaves separately so unchanged variation IDs cannot hide
// rollout weights or audience settings behind the display limit.
func addFeatureStrategyChanges(
	label string,
	before, after *featureproto.Strategy,
	l locale.Localizer,
	add func(string, interface{}, interface{}),
) {
	var oldType, newType interface{}
	if before != nil {
		oldType = before.Type.String()
	}
	if after != nil {
		newType = after.Type.String()
	}
	add(label+" / "+l.MustLocalize(locale.NotificationFeatureStrategyType), oldType, newType)
	add(label+" / "+l.MustLocalize(locale.Variation),
		before.GetFixedStrategy().GetVariation(), after.GetFixedStrategy().GetVariation())
	oldRollout, newRollout := before.GetRolloutStrategy(), after.GetRolloutStrategy()
	oldWeights, newWeights := featureRolloutWeights(oldRollout), featureRolloutWeights(newRollout)
	for _, id := range sortedKeys(oldWeights, newWeights) {
		var oldWeight, newWeight interface{}
		if weight, ok := oldWeights[id]; ok {
			oldWeight = weight
		}
		if weight, ok := newWeights[id]; ok {
			newWeight = weight
		}
		weightLabel := label + " / " + l.MustLocalize(locale.NotificationFeatureRolloutWeight) +
			" (" + slackEscape(limitChangeValue(id)) + ")"
		add(weightLabel, oldWeight, newWeight)
	}
	oldAudience, newAudience := oldRollout.GetAudience(), newRollout.GetAudience()
	var oldPercentage, newPercentage interface{}
	if oldAudience != nil {
		oldPercentage = oldAudience.Percentage
	}
	if newAudience != nil {
		newPercentage = newAudience.Percentage
	}
	add(label+" / "+l.MustLocalize(locale.NotificationFeatureAudiencePercentage), oldPercentage, newPercentage)
	add(label+" / "+l.MustLocalize(locale.NotificationFeatureAudienceDefaultVariation),
		oldAudience.GetDefaultVariation(), newAudience.GetDefaultVariation())
}

func featureRolloutWeights(rollout *featureproto.RolloutStrategy) map[string]int32 {
	weights := make(map[string]int32)
	for _, variation := range rollout.GetVariations() {
		if variation != nil {
			weights[variation.Variation] = variation.Weight
		}
	}
	return weights
}

// Compare rules and clauses before shortening values. Serializing a complete
// rule can otherwise use the entire display budget on unchanged IDs/strategy.
func addFeatureRuleChanges(
	before, after []*featureproto.Rule,
	l locale.Localizer,
	add func(string, interface{}, interface{}),
) {
	oldRules, newRules := featureItemsByID(before), featureItemsByID(after)
	for _, id := range sortedKeys(oldRules, newRules) {
		oldRule, newRule := oldRules[id], newRules[id]
		label := l.MustLocalize(locale.Rule) + " (" + slackEscape(limitChangeValue(id)) + ")"
		if oldRule == nil || newRule == nil {
			// Keep the existence change small; configuration details follow in
			// their own entries instead of one truncated object.
			add(l.MustLocalize(locale.Rule), featureItemID(oldRule), featureItemID(newRule))
		}
		addFeatureStrategyChanges(label+" / "+l.MustLocalize(locale.RuleStrategy),
			oldRule.GetStrategy(), newRule.GetStrategy(), l, add)
		oldClauses := featureItemsByID(oldRule.GetClauses())
		newClauses := featureItemsByID(newRule.GetClauses())
		for _, clauseID := range sortedKeys(oldClauses, newClauses) {
			oldClause, newClause := oldClauses[clauseID], newClauses[clauseID]
			clauseLabel := label + " / " + l.MustLocalize(locale.RuleCondition) +
				" (" + slackEscape(limitChangeValue(clauseID)) + ")"
			if oldClause == nil || newClause == nil {
				add(label+" / "+l.MustLocalize(locale.RuleCondition),
					featureItemID(oldClause), featureItemID(newClause))
			}
			add(clauseLabel+" / "+l.MustLocalize(locale.RuleAttribute),
				oldClause.GetAttribute(), newClause.GetAttribute())
			add(clauseLabel+" / "+l.MustLocalize(locale.RuleOperator),
				oldClause.GetOperator().String(), newClause.GetOperator().String())
			// Values are a set. Omit shared values so a late change in a large
			// condition remains visible within the per-value display limit.
			removed := addedUsers(featureValueSet(oldClause.GetValues()), featureValueSet(newClause.GetValues()))
			added := addedUsers(featureValueSet(newClause.GetValues()), featureValueSet(oldClause.GetValues()))
			attribute := newClause.GetAttribute()
			if newClause == nil {
				attribute = oldClause.GetAttribute()
			}
			add(clauseLabel+" / "+slackEscape(limitChangeValue(attribute)), removed, added)
		}
	}
	// Compare relative priority among surviving rules. Insertions/deletions
	// alone do not change their relative order, but may accompany a real swap.
	oldPositions := featureItemPositions(before, func(id string) bool { return newRules[id] != nil })
	newPositions := featureItemPositions(after, func(id string) bool { return oldRules[id] != nil })
	for _, id := range sortedKeys(oldPositions, newPositions) {
		label := l.MustLocalize(locale.RulesOrder) + " (" + slackEscape(limitChangeValue(id)) + ")"
		add(label, oldPositions[id], newPositions[id])
	}
}

func featureItemPositions[T interface{ GetId() string }](items []T, include func(string) bool) map[string]int {
	result := make(map[string]int, len(items))
	for _, item := range items {
		if id := item.GetId(); id != "" && include(id) {
			result[id] = len(result) + 1
		}
	}
	return result
}

func featureItemsByID[T interface{ GetId() string }](items []T) map[string]T {
	result := make(map[string]T, len(items))
	for _, item := range items {
		if id := item.GetId(); id != "" {
			result[id] = item
		}
	}
	return result
}

func featureItemID(item interface{ GetId() string }) interface{} {
	if id := item.GetId(); id != "" {
		return id
	}
	return nil
}

func featureValueSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}

func featureChangeValue(value interface{}) string {
	// Protobuf entities only contain JSON-compatible values.
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return "…"
	}
	return slackEscape(limitChangeValue(strings.TrimSuffix(data.String(), "\n")))
}

func limitChangeValue(value string) string {
	runes := []rune(value)
	if len(runes) > maxFeatureChangeValueRunes {
		return string(runes[:maxFeatureChangeValueRunes]) + "…"
	}
	return value
}

func slackEscape(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(value)
}

func variationMap(feature *featureproto.Feature) map[string]*featureproto.Variation {
	result := make(map[string]*featureproto.Variation)
	for _, variation := range feature.Variations {
		if variation != nil {
			result[variation.Id] = variation
		}
	}
	return result
}

func variationLabel(feature *featureproto.Feature, id string) string {
	for _, variation := range feature.Variations {
		if variation != nil && variation.Id == id {
			if variation.Name != "" {
				return variation.Name + " (" + variation.Value + ")"
			}
			return variation.Value
		}
	}
	return id
}

func targetMap(feature *featureproto.Feature) map[string]map[string]bool {
	result := make(map[string]map[string]bool)
	for _, target := range feature.Targets {
		if target == nil {
			continue
		}
		if result[target.Variation] == nil {
			result[target.Variation] = make(map[string]bool)
		}
		for _, user := range target.Users {
			result[target.Variation][user] = true
		}
	}
	return result
}

func addedUsers(users, previous map[string]bool) []string {
	var result []string
	for user := range users {
		if !previous[user] {
			result = append(result, user)
		}
	}
	sort.Strings(result)
	return result
}

func sortedKeys[V any](before, after map[string]V) []string {
	keys := make(map[string]bool, len(before)+len(after))
	for key := range before {
		keys[key] = true
	}
	for key := range after {
		keys[key] = true
	}
	result := make([]string, 0, len(keys))
	for key := range keys {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}
