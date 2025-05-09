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
	"github.com/bucketeer-io/bucketeer/proto/common"
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
	t.Parallel()
	patterns := []*struct {
		desc                     string
		id                       string
		name                     string
		description              string
		variationType            feature.Feature_VariationType
		variations               []*feature.Variation
		tags                     []string
		defaultOnVariationIndex  int
		defaultOffVariationIndex int
		maintainer               string
		expected                 error
	}{
		{
			desc:          "err: variations must have at least two variations",
			id:            "test-feature",
			name:          "test feature",
			description:   "test feature description",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{
					Value:       "true",
					Name:        "Variation A",
					Description: "Thing does A",
				},
			},
			tags:                     []string{},
			defaultOnVariationIndex:  0,
			defaultOffVariationIndex: 0,
			maintainer:               "test@example.com",
			expected:                 errVariationsMustHaveAtLeastTwoVariations,
		},
		{
			desc:          "err: invalid default on variation index",
			id:            "test-feature",
			name:          "test feature",
			description:   "test feature description",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{
					Value:       "true",
					Name:        "Variation A",
					Description: "Thing does A",
				},
				{
					Value:       "false",
					Name:        "Variation B",
					Description: "Thing does B",
				},
			},
			tags:                     []string{},
			defaultOnVariationIndex:  2, // Out of range
			defaultOffVariationIndex: 0,
			maintainer:               "test@example.com",
			expected:                 errInvalidDefaultOnVariationIndex,
		},
		{
			desc:          "err: invalid default off variation index",
			id:            "test-feature",
			name:          "test feature",
			description:   "test feature description",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{
					Value:       "true",
					Name:        "Variation A",
					Description: "Thing does A",
				},
				{
					Value:       "false",
					Name:        "Variation B",
					Description: "Thing does B",
				},
			},
			tags:                     []string{},
			defaultOnVariationIndex:  0,
			defaultOffVariationIndex: 2, // Out of range
			maintainer:               "test@example.com",
			expected:                 errInvalidDefaultOffVariationIndex,
		},
		{
			desc:          "success",
			id:            "test-feature",
			name:          "test feature",
			description:   "test feature description",
			variationType: feature.Feature_BOOLEAN,
			variations: []*feature.Variation{
				{
					Value:       "true",
					Name:        "Variation A",
					Description: "Thing does A",
				},
				{
					Value:       "false",
					Name:        "Variation B",
					Description: "Thing does B",
				},
			},
			tags:                     []string{"tag1", "tag2"},
			defaultOnVariationIndex:  0,
			defaultOffVariationIndex: 1,
			maintainer:               "test@example.com",
			expected:                 nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f, err := NewFeature(
				p.id,
				p.name,
				p.description,
				p.variationType,
				p.variations,
				p.tags,
				p.defaultOnVariationIndex,
				p.defaultOffVariationIndex,
				p.maintainer,
			)
			assert.Equal(t, p.expected, err)
			if err == nil {
				assert.Equal(t, p.id, f.Id)
				assert.Equal(t, p.name, f.Name)
				assert.Equal(t, p.description, f.Description)
				assert.Equal(t, p.variationType, f.VariationType)
				assert.Equal(t, p.tags, f.Tags)
				assert.Equal(t, p.maintainer, f.Maintainer)
				assert.Equal(t, int32(1), f.Version)
				assert.NotEmpty(t, f.CreatedAt)
				assert.NotEmpty(t, f.UpdatedAt)
				assert.False(t, f.Enabled)
				assert.False(t, f.Deleted)
				assert.False(t, f.Archived)
				assert.Empty(t, f.Prerequisites)
				assert.Empty(t, f.Rules)
				assert.NotEmpty(t, f.Variations)
				assert.NotEmpty(t, f.Targets)
				assert.NotEmpty(t, f.DefaultStrategy)
				assert.NotEmpty(t, f.OffVariation)
			}
		})
	}
}

func TestAddVariation(t *testing.T) {
	t.Parallel()
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	patterns := []struct {
		desc          string
		variationType feature.Feature_VariationType
		id            string
		name          string
		value         string
		description   string
		expectedErr   error
		variations    []*feature.Variation
	}{
		{
			desc:          "fail: empty name",
			variationType: feature.Feature_BOOLEAN,
			id:            id1.String(),
			name:          "",
			value:         "true",
			description:   "first variation",
			expectedErr:   errVariationNameRequired,
			variations: []*feature.Variation{
				{Id: id1.String(), Name: "v1", Value: "true", Description: "first variation"},
			},
		},
		{
			desc:          "fail: empty value",
			variationType: feature.Feature_BOOLEAN,
			id:            id1.String(),
			name:          "v1",
			value:         "",
			description:   "first variation",
			expectedErr:   errVariationValueRequired,
		},
		{
			desc:          "fail: duplicate value",
			variationType: feature.Feature_BOOLEAN,
			id:            id2.String(),
			name:          "v2",
			value:         "true", // same value as first variation
			description:   "second variation",
			expectedErr:   errVariationValueUnique,
			variations: []*feature.Variation{
				{Id: id1.String(), Name: "v1", Value: "true", Description: "first variation"},
			},
		},
		{
			desc:          "success: valid variation",
			variationType: feature.Feature_BOOLEAN,
			id:            id1.String(),
			name:          "v1",
			value:         "true",
			description:   "first variation",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f := &Feature{Feature: &feature.Feature{
				VariationType: p.variationType,
				Variations:    p.variations,
			}}
			err := f.AddVariation(p.id, p.value, p.name, p.description)
			assert.Equal(t, p.expectedErr, err)
		})
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
			expectedErr: nil,
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
			expectedErr: nil,
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
	expected := "variation-D"
	f.AddVariation(expected, "D", "Variation D", "Thing does D")
	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       "variation-A",
			expected: errVariationInUse, // Used in default strategy
		},
		{
			id:       "variation-B",
			expected: errVariationInUse, // Used in default strategy
		},
		{
			id:       "variation-C",
			expected: errVariationInUse, // Has users in target
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
	expectedSize := 3
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %d, actual: %d", expectedSize, actualSize)
	}
}

func TestRemoveVariationUsingRolloutStrategy(t *testing.T) {
	f := makeFeature("test-feature")
	expected := "variation-D"
	f.AddVariation(expected, "D", "Variation D", "Thing does D")
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
				{
					Variation: expected,
					Weight:    0,
				},
			},
		},
	})
	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       "variation-A",
			expected: errVariationInUse, // Used in default strategy with weight > 0
		},
		{
			id:       "variation-B",
			expected: errVariationInUse, // Used in default strategy with weight > 0
		},
		{
			id:       "variation-C",
			expected: errVariationInUse, // Has users in target
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
	expectedSize := 3
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
			expected:      errVariationValueRequired,
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
			t.Parallel()
			f := &Feature{Feature: &feature.Feature{VariationType: p.variationType}}
			assert.Equal(t, p.expected, f.validateVariationValue("", p.value))
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
	t.Parallel()
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	ruleID, _ := uuid.NewUUID() // For rule ID
	patterns := []struct {
		desc                string
		feature             *Feature
		name                *wrapperspb.StringValue
		description         *wrapperspb.StringValue
		tags                *common.StringListValue
		enabled             *wrapperspb.BoolValue
		archived            *wrapperspb.BoolValue
		defaultStrategy     *feature.Strategy
		offVariation        *wrapperspb.StringValue
		resetSamplingSeed   bool
		prerequisiteChanges []*feature.PrerequisiteChange
		targetChanges       []*feature.TargetChange
		ruleChanges         []*feature.RuleChange
		variationChanges    []*feature.VariationChange
		tagChanges          []*feature.TagChange
		expected            *Feature
		expectedErr         error
	}{
		{
			desc: "success: no changes when updating with same values",
			feature: &Feature{Feature: &feature.Feature{
				Name:        "test-feature",
				Description: "test description",
				Tags:        []string{"tag1"},
				Enabled:     true,
				Archived:    false,
				DefaultStrategy: &feature.Strategy{
					Type: feature.Strategy_FIXED,
					FixedStrategy: &feature.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id2.String(),
				Variations: []*feature.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
			}},
			name:        wrapperspb.String("test-feature"),
			description: wrapperspb.String("test description"),
			tags:        &common.StringListValue{Values: []string{"tag1"}},
			enabled:     wrapperspb.Bool(true),
			archived:    wrapperspb.Bool(false),
			defaultStrategy: &feature.Strategy{
				Type: feature.Strategy_FIXED,
				FixedStrategy: &feature.FixedStrategy{
					Variation: id1.String(),
				},
			},
			offVariation: wrapperspb.String(id2.String()),
			expected: &Feature{Feature: &feature.Feature{
				Name:        "test-feature",
				Description: "test description",
				Tags:        []string{"tag1"},
				Enabled:     true,
				Archived:    false,
				DefaultStrategy: &feature.Strategy{
					Type: feature.Strategy_FIXED,
					FixedStrategy: &feature.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id2.String(),
				Variations: []*feature.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
			}},
		},
		{
			desc: "success: version incremented with full replacement fields",
			feature: &Feature{Feature: &feature.Feature{
				Name:        "old-name",
				Description: "old description",
				Tags:        []string{"old-tag"},
				Enabled:     false,
				Archived:    false,
				DefaultStrategy: &feature.Strategy{
					Type: feature.Strategy_FIXED,
					FixedStrategy: &feature.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id2.String(),
				Variations: []*feature.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
			}},
			name:        wrapperspb.String("new-name"),
			description: wrapperspb.String("new description"),
			tags:        &common.StringListValue{Values: []string{"new-tag"}},
			enabled:     wrapperspb.Bool(true),
			archived:    wrapperspb.Bool(false),
			defaultStrategy: &feature.Strategy{
				Type: feature.Strategy_FIXED,
				FixedStrategy: &feature.FixedStrategy{
					Variation: id1.String(),
				},
			},
			offVariation: wrapperspb.String(id1.String()),
			expected: &Feature{Feature: &feature.Feature{
				Name:        "new-name",
				Description: "new description",
				Tags:        []string{"new-tag"},
				Enabled:     true,
				Archived:    false,
				DefaultStrategy: &feature.Strategy{
					Type: feature.Strategy_FIXED,
					FixedStrategy: &feature.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id1.String(),
				Variations: []*feature.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
			}},
		},
		{
			desc: "success: granular updates for variations, rules, prerequisites, targets, and tags",
			feature: &Feature{Feature: &feature.Feature{
				Name: "test-feature",
				Variations: []*feature.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
				Targets: []*feature.Target{
					{Variation: id1.String(), Users: []string{}},
					{Variation: id2.String(), Users: []string{}},
				},
			}},
			variationChanges: []*feature.VariationChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Variation: &feature.Variation{
						Name:  "v3",
						Value: "new-value",
					},
				},
			},
			ruleChanges: []*feature.RuleChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Rule: &feature.Rule{
						Id: ruleID.String(),
						Strategy: &feature.Strategy{
							Type: feature.Strategy_FIXED,
							FixedStrategy: &feature.FixedStrategy{
								Variation: id1.String(),
							},
						},
						Clauses: []*feature.Clause{
							{
								Id:        id1.String(),
								Operator:  feature.Clause_EQUALS,
								Values:    []string{"user1"},
								Attribute: "name",
							},
						},
					},
				},
			},
			prerequisiteChanges: []*feature.PrerequisiteChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Prerequisite: &feature.Prerequisite{
						FeatureId:   "feature-1",
						VariationId: id1.String(),
					},
				},
			},
			targetChanges: []*feature.TargetChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Target: &feature.Target{
						Variation: id1.String(),
						Users:     []string{"user1"},
					},
				},
			},
			tagChanges: []*feature.TagChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Tag:        "new-tag",
				},
			},
			expected: &Feature{Feature: &feature.Feature{
				Name: "test-feature",
				Tags: []string{"new-tag"},
				Variations: []*feature.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
					{Name: "v3", Value: "new-value"},
				},
				Rules: []*feature.Rule{
					{
						Id: ruleID.String(),
						Strategy: &feature.Strategy{
							Type: feature.Strategy_FIXED,
							FixedStrategy: &feature.FixedStrategy{
								Variation: id1.String(),
							},
						},
						Clauses: []*feature.Clause{
							{
								Id:        id1.String(),
								Operator:  feature.Clause_EQUALS,
								Values:    []string{"user1"},
								Attribute: "name",
							},
						},
					},
				},
				Prerequisites: []*feature.Prerequisite{
					{
						FeatureId:   "feature-1",
						VariationId: id1.String(),
					},
				},
				Targets: []*feature.Target{
					{
						Variation: id1.String(),
						Users:     []string{"user1"},
					},
					{
						Variation: id2.String(),
						Users:     []string{},
					},
					{
						Users: []string{},
					},
				},
			}},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			updated, err := p.feature.Update(
				p.name,
				p.description,
				p.tags,
				p.enabled,
				p.archived,
				p.defaultStrategy,
				p.offVariation,
				p.resetSamplingSeed,
				p.prerequisiteChanges,
				p.targetChanges,
				p.ruleChanges,
				p.variationChanges,
				p.tagChanges,
			)
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, p.expected.Name, updated.Name)
			assert.Equal(t, p.expected.Description, updated.Description)
			assert.Equal(t, p.expected.Tags, updated.Tags)
			assert.Equal(t, p.expected.Enabled, updated.Enabled)
			assert.Equal(t, p.expected.Archived, updated.Archived)
			assert.Equal(t, p.expected.OffVariation, updated.OffVariation)

			// Check variations
			var createdVariationID string
			hasVariationCreate := len(p.variationChanges) > 0 && p.variationChanges[0].ChangeType == feature.ChangeType_CREATE
			if hasVariationCreate {
				assert.Equal(t, len(p.expected.Variations), len(updated.Variations))
				for i := range updated.Variations {
					assert.Equal(t, p.expected.Variations[i].Name, updated.Variations[i].Name)
					assert.Equal(t, p.expected.Variations[i].Value, updated.Variations[i].Value)
					assert.Equal(t, p.expected.Variations[i].Description, updated.Variations[i].Description)
				}
				createdVariationID = updated.Variations[len(updated.Variations)-1].Id
			} else {
				assert.Equal(t, p.expected.Variations, updated.Variations)
			}

			// Check targets
			// When creating a new variation, it also adds a new target.
			// So, we need to check the targets after the new variation is created.
			if createdVariationID != "" {
				// New target should be at the end of the array
				newTarget := updated.Targets[len(updated.Targets)-1]
				assert.Equal(t, createdVariationID, newTarget.Variation, "New target should have the created variation ID")
				assert.Empty(t, newTarget.Users, "New target should have no users")

				// Verify other targets remain unchanged
				for i := range updated.Targets[:len(updated.Targets)-1] {
					assert.Equal(t, p.expected.Targets[i].Users, updated.Targets[i].Users)
				}
			} else {
				assert.Equal(t, p.expected.Targets, updated.Targets)
			}
			// Check prerequisites and rules
			assert.Equal(t, p.expected.Prerequisites, updated.Prerequisites)
			assert.Equal(t, p.expected.Rules, updated.Rules)
		})
	}
}

func TestUpdatePrerequisitesGranular(t *testing.T) {
	// Generate valid UUIDs for variations.
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)

	// Baseline feature with no prerequisites.
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
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*proto.Prerequisite{},
				Targets:       []*proto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*proto.Rule{},
				DefaultStrategy: &proto.Strategy{
					Type:          proto.Strategy_FIXED,
					FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	patterns := []struct {
		desc                string
		inputFunc           func() *Feature
		prerequisiteChanges []*proto.PrerequisiteChange
		expectedFunc        func() *Feature
		expectedErr         error
	}{
		{
			desc:      "Prerequisite Create - success",
			inputFunc: genF,
			prerequisiteChanges: []*proto.PrerequisiteChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Prerequisite: &proto.Prerequisite{
						FeatureId:   "f1",
						VariationId: v1.String(),
					},
				},
			},
			expectedFunc: func() *Feature {
				f := genF()
				_ = f.AddPrerequisite("f1", v1.String())
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "Prerequisite Create - error if already exists",
			inputFunc: func() *Feature {
				f := genF()
				_ = f.AddPrerequisite("f1", v1.String())
				return f
			},
			prerequisiteChanges: []*proto.PrerequisiteChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Prerequisite: &proto.Prerequisite{
						FeatureId:   "f1",
						VariationId: v1.String(),
					},
				},
			},
			expectedFunc: func() *Feature {
				// Expect no change because duplicate creation returns an error.
				return genF()
			},
			expectedErr: errPrerequisiteAlreadyExists,
		},
		{
			desc: "Prerequisite Update - error if not found",
			inputFunc: func() *Feature {
				return genF()
			},
			prerequisiteChanges: []*proto.PrerequisiteChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Prerequisite: &proto.Prerequisite{
						FeatureId:   "non-existent",
						VariationId: v2.String(),
					},
				},
			},
			expectedFunc: func() *Feature {
				return genF()
			},
			expectedErr: errPrerequisiteNotFound,
		},
		{
			desc: "Prerequisite Delete - error if not found",
			inputFunc: func() *Feature {
				return genF()
			},
			prerequisiteChanges: []*proto.PrerequisiteChange{
				{
					ChangeType: feature.ChangeType_DELETE,
					Prerequisite: &proto.Prerequisite{
						FeatureId:   "non-existent",
						VariationId: v1.String(),
					},
				},
			},
			expectedFunc: func() *Feature {
				return genF()
			},
			expectedErr: errPrerequisiteNotFound,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.inputFunc().Update(
				nil, nil, nil, nil, nil, nil, nil, false,
				p.prerequisiteChanges, nil, nil, nil, nil,
			)
			assert.Equal(t, p.expectedErr, err, p.desc)
			if err == nil {
				assert.Equal(t, p.expectedFunc().Prerequisites, actual.Prerequisites, p.desc)
			}
		})
	}
}

func TestUpdateTargetsGranular(t *testing.T) {
	// Generate valid UUIDs for variations.
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)

	genF := func() *Feature {
		return &Feature{
			Feature: &proto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1"},
				VariationType: feature.Feature_BOOLEAN,
				Variations: []*proto.Variation{
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*proto.Prerequisite{},
				Targets: []*proto.Target{
					{Variation: v1.String(), Users: []string{"u1"}},
					{Variation: v2.String(), Users: []string{"u2"}},
				},
				Rules: []*proto.Rule{},
				DefaultStrategy: &proto.Strategy{
					Type:          proto.Strategy_FIXED,
					FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	patterns := []struct {
		desc          string
		inputFunc     func() *Feature
		targetChanges []*proto.TargetChange
		expectedFunc  func() *Feature
		expectedErr   error
	}{
		{
			desc:      "Target Create - error: empty target fields",
			inputFunc: genF,
			targetChanges: []*proto.TargetChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Target:     &proto.Target{Variation: "", Users: []string{}},
				},
			},
			expectedFunc: func() *Feature {
				return genF()
			},
			expectedErr: errTargetNotFound,
		},
		{
			desc:      "Target Update - error: target not found",
			inputFunc: genF,
			targetChanges: []*proto.TargetChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Target:     &proto.Target{Variation: "non-existent", Users: []string{"u-new"}},
				},
			},
			expectedFunc: func() *Feature {
				return genF()
			},
			expectedErr: errTargetNotFound,
		},
		{
			desc:      "Target Delete - error: target not found",
			inputFunc: genF,
			targetChanges: []*proto.TargetChange{
				{
					ChangeType: feature.ChangeType_DELETE,
					Target:     &proto.Target{Variation: "non-existent"},
				},
			},
			expectedFunc: func() *Feature {
				return genF()
			},
			expectedErr: errTargetNotFound,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.inputFunc().Update(
				nil, nil, nil, nil, nil, nil, nil, false,
				nil, p.targetChanges, nil, nil, nil,
			)
			assert.Equal(t, p.expectedErr, err, p.desc)
			if err == nil {
				assert.Equal(t, p.expectedFunc().Targets, actual.Targets, p.desc)
			}
		})
	}
}

func TestUpdateRulesGranular(t *testing.T) {
	// Generate valid UUIDs.
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)
	ruleID, err := uuid.NewUUID()
	require.NoError(t, err)
	clauseID, err := uuid.NewUUID()
	require.NoError(t, err)

	// genF returns a baseline Feature with no rules.
	genF := func() *Feature {
		return &Feature{
			Feature: &proto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1"},
				VariationType: feature.Feature_BOOLEAN,
				Variations: []*proto.Variation{
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*proto.Prerequisite{},
				Targets:       []*proto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*proto.Rule{},
				DefaultStrategy: &proto.Strategy{
					Type:          proto.Strategy_FIXED,
					FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	// Define patterns for rule granular updates including clause validations.
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		ruleChanges  []*proto.RuleChange
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc:      "Rule Create - error: rule required",
			inputFunc: genF,
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Rule:       nil,
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errRuleRequired,
		},
		{
			desc:      "Rule Create - error: nil strategy",
			inputFunc: genF,
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Rule: &proto.Rule{
						Id:       ruleID.String(),
						Strategy: nil, // This should trigger errStrategyRequired
						Clauses: []*proto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "attr",
								Operator:  feature.Clause_EQUALS,
								Values:    []string{"val"},
							},
						},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errStrategyRequired,
		},
		{
			desc:      "Rule Create - error: no clauses",
			inputFunc: genF,
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Rule: &proto.Rule{
						Id: ruleID.String(),
						Strategy: &proto.Strategy{
							Type:          feature.Strategy_FIXED,
							FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*proto.Clause{},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errors.New("feature: rule must have at least one clause"),
		},
		{
			desc:      "Rule Create - error: clause attribute not empty for SEGMENT operator",
			inputFunc: genF,
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Rule: &proto.Rule{
						Id: ruleID.String(),
						Strategy: &proto.Strategy{
							Type:          feature.Strategy_FIXED,
							FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*proto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "non-empty", // Not allowed for SEGMENT operator
								Operator:  feature.Clause_SEGMENT,
								Values:    []string{"val"},
							},
						},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errClauseAttributeNotEmpty,
		},
		{
			desc:      "Rule Create - error: clause values empty for SEGMENT operator",
			inputFunc: genF,
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Rule: &proto.Rule{
						Id: ruleID.String(),
						Strategy: &proto.Strategy{
							Type:          feature.Strategy_FIXED,
							FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*proto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "", // Correct for SEGMENT operator
								Operator:  feature.Clause_SEGMENT,
								Values:    []string{}, // Missing values
							},
						},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errClauseValuesEmpty,
		},
		{
			desc:      "Rule Update - error: rule not found",
			inputFunc: genF,
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Rule: &proto.Rule{
						Id: "non-existent",
						Strategy: &proto.Strategy{
							Type:          feature.Strategy_FIXED,
							FixedStrategy: &proto.FixedStrategy{Variation: v2.String()},
						},
						Clauses: []*proto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "attr",
								Operator:  feature.Clause_EQUALS,
								Values:    []string{"val-updated"},
							},
						},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errors.New("uuid: format must be an uuid version 4"),
		},
		{
			desc: "Rule Update - error: clause attribute empty for non-SEGMENT operator",
			inputFunc: func() *Feature {
				f := genF()
				// Pre-add a valid rule.
				_ = f.AddRule(&proto.Rule{
					Id: ruleID.String(),
					Strategy: &proto.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
					},
					Clauses: []*proto.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"val"},
						},
					},
				})
				return f
			},
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Rule: &proto.Rule{
						Id: ruleID.String(),
						Strategy: &proto.Strategy{
							Type:          feature.Strategy_FIXED,
							FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*proto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "", // empty attribute for non-SEGMENT operator not allowed
								Operator:  feature.Clause_EQUALS,
								Values:    []string{"val-updated"},
							},
						},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errClauseAttributeEmpty,
		},
		{
			desc: "Rule Update - error: clause values empty for non-SEGMENT operator",
			inputFunc: func() *Feature {
				f := genF()
				// Pre-add a valid rule.
				_ = f.AddRule(&proto.Rule{
					Id: ruleID.String(),
					Strategy: &proto.Strategy{
						Type:          feature.Strategy_FIXED,
						FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
					},
					Clauses: []*proto.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  feature.Clause_EQUALS,
							Values:    []string{"val"},
						},
					},
				})
				return f
			},
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Rule: &proto.Rule{
						Id: ruleID.String(),
						Strategy: &proto.Strategy{
							Type:          feature.Strategy_FIXED,
							FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*proto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "attr",
								Operator:  feature.Clause_EQUALS,
								Values:    []string{}, // empty values not allowed
							},
						},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errClauseValuesEmpty,
		},
		{
			desc:      "Rule Delete - error: empty rule id",
			inputFunc: genF,
			ruleChanges: []*proto.RuleChange{
				{
					ChangeType: feature.ChangeType_DELETE,
					Rule:       &proto.Rule{Id: ""},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errRuleIDRequired,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.inputFunc().Update(
				nil, nil, nil, nil, nil, nil, nil, false, // basic fields
				nil, nil, p.ruleChanges, nil, nil, // granular change lists
			)
			if p.expectedErr != nil {
				require.Error(t, err, p.desc)
				assert.Contains(t, err.Error(), p.expectedErr.Error(), p.desc)
			} else {
				require.NoError(t, err, p.desc)
				assert.Equal(t, p.expectedFunc().Rules, actual.Rules, p.desc)
			}
		})
	}
}

func TestUpdateVariationsGranular(t *testing.T) {
	// Generate valid UUIDs for variations.
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)
	v3, err := uuid.NewUUID() // for create success
	require.NoError(t, err)

	// Baseline generator for JSON type.
	genFJSON := func() *Feature {
		return &Feature{
			Feature: &proto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1"},
				VariationType: feature.Feature_JSON,
				Variations: []*proto.Variation{
					{Id: v1.String(), Value: `{"key": "value1"}`, Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: `{"key": "value2"}`, Name: "n2", Description: "d2"},
				},
				Prerequisites: []*proto.Prerequisite{},
				Targets:       []*proto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*proto.Rule{},
				DefaultStrategy: &proto.Strategy{
					Type:          proto.Strategy_FIXED,
					FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	// Define test patterns.
	patterns := []struct {
		desc             string
		inputFunc        func() *Feature
		variationChanges []*proto.VariationChange
		expectedFunc     func() *Feature
		expectedErr      error
	}{
		{
			desc:      "Variation Create - success",
			inputFunc: genFJSON,
			variationChanges: []*proto.VariationChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Variation: &proto.Variation{
						Value:       `{"key": "value3"}`,
						Name:        "n3",
						Description: "d3",
					},
				},
			},
			expectedFunc: func() *Feature {
				f := genFJSON()
				// Add the new variation directly
				f.Variations = append(f.Variations, &proto.Variation{
					Id:          v3.String(),
					Value:       `{"key": "value3"}`,
					Name:        "n3",
					Description: "d3",
				})
				// Add corresponding target for the new variation
				f.Targets = append(f.Targets, &proto.Target{
					Variation: v3.String(),
					Users:     []string{},
				})
				return f
			},
			expectedErr: nil,
		},
		{
			desc:      "Variation Update - success",
			inputFunc: genFJSON,
			variationChanges: []*proto.VariationChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Variation: &proto.Variation{
						Id:          v1.String(),
						Value:       `{"key": "updated-value1"}`,
						Name:        "n1-updated",
						Description: "d1-updated",
					},
				},
			},
			expectedFunc: func() *Feature {
				f := genFJSON()
				// Update the variation directly
				f.Variations[0].Value = `{"key": "updated-value1"}`
				f.Variations[0].Name = "n1-updated"
				f.Variations[0].Description = "d1-updated"
				return f
			},
			expectedErr: nil,
		},
		{
			desc: "Variation Delete - success",
			inputFunc: func() *Feature {
				f := genFJSON()
				// Add the variation to be deleted
				f.Variations = append(f.Variations, &proto.Variation{
					Id:          v3.String(),
					Value:       `{"key": "value3"}`,
					Name:        "n3",
					Description: "d3",
				})
				f.Targets = append(f.Targets, &proto.Target{
					Variation: v3.String(),
					Users:     []string{},
				})
				return f
			},
			variationChanges: []*proto.VariationChange{
				{
					ChangeType: feature.ChangeType_DELETE,
					Variation:  &proto.Variation{Id: v3.String()},
				},
			},
			expectedFunc: func() *Feature {
				return genFJSON()
			},
			expectedErr: nil,
		},
		{
			desc:      "Variation Update - error: nil variation",
			inputFunc: genFJSON,
			variationChanges: []*proto.VariationChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Variation:  nil,
				},
			},
			expectedFunc: func() *Feature { return genFJSON() },
			expectedErr:  errVariationRequired,
		},
		{
			desc:      "Variation Update - error: empty name",
			inputFunc: genFJSON,
			variationChanges: []*proto.VariationChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Variation: &proto.Variation{
						Id:          v1.String(),
						Value:       `{"key": "value1"}`,
						Name:        "",
						Description: "d1-updated",
					},
				},
			},
			expectedFunc: func() *Feature { return genFJSON() },
			expectedErr:  errVariationNameRequired,
		},
		{
			desc:      "Variation Update - success: valid JSON object",
			inputFunc: genFJSON,
			variationChanges: []*proto.VariationChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Variation: &proto.Variation{
						Id:          v1.String(),
						Value:       `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`,
						Name:        "n1-updated",
						Description: "d1-updated",
					},
				},
			},
			expectedFunc: func() *Feature {
				f := genFJSON()
				// Update the variation directly
				f.Variations[0].Value = `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`
				f.Variations[0].Name = "n1-updated"
				f.Variations[0].Description = "d1-updated"
				return f
			},
			expectedErr: nil,
		},
		{
			desc:      "Variation Update - success: valid JSON array",
			inputFunc: genFJSON,
			variationChanges: []*proto.VariationChange{
				{
					ChangeType: feature.ChangeType_UPDATE,
					Variation: &proto.Variation{
						Id:          v1.String(),
						Value:       `[{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}]`,
						Name:        "n1-updated",
						Description: "d1-updated",
					},
				},
			},
			expectedFunc: func() *Feature {
				f := genFJSON()
				// Update the variation directly
				f.Variations[0].Value = `[{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}]`
				f.Variations[0].Name = "n1-updated"
				f.Variations[0].Description = "d1-updated"
				return f
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.inputFunc().Update(
				nil, nil, nil, nil, nil, nil, nil, false, // basic fields
				nil, nil, nil, p.variationChanges, nil, // granular change lists
			)
			if p.expectedErr != nil {
				require.Error(t, err, p.desc)
				assert.Contains(t, err.Error(), p.expectedErr.Error(), p.desc)
			} else {
				require.NoError(t, err, p.desc)
				expected := p.expectedFunc()
				// Compare variations without considering IDs for newly created ones
				assert.Equal(t, len(expected.Variations), len(actual.Variations), p.desc)
				for i, expectedVar := range expected.Variations {
					actualVar := actual.Variations[i]
					assert.Equal(t, expectedVar.Value, actualVar.Value, p.desc)
					assert.Equal(t, expectedVar.Name, actualVar.Name, p.desc)
					assert.Equal(t, expectedVar.Description, actualVar.Description, p.desc)
				}
			}
		})
	}
}

func TestUpdateTagsGranular(t *testing.T) {
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)
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
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*proto.Prerequisite{},
				Targets:       []*proto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*proto.Rule{},
				DefaultStrategy: &proto.Strategy{
					Type:          proto.Strategy_FIXED,
					FixedStrategy: &proto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		tagChanges   []*proto.TagChange
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc:      "Tag Create - error: duplicate create should not add duplicate",
			inputFunc: genF,
			tagChanges: []*proto.TagChange{
				{
					ChangeType: feature.ChangeType_CREATE,
					Tag:        "new-tag",
				},
				{
					ChangeType: feature.ChangeType_CREATE,
					Tag:        "new-tag",
				},
			},
			expectedFunc: func() *Feature {
				f := genF()
				_ = f.AddTag("new-tag")
				return f
			},
			expectedErr: nil,
		},
		{
			desc:      "Tag Delete - error: tag not found",
			inputFunc: genF,
			tagChanges: []*proto.TagChange{
				{
					ChangeType: feature.ChangeType_DELETE,
					Tag:        "non-existent-tag",
				},
			},
			expectedFunc: func() *Feature {
				return genF()
			},
			expectedErr: errors.New("feature: value not found"),
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := p.inputFunc().Update(
				nil, nil, nil, nil, nil, nil, nil, false,
				nil, nil, nil, nil, p.tagChanges,
			)
			if p.expectedErr != nil {
				assert.Error(t, err, p.desc)
				assert.Contains(t, err.Error(), p.expectedErr.Error(), p.desc)
			} else {
				assert.NoError(t, err, p.desc)
				tagSet := make(map[string]struct{})
				for _, tag := range actual.Tags {
					tagSet[tag] = struct{}{}
				}
				assert.Equal(t, len(actual.Tags), len(tagSet), p.desc)
				assert.Equal(t, p.expectedFunc().Tags, actual.Tags, p.desc)
			}
		})
	}
}

func TestValidateStrategy(t *testing.T) {
	t.Parallel()
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	variations := []*feature.Variation{
		{Id: id1.String(), Value: "true", Name: "n1", Description: "d1"},
		{Id: id2.String(), Value: "false", Name: "n2", Description: "d2"},
	}
	tests := []struct {
		desc        string
		strategy    *feature.Strategy
		variations  []*feature.Variation
		expectedErr error
	}{
		{
			desc: "success: fixed strategy",
			strategy: &feature.Strategy{
				Type:          feature.Strategy_FIXED,
				FixedStrategy: &feature.FixedStrategy{Variation: id1.String()},
			},
			variations:  variations,
			expectedErr: nil,
		},
		{
			desc: "success: rollout strategy",
			strategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100},
					},
				},
			},
			variations:  variations,
			expectedErr: nil,
		},
		{
			desc:        "fail: strategy is nil",
			strategy:    nil,
			variations:  variations,
			expectedErr: errStrategyRequired,
		},
		{
			desc: "fail: fixed strategy with non-existent variation",
			strategy: &feature.Strategy{
				Type:          feature.Strategy_FIXED,
				FixedStrategy: &feature.FixedStrategy{Variation: "non-existent"},
			},
			variations:  variations,
			expectedErr: errVariationNotFound,
		},
		{
			desc: "fail: rollout strategy with non-existent variation",
			strategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{
						{Variation: "non-existent", Weight: 100},
					},
				},
			},
			variations:  variations,
			expectedErr: errVariationNotFound,
		},
		{
			desc: "fail: unsupported strategy type",
			strategy: &feature.Strategy{
				Type: 999,
			},
			variations:  variations,
			expectedErr: errUnsupportedStrategy,
		},
		{
			desc: "fail: fixed strategy is nil",
			strategy: &feature.Strategy{
				Type: feature.Strategy_FIXED,
			},
			variations:  variations,
			expectedErr: ErrRuleStrategyCannotBeEmpty,
		},
		{
			desc: "fail: rollout strategy is nil",
			strategy: &feature.Strategy{
				Type: feature.Strategy_ROLLOUT,
			},
			variations:  variations,
			expectedErr: ErrRuleStrategyCannotBeEmpty,
		},
		{
			desc: "fail: both strategies are set",
			strategy: &feature.Strategy{
				Type:          feature.Strategy_FIXED,
				FixedStrategy: &feature.FixedStrategy{Variation: id1.String()},
				RolloutStrategy: &feature.RolloutStrategy{
					Variations: []*feature.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100},
					},
				},
			},
			variations:  variations,
			expectedErr: ErrDefaultStrategyCannotBeBothFixedAndRollout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			err := validateStrategy(tt.strategy, tt.variations)
			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
