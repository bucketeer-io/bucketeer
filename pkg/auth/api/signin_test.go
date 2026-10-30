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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/ratelimit"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const (
	demoEmail    = "demo@example.com"
	demoPassword = "demo-password"
)

func newSignInTestService(t *testing.T, enabled bool) *authService {
	t.Helper()
	signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
	require.NoError(t, err)
	opts := defaultOptions
	return &authService{
		issuer:   "test-issuer",
		audience: "test-audience",
		signer:   signer,
		config: &auth.OAuthConfig{
			DemoSignIn: auth.DemoSignInConfig{
				Enabled:        enabled,
				Email:          demoEmail,
				Password:       demoPassword,
				OrganizationId: "org-demo",
			},
		},
		opts:   &opts,
		logger: zap.NewNop(),
	}
}

func expectDemoAccountLookup(ac *accountclientmock.MockClient) {
	ac.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), &acproto.GetMyOrganizationsByEmailRequest{
		Email: demoEmail,
	}).Return(&acproto.GetMyOrganizationsResponse{
		Organizations: []*envproto.Organization{{Id: "org-demo"}},
	}, nil)
	ac.EXPECT().GetAccountV2(gomock.Any(), &acproto.GetAccountV2Request{
		Email:          demoEmail,
		OrganizationId: "org-demo",
	}).Return(&acproto.GetAccountV2Response{
		Account: &acproto.AccountV2{
			Email:          demoEmail,
			OrganizationId: "org-demo",
		},
	}, nil)
}

func TestAuthService_SignIn(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		enabled     bool
		req         *authproto.SignInRequest
		setup       func(*accountclientmock.MockClient)
		expectedErr error
	}{
		{
			desc:        "err: empty email",
			enabled:     true,
			req:         &authproto.SignInRequest{Password: demoPassword},
			expectedErr: statusMissingUsername.Err(),
		},
		{
			desc:        "err: empty password",
			enabled:     true,
			req:         &authproto.SignInRequest{Email: demoEmail},
			expectedErr: statusMissingPassword.Err(),
		},
		{
			desc:        "err: demo sign in disabled",
			enabled:     false,
			req:         &authproto.SignInRequest{Email: demoEmail, Password: demoPassword},
			expectedErr: statusAccessDenied.Err(),
		},
		{
			desc:        "err: wrong password",
			enabled:     true,
			req:         &authproto.SignInRequest{Email: demoEmail, Password: "wrong-password"},
			expectedErr: statusAccessDenied.Err(),
		},
		{
			desc:        "err: password with matching prefix",
			enabled:     true,
			req:         &authproto.SignInRequest{Email: demoEmail, Password: demoPassword + "x"},
			expectedErr: statusAccessDenied.Err(),
		},
		{
			desc:        "err: wrong email",
			enabled:     true,
			req:         &authproto.SignInRequest{Email: "other@example.com", Password: demoPassword},
			expectedErr: statusAccessDenied.Err(),
		},
		{
			desc:        "success: valid credentials",
			enabled:     true,
			req:         &authproto.SignInRequest{Email: demoEmail, Password: demoPassword},
			setup:       expectDemoAccountLookup,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service := newSignInTestService(t, p.enabled)
			accountClient := accountclientmock.NewMockClient(mockController)
			service.accountClient = accountClient
			if p.setup != nil {
				p.setup(accountClient)
			}
			resp, err := service.SignIn(context.Background(), p.req)
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Token)
			assert.NotEmpty(t, resp.Token.AccessToken)
			assert.NotEmpty(t, resp.Token.RefreshToken)
		})
	}
}

func TestAuthService_SignIn_RateLimit(t *testing.T) {
	t.Parallel()
	const burst = 3
	newLimiter := func(t *testing.T) *ratelimit.Limiter {
		t.Helper()
		// A very low refill rate so only the burst is available during the test.
		return ratelimit.NewLimiter(t.Context(), ratelimit.Config{RequestsPerMinute: 1, BurstSize: burst})
	}
	patterns := []struct {
		desc    string
		limiter func(*testing.T) *ratelimit.Limiter
		run     func(*testing.T, *authService, *accountclientmock.MockClient)
	}{
		{
			desc:    "err: repeated wrong passwords for the demo account are rejected after the burst",
			limiter: newLimiter,
			run: func(t *testing.T, s *authService, _ *accountclientmock.MockClient) {
				t.Helper()
				req := &authproto.SignInRequest{Email: demoEmail, Password: "wrong-password"}
				for i := 0; i < burst; i++ {
					_, err := s.SignIn(context.Background(), req)
					require.Error(t, err)
					assert.Equal(t, statusAccessDenied.Err().Error(), err.Error())
				}
				_, err := s.SignIn(context.Background(), req)
				require.Error(t, err)
				assert.Equal(t, statusTooManyRequests.Err().Error(), err.Error())
			},
		},
		{
			desc:    "err: attempts for unknown emails share one bucket regardless of the email",
			limiter: newLimiter,
			run: func(t *testing.T, s *authService, _ *accountclientmock.MockClient) {
				t.Helper()
				for i := 0; i < burst; i++ {
					_, err := s.SignIn(context.Background(), &authproto.SignInRequest{
						Email:    fmt.Sprintf("user-%d@example.com", i),
						Password: "wrong-password",
					})
					require.Error(t, err)
					assert.Equal(t, statusAccessDenied.Err().Error(), err.Error())
				}
				_, err := s.SignIn(context.Background(), &authproto.SignInRequest{
					Email:    "another@example.com",
					Password: "wrong-password",
				})
				require.Error(t, err)
				assert.Equal(t, statusTooManyRequests.Err().Error(), err.Error())
			},
		},
		{
			desc:    "err: valid credentials are also rejected once the demo account is throttled",
			limiter: newLimiter,
			run: func(t *testing.T, s *authService, ac *accountclientmock.MockClient) {
				t.Helper()
				wrong := &authproto.SignInRequest{Email: demoEmail, Password: "wrong-password"}
				for i := 0; i < burst; i++ {
					_, err := s.SignIn(context.Background(), wrong)
					require.Error(t, err)
				}
				// No account lookup is expected because the request is rejected first.
				valid := &authproto.SignInRequest{Email: demoEmail, Password: demoPassword}
				_, err := s.SignIn(context.Background(), valid)
				require.Error(t, err)
				assert.Equal(t, statusTooManyRequests.Err().Error(), err.Error())
			},
		},
		{
			desc:    "success: attempts for unknown emails do not throttle the demo account",
			limiter: newLimiter,
			run: func(t *testing.T, s *authService, ac *accountclientmock.MockClient) {
				t.Helper()
				other := &authproto.SignInRequest{Email: "other@example.com", Password: "wrong-password"}
				for i := 0; i < burst+1; i++ {
					_, err := s.SignIn(context.Background(), other)
					require.Error(t, err)
				}
				expectDemoAccountLookup(ac)
				resp, err := s.SignIn(context.Background(), &authproto.SignInRequest{
					Email:    demoEmail,
					Password: demoPassword,
				})
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token.AccessToken)
			},
		},
		{
			desc:    "success: no limiter configured keeps sign in unthrottled",
			limiter: func(*testing.T) *ratelimit.Limiter { return nil },
			run: func(t *testing.T, s *authService, ac *accountclientmock.MockClient) {
				t.Helper()
				wrong := &authproto.SignInRequest{Email: demoEmail, Password: "wrong-password"}
				for i := 0; i < burst*2; i++ {
					_, err := s.SignIn(context.Background(), wrong)
					require.Error(t, err)
					assert.Equal(t, statusAccessDenied.Err().Error(), err.Error())
				}
				expectDemoAccountLookup(ac)
				resp, err := s.SignIn(context.Background(), &authproto.SignInRequest{
					Email:    demoEmail,
					Password: demoPassword,
				})
				require.NoError(t, err)
				require.NotNil(t, resp)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service := newSignInTestService(t, true)
			service.opts.signInRateLimiter = p.limiter(t)
			accountClient := accountclientmock.NewMockClient(mockController)
			service.accountClient = accountClient
			p.run(t, service, accountClient)
		})
	}
}
