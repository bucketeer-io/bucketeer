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
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Config holds rate limiter configuration.
type Config struct {
	// RequestsPerMinute is the maximum number of requests per minute per key.
	RequestsPerMinute int
	// BurstSize is the maximum number of requests that can be made at once.
	BurstSize int
}

// DefaultConfig returns a default rate limit configuration.
func DefaultConfig() Config {
	return Config{
		RequestsPerMinute: 20,
		BurstSize:         5,
	}
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Limiter provides per-key rate limiting using a token bucket algorithm.
type Limiter struct {
	mu       sync.Mutex
	limiters map[string]*limiterEntry
	config   Config
}

// NewLimiter creates a new rate limiter with the given configuration.
func NewLimiter(cfg Config) *Limiter {
	if cfg.RequestsPerMinute <= 0 {
		cfg.RequestsPerMinute = 20
	}
	if cfg.BurstSize <= 0 {
		cfg.BurstSize = 5
	}
	return &Limiter{
		limiters: make(map[string]*limiterEntry),
		config:   cfg,
	}
}

// Allow checks whether a request from the given key is allowed.
// Returns true if the request is within rate limits, false otherwise.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	entry, exists := l.limiters[key]
	if !exists {
		// rate.Limit is in events per second
		r := rate.Limit(float64(l.config.RequestsPerMinute) / 60.0)
		entry = &limiterEntry{
			limiter:  rate.NewLimiter(r, l.config.BurstSize),
			lastSeen: time.Now(),
		}
		l.limiters[key] = entry
	} else {
		entry.lastSeen = time.Now()
	}
	l.mu.Unlock()
	return entry.limiter.Allow()
}

// Cleanup removes limiters that have been idle for more than 10 minutes.
// This preserves rate limit state for active users while preventing
// unbounded growth from ephemeral keys.
//
// Callers must invoke Cleanup periodically (e.g., via time.Ticker) to prevent
// unbounded memory growth. Each new key that calls Allow creates an entry in
// the internal map; without periodic cleanup, long-running servers with many
// distinct keys will accumulate entries indefinitely.
func (l *Limiter) Cleanup() {
	l.mu.Lock()
	threshold := time.Now().Add(-10 * time.Minute)
	for key, entry := range l.limiters {
		if entry.lastSeen.Before(threshold) {
			delete(l.limiters, key)
		}
	}
	l.mu.Unlock()
}
