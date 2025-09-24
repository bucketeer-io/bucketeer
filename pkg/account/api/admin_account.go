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
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func (s *AccountService) GetMe(
	ctx context.Context,
	req *accountproto.GetMeRequest,
) (*accountproto.GetMeResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	t, ok := rpc.GetAccessToken(ctx)
	if !ok {
		dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.UnauthenticatedError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if !verifyEmailFormat(t.Email) {
		s.logger.Error(
			"Email inside IDToken has an invalid format",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", t.Email))...,
		)
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	projects, err := s.listProjectsByOrganizationID(ctx, req.OrganizationId)
	if err != nil {
		s.logger.Error(
			"Failed to get project list",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	environments, err := s.listEnvironmentsByOrganizationID(ctx, req.OrganizationId)
	if err != nil {
		s.logger.Error(
			"Failed to get environment list",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	organization, err := s.getOrganization(ctx, req.OrganizationId)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			s.logger.Error(
				"Organization not found",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", req.OrganizationId),
				)...,
			)
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get organization",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	// system admin account response
	sysAdminAccount, err := s.getSystemAdminAccountV2(ctx, t.Email, localizer)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	if sysAdminAccount != nil && !sysAdminAccount.Disabled {
		adminEnvRoles := s.getAdminConsoleAccountEnvironmentRoles(environments, projects)

		// update system admin user last seen
		err := s.updateLastSeen(ctx, sysAdminAccount.Email, sysAdminAccount.OrganizationId)
		if err != nil {
			s.logger.Error(
				"Failed to update system admin user last seen",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("email", sysAdminAccount.Email),
					zap.String("organizationId", req.OrganizationId),
				)...,
			)
		}

		return &accountproto.GetMeResponse{Account: &accountproto.ConsoleAccount{
			Email:            sysAdminAccount.Email,
			Name:             sysAdminAccount.Name,
			AvatarUrl:        sysAdminAccount.AvatarImageUrl,
			AvatarFileType:   sysAdminAccount.AvatarFileType,
			AvatarImage:      sysAdminAccount.AvatarImage,
			IsSystemAdmin:    true,
			Organization:     organization,
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			EnvironmentRoles: adminEnvRoles,
			SearchFilters:    sysAdminAccount.SearchFilters,
			FirstName:        sysAdminAccount.FirstName,
			LastName:         sysAdminAccount.LastName,
			Language:         sysAdminAccount.Language,
			LastSeen:         sysAdminAccount.LastSeen,
		}}, nil
	}
	// non admin account response
	account, err := s.getAccount(ctx, t.Email, req.OrganizationId, localizer)
	if err != nil {
		return nil, err
	}
	var envRoles []*accountproto.ConsoleAccount_EnvironmentRole
	if account.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
		envRoles = s.getConsoleAccountEnvironmentRoles(account.EnvironmentRoles, environments, projects)
	} else {
		// If the user is an admin or owner, no need to filter environments.
		envRoles = s.getAdminConsoleAccountEnvironmentRoles(environments, projects)
	}

	// update user last seen
	err = s.updateLastSeen(ctx, account.Email, req.OrganizationId)
	if err != nil {
		s.logger.Error(
			"Failed to update user last seen",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", account.Email),
				zap.String("organizationId", req.OrganizationId),
			)...,
		)
	}

	return &accountproto.GetMeResponse{Account: &accountproto.ConsoleAccount{
		Email:            account.Email,
		Name:             account.Name,
		AvatarUrl:        account.AvatarImageUrl,
		AvatarFileType:   account.AvatarFileType,
		AvatarImage:      account.AvatarImage,
		IsSystemAdmin:    false,
		Organization:     organization,
		OrganizationRole: account.OrganizationRole,
		EnvironmentRoles: envRoles,
		SearchFilters:    account.SearchFilters,
		FirstName:        account.FirstName,
		LastName:         account.LastName,
		Language:         account.Language,
		LastSeen:         account.LastSeen,
	}}, nil
}

// getAccount also checks if the account exists or is disabled
func (s *AccountService) getAccount(
	ctx context.Context,
	email string,
	organizationID string,
	localizer locale.Localizer,
) (*accountproto.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
			s.logger.Error("Account not found",
				zap.String("email", email),
				zap.String("organizationId", organizationID),
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if account.Disabled {
		s.logger.Error("Account is disabled",
			zap.String("email", email),
			zap.String("organizationId", organizationID),
		)
		dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.UnauthenticatedError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return account.AccountV2, nil
}

func (s *AccountService) getAdminConsoleAccountEnvironmentRoles(
	environments []*environmentproto.EnvironmentV2,
	projects []*environmentproto.Project,
) []*accountproto.ConsoleAccount_EnvironmentRole {
	projectSet := s.makeProjectSet(projects)
	environmentRoles := make([]*accountproto.ConsoleAccount_EnvironmentRole, 0, len(environments))
	for _, e := range environments {
		if e.Archived {
			continue
		}
		p, ok := projectSet[e.ProjectId]
		if !ok || p.Disabled {
			continue
		}
		er := &accountproto.ConsoleAccount_EnvironmentRole{
			Environment: e,
			Project:     p,
			Role:        accountproto.AccountV2_Role_Environment_EDITOR,
		}
		environmentRoles = append(environmentRoles, er)
	}
	return environmentRoles
}

func (s *AccountService) getConsoleAccountEnvironmentRoles(
	roles []*accountproto.AccountV2_EnvironmentRole,
	environments []*environmentproto.EnvironmentV2,
	projects []*environmentproto.Project,
) []*accountproto.ConsoleAccount_EnvironmentRole {
	envSet := s.makeEnvironmentSet(environments)
	projectSet := s.makeProjectSet(projects)
	environmentRoles := make([]*accountproto.ConsoleAccount_EnvironmentRole, 0, len(roles))
	for _, r := range roles {
		env, ok := envSet[r.EnvironmentId]
		if !ok || env.Archived {
			continue
		}
		project, ok := projectSet[env.ProjectId]
		if !ok || project.Disabled {
			continue
		}
		// TODO: Remove this checking after the web console 3.0 is ready
		// If the account is enabled in any environment in this organization,
		// we append the organization.
		// Note: When we disable an account on the web console,
		// we are updating the role to UNASSIGNED, not the `disabled` column.
		// When the new console is ready, we will use the DisableAccount API instead,
		// which will update the `disabled` column in the DB.
		if r.Role == accountproto.AccountV2_Role_Environment_UNASSIGNED {
			continue
		}
		environmentRoles = append(environmentRoles, &accountproto.ConsoleAccount_EnvironmentRole{
			Environment: env,
			Role:        r.Role,
			Project:     project,
		})
	}
	return environmentRoles
}

func (s *AccountService) GetMyOrganizations(
	ctx context.Context,
	_ *accountproto.GetMyOrganizationsRequest,
) (*accountproto.GetMyOrganizationsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	t, ok := rpc.GetAccessToken(ctx)
	if !ok {
		dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.UnauthenticatedError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	myOrgs, err := s.getMyOrganizations(ctx, t.Email, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetMyOrganizationsResponse{Organizations: myOrgs}, nil
}

func (s *AccountService) GetMyOrganizationsByEmail(
	ctx context.Context,
	req *accountproto.GetMyOrganizationsByEmailRequest,
) (*accountproto.GetMyOrganizationsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if !verifyEmailFormat(req.Email) {
		s.logger.Error(
			"Email inside request has an invalid format",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", req.Email))...,
		)
		dt, err := statusInvalidEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "email"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	myOrgs, err := s.getMyOrganizations(ctx, req.Email, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetMyOrganizationsResponse{Organizations: myOrgs}, nil
}

func (s *AccountService) getMyOrganizations(
	ctx context.Context,
	email string,
	localizer locale.Localizer,
) ([]*environmentproto.Organization, error) {
	accountsWithOrg, err := s.accountStorage.GetAccountsWithOrganization(ctx, email)
	if err != nil {
		s.logger.Error(
			"Failed to get accounts with organization",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if s.containsSystemAdminOrganization(accountsWithOrg) {
		resp, err := s.environmentClient.ListOrganizations(
			ctx,
			&environmentproto.ListOrganizationsRequest{
				Disabled: wrapperspb.Bool(false),
				Archived: wrapperspb.Bool(false),
			})
		if err != nil {
			s.logger.Error(
				"Failed to get organizations",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return resp.Organizations, nil
	}
	myOrgs := make([]*environmentproto.Organization, 0, len(accountsWithOrg))
	for _, accWithOrg := range accountsWithOrg {
		if accWithOrg.AccountV2.Disabled || accWithOrg.Organization.Disabled || accWithOrg.Archived {
			continue
		}
		// Add the organization if the account is an admin or owner.
		// Otherwise, we check if the account is enabled in any environment in this organization.
		if accWithOrg.OrganizationRole >= accountproto.AccountV2_Role_Organization_ADMIN {
			myOrgs = append(myOrgs, accWithOrg.Organization)
			continue
		}
		// TODO: Remove this loop after the web console 3.0 is ready
		// If the account is enabled in any environment in this organization,
		// we append the organization.
		// Note: When we disable an account on the web console,
		// we are updating the role to UNASSIGNED, not the `disabled` column.
		// When the new console is ready, we will use the DisableAccount API instead,
		// which will update the `disabled` column in the DB.
		var enabled bool
		for _, role := range accWithOrg.EnvironmentRoles {
			if role.Role != accountproto.AccountV2_Role_Environment_UNASSIGNED {
				enabled = true
			}
		}
		if !enabled {
			continue
		}
		myOrgs = append(myOrgs, accWithOrg.Organization)
	}
	return myOrgs, nil
}

func (s *AccountService) containsSystemAdminOrganization(
	organizations []*domain.AccountWithOrganization,
) bool {
	for _, org := range organizations {
		if org.SystemAdmin {
			return true
		}
	}
	return false
}

func (s *AccountService) getSystemAdminAccountV2(
	ctx context.Context,
	email string,
	localizer locale.Localizer,
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetSystemAdminAccountV2(ctx, email)
	if err != nil {
		if errors.Is(err, v2as.ErrSystemAdminAccountNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get system admin account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", email),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return account, nil
}

func (s *AccountService) updateLastSeen(ctx context.Context, email, organizationID string) error {
	// First get the existing account
	account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	// Update only the LastSeen field
	account.LastSeen = now
	account.UpdatedAt = now

	return s.accountStorage.UpdateAccountV2(ctx, account)
}
