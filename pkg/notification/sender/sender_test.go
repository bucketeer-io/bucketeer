// Copyright 2025 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	ncmock "github.com/bucketeer-io/bucketeer/v2/pkg/notification/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/notification/sender/notifier"
	nmock "github.com/bucketeer-io/bucketeer/v2/pkg/notification/sender/notifier/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	notificationproto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/notification/sender"
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
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
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
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
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
				Id:            "id",
				EnvironmentId: storage.AdminEnvironmentID,
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_ENVIRONMENT,
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
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
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
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
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
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
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
				Id:            "id",
				EnvironmentId: storage.AdminEnvironmentID,
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_PROJECT,
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
				Id:            "id",
				EnvironmentId: storage.AdminEnvironmentID,
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_PROJECT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: true,
			},
			expected: nil,
		},
		{
			desc: "success: multiple subscriptions with feature domain event filtering",
			setup: func(t *testing.T, s *sender) {
				// Create 3 subscriptions all for DOMAIN_EVENT_FEATURE:
				// 1. Feature domain subscription with matching tags (should receive notification)
				// 2. Feature domain subscription with non-matching tags (should NOT receive notification)
				// 3. Feature domain subscription with empty tags (should receive notification)
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{
							Id:              "sid0-feature-matching",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{"ios", "production"},
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
						{
							Id:              "sid1-feature-not-matching",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{"android", "staging"},
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
						{
							Id:              "sid2-feature-empty-tags",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{}, // Empty tags means all feature events
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
					}}, nil)
				// Two subscriptions should be notified:
				// 1. The one with matching tags
				// 2. The one with empty tags (receives all feature events)
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(2)
			},
			input: &senderproto.NotificationEvent{
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
				Notification: &senderproto.Notification{
					Type: senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{
						EntityData: `{
							"id": "feature-123",
							"name": "test-feature",
							"tags": ["ios", "production", "v2"]
						}`,
					},
				},
				IsAdminEvent: false,
			},
			expected: nil,
		},
		{
			desc: "success: multiple subscriptions with mixed domain types",
			setup: func(t *testing.T, s *sender) {
				// Test with 3 subscriptions where we're sending an ACCOUNT domain event
				// Only subscriptions for DOMAIN_EVENT_ACCOUNT will be returned by listEnabledSubscriptions
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), gomock.Any()).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{
							Id:              "sid0-for-account",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{"ios"}, // Tags are ignored for non-feature events
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
						{
							Id:              "sid1-for-account",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{"android"},
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
						{
							Id:              "sid2-for-account",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{},
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
					}}, nil)
				// All 3 subscriptions should be notified for account domain events
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(3)
			},
			input: &senderproto.NotificationEvent{
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
				Notification: &senderproto.Notification{
					Type:                    senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{},
				},
				IsAdminEvent: false,
			},
			expected: nil,
		},
		{
			desc: "success: feature event with mixed subscription types in system",
			setup: func(t *testing.T, s *sender) {
				// Simulates a realistic scenario where the system has subscriptions for different source types
				// but only feature subscriptions are returned when querying for DOMAIN_EVENT_FEATURE
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(
					gomock.Any(),
					&notificationproto.ListEnabledSubscriptionsRequest{
						EnvironmentId: "ns0",
						SourceTypes:   []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_FEATURE},
						PageSize:      listRequestSize,
						Cursor:        "",
					},
				).Return(
					&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: []*notificationproto.Subscription{
						{
							Id:              "sid0-feature-web",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{"web"},
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
						{
							Id:              "sid1-feature-mobile",
							EnvironmentId:   "ns0",
							FeatureFlagTags: []string{"mobile", "ios"},
							Recipient:       &notificationproto.Recipient{Language: notificationproto.Recipient_ENGLISH},
						},
						// Note: Account/Project subscriptions won't be included here because
						// listEnabledSubscriptions filters by sourceType
					}}, nil)
				// Only the mobile subscription should be notified (has "ios" tag)
				s.notifiers[0].(*nmock.MockNotifier).EXPECT().Notify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			input: &senderproto.NotificationEvent{
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
				Notification: &senderproto.Notification{
					Type: senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{
						EntityData: `{
							"id": "feature-456",
							"name": "mobile-feature",
							"tags": ["ios", "swift"]
						}`,
					},
				},
				IsAdminEvent: false,
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
					EnvironmentId: "ns0",
					SourceTypes:   []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:      listRequestSize,
					Cursor:        "",
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
					EnvironmentId: "ns0",
					SourceTypes:   []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:      listRequestSize,
					Cursor:        "",
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
					EnvironmentId: "ns0",
					SourceTypes:   []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:      listRequestSize,
					Cursor:        "",
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
					EnvironmentId: "ns0",
					SourceTypes:   []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:      listRequestSize,
					Cursor:        "",
				}).Return(&notificationproto.ListEnabledSubscriptionsResponse{Subscriptions: subs[:listRequestSize]}, nil)
				s.notificationClient.(*ncmock.MockClient).EXPECT().ListEnabledSubscriptions(gomock.Any(), &notificationproto.ListEnabledSubscriptionsRequest{
					EnvironmentId: "ns0",
					SourceTypes:   []notificationproto.Subscription_SourceType{notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT},
					PageSize:      listRequestSize,
					Cursor:        "",
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

func TestCheckForFeatureDomainEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	type inputTest struct {
		subscription *notificationproto.Subscription
		sourceType   notificationproto.Subscription_SourceType
		entityData   string
	}

	patterns := []struct {
		desc     string
		input    inputTest
		expected bool
	}{
		{
			desc: "err: failed to unmarshal",
			input: inputTest{
				subscription: &notificationproto.Subscription{
					Id:              "sub-id",
					Name:            "sub-name",
					EnvironmentId:   "env-id",
					FeatureFlagTags: []string{"ios"},
				},
				sourceType: notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
				entityData: "random-string",
			},
			expected: false,
		},
		{
			desc: "err: feature flag tag not found",
			input: inputTest{
				subscription: &notificationproto.Subscription{
					Id:              "sub-id",
					Name:            "sub-name",
					EnvironmentId:   "env-id",
					FeatureFlagTags: []string{"web"},
				},
				sourceType: notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
				entityData: `{
					"id": "feature-id-1",
					"tags": ["android", "ios"]
				}`,
			},
			expected: false,
		},
		{
			desc: "success: feature flag tag found",
			input: inputTest{
				subscription: &notificationproto.Subscription{
					Id:              "sub-id",
					Name:            "sub-name",
					EnvironmentId:   "env-id",
					FeatureFlagTags: []string{"ios"},
				},
				sourceType: notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
				entityData: `{
					"id": "feature-id-1",
					"tags": ["android", "ios"]
				}`,
			},
			expected: true,
		},
		{
			desc: "success: both subscription and flag have android and ios",
			input: inputTest{
				subscription: &notificationproto.Subscription{
					Id:              "sub-id",
					Name:            "sub-name",
					EnvironmentId:   "env-id",
					FeatureFlagTags: []string{"android", "ios"},
				},
				sourceType: notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
				entityData: `{
					"id": "feature-id-1",
					"tags": ["android", "ios"]
				}`,
			},
			expected: true,
		},
		{
			desc: "success: not a feature domain event",
			input: inputTest{
				subscription: &notificationproto.Subscription{
					Id:              "sub-id",
					Name:            "sub-name",
					EnvironmentId:   "env-id",
					FeatureFlagTags: []string{"ios"},
				},
				sourceType: notificationproto.Subscription_DOMAIN_EVENT_ACCOUNT,
			},
			expected: true,
		},
		{
			desc: "success: no feature flag tags configured",
			input: inputTest{
				subscription: &notificationproto.Subscription{
					Id:              "sub-id",
					Name:            "sub-name",
					EnvironmentId:   "env-id",
					FeatureFlagTags: []string{},
				},
				sourceType: notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
			},
			expected: true,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			sender := createSender(t, mockController)
			send := sender.checkForFeatureDomainEvent(
				p.input.subscription,
				p.input.sourceType,
				p.input.entityData,
			)
			assert.Equal(t, p.expected, send)
		})
	}
}

func TestContainsTags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc     string
		subTags  []string
		ftTags   []string
		expected bool
	}{
		{
			desc:     "no tags in either",
			subTags:  []string{},
			ftTags:   []string{},
			expected: false,
		},
		{
			desc:     "subscription has tags, feature has none",
			subTags:  []string{"android", "ios"},
			ftTags:   []string{},
			expected: false,
		},
		{
			desc:     "feature has tags, subscription has none",
			subTags:  []string{},
			ftTags:   []string{"android", "ios"},
			expected: false,
		},
		{
			desc:     "exact match - single tag",
			subTags:  []string{"android"},
			ftTags:   []string{"android"},
			expected: true,
		},
		{
			desc:     "exact match - multiple tags",
			subTags:  []string{"android", "ios"},
			ftTags:   []string{"android", "ios"},
			expected: true,
		},
		{
			desc:     "partial match - subscription has subset",
			subTags:  []string{"android"},
			ftTags:   []string{"android", "ios"},
			expected: true,
		},
		{
			desc:     "partial match - feature has subset",
			subTags:  []string{"android", "ios"},
			ftTags:   []string{"android"},
			expected: true,
		},
		{
			desc:     "no match",
			subTags:  []string{"web"},
			ftTags:   []string{"android", "ios"},
			expected: false,
		},
		{
			desc:     "case sensitive - no match",
			subTags:  []string{"Android"},
			ftTags:   []string{"android"},
			expected: false,
		},
		{
			desc:     "whitespace in tags - no match",
			subTags:  []string{" android ", " ios "},
			ftTags:   []string{"android", "ios"},
			expected: false, // Won't match because whitespace is not trimmed
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			result := containsTags(p.subTags, p.ftTags)
			assert.Equal(t, p.expected, result)
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
