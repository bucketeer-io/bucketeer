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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tokentesting "github.com/bucketeer-io/bucketeer/v2/pkg/token/testing"
)

func TestNewSigner(t *testing.T) {
	privateKeyPath, _ := tokentesting.WriteKeyPair(t)
	testcases := []struct {
		path string
		ok   bool
	}{
		{privateKeyPath, true},
		{"testdata/invalid-private.pem", false},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		s, err := NewSigner(tc.path)
		assert.Equal(t, err == nil, tc.ok, des)
		assert.Equal(t, s != nil, tc.ok, des)
	}
}

func TestSign(t *testing.T) {
	issuer := "test_issuer"
	audience := "test_client_id"
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := NewVerifier(publicKeyPath, issuer, audience)
	require.NoError(t, err)
	testcases := []struct {
		token *AccessToken
		ok    bool
	}{
		{
			&AccessToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    "test@email",
				Expiry:   time.Now().Add(time.Hour),
			},
			true,
		},
		{
			&AccessToken{
				Issuer:   issuer,
				Audience: audience,
				Expiry:   time.Now().Add(time.Hour),
			},
			false,
		},
		{
			&AccessToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    "test@email",
				Expiry:   time.Now().Add(-time.Hour),
			},
			false,
		},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		signedToken, err := signer.SignAccessToken(tc.token)
		require.NoError(t, err, des)
		require.True(t, len(signedToken) > 0, des)
		parsedToken, err := verifier.VerifyAccessToken(signedToken)
		if tc.ok {
			require.NoError(t, err, fmt.Sprintf("index: %d, error: %v", i, err))
			require.Equal(t, tc.token.Issuer, parsedToken.Issuer, des)
			require.Equal(t, tc.token.Audience, parsedToken.Audience, des)
			require.True(t, tc.token.Expiry.Equal(parsedToken.Expiry), des)
			require.True(t, tc.token.IssuedAt.Equal(parsedToken.IssuedAt), des)
			require.Equal(t, tc.token.Email, parsedToken.Email, des)
			require.Equal(t, tc.token.IsSystemAdmin, parsedToken.IsSystemAdmin, des)
			require.Equal(t, TokenUseAccess, parsedToken.TokenUse, des)
		} else {
			require.Error(t, err, des)
		}
	}
}

func TestSignDemoCreationToken(t *testing.T) {
	t.Parallel()
	issuer := "test_issuer"
	audience := "test_client_id"
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := NewVerifier(publicKeyPath, issuer, audience)
	require.NoError(t, err)
	testcases := []struct {
		desc  string
		token *DemoCreationToken
		ok    bool
	}{
		{
			desc: "success: valid demo token",
			token: &DemoCreationToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    "test@email.com",
				Expiry:   time.Now().Add(time.Hour),
				IssuedAt: time.Now(),
			},
			ok: true,
		},
		{
			desc: "error: missing email",
			token: &DemoCreationToken{
				Issuer:   issuer,
				Audience: audience,
				Expiry:   time.Now().Add(time.Hour),
				IssuedAt: time.Now(),
			},
			ok: false,
		},
		{
			desc: "error: expired token",
			token: &DemoCreationToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    "test@email.com",
				Expiry:   time.Now().Add(-time.Hour), // expired
				IssuedAt: time.Now().Add(-2 * time.Hour),
			},
			ok: false,
		},
	}
	for i, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			des := fmt.Sprintf("index: %d, desc: %s", i, tc.desc)
			signedToken, err := signer.SignDemoCreationToken(tc.token)
			require.NoError(t, err, des)
			require.True(t, len(signedToken) > 0, des)
			parsedToken, err := verifier.VerifyDemoCreationToken(signedToken)
			if tc.ok {
				require.NoError(t, err, fmt.Sprintf("index: %d, error: %v", i, err))
				require.Equal(t, tc.token.Issuer, parsedToken.Issuer, des)
				require.Equal(t, tc.token.Audience, parsedToken.Audience, des)
				require.True(t, tc.token.Expiry.Equal(parsedToken.Expiry), des)
				require.True(t, tc.token.IssuedAt.Equal(parsedToken.IssuedAt), des)
				require.Equal(t, tc.token.Email, parsedToken.Email, des)
				require.Equal(t, TokenUseDemo, parsedToken.TokenUse, des)
			} else {
				require.Error(t, err, des)
			}
		})
	}
}

func TestSignServiceToken(t *testing.T) {
	t.Parallel()
	issuer := "test_issuer"
	audience := "test_client_id"
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := NewVerifier(publicKeyPath, issuer, audience)
	require.NoError(t, err)
	now := time.Now()
	testcases := []struct {
		desc        string
		token       *ServiceToken
		expectedErr string
	}{
		{
			desc: "success: valid service token",
			token: &ServiceToken{
				Issuer:   issuer,
				Subject:  "batch",
				Audience: audience,
				Expiry:   now.Add(time.Hour),
				IssuedAt: now,
				Email:    "service@example.com",
				ID:       "token-1",
				// A caller-supplied token_use must not survive signing.
				TokenUse: TokenUseAccess,
			},
		},
		{
			desc: "success: without an email",
			token: &ServiceToken{
				Issuer:   issuer,
				Subject:  "batch",
				Audience: audience,
				Expiry:   now.Add(time.Hour),
				IssuedAt: now,
				ID:       "token-1",
			},
		},
		{
			desc: "error: missing sub",
			token: &ServiceToken{
				Issuer:   issuer,
				Audience: audience,
				Expiry:   now.Add(time.Hour),
				IssuedAt: now,
				ID:       "token-1",
			},
			expectedErr: "sub must be not empty",
		},
		{
			desc: "error: missing jti",
			token: &ServiceToken{
				Issuer:   issuer,
				Subject:  "batch",
				Audience: audience,
				Expiry:   now.Add(time.Hour),
				IssuedAt: now,
			},
			expectedErr: "jti must be not empty",
		},
		{
			desc: "error: wrong issuer",
			token: &ServiceToken{
				Issuer:   "other",
				Subject:  "batch",
				Audience: audience,
				Expiry:   now.Add(time.Hour),
				IssuedAt: now,
				ID:       "token-1",
			},
			expectedErr: "different provider",
		},
		{
			desc: "error: wrong audience",
			token: &ServiceToken{
				Issuer:   issuer,
				Subject:  "batch",
				Audience: "other",
				Expiry:   now.Add(time.Hour),
				IssuedAt: now,
				ID:       "token-1",
			},
			expectedErr: "audience",
		},
		{
			desc: "error: expired token",
			token: &ServiceToken{
				Issuer:   issuer,
				Subject:  "batch",
				Audience: audience,
				Expiry:   now.Add(-time.Hour),
				IssuedAt: now.Add(-2 * time.Hour),
				ID:       "token-1",
			},
			expectedErr: "expired",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			signedToken, err := signer.SignServiceToken(tc.token)
			require.NoError(t, err)
			parsedToken, err := verifier.VerifyServiceToken(signedToken)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.token.Issuer, parsedToken.Issuer)
			assert.Equal(t, tc.token.Subject, parsedToken.Subject)
			assert.Equal(t, tc.token.Audience, parsedToken.Audience)
			assert.True(t, tc.token.Expiry.Equal(parsedToken.Expiry))
			assert.True(t, tc.token.IssuedAt.Equal(parsedToken.IssuedAt))
			assert.Equal(t, tc.token.Email, parsedToken.Email)
			assert.Equal(t, tc.token.ID, parsedToken.ID)
			assert.Equal(t, TokenUseService, parsedToken.TokenUse)
		})
	}
}
