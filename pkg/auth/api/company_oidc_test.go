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
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
)

func TestCreateTemporaryToken(t *testing.T) {
	t.Parallel()

	// Create test service
	issuer := "test-issuer"
	audience := "test-audience"
	signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
	require.NoError(t, err)
	verifier, err := token.NewVerifier("../../token/testdata/valid-public.pem", issuer, audience)
	require.NoError(t, err)

	service := &authService{
		issuer:   issuer,
		audience: audience,
		signer:   signer,
		verifier: verifier,
		opts:     &defaultOptions,
		logger:   zap.NewNop(),
	}

	ctx := context.Background()
	userInfo := &auth.UserInfo{
		Email:         "test@example.com",
		Name:          "Test User",
		VerifiedEmail: true,
	}

	// Create temporary token
	token, err := service.createTemporaryToken(ctx, userInfo)
	require.NoError(t, err)
	require.NotNil(t, token)

	// Verify token properties
	assert.NotEmpty(t, token.AccessToken, "access token should not be empty")
	assert.Empty(t, token.RefreshToken, "refresh token should be empty for temporary tokens")
	assert.Equal(t, "Bearer", token.TokenType)
	assert.Greater(t, token.Expiry, time.Now().Unix(), "expiry should be in the future")

	// Verify token expires in approximately 5 minutes
	expectedExpiry := time.Now().Add(5 * time.Minute).Unix()
	assert.InDelta(t, expectedExpiry, token.Expiry, 5, "expiry should be ~5 minutes from now")

	// Verify token can be decoded
	accessToken, err := verifier.VerifyAccessToken(token.AccessToken)
	require.NoError(t, err)
	require.NotNil(t, accessToken)

	// Verify token claims
	assert.Equal(t, issuer, accessToken.Issuer)
	assert.Equal(t, audience, accessToken.Audience)
	assert.Equal(t, userInfo.Email, accessToken.Email)
	assert.Equal(t, userInfo.Name, accessToken.Name)
	assert.Empty(t, accessToken.OrganizationID, "organization ID should be empty for temporary token")
	assert.False(t, accessToken.IsSystemAdmin, "should not be system admin")

	// Verify token expiry
	assert.True(t, accessToken.Expiry.After(time.Now()), "token should not be expired")
	assert.True(t, accessToken.Expiry.Before(time.Now().Add(6*time.Minute)), "token should expire within 6 minutes")
}

func TestCreateTemporaryToken_TokenStructure(t *testing.T) {
	t.Parallel()

	issuer := "test-issuer"
	audience := "test-audience"
	signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
	require.NoError(t, err)
	verifier, err := token.NewVerifier("../../token/testdata/valid-public.pem", issuer, audience)
	require.NoError(t, err)

	service := &authService{
		issuer:   issuer,
		audience: audience,
		signer:   signer,
		verifier: verifier,
		opts:     &defaultOptions,
		logger:   zap.NewNop(),
	}

	ctx := context.Background()
	userInfo := &auth.UserInfo{
		Email:         "test@example.com",
		Name:          "Test User",
		FirstName:     "Test",
		LastName:      "User",
		Avatar:        "https://example.com/avatar.jpg",
		VerifiedEmail: true,
	}

	// Create token
	tokenResp, err := service.createTemporaryToken(ctx, userInfo)
	require.NoError(t, err)

	// Decode and verify structure
	accessToken, err := verifier.VerifyAccessToken(tokenResp.AccessToken)
	require.NoError(t, err)

	// Verify all required fields are present
	assert.NotEmpty(t, accessToken.Issuer)
	assert.NotEmpty(t, accessToken.Audience)
	assert.NotEmpty(t, accessToken.Email)
	assert.NotEmpty(t, accessToken.Name)
	assert.False(t, accessToken.Expiry.IsZero())
	assert.False(t, accessToken.IssuedAt.IsZero())

	// Verify temporary token specific fields
	assert.Empty(t, accessToken.OrganizationID, "temporary token should not have organization ID")
	assert.False(t, accessToken.IsSystemAdmin, "temporary token should not have system admin flag")
}
