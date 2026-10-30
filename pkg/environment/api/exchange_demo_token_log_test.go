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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	acmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

// testDemoAuthorizationCode is a request value that must never be written to the log.
const testDemoAuthorizationCode = "4/0AX4XfWgDemoAuthorizationCode1234567890"

func TestExchangeDemoToken_LogsOmitCode(t *testing.T) {
	t.Parallel()
	const email = "demo@example.com"
	patterns := []struct {
		desc           string
		req            *environmentproto.ExchangeDemoTokenRequest
		authenticator  *fakeAuthenticator
		setup          func(*testing.T, *EnvironmentService)
		expectedErr    error
		expectedLogMsg string
	}{
		{
			desc: "err: disabled demo site is logged without the code",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        testDemoAuthorizationCode,
				RedirectUrl: testDemoRedirectURL,
			},
			authenticator: &fakeAuthenticator{},
			setup: func(t *testing.T, s *EnvironmentService) {
				t.Helper()
				s.opts.isDemoSiteEnabled = false
			},
			expectedErr:    statusDemoSiteDisabled.Err(),
			expectedLogMsg: "Demo site is not enabled",
		},
		{
			desc: "err: validation failure is logged without the code",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type: authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code: testDemoAuthorizationCode,
			},
			authenticator:  &fakeAuthenticator{},
			expectedErr:    auth.StatusMissingRedirectURL.Err(),
			expectedLogMsg: "Failed to validate the exchange demo token request",
		},
		{
			desc: "err: unknown auth type is logged without the code",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GITHUB,
				Code:        testDemoAuthorizationCode,
				RedirectUrl: testDemoRedirectURL,
			},
			authenticator:  &fakeAuthenticator{},
			expectedErr:    statusInternal.Err(),
			expectedLogMsg: "Failed to get the authenticator",
		},
		{
			desc: "err: exchange failure is logged without the code",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        testDemoAuthorizationCode,
				RedirectUrl: testDemoRedirectURL,
			},
			authenticator:  &fakeAuthenticator{err: errors.New("exchange failed")},
			expectedErr:    statusInternal.Err(),
			expectedLogMsg: "Failed to exchange",
		},
		{
			desc: "err: already registered email is logged without the code",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        testDemoAuthorizationCode,
				RedirectUrl: testDemoRedirectURL,
			},
			authenticator: &fakeAuthenticator{userInfo: &auth.UserInfo{Email: email}},
			setup: func(t *testing.T, s *EnvironmentService) {
				t.Helper()
				s.accountClient.(*acmock.MockClient).EXPECT().
					GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).
					Return(&accproto.GetMyOrganizationsResponse{
						Organizations: []*environmentproto.Organization{{Id: "org-1"}},
					}, nil)
			},
			expectedErr:    statusUserAlreadyInOrganization.Err(),
			expectedLogMsg: "Email already exists in the system",
		},
		{
			desc: "success: the code is exchanged for a demo creation token",
			req: &environmentproto.ExchangeDemoTokenRequest{
				Type:        authproto.AuthType_AUTH_TYPE_GOOGLE,
				Code:        testDemoAuthorizationCode,
				RedirectUrl: testDemoRedirectURL,
			},
			authenticator: &fakeAuthenticator{userInfo: &auth.UserInfo{Email: email}},
			setup: func(t *testing.T, s *EnvironmentService) {
				t.Helper()
				signer, err := token.NewSigner("../../token/testdata/valid-private.pem")
				require.NoError(t, err)
				s.signer = signer
				s.accountClient.(*acmock.MockClient).EXPECT().
					GetMyOrganizationsByEmail(gomock.Any(), gomock.Any()).
					Return(&accproto.GetMyOrganizationsResponse{}, nil)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			core, logs := observer.New(zapcore.DebugLevel)
			service := newDemoEnvironmentService(t, mockController, nil)
			service.logger = zap.New(core)
			service.issuer = "test-issuer"
			service.audience = "test-audience"
			service.googleAuthenticator = p.authenticator
			if p.setup != nil {
				p.setup(t, service)
			}
			resp, err := service.ExchangeDemoToken(context.Background(), p.req)
			for _, entry := range logs.All() {
				assert.NotContains(t, entry.Message, testDemoAuthorizationCode)
				for key, field := range entry.ContextMap() {
					assert.NotContains(t, fmt.Sprint(field), testDemoAuthorizationCode,
						"field %q of %q", key, entry.Message)
				}
			}
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
				// The request context fields are still logged alongside the error.
				entries := logs.FilterMessage(p.expectedLogMsg).
					FilterField(zap.String("redirect_url", p.req.RedirectUrl)).All()
				require.Len(t, entries, 1, "expected exactly one %q entry for the request", p.expectedLogMsg)
				assert.Equal(t, p.req.Type.String(), fmt.Sprint(entries[0].ContextMap()["type"]))
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.DemoCreationToken)
			assert.NotEmpty(t, resp.DemoCreationToken.AccessToken)
			assert.Equal(t, 1, p.authenticator.calls)
		})
	}
}
