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

package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)


// RequestMagicLink initiates the magic link authentication flow
func (s *authService) RequestMagicLink(
	ctx context.Context,
	req *authproto.RequestMagicLinkRequest,
) (*authproto.RequestMagicLinkResponse, error) {
	localizer := locale.NewLocalizer(ctx)

	if err := validateRequestMagicLinkRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to validate the request magic link request",
			zap.Error(err),
			zap.String("email", req.Email),
		)
		return nil, err
	}

	email := req.Email
	token, err := auth.GenerateSecureToken()
	if err != nil {
		s.logger.Error(
			"Failed to generate magic link token",
			zap.Error(err),
			zap.String("email", email),
		)
		return nil, auth.StatusInternal.Err()
	}

	expiresAt := time.Now().Add(s.config.MagicLink.Tokens.VerifyTTL).Unix()
	ipAddress := getIPAddress(ctx)
	userAgent := getUserAgent(ctx)

	err = s.emailVerificationStorage.CreateVerificationToken(ctx, email, token, expiresAt, ipAddress, userAgent)
	if err != nil {
		s.logger.Error(
			"Failed to create verification token",
			zap.Error(err),
			zap.String("email", email),
		)
	}

	orgsResp, err := s.accountClient.GetMyOrganizationsByEmail(
		ctx,
		&acproto.GetMyOrganizationsByEmailRequest{Email: email},
	)
	if err == nil && orgsResp != nil && len(orgsResp.Organizations) > 0 {
		if s.emailService != nil {
			verifyPath := s.config.MagicLink.URLs.VerifyPath
			if verifyPath == "" {
				s.logger.Error("Magic link verify path not configured")
				return nil, auth.StatusInternal.Err()
			}
			verifyParam := s.config.MagicLink.URLs.TokenParam
			if verifyParam == "" {
				s.logger.Error("Magic link verify parameter not configured")
				return nil, auth.StatusInternal.Err()
			}
			magicLinkURL := fmt.Sprintf("%s%s?%s=%s",
				s.emailConfig.BaseURL, verifyPath, verifyParam, token)
			sendErr := s.emailService.SendMagicLinkEmail(
				ctx, email, magicLinkURL, s.config.MagicLink.Tokens.VerifyTTL, localizer.GetLocale(),
			)
			if sendErr != nil {
				s.logger.Error(
					"Failed to send magic link email",
					zap.Error(sendErr),
					zap.String("email", email),
				)
			} else {
				s.logger.Info(
					"Magic link email sent",
					zap.String("email", email),
				)
			}
		} else {
			s.logger.Warn(
				"Email service not configured, magic link not sent",
				zap.String("email", email),
			)
		}
	} else {
		s.logger.Info(
			"Magic link requested for non-existent email",
			zap.String("email", email),
		)
	}

	return &authproto.RequestMagicLinkResponse{
		Message: "If this email is registered, we've sent instructions to your inbox.",
	}, nil
}

// VerifyMagicLink verifies a magic link token and returns the verified email and organizations
func (s *authService) VerifyMagicLink(
	ctx context.Context,
	req *authproto.VerifyMagicLinkRequest,
) (*authproto.VerifyMagicLinkResponse, error) {
	localizer := locale.NewLocalizer(ctx)

	if err := validateVerifyMagicLinkRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to validate the verify magic link request",
			zap.Error(err),
		)
		return nil, err
	}

	vToken, err := s.emailVerificationStorage.GetVerificationToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, storage.ErrEmailVerificationTokenNotFound) {
			s.logger.Error(
				"Magic link token not found",
				zap.String("token", req.Token),
			)
			dt, err := auth.StatusInvalidMagicLinkToken.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return nil, err
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get verification token",
			zap.Error(err),
		)
		return nil, auth.StatusInternal.Err()
	}

	if vToken.IsExpired() {
		s.logger.Error(
			"Magic link token expired",
			zap.String("email", vToken.Email),
			zap.Int64("expires_at", vToken.ExpiresAt),
		)
		dt, err := auth.StatusExpiredMagicLinkToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InvalidArgumentError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	if vToken.IsVerified() {
		if vToken.WasRecentlyVerified() {
			s.logger.Info(
				"Magic link token already verified recently, allowing retry",
				zap.String("email", vToken.Email),
			)
		} else {
			s.logger.Error(
				"Magic link token already used",
				zap.String("email", vToken.Email),
				zap.Int64("verified_at", *vToken.VerifiedAt),
			)
			dt, err := auth.StatusMagicLinkTokenAlreadyUsed.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return nil, err
			}
			return nil, dt.Err()
		}
	} else {
		verifiedAt := time.Now().Unix()
		err = s.emailVerificationStorage.MarkVerified(ctx, req.Token, verifiedAt)
		if err != nil {
			s.logger.Error(
				"Failed to mark token as verified",
				zap.Error(err),
			)
			return nil, auth.StatusInternal.Err()
		}
	}

	envOrganizations, err := s.getOrganizationsByEmail(ctx, vToken.Email, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to get organizations for verified email",
			zap.Error(err),
			zap.String("email", vToken.Email),
		)
		return nil, err
	}

	// Convert environment organizations to auth organizations
	organizations := convertEnvOrgsToAuthOrgs(envOrganizations)

	s.logger.Info(
		"Magic link verified successfully",
		zap.String("email", vToken.Email),
		zap.Int("organization_count", len(organizations)),
	)

	return &authproto.VerifyMagicLinkResponse{
		Email:         vToken.Email,
		Organizations: organizations,
	}, nil
}

func getIPAddress(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if vals := md.Get("x-forwarded-for"); len(vals) > 0 {
		return vals[0]
	}
	if vals := md.Get("x-real-ip"); len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func getUserAgent(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if vals := md.Get("user-agent"); len(vals) > 0 {
		return vals[0]
	}
	return ""
}
