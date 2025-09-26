// Copyright 2025 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evaluation

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	featuredoman "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/proto/feature"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestEvaluateFeature(t *testing.T) {
	t.Parallel()
	f := makeFeature("fID0")
	f.Tags = append(f.Tags, "tag1")
	f1 := makeFeature("fID1")
	f1.Tags = append(f1.Tags, "tag1")
	f1.Enabled = false
	f1.OffVariation = f1.Variations[0].Id
	f2 := makeFeature("fID2")
	f2.Tags = append(f2.Tags, "tag1")
	patterns := []struct {
		enabled       bool
		offVariation  string
		userID        string
		prerequisite  []*ftproto.Prerequisite
		expected      *ftproto.Evaluation
		expectedError error
	}{
		{
			enabled:       false,
			offVariation:  "notfound",
			userID:        "uID0",
			prerequisite:  []*ftproto.Prerequisite{},
			expected:      nil,
			expectedError: ErrVariationNotFound,
		},
		{
			enabled:      false,
			offVariation: "variation-A",
			userID:       "uID0",
			prerequisite: []*ftproto.Prerequisite{},
			expected: &ftproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID0"),
				FeatureId:      "fID0",
				FeatureVersion: 1,
				UserId:         "uID0",
				VariationId:    "variation-A",
				VariationName:  "Variation A",
				VariationValue: "A",
				Variation: &ftproto.Variation{
					Id:    "variation-A",
					Name:  "Variation A",
					Value: "A",
				},
				Reason: &ftproto.Reason{Type: ftproto.Reason_OFF_VARIATION},
			},
			expectedError: nil,
		},
		{
			enabled:      false,
			offVariation: "",
			userID:       "uID0",
			prerequisite: []*ftproto.Prerequisite{},
			expected: &ftproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID0"),
				FeatureId:      "fID0",
				FeatureVersion: 1,
				UserId:         "uID0",
				VariationId:    "variation-B",
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &ftproto.Variation{
					Id:    "variation-B",
					Name:  "Variation B",
					Value: "B",
				},
				Reason: &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "",
			userID:       "uID2",
			prerequisite: []*ftproto.Prerequisite{},
			expected: &ftproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID2"),
				FeatureId:      "fID0",
				FeatureVersion: 1,
				UserId:         "uID2",
				VariationId:    "variation-B",
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &ftproto.Variation{
					Id:    "variation-B",
					Name:  "Variation B",
					Value: "B",
				},
				Reason: &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "v1",
			userID:       "uID2",
			prerequisite: []*ftproto.Prerequisite{},
			expected: &ftproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID2"),
				FeatureId:      "fID0",
				FeatureVersion: 1,
				UserId:         "uID2",
				VariationId:    "variation-B",
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &ftproto.Variation{
					Id:    "variation-B",
					Name:  "Variation B",
					Value: "B",
				},
				Reason: &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "variation-A",
			userID:       "uID2",
			prerequisite: []*ftproto.Prerequisite{
				{
					FeatureId:   f1.Id,
					VariationId: f1.Variations[1].Id,
				},
			},
			expected: &ftproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID2"),
				FeatureId:      "fID0",
				FeatureVersion: 1,
				UserId:         "uID2",
				VariationId:    "variation-A",
				VariationName:  "Variation A",
				VariationValue: "A",
				Variation: &ftproto.Variation{
					Id:    "variation-A",
					Name:  "Variation A",
					Value: "A",
				},
				Reason: &ftproto.Reason{Type: ftproto.Reason_PREREQUISITE},
			},
			expectedError: nil,
		},
		{
			enabled:      true,
			offVariation: "",
			userID:       "uID2",
			prerequisite: []*ftproto.Prerequisite{
				{
					FeatureId:   f2.Id,
					VariationId: f2.Variations[0].Id,
				},
			},
			expected: &ftproto.Evaluation{
				Id:             EvaluationID(f.Id, f.Version, "uID2"),
				FeatureId:      "fID0",
				FeatureVersion: 1,
				UserId:         "uID2",
				VariationId:    "variation-B",
				VariationName:  "Variation B",
				VariationValue: "B",
				Variation: &ftproto.Variation{
					Id:    "variation-B",
					Value: "B",
					Name:  "Variation B",
				},
				Reason: &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
	}

	for _, p := range patterns {
		evaluator := NewEvaluator()
		user := &userproto.User{Id: p.userID}
		f.Enabled = p.enabled
		f.OffVariation = p.offVariation
		f.Prerequisites = p.prerequisite
		segmentUser := map[string][]*ftproto.SegmentUser{}
		evaluation, err := evaluator.EvaluateFeatures([]*ftproto.Feature{f, f1, f2}, user, segmentUser, "tag1")
		assert.Equal(t, p.expectedError, err)
		if evaluation != nil {
			actual, err := findEvaluation(evaluation.Evaluations, f.Id)
			assert.NoError(t, err)
			proto.Equal(p.expected, actual)
		}
	}
}

func TestEvaluateFeaturesByEvaluatedAt_MissingPrerequisite(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluator()

	features := []*ftproto.Feature{
		makeDependentFeature(),
	}

	user := &userproto.User{Id: "user-1"}
	segmentUsersMap := map[string][]*ftproto.SegmentUser{}

	_, err := evaluator.EvaluateFeaturesByEvaluatedAt(
		features,
		user,
		segmentUsersMap,
		"prev-ueid",
		time.Now().Unix(),
		false,
		"test",
	)

	require.Error(t, err)
	require.ErrorIs(t, err, featuredoman.ErrFeatureNotFound)
}

func findEvaluation(es []*ftproto.Evaluation, fId string) (*ftproto.Evaluation, error) {
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
	segmentUser := map[string][]*ftproto.SegmentUser{}

	patterns := []struct {
		desc                   string
		prevUEID               string
		evaluatedAt            int64
		userAttributesUpdated  bool
		tag                    string
		createFeatures         func() []*ftproto.Feature
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
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true
				return []*ftproto.Feature{f1, f2, f3}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
					{
						Id:             "feature2:1:user1",
						FeatureId:      "feature2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
				},
				[]string{"feature3"},
				true,
			),
			expectedEvalFeatureIDs: []string{"feature1", "feature2"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate all features since the previous evaluation was over a month ago",
			prevUEID:              "prevUEID",
			evaluatedAt:           thirtyOneDaysAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true
				return []*ftproto.Feature{f1, f2, f3}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
					{
						Id:             "feature2:1:user1",
						FeatureId:      "feature2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
				},
				[]string{"feature3"},
				true,
			),
			expectedEvalFeatureIDs: []string{"feature1", "feature2"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate all features since both feature flags and user attributes have not been updated (although the UEID has been updated)",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature-1")
				f1.UpdatedAt = oneHourAgo.Unix()

				f2 := makeFeature("feature-2")
				f2.UpdatedAt = oneHourAgo.Unix()

				f3 := makeFeature("feature-3")
				f3.UpdatedAt = oneHourAgo.Unix()
				f3.Archived = true
				return []*ftproto.Feature{f1, f2, f3}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature-1:1:user-1",
						FeatureId:      "feature-1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
					{
						Id:             "feature-2:1:user-1",
						FeatureId:      "feature-2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
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
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature2")
				f2.UpdatedAt = oneHourAgo.Unix()

				f3 := makeFeature("feature3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true

				f4 := makeFeature("feature4")
				f4.UpdatedAt = oneHourAgo.Unix()
				f4.Archived = true
				return []*ftproto.Feature{f1, f2, f3, f4}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
				},
				[]string{"feature3"},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: check the adjustment seconds",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = tenMinutesAndNineSecondsAgo.Unix()

				f2 := makeFeature("feature2")
				f2.UpdatedAt = tenMinutesAndElevenSecondsAgo.Unix()
				return []*ftproto.Feature{f1, f2}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate only features has rules when user attributes updated",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: true,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = thirtyOneDaysAgo.Unix()

				f2 := makeFeature("feature2")
				f2.UpdatedAt = thirtyOneDaysAgo.Unix()
				f2.Rules = []*ftproto.Rule{}

				f3 := makeFeature("feature3")
				f3.UpdatedAt = thirtyOneDaysAgo.Unix()
				f3.Archived = true
				return []*ftproto.Feature{f1, f2, f3}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: evaluate only the features that have been updated since the previous evaluation, or the features that have rules when user attributes are updated",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: true,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()
				f1.Rules = []*ftproto.Rule{}

				f2 := makeFeature("feature2")
				f2.UpdatedAt = thirtyOneDaysAgo.Unix()
				f2.Rules = []*ftproto.Rule{}

				f3 := makeFeature("feature3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()
				f3.Archived = true

				f4 := makeFeature("feature4")
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				f4.Rules = []*ftproto.Rule{}
				return []*ftproto.Feature{f1, f2, f3, f4}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
					{
						Id:             "feature4:1:user1",
						FeatureId:      "feature4",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
				},
				[]string{"feature3"},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1", "feature4"},
			expectedError:          nil,
		},
		{
			desc:                  "success: prerequisite",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = thirtyOneDaysAgo.Unix()
				f1.Prerequisites = append(f1.Prerequisites, &ftproto.Prerequisite{
					FeatureId:   "feature4",
					VariationId: "B",
				})

				f2 := makeFeature("feature2")
				f2.UpdatedAt = thirtyOneDaysAgo.Unix()

				f3 := makeFeature("feature3")
				f3.UpdatedAt = thirtyOneDaysAgo.Unix()

				f4 := makeFeature("feature4")
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				return []*ftproto.Feature{f1, f2, f3, f4}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
					{
						Id:             "feature4:1:user1",
						FeatureId:      "feature4",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1", "feature4"},
			expectedError:          nil,
		},
		{
			desc:                  "success: When a tag is specified, it excludes the evaluations that don't have that tag. But archived features are not excluded",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "tag1",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.Tags = append(f1.Tags, "tag1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature2")
				f2.Tags = append(f2.Tags, "tag2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature3")
				f3.Tags = append(f3.Tags, "tag1")
				f3.Archived = true
				f3.UpdatedAt = fiveMinutesAgo.Unix()

				f4 := makeFeature("feature4")
				f4.Tags = append(f4.Tags, "tag2")
				f4.Archived = true
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				return []*ftproto.Feature{f1, f2, f3, f4}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
				},
				[]string{"feature3", "feature4"},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1"},
			expectedError:          nil,
		},
		{
			desc:                  "success: When a tag is not specified, it does not exclude evaluations that have tags.",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.Tags = append(f1.Tags, "tag1")
				f1.UpdatedAt = fiveMinutesAgo.Unix()

				f2 := makeFeature("feature2")
				f2.Tags = append(f2.Tags, "tag2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()

				f3 := makeFeature("feature3")
				f3.UpdatedAt = fiveMinutesAgo.Unix()

				f4 := makeFeature("feature4")
				f4.Tags = append(f4.Tags, "tag1")
				f4.Tags = append(f4.Tags, "tag2")
				f4.UpdatedAt = fiveMinutesAgo.Unix()
				return []*ftproto.Feature{f1, f2, f3, f4}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
					{
						Id:             "feature2:1:user1",
						FeatureId:      "feature2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
					{
						Id:             "feature3:1:user1",
						FeatureId:      "feature3",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
					{
						Id:             "feature4:1:user1",
						FeatureId:      "feature4",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1", "feature2", "feature3", "feature4"},
			expectedError:          nil,
		},
		{
			desc:                  "success: including up/downwards features of target feature with prerequisite",
			prevUEID:              "prevUEID",
			evaluatedAt:           tenMinutesAgo.Unix(),
			userAttributesUpdated: false,
			tag:                   "",
			createFeatures: func() []*ftproto.Feature {
				f1 := makeFeature("feature1")
				f1.UpdatedAt = oneHourAgo.Unix()
				f1.Prerequisites = []*ftproto.Prerequisite{{
					FeatureId:   "feature2",
					VariationId: "B",
				}}

				f2 := makeFeature("feature2")
				f2.UpdatedAt = fiveMinutesAgo.Unix()
				f2.Prerequisites = []*ftproto.Prerequisite{{
					FeatureId:   "feature3",
					VariationId: "B",
				}}
				f3 := makeFeature("feature3")
				f3.UpdatedAt = oneHourAgo.Unix()
				return []*ftproto.Feature{f1, f2, f3}
			},
			expectedEvals: NewUserEvaluations(
				"dummy",
				[]*ftproto.Evaluation{
					{
						Id:             "feature1:1:user1",
						FeatureId:      "feature1",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
					{
						Id:             "feature2:1:user1",
						FeatureId:      "feature2",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
					{
						Id:             "feature3:1:user1",
						FeatureId:      "feature3",
						VariationId:    "variation-B",
						VariationName:  "Variation B",
						VariationValue: "B",
						Reason:         &ftproto.Reason{Type: ftproto.Reason_RULE},
					},
				},
				[]string{},
				false,
			),
			expectedEvalFeatureIDs: []string{"feature1", "feature2", "feature3"},
			expectedError:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			actual, err := evaluator.EvaluateFeaturesByEvaluatedAt(
				p.createFeatures(),
				user,
				segmentUser,
				p.prevUEID,
				p.evaluatedAt,
				p.userAttributesUpdated,
				p.tag,
			)
			assert.Equal(t, p.expectedError, err)
			if p.expectedError != nil {
				return
			}
			assert.Equal(t, p.expectedEvals.UserEvaluations.ForceUpdate, actual.ForceUpdate)
			assert.ElementsMatch(t, p.expectedEvals.UserEvaluations.ArchivedFeatureIds, actual.ArchivedFeatureIds)
			assert.Equal(t, len(p.expectedEvals.UserEvaluations.Evaluations), len(actual.Evaluations))
			for _, e := range actual.Evaluations {
				assert.Contains(t, p.expectedEvalFeatureIDs, e.FeatureId)
			}
		})
	}
}

var allFeaturesForPrerequisiteTest = map[string]*ftproto.Feature{
	"featureA": {
		Id:   "featureA",
		Name: "featureA",
		Prerequisites: []*ftproto.Prerequisite{
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
		Prerequisites: []*ftproto.Prerequisite{
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
		Prerequisites: []*ftproto.Prerequisite{
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
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureH",
			},
		},
	},
	"featureH": {
		Id:   "featureH",
		Name: "featureH",
		Prerequisites: []*ftproto.Prerequisite{
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
		Prerequisites: []*ftproto.Prerequisite{
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
		Prerequisites: []*ftproto.Prerequisite{
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
		target      []*ftproto.Feature
		expected    []*ftproto.Feature
		expectedErr error
	}{
		{
			desc: "success: No prerequisites",
			target: []*ftproto.Feature{
				allFeaturesForPrerequisiteTest["featureB"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
			expected: []*ftproto.Feature{
				allFeaturesForPrerequisiteTest["featureB"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
			expectedErr: nil,
		},
		{
			desc: "success: Get prerequisites pattern1",
			target: []*ftproto.Feature{
				allFeaturesForPrerequisiteTest["featureA"],
			},
			expected: []*ftproto.Feature{
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
			target: []*ftproto.Feature{
				allFeaturesForPrerequisiteTest["featureC"],
				allFeaturesForPrerequisiteTest["featureD"],
			},
			expected: []*ftproto.Feature{
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
			target: []*ftproto.Feature{
				allFeaturesForPrerequisiteTest["featureD"],
				allFeaturesForPrerequisiteTest["featureH"],
			},
			expected: []*ftproto.Feature{
				allFeaturesForPrerequisiteTest["featureD"],
				allFeaturesForPrerequisiteTest["featureH"],
				allFeaturesForPrerequisiteTest["featureI"],
				allFeaturesForPrerequisiteTest["featureJ"],
				allFeaturesForPrerequisiteTest["featureK"],
			},
			expectedErr: nil,
		},
	}
	allFeatures := make([]*ftproto.Feature, 0, len(allFeaturesForPrerequisiteTest))
	for _, v := range allFeaturesForPrerequisiteTest {
		allFeatures = append(allFeatures, v)
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			actual, err := evaluator.GetPrerequisiteDownwards(p.target, allFeatures)
			assert.ElementsMatch(t, p.expected, actual)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetEvalFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	evaluator := NewEvaluator()

	patterns := []struct {
		desc        string
		targets     []*ftproto.Feature
		all         []*ftproto.Feature
		expectedIDs []string
	}{
		{
			desc: "success: No prerequisites",
			targets: []*ftproto.Feature{
				{Id: "featureA"},
			},
			all: []*ftproto.Feature{
				{Id: "featureA"},
				{Id: "featureB"},
			},
			expectedIDs: []string{"featureA"},
		},
		{
			desc: "success: one feature depends on target",
			targets: []*ftproto.Feature{
				{Id: "featureA"},
			},
			all: []*ftproto.Feature{
				{Id: "featureA"},
				{
					Id: "featureB",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureA"},
					},
				},
				{Id: "featureC"},
			},
			expectedIDs: []string{"featureA", "featureB"},
		},
		{
			desc: "success: multiple features depends on target",
			targets: []*ftproto.Feature{
				{Id: "featureA"},
			},
			all: []*ftproto.Feature{
				{Id: "featureA"},
				{
					Id: "featureB",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureA"},
					},
				},
				{
					Id: "featureC",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
					},
				},
				{
					Id: "featureD",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureA"},
					},
				},
				{Id: "featureE"},
			},
			expectedIDs: []string{"featureA", "featureB", "featureC", "featureD"},
		},
		{
			desc: "success: target depends on one feature",
			targets: []*ftproto.Feature{
				{
					Id: "featureA",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
					},
				},
			},
			all: []*ftproto.Feature{
				{
					Id: "featureA",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
					},
				},
				{
					Id: "featureB",
				},
				{Id: "featureC"},
			},
			expectedIDs: []string{"featureA", "featureB"},
		},
		{
			desc: "success: target depends on multiple features",
			targets: []*ftproto.Feature{
				{
					Id: "featureA",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
						{FeatureId: "featureC"},
					},
				},
			},
			all: []*ftproto.Feature{
				{
					Id: "featureA",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
					},
				},
				{
					Id: "featureB",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureD"},
					},
				},
				{Id: "featureC"},
				{Id: "featureD"},
				{Id: "featureE"},
			},
			expectedIDs: []string{"featureA", "featureB", "featureC", "featureD"},
		},
		{
			desc: "success: complex pattern 1",
			targets: []*ftproto.Feature{
				{
					Id: "featureD",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
					},
				},
			},
			all: []*ftproto.Feature{
				{
					Id: "featureA",
				},
				{
					Id: "featureB",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureA"},
					},
				},
				{
					Id: "featureC",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
					},
				},
				{
					Id: "featureD",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureB"},
					},
				},
				{
					Id: "featureE",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureC"},
						{FeatureId: "featureD"},
					},
				},
				{
					Id: "featureF",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureE"},
					},
				},
				{
					Id: "featureG",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "featureA"},
					},
				},
				{
					Id: "featureH",
				},
			},
			expectedIDs: []string{
				"featureA", "featureB", "featureC", "featureD", "featureE", "featureF",
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, _ := evaluator.getEvalFeatures(p.targets, p.all)
			assert.Equal(t, len(p.expectedIDs), len(actual))
			actualIDs := make([]string, 0, len(actual))
			for _, v := range actual {
				actualIDs = append(actualIDs, v.Id)
			}
			assert.ElementsMatch(t, p.expectedIDs, actualIDs)
		})
	}
}

func makeFeature(id string) *ftproto.Feature {
	return &ftproto.Feature{
		Id:            id,
		Name:          "test feature",
		Version:       1,
		Enabled:       true,
		CreatedAt:     time.Now().Unix(),
		VariationType: feature.Feature_STRING,
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
	}
}

func TestAssignUserOffVariation(t *testing.T) {
	t.Parallel()
	f := makeFeature("test-feature")
	evaluator := NewEvaluator()
	patterns := []struct {
		enabled           bool
		offVariation      string
		userID            string
		Flagvariations    map[string]string
		prerequisite      []*ftproto.Prerequisite
		expectedReason    *ftproto.Reason
		expectedVariation *ftproto.Variation
		expectedError     error
	}{
		{
			enabled:           false,
			offVariation:      "variation-C",
			userID:            "user5",
			Flagvariations:    map[string]string{},
			prerequisite:      []*ftproto.Prerequisite{},
			expectedReason:    &ftproto.Reason{Type: ftproto.Reason_OFF_VARIATION},
			expectedVariation: f.Variations[2],
			expectedError:     nil,
		},
		{
			enabled:           false,
			offVariation:      "",
			userID:            "user5",
			Flagvariations:    map[string]string{},
			prerequisite:      []*ftproto.Prerequisite{},
			expectedReason:    &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			expectedVariation: f.Variations[1],
			expectedError:     nil,
		},
		{
			enabled:           false,
			offVariation:      "variation-E",
			userID:            "user5",
			Flagvariations:    map[string]string{},
			prerequisite:      []*ftproto.Prerequisite{},
			expectedReason:    &ftproto.Reason{Type: ftproto.Reason_OFF_VARIATION},
			expectedVariation: nil,
			expectedError:     ErrVariationNotFound,
		},
		{
			enabled:           true,
			offVariation:      "",
			userID:            "user4",
			Flagvariations:    map[string]string{},
			prerequisite:      []*ftproto.Prerequisite{},
			expectedReason:    &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			expectedVariation: f.Variations[1],
			expectedError:     nil,
		},
		{
			enabled:           true,
			offVariation:      "variation-C",
			userID:            "user4",
			Flagvariations:    map[string]string{},
			prerequisite:      []*ftproto.Prerequisite{},
			expectedReason:    &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
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
			prerequisite: []*ftproto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expectedReason:    &ftproto.Reason{Type: ftproto.Reason_PREREQUISITE},
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
			prerequisite: []*ftproto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expectedReason:    &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			expectedVariation: f.Variations[1],
			expectedError:     nil,
		},
		{
			enabled:        true,
			offVariation:   "variation-C",
			userID:         "user4",
			Flagvariations: map[string]string{}, // not found prerequisite vatiation
			prerequisite: []*ftproto.Prerequisite{
				{
					FeatureId:   "test-feature2",
					VariationId: "variation D",
				},
			},
			expectedReason:    nil,
			expectedVariation: nil,
			expectedError:     ErrPrerequisiteVariationNotFound,
		},
	}
	for _, p := range patterns {
		user := &userproto.User{Id: p.userID}
		f.Enabled = p.enabled
		f.OffVariation = p.offVariation
		f.Prerequisites = p.prerequisite
		reason, variation, err := evaluator.assignUser(f, user, nil, p.Flagvariations)
		assert.Equal(t, p.expectedReason, reason)
		assert.Equal(t, p.expectedVariation, variation)
		assert.Equal(t, p.expectedError, err)
	}
}

func TestAssignUserTarget(t *testing.T) {
	f := makeFeature("test-feature")
	evaluator := NewEvaluator()
	patterns := []struct {
		userID              string
		expectedReason      ftproto.Reason_Type
		expectedVariationID string
	}{
		{
			userID:              "user1",
			expectedReason:      ftproto.Reason_TARGET,
			expectedVariationID: "variation-A",
		},
		{
			userID:              "user2",
			expectedReason:      ftproto.Reason_TARGET,
			expectedVariationID: "variation-B",
		},
		{
			userID:              "user3",
			expectedReason:      ftproto.Reason_TARGET,
			expectedVariationID: "variation-C",
		},
		{
			userID:              "user4",
			expectedReason:      ftproto.Reason_DEFAULT,
			expectedVariationID: "variation-B",
		},
	}
	for _, p := range patterns {
		user := &userproto.User{Id: p.userID}
		reason, variation, err := evaluator.assignUser(f, user, nil, nil)
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
	evaluator := NewEvaluator()
	reason, variation, err := evaluator.assignUser(f, user, nil, nil)
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

	evaluator := NewEvaluator()
	reason, variation, err := evaluator.assignUser(f, user, nil, nil)
	if reason != nil {
		t.Fatalf("Failed to assign user. Reason should be nil: %v", reason)
	}
	if variation != nil {
		t.Fatalf("Failed to assign user. Variation should be nil: %v", variation)
	}
	if err != ErrDefaultStrategyNotFound {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
}

func TestAssignUserDefaultStrategy(t *testing.T) {
	user := &userproto.User{
		Id:   "user-id1",
		Data: map[string]string{"name3": "user3"},
	}
	f := makeFeature("test-feature")
	evaluator := NewEvaluator()
	reason, variation, err := evaluator.assignUser(f, user, nil, nil)
	if err != nil {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
	if reason.Type != ftproto.Reason_DEFAULT {
		t.Fatalf("Failed to assign user. Reason type does not match. Current: %s, target: %v", reason.Type, ftproto.Reason_DEFAULT)
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
	f.DefaultStrategy = &ftproto.Strategy{
		Type: ftproto.Strategy_ROLLOUT,
		RolloutStrategy: &ftproto.RolloutStrategy{
			Variations: []*ftproto.RolloutStrategy_Variation{
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
	evaluator := NewEvaluator()
	reason, variation, err := evaluator.assignUser(f, user, nil, nil)
	if err != nil {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
	if reason.Type != ftproto.Reason_DEFAULT {
		t.Fatalf("Failed to assign user. Reason type does not match. Current: %s, target: %v", reason.Type, ftproto.Reason_DEFAULT)
	}
	if variation.Id != f.DefaultStrategy.RolloutStrategy.Variations[2].Variation {
		t.Fatalf("Failed to assign user. Variation id does not match. Current: %s, target: %s", variation.Id, f.DefaultStrategy.RolloutStrategy.Variations[2].Variation)
	}
	// Channge sampling seed to change assigned variation.
	f.SamplingSeed = "sampling-seed"
	reason, variation, err = evaluator.assignUser(f, user, nil, nil)
	if err != nil {
		t.Fatalf("Failed to assign user. Error: %v", err)
	}
	if reason.Type != ftproto.Reason_DEFAULT {
		t.Fatalf("Failed to assign user. Reason type does not match. Current: %s, target: %v", reason.Type, ftproto.Reason_DEFAULT)
	}
	if variation.Id != f.DefaultStrategy.RolloutStrategy.Variations[0].Variation {
		t.Fatalf("Failed to assign user. Variation id does not match. Current: %s, target: %s", variation.Id, f.DefaultStrategy.RolloutStrategy.Variations[0].Variation)
	}
}

func TestEvaluateFeaturesByEvaluatedAt_MissingPrerequisiteActual(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc                  string
		setupFunc             func() ([]*ftproto.Feature, *userproto.User)
		prevUEID              string
		evaluatedAt           int64
		userAttributesUpdated bool
		tag                   string
		expectedErr           string // Expected error substring, empty if no error expected
	}{
		{
			desc: "success: incremental evaluation with old prerequisites should not fail with 'feature not found'",
			setupFunc: func() ([]*ftproto.Feature, *userproto.User) {
				// Test the REAL production scenario:
				// - Main feature was updated recently
				// - Prerequisites were NOT updated recently
				// - ALL features are present in the input
				mainFeature := makeDependentFeature()
				mainFeature.UpdatedAt = time.Now().Unix() - 30 // Recently updated

				prereq1 := makeTestPrereqA()
				prereq1.UpdatedAt = time.Now().Unix() - 7200 // Updated 2 hours ago (old)

				prereq2 := makeTestPrereqB()
				prereq2.UpdatedAt = time.Now().Unix() - 7200 // Updated 2 hours ago (old)

				// Include ALL features - this simulates what the API layer now passes
				features := []*ftproto.Feature{mainFeature, prereq1, prereq2}
				user := &userproto.User{Id: "user-1"}

				return features, user
			},
			prevUEID:              "prev-ueid",
			evaluatedAt:           time.Now().Unix() - 60, // 1 minute ago
			userAttributesUpdated: false,
			tag:                   "test",
			expectedErr:           "", // Should not contain "feature not found"
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			features, user := p.setupFunc()
			segmentUsersMap := map[string][]*ftproto.SegmentUser{}

			result, err := evaluator.EvaluateFeaturesByEvaluatedAt(
				features,
				user,
				segmentUsersMap,
				p.prevUEID,
				p.evaluatedAt,
				p.userAttributesUpdated,
				p.tag,
			)

			// After our fix, dependency resolution should work and evaluation should succeed
			assert.NoError(t, err, "Evaluation should succeed with proper dependency resolution")
			assert.NotNil(t, result, "Result should not be nil")
			assert.NotNil(t, result.Evaluations, "Evaluations should not be nil")
		})
	}
}

func TestGetEvalFeatures_IncrementalEvaluationTransitiveDependencies(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		setupFunc   func() ([]*ftproto.Feature, []*ftproto.Feature)
		expectedIDs []string
		expectedErr error
	}{
		{
			desc: "success: transitive dependency resolution in incremental evaluation",
			setupFunc: func() ([]*ftproto.Feature, []*ftproto.Feature) {
				// Test scenario: dependent feature depends on BOTH prerequisites
				mainFeature := &ftproto.Feature{
					Id: "test-dependent-feature",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "test-prereq-a", VariationId: "var1"},
						{FeatureId: "test-prereq-b", VariationId: "var2"},
					},
				}

				prereq1 := &ftproto.Feature{
					Id: "test-prereq-a",
				}

				prereq2 := &ftproto.Feature{
					Id: "test-prereq-b",
				}

				// Simulate incremental evaluation: test-prereq-a was updated recently
				targets := []*ftproto.Feature{prereq1} // Only this one is "updated"
				allFeatures := []*ftproto.Feature{mainFeature, prereq1, prereq2}

				return targets, allFeatures
			},
			expectedIDs: []string{"test-prereq-a", "test-dependent-feature", "test-prereq-b"},
			expectedErr: nil,
		},
		{
			desc: "success: handles deep dependency chains within iteration limit",
			setupFunc: func() ([]*ftproto.Feature, []*ftproto.Feature) {
				// Create a chain of 10 features to test iteration limit
				features := make([]*ftproto.Feature, 10)
				for i := 0; i < 10; i++ {
					id := fmt.Sprintf("chain-feature-%d", i)
					features[i] = &ftproto.Feature{Id: id}

					// Each feature depends on the next one (except the last)
					if i < 9 {
						features[i].Prerequisites = []*ftproto.Prerequisite{
							{FeatureId: fmt.Sprintf("chain-feature-%d", i+1), VariationId: "var1"},
						}
					}
				}

				// Target is the first feature in the chain
				targets := []*ftproto.Feature{features[0]}

				return targets, features
			},
			expectedIDs: []string{
				"chain-feature-0", "chain-feature-1", "chain-feature-2", "chain-feature-3", "chain-feature-4",
				"chain-feature-5", "chain-feature-6", "chain-feature-7", "chain-feature-8", "chain-feature-9",
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			targets, allFeatures := p.setupFunc()

			result, err := evaluator.getEvalFeatures(targets, allFeatures)

			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				// Should include ALL 3 features:
				// 1. test-prereq-a (target)
				// 2. test-dependent-feature (depends on target)
				// 3. test-prereq-b (transitive dependency of #2)
				assert.Len(t, result, len(p.expectedIDs))

				resultIDs := make([]string, len(result))
				for i, f := range result {
					resultIDs[i] = f.Id
				}

				assert.ElementsMatch(t, p.expectedIDs, resultIDs)
			}
		})
	}
}

func TestGetEvalFeatures_FeatureFlagRuleDependencies(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		setupFunc   func() ([]*ftproto.Feature, []*ftproto.Feature)
		expectedIDs []string
		expectedErr error
	}{
		{
			desc: "success: FEATURE_FLAG rule dependency resolution with transitive prerequisites",
			setupFunc: func() ([]*ftproto.Feature, []*ftproto.Feature) {
				// Feature with FEATURE_FLAG rule dependency
				mainFeature := &ftproto.Feature{
					Id: "feature-with-rule",
					Rules: []*ftproto.Rule{
						{
							Clauses: []*ftproto.Clause{
								{
									Operator:  feature.Clause_FEATURE_FLAG,
									Attribute: "dependency-flag", // References another feature
									Values:    []string{"true"},
								},
							},
						},
					},
				}

				// The dependency has its own prerequisites
				dependencyFlag := &ftproto.Feature{
					Id: "dependency-flag",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "deep-dependency", VariationId: "var1"},
					},
				}

				deepDependency := &ftproto.Feature{
					Id: "deep-dependency",
				}

				// Simulate: dependency-flag was updated recently
				targets := []*ftproto.Feature{dependencyFlag}
				allFeatures := []*ftproto.Feature{mainFeature, dependencyFlag, deepDependency}

				return targets, allFeatures
			},
			expectedIDs: []string{"dependency-flag", "feature-with-rule", "deep-dependency"},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			targets, allFeatures := p.setupFunc()

			result, err := evaluator.getEvalFeatures(targets, allFeatures)

			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				resultIDs := make([]string, len(result))
				for i, f := range result {
					resultIDs[i] = f.Id
				}

				assert.ElementsMatch(t, p.expectedIDs, resultIDs)
			}
		})
	}
}

func TestGetFeaturesDependedOnTargets_MissingDependency(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		setupFunc   func() ([]*ftproto.Feature, []*ftproto.Feature)
		expectedLen int
		expectedIDs []string
		expectedErr error
	}{
		{
			desc: "success: graceful handling of missing dependencies",
			setupFunc: func() ([]*ftproto.Feature, []*ftproto.Feature) {
				// Feature that depends on missing prerequisite (simulates data corruption)
				mainFeature := &ftproto.Feature{
					Id: "main-feature",
					Prerequisites: []*ftproto.Prerequisite{
						{FeatureId: "missing-prereq", VariationId: "variation-1"},
					},
				}

				// The prerequisite is missing from allFeatures (simulates cache miss/corruption)
				targets := []*ftproto.Feature{mainFeature}
				allFeatures := []*ftproto.Feature{mainFeature}

				return targets, allFeatures
			},
			expectedLen: 1,
			expectedIDs: []string{"main-feature"},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			targets, allFeatures := p.setupFunc()

			evalFeatures, err := evaluator.getEvalFeatures(targets, allFeatures)

			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				// Should handle gracefully (not panic) due to our nil pointer fix
				assert.Len(t, evalFeatures, p.expectedLen)

				resultIDs := make([]string, len(evalFeatures))
				for i, f := range evalFeatures {
					resultIDs[i] = f.Id
				}

				assert.ElementsMatch(t, p.expectedIDs, resultIDs)
			}
		})
	}
}

// makeDependentFeature creates a feature that requires prerequisites for testing
func makeDependentFeature() *ftproto.Feature {
	return &ftproto.Feature{
		Id:            "test-dependent-feature",
		Name:          "Test Feature with Dependencies",
		Version:       1,
		Enabled:       true, // Enable the feature so it can be properly evaluated
		Archived:      false,
		CreatedAt:     1700000000,
		UpdatedAt:     1700000100,
		Tags:          []string{"test"},
		VariationType: feature.Feature_BOOLEAN,
		OffVariation:  "variation-false",
		Variations: []*ftproto.Variation{
			{
				Id:    "variation-true",
				Name:  "On",
				Value: "true",
			},
			{
				Id:    "variation-false",
				Name:  "Off",
				Value: "false",
			},
		},
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId:   "test-prereq-a",
				VariationId: "variation-true",
			},
			{
				FeatureId:   "test-prereq-b",
				VariationId: "variation-true",
			},
		},
		Rules: []*ftproto.Rule{},
		Targets: []*ftproto.Target{
			{
				Variation: "variation-true",
				Users:     []string{},
			},
			{
				Variation: "variation-false",
				Users:     []string{},
			},
		},
		DefaultStrategy: &ftproto.Strategy{
			Type: ftproto.Strategy_FIXED,
			FixedStrategy: &ftproto.FixedStrategy{
				Variation: "variation-false",
			},
		},
	}
}

// makeTestPrereqA creates the first prerequisite for testing
func makeTestPrereqA() *ftproto.Feature {
	return &ftproto.Feature{
		Id:            "test-prereq-a",
		Name:          "Test Prerequisite A",
		Version:       1,
		Enabled:       true,
		Archived:      false,
		CreatedAt:     1700000000,
		UpdatedAt:     1700000200,
		Tags:          []string{"test"},
		VariationType: feature.Feature_BOOLEAN,
		OffVariation:  "variation-false",
		Variations: []*ftproto.Variation{
			{
				Id:    "variation-true",
				Name:  "On",
				Value: "true",
			},
			{
				Id:    "variation-false",
				Name:  "Off",
				Value: "false",
			},
		},
		Prerequisites: []*ftproto.Prerequisite{},
		Rules: []*ftproto.Rule{
			{
				Id: "rule-1",
				Clauses: []*ftproto.Clause{
					{
						Id:        "clause-1",
						Values:    []string{"1.0.0"},
						Operator:  ftproto.Clause_GREATER_OR_EQUAL,
						Attribute: "app_version",
					},
				},
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{Weight: 50000, Variation: "variation-true"},
							{Weight: 50000, Variation: "variation-false"},
						},
					},
				},
			},
		},
		Targets: []*ftproto.Target{
			{
				Variation: "variation-true",
				Users:     []string{"test-user-1", "test-user-2"},
			},
			{
				Variation: "variation-false",
				Users:     []string{},
			},
		},
		DefaultStrategy: &ftproto.Strategy{
			Type: ftproto.Strategy_FIXED,
			FixedStrategy: &ftproto.FixedStrategy{
				Variation: "variation-true",
			},
		},
	}
}

// makeTestPrereqB creates the second prerequisite for testing
func makeTestPrereqB() *ftproto.Feature {
	return &ftproto.Feature{
		Id:            "test-prereq-b",
		Name:          "Test Prerequisite B",
		Version:       1,
		Enabled:       true,
		Archived:      false,
		CreatedAt:     1700000000,
		UpdatedAt:     1700000300,
		Tags:          []string{"test"},
		VariationType: feature.Feature_BOOLEAN,
		OffVariation:  "variation-false",
		Variations: []*ftproto.Variation{
			{
				Id:    "variation-true",
				Name:  "On",
				Value: "true",
			},
			{
				Id:    "variation-false",
				Name:  "Off",
				Value: "false",
			},
		},
		Prerequisites: []*ftproto.Prerequisite{},
		Rules: []*ftproto.Rule{
			{
				Id: "rule-2",
				Clauses: []*ftproto.Clause{
					{
						Id:        "clause-2",
						Values:    []string{"2.0.0"},
						Operator:  ftproto.Clause_GREATER_OR_EQUAL,
						Attribute: "app_version",
					},
				},
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_ROLLOUT,
					RolloutStrategy: &ftproto.RolloutStrategy{
						Variations: []*ftproto.RolloutStrategy_Variation{
							{Weight: 25000, Variation: "variation-true"},
							{Weight: 75000, Variation: "variation-false"},
						},
					},
				},
			},
		},
		Targets: []*ftproto.Target{
			{
				Variation: "variation-true",
				Users:     []string{"test-user-3", "test-user-4"},
			},
			{
				Variation: "variation-false",
				Users:     []string{},
			},
		},
		DefaultStrategy: &ftproto.Strategy{
			Type: ftproto.Strategy_FIXED,
			FixedStrategy: &ftproto.FixedStrategy{
				Variation: "variation-true",
			},
		},
	}
}

func TestEvaluateFeaturesByEvaluatedAt_TagMismatchScenario(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc                  string
		setupFunc             func() ([]*ftproto.Feature, *userproto.User)
		prevUEID              string
		evaluatedAt           int64
		userAttributesUpdated bool
		tag                   string
		expectedErr           string // Expected error substring, empty if no error expected
	}{
		{
			desc: "success: tag mismatch should not cause 'feature not found' errors",
			setupFunc: func() ([]*ftproto.Feature, *userproto.User) {
				// Test scenario: main feature has "test" tag, but prerequisite doesn't
				mainFeature := makeDependentFeature()
				mainFeature.UpdatedAt = time.Now().Unix() - 30 // Recently updated
				mainFeature.Tags = []string{"test"}

				prereq1 := makeTestPrereqA()
				prereq1.UpdatedAt = time.Now().Unix() - 3600 // Updated 1 hour ago (old)
				prereq1.Tags = []string{"mobile"}            // DIFFERENT TAG!

				prereq2 := makeTestPrereqB()
				prereq2.UpdatedAt = time.Now().Unix() - 3600 // Updated 1 hour ago (old)
				prereq2.Tags = []string{"test"}

				// Include all features but with tag mismatch
				features := []*ftproto.Feature{mainFeature, prereq1, prereq2}
				user := &userproto.User{Id: "user-1"}

				return features, user
			},
			prevUEID:              "prev-ueid",
			evaluatedAt:           time.Now().Unix() - 60, // 1 minute ago
			userAttributesUpdated: false,
			tag:                   "test", // Requesting test features only
			expectedErr:           "",     // Should not contain "feature not found"
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			features, user := p.setupFunc()
			segmentUsersMap := map[string][]*ftproto.SegmentUser{}

			result, err := evaluator.EvaluateFeaturesByEvaluatedAt(
				features,
				user,
				segmentUsersMap,
				p.prevUEID,
				p.evaluatedAt,
				p.userAttributesUpdated,
				p.tag,
			)

			// After our fix, evaluation should succeed even with tag mismatches in prerequisites
			assert.NoError(t, err, "Evaluation should succeed despite tag mismatches")
			assert.NotNil(t, result, "Result should not be nil")
			assert.NotNil(t, result.Evaluations, "Evaluations should not be nil")
		})
	}
}
