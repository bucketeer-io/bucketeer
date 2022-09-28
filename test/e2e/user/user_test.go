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

package feature

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"

	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	gatewayclient "github.com/bucketeer-io/bucketeer/pkg/gateway/client"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	userclient "github.com/bucketeer-io/bucketeer/pkg/user/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
	"github.com/bucketeer-io/bucketeer/test/util"
)

const (
	prefixTestName = "e2e-test"
	retryTimes     = 60
	timeout        = 10 * time.Second
)

var (
	// FIXME: To avoid compiling the test many times, webGatewayAddr, webGatewayPort & apiKey has been also added here to prevent from getting: "flag provided but not defined" error during the test. These 3 are being use  in the Gateway test
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

func TestGetUser(t *testing.T) {
	t.Parallel()
	uuid := newUUID(t)
	tag := fmt.Sprintf("%s-tag-%s", prefixTestName, uuid)
	userID := newUserID(t, uuid)
	now := time.Now()
	featureID := newFeatureID(t, uuid)
	featureclient := newFeatureClient(t)
	defer featureclient.Close()
	userClient := newUserClient(t)
	defer userClient.Close()
	feature := createFeatureWithTag(t, featureclient, featureID, tag)
	time.Sleep(3 * time.Second)
	user := &userproto.User{
		Id: userID,
	}

	// Check evaluations
	for i := 0; i < retryTimes; i++ {
		resp := getEvaluations(t, tag, user)
		if resp.State == featureproto.UserEvaluations_FULL {
			evaluationsSize := len(resp.Evaluations.Evaluations)
			if evaluationsSize != 1 {
				t.Fatalf("The number is evaluations is not correct. Expected: 0, actual: %d", evaluationsSize)
			}
			if resp.Evaluations == nil {
				t.Fatal("Evaluations field is nil")
			}
			variationID := resp.Evaluations.Evaluations[0].Variation.Id
			if feature.Variations[0].Id != variationID {
				t.Fatalf("Variation doesn't match. Expected: %s, actual: %s", feature.Variations[0].Id, variationID)
			}
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("State did not change. Expected: %v, actual: %v", featureproto.UserEvaluations_FULL, resp.State)
		}
		time.Sleep(time.Second)
	}

	// Check user
	var latestSeen int64
	for i := 0; i < retryTimes; i++ {
		actual, _ := getUser(t, userClient, userID)
		if actual != nil {
			if actual.Id != userID {
				t.Fatalf("User ID is not correct: expected: %s, actual: %s", userID, actual.Id)
			}
			if len(actual.TaggedData[tag].Value) != 0 {
				t.Fatalf("The user metadata should be zero. Actual: %v", actual.TaggedData[tag].Value)
			}
			if actual.LastSeen < now.Unix() {
				t.Fatalf("Last seen is not correct: expected: %d, actual: %d", now.Unix(), actual.LastSeen)
			}
			latestSeen = actual.LastSeen
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("User cannot be fetched.")
		}
		time.Sleep(time.Second)
	}

	// Change user attribute
	user.Data = map[string]string{"k-0": "v-0"}

	// Check evaluations change
	for i := 0; i < retryTimes; i++ {
		resp := getEvaluations(t, tag, user)
		if resp.State == featureproto.UserEvaluations_FULL {
			evaluationsSize := len(resp.Evaluations.Evaluations)
			if evaluationsSize != 1 {
				t.Fatalf("The evaluations size is not correct. Expected: 1, actual: %d. Data: %v", evaluationsSize, resp.Evaluations.Evaluations)
			}
			if resp.Evaluations == nil {
				t.Fatal("Evaluations field is nil")
			}
			variationID := resp.Evaluations.Evaluations[0].VariationId
			if feature.Variations[1].Id == variationID {
				break
			}
		}
		if i == retryTimes-1 {
			t.Fatalf("Evaluations did not change. Variation Expected: %s, actual: %s", feature.Variations[1].Id, resp.Evaluations.Evaluations[0].VariationId)
		}
		time.Sleep(time.Second)
	}

	// Check user changes
	for i := 0; i < retryTimes; i++ {
		actual, _ := getUser(t, userClient, userID)
		if actual != nil {
			if actual.Id != userID {
				t.Fatalf("User ID is not correct: expected: %s, actual: %s", userID, actual.Id)
			}
			if len(actual.TaggedData[tag].Value) == 1 {
				for _, data := range actual.TaggedData {
					for k, v := range data.Value {
						if k != "k-0" {
							t.Fatalf("Data key is different. Expected: %s, Actual: %s", k, "k-0")
						}
						if v != "v-0" {
							t.Fatalf("Data value is different. Expected: %s, Actual: %s", v, "v-0")
						}
						if actual.LastSeen < latestSeen {
							t.Fatalf("Last seen is not correct: expected: %d, actual: %d", now.Unix(), actual.LastSeen)
						}
					}
				}
				break
			}
		}
		if i == retryTimes-1 {
			t.Fatalf("User did not change when adding the user data.")
		}
		time.Sleep(time.Second)
	}

	// Use different tag
	tagServer := fmt.Sprintf("%s-tag-server-%s", prefixTestName, uuid)
	resp := getEvaluations(t, tagServer, user)
	assert.NotNil(t, resp)
	time.Sleep(time.Second)

	// Check if user's data has changed
	for i := 0; i < retryTimes; i++ {
		actual, _ := getUser(t, userClient, userID)
		if actual != nil {
			if actual.Id != userID {
				t.Fatalf("User ID is not correct: expected: %s, actual: %s", userID, actual.Id)
			}
			// At this point it has one tagged data waiting for the second one be persisted
			if len(actual.TaggedData) == 1 {
				continue
			}
			if len(actual.TaggedData[tagServer].Value) != 1 {
				t.Fatalf("User data size should not be different than one at this point. Actual: %d. Data: %v",
					len(actual.TaggedData[tagServer].Value),
					actual.TaggedData[tagServer].Value,
				)
			}
			for _, data := range actual.TaggedData {
				for k, v := range data.Value {
					if k != "k-0" {
						t.Fatalf("Data key is different. Expected: %s, Actual: %s", k, "k-0")
					}
					if v != "v-0" {
						t.Fatalf("Data value is different. Expected: %s, Actual: %s", v, "v-0")
					}
					if actual.LastSeen < latestSeen {
						t.Fatalf("Last seen is not correct: expected: %d, actual: %d", now.Unix(), actual.LastSeen)
					}
				}
			}
			break
		}
		if i == retryTimes-1 {
			t.Fatalf("User did not change when using different tags.")
		}
		time.Sleep(time.Second)
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

func createFeatureWithTag(t *testing.T, client featureclient.Client, featureID, tag string) *featureproto.Feature {
	cmd := newCreateFeatureCommand(featureID, "a", "b", []string{tag})
	createFeature(t, client, cmd)
	f := getFeature(t, featureID, client)
	rule := newFixedStrategyRule(f.Variations[1].Id, "k-0", "v-0")
	addCmd, _ := util.MarshalCommand(&featureproto.AddRuleCommand{Rule: rule})
	updateFeatureTargeting(t, client, addCmd, featureID)
	enableFeature(t, featureID, client)
	return getFeature(t, featureID, client)
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

func newCreateFeatureCommand(featureID string, varA, varB string, tags []string) *featureproto.CreateFeatureCommand {
	return &featureproto.CreateFeatureCommand{
		Id:          featureID,
		Name:        "e2e-test-gateway-feature-name",
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
		Tags:                     tags,
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

func newUserClient(t *testing.T) userclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	userClient, err := userclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create user client:", err)
	}
	return userClient
}

func getUser(t *testing.T, client userclient.Client, userID string) (*userproto.User, error) {
	t.Helper()
	req := &userproto.GetUserRequest{
		UserId:               userID,
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func getEvaluations(t *testing.T, tag string, user *userproto.User) *gatewayproto.GetEvaluationsResponse {
	t.Helper()
	c := newGatewayClient(t)
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := &gatewayproto.GetEvaluationsRequest{
		Tag:  tag,
		User: user,
	}
	response, err := c.GetEvaluations(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func newFixedStrategyRule(variationID string, attr string, value string) *featureproto.Rule {
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
				Attribute: attr,
				Operator:  featureproto.Clause_EQUALS,
				Values:    []string{value},
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
		EnvironmentNamespace: *environmentNamespace,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if _, err := client.UpdateFeatureTargeting(ctx, updateReq); err != nil {
		t.Fatal(err)
	}
}

func getFeature(t *testing.T, featureID string, client featureclient.Client) *featureproto.Feature {
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

func newFeatureID(t *testing.T, uuid string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-feature-id-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-feature-id-%s", prefixTestName, uuid)
}

func newUserID(t *testing.T, uuid string) string {
	if *testID != "" {
		return fmt.Sprintf("%s-%s-user-%s", prefixTestName, *testID, uuid)
	}
	return fmt.Sprintf("%s-user-%s", prefixTestName, uuid)
}
