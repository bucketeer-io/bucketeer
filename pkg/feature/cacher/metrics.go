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

package cacher

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

const (
	// code values for metrics labels
	codeSuccess = "Success"
	codeFail    = "Fail"

	// scope values for list features operations
	// scopeBatch is used for RefreshAllEnvironmentCaches (batch job fetching all environments)
	scopeBatch = "batch"
	// scopeSingle is used for RefreshEnvironmentCache (single environment refresh)
	scopeSingle = "single"

	// environmentIDAll is used for batch operations that cover all environments
	environmentIDAll = "all"
	// environmentIDProduction is used as a fallback for empty environment IDs
	// TODO: Remove this after the empty environment ID migration is complete
	environmentIDProduction = "production"

	// cacherTypeFeatureFlag is the cacher type for feature flags
	cacherTypeFeatureFlag = "feature_flag"
	// cacherTypeSegmentUser is the cacher type for segment users
	cacherTypeSegmentUser = "segment_user"
)

var (
	registerOnce sync.Once

	// listFeaturesCounter tracks DB fetch operations
	// cacher: "feature_flag" or "segment_user"
	// scope: "batch" for all-environments fetch, "single" for per-environment fetch
	// environment_id: "all" for batch, actual ID for single
	listFeaturesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "cacher",
			Name:      "list_features_total",
			Help:      "Total number of list features operations from DB",
		}, []string{"cacher", "scope", "environment_id", "code"})

	listFeaturesDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "cacher",
			Name:      "list_features_duration_seconds",
			Help:      "Duration of list features operations from DB in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"cacher", "scope", "environment_id"})

	// cachePutCounter tracks Redis put operations per environment
	cachePutCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "cacher",
			Name:      "cache_put_total",
			Help:      "Total number of cache put operations to Redis",
		}, []string{"cacher", "environment_id", "code"})

	// featuresUpdatedGauge tracks the number of features in the last successful cache update
	featuresUpdatedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "bucketeer",
			Subsystem: "cacher",
			Name:      "features_updated",
			Help:      "Number of features in the last successful cache update",
		}, []string{"cacher", "environment_id"})
)

// RegisterMetrics registers the cacher metrics.
func RegisterMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			listFeaturesCounter,
			listFeaturesDuration,
			cachePutCounter,
			featuresUpdatedGauge,
		)
	})
}

// normalizeEnvironmentID returns the environment ID or a fallback for empty values.
// TODO: Remove this function after the empty environment ID migration is complete
func normalizeEnvironmentID(environmentID string) string {
	if environmentID == "" {
		return environmentIDProduction
	}
	return environmentID
}

// recordListFeatures records a list features operation from DB.
// cacherType: cacherTypeFeatureFlag or cacherTypeSegment
// scope: scopeBatch for all-environments fetch, scopeSingle for per-environment fetch
// environmentID: environmentIDAll for batch, actual environment ID for single
func recordListFeatures(cacherType, scope, environmentID, code string, durationSeconds float64) {
	envID := normalizeEnvironmentID(environmentID)
	listFeaturesCounter.WithLabelValues(cacherType, scope, envID, code).Inc()
	listFeaturesDuration.WithLabelValues(cacherType, scope, envID).Observe(durationSeconds)
}

// recordCachePut records a cache put operation to Redis.
func recordCachePut(cacherType, environmentID, code string) {
	envID := normalizeEnvironmentID(environmentID)
	cachePutCounter.WithLabelValues(cacherType, envID, code).Inc()
}

// recordFeaturesUpdated records the number of features in the last successful cache update.
func recordFeaturesUpdated(cacherType, environmentID string, count int) {
	envID := normalizeEnvironmentID(environmentID)
	featuresUpdatedGauge.WithLabelValues(cacherType, envID).Set(float64(count))
}
