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

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/processor"
)

func TestResolveCacheInvalidationConfig(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc           string
		configs        map[string]subscriber.Configuration
		wantEnabled    bool
		wantTopic      string
		wantPubSubType string
	}{
		{
			desc:        "nil config map disables announcements",
			configs:     nil,
			wantEnabled: false,
		},
		{
			desc:        "empty config map disables announcements",
			configs:     map[string]subscriber.Configuration{},
			wantEnabled: false,
		},
		{
			desc: "no cacheRefresher entry disables announcements",
			configs: map[string]subscriber.Configuration{
				"someOtherProcessor": {
					PubSubType:             "google",
					CacheInvalidationTopic: "bucketeer-cache-invalidation",
				},
			},
			wantEnabled: false,
		},
		{
			desc: "cacheRefresher entry with empty CacheInvalidationTopic disables announcements (legacy evict-only)",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: {
					PubSubType:             "google",
					Topic:                  "bucketeer-domain-events",
					CacheInvalidationTopic: "",
				},
			},
			wantEnabled: false,
		},
		{
			desc: "cacheRefresher entry with CacheInvalidationTopic enables announcements",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: {
					PubSubType:             "google",
					Project:                "bucketeer-prj",
					Topic:                  "bucketeer-domain-events",
					CacheInvalidationTopic: "bucketeer-cache-invalidation",
				},
			},
			wantEnabled:    true,
			wantTopic:      "bucketeer-cache-invalidation",
			wantPubSubType: "google",
		},
		{
			desc: "redis-stream backend is also returned verbatim",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: {
					PubSubType:             "redis-stream",
					RedisAddr:              "redis:6379",
					RedisPartitionCount:    16,
					Topic:                  "bucketeer-domain-events",
					CacheInvalidationTopic: "bucketeer-cache-invalidation",
				},
			},
			wantEnabled:    true,
			wantTopic:      "bucketeer-cache-invalidation",
			wantPubSubType: "redis-stream",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			conf, enabled := resolveCacheInvalidationConfig(p.configs)
			assert.Equal(t, p.wantEnabled, enabled)
			if !p.wantEnabled {
				// On the disabled path the returned conf must be the
				// zero value so callers can't accidentally use stale
				// fields for backend selection.
				assert.Equal(t, subscriber.Configuration{}, conf)
				return
			}
			assert.Equal(t, p.wantTopic, conf.CacheInvalidationTopic)
			assert.Equal(t, p.wantPubSubType, conf.PubSubType)
		})
	}
}
