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
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	authstorage "github.com/bucketeer-io/bucketeer/v2/pkg/auth/storage"
)

func (s *AccountService) GetMe(
	ctx context.Context,
	req *accountproto.GetMeRequest,
) (*accountproto.GetMeResponse, error) {
	t, ok := rpc.GetAccessToken(ctx)
	if !ok {
		return nil, statusUnauthenticated.Err()
	}
	if !verifyEmailFormat(t.Email) {
		s.logger.Error(
			"Email inside IDToken has an invalid format",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", t.Email))...,
		)
		return nil, statusInvalidEmail.Err()
	}
	projects, err := s.listProjectsByOrganizationID(ctx, req.OrganizationId)
	if err != nil {
		s.logger.Error(
			"Failed to get project list",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	environments, err := s.listEnvironmentsByOrganizationID(ctx, req.OrganizationId)
	if err != nil {
		s.logger.Error(
			"Failed to get environment list",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)

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
			return nil, statusOrganizationNotFound.Err()
		}
		s.logger.Error(
			"Failed to get organization",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	// system admin account response
	sysAdminAccount, err := s.getSystemAdminAccountV2(ctx, t.Email)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	if sysAdminAccount != nil && !sysAdminAccount.Disabled {
		adminEnvRoles := s.getAdminConsoleAccountEnvironmentRoles(environments, projects)
		lastSeen := sysAdminAccount.LastSeen

		if updated, err := s.updateLastSeen(ctx, sysAdminAccount.Email, req.OrganizationId); err != nil {
			if errors.Is(err, v2as.ErrAccountNotFound) {
				s.logger.Warn(
					"System admin user not found in organization when updating last seen",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.String("email", sysAdminAccount.Email),
						zap.String("organizationId", req.OrganizationId),
					)...,
				)
			} else {
				s.logger.Error(
					"Failed to update system admin user last seen",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("email", sysAdminAccount.Email),
						zap.String("organizationId", req.OrganizationId),
					)...,
				)
			}
		} else {
			lastSeen = updated
		}
		if req.OrganizationId != sysAdminAccount.OrganizationId {
			// req.OrganizationId is the org currently being viewed in the console, while
			// sysAdminAccount.OrganizationId is the system admin's dedicated organization.
			// When these differ, the system admin is looking at a regular organization,
			// so we also update the last-seen timestamp for the system admin organization.
			if _, err := s.updateLastSeen(ctx, sysAdminAccount.Email, sysAdminAccount.OrganizationId); err != nil {
				s.logger.Error(
					"Failed to update system admin user last seen",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("email", sysAdminAccount.Email),
						zap.String("organizationId", sysAdminAccount.OrganizationId),
					)...,
				)
			}
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
			LastSeen:         lastSeen,
			PasswordSetupRequired: s.checkPasswordSetupRequired(ctx, t.Email, organization),
		}}, nil
	}
	// non admin account response
	account, err := s.getAccount(ctx, t.Email, req.OrganizationId)
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
	lastSeen := account.LastSeen
	var updated int64
	updated, err = s.updateLastSeen(ctx, account.Email, req.OrganizationId)
	if err != nil {
		s.logger.Error(
			"Failed to update user last seen",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", account.Email),
				zap.String("organizationId", req.OrganizationId),
			)...,
		)
	} else {
		lastSeen = updated
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
		LastSeen:         lastSeen,
		PasswordSetupRequired: s.checkPasswordSetupRequired(ctx, t.Email, organization),
	}}, nil
}

// isPasswordAuthenticationEnabled checks if password authentication is enabled in the organization
func (s *AccountService) isPasswordAuthenticationEnabled(organization *environmentproto.Organization) bool {
	if organization.AuthenticationSettings == nil {
		return false
	}

	// Check if PASSWORD authentication type is in the enabled types
	for _, authType := range organization.AuthenticationSettings.EnabledTypes {
		if authType == environmentproto.AuthenticationType_AUTHENTICATION_TYPE_PASSWORD {
			return true
		}
	}
	return false
}

// checkPasswordSetupRequired determines if the user needs to set up a password
// It checks both whether the user has existing credentials AND if the organization allows password authentication
// When setup is required, it proactively creates empty credentials for the frontend password setup flow
func (s *AccountService) checkPasswordSetupRequired(ctx context.Context, email string,
	organization *environmentproto.Organization) bool {
	// First check if the organization allows password authentication
	if organization != nil && !s.isPasswordAuthenticationEnabled(organization) {
		// Organization has disabled password authentication, no setup required
		return false
	}

	// Check if user already has password credentials
	credentials, err := s.credentialsStorage.GetCredentials(ctx, email)
	if err == nil && credentials.PasswordHash != "" {
		// User already has password, no setup required
		return false
	}
	if err != nil && !errors.Is(err, authstorage.ErrCredentialsNotFound) {
		// Real error occurred, log and assume no setup required to be safe
		s.logger.Warn("Failed to check credentials for password setup status",
			zap.Error(err),
			zap.String("email", email))
		return false
	}

	// At this point: credentials either don't exist OR exist with empty password hash
	// If credentials don't exist, create empty credentials record for frontend password setup flow
	if err != nil && errors.Is(err, authstorage.ErrCredentialsNotFound) {
		err = s.credentialsStorage.CreateCredentials(ctx, email, "")
		if err != nil {
			// If creation fails, log and assume no setup required to be safe
			s.logger.Warn("Failed to create empty credentials for password setup",
				zap.Error(err),
				zap.String("email", email))
			return false
		}
		s.logger.Info("Created empty credentials for password setup", zap.String("email", email))
	}

	// User needs password setup (either new credentials created or existing empty credentials)
	return true
}

// getAccount also checks if the account exists or is disabled
func (s *AccountService) getAccount(
	ctx context.Context,
	email string,
	organizationID string,
) (*accountproto.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
			s.logger.Error("Account not found",
				zap.String("email", email),
				zap.String("organizationId", organizationID),
			)
			return nil, statusUnauthenticated.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	if account.Disabled {
		s.logger.Error("Account is disabled",
			zap.String("email", email),
			zap.String("organizationId", organizationID),
		)
		return nil, statusUnauthenticated.Err()
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
	t, ok := rpc.GetAccessToken(ctx)
	if !ok {
		return nil, statusUnauthenticated.Err()
	}
	myOrgs, err := s.getMyOrganizations(ctx, t.Email)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetMyOrganizationsResponse{Organizations: myOrgs}, nil
}

func (s *AccountService) GetMyOrganizationsByEmail(
	ctx context.Context,
	req *accountproto.GetMyOrganizationsByEmailRequest,
) (*accountproto.GetMyOrganizationsResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !verifyEmailFormat(req.Email) {
		s.logger.Error(
			"Email inside request has an invalid format",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.String("email", req.Email))...,
		)
		return nil, statusInvalidEmail.Err()
	}
	myOrgs, err := s.getMyOrganizations(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetMyOrganizationsResponse{Organizations: myOrgs}, nil
}

func (s *AccountService) getMyOrganizations(
	ctx context.Context,
	email string,
) ([]*environmentproto.Organization, error) {
	accountsWithOrg, err := s.accountStorage.GetAccountsWithOrganization(ctx, email)
	if err != nil {
		s.logger.Error(
			"Failed to get accounts with organization",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
			return nil, api.NewGRPCStatus(err).Err()
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
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetSystemAdminAccountV2(ctx, email)
	if err != nil {
		if errors.Is(err, v2as.ErrSystemAdminAccountNotFound) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to get system admin account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", email),
			)...,
		)
	}
	return account, nil
}

func (s *AccountService) updateLastSeen(ctx context.Context, email, organizationID string) (int64, error) {
	// First get the existing account
	account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
	if err != nil {
		return 0, err
	}

	cloned, ok := proto.Clone(account.AccountV2).(*accountproto.AccountV2)
	if !ok {
		return 0, statusInternal.Err()
	}

	now := time.Now().Unix()
	// Update only the LastSeen field
	accountForUpdate := &domain.AccountV2{AccountV2: cloned}
	accountForUpdate.LastSeen = now
	accountForUpdate.UpdatedAt = now

	if err := s.accountStorage.UpdateAccountV2(ctx, accountForUpdate); err != nil {
		return 0, err
	}
	return now, nil
}
