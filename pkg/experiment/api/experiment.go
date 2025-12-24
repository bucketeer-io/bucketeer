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
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	maxExperimentPeriodDays = 30
	maxExperimentPeriod     = maxExperimentPeriodDays * 24 * 60 * 60
)

func (s *experimentService) GetExperiment(
	ctx context.Context,
	req *proto.GetExperimentRequest,
) (*proto.GetExperimentResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetExperimentRequest(req); err != nil {
		return nil, err
	}
	experiment, err := s.experimentStorage.GetExperiment(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentNotFound) {
			return nil, statusExperimentNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.GetExperimentResponse{
		Experiment: experiment.Experiment,
	}, nil
}

func validateGetExperimentRequest(req *proto.GetExperimentRequest) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	return nil
}

func (s *experimentService) ListExperiments(
	ctx context.Context,
	req *proto.ListExperimentsRequest,
) (*proto.ListExperimentsResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "deleted",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "environment_id",
			Operator: mysql.OperatorEqual,
			Value:    req.EnvironmentId,
		},
	}
	if req.Archived != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "archived",
			Operator: mysql.OperatorEqual,
			Value:    req.Archived.Value,
		})
	}
	if req.FeatureId != "" {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature_id",
			Operator: mysql.OperatorEqual,
			Value:    req.FeatureId,
		})
	}
	if req.FeatureVersion != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "feature_version",
			Operator: mysql.OperatorEqual,
			Value:    req.FeatureVersion.Value,
		})
	}
	if req.StartAt != 0 {
		// When a start timestamp is provided,
		// use it as the lower bound for filtering.
		filters = append(filters, &mysql.FilterV2{
			Column:   "start_at",
			Operator: mysql.OperatorGreaterThanOrEqual,
			Value:    req.StartAt,
		})
	}
	if req.StopAt != 0 {
		// When a stop timestamp is provided:
		// - If req.StartAt is also provided, treat req.StopAt as an absolute upper bound.
		// (This selects experiments with stop_at <= req.StopAt.)
		// - If req.StartAt is not provided, treat req.StopAt as a relative cutoff timestamp.
		// (This selects experiments with stop_at >= req.StopAt.)
		// It treats it as a relative duration when the `startAt` is not provide
		if req.StartAt != 0 {
			filters = append(filters, &mysql.FilterV2{
				Column:   "stop_at",
				Operator: mysql.OperatorLessThanOrEqual,
				Value:    req.StopAt,
			})
		} else {
			filters = append(filters, &mysql.FilterV2{
				Column:   "stop_at",
				Operator: mysql.OperatorGreaterThanOrEqual,
				Value:    req.StopAt,
			})
		}
	}
	if req.Maintainer != "" {
		filters = append(filters, &mysql.FilterV2{
			Column:   "maintainer",
			Operator: mysql.OperatorEqual,
			Value:    req.Maintainer,
		})
	}
	var inFilters []*mysql.InFilter
	if len(req.Statuses) > 0 {
		statuses := make([]interface{}, 0, len(req.Statuses))
		for _, sts := range req.Statuses {
			statuses = append(statuses, sts)
		}
		inFilters = append(inFilters, &mysql.InFilter{
			Column: "status",
			Values: statuses,
		})
	}
	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"name", "description"},
			Keyword: req.SearchKeyword,
		}
	}
	orders, err := s.newExperimentListOrders(req.OrderBy, req.OrderDirection)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
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
		return nil, statusInvalidCursor.Err()
	}
	options := &mysql.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		Orders:      orders,
		InFilters:   inFilters,
		SearchQuery: searchQuery,
		NullFilters: nil,
		JSONFilters: nil,
	}
	experiments, nextCursor, totalCount, err := s.experimentStorage.ListExperiments(
		ctx,
		options,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list experiments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	summary, err := s.experimentStorage.GetExperimentSummary(ctx, req.EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to get experiment summary",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId",
					req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.ListExperimentsResponse{
		Experiments: experiments,
		Cursor:      strconv.Itoa(nextCursor),
		TotalCount:  totalCount,
		Summary: &proto.ListExperimentsResponse_Summary{
			TotalWaitingCount: summary.TotalWaitingCount,
			TotalRunningCount: summary.TotalRunningCount,
			TotalStoppedCount: summary.TotalStoppedCount,
		},
	}, nil
}

func (s *experimentService) newExperimentListOrders(
	orderBy proto.ListExperimentsRequest_OrderBy,
	orderDirection proto.ListExperimentsRequest_OrderDirection,
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
	case proto.ListExperimentsRequest_START_AT:
		column = "ex.start_at"
	case proto.ListExperimentsRequest_STOP_AT:
		column = "ex.stop_at"
	case proto.ListExperimentsRequest_STATUS:
		column = "ex.status"
	case proto.ListExperimentsRequest_GOALS_COUNT:
		column = "JSON_LENGTH(ex.goal_ids)"
	default:
		return nil, statusInvalidOrderBy.Err()
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
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createExperimentNoCommand(ctx, req, editor)
	}
	if err := validateCreateExperimentRequest(req); err != nil {
		return nil, err
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:            req.Command.FeatureId,
		EnvironmentId: req.EnvironmentId,
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, statusFeatureNotFound.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	for _, gid := range req.Command.GoalIds {
		_, err := s.getGoalMySQL(ctx, gid, req.EnvironmentId)
		if err != nil {
			if errors.Is(err, v2es.ErrGoalNotFound) {
				return nil, statusGoalNotFound.Err()
			}
			return nil, api.NewGRPCStatus(err).Err()
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
			"Failed to create experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("featureId", resp.Feature.Id),
				zap.String("baseVariationId", req.BaseVariationId),
				zap.Any("featureVariations", resp.Feature.Variations),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
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
		return s.experimentStorage.CreateExperiment(contextWithTx, experiment, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.CreateExperimentResponse{
		Experiment: experiment.Experiment,
	}, nil
}

func (s *experimentService) createExperimentNoCommand(
	ctx context.Context,
	req *proto.CreateExperimentRequest,
	editor *eventproto.Editor,
) (*proto.CreateExperimentResponse, error) {
	err := validateCreateExperimentRequestNoCommand(req)
	if err != nil {
		return nil, err
	}
	getFeatureResp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		Id:            req.FeatureId,
		EnvironmentId: req.EnvironmentId,
	})
	if err != nil {
		if code := status.Code(err); code == codes.NotFound {
			return nil, statusFeatureNotFound.Err()
		}
		s.logger.Error(
			"Failed to get feature",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
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
			"Failed to create experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("featureId", getFeatureResp.Feature.Id),
				zap.String("baseVariationId", req.BaseVariationId),
				zap.Any("featureVariations", getFeatureResp.Feature.Variations),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		for _, gid := range req.GoalIds {
			goal, err := s.getGoalMySQL(ctxWithTx, gid, req.EnvironmentId)
			if err != nil {
				return err
			}
			if goal.ConnectionType != proto.Goal_EXPERIMENT {
				return statusGoalTypeMismatch.Err()
			}
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
			nil,
		)
		if err != nil {
			return err
		}
		err = s.publisher.Publish(ctx, e)
		if err != nil {
			return err
		}
		return s.experimentStorage.CreateExperiment(ctxWithTx, experiment, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrGoalNotFound) {
			return nil, statusInvalidGoalID.Err()
		}
		if errors.Is(err, statusGoalTypeMismatch.Err()) {
			return nil, statusGoalTypeMismatch.Err()
		}
		if errors.Is(err, v2es.ErrExperimentAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.CreateExperimentResponse{
		Experiment: experiment.Experiment,
	}, nil
}

func validateCreateExperimentRequest(req *proto.CreateExperimentRequest) error {
	if req.Command.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if len(req.Command.GoalIds) == 0 {
		return statusGoalIDRequired.Err()
	}
	for _, gid := range req.Command.GoalIds {
		if gid == "" {
			return statusGoalIDRequired.Err()
		}
	}
	if err := validateExperimentPeriod(req.Command.StartAt, req.Command.StopAt); err != nil {
		return err
	}
	// TODO: validate name empty check
	return nil
}

func validateCreateExperimentRequestNoCommand(
	req *proto.CreateExperimentRequest,
) error {
	if req.FeatureId == "" {
		return statusFeatureIDRequired.Err()
	}
	if len(req.GoalIds) == 0 {
		return statusGoalIDRequired.Err()
	}
	for _, gid := range req.GoalIds {
		if gid == "" {
			return statusGoalIDRequired.Err()
		}
	}
	if err := validateExperimentPeriod(req.StartAt, req.StopAt); err != nil {
		return err
	}
	if req.Name == "" {
		return statusExperimentNameRequired.Err()
	}
	return nil
}

func validateExperimentPeriod(startAt, stopAt int64) error {
	period := stopAt - startAt
	if period <= 0 || period > int64(maxExperimentPeriod) {
		return statusExperimentPeriodOutOfRange.Err()
	}
	return nil
}

func (s *experimentService) UpdateExperiment(
	ctx context.Context,
	req *proto.UpdateExperimentRequest,
) (*proto.UpdateExperimentResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.ChangeExperimentPeriodCommand == nil &&
		req.ChangeNameCommand == nil &&
		req.ChangeDescriptionCommand == nil {
		return s.updateExperimentNoCommand(ctx, req, editor)
	}
	if err := validateUpdateExperimentRequest(req); err != nil {
		return nil, err
	}
	var experimentPb *proto.Experiment
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		experiment, err := s.experimentStorage.GetExperiment(contextWithTx, req.Id, req.EnvironmentId)
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
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", req.EnvironmentId),
					)...,
				)
				return err
			}
			return s.experimentStorage.UpdateExperiment(contextWithTx, experiment, req.EnvironmentId)
		}
		if req.ChangeNameCommand != nil {
			if err = handler.Handle(ctx, req.ChangeNameCommand); err != nil {
				s.logger.Error(
					"Failed to change Name",
					log.FieldsFromIncomingContext(ctx).AddFields(
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
					log.FieldsFromIncomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("environmentId", req.EnvironmentId),
					)...,
				)
				return err
			}
		}
		experimentPb = experiment.Experiment
		return s.experimentStorage.UpdateExperiment(contextWithTx, experiment, req.EnvironmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentNotFound) || errors.Is(err, v2es.ErrExperimentUnexpectedAffectedRows) {
			return nil, statusExperimentNotFound.Err()
		}
		s.logger.Error(
			"Failed to update experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.UpdateExperimentResponse{
		Experiment: experimentPb,
	}, nil
}

func (s *experimentService) updateExperimentNoCommand(
	ctx context.Context,
	req *proto.UpdateExperimentRequest,
	editor *eventproto.Editor,
) (*proto.UpdateExperimentResponse, error) {
	err := validateUpdateExperimentNoCommandRequest(req)
	if err != nil {
		s.logger.Error(
			"Failed validate update experiment no command req",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
			return nil, statusExperimentNotFound.Err()
		}
		s.logger.Error(
			"Failed to update experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &proto.UpdateExperimentResponse{
		Experiment: experimentPb,
	}, nil
}

func validateUpdateExperimentRequest(req *proto.UpdateExperimentRequest) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	if req.ChangeExperimentPeriodCommand != nil {
		if err := validateExperimentPeriod(
			req.ChangeExperimentPeriodCommand.StartAt,
			req.ChangeExperimentPeriodCommand.StopAt,
		); err != nil {
			return err
		}
	}
	return nil
}

func validateUpdateExperimentNoCommandRequest(
	req *proto.UpdateExperimentRequest,
) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	if req.Name != nil && req.Name.Value == "" {
		return statusExperimentNameRequired.Err()
	}
	if (req.StartAt != nil && req.StopAt == nil) ||
		(req.StartAt == nil && req.StopAt != nil) {
		return statusExperimentPeriodInvalid.Err()
	}
	if req.StartAt != nil && req.StopAt != nil {
		if err := validateExperimentPeriod(
			req.StartAt.Value,
			req.StopAt.Value,
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
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateStartExperimentRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId); err != nil {
		return nil, err
	}
	return &proto.StartExperimentResponse{}, nil
}

func validateStartExperimentRequest(req *proto.StartExperimentRequest) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *experimentService) FinishExperiment(
	ctx context.Context,
	req *proto.FinishExperimentRequest,
) (*proto.FinishExperimentResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateFinishExperimentRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId); err != nil {
		return nil, err
	}
	return &proto.FinishExperimentResponse{}, nil
}

func validateFinishExperimentRequest(req *proto.FinishExperimentRequest) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *experimentService) StopExperiment(
	ctx context.Context,
	req *proto.StopExperimentRequest,
) (*proto.StopExperimentResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateStopExperimentRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId); err != nil {
		return nil, err
	}
	return &proto.StopExperimentResponse{}, nil
}

func validateStopExperimentRequest(req *proto.StopExperimentRequest) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *experimentService) ArchiveExperiment(
	ctx context.Context,
	req *proto.ArchiveExperimentRequest,
) (*proto.ArchiveExperimentResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, statusExperimentIDRequired.Err()
	}
	if req.Command == nil {
		return nil, statusNoCommand.Err()
	}
	err = s.updateExperiment(
		ctx,
		editor,
		req.Command,
		req.Id,
		req.EnvironmentId,
	)
	if err != nil {
		s.logger.Error(
			"Failed to archive experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteExperimentRequest(req); err != nil {
		return nil, err
	}
	if err := s.updateExperiment(ctx, editor, req.Command, req.Id, req.EnvironmentId); err != nil {
		return nil, err
	}
	return &proto.DeleteExperimentResponse{}, nil
}

func validateDeleteExperimentRequest(req *proto.DeleteExperimentRequest) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	if req.Command == nil {
		return statusNoCommand.Err()
	}
	return nil
}

func (s *experimentService) updateExperiment(
	ctx context.Context,
	editor *eventproto.Editor,
	cmd command.Command,
	id, environmentId string,
) error {
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		experiment, err := s.experimentStorage.GetExperiment(contextWithTx, id, environmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get experiment",
				log.FieldsFromIncomingContext(ctx).AddFields(
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
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return err
		}
		return s.experimentStorage.UpdateExperiment(contextWithTx, experiment, environmentId)
	})
	if err != nil {
		if errors.Is(err, v2es.ErrExperimentNotFound) || errors.Is(err, v2es.ErrExperimentUnexpectedAffectedRows) {
			return statusExperimentNotFound.Err()
		}
		s.logger.Error(
			"Failed to update experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}
