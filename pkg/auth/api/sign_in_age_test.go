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
	"go.uber.org/mock/gomock"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	tokentesting "github.com/bucketeer-io/bucketeer/v2/pkg/token/testing"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func TestWithMaxSignInAge(t *testing.T) {
	t.Parallel()
	opts := defaultOptions
	WithMaxSignInAge(48 * time.Hour)(&opts)
	assert.Equal(t, 48*time.Hour, opts.maxSignInAge)
	assert.Equal(t, 30*24*time.Hour, defaultOptions.maxSignInAge)
}

// newSignInAgeTestService returns a service with real tokens and a mocked
// account client that knows one enabled account.
func newSignInAgeTestService(
	t *testing.T,
	mockController *gomock.Controller,
	email string,
) (*authService, *accountclientmock.MockClient) {
	t.Helper()
	privateKeyPath, publicKeyPath := tokentesting.WriteKeyPair(t)
	signer, err := token.NewSigner(privateKeyPath)
	require.NoError(t, err)
	verifier, err := token.NewVerifier(publicKeyPath, "test-issuer", "test-audience")
	require.NoError(t, err)
	service := newAuthTypeTestService(t, nil)
	service.signer = signer
	service.verifier = verifier
	accountClient := accountclientmock.NewMockClient(mockController)
	service.accountClient = accountClient
	return service, accountClient
}

func expectEnabledAccount(ac *accountclientmock.MockClient, email string) {
	ac.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).Return(
		&acproto.GetMyOrganizationsResponse{
			Organizations: []*envproto.Organization{{Id: "org-1"}},
		}, nil)
	ac.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(
		&acproto.GetAccountV2Response{
			Account: &acproto.AccountV2{Email: email, OrganizationId: "org-1"},
		}, nil)
}

func TestAuthService_RefreshToken_SignInAge(t *testing.T) {
	t.Parallel()
	const email = "user@example.com"
	now := time.Now()
	patterns := []struct {
		desc           string
		issuedAt       time.Time
		authTime       time.Time
		expectedSignIn time.Time
		expectedErr    error
	}{
		{
			desc:           "success: within the age, the sign-in time is carried over",
			issuedAt:       now,
			authTime:       now.Add(-29 * 24 * time.Hour),
			expectedSignIn: now.Add(-29 * 24 * time.Hour),
		},
		{
			desc:        "err: past the age",
			issuedAt:    now,
			authTime:    now.Add(-31 * 24 * time.Hour),
			expectedErr: statusUnauthenticated.Err(),
		},
		{
			desc:           "success: without auth_time, the issued time counts as the sign-in",
			issuedAt:       now.Add(-time.Hour),
			expectedSignIn: now.Add(-time.Hour),
		},
		{
			desc:        "err: without auth_time and issued past the age",
			issuedAt:    now.Add(-31 * 24 * time.Hour),
			expectedErr: statusUnauthenticated.Err(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service, accountClient := newSignInAgeTestService(t, mockController, email)
			if p.expectedErr == nil {
				expectEnabledAccount(accountClient, email)
			}
			refreshToken, err := service.signer.SignRefreshToken(&token.RefreshToken{
				Email:    email,
				Expiry:   now.Add(time.Hour),
				IssuedAt: p.issuedAt,
				AuthTime: p.authTime,
			})
			require.NoError(t, err)

			resp, err := service.RefreshToken(context.Background(), &authproto.RefreshTokenRequest{RefreshToken: refreshToken})
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assertSignInTime(t, service, resp.Token, email, p.expectedSignIn)
		})
	}
}

func TestAuthService_SwitchOrganization_SignInAge(t *testing.T) {
	t.Parallel()
	const email = "user@example.com"
	now := time.Now()
	patterns := []struct {
		desc        string
		authTime    time.Time
		expectedErr error
	}{
		{
			desc:     "success: within the age, the sign-in time is carried over",
			authTime: now.Add(-29 * 24 * time.Hour),
		},
		{
			desc:        "err: past the age",
			authTime:    now.Add(-31 * 24 * time.Hour),
			expectedErr: statusUnauthenticated.Err(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service, accountClient := newSignInAgeTestService(t, mockController, email)
			if p.expectedErr == nil {
				expectEnabledAccount(accountClient, email)
			}
			accessToken, err := service.signer.SignAccessToken(&token.AccessToken{
				Issuer:   "test-issuer",
				Audience: "test-audience",
				Email:    email,
				Expiry:   now.Add(time.Hour),
				IssuedAt: now,
				AuthTime: p.authTime,
			})
			require.NoError(t, err)

			resp, err := service.SwitchOrganization(context.Background(), &authproto.SwitchOrganizationRequest{
				AccessToken:    accessToken,
				OrganizationId: "org-1",
			})
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			assertSignInTime(t, service, resp.Token, email, p.authTime)
		})
	}
}

// assertSignInTime checks that both issued tokens name the user and carry the sign-in time.
func assertSignInTime(t *testing.T, service *authService, issued *authproto.Token, email string, signedInAt time.Time) {
	t.Helper()
	access, err := service.verifier.VerifyAccessToken(issued.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, email, access.Subject)
	assert.True(t, signedInAt.Equal(access.AuthTime), "access auth_time %v, want %v", access.AuthTime, signedInAt)
	refresh, err := service.verifier.VerifyRefreshToken(issued.RefreshToken)
	require.NoError(t, err)
	assert.True(t, signedInAt.Equal(refresh.AuthTime), "refresh auth_time %v, want %v", refresh.AuthTime, signedInAt)
}
