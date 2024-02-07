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
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/proto/feature"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
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

/*
func TestUserAssignment(t *testing.T) {
	f := feature("test-feature")

	fmt.Println(f.assignUser("user1"))
	fmt.Println(f.assignUser("user2"))
	fmt.Println(f.assignUser("user3"))

	user1hash := f.Hash("user1")
	user2hash := f.Hash("user2")
	fmt.Println(hex.EncodeToString(user1hash[:8]))
	fmt.Println(hex.EncodeToString(user2hash[:8]))

	fmt.Println(f.Bucket("user1"))
	fmt.Println(f.Bucket("user2"))
func TestProportions(t *testing.T) {
	f := feature("test-feature")
	bucketA := 0
	bucketB := 0
	for i := 10000; i < 20000; i++ {
		user := fmt.Sprintf("user-%d", i)
		bucket := f.Bucket(user)
		if bucket < 0.2 {
			bucketA++
		} else {
			bucketB++
		}
	}
	a := float64(bucketA) / 10000.0 // should be close to 0.2
	b := float64(bucketB) / 10000.0 // should be close to 0.8
	assert.InEpsilon(t, 0.2, a, 0.05)
	assert.InEpsilon(t, 0.8, b, 0.05)
}

func TestCorrelation(t *testing.T) {
	// create some hundred tests
	// assign all people in each test
	// compute mutual information between each test pair
	// should be high with itself and else low
}
*/

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

func TestAssignUserOffVariation(t *testing.T) {
	t.Parallel()
	f := makeFeature("test-feature")
	patterns := []struct {
		enabled           bool
		offVariation      string
		userID            string
		Flagvariations    map[string]string
		prerequisite      []*proto.Prerequisite
		expectedReason    *proto.Reason
		expectedVariation *proto.Variation
		expectedError     error
	}{
		{
			enabled:           false,
			offVariation:      "variation-C",
			userID:            "user5",
			Flagvariations:    map[string]string{},
			prerequisite:      []*proto.Prerequisite{},
			expectedReason:    &proto.Reason{Type: proto.Reason_OFF_VARIATION},
			expectedVariation: f.Variations[2],
			expectedError:     nil,
		},
		{
			enabled:           false,
			offVariation:      "",
			userID:            "user5",
			Flagvariations:    map[string]string{},
			prerequisite:      []*proto.Prerequisite{},
			expectedReason:    &proto.Reason{Type: proto.Reason_DEFAULT},
			expectedVariation: f.Variations[1],
			expectedError:     nil,
		},
		{
			enabled:           false,
			offVariation:      "variation-E",
			userID:            "user5",
			Flagvariations:    map[string]string{},
			prerequisite:      []*proto.Prerequisite{},
			expectedReason:    &proto.Reason{Type: proto.Reason_OFF_VARIATION},
			expectedVariation: nil,
			expectedError:     errVariationNotFound,
		},
		{
			enabled:           true,
			offVariation:      "",
			userID:            "user4",
			Flagvariations:    map[string]string{},
			prerequisite:      []*proto.Prerequisite{},
			expectedReason:    &proto.Reason{Type: proto.Reason_DEFAULT},
			expectedVariation: f.Variations[1],
			expectedError:     nil,
		},
		{
			enabled:           true,
			offVariation:      "variation-C",
			userID:            "user4",
			Flagvariations:    map[string]string{},
			prerequisite:      []*proto.Prerequisite{},
			expectedReason:    &proto.Reason{Type: proto.Reason_DEFAULT},
			expectedVariation: f.Variations[1],
			expectedError:     nil,
		},
		{
			enabled:      true,
			offVariation: "variation-C",
			userID:       "user4",
			Flagvariations: map[string]string{
				"test-feature2": "variation A", // not matched with expected prerequisites variations
			},
			prerequisite: []*proto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expectedReason:    &proto.Reason{Type: proto.Reason_PREREQUISITE},
			expectedVariation: f.Variations[2],
			expectedError:     nil,
		},
		{
			enabled:      true,
			offVariation: "variation-C",
			userID:       "user4",
			Flagvariations: map[string]string{
				"test-feature2": "variation D", // matched with expected prerequisites variations
			},
			prerequisite: []*proto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expectedReason:    &proto.Reason{Type: proto.Reason_DEFAULT},
			expectedVariation: f.Variations[1],
			expectedError:     nil,
		},
		{
			enabled:        true,
			offVariation:   "variation-C",
			userID:         "user4",
			Flagvariations: map[string]string{}, // not found prerequisite vatiation
			prerequisite: []*proto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expectedReason:    nil,
			expectedVariation: nil,
			expectedError:     errPrerequisiteVariationNotFound,
		},
	}
	for _, p := range patterns {
		user := &userproto.User{Id: p.userID}
		f.Enabled = p.enabled
		f.OffVariation = p.offVariation
		f.Prerequisites = p.prerequisite
		reason, variation, err := f.assignUser(user, nil, p.Flagvariations)
		assert.Equal(t, p.expectedReason, reason)
		assert.Equal(t, p.expectedVariation, variation)
		assert.Equal(t, p.expectedError, err)
	}
}

func TestAssignUserTarget(t *testing.T) {
	f := makeFeature("test-feature")
	patterns := []struct {
		userID              string
		expectedReason      proto.Reason_Type
		expectedVariationID string
	}{
		{
			userID:              "user1",
			expectedReason:      proto.Reason_TARGET,
			expectedVariationID: "variation-A",
		},
		{
			userID:              "user2",
			expectedReason:      proto.Reason_TARGET,
			expectedVariationID: "variation-B",
		},
		{
			userID:              "user3",
			expectedReason:      proto.Reason_TARGET,
			expectedVariationID: "variation-C",
		},
		{
			userID:              "user4",
			expectedReason:      proto.Reason_DEFAULT,
			expectedVariationID: "variation-B",
		},
	}
	for _, p := range patterns {
		user := &userproto.User{Id: p.userID}
		reason, variation, err := f.assignUser(user, nil, nil)
		assert.Equal(t, p.expectedReason, reason.Type)
		assert.Equal(t, p.expectedVariationID, variation.Id)
		assert.NoError(t, err)
	}
}

func TestAssignUserRuleSet(t *testing.T) {
	user := &userproto.User{
		Id:   "user-id",
		Data: map[string]string{"name": "user3"},
	}
	f := makeFeature("test-feature")
	reason, variation, err := f.assignUser(user, nil, nil)
	if err != nil {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
	if reason.RuleId != "rule-2" {
		t.Fatalf("Failed to assign user. Reason id does not match. ID: %s", reason.RuleId)
	}
	if variation.Id != "variation-B" {
		t.Fatalf("Failed to assign user. Variation id does not match. ID: %s", variation.Id)
	}
}

func TestAssignUserWithNoDefaultStrategy(t *testing.T) {
	user := &userproto.User{
		Id:   "user-id1",
		Data: map[string]string{"name3": "user3"},
	}
	f := makeFeature("test-feature")
	f.DefaultStrategy = nil

	reason, variation, err := f.assignUser(user, nil, nil)
	if reason != nil {
		t.Fatalf("Failed to assign user. Reason should be nil: %v", reason)
	}
	if variation != nil {
		t.Fatalf("Failed to assign user. Variation should be nil: %v", variation)
	}
	if err != errDefaultStrategyNotFound {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
}

func TestAssignUserDefaultStrategy(t *testing.T) {
	user := &userproto.User{
		Id:   "user-id1",
		Data: map[string]string{"name3": "user3"},
	}
	f := makeFeature("test-feature")
	reason, variation, err := f.assignUser(user, nil, nil)
	if err != nil {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
	if reason.Type != proto.Reason_DEFAULT {
		t.Fatalf("Failed to assign user. Reason type does not match. Current: %s, target: %v", reason.Type, proto.Reason_DEFAULT)
	}
	targetVariationID := "variation-B"
	if variation.Id != targetVariationID {
		t.Fatalf("Failed to assign user. Variation id does not match. Current: %s, target: %s", variation.Id, targetVariationID)
	}
}

func TestAssignUserSamplingSeed(t *testing.T) {
	user := &userproto.User{
		Id:   "uid",
		Data: map[string]string{},
	}
	f := makeFeature("fid")
	f.DefaultStrategy = &proto.Strategy{
		Type: proto.Strategy_ROLLOUT,
		RolloutStrategy: &proto.RolloutStrategy{
			Variations: []*proto.RolloutStrategy_Variation{
				{
					Variation: f.Variations[0].Id,
					Weight:    30000,
				},
				{
					Variation: f.Variations[1].Id,
					Weight:    40000,
				},
				{
					Variation: f.Variations[2].Id,
					Weight:    30000,
				},
			},
		},
	}
	reason, variation, err := f.assignUser(user, nil, nil)
	if err != nil {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
	if reason.Type != proto.Reason_DEFAULT {
		t.Fatalf("Failed to assign user. Reason type does not match. Current: %s, target: %v", reason.Type, proto.Reason_DEFAULT)
	}
	if variation.Id != f.DefaultStrategy.RolloutStrategy.Variations[1].Variation {
		t.Fatalf("Failed to assign user. Variation id does not match. Current: %s, target: %s", variation.Id, f.DefaultStrategy.RolloutStrategy.Variations[1].Variation)
	}
	// Channge sampling seed to change assigned variation.
	f.SamplingSeed = "test"
	reason, variation, err = f.assignUser(user, nil, nil)
	if err != nil {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
	if reason.Type != proto.Reason_DEFAULT {
		t.Fatalf("Failed to assign user. Reason type does not match. Current: %s, target: %v", reason.Type, proto.Reason_DEFAULT)
	}
	if variation.Id != f.DefaultStrategy.RolloutStrategy.Variations[0].Variation {
		t.Fatalf("Failed to assign user. Variation id does not match. Current: %s, target: %s", variation.Id, f.DefaultStrategy.RolloutStrategy.Variations[0].Variation)
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
	if f.Tags[0] == tag1 {
		t.Fatalf("Failed to remove tag %s. Tags: %v", tag1, f.Tags)
	}
	if len(f.Tags) != 1 {
		t.Fatalf("Failed to remove tag. It should remove only 1: %v", f.Tags)
	}
	if err := f.RemoveTag(tag2); err != errTagsMustHaveAtLeastOneTag {
		t.Fatalf("Failed to remove tag. It must keep at least 1 tag %v", f.Tags)
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
	varitions := f.Variations
	patterns := []struct {
		id       string
		strategy *proto.Strategy
		expected error
	}{
		{
			id:       "rule-2",
			strategy: nil,
			expected: errRuleAlreadyExists,
		},
		{
			id: "rule-3",
			strategy: &proto.Strategy{
				Type:          proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{Variation: ""},
			},
			expected: errVariationNotFound,
		},
		{
			id: "rule-3",
			strategy: &proto.Strategy{
				Type:          proto.Strategy_FIXED,
				FixedStrategy: &proto.FixedStrategy{Variation: varitions[0].Id},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		rule := &proto.Rule{
			Id:       p.id,
			Strategy: p.strategy,
		}
		err := f.AddRule(rule)
		assert.Equal(t, p.expected, err)
	}
	rule := &proto.Rule{
		Id:       patterns[2].id,
		Strategy: patterns[2].strategy,
	}
	assert.Equal(t, rule, f.Rules[2])
}

func TestAddRolloutStrategyRule(t *testing.T) {
	f := makeFeature("test-feature")
	varitions := f.Variations
	patterns := []struct {
		id       string
		strategy *proto.Strategy
		expected error
	}{
		{
			id:       "rule-2",
			strategy: nil,
			expected: errRuleAlreadyExists,
		},
		{
			id: "rule-3",
			strategy: &proto.Strategy{
				Type: proto.Strategy_ROLLOUT,
				RolloutStrategy: &proto.RolloutStrategy{
					Variations: []*proto.RolloutStrategy_Variation{
						{
							Variation: varitions[0].Id,
							Weight:    30000,
						},
						{
							Variation: "",
							Weight:    70000,
						},
					},
				},
			},
			expected: errVariationNotFound,
		},
		{
			id: "rule-3",
			strategy: &proto.Strategy{
				Type: proto.Strategy_ROLLOUT,
				RolloutStrategy: &proto.RolloutStrategy{
					Variations: []*proto.RolloutStrategy_Variation{
						{
							Variation: varitions[0].Id,
							Weight:    30000,
						},
						{
							Variation: varitions[1].Id,
							Weight:    70000,
						},
					},
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		rule := &proto.Rule{
			Id:       p.id,
			Strategy: p.strategy,
		}
		err := f.AddRule(rule)
		assert.Equal(t, p.expected, err)
	}
	rule := &proto.Rule{
		Id:       patterns[2].id,
		Strategy: patterns[2].strategy,
	}
	assert.Equal(t, rule, f.Rules[2])
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

func TestListSegmentIDs(t *testing.T) {
	f := makeFeature("test-feature")
	expected := []string{"newUser1", "newUser2"}
	newRule := &proto.Rule{
		Clauses: []*proto.Clause{
			{Operator: proto.Clause_SEGMENT, Values: expected},
		},
	}
	f.Rules = append(f.Rules, newRule)
	actual := f.ListSegmentIDs()
	sort.Strings(actual)
	assert.Equal(t, expected, actual)
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

func TestValidateVariation(t *testing.T) {
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
			assert.Equal(t, p.expected, validateVariation(p.variationType, p.value))
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
