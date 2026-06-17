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
	"sync"

	"go.uber.org/zap"

	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// Dispatcher forwards relevant domain events to SSE connections.
type Dispatcher struct {
	mu sync.Mutex
	// envID -> tag -> set of conns
	conns  map[string]map[string]map[*conn]struct{}
	logger *zap.Logger
}

type event struct {
	environmentID string
	tags          []string
}

type conn struct {
	ch  chan event
	tag string
}

func NewDispatcher(logger *zap.Logger) *Dispatcher {
	return &Dispatcher{
		conns:  make(map[string]map[string]map[*conn]struct{}),
		logger: logger.Named("stream-dispatcher"),
	}
}

// register adds a connection to the dispatcher. The caller must invoke the returned
// deregister func on disconnect to free the slot.
func (d *Dispatcher) register(envID, tag string) (events <-chan event, deregister func()) {
	c := &conn{
		// Only the latest event is needed to evaluate with the latest config.
		ch:  make(chan event, 1),
		tag: tag,
	}
	d.mu.Lock()
	tagConns, ok := d.conns[envID]
	if !ok {
		tagConns = make(map[string]map[*conn]struct{})
		d.conns[envID] = tagConns
	}
	conns, ok := tagConns[tag]
	if !ok {
		conns = make(map[*conn]struct{})
		tagConns[tag] = conns
	}
	conns[c] = struct{}{}
	// Update the gauge inside the lock so a concurrent last-conn deregister
	// cannot delete the series after this Inc.
	sseActiveConnectionsGauge.WithLabelValues(envID, tag).Inc()
	d.mu.Unlock()

	return c.ch, func() { d.deregister(envID, c) }
}

func (d *Dispatcher) deregister(envID string, target *conn) {
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

// HandleEvent dispatches a domain event to the affected connections if it can
// change evaluation results.
func (d *Dispatcher) HandleEvent(e *domaineventproto.Event) {
	switch e.EntityType {
	case domaineventproto.Event_FEATURE:
		if e.Type != domaineventproto.Event_FEATURE_UPDATED &&
			e.Type != domaineventproto.Event_FEATURE_ENABLED &&
			e.Type != domaineventproto.Event_FEATURE_DISABLED {
			return
		}
		d.dispatch(event{
			environmentID: e.EnvironmentId,
			tags:          d.affectedTags(e),
		})
	case domaineventproto.Event_SEGMENT:
		// Segment membership only changes through a bulk upload.
		if e.Type != domaineventproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED {
			return
		}
		if e.Data == nil {
			return
		}
		payload := &domaineventproto.SegmentBulkUploadUsersStatusChangedEvent{}
		if err := e.Data.UnmarshalTo(payload); err != nil {
			d.logger.Warn("Failed to unmarshal segment bulk upload event", zap.Error(err))
			return
		}
		// Membership changes only on a succeeded upload.
		if payload.Status != featureproto.Segment_SUCEEDED {
			return
		}
		// TODO: resolve the affected tags from the segment.
		// Currently, it fans out env-wide (all tags).
		d.dispatch(event{
			environmentID: e.EnvironmentId,
		})
	}
}

// affectedTags unions the flag's tags before and after the update so that
// removing a tag still notifies that tag's subscribers.
//
// TODO: also pull the tags of flags that depend on this one.
func (d *Dispatcher) affectedTags(e *domaineventproto.Event) []string {
	seen := make(map[string]struct{})
	var tags []string
	for _, data := range []string{e.EntityData, e.PreviousEntityData} {
		for _, tag := range d.parseTags(data) {
			if _, ok := seen[tag]; ok {
				continue
			}
			seen[tag] = struct{}{}
			tags = append(tags, tag)
		}
	}
	return tags
}

func (d *Dispatcher) parseTags(data string) []string {
	if data == "" {
		return nil
	}
	var payload struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		d.logger.Warn("Failed to extract tags from feature entity data", zap.Error(err))
		return nil
	}
	return payload.Tags
}

// dispatch fans an event out to matching tag connections in the environment, or to
// all of them when tags is empty.
// Sends are non-blocking.
func (d *Dispatcher) dispatch(ev event) {
	d.mu.Lock()
	tagConns := d.conns[ev.environmentID]
	if len(tagConns) == 0 {
		d.mu.Unlock()
		return
	}

	// Snapshot conns to release the lock early to avoid blocking register/deregister.
	var targetConns []*conn
	if len(ev.tags) == 0 {
		n := 0
		for _, conns := range tagConns {
			n += len(conns)
		}
		targetConns = make([]*conn, 0, n)
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
		targetConns = make([]*conn, 0, n)
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
