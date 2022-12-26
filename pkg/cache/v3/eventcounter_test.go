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

package v3

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestGetEventValues(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		vals     []interface{}
		expected []float64
		inValid  bool
	}{
		{
			desc: "success: all vals are numbers",
			vals: []interface{}{
				"1", "2", "32", "91",
			},
			expected: []float64{
				1, 2, 32, 91,
			},
			inValid: false,
		},
		{
			desc: "success: some vals are nil",
			vals: []interface{}{
				"1", "2", nil, "32", "91", nil,
			},
			expected: []float64{
				1, 2, 0, 32, 91, 0,
			},
			inValid: false,
		},
		{
			desc: "fail: invalid vals: bool",
			vals: []interface{}{
				true,
			},
			expected: []float64{},
			inValid:  true,
		},
		{
			desc: "fail: invalid vals: alphabet",
			vals: []interface{}{
				"a",
			},
			expected: []float64{},
			inValid:  true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := getEventValues(p.vals)
			assert.Equal(t, p.expected, actual)
			if p.inValid {
				assert.Error(t, err)
			}
		})
	}
}

func TestGetUserValues(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		cmds     []*redis.IntCmd
		expected []float64
		inValid  bool
	}{
		{
			desc: "success",
			cmds: []*redis.IntCmd{
				redis.NewIntCmd(),
				redis.NewIntCmd(),
				redis.NewIntCmd(),
				redis.NewIntCmd(),
			},
			expected: []float64{
				0, 0, 0, 0,
			},
			inValid: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := getUserValues(p.cmds)
			assert.Equal(t, p.expected, actual)
			if p.inValid {
				assert.Error(t, err)
			}
		})
	}
}
