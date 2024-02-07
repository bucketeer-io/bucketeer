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

package rpc

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
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
		}, []string{"type", "service", "method"})

	serverHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "server_handled_total",
			Help:      "Total number of RPCs completed on the server, regardless of success or failure.",
		}, []string{"type", "service", "method", "code"})

	serverHandledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "server_handling_seconds",
			Help:      "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"type", "service", "method"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			serverStartedCounter,
			serverHandledCounter,
			serverHandledHistogram,
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
		serverStartedCounter.WithLabelValues(rpcTypeUnary, serviceName, methodName).Inc()
		resp, err := handler(ctx, req)
		serverHandledCounter.WithLabelValues(rpcTypeUnary, serviceName, methodName, status.Code(err).String()).Inc()
		serverHandledHistogram.WithLabelValues(rpcTypeUnary, serviceName, methodName).Observe(time.Since(startTime).Seconds())
		return resp, err
	}
}

func splitFullMethodName(fullMethodName string) (string, string) {
	// format: /package.service/method
	parts := strings.Split(fullMethodName, "/")
	if len(parts) != 3 {
		return "unknown", "unknown"
	}
	return parts[1], parts[2]
}
