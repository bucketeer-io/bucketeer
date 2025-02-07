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
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	"github.com/bucketeer-io/bucketeer/proto/feature"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func makeFeature(id string) *Feature {
	return &Feature{
		Feature: &proto.Feature{
			Id:            id,
			Name:          "test feature",
			Version:       1,
			Enabled:       true,
			CreatedAt:     time.Now().Unix(),
			VariationType: feature.Feature_STRING,
			Variations: []*proto.Variation{
				{
					Id:          "variation-A",
					Value:       "A",
					Name:        "Variation A",
					Description: "Thing does A",
				},
				{
					Id:          "variation-B",
					Value:       "B",
					Name:        "Variation B",
					Description: "Thing does B",
				},
				{
					Id:          "variation-C",
					Value:       "C",
					Name:        "Variation C",
					Description: "Thing does C",
				},
			},
			Targets: []*proto.Target{
				{
					Variation: "variation-A",
					Users: []string{
						"user1",
					},
				},
				{
					Variation: "variation-B",
					Users: []string{
						"user2",
					},
				},
				{
					Variation: "variation-C",
					Users: []string{
						"user3",
					},
				},
			},
			Rules: []*proto.Rule{
				{
					Id: "rule-1",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*proto.Clause{
						{
							Id:        "clause-1",
							Attribute: "name",
							Operator:  proto.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
				{
					Id: "rule-2",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{
							Variation: "variation-B",
						},
					},
					Clauses: []*proto.Clause{
						{
							Id:        "clause-2",
							Attribute: "name",
							Operator:  proto.Clause_EQUALS,
							Values: []string{
								"user3",
								"user4",
							},
						},
					},
				},
			},
			DefaultStrategy: &proto.Strategy{
				Type: proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{
					Variation: "variation-B",
				},
			},
		},
	}
}

func TestNewFeature(t *testing.T) {
	id := "id"
	name := "name"
	description := "description"
	variations := []*proto.Variation{
		{
			Value:       "A",
			Name:        "Variation A",
			Description: "Thing does A",
		},
		{
			Value:       "B",
			Name:        "Variation B",
			Description: "Thing does B",
		},
		{
			Value:       "C",
			Name:        "Variation C",
			Description: "Thing does C",
		},
	}
	variationType := feature.Feature_STRING
	tags := []string{"android", "ios", "web"}
	defaultOnVariationIndex := 0
	defaultOffVariationIndex := 2
	maintainer := "bucketeer@example.com"
	f, err := NewFeature(
		id,
		name,
		description,
		variationType,
		variations,
		tags,
		defaultOnVariationIndex,
		defaultOffVariationIndex,
		maintainer,
	)
	strategy := &feature.Strategy{
		Type:          feature.Strategy_FIXED,
		FixedStrategy: &feature.FixedStrategy{Variation: f.Variations[defaultOnVariationIndex].Id},
	}
	assert.NoError(t, err)
	assert.Equal(t, id, f.Id)
	assert.Equal(t, name, f.Name)
	assert.Equal(t, description, f.Description)
	for i := range variations {
		assert.Equal(t, variations[i].Name, f.Variations[i].Name)
		assert.Equal(t, variations[i].Description, f.Variations[i].Description)
	}
	assert.Equal(t, tags, f.Tags)
	assert.Equal(t, tags, f.Tags)
	assert.Equal(t, f.Variations[defaultOffVariationIndex].Id, f.OffVariation)
	assert.Equal(t, strategy, f.DefaultStrategy)
	assert.Equal(t, maintainer, f.Maintainer)
}

func TestAddVariation(t *testing.T) {
	createFeature := func() *Feature {
		f := makeFeature("test-feature")
		f.Rules = []*proto.Rule{
			{
				Id: "rule-0",
				Strategy: &proto.Strategy{
					Type: proto.Strategy_ROLLOUT,
					RolloutStrategy: &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
						{Variation: "A", Weight: 70000},
						{Variation: "B", Weight: 30000},
						{Variation: "C", Weight: 0},
					},
					},
				},
				Clauses: []*proto.Clause{
					{
						Id:        "clause-0",
						Attribute: "name",
						Operator:  proto.Clause_EQUALS,
						Values: []string{
							"user1",
						},
					},
				},
			},
		}
		f.DefaultStrategy = &proto.Strategy{
			Type: proto.Strategy_ROLLOUT,
			RolloutStrategy: &feature.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
				{Variation: "A", Weight: 70000},
				{Variation: "B", Weight: 30000},
				{Variation: "C", Weight: 0},
			}},
		}
		return f
	}
	patterns := []struct {
		desc  string
		input string
	}{
		{
			desc:  "Add D",
			input: "variation-D",
		},
	}
	for _, p := range patterns {
		f := createFeature()
		f.AddVariation(p.input, p.input, "", "")
		assert.Equal(t, p.input, f.Targets[3].Variation)
		assert.Equal(t, p.input, f.Rules[0].Strategy.RolloutStrategy.Variations[3].Variation)
		assert.Equal(t, p.input, f.DefaultStrategy.RolloutStrategy.Variations[3].Variation)
	}
}

func TestRename(t *testing.T) {
	f := makeFeature("test-feature")
	name := "new name"
	f.Rename(name)
	assert.Equal(t, name, f.Name)
}

func TestChangeDescription(t *testing.T) {
	f := makeFeature("test-feature")
	desc := "new desc"
	f.ChangeDescription(desc)
	assert.Equal(t, desc, f.Description)
}

func TestAddTag(t *testing.T) {
	tag := "test-tag"
	f := makeFeature("test-feature")
	if len(f.Tags) > 0 {
		t.Fatalf("Failed to add tag. It should be empty before add a tag: %v", f.Tags)
	}
	f.AddTag(tag)
	if len(f.Tags) == 0 {
		t.Fatal("Failed to add tag. Tags is empty.")
	}
	if len(f.Tags) != 1 {
		t.Fatalf("Failed to add tag. Tags has more than one element: %v", f.Tags)
	}
	if f.Tags[0] != tag {
		t.Fatalf("Failed to add tag. Tag does not match, current: %s, target: %s", f.Tags[0], tag)
	}
}

func TestRemoveTag(t *testing.T) {
	tag1 := "test-tag1"
	tag2 := "test-tag2"
	f := makeFeature("test-feature")
	f.AddTag(tag1)
	f.AddTag(tag2)
	f.RemoveTag(tag1)
	f.RemoveTag("not-exist-tag")
	if f.Tags[0] == tag1 {
		t.Fatalf("Failed to remove tag %s. Tags: %v", tag1, f.Tags)
	}
	if len(f.Tags) != 1 {
		t.Fatalf("Failed to remove tag. It should remove only 1: %v", f.Tags)
	}
	f.RemoveTag(tag2)
	if len(f.Tags) != 0 {
		t.Fatalf("Failed to remove tag. It should remove: %s. Actual: %v", tag2, f.Tags)
	}
}

func TestEnable(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		origin      bool
		expectedErr error
	}{
		{
			origin:      true,
			expectedErr: ErrAlreadyEnabled,
		},
		{
			origin:      false,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		f := makeFeature("test-feature")
		f.Enabled = p.origin
		err := f.Enable()
		assert.Equal(t, p.expectedErr, err)
		assert.True(t, f.Enabled)
	}
}

func TestDisable(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		origin      bool
		expectedErr error
	}{
		{
			origin:      false,
			expectedErr: ErrAlreadyDisabled,
		},
		{
			origin:      true,
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		f := makeFeature("test-feature")
		f.Enabled = p.origin
		err := f.Disable()
		assert.Equal(t, p.expectedErr, err)
		assert.False(t, f.Enabled)
	}
}

func TestArchive(t *testing.T) {
	f := makeFeature("test-feature")
	f.Archived = false
	f.Archive()
	assert.True(t, f.Archived)
}

func TestUnarchive(t *testing.T) {
	t.Parallel()
	f := makeFeature("test-feature")
	f.Archive()
	assert.True(t, f.Archived)
	f.Unarchive()
	assert.False(t, f.Archived)
}

func TestDelete(t *testing.T) {
	f := makeFeature("test-feature")
	f.Deleted = false
	f.Delete()
	assert.True(t, f.Deleted)
}

func TestAddUserToVariation(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		variation   string
		user        string
		idx         int
		expectedLen int
		expectedErr error
	}{
		{
			variation:   "",
			user:        "",
			idx:         -1,
			expectedLen: -1,
			expectedErr: errTargetNotFound,
		},
		{
			variation:   "variation-A",
			user:        "user1",
			idx:         0,
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			variation:   "variation-A",
			user:        "newUser1",
			idx:         0,
			expectedLen: 2,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.AddUserToVariation(p.variation, p.user)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.expectedLen, len(f.Targets[p.idx].Users))
		}
	}
}

func TestRemoveUserFromVariation(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		variation   string
		user        string
		idx         int
		expectedLen int
		expectedErr error
	}{
		{
			variation:   "",
			user:        "",
			idx:         -1,
			expectedLen: -1,
			expectedErr: errTargetNotFound,
		},
		{
			variation:   "variation-A",
			user:        "newUser1",
			idx:         -1,
			expectedLen: -1,
			expectedErr: errValueNotFound,
		},
		{
			variation:   "variation-A",
			user:        "user1",
			idx:         0,
			expectedLen: 0,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.RemoveUserFromVariation(p.variation, p.user)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.expectedLen, len(f.Targets[p.idx].Users))
		}
	}
}

func TestAddFixedStrategyRule(t *testing.T) {
	f := makeFeature("test-feature")
	id1, _ := uuid.NewUUID()
	patterns := []struct {
		desc        string
		id          string
		strategy    *proto.Strategy
		clauses     []*feature.Clause
		expectedErr bool
	}{
		{
			desc:     "fail: add rule with nil strategy",
			id:       "rule-2",
			strategy: nil,
			clauses: []*feature.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  feature.Clause_EQUALS,
					Values: []string{
						"user1",
					},
				},
			},
			expectedErr: true,
		},
		{
			desc: "fail: add rule with nil clauses",
			id:   "rule-3",
			strategy: &proto.Strategy{
				Type:          proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{Variation: ""},
			},
			clauses: []*feature.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  feature.Clause_EQUALS,
					Values: []string{
						"user1",
					},
				},
			},
			expectedErr: true,
		},
		{
			desc: "success",
			id:   "rule-3",
			strategy: &proto.Strategy{
				Type:          proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{Variation: f.Variations[0].Id},
			},
			clauses: []*feature.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  feature.Clause_EQUALS,
					Values: []string{
						"user1",
					},
				},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		rule := &proto.Rule{
			Id:       p.id,
			Strategy: p.strategy,
			Clauses:  p.clauses,
		}
		err := f.AddRule(rule)
		assert.Equal(t, p.expectedErr, err != nil, "%s", p.desc)
		if !p.expectedErr {
			assert.Equal(t, rule, f.Rules[2], "%s", p.desc)
		}
	}
}

func TestAddRolloutStrategyRule(t *testing.T) {
	f := makeFeature("test-feature")
	id1, _ := uuid.NewUUID()
	patterns := []struct {
		desc        string
		rule        *proto.Rule
		expectedErr bool
	}{
		{
			desc: "fail: rule already exists",
			rule: &feature.Rule{
				Id:       "rule-2",
				Strategy: nil,
			},
			expectedErr: true,
		},
		{
			desc: "fail: variation not found",
			rule: &feature.Rule{
				Id: "rule-3",
				Strategy: &proto.Strategy{
					Type: proto.Strategy_ROLLOUT,
					RolloutStrategy: &proto.RolloutStrategy{
						Variations: []*proto.RolloutStrategy_Variation{
							{
								Variation: f.Variations[0].Id,
								Weight:    30000,
							},
							{
								Variation: "",
								Weight:    70000,
							},
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			desc: "success",
			rule: &feature.Rule{
				Id: "rule-3",
				Strategy: &proto.Strategy{
					Type: proto.Strategy_ROLLOUT,
					RolloutStrategy: &proto.RolloutStrategy{
						Variations: []*proto.RolloutStrategy_Variation{
							{
								Variation: f.Variations[0].Id,
								Weight:    30000,
							},
							{
								Variation: f.Variations[1].Id,
								Weight:    70000,
							},
						},
					},
				},
				Clauses: []*proto.Clause{{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  proto.Clause_EQUALS,
					Values: []string{
						"user1",
					},
				}},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		err := f.AddRule(p.rule)
		assert.Equal(t, p.expectedErr, err != nil, "%s", p.desc)
		if !p.expectedErr {
			assert.Equal(t, p.rule, f.Rules[2], "%s", p.desc)
		}
	}
}

func TestChangeRuleStrategyToFixed(t *testing.T) {
	f := makeFeature("test-feature")
	r := f.Rules[0]
	rID := r.Id
	vID := f.Variations[1].Id
	expected := &proto.Strategy{
		Type:          proto.Strategy_FIXED,
		FixedStrategy: &proto.FixedStrategy{Variation: vID},
	}
	patterns := []*struct {
		ruleID   string
		strategy *proto.Strategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: expected,
			expected: errRuleNotFound,
		},
		{
			ruleID: rID,
			strategy: &proto.Strategy{
				Type:          proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{Variation: ""},
			},
			expected: errVariationNotFound,
		},
		{
			ruleID: rID,
			strategy: &proto.Strategy{
				Type:          proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{Variation: "variation-D"},
			},
			expected: errVariationNotFound,
		},
		{
			ruleID:   "",
			expected: errRuleNotFound,
		},
		{
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := f.ChangeRuleStrategy(p.ruleID, p.strategy)
		assert.Equal(t, p.expected, err)
	}
	if !reflect.DeepEqual(expected, r.Strategy) {
		t.Fatalf("Strategy is not equal. Expected: %s, actual: %s", expected, r.Strategy)
	}
}

func TestChangeRulesOrder(t *testing.T) {
	t.Helper()
	f := makeFeature("test-feature")
	patterns := []*struct {
		ruleIDs          []string
		expected         []string
		expectedUpdateAt int64
		expectedError    error
	}{
		{
			ruleIDs:          []string{f.Rules[0].Id, "not-found-id"},
			expected:         []string{f.Rules[0].Id, f.Rules[1].Id},
			expectedUpdateAt: f.UpdatedAt,
			expectedError:    errRuleNotFound,
		},
		{
			ruleIDs:          []string{f.Rules[0].Id},
			expected:         []string{f.Rules[0].Id, f.Rules[1].Id},
			expectedUpdateAt: f.UpdatedAt,
			expectedError:    errRulesOrderSizeNotEqual,
		},
		{
			ruleIDs:          []string{f.Rules[1].Id, f.Rules[1].Id},
			expected:         []string{f.Rules[0].Id, f.Rules[1].Id},
			expectedUpdateAt: f.UpdatedAt,
			expectedError:    errRulesOrderDuplicateIDs,
		},
		{
			ruleIDs:          []string{f.Rules[1].Id, f.Rules[0].Id},
			expected:         []string{f.Rules[1].Id, f.Rules[0].Id},
			expectedUpdateAt: time.Now().Unix(),
			expectedError:    nil,
		},
	}
	for _, p := range patterns {
		err := f.ChangeRulesOrder(p.ruleIDs)
		assert.Equal(t, p.expectedError, err)
		assert.Equal(t, p.expectedUpdateAt, f.UpdatedAt)
		for i := range f.Rules {
			if p.expected[i] != f.Rules[i].Id {
				t.Fatalf("Incorrect rules order. Expected: %s, actual: %s", p.expected[i], f.Rules[i].Id)
			}
		}
	}
}

func TestChangeRuleToRolloutStrategy(t *testing.T) {
	f := makeFeature("test-feature")
	r := f.Rules[0]
	rID := r.Id
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	expected := &proto.Strategy{
		Type: proto.Strategy_ROLLOUT,
		RolloutStrategy: &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
			{
				Variation: vID1,
				Weight:    30,
			},
			{
				Variation: vID2,
				Weight:    70,
			},
		}},
	}
	patterns := []*struct {
		ruleID   string
		strategy *proto.Strategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: expected,
			expected: errRuleNotFound,
		},
		{
			ruleID: rID,
			strategy: &proto.Strategy{
				Type: proto.Strategy_ROLLOUT,
				RolloutStrategy: &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
					{
						Variation: "",
						Weight:    30000,
					},
					{
						Variation: vID2,
						Weight:    70000,
					},
				}},
			},
			expected: errVariationNotFound,
		},
		{
			ruleID: rID,
			strategy: &proto.Strategy{
				Type: proto.Strategy_ROLLOUT,
				RolloutStrategy: &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
					{
						Variation: vID1,
						Weight:    30000,
					},
					{
						Variation: "",
						Weight:    70000,
					},
				}},
			},
			expected: errVariationNotFound,
		},
		{
			ruleID: rID,
			strategy: &proto.Strategy{
				Type: proto.Strategy_ROLLOUT,
				RolloutStrategy: &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
					{
						Variation: vID1,
						Weight:    30000,
					},
					{
						Variation: "variation-D",
						Weight:    70000,
					},
				}},
			},
			expected: errVariationNotFound,
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: errRuleNotFound,
		},
		{
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := f.ChangeRuleStrategy(p.ruleID, p.strategy)
		assert.Equal(t, p.expected, err)
	}
	if !reflect.DeepEqual(expected, r.Strategy) {
		t.Fatalf("Strategy is not equal. Expected: %v, actual: %v", expected, r.Strategy)
	}
}

func TestDeleteRule(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		rule        string
		expectedLen int
		expectedErr error
	}{
		{
			rule:        "",
			expectedLen: -1,
			expectedErr: errRuleNotFound,
		},
		{
			rule:        "rule-1",
			expectedLen: 1,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.DeleteRule(p.rule)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.expectedLen, len(f.Rules))
		}
	}
}

func TestDeleteClause(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		rule        string
		clause      string
		ruleIdx     int
		expectedLen int
		expectedErr error
	}{
		{
			rule:        "",
			clause:      "",
			ruleIdx:     -1,
			expectedLen: -1,
			expectedErr: errRuleNotFound,
		},
		{
			rule:        "rule-1",
			clause:      "clause-1",
			ruleIdx:     0,
			expectedLen: 0,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.DeleteClause(p.rule, p.clause)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.expectedLen, len(f.Rules[p.ruleIdx].Clauses))
		}
	}
}

func TestChangeClauseAttribute(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		rule        string
		clause      string
		attribute   string
		ruleIdx     int
		idx         int
		expectedErr error
	}{
		{
			rule:        "rule-1",
			clause:      "clause-1",
			attribute:   "newAttribute",
			ruleIdx:     0,
			idx:         0,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.ChangeClauseAttribute(p.rule, p.clause, p.attribute)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.attribute, f.Rules[p.ruleIdx].Clauses[p.idx].Attribute)
		}
	}
}

func TestChangeClauseOperator(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		rule        string
		clause      string
		operator    proto.Clause_Operator
		ruleIdx     int
		idx         int
		expectedErr error
	}{
		{
			rule:        "rule-1",
			clause:      "clause-1",
			operator:    proto.Clause_IN,
			ruleIdx:     0,
			idx:         0,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.ChangeClauseOperator(p.rule, p.clause, p.operator)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.operator, f.Rules[p.ruleIdx].Clauses[p.idx].Operator)
		}
	}
}

func TestAddClauseValueToFeature(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		rule        string
		clause      string
		value       string
		ruleIdx     int
		idx         int
		expectedLen int
		expectedErr error
	}{
		{
			rule:        "rule-1",
			clause:      "clause-1",
			value:       "newUser1",
			ruleIdx:     0,
			idx:         0,
			expectedLen: 3,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.AddClauseValue(p.rule, p.clause, p.value)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.expectedLen, len(f.Rules[p.ruleIdx].Clauses[p.idx].Values))
		}
	}
}

func TestRemoveClauseValueFromFeature(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		rule        string
		clause      string
		value       string
		ruleIdx     int
		idx         int
		expectedLen int
		expectedErr error
	}{
		{
			rule:        "rule-1",
			clause:      "clause-1",
			value:       "user1",
			ruleIdx:     0,
			idx:         0,
			expectedLen: 1,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.RemoveClauseValue(p.rule, p.clause, p.value)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.expectedLen, len(f.Rules[p.ruleIdx].Clauses[p.idx].Values))
		}
	}
}

func TestChangeVariationValue(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		id          string
		value       string
		idx         int
		expectedErr error
	}{
		{
			id:          "variation-A",
			value:       "newValue",
			idx:         0,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.ChangeVariationValue(p.id, p.value)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.value, f.Variations[p.idx].Value)
		}
	}
}

func TestChangeVariationName(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		id          string
		name        string
		idx         int
		expectedErr error
	}{
		{
			id:          "variation-A",
			name:        "newName",
			idx:         0,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.ChangeVariationName(p.id, p.name)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.name, f.Variations[p.idx].Name)
		}
	}
}

func TestChangeVariationDescription(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		id          string
		desc        string
		idx         int
		expectedErr error
	}{
		{
			id:          "variation-A",
			desc:        "newDesc",
			idx:         0,
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		err := f.ChangeVariationDescription(p.id, p.desc)
		if err != nil {
			assert.Equal(t, p.expectedErr, err)
		} else {
			assert.Equal(t, p.desc, f.Variations[p.idx].Description)
		}
	}
}

func TestRemoveVariationUsingFixedStrategy(t *testing.T) {
	f := makeFeature("test-feature")
	expected := "variation-C"
	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       "variation-A",
			expected: errVariationInUse,
		},
		{
			id:       "variation-B",
			expected: errVariationInUse,
		},
		{
			id:       expected,
			expected: nil,
		},
	}
	for i, p := range patterns {
		err := f.RemoveVariation(p.id)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}
	if _, err := f.findVariationIndex(expected); err == nil {
		t.Fatalf("Variation not deleted. Actual: %v", f.Variations)
	}
	actualSize := len(f.Variations)
	expectedSize := 2
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %d, actual: %d", expectedSize, actualSize)
	}
}

func TestRemoveVariationUsingRolloutStrategy(t *testing.T) {
	f := makeFeature("test-feature")
	f.ChangeDefaultStrategy(&proto.Strategy{
		Type: proto.Strategy_ROLLOUT,
		RolloutStrategy: &proto.RolloutStrategy{
			Variations: []*proto.RolloutStrategy_Variation{
				{
					Variation: "variation-A",
					Weight:    100000,
				},
				{
					Variation: "variation-B",
					Weight:    70000,
				},
				{
					Variation: "variation-C",
					Weight:    0,
				},
			},
		},
	})
	expected := "variation-C"
	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       "variation-A",
			expected: errVariationInUse,
		},
		{
			id:       "variation-B",
			expected: errVariationInUse,
		},
		{
			id:       expected,
			expected: nil,
		},
	}
	for i, p := range patterns {
		err := f.RemoveVariation(p.id)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}
	if _, err := f.findVariationIndex(expected); err == nil {
		t.Fatalf("Variation not deleted. Actual: %v", f.Variations)
	}
	actualSize := len(f.Variations)
	expectedSize := 2
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %d, actual: %d", expectedSize, actualSize)
	}
}

func TestRemoveVariationUsingOffVariation(t *testing.T) {
	f := makeFeature("test-feature")
	err := f.ChangeOffVariation("variation-C")
	assert.NoError(t, err)
	expected := "variation-D"
	f.AddVariation(
		expected,
		"value",
		"name",
		"description",
	)
	patterns := []*struct {
		des, id  string
		expected error
	}{
		{
			des:      "in use",
			id:       "variation-C",
			expected: errVariationInUse,
		},
		{
			des:      "success",
			id:       expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := f.RemoveVariation(p.id)
		assert.Equal(t, p.expected, err, p.des)
	}
	if _, err := f.findVariationIndex(expected); err == nil {
		t.Fatalf("Variation not deleted. Actual: %v", f.Variations)
	}
	actualSize := len(f.Variations)
	expectedSize := 3
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %d, actual: %d", expectedSize, actualSize)
	}
}

func TestChangeFixedStrategy(t *testing.T) {
	f := makeFeature("test-feature")
	r := f.Rules[0]
	rID := r.Id
	vID := f.Variations[1].Id
	patterns := []*struct {
		ruleID, variationID string
		expected            error
	}{
		{
			ruleID:      "",
			variationID: vID,
			expected:    errRuleNotFound,
		},
		{
			ruleID:      rID,
			variationID: "",
			expected:    errVariationNotFound,
		},
		{
			ruleID:      "",
			variationID: "",
			expected:    errRuleNotFound,
		},
		{
			ruleID:      rID,
			variationID: vID,
			expected:    nil,
		},
	}
	for _, p := range patterns {
		err := f.ChangeFixedStrategy(p.ruleID, &proto.FixedStrategy{Variation: p.variationID})
		assert.Equal(t, p.expected, err)
	}
	if r.Strategy.FixedStrategy.Variation != vID {
		t.Fatalf("Wrong variation id has been saved. Expected: %s, actual: %s", vID, r.Strategy.FixedStrategy.Variation)
	}
}

func TestChangeRolloutStrategy(t *testing.T) {
	f := makeFeature("test-feature")
	r := f.Rules[0]
	rID := r.Id
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	expected := &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
		{
			Variation: vID1,
			Weight:    30,
		},
		{
			Variation: vID2,
			Weight:    70,
		},
	}}
	patterns := []*struct {
		ruleID   string
		strategy *proto.RolloutStrategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: &proto.RolloutStrategy{},
			expected: errRuleNotFound,
		},
		{
			ruleID: rID,
			strategy: &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
				{
					Variation: "",
					Weight:    30,
				},
				{
					Variation: vID2,
					Weight:    70,
				},
			}},
			expected: errVariationNotFound,
		},
		{
			ruleID: rID,
			strategy: &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    30,
				},
				{
					Variation: "",
					Weight:    70,
				},
			}},
			expected: errVariationNotFound,
		},
		{
			ruleID:   "",
			strategy: nil,
			expected: errRuleNotFound,
		},
		{
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	for _, p := range patterns {
		err := f.ChangeRolloutStrategy(p.ruleID, p.strategy)
		assert.Equal(t, p.expected, err)
	}
	if !reflect.DeepEqual(expected, r.Strategy.RolloutStrategy) {
		t.Fatalf("Different rollout strategies. Expected: %v, actual: %v", expected, r.Strategy.RolloutStrategy)
	}
}

func TestIsStale(t *testing.T) {
	t.Parallel()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-01 0:00:00 +0000 UTC")
	require.NoError(t, err)
	t2, err := time.Parse(layout, "2014-03-31 23:59:59 +0000 UTC")
	require.NoError(t, err)
	t3, err := time.Parse(layout, "2014-04-01 0:00:00 +0000 UTC")
	require.NoError(t, err)
	patterns := []struct {
		desc     string
		feature  *Feature
		input    time.Time
		expected bool
	}{
		{
			desc: "false",
			feature: &Feature{Feature: &proto.Feature{
				LastUsedInfo: &proto.FeatureLastUsedInfo{
					LastUsedAt: t1.Unix(),
				},
			}},
			input:    t2,
			expected: false,
		},
		{
			desc: "true",
			feature: &Feature{Feature: &proto.Feature{
				LastUsedInfo: &proto.FeatureLastUsedInfo{
					LastUsedAt: t1.Unix(),
				},
			}},
			input:    t3,
			expected: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, p.feature.IsStale(p.input))
		})
	}
}

func TestValidateVariations(t *testing.T) {
	t.Parallel()
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	patterns := []struct {
		desc          string
		variationType feature.Feature_VariationType
		variations    []*feature.Variation
		expectedErr   bool
	}{
		{
			desc:          "fail: only one variation",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{Id: id1.String(), Name: "v1", Value: "true"},
			},
			expectedErr: true,
		},
		{
			desc:          "fail: empty id",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{Id: "", Name: "v1", Value: "true"},
				{Id: id2.String(), Name: "v2", Value: "false"},
			},
			expectedErr: true,
		},
		{
			desc:          "fail: empty name",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{Id: id1.String(), Name: "", Value: "true"},
				{Id: id2.String(), Name: "v2", Value: "false"},
			},
			expectedErr: true,
		},
		{
			desc:          "fail: invalid value",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{Id: id1.String(), Name: "v1", Value: "foo"},
				{Id: id2.String(), Name: "v2", Value: "false"},
			},
			expectedErr: true,
		},
		{
			desc:          "fail: id is duplicated",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{Id: id1.String(), Name: "v1", Value: "true"},
				{Id: id1.String(), Name: "v2", Value: "false"},
			},
			expectedErr: true,
		},
		{
			desc:          "fail: invalid id",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{Id: "v1", Name: "", Value: "true"},
				{Id: id2.String(), Name: "v2", Value: "false"},
			},
			expectedErr: true,
		},
		{
			desc:          "success",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{Id: id1.String(), Name: "v1", Value: "true"},
				{Id: id2.String(), Name: "v2", Value: "false"},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateVariations(p.variationType, p.variations)
			assert.Equal(t, p.expectedErr, err != nil)
		})
	}
}

func TestValidateVariationValue(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		variationType feature.Feature_VariationType
		value         string
		expected      error
	}{
		{
			desc:          "invalid bool",
			variationType: feature.Feature_BOOLEAN,
			value:         "hoge",
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "empty string",
			variationType: feature.Feature_JSON,
			value:         "",
			expected:      errors.New("feature: variation value cannot be empty"),
		},
		{
			desc:          "invalid number",
			variationType: feature.Feature_NUMBER,
			value:         `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`,
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "invalid json",
			variationType: feature.Feature_JSON,
			value:         "true",
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "valid bool",
			variationType: feature.Feature_BOOLEAN,
			value:         "true",
			expected:      nil,
		},
		{
			desc:          "valid number float",
			variationType: feature.Feature_NUMBER,
			value:         "1.23",
			expected:      nil,
		},
		{
			desc:          "valid number int",
			variationType: feature.Feature_NUMBER,
			value:         "123",
			expected:      nil,
		},
		{
			desc:          "valid json",
			variationType: feature.Feature_JSON,
			value:         `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`,
			expected:      nil,
		},
		{
			desc:          "valid json array",
			variationType: feature.Feature_JSON,
			value:         `[{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}]`,
			expected:      nil,
		},
		{
			desc:          "valid string",
			variationType: feature.Feature_STRING,
			value:         `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`,
			expected:      nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, validateVariationValue(p.variationType, p.value))
		})
	}
}

func TestNewClonedFeature(t *testing.T) {
	t.Parallel()
	pattenrs := []struct {
		maintainer        string
		offVariationIndex int
		expectedEnabled   bool
		expectedVersion   int32
		defaultStrategy   *feature.Strategy
		rules             []*feature.Rule
	}{
		{
			maintainer:        "sample@example.com",
			offVariationIndex: 2,
			expectedEnabled:   false,
			expectedVersion:   int32(1),
			defaultStrategy: &proto.Strategy{
				Type: proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{
					Variation: "variation-B",
				},
			},
			rules: []*proto.Rule{
				{
					Id: "rule-1",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*proto.Clause{
						{
							Id:        "clause-1",
							Attribute: "name",
							Operator:  proto.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
				{
					Id: "rule-2",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{
							Variation: "variation-B",
						},
					},
					Clauses: []*proto.Clause{
						{
							Id:        "clause-2",
							Attribute: "name",
							Operator:  proto.Clause_EQUALS,
							Values: []string{
								"user3",
								"user4",
							},
						},
					},
				},
			},
		},
		{
			maintainer:        "sample@example.com",
			offVariationIndex: 2,
			expectedEnabled:   false,
			expectedVersion:   int32(1),
			defaultStrategy: &proto.Strategy{
				Type: proto.Strategy_ROLLOUT,
				RolloutStrategy: &proto.RolloutStrategy{
					Variations: []*proto.RolloutStrategy_Variation{
						{
							Variation: "variation-A",
							Weight:    100000,
						},
						{
							Variation: "variation-B",
							Weight:    70000,
						},
						{
							Variation: "variation-C",
							Weight:    0,
						},
					},
				},
			},
			rules: []*proto.Rule{
				{
					Id: "rule-1",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_ROLLOUT,
						RolloutStrategy: &proto.RolloutStrategy{
							Variations: []*proto.RolloutStrategy_Variation{
								{
									Variation: "variation-A",
									Weight:    100000,
								},
								{
									Variation: "variation-B",
									Weight:    70000,
								},
								{
									Variation: "variation-C",
									Weight:    0,
								},
							},
						},
					},
				},
				{
					Id: "rule-2",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_ROLLOUT,
						RolloutStrategy: &proto.RolloutStrategy{
							Variations: []*proto.RolloutStrategy_Variation{
								{
									Variation: "variation-A",
									Weight:    100,
								},
								{
									Variation: "variation-B",
									Weight:    500,
								},
								{
									Variation: "variation-C",
									Weight:    300,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, p := range pattenrs {
		f := makeFeature("test-feature")
		f.Maintainer = "bucketeer@example.com"
		f.OffVariation = f.Variations[p.offVariationIndex].Id
		f.DefaultStrategy = p.defaultStrategy
		f.Rules = p.rules
		actual, err := f.Clone(p.maintainer)
		assert.NoError(t, err)
		assert.Equal(t, p.maintainer, actual.Maintainer)
		assert.Equal(t, p.expectedEnabled, actual.Enabled)
		assert.Equal(t, p.expectedVersion, actual.Version)
		assert.Equal(t, actual.OffVariation, actual.Variations[p.offVariationIndex].Id)
		for i := range actual.Variations {
			assert.Equal(t, actual.Variations[i].Id, actual.Targets[i].Variation)
		}
		if actual.DefaultStrategy.Type == feature.Strategy_FIXED {
			assert.Equal(t, actual.Variations[1].Id, actual.DefaultStrategy.FixedStrategy.Variation)
		} else {
			for i := range actual.Variations {
				assert.Equal(t, actual.Variations[i].Id, actual.DefaultStrategy.RolloutStrategy.Variations[i].Variation)
			}
		}
		assert.NotNil(t, actual.Prerequisites)
		assert.Equal(t, len(actual.Prerequisites), 0)
		for i := range actual.Rules {
			if actual.Rules[i].Strategy.Type == feature.Strategy_FIXED {
				assert.Equal(t, actual.Rules[i].Strategy.FixedStrategy.Variation, actual.Variations[i].Id)
			} else {
				for idx := range actual.Variations {
					assert.Equal(t, actual.Rules[i].Strategy.RolloutStrategy.Variations[idx].Variation, actual.Variations[idx].Id)
				}
			}
		}
	}
}

func TestResetSamplingSeed(t *testing.T) {
	f := makeFeature("test-feature")
	assert.Empty(t, f.SamplingSeed)
	err := f.ResetSamplingSeed()
	assert.NoError(t, err)
	assert.NotEmpty(t, f.SamplingSeed)
}

func TestFeatureIDsDependsOn(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		feature  *feature.Feature
		expected []string
	}{
		{
			feature:  &feature.Feature{},
			expected: []string{},
		},
		{
			feature: &feature.Feature{
				Prerequisites: []*feature.Prerequisite{
					{FeatureId: "feature-1"},
				},
			},
			expected: []string{"feature-1"},
		},
		{
			feature: &feature.Feature{
				Prerequisites: []*feature.Prerequisite{
					{FeatureId: "feature-1"},
					{FeatureId: "feature-2"},
				},
			},
			expected: []string{"feature-1", "feature-2"},
		},
		{
			feature: &feature.Feature{
				Rules: []*feature.Rule{
					{
						Clauses: []*feature.Clause{
							{Attribute: "feature-1", Operator: feature.Clause_FEATURE_FLAG},
						},
					},
				},
			},
			expected: []string{"feature-1"},
		},
		{
			feature: &feature.Feature{
				Rules: []*feature.Rule{
					{
						Clauses: []*feature.Clause{
							{Attribute: "feature-1", Operator: feature.Clause_FEATURE_FLAG},
							{Attribute: "feature-2", Operator: feature.Clause_FEATURE_FLAG},
						},
					},
				},
			},
			expected: []string{"feature-1", "feature-2"},
		},
		{
			feature: &feature.Feature{
				Prerequisites: []*feature.Prerequisite{
					{FeatureId: "feature-1"},
				},
				Rules: []*feature.Rule{
					{
						Clauses: []*feature.Clause{
							{Attribute: "feature-2", Operator: feature.Clause_FEATURE_FLAG},
						},
					},
				},
			},
			expected: []string{"feature-1", "feature-2"},
		},
	}
	for _, p := range patterns {
		f := &Feature{Feature: p.feature}
		assert.Equal(t, p.expected, f.FeatureIDsDependsOn())
	}
}

func TestValidateClauses(t *testing.T) {
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	id3, _ := uuid.NewUUID()
	patterns := []struct {
		desc     string
		clauses  []*proto.Clause
		expected error
	}{
		{
			desc: "success: 2 clauses",
			clauses: []*proto.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  proto.Clause_EQUALS,
					Values: []string{
						"user1",
						"user2",
					},
				},
				{
					Id:       id2.String(),
					Operator: proto.Clause_SEGMENT,
					Values: []string{
						"value",
					},
				},
				{
					Id:        id3.String(),
					Operator:  proto.Clause_FEATURE_FLAG,
					Attribute: "feature-1",
					Values: []string{
						"true",
					},
				},
			},
			expected: nil,
		},
		{
			desc:     "err: zero clause",
			clauses:  []*proto.Clause{},
			expected: fmt.Errorf("feature: rule must have at least one clause"),
		},
		{
			desc: "err: id is empty",
			clauses: []*proto.Clause{
				{
					Id:        "",
					Attribute: "name",
					Operator:  proto.Clause_EQUALS,
					Values:    []string{"user1"},
				},
			},
			expected: errors.New("feature: clause id cannot be empty"),
		},
		{
			desc: "err: compare missing attribute",
			clauses: []*proto.Clause{
				{
					Id:       id1.String(),
					Operator: proto.Clause_EQUALS,
					Values: []string{
						"user1",
						"user2",
					},
				},
			},
			expected: errClauseAttributeEmpty,
		},
		{
			desc: "err: compare missing values",
			clauses: []*proto.Clause{
				{
					Id:        id1.String(),
					Operator:  proto.Clause_EQUALS,
					Attribute: "name",
				},
			},
			expected: errClauseValuesEmpty,
		},
		{
			desc: "err: segment attribute not empty",
			clauses: []*proto.Clause{
				{
					Id:        id1.String(),
					Operator:  proto.Clause_SEGMENT,
					Attribute: "name",
					Values: []string{
						"user1",
					},
				},
			},
			expected: errClauseAttributeNotEmpty,
		},
		{
			desc: "err: segment value empty",
			clauses: []*proto.Clause{
				{
					Id:       id1.String(),
					Operator: proto.Clause_SEGMENT,
				},
			},
			expected: errClauseValuesEmpty,
		},
		{
			desc: "err: feature flag attribute empty",
			clauses: []*proto.Clause{
				{
					Id:       id1.String(),
					Operator: proto.Clause_FEATURE_FLAG,
				},
			},
			expected: errClauseAttributeEmpty,
		},
		{
			desc: "err: feature flag values empty",
			clauses: []*proto.Clause{
				{
					Id:        id1.String(),
					Operator:  proto.Clause_FEATURE_FLAG,
					Attribute: "feature-1",
				},
			},
			expected: errClauseValuesEmpty,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, validateClauses(p.clauses))
		})
	}
}

func TestUpdate(t *testing.T) {
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	id3, _ := uuid.NewUUID()
	genF := func() *Feature {
		return &Feature{
			Feature: &proto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1", "t2"},
				VariationType: feature.Feature_BOOLEAN,
				Variations: []*proto.Variation{
					{Id: id1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: id2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*proto.Prerequisite{},
				Targets: []*proto.Target{
					{Variation: id1.String()},
					{Variation: id2.String()},
				},
				Rules: []*proto.Rule{},
				DefaultStrategy: &proto.Strategy{
					Type:          proto.Strategy_FIXED,
					FixedStrategy: &proto.FixedStrategy{Variation: id1.String()},
				},
				OffVariation: id1.String(),
			},
		}
	}
	patterns := []struct {
		desc            string
		inputFunc       func() *Feature
		name            *wrapperspb.StringValue
		description     *wrapperspb.StringValue
		enabled         *wrapperspb.BoolValue
		tags            *proto.StringListValue
		archived        *wrapperspb.BoolValue
		variations      *proto.VariationListValue
		prerequisites   *proto.PrerequisiteListValue
		targets         *proto.TargetListValue
		rules           *proto.RuleListValue
		defaultStrategy *proto.Strategy
		offVariation    *wrapperspb.StringValue
		expectedFunc    func() *Feature
		expectedErr     error
	}{
		{
			desc: "fail: name is empty",
			inputFunc: func() *Feature {
				return genF()
			},
			name:         &wrapperspb.StringValue{Value: ""},
			description:  &wrapperspb.StringValue{Value: "description"},
			expectedFunc: func() *Feature { return nil },
			expectedErr:  errNameEmpty,
		},
		{
			desc: "fail: already enabled",
			inputFunc: func() *Feature {
				f := genF()
				f.Enabled = true
				return f
			},
			enabled:      &wrapperspb.BoolValue{Value: true},
			expectedFunc: func() *Feature { return nil },
			expectedErr:  ErrAlreadyEnabled,
		},
		{
			desc: "fail: already disabled",
			inputFunc: func() *Feature {
				f := genF()
				f.Enabled = false
				return f
			},
			enabled:      &wrapperspb.BoolValue{Value: false},
			expectedFunc: func() *Feature { return nil },
			expectedErr:  ErrAlreadyDisabled,
		},
		{
			desc: "success: version incremented",
			inputFunc: func() *Feature {
				f := genF()
				return f
			},
			name:        &wrapperspb.StringValue{Value: "n2"},
			description: &wrapperspb.StringValue{Value: "d2"},
			enabled:     &wrapperspb.BoolValue{Value: true},
			archived:    &wrapperspb.BoolValue{Value: true},
			tags:        &proto.StringListValue{Values: []string{"t3"}},
			variations: &proto.VariationListValue{Values: []*feature.Variation{
				{Id: id1.String(), Value: "true", Name: "n3"},
				{Id: id2.String(), Value: "false", Name: "n4"},
			}},
			prerequisites: &proto.PrerequisiteListValue{Values: []*feature.Prerequisite{
				{FeatureId: "f1", VariationId: "v1"},
			}},
			targets: &proto.TargetListValue{Values: []*feature.Target{
				{Variation: id1.String(), Users: []string{"uid1"}},
				{Variation: id2.String(), Users: []string{"uid2"}},
			}},
			rules: &proto.RuleListValue{Values: []*feature.Rule{
				{
					Id: id3.String(),
					Strategy: &feature.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{Variation: id1.String()},
					},
					Clauses: []*feature.Clause{
						{
							Id:        id1.String(),
							Attribute: "name",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"user1", "user2"},
						},
					},
				},
			}},
			defaultStrategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{{Variation: id1.String(), Weight: 100}},
				},
			},
			offVariation: &wrapperspb.StringValue{Value: id1.String()},
			expectedFunc: func() *Feature {
				f := genF()
				f.Name = "n2"
				f.Description = "d2"
				f.Archived = true
				f.Enabled = true
				f.Tags = []string{"t3"}
				f.UpdatedAt = time.Now().Unix()
				f.Version = 1
				f.Variations = []*feature.Variation{
					{Id: id1.String(), Value: "true", Name: "n3"},
					{Id: id2.String(), Value: "false", Name: "n4"},
				}
				f.Prerequisites = []*feature.Prerequisite{{FeatureId: "f1", VariationId: "v1"}}
				f.Targets =
					[]*feature.Target{
						{Variation: id1.String(), Users: []string{"uid1"}},
						{Variation: id2.String(), Users: []string{"uid2"}},
					}
				f.Rules = []*feature.Rule{
					{
						Id: id3.String(),
						Strategy: &feature.Strategy{
							Type:          feature.Strategy_FIXED,
							FixedStrategy: &feature.FixedStrategy{Variation: id1.String()},
						},
						Clauses: []*feature.Clause{
							{
								Id:        id1.String(),
								Attribute: "name",
								Operator:  feature.Clause_EQUALS,
								Values:    []string{"user1", "user2"},
							},
						},
					},
				}
				f.DefaultStrategy = &feature.Strategy{
					Type: feature.Strategy_ROLLOUT,
					RolloutStrategy: &feature.RolloutStrategy{
						Variations: []*feature.RolloutStrategy_Variation{{Variation: id1.String(), Weight: 100}},
					},
				}
				f.OffVariation = id1.String()
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "success: version not incremented",
			inputFunc: func() *Feature {
				return genF()
			},
			description: &wrapperspb.StringValue{Value: "d2"},
			expectedFunc: func() *Feature {
				f := genF()
				f.Description = "d2"
				f.Version = 0
				return f
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.inputFunc().Update(
				p.name, p.description,
				p.tags, p.enabled, p.archived,
				p.variations, p.prerequisites,
				p.targets, p.rules, p.defaultStrategy,
				p.offVariation,
			)
			if p.expectedErr == nil && actual != nil {
				assert.Equal(t, p.expectedFunc().Version, actual.Version, p.desc)
				assert.Equal(t, p.expectedFunc().Name, actual.Name, p.desc)
				assert.Equal(t, p.expectedFunc().Description, actual.Description, p.desc)
				assert.Equal(t, p.expectedFunc().Enabled, actual.Enabled, p.desc)
				assert.Equal(t, p.expectedFunc().Tags, actual.Tags, p.desc)
				assert.Equal(t, p.expectedFunc().Archived, actual.Archived, p.desc)
				assert.Equal(t, p.expectedFunc().Variations, actual.Variations, p.desc)
				assert.Equal(t, p.expectedFunc().Prerequisites, actual.Prerequisites, p.desc)
				assert.Equal(t, p.expectedFunc().Targets, actual.Targets, p.desc)
				assert.Equal(t, p.expectedFunc().Rules, actual.Rules, p.desc)
				assert.Equal(t, p.expectedFunc().DefaultStrategy, actual.DefaultStrategy, p.desc)
				assert.Equal(t, p.expectedFunc().OffVariation, actual.OffVariation, p.desc)
				assert.LessOrEqual(t, p.expectedFunc().UpdatedAt, actual.UpdatedAt)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestTopologicalSort(t *testing.T) {
	t.Parallel()
	makeFeature := func(id string) *proto.Feature {
		return &proto.Feature{
			Id:            id,
			Name:          "test feature",
			Version:       1,
			Enabled:       true,
			CreatedAt:     time.Now().Unix(),
			VariationType: feature.Feature_STRING,
			Variations: []*proto.Variation{
				{
					Id:          "variation-A",
					Value:       "A",
					Name:        "Variation A",
					Description: "Thing does A",
				},
				{
					Id:          "variation-B",
					Value:       "B",
					Name:        "Variation B",
					Description: "Thing does B",
				},
				{
					Id:          "variation-C",
					Value:       "C",
					Name:        "Variation C",
					Description: "Thing does C",
				},
			},
			Targets: []*proto.Target{
				{
					Variation: "variation-A",
					Users: []string{
						"user1",
					},
				},
				{
					Variation: "variation-B",
					Users: []string{
						"user2",
					},
				},
				{
					Variation: "variation-C",
					Users: []string{
						"user3",
					},
				},
			},
			Rules: []*proto.Rule{
				{
					Id: "rule-1",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*proto.Clause{
						{
							Id:        "clause-1",
							Attribute: "name",
							Operator:  proto.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
				{
					Id: "rule-2",
					Strategy: &proto.Strategy{
						Type: proto.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{
							Variation: "variation-B",
						},
					},
					Clauses: []*proto.Clause{
						{
							Id:        "clause-2",
							Attribute: "name",
							Operator:  proto.Clause_EQUALS,
							Values: []string{
								"user3",
								"user4",
							},
						},
					},
				},
			},
			DefaultStrategy: &proto.Strategy{
				Type: proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{
					Variation: "variation-B",
				},
			},
		}
	}
	f0 := makeFeature("fID0")
	f1 := makeFeature("fID1")
	f2 := makeFeature("fID2")
	f3 := makeFeature("fID3")
	f4 := makeFeature("fID4")
	f5 := makeFeature("fID5")
	patterns := []struct {
		f0Prerequisite []*proto.Prerequisite
		f1Prerequisite []*proto.Prerequisite
		f2Prerequisite []*proto.Prerequisite
		f3Prerequisite []*proto.Prerequisite
		f4Prerequisite []*proto.Prerequisite
		f5Prerequisite []*proto.Prerequisite
		expected       []*proto.Feature
		expectedError  error
	}{
		{
			f0Prerequisite: []*proto.Prerequisite{},
			f1Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
			},
			f2Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
			},
			f3Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
				{
					FeatureId: f2.Id,
				},
			},
			f4Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			f5Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f4.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			expected: []*proto.Feature{
				f0, f1, f2, f3, f4, f5,
			},
			expectedError: nil,
		},
		{
			f0Prerequisite: []*proto.Prerequisite{},
			f1Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
			},
			f2Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
			},
			f3Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
				{
					FeatureId: f2.Id,
				},
			},
			f4Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			f5Prerequisite: []*proto.Prerequisite{},
			expected: []*proto.Feature{
				f0, f1, f2, f5, f3, f4,
			},
			expectedError: nil,
		},
		{
			f0Prerequisite: []*proto.Prerequisite{},
			f1Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
			},
			f2Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f3.Id,
				},
			},
			f3Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f2.Id,
				},
			},
			f4Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			f5Prerequisite: []*proto.Prerequisite{
				{
					FeatureId: f4.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			expected:      nil,
			expectedError: ErrCycleExists,
		},
		{
			f0Prerequisite: []*proto.Prerequisite{},
			f1Prerequisite: []*proto.Prerequisite{},
			f2Prerequisite: []*proto.Prerequisite{},
			f3Prerequisite: []*proto.Prerequisite{},
			f4Prerequisite: []*proto.Prerequisite{},
			f5Prerequisite: []*proto.Prerequisite{},
			expected: []*proto.Feature{
				f2, f0, f5, f3, f1, f4,
			},
			expectedError: nil,
		},
	}
	for _, p := range patterns {
		f0.Prerequisites = p.f0Prerequisite
		f1.Prerequisites = p.f1Prerequisite
		f2.Prerequisites = p.f2Prerequisite
		f3.Prerequisites = p.f3Prerequisite
		f4.Prerequisites = p.f4Prerequisite
		f5.Prerequisites = p.f5Prerequisite
		fs := []*proto.Feature{
			f2, f0, f5, f3, f1, f4,
		}
		actual, err := TopologicalSort(fs)
		assert.Equal(t, p.expectedError, err)
		assert.Equal(t, p.expected, actual)
	}
}

func TestHasFeaturesDependsOnTargets(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		targets  []*feature.Feature
		all      []*feature.Feature
		expected bool
	}{
		{
			targets:  []*feature.Feature{},
			all:      []*feature.Feature{},
			expected: false,
		},
		{
			targets: []*feature.Feature{
				{Id: "1"},
			},
			all: []*feature.Feature{
				{Id: "1"},
			},
			expected: false,
		},
		{
			targets: []*feature.Feature{
				{Id: "1"},
			},
			all: []*feature.Feature{
				{Id: "1"},
				{Id: "2", Prerequisites: []*feature.Prerequisite{{FeatureId: "1"}}},
			},
			expected: true,
		},
	}
	for _, p := range patterns {
		assert.Equal(t, p.expected, HasFeaturesDependsOnTargets(p.targets, p.all))
	}
}

func TestValidateOffVariation(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		id          string
		variations  []*feature.Variation
		expectedErr bool
	}{
		{
			desc: "fails: id is empty",
			id:   "",
			variations: []*feature.Variation{
				{Id: "v1"},
				{Id: "v2"},
			},
			expectedErr: true,
		},
		{
			desc: "fails: id not found",
			id:   "v1",
			variations: []*feature.Variation{
				{Id: "v2"},
				{Id: "v3"},
			},
			expectedErr: true,
		},
		{
			desc: "success",
			id:   "v1",
			variations: []*feature.Variation{
				{Id: "v1"},
				{Id: "v2"},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expectedErr, validateOffVariation(p.id, p.variations) != nil)
		})
	}
}

func TestValidateTargets(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		targets     []*feature.Target
		variations  []*feature.Variation
		expectedErr bool
	}{
		{
			desc: "fail: variation not found",
			targets: []*feature.Target{
				{Variation: "v1"},
				{Variation: "v3"},
			},
			variations: []*feature.Variation{
				{Id: "v1"},
				{Id: "v2"},
			},
			expectedErr: true,
		},
		{
			desc: "success",
			targets: []*feature.Target{
				{Variation: "v1"},
				{Variation: "v2"},
			},
			variations: []*feature.Variation{
				{Id: "v1"},
				{Id: "v2"},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expectedErr, validateTargets(p.targets, p.variations) != nil)
		})
	}
}

func TestValidatePrerequisites(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		prerequisites []*feature.Prerequisite
		expectedErr   bool
	}{
		{
			desc: "fail: feature id is empty",
			prerequisites: []*feature.Prerequisite{
				{FeatureId: "", VariationId: "v1"},
			},
			expectedErr: true,
		},
		{
			desc: "fail: variation id is empty",
			prerequisites: []*feature.Prerequisite{
				{FeatureId: "f1", VariationId: ""},
			},
			expectedErr: true,
		},
		{
			desc: "success",
			prerequisites: []*feature.Prerequisite{
				{FeatureId: "f1", VariationId: "v1"},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expectedErr, validatePrerequisites(p.prerequisites) != nil)
		})
	}
}

func TestValidateRules(t *testing.T) {
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	t.Parallel()
	patterns := []struct {
		desc        string
		rules       []*feature.Rule
		variations  []*feature.Variation
		expectedErr bool
	}{
		{
			desc:        "fail: rule id is empty",
			rules:       []*feature.Rule{{Id: ""}},
			expectedErr: true,
		},
		{
			desc:        "fail: rule id is not uuid",
			rules:       []*feature.Rule{{Id: "r1"}},
			expectedErr: true,
		},
		{
			desc: "fail: rule id is duplicated",
			rules: []*feature.Rule{
				{Id: id1.String()},
				{Id: id1.String()},
			},
			expectedErr: true,
		},
		{
			desc: "fail: rule id is duplicated",
			rules: []*feature.Rule{
				{
					Id: id1.String(),
					Strategy: &feature.Strategy{
						Type: feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{
							Variation: "v1",
						},
					},
					Clauses: []*feature.Clause{
						{
							Id:        id2.String(),
							Attribute: "name",
							Operator:  feature.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
				{
					Id: id1.String(),
					Strategy: &feature.Strategy{
						Type: feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{
							Variation: "v1",
						},
					},
					Clauses: []*feature.Clause{
						{
							Id:        id2.String(),
							Attribute: "name",
							Operator:  feature.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			desc: "fail: invalid strategy",
			rules: []*feature.Rule{
				{
					Id: id1.String(),
					Strategy: &feature.Strategy{
						Type: feature.Strategy_FIXED,
					},
				},
			},
			expectedErr: true,
		},
		{
			desc: "fail: invalid clause",
			rules: []*feature.Rule{
				{
					Id: id1.String(),
					Strategy: &feature.Strategy{
						Type: feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{
							Variation: "v1",
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			desc: "success",
			rules: []*feature.Rule{
				{
					Id: id1.String(),
					Strategy: &feature.Strategy{
						Type: feature.Strategy_FIXED,
						FixedStrategy: &feature.FixedStrategy{
							Variation: "v1",
						},
					},
					Clauses: []*feature.Clause{
						{
							Id:        id2.String(),
							Attribute: "name",
							Operator:  feature.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
			},
			variations: []*feature.Variation{
				{Id: "v1"},
				{Id: "v2"},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expectedErr, validateRules(p.rules, p.variations) != nil)
		})
	}
}
