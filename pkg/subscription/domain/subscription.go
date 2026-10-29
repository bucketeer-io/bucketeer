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
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	pb "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

// Slack incoming webhooks always have the form
// https://hooks.slack.com/services/<team>/<channel>/<token>.
const (
	slackWebhookScheme     = "https"
	slackWebhookHost       = "hooks.slack.com"
	slackWebhookPathPrefix = "/services/"
)

// webhookURLMask replaces the path of a webhook URL when it is returned by the
// API or recorded in the audit log. The path is the credential: anyone holding
// the full URL can post to the channel.
const webhookURLMask = "****"

var (
	ErrUnknownRecipient = err.NewErrorInvalidArgUnknown(
		err.SubscriptionPackageName,
		"unknown recipient",
		"recipient",
	)
	ErrSlackWebhookURLInvalid = err.NewErrorInvalidArgNotMatchFormat(
		err.SubscriptionPackageName,
		"webhook URL must be a Slack incoming webhook URL (https://hooks.slack.com/services/...)",
		"webhook_url",
	)
	ErrSourceTypesMustHaveAtLeastOne = err.NewErrorInvalidArgNotMatchFormat(
		err.SubscriptionPackageName,
		"notification types must have at least one",
		"notification_types",
	)
	ErrSourceTypeNotFound = err.NewErrorNotFound(
		err.SubscriptionPackageName,
		"notification not found", "notification_type",
	)
	ErrAlreadyEnabled              = err.NewErrorAlreadyExists(err.SubscriptionPackageName, "already enabled")
	ErrAlreadyDisabled             = err.NewErrorAlreadyExists(err.SubscriptionPackageName, "already disabled")
	ErrCannotUpdateFeatureFlagTags = err.NewErrorNotFound(
		err.SubscriptionPackageName,
		"cannot update the feature flag tags when there is feature source type",
		"feature_source_type",
	)
)

type Subscription struct {
	*proto.Subscription
}

func NewSubscription(
	name string,
	sourceTypes []proto.Subscription_SourceType,
	recipient *proto.Recipient,
	featureFlagTags []string,
) (*Subscription, error) {

	sid, err := ID(recipient)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	s := &Subscription{&proto.Subscription{
		Id:              sid,
		Name:            name,
		SourceTypes:     sourceTypes,
		Recipient:       recipient,
		FeatureFlagTags: featureFlagTags,
		CreatedAt:       now,
		UpdatedAt:       now,
	}}
	return s, nil
}

func ID(recipient *proto.Recipient) (string, error) {
	if recipient.Type == proto.Recipient_SlackChannel {
		return SlackChannelRecipientID(recipient.SlackChannelRecipient.WebhookUrl), nil
	}
	return "", ErrUnknownRecipient
}

func SlackChannelRecipientID(webhookURL string) string {
	hashed := sha256.Sum256([]byte(webhookURL))
	return hex.EncodeToString(hashed[:])
}

// ValidateSlackWebhookURL ensures the URL is a Slack incoming webhook:
// the https scheme, exactly the hooks.slack.com host (no credentials, no port)
// and a non-empty /services/ path. Notifications are only ever delivered to
// Slack, so any other destination is rejected.
func ValidateSlackWebhookURL(webhookURL string) error {
	u, e := url.Parse(webhookURL)
	if e != nil {
		return ErrSlackWebhookURLInvalid
	}
	if u.Scheme != slackWebhookScheme || u.User != nil || u.Host != slackWebhookHost {
		return ErrSlackWebhookURLInvalid
	}
	if !strings.HasPrefix(u.Path, slackWebhookPathPrefix) || len(u.Path) == len(slackWebhookPathPrefix) {
		return ErrSlackWebhookURLInvalid
	}
	return nil
}

// MaskWebhookURL hides the credential part of a webhook URL, keeping only the
// scheme and host so the destination stays recognizable
// (https://hooks.slack.com/services/T000/B000/token -> https://hooks.slack.com/****).
// Values that are not URLs are replaced entirely, and masking an already masked
// URL leaves it unchanged. It must be used everywhere a stored webhook URL is
// returned or recorded; only the notification sender needs the original value.
func MaskWebhookURL(webhookURL string) string {
	if webhookURL == "" {
		return ""
	}
	u, e := url.Parse(webhookURL)
	if e != nil || u.Scheme == "" || u.Host == "" {
		return webhookURLMask
	}
	return u.Scheme + "://" + u.Host + "/" + webhookURLMask
}

// MaskSubscriptionWebhookURL masks the webhook URL of the subscription in place.
// Use it on subscriptions that are only returned, never persisted afterwards.
func MaskSubscriptionWebhookURL(subscription *proto.Subscription) {
	if subscription == nil || subscription.Recipient == nil || subscription.Recipient.SlackChannelRecipient == nil {
		return
	}
	subscription.Recipient.SlackChannelRecipient.WebhookUrl = MaskWebhookURL(
		subscription.Recipient.SlackChannelRecipient.WebhookUrl,
	)
}

// CopyWithMaskedWebhookURL returns a deep copy of the subscription whose webhook
// URL is masked, leaving the original untouched for persistence.
func CopyWithMaskedWebhookURL(subscription *proto.Subscription) *proto.Subscription {
	if subscription == nil {
		return nil
	}
	masked := pb.Clone(subscription).(*proto.Subscription)
	MaskSubscriptionWebhookURL(masked)
	return masked
}

func (s *Subscription) UpdateSubscription(
	name *wrapperspb.StringValue,
	sourceTypes []proto.Subscription_SourceType,
	disabled *wrapperspb.BoolValue,
	featureFlagTags []string,
) (*Subscription, error) {
	updated := &Subscription{}
	if err := copier.Copy(updated, s); err != nil {
		return nil, err
	}

	if name != nil {
		updated.Name = name.Value
	}
	if len(sourceTypes) > 0 {
		// We must check the feature source type is being deleted.
		// If so, we must reset the tags
		for _, sourceType := range updated.SourceTypes {
			if sourceType == proto.Subscription_DOMAIN_EVENT_FEATURE {
				updated.FeatureFlagTags = []string{}
				break
			}
		}
		updated.SourceTypes = sourceTypes
	}
	if disabled != nil {
		updated.Disabled = disabled.Value
	}
	// The tags updating must be updated after the source types updating
	if featureFlagTags != nil {
		var found bool
		for _, sourceType := range updated.SourceTypes {
			// Because the feature flag tags belong to the feature domain event
			// we must ensure the feature source type is in the list.
			if sourceType == proto.Subscription_DOMAIN_EVENT_FEATURE {
				found = true
				break
			}
		}
		if !found {
			return nil, ErrCannotUpdateFeatureFlagTags
		}
		updated.FeatureFlagTags = featureFlagTags
	}
	updated.UpdatedAt = time.Now().Unix()
	return updated, nil
}

func (s *Subscription) Enable() error {
	if !s.Disabled {
		return ErrAlreadyEnabled
	}
	s.Disabled = false
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Subscription) Disable() error {
	if s.Disabled {
		return ErrAlreadyDisabled
	}
	s.Disabled = true
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Subscription) Rename(name string) error {
	s.Name = name
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Subscription) AddSourceTypes(sourceTypes []proto.Subscription_SourceType) error {
	for _, nt := range sourceTypes {
		if containsSourceType(nt, s.SourceTypes) {
			continue
		}
		s.SourceTypes = append(s.SourceTypes, nt)
	}
	sortSourceType(s.SourceTypes)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Subscription) DeleteSourceTypes(sourceTypes []proto.Subscription_SourceType) error {
	if len(s.SourceTypes) <= 1 {
		return ErrSourceTypesMustHaveAtLeastOne
	}
	for _, nt := range sourceTypes {
		// Reset the tags if the feature source type is being deleted
		if nt == proto.Subscription_DOMAIN_EVENT_FEATURE {
			s.FeatureFlagTags = []string{}
		}
		idx, err := indexSourceType(nt, s.SourceTypes)
		if err != nil {
			return err
		}
		s.SourceTypes = append(s.SourceTypes[:idx], s.SourceTypes[idx+1:]...)
	}
	sortSourceType(s.SourceTypes)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

// The tags updating must be updated after the source types updating
func (s *Subscription) UpdateFeatureFlagTags(tags []string) error {
	var found bool
	for _, sourceType := range s.SourceTypes {
		// Because the feature flag tags belong to the feature domain event
		// we must ensure the feature source type is in the list.
		if sourceType == proto.Subscription_DOMAIN_EVENT_FEATURE {
			found = true
			break
		}
	}
	if !found {
		return ErrCannotUpdateFeatureFlagTags
	}
	s.FeatureFlagTags = tags
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func indexSourceType(needle proto.Subscription_SourceType, haystack []proto.Subscription_SourceType) (int, error) {
	for i := range haystack {
		if haystack[i] == needle {
			return i, nil
		}
	}
	return -1, ErrSourceTypeNotFound
}

func containsSourceType(needle proto.Subscription_SourceType, haystack []proto.Subscription_SourceType) bool {
	for i := range haystack {
		if haystack[i] == needle {
			return true
		}
	}
	return false
}

func sortSourceType(sourceTypes []proto.Subscription_SourceType) {
	sort.Slice(sourceTypes, func(i, j int) bool {
		return sourceTypes[i] < sourceTypes[j]
	})
}
