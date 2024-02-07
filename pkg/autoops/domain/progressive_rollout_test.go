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

package domain

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoiface"

	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestNewProgressiveRollout(t *testing.T) {
	t.Parallel()
	aor := createProgressiveRollout(t)
	assert.IsType(t, &ProgressiveRollout{}, aor)
	assert.Equal(t, "feature-id", aor.FeatureId)
	assert.Equal(t, autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE, aor.Type)
	assert.NotNil(t, aor.Clause)
	assert.NotZero(t, aor.CreatedAt)
	assert.NotZero(t, aor.UpdatedAt)
}

func createProgressiveRollout(t *testing.T) *ProgressiveRollout {
	aor, err := NewProgressiveRollout(
		"feature-id",
		nil,
		&autoopsproto.ProgressiveRolloutTemplateScheduleClause{
			Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
				{
					ScheduleId: "schedule-id-0",
					ExecuteAt:  time.Now().Unix(),
					Weight:     0,
				},
				{
					ScheduleId: "schedule-id-1",
					ExecuteAt:  time.Now().AddDate(1, 0, 0).Unix(),
					Weight:     20,
				},
				{
					ScheduleId: "schedule-id-2",
					ExecuteAt:  time.Now().AddDate(2, 0, 0).Unix(),
					Weight:     40,
				},
				{
					ScheduleId: "schedule-id-3",
					ExecuteAt:  time.Now().AddDate(3, 0, 0).Unix(),
					Weight:     60,
				},
				{
					ScheduleId: "schedule-id-4",
					ExecuteAt:  time.Now().AddDate(4, 0, 0).Unix(),
					Weight:     80,
				},
				{
					ScheduleId: "schedule-id-5",
					ExecuteAt:  time.Now().AddDate(5, 0, 0).Unix(),
					Weight:     100,
				},
			},
			Interval:    autoopsproto.ProgressiveRolloutTemplateScheduleClause_DAILY,
			Increments:  20,
			VariationId: "vid-1",
		},
	)
	require.NoError(t, err)
	return aor
}

func TestAlreadyTriggered(t *testing.T) {
	patterns := []struct {
		desc                   string
		progressiveRolloutType autoopsproto.ProgressiveRollout_Type
		clause                 protoiface.MessageV1
		targetScheduleID       string
		expected               bool
		expectedErr            error
	}{
		{
			desc:                   "err: template ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: manual ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: ErrProgressiveRolloutInvalidType",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_Type(10),
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutInvalidType,
		},
		{
			desc:                   "success: false",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-2",
			expected:         false,
		},
		{
			desc:                   "success: true",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId:  "sid-2",
						TriggeredAt: time.Now().AddDate(0, -1, 0).Unix(),
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-2",
			expected:         true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createProgressiveRollout(t)
			ac, err := ptypes.MarshalAny(p.clause)
			assert.NoError(t, err)
			s.Clause = ac
			s.Type = p.progressiveRolloutType
			triggered, err := s.AlreadyTriggered(p.targetScheduleID)
			assert.Equal(t, triggered, p.expected)
			assert.Equal(t, err, p.expectedErr)
		})
	}
}

func TestProgressiveRolloutSetTriggeredAt(t *testing.T) {
	patterns := []struct {
		desc                   string
		progressiveRolloutType autoopsproto.ProgressiveRollout_Type
		clause                 protoiface.MessageV1
		targetScheduleID       string
		expectedErr            error
		expectedStatus         autoopsproto.ProgressiveRollout_Status
	}{
		{
			desc:                   "err: template ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: manual ErrProgressiveRolloutScheduleNotFound",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutScheduleNotFound,
		},
		{
			desc:                   "err: ErrProgressiveRolloutInvalidType",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_Type(10),
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-1",
			expectedErr:      ErrProgressiveRolloutInvalidType,
		},
		{
			desc:                   "success",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-2",
			expectedStatus:   autoopsproto.ProgressiveRollout_RUNNING,
		},
		{
			desc:                   "success last schedule is executed",
			progressiveRolloutType: autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE,
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-2",
					},
					{
						ScheduleId: "sid-3",
					},
					{
						ScheduleId: "sid-4",
					},
				},
			},
			targetScheduleID: "sid-4",
			expectedStatus:   autoopsproto.ProgressiveRollout_FINISHED,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createProgressiveRollout(t)
			ac, err := ptypes.MarshalAny(p.clause)
			assert.NoError(t, err)
			s.Clause = ac
			s.Type = p.progressiveRolloutType
			err = s.SetTriggeredAt(p.targetScheduleID)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr == nil {
				c, err := unmarshalProgressiveRolloutManualClause(s.Clause)
				assert.NoError(t, err)
				s, err := findTargetSchedule(c.Schedules, p.targetScheduleID)
				assert.NoError(t, err)
				assert.NotZero(t, s.TriggeredAt)
			}
			assert.Equal(t, p.expectedStatus, s.Status)
		})
	}
}

func TestAddManualScheduleClause(t *testing.T) {
	patterns := []struct {
		desc   string
		clause *autoopsproto.ProgressiveRolloutManualScheduleClause
	}{
		{
			desc: "success",
			clause: &autoopsproto.ProgressiveRolloutManualScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-1",
						Weight:     10,
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pro := createProgressiveRollout(t)
			pro.Clause = nil
			assert.Nil(t, pro.Clause)
			pro.addManualScheduleClause(p.clause)
			assert.NotNil(t, pro.Clause)
		})
	}
}

func TestAddTemplateScheduleClause(t *testing.T) {
	patterns := []struct {
		desc   string
		clause *autoopsproto.ProgressiveRolloutTemplateScheduleClause
	}{
		{
			desc: "success",
			clause: &autoopsproto.ProgressiveRolloutTemplateScheduleClause{
				Schedules: []*autoopsproto.ProgressiveRolloutSchedule{
					{
						ScheduleId: "sid-1",
						Weight:     10,
					},
				},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pro := createProgressiveRollout(t)
			pro.Clause = nil
			assert.Nil(t, pro.Clause)
			pro.addTemplatelScheduleClause(p.clause)
			assert.NotNil(t, pro.Clause)
		})
	}
}

func TestExtractSchedules(t *testing.T) {
	p := createProgressiveRollout(t)
	actual, err := p.ExtractSchedules()
	assert.NoError(t, err)
	assert.Len(t, actual, 6)
	assert.Equal(t, actual[1].Weight, int32(20))
	assert.Equal(t, actual[5].Weight, int32(100))
}

func TestStop(t *testing.T) {
	patterns := []struct {
		desc     string
		input    autoopsproto.ProgressiveRollout_StoppedBy
		expected error
	}{
		{
			desc:     "err: stopped by is required",
			input:    autoopsproto.ProgressiveRollout_UNKNOWN,
			expected: ErrProgressiveRolloutStoopedByRequired,
		},
		{
			desc:     "success: by user",
			input:    autoopsproto.ProgressiveRollout_USER,
			expected: nil,
		},
		{
			desc:     "success: by schedule",
			input:    autoopsproto.ProgressiveRollout_OPS_SCHEDULE,
			expected: nil,
		},
		{
			desc:     "success: by kill switch",
			input:    autoopsproto.ProgressiveRollout_OPS_KILL_SWITCH,
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			pr := createProgressiveRollout(t)
			err := pr.Stop(p.input)
			if err != nil {
				assert.Equal(t, p.expected, err, p.desc)
				assert.Equal(t, autoopsproto.ProgressiveRollout_WAITING, pr.Status, p.desc)
				assert.Equal(t, autoopsproto.ProgressiveRollout_UNKNOWN, pr.StoppedBy, p.desc)
				assert.Zero(t, pr.StoppedAt, p.desc)
				assert.NotZero(t, pr.UpdatedAt, p.desc)
			} else {
				assert.Equal(t, p.expected, err, p.desc)
				assert.Equal(t, autoopsproto.ProgressiveRollout_STOPPED, pr.Status, p.desc)
				assert.Equal(t, p.input, pr.StoppedBy, p.desc)
				assert.NotZero(t, pr.StoppedAt, p.desc)
				assert.True(t, pr.UpdatedAt > time.Now().Add(time.Second*-2).Unix(), p.desc)
			}
		})
	}
}
