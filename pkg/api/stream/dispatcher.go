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

const (
	maxRefetchAttempts          = 3
	defaultRefetchRetryInterval = time.Second
)

// FeaturesFetcher returns all features for the given environment.
type FeaturesFetcher func(envID string) ([]*featureproto.Feature, error)

// Dispatcher forwards relevant domain events to SSE connections.
type Dispatcher struct {
	mu sync.Mutex
	// envID -> tag -> set of conns
	conns           map[string]map[string]map[*conn]struct{}
	pending         map[string]*pendingDispatch
	draining        map[string]bool
	totalConns      int
	maxConns        int
	fetchFeatures   FeaturesFetcher
	refetchFeatures FeaturesFetcher
	// Base backoff between refetch attempts; doubled on each retry.
	refetchRetryInterval time.Duration
	shutdownCh           chan struct{}
	shutdownOnce         sync.Once
	logger               *zap.Logger
}

type DispatcherOption func(*Dispatcher)

// WithFeaturesRefetcher sets the source used when the cache is older than the event.
func WithFeaturesRefetcher(f FeaturesFetcher) DispatcherOption {
	return func(d *Dispatcher) {
		d.refetchFeatures = f
	}
}

// pendingDispatch merges an environment's events until its drain goroutine processes them.
type pendingDispatch struct {
	eventType domaineventproto.Event_Type
	allTags   bool
	// featureID -> change
	changes  map[string]*featureChange
	attempts int
}

type featureChange struct {
	tags       []string
	version    int32
	hasVersion bool
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
		conns:                make(map[string]map[string]map[*conn]struct{}),
		pending:              make(map[string]*pendingDispatch),
		draining:             make(map[string]bool),
		maxConns:             maxConns,
		fetchFeatures:        fetchFeatures,
		refetchRetryInterval: defaultRefetchRetryInterval,
		shutdownCh:           make(chan struct{}),
		logger:               logger.Named("stream-dispatcher"),
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
		tags := d.changeTags(e)
		version, hasVersion := d.parseVersion(e.EntityData)
		d.enqueue(e.EnvironmentId, e.Type, func(p *pendingDispatch) {
			p.addChange(e.EntityId, &featureChange{tags: tags, version: version, hasVersion: hasVersion})
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
		d.enqueue(e.EnvironmentId, e.Type, func(p *pendingDispatch) {
			p.allTags = true
		})
	}
}

func (p *pendingDispatch) addChange(id string, change *featureChange) {
	c, ok := p.changes[id]
	if !ok {
		c = &featureChange{}
		p.changes[id] = c
	}
	c.tags = append(c.tags, change.tags...)
	if change.hasVersion && (!c.hasVersion || change.version > c.version) {
		c.version, c.hasVersion = change.version, true
	}
}

// enqueue merges an event into the env's pending dispatch so refetches never block the caller.
func (d *Dispatcher) enqueue(envID string, eventType domaineventproto.Event_Type, merge func(*pendingDispatch)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Pods without connections in the environment have nothing to fan out.
	if len(d.conns[envID]) == 0 {
		return
	}
	p, ok := d.pending[envID]
	if !ok {
		p = &pendingDispatch{changes: make(map[string]*featureChange)}
		d.pending[envID] = p
	}
	p.eventType = eventType
	merge(p)
	if d.draining[envID] {
		return
	}
	d.draining[envID] = true
	go d.drain(envID)
}

// drain processes the environment's pending dispatches one at a time, preserving their order.
func (d *Dispatcher) drain(envID string) {
	for {
		d.mu.Lock()
		p, ok := d.pending[envID]
		if !ok {
			delete(d.draining, envID)
			d.mu.Unlock()
			return
		}
		delete(d.pending, envID)
		d.mu.Unlock()
		if d.process(envID, p) {
			continue
		}
		// Never send a snapshot known to be stale: retry, then give up.
		p.attempts++
		if p.attempts >= maxRefetchAttempts {
			d.logger.Error("Dropped dispatch after failing to refetch stale features",
				zap.String("environmentID", envID),
				zap.Int("attempts", p.attempts),
			)
			continue
		}
		d.requeue(envID, p)
		select {
		case <-time.After(d.refetchRetryInterval << (p.attempts - 1)):
		case <-d.shutdownCh:
			d.mu.Lock()
			delete(d.pending, envID)
			delete(d.draining, envID)
			d.mu.Unlock()
			return
		}
	}
}

// requeue merges a failed dispatch with any events that arrived meanwhile.
func (d *Dispatcher) requeue(envID string, p *pendingDispatch) {
	d.mu.Lock()
	defer d.mu.Unlock()
	cur, ok := d.pending[envID]
	if !ok {
		d.pending[envID] = p
		return
	}
	cur.allTags = cur.allTags || p.allTags
	cur.attempts = max(cur.attempts, p.attempts)
	for id, c := range p.changes {
		cur.addChange(id, c)
	}
}

// process dispatches the pending events and returns false if the refetch failed.
func (d *Dispatcher) process(envID string, p *pendingDispatch) bool {
	ev := event{
		environmentID: envID,
		eventType:     p.eventType,
	}
	if len(p.changes) > 0 {
		features, ok := d.featuresFor(envID, p.changes)
		if !ok {
			return false
		}
		ev.features = features
	}
	if !p.allTags {
		ev.tags = affectedTags(p.changes, ev.features)
	}
	d.dispatch(ev)
	return true
}

// changeTags unions the flag's tags before and after the update so that
// removing a tag still notifies that tag's subscribers.
func (d *Dispatcher) changeTags(e *domaineventproto.Event) []string {
	var tags []string
	for _, data := range []string{e.EntityData, e.PreviousEntityData} {
		tags = append(tags, d.parseTags(data)...)
	}
	return tags
}

// affectedTags unions the changed flags' tags and the tags of flags that transitively depend on them.
// It returns nil (all tags) if any changed flag has no tags to target.
func affectedTags(changes map[string]*featureChange, features []*featureproto.Feature) []string {
	seen := make(map[string]struct{})
	for id, c := range changes {
		dependents := dependentTags(id, features)
		if len(c.tags) == 0 && len(dependents) == 0 {
			return nil
		}
		for _, tag := range c.tags {
			seen[tag] = struct{}{}
		}
		for _, tag := range dependents {
			seen[tag] = struct{}{}
		}
	}
	if len(seen) == 0 {
		return nil
	}
	out := make([]string, 0, len(seen))
	for tag := range seen {
		out = append(out, tag)
	}
	return out
}

// featuresFor returns the env's features, refetching them if the cache predates any change.
// It returns false if the cache is stale and the refetch failed.
func (d *Dispatcher) featuresFor(
	envID string,
	changes map[string]*featureChange,
) ([]*featureproto.Feature, bool) {
	if d.fetchFeatures == nil {
		return nil, true
	}
	features, err := d.fetchFeatures(envID)
	if err != nil {
		d.logger.Warn("Failed to fetch features",
			zap.Error(err),
			zap.String("environmentID", envID),
		)
		features = nil
	}
	if d.refetchFeatures == nil || includesChanges(features, changes) {
		return features, true
	}
	sseStaleFeaturesCounter.WithLabelValues(envID).Inc()
	fresh, err := d.refetchFeatures(envID)
	if err != nil {
		d.logger.Warn("Failed to refetch stale features",
			zap.Error(err),
			zap.String("environmentID", envID),
		)
		return nil, false
	}
	return fresh, true
}

func includesChanges(features []*featureproto.Feature, changes map[string]*featureChange) bool {
	versions := make(map[string]int32, len(features))
	for _, f := range features {
		versions[f.Id] = f.Version
	}
	for id, c := range changes {
		if !c.hasVersion {
			continue
		}
		if v, ok := versions[id]; !ok || v < c.version {
			return false
		}
	}
	return true
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
	// Flags excluded from every features cache can never be found there, so skip the check.
	f := &featureproto.Feature{}
	if err := json.Unmarshal([]byte(data), f); err == nil {
		ff := ftdomain.Feature{Feature: f}
		if ff.IsDisabledAndOffVariationEmpty() || ff.IsArchivedBeforeLastThirtyDays() {
			return 0, false
		}
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
