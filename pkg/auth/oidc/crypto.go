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

package oidc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	encryptedPrefix = "encrypted:"
)

var (
	ErrInvalidEncryptionKey = errors.New("encryption key must be 32 bytes for AES-256")
	ErrInvalidCiphertext    = errors.New("invalid ciphertext")
)

// EncryptClientSecret encrypts a client secret using AES-256-GCM
func EncryptClientSecret(plaintext string, encryptionKey []byte) (string, error) {
	if len(encryptionKey) != 32 {
		return "", ErrInvalidEncryptionKey
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encryptedPrefix + encoded, nil
}

// DecryptClientSecret decrypts a client secret using AES-256-GCM
func DecryptClientSecret(encrypted string, encryptionKey []byte) (string, error) {
	// If not encrypted, return as-is (for backward compatibility during migration)
	if !strings.HasPrefix(encrypted, encryptedPrefix) {
		return encrypted, nil
	}

	if len(encryptionKey) != 32 {
		return "", ErrInvalidEncryptionKey
	}

	// Remove prefix and decode
	encoded := strings.TrimPrefix(encrypted, encryptedPrefix)
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a value is encrypted
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}
