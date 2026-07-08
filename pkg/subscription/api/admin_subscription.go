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
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	v2ss "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

func (s *SubscriptionService) CreateAdminSubscription(
	ctx context.Context,
	req *subscriptionproto.CreateAdminSubscriptionRequest,
) (*subscriptionproto.CreateAdminSubscriptionResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	subscription, err := domain.NewSubscription(req.Command.Name, req.Command.SourceTypes, req.Command.Recipient, nil)
	if err != nil {
		s.logger.Error(
			"Failed to create a new admin subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Any("sourceType", req.Command.SourceTypes),
				zap.String("recipientType", req.Command.Recipient.GetType().String()),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	var handler = command.NewEmptyAdminSubscriptionCommandHandler()
	err = s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		if err := s.adminSubscriptionStorage.CreateAdminSubscription(contextWithTx, subscription); err != nil {
			return err
		}
		handler, err = command.NewAdminSubscriptionCommandHandler(editor, subscription)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionAlreadyExists {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create admin subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events()); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Any("errors", errs),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &subscriptionproto.CreateAdminSubscriptionResponse{}, nil
}

func (s *SubscriptionService) validateCreateAdminSubscriptionRequest(
	req *subscriptionproto.CreateAdminSubscriptionRequest,
) error {
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	if req.Command.Name == "" {
		return statusNameRequired.Err()
	}
	if len(req.Command.SourceTypes) == 0 {
		return statusSourceTypesRequired.Err()
	}
	if err := s.validateRecipient(req.Command.Recipient); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateAdminSubscription(
	ctx context.Context,
	req *subscriptionproto.UpdateAdminSubscriptionRequest,
) (*subscriptionproto.UpdateAdminSubscriptionResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateUpdateAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	commands := s.createUpdateAdminSubscriptionCommands(req)
	if err := s.updateAdminSubscription(ctx, commands, req.Id, editor); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to update feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
				)...,
			)
		}
		return nil, err
	}
	return &subscriptionproto.UpdateAdminSubscriptionResponse{}, nil
}

func (s *SubscriptionService) validateUpdateAdminSubscriptionRequest(
	req *subscriptionproto.UpdateAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	if s.isNoUpdateAdminSubscriptionCommand(req) {
		return statusNoCommand.Err()
	}
	if req.AddSourceTypesCommand != nil && len(req.AddSourceTypesCommand.SourceTypes) == 0 {
		return statusSourceTypesRequired.Err()
	}
	if req.DeleteSourceTypesCommand != nil && len(req.DeleteSourceTypesCommand.SourceTypes) == 0 {
		return statusSourceTypesRequired.Err()
	}
	if req.RenameSubscriptionCommand != nil && req.RenameSubscriptionCommand.Name == "" {
		return statusNameRequired.Err()
	}
	return nil
}

func (s *SubscriptionService) isNoUpdateAdminSubscriptionCommand(
	req *subscriptionproto.UpdateAdminSubscriptionRequest,
) bool {
	return req.AddSourceTypesCommand == nil &&
		req.DeleteSourceTypesCommand == nil &&
		req.RenameSubscriptionCommand == nil
}

func (s *SubscriptionService) EnableAdminSubscription(
	ctx context.Context,
	req *subscriptionproto.EnableAdminSubscriptionRequest,
) (*subscriptionproto.EnableAdminSubscriptionResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateEnableAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateAdminSubscription(ctx, []command.Command{req.Command}, req.Id, editor); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to enable feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
		}
		return nil, err
	}
	return &subscriptionproto.EnableAdminSubscriptionResponse{}, nil
}

func (s *SubscriptionService) validateEnableAdminSubscriptionRequest(
	req *subscriptionproto.EnableAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *SubscriptionService) DisableAdminSubscription(
	ctx context.Context,
	req *subscriptionproto.DisableAdminSubscriptionRequest,
) (*subscriptionproto.DisableAdminSubscriptionResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.validateDisableAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateAdminSubscription(ctx, []command.Command{req.Command}, req.Id, editor); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to disable feature",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
		}
		return nil, err
	}
	return &subscriptionproto.DisableAdminSubscriptionResponse{}, nil
}

func (s *SubscriptionService) validateDisableAdminSubscriptionRequest(
	req *subscriptionproto.DisableAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *SubscriptionService) updateAdminSubscription(
	ctx context.Context,
	commands []command.Command,
	id string,
	editor *eventproto.Editor,
) error {
	var handler = command.NewEmptyAdminSubscriptionCommandHandler()
	err := s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		subscription, err := s.adminSubscriptionStorage.GetAdminSubscription(contextWithTx, id)
		if err != nil {
			return err
		}
		handler, err = command.NewAdminSubscriptionCommandHandler(editor, subscription)
		if err != nil {
			return err
		}
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		if err = s.adminSubscriptionStorage.UpdateAdminSubscription(contextWithTx, subscription); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionNotFound || err == v2ss.ErrAdminSubscriptionUnexpectedAffectedRows {
			return statusNotFound.Err()
		}
		s.logger.Error(
			"Failed to update admin subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", id),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events()); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("id", id),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

func (s *SubscriptionService) DeleteAdminSubscription(
	ctx context.Context,
	req *subscriptionproto.DeleteAdminSubscriptionRequest,
) (*subscriptionproto.DeleteAdminSubscriptionResponse, error) {
	editor, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	var handler = command.NewEmptyAdminSubscriptionCommandHandler()
	err = s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		subscription, err := s.adminSubscriptionStorage.GetAdminSubscription(contextWithTx, req.Id)
		if err != nil {
			return err
		}
		handler, err = command.NewAdminSubscriptionCommandHandler(editor, subscription)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		if err = s.adminSubscriptionStorage.DeleteAdminSubscription(contextWithTx, req.Id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionNotFound || err == v2ss.ErrAdminSubscriptionUnexpectedAffectedRows {
			return nil, statusNotFound.Err()
		}
		s.logger.Error(
			"Failed to delete admin subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if errs := s.publishDomainEvents(ctx, handler.Events()); len(errs) > 0 {
		s.logger.Error(
			"Failed to publish events",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Any("errors", errs),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &subscriptionproto.DeleteAdminSubscriptionResponse{}, nil
}

func validateDeleteAdminSubscriptionRequest(
	req *subscriptionproto.DeleteAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *SubscriptionService) createUpdateAdminSubscriptionCommands(
	req *subscriptionproto.UpdateAdminSubscriptionRequest,
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

func (s *SubscriptionService) GetAdminSubscription(
	ctx context.Context,
	req *subscriptionproto.GetAdminSubscriptionRequest,
) (*subscriptionproto.GetAdminSubscriptionResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateGetAdminSubscriptionRequest(req); err != nil {
		return nil, err
	}
	subscription, err := s.adminSubscriptionStorage.GetAdminSubscription(ctx, req.Id)
	if err != nil {
		if err == v2ss.ErrAdminSubscriptionNotFound {
			return nil, statusNotFound.Err()
		}
		s.logger.Error(
			"Failed to get admin subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &subscriptionproto.GetAdminSubscriptionResponse{Subscription: subscription.Subscription}, nil
}

func validateGetAdminSubscriptionRequest(
	req *subscriptionproto.GetAdminSubscriptionRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	return nil
}

func (s *SubscriptionService) ListAdminSubscriptions(
	ctx context.Context,
	req *subscriptionproto.ListAdminSubscriptionsRequest,
) (*subscriptionproto.ListAdminSubscriptionsResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	var disabled *bool
	if req.Disabled != nil {
		disabled = &req.Disabled.Value
	}

	subscriptions, cursor, totalCount, err := s.listAdminSubscriptions(
		ctx,
		v2ss.ListAdminSubscriptionsParams{
			SourceTypes:    req.SourceTypes,
			Disabled:       disabled,
			SearchKeyword:  req.SearchKeyword,
			OrderBy:        req.OrderBy,
			OrderDirection: req.OrderDirection,
			PageSize:       req.PageSize,
			Cursor:         req.Cursor,
		},
	)

	if err != nil {
		return nil, err
	}
	return &subscriptionproto.ListAdminSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
		TotalCount:    totalCount,
	}, nil
}

func (s *SubscriptionService) ListEnabledAdminSubscriptions(
	ctx context.Context,
	req *subscriptionproto.ListEnabledAdminSubscriptionsRequest,
) (*subscriptionproto.ListEnabledAdminSubscriptionsResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	var disabled = false
	subscriptions, cursor, _, err := s.listAdminSubscriptions(
		ctx,
		v2ss.ListAdminSubscriptionsParams{
			SourceTypes: req.SourceTypes,
			Disabled:    &disabled,
			PageSize:    req.PageSize,
			Cursor:      req.Cursor,
		},
	)

	if err != nil {
		return nil, err
	}
	return &subscriptionproto.ListEnabledAdminSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
	}, nil
}

func (s *SubscriptionService) listAdminSubscriptions(
	ctx context.Context,
	params v2ss.ListAdminSubscriptionsParams,
) ([]*subscriptionproto.Subscription, string, int64, error) {
	subscriptions, nextCursor, totalCount, err := s.adminSubscriptionStorage.ListAdminSubscriptions(ctx, params)
	if err != nil {
		if errors.Is(err, v2ss.ErrInvalidCursor) {
			return nil, "", 0, statusInvalidCursor.Err()
		}
		if errors.Is(err, v2ss.ErrInvalidOrderBy) {
			return nil, "", 0, statusInvalidOrderBy.Err()
		}
		s.logger.Error(
			"Failed to list admin subscriptions",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, "", 0, api.NewGRPCStatus(err).Err()
	}
	return subscriptions, strconv.Itoa(nextCursor), totalCount, nil
}
