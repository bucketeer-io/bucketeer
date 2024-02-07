// Copyright 2024 The Bucketeer Authors.
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

package mau

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/mau/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

type mauPartitionCreator struct {
	mauStorage storage.MAUStorage
	location   *time.Location
	opts       *jobs.Options
	logger     *zap.Logger
}

func NewMAUPartitionCreator(
	mysqlClient mysql.Client,
	location *time.Location,
	opts ...jobs.Option) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 5 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &mauPartitionCreator{
		mauStorage: storage.NewMAUStorage(mysqlClient),
		location:   location,
		opts:       dopts,
		logger:     dopts.Logger.Named("mau-partition-creator"),
	}
}

func (d *mauPartitionCreator) Run(ctx context.Context) error {
	d.logger.Info("mauPartitionCreator start running")
	ctx, cancel := context.WithTimeout(ctx, d.opts.Timeout)
	defer cancel()
	now := time.Now().In(d.location)
	targetPartition := d.newPartitionName(now.AddDate(0, 1, 0))
	lessThan := d.newLessThan(now.AddDate(0, 2, 0))
	if err := d.mauStorage.CreatePartition(ctx, targetPartition, lessThan); err != nil {
		d.logger.Error("Failed to create MAU partition",
			zap.Error(err),
			zap.String("targetPartition", targetPartition),
			zap.String("lessThan", lessThan),
		)
		return err
	}
	d.logger.Info("Succeeded to create MAU partition",
		zap.String("targetPartition", targetPartition),
		zap.String("lessThan", lessThan),
		zap.Duration("elapsedTime", time.Since(now)),
	)
	d.logger.Info("mauPartitionCreator start stopping")
	return nil
}

func (d *mauPartitionCreator) newPartitionName(target time.Time) string {
	return fmt.Sprintf("p%d%02d", target.Year(), target.Month())
}

func (d *mauPartitionCreator) newLessThan(target time.Time) string {
	return fmt.Sprintf("%d%02d", target.Year(), target.Month())
}
