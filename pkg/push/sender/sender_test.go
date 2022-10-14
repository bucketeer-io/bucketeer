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

package sender

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestIsFeaturesLatest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		features       *featureproto.Features
		featureID      string
		featureVersion int32
		expected       bool
	}{
		{
			desc: "no feature",
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "wrong", Version: int32(1)}},
			},
			featureID:      "fid",
			featureVersion: int32(1),
			expected:       false,
		},
		{
			desc: "not the latest version",
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "fid", Version: int32(1)}},
			},
			featureID:      "fid",
			featureVersion: int32(2),
			expected:       false,
		},
		{
			desc: "the latest version",
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "fid", Version: int32(2)}},
			},
			featureID:      "fid",
			featureVersion: int32(2),
			expected:       true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &sender{}
			actual := s.isFeaturesLatest(p.features, p.featureID, p.featureVersion)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestExtractFeatureID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc             string
		input            *domaineventproto.Event
		expectedID       string
		expectedIsTarget bool
	}{
		{
			desc: "not feature entity",
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_EXPERIMENT,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_VERSION_INCREMENTED,
			},
			expectedID:       "",
			expectedIsTarget: false,
		},
		{
			desc: "not version incremented",
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_EXPERIMENT,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_DESCRIPTION_CHANGED,
			},
			expectedID:       "",
			expectedIsTarget: false,
		},
		{
			desc: "is target",
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_FEATURE,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_VERSION_INCREMENTED,
			},
			expectedID:       "fid",
			expectedIsTarget: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &sender{}
			actualID, actualIsTarget := s.extractFeatureID(p.input)
			assert.Equal(t, p.expectedID, actualID)
			assert.Equal(t, p.expectedIsTarget, actualIsTarget)
		})
	}
}
