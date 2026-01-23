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
//

package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	pb "github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	numOfSecretCharsToShow = 5
	maskURI                = "********"
)

var webhookEditor = &eventproto.Editor{
	Email:   "webhook",
	IsAdmin: false,
}

func (s *FeatureService) CreateFlagTrigger(
	ctx context.Context,
	request *featureproto.CreateFlagTriggerRequest,
) (*featureproto.CreateFlagTriggerResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
	)
	if err != nil {
		return nil, err
	}

	if err := validateCreateFlagTriggerRequest(request); err != nil {
		s.logger.Error(
			"Error validating create flag trigger request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}
	flagTrigger, err := domain.NewFlagTrigger(
		request.EnvironmentId,
		request.FeatureId,
		request.Type,
		request.Action,
		request.Description,
	)
	if err != nil {
		s.logger.Error(
			"Error creating flag trigger",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		if err := flagTrigger.GenerateToken(); err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FEATURE_CREATED,
			&eventproto.FlagTriggerCreatedEvent{
				Id:            flagTrigger.Id,
				FeatureId:     flagTrigger.FeatureId,
				Type:          flagTrigger.Type,
				Action:        flagTrigger.Action,
				Description:   flagTrigger.Description,
				Token:         flagTrigger.Token,
				CreatedAt:     flagTrigger.CreatedAt,
				UpdatedAt:     flagTrigger.UpdatedAt,
				EnvironmentId: flagTrigger.EnvironmentId,
			},
			request.EnvironmentId,
			flagTrigger,
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.flagTriggerStorage.CreateFlagTrigger(contextWithTx, flagTrigger); err != nil {
			s.logger.Error(
				"Failed to create flag trigger",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	triggerURL := s.generateTriggerURL(ctx, flagTrigger.Token, false)
	flagTrigger.Token = ""

	if err = s.domainPublisher.Publish(ctx, event); err != nil {
		return nil, api.NewGRPCStatus(err).Err()
	}

	return &featureproto.CreateFlagTriggerResponse{
		FlagTrigger: flagTrigger.FlagTrigger,
		Url:         triggerURL,
	}, nil
}

func (s *FeatureService) UpdateFlagTrigger(
	ctx context.Context,
	request *featureproto.UpdateFlagTriggerRequest,
) (*featureproto.UpdateFlagTriggerResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
	)
	if err != nil {
		return nil, err
	}
	var event *eventproto.Event
	var resetURL string
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		flagTrigger, err := s.flagTriggerStorage.GetFlagTrigger(
			contextWithTx,
			request.Id,
			request.EnvironmentId,
		)
		if err != nil {
			return err
		}
		updated, err := flagTrigger.UpdateFlagTrigger(
			request.Description,
			request.Reset_,
			request.Disabled,
		)
		if err != nil {
			return err
		}
		if request.Reset_ {
			resetURL = s.generateTriggerURL(ctx, updated.Token, false)
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FLAG_TRIGGER_UPDATED,
			&eventproto.FlagTriggerUpdateEvent{
				Id:          updated.Id,
				FeatureId:   updated.FeatureId,
				Description: request.Description,
				Reset_:      request.Reset_,
				Disabled:    request.Disabled,
			},
			request.EnvironmentId,
			updated,
			flagTrigger,
		)
		if err != nil {
			return err
		}
		if err := s.flagTriggerStorage.UpdateFlagTrigger(
			contextWithTx,
			updated,
		); err != nil {
			s.logger.Error(
				"Failed to update flag trigger",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) ||
			errors.Is(err, v2fs.ErrFlagTriggerNotFound) {
			return nil, statusTriggerNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	if err = s.domainPublisher.Publish(ctx, event); err != nil {
		return nil, api.NewGRPCStatus(err).Err()
	}

	return &featureproto.UpdateFlagTriggerResponse{
		Url: resetURL,
	}, nil
}

func (s *FeatureService) DeleteFlagTrigger(
	ctx context.Context,
	request *featureproto.DeleteFlagTriggerRequest,
) (*featureproto.DeleteFlagTriggerResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
	)
	if err != nil {
		return nil, err
	}
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		flagTrigger, err := s.flagTriggerStorage.GetFlagTrigger(
			contextWithTx,
			request.Id,
			request.EnvironmentId,
		)
		if err != nil {
			s.logger.Error(
				"Failed to get flag trigger",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FLAG_TRIGGER_DELETED,
			&eventproto.FlagTriggerDeletedEvent{
				Id:            flagTrigger.Id,
				FeatureId:     flagTrigger.FeatureId,
				EnvironmentId: flagTrigger.EnvironmentId,
			},
			request.EnvironmentId,
			nil,         // Current state: entity no longer exists
			flagTrigger, // Previous state: what was deleted
		)
		if err != nil {
			return err
		}
		if err := s.flagTriggerStorage.DeleteFlagTrigger(
			contextWithTx,
			flagTrigger.Id,
			flagTrigger.EnvironmentId,
		); err != nil {
			s.logger.Error(
				"Failed to delete flag trigger",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) ||
			errors.Is(err, v2fs.ErrFlagTriggerNotFound) {
			return nil, statusTriggerNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.domainPublisher.Publish(ctx, event)
	if err != nil {
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.DeleteFlagTriggerResponse{}, nil
}

func (s *FeatureService) GetFlagTrigger(
	ctx context.Context,
	request *featureproto.GetFlagTriggerRequest,
) (*featureproto.GetFlagTriggerResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		request.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetFlagTriggerRequest(request); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}
	trigger, err := s.flagTriggerStorage.GetFlagTrigger(
		ctx,
		request.Id,
		request.EnvironmentId,
	)
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerNotFound) {
			return nil, statusTriggerNotFound.Err()
		}
		return nil, err
	}
	triggerURL := s.generateTriggerURL(ctx, trigger.Token, true)
	trigger.Token = ""
	return &featureproto.GetFlagTriggerResponse{
		FlagTrigger: trigger.FlagTrigger,
		Url:         triggerURL,
	}, nil
}

func (s *FeatureService) ListFlagTriggers(
	ctx context.Context,
	request *featureproto.ListFlagTriggersRequest,
) (*featureproto.ListFlagTriggersResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		request.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateListFlagTriggersRequest(request); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "feature_id",
			Operator: mysql.OperatorEqual,
			Value:    request.FeatureId,
		},
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    request.EnvironmentId,
		},
	}
	orders, err := s.newListFlagTriggerOrders(request.OrderBy, request.OrderDirection)
	if err != nil {
		s.logger.Error(
			"Failed to create order",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	limit := int(request.PageSize)
	cursor := request.Cursor
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
		Orders:      orders,
		Filters:     filters,
		NullFilters: nil,
		JSONFilters: nil,
		InFilters:   nil,
		SearchQuery: nil,
	}
	flagTriggers, nextOffset, totalCount, _ := s.flagTriggerStorage.ListFlagTriggers(ctx, options)
	triggerWithUrls := make([]*featureproto.ListFlagTriggersResponse_FlagTriggerWithUrl, 0, len(flagTriggers))
	for _, trigger := range flagTriggers {
		triggerURL := s.generateTriggerURL(ctx, trigger.Token, true)
		trigger.Token = ""
		triggerWithUrls = append(triggerWithUrls, &featureproto.ListFlagTriggersResponse_FlagTriggerWithUrl{
			FlagTrigger: trigger,
			Url:         triggerURL,
		})
	}
	return &featureproto.ListFlagTriggersResponse{
		FlagTriggers: triggerWithUrls,
		Cursor:       strconv.Itoa(nextOffset),
		TotalCount:   totalCount,
	}, nil
}

func (s *FeatureService) newListFlagTriggerOrders(
	orderBy featureproto.ListFlagTriggersRequest_OrderBy,
	orderDirection featureproto.ListFlagTriggersRequest_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case featureproto.ListFlagTriggersRequest_DEFAULT, featureproto.ListFlagTriggersRequest_CREATED_AT:
		column = "created_at"
	case featureproto.ListFlagTriggersRequest_UPDATED_AT:
		column = "updated_at"
	default:
		return nil, statusInvalidOrderBy.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == featureproto.ListFlagTriggersRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{
		mysql.NewOrder(column, direction),
	}, nil
}

func (s *FeatureService) FlagTriggerWebhook(
	ctx context.Context,
	request *featureproto.FlagTriggerWebhookRequest,
) (*featureproto.FlagTriggerWebhookResponse, error) {
	token := request.GetToken()
	resp := &featureproto.FlagTriggerWebhookResponse{}
	if token == "" {
		s.logger.Error(
			"Failed to get secret from query",
			log.FieldsFromIncomingContext(ctx)...,
		)
		return nil, statusSecretRequired.Err()
	}
	trigger, err := s.flagTriggerStorage.GetFlagTriggerByToken(ctx, token)
	if err != nil {
		s.logger.Error(
			"Failed to get flag trigger",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, statusTriggerNotFound.Err()
	}
	if trigger.GetDisabled() {
		s.logger.Error(
			"Flag trigger is disabled",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, statusTriggerAlreadyDisabled.Err()
	}
	feature, err := s.featureStorage.GetFeature(ctx, trigger.FeatureId, trigger.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2fs.ErrFeatureNotFound) {
			return nil, statusFeatureNotFound.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", trigger.FeatureId),
				zap.String("environmentId", trigger.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if trigger.GetAction() == featureproto.FlagTrigger_Action_ON {
		// check if feature is already enabled
		if !feature.GetEnabled() {
			err := s.updateEnableFeature(ctx, trigger.GetFeatureId(), trigger.GetEnvironmentId(), true)
			if err != nil {
				return nil, statusTriggerEnableFailed.Err()
			}
		}
	} else if trigger.GetAction() == featureproto.FlagTrigger_Action_OFF {
		// check if feature is already disabled
		if feature.GetEnabled() {
			err := s.updateEnableFeature(ctx, trigger.GetFeatureId(), trigger.GetEnvironmentId(), false)
			if err != nil {
				return nil, statusTriggerDisableFailed.Err()
			}
		}
	} else {
		s.logger.Error(
			"Invalid action",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, statusTriggerActionInvalid.Err()
	}
	err = s.updateTriggerUsageInfo(ctx, webhookEditor, trigger)
	if err != nil {
		return nil, statusTriggerUsageUpdateFailed.Err()
	}
	return resp, nil
}

func (s *FeatureService) updateTriggerUsageInfo(
	ctx context.Context,
	editor *eventproto.Editor,
	flagTrigger *domain.FlagTrigger,
) error {
	var event *eventproto.Event
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		prev := &domain.FlagTrigger{}
		if err := copier.Copy(prev, flagTrigger); err != nil {
			return err
		}
		err := flagTrigger.UpdateTriggerUsage()
		if err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FLAG_TRIGGER_USAGE_UPDATED,
			&eventproto.FlagTriggerUsageUpdatedEvent{
				Id:              flagTrigger.Id,
				FeatureId:       flagTrigger.FeatureId,
				LastTriggeredAt: flagTrigger.LastTriggeredAt,
				TriggerTimes:    flagTrigger.TriggerCount,
				EnvironmentId:   flagTrigger.EnvironmentId,
			},
			flagTrigger.EnvironmentId,
			flagTrigger,
			prev,
		)
		if err != nil {
			return err
		}
		err = s.flagTriggerStorage.UpdateFlagTrigger(contextWithTx, flagTrigger)
		if err != nil {
			s.logger.Error(
				"Failed to update flag trigger usage",
				log.FieldsFromIncomingContext(ctx).
					AddFields(zap.Error(err))...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		s.logger.Error(
			"Failed to update flag trigger usage",
			log.FieldsFromIncomingContext(ctx).
				AddFields(zap.Error(err))...,
		)
		return err
	}
	err = s.domainPublisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).
				AddFields(zap.Error(err))...,
		)
		return err
	}
	return nil
}

func (s *FeatureService) updateEnableFeature(
	ctx context.Context,
	featureId, environmentId string,
	enabled bool,
) error {
	var event *eventproto.Event
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		feature, err := s.featureStorage.GetFeature(contextWithTx, featureId, environmentId)
		if err != nil {
			return err
		}

		if enabled {
			err = feature.Enable()
		} else {
			err = feature.Disable()
		}
		if err != nil {
			return err
		}
		prev := &domain.Feature{}
		if err := copier.Copy(prev, feature); err != nil {
			return err
		}
		eventType := eventproto.Event_FEATURE_DISABLED
		var eventData pb.Message = &eventproto.FeatureDisabledEvent{Id: featureId}
		if enabled {
			eventType = eventproto.Event_FEATURE_ENABLED
			eventData = &eventproto.FeatureEnabledEvent{Id: featureId}
		}
		event, err = domainevent.NewEvent(
			webhookEditor,
			eventproto.Event_FEATURE,
			featureId,
			eventType,
			eventData,
			environmentId,
			feature.Feature,
			prev,
		)
		if err != nil {
			s.logger.Error(
				"Failed to create event",
				log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		return s.featureStorage.UpdateFeature(contextWithTx, feature, environmentId)
	})
	if err != nil {
		s.logger.Error(
			"Failed to update feature enabled state",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return err
	}
	err = s.domainPublisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return err
	}
	return nil
}

func (s *FeatureService) generateTriggerURL(
	ctx context.Context,
	token string,
	masked bool,
) string {
	if masked {
		return fmt.Sprintf("%s/%s", s.triggerURL, token[:numOfSecretCharsToShow]+maskURI)
	} else {
		return fmt.Sprintf("%s/%s", s.triggerURL, token)
	}
}
