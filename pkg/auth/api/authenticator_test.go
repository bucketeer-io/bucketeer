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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/google"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const testRedirectURL = "https://console.example.com/auth/callback"

// fakeAuthenticator stands in for the Google authenticator so that the token
// exchange can be exercised without a network round trip.
type fakeAuthenticator struct {
	userInfo *auth.UserInfo
	calls    int
}

func (a *fakeAuthenticator) Login(context.Context, string, string) (string, error) {
	a.calls++
	return "https://accounts.example.com/login", nil
}

func (a *fakeAuthenticator) Exchange(context.Context, string, string) (*auth.UserInfo, error) {
	a.calls++
	return a.userInfo, nil
}

func newAuthTypeTestService(t *testing.T, authenticator auth.Authenticator) *authService {
	t.Helper()
	opts := defaultOptions
	return &authService{
		issuer:              "test-issuer",
		audience:            "test-audience",
		googleAuthenticator: authenticator,
		opts:                &opts,
		logger:              zap.NewNop(),
	}
}

func TestAuthService_GetAuthenticator(t *testing.T) {
	t.Parallel()
	googleAuthenticator := google.NewAuthenticator(&auth.GoogleConfig{}, zap.NewNop())
	service := newAuthTypeTestService(t, googleAuthenticator)
	patterns := []struct {
		desc        string
		authType    authproto.AuthType
		expected    auth.Authenticator
		expectedErr error
	}{
		{
			desc:     "success: google returns the configured authenticator",
			authType: authproto.AuthType_AUTH_TYPE_GOOGLE,
			expected: googleAuthenticator,
		},
		{
			desc:        "err: github is not supported",
			authType:    authproto.AuthType_AUTH_TYPE_GITHUB,
			expectedErr: statusUnknownAuthType.Err(),
		},
		{
			desc:        "err: deprecated user password is not supported",
			authType:    authproto.AuthType_AUTH_TYPE_USER_PASSWORD,
			expectedErr: statusUnknownAuthType.Err(),
		},
		{
			desc:        "err: unspecified is not supported",
			authType:    authproto.AuthType_AUTH_TYPE_UNSPECIFIED,
			expectedErr: statusUnknownAuthType.Err(),
		},
		{
			desc:        "err: undefined enum value is not supported",
			authType:    authproto.AuthType(99),
			expectedErr: statusUnknownAuthType.Err(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := service.getAuthenticator(p.authType)
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, actual)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, actual)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestAuthService_GetAuthenticationURL_AuthType(t *testing.T) {
	t.Parallel()
	googleAuthenticator := google.NewAuthenticator(&auth.GoogleConfig{
		ClientID:     "client-id",
		RedirectURLs: []string{testRedirectURL},
	}, zap.NewNop())
	service := newAuthTypeTestService(t, googleAuthenticator)
	patterns := []struct {
		desc        string
		req         *authproto.GetAuthenticationURLRequest
		expectedErr error
	}{
		{
			desc: "err: github auth type is rejected",
			req: &authproto.GetAuthenticationURLRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GITHUB,
				State:       "state",
				RedirectUrl: testRedirectURL,
			},
			expectedErr: statusUnknownAuthType.Err(),
		},
		{
			desc: "err: undefined auth type is rejected",
			req: &authproto.GetAuthenticationURLRequest{
				Type:        authproto.AuthType(99),
				State:       "state",
				RedirectUrl: testRedirectURL,
			},
			expectedErr: statusUnknownAuthType.Err(),
		},
		{
			desc: "err: unspecified auth type is rejected by validation",
			req: &authproto.GetAuthenticationURLRequest{
				Type:        authproto.AuthType_AUTH_TYPE_UNSPECIFIED,
				State:       "state",
				RedirectUrl: testRedirectURL,
			},
			expectedErr: statusMissingAuthType.Err(),
		},
		{
			desc: "success: google auth type returns the login url",
			req: &authproto.GetAuthenticationURLRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				State:       "state",
				RedirectUrl: testRedirectURL,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			resp, err := service.GetAuthenticationURL(context.Background(), p.req)
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			loginURL, err := url.Parse(resp.Url)
			require.NoError(t, err)
			assert.Equal(t, "client-id", loginURL.Query().Get("client_id"))
			assert.Equal(t, "state", loginURL.Query().Get("state"))
			assert.Equal(t, testRedirectURL, loginURL.Query().Get("redirect_uri"))
		})
	}
}

func TestAuthService_ExchangeToken_AuthType(t *testing.T) {
	t.Parallel()
	const email = "user@example.com"
	patterns := []struct {
		desc          string
		req           *authproto.ExchangeTokenRequest
		setup         func(*testing.T, *authService, *accountclientmock.MockClient)
		expectedErr   error
		expectedCalls int
	}{
		{
			desc: "err: github auth type is rejected",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GITHUB,
				Code:        "code",
				RedirectUrl: testRedirectURL,
			},
			expectedErr: statusUnknownAuthType.Err(),
		},
		{
			desc: "err: undefined auth type is rejected",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType(99),
				Code:        "code",
				RedirectUrl: testRedirectURL,
			},
			expectedErr: statusUnknownAuthType.Err(),
		},
		{
			desc: "err: unspecified auth type is rejected by validation",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_UNSPECIFIED,
				Code:        "code",
				RedirectUrl: testRedirectURL,
			},
			expectedErr: statusMissingAuthType.Err(),
		},
		{
			desc: "success: google auth type exchanges the code for a token",
			req: &authproto.ExchangeTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        "code",
				RedirectUrl: testRedirectURL,
			},
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
			expectedCalls: 1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			authenticator := &fakeAuthenticator{userInfo: &auth.UserInfo{Email: email}}
			service := newAuthTypeTestService(t, authenticator)
			accountClient := accountclientmock.NewMockClient(mockController)
			service.accountClient = accountClient
			if p.setup != nil {
				p.setup(t, service, accountClient)
			}
			resp, err := service.ExchangeToken(context.Background(), p.req)
			assert.Equal(t, p.expectedCalls, authenticator.calls)
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
