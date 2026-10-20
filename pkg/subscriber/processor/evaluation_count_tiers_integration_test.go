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
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	redis "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
)

// TestClusterEvaluationCountTiersIntegration is the end-to-end verification
// of the evaluation timeseries key redesign against a real Redis cluster
// (miniredis cannot round-trip HLL blobs and sums multi-key PFCOUNT instead
// of unioning it). It writes events through the real persister key-pair and
// flush functions, materializes "backfilled" days with the same cache
// operations the backfill job uses, and asserts:
//
//  1. For a sampled live day: ec1d == sum of the day's 24 hourly ec keys and
//     PFCOUNT(uc1d) == PFCOUNT of the day's 24 hourly uc keys.
//  2. A 7-day read mixing backfilled and live daily keys equals the pure
//     hourly computation, including the uc total as ONE PFCount over the
//     mixed key set.
//
// Skipped unless REDIS_CLUSTER_ADDR points at a cluster node, e.g.:
//
//	docker run -d --rm --name bkt-redis-cluster -e IP=0.0.0.0 \
//	  -p 7000-7005:7000-7005 grokzen/redis-cluster:7.0.10
//	REDIS_CLUSTER_ADDR=localhost:7000 go test ./pkg/subscriber/processor/ \
//	  -run TestClusterEvaluationCountTiersIntegration -v -count=1
func TestClusterEvaluationCountTiersIntegration(t *testing.T) {
	addr := os.Getenv("REDIS_CLUSTER_ADDR")
	if addr == "" {
		t.Skip("REDIS_CLUSTER_ADDR not set; skipping real-cluster integration test")
	}
	client, err := redis.NewClient(addr, redis.WithRedisMode(redis.RedisModeCluster))
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	redisCache := cachev3.NewRedisCache(client)

	// A unique environment isolates this run's keys; every written key is
	// tracked and deleted afterwards.
	environmentId := fmt.Sprintf("tiers-it-%d", time.Now().UnixNano())
	writtenKeys := make(map[string]struct{})
	track := func(keys ...string) {
		for _, k := range keys {
			writtenKeys[k] = struct{}{}
		}
	}
	t.Cleanup(func() {
		for k := range writtenKeys {
			_ = redisCache.Delete(k)
		}
	})

	persister := &evaluationCountEventPersister{
		evaluationCountCacher: redisCache,
		logger:                zap.NewNop(),
	}

	const (
		featureID   = "feature-it"
		variationID = "variation-a"
	)
	// 7 UTC days: the first 3 are "backfilled" (hourly history materialized
	// into daily keys the way the backfill job does), the last 4 are "live"
	// (written through the persister's dual-write path).
	baseDay := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	days := make([]time.Time, 7)
	for i := range days {
		days[i] = baseDay.AddDate(0, 0, i)
	}
	const liveDaysStart = 3

	// Per day: 6 events over hours 0, 11, and 23. "shared-user" repeats
	// within each day AND across all days, so correct deduplication across
	// hours, days, and tiers is observable: per-day uniques = 3, 7-day
	// uniques = 1 + 7*2 = 15.
	type hourEvents struct {
		hour  int
		users []string
	}
	dayEvents := func(day int) []hourEvents {
		u1 := fmt.Sprintf("day%d-u1", day)
		u2 := fmt.Sprintf("day%d-u2", day)
		return []hourEvents{
			{hour: 0, users: []string{"shared-user", u1}},
			{hour: 11, users: []string{u1, u2}},
			{hour: 23, users: []string{"shared-user", u2}},
		}
	}
	const (
		expectedDayEvents   = int64(6)
		expectedDayUniques  = int64(3)
		expectedWeekEvents  = 7 * expectedDayEvents
		expectedWeekUniques = int64(1 + 7*2)
	)

	hourlyECKeys := func(day time.Time) []string {
		keys := make([]string, 0, 24)
		for h := 0; h < 24; h++ {
			ts := day.Add(time.Duration(h) * time.Hour).Unix()
			keys = append(keys, newEvaluationCountKey(cachev3.EventCountHourKind, featureID, variationID, environmentId, ts))
		}
		return keys
	}
	hourlyUCKeys := func(day time.Time) []string {
		keys := make([]string, 0, 24)
		for h := 0; h < 24; h++ {
			ts := day.Add(time.Duration(h) * time.Hour).Unix()
			keys = append(keys, newEvaluationCountKey(cachev3.UserCountHourKind, featureID, variationID, environmentId, ts))
		}
		return keys
	}
	dailyECKey := func(day time.Time) string {
		return newEvaluationCountKey(cachev3.EventCountDailyKind, featureID, variationID, environmentId, day.Unix())
	}
	dailyUCKey := func(day time.Time) string {
		return newEvaluationCountKey(cachev3.UserCountDailyKind, featureID, variationID, environmentId, day.Unix())
	}

	// --- Backfilled days: pre-redesign hourly history only. ---
	for d := 0; d < liveDaysStart; d++ {
		for _, he := range dayEvents(d) {
			ts := days[d].Add(time.Duration(he.hour) * time.Hour).Unix()
			ecKey := newEvaluationCountKey(cachev3.EventCountHourKind, featureID, variationID, environmentId, ts)
			ucKey := newEvaluationCountKey(cachev3.UserCountHourKind, featureID, variationID, environmentId, ts)
			track(ecKey, ucKey)
			_, err := redisCache.PFAdd(ucKey, he.users...)
			require.NoError(t, err)
			_, err = redisCache.IncrementBy(ecKey, int64(len(he.users)))
			require.NoError(t, err)
		}
	}

	// --- Live days: events written through the real persister path. ---
	aggregator := newEvaluationCountAggregator()
	for d := liveDaysStart; d < len(days); d++ {
		for _, he := range dayEvents(d) {
			// Mid-hour timestamp: truncation to minute/hour/day is the
			// persister's job, not the test's.
			ts := days[d].Add(time.Duration(he.hour)*time.Hour + 42*time.Minute + 7*time.Second).Unix()
			for _, user := range he.users {
				for _, pair := range persister.newEvaluationCountKeyPairs(featureID, variationID, environmentId, ts, time.Unix(ts, 0)) {
					aggregator.addEvent(pair.ecKey, pair.ucKey, user)
					track(pair.ecKey, pair.ucKey)
				}
			}
		}
	}
	eventCounts, userCounts := aggregator.flush()
	failed, err := persister.flushAggregatedCounts(eventCounts, userCounts, nil)
	require.NoError(t, err)
	require.Empty(t, failed)

	// asString tolerates the two representations the cache returns ([]byte
	// from Get, string from GetMulti).
	asString := func(v interface{}) string {
		switch v := v.(type) {
		case []byte:
			return string(v)
		case string:
			return v
		case nil:
			return ""
		default:
			t.Fatalf("unexpected cache value type %T", v)
			return ""
		}
	}
	sumHourlyEC := func(day time.Time) int64 {
		values, err := redisCache.GetMulti(hourlyECKeys(day), true)
		require.NoError(t, err)
		var sum int64
		for _, v := range values {
			s := asString(v)
			if s == "" {
				continue
			}
			n, err := strconv.ParseInt(s, 10, 64)
			require.NoError(t, err)
			sum += n
		}
		return sum
	}

	// --- Backfill materialization: the same cache operations the backfill
	// job performs per (day, feature, variation) after its scan. ---
	for d := 0; d < liveDaysStart; d++ {
		sum := sumHourlyEC(days[d])
		track(dailyECKey(days[d]), dailyUCKey(days[d]))
		require.NoError(t, redisCache.Put(dailyECKey(days[d]), strconv.FormatInt(sum, 10), cachev3.DailyCountTTL))
		require.NoError(t, redisCache.PFMerge(dailyUCKey(days[d]), cachev3.DailyCountTTL, hourlyUCKeys(days[d])...))
	}

	getInt := func(key string) int64 {
		value, err := redisCache.Get(key)
		require.NoError(t, err, "key %s must exist", key)
		n, err := strconv.ParseInt(asString(value), 10, 64)
		require.NoError(t, err)
		return n
	}

	// --- Assertion 1: sampled days. The daily key must equal the hourly
	// computation, on a live (dual-written) and a backfilled day alike. ---
	for _, d := range []int{1 /* backfilled */, 5 /* live */} {
		day := days[d]
		assert.Equal(t, expectedDayEvents, sumHourlyEC(day), "day %d hourly ec sum", d)
		assert.Equal(t, expectedDayEvents, getInt(dailyECKey(day)), "day %d: ec1d must equal sum(24 hourly ec)", d)

		hourlyUniques, err := redisCache.PFCount(hourlyUCKeys(day)...)
		require.NoError(t, err)
		dailyUniques, err := redisCache.PFCount(dailyUCKey(day))
		require.NoError(t, err)
		assert.Equal(t, expectedDayUniques, hourlyUniques, "day %d hourly uc union", d)
		assert.Equal(t, hourlyUniques, dailyUniques, "day %d: PFCOUNT(uc1d) must equal PFCOUNT(24 hourly uc)", d)
	}

	// --- Assertion 2: the 7-day read mixing backfilled + live daily keys
	// must equal the pure-hourly computation. ---
	var weekECDaily, weekECHourly int64
	mixedUCKeys := make([]string, 0, len(days))
	allHourlyUCKeys := make([]string, 0, len(days)*24)
	for _, day := range days {
		weekECDaily += getInt(dailyECKey(day))
		weekECHourly += sumHourlyEC(day)
		mixedUCKeys = append(mixedUCKeys, dailyUCKey(day))
		allHourlyUCKeys = append(allHourlyUCKeys, hourlyUCKeys(day)...)
	}
	assert.Equal(t, expectedWeekEvents, weekECHourly, "7-day hourly ec total")
	assert.Equal(t, weekECHourly, weekECDaily, "7-day ec total from daily keys must equal the hourly computation")

	// The uc total is ONE PFCount over the mixed key set, exactly like the
	// reader's getTotalUserCounts.
	mixedUniques, err := redisCache.PFCount(mixedUCKeys...)
	require.NoError(t, err)
	hourlyWeekUniques, err := redisCache.PFCount(allHourlyUCKeys...)
	require.NoError(t, err)
	assert.Equal(t, expectedWeekUniques, hourlyWeekUniques, "7-day hourly uc union")
	assert.Equal(t, hourlyWeekUniques, mixedUniques,
		"7-day uc total over mixed daily keys must equal the pure-hourly union")
}
