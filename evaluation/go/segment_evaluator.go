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
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

type segmentEvaluator struct {
}

// Evaluate reports whether the user belongs to ANY of the given segments (OR).
func (e *segmentEvaluator) Evaluate(
	segmentIDs []string,
	user *userproto.User,
	segments map[string]*featureproto.Segment,
	segmentUsers []*featureproto.SegmentUser,
) (bool, error) {
	for _, segmentID := range segmentIDs {
		inSegment, err := e.isUserInSegment(segmentID, user, segments[segmentID], segmentUsers)
		if err != nil {
			return false, err
		}
		if inSegment {
			return true, nil
		}
	}
	return false, nil
}

// isUserInSegment reports whether the user belongs to the segment:
// the user is in the included-user list OR matches any of the segment rules.
func (e *segmentEvaluator) isUserInSegment(
	segmentID string,
	user *userproto.User,
	segment *featureproto.Segment,
	segmentUsers []*featureproto.SegmentUser,
) (bool, error) {
	// 1. Explicit include list — existing behavior, unchanged.
	if e.containsSegmentUser(segmentID, user.Id, featureproto.SegmentUser_INCLUDED, segmentUsers) {
		return true, nil
	}
	// 2. Rules. Segment.Rules is []*featureproto.Rule — the same type as Feature.Rules —
	// so the flag-side rule evaluator (OR across rules, AND across clauses, attribute
	// resolution, all operators) runs on it as is.
	// segmentUsers/segments/flagVariations are nil: validation rejects SEGMENT and
	// FEATURE_FLAG operators inside segment rules, and nil makes any such clause
	// fail closed.
	if segment == nil || len(segment.Rules) == 0 {
		return false, nil
	}
	ruleEval := &ruleEvaluator{} // local value: avoids the embedding cycle
	matchedRule, err := ruleEval.Evaluate(segment.Rules, user, nil, nil, nil)
	if err != nil {
		return false, err
	}
	return matchedRule != nil, nil
}

func (e *segmentEvaluator) containsSegmentUser(
	segmentID, userID string,
	state featureproto.SegmentUser_State,
	segmentUsers []*featureproto.SegmentUser,
) bool {
	for _, user := range segmentUsers {
		if user.SegmentId != segmentID {
			continue
		}
		if user.UserId != userID {
			continue
		}
		if user.State != state {
			continue
		}
		return true
	}
	return false
}
