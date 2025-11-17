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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	"github.com/bucketeer-io/bucketeer/v2/proto/common"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func makeFeature(id string) *Feature {
	return &Feature{
		Feature: &ftproto.Feature{
			Id:            id,
			Name:          "test feature",
			Version:       1,
			Enabled:       true,
			CreatedAt:     time.Now().Unix(),
			VariationType: ftproto.Feature_STRING,
			Variations: []*ftproto.Variation{
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
			Targets: []*ftproto.Target{
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
			Rules: []*ftproto.Rule{
				{
					Id: "rule-1",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        "clause-1",
							Attribute: "name",
							Operator:  ftproto.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
				{
					Id: "rule-2",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "variation-B",
						},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        "clause-2",
							Attribute: "name",
							Operator:  ftproto.Clause_EQUALS,
							Values: []string{
								"user3",
								"user4",
							},
						},
					},
				},
			},
			DefaultStrategy: &ftproto.Strategy{
				Type: ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{
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
		variationType            ftproto.Feature_VariationType
		variations               []*ftproto.Variation
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
			variationType: ftproto.Feature_BOOLEAN,
			variations: []*ftproto.Variation{
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
			variationType: ftproto.Feature_BOOLEAN,
			variations: []*ftproto.Variation{
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
			variationType: ftproto.Feature_BOOLEAN,
			variations: []*ftproto.Variation{
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
			variationType: ftproto.Feature_BOOLEAN,
			variations: []*ftproto.Variation{
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
		variationType ftproto.Feature_VariationType
		id            string
		name          string
		value         string
		description   string
		expectedErr   error
		variations    []*ftproto.Variation
	}{
		{
			desc:          "fail: empty name",
			variationType: ftproto.Feature_BOOLEAN,
			id:            id1.String(),
			name:          "",
			value:         "true",
			description:   "first variation",
			expectedErr:   errVariationNameRequired,
			variations: []*ftproto.Variation{
				{Id: id1.String(), Name: "v1", Value: "true", Description: "first variation"},
			},
		},
		{
			desc:          "fail: empty value",
			variationType: ftproto.Feature_BOOLEAN,
			id:            id1.String(),
			name:          "v1",
			value:         "",
			description:   "first variation",
			expectedErr:   errVariationValueRequired,
		},
		{
			desc:          "fail: duplicate value",
			variationType: ftproto.Feature_BOOLEAN,
			id:            id2.String(),
			name:          "v2",
			value:         "true", // same value as first variation
			description:   "second variation",
			expectedErr:   errVariationValueUnique,
			variations: []*ftproto.Variation{
				{Id: id1.String(), Name: "v1", Value: "true", Description: "first variation"},
			},
		},
		{
			desc:          "success: valid variation",
			variationType: ftproto.Feature_BOOLEAN,
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
			f := &Feature{Feature: &ftproto.Feature{
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
		strategy    *ftproto.Strategy
		clauses     []*ftproto.Clause
		expectedErr bool
	}{
		{
			desc:     "fail: add rule with nil strategy",
			id:       "rule-2",
			strategy: nil,
			clauses: []*ftproto.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  ftproto.Clause_EQUALS,
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
			strategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: ""},
			},
			clauses: []*ftproto.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  ftproto.Clause_EQUALS,
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
			strategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: f.Variations[0].Id},
			},
			clauses: []*ftproto.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  ftproto.Clause_EQUALS,
					Values: []string{
						"user1",
					},
				},
			},
			expectedErr: false,
		},
	}
	for _, p := range patterns {
		rule := &ftproto.Rule{
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
		rule        *ftproto.Rule
		expectedErr bool
	}{
		{
			desc: "fail: rule already exists",
			rule: &ftproto.Rule{
				Id:       "rule-2",
				Strategy: nil,
			},
			expectedErr: true,
		},
		{
			desc: "fail: variation not found",
			rule: &ftproto.Rule{
				Id: "rule-3",
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
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
			rule: &ftproto.Rule{
				Id: "rule-3",
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
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
				Clauses: []*ftproto.Clause{{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  ftproto.Clause_EQUALS,
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
	expected := &ftproto.Strategy{
		Type:          ftproto.Strategy_FIXED,
		FixedStrategy: &ftproto.FixedStrategy{Variation: vID},
	}
	patterns := []*struct {
		ruleID   string
		strategy *ftproto.Strategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: expected,
			expected: errRuleNotFound,
		},
		{
			ruleID: rID,
			strategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: ""},
			},
			expected: errVariationNotFound,
		},
		{
			ruleID: rID,
			strategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: "variation-D"},
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
	expected := &ftproto.Strategy{
		Type: ftproto.Strategy_ROLLOUT,
		RolloutStrategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
			{
				Variation: vID1,
				Weight:    30000,
			},
			{
				Variation: vID2,
				Weight:    70000,
			},
		}},
	}
	patterns := []*struct {
		ruleID   string
		strategy *ftproto.Strategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: expected,
			expected: errRuleNotFound,
		},
		{
			ruleID: rID,
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
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
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
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
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
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
	if !reflect.DeepEqual(expected.RolloutStrategy, r.Strategy.RolloutStrategy) {
		t.Fatalf("Strategy is not equal. Expected: %v, actual: %v", expected.RolloutStrategy, r.Strategy.RolloutStrategy)
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
		operator    ftproto.Clause_Operator
		ruleIdx     int
		idx         int
		expectedErr error
	}{
		{
			rule:        "rule-1",
			clause:      "clause-1",
			operator:    ftproto.Clause_IN,
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
			expected: ErrVariationInUse, // Used in default strategy
		},
		{
			id:       "variation-B",
			expected: ErrVariationInUse, // Used in default strategy
		},
		{
			id:       "variation-C",
			expected: ErrVariationInUse, // Has users in target
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
	f.ChangeDefaultStrategy(&ftproto.Strategy{
		Type: ftproto.Strategy_ROLLOUT,
		RolloutStrategy: &ftproto.RolloutStrategy{
			Variations: []*ftproto.RolloutStrategy_Variation{
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
			expected: ErrVariationInUse, // Used in default strategy with weight > 0
		},
		{
			id:       "variation-B",
			expected: ErrVariationInUse, // Used in default strategy with weight > 0
		},
		{
			id:       "variation-C",
			expected: ErrVariationInUse, // Has users in target
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
			expected: ErrVariationInUse,
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

func TestRemoveVariationComprehensiveCleanup(t *testing.T) {
	t.Parallel()
	f := makeFeature("test-feature")
	expected := "variation-D"
	f.AddVariation(expected, "D", "Variation D", "Thing does D")

	// Set up rollout strategy in default strategy with the variation (weight=0 so it can be removed)
	f.ChangeDefaultStrategy(&ftproto.Strategy{
		Type: ftproto.Strategy_ROLLOUT,
		RolloutStrategy: &ftproto.RolloutStrategy{
			Variations: []*ftproto.RolloutStrategy_Variation{
				{
					Variation: "variation-A",
					Weight:    100000,
				},
				{
					Variation: expected,
					Weight:    0, // Weight 0 means not "in use" so can be removed
				},
			},
		},
	})

	// Add a rule with rollout strategy containing the variation
	rule := &ftproto.Rule{
		Id: "test-rule-rollout",
		Strategy: &ftproto.Strategy{
			Type: ftproto.Strategy_ROLLOUT,
			RolloutStrategy: &ftproto.RolloutStrategy{
				Variations: []*ftproto.RolloutStrategy_Variation{
					{
						Variation: "variation-B",
						Weight:    50000,
					},
					{
						Variation: expected,
						Weight:    0, // Weight 0 means not "in use" so can be removed
					},
				},
			},
		},
		Clauses: []*ftproto.Clause{
			{
				Id:        "clause-1",
				Attribute: "user_id",
				Operator:  ftproto.Clause_EQUALS,
				Values:    []string{"user-1"},
			},
		},
	}
	f.AddRule(rule)

	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       "variation-A",
			expected: ErrVariationInUse, // Used in default strategy with weight > 0
		},
		{
			id:       "variation-B",
			expected: ErrVariationInUse, // Used in rule strategy with weight > 0
		},
		{
			id:       "variation-C",
			expected: ErrVariationInUse, // Has users in target
		},
		{
			id:       expected,
			expected: nil, // Can be removed (weight=0 in all strategies)
		},
	}

	for i, p := range patterns {
		err := f.RemoveVariation(p.id)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}

	// Verify complete cleanup for successfully removed variation
	if _, err := f.findVariationIndex(expected); err == nil {
		t.Fatalf("Variation not deleted from Variations. Actual: %v", f.Variations)
	}
	if _, err := f.findTarget(expected); err == nil {
		t.Fatalf("Target not deleted. Actual: %v", f.Targets)
	}

	// Verify variation removed from default strategy rollout
	for _, v := range f.DefaultStrategy.RolloutStrategy.Variations {
		if v.Variation == expected {
			t.Fatalf("Variation not removed from default strategy. Actual: %v", f.DefaultStrategy.RolloutStrategy.Variations)
		}
	}

	// Verify variation removed from rule rollout strategy
	for _, r := range f.Rules {
		if r.Id == "test-rule-rollout" && r.Strategy.Type == ftproto.Strategy_ROLLOUT {
			for _, v := range r.Strategy.RolloutStrategy.Variations {
				if v.Variation == expected {
					t.Fatalf("Variation not removed from rule strategy. Actual: %v", r.Strategy.RolloutStrategy.Variations)
				}
			}
		}
	}

	actualSize := len(f.Variations)
	expectedSize := 3
	if expectedSize != actualSize {
		t.Fatalf("Different sizes. Expected: %d, actual: %d", expectedSize, actualSize)
	}
}

func TestRemoveVariationMultipleInstancesInRollout(t *testing.T) {
	t.Parallel()
	f := makeFeature("test-feature")
	expected := "variation-D"
	f.AddVariation(expected, "D", "Variation D", "Thing does D")

	// Create a rollout strategy with MULTIPLE instances of the same variation (edge case)
	f.ChangeDefaultStrategy(&ftproto.Strategy{
		Type: ftproto.Strategy_ROLLOUT,
		RolloutStrategy: &ftproto.RolloutStrategy{
			Variations: []*ftproto.RolloutStrategy_Variation{
				{
					Variation: "variation-A",
					Weight:    50000,
				},
				{
					Variation: expected,
					Weight:    0, // First instance with weight 0
				},
				{
					Variation: "variation-B",
					Weight:    50000,
				},
				{
					Variation: expected,
					Weight:    0, // Second instance with weight 0
				},
			},
		},
	})

	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       expected,
			expected: nil, // Can be removed (all instances have weight=0)
		},
	}

	// Verify multiple instances exist before removal
	instanceCount := 0
	for _, v := range f.DefaultStrategy.RolloutStrategy.Variations {
		if v.Variation == expected {
			instanceCount++
		}
	}
	if instanceCount != 2 {
		t.Fatalf("Expected 2 instances before removal, got %d", instanceCount)
	}

	for i, p := range patterns {
		err := f.RemoveVariation(p.id)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}

	// Verify ALL instances are removed (this would fail with the old single-remove bug)
	instanceCount = 0
	for _, v := range f.DefaultStrategy.RolloutStrategy.Variations {
		if v.Variation == expected {
			instanceCount++
		}
	}
	if instanceCount != 0 {
		t.Fatalf("Expected 0 instances after removal, got %d", instanceCount)
	}

	// Verify other variations are still present
	actualRolloutSize := len(f.DefaultStrategy.RolloutStrategy.Variations)
	expectedRolloutSize := 2
	if expectedRolloutSize != actualRolloutSize {
		t.Fatalf("Different rollout sizes. Expected: %d, actual: %d", expectedRolloutSize, actualRolloutSize)
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
		err := f.ChangeFixedStrategy(p.ruleID, &ftproto.FixedStrategy{Variation: p.variationID})
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
	expected := &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
		{
			Variation: vID1,
			Weight:    30000,
		},
		{
			Variation: vID2,
			Weight:    70000,
		},
	}}
	patterns := []*struct {
		ruleID   string
		strategy *ftproto.RolloutStrategy
		expected error
	}{
		{
			ruleID:   "",
			strategy: &ftproto.RolloutStrategy{},
			expected: errRuleNotFound,
		},
		{
			ruleID: rID,
			strategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
				{
					Variation: "",
					Weight:    30000,
				},
				{
					Variation: vID2,
					Weight:    70000,
				},
			}},
			expected: errVariationNotFound,
		},
		{
			ruleID: rID,
			strategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    30000,
				},
				{
					Variation: "",
					Weight:    70000,
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
		{
			ruleID: rID,
			strategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    30000, // 30%
				},
				{
					Variation: vID2,
					Weight:    40000, // 40%
				},
				// Total: 70000 (70%) - invalid!
			}},
			expected: ErrInvalidVariationWeightTotal,
		},
		{
			ruleID: rID,
			strategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    60000, // 60%
				},
				{
					Variation: vID2,
					Weight:    50000, // 50%
				},
				// Total: 110000 (110%) - invalid!
			}},
			expected: ErrInvalidVariationWeightTotal,
		},
		{
			ruleID: rID,
			strategy: &ftproto.RolloutStrategy{Variations: []*ftproto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    30, // Old test format - invalid!
				},
				{
					Variation: vID2,
					Weight:    70, // Old test format - invalid!
				},
			}},
			expected: ErrInvalidVariationWeightTotal,
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
			feature: &Feature{Feature: &ftproto.Feature{
				LastUsedInfo: &ftproto.FeatureLastUsedInfo{
					LastUsedAt: t1.Unix(),
				},
			}},
			input:    t2,
			expected: false,
		},
		{
			desc: "true",
			feature: &Feature{Feature: &ftproto.Feature{
				LastUsedInfo: &ftproto.FeatureLastUsedInfo{
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
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)
	patterns := []struct {
		desc          string
		variationType ftproto.Feature_VariationType
		value         string
		expected      error
	}{
		{
			desc:          "invalid bool",
			variationType: ftproto.Feature_BOOLEAN,
			value:         "hoge",
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "empty string",
			variationType: ftproto.Feature_JSON,
			value:         "",
			expected:      errVariationValueRequired,
		},
		{
			desc:          "invalid number",
			variationType: ftproto.Feature_NUMBER,
			value:         `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`,
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "invalid json",
			variationType: ftproto.Feature_JSON,
			value:         "true",
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "valid bool",
			variationType: ftproto.Feature_BOOLEAN,
			value:         "true",
			expected:      nil,
		},
		{
			desc:          "valid number float",
			variationType: ftproto.Feature_NUMBER,
			value:         "1.23",
			expected:      nil,
		},
		{
			desc:          "valid number int",
			variationType: ftproto.Feature_NUMBER,
			value:         "123",
			expected:      nil,
		},
		{
			desc:          "valid json",
			variationType: ftproto.Feature_JSON,
			value:         `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`,
			expected:      nil,
		},
		{
			desc:          "valid json array",
			variationType: ftproto.Feature_JSON,
			value:         `[{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}]`,
			expected:      nil,
		},
		{
			desc:          "valid string",
			variationType: ftproto.Feature_STRING,
			value:         `{"foo":"foo","fee":20,"hoo": [1, "lee", null], "boo": true}`,
			expected:      nil,
		},
		{
			desc:          "valid yaml - simple",
			variationType: ftproto.Feature_YAML,
			value: `name: John Doe
age: 30
active: true`,
			expected: nil,
		},
		{
			desc:          "valid yaml - nested objects",
			variationType: ftproto.Feature_YAML,
			value: `config:
  database:
    host: localhost
    port: 5432
  cache:
    enabled: true
    ttl: 3600`,
			expected: nil,
		},
		{
			desc:          "valid yaml - arrays",
			variationType: ftproto.Feature_YAML,
			value: `items:
  - id: 1
    name: Item 1
  - id: 2
    name: Item 2`,
			expected: nil,
		},
		{
			desc:          "valid yaml - with comments",
			variationType: ftproto.Feature_YAML,
			value: `# Configuration
name: Test Config
# Settings
settings:
  enabled: true  # Enable feature
  timeout: 30    # Timeout in seconds`,
			expected: nil,
		},
		{
			desc:          "valid yaml - mixed types",
			variationType: ftproto.Feature_YAML,
			value: `string: hello
number: 42
float: 3.14
boolean: true
null_value: null
list:
  - one
  - two
  - three
object:
  nested: value`,
			expected: nil,
		},
		{
			desc:          "invalid yaml - malformed",
			variationType: ftproto.Feature_YAML,
			value:         `invalid: yaml: [unclosed`,
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "invalid yaml - tab indentation",
			variationType: ftproto.Feature_YAML,
			value:         "config:\n\tkey: value",
			expected:      errVariationTypeUnmatched,
		},
		{
			desc:          "invalid yaml - unbalanced brackets",
			variationType: ftproto.Feature_YAML,
			value: `list: [1, 2, 3
incomplete`,
			expected: errVariationTypeUnmatched,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f := &Feature{Feature: &ftproto.Feature{
				VariationType: p.variationType,
				Variations: []*ftproto.Variation{
					{Id: v1.String(), Value: "value-1"},
					{Id: v2.String(), Value: "value-2"},
				},
			}}
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
		defaultStrategy   *ftproto.Strategy
		rules             []*ftproto.Rule
	}{
		{
			maintainer:        "sample@example.com",
			offVariationIndex: 2,
			expectedEnabled:   false,
			expectedVersion:   int32(1),
			defaultStrategy: &ftproto.Strategy{
				Type: ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{
					Variation: "variation-B",
				},
			},
			rules: []*ftproto.Rule{
				{
					Id: "rule-1",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        "clause-1",
							Attribute: "name",
							Operator:  ftproto.Clause_EQUALS,
							Values: []string{
								"user1",
								"user2",
							},
						},
					},
				},
				{
					Id: "rule-2",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "variation-B",
						},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        "clause-2",
							Attribute: "name",
							Operator:  ftproto.Clause_EQUALS,
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
			defaultStrategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
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
			rules: []*ftproto.Rule{
				{
					Id: "rule-1",
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_ROLLOUT,
						RolloutStrategy: &ftproto.RolloutStrategy{
							Variations: []*ftproto.RolloutStrategy_Variation{
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
					Strategy: &ftproto.Strategy{
						Type: ftproto.Strategy_ROLLOUT,
						RolloutStrategy: &ftproto.RolloutStrategy{
							Variations: []*ftproto.RolloutStrategy_Variation{
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
		if actual.DefaultStrategy.Type == ftproto.Strategy_FIXED {
			assert.Equal(t, actual.Variations[1].Id, actual.DefaultStrategy.FixedStrategy.Variation)
		} else {
			for i := range actual.Variations {
				assert.Equal(t, actual.Variations[i].Id, actual.DefaultStrategy.RolloutStrategy.Variations[i].Variation)
			}
		}
		assert.NotNil(t, actual.Prerequisites)
		assert.Equal(t, len(actual.Prerequisites), 0)
		for i := range actual.Rules {
			if actual.Rules[i].Strategy.Type == ftproto.Strategy_FIXED {
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
		feature  *ftproto.Feature
		expected []string
	}{
		{
			feature:  &ftproto.Feature{},
			expected: []string{},
		},
		{
			feature: &ftproto.Feature{
				Prerequisites: []*ftproto.Prerequisite{
					{FeatureId: "feature-1"},
				},
			},
			expected: []string{"feature-1"},
		},
		{
			feature: &ftproto.Feature{
				Prerequisites: []*ftproto.Prerequisite{
					{FeatureId: "feature-1"},
					{FeatureId: "feature-2"},
				},
			},
			expected: []string{"feature-1", "feature-2"},
		},
		{
			feature: &ftproto.Feature{
				Rules: []*ftproto.Rule{
					{
						Clauses: []*ftproto.Clause{
							{Attribute: "feature-1", Operator: ftproto.Clause_FEATURE_FLAG},
						},
					},
				},
			},
			expected: []string{"feature-1"},
		},
		{
			feature: &ftproto.Feature{
				Rules: []*ftproto.Rule{
					{
						Clauses: []*ftproto.Clause{
							{Attribute: "feature-1", Operator: ftproto.Clause_FEATURE_FLAG},
							{Attribute: "feature-2", Operator: ftproto.Clause_FEATURE_FLAG},
						},
					},
				},
			},
			expected: []string{"feature-1", "feature-2"},
		},
		{
			feature: &ftproto.Feature{
				Prerequisites: []*ftproto.Prerequisite{
					{FeatureId: "feature-1"},
				},
				Rules: []*ftproto.Rule{
					{
						Clauses: []*ftproto.Clause{
							{Attribute: "feature-2", Operator: ftproto.Clause_FEATURE_FLAG},
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

func TestValidateVariationUsage(t *testing.T) {
	t.Parallel()

	variationID1 := "variation-1"
	variationID2 := "variation-2"
	variationValue1 := "true"

	patterns := []struct {
		desc              string
		features          []*ftproto.Feature
		targetFeatureID   string
		deletedVariations map[string]string // variationID -> variationValue
		expected          error
	}{
		{
			desc:            "success: no features using deleted variations",
			features:        []*ftproto.Feature{},
			targetFeatureID: "feature-1",
			deletedVariations: map[string]string{
				variationID1: variationValue1,
			},
			expected: nil,
		},
		{
			desc: "success: target feature uses variation (should be excluded)",
			features: []*ftproto.Feature{
				{
					Id: "feature-1",
					Prerequisites: []*ftproto.Prerequisite{
						{
							FeatureId:   "feature-1",
							VariationId: variationID1,
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			deletedVariations: map[string]string{
				variationID1: variationValue1,
			},
			expected: nil,
		},
		{
			desc: "error: other feature has prerequisite using deleted variation",
			features: []*ftproto.Feature{
				{
					Id: "feature-2",
					Prerequisites: []*ftproto.Prerequisite{
						{
							FeatureId:   "feature-1",
							VariationId: variationID1,
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			deletedVariations: map[string]string{
				variationID1: variationValue1,
			},
			expected: ErrVariationInUse,
		},
		{
			desc: "error: other feature has FEATURE_FLAG rule using deleted variation ID",
			features: []*ftproto.Feature{
				{
					Id: "feature-2",
					Rules: []*ftproto.Rule{
						{
							Clauses: []*ftproto.Clause{
								{
									Operator:  ftproto.Clause_FEATURE_FLAG,
									Attribute: "feature-1",
									Values:    []string{variationID1}, // Fixed: Use variation ID, not value
								},
							},
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			deletedVariations: map[string]string{
				variationID1: variationValue1,
			},
			expected: ErrVariationInUse,
		},
		{
			desc: "success: no variations to delete",
			features: []*ftproto.Feature{
				{
					Id: "feature-2",
					Prerequisites: []*ftproto.Prerequisite{
						{
							FeatureId:   "feature-1",
							VariationId: variationID1,
						},
					},
				},
			},
			targetFeatureID:   "feature-1",
			deletedVariations: map[string]string{},
			expected:          nil,
		},
		{
			desc: "success: different feature ID in prerequisite",
			features: []*ftproto.Feature{
				{
					Id: "feature-2",
					Prerequisites: []*ftproto.Prerequisite{
						{
							FeatureId:   "feature-3", // Different feature
							VariationId: variationID1,
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			deletedVariations: map[string]string{
				variationID1: variationValue1,
			},
			expected: nil,
		},
		{
			desc: "success: different variation ID in FEATURE_FLAG rule",
			features: []*ftproto.Feature{
				{
					Id: "feature-2",
					Rules: []*ftproto.Rule{
						{
							Clauses: []*ftproto.Clause{
								{
									Operator:  ftproto.Clause_FEATURE_FLAG,
									Attribute: "feature-1",
									Values:    []string{variationID2}, // Fixed: Different variation ID
								},
							},
						},
					},
				},
			},
			targetFeatureID: "feature-1",
			deletedVariations: map[string]string{
				variationID1: variationValue1,
			},
			expected: nil,
		},
		{
			desc: "error: detailed FEATURE_FLAG rule test",
			features: []*ftproto.Feature{
				{
					Id: "feature-A",
					Variations: []*ftproto.Variation{
						{Id: "var-true", Value: "true"},
						{Id: "var-false", Value: "false"},
					},
				},
				{
					Id: "feature-B",
					Rules: []*ftproto.Rule{
						{
							Id: "test-rule",
							Clauses: []*ftproto.Clause{
								{
									Id:        "test-clause",
									Operator:  ftproto.Clause_FEATURE_FLAG,
									Attribute: "feature-A",          // References feature being updated
									Values:    []string{"var-true"}, // Fixed: Use variation ID, not value
								},
							},
						},
					},
				},
			},
			targetFeatureID: "feature-A",
			deletedVariations: map[string]string{
				"var-true": "true", // Deleting variation with value "true"
			},
			expected: ErrVariationInUse,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			err := ValidateVariationUsage(p.features, p.targetFeatureID, p.deletedVariations)
			if p.expected != nil {
				assert.Error(t, err)
				assert.Equal(t, p.expected, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateClauses(t *testing.T) {
	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	id3, _ := uuid.NewUUID()
	patterns := []struct {
		desc     string
		clauses  []*ftproto.Clause
		expected error
	}{
		{
			desc: "success: 2 clauses",
			clauses: []*ftproto.Clause{
				{
					Id:        id1.String(),
					Attribute: "name",
					Operator:  ftproto.Clause_EQUALS,
					Values: []string{
						"user1",
						"user2",
					},
				},
				{
					Id:       id2.String(),
					Operator: ftproto.Clause_SEGMENT,
					Values: []string{
						"value",
					},
				},
				{
					Id:        id3.String(),
					Operator:  ftproto.Clause_FEATURE_FLAG,
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
			clauses:  []*ftproto.Clause{},
			expected: fmt.Errorf("feature: rule must have at least one clause"),
		},
		{
			desc: "err: id is empty",
			clauses: []*ftproto.Clause{
				{
					Id:        "",
					Attribute: "name",
					Operator:  ftproto.Clause_EQUALS,
					Values:    []string{"user1"},
				},
			},
			expected: errors.New("feature: clause id cannot be empty"),
		},
		{
			desc: "err: compare missing attribute",
			clauses: []*ftproto.Clause{
				{
					Id:       id1.String(),
					Operator: ftproto.Clause_EQUALS,
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
			clauses: []*ftproto.Clause{
				{
					Id:        id1.String(),
					Operator:  ftproto.Clause_EQUALS,
					Attribute: "name",
				},
			},
			expected: errClauseValuesEmpty,
		},
		{
			desc: "err: segment attribute not empty",
			clauses: []*ftproto.Clause{
				{
					Id:        id1.String(),
					Operator:  ftproto.Clause_SEGMENT,
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
			clauses: []*ftproto.Clause{
				{
					Id:       id1.String(),
					Operator: ftproto.Clause_SEGMENT,
				},
			},
			expected: errClauseValuesEmpty,
		},
		{
			desc: "err: feature flag attribute empty",
			clauses: []*ftproto.Clause{
				{
					Id:       id1.String(),
					Operator: ftproto.Clause_FEATURE_FLAG,
				},
			},
			expected: errClauseAttributeEmpty,
		},
		{
			desc: "err: feature flag values empty",
			clauses: []*ftproto.Clause{
				{
					Id:        id1.String(),
					Operator:  ftproto.Clause_FEATURE_FLAG,
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
	id3, _ := uuid.NewUUID()
	ruleID, _ := uuid.NewUUID() // For rule ID
	patterns := []struct {
		desc                string
		feature             *Feature
		name                *wrapperspb.StringValue
		description         *wrapperspb.StringValue
		tags                *common.StringListValue
		enabled             *wrapperspb.BoolValue
		archived            *wrapperspb.BoolValue
		defaultStrategy     *ftproto.Strategy
		offVariation        *wrapperspb.StringValue
		resetSamplingSeed   bool
		prerequisiteChanges []*ftproto.PrerequisiteChange
		targetChanges       []*ftproto.TargetChange
		ruleChanges         []*ftproto.RuleChange
		variationChanges    []*ftproto.VariationChange
		tagChanges          []*ftproto.TagChange
		expected            *Feature
		expectedErr         error
	}{
		{
			desc: "success: no changes when updating with same values",
			feature: &Feature{Feature: &ftproto.Feature{
				Name:        "test-feature",
				Description: "test description",
				Tags:        []string{"tag1"},
				Enabled:     true,
				Archived:    false,
				DefaultStrategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id2.String(),
				Variations: []*ftproto.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
			}},
			name:        wrapperspb.String("test-feature"),
			description: wrapperspb.String("test description"),
			tags:        &common.StringListValue{Values: []string{"tag1"}},
			enabled:     wrapperspb.Bool(true),
			archived:    wrapperspb.Bool(false),
			defaultStrategy: &ftproto.Strategy{
				Type: ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{
					Variation: id1.String(),
				},
			},
			offVariation: wrapperspb.String(id2.String()),
			expected: &Feature{Feature: &ftproto.Feature{
				Name:        "test-feature",
				Description: "test description",
				Tags:        []string{"tag1"},
				Enabled:     true,
				Archived:    false,
				DefaultStrategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id2.String(),
				Variations: []*ftproto.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets:       []*ftproto.Target{},
				Rules:         []*ftproto.Rule{},
			}},
		},
		{
			desc: "success: version incremented with full replacement fields",
			feature: &Feature{Feature: &ftproto.Feature{
				Name:        "old-name",
				Description: "old description",
				Tags:        []string{"old-tag"},
				Enabled:     false,
				Archived:    false,
				DefaultStrategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id2.String(),
				Variations: []*ftproto.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets:       []*ftproto.Target{},
				Rules:         []*ftproto.Rule{},
			}},
			name:        wrapperspb.String("new-name"),
			description: wrapperspb.String("new description"),
			tags:        &common.StringListValue{Values: []string{"new-tag"}},
			enabled:     wrapperspb.Bool(true),
			archived:    wrapperspb.Bool(false),
			defaultStrategy: &ftproto.Strategy{
				Type: ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{
					Variation: id1.String(),
				},
			},
			offVariation: wrapperspb.String(id1.String()),
			expected: &Feature{Feature: &ftproto.Feature{
				Name:        "new-name",
				Description: "new description",
				Tags:        []string{"new-tag"},
				Enabled:     true,
				Archived:    false,
				DefaultStrategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: id1.String(),
					},
				},
				OffVariation: id1.String(),
				Variations: []*ftproto.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets:       []*ftproto.Target{},
				Rules:         []*ftproto.Rule{},
			}},
		},
		{
			desc: "success: granular updates for variations, rules, prerequisites, targets, and tags",
			feature: &Feature{Feature: &ftproto.Feature{
				Name: "test-feature",
				Variations: []*ftproto.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
				},
				Targets: []*ftproto.Target{
					{Variation: id1.String(), Users: []string{}},
					{Variation: id2.String(), Users: []string{}},
				},
			}},
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Variation: &ftproto.Variation{
						Id:    id3.String(),
						Name:  "v3",
						Value: "new-value",
					},
				},
			},
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Rule: &ftproto.Rule{
						Id: ruleID.String(),
						Strategy: &ftproto.Strategy{
							Type: ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{
								Variation: id1.String(),
							},
						},
						Clauses: []*ftproto.Clause{
							{
								Id:        id1.String(),
								Operator:  ftproto.Clause_EQUALS,
								Values:    []string{"user1"},
								Attribute: "name",
							},
						},
					},
				},
			},
			prerequisiteChanges: []*ftproto.PrerequisiteChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Prerequisite: &ftproto.Prerequisite{
						FeatureId:   "feature-1",
						VariationId: id1.String(),
					},
				},
			},
			targetChanges: []*ftproto.TargetChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Target: &ftproto.Target{
						Variation: id1.String(),
						Users:     []string{"user1"},
					},
				},
			},
			tagChanges: []*ftproto.TagChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Tag:        "new-tag",
				},
			},
			expected: &Feature{Feature: &ftproto.Feature{
				Name: "test-feature",
				Tags: []string{"new-tag"},
				Variations: []*ftproto.Variation{
					{Id: id1.String(), Name: "v1", Value: "true"},
					{Id: id2.String(), Name: "v2", Value: "false"},
					{Name: "v3", Value: "new-value"},
				},
				Rules: []*ftproto.Rule{
					{
						Id: ruleID.String(),
						Strategy: &ftproto.Strategy{
							Type: ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{
								Variation: id1.String(),
							},
						},
						Clauses: []*ftproto.Clause{
							{
								Id:        id1.String(),
								Operator:  ftproto.Clause_EQUALS,
								Values:    []string{"user1"},
								Attribute: "name",
							},
						},
					},
				},
				Prerequisites: []*ftproto.Prerequisite{
					{
						FeatureId:   "feature-1",
						VariationId: id1.String(),
					},
				},
				Targets: []*ftproto.Target{
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
			hasVariationCreate := len(p.variationChanges) > 0 && p.variationChanges[0].ChangeType == ftproto.ChangeType_CREATE
			if hasVariationCreate {
				assert.Equal(t, len(p.expected.Variations), len(updated.Variations))
				for i := range updated.Variations {
					assert.Equal(t, p.expected.Variations[i].Name, updated.Variations[i].Name)
					assert.Equal(t, p.expected.Variations[i].Value, updated.Variations[i].Value)
					assert.Equal(t, p.expected.Variations[i].Description, updated.Variations[i].Description)
				}
				createdVariationID = updated.Variations[len(updated.Variations)-1].Id
			} else {
				// Use proto.Equal for protobuf message comparison to avoid internal structure differences
				assert.Equal(t, len(p.expected.Variations), len(updated.Variations), "Variations length should match")
				for i := range p.expected.Variations {
					assert.True(t, proto.Equal(p.expected.Variations[i], updated.Variations[i]), "Variation %d should be equal", i)
				}
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
				// Use proto.Equal for protobuf message comparison
				assert.Equal(t, len(p.expected.Targets), len(updated.Targets), "Targets length should match")
				for i := range p.expected.Targets {
					assert.True(t, proto.Equal(p.expected.Targets[i], updated.Targets[i]), "Target %d should be equal", i)
				}
			}
			// Check prerequisites and rules using proto.Equal
			assert.Equal(t, len(p.expected.Prerequisites), len(updated.Prerequisites), "Prerequisites length should match")
			for i := range p.expected.Prerequisites {
				assert.True(t, proto.Equal(p.expected.Prerequisites[i], updated.Prerequisites[i]), "Prerequisite %d should be equal", i)
			}
			assert.Equal(t, len(p.expected.Rules), len(updated.Rules), "Rules length should match")
			for i := range p.expected.Rules {
				assert.True(t, proto.Equal(p.expected.Rules[i], updated.Rules[i]), "Rule %d should be equal", i)
			}
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
			Feature: &ftproto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1", "t2"},
				VariationType: ftproto.Feature_BOOLEAN,
				Variations: []*ftproto.Variation{
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets:       []*ftproto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*ftproto.Rule{},
				DefaultStrategy: &ftproto.Strategy{
					Type:          ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	patterns := []struct {
		desc                string
		inputFunc           func() *Feature
		prerequisiteChanges []*ftproto.PrerequisiteChange
		expectedFunc        func() *Feature
		expectedErr         error
	}{
		{
			desc:      "Prerequisite Create - success",
			inputFunc: genF,
			prerequisiteChanges: []*ftproto.PrerequisiteChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Prerequisite: &ftproto.Prerequisite{
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
			prerequisiteChanges: []*ftproto.PrerequisiteChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Prerequisite: &ftproto.Prerequisite{
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
			prerequisiteChanges: []*ftproto.PrerequisiteChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Prerequisite: &ftproto.Prerequisite{
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
			prerequisiteChanges: []*ftproto.PrerequisiteChange{
				{
					ChangeType: ftproto.ChangeType_DELETE,
					Prerequisite: &ftproto.Prerequisite{
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
			Feature: &ftproto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1"},
				VariationType: ftproto.Feature_BOOLEAN,
				Variations: []*ftproto.Variation{
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets: []*ftproto.Target{
					{Variation: v1.String(), Users: []string{"u1"}},
					{Variation: v2.String(), Users: []string{"u2"}},
				},
				Rules: []*ftproto.Rule{},
				DefaultStrategy: &ftproto.Strategy{
					Type:          ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	patterns := []struct {
		desc          string
		inputFunc     func() *Feature
		targetChanges []*ftproto.TargetChange
		expectedFunc  func() *Feature
		expectedErr   error
	}{
		{
			desc:      "Target Create - error: empty target fields",
			inputFunc: genF,
			targetChanges: []*ftproto.TargetChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Target:     &ftproto.Target{Variation: "", Users: []string{}},
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
			targetChanges: []*ftproto.TargetChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Target:     &ftproto.Target{Variation: "non-existent", Users: []string{"u-new"}},
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
			targetChanges: []*ftproto.TargetChange{
				{
					ChangeType: ftproto.ChangeType_DELETE,
					Target:     &ftproto.Target{Variation: "non-existent"},
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
			Feature: &ftproto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1"},
				VariationType: ftproto.Feature_BOOLEAN,
				Variations: []*ftproto.Variation{
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets:       []*ftproto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*ftproto.Rule{},
				DefaultStrategy: &ftproto.Strategy{
					Type:          ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	// Define patterns for rule granular updates including clause validations.
	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		ruleChanges  []*ftproto.RuleChange
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc:      "Rule Create - error: rule required",
			inputFunc: genF,
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Rule:       nil,
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errRuleRequired,
		},
		{
			desc:      "Rule Create - error: nil strategy",
			inputFunc: genF,
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Rule: &ftproto.Rule{
						Id:       ruleID.String(),
						Strategy: nil, // This should trigger errStrategyRequired
						Clauses: []*ftproto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "attr",
								Operator:  ftproto.Clause_EQUALS,
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
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Rule: &ftproto.Rule{
						Id: ruleID.String(),
						Strategy: &ftproto.Strategy{
							Type:          ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*ftproto.Clause{},
					},
				},
			},
			expectedFunc: func() *Feature { return genF() },
			expectedErr:  errors.New("feature: rule must have at least one clause"),
		},
		{
			desc:      "Rule Create - error: clause attribute not empty for SEGMENT operator",
			inputFunc: genF,
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Rule: &ftproto.Rule{
						Id: ruleID.String(),
						Strategy: &ftproto.Strategy{
							Type:          ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*ftproto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "non-empty", // Not allowed for SEGMENT operator
								Operator:  ftproto.Clause_SEGMENT,
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
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Rule: &ftproto.Rule{
						Id: ruleID.String(),
						Strategy: &ftproto.Strategy{
							Type:          ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*ftproto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "", // Correct for SEGMENT operator
								Operator:  ftproto.Clause_SEGMENT,
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
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Rule: &ftproto.Rule{
						Id: "non-existent",
						Strategy: &ftproto.Strategy{
							Type:          ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{Variation: v2.String()},
						},
						Clauses: []*ftproto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "attr",
								Operator:  ftproto.Clause_EQUALS,
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
				_ = f.AddRule(&ftproto.Rule{
					Id: ruleID.String(),
					Strategy: &ftproto.Strategy{
						Type:          ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  ftproto.Clause_EQUALS,
							Values:    []string{"val"},
						},
					},
				})
				return f
			},
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Rule: &ftproto.Rule{
						Id: ruleID.String(),
						Strategy: &ftproto.Strategy{
							Type:          ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*ftproto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "", // empty attribute for non-SEGMENT operator not allowed
								Operator:  ftproto.Clause_EQUALS,
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
				_ = f.AddRule(&ftproto.Rule{
					Id: ruleID.String(),
					Strategy: &ftproto.Strategy{
						Type:          ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
					},
					Clauses: []*ftproto.Clause{
						{
							Id:        clauseID.String(),
							Attribute: "attr",
							Operator:  ftproto.Clause_EQUALS,
							Values:    []string{"val"},
						},
					},
				})
				return f
			},
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Rule: &ftproto.Rule{
						Id: ruleID.String(),
						Strategy: &ftproto.Strategy{
							Type:          ftproto.Strategy_FIXED,
							FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
						},
						Clauses: []*ftproto.Clause{
							{
								Id:        clauseID.String(),
								Attribute: "attr",
								Operator:  ftproto.Clause_EQUALS,
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
			ruleChanges: []*ftproto.RuleChange{
				{
					ChangeType: ftproto.ChangeType_DELETE,
					Rule:       &ftproto.Rule{Id: ""},
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
			Feature: &ftproto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1"},
				VariationType: ftproto.Feature_JSON,
				Variations: []*ftproto.Variation{
					{Id: v1.String(), Value: `{"key": "value1"}`, Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: `{"key": "value2"}`, Name: "n2", Description: "d2"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets:       []*ftproto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*ftproto.Rule{},
				DefaultStrategy: &ftproto.Strategy{
					Type:          ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	// Define test patterns.
	patterns := []struct {
		desc             string
		inputFunc        func() *Feature
		variationChanges []*ftproto.VariationChange
		expectedFunc     func() *Feature
		expectedErr      error
	}{
		{
			desc:      "Variation Create - success",
			inputFunc: genFJSON,
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Variation: &ftproto.Variation{
						Id:          v3.String(),
						Value:       `{"key": "value3"}`,
						Name:        "n3",
						Description: "d3",
					},
				},
			},
			expectedFunc: func() *Feature {
				f := genFJSON()
				// Add the new variation directly
				f.Variations = append(f.Variations, &ftproto.Variation{
					Id:          v3.String(),
					Value:       `{"key": "value3"}`,
					Name:        "n3",
					Description: "d3",
				})
				// Add corresponding target for the new variation
				f.Targets = append(f.Targets, &ftproto.Target{
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
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Variation: &ftproto.Variation{
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
				// Add the variation to be deleted and check for errors
				err := f.AddVariation(v3.String(), `{"key": "value3"}`, "n3", "d3")
				require.NoError(t, err)
				return f
			},
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_DELETE,
					Variation:  &ftproto.Variation{Id: v3.String()},
				},
			},
			expectedFunc: genFJSON,
			expectedErr:  nil,
		},
		{
			desc:      "Variation Update - error: nil variation",
			inputFunc: genFJSON,
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Variation:  nil,
				},
			},
			expectedFunc: func() *Feature { return genFJSON() },
			expectedErr:  errVariationRequired,
		},
		{
			desc:      "Variation Update - error: empty name",
			inputFunc: genFJSON,
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Variation: &ftproto.Variation{
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
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Variation: &ftproto.Variation{
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
			variationChanges: []*ftproto.VariationChange{
				{
					ChangeType: ftproto.ChangeType_UPDATE,
					Variation: &ftproto.Variation{
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
			t.Parallel()
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
			Feature: &ftproto.Feature{
				Id:            "i",
				Name:          "n",
				Description:   "d",
				Archived:      false,
				Enabled:       false,
				Tags:          []string{"t1", "t2"},
				VariationType: ftproto.Feature_BOOLEAN,
				Variations: []*ftproto.Variation{
					{Id: v1.String(), Value: "true", Name: "n1", Description: "d1"},
					{Id: v2.String(), Value: "false", Name: "n2", Description: "d2"},
				},
				Prerequisites: []*ftproto.Prerequisite{},
				Targets:       []*ftproto.Target{{Variation: v1.String()}, {Variation: v2.String()}},
				Rules:         []*ftproto.Rule{},
				DefaultStrategy: &ftproto.Strategy{
					Type:          ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{Variation: v1.String()},
				},
				OffVariation: v1.String(),
			},
		}
	}

	patterns := []struct {
		desc         string
		inputFunc    func() *Feature
		tagChanges   []*ftproto.TagChange
		expectedFunc func() *Feature
		expectedErr  error
	}{
		{
			desc:      "Tag Create - error: duplicate create should not add duplicate",
			inputFunc: genF,
			tagChanges: []*ftproto.TagChange{
				{
					ChangeType: ftproto.ChangeType_CREATE,
					Tag:        "new-tag",
				},
				{
					ChangeType: ftproto.ChangeType_CREATE,
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
			tagChanges: []*ftproto.TagChange{
				{
					ChangeType: ftproto.ChangeType_DELETE,
					Tag:        "non-existent-tag",
				},
			},
			expectedFunc: func() *Feature {
				return genF()
			},
			expectedErr: errors.New("feature: tag not found"),
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
	variations := []*ftproto.Variation{
		{Id: id1.String(), Value: "true", Name: "n1", Description: "d1"},
		{Id: id2.String(), Value: "false", Name: "n2", Description: "d2"},
	}
	tests := []struct {
		desc        string
		strategy    *ftproto.Strategy
		variations  []*ftproto.Variation
		expectedErr error
	}{
		{
			desc: "success: fixed strategy",
			strategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: id1.String()},
			},
			variations:  variations,
			expectedErr: nil,
		},
		{
			desc: "success: rollout strategy",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100000},
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
			strategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: "non-existent"},
			},
			variations:  variations,
			expectedErr: errVariationNotFound,
		},
		{
			desc: "fail: rollout strategy with non-existent variation",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: "non-existent", Weight: 100},
					},
				},
			},
			variations:  variations,
			expectedErr: errVariationNotFound,
		},
		{
			desc: "fail: unsupported strategy type",
			strategy: &ftproto.Strategy{
				Type: 999,
			},
			variations:  variations,
			expectedErr: errUnsupportedStrategy,
		},
		{
			desc: "fail: fixed strategy is nil",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_FIXED,
			},
			variations:  variations,
			expectedErr: ErrRuleStrategyCannotBeEmpty,
		},
		{
			desc: "fail: rollout strategy is nil",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
			},
			variations:  variations,
			expectedErr: ErrRuleStrategyCannotBeEmpty,
		},
		{
			desc: "fail: both strategies are set",
			strategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: id1.String()},
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100},
					},
				},
			},
			variations:  variations,
			expectedErr: ErrDefaultStrategyCannotBeBothFixedAndRollout,
		},
		{
			desc: "success: rollout strategy with valid audience",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 50000},
						{Variation: id2.String(), Weight: 50000},
					},
					Audience: &ftproto.Audience{
						Percentage:       50,
						DefaultVariation: id1.String(),
					},
				},
			},
			variations:  variations,
			expectedErr: nil,
		},
		{
			desc: "success: rollout strategy with 0% audience",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100000},
					},
					Audience: &ftproto.Audience{
						Percentage:       0,
						DefaultVariation: "",
					},
				},
			},
			variations:  variations,
			expectedErr: nil,
		},
		{
			desc: "success: rollout strategy with 100% audience",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100000},
					},
					Audience: &ftproto.Audience{
						Percentage:       100,
						DefaultVariation: "",
					},
				},
			},
			variations:  variations,
			expectedErr: nil,
		},
		{
			desc: "fail: audience percentage below 0",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100000},
					},
					Audience: &ftproto.Audience{
						Percentage:       -1,
						DefaultVariation: id1.String(),
					},
				},
			},
			variations:  variations,
			expectedErr: ErrInvalidAudiencePercentage,
		},
		{
			desc: "fail: audience percentage above 100",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100000},
					},
					Audience: &ftproto.Audience{
						Percentage:       101,
						DefaultVariation: id1.String(),
					},
				},
			},
			variations:  variations,
			expectedErr: ErrInvalidAudiencePercentage,
		},
		{
			desc: "fail: audience percentage between 1-99 without default variation",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100000},
					},
					Audience: &ftproto.Audience{
						Percentage:       50,
						DefaultVariation: "",
					},
				},
			},
			variations:  variations,
			expectedErr: ErrInvalidAudienceDefaultVariation,
		},
		{
			desc: "fail: audience with non-existent default variation",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 100000},
					},
					Audience: &ftproto.Audience{
						Percentage:       50,
						DefaultVariation: "non-existent",
					},
				},
			},
			variations:  variations,
			expectedErr: ErrDefaultVariationNotFound,
		},
		{
			desc: "fail: rollout strategy weights sum less than 100000",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 30000}, // 30%
						{Variation: id2.String(), Weight: 40000}, // 40%
						// Total: 70000 (70%) - invalid!
					},
				},
			},
			variations:  variations,
			expectedErr: ErrInvalidVariationWeightTotal,
		},
		{
			desc: "fail: rollout strategy weights sum more than 100000",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 60000}, // 60%
						{Variation: id2.String(), Weight: 50000}, // 50%
						// Total: 110000 (110%) - invalid!
					},
				},
			},
			variations:  variations,
			expectedErr: ErrInvalidVariationWeightTotal,
		},
		{
			desc: "fail: rollout strategy weights zero total",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 0},
						{Variation: id2.String(), Weight: 0},
						// Total: 0 (0%) - invalid!
					},
				},
			},
			variations:  variations,
			expectedErr: ErrInvalidVariationWeightTotal,
		},
		{
			desc: "fail: rollout strategy weights using old test values",
			strategy: &ftproto.Strategy{
				Type: ftproto.Strategy_ROLLOUT,
				RolloutStrategy: &ftproto.RolloutStrategy{
					Variations: []*ftproto.RolloutStrategy_Variation{
						{Variation: id1.String(), Weight: 30}, // Old test format
						{Variation: id2.String(), Weight: 70}, // Old test format
						// Total: 100 instead of 100000 - invalid!
					},
				},
			},
			variations:  variations,
			expectedErr: ErrInvalidVariationWeightTotal,
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

func TestRemoveVariationMinimumVariationConstraint(t *testing.T) {
	// Test case 1: Feature with exactly 3 variations - should allow removal (leaves 2)
	f := makeFeature("test-feature")
	// makeFeature creates 3 variations: A, B, C
	// Remove variation-C (which has users, so we need to remove them first)
	f.Targets[2].Users = []string{} // Remove users from variation-C

	patterns := []*struct {
		id       string
		expected error
	}{
		{
			id:       "variation-C",
			expected: nil, // Should succeed - leaves 2 variations
		},
	}

	for i, p := range patterns {
		err := f.RemoveVariation(p.id)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}

	// Verify we now have 2 variations
	if len(f.Variations) != 2 {
		t.Fatalf("Expected 2 variations after removal, got %d", len(f.Variations))
	}

	// Test case 2: Now try to remove another variation - should fail (would leave 1)
	patterns2 := []*struct {
		id       string
		expected error
	}{
		{
			id:       "variation-A",
			expected: errVariationsMustHaveAtLeastTwoVariations, // Should fail - would leave 1 variation
		},
		{
			id:       "variation-B",
			expected: errVariationsMustHaveAtLeastTwoVariations, // Should fail - would leave 1 variation
		},
	}

	for i, p := range patterns2 {
		err := f.RemoveVariation(p.id)
		des := fmt.Sprintf("constraint_test_index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}

	// Verify we still have 2 variations (removal should have failed)
	if len(f.Variations) != 2 {
		t.Fatalf("Expected 2 variations after failed removal attempts, got %d", len(f.Variations))
	}
}
