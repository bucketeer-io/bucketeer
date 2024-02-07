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

package publisher

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	methodPublish      = "Publish"
	methodPublishMulti = "PublishMulti"

	codeOK               = "OK"
	codeBadMessage       = "BadMessage"
	codeDeadlineExceeded = "DeadlineExceeded"
	codeCanceled         = "Canceled"
	codeUnknown          = "Unknown"
)

var (
	registerOnce sync.Once

	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "pubsub_publisher",
			Name:      "handled_total",
			Help:      "Total number of handled messages",
		}, []string{"topic", "method", "code"},
	)

	handledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "pubsub_publisher",
			Name:      "handled_seconds",
			Help:      "Histogram of message handling duration (seconds)",
			Buckets:   prometheus.DefBuckets,
		}, []string{"topic", "method", "code"})
)

func convertErrorToCode(err error) string {
	switch err {
	case nil:
		return codeOK
	case ErrBadMessage:
		return codeBadMessage
	case context.DeadlineExceeded:
		return codeDeadlineExceeded
	case context.Canceled:
		return codeCanceled
	default:
		return codeUnknown
	}
}

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			handledCounter,
			handledHistogram,
		)
	})
}
