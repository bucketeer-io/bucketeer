// Copyright 2025 The Bucketeer Authors.
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

package deleter

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	tagstorage "github.com/bucketeer-io/bucketeer/pkg/tag/storage"
	ftproto "github.com/bucketeer-io/bucketeer/proto/feature"
	tagproto "github.com/bucketeer-io/bucketeer/proto/tag"
)

var (
	errInternal = errors.New("batch: internal error")
)

type tagDeleter struct {
	tagStorage tagstorage.TagStorage
	ftStorage  ftstorage.FeatureStorage
	opts       *jobs.Options
	logger     *zap.Logger
}

func NewTagDeleter(
	mysqlClient mysql.Client,
	opts ...jobs.Option) jobs.Job {

	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &tagDeleter{
		tagStorage: tagstorage.NewTagStorage(mysqlClient),
		ftStorage:  ftstorage.NewFeatureStorage(mysqlClient),
		opts:       dopts,
		logger:     dopts.Logger.Named("tag-deleter"),
	}
}

func (td *tagDeleter) Run(ctx context.Context) (lastErr error) {
	ctx, cancel := context.WithTimeout(ctx, td.opts.Timeout)
	defer cancel()

	td.logger.Info("Starting to delete unused tags")
	startTime := time.Now()
	// List all the tags by environment
	envTags, err := td.tagStorage.ListAllEnvironmentTags(ctx)
	if err != nil {
		td.logger.Error("Failed to list all environment tags", zap.Error(err))
		return errInternal
	}
	// List all the features by environment
	envFts, err := td.ftStorage.ListAllEnvironmentFeatures(ctx)
	if err != nil {
		td.logger.Error("Failed to list all environment features", zap.Error(err))
		return errInternal
	}

	var deletedSize int
	for _, envTag := range envTags {
		// Delete all the tags when there are no flags
		if len(envFts) == 0 {
			for _, tag := range envTag.Tags {
				// Delete unused tag
				if err := td.deleteTag(ctx, tag); err != nil {
					lastErr = err
					continue
				}
				deletedSize++
			}
		} else {
			// Check if the tags are in use in all the flags
			for _, envFt := range envFts {
				if envFt.EnvironmentId == envTag.EnvironmentId {
					deletedCount, err := td.deleteUnusedTags(ctx, envTag.Tags, envFt.Features)
					if err != nil {
						lastErr = err
						continue
					}
					deletedSize += deletedCount
				}
			}
		}
	}

	if lastErr != nil {
		td.logger.Error("Finished deleting unused tags with errors",
			zap.Error(lastErr),
			zap.Duration("elapsedTime", time.Since(startTime)),
			zap.Int("deletedSize", deletedSize),
		)
	} else {
		td.logger.Info("Finished deleting unused tags",
			zap.Duration("elapsedTime", time.Since(startTime)),
			zap.Int("deletedSize", deletedSize),
		)
	}
	return
}

func (td *tagDeleter) deleteUnusedTags(
	ctx context.Context,
	tags []*tagproto.Tag,
	fts []*ftproto.Feature,
) (int, error) {
	var deletedCount int
	for _, tag := range tags {
		var inUse bool
		for _, ft := range fts {
			// Check if the tag is in use
			if td.contains(tag.Name, ft.Tags) {
				inUse = true
				break
			}
		}
		// Skip deleting if the tag is in use
		if inUse {
			continue
		}
		// Delete unused tag
		if err := td.deleteTag(ctx, tag); err != nil {
			return deletedCount, err
		}
		deletedCount++
	}
	return deletedCount, nil
}

func (td *tagDeleter) deleteTag(ctx context.Context, tag *tagproto.Tag) error {
	if err := td.tagStorage.DeleteTag(ctx, tag.Id); err != nil {
		td.logger.Error("Failed to delete the tag",
			zap.Error(err),
			zap.String("tagId", tag.Id),
			zap.String("tagName", tag.Name),
			zap.String("environmentId", tag.EnvironmentId),
		)
		return errInternal
	}
	td.logger.Debug("Deleted tag successfully",
		zap.String("tagId", tag.Id),
		zap.String("tagName", tag.Name),
		zap.String("environmentId", tag.EnvironmentId),
	)
	return nil
}

func (td *tagDeleter) contains(needle string, haystack []string) bool {
	for i := range haystack {
		if haystack[i] == needle {
			return true
		}
	}
	return false
}
