// Copyright 2024 The Bucketeer Authors.
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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	oidc "github.com/coreos/go-oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var (
	ErrUnregisteredRedirectURL = errors.New("oidc: unregistered redirectURL")
	ErrBadRequest              = errors.New("oidc: bad request")
)

type Claims struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
	Exp           int64  `json:"exp"`
	Iat           int64  `json:"iat"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

type options struct {
	scopes      []string
	httpTimeout time.Duration
	logger      *zap.Logger
}

var defaultOptions = options{
	scopes:      []string{"openid", "profile", "email"},
	httpTimeout: 10 * time.Second,
	logger:      zap.NewNop(),
}

type Option func(*options)

func WithScopes(scopes []string) Option {
	return func(opts *options) {
		opts.scopes = scopes
	}
}

func WithHTTPTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.httpTimeout = timeout
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type OIDC struct {
	clientID       string
	clientSecret   string
	redirectURLs   []string
	provider       *oidc.Provider
	verifier       *oidc.IDTokenVerifier
	offlineAsScope bool
	client         *http.Client
	opts           *options
	logger         *zap.Logger
}

func NewOIDC(
	ctx context.Context,
	issuerURL, issuerCertPath, clientID, clientSecret string,
	redirectURLs []string,
	opts ...Option,
) (*OIDC, error) {
	dopts := defaultOptions
	for _, opt := range opts {
		opt(&dopts)
	}
	cert, err := ioutil.ReadFile(issuerCertPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, errors.New("oidc: Failed to parse issuer cert")
	}
	httpClient := &http.Client{
		Timeout: dopts.httpTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: os.Getenv("BUCKETEER_TEST_ENABLED") == "true",
			},
		},
	}
	logger := dopts.logger.Named("oidc")
	ctx = oidc.ClientContext(ctx, httpClient)
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		logger.Error("Failed to query provider", zap.Error(err), zap.String("issuerURL", issuerURL))
		return nil, err
	}
	offlineScope, err := checkOfflineScope(provider)
	if err != nil {
		logger.Error("Failed to check offline scope", zap.Error(err))
		return nil, err
	}
	return &OIDC{
		clientID:       clientID,
		clientSecret:   clientSecret,
		redirectURLs:   redirectURLs,
		provider:       provider,
		verifier:       provider.Verifier(&oidc.Config{ClientID: clientID}),
		offlineAsScope: offlineScope,
		client:         httpClient,
		opts:           &dopts,
		logger:         logger,
	}, nil
}

func checkOfflineScope(provider *oidc.Provider) (bool, error) {
	var s struct {
		ScopesSupported []string `json:"scopes_supported"`
	}
	if err := provider.Claims(&s); err != nil {
		return false, err
	}
	if len(s.ScopesSupported) == 0 {
		return true, nil
	}
	for _, scope := range s.ScopesSupported {
		if scope == oidc.ScopeOfflineAccess {
			return true, nil
		}
	}
	return false, nil
}

func (o *OIDC) AuthCodeURL(state, redirectURL string) (string, error) {
	if err := o.validateRedirectURL(redirectURL); err != nil {
		return "", err
	}
	scopes := o.opts.scopes
	if o.offlineAsScope {
		scopes = append(scopes, "offline_access")
		return o.oauth2Config(scopes, redirectURL).AuthCodeURL(state), nil
	}
	return o.oauth2Config(scopes, redirectURL).AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (o *OIDC) Exchange(ctx context.Context, code, redirectURL string) (*oauth2.Token, error) {
	if err := o.validateRedirectURL(redirectURL); err != nil {
		return nil, err
	}
	ctx = oidc.ClientContext(ctx, o.client)
	token, err := o.oauth2Config(nil, redirectURL).Exchange(ctx, code)
	if err == nil {
		return token, nil
	}
	if isBadRequestError(err) {
		o.logger.Info("failed to exchange token", zap.Error(err))
		return nil, ErrBadRequest
	}
	return nil, err
}

func (o *OIDC) RefreshToken(
	ctx context.Context,
	token string,
	expires time.Duration,
	redirectURL string,
) (*oauth2.Token, error) {
	if err := o.validateRedirectURL(redirectURL); err != nil {
		return nil, err
	}
	t := &oauth2.Token{
		RefreshToken: token,
		Expiry:       time.Now().Add(expires),
	}
	ctx = oidc.ClientContext(ctx, o.client)
	newToken, err := o.oauth2Config(nil, redirectURL).TokenSource(ctx, t).Token()
	if err == nil {
		return newToken, nil
	}
	if isBadRequestError(err) {
		o.logger.Info("failed to refresh token", zap.Error(err))
		return nil, ErrBadRequest
	}
	return nil, err
}

func (o *OIDC) Verify(ctx context.Context, rawIDToken string) (*Claims, error) {
	idToken, err := o.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}
	claims := &Claims{}
	if err := idToken.Claims(claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func ExtractRawIDToken(token *oauth2.Token) string {
	rawIDToken, _ := token.Extra("id_token").(string)
	return rawIDToken
}

func (o *OIDC) validateRedirectURL(url string) error {
	for _, r := range o.redirectURLs {
		if r == url {
			return nil
		}
	}
	return ErrUnregisteredRedirectURL
}

func (o *OIDC) oauth2Config(scopes []string, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.clientID,
		ClientSecret: o.clientSecret,
		Endpoint:     o.provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  redirectURL,
	}
}

func isBadRequestError(err error) bool {
	if retrieveErr, ok := err.(*oauth2.RetrieveError); ok {
		if code := retrieveErr.Response.StatusCode; code > 200 && code < 500 {
			return true
		}
	}
	return false
}
