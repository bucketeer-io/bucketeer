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
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber"
	"github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/processor"
)

type stubCloser struct {
	err error
}

func (s stubCloser) Close() error {
	return s.err
}

func TestCloseDWHWriter(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc       string
		closer     stubCloser
		wantErrors int
		wantMsg    string
	}{
		{
			desc:       "io.EOF on close is not logged as error",
			closer:     stubCloser{err: io.EOF},
			wantErrors: 0,
		},
		{
			desc:       "nil on close is not logged as error",
			closer:     stubCloser{err: nil},
			wantErrors: 0,
		},
		{
			desc:       "real close errors are logged",
			closer:     stubCloser{err: errors.New("connection reset")},
			wantErrors: 1,
			wantMsg:    "subscriber: failed to close evaluation event writer",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
			core := zapcore.NewCore(enc, zapcore.AddSync(&buf), zapcore.ErrorLevel)
			logger := zap.New(core)

			closeDWHWriter(p.closer, "evaluation event writer", logger)

			logOutput := strings.TrimSpace(buf.String())
			if p.wantErrors == 0 {
				assert.Empty(t, logOutput)
				return
			}
			assert.Contains(t, logOutput, p.wantMsg)
		})
	}
}

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

func TestValidatePushSenderConfig(t *testing.T) {
	t.Parallel()

	refresher := subscriber.Configuration{
		Topic:                  "domain-events",
		Subscription:           "domain-events-cache-refresher",
		CacheInvalidationTopic: "cache-invalidation",
	}
	patterns := []struct {
		desc    string
		configs map[string]subscriber.Configuration
		wantErr string
	}{
		{
			desc:    "nil config: nothing to validate",
			configs: nil,
		},
		{
			desc: "no pushSender block: nothing to validate",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: refresher,
			},
		},
		{
			desc: "pushSender on the invalidation topic with its own subscription",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: refresher,
				processor.PushSenderName: {
					Topic:        "cache-invalidation",
					Subscription: "cache-invalidation-push-sender",
				},
			},
		},
		{
			desc: "pushSender still on the domain-event topic",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: refresher,
				processor.PushSenderName: {
					Topic:        "domain-events",
					Subscription: "domain-events-push-sender",
				},
			},
			wantErr: `pushSender.topic must be the cacheRefresher cacheInvalidationTopic "cache-invalidation", ` +
				`got "domain-events"`,
		},
		{
			desc: "pushSender without a cacheRefresher block",
			configs: map[string]subscriber.Configuration{
				processor.PushSenderName: {
					Topic:        "cache-invalidation",
					Subscription: "cache-invalidation-push-sender",
				},
			},
			wantErr: "pushSender requires the cacheRefresher processor with a non-empty cacheInvalidationTopic",
		},
		{
			desc: "pushSender with cacheRefresher announcements disabled",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: {Topic: "domain-events", Subscription: "domain-events-cache-refresher"},
				processor.PushSenderName: {
					Topic:        "cache-invalidation",
					Subscription: "cache-invalidation-push-sender",
				},
			},
			wantErr: "pushSender requires the cacheRefresher processor with a non-empty cacheInvalidationTopic",
		},
		{
			desc: "pushSender sharing the cacheRefresher subscription would steal its events",
			configs: map[string]subscriber.Configuration{
				processor.CacheRefresherName: refresher,
				processor.PushSenderName: {
					Topic:        "cache-invalidation",
					Subscription: "domain-events-cache-refresher",
				},
			},
			wantErr: `pushSender.subscription "domain-events-cache-refresher" must not be shared with the cacheRefresher`,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			err := validatePushSenderConfig(p.configs)
			if p.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, p.wantErr)
		})
	}
}
