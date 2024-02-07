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
		code := getCodeFromError(*err)
		handledCounter.WithLabelValues(operation, code).Inc()
		handledHistogram.WithLabelValues(operation, code).Observe(time.Since(startTime).Seconds())
	}
}

func getCodeFromError(err error) string {
	if err == nil {
		return codeOK
	}
	var e *googleapi.Error
	if ok := errors.As(err, &e); !ok {
		return codeUnknown
	}
	switch e.Code {
	case http.StatusBadRequest:
		return codeBadRequest
	case http.StatusForbidden:
		return codeForbidden
	case http.StatusNotFound:
		return codeNotFound
	case http.StatusConflict:
		return codeConflict
	case http.StatusInternalServerError:
		return codeInternalServerError
	case http.StatusNotImplemented:
		return codeNotImplemented
	case http.StatusServiceUnavailable:
		return codeServiceUnavailable
	default:
		return codeUnknown
	}
}

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		handledCounter,
		handledHistogram,
	)
}
