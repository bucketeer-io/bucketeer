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
	"github.com/bucketeer-io/bucketeer/pkg/token"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func TestCreateAdminSubscriptionMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.CreateAdminSubscriptionRequest
		expectedErr error
	}{
		"err: ErrUnauthenticated": {
			setup: nil,
			token: createOwnerToken(t),
			input: &proto.CreateAdminSubscriptionRequest{
				Command: nil,
			},
			expectedErr: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"err: ErrNoCommand": {
			setup: nil,
			token: createAdminToken(t),
			input: &proto.CreateAdminSubscriptionRequest{
				Command: nil,
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: ErrSourceTypesRequired": {
			token: createAdminToken(t),
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
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
			token: createAdminToken(t),
			input: &proto.CreateAdminSubscriptionRequest{
				Command: &proto.CreateAdminSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
					},
				},
			},
			expectedErr: localizedError(statusRecipientRequired, locale.JaJP),
		},
		"err: ErrSlackRecipientRequired": {
			token: createAdminToken(t),
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
			expectedErr: localizedError(statusSlackRecipientRequired, locale.JaJP),
		},
		"err: ErrSlackRecipientWebhookURLRequired": {
			token: createAdminToken(t),
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
			expectedErr: localizedError(statusSlackRecipientWebhookURLRequired, locale.JaJP),
		},
		"err: ErrNameRequired": {
			token: createAdminToken(t),
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
			token: createAdminToken(t),
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
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, p.token)
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.UpdateAdminSubscriptionRequest
		expectedErr error
	}{
		"err: ErrPermissionDenied": {
			token:       createOwnerToken(t),
			input:       &proto.UpdateAdminSubscriptionRequest{},
			expectedErr: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"err: ErrIDRequired": {
			token:       createAdminToken(t),
			input:       &proto.UpdateAdminSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			token: createAdminToken(t),
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: localizedError(statusNoCommand, locale.JaJP),
		},
		"err: add notification types: ErrSourceTypesRequired": {
			token: createAdminToken(t),
			input: &proto.UpdateAdminSubscriptionRequest{
				Id:                    "key-0",
				AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{},
			},
			expectedErr: localizedError(statusSourceTypesRequired, locale.JaJP),
		},
		"err: delete notification types: ErrSourceTypesRequired": {
			token: createAdminToken(t),
			input: &proto.UpdateAdminSubscriptionRequest{
				Id:                       "key-0",
				DeleteSourceTypesCommand: &proto.DeleteAdminSubscriptionSourceTypesCommand{},
			},
			expectedErr: localizedError(statusSourceTypesRequired, locale.JaJP),
		},
		"err: ErrNotFound": {
			setup: func(s *NotificationService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ss.ErrAdminSubscriptionNotFound)
			},
			token: createAdminToken(t),
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-1",
				AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT,
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
			token: createAdminToken(t),
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
			token: createAdminToken(t),
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
			token: createAdminToken(t),
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
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, p.token)
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.EnableAdminSubscriptionRequest
		expectedErr error
	}{
		"err: ErrPermissionDenied": {
			token:       createOwnerToken(t),
			input:       &proto.EnableAdminSubscriptionRequest{},
			expectedErr: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"err: ErrIDRequired": {
			token:       createAdminToken(t),
			input:       &proto.EnableAdminSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			token: createAdminToken(t),
			input: &proto.EnableAdminSubscriptionRequest{
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
			token: createAdminToken(t),
			input: &proto.EnableAdminSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.EnableAdminSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, p.token)
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.DisableAdminSubscriptionRequest
		expectedErr error
	}{
		"err: ErrIDRequired": {
			token:       createAdminToken(t),
			input:       &proto.DisableAdminSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			token: createAdminToken(t),
			input: &proto.DisableAdminSubscriptionRequest{
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
			token: createAdminToken(t),
			input: &proto.DisableAdminSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DisableAdminSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, p.token)
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.DeleteAdminSubscriptionRequest
		expectedErr error
	}{
		"err: ErrIDRequired": {
			token:       createAdminToken(t),
			input:       &proto.DeleteAdminSubscriptionRequest{},
			expectedErr: localizedError(statusIDRequired, locale.JaJP),
		},
		"err: ErrNoCommand": {
			token: createAdminToken(t),
			input: &proto.DeleteAdminSubscriptionRequest{
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
			token: createAdminToken(t),
			input: &proto.DeleteAdminSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DeleteAdminSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			ctx := createContextWithToken(t, p.token)
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.GetAdminSubscriptionRequest
		expectedErr error
	}{
		"err: ErrPermissionDenied": {
			token:       createOwnerToken(t),
			input:       &proto.GetAdminSubscriptionRequest{},
			expectedErr: localizedError(statusPermissionDenied, locale.JaJP),
		},
		"err: ErrIDRequired": {
			token:       createAdminToken(t),
			input:       &proto.GetAdminSubscriptionRequest{},
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
			token:       createAdminToken(t),
			input:       &proto.GetAdminSubscriptionRequest{Id: "key-0"},
			expectedErr: nil,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			ctx := createContextWithToken(t, p.token)
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.ListAdminSubscriptionsRequest
		expected    *proto.ListAdminSubscriptionsResponse
		expectedErr error
	}{
		"err: ErrPermissionDenied": {
			setup:       nil,
			token:       createOwnerToken(t),
			input:       &proto.ListAdminSubscriptionsRequest{PageSize: 2, Cursor: ""},
			expected:    nil,
			expectedErr: localizedError(statusPermissionDenied, locale.JaJP),
		},
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
			token: createAdminToken(t),
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
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx := createContextWithToken(t, p.token)
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

	patterns := map[string]struct {
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.ListEnabledAdminSubscriptionsRequest
		expected    *proto.ListEnabledAdminSubscriptionsResponse
		expectedErr error
	}{
		"err: ErrPermissionDenied": {
			setup:       nil,
			token:       createOwnerToken(t),
			input:       &proto.ListEnabledAdminSubscriptionsRequest{PageSize: 2, Cursor: ""},
			expected:    nil,
			expectedErr: localizedError(statusPermissionDenied, locale.JaJP),
		},
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
			token: createAdminToken(t),
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
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx := createContextWithToken(t, p.token)
			actual, err := s.ListEnabledAdminSubscriptions(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
