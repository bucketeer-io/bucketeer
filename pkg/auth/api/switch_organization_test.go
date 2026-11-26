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

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	tokentesting "github.com/bucketeer-io/bucketeer/v2/pkg/token/testing"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func TestAuthService_SwitchOrganization_ServiceToken(t *testing.T) {
	t.Parallel()
	const email = "service@example.com"
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := token.NewSigner(privateKeyPath)
	require.NoError(t, err)
	service := newAuthTypeTestService(t, nil)
	service.verifier, err = token.NewVerifier(publicKeyPath, service.issuer, service.audience)
	require.NoError(t, err)
	now := time.Now()
	legacyServiceToken, err := signer.SignAccessToken(&token.AccessToken{
		Issuer:         service.issuer,
		Audience:       service.audience,
		Email:          email,
		Expiry:         now.Add(time.Hour),
		IssuedAt:       now,
		IsSystemAdmin:  true,
		IsServiceToken: true,
	})
	require.NoError(t, err)
	serviceToken, err := signer.SignServiceToken(&token.ServiceToken{
		Issuer:   service.issuer,
		Subject:  "batch",
		Audience: service.audience,
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
		ID:       "token-1",
	})
	require.NoError(t, err)

	patterns := []struct {
		desc        string
		accessToken string
		expectedErr error
	}{
		{
			desc:        "err: legacy service token is refused",
			accessToken: legacyServiceToken,
			expectedErr: statusAccessDenied.Err(),
		},
		{
			desc:        "err: service token does not verify as an access token",
			accessToken: serviceToken,
			expectedErr: statusUnauthenticated.Err(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			resp, err := service.SwitchOrganization(context.Background(), &authproto.SwitchOrganizationRequest{
				AccessToken:    p.accessToken,
				OrganizationId: "org-1",
			})
			require.Error(t, err)
			assert.Equal(t, p.expectedErr.Error(), err.Error())
			assert.Nil(t, resp)
		})
	}
}
