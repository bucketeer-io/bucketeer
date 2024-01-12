// Copyright 2023 The Bucketeer Authors.
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
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

var (
	// nolint:lll
	emailRegex           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	maxAccountNameLength = 250
)

func validateGetAdminAccountRequest(req *accountproto.GetAdminAccountRequest, localizer locale.Localizer) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateCreateAdminAccountRequest(req *accountproto.CreateAdminAccountRequest, localizer locale.Localizer) error {
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Command.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateEnableAdminAccountRequest(req *accountproto.EnableAdminAccountRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAccountID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateDisableAdminAccountRequest(
	req *accountproto.DisableAdminAccountRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusMissingAccountID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateGetAccountRequest(req *accountproto.GetAccountRequest, localizer locale.Localizer) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateCreateAccountRequest(req *accountproto.CreateAccountRequest, localizer locale.Localizer) error {
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Command.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateChangeAccountRoleRequest(req *accountproto.ChangeAccountRoleRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAccountID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateConvertAccountRequest(req *accountproto.ConvertAccountRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAccountID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateEnableAccountRequest(req *accountproto.EnableAccountRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAccountID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateDisableAccountRequest(req *accountproto.DisableAccountRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAccountID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "account_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func verifyEmailFormat(email string) bool {
	return emailRegex.MatchString(email)
}

func validateCreateAPIKeyRequest(req *accountproto.CreateAPIKeyRequest, localizer locale.Localizer) error {
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.Name == "" {
		dt, err := statusMissingAPIKeyName.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateChangeAPIKeyNameRequest(req *accountproto.ChangeAPIKeyNameRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAPIKeyID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateEnableAPIKeyRequest(req *accountproto.EnableAPIKeyRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAPIKeyID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateDisableAPIKeyRequest(req *accountproto.DisableAPIKeyRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusMissingAPIKeyID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateCreateAccountV2Request(req *accountproto.CreateAccountV2Request, localizer locale.Localizer) error {
	if req.OrganizationId == "" {
		dt, err := statusMissingOrganizationID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Command.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.Name == "" {
		dt, err := statusNameIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.Command.Name) > maxAccountNameLength {
		dt, err := statusInvalidName.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.OrganizationRole == accountproto.AccountV2_Role_Organization_UNASSIGNED {
		dt, err := statusInvalidOrganizationRole.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateUpdateAccountV2Request(
	req *accountproto.UpdateAccountV2Request,
	commands []command.Command,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.OrganizationId == "" {
		dt, err := statusMissingOrganizationID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(commands) == 0 {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, cmd := range commands {
		if c, ok := cmd.(*accountproto.ChangeAccountV2NameCommand); ok {
			newName := strings.TrimSpace(c.Name)
			if newName == "" {
				dt, err := statusNameIsEmpty.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
			if len(newName) > maxAccountNameLength {
				dt, err := statusInvalidName.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2OrganizationRoleCommand); ok {
			if c.Role == accountproto.AccountV2_Role_Organization_UNASSIGNED {
				dt, err := statusInvalidOrganizationRole.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
		}
	}
	return nil
}

func validateEnableAccountV2Request(req *accountproto.EnableAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.OrganizationId == "" {
		dt, err := statusMissingOrganizationID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateDisableAccountV2Request(req *accountproto.DisableAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.OrganizationId == "" {
		dt, err := statusMissingOrganizationID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateDeleteAccountV2Request(req *accountproto.DeleteAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.OrganizationId == "" {
		dt, err := statusMissingOrganizationID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateGetAccountV2Request(req *accountproto.GetAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.OrganizationId == "" {
		dt, err := statusMissingOrganizationID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateGetAccountV2ByEnvironmentIDRequest(
	req *accountproto.GetAccountV2ByEnvironmentIDRequest,
	localizer locale.Localizer,
) error {
	// We don't check the environmentID because there is environment with empty ID.
	if req.Email == "" {
		dt, err := statusEmailIsEmpty.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !verifyEmailFormat(req.Email) {
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}
