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
	"context"
	"flag"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	gatewayclient "github.com/bucketeer-io/bucketeer/pkg/api/client"
	btclient "github.com/bucketeer-io/bucketeer/pkg/batch/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	btproto "github.com/bucketeer-io/bucketeer/proto/batch"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/util"
)

const (
	prefixTestName = "e2e-test"
	timeout        = 60 * time.Second
)

var (
	webGatewayAddr   = flag.String("web-gateway-addr", "", "Web gateway endpoint address")
	webGatewayPort   = flag.Int("web-gateway-port", 443, "Web gateway endpoint port")
	webGatewayCert   = flag.String("web-gateway-cert", "", "Web gateway crt file")
	apiKeyPath       = flag.String("api-key", "", "Client SDK API key for api-gateway")
	apiKeyServerPath = flag.String("api-key-server", "", "Server SDK API key for api-gateway")
	gatewayAddr      = flag.String("gateway-addr", "", "Gateway endpoint address")
	gatewayPort      = flag.Int("gateway-port", 443, "Gateway endpoint port")
	gatewayCert      = flag.String("gateway-cert", "", "Gateway crt file")
	serviceTokenPath = flag.String("service-token", "", "Service token path")
	environmentID    = flag.String("environment-id", "", "Environment id")
	organizationID   = flag.String("organization-id", "", "Organization ID")
	testID           = flag.String("test-id", "", "test ID")
)

func TestAPIKey(t *testing.T) {
	t.Parallel()
	creds, err := gatewayclient.NewPerRPCCredentials("testdata/invalid-apikey")
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	c, err := gatewayclient.NewClient(
		fmt.Sprintf("%s:%d", *gatewayAddr, *gatewayPort),
		*gatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(timeout),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create gateway client:", err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &gatewayproto.GetEvaluationsRequest{
		Tag:  "tag",
		User: &userproto.User{Id: "userID"},
	}
	response, err := c.GetEvaluations(ctx, req)
	assert.Nil(t, response)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("not ok")
	}
	assert.Equal(t, st.Code(), codes.PermissionDenied)
}

func TestGrpcGetFeatureFlags(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	// Feature 1
	featureID1 := newFeatureID(t, uuid)
	cmd1 := createFeatureWithTag(t, tag, featureID1)

	// Feature 2
	uuid = newUUID(t)
	featureID2 := newFeatureID(t, uuid)
	cmd2 := createFeature(t, client, newCreateFeatureCommand(featureID2))

	// Feature 3
	uuid = newUUID(t)
	featureID3 := newFeatureID(t, uuid)
	cmd3 := createFeatureWithTag(t, tag, featureID3)

	time.Sleep(15 * time.Second) // It must be higher than the `secondsForAdjustment`

	// Find feature with no tag and no features ID
	response := grpcGetFeatureFlags(t, "", "", 0)
	assert.NotEmpty(t, response.FeatureFlagsId)
	assert.True(t, len(response.Features) >= 3)
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.True(t, response.ForceUpdate)
	assert.True(t, findFeatureByID(t, cmd1.Id, response.Features))
	assert.True(t, findFeatureByID(t, cmd2.Id, response.Features))
	assert.True(t, findFeatureByID(t, cmd3.Id, response.Features))

	// Find feature with tag and no features ID
	response = grpcGetFeatureFlags(t, tag, "", 0)
	assert.NotEmpty(t, response.FeatureFlagsId)
	assert.Equal(t, 2, len(response.Features))
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.True(t, response.ForceUpdate)
	assert.True(t, findFeatureByID(t, cmd1.Id, response.Features))
	assert.True(t, findFeatureByID(t, cmd3.Id, response.Features))

	// Find feature with tag, with the same features ID, and requested at
	ffid := response.FeatureFlagsId
	response = grpcGetFeatureFlags(t, tag, response.FeatureFlagsId, response.RequestedAt)
	assert.Equal(t, ffid, response.FeatureFlagsId)
	assert.Equal(t, 0, len(response.Features))
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.False(t, response.ForceUpdate)

	// Find feature with tag, with the different features ID, and requested at
	response = grpcGetFeatureFlags(t, tag, "random-id", response.RequestedAt)
	assert.Equal(t, ffid, response.FeatureFlagsId)
	assert.Equal(t, 0, len(response.Features))
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.False(t, response.ForceUpdate)
}

func TestGrpcGetFeatureFlagsWithArchivedIDs(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	// Feature 1
	featureID1 := newFeatureID(t, uuid)
	cmd1 := createFeatureWithTag(t, tag, featureID1)

	// Feature 2
	uuid = newUUID(t)
	featureID2 := newFeatureID(t, uuid)
	createFeature(t, client, newCreateFeatureCommand(featureID2))

	// Feature 3
	uuid = newUUID(t)
	featureID3 := newFeatureID(t, uuid)
	cmd3 := createFeatureWithTag(t, tag, featureID3)

	time.Sleep(15 * time.Second) // It must be higher than the `secondsForAdjustment`

	// Find feature by tag with tag and no features ID
	requestFFID := ""
	response := grpcGetFeatureFlags(t, tag, requestFFID, 0)
	assert.NotEmpty(t, response.FeatureFlagsId)
	assert.Equal(t, 2, len(response.Features))
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.True(t, response.ForceUpdate)
	assert.True(t, findFeatureByID(t, cmd1.Id, response.Features))
	assert.True(t, findFeatureByID(t, cmd3.Id, response.Features))

	// Archive feature
	archiveFeature(t, cmd1.Id, client)

	// Update feature flag cache
	updateFeatueFlagCache(t)

	// Find the archived flag
	requestFFID = response.FeatureFlagsId
	response = grpcGetFeatureFlags(t, tag, requestFFID, response.RequestedAt)
	assert.True(t, requestFFID != response.FeatureFlagsId)
	assert.Equal(t, 0, len(response.Features))
	assert.Equal(t, 1, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.False(t, response.ForceUpdate)
	assert.Equal(t, cmd1.Id, response.ArchivedFeatureFlagIds[0])
}

func TestGrpcGetFeatureFlagsWithRequestedAt(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	// Feature 1
	featureID1 := newFeatureID(t, uuid)
	cmd1 := createFeatureWithTag(t, tag, featureID1)

	// Feature 2
	uuid = newUUID(t)
	featureID2 := newFeatureID(t, uuid)
	createFeature(t, client, newCreateFeatureCommand(featureID2))

	// Feature 3
	uuid = newUUID(t)
	featureID3 := newFeatureID(t, uuid)
	cmd3 := createFeatureWithTag(t, tag, featureID3)

	time.Sleep(15 * time.Second) // It must be higher than the `secondsForAdjustment`

	// Find feature by tag with tag and no features ID
	requestFFID := ""
	response := grpcGetFeatureFlags(t, tag, requestFFID, 0)
	assert.NotEmpty(t, response.FeatureFlagsId)
	assert.Equal(t, 2, len(response.Features))
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.True(t, response.ForceUpdate)
	assert.True(t, findFeatureByID(t, cmd1.Id, response.Features))
	assert.True(t, findFeatureByID(t, cmd3.Id, response.Features))

	// Create another Flag
	// Feature 4
	uuid = newUUID(t)
	featureID4 := newFeatureID(t, uuid)
	cmd4 := createFeatureWithTag(t, tag, featureID4)

	// Find the flag 4
	requestFFID = response.FeatureFlagsId
	response = grpcGetFeatureFlags(t, tag, requestFFID, response.RequestedAt)
	assert.True(t, requestFFID != response.FeatureFlagsId)
	assert.Equal(t, 1, len(response.Features))
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.False(t, response.ForceUpdate)
	assert.True(t, findFeatureByID(t, cmd4.Id, response.Features))
}

func TestGrpcGetFeatureFlagsWithRequestedAt31daysAgo(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	// Feature 1
	featureID1 := newFeatureID(t, uuid)
	cmd1 := createFeatureWithTag(t, tag, featureID1)

	// Feature 2
	uuid = newUUID(t)
	featureID2 := newFeatureID(t, uuid)
	createFeature(t, client, newCreateFeatureCommand(featureID2))

	// Feature 3
	uuid = newUUID(t)
	featureID3 := newFeatureID(t, uuid)
	cmd3 := createFeatureWithTag(t, tag, featureID3)

	// Find feature by tag with tag with random id, and old requested at
	requestFFID := "random-id"
	requestedAt := time.Now().Add(-31 * 24 * time.Hour).Unix()
	response := grpcGetFeatureFlags(t, tag, requestFFID, requestedAt)
	assert.True(t, requestFFID != response.FeatureFlagsId)
	assert.Equal(t, 2, len(response.Features))
	assert.Equal(t, 0, len(response.ArchivedFeatureFlagIds))
	assert.True(t, response.RequestedAt >= time.Now().Add(-30*time.Second).Unix())
	assert.True(t, response.ForceUpdate)
	assert.True(t, findFeatureByID(t, cmd1.Id, response.Features))
	assert.True(t, findFeatureByID(t, cmd3.Id, response.Features))
}

func findFeatureByID(t *testing.T, id string, features []*featureproto.Feature) bool {
	t.Helper()
	for _, f := range features {
		if id == f.Id {
			return true
		}
	}
	return false
}

func TestGrpcGetEvaluationsFeatureFlagEnabled(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	cmd := createFeatureWithTag(t, tag, featureID)
	response := grpcGetEvaluations(t, tag, userID)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
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
	cmdVariation := cmd.Variations[0]
	variationValue := eval.VariationValue
	if variationValue != cmdVariation.Value {
		t.Fatalf("Variation value doesn't match. Expected: %s, actual: %s", variationValue, cmdVariation.Value)
	}
	variationName := eval.VariationName
	if variationName != cmdVariation.Name {
		t.Fatalf("Variation name doesn't match. Expected: %s, actual: %s", variationName, cmdVariation.Name)
	}
	valueDescription := eval.Variation.Description
	if valueDescription != "" {
		t.Fatalf("Variation description is not empty. Actual: %s", valueDescription)
	}
}

func TestGrpcGetEvaluationsFeatureFlagDisabled(t *testing.T) {
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

	response := grpcGetEvaluations(t, tag, userID)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
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
	cmdVariation := cmd.Variations[1]
	variationValue := eval.VariationValue
	if variationValue != cmdVariation.Value {
		t.Fatalf("Variation value doesn't match. Expected: %s, actual: %s", variationValue, cmdVariation.Value)
	}
	variationName := eval.VariationName
	if variationName != cmdVariation.Name {
		t.Fatalf("Variation name doesn't match. Expected: %s, actual: %s", variationName, cmdVariation.Name)
	}
	valueDescription := eval.Variation.Description
	if valueDescription != "" {
		t.Fatalf("Variation description is not empty. Actual: %s", valueDescription)
	}
}

func TestGrpcGetEvaluationsFullState(t *testing.T) {
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
	response := grpcGetEvaluations(t, tag, userID)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	evaluationSize := len(response.Evaluations.Evaluations)
	if evaluationSize < 2 {
		t.Fatalf("Wrong evaluation size. Expected more than two, actual: %d", evaluationSize)
	}
}

func TestGrpcGetEvaluationsByEvaluatedAt(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	time.Sleep(10 * time.Second) // It must be equal or higher than the `secondsForAdjustment`
	featureID2 := newFeatureID(t, newUUID(t))
	cmd := createFeatureWithTag(t, tag, featureID2)
	prevEvalAt := time.Now().Unix()
	response := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, "userEvaluationsID", prevEvalAt, false)
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if len(response.Evaluations.Evaluations) == 0 {
		t.Fatalf("Wrong evaluation size. Expected more than one, actual zero")
	}
	if contains(response.Evaluations.Evaluations, featureID) {
		t.Fatalf("Evaluation should not contain the evaluation of feature: %s", featureID)
	}
	if !contains(response.Evaluations.Evaluations, featureID2) {
		t.Fatalf("Evaluation should contain the evaluation of feature: %s", featureID2)
	}
	if response.Evaluations.ForceUpdate {
		t.Fatal("ForceUpdate should be false")
	}
	eval, err := findFeature(response.Evaluations.Evaluations, featureID2)
	if err != nil {
		t.Fatalf("Failed to find evaluation. Error: %v", err)
	}
	cmdVariation := cmd.Variations[0]
	variationValue := eval.VariationValue
	if variationValue != cmdVariation.Value {
		t.Fatalf("Variation value doesn't match. Expected: %s, actual: %s", variationValue, cmdVariation.Value)
	}
	variationName := eval.VariationName
	if variationName != cmdVariation.Name {
		t.Fatalf("Variation name doesn't match. Expected: %s, actual: %s", variationName, cmdVariation.Name)
	}
	valueDescription := eval.Variation.Description
	if valueDescription != "" {
		t.Fatalf("Variation description is not empty. Actual: %s", valueDescription)
	}
}

func TestGrpcGetEvaluationsByEvaluatedAtIncludingArchivedFeature(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	fc := newFeatureClient(t)
	defer fc.Close()

	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, fc, cmd)
	addTag(t, tag, featureID, fc)
	enableFeature(t, featureID, fc)
	archiveFeature(t, featureID, fc)
	time.Sleep(10 * time.Second) // It must be equal or higher than the `secondsForAdjustment`

	uuid2 := newUUID(t)
	featureID2 := newFeatureID(t, uuid2)
	cmd2 := newCreateFeatureCommand(featureID2)
	createFeature(t, fc, cmd2)
	addTag(t, tag, featureID2, fc)
	enableFeature(t, featureID2, fc)
	archiveFeature(t, featureID2, fc)

	// Update feature flag cache
	updateFeatueFlagCache(t)

	prevEvalAt := time.Now().Unix()
	response := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, "userEvaluationsID", prevEvalAt, false)
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if len(response.Evaluations.ArchivedFeatureIds) == 0 {
		t.Fatal("Evaluation is empty")
	}
	containsFeatureID := false
	containsFeatureID2 := false
	for _, archivedID := range response.Evaluations.ArchivedFeatureIds {
		if archivedID == featureID {
			containsFeatureID = true
		}
		if archivedID == featureID2 {
			containsFeatureID2 = true
		}
	}
	if containsFeatureID {
		t.Fatalf("ArchivedFeaturesIds should not contain %s", featureID)
	}
	if !containsFeatureID2 {
		t.Fatalf("ArchivedFeaturesIds should contain %s", featureID2)
	}
	if response.Evaluations.ForceUpdate {
		t.Fatal("ForceUpdate should be false")
	}
}

func TestGrpcGetEvaluationsByUserAttributesUpdated(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	time.Sleep(10 * time.Second) // It must be equal or higher than the `secondsForAdjustment`
	featureID2 := newFeatureID(t, newUUID(t))
	createFeatureWithRule(t, tag, featureID2)
	prevEvalAt := time.Now().Unix()
	response := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, "userEvaluationsID", prevEvalAt, true)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if len(response.Evaluations.Evaluations) == 0 {
		t.Fatal("Evaluation is empty")
	}
	if contains(response.Evaluations.Evaluations, featureID) {
		t.Fatalf("Evaluation should not contain the evaluation of feature that doesn't have rules: %s", featureID)
	}
	if !contains(response.Evaluations.Evaluations, featureID2) {
		t.Fatalf("Evaluation should contain the evaluation of feature that has rules: %s", featureID2)
	}
	if response.Evaluations.ForceUpdate {
		t.Fatal("ForceUpdate should be false")
	}
}

func TestGrpcGetEvaluationsWithPreviousEvaluation31daysAgo(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	userID := newUserID(t, uuid)
	prevEvalAt := time.Now().Add(-31 * 24 * time.Hour).Unix()
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	updateFeatueFlagCache(t)
	response := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, "userEvaluationsID", prevEvalAt, false)
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if !response.Evaluations.ForceUpdate {
		t.Fatal("ForceUpdate should be true because the previous evaluation is performed 31days ago")
	}
}

func TestGrpcGetEvaluationsWithEvaluatedAtIsZero(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	userID := newUserID(t, uuid)
	var prevEvalAt int64 = 0
	userEvaluationsID := ""
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	updateFeatueFlagCache(t)
	response := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, userEvaluationsID, prevEvalAt, false)
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if !response.Evaluations.ForceUpdate {
		t.Fatal("ForceUpdate should be true because the evaluatedAt is zero, which means that previous evaluation is performed more than 30 days ago.")
	}
}

func TestGrpcGetEvaluationsWithEmptyUserEvaluationsID(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	userID := newUserID(t, uuid)
	prevEvalAt := time.Now().Add(-1 * time.Second).Unix()
	userEvaluationsID := ""
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	updateFeatueFlagCache(t)
	response := grpcGetEvaluationsByEvaluatedAt(t, tag, userID, userEvaluationsID, prevEvalAt, false)
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if !response.Evaluations.ForceUpdate {
		t.Fatal("ForceUpdate should be true because the UserEvaluationsID is empty")
	}
}

func TestGrpcGetEvaluationsWithoutTag(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	uuid := newUUID(t)
	userID := newUserID(t, uuid)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)
	uuid2 := newUUID(t)
	tag2 := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid2)
	featureID2 := newFeatureID(t, uuid2)
	createFeatureWithTag(t, tag2, featureID2)

	prevEvalAt := time.Now().Add(-5 * time.Minute).Unix()
	response := grpcGetEvaluationsByEvaluatedAt(t, "", userID, "userEvaluationsID", prevEvalAt, false)
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	if len(response.Evaluations.Evaluations) == 0 {
		t.Fatalf("Wrong evaluation size. Expected more than one, actual zero")
	}
	if !contains(response.Evaluations.Evaluations, featureID) {
		t.Fatalf("Evaluation should contain the evaluation of feature: %s", featureID)
	}
	if !contains(response.Evaluations.Evaluations, featureID2) {
		t.Fatalf("Evaluation should contain the evaluation of feature: %s", featureID2)
	}
	if response.Evaluations.ForceUpdate {
		t.Fatal("ForceUpdate should be false")
	}
}

func TestGrpcGetEvaluation(t *testing.T) {
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
	response := grpcGetEvaluation(t, tag, featureID2, userID)
	if response.Evaluation == nil {
		t.Fatal("Evaluation field is nil")
	}
	targetFeatureID := response.Evaluation.FeatureId
	if targetFeatureID != featureID2 {
		t.Fatalf("Wrong feature id. Expected: %s, actual: %s", featureID2, targetFeatureID)
	}
}

func TestGrpcRegisterEvents(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Evaluation Event
	evaluation, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
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
	// GoalEvent
	goal, err := ptypes.MarshalAny(&eventproto.GoalEvent{
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
	// InternalSDKErrorMetricsEvent
	internalSDKErr, err := ptypes.MarshalAny(&eventproto.InternalSdkErrorMetricsEvent{
		ApiId:  eventproto.ApiId_GET_EVALUATIONS,
		Labels: map[string]string{"tag": "IOS"},
	})
	if err != nil {
		t.Fatal(err)
	}
	metricsInternalSDK, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
		Timestamp:  time.Now().Unix(),
		Event:      internalSDKErr,
		SdkVersion: "v0.0.1-e2e",
		SourceId:   eventproto.SourceId_IOS,
	})
	if err != nil {
		t.Fatal(err)
	}
	// BadRequestErrorMetricsEvent
	badRequestErr, err := ptypes.MarshalAny(&eventproto.BadRequestErrorMetricsEvent{
		ApiId:  eventproto.ApiId_REGISTER_EVENTS,
		Labels: map[string]string{"tag": "ANDROID"},
	})
	if err != nil {
		t.Fatal(err)
	}
	metricsBadRequest, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
		Timestamp:  time.Now().Unix(),
		Event:      badRequestErr,
		SdkVersion: "v0.0.1-e2e",
		SourceId:   eventproto.SourceId_ANDROID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// SizeMetricsEvent
	size, err := ptypes.MarshalAny(&eventproto.SizeMetricsEvent{
		ApiId:    eventproto.ApiId_REGISTER_EVENTS,
		Labels:   map[string]string{"tag": "JAVASCRIPT"},
		SizeByte: 99,
	})
	if err != nil {
		t.Fatal(err)
	}
	metricsSize, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
		Timestamp:  time.Now().Unix(),
		Event:      size,
		SdkVersion: "v0.0.1-e2e",
		SourceId:   eventproto.SourceId_JAVASCRIPT,
	})
	if err != nil {
		t.Fatal(err)
	}
	// LatencyMetricsEvent
	latency, err := ptypes.MarshalAny(&eventproto.LatencyMetricsEvent{
		ApiId:    eventproto.ApiId_REGISTER_EVENTS,
		Labels:   map[string]string{"tag": "GO_SERVER"},
		Duration: durationpb.New(time.Duration(99)),
	})
	if err != nil {
		t.Fatal(err)
	}
	metricsLatency, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
		Timestamp:  time.Now().Unix(),
		Event:      latency,
		SdkVersion: "v0.0.1-e2e",
		SourceId:   eventproto.SourceId_GO_SERVER,
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
			{
				Id:    newUUID(t),
				Event: goal,
			},
			{
				Id:    newUUID(t),
				Event: metricsInternalSDK,
			},
			{
				Id:    newUUID(t),
				Event: metricsBadRequest,
			},
			{
				Id:    newUUID(t),
				Event: metricsSize,
			},
			{
				Id:    newUUID(t),
				Event: metricsLatency,
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

func TestRegisterEventsForMetricsEvent(t *testing.T) {
	t.Parallel()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	sdkVersion := "v0.0.1-e2e"
	apiIDs := []eventproto.ApiId{
		eventproto.ApiId_GET_EVALUATION,
		eventproto.ApiId_GET_EVALUATIONS,
		eventproto.ApiId_REGISTER_EVENTS,
	}
	sourceIds := []eventproto.SourceId{
		eventproto.SourceId_ANDROID,
		eventproto.SourceId_IOS,
		eventproto.SourceId_GO_SERVER,
		eventproto.SourceId_NODE_SERVER,
		eventproto.SourceId_JAVASCRIPT,
	}
	events := make([]*eventproto.Event, 0, 0)
	rand.Seed(time.Now().UnixNano())
	for _, apiID := range apiIDs {
		for _, sourceID := range sourceIds {
			// InternalSDKErrorMetricsEvent
			internalSDKErr, err := ptypes.MarshalAny(&eventproto.InternalSdkErrorMetricsEvent{
				ApiId:  apiID,
				Labels: map[string]string{"tag": sourceID.String()},
			})
			if err != nil {
				t.Fatal(err)
			}
			metricsInternalSDK, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
				Timestamp:  time.Now().Unix(),
				Event:      internalSDKErr,
				SdkVersion: sdkVersion,
				SourceId:   sourceID,
			})
			if err != nil {
				t.Fatal(err)
			}
			events = append(events, &eventproto.Event{Id: newUUID(t), Event: metricsInternalSDK})
			// BadRequestErrorMetricsEvent
			badRequestErr, err := ptypes.MarshalAny(&eventproto.BadRequestErrorMetricsEvent{
				ApiId:  apiID,
				Labels: map[string]string{"tag": sourceID.String()},
			})
			if err != nil {
				t.Fatal(err)
			}
			metricsBadRequest, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
				Timestamp:  time.Now().Unix(),
				Event:      badRequestErr,
				SdkVersion: sdkVersion,
				SourceId:   sourceID,
			})
			if err != nil {
				t.Fatal(err)
			}
			events = append(events, &eventproto.Event{Id: newUUID(t), Event: metricsBadRequest})
			// SizeMetricsEvent
			size, err := ptypes.MarshalAny(&eventproto.SizeMetricsEvent{
				ApiId:    apiID,
				Labels:   map[string]string{"tag": sourceID.String()},
				SizeByte: rand.Int31n(100),
			})
			if err != nil {
				t.Fatal(err)
			}
			metricsSize, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
				Timestamp:  time.Now().Unix(),
				Event:      size,
				SdkVersion: sdkVersion,
				SourceId:   sourceID,
			})
			if err != nil {
				t.Fatal(err)
			}
			events = append(events, &eventproto.Event{Id: newUUID(t), Event: metricsSize})
			// LatencyMetricsEvent
			latency, err := ptypes.MarshalAny(&eventproto.LatencyMetricsEvent{
				ApiId:         apiID,
				Labels:        map[string]string{"tag": sourceID.String()},
				LatencySecond: rand.Float64(),
			})
			if err != nil {
				t.Fatal(err)
			}
			metricsLatency, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
				Timestamp:  time.Now().Unix(),
				Event:      latency,
				SdkVersion: sdkVersion,
				SourceId:   sourceID,
			})
			if err != nil {
				t.Fatal(err)
			}
			events = append(events, &eventproto.Event{Id: newUUID(t), Event: metricsLatency})
		}
	}

	req := &gatewayproto.RegisterEventsRequest{
		Events: events,
	}
	response, err := c.RegisterEvents(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Errors) > 0 {
		t.Fatalf("Failed to register events. Error: %v", response.Errors)
	}
}

func TestGetUserAttributeKeys(t *testing.T) {
	t.Parallel()
	client := newFeatureClient(t)
	defer client.Close()
	uuid := newUUID(t)
	environmentId := *environmentID
	maxRetryCount := 5
	sleepSecond := 30

	// Create some evaluation events to populate user attributes cache
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	featureID := newFeatureID(t, uuid)
	createFeatureWithTag(t, tag, featureID)

	// Get the feature to retrieve correct variation ID and feature version
	feature := getFeature(t, featureID, client)
	if len(feature.Variations) == 0 {
		t.Fatal("Feature has no variations")
	}
	variationID := feature.Variations[0].Id
	featureVersion := feature.Version

	// Register evaluation events with user attributes
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(maxRetryCount*sleepSecond)*time.Second)
	defer cancel()

	testUserDataKeySuffix := "testGetUserAttributeKeys-"

	// First evaluation event with 3 attributes
	data1 := map[string]string{
		testUserDataKeySuffix + "attr1-" + uuid: "value1",
		testUserDataKeySuffix + "attr2-" + uuid: "value2",
		testUserDataKeySuffix + "attr3-" + uuid: "value3",
	}

	evaluation1, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User: &userproto.User{
			Id:   userID,
			Data: data1,
		},
		Reason: &featureproto.Reason{
			Type: featureproto.Reason_CLIENT,
		},
		Tag:        tag,
		SdkVersion: "v0.0.1-e2e",
		SourceId:   eventproto.SourceId_ANDROID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Second evaluation event with different attributes and source ID
	data2 := map[string]string{
		testUserDataKeySuffix + "attr4-" + uuid: "value4",
		testUserDataKeySuffix + "attr5-" + uuid: "value5",
		testUserDataKeySuffix + "attr6-" + uuid: "value6",
	}

	evaluation2, err := ptypes.MarshalAny(&eventproto.EvaluationEvent{
		Timestamp:      time.Now().Unix(),
		FeatureId:      featureID,
		FeatureVersion: featureVersion,
		UserId:         userID,
		VariationId:    variationID,
		User: &userproto.User{
			Id:   userID,
			Data: data2,
		},
		Reason: &featureproto.Reason{
			Type: featureproto.Reason_CLIENT,
		},
		Tag:        tag,
		SdkVersion: "v0.0.1-e2e",
		SourceId:   eventproto.SourceId_IOS,
	})
	if err != nil {
		t.Fatal(err)
	}

	req := &gatewayproto.RegisterEventsRequest{
		Events: []*eventproto.Event{
			{
				Id:    newUUID(t),
				Event: evaluation1,
			},
			{
				Id:    newUUID(t),
				Event: evaluation2,
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

	// Check that all 6 attributes are found
	expectedAttributes := []string{
		testUserDataKeySuffix + "attr1-" + uuid,
		testUserDataKeySuffix + "attr2-" + uuid,
		testUserDataKeySuffix + "attr3-" + uuid,
		testUserDataKeySuffix + "attr4-" + uuid,
		testUserDataKeySuffix + "attr5-" + uuid,
		testUserDataKeySuffix + "attr6-" + uuid,
	}

	foundAttributes := make(map[string]bool)
	for i := 0; i < maxRetryCount; i++ {
		// Test GetUserAttributeKeys API
		userAttrReq := &featureproto.GetUserAttributeKeysRequest{
			EnvironmentId: environmentId,
		}
		userAttrResp, err := client.GetUserAttributeKeys(ctx, userAttrReq)
		if err != nil {
			t.Fatal("Failed to get user attribute keys:", err)
		}

		// Check for all expected attributes
		for _, key := range userAttrResp.UserAttributeKeys {
			for _, expectedAttr := range expectedAttributes {
				if key == expectedAttr {
					foundAttributes[expectedAttr] = true
				}
			}
		}

		// If all attributes are found, break
		if len(foundAttributes) == len(expectedAttributes) {
			break
		}

		time.Sleep(time.Duration(sleepSecond) * time.Second) // Wait for cache to update
	}

	// Verify all expected attributes were found
	for _, expectedAttr := range expectedAttributes {
		if !foundAttributes[expectedAttr] {
			t.Errorf("User attribute key '%s' not found after %d retries", expectedAttr, maxRetryCount)
		}
	}

	if len(foundAttributes) != len(expectedAttributes) {
		t.Fatalf("Expected %d attributes, but found %d", len(expectedAttributes), len(foundAttributes))
	}
}

func newGatewayClient(t *testing.T, apiKey string) gatewayclient.Client {
	t.Helper()
	creds, err := gatewayclient.NewPerRPCCredentials(apiKey)
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

func createFeatureWithTag(t *testing.T, tag, featureID string) *featureproto.CreateFeatureCommand {
	client := newFeatureClient(t)
	defer client.Close()
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	addTag(t, tag, cmd.Id, client)
	enableFeature(t, featureID, client)
	// Update feature flag cache
	updateFeatueFlagCache(t)
	return cmd
}

func createFeatureWithRule(t *testing.T, tag, featureID string) {
	client := newFeatureClient(t)
	defer client.Close()
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	addTag(t, tag, cmd.Id, client)
	addRule(t, cmd.Id, getFeature(t, featureID, client).Variations[1].Id, client)
	enableFeature(t, featureID, client)
	// Update feature flag cache
	updateFeatueFlagCache(t)
}

func updateFeatueFlagCache(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	batchClient := newBatchClient(t)
	defer batchClient.Close()
	numRetries := 5
	var err error
	for i := 0; i < numRetries; i++ {
		_, err = batchClient.ExecuteBatchJob(
			ctx,
			&btproto.BatchJobRequest{Job: btproto.BatchJob_FeatureFlagCacher})
		if err == nil {
			break
		}
		st, _ := status.FromError(err)
		if st.Code() == codes.Internal {
			t.Fatal(err)
		}
		fmt.Printf("Failed to execute feature flag cacher batch. Error code: %d\n. Retrying in 5 seconds.", st.Code())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		t.Fatal(err)
	}
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

func newCreateFeatureCommand(featureID string) *featureproto.CreateFeatureCommand {
	return &featureproto.CreateFeatureCommand{
		Id:          featureID,
		Name:        featureID,
		Description: "e2e-test-gateway-feature-description",
		Variations: []*featureproto.Variation{
			{
				Value:       "A",
				Name:        "Variation A",
				Description: "Thing does A",
			},
			{
				Value:       "B",
				Name:        "Variation B",
				Description: "Thing does B",
			},
		},
		Tags: []string{
			"e2e-test-tag-1",
			"e2e-test-tag-2",
			"e2e-test-tag-3",
		},
		DefaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
		DefaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
	}
}

func createFeature(
	t *testing.T,
	client featureclient.Client,
	cmd *featureproto.CreateFeatureCommand,
) *featureproto.CreateFeatureCommand {
	t.Helper()
	createReq := &featureproto.CreateFeatureRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
	return cmd
}

func getFeature(t *testing.T, featureID string, client featureclient.Client) *featureproto.Feature {
	t.Helper()
	getReq := &featureproto.GetFeatureRequest{
		Id:            featureID,
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.GetFeature(ctx, getReq)
	if err != nil {
		t.Fatal("Failed to get feature:", err)
	}
	return response.Feature
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

func addRule(t *testing.T, featureID, variationID string, client featureclient.Client) {
	t.Helper()
	rule := newFixedStrategyRule(variationID)
	addCmd, _ := util.MarshalCommand(&featureproto.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
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

func archiveFeature(t *testing.T, featureID string, client featureclient.Client) {
	t.Helper()
	req := &featureproto.ArchiveFeatureRequest{
		Id:            featureID,
		Command:       &featureproto.ArchiveFeatureCommand{},
		EnvironmentId: *environmentID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.ArchiveFeature(ctx, req); err != nil {
		t.Fatal(err)
	}
}

func grpcGetEvaluations(t *testing.T, tag, userID string) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	c := newGatewayClient(t, *apiKeyPath)
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

func grpcGetFeatureFlags(t *testing.T, tag, featuresID string, requestedAt int64) *gatewayproto.GetFeatureFlagsResponse {
	t.Helper()
	c := newGatewayClient(t, *apiKeyServerPath)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &gatewayproto.GetFeatureFlagsRequest{
		Tag:            tag,
		FeatureFlagsId: featuresID,
		RequestedAt:    requestedAt,
		SourceId:       eventproto.SourceId_GO_SERVER,
		SdkVersion:     "v0.0.1-e2e-test",
	}
	response, err := c.GetFeatureFlags(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func grpcGetEvaluationsByEvaluatedAt(
	t *testing.T,
	tag, userID, userEvaluationsID string,
	evaluatedAt int64,
	userAttributesUpdated bool,
) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &gatewayproto.GetEvaluationsRequest{
		UserEvaluationsId: userEvaluationsID,
		User:              &userproto.User{Id: userID},
		UserEvaluationCondition: &gatewayproto.GetEvaluationsRequest_UserEvaluationCondition{
			EvaluatedAt:           evaluatedAt,
			UserAttributesUpdated: userAttributesUpdated,
		},
		Tag: tag,
	}
	response, err := c.GetEvaluations(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func grpcGetEvaluation(t *testing.T, tag, featureID, userID string) *gatewayproto.GetEvaluationResponse {
	t.Helper()
	c := newGatewayClient(t, *apiKeyPath)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &gatewayproto.GetEvaluationRequest{
		Tag:       tag,
		User:      &userproto.User{Id: userID},
		FeatureId: featureID,
	}
	response, err := c.GetEvaluation(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func newUserID(t *testing.T, uuid string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-user-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-user-%s", prefixTestName, uuid)
}

func newFeatureID(t *testing.T, uuid string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, uuid)
}

func newFixedStrategyRule(variationID string) *featureproto.Rule {
	uuid, _ := uuid.NewUUID()
	return &featureproto.Rule{
		Id: uuid.String(),
		Strategy: &featureproto.Strategy{
			Type: featureproto.Strategy_FIXED,
			FixedStrategy: &featureproto.FixedStrategy{
				Variation: variationID,
			},
		},
		Clauses: []*featureproto.Clause{
			{
				Attribute: "attribute-1",
				Operator:  featureproto.Clause_EQUALS,
				Values:    []string{"value-1", "value-2"},
			},
			{
				Attribute: "attribute-2",
				Operator:  featureproto.Clause_IN,
				Values:    []string{"value-1", "value-2"},
			},
		},
	}
}

func updateFeatureTargeting(t *testing.T, client featureclient.Client, cmd *any.Any, featureID string) {
	t.Helper()
	updateReq := &featureproto.UpdateFeatureTargetingRequest{
		Id: featureID,
		Commands: []*featureproto.Command{
			{Command: cmd},
		},
		EnvironmentId: *environmentID,
		From:          featureproto.UpdateFeatureTargetingRequest_USER,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureTargeting(ctx, updateReq); err != nil {
		t.Fatal(err)
	}
}

func contains(evaluations []*featureproto.Evaluation, id string) bool {
	for _, e := range evaluations {
		if e.FeatureId == id {
			return true
		}
	}
	return false
}

func findFeature(
	evaluations []*featureproto.Evaluation,
	id string,
) (*featureproto.Evaluation, error) {
	for _, e := range evaluations {
		if e.FeatureId == id {
			return e, nil
		}
	}
	return nil, fmt.Errorf("evaluation not found")
}
