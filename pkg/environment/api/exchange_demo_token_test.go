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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	acmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/google"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const testDemoRedirectURL = "https://demo.example.com/auth/callback"

// fakeAuthenticator stands in for the Google authenticator so that the demo
// token exchange can be exercised without a network round trip.
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

func TestEnvironmentService_GetAuthenticator(t *testing.T) {
	t.Parallel()
	googleAuthenticator := google.NewAuthenticator(&auth.GoogleConfig{}, zap.NewNop())
	service := &EnvironmentService{
		googleAuthenticator: googleAuthenticator,
		logger:              zap.NewNop(),
	}
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
			expectedErr: auth.StatusUnknownAuthType.Err(),
		},
		{
			desc:        "err: deprecated user password is not supported",
			authType:    authproto.AuthType_AUTH_TYPE_USER_PASSWORD,
			expectedErr: auth.StatusUnknownAuthType.Err(),
		},
		{
			desc:        "err: unspecified is not supported",
			authType:    authproto.AuthType_AUTH_TYPE_UNSPECIFIED,
			expectedErr: auth.StatusUnknownAuthType.Err(),
		},
		{
			desc:        "err: undefined enum value is not supported",
			authType:    authproto.AuthType(99),
			expectedErr: auth.StatusUnknownAuthType.Err(),
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

func TestExchangeDemoToken_AuthType(t *testing.T) {
	t.Parallel()
	const email = "demo@example.com"
	patterns := []struct {
		desc          string
		req           *environmentproto.ExchangeDemoTokenRequest
		setup         func(*testing.T, *EnvironmentService)
		expectedErr   error
		expectedCalls int
	}{
		{
			desc: "err: github auth type is rejected",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GITHUB,
				Code:        "code",
				RedirectUrl: testDemoRedirectURL,
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc: "err: undefined auth type is rejected",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType(99),
				Code:        "code",
				RedirectUrl: testDemoRedirectURL,
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc: "err: unspecified auth type is rejected by validation",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_UNSPECIFIED,
				Code:        "code",
				RedirectUrl: testDemoRedirectURL,
			},
			expectedErr: auth.StatusMissingAuthType.Err(),
		},
		{
			desc: "success: google auth type returns a demo creation token",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        "code",
				RedirectUrl: testDemoRedirectURL,
			},
			setup: func(t *testing.T, s *EnvironmentService) {
				t.Helper()
				signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
				require.NoError(t, err)
				s.signer = signer
				s.accountClient.(*acmock.MockClient).EXPECT().
					GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).
					Return(&accproto.GetMyOrganizationsResponse{}, nil)
			},
			expectedCalls: 1,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			authenticator := &fakeAuthenticator{userInfo: &auth.UserInfo{Email: email}}
			service := newDemoEnvironmentService(t, mockController, nil)
			service.issuer = "test-issuer"
			service.audience = "test-audience"
			service.googleAuthenticator = authenticator
			if p.setup != nil {
				p.setup(t, service)
			}
			resp, err := service.ExchangeDemoToken(context.Background(), p.req)
			assert.Equal(t, p.expectedCalls, authenticator.calls)
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.DemoCreationToken)
			assert.NotEmpty(t, resp.DemoCreationToken.AccessToken)
		})
	}
}
