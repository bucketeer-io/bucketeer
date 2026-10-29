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

package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

const (
	rawWebhookURL    = "https://hooks.slack.com/services/T000/B000/token"
	maskedWebhookURL = "https://hooks.slack.com/****"
)

// subscriptionEntityData encodes a subscription the same way the domain event does.
func subscriptionEntityData(t *testing.T, webhookURL string) string {
	t.Helper()
	data, err := json.MarshalIndent(&subscriptionproto.Subscription{
		Id:        "sub-1",
		Name:      "name-1",
		CreatedAt: 1739746800,
		Recipient: &subscriptionproto.Recipient{
			Type:                  subscriptionproto.Recipient_SlackChannel,
			SlackChannelRecipient: &subscriptionproto.SlackChannelRecipient{WebhookUrl: webhookURL},
		},
	}, "", "  ")
	require.NoError(t, err)
	return string(data)
}

func TestMaskWebhookURLEntityData(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		input       string
		expected    string
		expectedErr bool
	}{
		{
			desc:     "empty",
			input:    "",
			expected: "",
		},
		{
			desc:        "invalid json: dropped",
			input:       "not a json",
			expected:    "",
			expectedErr: true,
		},
		{
			desc:     "no recipient: returned as it is",
			input:    "{\n  \"id\": \"sub-1\"\n}",
			expected: "{\n  \"id\": \"sub-1\"\n}",
		},
		{
			desc:     "no slack channel recipient: returned as it is",
			input:    "{\n  \"recipient\": {\n    \"type\": 0\n  }\n}",
			expected: "{\n  \"recipient\": {\n    \"type\": 0\n  }\n}",
		},
		{
			desc:  "webhook url masked and numbers kept as they are",
			input: subscriptionEntityData(t, rawWebhookURL),
			// The fields are sorted because the entity data is re-encoded from a map.
			expected: "{\n" +
				"  \"created_at\": 1739746800,\n" +
				"  \"disabled\": false,\n" +
				"  \"environment_id\": \"\",\n" +
				"  \"environment_name\": \"\",\n" +
				"  \"feature_flag_tags\": null,\n" +
				"  \"id\": \"sub-1\",\n" +
				"  \"name\": \"name-1\",\n" +
				"  \"recipient\": {\n" +
				"    \"language\": 0,\n" +
				"    \"slack_channel_recipient\": {\n" +
				"      \"webhook_url\": \"" + maskedWebhookURL + "\"\n" +
				"    },\n" +
				"    \"type\": 0\n" +
				"  },\n" +
				"  \"source_types\": null,\n" +
				"  \"updated_at\": 0\n" +
				"}",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := MaskWebhookURLEntityData(p.input)
			assert.Equal(t, p.expected, actual)
			if p.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestMaskWebhookURL(t *testing.T) {
	t.Parallel()
	createdEvent, err := anypb.New(&eventproto.SubscriptionCreatedEvent{
		Name: "name-1",
		Recipient: &subscriptionproto.Recipient{
			Type:                  subscriptionproto.Recipient_SlackChannel,
			SlackChannelRecipient: &subscriptionproto.SlackChannelRecipient{WebhookUrl: rawWebhookURL},
		},
	})
	require.NoError(t, err)
	updatedEvent, err := anypb.New(&eventproto.SubscriptionUpdatedEvent{Id: "sub-1"})
	require.NoError(t, err)
	featureEvent, err := anypb.New(&eventproto.FeatureCreatedEvent{Id: "feature-1"})
	require.NoError(t, err)

	patterns := []struct {
		desc  string
		input *proto.AuditLog
	}{
		{
			desc:  "nil",
			input: nil,
		},
		{
			desc: "not a subscription entity: untouched",
			input: &proto.AuditLog{
				EntityType: eventproto.Event_FEATURE,
				Type:       eventproto.Event_FEATURE_CREATED,
				Event:      featureEvent,
				EntityData: subscriptionEntityData(t, rawWebhookURL),
			},
		},
		{
			desc: "subscription created: entity data and event data masked",
			input: &proto.AuditLog{
				EntityType: eventproto.Event_SUBSCRIPTION,
				Type:       eventproto.Event_SUBSCRIPTION_CREATED,
				Event:      createdEvent,
				EntityData: subscriptionEntityData(t, rawWebhookURL),
			},
		},
		{
			desc: "subscription updated: both entity data masked",
			input: &proto.AuditLog{
				EntityType:         eventproto.Event_SUBSCRIPTION,
				Type:               eventproto.Event_SUBSCRIPTION_UPDATED,
				Event:              updatedEvent,
				EntityData:         subscriptionEntityData(t, rawWebhookURL),
				PreviousEntityData: subscriptionEntityData(t, rawWebhookURL),
			},
		},
		{
			desc: "subscription deleted: previous entity data masked",
			input: &proto.AuditLog{
				EntityType:         eventproto.Event_SUBSCRIPTION,
				Type:               eventproto.Event_SUBSCRIPTION_DELETED,
				PreviousEntityData: subscriptionEntityData(t, rawWebhookURL),
			},
		},
		{
			desc: "already masked: unchanged",
			input: &proto.AuditLog{
				EntityType: eventproto.Event_SUBSCRIPTION,
				Type:       eventproto.Event_SUBSCRIPTION_UPDATED,
				Event:      updatedEvent,
				EntityData: subscriptionEntityData(t, maskedWebhookURL),
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			require.NoError(t, MaskWebhookURL(p.input))
			if p.input == nil {
				return
			}
			if p.input.EntityType != eventproto.Event_SUBSCRIPTION {
				assert.Contains(t, p.input.EntityData, rawWebhookURL)
				assert.Equal(t, featureEvent, p.input.Event)
				return
			}
			assert.NotContains(t, p.input.EntityData, rawWebhookURL)
			assert.NotContains(t, p.input.PreviousEntityData, rawWebhookURL)
			if p.input.EntityData != "" {
				assert.Contains(t, p.input.EntityData, maskedWebhookURL)
			}
			if p.input.PreviousEntityData != "" {
				assert.Contains(t, p.input.PreviousEntityData, maskedWebhookURL)
			}
			if p.input.Type == eventproto.Event_SUBSCRIPTION_CREATED {
				created := &eventproto.SubscriptionCreatedEvent{}
				require.NoError(t, p.input.Event.UnmarshalTo(created))
				assert.Equal(t, maskedWebhookURL, created.Recipient.SlackChannelRecipient.WebhookUrl)
				assert.Equal(t, "name-1", created.Name)
				return
			}
			if p.input.Event != nil {
				assert.Equal(t, updatedEvent, p.input.Event)
			}
		})
	}
}

func TestNewAuditLogMasksWebhookURL(t *testing.T) {
	t.Parallel()
	createdEvent, err := anypb.New(&eventproto.SubscriptionCreatedEvent{
		Recipient: &subscriptionproto.Recipient{
			Type:                  subscriptionproto.Recipient_SlackChannel,
			SlackChannelRecipient: &subscriptionproto.SlackChannelRecipient{WebhookUrl: rawWebhookURL},
		},
	})
	require.NoError(t, err)
	entityData := subscriptionEntityData(t, rawWebhookURL)

	actual := NewAuditLog(&eventproto.Event{
		Id:                 "auditlog-1",
		EntityType:         eventproto.Event_SUBSCRIPTION,
		EntityId:           "sub-1",
		Type:               eventproto.Event_SUBSCRIPTION_CREATED,
		Editor:             &eventproto.Editor{Email: "bucketeer@bucketeer.io"},
		Data:               createdEvent,
		EntityData:         entityData,
		PreviousEntityData: entityData,
	}, "env-1")

	assert.NotContains(t, actual.EntityData, rawWebhookURL)
	assert.Contains(t, actual.EntityData, maskedWebhookURL)
	assert.NotContains(t, actual.PreviousEntityData, rawWebhookURL)
	assert.Contains(t, actual.PreviousEntityData, maskedWebhookURL)
	created := &eventproto.SubscriptionCreatedEvent{}
	require.NoError(t, actual.Event.UnmarshalTo(created))
	assert.Equal(t, maskedWebhookURL, created.Recipient.SlackChannelRecipient.WebhookUrl)
}
