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
)

var (
	sseActiveConnectionsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_active_connections",
			Help:      "Current number of active SSE connections held in the stream dispatcher.",
		}, []string{"environment_id", "tag"})
	sseDispatchDroppedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "gateway",
			Name:      "sse_dispatch_dropped_total",
			Help:      "Total dispatched events dropped because a connection's buffer was full.",
		}, []string{"environment_id", "tag"})
)

func RegisterMetrics(r metrics.Registerer) {
	r.MustRegister(
		sseActiveConnectionsGauge,
		sseDispatchDroppedCounter,
	)
}
