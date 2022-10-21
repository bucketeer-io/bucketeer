// Copyright 2022 The Bucketeer Authors.
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

	typeFeatures      = "Features"
	typeSegmentUsers  = "SegmentUsers"
	typeAPIKey        = "APIKey"
	typeRegisterEvent = "RegisterEvent"
	typeEvaluation    = "Evaluation"
	typeGoal          = "Goal"
	typeGoalBatch     = "GoalBatch"
	typeMetrics       = "Metrics"
	typeUnknown       = "Unknown"
	typeTrack         = "Track"

	cacheLayerExternal = "External"

	codeHit  = "Hit"
	codeMiss = "Miss"

	codeOK                         = "OK"
	codeInvalidID                  = "InvalidID"
	codeInvalidTimestamp           = "InvalidTimestamp"
	codeEmptyTag                   = "EmptyTag"
	codeEmptyUserID                = "EmptyUserID"
	codeInvalidTimestampRequest    = "InvalidTimestampRequest"
	codeUpsertUserEvaluationFailed = "UpsertUserEvaluationFailed"
	codeUnmarshalFailed            = "UnmarshalFailed"
	codeMarshalAnyFailed           = "MarshalAnyFailed"
	codeEvaluationConversionFailed = "EvaluationConversionFailed"
	codeInvalidType                = "InvalidType"
	codeNonRepeatableError         = "NonRepeatableError"
	codeRepeatableError            = "RepeatableError"
	codeInvalidURLParams           = "InvalidURLParams"
)

var (
	registerOnce sync.Once

	/* TODO: After deleting "gateway" service, we need to do the following things:
	1. Rename cacheCounter to grpccacheCounter
	2. Rename api_cache_requests_total to api_grpc_cache_requests_total
	3. Rename api_register_events_total to api_grpc_register_events_total
	4. Rename restCacheCounter to cacheCounter
	5. Rename api_rest_cache_requests_total to api_cache_requests_total
	6. Rename api_rest_register_events_total to api_register_events_total
	*/

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

	restCacheCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_rest_cache_requests_total",
			Help:      "Total number of cache requests",
		}, []string{"caller", "type", "layer", "code"})

	restEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "api_rest_register_events_total",
			Help:      "Total number of registered events",
		}, []string{"caller", "type", "code"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(cacheCounter, eventCounter)
	})
}
