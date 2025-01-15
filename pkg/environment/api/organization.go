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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	accdomain "github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2acc "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
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
		// A member can access the organization details but can't update it
		_, err = s.checkOrganizationRole(ctx, req.Id, accountproto.AccountV2_Role_Organization_MEMBER, localizer)
		if err != nil {
			return nil, err
		}
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
		whereParts = append(whereParts, mysql.NewFilter("organization.disabled", "=", req.Disabled.Value))
	}
	if req.Archived != nil {
		whereParts = append(whereParts, mysql.NewFilter("organization.archived", "=", req.Archived.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(
			whereParts,
			mysql.NewSearchQuery(
				[]string{"organization.id", "organization.name", "organization.url_code"},
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
		column = "organization.name"
	case environmentproto.ListOrganizationsRequest_URL_CODE:
		column = "organization.url_code"
	case environmentproto.ListOrganizationsRequest_ID:
		column = "organization.id"
	case environmentproto.ListOrganizationsRequest_CREATED_AT:
		column = "organization.created_at"
	case environmentproto.ListOrganizationsRequest_UPDATED_AT:
		column = "organization.updated_at"
	case environmentproto.ListOrganizationsRequest_ENVIRONMENT_COUNT:
		column = "environments"
	case environmentproto.ListOrganizationsRequest_PROJECT_COUNT:
		column = "projects"
	case environmentproto.ListOrganizationsRequest_USER_COUNT:
		column = "users"

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
	if req.Command == nil {
		return s.createOrganizationNoCommand(
			ctx,
			req,
			editor,
			localizer,
		)
	}
	if err := s.validateCreateOrganizationRequest(req, localizer); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Command.Name)
	urlCode := strings.TrimSpace(req.Command.UrlCode)
	organization, err := domain.NewOrganization(
		name,
		urlCode,
		req.Command.OwnerEmail,
		req.Command.Description,
		req.Command.IsTrial,
		req.Command.IsSystemAdmin,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create an organization",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, statusInternal.Err()
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
	if !emailRegex.MatchString(req.Command.OwnerEmail) {
		dt, err := statusInvalidOrganizationCreatorEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "owner_email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *EnvironmentService) createOrganizationNoCommand(
	ctx context.Context,
	req *environmentproto.CreateOrganizationRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*environmentproto.CreateOrganizationResponse, error) {
	if err := s.validateCreateOrganizationRequestNoCommand(req, localizer); err != nil {
		return nil, err
	}
	// Create the organization
	name := strings.TrimSpace(req.Name)
	urlCode := strings.TrimSpace(req.UrlCode)
	organization, err := domain.NewOrganization(
		name,
		urlCode,
		req.OwnerEmail,
		req.Description,
		req.IsTrial,
		req.IsSystemAdmin,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create an organization",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, statusInternal.Err()
	}
	var envRoles []*accountproto.AccountV2_EnvironmentRole
	// Begin the SQL transaction
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
		orgStorage := v2es.NewOrganizationStorage(tx)
		// Check if there is already a system admin organization
		if organization.Organization.SystemAdmin {
			org, err := orgStorage.GetSystemAdminOrganization(ctx)
			if err != nil {
				return err
			}
			if org != nil {
				return v2es.ErrOrganizationAlreadyExists
			}
		}
		if err := orgStorage.CreateOrganization(ctx, organization); err != nil {
			return err
		}
		// Create a default project
		project, err := s.createDefaultProject(
			ctx,
			tx,
			organization.Id,
			organization.OwnerEmail,
		)
		if err != nil {
			s.logger.Error(
				"Failed to create the default project",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
			return err
		}
		// Create default environments
		envRoles, err = s.createDefaultEnvironments(
			ctx,
			tx,
			organization.Id,
			project,
		)
		if err != nil {
			s.logger.Error(
				"Failed to create the default environments",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2es.ErrOrganizationAlreadyExists) {
			dt, err := statusOrganizationAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create an organization",
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
	// Create the admin account using the environment roles created in the last step
	// Because the account storage has a different implementation,
	// we can't create the account using the same transaction when creating the organization
	if err := s.createOwnerAccount(
		ctx,
		organization.Id,
		organization.OwnerEmail,
		envRoles,
	); err != nil {
		s.logger.Error(
			"Failed to create the owner account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, statusInternal.Err()
	}
	// Publish the auditlog event
	prev := &domain.Organization{}
	if err = copier.Copy(prev, organization); err != nil {
		return nil, statusInternal.Err()
	}
	event, err := domainevent.NewAdminEvent(
		editor,
		eventproto.Event_ORGANIZATION,
		organization.Id,
		eventproto.Event_ORGANIZATION_CREATED,
		&eventproto.OrganizationCreatedEvent{
			Id:          organization.Id,
			Name:        organization.Name,
			UrlCode:     organization.UrlCode,
			OwnerEmail:  organization.OwnerEmail,
			Description: organization.Description,
			Disabled:    organization.Disabled,
			Archived:    organization.Archived,
			Trial:       organization.Trial,
			CreatedAt:   organization.CreatedAt,
			UpdatedAt:   organization.UpdatedAt,
		},
		organization,
		prev,
	)
	if err != nil {
		return nil, statusInternal.Err()
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		return nil, statusInternal.Err()
	}
	return &environmentproto.CreateOrganizationResponse{
		Organization: organization.Organization,
	}, nil
}

func (s *EnvironmentService) validateCreateOrganizationRequestNoCommand(
	req *environmentproto.CreateOrganizationRequest,
	localizer locale.Localizer,
) error {
	name := strings.TrimSpace(req.Name)
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
	urlCode := strings.TrimSpace(req.UrlCode)
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
	if !emailRegex.MatchString(req.OwnerEmail) {
		dt, err := statusInvalidOrganizationCreatorEmail.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "owner_email"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

// Deprecated
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
		handler, err := command.NewOrganizationCommandHandler(editor, organization, s.publisher)
		if err != nil {
			return err
		}
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

// Create a default project.
// We must create a project when creating an Organization
// because we also need to create the organization owner in the account table.
// To create it we need the project, so we can also create the environment.
func (s *EnvironmentService) createDefaultProject(
	ctx context.Context,
	tx mysql.Transaction,
	organizationID, email string,
) (*domain.Project, error) {
	project, err := domain.NewProject(
		"Default Project",
		"default",
		"",
		email,
		organizationID,
		true,
	)
	if err != nil {
		return nil, err
	}
	projectStorage := v2es.NewProjectStorage(tx)
	if err := projectStorage.CreateProject(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

// Create Development and Production default environments.
// We must create the environment when creating an Organization
// because we also need to create the organization owner in the account table
// and to create it we need the organization and environment roles.
func (s *EnvironmentService) createDefaultEnvironments(
	ctx context.Context,
	tx mysql.Transaction,
	organizationID string,
	project *domain.Project,
) ([]*accountproto.AccountV2_EnvironmentRole, error) {
	envRoles := make([]*accountproto.AccountV2_EnvironmentRole, 0, 2)
	envNames := []string{
		"Development",
		"Production",
	}
	for _, name := range envNames {
		envURLCode := fmt.Sprintf("%s-%s", project.UrlCode, strings.ToLower(name))
		env, err := domain.NewEnvironmentV2(
			name,
			envURLCode,
			"",
			project.Id,
			organizationID,
			false,
			s.logger,
		)
		if err != nil {
			return nil, err
		}
		environmentStorage := v2es.NewEnvironmentStorage(tx)
		if err := environmentStorage.CreateEnvironmentV2(ctx, env); err != nil {
			return nil, err
		}
		envRoles = append(envRoles, &accountproto.AccountV2_EnvironmentRole{
			EnvironmentId: env.Id,
			Role:          accountproto.AccountV2_Role_Environment_EDITOR,
		})
	}
	return envRoles, nil
}

func (s *EnvironmentService) createOwnerAccount(
	ctx context.Context,
	organizationID, ownerEmail string,
	envRoles []*accountproto.AccountV2_EnvironmentRole,
) error {
	account := accdomain.NewAccountV2(
		ownerEmail,
		strings.Split(ownerEmail, "@")[0],
		"",
		"",
		"",
		"",
		[]string{},
		organizationID,
		accountproto.AccountV2_Role_Organization_OWNER,
		envRoles,
	)
	accountStorage := v2acc.NewAccountStorage(s.mysqlClient)
	if err := accountStorage.CreateAccountV2(ctx, account); err != nil {
		return err
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
		// If not system admin, check if user is organization owner
		editor, err = s.checkOrganizationRole(ctx, req.Id, accountproto.AccountV2_Role_Organization_OWNER, localizer)
		if err != nil {
			return nil, err
		}
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
	if req.ChangeOwnerEmailCommand != nil {
		commands = append(commands, req.ChangeOwnerEmailCommand)
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
	var prevOwnerEmail string
	var newOwnerEmail string
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctx, id)
		if err != nil {
			return err
		}
		prevOwnerEmail = organization.OwnerEmail
		handler, err := command.NewOrganizationCommandHandler(editor, organization, s.publisher)
		if err != nil {
			return err
		}
		for _, c := range commands {
			if err := handler.Handle(ctx, c); err != nil {
				return err
			}
		}
		// Set the new owner email if it changes
		if prevOwnerEmail != organization.OwnerEmail {
			newOwnerEmail = organization.OwnerEmail
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
	// Update the organization role when the owner email changes
	if prevOwnerEmail != "" && newOwnerEmail != "" {
		if err := s.updateOwnerRole(ctx, id, prevOwnerEmail, newOwnerEmail); err != nil {
			s.logger.Error("Failed to update the new owner's role",
				zap.Error(err),
				zap.String("organizationId", id),
				zap.String("prevOwnerEmail", prevOwnerEmail),
				zap.String("newOwnerEmail", newOwnerEmail),
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
	}
	return nil
}

func (s *EnvironmentService) updateOwnerRole(
	ctx context.Context,
	organizationID, prevOwnerEmail, newOwnerEmail string,
) error {
	accStorage := v2acc.NewAccountStorage(s.mysqlClient)
	// Update the old owner organization role
	prevOwnerAcc, err := accStorage.GetAccountV2(ctx, prevOwnerEmail, organizationID)
	if err != nil {
		return err
	}
	if err := prevOwnerAcc.ChangeOrganizationRole(accountproto.AccountV2_Role_Organization_ADMIN); err != nil {
		return err
	}
	if err := accStorage.UpdateAccountV2(ctx, prevOwnerAcc); err != nil {
		return err
	}
	// Update the new owner organization role
	newOwnerAcc, err := accStorage.GetAccountV2(ctx, newOwnerEmail, organizationID)
	if err != nil {
		return err
	}
	if err := newOwnerAcc.ChangeOrganizationRole(accountproto.AccountV2_Role_Organization_OWNER); err != nil {
		return err
	}
	if err := accStorage.UpdateAccountV2(ctx, newOwnerAcc); err != nil {
		return err
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
