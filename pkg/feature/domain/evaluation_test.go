// Copyright 2022 The Bucketeer Authors.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestEvaluateFeature(t *testing.T) {
	t.Parallel()
	f := makeFeature("fID-0")
	f.Tags = append(f.Tags, "tag-1")
	f1 := makeFeature("fID-1")
	f1.Tags = append(f1.Tags, "tag-1")
	f1.Enabled = false
	f1.OffVariation = f1.Variations[0].Id
	f2 := makeFeature("fID-2")
	f2.Tags = append(f2.Tags, "tag-1")
	patterns := []struct {
		enabled       bool
		offVariation  string
		userID        string
		prerequisite  []*featureproto.Prerequisite
		expected      *featureproto.Evaluation
		expectedError error
	}{
		{
			enabled:       false,
			offVariation:  "not-found",
			userID:        "uID-0",
			prerequisite:  []*featureproto.Prerequisite{},
			expected:      nil,
			expectedError: errVariationNotFound,
		},
		{
			enabled:      false,
			offVariation: "variation-A",
			userID:       "uID-0",
			prerequisite: []*featureproto.Prerequisite{},
			expected: &featureproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID-0"),
				FeatureId:      "fID-0",
				FeatureVersion: 1,
				UserId:         "uID-0",
				VariationId:    "variation-A",
				VariationValue: "A",
				Variation: &featureproto.Variation{
					Id:    "variation-A",
					Value: "A",
				},
				Reason: &featureproto.Reason{Type: featureproto.Reason_OFF_VARIATION},
			},
			expectedError: nil,
		},
		{
			enabled:      false,
			offVariation: "",
			userID:       "uID-0",
			prerequisite: []*featureproto.Prerequisite{},
			expected: &featureproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID-0"),
				FeatureId:      "fID-0",
				FeatureVersion: 1,
				UserId:         "uID-0",
				VariationId:    "variation-B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Value: "B",
				},
				Reason: &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "",
			userID:       "uID-2",
			prerequisite: []*featureproto.Prerequisite{},
			expected: &featureproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID-2"),
				FeatureId:      "fID-0",
				FeatureVersion: 1,
				UserId:         "uID-2",
				VariationId:    "variation-B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Value: "B",
				},
				Reason: &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "v1",
			userID:       "uID-2",
			prerequisite: []*featureproto.Prerequisite{},
			expected: &featureproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID-2"),
				FeatureId:      "fID-0",
				FeatureVersion: 1,
				UserId:         "uID-2",
				VariationId:    "variation-B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Value: "B",
				},
				Reason: &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "variation-A",
			userID:       "uID-2",
			prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId:   f1.Id,
					VariationId: f1.Variations[1].Id,
				},
			},
			expected: &featureproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID-2"),
				FeatureId:      "fID-0",
				FeatureVersion: 1,
				UserId:         "uID-2",
				VariationId:    "variation-A",
				VariationValue: "A",
				Variation: &featureproto.Variation{
					Id:    "variation-A",
					Value: "A",
				},
				Reason: &featureproto.Reason{Type: featureproto.Reason_PREREQUISITE},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "",
			userID:       "uID-2",
			prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId:   f2.Id,
					VariationId: f2.Variations[0].Id,
				},
			},
			expected: &featureproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID-2"),
				FeatureId:      "fID-0",
				FeatureVersion: 1,
				UserId:         "uID-2",
				VariationId:    "variation-B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Value: "B",
				},
				Reason: &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
	}

	for _, p := range patterns {
		user := &userproto.User{Id: p.userID}
		f.Enabled = p.enabled
		f.OffVariation = p.offVariation
		f.Prerequisites = p.prerequisite
		segmentUser := map[string][]*featureproto.SegmentUser{}
		evaluation, err := EvaluateFeatures([]*featureproto.Feature{f.Feature, f1.Feature, f2.Feature}, user, segmentUser, "tag-1")
		assert.Equal(t, p.expectedError, err)
		if evaluation != nil {
			actual, err := findEvaluation(evaluation.Evaluations, f.Id)
			assert.NoError(t, err)
			assert.True(t, proto.Equal(p.expected, actual))
		}
	}
}

func findEvaluation(es []*featureproto.Evaluation, fId string) (*featureproto.Evaluation, error) {
	for _, e := range es {
		if fId == e.FeatureId {
			return e, nil
		}
	}
	return nil, fmt.Errorf("%s was not found", fId)
}

func TestTopologicalSort(t *testing.T) {
	t.Parallel()
	f0 := makeFeature("fID-0")
	f1 := makeFeature("fID-1")
	f2 := makeFeature("fID-2")
	f3 := makeFeature("fID-3")
	f4 := makeFeature("fID-4")
	f5 := makeFeature("fID-5")
	patterns := []struct {
		f0Prerequisite []*featureproto.Prerequisite
		f1Prerequisite []*featureproto.Prerequisite
		f2Prerequisite []*featureproto.Prerequisite
		f3Prerequisite []*featureproto.Prerequisite
		f4Prerequisite []*featureproto.Prerequisite
		f5Prerequisite []*featureproto.Prerequisite
		expected       []*featureproto.Feature
		expectedError  error
	}{
		{
			f0Prerequisite: []*featureproto.Prerequisite{},
			f1Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
			},
			f2Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
			},
			f3Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
				{
					FeatureId: f2.Id,
				},
			},
			f4Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			f5Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f4.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			expected: []*featureproto.Feature{
				f0.Feature, f1.Feature, f2.Feature, f3.Feature, f4.Feature, f5.Feature,
			},
			expectedError: nil,
		},
		{
			f0Prerequisite: []*featureproto.Prerequisite{},
			f1Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
			},
			f2Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
			},
			f3Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f1.Id,
				},
				{
					FeatureId: f2.Id,
				},
			},
			f4Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			f5Prerequisite: []*featureproto.Prerequisite{},
			expected: []*featureproto.Feature{
				f0.Feature, f1.Feature, f2.Feature, f5.Feature, f3.Feature, f4.Feature,
			},
			expectedError: nil,
		},
		{
			f0Prerequisite: []*featureproto.Prerequisite{},
			f1Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
			},
			f2Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f3.Id,
				},
			},
			f3Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f2.Id,
				},
			},
			f4Prerequisite: []*featureproto.Prerequisite{
				{
					FeatureId: f0.Id,
				},
				{
					FeatureId: f3.Id,
				},
			},
			f5Prerequisite: []*featureproto.Prerequisite{
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
			f0Prerequisite: []*featureproto.Prerequisite{},
			f1Prerequisite: []*featureproto.Prerequisite{},
			f2Prerequisite: []*featureproto.Prerequisite{},
			f3Prerequisite: []*featureproto.Prerequisite{},
			f4Prerequisite: []*featureproto.Prerequisite{},
			f5Prerequisite: []*featureproto.Prerequisite{},
			expected: []*featureproto.Feature{
				f2.Feature, f0.Feature, f5.Feature, f3.Feature, f1.Feature, f4.Feature,
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
		fs := []*featureproto.Feature{
			f2.Feature, f0.Feature, f5.Feature, f3.Feature, f1.Feature, f4.Feature,
		}
		actual, err := TopologicalSort(fs)
		assert.Equal(t, p.expectedError, err)
		assert.Equal(t, p.expected, actual)
	}
}
