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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	"github.com/bucketeer-io/bucketeer/proto/feature"
)

type Mark int

const (
	unvisited Mark = iota
	temporary
	permanently
)

const (
	SecondsToStale         = 90 * 24 * 60 * 60 // 90 days
	secondsToReEvaluateAll = 30 * 24 * 60 * 60 // 30 days
)

var (
	errNameEmpty = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: name cannot be empty", "name")
	errClauseNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: clause not found", "clause")
	errClauseAttributeNotEmpty = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: clause attribute must be empty", "clause")
	errClauseAttributeEmpty = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: clause attribute cannot be empty", "clause")
	errClauseValuesEmpty = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: clause values cannot be empty", "clause_values")
	errClauseAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.FeaturePackageName, "feature: clause already exists")
	errRuleMustHaveAtLeastOneClause = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: rule must have at least one clause", "rule")
	errClauseMustHaveAtLeastOneValue = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: clause must have at least one value", "clause")
	errRuleAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.FeaturePackageName, "feature: rule already exists")
	errRuleIDRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: rule id required", "rule_id")
	errRuleRequired = pkgErr.NewErrorInvalidArgNil(
		pkgErr.FeaturePackageName, "feature: rule required", "rule")
	errRuleNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: rule not found", "rule")
	errTargetNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: target not found", "target")
	errVariationRequired = pkgErr.NewErrorInvalidArgNil(
		pkgErr.FeaturePackageName, "feature: variation required", "variation")
	errValueNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: value not found", "value")
	errVariationIDRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: variation id required", "variation_id")
	errVariationNameRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: variation name required", "variation_name")
	errVariationValueRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: variation value required", "variation_value")
	errVariationValueUnique = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: variation value must be unique", "variation")
	errVariationsMustHaveAtLeastTwoVariations = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: variations must have at least two variations", "variations")
	errInvalidDefaultOnVariationIndex = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: invalid default on variation index. Index is out of range", "variation")
	errInvalidDefaultOffVariationIndex = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: invalid default off variation index. Index is out of range", "variation")
	errTargetUsersRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: target users required", "target_users")
	errTargetUserRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: target user required", "target_user")
	ErrVariationInUse = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: variation in use", "variation")
	errVariationNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: variation not found", "variation")
	errVariationTypeUnmatched = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: variation value and type are unmatched", "variation")
	errStrategyRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: strategy required", "strategy")
	errUnsupportedStrategy = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: unsupported strategy", "strategy")
	errPrerequisiteAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.FeaturePackageName, "feature: prerequisite already exists")
	errPrerequisiteNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: prerequisite not found", "prerequisite")
	ErrLastUsedInfoNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: last used info not found", "last_used_info")
	errRulesOrderSizeNotEqual = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: rules order size not equal", "rules")
	errRulesOrderDuplicateIDs = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: rules order contains duplicate ids", "rules")
	ErrCycleExists = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: cycle exists in features", "features")
	errFeatureIDRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: feature id required", "feature_id")
	ErrFeatureNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: feature not found", "feature")
	ErrDefaultStrategyCannotBeBothFixedAndRollout = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: default strategy cannot be both fixed and rollout", "strategy")
	ErrRuleStrategyCannotBeEmpty = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.FeaturePackageName, "feature: rule strategy cannot be empty", "strategy")
	ErrInvalidAudiencePercentage = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: audience percentage must be between 0 and 100", "audience")
	ErrDefaultVariationNotFound = pkgErr.NewErrorNotFound(
		pkgErr.FeaturePackageName, "feature: default variation not found", "default_variation")
	ErrInvalidAudienceDefaultVariation = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName,
		"feature: default variation required when audience percentage is between 1 and 99",
		"variation",
	)
	ErrInvalidVariationWeightTotal = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.FeaturePackageName, "feature: variation weights must sum to 100%", "variation")
)

const (
	totalVariationWeight = int32(100000)
)

// TODO: think about splitting out ruleset / variation

type Feature struct {
	*feature.Feature
}

func NewFeature(
	id, name, description string,
	variationType feature.Feature_VariationType,
	variations []*feature.Variation,
	tags []string,
	defaultOnVariationIndex, defaultOffVariationIndex int,
	maintainer string,
) (*Feature, error) {
	f := &Feature{Feature: &feature.Feature{
		Id:            id,
		Name:          name,
		Description:   description,
		Version:       1,
		VariationType: variationType,
		CreatedAt:     time.Now().Unix(),
		Maintainer:    maintainer,
		Prerequisites: []*feature.Prerequisite{},
		Targets:       []*feature.Target{},
		Rules:         []*feature.Rule{},
	}}
	if len(variations) < 2 {
		return nil, errVariationsMustHaveAtLeastTwoVariations
	}
	if defaultOnVariationIndex < 0 || defaultOnVariationIndex >= len(variations) {
		return nil, errInvalidDefaultOnVariationIndex
	}
	if defaultOffVariationIndex < 0 || defaultOffVariationIndex >= len(variations) {
		return nil, errInvalidDefaultOffVariationIndex
	}
	for _, variation := range variations {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		if err = f.AddVariation(id.String(), variation.Value, variation.Name, variation.Description); err != nil {
			return nil, err
		}
	}
	f.Tags = append(f.Tags, tags...)
	if err := f.ChangeDefaultStrategy(&feature.Strategy{
		Type: feature.Strategy_FIXED,
		FixedStrategy: &feature.FixedStrategy{
			Variation: f.Variations[defaultOnVariationIndex].Id,
		},
	}); err != nil {
		return nil, err
	}
	if err := f.ChangeOffVariation(f.Variations[defaultOffVariationIndex].Id); err != nil {
		return nil, err
	}
	return f, nil
}

func findVariation(v string, vs []*feature.Variation) (*feature.Variation, error) {
	for i := range vs {
		if vs[i].Id == v {
			return vs[i], nil
		}
	}
	return nil, errVariationNotFound
}

func (f *Feature) Rename(name string) error {
	if name == "" {
		return errNameEmpty
	}
	f.Name = name
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeDescription(description string) error {
	f.Description = description
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeOffVariation(id string) error {
	_, err := findVariation(id, f.Variations)
	if err != nil {
		return err
	}
	f.OffVariation = id
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) AddTag(tag string) error {
	if contains(tag, f.Tags) {
		// output info log
		return nil
	}
	f.Tags = append(f.Tags, tag)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemoveTag(tag string) error {
	idx, err := index(tag, f.Tags)
	if err != nil {
		return err
	}
	if len(f.Tags) == 1 {
		f.Tags = []string{}
	} else {
		f.Tags = append(f.Tags[:idx], f.Tags[idx+1:]...)
	}
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) Enable() error {
	f.Enabled = true
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) Disable() error {
	f.Enabled = false
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) Archive() error {
	f.Archived = true
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) Unarchive() error {
	f.Archived = false
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) Delete() error {
	f.Deleted = true
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) AddUserToVariation(variation string, user string) error {
	idx, err := f.findTarget(variation)
	if err != nil {
		return err
	}
	if contains(user, f.Targets[idx].Users) {
		return nil
	}
	f.Targets[idx].Users = append(f.Targets[idx].Users, user)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemoveUserFromVariation(variation string, user string) error {
	idx, err := f.findTarget(variation)
	if err != nil {
		return err
	}
	uidx, err := index(user, f.Targets[idx].Users)
	if err != nil {
		return err
	}
	f.Targets[idx].Users = append(f.Targets[idx].Users[:uidx], f.Targets[idx].Users[uidx+1:]...)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) AddRule(rule *feature.Rule) error {
	if rule == nil {
		return errRuleRequired
	}
	if _, err := f.findRule(rule.Id); err == nil {
		return errRuleAlreadyExists
	}
	if err := validateClauses(rule.Clauses); err != nil {
		return err
	}
	if err := validateStrategy(rule.Strategy, f.Variations); err != nil {
		return err
	}
	f.Rules = append(f.Rules, rule)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeRule(rule *feature.Rule) error {
	if rule == nil {
		return errRuleRequired
	}
	idx, err := f.findRule(rule.Id)
	if err != nil {
		return err
	}
	if err := validateClauses(rule.Clauses); err != nil {
		return err
	}
	if err := validateStrategy(rule.Strategy, f.Variations); err != nil {
		return err
	}
	f.Rules[idx] = rule
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemoveRule(ruleID string) error {
	if ruleID == "" {
		return errRuleIDRequired
	}
	idx, err := f.findRule(ruleID)
	if err != nil {
		return err
	}
	f.Rules = append(f.Rules[:idx], f.Rules[idx+1:]...)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeRuleStrategy(ruleID string, strategy *feature.Strategy) error {
	idx, err := f.findRule(ruleID)
	if err != nil {
		return errRuleNotFound
	}
	if err := validateStrategy(strategy, f.Variations); err != nil {
		return err
	}
	f.Rules[idx].Strategy = strategy
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeRulesOrder(ruleIDs []string) error {
	if len(ruleIDs) != len(f.Rules) {
		return errRulesOrderSizeNotEqual
	}
	rules := make([]*feature.Rule, 0, len(ruleIDs))
	for _, ruleID := range ruleIDs {
		for _, r := range rules {
			if r.Id == ruleID {
				return errRulesOrderDuplicateIDs
			}
		}
		rule, err := f.getRule(ruleID)
		if err != nil {
			return errRuleNotFound
		}
		rules = append(rules, rule)
	}
	f.Rules = rules
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) getRule(id string) (*feature.Rule, error) {
	for _, rule := range f.Rules {
		if rule.Id == id {
			return rule, nil
		}
	}
	return nil, errRuleNotFound
}

func validateClauses(clauses []*feature.Clause) error {
	if len(clauses) == 0 {
		return fmt.Errorf("feature: rule must have at least one clause")
	}
	ids := make(map[string]struct{}, len(clauses))
	for _, c := range clauses {
		if c.Id == "" {
			return errors.New("feature: clause id cannot be empty")
		}
		if err := uuid.ValidateUUID(c.Id); err != nil {
			return fmt.Errorf("feature: clause %s id is invalid: %w", c.Id, err)
		}
		if _, ok := ids[c.Id]; ok {
			return fmt.Errorf("feature: clause %s already exists", c.Id)
		}
		ids[c.Id] = struct{}{}
		if err := validateClauseValues(c); err != nil {
			return err
		}
	}
	return nil
}

// validateClauseValues checks if attribute and values are valid according to type.
// WARNING: This function does not check if the feature id and variation id are valid,
// Call ValidateFeatureDependencies to check if they are so.
func validateClauseValues(c *feature.Clause) error {
	switch c.Operator {
	case feature.Clause_SEGMENT:
		if c.Attribute != "" {
			return errClauseAttributeNotEmpty
		}
		if len(c.Values) == 0 {
			return errClauseValuesEmpty
		}
	default:
		if c.Attribute == "" {
			return errClauseAttributeEmpty
		}
		if len(c.Values) == 0 {
			return errClauseValuesEmpty
		}
	}
	return nil
}

func validateStrategy(strategy *feature.Strategy, variations []*feature.Variation) error {
	if strategy == nil {
		return errStrategyRequired
	}
	switch strategy.Type {
	case feature.Strategy_FIXED:
		if strategy.FixedStrategy == nil {
			return ErrRuleStrategyCannotBeEmpty
		}
		if strategy.RolloutStrategy != nil {
			return ErrDefaultStrategyCannotBeBothFixedAndRollout
		}
		return validateFixedStrategy(strategy.FixedStrategy, variations)
	case feature.Strategy_ROLLOUT:
		if strategy.RolloutStrategy == nil {
			return ErrRuleStrategyCannotBeEmpty
		}
		if strategy.FixedStrategy != nil {
			return ErrDefaultStrategyCannotBeBothFixedAndRollout
		}
		return validateRolloutStrategy(strategy.RolloutStrategy, variations)
	default:
		return errUnsupportedStrategy
	}
}

func validateRolloutStrategy(strategy *feature.RolloutStrategy, variations []*feature.Variation) error {
	for _, v := range strategy.Variations {
		if _, err := findVariation(v.Variation, variations); err != nil {
			return errVariationNotFound
		}
	}

	// Validate audience if present
	if strategy.Audience != nil {
		// Validate percentage range
		if strategy.Audience.Percentage < 0 || strategy.Audience.Percentage > 100 {
			return ErrInvalidAudiencePercentage
		}

		// If percentage is between 1-99, default variation must be specified and exist
		if strategy.Audience.Percentage > 0 && strategy.Audience.Percentage < 100 {
			if strategy.Audience.DefaultVariation == "" {
				return ErrInvalidAudienceDefaultVariation
			}
			if _, err := findVariation(strategy.Audience.DefaultVariation, variations); err != nil {
				return ErrDefaultVariationNotFound
			}
		}
	}

	// Validate variation weights
	totalWeight := int32(0)
	for _, v := range strategy.Variations {
		totalWeight += v.Weight
	}
	if totalWeight != totalVariationWeight {
		return ErrInvalidVariationWeightTotal
	}

	return nil
}

func validateFixedStrategy(strategy *feature.FixedStrategy, variations []*feature.Variation) error {
	if _, err := findVariation(strategy.Variation, variations); err != nil {
		return errVariationNotFound
	}
	return nil
}

func (f *Feature) DeleteRule(rule string) error {
	idx, err := f.findRule(rule)
	if err != nil {
		return err
	}
	f.Rules = append(f.Rules[:idx], f.Rules[idx+1:]...)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) AddClause(rule string, clause *feature.Clause) error {
	if err := validateClauseValues(clause); err != nil {
		return err
	}
	// TODO: do same validation as in addrule?
	idx, err := f.findRule(rule)
	if err != nil {
		return err
	}
	f.Rules[idx].Clauses = append(f.Rules[idx].Clauses, clause)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) DeleteClause(rule string, clause string) error {
	ruleIdx, err := f.findRule(rule)
	if err != nil {
		return err
	}
	idx, err := f.findClause(clause, f.Rules[ruleIdx].Clauses)
	if err != nil {
		return err
	}
	f.Rules[ruleIdx].Clauses = append(f.Rules[ruleIdx].Clauses[:idx], f.Rules[ruleIdx].Clauses[idx+1:]...)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeClauseAttribute(rule string, clause string, attribute string) error {
	ruleIdx, err := f.findRule(rule)
	if err != nil {
		return err
	}
	idx, err := f.findClause(clause, f.Rules[ruleIdx].Clauses)
	if err != nil {
		return err
	}
	if f.Rules[ruleIdx].Clauses[idx].Attribute == attribute {
		// TODO: should something be returned so no event is created?
		return nil
	}
	f.Rules[ruleIdx].Clauses[idx].Attribute = attribute
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeClauseOperator(rule string, clause string, operator feature.Clause_Operator) error {
	ruleIdx, err := f.findRule(rule)
	if err != nil {
		return err
	}
	idx, err := f.findClause(clause, f.Rules[ruleIdx].Clauses)
	if err != nil {
		return err
	}
	if f.Rules[ruleIdx].Clauses[idx].Operator == operator {
		// TODO: same as attribute. maybe stop event from being generated
		return nil
	}
	f.Rules[ruleIdx].Clauses[idx].Operator = operator
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) AddClauseValue(rule string, clause string, value string) error {
	ruleIdx, err := f.findRule(rule)
	if err != nil {
		return err
	}
	idx, err := f.findClause(clause, f.Rules[ruleIdx].Clauses)
	if err != nil {
		return err
	}
	_, err = index(value, f.Rules[ruleIdx].Clauses[idx].Values)
	if err == nil {
		// TODO: same as attribute. maybe stop event from being generated
		return nil
	}
	f.Rules[ruleIdx].Clauses[idx].Values = append(f.Rules[ruleIdx].Clauses[idx].Values, value)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemoveClauseValue(rule string, clause string, value string) error {
	ruleIdx, err := f.findRule(rule)
	if err != nil {
		return err
	}
	clauseIdx, err := f.findClause(clause, f.Rules[ruleIdx].Clauses)
	if err != nil {
		return err
	}
	idx, err := index(value, f.Rules[ruleIdx].Clauses[clauseIdx].Values)
	if err != nil {
		// TODO: same as attribute. maybe stop event from being generated
		return nil
	}
	f.Rules[ruleIdx].Clauses[clauseIdx].Values = append(
		f.Rules[ruleIdx].Clauses[clauseIdx].Values[:idx],
		f.Rules[ruleIdx].Clauses[clauseIdx].Values[idx+1:]...,
	)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) AddVariation(id string, value string, name string, description string) error {
	if id == "" {
		return errVariationIDRequired
	}
	if name == "" {
		return errVariationNameRequired
	}
	if err := f.validateVariationValue(id, value); err != nil {
		return err
	}
	variation := &feature.Variation{
		Id:          id,
		Value:       value,
		Name:        name,
		Description: description,
	}
	f.Variations = append(f.Variations, variation)
	f.addTarget(id)
	f.addVariationToRules(id)
	f.addVariationToDefaultStrategy(id)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) validateVariationValue(id, value string) error {
	if value == "" {
		return errVariationValueRequired
	}
	// Check for duplicate values
	for _, existingVar := range f.Variations {
		if existingVar.Id != id && existingVar.Value == value {
			return errVariationValueUnique
		}
	}
	switch f.VariationType {
	case feature.Feature_BOOLEAN:
		if value != "true" && value != "false" {
			return errVariationTypeUnmatched
		}
	case feature.Feature_NUMBER:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return errVariationTypeUnmatched
		}
	case feature.Feature_JSON:
		var js map[string]interface{}
		var jsArray []interface{}
		if json.Unmarshal([]byte(value), &js) == nil || json.Unmarshal([]byte(value), &jsArray) == nil {
			return nil
		}
		return errVariationTypeUnmatched
	}
	return nil
}

func (f *Feature) addTarget(variationID string) {
	target := &feature.Target{
		Variation: variationID,
		Users:     []string{},
	}
	f.Targets = append(f.Targets, target)
}

func (f *Feature) addVariationToRules(variationID string) {
	for _, rule := range f.Rules {
		if rule.Strategy.Type == feature.Strategy_ROLLOUT {
			f.addVariationToRolloutStrategy(rule.Strategy.RolloutStrategy, variationID)
		}
	}
}

func (f *Feature) addVariationToDefaultStrategy(variationID string) {
	if f.DefaultStrategy != nil && f.DefaultStrategy.Type == feature.Strategy_ROLLOUT {
		f.addVariationToRolloutStrategy(f.DefaultStrategy.RolloutStrategy, variationID)
	}
}

func (f *Feature) addVariationToRolloutStrategy(strategy *feature.RolloutStrategy, variationID string) {
	strategy.Variations = append(strategy.Variations, &feature.RolloutStrategy_Variation{
		Variation: variationID,
		Weight:    0,
	})
}

func (f *Feature) ChangeVariation(variation *feature.Variation) error {
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
	f.Variations[idx] = variation
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemoveVariation(id string) error {
	idx, err := f.findVariationIndex(id)
	if err != nil {
		return err
	}
	if err = f.validateRemoveVariation(id); err != nil {
		return err
	}
	if err = f.removeTarget(id); err != nil {
		return err
	}
	f.removeVariationFromRules(id)
	f.removeVariationFromDefaultStrategy(id)
	f.Variations = append(f.Variations[:idx], f.Variations[idx+1:]...)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) validateRemoveVariation(id string) error {
	if len(f.Variations) <= 2 {
		return errVariationsMustHaveAtLeastTwoVariations
	}
	if f.OffVariation == id {
		return ErrVariationInUse
	}
	// Check if the individual targeting has any users
	idx, err := f.findTarget(id)
	if err != nil {
		return err
	}
	if len(f.Targets[idx].Users) > 0 {
		return ErrVariationInUse
	}
	if strategyContainsVariation(id, f.DefaultStrategy) {
		return ErrVariationInUse
	}
	if f.rulesContainsVariation(id) {
		return ErrVariationInUse
	}
	return nil
}

func (f *Feature) rulesContainsVariation(id string) bool {
	for _, r := range f.Rules {
		if ok := strategyContainsVariation(id, r.Strategy); ok {
			return true
		}
	}
	return false
}

func strategyContainsVariation(id string, strategy *feature.Strategy) bool {
	switch strategy.Type {
	case feature.Strategy_FIXED:
		if strategy.FixedStrategy.Variation == id {
			return true
		}
	case feature.Strategy_ROLLOUT:
		for _, v := range strategy.RolloutStrategy.Variations {
			if v.Variation == id && v.Weight > 0 {
				return true
			}
		}
	}
	return false
}

func (f *Feature) removeTarget(variationID string) error {
	idx, err := f.findTarget(variationID)
	if err != nil {
		return err
	}
	f.Targets = append(f.Targets[:idx], f.Targets[idx+1:]...)
	return nil
}

func (f *Feature) removeVariationFromRules(variationID string) {
	for _, rule := range f.Rules {
		if rule.Strategy.Type == feature.Strategy_ROLLOUT {
			f.removeVariationFromRolloutStrategy(rule.Strategy.RolloutStrategy, variationID)
		}
	}
}

func (f *Feature) removeVariationFromDefaultStrategy(variationID string) {
	if f.DefaultStrategy != nil && f.DefaultStrategy.Type == feature.Strategy_ROLLOUT {
		f.removeVariationFromRolloutStrategy(f.DefaultStrategy.RolloutStrategy, variationID)
	}
}

func (f *Feature) removeVariationFromRolloutStrategy(strategy *feature.RolloutStrategy, variationID string) {
	// Remove all instances of the variation, regardless of weight
	filteredVariations := make([]*feature.RolloutStrategy_Variation, 0, len(strategy.Variations))
	for _, v := range strategy.Variations {
		if v.Variation != variationID {
			filteredVariations = append(filteredVariations, v)
		}
	}
	strategy.Variations = filteredVariations
}

func (f *Feature) ChangeVariationValue(id string, value string) error {
	idx, err := f.findVariationIndex(id)
	if err != nil {
		return err
	}
	f.Variations[idx].Value = value
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeVariationName(id string, name string) error {
	idx, err := f.findVariationIndex(id)
	if err != nil {
		return err
	}
	f.Variations[idx].Name = name
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeVariationDescription(id string, description string) error {
	idx, err := f.findVariationIndex(id)
	if err != nil {
		return err
	}
	f.Variations[idx].Description = description
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeDefaultStrategy(s *feature.Strategy) error {
	if err := validateStrategy(s, f.Variations); err != nil {
		return err
	}
	f.DefaultStrategy = s
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeFixedStrategy(ruleID string, strategy *feature.FixedStrategy) error {
	ruleIdx, err := f.findRule(ruleID)
	if err != nil {
		return err
	}
	if _, err := findVariation(strategy.Variation, f.Variations); err != nil {
		return err
	}
	f.Rules[ruleIdx].Strategy.FixedStrategy = strategy
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangeRolloutStrategy(ruleID string, strategy *feature.RolloutStrategy) error {
	ruleIdx, err := f.findRule(ruleID)
	if err != nil {
		return err
	}
	for _, v := range strategy.Variations {
		if _, err := findVariation(v.Variation, f.Variations); err != nil {
			return err
		}
	}

	// Validate rollout strategy weights
	if err := validateRolloutStrategy(strategy, f.Variations); err != nil {
		return err
	}

	f.Rules[ruleIdx].Strategy.RolloutStrategy = strategy
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ListSegmentIDs() []string {
	mapIDs := make(map[string]struct{})
	for _, r := range f.Rules {
		for _, c := range r.Clauses {
			if c.Operator == feature.Clause_SEGMENT {
				for _, v := range c.Values {
					mapIDs[v] = struct{}{}
				}
			}
		}
	}
	ids := make([]string, 0, len(mapIDs))
	for id := range mapIDs {
		ids = append(ids, id)
	}
	return ids
}

func (f *Feature) IncrementVersion() error {
	f.Version++
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ResetSamplingSeed() error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	f.SamplingSeed = id.String()
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) AddPrerequisite(fID, variationID string) error {
	if _, err := f.findPrerequisite(fID); err == nil {
		return errPrerequisiteAlreadyExists
	}
	if err := validatePrerequisite(fID, variationID); err != nil {
		return err
	}
	p := &feature.Prerequisite{FeatureId: fID, VariationId: variationID}
	f.Prerequisites = append(f.Prerequisites, p)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangePrerequisiteVariation(fID, variationID string) error {
	if err := validatePrerequisite(fID, variationID); err != nil {
		return err
	}
	idx, err := f.findPrerequisite(fID)
	if err != nil {
		return err
	}
	f.Prerequisites[idx].VariationId = variationID
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemovePrerequisite(fID string) error {
	if fID == "" {
		return errFeatureIDRequired
	}
	idx, err := f.findPrerequisite(fID)
	if err != nil {
		return err
	}
	f.Prerequisites = append(f.Prerequisites[:idx], f.Prerequisites[idx+1:]...)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func validatePrerequisite(featureID, variationID string) error {
	if featureID == "" {
		return errFeatureIDRequired
	}
	if variationID == "" {
		return errVariationIDRequired
	}
	return nil
}

func (f *Feature) findPrerequisite(fID string) (int, error) {
	for i := range f.Prerequisites {
		if f.Prerequisites[i].FeatureId == fID {
			return i, nil
		}
	}
	return -1, errPrerequisiteNotFound
}

func (f *Feature) IsStale(t time.Time) bool {
	if f.LastUsedInfo == nil {
		return false
	}
	if (t.Unix() - f.LastUsedInfo.LastUsedAt) < SecondsToStale {
		return false
	}
	return true
}

func (f *Feature) IsDisabledAndOffVariationEmpty() bool {
	if f.Enabled {
		return false
	}
	return f.OffVariation == ""
}

/*
IsArchivedBeforeLastThirtyDays returns a bool value
indicating whether the feature flag was archived within the last thirty days.
*/
func (f *Feature) IsArchivedBeforeLastThirtyDays() bool {
	if !f.Archived {
		return false
	}
	now := time.Now()
	return f.UpdatedAt < now.Unix()-secondsToReEvaluateAll
}

func (f *Feature) findTarget(id string) (int, error) {
	for i := range f.Targets {
		if f.Targets[i].Variation == id {
			return i, nil
		}
	}
	return -1, errTargetNotFound
}

func (f *Feature) findVariationIndex(id string) (int, error) {
	for i := range f.Variations {
		if f.Variations[i].Id == id {
			return i, nil
		}
	}
	return -1, errVariationNotFound
}

func (f *Feature) findRule(id string) (int, error) {
	for i := range f.Rules {
		if f.Rules[i].Id == id {
			return i, nil
		}
	}
	return -1, errRuleNotFound
}

// TODO: this should be on Rule.. should wrap Rule..
// or maybe just find clause directly without finding the rule first.
func (f *Feature) findClause(id string, clauses []*feature.Clause) (int, error) {
	for i := range clauses {
		if clauses[i].Id == id {
			return i, nil
		}
	}
	return -1, nil
}

// TODO: this should be on Clause.. should wrap Clause.. do you see a pattern here?
func index(needle string, haystack []string) (int, error) {
	for i := range haystack {
		if haystack[i] == needle {
			return i, nil
		}
	}
	return -1, errValueNotFound
}

func contains(needle string, haystack []string) bool {
	for i := range haystack {
		if haystack[i] == needle {
			return true
		}
	}
	return false
}

// FeatureIDsDependsOn returns the ids of the features that this feature depends on.
func (f *Feature) FeatureIDsDependsOn() []string {
	ids := []string{}
	for _, p := range f.Prerequisites {
		ids = append(ids, p.FeatureId)
	}
	for _, p := range f.Rules {
		for _, c := range p.Clauses {
			if c.Operator == feature.Clause_FEATURE_FLAG {
				ids = append(ids, c.Attribute)
			}
		}
	}
	return ids
}

func (f *Feature) Clone(
	maintainer string,
) (*Feature, error) {
	now := time.Now().Unix()
	newFeature := &Feature{Feature: &feature.Feature{
		Id:              f.Id,
		Name:            f.Name,
		Description:     f.Description,
		Enabled:         false,
		Deleted:         false,
		Version:         1,
		CreatedAt:       now,
		UpdatedAt:       now,
		Variations:      f.Variations,
		Prerequisites:   []*feature.Prerequisite{},
		Targets:         f.Targets,
		Rules:           f.Rules,
		DefaultStrategy: f.DefaultStrategy,
		OffVariation:    f.OffVariation,
		Tags:            f.Tags,
		Maintainer:      maintainer,
		VariationType:   f.VariationType,
		Archived:        false,
	}}
	for i := range newFeature.Variations {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		if newFeature.Variations[i].Id == newFeature.OffVariation {
			newFeature.OffVariation = id.String()
		}
		for idx := range newFeature.Targets {
			if newFeature.Targets[idx].Variation == newFeature.Variations[i].Id {
				newFeature.Targets[idx].Variation = id.String()
				break
			}
		}
		if err = updateStrategyVariationID(newFeature.Variations[i].Id, id.String(), newFeature.DefaultStrategy); err != nil {
			return nil, err
		}
		for idx := range newFeature.Rules {
			err = updateStrategyVariationID(newFeature.Variations[i].Id, id.String(), newFeature.Rules[idx].Strategy)
			if err != nil {
				return nil, err
			}
		}
		newFeature.Variations[i].Id = id.String()
	}
	return newFeature, nil
}

func updateStrategyVariationID(varID, uID string, s *feature.Strategy) error {
	switch s.Type {
	case feature.Strategy_FIXED:
		if varID == s.FixedStrategy.Variation {
			s.FixedStrategy.Variation = uID
		}
	case feature.Strategy_ROLLOUT:
		for i := range s.RolloutStrategy.Variations {
			if s.RolloutStrategy.Variations[i].Variation == varID {
				s.RolloutStrategy.Variations[i].Variation = uID
				break
			}
		}
	default:
		return errUnsupportedStrategy
	}
	return nil
}

func (f *Feature) AddTargetUsers(target *feature.Target) error {
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
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemoveTargetUsers(target *feature.Target) error {
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
			return err
		}
		f.Targets[idx].Users = append(f.Targets[idx].Users[:uidx], f.Targets[idx].Users[uidx+1:]...)
	}
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func validateRules(rules []*feature.Rule, variations []*feature.Variation) error {
	ids := make(map[string]struct{}, len(variations))
	for _, r := range rules {
		if r == nil {
			return errRuleRequired
		}
		if r.Id == "" {
			return errRuleIDRequired
		}
		if err := uuid.ValidateUUID(r.Id); err != nil {
			return fmt.Errorf("feature: rule %s id is invalid: %w", r.Id, err)
		}
		if _, ok := ids[r.Id]; ok {
			return errRuleAlreadyExists
		}
		ids[r.Id] = struct{}{}
		if r.Strategy == nil {
			return errStrategyRequired
		}
		if err := validateStrategy(r.Strategy, variations); err != nil {
			return err
		}
		if err := validateClauses(r.Clauses); err != nil {
			return err
		}
	}
	return nil
}

func (f *Feature) validateVariationChanges(variationChanges []*feature.VariationChange) error {
	for _, change := range variationChanges {
		// For individual variation changes, only validate the variation value and type
		if change.Variation == nil {
			return errVariationRequired
		}

		switch change.ChangeType {
		case feature.ChangeType_CREATE:
			if err := uuid.ValidateUUID(change.Variation.Id); err != nil {
				return err
			}
			if change.Variation.Name == "" {
				return errVariationNameRequired
			}
			if err := f.validateVariationValue(change.Variation.Id, change.Variation.Value); err != nil {
				return err
			}
		case feature.ChangeType_UPDATE:
			// For UPDATE operations, validate ID exists
			if change.Variation.Id == "" {
				return errVariationIDRequired
			}
			if _, err := f.findVariationIndex(change.Variation.Id); err != nil {
				return err
			}
			if change.Variation.Name == "" {
				return errVariationNameRequired
			}
			if err := f.validateVariationValue(change.Variation.Id, change.Variation.Value); err != nil {
				return err
			}
		case feature.ChangeType_DELETE:
			// For DELETE operations, just validate the variation exists
			if change.Variation.Id == "" {
				return errVariationIDRequired
			}
			if _, err := f.findVariationIndex(change.Variation.Id); err != nil {
				return err
			}
		}
	}

	// Validate that the feature will have at least 2 variations after applying changes
	deletedCount := 0
	for _, change := range variationChanges {
		if change.ChangeType == feature.ChangeType_DELETE {
			deletedCount++
		}
	}
	if len(f.Variations)-deletedCount < 2 {
		return errVariationsMustHaveAtLeastTwoVariations
	}

	return nil
}

func ValidateFeatureDependencies(fs []*feature.Feature) error {
	_, err := TopologicalSort(fs)
	return err
}

// ValidateVariationUsage validates that the given variations are not used as cross-feature dependencies.
func ValidateVariationUsage(
	features []*feature.Feature,
	targetFeatureID string,
	deletedVariations map[string]string, // variationID -> variationValue
) error {
	for _, f := range features {
		if f.Id == targetFeatureID {
			continue // Skip the feature being updated
		}

		// Check if any deleted variations are used in prerequisites
		for _, prerequisite := range f.Prerequisites {
			if prerequisite.FeatureId == targetFeatureID {
				if _, found := deletedVariations[prerequisite.VariationId]; found {
					return ErrVariationInUse
				}
			}
		}

		// Check if any deleted variations are used in FEATURE_FLAG rules
		for _, rule := range f.Rules {
			for _, clause := range rule.Clauses {
				if clause.Operator == feature.Clause_FEATURE_FLAG && clause.Attribute == targetFeatureID {
					// FEATURE_FLAG clause values contain variation IDs, not values
					// We should check if any clause values match deleted variation IDs
					for _, clValue := range clause.Values {
						if _, found := deletedVariations[clValue]; found {
							return ErrVariationInUse
						}
					}
				}
			}
		}
	}
	return nil
}

// This logic is based on https://en.wikipedia.org/wiki/Topological_sorting.
// Note: This algorithm is not an exact topological sort because the order is reversed (=from upstream to downstream).
func TopologicalSort(features []*feature.Feature) ([]*feature.Feature, error) {
	marks := map[string]Mark{}
	mapFeatures := map[string]*feature.Feature{}
	for _, f := range features {
		marks[f.Id] = unvisited
		mapFeatures[f.Id] = f
	}
	var sortedFeatures []*feature.Feature
	var sort func(f *feature.Feature) error
	sort = func(f *feature.Feature) error {
		if marks[f.Id] == permanently {
			return nil
		}
		if marks[f.Id] == temporary {
			return ErrCycleExists
		}
		marks[f.Id] = temporary
		df := &Feature{Feature: f}
		for _, fid := range df.FeatureIDsDependsOn() {
			pf, ok := mapFeatures[fid]
			if !ok {
				return ErrFeatureNotFound
			}
			if err := sort(pf); err != nil {
				return err
			}
		}
		marks[f.Id] = permanently
		sortedFeatures = append(sortedFeatures, f)
		return nil
	}
	for _, f := range features {
		if marks[f.Id] != unvisited {
			continue
		}
		if err := sort(f); err != nil {
			return nil, err
		}
	}
	return sortedFeatures, nil
}

// getFeaturesDependedOnTargets returns the features that are depended on the target features.
// targetFeatures are included in the result.
func GetFeaturesDependedOnTargets(
	targets []*feature.Feature, all map[string]*feature.Feature,
) map[string]*feature.Feature {
	evals := make(map[string]*feature.Feature)
	var dfs func(f *feature.Feature)
	dfs = func(f *feature.Feature) {
		if _, ok := evals[f.Id]; ok {
			return
		}
		evals[f.Id] = f
		dmn := &Feature{Feature: f}
		for _, fid := range dmn.FeatureIDsDependsOn() {
			dfs(all[fid])
		}
	}
	for _, f := range targets {
		dfs(f)
	}
	return evals
}

// getFeaturesDependsOnTargets returns the features that depend on the target features.
// targetFeatures are included in the result.
func GetFeaturesDependsOnTargets(
	targets []*feature.Feature, all map[string]*feature.Feature,
) map[string]*feature.Feature {
	evals := make(map[string]*feature.Feature)
	for _, f := range targets {
		evals[f.Id] = f
	}
	var dfs func(f *feature.Feature) bool
	dfs = func(f *feature.Feature) bool {
		if _, ok := evals[f.Id]; ok {
			return true
		}
		dmn := &Feature{Feature: f}
		for _, fid := range dmn.FeatureIDsDependsOn() {
			if dfs(all[fid]) {
				evals[f.Id] = f
				return true
			}
		}
		return false
	}
	for _, f := range all {
		// Skip if the f is target feature.
		dfs(f)
	}
	return evals
}

// HasFeaturesDependsOnTargets returns true if there are features that depend on the target features.
// This is a thin wrapper of GetFeaturesDependsOnTargets.
func HasFeaturesDependsOnTargets(
	targets []*feature.Feature, all []*feature.Feature,
) bool {
	allfs := make(map[string]*feature.Feature, len(all))
	for _, f := range all {
		allfs[f.Id] = f
	}
	deps := GetFeaturesDependsOnTargets(targets, allfs)
	for _, tgt := range targets {
		delete(deps, tgt.Id)
	}
	return len(deps) > 0
}

func validateOffVariation(id string, variations []*feature.Variation) error {
	if id == "" {
		return errors.New("feature: off variation cannot be empty")
	}
	if _, err := findVariation(id, variations); err != nil {
		return err
	}
	return nil
}

func validateTargets(targets []*feature.Target, variations []*feature.Variation) error {
	for _, t := range targets {
		if _, err := findVariation(t.Variation, variations); err != nil {
			return errTargetNotFound
		}
	}
	return nil
}

// unique returns a new slice with duplicate values removed
func unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
