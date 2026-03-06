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
	mauTTL = 65 * 24 * time.Hour // 65 days
)

// MAUCache provides MAU counting operations.
type MAUCache interface {
	// MergeIntoMAUBatch merges DAUs into MAUs for multiple source IDs using pipeline.
	// Returns a map of sourceID to MAU count.
	MergeIntoMAUBatch(envID string, sourceIDs []string, date time.Time) (map[string]int64, error)
}

type mauCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewMAUCache(c cache.MultiGetDeleteCountCache) MAUCache {
	return &mauCache{cache: c}
}

// mauKey builds the Redis key for MAU HyperLogLog.
// Format: envId:mau:sourceId:yyyyMM
func mauKey(envID, sourceID string, month time.Time) string {
	monthStr := month.Format("200601")
	return fmt.Sprintf("%s:mau:%s:%s", envID, sourceID, monthStr)
}

// MergeIntoMAUBatch merges DAUs into MAUs for multiple source IDs using pipeline.
func (c *mauCache) MergeIntoMAUBatch(envID string, sourceIDs []string, date time.Time) (map[string]int64, error) {
	if len(envID) == 0 || len(sourceIDs) == 0 {
		return make(map[string]int64), nil
	}

	pipe := c.cache.Pipeline(false)

	for _, sourceID := range sourceIDs {
		mk := mauKey(envID, sourceID, date)
		dk := dauKey(envID, sourceID, date)
		pipe.PFMerge(mk, mk, dk)
		pipe.Expire(mk, mauTTL)
	}

	countCmds := make([]any, len(sourceIDs))
	for i, sourceID := range sourceIDs {
		mk := mauKey(envID, sourceID, date)
		countCmds[i] = pipe.PFCount(mk)
	}

	if _, err := pipe.Exec(); err != nil {
		return nil, fmt.Errorf("failed to execute merging MAU batch: %w", err)
	}

	result := make(map[string]int64, len(sourceIDs))
	for i, sourceID := range sourceIDs {
		cmd := countCmds[i].(interface{ Val() int64 })
		result[sourceID] = cmd.Val()
	}
	return result, nil
}
