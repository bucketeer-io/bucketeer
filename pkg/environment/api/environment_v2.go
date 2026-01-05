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
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

var (
	maxEnvironmentNameLength = 50
	environmentUrlCodeRegex  = regexp.MustCompile("^[a-z0-9-]{1,50}$")
)

func (s *EnvironmentService) GetEnvironmentV2(
	ctx context.Context,
	req *environmentproto.GetEnvironmentV2Request,
) (*environmentproto.GetEnvironmentV2Response, error) {
	_, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.Id,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetEnvironmentV2Request(req); err != nil {
		return nil, err
	}
	environment, err := s.environmentStorage.GetEnvironmentV2(ctx, req.Id)
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound {
			return nil, statusEnvironmentNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.GetEnvironmentV2Response{
		Environment: environment.EnvironmentV2,
	}, nil
}

func validateGetEnvironmentV2Request(
	req *environmentproto.GetEnvironmentV2Request,
) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
	return nil
}

func (s *EnvironmentService) ListEnvironmentsV2(
	ctx context.Context,
	req *environmentproto.ListEnvironmentsV2Request,
) (*environmentproto.ListEnvironmentsV2Response, error) {
	_, err := s.checkOrganizationRole(
		ctx,
		req.OrganizationId,
		accountproto.AccountV2_Role_Organization_MEMBER,
	)
	if err != nil {
		return nil, err
	}
	var filters []*mysql.FilterV2
	if req.ProjectId != "" {
		filters = append(filters, &mysql.FilterV2{
			Column:   "environment_v2.project_id",
			Operator: mysql.OperatorEqual,
			Value:    req.ProjectId,
		})
	}
	if req.OrganizationId != "" {
		filters = append(filters, &mysql.FilterV2{
			Column:   "environment_v2.organization_id",
			Operator: mysql.OperatorEqual,
			Value:    req.OrganizationId,
		})
	}
	if req.Archived != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "environment_v2.archived",
			Operator: mysql.OperatorEqual,
			Value:    req.Archived.Value,
		})
	}
	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{
				"environment_v2.id",
				"environment_v2.name",
				"environment_v2.url_code",
				"environment_v2.description",
			},
			Keyword: req.SearchKeyword,
		}
	}
	orders, err := s.newEnvironmentV2ListOrders(req.OrderBy, req.OrderDirection)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
		return nil, statusInvalidCursor.Err()
	}
	options := &mysql.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		Orders:      orders,
		SearchQuery: searchQuery,
		InFilters:   nil,
		NullFilters: nil,
		JSONFilters: nil,
	}
	environments, nextCursor, totalCount, err := s.environmentStorage.ListEnvironmentsV2(
		ctx,
		options,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list environments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.ListEnvironmentsV2Response{
		Environments: environments,
		Cursor:       strconv.Itoa(nextCursor),
		TotalCount:   totalCount,
	}, nil
}

func (s *EnvironmentService) newEnvironmentV2ListOrders(
	orderBy environmentproto.ListEnvironmentsV2Request_OrderBy,
	orderDirection environmentproto.ListEnvironmentsV2Request_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case environmentproto.ListEnvironmentsV2Request_DEFAULT,
		environmentproto.ListEnvironmentsV2Request_NAME:
		column = "environment_v2.name"
	case environmentproto.ListEnvironmentsV2Request_ID:
		column = "environment_v2.id"
	case environmentproto.ListEnvironmentsV2Request_URL_CODE:
		column = "environment_v2.url_code"
	case environmentproto.ListEnvironmentsV2Request_CREATED_AT:
		column = "environment_v2.created_at"
	case environmentproto.ListEnvironmentsV2Request_UPDATED_AT:
		column = "environment_v2.updated_at"
	case environmentproto.ListEnvironmentsV2Request_FEATURE_COUNT:
		column = "feature_count"
	default:
		return nil, statusInvalidOrderBy.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == environmentproto.ListEnvironmentsV2Request_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *EnvironmentService) CreateEnvironmentV2(
	ctx context.Context,
	req *environmentproto.CreateEnvironmentV2Request,
) (*environmentproto.CreateEnvironmentV2Response, error) {
	if req.Command != nil {
		if err := validateCreateEnvironmentV2Request(req); err != nil {
			return nil, err
		}
	} else {
		if err := validateCreateEnvironmentV2RequestNoCommand(req); err != nil {
			return nil, err
		}
	}

	// Get project ID from request
	var projectID string
	if req.Command != nil {
		projectID = req.Command.ProjectId
	} else {
		projectID = req.ProjectId
	}

	// Validate the project and get the actual organization ID
	orgID, err := s.getOrganizationID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Check if the user has admin role for the validated organization
	editor, err := s.checkOrganizationRole(
		ctx,
		orgID,
		accountproto.AccountV2_Role_Organization_ADMIN,
	)
	if err != nil {
		return nil, err
	}

	if req.Command == nil {
		return s.createEnvironmentV2NoCommand(ctx, req, editor, orgID)
	}

	name := strings.TrimSpace(req.Command.Name)
	newEnvironment, err := domain.NewEnvironmentV2(
		name,
		req.Command.UrlCode,
		req.Command.Description,
		req.Command.ProjectId,
		orgID,
		req.Command.RequireComment,
		s.logger,
	)
	if err != nil {
		return nil, err
	}
	if err := s.createEnvironmentV2(ctx, req.Command, newEnvironment, editor); err != nil {
		return nil, err
	}
	return &environmentproto.CreateEnvironmentV2Response{
		Environment: newEnvironment.EnvironmentV2,
	}, nil
}

func (s *EnvironmentService) createEnvironmentV2NoCommand(
	ctx context.Context,
	req *environmentproto.CreateEnvironmentV2Request,
	editor *eventproto.Editor,
	orgID string,
) (*environmentproto.CreateEnvironmentV2Response, error) {
	if err := validateCreateEnvironmentV2RequestNoCommand(req); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	newEnvironment, err := domain.NewEnvironmentV2(
		name,
		req.UrlCode,
		req.Description,
		req.ProjectId,
		orgID,
		req.RequireComment,
		s.logger,
	)
	if err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		e, err := domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ENVIRONMENT,
			newEnvironment.Id,
			eventproto.Event_ENVIRONMENT_V2_CREATED,
			&eventproto.EnvironmentV2CreatedEvent{
				Id:             newEnvironment.Id,
				Name:           newEnvironment.Name,
				UrlCode:        newEnvironment.UrlCode,
				Description:    newEnvironment.Description,
				ProjectId:      newEnvironment.ProjectId,
				Archived:       newEnvironment.Archived,
				RequireComment: newEnvironment.RequireComment,
				CreatedAt:      newEnvironment.CreatedAt,
				UpdatedAt:      newEnvironment.UpdatedAt,
			},
			newEnvironment.EnvironmentV2,
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, e); err != nil {
			return err
		}
		return s.environmentStorage.CreateEnvironmentV2(ctxWithTx, newEnvironment)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrEnvironmentAlreadyExists) {
			return nil, statusEnvironmentAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create environment",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.CreateEnvironmentV2Response{
		Environment: newEnvironment.EnvironmentV2,
	}, nil
}

func validateCreateEnvironmentV2Request(
	req *environmentproto.CreateEnvironmentV2Request,
) error {
	name := strings.TrimSpace(req.Command.Name)
	if name == "" {
		return statusEnvironmentNameRequired.Err()
	}
	if len(name) > maxEnvironmentNameLength {
		return statusInvalidEnvironmentName.Err()
	}
	if !environmentUrlCodeRegex.MatchString(req.Command.UrlCode) {
		return statusInvalidEnvironmentUrlCode.Err()
	}
	if req.Command.ProjectId == "" {
		return statusProjectIDRequired.Err()
	}
	return nil
}

func validateCreateEnvironmentV2RequestNoCommand(
	req *environmentproto.CreateEnvironmentV2Request,
) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return statusEnvironmentNameRequired.Err()
	}
	if len(name) > maxEnvironmentNameLength {
		return statusInvalidEnvironmentName.Err()
	}
	if !environmentUrlCodeRegex.MatchString(req.UrlCode) {
		return statusInvalidEnvironmentUrlCode.Err()
	}
	if req.ProjectId == "" {
		return statusProjectIDRequired.Err()
	}
	return nil
}

func (s *EnvironmentService) getOrganizationID(
	ctx context.Context,
	projectID string,
) (string, error) {
	// enabled project must exist
	existingProject, err := s.getProject(ctx, projectID)
	if err != nil {
		return "", err
	}
	if existingProject.Disabled {
		return "", statusProjectDisabled.Err()
	}
	return existingProject.OrganizationId, nil
}

func (s *EnvironmentService) createEnvironmentV2(
	ctx context.Context,
	cmd command.Command,
	environment *domain.EnvironmentV2,
	editor *eventproto.Editor,
) error {
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		handler, err := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return s.environmentStorage.CreateEnvironmentV2(contextWithTx, environment)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrEnvironmentAlreadyExists) {
			return statusEnvironmentAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create environment",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

func (s *EnvironmentService) UpdateEnvironmentV2(
	ctx context.Context,
	req *environmentproto.UpdateEnvironmentV2Request,
) (*environmentproto.UpdateEnvironmentV2Response, error) {
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.Id,
	)
	if err != nil {
		return nil, err
	}
	commands := getUpdateEnvironmentV2Commands(req)

	if len(commands) == 0 {
		return s.updateEnvironmentV2NoCommand(ctx, req, editor)
	}

	if err := validateUpdateEnvironmentV2Request(req.Id, commands); err != nil {
		return nil, err
	}
	if err := s.updateEnvironmentV2(ctx, req.Id, commands, editor); err != nil {
		return nil, err
	}
	return &environmentproto.UpdateEnvironmentV2Response{}, nil
}

func (s *EnvironmentService) updateEnvironmentV2NoCommand(
	ctx context.Context,
	req *environmentproto.UpdateEnvironmentV2Request,
	editor *eventproto.Editor,
) (*environmentproto.UpdateEnvironmentV2Response, error) {
	if err := validateUpdateEnvironmentV2RequestNoCommand(req); err != nil {
		return nil, err
	}

	err := s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		environment, err := s.environmentStorage.GetEnvironmentV2(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		updated, err := environment.Update(req.Name, req.Description, req.RequireComment, req.Archived)
		if err != nil {
			return err
		}
		event, err := domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ENVIRONMENT,
			environment.Id,
			eventproto.Event_ENVIRONMENT_V2_UPDATED,
			&eventproto.EnvironmentV2UpdatedEvent{
				Id:             updated.Id,
				Name:           req.Name,
				Description:    req.Description,
				RequireComment: req.RequireComment,
			},
			updated,
			environment,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, event); err != nil {
			return err
		}
		return s.environmentStorage.UpdateEnvironmentV2(ctxWithTx, updated)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrEnvironmentNotFound) || errors.Is(err, v2es.ErrEnvironmentUnexpectedAffectedRows) {
			return nil, statusEnvironmentNotFound.Err()
		}
		s.logger.Error(
			"Failed to update environment",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.UpdateEnvironmentV2Response{}, nil
}

func (s *EnvironmentService) updateEnvironmentV2(
	ctx context.Context,
	envId string,
	commands []command.Command,
	editor *eventproto.Editor,
) error {
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		environment, err := s.environmentStorage.GetEnvironmentV2(contextWithTx, envId)
		if err != nil {
			return err
		}
		handler, err := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		if err != nil {
			return err
		}
		for _, c := range commands {
			if err := handler.Handle(ctx, c); err != nil {
				return err
			}
		}
		return s.environmentStorage.UpdateEnvironmentV2(contextWithTx, environment)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrEnvironmentNotFound) || errors.Is(err, v2es.ErrEnvironmentUnexpectedAffectedRows) {
			return statusEnvironmentNotFound.Err()
		}
		s.logger.Error(
			"Failed to update environment",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

func getUpdateEnvironmentV2Commands(req *environmentproto.UpdateEnvironmentV2Request) []command.Command {
	commands := make([]command.Command, 0)
	if req.RenameCommand != nil {
		commands = append(commands, req.RenameCommand)
	}
	if req.ChangeDescriptionCommand != nil {
		commands = append(commands, req.ChangeDescriptionCommand)
	}
	if req.ChangeRequireCommentCommand != nil {
		commands = append(commands, req.ChangeRequireCommentCommand)
	}
	return commands
}

func validateUpdateEnvironmentV2Request(id string, commands []command.Command) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
	for _, cmd := range commands {
		if c, ok := cmd.(*environmentproto.RenameEnvironmentV2Command); ok {
			newName := strings.TrimSpace(c.Name)
			if newName == "" {
				return statusEnvironmentNameRequired.Err()
			}
			if len(newName) > maxEnvironmentNameLength {
				return statusInvalidEnvironmentName.Err()
			}
		}
	}
	return nil
}

func validateUpdateEnvironmentV2RequestNoCommand(
	req *environmentproto.UpdateEnvironmentV2Request,
) error {
	if req.Name != nil {
		newName := strings.TrimSpace(req.Name.Value)
		if newName == "" {
			return statusEnvironmentNameRequired.Err()
		}
		if len(newName) > maxEnvironmentNameLength {
			return statusInvalidEnvironmentName.Err()
		}
	}
	return nil
}

func (s *EnvironmentService) ArchiveEnvironmentV2(
	ctx context.Context,
	req *environmentproto.ArchiveEnvironmentV2Request,
) (*environmentproto.ArchiveEnvironmentV2Response, error) {
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.Id,
	)
	if err != nil {
		return nil, err
	}
	if err := validateArchiveEnvironmentV2Request(req); err != nil {
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		environment, err := s.environmentStorage.GetEnvironmentV2(contextWithTx, req.Id)
		if err != nil {
			return err
		}
		handler, err := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.environmentStorage.UpdateEnvironmentV2(contextWithTx, environment)
	})
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound || err == v2es.ErrEnvironmentUnexpectedAffectedRows {
			return nil, statusEnvironmentNotFound.Err()
		}
		s.logger.Error(
			"Failed to archive environment",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.ArchiveEnvironmentV2Response{}, nil
}

func validateArchiveEnvironmentV2Request(
	req *environmentproto.ArchiveEnvironmentV2Request,
) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *EnvironmentService) UnarchiveEnvironmentV2(
	ctx context.Context,
	req *environmentproto.UnarchiveEnvironmentV2Request,
) (*environmentproto.UnarchiveEnvironmentV2Response, error) {
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.Id,
	)
	if err != nil {
		return nil, err
	}
	if err := validateUnarchiveEnvironmentV2Request(req); err != nil {
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		environment, err := s.environmentStorage.GetEnvironmentV2(contextWithTx, req.Id)
		if err != nil {
			return err
		}
		handler, err := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.environmentStorage.UpdateEnvironmentV2(contextWithTx, environment)
	})
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound || err == v2es.ErrEnvironmentUnexpectedAffectedRows {
			return nil, statusEnvironmentNotFound.Err()
		}
		s.logger.Error(
			"Failed to unarchive environment",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.UnarchiveEnvironmentV2Response{}, nil
}

func validateUnarchiveEnvironmentV2Request(
	req *environmentproto.UnarchiveEnvironmentV2Request,
) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}
