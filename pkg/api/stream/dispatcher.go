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
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var errTooManyConnections = errors.New("stream: too many connections")

// FeaturesFetcher returns all features for the given environment.
type FeaturesFetcher func(envID string) ([]*featureproto.Feature, error)

// Dispatcher forwards relevant domain events to SSE connections.
type Dispatcher struct {
	mu sync.Mutex
	// envID -> tag -> set of conns
	conns           map[string]map[string]map[*conn]struct{}
	totalConns      int
	maxConns        int
	fetchFeatures   FeaturesFetcher
	refetchFeatures FeaturesFetcher
	shutdownCh      chan struct{}
	shutdownOnce    sync.Once
	logger          *zap.Logger
}

type DispatcherOption func(*Dispatcher)

// WithFeaturesRefetcher sets the source used when the cache is older than the event.
func WithFeaturesRefetcher(f FeaturesFetcher) DispatcherOption {
	return func(d *Dispatcher) {
		d.refetchFeatures = f
	}
}

type event struct {
	environmentID string
	tags          []string
	eventType     domaineventproto.Event_Type
	// features is a snapshot verified to include the event's change; nil means use the cache.
	features     []*featureproto.Feature
	dispatchedAt time.Time
}

type conn struct {
	ch        chan event
	tag       string
	sourceID  string
	createdAt time.Time
}

func NewDispatcher(
	maxConns int,
	fetchFeatures FeaturesFetcher,
	logger *zap.Logger,
	opts ...DispatcherOption,
) *Dispatcher {
	d := &Dispatcher{
		conns:         make(map[string]map[string]map[*conn]struct{}),
		maxConns:      maxConns,
		fetchFeatures: fetchFeatures,
		shutdownCh:    make(chan struct{}),
		logger:        logger.Named("stream-dispatcher"),
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// ActiveConns returns the current number of SSE connections.
func (d *Dispatcher) ActiveConns() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.totalConns
}

// Shutdown signals all active SSE handlers to exit immediately.
func (d *Dispatcher) Shutdown() {
	d.shutdownOnce.Do(func() { close(d.shutdownCh) })
}

// register adds a connection to the dispatcher. The caller must invoke the returned
// deregister func on disconnect to free the slot.
// Returns errTooManyConnections when maxConns is set and already reached.
func (d *Dispatcher) register(envID, tag, sourceID string) (events <-chan event, deregister func(), err error) {
	d.mu.Lock()
	if d.totalConns >= d.maxConns {
		d.mu.Unlock()
		sseErrorsCounter.WithLabelValues(envID, tag, sourceID, errorTypeConnectionRefusedByLimit).Inc()
		return nil, nil, errTooManyConnections
	}
	c := &conn{
		ch:        make(chan event, 1),
		tag:       tag,
		sourceID:  sourceID,
		createdAt: time.Now(),
	}
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
	d.totalConns++
	// Update the gauge inside the lock so it stays consistent with the conn map.
	sseActiveConnectionsGauge.WithLabelValues(envID, tag, sourceID).Inc()
	d.mu.Unlock()

	return c.ch, func() { d.deregister(envID, c) }, nil
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
	d.totalConns--
	sseActiveConnectionsGauge.WithLabelValues(envID, target.tag, target.sourceID).Dec()
	if len(conns) == 0 {
		delete(tagConns, target.tag)
	}
	if len(tagConns) == 0 {
		delete(d.conns, envID)
	}
	d.mu.Unlock()

	sseConnectionDurationHistogram.WithLabelValues(envID, target.tag, target.sourceID).
		Observe(time.Since(target.createdAt).Seconds())

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
		features := d.featuresFor(e)
		d.dispatch(event{
			environmentID: e.EnvironmentId,
			eventType:     e.Type,
			tags:          d.affectedTags(e, features),
			features:      features,
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
			eventType:     e.Type,
		})
	}
}

// affectedTags unions the flag's tags before and after the update so that
// removing a tag still notifies that tag's subscribers.
// And the result includes the tags of flags that transitively depend on this one.
func (d *Dispatcher) affectedTags(e *domaineventproto.Event, features []*featureproto.Feature) []string {
	seen := make(map[string]struct{})
	for _, data := range []string{e.EntityData, e.PreviousEntityData} {
		for _, tag := range d.parseTags(data) {
			seen[tag] = struct{}{}
		}
	}
	for _, tag := range dependentTags(e.EntityId, features) {
		seen[tag] = struct{}{}
	}
	if len(seen) == 0 {
		return nil
	}
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	return tags
}

// featuresFor returns the env's features, refetching them if the cache predates the event.
func (d *Dispatcher) featuresFor(e *domaineventproto.Event) []*featureproto.Feature {
	if d.fetchFeatures == nil {
		return nil
	}
	features, err := d.fetchFeatures(e.EnvironmentId)
	if err != nil {
		d.logger.Warn("Failed to fetch features",
			zap.Error(err),
			zap.String("environmentID", e.EnvironmentId),
		)
		features = nil
	}
	version, ok := d.parseVersion(e.EntityData)
	if !ok || d.refetchFeatures == nil || hasVersion(features, e.EntityId, version) {
		return features
	}
	sseStaleFeaturesCounter.WithLabelValues(e.EnvironmentId).Inc()
	fresh, err := d.refetchFeatures(e.EnvironmentId)
	if err != nil {
		d.logger.Warn("Failed to refetch stale features",
			zap.Error(err),
			zap.String("environmentID", e.EnvironmentId),
			zap.String("featureID", e.EntityId),
		)
		return features
	}
	return fresh
}

func hasVersion(features []*featureproto.Feature, id string, version int32) bool {
	for _, f := range features {
		if f.Id == id {
			return f.Version >= version
		}
	}
	return false
}

func dependentTags(entityID string, features []*featureproto.Feature) []string {
	if len(features) == 0 {
		return nil
	}
	featuresMap := make(map[string]*featureproto.Feature, len(features))
	for _, f := range features {
		featuresMap[f.Id] = f
	}
	target := featuresMap[entityID]
	if target == nil {
		return nil
	}
	var tags []string
	for _, f := range ftdomain.GetDependentsOfTargets(
		[]*featureproto.Feature{target}, featuresMap,
	) {
		tags = append(tags, f.Tags...)
	}
	return tags
}

func (d *Dispatcher) parseVersion(data string) (int32, bool) {
	if data == "" {
		return 0, false
	}
	var payload struct {
		Version *int32 `json:"version"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil || payload.Version == nil {
		return 0, false
	}
	return *payload.Version, true
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
	var dispatchTagCount float64
	if len(ev.tags) == 0 {
		dispatchTagCount = float64(len(tagConns))
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
		dispatchTagCount = float64(len(ev.tags))
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

	sseDispatchTagsHistogram.WithLabelValues(ev.environmentID, ev.eventType.String()).
		Observe(dispatchTagCount)

	ev.dispatchedAt = time.Now()
	for _, c := range targetConns {
		select {
		case c.ch <- ev:
			continue
		default:
		}
		// Replace the pending event so the conn evaluates the newest snapshot.
		sseDispatchDroppedCounter.WithLabelValues(ev.environmentID, c.tag, c.sourceID).Inc()
		next := ev
		select {
		case pending := <-c.ch:
			if next.features == nil {
				next.features = pending.features
			}
		default:
		}
		select {
		case c.ch <- next:
		default:
		}
	}
}
