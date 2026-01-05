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

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func (s *grpcGatewayService) GetEnvironmentV2(
	ctx context.Context,
	req *gwproto.GetEnvironmentV2Request,
) (*gwproto.GetEnvironmentV2Response, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetEnvironmentV2 request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetEnvironmentV2, "").Inc()

	resp, err := s.environmentClient.GetEnvironmentV2(ctx, &environmentproto.GetEnvironmentV2Request{
		Id: req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to get environment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.Id),
			)...,
		)
		return nil, err
	}

	if resp == nil {
		s.logger.Error("Environment resp is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environmentId", req.Id),
			)...,
		)
		return nil, ErrInternal
	}

	// for security, ensure the environment belongs to the same organization as the API key
	if resp.Environment.OrganizationId != envAPIKey.Environment.OrganizationId {
		s.logger.Error("Environment does not belong to the same organization as the API key",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environmentId", req.Id),
				zap.String("OrganizationId", resp.Environment.OrganizationId),
				zap.String("apiKeyOrgId", envAPIKey.Environment.OrganizationId),
			)...,
		)
		return nil, ErrNotFound
	}

	return &gwproto.GetEnvironmentV2Response{
		Environment: resp.Environment,
	}, nil
}

func (s *grpcGatewayService) ListEnvironmentsV2(
	ctx context.Context,
	req *gwproto.ListEnvironmentsV2Request,
) (*gwproto.ListEnvironmentsV2Response, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListEnvironmentV2 request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListEnvironmentsV2, "").Inc()

	resp, err := s.environmentClient.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		ProjectId:      req.ProjectId,
		Archived:       req.Archived,
		SearchKeyword:  req.SearchKeyword,
		OrganizationId: envAPIKey.Environment.OrganizationId, // restrict to the same organization as the API key
	})
	if err != nil {
		s.logger.Error("Failed to list environments",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("Environments resp is nil",
			log.FieldsFromIncomingContext(ctx)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListEnvironmentsV2Response{
		Environments: resp.Environments,
		TotalCount:   resp.TotalCount,
		Cursor:       resp.Cursor,
	}, nil
}

func (s *grpcGatewayService) GetProject(
	ctx context.Context,
	req *gwproto.GetProjectRequest,
) (*gwproto.GetProjectResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetProject request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetProject, "",
	).Inc()

	resp, err := s.environmentClient.GetProject(ctx, &environmentproto.GetProjectRequest{
		Id: req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to get project",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("projectId", req.Id),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("Project resp is nil",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("projectId", req.Id),
			)...,
		)
		return nil, ErrInternal
	}

	// for security, ensure the project belongs to the same organization as the API key
	if resp.Project.OrganizationId != envAPIKey.Environment.OrganizationId {
		s.logger.Error("Project does not belong to the same organization as the API key",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("projectId", req.Id),
				zap.String("OrganizationId", resp.Project.OrganizationId),
				zap.String("apiKeyOrgId", envAPIKey.Environment.OrganizationId),
			)...,
		)
		return nil, ErrNotFound
	}

	return &gwproto.GetProjectResponse{
		Project: resp.Project,
	}, nil
}

func (s *grpcGatewayService) ListProjects(
	ctx context.Context,
	req *gwproto.ListProjectsRequest,
) (*gwproto.ListProjectsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListProjects request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListProjects, "",
	).Inc()

	resp, err := s.environmentClient.ListProjects(ctx, &environmentproto.ListProjectsRequest{
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		SearchKeyword:  req.SearchKeyword,
		Disabled:       req.Disabled,
		// restrict to the same organization as the API key
		OrganizationIds: []string{envAPIKey.Environment.OrganizationId},
	})
	if err != nil {
		s.logger.Error("Failed to list projects",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if resp == nil {
		s.logger.Error("Projects resp is nil",
			log.FieldsFromIncomingContext(ctx)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListProjectsResponse{
		Projects:   resp.Projects,
		TotalCount: resp.TotalCount,
		Cursor:     resp.Cursor,
	}, nil
}
