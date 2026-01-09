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
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

func (s *grpcGatewayService) ListPushes(
	ctx context.Context,
	req *gwproto.ListPushesRequest,
) (*gwproto.ListPushesResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListPushes request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListPushes, "").Inc()

	res, err := s.pushClient.ListPushes(
		ctx,
		&pushproto.ListPushesRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			PageSize:       req.PageSize,
			Cursor:         req.Cursor,
			OrderBy:        req.OrderBy,
			OrderDirection: req.OrderDirection,
			SearchKeyword:  req.SearchKeyword,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to list pushes: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environment_id", envAPIKey.Environment.Id),
				zap.String("search_keyword", req.SearchKeyword),
			)...,
		)
		return nil, ErrPushNotFound
	}

	return &gwproto.ListPushesResponse{
		Pushes:     res.Pushes,
		Cursor:     res.Cursor,
		TotalCount: res.TotalCount,
	}, nil
}

func (s *grpcGatewayService) CreatePush(
	ctx context.Context,
	req *gwproto.CreatePushRequest,
) (*gwproto.CreatePushResponse, error) {
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
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreatePush, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	res, err := s.pushClient.CreatePush(
		ctx,
		&pushproto.CreatePushRequest{
			EnvironmentId:     envAPIKey.Environment.Id,
			Name:              req.Name,
			Tags:              req.Tags,
			FcmServiceAccount: req.FcmServiceAccount,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Not found created push",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("name", req.Name),
			)...,
		)
		return nil, errInternal
	}

	return &gwproto.CreatePushResponse{
		Push: res.Push,
	}, nil
}

func (s *grpcGatewayService) GetPush(
	ctx context.Context,
	req *gwproto.GetPushRequest,
) (*gwproto.GetPushResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetPush request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("pushId", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetPush, "").Inc()

	res, err := s.pushClient.GetPush(
		ctx,
		&pushproto.GetPushRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Push not found",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("name", req.Id),
			)...,
		)
		return nil, ErrPushNotFound
	}

	return &gwproto.GetPushResponse{
		Push: res.Push,
	}, nil
}

func (s *grpcGatewayService) UpdatePush(
	ctx context.Context,
	req *gwproto.UpdatePushRequest,
) (*gwproto.UpdatePushResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check UpdatePush request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("pushId", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdatePush, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)

	if req.Deleted != nil && req.Deleted.Value {
		_, err := s.pushClient.DeletePush(
			ctx,
			&pushproto.DeletePushRequest{
				EnvironmentId: envAPIKey.Environment.Id,
				Id:            req.Id,
			},
		)
		if err != nil {
			return nil, err
		}

		return &gwproto.UpdatePushResponse{}, nil
	}

	res, err := s.pushClient.UpdatePush(
		ctx,
		&pushproto.UpdatePushRequest{
			EnvironmentId: envAPIKey.Environment.Id,
			Id:            req.Id,
			Name:          req.Name,
			Tags:          req.Tags,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Not found updated push",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("pushId", req.Id),
			)...,
		)
		return nil, errInternal
	}

	return &gwproto.UpdatePushResponse{
		Push: res.Push,
	}, nil
}
