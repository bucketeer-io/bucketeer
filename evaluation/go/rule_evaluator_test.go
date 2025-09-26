// Copyright 2025 The Bucketeer Authors.
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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestRuleEvaluator(t *testing.T) {
	f := newFeature()
	testcases := []struct {
		user     *userproto.User
		expected *ftproto.Rule
	}{
		{
			user: &userproto.User{
				Id:   "user-id-1",
				Data: map[string]string{"full-name": "bucketeer project"},
			},
			expected: f.Rules[0],
		},
		{
			user: &userproto.User{
				Id:   "user-id-1",
				Data: map[string]string{"first-name": "bucketeer"},
			},
			expected: f.Rules[1],
		},
		{
			user: &userproto.User{
				Id:   "user-id-1",
				Data: map[string]string{"last-name": "project"},
			},
			expected: f.Rules[2],
		},
		{
			user: &userproto.User{
				Id:   "user-id-3",
				Data: map[string]string{"email": "bucketeer@gmail.com"},
			},
			expected: f.Rules[4],
		},
		{
			user: &userproto.User{
				Id:   "user-id-1",
				Data: nil,
			},
			expected: f.Rules[3],
		},
		{
			user: &userproto.User{
				Id:   "user-id-2",
				Data: nil,
			},
			expected: f.Rules[3],
		},
		{
			user: &userproto.User{
				Id:   "user-id-3",
				Data: nil,
			},
			expected: nil,
		},
		{
			user: &userproto.User{
				Id:   "user-id-4",
				Data: nil,
			},
			expected: nil,
		},
	}
	values := newSegmentUserIDs()
	ruleEvaluator := &ruleEvaluator{}
	for i, tc := range testcases {
		des := fmt.Sprintf("index: %d", i)
		actual, _ := ruleEvaluator.Evaluate(f.Rules, tc.user, values, nil)
		assert.Equal(t, tc.expected, actual, des)
	}
}

func newFeature() *ftproto.Feature {
	return &ftproto.Feature{
		Id:        "feature-id",
		Name:      "test feature",
		Version:   1,
		CreatedAt: time.Now().Unix(),
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
		},
		Rules: []*ftproto.Rule{
			{
				Id: "rule-id-1",
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: "variation-A",
					},
				},
				Clauses: []*ftproto.Clause{
					{
						Id:        "clause-id-1",
						Attribute: "full-name",
						Operator:  ftproto.Clause_EQUALS,
						Values:    []string{"bucketeer project"},
					},
				},
			},
			{
				Id: "rule-id-2",
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: "variation-A",
					},
				},
				Clauses: []*ftproto.Clause{
					{
						Id:        "clause-id-2",
						Attribute: "first-name",
						Operator:  ftproto.Clause_STARTS_WITH,
						Values:    []string{"buck"},
					},
				},
			},
			{
				Id: "rule-id-3",
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: "variation-A",
					},
				},
				Clauses: []*ftproto.Clause{
					{
						Id:        "clause-id-3",
						Attribute: "last-name",
						Operator:  ftproto.Clause_ENDS_WITH,
						Values:    []string{"ject"},
					},
				},
			},
			{
				Id: "rule-id-4",
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: "variation-B",
					},
				},
				Clauses: []*ftproto.Clause{
					{
						Id:        "clause-id-4",
						Attribute: "",
						Operator:  ftproto.Clause_SEGMENT,
						Values: []string{
							"segment-id-1",
							"segment-id-2",
						},
					},
				},
			},
			{
				Id: "rule-id-5",
				Strategy: &ftproto.Strategy{
					Type: ftproto.Strategy_FIXED,
					FixedStrategy: &ftproto.FixedStrategy{
						Variation: "variation-B",
					},
				},
				Clauses: []*ftproto.Clause{
					{
						Id:        "clause-id-5",
						Attribute: "email",
						Operator:  ftproto.Clause_IN,
						Values:    []string{"bucketeer@gmail.com"},
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

func newSegmentUserIDs() (values []*ftproto.SegmentUser) {
	values = append(values, &ftproto.SegmentUser{
		UserId:    "user-id-1",
		SegmentId: "segment-id-1",
		State:     ftproto.SegmentUser_INCLUDED,
	})
	values = append(values, &ftproto.SegmentUser{
		UserId:    "user-id-1",
		SegmentId: "segment-id-2",
		State:     ftproto.SegmentUser_INCLUDED,
	})
	values = append(values, &ftproto.SegmentUser{
		UserId:    "user-id-2",
		SegmentId: "segment-id-1",
		State:     ftproto.SegmentUser_INCLUDED,
	})
	values = append(values, &ftproto.SegmentUser{
		UserId:    "user-id-2",
		SegmentId: "segment-id-2",
		State:     ftproto.SegmentUser_INCLUDED,
	})
	values = append(values, &ftproto.SegmentUser{
		UserId:    "user-id-3",
		SegmentId: "segment-id-1",
		State:     ftproto.SegmentUser_INCLUDED,
	})
	values = append(values, &ftproto.SegmentUser{
		UserId:    "user-id-4",
		SegmentId: "segment-id-2",
		State:     ftproto.SegmentUser_INCLUDED,
	})
	return values
}
