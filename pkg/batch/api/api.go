// Copyright 2023 The Bucketeer Authors.
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

package api

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/experiment"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/notification"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/opsevent"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/proto/batch"
)

var (
	errUnknownJob = errors.New("batch: unknown job")
)

type BatchService struct {
	experimentStatusUpdaterJob *experiment.ExperimentStatusUpdaterJob
	experimentRunningWatcher   *notification.ExperimentRunningWatcherJob
	featureWatcherJob          *notification.FeatureWatcherJob
	mauCountWatcherJob         *notification.MauCountWatcherJob
	datetimeWatcherJob         *opsevent.DatetimeWatcherJob
	eventCountWatcherJob       *opsevent.EventCountWatcherJob
	logger                     *zap.Logger
}

func NewBatchService(
	experimentStatusUpdaterJob *experiment.ExperimentStatusUpdaterJob,
	experimentRunningWatcher *notification.ExperimentRunningWatcherJob,
	featureWatcherJob *notification.FeatureWatcherJob,
	mauCountWatcherJob *notification.MauCountWatcherJob,
	eventWatcherJob *opsevent.DatetimeWatcherJob,
	eventCountWatcherJob *opsevent.EventCountWatcherJob,
	logger *zap.Logger,
) *BatchService {
	return &BatchService{
		experimentStatusUpdaterJob: experimentStatusUpdaterJob,
		experimentRunningWatcher:   experimentRunningWatcher,
		featureWatcherJob:          featureWatcherJob,
		mauCountWatcherJob:         mauCountWatcherJob,
		datetimeWatcherJob:         eventWatcherJob,
		eventCountWatcherJob:       eventCountWatcherJob,
		logger:                     logger.Named("batch-service"),
	}
}

func (s *BatchService) ExecuteBatchJob(
	ctx context.Context, req *batch.BatchJobRequest) (*batch.BatchJobResponse, error) {
	var err error
	resp := &batch.BatchJobResponse{}
	switch req.Job {
	case batch.BatchJob_ExprimentStatusUpdater:
		err = s.experimentStatusUpdaterJob.Run(ctx)
	case batch.BatchJob_ExperimentRunningWatcher:
		err = s.experimentRunningWatcher.Run(ctx)
	case batch.BatchJob_FeatureStateWatcher:
		err = s.featureWatcherJob.Run(ctx)
	case batch.BatchJob_MauCountWatcher:
		err = s.mauCountWatcherJob.Run(ctx)
	case batch.BatchJob_DatetimeWatcher:
		err = s.datetimeWatcherJob.Run(ctx)
	case batch.BatchJob_EventCountWatcher:
		err = s.eventCountWatcherJob.Run(ctx)
	default:
		s.logger.Error("Batch Service unknown job",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("job_name", req.Job.String()),
			)...)
		err = errUnknownJob
	}
	return resp, err
}

func (s *BatchService) Register(server *grpc.Server) {
	batch.RegisterBatchServiceServer(server, s)
}
