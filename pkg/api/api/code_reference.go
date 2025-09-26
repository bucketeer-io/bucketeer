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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	coderefproto "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func (s *grpcGatewayService) GetCodeReference(
	ctx context.Context,
	req *gatewayproto.GetCodeReferenceRequest,
) (*gatewayproto.GetCodeReferenceResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetCodeReference request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetCodeReference, "").Inc()

	resp, err := s.codeRefClient.GetCodeReference(ctx, &coderefproto.GetCodeReferenceRequest{
		Id:            req.Id,
		EnvironmentId: envAPIKey.Environment.Id,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			s.logger.Info(
				"Code reference not found",
				zap.String("id", req.Id),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, "not found")
		}
		s.logger.Error(
			"Failed to get code reference",
			zap.String("id", req.Id),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gatewayproto.GetCodeReferenceResponse{
		CodeReference: resp.CodeReference,
	}, nil
}

func (s *grpcGatewayService) ListCodeReferences(
	ctx context.Context,
	req *gatewayproto.ListCodeReferencesRequest,
) (*gatewayproto.ListCodeReferencesResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListCodeReferences request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListCodeReferences, "").Inc()

	resp, err := s.codeRefClient.ListCodeReferences(ctx, &coderefproto.ListCodeReferencesRequest{
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		FeatureId:      req.FeatureId,
		EnvironmentId:  envAPIKey.Environment.Id,
	})
	if err != nil {
		s.logger.Error(
			"Failed to list code references",
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gatewayproto.ListCodeReferencesResponse{
		CodeReferences: resp.CodeReferences,
		Cursor:         resp.Cursor,
		TotalCount:     resp.TotalCount,
	}, nil
}

func (s *grpcGatewayService) CreateCodeReference(
	ctx context.Context,
	req *gatewayproto.CreateCodeReferenceRequest,
) (*gatewayproto.CreateCodeReferenceResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateCodeReference request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateCodeReference, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)

	resp, err := s.codeRefClient.CreateCodeReference(ctx, &coderefproto.CreateCodeReferenceRequest{
		FeatureId:        req.FeatureId,
		EnvironmentId:    envAPIKey.Environment.Id,
		FilePath:         req.FilePath,
		FileExtension:    req.FileExtension,
		LineNumber:       req.LineNumber,
		CodeSnippet:      req.CodeSnippet,
		ContentHash:      req.ContentHash,
		Aliases:          req.Aliases,
		RepositoryName:   req.RepositoryName,
		RepositoryOwner:  req.RepositoryOwner,
		RepositoryType:   req.RepositoryType,
		RepositoryBranch: req.RepositoryBranch,
		CommitHash:       req.CommitHash,
	})
	if err != nil {
		s.logger.Error(
			"Failed to create code reference",
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gatewayproto.CreateCodeReferenceResponse{
		CodeReference: resp.CodeReference,
	}, nil
}

func (s *grpcGatewayService) UpdateCodeReference(
	ctx context.Context,
	req *gatewayproto.UpdateCodeReferenceRequest,
) (*gatewayproto.UpdateCodeReferenceResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check UpdateCodeReference request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdateCodeReference, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)

	resp, err := s.codeRefClient.UpdateCodeReference(ctx, &coderefproto.UpdateCodeReferenceRequest{
		Id:               req.Id,
		EnvironmentId:    envAPIKey.Environment.Id,
		FilePath:         req.FilePath,
		FileExtension:    req.FileExtension,
		LineNumber:       req.LineNumber,
		CodeSnippet:      req.CodeSnippet,
		ContentHash:      req.ContentHash,
		Aliases:          req.Aliases,
		RepositoryBranch: req.RepositoryBranch,
		CommitHash:       req.CommitHash,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			s.logger.Info(
				"Code reference not found",
				zap.String("id", req.Id),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, "not found")
		}
		s.logger.Error(
			"Failed to update code reference",
			zap.String("id", req.Id),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gatewayproto.UpdateCodeReferenceResponse{
		CodeReference: resp.CodeReference,
	}, nil
}

func (s *grpcGatewayService) DeleteCodeReference(
	ctx context.Context,
	req *gatewayproto.DeleteCodeReferenceRequest,
) (*gatewayproto.DeleteCodeReferenceResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteCodeReference request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodDeleteCodeReference, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)

	_, err = s.codeRefClient.DeleteCodeReference(ctx, &coderefproto.DeleteCodeReferenceRequest{
		Id:            req.Id,
		EnvironmentId: envAPIKey.Environment.Id,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			s.logger.Info(
				"Code reference not found",
				zap.String("id", req.Id),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, "not found")
		}
		s.logger.Error(
			"Failed to delete code reference",
			zap.String("id", req.Id),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gatewayproto.DeleteCodeReferenceResponse{}, nil
}
