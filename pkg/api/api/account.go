package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) CreateAccountV2(
	ctx context.Context,
	request *gwproto.CreateAccountV2Request,
) (*gwproto.CreateAccountV2Response, error) {
	//TODO implement me
	panic("implement me")
}

func (s *grpcGatewayService) UpdateAccountV2(
	ctx context.Context,
	request *gwproto.UpdateAccountV2Request,
) (*gwproto.UpdateAccountV2Response, error) {
	//TODO implement me
	panic("implement me")
}

func (s *grpcGatewayService) GetAccountV2(
	ctx context.Context,
	request *gwproto.GetAccountV2Request,
) (*gwproto.GetAccountV2Response, error) {
	_, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get account request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", request.Email),
				zap.String("organizationId", request.OrganizationId),
			)...,
		)
		return nil, err
	}

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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
	_, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get account by environment request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", request.Email),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}

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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
	_, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check get my account request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

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
			log.FieldsFromImcomingContext(ctx).AddFields(
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
	_, err := s.checkRequest(ctx, []accountproto.APIKey_Role{
		accountproto.APIKey_PUBLIC_API_READ_ONLY,
		accountproto.APIKey_PUBLIC_API_WRITE,
		accountproto.APIKey_PUBLIC_API_ADMIN,
	})
	if err != nil {
		s.logger.Error("Failed to check list accounts request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, err
	}

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
		},
	)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Error("Failed to list accounts: nil response",
			log.FieldsFromImcomingContext(ctx).AddFields(
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
