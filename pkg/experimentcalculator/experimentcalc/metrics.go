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
//

package experimentcalc

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	calculationFail    = "Fail"
	calculationSuccess = "Success"

	binomialModelSampleMethod = "binomialModelSample"
	normalInverseGammaMethod  = "normalInverseGamma"

	valuesAreZero                    = "valuesAreZero"
	evalVariationCountNotFound       = "evalVariationCountNotFound"
	evaluationCountLessThanGoalEvent = "evaluationCountLessThanGoalEvent"
)

var (
	calculationCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "experiment_calculator",
			Name:      "calculate_calls_total",
			Help:      "Total number of calculate calls",
		}, []string{"code"})

	calculationExceptionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "experiment_calculator",
			Name:      "calculate_skipped_calls_total",
			Help:      "Total number of calculate skipped calls",
		}, []string{"exception"})

	calculationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "experiment_calculator",
			Name:      "calculate_duration_seconds",
			Help:      "The duration of calculate method (in seconds)",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method"})
)

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		calculationCounter,
		calculationExceptionCounter,
		calculationHistogram,
	)
}
