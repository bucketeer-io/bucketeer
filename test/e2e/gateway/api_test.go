// Copyright 2024 The Bucketeer Authors.
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

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	gwapi "github.com/bucketeer-io/bucketeer/pkg/gateway/api"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/e2e/util"
)

func TestGetEvaluationsWithoutCreatingFeature(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	response := util.GetEvaluations(t, tag, userID, *gatewayAddr, *apiKeyPath)

	if response.Evaluations != nil {
		evaluationSize := len(response.Evaluations.Evaluations)
		if evaluationSize > 0 {
			t.Fatalf("Different sizes. Expected: 0, actual: %v", evaluationSize)
		}
	}
}

func TestGetEvaluationsFeatureFlagEnabled(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	addTag(t, tag, featureID, client)
	enableFeature(t, featureID, client)
	time.Sleep(3 * time.Second)
	response := util.GetEvaluations(t, tag, userID, *gatewayAddr, *apiKeyPath)

	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	evaluationSize := len(response.Evaluations.Evaluations)
	if evaluationSize != 1 {
		t.Fatalf("Wrong evaluation size. Expected 1, actual: %d", evaluationSize)
	}
	reason := response.Evaluations.Evaluations[0].Reason.Type
	if reason != featureproto.Reason_DEFAULT {
		t.Fatalf("Reason doesn't match. Expected: %v, actual: %v", featureproto.Reason_DEFAULT, reason)
	}
}

func TestGetEvaluationsFeatureFlagDisabled(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	addTag(t, tag, featureID, client)
	time.Sleep(3 * time.Second)
	response := util.GetEvaluations(t, tag, userID, *gatewayAddr, *apiKeyPath)

	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	evaluationSize := len(response.Evaluations.Evaluations)
	if evaluationSize != 1 {
		t.Fatalf("Wrong evaluation size. Expected 1, actual: %d", evaluationSize)
	}
	reason := response.Evaluations.Evaluations[0].Reason.Type
	if reason != featureproto.Reason_OFF_VARIATION {
		t.Fatalf("Reason doesn't match. Expected: %v, actual: %v", featureproto.Reason_OFF_VARIATION, reason)
	}
}

func TestGetEvaluationsFullState(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t)
	defer c.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	featureID2 := fmt.Sprintf("%s-feature-id-%s", prefixTestName, newUUID(t))
	createFeatureWithTag(t, tag, featureID2)
	time.Sleep(3 * time.Second)
	response := util.GetEvaluations(t, tag, userID, *gatewayAddr, *apiKeyPath)

	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	evaluationSize := len(response.Evaluations.Evaluations)
	if evaluationSize != 2 {
		t.Fatalf("Wrong evaluation size. Expected 2, actual: %d", evaluationSize)
	}
}

func TestGetEvaluation(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t)
	defer c.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	featureID2 := fmt.Sprintf("%s-feature-id-%s", prefixTestName, newUUID(t))
	createFeatureWithTag(t, tag, featureID2)
	time.Sleep(3 * time.Second)
	response := util.GetEvaluation(t, tag, featureID2, userID, *gatewayAddr, *apiKeyPath)
	if response.Evaluation == nil {
		t.Fatal("Evaluation field is nil")
	}
	targetFeatureID := response.Evaluation.FeatureId
	if targetFeatureID != featureID2 {
		t.Fatalf("Wrong feature id. Expected: %s, actual: %s", featureID2, targetFeatureID)
	}
}

func TestRegisterEvents(t *testing.T) {
	t.Parallel()
	evaluation, err := protojson.Marshal(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      "feature-id",
		FeatureVersion: 1,
		UserId:         "user-id",
		VariationId:    "variation-id",
		User: &userproto.User{
			Id: "user-id",
		},
		Reason: &featureproto.Reason{},
		Tag:    "tag",
	})
	if err != nil {
		t.Fatal(err)
	}
	goal, err := protojson.Marshal(&eventproto.GoalEvent{
		Timestamp: time.Now().Unix(),
		GoalId:    "goal-id",
		UserId:    "user-id",
		Value:     0.3,
		User: &userproto.User{
			Id: "user-id",
		},
		Tag: "tag",
	})
	if err != nil {
		t.Fatal(err)
	}
	response := util.RegisterEvents(
		t,
		[]util.Event{
			{
				ID:                   newUUID(t),
				Event:                json.RawMessage(evaluation),
				EnvironmentNamespace: "",
				Type:                 gwapi.EvaluationEventType,
			},
			{
				ID:                   newUUID(t),
				Event:                json.RawMessage(goal),
				EnvironmentNamespace: "",
				Type:                 gwapi.GoalEventType,
			},
		},
		*gatewayAddr,
		*apiKeyPath,
	)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}
