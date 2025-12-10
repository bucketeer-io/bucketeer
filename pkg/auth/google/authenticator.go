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

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/oidc"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

const (
	// Google's OIDC issuer URL
	googleIssuer = "https://accounts.google.com"
)

var (
	ErrUnregisteredRedirectURL = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.AuthPackageName,
		"unregistered redirectURL",
		"redirectURL",
	)
)

type Authenticator struct {
	config *auth.GoogleConfig
	logger *zap.Logger
}

func NewAuthenticator(
	config *auth.GoogleConfig,
	logger *zap.Logger,
) *Authenticator {
	return &Authenticator{
		config: config,
		logger: logger.Named("google"),
	}
}

func (a Authenticator) Login(
	ctx context.Context,
	state, redirectURL string,
) (string, error) {
	if err := a.validateRedirectURL(redirectURL); err != nil {
		a.logger.Error("auth/google: failed to validate redirect url", zap.Error(err))
		return "", err
	}

	// Create OIDC provider configuration for Google
	oidcConfig := &authproto.CompanyOidcOption{
		Issuer:       googleIssuer,
		ClientId:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		Scopes:       []string{"openid", "email", "profile"},
	}

	// Create OIDC provider
	provider, err := oidc.NewProvider(ctx, oidcConfig, redirectURL, a.logger)
	if err != nil {
		a.logger.Error("auth/google: failed to create OIDC provider", zap.Error(err))
		return "", err
	}

	// Generate auth URL without PKCE (Google doesn't require it for confidential clients)
	// Add prompt=select_account to match previous behavior
	authURL := provider.GenerateAuthURL(state, "", "", "")

	// Append prompt parameter (OIDC provider doesn't support this directly)
	authURL += "&prompt=select_account"

	return authURL, nil
}

func (a Authenticator) Exchange(
	ctx context.Context,
	code, redirectURL string,
) (*auth.UserInfo, error) {
	if err := a.validateRedirectURL(redirectURL); err != nil {
		a.logger.Error("auth/google: failed to validate redirect url", zap.Error(err))
		return nil, err
	}

	// Create OIDC provider configuration for Google
	oidcConfig := &authproto.CompanyOidcOption{
		Issuer:       googleIssuer,
		ClientId:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		Scopes:       []string{"openid", "email", "profile"},
	}

	// Create OIDC provider
	provider, err := oidc.NewProvider(ctx, oidcConfig, redirectURL, a.logger)
	if err != nil {
		a.logger.Error("auth/google: failed to create OIDC provider", zap.Error(err))
		return nil, err
	}

	// Exchange code for token and get user info
	// No PKCE code verifier, no nonce validation (we didn't send nonce in Login)
	userInfo, err := provider.ExchangeToken(ctx, code, "", "")
	if err != nil {
		a.logger.Error("auth/google: failed to exchange token", zap.Error(err))
		return nil, err
	}

	return userInfo, nil
}

func (a Authenticator) validateRedirectURL(url string) error {
	for _, r := range a.config.RedirectURLs {
		if r == url {
			return nil
		}
	}
	return ErrUnregisteredRedirectURL
}
