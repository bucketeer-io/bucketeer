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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testConnSpec struct {
	envID string
	tag   string
}

// snapshotConnCounts captures d.conns as env -> tag -> conn count under lock so tests
// can assert the internal shape without races.
func snapshotConnCounts(d *StreamDispatcher) map[string]map[string]int {
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

func TestStreamDispatcherRegister(t *testing.T) {
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
			d := NewStreamDispatcher()
			for _, r := range tc.clients {
				_, cancel := d.register(r.envID, r.tag)
				defer cancel()
			}
			assert.Equal(t, tc.wantConnCounts, snapshotConnCounts(d))
		})
	}
}

func TestStreamDispatcherDeregister(t *testing.T) {
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
			d := NewStreamDispatcher()
			cancels := make(map[testConnSpec]func())
			for env, tagConns := range tc.conns {
				for tag := range tagConns {
					for i := 0; i < tagConns[tag]; i++ {
						_, cancels[testConnSpec{env, tag}] = d.register(env, tag)
					}
				}
			}
			cancels[tc.deregister]()
			assert.Equal(t, tc.wantConnCounts, snapshotConnCounts(d))
		})
	}
}

func TestStreamDispatcherDispatch(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		conns []testConnSpec
		event streamEvent
		// expected: parallel to conns; true if the conn at the same index should receive the event
		wantRecv []bool
	}{
		{
			name:     "matches env+tag",
			conns:    []testConnSpec{{"env-1", "tag-A"}},
			event:    streamEvent{environmentID: "env-1", tags: []string{"tag-A"}},
			wantRecv: []bool{true},
		},
		{
			name:     "other env not reached",
			conns:    []testConnSpec{{"env-1", "tag-A"}, {"env-2", "tag-A"}},
			event:    streamEvent{environmentID: "env-1", tags: []string{"tag-A"}},
			wantRecv: []bool{true, false},
		},
		{
			name:     "empty tags fan out within env",
			conns:    []testConnSpec{{"env-1", "tag-A"}, {"env-1", "tag-B"}},
			event:    streamEvent{environmentID: "env-1"},
			wantRecv: []bool{true, true},
		},
		{
			name:     "other tag not reached",
			conns:    []testConnSpec{{"env-1", "tag-A"}, {"env-1", "tag-B"}},
			event:    streamEvent{environmentID: "env-1", tags: []string{"tag-A"}},
			wantRecv: []bool{true, false},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := NewStreamDispatcher()
			chs := make([]<-chan streamEvent, len(tc.conns))
			for i, c := range tc.conns {
				ch, cancel := d.register(c.envID, c.tag)
				defer cancel()
				chs[i] = ch
			}

			d.dispatch(tc.event)

			for i, c := range tc.conns {
				if tc.wantRecv[i] {
					select {
					case got := <-chs[i]:
						assert.Equal(t, tc.event, got)
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
