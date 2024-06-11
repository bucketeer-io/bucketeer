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

package storage

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

var (
	// TODO: Remove this event after grpc server is removed.
	sdkGetEvaluationsLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_get_evaluations_handling_seconds",
			Help:      "Histogram of get evaluations response latency (seconds).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"tag", "state"})

	// TODO: Remove this event after grpc server is removed.
	sdkGetEvaluationsSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_get_evaluations_size",
			Help:      "Histogram of get evaluations response size (byte).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"tag", "state"})

	// TODO: Remove this event after grpc server is removed.
	sdkTimeoutErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_timeout_error_total",
			Help:      "Total number of sdk timeout errors",
		}, []string{"tag"})

	// TODO: Remove this event after grpc server is removed.
	sdkInternalErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_internal_error_total",
			Help:      "Total number of sdk internal errors",
		}, []string{"tag"})

	// We added prefix "sdk_api" because of avoiding name collision
	sdkLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_handling_seconds",
			Help:      "Histogram of get evaluations response latency (seconds).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"tag", "state", "api", "sdk_version"})

	sdkSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_response_size",
			Help:      "Histogram of get evaluations response size (byte).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"tag", "state", "api", "sdk_version"})

	sdkTimeoutError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_timeout_error_total",
			Help:      "Total number of sdk timeout errors",
		}, []string{"tag", "api", "sdk_version"})

	sdkInternalError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_internal_error_total",
			Help:      "Total number of sdk internal errors",
		}, []string{"tag", "api", "sdk_version"})

	sdkNetworkError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_network_error_total",
			Help:      "Total number of sdk network errors",
		}, []string{"tag", "api", "sdk_version"})

	sdkInternalSdkError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_internal_sdk_error_total",
			Help:      "Total number of sdk internal sdk errors",
		}, []string{"tag", "api", "sdk_version"})
)

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		sdkGetEvaluationsLatencyHistogram,
		sdkGetEvaluationsSizeHistogram,
		sdkTimeoutErrorCounter,
		sdkInternalErrorCounter,
		sdkLatencyHistogram,
		sdkSizeHistogram,
		sdkTimeoutError,
		sdkInternalError,
		sdkNetworkError,
		sdkInternalSdkError,
	)
}
