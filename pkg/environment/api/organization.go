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

	accdomain "github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2acc "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/environment/storage/v2"
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
)

func (s *EnvironmentService) GetOrganization(
	ctx context.Context,
	req *environmentproto.GetOrganizationRequest,
) (*environmentproto.GetOrganizationResponse, error) {
	_, err := s.checkOrganizationRole(ctx, req.Id, accountproto.AccountV2_Role_Organization_MEMBER)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetOrganizationRequest(req); err != nil {
		return nil, err
	}
	org, err := s.getOrganization(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &environmentproto.GetOrganizationResponse{
		Organization: org.Organization,
	}, nil
}

func (s *EnvironmentService) validateGetOrganizationRequest(
	req *environmentproto.GetOrganizationRequest,
) error {
	if req.Id == "" {
		return statusOrganizationIDRequired.Err()
	}
	return nil
}

func (s *EnvironmentService) getOrganization(
	ctx context.Context,
	id string,
) (*domain.Organization, error) {
	org, err := s.orgStorage.GetOrganization(ctx, id)
	if err != nil {
		s.logger.Error("failed to get organization",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("organizationId", id),
				zap.Error(err),
			)...)
		if errors.Is(err, v2es.ErrOrganizationNotFound) {
			return nil, statusOrganizationNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	return org, nil
}

func (s *EnvironmentService) ListOrganizations(
	ctx context.Context,
	req *environmentproto.ListOrganizationsRequest,
) (*environmentproto.ListOrganizationsResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
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
	orders, err := s.newOrganizationListOrders(req.OrderBy, req.OrderDirection)
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
		return nil, statusInvalidCursor.Err()
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
	if !s.opts.isDemoSiteEnabled {
		return nil, statusDemoSiteDisabled.Err()
	}
	demoToken, ok := rpc.GetDemoCreationToken(ctx)
	if !ok {
		s.logger.Error("failed to get access demoToken",
			log.FieldsFromIncomingContext(ctx)...,
		)
		return nil, statusUnauthenticated.Err()
	}
	editor := &eventproto.Editor{
		Email:   demoToken.Email,
		IsAdmin: false,
	}
	if err := validateCreateDemoOrganizationRequest(req, demoToken.Email); err != nil {
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
) error {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return statusOrganizationNameRequired.Err()
	}
	if len(req.Name) > maxOrganizationNameLength {
		return statusInvalidOrganizationName.Err()
	}

	req.UrlCode = strings.TrimSpace(req.UrlCode)
	if !organizationUrlCodeRegex.MatchString(req.UrlCode) {
		return statusInvalidOrganizationUrlCode.Err()
	}
	return nil
}

func (s *EnvironmentService) newOrganizationListOrders(
	orderBy environmentproto.ListOrganizationsRequest_OrderBy,
	orderDirection environmentproto.ListOrganizationsRequest_OrderDirection,
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
		return nil, statusInvalidOrderBy.Err()
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
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateOrganizationRequest(req); err != nil {
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

func (s *EnvironmentService) validateCreateOrganizationRequest(
	req *environmentproto.CreateOrganizationRequest,
) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return statusOrganizationNameRequired.Err()
	}
	if len(name) > maxOrganizationNameLength {
		return statusInvalidOrganizationName.Err()
	}
	urlCode := strings.TrimSpace(req.UrlCode)
	if !organizationUrlCodeRegex.MatchString(urlCode) {
		return statusInvalidOrganizationUrlCode.Err()
	}
	if !emailRegex.MatchString(req.OwnerEmail) {
		return statusInvalidOrganizationCreatorEmail.Err()
	}
	return nil
}

func (s *EnvironmentService) createOrganizationMySQL(
	ctx context.Context,
	name string,
	urlCode string,
	ownerEmail string,
	description string,
	isTrial bool,
	isSystemAdmin bool,
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
			return nil, statusOrganizationAlreadyExists.Err()
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
	editor, err := s.checkOrganizationRole(ctx, req.Id, accountproto.AccountV2_Role_Organization_OWNER)
	if err != nil {
		return nil, err
	}
	// Additional security validations for ownership transfer
	if req.OwnerEmail != nil {
		if err := s.validateOwnershipTransfer(ctx, req, editor); err != nil {
			return nil, err
		}
	}
	if err := s.validateUpdateOrganizationRequest(req); err != nil {
		return nil, err
	}
	var prevOwnerEmail string
	var newOwnerEmail string
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
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
		return nil, s.reportUpdateOrganizationError(ctx, err)
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
) error {
	s.logger.Error(
		"Failed to update organization",
		log.FieldsFromIncomingContext(ctx).AddFields(
			zap.Error(err),
		)...,
	)
	if errors.Is(err, domain.ErrCannotArchiveSystemAdmin) || errors.Is(err, domain.ErrCannotDisableSystemAdmin) {
		return statusCannotUpdateSystemAdmin.Err()
	}
	if errors.Is(err, v2es.ErrOrganizationNotFound) || errors.Is(err, v2es.ErrOrganizationUnexpectedAffectedRows) {
		return statusOrganizationNotFound.Err()
	}
	return api.NewGRPCStatus(err).Err()
}

func (s *EnvironmentService) validateUpdateOrganizationRequest(
	req *environmentproto.UpdateOrganizationRequest,
) error {
	if req.Id == "" {
		return statusOrganizationIDRequired.Err()
	}
	if req.Name != nil {
		name := strings.TrimSpace(req.Name.Value)
		if name == "" {
			return statusOrganizationNameRequired.Err()
		}
		if len(name) > maxOrganizationNameLength {
			return statusInvalidOrganizationName.Err()
		}
	}
	return nil
}

// validateOwnershipTransfer performs additional security validations for ownership transfer
func (s *EnvironmentService) validateOwnershipTransfer(
	ctx context.Context,
	req *environmentproto.UpdateOrganizationRequest,
	editor *eventproto.Editor,
) error {
	// Get current organization to validate against
	organization, err := s.orgStorage.GetOrganization(ctx, req.Id)
	if err != nil {
		return err
	}

	newOwnerEmail := req.OwnerEmail.Value

	// Don't allow no-op updates (setting same owner)
	if newOwnerEmail == organization.OwnerEmail {
		return statusNoCommand.Err()
	}

	// If not system admin, ensure current user is actually the current owner
	if !editor.IsAdmin && editor.Email != organization.OwnerEmail {
		return statusPermissionDenied.Err()
	}

	// New owner must exist and be a member of the organization
	newOwnerAccount, err := s.accountStorage.GetAccountV2(ctx, newOwnerEmail, req.Id)
	if err != nil {
		if errors.Is(err, v2acc.ErrAccountNotFound) {
			return statusAccountNotFound.Err()
		}
		return err
	}

	// New owner account must be enabled
	if newOwnerAccount.Disabled {
		return statusPermissionDenied.Err()
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
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateEnableOrganizationRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		prev := &domain.Organization{}
		if err := copier.Copy(prev, organization); err != nil {
			return err
		}
		organization.Enable()
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ORGANIZATION,
			organization.Id,
			eventproto.Event_ORGANIZATION_ENABLED,
			&eventproto.OrganizationEnabledEvent{
				Id: organization.Id,
			},
			organization,
			prev,
		)
		if err != nil {
			return err
		}
		return orgStorage.UpdateOrganization(ctxWithTx, organization)
	})
	if err != nil {
		return nil, s.reportUpdateOrganizationError(ctx, err)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish enable organization event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("event", event),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.EnableOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateEnableOrganizationRequest(
	req *environmentproto.EnableOrganizationRequest,
) error {
	if req.Id == "" {
		return statusOrganizationIDRequired.Err()
	}
	return nil
}

func (s *EnvironmentService) DisableOrganization(
	ctx context.Context,
	req *environmentproto.DisableOrganizationRequest,
) (*environmentproto.DisableOrganizationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateDisableOrganizationRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		prev := &domain.Organization{}
		if err := copier.Copy(prev, organization); err != nil {
			return err
		}
		if err := organization.Disable(); err != nil {
			return err
		}
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ORGANIZATION,
			organization.Id,
			eventproto.Event_ORGANIZATION_DISABLED,
			&eventproto.OrganizationDisabledEvent{
				Id: organization.Id,
			},
			organization,
			prev,
		)
		if err != nil {
			return err
		}
		return orgStorage.UpdateOrganization(ctxWithTx, organization)
	})
	if err != nil {
		return nil, s.reportUpdateOrganizationError(ctx, err)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish disable organization event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("event", event),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.DisableOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateDisableOrganizationRequest(
	req *environmentproto.DisableOrganizationRequest,
) error {
	if req.Id == "" {
		return statusOrganizationIDRequired.Err()
	}
	return nil
}

func (s *EnvironmentService) ArchiveOrganization(
	ctx context.Context,
	req *environmentproto.ArchiveOrganizationRequest,
) (*environmentproto.ArchiveOrganizationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateArchiveOrganizationRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		prev := &domain.Organization{}
		if err := copier.Copy(prev, organization); err != nil {
			return err
		}
		if err := organization.Archive(); err != nil {
			return err
		}
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ORGANIZATION,
			organization.Id,
			eventproto.Event_ORGANIZATION_ARCHIVED,
			&eventproto.OrganizationArchivedEvent{
				Id: organization.Id,
			},
			organization,
			prev,
		)
		if err != nil {
			return err
		}
		return orgStorage.UpdateOrganization(ctxWithTx, organization)
	})
	if err != nil {
		return nil, s.reportUpdateOrganizationError(ctx, err)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish archive organization event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("event", event),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.ArchiveOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateArchiveOrganizationRequest(
	req *environmentproto.ArchiveOrganizationRequest,
) error {
	if req.Id == "" {
		return statusOrganizationIDRequired.Err()
	}
	return nil
}

func (s *EnvironmentService) UnarchiveOrganization(
	ctx context.Context,
	req *environmentproto.UnarchiveOrganizationRequest,
) (*environmentproto.UnarchiveOrganizationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateUnarchiveOrganizationRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		prev := &domain.Organization{}
		if err := copier.Copy(prev, organization); err != nil {
			return err
		}
		organization.Unarchive()
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ORGANIZATION,
			organization.Id,
			eventproto.Event_ORGANIZATION_UNARCHIVED,
			&eventproto.OrganizationUnarchivedEvent{
				Id: organization.Id,
			},
			organization,
			prev,
		)
		if err != nil {
			return err
		}
		return orgStorage.UpdateOrganization(ctxWithTx, organization)
	})
	if err != nil {
		return nil, s.reportUpdateOrganizationError(ctx, err)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish unarchive organization event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("event", event),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.UnarchiveOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateUnarchiveOrganizationRequest(
	req *environmentproto.UnarchiveOrganizationRequest,
) error {
	if req.Id == "" {
		return statusOrganizationIDRequired.Err()
	}
	return nil
}

func (s *EnvironmentService) ConvertTrialOrganization(
	ctx context.Context,
	req *environmentproto.ConvertTrialOrganizationRequest,
) (*environmentproto.ConvertTrialOrganizationResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateConvertTrialOrganizationRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		orgStorage := v2es.NewOrganizationStorage(tx)
		organization, err := orgStorage.GetOrganization(ctxWithTx, req.Id)
		if err != nil {
			return err
		}
		prev := &domain.Organization{}
		if err := copier.Copy(prev, organization); err != nil {
			return err
		}
		organization.ConvertTrial()
		event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ORGANIZATION,
			organization.Id,
			eventproto.Event_ORGANIZATION_TRIAL_CONVERTED,
			&eventproto.OrganizationTrialConvertedEvent{
				Id: organization.Id,
			},
			organization,
			prev,
		)
		if err != nil {
			return err
		}
		return orgStorage.UpdateOrganization(ctxWithTx, organization)
	})
	if err != nil {
		return nil, s.reportUpdateOrganizationError(ctx, err)
	}
	if err = s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish convert trial organization event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("event", event),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &environmentproto.ConvertTrialOrganizationResponse{}, nil
}

func (s *EnvironmentService) validateConvertTrialOrganizationRequest(
	req *environmentproto.ConvertTrialOrganizationRequest,
) error {
	if req.Id == "" {
		return statusOrganizationIDRequired.Err()
	}
	return nil
}
