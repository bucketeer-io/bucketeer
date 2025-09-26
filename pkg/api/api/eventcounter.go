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

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	"github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func (s *grpcGatewayService) GetExperimentEvaluationCount(
	ctx context.Context,
	req *gwproto.GetExperimentEvaluationCountRequest,
) (*gwproto.GetExperimentEvaluationCountResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetExperimentEvaluationCount request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetExperimentEvaluationResult, "").Inc()

	resp, err := s.eventCounterClient.GetExperimentEvaluationCount(
		ctx, &eventcounter.GetExperimentEvaluationCountRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			StartAt:        req.StartAt,
			EndAt:          req.EndAt,
			FeatureId:      req.FeatureId,
			FeatureVersion: req.FeatureVersion,
			VariationIds:   req.VariationIds,
		})
	if err != nil {
		s.logger.Error("Failed to get experiment evaluation count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
				zap.Int64("startAt", req.StartAt),
				zap.Int64("endAt", req.EndAt),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ExperimentEvaluationCount response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", req.FeatureId),
				zap.Int64("startAt", req.StartAt),
				zap.Int64("endAt", req.EndAt),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetExperimentEvaluationCountResponse{
		FeatureId:       resp.FeatureId,
		FeatureVersion:  resp.FeatureVersion,
		VariationCounts: resp.VariationCounts,
	}, nil
}

func (s *grpcGatewayService) GetEvaluationTimeseriesCount(
	ctx context.Context,
	req *gwproto.GetEvaluationTimeseriesCountRequest,
) (*gwproto.GetEvaluationTimeseriesCountResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetEvaluationTimeseriesCount request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetEvaluationTimeseriesCount, "",
	).Inc()

	resp, err := s.eventCounterClient.GetEvaluationTimeseriesCount(
		ctx, &eventcounter.GetEvaluationTimeseriesCountRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			FeatureId:     req.FeatureId,
			TimeRange:     req.TimeRange,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get evaluation timeseries count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("EvaluationTimeseriesCount response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetEvaluationTimeseriesCountResponse{
		UserCounts:  resp.UserCounts,
		EventCounts: resp.EventCounts,
	}, nil
}

func (s *grpcGatewayService) GetExperimentResult(
	ctx context.Context,
	req *gwproto.GetExperimentResultRequest,
) (*gwproto.GetExperimentResultResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetExperimentResult request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetExperimentResult, "",
	).Inc()

	resp, err := s.eventCounterClient.GetExperimentResult(
		ctx, &eventcounter.GetExperimentResultRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			ExperimentId:  req.ExperimentId,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get experiment result",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("experimentId", req.ExperimentId),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ExperimentResult response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("experimentId", req.ExperimentId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetExperimentResultResponse{
		ExperimentResult: resp.ExperimentResult,
	}, nil
}
func (s *grpcGatewayService) ListExperimentResults(
	ctx context.Context,
	req *gwproto.ListExperimentResultsRequest,
) (*gwproto.ListExperimentResultsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListExperimentResults request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListExperimentResults, "",
	).Inc()

	resp, err := s.eventCounterClient.ListExperimentResults(
		ctx, &eventcounter.ListExperimentResultsRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			FeatureId:      req.FeatureId,
			FeatureVersion: req.FeatureVersion,
		},
	)
	if err != nil {
		s.logger.Error("Failed to list experiment results",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ListExperimentResults response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields()...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListExperimentResultsResponse{
		Results: resp.Results,
	}, nil
}

func (s *grpcGatewayService) GetExperimentGoalCount(
	ctx context.Context,
	req *gwproto.GetExperimentGoalCountRequest,
) (*gwproto.GetExperimentGoalCountResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetExperimentGoalCount request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetExperimentGoalCount, "",
	).Inc()

	resp, err := s.eventCounterClient.GetExperimentGoalCount(
		ctx, &eventcounter.GetExperimentGoalCountRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			GoalId:         req.GoalId,
			StartAt:        req.StartAt,
			EndAt:          req.EndAt,
			FeatureId:      req.FeatureId,
			FeatureVersion: req.FeatureVersion,
			VariationIds:   req.VariationIds,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get experiment goal count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("experimentId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
				zap.Int64("startAt", req.StartAt),
				zap.Int64("endAt", req.EndAt),
				zap.String("goalId", req.GoalId),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ExperimentGoalCount response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("experimentId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
				zap.Int64("startAt", req.StartAt),
				zap.Int64("endAt", req.EndAt),
				zap.String("goalId", req.GoalId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetExperimentGoalCountResponse{
		GoalId:          resp.GoalId,
		VariationCounts: resp.VariationCounts,
	}, nil
}

func (s *grpcGatewayService) GetOpsEvaluationUserCount(
	ctx context.Context,
	req *gwproto.GetOpsEvaluationUserCountRequest,
) (*gwproto.GetOpsEvaluationUserCountResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetOpsEvaluationUserCount request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetOpsEvaluationUserCount, "",
	).Inc()

	resp, err := s.eventCounterClient.GetOpsEvaluationUserCount(
		ctx, &eventcounter.GetOpsEvaluationUserCountRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			FeatureId:      req.FeatureId,
			FeatureVersion: req.FeatureVersion,
			ClauseId:       req.ClauseId,
			OpsRuleId:      req.OpsRuleId,
			VariationId:    req.VariationId,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get ops evaluation user count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("OpsEvaluationUserCount response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetOpsEvaluationUserCountResponse{
		OpsRuleId: resp.OpsRuleId,
		ClauseId:  resp.ClauseId,
		Count:     resp.Count,
	}, nil
}

func (s *grpcGatewayService) GetOpsGoalUserCount(
	ctx context.Context,
	req *gwproto.GetOpsGoalUserCountRequest,
) (*gwproto.GetOpsGoalUserCountResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetOpsGoalUserCount request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetOpsGoalUserCount, "",
	).Inc()

	resp, err := s.eventCounterClient.GetOpsGoalUserCount(
		ctx, &eventcounter.GetOpsGoalUserCountRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			OpsRuleId:      req.OpsRuleId,
			ClauseId:       req.ClauseId,
			FeatureId:      req.FeatureId,
			FeatureVersion: req.FeatureVersion,
			VariationId:    req.VariationId,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get ops goal user count",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("OpsGoalUserCount response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", req.FeatureId),
				zap.Int32("featureVersion", req.FeatureVersion),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetOpsGoalUserCountResponse{
		OpsRuleId: resp.OpsRuleId,
		ClauseId:  resp.ClauseId,
		Count:     resp.Count,
	}, nil
}
