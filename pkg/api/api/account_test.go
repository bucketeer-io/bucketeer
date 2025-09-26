package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
			actual, err := gs.GetAccountV2ByEnvironmentID(ctx, &gwproto.GetAccountV2ByEnvironmentIDRequest{})
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
						Environment: &environmentproto.EnvironmentV2{Id: "ns0"},
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
