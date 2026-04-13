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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
)

func TestComputeBackoffAndTTL_NoCap(t *testing.T) {
	t.Parallel()

	logger, err := log.NewLogger()
	require.NoError(t, err)

	w := &goalEvtWriter{
		logger:                  logger,
		retryGoalEventInterval:  1 * time.Second,
		maxRetryGoalEventPeriod: 24 * time.Hour,
		maxRetryBackoffInterval: 0, // No cap
	}

	tests := []struct {
		name                string
		retryCount          int
		expectedMinInterval time.Duration
		expectedMaxInterval time.Duration
		maxRetryBackoffCap  time.Duration
	}{
		{
			name:                "first retry (count=0)",
			retryCount:          0,
			expectedMinInterval: 500 * time.Millisecond, // backoff library adds randomization (0.5x-1.5x)
			expectedMaxInterval: 2 * time.Second,
			maxRetryBackoffCap:  0, // no cap
		},
		{
			name:                "second retry (count=1)",
			retryCount:          1,
			expectedMinInterval: 1 * time.Second, // 0.5 * 2s
			expectedMaxInterval: 4 * time.Second, // 1.5 * 2s (with some jitter)
		},
		{
			name:                "third retry (count=2)",
			retryCount:          2,
			expectedMinInterval: 2 * time.Second,
			expectedMaxInterval: 8 * time.Second,
		},
		{
			name:                "high retry count (count=10)",
			retryCount:          10,
			expectedMinInterval: 256 * time.Second, // 2^8 * 0.5 (randomization)
			expectedMaxInterval: 2048 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextInterval, ttl, err := w.computeBackoffAndTTL(
				tt.retryCount,
				0, // firstRetryAt = 0 (first time)
				w.retryGoalEventInterval,
				w.maxRetryGoalEventPeriod,
				tt.maxRetryBackoffCap,
			)

			require.NoError(t, err)
			assert.GreaterOrEqual(t, nextInterval, tt.expectedMinInterval, "interval should be at least expectedMin")
			assert.LessOrEqual(t, nextInterval, tt.expectedMaxInterval, "interval should be at most expectedMax")
			assert.Equal(t, w.maxRetryGoalEventPeriod, ttl, "TTL should equal maxRetryPeriod on first retry")
		})
	}
}

func TestComputeBackoffAndTTL_WithCap(t *testing.T) {
	t.Parallel()

	logger, err := log.NewLogger()
	require.NoError(t, err)

	w := &goalEvtWriter{
		logger:                  logger,
		retryGoalEventInterval:  1 * time.Second,
		maxRetryGoalEventPeriod: 12 * time.Hour,
		maxRetryBackoffInterval: 5 * time.Second, // Cap at 5 seconds
	}

	tests := []struct {
		name                string
		retryCount          int
		expectedMinInterval time.Duration
		expectedMaxInterval time.Duration
	}{
		{
			name:                "first retry (count=0) - below cap",
			retryCount:          0,
			expectedMinInterval: 500 * time.Millisecond,
			expectedMaxInterval: 2 * time.Second,
		},
		{
			name:                "second retry (count=1) - below cap",
			retryCount:          1,
			expectedMinInterval: 1 * time.Second,
			expectedMaxInterval: 4 * time.Second,
		},
		{
			name:                "third retry (count=2) - approaching cap",
			retryCount:          2,
			expectedMinInterval: 2 * time.Second,
			expectedMaxInterval: 8 * time.Second, // Could be up to 8s with randomization before cap applies
		},
		{
			name:                "fourth retry (count=3) - capped",
			retryCount:          3,
			expectedMinInterval: 2 * time.Second,  // 0.5 * 5s (randomized min)
			expectedMaxInterval: 10 * time.Second, // 1.5 * 5s (randomized max), but note: library will clamp
		},
		{
			name:                "tenth retry (count=10) - capped",
			retryCount:          10,
			expectedMinInterval: 2 * time.Second,
			expectedMaxInterval: 10 * time.Second, // Should be capped at ~5s with randomization
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextInterval, ttl, err := w.computeBackoffAndTTL(
				tt.retryCount,
				0, // firstRetryAt = 0 (first time)
				w.retryGoalEventInterval,
				w.maxRetryGoalEventPeriod,
				w.maxRetryBackoffInterval,
			)

			require.NoError(t, err)
			assert.GreaterOrEqual(t, nextInterval, tt.expectedMinInterval, "interval should be at least expectedMin")
			assert.LessOrEqual(t, nextInterval, tt.expectedMaxInterval, "interval should be at most expectedMax")
			assert.Equal(t, w.maxRetryGoalEventPeriod, ttl, "TTL should equal maxRetryPeriod on first retry")
		})
	}
}

func TestComputeBackoffAndTTL_WithCapE2EScenario(t *testing.T) {
	t.Parallel()

	logger, err := log.NewLogger()
	require.NoError(t, err)

	// E2E test configuration
	w := &goalEvtWriter{
		logger:                  logger,
		retryGoalEventInterval:  1 * time.Second,
		maxRetryGoalEventPeriod: 12 * time.Hour,
		maxRetryBackoffInterval: 5 * time.Second, // E2E test cap
	}

	// Simulate 50 retries like in the E2E test
	totalTime := time.Duration(0)
	for i := 0; i < 50; i++ {
		nextInterval, _, err := w.computeBackoffAndTTL(
			i,
			0,
			w.retryGoalEventInterval,
			w.maxRetryGoalEventPeriod,
			w.maxRetryBackoffInterval,
		)
		require.NoError(t, err)

		// With randomization factor (up to 1.5x), allow up to 7.5 seconds
		// but this is still much better than exponential growth without cap
		if i >= 3 {
			assert.LessOrEqual(t, nextInterval, 10*time.Second, "should be roughly capped around 5s (with randomization) for retry %d", i)
		}

		totalTime += nextInterval
	}

	// With 5 second cap (and randomization), 50 retries should complete well within 10 minutes
	// Even with max randomization (1.5x * 5s = 7.5s), 50 * 7.5s = 375s = 6.25 minutes
	maxExpectedTime := 8 * time.Minute // 480 seconds (generous buffer for randomization)
	assert.Less(t, totalTime, maxExpectedTime, "50 retries should complete within 8 minutes with 5s cap")

	t.Logf("Total time for 50 retries with 5s cap: %v (should be < 8 minutes)", totalTime)
}

func TestComputeBackoffAndTTL_ExceedsMaxPeriod(t *testing.T) {
	t.Parallel()

	logger, err := log.NewLogger()
	require.NoError(t, err)

	w := &goalEvtWriter{
		logger:                  logger,
		retryGoalEventInterval:  1 * time.Second,
		maxRetryGoalEventPeriod: 10 * time.Second, // Very short max period
		maxRetryBackoffInterval: 0,                // No cap
	}

	// First retry at 11 seconds ago (past the max period)
	now := time.Now().Add(-11 * time.Second).Unix()

	// Should fail because we're already past the max retry period
	_, _, err = w.computeBackoffAndTTL(
		5,
		now,
		w.retryGoalEventInterval,
		w.maxRetryGoalEventPeriod,
		w.maxRetryBackoffInterval,
	)

	// Should return an error because max retry period exceeded
	require.Error(t, err)
	assert.Contains(t, err.Error(), "retry period exceeded")
}

func TestComputeBackoffAndTTL_TTLDecreases(t *testing.T) {
	t.Parallel()

	logger, err := log.NewLogger()
	require.NoError(t, err)

	w := &goalEvtWriter{
		logger:                  logger,
		retryGoalEventInterval:  1 * time.Second,
		maxRetryGoalEventPeriod: 1 * time.Hour,
		maxRetryBackoffInterval: 5 * time.Second,
	}

	// First retry
	_, ttl1, err := w.computeBackoffAndTTL(
		0,
		0, // No previous retry
		w.retryGoalEventInterval,
		w.maxRetryGoalEventPeriod,
		w.maxRetryBackoffInterval,
	)
	require.NoError(t, err)
	assert.Equal(t, 1*time.Hour, ttl1, "First retry TTL should be full period")

	// Subsequent retry after 30 minutes
	firstRetryAt := time.Now().Add(-30 * time.Minute).Unix()
	_, ttl2, err := w.computeBackoffAndTTL(
		5,
		firstRetryAt,
		w.retryGoalEventInterval,
		w.maxRetryGoalEventPeriod,
		w.maxRetryBackoffInterval,
	)
	require.NoError(t, err)
	assert.Less(t, ttl2, ttl1, "TTL should decrease on subsequent retries")
	assert.Greater(t, ttl2, 20*time.Minute, "TTL should be remaining time (~30 minutes)")
	assert.Less(t, ttl2, 35*time.Minute, "TTL should be remaining time (~30 minutes)")
}

func TestComputeBackoffAndTTL_ZeroVsPositiveCap(t *testing.T) {
	t.Parallel()

	logger, err := log.NewLogger()
	require.NoError(t, err)

	w := &goalEvtWriter{
		logger:                  logger,
		retryGoalEventInterval:  1 * time.Second,
		maxRetryGoalEventPeriod: 24 * time.Hour,
	}

	// Test with no cap (0)
	intervalNoCap, _, err := w.computeBackoffAndTTL(
		10, // High retry count
		0,
		w.retryGoalEventInterval,
		w.maxRetryGoalEventPeriod,
		0, // No cap
	)
	require.NoError(t, err)

	// Test with cap
	intervalWithCap, _, err := w.computeBackoffAndTTL(
		10, // Same retry count
		0,
		w.retryGoalEventInterval,
		w.maxRetryGoalEventPeriod,
		5*time.Second, // 5 second cap
	)
	require.NoError(t, err)

	// Without cap should be much larger (exponential)
	assert.Greater(t, intervalNoCap, intervalWithCap, "uncapped interval should be larger than capped")
	assert.LessOrEqual(t, intervalWithCap, 10*time.Second, "capped interval should not exceed ~5 seconds (with randomization)")
	assert.Greater(t, intervalNoCap, 100*time.Second, "uncapped interval should be exponential (>100s for retry 10)")

	t.Logf("Uncapped interval (retry 10): %v", intervalNoCap)
	t.Logf("Capped interval (retry 10): %v", intervalWithCap)
}
