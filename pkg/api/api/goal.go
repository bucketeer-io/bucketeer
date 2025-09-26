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

func (s *grpcGatewayService) GetGoal(
	ctx context.Context,
	req *gwproto.GetGoalRequest,
) (*gwproto.GetGoalResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetGoal request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("goalId", req.Id),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetGoal, "").Inc()
	resp, err := s.experimentClient.GetGoal(ctx, &experimentproto.GetGoalRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to get goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("goalId", req.Id),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("Get goal response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("goalId", req.Id),
			)...,
		)
		return nil, ErrInternal
	}

	return &gwproto.GetGoalResponse{
		Goal: resp.Goal,
	}, nil
}

func (s *grpcGatewayService) ListGoals(
	ctx context.Context,
	req *gwproto.ListGoalsRequest,
) (*gwproto.ListGoalsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListGoals request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListGoals, "").Inc()
	resp, err := s.experimentClient.ListGoals(ctx, &experimentproto.ListGoalsRequest{
		EnvironmentId:  envAPIKey.Environment.Id,
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		SearchKeyword:  req.SearchKeyword,
		IsInUseStatus:  req.IsInUseStatus,
		Archived:       req.Archived,
		ConnectionType: req.ConnectionType,
	})
	if err != nil {
		s.logger.Error("Failed to list goals",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("List goals response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields()...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListGoalsResponse{
		Goals:      resp.Goals,
		Cursor:     resp.Cursor,
		TotalCount: resp.TotalCount,
	}, nil
}

func (s *grpcGatewayService) CreateGoal(
	ctx context.Context,
	req *gwproto.CreateGoalRequest,
) (*gwproto.CreateGoalResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateGoal request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateGoal, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	resp, err := s.experimentClient.CreateGoal(ctx, &experimentproto.CreateGoalRequest{
		EnvironmentId:  envAPIKey.Environment.Id,
		Id:             req.Id,
		Name:           req.Name,
		Description:    req.Description,
		ConnectionType: req.ConnectionType,
	})
	if err != nil {
		s.logger.Error("Failed to create goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("Create goal response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields()...,
		)
		return nil, ErrInternal
	}
	return &gwproto.CreateGoalResponse{
		Goal: resp.Goal,
	}, nil
}

func (s *grpcGatewayService) DeleteGoal(
	ctx context.Context,
	req *gwproto.DeleteGoalRequest,
) (*gwproto.DeleteGoalResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteGoal request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodDeleteGoal, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.experimentClient.DeleteGoal(ctx, &experimentproto.DeleteGoalRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to delete goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	return &gwproto.DeleteGoalResponse{}, nil
}

func (s *grpcGatewayService) UpdateGoal(
	ctx context.Context,
	req *gwproto.UpdateGoalRequest,
) (*gwproto.UpdateGoalResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check UpdateGoal request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdateGoal, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.experimentClient.UpdateGoal(ctx, &experimentproto.UpdateGoalRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
		Name:          req.Name,
		Description:   req.Description,
		Archived:      req.Archived,
	})
	if err != nil {
		s.logger.Error("Failed to update goal",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	return &gwproto.UpdateGoalResponse{}, nil
}
