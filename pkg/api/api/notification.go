package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	notificationproto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func (s *grpcGatewayService) ListSubscriptions(
	ctx context.Context,
	request *gwproto.ListSubscriptionsRequest,
) (*gwproto.ListSubscriptionsResponse, error) {
	envAPIKey, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check ListSubscriptions request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	res, err := s.notificationClient.ListSubscriptions(
		ctx,
		&notificationproto.ListSubscriptionsRequest{
			EnvironmentId:  envAPIKey.Environment.Id,
			PageSize:       request.PageSize,
			Cursor:         request.Cursor,
			OrderBy:        request.OrderBy,
			OrderDirection: request.OrderDirection,
			SearchKeyword:  request.SearchKeyword,
		},
	)
	if err != nil {
		return nil, err
	}

	if res == nil {
		s.logger.Error("Failed to list subscriptions: nil response",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("environment_namespace", envAPIKey.Environment.Id),
				zap.String("search_keyword", request.SearchKeyword),
			)...,
		)
		return nil, ErrSubscriptionNotFound
	}

	return &gwproto.ListSubscriptionsResponse{
		Subscriptions: res.Subscriptions,
	}, err
}

func (s *grpcGatewayService) CreateSubscription(
	ctx context.Context,
	request *gwproto.CreateSubscriptionRequest,
) (*gwproto.CreateSubscriptionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *grpcGatewayService) GetSubscription(
	ctx context.Context,
	request *gwproto.GetSubscriptionRequest,
) (*gwproto.GetSubscriptionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *grpcGatewayService) UpdateSubscription(
	ctx context.Context,
	request *gwproto.UpdateSubscriptionRequest,
) (*gwproto.UpdateSubscriptionResponse, error) {
	//TODO implement me
	panic("implement me")
}
