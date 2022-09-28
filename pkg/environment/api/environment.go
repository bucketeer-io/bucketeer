// Copyright 2022 The Bucketeer Authors.
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
	"regexp"
	"strconv"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/environment/command"
	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

var environmentIDRegex = regexp.MustCompile("^[a-z0-9-]{1,50}$")

func (s *EnvironmentService) GetEnvironment(
	ctx context.Context,
	req *environmentproto.GetEnvironmentRequest,
) (*environmentproto.GetEnvironmentResponse, error) {
	_, err := s.checkAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateGetEnvironmentRequest(req); err != nil {
		return nil, err
	}
	environmentStorage := v2es.NewEnvironmentStorage(s.mysqlClient)
	environment, err := environmentStorage.GetEnvironment(ctx, req.Id)
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound {
			return nil, localizedError(statusEnvironmentNotFound, locale.JaJP)
		}
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	if environment.Deleted {
		return nil, localizedError(statusEnvironmentAlreadyDeleted, locale.JaJP)
	}
	return &environmentproto.GetEnvironmentResponse{
		Environment: environment.Environment,
	}, nil
}

func validateGetEnvironmentRequest(req *environmentproto.GetEnvironmentRequest) error {
	if req.Id == "" {
		return localizedError(statusEnvironmentIDRequired, locale.JaJP)
	}
	return nil
}

func (s *EnvironmentService) GetEnvironmentByNamespace(
	ctx context.Context,
	req *environmentproto.GetEnvironmentByNamespaceRequest,
) (*environmentproto.GetEnvironmentByNamespaceResponse, error) {
	_, err := s.checkAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	environment, err := s.getEnvironmentByNamespace(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	return &environmentproto.GetEnvironmentByNamespaceResponse{
		Environment: environment,
	}, nil
}

func (s *EnvironmentService) getEnvironmentByNamespace(
	ctx context.Context,
	namespace string,
) (*environmentproto.Environment, error) {
	environmentStorage := v2es.NewEnvironmentStorage(s.mysqlClient)
	environment, err := environmentStorage.GetEnvironmentByNamespace(ctx, namespace, false)
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound {
			return nil, localizedError(statusEnvironmentNotFound, locale.JaJP)
		}
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return environment.Environment, nil
}

func (s *EnvironmentService) ListEnvironments(
	ctx context.Context,
	req *environmentproto.ListEnvironmentsRequest,
) (*environmentproto.ListEnvironmentsResponse, error) {
	_, err := s.checkAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{mysql.NewFilter("deleted", "=", false)}
	if req.ProjectId != "" {
		whereParts = append(whereParts, mysql.NewFilter("project_id", "=", req.ProjectId))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"id", "description"}, req.SearchKeyword))
	}
	orders, err := s.newEnvironmentListOrders(req.OrderBy, req.OrderDirection)
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
		return nil, localizedError(statusInvalidCursor, locale.JaJP)
	}
	environmentStorage := v2es.NewEnvironmentStorage(s.mysqlClient)
	environments, nextCursor, totalCount, err := environmentStorage.ListEnvironments(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list environments",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return &environmentproto.ListEnvironmentsResponse{
		Environments: environments,
		Cursor:       strconv.Itoa(nextCursor),
		TotalCount:   totalCount,
	}, nil
}

func (s *EnvironmentService) newEnvironmentListOrders(
	orderBy environmentproto.ListEnvironmentsRequest_OrderBy,
	orderDirection environmentproto.ListEnvironmentsRequest_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case environmentproto.ListEnvironmentsRequest_DEFAULT,
		environmentproto.ListEnvironmentsRequest_ID:
		column = "id"
	case environmentproto.ListEnvironmentsRequest_CREATED_AT:
		column = "created_at"
	case environmentproto.ListEnvironmentsRequest_UPDATED_AT:
		column = "updated_at"
	default:
		return nil, localizedError(statusInvalidOrderBy, locale.JaJP)
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == environmentproto.ListEnvironmentsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *EnvironmentService) CreateEnvironment(
	ctx context.Context,
	req *environmentproto.CreateEnvironmentRequest,
) (*environmentproto.CreateEnvironmentResponse, error) {
	editor, err := s.checkAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateCreateEnvironmentRequest(req); err != nil {
		return nil, err
	}
	if err := s.checkProjectExistence(ctx, req.Command.ProjectId); err != nil {
		return nil, err
	}
	newEnvironment := domain.NewEnvironment(req.Command.Id, req.Command.Description, req.Command.ProjectId)
	if err := s.createEnvironment(ctx, req.Command, newEnvironment, editor); err != nil {
		return nil, err
	}
	return &environmentproto.CreateEnvironmentResponse{}, nil
}

func (s *EnvironmentService) createEnvironment(
	ctx context.Context,
	cmd command.Command,
	environment *domain.Environment,
	editor *eventproto.Editor,
) error {
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return localizedError(statusInternal, locale.JaJP)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		environmentStorage := v2es.NewEnvironmentStorage(tx)
		handler := command.NewEnvironmentCommandHandler(editor, environment, s.publisher)
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return environmentStorage.CreateEnvironment(ctx, environment)
	})
	if err != nil {
		if err == v2es.ErrEnvironmentAlreadyExists {
			return localizedError(statusEnvironmentAlreadyExists, locale.JaJP)
		}
		s.logger.Error(
			"Failed to create environment",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return localizedError(statusInternal, locale.JaJP)
	}
	return nil
}

func validateCreateEnvironmentRequest(req *environmentproto.CreateEnvironmentRequest) error {
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if !environmentIDRegex.MatchString(req.Command.Id) {
		return localizedError(statusInvalidEnvironmentID, locale.JaJP)
	}
	if req.Command.ProjectId == "" {
		return localizedError(statusProjectIDRequired, locale.JaJP)
	}
	return nil
}

func (s *EnvironmentService) checkProjectExistence(ctx context.Context, projectID string) error {
	// enabled project must exist
	existingProject, err := s.getProject(ctx, projectID)
	if err != nil {
		return err
	}
	if existingProject.Disabled {
		return localizedError(statusProjectDisabled, locale.JaJP)
	}
	return nil
}

func (s *EnvironmentService) UpdateEnvironment(
	ctx context.Context,
	req *environmentproto.UpdateEnvironmentRequest,
) (*environmentproto.UpdateEnvironmentResponse, error) {
	editor, err := s.checkAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	commands := getUpdateEnvironmentCommands(req)
	if err := validateUpdateEnvironmentRequest(req.Id, commands); err != nil {
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		environmentStorage := v2es.NewEnvironmentStorage(tx)
		environment, err := environmentStorage.GetEnvironment(ctx, req.Id)
		if err != nil {
			return err
		}
		handler := command.NewEnvironmentCommandHandler(editor, environment, s.publisher)
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		return environmentStorage.UpdateEnvironment(ctx, environment)
	})
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound || err == v2es.ErrEnvironmentUnexpectedAffectedRows {
			return nil, localizedError(statusEnvironmentNotFound, locale.JaJP)
		}
		s.logger.Error(
			"Failed to update environment",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return &environmentproto.UpdateEnvironmentResponse{}, nil
}

func getUpdateEnvironmentCommands(req *environmentproto.UpdateEnvironmentRequest) []command.Command {
	commands := make([]command.Command, 0)
	if req.RenameCommand != nil {
		commands = append(commands, req.RenameCommand)
	}
	if req.ChangeDescriptionCommand != nil {
		commands = append(commands, req.ChangeDescriptionCommand)
	}
	return commands
}

func validateUpdateEnvironmentRequest(id string, commands []command.Command) error {
	if len(commands) == 0 {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if id == "" {
		return localizedError(statusEnvironmentIDRequired, locale.JaJP)
	}
	return nil
}

func (s *EnvironmentService) DeleteEnvironment(
	ctx context.Context,
	req *environmentproto.DeleteEnvironmentRequest,
) (*environmentproto.DeleteEnvironmentResponse, error) {
	editor, err := s.checkAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteEnvironmentRequest(req); err != nil {
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		environmentStorage := v2es.NewEnvironmentStorage(tx)
		environment, err := environmentStorage.GetEnvironment(ctx, req.Id)
		if err != nil {
			return err
		}
		handler := command.NewEnvironmentCommandHandler(editor, environment, s.publisher)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return environmentStorage.UpdateEnvironment(ctx, environment)
	})
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound || err == v2es.ErrEnvironmentUnexpectedAffectedRows {
			return nil, localizedError(statusEnvironmentNotFound, locale.JaJP)
		}
		s.logger.Error(
			"Failed to update environment",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, localizedError(statusInternal, locale.JaJP)
	}
	return &environmentproto.DeleteEnvironmentResponse{}, nil
}

func validateDeleteEnvironmentRequest(req *environmentproto.DeleteEnvironmentRequest) error {
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if req.Id == "" {
		return localizedError(statusEnvironmentIDRequired, locale.JaJP)
	}
	return nil
}
