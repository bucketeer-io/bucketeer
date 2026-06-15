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
	pb "google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
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
	params := v2es.ListExperimentsParams{
		EnvironmentID:  req.EnvironmentId,
		FeatureID:      req.FeatureId,
		StartAt:        req.StartAt,
		StopAt:         req.StopAt,
		Maintainer:     req.Maintainer,
		Statuses:       req.Statuses,
		SearchKeyword:  req.SearchKeyword,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		PageSize:       int(req.PageSize),
		Cursor:         req.Cursor,
	}
	if req.Archived != nil {
		params.Archived = &req.Archived.Value
	}
	if req.FeatureVersion != nil {
		params.FeatureVersion = &req.FeatureVersion.Value
	}
	experiments, nextCursor, totalCount, err := s.experimentStorage.ListExperiments(
		ctx,
		params,
	)
	if err != nil {
		if errors.Is(err, v2es.ErrInvalidCursor) {
			return nil, statusInvalidCursor.Err()
		}
		if errors.Is(err, v2es.ErrInvalidOrderBy) {
			return nil, statusInvalidOrderBy.Err()
		}
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
	err = validateCreateExperimentRequest(req)
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
	err = s.dbClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context) error {
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

func validateCreateExperimentRequest(
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
	err = validateUpdateExperimentRequest(req)
	if err != nil {
		s.logger.Error(
			"Failed validate update experiment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}

	var experimentPb *proto.Experiment
	err = s.dbClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context) error {
		experiment, err := s.experimentStorage.GetExperiment(ctxWithTx, req.Id, req.EnvironmentId)
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
		return s.experimentStorage.UpdateExperiment(ctxWithTx, updated, req.EnvironmentId)
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

func validateUpdateExperimentRequest(
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
	if err = validateDeleteExperimentRequest(req); err != nil {
		return nil, err
	}

	var experimentPb *domain.Experiment
	err = s.dbClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context) error {
		experiment, err := s.experimentStorage.GetExperiment(ctxWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get experiment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("experimentId", req.Id),
				)...,
			)
			return err
		}
		experimentPb = experiment

		err = experiment.SetDeleted()
		if err != nil {
			s.logger.Error(
				"Failed to set deleted",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
					zap.String("experimentId", req.Id),
				)...,
			)
			return err
		}

		return s.experimentStorage.UpdateExperiment(ctxWithTx, experiment, req.EnvironmentId)
	})
	if err != nil {
		s.logger.Error(
			"Failed to delete experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("experimentId", req.Id),
			)...,
		)
		if errors.Is(err, v2es.ErrExperimentNotFound) ||
			errors.Is(err, v2es.ErrExperimentUnexpectedAffectedRows) {
			return nil, statusExperimentNotFound.Err()
		}
		return nil, api.NewGRPCStatus(err).Err()
	}

	event, err := domainevent.NewEvent(
		editor,
		eventproto.Event_EXPERIMENT,
		req.Id,
		eventproto.Event_EXPERIMENT_DELETED,
		&eventproto.ExperimentDeletedEvent{
			Id: req.Id,
		},
		req.EnvironmentId,
		nil,
		experimentPb,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("experimentId", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("experimentId", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	return &proto.DeleteExperimentResponse{}, nil
}

func validateDeleteExperimentRequest(req *proto.DeleteExperimentRequest) error {
	if req.Id == "" {
		return statusExperimentIDRequired.Err()
	}
	return nil
}
