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

package v2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/encryption"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	v2ns "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

const (
	testEnvironmentID = "env-0"
	testWebhookURL    = "https://hooks.slack.com/services/T000/B000/secret-token"
)

func newTestEncryptor(t *testing.T) encryption.Encryptor {
	t.Helper()
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	enc, err := encryption.NewAEADLocal(context.Background(), "key-1", key)
	require.NoError(t, err)
	return enc
}

func encryptForTest(t *testing.T, enc encryption.Encryptor, plaintext string) string {
	t.Helper()
	sealed, err := enc.Encrypt(context.Background(), plaintext)
	require.NoError(t, err)
	return sealed
}

func newSubscription(id, webhookURL string) *proto.Subscription {
	return &proto.Subscription{
		Id:   id,
		Name: "sub",
		Recipient: &proto.Recipient{
			Type:                  proto.Recipient_SlackChannel,
			SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: webhookURL},
			Language:              proto.Recipient_JAPANESE,
		},
	}
}

func webhookURLOf(s *proto.Subscription) string {
	return s.GetRecipient().GetSlackChannelRecipient().GetWebhookUrl()
}

type writePattern struct {
	desc      string
	encryptor encryption.Encryptor
	input     *proto.Subscription
	check     func(t *testing.T, stored *proto.Subscription)
}

func writePatterns(t *testing.T, ctx context.Context, enc encryption.Encryptor) []writePattern {
	t.Helper()
	alreadySealed := encryptForTest(t, enc, testWebhookURL)
	return []writePattern{
		{
			desc:      "plaintext url is encrypted before it is stored",
			encryptor: enc,
			input:     newSubscription("id-0", testWebhookURL),
			check: func(t *testing.T, stored *proto.Subscription) {
				url := webhookURLOf(stored)
				assert.True(t, encryption.IsEncrypted(url))
				assert.NotContains(t, url, "hooks.slack.com")
				assert.NotContains(t, url, "secret-token")
				opened, err := enc.Decrypt(ctx, url)
				require.NoError(t, err)
				assert.Equal(t, testWebhookURL, opened)
				// The rest of the recipient is kept.
				assert.Equal(t, proto.Recipient_SlackChannel, stored.Recipient.Type)
				assert.Equal(t, proto.Recipient_JAPANESE, stored.Recipient.Language)
			},
		},
		{
			desc:      "an already encrypted url is stored as it is",
			encryptor: enc,
			input:     newSubscription("id-0", alreadySealed),
			check: func(t *testing.T, stored *proto.Subscription) {
				assert.Equal(t, alreadySealed, webhookURLOf(stored))
			},
		},
		{
			desc:      "a subscription without a recipient is stored as it is",
			encryptor: enc,
			input:     &proto.Subscription{Id: "id-0", Name: "sub"},
			check: func(t *testing.T, stored *proto.Subscription) {
				assert.Nil(t, stored.Recipient)
			},
		},
		{
			desc:      "none keeps the plaintext url",
			encryptor: encryption.NewNone(),
			input:     newSubscription("id-0", testWebhookURL),
			check: func(t *testing.T, stored *proto.Subscription) {
				assert.Equal(t, testWebhookURL, webhookURLOf(stored))
			},
		},
	}
}

type readPattern struct {
	desc        string
	encryptor   encryption.Encryptor
	stored      string
	expected    string
	expectedErr error
}

func readPatterns(t *testing.T, ctx context.Context, enc encryption.Encryptor) []readPattern {
	t.Helper()
	sealed := encryptForTest(t, enc, testWebhookURL)
	otherKeyEnc, err := encryption.NewAEADLocal(ctx, "key-2", make([]byte, 32))
	require.NoError(t, err)
	return []readPattern{
		{
			desc:      "an encrypted url is decrypted",
			encryptor: enc,
			stored:    sealed,
			expected:  testWebhookURL,
		},
		{
			desc:      "a legacy plaintext url is returned unchanged",
			encryptor: enc,
			stored:    testWebhookURL,
			expected:  testWebhookURL,
		},
		{
			desc:      "none returns the plaintext url unchanged",
			encryptor: encryption.NewNone(),
			stored:    testWebhookURL,
			expected:  testWebhookURL,
		},
		{
			desc:        "none refuses to return an encrypted url",
			encryptor:   encryption.NewNone(),
			stored:      sealed,
			expectedErr: encryption.ErrNotConfigured,
		},
		{
			desc:        "a url encrypted with another key is an error",
			encryptor:   otherKeyEnc,
			stored:      sealed,
			expectedErr: encryption.ErrUnknownKeyID,
		},
	}
}

func TestEncryptedSubscriptionStorageWrite(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, p := range writePatterns(t, ctx, newTestEncryptor(t)) {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			inner := mock.NewMockSubscriptionStorage(mockController)
			var created, updated *proto.Subscription
			inner.EXPECT().CreateSubscription(gomock.Any(), gomock.Any(), testEnvironmentID).DoAndReturn(
				func(_ context.Context, e *domain.Subscription, _ string) error {
					created = e.Subscription
					return nil
				},
			)
			inner.EXPECT().UpdateSubscription(gomock.Any(), gomock.Any(), testEnvironmentID).DoAndReturn(
				func(_ context.Context, e *domain.Subscription, _ string) error {
					updated = e.Subscription
					return nil
				},
			)
			storage := v2ns.NewEncryptedSubscriptionStorage(inner, p.encryptor)

			input := &domain.Subscription{Subscription: p.input}
			original := webhookURLOf(p.input)
			require.NoError(t, storage.CreateSubscription(ctx, input, testEnvironmentID))
			require.NoError(t, storage.UpdateSubscription(ctx, input, testEnvironmentID))

			p.check(t, created)
			p.check(t, updated)
			// The caller's entity is never modified.
			assert.Equal(t, original, webhookURLOf(input.Subscription))
		})
	}
}

func TestEncryptedSubscriptionStorageRead(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, p := range readPatterns(t, ctx, newTestEncryptor(t)) {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			inner := mock.NewMockSubscriptionStorage(mockController)
			inner.EXPECT().GetSubscription(gomock.Any(), "id-0", testEnvironmentID).Return(
				&domain.Subscription{Subscription: newSubscription("id-0", p.stored)}, nil,
			)
			inner.EXPECT().ListSubscriptions(gomock.Any(), gomock.Any()).Return(
				[]*proto.Subscription{
					newSubscription("id-0", p.stored),
					newSubscription("id-1", p.stored),
					{Id: "id-2", Name: "no recipient"},
				}, 3, int64(3), nil,
			)
			storage := v2ns.NewEncryptedSubscriptionStorage(inner, p.encryptor)

			subscription, err := storage.GetSubscription(ctx, "id-0", testEnvironmentID)
			subscriptions, nextOffset, totalCount, listErr := storage.ListSubscriptions(
				ctx,
				v2ns.ListSubscriptionsParams{EnvironmentIDs: []string{testEnvironmentID}},
			)
			if p.expectedErr != nil {
				assert.ErrorIs(t, err, p.expectedErr)
				assert.Nil(t, subscription)
				assert.ErrorIs(t, listErr, p.expectedErr)
				assert.Nil(t, subscriptions)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, p.expected, webhookURLOf(subscription.Subscription))
			require.NoError(t, listErr)
			assert.Equal(t, 3, nextOffset)
			assert.Equal(t, int64(3), totalCount)
			require.Len(t, subscriptions, 3)
			assert.Equal(t, p.expected, webhookURLOf(subscriptions[0]))
			assert.Equal(t, p.expected, webhookURLOf(subscriptions[1]))
			assert.Nil(t, subscriptions[2].Recipient)
		})
	}
}

func TestEncryptedSubscriptionStoragePassesThroughErrorsAndDelete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	inner := mock.NewMockSubscriptionStorage(mockController)
	storage := v2ns.NewEncryptedSubscriptionStorage(inner, newTestEncryptor(t))

	inner.EXPECT().GetSubscription(gomock.Any(), "id-0", testEnvironmentID).Return(nil, v2ns.ErrSubscriptionNotFound)
	subscription, err := storage.GetSubscription(ctx, "id-0", testEnvironmentID)
	assert.ErrorIs(t, err, v2ns.ErrSubscriptionNotFound)
	assert.Nil(t, subscription)

	inner.EXPECT().ListSubscriptions(gomock.Any(), gomock.Any()).Return(nil, 0, int64(0), v2ns.ErrInvalidCursor)
	subscriptions, _, _, err := storage.ListSubscriptions(ctx, v2ns.ListSubscriptionsParams{})
	assert.ErrorIs(t, err, v2ns.ErrInvalidCursor)
	assert.Nil(t, subscriptions)

	inner.EXPECT().CreateSubscription(gomock.Any(), gomock.Any(), testEnvironmentID).Return(
		v2ns.ErrSubscriptionAlreadyExists,
	)
	err = storage.CreateSubscription(
		ctx,
		&domain.Subscription{Subscription: newSubscription("id-0", testWebhookURL)},
		testEnvironmentID,
	)
	assert.ErrorIs(t, err, v2ns.ErrSubscriptionAlreadyExists)

	inner.EXPECT().DeleteSubscription(gomock.Any(), "id-0", testEnvironmentID).Return(nil)
	assert.NoError(t, storage.DeleteSubscription(ctx, "id-0", testEnvironmentID))
}

func TestEncryptedAdminSubscriptionStorageWrite(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, p := range writePatterns(t, ctx, newTestEncryptor(t)) {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			inner := mock.NewMockAdminSubscriptionStorage(mockController)
			var created, updated *proto.Subscription
			inner.EXPECT().CreateAdminSubscription(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, e *domain.Subscription) error {
					created = e.Subscription
					return nil
				},
			)
			inner.EXPECT().UpdateAdminSubscription(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, e *domain.Subscription) error {
					updated = e.Subscription
					return nil
				},
			)
			storage := v2ns.NewEncryptedAdminSubscriptionStorage(inner, p.encryptor)

			input := &domain.Subscription{Subscription: p.input}
			original := webhookURLOf(p.input)
			require.NoError(t, storage.CreateAdminSubscription(ctx, input))
			require.NoError(t, storage.UpdateAdminSubscription(ctx, input))

			p.check(t, created)
			p.check(t, updated)
			assert.Equal(t, original, webhookURLOf(input.Subscription))
		})
	}
}

func TestEncryptedAdminSubscriptionStorageRead(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, p := range readPatterns(t, ctx, newTestEncryptor(t)) {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			inner := mock.NewMockAdminSubscriptionStorage(mockController)
			inner.EXPECT().GetAdminSubscription(gomock.Any(), "id-0").Return(
				&domain.Subscription{Subscription: newSubscription("id-0", p.stored)}, nil,
			)
			inner.EXPECT().ListAdminSubscriptions(gomock.Any(), gomock.Any()).Return(
				[]*proto.Subscription{newSubscription("id-0", p.stored)}, 1, int64(1), nil,
			)
			storage := v2ns.NewEncryptedAdminSubscriptionStorage(inner, p.encryptor)

			subscription, err := storage.GetAdminSubscription(ctx, "id-0")
			subscriptions, nextOffset, totalCount, listErr := storage.ListAdminSubscriptions(
				ctx,
				v2ns.ListAdminSubscriptionsParams{},
			)
			if p.expectedErr != nil {
				assert.ErrorIs(t, err, p.expectedErr)
				assert.Nil(t, subscription)
				assert.ErrorIs(t, listErr, p.expectedErr)
				assert.Nil(t, subscriptions)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, p.expected, webhookURLOf(subscription.Subscription))
			require.NoError(t, listErr)
			assert.Equal(t, 1, nextOffset)
			assert.Equal(t, int64(1), totalCount)
			require.Len(t, subscriptions, 1)
			assert.Equal(t, p.expected, webhookURLOf(subscriptions[0]))
		})
	}
}

func TestEncryptedAdminSubscriptionStoragePassesThroughErrorsAndDelete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	inner := mock.NewMockAdminSubscriptionStorage(mockController)
	storage := v2ns.NewEncryptedAdminSubscriptionStorage(inner, newTestEncryptor(t))

	inner.EXPECT().GetAdminSubscription(gomock.Any(), "id-0").Return(nil, v2ns.ErrAdminSubscriptionNotFound)
	subscription, err := storage.GetAdminSubscription(ctx, "id-0")
	assert.ErrorIs(t, err, v2ns.ErrAdminSubscriptionNotFound)
	assert.Nil(t, subscription)

	inner.EXPECT().ListAdminSubscriptions(gomock.Any(), gomock.Any()).Return(nil, 0, int64(0), v2ns.ErrInvalidCursor)
	subscriptions, _, _, err := storage.ListAdminSubscriptions(ctx, v2ns.ListAdminSubscriptionsParams{})
	assert.ErrorIs(t, err, v2ns.ErrInvalidCursor)
	assert.Nil(t, subscriptions)

	inner.EXPECT().DeleteAdminSubscription(gomock.Any(), "id-0").Return(nil)
	assert.NoError(t, storage.DeleteAdminSubscription(ctx, "id-0"))
}
