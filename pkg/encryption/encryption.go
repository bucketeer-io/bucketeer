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

// Package encryption protects secret values that are stored at rest and read
// back later, such as FCM service accounts and Slack webhook URLs.
//
// Encrypted values are self-describing text envelopes
// ("enc:v1:<key-id>:<base64url payload>") so that one column can hold values
// written with different keys and a reader can tell an encrypted value from a
// legacy plaintext one. The backend is selected by configuration; the "none"
// type keeps values as they are.
package encryption

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

const (
	// TypeNone stores values as they are.
	TypeNone = "none"
	// TypeAEADLocal encrypts values with AES-256-GCM using a key that is
	// provided locally (a mounted Secret or file), without any external service.
	TypeAEADLocal = "aead-local"
)

var (
	ErrUnsupportedType = errors.New("encryption: unsupported encryption type")
	ErrInvalidKeyID    = errors.New(
		"encryption: key id must be non-empty and contain only letters, digits, '.', '_' or '-'",
	)
	ErrInvalidKey      = errors.New("encryption: key must be 32 bytes encoded in base64")
	ErrMalformedValue  = errors.New("encryption: malformed encrypted value")
	ErrUnknownKeyID    = errors.New("encryption: value was encrypted with an unknown key id")
	ErrNotConfigured   = errors.New("encryption: value is encrypted but encryption type is none")
	errKeyPathRequired = errors.New("encryption: key path is required")
)

var keyIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Config selects the encryption backend.
type Config struct {
	// Type is one of TypeNone or TypeAEADLocal.
	Type string
	// KeyID identifies the key in every encrypted value so that keys can be
	// rotated without re-encrypting everything at once. Required unless Type
	// is TypeNone.
	KeyID string
	// KeyPath is the file holding the base64-encoded 32-byte key used by
	// TypeAEADLocal.
	KeyPath string
}

// Encryptor encrypts and decrypts secret values stored at rest.
type Encryptor interface {
	// Encrypt returns the envelope to store for plaintext. An empty plaintext
	// is returned unchanged.
	Encrypt(ctx context.Context, plaintext string) (string, error)
	// Decrypt returns the plaintext of a stored value. Values that are not
	// envelopes are legacy plaintext and are returned unchanged.
	Decrypt(ctx context.Context, value string) (string, error)
}

// New returns the Encryptor selected by cfg.
func New(ctx context.Context, cfg Config) (Encryptor, error) {
	switch cfg.Type {
	case TypeNone:
		return NewNone(), nil
	case TypeAEADLocal:
		if cfg.KeyPath == "" {
			return nil, errKeyPathRequired
		}
		key, err := loadKey(cfg.KeyPath)
		if err != nil {
			return nil, err
		}
		return NewAEADLocal(ctx, cfg.KeyID, key)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedType, cfg.Type)
	}
}

func validateKeyID(keyID string) error {
	if !keyIDPattern.MatchString(keyID) {
		return ErrInvalidKeyID
	}
	return nil
}
