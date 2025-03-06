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

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
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
	config := s.config.DemoSignIn
	if !config.Enabled ||
		request.Email != config.Email ||
		request.Password != config.Password {
		s.logger.Error(
			"Sign in failed",
			zap.Bool("enabled", config.Enabled),
			zap.String("email", request.Email),
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
	organizations, err := s.getOrganizationsByEmail(ctx, config.Email, localizer)
	if err != nil {
		return nil, err
	}

	// Check if the user has at least one account enabled in any Organization
	account, err := s.checkAccountStatus(ctx, config.Email, organizations, localizer)
	if err != nil {
		s.logger.Error("Failed to check account",
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
	return &authproto.SignInResponse{Token: token}, nil
}
