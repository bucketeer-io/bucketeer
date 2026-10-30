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

package processor

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	notifiermock "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/sender/notifier/mock"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
	senderproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription/sender"
)

// testSlackWebhookURL is a credential-bearing URL that must never be written to the log.
const testSlackWebhookURL = "https://hooks.slack.com/services/T000/B000/token"

func TestDemoOrganizationCreationNotifierHandleMessageLogsOmitWebhookURL(t *testing.T) {
	t.Parallel()
	notifyErr := errors.New("slack returned 500")
	patterns := []struct {
		desc      string
		notifyErr error
	}{
		{desc: "err: notification failure is logged without the webhook url", notifyErr: notifyErr},
		{desc: "success: notification is sent to the configured webhook"},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			core, logs := observer.New(zapcore.DebugLevel)
			slackNotifier := notifiermock.NewMockNotifier(mockController)
			slackNotifier.EXPECT().Notify(
				gomock.Any(),
				&senderproto.Notification{
					Type: senderproto.Notification_DemoOrganizationCreation,
					DemoOrganizationCreationNotification: &senderproto.DemoOrganizationCreationNotification{
						OwnerEmail:       "owner@example.com",
						OrganizationId:   "org-1",
						OrganizationName: "Demo Org",
					},
				},
				&subscriptionproto.Recipient{
					Type:     subscriptionproto.Recipient_SlackChannel,
					Language: subscriptionproto.Recipient_ENGLISH,
					SlackChannelRecipient: &subscriptionproto.SlackChannelRecipient{
						WebhookUrl: testSlackWebhookURL,
					},
				},
				subscriptionproto.Recipient_ENGLISH,
			).Return(p.notifyErr)
			d := demoOrganizationCreationNotifier{
				slackNotifier: slackNotifier,
				demoOrganizationCreationNotifierConfig: DemoOrganizationCreationNotifierConfig{
					Notifier: NotifierConfig{Slack: SlackNotifierConfig{WebHookURL: testSlackWebhookURL}},
				},
				logger: zap.New(core),
			}
			acked := false
			d.handleMessage(newDemoOrganizationCreatedMessage(t, func() { acked = true }))

			assert.True(t, acked)
			for _, entry := range logs.All() {
				assert.NotContains(t, entry.Message, testSlackWebhookURL)
				for key, field := range entry.ContextMap() {
					assert.NotContains(t, fmt.Sprint(field), testSlackWebhookURL,
						"field %q of %q", key, entry.Message)
				}
			}
			entries := logs.FilterMessage("Failed to send notification").All()
			if p.notifyErr == nil {
				assert.Empty(t, entries)
				return
			}
			require.Len(t, entries, 1)
			fields := entries[0].ContextMap()
			assert.Equal(t, notifyErr.Error(), fields["error"])
			assert.Equal(t, "event-1", fields["event id"])
			assert.Equal(t, "org-1", fields["organizationId"])
		})
	}
}

func newDemoOrganizationCreatedMessage(t *testing.T, ack func()) *puller.Message {
	t.Helper()
	data, err := anypb.New(&domaineventproto.OrganizationCreatedEvent{
		Id:         "org-1",
		Name:       "Demo Org",
		OwnerEmail: "owner@example.com",
	})
	require.NoError(t, err)
	event, err := proto.Marshal(&domaineventproto.Event{
		Id:   "event-1",
		Type: domaineventproto.Event_DEMO_ORGANIZATION_CREATED,
		Data: data,
	})
	require.NoError(t, err)
	return &puller.Message{
		ID:         "msg-1",
		Data:       event,
		Attributes: map[string]string{"id": "event-1"},
		Ack:        ack,
		Nack:       func() { t.Fatal("message must not be nacked") },
	}
}
