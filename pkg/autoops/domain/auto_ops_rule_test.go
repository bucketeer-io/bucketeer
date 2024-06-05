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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestNewAutoOpsRule(t *testing.T) {
	t.Parallel()

	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_ENABLE_FEATURE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{Time: 0},
			{Time: 0, ActionType: autoopsproto.ActionType_ENABLE},
		},
	)
	require.NoError(t, err)

	assert.IsType(t, &AutoOpsRule{}, aor)
	assert.Equal(t, "feature-id", aor.FeatureId)
	assert.Equal(t, autoopsproto.OpsType_ENABLE_FEATURE, aor.OpsType)
	assert.NotZero(t, aor.Clauses)
	assert.Zero(t, aor.TriggeredAt)
	assert.NotZero(t, aor.CreatedAt)
	assert.NotZero(t, aor.UpdatedAt)
}

func TestNewAutoOpsRule_V2(t *testing.T) {
	patterns := []struct {
		featureId        string
		desc             string
		opsType          autoopsproto.OpsType
		datetimeClauses  []*autoopsproto.DatetimeClause
		eventRateClauses []*autoopsproto.OpsEventRateClause
	}{
		{
			desc:      "OpsType: Schedule",
			featureId: "feature-id",
			opsType:   autoopsproto.OpsType_SCHEDULE,
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{Time: 0, ActionType: autoopsproto.ActionType_ENABLE},
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
		},
		{
			desc:      "OpsType: EventRate",
			featureId: "feature-id",
			opsType:   autoopsproto.OpsType_EVENT_RATE,
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{Time: 0, ActionType: autoopsproto.ActionType_ENABLE},
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
		},
	}
	for _, p := range patterns {
		aor, err := NewAutoOpsRule_V2(
			p.featureId,
			p.opsType,
			p.eventRateClauses,
			p.datetimeClauses,
		)
		require.NoError(t, err)
		assert.IsType(t, &AutoOpsRule{}, aor)
		assert.Equal(t, p.featureId, aor.FeatureId)
		assert.Equal(t, p.opsType, aor.OpsType)
		assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)
		assert.NotZero(t, aor.CreatedAt)
		assert.NotZero(t, aor.UpdatedAt)
		assert.Zero(t, aor.TriggeredAt)
		assert.Zero(t, aor.StoppedAt)

		if aor.OpsType == autoopsproto.OpsType_EVENT_RATE {
			assert.Equal(t, len(p.eventRateClauses), len(aor.Clauses))
			for i, c := range aor.Clauses {
				eventRateClause, err := aor.unmarshalOpsEventRateClause(c)
				require.NoError(t, err)
				assert.Equal(t, p.eventRateClauses[i].GoalId, eventRateClause.GoalId)
				assert.Equal(t, p.eventRateClauses[i].MinCount, eventRateClause.MinCount)
				assert.Equal(t, p.eventRateClauses[i].ThreadsholdRate, eventRateClause.ThreadsholdRate)
				assert.Equal(t, p.eventRateClauses[i].Operator, eventRateClause.Operator)
				assert.Equal(t, p.eventRateClauses[i].ActionType, eventRateClause.ActionType)
			}
		}
		if aor.OpsType == autoopsproto.OpsType_SCHEDULE {
			assert.Equal(t, len(p.datetimeClauses), len(aor.Clauses))
			for i, c := range aor.Clauses {
				datetimeClause, err := aor.unmarshalDatetimeClause(c)
				require.NoError(t, err)
				assert.Equal(t, p.datetimeClauses[i].Time, datetimeClause.Time)
				assert.Equal(t, p.datetimeClauses[i].ActionType, datetimeClause.ActionType)
			}
		}
	}
}

func TestSetDeleted(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.NotEqual(t, autoopsproto.AutoOpsStatus_DELETED, aor.AutoOpsStatus)
	aor.SetDeleted()
	assert.Equal(t, autoopsproto.AutoOpsStatus_DELETED, aor.AutoOpsStatus)
	assert.NotZero(t, aor.UpdatedAt)
	assert.Zero(t, aor.StoppedAt)
	assert.Equal(t, true, aor.Deleted)
}

func TestSetTriggeredAt(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.SetTriggeredAt()
	assert.NotZero(t, aor.TriggeredAt)
}

func TestAlreadyTriggeredAt(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.False(t, aor.AlreadyTriggered())
	aor.SetTriggeredAt()
	assert.True(t, aor.AlreadyTriggered())
}

func TestSetOpsType(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_COMPLETED
	aor.SetOpsType(autoopsproto.OpsType_SCHEDULE)
	assert.Equal(t, autoopsproto.OpsType_SCHEDULE, aor.OpsType)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)
	assert.Zero(t, aor.TriggeredAt)
}

func TestAddOpsEventRateClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_COMPLETED
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
	assert.Zero(t, aor.TriggeredAt)

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
	aor.TriggeredAt = 1
	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_COMPLETED
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
	assert.Zero(t, aor.TriggeredAt)
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
	aor.TriggeredAt = 1
	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_COMPLETED
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
	assert.Zero(t, aor.TriggeredAt)
	assert.Equal(t, c.GoalId, eventRateClause.GoalId)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)
}

func TestChangeDatetimeClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_COMPLETED
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
	assert.Zero(t, aor.TriggeredAt)
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
	aor.TriggeredAt = 1
	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_COMPLETED
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
	assert.Zero(t, aor.TriggeredAt)
	assert.Equal(t, addClause.Id, aor.Clauses[0].Id)
	assert.Equal(t, autoopsproto.AutoOpsStatus_WAITING, aor.AutoOpsStatus)
}

func createAutoOpsRule(t *testing.T) *AutoOpsRule {
	aor, err := NewAutoOpsRule_V2(
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
	assert.Zero(t, aor.StoppedAt)
	assert.NotEqual(t, autoopsproto.AutoOpsStatus_STOPPED, aor.AutoOpsStatus)
	aor.SetStopped()
	assert.NotZero(t, aor.UpdatedAt)
	assert.NotZero(t, aor.StoppedAt)
	assert.Equal(t, autoopsproto.AutoOpsStatus_STOPPED, aor.AutoOpsStatus)
}

func TestSetCompleted(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.NotEqual(t, autoopsproto.AutoOpsStatus_COMPLETED, aor.AutoOpsStatus)
	aor.SetCompleted()
	assert.NotZero(t, aor.UpdatedAt)
	assert.Equal(t, autoopsproto.AutoOpsStatus_COMPLETED, aor.AutoOpsStatus)
}

func TestHasExecuteClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.True(t, aor.HasExecuteClause())

	aor.AutoOpsStatus = autoopsproto.AutoOpsStatus_COMPLETED
	assert.False(t, aor.HasExecuteClause())

	dc := &autoopsproto.DatetimeClause{
		Time: 1,
	}
	_, err := aor.AddDatetimeClause(dc)
	require.NoError(t, err)
	assert.True(t, aor.HasExecuteClause())
}

func TestSetAutoOpsStatus(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	assert.NotEqual(t, autoopsproto.AutoOpsStatus_COMPLETED, aor.AutoOpsStatus)
	aor.SetAutoOpsStatus(autoopsproto.AutoOpsStatus_COMPLETED)
	assert.Equal(t, autoopsproto.AutoOpsStatus_COMPLETED, aor.AutoOpsStatus)
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
