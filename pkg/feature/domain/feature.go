// Copyright 2024 The Bucketeer Authors.
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
	"strconv"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	"github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	SecondsToStale         = 90 * 24 * 60 * 60 // 90 days
	secondsToReEvaluateAll = 30 * 24 * 60 * 60 // 30 days
	secondsForAdjustment   = 10                // 10 seconds
)

var (
	errClauseNotFound                = errors.New("feature: clause not found")
	errDefaultStrategyNotFound       = errors.New("feature: default strategy not found")
	errClauseAlreadyExists           = errors.New("feature: clause already exists")
	errRuleMustHaveAtLeastOneClause  = errors.New("feature: rule must have at least one clause")
	errClauseMustHaveAtLeastOneValue = errors.New("feature: clause must have at least one value")
	errRuleAlreadyExists             = errors.New("feature: rule already exists")
	errRuleNotFound                  = errors.New("feature: rule not found")
	errTargetNotFound                = errors.New("feature: target not found")
	errValueNotFound                 = errors.New("feature: value not found")
	errVariationInUse                = errors.New("feature: variation in use")
	errVariationNotFound             = errors.New("feature: variation not found")
	errVariationTypeUnmatched        = errors.New("feature: variation value and type are unmatched")
	errTagsMustHaveAtLeastOneTag     = errors.New("feature: tags must have at least one tag set")
	errUnsupportedStrategy           = errors.New("feature: unsupported strategy")
	errFeatureNotFound               = errors.New("feature: feature not found")
	errPrerequisiteVariationNotFound = errors.New("feature: prerequisite variation not found")
	ErrCycleExists                   = errors.New("feature: cycle exists in features")
	errPrerequisiteNotFound          = errors.New("feature: prerequisite not found")
	ErrAlreadyEnabled                = errors.New("feature: already enabled")
	ErrAlreadyDisabled               = errors.New("feature: already disabled")
	ErrLastUsedInfoNotFound          = errors.New("feature: last used info not found")
	errRulesOrderSizeNotEqual        = errors.New("feature: rules order size not equal")
	errRulesOrderDuplicateIDs        = errors.New("feature: rules order contains duplicate ids")
)

// TODO: think about splitting out ruleset / variation

type Feature struct {
	*feature.Feature
	ruleEvaluator
	strategyEvaluator
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
	for i := range variations {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		if err = f.AddVariation(id.String(), variations[i].Value, variations[i].Name, variations[i].Description); err != nil {
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

func (f *Feature) assignUser(
	user *userproto.User,
	segmentUsers []*feature.SegmentUser,
	flagVariations map[string]string,
) (*feature.Reason, *feature.Variation, error) {
	for _, pf := range f.Prerequisites {
		variation, ok := flagVariations[pf.FeatureId]
		if !ok {
			return nil, nil, errPrerequisiteVariationNotFound
		}
		if pf.VariationId != variation {
			if f.OffVariation != "" {
				variation, err := findVariation(f.OffVariation, f.Variations)
				return &feature.Reason{Type: feature.Reason_PREREQUISITE}, variation, err
			}
		}
	}
	// It doesn't assign the user in case of the feature is disabled and OffVariation is not set
	if !f.Enabled && f.OffVariation != "" {
		variation, err := findVariation(f.OffVariation, f.Variations)
		return &feature.Reason{Type: feature.Reason_OFF_VARIATION}, variation, err
	}
	// evaluate from top to bottom, return if one rule matches
	// evaluate targeting rules
	for i := range f.Targets {
		if contains(user.Id, f.Targets[i].Users) {
			variation, err := findVariation(f.Targets[i].Variation, f.Variations)
			return &feature.Reason{Type: feature.Reason_TARGET}, variation, err
		}
	}
	// evaluate ruleset
	rule := f.ruleEvaluator.Evaluate(f.Rules, user, segmentUsers)
	if rule != nil {
		variation, err := f.strategyEvaluator.Evaluate(
			rule.Strategy,
			user.Id,
			f.Variations,
			f.Feature.Id,
			f.Feature.SamplingSeed,
		)
		return &feature.Reason{
			Type:   feature.Reason_RULE,
			RuleId: rule.Id,
		}, variation, err
	}
	// use default strategy
	if f.DefaultStrategy == nil {
		return nil, nil, errDefaultStrategyNotFound
	}
	variation, err := f.strategyEvaluator.Evaluate(
		f.DefaultStrategy,
		user.Id,
		f.Variations,
		f.Feature.Id,
		f.Feature.SamplingSeed,
	)
	if err != nil {
		return nil, nil, err
	}
	return &feature.Reason{Type: feature.Reason_DEFAULT}, variation, nil
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
	if len(f.Tags) <= 1 {
		return errTagsMustHaveAtLeastOneTag
	}
	idx, err := index(tag, f.Tags)
	if err != nil {
		return err
	}
	f.Tags = append(f.Tags[:idx], f.Tags[idx+1:]...)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) Enable() error {
	if f.Enabled {
		return ErrAlreadyEnabled
	}
	f.Enabled = true
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) Disable() error {
	if !f.Enabled {
		return ErrAlreadyDisabled
	}
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
	if _, err := f.findRule(rule.Id); err == nil {
		return errRuleAlreadyExists
	}
	if err := validateStrategy(rule.Strategy, f.Variations); err != nil {
		return err
	}
	// TODO: rule validation needed?
	// - maybe check if 2 rules are the same (not id but logic)
	// - check if two rules are the same but have different targets
	f.Rules = append(f.Rules, rule)
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

func validateStrategy(strategy *feature.Strategy, variations []*feature.Variation) error {
	switch strategy.Type {
	case feature.Strategy_FIXED:
		return validateFixedStrategy(strategy.FixedStrategy, variations)
	case feature.Strategy_ROLLOUT:
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
	if err := validateVariation(f.VariationType, value); err != nil {
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

func validateVariation(variationType feature.Feature_VariationType, value string) error {
	switch variationType {
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
	if strategyContainsVariation(id, f.Feature.DefaultStrategy) {
		return errVariationInUse
	}
	if f.rulesContainsVariation(id) {
		return errVariationInUse
	}
	if f.OffVariation == id {
		return errVariationInUse
	}
	return nil
}

func (f *Feature) rulesContainsVariation(id string) bool {
	for _, r := range f.Feature.Rules {
		if ok := strategyContainsVariation(id, r.Strategy); ok {
			return true
		}
	}
	return false
}

func strategyContainsVariation(id string, strategy *feature.Strategy) bool {
	if strategy.Type == feature.Strategy_FIXED {
		if strategy.FixedStrategy.Variation == id {
			return true
		}
	} else if strategy.Type == feature.Strategy_ROLLOUT {
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
			return
		}
	}
}

func (f *Feature) removeVariationFromDefaultStrategy(variationID string) {
	if f.DefaultStrategy != nil && f.DefaultStrategy.Type == feature.Strategy_ROLLOUT {
		f.removeVariationFromRolloutStrategy(f.DefaultStrategy.RolloutStrategy, variationID)
	}
}

func (f *Feature) removeVariationFromRolloutStrategy(strategy *feature.RolloutStrategy, variationID string) {
	for i, v := range strategy.Variations {
		if v.Variation == variationID {
			strategy.Variations = append(strategy.Variations[:i], strategy.Variations[i+1:]...)
			return
		}
	}
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
	p := &feature.Prerequisite{FeatureId: fID, VariationId: variationID}
	f.Prerequisites = append(f.Prerequisites, p)
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) ChangePrerequisiteVariation(fID, variationID string) error {
	idx, err := f.findPrerequisite(fID)
	if err != nil {
		return err
	}
	f.Prerequisites[idx].VariationId = variationID
	f.UpdatedAt = time.Now().Unix()
	return nil
}

func (f *Feature) RemovePrerequisite(fID string) error {
	idx, err := f.findPrerequisite(fID)
	if err != nil {
		return err
	}
	f.Prerequisites = append(f.Prerequisites[:idx], f.Prerequisites[idx+1:]...)
	f.UpdatedAt = time.Now().Unix()
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
