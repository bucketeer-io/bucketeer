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

package api

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func newTestStreamSubscriber(t *testing.T, dispatcher *StreamDispatcher) *StreamSubscriber {
	t.Helper()
	inMemoryCache := cachev3.NewInMemoryCache()
	return NewStreamSubscriber(
		dispatcher,
		cachev3.NewFeaturesCache(inMemoryCache, 0),
		cachev3.NewSegmentUsersCache(inMemoryCache, 0),
		zap.NewNop(),
	)
}

func TestHandleMessage(t *testing.T) {
	t.Parallel()
	const envID = "env-1"
	patterns := []struct {
		desc         string
		event        *domaineventproto.Event
		clientTag    string
		wantDispatch bool
		wantTags     []string
	}{
		{
			desc: "feature updated dispatches with affected tags",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				Type:          domaineventproto.Event_FEATURE_UPDATED,
				EnvironmentId: envID,
				EntityData:    featureTagsJSON(t, "android"),
			},
			clientTag:    "android",
			wantDispatch: true,
			wantTags:     []string{"android"},
		},
		{
			desc: "feature enabled dispatches",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				Type:          domaineventproto.Event_FEATURE_ENABLED,
				EnvironmentId: envID,
				EntityData:    featureTagsJSON(t, "android"),
			},
			clientTag:    "android",
			wantDispatch: true,
			wantTags:     []string{"android"},
		},
		{
			desc: "feature disabled dispatches",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				Type:          domaineventproto.Event_FEATURE_DISABLED,
				EnvironmentId: envID,
				EntityData:    featureTagsJSON(t, "android"),
			},
			clientTag:    "android",
			wantDispatch: true,
			wantTags:     []string{"android"},
		},
		{
			desc: "feature created is ignored",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				Type:          domaineventproto.Event_FEATURE_CREATED,
				EnvironmentId: envID,
				EntityData:    featureTagsJSON(t, "android"),
			},
			clientTag:    "android",
			wantDispatch: false,
		},
		{
			desc: "feature archived is ignored",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				Type:          domaineventproto.Event_FEATURE_ARCHIVED,
				EnvironmentId: envID,
				EntityData:    featureTagsJSON(t, "android"),
			},
			clientTag:    "android",
			wantDispatch: false,
		},
		{
			desc: "segment succeeded upload dispatches env-wide",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_SEGMENT,
				Type:          domaineventproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED,
				EnvironmentId: envID,
				Data:          segmentStatusData(t, featureproto.Segment_SUCEEDED),
			},
			clientTag:    "android",
			wantDispatch: true,
			wantTags:     nil,
		},
		{
			desc: "segment upload still uploading is ignored",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_SEGMENT,
				Type:          domaineventproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED,
				EnvironmentId: envID,
				Data:          segmentStatusData(t, featureproto.Segment_UPLOADING),
			},
			clientTag: "android",
		},
		{
			desc: "segment upload failed is ignored",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_SEGMENT,
				Type:          domaineventproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED,
				EnvironmentId: envID,
				Data:          segmentStatusData(t, featureproto.Segment_FAILED),
			},
			clientTag: "android",
		},
		{
			desc: "segment non-bulk-upload type is ignored",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_SEGMENT,
				Type:          domaineventproto.Event_SEGMENT_CREATED,
				EnvironmentId: envID,
				Data:          segmentStatusData(t, featureproto.Segment_SUCEEDED),
			},
			clientTag: "android",
		},
		{
			desc: "segment event with nil data is ignored",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_SEGMENT,
				Type:          domaineventproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED,
				EnvironmentId: envID,
				Data:          nil,
			},
			clientTag: "android",
		},
		{
			desc: "unknown entity type is ignored",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_GOAL,
				Type:          domaineventproto.Event_GOAL_CREATED,
				EnvironmentId: envID,
			},
			clientTag: "android",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			dispatcher := NewStreamDispatcher()
			ch, cancel := dispatcher.register(envID, p.clientTag)
			defer cancel()
			s := newTestStreamSubscriber(t, dispatcher)

			data, err := proto.Marshal(p.event)
			require.NoError(t, err)
			s.handleMessage(&puller.Message{Data: data})

			select {
			case got := <-ch:
				require.True(t, p.wantDispatch, "unexpected dispatch")
				assert.Equal(t, envID, got.environmentID)
				assert.Equal(t, p.wantTags, got.tags)
			default:
				require.False(t, p.wantDispatch, "expected dispatch but none occurred")
			}
		})
	}
}

func TestHandleMessageIgnoresMalformedData(t *testing.T) {
	t.Parallel()
	dispatcher := NewStreamDispatcher()
	ch, cancel := dispatcher.register("env-1", "android")
	defer cancel()
	s := newTestStreamSubscriber(t, dispatcher)

	// Truncated varint: not a decodable domain event. Must be dropped without
	// panicking or dispatching.
	s.handleMessage(&puller.Message{Data: []byte{0xff, 0xff, 0xff}})

	select {
	case <-ch:
		t.Fatal("malformed message must not dispatch")
	default:
	}
}

func segmentStatusData(t *testing.T, status featureproto.Segment_Status) *anypb.Any {
	t.Helper()
	data, err := anypb.New(&domaineventproto.SegmentBulkUploadUsersStatusChangedEvent{
		SegmentId: "seg-1",
		Status:    status,
	})
	require.NoError(t, err)
	return data
}

func TestAffectedTags(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		entity   string
		previous string
		expected []string
	}{
		{
			desc:     "both empty",
			entity:   "",
			previous: "",
			expected: nil,
		},
		{
			desc:     "only current",
			entity:   featureTagsJSON(t, "android", "ios"),
			previous: "",
			expected: []string{"android", "ios"},
		},
		{
			desc:     "only previous",
			entity:   "",
			previous: featureTagsJSON(t, "android"),
			expected: []string{"android"},
		},
		{
			desc:     "union of added and removed tags",
			entity:   featureTagsJSON(t, "android"),
			previous: featureTagsJSON(t, "android", "ios"),
			expected: []string{"android", "ios"},
		},
		{
			desc:     "deduplicates overlapping tags",
			entity:   featureTagsJSON(t, "android", "ios"),
			previous: featureTagsJSON(t, "ios", "web"),
			expected: []string{"android", "ios", "web"},
		},
		{
			desc:     "invalid json yields no tags",
			entity:   "{not json",
			previous: "",
			expected: nil,
		},
	}
	s := &StreamSubscriber{logger: zap.NewNop()}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			got := s.affectedTags(&domaineventproto.Event{
				EntityData:         p.entity,
				PreviousEntityData: p.previous,
			})
			sort.Strings(got)
			sort.Strings(p.expected)
			assert.Equal(t, p.expected, got)
		})
	}
}

func featureTagsJSON(t *testing.T, tags ...string) string {
	t.Helper()
	b, err := json.Marshal(&featureproto.Feature{Tags: tags})
	require.NoError(t, err)
	return string(b)
}
