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

package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tokentesting "github.com/bucketeer-io/bucketeer/v2/pkg/token/testing"
)

func TestSignedInAt(t *testing.T) {
	t.Parallel()
	issuedAt := time.Date(2026, 11, 24, 9, 0, 0, 0, time.UTC)
	authTime := issuedAt.Add(-48 * time.Hour)
	patterns := []struct {
		desc     string
		authTime time.Time
		expected time.Time
	}{
		{
			desc:     "auth_time present",
			authTime: authTime,
			expected: authTime,
		},
		{
			desc:     "auth_time absent: issued time",
			expected: issuedAt,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			access := &AccessToken{IssuedAt: issuedAt, AuthTime: p.authTime}
			refresh := &RefreshToken{IssuedAt: issuedAt, AuthTime: p.authTime}
			assert.True(t, p.expected.Equal(access.SignedInAt()))
			assert.True(t, p.expected.Equal(refresh.SignedInAt()))
		})
	}
}

func TestSignInClaimsRoundTrip(t *testing.T) {
	t.Parallel()
	issuer := "test_issuer"
	audience := "test_audience"
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := NewVerifier(publicKeyPath, issuer, audience)
	require.NoError(t, err)
	now := time.Now()
	authTime := now.Add(-time.Hour)

	accessRaw, err := signer.SignAccessToken(&AccessToken{
		Issuer:   issuer,
		Subject:  "user@example.com",
		Audience: audience,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
		Email:    "user@example.com",
		AuthTime: authTime,
	})
	require.NoError(t, err)
	access, err := verifier.VerifyAccessToken(accessRaw)
	require.NoError(t, err)
	assert.Equal(t, "user@example.com", access.Subject)
	assert.True(t, authTime.Equal(access.AuthTime))

	refreshRaw, err := signer.SignRefreshToken(&RefreshToken{
		Email:    "user@example.com",
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
		AuthTime: authTime,
	})
	require.NoError(t, err)
	refresh, err := verifier.VerifyRefreshToken(refreshRaw)
	require.NoError(t, err)
	assert.True(t, authTime.Equal(refresh.AuthTime))

	// A token issued before the claims existed verifies without them.
	legacyRaw := signRawClaims(t, privateKeyPath, map[string]any{
		"iss":   issuer,
		"aud":   audience,
		"exp":   now.Add(time.Hour),
		"iat":   now,
		"email": "user@example.com",
	})
	legacy, err := verifier.VerifyAccessToken(legacyRaw)
	require.NoError(t, err)
	assert.Empty(t, legacy.Subject)
	assert.True(t, legacy.AuthTime.IsZero())
	assert.True(t, now.Equal(legacy.SignedInAt()))
}
