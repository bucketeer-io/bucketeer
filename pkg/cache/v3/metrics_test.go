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
	beforeRefresher := testutil.ToFloat64(FeaturesPutRejectedStaleCounter(FeaturesWriterRefresher))
	beforeBatch := testutil.ToFloat64(FeaturesPutRejectedStaleCounter(FeaturesWriterBatchCacher))

	RecordFeaturesPutRejectedStale(FeaturesWriterRefresher)
	RecordFeaturesPutRejectedStale(FeaturesWriterBatchCacher)
	RecordFeaturesPutRejectedStale(FeaturesWriterBatchCacher)

	afterRefresher := testutil.ToFloat64(FeaturesPutRejectedStaleCounter(FeaturesWriterRefresher))
	afterBatch := testutil.ToFloat64(FeaturesPutRejectedStaleCounter(FeaturesWriterBatchCacher))
	assert.Equal(t, float64(1), afterRefresher-beforeRefresher)
	assert.Equal(t, float64(2), afterBatch-beforeBatch)

	count, err := testutil.GatherAndCount(mustRegistry(t), "bucketeer_cache_features_put_rejected_stale_total")
	require.NoError(t, err)
	assert.Equal(t, 2, count, "one series per writer label")
}

func mustRegistry(t *testing.T) *prometheus.Registry {
	t.Helper()
	reg := prometheus.NewRegistry()
	require.NoError(t, reg.Register(featuresPutRejectedStaleCounter))
	return reg
}
