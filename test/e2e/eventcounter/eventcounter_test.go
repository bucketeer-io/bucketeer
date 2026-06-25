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

package eventcounter

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	gwapi "github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	gatewayclient "github.com/bucketeer-io/bucketeer/v2/pkg/api/client"
	btclient "github.com/bucketeer-io/bucketeer/v2/pkg/batch/client"
	ecclient "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	btproto "github.com/bucketeer-io/bucketeer/v2/proto/batch"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	ecproto "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
	"github.com/bucketeer-io/bucketeer/v2/test/e2e/util"
)

const (
	prefixTestName        = "e2e-test"
	timeout               = 3 * time.Minute  // Overall test timeout
	grpcTimeout           = 10 * time.Second // Individual gRPC call timeout
	retryTimes            = 50
	deadlockRetryAttempts = 3
)

const defaultVariationID = "default"

// experimentResultUsersPerVariation is the per-arm sample size for
// TestExperimentResult / TestGrpcExperimentResult. At n=50 the empirical-Bayes
// prior weight is ~2%, so value-metric goldens reflect posterior signal from
// the 10 vs 15 per-user means rather than prior arithmetic.
const (
	experimentResultUsersPerVariation = 50
	experimentResultValueVariationA   = 10.0
	experimentResultValueVariationB   = 15.0
)

var (
	webGatewayAddr                 = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort                 = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert                 = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath                     = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath               = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr                    = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort                    = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert                    = flag.String("gateway-cert", "", "Gateway crt file")
	sysAdminAccessTokenPath        = flag.String("sys-admin-access-token", "", "System admin access token path")
	orgOwnerDefaultAccessTokenPath = flag.String("org-owner-default-access-token", "", "Organization admin access token path")
	orgOwnerE2EAccessTokenPath     = flag.String("org-owner-e2e-access-token", "", "Organization admin (e2e org) access token path")
	envEditorAccessTokenPath       = flag.String("env-editor-access-token", "", "Environment editor access token path")
	envViewerAccessTokenPath       = flag.String("env-viewer-access-token", "", "Environment viewer access token path")
	environmentID                  = flag.String("environment-id", "", "Environment id")
	organizationID                 = flag.String("organization-id", "", "Organization ID")
	testID                         = flag.String("test-id", "", "test ID")
	compareFloatOpt                = cmpopts.EquateApprox(0, 0.0001)
	compareFloatBayesian           = cmpopts.EquateApprox(0, 0.03) // 3% tolerance for Bayesian results to avoid flaky tests
)

func TestGrpcExperimentGoalCount(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	createFeature(t, featureClient, featureID, tag, variationVarA, variationVarB)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}
	// Because we set the user to the individual targeting when creating the flag,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)

	// Get the latest version after all changes in the feature flag
	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := startAt.Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcExperimentGoalCount", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// Evaluation events must always be sent before goal events
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, f.Variations[0].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.2), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Add(-2*time.Hour).Unix())

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA := getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, 0.5, compareFloatOpt); diff != "" {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSumPerUserMean, 0.5, compareFloatOpt); diff != "" {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSumPerUserVariance, float64(0), compareFloatOpt); diff != "" {
			continue
		}
		vcB := getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if vcB.ValueSum != float64(0) {
			continue
		}
		if vcB.ValueSumPerUserMean != float64(0.0) {
			continue
		}
		if vcB.ValueSumPerUserVariance != float64(0.0) {
			continue
		}
		break
	}
}

func TestExperimentGoalCount(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	createFeature(t, featureClient, featureID, tag, variationVarA, variationVarB)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Because we set the user to the individual targeting when creating the flag,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)

	// Get the latest version after all changes in the feature flag
	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := startAt.Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentGoalCount", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// Evaluation events must always be sent before goal events
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, f.Variations[0].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.2), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Add(-2*time.Hour).Unix())

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA := getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, 0.5, compareFloatOpt); diff != "" {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSumPerUserMean, 0.5, compareFloatOpt); diff != "" {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSumPerUserVariance, float64(0), compareFloatOpt); diff != "" {
			continue
		}

		vcB := getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSumPerUserMean, float64(0), compareFloatOpt); diff != "" {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSumPerUserVariance, float64(0), compareFloatOpt); diff != "" {
			continue
		}
		break
	}
}

func TestGrpcExperimentResult(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDsA, userIDsB := makeExperimentResultUserIDs(t, uuid)
	featureID := createFeatureID(t, uuid)

	createFeature(t, featureClient, featureID, tag, "a", "b")
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargetingBulkSplit(
		t, featureID, f.Variations[0].Id, f.Variations[1].Id, userIDsA, userIDsB, featureClient)

	updateFeatueFlagCache(t)

	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := time.Now().Add(2 * time.Hour)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	time.Sleep(15 * time.Second)

	grpcRegisterEvaluationEventsBatch(
		t, featureID, experiment.FeatureVersion, userIDsA, f.Variations[0].Id, tag, reason)
	grpcRegisterEvaluationEventsBatch(
		t, featureID, experiment.FeatureVersion, userIDsB, f.Variations[1].Id, tag, reason)

	time.Sleep(10 * time.Second)

	goalTimestamp := time.Now().Unix()
	grpcRegisterGoalEventsBatch(t, goalIDs[0], tag, goalValuesForUsers(userIDsA, experimentResultValueVariationA), goalTimestamp)
	grpcRegisterGoalEventsBatch(t, goalIDs[0], tag, goalValuesForUsers(userIDsB, experimentResultValueVariationB), goalTimestamp)

	waitAndCheckExperimentResult(t, ecClient, experiment, goalIDs[0])
	res, err := getExperiment(t, experimentClient, experiment.Id)
	if err != nil {
		t.Fatalf("Failed to get experiment. ID: %s. Error: %v", experiment.Id, err)
	}
	if res.Experiment.Status != experimentproto.Experiment_RUNNING {
		expected, _ := experimentproto.Experiment_Status_name[int32(experimentproto.Experiment_RUNNING)]
		actual, _ := experimentproto.Experiment_Status_name[int32(res.Experiment.Status)]
		t.Fatalf("the status of experiment is not correct. expected: %s, but got %s", expected, actual)
	}
}

func TestExperimentResult(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDsA, userIDsB := makeExperimentResultUserIDs(t, uuid)
	featureID := createFeatureID(t, uuid)

	createFeature(t, featureClient, featureID, tag, "a", "b")
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargetingBulkSplit(
		t, featureID, f.Variations[0].Id, f.Variations[1].Id, userIDsA, userIDsB, featureClient)

	updateFeatueFlagCache(t)

	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := time.Now().Add(2 * time.Hour)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	time.Sleep(15 * time.Second)

	registerEvaluationEventsBatch(
		t, featureID, experiment.FeatureVersion, userIDsA, f.Variations[0].Id, tag, reason)
	registerEvaluationEventsBatch(
		t, featureID, experiment.FeatureVersion, userIDsB, f.Variations[1].Id, tag, reason)

	time.Sleep(10 * time.Second)

	goalTimestamp := time.Now().Unix()
	registerGoalEventsBatch(t, goalIDs[0], tag, goalValuesForUsers(userIDsA, experimentResultValueVariationA), goalTimestamp)
	registerGoalEventsBatch(t, goalIDs[0], tag, goalValuesForUsers(userIDsB, experimentResultValueVariationB), goalTimestamp)

	waitAndCheckExperimentResult(t, ecClient, experiment, goalIDs[0])
	res, err := getExperiment(t, experimentClient, experiment.Id)
	if err != nil {
		t.Fatalf("Failed to get experiment. ID: %s. Error: %v", experiment.Id, err)
	}
	if res.Experiment.Status != experimentproto.Experiment_RUNNING {
		expected, _ := experimentproto.Experiment_Status_name[int32(experimentproto.Experiment_RUNNING)]
		actual, _ := experimentproto.Experiment_Status_name[int32(res.Experiment.Status)]
		t.Fatalf("the status of experiment is not correct. expected: %s, but got %s", expected, actual)
	}
}

// TestExperimentGoalCountWinsorization verifies the per-user p99 winsorization
// applied in the goal-count DWH query (goal_event.sql / goal_count.sql). The
// cap only bites once there are enough users for the top 1% to exclude the
// whale (NTILE(100) needs >= 100 rows on MySQL), so this fixture uses 120
// users: 119 normal goal values plus one "whale" several orders of magnitude
// larger. The whale's per-user value sits in the top percentile and must be
// capped down to the bulk level before aggregation.
//
//	capped   ValueSum ≈ 119*10 + cap(=10) = 1,200
//	uncapped ValueSum ≈ 119*10 + 100,000  = 101,190
//
// We wait until all 120 users are linked (so the whale is definitely counted)
// and then assert ValueSum is far below the uncapped figure. The threshold is
// deliberately loose (< 50,000) so it holds across all three DWH backends,
// whose percentile implementations differ slightly — only an uncapped result
// could exceed it.
func TestExperimentGoalCountWinsorization(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	const (
		totalUsers  = 120
		normalValue = float64(10)
		whaleValue  = float64(100000)
	)

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := make([]string, 0, totalUsers)
	for i := 0; i < totalUsers; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	createFeature(t, featureClient, featureID, tag, "a", "b")
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Route every user to variation A via individual targeting (one update for
	// the whole list), so all goal values land in a single pooled distribution
	// with a known per-user breakdown and a TARGET evaluation reason.
	reason := &featureproto.Reason{Type: featureproto.Reason_TARGET}
	addFeatureIndividualTargetingBulk(t, featureID, f.Variations[0].Id, userIDs, featureClient)

	updateFeatueFlagCache(t)

	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := time.Now().Add(2 * time.Hour)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentGoalCountWinsorization", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	variationIDs := make([]string, 0, len(experiment.Variations))
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
	}

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// Evaluation events must be sent (and land) before goal events so the
	// subscriber can link each goal event to the user's evaluation.
	registerEvaluationEventsBatch(t, featureID, experiment.FeatureVersion, userIDs, f.Variations[0].Id, tag, reason)

	// Wait a few seconds so the evaluations become available for linking.
	time.Sleep(10 * time.Second)

	// One goal event per user: the last user is the whale.
	goalValues := make(map[string]float64, totalUsers)
	for i, userID := range userIDs {
		if i == totalUsers-1 {
			goalValues[userID] = whaleValue
		} else {
			goalValues[userID] = normalValue
		}
	}
	registerGoalEventsBatch(t, goalIDs[0], tag, goalValues, time.Now().Unix())

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout after %d attempts", retryTimes)
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, experiment.FeatureVersion, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			}
			continue
		}
		if resp == nil || len(resp.VariationCounts) == 0 {
			continue
		}
		vcA := getVariationCount(resp.VariationCounts, f.Variations[0].Id)
		if vcA == nil {
			continue
		}
		// Wait until every user (including the whale) has been linked, so the
		// whale's value is part of the aggregate we assert on.
		if vcA.UserCount != totalUsers {
			t.Logf("Retry %d/%d: linked users %d/%d", i+1, retryTimes, vcA.UserCount, totalUsers)
			continue
		}
		// With the whale present, an uncapped ValueSum would be ~101,190; the
		// winsorized ValueSum is ~1,200. Anything below 50,000 proves the whale
		// was capped.
		if vcA.ValueSum >= 50000 {
			t.Fatalf(
				"winsorization did not cap the whale: ValueSum=%f (expected ~1200 capped, ~101190 uncapped)",
				vcA.ValueSum,
			)
		}
		t.Logf("winsorization OK: capped ValueSum=%f for %d users", vcA.ValueSum, vcA.UserCount)
		break
	}
}

func TestGrpcMultiGoalsEventCounter(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	createFeature(t, featureClient, featureID, tag, variationVarA, variationVarB)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userIDs[0], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[1], f.Variations[1].Id, featureClient)

	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)

	// Get the latest version after all changes in the feature flag
	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 3)
	startAt := time.Now().Add(-time.Hour)
	stopAt := startAt.Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcMultiGoalsEventCounter", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// Evaluation events must always be sent before goal events
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[0], f.Variations[0].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[1], f.Variations[1].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Add(-2*time.Hour).Unix())

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		// Goal 0.
		resp, err := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA := getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, 0.6, compareFloatOpt); diff != "" {
			continue
		}

		vcB := getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}

		// Goal 1.
		resp, err = getExperimentGoalCount(t, ecClient, goalIDs[1], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[1] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA = getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}

		vcB = getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 1 {
			continue
		}
		if vcB.EventCount != 2 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, 0.4, compareFloatOpt); diff != "" {
			continue
		}

		// Goal 2.
		resp, err = getExperimentGoalCount(t, ecClient, goalIDs[2], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[2] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA = getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}

		vcB = getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}
		break
	}
}

func TestMultiGoalsEventCounter(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	createFeature(t, featureClient, featureID, tag, variationVarA, variationVarB)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userIDs[0], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[1], f.Variations[1].Id, featureClient)

	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)

	// Get the latest version after all changes in the feature flag
	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 3)
	startAt := time.Now().Add(-time.Hour)
	stopAt := startAt.Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestMultiGoalsEventCounter", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// Evaluation events must always be sent before goal events
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[0], f.Variations[0].Id, tag, reason)
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[1], f.Variations[1].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Add(-2*time.Hour).Unix())

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		// Goal 0.
		resp, err := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA := getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, 0.6, compareFloatOpt); diff != "" {
			continue
		}

		vcB := getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}

		// Goal 1.
		resp, err = getExperimentGoalCount(t, ecClient, goalIDs[1], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[1] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA = getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}

		vcB = getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 1 {
			continue
		}
		if vcB.EventCount != 2 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, 0.4, compareFloatOpt); diff != "" {
			continue
		}

		// Goal 2.
		resp, err = getExperimentGoalCount(t, ecClient, goalIDs[2], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[2] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		vcA = getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}

		vcB = getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}
		break
	}
}

func TestHTTPTrack(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)
	value := float64(1.23)

	variationVarA := "a"
	variationVarB := "b"
	createFeature(t, featureClient, featureID, tag, variationVarA, variationVarB)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)

	// Get the latest version after all changes in the feature flag
	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := startAt.Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestHTTPTrack", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// Evaluation events must always be sent before goal events
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, f.Variations[0].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	// Send track events
	// This track event will be converted to a goal event on the backend.
	sendHTTPTrack(t, userID, goalIDs[0], tag, value)

	// Check the count
	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment goal count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if len(resp.VariationCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalId)
		}

		// variation a
		vcA := getVariationCount(resp.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 1 {
			continue
		}
		if diff := cmp.Diff(vcA.ValueSum, value, compareFloatOpt); diff != "" {
			continue
		}
		// variation b
		vcB := getVariationCount(resp.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation a is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if diff := cmp.Diff(vcB.ValueSum, float64(0), compareFloatOpt); diff != "" {
			continue
		}
		break
	}
}

func TestGrpcExperimentEvaluationEventCount(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()
	uuid := newUUID(t)

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)
	variationVarA := "a"
	variationVarB := "b"

	createFeature(t, featureClient, featureID, tag, variationVarA, variationVarB)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)

	// Get the latest version after all changes in the feature flag
	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range f.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := startAt.Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx,
		t,
		experimentClient,
		"TestGrpcExperimentEvaluationEventCount",
		featureID,
		goalIDs,
		f.Variations[0].Id,
		startAt, stopAt,
	)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, variations[variationVarA].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, variations[variationVarA].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, variations[variationVarA].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentEvaluationCount(t, ecClient, featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment evaluation count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}

		if resp == nil {
			continue
		}
		ec := resp
		if ec.FeatureId != featureID {
			t.Fatalf("feature ID is not correct: %s", ec.FeatureId)
		}
		if ec.FeatureVersion != f.Version {
			t.Fatalf("feature version is not correct: %d", ec.FeatureVersion)
		}

		vcA := getVariationCount(ec.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 3 {
			continue
		}
		if vcA.ValueSum != float64(0) {
			continue
		}

		vcB := getVariationCount(ec.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if vcB.ValueSum != float64(0) {
			continue
		}
		break
	}
}

func TestExperimentEvaluationEventCount(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)
	experimentClient := newExperimentClient(t)
	defer experimentClient.Close()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)
	variationVarA := "a"
	variationVarB := "b"

	createFeature(t, featureClient, featureID, tag, variationVarA, variationVarB)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)

	// Get the latest version after all changes in the feature flag
	f, err = getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range f.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}
	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Add(-time.Hour)
	stopAt := startAt.Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx,
		t,
		experimentClient,
		"TestExperimentEvaluationEventCount",
		featureID,
		goalIDs,
		f.Variations[0].Id,
		startAt, stopAt,
	)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		stopExperiment(cleanupCtx, t, experimentClient, experiment.Id)
	})

	// Wait for the on-demand subscriber to create PubSub subscriptions.
	time.Sleep(15 * time.Second)

	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, variations[variationVarA].Id, tag, reason)
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, variations[variationVarA].Id, tag, reason)
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userID, variations[variationVarA].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentEvaluationCount(t, ecClient, featureID, f.Version, variationIDs)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment evaluation count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}
		if resp == nil {
			continue
		}
		ec := resp
		if ec.FeatureId != featureID {
			t.Fatalf("feature ID is not correct: %s", ec.FeatureId)
		}
		if ec.FeatureVersion != f.Version {
			t.Fatalf("feature version is not correct: %d", ec.FeatureVersion)
		}

		vcA := getVariationCount(ec.VariationCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 3 {
			continue
		}
		if vcA.ValueSum != float64(0) {
			continue
		}

		vcB := getVariationCount(ec.VariationCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 0 {
			continue
		}
		if vcB.EventCount != 0 {
			continue
		}
		if vcB.ValueSum != float64(0) {
			continue
		}
		break
	}
}

func TestGetEvaluationTimeseriesCount(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
	uuid := newUUID(t)
	featureID := createFeatureID(t, uuid)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	createFeature(t, featureClient, featureID, tag, "a", "b")
	userIDs := []string{}
	for i := 0; i < 8; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	// Because the event-persister-dwh calls the EvaluateFeatures API
	// and it uses the feature flag cache, we must update it before sending events.
	updateFeatueFlagCache(t)
	f, err := getFeature(t, featureClient, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}

	// Register variation
	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], f.Variations[0].Id, tag, nil)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[1], f.Variations[0].Id, tag, nil)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[2], f.Variations[0].Id, tag, nil)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], f.Variations[0].Id, tag, nil)
	// Register variation
	registerEvaluationEvent(t, featureID, f.Version, userIDs[3], f.Variations[1].Id, tag, nil)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[4], f.Variations[1].Id, tag, nil)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, f.Version, userIDs[3], f.Variations[1].Id, tag, nil)
	// Register variation
	registerEvaluationEvent(t, featureID, f.Version, userIDs[5], defaultVariationID, tag, nil)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[6], defaultVariationID, tag, nil)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[7], defaultVariationID, tag, nil)
	expectedUserVal := []float64{3, 2, 3}
	expectedEventVal := []float64{4, 3, 3}
	i := 0
LOOP:
	for {
		time.Sleep(10 * time.Second)
		i++
		if i == retryTimes {
			t.Fatalf("retry timeout")
		}
		res, err := getEvaluationTimeseriesCount(t, featureID, ecClient, ecproto.GetEvaluationTimeseriesCountRequest_THIRTY_DAYS)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the evaluation timeseries count. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				continue
			}
		}
		if len(res.UserCounts) != 3 {
			t.Fatalf("the number of user counts is not correct: %d", len(res.UserCounts))
		}
		expectedVIDs := []string{f.Variations[0].Id, f.Variations[1].Id, defaultVariationID}
		if len(res.EventCounts) != len(expectedVIDs) {
			t.Fatalf("the number of event counts is not correct: %d", len(res.UserCounts))
		}
		for idx, uc := range res.UserCounts {
			if uc.VariationId != expectedVIDs[idx] {
				t.Fatalf("variation ID is not correct: %s", uc.VariationId)
			}
			if len(uc.Timeseries.Timestamps) != 30 {
				t.Fatalf("the number of user counts is not correct: %d", len(uc.Timeseries.Timestamps))
			}
			if uc.Timeseries.Values[len(uc.Timeseries.Values)-1] != expectedUserVal[idx] {
				continue LOOP
			}
		}
		for idx, ec := range res.EventCounts {
			if ec.VariationId != expectedVIDs[idx] {
				t.Fatalf("variation ID is not correct: %s", ec.VariationId)
			}
			if len(ec.Timeseries.Timestamps) != 30 {
				t.Fatalf("the number of event counts is not correct: %d", len(ec.Timeseries.Timestamps))
			}
			if ec.Timeseries.Values[len(ec.Timeseries.Values)-1] != expectedEventVal[idx] {
				continue LOOP
			}
		}
		break
	}
}

func getVariationCount(vcs []*ecproto.VariationCount, id string) *ecproto.VariationCount {
	for _, vc := range vcs {
		if vc.VariationId == id {
			return vc
		}
	}
	return nil
}

func newExperimentClient(t *testing.T) experimentclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*orgOwnerDefaultAccessTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := experimentclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create experiment client:", err)
	}
	return client
}

func createGoals(ctx context.Context, t *testing.T, client experimentclient.Client, total int) []string {
	t.Helper()
	goalIDs := make([]string, 0)
	for i := 0; i < total; i++ {
		uuid := newUUID(t)
		resp, err := client.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
			Id:             createGoalID(t, uuid),
			Name:           createGoalID(t, uuid),
			Description:    fmt.Sprintf("%s-goal-description", prefixTestName),
			ConnectionType: experimentproto.Goal_EXPERIMENT,
			EnvironmentId:  *environmentID,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp.Goal == nil {
			t.Fatal("Goal is nil")
		}
		goalIDs = append(goalIDs, resp.Goal.Id)
	}
	return goalIDs
}

func createExperimentWithMultiGoals(
	ctx context.Context,
	t *testing.T,
	client experimentclient.Client,
	name string,
	featureID string,
	goalIDs []string,
	baseVariationID string,
	startAt, stopAt time.Time,
) *experimentproto.Experiment {
	resp, err := client.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{
		Name:            fmt.Sprintf("%s - %v", name, strings.Join(goalIDs, ",")),
		FeatureId:       featureID,
		GoalIds:         goalIDs,
		StartAt:         startAt.Unix(),
		StopAt:          stopAt.Unix(),
		BaseVariationId: baseVariationID,
		EnvironmentId:   *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		EnvironmentId: *environmentID,
		Id:            resp.Experiment.Id,
		Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
			Status: experimentproto.Experiment_RUNNING,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Update experiment cache
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		_, err = batchClient.ExecuteBatchJob(
			ctx,
			&btproto.BatchJobRequest{Job: btproto.BatchJob_ExperimentCacher})
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Unavailable {
			t.Fatalf("Failed to execute experiment cacher batch. Error code: %d. Error: %v\n", st.Code(), err)
		}
		fmt.Printf("Failed to execute experiment cacher batch. Error code: %d. Retrying in 5 seconds.\n", st.Code())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		t.Fatal(err)
	}
	return resp.Experiment
}

// This helper tries to stop the running experiments
// that are finished testing and waiting for deletion.
// This will improve the load on the http-stan while analysing the other experiments
// speeding up and improve timeout flaky tests.
// Since this is optional, it will ignore any errors.
func stopExperiment(
	ctx context.Context,
	t *testing.T,
	client experimentclient.Client,
	id string,
) {
	t.Helper()
	_, err := client.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		EnvironmentId: *environmentID,
		Id:            id,
		Status: &experimentproto.UpdateExperimentRequest_UpdatedStatus{
			Status: experimentproto.Experiment_FORCE_STOPPED,
		},
	})
	if err != nil {
		// Ignore
		return
	}
	// Update experiment cache
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		_, err = batchClient.ExecuteBatchJob(
			ctx,
			&btproto.BatchJobRequest{Job: btproto.BatchJob_ExperimentCacher})
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Unavailable {
			return
		}
		fmt.Printf("Failed to execute experiment cacher batch (Called by stopExperiment). Error code: %d. Retrying in 5 seconds.\n", st.Code())
		time.Sleep(5 * time.Second)
	}
}

func triggerExperimentCalculator(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	_, err := batchClient.ExecuteBatchJob(
		ctx,
		&btproto.BatchJobRequest{Job: btproto.BatchJob_ExperimentCalculator})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.ResourceExhausted {
			return
		}
		// Log but continue - this is best-effort
		if st.Code() == codes.DeadlineExceeded {
			t.Logf("Experiment calculator timed out after %v (best-effort)", grpcTimeout)
			return
		}
		t.Logf("Failed to trigger experiment calculator (best-effort). Error code: %d. Error: %v", st.Code(), err)
	}
}

func updateFeatueFlagCache(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	numRetries := 3
	var err error
	for i := 0; i < numRetries; i++ {
		_, err = batchClient.ExecuteBatchJob(
			ctx,
			&btproto.BatchJobRequest{Job: btproto.BatchJob_FeatureFlagCacher})
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Unavailable {
			t.Fatalf("Failed to execute feature flag cacher batch. Error code: %d. Error: %v\n", st.Code(), err)
		}
		fmt.Printf("Failed to execute feature flag cacher batch. Error code: %d. Retrying in 5 seconds.\n", st.Code())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func grpcRegisterGoalEvent(
	t *testing.T,
	goalID, userID, tag string,
	value float64,
	timestamp int64,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	goal, err := anypb.New(&eventproto.GoalEvent{
		Timestamp: timestamp,
		GoalId:    goalID,
		UserId:    userID,
		Value:     value,
		User: &userproto.User{
			Id:   userID,
			Data: map[string]string{"appVersion": "0.1.0"}},
		Tag: tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &gatewayproto.RegisterEventsRequest{
		Events: []*eventproto.Event{
			{
				Id:    newUUID(t),
				Event: goal,
			},
		},
	}
	response, err := c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func registerGoalEvent(
	t *testing.T,
	goalID, userID, tag string,
	value float64,
	timestamp int64,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	goal, err := protojson.Marshal(&eventproto.GoalEvent{
		Timestamp: timestamp,
		GoalId:    goalID,
		UserId:    userID,
		Value:     value,
		User: &userproto.User{
			Id:   userID,
			Data: map[string]string{"appVersion": "0.1.0"}},
		Tag: tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	events := []util.Event{
		{
			ID:    newUUID(t),
			Event: goal,
			Type:  gwapi.GoalEventType,
		},
	}
	response := util.RegisterEvents(t, events, *gatewayAddr, *apiKeyPath)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func sendHTTPTrack(t *testing.T, userID, goalID, tag string, value float64) {
	data, err := os.ReadFile(*apiKeyPath)
	if err != nil {
		t.Fatal(err)
	}
	url := fmt.Sprintf("https://%s/track?timestamp=%d&apikey=%s&userid=%s&goalid=%s&tag=%s&value=%f",
		*gatewayAddr,
		time.Now().Unix(),
		strings.TrimSpace(string(data)),
		userID, goalID, tag, value)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Send HTTP track request failed: %d", resp.StatusCode)
	}
}

func grpcRegisterEvaluationEvent(
	t *testing.T,
	featureID string,
	featureVersion int32,
	userID, variationID, tag string,
	reason *featureproto.Reason,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	if reason == nil {
		reason = &featureproto.Reason{}
	}
	evaluation, err := anypb.New(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User: &userproto.User{
			Id:   userID,
			Data: map[string]string{"appVersion": "0.1.0"}},
		Reason: reason,
		Tag:    tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := &gatewayproto.RegisterEventsRequest{
		Events: []*eventproto.Event{
			{
				Id:    newUUID(t),
				Event: evaluation,
			},
		},
	}
	response, err := c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func grpcRegisterEventsInChunks(t *testing.T, events []*eventproto.Event) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	const chunkSize = 50
	for start := 0; start < len(events); start += chunkSize {
		end := start + chunkSize
		if end > len(events) {
			end = len(events)
		}
		ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
		req := &gatewayproto.RegisterEventsRequest{Events: events[start:end]}
		response, err := c.RegisterEvents(ctx, req)
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		if len(response.Errors) > 0 {
			t.Fatalf("Failed to register events. Error: %v", response.Errors)
		}
	}
}

func grpcRegisterEvaluationEventsBatch(
	t *testing.T,
	featureID string,
	featureVersion int32,
	userIDs []string,
	variationID, tag string,
	reason *featureproto.Reason,
) {
	t.Helper()
	if reason == nil {
		reason = &featureproto.Reason{}
	}
	events := make([]*eventproto.Event, 0, len(userIDs))
	for _, userID := range userIDs {
		evaluation, err := anypb.New(&eventproto.EvaluationEvent{
			Timestamp:      time.Now().Unix(),
			FeatureId:      featureID,
			FeatureVersion: featureVersion,
			UserId:         userID,
			VariationId:    variationID,
			User: &userproto.User{
				Id:   userID,
				Data: map[string]string{"appVersion": "0.1.0"}},
			Reason: reason,
			Tag:    tag,
		})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, &eventproto.Event{
			Id:    newUUID(t),
			Event: evaluation,
		})
	}
	grpcRegisterEventsInChunks(t, events)
}

func grpcRegisterGoalEventsBatch(
	t *testing.T,
	goalID, tag string,
	userValues map[string]float64,
	timestamp int64,
) {
	t.Helper()
	events := make([]*eventproto.Event, 0, len(userValues))
	for userID, value := range userValues {
		goal, err := anypb.New(&eventproto.GoalEvent{
			Timestamp: timestamp,
			GoalId:    goalID,
			UserId:    userID,
			Value:     value,
			User: &userproto.User{
				Id:   userID,
				Data: map[string]string{"appVersion": "0.1.0"}},
			Tag: tag,
		})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, &eventproto.Event{
			Id:    newUUID(t),
			Event: goal,
		})
	}
	grpcRegisterEventsInChunks(t, events)
}

func registerEvaluationEvent(
	t *testing.T,
	featureID string,
	featureVersion int32,
	userID, variationID, tag string,
	reason *featureproto.Reason,
) {
	t.Helper()
	if reason == nil {
		reason = &featureproto.Reason{}
	}
	evaluation, err := protojson.Marshal(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User: &userproto.User{
			Id:   userID,
			Data: map[string]string{"appVersion": "0.1.0"}},
		Reason: reason,
		Tag:    tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	events := []util.Event{
		{
			ID:    newUUID(t),
			Event: evaluation,
			Type:  gwapi.EvaluationEventType,
		},
	}
	response := util.RegisterEvents(t, events, *gatewayAddr, *apiKeyPath)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

// registerEventsInChunks registers events in batches so large fixtures don't
// require one gateway round-trip per event.
func registerEventsInChunks(t *testing.T, events []util.Event) {
	t.Helper()
	const chunkSize = 50
	for start := 0; start < len(events); start += chunkSize {
		end := start + chunkSize
		if end > len(events) {
			end = len(events)
		}
		response := util.RegisterEvents(t, events[start:end], *gatewayAddr, *apiKeyPath)
		if len(response.Errors) > 0 {
			t.Fatalf("Failed to register events. Error: %v", response.Errors)
		}
	}
}

// registerEvaluationEventsBatch registers one evaluation event per user (all
// mapped to variationID) in batched gateway calls.
func registerEvaluationEventsBatch(
	t *testing.T,
	featureID string,
	featureVersion int32,
	userIDs []string,
	variationID, tag string,
	reason *featureproto.Reason,
) {
	t.Helper()
	if reason == nil {
		reason = &featureproto.Reason{}
	}
	events := make([]util.Event, 0, len(userIDs))
	for _, userID := range userIDs {
		evaluation, err := protojson.Marshal(&eventproto.EvaluationEvent{
			Timestamp:      time.Now().Unix(),
			FeatureId:      featureID,
			FeatureVersion: featureVersion,
			UserId:         userID,
			VariationId:    variationID,
			User: &userproto.User{
				Id:   userID,
				Data: map[string]string{"appVersion": "0.1.0"}},
			Reason: reason,
			Tag:    tag,
		})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, util.Event{
			ID:    newUUID(t),
			Event: evaluation,
			Type:  gwapi.EvaluationEventType,
		})
	}
	registerEventsInChunks(t, events)
}

// registerGoalEventsBatch registers one goal event per user (value taken from
// userValues) in batched gateway calls.
func registerGoalEventsBatch(
	t *testing.T,
	goalID, tag string,
	userValues map[string]float64,
	timestamp int64,
) {
	t.Helper()
	events := make([]util.Event, 0, len(userValues))
	for userID, value := range userValues {
		goal, err := protojson.Marshal(&eventproto.GoalEvent{
			Timestamp: timestamp,
			GoalId:    goalID,
			UserId:    userID,
			Value:     value,
			User: &userproto.User{
				Id:   userID,
				Data: map[string]string{"appVersion": "0.1.0"}},
			Tag: tag,
		})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, util.Event{
			ID:    newUUID(t),
			Event: goal,
			Type:  gwapi.GoalEventType,
		})
	}
	registerEventsInChunks(t, events)
}

// addFeatureIndividualTargetingBulk sets the full individual-targeting user
// list for a variation in a single UpdateFeature call (instead of one call per
// user, which is prohibitive for large fixtures).
func addFeatureIndividualTargetingBulk(
	t *testing.T,
	featureID, variationID string,
	users []string,
	client featureclient.Client,
) {
	t.Helper()
	for i := 0; i < deadlockRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
		_, err := client.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
			Id:            featureID,
			EnvironmentId: *environmentID,
			TargetChanges: []*featureproto.TargetChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Target: &featureproto.Target{
						Variation: variationID,
						Users:     users,
					},
				},
			},
		})
		cancel()
		if err == nil {
			return
		}
		if i < deadlockRetryAttempts-1 && util.IsDeadlockError(err) {
			t.Logf("Retrying addFeatureIndividualTargetingBulk (attempt %d/%d) for %s: %v",
				i+1, deadlockRetryAttempts, featureID, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		t.Fatalf("Failed to bulk add individual targeting for feature %s: %v", featureID, err)
	}
}

func makeExperimentResultUserIDs(t *testing.T, uuid string) (userIDsA, userIDsB []string) {
	t.Helper()
	base := createUserID(t, uuid)
	userIDsA = make([]string, 0, experimentResultUsersPerVariation)
	userIDsB = make([]string, 0, experimentResultUsersPerVariation)
	for i := 0; i < experimentResultUsersPerVariation; i++ {
		userIDsA = append(userIDsA, fmt.Sprintf("%s-a-%d", base, i))
		userIDsB = append(userIDsB, fmt.Sprintf("%s-b-%d", base, i))
	}
	return userIDsA, userIDsB
}

func goalValuesForUsers(userIDs []string, value float64) map[string]float64 {
	out := make(map[string]float64, len(userIDs))
	for i, userID := range userIDs {
		// Alternate by 0.1 so VAR_SAMP > 0 (the calculator skips value metrics
		// when per-variation variance is exactly zero). The spread is tiny
		// relative to the 10 vs 15 arm separation and stays well below any
		// winsorization cap.
		out[userID] = value + float64(i%2)*0.1
	}
	return out
}

func experimentResultExpectedValueSum(baseValue float64) float64 {
	n := float64(experimentResultUsersPerVariation)
	half := n / 2
	return baseValue*n + 0.1*half
}

func addFeatureIndividualTargetingBulkSplit(
	t *testing.T,
	featureID, variationAID, variationBID string,
	usersA, usersB []string,
	client featureclient.Client,
) {
	t.Helper()
	for i := 0; i < deadlockRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
		_, err := client.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
			Id:            featureID,
			EnvironmentId: *environmentID,
			TargetChanges: []*featureproto.TargetChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Target: &featureproto.Target{
						Variation: variationAID,
						Users:     usersA,
					},
				},
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Target: &featureproto.Target{
						Variation: variationBID,
						Users:     usersB,
					},
				},
			},
		})
		cancel()
		if err == nil {
			return
		}
		if i < deadlockRetryAttempts-1 && util.IsDeadlockError(err) {
			t.Logf("Retrying addFeatureIndividualTargetingBulkSplit (attempt %d/%d) for %s: %v",
				i+1, deadlockRetryAttempts, featureID, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		t.Fatalf("Failed to bulk add split individual targeting for feature %s: %v", featureID, err)
	}
}

func experimentResultCountsReady(vsA, vsB *ecproto.VariationResult) bool {
	if vsA.EvaluationCount == nil || vsA.ExperimentCount == nil ||
		vsB.EvaluationCount == nil || vsB.ExperimentCount == nil {
		return false
	}
	n := int64(experimentResultUsersPerVariation)
	return vsA.EvaluationCount.EventCount == n &&
		vsA.EvaluationCount.UserCount == n &&
		vsB.EvaluationCount.EventCount == n &&
		vsB.EvaluationCount.UserCount == n &&
		vsA.ExperimentCount.EventCount == n &&
		vsA.ExperimentCount.UserCount == n &&
		vsB.ExperimentCount.EventCount == n &&
		vsB.ExperimentCount.UserCount == n
}

// experimentResultProbFieldsReady reports whether the calculator has populated
// the CVR and value-metric distribution summaries for a variation.
func experimentResultProbFieldsReady(vr *ecproto.VariationResult) bool {
	return vr.CvrProbBest != nil &&
		vr.CvrProbBeatBaseline != nil &&
		vr.GoalValueSumPerUserProbBest != nil &&
		vr.GoalValueSumPerUserProbBeatBaseline != nil
}

func requireDistributionSummary(
	t *testing.T,
	variationValue, fieldName string,
	ds *ecproto.DistributionSummary,
) *ecproto.DistributionSummary {
	t.Helper()
	if ds == nil {
		t.Fatalf("variation: %s: %s should be populated, got nil", variationValue, fieldName)
	}
	return ds
}

func requireVariationCounts(t *testing.T, variationValue string, vr *ecproto.VariationResult) {
	t.Helper()
	if vr.EvaluationCount == nil {
		t.Fatalf("variation: %s: evaluation count should be populated, got nil", variationValue)
	}
	if vr.ExperimentCount == nil {
		t.Fatalf("variation: %s: experiment count should be populated, got nil", variationValue)
	}
}

func getVariationResult(vrs []*ecproto.VariationResult, id string) *ecproto.VariationResult {
	for _, vr := range vrs {
		if vr.VariationId == id {
			return vr
		}
	}
	return nil
}

func waitAndCheckExperimentResult(
	t *testing.T,
	ecClient ecclient.Client,
	experiment *experimentproto.Experiment,
	goalID string,
) {
	t.Helper()
	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout after %d attempts", retryTimes)
		}
		time.Sleep(10 * time.Second)
		triggerExperimentCalculator(t)

		resp, err := getExperimentResult(t, ecClient, experiment.Id)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment result. Error code: %d. Error: %v\n", st.Code(), err)
			}
			t.Logf("Retry %d/%d: ExperimentResult not found yet (NotFound error)", i+1, retryTimes)
			continue
		}
		if resp == nil {
			continue
		}
		er := resp.ExperimentResult
		if er.Id != experiment.Id {
			t.Fatalf("experiment ID is not correct: %s", er.Id)
		}
		if len(er.GoalResults) == 0 {
			continue
		}
		if len(er.GoalResults) != 1 {
			t.Fatalf("the number of goal results is not correct: %d", len(er.GoalResults))
		}
		gr := er.GoalResults[0]
		if gr.GoalId != goalID {
			t.Fatalf("goal ID is not correct: %s", gr.GoalId)
		}
		if len(gr.VariationResults) != 2 {
			t.Fatalf("the number of variation results is not correct: %d", len(gr.VariationResults))
		}
		vsA := getVariationResult(gr.VariationResults, experiment.Variations[0].Id)
		vsB := getVariationResult(gr.VariationResults, experiment.Variations[1].Id)
		if vsA == nil || vsB == nil {
			t.Fatalf("missing variation result for experiment variations")
		}
		if !experimentResultCountsReady(vsA, vsB) {
			linkedA, linkedB := int64(0), int64(0)
			if vsA.ExperimentCount != nil {
				linkedA = vsA.ExperimentCount.UserCount
			}
			if vsB.ExperimentCount != nil {
				linkedB = vsB.ExperimentCount.UserCount
			}
			t.Logf("Retry %d/%d: waiting for linked users A=%d/%d B=%d/%d",
				i+1, retryTimes, linkedA, experimentResultUsersPerVariation,
				linkedB, experimentResultUsersPerVariation)
			continue
		}
		if !experimentResultProbFieldsReady(vsA) || !experimentResultProbFieldsReady(vsB) {
			t.Logf("Retry %d/%d: waiting for calculator prob fields on both variations", i+1, retryTimes)
			continue
		}
		checkExperimentVariationResultA(t, vsA, experiment.Variations[0].Value)
		checkExperimentVariationResultB(t, vsB, experiment.Variations[1].Value)
		checkExperimentSummary(t, gr, experiment.Variations[1].Id)
		checkExperimentSrmResult(t, er)
		return
	}
}

func newGatewayClient(t *testing.T) gatewayclient.Client {
	t.Helper()
	creds, err := gatewayclient.NewPerRPCCredentials(*apiKeyPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := gatewayclient.NewClient(
		fmt.Sprintf("%s:%d", *gatewayAddr, *gatewayPort),
		*gatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create gateway client:", err)
	}
	return client
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func newFeatureClient(t *testing.T) featureclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*orgOwnerDefaultAccessTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	featureClient, err := featureclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create feature client:", err)
	}
	return featureClient
}

func newEventCounterClient(t *testing.T) ecclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*orgOwnerDefaultAccessTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := ecclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create experiment client:", err)
	}
	return client
}

func newBatchClient(t *testing.T) btclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*orgOwnerDefaultAccessTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := btclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create batch client:", err)
	}
	return client
}

func newCreateFeatureReq(featureID string, variations []string) *featureproto.CreateFeatureRequest {
	req := &featureproto.CreateFeatureRequest{
		Id:          featureID,
		Name:        featureID,
		Description: "e2e-test-eventcounter-feature-description",
		Variations:  []*featureproto.Variation{},
		Tags: []string{
			"e2e-test-tag-1",
			"e2e-test-tag-2",
			"e2e-test-tag-3",
		},
		DefaultOnVariationIndex:  &wrapperspb.Int32Value{Value: int32(0)},
		DefaultOffVariationIndex: &wrapperspb.Int32Value{Value: int32(1)},
		EnvironmentId:            *environmentID,
	}
	for _, v := range variations {
		req.Variations = append(req.Variations, &featureproto.Variation{
			Value:       v,
			Name:        "Variation " + v,
			Description: "Thing does " + v,
		})
	}
	return req
}

func createFeature(
	t *testing.T,
	client featureclient.Client,
	featureID, tag, variationA, variationB string,
) {
	t.Helper()
	createReq := newCreateFeatureReq(featureID, []string{variationA, variationB})
	for i := 0; i < deadlockRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
		_, err := client.CreateFeature(ctx, createReq)
		cancel()
		if err == nil {
			break
		}
		if i < deadlockRetryAttempts-1 && util.IsDeadlockError(err) {
			t.Logf("Retrying createFeature (attempt %d/%d) for %s: %v", i+1, deadlockRetryAttempts, featureID, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		t.Fatal(err)
	}
	addTag(t, tag, featureID, client)
	enableFeature(t, featureID, client)
}

func addTag(t *testing.T, tag string, featureID string, client featureclient.Client) {
	t.Helper()
	addReq := &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
		TagChanges: []*featureproto.TagChange{
			{
				ChangeType: featureproto.ChangeType_CREATE,
				Tag:        tag,
			},
		},
	}
	for i := 0; i < deadlockRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
		_, err := client.UpdateFeature(ctx, addReq)
		cancel()
		if err == nil {
			return
		}
		if i < deadlockRetryAttempts-1 && util.IsDeadlockError(err) {
			t.Logf("Retrying addTag (attempt %d/%d) for %s: %v", i+1, deadlockRetryAttempts, featureID, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		t.Fatal(err)
	}
}

func enableFeature(t *testing.T, featureID string, client featureclient.Client) {
	t.Helper()
	enableReq := &featureproto.UpdateFeatureRequest{
		Id:            featureID,
		Enabled:       wrapperspb.Bool(true),
		EnvironmentId: *environmentID,
	}
	for i := 0; i < deadlockRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
		_, err := client.UpdateFeature(ctx, enableReq)
		cancel()
		if err == nil {
			return
		}
		if i < deadlockRetryAttempts-1 && util.IsDeadlockError(err) {
			t.Logf("Retrying enableFeature (attempt %d/%d) for %s: %v", i+1, deadlockRetryAttempts, featureID, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		t.Fatalf("Failed to enable feature id: %s. Error: %v", featureID, err)
	}
}

func addFeatureIndividualTargeting(t *testing.T, featureID, userID, variationID string, client featureclient.Client) {
	t.Helper()
	f, err := getFeature(t, client, featureID)
	if err != nil {
		t.Fatalf("Failed to get feature. ID: %s. Error: %v", featureID, err)
	}
	var users []string
	for _, v := range f.Targets {
		if v.Variation == variationID {
			users = append(v.Users, userID)
			break
		}
	}
	for i := 0; i < deadlockRetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
		_, err = client.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
			Id:            featureID,
			EnvironmentId: *environmentID,
			TargetChanges: []*featureproto.TargetChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Target: &featureproto.Target{
						Variation: variationID,
						Users:     users,
					},
				},
			},
		})
		cancel()
		if err == nil {
			return
		}
		if i < deadlockRetryAttempts-1 && util.IsDeadlockError(err) {
			t.Logf("Retrying addFeatureIndividualTargeting (attempt %d/%d) for %s: %v", i+1, deadlockRetryAttempts, featureID, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		t.Fatalf("Failed to add individual targeting for feature %s: %v", featureID, err)
	}
}

func getEvaluation(t *testing.T, tag string, userID string) (*gatewayproto.GetEvaluationsResponse, error) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	req := &gatewayproto.GetEvaluationsRequest{
		Tag:  tag,
		User: &userproto.User{Id: userID},
	}
	var response *gatewayproto.GetEvaluationsResponse
	var err error
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		response, err = c.GetEvaluations(ctx, req)
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Unavailable {
			return nil, err
		}
		fmt.Printf("Failed to get evaluations. Error code: %d. Retrying in 5 seconds.\n", st.Code())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getExperiment(t *testing.T, c experimentclient.Client, id string) (*experimentproto.GetExperimentResponse, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	req := &experimentproto.GetExperimentRequest{
		EnvironmentId: *environmentID,
		Id:            id,
	}
	var response *experimentproto.GetExperimentResponse
	var err error
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		response, err = c.GetExperiment(ctx, req)
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Unavailable {
			return nil, err
		}
		fmt.Printf("Failed to get experiment. Experiment ID: %s. Error code: %d. Retrying in 5 seconds.\n", id, st.Code())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getExperimentResult(t *testing.T, c ecclient.Client, experimentID string) (*ecproto.GetExperimentResultResponse, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	req := &ecproto.GetExperimentResultRequest{
		ExperimentId:  experimentID,
		EnvironmentId: *environmentID,
	}
	var response *ecproto.GetExperimentResultResponse
	var err error
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		response, err = c.GetExperimentResult(ctx, req)
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		// Retry on transient BigQuery emulator errors:
		// - "database table is locked": SQLite backend used by the emulator can't handle concurrent writes
		// - "job is not found": job may not be registered yet in the emulator
		if st.Code() == codes.Unavailable || st.Code() == codes.Internal ||
			(st.Code() == codes.Unknown && (strings.Contains(err.Error(), "database table is locked") ||
				strings.Contains(err.Error(), "job") && strings.Contains(err.Error(), "is not found"))) {
			fmt.Printf("Failed to get experiment result. Experiment ID: %s. Error code: %d. Retrying in 5 seconds.\n", experimentID, st.Code())
			time.Sleep(5 * time.Second)
			continue
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getExperimentEvaluationCount(
	t *testing.T, c ecclient.Client, featureID string, featureVersion int32, variationIDs []string,
) (*ecproto.GetExperimentEvaluationCountResponse, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	now := time.Now()
	req := &ecproto.GetExperimentEvaluationCountRequest{
		EnvironmentId:  *environmentID,
		StartAt:        now.Add(-30 * 24 * time.Hour).Unix(),
		EndAt:          now.Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		VariationIds:   variationIDs,
	}
	var response *ecproto.GetExperimentEvaluationCountResponse
	var err error
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		response, err = c.GetExperimentEvaluationCount(ctx, req)
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		// Retry on transient BigQuery emulator errors:
		// - "database table is locked": SQLite backend used by the emulator can't handle concurrent writes
		// - "job is not found": job may not be registered yet in the emulator
		if st.Code() == codes.Unavailable || st.Code() == codes.Internal ||
			(st.Code() == codes.Unknown && (strings.Contains(err.Error(), "database table is locked") ||
				strings.Contains(err.Error(), "job") && strings.Contains(err.Error(), "is not found"))) {
			fmt.Printf("Failed to get experiment evaluation count. Error code: %d. Retrying in 5 seconds.\n", st.Code())
			time.Sleep(5 * time.Second)
			continue
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getExperimentGoalCount(
	t *testing.T, c ecclient.Client, goalID, featureID string, featureVersion int32, variationIDs []string,
) (*ecproto.GetExperimentGoalCountResponse, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	now := time.Now()
	req := &ecproto.GetExperimentGoalCountRequest{
		EnvironmentId:  *environmentID,
		StartAt:        now.Add(-30 * 24 * time.Hour).Unix(),
		EndAt:          now.Unix(),
		GoalId:         goalID,
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		VariationIds:   variationIDs,
	}
	var response *ecproto.GetExperimentGoalCountResponse
	var err error
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		response, err = c.GetExperimentGoalCount(ctx, req)
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		// Retry on transient BigQuery emulator errors:
		// - "database table is locked": SQLite backend used by the emulator can't handle concurrent writes
		// - "job is not found": job may not be registered yet in the emulator
		if st.Code() == codes.Unavailable || st.Code() == codes.Internal ||
			(st.Code() == codes.Unknown && (strings.Contains(err.Error(), "database table is locked") ||
				strings.Contains(err.Error(), "job") && strings.Contains(err.Error(), "is not found"))) {
			fmt.Printf("Failed to get experiment goal count. Error code: %d. Retrying in 5 seconds.\n", st.Code())
			time.Sleep(5 * time.Second)
			continue
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getFeature(t *testing.T, client featureclient.Client, featureID string) (*featureproto.Feature, error) {
	t.Helper()
	getReq := &featureproto.GetFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	var response *featureproto.GetFeatureResponse
	var err error
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		response, err = client.GetFeature(ctx, getReq)
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Unavailable {
			return nil, err
		}
		fmt.Printf("Failed to get feature. ID: %s. Error code: %d. Retrying in 5 seconds.\n", featureID, st.Code())
		time.Sleep(5 * time.Second)
	}
	return response.Feature, err
}

func getEvaluationTimeseriesCount(
	t *testing.T,
	featureID string,
	c ecclient.Client,
	timeRange ecproto.GetEvaluationTimeseriesCountRequest_TimeRange,
) (*ecproto.GetEvaluationTimeseriesCountResponse, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()
	req := &ecproto.GetEvaluationTimeseriesCountRequest{
		EnvironmentId: *environmentID,
		FeatureId:     featureID,
		TimeRange:     timeRange,
	}
	var response *ecproto.GetEvaluationTimeseriesCountResponse
	var err error
	numRetries := 3
	for i := 0; i < numRetries; i++ {
		response, err = c.GetEvaluationTimeseriesCount(
			ctx,
			req,
		)
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		// Retry on transient BigQuery emulator errors:
		// - "database table is locked": SQLite backend used by the emulator can't handle concurrent writes
		// - "job is not found": job may not be registered yet in the emulator
		if st.Code() == codes.Unavailable || st.Code() == codes.Internal ||
			(st.Code() == codes.Unknown && (strings.Contains(err.Error(), "database table is locked") ||
				strings.Contains(err.Error(), "job") && strings.Contains(err.Error(), "is not found"))) {
			fmt.Printf("Failed to get evaluation timeseries count. ID: %s. Error code: %d. Retrying in 5 seconds.\n", featureID, st.Code())
			time.Sleep(5 * time.Second)
			continue
		}
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func createFeatureID(t *testing.T, uuid string) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, uuid)
}

func createGoalID(t *testing.T, uuid string) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-goal-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-goal-id-%s", prefixTestName, uuid)
}

func createUserID(t *testing.T, uuid string) string {
	t.Helper()
	if *testID != "" {
		return fmt.Sprintf("%s-%s-user-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-user-id-%s", prefixTestName, uuid)
}

// checkExperimentVariationResultA validates variation A for the 50-user-per-arm
// experiment-result fixture (baseline: per-user value 10 vs treatment 15).
func checkExperimentVariationResultA(t *testing.T, vsA *ecproto.VariationResult, variationValue string) {
	t.Helper()
	requireVariationCounts(t, variationValue, vsA)
	n := int64(experimentResultUsersPerVariation)
	valueSum := experimentResultExpectedValueSum(experimentResultValueVariationA)
	if vsA.EvaluationCount.EventCount != n {
		t.Fatalf("variation: %s: evaluation event count is not correct: %d", variationValue, vsA.EvaluationCount.EventCount)
	}
	if vsA.EvaluationCount.UserCount != n {
		t.Fatalf("variation: %s: evaluation user count is not correct: %d", variationValue, vsA.EvaluationCount.UserCount)
	}
	if vsA.ExperimentCount.EventCount != n {
		t.Fatalf("variation: %s: experiment event count is not correct: %d", variationValue, vsA.ExperimentCount.EventCount)
	}
	if vsA.ExperimentCount.UserCount != n {
		t.Fatalf("variation: %s: experiment user count is not correct: %d", variationValue, vsA.ExperimentCount.UserCount)
	}
	if diff := cmp.Diff(vsA.ExperimentCount.ValueSum, valueSum, compareFloatOpt); diff != "" {
		t.Fatalf("variation: %s: experiment value sum is not correct: %f", variationValue, vsA.ExperimentCount.ValueSum)
	}
	cvrProbBest := requireDistributionSummary(t, variationValue, "cvr_prob_best", vsA.CvrProbBest)
	if diff := cmp.Diff(cvrProbBest.Mean, 0.50, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", variationValue, cvrProbBest.Mean)
	}
	if diff := cmp.Diff(cvrProbBest.Sd, 0.50, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", variationValue, cvrProbBest.Sd)
	}
	if diff := cmp.Diff(cvrProbBest.Rhat, 0.99, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", variationValue, cvrProbBest.Rhat)
	}
	cvrProbBeatBaseline := requireDistributionSummary(t, variationValue, "cvr_prob_beat_baseline", vsA.CvrProbBeatBaseline)
	if diff := cmp.Diff(cvrProbBeatBaseline.Mean, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", variationValue, cvrProbBeatBaseline.Mean)
	}
	if diff := cmp.Diff(cvrProbBeatBaseline.Sd, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", variationValue, cvrProbBeatBaseline.Sd)
	}
	if diff := cmp.Diff(cvrProbBeatBaseline.Rhat, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", variationValue, cvrProbBeatBaseline.Rhat)
	}
	valueProbBest := requireDistributionSummary(t, variationValue, "goal_value_sum_per_user_prob_best", vsA.GoalValueSumPerUserProbBest)
	if diff := cmp.Diff(valueProbBest.Mean, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", variationValue, valueProbBest.Mean)
	}
	valueProbBeatBaseline := requireDistributionSummary(t, variationValue, "goal_value_sum_per_user_prob_beat_baseline", vsA.GoalValueSumPerUserProbBeatBaseline)
	if diff := cmp.Diff(valueProbBeatBaseline.Mean, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", variationValue, valueProbBeatBaseline.Mean)
	}
}

// checkExperimentVariationResultB validates variation B for the 50-user-per-arm
// experiment-result fixture (treatment: per-user value 15 vs baseline 10).
func checkExperimentVariationResultB(t *testing.T, vsB *ecproto.VariationResult, variationValue string) {
	t.Helper()
	requireVariationCounts(t, variationValue, vsB)
	n := int64(experimentResultUsersPerVariation)
	valueSum := experimentResultExpectedValueSum(experimentResultValueVariationB)
	if vsB.EvaluationCount.EventCount != n {
		t.Fatalf("variation: %s: evaluation event count is not correct: %d", variationValue, vsB.EvaluationCount.EventCount)
	}
	if vsB.EvaluationCount.UserCount != n {
		t.Fatalf("variation: %s: evaluation user count is not correct: %d", variationValue, vsB.EvaluationCount.UserCount)
	}
	if vsB.ExperimentCount.EventCount != n {
		t.Fatalf("variation: %s: experiment event count is not correct: %d", variationValue, vsB.ExperimentCount.EventCount)
	}
	if vsB.ExperimentCount.UserCount != n {
		t.Fatalf("variation: %s: experiment user count is not correct: %d", variationValue, vsB.ExperimentCount.UserCount)
	}
	if diff := cmp.Diff(vsB.ExperimentCount.ValueSum, valueSum, compareFloatOpt); diff != "" {
		t.Fatalf("variation: %s: experiment value sum is not correct: %f", variationValue, vsB.ExperimentCount.ValueSum)
	}
	cvrProbBest := requireDistributionSummary(t, variationValue, "cvr_prob_best", vsB.CvrProbBest)
	if diff := cmp.Diff(cvrProbBest.Mean, 0.50, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", variationValue, cvrProbBest.Mean)
	}
	if diff := cmp.Diff(cvrProbBest.Sd, 0.50, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", variationValue, cvrProbBest.Sd)
	}
	if diff := cmp.Diff(cvrProbBest.Rhat, 0.99, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", variationValue, cvrProbBest.Rhat)
	}
	cvrProbBeatBaseline := requireDistributionSummary(t, variationValue, "cvr_prob_beat_baseline", vsB.CvrProbBeatBaseline)
	if diff := cmp.Diff(cvrProbBeatBaseline.Mean, 0.50, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", variationValue, cvrProbBeatBaseline.Mean)
	}
	if diff := cmp.Diff(cvrProbBeatBaseline.Sd, 0.50, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", variationValue, cvrProbBeatBaseline.Sd)
	}
	if diff := cmp.Diff(cvrProbBeatBaseline.Rhat, 0.99, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", variationValue, cvrProbBeatBaseline.Rhat)
	}
	valueProbBest := requireDistributionSummary(t, variationValue, "goal_value_sum_per_user_prob_best", vsB.GoalValueSumPerUserProbBest)
	if diff := cmp.Diff(valueProbBest.Mean, 1.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", variationValue, valueProbBest.Mean)
	}
	valueProbBeatBaseline := requireDistributionSummary(t, variationValue, "goal_value_sum_per_user_prob_beat_baseline", vsB.GoalValueSumPerUserProbBeatBaseline)
	if diff := cmp.Diff(valueProbBeatBaseline.Mean, 1.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", variationValue, valueProbBeatBaseline.Mean)
	}
}

// checkExperimentSummary asserts the metric-aware best-variations lists end to
// end (Follow-up E). With this fixture the value-per-user metric is decisive
// (treatment ~15 vs baseline ~10) while CVR is a tie (100% conversion in both
// arms), so the value list must contain the treatment as the winner and the
// CVR list must be empty. This exercises the `best_variations_value` proto
// field through calculator → storage → API, complementing the in-process
// TestPickBestVariations unit coverage.
func checkExperimentSummary(t *testing.T, gr *ecproto.GoalResult, treatmentVariationID string) {
	t.Helper()
	if gr.Summary == nil {
		t.Fatalf("goal result summary should be populated, got nil")
	}
	// CVR is a tie, so no variation clears the 0.95 best-variation threshold.
	if len(gr.Summary.BestVariations) != 0 {
		t.Fatalf("expected no CVR best variations (CVR is a tie), got %d", len(gr.Summary.BestVariations))
	}
	// The value metric is decisive: the treatment is the value winner.
	if len(gr.Summary.BestVariationsValue) != 1 {
		t.Fatalf("expected exactly 1 value best variation, got %d", len(gr.Summary.BestVariationsValue))
	}
	best := gr.Summary.BestVariationsValue[0]
	if best.Id != treatmentVariationID {
		t.Fatalf("value winner should be treatment %s, got %s", treatmentVariationID, best.Id)
	}
	if !best.IsBest {
		t.Fatalf("value winner should be marked IsBest")
	}
	if best.Probability < 0.95 {
		t.Fatalf("value winner probability should exceed 0.95, got %f", best.Probability)
	}
}

// checkExperimentSrmResult asserts that the SRM check ran end-to-end (proto
// field populated by the calculator and serialized through the API).
//
// The fixture above creates the feature via CreateFeatureRequest with
// DefaultOnVariationIndex/DefaultOffVariationIndex, which produces a FIXED
// default strategy (no rollout weights to test against). The calculator's
// SRM path therefore short-circuits in extractRolloutWeights with
// "default strategy is not a rollout" before any sample-size check runs, so
// the expected outcome here is Status == SKIPPED with a non-empty
// skip_reason. (The minSRMSampleSize=100 floor is also satisfied by this
// fixture, but the strategy check wins first.)
//
// This is enough to catch wiring regressions (feature client not injected,
// proto field not exposed, calculator path not executed) without requiring
// a rollout-strategy feature or a large fixture; the in-process tests in
// srm_test.go cover the per-branch SKIPPED reasons directly.
func checkExperimentSrmResult(t *testing.T, er *ecproto.ExperimentResult) {
	t.Helper()
	if er.SrmResult == nil {
		t.Fatalf("SRM result should be populated on every ExperimentResult, got nil")
	}
	if er.SrmResult.Status != ecproto.SrmResult_SKIPPED {
		t.Fatalf("with a FIXED default strategy, SRM status should be SKIPPED, got %v (skip_reason=%q)",
			er.SrmResult.Status, er.SrmResult.SkipReason)
	}
	if er.SrmResult.SkipReason == "" {
		t.Fatalf("SRM SKIPPED status must come with a non-empty skip_reason")
	}
	if diff := cmp.Diff(er.SrmResult.Threshold, 0.001, compareFloatOpt); diff != "" {
		t.Fatalf("SRM threshold should default to 0.001, got %f", er.SrmResult.Threshold)
	}
}
