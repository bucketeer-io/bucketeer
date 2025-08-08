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
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/environment/command"
	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

var (
	maxEnvironmentNameLength = 50
	environmentUrlCodeRegex  = regexp.MustCompile("^[a-z0-9-]{1,50}$")
)

func (s *EnvironmentService) GetEnvironmentV2(
	ctx context.Context,
	req *environmentproto.GetEnvironmentV2Request,
) (*environmentproto.GetEnvironmentV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.Id,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetEnvironmentV2Request(req, localizer); err != nil {
		return nil, err
	}
	environment, err := s.environmentStorage.GetEnvironmentV2(ctx, req.Id)
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
	_, err := s.checkOrganizationRole(
		ctx,
		req.OrganizationId,
		accountproto.AccountV2_Role_Organization_MEMBER,
		localizer,
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
	orders, err := s.newEnvironmentV2ListOrders(req.OrderBy, req.OrderDirection, localizer)
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
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
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
	if req.Command != nil {
		if err := validateCreateEnvironmentV2Request(req, localizer); err != nil {
			return nil, err
		}
	} else {
		if err := validateCreateEnvironmentV2RequestNoCommand(req, localizer); err != nil {
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
	orgID, err := s.getOrganizationID(ctx, projectID, localizer)
	if err != nil {
		return nil, err
	}

	// Check if the user has admin role for the validated organization
	editor, err := s.checkOrganizationRole(
		ctx,
		orgID,
		accountproto.AccountV2_Role_Organization_ADMIN,
		localizer,
	)
	if err != nil {
		return nil, err
	}

	if req.Command == nil {
		return s.createEnvironmentV2NoCommand(ctx, req, localizer, editor, orgID)
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
	if err := s.createEnvironmentV2(ctx, req.Command, newEnvironment, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.CreateEnvironmentV2Response{
		Environment: newEnvironment.EnvironmentV2,
	}, nil
}

func (s *EnvironmentService) createEnvironmentV2NoCommand(
	ctx context.Context,
	req *environmentproto.CreateEnvironmentV2Request,
	localizer locale.Localizer,
	editor *eventproto.Editor,
	orgID string,
) (*environmentproto.CreateEnvironmentV2Response, error) {
	if err := validateCreateEnvironmentV2RequestNoCommand(req, localizer); err != nil {
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
			dt, err := statusEnvironmentAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create environment",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
	return &environmentproto.CreateEnvironmentV2Response{
		Environment: newEnvironment.EnvironmentV2,
	}, nil
}

func validateCreateEnvironmentV2Request(
	req *environmentproto.CreateEnvironmentV2Request,
	localizer locale.Localizer,
) error {
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

func validateCreateEnvironmentV2RequestNoCommand(
	req *environmentproto.CreateEnvironmentV2Request,
	localizer locale.Localizer,
) error {
	name := strings.TrimSpace(req.Name)
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
	if !environmentUrlCodeRegex.MatchString(req.UrlCode) {
		dt, err := statusInvalidEnvironmentUrlCode.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.ProjectId == "" {
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

func (s *EnvironmentService) getOrganizationID(
	ctx context.Context,
	projectID string,
	localizer locale.Localizer,
) (string, error) {
	// enabled project must exist
	existingProject, err := s.getProject(ctx, projectID, localizer)
	if err != nil {
		return "", err
	}
	if existingProject.Disabled {
		dt, err := statusProjectDisabled.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.ProjectDisabled),
		})
		if err != nil {
			return "", statusInternal.Err()
		}
		return "", dt.Err()
	}
	return existingProject.OrganizationId, nil
}

func (s *EnvironmentService) createEnvironmentV2(
	ctx context.Context,
	cmd command.Command,
	environment *domain.EnvironmentV2,
	editor *eventproto.Editor,
	localizer locale.Localizer,
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
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.Id,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	commands := getUpdateEnvironmentV2Commands(req)

	if len(commands) == 0 {
		return s.updateEnvironmentV2NoCommand(ctx, req, localizer, editor)
	}

	if err := validateUpdateEnvironmentV2Request(req.Id, commands, localizer); err != nil {
		return nil, err
	}
	if err := s.updateEnvironmentV2(ctx, req.Id, commands, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.UpdateEnvironmentV2Response{}, nil
}

func (s *EnvironmentService) updateEnvironmentV2NoCommand(
	ctx context.Context,
	req *environmentproto.UpdateEnvironmentV2Request,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) (*environmentproto.UpdateEnvironmentV2Response, error) {
	if err := validateUpdateEnvironmentV2RequestNoCommand(req, localizer); err != nil {
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
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
	return &environmentproto.UpdateEnvironmentV2Response{}, nil
}

func (s *EnvironmentService) updateEnvironmentV2(
	ctx context.Context,
	envId string,
	commands []command.Command,
	editor *eventproto.Editor,
	localizer locale.Localizer,
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
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
	if req.ChangeRequireCommentCommand != nil {
		commands = append(commands, req.ChangeRequireCommentCommand)
	}
	return commands
}

func validateUpdateEnvironmentV2Request(id string, commands []command.Command, localizer locale.Localizer) error {
	// Essentially, the id field is required, but no validation is performed because some older services do not have ID.
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

func validateUpdateEnvironmentV2RequestNoCommand(
	req *environmentproto.UpdateEnvironmentV2Request,
	localizer locale.Localizer,
) error {
	if req.Name != nil {
		newName := strings.TrimSpace(req.Name.Value)
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
	return nil
}

func (s *EnvironmentService) ArchiveEnvironmentV2(
	ctx context.Context,
	req *environmentproto.ArchiveEnvironmentV2Request,
) (*environmentproto.ArchiveEnvironmentV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.Id,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateArchiveEnvironmentV2Request(req, localizer); err != nil {
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
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.Id,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateUnarchiveEnvironmentV2Request(req, localizer); err != nil {
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
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
