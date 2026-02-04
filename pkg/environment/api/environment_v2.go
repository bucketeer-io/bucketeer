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

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
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
	if err := validateCreateEnvironmentV2Request(req); err != nil {
		return nil, err
	}

	// Validate the project and get the actual organization ID
	orgID, err := s.getOrganizationID(ctx, req.ProjectId)
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
	localizer := locale.NewLocalizer(ctx)
	if err := validateUpdateEnvironmentV2Request(ctx, req); err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		environment, err := s.environmentStorage.GetEnvironmentV2(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		updated, err := environment.Update(
			req.Name,
			req.Description,
			req.RequireComment,
			req.Archived,
			req.AutoArchiveEnabled,
			req.AutoArchiveUnusedDays,
			req.AutoArchiveCheckCodeRefs,
		)
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
		if errors.Is(err, domain.ErrAutoArchiveUnusedDaysRequired) {
			dt, err := statusInvalidAutoArchiveUnusedDays.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "auto_archive_unused_days"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if errors.Is(err, domain.ErrAutoArchiveNotEnabled) {
			dt, err := statusAutoArchiveNotEnabled.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "auto_archive_settings"),
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
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.UpdateEnvironmentV2Response{}, nil
}

func validateUpdateEnvironmentV2Request(
	ctx context.Context,
	req *environmentproto.UpdateEnvironmentV2Request,
) error {
	localizer := locale.NewLocalizer(ctx)
	if req.Name != nil {
		newName := strings.TrimSpace(req.Name.Value)
		if newName == "" {
			return statusEnvironmentNameRequired.Err()
		}
		if len(newName) > maxEnvironmentNameLength {
			return statusInvalidEnvironmentName.Err()
		}
	}
	// Auto-archive validation
	if req.AutoArchiveEnabled != nil && req.AutoArchiveEnabled.Value {
		if req.AutoArchiveUnusedDays == nil || req.AutoArchiveUnusedDays.Value <= 0 {
			dt, err := statusInvalidAutoArchiveUnusedDays.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "auto_archive_unused_days"),
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
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.Id,
	)
	if err != nil {
		s.logger.Error("Failed to check organization role",
			zap.Error(err),
			zap.String("id", req.Id),
		)
		return nil, err
	}

	event := &eventproto.Event{}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		environment, err := s.environmentStorage.GetEnvironmentV2(contextWithTx, req.Id)
		if err != nil {
			return err
		}
		prev := &domain.EnvironmentV2{}
		if err := copier.Copy(prev, environment); err != nil {
			return err
		}
		environment.SetArchived()
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ENVIRONMENT,
			environment.Id,
			eventproto.Event_ENVIRONMENT_V2_ARCHIVED,
			&eventproto.EnvironmentV2ArchivedEvent{
				Id:        environment.Id,
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
			},
			environment,
			prev,
		)
		if err != nil {
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
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish archive environment event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("event", event),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.ArchiveEnvironmentV2Response{}, nil
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
		s.logger.Error("Failed to check organization role",
			zap.Error(err),
			zap.String("id", req.Id),
		)
		return nil, err
	}

	event := &eventproto.Event{}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		environment, err := s.environmentStorage.GetEnvironmentV2(contextWithTx, req.Id)
		if err != nil {
			return err
		}
		prev := &domain.EnvironmentV2{}
		if err := copier.Copy(prev, environment); err != nil {
			return err
		}
		environment.SetUnarchived()
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ENVIRONMENT,
			environment.Id,
			eventproto.Event_ENVIRONMENT_V2_UNARCHIVED,
			&eventproto.EnvironmentV2UnarchivedEvent{
				Id:        environment.Id,
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
			},
			environment,
			prev,
		)
		if err != nil {
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
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish unarchive environment event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("event", event),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.UnarchiveEnvironmentV2Response{}, nil
}
