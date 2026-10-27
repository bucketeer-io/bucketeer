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

package evaluation

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"time"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type UserEvaluations struct {
	*ftproto.UserEvaluations
}

func NewUserEvaluations(
	id string,
	evaluations []*ftproto.Evaluation,
	archivedFeaturesIds []string,
	forceUpdate bool,
) *UserEvaluations {
	now := time.Now().Unix()
	return &UserEvaluations{&ftproto.UserEvaluations{
		Id:                 id,
		Evaluations:        evaluations,
		CreatedAt:          now,
		ArchivedFeatureIds: archivedFeaturesIds,
		ForceUpdate:        forceUpdate,
	}}
}

func UserEvaluationsID(
	userID string,
	userMetadata map[string]string,
	features []*ftproto.Feature,
	groups ...*ftproto.HoldoutGroup,
) string {
	sort.SliceStable(features, func(i, j int) bool {
		return features[i].Id < features[j].Id
	})
	// FNV-1a 64-bit is sufficient for change detection:
	// - Fast enough for SDK polling frequency
	// - 64-bit provides negligible collision risk for feature flag counts
	// - No cryptographic requirements
	h := fnv.New64a()
	h.Write([]byte(userID)) // nolint:errcheck
	keys := sortMapKeys(userMetadata)
	for _, key := range keys {
		fmt.Fprintf(h, "%s:%s", key, userMetadata[key])
	}
	for _, feature := range features {
		fmt.Fprintf(h, "%s:%d", feature.Id, feature.UpdatedAt)
	}
	// Hash live groups after features so a nil group leaves the ID unchanged.
	// UpdatedAt is whole seconds, so the current target rows are part of the
	// ID: feature, variation, when the row was added, and the row id. Two
	// edits in the same second then cannot share an ID.
	for _, group := range sortedHoldoutGroups(groups) {
		fmt.Fprintf(h, "%s:%d", group.Id, group.UpdatedAt) // nolint:errcheck
		for _, target := range sortedHoldoutTargets(group) {
			fmt.Fprintf( // nolint:errcheck
				h, ":%s:%s:%d:%s",
				target.FeatureId, target.HoldoutVariationId, target.AddedAt, target.Id,
			)
		}
	}
	id := strconv.FormatUint(h.Sum64(), 10)
	// An ID issued without a readable holdout snapshot is marked so the next
	// evaluation after the snapshot is back runs in full. See IssuedWithoutHoldout.
	if holdoutUnavailable(groups) {
		return unavailableUEIDPrefix + id
	}
	return id
}

func sortedHoldoutTargets(group *ftproto.HoldoutGroup) []*ftproto.HoldoutGroupFeature {
	if group == nil || len(group.Features) == 0 {
		return nil
	}
	targets := make([]*ftproto.HoldoutGroupFeature, 0, len(group.Features))
	for _, feature := range group.Features {
		if feature == nil || feature.FeatureId == "" || feature.RemovedAt != 0 {
			continue
		}
		targets = append(targets, feature)
	}
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].FeatureId != targets[j].FeatureId {
			return targets[i].FeatureId < targets[j].FeatureId
		}
		return targets[i].Id < targets[j].Id
	})
	return targets
}

func sortedHoldoutGroups(groups []*ftproto.HoldoutGroup) []*ftproto.HoldoutGroup {
	if len(groups) == 0 {
		return nil
	}
	live := make([]*ftproto.HoldoutGroup, 0, len(groups))
	for _, group := range groups {
		if group != nil && group.Id != "" {
			live = append(live, group)
		}
	}
	sort.Slice(live, func(i, j int) bool {
		return live[i].Id < live[j].Id
	})
	return live
}

func sortMapKeys(data map[string]string) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func GenerateFeaturesID(features []*ftproto.Feature) string {
	sort.SliceStable(features, func(i, j int) bool {
		return features[i].Id < features[j].Id
	})
	// FNV-1a 64-bit is sufficient for change detection:
	// - Fast enough for SDK polling frequency
	// - 64-bit provides negligible collision risk for feature flag counts
	// - No cryptographic requirements
	h := fnv.New64a()
	for _, feature := range features {
		fmt.Fprintf(h, "%s:%d", feature.Id, feature.UpdatedAt)
	}
	return strconv.FormatUint(h.Sum64(), 10)
}
