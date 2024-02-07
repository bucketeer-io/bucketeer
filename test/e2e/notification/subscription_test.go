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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	rpcclient "github.com/bucketeer-io/bucketeer/pkg/rpc/client"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
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

func TestCreateGetDeleteSubscription(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	notificationClient := newNotificationClient(t)
	defer notificationClient.Close()

	name := fmt.Sprintf("%s-name-%s", prefixTestName, newUUID(t))
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_DOMAIN_EVENT_ACCOUNT,
	}
	webhookURL := fmt.Sprintf("%s-webhook-url-%s", prefixTestName, newUUID(t))
	recipient := &proto.Recipient{
		Type:                  proto.Recipient_SlackChannel,
		SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: webhookURL},
	}
	id, err := domain.ID(recipient)
	if err != nil {
		t.Fatal(err)
	}
	createSubscription(ctx, t, notificationClient, name, sourceTypes, recipient)
	resp, err := notificationClient.GetSubscription(ctx, &proto.GetSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		t.Fatal(err)
	}
	subscription := resp.Subscription
	if subscription == nil {
		t.Fatalf("Subscription not found")
	}
	if subscription.Name != name {
		t.Fatalf("Incorrect name. Expected: %s actual: %s", name, subscription.Name)
	}
	if len(subscription.SourceTypes) != 1 {
		t.Fatalf("The number of notification types is incorrect. Expected: %d actual: %d", 1, len(subscription.SourceTypes))
	}
	if subscription.SourceTypes[0] != sourceTypes[0] {
		t.Fatalf("Incorrect notification type. Expected: %s actual: %s", sourceTypes[0], subscription.SourceTypes[0])
	}
	if subscription.Recipient.Type != proto.Recipient_SlackChannel {
		t.Fatalf("Incorrect recipient type. Expected: %s actual: %s", proto.Recipient_SlackChannel, subscription.Recipient.Type)
	}
	if subscription.Recipient.SlackChannelRecipient.WebhookUrl != webhookURL {
		t.Fatalf("Incorrect webhook URL. Expected: %s actual: %s", webhookURL, subscription.Recipient.SlackChannelRecipient.WebhookUrl)
	}
	if subscription.Disabled != false {
		t.Fatalf("Incorrect deleted. Expected: %t actual: %t", false, subscription.Disabled)
	}
	_, err = notificationClient.DeleteSubscription(ctx, &proto.DeleteSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command:              &proto.DeleteSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = notificationClient.GetSubscription(ctx, &proto.GetSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Fatal(err)
		}
	}
}

func TestCreateListDeleteSubscription(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	notificationClient := newNotificationClient(t)
	defer notificationClient.Close()

	name := fmt.Sprintf("%s-name-%s", prefixTestName, newUUID(t))
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_DOMAIN_EVENT_ACCOUNT,
	}
	webhookURL := fmt.Sprintf("%s-webhook-url-%s", prefixTestName, newUUID(t))
	recipient := &proto.Recipient{
		Type:                  proto.Recipient_SlackChannel,
		SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: webhookURL},
	}
	id, err := domain.ID(recipient)
	if err != nil {
		t.Fatal(err)
	}
	createSubscription(ctx, t, notificationClient, name, sourceTypes, recipient)
	subscriptions := listSubscriptions(t, notificationClient, []proto.Subscription_SourceType{proto.Subscription_DOMAIN_EVENT_ACCOUNT})
	var subscription *proto.Subscription
	for _, s := range subscriptions {
		if s.Id == id {
			subscription = s
			break
		}
	}
	if subscription == nil {
		t.Fatalf("Subscription not found")
	}
	if subscription.Name != name {
		t.Fatalf("Incorrect name. Expected: %s actual: %s", name, subscription.Name)
	}
	if len(subscription.SourceTypes) != 1 {
		t.Fatalf("The number of notification types is incorrect. Expected: %d actual: %d", 1, len(subscription.SourceTypes))
	}
	if subscription.SourceTypes[0] != sourceTypes[0] {
		t.Fatalf("Incorrect notification type. Expected: %s actual: %s", sourceTypes[0], subscription.SourceTypes[0])
	}
	if subscription.Recipient.Type != proto.Recipient_SlackChannel {
		t.Fatalf("Incorrect recipient type. Expected: %s actual: %s", proto.Recipient_SlackChannel, subscription.Recipient.Type)
	}
	if subscription.Recipient.SlackChannelRecipient.WebhookUrl != webhookURL {
		t.Fatalf("Incorrect webhook URL. Expected: %s actual: %s", webhookURL, subscription.Recipient.SlackChannelRecipient.WebhookUrl)
	}
	if subscription.Disabled != false {
		t.Fatalf("Incorrect deleted. Expected: %t actual: %t", false, subscription.Disabled)
	}
	_, err = notificationClient.DeleteSubscription(ctx, &proto.DeleteSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command:              &proto.DeleteSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = notificationClient.GetSubscription(ctx, &proto.GetSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Fatal(err)
		}
	}
}

func TestUpdateSubscription(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	notificationClient := newNotificationClient(t)
	defer notificationClient.Close()

	name := fmt.Sprintf("%s-name-%s", prefixTestName, newUUID(t))
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_DOMAIN_EVENT_ACCOUNT,
	}
	webhookURL := fmt.Sprintf("%s-webhook-url-%s", prefixTestName, newUUID(t))
	recipient := &proto.Recipient{
		Type:                  proto.Recipient_SlackChannel,
		SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: webhookURL},
	}
	id, err := domain.ID(recipient)
	if err != nil {
		t.Fatal(err)
	}
	createSubscription(ctx, t, notificationClient, name, sourceTypes, recipient)
	_, err = notificationClient.UpdateSubscription(ctx, &proto.UpdateSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		AddSourceTypesCommand: &proto.AddSourceTypesCommand{
			SourceTypes: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
			},
		},
		DeleteSourceTypesCommand: &proto.DeleteSourceTypesCommand{
			SourceTypes: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_ACCOUNT,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := notificationClient.GetSubscription(ctx, &proto.GetSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		t.Fatal(err)
	}
	subscription := resp.Subscription
	if subscription == nil {
		t.Fatalf("Subscription not found")
	}
	if subscription.Name != name {
		t.Fatalf("Incorrect name. Expected: %s actual: %s", name, subscription.Name)
	}
	if len(subscription.SourceTypes) != 1 {
		t.Fatalf("The number of notification types is incorrect. Expected: %d actual: %d", 1, len(subscription.SourceTypes))
	}
	if subscription.SourceTypes[0] != proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT {
		t.Fatalf("Incorrect notification type. Expected: %s actual: %s", sourceTypes[0], subscription.SourceTypes[0])
	}
	if subscription.Recipient.Type != proto.Recipient_SlackChannel {
		t.Fatalf("Incorrect recipient type. Expected: %s actual: %s", proto.Recipient_SlackChannel, subscription.Recipient.Type)
	}
	if subscription.Recipient.SlackChannelRecipient.WebhookUrl != webhookURL {
		t.Fatalf("Incorrect webhook URL. Expected: %s actual: %s", webhookURL, subscription.Recipient.SlackChannelRecipient.WebhookUrl)
	}
	if subscription.Disabled != false {
		t.Fatalf("Incorrect deleted. Expected: %t actual: %t", false, subscription.Disabled)
	}
	_, err = notificationClient.DeleteSubscription(ctx, &proto.DeleteSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command:              &proto.DeleteSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = notificationClient.GetSubscription(ctx, &proto.GetSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Fatal(err)
		}
	}
}

func TestListEnabledSubscriptions(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	notificationClient := newNotificationClient(t)
	defer notificationClient.Close()

	name := fmt.Sprintf("%s-name-%s", prefixTestName, newUUID(t))
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_MAU_COUNT,
	}
	webhookURL := fmt.Sprintf("%s-webhook-url-%s", prefixTestName, newUUID(t))
	recipient := &proto.Recipient{
		Type:                  proto.Recipient_SlackChannel,
		SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: webhookURL},
	}
	id, err := domain.ID(recipient)
	if err != nil {
		t.Fatal(err)
	}
	createSubscription(ctx, t, notificationClient, name, sourceTypes, recipient)
	_, err = notificationClient.DisableSubscription(ctx, &proto.DisableSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command:              &proto.DisableSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	subscriptions := listEnabledSubscriptions(
		t,
		notificationClient,
		[]proto.Subscription_SourceType{proto.Subscription_MAU_COUNT},
	)
	var contains bool
	for _, s := range subscriptions {
		if s.Id == id {
			contains = true
			break
		}
	}
	if contains {
		t.Fatal("List enabled subscriptions include disabled subscription")
	}
	_, err = notificationClient.DeleteSubscription(ctx, &proto.DeleteSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Id:                   id,
		Command:              &proto.DeleteSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
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

func newUUID(t *testing.T) string {
	t.Helper()
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

func createSubscription(
	ctx context.Context,
	t *testing.T,
	client notificationclient.Client,
	name string,
	sourceTypes []proto.Subscription_SourceType,
	recipient *proto.Recipient) {

	t.Helper()
	cmd := newCreateSubscriptionCommand(name, sourceTypes, recipient)
	createReq := &proto.CreateSubscriptionRequest{
		EnvironmentNamespace: *environmentNamespace,
		Command:              cmd,
	}
	if _, err := client.CreateSubscription(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func newCreateSubscriptionCommand(
	name string,
	sourceTypes []proto.Subscription_SourceType,
	recipient *proto.Recipient) *proto.CreateSubscriptionCommand {

	return &proto.CreateSubscriptionCommand{
		Name:        name,
		SourceTypes: sourceTypes,
		Recipient:   recipient,
	}
}

func listSubscriptions(
	t *testing.T,
	client notificationclient.Client,
	sourceTypes []proto.Subscription_SourceType) []*proto.Subscription {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListSubscriptions(ctx, &proto.ListSubscriptionsRequest{
		EnvironmentNamespace: *environmentNamespace,
		PageSize:             int64(500),
		SourceTypes:          sourceTypes,
	})
	if err != nil {
		t.Fatal("failed to list subscriptions", err)
	}
	return resp.Subscriptions
}

func listEnabledSubscriptions(
	t *testing.T,
	client notificationclient.Client,
	sourceTypes []proto.Subscription_SourceType) []*proto.Subscription {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListEnabledSubscriptions(ctx, &proto.ListEnabledSubscriptionsRequest{
		EnvironmentNamespace: *environmentNamespace,
		PageSize:             int64(500),
		SourceTypes:          sourceTypes,
	})
	if err != nil {
		t.Fatal("failed to list enabled subscriptions", err)
	}
	return resp.Subscriptions
}
