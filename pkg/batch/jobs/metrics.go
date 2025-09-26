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

package jobs

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

const (
	codeOK = "OK"

	// Job names
	JobExperimentCalculator      = "experiment_calculator"
	JobExperimentStatusUpdater   = "experiment_status_updater"
	JobFeatureStaleWatcher       = "feature_stale_watcher"
	JobMAUSummarizer             = "mau_summarizer"
	JobMAUPartitionDeleter       = "mau_partition_deleter"
	JobMAUPartitionCreator       = "mau_partition_creator"
	JobSegmentUsersUploader      = "segment_users_uploader"
	JobRedisCounterDeleter       = "redis_counter_deleter"
	JobFeatureFlagCacher         = "feature_flag_cacher"
	JobExperimentCacher          = "experiment_cacher"
	JobAPIKeyCacher              = "api_key_cacher"
	JobAutoOpsRulesCacher        = "auto_ops_rules_cacher"
	JobSegmentUserCacher         = "segment_user_cacher"
	JobTagDeleter                = "tag_deleter"
	JobExperimentRunningWatcher  = "experiment_running_watcher"
	JobMAUCountWatcher           = "mau_count_watcher"
	JobDatetimeWatcher           = "datetime_watcher"
	JobEventCountWatcher         = "event_count_watcher"
	JobProgressiveRolloutWatcher = "progressive_rollout_watcher"

	// Error types
	ErrorTypeTimeout    = "Timeout"
	ErrorTypeInternal   = "Internal"
	ErrorTypeNotFound   = "NotFound"
	ErrorTypeValidation = "Validation"
)

var (
	registerOnce    sync.Once
	batchJobCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "batch",
			Name:      "job_executions_total",
			Help:      "Total number of batch job executions",
		}, []string{"job", "code"})

	batchJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "batch",
			Name:      "job_duration_seconds",
			Help:      "Duration of batch job execution",
			Buckets:   prometheus.DefBuckets,
		}, []string{"job"})

	batchJobErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "batch",
			Name:      "job_errors_total",
			Help:      "Total number of batch job errors by type",
		}, []string{"job", "error_type"})
)

// RegisterMetrics registers batch metrics
func RegisterMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			batchJobCounter,
			batchJobDuration,
			batchJobErrorCounter,
		)
	})
}

func GetErrorType(err error) string {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return ErrorTypeTimeout
	}

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			return ErrorTypeNotFound
		case codes.InvalidArgument, codes.FailedPrecondition:
			return ErrorTypeValidation
		case codes.DeadlineExceeded:
			return ErrorTypeTimeout
		default:
			return ErrorTypeInternal
		}
	}

	if errors.Is(err, cache.ErrNotFound) {
		return ErrorTypeNotFound
	}

	if errors.Is(err, domain.ErrExperimentBeforeStart) ||
		errors.Is(err, domain.ErrExperimentBeforeStop) {
		return ErrorTypeValidation
	}

	return ErrorTypeInternal
}

// RecordJob is a helper that records job execution with error handling
func RecordJob(jobName string, err error, duration time.Duration) {
	recordJobExecution(jobName, codeOK, duration)
	if err != nil {
		errorType := GetErrorType(err)
		recordJobError(jobName, errorType)
	}
}

// recordJobExecution records the execution of a batch job
func recordJobExecution(jobName string, code string, duration time.Duration) {
	batchJobCounter.WithLabelValues(jobName, code).Inc()
	batchJobDuration.WithLabelValues(jobName).Observe(duration.Seconds())
}

// recordJobError records an error for a batch job
func recordJobError(jobName string, errorType string) {
	batchJobErrorCounter.WithLabelValues(jobName, errorType).Inc()
}
