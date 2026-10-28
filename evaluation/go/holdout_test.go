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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestHoldoutArcContains(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		bucket      float64
		bucketStart float64
		percentage  float64
		expected    bool
	}{
		{desc: "inside a contiguous arc", bucket: 0.15, bucketStart: 10, percentage: 20, expected: true},
		{desc: "start is inclusive", bucket: 0.10, bucketStart: 10, percentage: 20, expected: true},
		{desc: "end is exclusive", bucket: 0.30, bucketStart: 10, percentage: 20, expected: false},
		{desc: "below the arc", bucket: 0.09, bucketStart: 10, percentage: 20, expected: false},
		{desc: "wrap high side", bucket: 0.97, bucketStart: 97, percentage: 5, expected: true},
		{desc: "wrap low side", bucket: 0.01, bucketStart: 97, percentage: 5, expected: true},
		{desc: "wrap end is exclusive", bucket: 0.02, bucketStart: 97, percentage: 5, expected: false},
		{desc: "comparison band is not the holdout", bucket: 0.05, bucketStart: 97, percentage: 5, expected: false},
		{desc: "50 percent lower half", bucket: 0, bucketStart: 0, percentage: 50, expected: true},
		{desc: "50 percent end is exclusive", bucket: 0.5, bucketStart: 0, percentage: 50, expected: false},
		{desc: "50 percent upper half", bucket: 0.75, bucketStart: 0, percentage: 50, expected: false},
		{
			desc:        "50 percent wrap partitions the circle",
			bucket:      0.1,
			bucketStart: 75,
			percentage:  50,
			expected:    true,
		},
		{
			desc:        "50 percent wrap comparison half",
			bucket:      0.4,
			bucketStart: 75,
			percentage:  50,
			expected:    false,
		},
	}
	for _, p := range patterns {
		p := p
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.expected, holdoutArcContains(p.bucket, p.bucketStart, p.percentage))
		})
	}
	t.Run("three decimal percentage rounds onto the stored grid", func(t *testing.T) {
		t.Parallel()
		const percentage = 33.333
		end := float64(int64(math.Round(percentage*1000))) / float64(holdoutPercentScale)
		assert.True(t, holdoutArcContains(end-1/float64(holdoutPercentScale), 0, percentage))
		assert.False(t, holdoutArcContains(end, 0, percentage))
		// 33.3334 is stored as 33.333 on DECIMAL(6,3).
		assert.False(t, holdoutArcContains(end, 0, 33.3334))
	})
}

func TestIsHeldOut(t *testing.T) {
	t.Parallel()
	user := &userproto.User{Id: "user-1"}
	hashInput := holdoutHashInput("env", user.Id, "seed")
	bucket := (&bucketeer{}).bucket(hashInput)
	start := math.Floor(bucket * 100)
	if start > 90 {
		start -= 10
	}

	t.Run("hash user falls inside the arc", func(t *testing.T) {
		t.Parallel()
		group := hashGroup(start, 10)
		assert.True(t, IsHeldOut(group, user, nil))
	})
	t.Run("hash user falls outside the arc", func(t *testing.T) {
		t.Parallel()
		group := hashGroup(math.Mod(start+50, 100), 10)
		assert.False(t, IsHeldOut(group, user, nil))
	})
	t.Run("hash without a seed is not held out", func(t *testing.T) {
		t.Parallel()
		group := hashGroup(0, 50)
		group.SamplingSeed = ""
		assert.False(t, IsHeldOut(group, user, nil))
	})
	t.Run("segment included user", func(t *testing.T) {
		t.Parallel()
		group := &ftproto.HoldoutGroup{
			AssignmentType: ftproto.HoldoutGroup_SEGMENT,
			SegmentId:      "seg-1",
		}
		users := []*ftproto.SegmentUser{{
			SegmentId: "seg-1",
			UserId:    user.Id,
			State:     ftproto.SegmentUser_INCLUDED,
		}}
		assert.True(t, IsHeldOut(group, user, users))
	})
	t.Run("segment excluded user", func(t *testing.T) {
		t.Parallel()
		group := &ftproto.HoldoutGroup{
			AssignmentType: ftproto.HoldoutGroup_SEGMENT,
			SegmentId:      "seg-1",
		}
		users := []*ftproto.SegmentUser{{
			SegmentId: "seg-1",
			UserId:    user.Id,
			State:     ftproto.SegmentUser_EXCLUDED,
		}}
		assert.False(t, IsHeldOut(group, user, users))
	})
	t.Run("segment absent user", func(t *testing.T) {
		t.Parallel()
		group := &ftproto.HoldoutGroup{
			AssignmentType: ftproto.HoldoutGroup_SEGMENT,
			SegmentId:      "seg-1",
		}
		assert.False(t, IsHeldOut(group, user, nil))
	})
	t.Run("nil group", func(t *testing.T) {
		t.Parallel()
		assert.False(t, IsHeldOut(nil, user, nil))
	})
	t.Run("nil user", func(t *testing.T) {
		t.Parallel()
		assert.False(t, IsHeldOut(hashGroup(0, 50), nil, nil))
	})
}

func TestHoldoutSegmentIDs(t *testing.T) {
	t.Parallel()
	assert.Nil(t, HoldoutSegmentIDs(nil))
	assert.Nil(t, HoldoutSegmentIDs([]*ftproto.HoldoutGroup{hashGroup(0, 5)}))
	ids := HoldoutSegmentIDs([]*ftproto.HoldoutGroup{
		{AssignmentType: ftproto.HoldoutGroup_SEGMENT, SegmentId: "seg-1"},
		{AssignmentType: ftproto.HoldoutGroup_SEGMENT, SegmentId: "seg-1"},
		{AssignmentType: ftproto.HoldoutGroup_SEGMENT, SegmentId: "seg-2"},
		{AssignmentType: ftproto.HoldoutGroup_HASH, SegmentId: "ignored"},
	})
	assert.Equal(t, []string{"seg-1", "seg-2"}, ids)
}

func TestHoldoutStale(t *testing.T) {
	t.Parallel()
	group := &ftproto.HoldoutGroup{Id: "hg", UpdatedAt: 2000}
	assert.True(t, holdoutStale([]*ftproto.HoldoutGroup{group}, 1999))
	assert.False(t, holdoutStale([]*ftproto.HoldoutGroup{group}, 2000))
	// Empty id with a timestamp is the snapshot written when serving stopped.
	assert.True(t, holdoutStale([]*ftproto.HoldoutGroup{{UpdatedAt: 2000}}, 1000))
	assert.False(t, holdoutStale([]*ftproto.HoldoutGroup{{UpdatedAt: 2000}}, 2000))
	assert.False(t, holdoutStale([]*ftproto.HoldoutGroup{{}}, 1000))
	assert.False(t, holdoutStale(nil, 1000))
}

func hashGroup(bucketStart, percentage float64) *ftproto.HoldoutGroup {
	return &ftproto.HoldoutGroup{
		Id:             "hg",
		EnvironmentId:  "env",
		AssignmentType: ftproto.HoldoutGroup_HASH,
		SamplingSeed:   "seed",
		BucketStart:    bucketStart,
		Percentage:     percentage,
	}
}

func TestHoldoutHashInput(t *testing.T) {
	t.Parallel()
	require.Equal(t, "holdout-env-user-seed", holdoutHashInput("env", "user", "seed"))
}

func TestServingHoldout(t *testing.T) {
	t.Parallel()
	live := &ftproto.HoldoutGroup{Id: "hg", Status: ftproto.HoldoutGroup_RUNNING}
	analysis := &ftproto.HoldoutGroup{Id: "hg", Status: ftproto.HoldoutGroup_ANALYSIS}
	assert.Same(t, live, ServingHoldout(live))
	assert.Same(t, analysis, ServingHoldout(analysis))
	assert.Nil(t, ServingHoldout(nil))
	assert.Nil(t, ServingHoldout(&ftproto.HoldoutGroup{Status: ftproto.HoldoutGroup_RUNNING}))
	assert.Nil(t, ServingHoldout(&ftproto.HoldoutGroup{Id: "hg", Deleted: true, Status: ftproto.HoldoutGroup_RUNNING}))
	assert.Nil(t, ServingHoldout(&ftproto.HoldoutGroup{Id: "hg", Archived: true, Status: ftproto.HoldoutGroup_ANALYSIS}))
	assert.Nil(t, ServingHoldout(&ftproto.HoldoutGroup{Id: "hg", Status: ftproto.HoldoutGroup_WAITING}))
	assert.Nil(t, ServingHoldout(&ftproto.HoldoutGroup{Id: "hg", Status: ftproto.HoldoutGroup_STOPPED}))
}

func TestEvaluationHoldout(t *testing.T) {
	t.Parallel()
	live := &ftproto.HoldoutGroup{Id: "hg", Status: ftproto.HoldoutGroup_RUNNING}
	cleared := &ftproto.HoldoutGroup{UpdatedAt: 50}
	assert.Same(t, live, EvaluationHoldout(live))
	assert.Same(t, cleared, EvaluationHoldout(cleared))
	assert.Nil(t, EvaluationHoldout(nil))
	neverStopped := &ftproto.HoldoutGroup{}
	assert.Same(t, neverStopped, EvaluationHoldout(neverStopped))
	assert.Nil(t, EvaluationHoldout(&ftproto.HoldoutGroup{
		Id:        "hg",
		Status:    ftproto.HoldoutGroup_STOPPED,
		UpdatedAt: 50,
	}))
}

func TestHoldoutMembership(t *testing.T) {
	t.Parallel()
	user := &userproto.User{Id: "user-1"}
	segmentGroup := func(id, segmentID string, features ...*ftproto.HoldoutGroupFeature) *ftproto.HoldoutGroup {
		return &ftproto.HoldoutGroup{
			Id:             id,
			EnvironmentId:  "env",
			AssignmentType: ftproto.HoldoutGroup_SEGMENT,
			SegmentId:      segmentID,
			Features:       features,
		}
	}
	target := &ftproto.HoldoutGroupFeature{FeatureId: "flag", HoldoutVariationId: "hold"}
	removed := &ftproto.HoldoutGroupFeature{FeatureId: "flag", HoldoutVariationId: "hold", RemovedAt: 1}
	included := map[string][]*ftproto.SegmentUser{
		"seg": {{SegmentId: "seg", UserId: user.Id, State: ftproto.SegmentUser_INCLUDED}},
		// Another flag's segment also lists the user; it must not be consulted.
		"other": {{SegmentId: "seg", UserId: user.Id, State: ftproto.SegmentUser_INCLUDED}},
	}

	t.Run("nil receiver and empty input", func(t *testing.T) {
		t.Parallel()
		var none *holdoutMembership
		_, _, ok := none.variationFor("flag")
		assert.False(t, ok)
		assert.Nil(t, newHoldoutMembership(nil, user, included))
	})
	t.Run("segment member of the group's own segment", func(t *testing.T) {
		t.Parallel()
		group := segmentGroup("hg", "seg", target)
		got, variationID, ok := newHoldoutMembership([]*ftproto.HoldoutGroup{group}, user, included).variationFor("flag")
		require.True(t, ok)
		assert.Same(t, group, got)
		assert.Equal(t, "hold", variationID)
	})
	t.Run("only the group's segment is scanned", func(t *testing.T) {
		t.Parallel()
		// The user row is filed under key "other", not under the group's segment.
		group := segmentGroup("hg", "seg", target)
		users := map[string][]*ftproto.SegmentUser{"other": included["other"]}
		_, _, ok := newHoldoutMembership([]*ftproto.HoldoutGroup{group}, user, users).variationFor("flag")
		assert.False(t, ok)
	})
	t.Run("untargeted flag does not resolve membership", func(t *testing.T) {
		t.Parallel()
		group := segmentGroup("hg", "seg", target)
		m := newHoldoutMembership([]*ftproto.HoldoutGroup{group}, user, included)
		_, _, ok := m.variationFor("other-flag")
		assert.False(t, ok)
		assert.Equal(t, uint8(0), m.state[0], "no lookup for a flag the group does not target")
		_, _, ok = m.variationFor("flag")
		assert.True(t, ok)
		assert.Equal(t, uint8(1), m.state[0])
	})
	t.Run("membership is resolved once and reused", func(t *testing.T) {
		t.Parallel()
		group := segmentGroup("hg", "seg", target, &ftproto.HoldoutGroupFeature{FeatureId: "flag2", HoldoutVariationId: "h2"})
		// Deliberately give the lookup a map we mutate after the first call: a
		// second resolution would observe the removal.
		users := map[string][]*ftproto.SegmentUser{"seg": included["seg"]}
		m := newHoldoutMembership([]*ftproto.HoldoutGroup{group}, user, users)
		_, _, ok := m.variationFor("flag")
		require.True(t, ok)
		delete(users, "seg")
		_, variationID, ok := m.variationFor("flag2")
		assert.True(t, ok)
		assert.Equal(t, "h2", variationID)
	})
	t.Run("removed target and clearance snapshot are not targets", func(t *testing.T) {
		t.Parallel()
		_, _, ok := newHoldoutMembership([]*ftproto.HoldoutGroup{segmentGroup("hg", "seg", removed)}, user, included).variationFor("flag")
		assert.False(t, ok)
		_, _, ok = newHoldoutMembership([]*ftproto.HoldoutGroup{ClearanceSnapshot(42)}, user, included).variationFor("flag")
		assert.False(t, ok)
	})
	t.Run("a miss on the targeting group does not fall through to another group", func(t *testing.T) {
		t.Parallel()
		miss := segmentGroup("miss", "none", target)
		hit := segmentGroup("hit", "seg", target)
		_, _, ok := newHoldoutMembership([]*ftproto.HoldoutGroup{miss, hit}, user, included).variationFor("flag")
		assert.False(t, ok, "at most one serving group targets a flag")
	})
}

func TestClearanceSnapshot(t *testing.T) {
	t.Parallel()
	snapshot := ClearanceSnapshot(42)
	assert.Empty(t, snapshot.Id)
	assert.Empty(t, snapshot.Features)
	assert.Equal(t, int64(42), snapshot.UpdatedAt)
	// It is the signal EvaluationHoldout passes through, and it is stale for an older cursor.
	assert.Same(t, snapshot, EvaluationHoldout(snapshot))
	assert.True(t, holdoutStale([]*ftproto.HoldoutGroup{snapshot}, 41))
	assert.False(t, holdoutStale([]*ftproto.HoldoutGroup{snapshot}, 42))
	// A never-stopped environment is not a re-evaluation signal, but the
	// snapshot still reaches the evaluator so a disabled dependency is omitted.
	neverStopped := ClearanceSnapshot(0)
	assert.Same(t, neverStopped, EvaluationHoldout(neverStopped))
	assert.False(t, holdoutStale([]*ftproto.HoldoutGroup{neverStopped}, 0))
	assert.True(t, holdoutStale([]*ftproto.HoldoutGroup{UnavailableHoldoutSnapshot()}, 1_700_000_000))
}
