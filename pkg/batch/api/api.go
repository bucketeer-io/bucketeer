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

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/proto/batch"
)

var (
	errUnknownJob = errors.New("batch: unknown job")
)

type BatchService struct {
	experimentStatusUpdater  jobs.Job
	experimentRunningWatcher jobs.Job
	featureWatcher           jobs.Job
	mauCountWatcher          jobs.Job
	datetimeWatcher          jobs.Job
	countWatcher             jobs.Job
	logger                   *zap.Logger
}

func NewBatchService(
	experimentStatusUpdater, experimentRunningWatcher,
	featureWatcher, mauCountWatcher,
	datetimeWatcher, eventCountWatcher jobs.Job,
	logger *zap.Logger,
) *BatchService {
	return &BatchService{
		experimentStatusUpdater:  experimentStatusUpdater,
		experimentRunningWatcher: experimentRunningWatcher,
		featureWatcher:           featureWatcher,
		mauCountWatcher:          mauCountWatcher,
		datetimeWatcher:          datetimeWatcher,
		countWatcher:             eventCountWatcher,
		logger:                   logger.Named("batch-service"),
	}
}

func (s *BatchService) ExecuteBatchJob(
	ctx context.Context, req *batch.BatchJobRequest) (*batch.BatchJobResponse, error) {
	var err error
	resp := &batch.BatchJobResponse{}
	switch req.Job {
	case batch.BatchJob_ExperimentStatusUpdater:
		err = s.experimentStatusUpdater.Run(ctx)
	case batch.BatchJob_ExperimentRunningWatcher:
		err = s.experimentRunningWatcher.Run(ctx)
	case batch.BatchJob_FeatureStateWatcher:
		err = s.featureWatcher.Run(ctx)
	case batch.BatchJob_MauCountWatcher:
		err = s.mauCountWatcher.Run(ctx)
	case batch.BatchJob_DatetimeWatcher:
		err = s.datetimeWatcher.Run(ctx)
	case batch.BatchJob_EventCountWatcher:
		err = s.countWatcher.Run(ctx)
	default:
		s.logger.Error("Batch Service unknown job",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("job_name", req.Job.String()),
			)...,
		)
		err = errUnknownJob
	}
	return resp, err
}

func (s *BatchService) Register(server *grpc.Server) {
	batch.RegisterBatchServiceServer(server, s)
}
