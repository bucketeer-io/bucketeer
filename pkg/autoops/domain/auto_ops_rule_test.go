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
	aor := createAutoOpsRule(t)
	assert.IsType(t, &AutoOpsRule{}, aor)
	assert.Equal(t, "feature-id", aor.FeatureId)
	assert.Equal(t, autoopsproto.OpsType_ENABLE_FEATURE, aor.OpsType)
	assert.NotZero(t, aor.Clauses)
	assert.Zero(t, aor.TriggeredAt)
	assert.NotZero(t, aor.CreatedAt)
	assert.NotZero(t, aor.UpdatedAt)
}

func TestSetDeleted(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.SetDeleted()
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
	aor.SetOpsType(autoopsproto.OpsType_DISABLE_FEATURE)
	assert.Equal(t, autoopsproto.OpsType_DISABLE_FEATURE, aor.OpsType)
	assert.Zero(t, aor.TriggeredAt)
}

func TestAddOpsEventRateClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
	l := len(aor.Clauses)
	c := &autoopsproto.OpsEventRateClause{
		GoalId:          "goalid01",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	clause, err := aor.AddOpsEventRateClause(c)
	require.NoError(t, err)
	assert.NotNil(t, clause)
	assert.NotEmpty(t, aor.Clauses[l].Id)
	_, err = aor.unmarshalOpsEventRateClause(aor.Clauses[l])
	require.NoError(t, err)
	assert.Zero(t, aor.TriggeredAt)
}

func TestAddDatetimeClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
	l := len(aor.Clauses)
	c := &autoopsproto.DatetimeClause{
		Time: 1000000001,
	}
	clause, err := aor.AddDatetimeClause(c)
	require.NoError(t, err)
	assert.NotNil(t, clause)
	assert.NotEmpty(t, aor.Clauses[l].Id)
	dc, err := aor.unmarshalDatetimeClause(aor.Clauses[l])
	require.NoError(t, err)
	assert.Equal(t, c.Time, dc.Time)
	assert.Zero(t, aor.TriggeredAt)
}

func TestChangeOpsEventRateClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
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
	_, err = aor.unmarshalOpsEventRateClause(aor.Clauses[0])
	require.NoError(t, err)
	assert.Zero(t, aor.TriggeredAt)
}

func TestDatetimeClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
	l := len(aor.Clauses)
	c := &autoopsproto.DatetimeClause{
		Time: 1000000001,
	}
	err := aor.ChangeDatetimeClause(aor.Clauses[0].Id, c)
	require.NoError(t, err)
	assert.Equal(t, l, len(aor.Clauses))
	dc, err := aor.unmarshalDatetimeClause(aor.Clauses[0])
	require.NoError(t, err)
	assert.Equal(t, c.Time, dc.Time)
	assert.Zero(t, aor.TriggeredAt)
}

func TestDeleteClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)
	aor.TriggeredAt = 1
	l := len(aor.Clauses)
	c := &autoopsproto.OpsEventRateClause{
		GoalId:          "goalid01",
		MinCount:        10,
		ThreadsholdRate: 0.5,
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	}
	_, err := aor.AddOpsEventRateClause(c)
	require.NoError(t, err)
	err = aor.DeleteClause(aor.Clauses[l].Id)
	require.NoError(t, err)
	assert.Equal(t, l, len(aor.Clauses))
	assert.Zero(t, aor.TriggeredAt)
}

func createAutoOpsRule(t *testing.T) *AutoOpsRule {
	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_ENABLE_FEATURE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{Time: 0},
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
		Clauses:   []*autoopsproto.Clause{{Clause: c1}, {Clause: c2}, {Clause: c3}},
	}}
	expected := []*autoopsproto.DatetimeClause{dc1, dc2}
	actual, err := autoOpsRule.ExtractDatetimeClauses()
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(actual))
	for i, a := range actual {
		assert.True(t, proto.Equal(expected[i], a))
	}
}
