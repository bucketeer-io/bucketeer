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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
				VariationName:  "Variation A",
				VariationValue: "A",
				Variation: &featureproto.Variation{
					Id:    "variation-A",
					Name:  "Variation A",
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
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Name:  "Variation B",
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
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Name:  "Variation B",
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
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Name:  "Variation B",
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
				VariationName:  "Variation A",
				VariationValue: "A",
				Variation: &featureproto.Variation{
					Id:    "variation-A",
					Name:  "Variation A",
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
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &featureproto.Variation{
					Id:    "variation-B",
					Value: "B",
					Name:  "Variation B",
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

func TestEvaluateFeaturesByEvaluatedAt(t *testing.T) {
	t.Parallel()
	now := time.Now()
	thirtyOneDaysAgo := now.Add(-31 * 24 * time.Hour)
	fiveMinutesAgo := now.Add(-5 * time.Minute)
	tenMinutesAgo := now.Add(-10 * time.Minute)
	tenMinutesAndNineSecondsAgo := now.Add(-609 * time.Second)
	tenMinutesAndElevenSecondsAgo := now.Add(-611 * time.Second)
	oneHourAgo := now.Add(-1 * time.Hour)
	user := &userproto.User{Id: "user-1"}
	segmentUser := map[string][]*featureproto.SegmentUser{}

	patterns := []struct {
		desc                   string
		prevUEID               string
		evaluatedAt            int64
		userAttributesUpdated  bool
		tag                    string
		createFeatures         func() []*featureproto.Feature
		expectedEvals          *UserEvaluations
		expectedEvalFeatureIDs []string
		expectedError          error
	}{
		{
			desc:                  "success: evaluate all features since the previous UserEvaluationsID is empty",
			prevUEID:              "",
			evaluatedAt:           thirtyOneDaysAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
					{
						Id:             "feature-2:1:user-1",
						FeatureId:      "feature-2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
				},
				[]string{"feature-3"},
				true,
			),
			expectedEvalFeatureIDs: []string{"feature-1", "feature-2"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate all features since the previous evaluation was over a month ago",
			prevUEID:              "prevUEID",
			evaluatedAt:           thirtyOneDaysAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
					{
						Id:             "feature-2:1:user-1",
						FeatureId:      "feature-2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
				},
				[]string{"feature-3"},
				true,
			),
			expectedEvalFeatureIDs: []string{"feature-1", "feature-2"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate all features since both feature flags and user attributes have not been updated (although the UEID has been updated)",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = oneHourAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = oneHourAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = oneHourAgo.Unix()
				f3.Archived = true
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
					{
						Id:             "feature-2:1:user-1",
						FeatureId:      "feature-2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
				},
				[]string{"feature-3"},
				true,
			),
			expectedEvalFeatureIDs: []string{"feature-1", "feature-2"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate only features updated since the previous evaluations",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = oneHourAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true

				f4 := makeFeature("feature-4")
				f4.UpdatedAt = oneHourAgo.Unix()
				f4.Archived = true
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature, f4.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
				},
				[]string{"feature-3"},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature-1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: check the adjustment seconds",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = tenMinutesAndNineSecondsAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = tenMinutesAndElevenSecondsAgo.Unix()
				return []*featureproto.Feature{f1.Feature, f2.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature-1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate only features has rules when user attributes updated",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: true,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = thirtyOneDaysAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = thirtyOneDaysAgo.Unix()
				f2.Rules = []*featureproto.Rule{}

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = thirtyOneDaysAgo.Unix()
				f3.Archived = true
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_RULE},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature-1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate only the features that have been updated since the previous evaluation, or the features that have rules when user attributes are updated",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: true,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()
				f1.Rules = []*featureproto.Rule{}

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = thirtyOneDaysAgo.Unix()
				f2.Rules = []*featureproto.Rule{}

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true

				f4 := makeFeature("feature-4")
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				f4.Rules = []*featureproto.Rule{}
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature, f4.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_RULE},
					},
					{
						Id:             "feature-4:1:user-1",
						FeatureId:      "feature-4",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_RULE},
					},
				},
				[]string{"feature-3"},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature-1", "feature-4"},
			expectedError:          nil,
		},
		{
			desc:                  "success: prerequisite",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = thirtyOneDaysAgo.Unix()
				f1.Prerequisites = append(f1.Prerequisites, &featureproto.Prerequisite{
					FeatureId:   "feature-4",
					VariationId: "B",
				})

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = thirtyOneDaysAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = thirtyOneDaysAgo.Unix()

				f4 := makeFeature("feature-4")
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature, f4.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_RULE},
					},
					{
						Id:             "feature-4:1:user-1",
						FeatureId:      "feature-4",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_RULE},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature-1", "feature-4"},
			expectedError:          nil,
		},
		{
			desc:                  "success: When a tag is specified, it excludes the evaluations that don't have that tag. But archived features are not excluded",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "tag-1",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.Tags = append(f1.Tags, "tag-1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.Tags = append(f2.Tags, "tag-2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.Tags = append(f3.Tags, "tag-1")
				f3.Archived = true
				f3.UpdatedAt = fiveMinutesAgo.Unix()

				f4 := makeFeature("feature-4")
				f4.Tags = append(f4.Tags, "tag-2")
				f4.Archived = true
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature, f4.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
				},
				[]string{"feature-3", "feature-4"},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature-1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: When a tag is not specified, it does not exclude evaluations that have tags.",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*featureproto.Feature {
				f1 := makeFeature("feature-1")
				f1.Tags = append(f1.Tags, "tag-1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.Tags = append(f2.Tags, "tag-2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()

				f4 := makeFeature("feature-4")
				f4.Tags = append(f4.Tags, "tag-1")
				f4.Tags = append(f4.Tags, "tag-2")
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				return []*featureproto.Feature{f1.Feature, f2.Feature, f3.Feature, f4.Feature}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*featureproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
					{
						Id:             "feature-2:1:user-1",
						FeatureId:      "feature-2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
					{
						Id:             "feature-3:1:user-1",
						FeatureId:      "feature-3",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
					{
						Id:             "feature-4:1:user-1",
						FeatureId:      "feature-4",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &featureproto.Reason{Type: featureproto.Reason_DEFAULT},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature-1", "feature-2", "feature-3", "feature-4"},
			expectedError:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := EvaluateFeaturesByEvaluatedAt(
				p.createFeatures(),
				user,
				segmentUser,
				p.prevUEID,
				p.evaluatedAt,
				p.userAttributesUpdated,
				p.tag,
			)
			assert.Equal(t, p.expectedError, err)
			assert.Equal(t, p.expectedEvals.UserEvaluations.ForceUpdate, actual.ForceUpdate)
			assert.Equal(t, p.expectedEvals.UserEvaluations.ArchivedFeatureIds, actual.ArchivedFeatureIds)
			assert.Equal(t, len(p.expectedEvals.UserEvaluations.Evaluations), len(actual.Evaluations))
			for _, e := range actual.Evaluations {
				assert.Contains(t, p.expectedEvalFeatureIDs, e.FeatureId)
			}
		})
	}

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

var allFeaturesForPrerequisiteTest = map[string]*featureproto.Feature{
	"featureA": {
		Id:   "featureA",
		Name: "featureA",
		Prerequisites: []*featureproto.Prerequisite{
			{
				FeatureId: "featureE",
			},
			{
				FeatureId: "featureF",
			},
		},
	},
	"featureB": {
		Id:   "featureB",
		Name: "featureB",
	},
	"featureC": {
		Id:   "featureC",
		Name: "featureC",
		Prerequisites: []*featureproto.Prerequisite{
			{
				FeatureId: "featureL",
			},
		},
	},
	"featureD": {
		Id:   "featureD",
		Name: "featureD",
	},
	"featureE": {
		Id:   "featureE",
		Name: "featureE",
		Prerequisites: []*featureproto.Prerequisite{
			{
				FeatureId: "featureG",
			},
		},
	},
	"featureF": {
		Id:   "featureF",
		Name: "featureF",
	},
	"featureG": {
		Id:   "featureG",
		Name: "featureG",
		Prerequisites: []*featureproto.Prerequisite{
			{
				FeatureId: "featureH",
			},
		},
	},
	"featureH": {
		Id:   "featureH",
		Name: "featureH",
		Prerequisites: []*featureproto.Prerequisite{
			{
				FeatureId: "featureI",
			},
			{
				FeatureId: "featureJ",
			},
		},
	},
	"featureI": {
		Id:   "featureI",
		Name: "featureI",
		Prerequisites: []*featureproto.Prerequisite{
			{
				FeatureId: "featureK",
			},
		},
	},
	"featureJ": {
		Id:   "featureJ",
		Name: "featureJ",
	},
	"featureK": {
		Id:   "featureK",
		Name: "featureK",
	},
	"featureL": {
		Id:   "featureL",
		Name: "featureL",
		Prerequisites: []*featureproto.Prerequisite{
			{
				FeatureId: "featureM",
			},
			{
				FeatureId: "featureN",
			},
		},
	},
	"featureM": {
		Id:   "featureM",
		Name: "featureM",
	},
	"featureN": {
		Id:   "featureN",
		Name: "featureN",
	},
}

func TestGetPrerequisiteDownwards(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		target      []*featureproto.Feature
		expected    []*featureproto.Feature
		expectedErr error
	}{
		{
			desc: "success: No prerequisites",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureB"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureB"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
			expectedErr: nil,
		},
		{
			desc: "success: Get prerequisites pattern1",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureA"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureA"],
				allFeaturesForPrerequisiteTest["featureE"],
				allFeaturesForPrerequisiteTest["featureF"],
				allFeaturesForPrerequisiteTest["featureG"],
				allFeaturesForPrerequisiteTest["featureH"],
				allFeaturesForPrerequisiteTest["featureI"],
				allFeaturesForPrerequisiteTest["featureJ"],
				allFeaturesForPrerequisiteTest["featureK"],
			},
			expectedErr: nil,
		},
		{
			desc: "success: Get prerequisites pattern2",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureC"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureC"],
				allFeaturesForPrerequisiteTest["featureD"],
				allFeaturesForPrerequisiteTest["featureL"],
				allFeaturesForPrerequisiteTest["featureM"],
				allFeaturesForPrerequisiteTest["featureN"],
			},
			expectedErr: nil,
		},
		{
			desc: "success: Get prerequisites pattern3",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureD"],
				allFeaturesForPrerequisiteTest["featureH"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureD"],
				allFeaturesForPrerequisiteTest["featureH"],
				allFeaturesForPrerequisiteTest["featureI"],
				allFeaturesForPrerequisiteTest["featureJ"],
				allFeaturesForPrerequisiteTest["featureK"],
			},
			expectedErr: nil,
		},
	}
	allFeatures := make([]*featureproto.Feature, 0, len(allFeaturesForPrerequisiteTest))
	for _, v := range allFeaturesForPrerequisiteTest {
		allFeatures = append(allFeatures, v)
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := GetPrerequisiteDownwards(p.target, allFeatures)
			assert.ElementsMatch(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetPrerequisiteUpwards(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc     string
		target   []*featureproto.Feature
		expected []*featureproto.Feature
	}{
		{
			desc: "success: No prerequisites",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureA"],
				allFeaturesForPrerequisiteTest["featureB"],
				allFeaturesForPrerequisiteTest["featureC"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureA"],
				allFeaturesForPrerequisiteTest["featureB"],
				allFeaturesForPrerequisiteTest["featureC"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
		},
		{
			desc: "success: Get prerequisites pattern1",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureF"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureA"],
				allFeaturesForPrerequisiteTest["featureF"],
			},
		},
		{
			desc: "success: Get prerequisites pattern2",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureK"],
				allFeaturesForPrerequisiteTest["featureE"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureA"],
				allFeaturesForPrerequisiteTest["featureE"],
				allFeaturesForPrerequisiteTest["featureG"],
				allFeaturesForPrerequisiteTest["featureH"],
				allFeaturesForPrerequisiteTest["featureI"],
				allFeaturesForPrerequisiteTest["featureK"],
			},
		},
		{
			desc: "success: Get prerequisites pattern3",
			target: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureM"],
				allFeaturesForPrerequisiteTest["featureN"],
			},
			expected: []*featureproto.Feature{
				allFeaturesForPrerequisiteTest["featureC"],
				allFeaturesForPrerequisiteTest["featureL"],
				allFeaturesForPrerequisiteTest["featureM"],
				allFeaturesForPrerequisiteTest["featureN"],
			},
		},
	}
	allFeatures := make([]*featureproto.Feature, 0, len(allFeaturesForPrerequisiteTest))
	for _, v := range allFeaturesForPrerequisiteTest {
		allFeatures = append(allFeatures, v)
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			featuresHavePrerequisite := getFeaturesHavePrerequisite(allFeatures)
			actual := GetPrerequisiteUpwards(p.target, featuresHavePrerequisite)
			assert.ElementsMatch(t, p.expected, actual)
		})
	}
}
