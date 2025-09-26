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
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	tagproto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

func (s *grpcGatewayService) CreateTag(
	ctx context.Context,
	req *gwproto.CreateTagRequest,
) (*gwproto.CreateTagResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreatePush request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", req.Name),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateTag, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	res, err := s.tagClient.CreateTag(
		ctx,
		&tagproto.CreateTagRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Name:          req.Name,
			EntityType:    tagproto.Tag_FEATURE_FLAG, // deprecate type account
		},
	)
	if err != nil {
		s.logger.Error("Failed to create tag",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", req.Name),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to create tag: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("name", req.Name),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.CreateTagResponse{
		Tag: res.Tag,
	}, nil
}

func (s *grpcGatewayService) DeleteTag(
	ctx context.Context,
	req *gwproto.DeleteTagRequest,
) (*gwproto.DeleteTagResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteTag request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodDeleteTag, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.tagClient.DeleteTag(
		ctx,
		&tagproto.DeleteTagRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to delete tag",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.DeleteTagResponse{}, nil
}

func (s *grpcGatewayService) ListTags(
	ctx context.Context,
	req *gwproto.ListTagsRequest,
) (*gwproto.ListTagsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListTags request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListTags, "").Inc()

	res, err := s.tagClient.ListTags(
		ctx,
		&tagproto.ListTagsRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			PageSize:       req.PageSize,
			Cursor:         req.Cursor,
			OrderBy:        req.OrderBy,
			OrderDirection: req.OrderDirection,
			SearchKeyword:  req.SearchKeyword,
			OrganizationId: req.OrganizationId,
		},
	)
	if err != nil {
		s.logger.Error("Failed to list tags",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("search_keyword", req.SearchKeyword),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to list tags: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environment_id", envAPIKey.Environment.Id),
				zap.String("search_keyword", req.SearchKeyword),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListTagsResponse{
		Tags:       res.Tags,
		Cursor:     res.Cursor,
		TotalCount: res.TotalCount,
	}, nil
}
