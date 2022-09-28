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

package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	v2ss "github.com/bucketeer-io/bucketeer/pkg/notification/storage/v2"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func TestCreateSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.CreateSubscriptionRequest
		expectedErr error
	}{
		"err: ErrNoCommand": {
			setup: nil,
			input: &proto.CreateSubscriptionRequest{
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrSourceTypesRequired": {
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
					Name: "sname",
					Recipient: &proto.Recipient{
						Type:                  proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
					},
				},
			},
			expectedErr: localizedError(statusSourceTypesRequired, locale.JaJP),
		},
		"err: ErrRecipientRequired": {
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
				},
			},
			expectedErr: localizedError(statusRecipientRequired, locale.JaJP),
		},
		"err: ErrSlackRecipientRequired": {
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &proto.Recipient{
						Type: proto.Recipient_SlackChannel,
					},
				},
			},
			expectedErr: localizedError(statusSlackRecipientRequired, locale.JaJP),
		},
		"err: ErrSlackRecipientWebhookURLRequired": {
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
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
			},
			expectedErr: localizedError(statusSlackRecipientWebhookURLRequired, locale.JaJP),
		},
		"err: ErrNameRequired": {
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
					Recipient: &proto.Recipient{
						Type:                  proto.Recipient_SlackChannel,
						SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "url"},
					},
				},
			},
			expectedErr: localizedError(statusNameRequired, locale.JaJP),
		},
		"success": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
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
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, createAdminToken(t))
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.UpdateSubscriptionRequest
		expectedErr error
	}{
		"err: ErrIDRequired": {
			input:       &proto.UpdateSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: add notification types: ErrSourceTypesRequired": {
			input: &proto.UpdateSubscriptionRequest{
				Id:                    "key-0",
				AddSourceTypesCommand: &proto.AddSourceTypesCommand{},
			},
			expectedErr: localizedError(statusSourceTypesRequired, locale.JaJP),
		},
		"err: delete notification types: ErrSourceTypesRequired": {
			input: &proto.UpdateSubscriptionRequest{
				Id:                       "key-0",
				DeleteSourceTypesCommand: &proto.DeleteSourceTypesCommand{},
			},
			expectedErr: localizedError(statusSourceTypesRequired, locale.JaJP),
		},
		"err: ErrNotFound": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ss.ErrSubscriptionNotFound)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-1",
				AddSourceTypesCommand: &proto.AddSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
				},
			},
			expectedErr: localizedError(statusNotFound, locale.JaJP),
		},
		"success: addSourceTypes": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-0",
				AddSourceTypesCommand: &proto.AddSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
				},
			},
			expectedErr: nil,
		},
		"success: deleteSourceTypes": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-0",
				DeleteSourceTypesCommand: &proto.DeleteSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
				},
			},
			expectedErr: nil,
		},
		"success: all commands": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-0",
				AddSourceTypesCommand: &proto.AddSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
				},
				DeleteSourceTypesCommand: &proto.DeleteSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					},
				},
				RenameSubscriptionCommand: &proto.RenameSubscriptionCommand{
					Name: "rename",
				},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, createAdminToken(t))
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdateSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestEnableSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.EnableSubscriptionRequest
		expectedErr error
	}{
		"err: ErrIDRequired": {
			input:       &proto.EnableSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			input: &proto.EnableSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"success": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.EnableSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.EnableSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, createAdminToken(t))
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.EnableSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDisableSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.DisableSubscriptionRequest
		expectedErr error
	}{
		"err: ErrIDRequired": {
			input:       &proto.DisableSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			input: &proto.DisableSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"success": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.DisableSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DisableSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, createAdminToken(t))
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DisableSubscription(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.DeleteSubscriptionRequest
		expectedErr error
	}{
		"err: ErrIDRequired": {
			input:       &proto.DeleteSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			input: &proto.DeleteSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"success": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.domainEventPublisher.(*publishermock.MockPublisher).EXPECT().PublishMulti(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.DeleteSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DeleteSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, createAdminToken(t))
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.GetSubscriptionRequest
		expectedErr error
	}{
		"err: ErrIDRequired": {
			input:       &proto.GetSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"success": {
			setup: func(s *NotificationService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:       &proto.GetSubscriptionRequest{Id: "key-0"},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			ctx := createContextWithToken(t, createAdminToken(t))
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.ListSubscriptionsRequest
		expected    *proto.ListSubscriptionsResponse
		expectedErr error
	}{
		"success": {
			setup: func(s *NotificationService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: &proto.ListSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
			},
			expected:    &proto.ListSubscriptionsResponse{Subscriptions: []*proto.Subscription{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx := createContextWithToken(t, createAdminToken(t))
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		input       *proto.ListEnabledSubscriptionsRequest
		expected    *proto.ListEnabledSubscriptionsResponse
		expectedErr error
	}{
		"success": {
			setup: func(s *NotificationService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: &proto.ListEnabledSubscriptionsRequest{
				PageSize: 2,
				Cursor:   "1",
				SourceTypes: []proto.Subscription_SourceType{
					proto.Subscription_DOMAIN_EVENT_ACCOUNT,
					proto.Subscription_DOMAIN_EVENT_SUBSCRIPTION,
				},
			},
			expected:    &proto.ListEnabledSubscriptionsResponse{Subscriptions: []*proto.Subscription{}, Cursor: "1"},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx := createContextWithToken(t, createAdminToken(t))
			actual, err := s.ListEnabledSubscriptions(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
