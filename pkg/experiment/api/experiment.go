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

	pb "github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/experiment/command"
	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	maxExperimentPeriodDays = 30
	maxExperimentPeriod     = maxExperimentPeriodDays * 24 * 60 * 60
)

func (s *experimentService) GetExperiment(
	ctx context.Context,
	req *proto.GetExperimentRequest,
) (*proto.GetExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetExperimentRequest(req, localizer); err != nil {
		return nil, err
	}
	experimentStorage := v2es.NewExperimentStorage(s.mysqlClient)
	experiment, err := experimentStorage.GetExperiment(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &proto.GetExperimentResponse{
		Experiment: experiment.Experiment,
	}, nil
}

func validateGetExperimentRequest(req *proto.GetExperimentRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *experimentService) ListExperiments(
	ctx context.Context,
	req *proto.ListExperimentsRequest,
) (*proto.ListExperimentsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	if req.Archived != nil {
		whereParts = append(whereParts, mysql.NewFilter("archived", "=", req.Archived.Value))
	}
	if req.FeatureId != "" {
		whereParts = append(whereParts, mysql.NewFilter("feature_id", "=", req.FeatureId))
	}
	if req.FeatureVersion != nil {
		whereParts = append(whereParts, mysql.NewFilter("feature_version", "=", req.FeatureVersion.Value))
	}
	if req.From != 0 {
		whereParts = append(whereParts, mysql.NewFilter("stopped_at", ">=", req.From))
	}
	if req.To != 0 {
		whereParts = append(whereParts, mysql.NewFilter("start_at", "<=", req.To))
	}
	if req.Status != nil {
		whereParts = append(whereParts, mysql.NewFilter("status", "=", req.Status.Value))
	} else if len(req.Statuses) > 0 {
		statuses := make([]interface{}, 0, len(req.Statuses))
		for _, sts := range req.Statuses {
			statuses = append(statuses, sts)
		}
		whereParts = append(whereParts, mysql.NewInFilter("status", statuses))
	}
	if req.Maintainer != "" {
		whereParts = append(whereParts, mysql.NewFilter("maintainer", "=", req.Maintainer))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"name", "description"}, req.SearchKeyword))
	}
	orders, err := s.newExperimentListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
	experimentStorage := v2es.NewExperimentStorage(s.mysqlClient)
	experiments, nextCursor, totalCount, summary, err := experimentStorage.ListExperiments(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list experiments",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	return &proto.ListExperimentsResponse{
		Experiments: experiments,
		Cursor:      strconv.Itoa(nextCursor),
		TotalCount:  totalCount,
		Summary:     summary,
	}, nil
}

func (s *experimentService) newExperimentListOrders(
	orderBy proto.ListExperimentsRequest_OrderBy,
	orderDirection proto.ListExperimentsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case proto.ListExperimentsRequest_DEFAULT,
		proto.ListExperimentsRequest_NAME:
		column = "ex.name"
	case proto.ListExperimentsRequest_CREATED_AT:
		column = "ex.created_at"
	case proto.ListExperimentsRequest_UPDATED_AT:
		column = "ex.updated_at"
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
	if orderDirection == proto.ListExperimentsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *experimentService) CreateExperiment(
	ctx context.Context,
	req *proto.CreateExperimentRequest,
) (*proto.CreateExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createExperimentNoCommand(ctx, req, editor, localizer)
	}
	if err := validateCreateExperimentRequest(req, localizer); err != nil {
		return nil, err
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:            req.Command.FeatureId,
		EnvironmentId: req.EnvironmentId,
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			dt, err := statusFeatureNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	for _, gid := range req.Command.GoalIds {
		goal, err := s.getGoalMySQL(ctx, gid, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2es.ErrGoalNotFound) {
				dt, err := statusGoalNotFound.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalize(locale.NotFoundError),
				})
				if err != nil {
					return nil, statusInternal.Err()
				}
				return nil, dt.Err()
			}
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if goal.ConnectionType != proto.Goal_EXPERIMENT {
			dt, err := statusGoalTypeMismatch.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "goal_ids"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	experiment, err := domain.NewExperiment(
		req.Command.FeatureId,
		resp.Feature.Version,
		resp.Feature.Variations,
		req.Command.GoalIds,
		req.Command.StartAt,
		req.Command.StopAt,
		req.Command.Name,
		req.Command.Description,
		req.Command.BaseVariationId,
		editor.Email,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
		experimentStorage := v2es.NewExperimentStorage(tx)
		handler, err := command.NewExperimentCommandHandler(
			editor,
			experiment,
			s.publisher,
			req.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return experimentStorage.CreateExperiment(ctx, experiment, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentAlreadyExists) {
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
			"Failed to create experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	return &proto.CreateExperimentResponse{
		Experiment: experiment.Experiment,
	}, nil
}

func (s *experimentService) createExperimentNoCommand(
	ctx context.Context,
	req *proto.CreateExperimentRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*proto.CreateExperimentResponse, error) {
	err := validateCreateExperimentRequestNoCommand(req, localizer)
	if err != nil {
		return nil, err
	}
	getFeatureResp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:            req.FeatureId,
		EnvironmentId: req.EnvironmentId,
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			dt, err := statusFeatureNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	for _, gid := range req.GoalIds {
		goal, err := s.getGoalMySQL(ctx, gid, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2es.ErrGoalNotFound) {
				dt, err := statusGoalNotFound.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalize(locale.NotFoundError),
				})
				if err != nil {
					return nil, statusInternal.Err()
				}
				return nil, dt.Err()
			}
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if goal.ConnectionType != proto.Goal_EXPERIMENT {
			dt, err := statusGoalTypeMismatch.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "goal_ids"),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	experiment, err := domain.NewExperiment(
		req.FeatureId,
		getFeatureResp.Feature.Version,
		getFeatureResp.Feature.Variations,
		req.GoalIds,
		req.StartAt,
		req.StopAt,
		req.Name,
		req.Description,
		req.BaseVariationId,
		editor.Email,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		experimentStorage := v2es.NewExperimentStorage(s.mysqlClient)
		prev := &domain.Experiment{}
		if err = copier.Copy(prev, experiment); err != nil {
			return err
		}
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_EXPERIMENT,
			experiment.Id,
			eventproto.Event_EXPERIMENT_CREATED,
			&eventproto.ExperimentCreatedEvent{
				Id:              experiment.Id,
				FeatureId:       experiment.FeatureId,
				FeatureVersion:  experiment.FeatureVersion,
				Variations:      experiment.Variations,
				GoalIds:         experiment.GoalIds,
				StartAt:         experiment.StartAt,
				StopAt:          experiment.StopAt,
				StoppedAt:       experiment.StoppedAt,
				CreatedAt:       experiment.CreatedAt,
				UpdatedAt:       experiment.UpdatedAt,
				Name:            experiment.Name,
				Description:     experiment.Description,
				BaseVariationId: experiment.BaseVariationId,
			},
			req.EnvironmentId,
			experiment.Experiment,
			prev,
		)
		if err != nil {
			return err
		}
		err = s.publisher.Publish(ctx, e)
		if err != nil {
			return err
		}
		return experimentStorage.CreateExperiment(ctxWithTx, experiment, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentAlreadyExists) {
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
			"Failed to create experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	return &proto.CreateExperimentResponse{
		Experiment: experiment.Experiment,
	}, nil
}

func validateCreateExperimentRequest(req *proto.CreateExperimentRequest, localizer locale.Localizer) error {
	if req.Command.FeatureId == "" {
		dt, err := statusFeatureIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.Command.GoalIds) == 0 {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, gid := range req.Command.GoalIds {
		if gid == "" {
			dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
	}
	if err := validateExperimentPeriod(req.Command.StartAt, req.Command.StopAt, localizer); err != nil {
		return err
	}
	// TODO: validate name empty check
	return nil
}

func validateCreateExperimentRequestNoCommand(
	req *proto.CreateExperimentRequest,
	localizer locale.Localizer,
) error {
	if req.FeatureId == "" {
		dt, err := statusFeatureIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "feature_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.GoalIds) == 0 {
		dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	for _, gid := range req.GoalIds {
		if gid == "" {
			dt, err := statusGoalIDRequired.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "goal_id"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
	}
	if err := validateExperimentPeriod(req.StartAt, req.StopAt, localizer); err != nil {
		return err
	}
	if req.Name == "" {
		dt, err := statusExperimentNameRequired.WithDetails(&errdetails.LocalizedMessage{
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

func validateExperimentPeriod(startAt, stopAt int64, localizer locale.Localizer) error {
	period := stopAt - startAt
	if period <= 0 || period > int64(maxExperimentPeriod) {
		dt, err := statusPeriodTooLong.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "period"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *experimentService) UpdateExperiment(
	ctx context.Context,
	req *proto.UpdateExperimentRequest,
) (*proto.UpdateExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.ChangeExperimentPeriodCommand == nil &&
		req.ChangeNameCommand == nil &&
		req.ChangeDescriptionCommand == nil {
		return s.updateExperimentNoCommand(ctx, req, editor, localizer)
	}
	if err := validateUpdateExperimentRequest(req, localizer); err != nil {
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
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	var experimentPb *proto.Experiment
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		experimentStorage := v2es.NewExperimentStorage(tx)
		experiment, err := experimentStorage.GetExperiment(ctx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		handler, err := command.NewExperimentCommandHandler(
			editor,
			experiment,
			s.publisher,
			req.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if req.ChangeExperimentPeriodCommand != nil {
			if err = handler.Handle(ctx, req.ChangeExperimentPeriodCommand); err != nil {
				s.logger.Error(
					"Failed to change period",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", req.EnvironmentId),
					)...,
				)
				return err
			}
			return experimentStorage.UpdateExperiment(ctx, experiment, req.EnvironmentId)
		}
		if req.ChangeNameCommand != nil {
			if err = handler.Handle(ctx, req.ChangeNameCommand); err != nil {
				s.logger.Error(
					"Failed to change Name",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", req.EnvironmentId),
					)...,
				)
				return err
			}
		}
		if req.ChangeDescriptionCommand != nil {
			if err = handler.Handle(ctx, req.ChangeDescriptionCommand); err != nil {
				s.logger.Error(
					"Failed to change Description",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", req.EnvironmentId),
					)...,
				)
				return err
			}
		}
		experimentPb = experiment.Experiment
		return experimentStorage.UpdateExperiment(ctx, experiment, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentNotFound) || errors.Is(err, v2es.ErrExperimentUnexpectedAffectedRows) {
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
			"Failed to update experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	return &proto.UpdateExperimentResponse{
		Experiment: experimentPb,
	}, nil
}

func (s *experimentService) updateExperimentNoCommand(
	ctx context.Context,
	req *proto.UpdateExperimentRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*proto.UpdateExperimentResponse, error) {
	err := validateUpdateExperimentNoCommandRequest(req, localizer)
	if err != nil {
		s.logger.Error(
			"Failed validate update experiment no command req",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}

	var experimentPb *proto.Experiment
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		experimentStorage := v2es.NewExperimentStorage(s.mysqlClient)
		experiment, err := experimentStorage.GetExperiment(ctxWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		updated, err := experiment.Update(
			req.Name,
			req.Description,
			req.StartAt,
			req.StopAt,
			req.Status,
			req.Archived,
		)
		if err != nil {
			return err
		}

		var eventMsg pb.Message
		if req.Archived != nil {
			if experiment.Status == proto.Experiment_RUNNING {
				return v2es.ErrExperimentCannotBeArchived
			}
			eventMsg = &eventproto.ExperimentArchivedEvent{
				Id: req.Id,
			}
		} else {
			eventMsg = &eventproto.ExperimentUpdatedEvent{
				Id:          experiment.Id,
				Name:        updated.Name,
				Description: updated.Description,
				StartAt:     updated.StartAt,
				StopAt:      updated.StopAt,
				Status:      updated.Status,
			}
		}
		event, err := domainevent.NewEvent(
			editor,
			eventproto.Event_EXPERIMENT,
			experiment.Id,
			eventproto.Event_EXPERIMENT_UPDATED,
			eventMsg,
			req.EnvironmentId,
			updated,
			experiment,
		)
		if err != nil {
			return err
		}
		if err := s.publisher.Publish(ctxWithTx, event); err != nil {
			return err
		}
		experimentPb = updated.Experiment
		return experimentStorage.UpdateExperiment(ctxWithTx, updated, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentNotFound) || errors.Is(err, v2es.ErrExperimentUnexpectedAffectedRows) {
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
			"Failed to update experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	return &proto.UpdateExperimentResponse{
		Experiment: experimentPb,
	}, nil
}

func validateUpdateExperimentRequest(req *proto.UpdateExperimentRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.ChangeExperimentPeriodCommand != nil {
		if err := validateExperimentPeriod(
			req.ChangeExperimentPeriodCommand.StartAt,
			req.ChangeExperimentPeriodCommand.StopAt,
			localizer,
		); err != nil {
			return err
		}
	}
	return nil
}

func validateUpdateExperimentNoCommandRequest(
	req *proto.UpdateExperimentRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Name != nil && req.Name.Value == "" {
		dt, err := statusExperimentNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if (req.StartAt != nil && req.StopAt == nil) ||
		(req.StartAt == nil && req.StopAt != nil) {
		dt, err := statusPeriodInvalid.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "period"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.StartAt != nil && req.StopAt != nil {
		if err := validateExperimentPeriod(
			req.StartAt.Value,
			req.StopAt.Value,
			localizer,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *experimentService) StartExperiment(
	ctx context.Context,
	req *proto.StartExperimentRequest,
) (*proto.StartExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateStartExperimentRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId, localizer); err != nil {
		return nil, err
	}
	return &proto.StartExperimentResponse{}, nil
}

func validateStartExperimentRequest(req *proto.StartExperimentRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *experimentService) FinishExperiment(
	ctx context.Context,
	req *proto.FinishExperimentRequest,
) (*proto.FinishExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateFinishExperimentRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId, localizer); err != nil {
		return nil, err
	}
	return &proto.FinishExperimentResponse{}, nil
}

func validateFinishExperimentRequest(req *proto.FinishExperimentRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *experimentService) StopExperiment(
	ctx context.Context,
	req *proto.StopExperimentRequest,
) (*proto.StopExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateStopExperimentRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId, localizer); err != nil {
		return nil, err
	}
	return &proto.StopExperimentResponse{}, nil
}

func validateStopExperimentRequest(req *proto.StopExperimentRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *experimentService) ArchiveExperiment(
	ctx context.Context,
	req *proto.ArchiveExperimentRequest,
) (*proto.ArchiveExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "experiment_id"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if req.Command == nil {
		dt, err := statusNoCommand.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.updateExperiment(
		ctx,
		editor,
		req.Command,
		req.Id,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		s.logger.Error(
			"Failed to archive experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	return &proto.ArchiveExperimentResponse{}, nil
}

func (s *experimentService) DeleteExperiment(
	ctx context.Context,
	req *proto.DeleteExperimentRequest,
) (*proto.DeleteExperimentResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteExperimentRequest(req, localizer); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId, localizer); err != nil {
		return nil, err
	}
	return &proto.DeleteExperimentResponse{}, nil
}

func validateDeleteExperimentRequest(req *proto.DeleteExperimentRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusExperimentIDRequired.WithDetails(&errdetails.LocalizedMessage{
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

func (s *experimentService) updateExperiment(
	ctx context.Context,
	editor *eventproto.Editor,
	cmd command.Command,
	id, environmentId string,
	localizer locale.Localizer,
) error {
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
		experimentStorage := v2es.NewExperimentStorage(tx)
		experiment, err := experimentStorage.GetExperiment(ctx, id, environmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get experiment",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return err
		}
		handler, err := command.NewExperimentCommandHandler(editor, experiment, s.publisher, environmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, cmd); err != nil {
			s.logger.Error(
				"Failed to handle command",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return err
		}
		return experimentStorage.UpdateExperiment(ctx, experiment, environmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentNotFound) || errors.Is(err, v2es.ErrExperimentUnexpectedAffectedRows) {
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
			"Failed to update experiment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
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
