// Copyright 2023 The Bucketeer Authors.
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
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/wrapperspb"

	aoclientemock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	envclientemock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	executormock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestNewProgressiveRolloutWacher(t *testing.T) {
	w := NewProgressiveRolloutWacher(nil, nil, nil)
	assert.IsType(t, &progressiveRolloutWatcher{}, w)
}

func TestRunProgressiveRolloutWatcher(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	internalErr := errors.New("error")

	patterns := []struct {
		desc        string
		setup       func(*progressiveRolloutWatcher)
		expectedErr error
	}{
		{
			desc: "fail: internal error",
			setup: func(w *progressiveRolloutWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "eID", ProjectId: "pID"},
						},
					},
					nil,
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId: "sID",
							ExecuteAt:  time.Now().Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&autoopsproto.ListProgressiveRolloutsRequest{
						EnvironmentNamespace: "eID",
						PageSize:             0,
					},
				).Return(
					&autoopsproto.ListProgressiveRolloutsResponse{
						ProgressiveRollouts: []*autoopsproto.ProgressiveRollout{
							{
								Id:        "prID",
								FeatureId: "fID",
								Clause:    c,
								Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
							},
						},
					},
					nil,
				)
				w.progressiveRolloutExecutor.(*executormock.MockProgressiveRolloutExecutor).EXPECT().ExecuteProgressiveRollout(
					gomock.Any(), "eID", "prID", "sID",
				).Return(internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc: "success: executed_at is not past time",
			setup: func(w *progressiveRolloutWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "eID", ProjectId: "pID"},
						},
					},
					nil,
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId: "sID",
							ExecuteAt:  time.Now().AddDate(300, 0, 0).Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&autoopsproto.ListProgressiveRolloutsRequest{
						EnvironmentNamespace: "eID",
						PageSize:             0,
					},
				).Return(
					&autoopsproto.ListProgressiveRolloutsResponse{
						ProgressiveRollouts: []*autoopsproto.ProgressiveRollout{
							{
								Id:        "sID",
								FeatureId: "fID",
								Clause:    c,
								Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
							},
						},
					},
					nil,
				)
			},
			expectedErr: nil,
		},
		{
			desc: "success: executed_at is past time",
			setup: func(w *progressiveRolloutWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "eID", ProjectId: "pID"},
						},
					},
					nil,
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId: "sID",
							ExecuteAt:  time.Now().Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&autoopsproto.ListProgressiveRolloutsRequest{
						EnvironmentNamespace: "eID",
						PageSize:             0,
					},
				).Return(
					&autoopsproto.ListProgressiveRolloutsResponse{
						ProgressiveRollouts: []*autoopsproto.ProgressiveRollout{
							{
								Id:        "sID",
								FeatureId: "fID",
								Clause:    c,
								Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
							},
						},
					},
					nil,
				)
				w.progressiveRolloutExecutor.(*executormock.MockProgressiveRolloutExecutor).EXPECT().ExecuteProgressiveRollout(
					gomock.Any(), "eID", "sID", "sID",
				).Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: executed_at is past time and already triggerred",
			setup: func(w *progressiveRolloutWatcher) {
				w.envClient.(*envclientemock.MockClient).EXPECT().ListEnvironmentsV2(
					gomock.Any(),
					&environmentproto.ListEnvironmentsV2Request{
						PageSize: 0,
						Archived: wrapperspb.Bool(false),
					},
				).Return(
					&environmentproto.ListEnvironmentsV2Response{
						Environments: []*environmentproto.EnvironmentV2{
							{Id: "eID", ProjectId: "pID"},
						},
					},
					nil,
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId:  "sID",
							ExecuteAt:   time.Now().Unix(),
							TriggeredAt: time.Now().Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.aoClient.(*aoclientemock.MockClient).EXPECT().ListProgressiveRollouts(
					gomock.Any(),
					&autoopsproto.ListProgressiveRolloutsRequest{
						EnvironmentNamespace: "eID",
						PageSize:             0,
					},
				).Return(
					&autoopsproto.ListProgressiveRolloutsResponse{
						ProgressiveRollouts: []*autoopsproto.ProgressiveRollout{
							{
								Id:        "sID",
								FeatureId: "fID",
								Clause:    c,
								Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
							},
						},
					},
					nil,
				)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			w := newProgressiveRolloutWacherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(w)
			}
			err := w.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newProgressiveRolloutWacherWithMock(t *testing.T, mockController *gomock.Controller) *progressiveRolloutWatcher {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &progressiveRolloutWatcher{
		envClient:                  envclientemock.NewMockClient(mockController),
		aoClient:                   aoclientemock.NewMockClient(mockController),
		progressiveRolloutExecutor: executormock.NewMockProgressiveRolloutExecutor(mockController),
		logger:                     logger,
		opts: &jobs.Options{
			Timeout: time.Minute,
		},
	}
}
