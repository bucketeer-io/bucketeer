package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	environmentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	gwproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
)

func TestGrpcCreateAccountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.CreateAccountV2Response
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().CreateAccountV2(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "Success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().CreateAccountV2(gomock.Any(), gomock.Any()).Return(
					&accountproto.CreateAccountV2Response{
						Account: &accountproto.AccountV2{
							Email:            "demo@bucketeer.io",
							FirstName:        "firstName",
							LastName:         "lastName",
							OrganizationId:   "org-0",
							OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						},
					}, nil)
			},
			expected: &gwproto.CreateAccountV2Response{
				Account: &accountproto.AccountV2{
					Email:            "demo@bucketeer.io",
					FirstName:        "firstName",
					LastName:         "lastName",
					OrganizationId:   "org-0",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.CreateAccountV2(ctx, &gwproto.CreateAccountV2Request{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcUpdateAccountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		req         *gwproto.UpdateAccountV2Request
		setup       func(*grpcGatewayService)
		expected    *gwproto.UpdateAccountV2Response
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			req:         &gwproto.UpdateAccountV2Request{},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			req:         &gwproto.UpdateAccountV2Request{},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "Account not found",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(
					nil, ErrAccountNotFound)
			},
			req:         &gwproto.UpdateAccountV2Request{},
			expected:    nil,
			expectedErr: ErrAccountNotFound,
		},
		{
			desc: "Success: delete account",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().DeleteAccountV2(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
			req: &gwproto.UpdateAccountV2Request{
				Deleted:        wrapperspb.Bool(true),
				Email:          "test@bucketeer.io",
				OrganizationId: "org-0",
			},
			expected:    &gwproto.UpdateAccountV2Response{},
			expectedErr: nil,
		},
		{
			desc: "Success: disable account",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(
					&accountproto.UpdateAccountV2Response{
						Account: &accountproto.AccountV2{
							Disabled:       true,
							Email:          "test@bucketeer.io",
							OrganizationId: "org-0",
						},
					}, nil)
			},
			req: &gwproto.UpdateAccountV2Request{
				Disabled:       wrapperspb.Bool(true),
				Email:          "test@bucketeer.io",
				OrganizationId: "org-0",
			},
			expected: &gwproto.UpdateAccountV2Response{
				Account: &accountproto.AccountV2{
					Disabled:       true,
					Email:          "test@bucketeer.io",
					OrganizationId: "org-0",
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Success: enable account",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(
					&accountproto.UpdateAccountV2Response{
						Account: &accountproto.AccountV2{
							Disabled:       false,
							Email:          "test@bucketeer.io",
							OrganizationId: "org-0",
						},
					}, nil)
			},
			req: &gwproto.UpdateAccountV2Request{
				Disabled:       wrapperspb.Bool(false),
				Email:          "test@bucketeer.io",
				OrganizationId: "org-0",
			},
			expected: &gwproto.UpdateAccountV2Response{
				Account: &accountproto.AccountV2{
					Disabled:       false,
					Email:          "test@bucketeer.io",
					OrganizationId: "org-0",
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Success: update account",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_WRITE,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().UpdateAccountV2(gomock.Any(), gomock.Any()).Return(
					&accountproto.UpdateAccountV2Response{
						Account: &accountproto.AccountV2{
							Email:            "demo@bucketeer.io",
							FirstName:        "newFirstName",
							LastName:         "lastName",
							OrganizationId:   "org-0",
							OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						},
					}, nil)
			},
			req: &gwproto.UpdateAccountV2Request{
				Email:          "demo@bucketeer.io",
				OrganizationId: "org-0",
				FirstName:      wrapperspb.String("newFirstName"),
			},
			expected: &gwproto.UpdateAccountV2Response{
				Account: &accountproto.AccountV2{
					Email:            "demo@bucketeer.io",
					FirstName:        "newFirstName",
					LastName:         "lastName",
					OrganizationId:   "org-0",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.UpdateAccountV2(ctx, p.req)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetAccountV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetAccountV2Response
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "Account not found",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
			expected:    nil,
			expectedErr: ErrAccountNotFound,
		},
		{
			desc: "Success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(gomock.Any(), gomock.Any()).Return(
					&accountproto.GetAccountV2Response{
						Account: &accountproto.AccountV2{
							Email:            "demo@bucketeer.io",
							Name:             "demo",
							OrganizationId:   "org-0",
							OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						},
					}, nil)
			},
			expected: &gwproto.GetAccountV2Response{
				Account: &accountproto.AccountV2{
					Email:            "demo@bucketeer.io",
					Name:             "demo",
					OrganizationId:   "org-0",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetAccountV2(ctx, &gwproto.GetAccountV2Request{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetAccountV2ByEnvironmentID(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetAccountV2ByEnvironmentIDResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "Account not found",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
			expected:    nil,
			expectedErr: ErrAccountNotFound,
		},
		{
			desc: "Success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(
					&accountproto.GetAccountV2ByEnvironmentIDResponse{
						Account: &accountproto.AccountV2{
							Email:            "demo@bucketeer.io",
							Name:             "demo",
							OrganizationId:   "org-0",
							OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						},
					}, nil)
			},
			expected: &gwproto.GetAccountV2ByEnvironmentIDResponse{
				Account: &accountproto.AccountV2{
					Email:            "demo@bucketeer.io",
					Name:             "demo",
					OrganizationId:   "org-0",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetAccountV2ByEnvironmentID(ctx, &gwproto.GetAccountV2ByEnvironmentIDRequest{EnvironmentId: "ns0"})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetMe(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.GetMeResponse
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetMe(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "Account not found",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetMe(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
			expected:    nil,
			expectedErr: ErrAccountNotFound,
		},
		{
			desc: "Success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetMe(gomock.Any(), gomock.Any()).Return(
					&accountproto.GetMeResponse{
						Account: &accountproto.ConsoleAccount{
							Email:            "demo@bucketeer.io",
							Name:             "demo",
							OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						},
					}, nil)
			},
			expected: &gwproto.GetMeResponse{
				Account: &accountproto.ConsoleAccount{
					Email:            "demo@bucketeer.io",
					Name:             "demo",
					OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.GetMe(ctx, &gwproto.GetMeRequest{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcGetMeOrganizationScope(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name              string
		organizationID    string
		keyOrganizationID string
		wantErr           error
	}{
		{name: "same organization", organizationID: "org-a", keyOrganizationID: "org-a"},
		{name: "omitted organization", keyOrganizationID: "org-a"},
		{name: "other organization", organizationID: "org-b", keyOrganizationID: "org-a", wantErr: ErrNotFound},
		{name: "key without organization", wantErr: ErrInvalidAPIKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gs, ctx := newAccountReadGateway(t, tc.keyOrganizationID)
			account := &accountproto.ConsoleAccount{Email: "member@example.com", Name: "member"}
			if tc.wantErr == nil {
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetMe(
					accountReadMetadataMatcher(), &accountproto.GetMeRequest{OrganizationId: "org-a"},
				).Return(&accountproto.GetMeResponse{Account: account}, nil)
			}
			res, err := gs.GetMe(ctx, &gwproto.GetMeRequest{OrganizationId: tc.organizationID})
			require.Equal(t, tc.wantErr, err)
			if tc.wantErr != nil {
				require.Nil(t, res)
				return
			}
			require.Equal(t, account, res.Account)
		})
	}
}

func TestGrpcListAccountsV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		ctx         context.Context
		setup       func(*grpcGatewayService)
		expected    *gwproto.ListAccountsV2Response
		expectedErr error
	}{
		{
			desc: "fails: bad role",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_SDK_CLIENT,
							Disabled: false,
						},
					}, nil)
			},
			expected:    nil,
			expectedErr: ErrBadRole,
		},
		{
			desc: "Internal account grpc error",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().ListAccountsV2(gomock.Any(), gomock.Any()).Return(
					nil, ErrInternal)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "Response is nil",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().ListAccountsV2(gomock.Any(), gomock.Any()).Return(
					nil, nil)
			},
			expected:    nil,
			expectedErr: ErrInternal,
		},
		{
			desc: "Success",
			setup: func(gs *grpcGatewayService) {
				gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(gomock.Any()).Return(
					&accountproto.EnvironmentAPIKey{
						Environment: &environmentproto.EnvironmentV2{Id: "ns0", OrganizationId: "org-0"},
						ApiKey: &accountproto.APIKey{
							Id:       "id-0",
							Role:     accountproto.APIKey_PUBLIC_API_READ_ONLY,
							Disabled: false,
						},
					}, nil)
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().ListAccountsV2(gomock.Any(), gomock.Any()).Return(
					&accountproto.ListAccountsV2Response{
						Accounts: []*accountproto.AccountV2{
							{
								Email:            "demo@bucketeer.io",
								Name:             "demo",
								OrganizationId:   "org-0",
								OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
							},
							{
								Email:            "bucketeer@demo.io",
								Name:             "bucketeer",
								OrganizationId:   "org-0",
								OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
							},
						},
						Cursor:     "0",
						TotalCount: 2,
					}, nil)
			},
			expected: &gwproto.ListAccountsV2Response{
				Accounts: []*accountproto.AccountV2{
					{
						Email:            "demo@bucketeer.io",
						Name:             "demo",
						OrganizationId:   "org-0",
						OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
					},
					{
						Email:            "bucketeer@demo.io",
						Name:             "bucketeer",
						OrganizationId:   "org-0",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					},
				},
				Cursor:     "0",
				TotalCount: 2,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			gs := newGrpcGatewayServiceWithMock(t, mockController)
			p.setup(gs)
			ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{
				"authorization": []string{"test-key"},
			})
			actual, err := gs.ListAccountsV2(ctx, &gwproto.ListAccountsV2Request{})
			assert.Equal(t, p.expected, actual, "%s", p.desc)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func TestGrpcAccountReadsOrganizationScope(t *testing.T) {
	t.Parallel()

	for _, method := range []string{"get", "list"} {
		for _, tc := range []struct {
			name              string
			organizationID    string
			keyOrganizationID string
			wantErr           error
		}{
			{name: "same organization", organizationID: "org-a", keyOrganizationID: "org-a"},
			{name: "omitted organization", keyOrganizationID: "org-a"},
			{name: "other organization", organizationID: "org-b", keyOrganizationID: "org-a", wantErr: ErrNotFound},
			{name: "key without organization", wantErr: ErrInvalidAPIKey},
		} {
			t.Run(method+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				gs, ctx := newAccountReadGateway(t, tc.keyOrganizationID)
				account := &accountproto.AccountV2{Email: "member@example.com", OrganizationId: "org-a"}
				client := gs.accountClient.(*accountclientmock.MockClient)
				if method == "get" {
					if tc.wantErr == nil {
						client.EXPECT().GetAccountV2(accountReadMetadataMatcher(), &accountproto.GetAccountV2Request{
							Email: account.Email, OrganizationId: "org-a",
						}).Return(&accountproto.GetAccountV2Response{Account: account}, nil)
					}
					res, err := gs.GetAccountV2(ctx, &gwproto.GetAccountV2Request{
						Email: account.Email, OrganizationId: tc.organizationID,
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, account, res.Account)
				} else {
					if tc.wantErr == nil {
						client.EXPECT().ListAccountsV2(accountReadMetadataMatcher(), &accountproto.ListAccountsV2Request{
							OrganizationId: "org-a", PageSize: 10,
						}).Return(&accountproto.ListAccountsV2Response{
							Accounts: []*accountproto.AccountV2{account}, Cursor: "1", TotalCount: 1,
						}, nil)
					}
					res, err := gs.ListAccountsV2(ctx, &gwproto.ListAccountsV2Request{
						OrganizationId: tc.organizationID, PageSize: 10,
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, []*accountproto.AccountV2{account}, res.Accounts)
					require.Equal(t, "1", res.Cursor)
					require.Equal(t, int64(1), res.TotalCount)
				}
			})
		}
	}
}

func TestGrpcGetAccountV2ByEnvironmentIDOrganizationScope(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name              string
		environmentID     string
		keyOrganizationID string
		environment       *environmentproto.GetEnvironmentV2Response
		environmentErr    error
		wantErr           error
	}{
		{name: "key environment", environmentID: "env-a", keyOrganizationID: "org-a"},
		{
			name: "another environment in same organization", environmentID: "env-a2", keyOrganizationID: "org-a",
			environment: &environmentproto.GetEnvironmentV2Response{
				Environment: &environmentproto.EnvironmentV2{Id: "env-a2", OrganizationId: "org-a"},
			},
		},
		{
			name: "other organization even when maintainer can read it", environmentID: "env-b", keyOrganizationID: "org-a",
			environment: &environmentproto.GetEnvironmentV2Response{
				Environment: &environmentproto.EnvironmentV2{Id: "env-b", OrganizationId: "org-b"},
			},
			wantErr: ErrNotFound,
		},
		{name: "environment lookup error", environmentID: "env-b", keyOrganizationID: "org-a", environmentErr: ErrNotFound, wantErr: ErrNotFound},
		{
			// An unknown environment is reported exactly like one in another organization.
			name: "unknown environment", environmentID: "env-b", keyOrganizationID: "org-a",
			environmentErr: status.Error(codes.NotFound, "environment: not found"), wantErr: ErrNotFound,
		},
		{name: "environment lookup internal error", environmentID: "env-b", keyOrganizationID: "org-a", environmentErr: ErrInternal, wantErr: ErrInternal},
		{name: "nil environment response", environmentID: "env-b", keyOrganizationID: "org-a", wantErr: ErrInternal},
		{name: "nil environment", environmentID: "env-b", keyOrganizationID: "org-a", environment: &environmentproto.GetEnvironmentV2Response{}, wantErr: ErrInternal},
		{name: "key without organization", environmentID: "env-a", wantErr: ErrInvalidAPIKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gs, ctx := newAccountReadGateway(t, tc.keyOrganizationID)
			if tc.keyOrganizationID != "" && tc.environmentID != "env-a" {
				gs.environmentClient.(*environmentclientmock.MockClient).EXPECT().GetEnvironmentV2(
					accountReadMetadataMatcher(), &environmentproto.GetEnvironmentV2Request{Id: tc.environmentID},
				).Return(tc.environment, tc.environmentErr)
			}
			account := &accountproto.AccountV2{Email: "member@example.com", OrganizationId: "org-a"}
			if tc.wantErr == nil {
				gs.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2ByEnvironmentID(
					accountReadMetadataMatcher(), &accountproto.GetAccountV2ByEnvironmentIDRequest{
						Email: account.Email, EnvironmentId: tc.environmentID,
					},
				).Return(&accountproto.GetAccountV2ByEnvironmentIDResponse{Account: account}, nil)
			}
			res, err := gs.GetAccountV2ByEnvironmentID(ctx, &gwproto.GetAccountV2ByEnvironmentIDRequest{
				Email: account.Email, EnvironmentId: tc.environmentID,
			})
			require.Equal(t, tc.wantErr, err)
			if tc.wantErr != nil {
				require.Nil(t, res)
				return
			}
			require.Equal(t, account, res.Account)
		})
	}
}

func TestGrpcAccountWritesOrganizationScope(t *testing.T) {
	t.Parallel()

	for _, method := range []string{"create", "update", "delete"} {
		for _, tc := range []struct {
			name              string
			organizationID    string
			keyOrganizationID string
			wantErr           error
		}{
			{name: "same organization", organizationID: "org-a", keyOrganizationID: "org-a"},
			{name: "omitted organization", keyOrganizationID: "org-a"},
			{name: "other organization even when maintainer can write it", organizationID: "org-b", keyOrganizationID: "org-a", wantErr: ErrNotFound},
			{name: "key without organization", wantErr: ErrInvalidAPIKey},
		} {
			t.Run(method+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				gs, ctx := newAccountWriteGateway(t, tc.keyOrganizationID)
				account := &accountproto.AccountV2{Email: "member@example.com", OrganizationId: "org-a"}
				client := gs.accountClient.(*accountclientmock.MockClient)
				switch method {
				case "create":
					if tc.wantErr == nil {
						client.EXPECT().CreateAccountV2(accountWriteMetadataMatcher(), &accountproto.CreateAccountV2Request{
							Email: account.Email, OrganizationId: "org-a",
							OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						}).Return(&accountproto.CreateAccountV2Response{Account: account}, nil)
					}
					res, err := gs.CreateAccountV2(ctx, &gwproto.CreateAccountV2Request{
						Email: account.Email, OrganizationId: tc.organizationID,
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, account, res.Account)
				case "update":
					if tc.wantErr == nil {
						client.EXPECT().UpdateAccountV2(accountWriteMetadataMatcher(), &accountproto.UpdateAccountV2Request{
							Email: account.Email, OrganizationId: "org-a", Name: wrapperspb.String("renamed"),
						}).Return(&accountproto.UpdateAccountV2Response{Account: account}, nil)
					}
					res, err := gs.UpdateAccountV2(ctx, &gwproto.UpdateAccountV2Request{
						Email: account.Email, OrganizationId: tc.organizationID, Name: wrapperspb.String("renamed"),
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, account, res.Account)
				case "delete":
					if tc.wantErr == nil {
						client.EXPECT().DeleteAccountV2(accountWriteMetadataMatcher(), &accountproto.DeleteAccountV2Request{
							Email: account.Email, OrganizationId: "org-a",
						}).Return(&accountproto.DeleteAccountV2Response{}, nil)
					}
					res, err := gs.UpdateAccountV2(ctx, &gwproto.UpdateAccountV2Request{
						Email: account.Email, OrganizationId: tc.organizationID, Deleted: wrapperspb.Bool(true),
					})
					require.Equal(t, tc.wantErr, err)
					if tc.wantErr != nil {
						require.Nil(t, res)
						return
					}
					require.Equal(t, &gwproto.UpdateAccountV2Response{}, res)
				}
			})
		}
	}
}

func newAccountReadGateway(t *testing.T, organizationID string) (*grpcGatewayService, context.Context) {
	t.Helper()
	return newAccountGateway(t, organizationID, accountproto.APIKey_PUBLIC_API_READ_ONLY, "read-only-key")
}

func newAccountWriteGateway(t *testing.T, organizationID string) (*grpcGatewayService, context.Context) {
	t.Helper()
	return newAccountGateway(t, organizationID, accountproto.APIKey_PUBLIC_API_WRITE, "write-key")
}

// newAccountGateway returns a gateway whose cached key has the given role and
// belongs to the given organization, and a request context authenticated
// with that key.
func newAccountGateway(
	t *testing.T,
	organizationID string,
	keyRole accountproto.APIKey_Role,
	keyName string,
) (*grpcGatewayService, context.Context) {
	t.Helper()
	gs := newGrpcGatewayServiceWithMock(t, gomock.NewController(t))
	key := &accountproto.EnvironmentAPIKey{
		Environment: &environmentproto.EnvironmentV2{Id: "env-a", OrganizationId: organizationID},
		ApiKey: &accountproto.APIKey{
			Id: "key-a", ApiKey: "0123456789abcdef0123456789abcdef", Name: keyName,
			Role: keyRole, Maintainer: "maintainer@example.com",
		},
	}
	gs.environmentAPIKeyCache.(*cachev3mock.MockEnvironmentAPIKeyCache).EXPECT().Get(key.ApiKey.ApiKey).Return(key, nil)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", key.ApiKey.ApiKey))
	// Caller-supplied editor metadata must not replace the verified key's maintainer.
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(role.APIKeyMaintainerMDKey, "other-org-admin@example.com"))
	return gs, ctx
}

func accountReadMetadataMatcher() gomock.Matcher {
	return accountMetadataMatcher("read-only-key")
}

func accountWriteMetadataMatcher() gomock.Matcher {
	return accountMetadataMatcher("write-key")
}

// accountMetadataMatcher matches a context carrying the verified key's editor
// metadata for the key with the given name.
func accountMetadataMatcher(keyName string) gomock.Matcher {
	return gomock.Cond(func(ctx context.Context) bool {
		md, ok := metadata.FromOutgoingContext(ctx)
		return ok &&
			len(md.Get(role.APIKeyTokenMDKey)) == 1 && md.Get(role.APIKeyTokenMDKey)[0] == "0123....cdef" &&
			len(md.Get(role.APIKeyMaintainerMDKey)) == 1 && md.Get(role.APIKeyMaintainerMDKey)[0] == "maintainer@example.com" &&
			len(md.Get(role.APIKeyNameMDKey)) == 1 && md.Get(role.APIKeyNameMDKey)[0] == keyName
	})
}
