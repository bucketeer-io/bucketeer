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
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	callerGatewayService = "GatewayService"

	methodGetEvaluations = "GetEvaluations"
	methodGetEvaluation  = "GetEvaluation"
	methodRegisterEvents = "RegisterEvent"
	methodTrack          = "Track"

	typeFeatures      = "Features"
	typeSegmentUsers  = "SegmentUsers"
	typeAPIKey        = "APIKey"
	typeRegisterEvent = "RegisterEvent"
	typeEvaluation    = "Evaluation"
	typeGoal          = "Goal"
	typeMetrics       = "Metrics"
	typeUnknown       = "Unknown"
	typeTrack         = "Track"

	cacheLayerInMemory = "InMemory"
	cacheLayerExternal = "External"

	codeHit  = "Hit"
	codeMiss = "Miss"

	codeOK                      = "OK"
	codeInvalidID               = "InvalidID"
	codeInvalidTimestamp        = "InvalidTimestamp"
	codeInvalidTimestampRequest = "InvalidTimestampRequest"
	codeUnmarshalFailed         = "UnmarshalFailed"
	codeMarshalAnyFailed        = "MarshalAnyFailed"
	codeInvalidType             = "InvalidType"
	codeNonRepeatableError      = "NonRepeatableError"
	codeRepeatableError         = "RepeatableError"
	codeInvalidURLParams        = "InvalidURLParams"
	codeEmptyField              = "EmptyField"

	codeAll           = "All"
	codeDiff          = "Diff"
	codeNone          = "None"
	codeNoFeatures    = "NoFeatures"
	codeNoSegments    = "NoSegments"
	codeOld           = "Old"
	codeInternalError = "InternalError"
	codeBadRequest    = "BadRequest"
)

var (
	registerOnce sync.Once
	cacheCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_cache_requests_total",
			Help:      "Total number of cache requests",
		}, []string{"caller", "type", "layer", "code"})
	eventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_register_events_total",
			Help:      "Total number of registered events",
		}, []string{"caller", "type", "code"})
	evaluationsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_evaluations_total",
			Help:      "Total number of evaluations",
			// TODO: Remove project_id.
		}, []string{
			"project_id", "project_url_code",
			"environment_id", "environment_url_code", "tag", "evaluation_type",
		})
	getFeatureFlagsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_get_feature_flags_total",
			Help:      "Total number of get feature flags",
			// TODO: Remove project_id.
		}, []string{
			"project_id", "project_url_code",
			"environment_id", "environment_url_code", "tag", "response_type",
		})
	getSegmentUsersCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_get_segment_users_total",
			Help:      "Total number of get segment users api requests",
			// TODO: Remove project_id.
		}, []string{
			"project_id", "project_url_code", "environment_id", "environment_url_code",
			"source_id", "sdk_version", "response_type",
		})
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_request_total",
			Help:      "Total number of request",
			// TODO: Remove project_id.
		}, []string{
			"organization_id", "project_id", "project_url_code",
			"environment_id", "environment_url_code", "method", "source_id",
		})
	// TODO: Remove after deleting api-gateway REST server
	restCacheCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_rest_cache_requests_total",
			Help:      "Total number of cache requests",
		}, []string{"caller", "type", "layer", "code"})
	// TODO: Remove after deleting api-gateway REST server
	restEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_rest_register_events_total",
			Help:      "Total number of registered events",
		}, []string{"caller", "type", "code"})
	// TODO: Remove after deleting api-gateway REST server
	sdkGetEvaluationsLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_get_evaluations_handling_seconds",
			Help:      "Histogram of get evaluations response latency (seconds).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"environment_id", "tag", "state"})
	// TODO: Remove after deleting api-gateway REST server
	sdkGetEvaluationsSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_get_evaluations_size",
			Help:      "Histogram of get evaluations response size (byte).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"environment_id", "tag", "state"})
	// TODO: Remove after deleting api-gateway REST server
	sdkTimeoutErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_timeout_error_total",
			Help:      "Total number of sdk timeout errors",
		}, []string{"environment_id", "tag"})
	// TODO: Remove after deleting api-gateway REST server
	sdkInternalErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_internal_error_total",
			Help:      "Total number of sdk internal errors",
		}, []string{"environment_id", "tag"})
	sdkLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_handling_seconds",
			Help:      "Histogram of get evaluations response latency (seconds).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"project_id", "environment_id", "tag", "api", "sdk_version", "source_id"})
	sdkSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_response_size",
			Help:      "Histogram of get evaluations response size (byte).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"project_id", "environment_id", "tag", "api", "sdk_version", "source_id"})
	sdkErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "metrics_event",
			Name:      "sdk_api_error_total",
			Help:      "Total number of sdk errors",
		}, []string{"project_id", "environment_id", "tag", "error_type", "api", "sdk_version", "source_id"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			cacheCounter,
			eventCounter,
			evaluationsCounter,
			getFeatureFlagsCounter,
			getSegmentUsersCounter,
			requestTotal,
			sdkGetEvaluationsLatencyHistogram,
			sdkGetEvaluationsSizeHistogram,
			sdkTimeoutErrorCounter,
			sdkInternalErrorCounter,
			sdkLatencyHistogram,
			sdkSizeHistogram,
			sdkErrorCounter,
		)
	})
}
