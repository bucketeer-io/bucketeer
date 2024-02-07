// Copyright 2024 The Bucketeer Authors.
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
)

func TestNewSigner(t *testing.T) {
	testcases := []struct {
		path string
		ok   bool
	}{
		{"testdata/valid-private.pem", true},
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
	clientID := "test_client_id"
	signer, err := NewSigner("testdata/valid-private.pem")
	require.NoError(t, err)
	verifier, err := NewVerifier("testdata/valid-public.pem", issuer, clientID)
	require.NoError(t, err)
	testcases := []struct {
		token *IDToken
		ok    bool
	}{
		{
			&IDToken{
				Issuer:   issuer,
				Subject:  "subject",
				Audience: clientID,
				Email:    "test@email",
				Expiry:   time.Now().Add(time.Hour),
			},
			true,
		},
		{
			&IDToken{
				Issuer:   issuer,
				Subject:  "subject",
				Audience: clientID,
				Expiry:   time.Now().Add(time.Hour),
			},
			false,
		},
		{
			&IDToken{
				Issuer:   issuer,
				Subject:  "subject",
				Audience: clientID,
				Email:    "test@email",
				Expiry:   time.Now().Add(-time.Hour),
			},
			false,
		},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		signedToken, err := signer.Sign(tc.token)
		require.NoError(t, err, des)
		require.True(t, len(signedToken) > 0, des)
		parsedToken, err := verifier.Verify(signedToken)
		if tc.ok {
			require.NoError(t, err, fmt.Sprintf("index: %d, error: %v", i, err))
			require.Equal(t, tc.token.Issuer, parsedToken.Issuer, des)
			require.Equal(t, tc.token.Subject, parsedToken.Subject, des)
			require.Equal(t, tc.token.Audience, parsedToken.Audience, des)
			require.True(t, tc.token.Expiry.Equal(parsedToken.Expiry), des)
			require.True(t, tc.token.IssuedAt.Equal(parsedToken.IssuedAt), des)
			require.Equal(t, tc.token.Email, parsedToken.Email, des)
			require.Equal(t, tc.token.IsSystemAdmin, parsedToken.IsSystemAdmin, des)
		} else {
			require.Error(t, err, des)
		}
	}
}
