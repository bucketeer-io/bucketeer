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
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/account/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
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
		return s.accountStorage.CreateAPIKey(contextWithTx, key, req.EnvironmentId)
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.String("role", req.Role.String()),
				zap.String("maintainer", req.Maintainer),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		return s.accountStorage.CreateAPIKey(contextWithTx, key, req.EnvironmentId)
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.String("role", req.Role.String()),
				zap.String("maintainer", req.Maintainer),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
		nil,
	)
	if err != nil {
		return nil, err
	}
	if err := s.publisher.Publish(ctx, e); err != nil {
		s.logger.Error(
			"Failed to publish create account event",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
			log.FieldsFromIncomingContext(ctx).AddFields(
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
				zap.String("name", req.Command.Name),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
			log.FieldsFromIncomingContext(ctx).AddFields(
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
			log.FieldsFromIncomingContext(ctx).AddFields(
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.DisableAPIKeyResponse{}, nil
}

func (s *AccountService) updateAPIKeyMySQL(
	ctx context.Context,
	editor *eventproto.Editor,
	id, environmentID string,
	cmd command.Command,
) error {
	return s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		apiKey, err := s.accountStorage.GetAPIKey(contextWithTx, id, environmentID)
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
		return s.accountStorage.UpdateAPIKey(contextWithTx, apiKey, environmentID)
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if apiKey == nil {
		s.logger.Error(
			"Failed to get api key",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	editor, err := s.checkOrganizationRole(
		ctx, proto.AccountV2_Role_Organization_MEMBER,
		req.OrganizationId, localizer)
	if err != nil {
		return nil, err
	}
	filterEnvironmentIDs := s.getAllowedEnvironments(req.EnvironmentIds, editor)
	if req.OrganizationId == "" {
		dt, err := statusInvalidListAPIKeyRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "organization_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_v2.organization_id",
			Operator: mysql.OperatorEqual,
			Value:    req.OrganizationId,
		},
	}
	var inFilters []*mysql.InFilter
	if len(filterEnvironmentIDs) > 0 {
		environmentIds := make([]interface{}, 0, len(filterEnvironmentIDs))
		for _, id := range filterEnvironmentIDs {
			environmentIds = append(environmentIds, id)
		}
		inFilters = append(inFilters, &mysql.InFilter{
			Column: "api_key.environment_id",
			Values: environmentIds,
		})
	}
	if req.Disabled != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "api_key.disabled",
			Operator: mysql.OperatorEqual,
			Value:    req.Disabled.Value,
		})
	}
	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"api_key.name"},
			Keyword: req.SearchKeyword,
		}
	}
	orders, err := s.newAPIKeyListOrders(req.OrderBy, req.OrderDirection, localizer)
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
	listOptions := &mysql.ListOptions{
		Filters:     filters,
		InFilters:   inFilters,
		SearchQuery: searchQuery,
		Limit:       limit,
		Offset:      offset,
		Orders:      orders,
		NullFilters: nil,
		JSONFilters: nil,
	}
	apiKeys, nextCursor, totalCount, err := s.accountStorage.ListAPIKeys(ctx, listOptions)
	if err != nil {
		s.logger.Error(
			"Failed to list api keys",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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

func (s *AccountService) getAllowedEnvironments(
	reqEnvironmentIDs []string,
	editor *eventproto.Editor,
) []string {
	filterEnvironmentIDs := make([]string, 0)
	if editor.OrganizationRole == proto.AccountV2_Role_Organization_MEMBER {
		// only show API keys in allowed environments for member.
		if len(reqEnvironmentIDs) > 0 {
			for _, id := range reqEnvironmentIDs {
				for _, e := range editor.EnvironmentRoles {
					if e.EnvironmentId == id {
						filterEnvironmentIDs = append(filterEnvironmentIDs, id)
						break
					}
				}
			}
		} else {
			for _, e := range editor.EnvironmentRoles {
				filterEnvironmentIDs = append(filterEnvironmentIDs, e.EnvironmentId)
			}
		}
	} else {
		// if the user is an admin or owner, no need to filter environments.
		filterEnvironmentIDs = append(filterEnvironmentIDs, reqEnvironmentIDs...)
	}
	return filterEnvironmentIDs
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
	case proto.ListAPIKeysRequest_STATE:
		column = "api_key.disabled"
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

func (s *AccountService) GetEnvironmentAPIKey(
	ctx context.Context,
	req *proto.GetEnvironmentAPIKeyRequest,
) (*proto.GetEnvironmentAPIKeyResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkSystemAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
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
	envAPIKey, err := s.accountStorage.GetEnvironmentAPIKey(ctx, req.ApiKey)
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
			"Failed to get environment api key",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKey", req.ApiKey),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	// for security, obfuscate the returned key
	shadowLen := int(float64(len(envAPIKey.ApiKey.ApiKey)) * apiKeyShadowPercentage)
	envAPIKey.ApiKey.ApiKey = envAPIKey.ApiKey.ApiKey[shadowLen:]

	return &proto.GetEnvironmentAPIKeyResponse{
		EnvironmentApiKey: envAPIKey.EnvironmentAPIKey,
	}, nil
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
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		apiKey, err := s.accountStorage.GetAPIKey(contextWithTx, req.Id, req.EnvironmentId)
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

		return s.accountStorage.UpdateAPIKey(contextWithTx, updated, req.EnvironmentId)
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}

	return &proto.UpdateAPIKeyResponse{}, nil
}

func (s *AccountService) UpdateAPIKeyLastUsedAt(
	ctx context.Context,
	req *proto.UpdateAPIKeyLastUsedAtRequest,
) (*proto.UpdateAPIKeyLastUsedAtResponse, error) {
	// No need to check roles since this is only allowed to be called internally by other services.
	err := validateUpdateAPIKeyLastUsedAtRequest(req)
	if err != nil {
		s.logger.Error("Invalid UpdateAPIKeyLastUsedAt request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("apiKeyId", req.ApiKeyId),
				zap.Int64("lastUsedAt", req.LastUsedAt),
			)...)
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		envAPIKey, err := s.accountStorage.GetEnvironmentAPIKey(contextWithTx, req.ApiKeyId)
		if err != nil {
			return err
		}
		apiKey := &domain.APIKey{
			APIKey: envAPIKey.ApiKey,
		}
		err = apiKey.SetUsedAt(req.LastUsedAt)
		if err != nil {
			return err
		}
		return s.accountStorage.UpdateAPIKey(contextWithTx, apiKey, envAPIKey.Environment.Id)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAPIKeyNotFound) {
			return nil, statusNotFound.Err()
		}
		s.logger.Error(
			"Failed to update api key last used at",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.ApiKeyId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.UpdateAPIKeyLastUsedAtResponse{}, nil
}

func validateUpdateAPIKeyLastUsedAtRequest(
	req *proto.UpdateAPIKeyLastUsedAtRequest,
) error {
	if req.ApiKeyId == "" {
		return statusMissingAPIKeyID.Err()
	}

	if req.LastUsedAt <= 0 {
		return statusInvalidLastUsedAt.Err()
	}
	return nil
}
