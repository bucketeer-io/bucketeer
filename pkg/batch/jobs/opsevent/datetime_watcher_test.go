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
// See the License for the specific job.Options job.Options permissions and
// limitations under the TicenseT

package opsevent

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	aoclientemock "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client/mock"
	autoopsdomain "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	envclientemock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	ftcachermock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/cacher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	executormock "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/batch/executor/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func TestNewDatetimeWatcher(t *testing.T) {
	w := NewDatetimeWatcher(nil, nil, nil, nil)
	assert.IsType(t, &datetimeWatcher{}, w)
}

func newNewDatetimeWatcherWithMock(t *testing.T, mockController *gomock.Controller) *datetimeWatcher {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &datetimeWatcher{
		envClient:       envclientemock.NewMockClient(mockController),
		aoClient:        aoclientemock.NewMockClient(mockController),
		autoOpsExecutor: executormock.NewMockAutoOpsExecutor(mockController),
		ftCacher:        ftcachermock.NewMockFeatureFlagCacher(mockController),
		logger:          logger,
		opts: &jobs.Options{
			Timeout: time.Minute,
		},
	}
}

func TestRunDatetimeWatcher(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*datetimeWatcher)
		expectedErr error
	}{
		{
			desc: "success: scheduled clause does not satisfy the time condition",
			setup: func(w *datetimeWatcher) {
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
				dc := &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix()}
				c, err := anypb.New(dc)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:      0,
						EnvironmentId: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:            "id-0",
								FeatureId:     "fid-0",
								Clauses:       []*autoopsproto.Clause{{Clause: c}},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_WAITING,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
							{
								Id:            "id-1",
								FeatureId:     "fid-1",
								Clauses:       []*autoopsproto.Clause{{Clause: c}},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_FINISHED,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
						},
					},
					nil,
				)
			},
			expectedErr: nil,
		},
		{
			desc: "success: scheduled clause satisfies the time condition",
			setup: func(w *datetimeWatcher) {
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
				clause1, err := anypb.New(&autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, -2).Unix()})
				require.NoError(t, err)
				clause2, err := anypb.New(&autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, -1).Unix()})
				require.NoError(t, err)

				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:      0,
						EnvironmentId: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:            "id-0",
								FeatureId:     "fid-0",
								Clauses:       []*autoopsproto.Clause{{Id: "clause-id-0", ExecutedAt: 0, Clause: clause1}},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_WAITING,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
							{
								Id:            "id-1",
								FeatureId:     "fid-1",
								Clauses:       []*autoopsproto.Clause{{Id: "clause-id-1", ExecutedAt: 1, Clause: clause1}},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_FINISHED,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
							{
								Id:        "id-2",
								FeatureId: "fid-2",
								Clauses: []*autoopsproto.Clause{
									{Id: "clause-id-2", ExecutedAt: 1, Clause: clause1},
									{Id: "clause-id-3", ExecutedAt: 0, Clause: clause2},
								},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
						},
					},
					nil,
				)
				w.autoOpsExecutor.(*executormock.MockAutoOpsExecutor).
					EXPECT().Execute(gomock.Any(), "ns0", "id-0", "clause-id-0").Return(nil)
				w.autoOpsExecutor.(*executormock.MockAutoOpsExecutor).
					EXPECT().Execute(gomock.Any(), "ns0", "id-2", "clause-id-3").Return(nil)
				// Cache should be refreshed after successful execution
				w.ftCacher.(*ftcachermock.MockFeatureFlagCacher).
					EXPECT().RefreshEnvironmentCache(gomock.Any(), "ns0").Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: recurring clause satisfies the time condition",
			setup: func(w *datetimeWatcher) {
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
				recurringClause, err := anypb.New(&autoopsproto.DatetimeClause{
					Time:       36000, // 10:00 AM (seconds since midnight)
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().AddDate(0, 0, -7).Unix(),
					},
					NextExecutionAt: time.Now().Add(-1 * time.Hour).Unix(), // Past
					ExecutionCount:  1,
				})
				require.NoError(t, err)

				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:      0,
						EnvironmentId: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:        "id-recurring",
								FeatureId: "fid-recurring",
								Clauses: []*autoopsproto.Clause{{
									Id:          "clause-recurring",
									ExecutedAt:  0,
									IsRecurring: true,
									Clause:      recurringClause,
								}},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
						},
					},
					nil,
				)
				w.autoOpsExecutor.(*executormock.MockAutoOpsExecutor).
					EXPECT().Execute(gomock.Any(), "ns0", "id-recurring", "clause-recurring").Return(nil)
				w.ftCacher.(*ftcachermock.MockFeatureFlagCacher).
					EXPECT().RefreshEnvironmentCache(gomock.Any(), "ns0").Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: recurring clause with future NextExecutionAt is skipped",
			setup: func(w *datetimeWatcher) {
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
				recurringClause, err := anypb.New(&autoopsproto.DatetimeClause{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().AddDate(0, 0, -7).Unix(),
					},
					NextExecutionAt: time.Now().Add(48 * time.Hour).Unix(), // Future
					ExecutionCount:  1,
				})
				require.NoError(t, err)

				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:      0,
						EnvironmentId: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:        "id-recurring",
								FeatureId: "fid-recurring",
								Clauses: []*autoopsproto.Clause{{
									Id:          "clause-recurring",
									ExecutedAt:  0,
									IsRecurring: true,
									Clause:      recurringClause,
								}},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
						},
					},
					nil,
				)
				// No Execute call expected — clause is not ready
			},
			expectedErr: nil,
		},
		{
			desc: "success: exhausted recurring clause (NextExecutionAt=0) is skipped",
			setup: func(w *datetimeWatcher) {
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
				exhaustedClause, err := anypb.New(&autoopsproto.DatetimeClause{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Now().AddDate(0, 0, -30).Unix(),
					},
					NextExecutionAt: 0, // Exhausted
					ExecutionCount:  5,
				})
				require.NoError(t, err)

				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListAutoOpsRules(
					gomock.Any(),
					&autoopsproto.ListAutoOpsRulesRequest{
						PageSize:      0,
						EnvironmentId: "ns0",
					},
				).Return(
					&autoopsproto.ListAutoOpsRulesResponse{
						AutoOpsRules: []*autoopsproto.AutoOpsRule{
							{
								Id:        "id-exhausted",
								FeatureId: "fid-exhausted",
								Clauses: []*autoopsproto.Clause{{
									Id:          "clause-exhausted",
									ExecutedAt:  0,
									IsRecurring: true,
									Clause:      exhaustedClause,
								}},
								AutoOpsStatus: autoopsproto.AutoOpsStatus_RUNNING,
								OpsType:       autoopsproto.OpsType_SCHEDULE,
							},
						},
					},
					nil,
				)
				// No Execute call expected — clause is exhausted
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			w := newNewDatetimeWatcherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(w)
			}
			err := w.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetExecuteClauseId(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc       string
		rule       *autoopsproto.AutoOpsRule
		expectedID string
		expectErr  bool
	}{
		{
			desc: "one-time clause: not ready (future time)",
			rule: func() *autoopsproto.AutoOpsRule {
				dc := &autoopsproto.DatetimeClause{Time: time.Now().Add(24 * time.Hour).Unix()}
				c, err := anypb.New(dc)
				require.NoError(t, err)
				return &autoopsproto.AutoOpsRule{
					Id:        "rule-1",
					FeatureId: "feat-1",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					Clauses:   []*autoopsproto.Clause{{Id: "c1", ExecutedAt: 0, Clause: c}},
				}
			}(),
			expectedID: "",
			expectErr:  false,
		},
		{
			desc: "one-time clause: ready (past time)",
			rule: func() *autoopsproto.AutoOpsRule {
				dc := &autoopsproto.DatetimeClause{Time: time.Now().Add(-1 * time.Hour).Unix()}
				c, err := anypb.New(dc)
				require.NoError(t, err)
				return &autoopsproto.AutoOpsRule{
					Id:        "rule-1",
					FeatureId: "feat-1",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					Clauses:   []*autoopsproto.Clause{{Id: "c1", ExecutedAt: 0, Clause: c}},
				}
			}(),
			expectedID: "c1",
			expectErr:  false,
		},
		{
			desc: "one-time clause: already executed",
			rule: func() *autoopsproto.AutoOpsRule {
				dc := &autoopsproto.DatetimeClause{Time: time.Now().Add(-1 * time.Hour).Unix()}
				c, err := anypb.New(dc)
				require.NoError(t, err)
				return &autoopsproto.AutoOpsRule{
					Id:        "rule-1",
					FeatureId: "feat-1",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					Clauses:   []*autoopsproto.Clause{{Id: "c1", ExecutedAt: 1, Clause: c}},
				}
			}(),
			expectedID: "",
			expectErr:  false,
		},
		{
			desc: "recurring clause: ready (NextExecutionAt in past)",
			rule: func() *autoopsproto.AutoOpsRule {
				dc := &autoopsproto.DatetimeClause{
					Time:            36000,
					NextExecutionAt: time.Now().Add(-1 * time.Hour).Unix(),
					ExecutionCount:  1,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
					},
				}
				c, err := anypb.New(dc)
				require.NoError(t, err)
				return &autoopsproto.AutoOpsRule{
					Id:        "rule-1",
					FeatureId: "feat-1",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					Clauses:   []*autoopsproto.Clause{{Id: "c1", ExecutedAt: 0, IsRecurring: true, Clause: c}},
				}
			}(),
			expectedID: "c1",
			expectErr:  false,
		},
		{
			desc: "recurring clause: not ready (NextExecutionAt in future)",
			rule: func() *autoopsproto.AutoOpsRule {
				dc := &autoopsproto.DatetimeClause{
					Time:            36000,
					NextExecutionAt: time.Now().Add(48 * time.Hour).Unix(),
					ExecutionCount:  1,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
					},
				}
				c, err := anypb.New(dc)
				require.NoError(t, err)
				return &autoopsproto.AutoOpsRule{
					Id:        "rule-1",
					FeatureId: "feat-1",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					Clauses:   []*autoopsproto.Clause{{Id: "c1", ExecutedAt: 0, IsRecurring: true, Clause: c}},
				}
			}(),
			expectedID: "",
			expectErr:  false,
		},
		{
			desc: "recurring clause: exhausted (NextExecutionAt=0)",
			rule: func() *autoopsproto.AutoOpsRule {
				dc := &autoopsproto.DatetimeClause{
					Time:            36000,
					NextExecutionAt: 0,
					ExecutionCount:  5,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
					},
				}
				c, err := anypb.New(dc)
				require.NoError(t, err)
				return &autoopsproto.AutoOpsRule{
					Id:        "rule-1",
					FeatureId: "feat-1",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					Clauses:   []*autoopsproto.Clause{{Id: "c1", ExecutedAt: 0, IsRecurring: true, Clause: c}},
				}
			}(),
			expectedID: "",
			expectErr:  false,
		},
		{
			desc: "mixed: recurring ready and one-time ready, returns earliest",
			rule: func() *autoopsproto.AutoOpsRule {
				earlierTime := time.Now().Add(-2 * time.Hour).Unix()
				laterTime := time.Now().Add(-1 * time.Hour).Unix()

				recurringDC := &autoopsproto.DatetimeClause{
					Time:            36000,
					NextExecutionAt: laterTime,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency: autoopsproto.RecurrenceRule_DAILY,
						Timezone:  "UTC",
					},
				}
				oneTimeDC := &autoopsproto.DatetimeClause{
					Time: earlierTime,
				}
				rc, err := anypb.New(recurringDC)
				require.NoError(t, err)
				oc, err := anypb.New(oneTimeDC)
				require.NoError(t, err)

				return &autoopsproto.AutoOpsRule{
					Id:        "rule-mixed",
					FeatureId: "feat-mixed",
					OpsType:   autoopsproto.OpsType_SCHEDULE,
					Clauses: []*autoopsproto.Clause{
						{Id: "recurring-c", ExecutedAt: 0, IsRecurring: true, Clause: rc},
						{Id: "onetime-c", ExecutedAt: 0, IsRecurring: false, Clause: oc},
					},
				}
			}(),
			expectedID: "onetime-c", // earlier execution time
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			w := newNewDatetimeWatcherWithMock(t, mockController)

			aor := &autoopsdomain.AutoOpsRule{AutoOpsRule: tt.rule}
			id, err := w.getExecuteClauseId("env-test", aor)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)
		})
	}
}
