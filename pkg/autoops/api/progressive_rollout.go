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
	"time"

	"github.com/jinzhu/copier"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	exprpto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	fiveMinutes = 5 * time.Minute
)

func (s *AutoOpsService) CreateProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
) (*autoopsproto.CreateProgressiveRolloutResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.Command == nil {
		return s.createProgressiveRolloutNoCommand(ctx, req, editor)
	}

	if err := s.validateCreateProgressiveRolloutRequest(ctx, req); err != nil {
		return nil, err
	}
	progressiveRollout, err := domain.NewProgressiveRollout(
		req.Command.FeatureId,
		req.Command.ProgressiveRolloutManualScheduleClause,
		req.Command.ProgressiveRolloutTemplateScheduleClause,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create domain ProgressiveRollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		handler, err := command.NewProgressiveRolloutCommandHandler(
			editor,
			progressiveRollout,
			s.publisher,
			req.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.prStorage.CreateProgressiveRollout(contextWithTx, progressiveRollout, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrProgressiveRolloutAlreadyExists) {
			return nil, statusProgressiveRolloutAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create ProgressiveRollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.CreateProgressiveRolloutResponse{
		ProgressiveRollout: progressiveRollout.ProgressiveRollout,
	}, nil
}

func (s *AutoOpsService) createProgressiveRolloutNoCommand(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
	editor *eventproto.Editor,
) (*autoopsproto.CreateProgressiveRolloutResponse, error) {
	if err := s.validateCreateProgressiveRolloutRequestNoCommand(ctx, req); err != nil {
		return nil, err
	}
	progressiveRollout, err := domain.NewProgressiveRollout(
		req.FeatureId,
		req.ProgressiveRolloutManualScheduleClause,
		req.ProgressiveRolloutTemplateScheduleClause,
	)
	if err != nil {
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		event, err := domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			progressiveRollout.Id,
			eventproto.Event_PROGRESSIVE_ROLLOUT_CREATED,
			&eventproto.ProgressiveRolloutCreatedEvent{
				Id:        progressiveRollout.Id,
				FeatureId: progressiveRollout.FeatureId,
				Clause:    progressiveRollout.Clause,
				CreatedAt: progressiveRollout.CreatedAt,
				UpdatedAt: progressiveRollout.UpdatedAt,
				Type:      progressiveRollout.Type,
			},
			req.EnvironmentId,
			progressiveRollout.ProgressiveRollout,
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctx, event); err != nil {
			return err
		}
		return s.prStorage.CreateProgressiveRollout(contextWithTx, progressiveRollout, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrProgressiveRolloutAlreadyExists) {
			return nil, statusProgressiveRolloutAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create ProgressiveRollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.CreateProgressiveRolloutResponse{
		ProgressiveRollout: progressiveRollout.ProgressiveRollout,
	}, nil
}

func (s *AutoOpsService) GetProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.GetProgressiveRolloutRequest,
) (*autoopsproto.GetProgressiveRolloutResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetProgressiveRolloutRequest(req); err != nil {
		return nil, err
	}
	progressiveRollout, err := s.prStorage.GetProgressiveRollout(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to get progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2as.ErrProgressiveRolloutNotFound) {
			return nil, statusProgressiveRolloutNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.GetProgressiveRolloutResponse{
		ProgressiveRollout: progressiveRollout.ProgressiveRollout,
	}, nil
}

func (s *AutoOpsService) StopProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.StopProgressiveRolloutRequest,
) (*autoopsproto.StopProgressiveRolloutResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if req.Command == nil {
		return s.stopProgressiveRolloutNoCommand(ctx, req, editor)
	}

	if err := s.validateStopProgressiveRolloutRequest(req); err != nil {
		return nil, err
	}
	err = s.updateProgressiveRollout(
		ctx,
		req.Id,
		req.EnvironmentId,
		req.Command,
		editor,
	)
	if err != nil {
		return nil, err
	}
	return &autoopsproto.StopProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) stopProgressiveRolloutNoCommand(
	ctx context.Context,
	req *autoopsproto.StopProgressiveRolloutRequest,
	editor *eventproto.Editor,
) (*autoopsproto.StopProgressiveRolloutResponse, error) {
	if err := s.validateStopProgressiveRolloutRequest(req); err != nil {
		return nil, err
	}
	var event *eventproto.Event
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		progressiveRollout, err := s.prStorage.GetProgressiveRollout(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		prev := &domain.ProgressiveRollout{}
		if err := copier.Copy(prev, progressiveRollout); err != nil {
			return err
		}
		err = progressiveRollout.Stop(req.StoppedBy)
		if err != nil {
			return err
		}

		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			req.Id,
			eventproto.Event_PROGRESSIVE_ROLLOUT_STOPPED,
			&eventproto.ProgressiveRolloutStoppedEvent{
				Id: req.Id,
			},
			req.EnvironmentId,
			progressiveRollout,
			prev,
		)
		if err != nil {
			return err
		}
		return s.prStorage.UpdateProgressiveRollout(contextWithTx, progressiveRollout, req.EnvironmentId)
	})
	if err != nil {
		s.logger.Error(
			"Failed to stop the progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2as.ErrProgressiveRolloutNotFound) ||
			errors.Is(err, v2as.ErrProgressiveRolloutUnexpectedAffectedRows) {
			return nil, statusProgressiveRolloutNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to push stop progressive rollout event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.StopProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) updateProgressiveRollout(
	ctx context.Context,
	progressiveRolloutID, environmentId string,
	cmd command.Command,
	editor *eventproto.Editor,
) error {
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		progressiveRollout, err := s.prStorage.GetProgressiveRollout(contextWithTx, progressiveRolloutID, environmentId)
		if err != nil {
			return err
		}
		handler, err := command.NewProgressiveRolloutCommandHandler(
			editor,
			progressiveRollout,
			s.publisher,
			environmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, cmd); err != nil {
			return err
		}
		return s.prStorage.UpdateProgressiveRollout(contextWithTx, progressiveRollout, environmentId)
	})
	if err != nil {
		s.logger.Error(
			"Failed to stop the progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", progressiveRolloutID),
				zap.String("environmentId", environmentId),
			)...,
		)
		if errors.Is(err, v2as.ErrProgressiveRolloutNotFound) ||
			errors.Is(err, v2as.ErrProgressiveRolloutUnexpectedAffectedRows) {
			return statusProgressiveRolloutNotFound.Err()
		}
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

func (s *AutoOpsService) DeleteProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.DeleteProgressiveRolloutRequest,
) (*autoopsproto.DeleteProgressiveRolloutResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := s.validateDeleteProgressiveRolloutRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		progressiveRollout, err := s.prStorage.GetProgressiveRollout(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			req.Id,
			eventproto.Event_PROGRESSIVE_ROLLOUT_DELETED,
			&eventproto.ProgressiveRolloutDeletedEvent{
				Id: req.Id,
			},
			req.EnvironmentId,
			nil,                // Current state: entity no longer exists
			progressiveRollout, // Previous state: what was deleted
		)
		if err != nil {
			return err
		}
		return s.prStorage.DeleteProgressiveRollout(contextWithTx, req.Id, req.EnvironmentId)
	})
	if err != nil {
		s.logger.Error(
			"Failed to delete ProgressiveRollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2as.ErrProgressiveRolloutNotFound) ||
			errors.Is(err, v2as.ErrProgressiveRolloutUnexpectedAffectedRows) {
			return nil, statusProgressiveRolloutNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.publisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to push delete progressive rollout event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &autoopsproto.DeleteProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) ListProgressiveRollouts(
	ctx context.Context,
	req *autoopsproto.ListProgressiveRolloutsRequest,
) (*autoopsproto.ListProgressiveRolloutsResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	progressiveRollout, totalCount, nextOffset, err := s.listProgressiveRollouts(
		ctx,
		req,
	)
	if err != nil {
		return nil, err
	}
	return &autoopsproto.ListProgressiveRolloutsResponse{
		ProgressiveRollouts: progressiveRollout,
		TotalCount:          totalCount,
		Cursor:              strconv.Itoa(nextOffset),
	}, nil
}

func (s *AutoOpsService) ExecuteProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.ExecuteProgressiveRolloutRequest,
) (*autoopsproto.ExecuteProgressiveRolloutResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.ChangeProgressiveRolloutTriggeredAtCommand == nil {
		return s.executeProgressiveRolloutNoCommand(ctx, req, editor)
	}

	if err := s.validateExecuteProgressiveRolloutRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		progressiveRollout, err := s.prStorage.GetProgressiveRollout(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		feature, err := s.featureStorage.GetFeature(contextWithTx, progressiveRollout.FeatureId, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := s.checkStopStatus(progressiveRollout); err != nil {
			// If skip if it's already stopped
			return nil
		}
		triggered, err := s.checkAlreadyTriggered(
			req.ChangeProgressiveRolloutTriggeredAtCommand.ScheduleId,
			progressiveRollout,
		)
		if err != nil {
			return err
		}
		if triggered {
			s.logger.Warn(
				"Progressive Rollout is already triggered",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("ruleID", req.ChangeProgressiveRolloutTriggeredAtCommand.ScheduleId),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil
		}
		// Enable the flag if it is disabled and it is the first rollout execution
		var enabled *wrapperspb.BoolValue
		if !feature.Enabled && progressiveRollout.IsWaiting() {
			enabled = &wrapperspb.BoolValue{Value: true}
		}
		handler, err := command.NewProgressiveRolloutCommandHandler(
			editor,
			progressiveRollout,
			s.publisher,
			req.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.ChangeProgressiveRolloutTriggeredAtCommand); err != nil {
			return err
		}
		if err := s.prStorage.UpdateProgressiveRollout(contextWithTx, progressiveRollout, req.EnvironmentId); err != nil {
			return err
		}
		defaultStrategy, err := ExecuteProgressiveRolloutOperation(
			progressiveRollout,
			feature,
			req.ChangeProgressiveRolloutTriggeredAtCommand.ScheduleId,
		)
		if err != nil {
			return err
		}
		// Check if feature already has the target strategy to avoid unnecessary updates
		if proto.Equal(feature.DefaultStrategy, defaultStrategy) && (enabled == nil || feature.Enabled == enabled.Value) {
			s.logger.Warn(
				"Feature already has target strategy, skipping update",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("environmentId", req.EnvironmentId),
					zap.String("id", progressiveRollout.Id),
					zap.String("featureId", progressiveRollout.FeatureId),
				)...,
			)
			return nil
		}
		updated, err := feature.Update(
			nil, // name
			nil, // description
			nil, // tags
			enabled,
			nil, // archived
			defaultStrategy,
			nil,   // offVariation
			false, // resetSamplingSeed
			nil,   // prerequisiteChanges
			nil,   // targetChanges
			nil,   // ruleChanges
			nil,   // variationChanges
			nil,   // tagChanges
		)
		if err != nil {
			return err
		}
		if err := s.featureStorage.UpdateFeature(contextWithTx, updated, req.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to update feature flag",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("id", progressiveRollout.Id),
					zap.String("featureId", progressiveRollout.FeatureId),
				)...,
			)
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FEATURE,
			updated.Id,
			eventproto.Event_FEATURE_UPDATED,
			&eventproto.FeatureUpdatedEvent{
				Id: updated.Id,
			},
			req.EnvironmentId,
			updated.Feature,
			feature.Feature,
			domainevent.WithComment("Progressive rollout executed"),
			domainevent.WithNewVersion(updated.Version),
		)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.logger.Error(
			"Failed to execute progressiveRollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2as.ErrProgressiveRolloutNotFound) ||
			errors.Is(err, v2as.ErrProgressiveRolloutUnexpectedAffectedRows) {
			return nil, statusProgressiveRolloutNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	if event != nil {
		if errs := s.publisher.PublishMulti(ctx, []publisher.Message{event}); len(errs) > 0 {
			s.logger.Error(
				"Failed to publish events",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Any("errors", errs),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, statusInternal.Err()
		}
	}
	return &autoopsproto.ExecuteProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) executeProgressiveRolloutNoCommand(
	ctx context.Context,
	req *autoopsproto.ExecuteProgressiveRolloutRequest,
	editor *eventproto.Editor,
) (*autoopsproto.ExecuteProgressiveRolloutResponse, error) {
	err := s.validateExecuteProgressiveRolloutRequestNoCommand(req)
	if err != nil {
		s.logger.Error(
			"Failed to validate execute progressive rollout request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	var events []publisher.Message
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		progressiveRollout, err := s.prStorage.GetProgressiveRollout(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		feature, err := s.featureStorage.GetFeature(contextWithTx, progressiveRollout.FeatureId, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := s.checkStopStatus(progressiveRollout); err != nil {
			// skip if it's already stopped
			return nil
		}
		triggered, err := s.checkAlreadyTriggered(
			req.ScheduleId,
			progressiveRollout,
		)
		if err != nil {
			return err
		}
		if triggered {
			s.logger.Warn(
				"Progressive Rollout is already triggered",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("ruleID", req.ScheduleId),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil
		}
		// Enable the flag if it is disabled and it is the first rollout execution
		var enabled *wrapperspb.BoolValue
		if !feature.Enabled && progressiveRollout.IsWaiting() {
			enabled = &wrapperspb.BoolValue{Value: true}
		}
		prev := &domain.ProgressiveRollout{}
		if err := copier.Copy(prev, progressiveRollout); err != nil {
			return err
		}
		err = progressiveRollout.SetTriggeredAt(req.ScheduleId)
		if err != nil {
			return err
		}
		if err := s.prStorage.UpdateProgressiveRollout(contextWithTx, progressiveRollout, req.EnvironmentId); err != nil {
			return err
		}
		defaultStrategy, err := ExecuteProgressiveRolloutOperation(
			progressiveRollout,
			feature,
			req.ScheduleId,
		)
		if err != nil {
			return err
		}
		// Check if feature already has the target strategy to avoid unnecessary updates
		if proto.Equal(feature.DefaultStrategy, defaultStrategy) && (enabled == nil || feature.Enabled == enabled.Value) {
			s.logger.Warn(
				"Feature already has target strategy, skipping update",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("environmentId", req.EnvironmentId),
					zap.String("id", progressiveRollout.Id),
					zap.String("featureId", progressiveRollout.FeatureId),
				)...,
			)
			return nil
		}
		updated, err := feature.Update(
			nil, // name
			nil, // description
			nil, // tags
			enabled,
			nil, // archived
			defaultStrategy,
			nil,   // offVariation
			false, // resetSamplingSeed
			nil,   // prerequisiteChanges
			nil,   // targetChanges
			nil,   // ruleChanges
			nil,   // variationChanges
			nil,   // tagChanges
		)
		if err != nil {
			return err
		}
		if err := s.featureStorage.UpdateFeature(contextWithTx, updated, req.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to update feature flag",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("id", progressiveRollout.Id),
					zap.String("featureId", progressiveRollout.FeatureId),
				)...,
			)
			return err
		}
		executeAutoOpsEvent, err := domainevent.NewEvent(
			editor,
			eventproto.Event_AUTOOPS_RULE,
			progressiveRollout.Id,
			eventproto.Event_PROGRESSIVE_ROLLOUT_SCHEDULE_TRIGGERED_AT_CHANGED,
			&eventproto.ProgressiveRolloutScheduleTriggeredAtChangedEvent{
				ScheduleId: req.ScheduleId,
			},
			req.EnvironmentId,
			progressiveRollout.ProgressiveRollout,
			prev,
		)
		if err != nil {
			return err
		}
		updateFeatureEvent, err := domainevent.NewEvent(
			editor,
			eventproto.Event_FEATURE,
			updated.Id,
			eventproto.Event_FEATURE_UPDATED,
			&eventproto.FeatureUpdatedEvent{
				Id: updated.Id,
			},
			req.EnvironmentId,
			updated.Feature,
			feature.Feature,
			domainevent.WithComment("Progressive rollout executed"),
			domainevent.WithNewVersion(updated.Version),
		)
		if err != nil {
			return err
		}
		events = []publisher.Message{executeAutoOpsEvent, updateFeatureEvent}
		return nil
	})
	if err != nil {
		s.logger.Error(
			"Failed to execute progressiveRollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if errors.Is(err, v2as.ErrProgressiveRolloutNotFound) ||
			errors.Is(err, v2as.ErrProgressiveRolloutUnexpectedAffectedRows) {
			return nil, statusProgressiveRolloutNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	if len(events) > 0 {
		if errs := s.publisher.PublishMulti(ctx, events); len(errs) > 0 {
			s.logger.Error(
				"Failed to publish events",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Any("errors", errs),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, statusInternal.Err()
		}
	}
	return &autoopsproto.ExecuteProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) checkStopStatus(p *domain.ProgressiveRollout) error {
	if p.IsStopped() {
		return statusProgressiveRolloutAlreadyStopped.Err()
	}
	return nil
}

func (s *AutoOpsService) checkAlreadyTriggered(
	scheduleID string,
	p *domain.ProgressiveRollout,
) (bool, error) {
	triggered, err := p.AlreadyTriggered(scheduleID)
	if err != nil {
		return false, err
	}
	return triggered, nil
}

func (s *AutoOpsService) listProgressiveRollouts(
	ctx context.Context,
	req *autoopsproto.ListProgressiveRolloutsRequest,
) ([]*autoopsproto.ProgressiveRollout, int64, int, error) {
	filters := []*mysql.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, 0, 0, statusProgressiveRolloutInvalidCursor.Err()
	}
	var inFilters []*mysql.InFilter = nil
	if len(req.FeatureIds) > 0 {
		fIDs := s.convToInterfaceSlice(req.FeatureIds)
		inFilters = append(inFilters, &mysql.InFilter{
			Column: "feature_id",
			Values: fIDs,
		})
	}
	orders, err := s.newListProgressiveRolloutsOrdersMySQL(
		req.OrderBy,
		req.OrderDirection,
	)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, 0, 0, err
	}
	if req.Type != nil {
		filters = append(filters, &mysql.FilterV2{Column: "type", Operator: mysql.OperatorEqual, Value: req.Type})
	}
	if req.Status != nil {
		filters = append(filters, &mysql.FilterV2{Column: "status", Operator: mysql.OperatorEqual, Value: req.Status})
	}
	listOptions := &mysql.ListOptions{
		Filters:     filters,
		Orders:      orders,
		InFilters:   inFilters,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Limit:       limit,
		Offset:      offset,
	}

	progressiveRollouts, totalCount, nextOffset, err := s.prStorage.ListProgressiveRollouts(ctx, listOptions)
	if err != nil {
		s.logger.Error(
			"Failed to list progressive rollouts",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, 0, 0, api.NewGRPCStatus(err).Err()
	}
	return progressiveRollouts, totalCount, nextOffset, nil
}

func (s *AutoOpsService) newListProgressiveRolloutsOrdersMySQL(
	orderBy autoopsproto.ListProgressiveRolloutsRequest_OrderBy,
	orderDirection autoopsproto.ListProgressiveRolloutsRequest_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case autoopsproto.ListProgressiveRolloutsRequest_DEFAULT:
		column = "id"
	case autoopsproto.ListProgressiveRolloutsRequest_CREATED_AT:
		column = "created_at"
	case autoopsproto.ListProgressiveRolloutsRequest_UPDATED_AT:
		column = "updated_at"
	default:
		return nil, statusProgressiveRolloutInvalidOrderBy.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == autoopsproto.ListProgressiveRolloutsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *AutoOpsService) convToInterfaceSlice(
	slice []string,
) []interface{} {
	result := make([]interface{}, 0, len(slice))
	for _, element := range slice {
		result = append(result, element)
	}
	return result
}

func (s *AutoOpsService) validateCreateProgressiveRolloutRequest(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
) error {
	if req.Command.FeatureId == "" {
		return statusProgressiveRolloutFeatureIDRequired.Err()
	}
	// This operation is not the atomic. We may have the problem.
	f, err := s.getFeature(ctx, req.EnvironmentId, req.Command.FeatureId)
	if err != nil {
		return api.NewGRPCStatus(err).Err()
	}
	if err := s.validateTargetFeature(ctx, f); err != nil {
		return err
	}
	if req.Command.ProgressiveRolloutManualScheduleClause == nil &&
		req.Command.ProgressiveRolloutTemplateScheduleClause == nil {
		return statusProgressiveRolloutClauseRequired.Err()
	}
	if req.Command.ProgressiveRolloutManualScheduleClause != nil &&
		req.Command.ProgressiveRolloutTemplateScheduleClause != nil {
		return statusIncorrectProgressiveRolloutClause.Err()
	}
	if req.Command.ProgressiveRolloutManualScheduleClause != nil {
		if err := s.validateProgressiveRolloutManualScheduleClause(
			req.Command.ProgressiveRolloutManualScheduleClause,
			f,
		); err != nil {
			return err
		}
	}
	if req.Command.ProgressiveRolloutTemplateScheduleClause != nil {
		if err := s.validateProgressiveRolloutTemplateScheduleClause(
			req.Command.ProgressiveRolloutTemplateScheduleClause,
			f,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateCreateProgressiveRolloutRequestNoCommand(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
) error {
	if req.FeatureId == "" {
		return statusProgressiveRolloutFeatureIDRequired.Err()
	}
	// This operation is not the atomic. We may have the problem.
	f, err := s.getFeature(ctx, req.EnvironmentId, req.FeatureId)
	if err != nil {
		return api.NewGRPCStatus(err).Err()
	}
	if err := s.validateTargetFeature(ctx, f); err != nil {
		return err
	}
	if req.ProgressiveRolloutManualScheduleClause == nil &&
		req.ProgressiveRolloutTemplateScheduleClause == nil {
		return statusProgressiveRolloutClauseRequired.Err()
	}
	if req.ProgressiveRolloutManualScheduleClause != nil &&
		req.ProgressiveRolloutTemplateScheduleClause != nil {
		return statusIncorrectProgressiveRolloutClause.Err()
	}
	if req.ProgressiveRolloutManualScheduleClause != nil {
		if err := s.validateProgressiveRolloutManualScheduleClause(
			req.ProgressiveRolloutManualScheduleClause,
			f,
		); err != nil {
			return err
		}
	}
	if req.ProgressiveRolloutTemplateScheduleClause != nil {
		if err := s.validateProgressiveRolloutTemplateScheduleClause(
			req.ProgressiveRolloutTemplateScheduleClause,
			f,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateGetProgressiveRolloutRequest(
	req *autoopsproto.GetProgressiveRolloutRequest,
) error {
	if err := s.validateProgressiveRolloutID(req.Id); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateStopProgressiveRolloutRequest(
	req *autoopsproto.StopProgressiveRolloutRequest,
) error {
	if err := s.validateProgressiveRolloutID(req.Id); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateDeleteProgressiveRolloutRequest(
	req *autoopsproto.DeleteProgressiveRolloutRequest,
) error {
	if err := s.validateProgressiveRolloutID(req.Id); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateExecuteProgressiveRolloutRequest(
	req *autoopsproto.ExecuteProgressiveRolloutRequest,
) error {
	if err := s.validateProgressiveRolloutID(req.Id); err != nil {
		return err
	}
	if req.ChangeProgressiveRolloutTriggeredAtCommand.ScheduleId == "" {
		return statusProgressiveRolloutScheduleIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) validateExecuteProgressiveRolloutRequestNoCommand(
	req *autoopsproto.ExecuteProgressiveRolloutRequest,
) error {
	if err := s.validateProgressiveRolloutID(req.Id); err != nil {
		return err
	}
	if req.ScheduleId == "" {
		return statusProgressiveRolloutScheduleIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) validateProgressiveRolloutID(
	id string,
) error {
	if id == "" {
		return statusProgressiveRolloutIDRequired.Err()
	}
	return nil
}

func (s *AutoOpsService) getFeature(
	ctx context.Context,
	environmentID string,
	featureID string,
) (*featureproto.Feature, error) {
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: environmentID,
		Id:            featureID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Feature, nil
}

func (s *AutoOpsService) validateTargetFeature(
	ctx context.Context,
	f *featureproto.Feature,
) error {
	if len(f.Variations) != 2 {
		return statusProgressiveRolloutInvalidVariationSize.Err()
	}
	if err := s.checkIfHasExperiment(ctx, f.Id); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) checkIfHasExperiment(
	ctx context.Context,
	featureID string,
) error {
	// Check if the feature has scheduled or running experiment
	resp, err := s.experimentClient.ListExperiments(ctx, &exprpto.ListExperimentsRequest{
		FeatureId: featureID,
		Statuses: []exprpto.Experiment_Status{
			exprpto.Experiment_WAITING,
			exprpto.Experiment_RUNNING,
		},
	})
	if err != nil {
		return api.NewGRPCStatus(err).Err()
	}
	if len(resp.Experiments) > 0 {
		return statusProgressiveRolloutWaitingOrRunningExperimentExists.Err()
	}
	return nil
}

func (s *AutoOpsService) validateProgressiveRolloutManualScheduleClause(
	clause *autoopsproto.ProgressiveRolloutManualScheduleClause,
	f *featureproto.Feature,
) error {
	if err := s.validateProgressiveRolloutClauseVariationID(
		clause.VariationId,
		f,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseSchedules(
		clause.Schedules,
	); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateProgressiveRolloutTemplateScheduleClause(
	clause *autoopsproto.ProgressiveRolloutTemplateScheduleClause,
	f *featureproto.Feature,
) error {
	if err := s.validateProgressiveRolloutClauseVariationID(
		clause.VariationId,
		f,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseSchedules(
		clause.Schedules,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseIncrements(
		clause.Increments,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseInterval(
		clause.Interval,
	); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateProgressiveRolloutClauseVariationID(
	variationID string,
	f *featureproto.Feature,
) error {
	if variationID == "" {
		return statusProgressiveRolloutClauseVariationIDRequired.Err()
	}
	if exist := s.existVariationID(f, variationID); !exist {
		return statusProgressiveRolloutClauseInvalidVariationID.Err()
	}
	return nil
}

func (s *AutoOpsService) existVariationID(
	f *featureproto.Feature,
	targetVID string,
) bool {
	for _, v := range f.Variations {
		if v.Id == targetVID {
			return true
		}
	}
	return false
}

func (s *AutoOpsService) validateProgressiveRolloutClauseSchedules(
	schedules []*autoopsproto.ProgressiveRolloutSchedule,
) error {
	if len(schedules) == 0 {
		return statusProgressiveRolloutClauseSchedulesRequired.Err()
	}
	for _, s := range schedules {
		if s.ExecuteAt == 0 {
			return statusProgressiveRolloutScheduleExecutedAtRequired.Err()
		}
		if s.Weight < 1 {
			return statusProgressiveRolloutScheduleInvalidWeight.Err()
		}
	}
	if err := s.validateProgressiveRolloutClauseScheduleSpans(schedules); err != nil {
		return err
	}
	return nil
}

func (*AutoOpsService) validateProgressiveRolloutClauseIncrements(
	increments int64,
) error {
	if increments < 1 {
		return statusProgressiveRolloutClauseInvalidIncrements.Err()
	}
	return nil
}

func (*AutoOpsService) validateProgressiveRolloutClauseInterval(
	interval autoopsproto.ProgressiveRolloutTemplateScheduleClause_Interval,
) error {
	if interval == autoopsproto.ProgressiveRolloutTemplateScheduleClause_UNKNOWN {
		return statusProgressiveRolloutClauseUnknownInterval.Err()
	}
	return nil
}

// The span of time for each scheduled time must be at least 5 minutes.
func (*AutoOpsService) validateProgressiveRolloutClauseScheduleSpans(
	schedules []*autoopsproto.ProgressiveRolloutSchedule,
) error {
	for i := 0; i < len(schedules); i++ {
		for j := i + 1; j < len(schedules); j++ {
			if schedules[j].ExecuteAt-schedules[i].ExecuteAt < int64(fiveMinutes.Seconds()) {
				return statusProgressiveRolloutInvalidScheduleSpans.Err()
			}
		}
	}
	return nil
}
