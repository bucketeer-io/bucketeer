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

func BenchmarkEvaluateFeatures(b *testing.B) {
	features := make([]*ftproto.Feature, 0, 32)
	for i := 0; i < 32; i++ {
		feature := flagWithVariations(fmt.Sprintf("flag-%d", i), true, "off", "hold", "on")
		features = append(features, feature)
	}
	features[0].Id = "targeted"
	user := &userproto.User{Id: "user-1"}
	evaluator := NewEvaluator()
	group := targetingGroup("hg", 0, 5, "targeted", "hold")
	untargeted := targetingGroup("hg", 0, 5, "missing", "hold")

	b.Run("no holdout", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := evaluator.EvaluateFeatures(features, user, nil, nil, ""); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("group does not target", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := evaluator.EvaluateFeatures(features, user, nil, nil, "", untargeted); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("group targets one flag", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := evaluator.EvaluateFeatures(features, user, nil, nil, "", group); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestDisabledHoldoutTargetWithoutOffVariation(t *testing.T) {
	t.Parallel()
	evaluator := NewEvaluator()
	disabled := &ftproto.Feature{
		Id:      "flag",
		Enabled: false,
		Variations: []*ftproto.Variation{
			{Id: "hold", Value: "hold"},
			{Id: "on", Value: "on"},
		},
		DefaultStrategy: &ftproto.Strategy{
			Type:          ftproto.Strategy_FIXED,
			FixedStrategy: &ftproto.FixedStrategy{Variation: "on"},
		},
	}
	sibling := flagWithVariations("sibling", true, "off", "hold", "on")
	user := &userproto.User{Id: "user-1"}

	member, err := evaluator.EvaluateFeatures(
		[]*ftproto.Feature{disabled, sibling},
		user, nil, nil, "",
		targetingGroup("hg", 0, 100, "flag", "hold"),
	)
	require.NoError(t, err)
	require.Len(t, member.Evaluations, 2)
	byID := map[string]*ftproto.Evaluation{}
	for _, evaluation := range member.Evaluations {
		byID[evaluation.FeatureId] = evaluation
	}
	assert.Equal(t, ftproto.Reason_HOLDOUT, byID["flag"].Reason.Type)
	assert.Equal(t, "hold", byID["flag"].VariationId)
	assert.Equal(t, ftproto.Reason_DEFAULT, byID["sibling"].Reason.Type)

	nonMember, err := evaluator.EvaluateFeatures(
		[]*ftproto.Feature{disabled, sibling},
		user, nil, nil, "",
		targetingGroup("hg", 0, 0, "flag", "hold"),
	)
	require.NoError(t, err)
	require.Len(t, nonMember.Evaluations, 1)
	assert.Equal(t, "sibling", nonMember.Evaluations[0].FeatureId)
	assert.Equal(t, []string{"flag"}, nonMember.ArchivedFeatureIds)

	now := time.Now().Unix()
	disabled.UpdatedAt = now
	partial, err := evaluator.EvaluateFeaturesByEvaluatedAt(
		[]*ftproto.Feature{disabled, sibling},
		user, nil, nil,
		"prev", now-5, false, "",
		targetingGroup("hg", 0, 0, "flag", "hold"),
	)
	require.NoError(t, err)
	assert.False(t, partial.ForceUpdate)
	assert.Equal(t, []string{"flag"}, partial.ArchivedFeatureIds)
	for _, evaluation := range partial.Evaluations {
		assert.NotEqual(t, "flag", evaluation.FeatureId)
	}
}

func TestOmittedHoldoutTargetDoesNotAbortDependents(t *testing.T) {
	t.Parallel()
	evaluator := NewEvaluator()
	user := &userproto.User{Id: "user-1"}
	target := &ftproto.Feature{
		Id:      "flag",
		Enabled: false,
		Variations: []*ftproto.Variation{
			{Id: "hold", Value: "hold"},
			{Id: "on", Value: "on"},
		},
		DefaultStrategy: &ftproto.Strategy{
			Type:          ftproto.Strategy_FIXED,
			FixedStrategy: &ftproto.FixedStrategy{Variation: "on"},
		},
	}
	dependent := flagWithVariations("dependent", true, "off", "hold", "on")
	dependent.Prerequisites = []*ftproto.Prerequisite{{FeatureId: "flag", VariationId: "on"}}
	unrelated := flagWithVariations("sibling", true, "off", "hold", "on")

	// Percentage 0 holds nobody, so the disabled target is omitted.
	evaluations, err := evaluator.EvaluateFeatures(
		[]*ftproto.Feature{dependent, target, unrelated},
		user, nil, nil, "",
		targetingGroup("hg", 0, 0, "flag", "hold"),
	)
	require.NoError(t, err)
	byID := map[string]*ftproto.Evaluation{}
	for _, evaluation := range evaluations.Evaluations {
		byID[evaluation.FeatureId] = evaluation
	}
	assert.NotContains(t, byID, "flag")
	assert.Equal(t, []string{"flag"}, evaluations.ArchivedFeatureIds)
	require.Contains(t, byID, "dependent")
	assert.Equal(t, ftproto.Reason_PREREQUISITE, byID["dependent"].Reason.Type)
	assert.Equal(t, "off", byID["dependent"].VariationId)
	require.Contains(t, byID, "sibling")
	assert.Equal(t, ftproto.Reason_DEFAULT, byID["sibling"].Reason.Type)

	// A clearance is a snapshot with no targets. The disabled flag is still
	// omitted, and a request with no snapshot keeps the historical default.
	cleared, err := evaluator.EvaluateFeatures(
		[]*ftproto.Feature{target, unrelated},
		user, nil, nil, "",
		ClearanceSnapshot(10),
	)
	require.NoError(t, err)
	assert.Equal(t, []string{"flag"}, cleared.ArchivedFeatureIds)
	require.Len(t, cleared.Evaluations, 1)
	assert.Equal(t, "sibling", cleared.Evaluations[0].FeatureId)

	historical, err := evaluator.EvaluateFeatures(
		[]*ftproto.Feature{target},
		user, nil, nil, "",
	)
	require.NoError(t, err)
	require.Len(t, historical.Evaluations, 1)
	assert.Equal(t, ftproto.Reason_DEFAULT, historical.Evaluations[0].Reason.Type)
	assert.Empty(t, historical.ArchivedFeatureIds)
}

func TestFeaturesForServerEvalKeepsDependencyNodes(t *testing.T) {
	t.Parallel()
	user := &userproto.User{Id: "user-1"}
	disabled := &ftproto.Feature{
		Id:      "flag",
		Enabled: false,
		Variations: []*ftproto.Variation{
			{Id: "hold", Value: "hold"},
			{Id: "on", Value: "on"},
		},
		DefaultStrategy: &ftproto.Strategy{
			Type:          ftproto.Strategy_FIXED,
			FixedStrategy: &ftproto.FixedStrategy{Variation: "on"},
		},
	}
	prerequisite := flagWithVariations("dependent", true, "off", "hold", "on")
	prerequisite.Prerequisites = []*ftproto.Prerequisite{{FeatureId: "flag", VariationId: "on"}}
	rule := flagWithVariations("rule-dependent", true, "off", "hold", "on")
	rule.Rules = []*ftproto.Rule{{
		Id: "rule-1",
		Clauses: []*ftproto.Clause{{
			Id:        "clause-1",
			Attribute: "flag",
			Operator:  ftproto.Clause_FEATURE_FLAG,
			Values:    []string{"on"},
		}},
	}}
	orphan := &ftproto.Feature{Id: "orphan", Enabled: false}
	features := []*ftproto.Feature{prerequisite, rule, disabled, orphan}

	kept, removed := FeaturesForServerEval(features, ClearanceSnapshot(10))
	require.Equal(t, []string{"orphan"}, removed)
	keptIDs := make([]string, 0, len(kept))
	for _, feature := range kept {
		keptIDs = append(keptIDs, feature.Id)
	}
	assert.ElementsMatch(t, []string{"dependent", "rule-dependent", "flag"}, keptIDs)

	evaluations, err := NewEvaluator().EvaluateFeatures(kept, user, nil, nil, "", ClearanceSnapshot(10))
	require.NoError(t, err)
	byID := map[string]*ftproto.Evaluation{}
	for _, evaluation := range evaluations.Evaluations {
		byID[evaluation.FeatureId] = evaluation
	}
	assert.NotContains(t, byID, "flag")
	assert.Contains(t, evaluations.ArchivedFeatureIds, "flag")
	require.Contains(t, byID, "dependent")
	assert.Equal(t, ftproto.Reason_PREREQUISITE, byID["dependent"].Reason.Type)
	assert.Equal(t, "off", byID["dependent"].VariationId)
	require.Contains(t, byID, "rule-dependent")
	assert.Equal(t, ftproto.Reason_DEFAULT, byID["rule-dependent"].Reason.Type)
}

func TestFeaturesForEvalFiltersClearedDisabledTarget(t *testing.T) {
	t.Parallel()
	disabled := &ftproto.Feature{Id: "flag", Enabled: false}
	enabled := &ftproto.Feature{Id: "on", Enabled: true, OffVariation: "on"}
	features := []*ftproto.Feature{disabled, enabled}

	kept, removed := FeaturesForServerEval(features, nil)
	require.Len(t, kept, 2)
	assert.Empty(t, removed)

	kept, removed = FeaturesForServerEval(features, ClearanceSnapshot(10))
	require.Len(t, kept, 1)
	assert.Equal(t, "on", kept[0].Id)
	assert.Equal(t, []string{"flag"}, removed)

	visible, hidden := FeaturesForLocalEval(features, &ftproto.HoldoutGroup{
		Id:       "hg",
		Features: []*ftproto.HoldoutGroupFeature{{FeatureId: "flag"}},
	})
	require.Len(t, visible, 1)
	assert.Equal(t, "on", visible[0].Id)
	require.Len(t, hidden, 1)
	assert.Equal(t, "flag", hidden[0].Id)

	visible, hidden = FeaturesForLocalEval(features, ClearanceSnapshot(10))
	require.Len(t, hidden, 1)
	assert.Equal(t, "flag", hidden[0].Id)

	// A former target of a group that is still running is hidden too.
	visible, hidden = FeaturesForLocalEval(features, &ftproto.HoldoutGroup{Id: "hg"})
	require.Len(t, visible, 1)
	assert.Equal(t, "on", visible[0].Id)
	require.Len(t, hidden, 1)
	assert.Equal(t, "flag", hidden[0].Id)

	visible, hidden = FeaturesForLocalEval(features, nil)
	require.Len(t, visible, 2)
	assert.Empty(t, hidden)
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
