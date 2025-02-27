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

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	maxProjectNameLength = 50
	projectUrlCodeRegex  = regexp.MustCompile("^[a-z0-9-]{1,50}$")

	//nolint:lll
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func (s *EnvironmentService) GetProject(
	ctx context.Context,
	req *environmentproto.GetProjectRequest,
) (*environmentproto.GetProjectResponse, error) {
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
	if err := validateGetProjectRequest(req, localizer); err != nil {
		return nil, err
	}
	project, err := s.getProject(ctx, req.Id, localizer)
	if err != nil {
		return nil, err
	}
	return &environmentproto.GetProjectResponse{
		Project: project.Project,
	}, nil
}

func validateGetProjectRequest(req *environmentproto.GetProjectRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusProjectIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) getProject(
	ctx context.Context,
	id string,
	localizer locale.Localizer,
) (*domain.Project, error) {
	projectStorage := v2es.NewProjectStorage(s.mysqlClient)
	project, err := projectStorage.GetProject(ctx, id)
	if err != nil {
		if err == v2es.ErrProjectNotFound {
			dt, err := statusProjectNotFound.WithDetails(&errdetails.LocalizedMessage{
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
	return project, nil
}

func (s *EnvironmentService) ListProjects(
	ctx context.Context,
	req *environmentproto.ListProjectsRequest,
) (*environmentproto.ListProjectsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{}
	if len(req.OrganizationIds) > 0 {
		oIDs := convToInterfaceSlice(req.OrganizationIds)
		whereParts = append(whereParts, mysql.NewInFilter("project.organization_id", oIDs))
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("project.disabled", "=", req.Disabled.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(
			whereParts,
			mysql.NewSearchQuery(
				[]string{"project.id", "project.name", "project.url_code", "project.creator_email"},
				req.SearchKeyword,
			),
		)
	}
	orders, err := s.newProjectListOrders(req.OrderBy, req.OrderDirection, localizer)
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
	projectStorage := v2es.NewProjectStorage(s.mysqlClient)
	projects, nextCursor, totalCount, err := projectStorage.ListProjects(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list projects",
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
	return &environmentproto.ListProjectsResponse{
		Projects:   projects,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func convToInterfaceSlice(
	slice []string,
) []interface{} {
	result := make([]interface{}, 0, len(slice))
	for _, element := range slice {
		result = append(result, element)
	}
	return result
}

func (s *EnvironmentService) newProjectListOrders(
	orderBy environmentproto.ListProjectsRequest_OrderBy,
	orderDirection environmentproto.ListProjectsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case environmentproto.ListProjectsRequest_DEFAULT,
		environmentproto.ListProjectsRequest_NAME:
		column = "project.name"
	case environmentproto.ListProjectsRequest_URL_CODE:
		column = "project.url_code"
	case environmentproto.ListProjectsRequest_ID:
		column = "project.id"
	case environmentproto.ListProjectsRequest_CREATED_AT:
		column = "project.created_at"
	case environmentproto.ListProjectsRequest_UPDATED_AT:
		column = "project.updated_at"
	case environmentproto.ListProjectsRequest_ENVIRONMENT_COUNT:
		column = "environment_count"
	case environmentproto.ListProjectsRequest_FEATURE_COUNT:
		column = "feature_count"
	case environmentproto.ListProjectsRequest_CREATOR_EMAIL:
		column = "project.creator_email"
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
	if orderDirection == environmentproto.ListProjectsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *EnvironmentService) CreateProject(
	ctx context.Context,
	req *environmentproto.CreateProjectRequest,
) (*environmentproto.CreateProjectResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		req.OrganizationId,
		accountproto.AccountV2_Role_Organization_ADMIN,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createProjectNoCommand(ctx, req, localizer, editor)
	}
	if err := validateCreateProjectRequest(req, localizer); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Command.Name)
	urlCode := strings.TrimSpace(req.Command.UrlCode)
	project, err := domain.NewProject(
		name,
		urlCode,
		req.Command.Description,
		editor.Email,
		req.Command.OrganizationId,
		false,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create project",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if err := s.createProject(ctx, req.Command, project, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.CreateProjectResponse{Project: project.Project}, nil
}

func (s *EnvironmentService) createProjectNoCommand(
	ctx context.Context,
	req *environmentproto.CreateProjectRequest,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) (*environmentproto.CreateProjectResponse, error) {
	if err := validateCreateProjectRequestNoCommand(req, localizer); err != nil {
		return nil, err
	}
	newProj, err := domain.NewProject(
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.UrlCode),
		req.Description,
		editor.Email,
		req.OrganizationId,
		false,
	)
	if err != nil {
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		storage := v2es.NewProjectStorage(s.mysqlClient)
		e, err := domainevent.NewAdminEvent(
			editor,
			eventproto.Event_PROJECT,
			newProj.Id,
			eventproto.Event_PROJECT_CREATED,
			&eventproto.ProjectCreatedEvent{
				Id:          newProj.Id,
				Name:        newProj.Name,
				UrlCode:     newProj.UrlCode,
				Description: newProj.Description,
				Trial:       false,
				CreatedAt:   newProj.CreatedAt,
				UpdatedAt:   newProj.UpdatedAt,
			},
			newProj.Project,
			&domain.Project{},
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, e); err != nil {
			return err
		}
		return storage.CreateProject(ctxWithTx, newProj)
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
			"Failed to create project",
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
	return &environmentproto.CreateProjectResponse{
		Project: newProj.Project,
	}, nil
}

func validateCreateProjectRequestNoCommand(
	req *environmentproto.CreateProjectRequest,
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
	if req.OrganizationId == "" {
		dt, err := statusProjectIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateCreateProjectRequest(req *environmentproto.CreateProjectRequest, localizer locale.Localizer) error {
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
		dt, err := statusProjectNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(name) > maxProjectNameLength {
		dt, err := statusInvalidProjectName.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	urlCode := strings.TrimSpace(req.Command.UrlCode)
	if !projectUrlCodeRegex.MatchString(urlCode) {
		dt, err := statusInvalidProjectUrlCode.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) createProject(
	ctx context.Context,
	cmd command.Command,
	project *domain.Project,
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
		projectStorage := v2es.NewProjectStorage(tx)
		handler, err := command.NewProjectCommandHandler(editor, project, s.publisher)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return projectStorage.CreateProject(ctx, project)
	})
	if err != nil {
		if err == v2es.ErrProjectAlreadyExists {
			dt, err := statusProjectAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		s.logger.Error(
			"Failed to create project",
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

func (s *EnvironmentService) CreateTrialProject(
	ctx context.Context,
	req *environmentproto.CreateTrialProjectRequest,
) (*environmentproto.CreateTrialProjectResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateCreateTrialProjectRequest(req, localizer); err != nil {
		return nil, err
	}
	editor := &eventproto.Editor{
		Email:   req.Command.Email,
		IsAdmin: false,
	}
	existingProject, err := s.getTrialProjectByEmail(ctx, editor.Email, localizer)
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	if existingProject != nil {
		dt, err := statusProjectAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AlreadyExistsError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	// TODO Once we support new create project API requiring name instead of id, we should remove this process.
	name := strings.TrimSpace(req.Command.Name)
	if req.Command.Name == "" {
		name = req.Command.Id
	}
	// TODO Once we support new create project API requiring urlCode instead of id, we should remove this process.
	urlCode := name
	if req.Command.UrlCode != "" {
		urlCode = req.Command.UrlCode
	}
	// TODO: Temporary implementations that create Organization at the same time as Project.
	// This should be removed when the Organization management page is added.
	organization, err := domain.NewOrganization(
		name,
		urlCode,
		req.Command.OwnerEmail,
		"",
		true,
		false,
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
	createOrgCmd := &environmentproto.CreateOrganizationCommand{
		Name:        organization.Name,
		UrlCode:     organization.UrlCode,
		Description: organization.Description,
		IsTrial:     true,
	}
	if err = s.createOrganization(ctx, createOrgCmd, organization, editor, localizer); err != nil {
		s.logger.Error(
			"Failed to save organization",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	s.logger.Info(
		`Trial organization is created at the same time as Project.
This is a temporary implementation during the transition period.`,
		zap.String("organization_id", organization.Id),
		zap.String("organization_name", organization.Name),
		zap.String("organization_url_code", organization.UrlCode),
	)
	project, err := domain.NewProject(name, urlCode, "", editor.Email, organization.Id, true)
	if err != nil {
		s.logger.Error(
			"Failed to create trial project",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if err := s.createProject(ctx, req.Command, project, editor, localizer); err != nil {
		return nil, err
	}
	if err := s.createTrialEnvironmentsAndAccounts(ctx, project, editor, localizer); err != nil {
		return nil, err
	}
	return &environmentproto.CreateTrialProjectResponse{}, nil
}

func validateCreateTrialProjectRequest(
	req *environmentproto.CreateTrialProjectRequest,
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
	// TODO Once we support new create project API requiring name instead of id, we should validate name using regex.
	name := strings.TrimSpace(req.Command.Name)
	if name == "" {
		dt, err := statusProjectNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(name) > maxProjectNameLength {
		dt, err := statusInvalidProjectName.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.UrlCode != "" && !projectUrlCodeRegex.MatchString(req.Command.UrlCode) {
		dt, err := statusInvalidProjectUrlCode.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "url_code"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if !emailRegex.MatchString(req.Command.Email) {
		dt, err := statusInvalidProjectCreatorEmail.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) getTrialProjectByEmail(
	ctx context.Context,
	email string,
	localizer locale.Localizer,
) (*environmentproto.Project, error) {
	projectStorage := v2es.NewProjectStorage(s.mysqlClient)
	project, err := projectStorage.GetTrialProjectByEmail(ctx, email, false, true)
	if err != nil {
		if err == v2es.ErrProjectNotFound {
			dt, err := statusProjectNotFound.WithDetails(&errdetails.LocalizedMessage{
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
	return project.Project, nil
}

func (s *EnvironmentService) createTrialEnvironmentsAndAccounts(
	ctx context.Context,
	project *domain.Project,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) error {
	envRoles := make([]*accountproto.AccountV2_EnvironmentRole, 0, 2)
	envNames := []string{
		"Development",
		"Production",
	}
	for _, name := range envNames {
		envURLCode := fmt.Sprintf("%s-%s", project.UrlCode, strings.ToLower(name))
		createEnvCmdV2 := &environmentproto.CreateEnvironmentV2Command{
			Name:        name,
			UrlCode:     envURLCode,
			ProjectId:   project.Id,
			Description: "",
		}
		envV2, err := domain.NewEnvironmentV2(name, envURLCode, "", project.Id, project.OrganizationId, false, s.logger)
		if err != nil {
			return err
		}
		if err := s.createEnvironmentV2(ctx, createEnvCmdV2, envV2, editor, localizer); err != nil {
			return err
		}
		envRoles = append(envRoles, &accountproto.AccountV2_EnvironmentRole{
			EnvironmentId: envV2.Id,
			Role:          accountproto.AccountV2_Role_Environment_EDITOR,
		})
	}
	createAccountReq := &accountproto.CreateAccountV2Request{
		OrganizationId: project.OrganizationId,
		Command: &accountproto.CreateAccountV2Command{
			Email:          editor.Email,
			Name:           strings.Split(editor.Email, "@")[0],
			AvatarImageUrl: "",
			// TODO Once we support new console design, we should set OWNER role.
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			EnvironmentRoles: envRoles,
		},
	}
	if _, err := s.accountClient.CreateAccountV2(ctx, createAccountReq); err != nil {
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

func (s *EnvironmentService) UpdateProject(
	ctx context.Context,
	req *environmentproto.UpdateProjectRequest,
) (*environmentproto.UpdateProjectResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	commands := getUpdateProjectCommands(req)
	if err := validateUpdateProjectRequest(req.Id, commands, localizer); err != nil {
		return nil, err
	}
	if err := s.updateProject(ctx, req.Id, editor, localizer, commands...); err != nil {
		return nil, err
	}
	return &environmentproto.UpdateProjectResponse{}, nil
}

func getUpdateProjectCommands(req *environmentproto.UpdateProjectRequest) []command.Command {
	commands := make([]command.Command, 0)
	if req.ChangeDescriptionCommand != nil {
		commands = append(commands, req.ChangeDescriptionCommand)
	}
	if req.RenameCommand != nil {
		commands = append(commands, req.RenameCommand)
	}
	return commands
}

func validateUpdateProjectRequest(id string, commands []command.Command, localizer locale.Localizer) error {
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
		dt, err := statusProjectIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, cmd := range commands {
		if c, ok := cmd.(*environmentproto.RenameProjectCommand); ok {
			name := strings.TrimSpace(c.Name)
			if name == "" {
				dt, err := statusProjectNameRequired.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
			if len(name) > maxProjectNameLength {
				dt, err := statusInvalidProjectName.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) updateProject(
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
		projectStorage := v2es.NewProjectStorage(tx)
		project, err := projectStorage.GetProject(ctx, id)
		if err != nil {
			return err
		}
		handler, err := command.NewProjectCommandHandler(editor, project, s.publisher)
		if err != nil {
			return err
		}
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		return projectStorage.UpdateProject(ctx, project)
	})
	if err != nil {
		if err == v2es.ErrProjectNotFound || err == v2es.ErrProjectUnexpectedAffectedRows {
			dt, err := statusProjectNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		s.logger.Error(
			"Failed to update project",
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

func (s *EnvironmentService) EnableProject(
	ctx context.Context,
	req *environmentproto.EnableProjectRequest,
) (*environmentproto.EnableProjectResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateEnableProjectRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateProject(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.EnableProjectResponse{}, nil
}

func validateEnableProjectRequest(req *environmentproto.EnableProjectRequest, localizer locale.Localizer) error {
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
		dt, err := statusProjectIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) DisableProject(
	ctx context.Context,
	req *environmentproto.DisableProjectRequest,
) (*environmentproto.DisableProjectResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDisableProjectRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateProject(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.DisableProjectResponse{}, nil
}

func validateDisableProjectRequest(req *environmentproto.DisableProjectRequest, localizer locale.Localizer) error {
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
		dt, err := statusProjectIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) ConvertTrialProject(
	ctx context.Context,
	req *environmentproto.ConvertTrialProjectRequest,
) (*environmentproto.ConvertTrialProjectResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateConvertTrialProjectRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateProject(ctx, req.Id, editor, localizer, req.Command); err != nil {
		return nil, err
	}
	return &environmentproto.ConvertTrialProjectResponse{}, nil
}

func validateConvertTrialProjectRequest(
	req *environmentproto.ConvertTrialProjectRequest,
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
		dt, err := statusProjectIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *EnvironmentService) ListProjectsV2(
	ctx context.Context,
	req *environmentproto.ListProjectsV2Request,
) (*environmentproto.ListProjectsV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	// Check if the user has at least Member role in the requested organization
	_, err := s.checkOrganizationRole(
		ctx,
		req.OrganizationId,
		accountproto.AccountV2_Role_Organization_MEMBER,
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to check organization role",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, err
	}
	var whereParts []mysql.WherePart
	if req.OrganizationId != "" {
		whereParts = append(whereParts,
			mysql.NewFilter("project.organization_id", "=", req.OrganizationId))
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("project.disabled", "=", req.Disabled.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(
			whereParts,
			mysql.NewSearchQuery(
				[]string{"project.id", "project.name", "project.url_code", "project.creator_email"},
				req.SearchKeyword,
			),
		)
	}
	orders, err := s.newProjectListV2Orders(req.OrderBy, req.OrderDirection, localizer)
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
	projectStorage := v2es.NewProjectStorage(s.mysqlClient)
	projects, nextCursor, totalCount, err := projectStorage.ListProjects(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list projects",
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
	return &environmentproto.ListProjectsV2Response{
		Projects:   projects,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *EnvironmentService) newProjectListV2Orders(
	orderBy environmentproto.ListProjectsV2Request_OrderBy,
	orderDirection environmentproto.ListProjectsV2Request_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case environmentproto.ListProjectsV2Request_DEFAULT,
		environmentproto.ListProjectsV2Request_NAME:
		column = "project.name"
	case environmentproto.ListProjectsV2Request_URL_CODE:
		column = "project.url_code"
	case environmentproto.ListProjectsV2Request_ID:
		column = "project.id"
	case environmentproto.ListProjectsV2Request_CREATED_AT:
		column = "project.created_at"
	case environmentproto.ListProjectsV2Request_UPDATED_AT:
		column = "project.updated_at"
	case environmentproto.ListProjectsV2Request_ENVIRONMENT_COUNT:
		column = "environment_count"
	case environmentproto.ListProjectsV2Request_FEATURE_COUNT:
		column = "feature_count"
	case environmentproto.ListProjectsV2Request_CREATOR_EMAIL:
		column = "project.creator_email"
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
	if orderDirection == environmentproto.ListProjectsV2Request_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}
