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

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func (s *authService) SignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
) (*authproto.SignInResponse, error) {
	err := validateSignInRequest(request)
	if err != nil {
		return nil, err
	}

	// Then, try password authentication if enabled
	if s.config.PasswordAuth.Enabled {
		return s.handlePasswordSignIn(ctx, request, localizer)
	}

	// If neither is enabled nor credentials don't match, deny access
	s.logger.Error("Sign in failed - no valid authentication method",
		zap.String("email", request.Email),
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
