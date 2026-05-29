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

func TestNewLimiter(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc                   string
		config                 Config
		expectedRequestsPerMin int
		expectedBurstSize      int
	}{
		{
			desc:                   "default values",
			config:                 Config{},
			expectedRequestsPerMin: 20,
			expectedBurstSize:      5,
		},
		{
			desc:                   "custom values",
			config:                 Config{RequestsPerMinute: 60, BurstSize: 10},
			expectedRequestsPerMin: 60,
			expectedBurstSize:      10,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			l := NewLimiter(t.Context(), p.config)
			assert.NotNil(t, l)
			assert.Equal(t, p.expectedRequestsPerMin, l.config.RequestsPerMinute)
			assert.Equal(t, p.expectedBurstSize, l.config.BurstSize)
		})
	}
}

func TestAllow(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc      string
		burstSize int
		setup     func(l *Limiter)
		key       string
		expected  bool
	}{
		{
			desc:      "within burst",
			burstSize: 5,
			setup:     func(l *Limiter) {},
			key:       "user1",
			expected:  true,
		},
		{
			desc:      "exceeds burst",
			burstSize: 3,
			setup: func(l *Limiter) {
				for i := 0; i < 3; i++ {
					l.Allow("user1")
				}
			},
			key:      "user1",
			expected: false,
		},
		{
			desc:      "different keys are independent",
			burstSize: 2,
			setup: func(l *Limiter) {
				// Exhaust user1
				l.Allow("user1")
				l.Allow("user1")
			},
			key:      "user2",
			expected: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			l := NewLimiter(t.Context(), Config{RequestsPerMinute: 60, BurstSize: p.burstSize})
			p.setup(l)
			assert.Equal(t, p.expected, l.Allow(p.key))
		})
	}
}

func TestAllow_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	l := NewLimiter(t.Context(), Config{RequestsPerMinute: 600, BurstSize: 50})

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

func TestCleanup(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		setup    func(l *Limiter)
		expected bool
	}{
		{
			desc: "preserves active entries",
			setup: func(l *Limiter) {
				l.Allow("user1")
				l.Cleanup()
			},
			expected: false, // burst=1, already used, still rate limited after cleanup
		},
		{
			desc: "removes idle entries",
			setup: func(l *Limiter) {
				l.Allow("user1")
				// Simulate idle entry by setting lastSeen to the past
				l.mu.Lock()
				l.limiters["user1"].lastSeen = time.Now().Add(-15 * time.Minute)
				l.mu.Unlock()
				l.Cleanup()
			},
			expected: true, // entry was removed, burst is available again
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			l := NewLimiter(t.Context(), Config{RequestsPerMinute: 60, BurstSize: 1})
			p.setup(l)
			assert.Equal(t, p.expected, l.Allow("user1"))
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	cfg := DefaultConfig()
	assert.Equal(t, 20, cfg.RequestsPerMinute)
	assert.Equal(t, 5, cfg.BurstSize)
}
