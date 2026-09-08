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
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	// Nonce length in bytes (32 bytes = 256 bits)
	nonceLength = 32
)

// GenerateNonce generates a cryptographically random nonce for OIDC requests
func GenerateNonce() (string, error) {
	bytes := make([]byte, nonceLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encode to base64 URL-safe without padding
	nonce := base64.RawURLEncoding.EncodeToString(bytes)
	return nonce, nil
}

// ValidateNonce validates that a nonce is well-formed
func ValidateNonce(nonce string) error {
	if len(nonce) == 0 {
		return fmt.Errorf("nonce cannot be empty")
	}

	// Decode to ensure it's valid base64
	_, err := base64.RawURLEncoding.DecodeString(nonce)
	if err != nil {
		return fmt.Errorf("invalid nonce encoding: %w", err)
	}

	return nil
}
