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
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) CreateFeature(
	ctx context.Context,
	req *gwproto.CreateFeatureRequest,
) (*gwproto.CreateFeatureResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateFeature request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.Id),
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

	res, err := s.featureClient.CreateFeature(ctx, &featureproto.CreateFeatureRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Command: &featureproto.CreateFeatureCommand{
			Id:                       req.Id,
			Name:                     req.Name,
			Description:              req.Description,
			Variations:               req.Variations,
			Tags:                     req.Tags,
			DefaultOnVariationIndex:  &wrapperspb.Int32Value{Value: req.OnVariationIndex},
			DefaultOffVariationIndex: &wrapperspb.Int32Value{Value: req.OffVariationIndex},
			VariationType:            req.VariationType,
		},
	})
	if err != nil {
		return nil, err
	}
	return &gwproto.CreateFeatureResponse{Feature: res.Feature}, nil
}

func (s *grpcGatewayService) GetFeature(
	ctx context.Context,
	req *gwproto.GetFeatureRequest,
) (*gwproto.GetFeatureResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetFeature request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.Id),
			)...,
		)
		return nil, err
	}
	resp, err := s.featureClient.GetFeature(ctx, &featureproto.GetFeatureRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &gwproto.GetFeatureResponse{
		Feature: resp.Feature,
	}, nil
}

func (s *grpcGatewayService) ListFeatures(
	ctx context.Context,
	req *gwproto.ListFeaturesRequest,
) (*gwproto.ListFeaturesResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListFeatures request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	resp, err := s.featureClient.ListFeatures(ctx, &featureproto.ListFeaturesRequest{
		EnvironmentId:  envAPIKey.Environment.Id,
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
	})
	if err != nil {
		return nil, err
	}
	return &gwproto.ListFeaturesResponse{
		Features: resp.Features,
	}, nil
}

func (s *grpcGatewayService) UpdateFeature(
	ctx context.Context,
	req *gwproto.UpdateFeatureRequest,
) (*gwproto.UpdateFeatureResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetFeature request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("featureId", req.Id),
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

	res, err := s.featureClient.UpdateFeature(ctx, &featureproto.UpdateFeatureRequest{
		Comment:         req.Comment,
		EnvironmentId:   envAPIKey.Environment.Id,
		Id:              req.Id,
		Name:            req.Name,
		Description:     req.Description,
		Tags:            req.Tags,
		Enabled:         req.Enabled,
		Archived:        req.Archived,
		Variations:      req.Variations,
		Prerequisites:   req.Prerequisites,
		Targets:         req.Targets,
		Rules:           req.Rules,
		DefaultStrategy: req.DefaultStrategy,
		OffVariation:    req.OffVariation,
	})
	if err != nil {
		return nil, err
	}
	return &gwproto.UpdateFeatureResponse{
		Feature: res.Feature,
	}, nil
}
