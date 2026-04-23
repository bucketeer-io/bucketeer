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
	"errors"
	"slices"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/proto/common"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// Update returns a new Feature with the updated values.
func (f *Feature) Update(
	name, description *wrapperspb.StringValue,
	tags *common.StringListValue,
	enabled *wrapperspb.BoolValue,
	archived *wrapperspb.BoolValue,
	defaultStrategy *feature.Strategy,
	offVariation *wrapperspb.StringValue,
	resetSamplingSeed bool,
	prerequisiteChanges []*feature.PrerequisiteChange,
	targetChanges []*feature.TargetChange,
	ruleChanges []*feature.RuleChange,
	variationChanges []*feature.VariationChange,
	tagChanges []*feature.TagChange,
	maintainer *wrapperspb.StringValue,
	ruleOrder []string,
) (*Feature, error) {
	// Use copier.CopyWithOption with DeepCopy: true to standardize empty slices as []
	// This ensures consistent JSON serialization in both API responses and audit logs
	var clonedFeature feature.Feature
	if err := copier.CopyWithOption(&clonedFeature, f.Feature, copier.Option{DeepCopy: true}); err != nil {
		return nil, err
	}

	updated := &Feature{
		Feature: &clonedFeature,
	}

	// We split variation changes into two separate steps to handle dependencies correctly:
	//
	// Step 1: Apply variation creations and updates first.
	// - This ensures newly created or updated variations exist when validating and applying other changes
	//   (e.g., rules, targets, offVariation, defaultStrategy).
	// - Without this step, validations referencing newly created variations would fail.
	//
	// Step 2: Apply variation deletions last.
	// - Deleting variations last ensures that any references to deleted variations (including references
	//   created or updated in the same request) are properly cleaned up.
	// - If deletions were processed earlier, we could end up with invalid references to non-existent variations.
	//
	// This two-step approach maintains data integrity and prevents invalid intermediate states.

	// Step 1: Apply variation creations and updates first
	var variationCreationsAndUpdates, variationDeletions []*feature.VariationChange
	for _, change := range variationChanges {
		if change.ChangeType == feature.ChangeType_DELETE {
			variationDeletions = append(variationDeletions, change)
		} else {
			variationCreationsAndUpdates = append(variationCreationsAndUpdates, change)
		}
	}

	if err := updated.validateVariationChanges(variationCreationsAndUpdates); err != nil {
		return nil, err
	}
	if err := updated.applyVariationChanges(variationCreationsAndUpdates); err != nil {
		return nil, err
	}

	// Step 2: Validate the rest against the updated variation set
	if err := updated.validateAllChanges(
		name,
		defaultStrategy,
		offVariation,
		prerequisiteChanges,
		targetChanges,
		ruleChanges,
		tagChanges,
		ruleOrder,
	); err != nil {
		return nil, err
	}

	if err := updated.applyGeneralUpdates(
		name,
		description,
		tags,
		enabled,
		archived,
		defaultStrategy,
		offVariation,
		resetSamplingSeed,
		maintainer,
	); err != nil {
		return nil, err
	}

	if err := updated.applyGranularCRUDChanges(
		prerequisiteChanges,
		targetChanges,
		ruleChanges,
		tagChanges,
	); err != nil {
		return nil, err
	}

	// Rule ordering is intentionally separate from CRUD.
	// It must run after creates/deletes/updates so it sees the final rule set.
	if len(ruleOrder) > 0 {
		if err := updated.applyRuleOrder(ruleOrder); err != nil {
			return nil, err
		}
	}

	// Step 3: variation deletions last
	if err := updated.validateVariationChanges(variationDeletions); err != nil {
		return nil, err
	}
	if err := updated.applyVariationChanges(variationDeletions); err != nil {
		return nil, err
	}

	// Increment version and update timestamp if there are changes
	if !proto.Equal(f.Feature, updated.Feature) {
		updated.Version++
		updated.UpdatedAt = time.Now().Unix()
	}

	return updated, nil
}

// validateAllChanges ensures all inputs are valid before applying any changes
func (f *Feature) validateAllChanges(
	name *wrapperspb.StringValue,
	defaultStrategy *feature.Strategy,
	offVariation *wrapperspb.StringValue,
	prerequisiteChanges []*feature.PrerequisiteChange,
	targetChanges []*feature.TargetChange,
	ruleChanges []*feature.RuleChange,
	tagChanges []*feature.TagChange,
	ruleOrder []string,
) error {
	if name != nil && name.Value == "" {
		return errNameEmpty
	}
	if defaultStrategy != nil {
		if err := validateStrategy(defaultStrategy, f.Variations); err != nil {
			return err
		}
	}
	if offVariation != nil {
		if err := validateOffVariation(offVariation.Value, f.Variations); err != nil {
			return err
		}
	}
	if err := f.validatePrerequisiteChanges(prerequisiteChanges); err != nil {
		return err
	}
	if err := f.validateTargetChanges(targetChanges); err != nil {
		return err
	}
	if err := f.validateRuleChanges(ruleChanges); err != nil {
		return err
	}
	if err := f.validateTagChanges(tagChanges); err != nil {
		return err
	}

	// Optional early validation of ruleOrder against current post-general-update state.
	// Since CRUD hasn't been applied yet, we validate ruleOrder structure only here.
	// Membership/length are still authoritatively checked in applyRuleOrder after CRUD.
	if err := f.validateRuleOrderShape(ruleOrder); err != nil {
		return err
	}

	return nil
}

// applyGeneralUpdates applies simple field updates
func (f *Feature) applyGeneralUpdates(
	name, description *wrapperspb.StringValue,
	tags *common.StringListValue,
	enabled *wrapperspb.BoolValue,
	archived *wrapperspb.BoolValue,
	defaultStrategy *feature.Strategy,
	offVariation *wrapperspb.StringValue,
	resetSamplingSeed bool,
	maintainer *wrapperspb.StringValue,
) error {
	if name != nil {
		f.Name = name.Value
	}
	// Optional field
	if description != nil {
		f.Description = description.Value
	}
	// Optional field
	if tags != nil {
		f.Tags = unique(tags.Values)
	}
	if enabled != nil {
		if enabled.Value {
			_ = f.updateEnable()
		} else {
			_ = f.updateDisable()
		}
	}
	if archived != nil {
		if archived.Value {
			_ = f.updateArchive()
		} else {
			_ = f.updateUnarchive()
		}
	}
	if defaultStrategy != nil {
		f.DefaultStrategy = defaultStrategy
	}
	if offVariation != nil {
		f.OffVariation = offVariation.Value
	}
	if resetSamplingSeed {
		if err := f.ResetSamplingSeed(); err != nil {
			return err
		}
	}
	if maintainer != nil {
		if maintainer.Value == "" {
			return errMaintainerCannotBeEmpty
		}
		f.Maintainer = maintainer.Value
	}
	return nil
}

// applyVariationChanges handles only variation creations, updates, and deletions.
func (f *Feature) applyVariationChanges(
	variationChanges []*feature.VariationChange,
) error {
	for _, change := range variationChanges {
		switch change.ChangeType {
		case feature.ChangeType_CREATE:
			if err := f.updateAddVariation(
				change.Variation.Id,
				change.Variation.Value,
				change.Variation.Name,
				change.Variation.Description,
			); err != nil {
				return err
			}
		case feature.ChangeType_UPDATE:
			if err := f.updateChangeVariation(change.Variation); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			if err := f.updateRemoveVariation(change.Variation.Id); err != nil {
				return err
			}
		}
	}
	return nil
}

// applyGranularCRUDChanges handles prerequisites, targets, rules, and tags CRUD.
// Rule ordering is intentionally not handled here.
func (f *Feature) applyGranularCRUDChanges(
	prerequisiteChanges []*feature.PrerequisiteChange,
	targetChanges []*feature.TargetChange,
	ruleChanges []*feature.RuleChange,
	tagChanges []*feature.TagChange,
) error {
	for _, change := range prerequisiteChanges {
		switch change.ChangeType {
		case feature.ChangeType_CREATE:
			if err := f.updateAddPrerequisite(
				change.Prerequisite.FeatureId,
				change.Prerequisite.VariationId,
			); err != nil {
				return err
			}
		case feature.ChangeType_UPDATE:
			if err := f.updateChangePrerequisiteVariation(
				change.Prerequisite.FeatureId,
				change.Prerequisite.VariationId,
			); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			if err := f.updateRemovePrerequisite(change.Prerequisite.FeatureId); err != nil {
				return err
			}
		}
	}

	for _, change := range targetChanges {
		switch change.ChangeType {
		case feature.ChangeType_CREATE, feature.ChangeType_UPDATE:
			if err := f.updateAddTargetUsers(change.Target); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			if err := f.updateRemoveTargetUsers(change.Target); err != nil {
				return err
			}
		}
	}

	for _, change := range ruleChanges {
		switch change.ChangeType {
		case feature.ChangeType_CREATE:
			if err := f.updateAddRule(change.Rule); err != nil {
				return err
			}
		case feature.ChangeType_UPDATE:
			if err := f.updateChangeRule(change.Rule); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			if err := f.updateRemoveRule(change.Rule.Id); err != nil {
				return err
			}
		}
	}

	for _, change := range tagChanges {
		switch change.ChangeType {
		case feature.ChangeType_CREATE, feature.ChangeType_UPDATE:
			if err := f.updateAddTag(change.Tag); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			if err := f.updateRemoveTag(change.Tag); err != nil {
				return err
			}
		}
	}

	return nil
}

// applyRuleOrder reorders f.Rules to match the given list of rule IDs.
// Must be called after all rule CREATE/UPDATE/DELETE changes have been applied.
func (f *Feature) applyRuleOrder(ruleIDs []string) error {
	return f.ChangeRulesOrder(ruleIDs)
}

func (f *Feature) updateEnable() error {
	if !f.Enabled {
		f.Enabled = true
	}
	return nil
}

func (f *Feature) updateDisable() error {
	if f.Enabled {
		f.Enabled = false
	}
	return nil
}

func (f *Feature) updateArchive() error {
	if !f.Archived {
		f.Archived = true
	}
	return nil
}

func (f *Feature) updateUnarchive() error {
	if f.Archived {
		f.Archived = false
	}
	return nil
}

func (f *Feature) updateAddVariation(id, value, name, description string) error {
	if id == "" {
		return errVariationIDRequired
	}
	if value == "" {
		return errVariationValueRequired
	}
	if name == "" {
		return errVariationNameRequired
	}
	if err := f.validateVariationValue(id, value); err != nil {
		return err
	}
	if _, err := f.findVariationIndex(id); err == nil {
		return errVariationValueUnique // variation already exists
	}
	f.Variations = append(f.Variations, &feature.Variation{
		Id:          id,
		Value:       value,
		Name:        name,
		Description: description,
	})
	f.addTarget(id)
	f.updateAddVariationToRules(id)
	f.updateAddVariationToDefaultStrategy(id)
	return nil
}

func (f *Feature) updateChangeVariation(variation *feature.Variation) error {
	if variation == nil {
		return errVariationRequired
	}
	idx, err := f.findVariationIndex(variation.Id)
	if err != nil {
		return err
	}
	if variation.Name == "" {
		return errVariationNameRequired
	}
	if err := f.validateVariationValue(variation.Id, variation.Value); err != nil {
		return err
	}

	// Only update if the variation actually changed
	if !proto.Equal(f.Variations[idx], variation) {
		f.Variations[idx] = variation
	}
	return nil
}

func (f *Feature) updateRemoveVariation(id string) error {
	idx, err := f.updateFindVariationIndex(id)
	if err != nil {
		return err
	}
	if err := f.updateValidateRemoveVariation(id); err != nil {
		return err
	}
	// Clean up references to this variation before removing it
	if err = f.updateRemoveTarget(id); err != nil {
		return err
	}
	f.updateRemoveVariationFromRules(id)
	f.updateRemoveVariationFromDefaultStrategy(id)
	f.Variations = slices.Delete(f.Variations, idx, idx+1)
	return nil
}

// updateFindVariationIndex finds the index of the variation with the specified ID
func (f *Feature) updateFindVariationIndex(id string) (int, error) {
	for i := range f.Variations {
		if f.Variations[i].Id == id {
			return i, nil
		}
	}
	return -1, errVariationNotFound
}

// updateFindTarget finds the index of the target with the specified variation ID
func (f *Feature) updateFindTarget(id string) (int, error) {
	for i := range f.Targets {
		if f.Targets[i].Variation == id {
			return i, nil
		}
	}
	return -1, errTargetNotFound
}

// updateValidateRemoveVariation validates that a variation can be safely removed
func (f *Feature) updateValidateRemoveVariation(id string) error {
	if len(f.Variations) <= 2 {
		return errVariationsMustHaveAtLeastTwoVariations
	}
	if f.OffVariation == id {
		return ErrVariationInUse
	}
	// Check if the individual targeting has any users
	idx, err := f.updateFindTarget(id)
	if err != nil {
		return err
	}
	if len(f.Targets[idx].Users) > 0 {
		return ErrVariationInUse
	}
	if strategyContainsVariation(id, f.DefaultStrategy) {
		return ErrVariationInUse
	}
	if f.updateRulesContainsVariation(id) {
		return ErrVariationInUse
	}
	return nil
}

// updateRulesContainsVariation checks if any rule contains the specified variation
func (f *Feature) updateRulesContainsVariation(id string) bool {
	for _, r := range f.Rules {
		if ok := strategyContainsVariation(id, r.Strategy); ok {
			return true
		}
	}
	return false
}

// updateRemoveTarget removes the target entry for the specified variation
func (f *Feature) updateRemoveTarget(variationID string) error {
	idx, err := f.updateFindTarget(variationID)
	if err != nil {
		return err
	}
	f.Targets = slices.Delete(f.Targets, idx, idx+1)
	return nil
}

// updateRemoveVariationFromRules removes the variation from all rollout strategies in rules
func (f *Feature) updateRemoveVariationFromRules(variationID string) {
	for _, rule := range f.Rules {
		if rule.Strategy.Type == feature.Strategy_ROLLOUT {
			f.updateRemoveVariationFromRolloutStrategy(rule.Strategy.RolloutStrategy, variationID)
		}
	}
}

// updateRemoveVariationFromDefaultStrategy removes the variation from the default strategy if it's a rollout
func (f *Feature) updateRemoveVariationFromDefaultStrategy(variationID string) {
	if f.DefaultStrategy != nil && f.DefaultStrategy.Type == feature.Strategy_ROLLOUT {
		f.updateRemoveVariationFromRolloutStrategy(f.DefaultStrategy.RolloutStrategy, variationID)
	}
}

// updateRemoveVariationFromRolloutStrategy removes all instances of the variation from a rollout strategy
func (f *Feature) updateRemoveVariationFromRolloutStrategy(strategy *feature.RolloutStrategy, variationID string) {
	// Remove all instances of the variation, regardless of weight
	filteredVariations := make([]*feature.RolloutStrategy_Variation, 0, len(strategy.Variations))
	for _, v := range strategy.Variations {
		if v.Variation != variationID {
			filteredVariations = append(filteredVariations, v)
		}
	}
	strategy.Variations = filteredVariations
}

func (f *Feature) updateAddPrerequisite(featureID, variationID string) error {
	if err := validatePrerequisite(featureID, variationID); err != nil {
		return err
	}
	if _, err := f.findPrerequisite(featureID); err == nil {
		return errPrerequisiteAlreadyExists
	}
	f.Prerequisites = append(f.Prerequisites, &feature.Prerequisite{
		FeatureId:   featureID,
		VariationId: variationID,
	})
	return nil
}

func (f *Feature) updateChangePrerequisiteVariation(featureID, variationID string) error {
	if err := validatePrerequisite(featureID, variationID); err != nil {
		return err
	}
	idx, err := f.findPrerequisiteIndex(featureID)
	if err != nil {
		return err
	}

	// Only update if the variation actually changed
	if f.Prerequisites[idx].VariationId != variationID {
		f.Prerequisites[idx].VariationId = variationID
	}
	return nil
}

func (f *Feature) updateRemovePrerequisite(featureID string) error {
	idx, err := f.findPrerequisiteIndex(featureID)
	if err != nil {
		return err
	}
	f.Prerequisites = slices.Delete(f.Prerequisites, idx, idx+1)
	return nil
}

func (f *Feature) updateAddTargetUsers(target *feature.Target) error {
	idx, err := f.findTarget(target.Variation)
	if err != nil {
		return err
	}
	if target.Users == nil {
		return errTargetUsersRequired
	}
	for _, user := range target.Users {
		if user == "" {
			return errTargetUserRequired
		}
		if !contains(user, f.Targets[idx].Users) {
			f.Targets[idx].Users = append(f.Targets[idx].Users, user)
		}
	}
	return nil
}

func (f *Feature) updateRemoveTargetUsers(target *feature.Target) error {
	idx, err := f.findTarget(target.Variation)
	if err != nil {
		return err
	}
	if target.Users == nil {
		return errTargetUsersRequired
	}
	for _, user := range target.Users {
		if user == "" {
			return errTargetUserRequired
		}
		uidx, err := index(user, f.Targets[idx].Users)
		if err != nil {
			// User not found, skip (don't return error for non-existent user removal)
			continue
		}
		f.Targets[idx].Users = append(f.Targets[idx].Users[:uidx], f.Targets[idx].Users[uidx+1:]...)
	}
	return nil
}

func (f *Feature) updateAddRule(rule *feature.Rule) error {
	if rule == nil {
		return errRuleRequired
	}
	if err := validateClauses(rule.Clauses); err != nil {
		return err
	}
	if err := validateStrategy(rule.Strategy, f.Variations); err != nil {
		return err
	}
	if _, err := f.findRule(rule.Id); err == nil {
		return errRuleAlreadyExists
	}
	f.Rules = append(f.Rules, rule)
	return nil
}

func (f *Feature) updateChangeRule(rule *feature.Rule) error {
	if rule == nil {
		return errRuleRequired
	}
	idx, err := f.findRuleIndex(rule.Id)
	if err != nil {
		return err
	}
	if err := validateClauses(rule.Clauses); err != nil {
		return err
	}
	if err := validateStrategy(rule.Strategy, f.Variations); err != nil {
		return err
	}

	// Only update if the rule actually changed
	existingRule := f.Rules[idx]
	if !proto.Equal(existingRule, rule) {
		f.Rules[idx] = rule
	}
	return nil
}

func (f *Feature) updateRemoveRule(id string) error {
	idx, err := f.findRuleIndex(id)
	if err != nil {
		return err
	}
	f.Rules = slices.Delete(f.Rules, idx, idx+1)
	return nil
}

func (f *Feature) updateAddTag(tag string) error {
	if slices.Contains(f.Tags, tag) {
		return nil
	}
	f.Tags = append(f.Tags, tag)
	return nil
}

func (f *Feature) updateRemoveTag(tag string) error {
	index := slices.Index(f.Tags, tag)
	if index == -1 {
		return errors.New("feature: tag not found")
	}
	f.Tags = slices.Delete(f.Tags, index, index+1)
	return nil
}

func (f *Feature) findRuleIndex(id string) (int, error) {
	for i, rule := range f.Rules {
		if rule.Id == id {
			return i, nil
		}
	}
	return -1, errRuleNotFound
}

func (f *Feature) findPrerequisiteIndex(featureID string) (int, error) {
	for i, prereq := range f.Prerequisites {
		if prereq.FeatureId == featureID {
			return i, nil
		}
	}
	return -1, errPrerequisiteNotFound
}

func (f *Feature) updateAddVariationToRules(variationID string) {
	for _, rule := range f.Rules {
		if rule.Strategy.Type == feature.Strategy_ROLLOUT {
			f.updateAddVariationToRolloutStrategy(rule.Strategy.RolloutStrategy, variationID)
		}
	}
}

func (f *Feature) updateAddVariationToDefaultStrategy(variationID string) {
	if f.DefaultStrategy != nil && f.DefaultStrategy.Type == feature.Strategy_ROLLOUT {
		f.updateAddVariationToRolloutStrategy(f.DefaultStrategy.RolloutStrategy, variationID)
	}
}

func (f *Feature) updateAddVariationToRolloutStrategy(strategy *feature.RolloutStrategy, variationID string) {
	strategy.Variations = append(strategy.Variations, &feature.RolloutStrategy_Variation{
		Variation: variationID,
		Weight:    0,
	})
}

func (f *Feature) validatePrerequisiteChanges(changes []*feature.PrerequisiteChange) error {
	for _, change := range changes {
		if change == nil {
			return errPrerequisiteRequired
		}
		if change.Prerequisite == nil {
			return errPrerequisiteRequired
		}

		switch change.ChangeType {
		case feature.ChangeType_CREATE, feature.ChangeType_UPDATE:
			if err := validatePrerequisite(
				change.Prerequisite.FeatureId,
				change.Prerequisite.VariationId,
			); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			if change.Prerequisite.FeatureId == "" {
				return errPrerequisiteFeatureIDRequired
			}
		default:
			return errUnknownChangeType
		}
	}
	return nil
}

func (f *Feature) validateTargetChanges(changes []*feature.TargetChange) error {
	for _, change := range changes {
		if change == nil {
			return errTargetRequired
		}
		if change.Target == nil {
			return errTargetRequired
		}

		switch change.ChangeType {
		case feature.ChangeType_CREATE, feature.ChangeType_UPDATE:
			if err := validateTargets([]*feature.Target{change.Target}, f.Variations); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			// delete only needs enough info to identify the target bucket and users to remove
			if change.Target.Variation == "" {
				return errTargetVariationRequired
			}
			if change.Target.Users == nil {
				return errTargetUsersRequired
			}
			for _, user := range change.Target.Users {
				if user == "" {
					return errTargetUserRequired
				}
			}
		default:
			return errUnknownChangeType
		}
	}
	return nil
}

func (f *Feature) validateRuleChanges(changes []*feature.RuleChange) error {
	for _, change := range changes {
		if change == nil {
			return errRuleRequired
		}
		if change.Rule == nil {
			return errRuleRequired
		}

		switch change.ChangeType {
		case feature.ChangeType_CREATE, feature.ChangeType_UPDATE:
			if err := validateRules([]*feature.Rule{change.Rule}, f.Variations); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			if change.Rule.Id == "" {
				return errRuleIDRequired
			}
		default:
			return errUnknownChangeType
		}
	}
	return nil
}

func (f *Feature) validateTagChanges(changes []*feature.TagChange) error {
	for _, change := range changes {
		if change == nil {
			return errTagRequired
		}
		switch change.ChangeType {
		case feature.ChangeType_CREATE, feature.ChangeType_UPDATE, feature.ChangeType_DELETE:
			if change.Tag == "" {
				return errTagRequired
			}
		default:
			return errUnknownChangeType
		}
	}
	return nil
}

func (f *Feature) validateRuleOrderShape(ruleOrder []string) error {
	if len(ruleOrder) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(ruleOrder))
	for _, id := range ruleOrder {
		if id == "" {
			return errRuleIDRequired
		}
		if _, ok := seen[id]; ok {
			return errRulesOrderDuplicateIDs
		}
		seen[id] = struct{}{}
	}
	return nil
}
