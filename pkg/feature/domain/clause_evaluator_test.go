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

	"github.com/stretchr/testify/assert"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestGreaterFloat(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		// Int
		{
			targetValue: "1",
			values:      []string{"1"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"1", "2", "3"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"a", "1", "2.0"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"0a", "1a"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"0.0", "1.0", "2.0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"0.9", "1.0", "2.0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"0", "1", "2"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"a", "0", "1.0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"a", "0", "1"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"0a", "0"},
			expected:    true,
		},
		// Float
		{
			targetValue: "1.0",
			values:      []string{"1.0", "2.0", "3.0"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"1", "2", "3"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"a", "1", "2.0"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"0", "1.0", "2.0"},
			expected:    true,
		},
		{
			targetValue: "1.0",
			values:      []string{"a", "0.0", "1.0"},
			expected:    true,
		},
		{
			targetValue: "1.2",
			values:      []string{"a", "1.1", "2.0"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_GREATER,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestGreaterSemver(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "1.0.0",
			values:      []string{"1.0.0", "0.0", "1.0.1"},
			expected:    false,
		},
		{
			targetValue: "1.0.0",
			values:      []string{"1.0.0", "1.0.1", "v0.0.7"},
			expected:    false,
		},
		{
			targetValue: "0.0.8",
			values:      []string{"1.0.0", "0.0.9", "1.0.1"},
			expected:    false,
		},
		{
			targetValue: "1.1.0",
			values:      []string{"1.1.0", "v1.0.9", "1.1.1"},
			expected:    false,
		},
		{
			targetValue: "2.1.0",
			values:      []string{"2.1.0", "v2.0.9", "2.1.1"},
			expected:    false,
		},
		{
			targetValue: "1.0.1",
			values:      []string{"1.0.1", "1.0.0", "v0.0.7"},
			expected:    true,
		},
		{
			targetValue: "1.1.1",
			values:      []string{"1.1.1", "v1.0.9", "1.1.0"},
			expected:    true,
		},
		{
			targetValue: "2.1.1",
			values:      []string{"2.1.1", "v2.0.9", "2.1.0"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_GREATER,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestGreaterString(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "b",
			values:      []string{"c", "d", "e"},
			expected:    false,
		},
		{
			targetValue: "v1.0.0",
			values:      []string{"v2.0.0", "v1.0.9", "v1.0.8"},
			expected:    false,
		},
		{
			targetValue: "b",
			values:      []string{"1", "a", "2.0"},
			expected:    true,
		},
		{
			targetValue: "b",
			values:      []string{"c", "d", "a"},
			expected:    true,
		},
		{
			targetValue: "v1.0.0",
			values:      []string{"v1.0.0", "v1.0.9", "v0.0.9"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_GREATER,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestGreaterOrEqualFloat(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		// Int
		{
			targetValue: "1",
			values:      []string{"2"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"2", "3", "4"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"2.0", "3.0", "4.0"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"a", "2", "3.0"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"0a", "1a"},
			expected:    false,
		},
		{
			targetValue: "1",
			values:      []string{"1"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"0", "1", "2"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"0.0", "1.0", "2.0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"1.0", "2.0", "3.0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"1", "2", "3"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"a", "0", "1.0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"a", "0", "1"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"a", "1", "2.0"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"a", "1.0", "2"},
			expected:    true,
		},
		{
			targetValue: "1",
			values:      []string{"0a", "0"},
			expected:    true,
		},
		// Float
		{
			targetValue: "1.0",
			values:      []string{"2.0", "3.0", "4.0"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"2", "3", "4"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"a", "1.1", "2.0"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "1.0",
			values:      []string{"0.9", "2.0", "3.0"},
			expected:    true,
		},
		{
			targetValue: "1.0",
			values:      []string{"a", "0", "2.0"},
			expected:    true,
		},
		{
			targetValue: "1.1",
			values:      []string{"1", "2.0", "3.0"},
			expected:    true,
		},
		{
			targetValue: "1.1",
			values:      []string{"1.1", "2.0", "3.0"},
			expected:    true,
		},
		{
			targetValue: "1.1",
			values:      []string{"a", "1.0", "2.0"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_GREATER_OR_EQUAL,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestGreaterOrEqualSemver(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "1.0.0",
			values:      []string{"1.0.1", "0.0", "1.0.2"},
			expected:    false,
		},
		{
			targetValue: "1.0.0",
			values:      []string{"1.0.1", "1.0.2", "v0.0.7"},
			expected:    false,
		},
		{
			targetValue: "0.0.8",
			values:      []string{"1.0.0", "0.0.9", "1.0.1"},
			expected:    false,
		},
		{
			targetValue: "1.1.0",
			values:      []string{"1.1.1", "v1.0.9", "1.1.2"},
			expected:    false,
		},
		{
			targetValue: "2.1.0",
			values:      []string{"2.1.1", "v2.0.9", "2.1.2"},
			expected:    false,
		},
		{
			targetValue: "1.0.0",
			values:      []string{"1.0.1", "1.0.0", "v0.0.7"},
			expected:    true,
		},
		{
			targetValue: "1.1.1",
			values:      []string{"1.1.2", "v1.0.9", "1.1.1"},
			expected:    true,
		},
		{
			targetValue: "2.1.1",
			values:      []string{"2.1.2", "v2.0.9", "2.1.1"},
			expected:    true,
		},
		{
			targetValue: "1.0.1",
			values:      []string{"1.0.2", "1.0.1", "v0.0.7"},
			expected:    true,
		},
		{
			targetValue: "1.1.1",
			values:      []string{"1.1.2", "v1.0.9", "1.1.0"},
			expected:    true,
		},
		{
			targetValue: "2.1.1",
			values:      []string{"2.1.2", "v2.0.9", "2.1.0"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_GREATER_OR_EQUAL,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestGreaterOrEqualString(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "b",
			values:      []string{"c", "d", "e"},
			expected:    false,
		},
		{
			targetValue: "v1.0.0",
			values:      []string{"v2.0.0", "v1.0.9", "v1.0.8"},
			expected:    false,
		},
		{
			targetValue: "b",
			values:      []string{"1", "a", "2.0"},
			expected:    true,
		},
		{
			targetValue: "b",
			values:      []string{"d", "c", "b"},
			expected:    true,
		},
		{
			targetValue: "b",
			values:      []string{"c", "d", "a"},
			expected:    true,
		},
		{
			targetValue: "b",
			values:      []string{"d", "c", "b"},
			expected:    true,
		},
		{
			targetValue: "v1.0.0",
			values:      []string{"v1.0.8", "v1.0.9", "v1.0.0"},
			expected:    true,
		},
		{
			targetValue: "v1.0.0",
			values:      []string{"v1.0.8", "v1.0.9", "v0.0.9"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_GREATER_OR_EQUAL,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestLessThanSemver(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "1.0.0",
			values:      []string{"1.0.0", "0.0", "0.0.9"},
			expected:    false,
		},
		{
			targetValue: "1.0.0",
			values:      []string{"1.0.0", "v0.0.8", "0.0.7"},
			expected:    false,
		},
		{
			targetValue: "0.0.8",
			values:      []string{"0.0.8", "0.0.7", "v0.0.9"},
			expected:    false,
		},
		{
			targetValue: "1.1.0",
			values:      []string{"1.1.0", "v1.0.9", "1.0.8"},
			expected:    false,
		},
		{
			targetValue: "2.1.0",
			values:      []string{"2.1.0", "v2.0.9", "2.0.9"},
			expected:    false,
		},
		{
			targetValue: "1.0.1",
			values:      []string{"1.0.1", "v0.0.7", "1.0.2"},
			expected:    true,
		},
		{
			targetValue: "1.1.1",
			values:      []string{"1.1.1", "v1.0.9", "1.1.2"},
			expected:    true,
		},
		{
			targetValue: "2.1.1",
			values:      []string{"2.1.1", "v2.0.9", "2.1.2"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_LESS,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestLessFloat(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		// Int
		{
			targetValue: "3",
			values:      []string{"3"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"1", "2", "3"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"a", "1", "2.0"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"0a", "1a"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"4"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"2.0", "3.0", "4.0"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"1.0", "2.0", "3.1"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"2", "3", "4"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"d", "3", "3.5"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"a", "0", "4"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"4a", "4"},
			expected:    true,
		},
		// Float
		{
			targetValue: "3.0",
			values:      []string{"1.0", "2.0", "3.0"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"1", "2", "3"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"a", "1", "2.0"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"2", "3.0", "3.1"},
			expected:    true,
		},
		{
			targetValue: "3.0",
			values:      []string{"a", "0.0", "3.9"},
			expected:    true,
		},
		{
			targetValue: "3.2",
			values:      []string{"a", "1.1", "3.5"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_LESS,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestLessString(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "c",
			values:      []string{"c", "b", "a"},
			expected:    false,
		},
		{
			targetValue: "c",
			values:      []string{"1", "a", "2.0"},
			expected:    false,
		},
		{
			targetValue: "v2.0.0",
			values:      []string{"v2.0.0", "v1.0.9", "v1.0.8"},
			expected:    false,
		},
		{
			targetValue: "c",
			values:      []string{"b", "c", "d"},
			expected:    true,
		},
		{
			targetValue: "c",
			values:      []string{"3", "1.0", "d"},
			expected:    true,
		},
		{
			targetValue: "v2.0.0",
			values:      []string{"v1.0.0", "v1.0.9", "v2.1.0"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_LESS,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestLessOrEqualFloat(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		// Int
		{
			targetValue: "3",
			values:      []string{"2"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"0", "1", "2"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"0", "1.0", "2.0"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"a", "1", "2.0"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"3a", "4a"},
			expected:    false,
		},
		{
			targetValue: "3",
			values:      []string{"3"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"2", "3", "4"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"1.0", "2.0", "3.0"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"1.0", "2.0", "3.1"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"1", "2", "4"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"a", "0", "3.0"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"a", "1.0", "4"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"a", "1", "3.5"},
			expected:    true,
		},
		{
			targetValue: "3",
			values:      []string{"3a", "3"},
			expected:    true,
		},
		// Float
		{
			targetValue: "3.0",
			values:      []string{"0", "1.0", "2.0"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"0", "1", "2"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"a", "1.1", "2.0"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "3.0",
			values:      []string{"0.9", "2.0", "3.0"},
			expected:    true,
		},
		{
			targetValue: "3.0",
			values:      []string{"a", "0", "3.1"},
			expected:    true,
		},
		{
			targetValue: "3.1",
			values:      []string{"1", "2.0", "3.9"},
			expected:    true,
		},
		{
			targetValue: "3.1",
			values:      []string{"1.1", "2.0", "4"},
			expected:    true,
		},
		{
			targetValue: "3.1",
			values:      []string{"a", "1.0", "3.1"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_LESS_OR_EQUAL,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestLessThanOrEqualSemver(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "1.0.1",
			values:      []string{"1.0.0", "0.0", "0.0.9"},
			expected:    false,
		},
		{
			targetValue: "1.0.1",
			values:      []string{"1.0.0", "v0.0.8", "0.0.7"},
			expected:    false,
		},
		{
			targetValue: "0.0.9",
			values:      []string{"0.0.8", "0.0.7", "v0.0.9"},
			expected:    false,
		},
		{
			targetValue: "1.1.1",
			values:      []string{"1.1.0", "v1.0.9", "1.0.8"},
			expected:    false,
		},
		{
			targetValue: "2.1.1",
			values:      []string{"2.1.0", "v2.0.9", "2.0.9"},
			expected:    false,
		},
		{
			targetValue: "1.0.1",
			values:      []string{"1.0.1", "v0.0.7", "1.0.0"},
			expected:    true,
		},
		{
			targetValue: "1.1.1",
			values:      []string{"1.1.1", "v1.0.9", "1.1.0"},
			expected:    true,
		},
		{
			targetValue: "2.1.1",
			values:      []string{"2.1.1", "v2.0.9", "2.1.0"},
			expected:    true,
		},
		{
			targetValue: "1.0.1",
			values:      []string{"1.0.0", "v0.0.7", "1.0.2"},
			expected:    true,
		},
		{
			targetValue: "1.1.1",
			values:      []string{"1.1.0", "v1.0.9", "1.1.2"},
			expected:    true,
		},
		{
			targetValue: "2.1.1",
			values:      []string{"2.1.0", "v2.0.9", "2.1.2"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_LESS_OR_EQUAL,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestLessOrEqualString(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		{
			targetValue: "d",
			values:      []string{"a", "b", "c"},
			expected:    false,
		},
		{
			targetValue: "c",
			values:      []string{"1", "a", "2.0"},
			expected:    false,
		},
		{
			targetValue: "v2.0.0",
			values:      []string{"v1.0.0", "v1.0.9", "v1.0.8"},
			expected:    false,
		},
		{
			targetValue: "c",
			values:      []string{"3.0", "c", "b"},
			expected:    true,
		},
		{
			targetValue: "c",
			values:      []string{"c", "b", "a"},
			expected:    true,
		},
		{
			targetValue: "c",
			values:      []string{"a", "b", "d"},
			expected:    true,
		},
		{
			targetValue: "v2.0.0",
			values:      []string{"v1.0.0", "v1.0.9", "v2.0.0"},
			expected:    true,
		},
		{
			targetValue: "v2.0.0",
			values:      []string{"v1.0.0", "v1.0.9", "v2.0.1"},
			expected:    true,
		},
	}
	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_LESS_OR_EQUAL,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestBeforeInt(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		// Int
		{
			targetValue: "1519223320",
			values:      []string{"1419223320"},
			expected:    false,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1619223320"},
			expected:    true,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1519223320", "1519200000"},
			expected:    false,
		},
		// Strings
		{
			targetValue: "15192XXX23320",
			values:      []string{"1519223330", "1519223311", "1519223300"},
			expected:    false,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1519223320", "1519200000", "15192XXX23300"},
			expected:    false,
		},
		// Float
		{
			targetValue: "15192233.30",
			values:      []string{"1519223330", "1519223311", "1519223300"},
			expected:    false,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1519223320", "1519200000", "15192233.00"},
			expected:    false,
		},
	}

	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_BEFORE,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}

func TestAfterInt(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		targetValue string
		values      []string
		expected    bool
	}{
		// Int
		{
			targetValue: "1519223320",
			values:      []string{"1419223320"},
			expected:    true,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1619223320"},
			expected:    false,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1519223320", "1519223319"},
			expected:    true,
		},
		// Strings
		{
			targetValue: "15192XXX23320",
			values:      []string{"1519223330", "1519223311", "1519223300"},
			expected:    false,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1519223320", "1519200000", "15192XXX23300"},
			expected:    true,
		},
		// Float
		{
			targetValue: "15192233.30",
			values:      []string{"1519223330", "1519223311", "1519223300"},
			expected:    false,
		},
		{
			targetValue: "1519223320",
			values:      []string{"1519223320", "1519200000", "15192233.00"},
			expected:    true,
		},
	}

	clauseEvaluator := &clauseEvaluator{}
	for i, tc := range testcases {
		clause := &featureproto.Clause{
			Operator: featureproto.Clause_AFTER,
			Values:   tc.values,
		}
		des := fmt.Sprintf("index: %d", i)
		res := clauseEvaluator.Evaluate(tc.targetValue, clause, "userId", nil)
		assert.Equal(t, tc.expected, res, des)
	}
}
