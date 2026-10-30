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
	"crypto/subtle"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

const (
	// signInDemoAccountKey is the throttling key for sign-in attempts that
	// target the configured demo account.
	signInDemoAccountKey = "demo-account"
	// signInOtherAccountKey is the throttling key for sign-in attempts that
	// target any other email. Such attempts can never succeed.
	signInOtherAccountKey = "other-account"
)

func (s *authService) SignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
) (*authproto.SignInResponse, error) {
	err := validateSignInRequest(request)
	if err != nil {
		return nil, err
	}
	config := s.config.DemoSignIn
	// Both credentials are compared in constant time, and both comparisons
	// always run, so the response time does not depend on how much of the
	// submitted values match the configured ones.
	emailMatches := subtle.ConstantTimeCompare([]byte(request.Email), []byte(config.Email)) == 1
	passwordMatches := subtle.ConstantTimeCompare([]byte(request.Password), []byte(config.Password)) == 1
	// Only the configured demo account can ever sign in, so attempts for that
	// account share one bucket and attempts for any other email share another.
	// This slows down repeated guesses against the demo account regardless of
	// where they come from while keeping the number of tracked keys bounded.
	key := signInOtherAccountKey
	if emailMatches {
		key = signInDemoAccountKey
	}
	if err := s.throttle(ctx, s.opts.signInRateLimiter, "sign_in", key); err != nil {
		return nil, err
	}
	if !config.Enabled || !emailMatches || !passwordMatches {
		s.logger.Error(
			"Sign in failed",
			zap.Bool("enabled", config.Enabled),
			zap.String("email", request.Email),
		)
		return nil, statusAccessDenied.Err()
	}
	organizations, err := s.getOrganizationsByEmail(ctx, config.Email)
	if err != nil {
		return nil, err
	}

	// Check if the user has at least one account enabled in any Organization
	account, err := s.checkAccountStatus(ctx, config.Email, organizations)
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

	token, err := s.generateToken(ctx, config.Email, accountDomain, isSystemAdmin)
	if err != nil {
		return nil, err
	}
	return &authproto.SignInResponse{Token: token}, nil
}
