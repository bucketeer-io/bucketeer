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

package druid

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

func TestConvToEnvSegments(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc           string
		inputNamespace string
		inputSegments  []string
		expected       []string
	}{
		{
			desc:           "empty environment namespace",
			inputNamespace: "",
			inputSegments: []string{
				"tag",
				"user.data.sgmt",
			},
			expected: []string{
				"tag",
				"user.data.sgmt",
			},
		},
		{
			desc:           "non empty environment namespace",
			inputNamespace: "ns",
			inputSegments: []string{
				"tag",
				"user.data.sgmt",
			},
			expected: []string{
				"tag",
				"ns.user.data.sgmt",
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := convToEnvSegments(p.inputNamespace, p.inputSegments)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestConvToEnvFilters(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc           string
		inputNamespace string
		inputFilters   []*ecproto.Filter
		expected       []*ecproto.Filter
	}{
		{
			desc:           "empty environment namespace",
			inputNamespace: "",
			inputFilters: []*ecproto.Filter{
				{Key: "tag", Operator: ecproto.Filter_EQUALS, Values: []string{"t0"}},
				{Key: "user.data.sgmt", Operator: ecproto.Filter_EQUALS, Values: []string{"d0"}},
			},
			expected: []*ecproto.Filter{
				{Key: "tag", Operator: ecproto.Filter_EQUALS, Values: []string{"t0"}},
				{Key: "user.data.sgmt", Operator: ecproto.Filter_EQUALS, Values: []string{"d0"}},
			},
		},
		{
			desc:           "non empty environment namespace",
			inputNamespace: "ns",
			inputFilters: []*ecproto.Filter{
				{Key: "tag", Operator: ecproto.Filter_EQUALS, Values: []string{"t0"}},
				{Key: "user.data.sgmt", Operator: ecproto.Filter_EQUALS, Values: []string{"d0"}},
			},
			expected: []*ecproto.Filter{
				{Key: "tag", Operator: ecproto.Filter_EQUALS, Values: []string{"t0"}},
				{Key: "ns.user.data.sgmt", Operator: ecproto.Filter_EQUALS, Values: []string{"d0"}},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := convToEnvFilters(p.inputNamespace, p.inputFilters)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestUserDataPattern(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc           string
		inputNamespace string
		expected       string
	}{
		{
			desc:           "empty environment namespace",
			inputNamespace: "",
			expected:       `^user\.data\.(.*)$`,
		},
		{
			desc:           "non empty environment namespace",
			inputNamespace: "ns",
			expected:       `^ns\.user\.data\.(.*)$`,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := userDataPattern(p.inputNamespace)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestRemoveEnvFromUserData(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		inputKey    string
		inputRegexp *regexp.Regexp
		expected    string
	}{
		{
			desc:        "empty environment namespace",
			inputKey:    "user.data.attr",
			inputRegexp: regexp.MustCompile(userDataPattern("")),
			expected:    "user.data.attr",
		},
		{
			desc:        "non empty environment namespace",
			inputKey:    "ns.user.data.attr",
			inputRegexp: regexp.MustCompile(userDataPattern("ns")),
			expected:    "user.data.attr",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := removeEnvFromUserData(p.inputKey, p.inputRegexp)
			assert.Equal(t, p.expected, actual)
		})
	}
}
