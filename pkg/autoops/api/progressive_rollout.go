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
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/command"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	exprpto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	fiveMinutes = 5 * time.Minute
)

func (s *AutoOpsService) CreateProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
) (*autoopsproto.CreateProgressiveRolloutResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateProgressiveRolloutRequest(ctx, req, localizer); err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		progressiveRollout, err := domain.NewProgressiveRollout(
			req.Command.FeatureId,
			req.Command.ProgressiveRolloutManualScheduleClause,
			req.Command.ProgressiveRolloutTemplateScheduleClause,
		)
		if err != nil {
			return err
		}
		storage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
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
		return storage.CreateProgressiveRollout(contextWithTx, progressiveRollout, req.EnvironmentId)
	})
	if err != nil {
		switch err {
		case v2as.ErrProgressiveRolloutAlreadyExists:
			dt, err := statusProgressiveRolloutAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusProgressiveRolloutInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create ProgressiveRollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusProgressiveRolloutInternal.Err()
		}
		return nil, dt.Err()
	}
	return &autoopsproto.CreateProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) GetProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.GetProgressiveRolloutRequest,
) (*autoopsproto.GetProgressiveRolloutResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetProgressiveRolloutRequest(req, localizer); err != nil {
		return nil, err
	}
	storage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
	progressiveRollout, err := storage.GetProgressiveRollout(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to get progressive rollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if err == v2as.ErrProgressiveRolloutNotFound {
			dt, err := statusProgressiveRolloutNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout),
			})
			if err != nil {
				return nil, statusProgressiveRolloutInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusProgressiveRolloutInternal.Err()
		}
		return nil, dt.Err()
	}
	return &autoopsproto.GetProgressiveRolloutResponse{
		ProgressiveRollout: progressiveRollout.ProgressiveRollout,
	}, nil
}

func (s *AutoOpsService) StopProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.StopProgressiveRolloutRequest,
) (*autoopsproto.StopProgressiveRolloutResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateStopProgressiveRolloutRequest(req, localizer); err != nil {
		return nil, err
	}
	err = s.updateProgressiveRollout(
		ctx,
		req.Id,
		req.EnvironmentId,
		req.Command,
		editor,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	return &autoopsproto.StopProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) updateProgressiveRollout(
	ctx context.Context,
	progressiveRolloutID, environmentId string,
	cmd command.Command,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) error {
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		storage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
		progressiveRollout, err := storage.GetProgressiveRollout(contextWithTx, progressiveRolloutID, environmentId)
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
		return storage.UpdateProgressiveRollout(contextWithTx, progressiveRollout, environmentId)
	})
	if err != nil {
		s.logger.Error(
			"Failed to stop the progressive rollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", progressiveRolloutID),
				zap.String("environmentId", environmentId),
			)...,
		)
		if err == v2as.ErrProgressiveRolloutNotFound || err == v2as.ErrProgressiveRolloutUnexpectedAffectedRows {
			dt, err := statusProgressiveRolloutNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout),
			})
			if err != nil {
				return statusProgressiveRolloutInternal.Err()
			}
			return dt.Err()
		}
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) DeleteProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.DeleteProgressiveRolloutRequest,
) (*autoopsproto.DeleteProgressiveRolloutResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateDeleteProgressiveRolloutRequest(req, localizer); err != nil {
		return nil, err
	}

	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		storage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
		progressiveRollout, err := storage.GetProgressiveRollout(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
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
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return storage.DeleteProgressiveRollout(contextWithTx, req.Id, req.EnvironmentId)
	})
	if err != nil {
		s.logger.Error(
			"Failed to delete ProgressiveRollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if err == v2as.ErrProgressiveRolloutNotFound || err == v2as.ErrProgressiveRolloutUnexpectedAffectedRows {
			dt, err := statusProgressiveRolloutNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout),
			})
			if err != nil {
				return nil, statusProgressiveRolloutInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusProgressiveRolloutInternal.Err()
		}
		return nil, dt.Err()
	}
	return &autoopsproto.DeleteProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) ListProgressiveRollouts(
	ctx context.Context,
	req *autoopsproto.ListProgressiveRolloutsRequest,
) (*autoopsproto.ListProgressiveRolloutsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	progressiveRollout, totalCount, nextOffset, err := s.listProgressiveRollouts(
		ctx,
		req,
		localizer,
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
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateExecuteProgressiveRolloutRequest(req, localizer); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, tx mysql.Transaction) error {
		storage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
		progressiveRollout, err := storage.GetProgressiveRollout(ctx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		ftStorage := ftstorage.NewFeatureStorage(tx)
		feature, err := ftStorage.GetFeature(ctx, progressiveRollout.FeatureId, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := s.checkStopStatus(progressiveRollout, localizer); err != nil {
			s.logger.Warn(
				"Progressive rollout is already stopped",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("id", progressiveRollout.Id),
					zap.String("featureId", progressiveRollout.FeatureId),
				)...,
			)
			return nil
		}
		triggered, err := s.checkAlreadyTriggered(
			req.ChangeProgressiveRolloutTriggeredAtCommand,
			progressiveRollout,
		)
		if err != nil {
			return err
		}
		if triggered {
			s.logger.Warn(
				"Progressive Rollout is already triggered",
				log.FieldsFromImcomingContext(ctx).AddFields(
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
		if err := storage.UpdateProgressiveRollout(ctx, progressiveRollout, req.EnvironmentId); err != nil {
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
		updated, err := feature.Update(
			nil, nil, nil,
			enabled,
			nil, nil, nil, nil, nil,
			defaultStrategy,
			nil,
		)
		if err != nil {
			return err
		}
		if err := ftStorage.UpdateFeature(ctx, updated, req.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to update feature flag",
				log.FieldsFromImcomingContext(ctx).AddFields(
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		if err == v2as.ErrProgressiveRolloutNotFound || err == v2as.ErrProgressiveRolloutUnexpectedAffectedRows {
			dt, err := statusProgressiveRolloutNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, locale.ProgressiveRollout),
			})
			if err != nil {
				return nil, statusProgressiveRolloutInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusProgressiveRolloutInternal.Err()
		}
		return nil, dt.Err()
	}
	if event != nil {
		if errs := s.publisher.PublishMulti(ctx, []publisher.Message{event}); len(errs) > 0 {
			s.logger.Error(
				"Failed to publish events",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Any("errors", errs),
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
	}
	return &autoopsproto.ExecuteProgressiveRolloutResponse{}, nil
}

func (s *AutoOpsService) checkStopStatus(p *domain.ProgressiveRollout, localizer locale.Localizer) error {
	if p.IsStopped() {
		dt, err := statusProgressiveRolloutAlreadyStopped.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) checkAlreadyTriggered(
	cmd *autoopsproto.ChangeProgressiveRolloutScheduleTriggeredAtCommand,
	p *domain.ProgressiveRollout,
) (bool, error) {
	triggered, err := p.AlreadyTriggered(cmd.ScheduleId)
	if err != nil {
		return false, err
	}
	return triggered, nil
}

func (s *AutoOpsService) listProgressiveRollouts(
	ctx context.Context,
	req *autoopsproto.ListProgressiveRolloutsRequest,
	localizer locale.Localizer,
) ([]*autoopsproto.ProgressiveRollout, int64, int, error) {
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		dt, err := statusProgressiveRolloutInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, 0, 0, statusProgressiveRolloutInternal.Err()
		}
		return nil, 0, 0, dt.Err()
	}
	if len(req.FeatureIds) > 0 {
		fIDs := s.convToInterfaceSlice(req.FeatureIds)
		whereParts = append(whereParts, mysql.NewInFilter("feature_id", fIDs))
	}
	orders, err := s.newListProgressiveRolloutsOrdersMySQL(
		req.OrderBy,
		req.OrderDirection,
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, 0, 0, err
	}
	if req.Type != nil {
		whereParts = append(whereParts, mysql.NewFilter("type", "=", req.Type))
	}
	if req.Status != nil {
		whereParts = append(whereParts, mysql.NewFilter("status", "=", req.Status))
	}
	storage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
	progressiveRollouts, totalCount, nextOffset, err := storage.ListProgressiveRollouts(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list progressive rollouts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, 0, 0, statusProgressiveRolloutInternal.Err()
		}
		return nil, 0, 0, dt.Err()
	}
	return progressiveRollouts, totalCount, nextOffset, nil
}

func (s *AutoOpsService) newListProgressiveRolloutsOrdersMySQL(
	orderBy autoopsproto.ListProgressiveRolloutsRequest_OrderBy,
	orderDirection autoopsproto.ListProgressiveRolloutsRequest_OrderDirection,
	localizer locale.Localizer,
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
		dt, err := statusProgressiveRolloutInvalidOrderBy.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by"),
		})
		if err != nil {
			return nil, statusProgressiveRolloutInternal.Err()
		}
		return nil, dt.Err()
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
	localizer locale.Localizer,
) error {
	if req.Command == nil {
		dt, err := statusProgressiveRolloutNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.FeatureId == "" {
		dt, err := statusProgressiveRolloutFeatureIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	// This operation is not the atomic. We may have the problem.
	f, err := s.getFeature(ctx, req, localizer)
	if err != nil {
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if err := s.validateTargetFeature(ctx, f, localizer); err != nil {
		return err
	}
	if req.Command.ProgressiveRolloutManualScheduleClause == nil &&
		req.Command.ProgressiveRolloutTemplateScheduleClause == nil {
		dt, err := statusProgressiveRolloutClauseRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "clause"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.ProgressiveRolloutManualScheduleClause != nil &&
		req.Command.ProgressiveRolloutTemplateScheduleClause != nil {
		dt, err := statusIncorrectProgressiveRolloutClause.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "clause"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.ProgressiveRolloutManualScheduleClause != nil {
		if err := s.validateProgressiveRolloutManualScheduleClause(
			req.Command.ProgressiveRolloutManualScheduleClause,
			f,
			localizer,
		); err != nil {
			return err
		}
	}
	if req.Command.ProgressiveRolloutTemplateScheduleClause != nil {
		if err := s.validateProgressiveRolloutTemplateScheduleClause(
			req.Command.ProgressiveRolloutTemplateScheduleClause,
			f,
			localizer,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *AutoOpsService) validateGetProgressiveRolloutRequest(
	req *autoopsproto.GetProgressiveRolloutRequest,
	localizer locale.Localizer,
) error {
	if err := s.validateID(req.Id, localizer); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateStopProgressiveRolloutRequest(
	req *autoopsproto.StopProgressiveRolloutRequest,
	localizer locale.Localizer,
) error {
	if err := s.validateID(req.Id, localizer); err != nil {
		return err
	}
	if err := s.validateCommand(req.Command, localizer); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateDeleteProgressiveRolloutRequest(
	req *autoopsproto.DeleteProgressiveRolloutRequest,
	localizer locale.Localizer,
) error {
	if err := s.validateID(req.Id, localizer); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateExecuteProgressiveRolloutRequest(
	req *autoopsproto.ExecuteProgressiveRolloutRequest,
	localizer locale.Localizer,
) error {
	if err := s.validateID(req.Id, localizer); err != nil {
		return err
	}
	if req.ChangeProgressiveRolloutTriggeredAtCommand == nil {
		dt, err := statusProgressiveRolloutNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if req.ChangeProgressiveRolloutTriggeredAtCommand.ScheduleId == "" {
		dt, err := statusProgressiveRolloutScheduleIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "schedule_id"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) validateID(
	id string,
	localizer locale.Localizer,
) error {
	if id == "" {
		dt, err := statusProgressiveRolloutIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) validateCommand(
	command interface{},
	localizer locale.Localizer,
) error {
	switch command := command.(type) {
	case *autoopsproto.StopProgressiveRolloutCommand:
		if command == nil {
			dt, err := statusProgressiveRolloutNoCommand.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
			})
			if err != nil {
				return statusProgressiveRolloutInternal.Err()
			}
			return dt.Err()
		}
	default:
		return nil
	}
	return nil
}

func (s *AutoOpsService) getFeature(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
	localizer locale.Localizer,
) (*featureproto.Feature, error) {
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: req.EnvironmentId,
		Id:            req.Command.FeatureId,
	})
	if err != nil {
		return nil, err
	}
	return resp.Feature, nil
}

func (s *AutoOpsService) validateTargetFeature(
	ctx context.Context,
	f *featureproto.Feature,
	localizer locale.Localizer,
) error {
	if len(f.Variations) != 2 {
		dt, err := statusProgressiveRolloutInvalidVariationSize.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AutoOpsInvalidVariationSize),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if err := s.checkIfHasExperiment(ctx, f.Id, localizer); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) checkIfHasExperiment(
	ctx context.Context,
	featureID string,
	localizer locale.Localizer,
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
		dt, err := statusProgressiveRolloutInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if len(resp.Experiments) > 0 {
		dt, err := statusProgressiveRolloutWaitingOrRunningExperimentExists.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AutoOpsWaitingOrRunningExperimentExists),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AutoOpsService) validateProgressiveRolloutManualScheduleClause(
	clause *autoopsproto.ProgressiveRolloutManualScheduleClause,
	f *featureproto.Feature,
	localizer locale.Localizer,
) error {
	if err := s.validateProgressiveRolloutClauseVariationID(
		clause.VariationId,
		f,
		localizer,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseSchedules(
		clause.Schedules,
		localizer,
	); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateProgressiveRolloutTemplateScheduleClause(
	clause *autoopsproto.ProgressiveRolloutTemplateScheduleClause,
	f *featureproto.Feature,
	localizer locale.Localizer,
) error {
	if err := s.validateProgressiveRolloutClauseVariationID(
		clause.VariationId,
		f,
		localizer,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseSchedules(
		clause.Schedules,
		localizer,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseIncrements(
		clause.Increments,
		localizer,
	); err != nil {
		return err
	}
	if err := s.validateProgressiveRolloutClauseInterval(
		clause.Interval,
		localizer,
	); err != nil {
		return err
	}
	return nil
}

func (s *AutoOpsService) validateProgressiveRolloutClauseVariationID(
	variationID string,
	f *featureproto.Feature,
	localizer locale.Localizer,
) error {
	if variationID == "" {
		dt, err := statusProgressiveRolloutClauseVariationIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "variation_id"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	if exist := s.existVariationID(f, variationID); !exist {
		dt, err := statusProgressiveRolloutClauseInvalidVariationID.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "variation_id"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
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
	localizer locale.Localizer,
) error {
	if len(schedules) == 0 {
		dt, err := statusProgressiveRolloutClauseSchedulesRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "schedule"),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	for _, s := range schedules {
		if s.ExecuteAt == 0 {
			dt, err := statusProgressiveRolloutScheduleExecutedAtRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "execute_at"),
			})
			if err != nil {
				return statusProgressiveRolloutInternal.Err()
			}
			return dt.Err()
		}
		if s.Weight < 1 {
			dt, err := statusProgressiveRolloutScheduleInvalidWeight.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "weight"),
			})
			if err != nil {
				return statusProgressiveRolloutInternal.Err()
			}
			return dt.Err()
		}
	}
	if err := s.validateProgressiveRolloutClauseScheduleSpans(schedules, localizer); err != nil {
		return err
	}
	return nil
}

func (*AutoOpsService) validateProgressiveRolloutClauseIncrements(
	increments int64,
	localizer locale.Localizer,
) error {
	if increments < 1 {
		dt, err := statusProgressiveRolloutClauseInvalidIncrements.WithDetails(
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "increments"),
			},
		)
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (*AutoOpsService) validateProgressiveRolloutClauseInterval(
	interval autoopsproto.ProgressiveRolloutTemplateScheduleClause_Interval,
	localizer locale.Localizer,
) error {
	if interval == autoopsproto.ProgressiveRolloutTemplateScheduleClause_UNKNOWN {
		dt, err := statusProgressiveRolloutClauseUnknownInterval.WithDetails(
			&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "interval"),
			},
		)
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

// The span of time for each scheduled time must be at least 5 minutes.
func (*AutoOpsService) validateProgressiveRolloutClauseScheduleSpans(
	schedules []*autoopsproto.ProgressiveRolloutSchedule,
	localizer locale.Localizer,
) error {
	for i := 0; i < len(schedules); i++ {
		for j := i + 1; j < len(schedules); j++ {
			if schedules[j].ExecuteAt-schedules[i].ExecuteAt < int64(fiveMinutes.Seconds()) {
				dt, err := statusProgressiveRolloutInvalidScheduleSpans.WithDetails(
					&errdetails.LocalizedMessage{
						Locale:  localizer.GetLocale(),
						Message: localizer.MustLocalize(locale.AutoOpsInvalidScheduleSpans),
					},
				)
				if err != nil {
					return statusProgressiveRolloutInternal.Err()
				}
				return dt.Err()
			}
		}
	}
	return nil
}
