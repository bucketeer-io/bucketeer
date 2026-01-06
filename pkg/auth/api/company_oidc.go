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
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	authdomain "github.com/bucketeer-io/bucketeer/v2/pkg/auth/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/oidc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func (s *authService) GetCompanyOidcAuthURL(
	ctx context.Context,
	request *authproto.GetCompanyOidcAuthURLRequest,
) (*authproto.GetCompanyOidcAuthURLResponse, error) {
	localizer := locale.NewLocalizer(ctx)

	// Validate request
	if err := validateGetCompanyOidcAuthURLRequest(request, localizer); err != nil {
		return nil, err
	}

	// Normalize email and extract domain
	normalizedEmail, err := authdomain.NormalizeEmail(request.Email)
	if err != nil {
		s.logger.Error("oidc: failed to normalize email",
			zap.String("email", request.Email),
			zap.Error(err))
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	domain, err := authdomain.ExtractDomain(normalizedEmail)
	if err != nil {
		s.logger.Error("oidc: failed to extract domain",
			zap.String("email", normalizedEmail),
			zap.Error(err))
		return nil, auth.StatusInvalidArguments.Err()
	}

	// Get domain policy
	policy, err := s.domainPolicyStorage.GetDomainPolicy(ctx, domain)
	if err != nil {
		s.logger.Error("oidc: failed to get domain policy",
			zap.String("domain", domain),
			zap.Error(err))
		return nil, auth.StatusNotFound.Err()
	}

	// Check if company OIDC is enabled
	if policy.AuthPolicy == nil || policy.AuthPolicy.CompanyOidc == nil || !policy.AuthPolicy.CompanyOidc.Enabled {
		s.logger.Error("oidc: company OIDC not enabled for domain",
			zap.String("domain", domain))
		return nil, auth.StatusNotFound.Err()
	}

	// Generate nonce
	nonce, err := oidc.GenerateNonce()
	if err != nil {
		s.logger.Error("oidc: failed to generate nonce", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Create OIDC provider
	provider, err := oidc.NewProvider(ctx, policy.AuthPolicy.CompanyOidc, request.RedirectUrl, s.logger)
	if err != nil {
		s.logger.Error("oidc: failed to create provider",
			zap.String("issuer", policy.AuthPolicy.CompanyOidc.Issuer),
			zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Generate authorization URL
	authURL := provider.GenerateAuthURL(
		request.State,
		nonce,
		request.CodeChallenge,
		request.CodeChallengeMethod,
	)

	return &authproto.GetCompanyOidcAuthURLResponse{
		Url:   authURL,
		Nonce: nonce,
	}, nil
}

func (s *authService) ExchangeCompanyOidcToken(
	ctx context.Context,
	request *authproto.ExchangeCompanyOidcTokenRequest,
) (*authproto.ExchangeCompanyOidcTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)

	// Validate request
	if err := validateExchangeCompanyOidcTokenRequest(request, localizer); err != nil {
		return nil, err
	}

	// Normalize email and extract domain
	normalizedEmail, err := authdomain.NormalizeEmail(request.Email)
	if err != nil {
		s.logger.Error("oidc: failed to normalize email",
			zap.String("email", request.Email),
			zap.Error(err))
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	domain, err := authdomain.ExtractDomain(normalizedEmail)
	if err != nil {
		s.logger.Error("oidc: failed to extract domain",
			zap.String("email", normalizedEmail),
			zap.Error(err))
		return nil, auth.StatusInvalidArguments.Err()
	}

	// Get domain policy
	policy, err := s.domainPolicyStorage.GetDomainPolicy(ctx, domain)
	if err != nil {
		s.logger.Error("oidc: failed to get domain policy",
			zap.String("domain", domain),
			zap.Error(err))
		return nil, auth.StatusNotFound.Err()
	}

	// Check if company OIDC is enabled
	if policy.AuthPolicy == nil || policy.AuthPolicy.CompanyOidc == nil || !policy.AuthPolicy.CompanyOidc.Enabled {
		s.logger.Error("oidc: company OIDC not enabled for domain",
			zap.String("domain", domain))
		return nil, auth.StatusNotFound.Err()
	}

	// Create OIDC provider
	provider, err := oidc.NewProvider(ctx, policy.AuthPolicy.CompanyOidc, request.RedirectUrl, s.logger)
	if err != nil {
		s.logger.Error("oidc: failed to create provider",
			zap.String("issuer", policy.AuthPolicy.CompanyOidc.Issuer),
			zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Exchange code for token and get user info
	userInfo, err := provider.ExchangeToken(ctx, request.Code, request.CodeVerifier, request.Nonce)
	if err != nil {
		s.logger.Error("oidc: token exchange failed", zap.Error(err))
		return nil, auth.StatusUnauthenticated.Err()
	}

	// Verify email domain matches
	if err := oidc.VerifyEmailDomain(userInfo.Email, domain); err != nil {
		s.logger.Error("oidc: email domain mismatch",
			zap.String("email", userInfo.Email),
			zap.String("expected_domain", domain),
			zap.Error(err))
		return nil, auth.StatusAccessDenied.Err()
	}

	// Create temporary token (not org-scoped yet)
	token, err := s.createTemporaryToken(ctx, userInfo)
	if err != nil {
		s.logger.Error("oidc: failed to create token", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	return &authproto.ExchangeCompanyOidcTokenResponse{
		Token: token,
	}, nil
}

// createTemporaryToken creates a temporary, not-yet-org-scoped token
func (s *authService) createTemporaryToken(ctx context.Context, userInfo *auth.UserInfo) (*authproto.Token, error) {
	// Check email filter
	if err := s.checkEmail(userInfo.Email); err != nil {
		s.logger.Error(
			"oidc: access denied email",
			zap.String("email", userInfo.Email),
		)
		return nil, err
	}

	// Create temporary access token without organization ID
	// This token can ONLY be used to:
	// 1. Call GetMyOrganizations to see available orgs
	// 2. Call SwitchOrganization to get an org-scoped token
	timeNow := time.Now()
	temporaryTokenTTL := 5 * time.Minute

	accessToken := &token.AccessToken{
		Issuer:         s.issuer,
		Audience:       s.audience,
		Expiry:         timeNow.Add(temporaryTokenTTL),
		IssuedAt:       timeNow,
		Email:          userInfo.Email,
		Name:           userInfo.Name,
		OrganizationID: "", // Empty - not org-scoped
		IsSystemAdmin:  false,
	}

	signedAccessToken, err := s.signer.SignAccessToken(accessToken)
	if err != nil {
		s.logger.Error(
			"oidc: failed to sign temporary access token",
			zap.Error(err),
			zap.String("email", userInfo.Email),
		)
		return nil, auth.StatusInternal.Err()
	}

	// No refresh token for temporary tokens - they must complete org selection quickly
	return &authproto.Token{
		AccessToken:  signedAccessToken,
		RefreshToken: "", // No refresh token for temporary auth
		TokenType:    "Bearer",
		Expiry:       timeNow.Add(temporaryTokenTTL).Unix(),
	}, nil
}

func validateGetCompanyOidcAuthURLRequest(
	req *authproto.GetCompanyOidcAuthURLRequest,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	if req.State == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "state"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	if req.RedirectUrl == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	// PKCE validation: If code_challenge is provided, code_challenge_method must be provided too
	if req.CodeChallenge != "" && req.CodeChallengeMethod == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "code_challenge_method"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	return nil
}

func validateExchangeCompanyOidcTokenRequest(
	req *authproto.ExchangeCompanyOidcTokenRequest,
	localizer locale.Localizer,
) error {
	if req.Code == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "code"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	if req.State == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "state"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	if req.RedirectUrl == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	if req.Email == "" {
		dt, err := auth.StatusInvalidArguments.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return err
		}
		return dt.Err()
	}

	return nil
}
