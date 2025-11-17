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
	"encoding/json"
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

	// Add YAML feature test case
	yamlFeature := makeFeature("yaml-feature-id")
	yamlFeature.Tags = append(yamlFeature.Tags, "tag1")
	yamlFeature.VariationType = ftproto.Feature_YAML
	yamlFeature.Variations = []*ftproto.Variation{
		{
			Id:   "yaml-variation-A",
			Name: "YAML Variation A",
			Value: `# Configuration A
config:
  enabled: true  # Enable feature
  timeout: 30    # Timeout in seconds`,
		},
		{
			Id:   "yaml-variation-B",
			Name: "YAML Variation B",
			Value: `# Configuration B
config:
  # Feature toggle
  enabled: false
  timeout: 60  # Longer timeout`,
		},
	}
	yamlFeature.DefaultStrategy = &ftproto.Strategy{
		Type: ftproto.Strategy_FIXED,
		FixedStrategy: &ftproto.FixedStrategy{
			Variation: "yaml-variation-B",
		},
	}

	yamlPatterns := []struct {
		enabled       bool
		offVariation  string
		userID        string
		prerequisite  []*ftproto.Prerequisite
		expected      *ftproto.Evaluation
		expectedError error
	}{
		{
			enabled:      true,
			offVariation: "",
			userID:       "yaml-user-1",
			prerequisite: []*ftproto.Prerequisite{},
			expected: &ftproto.Evaluation{
				Id:             EvaluationID(yamlFeature.Id, yamlFeature.Version, "yaml-user-1"),
				FeatureId:      "yaml-feature-id",
				FeatureVersion: 1,
				UserId:         "yaml-user-1",
				VariationId:    "yaml-variation-B",
				VariationName:  "YAML Variation B",
				VariationValue: `{"config":{"enabled":false,"timeout":60}}`,
				Variation: &ftproto.Variation{
					Id:    "yaml-variation-B",
					Name:  "YAML Variation B",
					Value: `{"config":{"enabled":false,"timeout":60}}`,
				},
				Reason: &ftproto.Reason{Type: ftproto.Reason_DEFAULT},
			},
			expectedError: nil,
		},
	}

	// Test YAML feature
	for _, p := range yamlPatterns {
		evaluator := NewEvaluator()
		user := &userproto.User{Id: p.userID}
		yamlFeature.Enabled = p.enabled
		yamlFeature.OffVariation = p.offVariation
		yamlFeature.Prerequisites = p.prerequisite
		segmentUser := map[string][]*ftproto.SegmentUser{}
		evaluation, err := evaluator.EvaluateFeatures([]*ftproto.Feature{yamlFeature}, user, segmentUser, "tag1")
		assert.Equal(t, p.expectedError, err)
		if evaluation != nil {
			actual, err := findEvaluation(evaluation.Evaluations, yamlFeature.Id)
			assert.NoError(t, err)
			assert.Equal(t, p.expected.VariationValue, actual.VariationValue, "YAML should be converted to JSON")
			assert.Equal(t, p.expected.Variation.Value, actual.Variation.Value, "Variation.Value should also be converted")
			// Verify it's valid JSON
			var jsonData interface{}
			jsonErr := json.Unmarshal([]byte(actual.VariationValue), &jsonData)
			assert.NoError(t, jsonErr, "YAML should be converted to valid JSON")
		}
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

func TestConvertVariationValue(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		variationType  ftproto.Feature_VariationType
		variationValue string
		expectedValue  string
	}{
		{
			desc:           "Non-YAML type returns original value",
			variationType:  ftproto.Feature_STRING,
			variationValue: "simple string",
			expectedValue:  "simple string",
		},
		{
			desc:           "JSON type returns original value",
			variationType:  ftproto.Feature_JSON,
			variationValue: `{"key": "value"}`,
			expectedValue:  `{"key": "value"}`,
		},
		{
			desc:          "YAML type converts to JSON",
			variationType: ftproto.Feature_YAML,
			variationValue: `name: John Doe
age: 30
active: true`,
			expectedValue: `{"active":true,"age":30,"name":"John Doe"}`,
		},
		{
			desc:          "YAML with nested objects converts to JSON",
			variationType: ftproto.Feature_YAML,
			variationValue: `user:
  name: Jane
  email: jane@example.com
settings:
  theme: dark
  notifications: true`,
			expectedValue: `{"settings":{"notifications":true,"theme":"dark"},"user":{"email":"jane@example.com","name":"Jane"}}`,
		},
		{
			desc:          "YAML with arrays converts to JSON",
			variationType: ftproto.Feature_YAML,
			variationValue: `items:
  - id: 1
    name: Item 1
  - id: 2
    name: Item 2`,
			expectedValue: `{"items":[{"id":1,"name":"Item 1"},{"id":2,"name":"Item 2"}]}`,
		},
		{
			desc:          "YAML with comments converts to JSON",
			variationType: ftproto.Feature_YAML,
			variationValue: `# This is a configuration
name: John Doe
# Age in years
age: 30
active: true # User is active`,
			expectedValue: `{"active":true,"age":30,"name":"John Doe"}`,
		},
		{
			desc:          "YAML with comments and nested objects converts to JSON",
			variationType: ftproto.Feature_YAML,
			variationValue: `# Database configuration
database:
  # Connection settings
  host: localhost
  port: 5432
  # Security
  ssl: true
# Application settings
app:
  debug: false # Disable in production`,
			expectedValue: `{"app":{"debug":false},"database":{"host":"localhost","port":5432,"ssl":true}}`,
		},
		{
			desc:           "Invalid YAML returns original value as fallback",
			variationType:  ftproto.Feature_YAML,
			variationValue: "invalid: yaml: [unclosed",
			expectedValue:  "invalid: yaml: [unclosed",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			feature := &ftproto.Feature{
				Id:            "feature-1",
				VariationType: p.variationType,
			}
			variation := &ftproto.Variation{
				Id:    "var-1",
				Value: p.variationValue,
			}

			result := evaluator.convertVariationValue(feature, variation)
			assert.Equal(t, p.expectedValue, result)
		})
	}
}

func TestConvertVariationValueCaching(t *testing.T) {
	t.Parallel()

	t.Run("Caches YAML to JSON conversion", func(t *testing.T) {
		evaluator := NewEvaluator()
		feature := &ftproto.Feature{
			Id:            "feature-1",
			VariationType: ftproto.Feature_YAML,
			UpdatedAt:     1234567890,
		}
		variation := &ftproto.Variation{
			Id: "var-1",
			Value: `name: Test
value: 123`,
		}

		// First call should convert and cache
		result1 := evaluator.convertVariationValue(feature, variation)
		assert.Equal(t, `{"name":"Test","value":123}`, result1)

		// Second call should use cache (verify by checking cache directly with correct key)
		cacheKey := fmt.Sprintf("%d:%s", feature.UpdatedAt, variation.Id)
		cached, ok := evaluator.variationCache.Load(cacheKey)
		assert.True(t, ok)
		assert.Equal(t, result1, cached)

		// Third call should return cached value
		result2 := evaluator.convertVariationValue(feature, variation)
		assert.Equal(t, result1, result2)
	})

	t.Run("Cache is keyed by UpdatedAt and variation ID", func(t *testing.T) {
		evaluator := NewEvaluator()
		feature := &ftproto.Feature{
			Id:            "feature-1",
			VariationType: ftproto.Feature_YAML,
			UpdatedAt:     1234567890,
		}

		variation1 := &ftproto.Variation{
			Id:    "var-1",
			Value: "key1: value1",
		}
		variation2 := &ftproto.Variation{
			Id:    "var-2",
			Value: "key2: value2",
		}

		result1 := evaluator.convertVariationValue(feature, variation1)
		result2 := evaluator.convertVariationValue(feature, variation2)

		// Different variations should have different results
		assert.NotEqual(t, result1, result2)
		assert.Equal(t, `{"key1":"value1"}`, result1)
		assert.Equal(t, `{"key2":"value2"}`, result2)

		// Both should be cached separately with correct keys
		cacheKey1 := fmt.Sprintf("%d:%s", feature.UpdatedAt, variation1.Id)
		cacheKey2 := fmt.Sprintf("%d:%s", feature.UpdatedAt, variation2.Id)
		cached1, ok1 := evaluator.variationCache.Load(cacheKey1)
		cached2, ok2 := evaluator.variationCache.Load(cacheKey2)
		assert.True(t, ok1)
		assert.True(t, ok2)
		assert.Equal(t, result1, cached1)
		assert.Equal(t, result2, cached2)
	})

	t.Run("Cache invalidates when feature is updated", func(t *testing.T) {
		evaluator := NewEvaluator()
		feature := &ftproto.Feature{
			Id:            "feature-1",
			VariationType: ftproto.Feature_YAML,
			UpdatedAt:     1234567890,
		}
		variation := &ftproto.Variation{
			Id:    "var-1",
			Value: "key: original_value",
		}

		// First call with original timestamp
		result1 := evaluator.convertVariationValue(feature, variation)
		assert.Equal(t, `{"key":"original_value"}`, result1)

		// Verify cache with original key
		cacheKey1 := fmt.Sprintf("%d:%s", feature.UpdatedAt, variation.Id)
		cached1, ok1 := evaluator.variationCache.Load(cacheKey1)
		assert.True(t, ok1)
		assert.Equal(t, result1, cached1)

		// Update feature timestamp (simulating feature update)
		feature.UpdatedAt = 1234567999
		variation.Value = "key: updated_value"

		// Second call with updated timestamp should create new cache entry
		result2 := evaluator.convertVariationValue(feature, variation)
		assert.Equal(t, `{"key":"updated_value"}`, result2)

		// Verify new cache key exists
		cacheKey2 := fmt.Sprintf("%d:%s", feature.UpdatedAt, variation.Id)
		cached2, ok2 := evaluator.variationCache.Load(cacheKey2)
		assert.True(t, ok2)
		assert.Equal(t, result2, cached2)

		// Results should be different
		assert.NotEqual(t, result1, result2)

		// Old cache entry still exists (no automatic cleanup)
		_, stillExists := evaluator.variationCache.Load(cacheKey1)
		assert.True(t, stillExists)
	})

	t.Run("Does not cache non-YAML types", func(t *testing.T) {
		evaluator := NewEvaluator()
		feature := &ftproto.Feature{
			Id:            "feature-1",
			VariationType: ftproto.Feature_STRING,
			UpdatedAt:     1234567890,
		}
		variation := &ftproto.Variation{
			Id:    "var-1",
			Value: "simple string",
		}

		result := evaluator.convertVariationValue(feature, variation)
		assert.Equal(t, "simple string", result)

		// Should not be in cache
		cacheKey := fmt.Sprintf("%d:%s", feature.UpdatedAt, variation.Id)
		_, ok := evaluator.variationCache.Load(cacheKey)
		assert.False(t, ok)
	})
}

func TestYAMLToJSON(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		yamlInput   string
		expected    string
		expectedErr bool
	}{
		{
			desc:        "Simple YAML",
			yamlInput:   "key: value",
			expected:    `{"key":"value"}`,
			expectedErr: false,
		},
		{
			desc: "Nested YAML",
			yamlInput: `parent:
  child: value
  number: 42`,
			expected:    `{"parent":{"child":"value","number":42}}`,
			expectedErr: false,
		},
		{
			desc: "YAML with array",
			yamlInput: `list:
  - item1
  - item2
  - item3`,
			expected:    `{"list":["item1","item2","item3"]}`,
			expectedErr: false,
		},
		{
			desc: "YAML with mixed types",
			yamlInput: `string: text
number: 123
float: 45.67
boolean: true
nullValue: null`,
			expected:    `{"boolean":true,"float":45.67,"nullValue":null,"number":123,"string":"text"}`,
			expectedErr: false,
		},
		{
			desc: "YAML with comments",
			yamlInput: `# Configuration file
key: value
# Number setting
count: 42`,
			expected:    `{"count":42,"key":"value"}`,
			expectedErr: false,
		},
		{
			desc:        "Invalid YAML",
			yamlInput:   "invalid: yaml: [unclosed",
			expected:    "",
			expectedErr: true,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			result, err := yamlToJSON(p.yamlInput)

			if p.expectedErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to convert YAML to JSON")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, p.expected, result)
			}
		})
	}
}

func TestEvaluateWithYAMLVariation(t *testing.T) {
	t.Parallel()

	t.Run("Evaluates YAML variation and converts to JSON", func(t *testing.T) {
		evaluator := NewEvaluator()

		// Create a feature with YAML variation type
		feature := &ftproto.Feature{
			Id:            "yaml-feature",
			Name:          "YAML Feature",
			Version:       1,
			Enabled:       true,
			VariationType: ftproto.Feature_YAML,
			Variations: []*ftproto.Variation{
				{
					Id:   "yaml-var-1",
					Name: "YAML Variation",
					Value: `config:
  enabled: true
  maxRetries: 3
  timeout: 30`,
				},
			},
			DefaultStrategy: &ftproto.Strategy{
				Type: ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{
					Variation: "yaml-var-1",
				},
			},
		}

		user := &userproto.User{Id: "user-1"}
		result, err := evaluator.EvaluateFeatures(
			[]*ftproto.Feature{feature},
			user,
			map[string][]*ftproto.SegmentUser{},
			"",
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Evaluations, 1)

		evaluation := result.Evaluations[0]
		// Verify the value is converted to JSON
		expectedJSON := `{"config":{"enabled":true,"maxRetries":3,"timeout":30}}`
		assert.Equal(t, expectedJSON, evaluation.VariationValue)
		assert.Equal(t, expectedJSON, evaluation.Variation.Value)
	})

	t.Run("Multiple evaluations with same YAML variation use cache", func(t *testing.T) {
		evaluator := NewEvaluator()

		yamlValue := `settings:
  theme: dark
  language: en`

		feature := &ftproto.Feature{
			Id:            "yaml-feature",
			Name:          "YAML Feature",
			Version:       1,
			Enabled:       true,
			CreatedAt:     1234567890,
			UpdatedAt:     1234567890,
			VariationType: ftproto.Feature_YAML,
			Variations: []*ftproto.Variation{
				{
					Id:    "yaml-var-shared",
					Name:  "Shared YAML Variation",
					Value: yamlValue,
				},
			},
			DefaultStrategy: &ftproto.Strategy{
				Type: ftproto.Strategy_FIXED,
				FixedStrategy: &ftproto.FixedStrategy{
					Variation: "yaml-var-shared",
				},
			},
		}

		// Evaluate for first user
		user1 := &userproto.User{Id: "user-1"}
		result1, err1 := evaluator.EvaluateFeatures(
			[]*ftproto.Feature{feature},
			user1,
			map[string][]*ftproto.SegmentUser{},
			"",
		)
		require.NoError(t, err1)
		require.Len(t, result1.Evaluations, 1)

		// Evaluate for second user
		user2 := &userproto.User{Id: "user-2"}
		result2, err2 := evaluator.EvaluateFeatures(
			[]*ftproto.Feature{feature},
			user2,
			map[string][]*ftproto.SegmentUser{},
			"",
		)
		require.NoError(t, err2)
		require.Len(t, result2.Evaluations, 1)

		// Both should have the same converted JSON value
		expectedJSON := `{"settings":{"language":"en","theme":"dark"}}`
		assert.Equal(t, expectedJSON, result1.Evaluations[0].VariationValue)
		assert.Equal(t, expectedJSON, result2.Evaluations[0].VariationValue)

		// Verify cache was used (with correct key format: updatedAt:variationId)
		cacheKey := fmt.Sprintf("%d:%s", feature.UpdatedAt, "yaml-var-shared")
		cached, ok := evaluator.variationCache.Load(cacheKey)
		assert.True(t, ok)
		assert.Equal(t, expectedJSON, cached)
	})
}

func TestEvaluate_YAMLConversion(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc                string
		setupFunc           func() ([]*ftproto.Feature, *userproto.User)
		expectedEvalCount   int
		expectedVariationID string
		validateValue       func(t *testing.T, value string)
	}{
		{
			desc: "Evaluate function converts YAML to JSON for single feature",
			setupFunc: func() ([]*ftproto.Feature, *userproto.User) {
				feature := &ftproto.Feature{
					Id:            "yaml-feature-1",
					Name:          "YAML Feature 1",
					Version:       1,
					Enabled:       true,
					VariationType: ftproto.Feature_YAML,
					Variations: []*ftproto.Variation{
						{
							Id:   "yaml-var-a",
							Name: "Config A",
							Value: `# Application config
app:
  name: MyApp
  version: 1.0.0
  # Feature flags
  features:
    - login
    - signup`,
						},
						{
							Id:   "yaml-var-b",
							Name: "Config B",
							Value: `app:
  name: MyApp
  version: 2.0.0`,
						},
					},
					DefaultStrategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "yaml-var-a",
						},
					},
				}
				user := &userproto.User{Id: "test-user-1"}
				return []*ftproto.Feature{feature}, user
			},
			expectedEvalCount:   1,
			expectedVariationID: "yaml-var-a",
			validateValue: func(t *testing.T, value string) {
				// Verify it's valid JSON
				var jsonData map[string]interface{}
				err := json.Unmarshal([]byte(value), &jsonData)
				assert.NoError(t, err, "Should be valid JSON")

				// Verify structure
				assert.Contains(t, jsonData, "app")
				app := jsonData["app"].(map[string]interface{})
				assert.Equal(t, "MyApp", app["name"])
				assert.Equal(t, "1.0.0", app["version"])
				assert.Contains(t, app, "features")
			},
		},
		{
			desc: "Evaluate function with multiple YAML features",
			setupFunc: func() ([]*ftproto.Feature, *userproto.User) {
				feature1 := &ftproto.Feature{
					Id:            "yaml-feature-1",
					Name:          "YAML Feature 1",
					Version:       1,
					Enabled:       true,
					VariationType: ftproto.Feature_YAML,
					Variations: []*ftproto.Variation{
						{
							Id:   "yaml-var-1",
							Name: "Config 1",
							Value: `database:
  host: localhost
  port: 5432`,
						},
					},
					DefaultStrategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "yaml-var-1",
						},
					},
				}
				feature2 := &ftproto.Feature{
					Id:            "yaml-feature-2",
					Name:          "YAML Feature 2",
					Version:       1,
					Enabled:       true,
					VariationType: ftproto.Feature_YAML,
					Variations: []*ftproto.Variation{
						{
							Id:   "yaml-var-2",
							Name: "Config 2",
							Value: `cache:
  enabled: true
  ttl: 3600`,
						},
					},
					DefaultStrategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "yaml-var-2",
						},
					},
				}
				user := &userproto.User{Id: "test-user-2"}
				return []*ftproto.Feature{feature1, feature2}, user
			},
			expectedEvalCount: 2,
			validateValue: func(t *testing.T, value string) {
				var jsonData map[string]interface{}
				err := json.Unmarshal([]byte(value), &jsonData)
				assert.NoError(t, err, "Should be valid JSON")
			},
		},
		{
			desc: "Evaluate function with mixed variation types (YAML and non-YAML)",
			setupFunc: func() ([]*ftproto.Feature, *userproto.User) {
				yamlFeature := &ftproto.Feature{
					Id:            "yaml-feature",
					Name:          "YAML Feature",
					Version:       1,
					Enabled:       true,
					VariationType: ftproto.Feature_YAML,
					Variations: []*ftproto.Variation{
						{
							Id:   "yaml-var",
							Name: "YAML Config",
							Value: `enabled: true
timeout: 30`,
						},
					},
					DefaultStrategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "yaml-var",
						},
					},
				}
				stringFeature := &ftproto.Feature{
					Id:            "string-feature",
					Name:          "String Feature",
					Version:       1,
					Enabled:       true,
					VariationType: ftproto.Feature_STRING,
					Variations: []*ftproto.Variation{
						{
							Id:    "string-var",
							Name:  "String Value",
							Value: "simple-string",
						},
					},
					DefaultStrategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "string-var",
						},
					},
				}
				jsonFeature := &ftproto.Feature{
					Id:            "json-feature",
					Name:          "JSON Feature",
					Version:       1,
					Enabled:       true,
					VariationType: ftproto.Feature_JSON,
					Variations: []*ftproto.Variation{
						{
							Id:    "json-var",
							Name:  "JSON Value",
							Value: `{"key":"value"}`,
						},
					},
					DefaultStrategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "json-var",
						},
					},
				}
				user := &userproto.User{Id: "test-user-3"}
				return []*ftproto.Feature{yamlFeature, stringFeature, jsonFeature}, user
			},
			expectedEvalCount: 3,
			validateValue: func(t *testing.T, value string) {
				// All values should be returned as-is or converted appropriately
				assert.NotEmpty(t, value)
			},
		},
		{
			desc: "Evaluate function caches YAML conversion across multiple calls",
			setupFunc: func() ([]*ftproto.Feature, *userproto.User) {
				feature := &ftproto.Feature{
					Id:            "cached-yaml-feature",
					Name:          "Cached YAML Feature",
					Version:       1,
					Enabled:       true,
					VariationType: ftproto.Feature_YAML,
					Variations: []*ftproto.Variation{
						{
							Id:   "cached-yaml-var",
							Name: "Cached Config",
							Value: `settings:
  theme: dark
  language: en
  notifications:
    email: true
    push: false`,
						},
					},
					DefaultStrategy: &ftproto.Strategy{
						Type: ftproto.Strategy_FIXED,
						FixedStrategy: &ftproto.FixedStrategy{
							Variation: "cached-yaml-var",
						},
					},
				}
				user := &userproto.User{Id: "test-user-4"}
				return []*ftproto.Feature{feature}, user
			},
			expectedEvalCount:   1,
			expectedVariationID: "cached-yaml-var",
			validateValue: func(t *testing.T, value string) {
				expectedJSON := `{"settings":{"language":"en","notifications":{"email":true,"push":false},"theme":"dark"}}`
				assert.Equal(t, expectedJSON, value)
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			evaluator := NewEvaluator()
			features, user := p.setupFunc()

			// First evaluation
			result1, err := evaluator.EvaluateFeatures(
				features,
				user,
				map[string][]*ftproto.SegmentUser{},
				"",
			)

			require.NoError(t, err)
			require.NotNil(t, result1)
			require.Len(t, result1.Evaluations, p.expectedEvalCount)

			// Validate each evaluation
			for _, eval := range result1.Evaluations {
				// Check that variation value is set
				assert.NotEmpty(t, eval.VariationValue)

				// Check that Variation.Value is also set correctly
				assert.Equal(t, eval.VariationValue, eval.Variation.Value)

				// Run custom validation if provided
				if p.validateValue != nil {
					p.validateValue(t, eval.VariationValue)
				}

				// If it's a YAML feature, ensure it's valid JSON
				for _, f := range features {
					if f.Id == eval.FeatureId && f.VariationType == ftproto.Feature_YAML {
						var jsonData interface{}
						err := json.Unmarshal([]byte(eval.VariationValue), &jsonData)
						assert.NoError(t, err, "YAML should be converted to valid JSON")
					}
				}
			}

			// Second evaluation with same features (should use cache)
			result2, err := evaluator.EvaluateFeatures(
				features,
				user,
				map[string][]*ftproto.SegmentUser{},
				"",
			)

			require.NoError(t, err)
			require.NotNil(t, result2)

			// Results should be identical (compare by feature ID)
			for _, eval2 := range result2.Evaluations {
				// Find matching evaluation in result1
				for _, eval1 := range result1.Evaluations {
					if eval1.FeatureId == eval2.FeatureId {
						assert.Equal(t, eval1.VariationValue, eval2.VariationValue,
							"Values should match for feature %s", eval1.FeatureId)
						break
					}
				}
			}
		})
	}
}
