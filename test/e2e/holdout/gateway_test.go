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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
	"github.com/bucketeer-io/bucketeer/v2/test/e2e/util"
)

const (
	// A patch after a holdout write arrives through the HOLDOUT_GROUP domain
	// event (seconds); cachePropagationTimeout is the hard upper bound.
	ssePatchTimeout = cachePropagationTimeout
	// How long a stream is watched to prove that no patch carried HOLDOUT.
	sseQuietWindow = 15 * time.Second
)

// lockServing serializes the serving tests (one serving group per
// environment). The unlock is registered before any group cleanup so it runs
// after the group has been stopped and archived.
func lockServing(t *testing.T) {
	t.Helper()
	servingMu.Lock()
	t.Cleanup(servingMu.Unlock)
}

// servingScenario is a tagged flag enrolled in a RUNNING 50% HASH group with
// one held-out and one non-held-out user.
type servingScenario struct {
	tag                 string
	flag                *featureproto.Feature
	group               *featureproto.HoldoutGroup
	heldOut, notHeldOut string
}

func (s *servingScenario) holdoutVariation() string { return variationID(s.flag, "C") }
func (s *servingScenario) onVariation() string      { return variationID(s.flag, "A") }
func (s *servingScenario) offVariation() string     { return variationID(s.flag, "B") }

// newServingScenario creates the flag and the group, enrolls the flag with
// variation C as the holdout variation and leaves the group WAITING so the
// test can observe the transition.
func newServingScenario(t *testing.T, client featureclient.Client) *servingScenario {
	t.Helper()
	s := &servingScenario{tag: newTag(t)}
	s.flag = createEnabledFeature(t, client, s.tag)
	refreshFeatureCache(t, s.tag, s.flag.Id)
	s.group = createHoldoutGroup(t, client, servingHashSpec().request(t))
	addHoldoutFeatures(t, client, s.group.Id, target(s.flag, "C"))
	s.heldOut, s.notHeldOut = pickUsers(t, s.group)
	return s
}

func (s *servingScenario) start(t *testing.T, client featureclient.Client) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	s.group = startHoldoutGroup(t, client, s.group.Id)
	return waitForHoldoutServing(t, s.tag, s.heldOut, s.flag.Id, s.holdoutVariation())
}

func (s *servingScenario) stop(t *testing.T, client featureclient.Client) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	s.group = stopHoldoutGroup(t, client, s.group.Id)
	return waitForHoldoutCleared(t, s.tag, s.heldOut)
}

// RFC §4.3.1 / §4.4: a held-out user receives the holdout variation with
// Reason HOLDOUT and the group id; everyone else is evaluated normally. The
// same result over gRPC, REST (/get_evaluations) and the SSE initial put.
func TestServingBasic(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	s := newServingScenario(t, client)

	// Before the start nobody is held out.
	pre := grpcGetEvaluations(t, s.tag, s.heldOut)
	requireNormalEvaluation(t, pre.Evaluations.Evaluations, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)

	held := s.start(t, client)
	requireHoldoutEvaluation(t, held.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
	assert.NotEqual(t, pre.UserEvaluationsId, held.UserEvaluationsId, "UEID changes on start")

	other := grpcGetEvaluations(t, s.tag, s.notHeldOut)
	requireNormalEvaluation(t, other.Evaluations.Evaluations, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)
	assert.False(t, hasHoldoutReason(other.Evaluations.Evaluations))

	t.Run("REST", func(t *testing.T) {
		resp := util.GetEvaluationsRaw(t, restEvaluationsBody(s.tag, s.heldOut, "", 0), *gatewayAddr, holdoutEnv.clientKeyPath)
		requireHoldoutEvaluation(t, resp.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
		assert.Equal(t, held.UserEvaluationsId, resp.UserEvaluationsID, "REST and gRPC agree on the UEID")
		resp = util.GetEvaluationsRaw(t, restEvaluationsBody(s.tag, s.notHeldOut, "", 0), *gatewayAddr, holdoutEnv.clientKeyPath)
		requireNormalEvaluation(t, resp.Evaluations.Evaluations, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)
	})
	t.Run("SSE initial put", func(t *testing.T) {
		stream := util.ConnectStream(t, *gatewayAddr, holdoutEnv.clientKeyPath, util.StreamBody(s.tag, s.heldOut))
		put := stream.WaitForEvent(t, util.SSEEventPut, timeout, nil)
		requireHoldoutEvaluation(t, put.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
		// The stream hashes the unfiltered flag list into its ID (see
		// evaluateFeaturesForStream), so it is not compared with the poll ID.
		assert.NotEmpty(t, put.UserEvaluationsId)
		stream.Close()
	})
	t.Run("GetEvaluation single flag", func(t *testing.T) {
		// RFC §4.3.1 B8: the single-flag path shares the evaluator.
		resp := grpcGetEvaluation(t, s.tag, s.flag.Id, s.heldOut)
		requireHoldoutEvaluation(t, []*featureproto.Evaluation{resp.Evaluation}, s.flag.Id, s.holdoutVariation(), s.group.Id)
		resp = grpcGetEvaluation(t, s.tag, s.flag.Id, s.notHeldOut)
		requireNormalEvaluation(t, []*featureproto.Evaluation{resp.Evaluation}, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)
		rest := util.GetEvaluation(t, s.tag, s.flag.Id, s.heldOut, *gatewayAddr, holdoutEnv.clientKeyPath)
		requireHoldoutEvaluation(t, []*featureproto.Evaluation{rest.Evaluation}, s.flag.Id, s.holdoutVariation(), s.group.Id)
	})

	cleared := s.stop(t, client)
	requireNormalEvaluation(t, cleared.Evaluations.Evaluations, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)
}

// RFC §4.3.1: holdout is decided before flag state. It wins over a disabled
// flag, an unsatisfied prerequisite and a matching targeting rule, and a
// downstream flag whose prerequisite is a held-out flag is resolved against
// the holdout variation.
func TestServingPrecedence(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	tag := newTag(t)
	disabled := createEnabledFeature(t, client, tag)
	updateFeature(t, client, &featureproto.UpdateFeatureRequest{Id: disabled.Id, Enabled: wrapperspb.Bool(false)})
	ruled := createEnabledFeature(t, client, tag)
	addFixedRule(t, client, ruled.Id, variationID(ruled, "B"), "plan", "gold")
	parent := createEnabledFeature(t, client, tag)
	// child requires parent == A; the holdout moves parent to C.
	child := createEnabledFeature(t, client, tag)
	addPrerequisite(t, client, child.Id, parent.Id, variationID(parent, "A"))
	// gated requires parent == B, which is never satisfied; the holdout
	// still serves its holdout variation.
	gated := createEnabledFeature(t, client, tag)
	addPrerequisite(t, client, gated.Id, parent.Id, variationID(parent, "B"))
	refreshFeatureCache(t, tag, disabled.Id, ruled.Id, parent.Id, child.Id, gated.Id)

	group := createHoldoutGroup(t, client, servingHashSpec().request(t))
	addHoldoutFeatures(t, client, group.Id,
		target(disabled, "C"), target(ruled, "C"), target(parent, "C"), target(gated, "C"))
	heldOut, notHeldOut := pickUsers(t, group)
	gold := &userproto.User{Id: heldOut, Data: map[string]string{"plan": "gold"}}
	goldNotHeld := &userproto.User{Id: notHeldOut, Data: map[string]string{"plan": "gold"}}

	// Normal path first, so the comparison below is against observed values.
	pre := grpcGetEvaluationsWithUser(t, tag, gold)
	requireNormalEvaluation(t, pre.Evaluations.Evaluations, disabled.Id, variationID(disabled, "B"), featureproto.Reason_OFF_VARIATION)
	requireNormalEvaluation(t, pre.Evaluations.Evaluations, ruled.Id, variationID(ruled, "B"), featureproto.Reason_RULE)
	requireNormalEvaluation(t, pre.Evaluations.Evaluations, child.Id, variationID(child, "A"), featureproto.Reason_DEFAULT)
	requireNormalEvaluation(t, pre.Evaluations.Evaluations, gated.Id, variationID(gated, "B"), featureproto.Reason_PREREQUISITE)

	group = startHoldoutGroup(t, client, group.Id)
	var evals []*featureproto.Evaluation
	require.Eventually(t, func() bool {
		evals = grpcGetEvaluationsWithUser(t, tag, gold).Evaluations.Evaluations
		e := findEvaluation(evals, parent.Id)
		return e != nil && e.GetReason().GetType() == featureproto.Reason_HOLDOUT
	}, cachePropagationTimeout, pollInterval, "holdout never reached the gateway")

	requireHoldoutEvaluation(t, evals, disabled.Id, variationID(disabled, "C"), group.Id)
	requireHoldoutEvaluation(t, evals, ruled.Id, variationID(ruled, "C"), group.Id)
	requireHoldoutEvaluation(t, evals, parent.Id, variationID(parent, "C"), group.Id)
	requireHoldoutEvaluation(t, evals, gated.Id, variationID(gated, "C"), group.Id)
	// child is not a target: its prerequisite (parent == A) now fails
	// because parent resolves to the holdout variation C.
	requireNormalEvaluation(t, evals, child.Id, variationID(child, "B"), featureproto.Reason_PREREQUISITE)

	other := grpcGetEvaluationsWithUser(t, tag, goldNotHeld).Evaluations.Evaluations
	assert.False(t, hasHoldoutReason(other))
	requireNormalEvaluation(t, other, disabled.Id, variationID(disabled, "B"), featureproto.Reason_OFF_VARIATION)
	requireNormalEvaluation(t, other, ruled.Id, variationID(ruled, "B"), featureproto.Reason_RULE)
	requireNormalEvaluation(t, other, child.Id, variationID(child, "A"), featureproto.Reason_DEFAULT)

	stopHoldoutGroup(t, client, group.Id)
	require.Eventually(t, func() bool {
		evals = grpcGetEvaluationsWithUser(t, tag, gold).Evaluations.Evaluations
		return !hasHoldoutReason(evals)
	}, cachePropagationTimeout, pollInterval, "holdout never cleared")
	requireNormalEvaluation(t, evals, child.Id, variationID(child, "A"), featureproto.Reason_DEFAULT)
}

// RFC §4.3.2 change detection on the polling path. A client that polls with
// the UEID and evaluatedAt of its last response must be told to replace its
// cache (ForceUpdate=true) after every holdout change, and must get an empty,
// non-forced diff once it has caught up. Exercised on Start, Add, Remove and
// Update of a RUNNING group.
func TestChangeDetectionOnPoll(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	s := newServingScenario(t, client)
	extra := createEnabledFeature(t, client, s.tag)
	refreshFeatureCache(t, s.tag, extra.Id)

	// The client's cached state before any holdout exists.
	pre := grpcGetEvaluations(t, s.tag, s.heldOut)
	assert.True(t, pre.Evaluations.ForceUpdate, "first request without an ID is a full update")
	caughtUp := func(t *testing.T, prev *gatewayproto.GetEvaluationsResponse) {
		t.Helper()
		again := grpcGetEvaluationsByEvaluatedAt(t, s.tag, s.heldOut, prev.UserEvaluationsId, prev.Evaluations.CreatedAt)
		assert.False(t, again.Evaluations.ForceUpdate, "caught-up client must not be forced")
		assert.Empty(t, again.Evaluations.Evaluations)
		assert.Empty(t, again.Evaluations.ArchivedFeatureIds)
		assert.Equal(t, prev.UserEvaluationsId, again.UserEvaluationsId)
	}
	caughtUp(t, pre)

	// Start.
	s.start(t, client)
	afterStart := grpcGetEvaluationsByEvaluatedAt(t, s.tag, s.heldOut, pre.UserEvaluationsId, pre.Evaluations.CreatedAt)
	assert.True(t, afterStart.Evaluations.ForceUpdate, "stale ID after start")
	assert.NotEqual(t, pre.UserEvaluationsId, afterStart.UserEvaluationsId)
	requireHoldoutEvaluation(t, afterStart.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
	caughtUp(t, afterStart)

	// Add a flag while RUNNING.
	addHoldoutFeatures(t, client, s.group.Id, target(extra, "C"))
	var afterAdd *gatewayproto.GetEvaluationsResponse
	require.Eventually(t, func() bool {
		afterAdd = grpcGetEvaluationsByEvaluatedAt(t, s.tag, s.heldOut, afterStart.UserEvaluationsId, afterStart.Evaluations.CreatedAt)
		return findEvaluation(afterAdd.Evaluations.Evaluations, extra.Id) != nil
	}, cachePropagationTimeout, pollInterval, "added flag never reached the gateway")
	assert.True(t, afterAdd.Evaluations.ForceUpdate, "stale ID after add")
	assert.NotEqual(t, afterStart.UserEvaluationsId, afterAdd.UserEvaluationsId)
	requireHoldoutEvaluation(t, afterAdd.Evaluations.Evaluations, extra.Id, variationID(extra, "C"), s.group.Id)
	requireHoldoutEvaluation(t, afterAdd.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
	caughtUp(t, afterAdd)

	// Remove it again.
	removeHoldoutFeatures(t, client, s.group.Id, extra.Id)
	var afterRemove *gatewayproto.GetEvaluationsResponse
	require.Eventually(t, func() bool {
		afterRemove = grpcGetEvaluationsByEvaluatedAt(t, s.tag, s.heldOut, afterAdd.UserEvaluationsId, afterAdd.Evaluations.CreatedAt)
		e := findEvaluation(afterRemove.Evaluations.Evaluations, extra.Id)
		return e != nil && e.GetReason().GetType() != featureproto.Reason_HOLDOUT
	}, cachePropagationTimeout, pollInterval, "removed flag never left the gateway")
	assert.True(t, afterRemove.Evaluations.ForceUpdate, "stale ID after remove")
	assert.NotEqual(t, afterAdd.UserEvaluationsId, afterRemove.UserEvaluationsId)
	requireNormalEvaluation(t, afterRemove.Evaluations.Evaluations, extra.Id, variationID(extra, "A"), featureproto.Reason_DEFAULT)
	requireHoldoutEvaluation(t, afterRemove.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
	caughtUp(t, afterRemove)

	// Update (stop_at) of a RUNNING group bumps updated_at and so the UEID.
	bumpSecond(t, getHoldoutGroup(t, client, s.group.Id))
	updateHoldoutGroup(t, client, &featureproto.UpdateHoldoutGroupRequest{
		Id: s.group.Id, StopAt: wrapperspb.Int64(s.group.StopAt + 600),
	})
	var afterUpdate *gatewayproto.GetEvaluationsResponse
	require.Eventually(t, func() bool {
		afterUpdate = grpcGetEvaluationsByEvaluatedAt(t, s.tag, s.heldOut, afterRemove.UserEvaluationsId, afterRemove.Evaluations.CreatedAt)
		return afterUpdate.UserEvaluationsId != afterRemove.UserEvaluationsId
	}, cachePropagationTimeout, pollInterval, "update never reached the gateway")
	assert.True(t, afterUpdate.Evaluations.ForceUpdate, "stale ID after update")
	requireHoldoutEvaluation(t, afterUpdate.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
	caughtUp(t, afterUpdate)

	t.Run("REST polling agrees", func(t *testing.T) {
		stale := util.GetEvaluationsRaw(t,
			restEvaluationsBody(s.tag, s.heldOut, pre.UserEvaluationsId, pre.Evaluations.CreatedAt),
			*gatewayAddr, holdoutEnv.clientKeyPath)
		assert.True(t, stale.Evaluations.ForceUpdate)
		requireHoldoutEvaluation(t, stale.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
		fresh := util.GetEvaluationsRaw(t,
			restEvaluationsBody(s.tag, s.heldOut, stale.UserEvaluationsID, stale.Evaluations.CreatedAt),
			*gatewayAddr, holdoutEnv.clientKeyPath)
		assert.False(t, fresh.Evaluations.ForceUpdate)
		assert.Empty(t, fresh.Evaluations.Evaluations)
	})

	s.stop(t, client)
}

// RFC §4.3.2 stop and clearance. A client holding an ID from the running
// period is forced to a full update that carries no HOLDOUT, and a fresh
// client is not forced on its second poll. Covered for a manual stop, a
// scheduled stop, and a second group cycle in the same environment.
func TestStopAndClearance(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	s := newServingScenario(t, client)

	assertCleared := func(t *testing.T, running *gatewayproto.GetEvaluationsResponse) {
		t.Helper()
		stale := grpcGetEvaluationsByEvaluatedAt(t, s.tag, s.heldOut, running.UserEvaluationsId, running.Evaluations.CreatedAt)
		assert.True(t, stale.Evaluations.ForceUpdate, "running-period ID must be replaced")
		assert.NotEqual(t, running.UserEvaluationsId, stale.UserEvaluationsId)
		assert.False(t, hasHoldoutReason(stale.Evaluations.Evaluations))
		requireNormalEvaluation(t, stale.Evaluations.Evaluations, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)

		// A client that first connects after the stop: its second poll,
		// chained from its own response, is not forced either.
		freshUser := newUserID(t)
		fresh := grpcGetEvaluations(t, s.tag, freshUser)
		assert.False(t, hasHoldoutReason(fresh.Evaluations.Evaluations))
		second := grpcGetEvaluationsByEvaluatedAt(t, s.tag, freshUser, fresh.UserEvaluationsId, fresh.Evaluations.CreatedAt)
		assert.False(t, second.Evaluations.ForceUpdate, "fresh client is not forced after clearance")
		assert.Empty(t, second.Evaluations.Evaluations)
		assert.Equal(t, fresh.UserEvaluationsId, second.UserEvaluationsId)
		// The previously held-out client, once caught up, is not forced.
		caughtUp := grpcGetEvaluationsByEvaluatedAt(t, s.tag, s.heldOut, stale.UserEvaluationsId, stale.Evaluations.CreatedAt)
		assert.False(t, caughtUp.Evaluations.ForceUpdate)
		assert.Empty(t, caughtUp.Evaluations.Evaluations)
	}

	t.Run("manual stop", func(t *testing.T) {
		running := s.start(t, client)
		s.stop(t, client)
		assertCleared(t, running)
	})

	t.Run("second group cycle", func(t *testing.T) {
		second := createHoldoutGroup(t, client, servingHashSpec().request(t))
		addHoldoutFeatures(t, client, second.Id, target(s.flag, "B"))
		heldOut, _ := pickUsers(t, second)
		before := grpcGetEvaluations(t, s.tag, heldOut)
		assert.False(t, hasHoldoutReason(before.Evaluations.Evaluations))
		second = startHoldoutGroup(t, client, second.Id)
		running := waitForHoldoutServing(t, s.tag, heldOut, s.flag.Id, s.offVariation())
		requireHoldoutEvaluation(t, running.Evaluations.Evaluations, s.flag.Id, s.offVariation(), second.Id)
		stale := grpcGetEvaluationsByEvaluatedAt(t, s.tag, heldOut, before.UserEvaluationsId, before.Evaluations.CreatedAt)
		assert.True(t, stale.Evaluations.ForceUpdate)
		stopHoldoutGroup(t, client, second.Id)
		waitForHoldoutCleared(t, s.tag, heldOut)
		after := grpcGetEvaluationsByEvaluatedAt(t, s.tag, heldOut, running.UserEvaluationsId, running.Evaluations.CreatedAt)
		assert.True(t, after.Evaluations.ForceUpdate)
		assert.False(t, hasHoldoutReason(after.Evaluations.Evaluations))
	})

	t.Run("scheduled stop", func(t *testing.T) {
		// Both transitions are taken by the HoldoutGroupStatusUpdater, run
		// explicitly here (the dev cron runs the same job every minute and
		// can only take the same step; see TestLifecycleScheduled). The
		// window is opened wide so cache propagation (up to 150 s) cannot
		// outlive it; stop_at is pulled in only once HOLDOUT is observed.
		start, stop := time.Now().Unix(), time.Now().Add(servingWindow).Unix()
		spec := groupSpec{assignment: featureproto.HoldoutGroup_HASH, percentage: 50, startAt: start, stopAt: stop}
		third := createHoldoutGroup(t, client, spec.request(t))
		addHoldoutFeatures(t, client, third.Id, target(s.flag, "C"))
		heldOut, _ := pickUsers(t, third)
		runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
		third = getHoldoutGroup(t, client, third.Id)
		require.Equal(t, featureproto.HoldoutGroup_RUNNING, third.Status)
		running := waitForHoldoutServing(t, s.tag, heldOut, s.flag.Id, s.holdoutVariation())
		requireHoldoutEvaluation(t, running.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), third.Id)

		stopAt := time.Now().Unix() + 5
		updateHoldoutGroup(t, client, &featureproto.UpdateHoldoutGroupRequest{Id: third.Id, StopAt: wrapperspb.Int64(stopAt)})
		waitUntilUnix(t, stopAt)
		runBatchJob(t, btproto.BatchJob_HoldoutGroupStatusUpdater)
		got := getHoldoutGroup(t, client, third.Id)
		require.Equal(t, featureproto.HoldoutGroup_STOPPED, got.Status)
		waitForHoldoutCleared(t, s.tag, heldOut)
		stale := grpcGetEvaluationsByEvaluatedAt(t, s.tag, heldOut, running.UserEvaluationsId, running.Evaluations.CreatedAt)
		assert.True(t, stale.Evaluations.ForceUpdate)
		assert.False(t, hasHoldoutReason(stale.Evaluations.Evaluations))
		second := grpcGetEvaluationsByEvaluatedAt(t, s.tag, heldOut, stale.UserEvaluationsId, stale.Evaluations.CreatedAt)
		assert.False(t, second.Evaluations.ForceUpdate)
		assert.Empty(t, second.Evaluations.Evaluations)
	})
}

// RFC §4.2.2 SEGMENT mode on the gateway: membership comes from the
// uploaded user list even though no flag rule references the segment, and
// the segment is locked while the group serves (RFC §4.5) and released after
// the stop.
func TestServingSegmentMode(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	tag := newTag(t)
	flag := createEnabledFeature(t, client, tag)
	refreshFeatureCache(t, tag, flag.Id)
	segment := createSegment(t, client)
	member, nonMember := newUserID(t), newUserID(t)
	uploadSegmentUsersAndWait(t, client, segment.Id, []string{member})

	start, stop := servingWindowNow()
	group := createHoldoutGroup(t, client, groupSpec{
		assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, startAt: start, stopAt: stop,
	}.request(t))
	addHoldoutFeatures(t, client, group.Id, target(flag, "C"))
	group = startHoldoutGroup(t, client, group.Id)

	running := waitForHoldoutServing(t, tag, member, flag.Id, variationID(flag, "C"))
	requireHoldoutEvaluation(t, running.Evaluations.Evaluations, flag.Id, variationID(flag, "C"), group.Id)
	other := grpcGetEvaluations(t, tag, nonMember)
	requireNormalEvaluation(t, other.Evaluations.Evaluations, flag.Id, variationID(flag, "A"), featureproto.Reason_DEFAULT)

	t.Run("segment is locked while serving", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.UpdateSegment(ctx, &featureproto.UpdateSegmentRequest{
			Id: segment.Id, EnvironmentId: holdoutEnv.id, Name: wrapperspb.String(segment.Name + "-renamed"),
		})
		requireCode(t, err, codes.FailedPrecondition, "UpdateSegment")
		_, err = client.DeleteSegment(ctx, &featureproto.DeleteSegmentRequest{Id: segment.Id, EnvironmentId: holdoutEnv.id})
		requireCode(t, err, codes.FailedPrecondition, "DeleteSegment")
		_, err = client.BulkUploadSegmentUsers(ctx, &featureproto.BulkUploadSegmentUsersRequest{
			EnvironmentId: holdoutEnv.id, SegmentId: segment.Id, Data: []byte(nonMember + "\n"),
			State: featureproto.SegmentUser_INCLUDED,
		})
		requireCode(t, err, codes.FailedPrecondition, "BulkUploadSegmentUsers")
		// Another group cannot reuse the locked segment.
		s2, e2 := farFutureWindow()
		_, err = client.CreateHoldoutGroup(ctx, groupSpec{
			assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, startAt: s2, stopAt: e2,
		}.request(t))
		requireCode(t, err, codes.FailedPrecondition, "CreateHoldoutGroup with a locked segment")
	})

	stopHoldoutGroup(t, client, group.Id)
	waitForHoldoutCleared(t, tag, member)

	t.Run("segment is released after the stop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_, err := client.UpdateSegment(ctx, &featureproto.UpdateSegmentRequest{
			Id: segment.Id, EnvironmentId: holdoutEnv.id, Name: wrapperspb.String(segment.Name + "-renamed"),
		})
		require.NoError(t, err)
		s2, e2 := farFutureWindow()
		other := createHoldoutGroup(t, client, groupSpec{
			assignment: featureproto.HoldoutGroup_SEGMENT, segmentID: segment.Id, startAt: s2, stopAt: e2,
		}.request(t))
		assert.Equal(t, segment.Id, other.SegmentId)
	})
}

// RFC §4.3.1: tag filtering is unchanged by holdouts. A request for a tag
// only sees the flags carrying it, held out or not, and the UEID of an
// unrelated tag does not depend on the group's targets.
func TestServingTagFiltering(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	tagA, tagB := newTag(t), newTag(t)
	flagA := createEnabledFeature(t, client, tagA)
	flagB := createEnabledFeature(t, client, tagB)
	refreshFeatureCache(t, tagA, flagA.Id)
	refreshFeatureCache(t, tagB, flagB.Id)

	group := createHoldoutGroup(t, client, servingHashSpec().request(t))
	addHoldoutFeatures(t, client, group.Id, target(flagA, "C"), target(flagB, "C"))
	heldOut, _ := pickUsers(t, group)
	group = startHoldoutGroup(t, client, group.Id)

	forA := waitForHoldoutServing(t, tagA, heldOut, flagA.Id, variationID(flagA, "C"))
	assert.Nil(t, findEvaluation(forA.Evaluations.Evaluations, flagB.Id), "tag A must not see flag B")
	forB := waitForHoldoutServing(t, tagB, heldOut, flagB.Id, variationID(flagB, "C"))
	assert.Nil(t, findEvaluation(forB.Evaluations.Evaluations, flagA.Id), "tag B must not see flag A")
	assert.NotEqual(t, forA.UserEvaluationsId, forB.UserEvaluationsId)

	// No tag: every flag of the environment.
	all := grpcGetEvaluations(t, "", heldOut)
	requireHoldoutEvaluation(t, all.Evaluations.Evaluations, flagA.Id, variationID(flagA, "C"), group.Id)
	requireHoldoutEvaluation(t, all.Evaluations.Evaluations, flagB.Id, variationID(flagB, "C"), group.Id)

	stopHoldoutGroup(t, client, group.Id)
	waitForHoldoutCleared(t, tagA, heldOut)
	waitForHoldoutCleared(t, tagB, heldOut)
}

// RFC §4.3.2 / §4.4 on the SSE stream: a held-out client receives a patch
// with HOLDOUT on Start and a patch restoring the normal variation on Stop; a
// non-held-out client never sees HOLDOUT; the stream is never closed with an
// INTERNAL error; a reconnect carrying a pre-start ID gets a forced put
// (mirrors TestStreamReconnect in test/e2e/gateway/stream_test.go).
func TestServingStream(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	s := newServingScenario(t, client)

	held := util.ConnectStream(t, *gatewayAddr, holdoutEnv.clientKeyPath, util.StreamBody(s.tag, s.heldOut))
	other := util.ConnectStream(t, *gatewayAddr, holdoutEnv.clientKeyPath, util.StreamBody(s.tag, s.notHeldOut))
	prePut := held.WaitForEvent(t, util.SSEEventPut, timeout, nil)
	requireNormalEvaluation(t, prePut.Evaluations.Evaluations, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)
	otherPut := other.WaitForEvent(t, util.SSEEventPut, timeout, nil)
	assert.False(t, hasHoldoutReason(otherPut.Evaluations.Evaluations))

	s.group = startHoldoutGroup(t, client, s.group.Id)
	startPatch := held.WaitForEvent(t, util.SSEEventPatch, ssePatchTimeout, func(ev *gatewayproto.StreamEvaluationsEvent) bool {
		return hasHoldoutReason(ev.Evaluations.Evaluations)
	})
	requireHoldoutEvaluation(t, startPatch.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
	assert.NotEqual(t, prePut.UserEvaluationsId, startPatch.UserEvaluationsId)
	for _, p := range other.DrainPatches(t, sseQuietWindow) {
		assert.False(t, hasHoldoutReason(p.Evaluations.Evaluations), "non-held-out stream saw HOLDOUT")
	}

	t.Run("reconnect with a pre-start ID is forced", func(t *testing.T) {
		resumed := util.ConnectStream(t, *gatewayAddr, holdoutEnv.clientKeyPath,
			util.StreamResumeBody(s.tag, s.heldOut, prePut.UserEvaluationsId, prePut.Evaluations.CreatedAt))
		put := resumed.WaitForEvent(t, util.SSEEventPut, timeout, nil)
		assert.True(t, put.Evaluations.ForceUpdate)
		requireHoldoutEvaluation(t, put.Evaluations.Evaluations, s.flag.Id, s.holdoutVariation(), s.group.Id)
		resumed.Close()
	})

	s.group = stopHoldoutGroup(t, client, s.group.Id)
	stopPatch := held.WaitForEvent(t, util.SSEEventPatch, ssePatchTimeout, func(ev *gatewayproto.StreamEvaluationsEvent) bool {
		e := findEvaluation(ev.Evaluations.Evaluations, s.flag.Id)
		return e != nil && e.GetReason().GetType() != featureproto.Reason_HOLDOUT
	})
	requireNormalEvaluation(t, stopPatch.Evaluations.Evaluations, s.flag.Id, s.onVariation(), featureproto.Reason_DEFAULT)
	assert.True(t, stopPatch.Evaluations.ForceUpdate, "stop patch replaces the cached holdout assignment")
	for _, p := range other.DrainPatches(t, sseQuietWindow) {
		assert.False(t, hasHoldoutReason(p.Evaluations.Evaluations), "non-held-out stream saw HOLDOUT")
	}
	assert.Empty(t, held.Errors, "held-out stream received an error event")
	assert.Empty(t, other.Errors, "non-held-out stream received an error event")
	held.Close()
	other.Close()
}

// RFC §4.4 server SDK path (phase 2 ships no holdout data to server SDKs):
// GetFeatureFlags and GetSegmentUsers keep working with the server key, the
// flag is byte-for-byte unchanged by the holdout, and the feature_flags_id
// is stable across Start and Stop.
func TestServerSDKUnaffected(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	s := newServingScenario(t, client)
	// A segment is delivered to server SDKs only while a flag references it.
	segment := createSegment(t, client)
	uploadSegmentUsersAndWait(t, client, segment.Id, []string{newUserID(t)})
	segmented := createEnabledFeature(t, client, s.tag)
	addSegmentRule(t, client, segmented.Id, variationID(segmented, "B"), segment.Id)
	refreshFeatureCache(t, s.tag, segmented.Id)

	var before *gatewayproto.GetFeatureFlagsResponse
	require.Eventually(t, func() bool {
		before = grpcGetFeatureFlags(t, s.tag, "", 0)
		return findFeature(before.Features, s.flag.Id) != nil
	}, cachePropagationTimeout, pollInterval, "flag not visible to the server SDK")
	flagBefore := findFeature(before.Features, s.flag.Id)
	assert.True(t, before.ForceUpdate)
	assert.NotEmpty(t, before.FeatureFlagsId)

	s.start(t, client)
	during := grpcGetFeatureFlags(t, s.tag, before.FeatureFlagsId, before.RequestedAt)
	assert.Equal(t, before.FeatureFlagsId, during.FeatureFlagsId, "ffID is stable across Start")
	assert.False(t, during.ForceUpdate)
	assert.Empty(t, during.Features, "no flag diff on Start")
	assert.Empty(t, during.ArchivedFeatureFlagIds)
	full := grpcGetFeatureFlags(t, s.tag, "", 0)
	assert.True(t, proto.Equal(flagBefore, findFeature(full.Features, s.flag.Id)), "flag unchanged while serving")

	// Segment users reach the server SDK path through the SegmentUserCacher
	// cron (1 min) and the gateway L1 TTL (1 min). Other segments of the
	// environment may be delivered too; only this one is checked.
	var delivered *featureproto.SegmentUsers
	require.Eventually(t, func() bool {
		for _, su := range grpcGetSegmentUsers(t, []string{segment.Id}, 0).SegmentUsers {
			if su.SegmentId == segment.Id {
				delivered = su
				return len(su.Users) == 1
			}
		}
		return false
	}, cachePropagationTimeout, pollInterval, "segment users not visible to the server SDK")
	assert.Equal(t, segment.Id, delivered.SegmentId)

	s.stop(t, client)
	after := grpcGetFeatureFlags(t, s.tag, before.FeatureFlagsId, before.RequestedAt)
	assert.Equal(t, before.FeatureFlagsId, after.FeatureFlagsId, "ffID is stable across Stop")
	assert.False(t, after.ForceUpdate)
	assert.Empty(t, after.Features)
}

func findFeature(features []*featureproto.Feature, id string) *featureproto.Feature {
	for _, f := range features {
		if f.Id == id {
			return f
		}
	}
	return nil
}
