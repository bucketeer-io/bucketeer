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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluationCountAggregator_AddEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		events             []struct{ ecKey, ucKey, userID string }
		expectedEventCount map[string]int64
		expectedUserCount  map[string]int // number of unique users per key
	}{
		{
			name: "single event",
			events: []struct{ ecKey, ucKey, userID string }{
				{"ec:key1", "uc:key1", "user1"},
			},
			expectedEventCount: map[string]int64{
				"ec:key1": 1,
			},
			expectedUserCount: map[string]int{
				"uc:key1": 1,
			},
		},
		{
			name: "multiple events same key",
			events: []struct{ ecKey, ucKey, userID string }{
				{"ec:key1", "uc:key1", "user1"},
				{"ec:key1", "uc:key1", "user2"},
				{"ec:key1", "uc:key1", "user3"},
			},
			expectedEventCount: map[string]int64{
				"ec:key1": 3,
			},
			expectedUserCount: map[string]int{
				"uc:key1": 3,
			},
		},
		{
			name: "duplicate user IDs deduplicated",
			events: []struct{ ecKey, ucKey, userID string }{
				{"ec:key1", "uc:key1", "user1"},
				{"ec:key1", "uc:key1", "user1"}, // duplicate
				{"ec:key1", "uc:key1", "user2"},
				{"ec:key1", "uc:key1", "user1"}, // duplicate
			},
			expectedEventCount: map[string]int64{
				"ec:key1": 4, // event count still 4
			},
			expectedUserCount: map[string]int{
				"uc:key1": 2, // only 2 unique users
			},
		},
		{
			name: "multiple keys",
			events: []struct{ ecKey, ucKey, userID string }{
				{"ec:key1", "uc:key1", "user1"},
				{"ec:key1", "uc:key1", "user2"},
				{"ec:key2", "uc:key2", "user3"},
				{"ec:key2", "uc:key2", "user4"},
				{"ec:key3", "uc:key3", "user5"},
			},
			expectedEventCount: map[string]int64{
				"ec:key1": 2,
				"ec:key2": 2,
				"ec:key3": 1,
			},
			expectedUserCount: map[string]int{
				"uc:key1": 2,
				"uc:key2": 2,
				"uc:key3": 1,
			},
		},
		{
			name: "same user different keys",
			events: []struct{ ecKey, ucKey, userID string }{
				{"ec:key1", "uc:key1", "user1"},
				{"ec:key2", "uc:key2", "user1"}, // same user, different key
			},
			expectedEventCount: map[string]int64{
				"ec:key1": 1,
				"ec:key2": 1,
			},
			expectedUserCount: map[string]int{
				"uc:key1": 1,
				"uc:key2": 1, // user1 counted in both keys
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			agg := newEvaluationCountAggregator()

			// Add all events
			for _, e := range tt.events {
				agg.addEvent(e.ecKey, e.ucKey, e.userID)
			}

			// Verify event counts
			eventCounts, userCounts := agg.flush()

			assert.Equal(t, tt.expectedEventCount, eventCounts, "event counts mismatch")

			// Verify user counts (check number of unique users per key)
			assert.Len(t, userCounts, len(tt.expectedUserCount), "user count keys mismatch")
			for key, expectedCount := range tt.expectedUserCount {
				assert.Len(t, userCounts[key], expectedCount, "unique user count mismatch for key %s", key)
			}
		})
	}
}

func TestEvaluationCountAggregator_Flush(t *testing.T) {
	t.Parallel()

	t.Run("flush returns accumulated data", func(t *testing.T) {
		t.Parallel()
		agg := newEvaluationCountAggregator()

		agg.addEvent("ec:key1", "uc:key1", "user1")
		agg.addEvent("ec:key1", "uc:key1", "user2")
		agg.addEvent("ec:key2", "uc:key2", "user3")

		eventCounts, userCounts := agg.flush()

		assert.Equal(t, int64(2), eventCounts["ec:key1"])
		assert.Equal(t, int64(1), eventCounts["ec:key2"])
		assert.Len(t, userCounts["uc:key1"], 2)
		assert.Len(t, userCounts["uc:key2"], 1)
	})

	t.Run("flush resets aggregator", func(t *testing.T) {
		t.Parallel()
		agg := newEvaluationCountAggregator()

		agg.addEvent("ec:key1", "uc:key1", "user1")
		agg.flush()

		// After flush, aggregator should be empty
		eventCounts, userCounts := agg.flush()
		assert.Empty(t, eventCounts)
		assert.Empty(t, userCounts)
	})

	t.Run("multiple flushes independent", func(t *testing.T) {
		t.Parallel()
		agg := newEvaluationCountAggregator()

		// First batch
		agg.addEvent("ec:key1", "uc:key1", "user1")
		eventCounts1, _ := agg.flush()
		assert.Equal(t, int64(1), eventCounts1["ec:key1"])

		// Second batch (should not include first batch)
		agg.addEvent("ec:key1", "uc:key1", "user2")
		agg.addEvent("ec:key1", "uc:key1", "user3")
		eventCounts2, _ := agg.flush()
		assert.Equal(t, int64(2), eventCounts2["ec:key1"])
	})
}

func TestEvaluationCountAggregator_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	t.Run("concurrent adds are thread-safe", func(t *testing.T) {
		t.Parallel()
		agg := newEvaluationCountAggregator()

		const numGoroutines = 100
		const eventsPerGoroutine = 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// Concurrently add events from multiple goroutines
		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < eventsPerGoroutine; j++ {
					agg.addEvent("ec:key1", "uc:key1", "user1")
				}
			}(i)
		}

		wg.Wait()

		eventCounts, userCounts := agg.flush()

		// Should have exactly numGoroutines * eventsPerGoroutine events
		assert.Equal(t, int64(numGoroutines*eventsPerGoroutine), eventCounts["ec:key1"])
		// But only 1 unique user
		assert.Len(t, userCounts["uc:key1"], 1)
	})

	t.Run("concurrent add and flush", func(t *testing.T) {
		t.Parallel()
		agg := newEvaluationCountAggregator()

		var wg sync.WaitGroup
		wg.Add(2)

		// Goroutine 1: Add events
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				agg.addEvent("ec:key1", "uc:key1", "user1")
			}
		}()

		// Goroutine 2: Flush periodically
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				agg.flush()
			}
		}()

		wg.Wait()

		// Should not panic (thread-safety test)
		// Final flush to verify state is consistent
		_, _ = agg.flush()
	})
}

func TestEvaluationCountAggregator_RealWorldScenario(t *testing.T) {
	t.Parallel()

	t.Run("simulate 125k events aggregating to few keys", func(t *testing.T) {
		t.Parallel()
		agg := newEvaluationCountAggregator()

		// Simulate realistic scenario:
		// - 1000 events in a batch
		// - 3 popular feature flags (70% of traffic)
		// - 10 other feature flags (30% of traffic)

		popularFeatures := []struct{ feature, variation string }{
			{"feature_login", "variant_A"},
			{"feature_checkout", "variant_B"},
			{"feature_sidebar", "variant_A"},
		}

		// 70% to popular keys
		for i := 0; i < 700; i++ {
			featIdx := i % len(popularFeatures)
			feature := popularFeatures[featIdx].feature
			variation := popularFeatures[featIdx].variation
			userID := string(rune('a' + (i % 26))) // 26 unique users cycling

			agg.addEvent(
				"ec:hour1:"+feature+":"+variation+":env_prod",
				"uc:hour1:"+feature+":"+variation+":env_prod",
				userID,
			)
		}

		// 30% to other keys (10 different features)
		for i := 0; i < 300; i++ {
			featureNum := i % 10 // 10 different features
			feature := "feature_other_" + string(rune('A'+featureNum))
			userID := string(rune('a' + (i % 26)))

			agg.addEvent(
				"ec:hour1:"+feature+":variant_A:env_prod",
				"uc:hour1:"+feature+":variant_A:env_prod",
				userID,
			)
		}

		eventCounts, userCounts := agg.flush()

		// Should have aggregated 1000 events into 13 unique keys (3 popular + 10 other)
		assert.Equal(t, 13, len(eventCounts), "should aggregate to 13 event count keys")
		assert.Equal(t, 13, len(userCounts), "should aggregate to 13 user count keys")

		// Popular keys should have high counts
		for _, pf := range popularFeatures {
			key := "ec:hour1:" + pf.feature + ":" + pf.variation + ":env_prod"
			count := eventCounts[key]
			assert.Greater(t, count, int64(200), "popular key %s should have >200 events", key)
		}

		// Total events should sum to 1000
		totalEvents := int64(0)
		for _, count := range eventCounts {
			totalEvents += count
		}
		assert.Equal(t, int64(1000), totalEvents, "total events should be 1000")
	})
}
