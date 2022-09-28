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
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewUserEvaluations(t *testing.T) {
	patterns := []struct {
		id          string
		evaluations []*proto.Evaluation
		expected    *proto.UserEvaluations
	}{
		{
			id:          "1234",
			evaluations: []*proto.Evaluation{{Id: "test-id"}},
			expected: &proto.UserEvaluations{
				Id:          "1234",
				Evaluations: []*proto.Evaluation{{Id: "test-id"}},
			},
		},
	}

	for _, p := range patterns {
		actual := NewUserEvaluations(p.id, p.evaluations)
		assert.Equal(t, p.expected.Id, actual.Id)
		assert.Equal(t, p.expected.Evaluations, actual.Evaluations)
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
