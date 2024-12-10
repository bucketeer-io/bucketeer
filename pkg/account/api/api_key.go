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
	"strconv"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	apiKeyShadowPercentage = 0.75 // hide a part of the api key
)

func (s *AccountService) CreateAPIKey(
	ctx context.Context,
	req *proto.CreateAPIKeyRequest,
) (*proto.CreateAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		proto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}

	if req.Command == nil {
		return s.createAPIKeyNoCommand(ctx, req, localizer, editor)
	}

	if err := validateCreateAPIKeyRequest(req, localizer); err != nil {
		return nil, err
	}
	if req.Maintainer == "" {
		req.Maintainer = editor.Email
	}

	key, err := domain.NewAPIKey(req.Command.Name, req.Command.Role, req.Maintainer, req.Description)
	if err != nil {
		s.logger.Error(
			"Failed to create a new api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		handler, err := command.NewAPIKeyCommandHandler(
			editor,
			key,
			s.publisher,
			req.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.accountStorage.CreateAPIKey(ctx, key, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAPIKeyAlreadyExists) {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	return &proto.CreateAPIKeyResponse{
		ApiKey: key.APIKey,
	}, nil
}

func (s *AccountService) createAPIKeyNoCommand(
	ctx context.Context,
	req *proto.CreateAPIKeyRequest,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) (*proto.CreateAPIKeyResponse, error) {
	if err := validateCreateAPIKeyRequestNoCommand(req, localizer); err != nil {
		return nil, err
	}
	if req.Maintainer == "" {
		req.Maintainer = editor.Email
	}

	key, err := domain.NewAPIKey(req.Name, req.Role, req.Maintainer, req.Description)
	if err != nil {
		s.logger.Error(
			"Failed to create a new api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.String("role", req.Role.String()),
				zap.String("maintainer", req.Maintainer),
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

	err = s.accountStorage.RunInTransaction(ctx, func() error {
		return s.accountStorage.CreateAPIKey(ctx, key, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAPIKeyAlreadyExists) {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.String("role", req.Role.String()),
				zap.String("maintainer", req.Maintainer),
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

	prev := &domain.APIKey{}
	if err := copier.Copy(prev, key); err != nil {
		return nil, err
	}
	e, err := domainevent.NewEvent(
		editor,
		eventproto.Event_APIKEY,
		key.Id,
		eventproto.Event_APIKEY_CREATED,
		&eventproto.APIKeyCreatedEvent{
			Id:         key.Id,
			Name:       key.Name,
			Role:       key.Role,
			Disabled:   key.Disabled,
			CreatedAt:  key.CreatedAt,
			UpdatedAt:  key.UpdatedAt,
			Maintainer: key.Maintainer,
			ApiKey:     key.ApiKey,
		},
		req.EnvironmentId,
		key.APIKey,
		prev,
	)
	if err != nil {
		return nil, err
	}
	if err := s.publisher.Publish(ctx, e); err != nil {
		s.logger.Error(
			"Failed to publish create account event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.String("role", req.Role.String()),
				zap.String("maintainer", req.Maintainer),
			)...,
		)
		return nil, err
	}

	return &proto.CreateAPIKeyResponse{
		ApiKey: key.APIKey,
	}, nil
}

func (s *AccountService) ChangeAPIKeyName(
	ctx context.Context,
	req *proto.ChangeAPIKeyNameRequest,
) (*proto.ChangeAPIKeyNameResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		proto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateChangeAPIKeyNameRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to change api key name",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := s.updateAPIKeyMySQL(
		ctx,
		editor,
		req.Id,
		req.EnvironmentId,
		req.Command,
	); err != nil {
		if err == v2as.ErrAPIKeyNotFound || err == v2as.ErrAPIKeyUnexpectedAffectedRows {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to change api key name",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
				zap.String("name", req.Command.Name),
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
	return &proto.ChangeAPIKeyNameResponse{}, nil
}

func (s *AccountService) EnableAPIKey(
	ctx context.Context,
	req *proto.EnableAPIKeyRequest,
) (*proto.EnableAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		proto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateEnableAPIKeyRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to enable api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := s.updateAPIKeyMySQL(
		ctx,
		editor,
		req.Id,
		req.EnvironmentId,
		req.Command,
	); err != nil {
		if err == v2as.ErrAPIKeyNotFound || err == v2as.ErrAPIKeyUnexpectedAffectedRows {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to enable api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
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
	return &proto.EnableAPIKeyResponse{}, nil
}

func (s *AccountService) DisableAPIKey(
	ctx context.Context,
	req *proto.DisableAPIKeyRequest,
) (*proto.DisableAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		proto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDisableAPIKeyRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to disable api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := s.updateAPIKeyMySQL(
		ctx,
		editor,
		req.Id,
		req.EnvironmentId,
		req.Command,
	); err != nil {
		if errors.Is(err, v2as.ErrAPIKeyNotFound) || errors.Is(err, v2as.ErrAPIKeyUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to disable api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
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
	return &proto.DisableAPIKeyResponse{}, nil
}

func (s *AccountService) updateAPIKeyMySQL(
	ctx context.Context,
	editor *eventproto.Editor,
	id, environmentID string,
	cmd command.Command,
) error {
	return s.accountStorage.RunInTransaction(ctx, func() error {
		apiKey, err := s.accountStorage.GetAPIKey(ctx, id, environmentID)
		if err != nil {
			return err
		}
		handler, err := command.NewAPIKeyCommandHandler(editor, apiKey, s.publisher, environmentID)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return s.accountStorage.UpdateAPIKey(ctx, apiKey, environmentID)
	})
}

func (s *AccountService) GetAPIKey(ctx context.Context, req *proto.GetAPIKeyRequest) (*proto.GetAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, proto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusMissingAPIKeyID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	apiKey, err := s.accountStorage.GetAPIKey(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2as.ErrAPIKeyNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
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
	if apiKey == nil {
		s.logger.Error(
			"Failed to get api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, statusNotFound.Err()
	}

	// for security, obfuscate the returned key
	shadowLen := int(float64(len(apiKey.ApiKey)) * apiKeyShadowPercentage)
	apiKey.ApiKey = apiKey.ApiKey[shadowLen:]

	return &proto.GetAPIKeyResponse{ApiKey: apiKey.APIKey}, nil
}

func (s *AccountService) ListAPIKeys(
	ctx context.Context,
	req *proto.ListAPIKeysRequest,
) (*proto.ListAPIKeysResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, proto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if len(req.EnvironmentIds) == 0 && req.OrganizationId == "" {
		dt, err := statusInvalidListAPIKeyRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "environment_ids or organization_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	whereParts := []mysql.WherePart{}
	if len(req.EnvironmentIds) > 0 {
		environmentIds := make([]interface{}, 0, len(req.EnvironmentIds))
		for _, id := range req.EnvironmentIds {
			environmentIds = append(environmentIds, id)
		}
		whereParts = append(whereParts, mysql.NewInFilter("api_key.environment_id", environmentIds))
	}
	if req.OrganizationId != "" {
		whereParts = append(whereParts, mysql.NewFilter("environment_v2.organization_id", "=", req.OrganizationId))
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("api_key.disabled", "=", req.Disabled.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"api_key.name"}, req.SearchKeyword))
	}
	orders, err := s.newAPIKeyListOrders(req.OrderBy, req.OrderDirection, localizer)
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
	apiKeys, nextCursor, totalCount, err := s.accountStorage.ListAPIKeys(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list api keys",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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

	// for security, obfuscate the returned key
	for i := 0; i < len(apiKeys); i++ {
		shadowLen := int(float64(len(apiKeys[i].ApiKey)) * apiKeyShadowPercentage)
		apiKeys[i].ApiKey = apiKeys[i].ApiKey[shadowLen:]
	}

	return &proto.ListAPIKeysResponse{
		ApiKeys:    apiKeys,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *AccountService) newAPIKeyListOrders(
	orderBy proto.ListAPIKeysRequest_OrderBy,
	orderDirection proto.ListAPIKeysRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case proto.ListAPIKeysRequest_DEFAULT,
		proto.ListAPIKeysRequest_NAME:
		column = "api_key.name"
	case proto.ListAPIKeysRequest_CREATED_AT:
		column = "api_key.created_at"
	case proto.ListAPIKeysRequest_UPDATED_AT:
		column = "api_key.updated_at"
	case proto.ListAPIKeysRequest_ROLE:
		column = "api_key.role"
	case proto.ListAPIKeysRequest_ENVIRONMENT:
		column = "environment_v2.name"
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
	if orderDirection == proto.ListAPIKeysRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *AccountService) GetAPIKeyBySearchingAllEnvironments(
	ctx context.Context,
	req *proto.GetAPIKeyBySearchingAllEnvironmentsRequest,
) (*proto.GetAPIKeyBySearchingAllEnvironmentsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}

	// TODO: support both fields, when migration finished, remove this block
	if req.ApiKey == "" {
		req.ApiKey = req.Id
	}

	if req.ApiKey == "" {
		dt, err := statusMissingAPIKeyID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "api_key_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	projects, err := s.listProjects(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to get project list",
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
	if len(projects) == 0 {
		s.logger.Error(
			"Could not find any projects",
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
	environments, err := s.listEnvironments(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to get environment list",
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
	if len(environments) == 0 {
		s.logger.Error(
			"Could not find any environments",
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
	projectSet := s.makeProjectSet(projects)
	for _, e := range environments {
		p, ok := projectSet[e.ProjectId]
		if !ok || p.Disabled {
			continue
		}
		apiKey, err := s.accountStorage.GetAPIKeyByAPIKey(ctx, req.ApiKey, e.Id)
		if err != nil {
			if errors.Is(err, v2as.ErrAPIKeyNotFound) {
				continue
			}
			s.logger.Error(
				"Failed to get api key",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", e.Id),
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

		// for security, obfuscate the returned key
		shadowLen := int(float64(len(apiKey.ApiKey)) * apiKeyShadowPercentage)
		apiKey.ApiKey = apiKey.ApiKey[shadowLen:]

		return &proto.GetAPIKeyBySearchingAllEnvironmentsResponse{
			EnvironmentApiKey: &proto.EnvironmentAPIKey{
				ApiKey:         apiKey.APIKey,
				ProjectId:      p.Id,
				ProjectUrlCode: p.UrlCode,
				Environment:    e,
			},
		}, nil
	}
	dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalize(locale.NotFoundError),
	})
	if err != nil {
		return nil, statusInternal.Err()
	}
	return nil, dt.Err()
}

func (s *AccountService) UpdateAPIKey(
	ctx context.Context,
	req *proto.UpdateAPIKeyRequest,
) (*proto.UpdateAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		proto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}

	if err := validateUpdateAPIKeyRequestNoCommand(req, localizer); err != nil {
		return nil, err
	}

	var prev, current *proto.APIKey
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		apiKey, err := s.accountStorage.GetAPIKey(ctx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		prev = apiKey.APIKey

		// Update fields
		updated, err := apiKey.Update(
			req.Name,
			req.Description,
			req.Role,
			req.Maintainer,
			req.Disabled,
		)
		if err != nil {
			return err
		}
		current = updated.APIKey

		return s.accountStorage.UpdateAPIKey(ctx, updated, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAPIKeyNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to update api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
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
	e, err := domainevent.NewEvent(
		editor,
		eventproto.Event_APIKEY,
		req.Id,
		eventproto.Event_APIKEY_CHANGED,
		&eventproto.APIKeyChangedEvent{
			Id: req.Id,
		},
		req.EnvironmentId,
		current,
		prev,
	)
	if err != nil {
		return nil, err
	}
	if err := s.publisher.Publish(ctx, e); err != nil {
		s.logger.Error(
			"Failed to publish update api key event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}

	return &proto.UpdateAPIKeyResponse{}, nil
}
