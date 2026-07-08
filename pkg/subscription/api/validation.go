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
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

func (s *SubscriptionService) validateCreateSubscriptionRequest(
	req *subscriptionproto.CreateSubscriptionRequest,
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
			if sourceType == subscriptionproto.Subscription_DOMAIN_EVENT_FEATURE {
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

func (s *SubscriptionService) validateRecipient(
	recipient *subscriptionproto.Recipient,
) error {
	if recipient == nil {
		return statusRecipientRequired.Err()
	}
	if recipient.Type == subscriptionproto.Recipient_SlackChannel {
		return s.validateSlackRecipient(recipient.SlackChannelRecipient)
	}
	return statusUnknownRecipient.Err()
}

func (s *SubscriptionService) validateSlackRecipient(
	sr *subscriptionproto.SlackChannelRecipient,
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

func (s *SubscriptionService) validateUpdateSubscriptionRequest(
	req *subscriptionproto.UpdateSubscriptionRequest,
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
	req *subscriptionproto.DeleteSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	return nil
}

func validateGetSubscriptionRequest(req *subscriptionproto.GetSubscriptionRequest) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	return nil
}
