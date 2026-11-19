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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewUserEvaluations(t *testing.T) {
	patterns := []struct {
		id          string
		evaluations []*proto.Evaluation
		archivedIDs []string
		forceUpdate bool
		expected    *proto.UserEvaluations
	}{
		{
			id:          "1234",
			evaluations: []*proto.Evaluation{{Id: "test-id1"}},
			archivedIDs: []string{"test-id2"},
			forceUpdate: false,
			expected: &proto.UserEvaluations{
				Id:                 "1234",
				Evaluations:        []*proto.Evaluation{{Id: "test-id1"}},
				ArchivedFeatureIds: []string{"test-id2"},
				ForceUpdate:        false,
			},
		},
		{
			id:          "5678",
			evaluations: []*proto.Evaluation{{Id: "test-id3"}},
			archivedIDs: []string{},
			forceUpdate: true,
			expected: &proto.UserEvaluations{
				Id:                 "5678",
				Evaluations:        []*proto.Evaluation{{Id: "test-id3"}},
				ArchivedFeatureIds: []string{},
				ForceUpdate:        true,
			},
		},
	}

	for _, p := range patterns {
		actual := NewUserEvaluations(p.id, p.evaluations, p.archivedIDs, p.forceUpdate)
		assert.Equal(t, p.expected.Id, actual.Id)
		assert.Equal(t, p.expected.Evaluations, actual.Evaluations)
		assert.Equal(t, p.expected.ArchivedFeatureIds, actual.ArchivedFeatureIds)
		assert.Equal(t, p.expected.ForceUpdate, actual.ForceUpdate)
		assert.NotZero(t, actual.CreatedAt)
	}
}

func TestSortMapKeys(t *testing.T) {
	patterns := []struct {
		input    map[string]string
		expected []string
		desc     string
	}{
		{
			input:    nil,
			expected: []string{},
			desc:     "nil",
		},
		{
			input:    map[string]string{},
			expected: []string{},
			desc:     "empty",
		},
		{
			input:    map[string]string{"b": "value-b", "c": "value-c", "a": "value-a", "d": "value-d"},
			expected: []string{"a", "b", "c", "d"},
			desc:     "success",
		},
	}
	for _, p := range patterns {
		keys := sortMapKeys(p.input)
		assert.Equal(t, p.expected, keys, p.desc)
	}
}

func TestUserEvaluationsID(t *testing.T) {
	patterns := []struct {
		desc         string
		userID       string
		userMetadata map[string]string
		features     []*ftproto.Feature
		expected     string
	}{
		{
			desc:         "empty user ID, empty metadata, empty features",
			userID:       "",
			userMetadata: nil,
			features:     nil,
			expected:     "14695981039346656037",
		},
		{
			desc:         "user ID only",
			userID:       "user-1",
			userMetadata: nil,
			features:     nil,
			expected:     "17891572797655370708",
		},
		{
			desc:         "user ID with metadata",
			userID:       "user-1",
			userMetadata: map[string]string{"age": "25", "country": "jp"},
			features:     nil,
			expected:     "15857499200645826216",
		},
		{
			desc:         "user ID with metadata and single feature",
			userID:       "user-1",
			userMetadata: map[string]string{"age": "25", "country": "jp"},
			features: []*ftproto.Feature{
				{
					Id:        "feature-1",
					UpdatedAt: 1000,
				},
			},
			expected: "10450974209164395423",
		},
		{
			desc:         "user ID with metadata and multiple features",
			userID:       "user-1",
			userMetadata: map[string]string{"age": "25", "country": "jp"},
			features: []*ftproto.Feature{
				{
					Id:        "feature-1",
					UpdatedAt: 1000,
				},
				{
					Id:        "feature-2",
					UpdatedAt: 2000,
				},
			},
			expected: "7257619227440290900",
		},
	}
	for _, p := range patterns {
		actual := UserEvaluationsID(p.userID, p.userMetadata, p.features)
		assert.Equal(t, p.expected, actual, p.desc)
	}

	features := []*ftproto.Feature{{Id: "feature-1", UpdatedAt: 1000}}
	without := UserEvaluationsID("user-1", nil, features)
	assert.Equal(t, without, UserEvaluationsID("user-1", nil, features, nil))
	groupA := &ftproto.HoldoutGroup{Id: "hg-a", UpdatedAt: 50}
	groupB := &ftproto.HoldoutGroup{Id: "hg-b", UpdatedAt: 60}
	withA := UserEvaluationsID("user-1", nil, features, groupA)
	assert.Equal(t, "13300710561474652895", withA)
	assert.NotEqual(t, without, withA)
	assert.Equal(t, withA, UserEvaluationsID("user-1", nil, features, groupA))
	assert.Equal(t,
		UserEvaluationsID("user-1", nil, features, groupA, groupB),
		UserEvaluationsID("user-1", nil, features, groupB, groupA),
	)
	assert.NotEqual(t, withA, UserEvaluationsID("user-1", nil, features, &ftproto.HoldoutGroup{Id: "hg-a", UpdatedAt: 51}))

	// Same UpdatedAt, different targets. A client polling between two edits in
	// one second must not treat the second edit as unchanged.
	targetA := &ftproto.HoldoutGroup{
		Id:        "hg-a",
		UpdatedAt: 50,
		Features: []*ftproto.HoldoutGroupFeature{
			{FeatureId: "flag-b", HoldoutVariationId: "hold-b"},
			{FeatureId: "flag-a", HoldoutVariationId: "hold-a"},
		},
	}
	targetB := &ftproto.HoldoutGroup{
		Id:        "hg-a",
		UpdatedAt: 50,
		Features: []*ftproto.HoldoutGroupFeature{
			{FeatureId: "flag-a", HoldoutVariationId: "hold-a"},
			{FeatureId: "flag-b", HoldoutVariationId: "other"},
		},
	}
	reordered := &ftproto.HoldoutGroup{
		Id:        "hg-a",
		UpdatedAt: 50,
		Features: []*ftproto.HoldoutGroupFeature{
			{FeatureId: "flag-a", HoldoutVariationId: "hold-a"},
			{FeatureId: "flag-b", HoldoutVariationId: "hold-b"},
		},
	}
	removed := &ftproto.HoldoutGroup{
		Id:        "hg-a",
		UpdatedAt: 50,
		Features: []*ftproto.HoldoutGroupFeature{
			{FeatureId: "flag-a", HoldoutVariationId: "hold-a", RemovedAt: 50},
		},
	}
	withTargets := UserEvaluationsID("user-1", nil, features, targetA)
	assert.NotEqual(t, withA, withTargets)
	assert.NotEqual(t, withTargets, UserEvaluationsID("user-1", nil, features, targetB))
	assert.Equal(t, withTargets, UserEvaluationsID("user-1", nil, features, reordered))
	assert.Equal(t, withA, UserEvaluationsID("user-1", nil, features, removed))
	// Replacing the target row in the same second keeps the feature and the
	// variation, and still has to change the ID the client sends back.
	replaced := &ftproto.HoldoutGroup{
		Id:        "hg-a",
		UpdatedAt: 50,
		Features: []*ftproto.HoldoutGroupFeature{
			{Id: "row-2", FeatureId: "flag-b", HoldoutVariationId: "hold-b", AddedAt: 50},
			{Id: "row-1", FeatureId: "flag-a", HoldoutVariationId: "hold-a"},
		},
	}
	assert.NotEqual(t, withTargets, UserEvaluationsID("user-1", nil, features, replaced))
}

func TestGenerateFeaturesID(t *testing.T) {
	// Note: GenerateFeaturesID uses UpdatedAt (not Version) to generate the hash
	patterns := []struct {
		desc     string
		input    []*ftproto.Feature
		expected string
	}{
		{
			desc:     "nil",
			input:    nil,
			expected: "14695981039346656037",
		},
		{
			desc: "success: single",
			input: []*ftproto.Feature{
				{
					Id:        "id-1",
					UpdatedAt: 1,
				},
			},
			expected: "5476413260388599211",
		},
		{
			desc: "success: multiple",
			input: []*ftproto.Feature{
				{
					Id:        "id-1",
					UpdatedAt: 1,
				},
				{
					Id:        "id-2",
					UpdatedAt: 2,
				},
			},
			expected: "17283374094628184689",
		},
	}
	for _, p := range patterns {
		id := GenerateFeaturesID(p.input)
		assert.Equal(t, p.expected, id, p.desc)
	}
}

func TestGenerateFeaturesIDHoldout(t *testing.T) {
	t.Parallel()
	newFeatures := func() []*ftproto.Feature {
		return []*ftproto.Feature{
			{Id: "id-2", UpdatedAt: 2},
			{Id: "id-1", UpdatedAt: 1},
		}
	}
	withoutGroup := GenerateFeaturesID(newFeatures())
	group := &ftproto.HoldoutGroup{Id: "hg-a", UpdatedAt: 50}
	withTarget := &ftproto.HoldoutGroup{
		Id:        "hg-a",
		UpdatedAt: 50,
		Features: []*ftproto.HoldoutGroupFeature{
			{Id: "row-2", FeatureId: "id-2", HoldoutVariationId: "v-off", AddedAt: 20},
			{Id: "row-1", FeatureId: "id-1", HoldoutVariationId: "v-on", AddedAt: 10},
		},
	}

	// The values are pinned so the TypeScript twin can be checked against them.
	// With a live group the ID is "<flags>:<holdout>".
	assert.Equal(t, "17283374094628184689", withoutGroup)
	assert.Equal(t, "17283374094628184689:16147347951762488293", GenerateFeaturesID(newFeatures(), group))
	assert.Equal(t, "17283374094628184689:6966737048624406302", GenerateFeaturesID(newFeatures(), withTarget))
	assert.Equal(t, "14695981039346656037:16147347951762488293", GenerateFeaturesID(nil, group))

	t.Run("the flags part does not depend on the groups", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, withoutGroup, FeaturesIDFlagsPart(withoutGroup))
		assert.Equal(t, withoutGroup, FeaturesIDFlagsPart(GenerateFeaturesID(newFeatures(), group)))
		assert.Equal(t, withoutGroup, FeaturesIDFlagsPart(GenerateFeaturesID(newFeatures(), withTarget)))
		assert.Equal(t, "", FeaturesIDFlagsPart(""))

		// A flag change moves the flags part and leaves the holdout part.
		edited := newFeatures()
		edited[0].UpdatedAt++
		editedID := GenerateFeaturesID(edited, group)
		assert.NotEqual(t, withoutGroup, FeaturesIDFlagsPart(editedID))
		assert.Equal(t, GenerateFeaturesID(edited), FeaturesIDFlagsPart(editedID))
		assert.Equal(t,
			strings.TrimPrefix(GenerateFeaturesID(newFeatures(), group), withoutGroup),
			strings.TrimPrefix(editedID, FeaturesIDFlagsPart(editedID)),
		)
	})

	t.Run("a group with no id leaves the id unchanged", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, withoutGroup, GenerateFeaturesID(newFeatures(), nil))
		assert.Equal(t, withoutGroup, GenerateFeaturesID(newFeatures(), ClearanceSnapshot(100)))
		assert.Equal(t, withoutGroup, GenerateFeaturesID(newFeatures(), UnavailableHoldoutSnapshot()))
	})

	t.Run("starting, stopping and editing a group change the id", func(t *testing.T) {
		t.Parallel()
		started := GenerateFeaturesID(newFeatures(), group)
		assert.NotEqual(t, withoutGroup, started)
		assert.NotEqual(t, started, GenerateFeaturesID(
			newFeatures(), &ftproto.HoldoutGroup{Id: "hg-a", UpdatedAt: 51},
		))
		assert.NotEqual(t, started, GenerateFeaturesID(
			newFeatures(), &ftproto.HoldoutGroup{Id: "hg-b", UpdatedAt: 50},
		))
	})

	t.Run("targets are part of the id even within the same second", func(t *testing.T) {
		t.Parallel()
		base := GenerateFeaturesID(newFeatures(), withTarget)
		assert.NotEqual(t, GenerateFeaturesID(newFeatures(), group), base)

		reordered := &ftproto.HoldoutGroup{
			Id:        "hg-a",
			UpdatedAt: 50,
			Features: []*ftproto.HoldoutGroupFeature{
				withTarget.Features[1],
				withTarget.Features[0],
			},
		}
		assert.Equal(t, base, GenerateFeaturesID(newFeatures(), reordered))

		otherVariation := &ftproto.HoldoutGroup{
			Id:        "hg-a",
			UpdatedAt: 50,
			Features: []*ftproto.HoldoutGroupFeature{
				{Id: "row-2", FeatureId: "id-2", HoldoutVariationId: "v-on", AddedAt: 20},
				{Id: "row-1", FeatureId: "id-1", HoldoutVariationId: "v-on", AddedAt: 10},
			},
		}
		assert.NotEqual(t, base, GenerateFeaturesID(newFeatures(), otherVariation))

		removed := &ftproto.HoldoutGroup{
			Id:        "hg-a",
			UpdatedAt: 50,
			Features: []*ftproto.HoldoutGroupFeature{
				{Id: "row-2", FeatureId: "id-2", HoldoutVariationId: "v-off", AddedAt: 20, RemovedAt: 30},
				{Id: "row-1", FeatureId: "id-1", HoldoutVariationId: "v-on", AddedAt: 10},
			},
		}
		onlyOne := &ftproto.HoldoutGroup{
			Id:        "hg-a",
			UpdatedAt: 50,
			Features: []*ftproto.HoldoutGroupFeature{
				{Id: "row-1", FeatureId: "id-1", HoldoutVariationId: "v-on", AddedAt: 10},
			},
		}
		assert.NotEqual(t, base, GenerateFeaturesID(newFeatures(), removed))
		assert.Equal(t, GenerateFeaturesID(newFeatures(), onlyOne), GenerateFeaturesID(newFeatures(), removed))
	})
}
