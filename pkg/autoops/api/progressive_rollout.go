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
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/golang/protobuf/ptypes"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/command"
	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	fiveMinutes     = 5 * time.Minute
	listRequestSize = 500
)

var (
	errProgressiveRolloutAutoOpsHasWebhook = errors.New(
		"autoops: can not create a progressive rollout when the webhook is set in the auto ops",
	)
	errProgressiveRolloutAutoOpsHasDatetime = errors.New(
		"autoops: can not create a progressive rollout when the schedule is set in the auto ops",
	)
)

func (s *AutoOpsService) CreateProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
) (*autoopsproto.CreateProgressiveRolloutResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreateProgressiveRolloutRequest(ctx, req, localizer); err != nil {
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		// We validate auto ops rules here since it's not possible to mock `ListAutoOpsRules`.
		autoOpsRuleStorage := v2as.NewAutoOpsRuleStorage(tx)
		if err := s.validateTargetAutoOpsRules(ctx, req, localizer, autoOpsRuleStorage); err != nil {
			return err
		}
		progressiveRollout, err := domain.NewProgressiveRollout(
			req.Command.FeatureId,
			req.Command.ProgressiveRolloutManualScheduleClause,
			req.Command.ProgressiveRolloutTemplateScheduleClause,
		)
		if err != nil {
			return err
		}
		storage := v2as.NewProgressiveRolloutStorage(tx)
		handler := command.NewProgressiveRolloutCommandHandler(
			editor,
			progressiveRollout,
			s.publisher,
			req.EnvironmentNamespace,
		)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return storage.CreateProgressiveRollout(ctx, progressiveRollout, req.EnvironmentNamespace)
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
		case errProgressiveRolloutAutoOpsHasWebhook:
			dt, err := statusProgressiveRolloutAutoOpsHasWebhook.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AutoOpsHasWebhook),
			})
			if err != nil {
				return nil, statusProgressiveRolloutInternal.Err()
			}
			return nil, dt.Err()
		case errProgressiveRolloutAutoOpsHasDatetime:
			dt, err := statusProgressiveRolloutAutoOpsHasDatetime.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AutoOpsHasDatetime),
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateGetProgressiveRolloutRequest(req, localizer); err != nil {
		return nil, err
	}
	storage := v2as.NewProgressiveRolloutStorage(s.mysqlClient)
	progressiveRollout, err := storage.GetProgressiveRollout(ctx, req.Id, req.EnvironmentNamespace)
	if err != nil {
		s.logger.Error(
			"Failed to get progressive rollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		if err == v2as.ErrProgressiveRolloutNotFound {
			dt, err := statusProgressiveRolloutNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "progressive_rollout"),
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

func (s *AutoOpsService) DeleteProgressiveRollout(
	ctx context.Context,
	req *autoopsproto.DeleteProgressiveRolloutRequest,
) (*autoopsproto.DeleteProgressiveRolloutResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateDeleteProgressiveRolloutRequest(req, localizer); err != nil {
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2as.NewProgressiveRolloutStorage(tx)
		progressiveRollout, err := storage.GetProgressiveRollout(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		handler := command.NewProgressiveRolloutCommandHandler(
			editor,
			progressiveRollout,
			s.publisher,
			req.EnvironmentNamespace,
		)
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return storage.DeleteProgressiveRollout(ctx, req.Id, req.EnvironmentNamespace)
	})
	if err != nil {
		s.logger.Error(
			"Failed to delete ProgressiveRollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		if err == v2as.ErrProgressiveRolloutNotFound || err == v2as.ErrProgressiveRolloutUnexpectedAffectedRows {
			dt, err := statusProgressiveRolloutNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "progressive_rollout"),
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
	_, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_VIEWER, req.EnvironmentNamespace, localizer)
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
	editor, err := s.checkRole(ctx, accountproto.AccountV2_Role_Environment_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := s.validateExecuteProgressiveRolloutRequest(req, localizer); err != nil {
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2as.NewProgressiveRolloutStorage(tx)
		progressiveRollout, err := storage.GetProgressiveRollout(ctx, req.Id, req.EnvironmentNamespace)
		if err != nil {
			return err
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
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return nil
		}
		handler := command.NewProgressiveRolloutCommandHandler(
			editor,
			progressiveRollout,
			s.publisher,
			req.EnvironmentNamespace,
		)
		if err := handler.Handle(ctx, req.ChangeProgressiveRolloutTriggeredAtCommand); err != nil {
			return err
		}
		if err := storage.UpdateProgressiveRollout(ctx, progressiveRollout, req.EnvironmentNamespace); err != nil {
			return err
		}
		return ExecuteProgressiveRolloutOperation(
			ctx,
			progressiveRollout,
			s.featureClient,
			req.ChangeProgressiveRolloutTriggeredAtCommand.ScheduleId,
			req.EnvironmentNamespace,
		)
	})
	if err != nil {
		s.logger.Error(
			"Failed to execute progressiveRollout",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		if err == v2as.ErrProgressiveRolloutNotFound || err == v2as.ErrProgressiveRolloutUnexpectedAffectedRows {
			dt, err := statusProgressiveRolloutNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.NotFoundError, "progressive_rollout"),
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
	return &autoopsproto.ExecuteProgressiveRolloutResponse{}, nil
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
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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

func (s *AutoOpsService) listAutoOpsRulesByFeatureID(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
	localizer locale.Localizer,
	storage v2as.AutoOpsRuleStorage,
) ([]*autoopsproto.AutoOpsRule, error) {
	allRules := []*autoopsproto.AutoOpsRule{}
	cursor := ""
	for {
		rules, c, err := s.listAutoOpsRules(
			ctx,
			listRequestSize,
			cursor,
			[]string{req.Command.FeatureId},
			req.EnvironmentNamespace,
			localizer,
			storage,
		)
		if err != nil {
			return nil, err
		}
		allRules = append(allRules, rules...)
		size := len(rules)
		if size == 0 || size < listRequestSize {
			return allRules, nil
		}
		cursor = c
	}
}

func (s *AutoOpsService) getFeature(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
	localizer locale.Localizer,
) (*featureproto.Feature, error) {
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentNamespace: req.EnvironmentNamespace,
		Id:                   req.Command.FeatureId,
	})
	if err != nil {
		return nil, err
	}
	return resp.Feature, nil
}

func (s *AutoOpsService) validateTargetAutoOpsRules(
	ctx context.Context,
	req *autoopsproto.CreateProgressiveRolloutRequest,
	localizer locale.Localizer,
	storage v2as.AutoOpsRuleStorage,
) error {
	rules, err := s.listAutoOpsRulesByFeatureID(
		ctx,
		req,
		localizer,
		storage,
	)
	if err != nil {
		return err
	}
	for _, r := range rules {
		if r.TriggeredAt > 0 {
			continue
		}
		for _, c := range r.Clauses {
			// Return an error when Clause is DatetimeClause or WebhookClause.
			if ptypes.Is(c.Clause, domain.DatetimeClause) {
				return errProgressiveRolloutAutoOpsHasDatetime
			}
			if ptypes.Is(c.Clause, domain.WebhookClause) {
				return errProgressiveRolloutAutoOpsHasWebhook
			}
		}
	}
	return nil
}

func (s *AutoOpsService) validateTargetFeature(
	ctx context.Context,
	f *featureproto.Feature,
	localizer locale.Localizer,
) error {
	if !f.Enabled {
		dt, err := statusProgressiveRolloutFeatureDisabled.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AutoOpsFeatureDisabled),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
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
	if len(f.Prerequisites) > 0 {
		dt, err := statusProgressiveRolloutFeatureHasPrerequisitess.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AutoOpsFeatureHasPrerequisites),
		})
		if err != nil {
			return statusProgressiveRolloutInternal.Err()
		}
		return dt.Err()
	}
	for _, t := range f.Targets {
		if len(t.Users) > 0 {
			dt, err := statusProgressiveRolloutFeatureHasIndividualTargeting.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AutoOpsFeatureHasIndividualTargeting),
			})
			if err != nil {
				return statusProgressiveRolloutInternal.Err()
			}
			return dt.Err()
		}
	}
	if len(f.Rules) > 0 {
		dt, err := statusProgressiveRolloutFeatureHasRules.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AutoOpsFeatureHasRules),
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
