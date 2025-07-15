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

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"
	grpcapi "github.com/bucketeer-io/bucketeer/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

var (
	// nolint:lll
	emailRegex           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	maxAccountNameLength = 250
)

func verifyEmailFormat(email string) bool {
	return emailRegex.MatchString(email)
}

func validateCreateAPIKeyRequest(req *accountproto.CreateAPIKeyRequest, localizer locale.Localizer) error {
	if req.Command.Name == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"api_key_name"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_name"),
			},
		).Err()
	}
	return nil
}

func validateCreateAPIKeyRequestNoCommand(req *accountproto.CreateAPIKeyRequest, localizer locale.Localizer) error {
	if req.Name == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"api_key_name"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_name"),
			},
		).Err()
	}
	if req.Maintainer != "" && !verifyEmailFormat(req.Maintainer) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"maintainer"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "maintainer"),
			},
		).Err()
	}
	return nil
}

func validateChangeAPIKeyNameRequest(req *accountproto.ChangeAPIKeyNameRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"api_key_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
			},
		).Err()
	}
	if req.Command == nil {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"command"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			},
		).Err()
	}
	return nil
}

func validateEnableAPIKeyRequest(req *accountproto.EnableAPIKeyRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"api_key_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
			},
		).Err()
	}
	if req.Command == nil {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"command"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			},
		).Err()
	}
	return nil
}

func validateDisableAPIKeyRequest(req *accountproto.DisableAPIKeyRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"api_key_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
			},
		).Err()
	}
	if req.Command == nil {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"command"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			},
		).Err()
	}
	return nil
}

func validateCreateAccountV2Request(req *accountproto.CreateAccountV2Request, localizer locale.Localizer) error {
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	if req.Command.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Command.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.Command.OrganizationRole == accountproto.AccountV2_Role_Organization_UNASSIGNED {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_role"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role"),
			},
		).Err()
	}
	if req.Command.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
		if len(req.Command.EnvironmentRoles) == 0 {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"environment_roles"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "environment_roles"),
				},
			).Err()
		}
	}
	return nil
}

func validateCreateAccountV2NoCommandRequest(
	req *accountproto.CreateAccountV2Request,
	localizer locale.Localizer,
) error {
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.OrganizationRole == accountproto.AccountV2_Role_Organization_UNASSIGNED {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_role"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role"),
			},
		).Err()
	}
	if req.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
		if len(req.EnvironmentRoles) == 0 {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"environment_roles"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "environment_roles"),
				},
			).Err()
		}
	}
	return nil
}

func validateUpdateAccountV2Request(
	req *accountproto.UpdateAccountV2Request,
	commands []command.Command,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	if len(commands) == 0 {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"command"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			},
		).Err()
	}
	for _, cmd := range commands {
		if c, ok := cmd.(*accountproto.ChangeAccountV2FirstNameCommand); ok {
			newFirstName := strings.TrimSpace(c.FirstName)
			if newFirstName == "" {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"first_name"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "first_name"),
					},
				).Err()
			}
			if len(newFirstName) > maxAccountNameLength {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"first_name"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "first_name"),
					},
				).Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2LastNameCommand); ok {
			newLastName := strings.TrimSpace(c.LastName)
			if newLastName == "" {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"last_name"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "last_name"),
					},
				).Err()
			}
			if len(newLastName) > maxAccountNameLength {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"last_name"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "last_name"),
					},
				).Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2LanguageCommand); ok {
			if c.Language == "" {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"language"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "language"),
					},
				).Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2OrganizationRoleCommand); ok {
			if c.Role == accountproto.AccountV2_Role_Organization_UNASSIGNED {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"organization_role"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role"),
					},
				).Err()
			}
		}
		if c, ok := cmd.(*accountproto.ChangeAccountV2EnvironmentRolesCommand); ok {
			if c.WriteType == accountproto.ChangeAccountV2EnvironmentRolesCommand_WriteType_UNSPECIFIED {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"environment_role_write_type"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "environment_role_write_type"),
					},
				).Err()
			}
			if len(c.Roles) == 0 {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"environment_roles"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "environment_roles"),
					},
				).Err()
			}
		}
	}
	return nil
}

func validateUpdateAccountV2NoCommandRequest(
	req *accountproto.UpdateAccountV2Request,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	if req.FirstName != nil {
		newFirstName := strings.TrimSpace(req.FirstName.Value)
		if newFirstName == "" {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"first_name"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "first_name"),
				},
			).Err()
		}
		if len(newFirstName) > maxAccountNameLength {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"first_name"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "first_name"),
				},
			).Err()
		}
	}
	if req.LastName != nil {
		newLastName := strings.TrimSpace(req.LastName.Value)
		if newLastName == "" {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"last_name"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "last_name"),
				},
			).Err()
		}
		if len(newLastName) > maxAccountNameLength {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"last_name"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "last_name"),
				},
			).Err()
		}
	}
	if req.Language != nil {
		if req.Language.Value == "" {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"language"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "language"),
				},
			).Err()
		}
	}
	if req.OrganizationRole != nil && req.OrganizationRole.Role == accountproto.AccountV2_Role_Organization_UNASSIGNED {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_role"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "organization_role"),
			},
		).Err()
	}
	for _, r := range req.EnvironmentRoles {
		if r.Role == accountproto.AccountV2_Role_Environment_UNASSIGNED {
			return grpcapi.NewGRPCStatus(
				pkgErr.NewErrorInvalidAugment("account", []string{"environment_role"}),
				"account",
				&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "environment_role"),
				},
			).Err()
		}
	}
	return nil
}

func validateEnableAccountV2Request(req *accountproto.EnableAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	return nil
}

func validateDisableAccountV2Request(req *accountproto.DisableAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	return nil
}

func validateDeleteAccountV2Request(req *accountproto.DeleteAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	return nil
}

func validateGetAccountV2Request(req *accountproto.GetAccountV2Request, localizer locale.Localizer) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	return nil
}

func validateGetAccountV2ByEnvironmentIDRequest(
	req *accountproto.GetAccountV2ByEnvironmentIDRequest,
	localizer locale.Localizer,
) error {
	// We don't check the environmentID because there is environment with empty ID.
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if !verifyEmailFormat(req.Email) {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
			},
		).Err()
	}
	return nil
}

func validateCreateSearchFilterRequest(
	req *accountproto.CreateSearchFilterRequest,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	if req.Command == nil {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"command"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			},
		).Err()
	}
	if req.Command.Name == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"name"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
			},
		).Err()
	}
	if req.Command.Query == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"query"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "query"),
			},
		).Err()
	}
	if req.Command.FilterTargetType == accountproto.FilterTargetType_UNKNOWN {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"filter_target_type"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "filter_target_type"),
			},
		).Err()
	}
	return nil
}

func validateUpdateSearchFilterRequest(
	req *accountproto.UpdateSearchFilterRequest,
	commands []command.Command,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}

	if len(commands) == 0 {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"command"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			},
		).Err()
	}
	for _, cmd := range commands {
		switch cmd := cmd.(type) {
		case *accountproto.ChangeSearchFilterNameCommand:
			if err := validateChangeSearchFilterId(cmd.Id, localizer); err != nil {
				return err
			}
			if cmd.Name == "" {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"name"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
					},
				).Err()
			}
		case *accountproto.ChangeSearchFilterQueryCommand:
			if err := validateChangeSearchFilterId(cmd.Id, localizer); err != nil {
				return err
			}
			if cmd.Query == "" {
				return grpcapi.NewGRPCStatus(
					pkgErr.NewErrorInvalidAugment("account", []string{"query"}),
					"account",
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "query"),
					},
				).Err()
			}
		case *accountproto.ChangeDefaultSearchFilterCommand:
			if err := validateChangeSearchFilterId(cmd.Id, localizer); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateChangeSearchFilterId(id string, localizer locale.Localizer) error {
	if id == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
			},
		).Err()
	}
	return nil
}

func validateDeleteSearchFilterRequest(
	req *accountproto.DeleteSearchFilterRequest,
	localizer locale.Localizer,
) error {
	if req.Email == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"email"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "email"),
			},
		).Err()
	}
	if req.OrganizationId == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"organization_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
			},
		).Err()
	}
	if req.Command == nil {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"command"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			},
		).Err()
	}
	if req.Command.Id == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"search_filter_id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "search_filter_id"),
			},
		).Err()
	}
	return nil
}

func validateUpdateAPIKeyRequestNoCommand(req *accountproto.UpdateAPIKeyRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		return grpcapi.NewGRPCStatus(
			pkgErr.NewErrorInvalidAugment("account", []string{"id"}),
			"account",
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
			},
		).Err()
	}
	return nil
}
