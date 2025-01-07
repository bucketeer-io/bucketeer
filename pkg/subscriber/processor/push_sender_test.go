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

package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domaineventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

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
				Type:       domaineventproto.Event_FEATURE_UPDATED,
			},
			expectedID:       "fid",
			expectedIsTarget: true,
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
			s := &pushSender{}
			actualID, actualIsTarget := s.extractFeatureID(p.input)
			assert.Equal(t, p.expectedID, actualID)
			assert.Equal(t, p.expectedIsTarget, actualIsTarget)
		})
	}
}
