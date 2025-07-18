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

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/api/api"

	accounterr "github.com/bucketeer-io/bucketeer/pkg/account"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func (s *AccountService) CreateSearchFilter(
	ctx context.Context,
	req *accountproto.CreateSearchFilterRequest,
) (*accountproto.CreateSearchFilterResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}

	if err := validateCreateSearchFilterRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to validate request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("request", req),
			)...,
		)
		return nil, err
	}

	// If the target account is a system admin, we must use the system admin organization ID
	// Otherwise, it will return a not found error
	// because the account doesn't exist in non-system admin organizations.
	sysAdminAccount, err := s.getSystemAdminAccountV2(ctx, req.Email, localizer)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	orgID := req.OrganizationId
	if sysAdminAccount != nil {
		orgID = sysAdminAccount.OrganizationId
	}

	if _, err := s.updateAccountV2MySQL(
		ctx,
		editor,
		[]command.Command{req.Command},
		req.Email,
		orgID); err != nil {
		s.logger.Error(
			"Failed to create search filter",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", orgID),
				zap.String("email", req.Email),
				zap.String("environmentID", req.EnvironmentId),
				zap.String("searchFilterName", req.Command.Name),
				zap.String("query", req.Command.Query),
				zap.String("filterTargetType", req.Command.FilterTargetType.String()),
			)...,
		)
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
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}

	return &accountproto.CreateSearchFilterResponse{}, nil
}

func (s *AccountService) UpdateSearchFilter(
	ctx context.Context,
	req *accountproto.UpdateSearchFilterRequest,
) (*accountproto.UpdateSearchFilterResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	commands := s.getUpdateSearchFilterCommands(req)

	if err := validateUpdateSearchFilterRequest(req, commands, localizer); err != nil {
		s.logger.Error(
			"Failed to validate request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("request", req),
			)...,
		)
		return nil, err
	}

	// If the target account is a system admin, we must use the system admin organization ID
	// Otherwise, it will return a not found error
	// because the account doesn't exist in non-system admin organizations.
	sysAdminAccount, err := s.getSystemAdminAccountV2(ctx, req.Email, localizer)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	orgID := req.OrganizationId
	if sysAdminAccount != nil {
		orgID = sysAdminAccount.OrganizationId
	}

	if _, err := s.updateAccountV2MySQL(
		ctx,
		editor,
		commands,
		req.Email,
		orgID); err != nil {
		s.logger.Error(
			"Failed to update search filter",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", orgID),
				zap.String("email", req.Email),
			)...,
		)
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		} else if errors.Is(err, accounterr.ErrSearchFilterNotFound) {
			dt, err := statusSearchFilterIDNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &accountproto.UpdateSearchFilterResponse{}, nil
}

func (s *AccountService) getUpdateSearchFilterCommands(req *accountproto.UpdateSearchFilterRequest) []command.Command {
	commands := make([]command.Command, 0)
	if req.ChangeNameCommand != nil {
		commands = append(commands, req.ChangeNameCommand)
	}
	if req.ChangeQueryCommand != nil {
		commands = append(commands, req.ChangeQueryCommand)
	}
	if req.ChangeDefaultFilterCommand != nil {
		commands = append(commands, req.ChangeDefaultFilterCommand)
	}
	return commands
}

func (s *AccountService) DeleteSearchFilter(
	ctx context.Context,
	req *accountproto.DeleteSearchFilterRequest,
) (*accountproto.DeleteSearchFilterResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}

	if err := validateDeleteSearchFilterRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to validate request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("request", req),
			)...,
		)
		return nil, err
	}

	// If the target account is a system admin, we must use the system admin organization ID
	// Otherwise, it will return a not found error
	// because the account doesn't exist in non-system admin organizations.
	sysAdminAccount, err := s.getSystemAdminAccountV2(ctx, req.Email, localizer)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	orgID := req.OrganizationId
	if sysAdminAccount != nil {
		orgID = sysAdminAccount.OrganizationId
	}

	if _, err := s.updateAccountV2MySQL(
		ctx,
		editor,
		[]command.Command{req.Command},
		req.Email,
		orgID); err != nil {
		s.logger.Error(
			"Failed to delete search filter",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", orgID),
				zap.String("email", req.Email),
				zap.String("searchFilterID", req.Command.Id),
			)...,
		)
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if errors.Is(err, accounterr.ErrSearchFilterNotFound) {
			dt, err := statusSearchFilterIDNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}

	return &accountproto.DeleteSearchFilterResponse{}, nil
}
