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

package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	tokentesting "github.com/bucketeer-io/bucketeer/v2/pkg/token/testing"
)

func TestAuthUnaryServerInterceptor(t *testing.T) {
	t.Parallel()
	const (
		issuer   = "test-issuer"
		audience = "test-audience"
		email    = "user@example.com"
	)
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := token.NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := token.NewVerifier(publicKeyPath, issuer, audience)
	require.NoError(t, err)
	now := time.Now()
	accessToken, err := signer.SignAccessToken(&token.AccessToken{
		Issuer:   issuer,
		Audience: audience,
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
	})
	require.NoError(t, err)
	legacyServiceToken, err := signer.SignAccessToken(&token.AccessToken{
		Issuer:         issuer,
		Audience:       audience,
		Email:          email,
		Expiry:         now.Add(time.Hour),
		IssuedAt:       now,
		IsSystemAdmin:  true,
		IsServiceToken: true,
	})
	require.NoError(t, err)
	serviceToken, err := signer.SignServiceToken(&token.ServiceToken{
		Issuer:   issuer,
		Subject:  "batch",
		Audience: audience,
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
		ID:       "token-1",
	})
	require.NoError(t, err)
	serviceTokenWithoutEmail, err := signer.SignServiceToken(&token.ServiceToken{
		Issuer:   issuer,
		Subject:  "batch",
		Audience: audience,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
		ID:       "token-2",
	})
	require.NoError(t, err)
	refreshToken, err := signer.SignRefreshToken(&token.RefreshToken{
		Email:    email,
		Expiry:   now.Add(time.Hour),
		IssuedAt: now,
	})
	require.NoError(t, err)

	patterns := []struct {
		desc          string
		authorization string
		expected      *token.AccessToken
		expectedCode  codes.Code
	}{
		{
			desc:          "success: access token",
			authorization: "bearer " + accessToken,
			expected: &token.AccessToken{
				Issuer:   issuer,
				Audience: audience,
				Email:    email,
				TokenUse: token.TokenUseAccess,
			},
		},
		{
			desc:          "success: legacy service token",
			authorization: "bearer " + legacyServiceToken,
			expected: &token.AccessToken{
				Issuer:         issuer,
				Audience:       audience,
				Email:          email,
				IsSystemAdmin:  true,
				IsServiceToken: true,
				TokenUse:       token.TokenUseAccess,
			},
		},
		{
			desc:          "success: service token is presented with the legacy flags",
			authorization: "bearer " + serviceToken,
			expected: &token.AccessToken{
				Issuer:         issuer,
				Audience:       audience,
				Email:          email,
				IsSystemAdmin:  true,
				IsServiceToken: true,
				TokenUse:       token.TokenUseService,
			},
		},
		{
			desc:          "success: service token without an email is identified by its subject",
			authorization: "bearer " + serviceTokenWithoutEmail,
			expected: &token.AccessToken{
				Issuer:         issuer,
				Audience:       audience,
				Email:          "batch",
				IsSystemAdmin:  true,
				IsServiceToken: true,
				TokenUse:       token.TokenUseService,
			},
		},
		{
			desc:          "err: refresh token",
			authorization: "bearer " + refreshToken,
			expectedCode:  codes.Unauthenticated,
		},
		{
			desc:          "err: malformed token",
			authorization: "bearer not-a-jwt",
			expectedCode:  codes.Unauthenticated,
		},
		{
			desc:         "err: missing token",
			expectedCode: codes.Unauthenticated,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			if p.authorization != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", p.authorization))
			}
			var got *token.AccessToken
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				got, _ = GetAccessToken(ctx)
				return nil, nil
			}
			info := &grpc.UnaryServerInfo{FullMethod: "/bucketeer.feature.FeatureService/GetFeature"}
			_, err := AuthUnaryServerInterceptor(verifier)(ctx, nil, info, handler)
			if p.expectedCode != codes.OK {
				require.Error(t, err)
				assert.Equal(t, p.expectedCode, status.Code(err))
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, p.expected.Issuer, got.Issuer)
			assert.Equal(t, p.expected.Audience, got.Audience)
			assert.Equal(t, p.expected.Email, got.Email)
			assert.Equal(t, p.expected.IsSystemAdmin, got.IsSystemAdmin)
			assert.Equal(t, p.expected.IsServiceToken, got.IsServiceToken)
			assert.Equal(t, p.expected.TokenUse, got.TokenUse)
			assert.True(t, now.Add(time.Hour).Equal(got.Expiry))
		})
	}
}
