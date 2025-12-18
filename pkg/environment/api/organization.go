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
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	accdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2acc "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

var (
	maxOrganizationNameLength = 50
	organizationUrlCodeRegex  = regexp.MustCompile("^[a-z0-9-]{1,50}$")
	featureActiveDays         = 7 * 24 * time.Hour

	targetEntitiesInEnvironment = []string{
		"subscription",
		"experiment_result",
		"push",
		"ops_count",
		"auto_ops_rule",
		"segment_user",
		"segment",
		"goal",
		"experiment",
		"tag",
		"ops_progressive_rollout",
		"flag_trigger",
		"code_reference",
		"feature",
		"api_key",
		"audit_log",
	}
	targetEntitiesInOrganization = []string{
		"account_v2",
	}
)

func (s *EnvironmentService) GetOrganization(
	ctx context.Context,
	req *environmentproto.GetOrganizationRequest,
) (*environmentproto.GetOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRole(ctx, req.Id, accountproto.AccountV2_Role_Organization_MEMBER, localizer)
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
	org, err := s.orgStorage.GetOrganization(ctx, id)
	if err != nil {
		s.logger.Error("failed to get organization",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("organizationId", id),
				zap.Error(err),
			)...)
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
		return nil, api.NewGRPCStatus(err).Err()
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
	var filters []*mysql.FilterV2
	if req.Disabled != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "organization.disabled",
			Operator: mysql.OperatorEqual,
			Value:    req.Disabled.Value,
		})
	}
	if req.Archived != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "organization.archived",
			Operator: mysql.OperatorEqual,
			Value:    req.Archived.Value,
		})
	}
	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"organization.id", "organization.name", "organization.url_code"},
			Keyword: req.SearchKeyword,
		}
	}
	orders, err := s.newOrganizationListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"failed to create OrganizationListOrders",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	options := &mysql.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		InFilters:   nil,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: searchQuery,
		Orders:      orders,
	}
	organizations, nextCursor, totalCount, err := s.orgStorage.ListOrganizations(ctx, options)
	if err != nil {
		s.logger.Error(
			"failed to list organizations",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.ListOrganizationsResponse{
		Organizations: organizations,
		Cursor:        strconv.Itoa(nextCursor),
		TotalCount:    totalCount,
	}, nil
}

func (s *EnvironmentService) CreateDemoOrganization(
	ctx context.Context,
	req *environmentproto.CreateDemoOrganizationRequest,
) (*environmentproto.CreateDemoOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if !s.opts.isDemoSiteEnabled {
		dt, err := statusDemoSiteDisabled.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.Organization),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	demoToken, ok := rpc.GetDemoCreationToken(ctx)
	if !ok {
		s.logger.Error("failed to get access demoToken",
			log.FieldsFromIncomingContext(ctx)...,
		)
		dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.UnauthenticatedError,
				"demo creation token",
			),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	editor := &eventproto.Editor{
		Email:   demoToken.Email,
		IsAdmin: false,
	}
	if err := validateCreateDemoOrganizationRequest(req, demoToken.Email, localizer); err != nil {
		s.logger.Error("failed to validate CreateDemoOrganizationRequest",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}

	organization, err := s.createOrganizationMySQL(
		ctx,
		req.Name,
		req.UrlCode,
		demoToken.Email,
		req.Description,
		false,
		false,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	// Publish the auditlog event
	event, err := domainevent.NewAdminEvent(
		editor,
		eventproto.Event_ORGANIZATION,
		organization.Id,
		eventproto.Event_DEMO_ORGANIZATION_CREATED,
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
		nil,
	)
	if err != nil {
		return nil, api.NewGRPCStatus(err).Err()
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error("failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.CreateDemoOrganizationResponse{
		Organization: organization.Organization,
	}, nil
}

func validateCreateDemoOrganizationRequest(
	req *environmentproto.CreateDemoOrganizationRequest,
	ownerEmail string,
	localizer locale.Localizer,
) error {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		dt, err := statusOrganizationNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.Name) > maxOrganizationNameLength {
		dt, err := statusInvalidOrganizationName.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}

	req.UrlCode = strings.TrimSpace(req.UrlCode)
	if !organizationUrlCodeRegex.MatchString(req.UrlCode) {
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
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	organization, err := s.createOrganizationMySQL(
		ctx,
		name,
		urlCode,
		req.OwnerEmail,
		req.Description,
		req.IsTrial,
		req.IsSystemAdmin,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	// Publish the auditlog event
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
		nil,
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

func (s *EnvironmentService) createOrganizationMySQL(
	ctx context.Context,
	name string,
	urlCode string,
	ownerEmail string,
	description string,
	isTrial bool,
	isSystemAdmin bool,
	localizer locale.Localizer,
) (*domain.Organization, error) {
	organization, err := domain.NewOrganization(
		name,
		urlCode,
		ownerEmail,
		description,
		isTrial,
		isSystemAdmin,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a domain organization",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", name),
				zap.String("urlCode", urlCode),
				zap.String("ownerEmail", ownerEmail),
				zap.Bool("isTrial", isTrial),
				zap.Bool("isSystemAdmin", isSystemAdmin),
			)...)
		return nil, api.NewGRPCStatus(err).Err()
	}
	var envRoles []*accountproto.AccountV2_EnvironmentRole
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		// Check if there is already a system admin organization
		if organization.SystemAdmin {
			org, err := s.orgStorage.GetSystemAdminOrganization(contextWithTx)
			if err != nil {
				return err
			}
			if org != nil {
				return v2es.ErrOrganizationAlreadyExists
			}
		}
		if err := s.orgStorage.CreateOrganization(contextWithTx, organization); err != nil {
			return err
		}
		// Create a default project
		project, err := s.createDefaultProject(
			contextWithTx,
			organization.Id,
			organization.OwnerEmail,
		)
		if err != nil {
			s.logger.Error(
				"Failed to create the default project",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("name", name),
					zap.String("urlCode", urlCode),
					zap.String("ownerEmail", ownerEmail),
					zap.Bool("isTrial", isTrial),
					zap.Bool("isSystemAdmin", isSystemAdmin),
				)...,
			)
			return err
		}
		// Create default environments
		envRoles, err = s.createDefaultEnvironments(
			contextWithTx,
			organization.Id,
			project,
		)
		if err != nil {
			s.logger.Error(
				"Failed to create the default environments",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("name", name),
					zap.String("urlCode", urlCode),
					zap.String("ownerEmail", ownerEmail),
					zap.Bool("isTrial", isTrial),
					zap.Bool("isSystemAdmin", isSystemAdmin),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		s.logger.Error(
			"Failed to create an organization",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", name),
				zap.String("urlCode", urlCode),
				zap.String("ownerEmail", ownerEmail),
				zap.Bool("isTrial", isTrial),
				zap.Bool("isSystemAdmin", isSystemAdmin),
			)...,
		)
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
		return nil, api.NewGRPCStatus(err).Err()
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, statusInternal.Err()
	}
	return organization, nil
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
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		if organization.SystemAdmin {
			org, err := s.orgStorage.GetSystemAdminOrganization(contextWithTx)
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
		return s.orgStorage.CreateOrganization(contextWithTx, organization)
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

// Create a default project.
// We must create a project when creating an Organization
// because we also need to create the organization owner in the account table.
// To create it we need the project, so we can also create the environment.
func (s *EnvironmentService) createDefaultProject(
	ctx context.Context,
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
	if err := s.projectStorage.CreateProject(ctx, project); err != nil {
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
	organizationID string,
	project *domain.Project,
) ([]*accountproto.AccountV2_EnvironmentRole, error) {
	envRoles := make([]*accountproto.AccountV2_EnvironmentRole, 0, 2)
	envNames := []string{
		"Development",
		"Production",
	}
	for _, name := range envNames {
		env, err := domain.NewEnvironmentV2(
			name,
			strings.ToLower(name),
			"",
			project.Id,
			organizationID,
			false,
			s.logger,
		)
		if err != nil {
			return nil, err
		}
		if err := s.environmentStorage.CreateEnvironmentV2(ctx, env); err != nil {
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
		[]string{},
		organizationID,
		accountproto.AccountV2_Role_Organization_OWNER,
		envRoles,
	)
	if err := s.accountStorage.CreateAccountV2(ctx, account); err != nil {
		return err
	}
	return nil
}

func (s *EnvironmentService) UpdateOrganization(
	ctx context.Context,
	req *environmentproto.UpdateOrganizationRequest,
) (*environmentproto.UpdateOrganizationResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(ctx, req.Id, accountproto.AccountV2_Role_Organization_OWNER, localizer)
	if err != nil {
		return nil, err
	}
	// Additional security validations for ownership transfer
	if req.OwnerEmail != nil || (req.ChangeOwnerEmailCommand != nil && req.ChangeOwnerEmailCommand.OwnerEmail != "") {
		if err := s.validateOwnershipTransfer(ctx, req, editor, localizer); err != nil {
			return nil, err
		}
	}

	commands := s.getUpdateOrganizationCommands(req)
	if len(commands) == 0 {
		return s.updateOrganizationNoCommand(ctx, req, editor, localizer)
	}

	if err := s.validateUpdateOrganizationRequest(req.Id, commands, localizer); err != nil {
		return nil, err
	}
	if err := s.updateOrganization(ctx, req.Id, editor, localizer, commands...); err != nil {
		return nil, err
	}
	return &environmentproto.UpdateOrganizationResponse{}, nil
}

func (s *EnvironmentService) updateOrganizationNoCommand(
	ctx context.Context,
	req *environmentproto.UpdateOrganizationRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*environmentproto.UpdateOrganizationResponse, error) {
	if err := s.validateUpdateOrganizationRequestNoCommand(req, localizer); err != nil {
		return nil, err
	}
	var prevOwnerEmail string
	var newOwnerEmail string
	var event *eventproto.Event
	err := s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		prevOwnerEmail = organization.OwnerEmail
		updated, err := organization.Update(
			req.Name,
			req.Description,
			req.OwnerEmail,
		)
		if err != nil {
			return err
		}
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ORGANIZATION,
			req.Id,
			eventproto.Event_ORGANIZATION_UPDATED,
			&eventproto.OrganizationUpdatedEvent{
				Id:          req.Id,
				Name:        req.Name,
				Description: req.Description,
				OwnerEmail:  req.OwnerEmail,
			},
			updated,
			organization,
		)
		if err != nil {
			return err
		}
		// Set the new owner email if it changes
		if prevOwnerEmail != updated.OwnerEmail {
			newOwnerEmail = updated.OwnerEmail
		}
		return orgStorage.UpdateOrganization(ctxWithTx, updated)
	})
	if err != nil {
		return nil, s.reportUpdateOrganizationError(ctx, err, localizer)
	}

	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish the event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	// Update the organization role when the owner email changes
	if prevOwnerEmail != "" && newOwnerEmail != "" {
		if err := s.updateOwnerRole(ctx, req.Id, prevOwnerEmail, newOwnerEmail); err != nil {
			s.logger.Error("Failed to update the new owner's role",
				zap.Error(err),
				zap.String("organizationId", req.Id),
				zap.String("prevOwnerEmail", prevOwnerEmail),
				zap.String("newOwnerEmail", newOwnerEmail),
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}

	return &environmentproto.UpdateOrganizationResponse{}, nil
}

func (s *EnvironmentService) reportUpdateOrganizationError(
	ctx context.Context,
	err error,
	localizer locale.Localizer,
) error {
	s.logger.Error(
		"Failed to update organization",
		log.FieldsFromIncomingContext(ctx).AddFields(
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
	return api.NewGRPCStatus(err).Err()
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

func (s *EnvironmentService) validateUpdateOrganizationRequestNoCommand(
	req *environmentproto.UpdateOrganizationRequest,
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
	if req.Name != nil {
		name := strings.TrimSpace(req.Name.Value)
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
	return nil
}

func (s *EnvironmentService) updateOrganization(
	ctx context.Context,
	id string,
	editor *eventproto.Editor,
	localizer locale.Localizer,
	commands ...command.Command,
) error {
	var prevOwnerEmail string
	var newOwnerEmail string
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		organization, err := s.orgStorage.GetOrganization(contextWithTx, id)
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
		return s.orgStorage.UpdateOrganization(contextWithTx, organization)
	})
	if err != nil {
		return s.reportUpdateOrganizationError(ctx, err, localizer)
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
			return api.NewGRPCStatus(err).Err()
		}
	}
	return nil
}

// validateOwnershipTransfer performs additional security validations for ownership transfer
func (s *EnvironmentService) validateOwnershipTransfer(
	ctx context.Context,
	req *environmentproto.UpdateOrganizationRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) error {
	// Get current organization to validate against
	organization, err := s.orgStorage.GetOrganization(ctx, req.Id)
	if err != nil {
		return err
	}

	// Determine the new owner email being requested
	var newOwnerEmail string
	if req.OwnerEmail != nil {
		newOwnerEmail = req.OwnerEmail.Value
	} else if req.ChangeOwnerEmailCommand != nil {
		newOwnerEmail = req.ChangeOwnerEmailCommand.OwnerEmail
	}

	// Don't allow no-op updates (setting same owner)
	if newOwnerEmail == organization.OwnerEmail {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale: localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(
				locale.InvalidArgumentError,
				"new owner email is the same as the current owner",
			),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}

	// If not system admin, ensure current user is actually the current owner
	if !editor.IsAdmin && editor.Email != organization.OwnerEmail {
		dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}

	// New owner must exist and be a member of the organization
	newOwnerAccount, err := s.accountStorage.GetAccountV2(ctx, newOwnerEmail, req.Id)
	if err != nil {
		if errors.Is(err, v2acc.ErrAccountNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "new owner account not found in organization"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		return err
	}

	// New owner account must be enabled
	if newOwnerAccount.Disabled {
		dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "new owner account is disabled"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}

	return nil
}

func (s *EnvironmentService) updateOwnerRole(
	ctx context.Context,
	organizationID, prevOwnerEmail, newOwnerEmail string,
) error {
	// Update the old owner organization role
	prevOwnerAcc, err := s.accountStorage.GetAccountV2(ctx, prevOwnerEmail, organizationID)
	if err != nil {
		return err
	}
	if err := prevOwnerAcc.ChangeOrganizationRole(accountproto.AccountV2_Role_Organization_ADMIN); err != nil {
		return err
	}
	if err := s.accountStorage.UpdateAccountV2(ctx, prevOwnerAcc); err != nil {
		return err
	}
	// Update the new owner organization role
	newOwnerAcc, err := s.accountStorage.GetAccountV2(ctx, newOwnerEmail, organizationID)
	if err != nil {
		return err
	}
	if err := newOwnerAcc.ChangeOrganizationRole(accountproto.AccountV2_Role_Organization_OWNER); err != nil {
		return err
	}
	if err := s.accountStorage.UpdateAccountV2(ctx, newOwnerAcc); err != nil {
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

func (s *EnvironmentService) DeleteOrganizationData(
	ctx context.Context,
	req *environmentproto.DeleteOrganizationDataRequest,
) (*environmentproto.DeleteOrganizationDataResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if len(req.OrganizationIds) == 0 {
		return nil, statusOrganizationIDRequired.Err()
	}

	// 1. Get all environments for the organization IDs
	inFilters := []*mysql.InFilter{
		{
			Column: "environment_v2.organization_id",
			Values: convToInterfaceSlice(req.OrganizationIds),
		},
	}

	options := &mysql.ListOptions{
		Limit:       mysql.QueryNoLimit,
		Offset:      mysql.QueryNoOffset,
		Filters:     nil,
		InFilters:   inFilters,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	environments, _, _, err := s.environmentStorage.ListEnvironmentsV2(ctx, options)
	if err != nil {
		s.logger.Error("Could not list environments", zap.Error(err))
		return nil, err
	}
	environmentIDs := make([]string, 0, len(environments))
	for _, environment := range environments {
		environmentIDs = append(environmentIDs, environment.Id)
	}

	deletionSummary := make([]*environmentproto.OrganizationDeletionSummary, 0, len(req.OrganizationIds))
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		// 2. Count all target entities in the organizations
		for _, orgID := range req.OrganizationIds {
			targetInOrg, err := s.countTargetDeletionEntitiesInOrganizations(ctxWithTx, orgID)
			if err != nil {
				s.logger.Error("Could not count target entities in organization",
					zap.Error(err),
					zap.String("organizationID", orgID),
				)
				return err
			}
			deletionSummary = append(deletionSummary, targetInOrg)
		}
		if req.DryRun && !req.Force {
			return nil
		}

		// 3. safety checks if features receiving request
		if !req.Force {
			isRunning, err := s.checkRunningFeaturesInEnvironments(ctxWithTx, environmentIDs)
			if err != nil {
				s.logger.Error("Could not check running features in environments",
					zap.Error(err),
					zap.Strings("organizationIDs", req.OrganizationIds),
					zap.Strings("environmentIDs", environmentIDs),
				)
				return err
			}
			if isRunning {
				return statusCannotDeleteOrganization.Err()
			}
		}

		// 4. Delete all target entities from the environments
		for _, environment := range environments {
			for _, target := range targetEntitiesInEnvironment {
				err := s.environmentStorage.DeleteTargetFromEnvironmentV2(ctxWithTx, environment.Id, target)
				if err != nil {
					s.logger.Error("Could not delete target from environment",
						zap.Error(err),
						zap.String("environmentID", environment.Id),
						zap.String("target", target),
						zap.Strings("organizationIDs", req.OrganizationIds),
					)
					return err
				}
			}
		}

		// 5. Delete all environments from organizations
		whereParts := []mysql.WherePart{
			mysql.NewInFilter("organization_id", convToInterfaceSlice(req.OrganizationIds)),
		}
		err = s.environmentStorage.DeleteEnvironmentV2(ctxWithTx, whereParts)
		if err != nil {
			s.logger.Error("Could not delete environment",
				zap.Error(err),
				zap.Strings("organizationIDs", req.OrganizationIds),
			)
			return err
		}

		// 6. Delete all projects from organizations
		whereParts = []mysql.WherePart{
			mysql.NewInFilter("organization_id", convToInterfaceSlice(req.OrganizationIds)),
		}
		err = s.projectStorage.DeleteProjects(ctxWithTx, whereParts)
		if err != nil {
			s.logger.Error("Could not delete projects",
				zap.Error(err),
				zap.Strings("organizationIDs", req.OrganizationIds),
			)
			return err
		}

		// 7. Delete all organizations with its data
		err = s.deleteOrganizations(ctxWithTx, req.OrganizationIds)
		if err != nil {
			s.logger.Error("Failed to delete organizations",
				zap.Error(err),
				zap.Strings("organizationIDs", req.OrganizationIds),
			)
			return err
		}
		return nil
	})
	if err != nil {
		s.logger.Error("Failed to delete organizations",
			zap.Error(err),
			zap.Strings("organizationIDs", req.OrganizationIds),
		)
		return nil, err
	}
	return &environmentproto.DeleteOrganizationDataResponse{
		Summaries: deletionSummary,
	}, nil
}

func (s *EnvironmentService) checkRunningFeaturesInEnvironments(
	ctx context.Context,
	environmentIDs []string,
) (bool, error) {
	listOptions := &mysql.ListOptions{
		Limit:  mysql.QueryNoLimit,
		Offset: mysql.QueryNoOffset,
		Filters: []*mysql.FilterV2{
			{
				Column:   "last_used_at",
				Operator: mysql.OperatorGreaterThanOrEqual,
				Value:    time.Now().Add(-featureActiveDays).Unix(),
			},
		},
		InFilters: []*mysql.InFilter{
			{
				Column: "environment_id",
				Values: convToInterfaceSlice(environmentIDs),
			},
		},
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	flui, err := s.fluiStorage.SelectFeatureLastUsedInfos(ctx, listOptions)
	if err != nil {
		s.logger.Error("Could not list feature last used infos",
			zap.Error(err),
			zap.Strings("environmentIDs", environmentIDs),
		)
		return false, err
	}
	return len(flui) > 0, nil
}

func (s *EnvironmentService) countTargetDeletionEntitiesInOrganizations(
	ctx context.Context,
	organizationID string,
) (*environmentproto.OrganizationDeletionSummary, error) {
	listEnvironmentOptions := &mysql.ListOptions{
		Limit:  mysql.QueryNoLimit,
		Offset: mysql.QueryNoOffset,
		Filters: []*mysql.FilterV2{
			{
				Column:   "organization_id",
				Operator: mysql.OperatorEqual,
				Value:    organizationID,
			},
		},
		InFilters:   nil,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	_, _, environmentsCount, err := s.environmentStorage.ListEnvironmentsV2(ctx, listEnvironmentOptions)
	if err != nil {
		s.logger.Error("Failed to count environments in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}
	listProjectsOptions := &mysql.ListOptions{
		Limit:  mysql.QueryNoLimit,
		Offset: mysql.QueryNoOffset,
		Filters: []*mysql.FilterV2{
			{
				Column:   "project.organization_id",
				Operator: mysql.OperatorEqual,
				Value:    organizationID,
			},
		},
		InFilters:   nil,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	_, _, projectsCount, err := s.projectStorage.ListProjects(ctx, listProjectsOptions)
	if err != nil {
		s.logger.Error("Failed to count projects in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}
	listAccountOptions := &mysql.ListOptions{
		Limit:  mysql.QueryNoLimit,
		Offset: mysql.QueryNoOffset,
		Filters: []*mysql.FilterV2{
			{
				Column:   "organization_id",
				Operator: mysql.OperatorEqual,
				Value:    organizationID,
			},
		},
		InFilters:   nil,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	_, _, accountsCount, err := s.accountStorage.ListAccountsV2(ctx, listAccountOptions)
	if err != nil {
		s.logger.Error("Failed to count account_v2 entities in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	listTeamsOptions := &mysql.ListOptions{
		Limit:  mysql.QueryNoLimit,
		Offset: mysql.QueryNoOffset,
		Filters: []*mysql.FilterV2{
			{
				Column:   "organization_id",
				Operator: mysql.OperatorEqual,
				Value:    organizationID,
			},
		},
		InFilters:   nil,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	_, _, teamsCount, err := s.teamStorage.ListTeams(ctx, listTeamsOptions)
	if err != nil {
		s.logger.Error("Failed to count teams in environment",
			zap.Error(err),
			zap.String("organization_id", organizationID),
		)
		return nil, err
	}

	featuresCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "feature")
	if err != nil {
		s.logger.Error("Failed to count feature entities in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	experimentsCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "experiment")
	if err != nil {
		s.logger.Error("Failed to count experiment entities in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	subscriptionsCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "subscription")
	if err != nil {
		s.logger.Error("Failed to count subscriptions in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	pushesCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "push")
	if err != nil {
		s.logger.Error("Failed to count pushes in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	tagsCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "tag")
	if err != nil {
		s.logger.Error("Failed to count tags in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	segmentsCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "segment")
	if err != nil {
		s.logger.Error("Failed to count segments in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	flagTriggersCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "flag_trigger")
	if err != nil {
		s.logger.Error("Failed to count flag triggers in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	apiKeysCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "api_key")
	if err != nil {
		s.logger.Error("Failed to count api keys in organization",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	operationsCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "auto_ops_rule")
	if err != nil {
		s.logger.Error("Failed to count operations in environment",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	featureLastUsedInfosCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(
		ctx, organizationID, "feature_last_used_info",
	)
	if err != nil {
		s.logger.Error("Failed to count feature last used infos in environment",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	goalsCount, err := s.orgStorage.CountEnvTargetEntitiesInOrganization(ctx, organizationID, "goal")
	if err != nil {
		s.logger.Error("Failed to count goals in environment",
			zap.Error(err),
			zap.String("organizationID", organizationID),
		)
		return nil, err
	}

	return &environmentproto.OrganizationDeletionSummary{
		OrganizationId:              organizationID,
		EnvironmentsDeleted:         environmentsCount,
		ProjectsDeleted:             projectsCount,
		FeaturesDeleted:             featuresCount,
		AccountsDeleted:             accountsCount,
		ExperimentsDeleted:          experimentsCount,
		SubscriptionsDeleted:        subscriptionsCount,
		PushesDeleted:               pushesCount,
		TagsDeleted:                 tagsCount,
		TeamsDeleted:                teamsCount,
		SegmentsDeleted:             segmentsCount,
		FlagTriggersDeleted:         flagTriggersCount,
		ApiKeysDeleted:              apiKeysCount,
		OperationsDeleted:           operationsCount,
		FeatureLastUsedInfosDeleted: featureLastUsedInfosCount,
		GoalsDeleted:                goalsCount,
	}, nil
}

// deleteOrganizations deletes organizations and their related data from various target entities.
func (s *EnvironmentService) deleteOrganizations(ctx context.Context, organizationIDs []string) error {
	whereParts := []mysql.WherePart{
		mysql.NewInFilter("organization_id", convToInterfaceSlice(organizationIDs)),
	}
	for _, target := range targetEntitiesInOrganization {
		err := s.orgStorage.DeleteOrganizationData(ctx, target, whereParts)
		if err != nil {
			s.logger.Error("Failed to delete data from organization",
				zap.Error(err),
				zap.Strings("organizationIDs", organizationIDs),
			)
			return err
		}
	}
	whereParts = []mysql.WherePart{
		mysql.NewInFilter("id", convToInterfaceSlice(organizationIDs)),
	}
	err := s.orgStorage.DeleteOrganizations(ctx, whereParts)
	if err != nil {
		s.logger.Error("Failed to delete organizations",
			zap.Error(err),
			zap.Strings("organizationIDs", organizationIDs),
		)
		return err
	}
	return nil
}
