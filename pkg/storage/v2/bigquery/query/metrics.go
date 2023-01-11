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

package query

import (
	"time"

	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/prometheus/client_golang/prometheus"
	storagepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1"

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
			Subsystem: "bigquery_query",
			Name:      "handled_total",
			Help:      "Total number of completed operations.",
		}, []string{"operation", "code"})

	handledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "bigquery_query",
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
	if apiErr, ok := apierror.FromError(err); ok {
		storageErr := &storagepb.StorageError{}
		if e := apiErr.Details().ExtractProtoMessage(storageErr); e != nil {
		}
		switch storageErr.GetCode() {
		case storagepb.StorageError_STORAGE_ERROR_CODE_UNSPECIFIED:
			return codeBadRequest
		case storagepb.StorageError_TABLE_NOT_FOUND:
			return codeForbidden
		case storagepb.StorageError_STREAM_ALREADY_COMMITTED:
			return codeNotFound
		case storagepb.StorageError_STREAM_NOT_FOUND:
			return codeConflict
		case storagepb.StorageError_INVALID_STREAM_TYPE:
			return codeInternalServerError
		case storagepb.StorageError_INVALID_STREAM_STATE:
			return codeNotImplemented
		case storagepb.StorageError_STREAM_FINALIZED:
			return codeServiceUnavailable
		case storagepb.StorageError_SCHEMA_MISMATCH_EXTRA_FIELDS:
			return codeNotImplemented
		case storagepb.StorageError_OFFSET_ALREADY_EXISTS:
			return codeServiceUnavailable
		case storagepb.StorageError_OFFSET_OUT_OF_RANGE:
			return codeServiceUnavailable
		default:
			return codeUnknown
		}
	}
	return ""
}

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		handledCounter,
		handledHistogram,
	)
}
