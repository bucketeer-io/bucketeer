// Copyright 2024 The Bucketeer Authors.
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

package opsevent

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	aoclientemock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	envclientemock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	eccmock "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	executormock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewEvaluationRealtimeCountPersister(t *testing.T) {
	g := NewEventCountWatcher(nil, nil, nil, nil, nil, nil)
	assert.IsType(t, &eventCountWatcher{}, g)
}

func newNewCountWatcherWithMock(t *testing.T, mockController *gomock.Controller) *eventCountWatcher {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &eventCountWatcher{
		mysqlClient:        mysqlmock.NewMockClient(mockController),
		envClient:          envclientemock.NewMockClient(mockController),
		aoClient:           aoclientemock.NewMockClient(mockController),
		eventCounterClient: eccmock.NewMockClient(mockController),
		featureClient:      ftmock.NewMockClient(mockController),
		autoOpsExecutor:    executormock.NewMockAutoOpsExecutor(mockController),
		logger:             logger,
		opts: &jobs.Options{
			Timeout: time.Minute,
		},
	}
}

func TestRunCountWatcher(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*testing.T, *eventCountWatcher)
		expectedErr error
	}{
		{
			desc: "error: GetFeature fails",
			setup: func(t *testing.T, w *eventCountWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "ns0", ProjectId: "pj0"},
						},
					},
					nil,
				)
				oerc1, _ := newOpsEventRateClauses(t)
				c1, err := anypb.New(oerc1)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:             0,
						EnvironmentNamespace: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:        "id-0",
								FeatureId: "fid-0",
								Clauses:   []*autoopsproto.Clause{{Clause: c1}},
							},
						},
					},
					nil,
				)
				w.featureClient.(*ftmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Internal, "test"))
			},
			expectedErr: status.Errorf(codes.Internal, "test"),
		},
		{
			desc: "error: GetOpsEvaluationUserCount fails",
			setup: func(t *testing.T, w *eventCountWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "ns0", ProjectId: "pj0"},
						},
					},
					nil,
				)
				oerc1, _ := newOpsEventRateClauses(t)
				c1, err := anypb.New(oerc1)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:             0,
						EnvironmentNamespace: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:        "id-0",
								FeatureId: "fid-0",
								Clauses:   []*autoopsproto.Clause{{Clause: c1}},
							},
						},
					},
					nil,
				)
				w.eventCounterClient.(*eccmock.MockClient).
					EXPECT().GetOpsEvaluationUserCount(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.NotFound, "test"))
				w.featureClient.(*ftmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					&ftproto.GetFeatureResponse{
						Feature: &ftproto.Feature{
							Version: 1,
						},
					}, nil)
			},
			expectedErr: status.Errorf(codes.NotFound, "test"),
		},
		{
			desc: "error: GetOpsGoalUserCount fails",
			setup: func(t *testing.T, w *eventCountWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "ns0", ProjectId: "pj0"},
						},
					},
					nil,
				)
				oerc1, _ := newOpsEventRateClauses(t)
				c1, err := anypb.New(oerc1)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:             0,
						EnvironmentNamespace: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:        "id-0",
								FeatureId: "fid-0",
								Clauses:   []*autoopsproto.Clause{{Clause: c1}},
							},
						},
					},
					nil,
				)
				w.eventCounterClient.(*eccmock.MockClient).
					EXPECT().GetOpsEvaluationUserCount(gomock.Any(), gomock.Any()).Return(
					&ecproto.GetOpsEvaluationUserCountResponse{
						OpsRuleId: "rule-id",
						ClauseId:  "clause-id",
						Count:     1,
					}, nil)
				w.eventCounterClient.(*eccmock.MockClient).
					EXPECT().GetOpsGoalUserCount(gomock.Any(), gomock.Any()).
					Return(nil, status.Errorf(codes.NotFound, "test"))
				w.featureClient.(*ftmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					&ftproto.GetFeatureResponse{
						Feature: &ftproto.Feature{
							Version: 1,
						},
					}, nil)
			},
			expectedErr: status.Errorf(codes.NotFound, "test"),
		},
		{
			desc: "success",
			setup: func(t *testing.T, w *eventCountWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "ns0", ProjectId: "pj0"},
						},
					},
					nil,
				)
				oerc1, _ := newOpsEventRateClauses(t)
				c1, err := anypb.New(oerc1)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:             0,
						EnvironmentNamespace: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:          "id-0",
								FeatureId:   "fid-0",
								Clauses:     []*autoopsproto.Clause{{Clause: c1}},
								TriggeredAt: 0,
							},
							{
								Id:          "id-1",
								FeatureId:   "fid-1",
								Clauses:     []*autoopsproto.Clause{{Clause: c1}},
								TriggeredAt: 1,
							},
						},
					},
					nil,
				)
				w.eventCounterClient.(*eccmock.MockClient).
					EXPECT().GetOpsEvaluationUserCount(gomock.Any(), gomock.Any()).Return(
					&ecproto.GetOpsEvaluationUserCountResponse{
						OpsRuleId: "rule-id",
						ClauseId:  "clause-id",
						Count:     15,
					}, nil)
				w.eventCounterClient.(*eccmock.MockClient).
					EXPECT().GetOpsGoalUserCount(gomock.Any(), gomock.Any()).
					Return(
						&ecproto.GetOpsGoalUserCountResponse{
							OpsRuleId: "rule-id",
							ClauseId:  "clause-id",
							Count:     15,
						},
						nil,
					)
				w.featureClient.(*ftmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					&ftproto.GetFeatureResponse{
						Feature: &ftproto.Feature{
							Version: 1,
						},
					}, nil)

				w.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, nil,
				)

				w.autoOpsExecutor.(*executormock.MockAutoOpsExecutor).
					EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNewCountWatcherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(t, s)
			}
			err := s.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCountWatcherAssessRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc               string
		opsEventRateClause *autoopsproto.OpsEventRateClause
		evaluationCount    int64
		opsCount           int64
		expected           bool
	}{
		{
			desc: "GREATER_OR_EQUAL: false: not enough count",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
			},
			evaluationCount: 10,
			opsCount:        4,
			expected:        false,
		},
		{
			desc: "GREATER_OR_EQUAL: false: less than",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
			},
			evaluationCount: 11,
			opsCount:        5,
			expected:        false,
		},
		{
			desc: "GREATER_OR_EQUAL: true: equal",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
			},
			evaluationCount: 10,
			opsCount:        5,
			expected:        true,
		},
		{
			desc: "GREATER_OR_EQUAL: true: greater",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
			},
			evaluationCount: 10,
			opsCount:        6,
			expected:        true,
		},
		{
			desc: "LESS_OR_EQUAL: false: not enough count",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_LESS_OR_EQUAL,
			},
			evaluationCount: 10,
			opsCount:        4,
			expected:        false,
		},
		{
			desc: "LESS_OR_EQUAL: false: greater than",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_LESS_OR_EQUAL,
			},
			evaluationCount: 10,
			opsCount:        6,
			expected:        false,
		},
		{
			desc: "LESS_OR_EQUAL: true: equal",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_LESS_OR_EQUAL,
			},
			evaluationCount: 10,
			opsCount:        5,
			expected:        true,
		},
		{
			desc: "LESS_OR_EQUAL: true: less",
			opsEventRateClause: &autoopsproto.OpsEventRateClause{
				VariationId:     "vid1",
				GoalId:          "gid1",
				MinCount:        int64(5),
				ThreadsholdRate: float64(0.5),
				Operator:        autoopsproto.OpsEventRateClause_LESS_OR_EQUAL,
			},
			evaluationCount: 11,
			opsCount:        5,
			expected:        true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newNewCountWatcherWithMock(t, mockController)
			actual := s.assessRule(p.opsEventRateClause, p.evaluationCount, p.opsCount)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newOpsEventRateClauses(t *testing.T) (*autoopsproto.OpsEventRateClause, *autoopsproto.OpsEventRateClause) {
	t.Helper()
	oerc1 := &autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid1",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	oerc2 := &autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid2",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	return oerc1, oerc2
}
