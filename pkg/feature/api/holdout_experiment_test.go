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

package api

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"

	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	experimentstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2/mock"
	featuredomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	featurestoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	holdoutdomain "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	holdoutmock "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	"github.com/bucketeer-io/bucketeer/v2/proto/common"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type holdoutTxKey struct{}

// inHoldoutTx matches the context the transaction hands to its callback, so a
// write that is expected with it is a write inside the transaction.
func inHoldoutTx() gomock.Matcher {
	return gomock.Cond(func(ctx context.Context) bool {
		marked, _ := ctx.Value(holdoutTxKey{}).(bool)
		return marked
	})
}

// measuredHoldoutTest is a holdout service whose transaction marks its
// context, and which records the group that is stored and the event that is
// published.
type measuredHoldoutTest struct {
	service     *FeatureService
	storage     *holdoutmock.MockHoldoutStorage
	experiments *experimentstoragemock.MockExperimentStorage
	stored      *featureproto.HoldoutGroup
	event       *eventproto.Event
}

// newMeasuredHoldoutTest expects one successful change of the group "hg":
// the transaction, the locks, the group write inside the transaction and one
// published event of eventType. transition adds the ListActiveGroups call of
// transitionHoldout. The experiment storage has no expectation, so a test
// that sets none proves that the change does not touch it.
func newMeasuredHoldoutTest(
	t *testing.T,
	group *holdoutdomain.HoldoutGroup,
	eventType eventproto.Event_Type,
	transition bool,
) *measuredHoldoutTest {
	t.Helper()
	controller := gomock.NewController(t)
	service, storage := newHoldoutService(t, controller)
	h := &measuredHoldoutTest{
		service:     service,
		storage:     storage,
		experiments: service.experimentStorage.(*experimentstoragemock.MockExperimentStorage),
	}
	service.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(context.WithValue(ctx, holdoutTxKey{}, true))
	})
	storage.EXPECT().LockEnvironment(inHoldoutTx(), "ns0").Return(nil)
	storage.EXPECT().GetHoldoutGroupForUpdate(inHoldoutTx(), "hg", "ns0").Return(group, nil)
	if transition {
		storage.EXPECT().ListActiveGroups(inHoldoutTx(), "ns0").Return(nil, nil)
	}
	storage.EXPECT().UpdateHoldoutGroup(inHoldoutTx(), gomock.Any()).DoAndReturn(
		func(_ context.Context, updated *holdoutdomain.HoldoutGroup) error {
			h.stored = updated.Clone().HoldoutGroup
			return nil
		},
	)
	storage.EXPECT().ListCurrentFeatures(gomock.Any(), "ns0", []string{"hg"}).Return(nil, nil)
	pub := publishermock.NewMockPublisher(controller)
	service.domainPublisher = pub
	pub.EXPECT().PublishMulti(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, messages []publisher.Message) map[string]error {
			require.Len(t, messages, 1)
			h.event = messages[0].(*eventproto.Event)
			assert.Equal(t, eventType, h.event.Type)
			assert.Equal(t, eventproto.Event_HOLDOUT_GROUP, h.event.EntityType)
			return nil
		},
	)
	return h
}

// expectCreate expects one internal experiment to be inserted inside the
// transaction and returns a pointer that holds it afterwards.
func (h *measuredHoldoutTest) expectCreate() **experimentproto.Experiment {
	created := new(*experimentproto.Experiment)
	h.experiments.EXPECT().CreateExperiment(inHoldoutTx(), gomock.Any(), "ns0").DoAndReturn(
		func(_ context.Context, e *experimentdomain.Experiment, _ string) error {
			*created = e.Experiment
			return nil
		},
	)
	return created
}

// expectRead expects the stored internal experiment to be read inside the
// transaction.
func (h *measuredHoldoutTest) expectRead(stored *experimentproto.Experiment) {
	h.experiments.EXPECT().GetExperiment(inHoldoutTx(), stored.Id, "ns0").Return(
		&experimentdomain.Experiment{Experiment: stored}, nil,
	)
}

// expectUpdate expects the internal experiment with the ID to be updated
// inside the transaction and returns a pointer that holds it afterwards.
func (h *measuredHoldoutTest) expectUpdate(id string) **experimentproto.Experiment {
	updated := new(*experimentproto.Experiment)
	h.experiments.EXPECT().UpdateExperiment(inHoldoutTx(), gomock.Any(), "ns0").DoAndReturn(
		func(_ context.Context, e *experimentdomain.Experiment, _ string) error {
			if e.Id != id {
				return errors.New("unexpected experiment " + e.Id)
			}
			*updated = e.Experiment
			return nil
		},
	)
	return updated
}

// eventGroup decodes the group that the published event carries.
func (h *measuredHoldoutTest) eventGroup(t *testing.T) *featureproto.HoldoutGroup {
	t.Helper()
	require.NotNil(t, h.event)
	var group featureproto.HoldoutGroup
	require.NoError(t, protojson.Unmarshal([]byte(h.event.EntityData), &group))
	return &group
}

func measuredHoldout(now int64, status featureproto.HoldoutGroup_Status) *holdoutdomain.HoldoutGroup {
	group := waitingHoldout(now)
	group.Status = status
	group.Description = "about"
	group.CreatedBy = "creator@example.com"
	group.GoalIds = []string{"goal-1"}
	if status != featureproto.HoldoutGroup_WAITING {
		group.StartAt = now - 100
	}
	if status == featureproto.HoldoutGroup_ANALYSIS {
		group.AnalysisStartedAt = now - 50
	}
	return group
}

// internalExperiment is a stored internal experiment of the group as the
// holdout service leaves it while the group serves.
func internalExperiment(id string, group *holdoutdomain.HoldoutGroup, startAt int64) *experimentproto.Experiment {
	return &experimentproto.Experiment{
		Id:             id,
		Kind:           experimentproto.Experiment_HOLDOUT,
		FeatureId:      group.Id,
		FeatureVersion: 1,
		Variations: []*featureproto.Variation{
			{Id: "in_holdout", Value: "in_holdout", Name: "in_holdout"},
			{Id: "not_in_holdout", Value: "not_in_holdout", Name: "not_in_holdout"},
		},
		BaseVariationId: "in_holdout",
		GoalIds:         append([]string(nil), group.GoalIds...),
		StartAt:         startAt,
		StopAt:          group.StopAt,
		StoppedAt:       math.MaxInt64,
		Status:          experimentproto.Experiment_RUNNING,
		Name:            group.Name,
		Description:     group.Description,
		Maintainer:      group.CreatedBy,
		CreatedAt:       startAt,
		UpdatedAt:       startAt,
	}
}

func assertInternalExperiment(t *testing.T, e *experimentproto.Experiment, group *featureproto.HoldoutGroup) {
	t.Helper()
	require.NotNil(t, e)
	assert.NotEmpty(t, e.Id)
	assert.Equal(t, experimentproto.Experiment_HOLDOUT, e.Kind)
	assert.Equal(t, group.Id, e.FeatureId)
	assert.Equal(t, int32(1), e.FeatureVersion)
	require.Len(t, e.Variations, 2)
	assert.Equal(t, "in_holdout", e.Variations[0].Id)
	assert.Equal(t, "not_in_holdout", e.Variations[1].Id)
	assert.Equal(t, "in_holdout", e.BaseVariationId)
	assert.Equal(t, group.GoalIds, e.GoalIds)
	assert.Equal(t, group.Name, e.Name)
	assert.Equal(t, group.Description, e.Description)
	assert.Equal(t, experimentproto.Experiment_RUNNING, e.Status)
	assert.Equal(t, group.StopAt, e.StopAt)
	assert.Equal(t, int64(math.MaxInt64), e.StoppedAt)
	assert.False(t, e.Archived)
	assert.False(t, e.Deleted)
}

func TestStartHoldoutGroupCreatesInternalExperiment(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()

	t.Run("manual start", func(t *testing.T) {
		t.Parallel()
		h := newMeasuredHoldoutTest(
			t, measuredHoldout(now, featureproto.HoldoutGroup_WAITING), eventproto.Event_HOLDOUT_GROUP_STARTED, true,
		)
		created := h.expectCreate()
		resp, err := h.service.StartHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StartHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		require.NoError(t, err)
		group := resp.HoldoutGroup
		assert.Equal(t, featureproto.HoldoutGroup_RUNNING, group.Status)
		assertInternalExperiment(t, *created, group)
		// The manual start clamped start_at to now, and the experiment starts with the group.
		assert.Equal(t, group.StartAt, (*created).StartAt)
		assert.Equal(t, "creator@example.com", (*created).Maintainer)
		assert.Equal(t, (*created).Id, group.ExperimentId)
		assert.Empty(t, group.AnalysisExperimentId)
		assert.Equal(t, (*created).Id, h.stored.ExperimentId)
		assert.Equal(t, (*created).Id, h.eventGroup(t).ExperimentId)
	})

	t.Run("scheduled start", func(t *testing.T) {
		t.Parallel()
		waiting := measuredHoldout(now, featureproto.HoldoutGroup_WAITING)
		waiting.StartAt = now - 30
		h := newMeasuredHoldoutTest(t, waiting, eventproto.Event_HOLDOUT_GROUP_STARTED, true)
		created := h.expectCreate()
		resp, err := h.service.StartHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StartHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0", Mode: featureproto.StartHoldoutGroupRequest_SCHEDULED,
			},
		)
		require.NoError(t, err)
		assertInternalExperiment(t, *created, resp.HoldoutGroup)
		assert.Equal(t, now-30, (*created).StartAt)
		assert.Equal(t, (*created).Id, resp.HoldoutGroup.ExperimentId)
		assert.Equal(t, (*created).Id, h.stored.ExperimentId)
	})

	t.Run("a failed experiment write fails the start and leaves the group unwritten", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").
			Return(measuredHoldout(now, featureproto.HoldoutGroup_WAITING), nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		service.experimentStorage.(*experimentstoragemock.MockExperimentStorage).EXPECT().
			CreateExperiment(gomock.Any(), gomock.Any(), "ns0").Return(errors.New("db"))
		// No UpdateHoldoutGroup and no published event are expected.
		_, err := service.StartHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StartHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		assert.Equal(t, codes.Unknown, status.Code(err))
	})
}

func TestStartHoldoutGroupAnalysisCreatesAnalysisExperiment(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	modes := map[string]featureproto.StartHoldoutGroupAnalysisRequest_Mode{
		"manual":    featureproto.StartHoldoutGroupAnalysisRequest_MANUAL,
		"scheduled": featureproto.StartHoldoutGroupAnalysisRequest_SCHEDULED,
	}
	for name, mode := range modes {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			running := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
			running.StartAnalysisAt = now - 10
			running.ExperimentId = "whole"
			whole := internalExperiment("whole", running, running.StartAt)
			h := newMeasuredHoldoutTest(t, running, eventproto.Event_HOLDOUT_GROUP_ANALYSIS_STARTED, true)
			created := h.expectCreate()
			// The experiment of the whole period is read and, being unchanged, not written.
			h.expectRead(whole)
			resp, err := h.service.StartHoldoutGroupAnalysis(
				createContextWithTokenRoleUnassigned(),
				&featureproto.StartHoldoutGroupAnalysisRequest{Id: "hg", EnvironmentId: "ns0", Mode: mode},
			)
			require.NoError(t, err)
			group := resp.HoldoutGroup
			assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, group.Status)
			assertInternalExperiment(t, *created, group)
			assert.Equal(t, group.AnalysisStartedAt, (*created).StartAt)
			assert.NotEqual(t, "whole", (*created).Id)
			assert.Equal(t, "whole", group.ExperimentId)
			assert.Equal(t, (*created).Id, group.AnalysisExperimentId)
			assert.Equal(t, (*created).Id, h.stored.AnalysisExperimentId)
			assert.Equal(t, "whole", h.stored.ExperimentId)
			assert.Equal(t, (*created).Id, h.eventGroup(t).AnalysisExperimentId)
		})
	}
}

func TestStopHoldoutGroupStopsInternalExperiments(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()

	t.Run("a manual stop ends both experiments at stopped_at, as STOPPED", func(t *testing.T) {
		t.Parallel()
		analysis := measuredHoldout(now, featureproto.HoldoutGroup_ANALYSIS)
		analysis.ExperimentId = "whole"
		analysis.AnalysisExperimentId = "analysis"
		h := newMeasuredHoldoutTest(t, analysis, eventproto.Event_HOLDOUT_GROUP_STOPPED, true)
		h.expectRead(internalExperiment("whole", analysis, analysis.StartAt))
		h.expectRead(internalExperiment("analysis", analysis, analysis.AnalysisStartedAt))
		wholeUpdated := h.expectUpdate("whole")
		analysisUpdated := h.expectUpdate("analysis")
		resp, err := h.service.StopHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StopHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		require.NoError(t, err)
		group := resp.HoldoutGroup
		assert.Equal(t, featureproto.HoldoutGroup_FORCE_STOPPED, group.Status)
		require.NotZero(t, group.StoppedAt)
		for _, updated := range []*experimentproto.Experiment{*wholeUpdated, *analysisUpdated} {
			require.NotNil(t, updated)
			assert.Equal(t, experimentproto.Experiment_STOPPED, updated.Status)
			assert.Equal(t, group.StoppedAt, updated.StopAt)
			assert.Equal(t, group.StoppedAt, updated.StoppedAt)
			assert.Less(t, updated.StopAt, group.StopAt)
			assert.False(t, updated.Stopped)
			assert.False(t, updated.Archived)
			assert.Equal(t, experimentproto.Experiment_HOLDOUT, updated.Kind)
		}
		assert.Equal(t, analysis.StartAt, (*wholeUpdated).StartAt)
		assert.Equal(t, now-50, (*analysisUpdated).StartAt)
	})

	t.Run("a scheduled stop ends the experiment at stop_at, as STOPPED", func(t *testing.T) {
		t.Parallel()
		running := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
		running.StopAt = now - 5
		running.ExperimentId = "whole"
		h := newMeasuredHoldoutTest(t, running, eventproto.Event_HOLDOUT_GROUP_STOPPED, true)
		h.expectRead(internalExperiment("whole", running, running.StartAt))
		updated := h.expectUpdate("whole")
		resp, err := h.service.StopHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StopHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0", Mode: featureproto.StopHoldoutGroupRequest_SCHEDULED,
			},
		)
		require.NoError(t, err)
		group := resp.HoldoutGroup
		assert.Equal(t, featureproto.HoldoutGroup_STOPPED, group.Status)
		require.NotNil(t, *updated)
		assert.Equal(t, experimentproto.Experiment_STOPPED, (*updated).Status)
		assert.Equal(t, now-5, (*updated).StopAt)
		assert.Equal(t, group.StoppedAt, (*updated).StoppedAt)
	})

	t.Run("a failed experiment write fails the stop and leaves the group unwritten", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		running := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
		running.ExperimentId = "whole"
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").Return(running, nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		experiments := service.experimentStorage.(*experimentstoragemock.MockExperimentStorage)
		experiments.EXPECT().GetExperiment(gomock.Any(), "whole", "ns0").Return(
			&experimentdomain.Experiment{Experiment: internalExperiment("whole", running, running.StartAt)}, nil,
		)
		experiments.EXPECT().UpdateExperiment(gomock.Any(), gomock.Any(), "ns0").Return(errors.New("db"))
		_, err := service.StopHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StopHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		assert.Equal(t, codes.Unknown, status.Code(err))
	})
}

func TestArchiveHoldoutGroupLeavesInternalExperiments(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	stopped := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
	stopped.ExperimentId = "whole"
	whole := internalExperiment("whole", stopped, stopped.StartAt)
	require.NoError(t, stopped.StopManual(now-20))
	require.True(t, stopped.MirrorToExperiment(&experimentdomain.Experiment{Experiment: whole}, now-20))

	h := newMeasuredHoldoutTest(t, stopped, eventproto.Event_HOLDOUT_GROUP_ARCHIVED, true)
	// The experiment is read and not written: no UpdateExperiment is expected.
	h.expectRead(whole)
	resp, err := h.service.ArchiveHoldoutGroup(
		createContextWithTokenRoleUnassigned(),
		&featureproto.ArchiveHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
	)
	require.NoError(t, err)
	assert.True(t, resp.HoldoutGroup.Archived)
	assert.False(t, whole.Archived)
	assert.Equal(t, "whole", resp.HoldoutGroup.ExperimentId)
}

func TestUpdateHoldoutGroupMirrorsInternalExperiments(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()

	t.Run("stop_at is written to both experiments", func(t *testing.T) {
		t.Parallel()
		analysis := measuredHoldout(now, featureproto.HoldoutGroup_ANALYSIS)
		analysis.ExperimentId = "whole"
		analysis.AnalysisExperimentId = "analysis"
		h := newMeasuredHoldoutTest(t, analysis, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		h.storage.EXPECT().ListActiveGroups(inHoldoutTx(), "ns0").Return(nil, nil)
		h.expectRead(internalExperiment("whole", analysis, analysis.StartAt))
		h.expectRead(internalExperiment("analysis", analysis, analysis.AnalysisStartedAt))
		wholeUpdated := h.expectUpdate("whole")
		analysisUpdated := h.expectUpdate("analysis")
		_, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0", StopAt: wrapperspb.Int64(now + 9_000),
			},
		)
		require.NoError(t, err)
		for _, updated := range []*experimentproto.Experiment{*wholeUpdated, *analysisUpdated} {
			require.NotNil(t, updated)
			assert.Equal(t, now+9_000, updated.StopAt)
			assert.Equal(t, experimentproto.Experiment_RUNNING, updated.Status)
		}
		assert.Equal(t, analysis.StartAt, (*wholeUpdated).StartAt)
		assert.Equal(t, now-50, (*analysisUpdated).StartAt)
	})

	t.Run("start_analysis_at alone writes no experiment", func(t *testing.T) {
		t.Parallel()
		running := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
		running.ExperimentId = "whole"
		h := newMeasuredHoldoutTest(t, running, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		h.storage.EXPECT().ListActiveGroups(inHoldoutTx(), "ns0").Return(nil, nil)
		h.expectRead(internalExperiment("whole", running, running.StartAt))
		_, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0", StartAnalysisAt: wrapperspb.Int64(now + 500),
			},
		)
		require.NoError(t, err)
	})

	t.Run("name and description are written to the experiment", func(t *testing.T) {
		t.Parallel()
		running := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
		running.ExperimentId = "whole"
		h := newMeasuredHoldoutTest(t, running, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		h.expectRead(internalExperiment("whole", running, running.StartAt))
		updated := h.expectUpdate("whole")
		_, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0",
				Name:        wrapperspb.String("renamed"),
				Description: wrapperspb.String("new description"),
			},
		)
		require.NoError(t, err)
		require.NotNil(t, *updated)
		assert.Equal(t, "renamed", (*updated).Name)
		assert.Equal(t, "new description", (*updated).Description)
	})

	t.Run("a stopped group can be renamed, and its experiment stays stopped", func(t *testing.T) {
		t.Parallel()
		stopped := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
		stopped.ExperimentId = "whole"
		whole := internalExperiment("whole", stopped, stopped.StartAt)
		require.NoError(t, stopped.StopManual(now-20))
		require.True(t, stopped.MirrorToExperiment(&experimentdomain.Experiment{Experiment: whole}, now-20))
		h := newMeasuredHoldoutTest(t, stopped, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		h.expectRead(whole)
		updated := h.expectUpdate("whole")
		_, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0", Name: wrapperspb.String("renamed"),
			},
		)
		require.NoError(t, err)
		require.NotNil(t, *updated)
		assert.Equal(t, "renamed", (*updated).Name)
		assert.Equal(t, experimentproto.Experiment_STOPPED, (*updated).Status)
		assert.Equal(t, now-20, (*updated).StopAt)
		assert.Equal(t, now-20, (*updated).StoppedAt)
	})

	t.Run("goals added to a group that has experiments are written to both", func(t *testing.T) {
		t.Parallel()
		analysis := measuredHoldout(now, featureproto.HoldoutGroup_ANALYSIS)
		analysis.ExperimentId = "whole"
		analysis.AnalysisExperimentId = "analysis"
		h := newMeasuredHoldoutTest(t, analysis, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		expectGoal(h.service, "goal-1")
		expectGoal(h.service, "goal-2")
		h.expectRead(internalExperiment("whole", analysis, analysis.StartAt))
		h.expectRead(internalExperiment("analysis", analysis, analysis.AnalysisStartedAt))
		wholeUpdated := h.expectUpdate("whole")
		analysisUpdated := h.expectUpdate("analysis")
		resp, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0",
				GoalIds: &common.StringListValue{Values: []string{"goal-1", "goal-2"}},
			},
		)
		require.NoError(t, err)
		for _, updated := range []*experimentproto.Experiment{*wholeUpdated, *analysisUpdated} {
			require.NotNil(t, updated)
			assert.Equal(t, []string{"goal-1", "goal-2"}, updated.GoalIds)
		}
		// No experiment is created, and the windows are left alone.
		assert.Equal(t, analysis.StartAt, (*wholeUpdated).StartAt)
		assert.Equal(t, now-50, (*analysisUpdated).StartAt)
		assert.Equal(t, "whole", resp.HoldoutGroup.ExperimentId)
		assert.Equal(t, "analysis", resp.HoldoutGroup.AnalysisExperimentId)
	})

	t.Run("first goals while RUNNING create the experiment, measuring from now", func(t *testing.T) {
		t.Parallel()
		running := measuredHoldout(now, featureproto.HoldoutGroup_RUNNING)
		running.GoalIds = []string{}
		h := newMeasuredHoldoutTest(t, running, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		expectGoal(h.service, "goal-1")
		created := h.expectCreate()
		resp, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0",
				GoalIds: &common.StringListValue{Values: []string{"goal-1"}},
			},
		)
		require.NoError(t, err)
		group := resp.HoldoutGroup
		assertInternalExperiment(t, *created, group)
		assert.Greater(t, (*created).StartAt, group.StartAt)
		assert.GreaterOrEqual(t, (*created).StartAt, now)
		assert.Equal(t, group.UpdatedAt, (*created).StartAt)
		assert.Equal(t, (*created).Id, group.ExperimentId)
		assert.Empty(t, group.AnalysisExperimentId)
		assert.Equal(t, (*created).Id, h.stored.ExperimentId)
		assert.Equal(t, (*created).Id, h.eventGroup(t).ExperimentId)
	})

	t.Run("first goals while in ANALYSIS create only the analysis experiment, measuring from now", func(t *testing.T) {
		t.Parallel()
		analysis := measuredHoldout(now, featureproto.HoldoutGroup_ANALYSIS)
		analysis.GoalIds = []string{}
		h := newMeasuredHoldoutTest(t, analysis, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		expectGoal(h.service, "goal-1")
		created := h.expectCreate()
		resp, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0",
				GoalIds: &common.StringListValue{Values: []string{"goal-1"}},
			},
		)
		require.NoError(t, err)
		group := resp.HoldoutGroup
		assertInternalExperiment(t, *created, group)
		assert.Greater(t, (*created).StartAt, group.AnalysisStartedAt)
		assert.Equal(t, group.UpdatedAt, (*created).StartAt)
		assert.Empty(t, group.ExperimentId)
		assert.Equal(t, (*created).Id, group.AnalysisExperimentId)
		assert.Empty(t, h.stored.ExperimentId)
		assert.Equal(t, (*created).Id, h.stored.AnalysisExperimentId)
	})

	t.Run("goals on a WAITING group create nothing", func(t *testing.T) {
		t.Parallel()
		waiting := waitingHoldout(now)
		h := newMeasuredHoldoutTest(t, waiting, eventproto.Event_HOLDOUT_GROUP_UPDATED, false)
		expectGoal(h.service, "goal-1")
		resp, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0",
				GoalIds: &common.StringListValue{Values: []string{"goal-1"}},
			},
		)
		require.NoError(t, err)
		assert.Equal(t, []string{"goal-1"}, resp.HoldoutGroup.GoalIds)
		assert.Empty(t, resp.HoldoutGroup.ExperimentId)
	})
}

// A group without goals has no internal experiment at any point of its life:
// the experiment storage has no expectation, so any call to it fails the test.
func TestHoldoutGroupWithoutGoalsNeverTouchesExperiments(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	withoutGoals := func(status featureproto.HoldoutGroup_Status) *holdoutdomain.HoldoutGroup {
		group := measuredHoldout(now, status)
		group.GoalIds = []string{}
		return group
	}
	assertNoExperiment := func(t *testing.T, h *measuredHoldoutTest, group *featureproto.HoldoutGroup) {
		t.Helper()
		assert.Empty(t, group.ExperimentId)
		assert.Empty(t, group.AnalysisExperimentId)
		assert.Empty(t, h.stored.ExperimentId)
		assert.Empty(t, h.stored.AnalysisExperimentId)
		assert.Empty(t, h.eventGroup(t).ExperimentId)
		assert.Empty(t, h.eventGroup(t).AnalysisExperimentId)
	}

	t.Run("start", func(t *testing.T) {
		t.Parallel()
		for _, mode := range []featureproto.StartHoldoutGroupRequest_Mode{
			featureproto.StartHoldoutGroupRequest_MANUAL, featureproto.StartHoldoutGroupRequest_SCHEDULED,
		} {
			waiting := withoutGoals(featureproto.HoldoutGroup_WAITING)
			waiting.StartAt = now - 30
			h := newMeasuredHoldoutTest(t, waiting, eventproto.Event_HOLDOUT_GROUP_STARTED, true)
			resp, err := h.service.StartHoldoutGroup(
				createContextWithTokenRoleUnassigned(),
				&featureproto.StartHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0", Mode: mode},
			)
			require.NoError(t, err)
			assert.Equal(t, featureproto.HoldoutGroup_RUNNING, resp.HoldoutGroup.Status)
			assertNoExperiment(t, h, resp.HoldoutGroup)
		}
	})

	t.Run("start analysis", func(t *testing.T) {
		t.Parallel()
		for _, mode := range []featureproto.StartHoldoutGroupAnalysisRequest_Mode{
			featureproto.StartHoldoutGroupAnalysisRequest_MANUAL, featureproto.StartHoldoutGroupAnalysisRequest_SCHEDULED,
		} {
			running := withoutGoals(featureproto.HoldoutGroup_RUNNING)
			running.StartAnalysisAt = now - 10
			h := newMeasuredHoldoutTest(t, running, eventproto.Event_HOLDOUT_GROUP_ANALYSIS_STARTED, true)
			resp, err := h.service.StartHoldoutGroupAnalysis(
				createContextWithTokenRoleUnassigned(),
				&featureproto.StartHoldoutGroupAnalysisRequest{Id: "hg", EnvironmentId: "ns0", Mode: mode},
			)
			require.NoError(t, err)
			assert.Equal(t, featureproto.HoldoutGroup_ANALYSIS, resp.HoldoutGroup.Status)
			assertNoExperiment(t, h, resp.HoldoutGroup)
		}
	})

	t.Run("stop", func(t *testing.T) {
		t.Parallel()
		for _, mode := range []featureproto.StopHoldoutGroupRequest_Mode{
			featureproto.StopHoldoutGroupRequest_MANUAL, featureproto.StopHoldoutGroupRequest_SCHEDULED,
		} {
			analysis := withoutGoals(featureproto.HoldoutGroup_ANALYSIS)
			analysis.StopAt = now - 5
			h := newMeasuredHoldoutTest(t, analysis, eventproto.Event_HOLDOUT_GROUP_STOPPED, true)
			resp, err := h.service.StopHoldoutGroup(
				createContextWithTokenRoleUnassigned(),
				&featureproto.StopHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0", Mode: mode},
			)
			require.NoError(t, err)
			assertNoExperiment(t, h, resp.HoldoutGroup)
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()
		h := newMeasuredHoldoutTest(
			t, withoutGoals(featureproto.HoldoutGroup_RUNNING), eventproto.Event_HOLDOUT_GROUP_UPDATED, false,
		)
		h.storage.EXPECT().ListActiveGroups(inHoldoutTx(), "ns0").Return(nil, nil)
		resp, err := h.service.UpdateHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.UpdateHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0",
				Name:   wrapperspb.String("renamed"),
				StopAt: wrapperspb.Int64(now + 9_000),
			},
		)
		require.NoError(t, err)
		assertNoExperiment(t, h, resp.HoldoutGroup)
	})

	t.Run("archive", func(t *testing.T) {
		t.Parallel()
		h := newMeasuredHoldoutTest(
			t, withoutGoals(featureproto.HoldoutGroup_STOPPED), eventproto.Event_HOLDOUT_GROUP_ARCHIVED, true,
		)
		resp, err := h.service.ArchiveHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.ArchiveHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		require.NoError(t, err)
		assertNoExperiment(t, h, resp.HoldoutGroup)
	})
}

func TestStartHoldoutGroupSnapshotsSegmentUsers(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	segmentGroup := func(goalIDs []string) *holdoutdomain.HoldoutGroup {
		group := measuredHoldout(now, featureproto.HoldoutGroup_WAITING)
		group.AssignmentType = featureproto.HoldoutGroup_SEGMENT
		group.SegmentId = "seg"
		group.GoalIds = goalIDs
		return group
	}
	expectSegment := func(h *measuredHoldoutTest) *gomock.Call {
		return h.service.segmentStorage.(*featurestoragemock.MockSegmentStorage).EXPECT().GetSegmentForUpdate(
			inHoldoutTx(), "seg", "ns0",
		).Return(&featuredomain.Segment{Segment: &featureproto.Segment{
			Id: "seg", Status: featureproto.Segment_SUCEEDED,
		}}, nil, nil)
	}

	t.Run("a SEGMENT group without goals gets a snapshot and no experiment", func(t *testing.T) {
		t.Parallel()
		h := newMeasuredHoldoutTest(t, segmentGroup([]string{}), eventproto.Event_HOLDOUT_GROUP_STARTED, true)
		locked := expectSegment(h)
		h.storage.EXPECT().CreateSegmentUserSnapshot(inHoldoutTx(), "hg", "seg", "ns0").
			Return(int64(3), nil).After(locked)
		resp, err := h.service.StartHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StartHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		require.NoError(t, err)
		assert.Empty(t, resp.HoldoutGroup.ExperimentId)
	})

	t.Run("a SEGMENT group with goals gets a snapshot and the experiment", func(t *testing.T) {
		t.Parallel()
		h := newMeasuredHoldoutTest(t, segmentGroup([]string{"goal-1"}), eventproto.Event_HOLDOUT_GROUP_STARTED, true)
		locked := expectSegment(h)
		h.storage.EXPECT().CreateSegmentUserSnapshot(inHoldoutTx(), "hg", "seg", "ns0").
			Return(int64(3), nil).After(locked)
		created := h.expectCreate()
		resp, err := h.service.StartHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StartHoldoutGroupRequest{
				Id: "hg", EnvironmentId: "ns0", Mode: featureproto.StartHoldoutGroupRequest_MANUAL,
			},
		)
		require.NoError(t, err)
		assert.Equal(t, (*created).Id, resp.HoldoutGroup.ExperimentId)
	})

	t.Run("a failed snapshot fails the start and leaves the group unwritten", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		service, storage := newHoldoutService(t, controller)
		expectTx(service)
		storage.EXPECT().LockEnvironment(gomock.Any(), "ns0").Return(nil)
		storage.EXPECT().GetHoldoutGroupForUpdate(gomock.Any(), "hg", "ns0").
			Return(segmentGroup([]string{"goal-1"}), nil)
		storage.EXPECT().ListActiveGroups(gomock.Any(), "ns0").Return(nil, nil)
		service.segmentStorage.(*featurestoragemock.MockSegmentStorage).EXPECT().GetSegmentForUpdate(
			gomock.Any(), "seg", "ns0",
		).Return(&featuredomain.Segment{Segment: &featureproto.Segment{Id: "seg"}}, nil, nil)
		storage.EXPECT().CreateSegmentUserSnapshot(gomock.Any(), "hg", "seg", "ns0").
			Return(int64(0), errors.New("db"))
		// No experiment write, no UpdateHoldoutGroup and no published event are expected.
		_, err := service.StartHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StartHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		assert.Equal(t, codes.Unknown, status.Code(err))
	})

	t.Run("later transitions of a SEGMENT group write no snapshot", func(t *testing.T) {
		t.Parallel()
		running := segmentGroup([]string{})
		running.Status = featureproto.HoldoutGroup_RUNNING
		running.StartAt = now - 100
		// CreateSegmentUserSnapshot has no expectation.
		h := newMeasuredHoldoutTest(t, running, eventproto.Event_HOLDOUT_GROUP_ANALYSIS_STARTED, true)
		_, err := h.service.StartHoldoutGroupAnalysis(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StartHoldoutGroupAnalysisRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		require.NoError(t, err)

		stopping := segmentGroup([]string{})
		stopping.Status = featureproto.HoldoutGroup_RUNNING
		stopping.StartAt = now - 100
		h = newMeasuredHoldoutTest(t, stopping, eventproto.Event_HOLDOUT_GROUP_STOPPED, true)
		_, err = h.service.StopHoldoutGroup(
			createContextWithTokenRoleUnassigned(),
			&featureproto.StopHoldoutGroupRequest{Id: "hg", EnvironmentId: "ns0"},
		)
		require.NoError(t, err)
	})
}
