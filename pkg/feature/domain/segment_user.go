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

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type SegmentUser struct {
	*featureproto.SegmentUser
}

func NewSegmentUser(segmentID string, userID string, state featureproto.SegmentUser_State, deleted bool) *SegmentUser {
	id := SegmentUserID(segmentID, userID, state)
	return &SegmentUser{
		SegmentUser: &featureproto.SegmentUser{
			Id:        id,
			SegmentId: segmentID,
			UserId:    userID,
			State:     state,
			Deleted:   deleted,
		},
	}
}

func SegmentUserID(segmentID string, userID string, state featureproto.SegmentUser_State) string {
	return fmt.Sprintf("%s:%s:%v", segmentID, userID, state)
}
