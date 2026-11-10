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
		{locale.DefaultStrategy, before.DefaultStrategy, after.DefaultStrategy},
		{locale.OffVariation,
			variationLabel(&before.Feature, before.OffVariation),
			variationLabel(&after.Feature, after.OffVariation)},
		{locale.Rule, before.Rules, after.Rules},
		{locale.Prerequisite, before.Prerequisites, after.Prerequisites},
		{locale.NotificationFeatureVariationType, before.VariationType.String(), after.VariationType.String()},
		{locale.NotificationFeatureVariationSchema, before.VariationValueSchema, after.VariationValueSchema},
	}
	for _, field := range fields {
		add(l.MustLocalize(field.label), field.before, field.after)
	}
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
