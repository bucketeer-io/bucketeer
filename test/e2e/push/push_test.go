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

package autoops

import (
	"context"
	"flag"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"

	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	pushclient "github.com/bucketeer-io/bucketeer/v2/pkg/push/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
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

	fcmServiceAccountDummy = `{
		"type": "service_account",
		"project_id": "%s-%s",
		"private_key_id": "private-key-id",
		"private_key": "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n",
		"client_email": "fcm-service-account@test.iam.gserviceaccount.com",
		"client_id": "client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fcm-service-account@test.iam.gserviceaccount.com",
		"universe_domain": "googleapis.com"
	}`
)

func TestCreateAndListPush(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	featureClient := newFeatureClient(t)
	defer featureClient.Close()
	pushClient := newPushClient(t)
	defer pushClient.Close()

	featureID := newFeatureID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, newUUID(t))
	fcmServiceAccount := fmt.Sprintf(fcmServiceAccountDummy, prefixTestName, newUUID(t))

	createFeature(ctx, t, featureClient, featureID, tag)
	createPush(ctx, t, pushClient, []byte(fcmServiceAccount), tag)
	pushes := listPushes(t, pushClient)
	var push *pushproto.Push
	for _, p := range pushes {
		// Search the push by tag
		for _, t := range p.Tags {
			if t == tag {
				push = p
				break
			}
		}
	}
	if push == nil {
		t.Fatalf("Push not found")
	}
	if push.FcmServiceAccount != "" {
		t.Fatalf("The FCM service account must be empty. Actual: %s", push.FcmServiceAccount)
	}
	if len(push.Tags) != 1 {
		t.Fatalf("The number of tags is incorrect. Expected: %d actual: %d", 1, len(push.Tags))
	}
	if push.Tags[0] != tag {
		t.Fatalf("Incorrect tag. Expected: %s actual: %s", tag, push.Tags[0])
	}
	if push.Deleted != false {
		t.Fatalf("Incorrect deleted. Expected: %t actual: %t", false, push.Deleted)
	}
}

func TestCreateAndListPush_NoCommand_WithoutTags(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	pushClient := newPushClient(t)
	defer pushClient.Close()

	fcmServiceAccount := fmt.Sprintf(fcmServiceAccountDummy, prefixTestName, newUUID(t))
	pushName := newPushName(t)

	// Create push without tags using no-command API
	req := &pushproto.CreatePushRequest{
		EnvironmentId:     *environmentID,
		Name:              pushName,
		FcmServiceAccount: []byte(fcmServiceAccount),
		Tags:              []string{}, // Empty tags
	}
	resp, err := pushClient.CreatePush(ctx, req)
	if err != nil {
		t.Fatal("Failed to create push without tags:", err)
	}

	// Verify the push was created correctly
	push := resp.Push
	if push == nil {
		t.Fatal("Push response is nil")
	}
	if push.Name != pushName {
		t.Fatalf("Incorrect push name. Expected: %s actual: %s", pushName, push.Name)
	}
	if push.FcmServiceAccount != "" {
		t.Fatalf("The FCM service account must be empty. Actual: %s", push.FcmServiceAccount)
	}
	if len(push.Tags) != 0 {
		t.Fatalf("The number of tags is incorrect. Expected: %d actual: %d", 0, len(push.Tags))
	}
	if push.Deleted != false {
		t.Fatalf("Incorrect deleted. Expected: %t actual: %t", false, push.Deleted)
	}

	// Verify the push appears in the list
	pushes := listPushes(t, pushClient)
	var foundPush *pushproto.Push
	for _, p := range pushes {
		if p.Id == push.Id {
			foundPush = p
			break
		}
	}
	if foundPush == nil {
		t.Fatalf("Push not found in list")
	}
}

func TestUpdatePush(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	pushClient := newPushClient(t)
	defer pushClient.Close()

	// Create initial push
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, newUUID(t))
	newTag := fmt.Sprintf("%s-new-tag-%s", prefixTestName, newUUID(t))
	fcmServiceAccount := fmt.Sprintf(fcmServiceAccountDummy, prefixTestName, newUUID(t))
	newName := fmt.Sprintf("%s-updated", newPushName(t))

	createPush(ctx, t, pushClient, []byte(fcmServiceAccount), tag)
	pushes := listPushes(t, pushClient)
	var push *pushproto.Push
	for _, p := range pushes {
		for _, t := range p.Tags {
			if t == tag {
				push = p
				break
			}
		}
	}
	if push == nil {
		t.Fatalf("Push not found")
	}

	// Update the push
	updateReq := &pushproto.UpdatePushRequest{
		Id:            push.Id,
		EnvironmentId: *environmentID,
		Name:          &wrappers.StringValue{Value: newName},
		TagChanges: []*pushproto.TagChange{
			{
				Tag:        tag,
				ChangeType: pushproto.ChangeType_CREATE, // This will be ignored
			},
			{
				Tag:        newTag,
				ChangeType: pushproto.ChangeType_CREATE,
			},
		},
		Disabled: wrapperspb.Bool(true),
	}
	if _, err := pushClient.UpdatePush(ctx, updateReq); err != nil {
		t.Fatal("Failed to update push:", err)
	}

	// Verify the update
	pushes = listPushes(t, pushClient)
	var updatedPush *pushproto.Push
	for _, p := range pushes {
		if p.Id == push.Id {
			updatedPush = p
			break
		}
	}
	if updatedPush == nil {
		t.Fatalf("Updated push not found")
	}
	if updatedPush.Name != newName {
		t.Errorf("Push name not updated. Expected: %s, got: %s", newName, updatedPush.Name)
	}
	if len(updatedPush.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(updatedPush.Tags))
	}
	if !slices.Contains(updatedPush.Tags, newTag) {
		t.Errorf("Push tag not updated. Expected: %s, got: %s", newTag, updatedPush.Tags)
	}
	if !updatedPush.Disabled {
		t.Errorf("Push disabled status not updated. Expected: true, got: %v", updatedPush.Disabled)
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

func createFeature(ctx context.Context, t *testing.T, client featureclient.Client, featureID, tag string) {
	t.Helper()
	cmd := newCreateFeatureCommand(featureID, tag)
	createReq := &featureproto.CreateFeatureRequest{
		Command:       cmd,
		EnvironmentId: *environmentID,
	}
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
	enableFeature(t, featureID, client)
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

func newCreateFeatureCommand(featureID, tag string) *featureproto.CreateFeatureCommand {
	return &featureproto.CreateFeatureCommand{
		Id:          featureID,
		Name:        "e2e-test-push-feature-name",
		Description: "e2e-test-push-feature-description",
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
		Tags:                     []string{tag},
		DefaultOnVariationIndex:  &wrappers.Int32Value{Value: int32(0)},
		DefaultOffVariationIndex: &wrappers.Int32Value{Value: int32(1)},
	}
}

func newPushClient(t *testing.T) pushclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := pushclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create auto ops client:", err)
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

func createPush(
	ctx context.Context,
	t *testing.T,
	client pushclient.Client,
	fcmServiceAccount []byte,
	tag string,
) {
	t.Helper()
	createReq := &pushproto.CreatePushRequest{
		EnvironmentId:     *environmentID,
		Name:              newPushName(t),
		Tags:              []string{tag},
		FcmServiceAccount: fcmServiceAccount,
	}
	if _, err := client.CreatePush(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func listPushes(t *testing.T, client pushclient.Client) []*pushproto.Push {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListPushes(ctx, &pushproto.ListPushesRequest{
		EnvironmentId: *environmentID,
		PageSize:      int64(500),
	})
	if err != nil {
		t.Fatal("failed to list pushes", err)
	}
	return resp.Pushes
}

func newFeatureID(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, newUUID(t))
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, newUUID(t))
}

func newPushName(t *testing.T) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-push-name", prefixTestName, *testID)
	}
	return fmt.Sprintf("%s-push-name", prefixTestName)
}
