// Copyright 2023 The Bucketeer Authors.
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
	"regexp"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/auth/oidc"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

type options struct {
	refreshTokenTTL time.Duration
	emailFilter     *regexp.Regexp
	logger          *zap.Logger
}

var defaultOptions = options{
	refreshTokenTTL: time.Hour,
	logger:          zap.NewNop(),
}

type Option func(*options)

func WithRefreshTokenTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.refreshTokenTTL = ttl
	}
}

func WithEmailFilter(regexp *regexp.Regexp) Option {
	return func(opts *options) {
		opts.emailFilter = regexp
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type authService struct {
	oidc          *oidc.OIDC
	signer        token.Signer
	accountClient accountclient.Client
	opts          *options
	logger        *zap.Logger
}

func NewAuthService(
	oidc *oidc.OIDC,
	signer token.Signer,
	accountClient accountclient.Client,
	opts ...Option,
) rpc.Service {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &authService{
		oidc:          oidc,
		signer:        signer,
		accountClient: accountClient,
		opts:          &options,
		logger:        options.logger.Named("api"),
	}
}

func (s *authService) Register(server *grpc.Server) {
	authproto.RegisterAuthServiceServer(server, s)
}

func (s *authService) GetAuthCodeURL(
	ctx context.Context,
	req *authproto.GetAuthCodeURLRequest,
) (*authproto.GetAuthCodeURLResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	// The state parameter is used to help mitigate CSRF attacks.
	// Before sending a request to get authCodeURL, the client has to generate a random string,
	// store it in local and set to the state parameter in GetAuthCodeURLRequest.
	// When the client is redirected back, the state value will be included in that redirect.
	// Client compares the returned state to the one generated before,
	// if the values match then send a new request to ExchangeToken, else deny it.
	if err := validateGetAuthCodeURLRequest(req, localizer); err != nil {
		return nil, err
	}
	url, err := s.oidc.AuthCodeURL(req.State, req.RedirectUrl)
	if err != nil {
		if err == oidc.ErrUnregisteredRedirectURL {
			dt, err := statusUnregisteredRedirectURL.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "redirect_url"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get auth code url",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &authproto.GetAuthCodeURLResponse{Url: url}, nil
}

func validateGetAuthCodeURLRequest(req *authproto.GetAuthCodeURLRequest, localizer locale.Localizer) error {
	if req.State == "" {
		dt, err := statusMissingState.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "state"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.RedirectUrl == "" {
		dt, err := statusMissingRedirectURL.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *authService) ExchangeToken(
	ctx context.Context,
	req *authproto.ExchangeTokenRequest,
) (*authproto.ExchangeTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateExchangeTokenRequest(req, localizer); err != nil {
		return nil, err
	}
	authToken, err := s.oidc.Exchange(ctx, req.Code, req.RedirectUrl)
	if err != nil {
		if err == oidc.ErrUnregisteredRedirectURL {
			dt, err := statusUnregisteredRedirectURL.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "redirect_url"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if err == oidc.ErrBadRequest {
			dt, err := statusInvalidCode.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to exchange token",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	token, err := s.generateToken(ctx, authToken, localizer)
	if err != nil {
		return nil, err
	}
	return &authproto.ExchangeTokenResponse{Token: token}, nil
}

func validateExchangeTokenRequest(req *authproto.ExchangeTokenRequest, localizer locale.Localizer) error {
	if req.Code == "" {
		dt, err := statusMissingCode.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "code"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.RedirectUrl == "" {
		dt, err := statusMissingRedirectURL.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *authService) RefreshToken(
	ctx context.Context,
	req *authproto.RefreshTokenRequest,
) (*authproto.RefreshTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateRefreshTokenRequest(req, localizer); err != nil {
		return nil, err
	}
	authToken, err := s.oidc.RefreshToken(ctx, req.RefreshToken, s.opts.refreshTokenTTL, req.RedirectUrl)
	if err != nil {
		if err == oidc.ErrUnregisteredRedirectURL {
			dt, err := statusUnregisteredRedirectURL.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "redirect_url"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if err == oidc.ErrBadRequest {
			dt, err := statusInvalidRefreshToken.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "refresh_token"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to refresh token",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	token, err := s.generateToken(ctx, authToken, localizer)
	if err != nil {
		return nil, err
	}
	return &authproto.RefreshTokenResponse{Token: token}, nil
}

func validateRefreshTokenRequest(req *authproto.RefreshTokenRequest, localizer locale.Localizer) error {
	if req.RefreshToken == "" {
		dt, err := statusMissingRefreshToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "refresh_token"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.RedirectUrl == "" {
		dt, err := statusMissingRedirectURL.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *authService) generateToken(
	ctx context.Context,
	t *oauth2.Token,
	localizer locale.Localizer,
) (*authproto.Token, error) {
	rawIDToken := oidc.ExtractRawIDToken(t)
	if len(rawIDToken) == 0 {
		s.logger.Error(
			"Token does not contain id_token",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Any("oauth2Token", t))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	claims, err := s.oidc.Verify(ctx, rawIDToken)
	if err != nil {
		s.logger.Error(
			"Failed to verify id token",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if err := s.maybeCheckEmail(ctx, claims.Email, localizer); err != nil {
		return nil, err
	}
	resp, err := s.accountClient.GetMeByEmailV2(ctx, &accountproto.GetMeByEmailV2Request{
		Email: claims.Email,
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			s.logger.Warn(
				"Unabled to generate token for an unapproved account",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", claims.Email))...,
			)
			dt, err := statusUnapprovedAccount.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", claims.Email),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	adminRole := accountproto.Account_UNASSIGNED
	if resp.IsAdmin {
		adminRole = accountproto.Account_OWNER
	}
	idToken := &token.IDToken{
		Issuer:    claims.Iss,
		Subject:   claims.Sub,
		Audience:  claims.Aud,
		Expiry:    time.Unix(claims.Exp, 0),
		IssuedAt:  time.Unix(claims.Iat, 0),
		Email:     claims.Email,
		AdminRole: adminRole,
	}
	signedIDToken, err := s.signer.Sign(idToken)
	if err != nil {
		s.logger.Error(
			"Failed to sign id token",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &authproto.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry.Unix(),
		IdToken:      signedIDToken,
	}, nil
}

func (s *authService) maybeCheckEmail(ctx context.Context, email string, localizer locale.Localizer) error {
	if s.opts.emailFilter == nil {
		return nil
	}
	if s.opts.emailFilter.MatchString(email) {
		return nil
	}
	s.logger.Info(
		"Access denied email",
		log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", email))...,
	)
	dt, err := statusAccessDeniedEmail.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.PermissionDenied),
	})
	if err != nil {
		return statusInternal.Err()
	}
	return dt.Err()
}
