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

func TestBackfillerSkipsDaysCoveredByCutover(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	cutover := today - 5*daySeconds
	coveredDay := today - 3*daySeconds   // >= cutover: live daily key is authoritative
	fallbackDay := today - 10*daySeconds // < cutover: needs backfilling

	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(cutover, 10), 0))

	// The covered day has hourly data AND a live-written daily key that the
	// backfill must not clobber.
	require.NoError(t, redisCache.Put(hourlyKey("ec", coveredDay, 2), "5", 0))
	require.NoError(t, redisCache.Put(dailyKey("ec1d", coveredDay), "999", 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", fallbackDay, 2), "21", 0))

	require.NoError(t, backfiller.Run(context.Background()))

	got, err := mr.Get(dailyKey("ec1d", coveredDay))
	require.NoError(t, err)
	assert.Equal(t, "999", got, "days on/after the cutover must not be rewritten")

	got, err = mr.Get(dailyKey("ec1d", fallbackDay))
	require.NoError(t, err)
	assert.Equal(t, "21", got, "days before the cutover must be materialized")

	marker, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(today-backfillDays*daySeconds, 10), marker)
}

func TestBackfillerNoWorkWhenMarkerCoversWindow(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mr, redisCache, backfiller := newMiniredisBackfiller(t, mockController)

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	// Marker already at (or before) the window start: the job is a no-op and
	// must never move the marker forward.
	marker := today - (backfillDays+5)*daySeconds
	require.NoError(t, redisCache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(marker, 10), 0))
	require.NoError(t, redisCache.Put(hourlyKey("ec", today-2*daySeconds, 1), "4", 0))

	require.NoError(t, backfiller.Run(context.Background()))

	assert.False(t, mr.Exists(dailyKey("ec1d", today-2*daySeconds)),
		"days already covered by the marker must not be materialized")
	got, err := mr.Get(cachev3.EvaluationDailyCutoverKey)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(marker, 10), got, "the marker must never move forward")
}
