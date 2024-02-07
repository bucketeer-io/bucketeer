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
	"github.com/stretchr/testify/require"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestAddClauseValueToSegment(t *testing.T) {
	testcases := []struct {
		ruleID, clauseID string
		size             int
		expected         error
	}{
		{
			ruleID:   "rule-id-2",
			clauseID: "clause-id-2",
			size:     2,
			expected: errRuleNotFound,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-2",
			size:     2,
			expected: errClauseNotFound,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-1",
			size:     3,
			expected: nil,
		},
	}
	s := newSegment(t)
	rule := &featureproto.Rule{
		Id: "rule-id-1",
		Clauses: []*featureproto.Clause{
			{
				Id:     "clause-id-1",
				Values: []string{"value-1", "value-2"},
			},
		},
	}
	err := s.AddRule(rule)
	require.NoError(t, err)
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		err := s.AddClauseValue(tc.ruleID, tc.clauseID, "value-3")
		assert.Equal(t, tc.expected, err, des)
		assert.Equal(t, tc.size, len(s.Rules[0].Clauses[0].Values), des)
	}
	clause, err := s.findClause(rule.Id, rule.Clauses[0].Id)
	require.NoError(t, err)
	idx, err := index("value-3", clause.Values)
	assert.Equal(t, idx != -1, err == nil)
}

func TestRemoveClauseValueFromSegment(t *testing.T) {
	testcases := []struct {
		ruleID, clauseID, value string
		size                    int
		expected                error
	}{
		{
			ruleID:   "rule-id-2",
			clauseID: "clause-id-2",
			value:    "value-3",
			size:     2,
			expected: errRuleNotFound,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-2",
			value:    "value-3",
			size:     2,
			expected: errClauseNotFound,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-1",
			value:    "value-3",
			size:     2,
			expected: errValueNotFound,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-1",
			value:    "value-2",
			size:     1,
			expected: nil,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-1",
			value:    "value-1",
			size:     1,
			expected: errClauseMustHaveAtLeastOneValue,
		},
	}
	s := newSegment(t)
	rule := &featureproto.Rule{
		Id: "rule-id-1",
		Clauses: []*featureproto.Clause{
			{
				Id:     "clause-id-1",
				Values: []string{"value-1", "value-2"},
			},
		},
	}
	err := s.AddRule(rule)
	require.NoError(t, err)
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		err := s.RemoveClauseValue(tc.ruleID, tc.clauseID, tc.value)
		assert.Equal(t, tc.expected, err, des)
		assert.Equal(t, tc.size, len(s.Rules[0].Clauses[0].Values), des)
	}
	clause, err := s.findClause(rule.Id, rule.Clauses[0].Id)
	require.NoError(t, err)
	idx, err := index("value-2", clause.Values)
	assert.Equal(t, idx == -1, err == errValueNotFound)
}

func TestFindRuleIndex(t *testing.T) {
	testcases := []struct {
		ruleID   string
		index    int
		expected error
	}{
		{
			ruleID:   "rule-id-2",
			index:    -1,
			expected: errRuleNotFound,
		},
		{
			ruleID:   "rule-id-1",
			index:    0,
			expected: nil,
		},
	}
	s := newSegment(t)
	rule := &featureproto.Rule{Id: "rule-id-1"}
	err := s.AddRule(rule)
	require.NoError(t, err)
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		index, err := s.findRuleIndex(tc.ruleID)
		assert.Equal(t, tc.expected, err, des)
		assert.Equal(t, tc.index, index, des)
	}
}

func TestFindClauseIndex(t *testing.T) {
	testcases := []struct {
		clauseID string
		index    int
		expected error
	}{
		{
			clauseID: "clause-id-3",
			index:    -1,
			expected: errClauseNotFound,
		},
		{
			clauseID: "clause-id-2",
			index:    1,
			expected: nil,
		},
		{
			clauseID: "clause-id-1",
			index:    0,
			expected: nil,
		},
	}
	clauses := []*featureproto.Clause{
		{Id: "clause-id-1"},
		{Id: "clause-id-2"},
	}
	s := newSegment(t)
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		index, err := s.findClauseIndex(tc.clauseID, clauses)
		assert.Equal(t, tc.expected, err, des)
		assert.Equal(t, tc.index, index, des)
	}
}

func TestFindClause(t *testing.T) {
	clause1 := &featureproto.Clause{
		Id: "clause-id-1",
	}
	clause2 := &featureproto.Clause{
		Id: "clause-id-2",
	}
	testcases := []struct {
		ruleID, clauseID string
		clause           *featureproto.Clause
		expected         error
	}{
		{
			ruleID:   "rule-id-2",
			clauseID: "clause-id-1",
			clause:   nil,
			expected: errRuleNotFound,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-3",
			clause:   nil,
			expected: errClauseNotFound,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-2",
			clause:   clause2,
			expected: nil,
		},
		{
			ruleID:   "rule-id-1",
			clauseID: "clause-id-1",
			clause:   clause1,
			expected: nil,
		},
	}
	s := newSegment(t)
	rule := &featureproto.Rule{
		Id: "rule-id-1",
		Clauses: []*featureproto.Clause{
			clause1,
			clause2,
		},
	}
	err := s.AddRule(rule)
	require.NoError(t, err)
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		clause, err := s.findClause(tc.ruleID, tc.clauseID)
		assert.Equal(t, tc.expected, err, des)
		assert.Equal(t, tc.clause, clause, des)
	}
}

func newSegment(t *testing.T) *Segment {
	s, err := NewSegment("name", "description")
	require.NoError(t, err)
	return s
}
