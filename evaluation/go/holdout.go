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

// holdoutMembership resolves membership lazily and at most once per group
// for one evaluation. Membership is flag-independent (RFC 0048 §4.2.1), so
// the first flag a group targets pays one hash (HASH) or one scan of the
// group's own segment (SEGMENT); later flags reuse the answer, and an
// evaluation that touches no targeted flag never hashes at all.
type holdoutMembership struct {
	groups          []*ftproto.HoldoutGroup
	user            *userproto.User
	mapSegmentUsers map[string][]*ftproto.SegmentUser
	// 0 unresolved, 1 member, 2 not a member; indexed like groups.
	state []uint8
}

func newHoldoutMembership(
	groups []*ftproto.HoldoutGroup,
	user *userproto.User,
	mapSegmentUsers map[string][]*ftproto.SegmentUser,
) *holdoutMembership {
	if len(groups) == 0 {
		return nil
	}
	return &holdoutMembership{
		groups:          groups,
		user:            user,
		mapSegmentUsers: mapSegmentUsers,
		state:           make([]uint8, len(groups)),
	}
}

// variationFor returns the group and holdout variation to serve for
// featureID, or false when no group the user belongs to targets it.
// A nil receiver means no groups were passed.
func (m *holdoutMembership) variationFor(featureID string) (*ftproto.HoldoutGroup, string, bool) {
	if m == nil {
		return nil, "", false
	}
	for i, group := range m.groups {
		variationID, targeted := holdoutVariationID(group, featureID)
		if !targeted {
			continue
		}
		if m.isMember(i) {
			return group, variationID, true
		}
		// At most one serving group targets a flag. A miss stays on the normal path.
		break
	}
	return nil, "", false
}

// hasSnapshot reports whether evaluation was given a holdout snapshot.
// A nil receiver means no snapshot: a clearance and an unavailable snapshot
// both count, because neither one should fall through to a normal variation.
func (m *holdoutMembership) hasSnapshot() bool {
	return m != nil && len(m.groups) > 0
}

// targets reports whether any group targets featureID, including when the
// user is not a member. A nil receiver means no groups were passed.
func (m *holdoutMembership) targets(featureID string) bool {
	if m == nil {
		return false
	}
	for _, group := range m.groups {
		if _, targeted := holdoutVariationID(group, featureID); targeted {
			return true
		}
	}
	return false
}

func (m *holdoutMembership) isMember(i int) bool {
	if m.state[i] == 0 {
		group := m.groups[i]
		var segmentUsers []*ftproto.SegmentUser
		if group.AssignmentType == ftproto.HoldoutGroup_SEGMENT {
			segmentUsers = m.mapSegmentUsers[group.SegmentId]
		}
		if IsHeldOut(group, m.user, segmentUsers) {
			m.state[i] = 1
		} else {
			m.state[i] = 2
		}
	}
	return m.state[i] == 1
}

// ClearanceSnapshot is the empty snapshot the cacher stores when no group is
// serving. It assigns nothing; stoppedAt is the updated_at of the newest
// stopped group, so a client whose cursor is behind it is forced to a full
// re-evaluation and drops cached HOLDOUT assignments (RFC 0048 §4.3.2).
// stoppedAt 0 means the environment has never stopped a group.
func ClearanceSnapshot(stoppedAt int64) *ftproto.HoldoutGroup {
	return &ftproto.HoldoutGroup{UpdatedAt: stoppedAt}
}

// unavailableHoldoutUpdatedAt marks a clearance used when the snapshot
// cannot be read. It is earlier than any real stop time, so it does not
// move a client's cursor. It only keeps the snapshot non-nil so a disabled
// flag with no off variation stays omitted instead of falling through to a
// normal variation.
const unavailableHoldoutUpdatedAt int64 = 1

// UnavailableHoldoutSnapshot is the clearance returned when the holdout
// cache cannot be read. A miss stays nil: that is "no group".
func UnavailableHoldoutSnapshot() *ftproto.HoldoutGroup {
	return ClearanceSnapshot(unavailableHoldoutUpdatedAt)
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

// HoldoutTargetFeatureIDs is the set of current target feature IDs.
// An empty id is a clearance snapshot and has no targets. Removed rows are skipped.
func HoldoutTargetFeatureIDs(group *ftproto.HoldoutGroup) map[string]struct{} {
	if group == nil || group.Id == "" {
		return nil
	}
	ids := make(map[string]struct{}, len(group.Features))
	for _, feature := range group.Features {
		if feature == nil || feature.FeatureId == "" || feature.RemovedAt != 0 {
			continue
		}
		ids[feature.FeatureId] = struct{}{}
	}
	if len(ids) == 0 {
		return nil
	}
	return ids
}

// HoldoutStale reports whether the client's cursor is behind a holdout
// revision. The caller passes evaluatedAt minus the feature-diff grace
// window. A live group and a clearance snapshot (empty id, updated_at of
// the stop) both count. UpdatedAt 0 is not a revision.
func HoldoutStale(groups []*ftproto.HoldoutGroup, evaluatedAt int64) bool {
	return holdoutStale(groups, evaluatedAt)
}

// FeaturesForServerEval drops disabled flags with no off variation that the
// snapshot no longer targets. Those flags stay in the feature cache until
// the next feature refresh, and without this filter the evaluator would
// assign them a normal variation after the group stops or the snapshot is a
// clearance. A nil snapshot leaves the slice unchanged.
// removed is the IDs a partial evaluation must report so clients delete them.
func FeaturesForServerEval(
	features []*ftproto.Feature,
	group *ftproto.HoldoutGroup,
) (kept []*ftproto.Feature, removed []string) {
	if group == nil {
		return features, nil
	}
	targets := HoldoutTargetFeatureIDs(group)
	kept = make([]*ftproto.Feature, 0, len(features))
	for _, feature := range features {
		if feature != nil && !feature.Enabled && feature.OffVariation == "" {
			if _, targeted := targets[feature.Id]; !targeted {
				if feature.Id != "" {
					removed = append(removed, feature.Id)
				}
				continue
			}
		}
		kept = append(kept, feature)
	}
	return kept, removed
}

// FeaturesForLocalEval splits flags for GetFeatureFlags. Local evaluation
// does not receive holdout metadata yet, so any snapshot — a serving group
// or a clearance — withholds every disabled flag with no off variation,
// including a former target of a group that is still running. A nil snapshot
// leaves the slice unchanged.
func FeaturesForLocalEval(
	features []*ftproto.Feature,
	group *ftproto.HoldoutGroup,
) (visible, hidden []*ftproto.Feature) {
	if group == nil {
		return features, nil
	}
	visible = make([]*ftproto.Feature, 0, len(features))
	for _, feature := range features {
		if feature != nil && !feature.Enabled && feature.OffVariation == "" {
			hidden = append(hidden, feature)
			continue
		}
		visible = append(visible, feature)
	}
	return visible, hidden
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
