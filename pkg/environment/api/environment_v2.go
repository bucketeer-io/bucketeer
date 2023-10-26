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
	"regexp"
	"strconv"
	"strings"

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

var (
	maxEnvironmentNameLength = 50
	environmentUrlCodeRegex  = regexp.MustCompile("^[a-z0-9-_.]{1,50}$")
)

func (s *EnvironmentService) GetEnvironmentV2(
	ctx context.Context,
	req *environmentproto.GetEnvironmentV2Request,
) (*environmentproto.GetEnvironmentV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetEnvironmentV2Request(req, localizer); err != nil {
		return nil, err
	}
	environmentStorage := v2es.NewEnvironmentStorage(s.mysqlClient)
	environment, err := environmentStorage.GetEnvironmentV2(ctx, req.Id)
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
	return &environmentproto.GetEnvironmentV2Response{
		Environment: environment.EnvironmentV2,
	}, nil
}

func validateGetEnvironmentV2Request(
	req *environmentproto.GetEnvironmentV2Request,
	localizer locale.Localizer,
) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
	return nil
}

func (s *EnvironmentService) ListEnvironmentsV2(
	ctx context.Context,
	req *environmentproto.ListEnvironmentsV2Request,
) (*environmentproto.ListEnvironmentsV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	var whereParts []mysql.WherePart
	if req.ProjectId != "" {
		whereParts = append(whereParts, mysql.NewFilter("project_id", "=", req.ProjectId))
	}
	if req.Archived != nil {
		whereParts = append(whereParts, mysql.NewFilter("archived", "=", req.Archived.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(
			whereParts,
			mysql.NewSearchQuery([]string{"id", "name", "url_code", "description"},
				req.SearchKeyword,
			),
		)
	}
	orders, err := s.newEnvironmentV2ListOrders(req.OrderBy, req.OrderDirection, localizer)
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
	environments, nextCursor, totalCount, err := environmentStorage.ListEnvironmentsV2(
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
	return &environmentproto.ListEnvironmentsV2Response{
		Environments: environments,
		Cursor:       strconv.Itoa(nextCursor),
		TotalCount:   totalCount,
	}, nil
}

func (s *EnvironmentService) newEnvironmentV2ListOrders(
	orderBy environmentproto.ListEnvironmentsV2Request_OrderBy,
	orderDirection environmentproto.ListEnvironmentsV2Request_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case environmentproto.ListEnvironmentsV2Request_DEFAULT,
		environmentproto.ListEnvironmentsV2Request_NAME:
		column = "name"
	case environmentproto.ListEnvironmentsV2Request_ID:
		column = "id"
	case environmentproto.ListEnvironmentsV2Request_URL_CODE:
		column = "url_code"
	case environmentproto.ListEnvironmentsV2Request_CREATED_AT:
		column = "created_at"
	case environmentproto.ListEnvironmentsV2Request_UPDATED_AT:
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
	if orderDirection == environmentproto.ListEnvironmentsV2Request_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *EnvironmentService) CreateEnvironmentV2(
	ctx context.Context,
	req *environmentproto.CreateEnvironmentV2Request,
) (*environmentproto.CreateEnvironmentV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateCreateEnvironmentV2Request(req, localizer); err != nil {
		return nil, err
	}
	if err := s.checkProjectExistence(ctx, req.Command.ProjectId, localizer); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Command.Name)
	newEnvironment, err := domain.NewEnvironmentV2(
		name,
		req.Command.UrlCode,
		req.Command.Description,
		req.Command.ProjectId,
		s.logger,
	)
	if err != nil {
		return nil, err
	}
	if err := s.createEnvironmentV2(ctx, req.Command, newEnvironment, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.CreateEnvironmentV2Response{
		Environment: newEnvironment.EnvironmentV2,
	}, nil
}

func validateCreateEnvironmentV2Request(
	req *environmentproto.CreateEnvironmentV2Request,
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
	name := strings.TrimSpace(req.Command.Name)
	if name == "" {
		dt, err := statusEnvironmentNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(name) > maxEnvironmentNameLength {
		dt, err := statusInvalidEnvironmentName.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !environmentUrlCodeRegex.MatchString(req.Command.UrlCode) {
		dt, err := statusInvalidEnvironmentUrlCode.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
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

func (s *EnvironmentService) createEnvironmentV2(
	ctx context.Context,
	cmd command.Command,
	environment *domain.EnvironmentV2,
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
		handler := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return environmentStorage.CreateEnvironmentV2(ctx, environment)
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

func (s *EnvironmentService) UpdateEnvironmentV2(
	ctx context.Context,
	req *environmentproto.UpdateEnvironmentV2Request,
) (*environmentproto.UpdateEnvironmentV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	commands := getUpdateEnvironmentV2Commands(req)
	if err := validateUpdateEnvironmentV2Request(req.Id, commands, localizer); err != nil {
		return nil, err
	}
	if err := s.updateEnvironmentV2(ctx, req.Id, commands, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.UpdateEnvironmentV2Response{}, nil
}

func (s *EnvironmentService) updateEnvironmentV2(
	ctx context.Context,
	envId string,
	commands []command.Command,
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
		environment, err := environmentStorage.GetEnvironmentV2(ctx, envId)
		if err != nil {
			return err
		}
		handler := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		for _, c := range commands {
			if err := handler.Handle(ctx, c); err != nil {
				return err
			}
		}
		return environmentStorage.UpdateEnvironmentV2(ctx, environment)
	})
	if err != nil {
		if err == v2es.ErrEnvironmentNotFound || err == v2es.ErrEnvironmentUnexpectedAffectedRows {
			dt, err := statusEnvironmentNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
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
			return statusInternal.Err()
		}
		return dt.Err()
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
	return commands
}

func validateUpdateEnvironmentV2Request(id string, commands []command.Command, localizer locale.Localizer) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
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
	for _, cmd := range commands {
		if c, ok := cmd.(*environmentproto.RenameEnvironmentV2Command); ok {
			newName := strings.TrimSpace(c.Name)
			if newName == "" {
				dt, err := statusEnvironmentNameRequired.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
			if len(newName) > maxEnvironmentNameLength {
				dt, err := statusInvalidEnvironmentName.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
		}
	}
	return nil
}

func (s *EnvironmentService) ArchiveEnvironmentV2(
	ctx context.Context,
	req *environmentproto.ArchiveEnvironmentV2Request,
) (*environmentproto.ArchiveEnvironmentV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateArchiveEnvironmentV2Request(req, localizer); err != nil {
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
		environment, err := environmentStorage.GetEnvironmentV2(ctx, req.Id)
		if err != nil {
			return err
		}
		handler := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return environmentStorage.UpdateEnvironmentV2(ctx, environment)
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
			"Failed to archive environment",
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
	return &environmentproto.ArchiveEnvironmentV2Response{}, nil
}

func validateArchiveEnvironmentV2Request(
	req *environmentproto.ArchiveEnvironmentV2Request,
	localizer locale.Localizer,
) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
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
	return nil
}

func (s *EnvironmentService) UnarchiveEnvironmentV2(
	ctx context.Context,
	req *environmentproto.UnarchiveEnvironmentV2Request,
) (*environmentproto.UnarchiveEnvironmentV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateUnarchiveEnvironmentV2Request(req, localizer); err != nil {
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
		environment, err := environmentStorage.GetEnvironmentV2(ctx, req.Id)
		if err != nil {
			return err
		}
		handler := command.NewEnvironmentV2CommandHandler(editor, environment, s.publisher)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return environmentStorage.UpdateEnvironmentV2(ctx, environment)
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
			"Failed to unarchive environment",
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
	return &environmentproto.UnarchiveEnvironmentV2Response{}, nil
}

func validateUnarchiveEnvironmentV2Request(
	req *environmentproto.UnarchiveEnvironmentV2Request,
	localizer locale.Localizer,
) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
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
	return nil
}
