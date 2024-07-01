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

package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	accountproto "github.com/bucketeer-io/bucketeer/proto/account"

	accountclientmock "github.com/bucketeer-io/bucketeer/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

const (
	adminSubscriptionKind = "AdminSubscription"
	subscriptionKind      = "Subscription"
)

func TestNewNotificationService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mysqlClient := mysqlmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	pm := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewNotificationService(mysqlClient, accountClientMock, pm, WithLogger(logger))
	assert.IsType(t, &NotificationService{}, s)
}

func newNotificationServiceWithMock(
	t *testing.T,
	c *gomock.Controller,
) *NotificationService {
	t.Helper()
	return &NotificationService{
		mysqlClient:          mysqlmock.NewMockClient(c),
		accountClient:        accountclientmock.NewMockClient(c),
		domainEventPublisher: publishermock.NewMockPublisher(c),
		logger:               zap.NewNop(),
	}
}

func newNotificationService(c *gomock.Controller, specifiedEnvironmentId *string, specifiedOrgRole *accountproto.AccountV2_Role_Organization, specifiedEnvRole *accountproto.AccountV2_Role_Environment) *NotificationService {
	var or accountproto.AccountV2_Role_Organization
	var er accountproto.AccountV2_Role_Environment
	var envId string
	if specifiedOrgRole != nil {
		or = *specifiedOrgRole
	} else {
		or = accountproto.AccountV2_Role_Organization_ADMIN
	}
	if specifiedEnvRole != nil {
		er = *specifiedEnvRole
	} else {
		er = accountproto.AccountV2_Role_Environment_EDITOR
	}
	if specifiedEnvironmentId != nil {
		envId = *specifiedEnvironmentId
	} else {
		envId = "ns0"
	}

	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: or,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: envId,
					Role:          er,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	mysqlClient := mysqlmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &NotificationService{
		mysqlClient:          mysqlClient,
		accountClient:        accountClientMock,
		domainEventPublisher: publishermock.NewMockPublisher(c),
		logger:               zap.NewNop(),
	}
}

func createContextWithToken(t *testing.T, token *token.AccessToken) context.Context {
	t.Helper()
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createAdminToken(t *testing.T) *token.AccessToken {
	t.Helper()
	return &token.AccessToken{
		Issuer:   "issuer",
		Subject:  "sub",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
}

func createOwnerToken(t *testing.T) *token.AccessToken {
	t.Helper()
	return &token.AccessToken{
		Issuer:   "issuer",
		Subject:  "sub",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
}

type msgLengthMatcher struct{ length int }

func newMsgLengthMatcher(length int) gomock.Matcher {
	return &msgLengthMatcher{length: length}
}

func (m *msgLengthMatcher) Matches(x interface{}) bool {
	return len(x.([]publisher.Message)) == m.length
}

func (m *msgLengthMatcher) String() string {
	return fmt.Sprintf("length: %d", m.length)
}

func putSubscription(t *testing.T, s storage.Client, kind, namespace string, disabled bool) {
	t.Helper()
	key := storage.NewKey("key-0", kind, namespace)
	sourceTypes := []proto.Subscription_SourceType{
		proto.Subscription_DOMAIN_EVENT_ACCOUNT,
		proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
	}
	recipient := &proto.Recipient{
		Type:                  proto.Recipient_SlackChannel,
		SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
	}
	subscription, err := domain.NewSubscription("sname", sourceTypes, recipient)
	subscription.Disabled = disabled
	require.NoError(t, err)
	err = s.Put(context.Background(), key, subscription.Subscription)
	require.NoError(t, err)
}

// convert to pointer
func toPtr[T any](value T) *T {
	return &value
}
