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
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/aichat/ratelimit"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const testBurst = 3

// newTestTokenLimiter returns a limiter that only allows the burst during the
// test, since the refill rate is one request per minute.
func newTestTokenLimiter(t *testing.T) *ratelimit.Limiter {
	t.Helper()
	return ratelimit.NewLimiter(t.Context(), ratelimit.Config{RequestsPerMinute: 1, BurstSize: testBurst})
}

func contextWithClientAddress(addr string) context.Context {
	return metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(xForwardedForKey, addr),
	)
}

func TestClientAddressFromContext(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		ctx      context.Context
		expected string
	}{
		{
			desc:     "no metadata and no peer",
			ctx:      context.Background(),
			expected: unknownClientAddress,
		},
		{
			desc: "single forwarded address",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs(xForwardedForKey, "203.0.113.10")),
			expected: "203.0.113.10",
		},
		{
			desc: "first address of a forwarded chain",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs(xForwardedForKey, "203.0.113.10, 10.0.0.1, 127.0.0.1")),
			expected: "203.0.113.10",
		},
		{
			desc: "first forwarded value that is an IP address",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs(xForwardedForKey, " ", xForwardedForKey, "198.51.100.7")),
			expected: "198.51.100.7",
		},
		{
			desc: "IPv6 forwarded address is normalized",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs(xForwardedForKey, "2001:DB8::1, 10.0.0.1")),
			expected: "2001:db8::1",
		},
		{
			desc: "forwarded value that is not an IP address falls back to the peer",
			ctx: peer.NewContext(
				metadata.NewIncomingContext(context.Background(),
					metadata.Pairs(xForwardedForKey, strings.Repeat("x", 4096)+", 10.0.0.1")),
				&peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("192.0.2.5"), Port: 50051}},
			),
			expected: "192.0.2.5",
		},
		{
			desc: "forwarded value that is not an IP address without a peer",
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.Pairs(xForwardedForKey, "not-an-address")),
			expected: unknownClientAddress,
		},
		{
			desc: "peer address without forwarded metadata",
			ctx: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.TCPAddr{IP: net.ParseIP("192.0.2.5"), Port: 50051},
			}),
			expected: "192.0.2.5",
		},
		{
			desc: "forwarded metadata takes precedence over the peer",
			ctx: peer.NewContext(
				metadata.NewIncomingContext(context.Background(),
					metadata.Pairs(xForwardedForKey, "203.0.113.10")),
				&peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("192.0.2.5"), Port: 50051}},
			),
			expected: "203.0.113.10",
		},
		{
			desc: "empty forwarded metadata falls back to the peer",
			ctx: peer.NewContext(
				metadata.NewIncomingContext(context.Background(),
					metadata.Pairs(xForwardedForKey, "")),
				&peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("192.0.2.5"), Port: 50051}},
			),
			expected: "192.0.2.5",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, clientAddressFromContext(p.ctx))
		})
	}
}

func TestAuthService_ExchangeToken_RateLimit(t *testing.T) {
	t.Parallel()
	const email = "user@example.com"
	req := &authproto.ExchangeTokenRequest{
		Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
		Code:        "code",
		RedirectUrl: testRedirectURL,
	}
	expectAccountLookup := func(ac *accountclientmock.MockClient) {
		ac.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).Return(
			&acproto.GetMyOrganizationsResponse{
				Organizations: []*envproto.Organization{{Id: "org-1"}},
			}, nil)
		ac.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(
			&acproto.GetAccountV2Response{
				Account: &acproto.AccountV2{Email: email, OrganizationId: "org-1", LastSeen: 1},
			}, nil).Times(2)
	}
	patterns := []struct {
		desc    string
		limiter func(*testing.T) *ratelimit.Limiter
		run     func(*testing.T, *authService, *fakeAuthenticator, *accountclientmock.MockClient)
	}{
		{
			desc:    "err: requests from one client are rejected after the burst without calling the authenticator",
			limiter: newTestTokenLimiter,
			run: func(t *testing.T, s *authService, a *fakeAuthenticator, ac *accountclientmock.MockClient) {
				t.Helper()
				ctx := contextWithClientAddress("203.0.113.10")
				for i := 0; i < testBurst; i++ {
					expectAccountLookup(ac)
					_, err := s.ExchangeToken(ctx, req)
					require.NoError(t, err)
				}
				_, err := s.ExchangeToken(ctx, req)
				require.Error(t, err)
				assert.Equal(t, statusTooManyRequests.Err().Error(), err.Error())
				assert.Equal(t, testBurst, a.calls)
			},
		},
		{
			desc:    "success: requests are counted per client address",
			limiter: newTestTokenLimiter,
			run: func(t *testing.T, s *authService, a *fakeAuthenticator, ac *accountclientmock.MockClient) {
				t.Helper()
				for i := 0; i < testBurst; i++ {
					expectAccountLookup(ac)
					_, err := s.ExchangeToken(contextWithClientAddress("203.0.113.10"), req)
					require.NoError(t, err)
				}
				expectAccountLookup(ac)
				resp, err := s.ExchangeToken(contextWithClientAddress("198.51.100.7"), req)
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token.AccessToken)
			},
		},
		{
			desc:    "success: no limiter configured keeps the exchange unthrottled",
			limiter: func(*testing.T) *ratelimit.Limiter { return nil },
			run: func(t *testing.T, s *authService, a *fakeAuthenticator, ac *accountclientmock.MockClient) {
				t.Helper()
				ctx := contextWithClientAddress("203.0.113.10")
				for i := 0; i < testBurst*2; i++ {
					expectAccountLookup(ac)
					_, err := s.ExchangeToken(ctx, req)
					require.NoError(t, err)
				}
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			authenticator := &fakeAuthenticator{userInfo: &auth.UserInfo{Email: email}}
			service := newAuthTypeTestService(t, authenticator)
			signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
			require.NoError(t, err)
			service.signer = signer
			service.opts.tokenRateLimiter = p.limiter(t)
			accountClient := accountclientmock.NewMockClient(mockController)
			service.accountClient = accountClient
			p.run(t, service, authenticator, accountClient)
		})
	}
}

func TestAuthService_RefreshToken_RateLimit(t *testing.T) {
	t.Parallel()
	const email = "user@example.com"
	expectAccountLookup := func(ac *accountclientmock.MockClient) {
		ac.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).Return(
			&acproto.GetMyOrganizationsResponse{
				Organizations: []*envproto.Organization{{Id: "org-1"}},
			}, nil)
		ac.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(
			&acproto.GetAccountV2Response{
				Account: &acproto.AccountV2{Email: email, OrganizationId: "org-1"},
			}, nil)
	}
	patterns := []struct {
		desc    string
		limiter func(*testing.T) *ratelimit.Limiter
		run     func(*testing.T, *authService, string, *accountclientmock.MockClient)
	}{
		{
			desc:    "err: invalid tokens from one client are rejected after the burst",
			limiter: newTestTokenLimiter,
			run: func(t *testing.T, s *authService, _ string, _ *accountclientmock.MockClient) {
				t.Helper()
				ctx := contextWithClientAddress("203.0.113.10")
				req := &authproto.RefreshTokenRequest{RefreshToken: "not-a-token"}
				for i := 0; i < testBurst; i++ {
					_, err := s.RefreshToken(ctx, req)
					require.Error(t, err)
					assert.Equal(t, statusUnauthenticated.Err().Error(), err.Error())
				}
				_, err := s.RefreshToken(ctx, req)
				require.Error(t, err)
				assert.Equal(t, statusTooManyRequests.Err().Error(), err.Error())
			},
		},
		{
			desc:    "success: a valid token is accepted even when the client address is throttled",
			limiter: newTestTokenLimiter,
			run: func(t *testing.T, s *authService, refreshToken string, ac *accountclientmock.MockClient) {
				t.Helper()
				ctx := contextWithClientAddress("203.0.113.10")
				for i := 0; i < testBurst+1; i++ {
					_, err := s.RefreshToken(ctx, &authproto.RefreshTokenRequest{RefreshToken: "not-a-token"})
					require.Error(t, err)
				}
				expectAccountLookup(ac)
				resp, err := s.RefreshToken(ctx, &authproto.RefreshTokenRequest{RefreshToken: refreshToken})
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token.AccessToken)
			},
		},
		{
			desc:    "success: valid tokens do not consume the limit",
			limiter: newTestTokenLimiter,
			run: func(t *testing.T, s *authService, refreshToken string, ac *accountclientmock.MockClient) {
				t.Helper()
				ctx := contextWithClientAddress("203.0.113.10")
				for i := 0; i < testBurst*2; i++ {
					expectAccountLookup(ac)
					_, err := s.RefreshToken(ctx, &authproto.RefreshTokenRequest{RefreshToken: refreshToken})
					require.NoError(t, err)
				}
				// The full burst is still available for rejected tokens.
				for i := 0; i < testBurst; i++ {
					_, err := s.RefreshToken(ctx, &authproto.RefreshTokenRequest{RefreshToken: "not-a-token"})
					require.Error(t, err)
					assert.Equal(t, statusUnauthenticated.Err().Error(), err.Error())
				}
			},
		},
		{
			desc:    "success: rejected tokens are counted per client address",
			limiter: newTestTokenLimiter,
			run: func(t *testing.T, s *authService, _ string, _ *accountclientmock.MockClient) {
				t.Helper()
				req := &authproto.RefreshTokenRequest{RefreshToken: "not-a-token"}
				for i := 0; i < testBurst+1; i++ {
					_, err := s.RefreshToken(contextWithClientAddress("203.0.113.10"), req)
					require.Error(t, err)
				}
				// The other address still has its own full burst before it is throttled.
				other := contextWithClientAddress("198.51.100.7")
				for i := 0; i < testBurst; i++ {
					_, err := s.RefreshToken(other, req)
					require.Error(t, err)
					assert.Equal(t, statusUnauthenticated.Err().Error(), err.Error())
				}
				_, err := s.RefreshToken(other, req)
				require.Error(t, err)
				assert.Equal(t, statusTooManyRequests.Err().Error(), err.Error())
			},
		},
		{
			desc:    "success: no limiter configured keeps the refresh unthrottled",
			limiter: func(*testing.T) *ratelimit.Limiter { return nil },
			run: func(t *testing.T, s *authService, refreshToken string, ac *accountclientmock.MockClient) {
				t.Helper()
				ctx := contextWithClientAddress("203.0.113.10")
				for i := 0; i < testBurst*2; i++ {
					expectAccountLookup(ac)
					_, err := s.RefreshToken(ctx, &authproto.RefreshTokenRequest{RefreshToken: refreshToken})
					require.NoError(t, err)
				}
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
			require.NoError(t, err)
			verifier, err := token.NewVerifier("../../token/testdata/valid-public.pem", "test-issuer", "test-audience")
			require.NoError(t, err)
			now := time.Now()
			refreshToken, err := signer.SignRefreshToken(&token.RefreshToken{
				Email:          email,
				OrganizationID: "org-1",
				Expiry:         now.Add(time.Hour),
				IssuedAt:       now,
			})
			require.NoError(t, err)
			service := newAuthTypeTestService(t, nil)
			service.signer = signer
			service.verifier = verifier
			service.opts.tokenRateLimiter = p.limiter(t)
			accountClient := accountclientmock.NewMockClient(mockController)
			service.accountClient = accountClient
			p.run(t, service, refreshToken, accountClient)
		})
	}
}
