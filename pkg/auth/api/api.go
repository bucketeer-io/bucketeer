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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/email"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/storage"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	accountstotage "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/google"
	envdomain "github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	envstotage "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

type options struct {
	accessTokenTTL    time.Duration
	refreshTokenTTL   time.Duration
	emailFilter       *regexp.Regexp
	logger            *zap.Logger
	isDemoSiteEnabled bool
}

var defaultOptions = options{
	accessTokenTTL:    10 * time.Minute,
	refreshTokenTTL:   7 * 24 * time.Hour,
	logger:            zap.NewNop(),
	isDemoSiteEnabled: false,
}

type Option func(*options)

func WithAccessTokenTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.accessTokenTTL = ttl
	}
}

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

func WithDemoSiteEnabled(isDemoSiteEnabled bool) Option {
	return func(opts *options) {
		opts.isDemoSiteEnabled = isDemoSiteEnabled
	}
}

type authService struct {
	issuer              string
	audience            string
	signer              token.Signer
	config              *auth.OAuthConfig
	mysqlClient         mysql.Client
	organizationStorage envstotage.OrganizationStorage
	projectStorage      envstotage.ProjectStorage
	environmentStorage  envstotage.EnvironmentStorage
	accountClient       accountclient.Client
	verifier            token.Verifier
	googleAuthenticator auth.Authenticator
	credentialsStorage  storage.CredentialsStorage
	emailService        email.EmailService
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

	// Initialize email service if password auth and email are enabled
	var emailService email.EmailService
	if config.PasswordAuth.Enabled && config.PasswordAuth.EmailServiceEnabled {
		var err error
		emailService, err = email.NewEmailService(config.PasswordAuth.EmailServiceConfig, logger)
		if err != nil {
			logger.Warn("Failed to initialize email service", zap.Error(err))
			emailService = email.NewNoOpEmailService(logger)
		}
	} else {
		emailService = email.NewNoOpEmailService(logger)
	}

	service := &authService{
		issuer:              issuer,
		audience:            audience,
		signer:              signer,
		config:              config,
		mysqlClient:         mysqlClient,
		organizationStorage: envstotage.NewOrganizationStorage(mysqlClient),
		environmentStorage:  envstotage.NewEnvironmentStorage(mysqlClient),
		projectStorage:      envstotage.NewProjectStorage(mysqlClient),
		accountClient:       accountClient,
		verifier:            verifier,
		googleAuthenticator: google.NewAuthenticator(
			&config.GoogleConfig, logger,
		),
		credentialsStorage: storage.NewCredentialsStorage(mysqlClient),
		emailService:       emailService,
		opts:               &options,
		logger:             logger,
	}
	service.PrepareDemoUser()
	return service
}

func (s *authService) Register(server *grpc.Server) {
	authproto.RegisterAuthServiceServer(server, s)
}

func (s *authService) GetDemoSiteStatus(
	_ context.Context,
	_ *authproto.GetDemoSiteStatusRequest,
) (*authproto.GetDemoSiteStatusResponse, error) {
	return &authproto.GetDemoSiteStatusResponse{
		IsDemoSiteEnabled: s.opts.isDemoSiteEnabled,
	}, nil
}

func (s *authService) GetAuthenticationURL(
	ctx context.Context,
	req *authproto.GetAuthenticationURLRequest,
) (*authproto.GetAuthenticationURLResponse, error) {
	// The state parameter is used to help mitigate CSRF attacks.
	// Before sending a request to get authCodeURL, the client has to generate a random string,
	// store it in local and set to the state parameter in GetAuthenticationURLRequest.
	// When the client is redirected back, the state value will be included in that redirect.
	// Client compares the returned state to the one generated before,
	// if the values match then send a new request to ExchangeToken, else deny it.
	if err := validateGetAuthenticationURLRequest(req); err != nil {
		s.logger.Error("Failed to validate the get authentication url request",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("state", req.State),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, err
	}
	authenticator, err := s.getAuthenticator(req.Type)
	if err != nil {
		s.logger.Error("Failed to get the authenticator",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("state", req.State),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	loginURL, err := authenticator.Login(ctx, req.State, req.RedirectUrl)
	if err != nil {
		s.logger.Error("Failed to get the login url",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("state", req.State),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &authproto.GetAuthenticationURLResponse{Url: loginURL}, nil
}

func (s *authService) ExchangeToken(
	ctx context.Context,
	req *authproto.ExchangeTokenRequest,
) (*authproto.ExchangeTokenResponse, error) {
	if err := validateExchangeTokenRequest(req); err != nil {
		s.logger.Error("Failed to validate the exchange token request",
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
		return nil, api.NewGRPCStatus(err).Err()
	}
	userInfo, err := authenticator.Exchange(ctx, req.Code, req.RedirectUrl)
	if err != nil {
		s.logger.Error("Failed to exchange",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	organizations, err := s.getOrganizationsByEmail(ctx, userInfo.Email)
	if err != nil {
		s.logger.Error("Failed to get organizations by email",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, err
	}

	s.updateUserInfoForOrganizations(ctx, userInfo, organizations)

	// Check if the user has at least one account enabled in any Organization
	account, err := s.checkAccountStatus(ctx, userInfo.Email, organizations)
	if err != nil {
		s.logger.Error("Failed to check account",
			zap.Error(err),
			zap.String("email", userInfo.Email),
			zap.Any("organizations", organizations),
		)
		return nil, err
	}
	accountDomain := domain.AccountV2{AccountV2: account.Account}
	isSystemAdmin := s.hasSystemAdminOrganization(organizations)

	token, err := s.generateToken(ctx, userInfo.Email, accountDomain, isSystemAdmin)
	if err != nil {
		s.logger.Error("Failed to generate token",
			zap.Error(err),
			zap.Any("type", req.Type),
			zap.String("code", req.Code),
			zap.String("redirect_url", req.RedirectUrl),
		)
		return nil, err
	}

	return &authproto.ExchangeTokenResponse{Token: token}, nil
}

func (s *authService) RefreshToken(
	ctx context.Context,
	req *authproto.RefreshTokenRequest,
) (*authproto.RefreshTokenResponse, error) {
	if err := validateRefreshTokenRequest(req); err != nil {
		s.logger.Error("Failed to validate refresh token request",
			zap.Error(err),
			zap.String("refresh_token", req.RefreshToken),
		)
		return nil, err
	}
	refreshToken, err := s.verifier.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		s.logger.Error("Refresh token is invalid", zap.Any("refresh_token", refreshToken))
		return nil, statusUnauthenticated.Err()
	}
	organizations, err := s.getOrganizationsByEmail(ctx, refreshToken.Email)
	if err != nil {
		s.logger.Error("Failed to get organizations by email",
			zap.Error(err),
			zap.String("email", refreshToken.Email),
			zap.String("refresh_token", req.RefreshToken),
		)
		return nil, err
	}

	// Check if the user has at least one account enabled in any Organization
	account, err := s.checkAccountStatus(ctx, refreshToken.Email, organizations)
	if err != nil {
		s.logger.Error("Failed to check account",
			zap.Error(err),
			zap.String("email", refreshToken.Email),
			zap.Any("organizations", organizations),
		)
		return nil, err
	}
	accountDomain := domain.AccountV2{AccountV2: account.Account}
	isSystemAdmin := s.hasSystemAdminOrganization(organizations)

	newToken, err := s.generateToken(ctx, refreshToken.Email, accountDomain, isSystemAdmin)
	if err != nil {
		s.logger.Error(
			"Failed to generate token",
			zap.Error(err),
			zap.String("email", refreshToken.Email),
			zap.Any("organizations", organizations),
			zap.Any("refresh_token", refreshToken),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &authproto.RefreshTokenResponse{Token: newToken}, nil
}

func (s *authService) SwitchOrganization(
	ctx context.Context,
	req *authproto.SwitchOrganizationRequest,
) (*authproto.SwitchOrganizationResponse, error) {
	newOrganizationID := req.OrganizationId

	// Verify the access token
	accessToken, err := s.verifier.VerifyAccessToken(req.AccessToken)
	if err != nil {
		fields := log.FieldsFromIncomingContext(ctx)
		s.logger.Error(
			"Failed to verify access token",
			append(fields, zap.Error(err))...,
		)
		return nil, statusUnauthenticated.Err()
	}

	// Get the organizations that the user belongs to
	organizations, err := s.getOrganizationsByEmail(ctx, accessToken.Email)
	if err != nil {
		return nil, err
	}

	isSystemAdmin := s.hasSystemAdminOrganization(organizations)

	if isSystemAdmin {
		account, err := s.checkAccountStatus(ctx, accessToken.Email, organizations)
		if err != nil {
			return nil, err
		}
		accountDomain := domain.AccountV2{AccountV2: account.Account}
		if account.Account.Disabled {
			s.logger.Error(
				"The account is disabled",
				zap.String("email", accessToken.Email),
				zap.String("organizationID", newOrganizationID),
			)
		}
		accountDomain.OrganizationId = newOrganizationID
		token, err := s.generateToken(
			ctx,
			accessToken.Email,
			accountDomain,
			isSystemAdmin,
		)
		if err != nil {
			return nil, err
		}
		return &authproto.SwitchOrganizationResponse{
			Token: token,
		}, nil
	}
	account, err := s.accountClient.GetAccountV2(ctx, &acproto.GetAccountV2Request{
		Email:          accessToken.Email,
		OrganizationId: newOrganizationID,
	})
	if err != nil {
		s.logger.Error(
			"Failed to get account",
			zap.Error(err),
			zap.String("email", accessToken.Email),
			zap.String("organizationID", newOrganizationID),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	accountDomain := domain.AccountV2{AccountV2: account.Account}
	if account.Account.Disabled {
		s.logger.Error(
			"The account is disabled",
			zap.String("email", accessToken.Email),
			zap.String("organizationID", newOrganizationID),
		)
		return nil, statusUnauthenticated.Err()
	}

	// Generate new tokens with the new organization ID
	newToken, err := s.generateToken(
		ctx,
		accessToken.Email,
		accountDomain,
		isSystemAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &authproto.SwitchOrganizationResponse{
		Token: newToken,
	}, nil
}

func (s *authService) getAuthenticator(
	authType authproto.AuthType,
) (auth.Authenticator, error) {
	var authenticator auth.Authenticator
	switch authType {
	case authproto.AuthType_AUTH_TYPE_GOOGLE:
		authenticator = s.googleAuthenticator
	case authproto.AuthType_AUTH_TYPE_GITHUB:

	default:
		s.logger.Error("Unknown auth type", zap.String("authType", authType.String()))
		return nil, statusUnknownAuthType.Err()
	}
	return authenticator, nil
}

func (s *authService) getOrganizationsByEmail(
	ctx context.Context,
	email string,
) ([]*envproto.Organization, error) {
	orgResp, err := s.accountClient.GetMyOrganizationsByEmail(
		ctx,
		&acproto.GetMyOrganizationsByEmailRequest{
			Email: email,
		},
	)
	if err != nil {
		s.logger.Error(
			"Failed to get account's organizations",
			zap.Error(err),
			zap.String("email", email),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if len(orgResp.Organizations) == 0 {
		s.logger.Error(
			"The account is not registered in any organization",
			zap.String("email", email),
		)
		return nil, statusUnapprovedAccount.Err()
	}
	return orgResp.Organizations, nil
}

func (s *authService) updateUserInfoForOrganizations(
	ctx context.Context,
	userInfo *auth.UserInfo,
	organizations []*envproto.Organization,
) {
	for _, org := range organizations {
		account, err := s.accountClient.GetAccountV2(ctx, &acproto.GetAccountV2Request{
			Email:          userInfo.Email,
			OrganizationId: org.Id,
		})
		if err != nil {
			// Because we don't know what organization the user belongs when exchanging the token,
			// we ignore logs that were not found to avoid unnecessary logging.
			if status.Code(err) != codes.NotFound {
				s.logger.Error(
					"Failed to get account",
					zap.Error(err),
					zap.String("email", userInfo.Email),
					zap.String("organizationId", org.Id),
				)
			}
			continue
		}

		if account.Account.LastSeen == 0 {
			// Download avatar image if URL exists
			var avatarBytes []byte
			if userInfo.Avatar != "" {
				avatarBytes, err = s.downloadAvatar(ctx, userInfo.Avatar)
				if err != nil {
					s.logger.Error(
						"Failed to download avatar image",
						zap.Error(err),
						zap.String("avatarUrl", userInfo.Avatar),
					)
					// Continue with update even if avatar download fails
				}
			}

			// Handle empty first/last names from Google by providing fallbacks
			firstName := userInfo.FirstName
			lastName := userInfo.LastName

			// If first name is empty, try to extract from full name or use email prefix
			if firstName == "" {
				if userInfo.Name != "" {
					parts := strings.Fields(userInfo.Name)
					if len(parts) > 0 {
						firstName = parts[0]
					}
				}
				if firstName == "" {
					// Use email prefix as fallback
					emailParts := strings.Split(userInfo.Email, "@")
					if len(emailParts) > 0 {
						firstName = emailParts[0]
					}
				}
			}

			// If last name is empty, try to extract from full name or use default
			if lastName == "" {
				if userInfo.Name != "" {
					parts := strings.Fields(userInfo.Name)
					if len(parts) > 1 {
						lastName = strings.Join(parts[1:], " ")
					}
				}
				if lastName == "" {
					s.logger.Warn("Last name is empty. Using default fallback",
						zap.String("name", userInfo.Name),
						zap.String("last_name", lastName),
					)
					lastName = "User" // Default fallback
				}
			}

			updateReq := &acproto.UpdateAccountV2Request{
				Email:          userInfo.Email,
				OrganizationId: org.Id,
				FirstName:      wrapperspb.String(firstName),
				LastName:       wrapperspb.String(lastName),
				AvatarImageUrl: wrapperspb.String(userInfo.Avatar),
				Language:       wrapperspb.String("en"), // Default language
			}

			if len(avatarBytes) > 0 {
				updateReq.Avatar = &acproto.UpdateAccountV2Request_AccountV2Avatar{
					AvatarImage:    avatarBytes,
					AvatarFileType: "image/png",
				}
			}

			_, err = s.accountClient.UpdateAccountV2(ctx, updateReq)
			if err != nil {
				s.logger.Error(
					"Failed to update account first name, last name or avatar",
					zap.Error(err),
					zap.String("email", userInfo.Email),
					zap.String("organizationId", org.Id),
				)
			}
		}
	}
}

func (s *authService) downloadAvatar(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download avatar: status code %d", resp.StatusCode)
	}

	// Read response body into byte slice
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *authService) generateToken(
	ctx context.Context,
	userEmail string,
	accountDomain domain.AccountV2,
	isSystemAdmin bool,
) (*authproto.Token, error) {
	if err := s.checkEmail(userEmail); err != nil {
		s.logger.Error(
			"Access denied email",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", userEmail))...,
		)
		return nil, err
	}

	// Use the account's organization ID
	organizationID := accountDomain.OrganizationId

	// Create access token
	timeNow := time.Now()
	accessToken := &token.AccessToken{
		Issuer:         s.issuer,
		Audience:       s.audience,
		Expiry:         timeNow.Add(s.opts.accessTokenTTL),
		IssuedAt:       timeNow,
		Email:          userEmail,
		OrganizationID: organizationID,
		Name:           accountDomain.GetAccountFullName(),
		IsSystemAdmin:  isSystemAdmin,
	}

	signedAccessToken, err := s.signer.SignAccessToken(accessToken)
	if err != nil {
		s.logger.Error(
			"Failed to sign access token",
			zap.Error(err),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	// Create refresh token
	refreshToken := &token.RefreshToken{
		Email:          userEmail,
		OrganizationID: organizationID,
		Expiry:         timeNow.Add(s.opts.refreshTokenTTL),
		IssuedAt:       timeNow,
	}

	signedRefreshToken, err := s.signer.SignRefreshToken(refreshToken)
	if err != nil {
		s.logger.Error(
			"Failed to sign refresh token",
			zap.Error(err),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	return &authproto.Token{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		TokenType:    "Bearer",
		Expiry:       timeNow.Add(s.opts.accessTokenTTL).Unix(),
	}, nil
}

func (s *authService) checkEmail(email string) error {
	if s.opts.emailFilter == nil {
		return nil
	}
	if s.opts.emailFilter.MatchString(email) {
		return nil
	}
	return statusAccessDeniedEmail.Err()
}

// Check if the user has at least one account enabled in any Organization
func (s *authService) checkAccountStatus(
	ctx context.Context,
	email string,
	organizations []*envproto.Organization,
) (*acproto.GetAccountV2Response, error) {
	for _, org := range organizations {
		resp, err := s.accountClient.GetAccountV2(ctx, &acproto.GetAccountV2Request{
			Email:          email,
			OrganizationId: org.Id,
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// System admin accounts have access to all organizations,
				// but they are registered only in the system admin organization.
				// So, to avoid false errors, we ignore them if the account wasn't found in non-system admin organizations.
				continue
			}
			return nil, api.NewGRPCStatus(err).Err()
		}
		if !resp.Account.Disabled {
			// The account must have at least one account enabled
			return resp, nil
		}
	}
	// The account wasn't found or doesn't belong to any organization
	return nil, statusUnauthenticated.Err()
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
	config := s.config.DemoSignIn
	if !config.Enabled {
		s.logger.Info("Demo sign in is disabled")
		return
	}
	if config.Email == "" ||
		config.Password == "" ||
		config.OrganizationId == "" ||
		config.ProjectId == "" ||
		config.EnvironmentId == "" {
		s.logger.Error("One or more demo sign-in configuration are missing")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	var err error
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		// Create a demo organization if not exists
		_, err = s.organizationStorage.GetOrganization(contextWithTx, config.OrganizationId)
		if err != nil {
			if errors.Is(err, envstotage.ErrOrganizationNotFound) {
				err = s.organizationStorage.CreateOrganization(contextWithTx, &envdomain.Organization{
					Organization: &envproto.Organization{
						Id:          config.OrganizationId,
						Name:        "Demo organization",
						UrlCode:     "demo",
						OwnerEmail:  config.OrganizationOwnerEmail,
						Description: "This organization is for demo users",
						Disabled:    false,
						Archived:    false,
						Trial:       false,
						CreatedAt:   now.Unix(),
						UpdatedAt:   now.Unix(),
						SystemAdmin: config.IsSystemAdmin,
						AuthenticationSettings: &envproto.AuthenticationSettings{
							EnabledTypes: []envproto.AuthenticationType{
								envproto.AuthenticationType_AUTHENTICATION_TYPE_GOOGLE,
								envproto.AuthenticationType_AUTHENTICATION_TYPE_PASSWORD,
							},
						},
					},
				})
			}
			if err != nil && !errors.Is(err, envstotage.ErrOrganizationAlreadyExists) {
				return err
			}
		}
		// Create a demo project if not exists
		_, err = s.projectStorage.GetProject(contextWithTx, config.ProjectId)
		if err != nil {
			if errors.Is(err, envstotage.ErrProjectNotFound) {
				err = s.projectStorage.CreateProject(contextWithTx, &envdomain.Project{
					Project: &envproto.Project{
						Id:             config.ProjectId,
						Description:    "This project is for demo users",
						Disabled:       false,
						Trial:          false,
						CreatorEmail:   config.Email,
						CreatedAt:      now.Unix(),
						UpdatedAt:      now.Unix(),
						Name:           "Demo",
						UrlCode:        "demo",
						OrganizationId: config.OrganizationId,
					},
				})
				if err != nil && !errors.Is(err, envstotage.ErrProjectAlreadyExists) {
					return err
				}
			}
		}
		// Create a demo environment if not exists
		_, err = s.environmentStorage.GetEnvironmentV2(contextWithTx, config.EnvironmentId)
		if err != nil {
			if errors.Is(err, envstotage.ErrEnvironmentNotFound) {
				err = s.environmentStorage.CreateEnvironmentV2(contextWithTx, &envdomain.EnvironmentV2{
					EnvironmentV2: &envproto.EnvironmentV2{
						Id:             config.EnvironmentId,
						Name:           "Demo",
						UrlCode:        "demo",
						Description:    "This environment is for demo users",
						ProjectId:      config.ProjectId,
						Archived:       false,
						CreatedAt:      now.Unix(),
						UpdatedAt:      now.Unix(),
						OrganizationId: config.OrganizationId,
						RequireComment: false,
					},
				})
				if err != nil && !errors.Is(err, envstotage.ErrEnvironmentAlreadyExists) {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Error("Failed to prepare demo environment", zap.Error(err))
		return
	}
	// Create a demo account if not exists
	accountStorage := accountstotage.NewAccountStorage(s.mysqlClient)
	_, err = accountStorage.GetAccountV2(
		ctx,
		config.Email,
		config.OrganizationId,
	)
	if err != nil {
		if errors.Is(err, accountstotage.ErrAccountNotFound) {
			err = accountStorage.CreateAccountV2(ctx, &domain.AccountV2{
				AccountV2: &acproto.AccountV2{
					OrganizationId:   config.OrganizationId,
					Email:            config.Email,
					FirstName:        "Bucketeer",
					LastName:         "Demo",
					Language:         "en",
					OrganizationRole: acproto.AccountV2_Role_Organization_ADMIN,
					EnvironmentRoles: []*acproto.AccountV2_EnvironmentRole{
						{
							EnvironmentId: config.EnvironmentId,
							Role:          acproto.AccountV2_Role_Environment_EDITOR,
						},
					},
					CreatedAt: now.Unix(),
					UpdatedAt: now.Unix(),
				},
			})
			if err != nil && !errors.Is(err, accountstotage.ErrAccountAlreadyExists) {
				s.logger.Error("Create account for demo user error", zap.Error(err))
			}
		}
	}

	// Create credentials for demo user if not exists
	_, err = s.credentialsStorage.GetCredentials(ctx, config.Email)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsNotFound) {
			passwordHash, hashErr := auth.HashPassword(config.Password)
			if hashErr != nil {
				s.logger.Error("Failed to hash demo user password", zap.Error(hashErr))
				return
			}
			err = s.credentialsStorage.CreateCredentials(ctx, config.Email, passwordHash)
			if err != nil && !errors.Is(err, storage.ErrCredentialsAlreadyExists) {
				s.logger.Error("Failed to create credentials for demo user", zap.Error(err))
				return
			}
		} else {
			s.logger.Error("Failed to check credentials for demo user", zap.Error(err))
			return
		}
	}
	s.logger.Info("Demo environment prepared successfully")
}
