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
	editor, err := s.checkRole(ctx, accountproto.Account_OWNER, req.EnvironmentNamespace, localizer)
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
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		accountStorage := v2as.NewAccountStorage(tx)
		handler := command.NewAccountCommandHandler(editor, account, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return accountStorage.CreateAccount(ctx, account, req.EnvironmentNamespace)
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
	editor, err := s.checkRole(ctx, accountproto.Account_OWNER, req.EnvironmentNamespace, localizer)
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
	if err := s.updateAccountMySQL(ctx, editor, req.Command, req.Id, req.EnvironmentNamespace); err != nil {
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
	editor, err := s.checkRole(ctx, accountproto.Account_OWNER, req.EnvironmentNamespace, localizer)
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
	if err := s.updateAccountMySQL(ctx, editor, req.Command, req.Id, req.EnvironmentNamespace); err != nil {
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
	editor, err := s.checkRole(ctx, accountproto.Account_OWNER, req.EnvironmentNamespace, localizer)
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
	if err := s.updateAccountMySQL(ctx, editor, req.Command, req.Id, req.EnvironmentNamespace); err != nil {
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

func (s *AccountService) updateAccountMySQL(
	ctx context.Context,
	editor *eventproto.Editor,
	cmd command.Command,
	id, environmentNamespace string,
) error {
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return err
	}
	return s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		accountStorage := v2as.NewAccountStorage(tx)
		account, err := accountStorage.GetAccount(ctx, id, environmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewAccountCommandHandler(editor, account, s.publisher, environmentNamespace)
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return accountStorage.UpdateAccount(ctx, account, environmentNamespace)
	})
}

func (s *AccountService) GetAccount(
	ctx context.Context,
	req *accountproto.GetAccountRequest,
) (*accountproto.GetAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
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
	accountStorage := v2as.NewAccountStorage(s.mysqlClient)
	account, err := accountStorage.GetAccount(ctx, email, environmentNamespace)
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
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
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
	accountStorage := v2as.NewAccountStorage(s.mysqlClient)
	accounts, nextCursor, totalCount, err := accountStorage.ListAccounts(
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

func (s *AccountService) CreateAccountV2(
	ctx context.Context,
	req *accountproto.CreateAccountV2Request,
) (*accountproto.CreateAccountV2Response, error) {
	// TODO: Implement the business logic.
	return &accountproto.CreateAccountV2Response{}, nil
}

func (s *AccountService) UpdateAccountV2(
	ctx context.Context,
	req *accountproto.UpdateAccountV2Request,
) (*accountproto.UpdateAccountV2Response, error) {
	// TODO: Implement the business logic.
	return &accountproto.UpdateAccountV2Response{}, nil
}

func (s *AccountService) EnableAccountV2(
	ctx context.Context,
	req *accountproto.EnableAccountV2Request,
) (*accountproto.EnableAccountV2Response, error) {
	// TODO: Implement the business logic.
	return &accountproto.EnableAccountV2Response{}, nil
}

func (s *AccountService) DisableAccountV2(
	ctx context.Context,
	req *accountproto.DisableAccountV2Request,
) (*accountproto.DisableAccountV2Response, error) {
	// TODO: Implement the business logic.
	return &accountproto.DisableAccountV2Response{}, nil
}

func (s *AccountService) DeleteAccountV2(
	ctx context.Context,
	req *accountproto.DeleteAccountV2Request,
) (*accountproto.DeleteAccountV2Response, error) {
	// TODO: Implement the business logic.
	return &accountproto.DeleteAccountV2Response{}, nil
}

func (s *AccountService) GetAccountV2(
	ctx context.Context,
	req *accountproto.GetAccountV2Request,
) (*accountproto.GetAccountV2Response, error) {
	// TODO: Implement the business logic.
	return &accountproto.GetAccountV2Response{}, nil
}

func (s *AccountService) ListAccountsV2(
	ctx context.Context,
	req *accountproto.ListAccountsV2Request,
) (*accountproto.ListAccountsV2Response, error) {
	// TODO: Implement the business logic.
	return &accountproto.ListAccountsV2Response{}, nil
}
