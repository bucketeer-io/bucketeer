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

package google

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
)

const (
	testRedirectURL         = "https://console.example.com/auth/callback"
	testAccessToken         = "access-token"
	testUserInfoURL         = "https://www.googleapis.com/oauth2/v2/userinfo"
	testTokenResponse       = `{"access_token":"` + testAccessToken + `","token_type":"Bearer","expires_in":3600}`
	testAuthorizationHeader = "Bearer " + testAccessToken
)

// roundTripFunc serves Google's token and user info endpoints in-process so
// that Exchange can be exercised without a network round trip.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newGoogleStubClient(t *testing.T, userInfoResponse string) *http.Client {
	t.Helper()
	return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		var body string
		switch {
		case req.Method == http.MethodPost && req.URL.String() == google.Endpoint.TokenURL:
			body = testTokenResponse
		case req.Method == http.MethodGet && req.URL.String() == testUserInfoURL:
			if req.Header.Get("Authorization") != testAuthorizationHeader {
				return newJSONResponse(req, http.StatusUnauthorized, `{"error":"invalid_token"}`), nil
			}
			body = userInfoResponse
		default:
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
			return newJSONResponse(req, http.StatusNotFound, `{}`), nil
		}
		return newJSONResponse(req, http.StatusOK, body), nil
	})}
}

func newJSONResponse(req *http.Request, statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func TestAuthenticator_Exchange(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		allowedDomains []string
		userInfo       string
		expected       *auth.UserInfo
		expectedErr    error
	}{
		{
			desc: "success: verified email when no domain restriction is configured",
			userInfo: `{"name":"Test User","given_name":"Test","family_name":"User",` +
				`"picture":"https://example.com/avatar.png","email":"user@example.com","verified_email":true}`,
			expected: &auth.UserInfo{
				Name:          "Test User",
				FirstName:     "Test",
				LastName:      "User",
				Avatar:        "https://example.com/avatar.png",
				Email:         "user@example.com",
				VerifiedEmail: true,
			},
		},
		{
			desc:     "success: verified email from a hosted domain when no domain restriction is configured",
			userInfo: `{"email":"user@example.com","verified_email":true,"hd":"example.com"}`,
			expected: &auth.UserInfo{Email: "user@example.com", VerifiedEmail: true},
		},
		{
			desc:           "success: verified email from the allowed hosted domain",
			allowedDomains: []string{"example.com"},
			userInfo:       `{"email":"user@example.com","verified_email":true,"hd":"example.com"}`,
			expected:       &auth.UserInfo{Email: "user@example.com", VerifiedEmail: true},
		},
		{
			desc:           "success: verified email from one of several allowed hosted domains",
			allowedDomains: []string{"example.com", "example.org"},
			userInfo:       `{"email":"user@example.org","verified_email":true,"hd":"example.org"}`,
			expected:       &auth.UserInfo{Email: "user@example.org", VerifiedEmail: true},
		},
		{
			desc:           "success: hosted domain comparison ignores case",
			allowedDomains: []string{"Example.com"},
			userInfo:       `{"email":"user@example.com","verified_email":true,"hd":"example.com"}`,
			expected:       &auth.UserInfo{Email: "user@example.com", VerifiedEmail: true},
		},
		{
			desc:        "err: unverified email is rejected",
			userInfo:    `{"email":"user@example.com","verified_email":false}`,
			expectedErr: auth.StatusAccessDeniedEmail.Err(),
		},
		{
			desc:        "err: missing verified_email is rejected",
			userInfo:    `{"email":"user@example.com"}`,
			expectedErr: auth.StatusAccessDeniedEmail.Err(),
		},
		{
			desc:           "err: unverified email is rejected even from the allowed hosted domain",
			allowedDomains: []string{"example.com"},
			userInfo:       `{"email":"user@example.com","verified_email":false,"hd":"example.com"}`,
			expectedErr:    auth.StatusAccessDeniedEmail.Err(),
		},
		{
			desc:           "err: hosted domain outside the allowed list is rejected",
			allowedDomains: []string{"example.com"},
			userInfo:       `{"email":"user@other.example","verified_email":true,"hd":"other.example"}`,
			expectedErr:    auth.StatusAccessDeniedEmail.Err(),
		},
		{
			desc:           "err: account without a hosted domain is rejected even when the email domain matches",
			allowedDomains: []string{"example.com"},
			userInfo:       `{"email":"user@example.com","verified_email":true}`,
			expectedErr:    auth.StatusAccessDeniedEmail.Err(),
		},
		{
			desc:           "err: empty allowed domain entry does not admit an account without a hosted domain",
			allowedDomains: []string{""},
			userInfo:       `{"email":"user@example.com","verified_email":true}`,
			expectedErr:    auth.StatusAccessDeniedEmail.Err(),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			authenticator := NewAuthenticator(&auth.GoogleConfig{
				ClientID:       "client-id",
				ClientSecret:   "client-secret",
				RedirectURLs:   []string{testRedirectURL},
				AllowedDomains: p.allowedDomains,
			}, zap.NewNop())
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, newGoogleStubClient(t, p.userInfo))
			actual, err := authenticator.Exchange(ctx, "code", testRedirectURL)
			if p.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, p.expectedErr.Error(), err.Error())
				assert.Nil(t, actual)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}
