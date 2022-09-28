// Copyright 2022 The Bucketeer Authors.
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
	"regexp"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

// nolint:lll
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func validateGetAdminAccountRequest(req *accountproto.GetAdminAccountRequest) error {
	if req.Email == "" {
		return localizedError(statusEmailIsEmpty, locale.JaJP)
	}
	if !verifyEmailFormat(req.Email) {
		return localizedError(statusInvalidEmail, locale.JaJP)
	}
	return nil
}

func validateCreateAdminAccountRequest(req *accountproto.CreateAdminAccountRequest) error {
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if req.Command.Email == "" {
		return localizedError(statusEmailIsEmpty, locale.JaJP)
	}
	if !verifyEmailFormat(req.Command.Email) {
		return localizedError(statusInvalidEmail, locale.JaJP)
	}
	return nil
}

func validateEnableAdminAccountRequest(req *accountproto.EnableAdminAccountRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAccountID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func validateDisableAdminAccountRequest(req *accountproto.DisableAdminAccountRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAccountID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func validateGetAccountRequest(req *accountproto.GetAccountRequest) error {
	if req.Email == "" {
		return localizedError(statusEmailIsEmpty, locale.JaJP)
	}
	if !verifyEmailFormat(req.Email) {
		return localizedError(statusInvalidEmail, locale.JaJP)
	}
	return nil
}

func validateCreateAccountRequest(req *accountproto.CreateAccountRequest) error {
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if req.Command.Email == "" {
		return localizedError(statusEmailIsEmpty, locale.JaJP)
	}
	if !verifyEmailFormat(req.Command.Email) {
		return localizedError(statusInvalidEmail, locale.JaJP)
	}
	return nil
}

func validateChangeAccountRoleRequest(req *accountproto.ChangeAccountRoleRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAccountID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func validateConvertAccountRequest(req *accountproto.ConvertAccountRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAccountID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func validateEnableAccountRequest(req *accountproto.EnableAccountRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAccountID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func validateDisableAccountRequest(req *accountproto.DisableAccountRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAccountID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func verifyEmailFormat(email string) bool {
	return emailRegex.MatchString(email)
}

func validateCreateAPIKeyRequest(req *accountproto.CreateAPIKeyRequest) error {
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if req.Command.Name == "" {
		return localizedError(statusMissingAPIKeyName, locale.JaJP)
	}
	return nil
}

func validateChangeAPIKeyNameRequest(req *accountproto.ChangeAPIKeyNameRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAPIKeyID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func validateEnableAPIKeyRequest(req *accountproto.EnableAPIKeyRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAPIKeyID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func validateDisableAPIKeyRequest(req *accountproto.DisableAPIKeyRequest) error {
	if req.Id == "" {
		return localizedError(statusMissingAPIKeyID, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}
