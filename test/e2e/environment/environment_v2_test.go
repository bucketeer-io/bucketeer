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

package environment

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	environmentclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	notificationclient "github.com/bucketeer-io/bucketeer/v2/pkg/notification/client"
	pushclient "github.com/bucketeer-io/bucketeer/v2/pkg/push/client"
	rpcclient "github.com/bucketeer-io/bucketeer/v2/pkg/rpc/client"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	notificationproto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

const (
	timeout = 60 * time.Second
)

var (
	// FIXME: To avoid compiling the test many times, webGatewayAddr, webGatewayPort & apiKey has been also added here to prevent from getting: "flag provided but not defined" error during the test. These 3 are being use in the Gateway test
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

const (
	environmentName = "E2E environment"
)

func TestGetEnvironmentV2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := getEnvironmentID(t)
	resp, err := c.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Environment.Id != id {
		t.Fatalf("different ids, expected: %v, actual: %v", id, resp.Environment.Id)
	}
	if resp.Environment.Name != environmentName {
		t.Fatalf("different name, expected: %v, actual: %v", environmentName, resp.Environment.Name)
	}
}

func TestListEnvironmentsV2ByProject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	resp, err := c.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{ProjectId: defaultProjectID})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Environments) == 0 {
		t.Fatal("environments is empty, expected at least 1")
	}
	for _, env := range resp.Environments {
		if env.ProjectId != defaultProjectID {
			t.Fatalf("different project id, expected: %s, actual: %s", defaultProjectID, env.ProjectId)
		}
	}
}

func TestListEnvironmentsV2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	pageSize := int64(1)
	resp, err := c.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{PageSize: pageSize})
	if err != nil {
		t.Fatal(err)
	}
	responseSize := int64(len(resp.Environments))
	if responseSize != pageSize {
		t.Fatalf("different sizes, expected: %d actual: %d", pageSize, responseSize)
	}
}

func TestUpdateEnvironmentV2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c := newEnvironmentClient(t)
	defer c.Close()
	id := getEnvironmentID(t)
	newDesc := fmt.Sprintf("This environment is for local development (Updated at %d)", time.Now().Unix())
	_, err := c.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
		Id:                       id,
		ChangeDescriptionCommand: &environmentproto.ChangeDescriptionEnvironmentV2Command{Description: newDesc},
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err := c.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp.Environment.Id != id {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp.Environment.Id)
	}
	if getResp.Environment.Name != environmentName {
		t.Fatalf("different name, expected: %v, actual: %v", environmentName, getResp.Environment.Name)
	}
	if getResp.Environment.Description != newDesc {
		t.Fatalf("different descriptions, expected: %v, actual: %v", newDesc, getResp.Environment.Description)
	}

	newDesc = fmt.Sprintf("This environment is for local development (Updated at %d with no command)", time.Now().Unix())
	_, err = c.UpdateEnvironmentV2(ctx, &environmentproto.UpdateEnvironmentV2Request{
		Id:          id,
		Description: wrapperspb.String(newDesc),
	})
	if err != nil {
		t.Fatal(err)
	}
	getResp, err = c.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{Id: id})
	if err != nil {
		t.Fatal(err)
	}
	if getResp.Environment.Id != id {
		t.Fatalf("different ids, expected: %v, actual: %v", id, getResp.Environment.Id)
	}
	if getResp.Environment.Name != environmentName {
		t.Fatalf("different name, expected: %v, actual: %v", environmentName, getResp.Environment.Name)
	}
	if getResp.Environment.Description != newDesc {
		t.Fatalf("different descriptions, expected: %v, actual: %v", newDesc, getResp.Environment.Description)
	}
}

func TestCreateDeleteEnvironmentV2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	envc := newEnvironmentClient(t)
	ftc := newFeatureClient(t)
	expc := newExperimentClient(t)
	pushc := newPushClient(t)
	notic := newNotificationClient(t)
	opsc := newAutoOpsClient(t)
	accountc := newAccountClient(t)
	defer envc.Close()
	defer ftc.Close()
	defer expc.Close()
	defer pushc.Close()
	defer notic.Close()
	defer opsc.Close()
	defer accountc.Close()

	// 1. create env
	createEnvResp, err := envc.CreateEnvironmentV2(ctx, &environmentproto.CreateEnvironmentV2Request{
		Name:           fmt.Sprintf("%s-%d", environmentName, time.Now().UnixNano()),
		UrlCode:        fmt.Sprintf("env-url-%d", time.Now().UnixNano()),
		ProjectId:      defaultProjectID,
		RequireComment: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if createEnvResp == nil || createEnvResp.Environment == nil {
		t.Fatal("CreateEnvironmentV2 returned nil response")
	}
	envID := createEnvResp.Environment.Id

	// 2. create data for the new environment
	// 2.1 create feature
	createFfResp, err := ftc.CreateFeature(ctx, newCreateFeatureReq(
		fmt.Sprintf("feature-e2e-%d", time.Now().UnixNano()),
		envID,
	))
	if err != nil {
		t.Fatal(err)
	}
	if createFfResp == nil || createFfResp.Feature == nil {
		t.Fatal("create feature response or feature is nil")
	}

	// 2.2 create goals
	createGoalResp, err := expc.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		EnvironmentId:  envID,
		Id:             fmt.Sprintf("goal-id-e2e-%d", time.Now().UnixNano()),
		Name:           fmt.Sprintf("goal-name-e2e-%d", time.Now().UnixNano()),
		ConnectionType: experimentproto.Goal_EXPERIMENT,
	})
	if err != nil {
		t.Fatal(err)
	}
	if createGoalResp == nil || createGoalResp.Goal == nil {
		t.Fatal("create goal response or goal is nil")
	}

	// 2.3 create push
	createPush, err := pushc.CreatePush(ctx, &pushproto.CreatePushRequest{
		EnvironmentId:     envID,
		Name:              fmt.Sprintf("push-name-e2e-%d", time.Now().UnixNano()),
		FcmServiceAccount: []byte(fmt.Sprintf(fcmServiceAccountDummy, time.Now().UnixNano())),
		Tags:              []string{"e2e-test-tag-1", "e2e-test-tag-2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if createPush == nil || createPush.Push == nil {
		t.Fatal("create push response or push is nil")
	}

	// 2.4 create subscription
	createSubscription, err := notic.CreateSubscription(ctx, &notificationproto.CreateSubscriptionRequest{
		EnvironmentId: envID,
		Name:          fmt.Sprintf("subscription-name-e2e-%d", time.Now().UnixNano()),
		SourceTypes: []notificationproto.Subscription_SourceType{
			notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
			notificationproto.Subscription_DOMAIN_EVENT_APIKEY,
		},
		Recipient: &notificationproto.Recipient{
			Type: notificationproto.Recipient_SlackChannel,
			SlackChannelRecipient: &notificationproto.SlackChannelRecipient{
				WebhookUrl: "https://hooks.slack",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if createSubscription == nil || createSubscription.Subscription == nil {
		t.Fatal("create subscription response or subscription is nil")
	}

	// 2.5 create operation
	createOpsResp, err := opsc.CreateAutoOpsRule(ctx, &autoopsproto.CreateAutoOpsRuleRequest{
		EnvironmentId: envID,
		FeatureId:     createFfResp.Feature.Id,
		OpsType:       autoopsproto.OpsType_SCHEDULE,
		DatetimeClauses: []*autoopsproto.DatetimeClause{
			{
				Time:       time.Now().Add(2 * time.Hour).Unix(),
				ActionType: autoopsproto.ActionType_ENABLE,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if createOpsResp == nil || createOpsResp.AutoOpsRule == nil {
		t.Fatal("create auto ops rule response or rule is nil")
	}

	// 2.6 create experiment
	createExpResp, err := expc.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{
		EnvironmentId:   envID,
		FeatureId:       createFfResp.Feature.Id,
		Name:            fmt.Sprintf("experiment-name-e2e-%d", time.Now().UnixNano()),
		Description:     "e2e test experiment description",
		GoalIds:         []string{createGoalResp.Goal.Id},
		StartAt:         time.Now().Unix(),
		StopAt:          time.Now().Add(24 * time.Hour).Unix(),
		BaseVariationId: createFfResp.Feature.Variations[0].Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if createExpResp == nil || createExpResp.Experiment == nil {
		t.Fatal("create experiment response or experiment is nil")
	}

	// 2.7 create API Key
	createAPIKeyResp, err := accountc.CreateAPIKey(ctx, &accountproto.CreateAPIKeyRequest{
		EnvironmentId: envID,
		Name:          fmt.Sprintf("api-key-name-e2e-%d", time.Now().UnixNano()),
		Role:          accountproto.APIKey_PUBLIC_API_WRITE,
	})
	if err != nil {
		t.Fatal(err)
	}
	if createAPIKeyResp == nil || createAPIKeyResp.ApiKey == nil {
		t.Fatal("create API key response or API key is nil")
	}

	// 3.0 dry-run delete env
	_, err = envc.DeleteEnvironmentData(ctx, &environmentproto.DeleteEnvironmentDataRequest{
		EnvironmentIds: []string{envID},
		DryRun:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 3.1 verify feature still exists
	getFeatureResp, err := ftc.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: envID,
		Id:            createFfResp.Feature.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getFeatureResp == nil || getFeatureResp.Feature == nil {
		t.Fatal("get feature response or feature is nil after dry-run delete")
	}

	// 3.2 verify goal still exists
	getGoalResp, err := expc.GetGoal(ctx, &experimentproto.GetGoalRequest{
		EnvironmentId: envID,
		Id:            createGoalResp.Goal.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getGoalResp == nil || getGoalResp.Goal == nil {
		t.Fatal("get goal response or goal is nil after dry-run delete")
	}

	// 3.3 verify push still exists
	getPushResp, err := pushc.GetPush(ctx, &pushproto.GetPushRequest{
		EnvironmentId: envID,
		Id:            createPush.Push.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getPushResp == nil || getPushResp.Push == nil {
		t.Fatal("get push response or push is nil after dry-run delete")
	}

	// 3.4 verify subscription still exists
	getSubscriptionResp, err := notic.GetSubscription(ctx, &notificationproto.GetSubscriptionRequest{
		EnvironmentId: envID,
		Id:            createSubscription.Subscription.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getSubscriptionResp == nil || getSubscriptionResp.Subscription == nil {
		t.Fatal("get subscription response or subscription is nil after dry-run delete")
	}

	// 3.5 verify auto ops rule still exists
	getAutoOpsRuleResp, err := opsc.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: envID,
		Id:            createOpsResp.AutoOpsRule.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getAutoOpsRuleResp == nil || getAutoOpsRuleResp.AutoOpsRule == nil {
		t.Fatal("get auto ops rule response or rule is nil after dry-run delete")
	}

	// 3.6 verify experiment still exists
	getExperimentResp, err := expc.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		EnvironmentId: envID,
		Id:            createExpResp.Experiment.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getExperimentResp == nil || getExperimentResp.Experiment == nil {
		t.Fatal("get experiment response or experiment is nil after dry-run delete")
	}

	// 3.7 verify API key still exists
	getAPIKeyResp, err := accountc.GetAPIKey(ctx, &accountproto.GetAPIKeyRequest{
		EnvironmentId: envID,
		Id:            createAPIKeyResp.ApiKey.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	if getAPIKeyResp == nil || getAPIKeyResp.ApiKey == nil {
		t.Fatal("get API key response or API key is nil after dry-run delete")
	}

	// 4. delete env
	_, err = envc.DeleteEnvironmentData(ctx, &environmentproto.DeleteEnvironmentDataRequest{
		EnvironmentIds: []string{envID},
	})
	if err != nil {
		t.Fatal(err)
	}

	// 5.0 verify data is deleted
	// 5.1 verify feature is deleted
	_, err = ftc.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: envID,
		Id:            createFfResp.Feature.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting feature from deleted environment, but got nil")
	}
	// 5.2 verify goal is deleted
	_, err = expc.GetGoal(ctx, &experimentproto.GetGoalRequest{
		EnvironmentId: envID,
		Id:            createGoalResp.Goal.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting goal from deleted organization, but got nil")
	}

	// 5.3 verify push is deleted
	_, err = pushc.GetPush(ctx, &pushproto.GetPushRequest{
		EnvironmentId: envID,
		Id:            createPush.Push.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting push from deleted organization, but got nil")
	}

	// 5.4 verify subscription is deleted
	_, err = notic.GetSubscription(ctx, &notificationproto.GetSubscriptionRequest{
		EnvironmentId: envID,
		Id:            createSubscription.Subscription.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting subscription from deleted organization, but got nil")
	}

	// 5.5 verify auto ops rule is deleted
	_, err = opsc.GetAutoOpsRule(ctx, &autoopsproto.GetAutoOpsRuleRequest{
		EnvironmentId: envID,
		Id:            createOpsResp.AutoOpsRule.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting auto ops rule from deleted organization, but got nil")
	}

	// 5.6 verify experiment is deleted
	_, err = expc.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		EnvironmentId: envID,
		Id:            createExpResp.Experiment.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting experiment from deleted organization, but got nil")
	}

	// 5.7 verify API key is deleted
	_, err = accountc.GetAPIKey(ctx, &accountproto.GetAPIKeyRequest{
		EnvironmentId: envID,
		Id:            createAPIKeyResp.ApiKey.Id,
	})
	if err == nil {
		t.Fatal("expected error when getting API key from deleted organization, but got nil")
	}
}

func getEnvironmentID(t *testing.T) string {
	t.Helper()
	if *environmentID == "" {
		return "production"
	}
	return *environmentID
}

func newEnvironmentClient(t *testing.T) environmentclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := environmentclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create environment client:", err)
	}
	return client
}

func newFeatureClient(t *testing.T) featureclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := featureclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create feature client:", err)
	}
	return client
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

func newNotificationClient(t *testing.T) notificationclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := notificationclient.NewClient(
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

func newAutoOpsClient(t *testing.T) autoopsclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := autoopsclient.NewClient(
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

func newAccountClient(t *testing.T) accountclient.Client {
	t.Helper()
	creds, err := rpcclient.NewPerRPCCredentials(*serviceTokenPath)
	if err != nil {
		t.Fatal("Failed to create RPC credentials:", err)
	}
	client, err := accountclient.NewClient(
		fmt.Sprintf("%s:%d", *webGatewayAddr, *webGatewayPort),
		*webGatewayCert,
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithDialTimeout(30*time.Second),
		rpcclient.WithBlock(),
	)
	if err != nil {
		t.Fatal("Failed to create environment client:", err)
	}
	return client
}
