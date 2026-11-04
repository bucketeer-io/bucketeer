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

package encryption

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-kms-wrapping/v2/aead"
	"google.golang.org/protobuf/proto"
)

// aeadKeySize is the AES-256 key length.
const aeadKeySize = 32

// wrapperEncryptor stores values encrypted by a go-kms-wrapping Wrapper. The
// wrapper output (ciphertext plus key metadata) is serialized into the
// envelope payload, so the same code path can later serve KMS-backed wrappers.
type wrapperEncryptor struct {
	keyID   string
	wrapper wrapping.Wrapper
}

// NewAEADLocal returns an Encryptor that uses AES-256-GCM with the given
// 32-byte key. keyID is recorded in every encrypted value.
func NewAEADLocal(ctx context.Context, keyID string, key []byte) (Encryptor, error) {
	if err := validateKeyID(keyID); err != nil {
		return nil, err
	}
	if len(key) != aeadKeySize {
		return nil, ErrInvalidKey
	}
	wrapper := aead.NewWrapper()
	if _, err := wrapper.SetConfig(ctx, wrapping.WithKeyId(keyID), aead.WithKey(key)); err != nil {
		return nil, fmt.Errorf("encryption: configure aead wrapper: %w", err)
	}
	return &wrapperEncryptor{keyID: keyID, wrapper: wrapper}, nil
}

// loadKey reads a base64-encoded 32-byte key from path (the format produced
// by `openssl rand -base64 32`).
func loadKey(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("encryption: read key file: %w", err)
	}
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(content)))
	if err != nil {
		return nil, ErrInvalidKey
	}
	if len(key) != aeadKeySize {
		return nil, ErrInvalidKey
	}
	return key, nil
}

func (e *wrapperEncryptor) Encrypt(ctx context.Context, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	blob, err := e.wrapper.Encrypt(ctx, []byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("encryption: encrypt: %w", err)
	}
	payload, err := proto.Marshal(blob)
	if err != nil {
		return "", fmt.Errorf("encryption: encode ciphertext: %w", err)
	}
	return encodeEnvelope(e.keyID, payload), nil
}

func (e *wrapperEncryptor) Decrypt(ctx context.Context, value string) (string, error) {
	if !IsEncrypted(value) {
		return value, nil
	}
	keyID, payload, err := decodeEnvelope(value)
	if err != nil {
		return "", err
	}
	if keyID != e.keyID {
		return "", fmt.Errorf("%w: %q", ErrUnknownKeyID, keyID)
	}
	blob := &wrapping.BlobInfo{}
	if err := proto.Unmarshal(payload, blob); err != nil {
		return "", fmt.Errorf("%w: %v", ErrMalformedValue, err)
	}
	plaintext, err := e.wrapper.Decrypt(ctx, blob)
	if err != nil {
		return "", fmt.Errorf("encryption: decrypt: %w", err)
	}
	return string(plaintext), nil
}
