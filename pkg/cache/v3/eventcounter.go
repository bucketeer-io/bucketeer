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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v3

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	"github.com/bucketeer-io/bucketeer/v2/pkg/errgroup"
)

// getEventCountsConcurrency bounds how many per-day count fetches run in
// parallel so a long time range cannot flood Redis with concurrent requests.
const getEventCountsConcurrency = 8

// Evaluation count keys come in three tiers sharing one format,
// "<environmentId>:<kind>:<bucketTimestamp>:<featureId>:<variationId>"
// (no environment prefix for the admin environment):
//
//   - minute tier (ec1m/uc1m): 1-minute UTC buckets serving the "past hour"
//     view, expired by TTL.
//   - hour tier (ec/uc): 1-hour UTC buckets serving the "past 24 hours" view,
//     deleted by the redis-counter-deleter batch job (no TTL).
//   - day tier (ec1d/uc1d): UTC-midnight buckets serving the multi-day views,
//     expired by TTL.
//
// The minute/day kind names deliberately do NOT extend "ec:"/"uc:" with a
// colon: the deleter scans "<env>:ec:*" / "<env>:uc:*" and aborts its whole
// run on keys whose timestamp segment it cannot parse, so the new kinds must
// stay invisible to those scan prefixes ("<env>:ec1m:" etc. never match).
const (
	EventCountMinuteKind = "ec1m"
	UserCountMinuteKind  = "uc1m"
	EventCountHourKind   = "ec"
	UserCountHourKind    = "uc"
	EventCountDailyKind  = "ec1d"
	UserCountDailyKind   = "uc1d"

	// MinuteCountTTL keeps minute buckets around for twice the 1-hour window
	// they serve, so a "past hour" read never hits an expired bucket.
	MinuteCountTTL = 2 * time.Hour
	// DailyCountTTL keeps daily buckets just beyond the longest supported
	// range (90 days).
	DailyCountTTL = 92 * 24 * time.Hour
)

// EvaluationDailyCutoverKey holds the unix timestamp (seconds) of the first
// UTC midnight that is fully covered by the daily tier: days starting on or
// after it are served from ec1d/uc1d keys, days before it fall back to the
// day's 24 hourly keys (§3.1 Option A of the timeseries key redesign).
//
// The persister writes it once with SETNX on startup, pointing at the UTC
// midnight after the deploy (the deploy day itself is incomplete in the daily
// tier). The backfill job moves it back to the earliest backfilled day after
// a verified run. When the marker is absent, readers serve every day from
// hourly keys, which is exactly the pre-redesign behavior.
const EvaluationDailyCutoverKey = "system:ec1d:cutover"

type EventCounterCache interface {
	GetEventCounts(keys []string) ([]float64, error)
	GetEventCountsV2(keys [][]string) ([]float64, error)
	GetUserCount(key string) (int64, error)
	GetUserCounts(keys []string) ([]float64, error)
	GetUserCountsV2(keys []string) ([]float64, error)
	// CountUniqueUsers returns the union cardinality of the user-count HLLs
	// stored at the given keys. In cluster mode, keys spanning multiple hash
	// slots are counted through short-lived staging keys.
	CountUniqueUsers(keys []string) (int64, error)
	UpdateUserCount(key, userID string) error
}

type eventCounterCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewEventCountCache(c cache.MultiGetDeleteCountCache) EventCounterCache {
	return &eventCounterCache{cache: c}
}

func (c *eventCounterCache) GetEventCounts(keys []string) ([]float64, error) {
	values, err := c.cache.GetMulti(keys, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple keys: %w", err)
	}
	return c.getEventValues(values)
}

func (*eventCounterCache) getEventValues(values []interface{}) ([]float64, error) {
	eventVals := make([]float64, 0, len(values))
	for _, v := range values {
		var str string
		switch v := v.(type) {
		case []byte:
			str = string(v)
		case string:
			str = v
		default:
			return nil, fmt.Errorf("unexpected value type: %v", v)
		}
		if str == "" {
			str = "0"
		}
		float, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float: %w", err)
		}
		eventVals = append(eventVals, float)
	}
	return eventVals, nil
}

func (c *eventCounterCache) GetEventCountsV2(keys [][]string) ([]float64, error) {
	// The days are independent, so fetch them in parallel. The local errgroup
	// recovers panics inside the goroutines, which plain x/sync/errgroup
	// would let crash the process.
	eventVals := make([]float64, len(keys))
	var eg errgroup.Group
	eg.SetLimit(getEventCountsConcurrency)
	for i, day := range keys {
		eg.Go(func() error {
			values, err := c.cache.GetMulti(day, true)
			if err != nil {
				return fmt.Errorf("failed to get multiple keys: %w", err)
			}
			dayTotal, err := c.getEventValues(values)
			if err != nil {
				return err
			}
			var totalVal float64
			for _, v := range dayTotal {
				totalVal += v
			}
			eventVals[i] = totalVal
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return eventVals, nil
}

func (c *eventCounterCache) GetUserCount(key string) (int64, error) {
	return c.cache.PFCount(key)
}

func (c *eventCounterCache) GetUserCounts(keys []string) ([]float64, error) {
	userVals := make([]float64, 0, len(keys))
	for _, key := range keys {
		count, err := c.cache.PFCount(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get PFCount for key %s: %w", key, err)
		}
		userVals = append(userVals, float64(count))
	}
	return userVals, nil
}

func (c *eventCounterCache) GetUserCountsV2(keys []string) ([]float64, error) {
	return c.GetUserCounts(keys)
}

func (c *eventCounterCache) UpdateUserCount(key, userID string) error {
	_, err := c.cache.PFAdd(key, userID)
	if err != nil {
		return fmt.Errorf("failed to update user count: %w", err)
	}
	return nil
}

func (c *eventCounterCache) CountUniqueUsers(keys []string) (int64, error) {
	count, err := c.cache.PFCount(keys...)
	if err != nil {
		return 0, fmt.Errorf("failed to count unique users: %w", err)
	}
	return count, nil
}
