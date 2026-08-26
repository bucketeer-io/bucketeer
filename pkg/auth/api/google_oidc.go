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

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func (s *authService) GetGoogleOidcAuthURL(
	ctx context.Context,
	request *authproto.GetGoogleOidcAuthURLRequest,
) (*authproto.GetGoogleOidcAuthURLResponse, error) {
	err := validateGetGoogleOidcAuthURLRequest(request)
	if err != nil {
		s.logger.Error("GetGoogleOidcAuthURL request validation failed", zap.Error(err))
		return nil, err
	}

	// Use existing Google authenticator to generate auth URL
	url, err := s.googleAuthenticator.Login(ctx, request.State, request.RedirectUrl)
	if err != nil {
		s.logger.Error("Failed to generate Google OIDC auth URL",
			zap.Error(err),
			zap.String("redirect_url", request.RedirectUrl),
		)
		return nil, auth.StatusInternal.Err()
	}

	s.logger.Info("Google OIDC auth URL generated successfully")
	return &authproto.GetGoogleOidcAuthURLResponse{
		Url: url,
	}, nil
}

func (s *authService) ExchangeGoogleOidcToken(
	ctx context.Context,
	request *authproto.ExchangeGoogleOidcTokenRequest,
) (*authproto.ExchangeGoogleOidcTokenResponse, error) {
	err := validateExchangeGoogleOidcTokenRequest(request)
	if err != nil {
		s.logger.Error("ExchangeGoogleOidcToken request validation failed", zap.Error(err))
		return nil, err
	}

	// Exchange code for user info using existing Google authenticator
	userInfo, err := s.googleAuthenticator.Exchange(ctx, request.Code, request.RedirectUrl)
	if err != nil {
		s.logger.Error("Failed to exchange Google OIDC code",
			zap.Error(err),
			zap.String("code", request.Code),
		)
		return nil, statusAccessDenied.Err()
	}

	// Create temporary token (5 minutes, no org scope)
	// User must call GetMyOrganizations then SwitchOrganization to get org-scoped token
	token, err := s.createTemporaryToken(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Google OIDC token exchange successful - temporary token issued",
		zap.String("email", userInfo.Email),
	)

	return &authproto.ExchangeGoogleOidcTokenResponse{
		Token: token,
	}, nil
}
