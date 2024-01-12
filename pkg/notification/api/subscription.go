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
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func (s *NotificationService) CreateSubscription(
	ctx context.Context,
	req *notificationproto.CreateSubscriptionRequest,
) (*notificationproto.CreateSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateSubscriptionRequest(req, localizer); err != nil {
		return nil, err
	}
	subscription, err := domain.NewSubscription(req.Command.Name, req.Command.SourceTypes, req.Command.Recipient)
	if err != nil {
		s.logger.Error(
			"Failed to create a new subscription",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	var handler command.Handler = command.NewEmptySubscriptionCommandHandler()
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
		subscriptionStorage := v2ss.NewSubscriptionStorage(tx)
		if err := subscriptionStorage.CreateSubscription(ctx, subscription, req.EnvironmentNamespace); err != nil {
			return err
		}
		handler = command.NewSubscriptionCommandHandler(editor, subscription, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrSubscriptionAlreadyExists {
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
			"Failed to create subscription",
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
	return &notificationproto.CreateSubscriptionResponse{}, nil
}

func (s *NotificationService) validateCreateSubscriptionRequest(
	req *notificationproto.CreateSubscriptionRequest,
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
	if req.Command.Name == "" {
		dt, err := statusNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.Command.SourceTypes) == 0 {
		dt, err := statusSourceTypesRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if err := s.validateRecipient(req.Command.Recipient, localizer); err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) validateRecipient(
	recipient *notificationproto.Recipient,
	localizer locale.Localizer,
) error {
	if recipient == nil {
		dt, err := statusRecipientRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "recipant"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if recipient.Type == notificationproto.Recipient_SlackChannel {
		return s.validateSlackRecipient(recipient.SlackChannelRecipient, localizer)
	}
	dt, err := statusUnknownRecipient.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "recipant"),
	})
	if err != nil {
		return statusInternal.Err()
	}
	return dt.Err()
}

func (s *NotificationService) validateSlackRecipient(
	sr *notificationproto.SlackChannelRecipient,
	localizer locale.Localizer,
) error {
	// TODO: Check ping to the webhook URL?
	if sr == nil {
		dt, err := statusSlackRecipientRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "slack_recipant"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if sr.WebhookUrl == "" {
		dt, err := statusSlackRecipientWebhookURLRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook_url"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *NotificationService) UpdateSubscription(
	ctx context.Context,
	req *notificationproto.UpdateSubscriptionRequest,
) (*notificationproto.UpdateSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateUpdateSubscriptionRequest(req, localizer); err != nil {
		return nil, err
	}
	commands := s.createUpdateSubscriptionCommands(req)
	if err := s.updateSubscription(ctx, commands, req.Id, req.EnvironmentNamespace, editor, localizer); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
					zap.String("id", req.Id),
				)...,
			)
		}
		return nil, err
	}
	return &notificationproto.UpdateSubscriptionResponse{}, nil
}

func (s *NotificationService) EnableSubscription(
	ctx context.Context,
	req *notificationproto.EnableSubscriptionRequest,
) (*notificationproto.EnableSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateEnableSubscriptionRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateSubscription(
		ctx,
		[]command.Command{req.Command},
		req.Id,
		req.EnvironmentNamespace,
		editor,
		localizer,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to enable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
		}
		return nil, err
	}
	return &notificationproto.EnableSubscriptionResponse{}, nil
}

func (s *NotificationService) validateEnableSubscriptionRequest(
	req *notificationproto.EnableSubscriptionRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
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
	return nil
}

func (s *NotificationService) DisableSubscription(
	ctx context.Context,
	req *notificationproto.DisableSubscriptionRequest,
) (*notificationproto.DisableSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateDisableSubscriptionRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateSubscription(
		ctx,
		[]command.Command{req.Command},
		req.Id,
		req.EnvironmentNamespace,
		editor,
		localizer,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to disable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
		}
		return nil, err
	}
	return &notificationproto.DisableSubscriptionResponse{}, nil
}

func (s *NotificationService) validateDisableSubscriptionRequest(
	req *notificationproto.DisableSubscriptionRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
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
	return nil
}

func (s *NotificationService) updateSubscription(
	ctx context.Context,
	commands []command.Command,
	id, environmentNamespace string,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) error {
	var handler command.Handler = command.NewEmptySubscriptionCommandHandler()
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
		subscriptionStorage := v2ss.NewSubscriptionStorage(tx)
		subscription, err := subscriptionStorage.GetSubscription(ctx, id, environmentNamespace)
		if err != nil {
			return err
		}
		handler = command.NewSubscriptionCommandHandler(editor, subscription, environmentNamespace)
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		if err = subscriptionStorage.UpdateSubscription(ctx, subscription, environmentNamespace); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrSubscriptionNotFound || err == v2ss.ErrSubscriptionUnexpectedAffectedRows {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		s.logger.Error(
			"Failed to update subscription",
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
				zap.String("environmentNamespace", environmentNamespace),
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

func (s *NotificationService) validateUpdateSubscriptionRequest(
	req *notificationproto.UpdateSubscriptionRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if s.isNoUpdateSubscriptionCommand(req) {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.AddSourceTypesCommand != nil && len(req.AddSourceTypesCommand.SourceTypes) == 0 {
		dt, err := statusSourceTypesRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.DeleteSourceTypesCommand != nil && len(req.DeleteSourceTypesCommand.SourceTypes) == 0 {
		dt, err := statusSourceTypesRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "SourceTypes"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.RenameSubscriptionCommand != nil && req.RenameSubscriptionCommand.Name == "" {
		dt, err := statusNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *NotificationService) isNoUpdateSubscriptionCommand(req *notificationproto.UpdateSubscriptionRequest) bool {
	return req.AddSourceTypesCommand == nil &&
		req.DeleteSourceTypesCommand == nil &&
		req.RenameSubscriptionCommand == nil
}

func (s *NotificationService) DeleteSubscription(
	ctx context.Context,
	req *notificationproto.DeleteSubscriptionRequest,
) (*notificationproto.DeleteSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteSubscriptionRequest(req, localizer); err != nil {
		return nil, err
	}
	var handler command.Handler = command.NewEmptySubscriptionCommandHandler()
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
		subscriptionStorage := v2ss.NewSubscriptionStorage(tx)
		subscription, err := subscriptionStorage.GetSubscription(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler = command.NewSubscriptionCommandHandler(editor, subscription, req.EnvironmentNamespace)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		if err = subscriptionStorage.DeleteSubscription(ctx, req.Id, req.EnvironmentNamespace); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrSubscriptionNotFound || err == v2ss.ErrSubscriptionUnexpectedAffectedRows {
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
			"Failed to delete subscription",
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
	return &notificationproto.DeleteSubscriptionResponse{}, nil
}

func validateDeleteSubscriptionRequest(
	req *notificationproto.DeleteSubscriptionRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
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
	return nil
}

func (s *NotificationService) createUpdateSubscriptionCommands(
	req *notificationproto.UpdateSubscriptionRequest,
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

func (s *NotificationService) GetSubscription(
	ctx context.Context,
	req *notificationproto.GetSubscriptionRequest,
) (*notificationproto.GetSubscriptionResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetSubscriptionRequest(req, localizer); err != nil {
		return nil, err
	}
	subscriptionStorage := v2ss.NewSubscriptionStorage(s.mysqlClient)
	subscription, err := subscriptionStorage.GetSubscription(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		if err == v2ss.ErrSubscriptionNotFound {
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
			"Failed to get subscription",
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
	return &notificationproto.GetSubscriptionResponse{Subscription: subscription.Subscription}, nil
}

func validateGetSubscriptionRequest(req *notificationproto.GetSubscriptionRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *NotificationService) ListSubscriptions(
	ctx context.Context,
	req *notificationproto.ListSubscriptionsRequest,
) (*notificationproto.ListSubscriptionsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	var whereParts []mysql.WherePart
	whereParts = append(whereParts, mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace))
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
	orders, err := s.newSubscriptionListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	subscriptions, cursor, totalCount, err := s.listSubscriptionsMySQL(
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
	return &notificationproto.ListSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
		TotalCount:    totalCount,
	}, nil
}

func (s *NotificationService) newSubscriptionListOrders(
	orderBy notificationproto.ListSubscriptionsRequest_OrderBy,
	orderDirection notificationproto.ListSubscriptionsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case notificationproto.ListSubscriptionsRequest_DEFAULT,
		notificationproto.ListSubscriptionsRequest_NAME:
		column = "name"
	case notificationproto.ListSubscriptionsRequest_CREATED_AT:
		column = "created_at"
	case notificationproto.ListSubscriptionsRequest_UPDATED_AT:
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
	if orderDirection == notificationproto.ListSubscriptionsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *NotificationService) ListEnabledSubscriptions(
	ctx context.Context,
	req *notificationproto.ListEnabledSubscriptionsRequest,
) (*notificationproto.ListEnabledSubscriptionsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	var whereParts []mysql.WherePart
	whereParts = append(
		whereParts,
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
		mysql.NewFilter("disabled", "=", false),
	)
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
	subscriptions, cursor, _, err := s.listSubscriptionsMySQL(
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
	return &notificationproto.ListEnabledSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
	}, nil
}

func (s *NotificationService) listSubscriptionsMySQL(
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
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, "", 0, statusInternal.Err()
		}
		return nil, "", 0, dt.Err()
	}
	subscriptionStorage := v2ss.NewSubscriptionStorage(s.mysqlClient)
	subscriptions, nextCursor, totalCount, err := subscriptionStorage.ListSubscriptions(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list subscriptions",
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
