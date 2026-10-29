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
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVerifier(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		path string
		ok   bool
	}{
		{"testdata/valid-public.pem", true},
		{"testdata/invalid-public.pem", false},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		s, err := NewVerifier(tc.path, "issuer", "client_id")
		assert.Equal(t, err == nil, tc.ok, fmt.Sprintf("index: %d, err: %v", i, err))
		assert.Equal(t, s != nil, tc.ok, des)
	}
}

func TestVerify(t *testing.T) {
	t.Parallel()
	issuer := "test_issuer"
	audience := "test_audience"
	signer, err := NewSigner("testdata/valid-private.pem")
	require.NoError(t, err)
	accessToken := &AccessToken{
		Issuer:   issuer,
		Audience: audience,
		Email:    "test@email",
		Expiry:   time.Now().Add(time.Hour),
	}
	testcases := []struct {
		desc           string
		rawAccessToken string
		valid          bool
	}{
		{
			desc:           "err: malformed jwt",
			rawAccessToken: "",
			valid:          false,
		},
		{
			desc:           "err: invalid jwt",
			rawAccessToken: createInvalidRawIDToken(t, signer, accessToken),
			valid:          false,
		},
		{
			desc:           "success",
			rawAccessToken: createValidRawIDToken(t, signer, accessToken),
			valid:          true,
		},
	}
	verifier, err := NewVerifier("testdata/valid-public.pem", issuer, audience)
	require.NoError(t, err)
	for _, p := range testcases {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			actualToken, err := verifier.VerifyAccessToken(p.rawAccessToken)
			if p.valid {
				assert.NotNil(t, actualToken)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, actualToken)
				assert.Error(t, err)
			}
		})
	}
}

func createValidRawIDToken(t *testing.T, signer Signer, accessToken *AccessToken) string {
	t.Helper()
	rawIDToken, err := signer.SignAccessToken(accessToken)
	require.NoError(t, err)
	return rawIDToken
}

func createInvalidRawIDToken(t *testing.T, signer Signer, accessToken *AccessToken) string {
	t.Helper()
	rawIDToken, err := signer.SignAccessToken(accessToken)
	require.NoError(t, err)
	parts := strings.Split(rawIDToken, ".")
	invalidSignature := base64.RawURLEncoding.EncodeToString([]byte("invalid-signature"))
	return fmt.Sprintf("%s.%s.%s", parts[0], parts[1], invalidSignature)
}

func TestVerifyDemoCreationToken(t *testing.T) {
	t.Parallel()
	issuer := "test_issuer"
	audience := "test_audience"
	signer, err := NewSigner("testdata/valid-private.pem")
	require.NoError(t, err)
	demoToken := &DemoCreationToken{
		Issuer:   issuer,
		Audience: audience,
		Email:    "test@email.com",
		Expiry:   time.Now().Add(time.Hour),
		IssuedAt: time.Now(),
	}
	testcases := []struct {
		desc         string
		rawDemoToken string
		valid        bool
	}{
		{
			desc:         "err: malformed jwt",
			rawDemoToken: "",
			valid:        false,
		},
		{
			desc:         "err: invalid jwt",
			rawDemoToken: createInvalidRawDemoToken(t, signer, demoToken),
			valid:        false,
		},
		{
			desc:         "success",
			rawDemoToken: createValidRawDemoToken(t, signer, demoToken),
			valid:        true,
		},
	}
	verifier, err := NewVerifier("testdata/valid-public.pem", issuer, audience)
	require.NoError(t, err)
	for _, p := range testcases {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			actualToken, err := verifier.VerifyDemoCreationToken(p.rawDemoToken)
			if p.valid {
				assert.NotNil(t, actualToken)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, actualToken)
				assert.Error(t, err)
			}
		})
	}
}

func TestVerifyDemoCreationTokenValidation(t *testing.T) {
	t.Parallel()
	issuer := "test_issuer"
	audience := "test_audience"
	signer, err := NewSigner("testdata/valid-private.pem")
	require.NoError(t, err)
	verifier, err := NewVerifier("testdata/valid-public.pem", issuer, audience)
	require.NoError(t, err)

	testcases := []struct {
		desc        string
		demoToken   *DemoCreationToken
		expectError bool
	}{
		{
			desc: "success: valid token",
			demoToken: &DemoCreationToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    "test@email.com",
				Expiry:   time.Now().Add(time.Hour),
				IssuedAt: time.Now(),
			},
			expectError: false,
		},
		{
			desc: "error: wrong issuer",
			demoToken: &DemoCreationToken{
				Issuer:   "wrong_issuer",
				Audience: audience,
				Email:    "test@email.com",
				Expiry:   time.Now().Add(time.Hour),
				IssuedAt: time.Now(),
			},
			expectError: true,
		},
		{
			desc: "error: wrong audience",
			demoToken: &DemoCreationToken{
				Issuer:   issuer,
				Audience: "wrong_audience",
				Email:    "test@email.com",
				Expiry:   time.Now().Add(time.Hour),
				IssuedAt: time.Now(),
			},
			expectError: true,
		},
		{
			desc: "error: expired token",
			demoToken: &DemoCreationToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    "test@email.com",
				Expiry:   time.Now().Add(-time.Hour),
				IssuedAt: time.Now().Add(-2 * time.Hour),
			},
			expectError: true,
		},
		{
			desc: "error: empty email",
			demoToken: &DemoCreationToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    "",
				Expiry:   time.Now().Add(time.Hour),
				IssuedAt: time.Now(),
			},
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			rawToken := createValidRawDemoToken(t, signer, tc.demoToken)
			actualToken, err := verifier.VerifyDemoCreationToken(rawToken)

			if tc.expectError {
				assert.Nil(t, actualToken)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, actualToken)
				assert.NoError(t, err)
				assert.Equal(t, tc.demoToken.Issuer, actualToken.Issuer)
				assert.Equal(t, tc.demoToken.Audience, actualToken.Audience)
				assert.Equal(t, tc.demoToken.Email, actualToken.Email)
				assert.True(t, tc.demoToken.Expiry.Equal(actualToken.Expiry))
				assert.True(t, tc.demoToken.IssuedAt.Equal(actualToken.IssuedAt))
			}
		})
	}
}

func createValidRawDemoToken(t *testing.T, signer Signer, demoToken *DemoCreationToken) string {
	t.Helper()
	rawDemoToken, err := signer.SignDemoCreationToken(demoToken)
	require.NoError(t, err)
	return rawDemoToken
}

func createInvalidRawDemoToken(t *testing.T, signer Signer, demoToken *DemoCreationToken) string {
	t.Helper()
	rawDemoToken, err := signer.SignDemoCreationToken(demoToken)
	require.NoError(t, err)
	parts := strings.Split(rawDemoToken, ".")
	invalidSignature := base64.RawURLEncoding.EncodeToString([]byte("invalid-signature"))
	return fmt.Sprintf("%s.%s.%s", parts[0], parts[1], invalidSignature)
}

// legacyAccessToken is the access-token payload from before token_use existed:
// iss, aud, exp, and email, with no token_use claim.
type legacyAccessToken struct {
	Issuer        string    `json:"iss"`
	Audience      string    `json:"aud"`
	Expiry        time.Time `json:"exp"`
	IssuedAt      time.Time `json:"iat"`
	Email         string    `json:"email"`
	IsSystemAdmin bool      `json:"is_system_admin"`
}

// legacyRefreshToken is the refresh-token payload from before token_use, iss, and aud existed.
type legacyRefreshToken struct {
	Email          string    `json:"email"`
	Expiry         time.Time `json:"exp"`
	IssuedAt       time.Time `json:"iat"`
	OrganizationID string    `json:"organization_id"`
}

func TestVerifierAcceptsOnlyMatchingTokenUse(t *testing.T) {
	t.Parallel()
	issuer := "test_issuer"
	audience := "test_audience"
	email := "user@example.com"
	signer, err := NewSigner("testdata/valid-private.pem")
	require.NoError(t, err)
	verifier, err := NewVerifier("testdata/valid-public.pem", issuer, audience)
	require.NoError(t, err)

	now := time.Now()
	// A caller-supplied token_use must not survive signing.
	accessRaw, err := signer.SignAccessToken(&AccessToken{
		Issuer:   issuer,
		Audience: audience,
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
		TokenUse: TokenUseDemo,
	})
	require.NoError(t, err)
	refreshRaw, err := signer.SignRefreshToken(&RefreshToken{
		Issuer:         issuer,
		Audience:       audience,
		Email:          email,
		Expiry:         now.Add(time.Hour),
		IssuedAt:       now,
		OrganizationID: "org-1",
	})
	require.NoError(t, err)
	demoRaw, err := signer.SignDemoCreationToken(&DemoCreationToken{
		Issuer:   issuer,
		Audience: audience,
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
	})
	require.NoError(t, err)
	wrongIssuerRefresh, err := signer.SignRefreshToken(&RefreshToken{
		Issuer:   "other-issuer",
		Audience: audience,
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
	})
	require.NoError(t, err)
	wrongAudienceRefresh, err := signer.SignRefreshToken(&RefreshToken{
		Issuer:   issuer,
		Audience: "other-audience",
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
	})
	require.NoError(t, err)
	missingIssuerRefresh, err := signer.SignRefreshToken(&RefreshToken{
		Audience: audience,
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
	})
	require.NoError(t, err)
	expiredRefresh, err := signer.SignRefreshToken(&RefreshToken{
		Issuer:   issuer,
		Audience: audience,
		Email:    email,
		Expiry:   now.Add(-time.Hour),
		IssuedAt: now.Add(-2 * time.Hour),
	})
	require.NoError(t, err)
	emptyEmailRefresh, err := signer.SignRefreshToken(&RefreshToken{
		Issuer:   issuer,
		Audience: audience,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
	})
	require.NoError(t, err)

	legacyAccess := signRawClaims(t, &legacyAccessToken{
		Issuer:        issuer,
		Audience:      audience,
		Email:         email,
		Expiry:        now.Add(time.Hour),
		IssuedAt:      now,
		IsSystemAdmin: true,
	})
	legacyRefresh := signRawClaims(t, &legacyRefreshToken{
		Email:          email,
		Expiry:         now.Add(time.Hour),
		IssuedAt:       now,
		OrganizationID: "org-1",
	})

	access, err := verifier.VerifyAccessToken(accessRaw)
	require.NoError(t, err)
	assert.Equal(t, TokenUseAccess, access.TokenUse)
	assert.Equal(t, email, access.Email)
	assert.Equal(t, issuer, access.Issuer)
	assert.Equal(t, audience, access.Audience)

	refresh, err := verifier.VerifyRefreshToken(refreshRaw)
	require.NoError(t, err)
	assert.Equal(t, TokenUseRefresh, refresh.TokenUse)
	assert.Equal(t, email, refresh.Email)
	assert.Equal(t, issuer, refresh.Issuer)
	assert.Equal(t, audience, refresh.Audience)
	assert.Equal(t, "org-1", refresh.OrganizationID)

	demo, err := verifier.VerifyDemoCreationToken(demoRaw)
	require.NoError(t, err)
	assert.Equal(t, TokenUseDemo, demo.TokenUse)
	assert.Equal(t, email, demo.Email)
	assert.Equal(t, issuer, demo.Issuer)
	assert.Equal(t, audience, demo.Audience)

	legacyAccessParsed, err := verifier.VerifyAccessToken(legacyAccess)
	require.NoError(t, err)
	assert.Empty(t, legacyAccessParsed.TokenUse)
	assert.Equal(t, email, legacyAccessParsed.Email)
	assert.Equal(t, issuer, legacyAccessParsed.Issuer)
	assert.True(t, legacyAccessParsed.IsSystemAdmin)

	legacyRefreshParsed, err := verifier.VerifyRefreshToken(legacyRefresh)
	require.NoError(t, err)
	assert.Empty(t, legacyRefreshParsed.TokenUse)
	assert.Equal(t, email, legacyRefreshParsed.Email)
	assert.Empty(t, legacyRefreshParsed.Issuer)
	assert.Empty(t, legacyRefreshParsed.Audience)
	assert.Equal(t, "org-1", legacyRefreshParsed.OrganizationID)

	verifyAccess := func(raw string) error {
		_, err := verifier.VerifyAccessToken(raw)
		return err
	}
	verifyRefresh := func(raw string) error {
		_, err := verifier.VerifyRefreshToken(raw)
		return err
	}
	verifyDemo := func(raw string) error {
		_, err := verifier.VerifyDemoCreationToken(raw)
		return err
	}

	cases := []struct {
		name          string
		raw           string
		verify        func(string) error
		wantSubstring string
	}{
		{
			name:          "access token is rejected as a refresh token",
			raw:           accessRaw,
			verify:        verifyRefresh,
			wantSubstring: "token_use",
		},
		{
			name:          "access token is rejected as a demo creation token",
			raw:           accessRaw,
			verify:        verifyDemo,
			wantSubstring: "token_use",
		},
		{
			name:          "refresh token is rejected as an access token",
			raw:           refreshRaw,
			verify:        verifyAccess,
			wantSubstring: "token_use",
		},
		{
			name:          "refresh token is rejected as a demo creation token",
			raw:           refreshRaw,
			verify:        verifyDemo,
			wantSubstring: "token_use",
		},
		{
			name:          "demo creation token is rejected as an access token",
			raw:           demoRaw,
			verify:        verifyAccess,
			wantSubstring: "token_use",
		},
		{
			name:          "demo creation token is rejected as a refresh token",
			raw:           demoRaw,
			verify:        verifyRefresh,
			wantSubstring: "token_use",
		},
		{
			name:          "refresh token with the wrong issuer is rejected",
			raw:           wrongIssuerRefresh,
			verify:        verifyRefresh,
			wantSubstring: "different provider",
		},
		{
			name:          "refresh token with the wrong audience is rejected",
			raw:           wrongAudienceRefresh,
			verify:        verifyRefresh,
			wantSubstring: "expected audience",
		},
		{
			name:          "refresh token without an issuer is rejected",
			raw:           missingIssuerRefresh,
			verify:        verifyRefresh,
			wantSubstring: "different provider",
		},
		{
			name:          "expired refresh token is rejected",
			raw:           expiredRefresh,
			verify:        verifyRefresh,
			wantSubstring: "expired",
		},
		{
			name:          "refresh token without an email is rejected",
			raw:           emptyEmailRefresh,
			verify:        verifyRefresh,
			wantSubstring: "email must be not empty",
		},
		{
			name:          "access token without token_use is rejected as a refresh token",
			raw:           legacyAccess,
			verify:        verifyRefresh,
			wantSubstring: "token_use",
		},
		{
			name:          "access token without token_use is rejected as a demo creation token",
			raw:           legacyAccess,
			verify:        verifyDemo,
			wantSubstring: "token_use",
		},
		{
			name:          "refresh token without token_use is rejected as an access token",
			raw:           legacyRefresh,
			verify:        verifyAccess,
			wantSubstring: "different provider",
		},
		{
			name:          "refresh token without token_use is rejected as a demo creation token",
			raw:           legacyRefresh,
			verify:        verifyDemo,
			wantSubstring: "token_use",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.verify(tc.raw)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantSubstring)
		})
	}
}

func signRawClaims(t *testing.T, claims any) string {
	t.Helper()
	data, err := os.ReadFile("testdata/valid-private.pem")
	require.NoError(t, err)
	key, err := parseRSAPrivateKey(data)
	require.NoError(t, err)
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, nil)
	require.NoError(t, err)
	raw, err := jwt.Signed(sig).Claims(claims).Serialize()
	require.NoError(t, err)
	return raw
}
