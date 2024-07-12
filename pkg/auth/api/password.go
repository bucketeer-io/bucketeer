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

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

func (s *authService) SignIn(
	ctx context.Context,
	request *authproto.SignInRequest,
) (*authproto.SignInResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validatePasswordLoginRequest(request, localizer)
	if err != nil {
		return nil, err
	}
	if request.Username != s.config.DemoSignInConfig.Username &&
		request.Password != s.config.DemoSignInConfig.Password {
		s.logger.Error(
			"Password login failed",
			zap.String("username", request.Username),
			zap.String("password", request.Password),
		)
		dt, err := auth.StatusPasswordAccessDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}
	token, err := s.generateToken(ctx, s.config.DemoSignInConfig.Email, localizer)
	if err != nil {
		return nil, err
	}
	return &authproto.SignInResponse{Token: token}, nil
}
