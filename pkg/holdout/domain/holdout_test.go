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

package domain

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestAssignHashCopiesEnvironmentSeed(t *testing.T) {
	restore := stubHashRandom("unused-new-seed", 12.5)
	defer restore()

	seed, bucket, err := AssignHash("env-seed")
	require.NoError(t, err)
	assert.Equal(t, "env-seed", seed)
	assert.Equal(t, 12.5, bucket)

	seed, bucket, err = AssignHash("")
	require.NoError(t, err)
	assert.Equal(t, "unused-new-seed", seed)
	assert.Equal(t, 12.5, bucket)
}

func TestNormalizePercentageRejectsNaN(t *testing.T) {
	t.Parallel()
	_, err := NormalizePercentage(math.NaN())
	assert.ErrorIs(t, err, ErrPercentageOutOfRange)
}

func TestRandomBucketStartInRange(t *testing.T) {
	t.Parallel()
	for i := 0; i < 20; i++ {
		bucket, err := defaultRandomBucketStart()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, bucket, 0.0)
		assert.Less(t, bucket, 100.0)
	}
}

func TestNewHoldoutGroupHash(t *testing.T) {
	restore := stubHashRandom("generated-seed", 3)
	defer restore()
	now := time.Now().Unix()

	group, err := NewHoldoutGroup(CreateParams{
		EnvironmentID:  "env",
		Name:           "Q1 holdout",
		AssignmentType: featureproto.HoldoutGroup_HASH,
		Percentage:     5.0004,
		StartAt:        now + 10,
		StopAt:         now + 100,
		CreatedBy:      "admin@example.com",
		Now:            now,
	})
	require.NoError(t, err)
	assert.Equal(t, featureproto.HoldoutGroup_WAITING, group.Status)
	assert.Equal(t, "generated-seed", group.SamplingSeed)
	assert.Equal(t, 3.0, group.BucketStart)
	assert.Equal(t, 5.0, group.Percentage)
	assert.Empty(t, group.ExperimentId)
	assert.Empty(t, group.AnalysisExperimentId)
	assert.Empty(t, group.GoalIds)

	_, err = NewHoldoutGroup(CreateParams{
		EnvironmentID:  "env",
		Name:           "too big",
		AssignmentType: featureproto.HoldoutGroup_HASH,
		Percentage:     50.001,
		StartAt:        now + 10,
		StopAt:         now + 100,
		Now:            now,
	})
	assert.ErrorIs(t, err, ErrPercentageOutOfRange)

	_, err = NewHoldoutGroup(CreateParams{
		EnvironmentID:  "env",
		Name:           "segment on hash",
		AssignmentType: featureproto.HoldoutGroup_HASH,
		Percentage:     5,
		SegmentID:      "seg",
		StartAt:        now + 10,
		StopAt:         now + 100,
		Now:            now,
	})
	assert.ErrorIs(t, err, ErrSegmentNotAllowed)
}

func TestNewHoldoutGroupSegment(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group, err := NewHoldoutGroup(CreateParams{
		EnvironmentID:   "env",
		Name:            "known cohort",
		AssignmentType:  featureproto.HoldoutGroup_SEGMENT,
		SegmentID:       "seg-1",
		StartAt:         now + 10,
		StopAt:          now + 100,
		StartAnalysisAt: now + 50,
		GoalIDs:         []string{"g1", "g1", "g2"},
		CreatedBy:       "admin@example.com",
		Now:             now,
	})
	require.NoError(t, err)
	assert.Equal(t, "seg-1", group.SegmentId)
	assert.Empty(t, group.SamplingSeed)
	assert.Equal(t, []string{"g1", "g2"}, group.GoalIds)

	_, err = NewHoldoutGroup(CreateParams{
		EnvironmentID:  "env",
		Name:           "missing segment",
		AssignmentType: featureproto.HoldoutGroup_SEGMENT,
		StartAt:        now + 10,
		StopAt:         now + 100,
		Now:            now,
	})
	assert.ErrorIs(t, err, ErrSegmentRequired)
}

func TestRejectClientAssignment(t *testing.T) {
	t.Parallel()
	assert.NoError(t, RejectClientAssignment(nil, nil))
	assert.ErrorIs(t, RejectClientAssignment(wrapperspb.Double(0), nil), ErrClientAssignmentRejected)
	assert.ErrorIs(t, RejectClientAssignment(nil, wrapperspb.String("")), ErrClientAssignmentRejected)
}

func TestPeriodsOverlap(t *testing.T) {
	t.Parallel()
	assert.False(t, PeriodsOverlap(0, 10, 10, 20))
	assert.True(t, PeriodsOverlap(0, 10, 9, 12))
	assert.True(t, PeriodsOverlap(0, 10, 0, 10))
	assert.False(t, PeriodsOverlap(20, 30, 0, 20))

	others := []GroupSnapshot{{
		ID: "other", Name: "other", Status: featureproto.HoldoutGroup_WAITING, StartAt: 100, StopAt: 200,
	}, {
		ID: "done", Name: "done", Status: featureproto.HoldoutGroup_STOPPED, StartAt: 0, StopAt: 50,
	}}
	conflict, ok := OverlappingGroup(150, 180, "self", others)
	assert.True(t, ok)
	assert.Equal(t, "other", conflict.ID)
	_, ok = OverlappingGroup(200, 250, "self", others)
	assert.False(t, ok)
	_, ok = OverlappingGroup(0, 50, "self", others)
	assert.False(t, ok)
	_, ok = OverlappingGroup(100, 200, "other", others)
	assert.False(t, ok)
}

func TestOtherServingGroup(t *testing.T) {
	t.Parallel()
	others := []GroupSnapshot{
		{ID: "waiting", Status: featureproto.HoldoutGroup_WAITING},
		{ID: "running", Name: "live", Status: featureproto.HoldoutGroup_RUNNING},
	}
	got, ok := OtherServingGroup("self", others)
	assert.True(t, ok)
	assert.Equal(t, "running", got.ID)
	_, ok = OtherServingGroup("running", others)
	assert.False(t, ok)

	locked, ok := ServingGroupForSegment("seg", []GroupSnapshot{{
		ID: "g", Name: "g", SegmentID: "seg", Status: featureproto.HoldoutGroup_ANALYSIS,
	}})
	assert.True(t, ok)
	assert.Equal(t, "g", locked.ID)
	_, ok = ServingGroupForSegment("seg", []GroupSnapshot{{
		ID: "g", SegmentID: "seg", Status: featureproto.HoldoutGroup_WAITING,
	}})
	assert.False(t, ok)
	locked, ok = ServingGroupForSegment("SEG", []GroupSnapshot{{
		ID: "g", Name: "g", SegmentID: "seg", Status: featureproto.HoldoutGroup_RUNNING,
	}})
	assert.True(t, ok)
	assert.Equal(t, "g", locked.ID)
}

func TestConflictingFeatureTarget(t *testing.T) {
	t.Parallel()
	targets := []FeatureTarget{{
		GroupID: "g1", GroupName: "one", FeatureID: "flag", Status: featureproto.HoldoutGroup_RUNNING,
	}, {
		GroupID: "g2", FeatureID: "flag", Status: featureproto.HoldoutGroup_STOPPED,
	}}
	got, ok := ConflictingFeatureTarget("flag", "self", targets)
	assert.True(t, ok)
	assert.Equal(t, "g1", got.GroupID)
	_, ok = ConflictingFeatureTarget("flag", "g1", targets)
	assert.False(t, ok)
	_, ok = ConflictingFeatureTarget("other", "self", targets)
	assert.False(t, ok)
	got, ok = ConflictingFeatureTarget("FLAG", "self", targets)
	assert.True(t, ok)
	assert.Equal(t, "g1", got.GroupID)
}

func TestStatusTransitions(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now+100, now+1000, now+500)

	require.NoError(t, group.Start(now))
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, group.Status)
	assert.Equal(t, now, group.StartAt)

	require.NoError(t, group.StartAnalysis(now+10))
	assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, group.Status)
	assert.Equal(t, now+10, group.AnalysisStartedAt)
	assert.ErrorIs(t, group.CanAddFeatures(), ErrAnalysisFrozen)
	assert.NoError(t, group.CanRemoveFeatures())

	require.NoError(t, group.StopManual(now+20))
	assert.Equal(t, featureproto.HoldoutGroup_FORCE_STOPPED, group.Status)
	assert.Equal(t, now+20, group.StoppedAt)
	assert.ErrorIs(t, group.Delete(now+20), ErrInvalidStatus)
	require.NoError(t, group.Archive(now+30))
	assert.True(t, group.Archived)
}

func TestStartScheduledKeepsSavedStart(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now-10, now+100, 0)
	require.NoError(t, group.StartScheduled(now))
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, group.Status)
	assert.Equal(t, now-10, group.StartAt)

	postponed := waitingGroup(now+100, now+1000, 0)
	assert.ErrorIs(t, postponed.StartScheduled(now), ErrBeforeStart)
	assert.Equal(t, featureproto.HoldoutGroup_WAITING, postponed.Status)
	assert.Equal(t, now+100, postponed.StartAt)

	expired := waitingGroup(now-200, now-100, 0)
	assert.ErrorIs(t, expired.StartScheduled(now), ErrAfterStop)
	assert.Equal(t, featureproto.HoldoutGroup_WAITING, expired.Status)
}

func TestStartAnalysisScheduledRequiresDueTime(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now-100, now+100, now-10)
	require.NoError(t, group.Start(now-50))
	require.NoError(t, group.StartAnalysisScheduled(now))
	assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, group.Status)
	assert.Equal(t, now, group.AnalysisStartedAt)

	cleared := waitingGroup(now-100, now+100, 0)
	require.NoError(t, cleared.Start(now-50))
	assert.ErrorIs(t, cleared.StartAnalysisScheduled(now), ErrBeforeAnalysis)
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, cleared.Status)

	postponed := waitingGroup(now-100, now+100, now+10)
	require.NoError(t, postponed.Start(now-50))
	assert.ErrorIs(t, postponed.StartAnalysisScheduled(now), ErrBeforeAnalysis)
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, postponed.Status)

	expired := waitingGroup(now-100, now+100, now-10)
	require.NoError(t, expired.Start(now-50))
	expired.StopAt = now
	assert.ErrorIs(t, expired.StartAnalysisScheduled(now), ErrAfterStop)
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, expired.Status)
}

func TestScheduledFinish(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now-10, now+100, 0)
	require.NoError(t, group.Start(now))
	assert.Equal(t, now-10, group.StartAt)
	assert.ErrorIs(t, group.FinishScheduled(now), ErrBeforeStop)
	group.StopAt = now
	require.NoError(t, group.FinishScheduled(now))
	assert.Equal(t, featureproto.HoldoutGroup_STOPPED, group.Status)
}

func TestScheduledFinishMissedWindow(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	// The batch never ran during [start_at, stop_at): the group is still WAITING after stop_at.
	group := waitingGroup(now-200, now-100, 0)
	assert.ErrorIs(t, group.Start(now), ErrStopAtBeforeNow)
	assert.Equal(t, featureproto.HoldoutGroup_WAITING, group.Status)

	assert.False(t, DueForScheduledStart(group.HoldoutGroup, time.Unix(now, 0)))
	assert.True(t, DueForScheduledStop(group.HoldoutGroup, time.Unix(now, 0)))
	require.NoError(t, group.FinishScheduled(now))
	assert.Equal(t, featureproto.HoldoutGroup_STOPPED, group.Status)
	assert.Equal(t, now, group.StoppedAt)

	// A WAITING group whose stop_at is still ahead must not be finished.
	pending := waitingGroup(now-10, now+100, 0)
	assert.ErrorIs(t, pending.FinishScheduled(now), ErrBeforeStop)
	assert.Equal(t, featureproto.HoldoutGroup_WAITING, pending.Status)
}

func TestStartAnalysisRejectsExpiredPeriod(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now-100, now+100, now)
	require.NoError(t, group.Start(now-50))

	// now == stop_at is already outside the half-open window.
	group.StopAt = now
	assert.ErrorIs(t, group.StartAnalysis(now), ErrAfterStop)
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, group.Status)
	assert.Zero(t, group.AnalysisStartedAt)

	group.StopAt = now + 1
	require.NoError(t, group.StartAnalysis(now))
	assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, group.Status)
	assert.Equal(t, now, group.AnalysisStartedAt)
}

func TestRejectedUpdateLeavesGroupUntouched(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now+100, now+1000, 0)
	group.Description = "before"
	group.GoalIds = []string{"g1"}
	before := group.Clone()

	// Name and description are valid, but the period is not: nothing may be applied.
	err := group.Update(UpdateParams{
		Name:        wrapperspb.String("renamed"),
		Description: wrapperspb.String("after"),
		StopAt:      wrapperspb.Int64(now + 50),
		GoalIDsSet:  true,
		GoalIDs:     []string{"g1", "g2"},
		UpdatedBy:   "admin",
		Now:         now,
	})
	assert.ErrorIs(t, err, ErrPeriodInvalid)
	assert.True(t, proto.Equal(before.HoldoutGroup, group.HoldoutGroup))

	// A manual start that fails validation must not clamp start_at.
	expired := waitingGroup(now+10, now+20, 0)
	assert.ErrorIs(t, expired.Start(now+20), ErrStopAtBeforeNow)
	assert.Equal(t, now+10, expired.StartAt)
}

func TestDeleteOnlyWhileWaiting(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now+10, now+100, 0)
	require.NoError(t, group.Delete(now))
	assert.True(t, group.Deleted)
}

func TestUpdateGoalAndScheduleRules(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	group := waitingGroup(now+100, now+1000, 0)
	group.GoalIds = []string{"g1"}

	err := group.Update(UpdateParams{
		GoalIDsSet: true,
		GoalIDs:    nil,
		UpdatedBy:  "admin",
		Now:        now,
	})
	require.NoError(t, err)
	assert.Empty(t, group.GoalIds)

	require.NoError(t, group.Start(now))
	group.GoalIds = []string{"g1"}
	err = group.Update(UpdateParams{
		StartAt:   wrapperspb.Int64(now + 1),
		UpdatedBy: "admin",
		Now:       now,
	})
	assert.ErrorIs(t, err, ErrStartAtImmutable)

	err = group.Update(UpdateParams{
		GoalIDsSet: true,
		GoalIDs:    []string{},
		UpdatedBy:  "admin",
		Now:        now,
	})
	assert.ErrorIs(t, err, ErrGoalsCannotBeRemoved)

	group.GoalIds = []string{"g1"}
	err = group.Update(UpdateParams{
		GoalIDsSet: true,
		GoalIDs:    []string{"g1", "g2"},
		StopAt:     wrapperspb.Int64(now + 2000),
		UpdatedBy:  "admin",
		Now:        now,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"g1", "g2"}, group.GoalIds)
	assert.Equal(t, now+2000, group.StopAt)
}

func TestDueHelpers(t *testing.T) {
	t.Parallel()
	now := time.Unix(1_000, 0)
	waiting := &featureproto.HoldoutGroup{
		Status: featureproto.HoldoutGroup_WAITING, StartAt: 1_000, StopAt: 2_000,
	}
	assert.True(t, DueForScheduledStart(waiting, now))
	waiting.StartAt = 1_001
	assert.False(t, DueForScheduledStart(waiting, now))

	running := &featureproto.HoldoutGroup{
		Status:          featureproto.HoldoutGroup_RUNNING,
		StartAnalysisAt: 1_000,
		StopAt:          3_000,
	}
	assert.True(t, DueForScheduledAnalysis(running, now))
	running.StopAt = 1_000
	assert.False(t, DueForScheduledAnalysis(running, now))
	assert.True(t, DueForScheduledStop(running, now))
}

func TestValidateFeatureInputs(t *testing.T) {
	t.Parallel()
	err := ValidateFeatureInputs([]*featureproto.HoldoutGroupFeatureInput{
		{FeatureId: "f", HoldoutVariationId: "v"},
		{FeatureId: "f", HoldoutVariationId: "v2"},
	})
	assert.ErrorIs(t, err, ErrDuplicateFeature)
	assert.True(t, VariationBelongsToFeature("v", []*featureproto.Variation{{Id: "v"}}))
	assert.False(t, VariationBelongsToFeature("missing", []*featureproto.Variation{{Id: "v"}}))
}

func waitingGroup(start, stop, analysis int64) *HoldoutGroup {
	return &HoldoutGroup{HoldoutGroup: &featureproto.HoldoutGroup{
		Id:              "g",
		Name:            "g",
		Status:          featureproto.HoldoutGroup_WAITING,
		StartAt:         start,
		StopAt:          stop,
		StartAnalysisAt: analysis,
		AssignmentType:  featureproto.HoldoutGroup_HASH,
		Percentage:      5,
		SamplingSeed:    "seed",
		BucketStart:     1,
	}}
}

func stubHashRandom(seed string, bucket float64) func() {
	prevSeed := newSamplingSeed
	prevBucket := randomBucketStart
	newSamplingSeed = func() (string, error) { return seed, nil }
	randomBucketStart = func() (float64, error) { return bucket, nil }
	return func() {
		newSamplingSeed = prevSeed
		randomBucketStart = prevBucket
	}
}
