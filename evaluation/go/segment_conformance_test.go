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

// The conformance fixtures are shared with evaluation/typescript so the two
// engines stay in lockstep. See evaluation/testdata/segment_rules_conformance.json.
const conformanceFixturePath = "../testdata/segment_rules_conformance.json"

type conformanceClause struct {
	ID        string   `json:"id"`
	Attribute string   `json:"attribute"`
	Operator  string   `json:"operator"`
	Values    []string `json:"values"`
}

type conformanceRule struct {
	ID      string              `json:"id"`
	Clauses []conformanceClause `json:"clauses"`
}

type conformanceSegment struct {
	ID    string            `json:"id"`
	Rules []conformanceRule `json:"rules"`
}

type conformanceSegmentUser struct {
	SegmentID string `json:"segmentId"`
	UserID    string `json:"userId"`
	State     string `json:"state"`
}

type conformanceUser struct {
	ID   string            `json:"id"`
	Data map[string]string `json:"data"`
}

type conformanceTestCase struct {
	Desc       string          `json:"desc"`
	SegmentIDs []string        `json:"segmentIds"`
	User       conformanceUser `json:"user"`
	Expected   bool            `json:"expected"`
}

type conformanceFixture struct {
	Segments     []conformanceSegment     `json:"segments"`
	SegmentUsers []conformanceSegmentUser `json:"segmentUsers"`
	TestCases    []conformanceTestCase    `json:"testCases"`
}

func TestSegmentRulesConformance(t *testing.T) {
	t.Parallel()
	fixture := loadConformanceFixture(t)
	segments := make(map[string]*ftproto.Segment, len(fixture.Segments))
	for _, s := range fixture.Segments {
		segments[s.ID] = &ftproto.Segment{
			Id:    s.ID,
			Rules: toProtoRules(t, s.Rules),
		}
	}
	segmentUsers := make([]*ftproto.SegmentUser, 0, len(fixture.SegmentUsers))
	for _, su := range fixture.SegmentUsers {
		state, ok := ftproto.SegmentUser_State_value[su.State]
		require.True(t, ok, "unknown segment user state: %s", su.State)
		segmentUsers = append(segmentUsers, &ftproto.SegmentUser{
			SegmentId: su.SegmentID,
			UserId:    su.UserID,
			State:     ftproto.SegmentUser_State(state),
		})
	}
	evaluator := &segmentEvaluator{}
	for _, tc := range fixture.TestCases {
		tc := tc
		t.Run(tc.Desc, func(t *testing.T) {
			t.Parallel()
			user := &userproto.User{Id: tc.User.ID, Data: tc.User.Data}
			actual, err := evaluator.Evaluate(tc.SegmentIDs, user, segments, segmentUsers)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func loadConformanceFixture(t *testing.T) *conformanceFixture {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(conformanceFixturePath))
	require.NoError(t, err)
	fixture := &conformanceFixture{}
	require.NoError(t, json.Unmarshal(data, fixture))
	return fixture
}

func toProtoRules(t *testing.T, rules []conformanceRule) []*ftproto.Rule {
	t.Helper()
	protoRules := make([]*ftproto.Rule, 0, len(rules))
	for _, r := range rules {
		clauses := make([]*ftproto.Clause, 0, len(r.Clauses))
		for _, c := range r.Clauses {
			operator, ok := ftproto.Clause_Operator_value[c.Operator]
			require.True(t, ok, "unknown clause operator: %s", c.Operator)
			clauses = append(clauses, &ftproto.Clause{
				Id:        c.ID,
				Attribute: c.Attribute,
				Operator:  ftproto.Clause_Operator(operator),
				Values:    c.Values,
			})
		}
		protoRules = append(protoRules, &ftproto.Rule{
			Id:      r.ID,
			Clauses: clauses,
		})
	}
	return protoRules
}
