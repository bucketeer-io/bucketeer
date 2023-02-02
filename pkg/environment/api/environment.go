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
	"google.golang.org/genproto/googleapis/rpc/errdetails"

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
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetEnvironmentRequest(req, localizer); err != nil {
		return nil, err
	}
	environmentStorage := v2es.NewEnvironmentStorage(s.mysqlClient)
	environment, err := environmentStorage.GetEnvironment(ctx, req.Id)
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound {
			dt, err := statusEnvironmentNotFound.WithDetails(&errdetails.LocalizedMessage{
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
	if environment.Deleted {
		dt, err := statusEnvironmentAlreadyDeleted.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AlreadyDeletedError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &environmentproto.GetEnvironmentResponse{
		Environment: environment.Environment,
	}, nil
}

func validateGetEnvironmentRequest(req *environmentproto.GetEnvironmentRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusEnvironmentIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *EnvironmentService) GetEnvironmentByNamespace(
	ctx context.Context,
	req *environmentproto.GetEnvironmentByNamespaceRequest,
) (*environmentproto.GetEnvironmentByNamespaceResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	environment, err := s.getEnvironmentByNamespace(ctx, req.Namespace, localizer)
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
	localizer locale.Localizer,
) (*environmentproto.Environment, error) {
	environmentStorage := v2es.NewEnvironmentStorage(s.mysqlClient)
	environment, err := environmentStorage.GetEnvironmentByNamespace(ctx, namespace, false)
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound {
			dt, err := statusEnvironmentNotFound.WithDetails(&errdetails.LocalizedMessage{
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
	return environment.Environment, nil
}

func (s *EnvironmentService) ListEnvironments(
	ctx context.Context,
	req *environmentproto.ListEnvironmentsRequest,
) (*environmentproto.ListEnvironmentsResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkAdminRole(ctx, localizer)
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
	orders, err := s.newEnvironmentListOrders(req.OrderBy, req.OrderDirection, localizer)
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
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
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
	localizer locale.Localizer,
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
	if orderDirection == environmentproto.ListEnvironmentsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *EnvironmentService) CreateEnvironment(
	ctx context.Context,
	req *environmentproto.CreateEnvironmentRequest,
) (*environmentproto.CreateEnvironmentResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateCreateEnvironmentRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.checkProjectExistence(ctx, req.Command.ProjectId, localizer); err != nil {
		return nil, err
	}
	newEnvironment := domain.NewEnvironment(req.Command.Id, req.Command.Description, req.Command.ProjectId)
	if err := s.createEnvironment(ctx, req.Command, newEnvironment, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.CreateEnvironmentResponse{}, nil
}

func (s *EnvironmentService) createEnvironment(
	ctx context.Context,
	cmd command.Command,
	environment *domain.Environment,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) error {
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
			return statusInternal.Err()
		}
		return dt.Err()
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
			dt, err := statusEnvironmentAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		s.logger.Error(
			"Failed to create environment",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateCreateEnvironmentRequest(
	req *environmentproto.CreateEnvironmentRequest,
	localizer locale.Localizer,
) error {
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !environmentIDRegex.MatchString(req.Command.Id) {
		dt, err := statusInvalidEnvironmentID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.ProjectId == "" {
		dt, err := statusProjectIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "project_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *EnvironmentService) checkProjectExistence(
	ctx context.Context,
	projectID string,
	localizer locale.Localizer,
) error {
	// enabled project must exist
	existingProject, err := s.getProject(ctx, projectID, localizer)
	if err != nil {
		return err
	}
	if existingProject.Disabled {
		dt, err := statusProjectDisabled.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.ProjectDisabled),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *EnvironmentService) UpdateEnvironment(
	ctx context.Context,
	req *environmentproto.UpdateEnvironmentRequest,
) (*environmentproto.UpdateEnvironmentResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	commands := getUpdateEnvironmentCommands(req)
	if err := validateUpdateEnvironmentRequest(req.Id, commands, localizer); err != nil {
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
			dt, err := statusEnvironmentNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to update environment",
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

func validateUpdateEnvironmentRequest(id string, commands []command.Command, localizer locale.Localizer) error {
	if len(commands) == 0 {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if id == "" {
		dt, err := statusEnvironmentIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *EnvironmentService) DeleteEnvironment(
	ctx context.Context,
	req *environmentproto.DeleteEnvironmentRequest,
) (*environmentproto.DeleteEnvironmentResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteEnvironmentRequest(req, localizer); err != nil {
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
			dt, err := statusEnvironmentNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to update environment",
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
	return &environmentproto.DeleteEnvironmentResponse{}, nil
}

func validateDeleteEnvironmentRequest(
	req *environmentproto.DeleteEnvironmentRequest,
	localizer locale.Localizer,
) error {
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Id == "" {
		dt, err := statusEnvironmentIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}
