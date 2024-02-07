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
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type segmentEvaluator struct {
}

func (e *segmentEvaluator) Evaluate(segmentIDs []string, userID string, segmentUsers []*featureproto.SegmentUser) bool {
	return e.findSegmentUser(segmentIDs, userID, featureproto.SegmentUser_INCLUDED, segmentUsers)
}

func (e *segmentEvaluator) findSegmentUser(
	segmentIDs []string,
	userID string,
	state featureproto.SegmentUser_State,
	segmentUsers []*featureproto.SegmentUser,
) bool {
	for _, segmentID := range segmentIDs {
		if !e.containsSegmentUser(segmentID, userID, state, segmentUsers) {
			return false
		}
	}
	return true
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
