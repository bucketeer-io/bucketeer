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
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	authclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/auth/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2ao "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	mockAutoOpsStorage "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2/mock"
	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	experimentclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client/mock"
	featureclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client/mock"
	mockFeatureStorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	mockOpsCountStorage "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

func TestNewAutoOpsService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mysqlClientMock := mysqlmock.NewMockClient(mockController)
	featureClientMock := featureclientmock.NewMockClient(mockController)
	experimentClientMock := experimentclientmock.NewMockClient(mockController)
	accountClientMock := accountclientmock.NewMockClient(mockController)
	authClientMock := authclientmock.NewMockClient(mockController)
	p := publishermock.NewMockPublisher(mockController)
	logger := zap.NewNop()
	s := NewAutoOpsService(
		mysqlClientMock,
		featureClientMock,
		experimentClientMock,
		accountClientMock,
		authClientMock,
		p,
		WithLogger(logger),
	)
	assert.IsType(t, &AutoOpsService{}, s)
}

func TestCreateAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.CreateAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc: "err: ErrFeatureIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				EnvironmentId: "env-id",
			},
			expectedErr: statusFeatureIDRequired.Err(),
		},
		{
			desc: "err: ErrClauseRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
			},
			expectedErr: statusClauseRequired.Err(),
		},
		{
			desc: "err: ErrIncompatibleOpsType",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_TYPE_UNKNOWN,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: statusIncompatibleOpsType.Err(),
		},
		{
			desc: "err: ErrOpsEventRateClauseVariationIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "",
						GoalId:          "gid1",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: statusOpsEventRateClauseVariationIDRequired.Err(),
		},
		{
			desc: "err: ErrOpsEventRateClauseGoalIDRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: statusOpsEventRateClauseGoalIDRequired.Err(),
		},
		{
			desc: "err: ErrOpsEventRateClauseMinCountRequired",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        0,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: statusOpsEventRateClauseMinCountRequired.Err(),
		},
		{
			desc: "err: ErrOpsEventRateClauseInvalidThredshold: less",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: -0.1,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: statusOpsEventRateClauseInvalidThredshold.Err(),
		},
		{
			desc: "err: ErrOpsEventRateClauseInvalidThredshold: greater",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 1.1,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			expectedErr: statusOpsEventRateClauseInvalidThredshold.Err(),
		},
		{
			desc: "err: ErrDatetimeClauseInvalidTime",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: 0},
				},
			},
			expectedErr: statusDatetimeClauseInvalidTime.Err(),
		},
		{
			desc: "err: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
			},
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "err: ErrDatetimeClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId:       "fid",
				OpsType:         autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: nil,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: statusClauseRequiredForDateTime.Err(),
		},
		{
			desc: "err: ErrDatetimeClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: statusIncompatibleOpsType.Err(),
		},
		{
			desc: "err: ErrOpsEventRateClauseMustSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
				OpsEventRateClauses: nil,
			},
			expectedErr: statusClauseRequiredForEventRate.Err(),
		},
		{
			desc: "err: ErrDatetimeClauseMustNotBeSpecified",
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: statusIncompatibleOpsType.Err(),
		},
		{
			desc: "err: internal error",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(nil, bkterr.NewErrorInternal(bkterr.AutoopsPackageName, "error"))
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: api.NewGRPCStatus(bkterr.NewErrorInternal(bkterr.AutoopsPackageName, "error")).Err(),
		},
		{
			desc: "success event rate",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(
					gomock.Any(), gomock.Any(),
				).Return(&experimentproto.GetGoalResponse{}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().CreateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_EVENT_RATE,
				OpsEventRateClauses: []*autoopsproto.OpsEventRateClause{
					{
						VariationId:     "vid",
						GoalId:          "gid",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success schedule",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().CreateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.CreateAutoOpsRuleRequest{
				FeatureId: "fid",
				OpsType:   autoopsproto.OpsType_SCHEDULE,
				DatetimeClauses: []*autoopsproto.DatetimeClause{
					{Time: time.Now().AddDate(0, 0, 1).Unix(), ActionType: autoopsproto.ActionType_ENABLE},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.CreateAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	selfMatchTime := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.UpdateAutoOpsRuleRequest
		expected    *autoopsproto.UpdateAutoOpsRuleResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.UpdateAutoOpsRuleRequest{},
			expected:    nil,
			expectedErr: statusAutoOpsRuleIDRequired.Err(),
		},
		{
			desc: "err: ErrOpsEventRateClauseRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:                        "aid1",
				OpsEventRateClauseChanges: []*autoopsproto.OpsEventRateClauseChange{{}},
			},
			expected:    nil,
			expectedErr: statusOpsEventRateClauseRequired.Err(),
		},
		{
			desc: "err: DeleteClause ErrClauseIdRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				OpsEventRateClauseChanges: []*autoopsproto.OpsEventRateClauseChange{{
					ChangeType: autoopsproto.ChangeType_DELETE,
				}},
			},
			expected:    nil,
			expectedErr: statusClauseIDRequired.Err(),
		},
		{
			desc: "err: ChangeOpsEventRateClauseCommand: ErrOpsEventRateClauseRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				OpsEventRateClauseChanges: []*autoopsproto.OpsEventRateClauseChange{{
					Id: "aid",
				}},
			},
			expected:    nil,
			expectedErr: statusOpsEventRateClauseRequired.Err(),
		},
		{
			desc: "err: ErrDatetimeClauseRequired",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{{
					ChangeType: autoopsproto.ChangeType_UPDATE,
				}},
			},
			expected:    nil,
			expectedErr: statusDatetimeClauseRequired.Err(),
		},
		{
			desc: "err: ChangeDatetimeClause: ErrDatetimeClauseInvalidTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{{
					Id:         "aid",
					Clause:     &autoopsproto.DatetimeClause{Time: 0},
					ChangeType: autoopsproto.ChangeType_UPDATE,
				}},
			},
			expected:    nil,
			expectedErr: statusDatetimeClauseInvalidTime.Err(),
		},
		{
			desc: "err: ChangeDatetimeClause: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{
					{
						Id: "aid",
						Clause: &autoopsproto.DatetimeClause{
							Time:       time.Now().AddDate(0, 0, 1).Unix(),
							ActionType: autoopsproto.ActionType_ENABLE,
						},
						ChangeType: autoopsproto.ChangeType_UPDATE,
					},
					{
						Id: "aid2",
						Clause: &autoopsproto.DatetimeClause{
							Time:       time.Now().AddDate(0, 0, 1).Unix(),
							ActionType: autoopsproto.ActionType_ENABLE,
						},
						ChangeType: autoopsproto.ChangeType_UPDATE,
					},
				},
			},
			expected:    nil,
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "err: AddDatetimeClause: ErrDatetimeClauseInvalidTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{
					{
						Clause:     &autoopsproto.DatetimeClause{Time: 0, ActionType: autoopsproto.ActionType_DISABLE},
						ChangeType: autoopsproto.ChangeType_CREATE,
					},
				},
			},
			expected:    nil,
			expectedErr: statusDatetimeClauseInvalidTime.Err(),
		},
		{
			desc: "err: AddDatetimeClauses: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{
					{
						Clause: &autoopsproto.DatetimeClause{
							Time:       time.Now().AddDate(0, 0, 1).Unix(),
							ActionType: autoopsproto.ActionType_DISABLE,
						},
						ChangeType: autoopsproto.ChangeType_CREATE,
					},
					{
						Clause: &autoopsproto.DatetimeClause{
							Time:       time.Now().AddDate(0, 0, 1).Unix(),
							ActionType: autoopsproto.ActionType_DISABLE,
						},
						ChangeType: autoopsproto.ChangeType_CREATE,
					},
				},
			},
			expected:    nil,
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "err: AddAndUpdateDatetimeClause: ErrDatetimeClauseDuplicateTime",
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id: "aid1",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{
					{
						Clause: &autoopsproto.DatetimeClause{
							Time:       time.Now().AddDate(0, 0, 1).Unix(),
							ActionType: autoopsproto.ActionType_DISABLE,
						},
						ChangeType: autoopsproto.ChangeType_CREATE,
					},
					{
						Id: "aid",
						Clause: &autoopsproto.DatetimeClause{
							Time:       time.Now().AddDate(0, 0, 1).Unix(),
							ActionType: autoopsproto.ActionType_DISABLE,
						},
						ChangeType: autoopsproto.ChangeType_UPDATE,
					},
				},
			},
			expected:    nil,
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.All(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id:            "aid1",
						OpsType:       autoopsproto.OpsType_SCHEDULE,
						AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING,
						Deleted:       false,
						Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
							{Id: "cid2", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil).AnyTimes()
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().UpdateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{
					{
						Clause: &autoopsproto.DatetimeClause{
							ActionType: autoopsproto.ActionType_ENABLE,
							Time:       time.Now().AddDate(0, 0, 1).Unix(),
						},
						ChangeType: autoopsproto.ChangeType_CREATE,
					},
					{
						Id: "cid2",
						Clause: &autoopsproto.DatetimeClause{
							ActionType: autoopsproto.ActionType_DISABLE,
							Time:       time.Now().AddDate(0, 0, 2).Unix(),
						},
						ChangeType: autoopsproto.ChangeType_UPDATE,
					},
					{
						Id:         "cid",
						ChangeType: autoopsproto.ChangeType_DELETE,
					},
				},
			},
			expected:    &autoopsproto.UpdateAutoOpsRuleResponse{},
			expectedErr: nil,
		},
		{
			desc: "success: UPDATE clause with same time does not trigger self-duplicate",
			setup: func(s *AutoOpsService) {
				existingClause, _ := anypb.New(&autoopsproto.DatetimeClause{
					Time:       selfMatchTime,
					ActionType: autoopsproto.ActionType_ENABLE,
				})
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.All(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id:            "aid1",
						OpsType:       autoopsproto.OpsType_SCHEDULE,
						AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING,
						Deleted:       false,
						Clauses: []*autoopsproto.Clause{
							{Id: "cid1", ActionType: autoopsproto.ActionType_ENABLE, Clause: existingClause},
						}},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil).AnyTimes()
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().UpdateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.UpdateAutoOpsRuleRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				DatetimeClauseChanges: []*autoopsproto.DatetimeClauseChange{
					{
						Id: "cid1",
						Clause: &autoopsproto.DatetimeClause{
							ActionType: autoopsproto.ActionType_DISABLE,
							Time:       selfMatchTime,
						},
						ChangeType: autoopsproto.ChangeType_UPDATE,
					},
				},
			},
			expected:    &autoopsproto.UpdateAutoOpsRuleResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.UpdateAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestStopAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.StopAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.StopAutoOpsRuleRequest{},
			expectedErr: statusAutoOpsRuleIDRequired.Err(),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().UpdateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.StopAutoOpsRuleRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.StopAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.DeleteAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.DeleteAutoOpsRuleRequest{},
			expectedErr: statusAutoOpsRuleIDRequired.Err(),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					_ = fn(ctx, nil)
				}).Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().Publish(
					gomock.Any(), gomock.Any(),
				).Return(nil)
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().UpdateAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopsproto.DeleteAutoOpsRuleRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.DeleteAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleUnassigned(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc        string
		service     *AutoOpsService
		setup       func(*AutoOpsService)
		req         *autoopsproto.GetAutoOpsRuleRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			service:     createAutoOpsService(mockController),
			req:         &autoopsproto.GetAutoOpsRuleRequest{EnvironmentId: "ns0"},
			expectedErr: statusAutoOpsRuleIDRequired.Err(),
		},
		{
			desc:    "err: ErrNotFound",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ao.ErrAutoOpsRuleNotFound)
			},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "wrongid", EnvironmentId: "ns0"},
			expectedErr: statusAutoOpsRuleNotFound.Err(),
		},
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *AutoOpsService) {},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "aid1", EnvironmentId: "ns0"},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "success",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
			},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "aid1", EnvironmentId: "ns0"},
			expectedErr: nil,
		},
		{
			desc:    "success with view account",
			service: createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "cid", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil)
			},
			req:         &autoopsproto.GetAutoOpsRuleRequest{Id: "aid1", EnvironmentId: "ns0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := p.service
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.GetAutoOpsRule(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAutoOpsRulesMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleUnassigned(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc        string
		service     *AutoOpsService
		setup       func(*AutoOpsService)
		req         *autoopsproto.ListAutoOpsRulesRequest
		expectedErr error
	}{
		{
			desc:    "success",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().ListAutoOpsRules(
					gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.AutoOpsRule{}, 0, nil)
			},
			req:         &autoopsproto.ListAutoOpsRulesRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *AutoOpsService) {},
			req:         &autoopsproto.ListAutoOpsRulesRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "success with viewer",
			service: createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().ListAutoOpsRules(
					gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.AutoOpsRule{}, 0, nil)
			},
			req:         &autoopsproto.ListAutoOpsRulesRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ListAutoOpsRules(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListOpsCountsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleUnassigned(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc        string
		service     *AutoOpsService
		setup       func(*AutoOpsService)
		req         *autoopsproto.ListOpsCountsRequest
		expectedErr error
	}{
		{
			desc:    "success",
			service: createAutoOpsService(mockController),
			setup: func(s *AutoOpsService) {
				s.opsCountStorage.(*mockOpsCountStorage.MockOpsCountStorage).EXPECT().ListOpsCounts(
					gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.OpsCount{}, 0, nil)
			},
			req:         &autoopsproto.ListOpsCountsRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
		{
			desc:        "errPermissionDenied",
			service:     createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			setup:       func(s *AutoOpsService) {},
			req:         &autoopsproto.ListOpsCountsRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: statusPermissionDenied.Err(),
		},
		{
			desc:    "success with view ",
			service: createServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			setup: func(s *AutoOpsService) {
				s.opsCountStorage.(*mockOpsCountStorage.MockOpsCountStorage).EXPECT().ListOpsCounts(
					gomock.Any(), gomock.Any(),
				).Return([]*autoopsproto.OpsCount{}, 0, nil)
			},
			req:         &autoopsproto.ListOpsCountsRequest{EnvironmentId: "ns0", Cursor: ""},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			service := p.service
			if p.setup != nil {
				p.setup(service)
			}
			_, err := service.ListOpsCounts(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestExecuteAutoOpsRuleMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopsproto.ExecuteAutoOpsRequest
		expectedErr error
	}{
		{
			desc:        "err: ErrIDRequired",
			req:         &autoopsproto.ExecuteAutoOpsRequest{},
			expectedErr: statusAutoOpsRuleIDRequired.Err(),
		},
		{
			desc: "err: ErrNoExecuteAutoOpsRuleCommand_ClauseId",
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				ClauseId:      "",
			},
			expectedErr: statusClauseIDRequired.Err(),
		},
		{
			desc: "err: ErrNotFound",
			setup: func(s *AutoOpsService) {
				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2ao.ErrAutoOpsRuleNotFound)
			},
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				ClauseId:      "id",
			},
			expectedErr: statusAutoOpsRuleNotFound.Err(),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(nil)

				s.autoOpsStorage.(*mockAutoOpsStorage.MockAutoOpsRuleStorage).EXPECT().GetAutoOpsRule(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.AutoOpsRule{
					AutoOpsRule: &autoopsproto.AutoOpsRule{
						Id: "aid1", OpsType: autoopsproto.OpsType_SCHEDULE, AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING, Deleted: false, Clauses: []*autoopsproto.Clause{
							{Id: "testClauseId", ActionType: autoopsproto.ActionType_ENABLE, Clause: &anypb.Any{}},
						}},
				}, nil).AnyTimes()
			},
			req: &autoopsproto.ExecuteAutoOpsRequest{
				Id:            "aid1",
				EnvironmentId: "ns0",
				ClauseId:      "testClauseId",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			_, err := s.ExecuteAutoOps(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestExistGoal(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		goalID      string
		expected    bool
		expectedErr error
	}{
		{
			desc: "not found",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(gomock.Any(), gomock.Any()).Return(nil, storage.ErrKeyNotFound)
			},
			goalID:      "gid-0",
			expected:    false,
			expectedErr: nil,
		},
		{
			desc: "fails",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))
			},
			goalID:      "gid-0",
			expected:    false,
			expectedErr: errors.New("test"),
		},
		{
			desc: "exists",
			setup: func(s *AutoOpsService) {
				s.experimentClient.(*experimentclientmock.MockClient).EXPECT().GetGoal(gomock.Any(), gomock.Any()).Return(&experimentproto.GetGoalResponse{}, nil)
			},
			goalID:      "gid-0",
			expected:    true,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.existGoal(context.Background(), "ns-0", p.goalID)
			assert.Equal(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func createAutoOpsService(c *gomock.Controller) *AutoOpsService {
	mysqlClientMock := mysqlmock.NewMockClient(c)
	featureClientMock := featureclientmock.NewMockClient(c)
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: accountproto.AccountV2_Role_Organization_ADMIN,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          accountproto.AccountV2_Role_Environment_EDITOR,
				},
				{
					EnvironmentId: "",
					Role:          accountproto.AccountV2_Role_Environment_EDITOR,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	experimentClientMock := experimentclientmock.NewMockClient(c)
	authClientMock := authclientmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &AutoOpsService{
		mysqlClient:      mysqlClientMock,
		featureStorage:   mockFeatureStorage.NewMockFeatureStorage(c),
		autoOpsStorage:   mockAutoOpsStorage.NewMockAutoOpsRuleStorage(c),
		prStorage:        mockAutoOpsStorage.NewMockProgressiveRolloutStorage(c),
		opsCountStorage:  mockOpsCountStorage.NewMockOpsCountStorage(c),
		featureClient:    featureClientMock,
		experimentClient: experimentClientMock,
		accountClient:    accountClientMock,
		authClient:       authClientMock,
		publisher:        p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

func createServiceWithGetAccountByEnvironmentMock(c *gomock.Controller, ro accountproto.AccountV2_Role_Organization, re accountproto.AccountV2_Role_Environment) *AutoOpsService {
	mysqlClientMock := mysqlmock.NewMockClient(c)
	featureClientMock := featureclientmock.NewMockClient(c)
	accountClientMock := accountclientmock.NewMockClient(c)
	ar := &accountproto.GetAccountV2ByEnvironmentIDResponse{
		Account: &accountproto.AccountV2{
			Email:            "email",
			OrganizationRole: ro,
			EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
				{
					EnvironmentId: "ns0",
					Role:          re,
				},
			},
		},
	}
	accountClientMock.EXPECT().GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).Return(ar, nil).AnyTimes()
	experimentClientMock := experimentclientmock.NewMockClient(c)
	authClientMock := authclientmock.NewMockClient(c)
	p := publishermock.NewMockPublisher(c)
	logger := zap.NewNop()
	return &AutoOpsService{
		mysqlClient:      mysqlClientMock,
		autoOpsStorage:   mockAutoOpsStorage.NewMockAutoOpsRuleStorage(c),
		prStorage:        mockAutoOpsStorage.NewMockProgressiveRolloutStorage(c),
		opsCountStorage:  mockOpsCountStorage.NewMockOpsCountStorage(c),
		featureClient:    featureClientMock,
		experimentClient: experimentClientMock,
		accountClient:    accountClientMock,
		authClient:       authClientMock,
		publisher:        p,
		opts: &options{
			logger: zap.NewNop(),
		},
		logger: logger,
	}
}

func createContextWithTokenRoleUnassigned(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}

func createContextWithTokenRoleOwner(t *testing.T) context.Context {
	t.Helper()
	token := &token.AccessToken{
		Issuer:        "issuer",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: true,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.AccessTokenKey, token)
}

func TestValidateDatetimeClause(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		clause      *autoopsproto.DatetimeClause
		expectedErr error
	}{
		{
			desc: "success: one-time with future time",
			clause: &autoopsproto.DatetimeClause{
				Time:       time.Now().Add(24 * time.Hour).Unix(),
				ActionType: autoopsproto.ActionType_ENABLE,
			},
			expectedErr: nil,
		},
		{
			desc: "err: one-time with past time",
			clause: &autoopsproto.DatetimeClause{
				Time:       time.Now().Add(-1 * time.Hour).Unix(),
				ActionType: autoopsproto.ActionType_ENABLE,
			},
			expectedErr: statusDatetimeClauseInvalidTime.Err(),
		},
		{
			desc: "err: unknown action type",
			clause: &autoopsproto.DatetimeClause{
				Time:       time.Now().Add(24 * time.Hour).Unix(),
				ActionType: autoopsproto.ActionType_UNKNOWN,
			},
			expectedErr: statusIncompatibleOpsType.Err(),
		},
		{
			desc: "success: recurring with valid time-of-day",
			clause: &autoopsproto.DatetimeClause{
				Time:       36000, // 10:00 AM
				ActionType: autoopsproto.ActionType_ENABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{1},
					Timezone:   "Asia/Tokyo",
					StartDate:  time.Now().Add(24 * time.Hour).Unix(),
				},
			},
			expectedErr: nil,
		},
		{
			desc: "err: recurring with invalid time-of-day (>= 86400)",
			clause: &autoopsproto.DatetimeClause{
				Time:       86400,
				ActionType: autoopsproto.ActionType_ENABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{1},
					Timezone:   "Asia/Tokyo",
					StartDate:  time.Now().Add(24 * time.Hour).Unix(),
				},
			},
			expectedErr: statusDatetimeClauseInvalidTimeOfDay.Err(),
		},
		{
			desc: "err: recurring with negative time",
			clause: &autoopsproto.DatetimeClause{
				Time:       -1,
				ActionType: autoopsproto.ActionType_ENABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{1},
					Timezone:   "Asia/Tokyo",
					StartDate:  time.Now().Add(24 * time.Hour).Unix(),
				},
			},
			expectedErr: statusDatetimeClauseInvalidTimeOfDay.Err(),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			err := s.validateDatetimeClause(p.clause)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateRecurrenceRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		recurrence  *autoopsproto.RecurrenceRule
		expectedErr error
	}{
		{
			desc:        "success: nil recurrence",
			recurrence:  nil,
			expectedErr: nil,
		},
		{
			desc: "success: ONCE frequency",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_ONCE,
			},
			expectedErr: nil,
		},
		{
			desc: "err: FREQUENCY_UNSPECIFIED",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_FREQUENCY_UNSPECIFIED,
			},
			expectedErr: statusInvalidRecurrenceFrequency.Err(),
		},
		{
			desc: "err: missing timezone",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_DAILY,
				StartDate: time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusRecurrenceTimezoneRequired.Err(),
		},
		{
			desc: "err: invalid timezone",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_DAILY,
				Timezone:  "Invalid/Timezone",
				StartDate: time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusInvalidRecurrenceTimezone.Err(),
		},
		{
			desc: "err: missing start date",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_DAILY,
				Timezone:  "UTC",
			},
			expectedErr: statusRecurrenceStartDateRequired.Err(),
		},
		{
			desc: "err: end date before start date",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_DAILY,
				Timezone:  "UTC",
				StartDate: time.Now().Add(48 * time.Hour).Unix(),
				EndDate:   time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusRecurrenceEndDateMustBeAfterStart.Err(),
		},
		{
			desc: "err: negative max occurrences",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency:      autoopsproto.RecurrenceRule_DAILY,
				Timezone:       "UTC",
				StartDate:      time.Now().Add(24 * time.Hour).Unix(),
				MaxOccurrences: -1,
			},
			expectedErr: statusRecurrenceMaxOccurrencesInvalid.Err(),
		},
		{
			desc: "success: valid weekly",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
				DaysOfWeek: []int32{1, 5},
				Timezone:   "Asia/Tokyo",
				StartDate:  time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: nil,
		},
		{
			desc: "err: weekly without days",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_WEEKLY,
				Timezone:  "UTC",
				StartDate: time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusRecurrenceDaysOfWeekRequired.Err(),
		},
		{
			desc: "err: weekly with invalid day (7)",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
				DaysOfWeek: []int32{1, 7},
				Timezone:   "UTC",
				StartDate:  time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusRecurrenceDaysOfWeekInvalid.Err(),
		},
		{
			desc: "success: valid monthly",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
				DayOfMonth: 15,
				Timezone:   "UTC",
				StartDate:  time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: nil,
		},
		{
			desc: "err: monthly with invalid day (0)",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
				DayOfMonth: 0,
				Timezone:   "UTC",
				StartDate:  time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusRecurrenceDayOfMonthInvalid.Err(),
		},
		{
			desc: "err: monthly with invalid day (32)",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
				DayOfMonth: 32,
				Timezone:   "UTC",
				StartDate:  time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusRecurrenceDayOfMonthInvalid.Err(),
		},
		{
			desc: "success: valid daily",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_DAILY,
				Timezone:  "America/New_York",
				StartDate: time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: nil,
		},
		{
			desc: "err: unsupported frequency value",
			recurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_Frequency(99),
				Timezone:  "UTC",
				StartDate: time.Now().Add(24 * time.Hour).Unix(),
			},
			expectedErr: statusInvalidRecurrenceFrequency.Err(),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			err := s.validateRecurrenceRule(p.recurrence)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestValidateDatetimeClauses_RecurringDuplicates(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		clauses     []*autoopsproto.DatetimeClause
		expectedErr error
	}{
		{
			desc: "success: two recurring clauses same time, different days",
			clauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_DISABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{5},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "err: two recurring clauses same time and same days",
			clauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
			},
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "err: two one-time clauses same time",
			clauses: func() []*autoopsproto.DatetimeClause {
				sameTime := time.Date(2030, 6, 15, 12, 0, 0, 0, time.UTC).Unix()
				return []*autoopsproto.DatetimeClause{
					{
						Time:       sameTime,
						ActionType: autoopsproto.ActionType_ENABLE,
					},
					{
						Time:       sameTime,
						ActionType: autoopsproto.ActionType_DISABLE,
					},
				}
			}(),
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "success: two one-time clauses different times",
			clauses: []*autoopsproto.DatetimeClause{
				{
					Time:       time.Now().Add(24 * time.Hour).Unix(),
					ActionType: autoopsproto.ActionType_ENABLE,
				},
				{
					Time:       time.Now().Add(48 * time.Hour).Unix(),
					ActionType: autoopsproto.ActionType_DISABLE,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "success: recurring daily and recurring weekly, same time",
			clauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency: autoopsproto.RecurrenceRule_DAILY,
						Timezone:  "UTC",
						StartDate: time.Now().Add(24 * time.Hour).Unix(),
					},
				},
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "err: weekly duplicates with reversed DaysOfWeek order",
			clauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1, 5},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{5, 1},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
			},
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "err: weekly duplicates differing only in irrelevant DayOfMonth",
			clauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						DayOfMonth: 0,
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						DayOfMonth: 15,
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
			},
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
		{
			desc: "err: same time and pattern with different action types are duplicates",
			clauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_DISABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
			},
			expectedErr: statusDatetimeClauseDuplicateTime.Err(),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController)
			err := s.validateDatetimeClauses(p.clauses)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
