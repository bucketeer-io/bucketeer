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

package stream

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

const (
	methodStreamEvaluations = "StreamEvaluations"

	errorTypeEvaluationPut            = "evaluation_put"
	errorTypeEvaluationPatch          = "evaluation_patch"
	errorTypeConnectionRefusedByLimit = "connection_refused_by_limit"

	patchCodeDiff = "Diff"
	patchCodeNone = "None"
)

var (
	sseActiveConnectionsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_active_connections",
			Help:      "Current number of active SSE connections held in the stream dispatcher.",
		}, []string{"environment_id", "tag", "source_id"})
	sseDispatchDroppedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_dispatch_dropped_total",
			Help:      "Total dispatched events dropped because a connection's buffer was full.",
		}, []string{"environment_id", "tag", "source_id"})
	sseConnectionDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_connection_duration_seconds",
			Help:      "Duration of SSE connections from register to deregister.",
			Buckets:   []float64{1, 5, 15, 30, 60, 300, 600, 1800, 3600, 7200, 21600, 86400},
		}, []string{"environment_id", "tag", "source_id"})
	sseDispatchTagsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_dispatch_tags",
			Help:      "Number of tags affected per dispatch call triggered by domain events.",
			Buckets:   []float64{1, 2, 3, 5, 7, 10, 15, 20, 30},
		}, []string{"environment_id", "event_type"})
	ssePatchCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_patch_total",
			Help:      "Total patch events by evaluation result type.",
		}, []string{"environment_id", "tag", "source_id", "evaluation_type"})
	sseErrorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_errors_total",
			Help:      "Total errors during SSE stream processing.",
		}, []string{"environment_id", "tag", "source_id", "error_type"})
	sseEvaluationDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_evaluation_duration_seconds",
			Help:      "Time spent evaluating features for an SSE put or patch event.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"environment_id", "tag", "source_id", "event_type"})
	sseInitialPutDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_initial_put_duration_seconds",
			Help:      "Time from request start to initial put event sent.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"environment_id", "tag", "source_id"})
	sseDispatchToSendDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_dispatch_to_send_duration_seconds",
			Help: "Time from dispatch to send completion per connection," +
				" including queuing, evaluation, marshaling, and flushing.",
			Buckets: prometheus.DefBuckets,
		}, []string{"environment_id", "tag", "source_id"})
)

func RegisterMetrics(r metrics.Registerer) {
	r.MustRegister(
		sseActiveConnectionsGauge,
		sseDispatchDroppedCounter,
		sseConnectionDurationHistogram,
		sseDispatchTagsHistogram,
		ssePatchCounter,
		sseErrorsCounter,
		sseEvaluationDurationHistogram,
		sseInitialPutDurationHistogram,
		sseDispatchToSendDurationHistogram,
	)
}
