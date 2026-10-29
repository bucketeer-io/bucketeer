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

package gateway

// Regression guards for environments that never had a holdout group
// (RFC 0048 §4.3.2, §4.7). The shared e2e environment is such an
// environment: nothing in test/e2e creates a holdout group here, the holdout
// suite runs in its own environment (test/e2e/holdout). These tests pin the
// behaviour that #2881 broke on the gateway: a caught-up client must see
// ForceUpdate=false with an empty diff, archived IDs must not leak into
// diffs, and the SSE stream must stay open.
//
// They complement the existing tests, which exercise the same endpoints
// with a synthetic or stale ID rather than the ID the gateway just issued:
//   - TestGrpcGetEvaluationsByEvaluatedAt (diff with a stale ID)
//   - TestGrpcGetEvaluationsWithEmptyUserEvaluationsID (full update)
//   - TestGrpcGetFeatureFlagsWithRequestedAt (diff after a new flag)
//   - TestStreamHeartbeat (first heartbeat), TestStreamPatchOnFlagChange

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	"github.com/bucketeer-io/bucketeer/v2/test/e2e/util"
)

// Two consecutive polls carrying the ID and evaluatedAt the gateway just
// returned: the second one is an empty, non-forced diff with no archived IDs.
func TestGrpcGetEvaluationsNoHoldoutSecondPollIsEmpty(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)

	first := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, "", 0, false)
	require.NotNil(t, first.Evaluations)
	require.True(t, contains(first.Evaluations.Evaluations, featureID))
	assert.True(t, first.Evaluations.ForceUpdate, "a client without an ID gets a full update")
	assert.NotEmpty(t, first.UserEvaluationsId)
	assertNoHoldoutReason(t, first.Evaluations.Evaluations)

	second := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, first.UserEvaluationsId, first.Evaluations.CreatedAt, false)
	assert.False(t, second.Evaluations.ForceUpdate, "caught-up client must not be forced")
	assert.Empty(t, second.Evaluations.Evaluations, "caught-up client must get an empty diff")
	assert.Empty(t, second.Evaluations.ArchivedFeatureIds, "no archived IDs without a change")
	assert.Equal(t, first.UserEvaluationsId, second.UserEvaluationsId, "the ID is stable without a change")
}

// Same as above over REST (/get_evaluations).
func TestGetEvaluationsNoHoldoutSecondPollIsEmpty(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)

	// An (empty) user_evaluation_condition selects the evaluatedAt path that
	// current SDKs use; without it the gateway takes the legacy path, which
	// does not set force_update.
	body := map[string]interface{}{
		"tag":                       tag,
		"user":                      map[string]interface{}{"id": userID},
		"user_evaluation_condition": map[string]interface{}{},
	}
	first := util.GetEvaluationsRaw(t, body, *gatewayAddr, *apiKeyPath)
	require.NotNil(t, first.Evaluations)
	require.True(t, contains(first.Evaluations.Evaluations, featureID))
	assert.True(t, first.Evaluations.ForceUpdate)
	assertNoHoldoutReason(t, first.Evaluations.Evaluations)

	body["user_evaluations_id"] = first.UserEvaluationsID
	body["user_evaluation_condition"] = map[string]interface{}{
		"evaluated_at":            fmt.Sprintf("%d", first.Evaluations.CreatedAt),
		"user_attributes_updated": false,
	}
	second := util.GetEvaluationsRaw(t, body, *gatewayAddr, *apiKeyPath)
	assert.False(t, second.Evaluations.ForceUpdate)
	assert.Empty(t, second.Evaluations.Evaluations)
	assert.Empty(t, second.Evaluations.ArchivedFeatureIds)
	assert.Equal(t, first.UserEvaluationsID, second.UserEvaluationsID)
}

// A server SDK that polls with the feature_flags_id and requested_at it
// just received gets an empty, non-forced diff.
func TestGrpcGetFeatureFlagsNoHoldoutSecondPollIsEmpty(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)

	var first *gatewayproto.GetFeatureFlagsResponse
	require.Eventually(t, func() bool {
		first = grpcGetFeatureFlags(t, tag, "", 0)
		return findFeatureByID(t, featureID, first.Features)
	}, 30*time.Second, 2*time.Second, "tagged feature should be visible")
	assert.True(t, first.ForceUpdate)
	assert.NotEmpty(t, first.FeatureFlagsId)

	second := grpcGetFeatureFlags(t, tag, first.FeatureFlagsId, first.RequestedAt)
	assert.False(t, second.ForceUpdate)
	assert.Empty(t, second.Features)
	assert.Empty(t, second.ArchivedFeatureFlagIds)
	assert.Equal(t, first.FeatureFlagsId, second.FeatureFlagsId)
}

// The stream stays open across a flag update and a two-minute heartbeat
// window: no error event, no close. #2881 closed streams with INTERNAL
// "holdout snapshot unavailable" in environments without a holdout.
func TestStreamNoHoldoutStaysOpen(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	req := newCreateFeatureReq(featureID)
	createFeature(t, client, req)
	addTag(t, tag, featureID, client)
	updateFeatueFlagCache(t)

	stream := connectStream(t, newStreamBody(tag, userID))
	put := stream.waitForEvent(t, sseEventPut, initialPutTimeout, nil)
	assertNoHoldoutReason(t, put.Evaluations.Evaluations)

	enableFeature(t, featureID, client)
	patch := stream.waitForEvent(t, sseEventPatch, patchTimeout, func(evt *gatewayproto.StreamEvaluationsEvent) bool {
		e, err := findFeature(evt.Evaluations.Evaluations, featureID)
		return err == nil && e.VariationValue == "A"
	})
	assert.False(t, patch.Evaluations.ForceUpdate)
	assertNoHoldoutReason(t, patch.Evaluations.Evaluations)

	// The server heartbeat interval defaults to 25 s; two minutes covers
	// several heartbeats and one full L1 cache TTL (1 minute) on the gateway,
	// which is where a wrong "unavailable" snapshot would be re-read.
	deadline := time.Now().Add(2 * time.Minute)
	heartbeats := 0
	for time.Now().Before(deadline) {
		stream.waitForEvent(t, sseEventHeartbeat, heartbeatTimeout, nil)
		heartbeats++
	}
	assert.GreaterOrEqual(t, heartbeats, 3, "expected several heartbeats in two minutes")
}

func assertNoHoldoutReason(t *testing.T, evaluations []*featureproto.Evaluation) {
	t.Helper()
	for _, e := range evaluations {
		assert.NotEqual(t, featureproto.Reason_HOLDOUT, e.GetReason().GetType(),
			"evaluation of %s carries Reason.HOLDOUT in an environment without holdouts", e.FeatureId)
		assert.Empty(t, e.GetReason().GetHoldoutGroupId())
	}
}
