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

package bigtable

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	operationReadItems = "ReadItems"
	operationReadRows  = "ReadRows"
	operationWriteRow  = "WriteRow"
	operationWriteRows = "WriteRows"
	operationClose     = "Close"

	codeOK                       = "OK"
	codeKeyNotFound              = "KeyNotFound"
	codeColumnFamilyNotFound     = "ColumnFamilyNotFound"
	codeColumnNotFound           = "ColumnNotFound"
	codeFailedToWritePartialRows = "FailedToWritePartialRows"
	codeDeadlineExceeded         = "DeadlineExceeded"
	codeCanceled                 = "Canceled"
	codeUnknown                  = "Unknown"
)

var (
	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "bigtable",
			Name:      "handled_total",
			Help:      "Total number of completed operations.",
		}, []string{"operation", "code"})

	handledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "bigtable",
			Name:      "handling_seconds",
			Help:      "Histogram of operation response latency (seconds).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"operation", "code"})
)

func record() func(operation string, err *error) {
	startTime := time.Now()
	return func(operation string, err *error) {
		var code string
		switch *err {
		case nil:
			code = codeOK
		case ErrKeyNotFound:
			code = codeKeyNotFound
		case ErrColumnFamilyNotFound:
			code = codeColumnFamilyNotFound
		case ErrColumnNotFound:
			code = codeColumnNotFound
		case errFailedToWritePartialRows:
			code = codeFailedToWritePartialRows
		case context.DeadlineExceeded:
			code = codeDeadlineExceeded
		case context.Canceled:
			code = codeCanceled
		default:
			code = codeUnknown
		}
		handledCounter.WithLabelValues(operation, code).Inc()
		handledHistogram.WithLabelValues(operation, code).Observe(time.Since(startTime).Seconds())
	}
}

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		handledCounter,
		handledHistogram,
	)
}
