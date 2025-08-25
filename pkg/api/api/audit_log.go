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

	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	auditlogproto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) GetAuditLog(
	ctx context.Context,
	request *gwproto.GetAuditLogRequest,
) (*gwproto.GetAuditLogResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get auditlog request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetAuditLog, "").Inc()

	res, err := s.auditLogClient.GetAuditLog(
		ctx,
		&auditlogproto.GetAuditLogRequest{
			Id:            request.Id,
			EnvironmentId: envAPIKey.Environment.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to get audit log",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", request.Id),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Get audit log response is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("id", request.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetAuditLogResponse{
		AuditLog: res.AuditLog,
	}, nil
}

func (s *grpcGatewayService) ListAuditLogs(
	ctx context.Context,
	request *gwproto.ListAuditLogsRequest,
) (*gwproto.ListAuditLogsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListAuditLogs request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListAuditLogs, "").Inc()

	res, err := s.auditLogClient.ListAuditLogs(
		ctx,
		&auditlogproto.ListAuditLogsRequest{
			PageSize:       request.PageSize,
			Cursor:         request.Cursor,
			OrderBy:        request.OrderBy,
			OrderDirection: request.OrderDirection,
			SearchKeyword:  request.SearchKeyword,
			From:           request.From,
			To:             request.To,
			EntityType:     request.EntityType,
			EnvironmentId:  envAPIKey.Environment.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to list audit logs",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("List audit logs response is nil")
		return nil, ErrInternal
	}
	return &gwproto.ListAuditLogsResponse{
		AuditLogs:  res.AuditLogs,
		Cursor:     res.Cursor,
		TotalCount: res.TotalCount,
	}, nil
}

func (s *grpcGatewayService) ListFeatureHistory(
	ctx context.Context,
	request *gwproto.ListFeatureHistoryRequest,
) (*gwproto.ListFeatureHistoryResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListFeatureHistory request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListFeatureHistory, "").Inc()

	res, err := s.auditLogClient.ListFeatureHistory(
		ctx,
		&auditlogproto.ListFeatureHistoryRequest{
			FeatureId:      request.FeatureId,
			PageSize:       request.PageSize,
			Cursor:         request.Cursor,
			OrderBy:        request.OrderBy,
			OrderDirection: request.OrderDirection,
			SearchKeyword:  request.SearchKeyword,
			From:           request.From,
			To:             request.To,
			EnvironmentId:  envAPIKey.Environment.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to list feature history",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("List feature history response is nil")
		return nil, ErrInternal
	}
	return &gwproto.ListFeatureHistoryResponse{
		AuditLogs:  res.AuditLogs,
		Cursor:     res.Cursor,
		TotalCount: res.TotalCount,
	}, nil
}
