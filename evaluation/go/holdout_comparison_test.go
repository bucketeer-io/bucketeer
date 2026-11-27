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
	"testing"

	"github.com/stretchr/testify/assert"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

// The comparison range is the arc of the size of the holdout arc that starts
// where the holdout arc ends. The expected answer is computed here on the
// integer grid, independently of holdoutArcContains: the distance of the
// user's position from the end of the holdout arc, going around the circle.
func TestIsInHoldoutComparisonRange(t *testing.T) {
	t.Parallel()
	const users = 20000
	patterns := []struct {
		desc        string
		bucketStart float64
		percentage  float64
	}{
		{desc: "both arcs contiguous", bucketStart: 10, percentage: 20},
		{desc: "holdout arc ends at the top of the circle", bucketStart: 95, percentage: 5},
		{desc: "holdout arc wraps", bucketStart: 97, percentage: 5},
		{desc: "comparison range wraps", bucketStart: 60, percentage: 33.333},
		{desc: "holdout arc starts near the top and both wrap past it", bucketStart: 99.999, percentage: 0.5},
		{desc: "50 percent partitions the circle", bucketStart: 0, percentage: 50},
		{desc: "50 percent wrapped partitions the circle", bucketStart: 75, percentage: 50},
		{desc: "smallest percentage", bucketStart: 42.123, percentage: 0.001},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			group := hashGroup(p.bucketStart, p.percentage)
			startUnits := int64(math.Round(p.bucketStart * 1000))
			lengthUnits := int64(math.Round(p.percentage * 1000))
			var heldOut, comparison int
			for i := 0; i < users; i++ {
				userID := fmt.Sprintf("user-%d", i)
				bucket := (&bucketeer{}).bucket(holdoutHashInput(group.EnvironmentId, userID, group.SamplingSeed))
				pos := int64(math.Round(bucket * float64(holdoutPercentScale)))
				if pos >= holdoutPercentScale {
					pos = 0
				}
				fromHoldoutEnd := ((pos-startUnits-lengthUnits)%holdoutPercentScale + holdoutPercentScale) %
					holdoutPercentScale
				expected := fromHoldoutEnd < lengthUnits

				actual := IsInHoldoutComparisonRange(group, userID)
				assert.Equal(t, expected, actual, "user %s at position %d", userID, pos)
				isHeldOut := IsHeldOut(group, &userproto.User{Id: userID}, nil)
				assert.False(t, isHeldOut && actual, "user %s is in both arcs", userID)
				if isHeldOut {
					heldOut++
				}
				if actual {
					comparison++
				}
			}
			if p.percentage == 50 {
				assert.Equal(t, users, heldOut+comparison, "the two arcs cover the circle")
			}
			// The two arcs have the same size, so about the same number of users.
			expectedUsers := float64(users) * p.percentage / 100
			assert.InDelta(t, expectedUsers, float64(comparison), 4*math.Sqrt(expectedUsers)+1)
			assert.InDelta(t, expectedUsers, float64(heldOut), 4*math.Sqrt(expectedUsers)+1)
		})
	}
}

func TestIsInHoldoutComparisonRangeNeverMatches(t *testing.T) {
	t.Parallel()
	// With 50 percent every user is in one of the two arcs, so a false answer
	// below comes from the guard and not from the hash.
	inComparison := ""
	for i := 0; inComparison == ""; i++ {
		userID := fmt.Sprintf("user-%d", i)
		if IsInHoldoutComparisonRange(hashGroup(0, 50), userID) {
			inComparison = userID
		}
	}
	patterns := []struct {
		desc   string
		group  *ftproto.HoldoutGroup
		userID string
	}{
		{desc: "nil group", group: nil, userID: inComparison},
		{desc: "empty user ID", group: hashGroup(0, 50), userID: ""},
		{
			desc: "segment group",
			group: &ftproto.HoldoutGroup{
				EnvironmentId:  "env",
				AssignmentType: ftproto.HoldoutGroup_SEGMENT,
				SegmentId:      "seg-1",
				SamplingSeed:   "seed",
				Percentage:     50,
			},
			userID: inComparison,
		},
		{
			desc: "hash group without a seed",
			group: func() *ftproto.HoldoutGroup {
				g := hashGroup(0, 50)
				g.SamplingSeed = ""
				return g
			}(),
			userID: inComparison,
		},
		{desc: "hash group without a percentage", group: hashGroup(0, 0), userID: inComparison},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.False(t, IsInHoldoutComparisonRange(p.group, p.userID))
		})
	}
}
