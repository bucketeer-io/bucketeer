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

package command

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestAddFixedStrategyRule(t *testing.T) {
	f := makeFeature("feature-id")
	id, _ := uuid.NewUUID()
	rID := id.String()
	vID := f.Variations[0].Id
	expected := &proto.Rule{
		Id: rID,
		Strategy: &proto.Strategy{
			Type:          proto.Strategy_FIXED,
			FixedStrategy: &proto.FixedStrategy{Variation: vID},
		},
	}
	patterns := []*struct {
		rule     *proto.Rule
		expected error
	}{
		{
			rule:     expected,
			expected: nil,
		},
	}
	targetingCmd := &FeatureCommandHandler{
		feature:      f,
		eventFactory: makeEventFactory(f),
	}
	for i, p := range patterns {
		cmd := &proto.AddRuleCommand{Rule: p.rule}
		err := targetingCmd.Handle(context.Background(), cmd)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}
	if !reflect.DeepEqual(expected, f.Rules[1]) {
		t.Fatalf("Rule is not equal. Expected: %v, actual: %v", expected, f.Rules[1])
	}
}

func TestAddRolloutStrategyRule(t *testing.T) {
	f := makeFeature("feature-id")
	id, _ := uuid.NewUUID()
	rID := id.String()
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	expected := &proto.Rule{
		Id: rID,
		Strategy: &proto.Strategy{
			Type: proto.Strategy_ROLLOUT,
			RolloutStrategy: &proto.RolloutStrategy{
				Variations: []*proto.RolloutStrategy_Variation{
					{
						Variation: vID1,
						Weight:    30000,
					},
					{
						Variation: vID2,
						Weight:    70000,
					},
				},
			},
		},
	}
	patterns := []*struct {
		rule     *proto.Rule
		expected error
	}{
		{
			rule:     expected,
			expected: nil,
		},
	}
	targetingCmd := &FeatureCommandHandler{
		feature:      f,
		eventFactory: makeEventFactory(f),
	}
	for i, p := range patterns {
		cmd := &proto.AddRuleCommand{Rule: p.rule}
		err := targetingCmd.Handle(context.Background(), cmd)
		des := fmt.Sprintf("index: %d", i)
		assert.Equal(t, p.expected, err, des)
	}
	if !reflect.DeepEqual(expected, f.Rules[1]) {
		t.Fatalf("Rule is not equal. Expected: %v, actual: %v", expected, f.Rules[1])
	}
}

func TestChangeRuleToFixedStrategy(t *testing.T) {
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID := f.Variations[0].Id
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
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	targetingCmd := &FeatureCommandHandler{
		feature:      f,
		eventFactory: makeEventFactory(f),
	}
	for _, p := range patterns {
		cmd := &proto.ChangeRuleStrategyCommand{
			RuleId:   p.ruleID,
			Strategy: p.strategy,
		}
		err := targetingCmd.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
	if !reflect.DeepEqual(expected, r.Strategy) {
		t.Fatalf("Strategy is not equal. Expected: %v, actual: %v", expected, r.Strategy)
	}
}

func TestChangeRuleToRolloutStrategy(t *testing.T) {
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	expected := &proto.Strategy{
		Type: proto.Strategy_ROLLOUT,
		RolloutStrategy: &proto.RolloutStrategy{
			Variations: []*proto.RolloutStrategy_Variation{
				{
					Variation: vID1,
					Weight:    30000,
				},
				{
					Variation: vID2,
					Weight:    70000,
				},
			},
		},
	}
	patterns := []struct {
		desc     string
		ruleID   string
		strategy *proto.Strategy
		expected error
	}{
		{
			desc:     "success",
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	targetingCmd := &FeatureCommandHandler{
		feature:      f,
		eventFactory: makeEventFactory(f),
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cmd := &proto.ChangeRuleStrategyCommand{
				RuleId:   p.ruleID,
				Strategy: p.strategy,
			}
			err := targetingCmd.Handle(context.Background(), cmd)
			assert.Equal(t, p.expected, err)
		})
	}
	if !reflect.DeepEqual(expected, r.Strategy) {
		t.Fatalf("Strategy is not equal. Expected: %v, actual: %v", expected, r.Strategy)
	}
}

func TestChangeFixedStrategy(t *testing.T) {
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID := f.Variations[0].Id
	patterns := []*struct {
		ruleID   string
		strategy *proto.FixedStrategy
		expected error
	}{
		{
			ruleID:   rID,
			strategy: &proto.FixedStrategy{Variation: vID},
			expected: nil,
		},
	}
	targetingCmd := &FeatureCommandHandler{
		feature:      f,
		eventFactory: makeEventFactory(f),
	}
	for _, p := range patterns {
		cmd := &proto.ChangeFixedStrategyCommand{
			RuleId:   p.ruleID,
			Strategy: p.strategy,
		}
		err := targetingCmd.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
	if r.Strategy.FixedStrategy.Variation != vID {
		t.Fatalf("Wrong variation id has been saved. Expected: %s, actual: %s", vID, r.Strategy.FixedStrategy.Variation)
	}
}

func TestChangeRolloutStrategy(t *testing.T) {
	f := makeFeature("feature-id")
	r := f.Rules[0]
	rID := r.Id
	vID1 := f.Variations[0].Id
	vID2 := f.Variations[1].Id
	expected := &proto.RolloutStrategy{Variations: []*proto.RolloutStrategy_Variation{
		{
			Variation: vID1,
			Weight:    70000,
		},
		{
			Variation: vID2,
			Weight:    30000,
		},
	}}
	patterns := []*struct {
		ruleID   string
		strategy *proto.RolloutStrategy
		expected error
	}{
		{
			ruleID:   rID,
			strategy: expected,
			expected: nil,
		},
	}
	targetingCmd := &FeatureCommandHandler{
		feature:      f,
		eventFactory: makeEventFactory(f),
	}
	for _, p := range patterns {
		cmd := &proto.ChangeRolloutStrategyCommand{
			RuleId:   p.ruleID,
			Strategy: p.strategy,
		}
		err := targetingCmd.Handle(context.Background(), cmd)
		assert.Equal(t, p.expected, err)
	}
	if !reflect.DeepEqual(expected, r.Strategy.RolloutStrategy) {
		t.Fatalf("Different rollout strategies. Expected: %v, actual: %v", expected, r.Strategy.RolloutStrategy)
	}
}

func TestChangeDefaultStrategy(t *testing.T) {
	patterns := []struct {
		desc        string
		strategy    *proto.Strategy
		expectedErr error
	}{
		{
			desc: "success",
			strategy: &proto.Strategy{
				Type: proto.Strategy_ROLLOUT,
				RolloutStrategy: &proto.RolloutStrategy{
					Variations: []*proto.RolloutStrategy_Variation{
						{
							Variation: "variation-A",
							Weight:    30000,
						},
						{
							Variation: "variation-B",
							Weight:    70000,
						},
					},
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			f := makeFeature("feature-id")
			targetingCmd := &FeatureCommandHandler{
				feature:      f,
				eventFactory: makeEventFactory(f),
			}
			cmd := &proto.ChangeDefaultStrategyCommand{
				Strategy: p.strategy,
			}
			err := targetingCmd.Handle(context.Background(), cmd)
			assert.Equal(t, p.expectedErr, err)
			if p.expectedErr != nil {
				return
			}
			assert.Equal(t, p.strategy, f.DefaultStrategy)
		})
	}
}

func TestEnableFeature(t *testing.T) {
	patterns := []struct {
		desc     string
		cmd      *proto.EnableFeatureCommand
		expected error
	}{
		{
			desc:     "success",
			cmd:      &proto.EnableFeatureCommand{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		f := makeFeature("feature-id")
		cmd := &FeatureCommandHandler{
			feature:      f,
			eventFactory: makeEventFactory(f),
		}
		err := cmd.Handle(ctx, p.cmd)
		assert.Equal(t, p.expected, err, p.desc)
		assert.True(t, f.Feature.Enabled, p.desc)
	}
}

func TestDisableFeature(t *testing.T) {
	patterns := []struct {
		desc     string
		cmd      *proto.DisableFeatureCommand
		expected error
	}{
		{
			desc:     "success",
			cmd:      &proto.DisableFeatureCommand{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		f := makeFeature("feature-id")
		f.Feature.Enabled = true
		cmd := &FeatureCommandHandler{
			feature:      f,
			eventFactory: makeEventFactory(f),
		}
		err := cmd.Handle(ctx, p.cmd)
		assert.Equal(t, p.expected, err, p.desc)
		assert.False(t, f.Feature.Enabled, p.desc)
	}
}

func TestResetSamplingSeed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	patterns := []struct {
		desc     string
		cmd      *proto.ResetSamplingSeedCommand
		expected error
	}{
		{
			desc:     "success",
			cmd:      &proto.ResetSamplingSeedCommand{},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			f := makeFeature("fid")
			assert.Empty(t, f.Feature.SamplingSeed)
			cmd := &FeatureCommandHandler{
				feature:      f,
				eventFactory: makeEventFactory(f),
			}
			err := cmd.Handle(ctx, p.cmd)
			assert.Equal(t, p.expected, err)
			assert.NotEmpty(t, f.Feature.SamplingSeed)
		})
	}
}

func TestAddPrerequisite(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	patterns := []struct {
		desc     string
		cmd      *proto.AddPrerequisiteCommand
		expected error
	}{
		{
			desc: "success",
			cmd: &proto.AddPrerequisiteCommand{
				Prerequisite: &proto.Prerequisite{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			f := makeFeature("fid")
			assert.Empty(t, f.Feature.Prerequisites)
			cmd := &FeatureCommandHandler{
				feature:      f,
				eventFactory: makeEventFactory(f),
			}
			err := cmd.Handle(ctx, p.cmd)
			assert.Equal(t, p.expected, err)
			assert.Equal(t, 1, len(f.Feature.Prerequisites))
			assert.True(t, protobuf.Equal(p.cmd.Prerequisite, f.Feature.Prerequisites[0]))
		})
	}
}

func TestRemovePrerequisite(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	patterns := []struct {
		desc         string
		cmd          *proto.RemovePrerequisiteCommand
		prerequisite []*proto.Prerequisite
		expected     error
	}{
		{
			desc: "success",
			cmd: &proto.RemovePrerequisiteCommand{
				FeatureId: "test-feature2",
			},
			prerequisite: []*proto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expected: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			f := makeFeature("fid")
			f.Prerequisites = p.prerequisite
			assert.NotEmpty(t, f.Feature.Prerequisites)
			cmd := &FeatureCommandHandler{
				feature:      f,
				eventFactory: makeEventFactory(f),
			}
			err := cmd.Handle(ctx, p.cmd)
			assert.Equal(t, p.expected, err)
			assert.Empty(t, f.Feature.Prerequisites)
		})
	}
}

func TestChangePrerequisiteVariation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	patterns := []struct {
		desc              string
		cmd               *proto.ChangePrerequisiteVariationCommand
		prerequisite      []*proto.Prerequisite
		expectedErr       error
		expectedVariation string
	}{
		{
			desc: "success",
			cmd: &proto.ChangePrerequisiteVariationCommand{
				Prerequisite: &proto.Prerequisite{
					FeatureId:   "test-feature2",
					VariationId: "variation A",
				},
			},
			prerequisite: []*proto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expectedErr:       nil,
			expectedVariation: "variation A",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			f := makeFeature("fid")
			f.Prerequisites = p.prerequisite
			assert.NotEmpty(t, f.Feature.Prerequisites)
			cmd := &FeatureCommandHandler{
				feature:      f,
				eventFactory: makeEventFactory(f),
			}
			err := cmd.Handle(ctx, p.cmd)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expectedVariation, f.Prerequisites[0].VariationId)
		})
	}
}

func makeFeature(id string) *domain.Feature {
	return &domain.Feature{
		Feature: &proto.Feature{
			Id:        id,
			Name:      "test feature",
			Version:   1,
			CreatedAt: time.Now().Unix(),
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
			},
			Targets: []*proto.Target{
				{
					Variation: "variation-B",
					Users: []string{
						"user1",
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

func makeEventFactory(feature *domain.Feature) *FeatureEventFactory {
	return &FeatureEventFactory{
		editor: &eventproto.Editor{
			Email: "email",
		},
		feature:              feature,
		environmentNamespace: "ns0",
	}
}
