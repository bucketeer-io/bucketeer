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

package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
	senderproto "github.com/bucketeer-io/bucketeer/proto/notification/sender"
)

func TestCreateNotificationEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	envName := "env-name"
	envURLCode := "env-url-code"
	patterns := []struct {
		desc               string
		input              *domaineventproto.Event
		environmentName    string
		environmentURLCode string
		expected           *senderproto.NotificationEvent
		expectedErr        error
	}{
		{
			desc: "success: DomainEvent",
			input: &domaineventproto.Event{
				Id:            "did",
				EntityType:    domaineventproto.Event_FEATURE,
				EntityId:      "fid",
				Type:          domaineventproto.Event_FEATURE_CREATED,
				Editor:        &domaineventproto.Editor{Email: "test@test.com"},
				EnvironmentId: "ns0",
				IsAdminEvent:  false,
			},
			environmentName:    envName,
			environmentURLCode: envURLCode,
			expected: &senderproto.NotificationEvent{
				Id:            "id",
				EnvironmentId: "ns0",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_FEATURE,
				Notification: &senderproto.Notification{
					Type: senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{
						EnvironmentName:    envName,
						EnvironmentUrlCode: envURLCode,
						Editor:             &domaineventproto.Editor{Email: "test@test.com"},
						EntityType:         domaineventproto.Event_FEATURE,
						EntityId:           "fid",
						Type:               domaineventproto.Event_FEATURE_CREATED,
					},
				},
				IsAdminEvent: false,
			},
			expectedErr: nil,
		},
		{
			desc: "success: Admin DomainEvent",
			input: &domaineventproto.Event{
				Id:            "did",
				EntityType:    domaineventproto.Event_PROJECT,
				EntityId:      "pid",
				Type:          domaineventproto.Event_PROJECT_CREATED,
				Editor:        &domaineventproto.Editor{Email: "test@test.com"},
				EnvironmentId: "",
				IsAdminEvent:  true,
			},
			environmentName:    envName,
			environmentURLCode: envURLCode,
			expected: &senderproto.NotificationEvent{
				Id:            "id",
				EnvironmentId: "",
				SourceType:    notificationproto.Subscription_DOMAIN_EVENT_PROJECT,
				Notification: &senderproto.Notification{
					Type: senderproto.Notification_DomainEvent,
					DomainEventNotification: &senderproto.DomainEventNotification{
						EnvironmentName:    envName,
						EnvironmentUrlCode: envURLCode,
						Editor:             &domaineventproto.Editor{Email: "test@test.com"},
						EntityType:         domaineventproto.Event_PROJECT,
						EntityId:           "pid",
						Type:               domaineventproto.Event_PROJECT_CREATED,
					},
				},
				IsAdminEvent: true,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			i := newDomainEventInformer(t, mockController)
			actual, err := i.createNotificationEvent(p.input, p.environmentName, p.environmentURLCode, p.input.IsAdminEvent)
			assert.Equal(t, p.expectedErr, err)
			if p.expected != nil {
				assert.Equal(t, p.expected.EnvironmentId, actual.EnvironmentId)
				assert.Equal(t, p.expected.SourceType, actual.SourceType)
				assert.Equal(t, p.expected.IsAdminEvent, actual.IsAdminEvent)
				assert.Equal(t, p.expected.Notification.Type, actual.Notification.Type)
				assert.Equal(t, p.expected.Notification.DomainEventNotification.EnvironmentName, actual.Notification.DomainEventNotification.EnvironmentName)
				assert.Equal(t, p.expected.Notification.DomainEventNotification.EnvironmentUrlCode, actual.Notification.DomainEventNotification.EnvironmentUrlCode)
				assert.Equal(t, p.expected.Notification.DomainEventNotification.Editor, actual.Notification.DomainEventNotification.Editor)
				assert.Equal(t, p.expected.Notification.DomainEventNotification.EntityType, actual.Notification.DomainEventNotification.EntityType)
				assert.Equal(t, p.expected.Notification.DomainEventNotification.EntityId, actual.Notification.DomainEventNotification.EntityId)
				assert.Equal(t, p.expected.Notification.DomainEventNotification.Type, actual.Notification.DomainEventNotification.Type)
			}
		})
	}
}

func TestConvSourceType(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	for k, v := range domaineventproto.Event_EntityType_name {
		t.Run(v, func(t *testing.T) {
			i := newDomainEventInformer(t, mockController)
			_, err := i.convSourceType(domaineventproto.Event_EntityType(k))
			assert.NoError(t, err)
		})
	}
}

func newDomainEventInformer(t *testing.T, c *gomock.Controller) *domainEventInformer {
	t.Helper()
	return &domainEventInformer{
		logger: zap.NewNop(),
	}
}
