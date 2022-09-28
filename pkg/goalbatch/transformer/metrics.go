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

package transformer

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	typeGoal = "Goal"

	codeOK   = "OK"
	codeFail = "Fail"
)

var (
	receivedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "goal_batch_event_transformer",
			Name:      "received_total",
			Help:      "Total number of received messages",
		},
	)

	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "goal_batch_event_transformer",
			Name:      "handled_total",
			Help:      "Total number of handled messages",
		}, []string{"code"},
	)

	handledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "goal_batch_event_transformer",
			Name:      "handled_seconds",
			Help:      "Histogram of message handling duration (seconds)",
			Buckets:   prometheus.DefBuckets,
		}, []string{"code"})

	cacheCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "goal_batch_event_transformer",
			Name:      "cache_requests_total",
			Help:      "Total number of cache requests",
		}, []string{"type", "code"})

	eventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "goal_batch_event_transformer",
			Name:      "register_events_total",
			Help:      "Total number of registered events",
		}, []string{"type", "code"})
)

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		receivedCounter,
		handledCounter,
		handledHistogram,
		cacheCounter,
		eventCounter,
	)
}
