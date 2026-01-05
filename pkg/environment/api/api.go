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
	"errors"
	"regexp"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	accdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2acc "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/google"
	v2 "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
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
	accountStorage      v2acc.AccountStorage
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
		accountStorage:     v2acc.NewAccountStorage(mysqlClient),
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
	if !s.opts.isDemoSiteEnabled {
		s.logger.Error("Demo site is not enabled",
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, statusDemoSiteDisabled.Err()
	}
	if err := validateExchangeDemoTokenRequest(req); err != nil {
		s.logger.Error("Failed to validate the exchange demo token request",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, err
	}
	authenticator, err := s.getAuthenticator(req.Type)
	if err != nil {
		s.logger.Error("Failed to get the authenticator",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, statusInternal.Err()
	}
	userInfo, err := authenticator.Exchange(ctx, req.Code, req.RedirectUrl)
	if err != nil {
		s.logger.Error("Failed to exchange",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, statusInternal.Err()
	}
	existedInSystem, err := s.checkEmailExistedInSystem(ctx, userInfo.Email)
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
		return nil, statusUserAlreadyInOrganization.Err()
	}

	demoToken, err := s.generateDemoToken(ctx, userInfo.Email)
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
) (bool, error) {
	getAccountOrgs, err := s.accountClient.GetMyOrganizationsByEmail(ctx, &accproto.GetMyOrganizationsByEmailRequest{
		Email: email,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get organizations by email",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", email), zap.Error(err))...,
		)
		return false, statusInternal.Err()
	}
	if len(getAccountOrgs.Organizations) > 0 {
		s.logger.Error(
			"Email already exists in the system",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", email))...)
		return true, nil
	}
	return false, nil
}

func (s *EnvironmentService) generateDemoToken(
	ctx context.Context,
	userEmail string,
) (*environmentproto.DemoCreationToken, error) {
	if err := s.checkEmail(userEmail); err != nil {
		s.logger.Error(
			"Access denied email",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", userEmail))...,
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
		return nil, statusInternal.Err()
	}

	return &environmentproto.DemoCreationToken{
		AccessToken: signedAccessToken,
		TokenType:   "Bearer",
		Expiry:      accessTokenTTL.Unix(),
	}, nil
}

func (s *EnvironmentService) checkEmail(
	email string,
) error {
	if s.opts.emailFilter == nil {
		return nil
	}
	if s.opts.emailFilter.MatchString(email) {
		return nil
	}
	return statusPermissionDenied.Err()
}

func (s *EnvironmentService) getAuthenticator(
	authType authproto.AuthType,
) (auth.Authenticator, error) {
	var authenticator auth.Authenticator
	switch authType {
	case authproto.AuthType_AUTH_TYPE_GOOGLE:
		authenticator = s.googleAuthenticator
	case authproto.AuthType_AUTH_TYPE_GITHUB:

	default:
		s.logger.Error("Unknown auth type", zap.String("authType", authType.String()))
		return nil, auth.StatusUnknownAuthType.Err()
	}
	return authenticator, nil
}

func validateExchangeDemoTokenRequest(
	req *environmentproto.ExchangeDemoTokenRequest,
) error {
	if req.Type == authproto.AuthType_AUTH_TYPE_UNSPECIFIED {
		return auth.StatusMissingAuthType.Err()
	}
	if req.Code == "" {
		return auth.StatusMissingCode.Err()
	}
	if req.RedirectUrl == "" {
		return auth.StatusMissingRedirectURL.Err()
	}
	return nil
}

func (s *EnvironmentService) checkSystemAdminRole(
	ctx context.Context,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}

func (s *EnvironmentService) checkOrganizationRole(
	ctx context.Context,
	organizationID string,
	requiredRole accproto.AccountV2_Role_Organization,
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
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}

func (s *EnvironmentService) checkOrganizationRoleByEnvironmentID(
	ctx context.Context,
	requiredRole accproto.AccountV2_Role_Organization,
	environmentID string,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(
		ctx,
		requiredRole,
		func(email string,
		) (*accproto.GetAccountV2Response, error) {
			account, err := s.getAccountV2ByEnvironmentID(ctx, email, environmentID)
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
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}

func (s *EnvironmentService) getAccountV2ByEnvironmentID(
	ctx context.Context,
	email, environmentID string,
) (*accdomain.AccountV2, error) {
	storage := v2acc.NewAccountStorage(s.mysqlClient)
	account, err := storage.GetAccountV2ByEnvironmentID(ctx, email, environmentID)
	if err != nil {
		if errors.Is(err, v2acc.ErrAccountNotFound) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to get account by environment id",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentID),
				zap.String("email", email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return account, nil
}
