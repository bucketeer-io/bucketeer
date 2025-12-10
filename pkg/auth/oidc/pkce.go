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
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const (
	// PKCE code verifier length (43-128 characters recommended by RFC 7636)
	codeVerifierLength = 64
)

// PKCEChallenge represents a PKCE code challenge and verifier pair
type PKCEChallenge struct {
	Verifier        string
	Challenge       string
	ChallengeMethod string
}

// GeneratePKCEChallenge generates a PKCE code verifier and challenge using S256 method
func GeneratePKCEChallenge() (*PKCEChallenge, error) {
	// Generate code verifier
	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	// Generate code challenge using S256 method
	challenge := generateS256Challenge(verifier)

	return &PKCEChallenge{
		Verifier:        verifier,
		Challenge:       challenge,
		ChallengeMethod: "S256",
	}, nil
}

// generateCodeVerifier generates a cryptographically random code verifier
func generateCodeVerifier() (string, error) {
	// Generate random bytes
	bytes := make([]byte, codeVerifierLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base64 URL-safe without padding
	verifier := base64.RawURLEncoding.EncodeToString(bytes)

	return verifier, nil
}

// generateS256Challenge generates a code challenge from a verifier using SHA256
func generateS256Challenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return challenge
}

// ValidateCodeVerifier validates that a code verifier is well-formed
func ValidateCodeVerifier(verifier string) error {
	if len(verifier) < 43 || len(verifier) > 128 {
		return fmt.Errorf("code verifier length must be between 43 and 128 characters")
	}
	return nil
}
