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

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

const (
	// backfillDays is how far back daily keys are materialized. It is one day
	// LESS than the hourly retention of the redis counter deleter (30d → 29).
	//
	// The deleter's cutoff is a rolling instant (`now - 720h`), while this
	// job sums whole UTC days. Example: a run at 15:00 UTC has a cutoff of
	// 15:00 UTC 30 days ago, so the UTC day that started 30 days ago at
	// 00:00 has already lost hours 00-14 (or may lose them mid-run); summing
	// it would publish a partial total as the complete daily count. The day
	// that started 29 days ago only becomes deletable at the NEXT UTC
	// midnight, i.e. after this run, whatever time of day it executes.
	//
	// Month lengths are irrelevant: both retention and this window are
	// fixed 24h multiples, not calendar months.
	backfillDays = int64(defaultCounterRetention/(24*time.Hour)) - 1
	daySeconds   = int64(24 * 60 * 60)
)

var (
	errCutoverMarkerMissing = errors.New(
		"batch: daily cutover marker missing; deploy the dual-writing persister before backfilling")
)

// evaluationDailyCountBackfiller is the one-shot, non-destructive migration
// job of the evaluation timeseries key redesign (§3.3). For every
// environment, feature, and variation it materializes the daily-tier keys
// from the existing hourly keys, for the last 29 fully elapsed UTC days (one
// less than the hourly retention, see backfillDays) that
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
		Timeout: 55 * time.Minute,
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

	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour).Unix()
	// The oldest day is chosen from the run's DEADLINE, not its start: the
	// backfillDays window is only safe until the next UTC midnight, and a run
	// started shortly before midnight may still be reading its oldest day
	// when the midnight deleter removes that day's hour 00. Using the
	// deadline's day shifts the window one day later for such runs, so every
	// hour read is protected for the whole run (the context deadline aborts
	// the run before the marker moves if it overruns anyway).
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = now
	}
	earliest := deadline.UTC().Truncate(24*time.Hour).Unix() - backfillDays*daySeconds
	// Right after the persister's deployment the marker points at a FUTURE UTC
	// midnight (the rollout window): the daily keys of the current days only
	// cover events written by the new pods, so those days must keep being
	// served from the hourly fallback. Moving the marker back now would put an
	// incomplete daily key in charge, so skip until the marker's day has
	// started (then every day before it is fully elapsed and can be
	// materialized from the complete hourly keys). This is the expected state
	// for the first days after a deploy, not a failure: returning an error
	// here would surface as a failed job (alerts, retries by the caller).
	if cutover > today {
		b.logger.Info("Skipping the evaluation daily count backfill: the cutover marker is still in the future "+
			"(persister rollout window); nothing to do until that UTC day",
			zap.Int64("cutover", cutover),
			zap.Int64("today", today),
		)
		return nil
	}
	// Every fully elapsed day in the window is (re)derived from its hourly
	// keys, including days on/after the cutover whose daily keys were written
	// live. Those are expected to already match (both tiers are written per
	// event), so for them this is a verification that normally rewrites the
	// same value; but it is also the recovery path: a day whose daily key
	// drifted from its hourly keys (the documented race between a live
	// cross-node daily INCRBY and this job's SET, or a uc1d PFMERGE racing a
	// live PFADD) is reconciled by simply re-running the job while its hourly
	// keys are still retained. Restricting the range to days before the
	// cutover would make re-runs a no-op once the marker has moved. The
	// current (partial) day is excluded.
	endExclusive := today
	if earliest >= endExclusive {
		b.logger.Info("Nothing to backfill: no fully elapsed day in the window")
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
	// Environments whose daily keys are materialized: every environment the
	// environment service knows, plus the admin environment. Its keys have
	// no "<env>:" prefix (MakeKey with an empty environment ID) and the
	// environment API never lists it, yet the cutover marker is global, so
	// skipping it would leave any admin hourly history unreadable through
	// the daily path once the marker moves back. Hourly keys of any other
	// environment that only exists in Redis are left alone.
	environmentIDs := make([]string, 0, len(envs)+1)
	for _, env := range envs {
		environmentIDs = append(environmentIDs, env.Id)
	}
	environmentIDs = append(environmentIDs, "")
	backfilled := 0
	for _, pair := range countKindPairs {
		if err := ctx.Err(); err != nil {
			return err
		}
		// One cluster-wide SCAN per kind, not one per environment: SCAN
		// walks the whole keyspace whatever the MATCH pattern, so N
		// environments would cost N full walks of a multi-million-key
		// cluster (the discovery phase alone would exceed the job timeout).
		daysByEnv, err := b.collectHourlyKeysByDay(pair.hourlyKind, earliest, endExclusive)
		if err != nil {
			return err
		}
		for _, environmentID := range environmentIDs {
			if err := ctx.Err(); err != nil {
				return err
			}
			n, err := b.backfillEnvironment(ctx, environmentID, pair, daysByEnv[environmentID])
			if err != nil {
				return err
			}
			backfilled += n
		}
	}

	// Every materialized day passed verification: move the cutover marker
	// back to the earliest backfilled day (never forward). Publishing is the
	// one step that must not happen after the deadline that the window was
	// derived from, so check it explicitly rather than relying on ctx.Err()
	// having fired in between.
	if time.Now().After(deadline) {
		return fmt.Errorf("batch: backfill ran past its deadline %s; cutover marker left unchanged", deadline)
	}
	if err := b.moveCutoverBack(earliest); err != nil {
		return err
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

// moveCutoverBack moves the cutover marker back to earliest, but never
// forward: the marker is re-read right before writing so that a concurrent or
// earlier run that already moved it further back is not undone.
func (b *evaluationDailyCountBackfiller) moveCutoverBack(earliest int64) error {
	current, err := b.getCutoverMarker()
	if err != nil {
		return err
	}
	if earliest >= current {
		return nil
	}
	if err := b.cache.Put(cachev3.EvaluationDailyCutoverKey, strconv.FormatInt(earliest, 10), 0); err != nil {
		return fmt.Errorf("batch: failed to move the daily cutover marker back: %w", err)
	}
	b.logger.Info("Daily cutover marker moved back",
		zap.Int64("from", current),
		zap.Int64("to", earliest),
	)
	return nil
}

// listEnvironments deliberately includes archived environments: the cutover
// marker is global, so an environment that is unarchived later would
// otherwise be left with pre-cutover days that have no daily keys.
func (b *evaluationDailyCountBackfiller) listEnvironments(ctx context.Context) ([]*envproto.EnvironmentV2, error) {
	resp, err := b.envClient.ListEnvironmentsV2(ctx, &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
	})
	if err != nil {
		b.logger.Error("Failed to list environments", zap.Error(err))
		return nil, err
	}
	return resp.Environments, nil
}

// countKindPair maps an hourly kind to the daily kind materialized from it.
type countKindPair struct {
	hourlyKind string
	dailyKind  string
}

var countKindPairs = []countKindPair{
	{cachev3.EventCountHourKind, cachev3.EventCountDailyKind},
	{cachev3.UserCountHourKind, cachev3.UserCountDailyKind},
}

// backfillEnvironment materializes and verifies the daily keys of one kind
// for one environment, given that environment's (day, suffix) groups from
// collectHourlyKeysByDay. It returns the number of materialized (feature,
// variation, day) daily keys.
func (b *evaluationDailyCountBackfiller) backfillEnvironment(
	ctx context.Context,
	environmentId string,
	pair countKindPair,
	days map[dayGroup][]string,
) (int, error) {
	backfilled := 0
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
			err = b.backfillEventCountDay(ctx, dailyKey, hourlyKeys)
		} else {
			err = b.backfillUserCountDay(ctx, dailyKey, hourlyKeys)
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
	return backfilled, nil
}

// dayGroup identifies one daily key to materialize: the UTC day bucket plus
// the "<featureId>:<variationId>" key suffix.
type dayGroup struct {
	dayTs  int64
	suffix string
}

// hourlyKeyPattern is the SCAN MATCH pattern for the hourly keys of one kind
// across every environment: "<env>:<kind>:<hourTs>:<suffix>" and, for the
// admin environment, "<kind>:<hourTs>:<suffix>". It also matches keys whose
// environment ID merely ends in the kind name; hourlyKeyRegex rejects those.
func hourlyKeyPattern(kind string) string {
	return "*" + kind + ":*"
}

// hourlyKeyRegex parses "<env>:<kind>:<hourTs>:<suffix>" (or the admin form
// without "<env>:") into (env, hourTs, suffix). The optional env segment
// cannot contain ":".
func hourlyKeyRegex(kind string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^(?:([^:]+):)?%s:([0-9]+):(.+)$", regexp.QuoteMeta(kind)))
}

// collectHourlyKeysByDay scans the hourly keys of one kind, cluster-wide, to
// discover which (environment, UTC day, feature:variation suffix) groups
// exist in [earliest, endExclusive), and returns, per environment and per
// group, ALL 24 hourly keys of that day (constructed, not scanned). Days
// without any hourly key get no entry: both the daily key and the hourly
// fallback correctly yield 0 for them.
//
// A single SCAN serves every environment. Redis evaluates MATCH only after
// walking the keyspace, so the cost of a SCAN is the total number of keys in
// the cluster regardless of the pattern; scanning per environment would
// multiply that cost by the number of environments.
//
// The scan is only used for discovery. Reading the constructed 24 keys rather
// than the scanned ones means an hourly key that is created after the SCAN
// but before the reads (a late event for a day that previously had no
// counter for that hour) is still included; the stabilization loop in
// backfillEventCountDay then has a complete view to converge on. Missing
// hourly keys read as absent (ignoreNotFound / PFMERGE of a missing key is a
// no-op), so constructing keys that never existed is harmless.
//
// ScanKeys is cluster-exhaustive (it enumerates every master), and because
// groups are a set, a key reported more than once by SCAN cannot inflate the
// result.
func (b *evaluationDailyCountBackfiller) collectHourlyKeysByDay(
	kind string,
	earliest, endExclusive int64,
) (map[string]map[dayGroup][]string, error) {
	keyRegex := hourlyKeyRegex(kind)
	scanStart := time.Now()
	keys, err := b.cache.ScanKeys(hourlyKeyPattern(kind), redisScanMaxSize)
	if err != nil {
		return nil, err
	}
	byEnv := make(map[string]map[dayGroup][]string)
	groupCount := 0
	for _, key := range keys {
		match := keyRegex.FindStringSubmatch(key)
		if len(match) != 4 {
			// Not an hourly counter key (e.g. a pfmerge staging key).
			continue
		}
		environmentId := match[1]
		hourTs, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			continue
		}
		dayTs := (hourTs / daySeconds) * daySeconds
		if dayTs < earliest || dayTs >= endExclusive {
			continue
		}
		groups := byEnv[environmentId]
		if groups == nil {
			groups = make(map[dayGroup][]string)
			byEnv[environmentId] = groups
		}
		group := dayGroup{dayTs: dayTs, suffix: match[3]}
		if _, known := groups[group]; known {
			continue
		}
		groups[group] = hourlyKeysOfDay(kind, environmentId, group)
		groupCount++
	}
	b.logger.Info("Scanned hourly keys",
		zap.String("kind", kind),
		zap.Int("scannedKeys", len(keys)),
		zap.Int("environments", len(byEnv)),
		zap.Int("dailyKeysToBackfill", groupCount),
		zap.Duration("elapsedTime", time.Since(scanStart)),
	)
	return byEnv, nil
}

// hourlyKeysOfDay constructs the 24 hourly keys of one (day, feature:variation)
// group, in hour order.
func hourlyKeysOfDay(kind, environmentId string, group dayGroup) []string {
	keys := make([]string, 0, 24)
	for hour := int64(0); hour < 24; hour++ {
		keys = append(keys, cache.MakeKey(
			kind,
			fmt.Sprintf("%d:%s", group.dayTs+hour*3600, group.suffix),
			environmentId,
		))
	}
	return keys
}

// backfillEventCountDayAttempts bounds the snapshot-stabilization loop below.
const backfillEventCountDayAttempts = 3

// backfillEventCountDay sets the daily ec key to the sum of the day's hourly
// values. Although the day has elapsed, late events bucketed by their original
// timestamps can still land in it: the persister increments the hourly AND
// daily keys, and a late event racing between our hourly read and the SET
// would have its daily increment overwritten by the stale snapshot. To narrow
// that window the hourly sum is re-read after the SET and the write is
// repeated until it is stable, then verified by reading the daily key back.
//
// Known residual race (accepted): the persister's hourly and daily INCRBYs
// for one event travel in the same pipeline but hash to different cluster
// nodes, so there is no ordering between them. If the hourly INCRBY lands
// before our read and the daily INCRBY lands after our SET and verification,
// the day ends up over-counted by that one event (daily = sum + 1) and the
// stable-sum check cannot see it. Closing this entirely would require the
// live writers to be paused while the job runs; with a running persister
// the window is one cross-node pipeline skew (sub-millisecond to a few
// milliseconds) for a late event, days old, that must arrive while this
// exact (day, feature, variation) is being processed. Because the job
// re-derives EVERY day in the window from its hourly keys on every run
// (including days already on/after the cutover marker, see Run) and lowers
// a daily key that exceeds the complete hourly sum, a later re-run while the
// hourly keys are still retained re-converges any day affected by the race.
// Operators who need an exact migration can scale
// the evaluation count persister to zero for the duration of the run; events
// buffer in Pub/Sub and nothing is lost.
func (b *evaluationDailyCountBackfiller) backfillEventCountDay(
	ctx context.Context,
	dailyKey string,
	hourlyKeys []string,
) error {
	sum, err := b.sumHourlyEventCounts(hourlyKeys)
	if err != nil {
		return err
	}
	// An existing daily count above the hourly sum means either (a) our
	// hourly snapshot is stale (late events landed in both tiers between the
	// two reads) or (b) the daily key over-counts (the documented race with a
	// live cross-node daily INCRBY). Re-reading the hourly sum resolves (a):
	// the fresh sum has caught up and the normal stabilization below takes
	// over. If the daily value is still higher it is case (b), and the hourly
	// sum wins: within the backfill window every hourly key of the day is
	// guaranteed retained (backfillDays is derived from the deleter's
	// retention, and hourly keys never carry a TTL), so hourly is the
	// complete, authoritative source and lowering the daily key is the
	// recovery for the race. This is logged because it should be rare.
	existing, err := b.getDailyEventCount(dailyKey)
	if err != nil {
		return err
	}
	if existing > sum {
		sum, err = b.sumHourlyEventCounts(hourlyKeys)
		if err != nil {
			return err
		}
		if existing > sum {
			b.logger.Warn("Daily event count exceeds the complete hourly sum; lowering it to the hourly sum",
				zap.String("dailyKey", dailyKey),
				zap.Int64("existing", existing),
				zap.Int64("hourlySum", sum),
			)
		}
	}
	stable := false
	for attempt := 0; attempt < backfillEventCountDayAttempts; attempt++ {
		if err := checkDeadline(ctx); err != nil {
			return err
		}
		if err := b.cache.Put(dailyKey, strconv.FormatInt(sum, 10), cachev3.DailyCountTTL); err != nil {
			return fmt.Errorf("set daily event count: %w", err)
		}
		again, err := b.sumHourlyEventCounts(hourlyKeys)
		if err != nil {
			return err
		}
		if again == sum {
			stable = true
			break
		}
		sum = again
	}
	if !stable {
		return fmt.Errorf("batch: daily event count for key %s kept changing (live writes); re-run the backfill",
			dailyKey)
	}
	// Verification: the daily key must hold exactly the (stable) hourly sum.
	got, err := b.getDailyEventCount(dailyKey)
	if err != nil {
		return fmt.Errorf("read back daily event count: %w", err)
	}
	if got != sum {
		return fmt.Errorf("batch: daily event count verification failed: key %s has %d, hourly sum is %d",
			dailyKey, got, sum)
	}
	return nil
}

// checkDeadline guards every write to a daily key. The hourly sum is only
// authoritative while the run is within the deadline its window was derived
// from: past it, the midnight deleter may have removed hours of the oldest
// day, and a read that was delayed across the deadline would yield a partial
// sum. Writing that sum would lower an already published, correct daily
// total (the marker may already cover the day). The reads themselves are
// harmless; the write is what must not happen. A missing deadline (no
// timeout configured) never blocks.
func checkDeadline(ctx context.Context) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil
	}
	if now := time.Now(); now.After(deadline) {
		return fmt.Errorf("batch: backfill deadline %s passed before writing a daily key; "+
			"hourly data is no longer guaranteed complete, write skipped", deadline)
	}
	return nil
}

// getDailyEventCount reads a daily ec key; a missing key reads as 0.
func (b *evaluationDailyCountBackfiller) getDailyEventCount(dailyKey string) (int64, error) {
	stored, err := b.cache.Get(dailyKey)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("get daily event count: %w", err)
	}
	s, ok := toString(stored)
	if !ok {
		return 0, fmt.Errorf("batch: daily event count has an unexpected type: %v", stored)
	}
	got, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse daily event count: %w", err)
	}
	return got, nil
}

func (b *evaluationDailyCountBackfiller) sumHourlyEventCounts(hourlyKeys []string) (int64, error) {
	values, err := b.cache.GetMulti(hourlyKeys, true)
	if err != nil {
		return 0, fmt.Errorf("get hourly event counts: %w", err)
	}
	var sum int64
	for _, v := range values {
		s, ok := toString(v)
		if !ok || s == "" {
			continue
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse hourly event count: %w", err)
		}
		sum += n
	}
	return sum, nil
}

// backfillUserCountDay merges the day's hourly HLLs into the daily uc key and
// verifies that the daily union covers the hourly union. PFMERGE is an
// idempotent union: re-running, and racing with live writes to the daily key,
// can only grow it, never lose users.
func (b *evaluationDailyCountBackfiller) backfillUserCountDay(
	ctx context.Context,
	dailyKey string,
	hourlyKeys []string,
) error {
	if err := checkDeadline(ctx); err != nil {
		return err
	}
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
