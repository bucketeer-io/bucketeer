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
	maxOrganizationNameLength = 50
	organizationUrlCodeRegex  = regexp.MustCompile("^[a-z0-9-]{1,50}$")
)

func (s *EnvironmentService) GetOrganization(
	ctx context.Context,
	req *environmentproto.GetOrganizationRequest,
) (*environmentproto.GetOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	org, err := s.getOrganization(ctx, req.Id, localizer)
	if err != nil {
		return nil, err
	}
	return &environmentproto.GetOrganizationResponse{
		Organization: org.Organization,
	}, nil
}

func (s *EnvironmentService) validateGetOrganizationRequest(
	req *environmentproto.GetOrganizationRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusOrganizationIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) getOrganization(
	ctx context.Context,
	id string,
	localizer locale.Localizer,
) (*domain.Organization, error) {
	orgStorage := v2es.NewOrganizationStorage(s.mysqlClient)
	org, err := orgStorage.GetOrganization(ctx, id)
	if err != nil {
		if errors.Is(err, v2es.ErrOrganizationNotFound) {
			dt, err := statusOrganizationNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError),
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
	return org, nil
}

func (s *EnvironmentService) ListOrganizations(
	ctx context.Context,
	req *environmentproto.ListOrganizationsRequest,
) (*environmentproto.ListOrganizationsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("disabled", "=", req.Disabled.Value))
	}
	if req.Archived != nil {
		whereParts = append(whereParts, mysql.NewFilter("archived", "=", req.Archived.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(
			whereParts,
			mysql.NewSearchQuery(
				[]string{"id", "name", "url_code"},
				req.SearchKeyword,
			),
		)
	}
	orders, err := s.newOrganizationListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"failed to create OrganizationListOrders",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
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
	orgStorage := v2es.NewOrganizationStorage(s.mysqlClient)
	organizations, nextCursor, totalCount, err := orgStorage.ListOrganizations(ctx, whereParts, orders, limit, offset)
	if err != nil {
		s.logger.Error(
			"failed to list organizations",
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
	return &environmentproto.ListOrganizationsResponse{
		Organizations: organizations,
		Cursor:        strconv.Itoa(nextCursor),
		TotalCount:    totalCount,
	}, nil
}

func (s *EnvironmentService) newOrganizationListOrders(
	orderBy environmentproto.ListOrganizationsRequest_OrderBy,
	orderDirection environmentproto.ListOrganizationsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case environmentproto.ListOrganizationsRequest_DEFAULT,
		environmentproto.ListOrganizationsRequest_NAME:
		column = "name"
	case environmentproto.ListOrganizationsRequest_URL_CODE:
		column = "url_code"
	case environmentproto.ListOrganizationsRequest_ID:
		column = "id"
	case environmentproto.ListOrganizationsRequest_CREATED_AT:
		column = "created_at"
	case environmentproto.ListOrganizationsRequest_UPDATED_AT:
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
	if orderDirection == environmentproto.ListOrganizationsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *EnvironmentService) CreateOrganization(
	ctx context.Context,
	req *environmentproto.CreateOrganizationRequest,
) (*environmentproto.CreateOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Command.Name)
	urlCode := strings.TrimSpace(req.Command.UrlCode)
	organization, err := domain.NewOrganization(
		name,
		urlCode,
		req.Command.Description,
		req.Command.IsTrial,
		req.Command.IsSystemAdmin,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create organization",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if err := s.createOrganization(ctx, req.Command, organization, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.CreateOrganizationResponse{
		Organization: organization.Organization,
	}, nil
}

func (s *EnvironmentService) validateCreateOrganizationRequest(
	req *environmentproto.CreateOrganizationRequest,
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
		dt, err := statusOrganizationNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(name) > maxOrganizationNameLength {
		dt, err := statusInvalidOrganizationName.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	urlCode := strings.TrimSpace(req.Command.UrlCode)
	if !organizationUrlCodeRegex.MatchString(urlCode) {
		dt, err := statusInvalidOrganizationUrlCode.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *EnvironmentService) createOrganization(
	ctx context.Context,
	cmd command.Command,
	organization *domain.Organization,
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
		orgStorage := v2es.NewOrganizationStorage(tx)
		if organization.Organization.SystemAdmin {
			org, err := orgStorage.GetSystemAdminOrganization(ctx)
			if err != nil {
				return err
			}
			if org != nil {
				return v2es.ErrOrganizationAlreadyExists
			}
		}
		handler := command.NewOrganizationCommandHandler(editor, organization, s.publisher)
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return orgStorage.CreateOrganization(ctx, organization)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrOrganizationAlreadyExists) {
			dt, err := statusOrganizationAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		s.logger.Error(
			"Failed to create organization",
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
	return nil
}

func (s *EnvironmentService) UpdateOrganization(
	ctx context.Context,
	req *environmentproto.UpdateOrganizationRequest,
) (*environmentproto.UpdateOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	commands := s.getUpdateOrganizationCommands(req)
	if err := s.validateUpdateOrganizationRequest(req.Id, commands, localizer); err != nil {
		return nil, err
	}
	if err := s.updateOrganization(ctx, req.Id, editor, localizer, commands...); err != nil {
		return nil, err
	}
	return &environmentproto.UpdateOrganizationResponse{}, nil
}

func (s *EnvironmentService) getUpdateOrganizationCommands(
	req *environmentproto.UpdateOrganizationRequest,
) []command.Command {
	commands := make([]command.Command, 0)
	if req.ChangeDescriptionCommand != nil {
		commands = append(commands, req.ChangeDescriptionCommand)
	}
	if req.RenameCommand != nil {
		commands = append(commands, req.RenameCommand)
	}
	return commands
}

func (s *EnvironmentService) validateUpdateOrganizationRequest(
	id string,
	commands []command.Command,
	localizer locale.Localizer,
) error {
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
		dt, err := statusOrganizationIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, cmd := range commands {
		if c, ok := cmd.(*environmentproto.ChangeNameOrganizationCommand); ok {
			name := strings.TrimSpace(c.Name)
			if name == "" {
				dt, err := statusOrganizationNameRequired.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
			if len(name) > maxOrganizationNameLength {
				dt, err := statusInvalidOrganizationName.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) updateOrganization(
	ctx context.Context,
	id string,
	editor *eventproto.Editor,
	localizer locale.Localizer,
	commands ...command.Command,
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
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctx, id)
		if err != nil {
			return err
		}
		handler := command.NewOrganizationCommandHandler(editor, organization, s.publisher)
		for _, c := range commands {
			if err := handler.Handle(ctx, c); err != nil {
				return err
			}
		}
		return orgStorage.UpdateOrganization(ctx, organization)
	})
	if err != nil {
		s.logger.Error(
			"Failed to update organization",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		if errors.Is(err, domain.ErrCannotArchiveSystemAdmin) || errors.Is(err, domain.ErrCannotDisableSystemAdmin) {
			dt, err := statusCannotUpdateSystemAdmin.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if errors.Is(err, v2es.ErrOrganizationNotFound) || errors.Is(err, v2es.ErrOrganizationUnexpectedAffectedRows) {
			dt, err := statusOrganizationNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
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

func (s *EnvironmentService) EnableOrganization(
	ctx context.Context,
	req *environmentproto.EnableOrganizationRequest,
) (*environmentproto.EnableOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateEnableOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateOrganization(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.EnableOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateEnableOrganizationRequest(
	req *environmentproto.EnableOrganizationRequest,
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
		dt, err := statusOrganizationIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) DisableOrganization(
	ctx context.Context,
	req *environmentproto.DisableOrganizationRequest,
) (*environmentproto.DisableOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateDisableOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateOrganization(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.DisableOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateDisableOrganizationRequest(
	req *environmentproto.DisableOrganizationRequest,
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
		dt, err := statusOrganizationIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) ArchiveOrganization(
	ctx context.Context,
	req *environmentproto.ArchiveOrganizationRequest,
) (*environmentproto.ArchiveOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateArchiveOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateOrganization(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.ArchiveOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateArchiveOrganizationRequest(
	req *environmentproto.ArchiveOrganizationRequest,
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
		dt, err := statusOrganizationIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) UnarchiveOrganization(
	ctx context.Context,
	req *environmentproto.UnarchiveOrganizationRequest,
) (*environmentproto.UnarchiveOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateUnarchiveOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateOrganization(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.UnarchiveOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateUnarchiveOrganizationRequest(
	req *environmentproto.UnarchiveOrganizationRequest,
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
		dt, err := statusOrganizationIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) ConvertTrialOrganization(
	ctx context.Context,
	req *environmentproto.ConvertTrialOrganizationRequest,
) (*environmentproto.ConvertTrialOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateConvertTrialOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateOrganization(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.ConvertTrialOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateConvertTrialOrganizationRequest(
	req *environmentproto.ConvertTrialOrganizationRequest,
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
		dt, err := statusOrganizationIDRequired.WithDetails(&errdetails.LocalizedMessage{
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
