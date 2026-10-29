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
	"errors"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	subscriptiondomain "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	domainevent "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

// JSON field names of the subscription entity data, as encoded by encoding/json
// from the proto.Subscription struct tags.
const (
	recipientJSONField             = "recipient"
	slackChannelRecipientJSONField = "slack_channel_recipient"
	webhookURLJSONField            = "webhook_url"
)

// MaskWebhookURL masks the webhook URL of a subscription audit log, in place.
// The subscription service masks it before publishing the domain event; this
// covers the rows saved before that was the case.
func MaskWebhookURL(auditlog *proto.AuditLog) error {
	if auditlog == nil || auditlog.EntityType != domainevent.Event_SUBSCRIPTION {
		return nil
	}
	entityData, entityErr := MaskWebhookURLEntityData(auditlog.EntityData)
	auditlog.EntityData = entityData
	previousEntityData, previousErr := MaskWebhookURLEntityData(auditlog.PreviousEntityData)
	auditlog.PreviousEntityData = previousEntityData
	event, eventErr := MaskWebhookURLEventData(auditlog.Event)
	auditlog.Event = event
	return errors.Join(entityErr, previousErr, eventErr)
}

// MaskWebhookURLEntityData masks the webhook URL of the JSON encoded subscription entity data.
// It returns an empty string on failure, so a raw URL is never stored nor returned.
func MaskWebhookURLEntityData(data string) (string, error) {
	if data == "" {
		return data, nil
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	// Avoid converting the numbers to scientific notation.
	decoder.UseNumber()
	entity := make(map[string]interface{})
	if err := decoder.Decode(&entity); err != nil {
		return "", err
	}
	recipient, ok := entity[recipientJSONField].(map[string]interface{})
	if !ok {
		return data, nil
	}
	slackRecipient, ok := recipient[slackChannelRecipientJSONField].(map[string]interface{})
	if !ok {
		return data, nil
	}
	webhookURL, ok := slackRecipient[webhookURLJSONField].(string)
	if !ok {
		return data, nil
	}
	slackRecipient[webhookURLJSONField] = subscriptiondomain.MaskWebhookURL(webhookURL)
	masked, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return "", err
	}
	return string(masked), nil
}

// MaskWebhookURLEventData handles the created event, the only one containing the recipient.
// It returns nil on failure, so a raw URL is never stored nor returned.
func MaskWebhookURLEventData(event *anypb.Any) (*anypb.Any, error) {
	if event == nil || !event.MessageIs(&domainevent.SubscriptionCreatedEvent{}) {
		return event, nil
	}
	created := &domainevent.SubscriptionCreatedEvent{}
	if err := event.UnmarshalTo(created); err != nil {
		return nil, err
	}
	if created.Recipient == nil || created.Recipient.SlackChannelRecipient == nil {
		return event, nil
	}
	created.Recipient.SlackChannelRecipient.WebhookUrl = subscriptiondomain.MaskWebhookURL(
		created.Recipient.SlackChannelRecipient.WebhookUrl,
	)
	masked, err := anypb.New(created)
	if err != nil {
		return nil, err
	}
	return masked, nil
}
