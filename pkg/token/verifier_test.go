// Copyright 2025 The Bucketeer Authors.
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
	"strings"
	"testing"
	"time"

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
