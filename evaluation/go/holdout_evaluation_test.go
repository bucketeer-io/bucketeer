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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestAssignUserHoldout(t *testing.T) {
	t.Parallel()
	user := &userproto.User{Id: "user-1"}
	bucket := (&bucketeer{}).bucket(holdoutHashInput("env", user.Id, "seed"))
	start := float64(int(bucket * 100))
	if start > 90 {
		start -= 10
	}
	outside := start + 50
	if outside >= 100 {
		outside -= 100
	}

	held := targetingGroup("hg", start, 10, "flag", "hold")
	miss := targetingGroup("hg", outside, 10, "flag", "hold")
	// Membership is memoized per evaluation and handed to assignUser. The
	// memo is single-use state, so each subtest builds its own.
	membership := func(group *ftproto.HoldoutGroup) *holdoutMembership {
		return newHoldoutMembership([]*ftproto.HoldoutGroup{group}, user, nil)
	}

	t.Run("disabled flag still serves the holdout variation", func(t *testing.T) {
		t.Parallel()
		flag := flagWithVariations("flag", false, "off", "hold", "on")
		reason, variation, err := NewEvaluator().assignUser(flag, user, nil, nil, nil, membership(held))
		require.NoError(t, err)
		assert.Equal(t, ftproto.Reason_HOLDOUT, reason.Type)
		assert.Equal(t, "hg", reason.HoldoutGroupId)
		assert.Equal(t, "hold", variation.Id)
	})
	t.Run("unsatisfied prerequisite does not override the holdout", func(t *testing.T) {
		t.Parallel()
		flag := flagWithVariations("flag", true, "off", "hold", "on")
		flag.Prerequisites = []*ftproto.Prerequisite{{FeatureId: "parent", VariationId: "expected"}}
		reason, variation, err := NewEvaluator().assignUser(
			flag, user, nil, nil, map[string]string{"parent": "other"}, membership(held),
		)
		require.NoError(t, err)
		assert.Equal(t, ftproto.Reason_HOLDOUT, reason.Type)
		assert.Equal(t, "hold", variation.Id)
	})
	t.Run("user outside the arc keeps the off variation", func(t *testing.T) {
		t.Parallel()
		flag := flagWithVariations("flag", false, "off", "hold", "on")
		reason, variation, err := NewEvaluator().assignUser(flag, user, nil, nil, nil, membership(miss))
		require.NoError(t, err)
		assert.Equal(t, ftproto.Reason_OFF_VARIATION, reason.Type)
		assert.Empty(t, reason.HoldoutGroupId)
		assert.Equal(t, "off", variation.Id)
	})
	t.Run("flag the group does not target is unchanged", func(t *testing.T) {
		t.Parallel()
		flag := flagWithVariations("other", true, "off", "hold", "on")
		withGroup, variation, err := NewEvaluator().assignUser(flag, user, nil, nil, nil, membership(held))
		require.NoError(t, err)
		without, baseline, err := NewEvaluator().assignUser(flag, user, nil, nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, without.Type, withGroup.Type)
		assert.Equal(t, baseline.Id, variation.Id)
	})
}

func TestEvaluateFeaturesHoldoutPrerequisite(t *testing.T) {
	t.Parallel()
	user := &userproto.User{Id: "user-1"}
	bucket := (&bucketeer{}).bucket(holdoutHashInput("env", user.Id, "seed"))
	start := float64(int(bucket * 100))
	if start > 90 {
		start -= 10
	}
	parent := flagWithVariations("parent", true, "off", "hold", "on")
	parent.DefaultStrategy.FixedStrategy.Variation = "on"
	child := flagWithVariations("child", true, "off", "hold", "on")
	child.Prerequisites = []*ftproto.Prerequisite{{FeatureId: "parent", VariationId: "hold"}}
	child.DefaultStrategy.FixedStrategy.Variation = "on"
	group := targetingGroup("hg", start, 10, "parent", "hold")

	evaluations, err := NewEvaluator().EvaluateFeatures(
		[]*ftproto.Feature{child, parent}, user, nil, nil, "", group,
	)
	require.NoError(t, err)
	byID := map[string]*ftproto.Evaluation{}
	for _, evaluation := range evaluations.Evaluations {
		byID[evaluation.FeatureId] = evaluation
	}
	require.Equal(t, "hold", byID["parent"].VariationId)
	assert.Equal(t, ftproto.Reason_HOLDOUT, byID["parent"].Reason.Type)
	assert.Equal(t, "hg", byID["parent"].Reason.HoldoutGroupId)
	require.Equal(t, "on", byID["child"].VariationId)
	assert.Equal(t, ftproto.Reason_DEFAULT, byID["child"].Reason.Type)
}

func TestEvaluateFeaturesByEvaluatedAtHoldoutForcesFullUpdate(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	user := &userproto.User{Id: "user-1"}
	bucket := (&bucketeer{}).bucket(holdoutHashInput("env", user.Id, "seed"))
	start := float64(int(bucket * 100))
	if start > 90 {
		start -= 10
	}
	evaluatedAt := now - 100
	stale := flagWithVariations("stale", true, "off", "hold", "on")
	stale.UpdatedAt = now - 1000
	fresh := flagWithVariations("fresh", true, "off", "hold", "on")
	fresh.UpdatedAt = now - 10
	live := func(updatedAt int64) *ftproto.HoldoutGroup {
		group := targetingGroup("hg", start, 10, "stale", "hold")
		group.UpdatedAt = updatedAt
		return group
	}

	// defaultSecondsForAdjustment is 10. The cursor is evaluatedAt-10.
	tests := []struct {
		name        string
		group       *ftproto.HoldoutGroup
		force       bool
		featureIDs  []string
		staleReason ftproto.Reason_Type
	}{
		{
			name:        "holdout newer than the cursor",
			group:       live(now - 50),
			force:       true,
			featureIDs:  []string{"stale", "fresh"},
			staleReason: ftproto.Reason_HOLDOUT,
		},
		{
			name:        "same second as evaluatedAt",
			group:       live(evaluatedAt),
			force:       true,
			featureIDs:  []string{"stale", "fresh"},
			staleReason: ftproto.Reason_HOLDOUT,
		},
		{
			name:        "update inside the grace window",
			group:       live(evaluatedAt - 9),
			force:       true,
			featureIDs:  []string{"stale", "fresh"},
			staleReason: ftproto.Reason_HOLDOUT,
		},
		{
			name:       "update older than the grace window",
			group:      live(evaluatedAt - 10),
			force:      false,
			featureIDs: []string{"fresh"},
		},
		{
			name:        "serving stopped",
			group:       &ftproto.HoldoutGroup{UpdatedAt: now - 50},
			force:       true,
			featureIDs:  []string{"stale", "fresh"},
			staleReason: ftproto.Reason_DEFAULT,
		},
		{
			name:       "environment never stopped a group",
			group:      &ftproto.HoldoutGroup{},
			force:      false,
			featureIDs: []string{"fresh"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewEvaluator().EvaluateFeaturesByEvaluatedAt(
				[]*ftproto.Feature{stale, fresh},
				user,
				nil,
				nil,
				"previous",
				evaluatedAt,
				false,
				"",
				tt.group,
			)
			require.NoError(t, err)
			assert.Equal(t, tt.force, got.ForceUpdate)
			ids := make([]string, 0, len(got.Evaluations))
			for _, evaluation := range got.Evaluations {
				ids = append(ids, evaluation.FeatureId)
				if evaluation.FeatureId != "stale" {
					continue
				}
				assert.Equal(t, tt.staleReason, evaluation.Reason.Type)
				if tt.staleReason == ftproto.Reason_HOLDOUT {
					assert.Equal(t, "hold", evaluation.VariationId)
				} else {
					assert.Equal(t, "on", evaluation.VariationId)
				}
			}
			assert.Equal(t, tt.featureIDs, ids)
		})
	}
}

func TestUnavailableHoldoutForcesFullUpdateWhenAnotherFlagChanged(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	user := &userproto.User{Id: "user-1"}
	evaluatedAt := now - 100
	target := flagWithVariations("target", true, "off", "hold", "on")
	target.UpdatedAt = now - 1000
	fresh := flagWithVariations("fresh", true, "off", "hold", "on")
	fresh.UpdatedAt = now - 10

	got, err := NewEvaluator().EvaluateFeaturesByEvaluatedAt(
		[]*ftproto.Feature{target, fresh},
		user, nil, nil,
		"previous",
		evaluatedAt,
		false,
		"",
		UnavailableHoldoutSnapshot(),
	)
	require.NoError(t, err)
	assert.True(t, got.ForceUpdate)
	byID := map[string]*ftproto.Evaluation{}
	for _, evaluation := range got.Evaluations {
		byID[evaluation.FeatureId] = evaluation
	}
	require.Contains(t, byID, "target")
	require.Contains(t, byID, "fresh")
	assert.Equal(t, ftproto.Reason_DEFAULT, byID["target"].Reason.Type)
	assert.Equal(t, "on", byID["target"].VariationId)
}

// A holdout-cache outage issues evaluations without holdout assignments. When
// the cache recovers, the live group's updated_at can predate the client's
// cursor, so the staleness check does not fire. The ID issued during the
// outage carries a marker and the next evaluation must be a full one that
// restores HOLDOUT for the target, instead of a diff with only the fresh flag.
func TestRecoveryAfterUnavailableHoldoutForcesFullUpdate(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	user := &userproto.User{Id: "user-1"}
	target := flagWithVariations("target", true, "off", "hold", "on")
	target.UpdatedAt = now - 1000
	fresh := flagWithVariations("fresh", true, "off", "hold", "on")
	fresh.UpdatedAt = now - 10
	features := []*ftproto.Feature{target, fresh}
	group := targetingGroup("group", 0, 100, "target", "hold")
	group.UpdatedAt = now - 5000 // started long before the client's cursor
	evaluatedAt := now - 100

	// During the outage the ID is marked and the member gets the normal variation.
	outageID := UserEvaluationsID(user.Id, user.Data, features, UnavailableHoldoutSnapshot())
	assert.True(t, IssuedWithoutHoldout(outageID))
	assert.Equal(t, outageID, UserEvaluationsID(user.Id, user.Data, features, UnavailableHoldoutSnapshot()),
		"the marked ID must be stable so polls during the outage still short-circuit")
	assert.NotEqual(t, outageID, UserEvaluationsID(user.Id, user.Data, features))
	assert.False(t, IssuedWithoutHoldout(UserEvaluationsID(user.Id, user.Data, features, group)))
	assert.False(t, IssuedWithoutHoldout(UserEvaluationsID(user.Id, user.Data, features, ClearanceSnapshot(now))))

	// Recovery: the live group is old, only "fresh" passes the timestamp diff.
	assert.False(t, holdoutStale([]*ftproto.HoldoutGroup{group}, evaluatedAt-10))
	got, err := NewEvaluator().EvaluateFeaturesByEvaluatedAt(
		features, user, nil, nil, outageID, evaluatedAt, false, "", group,
	)
	require.NoError(t, err)
	assert.True(t, got.ForceUpdate)
	byID := map[string]*ftproto.Evaluation{}
	for _, evaluation := range got.Evaluations {
		byID[evaluation.FeatureId] = evaluation
	}
	require.Contains(t, byID, "target")
	require.Contains(t, byID, "fresh")
	assert.Equal(t, ftproto.Reason_HOLDOUT, byID["target"].Reason.Type)
	assert.Equal(t, "hold", byID["target"].VariationId)

	// An unmarked previous ID with the same inputs still takes the diff path.
	diff, err := NewEvaluator().EvaluateFeaturesByEvaluatedAt(
		features, user, nil, nil, "previous", evaluatedAt, false, "", group,
	)
	require.NoError(t, err)
	assert.False(t, diff.ForceUpdate)
	require.Len(t, diff.Evaluations, 1)
	assert.Equal(t, "fresh", diff.Evaluations[0].FeatureId)
}

func flagWithVariations(id string, enabled bool, off, hold, on string) *ftproto.Feature {
	return &ftproto.Feature{
		Id:           id,
		Enabled:      enabled,
		OffVariation: off,
		Variations: []*ftproto.Variation{
			{Id: off, Value: off},
			{Id: hold, Value: hold},
			{Id: on, Value: on},
		},
		DefaultStrategy: &ftproto.Strategy{
			Type:          ftproto.Strategy_FIXED,
			FixedStrategy: &ftproto.FixedStrategy{Variation: on},
		},
	}
}

func targetingGroup(id string, bucketStart, percentage float64, featureID, variationID string) *ftproto.HoldoutGroup {
	return &ftproto.HoldoutGroup{
		Id:             id,
		EnvironmentId:  "env",
		AssignmentType: ftproto.HoldoutGroup_HASH,
		SamplingSeed:   "seed",
		BucketStart:    bucketStart,
		Percentage:     percentage,
		Features: []*ftproto.HoldoutGroupFeature{{
			FeatureId:          featureID,
			HoldoutVariationId: variationID,
		}},
	}
}
