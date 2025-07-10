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

package rpc

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/proto/gateway"
)

const rpcTypeUnary = "Unary"

var (
	registerOnce sync.Once

	serverStartedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "server_started_total",
			Help:      "Total number of RPCs started on the server.",
		}, []string{"type", "service", "method", "source_id", "sdk_version", "tag"})

	serverHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "server_handled_total",
			Help:      "Total number of RPCs completed on the server, regardless of success or failure.",
		}, []string{"type", "service", "method", "code", "source_id", "sdk_version", "tag"})

	serverHandledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "server_handling_seconds",
			Help:      "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"type", "service", "method", "source_id", "sdk_version", "tag"})

	serverResponseSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "server_response_size_bytes",
			Help:      "Histogram of response size (bytes) of gRPC.",
			Buckets:   prometheus.ExponentialBuckets(1024, 4, 6),
		}, []string{"type", "service", "method", "source_id", "sdk_version", "tag"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			serverStartedCounter,
			serverHandledCounter,
			serverHandledHistogram,
			serverResponseSizeHistogram,
		)
	})
}

func MetricsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		serviceName, methodName := splitFullMethodName(info.FullMethod)
		sourceID, sdkVersion, tag := extractRequestLabels(methodName, req)
		serverStartedCounter.WithLabelValues(rpcTypeUnary, serviceName, methodName, sourceID, sdkVersion, tag).Inc()
		resp, err := handler(ctx, req)
		if err == nil {
			if m, ok := resp.(proto.Message); ok {
				serverResponseSizeHistogram.WithLabelValues(
					rpcTypeUnary,
					serviceName,
					methodName,
					sourceID,
					sdkVersion,
					tag,
				).Observe(float64(proto.Size(m)))
			}
		}
		serverHandledCounter.WithLabelValues(
			rpcTypeUnary,
			serviceName,
			methodName,
			status.Code(err).String(),
			sourceID,
			sdkVersion,
			tag,
		).Inc()
		serverHandledHistogram.WithLabelValues(
			rpcTypeUnary,
			serviceName,
			methodName,
			sourceID,
			sdkVersion,
			tag,
		).Observe(time.Since(startTime).Seconds())
		return resp, err
	}
}

func extractRequestLabels(methodName string, req interface{}) (string, string, string) {
	switch methodName {
	case "GetEvaluations":
		if r, ok := req.(*gateway.GetEvaluationsRequest); ok {
			return r.SourceId.String(), r.SdkVersion, r.Tag
		}
	case "GetEvaluation":
		if r, ok := req.(*gateway.GetEvaluationRequest); ok {
			return r.SourceId.String(), r.SdkVersion, r.Tag
		}
	case "GetFeatureFlags":
		if r, ok := req.(*gateway.GetFeatureFlagsRequest); ok {
			return r.SourceId.String(), r.SdkVersion, ""
		}
	case "GetSegmentUsers":
		if r, ok := req.(*gateway.GetSegmentUsersRequest); ok {
			return r.SourceId.String(), r.SdkVersion, ""
		}
	case "RegisterEvents":
		if r, ok := req.(*gateway.RegisterEventsRequest); ok {
			return r.SourceId.String(), r.SdkVersion, ""
		}
	}
	return "", "", ""
}

func splitFullMethodName(fullMethodName string) (string, string) {
	// format: /package.service/method
	parts := strings.Split(fullMethodName, "/")
	if len(parts) != 3 {
		return "unknown", "unknown"
	}
	return parts[1], parts[2]
}
