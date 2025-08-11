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
	"strconv"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainauditlog "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	tagdomain "github.com/bucketeer-io/bucketeer/v2/pkg/tag/domain"
	teamdomain "github.com/bucketeer-io/bucketeer/v2/pkg/team/domain"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	"github.com/bucketeer-io/bucketeer/v2/proto/common"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	tagproto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

func (s *AccountService) CreateAccountV2(
	ctx context.Context,
	req *accountproto.CreateAccountV2Request,
) (*accountproto.CreateAccountV2Response, error) {
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
	)
	if err != nil {
		return nil, err
	}
	err = validateCreateAccountV2Request(req)
	if err != nil {
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, err
	}
	account := domain.NewAccountV2(
		req.Email,
		req.Name,
		req.FirstName,
		req.LastName,
		req.Language,
		req.AvatarImageUrl,
		req.Tags,
		req.Teams,
		req.OrganizationId,
		req.OrganizationRole,
		req.EnvironmentRoles,
	)
	var createAccountEvent *eventproto.Event
	err = s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		createAccountEvent, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ACCOUNT,
			account.Email,
			eventproto.Event_ACCOUNT_V2_CREATED,
			&eventproto.AccountV2CreatedEvent{
				Email:            account.Email,
				FirstName:        account.FirstName,
				LastName:         account.LastName,
				Language:         account.Language,
				AvatarImageUrl:   account.AvatarImageUrl,
				OrganizationId:   account.OrganizationId,
				OrganizationRole: account.OrganizationRole,
				EnvironmentRoles: account.EnvironmentRoles,
				Disabled:         account.Disabled,
				CreatedAt:        account.CreatedAt,
				UpdatedAt:        account.UpdatedAt,
			},
			account,
			nil,
		)
		if err != nil {
			return err
		}
		err = s.accountStorage.CreateAccountV2(contextWithTx, account)
		if err != nil {
			return err
		}
		return s.adminAuditLogStorage.CreateAdminAuditLog(
			contextWithTx,
			domainauditlog.NewAuditLog(createAccountEvent, storage.AdminEnvironmentID),
		)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAccountAlreadyExists) {
			return nil, statusAccountAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	if err = s.publisher.Publish(ctx, createAccountEvent); err != nil {
		s.logger.Error(
			"Failed to publish create account event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	// Upsert tags -- deprecated
	for _, envRole := range req.EnvironmentRoles {
		if err := s.upsertTags(ctx, req.Tags, envRole.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to upsert account tags",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationId", req.OrganizationId),
					zap.String("environmentId", envRole.EnvironmentId),
					zap.String("email", req.Email),
					zap.Any("tags", req.Tags),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}

	if err := s.upsertTeams(ctx, req.Teams, req.OrganizationId); err != nil {
		s.logger.Error(
			"Failed to upsert account teams",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationId", req.OrganizationId),
				zap.String("email", req.Email),
				zap.Any("teams", req.Teams),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	// Initiate password setup for newly created account
	s.initiatePasswordSetupForNewAccount(ctx, req.Email)

	return &accountproto.CreateAccountV2Response{Account: account.AccountV2}, nil
}

func (s *AccountService) upsertTags(
	ctx context.Context,
	tags []string,
	environmentID string,
) error {
	for _, tag := range tags {
		trimed := strings.TrimSpace(tag)
		if trimed == "" {
			continue
		}
		t, err := tagdomain.NewTag(trimed, environmentID, tagproto.Tag_ACCOUNT)
		if err != nil {
			s.logger.Error(
				"Failed to create domain tag",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentID),
					zap.String("tagId", tag),
				)...,
			)
			return err
		}
		if err := s.tagStorage.UpsertTag(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func (s *AccountService) upsertTeams(
	ctx context.Context,
	teams []string,
	organizationID string,
) error {
	for _, team := range teams {
		trimed := strings.TrimSpace(team)
		if trimed == "" {
			continue
		}
		t, err := teamdomain.NewTeam(trimed, trimed, organizationID)
		if err != nil {
			s.logger.Error(
				"Failed to create domain team",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationId", organizationID),
					zap.String("teamId", team),
				)...,
			)
			return err
		}
		if err := s.teamStorage.UpsertTeam(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func (s *AccountService) UpdateAccountV2(
	ctx context.Context,
	req *accountproto.UpdateAccountV2Request,
) (*accountproto.UpdateAccountV2Response, error) {
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
	)
	if err != nil {
		// If not admin, check if user is updating their own account
		editor, err = s.checkOrganizationRole(
			ctx,
			accountproto.AccountV2_Role_Organization_MEMBER,
			req.OrganizationId,
		)
		if err != nil {
			return nil, err
		}
		if editor.Email != req.Email {
			return nil, statusPermissionDenied.Err()
		}
	}

	err = validateUpdateAccountV2Request(req)
	if err != nil {
		return nil, err
	}
	updatedAccountPb, err := s.updateAccountV2NoCommandMysql(
		ctx,
		editor,
		req.Email,
		req.OrganizationId,
		req.Name,
		req.FirstName,
		req.LastName,
		req.Language,
		req.AvatarImageUrl,
		req.Avatar,
		req.Tags,
		req.TeamChanges,
		req.OrganizationRole,
		req.EnvironmentRoles,
		req.Disabled,
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to update account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	// Upsert tags
	if req.Tags != nil {
		for _, envRole := range updatedAccountPb.EnvironmentRoles {
			if req.Tags == nil {
				continue
			}
			if err := s.upsertTags(ctx, req.Tags.Values, envRole.EnvironmentId); err != nil {
				s.logger.Error(
					"Failed to upsert account tags",
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("organizationId", req.OrganizationId),
						zap.String("environmentId", envRole.EnvironmentId),
						zap.String("email", req.Email),
						zap.Any("tags", req.Tags),
					)...,
				)
				return nil, api.NewGRPCStatus(err).Err()
			}
		}
	}
	if updatedAccountPb.Teams != nil {
		if err := s.upsertTeams(ctx, updatedAccountPb.Teams, req.OrganizationId); err != nil {
			s.logger.Error(
				"Failed to upsert account teams",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationId", req.OrganizationId),
					zap.String("email", req.Email),
					zap.Any("teams", updatedAccountPb.Teams),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return &accountproto.UpdateAccountV2Response{
		Account: updatedAccountPb,
	}, nil
}

func (s *AccountService) EnableAccountV2(
	ctx context.Context,
	req *accountproto.EnableAccountV2Request,
) (*accountproto.EnableAccountV2Response, error) {
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
	)
	if err != nil {
		return nil, err
	}
	if err := validateEnableAccountV2Request(req); err != nil {
		s.logger.Error(
			"Failed to enable account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	accountV2Pb, err := s.updateAccountV2NoCommandMysql(
		ctx,
		editor,
		req.Email,
		req.OrganizationId,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		wrapperspb.Bool(false),
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to enable account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &accountproto.EnableAccountV2Response{
		Account: accountV2Pb,
	}, nil
}

func (s *AccountService) DisableAccountV2(
	ctx context.Context,
	req *accountproto.DisableAccountV2Request,
) (*accountproto.DisableAccountV2Response, error) {
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDisableAccountV2Request(req); err != nil {
		s.logger.Error(
			"Failed to disable account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	accountV2Pb, err := s.updateAccountV2NoCommandMysql(
		ctx,
		editor,
		req.Email,
		req.OrganizationId,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		wrapperspb.Bool(true),
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to disable account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &accountproto.DisableAccountV2Response{
		Account: accountV2Pb,
	}, nil
}

func (s *AccountService) updateAccountV2MySQL(
	ctx context.Context,
	editor *eventproto.Editor,
	commands []command.Command,
	email, organizationID string,
) (*accountproto.AccountV2, error) {
	var updatedAccountPb *accountproto.AccountV2
	err := s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		account, err := s.accountStorage.GetAccountV2(contextWithTx, email, organizationID)
		if err != nil {
			return err
		}
		handler, err := command.NewAccountV2CommandHandler(editor, account, s.publisher, organizationID)
		if err != nil {
			return err
		}
		for _, c := range commands {
			if err := handler.Handle(ctx, c); err != nil {
				return err
			}
		}
		updatedAccountPb = account.AccountV2
		return s.accountStorage.UpdateAccountV2(contextWithTx, account)
	})
	return updatedAccountPb, err
}

// updateAccountV2NoCommandMysql updates account properties, if the value is nil, it will not be updated.
func (s *AccountService) updateAccountV2NoCommandMysql(
	ctx context.Context,
	editor *eventproto.Editor,
	email, organizationID string,
	name, firstName, lastName, language, avatarImageURL *wrapperspb.StringValue,
	avatar *accountproto.UpdateAccountV2Request_AccountV2Avatar,
	tags *common.StringListValue,
	teamChanges []*accountproto.TeamChange,
	organizationRole *accountproto.UpdateAccountV2Request_OrganizationRoleValue,
	environmentRoles []*accountproto.AccountV2_EnvironmentRole,
	isDisabled *wrapperspb.BoolValue,
) (*accountproto.AccountV2, error) {
	var updatedAccountPb *accountproto.AccountV2
	var updateAccountV2Event *eventproto.Event
	err := s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		account, err := s.accountStorage.GetAccountV2(contextWithTx, email, organizationID)
		if err != nil {
			return err
		}
		updated, err := account.Update(
			name,
			firstName,
			lastName,
			language,
			avatarImageURL,
			avatar,
			tags,
			teamChanges,
			organizationRole,
			environmentRoles,
			isDisabled,
		)
		if err != nil {
			return err
		}
		updateAccountV2Event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ACCOUNT,
			account.Email,
			eventproto.Event_ACCOUNT_V2_UPDATED,
			&eventproto.AccountV2UpdatedEvent{
				Email:          updated.Email,
				OrganizationId: updated.OrganizationId,
			},
			updated,
			account,
		)
		if err != nil {
			return err
		}
		updatedAccountPb = updated.AccountV2
		err = s.accountStorage.UpdateAccountV2(contextWithTx, updated)
		if err != nil {
			return err
		}
		return s.adminAuditLogStorage.CreateAdminAuditLog(
			contextWithTx,
			domainauditlog.NewAuditLog(updateAccountV2Event, storage.AdminEnvironmentID),
		)
	})
	if err != nil {
		return nil, err
	}
	if err = s.publisher.Publish(ctx, updateAccountV2Event); err != nil {
		s.logger.Error(
			"Failed to publish update account event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", organizationID),
				zap.String("email", email),
			)...,
		)
		return nil, err
	}
	return updatedAccountPb, err
}

func (s *AccountService) DeleteAccountV2(
	ctx context.Context,
	req *accountproto.DeleteAccountV2Request,
) (*accountproto.DeleteAccountV2Response, error) {
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAccountV2Request(req); err != nil {
		s.logger.Error(
			"Failed to delete account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	err = s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		account, err := s.accountStorage.GetAccountV2(contextWithTx, req.Email, req.OrganizationId)
		if err != nil {
			return err
		}
		deleteAccountV2Event, err := domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ACCOUNT,
			account.Email,
			eventproto.Event_ACCOUNT_V2_DELETED,
			&eventproto.AccountV2UpdatedEvent{
				Email:          account.Email,
				OrganizationId: account.OrganizationId,
			},
			nil,     // Current state: entity no longer exists
			account, // Previous state: what was deleted
		)
		if err != nil {
			return err
		}
		if err = s.publisher.Publish(ctx, deleteAccountV2Event); err != nil {
			return err
		}
		return s.accountStorage.DeleteAccountV2(contextWithTx, account)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to delete account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &accountproto.DeleteAccountV2Response{}, nil
}

func (s *AccountService) GetAccountV2(
	ctx context.Context,
	req *accountproto.GetAccountV2Request,
) (*accountproto.GetAccountV2Response, error) {
	_, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.OrganizationId,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetAccountV2Request(req); err != nil {
		s.logger.Error(
			"Failed to get account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	account, err := s.getAccountV2(ctx, req.Email, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAccountV2Response{Account: account.AccountV2}, nil
}

func (s *AccountService) getAccountV2(
	ctx context.Context,
	email, organizationID string,
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to get account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", organizationID),
				zap.String("email", email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return account, nil
}

func (s *AccountService) GetAccountV2ByEnvironmentID(
	ctx context.Context,
	req *accountproto.GetAccountV2ByEnvironmentIDRequest,
) (*accountproto.GetAccountV2ByEnvironmentIDResponse, error) {
	_, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.EnvironmentId,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetAccountV2ByEnvironmentIDRequest(req); err != nil {
		s.logger.Error(
			"Failed to get account by environment id",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("EnvironmentId", req.EnvironmentId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	account, err := s.getAccountV2ByEnvironmentID(ctx, req.Email, req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAccountV2ByEnvironmentIDResponse{Account: account.AccountV2}, nil
}

func (s *AccountService) getAccountV2ByEnvironmentID(
	ctx context.Context,
	email, environmentID string,
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2ByEnvironmentID(ctx, email, environmentID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
			return nil, statusAccountNotFound.Err()
		}
		s.logger.Error(
			"Failed to get account by environment id",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentID),
				zap.String("email", email),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return account, nil
}

func (s *AccountService) ListAccountsV2(
	ctx context.Context,
	req *accountproto.ListAccountsV2Request,
) (*accountproto.ListAccountsV2Response, error) {
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.OrganizationId,
	)
	if err != nil {
		return nil, err
	}

	// Users with member role can only view accounts in the environments they have access to
	requestEnvironmentRoles := make([]*accountproto.AccountV2_EnvironmentRole, 0)
	if editor.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
		requestEnvironmentRoles, err = s.constructEnvironmentRoles(req, editor)
		if err != nil {
			return nil, err
		}
	}

	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	if _, err := strconv.Atoi(cursor); err != nil {
		return nil, statusInvalidCursor.Err()
	}

	params := v2as.ListAccountsV2Params{
		OrganizationID:   req.OrganizationId,
		SearchKeyword:    req.SearchKeyword,
		OrderBy:          req.OrderBy,
		OrderDirection:   req.OrderDirection,
		PageSize:         limit,
		Cursor:           cursor,
		EnvironmentRoles: requestEnvironmentRoles,
	}
	if req.Disabled != nil {
		disabled := req.Disabled.Value
		params.Disabled = &disabled
	}
	if len(req.Tags) > 0 {
		params.Tags = req.Tags
	}
	if len(req.Teams) > 0 {
		params.Teams = req.Teams
	}
	if req.OrganizationRole != nil {
		orgRole := req.OrganizationRole.Value
		params.OrganizationRole = &orgRole
	}
	// Only set EnvironmentID/EnvironmentRole when user is not a member with restricted roles
	if len(requestEnvironmentRoles) == 0 {
		if req.EnvironmentId != nil {
			params.EnvironmentID = &req.EnvironmentId.Value
		}
		if req.EnvironmentRole != nil {
			params.EnvironmentRole = &req.EnvironmentRole.Value
		}
	}

	accounts, nextCursor, totalCount, err := s.accountStorage.ListAccountsV2(ctx, params)
	if err != nil {
		if errors.Is(err, v2as.ErrInvalidOrderBy) {
			return nil, statusInvalidOrderBy.Err()
		}
		s.logger.Error(
			"Failed to list accounts",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &accountproto.ListAccountsV2Response{
		Accounts:   accounts,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *AccountService) constructEnvironmentRoles(
	req *accountproto.ListAccountsV2Request,
	editor *eventproto.Editor,
) ([]*accountproto.AccountV2_EnvironmentRole, error) {
	requestEnvironmentRoles := make([]*accountproto.AccountV2_EnvironmentRole, 0)
	// No allowed roles means the user has no access to any environment in the organization
	if len(editor.EnvironmentRoles) == 0 {
		return nil, nil
	}

	if req.EnvironmentId != nil && req.EnvironmentRole != nil {
		for _, role := range editor.EnvironmentRoles {
			if role.EnvironmentId == req.EnvironmentId.Value &&
				role.Role == accountproto.AccountV2_Role_Environment(req.EnvironmentRole.Value) {
				requestEnvironmentRoles = append(requestEnvironmentRoles, role)
				break
			}
		}
	} else if req.EnvironmentId != nil && req.EnvironmentRole == nil {
		for _, role := range editor.EnvironmentRoles {
			if role.EnvironmentId == req.EnvironmentId.Value {
				requestEnvironmentRoles = append(requestEnvironmentRoles, role)
				break
			}
		}
	} else {
		requestEnvironmentRoles = append(requestEnvironmentRoles, editor.EnvironmentRoles...)
	}
	return requestEnvironmentRoles, nil
}

// initiatePasswordSetupForNewAccount sends password setup email to newly created accounts
func (s *AccountService) initiatePasswordSetupForNewAccount(ctx context.Context, email string) {
	if s.authClient == nil {
		s.logger.Debug("Auth client not available, skipping password setup", zap.String("email", email))
		return
	}

	_, err := s.authClient.InitiatePasswordSetup(ctx, &authproto.InitiatePasswordSetupRequest{
		Email: email,
	})
	if err != nil {
		s.logger.Warn("Failed to initiate password setup for new account",
			zap.Error(err),
			zap.String("email", email),
		)
	} else {
		s.logger.Info("Password setup initiated for new account", zap.String("email", email))
	}
}
