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
	"sync"
)

// StreamDispatcher is the mapping of (environmentID, tag) to SSE connections.
type StreamDispatcher struct {
	mu sync.Mutex
	// envID -> tag -> set of conns
	conns map[string]map[string]map[*streamConn]struct{}
}

type streamEvent struct {
	environmentID string
	tags          []string
}

type streamConn struct {
	ch  chan streamEvent
	tag string
}

func NewStreamDispatcher() *StreamDispatcher {
	return &StreamDispatcher{
		conns: make(map[string]map[string]map[*streamConn]struct{}),
	}
}

// register adds a connection to the dispatcher. The caller must invoke the returned
// deregister func on disconnect to free the slot.
func (d *StreamDispatcher) register(envID, tag string) (events <-chan streamEvent, deregister func()) {
	c := &streamConn{
		// Only the latest event is needed to evaluate with the latest config.
		ch:  make(chan streamEvent, 1),
		tag: tag,
	}
	d.mu.Lock()
	tagConns, ok := d.conns[envID]
	if !ok {
		tagConns = make(map[string]map[*streamConn]struct{})
		d.conns[envID] = tagConns
	}
	conns, ok := tagConns[tag]
	if !ok {
		conns = make(map[*streamConn]struct{})
		tagConns[tag] = conns
	}
	conns[c] = struct{}{}
	// Update the gauge inside the lock so a concurrent last-conn deregister
	// cannot delete the series after this Inc.
	sseActiveConnectionsGauge.WithLabelValues(envID, tag).Inc()
	d.mu.Unlock()

	return c.ch, func() { d.deregister(envID, c) }
}

func (d *StreamDispatcher) deregister(envID string, target *streamConn) {
	d.mu.Lock()
	tagConns, ok := d.conns[envID]
	if !ok {
		d.mu.Unlock()
		return
	}
	conns := tagConns[target.tag]
	if _, ok := conns[target]; !ok {
		d.mu.Unlock()
		return
	}
	delete(conns, target)
	if len(conns) == 0 {
		delete(tagConns, target.tag)
		sseActiveConnectionsGauge.DeleteLabelValues(envID, target.tag)
	} else {
		sseActiveConnectionsGauge.WithLabelValues(envID, target.tag).Dec()
	}
	if len(tagConns) == 0 {
		delete(d.conns, envID)
	}
	d.mu.Unlock()

	// target.ch is intentionally not closed because closing it here would
	// race with dispatch and panic. The handler exits via ctx.Done().
}

// dispatch fans an event out to matching tag connections in the environment, or to
// all of them when tags is empty.
// Sends are non-blocking.
func (d *StreamDispatcher) dispatch(ev streamEvent) {
	d.mu.Lock()
	tagConns := d.conns[ev.environmentID]
	if len(tagConns) == 0 {
		d.mu.Unlock()
		return
	}

	// Snapshot conns to release the lock early to avoid blocking register/deregister.
	var targetConns []*streamConn
	if len(ev.tags) == 0 {
		n := 0
		for _, conns := range tagConns {
			n += len(conns)
		}
		targetConns = make([]*streamConn, 0, n)
		for _, conns := range tagConns {
			for c := range conns {
				targetConns = append(targetConns, c)
			}
		}
	} else {
		n := 0
		for _, t := range ev.tags {
			n += len(tagConns[t])
		}
		targetConns = make([]*streamConn, 0, n)
		for _, t := range ev.tags {
			for c := range tagConns[t] {
				targetConns = append(targetConns, c)
			}
		}
	}
	d.mu.Unlock()

	for _, c := range targetConns {
		select {
		case c.ch <- ev:
		default:
			sseDispatchDroppedCounter.WithLabelValues(ev.environmentID, c.tag).Inc()
		}
	}
}
