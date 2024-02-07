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
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func (s *AccountService) CreateAPIKey(
	ctx context.Context,
	req *proto.CreateAPIKeyRequest,
) (*proto.CreateAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		proto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentNamespace,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateCreateAPIKeyRequest(req, localizer); err != nil {
		return nil, err
	}
	key, err := domain.NewAPIKey(req.Command.Name, req.Command.Role)
	if err != nil {
		s.logger.Error(
			"Failed to create a new api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
		handler := command.NewAPIKeyCommandHandler(editor, key, s.publisher, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.accountStorage.CreateAPIKey(ctx, key, req.EnvironmentNamespace)
	})
	if err != nil {
		if err == v2as.ErrAPIKeyAlreadyExists {
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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

func (s *AccountService) ChangeAPIKeyName(
	ctx context.Context,
	req *proto.ChangeAPIKeyNameRequest,
) (*proto.ChangeAPIKeyNameResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		proto.AccountV2_Role_Organization_ADMIN,
		req.EnvironmentNamespace,
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err := s.updateAPIKeyMySQL(ctx, editor, req.Id, req.EnvironmentNamespace, req.Command); err != nil {
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
		req.EnvironmentNamespace,
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err := s.updateAPIKeyMySQL(ctx, editor, req.Id, req.EnvironmentNamespace, req.Command); err != nil {
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
		req.EnvironmentNamespace,
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err := s.updateAPIKeyMySQL(ctx, editor, req.Id, req.EnvironmentNamespace, req.Command); err != nil {
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
			"Failed to disable api key",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	id, environmentNamespace string,
	cmd command.Command,
) error {
	return s.accountStorage.RunInTransaction(ctx, func() error {
		apiKey, err := s.accountStorage.GetAPIKey(ctx, id, environmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewAPIKeyCommandHandler(editor, apiKey, s.publisher, environmentNamespace)
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return s.accountStorage.UpdateAPIKey(ctx, apiKey, environmentNamespace)
	})
}

func (s *AccountService) GetAPIKey(ctx context.Context, req *proto.GetAPIKeyRequest) (*proto.GetAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, proto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentNamespace, localizer)
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
	apiKey, err := s.accountStorage.GetAPIKey(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		if err == v2as.ErrAPIKeyNotFound {
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &proto.GetAPIKeyResponse{ApiKey: apiKey.APIKey}, nil
}

func (s *AccountService) ListAPIKeys(
	ctx context.Context,
	req *proto.ListAPIKeysRequest,
) (*proto.ListAPIKeysResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, proto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("disabled", "=", req.Disabled.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"name"}, req.SearchKeyword))
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
		column = "name"
	case proto.ListAPIKeysRequest_CREATED_AT:
		column = "created_at"
	case proto.ListAPIKeysRequest_UPDATED_AT:
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
		apiKey, err := s.accountStorage.GetAPIKey(ctx, req.Id, e.Id)
		if err != nil {
			if err == v2as.ErrAPIKeyNotFound {
				continue
			}
			s.logger.Error(
				"Failed to get api key",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", e.Id),
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
		return &proto.GetAPIKeyBySearchingAllEnvironmentsResponse{
			EnvironmentApiKey: &proto.EnvironmentAPIKey{
				ApiKey:      apiKey.APIKey,
				ProjectId:   p.Id,
				Environment: e,
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
