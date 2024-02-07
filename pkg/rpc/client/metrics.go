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

package client

import (
	"context"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const rpcTypeUnary = "Unary"

var (
	registerOnce sync.Once

	startedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "client_started_total",
			Help:      "Total number of RPCs started on the client.",
		}, []string{"type", "service", "method"})

	handledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "grpc",
			Name:      "client_handled_total",
			Help:      "Total number of RPCs completed by the client, regardless of success or failure.",
		}, []string{"type", "service", "method", "code"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(startedCounter, handledCounter)
	})
}

func MetricsUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		serviceName, methodName := splitFullMethodName(method)
		startedCounter.WithLabelValues(rpcTypeUnary, serviceName, methodName).Inc()
		err := invoker(ctx, method, req, reply, cc, opts...)
		handledCounter.WithLabelValues(rpcTypeUnary, serviceName, methodName, status.Code(err).String()).Inc()
		return err
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
