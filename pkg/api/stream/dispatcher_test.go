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

package stream

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type testConnSpec struct {
	envID string
	tag   string
}

// snapshotConnCounts captures d.conns as env -> tag -> conn count under lock so tests
// can assert the internal shape without races.
func snapshotConnCounts(d *Dispatcher) map[string]map[string]int {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := map[string]map[string]int{}
	for env, tagConns := range d.conns {
		inner := map[string]int{}
		for tag, conns := range tagConns {
			inner[tag] = len(conns)
		}
		out[env] = inner
	}
	return out
}

func TestDispatcherRegister(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		clients        []testConnSpec
		wantConnCounts map[string]map[string]int
	}{
		{
			name:           "single conn",
			clients:        []testConnSpec{{"env-1", "tag-A"}},
			wantConnCounts: map[string]map[string]int{"env-1": {"tag-A": 1}},
		},
		{
			name:           "multiple conns same env+tag",
			clients:        []testConnSpec{{"env-1", "tag-A"}, {"env-1", "tag-A"}},
			wantConnCounts: map[string]map[string]int{"env-1": {"tag-A": 2}},
		},
		{
			name:           "multiple tags same env",
			clients:        []testConnSpec{{"env-1", "tag-A"}, {"env-1", "tag-B"}},
			wantConnCounts: map[string]map[string]int{"env-1": {"tag-A": 1, "tag-B": 1}},
		},
		{
			name:    "multiple envs",
			clients: []testConnSpec{{"env-1", "tag-A"}, {"env-2", "tag-A"}},
			wantConnCounts: map[string]map[string]int{
				"env-1": {"tag-A": 1},
				"env-2": {"tag-A": 1},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := NewDispatcher(zap.NewNop())
			for _, r := range tc.clients {
				_, cancel := d.register(r.envID, r.tag, "source1")
				defer cancel()
			}
			assert.Equal(t, tc.wantConnCounts, snapshotConnCounts(d))
		})
	}
}

func TestDispatcherDeregister(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		conns          map[string]map[string]int
		deregister     testConnSpec
		wantConnCounts map[string]map[string]int
	}{
		{
			name:           "deregister single conn empties map",
			conns:          map[string]map[string]int{"env-1": {"tag-A": 1}},
			deregister:     testConnSpec{"env-1", "tag-A"},
			wantConnCounts: map[string]map[string]int{},
		},
		{
			name:           "deregister one of two same-key conns leaves other",
			conns:          map[string]map[string]int{"env-1": {"tag-A": 2}},
			deregister:     testConnSpec{"env-1", "tag-A"},
			wantConnCounts: map[string]map[string]int{"env-1": {"tag-A": 1}},
		},
		{
			name:           "emptying a tag removes the tag entry",
			conns:          map[string]map[string]int{"env-1": {"tag-A": 1, "tag-B": 1}},
			deregister:     testConnSpec{"env-1", "tag-A"},
			wantConnCounts: map[string]map[string]int{"env-1": {"tag-B": 1}},
		},
		{
			name:           "emptying an env removes the env entry",
			conns:          map[string]map[string]int{"env-1": {"tag-A": 1}, "env-2": {"tag-A": 1}},
			deregister:     testConnSpec{"env-1", "tag-A"},
			wantConnCounts: map[string]map[string]int{"env-2": {"tag-A": 1}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := NewDispatcher(zap.NewNop())
			cancels := make(map[testConnSpec]func())
			for env, tagConns := range tc.conns {
				for tag := range tagConns {
					for i := 0; i < tagConns[tag]; i++ {
						var cancel func()
						_, cancel = d.register(env, tag, "source1")
						defer cancel()
						cancels[testConnSpec{env, tag}] = cancel
					}
				}
			}
			cancels[tc.deregister]()
			assert.Equal(t, tc.wantConnCounts, snapshotConnCounts(d))
		})
	}
}

func TestDispatcherDispatch(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		conns []testConnSpec
		event event
		// expected: parallel to conns; true if the conn at the same index should receive the event
		wantRecv []bool
	}{
		{
			name:     "matches env+tag",
			conns:    []testConnSpec{{"env-1", "tag-A"}},
			event:    event{environmentID: "env-1", tags: []string{"tag-A"}},
			wantRecv: []bool{true},
		},
		{
			name:     "other env not reached",
			conns:    []testConnSpec{{"env-1", "tag-A"}, {"env-2", "tag-A"}},
			event:    event{environmentID: "env-1", tags: []string{"tag-A"}},
			wantRecv: []bool{true, false},
		},
		{
			name:     "empty tags fan out within env",
			conns:    []testConnSpec{{"env-1", "tag-A"}, {"env-1", "tag-B"}},
			event:    event{environmentID: "env-1"},
			wantRecv: []bool{true, true},
		},
		{
			name:     "other tag not reached",
			conns:    []testConnSpec{{"env-1", "tag-A"}, {"env-1", "tag-B"}},
			event:    event{environmentID: "env-1", tags: []string{"tag-A"}},
			wantRecv: []bool{true, false},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := NewDispatcher(zap.NewNop())
			chs := make([]<-chan event, len(tc.conns))
			for i, c := range tc.conns {
				ch, cancel := d.register(c.envID, c.tag, "source1")
				defer cancel()
				chs[i] = ch
			}

			d.dispatch(tc.event)

			for i, c := range tc.conns {
				if tc.wantRecv[i] {
					select {
					case got := <-chs[i]:
						assert.Equal(t, tc.event.environmentID, got.environmentID)
						assert.Equal(t, tc.event.tags, got.tags)
						assert.False(t, got.dispatchedAt.IsZero())
					case <-time.After(time.Second):
						t.Fatalf("conn[%d] (env=%s tag=%s) expected event", i, c.envID, c.tag)
					}
					continue
				}
				select {
				case <-chs[i]:
					t.Fatalf("conn[%d] (env=%s tag=%s) must not receive event", i, c.envID, c.tag)
				case <-time.After(50 * time.Millisecond):
				}
			}
		})
	}
}

func TestDispatcherHandleEvent(t *testing.T) {
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
			d := NewDispatcher(zap.NewNop())
			ch, cancel := d.register(envID, p.clientTag, "source1")
			defer cancel()

			d.HandleEvent(p.event)

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

func segmentStatusData(t *testing.T, status featureproto.Segment_Status) *anypb.Any {
	t.Helper()
	data, err := anypb.New(&domaineventproto.SegmentBulkUploadUsersStatusChangedEvent{
		SegmentId: "seg-1",
		Status:    status,
	})
	require.NoError(t, err)
	return data
}

func TestDispatcherAffectedTags(t *testing.T) {
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
	d := &Dispatcher{logger: zap.NewNop()}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			got := d.affectedTags(&domaineventproto.Event{
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
