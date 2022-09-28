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

package datastore

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	writerKafka = "Kafka"
	codeSuccess = "Success"
	codeFail    = "Fail"
)

var (
	registerOnce sync.Once

	writeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister",
			Name:      "write_total",
			Help:      "Total number of writes",
		}, []string{"writer", "code"})

	wroteHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "event_persister",
			Name:      "wrote_seconds",
			Help:      "Histogram of events handling duration (seconds)",
			Buckets:   prometheus.DefBuckets,
		}, []string{"writer", "code"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			writeCounter,
			wroteHistogram,
		)
	})
}
