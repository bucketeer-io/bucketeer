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

	aoclientemock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	envclientemock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	executormock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestNewDatetimeWatcher(t *testing.T) {
	w := NewDatetimeWatcher(nil, nil, nil)
	assert.IsType(t, &datetimeWatcher{}, w)
}

func newNewDatetimeWatcherWithMock(t *testing.T, mockController *gomock.Controller) *datetimeWatcher {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &datetimeWatcher{
		envClient:       envclientemock.NewMockClient(mockController),
		aoClient:        aoclientemock.NewMockClient(mockController),
		autoOpsExecutor: executormock.NewMockAutoOpsExecutor(mockController),
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
