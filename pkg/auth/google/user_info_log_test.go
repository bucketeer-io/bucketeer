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
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/oauth2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
)

// userInfoStubTransport answers every request with the configured user info
// body so that the user info lookup can be exercised without a network round
// trip.
type userInfoStubTransport struct {
	body string
}

func (s userInfoStubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(s.body)),
		Request:    req,
	}, nil
}

func TestGetGoogleUserInfo_LogsOmitResponseBody(t *testing.T) {
	t.Parallel()
	// Profile values that must never be written to the log.
	const (
		email     = "user@example.com"
		givenName = "Given"
		picture   = "https://lh3.googleusercontent.com/a/profile-picture"
	)
	body := fmt.Sprintf(
		`{"name":"%s Family","given_name":"%s","family_name":"Family","picture":"%s","email":"%s","verified_email":true}`,
		givenName, givenName, picture, email,
	)
	core, logs := observer.New(zapcore.DebugLevel)
	authenticator := NewAuthenticator(&auth.GoogleConfig{ClientID: "client-id"}, zap.New(core))
	ctx := context.WithValue(
		context.Background(),
		oauth2.HTTPClient,
		&http.Client{Transport: userInfoStubTransport{body: body}},
	)
	token := &oauth2.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}

	userInfo, err := authenticator.getGoogleUserInfo(ctx, token, authenticator.oauth2Config(defaultScopes, ""))

	require.NoError(t, err)
	assert.Equal(t, googleUserInfo{
		Name:          givenName + " Family",
		GivenName:     givenName,
		FamilyName:    "Family",
		Picture:       picture,
		Email:         email,
		EmailVerified: true,
	}, userInfo)
	for _, entry := range logs.All() {
		for _, value := range []string{email, givenName, picture} {
			assert.NotContains(t, entry.Message, value)
			for key, field := range entry.ContextMap() {
				assert.NotContains(t, fmt.Sprint(field), value, "field %q of %q", key, entry.Message)
			}
		}
	}
}
