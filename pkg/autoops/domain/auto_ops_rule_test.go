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

package domain

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
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
		updateOpsEventRateClauses []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateOpsEventRateClause
		updateDatetimeClauses     []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateDatetimeClause
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
			updateOpsEventRateClauses: []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateOpsEventRateClause{
				{
					Id:      "id-0",
					Deleted: wrapperspb.Bool(true),
				},
			},
			updateDatetimeClauses: []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateDatetimeClause{},
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
			updateOpsEventRateClauses: []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateOpsEventRateClause{
				{
					Id: "id-0",
					Clause: &autoopsproto.OpsEventRateClause{
						GoalId:          "goal-02",
						MinCount:        20,
						ThreadsholdRate: 0.6,
						Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					},
				},
			},
			updateDatetimeClauses: []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateDatetimeClause{},
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
			updateOpsEventRateClauses: []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateOpsEventRateClause{},
			updateDatetimeClauses: []*autoopsproto.UpdateAutoOpsRuleRequest_UpdateDatetimeClause{
				{
					Id: "id-0",
					Clause: &autoopsproto.DatetimeClause{
						Time:       1000000002,
						ActionType: autoopsproto.ActionType_ENABLE,
					},
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
			_, err := aor.Update(p.updateOpsEventRateClauses, p.updateDatetimeClauses)
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
