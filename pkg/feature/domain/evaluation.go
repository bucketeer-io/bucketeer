// Copyright 2022 The Bucketeer Authors.
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

func EvaluateFeatures(
	fs []*featureproto.Feature,
	user *userproto.User,
	mapSegmentUsers map[string][]*featureproto.SegmentUser,
	targetTag string,
) (*featureproto.UserEvaluations, error) {
	flagVariations := map[string]string{}
	// fs need to be sorted in order from upstream to downstream.
	sortedFs, err := TopologicalSort(fs)
	if err != nil {
		return nil, err
	}
	evaluations := make([]*featureproto.Evaluation, 0, len(fs))
	for _, f := range sortedFs {
		feature := &Feature{Feature: f}
		segmentUsers := []*featureproto.SegmentUser{}
		for _, id := range feature.ListSegmentIDs() {
			segmentUsers = append(segmentUsers, mapSegmentUsers[id]...)
		}
		reason, variation, err := feature.assignUser(user, segmentUsers, flagVariations)
		if err != nil {
			return nil, err
		}
		// VariationId is used to check if prerequisite flag's result is what user expects it to be.
		flagVariations[f.Id] = variation.Id

		// We need to filter evaluations because we fetch all features in the environment namespace.
		if exist := tagExist(f.Tags, targetTag); !exist {
			continue
		}
		// FIXME: Remove the next two lines when the Variation
		// no longer is being used
		// For security reasons, it removes the variation name and description
		variation.Name = ""
		variation.Description = ""
		evaluationID := EvaluationID(f.Id, f.Version, user.Id)
		evaluation := &featureproto.Evaluation{
			Id:             evaluationID,
			FeatureId:      f.Id,
			FeatureVersion: f.Version,
			UserId:         user.Id,
			VariationId:    variation.Id,
			VariationValue: variation.Value,
			Variation:      variation, // deprecated
			Reason:         reason,
		}
		evaluations = append(evaluations, evaluation)
	}
	id := UserEvaluationsID(user.Id, user.Data, fs)
	userEvaluations := NewUserEvaluations(id, evaluations)
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
