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

package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"time"

	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

var (
	ErrUnknownRecipient              = errors.New("subscription: unknown recipient")
	ErrSourceTypesMustHaveAtLeastOne = errors.New("subscription: notification types must have at least one")
	ErrSourceTypeNotFound            = errors.New("subscription: notification not found")
	ErrAlreadyEnabled                = errors.New("subscription: already enabled")
	ErrAlreadyDisabled               = errors.New("subscription: already disabled")
)

type Subscription struct {
	*proto.Subscription
}

func NewSubscription(
	name string,
	sourceTypes []proto.Subscription_SourceType,
	recipient *proto.Recipient) (*Subscription, error) {

	sid, err := ID(recipient)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	s := &Subscription{&proto.Subscription{
		Name:        name,
		CreatedAt:   now,
		UpdatedAt:   now,
		Id:          sid,
		SourceTypes: sourceTypes,
		Recipient:   recipient,
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
