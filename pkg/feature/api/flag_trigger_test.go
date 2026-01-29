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
//

package api

import (
	"context"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCreateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)

	patterns := []struct {
		desc        string
		setup       func(service *FeatureService)
		input       *proto.CreateFlagTriggerRequest
		expectedErr error
	}{
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
			},
			input: &proto.CreateFlagTriggerRequest{
				EnvironmentId: "namespace",
				FeatureId:     "id-1",
				Type:          proto.FlagTrigger_Type_WEBHOOK,
				Action:        proto.FlagTrigger_Action_ON,
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error")).Err(),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().CreateFlagTrigger(
					ctx, gomock.Any(),
				).Return(nil)
			},
			input: &proto.CreateFlagTriggerRequest{
				EnvironmentId: "namespace",
				FeatureId:     "id-1",
				Type:          proto.FlagTrigger_Type_WEBHOOK,
				Action:        proto.FlagTrigger_Action_ON,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.CreateFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if resp != nil {
				assert.True(t, len(resp.FlagTrigger.Id) > 0)
				assert.Equal(t, p.input.FeatureId, resp.FlagTrigger.FeatureId)
				assert.Equal(t, p.input.Type, resp.FlagTrigger.Type)
				assert.Equal(t, p.input.Action, resp.FlagTrigger.Action)
				assert.True(t, len(resp.Url) > 0)
			}
		})
	}
}

func TestGetFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	baseId := "1"
	baseEnvironmentId := "ns0"

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(service *FeatureService)
		input          *proto.GetFlagTriggerRequest
		getExpectedErr func() error
	}{
		{
			desc: "Error Validate",
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			service: createFeatureServiceNew(mockController),
			setup:   nil,
			input:   &proto.GetFlagTriggerRequest{},
			getExpectedErr: func() error {
				return statusMissingTriggerID.Err()
			},
		},
		{
			desc: "Error Not Found",
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			service: createFeatureServiceNew(mockController),
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2fs.ErrFlagTriggerNotFound)
			},
			input: &proto.GetFlagTriggerRequest{Id: "1", EnvironmentId: "namespace"},
			getExpectedErr: func() error {
				return statusTriggerNotFound.Err()
			},
		},
		{
			desc: "Success",
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			service: createFeatureServiceNew(mockController),
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: &proto.FlagTrigger{
						Id:              baseId,
						FeatureId:       "featureId",
						EnvironmentId:   baseEnvironmentId,
						Type:            proto.FlagTrigger_Type_WEBHOOK,
						Action:          proto.FlagTrigger_Action_ON,
						Description:     "base",
						TriggerCount:    100,
						LastTriggeredAt: 500,
						Token:           "test-token",
						Disabled:        false,
						CreatedAt:       200,
						UpdatedAt:       300,
					},
				}, nil)
			},
			input: &proto.GetFlagTriggerRequest{Id: baseId, EnvironmentId: baseEnvironmentId},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc: "Success with Viewer Account",
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: &proto.FlagTrigger{
						Id:              "1",
						FeatureId:       "featureId",
						EnvironmentId:   "ns0",
						Type:            proto.FlagTrigger_Type_WEBHOOK,
						Action:          proto.FlagTrigger_Action_ON,
						Description:     "base",
						TriggerCount:    100,
						LastTriggeredAt: 500,
						Token:           "test-token",
						Disabled:        false,
						CreatedAt:       200,
						UpdatedAt:       300,
					},
				}, nil)
			},
			input: &proto.GetFlagTriggerRequest{Id: baseId, EnvironmentId: baseEnvironmentId},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc: "errPermissionDenied",
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:   func(s *FeatureService) {},
			input:   &proto.GetFlagTriggerRequest{Id: baseId, EnvironmentId: baseEnvironmentId},
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := p.service
			if p.setup != nil {
				p.setup(s)
			}
			ctx := p.context

			resp, err := s.GetFlagTrigger(ctx, p.input)
			assert.Equal(t, p.getExpectedErr(), err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestUpdateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)

	patterns := []struct {
		desc        string
		setup       func(service *FeatureService)
		input       *proto.UpdateFlagTriggerRequest
		expectedErr error
	}{
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
			},
			input: &proto.UpdateFlagTriggerRequest{
				Id:            "id",
				EnvironmentId: "namespace",
				Description:   wrapperspb.String("description"),
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error")).Err(),
		},
		{
			desc: "Success update description",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					ctx, gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: &proto.FlagTrigger{
						Id: "id",
					},
				}, nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().UpdateFlagTrigger(
					ctx, gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateFlagTriggerRequest{
				Id:            "id",
				EnvironmentId: "namespace",
				Description:   wrapperspb.String("description"),
			},
			expectedErr: nil,
		},
		{
			desc: "Success reset flag trigger",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					ctx, gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: &proto.FlagTrigger{
						Id:    "id",
						Token: "token",
					},
				}, nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().UpdateFlagTrigger(
					ctx, gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateFlagTriggerRequest{
				Id:            "id",
				EnvironmentId: "namespace",
				Reset_:        true,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.UpdateFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestDeleteFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	patterns := []struct {
		desc        string
		setup       func(service *FeatureService)
		input       *proto.DeleteFlagTriggerRequest
		expectedErr error
	}{
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error"))
			},
			input: &proto.DeleteFlagTriggerRequest{
				Id:            "id",
				EnvironmentId: "namespace",
			},
			expectedErr: api.NewGRPCStatus(pkgErr.NewErrorInternal(pkgErr.FeaturePackageName, "error")).Err(),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					ctx, gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: &proto.FlagTrigger{Id: "id"},
				}, nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().DeleteFlagTrigger(
					ctx, gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input:       &proto.DeleteFlagTriggerRequest{Id: "id", EnvironmentId: "namespace"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.DeleteFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestFlagTriggerWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)

	baseFlagTrigger := &proto.FlagTrigger{
		Id:              "1",
		FeatureId:       "featureId",
		EnvironmentId:   "namespace",
		Type:            proto.FlagTrigger_Type_WEBHOOK,
		Action:          proto.FlagTrigger_Action_ON,
		Description:     "base",
		TriggerCount:    100,
		LastTriggeredAt: 500,
		Token:           "test-token",
		Disabled:        false,
		CreatedAt:       200,
		UpdatedAt:       300,
	}

	patterns := []struct {
		desc        string
		setup       func(service *FeatureService)
		input       *proto.FlagTriggerWebhookRequest
		expectedErr error
	}{{
		desc:        "Error Invalid Argument",
		setup:       func(s *FeatureService) {},
		input:       &proto.FlagTriggerWebhookRequest{},
		expectedErr: statusSecretRequired.Err(),
	},
		{
			desc: "Error Not Found",
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTriggerByToken(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2fs.ErrFlagTriggerNotFound)
			},
			input:       &proto.FlagTriggerWebhookRequest{Token: "token"},
			expectedErr: statusTriggerNotFound.Err(),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTriggerByToken(
					gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: baseFlagTrigger,
				}, nil)

				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Feature{
					Feature: &proto.Feature{
						Id:      "id",
						Name:    "test feature",
						Version: 1,
						Enabled: true,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.domainPublisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().UpdateFlagTrigger(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input:       &proto.FlagTriggerWebhookRequest{Token: "token"},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.FlagTriggerWebhook(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListFlagTriggers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(service *FeatureService)
		input          *proto.ListFlagTriggersRequest
		expected       *proto.ListFlagTriggersResponse
		getExpectedErr func() error
	}{
		{
			desc:    "Error Validate",
			service: createFeatureServiceNew(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:    nil,
			input:    &proto.ListFlagTriggersRequest{},
			expected: nil,
			getExpectedErr: func() error {
				return statusMissingTriggerFeatureID.Err()
			},
		},
		{
			desc:    "Error Invalid Argument",
			service: createFeatureServiceNew(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:    nil,
			input:    &proto.ListFlagTriggersRequest{FeatureId: "1", Cursor: "XXX"},
			expected: nil,
			getExpectedErr: func() error {
				return statusInvalidCursor.Err()
			},
		},
		{
			desc:    "Success",
			service: createFeatureServiceNew(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().ListFlagTriggers(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.FlagTrigger{}, 0, int64(0), nil)
			},
			input:    &proto.ListFlagTriggersRequest{FeatureId: "1", PageSize: 2, Cursor: ""},
			expected: &proto.ListFlagTriggersResponse{FlagTriggers: []*proto.ListFlagTriggersResponse_FlagTriggerWithUrl{}, Cursor: "0"},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc:    "Success with Viewer Account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().ListFlagTriggers(
					gomock.Any(), gomock.Any(),
				).Return([]*proto.FlagTrigger{}, 0, int64(0), nil)
			},
			input:    &proto.ListFlagTriggersRequest{FeatureId: "1", PageSize: 2, Cursor: "", EnvironmentId: "ns0"},
			expected: &proto.ListFlagTriggersResponse{FlagTriggers: []*proto.ListFlagTriggersResponse_FlagTriggerWithUrl{}, Cursor: "0"},
			getExpectedErr: func() error {
				return nil
			},
		},
		{
			desc: "errPermissionDenied",
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:   func(s *FeatureService) {},
			input:   &proto.ListFlagTriggersRequest{FeatureId: "1", PageSize: 2, Cursor: "", EnvironmentId: "ns0"},
			getExpectedErr: func() error {
				return statusPermissionDenied.Err()
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := p.service
			if p.setup != nil {
				p.setup(s)
			}
			ctx := p.context

			actual, err := s.ListFlagTriggers(ctx, p.input)
			assert.Equal(t, p.getExpectedErr(), err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestFeatureServiceGenerateTriggerURL(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureService := createFeatureServiceNew(mockController)
	trigger, err := domain.NewFlagTrigger(
		"test",
		"test",
		proto.FlagTrigger_Type_WEBHOOK,
		proto.FlagTrigger_Action_ON,
		"test",
	)
	if err != nil {
		t.Errorf("NewFlagTrigger() error = %v", err)
	}
	err = trigger.GenerateToken()
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}
	t.Logf("GenerateToken() token = %v", trigger.Token)
	triggerURL := featureService.generateTriggerURL(context.Background(), trigger.Token, false)
	if triggerURL == "" {
		t.Errorf("generateTriggerURL() [full] triggerURL is empty")
	}
	t.Logf("generateTriggerURL() [full] triggerURL = %v", triggerURL)
	triggerURL = featureService.generateTriggerURL(context.Background(), trigger.Token, true)
	if triggerURL == "" {
		t.Errorf("generateTriggerURL() [masked] triggerURL is empty")
	}
	t.Logf("generateTriggerURL() [masked] triggerURL = %v", triggerURL)
}
