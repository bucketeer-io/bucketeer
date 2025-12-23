// Copyright 2025 The Bucketeer Authors.
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

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

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
	prefixTestName = "e2e-test"
	timeout        = 3 * time.Minute
	retryTimes     = 50
)

const defaultVariationID = "default"

var (
	webGatewayAddr       = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort       = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert       = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath           = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath     = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr          = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort          = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert          = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath     = flag.String("service-token", "", "Service token path")
	environmentID        = flag.String("environment-id", "", "Environment id")
	organizationID       = flag.String("organization-id", "", "Organization ID")
	testID               = flag.String("test-id", "", "test ID")
	compareFloatOpt      = cmpopts.EquateApprox(0, 0.0001)
	compareFloatBayesian = cmpopts.EquateApprox(0, 0.03) // 3% tolerance for Bayesian results to avoid flaky tests
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
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

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
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

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
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	createFeature(t, featureClient, featureID, tag, "a", "b")
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
	addFeatureIndividualTargeting(t, featureID, userIDs[0], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[1], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[2], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[3], f.Variations[1].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[4], f.Variations[1].Id, featureClient)

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
	stopAt := time.Now().Add(2 * time.Hour)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	// Wait for event-persister-dwh to subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

	// CVRs is 3/4
	// Evaluation events must always be sent before goal events
	// Register 3 events and 2 user counts for the user index 1, 2 and 3
	// Register variation a
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	// where the feature version changes between getFeature() and createExperiment()
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[0], experiment.Variations[0].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[1], experiment.Variations[0].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[2], experiment.Variations[0].Id, tag, reason)
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[0], experiment.Variations[0].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	// Register goal variation
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[1], tag, float64(0.2), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Add(-2*time.Hour).Unix())
	// Increment experiment event count
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Unix())

	// CVRs is 2/3
	// Evaluation events must always be sent before goal events
	// Register 3 events and 2 user counts for the user index 4 and 5
	// Register variation
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[3], experiment.Variations[1].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[4], experiment.Variations[1].Id, tag, reason)
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[3], experiment.Variations[1].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	// Register goal
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Add(-2*time.Hour).Unix())
	// Increment experiment event count
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Unix())

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout after %d attempts", retryTimes)
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentResult(t, ecClient, experiment.Id)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment result. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				t.Logf("Retry %d/%d: ExperimentResult not found yet (NotFound error)", i+1, retryTimes)
				continue
			}
		}

		if resp != nil {
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
			if gr.GoalId != goalIDs[0] {
				t.Fatalf("goal ID is not correct: %s", gr.GoalId)
			}
			if len(gr.VariationResults) != 2 {
				t.Fatalf("the number of variation results is not correct: %d", len(gr.VariationResults))
			}
			vsA := gr.VariationResults[0]
			vsB := gr.VariationResults[1]
			// These counts are based on the number of events sent earlier
			if vsA.EvaluationCount.EventCount != 4 || // variation A
				vsA.EvaluationCount.UserCount != 3 ||
				vsB.EvaluationCount.EventCount != 3 || // variation B
				vsB.EvaluationCount.UserCount != 2 {
				continue
			}
			// These counts are based on the number of events sent earlier
			if vsA.ExperimentCount.EventCount != 5 || // variation A
				vsA.ExperimentCount.UserCount != 3 ||
				vsB.ExperimentCount.EventCount != 4 || // variation B
				vsB.ExperimentCount.UserCount != 2 {
				continue
			}
			checkExperimentVariationResultA(t, vsA, experiment.Variations[0].Value)
			checkExperimentVariationResultB(t, vsB, experiment.Variations[1].Value)
			break
		}
	}
	res, err := getExperiment(t, experimentClient, experiment.Id)
	if err != nil {
		t.Fatalf("Failed to get experiment. ID: %s. Error: %v", experiment.Id, err)
	}
	if res.Experiment.Status != experimentproto.Experiment_RUNNING {
		expected, _ := experimentproto.Experiment_Status_name[int32(experimentproto.Experiment_RUNNING)]
		actual, _ := experimentproto.Experiment_Status_name[int32(res.Experiment.Status)]
		t.Fatalf("the status of experiment is not correct. expected: %s, but got %s", expected, actual)
	}
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	createFeature(t, featureClient, featureID, tag, "a", "b")
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
	addFeatureIndividualTargeting(t, featureID, userIDs[1], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[2], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[3], f.Variations[1].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[4], f.Variations[1].Id, featureClient)

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
	stopAt := time.Now().Add(2 * time.Hour)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

	// CVRs is 3/4
	// Evaluation events must always be sent before goal events
	// Register 3 events and 2 user counts for user 1, 2 and 3
	// Register variation a
	// IMPORTANT: Use experiment.FeatureVersion instead of f.Version to avoid race condition
	// where the feature version changes between getFeature() and createExperiment()
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[0], f.Variations[0].Id, tag, reason)
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[1], f.Variations[0].Id, tag, reason)
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[2], f.Variations[0].Id, tag, reason)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[0], f.Variations[0].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	// Register goal variation
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[1], tag, float64(0.2), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Add(-2*time.Hour).Unix())
	// Increment experiment event count
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Unix())

	// CVRs is 2/3
	// Evaluation events must always be sent before goal events
	// Register 3 events and 2 user counts for user 4 and 5
	// Register variation
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[3], f.Variations[1].Id, tag, reason)
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[4], f.Variations[1].Id, tag, reason)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, experiment.FeatureVersion, userIDs[3], f.Variations[1].Id, tag, reason)

	// Wait a few seconds so the data in BigQuery becomes available for linking the goal event
	time.Sleep(10 * time.Second)

	// Register goal
	registerGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Add(-2*time.Hour).Unix())
	// Increment experiment event count
	registerGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Unix())

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout after %d attempts", retryTimes)
		}
		time.Sleep(10 * time.Second)

		resp, err := getExperimentResult(t, ecClient, experiment.Id)
		if err != nil {
			st, _ := status.FromError(err)
			if st.Code() != codes.NotFound {
				t.Fatalf("Failed to get the experiment result. Error code: %d. Error: %v\n", st.Code(), err)
			} else {
				t.Logf("Retry %d/%d: ExperimentResult not found yet (NotFound error)", i+1, retryTimes)
				continue
			}
		}

		if resp != nil {
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
			if gr.GoalId != goalIDs[0] {
				t.Fatalf("goal ID is not correct: %s", gr.GoalId)
			}
			if len(gr.VariationResults) != 2 {
				t.Fatalf("the number of variation results is not correct: %d", len(gr.VariationResults))
			}
			vsA := gr.VariationResults[0]
			vsB := gr.VariationResults[1]
			// These counts are based on the number of events sent earlier
			if vsA.EvaluationCount.EventCount != 4 || // variation A
				vsA.EvaluationCount.UserCount != 3 ||
				vsB.EvaluationCount.EventCount != 3 || // variation B
				vsB.EvaluationCount.UserCount != 2 {
				continue
			}
			// These counts are based on the number of events sent earlier
			if vsA.ExperimentCount.EventCount != 5 || // variation A
				vsA.ExperimentCount.UserCount != 3 ||
				vsB.ExperimentCount.EventCount != 4 || // variation B
				vsB.ExperimentCount.UserCount != 2 {
				continue
			}
			checkExperimentVariationResultA(t, vsA, experiment.Variations[0].Value)
			checkExperimentVariationResultB(t, vsB, experiment.Variations[1].Value)
			break
		}
	}
	res, err := getExperiment(t, experimentClient, experiment.Id)
	if err != nil {
		t.Fatalf("Failed to get experiment. ID: %s. Error: %v", experiment.Id, err)
	}
	if res.Experiment.Status != experimentproto.Experiment_RUNNING {
		expected, _ := experimentproto.Experiment_Status_name[int32(experimentproto.Experiment_RUNNING)]
		actual, _ := experimentproto.Experiment_Status_name[int32(res.Experiment.Status)]
		t.Fatalf("the status of experiment is not correct. expected: %s, but got %s", expected, actual)
	}
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

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
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

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
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

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
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

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
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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

	// Wait for the event-persister-dwh subscribe to the pubsub
	// The batch runs every minute, so we give a extra 10 seconds
	// to ensure that it will subscribe correctly.
	time.Sleep(70 * time.Second)

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
	stopExperiment(ctx, t, experimentClient, experiment.Id)
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
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
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
		cmd := &experimentproto.CreateGoalCommand{
			Id:          createGoalID(t, uuid),
			Name:        createGoalID(t, uuid),
			Description: fmt.Sprintf("%s-goal-description", prefixTestName),
		}
		_, err := client.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
			Command:       cmd,
			EnvironmentId: *environmentID,
		})
		if err != nil {
			t.Fatal(err)
		}
		goalIDs = append(goalIDs, cmd.Id)
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
	cmd := &experimentproto.CreateExperimentCommand{
		Name:            fmt.Sprintf("%s - %v", name, strings.Join(goalIDs, ",")),
		FeatureId:       featureID,
		GoalIds:         goalIDs,
		StartAt:         startAt.Unix(),
		StopAt:          stopAt.Unix(),
		BaseVariationId: baseVariationID,
	}
	resp, err := client.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.StartExperiment(ctx, &experimentproto.StartExperimentRequest{
		EnvironmentId: *environmentID,
		Id:            resp.Experiment.Id,
		Command:       &experimentproto.StartExperimentCommand{},
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	goal, err := ptypes.MarshalAny(&eventproto.GoalEvent{
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if reason == nil {
		reason = &featureproto.Reason{}
	}
	evaluation, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
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
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
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
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
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
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
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

func newCreateFeatureCommand(featureID string, variations []string) *featureproto.CreateFeatureCommand {
	cmd := &featureproto.CreateFeatureCommand{
		Id:          featureID,
		Name:        featureID,
		Description: "e2e-test-eventcounter-feature-description",
		Variations:  []*featureproto.Variation{},
		Tags: []string{
			"e2e-test-tag-1",
			"e2e-test-tag-2",
			"e2e-test-tag-3",
		},
		DefaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
		DefaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
	}
	for _, v := range variations {
		cmd.Variations = append(cmd.Variations, &featureproto.Variation{
			Value:       v,
			Name:        "Variation " + v,
			Description: "Thing does " + v,
		})
	}
	return cmd
}

func createFeature(
	t *testing.T,
	client featureclient.Client,
	featureID, tag, variationA, variationB string,
) {
	t.Helper()
	cmd := newCreateFeatureCommand(featureID, []string{variationA, variationB})
	createReq := &featureproto.CreateFeatureRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
	addTag(t, tag, featureID, client)
	enableFeature(t, featureID, client)
}

func addTag(t *testing.T, tag string, featureID string, client featureclient.Client) {
	t.Helper()
	addReq := &featureproto.UpdateFeatureDetailsRequest{
		Id: featureID,
		AddTagCommands: []*featureproto.AddTagCommand{
			{Tag: tag},
		},
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureDetails(ctx, addReq); err != nil {
		t.Fatal(err)
	}
}

func enableFeature(t *testing.T, featureID string, client featureclient.Client) {
	t.Helper()
	enableReq := &featureproto.EnableFeatureRequest{
		Id:            featureID,
		Command:       &featureproto.EnableFeatureCommand{},
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.EnableFeature(ctx, enableReq); err != nil {
		t.Fatalf("Failed to enable feature id: %s. Error: %v", featureID, err)
	}
}

func addFeatureIndividualTargeting(t *testing.T, featureID, userID, variationID string, client featureclient.Client) {
	t.Helper()
	c := &featureproto.AddUserToVariationCommand{
		Id:   variationID,
		User: userID,
	}
	cmd, err := ptypes.MarshalAny(c)
	assert.NoError(t, err)
	req := &featureproto.UpdateFeatureTargetingRequest{
		Id: featureID,
		Commands: []*featureproto.Command{
			{
				Command: cmd,
			},
		},
		EnvironmentId: *environmentID,
		From:          featureproto.UpdateFeatureTargetingRequest_USER,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureTargeting(ctx, req); err != nil {
		t.Fatalf("Failed add user to individual targeting: %s. Error: %v", featureID, err)
	}
}

func getEvaluation(t *testing.T, tag string, userID string) (*gatewayproto.GetEvaluationsResponse, error) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

// This check the result for variation A
func checkExperimentVariationResultA(t *testing.T, vsA *ecproto.VariationResult, variationValue string) {
	t.Helper()
	// Evaluation
	if vsA.EvaluationCount.EventCount != 4 {
		t.Fatalf("variation: %s: evaluation event count is not correct: %d", variationValue, vsA.EvaluationCount.EventCount)
	}
	if vsA.EvaluationCount.UserCount != 3 {
		t.Fatalf("variation: %s: evaluation user count is not correct: %d", variationValue, vsA.EvaluationCount.UserCount)
	}
	// Experiment
	if vsA.ExperimentCount.EventCount != 5 {
		t.Fatalf("variation: %s: experiment event count is not correct: %d", variationValue, vsA.ExperimentCount.EventCount)
	}
	if vsA.ExperimentCount.UserCount != 3 {
		t.Fatalf("variation: %s: experiment user count is not correct: %d", variationValue, vsA.ExperimentCount.UserCount)
	}
	if diff := cmp.Diff(vsA.ExperimentCount.ValueSum, 1.0, compareFloatOpt); diff != "" {
		t.Fatalf("variation: %s: experiment value sum is not correct: %f", variationValue, vsA.ExperimentCount.ValueSum)
	}
	// cvr prob best
	if diff := cmp.Diff(vsA.CvrProbBest.Mean, 0.57, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", variationValue, vsA.CvrProbBest.Mean)
	}
	if diff := cmp.Diff(vsA.CvrProbBest.Sd, 0.49, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", variationValue, vsA.CvrProbBest.Sd)
	}
	if diff := cmp.Diff(vsA.CvrProbBest.Rhat, 0.99, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", variationValue, vsA.CvrProbBest.Rhat)
	}
	// cvr prob beat baseline
	if diff := cmp.Diff(vsA.CvrProbBeatBaseline.Mean, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", variationValue, vsA.CvrProbBeatBaseline.Mean)
	}
	if diff := cmp.Diff(vsA.CvrProbBeatBaseline.Sd, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", variationValue, vsA.CvrProbBeatBaseline.Sd)
	}
	if diff := cmp.Diff(vsA.CvrProbBeatBaseline.Rhat, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", variationValue, vsA.CvrProbBeatBaseline.Rhat)
	}
	// value sum per user prob best
	if diff := cmp.Diff(vsA.GoalValueSumPerUserProbBest.Mean, 0.32, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", variationValue, vsA.GoalValueSumPerUserProbBest.Mean)
	}
	// value sum per user prob beat baseline
	if diff := cmp.Diff(vsA.GoalValueSumPerUserProbBeatBaseline.Mean, 0.0, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", variationValue, vsA.GoalValueSumPerUserProbBeatBaseline.Mean)
	}
}

// This check the result for variation B
func checkExperimentVariationResultB(t *testing.T, vsB *ecproto.VariationResult, variationValue string) {
	t.Helper()
	// Evaluation
	if vsB.EvaluationCount.EventCount != 3 {
		t.Fatalf("variation: %s: evaluation event count is not correct: %d", variationValue, vsB.EvaluationCount.EventCount)
	}
	if vsB.EvaluationCount.UserCount != 2 {
		t.Fatalf("variation: %s: evaluation user count is not correct: %d", variationValue, vsB.EvaluationCount.UserCount)
	}
	// Experiment
	if vsB.ExperimentCount.EventCount != 4 {
		t.Fatalf("variation: %s: experiment event count is not correct: %d", variationValue, vsB.ExperimentCount.EventCount)
	}
	if vsB.ExperimentCount.UserCount != 2 {
		t.Fatalf("variation: %s: experiment user count is not correct: %d", variationValue, vsB.ExperimentCount.UserCount)
	}
	if diff := cmp.Diff(vsB.ExperimentCount.ValueSum, 0.50, compareFloatOpt); diff != "" {
		t.Fatalf("variation: %s: experiment value sum is not correct: %f", variationValue, vsB.ExperimentCount.ValueSum)
	}
	// cvr prob best
	if diff := cmp.Diff(vsB.CvrProbBest.Mean, 0.42, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", variationValue, vsB.CvrProbBest.Mean)
	}
	if diff := cmp.Diff(vsB.CvrProbBest.Sd, 0.49, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", variationValue, vsB.CvrProbBest.Sd)
	}
	if diff := cmp.Diff(vsB.CvrProbBest.Rhat, 0.99, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", variationValue, vsB.CvrProbBest.Rhat)
	}
	// cvr prob beat baseline
	if diff := cmp.Diff(vsB.CvrProbBeatBaseline.Mean, 0.42, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", variationValue, vsB.CvrProbBeatBaseline.Mean)
	}
	if diff := cmp.Diff(vsB.CvrProbBeatBaseline.Sd, 0.49, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", variationValue, vsB.CvrProbBeatBaseline.Sd)
	}
	if diff := cmp.Diff(vsB.CvrProbBeatBaseline.Rhat, 0.99, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", variationValue, vsB.CvrProbBeatBaseline.Rhat)
	}
	// value sum per user prob best
	if diff := cmp.Diff(vsB.GoalValueSumPerUserProbBest.Mean, 0.67, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", variationValue, vsB.GoalValueSumPerUserProbBest.Mean)
	}
	// value sum per user prob beat baseline
	if diff := cmp.Diff(vsB.GoalValueSumPerUserProbBeatBaseline.Mean, 0.67, compareFloatBayesian); diff != "" {
		t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", variationValue, vsB.GoalValueSumPerUserProbBeatBaseline.Mean)
	}
}
