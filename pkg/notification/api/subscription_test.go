// Copyright 2023 The Bucketeer Authors.
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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	v2ss "github.com/bucketeer-io/bucketeer/pkg/notification/storage/v2"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/proto/account"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
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
	localizer := locale.NewLocalizer(ctx)
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
		input       *proto.CreateSubscriptionRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrNoCommand",
			setup: nil,
			input: &proto.CreateSubscriptionRequest{
				Command: nil,
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrSourceTypesRequired",
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
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
			desc: "err: ErrRecipientRequired",
			input: &proto.CreateSubscriptionRequest{
				Command: &proto.CreateSubscriptionCommand{
					Name: "sname",
					SourceTypes: []proto.Subscription_SourceType{
						proto.Subscription_DOMAIN_EVENT_ACCOUNT,
						proto.Subscription_DOMAIN_EVENT_FEATURE,
					},
				},
			},
			expectedErr: createError(statusRecipientRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "recipant")),
		},
		{
			desc: "err: ErrSlackRecipientRequired",
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
			expectedErr: createError(statusSlackRecipientRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "slack_recipant")),
		},
		{
			desc: "err: ErrSlackRecipientWebhookURLRequired",
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
			expectedErr: createError(statusSlackRecipientWebhookURLRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_url")),
		},
		{
			desc: "err: ErrNameRequired",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, account.Account_OWNER)
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
	localizer := locale.NewLocalizer(ctx)
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
		input       *proto.UpdateSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			input:       &proto.UpdateSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			input: &proto.UpdateSubscriptionRequest{
				Id: "key-0",
			},
			expectedErr: createError(statusNoCommand, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: add notification types: ErrSourceTypesRequired",
			input: &proto.UpdateSubscriptionRequest{
				Id:                    "key-0",
				AddSourceTypesCommand: &proto.AddSourceTypesCommand{},
			},
			expectedErr: createError(statusSourceTypesRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes")),
		},
		{
			desc: "err: delete notification types: ErrSourceTypesRequired",
			input: &proto.UpdateSubscriptionRequest{
				Id:                       "key-0",
				DeleteSourceTypesCommand: &proto.DeleteSourceTypesCommand{},
			},
			expectedErr: createError(statusSourceTypesRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes")),
		},
		{
			desc: "err: ErrNotFound",
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, account.Account_OWNER)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
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
		input       *proto.EnableSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			input:       &proto.EnableSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			input: &proto.EnableSubscriptionRequest{
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
			input: &proto.EnableSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.EnableSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, account.Account_OWNER)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
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
		input       *proto.DisableSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			input:       &proto.DisableSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			input: &proto.DisableSubscriptionRequest{
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
			input: &proto.DisableSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DisableSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, account.Account_OWNER)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
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
		input       *proto.DeleteSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			input:       &proto.DeleteSubscriptionRequest{},
			expectedErr: createError(statusIDRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			input: &proto.DeleteSubscriptionRequest{
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
			input: &proto.DeleteSubscriptionRequest{
				Id:      "key-0",
				Command: &proto.DeleteSubscriptionCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx = setToken(t, ctx, account.Account_OWNER)
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
	localizer := locale.NewLocalizer(ctx)
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
		input       *proto.GetSubscriptionRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			input:       &proto.GetSubscriptionRequest{},
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
			input:       &proto.GetSubscriptionRequest{Id: "key-0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			ctx = setToken(t, ctx, account.Account_OWNER)
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

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		input       *proto.ListSubscriptionsRequest
		expected    *proto.ListSubscriptionsResponse
		expectedErr error
	}{
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx = setToken(t, ctx, account.Account_OWNER)
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

	patterns := []struct {
		desc        string
		setup       func(*NotificationService)
		input       *proto.ListEnabledSubscriptionsRequest
		expected    *proto.ListEnabledSubscriptionsResponse
		expectedErr error
	}{
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
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNotificationServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(s)
			}
			ctx = setToken(t, ctx, account.Account_OWNER)
			actual, err := s.ListEnabledSubscriptions(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
