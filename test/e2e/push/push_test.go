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

package autoops

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	pushclient "github.com/bucketeer-io/bucketeer/pkg/push/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
)

const (
	prefixTestName = "e2e-test"
	timeout        = 10 * time.Second
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
	fcmAPIKey := fmt.Sprintf("%s-fcm-api-key-%s", prefixTestName, newUUID(t))
	createFeature(ctx, t, featureClient, featureID, tag)
	createPush(ctx, t, pushClient, fcmAPIKey, tag)
	pushes := listPushes(t, pushClient)
	var push *pushproto.Push
	for _, p := range pushes {
		if p.FcmApiKey == fcmAPIKey {
			push = p
			break
		}
	}
	if push == nil {
		t.Fatalf("Push not found")
	}
	if push.FcmApiKey != fcmAPIKey {
		t.Fatalf("Incorrect FcmApiKey. Expected: %s actual: %s", fcmAPIKey, push.FcmApiKey)
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
		Command:              cmd,
		EnvironmentNamespace: *environmentNamespace,
	}
	if _, err := client.CreateFeature(ctx, createReq); err != nil {
		t.Fatal(err)
	}
	enableFeature(t, featureID, client)
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

func createPush(ctx context.Context, t *testing.T, client pushclient.Client, fcmAPIKey, tag string) {
	t.Helper()
	cmd := newCreatePushCommand(t, fcmAPIKey, []string{tag})
	createReq := &pushproto.CreatePushRequest{
		EnvironmentNamespace: *environmentNamespace,
		Command:              cmd,
	}
	if _, err := client.CreatePush(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func newCreatePushCommand(t *testing.T, fcmAPIKey string, tags []string) *pushproto.CreatePushCommand {
	return &pushproto.CreatePushCommand{
		Name:      newPushName(t),
		FcmApiKey: fcmAPIKey,
		Tags:      tags,
	}
}

func listPushes(t *testing.T, client pushclient.Client) []*pushproto.Push {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListPushes(ctx, &pushproto.ListPushesRequest{
		EnvironmentNamespace: *environmentNamespace,
		PageSize:             int64(500),
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
