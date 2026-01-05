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
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/domain"
	v2ss "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage/v2"
	staragemock "github.com/bucketeer-io/bucketeer/v2/pkg/notification/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

func TestCreateAdminSubscriptionMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.CreateAdminSubscriptionRequest
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			setup:         nil,
			isSystemAdmin: false,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: nil,
			},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:          "err: ErrNoCommand",
			setup:         nil,
			isSystemAdmin: true,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: nil,
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc:          "err: ErrSourceTypesRequired",
			isSystemAdmin: true,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
					Name: "sname",
					Recipient: &proto.Recipient{
						Type:                  proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
					},
				},
			},
			expectedErr: statusSourceTypesRequired.Err(),
		},
		{
			desc:          "err: ErrRecipientRequired",
			isSystemAdmin: true,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					},
				},
			},
			expectedErr: statusRecipientRequired.Err(),
		},
		{
			desc:          "err: ErrSlackRecipientRequired",
			isSystemAdmin: true,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
					},
				},
			},
			expectedErr: statusSlackRecipientRequired.Err(),
		},
		{
			desc:          "err: ErrSlackRecipientWebhookURLRequired",
			isSystemAdmin: true,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type:                  proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: ""},
					},
				},
			},
			expectedErr: statusSlackRecipientWebhookURLRequired.Err(),
		},
		{
			desc:          "err: ErrNameRequired",
			isSystemAdmin: true,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type:                  proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
					},
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
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().CreateAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					},
					Recipient: &proto.Recipient{
						Type:                  proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, p.isSystemAdmin)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreateAdminSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAdminSubscriptionMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.UpdateAdminSubscriptionRequest
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			isSystemAdmin: false,
			input:         &proto.UpdateAdminSubscriptionRequest{},
			expectedErr:   statusPermissionDenied.Err(),
		},
		{
			desc:          "err: ErrIDRequired",
			isSystemAdmin: true,
			input:         &proto.UpdateAdminSubscriptionRequest{},
			expectedErr:   statusIDRequired.Err(),
		},
		{
			desc:          "err: ErrNoCommand",
			isSystemAdmin: true,
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc:          "err: add notification types: ErrSourceTypesRequired",
			isSystemAdmin: true,
			input: &proto.UpdateAdminSubscriptionRequest{
				Id:                    "key-0",
				AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{},
			},
			expectedErr: statusSourceTypesRequired.Err(),
		},
		{
			desc:          "err: delete notification types: ErrSourceTypesRequired",
			isSystemAdmin: true,
			input: &proto.UpdateAdminSubscriptionRequest{
				Id:                       "key-0",
				DeleteSourceTypesCommand: &proto.DeleteAdminSubscriptionSourceTypesCommand{},
			},
			expectedErr: statusSourceTypesRequired.Err(),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2ss.ErrAdminSubscriptionNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2ss.ErrAdminSubscriptionNotFound)
			},
			isSystemAdmin: true,
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-1",
				AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					},
				},
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc: "success: addSourceTypes",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
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
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().UpdateAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-0",
				AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: deleteSourceTypes",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
						SourceTypes: []proto.Subscription_SourceType{
							proto.Subscription_DOMAIN_EVENT_ACCOUNT,
							proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
						},
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().UpdateAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-0",
				DeleteSourceTypesCommand: &proto.DeleteAdminSubscriptionSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: all commands",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
						SourceTypes: []proto.Subscription_SourceType{
							proto.Subscription_DOMAIN_EVENT_ACCOUNT,
							proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
						},
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().UpdateAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-0",
				AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
				},
				DeleteSourceTypesCommand: &proto.DeleteAdminSubscriptionSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
				},
				RenameSubscriptionCommand: &proto.RenameAdminSubscriptionCommand{
					Name: "rename",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, p.isSystemAdmin)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateAdminSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestEnableAdminSubscriptionMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.EnableAdminSubscriptionRequest
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			isSystemAdmin: false,
			input:         &proto.EnableAdminSubscriptionRequest{},
			expectedErr:   statusPermissionDenied.Err(),
		},
		{
			desc:          "err: ErrIDRequired",
			isSystemAdmin: true,
			input:         &proto.EnableAdminSubscriptionRequest{},
			expectedErr:   statusIDRequired.Err(),
		},
		{
			desc:          "err: ErrNoCommand",
			isSystemAdmin: true,
			input: &proto.EnableAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id:       "key-0",
						Disabled: true,
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().UpdateAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			input: &proto.EnableAdminSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.EnableAdminSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, p.isSystemAdmin)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.EnableAdminSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableAdminSubscriptionMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.DisableAdminSubscriptionRequest
		expectedErr   error
	}{
		{
			desc:          "err: ErrIDRequired",
			isSystemAdmin: true,
			input:         &proto.DisableAdminSubscriptionRequest{},
			expectedErr:   statusIDRequired.Err(),
		},
		{
			desc:          "err: ErrNoCommand",
			isSystemAdmin: true,
			input: &proto.DisableAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id:       "key-0",
						Disabled: false,
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().UpdateAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			input: &proto.DisableAdminSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DisableAdminSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, p.isSystemAdmin)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DisableAdminSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAdminSubscriptionMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.DeleteAdminSubscriptionRequest
		expectedErr   error
	}{
		{
			desc:          "err: ErrIDRequired",
			isSystemAdmin: true,
			input:         &proto.DeleteAdminSubscriptionRequest{},
			expectedErr:   statusIDRequired.Err(),
		},
		{
			desc:          "err: ErrNoCommand",
			isSystemAdmin: true,
			input: &proto.DeleteAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
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
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().DeleteAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			isSystemAdmin: true,
			input: &proto.DeleteAdminSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DeleteAdminSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, p.isSystemAdmin)
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DeleteAdminSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAdminSubscriptionMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.GetAdminSubscriptionRequest
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			isSystemAdmin: false,
			input:         &proto.GetAdminSubscriptionRequest{},
			expectedErr:   statusPermissionDenied.Err(),
		},
		{
			desc:          "err: ErrIDRequired",
			isSystemAdmin: true,
			input:         &proto.GetAdminSubscriptionRequest{},
			expectedErr:   statusIDRequired.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().GetAdminSubscription(
					gomock.Any(), gomock.Any(),
				).Return(&domain.Subscription{
					Subscription: &proto.Subscription{
						Id: "key-0",
					},
				}, nil)
			},
			isSystemAdmin: true,
			input:         &proto.GetAdminSubscriptionRequest{Id: "key-0"},
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			ctx = setToken(t, ctx, p.isSystemAdmin)
			actual, err := service.GetAdminSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, actual)
			}
		})
	}
}

func TestListAdminSubscriptionsMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.ListAdminSubscriptionsRequest
		expected      *proto.ListAdminSubscriptionsResponse
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			setup:         nil,
			isSystemAdmin: false,
			input:         &proto.ListAdminSubscriptionsRequest{PageSize: 2, Cursor: ""},
			expected:      nil,
			expectedErr:   statusPermissionDenied.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().ListAdminSubscriptions(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Subscription{}, 0, int64(0), nil)

			},
			isSystemAdmin: true,
			input: &proto.ListAdminSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_ADMIN_SUBSCRIPTION,
				},
			},
			expected:    &proto.ListAdminSubscriptionsResponse{Subscriptions: []*proto.Subscription{}, Cursor: "0", TotalCount: 0},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx = setToken(t, ctx, p.isSystemAdmin)
			actual, err := s.ListAdminSubscriptions(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestListEnabledAdminSubscriptionsMySQL(t *testing.T) {
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
		setup         func(*NotificationService)
		isSystemAdmin bool
		input         *proto.ListEnabledAdminSubscriptionsRequest
		expected      *proto.ListEnabledAdminSubscriptionsResponse
		expectedErr   error
	}{
		{
			desc:          "err: ErrPermissionDenied",
			setup:         nil,
			isSystemAdmin: false,
			input:         &proto.ListEnabledAdminSubscriptionsRequest{PageSize: 2, Cursor: ""},
			expected:      nil,
			expectedErr:   statusPermissionDenied.Err(),
		},
		{
			desc: "success",
			setup: func(s *NotificationService) {
				s.adminSubscriptionStorage.(*staragemock.MockAdminSubscriptionStorage).EXPECT().ListAdminSubscriptions(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Subscription{}, 1, int64(1), nil)
			},
			isSystemAdmin: true,
			input: &proto.ListEnabledAdminSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "1",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_ADMIN_SUBSCRIPTION,
				},
			},
			expected:    &proto.ListEnabledAdminSubscriptionsResponse{Subscriptions: []*proto.Subscription{}, Cursor: "1"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx = setToken(t, ctx, p.isSystemAdmin)
			actual, err := s.ListEnabledAdminSubscriptions(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func setToken(t *testing.T, ctx context.Context, isSystemAdmin bool) context.Context {
	t.Helper()
	tokenID := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: isSystemAdmin,
	}
	return context.WithValue(ctx, rpc.AccessTokenKey, tokenID)
}
