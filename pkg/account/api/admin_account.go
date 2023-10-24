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
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func (s *AccountService) GetMeV2(
	ctx context.Context,
	req *accountproto.GetMeV2Request,
) (*accountproto.GetMeV2Response, error) {
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
	return s.getMeV2(ctx, t.Email, localizer)
}

func (s *AccountService) GetMeByEmailV2(
	ctx context.Context,
	req *accountproto.GetMeByEmailV2Request,
) (*accountproto.GetMeV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkAdminRole(ctx, localizer)
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
	return s.getMeV2(ctx, req.Email, localizer)
}

func (s *AccountService) getMeV2(
	ctx context.Context,
	email string,
	localizer locale.Localizer,
) (*accountproto.GetMeV2Response, error) {
	projects, err := s.listProjects(ctx)
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
	if len(projects) == 0 {
		s.logger.Error(
			"Could not find any projects",
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
	environments, err := s.listEnvironments(ctx)
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
	if len(environments) == 0 {
		s.logger.Error(
			"Could not find any environments",
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
	// admin account response
	adminAccount, err := s.getAdminAccount(ctx, email, localizer)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	if adminAccount != nil && !adminAccount.Disabled && !adminAccount.Deleted {
		environmentRoles, err := s.makeAdminEnvironmentRolesV2(
			projects,
			environments,
			accountproto.Account_OWNER,
			localizer,
		)
		if err != nil {
			return nil, err
		}
		return &accountproto.GetMeV2Response{
			Email:            adminAccount.Email,
			IsAdmin:          true,
			EnvironmentRoles: environmentRoles,
		}, nil
	}
	// environment account response
	environmentRoles, err := s.makeEnvironmentRolesV2(ctx, email, projects, environments, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetMeV2Response{
		Email:            email,
		IsAdmin:          false,
		EnvironmentRoles: environmentRoles,
	}, nil
}

func (s *AccountService) makeAdminEnvironmentRolesV2(
	projects []*environmentproto.Project,
	environments []*environmentproto.EnvironmentV2,
	adminRole accountproto.Account_Role,
	localizer locale.Localizer,
) ([]*accountproto.EnvironmentRoleV2, error) {
	projectSet := s.makeProjectSet(projects)
	environmentRoles := make([]*accountproto.EnvironmentRoleV2, 0)
	for _, e := range environments {
		p, ok := projectSet[e.ProjectId]
		if !ok || p.Disabled {
			continue
		}
		er := &accountproto.EnvironmentRoleV2{Environment: e, Role: adminRole}
		if p.Trial {
			er.TrialProject = true
			er.TrialStartedAt = p.CreatedAt
		}
		environmentRoles = append(environmentRoles, er)
	}
	if len(environmentRoles) == 0 {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return environmentRoles, nil
}

func (s *AccountService) makeEnvironmentRolesV2(
	ctx context.Context,
	email string,
	projects []*environmentproto.Project,
	environments []*environmentproto.EnvironmentV2,
	localizer locale.Localizer,
) ([]*accountproto.EnvironmentRoleV2, error) {
	projectSet := s.makeProjectSet(projects)
	environmentRoles := make([]*accountproto.EnvironmentRoleV2, 0, len(environments))
	for _, e := range environments {
		p, ok := projectSet[e.ProjectId]
		if !ok || p.Disabled {
			continue
		}
		account, err := s.getAccount(ctx, email, e.Id, localizer)
		if err != nil && status.Code(err) != codes.NotFound {
			return nil, err
		}
		if account == nil || account.Disabled || account.Deleted {
			continue
		}
		er := &accountproto.EnvironmentRoleV2{Environment: e, Role: account.Role}
		if p.Trial {
			er.TrialProject = true
			er.TrialStartedAt = p.CreatedAt
		}
		environmentRoles = append(environmentRoles, er)
	}
	if len(environmentRoles) == 0 {
		dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.NotFoundError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return environmentRoles, nil
}

func (s *AccountService) CreateAdminAccount(
	ctx context.Context,
	req *accountproto.CreateAdminAccountRequest,
) (*accountproto.CreateAdminAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateCreateAdminAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to create admin account",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	account, err := domain.NewAccount(req.Command.Email, accountproto.Account_OWNER)
	if err != nil {
		s.logger.Error(
			"Failed to create a new admin account",
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
	environments, err := s.listEnvironments(ctx)
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
	// check if an Account that has the same email already exists in any environment
	accountStorage := v2as.NewAccountStorage(s.mysqlClient)
	for _, env := range environments {
		_, err := accountStorage.GetAccount(ctx, account.Id, env.Id)
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
		if err != v2as.ErrAccountNotFound {
			return nil, err
		}
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
		adminAccountStorage := v2as.NewAdminAccountStorage(tx)
		handler := command.NewAdminAccountCommandHandler(editor, account, s.publisher)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return adminAccountStorage.CreateAdminAccount(ctx, account)
	})
	if err != nil {
		if err == v2as.ErrAdminAccountAlreadyExists {
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
			"Failed to create admin account",
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
	return &accountproto.CreateAdminAccountResponse{}, nil
}

func (s *AccountService) EnableAdminAccount(
	ctx context.Context,
	req *accountproto.EnableAdminAccountRequest,
) (*accountproto.EnableAdminAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateEnableAdminAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to enable admin account",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	if err := s.updateAdminAccountMySQL(ctx, editor, req.Id, req.Command); err != nil {
		if err == v2as.ErrAdminAccountNotFound || err == v2as.ErrAdminAccountUnexpectedAffectedRows {
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
			"Failed to enable admin account",
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
	return &accountproto.EnableAdminAccountResponse{}, nil
}

func (s *AccountService) DisableAdminAccount(
	ctx context.Context,
	req *accountproto.DisableAdminAccountRequest,
) (*accountproto.DisableAdminAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDisableAdminAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to disable admin account",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	if err := s.updateAdminAccountMySQL(ctx, editor, req.Id, req.Command); err != nil {
		if err == v2as.ErrAdminAccountNotFound || err == v2as.ErrAdminAccountUnexpectedAffectedRows {
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
			"Failed to disable admin account",
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
	return &accountproto.DisableAdminAccountResponse{}, nil
}

func (s *AccountService) updateAdminAccountMySQL(
	ctx context.Context,
	editor *eventproto.Editor,
	id string,
	cmd command.Command,
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
		adminAccountStorage := v2as.NewAdminAccountStorage(tx)
		account, err := adminAccountStorage.GetAdminAccount(ctx, id)
		if err != nil {
			return err
		}
		handler := command.NewAdminAccountCommandHandler(editor, account, s.publisher)
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return adminAccountStorage.UpdateAdminAccount(ctx, account)
	})
}

func (s *AccountService) ConvertAccount(
	ctx context.Context,
	req *accountproto.ConvertAccountRequest,
) (*accountproto.ConvertAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateConvertAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	account, err := domain.NewAccount(req.Id, accountproto.Account_OWNER)
	if err != nil {
		s.logger.Error(
			"Failed to create a new admin account",
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
	environments, err := s.listEnvironments(ctx)
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
	deleteAccountCommand := &accountproto.DeleteAccountCommand{}
	createAdminAccountCommand := &accountproto.CreateAdminAccountCommand{Email: req.Id}
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
		var existedAccountCount int
		for _, env := range environments {
			existedAccount, err := accountStorage.GetAccount(ctx, account.Id, env.Id)
			if err != nil {
				if err == v2as.ErrAccountNotFound {
					continue
				}
				return err
			}
			existedAccountCount++
			handler := command.NewAccountCommandHandler(
				editor,
				existedAccount,
				s.publisher,
				env.Id,
			)
			if err := handler.Handle(ctx, deleteAccountCommand); err != nil {
				return err
			}
			if err := accountStorage.UpdateAccount(ctx, existedAccount, env.Id); err != nil {
				return err
			}
		}
		if existedAccountCount == 0 {
			return v2as.ErrAccountNotFound
		}
		adminAccountStorage := v2as.NewAdminAccountStorage(tx)
		handler := command.NewAdminAccountCommandHandler(editor, account, s.publisher)
		if err := handler.Handle(ctx, createAdminAccountCommand); err != nil {
			return err
		}
		return adminAccountStorage.CreateAdminAccount(ctx, account)
	})
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
		if err == v2as.ErrAdminAccountAlreadyExists {
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
			"Failed to convert account",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	return &accountproto.ConvertAccountResponse{}, nil
}

func (s *AccountService) GetAdminAccount(
	ctx context.Context,
	req *accountproto.GetAdminAccountRequest,
) (*accountproto.GetAdminAccountResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetAdminAccountRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to get admin account",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	account, err := s.getAdminAccount(ctx, req.Email, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAdminAccountResponse{Account: account.Account}, nil
}

func (s *AccountService) getAdminAccount(
	ctx context.Context,
	email string,
	localizer locale.Localizer,
) (*domain.Account, error) {
	adminAccountStorage := v2as.NewAdminAccountStorage(s.mysqlClient)
	account, err := adminAccountStorage.GetAdminAccount(ctx, email)
	if err != nil {
		if err == v2as.ErrAdminAccountNotFound {
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
			"Failed to get admin account",
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

func (s *AccountService) ListAdminAccounts(
	ctx context.Context,
	req *accountproto.ListAdminAccountsRequest,
) (*accountproto.ListAdminAccountsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{mysql.NewFilter("deleted", "=", false)}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("disabled", "=", req.Disabled.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"email"}, req.SearchKeyword))
	}
	orders, err := s.newAdminAccountListOrders(req.OrderBy, req.OrderDirection, localizer)
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
	adminAccountStorage := v2as.NewAdminAccountStorage(s.mysqlClient)
	accounts, nextCursor, totalCount, err := adminAccountStorage.ListAdminAccounts(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list admin accounts",
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
	return &accountproto.ListAdminAccountsResponse{
		Accounts:   accounts,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *AccountService) newAdminAccountListOrders(
	orderBy accountproto.ListAdminAccountsRequest_OrderBy,
	orderDirection accountproto.ListAdminAccountsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case accountproto.ListAdminAccountsRequest_DEFAULT,
		accountproto.ListAdminAccountsRequest_EMAIL:
		column = "email"
	case accountproto.ListAdminAccountsRequest_CREATED_AT:
		column = "created_at"
	case accountproto.ListAdminAccountsRequest_UPDATED_AT:
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
	if orderDirection == accountproto.ListAdminAccountsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}
