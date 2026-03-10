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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	ftdomain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	eventCountKey      = "ec"
	userCountKey       = "uc"
	defaultVariationID = "default"
	userDataAppVersion = "app_version"
)

type lastUsedInfoCache map[string]*ftdomain.FeatureLastUsedInfo
type environmentLastUsedInfoCache map[string]lastUsedInfoCache
type eventMap map[string]*eventproto.EvaluationEvent
type environmentEventMap map[string]eventMap
type userAttributesCache map[string]*userproto.UserAttributes

// dauBufferKey groups DAU entries by date, environment, and source.
type dauBufferKey struct {
	dateStr  string // "20060102"
	envID    string
	sourceID string
}

// dauBuffer accumulates unique user IDs per (date, env, source).
type dauBuffer map[dauBufferKey]map[string]struct{}

type EvaluationCountEventPersisterConfig struct {
	FlushSize           int `json:"flushSize"`
	FlushInterval       int `json:"flushInterval"`
	WriteCacheInterval  int `json:"writeCacheInterval"`
	WriteDAUInterval    int `json:"writeDAUInterval"`
	UserAttributeKeyTTL int `json:"userAttributeKeyTtl"`
}

type evaluationCountEventPersister struct {
	evaluationCountEventPersisterConfig EvaluationCountEventPersisterConfig
	mysqlClient                         mysql.Client
	envLastUsedCache                    environmentLastUsedInfoCache
	evaluationCountCacher               cache.MultiGetDeleteCountCache
	userAttributesCacher                cachev3.UserAttributesCache
	dauCache                            cachev3.DAUCache
	dauBuf                              dauBuffer
	userAttributesCache                 userAttributesCache
	envLastUsedCacheMutex               sync.Mutex
	userAttributesCacheMutex            sync.Mutex
	dauBufferMutex                      sync.Mutex
	logger                              *zap.Logger
}

func NewEvaluationCountEventPersister(
	ctx context.Context,
	config interface{},
	mysqlClient mysql.Client,
	evaluationCountCacher cache.MultiGetDeleteCountCache,
	userAttributesCacher cachev3.UserAttributesCache,
	dauCache cachev3.DAUCache,
	logger *zap.Logger,
) (subscriber.PubSubProcessor, error) {
	evaluationCountEventPersisterJsonConfig, ok := config.(map[string]interface{})
	if !ok {
		logger.Error("EvaluationCountEventPersister: invalid config")
		return nil, ErrEvaluationCountInvalidConfig
	}
	configBytes, err := json.Marshal(evaluationCountEventPersisterJsonConfig)
	if err != nil {
		logger.Error("EvaluationCountEventPersisterJsonConfig: failed to marshal config", zap.Error(err))
		return nil, err
	}
	var evaluationCountEventPersisterConfig EvaluationCountEventPersisterConfig
	err = json.Unmarshal(configBytes, &evaluationCountEventPersisterConfig)
	if err != nil {
		logger.Error("EvaluationCountEventPersister: failed to unmarshal config", zap.Error(err))
		return nil, err
	}
	e := &evaluationCountEventPersister{
		evaluationCountEventPersisterConfig: evaluationCountEventPersisterConfig,
		mysqlClient:                         mysqlClient,
		envLastUsedCache:                    make(environmentLastUsedInfoCache),
		evaluationCountCacher:               evaluationCountCacher,
		userAttributesCacher:                userAttributesCacher,
		dauCache:                            dauCache,
		dauBuf:                              make(dauBuffer),
		userAttributesCache:                 make(userAttributesCache),
		envLastUsedCacheMutex:               sync.Mutex{},
		userAttributesCacheMutex:            sync.Mutex{},
		dauBufferMutex:                      sync.Mutex{},
		logger:                              logger,
	}
	// write flag last used info cache periodically
	//nolint:errcheck
	go e.writeFlagLastUsedInfoCache(ctx)
	// write user attributes cache periodically
	//nolint:errcheck
	go e.writeUserAttributesCache(ctx)
	// write DAU cache periodically
	//nolint:errcheck
	go e.writeDAUCache(ctx)
	return e, nil
}

func (p *evaluationCountEventPersister) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	batch := make(map[string]*puller.Message)
	ticker := time.NewTicker(time.Duration(p.evaluationCountEventPersisterConfig.FlushInterval) * time.Second)
	defer ticker.Stop()
	updateEvaluationCounter := func(envEvents environmentEventMap) {
		// Increment the evaluation event count in the Redis
		fails := p.incrementEnvEvents(envEvents)
		// Check to Ack or Nack the messages
		p.checkMessages(batch, fails)
		// Reset the maps and the timer
		batch = make(map[string]*puller.Message)
	}
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				p.logger.Error("Failed to pull message")
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberEvaluationCount).Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				// TODO: better log format for msg data
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationCount, codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				previous.Ack()
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationCount, codes.DuplicateID.String()).Inc()
			}
			batch[id] = msg
			if len(batch) < p.evaluationCountEventPersisterConfig.FlushSize {
				continue
			}
			envEvents := p.extractEvents(batch)
			// Update the feature flag last-used cache
			p.cacheLastUsedInfoPerEnv(envEvents)
			// Update the user attributes cache
			p.cacheUserAttributes(envEvents)
			// Buffer DAU entries
			p.bufferDAU(envEvents)
			updateEvaluationCounter(envEvents)
		case <-ticker.C:
			envEvents := p.extractEvents(batch)
			// Update the feature flag last-used cache
			p.cacheLastUsedInfoPerEnv(envEvents)
			// Update the user attributes cache
			p.cacheUserAttributes(envEvents)
			// Buffer DAU entries
			p.bufferDAU(envEvents)
			updateEvaluationCounter(envEvents)
		case <-ctx.Done():
			// Nack the messages to be redelivered
			for _, msg := range batch {
				msg.Nack()
			}
			p.logger.Debug("All the left messages were Nack successfully before shutting down",
				zap.Int("batchSize", len(batch)))
			return nil
		}
	}
}

func (p *evaluationCountEventPersister) incrementEnvEvents(envEvents environmentEventMap) map[string]bool {
	fails := make(map[string]bool, len(envEvents))
	aggregator := newEvaluationCountAggregator()

	type metricsKey struct {
		environmentId string
		sourceId      string
		featureId     string
		variationId   string
	}
	metricsAgg := make(map[metricsKey]int64)

	// Reverse mapping: ecKey → set of event IDs that contributed to it.
	// Migration ecKeys are excluded so migration-only failures don't cause retries.
	ecKeyToEventIDs := make(map[string]map[string]struct{})

	// Tracks per-metricsKey event counts for each ecKey, so we can correctly
	// subtract failed counts across multiple sourceIds that share an ecKey.
	ecKeyToMetricsCounts := make(map[string]map[metricsKey]int64)

	for environmentId, events := range envEvents {
		for id, e := range events {
			vID, err := getVariationID(e.Reason, e.VariationId)
			if err != nil {
				if errors.Is(err, ErrReasonNil) {
					fails[id] = false
				} else {
					fails[id] = true
				}
				continue
			}

			ucKey := p.newEvaluationCountkeyV2(userCountKey, e.FeatureId, vID, environmentId, e.Timestamp)
			ecKey := p.newEvaluationCountkeyV2(eventCountKey, e.FeatureId, vID, environmentId, e.Timestamp)
			userID := getUserID(e.UserId, e.User)

			aggregator.addEvent(ecKey, ucKey, userID)

			// Build reverse mapping (primary ecKeys only, not migration)
			if ecKeyToEventIDs[ecKey] == nil {
				ecKeyToEventIDs[ecKey] = make(map[string]struct{})
			}
			ecKeyToEventIDs[ecKey][id] = struct{}{}

			mKey := metricsKey{
				environmentId: environmentId,
				sourceId:      e.SourceId.String(),
				featureId:     e.FeatureId,
				variationId:   vID,
			}
			metricsAgg[mKey]++
			if ecKeyToMetricsCounts[ecKey] == nil {
				ecKeyToMetricsCounts[ecKey] = make(map[metricsKey]int64)
			}
			ecKeyToMetricsCounts[ecKey][mKey]++

			// Migration: best-effort double-write, not tracked in reverse mapping
			if targetEnvID := getMigrationTargetEnvironmentID(environmentId); targetEnvID != "" {
				ucKeyTarget := p.newEvaluationCountkeyV2(userCountKey, e.FeatureId, vID, targetEnvID, e.Timestamp)
				ecKeyTarget := p.newEvaluationCountkeyV2(eventCountKey, e.FeatureId, vID, targetEnvID, e.Timestamp)

				aggregator.addEvent(ecKeyTarget, ucKeyTarget, userID)
			}
		}
	}

	eventCounts, userCounts := aggregator.flush()
	failedECKeys, err := p.flushAggregatedCounts(eventCounts, userCounts)

	if err != nil {
		// Only mark events whose primary ecKey failed
		for ecKey := range failedECKeys {
			for eventID := range ecKeyToEventIDs[ecKey] {
				if _, alreadyFailed := fails[eventID]; !alreadyFailed {
					fails[eventID] = true
				}
			}
		}
	}

	failedMetricsCounts := make(map[metricsKey]int64)
	for ecKey := range failedECKeys {
		for mKey, count := range ecKeyToMetricsCounts[ecKey] {
			failedMetricsCounts[mKey] += count
		}
	}

	for key, count := range metricsAgg {
		succeeded := count - failedMetricsCounts[key]
		if succeeded > 0 {
			evaluationEventCounter.WithLabelValues(
				key.environmentId,
				key.sourceId,
				key.featureId,
				key.variationId,
			).Add(float64(succeeded))
		}
	}

	return fails
}

func (p *evaluationCountEventPersister) checkMessages(messages map[string]*puller.Message, fails map[string]bool) {
	for id, m := range messages {
		if repeatable, ok := fails[id]; ok {
			if repeatable {
				m.Nack()
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationCount, codes.RepeatableError.String()).Inc()
			} else {
				m.Ack()
				subscriberHandledCounter.WithLabelValues(subscriberEvaluationCount, codes.NonRepeatableError.String()).Inc()
			}
			continue
		}
		m.Ack()
		subscriberHandledCounter.WithLabelValues(subscriberEvaluationCount, codes.OK.String()).Inc()
	}
}

func (p *evaluationCountEventPersister) extractEvents(messages map[string]*puller.Message) environmentEventMap {
	envEvents := environmentEventMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		p.logger.Error("Bad proto message",
			zap.Error(err),
			zap.String("messageID", m.ID),
			zap.ByteString("data", m.Data),
			zap.Any("attributes", m.Attributes),
		)
		subscriberHandledCounter.WithLabelValues(subscriberEvaluationCount, codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		// Check if message data is empty
		if len(m.Data) == 0 {
			handleBadMessage(m, fmt.Errorf("message data is empty"))
			continue
		}
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		innerEvent := &eventproto.EvaluationEvent{}
		if err := ptypes.UnmarshalAny(event.Event, innerEvent); err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentId]; ok {
			innerEvents[event.Id] = innerEvent
			continue
		}
		envEvents[event.EnvironmentId] = eventMap{event.Id: innerEvent}
	}
	return envEvents
}

// isErrorReason returns true if the reason indicates the user received the default value
// due to an error (e.g., flag not found, cache miss). These should be counted toward
// the default variation. Must stay in sync with grpc_validation.go's isErrorReason.
func isErrorReason(reason *featureproto.Reason) bool {
	if reason == nil {
		return false
	}
	switch reason.Type {
	case featureproto.Reason_CLIENT, // deprecated, replaced by specific error types
		featureproto.Reason_ERROR_NO_EVALUATIONS,
		featureproto.Reason_ERROR_FLAG_NOT_FOUND,
		featureproto.Reason_ERROR_WRONG_TYPE,
		featureproto.Reason_ERROR_USER_ID_NOT_SPECIFIED,
		featureproto.Reason_ERROR_FEATURE_FLAG_ID_NOT_SPECIFIED,
		featureproto.Reason_ERROR_EXCEPTION,
		featureproto.Reason_ERROR_CACHE_NOT_FOUND:
		return true
	default:
		return false
	}
}

func getVariationID(reason *featureproto.Reason, vID string) (string, error) {
	if reason == nil {
		return "", ErrReasonNil
	}
	if isErrorReason(reason) {
		return defaultVariationID, nil
	}
	return vID, nil
}

// flushAggregatedCounts writes aggregated counts to Redis using individual calls
// with PFADD-before-INCRBY ordering per key pair.
//
// Why not pipelines? While go-redis pipelines do expose per-command errors via
// each Cmder's .Err() method after Exec(), individual calls were chosen as the
// simpler correctness-first approach. Pipeline optimization (batching all PFADDs
// then all INCRBYs with per-command error inspection) is planned as a follow-up.
//
// Within each key pair:
//   - PFADD runs first (idempotent). If it fails, INCRBY is skipped.
//   - INCRBY runs only after PFADD succeeds. If INCRBY fails, PFADD is
//     already done (idempotent) so retry is safe.
//
// Across key pairs: all pairs are attempted even if some fail. The caller
// receives the set of failed ecKeys so it can Nack only the affected events,
// preventing over-counting of already-succeeded key pairs on retry.
func (p *evaluationCountEventPersister) flushAggregatedCounts(
	eventCounts map[string]int64,
	userCounts map[string]map[string]struct{},
) (map[string]struct{}, error) {
	if len(eventCounts) == 0 && len(userCounts) == 0 {
		return nil, nil
	}

	// Match up ec and uc keys by extracting the common suffix
	// Admin keys: "ec:timestamp:featureID:variationID" / "uc:timestamp:featureID:variationID"
	// Non-admin keys: "envID:ec:timestamp:featureID:variationID" / "envID:uc:timestamp:featureID:variationID"
	type keyPair struct {
		ecKey    string
		ecCount  int64
		ucKey    string
		ucUsers  []string
		hasUsers bool
	}

	// Extract suffix to pair ec/uc keys
	// Keys format: Admin: "ec:timestamp:feat:var", Non-admin: "envID:ec:timestamp:feat:var"
	// We extract the common suffix after the kind to pair them correctly
	keyPairs := make(map[string]*keyPair)

	extractKey := func(key, kind string) string {
		// Extract a pairing key that includes environment ID to prevent collisions during migration.
		// Non-admin: "envID:kind:suffix" → "envID:suffix"
		// Admin: "kind:suffix" → "suffix"
		pattern := ":" + kind + ":"
		idx := strings.Index(key, pattern)
		if idx >= 0 {
			envPrefix := key[:idx]
			suffix := key[idx+len(pattern):]
			return envPrefix + ":" + suffix
		}
		// Admin format: "kind:suffix"
		return strings.TrimPrefix(key, kind+":")
	}

	// Build ec side of key pairs: extract pairing key and populate ecKey/ecCount.
	// If uc was already processed, add to existing pair; otherwise create new pair.
	for ecKey, count := range eventCounts {
		suffix := extractKey(ecKey, eventCountKey)
		if pair, exists := keyPairs[suffix]; exists {
			pair.ecKey = ecKey
			pair.ecCount = count
		} else {
			keyPairs[suffix] = &keyPair{
				ecKey:   ecKey,
				ecCount: count,
			}
		}
	}

	// Build uc side of key pairs: extract pairing key and populate ucKey/ucUsers.
	// If ec was already processed, add to existing pair; otherwise create new pair.
	for ucKey, userIDSet := range userCounts {
		suffix := extractKey(ucKey, userCountKey)
		userIDs := make([]string, 0, len(userIDSet))
		for userID := range userIDSet {
			userIDs = append(userIDs, userID)
		}

		if pair, exists := keyPairs[suffix]; exists {
			pair.ucKey = ucKey
			pair.ucUsers = userIDs
			pair.hasUsers = len(userIDs) > 0
		} else {
			keyPairs[suffix] = &keyPair{
				ucKey:    ucKey,
				ucUsers:  userIDs,
				hasUsers: len(userIDs) > 0,
			}
		}
	}

	failedECKeys := make(map[string]struct{})

	for _, pair := range keyPairs {
		// Step 1: PFADD first (idempotent - safe to retry)
		if pair.hasUsers && pair.ucKey != "" {
			_, err := p.evaluationCountCacher.PFAdd(pair.ucKey, pair.ucUsers...)
			if err != nil {
				if !strings.Contains(err.Error(), "client is closed") {
					p.logger.Error("Failed to add users to HyperLogLog",
						zap.Error(err),
						zap.String("ucKey", pair.ucKey),
					)
				}
				if pair.ecKey != "" {
					failedECKeys[pair.ecKey] = struct{}{}
				}
				continue
			}
		}

		// Step 2: INCRBY only if PFADD succeeded
		if pair.ecKey != "" {
			_, err := p.evaluationCountCacher.IncrementBy(pair.ecKey, pair.ecCount)
			if err != nil {
				if !strings.Contains(err.Error(), "client is closed") {
					p.logger.Error("Failed to increment event count",
						zap.Error(err),
						zap.String("ecKey", pair.ecKey),
						zap.Int64("count", pair.ecCount),
					)
				}
				failedECKeys[pair.ecKey] = struct{}{}
				continue
			}
		}
	}

	if len(failedECKeys) > 0 {
		p.logger.Error("Partial flush failure",
			zap.Int("failedKeyPairs", len(failedECKeys)),
			zap.Int("totalKeyPairs", len(keyPairs)),
		)
		return failedECKeys, fmt.Errorf("flush: %d/%d key pairs failed", len(failedECKeys), len(keyPairs))
	}

	p.logger.Debug("Flushed aggregated counts to Redis",
		zap.Int("keyPairs", len(keyPairs)),
	)

	return nil, nil
}

func (p *evaluationCountEventPersister) newEvaluationCountkeyV2(
	kind, featureID, variationID, environmentId string,
	timestamp int64,
) string {
	t := time.Unix(timestamp, 0)
	date := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%d:%s:%s", date.Unix(), featureID, variationID),
		environmentId,
	)
}

// bufferDAU buffers DAU entries in memory grouped by date and (envID, sourceID).
// User IDs are deduplicated in-memory using a set to reduce the PFADD payload
// sent to Redis, since a single RegisterEvent request can contain multiple
// evaluation events from the same user.
// The buffered entries are flushed to Redis periodically by writeDAUCache.
func (p *evaluationCountEventPersister) bufferDAU(envEvents environmentEventMap) {
	p.dauBufferMutex.Lock()
	defer p.dauBufferMutex.Unlock()
	for environmentId, events := range envEvents {
		for _, event := range events {
			userID := getUserID(event.UserId, event.User)
			if userID == "" {
				continue
			}
			key := dauBufferKey{
				dateStr:  time.Unix(event.Timestamp, 0).Format("20060102"),
				envID:    environmentId,
				sourceID: event.SourceId.String(),
			}
			if p.dauBuf[key] == nil {
				p.dauBuf[key] = make(map[string]struct{})
			}
			p.dauBuf[key][userID] = struct{}{}
		}
	}
}

// writeDAUCache periodically flushes the in-memory DAU buffer to Redis.
func (p *evaluationCountEventPersister) writeDAUCache(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(p.evaluationCountEventPersisterConfig.WriteDAUInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p.logger.Debug("Write DAU cache timer triggered")
			p.writeDAU()
		}
	}
}

func (p *evaluationCountEventPersister) writeDAU() {
	p.dauBufferMutex.Lock()
	buf := p.dauBuf
	p.dauBuf = make(dauBuffer)
	p.dauBufferMutex.Unlock()
	records := make([]cachev3.DAURecord, 0, len(buf))
	for key, userIDSet := range buf {
		date, err := time.Parse("20060102", key.dateStr)
		if err != nil {
			p.logger.Warn("Failed to parse DAU date",
				zap.Error(err),
				zap.String("date", key.dateStr),
			)
			continue
		}
		userIDs := make([]string, 0, len(userIDSet))
		for uid := range userIDSet {
			userIDs = append(userIDs, uid)
		}
		records = append(records, cachev3.DAURecord{
			Date:     date,
			EnvID:    key.envID,
			SourceID: key.sourceID,
			UserIDs:  userIDs,
		})
	}
	if err := p.dauCache.RecordDAUBatch(records); err != nil {
		if !strings.Contains(err.Error(), "client is closed") {
			p.logger.Warn("Failed to record DAU batch",
				zap.Error(err),
				zap.Int("recordCount", len(records)),
			)
		}
	}
}

func (p *evaluationCountEventPersister) cacheLastUsedInfoPerEnv(envEvents environmentEventMap) {
	for environmentId, events := range envEvents {
		for _, event := range events {
			p.cacheEnvLastUsedInfo(event, environmentId)
		}
		p.logger.Debug("Cache has been updated",
			zap.String("environmentId", environmentId),
			zap.Int("cacheSize", len(p.envLastUsedCache[environmentId])),
			zap.Int("eventSize", len(events)),
		)
	}
}

func (p *evaluationCountEventPersister) cacheEnvLastUsedInfo(
	event *eventproto.EvaluationEvent,
	environmentId string,
) {
	p.envLastUsedCacheMutex.Lock()
	defer p.envLastUsedCacheMutex.Unlock()
	var clientVersion string
	if event.User != nil {
		if version, ok := event.User.Data[userDataAppVersion]; ok {
			clientVersion = version
		}
	}
	id := ftdomain.FeatureLastUsedInfoID(event.FeatureId, event.FeatureVersion)
	if cache, ok := p.envLastUsedCache[environmentId]; ok {
		if info, ok := cache[id]; ok {
			info.UsedAt(event.Timestamp)
			if err := info.SetClientVersion(clientVersion); err != nil {
				p.logger.Error("Failed to set client version",
					zap.Error(err),
					zap.String("environmentId", environmentId),
					zap.String("featureId", info.FeatureId),
					zap.Int32("featureVersion", info.Version),
					zap.String("clientVersion", clientVersion))
			}
			return
		}
		cache[id] = ftdomain.NewFeatureLastUsedInfo(
			event.FeatureId,
			event.FeatureVersion,
			event.Timestamp,
			clientVersion,
		)
		return
	}
	cache := lastUsedInfoCache{}
	cache[id] = ftdomain.NewFeatureLastUsedInfo(
		event.FeatureId,
		event.FeatureVersion,
		event.Timestamp,
		clientVersion,
	)
	p.envLastUsedCache[environmentId] = cache
}

// Write the feature flag last-used cache in the MySQL and reset the cache
func (p *evaluationCountEventPersister) writeFlagLastUsedInfoCache(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(p.evaluationCountEventPersisterConfig.WriteCacheInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p.logger.Debug("Write FlagLastUsedInfo cache timer triggered")
			p.writeEnvLastUsedInfo()
		}
	}
}

func (p *evaluationCountEventPersister) writeEnvLastUsedInfo() {
	p.envLastUsedCacheMutex.Lock()
	defer p.envLastUsedCacheMutex.Unlock()

	for environmentId, cache := range p.envLastUsedCache {
		info := make([]*ftdomain.FeatureLastUsedInfo, 0, len(cache))
		for _, v := range cache {
			info = append(info, v)
		}
		if err := p.upsertMultiFeatureLastUsedInfo(context.Background(), info, environmentId); err != nil {
			p.logger.Error("Failed to write feature last-used info", zap.Error(err),
				zap.String("environmentId", environmentId))
			continue
		}
		p.logger.Debug("Cache has been written",
			zap.String("environmentId", environmentId),
			zap.Int("cacheSize", len(info)),
		)
	}
	// Reset the cache
	p.envLastUsedCache = make(environmentLastUsedInfoCache)
}

func (p *evaluationCountEventPersister) upsertMultiFeatureLastUsedInfo(
	ctx context.Context,
	featureLastUsedInfos []*ftdomain.FeatureLastUsedInfo,
	environmentId string,
) error {
	ids := make([]string, 0, len(featureLastUsedInfos))
	for _, f := range featureLastUsedInfos {
		ids = append(ids, f.ID())
	}
	storage := ftstorage.NewFeatureLastUsedInfoStorage(p.mysqlClient)
	updatedInfo := make([]*ftdomain.FeatureLastUsedInfo, 0, len(ids))
	currentInfo, err := storage.GetFeatureLastUsedInfos(ctx, ids, environmentId)
	if err != nil {
		return err
	}
	currentInfoMap := make(map[string]*ftdomain.FeatureLastUsedInfo, len(currentInfo))
	for _, c := range currentInfo {
		currentInfoMap[c.ID()] = c
	}
	for _, f := range featureLastUsedInfos {
		v, ok := currentInfoMap[f.ID()]
		if !ok {
			updatedInfo = append(updatedInfo, f)
			continue
		}
		var update bool
		if v.LastUsedAt < f.LastUsedAt {
			update = true
			v.LastUsedAt = f.LastUsedAt
		}
		if v.ClientOldestVersion != f.ClientOldestVersion {
			update = true
			v.ClientOldestVersion = f.ClientOldestVersion
		}
		if v.ClientLatestVersion != f.ClientLatestVersion {
			update = true
			v.ClientLatestVersion = f.ClientLatestVersion
		}
		if update {
			updatedInfo = append(updatedInfo, v)
		}
	}
	for _, info := range updatedInfo {
		if err := p.upsertFeatureLastUsedInfo(ctx, info, environmentId); err != nil {
			return err
		}
	}
	return nil
}

func (p *evaluationCountEventPersister) upsertFeatureLastUsedInfo(
	ctx context.Context,
	featureLastUsedInfo *ftdomain.FeatureLastUsedInfo,
	environmentId string,
) error {
	storage := ftstorage.NewFeatureLastUsedInfoStorage(p.mysqlClient)
	if err := storage.UpsertFeatureLastUsedInfo(
		ctx,
		featureLastUsedInfo,
		environmentId,
	); err != nil {
		return err
	}
	return nil
}

func (p *evaluationCountEventPersister) cacheUserAttributes(envEvents environmentEventMap) {
	p.userAttributesCacheMutex.Lock()
	defer p.userAttributesCacheMutex.Unlock()
	for environmentId, events := range envEvents {
		userAttributesMap := make(map[string]*userproto.UserAttribute)

		if existingCache, exists := p.userAttributesCache[environmentId]; exists {
			for _, attr := range existingCache.UserAttributes {
				userAttributesMap[attr.Key] = &userproto.UserAttribute{
					Key:    attr.Key,
					Values: make([]string, len(attr.Values)),
				}
				copy(userAttributesMap[attr.Key].Values, attr.Values)
			}
		}

		for _, event := range events {
			if event.User == nil || event.User.Data == nil {
				continue
			}

			// Extract user attributes from User.Data
			for key, value := range event.User.Data {
				if key == "" {
					continue
				}

				if attr, exists := userAttributesMap[key]; exists {
					// Check if value already exists to avoid duplicates
					found := false
					for _, existingValue := range attr.Values {
						if existingValue == value {
							found = true
							break
						}
					}
					if !found {
						attr.Values = append(attr.Values, value)
					}
				} else {
					userAttributesMap[key] = &userproto.UserAttribute{
						Key:    key,
						Values: []string{value},
					}
				}
			}
		}

		// Convert map to slice and save to cache
		if len(userAttributesMap) > 0 {
			userAttributes := &userproto.UserAttributes{
				EnvironmentId:  environmentId,
				UserAttributes: make([]*userproto.UserAttribute, 0, len(userAttributesMap)),
			}

			for _, attr := range userAttributesMap {
				userAttributes.UserAttributes = append(userAttributes.UserAttributes, attr)
			}
			p.userAttributesCache[environmentId] = userAttributes
		}
	}
}

func (p *evaluationCountEventPersister) writeUserAttributesCache(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(p.evaluationCountEventPersisterConfig.WriteCacheInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p.writeUserAttributes()
		}
	}
}

func (p *evaluationCountEventPersister) writeUserAttributes() {
	p.userAttributesCacheMutex.Lock()
	defer p.userAttributesCacheMutex.Unlock()

	for envID, cache := range p.userAttributesCache {
		if cache != nil && len(cache.UserAttributes) > 0 {
			if err := p.upsertUserAttributes(cache); err != nil {
				p.logger.Error(
					"Failed to save user attributes, will retry next cycle",
					zap.Error(err),
					zap.String("environmentId", envID),
				)
				continue
			}
			// If successful, delete it from the cache.
			// The failed items will remain for the next attempt.
			delete(p.userAttributesCache, envID)
		}
	}
}

func (p *evaluationCountEventPersister) upsertUserAttributes(
	userAttributes *userproto.UserAttributes,
) error {
	ttl := time.Duration(p.evaluationCountEventPersisterConfig.UserAttributeKeyTTL) * time.Second
	if err := p.userAttributesCacher.Put(userAttributes, ttl); err != nil {
		p.logger.Error("Failed to save user attributes to cache",
			zap.Error(err),
			zap.String("environmentId", userAttributes.EnvironmentId),
			zap.Any("attributes", userAttributes.UserAttributes),
			zap.Int("attributeCount", len(userAttributes.UserAttributes)),
		)
		return err
	}

	// MIGRATION: Double-write to the target environment ID if configured
	if targetEnvID := getMigrationTargetEnvironmentID(userAttributes.EnvironmentId); targetEnvID != "" {
		migrationAttrs := &userproto.UserAttributes{
			EnvironmentId:  targetEnvID,
			UserAttributes: userAttributes.UserAttributes,
		}
		if err := p.userAttributesCacher.Put(migrationAttrs, ttl); err != nil {
			p.logger.Error("Migration: Failed to save user attributes for target environment",
				zap.Error(err),
				zap.String("fromEnvironmentId", userAttributes.EnvironmentId),
				zap.String("toEnvironmentId", targetEnvID),
			)
		}
	}
	return nil
}
