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
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func validateGetAuthenticationURLRequest(
	req *authproto.GetAuthenticationURLRequest,
) error {
	if req.Type == authproto.AuthType_AUTH_TYPE_UNSPECIFIED {
		return statusMissingAuthType.Err()
	}
	if req.State == "" {
		return statusMissingState.Err()
	}
	if req.RedirectUrl == "" {
		return statusMissingRedirectURL.Err()
	}
	return nil
}

func validateExchangeTokenRequest(
	req *authproto.ExchangeTokenRequest,
) error {
	if req.Type == authproto.AuthType_AUTH_TYPE_UNSPECIFIED {
		return statusMissingAuthType.Err()
	}
	if req.Code == "" {
		return statusMissingCode.Err()
	}
	if req.RedirectUrl == "" {
		return statusMissingRedirectURL.Err()
	}
	return nil
}

func validateRefreshTokenRequest(
	req *authproto.RefreshTokenRequest,
) error {
	if req.RefreshToken == "" {
		return statusMissingRefreshToken.Err()
	}
	return nil
}

func validateSignInRequest(
	req *authproto.SignInRequest,
) error {
	if req.Email == "" {
		return statusMissingUsername.Err()
	}
	if req.Password == "" {
		return statusMissingPassword.Err()
	}
	return nil
}

func validateRequestMagicLinkRequest(
	req *authproto.RequestMagicLinkRequest,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		dt, err := auth.StatusMissingEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	// Basic email format validation
	if !strings.Contains(req.Email, "@") {
		dt, err := auth.StatusMissingEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateVerifyMagicLinkRequest(
	req *authproto.VerifyMagicLinkRequest,
	localizer locale.Localizer,
) error {
	if req.Token == "" {
		dt, err := auth.StatusMissingToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "token"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}
