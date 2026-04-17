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

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/puller"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCacheInvalidatorHandleMessage(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc             string
		event            *domaineventproto.Event
		setupCache       func(fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, akc cachev3.EnvironmentAPIKeyCache)
		verifyNotEvicted func(t *testing.T, fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, akc cachev3.EnvironmentAPIKeyCache)
		verifyEvicted    func(t *testing.T, fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, akc cachev3.EnvironmentAPIKeyCache)
	}{
		{
			desc: "feature event evicts features cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				EntityId:      "feature-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_FEATURE_UPDATED,
			},
			setupCache: func(fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				fc.Put(&featureproto.Features{
					Features: []*featureproto.Feature{{Id: "feature-id-1"}},
				}, "env-1")
			},
			verifyEvicted: func(t *testing.T, fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				_, err := fc.Get("env-1")
				assert.Error(t, err, "features cache should be evicted for env-1")
			},
		},
		{
			desc: "segment event evicts segment users cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_SEGMENT,
				EntityId:      "segment-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_SEGMENT_CREATED,
			},
			setupCache: func(fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				sc.Put(&featureproto.SegmentUsers{
					SegmentId: "segment-id-1",
					Users:     []*featureproto.SegmentUser{{UserId: "u1"}},
				}, "env-1")
			},
			verifyEvicted: func(t *testing.T, fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				_, err := sc.Get("segment-id-1", "env-1")
				assert.Error(t, err, "segment users cache should be evicted")
			},
		},
		{
			desc: "api key event evicts environment api key cache",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_APIKEY,
				EntityId:      "apikey-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_APIKEY_CHANGED,
				EntityData:    `{"api_key": "secret-123"}`,
			},
			setupCache: func(_ cachev3.FeaturesCache, _ cachev3.SegmentUsersCache, akc cachev3.EnvironmentAPIKeyCache) {
				akc.Put(&accountproto.EnvironmentAPIKey{
					ApiKey: &accountproto.APIKey{ApiKey: "secret-123"},
				})
			},
			verifyEvicted: func(t *testing.T, _ cachev3.FeaturesCache, _ cachev3.SegmentUsersCache, akc cachev3.EnvironmentAPIKeyCache) {
				_, err := akc.Get("secret-123")
				assert.Error(t, err, "api key cache should be evicted for secret-123")
			},
		},
		{
			desc: "api key event evicts both old and new secrets on rotation",
			event: &domaineventproto.Event{
				EntityType:         domaineventproto.Event_APIKEY,
				EntityId:           "apikey-id-1",
				EnvironmentId:      "env-1",
				Type:               domaineventproto.Event_APIKEY_CHANGED,
				EntityData:         `{"api_key": "new-secret"}`,
				PreviousEntityData: `{"api_key": "old-secret"}`,
			},
			setupCache: func(_ cachev3.FeaturesCache, _ cachev3.SegmentUsersCache, akc cachev3.EnvironmentAPIKeyCache) {
				akc.Put(&accountproto.EnvironmentAPIKey{
					ApiKey: &accountproto.APIKey{ApiKey: "old-secret"},
				})
				akc.Put(&accountproto.EnvironmentAPIKey{
					ApiKey: &accountproto.APIKey{ApiKey: "new-secret"},
				})
			},
			verifyEvicted: func(t *testing.T, _ cachev3.FeaturesCache, _ cachev3.SegmentUsersCache, akc cachev3.EnvironmentAPIKeyCache) {
				_, err := akc.Get("old-secret")
				assert.Error(t, err, "old api key cache should be evicted")
				_, err = akc.Get("new-secret")
				assert.Error(t, err, "new api key cache should be evicted")
			},
		},
		{
			desc: "unrelated entity type does not evict",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_ACCOUNT,
				EntityId:      "account-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_ACCOUNT_V2_CREATED,
			},
			setupCache: func(fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				fc.Put(&featureproto.Features{
					Features: []*featureproto.Feature{{Id: "feature-id-1"}},
				}, "env-1")
			},
			verifyNotEvicted: func(t *testing.T, fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				features, err := fc.Get("env-1")
				assert.NoError(t, err, "features cache should NOT be evicted for unrelated entity type")
				assert.Len(t, features.Features, 1)
			},
		},
		{
			desc: "feature event only evicts the target environment",
			event: &domaineventproto.Event{
				EntityType:    domaineventproto.Event_FEATURE,
				EntityId:      "feature-id-1",
				EnvironmentId: "env-1",
				Type:          domaineventproto.Event_FEATURE_VERSION_INCREMENTED,
			},
			setupCache: func(fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				fc.Put(&featureproto.Features{
					Features: []*featureproto.Feature{{Id: "f1"}},
				}, "env-1")
				fc.Put(&featureproto.Features{
					Features: []*featureproto.Feature{{Id: "f2"}},
				}, "env-2")
			},
			verifyEvicted: func(t *testing.T, fc cachev3.FeaturesCache, sc cachev3.SegmentUsersCache, _ cachev3.EnvironmentAPIKeyCache) {
				_, err := fc.Get("env-1")
				assert.Error(t, err, "env-1 should be evicted")
				features, err := fc.Get("env-2")
				assert.NoError(t, err, "env-2 should NOT be evicted")
				assert.Equal(t, "f2", features.Features[0].Id)
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			inMemoryCache := cachev3.NewInMemoryCache()
			featuresCache := cachev3.NewFeaturesCache(inMemoryCache, 10*time.Minute)
			segmentUsersCache := cachev3.NewSegmentUsersCache(inMemoryCache, 10*time.Minute)
			environmentAPIKeyCache := cachev3.NewEnvironmentAPIKeyCache(inMemoryCache, 10*time.Minute)

			invalidator := NewCacheInvalidator(featuresCache, segmentUsersCache, environmentAPIKeyCache, zap.NewNop())

			if p.setupCache != nil {
				p.setupCache(featuresCache, segmentUsersCache, environmentAPIKeyCache)
			}

			data, err := proto.Marshal(p.event)
			require.NoError(t, err)

			msg := &puller.Message{Data: data}
			invalidator.handleMessage(msg)

			if p.verifyEvicted != nil {
				p.verifyEvicted(t, featuresCache, segmentUsersCache, environmentAPIKeyCache)
			}
			if p.verifyNotEvicted != nil {
				p.verifyNotEvicted(t, featuresCache, segmentUsersCache, environmentAPIKeyCache)
			}
		})
	}
}
