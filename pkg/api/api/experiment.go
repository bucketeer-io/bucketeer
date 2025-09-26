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
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func (s *grpcGatewayService) GetExperiment(
	ctx context.Context,
	req *gwproto.GetExperimentRequest,
) (*gwproto.GetExperimentResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetExperiment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetExperiment, "").Inc()
	resp, err := s.experimentClient.GetExperiment(ctx, &experimentproto.GetExperimentRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to get experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("experimentId", req.Id),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("Experiment resp is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("experimentId", req.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetExperimentResponse{
		Experiment: resp.Experiment,
	}, nil
}

func (s *grpcGatewayService) ListExperiments(
	ctx context.Context,
	req *gwproto.ListExperimentsRequest,
) (*gwproto.ListExperimentsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListExperiments request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListExperiments, "",
	).Inc()
	resp, err := s.experimentClient.ListExperiments(ctx, &experimentproto.ListExperimentsRequest{
		EnvironmentId:  envAPIKey.Environment.Id,
		FeatureVersion: req.FeatureVersion,
		StartAt:        req.StartAt,
		StopAt:         req.StopAt,
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		Maintainer:     req.Maintainer,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		SearchKeyword:  req.SearchKeyword,
		Archived:       req.Archived,
		Statuses:       req.Statuses,
	})
	if err != nil {
		s.logger.Error("Failed to list experiments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ListExperiments resp is nil",
			log.FieldsFromIncomingContext(ctx)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListExperimentsResponse{
		Experiments: resp.Experiments,
		Cursor:      resp.Cursor,
		TotalCount:  resp.TotalCount,
	}, nil
}

func (s *grpcGatewayService) CreateExperiment(
	ctx context.Context,
	req *gwproto.CreateExperimentRequest,
) (*gwproto.CreateExperimentResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateExperiment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateExperiment, "",
	).Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	resp, err := s.experimentClient.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{
		EnvironmentId:   envAPIKey.Environment.Id,
		Name:            req.Name,
		Description:     req.Description,
		FeatureId:       req.FeatureId,
		GoalIds:         req.GoalIds,
		StartAt:         req.StartAt,
		StopAt:          req.StopAt,
		BaseVariationId: req.BaseVariationId,
	})
	if err != nil {
		s.logger.Error("Failed to create experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("CreateExperiment resp is nil",
			log.FieldsFromIncomingContext(ctx)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.CreateExperimentResponse{
		Experiment: resp.Experiment,
	}, nil
}

func (s *grpcGatewayService) UpdateExperiment(
	ctx context.Context,
	req *gwproto.UpdateExperimentRequest,
) (*gwproto.UpdateExperimentResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check UpdateExperiment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdateExperiment, "",
	).Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.experimentClient.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
		Name:          req.Name,
		Description:   req.Description,
		StartAt:       req.StartAt,
		StopAt:        req.StopAt,
		Archived:      req.Archived,
		Status:        req.Status,
	})
	if err != nil {
		s.logger.Error("Failed to update experiment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	return &gwproto.UpdateExperimentResponse{}, nil
}
