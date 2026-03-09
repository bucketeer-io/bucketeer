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
	dauTTL = 35 * 24 * time.Hour // 35 days
)

// DAURecord groups user IDs by date, environment, and source.
type DAURecord struct {
	Date     time.Time
	EnvID    string
	SourceID string
	UserIDs  []string
}

// DAUCache provides DAU counting operations.
type DAUCache interface {
	// RecordDAUBatch adds user IDs to DAU HyperLogLog in
	// a single Redis Pipeline (1 round-trip).
	RecordDAUBatch(records []DAURecord) error
}

type dauCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewDAUCache(c cache.MultiGetDeleteCountCache) DAUCache {
	return &dauCache{cache: c}
}

// dauKey builds the Redis key for DAU HyperLogLog.
// Format: envId:dau:sourceId:yyyyMMdd
//
// Hash tags ({}) are not required in the key because pkg/redis/v3/redis.go
// PFMerge handles cross-slot merging transparently.
func dauKey(envID, sourceID string, date time.Time) string {
	return fmt.Sprintf("%s:dau:%s:%s", envID, sourceID, date.Format("20060102"))
}

func (c *dauCache) RecordDAUBatch(records []DAURecord) error {
	// Filter out empty records to avoid creating a Pipeline unnecessarily.
	filtered := make([]DAURecord, 0, len(records))
	for _, r := range records {
		if len(r.UserIDs) > 0 {
			filtered = append(filtered, r)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	pipe := c.cache.Pipeline(false)
	for _, r := range filtered {
		key := dauKey(r.EnvID, r.SourceID, r.Date)
		pipe.PFAdd(key, r.UserIDs...)
		pipe.Expire(key, dauTTL)
	}
	if _, err := pipe.Exec(); err != nil {
		return fmt.Errorf("failed to record DAU batch (%d keys): %w", len(filtered), err)
	}
	return nil
}
