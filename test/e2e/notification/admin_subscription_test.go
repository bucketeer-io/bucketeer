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
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notificationclient "github.com/bucketeer-io/bucketeer/pkg/notification/client"
	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func TestCreateGetDeleteAdminSubscription(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	notificationClient := newNotificationClient(t)
	defer notificationClient.Close()

	name := fmt.Sprintf("%s-name-%s", prefixTestName, newUUID(t))
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
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
	createAdminSubscription(ctx, t, notificationClient, name, sourceTypes, recipient)
	resp, err := notificationClient.GetAdminSubscription(ctx, &proto.GetAdminSubscriptionRequest{
		Id: id,
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
	_, err = notificationClient.DeleteAdminSubscription(ctx, &proto.DeleteAdminSubscriptionRequest{
		Id:      id,
		Command: &proto.DeleteAdminSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = notificationClient.GetAdminSubscription(ctx, &proto.GetAdminSubscriptionRequest{
		Id: id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Fatal(err)
		}
	}
}

func TestCreateListDeleteAdminSubscription(t *testing.T) {
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
	createAdminSubscription(ctx, t, notificationClient, name, sourceTypes, recipient)
	subscriptions := listAdminSubscriptions(t, notificationClient, []proto.Subscription_SourceType{proto.Subscription_DOMAIN_EVENT_ACCOUNT})
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
	_, err = notificationClient.DeleteAdminSubscription(ctx, &proto.DeleteAdminSubscriptionRequest{
		Id:      id,
		Command: &proto.DeleteAdminSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = notificationClient.GetAdminSubscription(ctx, &proto.GetAdminSubscriptionRequest{
		Id: id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Fatal(err)
		}
	}
}

func TestUpdateAdminSubscription(t *testing.T) {
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
	createAdminSubscription(ctx, t, notificationClient, name, sourceTypes, recipient)
	_, err = notificationClient.UpdateAdminSubscription(ctx, &proto.UpdateAdminSubscriptionRequest{
		Id: id,
		AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{
			SourceTypes: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
			},
		},
		DeleteSourceTypesCommand: &proto.DeleteAdminSubscriptionSourceTypesCommand{
			SourceTypes: []proto.Subscription_SourceType{
				proto.Subscription_DOMAIN_EVENT_ACCOUNT,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := notificationClient.GetAdminSubscription(ctx, &proto.GetAdminSubscriptionRequest{
		Id: id,
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
	_, err = notificationClient.DeleteAdminSubscription(ctx, &proto.DeleteAdminSubscriptionRequest{
		Id:      id,
		Command: &proto.DeleteAdminSubscriptionCommand{},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = notificationClient.GetSubscription(ctx, &proto.GetSubscriptionRequest{
		Id: id,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Fatal(err)
		}
	}
}

func createAdminSubscription(
	ctx context.Context,
	t *testing.T,
	client notificationclient.Client,
	name string,
	sourceTypes []proto.Subscription_SourceType,
	recipient *proto.Recipient) {

	t.Helper()
	cmd := newCreateAdminSubscriptionCommand(name, sourceTypes, recipient)
	createReq := &proto.CreateAdminSubscriptionRequest{
		Command: cmd,
	}
	if _, err := client.CreateAdminSubscription(ctx, createReq); err != nil {
		t.Fatal(err)
	}
}

func newCreateAdminSubscriptionCommand(
	name string,
	sourceTypes []proto.Subscription_SourceType,
	recipient *proto.Recipient) *proto.CreateAdminSubscriptionCommand {

	return &proto.CreateAdminSubscriptionCommand{
		Name:        name,
		SourceTypes: sourceTypes,
		Recipient:   recipient,
	}
}

func listAdminSubscriptions(
	t *testing.T,
	client notificationclient.Client,
	sourceTypes []proto.Subscription_SourceType) []*proto.Subscription {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := client.ListAdminSubscriptions(ctx, &proto.ListAdminSubscriptionsRequest{
		PageSize:    int64(500),
		SourceTypes: sourceTypes,
	})
	if err != nil {
		t.Fatal("failed to list subscriptions", err)
	}
	return resp.Subscriptions
}
