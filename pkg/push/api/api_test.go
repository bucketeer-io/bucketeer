// Copyright 2026 The Bucketeer Authors.
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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	v2ps "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

var fcmServiceAccountDummy = []byte(`
	{
		"type": "service_account",
		"project_id": "test",
		"private_key_id": "private-key-id",
		"private_key": "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----\n",
		"client_email": "fcm-service-account@test.iam.gserviceaccount.com",
		"client_id": "client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fcm-service-account@test.iam.gserviceaccount.com",
		"universe_domain": "googleapis.com"
	}
`)

func TestNewPushService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mysqlClient := mysqlmock.NewMockClient(mockController)
	featureClientMock := featureclientmock.NewMockClient(mockController)
	experimentClientMock := experimentclientmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	pm := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewPushService(
		mysqlClient,
		featureClientMock,
		experimentClientMock,
		accountClientMock,
		pm,
		WithLogger(logger),
	)
	assert.IsType(t, &PushService{}, s)
}

func TestCreatePushMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)

	patterns := []struct {
		desc        string
		setup       func(*PushService)
		req         *pushproto.CreatePushRequest
		expectedErr error
	}{
		{
			desc:  "err: ErrFCMServiceAccountRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				FcmServiceAccount: nil,
			},
			expectedErr: statusFCMServiceAccountRequired.Err(),
		},

		{
			desc:  "err: ErrNameRequired",
			setup: nil,
			req: &pushproto.CreatePushRequest{
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{}, // Tags are now optional
			},
			expectedErr: statusNameRequired.Err(),
		},
		{
			desc: "err: ErrAlreadyExists",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2ps.ErrPushAlreadyExists)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().CreatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushAlreadyExists)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
				Name:              "name-1",
			},
			expectedErr: statusAlreadyExists.Err(),
		},
		{
			desc: "success: with tags",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().CreatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
				Name:              "name-1",
			},
			expectedErr: nil,
		},
		{
			desc: "success: without tags",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().CreatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{}, // Empty tags should be allowed
				Name:              "name-1",
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.CreatePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdatePushMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)

	patterns := []struct {
		desc        string
		setup       func(*PushService)
		req         *pushproto.UpdatePushRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &pushproto.UpdatePushRequest{},
			expectedErr: statusIDRequired.Err(),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ps.ErrPushNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.UpdatePushRequest{
				Id:   "key-1",
				Name: wrapperspb.String("push-0"),
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc: "success update name",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Push{
					Push: &proto.Push{
						Id:   "key-0",
						Name: "push-0",
						Tags: []string{"tag-0"},
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				Id:   "key-0",
				Name: wrapperspb.String("push-0"),
			},
			expectedErr: nil,
		},
		{
			desc: "success update tags",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Push{
					Push: &proto.Push{
						Id:   "key-0",
						Name: "push-0",
						Tags: []string{"tag-0"},
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				Id: "key-0",
				TagChanges: []*pushproto.TagChange{
					{
						ChangeType: pushproto.ChangeType_CREATE,
						Tag:        "tag-1",
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Push{
					Push: &proto.Push{
						Id:   "key-0",
						Name: "push-0",
						Tags: []string{"tag-0"},
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
				Name:          wrapperspb.String("name-1"),
				TagChanges: []*pushproto.TagChange{
					{
						ChangeType: pushproto.ChangeType_CREATE,
						Tag:        "tag-0",
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.UpdatePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCheckFCMServiceAccount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)

	patterns := []struct {
		desc              string
		fcmServiceAccount []byte
		pushes            []*pushproto.Push
		expected          error
	}{
		{
			desc:              "err: invalid service account",
			fcmServiceAccount: []byte(`"key":"value"`),
			pushes:            nil,
			expected:          statusFCMServiceAccountInvalid.Err(),
		},
		{
			desc:              "err: internal error",
			fcmServiceAccount: fcmServiceAccountDummy,
			pushes: []*pushproto.Push{
				{
					FcmServiceAccount: "`{\"key\":\"value\"}`",
				},
			},
			expected: statusInternal.Err(),
		},
		{
			desc:              "err: service account already exists",
			fcmServiceAccount: fcmServiceAccountDummy,
			pushes: []*pushproto.Push{
				{
					FcmServiceAccount: string(fcmServiceAccountDummy),
				},
			},
			expected: statusFCMServiceAccountAlreadyExists.Err(),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			err := service.checkFCMServiceAccount(
				ctx,
				p.pushes,
				p.fcmServiceAccount,
			)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestDeletePushMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)

	patterns := []struct {
		desc        string
		setup       func(*PushService)
		req         *pushproto.DeletePushRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &pushproto.DeletePushRequest{},
			expectedErr: statusIDRequired.Err(),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ps.ErrPushNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Push{
					Push: &proto.Push{
						Id: "key-0",
					},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().DeletePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.DeletePush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListPushesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, false)

	patterns := []struct {
		desc        string
		orgRole     *accountproto.AccountV2_Role_Organization
		envRole     *accountproto.AccountV2_Role_Environment
		setup       func(*PushService)
		input       *pushproto.ListPushesRequest
		expected    *pushproto.ListPushesResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrInvalidCursor",
			orgRole:     toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:     toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup:       nil,
			input:       &pushproto.ListPushesRequest{Cursor: "XXX", EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: statusInvalidCursor.Err(),
		},
		{
			desc: "err: ErrInternal",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("error"))
			},
			input:       &pushproto.ListPushesRequest{EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: statusInternal.Err(),
		},
		{
			desc:        "err: ErrPermissionDenied",
			orgRole:     toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:     toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			input:       &pushproto.ListPushesRequest{EnvironmentId: "ns0"},
			expected:    nil,
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "success",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
			},
			input:       &pushproto.ListPushesRequest{PageSize: 2, Cursor: "", EnvironmentId: "ns0"},
			expected:    &pushproto.ListPushesResponse{Pushes: []*pushproto.Push{}, Cursor: "0"},
			expectedErr: nil,
		},
		{
			desc:    "success: filter by environmentIDs",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_EDITOR),
			setup: func(s *PushService) {
				s.accountClient.(*accountclientmock.MockClient).EXPECT().GetAccountV2(
					gomock.Any(), gomock.Any(),
				).Return(&accountproto.GetAccountV2Response{
					Account: &accountproto.AccountV2{
						Email:            "email",
						OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
						EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
							{
								EnvironmentId: "ns0",
								Role:          accountproto.AccountV2_Role_Environment_EDITOR,
							},
							{
								EnvironmentId: "ns1",
								Role:          accountproto.AccountV2_Role_Environment_VIEWER,
							},
						},
					},
				}, nil)
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{
					{Id: "push-1", EnvironmentId: "ns0"},
					{Id: "push-2", EnvironmentId: "ns1"},
				}, 2, int64(2), nil)
			},
			input: &pushproto.ListPushesRequest{
				PageSize:       2,
				Cursor:         "",
				EnvironmentIds: []string{"ns0"},
				OrganizationId: "org-1",
			},
			expected: &pushproto.ListPushesResponse{
				Pushes: []*pushproto.Push{
					{Id: "push-1", EnvironmentId: "ns0"},
					{Id: "push-2", EnvironmentId: "ns1"},
				},
				Cursor:     "2",
				TotalCount: 2,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newPushService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(s)
			}

			actual, err := s.ListPushes(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGetPushMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken(t, true)

	patterns := []struct {
		desc        string
		setup       func(*PushService)
		req         *pushproto.GetPushRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &pushproto.GetPushRequest{},
			expectedErr: statusIDRequired.Err(),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ps.ErrPushNotFound)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Push{
					Push: &proto.Push{
						Id: "key-1",
					},
				}, nil)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := newPushServiceWithMock(t, mockController)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.GetPush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newPushServiceWithMock(t *testing.T, c *gomock.Controller) *PushService {
	t.Helper()
	return &PushService{
		mysqlClient:      mysqlmock.NewMockClient(c),
		pushStorage:      storagemock.NewMockPushStorage(c),
		featureClient:    featureclientmock.NewMockClient(c),
		experimentClient: experimentclientmock.NewMockClient(c),
		accountClient:    accountclientmock.NewMockClient(c),
		publisher:        publishermock.NewMockPublisher(c),
		logger:           zap.NewNop(),
	}
}

func newPushService(c *gomock.Controller, specifiedEnvironmentId *string, specifiedOrgRole *accountproto.AccountV2_Role_Organization, specifiedEnvRole *accountproto.AccountV2_Role_Environment) *PushService {
	var or accountproto.AccountV2_Role_Organization
	var er accountproto.AccountV2_Role_Environment
	var envId string
	if specifiedOrgRole != nil {
		or = *specifiedOrgRole
	} else {
		or = accountproto.AccountV2_Role_Organization_ADMIN
	}
	if specifiedEnvRole != nil {
		er = *specifiedEnvRole
	} else {
		er = accountproto.AccountV2_Role_Environment_EDITOR
	}
	if specifiedEnvironmentId != nil {
		envId = *specifiedEnvironmentId
	} else {
		envId = "ns0"
	}

	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: or,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: envId,
					Role:          er,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	mysqlClient := mysqlmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &PushService{
		mysqlClient:      mysqlClient,
		featureClient:    featureclientmock.NewMockClient(c),
		pushStorage:      storagemock.NewMockPushStorage(c),
		experimentClient: experimentclientmock.NewMockClient(c),
		accountClient:    accountClientMock,
		publisher:        publishermock.NewMockPublisher(c),
		logger:           zap.NewNop(),
	}
}

func createContextWithToken(t *testing.T, isSystemAdmin bool) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: isSystemAdmin,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}

// convert to pointer
func toPtr[T any](value T) *T {
	return &value
}
