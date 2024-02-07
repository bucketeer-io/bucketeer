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
	"fmt"
	"time"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

type Mark int

const (
	unvisited Mark = iota
	temporary
	permanently
)

func EvaluationID(featureID string, featureVersion int32, userID string) string {
	return fmt.Sprintf("%s:%d:%s", featureID, featureVersion, userID)
}

// Deprecated: use EvaluateFeaturesByEvaluatedAt instead.
// This function will be removed once all the SDK clients are updated.
func EvaluateFeatures(
	fs []*featureproto.Feature,
	user *userproto.User,
	mapSegmentUsers map[string][]*featureproto.SegmentUser,
	targetTag string,
) (*featureproto.UserEvaluations, error) {
	return evaluate(fs, user, mapSegmentUsers, false, targetTag)
}

func EvaluateFeaturesByEvaluatedAt(
	fs []*featureproto.Feature,
	user *userproto.User,
	mapSegmentUsers map[string][]*featureproto.SegmentUser,
	prevUEID string,
	evaluatedAt int64,
	userAttributesUpdated bool,
	targetTag string,
) (*featureproto.UserEvaluations, error) {
	if prevUEID == "" {
		return evaluate(fs, user, mapSegmentUsers, true, targetTag)
	}
	now := time.Now()
	if evaluatedAt < now.Unix()-secondsToReEvaluateAll {
		return evaluate(fs, user, mapSegmentUsers, true, targetTag)
	}
	adjustedEvalAt := evaluatedAt - secondsForAdjustment
	updatedFeatures := make([]*featureproto.Feature, 0, len(fs))
	for _, f := range fs {
		feature := &Feature{Feature: f}
		if feature.UpdatedAt > adjustedEvalAt {
			updatedFeatures = append(updatedFeatures, f)
			continue
		}
		if userAttributesUpdated && len(feature.Rules) != 0 {
			updatedFeatures = append(updatedFeatures, f)
		}
	}
	// If the UserEvaluationsID has changed, but both User Attributes and Feature Flags have not been updated,
	// it is considered unusual and a force update should be performed.
	if len(updatedFeatures) == 0 {
		return evaluate(fs, user, mapSegmentUsers, true, targetTag)
	}
	featuresHavePrerequisite := getFeaturesHavePrerequisite(fs)
	evalTargets := GetPrerequisiteUpwards(updatedFeatures, featuresHavePrerequisite)
	return evaluate(evalTargets, user, mapSegmentUsers, false, targetTag)
}

func evaluate(
	fs []*featureproto.Feature,
	user *userproto.User,
	mapSegmentUsers map[string][]*featureproto.SegmentUser,
	forceUpdate bool,
	targetTag string,
) (*featureproto.UserEvaluations, error) {
	flagVariations := map[string]string{}
	// fs need to be sorted in order from upstream to downstream.
	sortedFs, err := TopologicalSort(fs)
	if err != nil {
		return nil, err
	}
	evaluations := make([]*featureproto.Evaluation, 0, len(fs))
	archivedIDs := make([]string, 0)
	for _, f := range sortedFs {
		feature := &Feature{Feature: f}
		if feature.Archived {
			// To keep response size small, the feature flags archived long time ago are excluded.
			if !feature.IsArchivedBeforeLastThirtyDays() {
				archivedIDs = append(archivedIDs, f.Id)
			}
			continue
		}
		var segmentUsers []*featureproto.SegmentUser
		for _, id := range feature.ListSegmentIDs() {
			segmentUsers = append(segmentUsers, mapSegmentUsers[id]...)
		}
		reason, variation, err := feature.assignUser(user, segmentUsers, flagVariations)
		if err != nil {
			return nil, err
		}
		// VariationId is used to check if prerequisite flag's result is what user expects it to be.
		flagVariations[f.Id] = variation.Id
		// When the tag is set in the request,
		// it will return only the evaluations of flags that match the tag configured on the dashboard.
		// When empty, it will return all the evaluations of the flags in the environment.
		if targetTag != "" && !tagExist(f.Tags, targetTag) {
			continue
		}
		// FIXME: Remove the next line when the Variation
		// no longer is being used
		// For security reasons, it removes the variation description
		variation.Description = ""
		evaluationID := EvaluationID(f.Id, f.Version, user.Id)
		evaluation := &featureproto.Evaluation{
			Id:             evaluationID,
			FeatureId:      f.Id,
			FeatureVersion: f.Version,
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
func TopologicalSort(features []*featureproto.Feature) ([]*featureproto.Feature, error) {
	marks := map[string]Mark{}
	mapFeatures := map[string]*featureproto.Feature{}
	for _, f := range features {
		marks[f.Id] = unvisited
		mapFeatures[f.Id] = f
	}
	var sortedFeatures []*featureproto.Feature
	var sort func(f *featureproto.Feature) error
	sort = func(f *featureproto.Feature) error {
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
				return errFeatureNotFound
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
GetPrerequisiteDownwards gets the features specified as prerequisite by the targetFeatures.
*/
func GetPrerequisiteDownwards(
	targetFeatures, allFeatures []*featureproto.Feature,
) ([]*featureproto.Feature, error) {
	allFeaturesMap := make(map[string]*featureproto.Feature, len(allFeatures))
	for _, f := range allFeatures {
		allFeaturesMap[f.Id] = f
	}
	prerequisites := make(map[string]*featureproto.Feature)
	// depth first search
	queue := append([]*featureproto.Feature{}, targetFeatures...)
	for len(queue) > 0 {
		f := queue[0]
		for _, p := range f.Prerequisites {
			preFeature, ok := allFeaturesMap[p.FeatureId]
			if !ok {
				return nil, errFeatureNotFound
			}
			prerequisites[preFeature.Id] = preFeature
			queue = append(queue, preFeature)
		}
		queue = queue[1:]
	}
	return getPrerequisiteResult(targetFeatures, prerequisites), nil
}

/*
GetPrerequisiteUpwards gets the features that have the specified targetFeatures as the prerequisite.
*/
func GetPrerequisiteUpwards( // nolint:unused
	targetFeatures, featuresHavePrerequisite []*featureproto.Feature,
) []*featureproto.Feature {
	upwardsFeatures := make(map[string]*featureproto.Feature)
	// depth first search
	queue := append([]*featureproto.Feature{}, targetFeatures...)
	for len(queue) > 0 {
		f := queue[0]
		for _, newTarget := range featuresHavePrerequisite {
			for _, p := range newTarget.Prerequisites {
				if p.FeatureId == f.Id {
					if _, ok := upwardsFeatures[newTarget.Id]; ok {
						continue
					}
					upwardsFeatures[newTarget.Id] = newTarget
					queue = append(queue, newTarget)
				}
			}
		}
		queue = queue[1:]
	}
	return getPrerequisiteResult(targetFeatures, upwardsFeatures)
}

func getPrerequisiteResult(
	targetFeatures []*featureproto.Feature,
	featuresDependencies map[string]*featureproto.Feature,
) []*featureproto.Feature {
	if len(featuresDependencies) == 0 {
		return targetFeatures
	}
	targetFeaturesMap := make(map[string]*featureproto.Feature, len(targetFeatures))
	for _, f := range targetFeatures {
		targetFeaturesMap[f.Id] = f
	}
	merged := mapMerge(targetFeaturesMap, featuresDependencies)
	result := make([]*featureproto.Feature, 0, len(merged))
	for _, v := range merged {
		result = append(result, v)
	}
	return result
}

func mapMerge(m1, m2 map[string]*featureproto.Feature) map[string]*featureproto.Feature {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func getFeaturesHavePrerequisite(
	fs []*featureproto.Feature,
) []*featureproto.Feature {
	featuresHavePrerequisite := make(map[string]*featureproto.Feature)
	for _, f := range fs {
		if len(f.Prerequisites) == 0 {
			continue
		}
		if _, ok := featuresHavePrerequisite[f.Id]; ok {
			continue
		}
		featuresHavePrerequisite[f.Id] = f
	}
	result := make([]*featureproto.Feature, 0, len(featuresHavePrerequisite))
	for _, v := range featuresHavePrerequisite {
		result = append(result, v)
	}
	return result
}
