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

package rest

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

var (
	registerOnce sync.Once

	serverStartedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "rest",
			Name:      "server_started_total",
			Help:      "Total number of REST started on the server.",
		}, []string{"version", "service", "method"})

	serverHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "rest",
			Name:      "server_handled_total",
			Help:      "Total number of REST completed on the server, regardless of success or failure.",
		}, []string{"version", "service", "method", "code"})

	serverHandledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "rest",
			Name:      "server_handling_seconds",
			Help:      "Histogram of response latency (seconds) of REST that had been application-level handled by the server.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"version", "service", "method"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			serverStartedCounter,
			serverHandledCounter,
			serverHandledHistogram,
		)
	})
}

func MetricsServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			rr := &responseRecorder{ResponseWriter: w, body: new(bytes.Buffer)}
			apiVersion, serviceName, apiName := splitURLPath(r.URL.Path)
			serverStartedCounter.WithLabelValues(apiVersion, serviceName, apiName).Inc()
			next.ServeHTTP(rr, r)
			serverHandledCounter.WithLabelValues(apiVersion, serviceName, apiName, strconv.Itoa(rr.statusCode)).Inc()
			serverHandledHistogram.WithLabelValues(apiVersion, serviceName, apiName).Observe(time.Since(startTime).Seconds())
		},
	)
}
