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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

// The tests in this file play a locally evaluating server SDK: they keep the
// feature_flags_id and requested_at of the last response as the cursor of the
// next poll, and evaluate with the evaluation package the Go SDK uses.

// requireSameAsServer evaluates the flag locally from a full GetFeatureFlags
// response and requires the variation, the reason and the holdout group of
// the server-evaluated GetEvaluations result for the same user.
func requireSameAsServer(
	t *testing.T,
	tag, featureID, userID string,
	full *gatewayproto.GetFeatureFlagsResponse,
	segmentUsers map[string][]*featureproto.SegmentUser,
) *featureproto.Evaluation {
	t.Helper()
	local, err := evaluation.NewEvaluator().EvaluateFeatures(
		full.Features, &userproto.User{Id: userID}, segmentUsers, nil, tag, full.HoldoutGroups...,
	)
	require.NoError(t, err)
	localEval := findEvaluation(local.Evaluations, featureID)
	require.NotNil(t, localEval, "flag %s not evaluated locally", featureID)
	serverEval := findEvaluation(grpcGetEvaluations(t, tag, userID).Evaluations.Evaluations, featureID)
	require.NotNil(t, serverEval, "flag %s not evaluated by the gateway", featureID)
	assert.Equal(t, serverEval.VariationId, localEval.VariationId, "variation for %s", userID)
	assert.Equal(t, serverEval.GetReason().GetType(), localEval.GetReason().GetType(), "reason for %s", userID)
	assert.Equal(t,
		serverEval.GetReason().GetHoldoutGroupId(), localEval.GetReason().GetHoldoutGroupId(),
		"holdout group for %s", userID,
	)
	return localEval
}

// RFC §4.3.2 local-evaluation path: GetFeatureFlags with the server key
// returns holdout_groups for a serving group and an empty list otherwise;
// feature_flags_id changes on Start, Add, Remove and Stop; none of them forces
// a full resend, and a caught-up SDK gets an empty non-forced diff; the
// delivered definition holds out the same users as GetEvaluations.
func TestServerSDKHoldoutDelivery(t *testing.T) {
	lockServing(t)
	client := newFeatureClient(t)
	s := newServingScenario(t, client)
	extra := createEnabledFeature(t, client, s.tag)
	refreshFeatureCache(t, s.tag, extra.Id)

	// The SDK's state before any group serves. The previous test's group may
	// still be leaving a replica, so wait for an empty list.
	var before *gatewayproto.GetFeatureFlagsResponse
	require.Eventually(t, func() bool {
		before = grpcGetFeatureFlags(t, s.tag, "", 0)
		return findFeature(before.Features, s.flag.Id) != nil &&
			findFeature(before.Features, extra.Id) != nil &&
			len(before.HoldoutGroups) == 0
	}, cachePropagationTimeout, pollInterval, "flags not visible to the server SDK without a holdout group")
	assert.True(t, before.ForceUpdate)
	assert.NotEmpty(t, before.FeatureFlagsId)
	flagBefore := findFeature(before.Features, s.flag.Id)
	idle := waitFeatureFlagsPollCaughtUp(t, s.tag, before)
	assert.Empty(t, idle.HoldoutGroups, "a WAITING group is not delivered")

	// Start.
	s.start(t, client)
	afterStart := waitFeatureFlagsPoll(t, s.tag, before, "the started group never reached the server SDK",
		func(r *gatewayproto.GetFeatureFlagsResponse) bool {
			return len(r.HoldoutGroups) == 1 && r.HoldoutGroups[0].Id == s.group.Id
		})
	assert.NotEqual(t, before.FeatureFlagsId, afterStart.FeatureFlagsId, "Start changes feature_flags_id")
	assert.False(t, afterStart.ForceUpdate, "a holdout change does not force a full resend")
	delivered := afterStart.HoldoutGroups[0]
	assert.Equal(t, holdoutEnv.id, delivered.EnvironmentId)
	assert.Equal(t, featureproto.HoldoutGroup_RUNNING, delivered.Status)
	assert.Equal(t, featureproto.HoldoutGroup_HASH, delivered.AssignmentType)
	assert.Equal(t, s.group.SamplingSeed, delivered.SamplingSeed)
	assert.Equal(t, s.group.BucketStart, delivered.BucketStart)
	assert.Equal(t, s.group.Percentage, delivered.Percentage)
	assert.Equal(t, map[string]string{s.flag.Id: s.holdoutVariation()}, holdoutTargets(delivered))
	// Only what local evaluation reads leaves the server.
	assert.Empty(t, delivered.Name)
	assert.Empty(t, delivered.Description)
	assert.Empty(t, delivered.CreatedBy)
	assert.Empty(t, delivered.UpdatedBy)
	caughtUp := waitFeatureFlagsPollCaughtUp(t, s.tag, afterStart)
	require.Len(t, caughtUp.HoldoutGroups, 1, "the groups are listed even when nothing changed")
	assert.True(t, proto.Equal(delivered, caughtUp.HoldoutGroups[0]))

	// The delivered definition holds out the users GetEvaluations holds out.
	assert.True(t, evaluation.IsHeldOut(delivered, &userproto.User{Id: s.heldOut}, nil))
	assert.False(t, evaluation.IsHeldOut(delivered, &userproto.User{Id: s.notHeldOut}, nil))
	var full *gatewayproto.GetFeatureFlagsResponse
	require.Eventually(t, func() bool {
		full = grpcGetFeatureFlags(t, s.tag, "", 0)
		return len(full.HoldoutGroups) == 1 && full.HoldoutGroups[0].Id == s.group.Id
	}, cachePropagationTimeout, pollInterval, "a first request never carried the serving group")
	assert.True(t, full.ForceUpdate)
	assert.Equal(t, afterStart.FeatureFlagsId, full.FeatureFlagsId)
	assert.True(t, proto.Equal(flagBefore, findFeature(full.Features, s.flag.Id)), "the flag is unchanged while serving")
	heldOutEval := requireSameAsServer(t, s.tag, s.flag.Id, s.heldOut, full, nil)
	assert.Equal(t, s.holdoutVariation(), heldOutEval.VariationId)
	assert.Equal(t, featureproto.Reason_HOLDOUT, heldOutEval.GetReason().GetType())
	assert.Equal(t, s.group.Id, heldOutEval.GetReason().GetHoldoutGroupId())
	otherEval := requireSameAsServer(t, s.tag, s.flag.Id, s.notHeldOut, full, nil)
	assert.Equal(t, s.onVariation(), otherEval.VariationId)
	assert.Equal(t, featureproto.Reason_DEFAULT, otherEval.GetReason().GetType())

	// Add a flag while RUNNING.
	addHoldoutFeatures(t, client, s.group.Id, target(extra, "C"))
	afterAdd := waitFeatureFlagsPoll(t, s.tag, afterStart, "the added target never reached the server SDK",
		func(r *gatewayproto.GetFeatureFlagsResponse) bool {
			return len(r.HoldoutGroups) == 1 && len(r.HoldoutGroups[0].Features) == 2
		})
	assert.NotEqual(t, afterStart.FeatureFlagsId, afterAdd.FeatureFlagsId, "Add changes feature_flags_id")
	assert.False(t, afterAdd.ForceUpdate, "a holdout change does not force a full resend")
	assert.Equal(t,
		map[string]string{s.flag.Id: s.holdoutVariation(), extra.Id: variationID(extra, "C")},
		holdoutTargets(afterAdd.HoldoutGroups[0]),
	)
	waitFeatureFlagsPollCaughtUp(t, s.tag, afterAdd)

	// Remove it again.
	removeHoldoutFeatures(t, client, s.group.Id, extra.Id)
	afterRemove := waitFeatureFlagsPoll(t, s.tag, afterAdd, "the removed target never left the server SDK",
		func(r *gatewayproto.GetFeatureFlagsResponse) bool {
			return len(r.HoldoutGroups) == 1 && len(r.HoldoutGroups[0].Features) == 1
		})
	assert.NotEqual(t, afterAdd.FeatureFlagsId, afterRemove.FeatureFlagsId, "Remove changes feature_flags_id")
	assert.False(t, afterRemove.ForceUpdate, "a holdout change does not force a full resend")
	assert.Equal(t,
		map[string]string{s.flag.Id: s.holdoutVariation()},
		holdoutTargets(afterRemove.HoldoutGroups[0]),
	)
	waitFeatureFlagsPollCaughtUp(t, s.tag, afterRemove)

	// Stop.
	s.stop(t, client)
	afterStop := waitFeatureFlagsPoll(t, s.tag, afterRemove, "the stopped group never left the server SDK",
		func(r *gatewayproto.GetFeatureFlagsResponse) bool {
			return len(r.HoldoutGroups) == 0 && r.FeatureFlagsId != afterRemove.FeatureFlagsId
		})
	assert.False(t, afterStop.ForceUpdate, "a holdout change does not force a full resend")
	stopped := waitFeatureFlagsPollCaughtUp(t, s.tag, afterStop)
	assert.Empty(t, stopped.HoldoutGroups)
}

// RFC §4.3.1 / §4.3.2: a SEGMENT group is delivered with its segment_id, and
// GetSegmentUsers delivers that segment although no flag rule references it,
// so the SDK can decide membership locally.
func TestServerSDKHoldoutSegment(t *testing.T) {
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

	// No flag references the segment, so a first request (a full resend) does
	// not carry it while the group is WAITING.
	assert.Nil(t, findSegmentUsers(grpcGetSegmentUsers(t, nil, 0).SegmentUsers, segment.Id),
		"the segment of a WAITING group is not delivered")

	group = startHoldoutGroup(t, client, group.Id)
	waitForHoldoutServing(t, tag, member, flag.Id, variationID(flag, "C"))

	var full *gatewayproto.GetFeatureFlagsResponse
	require.Eventually(t, func() bool {
		full = grpcGetFeatureFlags(t, tag, "", 0)
		return findFeature(full.Features, flag.Id) != nil &&
			len(full.HoldoutGroups) == 1 && full.HoldoutGroups[0].Id == group.Id
	}, cachePropagationTimeout, pollInterval, "the serving group never reached the server SDK")
	delivered := full.HoldoutGroups[0]
	assert.Equal(t, featureproto.HoldoutGroup_SEGMENT, delivered.AssignmentType)
	assert.Equal(t, segment.Id, delivered.SegmentId)
	assert.Equal(t, map[string]string{flag.Id: variationID(flag, "C")}, holdoutTargets(delivered))

	// An SDK that holds no segment and whose cursor is already past the last
	// edit of the segment: the diff rule alone would never send it.
	cursor := grpcGetSegmentUsers(t, nil, 0).RequestedAt
	var segmentUsers *featureproto.SegmentUsers
	var afterStart *gatewayproto.GetSegmentUsersResponse
	require.Eventually(t, func() bool {
		afterStart = grpcGetSegmentUsers(t, nil, cursor)
		segmentUsers = findSegmentUsers(afterStart.SegmentUsers, segment.Id)
		return segmentUsers != nil
	}, cachePropagationTimeout, pollInterval, "the holdout segment never reached the server SDK")
	assert.False(t, afterStart.ForceUpdate)
	require.Len(t, segmentUsers.Users, 1)
	assert.Equal(t, member, segmentUsers.Users[0].UserId)
	assert.Equal(t, featureproto.SegmentUser_INCLUDED, segmentUsers.Users[0].State)

	// Membership from the delivered definition matches GetEvaluations.
	assert.True(t, evaluation.IsHeldOut(delivered, &userproto.User{Id: member}, segmentUsers.Users))
	assert.False(t, evaluation.IsHeldOut(delivered, &userproto.User{Id: nonMember}, segmentUsers.Users))
	local := map[string][]*featureproto.SegmentUser{segment.Id: segmentUsers.Users}
	memberEval := requireSameAsServer(t, tag, flag.Id, member, full, local)
	assert.Equal(t, variationID(flag, "C"), memberEval.VariationId)
	assert.Equal(t, featureproto.Reason_HOLDOUT, memberEval.GetReason().GetType())
	assert.Equal(t, group.Id, memberEval.GetReason().GetHoldoutGroupId())
	otherEval := requireSameAsServer(t, tag, flag.Id, nonMember, full, local)
	assert.Equal(t, variationID(flag, "A"), otherEval.VariationId)
	assert.Equal(t, featureproto.Reason_DEFAULT, otherEval.GetReason().GetType())

	// Once the SDK holds the segment it is neither resent nor deleted.
	require.Eventually(t, func() bool {
		again := grpcGetSegmentUsers(t, []string{segment.Id}, afterStart.RequestedAt)
		return findSegmentUsers(again.SegmentUsers, segment.Id) == nil &&
			!containsString(again.DeletedSegmentIds, segment.Id)
	}, cachePropagationTimeout, pollInterval, "the holdout segment kept being resent or was deleted while serving")

	// After the stop nothing references the segment, so the SDK drops it.
	stopHoldoutGroup(t, client, group.Id)
	waitForHoldoutCleared(t, tag, member)
	require.Eventually(t, func() bool {
		after := grpcGetSegmentUsers(t, []string{segment.Id}, afterStart.RequestedAt)
		return containsString(after.DeletedSegmentIds, segment.Id)
	}, cachePropagationTimeout, pollInterval, "the holdout segment was never reported as deleted after the stop")
	require.Eventually(t, func() bool {
		return len(grpcGetFeatureFlags(t, tag, "", 0).HoldoutGroups) == 0
	}, cachePropagationTimeout, pollInterval, "the stopped group never left the server SDK")
}

func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
