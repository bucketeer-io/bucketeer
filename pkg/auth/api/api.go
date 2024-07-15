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
	"errors"
	"regexp"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	accountstotage "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/auth/google"
	envdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	envstotage "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	acproto "github.com/bucketeer-io/bucketeer/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

const (
	day       = 24 * time.Hour
	sevenDays = 7 * 24 * time.Hour
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
	issuer              string
	audience            string
	signer              token.Signer
	config              *auth.OAuthConfig
	mysqlClient         mysql.Client
	accountClient       accountclient.Client
	verifier            token.Verifier
	googleAuthenticator auth.Authenticator
	opts                *options
	logger              *zap.Logger
}

func NewAuthService(
	issuer string,
	audience string,
	signer token.Signer,
	verifier token.Verifier,
	mysqlClient mysql.Client,
	accountClient accountclient.Client,
	config *auth.OAuthConfig,
	opts ...Option,
) rpc.Service {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	logger := options.logger.Named("api")
	service := &authService{
		issuer:        issuer,
		audience:      audience,
		signer:        signer,
		config:        config,
		mysqlClient:   mysqlClient,
		accountClient: accountClient,
		verifier:      verifier,
		googleAuthenticator: google.NewAuthenticator(
			&config.GoogleConfig, signer, logger,
		),
		opts:   &options,
		logger: logger,
	}
	service.PrepareDemoUser()
	return service
}

func (s *authService) Register(server *grpc.Server) {
	authproto.RegisterAuthServiceServer(server, s)
}

func (s *authService) GetAuthenticationURL(
	ctx context.Context,
	req *authproto.GetAuthenticationURLRequest,
) (*authproto.GetAuthenticationURLResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	// The state parameter is used to help mitigate CSRF attacks.
	// Before sending a request to get authCodeURL, the client has to generate a random string,
	// store it in local and set to the state parameter in GetAuthenticationURLRequest.
	// When the client is redirected back, the state value will be included in that redirect.
	// Client compares the returned state to the one generated before,
	// if the values match then send a new request to ExchangeToken, else deny it.
	if err := validateGetAuthenticationURLRequest(req, localizer); err != nil {
		return nil, err
	}
	authenticator, err := s.getAuthenticator(req.Type, localizer)
	if err != nil {
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	loginURL, err := authenticator.Login(ctx, req.State, req.RedirectUrl)
	if err != nil {
		s.logger.Error(
			"Failed to get authentication",
			zap.Error(err),
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
	return &authproto.GetAuthenticationURLResponse{Url: loginURL}, nil
}

func (s *authService) ExchangeToken(
	ctx context.Context,
	req *authproto.ExchangeTokenRequest,
) (*authproto.ExchangeTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateExchangeTokenRequest(req, localizer); err != nil {
		return nil, err
	}
	authenticator, err := s.getAuthenticator(req.Type, localizer)
	if err != nil {
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
		dt, err := auth.StatusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	token, err := s.generateToken(ctx, userInfo.Email, localizer)
	if err != nil {
		return nil, err
	}
	return &authproto.ExchangeTokenResponse{Token: token}, nil
}

func (s *authService) RefreshToken(
	ctx context.Context,
	req *authproto.RefreshTokenRequest,
) (*authproto.RefreshTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := validateRefreshTokenRequest(req, localizer); err != nil {
		return nil, err
	}
	refreshToken, err := s.verifier.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		s.logger.Error("Refresh token is invalid", zap.Any("refresh_token", refreshToken))
		dt, err := auth.StatusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.UnauthenticatedError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}
	newToken, err := s.generateToken(ctx, refreshToken.Email, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to generate token",
			zap.Error(err),
			zap.Any("refresh_token", refreshToken),
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
	return &authproto.RefreshTokenResponse{Token: newToken}, nil
}

func (s *authService) getAuthenticator(
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

func (s *authService) generateToken(
	ctx context.Context,
	userEmail string,
	localizer locale.Localizer,
) (*authproto.Token, error) {
	if err := s.checkEmail(userEmail, localizer); err != nil {
		s.logger.Error(
			"Access denied email",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", userEmail))...,
		)
		return nil, err
	}
	orgResp, err := s.accountClient.GetMyOrganizationsByEmail(
		ctx,
		&acproto.GetMyOrganizationsByEmailRequest{
			Email: userEmail,
		},
	)
	if err != nil {
		s.logger.Error(
			"Failed to get account's organizations",
			zap.Error(err),
			zap.String("email", userEmail),
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
	if len(orgResp.Organizations) == 0 {
		s.logger.Error(
			"Unable to generate token for an unapproved account",
			zap.String("email", userEmail),
		)
		dt, err := auth.StatusUnapprovedAccount.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return nil, auth.StatusInternal.Err()
		}
		return nil, dt.Err()
	}
	timeNow := time.Now()
	accessTokenTTL := timeNow.Add(day)
	accessToken := &token.AccessToken{
		Issuer:        s.issuer,
		Audience:      s.audience,
		Expiry:        accessTokenTTL,
		IssuedAt:      timeNow,
		Email:         userEmail,
		IsSystemAdmin: s.hasSystemAdminOrganization(orgResp.Organizations),
	}
	signedAccessToken, err := s.signer.SignAccessToken(accessToken)
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

	refreshToken := &token.RefreshToken{
		Email:    userEmail,
		Expiry:   timeNow.Add(sevenDays),
		IssuedAt: timeNow,
	}
	signRefreshToken, err := s.signer.SignRefreshToken(refreshToken)
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

	return &authproto.Token{
		AccessToken:  signedAccessToken,
		RefreshToken: signRefreshToken,
		TokenType:    "Bearer",
		Expiry:       accessTokenTTL.Unix(),
	}, nil
}

func (s *authService) checkEmail(
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

func (s *authService) hasSystemAdminOrganization(orgs []*envproto.Organization) bool {
	for _, org := range orgs {
		if org.SystemAdmin {
			return true
		}
	}
	return false
}

func (s *authService) PrepareDemoUser() {
	if s.config.DemoSignInConfig.Username == "" ||
		s.config.DemoSignInConfig.Password == "" ||
		s.config.DemoSignInConfig.Organization == "" ||
		s.config.DemoSignInConfig.Project == "" ||
		s.config.DemoSignInConfig.Environment == "" {
		s.logger.Info("Skip preparing demo user, password login config is not completed")
	}
	ctx := context.Background()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error("create mysql tx error", zap.Error(err))
		return
	}
	now := time.Now()
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		organizationStorage := envstotage.NewOrganizationStorage(tx)
		_, err = organizationStorage.GetOrganization(ctx, s.config.DemoSignInConfig.Organization)
		if err != nil && errors.Is(err, envstotage.ErrOrganizationNotFound) {
			err = organizationStorage.CreateOrganization(ctx, &envdomain.Organization{
				Organization: &envproto.Organization{
					Id:          s.config.DemoSignInConfig.Organization,
					Name:        "Demo organization",
					UrlCode:     "demo",
					Description: "This organization is for demo users",
					Disabled:    false,
					Archived:    false,
					Trial:       false,
					CreatedAt:   now.Unix(),
					UpdatedAt:   now.Unix(),
					SystemAdmin: false,
				}})
			if err != nil {
				return err
			}
		}
		projectStorage := envstotage.NewProjectStorage(tx)
		_, err = projectStorage.GetProject(ctx, s.config.DemoSignInConfig.Project)
		if err != nil && errors.Is(err, envstotage.ErrProjectNotFound) {
			err = projectStorage.CreateProject(ctx, &envdomain.Project{
				Project: &envproto.Project{
					Id:             s.config.DemoSignInConfig.Project,
					Description:    "This project is for demo users",
					Disabled:       false,
					Trial:          false,
					CreatorEmail:   s.config.DemoSignInConfig.Email,
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
					Name:           "Demo",
					UrlCode:        "demo",
					OrganizationId: s.config.DemoSignInConfig.Organization,
				}})
			if err != nil {
				return err
			}
		}
		environmentStorage := envstotage.NewEnvironmentStorage(tx)
		_, err = environmentStorage.GetEnvironmentV2(ctx, s.config.DemoSignInConfig.Environment)
		if err != nil && errors.Is(err, envstotage.ErrEnvironmentNotFound) {
			err = environmentStorage.CreateEnvironmentV2(ctx, &envdomain.EnvironmentV2{
				EnvironmentV2: &envproto.EnvironmentV2{
					Id:             s.config.DemoSignInConfig.Environment,
					Name:           "Demo",
					UrlCode:        "demo",
					Description:    "This environment is for demo users",
					ProjectId:      s.config.DemoSignInConfig.Project,
					Archived:       false,
					CreatedAt:      now.Unix(),
					UpdatedAt:      now.Unix(),
					OrganizationId: s.config.DemoSignInConfig.Organization,
					RequireComment: false,
				}})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Error("prepare demo user organization, project and environment error", zap.Error(err))
		return
	}

	accountStorage := accountstotage.NewAccountStorage(s.mysqlClient)
	_, err = accountStorage.GetAccountV2(
		ctx,
		s.config.DemoSignInConfig.Email,
		s.config.DemoSignInConfig.Organization,
	)
	if err != nil && errors.Is(err, accountstotage.ErrAccountNotFound) {
		err = accountStorage.CreateAccountV2(ctx, &domain.AccountV2{
			AccountV2: &acproto.AccountV2{
				OrganizationId:   s.config.DemoSignInConfig.Organization,
				Email:            s.config.DemoSignInConfig.Email,
				Name:             "demo",
				OrganizationRole: acproto.AccountV2_Role_Organization_ADMIN,
				EnvironmentRoles: []*acproto.AccountV2_EnvironmentRole{
					{
						EnvironmentId: s.config.DemoSignInConfig.Environment,
						Role:          acproto.AccountV2_Role_Environment_EDITOR,
					},
				},
			},
		})
		if err != nil {
			s.logger.Error("prepare demo user account error", zap.Error(err))
		}
	}
}
