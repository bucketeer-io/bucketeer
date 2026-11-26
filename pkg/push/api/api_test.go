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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	"github.com/bucketeer-io/bucketeer/v2/pkg/authz"
	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	v2ps "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
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

// errPublish is returned by the publisher mock to simulate a publish failure
// after the transaction has committed.
var errPublish = bkterr.NewErrorInternal(bkterr.PushPackageName, "publish error")

// errStorage is returned by storage mocks to simulate a failed write inside
// the transaction.
var errStorage = bkterr.NewErrorInternal(bkterr.PushPackageName, "storage error")

// errCommit simulates a transaction that fails to commit after the callback succeeded.
var errCommit = bkterr.NewErrorInternal(bkterr.PushPackageName, "commit error")

func TestNewPushService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	dbClient := databasemock.NewMockClient(mockController)
	pushStorageMock := storagemock.NewMockPushStorage(mockController)
	featureClientMock := featureclientmock.NewMockClient(mockController)
	experimentClientMock := experimentclientmock.NewMockClient(mockController)
	pm := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewPushService(
		dbClient,
		pushStorageMock,
		featureClientMock,
		experimentClientMock,
		pm,
		WithLogger(logger),
	)
	assert.IsType(t, &PushService{}, s)
}

func TestCreatePushMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithEditor(t, newEditor(nil, nil, nil))

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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
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
			desc: "err: storage write fails, event is not published",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().CreatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errStorage)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
				Name:              "name-1",
			},
			expectedErr: api.NewGRPCStatus(errStorage).Err(),
		},
		{
			desc: "err: commit fails, event is not published",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					if err := fn(ctx); err != nil {
						return err
					}
					return errCommit
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().CreatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
				Name:              "name-1",
			},
			expectedErr: api.NewGRPCStatus(errCommit).Err(),
		},
		{
			desc: "err: publish fails after commit",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().CreatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errPublish).Times(1)
			},
			req: &pushproto.CreatePushRequest{
				EnvironmentId:     "ns0",
				FcmServiceAccount: fcmServiceAccountDummy,
				Tags:              []string{"tag-0"},
				Name:              "name-1",
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc: "success: with tags",
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().ListPushes(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.Push{}, 0, int64(0), nil)
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
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

	ctx := createContextWithEditor(t, newEditor(nil, nil, nil))

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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
				}).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.UpdatePushRequest{
				Id:   "key-1",
				Name: wrapperspb.String("push-0"),
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc: "err: storage write fails, event is not published",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errStorage)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
				Name:          wrapperspb.String("name-1"),
			},
			expectedErr: api.NewGRPCStatus(errStorage).Err(),
		},
		{
			desc: "err: commit fails, event is not published",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					if err := fn(ctx); err != nil {
						return err
					}
					return errCommit
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
				Name:          wrapperspb.String("name-1"),
			},
			expectedErr: api.NewGRPCStatus(errCommit).Err(),
		},
		{
			desc: "err: publish fails after commit",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errPublish).Times(1)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
				Name:          wrapperspb.String("name-1"),
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc: "success: no rows affected, event is not published",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().UpdatePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushUnexpectedAffectedRows)
			},
			req: &pushproto.UpdatePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
				Name:          wrapperspb.String("name-1"),
			},
			expectedErr: nil,
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
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

	ctx := createContextWithEditor(t, newEditor(nil, nil, nil))

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

	ctx := createContextWithEditor(t, newEditor(nil, nil, nil))

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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
				}).Return(v2ps.ErrPushNotFound)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc: "err: storage write fails, event is not published",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().DeletePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errStorage)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
			},
			expectedErr: api.NewGRPCStatus(errStorage).Err(),
		},
		{
			desc: "err: commit fails, event is not published",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					if err := fn(ctx); err != nil {
						return err
					}
					return errCommit
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().DeletePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
			},
			expectedErr: api.NewGRPCStatus(errCommit).Err(),
		},
		{
			desc: "err: publish fails after commit",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().DeletePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errPublish).Times(1)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
			},
			expectedErr: statusInternal.Err(),
		},
		{
			desc: "success: no rows affected, event is not published",
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
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().DeletePush(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2ps.ErrPushUnexpectedAffectedRows)
			},
			req: &pushproto.DeletePushRequest{
				EnvironmentId: "ns0",
				Id:            "key-0",
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
						Id: "key-0",
					},
				}, nil)
				s.dbClient.(*databasemock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context) error) {
					_ = fn(ctx)
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
			ctx := createContextWithEditor(t, newEditor(nil, p.orgRole, p.envRole))
			s := newPushService(mockController)
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

	storedPush := func() *domain.Push {
		return &domain.Push{
			Push: &proto.Push{
				Id:                "key-1",
				Name:              "push-1",
				FcmServiceAccount: string(fcmServiceAccountDummy),
				Tags:              []string{"tag-1"},
				EnvironmentId:     "ns0",
			},
		}
	}
	// The service account must never be returned to the client.
	expectedPush := &proto.Push{
		Id:            "key-1",
		Name:          "push-1",
		Tags:          []string{"tag-1"},
		EnvironmentId: "ns0",
	}

	patterns := []struct {
		desc        string
		envID       *string
		orgRole     *accountproto.AccountV2_Role_Organization
		envRole     *accountproto.AccountV2_Role_Environment
		setup       func(*PushService)
		req         *pushproto.GetPushRequest
		expected    *pushproto.GetPushResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			orgRole:     toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:     toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			req:         &pushproto.GetPushRequest{EnvironmentId: "ns0"},
			expectedErr: statusIDRequired.Err(),
		},
		{
			desc:    "err: ErrNotFound",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), "key-1", "ns0",
				).Return(nil, v2ps.ErrPushNotFound)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: statusNotFound.Err(),
		},
		{
			desc:    "err: ErrInternal",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), "key-1", "ns0",
				).Return(nil, pkgErr.NewErrorInternal(pkgErr.PushPackageName, "test"))
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.PushPackageName, "test")).Err(),
		},
		{
			desc:    "success: viewer of the environment",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), "key-1", "ns0",
				).Return(storedPush(), nil)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expected: &pushproto.GetPushResponse{Push: expectedPush},
		},
		{
			desc:    "success: organization admin without an environment role",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_ADMIN),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), "key-1", "ns0",
				).Return(storedPush(), nil)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expected: &pushproto.GetPushResponse{Push: expectedPush},
		},
		{
			desc:    "success: system admin without an environment role",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup: func(s *PushService) {
				s.pushStorage.(*storagemock.MockPushStorage).EXPECT().GetPush(
					gomock.Any(), "key-1", "ns0",
				).Return(storedPush(), nil)
			},
			req: &pushproto.GetPushRequest{
				EnvironmentId: "ns0",
				Id:            "key-1",
			},
			expected: &pushproto.GetPushResponse{Push: expectedPush},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := createContextWithEditor(t, newEditor(p.envID, p.orgRole, p.envRole))
			service := newPushService(mockController)
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.GetPush(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newPushServiceWithMock(t *testing.T, c *gomock.Controller) *PushService {
	t.Helper()
	return &PushService{
		dbClient:         databasemock.NewMockClient(c),
		pushStorage:      storagemock.NewMockPushStorage(c),
		featureClient:    featureclientmock.NewMockClient(c),
		experimentClient: experimentclientmock.NewMockClient(c),
		publisher:        publishermock.NewMockPublisher(c),
		logger:           zap.NewNop(),
	}
}

func newPushService(c *gomock.Controller) *PushService {
	p := publishermock.NewMockPublisher(c)
	p.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	return &PushService{
		dbClient:         databasemock.NewMockClient(c),
		featureClient:    featureclientmock.NewMockClient(c),
		pushStorage:      storagemock.NewMockPushStorage(c),
		experimentClient: experimentclientmock.NewMockClient(c),
		publisher:        publishermock.NewMockPublisher(c),
		logger:           zap.NewNop(),
	}
}

// newEditor is the caller the authz interceptor would have set: an admin
// editor of "ns0" unless the roles or the environment are specified.
func newEditor(
	specifiedEnvironmentId *string,
	specifiedOrgRole *accountproto.AccountV2_Role_Organization,
	specifiedEnvRole *accountproto.AccountV2_Role_Environment,
) *eventproto.Editor {
	or := accountproto.AccountV2_Role_Organization_ADMIN
	er := accountproto.AccountV2_Role_Environment_EDITOR
	envId := "ns0"
	if specifiedOrgRole != nil {
		or = *specifiedOrgRole
	}
	if specifiedEnvRole != nil {
		er = *specifiedEnvRole
	}
	if specifiedEnvironmentId != nil {
		envId = *specifiedEnvironmentId
	}
	return &eventproto.Editor{
		Email:            "email",
		OrganizationRole: or,
		EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
			{EnvironmentId: envId, Role: er},
		},
	}
}

// createContextWithEditor is a context after the authz interceptor allowed the request.
func createContextWithEditor(t *testing.T, editor *eventproto.Editor) context.Context {
	t.Helper()
	return authz.WithEditor(context.TODO(), editor)
}

// convert to pointer
func toPtr[T any](value T) *T {
	return &value
}
