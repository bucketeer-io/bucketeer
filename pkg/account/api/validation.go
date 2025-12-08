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
	"regexp"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/command"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
)

var (
	// nolint:lll
	emailRegex           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	maxAccountNameLength = 250
)

func verifyEmailFormat(email string) bool {
	return emailRegex.MatchString(email)
}

func validateCreateAPIKeyRequest(req *accountproto.CreateAPIKeyRequest) error {
	if req.Command.Name == "" {
		return statusMissingAPIKeyName.Err()
	}
	return nil
}

func validateCreateAPIKeyRequestNoCommand(req *accountproto.CreateAPIKeyRequest) error {
	if req.Name == "" {
		return statusMissingAPIKeyName.Err()
	}
	if req.Maintainer != "" && !verifyEmailFormat(req.Maintainer) {
		return statusInvalidEmail.Err()
	}
	return nil
}

func validateChangeAPIKeyNameRequest(req *accountproto.ChangeAPIKeyNameRequest) error {
	if req.Id == "" {
		return statusMissingAPIKeyID.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func validateEnableAPIKeyRequest(req *accountproto.EnableAPIKeyRequest) error {
	if req.Id == "" {
		return statusMissingAPIKeyID.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func validateDisableAPIKeyRequest(req *accountproto.DisableAPIKeyRequest) error {
	if req.Id == "" {
		return statusMissingAPIKeyID.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func validateCreateAccountV2Request(req *accountproto.CreateAccountV2Request) error {
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	if req.Command.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Command.Email) {
		return statusInvalidEmail.Err()
	}
	if req.Command.OrganizationRole == accountproto.AccountV2_Role_Organization_UNASSIGNED {
		return statusInvalidOrganizationRole.Err()
	}
	if req.Command.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
		if len(req.Command.EnvironmentRoles) == 0 {
			return statusInvalidEnvironmentRole.Err()
		}
	}
	return nil
}

func validateCreateAccountV2NoCommandRequest(
	req *accountproto.CreateAccountV2Request,
) error {
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	if req.OrganizationRole == accountproto.AccountV2_Role_Organization_UNASSIGNED {
		return statusInvalidOrganizationRole.Err()
	}
	if req.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
		if len(req.EnvironmentRoles) == 0 {
			return statusInvalidEnvironmentRole.Err()
		}
	}
	return nil
}

func validateUpdateAccountV2Request(
	req *accountproto.UpdateAccountV2Request,
	commands []command.Command,
) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	if len(commands) == 0 {
		return statusNoCommand.Err()
	}
	for _, cmd := range commands {
		if c, ok := cmd.(*accountproto.ChangeAccountV2FirstNameCommand); ok {
			newFirstName := strings.TrimSpace(c.FirstName)
			if newFirstName == "" {
				return statusFirstNameIsEmpty.Err()
			}
			if len(newFirstName) > maxAccountNameLength {
				return statusInvalidFirstName.Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2LastNameCommand); ok {
			newLastName := strings.TrimSpace(c.LastName)
			if newLastName == "" {
				return statusLastNameIsEmpty.Err()
			}
			if len(newLastName) > maxAccountNameLength {
				return statusInvalidLastName.Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2LanguageCommand); ok {
			if c.Language == "" {
				return statusLanguageIsEmpty.Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2OrganizationRoleCommand); ok {
			if c.Role == accountproto.AccountV2_Role_Organization_UNASSIGNED {
				return statusInvalidOrganizationRole.Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2EnvironmentRolesCommand); ok {
			if c.WriteType == accountproto.ChangeAccountV2EnvironmentRolesCommand_WriteType_UNSPECIFIED {
				return statusInvalidUpdateEnvironmentRolesWriteType.Err()
			}
			if len(c.Roles) == 0 {
				return statusInvalidEnvironmentRole.Err()
			}
		}
	}
	return nil
}

func validateUpdateAccountV2NoCommandRequest(
	req *accountproto.UpdateAccountV2Request,
) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	if req.FirstName != nil {
		newFirstName := strings.TrimSpace(req.FirstName.Value)
		if newFirstName == "" {
			return statusFirstNameIsEmpty.Err()
		}
		if len(newFirstName) > maxAccountNameLength {
			return statusInvalidFirstName.Err()
		}
	}
	if req.LastName != nil {
		newLastName := strings.TrimSpace(req.LastName.Value)
		if newLastName == "" {
			return statusLastNameIsEmpty.Err()
		}
		if len(newLastName) > maxAccountNameLength {
			return statusInvalidLastName.Err()
		}
	}
	if req.Language != nil {
		if req.Language.Value == "" {
			return statusLanguageIsEmpty.Err()
		}
	}
	if req.OrganizationRole != nil && req.OrganizationRole.Role == accountproto.AccountV2_Role_Organization_UNASSIGNED {
		return statusInvalidOrganizationRole.Err()
	}
	for _, r := range req.EnvironmentRoles {
		if r.Role == accountproto.AccountV2_Role_Environment_UNASSIGNED {
			return statusInvalidEnvironmentRole.Err()
		}
	}
	return nil
}

func validateEnableAccountV2Request(req *accountproto.EnableAccountV2Request) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	return nil
}

func validateDisableAccountV2Request(req *accountproto.DisableAccountV2Request) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	return nil
}

func validateDeleteAccountV2Request(req *accountproto.DeleteAccountV2Request) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	return nil
}

func validateGetAccountV2Request(req *accountproto.GetAccountV2Request) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	return nil
}

func validateGetAccountV2ByEnvironmentIDRequest(
	req *accountproto.GetAccountV2ByEnvironmentIDRequest,
) error {
	// We don't check the environmentID because there is environment with empty ID.
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if !verifyEmailFormat(req.Email) {
		return statusInvalidEmail.Err()
	}
	return nil
}

func validateCreateSearchFilterRequest(
	req *accountproto.CreateSearchFilterRequest,
) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	if req.Command.Name == "" {
		return statusSearchFilterNameIsEmpty.Err()
	}
	if req.Command.Query == "" {
		return statusSearchFilterQueryIsEmpty.Err()
	}
	if req.Command.FilterTargetType == accountproto.FilterTargetType_UNKNOWN {
		return statusSearchFilterTargetTypeIsRequired.Err()
	}
	return nil
}

func validateUpdateSearchFilterRequest(
	req *accountproto.UpdateSearchFilterRequest,
	commands []command.Command,
) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}

	if len(commands) == 0 {
		return statusNoCommand.Err()
	}
	for _, cmd := range commands {
		switch cmd := cmd.(type) {
		case *accountproto.ChangeSearchFilterNameCommand:
			if err := validateChangeSearchFilterId(cmd.Id); err != nil {
				return err
			}
			if cmd.Name == "" {
				return statusSearchFilterNameIsEmpty.Err()
			}
		case *accountproto.ChangeSearchFilterQueryCommand:
			if err := validateChangeSearchFilterId(cmd.Id); err != nil {
				return err
			}
			if cmd.Query == "" {
				return statusSearchFilterQueryIsEmpty.Err()
			}
		case *accountproto.ChangeDefaultSearchFilterCommand:
			if err := validateChangeSearchFilterId(cmd.Id); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateChangeSearchFilterId(id string) error {
	if id == "" {
		return statusSearchFilterIDIsEmpty.Err()
	}
	return nil
}

func validateDeleteSearchFilterRequest(
	req *accountproto.DeleteSearchFilterRequest,
) error {
	if req.Email == "" {
		return statusEmailIsEmpty.Err()
	}
	if req.OrganizationId == "" {
		return statusMissingOrganizationID.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	if req.Command.Id == "" {
		return statusSearchFilterIDIsEmpty.Err()
	}
	return nil
}

func validateUpdateAPIKeyRequestNoCommand(req *accountproto.UpdateAPIKeyRequest) error {
	if req.Id == "" {
		return statusMissingAPIKeyID.Err()
	}
	return nil
}
