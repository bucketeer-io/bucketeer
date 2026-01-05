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

package api

import (
	notificationproto "github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

func (s *NotificationService) validateCreateSubscriptionRequest(
	req *notificationproto.CreateSubscriptionRequest,
) error {
	if req.Name == "" {
		return statusNameRequired.Err()
	}
	if len(req.SourceTypes) == 0 {
		return statusSourceTypesRequired.Err()
	}
	// We should only save the tags if the Feature source type is in the list
	if req.FeatureFlagTags != nil {
		var found bool
		for _, sourceType := range req.SourceTypes {
			if sourceType == notificationproto.Subscription_DOMAIN_EVENT_FEATURE {
				found = true
				break
			}
		}
		if !found {
			return statusSourceTypesRequired.Err()
		}
	}
	if err := s.validateRecipient(req.Recipient); err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) validateRecipient(
	recipient *notificationproto.Recipient,
) error {
	if recipient == nil {
		return statusRecipientRequired.Err()
	}
	if recipient.Type == notificationproto.Recipient_SlackChannel {
		return s.validateSlackRecipient(recipient.SlackChannelRecipient)
	}
	return statusUnknownRecipient.Err()
}

func (s *NotificationService) validateSlackRecipient(
	sr *notificationproto.SlackChannelRecipient,
) error {
	// TODO: Check ping to the webhook URL?
	if sr == nil {
		return statusSlackRecipientRequired.Err()
	}
	if sr.WebhookUrl == "" {
		return statusSlackRecipientWebhookURLRequired.Err()
	}
	return nil
}

func (s *NotificationService) validateUpdateSubscriptionRequest(
	req *notificationproto.UpdateSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	if req.Name != nil && req.Name.Value == "" {
		return statusNameRequired.Err()
	}
	return nil
}

func validateDeleteSubscriptionRequest(
	req *notificationproto.DeleteSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	return nil
}

func validateGetSubscriptionRequest(req *notificationproto.GetSubscriptionRequest) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	return nil
}
