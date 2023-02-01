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
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	gstatus "google.golang.org/grpc/status"

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

	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.CreateAdminSubscriptionRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrUnauthenticated",
			setup: nil,
			token: createOwnerToken(t),
			input: &proto.CreateAdminSubscriptionRequest{
				Command: nil,
			},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			token: createAdminToken(t),
			input: &proto.CreateAdminSubscriptionRequest{
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: ErrSourceTypesRequired",
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
			expectedErr: createError(statusSourceTypesRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes")),
		},
		{
			desc:  "err: ErrRecipientRequired",
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
			expectedErr: createError(statusRecipientRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "recipant")),
		},
		{
			desc:  "err: ErrSlackRecipientRequired",
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
			expectedErr: createError(statusSlackRecipientRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "slack_recipant")),
		},
		{
			desc:  "err: ErrSlackRecipientWebhookURLRequired",
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
			expectedErr: createError(statusSlackRecipientWebhookURLRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_url")),
		},
		{
			desc:  "err: ErrNameRequired",
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
			expectedErr: createError(statusNameRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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

	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.UpdateAdminSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrPermissionDenied",
			token:       createOwnerToken(t),
			input:       &proto.UpdateAdminSubscriptionRequest{},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:        "err: ErrIDRequired",
			token:       createAdminToken(t),
			input:       &proto.UpdateAdminSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "err: ErrNoCommand",
			token: createAdminToken(t),
			input: &proto.UpdateAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc:  "err: add notification types: ErrSourceTypesRequired",
			token: createAdminToken(t),
			input: &proto.UpdateAdminSubscriptionRequest{
				Id:                    "key-0",
				AddSourceTypesCommand: &proto.AddAdminSubscriptionSourceTypesCommand{},
			},
			expectedErr: createError(statusSourceTypesRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes")),
		},
		{
			desc:  "err: delete notification types: ErrSourceTypesRequired",
			token: createAdminToken(t),
			input: &proto.UpdateAdminSubscriptionRequest{
				Id:                       "key-0",
				DeleteSourceTypesCommand: &proto.DeleteAdminSubscriptionSourceTypesCommand{},
			},
			expectedErr: createError(statusSourceTypesRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes")),
		},
		{
			desc: "err: ErrNotFound",
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
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success: addSourceTypes",
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
		{
			desc: "success: deleteSourceTypes",
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
		{
			desc: "success: all commands",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.EnableAdminSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrPermissionDenied",
			token:       createOwnerToken(t),
			input:       &proto.EnableAdminSubscriptionRequest{},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:        "err: ErrIDRequired",
			token:       createAdminToken(t),
			input:       &proto.EnableAdminSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "err: ErrNoCommand",
			token: createAdminToken(t),
			input: &proto.EnableAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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

	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.DisableAdminSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			token:       createAdminToken(t),
			input:       &proto.DisableAdminSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "err: ErrNoCommand",
			token: createAdminToken(t),
			input: &proto.DisableAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.DeleteAdminSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			token:       createAdminToken(t),
			input:       &proto.DeleteAdminSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "err: ErrNoCommand",
			token: createAdminToken(t),
			input: &proto.DeleteAdminSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.GetAdminSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrPermissionDenied",
			token:       createOwnerToken(t),
			input:       &proto.GetAdminSubscriptionRequest{},
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:        "err: ErrIDRequired",
			token:       createAdminToken(t),
			input:       &proto.GetAdminSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.ListAdminSubscriptionsRequest
		expected    *proto.ListAdminSubscriptionsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrPermissionDenied",
			setup:       nil,
			token:       createOwnerToken(t),
			input:       &proto.ListAdminSubscriptionsRequest{PageSize: 2, Cursor: ""},
			expected:    nil,
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		token       *token.IDToken
		input       *proto.ListEnabledAdminSubscriptionsRequest
		expected    *proto.ListEnabledAdminSubscriptionsResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrPermissionDenied",
			setup:       nil,
			token:       createOwnerToken(t),
			input:       &proto.ListEnabledAdminSubscriptionsRequest{PageSize: 2, Cursor: ""},
			expected:    nil,
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc: "success",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
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
