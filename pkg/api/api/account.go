package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) CreateAccountV2(
	ctx context.Context,
	request *gwproto.CreateAccountV2Request,
) (*gwproto.CreateAccountV2Response, error) {
	//TODO implement me
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}

func (s *grpcGatewayService) GetAccountV2(
	ctx context.Context,
	request *gwproto.GetAccountV2Request,
) (*gwproto.GetAccountV2Response, error) {
	//TODO implement me
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}

func (s *grpcGatewayService) GetAccountV2ByEnvironmentID(
	ctx context.Context,
	request *gwproto.GetAccountV2ByEnvironmentIDRequest,
) (*gwproto.GetAccountV2ByEnvironmentIDResponse, error) {
	//TODO implement me
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}

func (s *grpcGatewayService) UpdateAccountV2(
	ctx context.Context,
	request *gwproto.UpdateAccountV2Request,
) (*gwproto.UpdateAccountV2Response, error) {
	//TODO implement me
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}

func (s *grpcGatewayService) ListAccountsV2(
	ctx context.Context,
	request *gwproto.ListAccountsV2Request,
) (*gwproto.ListAccountsV2Response, error) {
	//TODO implement me
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}
