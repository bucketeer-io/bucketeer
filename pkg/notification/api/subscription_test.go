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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	v2ss "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

func TestCreateSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		input       *proto.CreateSubscriptionRequest
		expectedErr error
	}{
		{
			desc: "err: ErrSourceTypesRequired",
			input: &proto.CreateSubscriptionRequest{
				Name: "sname",
				Recipient: &proto.Recipient{
					Type:                  proto.Recipient_SlackChannel,
					SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
				},
			},
			expectedErr: statusSourceTypesRequired.Err(),
		},
		{
			desc: "err: ErrRecipientRequired",
			input: &proto.CreateSubscriptionRequest{
				Name: "sname",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
			},
			expectedErr: statusRecipientRequired.Err(),
		},
		{
			desc: "err: ErrSlackRecipientRequired",
			input: &proto.CreateSubscriptionRequest{
				Name: "sname",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
				Recipient: &proto.Recipient{
					Type: proto.Recipient_SlackChannel,
				},
			},
			expectedErr: statusSlackRecipientRequired.Err(),
		},
		{
			desc: "err: ErrSlackRecipientWebhookURLRequired",
			input: &proto.CreateSubscriptionRequest{
				Name: "sname",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
				Recipient: &proto.Recipient{
					Type:                  proto.Recipient_SlackChannel,
					SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: ""},
				},
			},
			expectedErr: statusSlackRecipientWebhookURLRequired.Err(),
		},
		{
			desc: "err: ErrNameRequired",
			input: &proto.CreateSubscriptionRequest{
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
				Recipient: &proto.Recipient{
					Type:                  proto.Recipient_SlackChannel,
					SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
				},
			},
			expectedErr: statusNameRequired.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().CreateSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.CreateSubscriptionRequest{
				Name: "sname",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
				Recipient: &proto.Recipient{
					Type:                  proto.Recipient_SlackChannel,
					SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, true)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		input       *proto.UpdateSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			input:       &proto.UpdateSubscriptionRequest{},
			expectedErr: statusIDRequired.Err(),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ss.ErrSubscriptionNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2ss.ErrSubscriptionNotFound)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-1",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc: "success: update SourceTypes",
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
						SourceTypes: []proto.Subscription_SourceType{
							proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						},
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().UpdateSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-0",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: rename",
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
						SourceTypes: []proto.Subscription_SourceType{
							proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						},
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)

				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().UpdateSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-0",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_FEATURE,
				},
				Name: wrapperspb.String("rename"),
			},
			expectedErr: nil,
		},
		{
			desc: "success: disable",
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
						SourceTypes: []proto.Subscription_SourceType{
							proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						},
						Disabled: false,
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().UpdateSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id:       "key-0",
				Disabled: wrapperspb.Bool(true),
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, true)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		input       *proto.DeleteSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			input:       &proto.DeleteSubscriptionRequest{},
			expectedErr: statusIDRequired.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().DeleteSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.DeleteSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, true)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DeleteSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc          string
		orgRole       *accountproto.AccountV2_Role_Organization
		envRole       *accountproto.AccountV2_Role_Environment
		isSystemAdmin bool
		setup         func(*NotificationService)
		input         *proto.GetSubscriptionRequest
		expectedErr   error
	}{
		{
			desc:          "err: ErrIDRequired",
			isSystemAdmin: true,
			input:         &proto.GetSubscriptionRequest{},
			expectedErr:   statusIDRequired.Err(),
		},
		{
			desc:          "err: ErrPermissionDenied",
			isSystemAdmin: false,
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			input:         &proto.GetSubscriptionRequest{},
			expectedErr:   statusPermissionDenied.Err(),
		},
		{
			desc:          "success",
			isSystemAdmin: false,
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().GetSubscription(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
					},
				}, nil)
			},
			input:       &proto.GetSubscriptionRequest{Id: "key-0", EnvironmentId: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newNotificationService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(service)
			}
			ctx = setToken(t, ctx, p.isSystemAdmin)
			actual, err := service.GetSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, actual)
			}
		})
	}
}

func TestListSubscriptionsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc          string
		orgRole       *accountproto.AccountV2_Role_Organization
		envRole       *accountproto.AccountV2_Role_Environment
		isSystemAdmin bool
		setup         func(*NotificationService)
		input         *proto.ListSubscriptionsRequest
		expected      *proto.ListSubscriptionsResponse
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			isSystemAdmin: false,
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			input: &proto.ListSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
				EnvironmentId: "ns0",
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:          "success: filter by environmentIDs",
			isSystemAdmin: false,
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *NotificationService) {
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_EDITOR,
							},
							{
								EnvironmentId: "ns1",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().ListSubscriptions(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Subscription{
					{
						Id:            "key-0",
						Name:          "sname",
						EnvironmentId: "ns0",
					},
					{
						Id:            "key-1",
						Name:          "sname1",
						EnvironmentId: "ns1",
					},
				}, 2, int64(2), nil)
			},
			input: &proto.ListSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
				EnvironmentId:  "ns0",
				OrganizationId: "org-0",
			},
			expected: &proto.ListSubscriptionsResponse{Subscriptions: []*proto.Subscription{
				{
					Id:            "key-0",
					Name:          "sname",
					EnvironmentId: "ns0",
				},
				{
					Id:            "key-1",
					Name:          "sname1",
					EnvironmentId: "ns1",
				},
			}, Cursor: "2", TotalCount: 2},
			expectedErr: nil,
		},
		{
			desc:          "success",
			isSystemAdmin: false,
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().ListSubscriptions(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Subscription{}, 0, int64(0), nil)
			},
			input: &proto.ListSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
				EnvironmentId: "ns0",
			},
			expected:    &proto.ListSubscriptionsResponse{Subscriptions: []*proto.Subscription{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNotificationService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(s)
			}
			ctx = setToken(t, ctx, p.isSystemAdmin)
			actual, err := s.ListSubscriptions(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestListEnabledSubscriptionsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc          string
		orgRole       *accountproto.AccountV2_Role_Organization
		envRole       *accountproto.AccountV2_Role_Environment
		isSystemAdmin bool
		setup         func(*NotificationService)
		input         *proto.ListEnabledSubscriptionsRequest
		expected      *proto.ListEnabledSubscriptionsResponse
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			isSystemAdmin: false,
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			input: &proto.ListEnabledSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "1",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
				EnvironmentId: "ns0",
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:          "success",
			isSystemAdmin: false,
			orgRole:       toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:       toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *NotificationService) {
				s.subscriptionStorage.(*storagemock.MockSubscriptionStorage).EXPECT().ListSubscriptions(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Subscription{}, 1, int64(1), nil)
			},
			input: &proto.ListEnabledSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "1",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
				EnvironmentId: "ns0",
			},
			expected:    &proto.ListEnabledSubscriptionsResponse{Subscriptions: []*proto.Subscription{}, Cursor: "1"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNotificationService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(s)
			}
			ctx = setToken(t, ctx, p.isSystemAdmin)
			actual, err := s.ListEnabledSubscriptions(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
