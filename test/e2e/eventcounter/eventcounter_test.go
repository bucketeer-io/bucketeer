// Copyright 2022 The Bucketeer Authors.
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
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"

	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	gwapi "github.com/bucketeer-io/bucketeer/pkg/gateway/api"
	gatewayclient "github.com/bucketeer-io/bucketeer/pkg/gateway/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/e2e/util"
)

const (
	prefixTestName = "e2e-test"
	timeout        = 20 * time.Second
	retryTimes     = 30
)

const defaultVariationID = "default"

var (
	webGatewayAddr       = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort       = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert       = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath           = flag.String("api-key", "", "Api key path for web gateway")
	gatewayAddr          = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort          = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert          = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath     = flag.String("service-token", "", "Service token path")
	environmentNamespace = flag.String("environment-namespace", "", "Environment namespace")
	testID               = flag.String("test-id", "", "test ID")
	compareFloatOpt      = cmpopts.EquateApprox(0, 0.0001)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	cmd := newCreateFeatureCommand(featureID, []string{variationVarA, variationVarB})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local()
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcExperimentGoalCount", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.2), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Add(-time.Hour).Unix())

	grpcRegisterEvaluationEvent(t, featureID, f.Version, userID, f.Variations[0].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	cmd := newCreateFeatureCommand(featureID, []string{variationVarA, variationVarB})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local()
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentGoalCount", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.2), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.3), time.Now().Add(-time.Hour).Unix())

	registerEvaluationEvent(t, featureID, f.Version, userID, f.Variations[0].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	cmd := newCreateFeatureCommand(featureID, []string{"a", "b"})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

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

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local()
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	// CVRs is 3/4
	// Register goal variation
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[1], tag, float64(0.2), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Add(-time.Hour).Unix())
	// Increment experiment event count
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	// Register 3 events and 2 user counts for user 1, 2 and 3
	// Register variation a
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[1], experiment.Variations[0].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[2], experiment.Variations[0].Id, tag, reason)
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, tag, reason)

	// CVRs is 2/3
	// Register goal
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Add(-time.Hour).Unix())
	// Increment experiment event count
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	// Register 3 events and 2 user counts for user 4 and 5
	// Register variation
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[4], experiment.Variations[1].Id, tag, reason)
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)
		resp := getExperimentResult(t, ecClient, experiment.Id)
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
			if gr.VariationResults[0].EvaluationCount.EventCount == 0 || // variation A
				gr.VariationResults[0].EvaluationCount.UserCount == 0 ||
				gr.VariationResults[1].EvaluationCount.EventCount == 0 || // variation B
				gr.VariationResults[1].EvaluationCount.UserCount == 0 {
				continue
			}
			if gr.VariationResults[0].ExperimentCount.EventCount == 0 || // variation A
				gr.VariationResults[0].ExperimentCount.UserCount == 0 ||
				gr.VariationResults[1].ExperimentCount.EventCount == 0 || // variation B
				gr.VariationResults[1].ExperimentCount.UserCount == 0 {
				continue
			}
			for _, vr := range gr.VariationResults {
				// variation a
				if vr.VariationId == experiment.Variations[0].Id {
					vv := experiment.Variations[0].Value
					// Evaluation
					if vr.EvaluationCount.EventCount != 4 {
						t.Fatalf("variation: %s: evaluation event count is not correct: %d", vv, vr.EvaluationCount.EventCount)
					}
					if vr.EvaluationCount.UserCount != 3 {
						t.Fatalf("variation: %s: evaluation user count is not correct: %d", vv, vr.EvaluationCount.UserCount)
					}
					// Experiment
					if vr.ExperimentCount.EventCount != 4 {
						t.Fatalf("variation: %s: experiment event count is not correct: %d", vv, vr.ExperimentCount.EventCount)
					}
					if vr.ExperimentCount.UserCount != 3 {
						t.Fatalf("variation: %s: experiment user count is not correct: %d", vv, vr.ExperimentCount.UserCount)
					}
					if diff := cmp.Diff(vr.ExperimentCount.ValueSum, 0.9, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: experiment value sum is not correct: %f", vv, vr.ExperimentCount.ValueSum)
					}
					// cvr prob best
					if vr.CvrProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", vv, vr.CvrProbBest.Mean)
					}
					if vr.CvrProbBest.Sd <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", vv, vr.CvrProbBest.Sd)
					}
					if vr.CvrProbBest.Rhat <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", vv, vr.CvrProbBest.Rhat)
					}
					// cvr prob beat baseline
					if diff := cmp.Diff(vr.CvrProbBeatBaseline.Mean, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", vv, vr.CvrProbBeatBaseline.Mean)
					}
					if diff := cmp.Diff(vr.CvrProbBeatBaseline.Sd, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", vv, vr.CvrProbBeatBaseline.Sd)
					}
					if diff := cmp.Diff(vr.CvrProbBeatBaseline.Rhat, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", vv, vr.CvrProbBeatBaseline.Rhat)
					}
					// value sum per user prob best
					if vr.GoalValueSumPerUserProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBest.Mean)
					}
					// value sum per user prob beat baseline
					if diff := cmp.Diff(vr.GoalValueSumPerUserProbBeatBaseline.Mean, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBeatBaseline.Mean)
					}
					continue
				}
				// variation b
				if vr.VariationId == experiment.Variations[1].Id {
					vv := experiment.Variations[1].Value
					// Evaluation
					if vr.EvaluationCount.EventCount != 3 {
						t.Fatalf("variation: %s: evaluation event count is not correct: %d", vv, vr.EvaluationCount.EventCount)
					}
					if vr.EvaluationCount.UserCount != 2 {
						t.Fatalf("variation: %s: evaluation user count is not correct: %d", vv, vr.EvaluationCount.UserCount)
					}
					// Experiment
					if vr.ExperimentCount.EventCount != 3 {
						t.Fatalf("variation: %s: experiment event count is not correct: %d", vv, vr.ExperimentCount.EventCount)
					}
					if vr.ExperimentCount.UserCount != 2 {
						t.Fatalf("variation: %s: experiment user count is not correct: %d", vv, vr.ExperimentCount.UserCount)
					}
					if diff := cmp.Diff(vr.ExperimentCount.ValueSum, 0.35, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: experiment value sum is not correct: %f", vv, vr.ExperimentCount.ValueSum)
					}
					// cvr prob best
					if vr.CvrProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", vv, vr.CvrProbBest.Mean)
					}
					if vr.CvrProbBest.Sd <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", vv, vr.CvrProbBest.Sd)
					}
					if vr.CvrProbBest.Rhat <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", vv, vr.CvrProbBest.Rhat)
					}
					// cvr prob beat baseline
					if vr.CvrProbBeatBaseline.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", vv, vr.CvrProbBeatBaseline.Mean)
					}
					if vr.CvrProbBeatBaseline.Sd <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", vv, vr.CvrProbBeatBaseline.Sd)
					}
					if vr.CvrProbBeatBaseline.Rhat <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", vv, vr.CvrProbBeatBaseline.Rhat)
					}
					// value sum per user prob best
					if vr.GoalValueSumPerUserProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBest.Mean)
					}
					// value sum per user prob beat baseline
					if vr.GoalValueSumPerUserProbBeatBaseline.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBeatBaseline.Mean)
					}
					continue
				}
				t.Fatalf("unknown variation results: %s", vr.VariationId)
			}
			break
		}
	}
	res := getExperiment(t, experimentClient, experiment.Id)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	cmd := newCreateFeatureCommand(featureID, []string{"a", "b"})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

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

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local()
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	// CVRs is 3/4
	// Register goal variation
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[1], tag, float64(0.2), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1), time.Now().Add(-time.Hour).Unix())
	// Increment experiment event count
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	// Register 3 events and 2 user counts for user 1, 2 and 3
	// Register variation a
	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], f.Variations[0].Id, tag, reason)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[1], f.Variations[0].Id, tag, reason)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[2], f.Variations[0].Id, tag, reason)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], f.Variations[0].Id, tag, reason)

	// CVRs is 2/3
	// Register goal
	registerGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15), time.Now().Add(-time.Hour).Unix())
	// Increment experiment event count
	registerGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1), time.Now().Unix())
	// Register 3 events and 2 user counts for user 4 and 5
	// Register variation
	registerEvaluationEvent(t, featureID, f.Version, userIDs[3], f.Variations[1].Id, tag, reason)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[4], f.Variations[1].Id, tag, reason)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, f.Version, userIDs[3], f.Variations[1].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)
		resp := getExperimentResult(t, ecClient, experiment.Id)
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
			if gr.VariationResults[0].EvaluationCount.EventCount == 0 || // variation A
				gr.VariationResults[0].EvaluationCount.UserCount == 0 ||
				gr.VariationResults[1].EvaluationCount.EventCount == 0 || // variation B
				gr.VariationResults[1].EvaluationCount.UserCount == 0 {
				continue
			}
			if gr.VariationResults[0].ExperimentCount.EventCount == 0 || // variation A
				gr.VariationResults[0].ExperimentCount.UserCount == 0 ||
				gr.VariationResults[1].ExperimentCount.EventCount == 0 || // variation B
				gr.VariationResults[1].ExperimentCount.UserCount == 0 {
				continue
			}
			for _, vr := range gr.VariationResults {
				// variation a
				if vr.VariationId == experiment.Variations[0].Id {
					vv := experiment.Variations[0].Value
					// Evaluation
					if vr.EvaluationCount.EventCount != 4 {
						t.Fatalf("variation: %s: evaluation event count is not correct: %d", vv, vr.EvaluationCount.EventCount)
					}
					if vr.EvaluationCount.UserCount != 3 {
						t.Fatalf("variation: %s: evaluation user count is not correct: %d", vv, vr.EvaluationCount.UserCount)
					}
					// Experiment
					if vr.ExperimentCount.EventCount != 4 {
						t.Fatalf("variation: %s: experiment event count is not correct: %d", vv, vr.ExperimentCount.EventCount)
					}
					if vr.ExperimentCount.UserCount != 3 {
						t.Fatalf("variation: %s: experiment user count is not correct: %d", vv, vr.ExperimentCount.UserCount)
					}
					if diff := cmp.Diff(vr.ExperimentCount.ValueSum, 0.9, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: experiment value sum is not correct: %f", vv, vr.ExperimentCount.ValueSum)
					}
					// cvr prob best
					if vr.CvrProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", vv, vr.CvrProbBest.Mean)
					}
					if vr.CvrProbBest.Sd <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", vv, vr.CvrProbBest.Sd)
					}
					if vr.CvrProbBest.Rhat <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", vv, vr.CvrProbBest.Rhat)
					}
					// cvr prob beat baseline
					if diff := cmp.Diff(vr.CvrProbBeatBaseline.Mean, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", vv, vr.CvrProbBeatBaseline.Mean)
					}
					if diff := cmp.Diff(vr.CvrProbBeatBaseline.Sd, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", vv, vr.CvrProbBeatBaseline.Sd)
					}
					if diff := cmp.Diff(vr.CvrProbBeatBaseline.Rhat, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", vv, vr.CvrProbBeatBaseline.Rhat)
					}
					// value sum per user prob best
					if vr.GoalValueSumPerUserProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBest.Mean)
					}
					// value sum per user prob beat baseline
					if diff := cmp.Diff(vr.GoalValueSumPerUserProbBeatBaseline.Mean, 0.0, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBeatBaseline.Mean)
					}
					continue
				}
				// variation b
				if vr.VariationId == experiment.Variations[1].Id {
					vv := experiment.Variations[1].Value
					// Evaluation
					if vr.EvaluationCount.EventCount != 3 {
						t.Fatalf("variation: %s: evaluation event count is not correct: %d", vv, vr.EvaluationCount.EventCount)
					}
					if vr.EvaluationCount.UserCount != 2 {
						t.Fatalf("variation: %s: evaluation user count is not correct: %d", vv, vr.EvaluationCount.UserCount)
					}
					// Experiment
					if vr.ExperimentCount.EventCount != 3 {
						t.Fatalf("variation: %s: experiment event count is not correct: %d", vv, vr.ExperimentCount.EventCount)
					}
					if vr.ExperimentCount.UserCount != 2 {
						t.Fatalf("variation: %s: experiment user count is not correct: %d", vv, vr.ExperimentCount.UserCount)
					}
					if diff := cmp.Diff(vr.ExperimentCount.ValueSum, 0.35, compareFloatOpt); diff != "" {
						t.Fatalf("variation: %s: experiment value sum is not correct: %f", vv, vr.ExperimentCount.ValueSum)
					}
					// cvr prob best
					if vr.CvrProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best mean is not correct: %f", vv, vr.CvrProbBest.Mean)
					}
					if vr.CvrProbBest.Sd <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best sd is not correct: %f", vv, vr.CvrProbBest.Sd)
					}
					if vr.CvrProbBest.Rhat <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob best rhat is not correct: %f", vv, vr.CvrProbBest.Rhat)
					}
					// cvr prob beat baseline
					if vr.CvrProbBeatBaseline.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", vv, vr.CvrProbBeatBaseline.Mean)
					}
					if vr.CvrProbBeatBaseline.Sd <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", vv, vr.CvrProbBeatBaseline.Sd)
					}
					if vr.CvrProbBeatBaseline.Rhat <= float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", vv, vr.CvrProbBeatBaseline.Rhat)
					}
					// value sum per user prob best
					if vr.GoalValueSumPerUserProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBest.Mean)
					}
					// value sum per user prob beat baseline
					if vr.GoalValueSumPerUserProbBeatBaseline.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob beat baseline mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBeatBaseline.Mean)
					}
					continue
				}
				t.Fatalf("unknown variation results: %s", vr.VariationId)
			}
			break
		}
	}
	res := getExperiment(t, experimentClient, experiment.Id)
	if res.Experiment.Status != experimentproto.Experiment_RUNNING {
		expected, _ := experimentproto.Experiment_Status_name[int32(experimentproto.Experiment_RUNNING)]
		actual, _ := experimentproto.Experiment_Status_name[int32(res.Experiment.Status)]
		t.Fatalf("the status of experiment is not correct. expected: %s, but got %s", expected, actual)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	cmd := newCreateFeatureCommand(featureID, []string{variationVarA, variationVarB})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userIDs[0], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[1], f.Variations[1].Id, featureClient)

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	goalIDs := createGoals(ctx, t, experimentClient, 3)
	startAt := time.Now().Local()
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcMultiGoalsEventCounter", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Add(-time.Hour).Unix())

	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], f.Variations[0].Id, tag, reason)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[1], f.Variations[1].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		// Goal 0.
		resp := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
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
		resp = getExperimentGoalCount(t, ecClient, goalIDs[1], featureID, f.Version, variationIDs)
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
		resp = getExperimentGoalCount(t, ecClient, goalIDs[2], featureID, f.Version, variationIDs)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := []string{}
	for i := 0; i < 5; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	featureID := createFeatureID(t, uuid)

	variationVarA := "a"
	variationVarB := "b"
	cmd := newCreateFeatureCommand(featureID, []string{variationVarA, variationVarB})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userIDs[0], f.Variations[0].Id, featureClient)
	addFeatureIndividualTargeting(t, featureID, userIDs[1], f.Variations[1].Id, featureClient)

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	goalIDs := createGoals(ctx, t, experimentClient, 3)
	startAt := time.Now().Local()
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestMultiGoalsEventCounter", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3), time.Now().Unix())
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Unix())
	// This event will be ignored because the timestamp is older than the experiment startAt time stamp
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2), time.Now().Add(-time.Hour).Unix())

	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], f.Variations[0].Id, tag, reason)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[1], f.Variations[1].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		// Goal 0.
		resp := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
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
		resp = getExperimentGoalCount(t, ecClient, goalIDs[1], featureID, f.Version, variationIDs)
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
		resp = getExperimentGoalCount(t, ecClient, goalIDs[2], featureID, f.Version, variationIDs)
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := createUserID(t, uuid)
	featureID := createFeatureID(t, uuid)
	value := float64(1.23)

	variationVarA := "a"
	variationVarB := "b"
	cmd := newCreateFeatureCommand(featureID, []string{variationVarA, variationVarB})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestHTTPTrack", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	// Send track events.
	sendHTTPTrack(t, userID, goalIDs[0], tag, value)
	registerEvaluationEvent(t, featureID, f.Version, userID, f.Variations[0].Id, tag, reason)

	// Check the count
	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp := getExperimentGoalCount(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	cmd := newCreateFeatureCommand(featureID, []string{variationVarA, variationVarB})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range f.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	createExperimentWithMultiGoals(
		ctx,
		t,
		experimentClient,
		"TestGrpcExperimentEvaluationEventCount",
		featureID,
		goalIDs,
		f.Variations[0].Id,
		startAt, stopAt,
	)

	grpcRegisterEvaluationEvent(t, featureID, f.Version, userID, variations[variationVarA].Id, tag, reason)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp := getExperimentEvaluationCount(t, ecClient, featureID, f.Version, variationIDs)
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
		if vcA.EventCount != 1 {
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	cmd := newCreateFeatureCommand(featureID, []string{variationVarA, variationVarB})
	createFeature(t, featureClient, cmd)
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)

	// Because we set the user to the individual targeting,
	// We must ensure to set the correct reason. Otherwise, it will fail when the event persister
	// evaluates the user
	reason := &featureproto.Reason{
		Type: featureproto.Reason_TARGET,
	}
	addFeatureIndividualTargeting(t, featureID, userID, f.Variations[0].Id, featureClient)

	// Get the latest version after all changes in the feature flag
	f = getFeature(t, featureClient, featureID)

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range f.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}
	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	createExperimentWithMultiGoals(
		ctx,
		t,
		experimentClient,
		"TestExperimentEvaluationEventCount",
		featureID,
		goalIDs,
		f.Variations[0].Id,
		startAt, stopAt,
	)

	registerEvaluationEvent(t, featureID, f.Version, userID, variations[variationVarA].Id, tag, reason)
	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(10 * time.Second)

		resp := getExperimentEvaluationCount(t, ecClient, featureID, f.Version, variationIDs)
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
		if vcA.EventCount != 1 {
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
	cmd := newCreateFeatureCommand(featureID, []string{"a", "b"})
	createFeature(t, featureClient, cmd)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userIDs := []string{}
	for i := 0; i < 8; i++ {
		userIDs = append(userIDs, fmt.Sprintf("%s-%d", createUserID(t, uuid), i))
	}
	addTag(t, tag, featureID, featureClient)
	enableFeature(t, featureID, featureClient)
	f := getFeature(t, featureClient, featureID)
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
		res := getEvaluationTimeseriesCount(t, featureID, ecClient, ecproto.GetEvaluationTimeseriesCountRequest_THIRTY_DAYS)
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
			Command:              cmd,
			EnvironmentNamespace: *environmentNamespace,
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
		Name:            name + strings.Join(goalIDs, ","),
		FeatureId:       featureID,
		GoalIds:         goalIDs,
		StartAt:         startAt.Unix(),
		StopAt:          stopAt.Unix(),
		BaseVariationId: baseVariationID,
	}
	resp, err := client.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{
		Command:              cmd,
		EnvironmentNamespace: *environmentNamespace,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.StartExperiment(ctx, &experimentproto.StartExperimentRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   resp.Experiment.Id,
		Command:              &experimentproto.StartExperimentCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.Experiment
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

// Test for old SDK client
// Evaluation field in the GoalEvent is deprecated.
func registerGoalEventWithEvaluations(
	t *testing.T,
	featureID string,
	featureVersion int32,
	goalID, userID, variationID string,
	value float64,
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	goal, err := ptypes.MarshalAny(&eventproto.GoalEvent{
		Timestamp: time.Now().Unix(),
		GoalId:    goalID,
		UserId:    userID,
		Value:     value,
		User:      &userproto.User{},
		Evaluations: []*featureproto.Evaluation{
			{
				Id:             fmt.Sprintf("%s-evaluation-id-%s", prefixTestName, newUUID(t)),
				FeatureId:      featureID,
				FeatureVersion: featureVersion,
				UserId:         userID,
				VariationId:    variationID,
			},
		},
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

func sendHTTPTrack(t *testing.T, userID, goalID, tag string, value float64) {
	data, err := ioutil.ReadFile(*apiKeyPath)
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
	client := &http.Client{}
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

func createFeature(t *testing.T, client featureclient.Client, cmd *featureproto.CreateFeatureCommand) {
	t.Helper()
	createReq := &featureproto.CreateFeatureRequest{
		Command:              cmd,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func addTag(t *testing.T, tag string, featureID string, client featureclient.Client) {
	t.Helper()
	addReq := &featureproto.UpdateFeatureDetailsRequest{
		Id: featureID,
		AddTagCommands: []*featureproto.AddTagCommand{
			{Tag: tag},
		},
		EnvironmentNamespace: *environmentNamespace,
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
		Id:                   featureID,
		Command:              &featureproto.EnableFeatureCommand{},
		EnvironmentNamespace: *environmentNamespace,
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
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureTargeting(ctx, req); err != nil {
		t.Fatalf("Failed add user to individual targeting: %s. Error: %v", featureID, err)
	}
}

func getEvaluation(t *testing.T, tag string, userID string) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &gatewayproto.GetEvaluationsRequest{
		Tag:  tag,
		User: &userproto.User{Id: userID},
	}
	response, err := c.GetEvaluations(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func getExperiment(t *testing.T, c experimentclient.Client, id string) *experimentproto.GetExperimentResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &experimentproto.GetExperimentRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	}
	res, err := c.GetExperiment(ctx, req)
	if err != nil {
		// pass not found error
		if err.Error() != "rpc error: code = NotFound desc = eventcounter: not found" {
			t.Fatal(err)
		}
	}
	return res
}

func getExperimentResult(t *testing.T, c ecclient.Client, experimentID string) *ecproto.GetExperimentResultResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &ecproto.GetExperimentResultRequest{
		ExperimentId:         experimentID,
		EnvironmentNamespace: *environmentNamespace,
	}
	response, err := c.GetExperimentResult(ctx, req)
	if err != nil {
		// pass not found error
		if err.Error() != "rpc error: code = NotFound desc = eventcounter: not found" {
			t.Fatal(err)
		}
	}
	return response
}

func getExperimentEvaluationCount(t *testing.T, c ecclient.Client, featureID string, featureVersion int32, variationIDs []string) *ecproto.GetExperimentEvaluationCountResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	now := time.Now()
	req := &ecproto.GetExperimentEvaluationCountRequest{
		EnvironmentNamespace: *environmentNamespace,
		StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
		EndAt:                now.Unix(),
		FeatureId:            featureID,
		FeatureVersion:       featureVersion,
		VariationIds:         variationIDs,
	}
	response, err := c.GetExperimentEvaluationCount(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func getExperimentGoalCount(t *testing.T, c ecclient.Client, goalID, featureID string, featureVersion int32, variationIDs []string) *ecproto.GetExperimentGoalCountResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	now := time.Now()
	req := &ecproto.GetExperimentGoalCountRequest{
		EnvironmentNamespace: *environmentNamespace,
		StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
		EndAt:                now.Unix(),
		GoalId:               goalID,
		FeatureId:            featureID,
		FeatureVersion:       featureVersion,
		VariationIds:         variationIDs,
	}
	response, err := c.GetExperimentGoalCount(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func getFeature(t *testing.T, client featureclient.Client, featureID string) *featureproto.Feature {
	t.Helper()
	getReq := &featureproto.GetFeatureRequest{
		Id:                   featureID,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.GetFeature(ctx, getReq)
	if err != nil {
		t.Fatal("Failed to get feature:", err)
	}
	return response.Feature
}

func createFeatureID(t *testing.T, uuid string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, uuid)
}

func createGoalID(t *testing.T, uuid string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-goal-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-goal-id-%s", prefixTestName, uuid)
}

func createUserID(t *testing.T, uuid string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-user-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-user-id-%s", prefixTestName, uuid)
}

func getEvaluationTimeseriesCount(
	t *testing.T,
	featureID string,
	c ecclient.Client,
	timeRange ecproto.GetEvaluationTimeseriesCountRequest_TimeRange,
) *ecproto.GetEvaluationTimeseriesCountResponse {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &ecproto.GetEvaluationTimeseriesCountRequest{
		EnvironmentNamespace: *environmentNamespace,
		FeatureId:            featureID,
		TimeRange:            timeRange,
	}
	res, err := c.GetEvaluationTimeseriesCount(
		ctx,
		req,
	)
	if err != nil {
		t.Fatal(err)
	}
	return res
}
