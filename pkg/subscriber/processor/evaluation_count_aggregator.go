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
)

// evaluationCountAggregator batches evaluation counts in memory before flushing to Redis
//
// Example: 1000 events in 10 seconds for the same feature/variation/environment:
// - Without aggregation: 1000 INCR + 1000 PFADD = 2000 Redis calls
// - With aggregation: 1 INCRBY 1000 + 1 PFADD (1000 userIDs) = 2 Redis calls
type evaluationCountAggregator struct {
	mu          sync.Mutex
	eventCounts map[string]int64               // eventCountKey -> accumulated count
	userCounts  map[string]map[string]struct{} // userCountKey -> set of unique userIDs
}

func newEvaluationCountAggregator() *evaluationCountAggregator {
	return &evaluationCountAggregator{
		eventCounts: make(map[string]int64),
		userCounts:  make(map[string]map[string]struct{}),
	}
}

// addEvent accumulates one event count and one user ID for the given keys
// Multiple calls with the same keys will aggregate:
// - eventCountKey: increment the counter
// - userCountKey: add userID to the set (automatically deduplicated)
func (a *evaluationCountAggregator) addEvent(eventCountKey, userCountKey, userID string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Accumulate event count (same key → increment counter)
	a.eventCounts[eventCountKey]++

	// Accumulate unique user IDs (same key → add to set, deduped automatically)
	if _, exists := a.userCounts[userCountKey]; !exists {
		a.userCounts[userCountKey] = make(map[string]struct{})
	}
	a.userCounts[userCountKey][userID] = struct{}{}
}

// flush returns all accumulated data and resets the aggregator
// Returns: (eventCounts, userCounts)
func (a *evaluationCountAggregator) flush() (map[string]int64, map[string]map[string]struct{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	eventCounts := a.eventCounts
	userCounts := a.userCounts

	// Reset for next batch
	a.eventCounts = make(map[string]int64)
	a.userCounts = make(map[string]map[string]struct{})

	return eventCounts, userCounts
}
