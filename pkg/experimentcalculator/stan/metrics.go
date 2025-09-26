// Copyright 2025 The Bucketeer Authors.
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

package stan

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

const (
	// HTTPStan operations
	methodCompileModel        = "compile_model"
	methodCreateFit           = "create_fit"
	methodGetOperationDetails = "get_operation_details"
	methodGetFitResult        = "get_fit_result"
	methodStanParams          = "stan_params"

	// Status codes
	codeOK         = "OK"
	codeNotFound   = "NotFound"
	codeBadRequest = "BadRequest"
	codeTimeout    = "Timeout"
	codeInternal   = "Internal"
)

var (
	registerOnce           sync.Once
	httpstanRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "httpstan_client",
			Name:      "requests_total",
			Help:      "Total number of httpstan requests",
		}, []string{"method", "code"})

	httpstanRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "httpstan_client",
			Name:      "request_duration_seconds",
			Help:      "Duration of httpstan requests",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method"})
)

// RegisterMetrics registers HTTPStan client metrics
func RegisterMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			httpstanRequestCounter,
			httpstanRequestDuration,
		)
	})
}

// RecordRequest records an HTTPStan request
func RecordRequest(method string, code string, duration time.Duration) {
	httpstanRequestCounter.WithLabelValues(method, code).Inc()
	httpstanRequestDuration.WithLabelValues(method).Observe(duration.Seconds())
}

// RecordHTTPStan is a helper that records HTTPStan request with error handling
func RecordHTTPStan(method string, err error, statusCode int, duration time.Duration) {
	var code string
	if err != nil {
		code = codeInternal
	} else {
		code = getStatusCode(statusCode)
	}
	RecordRequest(method, code, duration)
}

// getStatusCode maps HTTP status codes to metric codes
func getStatusCode(httpStatus int) string {
	switch httpStatus {
	case http.StatusOK, http.StatusCreated:
		return codeOK
	case http.StatusBadRequest:
		return codeBadRequest
	case http.StatusNotFound:
		return codeNotFound
	case http.StatusRequestTimeout:
		return codeTimeout
	default:
		return codeInternal
	}
}
