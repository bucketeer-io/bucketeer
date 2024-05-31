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

package evaluation

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/maps"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

type Mark int

const (
	unvisited Mark = iota
	temporary
	permanently
)

const (
	secondsToReEvaluateAll = 30 * 24 * 60 * 60 // 30 days
	secondsForAdjustment   = 10                // 10 seconds
)

var (
	ErrCycleExists                   = errors.New("evaluator: cycle exists in features")
	ErrDefaultStrategyNotFound       = errors.New("evaluator: default strategy not found")
	ErrFeatureNotFound               = errors.New("evaluator: feature not found")
	ErrPrerequisiteVariationNotFound = errors.New("evaluator: prerequisite variation not found")
	ErrVariationNotFound             = errors.New("evaluator: variation not found")
	ErrUnsupportedStrategy           = errors.New("evaluator: unsupported strategy")
)

func EvaluationID(featureID string, featureVersion int32, userID string) string {
	return fmt.Sprintf("%s:%d:%s", featureID, featureVersion, userID)
}

type evaluator struct {
	ruleEvaluator
	strategyEvaluator
}

func NewEvaluator() *evaluator {
	return &evaluator{}
}

// This function will be removed once all the SDK clients are updated.
func (e *evaluator) EvaluateFeatures(
	fs []*ftproto.Feature,
	user *userproto.User,
	mapSegmentUsers map[string][]*ftproto.SegmentUser,
	targetTag string,
) (*ftproto.UserEvaluations, error) {
	return e.evaluate(fs, user, mapSegmentUsers, false, targetTag)
}

func (e *evaluator) EvaluateFeaturesByEvaluatedAt(
	fs []*ftproto.Feature,
	user *userproto.User,
	mapSegmentUsers map[string][]*ftproto.SegmentUser,
	prevUEID string,
	evaluatedAt int64,
	userAttributesUpdated bool,
	targetTag string,
) (*ftproto.UserEvaluations, error) {
	if prevUEID == "" {
		return e.evaluate(fs, user, mapSegmentUsers, true, targetTag)
	}
	now := time.Now()
	if evaluatedAt < now.Unix()-secondsToReEvaluateAll {
		return e.evaluate(fs, user, mapSegmentUsers, true, targetTag)
	}
	adjustedEvalAt := evaluatedAt - secondsForAdjustment
	updatedFeatures := make([]*ftproto.Feature, 0, len(fs))
	for _, feature := range fs {
		if feature.UpdatedAt > adjustedEvalAt {
			updatedFeatures = append(updatedFeatures, feature)
			continue
		}
		if userAttributesUpdated && len(feature.Rules) != 0 {
			updatedFeatures = append(updatedFeatures, feature)
		}
	}
	// If the UserEvaluationsID has changed, but both User Attributes and Feature Flags have not been updated,
	// it is considered unusual and a force update should be performed.
	if len(updatedFeatures) == 0 {
		return e.evaluate(fs, user, mapSegmentUsers, true, targetTag)
	}
	evalTargets, err := e.getEvalFeatures(updatedFeatures, fs)
	if err != nil {
		return nil, err
	}
	return e.evaluate(evalTargets, user, mapSegmentUsers, false, targetTag)
}

func (e *evaluator) evaluate(
	fs []*ftproto.Feature,
	user *userproto.User,
	mapSegmentUsers map[string][]*ftproto.SegmentUser,
	forceUpdate bool,
	targetTag string,
) (*ftproto.UserEvaluations, error) {

	flagVariations := map[string]string{}
	// fs need to be sorted in order from upstream to downstream.
	sortedFs, err := e.TopologicalSort(fs)
	if err != nil {
		return nil, err
	}
	evaluations := make([]*ftproto.Evaluation, 0, len(fs))
	archivedIDs := make([]string, 0)
	for _, feature := range sortedFs {
		if feature.Archived {
			// To keep response size small, the feature flags archived long time ago are excluded.
			if !e.isArchivedBeforeLastThirtyDays(feature) {
				archivedIDs = append(archivedIDs, feature.Id)
			}
			continue
		}
		var segmentUsers []*ftproto.SegmentUser
		for _, id := range e.ListSegmentIDs(feature) {
			segmentUsers = append(segmentUsers, mapSegmentUsers[id]...)
		}
		reason, variation, err := e.assignUser(feature, user, segmentUsers, flagVariations)
		if err != nil {
			return nil, err
		}
		// VariationId is used to check if prerequisite flag's result is what user expects it to be.
		flagVariations[feature.Id] = variation.Id
		// When the tag is set in the request,
		// it will return only the evaluations of flags that match the tag configured on the dashboard.
		// When empty, it will return all the evaluations of the flags in the environment.
		if targetTag != "" && !tagExist(feature.Tags, targetTag) {
			continue
		}
		// FIXME: Remove the next line when the Variation
		// no longer is being used
		// For security reasons, it removes the variation description
		variation.Description = ""
		evaluationID := EvaluationID(feature.Id, feature.Version, user.Id)
		evaluation := &ftproto.Evaluation{
			Id:             evaluationID,
			FeatureId:      feature.Id,
			FeatureVersion: feature.Version,
			UserId:         user.Id,
			VariationId:    variation.Id,
			VariationName:  variation.Name,
			VariationValue: variation.Value,
			Variation:      variation, // deprecated
			Reason:         reason,
		}
		evaluations = append(evaluations, evaluation)
	}
	// FIXME: Remove id once all SDKs will be updated.
	id := UserEvaluationsID(user.Id, user.Data, fs)
	userEvaluations := NewUserEvaluations(id, evaluations, archivedIDs, forceUpdate)
	return userEvaluations.UserEvaluations, nil
}

func tagExist(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

// This logic is based on https://en.wikipedia.org/wiki/Topological_sorting.
// Note: This algorithm is not an exact topological sort because the order is reversed (=from upstream to downstream).
func (e *evaluator) TopologicalSort(features []*ftproto.Feature) ([]*ftproto.Feature, error) {
	marks := map[string]Mark{}
	mapFeatures := map[string]*ftproto.Feature{}
	for _, f := range features {
		marks[f.Id] = unvisited
		mapFeatures[f.Id] = f
	}
	var sortedFeatures []*ftproto.Feature
	var sort func(f *ftproto.Feature) error
	sort = func(f *ftproto.Feature) error {
		if marks[f.Id] == permanently {
			return nil
		}
		if marks[f.Id] == temporary {
			return ErrCycleExists
		}
		marks[f.Id] = temporary
		for _, p := range f.Prerequisites {
			pf, ok := mapFeatures[p.FeatureId]
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

/*
IsArchivedBeforeLastThirtyDays returns a bool value
indicating whether the feature flag was archived within the last thirty days.
*/
func (e *evaluator) isArchivedBeforeLastThirtyDays(feature *ftproto.Feature) bool {
	if !feature.Archived {
		return false
	}
	now := time.Now()
	return feature.UpdatedAt < now.Unix()-secondsToReEvaluateAll
}

func (e *evaluator) ListSegmentIDs(feature *ftproto.Feature) []string {
	mapIDs := make(map[string]struct{})
	for _, r := range feature.Rules {
		for _, c := range r.Clauses {
			if c.Operator == ftproto.Clause_SEGMENT {
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

func (e *evaluator) assignUser(
	feature *ftproto.Feature,
	user *userproto.User,
	segmentUsers []*ftproto.SegmentUser,
	flagVariations map[string]string,
) (*ftproto.Reason, *ftproto.Variation, error) {
	for _, pf := range feature.Prerequisites {
		variation, ok := flagVariations[pf.FeatureId]
		if !ok {
			return nil, nil, ErrPrerequisiteVariationNotFound
		}
		if pf.VariationId != variation {
			if feature.OffVariation != "" {
				variation, err := findVariation(feature.OffVariation, feature.Variations)
				return &ftproto.Reason{Type: ftproto.Reason_PREREQUISITE}, variation, err
			}
		}
	}
	// It doesn't assign the user in case of the feature is disabled and OffVariation is not set
	if !feature.Enabled && feature.OffVariation != "" {
		variation, err := findVariation(feature.OffVariation, feature.Variations)
		return &ftproto.Reason{Type: ftproto.Reason_OFF_VARIATION}, variation, err
	}
	// evaluate from top to bottom, return if one rule matches
	// evaluate targeting rules
	for i := range feature.Targets {
		if contains(user.Id, feature.Targets[i].Users) {
			variation, err := findVariation(feature.Targets[i].Variation, feature.Variations)
			return &ftproto.Reason{Type: ftproto.Reason_TARGET}, variation, err
		}
	}
	// evaluate ruleset
	rule := e.ruleEvaluator.Evaluate(feature.Rules, user, segmentUsers)
	if rule != nil {
		variation, err := e.strategyEvaluator.Evaluate(
			rule.Strategy,
			user.Id,
			feature.Variations,
			feature.Id,
			feature.SamplingSeed,
		)
		return &ftproto.Reason{
			Type:   ftproto.Reason_RULE,
			RuleId: rule.Id,
		}, variation, err
	}
	// use default strategy
	if feature.DefaultStrategy == nil {
		return nil, nil, ErrDefaultStrategyNotFound
	}
	variation, err := e.strategyEvaluator.Evaluate(
		feature.DefaultStrategy,
		user.Id,
		feature.Variations,
		feature.Id,
		feature.SamplingSeed,
	)
	if err != nil {
		return nil, nil, err
	}
	return &ftproto.Reason{Type: ftproto.Reason_DEFAULT}, variation, nil
}

// GetPrerequisiteDownwards gets the features specified as prerequisite by the targetFeatures.
func (e *evaluator) GetPrerequisiteDownwards(
	targetFeatures, allFeatures []*ftproto.Feature,
) ([]*ftproto.Feature, error) {
	allFeaturesMap := make(map[string]*ftproto.Feature, len(allFeatures))
	for _, f := range allFeatures {
		allFeaturesMap[f.Id] = f
	}
	return maps.Values(e.getFeaturesDependedOnTargets(targetFeatures, allFeaturesMap)), nil
}

func (e *evaluator) getEvalFeatures(
	targetFeatures, allFeatures []*ftproto.Feature,
) ([]*ftproto.Feature, error) {
	all := make(map[string]*ftproto.Feature, len(allFeatures))
	for _, f := range allFeatures {
		all[f.Id] = f
	}

	evals1 := e.getFeaturesDependedOnTargets(targetFeatures, all)
	evals2 := e.getFeaturesDependsOnTargets(targetFeatures, all)
	evals := e.mapMerge(evals1, evals2)
	return maps.Values(evals), nil
}

// getFeaturesDependedOnTargets returns the features that are depended on the target features.
// targetFeatures are included in the result.
func (e *evaluator) getFeaturesDependedOnTargets(
	targets []*ftproto.Feature, all map[string]*ftproto.Feature,
) map[string]*ftproto.Feature {
	evals := make(map[string]*ftproto.Feature)
	var dfs func(f *ftproto.Feature)
	dfs = func(f *ftproto.Feature) {
		if _, ok := evals[f.Id]; ok {
			return
		}
		evals[f.Id] = f
		dmn := &domain.Feature{Feature: f}
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
func (e *evaluator) getFeaturesDependsOnTargets(
	targets []*ftproto.Feature, all map[string]*ftproto.Feature,
) map[string]*ftproto.Feature {
	evals := make(map[string]*ftproto.Feature)
	for _, f := range targets {
		evals[f.Id] = f
	}
	var dfs func(f *ftproto.Feature) bool
	dfs = func(f *ftproto.Feature) bool {
		if _, ok := evals[f.Id]; ok {
			return true
		}
		dmn := &domain.Feature{Feature: f}
		for _, fid := range dmn.FeatureIDsDependsOn() {
			if dfs(all[fid]) {
				evals[f.Id] = f
				return true
			}
		}
		for _, p := range f.Prerequisites {
			if dfs(all[p.FeatureId]) {
				evals[f.Id] = f
				return true
			}
		}
		return false
	}
	for _, f := range all {
		dfs(f)
	}
	return evals
}

func (e *evaluator) mapMerge(m1, m2 map[string]*ftproto.Feature) map[string]*ftproto.Feature {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func findVariation(v string, vs []*ftproto.Variation) (*ftproto.Variation, error) {
	for i := range vs {
		if vs[i].Id == v {
			return vs[i], nil
		}
	}
	return nil, ErrVariationNotFound
}

func contains(needle string, haystack []string) bool {
	for i := range haystack {
		if haystack[i] == needle {
			return true
		}
	}
	return false
}
