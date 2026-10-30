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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/google"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const (
	// Request values that must never be written to the log.
	testAuthorizationCode = "4/0AX4XfWgOauthAuthorizationCode1234567890"
	testOAuthState        = "oauth-state-9f8e7d6c5b4a"
)

// newObservedAuthService returns a service whose logger records every entry
// so that tests can inspect what would have been written.
func newObservedAuthService(
	t *testing.T,
	authenticator auth.Authenticator,
) (*authService, *observer.ObservedLogs) {
	t.Helper()
	core, logs := observer.New(zapcore.DebugLevel)
	service := newAuthTypeTestService(t, authenticator)
	service.logger = zap.New(core)
	return service, logs
}

// assertLogsOmit fails when any recorded entry mentions value in its message
// or in any of its fields.
func assertLogsOmit(t *testing.T, logs *observer.ObservedLogs, value string) {
	t.Helper()
	for _, entry := range logs.All() {
		assert.NotContains(t, entry.Message, value)
		for key, field := range entry.ContextMap() {
			assert.NotContains(t, fmt.Sprint(field), value, "field %q of %q", key, entry.Message)
		}
	}
}

// requireLogEntry checks that the error was logged once for the request and
// that the request context fields are still present.
func requireLogEntry(
	t *testing.T,
	logs *observer.ObservedLogs,
	message string,
	authType authproto.AuthType,
	redirectURL string,
) {
	t.Helper()
	entries := logs.FilterMessage(message).FilterField(zap.String("redirect_url", redirectURL)).All()
	require.Len(t, entries, 1, "expected exactly one %q entry for the request", message)
	assert.Equal(t, authType.String(), fmt.Sprint(entries[0].ContextMap()["type"]))
}

func TestAuthService_GetAuthenticationURL_LogsOmitState(t *testing.T) {
	t.Parallel()
	loginErr := status.Error(codes.Unavailable, "login failed")
	patterns := []struct {
		desc           string
		req            *authproto.GetAuthenticationURLRequest
		authenticator  auth.Authenticator
		expectedErr    error
		expectedLogMsg string
	}{
		{
			desc: "err: validation failure is logged without the state",
			req: &authproto.GetAuthenticationURLRequest{
				Type:  authproto.AuthType_AUTH_TYPE_GOOGLE,
				State: testOAuthState,
			},
			authenticator:  &fakeAuthenticator{},
			expectedErr:    statusMissingRedirectURL.Err(),
			expectedLogMsg: "Failed to validate the get authentication url request",
		},
		{
			desc: "err: unknown auth type is logged without the state",
			req: &authproto.GetAuthenticationURLRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GITHUB,
				State:       testOAuthState,
				RedirectUrl: testRedirectURL,
			},
			authenticator:  &fakeAuthenticator{},
			expectedErr:    statusUnknownAuthType.Err(),
			expectedLogMsg: "Failed to get the authenticator",
		},
		{
			desc: "err: login failure is logged without the state",
			req: &authproto.GetAuthenticationURLRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				State:       testOAuthState,
				RedirectUrl: testRedirectURL,
			},
			authenticator:  &fakeAuthenticator{err: loginErr},
			expectedErr:    loginErr,
			expectedLogMsg: "Failed to get the login url",
		},
		{
			desc: "success: the state is still forwarded to the provider",
			req: &authproto.GetAuthenticationURLRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				State:       testOAuthState,
				RedirectUrl: testRedirectURL,
			},
			authenticator: google.NewAuthenticator(&auth.GoogleConfig{
				ClientID:     "client-id",
				RedirectURLs: []string{testRedirectURL},
			}, zap.NewNop()),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service, logs := newObservedAuthService(t, p.authenticator)
			resp, err := service.GetAuthenticationURL(context.Background(), p.req)
			assertLogsOmit(t, logs, testOAuthState)
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
				requireLogEntry(t, logs, p.expectedLogMsg, p.req.Type, p.req.RedirectUrl)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			loginURL, err := url.Parse(resp.Url)
			require.NoError(t, err)
			assert.Equal(t, testOAuthState, loginURL.Query().Get("state"))
		})
	}
}

func TestAuthService_ExchangeToken_LogsOmitCode(t *testing.T) {
	t.Parallel()
	const email = "user@example.com"
	exchangeErr := status.Error(codes.Unavailable, "exchange failed")
	organizationsErr := status.Error(codes.Internal, "account service unavailable")
	patterns := []struct {
		desc           string
		req            *authproto.ExchangeTokenRequest
		authenticator  *fakeAuthenticator
		setup          func(*testing.T, *authService, *accountclientmock.MockClient)
		expectedErr    error
		expectedLogMsg string
	}{
		{
			desc: "err: validation failure is logged without the code",
			req: &authproto.ExchangeTokenRequest{
				Type: authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code: testAuthorizationCode,
			},
			authenticator:  &fakeAuthenticator{},
			expectedErr:    statusMissingRedirectURL.Err(),
			expectedLogMsg: "Failed to validate the exchange token request",
		},
		{
			desc: "err: unknown auth type is logged without the code",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GITHUB,
				Code:        testAuthorizationCode,
				RedirectUrl: testRedirectURL,
			},
			authenticator:  &fakeAuthenticator{},
			expectedErr:    statusUnknownAuthType.Err(),
			expectedLogMsg: "Failed to get the authenticator",
		},
		{
			desc: "err: exchange failure is logged without the code",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        testAuthorizationCode,
				RedirectUrl: testRedirectURL,
			},
			authenticator:  &fakeAuthenticator{err: exchangeErr},
			expectedErr:    exchangeErr,
			expectedLogMsg: "Failed to exchange",
		},
		{
			desc: "err: organization lookup failure is logged without the code",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        testAuthorizationCode,
				RedirectUrl: testRedirectURL,
			},
			authenticator: &fakeAuthenticator{userInfo: &auth.UserInfo{Email: email}},
			setup: func(t *testing.T, s *authService, ac *accountclientmock.MockClient) {
				t.Helper()
				ac.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).Return(nil, organizationsErr)
			},
			expectedErr:    organizationsErr,
			expectedLogMsg: "Failed to get organizations by email",
		},
		{
			desc: "success: the code is exchanged for a token",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        testAuthorizationCode,
				RedirectUrl: testRedirectURL,
			},
			authenticator: &fakeAuthenticator{userInfo: &auth.UserInfo{Email: email}},
			setup: func(t *testing.T, s *authService, ac *accountclientmock.MockClient) {
				t.Helper()
				signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
				require.NoError(t, err)
				s.signer = signer
				ac.EXPECT().GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).Return(
					&acproto.GetMyOrganizationsResponse{
						Organizations: []*envproto.Organization{{Id: "org-1"}},
					}, nil)
				// Called once to refresh the profile and once to check the account status.
				ac.EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(
					&acproto.GetAccountV2Response{
						Account: &acproto.AccountV2{
							Email:          email,
							OrganizationId: "org-1",
							LastSeen:       1,
						},
					}, nil).Times(2)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			service, logs := newObservedAuthService(t, p.authenticator)
			accountClient := accountclientmock.NewMockClient(mockController)
			service.accountClient = accountClient
			if p.setup != nil {
				p.setup(t, service, accountClient)
			}
			resp, err := service.ExchangeToken(context.Background(), p.req)
			assertLogsOmit(t, logs, testAuthorizationCode)
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
				requireLogEntry(t, logs, p.expectedLogMsg, p.req.Type, p.req.RedirectUrl)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Token)
			assert.NotEmpty(t, resp.Token.AccessToken)
			assert.Equal(t, 1, p.authenticator.calls)
		})
	}
}
