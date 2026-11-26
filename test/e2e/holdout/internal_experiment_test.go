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

package holdout

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"

	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	holdoutdomain "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	"github.com/bucketeer-io/bucketeer/v2/proto/common"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// A group with goals gets one internal experiment when it starts and a second
// one when its analysis period starts. Both follow the group and only the
// holdout service can write them.
func TestInternalExperimentManual(t *testing.T) {
	servingMu.Lock()
	defer servingMu.Unlock()
	client := newFeatureClient(t)
	experimentClient := newExperimentClient(t)
	goalA := createGoal(t, experimentClient)
	goalB := createGoal(t, experimentClient)

	req := servingHashSpec().request(t)
	req.GoalIds = []string{goalA}
	group := createHoldoutGroup(t, client, req)
	assert.Equal(t, []string{goalA}, group.GoalIds)
	assert.Empty(t, group.ExperimentId, "a WAITING group has no internal experiment")
	assert.Empty(t, group.AnalysisExperimentId)
	assert.Empty(t, listInternalExperiments(t, experimentClient, group.Id))

	running := startHoldoutGroup(t, client, group.Id)
	require.NotEmpty(t, running.ExperimentId, "Start creates the experiment of the whole period")
	assert.Empty(t, running.AnalysisExperimentId)
	whole := getExperiment(t, experimentClient, running.ExperimentId)
	assert.Equal(t, experimentproto.Experiment_HOLDOUT, whole.Kind)
	assert.Equal(t, group.Id, whole.FeatureId, "feature_id is the group ID")
	assert.Equal(t, holdoutdomain.InternalExperimentFeatureVersion, whole.FeatureVersion)
	assert.Equal(t,
		[]string{holdoutdomain.InHoldoutVariationID, holdoutdomain.NotInHoldoutVariationID},
		variationIDs(whole),
	)
	assert.Equal(t, holdoutdomain.InHoldoutVariationID, whole.BaseVariationId)
	assert.Equal(t, running.CreatedBy, whole.Maintainer)
	requireMirrorsServingGroup(t, whole, running)
	assert.Equal(t, running.StartAt, whole.StartAt, "starts at the clamped start_at of the group")
	assert.Equal(t, []string{running.ExperimentId}, experimentIDs(listInternalExperiments(t, experimentClient, group.Id)))

	t.Run("the experiment API rejects writes to an internal experiment", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := experimentClient.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
			Id:            running.ExperimentId,
			EnvironmentId: holdoutEnv.id,
			Name:          wrapperspb.String("renamed-through-the-experiment-api"),
		})
		requireCode(t, err, codes.FailedPrecondition, "UpdateExperiment")
		_, err = experimentClient.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
			Id:            running.ExperimentId,
			EnvironmentId: holdoutEnv.id,
			Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
				Status: experimentproto.Experiment_FORCE_STOPPED,
			},
		})
		requireCode(t, err, codes.FailedPrecondition, "UpdateExperiment status")
		_, err = experimentClient.DeleteExperiment(ctx, &experimentproto.DeleteExperimentRequest{
			Id:            running.ExperimentId,
			EnvironmentId: holdoutEnv.id,
		})
		requireCode(t, err, codes.FailedPrecondition, "DeleteExperiment")
		requireMirrorsServingGroup(t, getExperiment(t, experimentClient, running.ExperimentId), running)
	})

	// The name, the description, the goals and stop_at follow the group.
	bumpSecond(t, running)
	updated := updateHoldoutGroup(t, client, &featureproto.UpdateHoldoutGroupRequest{
		Id:          group.Id,
		Name:        wrapperspb.String(newGroupName(t)),
		Description: wrapperspb.String("e2e-test-holdout-group-renamed"),
		GoalIds:     &common.StringListValue{Values: []string{goalA, goalB}},
		StopAt:      wrapperspb.Int64(running.StopAt + 600),
	})
	assert.Equal(t, running.ExperimentId, updated.ExperimentId)
	assert.NotEqual(t, running.Name, updated.Name)
	mirrored := getExperiment(t, experimentClient, running.ExperimentId)
	requireMirrorsServingGroup(t, mirrored, updated)
	assert.Equal(t, whole.StartAt, mirrored.StartAt, "start_at is left as it is")
	assert.Greater(t, mirrored.UpdatedAt, whole.UpdatedAt)

	t.Run("goals cannot be removed while the group is serving", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.UpdateHoldoutGroup(ctx, &featureproto.UpdateHoldoutGroupRequest{
			Id:            group.Id,
			EnvironmentId: holdoutEnv.id,
			GoalIds:       &common.StringListValue{Values: []string{goalA}},
		})
		requireCode(t, err, codes.FailedPrecondition)
		assert.Equal(t, []string{goalA, goalB}, getExperiment(t, experimentClient, running.ExperimentId).GoalIds)
	})

	analysis := startHoldoutGroupAnalysis(t, client, group.Id)
	assert.Equal(t, running.ExperimentId, analysis.ExperimentId)
	require.NotEmpty(t, analysis.AnalysisExperimentId, "StartAnalysis creates the experiment of the analysis period")
	assert.NotEqual(t, analysis.ExperimentId, analysis.AnalysisExperimentId)
	analysisExperiment := getExperiment(t, experimentClient, analysis.AnalysisExperimentId)
	assert.Equal(t, experimentproto.Experiment_HOLDOUT, analysisExperiment.Kind)
	assert.Equal(t, group.Id, analysisExperiment.FeatureId)
	requireMirrorsServingGroup(t, analysisExperiment, analysis)
	assert.Equal(t, analysis.AnalysisStartedAt, analysisExperiment.StartAt, "starts when the analysis period starts")
	requireMirrorsServingGroup(t, getExperiment(t, experimentClient, analysis.ExperimentId), analysis)
	assert.ElementsMatch(t,
		[]string{analysis.ExperimentId, analysis.AnalysisExperimentId},
		experimentIDs(listInternalExperiments(t, experimentClient, group.Id)),
	)

	// A manual stop ends both experiments at stopped_at. They are STOPPED,
	// not FORCE_STOPPED, although the group is FORCE_STOPPED.
	stopped := stopHoldoutGroup(t, client, group.Id)
	require.Equal(t, featureproto.HoldoutGroup_FORCE_STOPPED, stopped.Status)
	for _, id := range []string{stopped.ExperimentId, stopped.AnalysisExperimentId} {
		experiment := getExperiment(t, experimentClient, id)
		assert.Equal(t, experimentproto.Experiment_STOPPED, experiment.Status, "experiment %s", id)
		assert.Equal(t, stopped.StoppedAt, experiment.StopAt, "experiment %s", id)
		assert.Equal(t, stopped.StoppedAt, experiment.StoppedAt, "experiment %s", id)
		assert.Equal(t, updated.Name, experiment.Name, "experiment %s", id)
	}
	assert.Equal(t, whole.StartAt, getExperiment(t, experimentClient, stopped.ExperimentId).StartAt)
}

// The HoldoutGroupStatusUpdater creates and stops the internal experiments as
// the manual operations do. A scheduled start does not clamp start_at and a
// scheduled stop keeps stop_at.
func TestInternalExperimentScheduled(t *testing.T) {
	servingMu.Lock()
	defer servingMu.Unlock()
	client := newFeatureClient(t)
	experimentClient := newExperimentClient(t)
	goal := createGoal(t, experimentClient)

	now := time.Now().Unix()
	spec := groupSpec{
		assignment: featureproto.HoldoutGroup_HASH, percentage: 50,
		startAt: now + 10, startAnalysisAt: now + 25, stopAt: now + 40,
	}
	req := spec.request(t)
	req.GoalIds = []string{goal}
	group := createHoldoutGroup(t, client, req)

	waitUntilUnix(t, spec.startAt)
	runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
	running := getHoldoutGroup(t, client, group.Id)
	require.Equal(t, featureproto.HoldoutGroup_RUNNING, running.Status)
	require.NotEmpty(t, running.ExperimentId)
	assert.Empty(t, running.AnalysisExperimentId)
	whole := getExperiment(t, experimentClient, running.ExperimentId)
	assert.Equal(t, experimentproto.Experiment_HOLDOUT, whole.Kind)
	assert.Equal(t, group.Id, whole.FeatureId)
	requireMirrorsServingGroup(t, whole, running)
	assert.Equal(t, spec.startAt, whole.StartAt)

	waitUntilUnix(t, spec.startAnalysisAt)
	runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
	analysis := getHoldoutGroup(t, client, group.Id)
	require.Equal(t, featureproto.HoldoutGroup_ANALYSIS, analysis.Status)
	assert.Equal(t, running.ExperimentId, analysis.ExperimentId)
	require.NotEmpty(t, analysis.AnalysisExperimentId)
	analysisExperiment := getExperiment(t, experimentClient, analysis.AnalysisExperimentId)
	requireMirrorsServingGroup(t, analysisExperiment, analysis)
	assert.Equal(t, analysis.AnalysisStartedAt, analysisExperiment.StartAt)

	waitUntilUnix(t, spec.stopAt)
	runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
	stopped := getHoldoutGroup(t, client, group.Id)
	require.Equal(t, featureproto.HoldoutGroup_STOPPED, stopped.Status)
	for _, id := range []string{stopped.ExperimentId, stopped.AnalysisExperimentId} {
		experiment := getExperiment(t, experimentClient, id)
		assert.Equal(t, experimentproto.Experiment_STOPPED, experiment.Status, "experiment %s", id)
		assert.Equal(t, spec.stopAt, experiment.StopAt, "experiment %s", id)
		assert.Equal(t, stopped.StoppedAt, experiment.StoppedAt, "experiment %s", id)
	}
}

// A group without goals has no internal experiment. When it gets its first
// goals while serving, the experiment of the current period is created and
// starts then, not at the start of the group.
func TestInternalExperimentGoalsAddedWhileServing(t *testing.T) {
	servingMu.Lock()
	defer servingMu.Unlock()
	client := newFeatureClient(t)
	experimentClient := newExperimentClient(t)
	goal := createGoal(t, experimentClient)

	t.Run("while RUNNING", func(t *testing.T) {
		group := createHoldoutGroup(t, client, servingHashSpec().request(t))
		running := startHoldoutGroup(t, client, group.Id)
		assert.Empty(t, running.ExperimentId, "a group without goals has no internal experiment")
		assert.Empty(t, listInternalExperiments(t, experimentClient, group.Id))

		bumpSecond(t, running)
		before := time.Now().Unix()
		withGoal := updateHoldoutGroup(t, client, &featureproto.UpdateHoldoutGroupRequest{
			Id:      group.Id,
			GoalIds: &common.StringListValue{Values: []string{goal}},
		})
		require.NotEmpty(t, withGoal.ExperimentId)
		assert.Empty(t, withGoal.AnalysisExperimentId)
		whole := getExperiment(t, experimentClient, withGoal.ExperimentId)
		requireMirrorsServingGroup(t, whole, withGoal)
		assert.GreaterOrEqual(t, whole.StartAt, before, "starts when the goals are added")
		assert.Greater(t, whole.StartAt, running.StartAt)

		analysis := startHoldoutGroupAnalysis(t, client, group.Id)
		assert.Equal(t, withGoal.ExperimentId, analysis.ExperimentId)
		require.NotEmpty(t, analysis.AnalysisExperimentId)

		stopped := stopHoldoutGroup(t, client, group.Id)
		for _, id := range []string{stopped.ExperimentId, stopped.AnalysisExperimentId} {
			assert.Equal(t, experimentproto.Experiment_STOPPED, getExperiment(t, experimentClient, id).Status)
		}
	})

	t.Run("while in ANALYSIS", func(t *testing.T) {
		group := createHoldoutGroup(t, client, servingHashSpec().request(t))
		startHoldoutGroup(t, client, group.Id)
		analysis := startHoldoutGroupAnalysis(t, client, group.Id)
		assert.Empty(t, analysis.ExperimentId)
		assert.Empty(t, analysis.AnalysisExperimentId)

		bumpSecond(t, analysis)
		before := time.Now().Unix()
		withGoal := updateHoldoutGroup(t, client, &featureproto.UpdateHoldoutGroupRequest{
			Id:      group.Id,
			GoalIds: &common.StringListValue{Values: []string{goal}},
		})
		assert.Empty(t, withGoal.ExperimentId, "an experiment of the whole period would have the same window")
		require.NotEmpty(t, withGoal.AnalysisExperimentId)
		analysisExperiment := getExperiment(t, experimentClient, withGoal.AnalysisExperimentId)
		requireMirrorsServingGroup(t, analysisExperiment, withGoal)
		assert.GreaterOrEqual(t, analysisExperiment.StartAt, before)
		assert.Equal(t,
			[]string{withGoal.AnalysisExperimentId},
			experimentIDs(listInternalExperiments(t, experimentClient, group.Id)),
		)

		stopped := stopHoldoutGroup(t, client, group.Id)
		assert.Empty(t, stopped.ExperimentId)
		assert.Equal(t,
			experimentproto.Experiment_STOPPED,
			getExperiment(t, experimentClient, stopped.AnalysisExperimentId).Status,
		)
	})

	t.Run("a group that never has goals never has an experiment", func(t *testing.T) {
		group := createHoldoutGroup(t, client, servingHashSpec().request(t))
		startHoldoutGroup(t, client, group.Id)
		startHoldoutGroupAnalysis(t, client, group.Id)
		stopped := stopHoldoutGroup(t, client, group.Id)
		assert.Empty(t, stopped.ExperimentId)
		assert.Empty(t, stopped.AnalysisExperimentId)
		assert.Empty(t, listInternalExperiments(t, experimentClient, group.Id))
	})
}

// The segment of a SEGMENT group includes at most MaxSegmentUsers users. The
// limit is checked when the group is created and again when it starts,
// because the segment can grow in between.
func TestSegmentModeSizeLimit(t *testing.T) {
	client := newFeatureClient(t)
	segment := createSegment(t, client)
	userIDs := make([]string, holdoutdomain.MaxSegmentUsers)
	for i := range userIDs {
		userIDs[i] = fmt.Sprintf("u%d", i)
	}
	uploadSegmentUsersAndWaitForCount(t, client, segment.Id, userIDs)

	start, stop := farFutureWindow()
	group := createHoldoutGroup(t, client, groupSpec{
		assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, startAt: start, stopAt: stop,
	}.request(t))
	require.Equal(t, featureproto.HoldoutGroup_WAITING, group.Status, "a segment at the limit is accepted")

	// An upload adds to the segment, so one more user puts it over the limit.
	uploadSegmentUsersAndWaitForCount(t, client, segment.Id, []string{"one-over-the-limit"})
	tooLarge := fmt.Sprintf("segment %s includes more than %d users", segment.Id, holdoutdomain.MaxSegmentUsers)

	t.Run("create is rejected", func(t *testing.T) {
		start, stop := farFutureWindow()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.CreateHoldoutGroup(ctx, groupSpec{
			assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, startAt: start, stopAt: stop,
		}.request(t))
		requireCode(t, err, codes.FailedPrecondition)
		assert.Contains(t, err.Error(), tooLarge)
	})

	t.Run("start is rejected and the group stays WAITING", func(t *testing.T) {
		servingMu.Lock()
		defer servingMu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{
			Id: group.Id, EnvironmentId: holdoutEnv.id,
		})
		requireCode(t, err, codes.FailedPrecondition)
		assert.Contains(t, err.Error(), tooLarge)
		got := getHoldoutGroup(t, client, group.Id)
		assert.Equal(t, featureproto.HoldoutGroup_WAITING, got.Status)
		assert.Equal(t, group.StartAt, got.StartAt)
	})
}

// ---- helpers ----

func newExperimentClient(t *testing.T) experimentclient.Client {
	t.Helper()
	opts, err := dialOptions(*orgOwnerDefaultAccessTokenPath)
	require.NoError(t, err)
	c, err := experimentclient.NewClient(webGatewayTarget(), *webGatewayCert, opts...)
	require.NoError(t, err)
	t.Cleanup(func() { c.Close() })
	return c
}

// createGoal creates a goal in the holdout environment and returns its ID.
func createGoal(t *testing.T, client experimentclient.Client) string {
	t.Helper()
	id := prefixed("goal-id", newUUID(t))
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		EnvironmentId:  holdoutEnv.id,
		Id:             id,
		Name:           id,
		Description:    "e2e-test-holdout-goal",
		ConnectionType: experimentproto.Goal_EXPERIMENT,
	})
	require.NoError(t, err)
	return id
}

func getExperiment(t *testing.T, client experimentclient.Client, id string) *experimentproto.Experiment {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		Id: id, EnvironmentId: holdoutEnv.id,
	})
	require.NoError(t, err)
	return resp.Experiment
}

// listInternalExperiments lists the experiments of kind HOLDOUT whose
// feature_id is the group ID.
func listInternalExperiments(
	t *testing.T,
	client experimentclient.Client,
	groupID string,
) []*experimentproto.Experiment {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
		EnvironmentId: holdoutEnv.id,
		FeatureId:     groupID,
		PageSize:      100,
		Kinds:         []experimentproto.Experiment_Kind{experimentproto.Experiment_HOLDOUT},
	})
	require.NoError(t, err)
	return resp.Experiments
}

func experimentIDs(experiments []*experimentproto.Experiment) []string {
	ids := make([]string, 0, len(experiments))
	for _, e := range experiments {
		ids = append(ids, e.Id)
	}
	return ids
}

func variationIDs(experiment *experimentproto.Experiment) []string {
	ids := make([]string, 0, len(experiment.Variations))
	for _, v := range experiment.Variations {
		ids = append(ids, v.Id)
	}
	return ids
}

// requireMirrorsServingGroup asserts what an internal experiment takes from a
// RUNNING or ANALYSIS group: it is RUNNING, has the name, the description and
// the goals of the group, and stops when the group does.
func requireMirrorsServingGroup(
	t *testing.T,
	experiment *experimentproto.Experiment,
	group *featureproto.HoldoutGroup,
) {
	t.Helper()
	require.Equal(t, experimentproto.Experiment_RUNNING, experiment.Status, "status of %s", experiment.Id)
	require.Equal(t, group.Name, experiment.Name, "name of %s", experiment.Id)
	require.Equal(t, group.Description, experiment.Description, "description of %s", experiment.Id)
	require.Equal(t, group.GoalIds, experiment.GoalIds, "goal_ids of %s", experiment.Id)
	require.Equal(t, group.StopAt, experiment.StopAt, "stop_at of %s", experiment.Id)
	require.Equal(t, int64(math.MaxInt64), experiment.StoppedAt, "stopped_at of %s", experiment.Id)
	require.False(t, experiment.Archived, "archived of %s", experiment.Id)
	require.False(t, experiment.Deleted, "deleted of %s", experiment.Id)
}

// uploadSegmentUsersAndWaitForCount uploads included users and waits until
// the upload has been stored. Unlike uploadSegmentUsersAndWait it does not
// list the users, so it also works for a segment too large for one response.
func uploadSegmentUsersAndWaitForCount(
	t *testing.T,
	client featureclient.Client,
	segmentID string,
	userIDs []string,
) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := client.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
		EnvironmentId: holdoutEnv.id,
		SegmentId:     segmentID,
		Data:          []byte(strings.Join(userIDs, "\n") + "\n"),
		State:         featureproto.SegmentUser_INCLUDED,
	})
	require.NoError(t, err)
	// The upload sets the segment to UPLOADING before it returns, and the
	// persister stores the number of users of this upload with SUCEEDED.
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		res, err := client.GetSegment(ctx, &featureproto.GetSegmentRequest{
			Id: segmentID, EnvironmentId: holdoutEnv.id,
		})
		return err == nil &&
			res.Segment.Status == featureproto.Segment_SUCEEDED &&
			res.Segment.IncludedUserCount == int64(len(userIDs))
	}, segmentUploadTimeout, 2*time.Second, "segment upload of %d users not stored", len(userIDs))
}
