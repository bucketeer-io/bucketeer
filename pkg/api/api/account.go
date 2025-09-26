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
)

func (s *grpcGatewayService) CreateAccountV2(
	ctx context.Context,
	request *gwproto.CreateAccountV2Request,
) (*gwproto.CreateAccountV2Response, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check create account request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", request.Email),
				zap.String("organizationId", request.OrganizationId),
				zap.String("role", request.OrganizationRole.String()),
			)...,
		)
		return nil, err
	}
	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateAccountV2, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	res, err := s.accountClient.CreateAccountV2(
		ctx,
		&accountproto.CreateAccountV2Request{
			OrganizationId:   request.OrganizationId,
			Email:            request.Email,
			Name:             request.Name,
			AvatarImageUrl:   request.AvatarImageUrl,
			OrganizationRole: request.OrganizationRole,
			EnvironmentRoles: request.EnvironmentRoles,
			FirstName:        request.FirstName,
			LastName:         request.LastName,
			Language:         request.Language,
		},
	)
	if err != nil {
		return nil, err
	}

	if res == nil {
		s.logger.Error("Failed to create account: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("email", request.Email),
				zap.String("organizationId", request.OrganizationId),
				zap.String("role", request.OrganizationRole.String()),
			)...,
		)
		return nil, ErrInternal
	}

	return &gwproto.CreateAccountV2Response{
		Account: res.Account,
	}, nil
}

func (s *grpcGatewayService) UpdateAccountV2(
	ctx context.Context,
	request *gwproto.UpdateAccountV2Request,
) (*gwproto.UpdateAccountV2Response, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check update account request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", request.Email),
				zap.String("organizationId", request.OrganizationId),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdateAccountV2, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)

	// delete account
	if request.Deleted != nil && request.Deleted.Value {
		_, err := s.accountClient.DeleteAccountV2(
			ctx,
			&accountproto.DeleteAccountV2Request{
				Email:          request.Email,
				OrganizationId: request.OrganizationId,
			},
		)
		if err != nil {
			return nil, err
		}

		return &gwproto.UpdateAccountV2Response{}, nil
	}

	res, err := s.accountClient.UpdateAccountV2(
		ctx,
		&accountproto.UpdateAccountV2Request{
			OrganizationId:   request.OrganizationId,
			Email:            request.Email,
			Name:             request.Name,
			AvatarImageUrl:   request.AvatarImageUrl,
			OrganizationRole: request.OrganizationRole,
			EnvironmentRoles: request.EnvironmentRoles,
			FirstName:        request.FirstName,
			LastName:         request.LastName,
			Language:         request.Language,
			LastSeen:         request.LastSeen,
			Avatar:           request.Avatar,
			Disabled:         request.Disabled,
		},
	)
	if err != nil {
		return nil, err
	}

	if res == nil {
		s.logger.Error("Not found updated account",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("email", request.Email),
				zap.String("organizationId", request.OrganizationId),
			)...,
		)
		return nil, ErrAccountNotFound
	}

	return &gwproto.UpdateAccountV2Response{
		Account: res.Account,
	}, nil
}

func (s *grpcGatewayService) GetAccountV2(
	ctx context.Context,
	request *gwproto.GetAccountV2Request,
) (*gwproto.GetAccountV2Response, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get account request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", request.Email),
				zap.String("organizationId", request.OrganizationId),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetAccountV2, "").Inc()

	res, err := s.accountClient.GetAccountV2(
		ctx,
		&accountproto.GetAccountV2Request{
			Email:          request.Email,
			OrganizationId: request.OrganizationId,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Account not found",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("email", request.Email),
				zap.String("organizationId", request.OrganizationId),
			)...,
		)
		return nil, ErrAccountNotFound
	}

	return &gwproto.GetAccountV2Response{
		Account: res.Account,
	}, nil
}

func (s *grpcGatewayService) GetAccountV2ByEnvironmentID(
	ctx context.Context,
	request *gwproto.GetAccountV2ByEnvironmentIDRequest,
) (*gwproto.GetAccountV2ByEnvironmentIDResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get account by environment request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", request.Email),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetAccountV2ByEnvironmentID, "").Inc()

	res, err := s.accountClient.GetAccountV2ByEnvironmentID(
		ctx,
		&accountproto.GetAccountV2ByEnvironmentIDRequest{
			Email:         request.Email,
			EnvironmentId: request.EnvironmentId,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Account not found",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("email", request.Email),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, ErrAccountNotFound
	}

	return &gwproto.GetAccountV2ByEnvironmentIDResponse{
		Account: res.Account,
	}, nil
}

func (s *grpcGatewayService) GetMe(
	ctx context.Context,
	request *gwproto.GetMeRequest,
) (*gwproto.GetMeResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get my account request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetMe, "").Inc()

	res, err := s.accountClient.GetMe(
		ctx,
		&accountproto.GetMeRequest{
			OrganizationId: request.OrganizationId,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Account not found",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("organizationId", request.OrganizationId),
			)...,
		)
		return nil, ErrAccountNotFound
	}

	return &gwproto.GetMeResponse{
		Account: res.Account,
	}, nil
}

func (s *grpcGatewayService) ListAccountsV2(
	ctx context.Context,
	request *gwproto.ListAccountsV2Request,
) (*gwproto.ListAccountsV2Response, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check list accounts request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListAccountsV2, "").Inc()

	res, err := s.accountClient.ListAccountsV2(
		ctx,
		&accountproto.ListAccountsV2Request{
			PageSize:         request.PageSize,
			Cursor:           request.Cursor,
			OrganizationId:   request.OrganizationId,
			OrderBy:          request.OrderBy,
			OrderDirection:   request.OrderDirection,
			SearchKeyword:    request.SearchKeyword,
			Disabled:         request.Disabled,
			OrganizationRole: request.OrganizationRole,
			EnvironmentId:    request.EnvironmentId,
			EnvironmentRole:  request.EnvironmentRole,
			Teams:            request.Teams,
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to list accounts: nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Int64("pageSize", request.PageSize),
				zap.String("cursor", request.Cursor),
				zap.String("organizationId", request.OrganizationId),
				zap.String("searchKeyword", request.SearchKeyword),
			)...,
		)
		return nil, ErrInternal
	}

	return &gwproto.ListAccountsV2Response{
		Accounts:   res.Accounts,
		Cursor:     res.Cursor,
		TotalCount: res.TotalCount,
	}, nil
}
