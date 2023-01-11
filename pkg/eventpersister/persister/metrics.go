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

package persister

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	codeFailedToExtractOpsEventRateClauses = "FailedToExtractOpsEventRateClauses"
	codeFailedToGetFeatures                = "FailedToGetFeatures"
	codeFailedToGetUserEvaluation          = "FailedToGetUserEvaluation"
	codeFailedToListAutoOpsRules           = "FailedToListAutoOpsRules"
	codeFailedToListExperiments            = "FailedToListExperiments"
	codeNothingToLink                      = "NothingToLink"
	codeUpsertUserEvaluationFailed         = "UpsertUserEvaluationFailed"
	codeUserEvaluationNotFound             = "UserEvaluationNotFound"
)

var (
	receivedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister",
			Name:      "received_total",
			Help:      "Total number of received messages",
		})

	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister",
			Name:      "handled_total",
			Help:      "Total number of handled messages",
		}, []string{"code"})

	cacheCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister",
			Name:      "cache_requests_total",
			Help:      "Total number of cache requests",
		}, []string{"type", "code"})
	dwhReceivedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister_dwh",
			Name:      "received_total",
			Help:      "Total number of received messages",
		})

	dwhHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister_dwh",
			Name:      "handled_total",
			Help:      "Total number of handled messages",
		}, []string{"code"})

	dwhCacheCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister_dwh",
			Name:      "cache_requests_total",
			Help:      "Total number of cache requests",
		}, []string{"type", "code"})
)

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(receivedCounter, handledCounter, cacheCounter)
}

func dwhRegisterMetrics(r metrics.Registerer) {
	r.MustRegister(dwhReceivedCounter, dwhHandledCounter, dwhCacheCounter)
}
