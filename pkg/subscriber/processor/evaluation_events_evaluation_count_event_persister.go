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

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/pkg/cache/v3"
	ftdomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
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

type EvaluationCountEventPersisterConfig struct {
	FlushSize           int `json:"flushSize"`
	FlushInterval       int `json:"flushInterval"`
	WriteCacheInterval  int `json:"writeCacheInterval"`
	UserAttributeKeyTTL int `json:"userAttributeKeyTtl"`
}

type evaluationCountEventPersister struct {
	evaluationCountEventPersisterConfig EvaluationCountEventPersisterConfig
	mysqlClient                         mysql.Client
	envLastUsedCache                    environmentLastUsedInfoCache
	evaluationCountCacher               cache.MultiGetDeleteCountCache
	userAttributesCacher                cachev3.UserAttributesCache
	userAttributesCache                 userAttributesCache
	envLastUsedCacheMutex               sync.Mutex
	userAttributesCacheMutex            sync.Mutex
	logger                              *zap.Logger
}

func NewEvaluationCountEventPersister(
	ctx context.Context,
	config interface{},
	mysqlClient mysql.Client,
	evaluationCountCacher cache.MultiGetDeleteCountCache,
	userAttributesCacher cachev3.UserAttributesCache,
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
		userAttributesCache:                 make(userAttributesCache),
		envLastUsedCacheMutex:               sync.Mutex{},
		userAttributesCacheMutex:            sync.Mutex{},
		logger:                              logger,
	}
	// write flag last used info cache periodically
	//nolint:errcheck
	go e.writeFlagLastUsedInfoCache(ctx)
	// write user attributes cache periodically
	//nolint:errcheck
	go e.writeUserAttributesCache(ctx)
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
			updateEvaluationCounter(envEvents)
		case <-ticker.C:
			envEvents := p.extractEvents(batch)
			// Update the feature flag last-used cache
			p.cacheLastUsedInfoPerEnv(envEvents)
			// Update the user attributes cache
			p.cacheUserAttributes(envEvents)
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
	for environmentId, events := range envEvents {
		for id, event := range events {
			// Increment the evaluation event count in the Redis
			if err := p.incrementEvaluationCount(id, event, environmentId); err != nil {
				if errors.Is(err, ErrReasonNil) {
					fails[id] = false
				} else {
					fails[id] = true
				}
			}
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

func getVariationID(reason *featureproto.Reason, vID string) (string, error) {
	if reason == nil {
		return "", ErrReasonNil
	}
	if reason.Type == featureproto.Reason_CLIENT {
		return defaultVariationID, nil
	}
	return vID, nil
}

func (p *evaluationCountEventPersister) incrementEvaluationCount(
	eventID string,
	event proto.Message,
	environmentId string,
) error {
	if e, ok := event.(*eventproto.EvaluationEvent); ok {
		vID, err := getVariationID(e.Reason, e.VariationId)
		if err != nil {
			return err
		}
		// To avoid duplication when the request fails, we increment the event count in the end
		// because the user count is an unique count, and there is no problem adding the same event more than once
		// We tried to use Pipeline and indeed it improves the response time,
		// but it also increases the Pod CPU usage. It's a trade-off.
		// Since this is a background service and it's not latency-sensitive, we split the requests.
		ucKey := p.newEvaluationCountkeyV2(userCountKey, e.FeatureId, vID, environmentId, e.Timestamp)
		userID := getUserID(e.UserId, e.User)
		if err := p.countUser(ucKey, userID); err != nil {
			if !strings.Contains(err.Error(), "client is closed") {
				p.logger.Error("Failed to increment the evaluation user event in the Redis",
					zap.Error(err),
					zap.String("environmentId", environmentId),
					zap.String("eventId", eventID),
					zap.String("userId", userID),
					zap.String("userCountKey", ucKey),
					zap.Any("evaluationEvent", e),
				)
			}
			return err
		}
		ecKey := p.newEvaluationCountkeyV2(eventCountKey, e.FeatureId, vID, environmentId, e.Timestamp)
		if err := p.countEvent(ecKey); err != nil {
			if !strings.Contains(err.Error(), "client is closed") {
				p.logger.Error("Failed to increment the evaluation event in the Redis",
					zap.Error(err),
					zap.String("environmentId", environmentId),
					zap.String("eventId", eventID),
					zap.String("userId", userID),
					zap.String("eventCountKey", ecKey),
					zap.Any("evaluationEvent", e),
				)
			}
			return err
		}
		evaluationEventCounter.WithLabelValues(
			environmentId,
			e.SdkVersion,
			e.FeatureId,
			e.Metadata[appVersion],
			e.VariationId,
		).Inc()
	}
	return nil
}

func (p *evaluationCountEventPersister) countEvent(key string) error {
	_, err := p.evaluationCountCacher.Increment(key)
	if err != nil {
		return err
	}
	return nil
}

func (p *evaluationCountEventPersister) countUser(key, userID string) error {
	_, err := p.evaluationCountCacher.PFAdd(key, userID)
	if err != nil {
		return err
	}
	return nil
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
	if event.User == nil {
		p.logger.Warn("Failed to cache last used info. User is nil.",
			zap.String("environmentId", environmentId))
	} else {
		clientVersion = event.User.Data[userDataAppVersion]
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
	if err := p.userAttributesCacher.Put(
		userAttributes,
		time.Duration(p.evaluationCountEventPersisterConfig.UserAttributeKeyTTL)*time.Second,
	); err != nil {
		p.logger.Error("Failed to save user attributes to cache",
			zap.Error(err),
			zap.String("environmentId", userAttributes.EnvironmentId),
			zap.Any("attributes", userAttributes.UserAttributes),
			zap.Int("attributeCount", len(userAttributes.UserAttributes)),
		)
		return err
	} else {
		p.logger.Debug("Successfully saved user attributes to cache",
			zap.String("environmentId", userAttributes.EnvironmentId),
			zap.Int("attributeCount", len(userAttributes.UserAttributes)),
		)
		return nil
	}
}
