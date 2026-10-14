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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterMetricsIsIdempotent(t *testing.T) {
	reg := prometheus.NewRegistry()
	RegisterMetrics(reg)
	// A second call must not panic with AlreadyRegisteredError.
	RegisterMetrics(reg)
	RegisterMetrics(prometheus.NewRegistry())
}

func TestRecordFeaturesPutRejectedStale(t *testing.T) {
	refresherEnvA := FeaturesPutRejectedStaleCounter(FeaturesWriterRefresher, "metrics-env-a")
	batchEnvA := FeaturesPutRejectedStaleCounter(FeaturesWriterBatchCacher, "metrics-env-a")
	batchEnvB := FeaturesPutRejectedStaleCounter(FeaturesWriterBatchCacher, "metrics-env-b")
	beforeRefresherA := testutil.ToFloat64(refresherEnvA)
	beforeBatchA := testutil.ToFloat64(batchEnvA)
	beforeBatchB := testutil.ToFloat64(batchEnvB)

	RecordFeaturesPutRejectedStale(FeaturesWriterRefresher, "metrics-env-a")
	RecordFeaturesPutRejectedStale(FeaturesWriterBatchCacher, "metrics-env-a")
	RecordFeaturesPutRejectedStale(FeaturesWriterBatchCacher, "metrics-env-a")
	RecordFeaturesPutRejectedStale(FeaturesWriterBatchCacher, "metrics-env-b")

	assert.Equal(t, float64(1), testutil.ToFloat64(refresherEnvA)-beforeRefresherA)
	assert.Equal(t, float64(2), testutil.ToFloat64(batchEnvA)-beforeBatchA)
	assert.Equal(t, float64(1), testutil.ToFloat64(batchEnvB)-beforeBatchB)

	// One series per (writer, environment_id) pair touched above. Other tests
	// in this package may add series, so assert at least these three exist.
	count, err := testutil.GatherAndCount(mustRegistry(t), "bucketeer_cache_features_put_rejected_stale_total")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 3, "one series per (writer, environment_id)")
}

func mustRegistry(t *testing.T) *prometheus.Registry {
	t.Helper()
	reg := prometheus.NewRegistry()
	require.NoError(t, reg.Register(featuresPutRejectedStaleCounter))
	return reg
}
