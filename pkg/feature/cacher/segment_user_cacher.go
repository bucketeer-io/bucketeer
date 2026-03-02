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

package cacher

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	ftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

// SegmentUserCacher provides functionality to sync segment users from MySQL to Redis.
// This is used by the batch job to periodically refresh the cache for all environments.
//
//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
type SegmentUserCacher interface {
	// RefreshAllEnvironmentCaches updates the Redis cache for all environments.
	// This is used by the periodic batch job.
	RefreshAllEnvironmentCaches(ctx context.Context) error
}

type segmentUserCacher struct {
	segStorage ftstorage.SegmentStorage
	caches     []cachev3.SegmentUsersCache
	logger     *zap.Logger
}

// NewSegmentUserCacher creates a new SegmentUserCacher.
func NewSegmentUserCacher(
	mysqlClient mysql.Client,
	multiCaches []cache.MultiGetCache,
	logger *zap.Logger,
) SegmentUserCacher {
	caches := make([]cachev3.SegmentUsersCache, 0, len(multiCaches))
	for _, c := range multiCaches {
		caches = append(caches, cachev3.NewSegmentUsersCache(c))
	}
	return &segmentUserCacher{
		segStorage: ftstorage.NewSegmentStorage(mysqlClient),
		caches:     caches,
		logger:     logger.Named("segment-user-cacher"),
	}
}

// RefreshAllEnvironmentCaches updates the Redis cache for all environments.
func (c *segmentUserCacher) RefreshAllEnvironmentCaches(ctx context.Context) error {
	startTime := time.Now()

	// First, get all in-use segments (lightweight query)
	inUseSegments, err := c.segStorage.ListAllInUseSegments(ctx)
	if err != nil {
		c.logger.Error("Failed to list all in-use segments")
		recordListFeatures(cacherTypeSegmentUser, scopeBatch, environmentIDAll, codeFail, time.Since(startTime).Seconds())
		return err
	}
	recordListFeatures(cacherTypeSegmentUser, scopeBatch, environmentIDAll, codeSuccess, time.Since(startTime).Seconds())

	// Then, for each segment, fetch its users and cache them
	// This avoids loading all users in a single query which could be problematic
	// for large datasets (200k+ users)
	for _, seg := range inUseSegments {
		segStartTime := time.Now()
		users, err := c.segStorage.ListSegmentUsersBySegment(ctx, seg.SegmentID, seg.EnvironmentID)
		if err != nil {
			c.logger.Error("Failed to list segment users",
				zap.String("environmentId", seg.EnvironmentID),
				zap.String("segmentId", seg.SegmentID),
				zap.Error(err),
			)
			recordListFeatures(
				cacherTypeSegmentUser, scopeSingle, seg.EnvironmentID, codeFail,
				time.Since(segStartTime).Seconds(),
			)
			// Continue with other segments even if one fails
			continue
		}
		recordListFeatures(
			cacherTypeSegmentUser, scopeSingle, seg.EnvironmentID, codeSuccess,
			time.Since(segStartTime).Seconds(),
		)

		segUsers := &ftproto.SegmentUsers{
			SegmentId: seg.SegmentID,
			Users:     users,
			UpdatedAt: seg.UpdatedAt,
		}
		c.putCache(segUsers, seg.EnvironmentID, len(users))
	}

	return nil
}

// putCache saves segment users to all Redis instances and records metrics.
func (c *segmentUserCacher) putCache(segmentUsers *ftproto.SegmentUsers, environmentID string, userCount int) {
	var wg sync.WaitGroup
	var hasError bool
	var mu sync.Mutex

	for _, cache := range c.caches {
		wg.Add(1)
		go func(cache cachev3.SegmentUsersCache) {
			defer wg.Done()
			if err := cache.Put(segmentUsers, environmentID); err != nil {
				c.logger.Error("Failed to cache segment users",
					zap.Error(err),
					zap.String("environmentId", environmentID),
					zap.String("segmentId", segmentUsers.SegmentId),
				)
				mu.Lock()
				hasError = true
				mu.Unlock()
			}
		}(cache)
	}
	wg.Wait()

	// Record metrics based on overall success/failure
	if hasError {
		recordCachePut(cacherTypeSegmentUser, environmentID, codeFail)
	} else {
		recordCachePut(cacherTypeSegmentUser, environmentID, codeSuccess)
		recordFeaturesUpdated(cacherTypeSegmentUser, environmentID, userCount)
	}
}
