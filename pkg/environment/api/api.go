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
	"regexp"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	accdomain "github.com/bucketeer-io/bucketeer/pkg/account/domain"
	accstorage "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/auth/google"
	v2 "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	accproto "github.com/bucketeer-io/bucketeer/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type options struct {
	emailFilter       *regexp.Regexp
	isDemoSiteEnabled bool
	logger            *zap.Logger
}

type Option func(*options)

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

func WithDemoSiteEnabled(isDemoSiteEnabled bool) Option {
	return func(opts *options) {
		opts.isDemoSiteEnabled = isDemoSiteEnabled
	}
}

type EnvironmentService struct {
	issuer              string
	audience            string
	signer              token.Signer
	accountClient       accountclient.Client
	mysqlClient         mysql.Client
	projectStorage      v2.ProjectStorage
	orgStorage          v2.OrganizationStorage
	environmentStorage  v2.EnvironmentStorage
	publisher           publisher.Publisher
	googleAuthenticator auth.Authenticator
	verifier            token.Verifier
	opts                *options
	logger              *zap.Logger
}

func NewEnvironmentService(
	ac accountclient.Client,
	mysqlClient mysql.Client,
	publisher publisher.Publisher,
	config *auth.OAuthConfig,
	issuer string,
	audience string,
	signer token.Signer,
	verifier token.Verifier,
	opts ...Option,
) *EnvironmentService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	logger := dopts.logger.Named("api")
	return &EnvironmentService{
		accountClient:      ac,
		mysqlClient:        mysqlClient,
		projectStorage:     v2.NewProjectStorage(mysqlClient),
		orgStorage:         v2.NewOrganizationStorage(mysqlClient),
		environmentStorage: v2.NewEnvironmentStorage(mysqlClient),
		publisher:          publisher,
		googleAuthenticator: google.NewAuthenticator(
			&config.GoogleConfig, logger,
		),
		issuer:   issuer,
		audience: audience,
		signer:   signer,
		verifier: verifier,
		opts:     dopts,
		logger:   logger,
	}
}

func (s *EnvironmentService) Register(server *grpc.Server) {
	environmentproto.RegisterEnvironmentServiceServer(server, s)
}

func (s *EnvironmentService) ExchangeDemoToken(
	ctx context.Context,
	req *environmentproto.ExchangeDemoTokenRequest,
) (*environmentproto.ExchangeDemoTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if !s.opts.isDemoSiteEnabled {
		s.logger.Error("Demo site is not enabled",
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		dt, err := statusDemoSiteDisabled.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.Organization),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	if err := validateExchangeDemoTokenRequest(req, localizer); err != nil {
		s.logger.Error("Failed to validate the exchange demo token request",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, err
	}
	authenticator, err := s.getAuthenticator(req.Type, localizer)
	if err != nil {
		s.logger.Error("Failed to get the authenticator",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	userInfo, err := authenticator.Exchange(ctx, req.Code, req.RedirectUrl)
	if err != nil {
		s.logger.Error("Failed to exchange",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	if !userInfo.VerifiedEmail {
		s.logger.Error("Email is not verified",
			zap.String("email", userInfo.Email),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		dt, err := auth.StatusEmailNotVerified.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.Account),
		})
		if err != nil {
			return nil, auth.StatusEmailNotVerified.Err()
		}
		return nil, dt.Err()
	}

	existedInSystem, err := s.checkEmailExistedInSystem(ctx, userInfo.Email, localizer)
	if err != nil {
		return nil, err
	}
	if existedInSystem {
		s.logger.Error("Email already exists in the system",
			zap.String("email", userInfo.Email),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		dt, err := statusUserAlreadyInOrganization.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.Organization),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}

	demoToken, err := s.generateDemoToken(ctx, userInfo.Email, localizer)
	if err != nil {
		s.logger.Error("Failed to generate demoToken",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, err
	}
	return &environmentproto.ExchangeDemoTokenResponse{
		DemoCreationToken: demoToken,
	}, nil
}

func (s *EnvironmentService) checkEmailExistedInSystem(
	ctx context.Context,
	email string,
	localizer locale.Localizer,
) (bool, error) {
	getAccountOrgs, err := s.accountClient.GetMyOrganizationsByEmail(ctx, &accproto.GetMyOrganizationsByEmailRequest{
		Email: email,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get organizations by email",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", email), zap.Error(err))...,
		)
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return false, auth.StatusInternal.Err()
		}
		return false, dt.Err()
	}
	if len(getAccountOrgs.Organizations) > 0 {
		s.logger.Error(
			"Email already exists in the system",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", email))...)
		return true, nil
	}
	return false, nil
}

func (s *EnvironmentService) generateDemoToken(
	ctx context.Context,
	userEmail string,
	localizer locale.Localizer,
) (*environmentproto.DemoCreationToken, error) {
	if err := s.checkEmail(userEmail, localizer); err != nil {
		s.logger.Error(
			"Access denied email",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", userEmail))...,
		)
		return nil, err
	}

	// Create demo access token
	timeNow := time.Now()
	accessTokenTTL := timeNow.Add(5 * time.Minute)
	accessToken := &token.DemoCreationToken{
		Issuer:   s.issuer,
		Audience: s.audience,
		Expiry:   accessTokenTTL,
		IssuedAt: timeNow,
		Email:    userEmail,
	}

	signedAccessToken, err := s.signer.SignDemoCreationToken(accessToken)
	if err != nil {
		s.logger.Error(
			"Failed to sign access token",
			zap.Error(err),
		)
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}

	return &environmentproto.DemoCreationToken{
		AccessToken: signedAccessToken,
		TokenType:   "Bearer",
		Expiry:      accessTokenTTL.Unix(),
	}, nil
}

func (s *EnvironmentService) checkEmail(
	email string,
	localizer locale.Localizer,
) error {
	if s.opts.emailFilter == nil {
		return nil
	}
	if s.opts.emailFilter.MatchString(email) {
		return nil
	}
	dt, err := auth.StatusAccessDeniedEmail.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.PermissionDenied),
	})
	if err != nil {
		return auth.StatusInternal.Err()
	}
	return dt.Err()
}

func (s *EnvironmentService) getAuthenticator(
	authType authproto.AuthType,
	localizer locale.Localizer,
) (auth.Authenticator, error) {
	var authenticator auth.Authenticator
	switch authType {
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

func validateExchangeDemoTokenRequest(
	req *environmentproto.ExchangeDemoTokenRequest,
	localizer locale.Localizer,
) error {
	if req.Type == authproto.AuthType_AUTH_TYPE_UNSPECIFIED {
		dt, err := auth.StatusMissingAuthType.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "auth_type"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.Code == "" {
		dt, err := auth.StatusMissingCode.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "code"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.RedirectUrl == "" {
		dt, err := auth.StatusMissingRedirectURL.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "redirect_url"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *EnvironmentService) checkSystemAdminRole(
	ctx context.Context,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
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
	}
	return editor, nil
}

func (s *EnvironmentService) checkOrganizationRole(
	ctx context.Context,
	organizationID string,
	requiredRole accproto.AccountV2_Role_Organization,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(
		ctx,
		requiredRole,
		func(email string) (*accproto.GetAccountV2Response, error) {
			return s.accountClient.GetAccountV2(ctx, &accproto.GetAccountV2Request{
				Email:          email,
				OrganizationId: organizationID,
			})
		},
	)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
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
	}
	return editor, nil
}

func (s *EnvironmentService) checkOrganizationRoleByEnvironmentID(
	ctx context.Context,
	requiredRole accproto.AccountV2_Role_Organization,
	environmentID string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(
		ctx,
		requiredRole,
		func(email string,
		) (*accproto.GetAccountV2Response, error) {
			account, err := s.getAccountV2ByEnvironmentID(ctx, email, environmentID, localizer)
			if err != nil {
				return nil, err
			}
			return &accproto.GetAccountV2Response{Account: account.AccountV2}, nil
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
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
	}
	return editor, nil
}

func (s *EnvironmentService) getAccountV2ByEnvironmentID(
	ctx context.Context,
	email, environmentID string,
	localizer locale.Localizer,
) (*accdomain.AccountV2, error) {
	storage := accstorage.NewAccountStorage(s.mysqlClient)
	account, err := storage.GetAccountV2ByEnvironmentID(ctx, email, environmentID)
	if err != nil {
		if errors.Is(err, accstorage.ErrAccountNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get account by environment id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentID),
				zap.String("email", email),
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
	return account, nil
}
