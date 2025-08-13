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
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) GetAutoOpsRule(
	ctx context.Context,
	request *gwproto.GetAutoOpsRuleRequest,
) (*gwproto.GetAutoOpsRuleResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetAutoOpsRule request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetAutoOpsRule, "").Inc()

	res, err := s.autoOpsClient.GetAutoOpsRule(
		ctx,
		&autoopsproto.GetAutoOpsRuleRequest{
			Id:            request.Id,
			EnvironmentId: envAPIKey.Environment.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get auto ops rule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Get auto ops rule response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("id", request.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetAutoOpsRuleResponse{
		AutoOpsRule: res.AutoOpsRule,
	}, nil
}

func (s *grpcGatewayService) CreateAutoOpsRule(
	ctx context.Context,
	request *gwproto.CreateAutoOpsRuleRequest,
) (*gwproto.CreateAutoOpsRuleResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateAutoOpsRule request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", request.FeatureId),
				zap.String("ops type", request.OpsType.String()),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateAutoOpsRule, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	res, err := s.autoOpsClient.CreateAutoOpsRule(
		ctx,
		&autoopsproto.CreateAutoOpsRuleRequest{
			EnvironmentId:       envAPIKey.Environment.Id,
			FeatureId:           request.FeatureId,
			OpsType:             request.OpsType,
			OpsEventRateClauses: request.OpsEventRateClauses,
			DatetimeClauses:     request.DatetimeClauses,
		},
	)
	if err != nil {
		s.logger.Error("Failed to create auto ops rule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", request.FeatureId),
				zap.String("ops type", request.OpsType.String()),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Create auto ops rule response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("featureId", request.FeatureId),
				zap.String("ops type", request.OpsType.String()),
			)...,
		)
		return nil, ErrInternal
	}

	return &gwproto.CreateAutoOpsRuleResponse{
		AutoOpsRule: res.AutoOpsRule,
	}, nil
}

func (s *grpcGatewayService) ListAutoOpsRules(
	ctx context.Context,
	request *gwproto.ListAutoOpsRulesRequest,
) (*gwproto.ListAutoOpsRulesResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListAutoOpsRules request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Strings("featureIds", request.FeatureIds),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListAutoOpsRules, "").Inc()

	res, err := s.autoOpsClient.ListAutoOpsRules(
		ctx,
		&autoopsproto.ListAutoOpsRulesRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			PageSize:      request.PageSize,
			Cursor:        request.Cursor,
			FeatureIds:    request.FeatureIds,
		},
	)
	if err != nil {
		s.logger.Error("Failed to list auto ops rules",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Strings("featureIds", request.FeatureIds),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("List auto ops rules response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Strings("featureIds", request.FeatureIds),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListAutoOpsRulesResponse{
		AutoOpsRules: res.AutoOpsRules,
		Cursor:       res.Cursor,
	}, nil
}

func (s *grpcGatewayService) StopAutoOpsRule(
	ctx context.Context,
	request *gwproto.StopAutoOpsRuleRequest,
) (*gwproto.StopAutoOpsRuleResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check StopAutoOpsRule request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodStopAutoOpsRule, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.autoOpsClient.StopAutoOpsRule(
		ctx,
		&autoopsproto.StopAutoOpsRuleRequest{
			Id:            request.Id,
			EnvironmentId: envAPIKey.Environment.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to stop auto ops rule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.StopAutoOpsRuleResponse{}, nil
}

func (s *grpcGatewayService) DeleteAutoOpsRule(
	ctx context.Context,
	request *gwproto.DeleteAutoOpsRuleRequest,
) (*gwproto.DeleteAutoOpsRuleResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteAutoOpsRule request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodDeleteAutoOpsRule, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.autoOpsClient.DeleteAutoOpsRule(
		ctx,
		&autoopsproto.DeleteAutoOpsRuleRequest{
			Id:            request.Id,
			EnvironmentId: envAPIKey.Environment.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to delete auto ops rule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.DeleteAutoOpsRuleResponse{}, nil
}

func (s *grpcGatewayService) UpdateAutoOpsRule(
	ctx context.Context,
	request *gwproto.UpdateAutoOpsRuleRequest,
) (*gwproto.UpdateAutoOpsRuleResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check UpdateAutoOpsRule request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdateAutoOpsRule, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.autoOpsClient.UpdateAutoOpsRule(
		ctx,
		&autoopsproto.UpdateAutoOpsRuleRequest{
			Id:                        request.Id,
			EnvironmentId:             envAPIKey.Environment.Id,
			OpsEventRateClauseChanges: request.OpsEventRateClauseChanges,
			DatetimeClauseChanges:     request.DatetimeClauseChanges,
		},
	)
	if err != nil {
		s.logger.Error("Failed to update auto ops rule",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.UpdateAutoOpsRuleResponse{}, nil
}

func (s *grpcGatewayService) ExecuteAutoOps(
	ctx context.Context,
	request *gwproto.ExecuteAutoOpsRequest,
) (*gwproto.ExecuteAutoOpsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ExecuteAutoOps request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodExecuteAutoOps, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	res, err := s.autoOpsClient.ExecuteAutoOps(
		ctx,
		&autoopsproto.ExecuteAutoOpsRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            request.Id,
			ClauseId:      request.ClauseId,
		},
	)
	if err != nil {
		s.logger.Error("Failed to execute auto ops",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("ID", request.Id),
				zap.String("ClauseID", request.ClauseId),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Execute auto ops response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("ID", request.Id),
				zap.String("ClauseID", request.ClauseId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ExecuteAutoOpsResponse{
		AlreadyTriggered: res.AlreadyTriggered,
	}, nil
}
