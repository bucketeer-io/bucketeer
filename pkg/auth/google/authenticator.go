// Copyright 2025 The Bucketeer Authors.
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
	"encoding/json"
	"errors"
	"io"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

var (
	ErrUnregisteredRedirectURL = errors.New("google: unregistered redirectURL")
)

var (
	defaultScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
)

type googleUserInfo struct {
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type Authenticator struct {
	config *auth.GoogleConfig
	signer token.Signer
	logger *zap.Logger
}

func NewAuthenticator(
	config *auth.GoogleConfig,
	signer token.Signer,
	logger *zap.Logger,
) *Authenticator {
	return &Authenticator{
		config: config,
		signer: signer,
		logger: logger.Named("auth"),
	}
}

func (a Authenticator) Login(
	ctx context.Context,
	state, redirectURL string,
) (string, error) {
	if err := a.validateRedirectURL(redirectURL); err != nil {
		return "", err
	}
	selectAccount := oauth2.SetAuthURLParam("prompt", "select_account")
	return a.oauth2Config(defaultScopes, redirectURL).AuthCodeURL(state, selectAccount), nil
}

func (a Authenticator) Exchange(
	ctx context.Context,
	code, redirectURL string,
) (*auth.UserInfo, error) {
	if err := a.validateRedirectURL(redirectURL); err != nil {
		return nil, err
	}
	oauth2Config := a.oauth2Config(defaultScopes, redirectURL)
	authToken, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		a.logger.Error("auth/google: failed to exchange token", zap.Error(err))
		return nil, err
	}
	userInfo, err := a.getGoogleUserInfo(ctx, authToken, oauth2Config)
	if err != nil {
		a.logger.Error("auth/google: failed to query user info", zap.Error(err))
		return nil, err
	}
	return &auth.UserInfo{
		FirstName: userInfo.GivenName,
		LastName:  userInfo.FamilyName,
		Avatar:    userInfo.Picture,
		Email:     userInfo.Email,
	}, nil
}

func (a Authenticator) getGoogleUserInfo(
	ctx context.Context,
	t *oauth2.Token,
	config *oauth2.Config,
) (googleUserInfo, error) {
	var userInfo googleUserInfo
	client := config.Client(ctx, t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		a.logger.Error("Google: failed to get user info", zap.Error(err))
		return userInfo, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			a.logger.Error("auth/google: failed to close response body", zap.Error(err))
		}
	}(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		a.logger.Error("auth/google: failed to decode user info", zap.Error(err))
		return userInfo, err
	}
	a.logger.Debug("auth/google: user info", zap.Any("user_info", userInfo))
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

func (a Authenticator) oauth2Config(scopes []string, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.config.ClientID,
		ClientSecret: a.config.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  redirectURL,
	}
}
