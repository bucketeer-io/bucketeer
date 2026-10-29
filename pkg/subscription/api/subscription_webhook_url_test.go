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

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

const maskedSlackWebhookURL = "https://hooks.slack.com/****"

func storedSubscription() *domain.Subscription {
	return &domain.Subscription{Subscription: &proto.Subscription{
		Id:            "sub-0",
		Name:          "sname",
		EnvironmentId: "ns0",
		SourceTypes:   []proto.Subscription_SourceType{proto.Subscription_DOMAIN_EVENT_ACCOUNT},
		Recipient: &proto.Recipient{
			Type:                  proto.Recipient_SlackChannel,
			SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: slackWebhookURL},
		},
	}}
}

func newIncomingContext(t *testing.T) context.Context {
	t.Helper()
	return metadata.NewIncomingContext(context.Background(), metadata.MD{
		"accept-language": []string{"ja"},
	})
}

func setServiceToken(t *testing.T, ctx context.Context) context.Context {
	t.Helper()
	return context.WithValue(ctx, rpc.AccessTokenKey, &token.AccessToken{
		Issuer:         "issuer",
		Audience:       "audience",
		Expiry:         time.Now().AddDate(100, 0, 0),
		IssuedAt:       time.Now(),
		Email:          "service@bucketeer.io",
		IsSystemAdmin:  true,
		IsServiceToken: true,
	})
}

func runTransaction(s *SubscriptionService) {
	s.dbClient.(*dbmock.MockClient).EXPECT().RunInTransactionV2(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	})
}

func TestGetSubscriptionMasksWebhookURL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	service := newSubscriptionService(
		mockController,
		nil,
		toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
		toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
	)
	service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
		gomock.Any(), "sub-0", "ns0",
	).Return(storedSubscription(), nil)

	ctx := setToken(t, newIncomingContext(t), false)
	actual, err := service.GetSubscription(ctx, &proto.GetSubscriptionRequest{Id: "sub-0", EnvironmentId: "ns0"})
	require.NoError(t, err)
	assert.Equal(t, maskedSlackWebhookURL, actual.Subscription.Recipient.SlackChannelRecipient.WebhookUrl)
	assert.Equal(t, "sub-0", actual.Subscription.Id)
	assert.Equal(t, "sname", actual.Subscription.Name)
}

func TestListSubscriptionsMasksWebhookURL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	service := newSubscriptionService(
		mockController,
		nil,
		toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
		toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
	)
	service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().ListSubscriptions(
		gomock.Any(), gomock.Any(),
	).Return([]*proto.Subscription{
		storedSubscription().Subscription,
		{Id: "sub-1", Name: "without recipient"},
	}, 0, int64(2), nil)

	ctx := setToken(t, newIncomingContext(t), false)
	actual, err := service.ListSubscriptions(ctx, &proto.ListSubscriptionsRequest{EnvironmentId: "ns0", PageSize: 10})
	require.NoError(t, err)
	require.Len(t, actual.Subscriptions, 2)
	assert.Equal(t, maskedSlackWebhookURL, actual.Subscriptions[0].Recipient.SlackChannelRecipient.WebhookUrl)
	assert.Equal(t, "sub-1", actual.Subscriptions[1].Id)
	assert.Nil(t, actual.Subscriptions[1].Recipient)
}

func TestListEnabledSubscriptionsWebhookURL(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		setToken func(*testing.T, context.Context) context.Context
		expected string
	}{
		{
			desc: "environment viewer: masked",
			setToken: func(t *testing.T, ctx context.Context) context.Context {
				return setToken(t, ctx, false)
			},
			expected: maskedSlackWebhookURL,
		},
		{
			desc: "system admin: masked",
			setToken: func(t *testing.T, ctx context.Context) context.Context {
				return setToken(t, ctx, true)
			},
			expected: maskedSlackWebhookURL,
		},
		{
			desc:     "service token used by the notification sender: original URL",
			setToken: setServiceToken,
			expected: slackWebhookURL,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := newSubscriptionService(
				mockController,
				nil,
				toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
				toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			)
			service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().ListSubscriptions(
				gomock.Any(), gomock.Any(),
			).Return([]*proto.Subscription{storedSubscription().Subscription}, 0, int64(1), nil)

			ctx := p.setToken(t, newIncomingContext(t))
			actual, err := service.ListEnabledSubscriptions(ctx, &proto.ListEnabledSubscriptionsRequest{
				EnvironmentId: "ns0",
				PageSize:      10,
			})
			require.NoError(t, err)
			require.Len(t, actual.Subscriptions, 1)
			assert.Equal(t, p.expected, actual.Subscriptions[0].Recipient.SlackChannelRecipient.WebhookUrl)
		})
	}
}

func TestCreateSubscriptionStoresWebhookURLButMasksEventAndResponse(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	service := newSubscriptionServiceWithMock(t, mockController)
	runTransaction(service)
	var stored *domain.Subscription
	service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().CreateSubscription(
		gomock.Any(), gomock.Any(), "ns0",
	).DoAndReturn(func(_ context.Context, s *domain.Subscription, _ string) error {
		stored = s
		return nil
	})
	var published *eventproto.Event
	service.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(_ context.Context, msg publisher.Message) error {
		published = msg.(*eventproto.Event)
		return nil
	})

	ctx := setToken(t, newIncomingContext(t), true)
	actual, err := service.CreateSubscription(ctx, &proto.CreateSubscriptionRequest{
		EnvironmentId: "ns0",
		Name:          "sname",
		SourceTypes:   []proto.Subscription_SourceType{proto.Subscription_DOMAIN_EVENT_ACCOUNT},
		Recipient: &proto.Recipient{
			Type:                  proto.Recipient_SlackChannel,
			SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: slackWebhookURL},
		},
	})
	require.NoError(t, err)

	// The original URL is persisted, so notifications can still be delivered.
	require.NotNil(t, stored)
	assert.Equal(t, slackWebhookURL, stored.Recipient.SlackChannelRecipient.WebhookUrl)

	// Neither the event payload nor the entity data recorded in the audit log carry it.
	require.NotNil(t, published)
	created := &eventproto.SubscriptionCreatedEvent{}
	require.NoError(t, published.Data.UnmarshalTo(created))
	assert.Equal(t, maskedSlackWebhookURL, created.Recipient.SlackChannelRecipient.WebhookUrl)
	assert.Equal(t, "sname", created.Name)
	assert.NotContains(t, published.EntityData, slackWebhookURL)
	assert.Contains(t, published.EntityData, maskedSlackWebhookURL)

	assert.Equal(t, maskedSlackWebhookURL, actual.Subscription.Recipient.SlackChannelRecipient.WebhookUrl)
	assert.Equal(t, stored.Id, actual.Subscription.Id)
}

func TestUpdateSubscriptionKeepsStoredWebhookURLAndMasksEventAndResponse(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	service := newSubscriptionServiceWithMock(t, mockController)
	runTransaction(service)
	service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
		gomock.Any(), "sub-0", "ns0",
	).Return(storedSubscription(), nil)
	var stored *domain.Subscription
	service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().UpdateSubscription(
		gomock.Any(), gomock.Any(), "ns0",
	).DoAndReturn(func(_ context.Context, s *domain.Subscription, _ string) error {
		stored = s
		return nil
	})
	var published *eventproto.Event
	service.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(_ context.Context, msg publisher.Message) error {
		published = msg.(*eventproto.Event)
		return nil
	})

	ctx := setToken(t, newIncomingContext(t), true)
	actual, err := service.UpdateSubscription(ctx, &proto.UpdateSubscriptionRequest{
		Id:            "sub-0",
		EnvironmentId: "ns0",
		Name:          wrapperspb.String("renamed"),
	})
	require.NoError(t, err)

	// The update must not overwrite the stored URL with the masked value.
	require.NotNil(t, stored)
	assert.Equal(t, slackWebhookURL, stored.Recipient.SlackChannelRecipient.WebhookUrl)
	assert.Equal(t, "renamed", stored.Name)

	require.NotNil(t, published)
	assert.NotContains(t, published.EntityData, slackWebhookURL)
	assert.Contains(t, published.EntityData, maskedSlackWebhookURL)
	assert.NotContains(t, published.PreviousEntityData, slackWebhookURL)
	assert.Contains(t, published.PreviousEntityData, maskedSlackWebhookURL)

	assert.Equal(t, maskedSlackWebhookURL, actual.Subscription.Recipient.SlackChannelRecipient.WebhookUrl)
	assert.Equal(t, "renamed", actual.Subscription.Name)
}

func TestDeleteSubscriptionMasksWebhookURLInEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	service := newSubscriptionServiceWithMock(t, mockController)
	runTransaction(service)
	service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
		gomock.Any(), "sub-0", "ns0",
	).Return(storedSubscription(), nil)
	service.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().DeleteSubscription(
		gomock.Any(), "sub-0", "ns0",
	).Return(nil)
	var published *eventproto.Event
	service.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
		gomock.Any(), gomock.Any(),
	).DoAndReturn(func(_ context.Context, msg publisher.Message) error {
		published = msg.(*eventproto.Event)
		return nil
	})

	ctx := setToken(t, newIncomingContext(t), true)
	_, err := service.DeleteSubscription(ctx, &proto.DeleteSubscriptionRequest{Id: "sub-0", EnvironmentId: "ns0"})
	require.NoError(t, err)

	require.NotNil(t, published)
	assert.Empty(t, published.EntityData)
	assert.NotContains(t, published.PreviousEntityData, slackWebhookURL)
	assert.Contains(t, published.PreviousEntityData, maskedSlackWebhookURL)
}
