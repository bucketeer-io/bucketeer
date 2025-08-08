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
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	teamproto "github.com/bucketeer-io/bucketeer/proto/team"
)

func (s *grpcGatewayService) CreateTeam(
	ctx context.Context,
	req *gwproto.CreateTeamRequest,
) (*gwproto.CreateTeamResponse, error) {
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

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	res, err := s.teamClient.CreateTeam(
		ctx,
		&teamproto.CreateTeamRequest{
			OrganizationId: req.OrganizationId,
			Name:           req.Name,
			Description:    req.Description,
		},
	)
	if err != nil {
		s.logger.Error("Failed to create team",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", req.Name),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to create team: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("name", req.Name),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.CreateTeamResponse{
		Team: res.Team,
	}, nil
}

func (s *grpcGatewayService) DeleteTeam(
	ctx context.Context,
	req *gwproto.DeleteTeamRequest,
) (*gwproto.DeleteTeamResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteTeam request",
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
	_, err = s.teamClient.DeleteTeam(
		ctx,
		&teamproto.DeleteTeamRequest{
			OrganizationId: req.OrganizationId,
			Id:             req.Id,
		},
	)
	if err != nil {
		s.logger.Error("Failed to delete team",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.DeleteTeamResponse{}, nil
}

func (s *grpcGatewayService) ListTeams(
	ctx context.Context,
	req *gwproto.ListTeamsRequest,
) (*gwproto.ListTeamsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListTeams request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	res, err := s.teamClient.ListTeams(
		ctx,
		&teamproto.ListTeamsRequest{
			PageSize:       req.PageSize,
			Cursor:         req.Cursor,
			OrderBy:        req.OrderBy,
			OrderDirection: req.OrderDirection,
			SearchKeyword:  req.SearchKeyword,
			OrganizationId: req.OrganizationId,
		},
	)
	if err != nil {
		s.logger.Error("Failed to list teams",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("search_keyword", req.SearchKeyword),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to list teams: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("environment_id", envAPIKey.Environment.Id),
				zap.String("search_keyword", req.SearchKeyword),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListTeamsResponse{
		Teams:      res.Teams,
		Cursor:     res.NextCursor,
		TotalCount: res.TotalCount,
	}, nil
}
