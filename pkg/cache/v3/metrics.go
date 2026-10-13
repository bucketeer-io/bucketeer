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
	//
	// The generation is the writer pod's wall clock captured before it read
	// the database, so occasional rejections are the expected outcome of two
	// writers racing on the same environment (typically the batch cacher losing
	// to the refresher). A sustained non-zero rate for writer="refresher", or
	// for a single environment_id, means snapshots are being written out of
	// order: check for clock skew between subscriber and batch nodes or for a
	// slow database read. Worth an alert.
	featuresPutRejectedStaleCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "cache",
			Name:      "features_put_rejected_stale_total",
			Help: "Total number of feature snapshot writes rejected because a newer snapshot was already " +
				"cached. Occasional rejections are expected; a sustained rate indicates clock skew or " +
				"out-of-order writes.",
		}, []string{"writer", "environment_id"})
)

// RegisterMetrics registers the cache metrics. Safe to call more than once.
func RegisterMetrics(r metrics.Registerer) {
	registerMetricsOnce.Do(func() {
		r.MustRegister(featuresPutRejectedStaleCounter)
	})
}

// RecordFeaturesPutRejectedStale records a rejected stale snapshot write for
// the given writer and environment.
func RecordFeaturesPutRejectedStale(writer, environmentID string) {
	FeaturesPutRejectedStaleCounter(writer, environmentID).Inc()
}

// FeaturesPutRejectedStaleCounter returns the counter series for writer and
// environment. It exists so writers' tests can assert on the metric;
// production code should use RecordFeaturesPutRejectedStale.
func FeaturesPutRejectedStaleCounter(writer, environmentID string) prometheus.Counter {
	return featuresPutRejectedStaleCounter.WithLabelValues(writer, environmentID)
}
