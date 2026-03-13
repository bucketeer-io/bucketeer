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
	// RecordDAUBatch adds user IDs to DAU HyperLogLog.
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
	for _, r := range records {
		if len(r.UserIDs) == 0 {
			continue
		}
		key := dauKey(r.EnvID, r.SourceID, r.Date)
		if _, err := c.cache.PFAdd(key, r.UserIDs...); err != nil {
			return fmt.Errorf("failed to PFAdd DAU for key %s: %w", key, err)
		}
		if _, err := c.cache.Expire(key, dauTTL); err != nil {
			return fmt.Errorf("failed to set TTL for key %s: %w", key, err)
		}
	}
	return nil
}
