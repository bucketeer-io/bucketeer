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

package evaluation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const holdoutConformanceFixturePath = "../testdata/holdout_conformance.json"

type holdoutConformanceTarget struct {
	FeatureID   string `json:"featureId"`
	VariationID string `json:"variationId"`
}

type holdoutConformanceGroup struct {
	ID             string                     `json:"id"`
	AssignmentType string                     `json:"assignmentType"`
	BucketStart    float64                    `json:"bucketStart"`
	Percentage     float64                    `json:"percentage"`
	SegmentID      string                     `json:"segmentId"`
	Targets        []holdoutConformanceTarget `json:"targets"`
}

type holdoutConformancePrerequisite struct {
	FeatureID   string `json:"featureId"`
	VariationID string `json:"variationId"`
}

type holdoutConformanceFeature struct {
	ID               string                           `json:"id"`
	Enabled          bool                             `json:"enabled"`
	OffVariation     string                           `json:"offVariation"`
	Variations       []string                         `json:"variations"`
	DefaultVariation string                           `json:"defaultVariation"`
	Prerequisites    []holdoutConformancePrerequisite `json:"prerequisites"`
}

type holdoutConformanceExpected struct {
	FeatureID      string `json:"featureId"`
	VariationID    string `json:"variationId"`
	Reason         string `json:"reason"`
	HoldoutGroupID string `json:"holdoutGroupId"`
}

type holdoutConformanceCase struct {
	Desc     string                       `json:"desc"`
	UserID   string                       `json:"userId"`
	Group    holdoutConformanceGroup      `json:"group"`
	Features []holdoutConformanceFeature  `json:"features"`
	Expected []holdoutConformanceExpected `json:"expected"`
}

type holdoutConformanceFixture struct {
	EnvironmentID string                   `json:"environmentId"`
	Seed          string                   `json:"seed"`
	SegmentUsers  []conformanceSegmentUser `json:"segmentUsers"`
	Cases         []holdoutConformanceCase `json:"cases"`
}

func TestHoldoutConformance(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Clean(holdoutConformanceFixturePath))
	require.NoError(t, err)
	fixture := &holdoutConformanceFixture{}
	require.NoError(t, json.Unmarshal(data, fixture))
	segmentUsers := map[string][]*ftproto.SegmentUser{}
	for _, su := range fixture.SegmentUsers {
		state, ok := ftproto.SegmentUser_State_value[su.State]
		require.True(t, ok, "unknown segment user state: %s", su.State)
		segmentUsers[su.SegmentID] = append(segmentUsers[su.SegmentID], &ftproto.SegmentUser{
			SegmentId: su.SegmentID,
			UserId:    su.UserID,
			State:     ftproto.SegmentUser_State(state),
		})
	}
	evaluator := NewEvaluator()
	for _, tc := range fixture.Cases {
		tc := tc
		t.Run(tc.Desc, func(t *testing.T) {
			t.Parallel()
			group := conformanceGroup(t, fixture, tc.Group)
			features := conformanceFeatures(tc.Features)
			evaluations, err := evaluator.EvaluateFeatures(
				features,
				&userproto.User{Id: tc.UserID},
				segmentUsers,
				nil,
				"",
				group,
			)
			require.NoError(t, err)
			assert.Equal(t, conformanceExpected(t, tc.Expected), conformanceActual(evaluations.Evaluations))
		})
	}
}

func conformanceGroup(
	t *testing.T,
	fixture *holdoutConformanceFixture,
	group holdoutConformanceGroup,
) *ftproto.HoldoutGroup {
	t.Helper()
	assignment, ok := ftproto.HoldoutGroup_AssignmentType_value[group.AssignmentType]
	require.True(t, ok, "unknown assignment type: %s", group.AssignmentType)
	features := make([]*ftproto.HoldoutGroupFeature, 0, len(group.Targets))
	for _, target := range group.Targets {
		features = append(features, &ftproto.HoldoutGroupFeature{
			FeatureId:          target.FeatureID,
			HoldoutVariationId: target.VariationID,
		})
	}
	return &ftproto.HoldoutGroup{
		Id:             group.ID,
		EnvironmentId:  fixture.EnvironmentID,
		AssignmentType: ftproto.HoldoutGroup_AssignmentType(assignment),
		SamplingSeed:   fixture.Seed,
		BucketStart:    group.BucketStart,
		Percentage:     group.Percentage,
		SegmentId:      group.SegmentID,
		Features:       features,
	}
}

func conformanceFeatures(features []holdoutConformanceFeature) []*ftproto.Feature {
	out := make([]*ftproto.Feature, 0, len(features))
	for _, feature := range features {
		variations := make([]*ftproto.Variation, 0, len(feature.Variations))
		for _, id := range feature.Variations {
			variations = append(variations, &ftproto.Variation{Id: id, Value: id})
		}
		prerequisites := make([]*ftproto.Prerequisite, 0, len(feature.Prerequisites))
		for _, prerequisite := range feature.Prerequisites {
			prerequisites = append(prerequisites, &ftproto.Prerequisite{
				FeatureId:   prerequisite.FeatureID,
				VariationId: prerequisite.VariationID,
			})
		}
		out = append(out, &ftproto.Feature{
			Id:            feature.ID,
			Enabled:       feature.Enabled,
			OffVariation:  feature.OffVariation,
			Variations:    variations,
			Prerequisites: prerequisites,
			DefaultStrategy: &ftproto.Strategy{
				Type:          ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{Variation: feature.DefaultVariation},
			},
		})
	}
	return out
}

type holdoutConformanceResult struct {
	FeatureID      string
	VariationID    string
	Reason         string
	HoldoutGroupID string
}

func conformanceExpected(t *testing.T, expected []holdoutConformanceExpected) []holdoutConformanceResult {
	t.Helper()
	out := make([]holdoutConformanceResult, 0, len(expected))
	for _, item := range expected {
		_, ok := ftproto.Reason_Type_value[item.Reason]
		require.True(t, ok, "unknown reason: %s", item.Reason)
		out = append(out, holdoutConformanceResult{
			FeatureID:      item.FeatureID,
			VariationID:    item.VariationID,
			Reason:         item.Reason,
			HoldoutGroupID: item.HoldoutGroupID,
		})
	}
	return out
}

func conformanceActual(evaluations []*ftproto.Evaluation) []holdoutConformanceResult {
	out := make([]holdoutConformanceResult, 0, len(evaluations))
	for _, evaluation := range evaluations {
		out = append(out, holdoutConformanceResult{
			FeatureID:      evaluation.FeatureId,
			VariationID:    evaluation.VariationId,
			Reason:         evaluation.Reason.Type.String(),
			HoldoutGroupID: evaluation.Reason.HoldoutGroupId,
		})
	}
	return out
}
