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
	"errors"

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
	t, ok := rpc.GetIDToken(ctx)
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
			log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", t.Email))...,
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
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
				log.FieldsFromImcomingContext(ctx).AddFields(
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
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
		return &accountproto.GetMeResponse{Account: &accountproto.ConsoleAccount{
			Email:            sysAdminAccount.Email,
			Name:             sysAdminAccount.Name,
			AvatarUrl:        sysAdminAccount.AvatarImageUrl,
			IsSystemAdmin:    true,
			Organization:     organization,
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			EnvironmentRoles: adminEnvRoles,
		}}, nil
	}
	// non admin account response
	account, err := s.getAccountV2(ctx, t.Email, req.OrganizationId, localizer)
	if err != nil {
		return nil, err
	}
	envRoles := s.getConsoleAccountEnvironmentRoles(account.EnvironmentRoles, environments, projects)
	return &accountproto.GetMeResponse{Account: &accountproto.ConsoleAccount{
		Email:            account.Email,
		Name:             account.Name,
		AvatarUrl:        account.AvatarImageUrl,
		IsSystemAdmin:    false,
		Organization:     organization,
		OrganizationRole: account.OrganizationRole,
		EnvironmentRoles: envRoles,
	}}, nil
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
	t, ok := rpc.GetIDToken(ctx)
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
			log.FieldsFromImcomingContext(ctx).AddFields(zap.String("email", req.Email))...,
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
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
		if accWithOrg.AccountV2.Disabled || accWithOrg.Organization.Disabled || accWithOrg.Organization.Archived {
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
		if org.Organization.SystemAdmin {
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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
