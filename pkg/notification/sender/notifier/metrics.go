// Copyright 2023 The Bucketeer Authors.
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

package notifier

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	typeSlack   = "Slack"
	codeSuccess = "Success"
	codeFail    = "Fail"
)

var (
	receivedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "notification",
			Name:      "notifier_received_total",
			Help:      "Total number of received messages",
		}, []string{"type"})

	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "notification",
			Name:      "notifier_handled_total",
			Help:      "Total number of handled messages",
		}, []string{"type", "code"})
)

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		receivedCounter,
		handledCounter,
	)
}
