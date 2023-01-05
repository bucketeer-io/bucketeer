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

package querier

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/api/googleapi"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	operationQuery = "Query"

	codeOK                  = "OK"
	codeBadRequest          = "BadRequest"
	codeForbidden           = "Forbidden"
	codeNotFound            = "NotFound"
	codeConflict            = "Conflict"
	codeInternalServerError = "InternalServerError"
	codeNotImplemented      = "NotImplemented"
	codeServiceUnavailable  = "ServiceUnavailable"
	codeUnknown             = "Unknown"
)

var (
	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "bigquery_querier",
			Name:      "handled_total",
			Help:      "Total number of completed operations.",
		}, []string{"operation", "code"})

	handledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "bigquery_querier",
			Name:      "handling_seconds",
			Help:      "Histogram of operation response latency (seconds).",
			Buckets:   prometheus.DefBuckets,
		}, []string{"operation", "code"})
)

func record() func(operation string, err *error) {
	startTime := time.Now()
	return func(operation string, err *error) {
		if err == nil {
			handledCounter.WithLabelValues(operation, codeOK).Inc()
			handledHistogram.WithLabelValues(operation, codeOK).Observe(time.Since(startTime).Seconds())
			return
		}
		var code string
		var e *googleapi.Error
		if ok := errors.As(*err, &e); !ok {
			handledCounter.WithLabelValues(operation, codeUnknown).Inc()
			handledHistogram.WithLabelValues(operation, codeUnknown).Observe(time.Since(startTime).Seconds())
			return
		}
		switch e.Code {
		case http.StatusBadRequest:
			code = codeBadRequest
		case http.StatusForbidden:
			code = codeForbidden
		case http.StatusNotFound:
			code = codeNotFound
		case http.StatusConflict:
			code = codeConflict
		case http.StatusInternalServerError:
			code = codeInternalServerError
		case http.StatusNotImplemented:
			code = codeNotImplemented
		case http.StatusServiceUnavailable:
			code = codeServiceUnavailable
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
