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

package v2

import (
	"context"

	pb "google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/encryption"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

// encryptedSubscriptionStorage protects the Slack webhook URL of a recipient
// at rest. It encrypts the URL before it is written and decrypts it after it
// is read, so callers keep working with the plaintext URL whatever the
// configured backend.
type encryptedSubscriptionStorage struct {
	inner     SubscriptionStorage
	encryptor encryption.Encryptor
}

// NewEncryptedSubscriptionStorage wraps inner so that the recipient webhook
// URL is stored through encryptor.
func NewEncryptedSubscriptionStorage(
	inner SubscriptionStorage,
	encryptor encryption.Encryptor,
) SubscriptionStorage {
	return &encryptedSubscriptionStorage{inner: inner, encryptor: encryptor}
}

func (s *encryptedSubscriptionStorage) CreateSubscription(
	ctx context.Context,
	e *domain.Subscription,
	environmentId string,
) error {
	sealed, err := sealRecipient(ctx, s.encryptor, e)
	if err != nil {
		return err
	}
	return s.inner.CreateSubscription(ctx, sealed, environmentId)
}

func (s *encryptedSubscriptionStorage) UpdateSubscription(
	ctx context.Context,
	e *domain.Subscription,
	environmentId string,
) error {
	sealed, err := sealRecipient(ctx, s.encryptor, e)
	if err != nil {
		return err
	}
	return s.inner.UpdateSubscription(ctx, sealed, environmentId)
}

func (s *encryptedSubscriptionStorage) DeleteSubscription(ctx context.Context, id, environmentId string) error {
	return s.inner.DeleteSubscription(ctx, id, environmentId)
}

func (s *encryptedSubscriptionStorage) GetSubscription(
	ctx context.Context,
	id, environmentId string,
) (*domain.Subscription, error) {
	subscription, err := s.inner.GetSubscription(ctx, id, environmentId)
	if err != nil {
		return nil, err
	}
	if err := openRecipient(ctx, s.encryptor, subscription.Subscription); err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *encryptedSubscriptionStorage) ListSubscriptions(
	ctx context.Context,
	params ListSubscriptionsParams,
) ([]*proto.Subscription, int, int64, error) {
	subscriptions, nextOffset, totalCount, err := s.inner.ListSubscriptions(ctx, params)
	if err != nil {
		return nil, 0, 0, err
	}
	for _, subscription := range subscriptions {
		if err := openRecipient(ctx, s.encryptor, subscription); err != nil {
			return nil, 0, 0, err
		}
	}
	return subscriptions, nextOffset, totalCount, nil
}

// encryptedAdminSubscriptionStorage is the admin subscription counterpart of
// encryptedSubscriptionStorage.
type encryptedAdminSubscriptionStorage struct {
	inner     AdminSubscriptionStorage
	encryptor encryption.Encryptor
}

// NewEncryptedAdminSubscriptionStorage wraps inner so that the recipient
// webhook URL is stored through encryptor.
func NewEncryptedAdminSubscriptionStorage(
	inner AdminSubscriptionStorage,
	encryptor encryption.Encryptor,
) AdminSubscriptionStorage {
	return &encryptedAdminSubscriptionStorage{inner: inner, encryptor: encryptor}
}

func (s *encryptedAdminSubscriptionStorage) CreateAdminSubscription(
	ctx context.Context,
	e *domain.Subscription,
) error {
	sealed, err := sealRecipient(ctx, s.encryptor, e)
	if err != nil {
		return err
	}
	return s.inner.CreateAdminSubscription(ctx, sealed)
}

func (s *encryptedAdminSubscriptionStorage) UpdateAdminSubscription(
	ctx context.Context,
	e *domain.Subscription,
) error {
	sealed, err := sealRecipient(ctx, s.encryptor, e)
	if err != nil {
		return err
	}
	return s.inner.UpdateAdminSubscription(ctx, sealed)
}

func (s *encryptedAdminSubscriptionStorage) DeleteAdminSubscription(ctx context.Context, id string) error {
	return s.inner.DeleteAdminSubscription(ctx, id)
}

func (s *encryptedAdminSubscriptionStorage) GetAdminSubscription(
	ctx context.Context,
	id string,
) (*domain.Subscription, error) {
	subscription, err := s.inner.GetAdminSubscription(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := openRecipient(ctx, s.encryptor, subscription.Subscription); err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *encryptedAdminSubscriptionStorage) ListAdminSubscriptions(
	ctx context.Context,
	params ListAdminSubscriptionsParams,
) ([]*proto.Subscription, int, int64, error) {
	subscriptions, nextOffset, totalCount, err := s.inner.ListAdminSubscriptions(ctx, params)
	if err != nil {
		return nil, 0, 0, err
	}
	for _, subscription := range subscriptions {
		if err := openRecipient(ctx, s.encryptor, subscription); err != nil {
			return nil, 0, 0, err
		}
	}
	return subscriptions, nextOffset, totalCount, nil
}

// sealRecipient returns a copy of e whose Slack webhook URL is encrypted. The
// caller's entity is left untouched. A URL that is already encrypted is stored
// as it is so that a re-saved entity is never encrypted twice.
func sealRecipient(
	ctx context.Context,
	encryptor encryption.Encryptor,
	e *domain.Subscription,
) (*domain.Subscription, error) {
	webhookURL := e.GetRecipient().GetSlackChannelRecipient().GetWebhookUrl()
	if webhookURL == "" || encryption.IsEncrypted(webhookURL) {
		return e, nil
	}
	sealed, err := encryptor.Encrypt(ctx, webhookURL)
	if err != nil {
		return nil, err
	}
	copied := pb.Clone(e.Subscription).(*proto.Subscription)
	copied.Recipient.SlackChannelRecipient.WebhookUrl = sealed
	return &domain.Subscription{Subscription: copied}, nil
}

// openRecipient replaces the stored Slack webhook URL of subscription with its
// plaintext. Legacy plaintext values are returned unchanged by the encryptor.
func openRecipient(ctx context.Context, encryptor encryption.Encryptor, subscription *proto.Subscription) error {
	slackRecipient := subscription.GetRecipient().GetSlackChannelRecipient()
	if slackRecipient == nil || slackRecipient.WebhookUrl == "" {
		return nil
	}
	plaintext, err := encryptor.Decrypt(ctx, slackRecipient.WebhookUrl)
	if err != nil {
		return err
	}
	slackRecipient.WebhookUrl = plaintext
	return nil
}
