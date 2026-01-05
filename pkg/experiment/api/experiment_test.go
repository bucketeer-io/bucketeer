// Copyright 2025 The Bucketeer Authors.
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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

func TestGetExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleUnassigned()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		setup         func(*experimentService)
		id            string
		environmentId string
		expectedErr   error
	}{
		{
			setup:         nil,
			id:            "",
			environmentId: "ns0",
			expectedErr:   statusExperimentIDRequired.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, v2es.ErrExperimentNotFound)
			},
			id:            "id-0",
			environmentId: "ns0",
			expectedErr:   statusExperimentNotFound.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Experiment{
					Experiment: &experimentproto.Experiment{Id: "id-1"},
				}, nil)
			},
			id:            "id-1",
			environmentId: "ns0",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		req := &experimentproto.GetExperimentRequest{Id: p.id, EnvironmentId: p.environmentId}
		_, err := service.GetExperiment(ctx, req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestListExperimentsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		orgRole     *accountproto.AccountV2_Role_Organization
		envRole     *accountproto.AccountV2_Role_Environment
		setup       func(*experimentService)
		req         *experimentproto.ListExperimentsRequest
		expectedErr error
	}{
		{
			desc:        "error: ErrPermissionDenied",
			orgRole:     toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole:     toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
			req:         &experimentproto.ListExperimentsRequest{FeatureId: "id-0", EnvironmentId: "ns0"},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "success",
			orgRole: toPtr(accountproto.AccountV2_Role_Organization_MEMBER),
			envRole: toPtr(accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().ListExperiments(
					gomock.Any(), gomock.Any(),
				).Return([]*experimentproto.Experiment{
					{Id: "id-1"},
				}, 0, int64(0), nil)
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperimentSummary(
					gomock.Any(), gomock.Any(),
				).Return(&v2es.ExperimentSummary{}, nil)

			},
			req:         &experimentproto.ListExperimentsRequest{FeatureId: "id-0", EnvironmentId: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController, nil, p.orgRole, p.envRole)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ListExperiments(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(s *experimentService)
		input       *experimentproto.CreateExperimentRequest
		expectedErr error
	}{
		{
			desc: "missing feature id",
			input: &experimentproto.CreateExperimentRequest{
				FeatureId:     "",
				GoalIds:       []string{"gid"},
				EnvironmentId: "ns0",
			},
			expectedErr: statusFeatureIDRequired.Err(),
		},
		{
			desc: "missing goal id",
			input: &experimentproto.CreateExperimentRequest{
				FeatureId:     "fid",
				GoalIds:       nil,
				EnvironmentId: "ns0",
			},
			expectedErr: statusGoalIDRequired.Err(),
		},
		{
			desc: "empty goal id",
			input: &experimentproto.CreateExperimentRequest{
				FeatureId:     "fid",
				GoalIds:       []string{""},
				EnvironmentId: "ns0",
			},
			expectedErr: statusGoalIDRequired.Err(),
		},
		{
			desc: "empty goal id",
			input: &experimentproto.CreateExperimentRequest{
				FeatureId:     "fid",
				GoalIds:       []string{"gid", ""},
				EnvironmentId: "ns0",
			},
			expectedErr: statusGoalIDRequired.Err(),
		},
		{
			desc: "period too long",
			input: &experimentproto.CreateExperimentRequest{
				FeatureId:     "fid",
				GoalIds:       []string{"gid0", "gid1"},
				StartAt:       1,
				StopAt:        30*24*60*60 + 2,
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentPeriodOutOfRange.Err(),
		},
		{
			desc: "success",
			setup: func(s *experimentService) {
				s.goalStorage.(*storagemock.MockGoalStorage).EXPECT().GetGoal(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Goal{
					Goal: &experimentproto.Goal{Id: "goalId", ConnectionType: experimentproto.Goal_EXPERIMENT},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().CreateExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			input: &experimentproto.CreateExperimentRequest{
				FeatureId:       "fid",
				GoalIds:         []string{"goalId"},
				Name:            "exp0",
				StartAt:         1,
				StopAt:          10,
				EnvironmentId:   "ns0",
				BaseVariationId: "variation-b-id",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.CreateExperiment(ctx, p.input)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestUpdateExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*experimentService)
		req         *experimentproto.UpdateExperimentRequest
		expectedErr error
	}{
		{
			desc:  "error id required",
			setup: nil,
			req: &experimentproto.UpdateExperimentRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentIDRequired.Err(),
		},
		{
			desc:  "period too long",
			setup: nil,
			req: &experimentproto.UpdateExperimentRequest{
				Id:            "id-1",
				StartAt:       wrapperspb.Int64(time.Now().Unix()),
				StopAt:        wrapperspb.Int64(time.Now().AddDate(0, 0, 31).Unix()),
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentPeriodOutOfRange.Err(),
		},
		{
			desc:  "invalid period input",
			setup: nil,
			req: &experimentproto.UpdateExperimentRequest{
				Id:            "id-1",
				StartAt:       wrapperspb.Int64(time.Now().Unix()),
				StopAt:        nil,
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentPeriodInvalid.Err(),
		},
		{
			desc: "experiment not found",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.UpdateExperimentRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
				Name:          wrapperspb.String("new-name"),
			},
			expectedErr: statusExperimentNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *experimentService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &experimentproto.UpdateExperimentRequest{
				Id:            "id-1",
				Name:          wrapperspb.String("new-name"),
				StartAt:       wrapperspb.Int64(time.Now().Unix()),
				StopAt:        wrapperspb.Int64(time.Now().AddDate(0, 0, 1).Unix()),
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.UpdateExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestStartExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*experimentService)
		req         *experimentproto.StartExperimentRequest
		expectedErr error
	}{
		{
			desc:  "error id required",
			setup: nil,
			req: &experimentproto.StartExperimentRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentIDRequired.Err(),
		},
		{
			desc:  "error no command",
			setup: nil,
			req: &experimentproto.StartExperimentRequest{
				Id:            "eid",
				EnvironmentId: "ns0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "error not found",
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, v2es.ErrExperimentNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.StartExperimentRequest{
				Id:            "noop",
				Command:       &experimentproto.StartExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Experiment{
					Experiment: &experimentproto.Experiment{Id: "id-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &experimentproto.StartExperimentRequest{
				Id:            "eid",
				Command:       &experimentproto.StartExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController, nil, nil, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.StartExperiment(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestFinishExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		desc        string
		setup       func(*experimentService)
		req         *experimentproto.FinishExperimentRequest
		expectedErr error
	}{
		{
			desc:  "error id required",
			setup: nil,
			req: &experimentproto.FinishExperimentRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentIDRequired.Err(),
		},
		{
			desc:  "error no command",
			setup: nil,
			req: &experimentproto.FinishExperimentRequest{
				Id:            "eid",
				EnvironmentId: "ns0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			desc: "error not found",
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, v2es.ErrExperimentNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.FinishExperimentRequest{
				Id:            "noop",
				Command:       &experimentproto.FinishExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Experiment{
					Experiment: &experimentproto.Experiment{Id: "id-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &experimentproto.FinishExperimentRequest{
				Id:            "eid",
				Command:       &experimentproto.FinishExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := createExperimentService(mockController, nil, nil, nil)
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.FinishExperiment(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestStopExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.StopExperimentRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.StopExperimentRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentIDRequired.Err(),
		},
		{
			setup: nil,
			req: &experimentproto.StopExperimentRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, v2es.ErrExperimentNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.StopExperimentRequest{
				Id:            "id-0",
				Command:       &experimentproto.StopExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentNotFound.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Experiment{
					Experiment: &experimentproto.Experiment{Id: "id-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &experimentproto.StopExperimentRequest{
				Id:            "id-1",
				Command:       &experimentproto.StopExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.StopExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestArchiveExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.ArchiveExperimentRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.ArchiveExperimentRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentIDRequired.Err(),
		},
		{
			setup: nil,
			req: &experimentproto.ArchiveExperimentRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, v2es.ErrExperimentNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.ArchiveExperimentRequest{
				Id:            "id-0",
				Command:       &experimentproto.ArchiveExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentNotFound.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Experiment{
					Experiment: &experimentproto.Experiment{Id: "id-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &experimentproto.ArchiveExperimentRequest{
				Id:            "id-1",
				Command:       &experimentproto.ArchiveExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.ArchiveExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestDeleteExperimentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenAndMetadata(metadata.MD{
		"accept-language": []string{"ja"},
	})

	patterns := []struct {
		setup       func(*experimentService)
		req         *experimentproto.DeleteExperimentRequest
		expectedErr error
	}{
		{
			setup: nil,
			req: &experimentproto.DeleteExperimentRequest{
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentIDRequired.Err(),
		},
		{
			setup: nil,
			req: &experimentproto.DeleteExperimentRequest{
				Id:            "id-0",
				EnvironmentId: "ns0",
			},
			expectedErr: statusNoCommand.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, v2es.ErrExperimentNotFound)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(v2es.ErrExperimentNotFound)
			},
			req: &experimentproto.DeleteExperimentRequest{
				Id:            "id-0",
				Command:       &experimentproto.DeleteExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: statusExperimentNotFound.Err(),
		},
		{
			setup: func(s *experimentService) {
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().GetExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Experiment{
					Experiment: &experimentproto.Experiment{Id: "id-1"},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.experimentStorage.(*storagemock.MockExperimentStorage).EXPECT().UpdateExperiment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			req: &experimentproto.DeleteExperimentRequest{
				Id:            "id-1",
				Command:       &experimentproto.DeleteExperimentCommand{},
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		service := createExperimentService(mockController, nil, nil, nil)
		if p.setup != nil {
			p.setup(service)
		}
		_, err := service.DeleteExperiment(ctx, p.req)
		assert.Equal(t, p.expectedErr, err)
	}
}

func TestExperimentPermissionDenied(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleUnassigned()
	// Use unassigned roles instead of default admin
	service := createExperimentService(
		mockController,
		toPtr("ns0"),
		toPtr(accountproto.AccountV2_Role_Organization_UNASSIGNED),
		toPtr(accountproto.AccountV2_Role_Environment_UNASSIGNED),
	)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc     string
		action   func(context.Context, *experimentService) error
		expected error
	}{
		{
			desc: "CreateExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.CreateExperiment(ctx, &experimentproto.CreateExperimentRequest{})
				return err
			},
			expected: statusPermissionDenied.Err(),
		},
		{
			desc: "UpdateExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.UpdateExperiment(ctx, &experimentproto.UpdateExperimentRequest{})
				return err
			},
			expected: statusPermissionDenied.Err(),
		},
		{
			desc: "StopExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.StopExperiment(ctx, &experimentproto.StopExperimentRequest{})
				return err
			},
			expected: statusPermissionDenied.Err(),
		},
		{
			desc: "DeleteExperiment",
			action: func(ctx context.Context, es *experimentService) error {
				_, err := es.DeleteExperiment(ctx, &experimentproto.DeleteExperimentRequest{})
				return err
			},
			expected: statusPermissionDenied.Err(),
		},
	}
	for _, p := range patterns {
		actual := p.action(ctx, service)
		assert.Equal(t, p.expected, actual, "%s", p.desc)
	}
}
