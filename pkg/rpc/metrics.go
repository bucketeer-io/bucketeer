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
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
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

	// Shutdown-related metrics
	shutdownStartedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "server",
			Name:      "shutdown_started_total",
			Help:      "Total number of server shutdowns initiated.",
		}, []string{"server_name", "shutdown_reason"})

	shutdownDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "server",
			Name:      "shutdown_duration_seconds",
			Help:      "Time taken for graceful shutdown by component.",
			Buckets:   []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 15.0, 20.0, 25.0, 30.0},
		}, []string{"server_name", "component", "status"})

	inflightRequestsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "bucketeer",
			Subsystem: "server",
			Name:      "inflight_requests",
			Help:      "Number of requests currently being processed.",
		}, []string{"server_name", "protocol"})

	shutdownRequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "server",
			Name:      "shutdown_requests_total",
			Help:      "Total number of requests processed during shutdown phase.",
		}, []string{"server_name", "protocol", "status"})

	shutdownRequestDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "server",
			Name:      "shutdown_request_duration_seconds",
			Help:      "Duration of requests processed during shutdown phase.",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
		}, []string{"server_name", "protocol", "status"})

	droppedRequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "server",
			Name:      "dropped_requests_total",
			Help:      "Total number of requests dropped during shutdown.",
		}, []string{"server_name", "protocol", "reason"})
)

func registerMetrics(r metrics.Registerer) {
	registerOnce.Do(func() {
		r.MustRegister(
			serverStartedCounter,
			serverHandledCounter,
			serverHandledHistogram,
			// Shutdown metrics
			shutdownStartedCounter,
			shutdownDurationHistogram,
			inflightRequestsGauge,
			shutdownRequestsCounter,
			shutdownRequestDurationHistogram,
			droppedRequestsCounter,
		)
	})
}

// ShutdownTracker tracks shutdown metrics for a server
type ShutdownTracker struct {
	serverName   string
	logger       *zap.Logger
	isShutdown   int64 // atomic boolean
	shutdownTime time.Time
}

// NewShutdownTracker creates a new shutdown tracker
func NewShutdownTracker(serverName string, logger *zap.Logger) *ShutdownTracker {
	return &ShutdownTracker{
		serverName: serverName,
		logger:     logger,
	}
}

// IsShuttingDown returns true if the server is in shutdown phase
func (st *ShutdownTracker) IsShuttingDown() bool {
	return atomic.LoadInt64(&st.isShutdown) == 1
}

// StartShutdown marks the beginning of shutdown phase
func (st *ShutdownTracker) StartShutdown(reason string) {
	if atomic.CompareAndSwapInt64(&st.isShutdown, 0, 1) {
		st.shutdownTime = time.Now()
		shutdownStartedCounter.WithLabelValues(st.serverName, reason).Inc()
		st.logger.Info("Shutdown tracking started",
			zap.String("server", st.serverName),
			zap.String("reason", reason),
			zap.Time("shutdown_time", st.shutdownTime))
	}
}

// TrackShutdownDuration records the duration of a shutdown component
func (st *ShutdownTracker) TrackShutdownDuration(component string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "timeout"
	}
	shutdownDurationHistogram.WithLabelValues(st.serverName, component, status).Observe(duration.Seconds())

	st.logger.Info("Shutdown component completed",
		zap.String("server", st.serverName),
		zap.String("component", component),
		zap.Duration("duration", duration),
		zap.String("status", status))
}

// TrackInflightRequests updates the count of in-flight requests
func (st *ShutdownTracker) TrackInflightRequests(protocol string, delta int) {
	inflightRequestsGauge.WithLabelValues(st.serverName, protocol).Add(float64(delta))
}

// TrackShutdownRequest records a request processed during shutdown
func (st *ShutdownTracker) TrackShutdownRequest(protocol string, duration time.Duration, success bool) {
	status := "completed"
	if !success {
		status = "failed"
	}

	shutdownRequestsCounter.WithLabelValues(st.serverName, protocol, status).Inc()
	shutdownRequestDurationHistogram.WithLabelValues(st.serverName, protocol, status).Observe(duration.Seconds())

	if st.IsShuttingDown() {
		st.logger.Debug("Request processed during shutdown",
			zap.String("server", st.serverName),
			zap.String("protocol", protocol),
			zap.Duration("duration", duration),
			zap.String("status", status),
			zap.Duration("time_since_shutdown", time.Since(st.shutdownTime)))
	}
}

// TrackDroppedRequest records a request that was dropped
func (st *ShutdownTracker) TrackDroppedRequest(protocol string, reason string) {
	droppedRequestsCounter.WithLabelValues(st.serverName, protocol, reason).Inc()

	st.logger.Warn("Request dropped during shutdown",
		zap.String("server", st.serverName),
		zap.String("protocol", protocol),
		zap.String("reason", reason),
		zap.Duration("time_since_shutdown", time.Since(st.shutdownTime)))
}

// ShutdownAwareUnaryInterceptor creates a gRPC interceptor that tracks shutdown metrics
func (st *ShutdownTracker) ShutdownAwareUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		// Track in-flight request
		st.TrackInflightRequests("grpc", 1)
		defer st.TrackInflightRequests("grpc", -1)

		// Check if we're shutting down when request starts
		isShutdownRequest := st.IsShuttingDown()

		// Process the request
		resp, err := handler(ctx, req)

		// Track the request if it was processed during shutdown
		if isShutdownRequest {
			duration := time.Since(startTime)
			success := err == nil || status.Code(err).String() != "Unavailable"
			st.TrackShutdownRequest("grpc", duration, success)
		}

		return resp, err
	}
}

// ShutdownAwareHTTPMiddleware creates HTTP middleware that tracks shutdown metrics
func (st *ShutdownTracker) ShutdownAwareHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Track in-flight request
		st.TrackInflightRequests("http", 1)
		defer st.TrackInflightRequests("http", -1)

		// Check if we're shutting down when request starts
		isShutdownRequest := st.IsShuttingDown()

		// Check if server is shutting down and reject new requests
		if isShutdownRequest {
			// Allow health check requests during shutdown for graceful coordination
			if r.URL.Path == "/health" || r.URL.Path == "/internal/shutdown-ready" {
				next.ServeHTTP(w, r)
				return
			}

			// For other requests during shutdown, we can either:
			// 1. Process them (current behavior)
			// 2. Reject them immediately (optional)

			// Option 1: Process the request and track it
			responseRecorder := &responseRecorder{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(responseRecorder, r)

			duration := time.Since(startTime)
			success := responseRecorder.statusCode < 400
			st.TrackShutdownRequest("http", duration, success)

			return
		}

		// Normal request processing
		next.ServeHTTP(w, r)
	})
}

// responseRecorder captures the HTTP status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
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
