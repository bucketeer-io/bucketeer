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
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func validateUpdatePasswordRequest(
	req *authproto.UpdatePasswordRequest,
	localizer locale.Localizer,
) error {
	if req.CurrentPassword == "" {
		dt, err := auth.StatusMissingCurrentPassword.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "current_password"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.NewPassword == "" {
		dt, err := auth.StatusMissingNewPassword.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "new_password"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.CurrentPassword == req.NewPassword {
		dt, err := auth.StatusPasswordsIdentical.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PasswordsIdentical),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateInitiatePasswordSetupRequest(
	req *authproto.InitiatePasswordSetupRequest,
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
	return nil
}

func validateSetupPasswordRequest(
	req *authproto.SetupPasswordRequest,
	localizer locale.Localizer,
) error {
	if req.SetupToken == "" {
		dt, err := auth.StatusMissingResetToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "setup_token"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.NewPassword == "" {
		dt, err := auth.StatusMissingNewPassword.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "new_password"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validatePasswordSetupTokenRequest(
	req *authproto.ValidatePasswordSetupTokenRequest,
	localizer locale.Localizer,
) error {
	if req.SetupToken == "" {
		dt, err := auth.StatusMissingResetToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "setup_token"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateInitiatePasswordResetRequest(
	req *authproto.InitiatePasswordResetRequest,
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
	return nil
}

func validateResetPasswordRequest(
	req *authproto.ResetPasswordRequest,
	localizer locale.Localizer,
) error {
	if req.ResetToken == "" {
		dt, err := auth.StatusMissingResetToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "reset_token"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	if req.NewPassword == "" {
		dt, err := auth.StatusMissingNewPassword.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "new_password"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validatePasswordResetTokenRequest(
	req *authproto.ValidatePasswordResetTokenRequest,
	localizer locale.Localizer,
) error {
	if req.ResetToken == "" {
		dt, err := auth.StatusMissingResetToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "reset_token"),
		})
		if err != nil {
			return auth.StatusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}
