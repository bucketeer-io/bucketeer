// Copyright 2022 The Bucketeer Authors.
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

package job

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	environmentdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	eccmock "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	executormock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor/mock"
	targetstoremock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewEvaluationRealtimeCountPersister(t *testing.T) {
	g := NewCountWatcher(nil, nil, nil, nil, nil)
	assert.IsType(t, &countWatcher{}, g)
}

func newNewCountWatcherWithMock(t *testing.T, mockController *gomock.Controller) *countWatcher {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &countWatcher{
		mysqlClient:        mysqlmock.NewMockClient(mockController),
		environmentLister:  targetstoremock.NewMockEnvironmentLister(mockController),
		autoOpsRuleLister:  targetstoremock.NewMockAutoOpsRuleLister(mockController),
		eventCounterClient: eccmock.NewMockClient(mockController),
		featureClient:      ftmock.NewMockClient(mockController),
		autoOpsExecutor:    executormock.NewMockAutoOpsExecutor(mockController),
		logger:             logger,
		opts: &options{
			timeout: time.Minute,
		},
	}
}

func TestRunCountWatcher(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*testing.T, *countWatcher)
		expectedErr error
	}{
		{
			desc: "error: GetFeature fails",
			setup: func(t *testing.T, w *countWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				oerc1, _ := newOpsEventRateClauses(t)
				c1, err := ptypes.MarshalAny(oerc1)
				require.NoError(t, err)
				w.autoOpsRuleLister.(*targetstoremock.MockAutoOpsRuleLister).EXPECT().GetAutoOpsRules(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.AutoOpsRule{
						{AutoOpsRule: &autoopsproto.AutoOpsRule{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clauses:   []*autoopsproto.Clause{{Clause: c1}},
						}},
					},
				)
				w.featureClient.(*ftmock.MockClient).EXPECT().GetFeature(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Internal, "test"))
			},
			expectedErr: status.Errorf(codes.Internal, "test"),
		},
		{
			desc: "error: GetEvaluationRealtimeCount fails",
			setup: func(t *testing.T, w *countWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				oerc1, _ := newOpsEventRateClauses(t)
				c1, err := ptypes.MarshalAny(oerc1)
				require.NoError(t, err)
				w.autoOpsRuleLister.(*targetstoremock.MockAutoOpsRuleLister).EXPECT().GetAutoOpsRules(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.AutoOpsRule{
						{AutoOpsRule: &autoopsproto.AutoOpsRule{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clauses:   []*autoopsproto.Clause{{Clause: c1}},
						}},
					},
				)
				w.eventCounterClient.(*eccmock.MockClient).EXPECT().GetEvaluationCountV2(gomock.Any(), gomock.Any()).Return(
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
			desc: "error: GetOpsRealtimeVariationCount fails",
			setup: func(t *testing.T, w *countWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				oerc1, _ := newOpsEventRateClauses(t)
				c1, err := ptypes.MarshalAny(oerc1)
				require.NoError(t, err)
				w.autoOpsRuleLister.(*targetstoremock.MockAutoOpsRuleLister).EXPECT().GetAutoOpsRules(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.AutoOpsRule{
						{AutoOpsRule: &autoopsproto.AutoOpsRule{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clauses:   []*autoopsproto.Clause{{Clause: c1}},
						}},
					},
				)
				w.eventCounterClient.(*eccmock.MockClient).EXPECT().GetEvaluationCountV2(gomock.Any(), gomock.Any()).Return(
					&ecproto.GetEvaluationCountV2Response{Count: &ecproto.EvaluationCount{
						RealtimeCounts: []*ecproto.VariationCount{{VariationId: "vid1", UserCount: 1}},
					}}, nil)
				w.eventCounterClient.(*eccmock.MockClient).EXPECT().GetGoalCountV2(gomock.Any(), gomock.Any()).Return(
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
		evaluationCount    *ecproto.VariationCount
		opsCount           *ecproto.VariationCount
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
			evaluationCount: &ecproto.VariationCount{UserCount: 10},
			opsCount:        &ecproto.VariationCount{UserCount: 4},
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
			evaluationCount: &ecproto.VariationCount{UserCount: 11},
			opsCount:        &ecproto.VariationCount{UserCount: 5},
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
			evaluationCount: &ecproto.VariationCount{UserCount: 10},
			opsCount:        &ecproto.VariationCount{UserCount: 5},
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
			evaluationCount: &ecproto.VariationCount{UserCount: 10},
			opsCount:        &ecproto.VariationCount{UserCount: 6},
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
			evaluationCount: &ecproto.VariationCount{UserCount: 10},
			opsCount:        &ecproto.VariationCount{UserCount: 4},
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
			evaluationCount: &ecproto.VariationCount{UserCount: 10},
			opsCount:        &ecproto.VariationCount{UserCount: 6},
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
			evaluationCount: &ecproto.VariationCount{UserCount: 10},
			opsCount:        &ecproto.VariationCount{UserCount: 5},
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
			evaluationCount: &ecproto.VariationCount{UserCount: 11},
			opsCount:        &ecproto.VariationCount{UserCount: 5},
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
