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
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
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

func TestDWHPersisterConfig(t *testing.T) {
	t.Parallel()

	onDemand := map[string]interface{}{"name": "onDemand"}
	evalDedicated := map[string]interface{}{"name": "evalDedicated"}
	goalDedicated := map[string]interface{}{"name": "goalDedicated"}
	patterns := []struct {
		desc       string
		persisters map[string]interface{}
		expected   interface{}
	}{
		{
			desc:       "no persisters",
			persisters: map[string]interface{}{},
			expected:   nil,
		},
		{
			desc: "prefers the on-demand evaluation persister",
			persisters: map[string]interface{}{
				processor.EvaluationCountEventDWHPersisterName: onDemand,
				processor.EvaluationEventDWHPersisterName:      evalDedicated,
				processor.GoalEventDWHPersisterName:            goalDedicated,
			},
			expected: onDemand,
		},
		{
			desc: "falls back to the dedicated evaluation persister",
			persisters: map[string]interface{}{
				processor.EvaluationEventDWHPersisterName: evalDedicated,
				processor.GoalEventDWHPersisterName:       goalDedicated,
			},
			expected: evalDedicated,
		},
		{
			desc: "falls back to the dedicated goal persister",
			persisters: map[string]interface{}{
				processor.GoalEventDWHPersisterName: goalDedicated,
			},
			expected: goalDedicated,
		},
		{
			desc: "falls back to the on-demand goal persister",
			persisters: map[string]interface{}{
				processor.GoalCountEventDWHPersisterName: onDemand,
			},
			expected: onDemand,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, p.expected, dwhPersisterConfig(p.persisters))
		})
	}
}

func TestEventPersisterConfigs(t *testing.T) {
	t.Parallel()

	legacyDWH := map[string]interface{}{"name": "legacyDWH"}
	legacyOPS := map[string]interface{}{"name": "legacyOPS"}
	dedicatedDWH := map[string]interface{}{"name": "dedicatedDWH"}
	dedicatedOPS := map[string]interface{}{"name": "dedicatedOPS"}
	patterns := []struct {
		desc        string
		processors  map[string]interface{}
		onDemand    map[string]interface{}
		expectedDWH map[string]interface{}
		expectedOPS map[string]interface{}
	}{
		{
			desc:        "nothing configured",
			expectedDWH: map[string]interface{}{},
			expectedOPS: map[string]interface{}{},
		},
		{
			desc: "dedicated only with an empty on-demand config",
			processors: map[string]interface{}{
				processor.EvaluationEventDWHPersisterName: dedicatedDWH,
				processor.EvaluationEventOPSPersisterName: dedicatedOPS,
				processor.AuditLogPersisterName:           map[string]interface{}{},
			},
			onDemand: map[string]interface{}{},
			expectedDWH: map[string]interface{}{
				processor.EvaluationEventDWHPersisterName: dedicatedDWH,
			},
			expectedOPS: map[string]interface{}{
				processor.EvaluationEventOPSPersisterName: dedicatedOPS,
			},
		},
		{
			desc: "legacy and dedicated together",
			processors: map[string]interface{}{
				processor.GoalEventDWHPersisterName: dedicatedDWH,
				processor.GoalEventOPSPersisterName: dedicatedOPS,
			},
			onDemand: map[string]interface{}{
				processor.GoalCountEventDWHPersisterName: legacyDWH,
				processor.GoalCountEventOPSPersisterName: legacyOPS,
			},
			expectedDWH: map[string]interface{}{
				processor.GoalCountEventDWHPersisterName: legacyDWH,
				processor.GoalEventDWHPersisterName:      dedicatedDWH,
			},
			expectedOPS: map[string]interface{}{
				processor.GoalCountEventOPSPersisterName: legacyOPS,
				processor.GoalEventOPSPersisterName:      dedicatedOPS,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			dwh, ops := eventPersisterConfigs(p.processors, p.onDemand)
			assert.Equal(t, p.expectedDWH, dwh)
			assert.Equal(t, p.expectedOPS, ops)
		})
	}
}

func TestInitDataWarehouseStoragesMySQL(t *testing.T) {
	t.Parallel()

	empty := ""
	zero := 0
	emptyFlags := &server{
		mysqlUser:   &empty,
		mysqlPass:   &empty,
		mysqlHost:   &empty,
		mysqlPort:   &zero,
		mysqlDBName: &empty,
	}
	patterns := []struct {
		desc        string
		mysqlConfig MySQLConfig
		mainClient  func(ctrl *gomock.Controller) mysql.Client
		expectedErr string
	}{
		{
			desc:        "main connection without main client or mysql flags fails",
			mysqlConfig: MySQLConfig{UseMainConnection: true},
			mainClient:  func(*gomock.Controller) mysql.Client { return nil },
			expectedErr: "needs the --mysql-* flags when storage-type is not mysql",
		},
		{
			desc:        "dedicated connection without host fails before dialling",
			mysqlConfig: MySQLConfig{User: "user", Database: "db"},
			mainClient:  func(*gomock.Controller) mysql.Client { return nil },
			expectedErr: "failed to create dedicated MySQL client",
		},
		{
			desc:        "main connection reuses the main client and has no cleanup",
			mysqlConfig: MySQLConfig{UseMainConnection: true},
			mainClient: func(ctrl *gomock.Controller) mysql.Client {
				return mysqlmock.NewMockClient(ctrl)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			dwhConfig := DataWarehouseConfig{Type: "mysql", MySQL: p.mysqlConfig}
			evalWriter, goalWriter, goalStorage, cleanup, err := emptyFlags.initDataWarehouseStorages(
				context.Background(), dwhConfig, p.mainClient(ctrl), nil, zap.NewNop(),
			)
			assert.Nil(t, cleanup)
			if p.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), p.expectedErr)
				assert.Nil(t, evalWriter)
				assert.Nil(t, goalWriter)
				assert.Nil(t, goalStorage)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, evalWriter)
			assert.NotNil(t, goalWriter)
			assert.NotNil(t, goalStorage)
		})
	}
}
