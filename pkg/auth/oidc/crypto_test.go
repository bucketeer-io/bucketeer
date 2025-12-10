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
	"strings"
	"testing"
)

func TestEncryptDecryptClientSecret(t *testing.T) {
	key := make([]byte, 32) // 32 bytes for AES-256
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := "my-client-secret-123"

	encrypted, err := EncryptClientSecret(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptClientSecret() error = %v", err)
	}

	if !strings.HasPrefix(encrypted, encryptedPrefix) {
		t.Errorf("EncryptClientSecret() result doesn't have encrypted prefix")
	}

	if encrypted == plaintext {
		t.Error("EncryptClientSecret() returned same value as input")
	}

	decrypted, err := DecryptClientSecret(encrypted, key)
	if err != nil {
		t.Fatalf("DecryptClientSecret() error = %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("DecryptClientSecret() = %v, want %v", decrypted, plaintext)
	}
}

func TestEncryptClientSecretInvalidKey(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
	}{
		{"key too short", 16},
		{"key too long", 48},
		{"empty key", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)
			_, err := EncryptClientSecret("test", key)
			if err == nil {
				t.Error("EncryptClientSecret() expected error for invalid key size")
			}
		})
	}
}

func TestDecryptClientSecretInvalidKey(t *testing.T) {
	key := make([]byte, 32)
	encrypted, _ := EncryptClientSecret("test", key)

	wrongKey := make([]byte, 32)
	for i := range wrongKey {
		wrongKey[i] = byte(255 - i)
	}

	_, err := DecryptClientSecret(encrypted, wrongKey)
	if err == nil {
		t.Error("DecryptClientSecret() expected error when using wrong key")
	}
}

func TestDecryptClientSecretBackwardCompatibility(t *testing.T) {
	key := make([]byte, 32)
	plainSecret := "unencrypted-secret"

	// Should return plaintext if not encrypted (backward compatibility)
	decrypted, err := DecryptClientSecret(plainSecret, key)
	if err != nil {
		t.Fatalf("DecryptClientSecret() error = %v", err)
	}

	if decrypted != plainSecret {
		t.Errorf("DecryptClientSecret() = %v, want %v", decrypted, plainSecret)
	}
}

func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "encrypted value",
			value: "encrypted:abc123",
			want:  true,
		},
		{
			name:  "plain value",
			value: "plain-secret",
			want:  false,
		},
		{
			name:  "empty value",
			value: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEncrypted(tt.value); got != tt.want {
				t.Errorf("IsEncrypted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryptionUniqueness(t *testing.T) {
	key := make([]byte, 32)
	plaintext := "test-secret"

	encrypted1, _ := EncryptClientSecret(plaintext, key)
	encrypted2, _ := EncryptClientSecret(plaintext, key)

	// Should produce different ciphertexts due to random nonce
	if encrypted1 == encrypted2 {
		t.Error("EncryptClientSecret() produced same ciphertext for same plaintext")
	}

	// But both should decrypt to same value
	decrypted1, _ := DecryptClientSecret(encrypted1, key)
	decrypted2, _ := DecryptClientSecret(encrypted2, key)

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Error("Decryption failed for unique ciphertexts")
	}
}
