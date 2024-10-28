package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	pushes, err := s.pushClient.ListPushes(
		ctx,
		&pushproto.ListPushesRequest{
			EnvironmentNamespace: envAPIKey.Environment.Id,
			PageSize:             req.PageSize,
			Cursor:               req.Cursor,
			OrderBy:              req.OrderBy,
			OrderDirection:       req.OrderDirection,
			SearchKeyword:        req.SearchKeyword,
		},
	)
	if err != nil {
		return nil, err
	}

	return &gwproto.ListPushesResponse{
		Pushes: pushes.Pushes,
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("name", req.Name),
			)...,
		)
		return nil, err
	}
	res, err := s.pushClient.CreatePush(
		ctx,
		&pushproto.CreatePushRequest{
			EnvironmentNamespace: envAPIKey.Environment.Id,
			Name:                 req.Name,
			Tags:                 req.Tags,
			FcmServiceAccount:    req.FcmServiceAccount,
		},
	)
	if err != nil {
		return nil, err
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("pushId", req.Id),
			)...,
		)
		return nil, err
	}

	res, err := s.pushClient.GetPush(
		ctx,
		&pushproto.GetPushRequest{
			EnvironmentNamespace: envAPIKey.Environment.Id,
			Id:                   req.Id,
		},
	)
	if err != nil {
		return nil, err
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("pushId", req.Id),
				zap.String("name", req.Name),
			)...,
		)
		return nil, err
	}

	if req.Deleted != nil && req.Deleted.Value {
		_, err := s.pushClient.DeletePush(
			ctx,
			&pushproto.DeletePushRequest{
				EnvironmentNamespace: envAPIKey.Environment.Id,
				Id:                   req.Id,
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
			EnvironmentNamespace: envAPIKey.Environment.Id,
			Id:                   req.Id,
			Name:                 req.Name,
			Tags:                 req.Tags,
		},
	)
	if err != nil {
		return nil, err
	}

	return &gwproto.UpdatePushResponse{
		Push: res.Push,
	}, nil
}
