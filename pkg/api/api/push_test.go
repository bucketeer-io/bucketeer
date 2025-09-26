package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	pushclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/push/client/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

func TestGrpcListPushes(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListPushesResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: listPushes error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().ListPushes(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().ListPushes(gomock.Any(), gomock.Any()).Return(
					&pushproto.ListPushesResponse{Pushes: []*pushproto.Push{
						{Id: "id-0"}, {Id: "id-1"},
					}}, nil)
			},
			expected: &gwproto.ListPushesResponse{Pushes: []*pushproto.Push{
				{Id: "id-0"}, {Id: "id-1"},
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.ListPushes(ctx, &gwproto.ListPushesRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcCreatePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		req         *gwproto.CreatePushRequest
		expected    *gwproto.CreatePushResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: create feature error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().CreatePush(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().CreatePush(gomock.Any(), gomock.Any()).Return(
					&pushproto.CreatePushResponse{
						Push: &pushproto.Push{Id: "id-0"},
					},
					nil)
			},
			expected: &gwproto.CreatePushResponse{
				Push: &pushproto.Push{Id: "id-0"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.CreatePush(ctx, &gwproto.CreatePushRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetPush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		req         *gwproto.GetPushRequest
		expected    *gwproto.GetPushResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_SERVER,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: getPush error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().GetPush(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().GetPush(gomock.Any(), gomock.Any()).Return(
					&pushproto.GetPushResponse{
						Push: &pushproto.Push{
							Id:      "id-0",
							Tags:    []string{"tag-0", "tag-1"},
							Name:    "name-0",
							Deleted: false,
						},
					}, nil)
			},
			expected: &gwproto.GetPushResponse{
				Push: &pushproto.Push{
					Id:      "id-0",
					Tags:    []string{"tag-0", "tag-1"},
					Name:    "name-0",
					Deleted: false,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetPush(ctx, &gwproto.GetPushRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcUpdatePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		req         *gwproto.UpdatePushRequest
		expected    *gwproto.UpdatePushResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "fail: update push error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().UpdatePush(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "success update",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.pushClient.(*pushclientmock.MockClient).EXPECT().UpdatePush(gomock.Any(), gomock.Any()).Return(
					&pushproto.UpdatePushResponse{
						Push: &pushproto.Push{
							Id:        "id-0",
							Tags:      []string{"tag-0"},
							Deleted:   false,
							Name:      "name-0",
							CreatedAt: 1729672123,
							UpdatedAt: 1729672123,
						},
					},
					nil)
			},
			expected: &gwproto.UpdatePushResponse{
				Push: &pushproto.Push{
					Id:        "id-0",
					Tags:      []string{"tag-0"},
					Deleted:   false,
					Name:      "name-0",
					CreatedAt: 1729672123,
					UpdatedAt: 1729672123,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.UpdatePush(ctx, &gwproto.UpdatePushRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

// TestGrpcDeletePush tests the case when the push is successfully deleted via updatePush api.
func TestGrpcDeletePush(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	setup := func(gs *grpcGatewayService) {
		gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
			&accountproto.EnvironmentAPIKey{
				Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
				ApiKey: &accountproto.APIKey{
					Id:       "id-0",
					Role:     accountproto.APIKey_PUBLIC_API_WRITE,
					Disabled: false,
				},
			},
			nil,
		)
		gs.pushClient.(*pushclientmock.MockClient).EXPECT().DeletePush(gomock.Any(), gomock.Any()).Return(
			&pushproto.DeletePushResponse{},
			nil,
		)
	}
	gs := newGrpcGatewayServiceWithMock(t, mockController)
	setup(gs)
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
		"authorization": []string{"test-key"},
	})
	expected := &gwproto.UpdatePushResponse{}
	var expectedErr error = nil

	actual, err := gs.UpdatePush(ctx, &gwproto.UpdatePushRequest{
		Deleted: wrapperspb.Bool(true),
	})
	assert.Equal(t, expected, actual, "%s", "success delete")
	assert.Equal(t, expectedErr, err, "%s", "success delete")
}
