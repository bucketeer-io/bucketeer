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

package api

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	codeFail    = "Fail"
	codeSuccess = "Success"
)

var (
	listExperimentCountsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_counter",
			Name:      "api_list_experiment_counts_calls_total",
			Help:      "Total number of ListExperimentCounts calls",
		}, []string{"code"})

	listExperimentResultsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_counter",
			Name:      "api_list_experiment_results_calls_total",
			Help:      "Total number of ListExperimentResults calls",
		}, []string{"code"})

	getExperimentCountsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_counter",
			Name:      "get_experiment_counts_calls_total",
			Help:      "Total number of GetExperimentCounts calls",
		}, []string{"code"})

	getExperimentResultCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "event_counter",
			Name:      "get_experiment_result_calls_total",
			Help:      "Total number of GetExperimentResult calls",
		}, []string{"code"})
)

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		listExperimentCountsCounter,
		listExperimentResultsCounter,
		getExperimentCountsCounter,
		getExperimentResultCounter,
	)
}
