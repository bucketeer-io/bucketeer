// Copyright 2023 The Bucketeer Authors.
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
	"encoding/base64"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/command"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountpb "github.com/bucketeer-io/bucketeer/proto/account"
	autoopspb "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventpb "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	webhookAuthKey = "auth"
)

func (s *AutoOpsService) CreateWebhook(
	ctx context.Context,
	req *autoopspb.CreateWebhookRequest,
) (*autoopspb.CreateWebhookResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountpb.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	status, err := validateCreateWebhook(req, localizer)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	if status != nil {
		s.logger.Error(
			"Failed to validate webhook create request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(status.Err()),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, status.Err()
	}
	resp, err := s.createWebhook(ctx, req, editor)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	return resp, nil
}

func (s *AutoOpsService) createWebhook(
	ctx context.Context,
	req *autoopspb.CreateWebhookRequest,
	editor *eventpb.Editor,
) (*autoopspb.CreateWebhookResponse, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	secret, err := s.generateWebhookSecret(ctx, id.String(), req.EnvironmentNamespace)
	if err != nil {
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	webhook := domain.NewWebhook(id.String(), req.Command.Name, req.Command.Description)
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		webhookStorage := v2as.NewWebhookStorage(tx)
		err := webhookStorage.CreateWebhook(ctx, webhook, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewWebhookCommandHandler(
			editor,
			s.publisher,
			webhook,
			req.EnvironmentNamespace,
		)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &autoopspb.CreateWebhookResponse{
		Webhook: webhook.Webhook,
		Url:     s.createWebhookURL(secret),
	}, nil
}

func (s *AutoOpsService) generateWebhookSecret(
	ctx context.Context,
	id, environmentNamespace string,
) (string, error) {
	ws := domain.NewWebhookSecret(id, environmentNamespace)
	encoded, err := ws.Marshal()
	if err != nil {
		s.logger.Error(
			"Failed to marshal webhook secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return "", err
	}
	encrypted, err := s.webhookCryptoUtil.Encrypt(ctx, encoded)
	if err != nil {
		s.logger.Error(
			"Failed to encrypt webhook secret",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encrypted), nil
}

func validateCreateWebhook(
	req *autoopspb.CreateWebhookRequest,
	localizer locale.Localizer,
) (*status.Status, error) {
	if req.Command == nil {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
	}
	if req.Command.Name == "" {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook name"),
		})
	}
	return nil, nil
}

func (s *AutoOpsService) GetWebhook(
	ctx context.Context,
	req *autoopspb.GetWebhookRequest,
) (*autoopspb.GetWebhookResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountpb.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	status, err := validateGetWebhookRequest(req, localizer)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	if status != nil {
		s.logger.Error(
			"Failed to validate webhook get request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(status.Err()),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, status.Err()
	}
	secret, err := s.generateWebhookSecret(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	webhookStorage := v2as.NewWebhookStorage(s.mysqlClient)
	webhook, err := webhookStorage.GetWebhook(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get webhook",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		if err == v2as.ErrWebhookNotFound {
			dt, err := statusWebhookNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "webhook"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	return &autoopspb.GetWebhookResponse{
		Webhook: webhook.Webhook,
		Url:     s.createWebhookURL(secret),
	}, nil
}

func validateGetWebhookRequest(
	req *autoopspb.GetWebhookRequest,
	localizer locale.Localizer,
) (*status.Status, error) {
	if req.Id == "" {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
	}
	return nil, nil
}

func (s *AutoOpsService) ListWebhooks(
	ctx context.Context,
	req *autoopspb.ListWebhooksRequest,
) (*autoopspb.ListWebhooksResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountpb.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"id", "name", "description"}, req.SearchKeyword))
	}
	orders, err := s.newListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	webhookStorage := v2as.NewWebhookStorage(s.mysqlClient)
	webhooks, nextCursor, totalCount, err := webhookStorage.ListWebhooks(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list webhooks",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	return &autoopspb.ListWebhooksResponse{
		Webhooks:   webhooks,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *AutoOpsService) newListOrders(
	orderBy autoopspb.ListWebhooksRequest_OrderBy,
	orderDirection autoopspb.ListWebhooksRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case autoopspb.ListWebhooksRequest_DEFAULT,
		autoopspb.ListWebhooksRequest_NAME:
		column = "webhook.name"
	case autoopspb.ListWebhooksRequest_CREATED_AT:
		column = "webhook.created_at"
	case autoopspb.ListWebhooksRequest_UPDATED_AT:
		column = "webhook.updated_at"
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
	if orderDirection == autoopspb.ListWebhooksRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *AutoOpsService) UpdateWebhook(
	ctx context.Context,
	req *autoopspb.UpdateWebhookRequest,
) (*autoopspb.UpdateWebhookResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountpb.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	status, err := validateUpdateWebhookRequest(req, localizer)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	if status != nil {
		s.logger.Error(
			"Failed to validate webhook update request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(status.Err()),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, status.Err()
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	var commands []command.Command
	if req.ChangeWebhookDescriptionCommand != nil {
		commands = append(commands, req.ChangeWebhookDescriptionCommand)
	}
	if req.ChangeWebhookNameCommand != nil {
		commands = append(commands, req.ChangeWebhookNameCommand)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		webhookStorage := v2as.NewWebhookStorage(tx)
		webhook, err := webhookStorage.GetWebhook(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewWebhookCommandHandler(editor, s.publisher, webhook, req.EnvironmentNamespace)
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		return webhookStorage.UpdateWebhook(ctx, webhook, req.EnvironmentNamespace)
	})
	if err != nil {
		if err == v2as.ErrWebhookNotFound || err == v2as.ErrAutoOpsRuleUnexpectedAffectedRows {
			dt, err := statusWebhookNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "webhook"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to update webhook",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	return &autoopspb.UpdateWebhookResponse{}, nil
}

func validateUpdateWebhookRequest(
	req *autoopspb.UpdateWebhookRequest,
	localizer locale.Localizer,
) (*status.Status, error) {
	if req.Id == "" {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
	}
	if req.ChangeWebhookNameCommand == nil && req.ChangeWebhookDescriptionCommand == nil {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
	}
	if req.ChangeWebhookNameCommand != nil && req.ChangeWebhookNameCommand.Name == "" {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook name"),
		})
	}
	return nil, nil
}

func (s *AutoOpsService) DeleteWebhook(
	ctx context.Context,
	req *autoopspb.DeleteWebhookRequest,
) (*autoopspb.DeleteWebhookResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountpb.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	status, err := validateDeleteWebhookRequest(req, localizer)
	if err != nil {
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	if status != nil {
		s.logger.Error(
			"Failed to validate webhook delete request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(status.Err()),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, status.Err()
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		webhookStorage := v2as.NewWebhookStorage(tx)
		webhook, err := webhookStorage.GetWebhook(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewWebhookCommandHandler(editor, s.publisher, webhook, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return webhookStorage.DeleteWebhook(ctx, req.Id, req.EnvironmentNamespace)
	})
	if err != nil {
		if err == v2as.ErrWebhookNotFound || err == v2as.ErrAutoOpsRuleUnexpectedAffectedRows {
			dt, err := statusWebhookNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "webhook"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to delete webhook",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.EnvironmentNamespace, localizer)
	}
	return &autoopspb.DeleteWebhookResponse{}, nil
}

func validateDeleteWebhookRequest(
	req *autoopspb.DeleteWebhookRequest,
	localizer locale.Localizer,
) (*status.Status, error) {
	if req.Id == "" {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
	}
	if req.Command == nil {
		return statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
	}
	return nil, nil
}

func (s *AutoOpsService) createWebhookURL(secret string) string {
	url := s.webhookBaseURL
	q := url.Query()
	q.Set(webhookAuthKey, secret)
	url.RawQuery = q.Encode()
	return url.String()
}
