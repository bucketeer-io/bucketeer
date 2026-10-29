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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

// ---- features ----

// newCreateFeatureReq builds a three-variation flag (A on, B off, C spare)
// in the holdout environment. The spare variation lets tests pick a holdout
// variation that is neither the on nor the off variation.
func newCreateFeatureReq(featureID string, tags ...string) *featureproto.CreateFeatureRequest {
	return &featureproto.CreateFeatureRequest{
		Id:          featureID,
		Name:        featureID,
		Description: "e2e-test-holdout-feature",
		Variations: []*featureproto.Variation{
			{Value: "A", Name: "Variation A", Description: "Thing does A"},
			{Value: "B", Name: "Variation B", Description: "Thing does B"},
			{Value: "C", Name: "Variation C", Description: "Thing does C"},
		},
		Tags:                     tags,
		DefaultOnVariationIndex:  wrapperspb.Int32(0),
		DefaultOffVariationIndex: wrapperspb.Int32(1),
		EnvironmentId:            holdoutEnv.id,
	}
}

func createFeature(t *testing.T, client featureclient.Client, req *featureproto.CreateFeatureRequest) *featureproto.Feature {
	t.Helper()
	var feature *featureproto.Feature
	err := withDeadlockRetry(t, "CreateFeature "+req.Id, func(ctx context.Context) error {
		resp, err := client.CreateFeature(ctx, req)
		if err != nil {
			return err
		}
		feature = resp.Feature
		return nil
	})
	require.NoError(t, err)
	return feature
}

// createEnabledFeature creates and enables a tagged flag and returns it. A
// flag needs at least one tag, so a private one is generated when none is
// given.
func createEnabledFeature(t *testing.T, client featureclient.Client, tags ...string) *featureproto.Feature {
	t.Helper()
	id := newFeatureID(t)
	if len(tags) == 0 {
		tags = []string{newTag(t)}
	}
	createFeature(t, client, newCreateFeatureReq(id, tags...))
	updateFeature(t, client, &featureproto.UpdateFeatureRequest{
		Id: id, EnvironmentId: holdoutEnv.id, Enabled: wrapperspb.Bool(true),
	})
	return getFeature(t, client, id)
}

func updateFeature(t *testing.T, client featureclient.Client, req *featureproto.UpdateFeatureRequest) {
	t.Helper()
	require.NoError(t, tryUpdateFeature(t, client, req))
}

func tryUpdateFeature(t *testing.T, client featureclient.Client, req *featureproto.UpdateFeatureRequest) error {
	t.Helper()
	req.EnvironmentId = holdoutEnv.id
	return withDeadlockRetry(t, "UpdateFeature "+req.Id, func(ctx context.Context) error {
		_, err := client.UpdateFeature(ctx, req)
		return err
	})
}

func getFeature(t *testing.T, client featureclient.Client, id string) *featureproto.Feature {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.GetFeature(ctx, &featureproto.GetFeatureRequest{Id: id, EnvironmentId: holdoutEnv.id})
	require.NoError(t, err)
	return resp.Feature
}

func variationID(f *featureproto.Feature, value string) string {
	for _, v := range f.Variations {
		if v.Value == value {
			return v.Id
		}
	}
	return ""
}

func addFixedRule(t *testing.T, client featureclient.Client, featureID, variationID, attribute, value string) {
	t.Helper()
	updateFeature(t, client, &featureproto.UpdateFeatureRequest{
		Id: featureID,
		RuleChanges: []*featureproto.RuleChange{{
			ChangeType: featureproto.ChangeType_CREATE,
			Rule: &featureproto.Rule{
				Id: newUUID(t),
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: variationID},
				},
				Clauses: []*featureproto.Clause{{
					Id:        newUUID(t),
					Attribute: attribute,
					Operator:  featureproto.Clause_EQUALS,
					Values:    []string{value},
				}},
			},
		}},
	})
}

// addSegmentRule makes the flag reference the segment, which is what marks a
// segment "in use" for GetSegmentUsers (see
// TestGrpcGetSegmentUsersRuleBasedSegment in test/e2e/gateway).
func addSegmentRule(t *testing.T, client featureclient.Client, featureID, variationID, segmentID string) {
	t.Helper()
	updateFeature(t, client, &featureproto.UpdateFeatureRequest{
		Id: featureID,
		RuleChanges: []*featureproto.RuleChange{{
			ChangeType: featureproto.ChangeType_CREATE,
			Rule: &featureproto.Rule{
				Id: newUUID(t),
				Strategy: &featureproto.Strategy{
					Type:          featureproto.Strategy_FIXED,
					FixedStrategy: &featureproto.FixedStrategy{Variation: variationID},
				},
				Clauses: []*featureproto.Clause{{
					Id:       newUUID(t),
					Operator: featureproto.Clause_SEGMENT,
					Values:   []string{segmentID},
				}},
			},
		}},
	})
}

func addPrerequisite(t *testing.T, client featureclient.Client, featureID, prerequisiteID, variationID string) {
	t.Helper()
	updateFeature(t, client, &featureproto.UpdateFeatureRequest{
		Id: featureID,
		PrerequisiteChanges: []*featureproto.PrerequisiteChange{{
			ChangeType:   featureproto.ChangeType_CREATE,
			Prerequisite: &featureproto.Prerequisite{FeatureId: prerequisiteID, VariationId: variationID},
		}},
	})
}

// refreshFeatureCache runs the FeatureFlagCacher batch job so the gateway's
// Redis snapshot contains the flags just written (the gateway's L1 is already
// evicted by the FEATURE domain event), then waits until the gateway serves
// every given flag for the tag.
func refreshFeatureCache(t *testing.T, tag string, featureIDs ...string) {
	t.Helper()
	runBatchJob(t, btproto.BatchJob_FeatureFlagCacher)
	require.Eventually(t, func() bool {
		resp := grpcGetEvaluations(t, tag, "cache-probe-"+newUUID(t))
		for _, id := range featureIDs {
			if findEvaluation(resp.Evaluations.Evaluations, id) == nil {
				return false
			}
		}
		return true
	}, cachePropagationTimeout, pollInterval, "flags %v not visible on the gateway", featureIDs)
}

// ---- segments ----

func createSegment(t *testing.T, client featureclient.Client) *featureproto.Segment {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.CreateSegment(ctx, &featureproto.CreateSegmentRequest{
		EnvironmentId: holdoutEnv.id,
		Name:          prefixed("segment", newUUID(t)),
		Description:   "e2e-test-holdout-segment",
	})
	require.NoError(t, err)
	return resp.Segment
}

// uploadSegmentUsersAndWait follows test/e2e/feature/segment_user_test.go:
// the upload is asynchronous, so poll until every user is listed.
func uploadSegmentUsersAndWait(t *testing.T, client featureclient.Client, segmentID string, userIDs []string) {
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
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		res, err := client.ListSegmentUsers(ctx, &featureproto.ListSegmentUsersRequest{
			SegmentId:     segmentID,
			State:         wrapperspb.Int32(int32(featureproto.SegmentUser_INCLUDED)),
			EnvironmentId: holdoutEnv.id,
		})
		return err == nil && len(res.Users) >= len(userIDs)
	}, segmentUploadTimeout, 2*time.Second, "segment users not uploaded")
	// The segment must leave UPLOADING before a SEGMENT group can start.
	require.Eventually(t, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		res, err := client.GetSegment(ctx, &featureproto.GetSegmentRequest{Id: segmentID, EnvironmentId: holdoutEnv.id})
		return err == nil && res.Segment.Status == featureproto.Segment_SUCEEDED
	}, segmentUploadTimeout, 2*time.Second, "segment did not reach SUCEEDED")
}

// ---- holdout groups ----

type groupSpec struct {
	assignment      featureproto.HoldoutGroup_AssignmentType
	percentage      float64
	segmentID       string
	startAt, stopAt int64
	startAnalysisAt int64
}

func (s groupSpec) request(t *testing.T) *featureproto.CreateHoldoutGroupRequest {
	t.Helper()
	return &featureproto.CreateHoldoutGroupRequest{
		EnvironmentId:   holdoutEnv.id,
		Name:            newGroupName(t),
		Description:     "e2e-test-holdout-group",
		AssignmentType:  s.assignment,
		Percentage:      s.percentage,
		SegmentId:       s.segmentID,
		StartAt:         s.startAt,
		StopAt:          s.stopAt,
		StartAnalysisAt: s.startAnalysisAt,
	}
}

// waitingHashSpec is a 50% HASH group on a far-future window.
func waitingHashSpec() groupSpec {
	start, stop := farFutureWindow()
	return groupSpec{assignment: featureproto.HoldoutGroup_HASH, percentage: 50, startAt: start, stopAt: stop}
}

// servingHashSpec is a 50% HASH group meant to be started manually (see
// servingWindowNow). Callers must hold servingMu.
func servingHashSpec() groupSpec {
	start, stop := servingWindowNow()
	return groupSpec{assignment: featureproto.HoldoutGroup_HASH, percentage: 50, startAt: start, stopAt: stop}
}

func createHoldoutGroup(t *testing.T, client featureclient.Client, req *featureproto.CreateHoldoutGroupRequest) *featureproto.HoldoutGroup {
	t.Helper()
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "CreateHoldoutGroup", func(ctx context.Context) error {
		resp, err := client.CreateHoldoutGroup(ctx, req)
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	registerGroupCleanup(t, client, group.Id)
	return group
}

// registerGroupCleanup leaves no group behind, whatever state the test
// failed in: WAITING groups are deleted, serving groups are stopped and
// archived, stopped groups are archived. Errors are logged, not asserted,
// because the test may already be failing.
func registerGroupCleanup(t *testing.T, client featureclient.Client, id string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		resp, err := client.GetHoldoutGroup(ctx, &featureproto.GetHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				t.Logf("cleanup: get holdout group %s: %v", id, err)
			}
			return
		}
		g := resp.HoldoutGroup
		if g.Archived || g.Deleted {
			return
		}
		switch g.Status {
		case featureproto.HoldoutGroup_WAITING:
			if _, err := client.DeleteHoldoutGroup(ctx, &featureproto.DeleteHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id}); err != nil {
				t.Logf("cleanup: delete holdout group %s: %v", id, err)
			}
			return
		case featureproto.HoldoutGroup_RUNNING, featureproto.HoldoutGroup_ANALYSIS:
			if _, err := client.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id}); err != nil {
				t.Logf("cleanup: stop holdout group %s: %v", id, err)
				return
			}
		}
		if _, err := client.ArchiveHoldoutGroup(ctx, &featureproto.ArchiveHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id}); err != nil {
			t.Logf("cleanup: archive holdout group %s: %v", id, err)
		}
	})
}

func getHoldoutGroup(t *testing.T, client featureclient.Client, id string) *featureproto.HoldoutGroup {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.GetHoldoutGroup(ctx, &featureproto.GetHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id})
	require.NoError(t, err)
	return resp.HoldoutGroup
}

func addHoldoutFeatures(t *testing.T, client featureclient.Client, id string, inputs ...*featureproto.HoldoutGroupFeatureInput) *featureproto.HoldoutGroup {
	t.Helper()
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "AddHoldoutGroupFeatures", func(ctx context.Context) error {
		resp, err := client.AddHoldoutGroupFeatures(ctx, &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: id, EnvironmentId: holdoutEnv.id, Features: inputs,
		})
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	return group
}

func tryAddHoldoutFeatures(t *testing.T, client featureclient.Client, id string, inputs ...*featureproto.HoldoutGroupFeatureInput) error {
	t.Helper()
	return withDeadlockRetry(t, "AddHoldoutGroupFeatures", func(ctx context.Context) error {
		_, err := client.AddHoldoutGroupFeatures(ctx, &featureproto.AddHoldoutGroupFeaturesRequest{
			Id: id, EnvironmentId: holdoutEnv.id, Features: inputs,
		})
		return err
	})
}

func removeHoldoutFeatures(t *testing.T, client featureclient.Client, id string, featureIDs ...string) *featureproto.HoldoutGroup {
	t.Helper()
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "RemoveHoldoutGroupFeatures", func(ctx context.Context) error {
		resp, err := client.RemoveHoldoutGroupFeatures(ctx, &featureproto.RemoveHoldoutGroupFeaturesRequest{
			Id: id, EnvironmentId: holdoutEnv.id, FeatureIds: featureIDs,
		})
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	return group
}

func target(f *featureproto.Feature, holdoutValue string) *featureproto.HoldoutGroupFeatureInput {
	return &featureproto.HoldoutGroupFeatureInput{FeatureId: f.Id, HoldoutVariationId: variationID(f, holdoutValue)}
}

func startHoldoutGroup(t *testing.T, client featureclient.Client, id string) *featureproto.HoldoutGroup {
	t.Helper()
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "StartHoldoutGroup", func(ctx context.Context) error {
		resp, err := client.StartHoldoutGroup(ctx, &featureproto.StartHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id})
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	return group
}

func startHoldoutGroupAnalysis(t *testing.T, client featureclient.Client, id string) *featureproto.HoldoutGroup {
	t.Helper()
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "StartHoldoutGroupAnalysis", func(ctx context.Context) error {
		resp, err := client.StartHoldoutGroupAnalysis(ctx, &featureproto.StartHoldoutGroupAnalysisRequest{Id: id, EnvironmentId: holdoutEnv.id})
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	return group
}

func stopHoldoutGroup(t *testing.T, client featureclient.Client, id string) *featureproto.HoldoutGroup {
	t.Helper()
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "StopHoldoutGroup", func(ctx context.Context) error {
		resp, err := client.StopHoldoutGroup(ctx, &featureproto.StopHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id})
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	return group
}

func archiveHoldoutGroup(t *testing.T, client featureclient.Client, id string) *featureproto.HoldoutGroup {
	t.Helper()
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "ArchiveHoldoutGroup", func(ctx context.Context) error {
		resp, err := client.ArchiveHoldoutGroup(ctx, &featureproto.ArchiveHoldoutGroupRequest{Id: id, EnvironmentId: holdoutEnv.id})
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	return group
}

func updateHoldoutGroup(t *testing.T, client featureclient.Client, req *featureproto.UpdateHoldoutGroupRequest) *featureproto.HoldoutGroup {
	t.Helper()
	req.EnvironmentId = holdoutEnv.id
	var group *featureproto.HoldoutGroup
	err := withDeadlockRetry(t, "UpdateHoldoutGroup", func(ctx context.Context) error {
		resp, err := client.UpdateHoldoutGroup(ctx, req)
		if err != nil {
			return err
		}
		group = resp.HoldoutGroup
		return nil
	})
	require.NoError(t, err)
	return group
}

// bumpSecond waits until the wall clock has moved past the group's updated_at
// so the next mutation observably bumps it (updated_at has second precision).
func bumpSecond(t *testing.T, g *featureproto.HoldoutGroup) {
	t.Helper()
	waitUntilUnix(t, g.UpdatedAt)
}

// pickUsers returns one user the group holds out and one it does not, using
// the shared evaluation library against the server-assigned seed and range
// (RFC §4.2.1), instead of guessing.
func pickUsers(t *testing.T, g *featureproto.HoldoutGroup) (heldOut, notHeldOut string) {
	t.Helper()
	require.Equal(t, featureproto.HoldoutGroup_HASH, g.AssignmentType)
	require.NotEmpty(t, g.SamplingSeed)
	for i := 0; i < 200 && (heldOut == "" || notHeldOut == ""); i++ {
		id := newUserID(t)
		if evaluation.IsHeldOut(g, &userproto.User{Id: id}, nil) {
			if heldOut == "" {
				heldOut = id
			}
		} else if notHeldOut == "" {
			notHeldOut = id
		}
	}
	require.NotEmpty(t, heldOut, "no held-out user found")
	require.NotEmpty(t, notHeldOut, "no non-held-out user found")
	return heldOut, notHeldOut
}

// ---- gateway ----

// grpcGetEvaluationsWithRetry mirrors test/e2e/gateway: codes.Internal is
// transient while the feature cache snapshot is inconsistent.
func grpcGetEvaluationsWithRetry(t *testing.T, apiKeyPath string, req *gatewayproto.GetEvaluationsRequest) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	c := newGatewayClientT(t, apiKeyPath)
	var lastErr error
	for i := 0; i < getEvaluationsRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		resp, err := c.GetEvaluations(ctx, req)
		cancel()
		if err == nil {
			return resp
		}
		lastErr = err
		if status.Code(err) != codes.Internal {
			break
		}
		if i < getEvaluationsRetryAttempts-1 {
			t.Logf("GetEvaluations returned Internal (attempt %d/%d). Retrying in %v: %v",
				i+1, getEvaluationsRetryAttempts, getEvaluationsRetryInterval, err)
			time.Sleep(getEvaluationsRetryInterval)
		}
	}
	require.NoError(t, lastErr)
	return nil
}

// grpcGetEvaluations is a first request from a current client SDK: it sends
// an empty UserEvaluationCondition (the path every maintained SDK uses; the
// legacy path without it requires a tag) and gets a full, forced update.
func grpcGetEvaluations(t *testing.T, tag, userID string) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	return grpcGetEvaluationsWithUser(t, tag, &userproto.User{Id: userID})
}

func grpcGetEvaluationsWithUser(t *testing.T, tag string, user *userproto.User) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	return grpcGetEvaluationsWithRetry(t, holdoutEnv.clientKeyPath, &gatewayproto.GetEvaluationsRequest{
		Tag:                     tag,
		User:                    user,
		UserEvaluationCondition: &gatewayproto.GetEvaluationsRequest_UserEvaluationCondition{},
	})
}

// grpcGetEvaluationsByEvaluatedAt polls as a client SDK does: with the
// UserEvaluationsID and evaluatedAt from the previous response.
func grpcGetEvaluationsByEvaluatedAt(t *testing.T, tag, userID, ueid string, evaluatedAt int64) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	return grpcGetEvaluationsWithRetry(t, holdoutEnv.clientKeyPath, &gatewayproto.GetEvaluationsRequest{
		Tag:               tag,
		User:              &userproto.User{Id: userID},
		UserEvaluationsId: ueid,
		UserEvaluationCondition: &gatewayproto.GetEvaluationsRequest_UserEvaluationCondition{
			EvaluatedAt: evaluatedAt,
		},
	})
}

func grpcGetEvaluation(t *testing.T, tag, featureID, userID string) *gatewayproto.GetEvaluationResponse {
	t.Helper()
	c := newGatewayClientT(t, holdoutEnv.clientKeyPath)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.GetEvaluation(ctx, &gatewayproto.GetEvaluationRequest{
		Tag:       tag,
		User:      &userproto.User{Id: userID},
		FeatureId: featureID,
	})
	require.NoError(t, err)
	return resp
}

func grpcGetFeatureFlags(t *testing.T, tag, ffID string, requestedAt int64) *gatewayproto.GetFeatureFlagsResponse {
	t.Helper()
	c := newGatewayClientT(t, holdoutEnv.serverKeyPath)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.GetFeatureFlags(ctx, &gatewayproto.GetFeatureFlagsRequest{
		Tag:            tag,
		FeatureFlagsId: ffID,
		RequestedAt:    requestedAt,
		SourceId:       eventproto.SourceId_GO_SERVER,
		SdkVersion:     "v0.0.1-e2e-test",
	})
	require.NoError(t, err)
	return resp
}

func grpcGetSegmentUsers(t *testing.T, segmentIDs []string, requestedAt int64) *gatewayproto.GetSegmentUsersResponse {
	t.Helper()
	c := newGatewayClientT(t, holdoutEnv.serverKeyPath)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := c.GetSegmentUsers(ctx, &gatewayproto.GetSegmentUsersRequest{
		SegmentIds:  segmentIDs,
		RequestedAt: requestedAt,
		SourceId:    eventproto.SourceId_GO_SERVER,
		SdkVersion:  "v0.0.1-e2e-test",
	})
	require.NoError(t, err)
	return resp
}

func findEvaluation(evaluations []*featureproto.Evaluation, featureID string) *featureproto.Evaluation {
	for _, e := range evaluations {
		if e.FeatureId == featureID {
			return e
		}
	}
	return nil
}

func hasHoldoutReason(evaluations []*featureproto.Evaluation) bool {
	for _, e := range evaluations {
		if e.GetReason().GetType() == featureproto.Reason_HOLDOUT {
			return true
		}
	}
	return false
}

// requireHoldoutEvaluation asserts that featureID was served from the group:
// the holdout variation, Reason.Type HOLDOUT and the group id on the reason.
func requireHoldoutEvaluation(t *testing.T, evaluations []*featureproto.Evaluation, featureID, holdoutVariationID, groupID string) {
	t.Helper()
	e := findEvaluation(evaluations, featureID)
	require.NotNil(t, e, "evaluation for %s missing", featureID)
	require.Equal(t, holdoutVariationID, e.VariationId, "variation of %s", featureID)
	require.Equal(t, featureproto.Reason_HOLDOUT, e.GetReason().GetType(), "reason of %s", featureID)
	require.Equal(t, groupID, e.GetReason().GetHoldoutGroupId(), "holdout_group_id of %s", featureID)
}

// requireNormalEvaluation asserts that featureID was served by the normal
// path with the given variation and reason.
func requireNormalEvaluation(t *testing.T, evaluations []*featureproto.Evaluation, featureID, variationID string, reason featureproto.Reason_Type) {
	t.Helper()
	e := findEvaluation(evaluations, featureID)
	require.NotNil(t, e, "evaluation for %s missing", featureID)
	require.Equal(t, variationID, e.VariationId, "variation of %s", featureID)
	require.Equal(t, reason, e.GetReason().GetType(), "reason of %s", featureID)
	require.Empty(t, e.GetReason().GetHoldoutGroupId(), "holdout_group_id of %s", featureID)
}

// waitForHoldoutServing polls GetEvaluations until the held-out user sees
// the holdout variation for featureID, bounded by the cache propagation
// time, and returns the response that observed it.
func waitForHoldoutServing(t *testing.T, tag, heldOutUser, featureID, holdoutVariationID string) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	var resp *gatewayproto.GetEvaluationsResponse
	require.Eventually(t, func() bool {
		resp = grpcGetEvaluations(t, tag, heldOutUser)
		e := findEvaluation(resp.Evaluations.Evaluations, featureID)
		return e != nil && e.VariationId == holdoutVariationID && e.GetReason().GetType() == featureproto.Reason_HOLDOUT
	}, cachePropagationTimeout, pollInterval, "holdout for %s never reached the gateway", featureID)
	return resp
}

// waitForHoldoutCleared polls GetEvaluations until no evaluation of the user
// carries Reason.HOLDOUT any more.
func waitForHoldoutCleared(t *testing.T, tag, userID string) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	var resp *gatewayproto.GetEvaluationsResponse
	require.Eventually(t, func() bool {
		resp = grpcGetEvaluations(t, tag, userID)
		return !hasHoldoutReason(resp.Evaluations.Evaluations)
	}, cachePropagationTimeout, pollInterval, "holdout for %s never cleared on the gateway", userID)
	return resp
}

// restEvaluationsBody builds the /get_evaluations request body for the REST
// helpers in test/e2e/util. The condition is always sent (evaluated_at as a
// string for protojson) so the gateway takes the same evaluatedAt path as
// current SDKs.
func restEvaluationsBody(tag, userID, ueid string, evaluatedAt int64) map[string]interface{} {
	body := map[string]interface{}{
		"tag":  tag,
		"user": map[string]interface{}{"id": userID},
		"user_evaluation_condition": map[string]interface{}{
			"evaluated_at":            fmt.Sprintf("%d", evaluatedAt),
			"user_attributes_updated": false,
		},
	}
	if ueid != "" {
		body["user_evaluations_id"] = ueid
	}
	return body
}
