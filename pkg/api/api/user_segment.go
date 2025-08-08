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
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) CreateSegment(
	ctx context.Context,
	req *gwproto.CreateSegmentRequest,
) (*gwproto.CreateSegmentResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateSegment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", req.Name),
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
	res, err := s.featureClient.CreateSegment(ctx, &featureproto.CreateSegmentRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Name:          req.Name,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Not found created segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environment_id", envAPIKey.Environment.Id),
				zap.String("name", req.Name),
			)...)
		return nil, errInternal
	}
	return &gwproto.CreateSegmentResponse{
		Segment: res.Segment,
	}, nil
}

func (s *grpcGatewayService) GetSegment(
	ctx context.Context,
	req *gwproto.GetSegmentRequest,
) (*gwproto.GetSegmentResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetSegment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("segmentId", req.Id),
			)...,
		)
		return nil, err
	}

	res, err := s.featureClient.GetSegment(
		ctx,
		&featureproto.GetSegmentRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Segment not found",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("segmentId", req.Id),
			)...,
		)
		return nil, errInternal
	}

	return &gwproto.GetSegmentResponse{
		Segment: res.Segment,
	}, nil
}

func (s *grpcGatewayService) ListSegments(
	ctx context.Context,
	req *gwproto.ListSegmentsRequest,
) (*gwproto.ListSegmentsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListSegments request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	res, err := s.featureClient.ListSegments(
		ctx,
		&featureproto.ListSegmentsRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			PageSize:       req.PageSize,
			Cursor:         req.Cursor,
			OrderBy:        req.OrderBy,
			OrderDirection: req.OrderDirection,
			SearchKeyword:  req.SearchKeyword,
			Status:         req.Status,
			IsInUseStatus:  req.IsInUseStatus,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to list segments: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environment_id", envAPIKey.Environment.Id),
				zap.String("search_keyword", req.SearchKeyword),
			)...,
		)
		return nil, errInternal
	}
	return &gwproto.ListSegmentsResponse{
		Segments:   res.Segments,
		Cursor:     res.Cursor,
		TotalCount: res.TotalCount,
	}, nil
}

func (s *grpcGatewayService) DeleteSegment(
	ctx context.Context,
	req *gwproto.DeleteSegmentRequest,
) (*gwproto.DeleteSegmentResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteSegment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("segmentId", req.Id),
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
	_, err = s.featureClient.DeleteSegment(
		ctx,
		&featureproto.DeleteSegmentRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		return nil, err
	}
	return &gwproto.DeleteSegmentResponse{}, nil
}

func (s *grpcGatewayService) UpdateSegment(
	ctx context.Context,
	req *gwproto.UpdateSegmentRequest,
) (*gwproto.UpdateSegmentResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check UpdateSegment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("segmentId", req.Id),
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
	res, err := s.featureClient.UpdateSegment(
		ctx,
		&featureproto.UpdateSegmentRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
			Name:          req.Name,
			Description:   req.Description,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Not found updated segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("segmentId", req.Id),
			)...,
		)
		return nil, errInternal
	}
	return &gwproto.UpdateSegmentResponse{
		Segment: res.Segment,
	}, nil
}

func (s *grpcGatewayService) BulkUploadSegmentUsers(
	ctx context.Context,
	req *gwproto.BulkUploadSegmentUsersRequest,
) (*gwproto.BulkUploadSegmentUsersResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check BulkUploadSegmentUsers request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("segmentId", req.SegmentId),
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
	_, err = s.featureClient.BulkUploadSegmentUsers(
		ctx,
		&featureproto.BulkUploadSegmentUsersRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			SegmentId:     req.SegmentId,
			Data:          req.Data,
			State:         req.State,
		},
	)
	if err != nil {
		return nil, err
	}
	return &gwproto.BulkUploadSegmentUsersResponse{}, nil
}
