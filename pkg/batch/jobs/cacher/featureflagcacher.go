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
//

package cacher

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	ftcacher "github.com/bucketeer-io/bucketeer/v2/pkg/feature/cacher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
)

type featureFlagCacherJob struct {
	cacher ftcacher.FeatureFlagCacher
	opts   *jobs.Options
	logger *zap.Logger
}

// NewFeatureFlagCacher creates a new feature flag cacher batch job.
func NewFeatureFlagCacher(
	mysqlClient mysql.Client,
	multiCaches []cache.MultiGetCache,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &featureFlagCacherJob{
		cacher: ftcacher.NewFeatureFlagCacher(mysqlClient, multiCaches, dopts.Logger),
		opts:   dopts,
		logger: dopts.Logger.Named("feature-flag-cacher-job"),
	}
}

// Run executes the batch job to update all environments' feature flag caches.
func (c *featureFlagCacherJob) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobFeatureFlagCacher, lastErr, time.Since(startTime))
	}()
	return c.cacher.RefreshAllEnvironmentCaches(ctx)
}
