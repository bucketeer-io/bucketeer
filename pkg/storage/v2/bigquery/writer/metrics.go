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

package writer

import (
	"time"

	"cloud.google.com/go/bigquery/storage/apiv1/storagepb"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	operationQuery              = "Query"
	codeOK                      = "OK"
	codeUnknown                 = "Unknown"
	storageErrorCodeUnspecified = "StorageErrorCodeUnspecified"
	tableNotFound               = "TableNotFound"
	streamAlreadyCommitted      = "StreamAlreadyCommitted"
	sreamNotFound               = "StreamNotFound"
	invalidStreamType           = "InvalidStreamType"
	invalidStreamState          = "InvalidStreamState"
	streamFinalized             = "streamFinalized"
	schemaMismatchExtraFields   = "SchemaMismatchExtraFields"
	offsetAlreadyExists         = "offsetAlreadyExists"
	offsetOutOfRange            = "offsetOutOfRange"
)

var (
	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "bigquery_writer",
			Name:      "handled_total",
			Help:      "Total number of completed operations.",
		}, []string{"operation", "code"})

	handledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "bigquery_writer",
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
			return codeUnknown
		}
		switch storageErr.GetCode() {
		case storagepb.StorageError_STORAGE_ERROR_CODE_UNSPECIFIED:
			return storageErrorCodeUnspecified
		case storagepb.StorageError_TABLE_NOT_FOUND:
			return tableNotFound
		case storagepb.StorageError_STREAM_ALREADY_COMMITTED:
			return streamAlreadyCommitted
		case storagepb.StorageError_STREAM_NOT_FOUND:
			return sreamNotFound
		case storagepb.StorageError_INVALID_STREAM_TYPE:
			return invalidStreamType
		case storagepb.StorageError_INVALID_STREAM_STATE:
			return invalidStreamState
		case storagepb.StorageError_STREAM_FINALIZED:
			return streamFinalized
		case storagepb.StorageError_SCHEMA_MISMATCH_EXTRA_FIELDS:
			return schemaMismatchExtraFields
		case storagepb.StorageError_OFFSET_ALREADY_EXISTS:
			return offsetAlreadyExists
		case storagepb.StorageError_OFFSET_OUT_OF_RANGE:
			return offsetOutOfRange
		}
	}
	return codeUnknown
}

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		handledCounter,
		handledHistogram,
	)
}
