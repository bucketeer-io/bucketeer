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

package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

var (
	ErrInvalidIssuer             = pkgErr.NewErrorInvalidArgEmpty(pkgErr.AuthPackageName, "invalid issuer", "issuer")
	ErrProviderDiscoveryFailed   = pkgErr.NewErrorInternal(pkgErr.AuthPackageName, "provider discovery failed")
	ErrTokenExchangeFailed       = pkgErr.NewErrorInternal(pkgErr.AuthPackageName, "token exchange failed")
	ErrIDTokenVerificationFailed = pkgErr.NewErrorUnauthenticated(pkgErr.AuthPackageName, "ID token verification failed")
	ErrEmailNotVerified          = pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "email not verified")
	ErrDomainMismatch            = pkgErr.NewErrorPermissionDenied(pkgErr.AuthPackageName, "email domain mismatch")
	ErrNonceValidationFailed     = pkgErr.NewErrorUnauthenticated(pkgErr.AuthPackageName, "nonce validation failed")
)

// Provider represents a company OIDC provider
type Provider struct {
	config          *authproto.CompanyOidcOption
	oidcProvider    *oidc.Provider
	oauth2Config    *oauth2.Config
	idTokenVerifier *oidc.IDTokenVerifier
	logger          *zap.Logger
	httpClient      *http.Client
}

// NewProvider creates a new OIDC provider with discovery
func NewProvider(
	ctx context.Context,
	config *authproto.CompanyOidcOption,
	redirectURL string,
	logger *zap.Logger,
) (*Provider, error) {
	if config == nil || config.Issuer == "" {
		return nil, ErrInvalidIssuer
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	ctx = oidc.ClientContext(ctx, httpClient)

	// Discover OIDC provider configuration
	provider, err := oidc.NewProvider(ctx, config.Issuer)
	if err != nil {
		logger.Error("oidc: failed to discover provider",
			zap.String("issuer", config.Issuer),
			zap.Error(err))
		return nil, ErrProviderDiscoveryFailed
	}

	// Setup OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       config.Scopes,
	}

	// Ensure openid scope is included
	if !contains(oauth2Config.Scopes, "openid") {
		oauth2Config.Scopes = append([]string{"openid"}, oauth2Config.Scopes...)
	}

	// Create ID token verifier
	idTokenVerifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientId,
	})

	return &Provider{
		config:          config,
		oidcProvider:    provider,
		oauth2Config:    oauth2Config,
		idTokenVerifier: idTokenVerifier,
		logger:          logger.Named("oidc"),
		httpClient:      httpClient,
	}, nil
}

// GenerateAuthURL generates the authorization URL for OIDC login
func (p *Provider) GenerateAuthURL(state, nonce, codeChallenge, codeChallengeMethod string) string {
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("nonce", nonce),
	}

	// Add PKCE parameters if provided
	if codeChallenge != "" {
		opts = append(opts,
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", codeChallengeMethod),
		)
	}

	return p.oauth2Config.AuthCodeURL(state, opts...)
}

// ExchangeToken exchanges authorization code for tokens and returns user info
func (p *Provider) ExchangeToken(
	ctx context.Context,
	code, codeVerifier, expectedNonce string,
) (*auth.UserInfo, error) {
	ctx = oidc.ClientContext(ctx, p.httpClient)

	// Exchange code for tokens
	var opts []oauth2.AuthCodeOption
	if codeVerifier != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}

	oauth2Token, err := p.oauth2Config.Exchange(ctx, code, opts...)
	if err != nil {
		p.logger.Error("oidc: failed to exchange token", zap.Error(err))
		return nil, ErrTokenExchangeFailed
	}

	// Extract and verify ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		p.logger.Error("oidc: no id_token in response")
		return nil, ErrIDTokenVerificationFailed
	}

	idToken, err := p.idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		p.logger.Error("oidc: failed to verify ID token", zap.Error(err))
		return nil, ErrIDTokenVerificationFailed
	}

	// Verify nonce
	if expectedNonce != "" {
		if err := verifyNonce(idToken, expectedNonce); err != nil {
			p.logger.Error("oidc: nonce validation failed", zap.Error(err))
			return nil, ErrNonceValidationFailed
		}
	}

	// Extract claims from ID token
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	if err := idToken.Claims(&claims); err != nil {
		p.logger.Error("oidc: failed to parse claims", zap.Error(err))
		return nil, ErrIDTokenVerificationFailed
	}

	// Verify email is present and verified
	if claims.Email == "" {
		p.logger.Error("oidc: email claim missing")
		return nil, ErrIDTokenVerificationFailed
	}

	if !claims.EmailVerified {
		p.logger.Error("oidc: email not verified", zap.String("email", claims.Email))
		return nil, ErrEmailNotVerified
	}

	return &auth.UserInfo{
		Email:         claims.Email,
		VerifiedEmail: claims.EmailVerified,
		Name:          claims.Name,
		FirstName:     claims.GivenName,
		LastName:      claims.FamilyName,
		Avatar:        claims.Picture,
	}, nil
}

// VerifyEmailDomain verifies that the email domain matches the expected domain
func VerifyEmailDomain(email, expectedDomain string) error {
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}

	domain := parts[1]
	if domain != strings.ToLower(expectedDomain) {
		return ErrDomainMismatch
	}

	return nil
}

// verifyNonce verifies the nonce in the ID token
func verifyNonce(token *oidc.IDToken, expectedNonce string) error {
	var claims struct {
		Nonce string `json:"nonce"`
	}

	if err := token.Claims(&claims); err != nil {
		return err
	}

	if claims.Nonce != expectedNonce {
		return fmt.Errorf("nonce mismatch: expected %s, got %s", expectedNonce, claims.Nonce)
	}

	return nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// DiscoverProviderConfig fetches and returns the OIDC provider configuration
func DiscoverProviderConfig(ctx context.Context, issuer string) (map[string]interface{}, error) {
	discoveryURL := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"

	req, err := http.NewRequestWithContext(ctx, "GET", discoveryURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, err
	}

	return config, nil
}
