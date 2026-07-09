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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscription/domain"
	v2ss "github.com/bucketeer-io/bucketeer/v2/pkg/subscription/storage/v2"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	subscriptionproto "github.com/bucketeer-io/bucketeer/v2/proto/subscription"
)

func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	req *subscriptionproto.CreateSubscriptionRequest,
) (*subscriptionproto.CreateSubscriptionResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if err := s.validateCreateSubscriptionRequest(req); err != nil {
		return nil, err
	}
	subscription, err := domain.NewSubscription(
		req.Name,
		req.SourceTypes,
		req.Recipient,
		req.FeatureFlagTags,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Any("sourceType", req.SourceTypes),
				zap.String("recipientType", req.Recipient.GetType().String()),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		return s.subscriptionStorage.CreateSubscription(contextWithTx, subscription, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2ss.ErrSubscriptionAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	event, err := domainevent.NewEvent(
		editor,
		eventproto.Event_SUBSCRIPTION,
		subscription.Id,
		eventproto.Event_SUBSCRIPTION_CREATED,
		&eventproto.SubscriptionCreatedEvent{
			SourceTypes:     subscription.SourceTypes,
			Recipient:       subscription.Recipient,
			Name:            subscription.Name,
			FeatureFlagTags: subscription.FeatureFlagTags,
		},
		req.EnvironmentId,
		subscription.Subscription,
		nil,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("subscriptionId", subscription.Id),
			)...,
		)
		return nil, err
	}
	err = s.domainEventPublisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("subscriptionId", subscription.Id),
			)...,
		)
		return nil, err
	}
	return &subscriptionproto.CreateSubscriptionResponse{
		Subscription: subscription.Subscription,
	}, nil
}

func (s *SubscriptionService) UpdateSubscription(
	ctx context.Context,
	req *subscriptionproto.UpdateSubscriptionRequest,
) (*subscriptionproto.UpdateSubscriptionResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := s.validateUpdateSubscriptionRequest(req); err != nil {
		return nil, err
	}

	updatedSubscription, err := s.updateSubscriptionMySQL(
		ctx,
		req.Id,
		req.EnvironmentId,
		req.Name,
		req.SourceTypes,
		req.Disabled,
		req.FeatureFlagTags,
		editor,
	)
	if err != nil {
		return nil, err
	}

	return &subscriptionproto.UpdateSubscriptionResponse{
		Subscription: updatedSubscription,
	}, nil
}

func (s *SubscriptionService) updateSubscriptionMySQL(
	ctx context.Context,
	ID, environmentID string,
	name *wrapperspb.StringValue,
	sourceTypes []subscriptionproto.Subscription_SourceType,
	disabled *wrapperspb.BoolValue,
	featureFlagTags []string,
	editor *eventproto.Editor,
) (*subscriptionproto.Subscription, error) {
	var updatedSubscription *subscriptionproto.Subscription
	var event *eventproto.Event
	err := s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		subscription, err := s.subscriptionStorage.GetSubscription(contextWithTx, ID, environmentID)
		if err != nil {
			return err
		}
		updated, err := subscription.UpdateSubscription(name, sourceTypes, disabled, featureFlagTags)
		if err != nil {
			return err
		}
		updatedSubscription = updated.Subscription
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_SUBSCRIPTION,
			subscription.Id,
			eventproto.Event_SUBSCRIPTION_UPDATED,
			&eventproto.SubscriptionUpdatedEvent{
				Id:              ID,
				SourceTypes:     sourceTypes,
				Name:            name,
				Disabled:        disabled,
				FeatureFlagTags: featureFlagTags,
			},
			ID,
			updatedSubscription,
			subscription,
		)
		if err != nil {
			return err
		}
		return s.subscriptionStorage.UpdateSubscription(contextWithTx, updated, environmentID)
	})
	if err != nil {
		if errors.Is(err, v2ss.ErrSubscriptionNotFound) || errors.Is(err, v2ss.ErrSubscriptionUnexpectedAffectedRows) {
			return nil, statusNotFound.Err()
		}
		s.logger.Error(
			"Failed to update subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", ID),
				zap.String("environmentID", environmentID),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	err = s.domainEventPublisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentID),
				zap.String("ID", ID),
			)...,
		)
		return nil, err
	}
	return updatedSubscription, nil
}

func (s *SubscriptionService) DeleteSubscription(
	ctx context.Context,
	req *subscriptionproto.DeleteSubscriptionRequest,
) (*subscriptionproto.DeleteSubscriptionResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteSubscriptionRequest(req); err != nil {
		return nil, err
	}

	var subscription *domain.Subscription
	var event *eventproto.Event
	err = s.dbClient.RunInTransactionV2(ctx, func(contextWithTx context.Context) error {
		subscription, err = s.subscriptionStorage.GetSubscription(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_SUBSCRIPTION,
			subscription.Id,
			eventproto.Event_SUBSCRIPTION_DELETED,
			&eventproto.SubscriptionDeletedEvent{},
			req.EnvironmentId,
			nil,                       // Current state: entity no longer exists
			subscription.Subscription, // Previous state: what was deleted
		)
		if err = s.subscriptionStorage.DeleteSubscription(contextWithTx, req.Id, req.EnvironmentId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2ss.ErrSubscriptionNotFound) || errors.Is(err, v2ss.ErrSubscriptionUnexpectedAffectedRows) {
			return nil, statusNotFound.Err()
		}
		s.logger.Error(
			"Failed to delete subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	err = s.domainEventPublisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("subscriptionId", subscription.Id),
			)...,
		)
		return nil, err
	}
	return &subscriptionproto.DeleteSubscriptionResponse{}, nil
}

func (s *SubscriptionService) GetSubscription(
	ctx context.Context,
	req *subscriptionproto.GetSubscriptionRequest,
) (*subscriptionproto.GetSubscriptionResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetSubscriptionRequest(req); err != nil {
		return nil, err
	}
	subscription, err := s.subscriptionStorage.GetSubscription(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2ss.ErrSubscriptionNotFound) {
			return nil, statusNotFound.Err()
		}
		s.logger.Error(
			"Failed to get subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &subscriptionproto.GetSubscriptionResponse{Subscription: subscription.Subscription}, nil
}

func (s *SubscriptionService) ListSubscriptions(
	ctx context.Context,
	req *subscriptionproto.ListSubscriptionsRequest,
) (*subscriptionproto.ListSubscriptionsResponse, error) {
	var filterEnvironmentIDs []string
	if req.OrganizationId != "" {
		// console v3
		editor, err := s.checkOrganizationRole(
			ctx, accountproto.AccountV2_Role_Organization_MEMBER,
			req.OrganizationId)
		if err != nil {
			return nil, err
		}
		filterEnvironmentIDs = s.getAllowedEnvironments(req.EnvironmentIds, editor)
	} else {
		// console v2
		_, err := s.checkEnvironmentRole(
			ctx, accountproto.AccountV2_Role_Environment_VIEWER,
			req.EnvironmentId)
		if err != nil {
			return nil, err
		}
		filterEnvironmentIDs = append(filterEnvironmentIDs, req.EnvironmentId)
	}

	var disabled *bool
	if req.Disabled != nil {
		disabled = &req.Disabled.Value
	}

	subscriptions, cursor, totalCount, err := s.listSubscriptions(
		ctx,
		v2ss.ListSubscriptionsParams{
			OrganizationID: req.OrganizationId,
			EnvironmentIDs: filterEnvironmentIDs,
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
	return &subscriptionproto.ListSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
		TotalCount:    totalCount,
	}, nil
}

func (s *SubscriptionService) getAllowedEnvironments(
	reqEnvironmentIDs []string,
	editor *eventproto.Editor,
) []string {
	filterEnvironmentIDs := make([]string, 0)
	if editor.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
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

func (s *SubscriptionService) ListEnabledSubscriptions(
	ctx context.Context,
	req *subscriptionproto.ListEnabledSubscriptionsRequest,
) (*subscriptionproto.ListEnabledSubscriptionsResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	var disabled = false
	subscriptions, cursor, _, err := s.listSubscriptions(
		ctx,
		v2ss.ListSubscriptionsParams{
			EnvironmentIDs: []string{req.EnvironmentId},
			SourceTypes:    req.SourceTypes,
			Disabled:       &disabled,
			PageSize:       req.PageSize,
			Cursor:         req.Cursor,
		},
	)
	if err != nil {
		return nil, err
	}
	return &subscriptionproto.ListEnabledSubscriptionsResponse{
		Subscriptions: subscriptions,
		Cursor:        cursor,
	}, nil
}

func (s *SubscriptionService) listSubscriptions(
	ctx context.Context,
	params v2ss.ListSubscriptionsParams,
) ([]*subscriptionproto.Subscription, string, int64, error) {
	subscriptions, nextCursor, totalCount, err := s.subscriptionStorage.ListSubscriptions(ctx, params)
	if err != nil {
		if errors.Is(err, v2ss.ErrInvalidCursor) {
			return nil, "", 0, statusInvalidCursor.Err()
		}
		if errors.Is(err, v2ss.ErrInvalidOrderBy) {
			return nil, "", 0, statusInvalidOrderBy.Err()
		}
		s.logger.Error(
			"Failed to list subscriptions",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, "", 0, api.NewGRPCStatus(err).Err()
	}
	return subscriptions, strconv.Itoa(nextCursor), totalCount, nil
}
