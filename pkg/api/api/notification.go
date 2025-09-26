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
	"github.com/bucketeer-io/bucketeer/v2/proto/notification"
)

func (s *grpcGatewayService) GetSubscription(
	ctx context.Context,
	req *gwproto.GetSubscriptionRequest,
) (*gwproto.GetSubscriptionResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check GetSubscription request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodGetSubscription, "").Inc()

	res, err := s.notificationClient.GetSubscription(ctx, &notification.GetSubscriptionRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to get subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("subscriptionId", req.Id),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("GetSubscription returned nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("ID", req.Id),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.GetSubscriptionResponse{
		Subscription: res.Subscription,
	}, nil
}

func (s *grpcGatewayService) ListSubscriptions(
	ctx context.Context,
	req *gwproto.ListSubscriptionsRequest,
) (*gwproto.ListSubscriptionsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListSubscriptions request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodListSubscriptions, "").Inc()

	res, err := s.notificationClient.ListSubscriptions(ctx, &notification.ListSubscriptionsRequest{
		PageSize:       req.PageSize,
		Cursor:         req.Cursor,
		OrderBy:        req.OrderBy,
		OrderDirection: req.OrderDirection,
		SourceTypes:    req.SourceTypes,
		SearchKeyword:  req.SearchKeyword,
		Disabled:       req.Disabled,
		EnvironmentIds: req.EnvironmentIds,
		OrganizationId: req.OrganizationId,
	})
	if err != nil {
		s.logger.Error("Failed to list subscriptions",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("ListSubscriptions returned nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("organizationId", req.OrganizationId),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.ListSubscriptionsResponse{
		Subscriptions: res.Subscriptions,
		Cursor:        res.Cursor,
		TotalCount:    res.TotalCount,
	}, nil
}

func (s *grpcGatewayService) CreateSubscription(
	ctx context.Context,
	req *gwproto.CreateSubscriptionRequest,
) (*gwproto.CreateSubscriptionResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check CreateSubscription request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodCreateSubscription, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	res, err := s.notificationClient.CreateSubscription(ctx, &notification.CreateSubscriptionRequest{
		EnvironmentId:   envAPIKey.Environment.Id,
		Name:            req.Name,
		SourceTypes:     req.SourceTypes,
		Recipient:       req.Recipient,
		FeatureFlagTags: req.FeatureFlagTags,
	})
	if err != nil {
		s.logger.Error("Failed to create subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", req.Name),
			)...,
		)
		return nil, err
	}
	if res == nil {
		s.logger.Error("CreateSubscription returned nil response",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("name", req.Name),
			)...,
		)
		return nil, ErrInternal
	}
	return &gwproto.CreateSubscriptionResponse{
		Subscription: res.Subscription,
	}, nil
}

func (s *grpcGatewayService) DeleteSubscription(
	ctx context.Context,
	req *gwproto.DeleteSubscriptionRequest,
) (*gwproto.DeleteSubscriptionResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check DeleteSubscription request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("subscriptionId", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodDeleteSubscription, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.notificationClient.DeleteSubscription(ctx, &notification.DeleteSubscriptionRequest{
		EnvironmentId: envAPIKey.Environment.Id,
		Id:            req.Id,
	})
	if err != nil {
		s.logger.Error("Failed to delete subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("subscriptionId", req.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.DeleteSubscriptionResponse{}, nil
}

func (s *grpcGatewayService) UpdateSubscription(
	ctx context.Context,
	req *gwproto.UpdateSubscriptionRequest,
) (*gwproto.UpdateSubscriptionResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check UpdateSubscription request",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("subscriptionId", req.Id),
			)...,
		)
		return nil, err
	}

	requestTotal.WithLabelValues(
		envAPIKey.Environment.OrganizationId, envAPIKey.ProjectId, envAPIKey.ProjectUrlCode,
		envAPIKey.Environment.Id, envAPIKey.Environment.UrlCode, methodUpdateSubscription, "").Inc()

	headerMetaData := metadata.New(map[string]string{
		role.APIKeyTokenMDKey:      envAPIKey.ApiKey.ApiKey,
		role.APIKeyMaintainerMDKey: envAPIKey.ApiKey.Maintainer,
		role.APIKeyNameMDKey:       envAPIKey.ApiKey.Name,
	})
	ctx = metadata.NewOutgoingContext(ctx, headerMetaData)
	_, err = s.notificationClient.UpdateSubscription(ctx, &notification.UpdateSubscriptionRequest{
		EnvironmentId:   envAPIKey.Environment.Id,
		Id:              req.Id,
		Name:            req.Name,
		SourceTypes:     req.SourceTypes,
		Disabled:        req.Disabled,
		FeatureFlagTags: req.FeatureFlagTags,
	})
	if err != nil {
		s.logger.Error("Failed to update subscription",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("subscriptionId", req.Id),
			)...,
		)
		return nil, err
	}
	return &gwproto.UpdateSubscriptionResponse{}, nil
}
