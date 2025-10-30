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
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/storage"
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

	// First, try demo authentication if enabled
	demoConfig := s.config.DemoSignIn
	if demoConfig.Enabled &&
		request.Email == demoConfig.Email &&
		request.Password == demoConfig.Password {
		return s.handleDemoSignIn(ctx, request, localizer)
	}

	// Then, try password authentication if enabled
	if s.config.PasswordAuth.Enabled {
		return s.handlePasswordSignIn(ctx, request, localizer)
	}

	// If neither is enabled or credentials don't match, deny access
	s.logger.Error("Sign in failed - no valid authentication method",
		zap.String("email", request.Email),
		zap.Bool("demoEnabled", demoConfig.Enabled),
		zap.Bool("passwordAuthEnabled", s.config.PasswordAuth.Enabled),
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

func (s *authService) handleDemoSignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
	localizer locale.Localizer,
) (*authproto.SignInResponse, error) {
	config := s.config.DemoSignIn
	organizations, err := s.getOrganizationsByEmail(ctx, config.Email, localizer)
	if err != nil {
		return nil, err
	}

	// Check if the user has at least one account enabled in any Organization
	account, err := s.checkAccountStatus(ctx, config.Email, organizations, localizer)
	if err != nil {
		s.logger.Error("Failed to check account for demo sign in",
			zap.Error(err),
			zap.String("email", config.Email),
			zap.Any("organizations", organizations),
		)
		return nil, err
	}
	accountDomain := domain.AccountV2{AccountV2: account.Account}
	isSystemAdmin := s.hasSystemAdminOrganization(organizations)

	token, err := s.generateToken(ctx, config.Email, accountDomain, isSystemAdmin, localizer)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Successful demo authentication", zap.String("email", config.Email))
	return &authproto.SignInResponse{Token: token}, nil
}

func (s *authService) handlePasswordSignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
	localizer locale.Localizer,
) (*authproto.SignInResponse, error) {
	// Sanitize email
	email := auth.SanitizeEmail(request.Email)
	if !auth.IsValidEmail(email) {
		s.logger.Error("Invalid email format for password sign in", zap.String("email", email))
		dt, err := auth.StatusAccessDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

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
