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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const (
	// backfillDays is how far back daily keys are materialized: the last 30
	// fully elapsed UTC days, matching the hourly retention horizon (hourly
	// data older than that no longer exists anywhere).
	backfillDays = 30
	daySeconds   = int64(24 * 60 * 60)
)

var (
	errCutoverMarkerMissing = errors.New(
		"batch: daily cutover marker missing; deploy the dual-writing persister before backfilling")
)

// evaluationDailyCountBackfiller is the one-shot, non-destructive migration
// job of the evaluation timeseries key redesign (§3.3). For every
// environment, feature, and variation it materializes the daily-tier keys
// from the existing hourly keys, for the last 30 fully elapsed UTC days that
// are still served by the hourly fallback:
//
//   - ec1d: SET of the summed hourly ec values (recomputed from data that is
//     immutable once the day has fully elapsed, so re-running is safe).
//   - uc1d: PFMERGE of the day's hourly uc HLLs (an idempotent union, safe to
//     re-run and safe under concurrent live writes).
//
// After every materialized day passes verification (ec1d equals the hourly
// sum, PFCOUNT(uc1d) covers the hourly union), the job moves the daily
// cutover marker back to the earliest backfilled day so the multi-day views
// are served purely from daily keys. If verification fails the marker is not
// moved and the job can simply be re-run.
//
// The job becomes useless (and exits without work) once the cutover marker
// already covers the whole backfill window, which happens naturally ~30 days
// after the persister rollout even if the job never runs.
type evaluationDailyCountBackfiller struct {
	envClient envclient.Client
	cache     cache.MultiGetDeleteCountCache
	opts      *jobs.Options
	logger    *zap.Logger
}

func NewEvaluationDailyCountBackfiller(
	redis cache.MultiGetDeleteCountCache,
	environmentClient envclient.Client,
	opts ...jobs.Option) jobs.Job {

	dopts := &jobs.Options{
		Timeout: 30 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &evaluationDailyCountBackfiller{
		cache:     redis,
		envClient: environmentClient,
		opts:      dopts,
		logger:    dopts.Logger.Named("evaluation-daily-count-backfiller"),
	}
}

func (b *evaluationDailyCountBackfiller) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobEvaluationDailyCountBackfiller, lastErr, time.Since(startTime))
	}()

	ctx, cancel := context.WithTimeout(ctx, b.opts.Timeout)
	defer cancel()

	cutover, err := b.getCutoverMarker()
	if err != nil {
		return err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour).Unix()
	earliest := today - backfillDays*daySeconds
	// Only fully elapsed days that the reader still serves from the hourly
	// fallback (days strictly before the cutover) need materializing. Days on
	// or after the cutover already have complete live-written daily keys.
	endExclusive := min64(today, cutover)
	if earliest >= endExclusive {
		b.logger.Info("Nothing to backfill: the cutover marker already covers the backfill window",
			zap.Int64("cutover", cutover),
			zap.Int64("earliestBackfillDay", earliest),
		)
		return nil
	}

	b.logger.Info("Starting the evaluation daily count backfill",
		zap.Int64("cutover", cutover),
		zap.Int64("earliestBackfillDay", earliest),
		zap.Int64("endExclusive", endExclusive),
	)

	envs, err := b.listEnvironments(ctx)
	if err != nil {
		return err
	}
	backfilled := 0
	for _, env := range envs {
		if err := ctx.Err(); err != nil {
			return err
		}
		n, err := b.backfillEnvironment(ctx, env.Id, earliest, endExclusive)
		if err != nil {
			return err
		}
		backfilled += n
	}

	// Every materialized day passed verification: move the cutover marker
	// back to the earliest backfilled day (never forward).
	if earliest < cutover {
		if err := b.cache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(earliest, 10), 0); err != nil {
			return fmt.Errorf("batch: failed to move the daily cutover marker back: %w", err)
		}
		b.logger.Info("Daily cutover marker moved back",
			zap.Int64("from", cutover),
			zap.Int64("to", earliest),
		)
	}

	b.logger.Info("Finished the evaluation daily count backfill",
		zap.Duration("elapsedTime", time.Since(startTime)),
		zap.Int("backfilledDailyKeys", backfilled),
	)
	return nil
}

// getCutoverMarker reads the marker written by the persister on startup. The
// backfill requires it: without the dual-writing persister deployed, freshly
// materialized daily keys would immediately go stale.
func (b *evaluationDailyCountBackfiller) getCutoverMarker() (int64, error) {
	value, err := b.cache.Get(cachev3.EvaluationDailyCutoverKey)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return 0, errCutoverMarkerMissing
		}
		return 0, err
	}
	bytes, err := cache.Bytes(value)
	if err != nil {
		return 0, fmt.Errorf("batch: unexpected daily cutover marker type: %w", err)
	}
	cutover, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("batch: failed to parse the daily cutover marker: %w", err)
	}
	return cutover, nil
}

func (b *evaluationDailyCountBackfiller) listEnvironments(ctx context.Context) ([]*envproto.EnvironmentV2, error) {
	resp, err := b.envClient.ListEnvironmentsV2(ctx, &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
		Archived: &wrapperspb.BoolValue{Value: false},
	})
	if err != nil {
		b.logger.Error("Failed to list environments", zap.Error(err))
		return nil, err
	}
	return resp.Environments, nil
}

// backfillEnvironment materializes and verifies the daily keys of one
// environment. It returns the number of materialized (feature, variation,
// day, kind) daily keys.
func (b *evaluationDailyCountBackfiller) backfillEnvironment(
	ctx context.Context,
	environmentId string,
	earliest, endExclusive int64,
) (int, error) {
	backfilled := 0
	for _, pair := range []struct {
		hourlyKind string
		dailyKind  string
	}{
		{cachev3.EventCountHourKind, cachev3.EventCountDailyKind},
		{cachev3.UserCountHourKind, cachev3.UserCountDailyKind},
	} {
		days, err := b.collectHourlyKeysByDay(environmentId, pair.hourlyKind, earliest, endExclusive)
		if err != nil {
			return backfilled, err
		}
		for group, hourlyKeys := range days {
			if err := ctx.Err(); err != nil {
				return backfilled, err
			}
			dailyKey := cache.MakeKey(
				pair.dailyKind,
				fmt.Sprintf("%d:%s", group.dayTs, group.suffix),
				environmentId,
			)
			var err error
			if pair.hourlyKind == cachev3.EventCountHourKind {
				err = b.backfillEventCountDay(dailyKey, hourlyKeys)
			} else {
				err = b.backfillUserCountDay(dailyKey, hourlyKeys)
			}
			if err != nil {
				b.logger.Error("Failed to backfill daily key",
					zap.Error(err),
					zap.String("environmentId", environmentId),
					zap.String("dailyKey", dailyKey),
					zap.Strings("hourlyKeys", hourlyKeys),
				)
				return backfilled, err
			}
			backfilled++
		}
	}
	return backfilled, nil
}

// dayGroup identifies one daily key to materialize: the UTC day bucket plus
// the "<featureId>:<variationId>" key suffix.
type dayGroup struct {
	dayTs  int64
	suffix string
}

// collectHourlyKeysByDay scans the environment's hourly keys of one kind and
// groups the existing keys by (UTC day, feature:variation suffix), limited to
// days in [earliest, endExclusive). Days without any hourly key get no entry:
// both the daily key and the hourly fallback correctly yield 0 for them.
func (b *evaluationDailyCountBackfiller) collectHourlyKeysByDay(
	environmentId, kind string,
	earliest, endExclusive int64,
) (map[dayGroup][]string, error) {
	pattern := cache.MakeKeyPrefix(kind, environmentId) + "*"
	var keyRegex *regexp.Regexp
	if environmentId == "" {
		keyRegex = regexp.MustCompile(fmt.Sprintf("^%s:([0-9]+):(.+)$", kind))
	} else {
		keyRegex = regexp.MustCompile(fmt.Sprintf("^%s:%s:([0-9]+):(.+)$", environmentId, kind))
	}

	groups := make(map[dayGroup][]string)
	var cursor uint64
	for {
		nextCursor, keys, err := b.cache.Scan(cursor, pattern, redisScanMaxSize)
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			match := keyRegex.FindStringSubmatch(key)
			if len(match) != 3 {
				// Not an hourly counter key (e.g. a pfmerge staging key).
				continue
			}
			hourTs, err := strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				continue
			}
			dayTs := (hourTs / daySeconds) * daySeconds
			if dayTs < earliest || dayTs >= endExclusive {
				continue
			}
			group := dayGroup{dayTs: dayTs, suffix: match[2]}
			groups[group] = append(groups[group], key)
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return groups, nil
}

// backfillEventCountDay sets the daily ec key to the sum of the day's hourly
// values and verifies the write by reading it back. The day has fully
// elapsed, so the hourly source data is immutable and re-running recomputes
// the same value.
func (b *evaluationDailyCountBackfiller) backfillEventCountDay(dailyKey string, hourlyKeys []string) error {
	values, err := b.cache.GetMulti(hourlyKeys, true)
	if err != nil {
		return fmt.Errorf("get hourly event counts: %w", err)
	}
	var sum int64
	for _, v := range values {
		s, ok := toString(v)
		if !ok || s == "" {
			continue
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("parse hourly event count: %w", err)
		}
		sum += n
	}
	if err := b.cache.Put(dailyKey, strconv.FormatInt(sum, 10), cachev3.DailyCountTTL); err != nil {
		return fmt.Errorf("set daily event count: %w", err)
	}
	// Verification: the daily key must hold exactly the hourly sum.
	stored, err := b.cache.Get(dailyKey)
	if err != nil {
		return fmt.Errorf("read back daily event count: %w", err)
	}
	s, ok := toString(stored)
	if !ok {
		return fmt.Errorf("batch: daily event count has an unexpected type: %v", stored)
	}
	got, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("parse daily event count: %w", err)
	}
	if got != sum {
		return fmt.Errorf("batch: daily event count verification failed: key %s has %d, hourly sum is %d",
			dailyKey, got, sum)
	}
	return nil
}

// backfillUserCountDay merges the day's hourly HLLs into the daily uc key and
// verifies that the daily union covers the hourly union. PFMERGE is an
// idempotent union: re-running, and racing with live writes to the daily key,
// can only grow it, never lose users.
func (b *evaluationDailyCountBackfiller) backfillUserCountDay(dailyKey string, hourlyKeys []string) error {
	if err := b.cache.PFMerge(dailyKey, cachev3.DailyCountTTL, hourlyKeys...); err != nil {
		return fmt.Errorf("merge hourly user counts: %w", err)
	}
	// Verification: merging HLLs produces exactly the union registers, so the
	// daily count can never be below the hourly union. It can exceed it when
	// the daily key already contained live-written users (late events).
	hourlyCount, err := b.cache.PFCount(hourlyKeys...)
	if err != nil {
		return fmt.Errorf("count hourly users: %w", err)
	}
	dailyCount, err := b.cache.PFCount(dailyKey)
	if err != nil {
		return fmt.Errorf("count daily users: %w", err)
	}
	if dailyCount < hourlyCount {
		return fmt.Errorf("batch: daily user count verification failed: key %s has %d, hourly union is %d",
			dailyKey, dailyCount, hourlyCount)
	}
	return nil
}

func toString(v interface{}) (string, bool) {
	switch v := v.(type) {
	case []byte:
		return string(v), true
	case string:
		return v, true
	case nil:
		return "", true
	default:
		return "", false
	}
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
