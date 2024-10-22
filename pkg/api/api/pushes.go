package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func (s *grpcGatewayService) ListPushes(ctx context.Context, request *gwproto.ListPushesRequest) (*gwproto.ListPushesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}

func (s *grpcGatewayService) CreatePush(ctx context.Context, request *gwproto.CreatePushRequest) (*gwproto.CreatePushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}

func (s *grpcGatewayService) DeletePush(ctx context.Context, request *gwproto.DeletePushRequest) (*gwproto.DeletePushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}

func (s *grpcGatewayService) UpdatePush(ctx context.Context, request *gwproto.UpdatePushRequest) (*gwproto.UpdatePushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unsupported method")
}
