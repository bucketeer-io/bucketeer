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
