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

	"github.com/bucketeer-io/bucketeer/pkg/account/command"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"

	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
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
			"Failed to create search filter",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	account, err := s.getAccountV2(ctx, req.Email, req.OrganizationId, localizer)
	if err != nil {
		return nil, err
	}
	// Since there is only one default setting for a filter target, set the existing default to OFF.
	commands := make([]command.Command, 0)
	if req.Command.DefaultFilter {
		changeDefaultFilterCmds := getChangeDefaultSearchFilterCommands(
			account,
			req.Command.FilterTargetType,
			req.Command.EnvironmentId,
		)
		for _, changeDefaultFilterCmd := range changeDefaultFilterCmds {
			commands = append(commands, changeDefaultFilterCmd)
		}
	}
	commands = append(commands, req.Command)

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
			"Failed to create search filter",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
				zap.String("searchFilterName", req.Command.Name),
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
	return &accountproto.CreateSearchFilterResponse{}, nil
}

func getChangeDefaultSearchFilterCommands(
	account *domain.AccountV2,
	filterTarget accountproto.FilterTargetType,
	environmentId string,
) []*accountproto.ChangeDefaultSearchFilterCommand {
	var changeDefaultFiltersCmds []*accountproto.ChangeDefaultSearchFilterCommand
	for _, filter := range account.SearchFilters {
		if filter.DefaultFilter &&
			filterTarget == filter.FilterTargetType &&
			environmentId == filter.EnvironmentId {
			changeDefaultFiltersCmds = append(changeDefaultFiltersCmds, &accountproto.ChangeDefaultSearchFilterCommand{
				Id:            filter.Id,
				DefaultFilter: false,
			})
		}
	}
	return changeDefaultFiltersCmds
}
