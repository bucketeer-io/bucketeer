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

package domain

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

func TestNewAutoOpsRule(t *testing.T) {
	patterns := []struct {
		desc             string
		featureId        string
		opsType          autoopsproto.OpsType
		datetimeClauses  []*autoopsproto.DatetimeClause
		eventRateClauses []*autoopsproto.OpsEventRateClause
		expected         *AutoOpsRule
		expectedErr      error
	}{
		{
			desc:      "OpsType: Schedule",
			featureId: "feature-id",
			opsType:   autoopsproto.OpsType_SCHEDULE,
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{Time: 2, ActionType: autoopsproto.ActionType_ENABLE},
				{Time: 1, ActionType: autoopsproto.ActionType_DISABLE},
			},
			eventRateClauses: []*autoopsproto.OpsEventRateClause{
				{
					GoalId:          "goalid01",
					MinCount:        10,
					ThreadsholdRate: 0.5,
					Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					ActionType:      autoopsproto.ActionType_DISABLE,
				},
			},
			expectedErr: nil,
			expected: &AutoOpsRule{&autoopsproto.AutoOpsRule{FeatureId: "feature-id",
				OpsType: autoopsproto.OpsType_SCHEDULE,
				Clauses: []*autoopsproto.Clause{
					newDateTimeClause(t, &autoopsproto.DatetimeClause{Time: 1, ActionType: autoopsproto.ActionType_DISABLE}),
					newDateTimeClause(t, &autoopsproto.DatetimeClause{Time: 2, ActionType: autoopsproto.ActionType_ENABLE}),
				},
				CreatedAt:     time.Now().Unix(),
				UpdatedAt:     time.Now().Unix(),
				Deleted:       false,
				AutoOpsStatus: autoopsproto.AutoOpsStatus_WAITING,
			}},
		},
		{
			desc:      "OpsType: EventRate",
			featureId: "feature-id",
			opsType:   autoopsproto.OpsType_EVENT_RATE,
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{Time: 1, ActionType: autoopsproto.ActionType_ENABLE},
				{Time: 0, ActionType: autoopsproto.ActionType_DISABLE},
			},
			eventRateClauses: []*autoopsproto.OpsEventRateClause{
				{
					GoalId:          "goalid01",
					MinCount:        10,
					ThreadsholdRate: 0.5,
					Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					ActionType:      autoopsproto.ActionType_DISABLE,
				},
			},
			expectedErr: nil,
			expected: &AutoOpsRule{&autoopsproto.AutoOpsRule{FeatureId: "feature-id",
				OpsType: autoopsproto.OpsType_EVENT_RATE,
				Clauses: []*autoopsproto.Clause{
					newEventRateClause(t, &autoopsproto.OpsEventRateClause{
						GoalId:          "goalid01",
						MinCount:        10,
						ThreadsholdRate: 0.5,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						ActionType:      autoopsproto.ActionType_DISABLE,
					}),
				},
				CreatedAt:     time.Now().Unix(),
				UpdatedAt:     time.Now().Unix(),
				Deleted:       false,
				AutoOpsStatus: autoopsproto.AutoOpsStatus_WAITING,
			}},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			aor, err := NewAutoOpsRule(
				p.featureId,
				p.opsType,
				p.eventRateClauses,
				p.datetimeClauses,
			)
			assert.Equal(t, p.expectedErr, err)

			if err == nil {
				assert.Equal(t, p.expected.FeatureId, aor.FeatureId)
				assert.Equal(t, p.expected.OpsType, aor.OpsType)
				assert.Equal(t, p.expected.AutoOpsStatus, aor.AutoOpsStatus)
				assert.Equal(t, p.expected.CreatedAt, aor.CreatedAt)
				assert.Equal(t, p.expected.UpdatedAt, aor.UpdatedAt)
				assert.Equal(t, p.expected.Deleted, aor.Deleted)

				for i, c := range aor.Clauses {
					assert.Equal(t, p.expected.Clauses[i].ActionType, c.ActionType)
					assert.Equal(t, p.expected.Clauses[i].Clause, c.Clause)
				}
			}
		})
	}
}

func TestSetDeleted(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.SetDeleted()
	assert.NotZero(t, aor.UpdatedAt)
	assert.Equal(t, true, aor.Deleted)
}

func TestAddOpsEventRateClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	l := len(aor.Clauses)
	c := &autoopsproto.OpsEventRateClause{
		GoalId:          "goalid01",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
		ActionType:      autoopsproto.ActionType_DISABLE,
	}
	clause, err := aor.AddOpsEventRateClause(c)
	require.NoError(t, err)
	assert.NotNil(t, clause)
	assert.NotEmpty(t, aor.Clauses[l].Id)
	eventRateClause, err := aor.unmarshalOpsEventRateClause(aor.Clauses[l])
	require.NoError(t, err)

	assert.Equal(t, c.GoalId, eventRateClause.GoalId)
	assert.Equal(t, c.MinCount, eventRateClause.MinCount)
	assert.Equal(t, c.ThreadsholdRate, eventRateClause.ThreadsholdRate)
	assert.Equal(t, c.Operator, eventRateClause.Operator)
	assert.Equal(t, c.ActionType, eventRateClause.ActionType)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)
}

func TestAddDatetimeClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	l := len(aor.Clauses)
	c1 := &autoopsproto.DatetimeClause{
		Time:       1000000001,
		ActionType: autoopsproto.ActionType_DISABLE,
	}
	c2 := &autoopsproto.DatetimeClause{
		Time:       1000000000,
		ActionType: autoopsproto.ActionType_DISABLE,
	}

	clause, err := aor.AddDatetimeClause(c1)
	require.NoError(t, err)
	assert.NotNil(t, clause)
	assert.NotEmpty(t, clause.Id)
	dc, err := aor.unmarshalDatetimeClause(aor.Clauses[l])
	assert.Equal(t, c1.ActionType, aor.Clauses[l].ActionType)
	require.NoError(t, err)
	assert.Equal(t, c1.Time, dc.Time)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)

	clause2, err := aor.AddDatetimeClause(c2)
	require.NoError(t, err)
	assert.NotNil(t, clause2)
	assert.NotEmpty(t, clause2.Id)
	dc2, err := aor.unmarshalDatetimeClause(aor.Clauses[l])
	require.NoError(t, err)
	assert.Equal(t, c2.Time, dc2.Time)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)
}

func TestChangeOpsEventRateClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	l := len(aor.Clauses)
	c := &autoopsproto.OpsEventRateClause{
		GoalId:          "goalid01",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	err := aor.ChangeOpsEventRateClause(aor.Clauses[0].Id, c)
	require.NoError(t, err)
	assert.Equal(t, l, len(aor.Clauses))
	eventRateClause, err := aor.unmarshalOpsEventRateClause(aor.Clauses[0])
	require.NoError(t, err)
	assert.Equal(t, c.GoalId, eventRateClause.GoalId)
}

func TestChangeDatetimeClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	l := len(aor.Clauses)
	c := &autoopsproto.DatetimeClause{
		Time:       1,
		ActionType: autoopsproto.ActionType_DISABLE,
	}
	err := aor.ChangeDatetimeClause(aor.Clauses[0].Id, c)
	require.NoError(t, err)
	assert.Equal(t, l, len(aor.Clauses))
	assert.Equal(t, c.ActionType, aor.Clauses[0].ActionType)
	dc, err := aor.unmarshalDatetimeClause(aor.Clauses[0])
	require.NoError(t, err)
	assert.Equal(t, c.Time, dc.Time)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)

	c1 := &autoopsproto.DatetimeClause{
		Time:       3,
		ActionType: autoopsproto.ActionType_DISABLE,
	}
	c2 := &autoopsproto.DatetimeClause{
		Time:       5,
		ActionType: autoopsproto.ActionType_ENABLE,
	}
	addClause1, err := aor.AddDatetimeClause(c1)
	addClause2, err := aor.AddDatetimeClause(c2)
	require.NoError(t, err)
	assert.Equal(t, 3, len(aor.Clauses))

	cc := &autoopsproto.DatetimeClause{
		Time:       2,
		ActionType: autoopsproto.ActionType_DISABLE,
	}

	err = aor.ChangeDatetimeClause(addClause2.Id, cc)
	require.NoError(t, err)
	assert.Equal(t, 3, len(aor.Clauses))

	dc1, err := aor.unmarshalDatetimeClause(aor.Clauses[0])
	require.NoError(t, err)
	assert.Equal(t, c.Time, dc1.Time)
	dc2, err := aor.unmarshalDatetimeClause(aor.Clauses[1])
	require.NoError(t, err)
	assert.Equal(t, aor.Clauses[1].Id, addClause2.Id)
	assert.Equal(t, cc.Time, dc2.Time)
	dc3, err := aor.unmarshalDatetimeClause(aor.Clauses[2])
	require.NoError(t, err)
	assert.Equal(t, aor.Clauses[2].Id, addClause1.Id)
	assert.Equal(t, c1.Time, dc3.Time)
}

func TestDeleteClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	l := len(aor.Clauses)
	c := &autoopsproto.OpsEventRateClause{
		GoalId:          "goalid01",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
		ActionType:      autoopsproto.ActionType_DISABLE,
	}
	addClause, err := aor.AddOpsEventRateClause(c)
	require.NoError(t, err)
	assert.Equal(t, l+1, len(aor.Clauses))
	err = aor.DeleteClause(aor.Clauses[0].Id)
	require.NoError(t, err)
	assert.Equal(t, l, len(aor.Clauses))
	assert.Equal(t, addClause.Id, aor.Clauses[0].Id)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)
}

func createAutoOpsRule(t *testing.T) *AutoOpsRule {
	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{Time: 0, ActionType: autoopsproto.ActionType_ENABLE},
		},
	)
	require.NoError(t, err)
	return aor
}

func TestExtractOpsEventRateClauses(t *testing.T) {
	oerc1 := &autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid1",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	c1, err := ptypes.MarshalAny(oerc1)
	require.NoError(t, err)
	oerc2 := &autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid2",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	c2, err := ptypes.MarshalAny(oerc2)
	require.NoError(t, err)
	dc1 := &autoopsproto.DatetimeClause{
		Time: 1000000001,
	}
	c3, err := ptypes.MarshalAny(dc1)
	require.NoError(t, err)
	autoOpsRule := &AutoOpsRule{&autoopsproto.AutoOpsRule{
		Id:        "id-0",
		FeatureId: "fid-0",
		Clauses:   []*autoopsproto.Clause{{Id: "c1", Clause: c1}, {Id: "c2", Clause: c2}, {Id: "c3", Clause: c3}},
	}}
	expected := map[string]*autoopsproto.OpsEventRateClause{"c1": oerc1, "c2": oerc2}
	actual, err := autoOpsRule.ExtractOpsEventRateClauses()
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(actual))
	for i, a := range actual {
		assert.True(t, proto.Equal(expected[i], a))
	}
}

func TestExtractDatetimeClauses(t *testing.T) {
	dc1 := &autoopsproto.DatetimeClause{
		Time: 1000000001,
	}
	c1, err := ptypes.MarshalAny(dc1)
	require.NoError(t, err)
	dc2 := &autoopsproto.DatetimeClause{
		Time: 1000000002,
	}
	c2, err := ptypes.MarshalAny(dc2)
	require.NoError(t, err)
	oerc1 := &autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid1",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	c3, err := ptypes.MarshalAny(oerc1)
	require.NoError(t, err)
	autoOpsRule := &AutoOpsRule{&autoopsproto.AutoOpsRule{
		Id:        "id-0",
		FeatureId: "fid-0",
		Clauses:   []*autoopsproto.Clause{{Id: "c1", Clause: c1}, {Id: "c2", Clause: c2}, {Id: "c3", Clause: c3}},
	}}
	expected := []*autoopsproto.DatetimeClause{dc1, dc2}
	actual, err := autoOpsRule.ExtractDatetimeClauses()
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(actual))
	act1, has := actual["c1"]
	assert.True(t, proto.Equal(dc1, act1))
	act2, has := actual["c2"]
	assert.True(t, proto.Equal(dc2, act2))
	_, has = actual["c3"]
	assert.False(t, has)
}

func TestSetStopped(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.NotEqual(t, autoopsproto.AutoOpsStatus_STOPPED, aor.AutoOpsStatus)
	aor.SetStopped()
	assert.NotZero(t, aor.UpdatedAt)
	assert.Equal(t, autoopsproto.AutoOpsStatus_STOPPED, aor.AutoOpsStatus)
}

func TestSetFinished(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.NotEqual(t, autoopsproto.AutoOpsStatus_FINISHED, aor.AutoOpsStatus)
	aor.SetFinished()
	assert.NotZero(t, aor.UpdatedAt)
	assert.Equal(t, autoopsproto.AutoOpsStatus_FINISHED, aor.AutoOpsStatus)
}

func TestIsFinished(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.False(t, aor.IsFinished())

	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_FINISHED
	assert.True(t, aor.IsFinished())
}

func TestIsStopped(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.False(t, aor.IsStopped())

	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_STOPPED
	assert.True(t, aor.IsStopped())
}

func TestSetAutoOpsStatus(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.NotEqual(t, autoopsproto.AutoOpsStatus_FINISHED, aor.AutoOpsStatus)
	aor.SetAutoOpsStatus(autoopsproto.AutoOpsStatus_FINISHED)
	assert.Equal(t, autoopsproto.AutoOpsStatus_FINISHED, aor.AutoOpsStatus)
}

func TestHasEventRateOps(t *testing.T) {
	t.Parallel()
	dc1 := &autoopsproto.DatetimeClause{
		Time: 1000000001,
	}
	c1, err := ptypes.MarshalAny(dc1)
	require.NoError(t, err)
	dc2 := &autoopsproto.DatetimeClause{
		Time: 1000000002,
	}
	c2, err := ptypes.MarshalAny(dc2)
	require.NoError(t, err)

	autoOpsRule := &AutoOpsRule{&autoopsproto.AutoOpsRule{
		Id:        "id-0",
		FeatureId: "fid-0",
		Clauses:   []*autoopsproto.Clause{{Clause: c1}, {Clause: c2}},
	}}

	hasEventRateOps, err := autoOpsRule.HasEventRateOps()
	require.NoError(t, err)
	assert.False(t, hasEventRateOps)

	oerc := &autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid1",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	_, err = autoOpsRule.AddOpsEventRateClause(oerc)
	require.NoError(t, err)
	hasEventRateOps2, err := autoOpsRule.HasEventRateOps()
	assert.True(t, hasEventRateOps2)
}

func TestHasScheduleOps(t *testing.T) {
	t.Parallel()
	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{
			{
				GoalId:          "goalid01",
				MinCount:        10,
				ThreadsholdRate: 0.5,
				Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
				ActionType:      autoopsproto.ActionType_DISABLE,
			},
		},
		[]*autoopsproto.DatetimeClause{},
	)
	require.NoError(t, err)
	hasDateTimeOps, err := aor.HasScheduleOps()
	require.NoError(t, err)
	assert.False(t, hasDateTimeOps)

	ac1 := &autoopsproto.DatetimeClause{
		Time: 5,
	}
	_, err = aor.AddDatetimeClause(ac1)

	hasDateTimeOps2, err := aor.HasScheduleOps()
	require.NoError(t, err)
	assert.True(t, hasDateTimeOps2)
}

func TestUnmarshalOpsEventRateClause(t *testing.T) {
	erc := &autoopsproto.OpsEventRateClause{
		GoalId:          "goalid01",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
		ActionType:      autoopsproto.ActionType_DISABLE,
	}
	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{erc},
		[]*autoopsproto.DatetimeClause{},
	)
	require.NoError(t, err)
	eventRateClause, err := aor.unmarshalOpsEventRateClause(aor.Clauses[0])
	require.NoError(t, err)
	assert.Equal(t, erc.GoalId, eventRateClause.GoalId)
	assert.Equal(t, erc.MinCount, eventRateClause.MinCount)
	assert.Equal(t, erc.ActionType, eventRateClause.ActionType)
	assert.Equal(t, erc.ThreadsholdRate, eventRateClause.ThreadsholdRate)
	assert.Equal(t, erc.Operator, eventRateClause.Operator)
}

func TestUnmarshalDatetimeClause(t *testing.T) {
	dtc := &autoopsproto.DatetimeClause{
		Time: 5,
	}
	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{dtc},
	)
	require.NoError(t, err)
	dataTimeClause, err := aor.unmarshalDatetimeClause(aor.Clauses[0])
	require.NoError(t, err)
	assert.Equal(t, dtc.Time, dataTimeClause.Time)
}

func TestUpdateAutoOpsRule(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc                      string
		setup                     func() *AutoOpsRule
		updateOpsEventRateClauses []*autoopsproto.OpsEventRateClauseChange
		updateDatetimeClauses     []*autoopsproto.DatetimeClauseChange
		expected                  func() *AutoOpsRule
		expectedErr               error
	}{
		{
			desc: "Error clause empty",
			setup: func() *AutoOpsRule {
				aor, err := NewAutoOpsRule(
					"feature-id",
					autoopsproto.OpsType_EVENT_RATE,
					[]*autoopsproto.OpsEventRateClause{
						{
							GoalId:          "goal-01",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
					[]*autoopsproto.DatetimeClause{},
				)
				require.NoError(t, err)
				aor.Clauses[0].Id = "id-0"
				return aor
			},
			updateOpsEventRateClauses: []*autoopsproto.OpsEventRateClauseChange{
				{
					Id:         "id-0",
					ChangeType: autoopsproto.ChangeType_DELETE,
				},
			},
			updateDatetimeClauses: []*autoopsproto.DatetimeClauseChange{},
			expectedErr:           errClauseEmpty,
		},
		{
			desc: "Update OpsEventRateClause",
			setup: func() *AutoOpsRule {
				aor, err := NewAutoOpsRule(
					"feature-id",
					autoopsproto.OpsType_EVENT_RATE,
					[]*autoopsproto.OpsEventRateClause{
						{
							GoalId:          "goal-01",
							MinCount:        10,
							ThreadsholdRate: 0.5,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
					[]*autoopsproto.DatetimeClause{},
				)
				require.NoError(t, err)
				aor.Clauses[0].Id = "id-0"
				return aor
			},
			updateOpsEventRateClauses: []*autoopsproto.OpsEventRateClauseChange{
				{
					Id: "id-0",
					Clause: &autoopsproto.OpsEventRateClause{
						GoalId:          "goal-02",
						MinCount:        20,
						ThreadsholdRate: 0.6,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
					ChangeType: autoopsproto.ChangeType_UPDATE,
				},
			},
			updateDatetimeClauses: []*autoopsproto.DatetimeClauseChange{},
			expected: func() *AutoOpsRule {
				aor, err := NewAutoOpsRule(
					"feature-id",
					autoopsproto.OpsType_EVENT_RATE,
					[]*autoopsproto.OpsEventRateClause{
						{
							GoalId:          "goal-02",
							MinCount:        20,
							ThreadsholdRate: 0.6,
							Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
						},
					},
					[]*autoopsproto.DatetimeClause{},
				)
				require.NoError(t, err)
				aor.Clauses[0].Id = "id-0"
				return aor
			},
			expectedErr: nil,
		},
		{
			desc: "Update DatetimeClause",
			setup: func() *AutoOpsRule {
				aor, err := NewAutoOpsRule(
					"feature-id",
					autoopsproto.OpsType_SCHEDULE,
					[]*autoopsproto.OpsEventRateClause{},
					[]*autoopsproto.DatetimeClause{
						{
							Time:       1000000001,
							ActionType: autoopsproto.ActionType_ENABLE,
						},
					},
				)
				require.NoError(t, err)
				aor.Clauses[0].Id = "id-0"
				return aor
			},
			updateOpsEventRateClauses: []*autoopsproto.OpsEventRateClauseChange{},
			updateDatetimeClauses: []*autoopsproto.DatetimeClauseChange{
				{
					Id: "id-0",
					Clause: &autoopsproto.DatetimeClause{
						Time:       1000000002,
						ActionType: autoopsproto.ActionType_ENABLE,
					},
					ChangeType: autoopsproto.ChangeType_UPDATE,
				},
			},
			expected: func() *AutoOpsRule {
				aor, err := NewAutoOpsRule(
					"feature-id",
					autoopsproto.OpsType_SCHEDULE,
					[]*autoopsproto.OpsEventRateClause{},
					[]*autoopsproto.DatetimeClause{
						{
							Time:       1000000002,
							ActionType: autoopsproto.ActionType_ENABLE,
						},
					},
				)
				require.NoError(t, err)
				aor.Clauses[0].Id = "id-0"
				return aor
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			aor := p.setup()
			_, err := aor.Update(nil, p.updateOpsEventRateClauses, p.updateDatetimeClauses)
			if p.expectedErr != nil {
				require.Equal(t, p.expectedErr, err)
				return
			}
			expected := p.expected()

			assert.Equal(t, expected.FeatureId, aor.FeatureId)
			assert.Equal(t, expected.Clauses, aor.Clauses)
			assert.Equal(t, expected.OpsType, aor.OpsType)
		})
	}
}

func newDateTimeClause(t *testing.T, c *autoopsproto.DatetimeClause) *autoopsproto.Clause {
	clause, err := anypb.New(c)
	require.NoError(t, err)
	newClause := &autoopsproto.Clause{
		Clause:     clause,
		ActionType: c.ActionType,
	}
	return newClause
}

func newEventRateClause(t *testing.T, c *autoopsproto.OpsEventRateClause) *autoopsproto.Clause {
	clause, err := anypb.New(c)
	require.NoError(t, err)
	newClause := &autoopsproto.Clause{
		Clause:     clause,
		ActionType: c.ActionType,
	}
	return newClause
}

// createRecurringRule is a helper that creates a rule with a single recurring
// clause and pre-sets it to the given ExecutionCount / LastExecutedAt / NextExecutionAt.
func createRecurringRule(
	t *testing.T,
	timeOfDay int64,
	recurrence *autoopsproto.RecurrenceRule,
	executionCount int32,
	lastExecutedAt int64,
	nextExecutionAt int64,
) *AutoOpsRule {
	t.Helper()
	dc := &autoopsproto.DatetimeClause{
		Time:       timeOfDay,
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: recurrence,
	}
	rule, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		nil,
		[]*autoopsproto.DatetimeClause{dc},
	)
	require.NoError(t, err)
	require.Len(t, rule.Clauses, 1)

	if executionCount > 0 || lastExecutedAt > 0 || nextExecutionAt > 0 {
		dtClauses, err := rule.ExtractDatetimeClauses()
		require.NoError(t, err)
		dtc := dtClauses[rule.Clauses[0].Id]
		require.NotNil(t, dtc)
		dtc.ExecutionCount = executionCount
		dtc.LastExecutedAt = lastExecutedAt
		dtc.NextExecutionAt = nextExecutionAt
		updatedAny, err := ptypes.MarshalAny(dtc)
		require.NoError(t, err)
		rule.Clauses[0].Clause = updatedAny
	}
	return rule
}

// TestChangeDatetimeClause_RecalculatesNextExecution tests that updating
// scheduling-relevant fields on an already-executed recurring clause triggers
// recalculation of NextExecutionAt.
func TestChangeDatetimeClause_RecalculatesNextExecution(t *testing.T) {
	t.Parallel()

	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	baseStartDate := time.Date(2026, 2, 9, 0, 0, 0, 0, jst)
	baseRecurrence := func() *autoopsproto.RecurrenceRule {
		return &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1}, // Monday
			Timezone:   "Asia/Tokyo",
			StartDate:  baseStartDate.Unix(),
		}
	}

	// Existing state: executed once on Monday Feb 9 10:00 JST, next is Feb 16 10:00 JST
	existingNextExec := time.Date(2026, 2, 16, 10, 0, 0, 0, jst).Unix()
	existingLastExec := time.Date(2026, 2, 9, 10, 0, 1, 0, jst).Unix()

	tests := []struct {
		desc               string
		existingTimeOfDay  int64
		existingRecurrence *autoopsproto.RecurrenceRule
		existingExecCount  int32
		existingLastExec   int64
		existingNextExec   int64
		updatedTimeOfDay   int64
		updatedActionType  autoopsproto.ActionType
		updatedRecurrence  *autoopsproto.RecurrenceRule
		expectRecalculated bool
		expectNextExecZero bool
		expectExecCount    int32
		expectLastExecAt   int64
	}{
		{
			desc:               "time-of-day changed: recalculates nextExecutionAt",
			existingTimeOfDay:  75300, // 20:55
			existingRecurrence: baseRecurrence(),
			existingExecCount:  1,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   68400, // 19:00
			updatedActionType:  autoopsproto.ActionType_ENABLE,
			updatedRecurrence:  baseRecurrence(),
			expectRecalculated: true,
			expectExecCount:    1,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:               "same time: preserves nextExecutionAt",
			existingTimeOfDay:  36000,
			existingRecurrence: baseRecurrence(),
			existingExecCount:  1,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   36000,
			updatedActionType:  autoopsproto.ActionType_ENABLE,
			updatedRecurrence:  baseRecurrence(),
			expectRecalculated: false,
			expectExecCount:    1,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:               "only actionType changed: preserves nextExecutionAt",
			existingTimeOfDay:  36000,
			existingRecurrence: baseRecurrence(),
			existingExecCount:  1,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   36000,
			updatedActionType:  autoopsproto.ActionType_DISABLE,
			updatedRecurrence:  baseRecurrence(),
			expectRecalculated: false,
			expectExecCount:    1,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:               "time changed on clause with endDate already passed: nextExec becomes 0",
			existingTimeOfDay:  36000,
			existingRecurrence: baseRecurrence(),
			existingExecCount:  3,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   32400, // 9:00 AM
			updatedActionType:  autoopsproto.ActionType_ENABLE,
			updatedRecurrence: func() *autoopsproto.RecurrenceRule {
				r := baseRecurrence()
				r.EndDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix() // already passed
				return r
			}(),
			expectRecalculated: true,
			expectNextExecZero: true,
			expectExecCount:    3,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:               "maxOccurrences reduced below current count: nextExec becomes 0",
			existingTimeOfDay:  36000,
			existingRecurrence: baseRecurrence(),
			existingExecCount:  5,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   36000,
			updatedActionType:  autoopsproto.ActionType_ENABLE,
			updatedRecurrence: func() *autoopsproto.RecurrenceRule {
				r := baseRecurrence()
				r.MaxOccurrences = 3 // less than current count of 5
				return r
			}(),
			expectRecalculated: true,
			expectNextExecZero: true,
			expectExecCount:    5,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:               "daysOfWeek changed: recalculates nextExecutionAt",
			existingTimeOfDay:  36000,
			existingRecurrence: baseRecurrence(), // Monday only
			existingExecCount:  1,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   36000,
			updatedActionType:  autoopsproto.ActionType_ENABLE,
			updatedRecurrence: func() *autoopsproto.RecurrenceRule {
				r := baseRecurrence()
				r.DaysOfWeek = []int32{3, 5} // Wed, Fri instead of Mon
				return r
			}(),
			expectRecalculated: true,
			expectExecCount:    1,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:               "timezone changed: recalculates nextExecutionAt",
			existingTimeOfDay:  36000,
			existingRecurrence: baseRecurrence(),
			existingExecCount:  1,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   36000,
			updatedActionType:  autoopsproto.ActionType_ENABLE,
			updatedRecurrence: func() *autoopsproto.RecurrenceRule {
				r := baseRecurrence()
				r.Timezone = "America/New_York"
				return r
			}(),
			expectRecalculated: true,
			expectExecCount:    1,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:               "frequency changed from weekly to daily: recalculates",
			existingTimeOfDay:  36000,
			existingRecurrence: baseRecurrence(),
			existingExecCount:  1,
			existingLastExec:   existingLastExec,
			existingNextExec:   existingNextExec,
			updatedTimeOfDay:   36000,
			updatedActionType:  autoopsproto.ActionType_ENABLE,
			updatedRecurrence: &autoopsproto.RecurrenceRule{
				Frequency: autoopsproto.RecurrenceRule_DAILY,
				Timezone:  "Asia/Tokyo",
				StartDate: baseStartDate.Unix(),
			},
			expectRecalculated: true,
			expectExecCount:    1,
			expectLastExecAt:   existingLastExec,
		},
		{
			desc:              "never-executed clause: initializes via InitializeRecurringClause",
			existingTimeOfDay: 36000,
			existingRecurrence: func() *autoopsproto.RecurrenceRule {
				r := baseRecurrence()
				r.StartDate = time.Now().Add(24 * time.Hour).Unix()
				return r
			}(),
			existingExecCount: 0,
			existingLastExec:  0,
			existingNextExec:  0,
			updatedTimeOfDay:  32400,
			updatedActionType: autoopsproto.ActionType_ENABLE,
			updatedRecurrence: func() *autoopsproto.RecurrenceRule {
				r := baseRecurrence()
				r.StartDate = time.Now().Add(24 * time.Hour).Unix()
				return r
			}(),
			expectRecalculated: true,
			expectExecCount:    0,
			expectLastExecAt:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			rule := createRecurringRule(
				t,
				tt.existingTimeOfDay,
				tt.existingRecurrence,
				tt.existingExecCount,
				tt.existingLastExec,
				tt.existingNextExec,
			)
			clauseID := rule.Clauses[0].Id

			updatedDtClause := &autoopsproto.DatetimeClause{
				Time:       tt.updatedTimeOfDay,
				ActionType: tt.updatedActionType,
				Recurrence: tt.updatedRecurrence,
			}

			err := rule.ChangeDatetimeClause(clauseID, updatedDtClause)
			require.NoError(t, err)

			dtClauses, err := rule.ExtractDatetimeClauses()
			require.NoError(t, err)
			dtClause := dtClauses[clauseID]
			require.NotNil(t, dtClause)

			assert.Equal(t, tt.expectExecCount, dtClause.ExecutionCount,
				"ExecutionCount should be preserved")
			assert.Equal(t, tt.expectLastExecAt, dtClause.LastExecutedAt,
				"LastExecutedAt should be preserved")

			if tt.expectNextExecZero {
				assert.Equal(t, int64(0), dtClause.NextExecutionAt,
					"NextExecutionAt should be 0 (exhausted)")
			} else if tt.expectRecalculated {
				if tt.existingExecCount > 0 {
					assert.NotEqual(t, tt.existingNextExec, dtClause.NextExecutionAt,
						"NextExecutionAt should have been recalculated")
				}
				assert.True(t, dtClause.NextExecutionAt > 0,
					"NextExecutionAt should be positive after recalculation")
			} else {
				assert.Equal(t, tt.existingNextExec, dtClause.NextExecutionAt,
					"NextExecutionAt should be preserved (no schedule change)")
			}

			assert.Equal(t, tt.updatedTimeOfDay, dtClause.Time,
				"Time should be updated")
			assert.Equal(t, tt.updatedActionType, dtClause.ActionType,
				"ActionType should be updated")
		})
	}
}

// TestChangeDatetimeClause_UpdateTimeOnExecutedClause_ViaUpdate tests the full
// Update() path (which uses copier.Copy) to verify recalculation survives the
// deep copy, simulating the actual API flow.
func TestChangeDatetimeClause_UpdateTimeOnExecutedClause_ViaUpdate(t *testing.T) {
	t.Parallel()

	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 9, 0, 0, 0, 0, jst)
	oldNextExec := time.Date(2026, 2, 16, 20, 55, 0, 0, jst).Unix()
	lastExec := time.Date(2026, 2, 9, 20, 55, 1, 0, jst).Unix()

	rule := createRecurringRule(
		t,
		75300, // 20:55
		&autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
		1,           // executed once
		lastExec,    // last executed at
		oldNextExec, // next execution at
	)
	clauseID := rule.Clauses[0].Id

	updated, err := rule.Update(nil, nil, []*autoopsproto.DatetimeClauseChange{
		{
			Id:         clauseID,
			ChangeType: autoopsproto.ChangeType_UPDATE,
			Clause: &autoopsproto.DatetimeClause{
				Time:       68400, // 19:00 (changed from 20:55)
				ActionType: autoopsproto.ActionType_ENABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{1},
					Timezone:   "Asia/Tokyo",
					StartDate:  startDate.Unix(),
				},
			},
		},
	})
	require.NoError(t, err)

	dtClauses, err := updated.ExtractDatetimeClauses()
	require.NoError(t, err)
	dtClause := dtClauses[clauseID]
	require.NotNil(t, dtClause)

	assert.Equal(t, int64(68400), dtClause.Time, "Time should be updated to 19:00")
	assert.Equal(t, int32(1), dtClause.ExecutionCount, "ExecutionCount should be preserved")
	assert.Equal(t, lastExec, dtClause.LastExecutedAt, "LastExecutedAt should be preserved")
	assert.NotEqual(t, oldNextExec, dtClause.NextExecutionAt,
		"NextExecutionAt should have been recalculated (not the old 20:55 timestamp)")
	assert.True(t, dtClause.NextExecutionAt > 0,
		"NextExecutionAt should be positive (recalculated for 19:00)")
}

// TestChangeDatetimeClause_MultiClause_OnlyTargetRecalculated verifies that
// updating one clause in a multi-clause rule only recalculates that clause,
// leaving the other clause's NextExecutionAt intact.
func TestChangeDatetimeClause_MultiClause_OnlyTargetRecalculated(t *testing.T) {
	t.Parallel()

	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 9, 0, 0, 0, 0, jst)
	recurrence := &autoopsproto.RecurrenceRule{
		Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
		DaysOfWeek: []int32{1},
		Timezone:   "Asia/Tokyo",
		StartDate:  startDate.Unix(),
	}

	enableClause := &autoopsproto.DatetimeClause{
		Time:       36000, // 10:00
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: recurrence,
	}
	disableClause := &autoopsproto.DatetimeClause{
		Time:       64800, // 18:00
		ActionType: autoopsproto.ActionType_DISABLE,
		Recurrence: recurrence,
	}

	rule, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		nil,
		[]*autoopsproto.DatetimeClause{enableClause, disableClause},
	)
	require.NoError(t, err)
	require.Len(t, rule.Clauses, 2)

	// Simulate both clauses having been executed once
	advanceTime := time.Date(2026, 2, 9, 18, 0, 1, 0, jst)
	err = rule.AdvanceRecurringClause(rule.Clauses[0].Id, advanceTime)
	require.NoError(t, err)
	err = rule.AdvanceRecurringClause(rule.Clauses[1].Id, advanceTime)
	require.NoError(t, err)

	// Save clause IDs before the update (sorting may reorder positions)
	enableClauseID := rule.Clauses[0].Id
	disableClauseID := rule.Clauses[1].Id

	dtClauses, err := rule.ExtractDatetimeClauses()
	require.NoError(t, err)
	disableNextExec := dtClauses[disableClauseID].NextExecutionAt

	// Update only the enable clause's time (same recurrence)
	err = rule.ChangeDatetimeClause(enableClauseID, &autoopsproto.DatetimeClause{
		Time:       32400, // Changed to 9:00
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: recurrence,
	})
	require.NoError(t, err)

	dtClauses, err = rule.ExtractDatetimeClauses()
	require.NoError(t, err)

	// Enable clause should have been recalculated
	enableDt := dtClauses[enableClauseID]
	require.NotNil(t, enableDt)
	assert.True(t, enableDt.NextExecutionAt > 0,
		"enable clause NextExecutionAt should be recalculated and positive")

	// Disable clause should be unchanged
	disableDt := dtClauses[disableClauseID]
	require.NotNil(t, disableDt)
	assert.Equal(t, disableNextExec, disableDt.NextExecutionAt,
		"disable clause NextExecutionAt should be unchanged")
}

// TestChangeDatetimeClause_EndDateMadeEarlier_ExhaustsClause tests that
// changing the end date to be in the past on an executed clause correctly
// sets NextExecutionAt to 0 (exhausted).
func TestChangeDatetimeClause_EndDateMadeEarlier_ExhaustsClause(t *testing.T) {
	t.Parallel()

	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 9, 0, 0, 0, 0, jst)
	oldNextExec := time.Date(2026, 2, 16, 10, 0, 0, 0, jst).Unix()
	lastExec := time.Date(2026, 2, 9, 10, 0, 1, 0, jst).Unix()

	rule := createRecurringRule(
		t,
		36000, // 10:00
		&autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
		1,           // executed once
		lastExec,    // last executed at
		oldNextExec, // next execution at
	)
	clauseID := rule.Clauses[0].Id

	// Update with same time but end date in the past
	err = rule.ChangeDatetimeClause(clauseID, &autoopsproto.DatetimeClause{
		Time:       36000,
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
			EndDate:    time.Date(2026, 2, 10, 0, 0, 0, 0, jst).Unix(), // already passed
		},
	})
	require.NoError(t, err)

	dtClauses, err := rule.ExtractDatetimeClauses()
	require.NoError(t, err)
	dtClause := dtClauses[clauseID]
	require.NotNil(t, dtClause)

	assert.Equal(t, int64(0), dtClause.NextExecutionAt,
		"NextExecutionAt should be 0 since endDate is in the past")
	assert.Equal(t, int32(1), dtClause.ExecutionCount, "ExecutionCount preserved")
	assert.Equal(t, lastExec, dtClause.LastExecutedAt, "LastExecutedAt preserved")
}
