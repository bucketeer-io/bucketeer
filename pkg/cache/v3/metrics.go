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

package v3

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

const (
	// FeaturesWriterRefresher labels writes from the subscriber cache refresher.
	FeaturesWriterRefresher = "refresher"
	// FeaturesWriterBatchCacher labels writes from the batch feature flag cacher.
	FeaturesWriterBatchCacher = "batch_cacher"
)

var (
	registerMetricsOnce sync.Once

	// featuresPutRejectedStaleCounter counts snapshot writes rejected by
	// FeaturesCache.PutIfNewer because a newer snapshot was already cached.
	// Rejections from the batch cacher are expected; sustained rejections from
	// the refresher indicate an ordering or clock problem.
	featuresPutRejectedStaleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "cache",
			Name:      "features_put_rejected_stale_total",
			Help:      "Total number of feature snapshot writes rejected because a newer snapshot was already cached",
		}, []string{"writer"})
)

// RegisterMetrics registers the cache metrics. Safe to call more than once.
func RegisterMetrics(r metrics.Registerer) {
	registerMetricsOnce.Do(func() {
		r.MustRegister(featuresPutRejectedStaleCounter)
	})
}

// RecordFeaturesPutRejectedStale records a rejected stale snapshot write.
func RecordFeaturesPutRejectedStale(writer string) {
	FeaturesPutRejectedStaleCounter(writer).Inc()
}

// FeaturesPutRejectedStaleCounter returns the counter series for writer. It
// exists so writers' tests can assert on the metric; production code should
// use RecordFeaturesPutRejectedStale.
func FeaturesPutRejectedStaleCounter(writer string) prometheus.Counter {
	return featuresPutRejectedStaleCounter.WithLabelValues(writer)
}
