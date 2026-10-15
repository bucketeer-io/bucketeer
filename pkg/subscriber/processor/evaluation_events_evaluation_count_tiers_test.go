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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestNewEvaluationCountKeyPairs(t *testing.T) {
	t.Parallel()

	// 2024-03-09 09:23:45 UTC
	ts := time.Date(2024, 3, 9, 9, 23, 45, 0, time.UTC).Unix()
	minuteTs := time.Date(2024, 3, 9, 9, 23, 0, 0, time.UTC).Unix()
	hourTs := time.Date(2024, 3, 9, 9, 0, 0, 0, time.UTC).Unix()
	dayTs := time.Date(2024, 3, 9, 0, 0, 0, 0, time.UTC).Unix()

	p := &evaluationCountEventPersister{logger: zap.NewNop()}
	pairs := p.newEvaluationCountKeyPairs("feature-1", "variation-A", "env-1", ts)
	require.Len(t, pairs, 3)

	assert.Equal(t, fmt.Sprintf("env-1:ec1m:%d:feature-1:variation-A", minuteTs), pairs[0].ecKey)
	assert.Equal(t, fmt.Sprintf("env-1:uc1m:%d:feature-1:variation-A", minuteTs), pairs[0].ucKey)
	assert.Equal(t, cachev3.MinuteCountTTL, pairs[0].ttl)

	assert.Equal(t, fmt.Sprintf("env-1:ec:%d:feature-1:variation-A", hourTs), pairs[1].ecKey)
	assert.Equal(t, fmt.Sprintf("env-1:uc:%d:feature-1:variation-A", hourTs), pairs[1].ucKey)
	assert.Equal(t, time.Duration(0), pairs[1].ttl, "hourly keys are owned by the deleter job, no TTL")

	assert.Equal(t, fmt.Sprintf("env-1:ec1d:%d:feature-1:variation-A", dayTs), pairs[2].ecKey)
	assert.Equal(t, fmt.Sprintf("env-1:uc1d:%d:feature-1:variation-A", dayTs), pairs[2].ucKey)
	assert.Equal(t, cachev3.DailyCountTTL, pairs[2].ttl)

	// The hourly key must stay byte-identical to the legacy key builder.
	assert.Equal(t,
		legacyHourlyCountKey(eventCountKey, "feature-1", "variation-A", "env-1", ts),
		pairs[1].ecKey,
	)
	assert.Equal(t,
		legacyHourlyCountKey(userCountKey, "feature-1", "variation-A", "env-1", ts),
		pairs[1].ucKey,
	)
}

// legacyHourlyCountKey reproduces the pre-redesign hourly key construction
// (newEvaluationCountkeyV2, now removed) independently of the production key
// builder, so tests keep proving the hourly tier stays byte-identical to the
// keys already in Redis.
func legacyHourlyCountKey(kind, featureID, variationID, environmentId string, timestamp int64) string {
	hourTs := time.Unix(timestamp, 0).UTC().Truncate(time.Hour).Unix()
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%d:%s:%s", hourTs, featureID, variationID),
		environmentId,
	)
}

// TestEvaluationCountKeysIgnoreLocalTimezone must not call t.Parallel(): it
// temporarily overrides time.Local to prove the bucket truncation no longer
// depends on the pod timezone (the pre-existing time.Unix wall-clock bug).
func TestEvaluationCountKeysIgnoreLocalTimezone(t *testing.T) {
	origLocal := time.Local
	// A fractional-offset zone (like Asia/Kolkata) breaks wall-clock-based
	// truncation for hours and minutes, not only for days.
	time.Local = time.FixedZone("UTC+5:30", 5*3600+30*60)
	defer func() { time.Local = origLocal }()

	ts := time.Date(2024, 3, 9, 9, 23, 45, 0, time.UTC).Unix()
	assert.Equal(t, time.Date(2024, 3, 9, 9, 23, 0, 0, time.UTC).Unix(), truncateUnixUTC(ts, time.Minute))
	assert.Equal(t, time.Date(2024, 3, 9, 9, 0, 0, 0, time.UTC).Unix(), truncateUnixUTC(ts, time.Hour))
	assert.Equal(t, time.Date(2024, 3, 9, 0, 0, 0, 0, time.UTC).Unix(), truncateUnixUTC(ts, 24*time.Hour))
}

func TestSplitCountKey(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc            string
		key             string
		expectedKind    string
		expectedPairing string
	}{
		{
			desc:            "hourly ec key",
			key:             "env-1:ec:1709974800:feature-1:variation-A",
			expectedKind:    "ec",
			expectedPairing: "env-1:1h:1709974800:feature-1:variation-A",
		},
		{
			desc:            "hourly uc key pairs with its ec key",
			key:             "env-1:uc:1709974800:feature-1:variation-A",
			expectedKind:    "uc",
			expectedPairing: "env-1:1h:1709974800:feature-1:variation-A",
		},
		{
			desc:            "minute ec key",
			key:             "env-1:ec1m:1709974800:feature-1:variation-A",
			expectedKind:    "ec1m",
			expectedPairing: "env-1:1m:1709974800:feature-1:variation-A",
		},
		{
			desc:            "daily uc key",
			key:             "env-1:uc1d:1709942400:feature-1:variation-A",
			expectedKind:    "uc1d",
			expectedPairing: "env-1:1d:1709942400:feature-1:variation-A",
		},
		{
			desc:            "admin environment key has no environment prefix",
			key:             "ec:1709974800:feature-1:variation-A",
			expectedKind:    "ec",
			expectedPairing: "1h:1709974800:feature-1:variation-A",
		},
		{
			desc:            "unknown kind falls back to the full key",
			key:             "something:else",
			expectedKind:    "",
			expectedPairing: "something:else",
		},
		// A feature ID that equals a kind name must not be mistaken for the
		// kind: the kind is positional and always followed by the bucket
		// timestamp. Before the positional parser, "ec1m" was matched inside
		// the feature segment and the daily key got the minute TTL.
		{
			desc:            "feature ID equal to a kind name inside a daily key",
			key:             "env-1:ec1d:1709942400:ec1m:variation-A",
			expectedKind:    "ec1d",
			expectedPairing: "env-1:1d:1709942400:ec1m:variation-A",
		},
		{
			desc:            "feature ID equal to a kind name inside an hourly key",
			key:             "env-1:uc:1709974800:ec1d:variation-A",
			expectedKind:    "uc",
			expectedPairing: "env-1:1h:1709974800:ec1d:variation-A",
		},
		{
			desc:            "environment ID equal to a kind name",
			key:             "ec:uc1d:1709942400:feature-1:variation-A",
			expectedKind:    "uc1d",
			expectedPairing: "ec:1d:1709942400:feature-1:variation-A",
		},
		{
			desc:            "admin key whose feature ID equals a kind name",
			key:             "ec1m:1709974800:ec:variation-A",
			expectedKind:    "ec1m",
			expectedPairing: "1m:1709974800:ec:variation-A",
		},
		{
			desc:            "kind-like segment not followed by a bucket timestamp is not a kind",
			key:             "ec:feature:variation",
			expectedKind:    "",
			expectedPairing: "ec:feature:variation",
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			kind, pairing := splitCountKey(p.key)
			assert.Equal(t, p.expectedKind, kind)
			assert.Equal(t, p.expectedPairing, pairing)
		})
	}
}

// At a UTC midnight the minute, hour, and day buckets share the same
// timestamp; the pairing keys must still keep the tiers apart.
func TestSplitCountKeyKeepsTiersApartAtMidnight(t *testing.T) {
	t.Parallel()

	midnight := time.Date(2024, 3, 9, 0, 0, 0, 0, time.UTC).Unix()
	pairings := map[string]struct{}{}
	for _, kind := range []string{"ec1m", "ec", "ec1d"} {
		key := fmt.Sprintf("env-1:%s:%d:feature-1:variation-A", kind, midnight)
		_, pairing := splitCountKey(key)
		pairings[pairing] = struct{}{}
	}
	assert.Len(t, pairings, 3, "each tier must produce a distinct pairing key")
}

func TestFlushAggregatedCountsAppliesTierTTLs(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 3, 9, 9, 23, 45, 0, time.UTC).Unix()
	mockCache := &mockEvaluationCountCache{}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	envEvents := environmentEventMap{
		"env-1": eventMap{
			"event-1": &eventproto.EvaluationEvent{
				FeatureId:      "feature-1",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-1",
				User:           &userproto.User{Id: "user-1"},
				Timestamp:      ts,
				FeatureVersion: 1,
			},
		},
	}
	fails := persister.incrementEnvEvents(envEvents)
	require.Empty(t, fails)

	pairs := persister.newEvaluationCountKeyPairs("feature-1", "variation-A", "env-1", ts)
	minute, hour, day := pairs[0], pairs[1], pairs[2]

	assert.Equal(t, cachev3.MinuteCountTTL, mockCache.expireCalls[minute.ecKey])
	assert.Equal(t, cachev3.MinuteCountTTL, mockCache.expireCalls[minute.ucKey])
	assert.Equal(t, cachev3.DailyCountTTL, mockCache.expireCalls[day.ecKey])
	assert.Equal(t, cachev3.DailyCountTTL, mockCache.expireCalls[day.ucKey])
	assert.NotContains(t, mockCache.expireCalls, hour.ecKey, "hourly keys must never get a TTL")
	assert.NotContains(t, mockCache.expireCalls, hour.ucKey, "hourly keys must never get a TTL")
}

// A feature whose ID equals a kind name must get the TTL of the tier it is
// written to, not the TTL of the kind its name happens to match.
func TestFlushAggregatedCountsTierTTLsWithKindLikeFeatureID(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 3, 9, 9, 23, 45, 0, time.UTC).Unix()
	mockCache := &mockEvaluationCountCache{}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	envEvents := environmentEventMap{
		"env-1": eventMap{
			"event-1": &eventproto.EvaluationEvent{
				FeatureId:      "ec1m",
				VariationId:    "uc",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-1",
				User:           &userproto.User{Id: "user-1"},
				Timestamp:      ts,
				FeatureVersion: 1,
			},
		},
	}
	fails := persister.incrementEnvEvents(envEvents)
	require.Empty(t, fails)

	pairs := persister.newEvaluationCountKeyPairs("ec1m", "uc", "env-1", ts)
	minute, hour, day := pairs[0], pairs[1], pairs[2]

	assert.Equal(t, cachev3.MinuteCountTTL, mockCache.expireCalls[minute.ecKey])
	assert.Equal(t, cachev3.MinuteCountTTL, mockCache.expireCalls[minute.ucKey])
	assert.Equal(t, cachev3.DailyCountTTL, mockCache.expireCalls[day.ecKey],
		"daily key must keep the daily TTL even though the feature ID is 'ec1m'")
	assert.Equal(t, cachev3.DailyCountTTL, mockCache.expireCalls[day.ucKey])
	assert.NotContains(t, mockCache.expireCalls, hour.ecKey, "hourly keys must never get a TTL")
	assert.NotContains(t, mockCache.expireCalls, hour.ucKey, "hourly keys must never get a TTL")
	// Every tier was written exactly once and each ec/uc pair stayed paired.
	assert.Equal(t, 3, mockCache.pfaddCallCount)
	assert.Equal(t, 3, mockCache.execCount)
	assert.Equal(t, int64(1), mockCache.lastEventCounts[day.ecKey])
}

// The three tiers of one event must be flushed in two pipelined round trips
// (all PFADDs, then all INCRBYs), not in up to twelve synchronous calls, and
// each EXPIRE must be queued directly behind the mutation of its key.
func TestFlushAggregatedCountsPipelinesTiers(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 3, 9, 9, 23, 45, 0, time.UTC).Unix()
	mockCache := &mockEvaluationCountCache{}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	envEvents := environmentEventMap{
		"env-1": eventMap{
			"event-1": &eventproto.EvaluationEvent{
				FeatureId:      "feature-1",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-1",
				User:           &userproto.User{Id: "user-1"},
				Timestamp:      ts,
				FeatureVersion: 1,
			},
			"event-2": &eventproto.EvaluationEvent{
				FeatureId:      "feature-2",
				VariationId:    "variation-B",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-2",
				User:           &userproto.User{Id: "user-2"},
				Timestamp:      ts,
				FeatureVersion: 1,
			},
		},
	}
	fails := persister.incrementEnvEvents(envEvents)
	require.Empty(t, fails)

	assert.Equal(t, 2, mockCache.pipelineCount, "one pipeline for PFADDs, one for INCRBYs")
	assert.Equal(t, 2, mockCache.execCallCount, "two round trips regardless of the number of tiers")
	assert.Equal(t, 6, mockCache.pfaddCallCount, "3 tiers x 2 events")
	assert.Equal(t, 6, mockCache.execCount, "3 tiers x 2 events")

	// All PFADDs precede all INCRBYs.
	lastPFAdd, firstIncrBy := -1, len(mockCache.commandOrder)
	for idx, cmd := range mockCache.commandOrder {
		switch {
		case strings.HasPrefix(cmd, "PFADD:"):
			lastPFAdd = idx
		case strings.HasPrefix(cmd, "INCRBY:") && idx < firstIncrBy:
			firstIncrBy = idx
		}
	}
	assert.Less(t, lastPFAdd, firstIncrBy, "every PFADD must run before any INCRBY")

	// Each EXPIRE directly follows the mutation of the same key, so the TTL
	// is applied in the same round trip (and on the same node) as the write.
	expires := 0
	for idx, cmd := range mockCache.commandOrder {
		if !strings.HasPrefix(cmd, "EXPIRE:") {
			continue
		}
		expires++
		require.Greater(t, idx, 0)
		prev := mockCache.commandOrder[idx-1]
		key := strings.TrimPrefix(cmd, "EXPIRE:")
		assert.True(t,
			prev == "PFADD:"+key || prev == "INCRBY:"+key,
			"EXPIRE of %s must directly follow its mutation, got %q", key, prev)
	}
	assert.Equal(t, 8, expires, "minute and day tiers x ec/uc x 2 events; hourly keys get no EXPIRE")
}

// A failed EXPIRE must be logged, not turned into a Nack: the counter write
// already succeeded, so redelivering the events would double count them.
func TestFlushAggregatedCountsExpireFailureDoesNotFailPair(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 3, 9, 9, 23, 45, 0, time.UTC).Unix()
	mockCache := &mockEvaluationCountCache{shouldFailExpire: true}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	envEvents := environmentEventMap{
		"env-1": eventMap{
			"event-1": &eventproto.EvaluationEvent{
				FeatureId:      "feature-1",
				VariationId:    "variation-A",
				Reason:         &featureproto.Reason{Type: featureproto.Reason_TARGET},
				UserId:         "user-1",
				User:           &userproto.User{Id: "user-1"},
				Timestamp:      ts,
				FeatureVersion: 1,
			},
		},
	}
	fails := persister.incrementEnvEvents(envEvents)
	assert.Empty(t, fails, "EXPIRE failures must not Nack events")

	pairs := persister.newEvaluationCountKeyPairs("feature-1", "variation-A", "env-1", ts)
	for _, pair := range pairs {
		assert.Equal(t, int64(1), mockCache.lastEventCounts[pair.ecKey])
		assert.Contains(t, mockCache.lastUserCounts[pair.ucKey], "user-1")
	}
}

func TestWriteDailyCutoverMarker(t *testing.T) {
	t.Parallel()

	mockCache := &mockEvaluationCountCache{}
	persister := &evaluationCountEventPersister{
		evaluationCountCacher: mockCache,
		logger:                zap.NewNop(),
	}

	// The marker must skip the deploy day AND the following day: during a
	// rolling update old hourly-only pods can still be consuming events when
	// the next UTC midnight passes, so that day is not complete in the daily
	// tier either.
	require.Equal(t, 2, dailyCutoverRolloutDays)
	before := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, 2).Unix()
	persister.writeDailyCutoverMarker()
	after := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, 2).Unix()

	assert.Equal(t, cachev3.EvaluationDailyCutoverKey, mockCache.putIfAbsentKey)
	value, err := strconv.ParseInt(mockCache.putIfAbsentValue.(string), 10, 64)
	require.NoError(t, err)
	// Guard against the test straddling a UTC midnight.
	assert.True(t, value == before || value == after,
		"marker must hold the second UTC midnight after startup, got %d", value)
	nextMidnight := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, 1).Unix()
	assert.Greater(t, value, nextMidnight,
		"the next midnight is not safe: old pods may still write hourly-only during the rollout")
}
