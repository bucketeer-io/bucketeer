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
	patterns := map[string]struct {
		features       *featureproto.Features
		featureID      string
		featureVersion int32
		expected       bool
	}{
		"no feature": {
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "wrong", Version: int32(1)}},
			},
			featureID:      "fid",
			featureVersion: int32(1),
			expected:       false,
		},
		"not the latest version": {
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "fid", Version: int32(1)}},
			},
			featureID:      "fid",
			featureVersion: int32(2),
			expected:       false,
		},
		"the latest version": {
			features: &featureproto.Features{
				Features: []*featureproto.Feature{{Id: "fid", Version: int32(2)}},
			},
			featureID:      "fid",
			featureVersion: int32(2),
			expected:       true,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := &sender{}
			actual := s.isFeaturesLatest(p.features, p.featureID, p.featureVersion)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestExtractFeatureID(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input            *domaineventproto.Event
		expectedID       string
		expectedIsTarget bool
	}{
		"not feature entity": {
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_EXPERIMENT,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_VERSION_INCREMENTED,
			},
			expectedID:       "",
			expectedIsTarget: false,
		},
		"not version incremented": {
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_EXPERIMENT,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_DESCRIPTION_CHANGED,
			},
			expectedID:       "",
			expectedIsTarget: false,
		},
		"is target": {
			input: &domaineventproto.Event{
				EntityType: domaineventproto.Event_FEATURE,
				EntityId:   "fid",
				Type:       domaineventproto.Event_FEATURE_VERSION_INCREMENTED,
			},
			expectedID:       "fid",
			expectedIsTarget: true,
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			s := &sender{}
			actualID, actualIsTarget := s.extractFeatureID(p.input)
			assert.Equal(t, p.expectedID, actualID)
			assert.Equal(t, p.expectedIsTarget, actualIsTarget)
		})
	}
}
