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

package sender

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	ncmock "github.com/bucketeer-io/bucketeer/pkg/notification/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier"
	nmock "github.com/bucketeer-io/bucketeer/pkg/notification/sender/notifier/mock"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

func TestHandle(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc     string
		setup    func(t *testing.T, s *sender)
		input    *senderproto.NotificationEvent
		expected error
	}{
		{
			desc: "error: list subscriptions",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), gomock.Any()).Return(
					nil, errors.New("test"))
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: "ns0",
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: false,
			},
			expected: errors.New("test"),
		},
		{
			desc: "success: 0 subscription",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{}, nil)
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: "ns0",
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: false,
			},
			expected: nil,
		},
		{
			desc: "success: 0 admin subscription",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledAdminSubscriptionsResponse{}, nil)
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: storage.AdminEnvironmentNamespace,
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_ENVIRONMENT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: true,
			},
			expected: nil,
		},
		{
			desc: "error: notify",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{Id: "sid0", Recipient: &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH}},
					}}, nil)
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: "ns0",
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: false,
			},
			expected: errors.New("test"),
		},
		{
			desc: "success: 1 subscription",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{Id: "sid0", Recipient: &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH}},
					}}, nil)
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: "ns0",
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: false,
			},
			expected: nil,
		},
		{
			desc: "success: 2 subscription",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{Id: "sid0", Recipient: &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH}}, {Id: "sid1", Recipient: &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH}},
					}}, nil)
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: "ns0",
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: false,
			},
			expected: nil,
		},
		{
			desc: "success: 1 admin subscription",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledAdminSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{Id: "sid0", Recipient: &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH}},
					}}, nil)
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: storage.AdminEnvironmentNamespace,
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_PROJECT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: true,
			},
			expected: nil,
		},
		{
			desc: "success: 2 admin subscription",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledAdminSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{Id: "sid0", Recipient: &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH}}, {Id: "sid1", Recipient: &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH}},
					}}, nil)
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			input: &senderproto.NotificationEvent{
				Id:                   "id",
				EnvironmentNamespace: storage.AdminEnvironmentNamespace,
				SourceType:           notificationproto.Subscription_DOMAIN_EVENT_PROJECT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: true,
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sender := createSender(t, mockController)
			if p.setup != nil {
				p.setup(t, sender)
			}
			err := sender.Send(context.Background(), p.input)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestListEnabledSubscriptions(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(t *testing.T, s *sender)
		input       notificationproto.Subscription_SourceType
		expected    []*notificationproto.Subscription
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), &notificationproto.ListEnabledSubscriptionsRequest{
					EnvironmentNamespace: "ns0",
					SourceTypes:          []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:             listRequestSize,
					Cursor:               "",
				}).Return(nil, errors.New("test"))
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    nil,
			expectedErr: errors.New("test"),
		},
		{
			desc: "success: 0 entity",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), &notificationproto.ListEnabledSubscriptionsRequest{
					EnvironmentNamespace: "ns0",
					SourceTypes:          []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:             listRequestSize,
					Cursor:               "",
				}).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{}}, nil)
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    []*notificationproto.Subscription{},
			expectedErr: nil,
		},
		{
			desc: "success: 1 entity",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), &notificationproto.ListEnabledSubscriptionsRequest{
					EnvironmentNamespace: "ns0",
					SourceTypes:          []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:             listRequestSize,
					Cursor:               "",
				}).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{Id: "sid0"},
					}}, nil)
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    []*notificationproto.Subscription{{Id: "sid0"}},
			expectedErr: nil,
		},
		{
			desc: "success: listRequestSize + 1 entity",
			setup: func(t *testing.T, s *sender) {
				subs := createSubscriptions(t, listRequestSize+1)
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), &notificationproto.ListEnabledSubscriptionsRequest{
					EnvironmentNamespace: "ns0",
					SourceTypes:          []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:             listRequestSize,
					Cursor:               "",
				}).Return(&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: subs[:listRequestSize]}, nil)
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), &notificationproto.ListEnabledSubscriptionsRequest{
					EnvironmentNamespace: "ns0",
					SourceTypes:          []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:             listRequestSize,
					Cursor:               "",
				}).Return(&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: subs[listRequestSize : listRequestSize+1]}, nil)
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    createSubscriptions(t, listRequestSize+1),
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sender := createSender(t, mockController)
			if p.setup != nil {
				p.setup(t, sender)
			}
			actual, err := sender.listEnabledSubscriptions(context.Background(), "ns0", p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListEnabledAdminSubscriptions(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(t *testing.T, s *sender)
		input       notificationproto.Subscription_SourceType
		expected    []*notificationproto.Subscription
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), &notificationproto.ListEnabledAdminSubscriptionsRequest{
					SourceTypes: []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:    listRequestSize,
					Cursor:      "",
				}).Return(nil, errors.New("test"))
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    nil,
			expectedErr: errors.New("test"),
		},
		{
			desc: "success: 0 entity",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), &notificationproto.ListEnabledAdminSubscriptionsRequest{
					SourceTypes: []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:    listRequestSize,
					Cursor:      "",
				}).Return(
					&notificationproto.ListEnabledAdminSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{}}, nil)
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    []*notificationproto.Subscription{},
			expectedErr: nil,
		},
		{
			desc: "success: 1 entity",
			setup: func(t *testing.T, s *sender) {
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), &notificationproto.ListEnabledAdminSubscriptionsRequest{
					SourceTypes: []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:    listRequestSize,
					Cursor:      "",
				}).Return(
					&notificationproto.ListEnabledAdminSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{Id: "sid0"},
					}}, nil)
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    []*notificationproto.Subscription{{Id: "sid0"}},
			expectedErr: nil,
		},
		{
			desc: "success: listRequestSize + 1 entity",
			setup: func(t *testing.T, s *sender) {
				subs := createSubscriptions(t, listRequestSize+1)
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), &notificationproto.ListEnabledAdminSubscriptionsRequest{
					SourceTypes: []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:    listRequestSize,
					Cursor:      "",
				}).Return(&notificationproto.ListEnabledAdminSubscriptionsResponse{Subscriptions: subs[:listRequestSize]}, nil)
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledAdminSubscriptions(gomock.Any(), &notificationproto.ListEnabledAdminSubscriptionsRequest{
					SourceTypes: []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:    listRequestSize,
					Cursor:      "",
				}).Return(&notificationproto.ListEnabledAdminSubscriptionsResponse{Subscriptions: subs[listRequestSize : listRequestSize+1]}, nil)
			},
			input:       notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			expected:    createSubscriptions(t, listRequestSize+1),
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sender := createSender(t, mockController)
			if p.setup != nil {
				p.setup(t, sender)
			}
			actual, err := sender.listEnabledAdminSubscriptions(context.Background(), p.input)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createSubscriptions(t *testing.T, size int) []*notificationproto.Subscription {
	subscriptions := []*notificationproto.Subscription{}
	for i := 0; i < size; i++ {
		subscriptions = append(subscriptions, &notificationproto.Subscription{Id: fmt.Sprintf("sid%d", i)})
	}
	return subscriptions
}

func createSender(t *testing.T, c *gomock.Controller) *sender {
	ncMock := ncmock.NewMockClient(c)
	nMock := nmock.NewMockNotifier(c)
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &sender{
		notificationClient: ncMock,
		notifiers:          []notifier.Notifier{nMock},
		logger:             logger,
	}
}
