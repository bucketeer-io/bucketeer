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

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func TestAuthService_GetDeploymentStatus(t *testing.T) {
	t.Parallel()
	options := defaultOptions
	service := &authService{
		logger: options.logger,
		opts:   &options,
	}
	patterns := []struct {
		desc        string
		setup       func(s *authService)
		expectedErr error
		expected    *authproto.GetDemoSiteStatusResponse
	}{
		{
			desc: "success: true",
			setup: func(s *authService) {
				s.opts.isDemoSiteEnabled = true
			},
			expectedErr: nil,
			expected: &authproto.GetDemoSiteStatusResponse{
				IsDemoSiteEnabled: true,
			},
		},
		{
			desc: "success: false",
			setup: func(s *authService) {
				s.opts.isDemoSiteEnabled = false
			},
			expectedErr: nil,
			expected: &authproto.GetDemoSiteStatusResponse{
				IsDemoSiteEnabled: false,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			p.setup(service)
			resp, err := service.GetDemoSiteStatus(context.Background(), nil)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, resp)
		})
	}
}

func TestWithAccessTokenTTL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		inputFunc    func() Option
		expectedFunc func(*testing.T, *options)
		expectedErr  error
	}{
		{
			desc: "success: set custom access token TTL",
			inputFunc: func() Option {
				return WithAccessTokenTTL(30 * time.Minute)
			},
			expectedFunc: func(t *testing.T, opts *options) {
				t.Helper()
				assert.Equal(t, 30*time.Minute, opts.accessTokenTTL)
			},
			expectedErr: nil,
		},
		{
			desc: "success: set access token TTL to 1 hour",
			inputFunc: func() Option {
				return WithAccessTokenTTL(1 * time.Hour)
			},
			expectedFunc: func(t *testing.T, opts *options) {
				t.Helper()
				assert.Equal(t, 1*time.Hour, opts.accessTokenTTL)
			},
			expectedErr: nil,
		},
		{
			desc: "success: set access token TTL to 5 minutes",
			inputFunc: func() Option {
				return WithAccessTokenTTL(5 * time.Minute)
			},
			expectedFunc: func(t *testing.T, opts *options) {
				t.Helper()
				assert.Equal(t, 5*time.Minute, opts.accessTokenTTL)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			opts := defaultOptions
			opt := p.inputFunc()
			opt(&opts)
			if p.expectedFunc != nil {
				p.expectedFunc(t, &opts)
			}
		})
	}
}

func TestWithRefreshTokenTTL(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		inputFunc    func() Option
		expectedFunc func(*testing.T, *options)
		expectedErr  error
	}{
		{
			desc: "success: set custom refresh token TTL",
			inputFunc: func() Option {
				return WithRefreshTokenTTL(14 * 24 * time.Hour)
			},
			expectedFunc: func(t *testing.T, opts *options) {
				t.Helper()
				assert.Equal(t, 14*24*time.Hour, opts.refreshTokenTTL)
			},
			expectedErr: nil,
		},
		{
			desc: "success: set refresh token TTL to 30 days",
			inputFunc: func() Option {
				return WithRefreshTokenTTL(30 * 24 * time.Hour)
			},
			expectedFunc: func(t *testing.T, opts *options) {
				t.Helper()
				assert.Equal(t, 30*24*time.Hour, opts.refreshTokenTTL)
			},
			expectedErr: nil,
		},
		{
			desc: "success: set refresh token TTL to 1 day",
			inputFunc: func() Option {
				return WithRefreshTokenTTL(24 * time.Hour)
			},
			expectedFunc: func(t *testing.T, opts *options) {
				t.Helper()
				assert.Equal(t, 24*time.Hour, opts.refreshTokenTTL)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			opts := defaultOptions
			opt := p.inputFunc()
			opt(&opts)
			if p.expectedFunc != nil {
				p.expectedFunc(t, &opts)
			}
		})
	}
}

func TestNewAuthService_WithTokenTTLs(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		setupFunc    func() []Option
		expectedFunc func(*testing.T, *authService)
		expectedErr  error
	}{
		{
			desc: "success: default TTLs",
			setupFunc: func() []Option {
				return []Option{}
			},
			expectedFunc: func(t *testing.T, s *authService) {
				t.Helper()
				assert.Equal(t, 10*time.Minute, s.opts.accessTokenTTL)
				assert.Equal(t, 7*24*time.Hour, s.opts.refreshTokenTTL)
			},
			expectedErr: nil,
		},
		{
			desc: "success: custom access token TTL only",
			setupFunc: func() []Option {
				return []Option{
					WithAccessTokenTTL(30 * time.Minute),
				}
			},
			expectedFunc: func(t *testing.T, s *authService) {
				t.Helper()
				assert.Equal(t, 30*time.Minute, s.opts.accessTokenTTL)
				assert.Equal(t, 7*24*time.Hour, s.opts.refreshTokenTTL) // default refresh token TTL
			},
			expectedErr: nil,
		},
		{
			desc: "success: custom refresh token TTL only",
			setupFunc: func() []Option {
				return []Option{
					WithRefreshTokenTTL(14 * 24 * time.Hour),
				}
			},
			expectedFunc: func(t *testing.T, s *authService) {
				t.Helper()
				assert.Equal(t, 10*time.Minute, s.opts.accessTokenTTL) // default access token TTL
				assert.Equal(t, 14*24*time.Hour, s.opts.refreshTokenTTL)
			},
			expectedErr: nil,
		},
		{
			desc: "success: both custom TTLs",
			setupFunc: func() []Option {
				return []Option{
					WithAccessTokenTTL(1 * time.Hour),
					WithRefreshTokenTTL(30 * 24 * time.Hour),
				}
			},
			expectedFunc: func(t *testing.T, s *authService) {
				t.Helper()
				assert.Equal(t, 1*time.Hour, s.opts.accessTokenTTL)
				assert.Equal(t, 30*24*time.Hour, s.opts.refreshTokenTTL)
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mysqlClient := mysqlmock.NewMockClient(mockController)
			accountClient := accountclientmock.NewMockClient(mockController)
			signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
			require.NoError(t, err)
			verifier, err := token.NewVerifier("../../token/testdata/valid-public.pem", "test-issuer", "test-audience")
			require.NoError(t, err)
			config := &auth.OAuthConfig{}

			service := NewAuthService(
				"test-issuer",
				"test-audience",
				signer,
				verifier,
				mysqlClient,
				accountClient,
				config,
				p.setupFunc()...,
			).(*authService)

			if p.expectedFunc != nil {
				p.expectedFunc(t, service)
			}
		})
	}
}

func TestAuthService_GenerateToken_WithCustomTTLs(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc         string
		setupFunc    func() []Option
		inputFunc    func() (context.Context, string, domain.AccountV2, bool, locale.Localizer)
		expectedFunc func(*testing.T, *authproto.Token, time.Duration, time.Duration)
		expectedErr  error
	}{
		{
			desc: "success: default TTLs",
			setupFunc: func() []Option {
				return []Option{}
			},
			inputFunc: func() (context.Context, string, domain.AccountV2, bool, locale.Localizer) {
				ctx := context.Background()
				email := "test@example.com"
				account := domain.AccountV2{
					AccountV2: &acproto.AccountV2{
						Email:          email,
						OrganizationId: "org-1",
						FirstName:      "Test",
						LastName:       "User",
					},
				}
				return ctx, email, account, false, locale.NewLocalizer(ctx)
			},
			expectedFunc: func(t *testing.T, token *authproto.Token, accessTokenTTL, refreshTokenTTL time.Duration) {
				t.Helper()
				require.NotNil(t, token)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)
				assert.Equal(t, "Bearer", token.TokenType)

				// Verify expiry time is approximately correct (within 1 second tolerance)
				expectedExpiry := time.Now().Add(accessTokenTTL).Unix()
				actualExpiry := token.Expiry
				assert.InDelta(t, expectedExpiry, actualExpiry, 1, "Token expiry should match configured TTL")
			},
			expectedErr: nil,
		},
		{
			desc: "success: custom access token TTL (30 minutes)",
			setupFunc: func() []Option {
				return []Option{
					WithAccessTokenTTL(30 * time.Minute),
				}
			},
			inputFunc: func() (context.Context, string, domain.AccountV2, bool, locale.Localizer) {
				ctx := context.Background()
				email := "test@example.com"
				account := domain.AccountV2{
					AccountV2: &acproto.AccountV2{
						Email:          email,
						OrganizationId: "org-1",
						FirstName:      "Test",
						LastName:       "User",
					},
				}
				return ctx, email, account, false, locale.NewLocalizer(ctx)
			},
			expectedFunc: func(t *testing.T, token *authproto.Token, accessTokenTTL, refreshTokenTTL time.Duration) {
				t.Helper()
				require.NotNil(t, token)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)

				// Verify expiry time matches the custom TTL
				expectedExpiry := time.Now().Add(accessTokenTTL).Unix()
				actualExpiry := token.Expiry
				assert.InDelta(t, expectedExpiry, actualExpiry, 1, "Token expiry should match configured TTL")
			},
			expectedErr: nil,
		},
		{
			desc: "success: custom refresh token TTL (14 days)",
			setupFunc: func() []Option {
				return []Option{
					WithRefreshTokenTTL(14 * 24 * time.Hour),
				}
			},
			inputFunc: func() (context.Context, string, domain.AccountV2, bool, locale.Localizer) {
				ctx := context.Background()
				email := "test@example.com"
				account := domain.AccountV2{
					AccountV2: &acproto.AccountV2{
						Email:          email,
						OrganizationId: "org-1",
						FirstName:      "Test",
						LastName:       "User",
					},
				}
				return ctx, email, account, false, locale.NewLocalizer(ctx)
			},
			expectedFunc: func(t *testing.T, token *authproto.Token, accessTokenTTL, refreshTokenTTL time.Duration) {
				t.Helper()
				require.NotNil(t, token)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)

				// Verify access token expiry (should use default)
				expectedExpiry := time.Now().Add(accessTokenTTL).Unix()
				actualExpiry := token.Expiry
				assert.InDelta(t, expectedExpiry, actualExpiry, 1, "Access token expiry should match configured TTL")
			},
			expectedErr: nil,
		},
		{
			desc: "success: both custom TTLs",
			setupFunc: func() []Option {
				return []Option{
					WithAccessTokenTTL(1 * time.Hour),
					WithRefreshTokenTTL(30 * 24 * time.Hour),
				}
			},
			inputFunc: func() (context.Context, string, domain.AccountV2, bool, locale.Localizer) {
				ctx := context.Background()
				email := "test@example.com"
				account := domain.AccountV2{
					AccountV2: &acproto.AccountV2{
						Email:          email,
						OrganizationId: "org-1",
						FirstName:      "Test",
						LastName:       "User",
					},
				}
				return ctx, email, account, false, locale.NewLocalizer(ctx)
			},
			expectedFunc: func(t *testing.T, token *authproto.Token, accessTokenTTL, refreshTokenTTL time.Duration) {
				t.Helper()
				require.NotNil(t, token)
				assert.NotEmpty(t, token.AccessToken)
				assert.NotEmpty(t, token.RefreshToken)

				// Verify access token expiry matches custom TTL
				expectedExpiry := time.Now().Add(accessTokenTTL).Unix()
				actualExpiry := token.Expiry
				assert.InDelta(t, expectedExpiry, actualExpiry, 1, "Access token expiry should match configured TTL")
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mysqlClient := mysqlmock.NewMockClient(mockController)
			accountClient := accountclientmock.NewMockClient(mockController)
			signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
			require.NoError(t, err)
			verifier, err := token.NewVerifier("../../token/testdata/valid-public.pem", "test-issuer", "test-audience")
			require.NoError(t, err)
			config := &auth.OAuthConfig{}

			opts := p.setupFunc()
			service := NewAuthService(
				"test-issuer",
				"test-audience",
				signer,
				verifier,
				mysqlClient,
				accountClient,
				config,
				opts...,
			).(*authService)

			// Get the configured TTLs for verification
			accessTokenTTL := service.opts.accessTokenTTL
			refreshTokenTTL := service.opts.refreshTokenTTL

			ctx, email, account, isSystemAdmin, localizer := p.inputFunc()
			token, err := service.generateToken(ctx, email, account, isSystemAdmin, localizer)

			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				if p.expectedFunc != nil {
					p.expectedFunc(t, token, accessTokenTTL, refreshTokenTTL)
				}
			}
		})
	}
}
