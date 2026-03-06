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

package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLimiter_DefaultValues(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{})
	assert.NotNil(t, l)
	assert.Equal(t, 20, l.config.RequestsPerMinute)
	assert.Equal(t, 5, l.config.BurstSize)
}

func TestNewLimiter_CustomValues(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{RequestsPerMinute: 60, BurstSize: 10})
	assert.Equal(t, 60, l.config.RequestsPerMinute)
	assert.Equal(t, 10, l.config.BurstSize)
}

func TestAllow_WithinBurst(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{RequestsPerMinute: 60, BurstSize: 5})

	// First 5 requests should be allowed (burst)
	for i := 0; i < 5; i++ {
		assert.True(t, l.Allow("user1"), "request %d should be allowed within burst", i)
	}
}

func TestAllow_ExceedsBurst(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{RequestsPerMinute: 60, BurstSize: 3})

	// First 3 requests should be allowed (burst)
	for i := 0; i < 3; i++ {
		assert.True(t, l.Allow("user1"))
	}

	// Next request should be rate limited
	assert.False(t, l.Allow("user1"))
}

func TestAllow_DifferentKeys(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{RequestsPerMinute: 60, BurstSize: 2})

	// Exhaust user1 burst
	assert.True(t, l.Allow("user1"))
	assert.True(t, l.Allow("user1"))
	assert.False(t, l.Allow("user1"))

	// user2 should still have full burst
	assert.True(t, l.Allow("user2"))
	assert.True(t, l.Allow("user2"))
	assert.False(t, l.Allow("user2"))
}

func TestAllow_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{RequestsPerMinute: 600, BurstSize: 50})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			key := fmt.Sprintf("user-%d", userID)
			for j := 0; j < 10; j++ {
				l.Allow(key) // Should not panic
			}
		}(i)
	}
	wg.Wait()
}

func TestCleanup_PreservesActiveEntries(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{RequestsPerMinute: 60, BurstSize: 1})

	// Exhaust limiter
	assert.True(t, l.Allow("user1"))
	assert.False(t, l.Allow("user1"))

	// Cleanup should preserve active entries (lastSeen is recent)
	l.Cleanup()

	// After cleanup, user1 is still rate limited (entry was preserved)
	assert.False(t, l.Allow("user1"))
}

func TestCleanup_RemovesIdleEntries(t *testing.T) {
	t.Parallel()
	l := NewLimiter(Config{RequestsPerMinute: 60, BurstSize: 1})

	// Exhaust limiter
	assert.True(t, l.Allow("user1"))
	assert.False(t, l.Allow("user1"))

	// Simulate idle entry by setting lastSeen to the past
	l.mu.Lock()
	l.limiters["user1"].lastSeen = time.Now().Add(-15 * time.Minute)
	l.mu.Unlock()

	l.Cleanup()

	// After cleanup, user1 entry was removed so burst is available again
	assert.True(t, l.Allow("user1"))
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	cfg := DefaultConfig()
	assert.Equal(t, 20, cfg.RequestsPerMinute)
	assert.Equal(t, 5, cfg.BurstSize)
}
