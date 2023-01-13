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
	"google.golang.org/protobuf/encoding/protojson"

	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
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
	retryTimes     = 360
)

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
)

func TestGrpcGoalCountV2(t *testing.T) {
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
	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcGoalCountV2", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.2))
	grpcRegisterGoalEvent(t, goalIDs[0], userID, tag, float64(0.3))

	grpcRegisterEvaluationEvent(t, featureID, f.Version, userID, f.Variations[0].Id, tag)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)

		resp := getGoalCountV2(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if vcA.ValueSum != float64(0.5) {
			continue
		}
		if vcA.ValueSumPerUserMean != float64(0.5) {
			continue
		}
		if vcA.ValueSumPerUserVariance != float64(0) {
			continue
		}

		vcB := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
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

func TestGoalCountV2(t *testing.T) {
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
	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGoalCountV2", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.2))
	registerGoalEvent(t, goalIDs[0], userID, tag, float64(0.3))

	registerEvaluationEvent(t, featureID, f.Version, userID, f.Variations[0].Id, tag)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)

		resp := getGoalCountV2(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if vcA.ValueSum != float64(0.5) {
			continue
		}
		if vcA.ValueSumPerUserMean != float64(0.5) {
			continue
		}
		if vcA.ValueSumPerUserVariance != float64(0) {
			continue
		}

		vcB := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
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

// Test for old SDK client. Tag is not set in the EvaluationEvent and GoalEvent
// Evaluation field in the GoalEvent is deprecated.
func TestExperimentResultWithoutTag(t *testing.T) {
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
	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentResultWithoutTag", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	// CVRs is 3/4
	// Register goal variation
	registerGoalEventWithEvaluations(t, featureID, f.Version, goalIDs[0], userIDs[0], experiment.Variations[0].Id, float64(0.3))
	registerGoalEventWithEvaluations(t, featureID, f.Version, goalIDs[0], userIDs[1], experiment.Variations[0].Id, float64(0.2))
	registerGoalEventWithEvaluations(t, featureID, f.Version, goalIDs[0], userIDs[2], experiment.Variations[0].Id, float64(0.1))
	// Register 3 events and 2 user counts for user 1, 2 and 3
	// Register variation a
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, "")
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[1], experiment.Variations[0].Id, "")
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[2], experiment.Variations[0].Id, "")
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, "")

	// CVRs is 2/3
	// Register goal
	registerGoalEventWithEvaluations(t, featureID, f.Version, goalIDs[0], userIDs[3], experiment.Variations[1].Id, float64(0.1))
	registerGoalEventWithEvaluations(t, featureID, f.Version, goalIDs[0], userIDs[4], experiment.Variations[1].Id, float64(0.15))
	// Register 3 events and 2 user counts for user 4 and 5
	// Register variation
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, "")
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[4], experiment.Variations[1].Id, "")
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, "")

	for i := 0; i < retryTimes; i++ {
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
			if gr.VariationResults[0].EvaluationCount.EventCount == 0 && // variation A
				gr.VariationResults[0].EvaluationCount.UserCount == 0 &&
				gr.VariationResults[1].EvaluationCount.EventCount == 0 && // variation B
				gr.VariationResults[1].EvaluationCount.UserCount == 0 {
				continue
			}
			if gr.VariationResults[0].ExperimentCount.EventCount == 0 && // variation A
				gr.VariationResults[0].ExperimentCount.UserCount == 0 &&
				gr.VariationResults[1].ExperimentCount.EventCount == 0 && // variation B
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
					if vr.ExperimentCount.EventCount != 3 {
						t.Fatalf("variation: %s: experiment event count is not correct: %d", vv, vr.ExperimentCount.EventCount)
					}
					if vr.ExperimentCount.UserCount != 3 {
						t.Fatalf("variation: %s: experiment user count is not correct: %d", vv, vr.ExperimentCount.UserCount)
					}
					if vr.ExperimentCount.ValueSum != float64(0.6) {
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
					if vr.CvrProbBeatBaseline.Mean != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", vv, vr.CvrProbBeatBaseline.Mean)
					}
					if vr.CvrProbBeatBaseline.Sd != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", vv, vr.CvrProbBeatBaseline.Sd)
					}
					if vr.CvrProbBeatBaseline.Rhat != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", vv, vr.CvrProbBeatBaseline.Rhat)
					}
					// value sum per user prob best
					if vr.GoalValueSumPerUserProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBest.Mean)
					}
					// value sum per user prob beat baseline
					if vr.GoalValueSumPerUserProbBeatBaseline.Mean != float64(0.0) {
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
					if vr.ExperimentCount.EventCount != 2 {
						t.Fatalf("variation: %s: experiment event count is not correct: %d", vv, vr.ExperimentCount.EventCount)
					}
					if vr.ExperimentCount.UserCount != 2 {
						t.Fatalf("variation: %s: experiment user count is not correct: %d", vv, vr.ExperimentCount.UserCount)
					}
					if vr.ExperimentCount.ValueSum != float64(0.25) {
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
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)
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
	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestGrpcExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	// CVRs is 3/4
	// Register goal variation
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[1], tag, float64(0.2))
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1))
	// Increment experiment event count
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	// Register 3 events and 2 user counts for user 1, 2 and 3
	// Register variation a
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, tag)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[1], experiment.Variations[0].Id, tag)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[2], experiment.Variations[0].Id, tag)
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, tag)

	// CVRs is 2/3
	// Register goal
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1))
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15))
	// Increment experiment event count
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1))
	// Register 3 events and 2 user counts for user 4 and 5
	// Register variation
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, tag)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[4], experiment.Variations[1].Id, tag)
	// Increment evaluation event count
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, tag)

	for i := 0; i < retryTimes; i++ {
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
			if gr.VariationResults[0].EvaluationCount.EventCount == 0 && // variation A
				gr.VariationResults[0].EvaluationCount.UserCount == 0 &&
				gr.VariationResults[1].EvaluationCount.EventCount == 0 && // variation B
				gr.VariationResults[1].EvaluationCount.UserCount == 0 {
				continue
			}
			if gr.VariationResults[0].ExperimentCount.EventCount == 0 && // variation A
				gr.VariationResults[0].ExperimentCount.UserCount == 0 &&
				gr.VariationResults[1].ExperimentCount.EventCount == 0 && // variation B
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
					if vr.ExperimentCount.ValueSum != float64(0.9) {
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
					if vr.CvrProbBeatBaseline.Mean != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", vv, vr.CvrProbBeatBaseline.Mean)
					}
					if vr.CvrProbBeatBaseline.Sd != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", vv, vr.CvrProbBeatBaseline.Sd)
					}
					if vr.CvrProbBeatBaseline.Rhat != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", vv, vr.CvrProbBeatBaseline.Rhat)
					}
					// value sum per user prob best
					if vr.GoalValueSumPerUserProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBest.Mean)
					}
					// value sum per user prob beat baseline
					if vr.GoalValueSumPerUserProbBeatBaseline.Mean != float64(0.0) {
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
					if vr.ExperimentCount.ValueSum != float64(0.35) {
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
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)
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
	goalIDs := createGoals(ctx, t, experimentClient, 1)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "TestExperimentResult", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	// CVRs is 3/4
	// Register goal variation
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	registerGoalEvent(t, goalIDs[0], userIDs[1], tag, float64(0.2))
	registerGoalEvent(t, goalIDs[0], userIDs[2], tag, float64(0.1))
	// Increment experiment event count
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	// Register 3 events and 2 user counts for user 1, 2 and 3
	// Register variation a
	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, tag)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[1], experiment.Variations[0].Id, tag)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[2], experiment.Variations[0].Id, tag)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, tag)

	// CVRs is 2/3
	// Register goal
	registerGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1))
	registerGoalEvent(t, goalIDs[0], userIDs[4], tag, float64(0.15))
	// Increment experiment event count
	registerGoalEvent(t, goalIDs[0], userIDs[3], tag, float64(0.1))
	// Register 3 events and 2 user counts for user 4 and 5
	// Register variation
	registerEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, tag)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[4], experiment.Variations[1].Id, tag)
	// Increment evaluation event count
	registerEvaluationEvent(t, featureID, f.Version, userIDs[3], experiment.Variations[1].Id, tag)

	for i := 0; i < retryTimes; i++ {
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
			if gr.VariationResults[0].EvaluationCount.EventCount == 0 && // variation A
				gr.VariationResults[0].EvaluationCount.UserCount == 0 &&
				gr.VariationResults[1].EvaluationCount.EventCount == 0 && // variation B
				gr.VariationResults[1].EvaluationCount.UserCount == 0 {
				continue
			}
			if gr.VariationResults[0].ExperimentCount.EventCount == 0 && // variation A
				gr.VariationResults[0].ExperimentCount.UserCount == 0 &&
				gr.VariationResults[1].ExperimentCount.EventCount == 0 && // variation B
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
					if vr.ExperimentCount.ValueSum != float64(0.9) {
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
					if vr.CvrProbBeatBaseline.Mean != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline mean is not correct: %f", vv, vr.CvrProbBeatBaseline.Mean)
					}
					if vr.CvrProbBeatBaseline.Sd != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best sd is not correct: %f", vv, vr.CvrProbBeatBaseline.Sd)
					}
					if vr.CvrProbBeatBaseline.Rhat != float64(0.0) {
						t.Fatalf("variation: %s: cvr prob beat baseline best rhat is not correct: %f", vv, vr.CvrProbBeatBaseline.Rhat)
					}
					// value sum per user prob best
					if vr.GoalValueSumPerUserProbBest.Mean <= float64(0.0) {
						t.Fatalf("variation: %s: value sum per user prob best mean is not correct: %f", vv, vr.GoalValueSumPerUserProbBest.Mean)
					}
					// value sum per user prob beat baseline
					if vr.GoalValueSumPerUserProbBeatBaseline.Mean != float64(0.0) {
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
					if vr.ExperimentCount.ValueSum != float64(0.35) {
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
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)
	}
	res := getExperiment(t, experimentClient, experiment.Id)
	if res.Experiment.Status != experimentproto.Experiment_RUNNING {
		expected, _ := experimentproto.Experiment_Status_name[int32(experimentproto.Experiment_RUNNING)]
		actual, _ := experimentproto.Experiment_Status_name[int32(res.Experiment.Status)]
		t.Fatalf("the status of experiment is not correct. expected: %s, but got %s", expected, actual)
	}
}

func TestGrpcMultiGoalsEventCounterRealtime(t *testing.T) {
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
	goalIDs := createGoals(ctx, t, experimentClient, 3)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "GrpcMultiGoalsEventCounterRealtime", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	grpcRegisterGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2))
	grpcRegisterGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2))

	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[0], experiment.Variations[0].Id, tag)
	grpcRegisterEvaluationEvent(t, featureID, f.Version, userIDs[1], experiment.Variations[1].Id, tag)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)

		// Goal 0.
		resp := getGoalCountV2(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if vcA.ValueSum != float64(0.6) {
			continue
		}

		vcB := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
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

		// Goal 1.
		resp = getGoalCountV2(t, ecClient, goalIDs[1], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[1] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if vcA.ValueSum != float64(0) {
			continue
		}

		vcB = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 1 {
			continue
		}
		if vcB.EventCount != 2 {
			continue
		}
		if vcB.ValueSum != float64(0.4) {
			continue
		}

		// Goal 2.
		resp = getGoalCountV2(t, ecClient, goalIDs[2], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[2] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if vcA.ValueSum != float64(0) {
			continue
		}

		vcB = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
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

func TestMultiGoalsEventCounterRealtime(t *testing.T) {
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
	goalIDs := createGoals(ctx, t, experimentClient, 3)
	startAt := time.Now().Local().Add(-1 * time.Hour)
	stopAt := startAt.Local().Add(time.Hour * 2)
	experiment := createExperimentWithMultiGoals(
		ctx, t, experimentClient, "GrpcMultiGoalsEventCounterRealtime", featureID, goalIDs, f.Variations[0].Id, startAt, stopAt)

	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range experiment.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	registerGoalEvent(t, goalIDs[0], userIDs[0], tag, float64(0.3))
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2))
	registerGoalEvent(t, goalIDs[1], userIDs[1], tag, float64(0.2))

	registerEvaluationEvent(t, featureID, f.Version, userIDs[0], f.Variations[0].Id, tag)
	registerEvaluationEvent(t, featureID, f.Version, userIDs[1], f.Variations[1].Id, tag)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)

		// Goal 0.
		resp := getGoalCountV2(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 2 {
			continue
		}
		if vcA.ValueSum != float64(0.6) {
			continue
		}

		vcB := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
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

		// Goal 1.
		resp = getGoalCountV2(t, ecClient, goalIDs[1], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[1] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if vcA.ValueSum != float64(0) {
			continue
		}

		vcB = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation b is missing")
		}
		if vcB.UserCount != 1 {
			continue
		}
		if vcB.EventCount != 2 {
			continue
		}
		if vcB.ValueSum != float64(0.4) {
			continue
		}

		// Goal 2.
		resp = getGoalCountV2(t, ecClient, goalIDs[2], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[2] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		vcA = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 0 {
			continue
		}
		if vcA.EventCount != 0 {
			continue
		}
		if vcA.ValueSum != float64(0) {
			continue
		}

		vcB = getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
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
	registerEvaluationEvent(t, featureID, f.Version, userID, f.Variations[0].Id, tag)

	// Check the count
	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)

		resp := getGoalCountV2(t, ecClient, goalIDs[0], featureID, f.Version, variationIDs)
		if len(resp.GoalCounts.RealtimeCounts) == 0 {
			t.Fatalf("no count returned")
		}
		if resp.GoalCounts.GoalId != goalIDs[0] {
			t.Fatalf("goal ID is not correct: %s", resp.GoalCounts.GoalId)
		}

		// variation a
		vcA := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarA].Id)
		if vcA == nil {
			t.Fatalf("variation a is missing")
		}
		if vcA.UserCount != 1 {
			continue
		}
		if vcA.EventCount != 1 {
			continue
		}
		if vcA.ValueSum != value {
			continue
		}
		// variation b
		vcB := getVariationCount(resp.GoalCounts.RealtimeCounts, variations[variationVarB].Id)
		if vcB == nil {
			t.Fatalf("variation a is missing")
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

func TestGrpcEvaluationEventCountV2(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
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
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range f.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	grpcRegisterEvaluationEvent(t, featureID, f.Version, userID, variations[variationVarA].Id, tag)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)

		resp := getEvaluationCountV2(t, ecClient, featureID, f.Version, variationIDs)
		if resp == nil {
			continue
		}
		ec := resp.Count
		if ec.FeatureId != featureID {
			t.Fatalf("feature ID is not correct: %s", ec.FeatureId)
		}
		if ec.FeatureVersion != f.Version {
			t.Fatalf("feature version is not correct: %d", ec.FeatureVersion)
		}

		vcA := getVariationCount(ec.RealtimeCounts, variations[variationVarA].Id)
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

		vcB := getVariationCount(ec.RealtimeCounts, variations[variationVarB].Id)
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

func TestEvaluationEventCountV2(t *testing.T) {
	t.Parallel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	ecClient := newEventCounterClient(t)
	defer ecClient.Close()
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
	variations := make(map[string]*featureproto.Variation)
	variationIDs := []string{}
	for _, v := range f.Variations {
		variationIDs = append(variationIDs, v.Id)
		variations[v.Value] = v
	}

	registerEvaluationEvent(t, featureID, f.Version, userID, variations[variationVarA].Id, tag)

	for i := 0; i < retryTimes; i++ {
		if i == retryTimes-1 {
			t.Fatalf("retry timeout")
		}
		time.Sleep(time.Second)

		resp := getEvaluationCountV2(t, ecClient, featureID, f.Version, variationIDs)
		if resp == nil {
			continue
		}
		ec := resp.Count
		if ec.FeatureId != featureID {
			t.Fatalf("feature ID is not correct: %s", ec.FeatureId)
		}
		if ec.FeatureVersion != f.Version {
			t.Fatalf("feature version is not correct: %d", ec.FeatureVersion)
		}

		vcA := getVariationCount(ec.RealtimeCounts, variations[variationVarA].Id)
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

		vcB := getVariationCount(ec.RealtimeCounts, variations[variationVarB].Id)
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
		Tag:       tag,
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
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	goal, err := protojson.Marshal(&eventproto.GoalEvent{
		Timestamp: time.Now().Unix(),
		GoalId:    goalID,
		UserId:    userID,
		Value:     value,
		User:      &userproto.User{},
		Tag:       tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	events := []util.Event{
		{
			ID:    newUUID(t),
			Event: goal,
			Type:  util.GoalEventType,
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
) {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	evaluation, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User:           &userproto.User{Data: map[string]string{"appVersion": "0.1.0"}},
		Reason:         &featureproto.Reason{},
		Tag:            tag,
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
) {
	t.Helper()
	evaluation, err := protojson.Marshal(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User:           &userproto.User{},
		Reason:         &featureproto.Reason{},
		Tag:            tag,
	})
	if err != nil {
		t.Fatal(err)
	}
	events := []util.Event{
		{
			ID:    newUUID(t),
			Event: evaluation,
			Type:  util.EvaluationEventType,
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

func getEvaluationCountV2(t *testing.T, c ecclient.Client, featureID string, featureVersion int32, variationIDs []string) *ecproto.GetEvaluationCountV2Response {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	now := time.Now()
	req := &ecproto.GetEvaluationCountV2Request{
		EnvironmentNamespace: *environmentNamespace,
		StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
		EndAt:                now.Unix(),
		FeatureId:            featureID,
		FeatureVersion:       featureVersion,
		VariationIds:         variationIDs,
	}
	response, err := c.GetEvaluationCountV2(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func getGoalCountV2(t *testing.T, c ecclient.Client, goalID, featureID string, featureVersion int32, variationIDs []string) *ecproto.GetGoalCountV2Response {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	now := time.Now()
	req := &ecproto.GetGoalCountV2Request{
		EnvironmentNamespace: *environmentNamespace,
		StartAt:              now.Add(-30 * 24 * time.Hour).Unix(),
		EndAt:                now.Unix(),
		GoalId:               goalID,
		FeatureId:            featureID,
		FeatureVersion:       featureVersion,
		VariationIds:         variationIDs,
	}
	response, err := c.GetGoalCountV2(ctx, req)
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
