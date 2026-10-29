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

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

func (s *AccountService) CreateSearchFilter(
	ctx context.Context,
	req *accountproto.CreateSearchFilterRequest,
) (*accountproto.CreateSearchFilterResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if err := validateCreateSearchFilterRequest(req); err != nil {
		s.logger.Error(
			"Failed to validate request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("request", req),
			)...,
		)
		return nil, err
	}

	if err := s.checkSearchFilterOwner(ctx, editor, req.Email, req.EnvironmentId); err != nil {
		return nil, err
	}

	if _, err := s.updateAccountV2MySQL(
		ctx,
		editor,
		[]command.Command{req.Command},
		req.EnvironmentId); err != nil {
		s.logger.Error(
			"Failed to create search filter",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", editor.Email),
				zap.String("environmentID", req.EnvironmentId),
				zap.String("searchFilterName", req.Command.Name),
				zap.String("query", req.Command.Query),
				zap.String("filterTargetType", req.Command.FilterTargetType.String()),
			)...,
		)
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			return nil, statusAccountNotFound.Err()
		}
		return nil, statusInternal.Err()
	}

	return &accountproto.CreateSearchFilterResponse{}, nil
}

func (s *AccountService) UpdateSearchFilter(
	ctx context.Context,
	req *accountproto.UpdateSearchFilterRequest,
) (*accountproto.UpdateSearchFilterResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	commands := s.getUpdateSearchFilterCommands(req)

	if err := validateUpdateSearchFilterRequest(req, commands); err != nil {
		s.logger.Error(
			"Failed to validate request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("request", req),
			)...,
		)
		return nil, err
	}

	if err := s.checkSearchFilterOwner(ctx, editor, req.Email, req.EnvironmentId); err != nil {
		return nil, err
	}

	if _, err := s.updateAccountV2MySQL(
		ctx,
		editor,
		commands,
		req.EnvironmentId); err != nil {
		s.logger.Error(
			"Failed to update search filter",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", editor.Email),
				zap.String("environmentID", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			return nil, statusAccountNotFound.Err()
		} else if errors.Is(err, domain.ErrSearchFilterNotFound) {
			return nil, statusSearchFilterIDNotFound.Err()
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
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if err := validateDeleteSearchFilterRequest(req); err != nil {
		s.logger.Error(
			"Failed to validate request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("request", req),
			)...,
		)
		return nil, err
	}

	if err := s.checkSearchFilterOwner(ctx, editor, req.Email, req.EnvironmentId); err != nil {
		return nil, err
	}

	if _, err := s.updateAccountV2MySQL(
		ctx,
		editor,
		[]command.Command{req.Command},
		req.EnvironmentId); err != nil {
		s.logger.Error(
			"Failed to delete search filter",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", editor.Email),
				zap.String("environmentID", req.EnvironmentId),
				zap.String("searchFilterID", req.Command.Id),
			)...,
		)
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
			return nil, statusAccountNotFound.Err()
		}
		if errors.Is(err, domain.ErrSearchFilterNotFound) {
			return nil, statusSearchFilterIDNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}

	return &accountproto.DeleteSearchFilterResponse{}, nil
}

// checkSearchFilterOwner rejects search filter requests that target an account
// other than the caller's own. Search filters are personal settings, and the
// environment role check only authorizes the caller, so a request email that
// differs from the authenticated editor is denied.
func (s *AccountService) checkSearchFilterOwner(
	ctx context.Context,
	editor *eventproto.Editor,
	email, environmentID string,
) error {
	if email == editor.Email {
		return nil
	}
	s.logger.Warn(
		"Rejected search filter request for another account",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.String("email", editor.Email),
			zap.String("requestEmail", email),
			zap.String("environmentID", environmentID),
		)...,
	)
	return statusPermissionDenied.Err()
}
