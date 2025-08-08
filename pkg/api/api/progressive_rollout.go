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

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	"github.com/bucketeer-io/bucketeer/proto/autoops"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) CreateProgressiveRollout(
	ctx context.Context,
	req *gwproto.CreateProgressiveRolloutRequest,
) (*gwproto.CreateProgressiveRolloutResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateProgressiveRollout request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, err
	}
	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	resp, err := s.autoOpsClient.CreateProgressiveRollout(
		ctx,
		&autoops.CreateProgressiveRolloutRequest{
			EnvironmentId:                            envAPIKey.Environment.Id,
			FeatureId:                                req.FeatureId,
			ProgressiveRolloutManualScheduleClause:   req.ProgressiveRolloutManualScheduleClause,
			ProgressiveRolloutTemplateScheduleClause: req.ProgressiveRolloutTemplateScheduleClause,
		},
	)
	if err != nil {
		s.logger.Error("Failed to create progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("CreateProgressiveRolloutResponse is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", req.FeatureId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.CreateProgressiveRolloutResponse{
		ProgressiveRollout: resp.ProgressiveRollout,
	}, nil
}

func (s *grpcGatewayService) GetProgressiveRollout(
	ctx context.Context,
	req *gwproto.GetProgressiveRolloutRequest,
) (*gwproto.GetProgressiveRolloutResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetProgressiveRollout request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	resp, err := s.autoOpsClient.GetProgressiveRollout(
		ctx,
		&autoops.GetProgressiveRolloutRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("GetProgressiveRolloutResponse is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("id", req.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetProgressiveRolloutResponse{
		ProgressiveRollout: resp.ProgressiveRollout,
	}, nil
}

func (s *grpcGatewayService) StopProgressiveRollout(
	ctx context.Context,
	req *gwproto.StopProgressiveRolloutRequest,
) (*gwproto.StopProgressiveRolloutResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check StopProgressiveRollout request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.autoOpsClient.StopProgressiveRollout(
		ctx,
		&autoops.StopProgressiveRolloutRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to stop progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.StopProgressiveRolloutResponse{}, nil
}

func (s *grpcGatewayService) DeleteProgressiveRollout(
	ctx context.Context,
	req *gwproto.DeleteProgressiveRolloutRequest,
) (*gwproto.DeleteProgressiveRolloutResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteProgressiveRollout request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.autoOpsClient.DeleteProgressiveRollout(
		ctx,
		&autoops.DeleteProgressiveRolloutRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to delete progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.DeleteProgressiveRolloutResponse{}, nil
}

func (s *grpcGatewayService) ListProgressiveRollouts(
	ctx context.Context,
	req *gwproto.ListProgressiveRolloutsRequest,
) (*gwproto.ListProgressiveRolloutsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListProgressiveRollouts request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	resp, err := s.autoOpsClient.ListProgressiveRollouts(
		ctx,
		&autoops.ListProgressiveRolloutsRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			PageSize:       req.PageSize,
			Cursor:         req.Cursor,
			OrderBy:        req.OrderBy,
			OrderDirection: req.OrderDirection,
			FeatureIds:     req.FeatureIds,
			Status:         req.Status,
			Type:           req.Type,
		},
	)
	if err != nil {
		s.logger.Error("Failed to list progressive rollouts",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ListProgressiveRolloutsResponse is nil")
		return nil, ErrInternal
	}
	return &gwproto.ListProgressiveRolloutsResponse{
		ProgressiveRollouts: resp.ProgressiveRollouts,
		Cursor:              resp.Cursor,
		TotalCount:          resp.TotalCount,
	}, nil
}

func (s *grpcGatewayService) ExecuteProgressiveRollout(
	ctx context.Context,
	req *gwproto.ExecuteProgressiveRolloutRequest,
) (*gwproto.ExecuteProgressiveRolloutResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ExecuteProgressiveRollout request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	resp, err := s.autoOpsClient.ExecuteProgressiveRollout(
		ctx,
		&autoops.ExecuteProgressiveRolloutRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to execute progressive rollout",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ExecuteProgressiveRolloutResponse is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("id", req.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ExecuteProgressiveRolloutResponse{}, nil
}
