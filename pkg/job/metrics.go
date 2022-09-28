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

package job

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	codeSuccess = "Success"
	codeFail    = "Fail"
)

var (
	startedJobCounter    *prometheus.CounterVec
	finishedJobCounter   *prometheus.CounterVec
	finishedJobHistogram *prometheus.HistogramVec
)

func setSubsystem(subsystem string) {
	startedJobCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: subsystem,
			Name:      "batch_started_jobs_total",
			Help:      "Total number of started jobs.",
		}, []string{"name"})

	finishedJobCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: subsystem,
			Name:      "batch_finished_jobs_total",
			Help:      "Total number of finished jobs.",
		}, []string{"name", "code"})

	finishedJobHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: subsystem,
			Name:      "batch_job_running_time_seconds",
			Help:      "Histogram of the time job takes to run in seconds.",
		}, []string{"name", "code"})

}
func registerMetrics(r metrics.Registerer, subsystem string) {
	setSubsystem(subsystem)
	r.MustRegister(
		startedJobCounter,
		finishedJobCounter,
		finishedJobHistogram,
	)
}
