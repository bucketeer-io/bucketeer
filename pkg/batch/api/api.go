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

package api

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/proto/batch"
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
	mauPartitionDeleter       jobs.Job
	mauPartitionCreator       jobs.Job
	featureFlagCacher         jobs.Job
	segmentUserCacher         jobs.Job
	apiKeyCacher              jobs.Job
	experimentCacher          jobs.Job
	autoOpsRulesCacher        jobs.Job
	tagDeleter                jobs.Job
	featureAutoArchiver       jobs.Job
	logger                    *zap.Logger
}

func NewBatchService(
	experimentStatusUpdater, experimentRunningWatcher,
	featureStaleWatcher, mauCountWatcher, datetimeWatcher,
	eventCountWatcher, progressiveRolloutWatcher,
	redisCounterDeleter, experimentCalculator,
	mauSummarizer, mauPartitionDeleter, mauPartitionCreator,
	featureFlagCacher, segmentUserCacher, apiKeyCacher,
	experimentCacher, autoOpsRulesCacher, tagDeleter,
	featureAutoArchiver jobs.Job,
	logger *zap.Logger,
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
		mauPartitionDeleter:       mauPartitionDeleter,
		mauPartitionCreator:       mauPartitionCreator,
		featureFlagCacher:         featureFlagCacher,
		segmentUserCacher:         segmentUserCacher,
		apiKeyCacher:              apiKeyCacher,
		experimentCacher:          experimentCacher,
		autoOpsRulesCacher:        autoOpsRulesCacher,
		tagDeleter:                tagDeleter,
		featureAutoArchiver:       featureAutoArchiver,
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
	case batch.BatchJob_RedisCounterDeleter:
		err = s.redisCounterDeleter.Run(ctx)
	case batch.BatchJob_ProgressiveRolloutWatcher:
		err = s.progressiveRolloutWatcher.Run(ctx)
	case batch.BatchJob_ExperimentCalculator:
		err = s.experimentCalculator.Run(ctx)
	case batch.BatchJob_MauSummarizer:
		err = s.mauSummarizer.Run(ctx)
	case batch.BatchJob_MauPartitionDeleter:
		err = s.mauPartitionDeleter.Run(ctx)
	case batch.BatchJob_MauPartitionCreator:
		err = s.mauPartitionCreator.Run(ctx)
	case batch.BatchJob_FeatureFlagCacher:
		err = s.featureFlagCacher.Run(ctx)
	case batch.BatchJob_SegmentUserCacher:
		err = s.segmentUserCacher.Run(ctx)
	case batch.BatchJob_ApiKeyCacher:
		err = s.apiKeyCacher.Run(ctx)
	case batch.BatchJob_ExperimentCacher:
		err = s.experimentCacher.Run(ctx)
	case batch.BatchJob_AutoOpsRulesCacher:
		err = s.autoOpsRulesCacher.Run(ctx)
	case batch.BatchJob_TagDeleter:
		err = s.tagDeleter.Run(ctx)
	case batch.BatchJob_FeatureAutoArchiver:
		err = s.featureAutoArchiver.Run(ctx)
	default:
		s.logger.Error("Unknown job",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.String("job_name", req.Job.String()),
			)...,
		)
		return nil, errUnknownJob
	}
	if err != nil {
		// Don't log ResourceExhausted as ERROR - it's expected when pod is busy
		if st, ok := status.FromError(err); ok && st.Code() == codes.ResourceExhausted {
			s.logger.Debug("Job skipped, pod busy",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("job_name", req.Job.String()),
				)...,
			)
		} else {
			s.logger.Error("Failed to run the job",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.String("job_name", req.Job.String()),
					zap.Error(err),
				)...,
			)
		}
		return nil, err
	}
	return &batch.BatchJobResponse{}, nil
}

func (s *batchService) Register(server *grpc.Server) {
	batch.RegisterBatchServiceServer(server, s)
}
