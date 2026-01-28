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
	// Code values for metrics labels
	CodeSuccess = "Success"
	CodeFail    = "Fail"
)

var (
	registerOnce sync.Once

	// listFeaturesCounter tracks DB fetch operations
	listFeaturesCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "feature_flag_cacher",
			Name:      "list_features_total",
			Help:      "Total number of list features operations from DB",
		}, []string{"environment_id", "code"})

	listFeaturesDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "feature_flag_cacher",
			Name:      "list_features_duration_seconds",
			Help:      "Duration of list features operations from DB in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"environment_id"})

	// cachePutCounter tracks Redis put operations per environment
	cachePutCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "feature_flag_cacher",
			Name:      "cache_put_total",
			Help:      "Total number of cache put operations to Redis",
		}, []string{"environment_id", "code"})

	// featuresUpdatedGauge tracks the number of features in the last successful cache update
	featuresUpdatedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "bucketeer",
			Subsystem: "feature_flag_cacher",
			Name:      "features_updated",
			Help:      "Number of features in the last successful cache update",
		}, []string{"environment_id"})
)

// RegisterMetrics registers the feature flag cacher metrics.
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

// recordListFeatures records a list features operation from DB.
func recordListFeatures(environmentID, code string, durationSeconds float64) {
	listFeaturesCounter.WithLabelValues(environmentID, code).Inc()
	listFeaturesDuration.WithLabelValues(environmentID).Observe(durationSeconds)
}

// recordCachePut records a cache put operation to Redis.
func recordCachePut(environmentID, code string) {
	cachePutCounter.WithLabelValues(environmentID, code).Inc()
}

// recordFeaturesUpdated records the number of features in the last successful cache update.
func recordFeaturesUpdated(environmentID string, count int) {
	featuresUpdatedGauge.WithLabelValues(environmentID).Set(float64(count))
}
