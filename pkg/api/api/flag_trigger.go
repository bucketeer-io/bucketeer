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

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func (s *grpcGatewayService) CreateFlagTrigger(
	ctx context.Context,
	req *gwproto.CreateFlagTriggerRequest,
) (*gwproto.CreateFlagTriggerResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check create flag trigger request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureID", req.FeatureId),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateFlagTrigger, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	resp, err := s.featureClient.CreateFlagTrigger(ctx, &featureproto.CreateFlagTriggerRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		FeatureId:     req.FeatureId,
		Type:          req.Type,
		Action:        req.Action,
		Description:   req.Description,
	})
	if err != nil {
		s.logger.Error("Failed to create flag trigger",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureID", req.FeatureId),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("CreateFlagTrigger returned nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureID", req.FeatureId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.CreateFlagTriggerResponse{
		FlagTrigger: resp.FlagTrigger,
	}, nil
}

func (s *grpcGatewayService) DeleteFlagTrigger(
	ctx context.Context,
	req *gwproto.DeleteFlagTriggerRequest,
) (*gwproto.DeleteFlagTriggerResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check delete flag trigger request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodDeleteFlagTrigger, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.featureClient.DeleteFlagTrigger(ctx, &featureproto.DeleteFlagTriggerRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to delete flag trigger",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", req.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.DeleteFlagTriggerResponse{}, nil
}

func (s *grpcGatewayService) UpdateFlagTrigger(
	ctx context.Context,
	req *gwproto.UpdateFlagTriggerRequest,
) (*gwproto.UpdateFlagTriggerResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check update flag trigger request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdateFlagTrigger, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	resp, err := s.featureClient.UpdateFlagTrigger(ctx, &featureproto.UpdateFlagTriggerRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
		Description:   req.Description,
		Reset_:        req.Reset_,
		Disabled:      req.Disabled,
	})
	if err != nil {
		s.logger.Error("Failed to update flag trigger",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", req.Id),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("UpdateFlagTrigger returned nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("ID", req.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.UpdateFlagTriggerResponse{
		Url: resp.Url,
	}, nil
}

func (s *grpcGatewayService) GetFlagTrigger(
	ctx context.Context,
	req *gwproto.GetFlagTriggerRequest,
) (*gwproto.GetFlagTriggerResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get flag trigger request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetFlagTrigger, "").Inc()

	resp, err := s.featureClient.GetFlagTrigger(ctx, &featureproto.GetFlagTriggerRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to get flag trigger",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", req.Id),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("GetFlagTrigger returned nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("ID", req.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetFlagTriggerResponse{
		FlagTrigger: resp.FlagTrigger,
	}, nil
}

func (s *grpcGatewayService) ListFlagTriggers(
	ctx context.Context,
	req *gwproto.ListFlagTriggersRequest,
) (*gwproto.ListFlagTriggersResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check list flag triggers request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListFlagTriggers, "").Inc()

	resp, err := s.featureClient.ListFlagTriggers(ctx, &featureproto.ListFlagTriggersRequest{
		EnvironmentId:  envAPIKey.Environment.Id,
		FeatureId:      req.FeatureId,
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
	})
	if err != nil {
		s.logger.Error("Failed to list flag triggers",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("ListFlagTriggers returned nil response",
			log.FieldsFromIncomingContext(ctx).AddFields()...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListFlagTriggersResponse{
		FlagTriggers: resp.FlagTriggers,
		Cursor:       resp.Cursor,
		TotalCount:   resp.TotalCount,
	}, nil
}
