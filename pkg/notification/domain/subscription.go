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

package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/wrapperspb"

	err "github.com/bucketeer-io/bucketeer/pkg/error"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

var (
	ErrUnknownRecipient = err.NewErrorInvalidArgUnknown(
		err.NotificationPackageName,
		"unknown recipient",
		"recipient",
	)
	ErrSourceTypesMustHaveAtLeastOne = err.NewErrorInvalidArgNotMatchFormat(
		err.NotificationPackageName,
		"notification types must have at least one",
		"notification_types",
	)
	ErrSourceTypeNotFound          = err.NewErrorNotFound(err.NotificationPackageName, "notification not found", "notification_type")
	ErrAlreadyEnabled              = err.NewErrorAlreadyExists(err.NotificationPackageName, "already enabled")
	ErrAlreadyDisabled             = err.NewErrorAlreadyExists(err.NotificationPackageName, "already disabled")
	ErrCannotUpdateFeatureFlagTags = err.NewErrorNotFound(
		err.NotificationPackageName,
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
