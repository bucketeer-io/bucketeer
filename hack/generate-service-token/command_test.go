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

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	tokentesting "github.com/bucketeer-io/bucketeer/v2/pkg/token/testing"
)

func TestCommandSign(t *testing.T) {
	t.Parallel()
	const (
		issuer   = "test-issuer"
		audience = "test-audience"
		email    = "service@example.com"
	)
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := token.NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := token.NewVerifier(publicKeyPath, issuer, audience)
	require.NoError(t, err)
	now := time.Now()

	patterns := []struct {
		desc           string
		service        string
		ttl            time.Duration
		expectedExpiry time.Time
	}{
		{
			desc:           "legacy form lasts 100 years",
			expectedExpiry: now.AddDate(100, 0, 0),
		},
		{
			desc:           "legacy form with a ttl",
			ttl:            time.Hour,
			expectedExpiry: now.Add(time.Hour),
		},
		{
			desc:           "service token lasts one year",
			service:        "batch",
			expectedExpiry: now.AddDate(1, 0, 0),
		},
		{
			desc:           "service token with a ttl",
			service:        "batch",
			ttl:            time.Hour,
			expectedExpiry: now.Add(time.Hour),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			issuerFlag, audienceFlag, emailFlag := issuer, audience, email
			serviceFlag, ttlFlag := p.service, p.ttl
			c := &command{
				issuer:   &issuerFlag,
				audience: &audienceFlag,
				email:    &emailFlag,
				service:  &serviceFlag,
				ttl:      &ttlFlag,
			}
			raw, err := c.sign(signer, now)
			require.NoError(t, err)
			if p.service == "" {
				accessToken, err := verifier.VerifyAccessToken(raw)
				require.NoError(t, err)
				assert.Equal(t, email, accessToken.Email)
				assert.True(t, accessToken.IsSystemAdmin)
				assert.True(t, accessToken.IsServiceToken)
				assert.True(t, p.expectedExpiry.Equal(accessToken.Expiry))
				_, err = verifier.VerifyServiceToken(raw)
				assert.Error(t, err)
				return
			}
			serviceToken, err := verifier.VerifyServiceToken(raw)
			require.NoError(t, err)
			assert.Equal(t, p.service, serviceToken.Subject)
			assert.Equal(t, email, serviceToken.Email)
			assert.NotEmpty(t, serviceToken.ID)
			assert.True(t, p.expectedExpiry.Equal(serviceToken.Expiry))
			_, err = verifier.VerifyAccessToken(raw)
			assert.Error(t, err)
		})
	}
}
