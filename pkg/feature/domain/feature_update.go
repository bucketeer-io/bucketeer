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
	"slices"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/proto/common"
	"github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrTagNotFound = errors.New("feature: tag not found")
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

	// Step 2: Validate all other changes now that new variations exist
	if err := updated.validateAllChanges(
		name,
		defaultStrategy,
		offVariation,
		prerequisiteChanges,
		targetChanges,
		ruleChanges,
	); err != nil {
		return nil, err
	}

	// Apply general updates
	if err := updated.applyGeneralUpdates(
		name,
		description,
		tags,
		enabled,
		archived,
		defaultStrategy,
		offVariation,
		resetSamplingSeed,
	); err != nil {
		return nil, err
	}

	// Apply remaining granular updates (prerequisites, targets, rules, tags)
	if err := updated.applyGranularChanges(
		prerequisiteChanges,
		targetChanges,
		ruleChanges,
		tagChanges,
	); err != nil {
		return nil, err
	}

	// Step 3: Apply variation deletions last
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
) error {
	// Validate name if provided
	if name != nil && name.Value == "" {
		return errNameEmpty
	}

	// Validate default strategy if provided
	if defaultStrategy != nil {
		if err := validateStrategy(defaultStrategy, f.Variations); err != nil {
			return err
		}
	}

	// Validate off variation if provided
	if offVariation != nil {
		if err := validateOffVariation(offVariation.Value, f.Variations); err != nil {
			return err
		}
	}

	// Validate granular changes
	for _, change := range prerequisiteChanges {
		if err := validatePrerequisite(change.Prerequisite.FeatureId, change.Prerequisite.VariationId); err != nil {
			return err
		}
	}

	for _, change := range targetChanges {
		if err := validateTargets([]*feature.Target{change.Target}, f.Variations); err != nil {
			return err
		}
	}

	for _, change := range ruleChanges {
		if err := validateRules([]*feature.Rule{change.Rule}, f.Variations); err != nil {
			return err
		}
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

// applyGranularChanges handles prerequisites, targets, rules, and tags updates.
func (f *Feature) applyGranularChanges(
	prerequisiteChanges []*feature.PrerequisiteChange,
	targetChanges []*feature.TargetChange,
	ruleChanges []*feature.RuleChange,
	tagChanges []*feature.TagChange,
) error {
	// Prerequisites
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
			if err := f.updateRemovePrerequisite(
				change.Prerequisite.FeatureId,
			); err != nil {
				return err
			}
		}
	}

	// Individual Targets
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

	// Custom Rules
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

	// Tags
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

// New functions for Update method that handle timestamp updates intelligently
// These functions only update the timestamp when actual changes occur
// TODO: Remove these duplicate functions once old console is deprecated

// updateEnable enables the feature only if it's not already enabled
func (f *Feature) updateEnable() error {
	if !f.Enabled {
		f.Enabled = true
	}
	return nil
}

// updateDisable disables the feature only if it's not already disabled
func (f *Feature) updateDisable() error {
	if f.Enabled {
		f.Enabled = false
	}
	return nil
}

// updateArchive archives the feature only if it's not already archived
func (f *Feature) updateArchive() error {
	if !f.Archived {
		f.Archived = true
	}
	return nil
}

// updateUnarchive unarchives the feature only if it's not already unarchived
func (f *Feature) updateUnarchive() error {
	if f.Archived {
		f.Archived = false
	}
	return nil
}

// updateAddVariation adds a variation, updating timestamp only if successful
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
	return nil
}

// updateChangeVariation changes a variation only if it actually differs
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

// updateRemoveVariation removes a variation, updating timestamp only if successful
func (f *Feature) updateRemoveVariation(id string) error {
	if len(f.Variations) == 1 {
		return errVariationInUse
	}
	idx, err := f.findVariationIndex(id)
	if err != nil {
		return err
	}
	if err := f.validateRemoveVariation(id); err != nil {
		return err
	}
	f.Variations = slices.Delete(f.Variations, idx, idx+1)
	return nil
}

// updateAddPrerequisite adds a prerequisite, updating timestamp only if successful
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

// updateChangePrerequisiteVariation changes a prerequisite variation only if it differs
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

// updateRemovePrerequisite removes a prerequisite, updating timestamp only if successful
func (f *Feature) updateRemovePrerequisite(featureID string) error {
	idx, err := f.findPrerequisiteIndex(featureID)
	if err != nil {
		return err
	}
	f.Prerequisites = slices.Delete(f.Prerequisites, idx, idx+1)
	return nil
}

// updateAddTargetUsers adds target users, updating timestamp only if users were actually added
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

// updateRemoveTargetUsers removes target users, updating timestamp only if users were actually removed
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

// updateAddRule adds a rule, updating timestamp only if successful
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

// updateChangeRule changes a rule only if it actually differs
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

// updateRemoveRule removes a rule, updating timestamp only if successful
func (f *Feature) updateRemoveRule(id string) error {
	idx, err := f.findRuleIndex(id)
	if err != nil {
		return err
	}
	f.Rules = slices.Delete(f.Rules, idx, idx+1)
	return nil
}

// updateAddTag adds a tag, updating timestamp only if tag was actually added
func (f *Feature) updateAddTag(tag string) error {
	if slices.Contains(f.Tags, tag) {
		return nil
	}
	f.Tags = append(f.Tags, tag)
	return nil
}

// updateRemoveTag removes a tag, updating timestamp only if tag was actually removed
func (f *Feature) updateRemoveTag(tag string) error {
	index := slices.Index(f.Tags, tag)
	if index == -1 {
		return ErrTagNotFound
	}
	f.Tags = slices.Delete(f.Tags, index, index+1)
	return nil
}

// Helper functions to support the update methods

// findRuleIndex finds the index of a rule by ID
func (f *Feature) findRuleIndex(id string) (int, error) {
	for i, rule := range f.Rules {
		if rule.Id == id {
			return i, nil
		}
	}
	return -1, errRuleNotFound
}

// findPrerequisiteIndex finds the index of a prerequisite by feature ID
func (f *Feature) findPrerequisiteIndex(featureID string) (int, error) {
	for i, prereq := range f.Prerequisites {
		if prereq.FeatureId == featureID {
			return i, nil
		}
	}
	return -1, errPrerequisiteNotFound
}
