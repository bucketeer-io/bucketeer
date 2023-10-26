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

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs/notification"
	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	notificationsender "github.com/bucketeer-io/bucketeer/pkg/notification/sender"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/proto/batch"
)

var (
	errUnknownJob = status.Error(codes.InvalidArgument, "batch: unknown job")
)

type batchService struct {
	experimentStatusUpdater   jobs.Job
	experimentRunningWatcher  jobs.Job
	featureStaleWatcher       jobs.Job
	mauCountWatcher           jobs.Job
	datetimeWatcher           jobs.Job
	countWatcher              jobs.Job
	progressiveRolloutWatcher jobs.Job
	redisCounterDeleter       jobs.Job
	experimentCalculator      jobs.Job
	mauSummarizer             jobs.Job
	environmentClient         envclient.Client
	domainEventPuller         puller.Puller
	notificationSender        notificationsender.Sender
	notificationOpts          []notification.Option
	logger                    *zap.Logger
}

func NewBatchService(
	experimentStatusUpdater, experimentRunningWatcher,
	featureStaleWatcher, mauCountWatcher, datetimeWatcher,
	eventCountWatcher, progressiveRolloutWatcher jobs.Job,
	environmentClient envclient.Client,
	domainEventPuller puller.Puller,
	notificationSender notificationsender.Sender,
	redisCounterDeleter jobs.Job,
	experimentCalculator jobs.Job,
	mauSummarizer jobs.Job,
	logger *zap.Logger,
	options ...notification.Option,
) *batchService {
	return &batchService{
		experimentStatusUpdater:   experimentStatusUpdater,
		experimentRunningWatcher:  experimentRunningWatcher,
		featureStaleWatcher:       featureStaleWatcher,
		mauCountWatcher:           mauCountWatcher,
		datetimeWatcher:           datetimeWatcher,
		countWatcher:              eventCountWatcher,
		progressiveRolloutWatcher: progressiveRolloutWatcher,
		redisCounterDeleter:       redisCounterDeleter,
		experimentCalculator:      experimentCalculator,
		mauSummarizer:             mauSummarizer,
		environmentClient:         environmentClient,
		domainEventPuller:         domainEventPuller,
		notificationSender:        notificationSender,
		notificationOpts:          options,
		logger:                    logger.Named("batch-service"),
	}
}

func (s *batchService) ExecuteBatchJob(
	ctx context.Context, req *batch.BatchJobRequest) (*batch.BatchJobResponse, error) {
	var err error
	switch req.Job {
	case batch.BatchJob_ExperimentStatusUpdater:
		err = s.experimentStatusUpdater.Run(ctx)
	case batch.BatchJob_ExperimentRunningWatcher:
		err = s.experimentRunningWatcher.Run(ctx)
	case batch.BatchJob_FeatureStaleWatcher:
		err = s.featureStaleWatcher.Run(ctx)
	case batch.BatchJob_MauCountWatcher:
		err = s.mauCountWatcher.Run(ctx)
	case batch.BatchJob_DatetimeWatcher:
		err = s.datetimeWatcher.Run(ctx)
	case batch.BatchJob_EventCountWatcher:
		err = s.countWatcher.Run(ctx)
	case batch.BatchJob_DomainEventInformer:
		domainEventInformer := notification.NewDomainEventInformer(
			s.environmentClient,
			s.domainEventPuller,
			s.notificationSender,
			s.notificationOpts...,
		)
		err = domainEventInformer.Run(ctx)
	case batch.BatchJob_RedisCounterDeleter:
		err = s.redisCounterDeleter.Run(ctx)
	case batch.BatchJob_ProgressiveRolloutWatcher:
		err = s.progressiveRolloutWatcher.Run(ctx)
	case batch.BatchJob_ExperimentCalculator:
		err = s.experimentCalculator.Run(ctx)
	case batch.BatchJob_MauSummarizer:
		err = s.mauSummarizer.Run(ctx)
	default:
		s.logger.Error("Unknown job",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("job_name", req.Job.String()),
			)...,
		)
		return nil, errUnknownJob
	}
	if err != nil {
		s.logger.Error("Failed to run the job",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("job_name", req.Job.String()),
				zap.Error(err),
			)...,
		)
		return nil, err
	}
	return &batch.BatchJobResponse{}, nil
}

func (s *batchService) Register(server *grpc.Server) {
	batch.RegisterBatchServiceServer(server, s)
}
