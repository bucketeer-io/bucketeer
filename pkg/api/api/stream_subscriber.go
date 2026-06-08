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
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// StreamSubscriber forwards relevant domain events to the StreamDispatcher.
// It assumes L2-cache is already fresh on arrival.
// It evicts L1-cache before dispatching to certainly evaluate with latest data.
type StreamSubscriber struct {
	dispatcher        *StreamDispatcher
	featuresCache     cachev3.FeaturesCache
	segmentUsersCache cachev3.SegmentUsersCache
	logger            *zap.Logger
}

func NewStreamSubscriber(
	dispatcher *StreamDispatcher,
	featuresCache cachev3.FeaturesCache,
	segmentUsersCache cachev3.SegmentUsersCache,
	logger *zap.Logger,
) *StreamSubscriber {
	return &StreamSubscriber{
		dispatcher:        dispatcher,
		featuresCache:     featuresCache,
		segmentUsersCache: segmentUsersCache,
		logger:            logger.Named("stream-subscriber"),
	}
}

func (s *StreamSubscriber) Run(ctx context.Context, msgChan <-chan *puller.Message) error {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			s.handleMessage(msg)
			msg.Ack()
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *StreamSubscriber) handleMessage(msg *puller.Message) {
	event := &domaineventproto.Event{}
	if err := proto.Unmarshal(msg.Data, event); err != nil {
		s.logger.Warn("Failed to unmarshal domain event", zap.Error(err))
		return
	}
	switch event.EntityType {
	case domaineventproto.Event_FEATURE:
		if event.Type != domaineventproto.Event_FEATURE_UPDATED &&
			event.Type != domaineventproto.Event_FEATURE_ENABLED &&
			event.Type != domaineventproto.Event_FEATURE_DISABLED {
			return
		}
		if err := s.featuresCache.Evict(event.EnvironmentId); err != nil {
			s.logger.Warn("Failed to evict features cache before SSE dispatch",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("entityId", event.EntityId),
				zap.String("type", event.Type.String()),
			)
		}
		s.dispatcher.dispatch(streamEvent{
			environmentID: event.EnvironmentId,
			tags:          s.affectedTags(event),
		})
	case domaineventproto.Event_SEGMENT:
		// Segment membership only changes through a bulk upload.
		if event.Type != domaineventproto.Event_SEGMENT_BULK_UPLOAD_USERS_STATUS_CHANGED {
			return
		}
		if event.Data == nil {
			return
		}
		payload := &domaineventproto.SegmentBulkUploadUsersStatusChangedEvent{}
		if err := event.Data.UnmarshalTo(payload); err != nil {
			s.logger.Warn("Failed to unmarshal segment bulk upload event", zap.Error(err))
			return
		}
		// Membership changes only on a succeeded upload.
		if payload.Status != featureproto.Segment_SUCEEDED {
			return
		}
		if err := s.segmentUsersCache.Evict(event.EntityId, event.EnvironmentId); err != nil {
			s.logger.Warn("Failed to evict segment users cache before SSE dispatch",
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
				zap.String("segmentId", event.EntityId),
				zap.String("type", event.Type.String()),
			)
		}
		// TODO: resolve the affected tags from the segment.
		// Currently, it fans out env-wide (all tags).
		s.dispatcher.dispatch(streamEvent{
			environmentID: event.EnvironmentId,
		})
	}
}

// affectedTags unions the flag's tags before and after the update so that
// removing a tag still notifies that tag's subscribers.
//
// TODO: also pull the tags of flags that depend on this one.
func (s *StreamSubscriber) affectedTags(event *domaineventproto.Event) []string {
	seen := make(map[string]struct{})
	var tags []string
	for _, data := range []string{event.EntityData, event.PreviousEntityData} {
		for _, tag := range s.parseTags(data) {
			if _, ok := seen[tag]; ok {
				continue
			}
			seen[tag] = struct{}{}
			tags = append(tags, tag)
		}
	}
	return tags
}

func (s *StreamSubscriber) parseTags(data string) []string {
	if data == "" {
		return nil
	}
	var payload struct {
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		s.logger.Warn("Failed to extract tags from feature entity data", zap.Error(err))
		return nil
	}
	return payload.Tags
}
