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

package gateway

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/stretchr/testify/assert"

	gwapi "github.com/bucketeer-io/bucketeer/pkg/api/api"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/e2e/util"
)

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
	// Update feature flag cache
	updateFeatueFlagCache(t)
	response := util.GetEvaluations(t, tag, userID, *gatewayAddr, *apiKeyPath)

	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if len(response.Evaluations.Evaluations) == 0 {
		t.Fatalf("Wrong evaluation size. Expected more than one, actual zero")
	}
	eval, err := findFeature(response.Evaluations.Evaluations, featureID)
	if err != nil {
		t.Fatalf("Failed to find evaluation. Error: %v", err)
	}
	reason := eval.Reason.Type
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
	// Update feature flag cache
	updateFeatueFlagCache(t)
	response := util.GetEvaluations(t, tag, userID, *gatewayAddr, *apiKeyPath)

	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if len(response.Evaluations.Evaluations) == 0 {
		t.Fatalf("Wrong evaluation size. Expected more than one, actual zero")
	}
	eval, err := findFeature(response.Evaluations.Evaluations, featureID)
	if err != nil {
		t.Fatalf("Failed to find evaluation. Error: %v", err)
	}
	reason := eval.Reason.Type
	if reason != featureproto.Reason_OFF_VARIATION {
		t.Fatalf("Reason doesn't match. Expected: %v, actual: %v", featureproto.Reason_OFF_VARIATION, reason)
	}
}

func TestGetEvaluationsFullState(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	featureID2 := newFeatureID(t, newUUID(t))
	createFeatureWithTag(t, tag, featureID2)
	time.Sleep(3 * time.Second)
	response := util.GetEvaluations(t, tag, userID, *gatewayAddr, *apiKeyPath)

	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	evaluationSize := len(response.Evaluations.Evaluations)
	if evaluationSize < 2 {
		t.Fatalf("Wrong evaluation size. Expected 2, actual: %d", evaluationSize)
	}
}

func TestGetEvaluation(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	featureID2 := newFeatureID(t, newUUID(t))
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
				ID:            newUUID(t),
				Event:         json.RawMessage(evaluation),
				EnvironmentId: "",
				Type:          gwapi.EvaluationEventType,
			},
			{
				ID:            newUUID(t),
				Event:         json.RawMessage(goal),
				EnvironmentId: "",
				Type:          gwapi.GoalEventType,
			},
		},
		*gatewayAddr,
		*apiKeyPath,
	)
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

// TODO: This is a temporary custom JSON unmarshaler to handle boolean fields sent as strings from the Android SDK.
// This addresses a compatibility issue that arose after switching from Envoy JSON transcoder to gRPC-Gateway.
// Reference: https://github.com/bucketeer-io/android-client-sdk/pull/230
func TestHTTPBooleanStringHandling(t *testing.T) {
	t.Parallel()
	// This test verifies that the custom JSON unmarshaler correctly handles
	// boolean strings sent from the Android SDK via HTTP requests,
	// specifically the userAttributesUpdated field in get_evaluations API

	testCases := []struct {
		name        string
		description string
		requestBody map[string]interface{}
	}{
		{
			name:        "boolean_string_true_lowercase",
			description: "userAttributesUpdated as 'true' string",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationsId": "test-evaluations-id",
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "true",
				},
			},
		},
		{
			name:        "boolean_string_false_lowercase",
			description: "userAttributesUpdated as 'false' string",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "false",
				},
			},
		},
		{
			name:        "boolean_string_true_titlecase",
			description: "userAttributesUpdated as 'True' string",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "True",
				},
			},
		},
		{
			name:        "boolean_string_false_titlecase",
			description: "userAttributesUpdated as 'False' string",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "False",
				},
			},
		},
		{
			name:        "boolean_string_true_uppercase",
			description: "userAttributesUpdated as 'TRUE' string",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "TRUE",
				},
			},
		},
		{
			name:        "boolean_string_false_uppercase",
			description: "userAttributesUpdated as 'FALSE' string",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "FALSE",
				},
			},
		},
		{
			name:        "boolean_string_numeric_1",
			description: "userAttributesUpdated as '1' string (true)",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "1",
				},
			},
		},
		{
			name:        "boolean_string_numeric_0",
			description: "userAttributesUpdated as '0' string (false)",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "0",
				},
			},
		},
		{
			name:        "empty_user_evaluation_condition",
			description: "Empty user evaluation condition object",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{},
			},
		},
		{
			name:        "null_user_data",
			description: "User data is null with boolean string",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id":   "test-user-id",
					"data": nil,
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": "false",
				},
			},
		},
		{
			name:        "actual_boolean_value",
			description: "userAttributesUpdated as actual boolean (not string)",
			requestBody: map[string]interface{}{
				"tag": "test-tag",
				"user": map[string]interface{}{
					"id": "test-user-id",
				},
				"userEvaluationCondition": map[string]interface{}{
					"evaluatedAt":           time.Now().Unix() - 60,
					"userAttributesUpdated": true, // actual boolean
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Make the HTTP request - should succeed for all cases
			response := util.GetEvaluationsRaw(t, tc.requestBody, *gatewayAddr, *apiKeyPath)

			// Verify we got a valid response (no error occurred)
			assert.NotNil(t, response, "Response should not be nil for test case: %s", tc.description)
		})
	}
}
