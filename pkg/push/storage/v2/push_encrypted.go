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
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

// encryptedPushStorage protects the FCM service account at rest. It encrypts
// the value before it is written and decrypts it after it is read, so callers
// keep working with the plaintext credential whatever the configured backend.
type encryptedPushStorage struct {
	inner     PushStorage
	encryptor encryption.Encryptor
}

// NewEncryptedPushStorage wraps inner so that the FCM service account is
// stored through encryptor.
func NewEncryptedPushStorage(inner PushStorage, encryptor encryption.Encryptor) PushStorage {
	return &encryptedPushStorage{inner: inner, encryptor: encryptor}
}

func (s *encryptedPushStorage) CreatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	sealed, err := s.seal(ctx, e)
	if err != nil {
		return err
	}
	return s.inner.CreatePush(ctx, sealed, environmentId)
}

func (s *encryptedPushStorage) UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	sealed, err := s.seal(ctx, e)
	if err != nil {
		return err
	}
	return s.inner.UpdatePush(ctx, sealed, environmentId)
}

func (s *encryptedPushStorage) GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error) {
	push, err := s.inner.GetPush(ctx, id, environmentId)
	if err != nil {
		return nil, err
	}
	if err := s.open(ctx, push.Push); err != nil {
		return nil, err
	}
	return push, nil
}

func (s *encryptedPushStorage) ListPushes(
	ctx context.Context,
	p ListPushesParams,
) ([]*proto.Push, int, int64, error) {
	pushes, nextOffset, totalCount, err := s.inner.ListPushes(ctx, p)
	if err != nil {
		return nil, 0, 0, err
	}
	for _, push := range pushes {
		if err := s.open(ctx, push); err != nil {
			return nil, 0, 0, err
		}
	}
	return pushes, nextOffset, totalCount, nil
}

func (s *encryptedPushStorage) DeletePush(ctx context.Context, id, environmentId string) error {
	return s.inner.DeletePush(ctx, id, environmentId)
}

// seal returns a copy of e whose FCM service account is encrypted. The
// caller's entity is left untouched. A value that is already encrypted is
// stored as it is so that a re-saved entity is never encrypted twice.
func (s *encryptedPushStorage) seal(ctx context.Context, e *domain.Push) (*domain.Push, error) {
	if e.FcmServiceAccount == "" || encryption.IsEncrypted(e.FcmServiceAccount) {
		return e, nil
	}
	sealed, err := s.encryptor.Encrypt(ctx, e.FcmServiceAccount)
	if err != nil {
		return nil, err
	}
	copied := pb.Clone(e.Push).(*proto.Push)
	copied.FcmServiceAccount = sealed
	return &domain.Push{Push: copied}, nil
}

// open replaces the stored FCM service account of push with its plaintext.
// Legacy plaintext values are returned unchanged by the encryptor.
func (s *encryptedPushStorage) open(ctx context.Context, push *proto.Push) error {
	if push.FcmServiceAccount == "" {
		return nil
	}
	plaintext, err := s.encryptor.Decrypt(ctx, push.FcmServiceAccount)
	if err != nil {
		return err
	}
	push.FcmServiceAccount = plaintext
	return nil
}
