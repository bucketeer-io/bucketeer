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

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	acproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func (s *authService) SignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
) (*authproto.SignInResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validateSignInRequest(request, localizer)
	if err != nil {
		return nil, err
	}

	// Check if password authentication is enabled
	if !s.config.Password.Enabled {
		s.logger.Error("Password authentication not enabled")
		dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// NEW FLOW: If organization_id is provided, use organization-specific authentication
	if request.OrganizationId != "" {
		return s.handlePasswordSignInWithOrganization(ctx, request, localizer)
	}

	// OLD FLOW: Backward compatibility - authenticate with first organization
	return s.handlePasswordSignIn(ctx, request, localizer)
}

func (s *authService) handlePasswordSignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
	localizer locale.Localizer,
) (*authproto.SignInResponse, error) {
	email := request.Email

	// Get credentials
	credentials, err := s.credentialsStorage.GetCredentials(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsNotFound) {
			s.logger.Error("Password sign in failed - no credentials found",
				zap.String("email", email),
			)
		} else {
			s.logger.Error("Password sign in failed - credentials lookup error",
				zap.Error(err),
				zap.String("email", email),
			)
		}
		dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Verify password
	if !auth.ValidatePassword(request.Password, credentials.PasswordHash) {
		s.logger.Error("Password sign in failed - invalid password",
			zap.String("email", email),
		)
		dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Get organizations for the user
	organizations, err := s.getOrganizationsByEmail(ctx, email, localizer)
	if err != nil {
		return nil, err
	}

	// Check account status - user should have at least one enabled account
	account, err := s.checkAccountStatus(ctx, email, organizations, localizer)
	if err != nil {
		s.logger.Error("Failed to check account status for password sign in",
			zap.Error(err),
			zap.String("email", email),
		)
		return nil, err
	}

	accountDomain := domain.AccountV2{AccountV2: account.Account}
	isSystemAdmin := s.hasSystemAdminOrganization(organizations)

	// Generate token for successful authentication
	token, err := s.generateToken(ctx, email, accountDomain, isSystemAdmin, localizer)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Successful password authentication", zap.String("email", email))
	return &authproto.SignInResponse{Token: token}, nil
}

// handlePasswordSignInWithOrganization handles password sign-in with organization selection (new flow)
func (s *authService) handlePasswordSignInWithOrganization(
	ctx context.Context,
	request *authproto.SignInRequest,
	localizer locale.Localizer,
) (*authproto.SignInResponse, error) {
	email := request.Email
	organizationID := request.OrganizationId

	// Step 1: Verify password
	credentials, err := s.credentialsStorage.GetCredentials(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsNotFound) {
			s.logger.Error("Password sign in failed - no credentials found",
				zap.String("email", email),
			)
		} else {
			s.logger.Error("Password sign in failed - credentials lookup error",
				zap.Error(err),
				zap.String("email", email),
			)
		}
		dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	if !auth.ValidatePassword(request.Password, credentials.PasswordHash) {
		s.logger.Error("Password sign in failed - invalid password",
			zap.String("email", email),
		)
		dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Step 2: Get all organizations for this user
	organizations, err := s.getOrganizationsByEmail(ctx, email, localizer)
	if err != nil {
		return nil, err
	}

	// Step 3: Verify the requested organization is in the user's organization list
	var selectedOrg *envproto.Organization
	for _, org := range organizations {
		if org.Id == organizationID {
			selectedOrg = org
			break
		}
	}
	if selectedOrg == nil {
		s.logger.Error("Organization not found for user",
			zap.String("email", email),
			zap.String("organization_id", organizationID),
		)
		dt, err := auth.StatusOrganizationNotFound.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.NotFoundError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Step 4: Verify PASSWORD auth is enabled for this organization
	if !isAuthMethodEnabled(selectedOrg, envproto.AuthenticationType_AUTHENTICATION_TYPE_PASSWORD) {
		s.logger.Error("Password authentication not enabled for organization",
			zap.String("email", email),
			zap.String("organization_id", organizationID),
		)
		dt, err := auth.StatusAuthMethodNotEnabled.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Step 5: Get account for the specific organization
	account, err := s.accountClient.GetAccountV2(ctx, &acproto.GetAccountV2Request{
		Email:          email,
		OrganizationId: organizationID,
	})
	if err != nil {
		s.logger.Error("Failed to get account for organization",
			zap.Error(err),
			zap.String("email", email),
			zap.String("organization_id", organizationID),
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	// Check if account is disabled
	if account.Account.Disabled {
		s.logger.Error("Account is disabled",
			zap.String("email", email),
			zap.String("organization_id", organizationID),
		)
		dt, err := auth.StatusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.UnauthenticatedError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	accountDomain := domain.AccountV2{AccountV2: account.Account}
	isSystemAdmin := s.hasSystemAdminOrganization(organizations)

	// Step 6: Generate token for the specified organization
	token, err := s.generateToken(ctx, email, accountDomain, isSystemAdmin, localizer)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Successful password authentication with organization selection",
		zap.String("email", email),
		zap.String("organization_id", organizationID),
	)
	return &authproto.SignInResponse{Token: token}, nil
}

// isAuthMethodEnabled checks if a specific authentication method is enabled for an organization
func isAuthMethodEnabled(org *envproto.Organization, authType envproto.AuthenticationType) bool {
	if org.AuthenticationSettings == nil {
		return false
	}
	for _, enabledType := range org.AuthenticationSettings.EnabledTypes {
		if enabledType == authType {
			return true
		}
	}
	return false
}
