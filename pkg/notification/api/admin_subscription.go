// Copyright 2022 The Bucketeer Authors.
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/notification/command"
	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	v2ss "github.com/bucketeer-io/bucketeer/pkg/notification/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func (s *NotificationService) CreateAdminSubscription(
	ctx context.Context,
	req *notificationproto.CreateAdminSubscriptionRequest,
) (*notificationproto.CreateAdminSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	subscription, err := domain.NewSubscription(req.Command.Name, req.Command.SourceTypes, req.Command.Recipient)
	if err != nil {
		s.logger.Error(
			"Failed to create a new admin subscription",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("sourceType", req.Command.SourceTypes),
				zap.Any("recipient", req.Command.Recipient),
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
	var handler command.Handler = command.NewEmptyAdminSubscriptionCommandHandler()
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
		adminSubscriptionStorage := v2ss.NewAdminSubscriptionStorage(tx)
		if err := adminSubscriptionStorage.CreateAdminSubscription(ctx, subscription); err != nil {
			return err
		}
		handler = command.NewAdminSubscriptionCommandHandler(editor, subscription)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionAlreadyExists {
			return nil, localizedError(statusAlreadyExists, locale.JaJP)
		}
		s.logger.Error(
			"Failed to create admin subscription",
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
	if errs := s.publishDomainEvents(ctx, handler.Events()); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
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
	return &notificationproto.CreateAdminSubscriptionResponse{}, nil
}

func (s *NotificationService) validateCreateAdminSubscriptionRequest(
	req *notificationproto.CreateAdminSubscriptionRequest,
) error {
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if req.Command.Name == "" {
		return localizedError(statusNameRequired, locale.JaJP)
	}
	if len(req.Command.SourceTypes) == 0 {
		return localizedError(statusSourceTypesRequired, locale.JaJP)
	}
	if err := s.validateRecipient(req.Command.Recipient); err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) UpdateAdminSubscription(
	ctx context.Context,
	req *notificationproto.UpdateAdminSubscriptionRequest,
) (*notificationproto.UpdateAdminSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateUpdateAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	commands := s.createUpdateAdminSubscriptionCommands(req)
	if err := s.updateAdminSubscription(ctx, commands, req.Id, editor, localizer); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
				)...,
			)
		}
		return nil, err
	}
	return &notificationproto.UpdateAdminSubscriptionResponse{}, nil
}

func (s *NotificationService) validateUpdateAdminSubscriptionRequest(
	req *notificationproto.UpdateAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	if s.isNoUpdateAdminSubscriptionCommand(req) {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	if req.AddSourceTypesCommand != nil && len(req.AddSourceTypesCommand.SourceTypes) == 0 {
		return localizedError(statusSourceTypesRequired, locale.JaJP)
	}
	if req.DeleteSourceTypesCommand != nil && len(req.DeleteSourceTypesCommand.SourceTypes) == 0 {
		return localizedError(statusSourceTypesRequired, locale.JaJP)
	}
	if req.RenameSubscriptionCommand != nil && req.RenameSubscriptionCommand.Name == "" {
		return localizedError(statusNameRequired, locale.JaJP)
	}
	return nil
}

func (s *NotificationService) isNoUpdateAdminSubscriptionCommand(
	req *notificationproto.UpdateAdminSubscriptionRequest,
) bool {
	return req.AddSourceTypesCommand == nil &&
		req.DeleteSourceTypesCommand == nil &&
		req.RenameSubscriptionCommand == nil
}

func (s *NotificationService) EnableAdminSubscription(
	ctx context.Context,
	req *notificationproto.EnableAdminSubscriptionRequest,
) (*notificationproto.EnableAdminSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateEnableAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateAdminSubscription(ctx, []command.Command{req.Command}, req.Id, editor, localizer); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to enable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
		}
		return nil, err
	}
	return &notificationproto.EnableAdminSubscriptionResponse{}, nil
}

func (s *NotificationService) validateEnableAdminSubscriptionRequest(
	req *notificationproto.EnableAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func (s *NotificationService) DisableAdminSubscription(
	ctx context.Context,
	req *notificationproto.DisableAdminSubscriptionRequest,
) (*notificationproto.DisableAdminSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateDisableAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateAdminSubscription(ctx, []command.Command{req.Command}, req.Id, editor, localizer); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to disable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
		}
		return nil, err
	}
	return &notificationproto.DisableAdminSubscriptionResponse{}, nil
}

func (s *NotificationService) validateDisableAdminSubscriptionRequest(
	req *notificationproto.DisableAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func (s *NotificationService) updateAdminSubscription(
	ctx context.Context,
	commands []command.Command,
	id string,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) error {
	var handler command.Handler = command.NewEmptyAdminSubscriptionCommandHandler()
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
		adminSubscriptionStorage := v2ss.NewAdminSubscriptionStorage(tx)
		subscription, err := adminSubscriptionStorage.GetAdminSubscription(ctx, id)
		if err != nil {
			return err
		}
		handler = command.NewAdminSubscriptionCommandHandler(editor, subscription)
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		if err = adminSubscriptionStorage.UpdateAdminSubscription(ctx, subscription); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionNotFound || err == v2ss.ErrAdminSubscriptionUnexpectedAffectedRows {
			return localizedError(statusNotFound, locale.JaJP)
		}
		s.logger.Error(
			"Failed to update admin subscription",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", id),
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
	if errs := s.publishDomainEvents(ctx, handler.Events()); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("id", id),
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

func (s *NotificationService) DeleteAdminSubscription(
	ctx context.Context,
	req *notificationproto.DeleteAdminSubscriptionRequest,
) (*notificationproto.DeleteAdminSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	editor, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	var handler command.Handler = command.NewEmptyAdminSubscriptionCommandHandler()
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
		adminSubscriptionStorage := v2ss.NewAdminSubscriptionStorage(tx)
		subscription, err := adminSubscriptionStorage.GetAdminSubscription(ctx, req.Id)
		if err != nil {
			return err
		}
		handler = command.NewAdminSubscriptionCommandHandler(editor, subscription)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		if err = adminSubscriptionStorage.DeleteAdminSubscription(ctx, req.Id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionNotFound || err == v2ss.ErrAdminSubscriptionUnexpectedAffectedRows {
			return nil, localizedError(statusNotFound, locale.JaJP)
		}
		s.logger.Error(
			"Failed to delete admin subscription",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	if errs := s.publishDomainEvents(ctx, handler.Events()); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Any("errors", errs),
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
	return &notificationproto.DeleteAdminSubscriptionResponse{}, nil
}

func validateDeleteAdminSubscriptionRequest(req *notificationproto.DeleteAdminSubscriptionRequest) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	if req.Command == nil {
		return localizedError(statusNoCommand, locale.JaJP)
	}
	return nil
}

func (s *NotificationService) createUpdateAdminSubscriptionCommands(
	req *notificationproto.UpdateAdminSubscriptionRequest,
) []command.Command {
	commands := make([]command.Command, 0)
	if req.AddSourceTypesCommand != nil {
		commands = append(commands, req.AddSourceTypesCommand)
	}
	if req.DeleteSourceTypesCommand != nil {
		commands = append(commands, req.DeleteSourceTypesCommand)
	}
	if req.RenameSubscriptionCommand != nil {
		commands = append(commands, req.RenameSubscriptionCommand)
	}
	return commands
}

func (s *NotificationService) GetAdminSubscription(
	ctx context.Context,
	req *notificationproto.GetAdminSubscriptionRequest,
) (*notificationproto.GetAdminSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	adminSubscriptionStorage := v2ss.NewAdminSubscriptionStorage(s.mysqlClient)
	subscription, err := adminSubscriptionStorage.GetAdminSubscription(ctx, req.Id)
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionNotFound {
			return nil, localizedError(statusNotFound, locale.JaJP)
		}
		s.logger.Error(
			"Failed to get admin subscription",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	return &notificationproto.GetAdminSubscriptionResponse{Subscription: subscription.Subscription}, nil
}

func validateGetAdminSubscriptionRequest(req *notificationproto.GetAdminSubscriptionRequest) error {
	if req.Id == "" {
		return localizedError(statusIDRequired, locale.JaJP)
	}
	return nil
}

func (s *NotificationService) ListAdminSubscriptions(
	ctx context.Context,
	req *notificationproto.ListAdminSubscriptionsRequest,
) (*notificationproto.ListAdminSubscriptionsResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	var whereParts []mysql.WherePart
	sourceTypesValues := make([]interface{}, len(req.SourceTypes))
	for i, st := range req.SourceTypes {
		sourceTypesValues[i] = int32(st)
	}
	if len(sourceTypesValues) > 0 {
		whereParts = append(
			whereParts,
			mysql.NewJSONFilter("source_types", mysql.JSONContainsNumber, sourceTypesValues),
		)
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("disabled", "=", req.Disabled.Value))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"name"}, req.SearchKeyword))
	}
	orders, err := s.newAdminSubscriptionListOrders(req.OrderBy, req.OrderDirection)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	subscriptions, cursor, totalCount, err := s.listAdminSubscriptionsMySQL(
		ctx,
		whereParts,
		orders,
		req.PageSize,
		req.Cursor,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	return &notificationproto.ListAdminSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
		TotalCount:    totalCount,
	}, nil
}

func (s *NotificationService) newAdminSubscriptionListOrders(
	orderBy notificationproto.ListAdminSubscriptionsRequest_OrderBy,
	orderDirection notificationproto.ListAdminSubscriptionsRequest_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case notificationproto.ListAdminSubscriptionsRequest_DEFAULT,
		notificationproto.ListAdminSubscriptionsRequest_NAME:
		column = "name"
	case notificationproto.ListAdminSubscriptionsRequest_CREATED_AT:
		column = "created_at"
	case notificationproto.ListAdminSubscriptionsRequest_UPDATED_AT:
		column = "updated_at"
	default:
		return nil, localizedError(statusInvalidOrderBy, locale.JaJP)
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == notificationproto.ListAdminSubscriptionsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *NotificationService) ListEnabledAdminSubscriptions(
	ctx context.Context,
	req *notificationproto.ListEnabledAdminSubscriptionsRequest,
) (*notificationproto.ListEnabledAdminSubscriptionsResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkAdminRole(ctx, localizer)
	if err != nil {
		return nil, err
	}
	var whereParts []mysql.WherePart
	whereParts = append(whereParts, mysql.NewFilter("disabled", "=", false))
	sourceTypesValues := make([]interface{}, len(req.SourceTypes))
	for i, st := range req.SourceTypes {
		sourceTypesValues[i] = int32(st)
	}
	if len(sourceTypesValues) > 0 {
		whereParts = append(
			whereParts,
			mysql.NewJSONFilter("source_types", mysql.JSONContainsNumber, sourceTypesValues),
		)
	}
	subscriptions, cursor, _, err := s.listAdminSubscriptionsMySQL(
		ctx,
		whereParts,
		nil,
		req.PageSize,
		req.Cursor,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	return &notificationproto.ListEnabledAdminSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
	}, nil
}

func (s *NotificationService) listAdminSubscriptionsMySQL(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	pageSize int64,
	cursor string,
	localizer locale.Localizer,
) ([]*notificationproto.Subscription, string, int64, error) {
	limit := int(pageSize)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, "", 0, localizedError(statusInvalidCursor, locale.JaJP)
	}
	adminSubscriptionStorage := v2ss.NewAdminSubscriptionStorage(s.mysqlClient)
	subscriptions, nextCursor, totalCount, err := adminSubscriptionStorage.ListAdminSubscriptions(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list admin subscriptions",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, "", 0, statusInternal.Err()
		}
		return nil, "", 0, dt.Err()
	}
	return subscriptions, strconv.Itoa(nextCursor), totalCount, nil
}
