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

package rediscounter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	evmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	redis "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	evproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

// newMiniredisDeleter runs the deleter against an in-process Redis so the
// real SCAN glob matching (not a mock of it) is what is tested.
func newMiniredisDeleter(
	t *testing.T,
	mockController *gomock.Controller,
	retentions map[string]time.Duration,
) (*miniredis.Miniredis, cache.MultiGetDeleteCountCache, *redisCounterDeleter) {
	t.Helper()
	mr := miniredis.RunT(t)
	client, err := redis.NewClient(
		mr.Addr(),
		redis.WithRedisMode(redis.RedisModeStandalone),
		redis.WithLogger(zap.NewNop()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	envClient := evmock.NewMockClient(mockController)
	envClient.EXPECT().ListEnvironmentsV2(gomock.Any(), gomock.Any()).Return(
		&evproto.ListEnvironmentsV2Response{
			Environments: []*evproto.EnvironmentV2{{Id: "dev"}},
		}, nil).AnyTimes()

	redisCache := cachev3.NewRedisCache(client)
	return mr, redisCache, &redisCounterDeleter{
		envClient:  envClient,
		cache:      redisCache,
		retentions: retentions,
		opts: &jobs.Options{
			Timeout: 5 * time.Second,
		},
		logger: zap.NewNop(),
	}
}

// TestRunIgnoresTieredKinds proves the production-safety property of the
// minute/day tier kind names: the deleter's scan prefixes ("<env>:ec:*",
// "<env>:uc:*") can never match ec1m/uc1m/ec1d/uc1d keys. If they could, the
// timestamp regex would find no submatch on those keys and the whole job
// would abort with errSubmatchStringNotFound.
func TestRunIgnoresTieredKinds(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, deleter := newMiniredisDeleter(t, mockController, nil)

	oldTs := time.Now().AddDate(0, 0, -40).Unix()
	recentTs := time.Now().Add(-2 * time.Hour).Unix()

	oldHourlyKeys := []string{
		fmt.Sprintf("dev:ec:%d:feature-1:variation-A", oldTs),
		fmt.Sprintf("dev:uc:%d:feature-1:variation-A", oldTs),
	}
	recentHourlyKeys := []string{
		fmt.Sprintf("dev:ec:%d:feature-1:variation-A", recentTs),
		fmt.Sprintf("dev:uc:%d:feature-1:variation-A", recentTs),
	}
	// Tiered keys share the environment and even the old timestamp: only the
	// kind name protects them from the scan.
	tieredKeys := []string{
		fmt.Sprintf("dev:ec1m:%d:feature-1:variation-A", oldTs),
		fmt.Sprintf("dev:uc1m:%d:feature-1:variation-A", oldTs),
		fmt.Sprintf("dev:ec1d:%d:feature-1:variation-A", oldTs),
		fmt.Sprintf("dev:uc1d:%d:feature-1:variation-A", oldTs),
	}
	for _, key := range append(append(append([]string{}, oldHourlyKeys...), recentHourlyKeys...), tieredKeys...) {
		require.NoError(t, redisCache.Put(key, "1", 0))
	}

	// The scan itself must never return a tiered key.
	for _, kind := range kinds {
		scanned, err := deleter.scan("dev", kind, deleter.newKeyPrefix("dev", kind))
		require.NoError(t, err)
		for _, key := range scanned {
			for _, tiered := range tieredKeys {
				assert.NotEqual(t, tiered, key, "scan prefix for %q must not match tiered keys", kind)
			}
		}
	}

	// The full job must complete (no errSubmatchStringNotFound abort), delete
	// only the expired hourly keys, and leave every tiered key untouched.
	require.NoError(t, deleter.Run(context.Background()))
	for _, key := range oldHourlyKeys {
		assert.False(t, mr.Exists(key), "expired hourly key %s must be deleted", key)
	}
	for _, key := range recentHourlyKeys {
		assert.True(t, mr.Exists(key), "recent hourly key %s must be kept", key)
	}
	for _, key := range tieredKeys {
		assert.True(t, mr.Exists(key), "tiered key %s must be invisible to the deleter", key)
	}
}

func TestRunPerKindRetention(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	// Shrink only the ec retention to 25 hours; uc keeps the 30-day default.
	mr, redisCache, deleter := newMiniredisDeleter(t, mockController, map[string]time.Duration{
		"ec": 25 * time.Hour,
	})

	ts26hAgo := time.Now().Add(-26 * time.Hour).Unix()
	ts24hAgo := time.Now().Add(-24 * time.Hour).Unix()

	ecExpired := fmt.Sprintf("dev:ec:%d:feature-1:variation-A", ts26hAgo)
	ecKept := fmt.Sprintf("dev:ec:%d:feature-1:variation-A", ts24hAgo)
	ucKept := fmt.Sprintf("dev:uc:%d:feature-1:variation-A", ts26hAgo)
	for _, key := range []string{ecExpired, ecKept, ucKept} {
		require.NoError(t, redisCache.Put(key, "1", 0))
	}

	require.NoError(t, deleter.Run(context.Background()))
	assert.False(t, mr.Exists(ecExpired), "ec key older than the 25h retention must be deleted")
	assert.True(t, mr.Exists(ecKept), "ec key within the 25h retention must be kept")
	assert.True(t, mr.Exists(ucKept), "uc key must keep the default 30-day retention")
}
