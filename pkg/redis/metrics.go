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

package redis

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	CodeSuccess     = "Success"
	CodeFail        = "Fail"
	CodeNotFound    = "NotFound"
	CodeInvalidType = "InvalidType"
)

var (
	ReceivedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "redis",
			Name:      "received_total",
			Help:      "Total number of received commands.",
		}, []string{"version", "server", "command"})

	HandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "redis",
			Name:      "handled_total",
			Help:      "Total number of completed commands.",
		}, []string{"version", "server", "command", "code"})

	HandledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "redis",
			Name:      "handling_seconds",
			Help:      "Histogram of command response latency (seconds).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"version", "server", "command", "code"})

	poolActiveConnectionsDesc = prometheus.NewDesc(
		"bucketeer_redis_pool_active_connections",
		"Number of connections in the pool.",
		[]string{"version", "server"},
		nil,
	)

	poolIdleConnectionsDesc = prometheus.NewDesc(
		"bucketeer_redis_pool_idle_connections",
		"Number of idle connections in the pool.",
		[]string{"version", "server"},
		nil,
	)

	registerOnce sync.Once
	clients      sync.Map
)

type PoolStater interface {
	Stats() PoolStats
}

type PoolStats interface {
	ActiveCount() int
	IdleCount() int
}

type metricsKey struct {
	version string
	server  string
}

func RegisterMetrics(r metrics.Registerer, version, server string, stater PoolStater) {
	clients.Store(metricsKey{version: version, server: server}, stater)
	registerOnce.Do(func() {
		r.MustRegister(
			ReceivedCounter,
			HandledCounter,
			HandledHistogram,
			&poolCollector{},
		)
	})
}

type poolCollector struct {
}

func (c *poolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- poolActiveConnectionsDesc
	ch <- poolIdleConnectionsDesc
}

func (c *poolCollector) Collect(ch chan<- prometheus.Metric) {
	clients.Range(func(key, value interface{}) bool {
		mKey := key.(metricsKey)
		stater := value.(PoolStater)
		stats := stater.Stats()
		ch <- prometheus.MustNewConstMetric(
			poolActiveConnectionsDesc,
			prometheus.GaugeValue,
			float64(stats.ActiveCount()),
			mKey.version,
			mKey.server,
		)
		ch <- prometheus.MustNewConstMetric(
			poolIdleConnectionsDesc,
			prometheus.GaugeValue,
			float64(stats.IdleCount()),
			mKey.version,
			mKey.server,
		)
		return true
	})
}
