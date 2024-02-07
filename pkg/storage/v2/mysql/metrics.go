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

package mysql

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	operationExec             = "Exec"
	operationQuery            = "Query"
	operationQueryRow         = "QueryRow"
	operationBeginTx          = "BeginTx"
	operationRunInTransaction = "RunInTransaction"
	operationCommit           = "Commit"
	operationRollback         = "Rollback"

	codeOK             = "OK"
	codeNoRows         = "NoRows"
	codeTxDone         = "TxDone"
	codeDuplicateEntry = "DuplicateEntry"
	codeUnknown        = "Unknown"
)

var (
	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "mysql",
			Name:      "handled_total",
			Help:      "Total number of completed operations.",
		}, []string{"operation", "code"})

	handledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "mysql",
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
		case ErrNoRows:
			code = codeNoRows
		case ErrTxDone:
			code = codeTxDone
		case ErrDuplicateEntry:
			code = codeDuplicateEntry
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
