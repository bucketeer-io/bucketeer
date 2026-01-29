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
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
)

const (
	dauKeyPrefix = "dau"
	dauTTL       = 60 * 24 * time.Hour // 60 days
)

// MAUCache provides DAU/MAU counting operations for Insights Dashboard.
type MAUCache interface {
	// RecordDAU adds a user ID to the DAU HyperLogLog for the given date.
	// Key pattern: {envId}:{sourceId}:dau:{yyyyMMdd}
	RecordDAU(envID, sourceID, userID string, date time.Time) error
}

type mauCache struct {
	cache cache.MultiGetDeleteCountCache
}

// NewMAUCache creates a new MAUCache.
func NewMAUCache(c cache.MultiGetDeleteCountCache) MAUCache {
	return &mauCache{cache: c}
}

// dauKey builds the Redis key for DAU HyperLogLog.
// Format: {envId}:{sourceId}:dau:{yyyyMMdd}
func (*mauCache) dauKey(envID, sourceID string, date time.Time) string {
	dateStr := date.Format("20060102")
	return fmt.Sprintf("%s:%s:%s:%s", envID, sourceID, dauKeyPrefix, dateStr)
}

// RecordDAU adds a user ID to the DAU HyperLogLog using PFADD.
// Uses Pipeline for atomic PFAdd + Expire execution (same pattern as userAttributesCache).
func (c *mauCache) RecordDAU(envID, sourceID, userID string, date time.Time) error {
	if len(userID) == 0 {
		return nil
	}
	key := c.dauKey(envID, sourceID, date)
	pipe := c.cache.Pipeline(false)
	pipe.PFAdd(key, userID)
	pipe.Expire(key, dauTTL)
	if _, err := pipe.Exec(); err != nil {
		return fmt.Errorf("failed to record DAU: %w", err)
	}
	return nil
}
