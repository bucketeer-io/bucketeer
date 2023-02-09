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

package gateway

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	gatewayclient "github.com/bucketeer-io/bucketeer/pkg/gateway/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	prefixTestName = "e2e-test"
	timeout        = 20 * time.Second
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

func TestGrpcGetEvaluationsWithoutCreatingFeature(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	response := grpcGetEvaluations(t, tag, userID)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
	if response.Evaluations != nil {
		evaluationSize := len(response.Evaluations.Evaluations)
		if evaluationSize > 0 {
			t.Fatalf("Different sizes. Expected: 0, actual: %v", evaluationSize)
		}
	}
}

func TestGrpcGetEvaluationsFeatureFlagEnabled(t *testing.T) {
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
	response := grpcGetEvaluations(t, tag, userID)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
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
	time.Sleep(3 * time.Second)
	response := grpcGetEvaluations(t, tag, userID)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
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

func TestGrpcGetEvaluationsFullState(t *testing.T) {
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
	response := grpcGetEvaluations(t, tag, userID)
	if response.State != featureproto.UserEvaluations_FULL {
		t.Fatalf("Different states. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, response.State)
	}
	if response.Evaluations == nil {
		t.Fatal("Evaluations field is nil")
	}
	evaluationSize := len(response.Evaluations.Evaluations)
	if evaluationSize != 2 {
		t.Fatalf("Wrong evaluation size. Expected 2, actual: %d", evaluationSize)
	}
}

func TestGrpcGetEvaluation(t *testing.T) {
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
	c := newGatewayClient(t)
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
		Labels: map[string]string{"tag": "iOS"},
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
		Labels: map[string]string{"tag": "Android"},
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
		Labels:   map[string]string{"tag": "web"},
		SizeByte: 99,
	})
	if err != nil {
		t.Fatal(err)
	}
	metricsSize, err := ptypes.MarshalAny(&eventproto.MetricsEvent{
		Timestamp:  time.Now().Unix(),
		Event:      size,
		SdkVersion: "v0.0.1-e2e",
		SourceId:   eventproto.SourceId_WEB,
	})
	if err != nil {
		t.Fatal(err)
	}
	// LatencyMetricsEvent
	latency, err := ptypes.MarshalAny(&eventproto.LatencyMetricsEvent{
		ApiId:    eventproto.ApiId_REGISTER_EVENTS,
		Labels:   map[string]string{"tag": "go-server-sdk"},
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

func createFeatureWithTag(t *testing.T, tag, featureID string) {
	client := newFeatureClient(t)
	defer client.Close()
	cmd := newCreateFeatureCommand(featureID)
	createFeature(t, client, cmd)
	addTag(t, tag, cmd.Id, client)
	enableFeature(t, featureID, client)
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

func grpcGetEvaluations(t *testing.T, tag, userID string) *gatewayproto.GetEvaluationsResponse {
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

func grpcGetEvaluation(t *testing.T, tag, featureID, userID string) *gatewayproto.GetEvaluationResponse {
	t.Helper()
	c := newGatewayClient(t)
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
