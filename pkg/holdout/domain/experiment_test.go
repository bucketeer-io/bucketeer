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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const experimentTestNow int64 = 1_800_000_000

func measuredGroup(status featureproto.HoldoutGroup_Status) *HoldoutGroup {
	return &HoldoutGroup{HoldoutGroup: &featureproto.HoldoutGroup{
		Id:            "hg",
		EnvironmentId: "env",
		Name:          "holdout",
		Description:   "about",
		Status:        status,
		StartAt:       experimentTestNow - 1_000,
		StopAt:        experimentTestNow + 1_000,
		GoalIds:       []string{"goal-1", "goal-2"},
		CreatedBy:     "creator@example.com",
	}}
}

func TestNewInternalExperimentNeedsAServingGroupWithGoals(t *testing.T) {
	t.Parallel()
	statuses := []featureproto.HoldoutGroup_Status{
		featureproto.HoldoutGroup_WAITING,
		featureproto.HoldoutGroup_RUNNING,
		featureproto.HoldoutGroup_ANALYSIS,
		featureproto.HoldoutGroup_STOPPED,
		featureproto.HoldoutGroup_FORCE_STOPPED,
	}
	for _, status := range statuses {
		t.Run("no goals while "+status.String(), func(t *testing.T) {
			t.Parallel()
			group := measuredGroup(status)
			group.GoalIds = nil
			experiment, err := group.NewInternalExperiment(featureproto.HoldoutGroup_WAITING, experimentTestNow)
			require.NoError(t, err)
			assert.Nil(t, experiment)
			assert.Empty(t, group.ExperimentId)
			assert.Empty(t, group.AnalysisExperimentId)
		})
	}
	notServing := []featureproto.HoldoutGroup_Status{
		featureproto.HoldoutGroup_WAITING,
		featureproto.HoldoutGroup_STOPPED,
		featureproto.HoldoutGroup_FORCE_STOPPED,
	}
	for _, status := range notServing {
		t.Run("goals while "+status.String(), func(t *testing.T) {
			t.Parallel()
			group := measuredGroup(status)
			experiment, err := group.NewInternalExperiment(status, experimentTestNow)
			require.NoError(t, err)
			assert.Nil(t, experiment)
			assert.Empty(t, group.ExperimentId)
			assert.Empty(t, group.AnalysisExperimentId)
		})
	}
}

func TestNewInternalExperimentOfTheWholePeriod(t *testing.T) {
	t.Parallel()

	t.Run("a starting group measures from its start_at", func(t *testing.T) {
		t.Parallel()
		group := measuredGroup(featureproto.HoldoutGroup_RUNNING)
		experiment, err := group.NewInternalExperiment(featureproto.HoldoutGroup_WAITING, experimentTestNow)
		require.NoError(t, err)
		require.NotNil(t, experiment)

		assert.Equal(t, experiment.Id, group.ExperimentId)
		assert.Empty(t, group.AnalysisExperimentId)
		assert.NotEmpty(t, experiment.Id)
		assert.Equal(t, experimentproto.Experiment_HOLDOUT, experiment.Kind)
		assert.Equal(t, "hg", experiment.FeatureId)
		assert.Equal(t, InternalExperimentFeatureVersion, experiment.FeatureVersion)
		assert.Equal(t, int32(1), experiment.FeatureVersion)
		require.Len(t, experiment.Variations, 2)
		assert.Equal(t, "in_holdout", experiment.Variations[0].Id)
		assert.Equal(t, "not_in_holdout", experiment.Variations[1].Id)
		assert.Equal(t, "in_holdout", experiment.BaseVariationId)
		assert.Equal(t, []string{"goal-1", "goal-2"}, experiment.GoalIds)
		assert.Equal(t, group.StartAt, experiment.StartAt)
		assert.Equal(t, group.StopAt, experiment.StopAt)
		assert.Equal(t, int64(math.MaxInt64), experiment.StoppedAt)
		assert.Equal(t, experimentproto.Experiment_RUNNING, experiment.Status)
		assert.Equal(t, "holdout", experiment.Name)
		assert.Equal(t, "about", experiment.Description)
		assert.Equal(t, "creator@example.com", experiment.Maintainer)
		assert.Equal(t, experimentTestNow, experiment.CreatedAt)
		assert.Equal(t, experimentTestNow, experiment.UpdatedAt)
		assert.False(t, experiment.Archived)
		assert.False(t, experiment.Deleted)
		assert.False(t, experiment.Stopped)
	})

	t.Run("a group that was already running measures from now", func(t *testing.T) {
		t.Parallel()
		group := measuredGroup(featureproto.HoldoutGroup_RUNNING)
		experiment, err := group.NewInternalExperiment(featureproto.HoldoutGroup_RUNNING, experimentTestNow)
		require.NoError(t, err)
		require.NotNil(t, experiment)
		assert.Equal(t, experiment.Id, group.ExperimentId)
		assert.Equal(t, experimentTestNow, experiment.StartAt)
		assert.Equal(t, group.StopAt, experiment.StopAt)
	})

	t.Run("a group that has one gets no second one", func(t *testing.T) {
		t.Parallel()
		group := measuredGroup(featureproto.HoldoutGroup_RUNNING)
		group.ExperimentId = "existing"
		experiment, err := group.NewInternalExperiment(featureproto.HoldoutGroup_RUNNING, experimentTestNow)
		require.NoError(t, err)
		assert.Nil(t, experiment)
		assert.Equal(t, "existing", group.ExperimentId)
	})
}

func TestNewInternalExperimentOfTheAnalysisPeriod(t *testing.T) {
	t.Parallel()

	t.Run("a group entering ANALYSIS measures from analysis_started_at", func(t *testing.T) {
		t.Parallel()
		group := measuredGroup(featureproto.HoldoutGroup_ANALYSIS)
		group.ExperimentId = "whole"
		group.AnalysisStartedAt = experimentTestNow - 5
		experiment, err := group.NewInternalExperiment(featureproto.HoldoutGroup_RUNNING, experimentTestNow)
		require.NoError(t, err)
		require.NotNil(t, experiment)

		assert.Equal(t, "whole", group.ExperimentId)
		assert.Equal(t, experiment.Id, group.AnalysisExperimentId)
		assert.NotEqual(t, "whole", experiment.Id)
		assert.Equal(t, experimentproto.Experiment_HOLDOUT, experiment.Kind)
		assert.Equal(t, "hg", experiment.FeatureId)
		assert.Equal(t, int32(1), experiment.FeatureVersion)
		assert.Equal(t, "in_holdout", experiment.BaseVariationId)
		assert.Equal(t, group.AnalysisStartedAt, experiment.StartAt)
		assert.Equal(t, group.StopAt, experiment.StopAt)
		assert.Equal(t, experimentproto.Experiment_RUNNING, experiment.Status)
	})

	t.Run("first goals during ANALYSIS give one experiment that measures from now", func(t *testing.T) {
		t.Parallel()
		group := measuredGroup(featureproto.HoldoutGroup_ANALYSIS)
		group.AnalysisStartedAt = experimentTestNow - 500
		experiment, err := group.NewInternalExperiment(featureproto.HoldoutGroup_ANALYSIS, experimentTestNow)
		require.NoError(t, err)
		require.NotNil(t, experiment)
		assert.Empty(t, group.ExperimentId)
		assert.Equal(t, experiment.Id, group.AnalysisExperimentId)
		assert.Equal(t, experimentTestNow, experiment.StartAt)
		assert.Equal(t, group.StopAt, experiment.StopAt)

		again, err := group.NewInternalExperiment(featureproto.HoldoutGroup_ANALYSIS, experimentTestNow+1)
		require.NoError(t, err)
		assert.Nil(t, again)
		assert.Empty(t, group.ExperimentId)
	})
}

func TestMirrorToExperiment(t *testing.T) {
	t.Parallel()
	newExperiment := func(t *testing.T) (*HoldoutGroup, *experimentdomain.Experiment) {
		t.Helper()
		group := measuredGroup(featureproto.HoldoutGroup_RUNNING)
		experiment, err := group.NewInternalExperiment(featureproto.HoldoutGroup_WAITING, experimentTestNow)
		require.NoError(t, err)
		require.NotNil(t, experiment)
		return group, experiment
	}
	later := experimentTestNow + 60

	t.Run("an unchanged group changes nothing", func(t *testing.T) {
		t.Parallel()
		group, experiment := newExperiment(t)
		before := proto.Clone(experiment.Experiment)
		assert.False(t, group.MirrorToExperiment(experiment, later))
		assert.True(t, proto.Equal(before, experiment.Experiment))
	})

	t.Run("an edit of start_analysis_at changes nothing", func(t *testing.T) {
		t.Parallel()
		group, experiment := newExperiment(t)
		group.StartAnalysisAt = experimentTestNow + 500
		assert.False(t, group.MirrorToExperiment(experiment, later))
	})

	t.Run("name, description, goals and stop_at follow a serving group", func(t *testing.T) {
		t.Parallel()
		group, experiment := newExperiment(t)
		startAt := experiment.StartAt
		group.Name = "renamed"
		group.Description = "new description"
		group.GoalIds = []string{"goal-1", "goal-2", "goal-3"}
		group.StopAt = experimentTestNow + 5_000
		group.Status = featureproto.HoldoutGroup_ANALYSIS
		assert.True(t, group.MirrorToExperiment(experiment, later))
		assert.Equal(t, "renamed", experiment.Name)
		assert.Equal(t, "new description", experiment.Description)
		assert.Equal(t, []string{"goal-1", "goal-2", "goal-3"}, experiment.GoalIds)
		assert.Equal(t, experimentTestNow+5_000, experiment.StopAt)
		assert.Equal(t, experimentproto.Experiment_RUNNING, experiment.Status)
		assert.Equal(t, int64(math.MaxInt64), experiment.StoppedAt)
		assert.Equal(t, startAt, experiment.StartAt)
		assert.Equal(t, later, experiment.UpdatedAt)
		assert.Equal(t, experimentTestNow, experiment.CreatedAt)
		assert.Equal(t, experimentproto.Experiment_HOLDOUT, experiment.Kind)

		// The goals are a copy: a later change of the group does not reach the experiment by itself.
		group.GoalIds[0] = "changed"
		assert.Equal(t, "goal-1", experiment.GoalIds[0])
	})

	t.Run("a manual stop stops the experiment at stopped_at", func(t *testing.T) {
		t.Parallel()
		group, experiment := newExperiment(t)
		require.NoError(t, group.StopManual(later))
		require.Equal(t, featureproto.HoldoutGroup_FORCE_STOPPED, group.Status)
		assert.True(t, group.MirrorToExperiment(experiment, later))
		assert.Equal(t, experimentproto.Experiment_STOPPED, experiment.Status)
		assert.Equal(t, later, experiment.StopAt)
		assert.Equal(t, later, experiment.StoppedAt)
		assert.False(t, experiment.Stopped)
		assert.False(t, experiment.Archived)
		assert.Less(t, experiment.StopAt, group.StopAt)
	})

	t.Run("a scheduled stop stops the experiment at stop_at", func(t *testing.T) {
		t.Parallel()
		group, experiment := newExperiment(t)
		finishedAt := group.StopAt + 120
		require.NoError(t, group.FinishScheduled(finishedAt))
		require.Equal(t, featureproto.HoldoutGroup_STOPPED, group.Status)
		assert.True(t, group.MirrorToExperiment(experiment, finishedAt))
		assert.Equal(t, experimentproto.Experiment_STOPPED, experiment.Status)
		assert.Equal(t, group.StopAt, experiment.StopAt)
		assert.Equal(t, finishedAt, experiment.StoppedAt)
		assert.Equal(t, finishedAt, experiment.UpdatedAt)
	})

	t.Run("archiving the group leaves the experiment as it is", func(t *testing.T) {
		t.Parallel()
		group, experiment := newExperiment(t)
		require.NoError(t, group.StopManual(later))
		require.True(t, group.MirrorToExperiment(experiment, later))
		require.NoError(t, group.Archive(later+10))
		assert.False(t, group.MirrorToExperiment(experiment, later+10))
		assert.False(t, experiment.Archived)
		assert.Equal(t, later, experiment.UpdatedAt)
	})

	t.Run("a stopped group can still be renamed", func(t *testing.T) {
		t.Parallel()
		group, experiment := newExperiment(t)
		require.NoError(t, group.StopManual(later))
		require.True(t, group.MirrorToExperiment(experiment, later))
		group.Name = "renamed after stop"
		assert.True(t, group.MirrorToExperiment(experiment, later+10))
		assert.Equal(t, "renamed after stop", experiment.Name)
		assert.Equal(t, experimentproto.Experiment_STOPPED, experiment.Status)
		assert.Equal(t, later, experiment.StopAt)
		assert.Equal(t, later, experiment.StoppedAt)
	})
}

func TestInternalExperimentStopAt(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		status   featureproto.HoldoutGroup_Status
		expected int64
	}{
		{desc: "waiting", status: featureproto.HoldoutGroup_WAITING, expected: 200},
		{desc: "running", status: featureproto.HoldoutGroup_RUNNING, expected: 200},
		{desc: "analysis", status: featureproto.HoldoutGroup_ANALYSIS, expected: 200},
		{desc: "a scheduled stop is dated by stop_at", status: featureproto.HoldoutGroup_STOPPED, expected: 200},
		{desc: "a manual stop is dated by stopped_at", status: featureproto.HoldoutGroup_FORCE_STOPPED, expected: 150},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			group := &featureproto.HoldoutGroup{Status: p.status, StopAt: 200, StoppedAt: 150}
			assert.Equal(t, p.expected, InternalExperimentStopAt(group))
		})
	}
}
