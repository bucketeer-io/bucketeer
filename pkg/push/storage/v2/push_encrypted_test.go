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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/encryption"
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	v2ps "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

const (
	testEnvironmentID = "env-0"
	testServiceAcct   = `{"type":"service_account","project_id":"p","private_key":"-----BEGIN PRIVATE KEY-----\nMIIB"}`
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

func TestEncryptedPushStorageWrite(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	enc := newTestEncryptor(t)
	alreadySealed := encryptForTest(t, enc, testServiceAcct)

	patterns := []struct {
		desc      string
		encryptor encryption.Encryptor
		input     string
		check     func(t *testing.T, stored string)
	}{
		{
			desc:      "plaintext is encrypted before it is stored",
			encryptor: enc,
			input:     testServiceAcct,
			check: func(t *testing.T, stored string) {
				assert.True(t, encryption.IsEncrypted(stored))
				assert.NotContains(t, stored, "service_account")
				assert.NotContains(t, stored, "PRIVATE KEY")
				opened, err := enc.Decrypt(ctx, stored)
				require.NoError(t, err)
				assert.Equal(t, testServiceAcct, opened)
			},
		},
		{
			desc:      "an already encrypted value is stored as it is",
			encryptor: enc,
			input:     alreadySealed,
			check: func(t *testing.T, stored string) {
				assert.Equal(t, alreadySealed, stored)
			},
		},
		{
			desc:      "an empty value stays empty",
			encryptor: enc,
			input:     "",
			check: func(t *testing.T, stored string) {
				assert.Equal(t, "", stored)
			},
		},
		{
			desc:      "none keeps the plaintext",
			encryptor: encryption.NewNone(),
			input:     testServiceAcct,
			check: func(t *testing.T, stored string) {
				assert.Equal(t, testServiceAcct, stored)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			inner := mock.NewMockPushStorage(mockController)
			var created, updated string
			inner.EXPECT().CreatePush(gomock.Any(), gomock.Any(), testEnvironmentID).DoAndReturn(
				func(_ context.Context, e *domain.Push, _ string) error {
					created = e.FcmServiceAccount
					return nil
				},
			)
			inner.EXPECT().UpdatePush(gomock.Any(), gomock.Any(), testEnvironmentID).DoAndReturn(
				func(_ context.Context, e *domain.Push, _ string) error {
					updated = e.FcmServiceAccount
					return nil
				},
			)
			storage := v2ps.NewEncryptedPushStorage(inner, p.encryptor)

			push := &domain.Push{Push: &proto.Push{Id: "id-0", Name: "push", FcmServiceAccount: p.input}}
			require.NoError(t, storage.CreatePush(ctx, push, testEnvironmentID))
			require.NoError(t, storage.UpdatePush(ctx, push, testEnvironmentID))

			p.check(t, created)
			p.check(t, updated)
			// The caller's entity is never modified.
			assert.Equal(t, p.input, push.FcmServiceAccount)
		})
	}
}

func TestEncryptedPushStorageRead(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	enc := newTestEncryptor(t)
	sealed := encryptForTest(t, enc, testServiceAcct)
	otherKeyEnc, err := encryption.NewAEADLocal(ctx, "key-2", make([]byte, 32))
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		encryptor   encryption.Encryptor
		stored      string
		expected    string
		expectedErr error
	}{
		{
			desc:      "an encrypted value is decrypted",
			encryptor: enc,
			stored:    sealed,
			expected:  testServiceAcct,
		},
		{
			desc:      "a legacy plaintext value is returned unchanged",
			encryptor: enc,
			stored:    testServiceAcct,
			expected:  testServiceAcct,
		},
		{
			desc:      "an empty value stays empty",
			encryptor: enc,
			stored:    "",
			expected:  "",
		},
		{
			desc:      "none returns the plaintext unchanged",
			encryptor: encryption.NewNone(),
			stored:    testServiceAcct,
			expected:  testServiceAcct,
		},
		{
			desc:        "none refuses to return an encrypted value",
			encryptor:   encryption.NewNone(),
			stored:      sealed,
			expectedErr: encryption.ErrNotConfigured,
		},
		{
			desc:        "a value encrypted with another key is an error",
			encryptor:   otherKeyEnc,
			stored:      sealed,
			expectedErr: encryption.ErrUnknownKeyID,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			inner := mock.NewMockPushStorage(mockController)
			inner.EXPECT().GetPush(gomock.Any(), "id-0", testEnvironmentID).Return(
				&domain.Push{Push: &proto.Push{Id: "id-0", FcmServiceAccount: p.stored}}, nil,
			)
			inner.EXPECT().ListPushes(gomock.Any(), gomock.Any()).Return(
				[]*proto.Push{
					{Id: "id-0", FcmServiceAccount: p.stored},
					{Id: "id-1", FcmServiceAccount: p.stored},
				}, 2, int64(2), nil,
			)
			storage := v2ps.NewEncryptedPushStorage(inner, p.encryptor)

			push, err := storage.GetPush(ctx, "id-0", testEnvironmentID)
			pushes, nextOffset, totalCount, listErr := storage.ListPushes(ctx, v2ps.ListPushesParams{
				EnvironmentIDs: []string{testEnvironmentID},
			})
			if p.expectedErr != nil {
				assert.ErrorIs(t, err, p.expectedErr)
				assert.Nil(t, push)
				assert.ErrorIs(t, listErr, p.expectedErr)
				assert.Nil(t, pushes)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, p.expected, push.FcmServiceAccount)
			require.NoError(t, listErr)
			assert.Equal(t, 2, nextOffset)
			assert.Equal(t, int64(2), totalCount)
			require.Len(t, pushes, 2)
			for _, listed := range pushes {
				assert.Equal(t, p.expected, listed.FcmServiceAccount)
			}
		})
	}
}

func TestEncryptedPushStoragePassesThroughErrorsAndDelete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	inner := mock.NewMockPushStorage(mockController)
	storage := v2ps.NewEncryptedPushStorage(inner, newTestEncryptor(t))

	inner.EXPECT().GetPush(gomock.Any(), "id-0", testEnvironmentID).Return(nil, v2ps.ErrPushNotFound)
	push, err := storage.GetPush(ctx, "id-0", testEnvironmentID)
	assert.ErrorIs(t, err, v2ps.ErrPushNotFound)
	assert.Nil(t, push)

	listErr := errors.New("list error")
	inner.EXPECT().ListPushes(gomock.Any(), gomock.Any()).Return(nil, 0, int64(0), listErr)
	pushes, _, _, err := storage.ListPushes(ctx, v2ps.ListPushesParams{})
	assert.ErrorIs(t, err, listErr)
	assert.Nil(t, pushes)

	inner.EXPECT().CreatePush(gomock.Any(), gomock.Any(), testEnvironmentID).Return(v2ps.ErrPushAlreadyExists)
	err = storage.CreatePush(ctx, &domain.Push{Push: &proto.Push{FcmServiceAccount: testServiceAcct}}, testEnvironmentID)
	assert.ErrorIs(t, err, v2ps.ErrPushAlreadyExists)

	inner.EXPECT().DeletePush(gomock.Any(), "id-0", testEnvironmentID).Return(nil)
	assert.NoError(t, storage.DeletePush(ctx, "id-0", testEnvironmentID))
}
