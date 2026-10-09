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
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func validateUpdatePasswordRequest(
	req *authproto.UpdatePasswordRequest,
) error {
	if req.CurrentPassword == "" {
		return statusMissingCurrentPassword.Err()
	}
	if req.NewPassword == "" {
		return statusMissingNewPassword.Err()
	}
	if req.CurrentPassword == req.NewPassword {
		return statusPasswordsIdentical.Err()
	}
	return nil
}

func validateInitiatePasswordSetupRequest(
	req *authproto.InitiatePasswordSetupRequest,
) error {
	if req.Email == "" {
		return statusMissingEmail.Err()
	}
	return nil
}

func validateSetupPasswordRequest(
	req *authproto.SetupPasswordRequest,
) error {
	if req.SetupToken == "" {
		return statusMissingResetToken.Err()
	}
	if req.NewPassword == "" {
		return statusMissingNewPassword.Err()
	}
	return nil
}

func validatePasswordSetupTokenRequest(
	req *authproto.ValidatePasswordSetupTokenRequest,
) error {
	if req.SetupToken == "" {
		return statusMissingResetToken.Err()
	}
	return nil
}

func validateInitiatePasswordResetRequest(
	req *authproto.InitiatePasswordResetRequest,
) error {
	if req.Email == "" {
		return statusMissingEmail.Err()
	}
	return nil
}

func validateResetPasswordRequest(
	req *authproto.ResetPasswordRequest,
) error {
	if req.ResetToken == "" {
		return statusMissingResetToken.Err()
	}
	if req.NewPassword == "" {
		return statusMissingNewPassword.Err()
	}
	return nil
}

func validatePasswordResetTokenRequest(
	req *authproto.ValidatePasswordResetTokenRequest,
) error {
	if req.ResetToken == "" {
		return statusMissingResetToken.Err()
	}
	return nil
}
