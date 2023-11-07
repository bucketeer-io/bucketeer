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

const partitionRetentionMonths = 3

type mauPartitionDeleter struct {
	mauStorage storage.MAUStorage
	location   *time.Location
	opts       *jobs.Options
	logger     *zap.Logger
}

func NewMAUPartitionDeleter(
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
	return &mauPartitionDeleter{
		mauStorage: storage.NewMAUStorage(mysqlClient),
		location:   location,
		opts:       dopts,
		logger:     dopts.Logger.Named("mau-partition-deleter"),
	}
}

func (d *mauPartitionDeleter) Run(ctx context.Context) error {
	d.logger.Info("MAUPartitionDeleter start running")
	ctx, cancel := context.WithTimeout(ctx, d.opts.Timeout)
	defer cancel()
	now := time.Now().In(d.location)
	targetPartition := d.newPartitionName(d.getDeleteTargetPartition(now))
	d.logger.Info("Start to delete MAU records",
		zap.String("targetPartition", targetPartition),
	)
	// Delete records from the target partition before dropping partition.
	// This process doesn't block the reads and writes.
	if err := d.mauStorage.DeleteRecords(ctx, targetPartition); err != nil {
		d.logger.Error("Failed to delete MAU records",
			zap.Error(err),
			zap.String("targetPartition", targetPartition),
		)
		return err
	}
	d.logger.Info("Succeeded to delete MAU records",
		zap.String("targetPartition", targetPartition),
		zap.Duration("elapsedTime", time.Since(now)),
	)
	// Delete the actual files from the disk before dropping partition.
	// This process blocks only writes, but the duration is short.
	if err := d.mauStorage.RebuildPartition(ctx, targetPartition); err != nil {
		d.logger.Error("Failed to rebuild MAU partition",
			zap.Error(err),
			zap.String("targetPartition", targetPartition),
		)
		return err
	}
	d.logger.Info("Succeeded to rebuild MAU partition",
		zap.String("targetPartition", targetPartition),
		zap.Duration("elapsedTime", time.Since(now)),
	)
	// Drop the partition.
	// This process doesn't block the reads and writes.
	if err := d.mauStorage.DropPartition(ctx, targetPartition); err != nil {
		d.logger.Error("Failed to drop MAU partition",
			zap.Error(err),
			zap.String("targetPartition", targetPartition),
		)
		return err
	}
	d.logger.Info("Succeeded to drop MAU partition",
		zap.String("targetPartition", targetPartition),
		zap.Duration("elapsedTime", time.Since(now)),
	)
	d.logger.Info("MAUPartitionDeleter start stopping")
	return nil
}

func (d *mauPartitionDeleter) getDeleteTargetPartition(now time.Time) time.Time {
	m := -1 * partitionRetentionMonths
	base := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, d.location)
	return base.AddDate(0, m, 0)
}

func (d *mauPartitionDeleter) newPartitionName(target time.Time) string {
	return fmt.Sprintf("p%d%02d", target.Year(), target.Month())
}
