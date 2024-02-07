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

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

func TestRuleEvaluator(t *testing.T) {
	f := newFeature()
	testcases := []struct {
		user     *userproto.User
		expected *featureproto.Rule
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
		assert.Equal(t, tc.expected, ruleEvaluator.Evaluate(f.Rules, tc.user, values), des)
	}
}

func newFeature() *Feature {
	return &Feature{
		Feature: &featureproto.Feature{
			Id:        "feature-id",
			Name:      "test feature",
			Version:   1,
			CreatedAt: time.Now().Unix(),
			Variations: []*featureproto.Variation{
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
			Rules: []*featureproto.Rule{
				{
					Id: "rule-id-1",
					Strategy: &featureproto.Strategy{
						Type: featureproto.Strategy_FIXED,
						FixedStrategy: &featureproto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*featureproto.Clause{
						{
							Id:        "clause-id-1",
							Attribute: "full-name",
							Operator:  featureproto.Clause_EQUALS,
							Values:    []string{"bucketeer project"},
						},
					},
				},
				{
					Id: "rule-id-2",
					Strategy: &featureproto.Strategy{
						Type: featureproto.Strategy_FIXED,
						FixedStrategy: &featureproto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*featureproto.Clause{
						{
							Id:        "clause-id-2",
							Attribute: "first-name",
							Operator:  featureproto.Clause_STARTS_WITH,
							Values:    []string{"buck"},
						},
					},
				},
				{
					Id: "rule-id-3",
					Strategy: &featureproto.Strategy{
						Type: featureproto.Strategy_FIXED,
						FixedStrategy: &featureproto.FixedStrategy{
							Variation: "variation-A",
						},
					},
					Clauses: []*featureproto.Clause{
						{
							Id:        "clause-id-3",
							Attribute: "last-name",
							Operator:  featureproto.Clause_ENDS_WITH,
							Values:    []string{"ject"},
						},
					},
				},
				{
					Id: "rule-id-4",
					Strategy: &featureproto.Strategy{
						Type: featureproto.Strategy_FIXED,
						FixedStrategy: &featureproto.FixedStrategy{
							Variation: "variation-B",
						},
					},
					Clauses: []*featureproto.Clause{
						{
							Id:        "clause-id-4",
							Attribute: "",
							Operator:  featureproto.Clause_SEGMENT,
							Values: []string{
								"segment-id-1",
								"segment-id-2",
							},
						},
					},
				},
				{
					Id: "rule-id-5",
					Strategy: &featureproto.Strategy{
						Type: featureproto.Strategy_FIXED,
						FixedStrategy: &featureproto.FixedStrategy{
							Variation: "variation-B",
						},
					},
					Clauses: []*featureproto.Clause{
						{
							Id:        "clause-id-5",
							Attribute: "email",
							Operator:  featureproto.Clause_IN,
							Values:    []string{"bucketeer@gmail.com"},
						},
					},
				},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type: featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{
					Variation: "variation-B",
				},
			},
		},
	}
}

func newSegmentUserIDs() (values []*featureproto.SegmentUser) {
	values = append(values, &featureproto.SegmentUser{
		UserId:    "user-id-1",
		SegmentId: "segment-id-1",
		State:     featureproto.SegmentUser_INCLUDED,
	})
	values = append(values, &featureproto.SegmentUser{
		UserId:    "user-id-1",
		SegmentId: "segment-id-2",
		State:     featureproto.SegmentUser_INCLUDED,
	})
	values = append(values, &featureproto.SegmentUser{
		UserId:    "user-id-2",
		SegmentId: "segment-id-1",
		State:     featureproto.SegmentUser_INCLUDED,
	})
	values = append(values, &featureproto.SegmentUser{
		UserId:    "user-id-2",
		SegmentId: "segment-id-2",
		State:     featureproto.SegmentUser_INCLUDED,
	})
	values = append(values, &featureproto.SegmentUser{
		UserId:    "user-id-3",
		SegmentId: "segment-id-1",
		State:     featureproto.SegmentUser_INCLUDED,
	})
	values = append(values, &featureproto.SegmentUser{
		UserId:    "user-id-4",
		SegmentId: "segment-id-2",
		State:     featureproto.SegmentUser_INCLUDED,
	})
	return values
}
