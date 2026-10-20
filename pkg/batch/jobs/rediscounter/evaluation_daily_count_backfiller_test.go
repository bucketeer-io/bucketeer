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
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachemock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	evmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	redis "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	evproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

func newMiniredisBackfiller(
	t *testing.T,
	mockController *gomock.Controller,
) (*miniredis.Miniredis, cache.MultiGetDeleteCountCache, *evaluationDailyCountBackfiller) {
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
	return mr, redisCache, &evaluationDailyCountBackfiller{
		envClient: envClient,
		cache:     redisCache,
		opts: &jobs.Options{
			Timeout: 10 * time.Second,
		},
		logger: zap.NewNop(),
	}
}

func hourlyKey(kind string, dayTs int64, hour int) string {
	return fmt.Sprintf("dev:%s:%d:feature-1:variation-A", kind, dayTs+int64(hour)*3600)
}

func dailyKey(kind string, dayTs int64) string {
	return fmt.Sprintf("dev:%s:%d:feature-1:variation-A", kind, dayTs)
}

func TestBackfillerRequiresCutoverMarker(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	_, _, backfiller := newMiniredisBackfiller(t, mockController)
	err := backfiller.Run(context.Background())
	assert.ErrorIs(t, err, errCutoverMarkerMissing)
}

func TestBackfillerMaterializesDailyKeys(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	day1 := today - 2*daySeconds
	day2 := today - daySeconds
	oldDay := today - 40*daySeconds // outside the 30-day window

	// The persister wrote the marker on its (simulated) deploy day: days
	// before today are still served from the hourly fallback.
	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(today, 10), 0))

	// Hourly event counts: day1 = 3+5, day2 = 7. The current (partial) day
	// and days older than the window must not be materialized.
	for key, value := range map[string]string{
		hourlyKey("ec", day1, 0):   "3",
		hourlyKey("ec", day1, 23):  "5",
		hourlyKey("ec", day2, 12):  "7",
		hourlyKey("ec", today, 1):  "11",
		hourlyKey("ec", oldDay, 1): "13",
	} {
		require.NoError(t, redisCache.Put(key, value, 0))
	}
	// Hourly user counts. The users are disjoint across hours because
	// miniredis sums multi-key PFCOUNT instead of computing the union (real
	// Redis unions), which would trip the job's daily >= hourly-union
	// verification. Cross-hour deduplication is covered by the real-Redis
	// integration test in pkg/redis/v3.
	_, err := redisCache.PFAdd(hourlyKey("uc", day1, 0), "u1", "u2", "u1")
	require.NoError(t, err)
	_, err = redisCache.PFAdd(hourlyKey("uc", day1, 23), "u3")
	require.NoError(t, err)
	_, err = redisCache.PFAdd(hourlyKey("uc", day2, 12), "u4")
	require.NoError(t, err)
	// A staging key in the scan range must be skipped, not abort the job.
	_, err = redisCache.PFAdd("dev:uc:pfmerge-key:feature-1", "u9")
	require.NoError(t, err)

	require.NoError(t, backfiller.Run(context.Background()))

	// ec1d: sum of the day's hourly values, with the daily TTL.
	for _, p := range []struct {
		day      int64
		expected string
	}{
		{day1, "8"},
		{day2, "7"},
	} {
		key := dailyKey("ec1d", p.day)
		got, err := mr.Get(key)
		require.NoError(t, err, "daily key %s must exist", key)
		assert.Equal(t, p.expected, got)
		assert.Equal(t, cachev3.DailyCountTTL, mr.TTL(key))
	}
	// uc1d: union of the day's hourly HLLs.
	count, err := redisCache.PFCount(dailyKey("uc1d", day1))
	require.NoError(t, err)
	assert.Equal(t, int64(3), count, "daily key must union the day's hourly HLLs")
	count, err = redisCache.PFCount(dailyKey("uc1d", day2))
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, cachev3.DailyCountTTL, mr.TTL(dailyKey("uc1d", day1)))

	// Out-of-window and current-day daily keys must not exist.
	assert.False(t, mr.Exists(dailyKey("ec1d", oldDay)))
	assert.False(t, mr.Exists(dailyKey("ec1d", today)))

	// The cutover marker moved back to the earliest backfilled day.
	marker, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(today-backfillDays*daySeconds, 10), marker)

	// Re-running is idempotent: same values, no error, marker unchanged.
	require.NoError(t, backfiller.Run(context.Background()))
	got, err := mr.Get(dailyKey("ec1d", day1))
	require.NoError(t, err)
	assert.Equal(t, "8", got)
	count, err = redisCache.PFCount(dailyKey("uc1d", day1))
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
	marker, err = mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(today-backfillDays*daySeconds, 10), marker)
}

// Days on/after the cutover are re-verified against their hourly keys too:
// hourly and daily are both written per event, so a live-written daily key
// that agrees with the hourly sum is rewritten with the same value, and one
// that drifted is reconciled. This is what makes a re-run a recovery path.
// End to end through Run with several environments sharing one Redis: the
// single cluster-wide scan must route every (day, suffix) group to its own
// environment's daily key (never mixing environments that share a feature
// ID), handle the prefix-less admin environment even though the environment
// service never lists it, and leave alone the keys of any other environment
// Redis knows but the environment service does not list.
func TestBackfillerRunMaterializesEveryEnvironmentFromOneScan(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr := miniredis.RunT(t)
	client, err := redis.NewClient(
		mr.Addr(),
		redis.WithRedisMode(redis.RedisModeStandalone),
		redis.WithLogger(zap.NewNop()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	redisCache := cachev3.NewRedisCache(client)

	envClient := evmock.NewMockClient(mockController)
	envClient.EXPECT().ListEnvironmentsV2(gomock.Any(), gomock.Any()).Return(
		&evproto.ListEnvironmentsV2Response{
			// As in production, the admin environment ("" / prefix-less
			// keys) is NOT listed; "ghost" exists only in Redis.
			Environments: []*evproto.EnvironmentV2{{Id: "dev"}, {Id: "prod"}},
		}, nil).AnyTimes()
	backfiller := &evaluationDailyCountBackfiller{
		envClient: envClient,
		cache:     redisCache,
		opts:      &jobs.Options{Timeout: 10 * time.Second},
		logger:    zap.NewNop(),
	}

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	day := today - 2*daySeconds
	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(today, 10), 0))

	key := func(env, kind string, hour int) string {
		return cache.MakeKey(kind, fmt.Sprintf("%d:feature-1:variation-A", day+int64(hour)*3600), env)
	}
	// Same feature/variation/day in every environment, different values.
	for k, v := range map[string]string{
		key("dev", "ec", 1):   "1",
		key("dev", "ec", 2):   "1",
		key("prod", "ec", 1):  "10",
		key("prod", "ec", 20): "10",
		key("", "ec", 5):      "100",
		key("ghost", "ec", 5): "1000",
	} {
		require.NoError(t, redisCache.Put(k, v, 0))
	}
	for k, users := range map[string][]string{
		key("dev", "uc", 1):   {"d1"},
		key("prod", "uc", 1):  {"p1", "p2"},
		key("prod", "uc", 20): {"p3"},
		key("", "uc", 5):      {"a1", "a2", "a3"},
		key("ghost", "uc", 5): {"g1"},
	} {
		_, err := redisCache.PFAdd(k, users...)
		require.NoError(t, err)
	}

	require.NoError(t, backfiller.Run(context.Background()))

	dailyOf := func(env, kind string) string {
		return cache.MakeKey(kind, fmt.Sprintf("%d:feature-1:variation-A", day), env)
	}
	for _, p := range []struct {
		env   string
		ec    string
		users int64
	}{
		{"dev", "2", 1},
		{"prod", "20", 3},
		{"", "100", 3},
	} {
		got, err := mr.Get(dailyOf(p.env, "ec1d"))
		require.NoError(t, err, "env %q", p.env)
		assert.Equal(t, p.ec, got, "ec1d of env %q", p.env)
		count, err := redisCache.PFCount(dailyOf(p.env, "uc1d"))
		require.NoError(t, err)
		assert.Equal(t, p.users, count, "uc1d of env %q", p.env)
	}
	assert.False(t, mr.Exists(dailyOf("ghost", "ec1d")), "unlisted environment is not backfilled")
	assert.False(t, mr.Exists(dailyOf("ghost", "uc1d")))
	assert.Equal(t, "1000", mustGet(t, mr, key("ghost", "ec", 5)), "and its hourly keys are untouched")

	marker, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(today-backfillDays*daySeconds, 10), marker)
}

func mustGet(t *testing.T, mr *miniredis.Miniredis, key string) string {
	t.Helper()
	v, err := mr.Get(key)
	require.NoError(t, err)
	return v
}

func TestNewEvaluationDailyCountBackfillerDefaultTimeoutFitsTriggerPath(t *testing.T) {
	t.Parallel()
	job := NewEvaluationDailyCountBackfiller(nil, nil)
	b, ok := job.(*evaluationDailyCountBackfiller)
	require.True(t, ok)
	// The manual trigger path (cronjob curl -m 3600, envoy route timeouts
	// 3600s) must never be the deadline that fires: the job's own must.
	assert.Equal(t, 55*time.Minute, b.opts.Timeout)
	assert.Less(t, b.opts.Timeout, time.Hour)
}

func TestBackfillerReverifiesDaysCoveredByCutover(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	cutover := today - 5*daySeconds
	consistentDay := today - 3*daySeconds // >= cutover, daily == hourly sum
	driftedDay := today - 4*daySeconds    // >= cutover, daily over-counts by one
	fallbackDay := today - 10*daySeconds  // < cutover: needs backfilling

	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(cutover, 10), 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", consistentDay, 2), "5", 0))
	require.NoError(t, redisCache.Put(dailyKey("ec1d", consistentDay), "5", 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", driftedDay, 2), "6", 0))
	require.NoError(t, redisCache.Put(dailyKey("ec1d", driftedDay), "7", 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", fallbackDay, 2), "21", 0))

	require.NoError(t, backfiller.Run(context.Background()))

	for _, p := range []struct {
		day      int64
		expected string
	}{
		{consistentDay, "5"},
		{driftedDay, "6"},
		{fallbackDay, "21"},
	} {
		got, err := mr.Get(dailyKey("ec1d", p.day))
		require.NoError(t, err)
		assert.Equal(t, p.expected, got, "day %d must equal its hourly sum", p.day)
	}

	marker, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(today-backfillDays*daySeconds, 10), marker)
}

// While the marker is in the future (rollout window) the job is a clean
// no-op, not a failure: it writes nothing, leaves the marker alone, and
// returns nil so the batch API does not report/retry it as an error.
func TestBackfillerSkipsWhileCutoverInFuture(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	// Deployment-day situation: the persister just started, so the marker
	// points at the NEXT UTC midnight. Backfilling now and moving the marker
	// back would hand today's reads to its incomplete daily key.
	futureCutover := today + daySeconds
	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(futureCutover, 10), 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", today-2*daySeconds, 1), "4", 0))

	require.NoError(t, backfiller.Run(context.Background()), "expected state, must not surface as a failed job")

	assert.False(t, mr.Exists(dailyKey("ec1d", today-2*daySeconds)), "nothing must be materialized")
	marker, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(futureCutover, 10), marker, "the marker must stay untouched")
}

func TestMoveCutoverBackNeverForward(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	current := today - 20*daySeconds
	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(current, 10), 0))

	// The marker is re-read right before writing, so a marker that was
	// already moved further back (e.g. by a concurrent run finishing first)
	// is never overwritten with a later day.
	require.NoError(t, backfiller.moveCutoverBack(today-10*daySeconds))
	marker, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(current, 10), marker, "a later day must never replace the marker")

	require.NoError(t, backfiller.moveCutoverBack(today-25*daySeconds))
	marker, err = mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(today-25*daySeconds, 10), marker, "an earlier day must move the marker back")
}

// The scan only discovers (day, feature:variation) groups; the keys read for a
// group are all 24 constructed hourly keys of that day, each exactly once,
// regardless of which (or how many times) keys were reported by SCAN.
func TestCollectHourlyKeysByDayExpandsToAllHoursExactlyOnce(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	day := today - 2*daySeconds
	keyA := hourlyKey("ec", day, 3)
	keyB := hourlyKey("ec", day, 7)

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	// SCAN may report the same key more than once (rehashing, or slot
	// migration across cluster nodes); counting it twice would inflate the
	// summed daily value.
	mockCache.EXPECT().ScanKeys("*ec:*", redisScanMaxSize).Return(
		[]string{keyA, keyA, keyB, "dev:ec:pfmerge-staging"}, nil)

	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	byEnv, err := backfiller.collectHourlyKeysByDay("ec", today-30*daySeconds, today)
	require.NoError(t, err)
	require.Len(t, byEnv, 1)
	groups := byEnv["dev"]
	require.Len(t, groups, 1)

	expected := make([]string, 0, 24)
	for hour := 0; hour < 24; hour++ {
		expected = append(expected, hourlyKey("ec", day, hour))
	}
	got := groups[dayGroup{dayTs: day, suffix: "feature-1:variation-A"}]
	assert.Equal(t, expected, got, "all 24 hourly keys of the day, in order, each exactly once")
	assert.Contains(t, got, keyA)
	assert.Contains(t, got, keyB)
}

// An hourly key created after the SCAN (a late event for an hour that had no
// counter when the job started) must still be counted: the job reads the
// constructed 24 hourly keys, not the scanned snapshot.
func TestBackfillEventCountDayIncludesHourlyKeysCreatedAfterScan(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	day := today - 2*daySeconds
	dKey := dailyKey("ec1d", day)
	allHours := make([]string, 0, 24)
	for hour := 0; hour < 24; hour++ {
		allHours = append(allHours, hourlyKey("ec", day, hour))
	}
	// Values for the 24 hourly keys: hour 0 existed at scan time (5), hour 1
	// is created between the scan and the read (1), the rest are absent.
	snapshot := func(hour1 interface{}) []interface{} {
		values := make([]interface{}, 24)
		values[0] = "5"
		values[1] = hour1
		return values
	}

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	gomock.InOrder(
		// Only hour 0 existed when the scan ran.
		mockCache.EXPECT().ScanKeys("*ec:*", redisScanMaxSize).Return(
			[]string{hourlyKey("ec", day, 0)}, nil),
		// By the time the day is read, hour 1 exists too.
		mockCache.EXPECT().GetMulti(allHours, true).Return(snapshot("1"), nil),
		mockCache.EXPECT().Get(dKey).Return(nil, cache.ErrNotFound),
		mockCache.EXPECT().Put(dKey, "6", cachev3.DailyCountTTL).Return(nil),
		mockCache.EXPECT().GetMulti(allHours, true).Return(snapshot("1"), nil),
		mockCache.EXPECT().Get(dKey).Return("6", nil),
	)

	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	byEnv, err := backfiller.collectHourlyKeysByDay("ec", today-30*daySeconds, today)
	require.NoError(t, err)
	groups := byEnv["dev"]
	require.Len(t, groups, 1)
	hourlyKeys := groups[dayGroup{dayTs: day, suffix: "feature-1:variation-A"}]
	require.NoError(t, backfiller.backfillEventCountDay(context.Background(), dKey, hourlyKeys))
}

// One cluster-wide SCAN per kind discovers every environment's days at once
// (including the admin environment, whose keys have no "<env>:" prefix),
// and only the keys of the requested kind: other tiers, other kinds, staging
// keys, and environments whose ID merely ends in the kind name are ignored.
func TestCollectHourlyKeysByDayGroupsAllEnvironmentsFromOneScan(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	day := today - 2*daySeconds
	hour := day + 5*3600

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	// Exactly one SCAN, with the cluster-wide pattern, however many
	// environments exist.
	mockCache.EXPECT().ScanKeys("*ec:*", redisScanMaxSize).Times(1).Return([]string{
		fmt.Sprintf("dev:ec:%d:feature-1:variation-A", hour),
		fmt.Sprintf("prod:ec:%d:feature-1:variation-A", hour),
		fmt.Sprintf("prod:ec:%d:feature-2:variation-B", hour+3600),
		fmt.Sprintf("ec:%d:feature-1:variation-A", hour), // admin environment
		fmt.Sprintf("dev:uc:%d:feature-1:variation-A", hour),
		fmt.Sprintf("dev:ec1d:%d:feature-1:variation-A", day),
		fmt.Sprintf("dev:ec1m:%d:feature-1:variation-A", hour),
		fmt.Sprintf("myspec:%d:feature-1:variation-A", hour), // env "myspec" has no kind segment
		"dev:ec:pfmerge-staging",
		"system:ec1d:cutover",
	}, nil)

	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	byEnv, err := backfiller.collectHourlyKeysByDay("ec", today-30*daySeconds, today)
	require.NoError(t, err)
	require.Len(t, byEnv, 3, "dev, prod and the admin environment")

	require.Len(t, byEnv["dev"], 1)
	require.Len(t, byEnv["prod"], 2)
	require.Len(t, byEnv[""], 1)

	adminKeys := byEnv[""][dayGroup{dayTs: day, suffix: "feature-1:variation-A"}]
	require.Len(t, adminKeys, 24)
	assert.Equal(t, fmt.Sprintf("ec:%d:feature-1:variation-A", day), adminKeys[0], "admin keys keep the prefix-less form")

	prodKeys := byEnv["prod"][dayGroup{dayTs: day, suffix: "feature-2:variation-B"}]
	require.Len(t, prodKeys, 24)
	assert.Equal(t, fmt.Sprintf("prod:ec:%d:feature-2:variation-B", day+23*3600), prodKeys[23])
}

func TestHourlyKeyRegexRejectsOtherKindsAndTiers(t *testing.T) {
	t.Parallel()

	re := hourlyKeyRegex("ec")
	patterns := []struct {
		key      string
		matches  bool
		env      string
		hourTs   string
		suffixOf string
	}{
		{key: "dev:ec:1792454400:f:v", matches: true, env: "dev", hourTs: "1792454400", suffixOf: "f:v"},
		{key: "ec:1792454400:f:v", matches: true, env: "", hourTs: "1792454400", suffixOf: "f:v"},
		{key: "ec:ec:1792454400:f:v", matches: true, env: "ec", hourTs: "1792454400", suffixOf: "f:v"},
		{key: "dev:uc:1792454400:f:v", matches: false},
		{key: "dev:ec1d:1792454400:f:v", matches: false},
		{key: "dev:ec1m:1792454400:f:v", matches: false},
		{key: "spec:1792454400:f:v", matches: false},
		{key: "dev:ec:notatimestamp:f:v", matches: false},
		{key: "dev:ec:1792454400", matches: false},
	}
	for _, p := range patterns {
		t.Run(p.key, func(t *testing.T) {
			m := re.FindStringSubmatch(p.key)
			if !p.matches {
				assert.Nil(t, m)
				return
			}
			require.Len(t, m, 4)
			assert.Equal(t, p.env, m[1])
			assert.Equal(t, p.hourTs, m[2])
			assert.Equal(t, p.suffixOf, m[3])
		})
	}
}

func TestBackfillEventCountDayStabilizesAgainstLiveWrites(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const dKey = "dev:ec1d:0:feature-1:variation-A"
	hourly := []string{"dev:ec:0:feature-1:variation-A"}

	// A late event increments the hourly key between the first snapshot and
	// the SET; the job must notice the changed sum and rewrite the daily key
	// before verifying.
	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	gomock.InOrder(
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"5"}, nil),
		mockCache.EXPECT().Get(dKey).Return(nil, cache.ErrNotFound),
		mockCache.EXPECT().Put(dKey, "5", cachev3.DailyCountTTL).Return(nil),
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"6"}, nil),
		mockCache.EXPECT().Put(dKey, "6", cachev3.DailyCountTTL).Return(nil),
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"6"}, nil),
		mockCache.EXPECT().Get(dKey).Return([]byte("6"), nil),
	)
	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	require.NoError(t, backfiller.backfillEventCountDay(context.Background(), dKey, hourly))
}

func TestBackfillEventCountDayGivesUpWhenNeverStable(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const dKey = "dev:ec1d:0:feature-1:variation-A"
	hourly := []string{"dev:ec:0:feature-1:variation-A"}

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	// The hourly sum changes on every read: after the bounded attempts the
	// job must fail (and be re-run) instead of publishing a stale snapshot.
	sum := 0
	mockCache.EXPECT().GetMulti(hourly, true).DoAndReturn(
		func(interface{}, bool) ([]interface{}, error) {
			sum++
			return []interface{}{strconv.Itoa(sum)}, nil
		}).Times(backfillEventCountDayAttempts + 1)
	mockCache.EXPECT().Get(dKey).Return(nil, cache.ErrNotFound)
	mockCache.EXPECT().Put(dKey, gomock.Any(), cachev3.DailyCountTTL).
		Return(nil).Times(backfillEventCountDayAttempts)

	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	err := backfiller.backfillEventCountDay(context.Background(), dKey, hourly)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "kept changing")
}

// Once a run has moved the marker to the window start, a later run must
// still do work: it is the only way to recover a day whose daily key drifted
// from its hourly keys after publication (the documented live-write race).
// The marker must never move forward in the process.
func TestBackfillerRerunRecoversDriftedDayAfterMarkerMoved(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	day := today - 2*daySeconds
	// Marker already before the window start (a previous run published it).
	marker := today - (backfillDays+5)*daySeconds
	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(marker, 10), 0))

	// ec1d over-counts by one (race: the daily INCRBY landed after the
	// previous run's SET + verification).
	require.NoError(t, redisCache.Put(hourlyKey("ec", day, 1), "4", 0))
	require.NoError(t, redisCache.Put(dailyKey("ec1d", day), "5", 0))
	// uc1d misses a user (race: a live PFADD to uc1d was overwritten by the
	// previous run's PFMERGE while the hourly PFADD landed after it).
	_, err := redisCache.PFAdd(hourlyKey("uc", day, 1), "u1", "u2")
	require.NoError(t, err)
	_, err = redisCache.PFAdd(dailyKey("uc1d", day), "u1")
	require.NoError(t, err)

	require.NoError(t, backfiller.Run(context.Background()))

	got, err := mr.Get(dailyKey("ec1d", day))
	require.NoError(t, err)
	assert.Equal(t, "4", got, "re-run must lower the over-counted daily key to the complete hourly sum")
	count, err := redisCache.PFCount(dailyKey("uc1d", day))
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "re-run must merge the missing user into uc1d")

	markerAfter, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(marker, 10), markerAfter, "the marker must never move forward")
}

// The backfill window must stay strictly inside the deleter's hourly
// retention: the deleter removes hourly keys older than a rolling
// `now - retention`, so the day that started exactly retention days ago is
// already losing hours while the job runs, and summing it would publish a
// partial total. Every hour of the oldest backfilled day must survive the
// deleter until the next UTC midnight.
func TestBackfillWindowIsInsideHourlyRetention(t *testing.T) {
	t.Parallel()

	assert.Equal(t, int64(29), backfillDays)

	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour).Unix()
	oldestDay := today - backfillDays*daySeconds
	nextMidnight := today + daySeconds

	deleterCutoff := func(at int64) int64 { return at - int64(defaultCounterRetention/time.Second) }
	// At any point before the next midnight the deleter's cutoff is still
	// before the first hour of the oldest day.
	assert.Less(t, deleterCutoff(nextMidnight-1), oldestDay,
		"deleter must not be able to remove any hour of the oldest backfilled day during the run")
	// Whereas a 30-day window would include a day whose early hours are
	// already deletable now (unless the run starts exactly at midnight).
	thirtyDaysAgo := today - 30*daySeconds
	assert.GreaterOrEqual(t, deleterCutoff(now.Unix()), thirtyDaysAgo,
		"the day retention days ago is already (partially) deletable")
}

// Within the backfill window every hourly key is retained, so the hourly sum
// is authoritative: a daily key that still exceeds it after a re-read is an
// over-count (live-write race) and is lowered. This is the recovery path for
// the documented race; it must not be mistaken for deleted history.
func TestBackfillEventCountDayLowersOverCountedDailyKey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const dKey = "dev:ec1d:0:feature-1:variation-A"
	hourly := []string{"dev:ec:0:feature-1:variation-A", "dev:ec:3600:feature-1:variation-A"}

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	gomock.InOrder(
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"5", "6"}, nil),
		mockCache.EXPECT().Get(dKey).Return([]byte("12"), nil),
		// Re-read confirms hourly is genuinely 11: the daily key over-counts.
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"5", "6"}, nil),
		mockCache.EXPECT().Put(dKey, "11", cachev3.DailyCountTTL).Return(nil),
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"5", "6"}, nil),
		mockCache.EXPECT().Get(dKey).Return([]byte("11"), nil),
	)

	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	require.NoError(t, backfiller.backfillEventCountDay(context.Background(), dKey, hourly))
}

// A daily value above the hourly snapshot can also mean the snapshot is
// stale: late events landed in both tiers between the hourly read and the
// daily read. The job must re-read and then reconcile to the true total (3),
// not keep the intermediate daily value (2).
func TestBackfillEventCountDayReconcilesStaleSnapshotAboveDaily(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const dKey = "dev:ec1d:0:feature-1:variation-A"
	hourly := []string{"dev:ec:0:feature-1:variation-A"}

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	gomock.InOrder(
		// Pre-deployment event: hourly=1, daily absent.
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"1"}, nil),
		// Two late events land in both tiers before the daily read: daily=2.
		mockCache.EXPECT().Get(dKey).Return([]byte("2"), nil),
		// Re-read: hourly has caught up to 3.
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"3"}, nil),
		mockCache.EXPECT().Put(dKey, "3", cachev3.DailyCountTTL).Return(nil),
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"3"}, nil),
		mockCache.EXPECT().Get(dKey).Return([]byte("3"), nil),
	)
	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	require.NoError(t, backfiller.backfillEventCountDay(context.Background(), dKey, hourly))
}

// A run that may cross UTC midnight (deadline on the next day) must choose
// its oldest day from the deadline, because the midnight deleter removes
// hour 00 of the day that is oldest-safe only for the start day.
func TestBackfillerWindowUsesRunDeadlineAcrossMidnight(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour).Unix()
	nextMidnight := time.Unix(today+daySeconds, 0).UTC()
	oldestForStartDay := today - backfillDays*daySeconds // unsafe if the run crosses midnight
	oldestForDeadline := oldestForStartDay + daySeconds

	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(today, 10), 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", oldestForStartDay, 5), "7", 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", oldestForDeadline, 5), "9", 0))

	// Effective deadline just after the next UTC midnight: the run may
	// straddle it. (Run caps the parent deadline with its own timeout, so the
	// job timeout must reach past midnight too.)
	backfiller.opts.Timeout = 48 * time.Hour
	ctx, cancel := context.WithDeadline(context.Background(), nextMidnight.Add(10*time.Minute))
	defer cancel()
	require.NoError(t, backfiller.Run(ctx))

	assert.False(t, mr.Exists(dailyKey("ec1d", oldestForStartDay)),
		"the day whose hour 00 the midnight deleter may remove mid-run must be excluded")
	assert.True(t, mr.Exists(dailyKey("ec1d", oldestForDeadline)))
	marker, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(oldestForDeadline, 10), marker,
		"cutover must not move over the excluded day")
}

// If the run overshoots the deadline its window was derived from, the marker
// must not be published.
func TestBackfillerDoesNotPublishCutoverPastDeadline(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	_, redisCache, backfiller := newMiniredisBackfiller(t, mockController)
	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(today, 10), 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", today-2*daySeconds, 5), "7", 0))

	// A deadline already in the past (but not yet enforced through ctx.Done
	// because Run installs its own timeout on top).
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	err := backfiller.Run(ctx)
	require.Error(t, err)
	marker, getErr := redisCache.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, getErr)
	assert.Equal(t, strconv.FormatInt(today, 10), string(marker.([]byte)), "marker unchanged")
}

// A smaller existing daily count (live-written partial day before the
// cutover, or late events since the last run) is still raised to the sum.
func TestBackfillEventCountDayRaisesSmallerExistingDailyCount(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const dKey = "dev:ec1d:0:feature-1:variation-A"
	hourly := []string{"dev:ec:0:feature-1:variation-A"}

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	gomock.InOrder(
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"24"}, nil),
		mockCache.EXPECT().Get(dKey).Return([]byte("11"), nil),
		mockCache.EXPECT().Put(dKey, "24", cachev3.DailyCountTTL).Return(nil),
		mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"24"}, nil),
		mockCache.EXPECT().Get(dKey).Return([]byte("24"), nil),
	)
	backfiller := &evaluationDailyCountBackfiller{
		cache:  mockCache,
		logger: zap.NewNop(),
	}
	require.NoError(t, backfiller.backfillEventCountDay(context.Background(), dKey, hourly))
}

// A hourly read that is delayed across both the job deadline and UTC midnight
// can observe a day the midnight deleter has already trimmed. The partial sum
// must not be written: the day was already published at its complete value
// and the marker covers it, so a write would hand readers the lower total.
// The deadline is checked before every daily write, not only before moving
// the marker.
func TestBackfillEventCountDayDoesNotWritePastDeadline(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const dKey = "dev:ec1d:0:feature-1:variation-A"
	hourly := []string{"dev:ec:0:feature-1:variation-A", "dev:ec:3600:feature-1:variation-A"}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
	defer cancel()

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	// The read crosses the deadline (and, in production, midnight: hour 00
	// is gone, so the sum is 23 instead of the published 24).
	mockCache.EXPECT().GetMulti(hourly, true).DoAndReturn(
		func(interface{}, bool) ([]interface{}, error) {
			<-ctx.Done()
			return []interface{}{"11", "12"}, nil
		})
	mockCache.EXPECT().Get(dKey).Return([]byte("24"), nil)
	// Re-read (daily > sum) still sees the trimmed day.
	mockCache.EXPECT().GetMulti(hourly, true).Return([]interface{}{"11", "12"}, nil)
	mockCache.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	backfiller := &evaluationDailyCountBackfiller{cache: mockCache, logger: zap.NewNop()}
	err := backfiller.backfillEventCountDay(ctx, dKey, hourly)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "deadline")
}

func TestBackfillUserCountDayDoesNotWritePastDeadline(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	mockCache := cachemock.NewMockMultiGetDeleteCountCache(mockController)
	mockCache.EXPECT().PFMerge(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	backfiller := &evaluationDailyCountBackfiller{cache: mockCache, logger: zap.NewNop()}
	err := backfiller.backfillUserCountDay(ctx, "dev:uc1d:0:feature-1:variation-A",
		[]string{"dev:uc:0:feature-1:variation-A"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "deadline")
}
