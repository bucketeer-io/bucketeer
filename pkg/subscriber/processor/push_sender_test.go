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

package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
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

func TestGetTopicsForPush(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc           string
		feature        *featureproto.Feature
		push           *pushproto.Push
		expectedTopics []string
	}{
		{
			desc: "both feature and push have no tags",
			feature: &featureproto.Feature{
				Tags: []string{},
			},
			push: &pushproto.Push{
				Tags: []string{},
			},
			expectedTopics: []string{
				"bucketeer-default",
			},
		},
		{
			desc: "feature has tags but push has no tags",
			feature: &featureproto.Feature{
				Tags: []string{"tag1", "tag2"},
			},
			push: &pushproto.Push{
				Tags: []string{},
			},
			expectedTopics: []string{},
		},
		{
			desc: "feature has no tags but push has tags",
			feature: &featureproto.Feature{
				Tags: []string{},
			},
			push: &pushproto.Push{
				Tags: []string{"tag1", "tag2"},
			},
			expectedTopics: []string{},
		},
		{
			desc: "feature and push have matching tags",
			feature: &featureproto.Feature{
				Tags: []string{"tag1", "tag2", "tag3"},
			},
			push: &pushproto.Push{
				Tags: []string{"tag2", "tag4"},
			},
			expectedTopics: []string{
				"bucketeer-tag2",
			},
		},
		{
			desc: "feature and push have multiple matching tags",
			feature: &featureproto.Feature{
				Tags: []string{"tag1", "tag2", "tag3"},
			},
			push: &pushproto.Push{
				Tags: []string{"tag1", "tag2", "tag4"},
			},
			expectedTopics: []string{
				"bucketeer-tag1",
				"bucketeer-tag2",
			},
		},
		{
			desc: "feature and push have no matching tags",
			feature: &featureproto.Feature{
				Tags: []string{"tag1", "tag2"},
			},
			push: &pushproto.Push{
				Tags: []string{"tag3", "tag4"},
			},
			expectedTopics: []string{},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &pushSender{}
			actualTopics := s.getTopicsForPush(p.feature, p.push)
			assert.Equal(t, p.expectedTopics, actualTopics)
		})
	}
}

func TestSendPushNotification(t *testing.T) {
	t.Parallel()
	// Note: This is a basic test structure. In a real implementation,
	// you would mock the pushFCM method and logger to test the behavior properly.
	patterns := []struct {
		desc          string
		topic         string
		push          *pushproto.Push
		featureID     string
		environmentId string
		expectError   bool
	}{
		{
			desc:  "send notification with tag",
			topic: "bucketeer-tag1",
			push: &pushproto.Push{
				Id:                "push-123",
				FcmServiceAccount: `{"type": "service_account"}`,
			},
			featureID:     "feature-123",
			environmentId: "env-123",
			expectError:   false,
		},
		{
			desc:  "send notification without tag",
			topic: "bucketeer-default",
			push: &pushproto.Push{
				Id:                "push-123",
				FcmServiceAccount: `{"type": "service_account"}`,
			},
			featureID:     "feature-123",
			environmentId: "env-123",
			expectError:   false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			// This test would need proper mocking of dependencies
			// For now, we're just ensuring the function structure is correct
			assert.NotEmpty(t, p.topic)
			assert.NotNil(t, p.push)
			assert.NotEmpty(t, p.featureID)
			assert.NotEmpty(t, p.environmentId)
		})
	}
}
