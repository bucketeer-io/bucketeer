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
	"math"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

// holdoutPercentScale is the DECIMAL(6,3) grid on [0, 100): one unit is
// 0.001 of a percent point. Both endpoints and the user's bucket are rounded
// onto it so Go and TypeScript agree at the boundary.
const holdoutPercentScale int64 = 100000

// IsHeldOut reports whether user is a member of group.
// HASH buckets "holdout-{environmentID}-{userID}-{seed}" with the murmur3
// bucketer and tests the arc [bucket_start, bucket_start+percentage)/100
// on [0, 1), wrapping at 1. SEGMENT checks the included user list only.
// A nil group or user is not held out. The comparison range is not a
// membership test: users there take the normal evaluation path.
func IsHeldOut(
	group *ftproto.HoldoutGroup,
	user *userproto.User,
	segmentUsers []*ftproto.SegmentUser,
) bool {
	if group == nil || user == nil {
		return false
	}
	switch group.AssignmentType {
	case ftproto.HoldoutGroup_SEGMENT:
		return (&segmentEvaluator{}).containsSegmentUser(
			group.SegmentId,
			user.Id,
			ftproto.SegmentUser_INCLUDED,
			segmentUsers,
		)
	case ftproto.HoldoutGroup_HASH:
		if group.SamplingSeed == "" || group.Percentage <= 0 {
			return false
		}
		bucket := (&bucketeer{}).bucket(holdoutHashInput(group.EnvironmentId, user.Id, group.SamplingSeed))
		return holdoutArcContains(bucket, group.BucketStart, group.Percentage)
	default:
		return false
	}
}

func holdoutHashInput(environmentID, userID, seed string) string {
	return fmt.Sprintf("holdout-%s-%s-%s", environmentID, userID, seed)
}

// holdoutArcContains reports whether bucket, in [0, 1), lies in
// [bucketStart, bucketStart+percentage)/100, wrapping past 1.
// The comparison is made on the DECIMAL(6,3) grid so a value that lands on
// the exclusive end stays outside the arc.
func holdoutArcContains(bucket, bucketStart, percentage float64) bool {
	start := int64(math.Round(bucketStart * 1000))
	end := int64(math.Round((bucketStart + percentage) * 1000))
	pos := int64(math.Round(bucket * float64(holdoutPercentScale)))
	if pos >= holdoutPercentScale {
		pos = 0
	}
	if end <= holdoutPercentScale {
		return pos >= start && pos < end
	}
	return pos >= start || pos < end-holdoutPercentScale
}

// HoldoutSegmentIDs returns the segment IDs of SEGMENT-mode groups.
// HASH groups and a nil list contribute nothing, so the no-holdout path
// does not allocate.
func HoldoutSegmentIDs(groups []*ftproto.HoldoutGroup) []string {
	var ids []string
	for _, group := range groups {
		if group == nil || group.AssignmentType != ftproto.HoldoutGroup_SEGMENT || group.SegmentId == "" {
			continue
		}
		if contains(group.SegmentId, ids) {
			continue
		}
		ids = append(ids, group.SegmentId)
	}
	return ids
}

func holdoutVariationID(group *ftproto.HoldoutGroup, featureID string) (string, bool) {
	if group == nil {
		return "", false
	}
	for _, feature := range group.Features {
		if feature == nil || feature.RemovedAt != 0 || feature.FeatureId != featureID {
			continue
		}
		return feature.HoldoutVariationId, true
	}
	return "", false
}

// ServingHoldout returns group when it is a live serving snapshot.
// An empty id, a deleted or archived group, and any status other than
// RUNNING or ANALYSIS is not serving.
func ServingHoldout(group *ftproto.HoldoutGroup) *ftproto.HoldoutGroup {
	if group == nil || group.Id == "" || group.Deleted || group.Archived {
		return nil
	}
	if group.Status != ftproto.HoldoutGroup_RUNNING && group.Status != ftproto.HoldoutGroup_ANALYSIS {
		return nil
	}
	return group
}

// EvaluationHoldout is the snapshot the evaluator should see.
// A serving group is returned as-is. An empty snapshot whose updated_at is
// set is the clearance signal written when serving stopped: it assigns
// nothing, and it forces a full re-evaluation while the client's cursor is
// behind that revision. UpdatedAt 0 means the environment has never stopped
// a group and is not a signal.
func EvaluationHoldout(group *ftproto.HoldoutGroup) *ftproto.HoldoutGroup {
	if serving := ServingHoldout(group); serving != nil {
		return serving
	}
	if group != nil && group.Id == "" && group.UpdatedAt > 0 {
		return group
	}
	return nil
}

// holdoutStale reports whether the client's cursor is behind a holdout
// revision. The caller passes evaluatedAt minus the feature-diff grace
// window. A live group and a clearance snapshot (empty id, updated_at of
// the stop) both count. UpdatedAt 0 is not a revision.
func holdoutStale(groups []*ftproto.HoldoutGroup, evaluatedAt int64) bool {
	for _, group := range groups {
		if group == nil || group.UpdatedAt == 0 {
			continue
		}
		if evaluatedAt < group.UpdatedAt {
			return true
		}
	}
	return false
}
