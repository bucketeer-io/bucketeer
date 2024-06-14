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

package api

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

func (s *authService) GetAuthenticationURL(
	ctx context.Context,
	req *authproto.GetAuthenticationURLRequest,
) (*authproto.GetAuthenticationURLResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	// The state parameter is used to help mitigate CSRF attacks.
	// Before sending a request to get authCodeURL, the client has to generate a random string,
	// store it in local and set to the state parameter in GetAuthCodeURLRequest.
	// When the client is redirected back, the state value will be included in that redirect.
	// Client compares the returned state to the one generated before,
	// if the values match then send a new request to ExchangeToken, else deny it.
	if err := validateGetAuthenticationURLRequest(req, localizer); err != nil {
		return nil, err
	}
	authenticator, err := s.getAuthenticator(ctx, req.Type, localizer)
	if err != nil {
		return nil, err
	}
	loginURL := authenticator.Login(ctx, req.State, localizer)
	return &authproto.GetAuthenticationURLResponse{Url: loginURL}, nil
}

func (s *authService) ExchangeBucketeerToken(
	ctx context.Context,
	req *authproto.ExchangeBucketeerTokenRequest,
) (*authproto.ExchangeBucketeerTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateExchangeBucketeerTokenRequest(req, localizer); err != nil {
		return nil, err
	}
	authenticator, err := s.getAuthenticator(ctx, req.Type, localizer)
	if err != nil {
		return nil, err
	}
	authToken, err := authenticator.Exchange(ctx, req.Code, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to exchange token",
			zap.Error(err),
			zap.String("auth_type", req.Type.String()),
		)
		return nil, err
	}
	return &authproto.ExchangeBucketeerTokenResponse{Token: authToken}, nil
}

func (s *authService) RefreshBucketeerToken(
	ctx context.Context,
	req *authproto.RefreshBucketeerTokenRequest,
) (*authproto.RefreshBucketeerTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateRefreshBucketeerTokenRequest(req, localizer); err != nil {
		return nil, err
	}
	authenticator, err := s.getAuthenticator(ctx, req.Type, localizer)
	if err != nil {
		return nil, err
	}
	newToken, err := authenticator.Refresh(ctx, req.RefreshToken, s.opts.refreshTokenTTL, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to refresh token",
			zap.Error(err),
			zap.String("auth_type", req.Type.String()),
		)
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}
	return &authproto.RefreshBucketeerTokenResponse{Token: newToken}, nil
}

func (s *authService) getAuthenticator(
	ctx context.Context,
	authType authproto.AuthType,
	localizer locale.Localizer,
) (auth.Authenticator, error) {
	var authenticator auth.Authenticator
	switch authType {
	case authproto.AuthType_AUTH_TYPE_USER_PASSWORD:

	case authproto.AuthType_AUTH_TYPE_GOOGLE:
		authenticator = s.googleAuthenticator
	case authproto.AuthType_AUTH_TYPE_GITHUB:

	default:
		s.logger.Error("Unknown auth type", zap.String("authType", authType.String()))
		dt, err := auth.StatusUnknownAuthType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "auth_type"),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}
	return authenticator, nil
}
