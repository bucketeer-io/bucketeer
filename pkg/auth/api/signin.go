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
	if s.config.Password.Enabled {
		return s.handlePasswordSignIn(ctx, request)
	}

	// If neither is enabled nor credentials don't match, deny access
	s.logger.Error("Sign in failed - no valid authentication method",
		zap.String("email", request.Email),
		zap.Bool("passwordAuthEnabled", s.config.Password.Enabled),
	)
	return nil, statusAccessDenied.Err()
}

func (s *authService) handlePasswordSignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
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
		return nil, statusAccessDenied.Err()
	}

	// Verify password
	if !auth.ValidatePassword(request.Password, credentials.PasswordHash) {
		s.logger.Error("Password sign in failed - invalid password",
			zap.String("email", email),
		)
		return nil, statusAccessDenied.Err()
	}

	// Get organizations for the user
	organizations, err := s.getOrganizationsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Check account status - user should have at least one enabled account
	account, err := s.checkAccountStatus(ctx, email, organizations)
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
	token, err := s.generateToken(ctx, email, accountDomain, isSystemAdmin)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Successful password authentication", zap.String("email", email))
	return &authproto.SignInResponse{Token: token}, nil
}

func (s *authService) SignInPassword(
	ctx context.Context,
	request *authproto.SignInPasswordRequest,
) (*authproto.SignInPasswordResponse, error) {
	err := validateSignInPasswordRequest(request)
	if err != nil {
		s.logger.Error("SignInPassword request validation failed", zap.Error(err))
		return nil, err
	}

	// Get credentials
	credentials, err := s.credentialsStorage.GetCredentials(ctx, request.Email)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsNotFound) {
			s.logger.Error("SignInPassword failed - no credentials found",
				zap.String("email", request.Email),
			)
		} else {
			s.logger.Error("SignInPassword failed - credentials lookup error",
				zap.Error(err),
				zap.String("email", request.Email),
			)
		}
		return nil, statusAccessDenied.Err()
	}

	// Verify password
	if !auth.ValidatePassword(request.Password, credentials.PasswordHash) {
		s.logger.Error("SignInPassword failed - invalid password",
			zap.String("email", request.Email),
		)
		return nil, statusAccessDenied.Err()
	}

	// Create temporary token using existing helper
	userInfo := &auth.UserInfo{
		Email: request.Email,
		Name:  request.Email, // Use email as fallback name
	}

	token, err := s.createTemporaryToken(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	s.logger.Info("SignInPassword successful - temporary token issued",
		zap.String("email", request.Email),
	)

	return &authproto.SignInPasswordResponse{Token: token}, nil
}
