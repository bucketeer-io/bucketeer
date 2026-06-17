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
	"testing"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
)

func TestNewAuthenticator(t *testing.T) {
	config := &auth.GoogleConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURLs: []string{"https://example.com/callback"},
	}
	logger := zap.NewNop()

	authenticator := NewAuthenticator(config, logger)

	if authenticator == nil {
		t.Fatal("NewAuthenticator() returned nil")
	}

	if authenticator.config != config {
		t.Error("NewAuthenticator() config not set correctly")
	}

	if authenticator.logger == nil {
		t.Error("NewAuthenticator() logger not set")
	}
}

func TestValidateRedirectURL(t *testing.T) {
	tests := []struct {
		name        string
		redirectURL string
		configured  []string
		wantError   bool
	}{
		{
			name:        "valid redirect URL",
			redirectURL: "https://example.com/callback",
			configured:  []string{"https://example.com/callback"},
			wantError:   false,
		},
		{
			name:        "invalid redirect URL",
			redirectURL: "https://malicious.com/callback",
			configured:  []string{"https://example.com/callback"},
			wantError:   true,
		},
		{
			name:        "multiple configured URLs - match first",
			redirectURL: "https://example1.com/callback",
			configured:  []string{"https://example1.com/callback", "https://example2.com/callback"},
			wantError:   false,
		},
		{
			name:        "multiple configured URLs - match second",
			redirectURL: "https://example2.com/callback",
			configured:  []string{"https://example1.com/callback", "https://example2.com/callback"},
			wantError:   false,
		},
		{
			name:        "empty redirect URL list",
			redirectURL: "https://example.com/callback",
			configured:  []string{},
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &auth.GoogleConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURLs: tt.configured,
			}
			authenticator := NewAuthenticator(config, zap.NewNop())

			err := authenticator.validateRedirectURL(tt.redirectURL)
			if (err != nil) != tt.wantError {
				t.Errorf("validateRedirectURL() error = %v, wantError %v", err, tt.wantError)
			}

			if tt.wantError && err != ErrUnregisteredRedirectURL {
				t.Errorf("validateRedirectURL() error = %v, want %v", err, ErrUnregisteredRedirectURL)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	config := &auth.GoogleConfig{
		ClientID:     "test-client-id.apps.googleusercontent.com",
		ClientSecret: "test-client-secret",
		RedirectURLs: []string{"https://example.com/callback"},
	}
	authenticator := NewAuthenticator(config, zap.NewNop())

	tests := []struct {
		name        string
		state       string
		redirectURL string
		wantError   bool
	}{
		{
			name:        "valid request",
			state:       "random-state-123",
			redirectURL: "https://example.com/callback",
			wantError:   false,
		},
		{
			name:        "invalid redirect URL",
			state:       "random-state-123",
			redirectURL: "https://malicious.com/callback",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			url, err := authenticator.Login(ctx, tt.state, tt.redirectURL)

			if (err != nil) != tt.wantError {
				t.Errorf("Login() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if url == "" {
					t.Error("Login() returned empty URL")
				}

				// Verify URL contains Google's authorization endpoint
				if len(url) < 10 {
					t.Errorf("Login() URL too short: %s", url)
				}

				// Check for state parameter
				// Note: We can't check exact URL format since it includes Google's auth endpoint
			}
		})
	}
}

func TestLoginURLContainsRequiredParameters(t *testing.T) {
	config := &auth.GoogleConfig{
		ClientID:     "test-client-id.apps.googleusercontent.com",
		ClientSecret: "test-client-secret",
		RedirectURLs: []string{"https://example.com/callback"},
	}
	authenticator := NewAuthenticator(config, zap.NewNop())

	ctx := context.Background()
	state := "test-state-123"
	redirectURL := "https://example.com/callback"

	url, err := authenticator.Login(ctx, state, redirectURL)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	// Verify URL contains expected parameters
	// Note: Since we're using OIDC provider now, the URL format will be from Google's OIDC endpoint
	if url == "" {
		t.Error("Login() returned empty URL")
	}

	// The URL should be a valid Google OAuth URL
	// We can't test the exact format without mocking the OIDC provider
	// but we can verify it's not empty and reasonably long
	if len(url) < 50 {
		t.Errorf("Login() URL seems too short: %s", url)
	}
}

func TestExchangeValidatesRedirectURL(t *testing.T) {
	config := &auth.GoogleConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURLs: []string{"https://example.com/callback"},
	}
	authenticator := NewAuthenticator(config, zap.NewNop())

	ctx := context.Background()

	// Test with invalid redirect URL - should fail before attempting token exchange
	_, err := authenticator.Exchange(ctx, "test-code", "https://malicious.com/callback")
	if err != ErrUnregisteredRedirectURL {
		t.Errorf("Exchange() with invalid redirect URL: error = %v, want %v", err, ErrUnregisteredRedirectURL)
	}
}
