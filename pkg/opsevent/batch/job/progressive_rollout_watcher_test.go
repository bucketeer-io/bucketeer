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

package job

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	environmentdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	executormock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor/mock"
	targetstoremock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestNewProgressiveRolloutWacher(t *testing.T) {
	w := NewProgressiveRolloutWacher(nil, nil)
	assert.IsType(t, &progressiveRolloutWatcher{}, w)
}

func newProgressiveRolloutWacherWithMock(t *testing.T, mockController *gomock.Controller) *progressiveRolloutWatcher {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &progressiveRolloutWatcher{
		environmentLister:          targetstoremock.NewMockEnvironmentLister(mockController),
		progressiveRolloutLister:   targetstoremock.NewMockProgressiveRolloutLister(mockController),
		progressiveRolloutExecutor: executormock.NewMockProgressiveRolloutExecutor(mockController),
		logger:                     logger,
		opts: &options{
			timeout: time.Minute,
		},
	}
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
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId: "sid1",
							ExecuteAt:  time.Now().Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.progressiveRolloutLister.(*targetstoremock.MockProgressiveRolloutLister).EXPECT().GetProgressiveRollouts(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.ProgressiveRollout{
						{ProgressiveRollout: &autoopsproto.ProgressiveRollout{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clause:    c,
							Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
						}},
					},
				)
				w.progressiveRolloutExecutor.(*executormock.MockProgressiveRolloutExecutor).EXPECT().ExecuteProgressiveRollout(gomock.Any(), "ns0", "id-0", "sid1").Return(internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc: "success: executed_at is not past time",
			setup: func(w *progressiveRolloutWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId: "sid1",
							ExecuteAt:  time.Now().AddDate(300, 0, 0).Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.progressiveRolloutLister.(*targetstoremock.MockProgressiveRolloutLister).EXPECT().GetProgressiveRollouts(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.ProgressiveRollout{
						{ProgressiveRollout: &autoopsproto.ProgressiveRollout{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clause:    c,
							Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
						}},
					},
				)
			},
			expectedErr: nil,
		},
		{
			desc: "success: executed_at is past time",
			setup: func(w *progressiveRolloutWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId: "sid1",
							ExecuteAt:  time.Now().Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.progressiveRolloutLister.(*targetstoremock.MockProgressiveRolloutLister).EXPECT().GetProgressiveRollouts(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.ProgressiveRollout{
						{ProgressiveRollout: &autoopsproto.ProgressiveRollout{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clause:    c,
							Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
						}},
					},
				)
				w.progressiveRolloutExecutor.(*executormock.MockProgressiveRolloutExecutor).EXPECT().ExecuteProgressiveRollout(gomock.Any(), "ns0", "id-0", "sid1").Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: executed_at is past time and already triggerred",
			setup: func(w *progressiveRolloutWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				dc := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
					Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
						{
							ScheduleId:  "sid1",
							ExecuteAt:   time.Now().Unix(),
							TriggeredAt: time.Now().Unix(),
						},
					},
				}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.progressiveRolloutLister.(*targetstoremock.MockProgressiveRolloutLister).EXPECT().GetProgressiveRollouts(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.ProgressiveRollout{
						{ProgressiveRollout: &autoopsproto.ProgressiveRollout{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clause:    c,
							Type:      autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
						}},
					},
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
