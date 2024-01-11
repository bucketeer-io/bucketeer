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
	"context"
	"errors"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func (s *AccountService) CreateAccount(
	ctx context.Context,
	req *accountproto.CreateAccountRequest,
) (*accountproto.CreateAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentNamespace,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateCreateAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	account, err := domain.NewAccount(req.Command.Email, req.Command.Role)
	if err != nil {
		s.logger.Error(
			"Failed to create a new account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	orgID, err := s.getOrganizationID(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get organization id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	// check if an Admin Account that has the same email already exists
	_, err = s.getAdminAccount(ctx, account.Id, localizer)
	if status.Code(err) != codes.NotFound {
		if err == nil {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
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
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		handler := command.NewAccountCommandHandler(editor, account, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		err := s.accountStorage.CreateAccount(ctx, account, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		// TODO: temporary implementation: double write account v2
		accountV2 := domain.ConvertAccountV2(account, req.EnvironmentNamespace, orgID)
		return s.accountStorage.CreateAccountV2(ctx, accountV2)
	})
	if err != nil {
		if err == v2as.ErrAccountAlreadyExists {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &accountproto.CreateAccountResponse{}, nil
}

func (s *AccountService) ChangeAccountRole(
	ctx context.Context,
	req *accountproto.ChangeAccountRoleRequest,
) (*accountproto.ChangeAccountRoleResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentNamespace,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateChangeAccountRoleRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to change account role",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	orgID, err := s.getOrganizationID(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get organization id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		account, err := s.accountStorage.GetAccount(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewAccountCommandHandler(editor, account, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		err = s.accountStorage.UpdateAccount(ctx, account, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		// TODO: temporary implementation: double write account v2
		accountV2, err := s.accountStorage.GetAccountV2(ctx, req.Id, orgID)
		if err != nil {
			return err
		}
		accountV2.PatchAccountV2EnvironmentRoles(req.EnvironmentNamespace, req.Command.Role)
		return s.accountStorage.UpdateAccountV2(ctx, accountV2)
	})
	if err != nil {
		if err == v2as.ErrAccountNotFound || err == v2as.ErrAccountUnexpectedAffectedRows {
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
			"Failed to change account role",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &accountproto.ChangeAccountRoleResponse{}, nil
}

func (s *AccountService) EnableAccount(
	ctx context.Context,
	req *accountproto.EnableAccountRequest,
) (*accountproto.EnableAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentNamespace,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateEnableAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to enable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	orgID, err := s.getOrganizationID(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get organization id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		accountV1, err := s.accountStorage.GetAccount(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewAccountCommandHandler(editor, accountV1, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		err = s.accountStorage.UpdateAccount(ctx, accountV1, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		// TODO: temporary implementation: double write account v2
		accountV2, err := s.accountStorage.GetAccountV2(ctx, req.Id, orgID)
		if err != nil {
			return err
		}
		accountV2.PatchAccountV2EnvironmentRoles(req.EnvironmentNamespace, accountV1.Role)
		return s.accountStorage.UpdateAccountV2(ctx, accountV2)
	})
	if err != nil {
		if err == v2as.ErrAccountNotFound || err == v2as.ErrAccountUnexpectedAffectedRows {
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
			"Failed to enable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &accountproto.EnableAccountResponse{}, nil
}

func (s *AccountService) DisableAccount(
	ctx context.Context,
	req *accountproto.DisableAccountRequest,
) (*accountproto.DisableAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentNamespace,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDisableAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to disable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	orgID, err := s.getOrganizationID(ctx, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get organization id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		accountV1, err := s.accountStorage.GetAccount(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewAccountCommandHandler(editor, accountV1, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		err = s.accountStorage.UpdateAccount(ctx, accountV1, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		// TODO: temporary implementation: double write account v2
		accountV2, err := s.accountStorage.GetAccountV2(ctx, req.Id, orgID)
		if err != nil {
			return err
		}
		accountV2.RemoveAccountV2EnvironmentRole(req.EnvironmentNamespace)
		return s.accountStorage.UpdateAccountV2(ctx, accountV2)
	})
	if err != nil {
		if err == v2as.ErrAccountNotFound || err == v2as.ErrAccountUnexpectedAffectedRows {
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
			"Failed to disable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &accountproto.DisableAccountResponse{}, nil
}

func (s *AccountService) GetAccount(
	ctx context.Context,
	req *accountproto.GetAccountRequest,
) (*accountproto.GetAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	account, err := s.getAccount(ctx, req.Email, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAccountResponse{Account: account.Account}, nil
}

func (s *AccountService) getAccount(
	ctx context.Context,
	email, environmentNamespace string,
	localizer locale.Localizer,
) (*domain.Account, error) {
	account, err := s.accountStorage.GetAccount(ctx, email, environmentNamespace)
	if err != nil {
		if err == v2as.ErrAccountNotFound {
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
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
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

func (s *AccountService) ListAccounts(
	ctx context.Context,
	req *accountproto.ListAccountsRequest,
) (*accountproto.ListAccountsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("disabled", "=", req.Disabled.Value))
	}
	if req.Role != nil {
		whereParts = append(whereParts, mysql.NewFilter("role", "=", req.Role.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"email"}, req.SearchKeyword))
	}
	orders, err := s.newAccountListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	accounts, nextCursor, totalCount, err := s.accountStorage.ListAccounts(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list accounts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &accountproto.ListAccountsResponse{
		Accounts:   accounts,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *AccountService) newAccountListOrders(
	orderBy accountproto.ListAccountsRequest_OrderBy,
	orderDirection accountproto.ListAccountsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case accountproto.ListAccountsRequest_DEFAULT,
		accountproto.ListAccountsRequest_EMAIL:
		column = "email"
	case accountproto.ListAccountsRequest_CREATED_AT:
		column = "created_at"
	case accountproto.ListAccountsRequest_UPDATED_AT:
		column = "updated_at"
	default:
		dt, err := statusInvalidOrderBy.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == accountproto.ListAccountsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *AccountService) getOrganizationID(
	ctx context.Context,
	environmentNamespace string,
) (string, error) {
	environments, err := s.listEnvironments(ctx)
	if err != nil {
		return "", err
	}
	for _, e := range environments {
		if e.Id == environmentNamespace {
			return e.OrganizationId, nil
		}
	}
	return "", statusEnvironmentNotFound.Err()
}

func (s *AccountService) CreateAccountV2(
	ctx context.Context,
	req *accountproto.CreateAccountV2Request,
) (*accountproto.CreateAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateCreateAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, err
	}
	account := domain.NewAccountV2(
		req.Command.Email, req.Command.Name, req.Command.AvatarImageUrl, req.OrganizationId,
		req.Command.OrganizationRole, req.Command.EnvironmentRoles,
	)
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		handler := command.NewAccountV2CommandHandler(editor, account, s.publisher, req.OrganizationId)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.accountStorage.CreateAccountV2(ctx, account)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAccountAlreadyExists) {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
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
	return &accountproto.CreateAccountV2Response{Account: account.AccountV2}, nil
}

func (s *AccountService) UpdateAccountV2(
	ctx context.Context,
	req *accountproto.UpdateAccountV2Request,
) (*accountproto.UpdateAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	commands := s.getUpdateAccountV2Commands(req)
	if err := validateUpdateAccountV2Request(req, commands, localizer); err != nil {
		s.logger.Error(
			"Failed to update account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	if err := s.updateAccountV2MySQL(ctx, editor, commands, req.Email, req.OrganizationId); err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to update account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	return &accountproto.UpdateAccountV2Response{}, nil
}

func (s *AccountService) getUpdateAccountV2Commands(req *accountproto.UpdateAccountV2Request) []command.Command {
	commands := make([]command.Command, 0)
	if req.ChangeNameCommand != nil {
		commands = append(commands, req.ChangeNameCommand)
	}
	if req.ChangeAvatarUrlCommand != nil {
		commands = append(commands, req.ChangeAvatarUrlCommand)
	}
	if req.ChangeOrganizationRoleCommand != nil {
		commands = append(commands, req.ChangeOrganizationRoleCommand)
	}
	if req.ChangeEnvironmentRolesCommand != nil {
		commands = append(commands, req.ChangeEnvironmentRolesCommand)
	}
	return commands
}

func (s *AccountService) EnableAccountV2(
	ctx context.Context,
	req *accountproto.EnableAccountV2Request,
) (*accountproto.EnableAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateEnableAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to enable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	err = s.updateAccountV2MySQL(
		ctx,
		editor,
		[]command.Command{req.Command},
		req.Email,
		req.OrganizationId,
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to enable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	return &accountproto.EnableAccountV2Response{}, nil
}

func (s *AccountService) DisableAccountV2(
	ctx context.Context,
	req *accountproto.DisableAccountV2Request,
) (*accountproto.DisableAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDisableAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to disable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	err = s.updateAccountV2MySQL(
		ctx,
		editor,
		[]command.Command{req.Command},
		req.Email,
		req.OrganizationId,
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to disable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	return &accountproto.DisableAccountV2Response{}, nil
}

func (s *AccountService) updateAccountV2MySQL(
	ctx context.Context,
	editor *eventproto.Editor,
	commands []command.Command,
	email, organizationID string,
) error {
	return s.accountStorage.RunInTransaction(ctx, func() error {
		account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
		if err != nil {
			return err
		}
		handler := command.NewAccountV2CommandHandler(editor, account, s.publisher, organizationID)
		for _, c := range commands {
			if err := handler.Handle(ctx, c); err != nil {
				return err
			}
		}
		return s.accountStorage.UpdateAccountV2(ctx, account)
	})
}

func (s *AccountService) DeleteAccountV2(
	ctx context.Context,
	req *accountproto.DeleteAccountV2Request,
) (*accountproto.DeleteAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to delete account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		account, err := s.accountStorage.GetAccountV2(ctx, req.Email, req.OrganizationId)
		if err != nil {
			return err
		}
		handler := command.NewAccountV2CommandHandler(editor, account, s.publisher, req.OrganizationId)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.accountStorage.DeleteAccountV2(ctx, account)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to delete account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	return &accountproto.DeleteAccountV2Response{}, nil
}

func (s *AccountService) GetAccountV2(
	ctx context.Context,
	req *accountproto.GetAccountV2Request,
) (*accountproto.GetAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	account, err := s.getAccountV2(ctx, req.Email, req.OrganizationId, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAccountV2Response{Account: account.AccountV2}, nil
}

func (s *AccountService) getAccountV2(
	ctx context.Context,
	email, organizationID string,
	localizer locale.Localizer,
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
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
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", organizationID),
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

func (s *AccountService) GetAccountV2ByEnvironmentID(
	ctx context.Context,
	req *accountproto.GetAccountV2ByEnvironmentIDRequest,
) (*accountproto.GetAccountV2ByEnvironmentIDResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetAccountV2ByEnvironmentIDRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to get account by environment id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("EnvironmentId", req.EnvironmentId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	account, err := s.getAccountV2ByEnvironmentID(ctx, req.Email, req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAccountV2ByEnvironmentIDResponse{Account: account.AccountV2}, nil
}

func (s *AccountService) getAccountV2ByEnvironmentID(
	ctx context.Context,
	email, environmentID string,
	localizer locale.Localizer,
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2ByEnvironmentID(ctx, email, environmentID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
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
			"Failed to get account by environment id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentID),
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

func (s *AccountService) ListAccountsV2(
	ctx context.Context,
	req *accountproto.ListAccountsV2Request,
) (*accountproto.ListAccountsV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("organization_id", "=", req.OrganizationId),
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("disabled", "=", req.Disabled.Value))
	}
	if req.Role != nil {
		whereParts = append(whereParts, mysql.NewFilter("role", "=", req.Role.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"email"}, req.SearchKeyword))
	}
	orders, err := s.newAccountV2ListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	accounts, nextCursor, totalCount, err := s.accountStorage.ListAccountsV2(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list accounts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
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
	return &accountproto.ListAccountsV2Response{
		Accounts:   accounts,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *AccountService) newAccountV2ListOrders(
	orderBy accountproto.ListAccountsV2Request_OrderBy,
	orderDirection accountproto.ListAccountsV2Request_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case accountproto.ListAccountsV2Request_DEFAULT,
		accountproto.ListAccountsV2Request_EMAIL:
		column = "email"
	case accountproto.ListAccountsV2Request_CREATED_AT:
		column = "created_at"
	case accountproto.ListAccountsV2Request_UPDATED_AT:
		column = "updated_at"
	default:
		dt, err := statusInvalidOrderBy.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == accountproto.ListAccountsV2Request_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}
