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
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

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
		// TODO: temporary implementation: double write account v2 ---
		exist, err := s.accountStorage.GetAccountV2(ctx, account.Email, req.OrganizationId)
		if err != nil && !errors.Is(err, v2as.ErrAccountNotFound) {
			return err
		}
		if exist != nil {
			handler := command.NewAccountV2CommandHandler(editor, exist, s.publisher, req.OrganizationId)
			cmd := &accountproto.ChangeAccountV2EnvironmentRolesCommand{
				Roles:     account.EnvironmentRoles,
				WriteType: accountproto.ChangeAccountV2EnvironmentRolesCommand_WriteType_PATCH,
			}
			if err := handler.Handle(ctx, cmd); err != nil {
				return err
			}
			return s.accountStorage.UpdateAccountV2(ctx, exist)
		}
		// TODO: temporary implementation end ---
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
	if req.OrganizationRole != nil {
		whereParts = append(whereParts, mysql.NewFilter("organization_role", "=", req.OrganizationRole.Value))
	}
	if req.EnvironmentId != nil && req.EnvironmentRole != nil {
		values := make([]interface{}, 1)
		values[0] = fmt.Sprintf("{\"environment_id\": \"%s\", \"role\": %d}", req.EnvironmentId.Value, req.EnvironmentRole.Value) // nolint:lll
		whereParts = append(whereParts, mysql.NewJSONFilter("environment_roles", mysql.JSONContainsJSON, values))
	} else if req.EnvironmentId != nil {
		values := make([]interface{}, 1)
		values[0] = fmt.Sprintf("{\"environment_id\": \"%s\"}", req.EnvironmentId.Value)
		whereParts = append(whereParts, mysql.NewJSONFilter("environment_roles", mysql.JSONContainsJSON, values))
	} else if req.EnvironmentRole != nil {
		values := make([]interface{}, 1)
		values[0] = fmt.Sprintf("{\"role\": %d}", req.EnvironmentRole.Value)
		whereParts = append(whereParts, mysql.NewJSONFilter("environment_roles", mysql.JSONContainsJSON, values))
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
