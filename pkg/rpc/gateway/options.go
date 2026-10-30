// Copyright 2026 The Bucketeer Authors.
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

package gateway

import (
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
)

type Option func(*options)

type options struct {
	metrics                   metrics.Registerer
	logger                    *zap.Logger
	timeout                   time.Duration
	keepaliveTime             time.Duration
	keepaliveTimeout          time.Duration
	keepaliveMinTime          time.Duration
	maxConnectionIdle         time.Duration
	maxConnectionAge          time.Duration
	maxConnectionAgeGrace     time.Duration
	time                      time.Duration
	maxConcurrentStreams      uint32
	maxConnectionIdleTime     time.Duration
	maxConnectionAgeTime      time.Duration
	maxConnectionAgeGraceTime time.Duration
	timeoutTime               time.Duration
	permitWithoutStream       bool
	initialWindowSize         int32
	initialConnWindowSize     int32
	certPath                  string
	keyPath                   string
	httpReadTimeout           time.Duration
	httpWriteTimeout          time.Duration
	httpIdleTimeout           time.Duration
	maxRequestBodyBytes       int64
	incomingHeaderMatcher     runtime.HeaderMatcherFunc
}

var defaultOptions = options{
	timeout:                   30 * time.Second,
	keepaliveTime:             30 * time.Second,
	keepaliveTimeout:          10 * time.Second,
	keepaliveMinTime:          30 * time.Second,
	maxConnectionIdle:         300 * time.Second,
	maxConnectionAge:          300 * time.Second,
	maxConnectionAgeGrace:     10 * time.Second,
	time:                      30 * time.Second,
	maxConcurrentStreams:      100,
	maxConnectionIdleTime:     300 * time.Second,
	maxConnectionAgeTime:      300 * time.Second,
	maxConnectionAgeGraceTime: 10 * time.Second,
	timeoutTime:               30 * time.Second,
	permitWithoutStream:       true,
	initialWindowSize:         1024 * 1024 * 2, // 2MB
	initialConnWindowSize:     1024 * 1024 * 2, // 2MB
	certPath:                  "",
	keyPath:                   "",
	httpReadTimeout:           30 * time.Second,
	httpWriteTimeout:          30 * time.Second,
	httpIdleTimeout:           60 * time.Second,
	// The gRPC servers behind the gateway accept messages of at most 4MB
	// (grpc-go default). The JSON form of a message is larger than its binary
	// form (field names, base64 for bytes fields), so twice that keeps the
	// gateway from rejecting a body the backend would still accept, while
	// bounding how much of a request body is buffered before it is decoded.
	maxRequestBodyBytes: 8 * 1024 * 1024, // 8MB
}

func WithMetrics(metrics metrics.Registerer) Option {
	return func(o *options) {
		o.metrics = metrics
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

func WithKeepaliveTime(keepaliveTime time.Duration) Option {
	return func(o *options) {
		o.keepaliveTime = keepaliveTime
	}
}

func WithKeepaliveTimeout(keepaliveTimeout time.Duration) Option {
	return func(o *options) {
		o.keepaliveTimeout = keepaliveTimeout
	}
}

func WithKeepaliveMinTime(keepaliveMinTime time.Duration) Option {
	return func(o *options) {
		o.keepaliveMinTime = keepaliveMinTime
	}
}

func WithMaxConnectionIdle(maxConnectionIdle time.Duration) Option {
	return func(o *options) {
		o.maxConnectionIdle = maxConnectionIdle
	}
}

func WithMaxConnectionAge(maxConnectionAge time.Duration) Option {
	return func(o *options) {
		o.maxConnectionAge = maxConnectionAge
	}
}

func WithMaxConnectionAgeGrace(maxConnectionAgeGrace time.Duration) Option {
	return func(o *options) {
		o.maxConnectionAgeGrace = maxConnectionAgeGrace
	}
}

func WithPermitWithoutStream(permitWithoutStream bool) Option {
	return func(o *options) {
		o.permitWithoutStream = permitWithoutStream
	}
}

func WithInitialWindowSize(initialWindowSize int32) Option {
	return func(o *options) {
		o.initialWindowSize = initialWindowSize
	}
}

func WithInitialConnWindowSize(initialConnWindowSize int32) Option {
	return func(o *options) {
		o.initialConnWindowSize = initialConnWindowSize
	}
}

func WithCertPath(certPath string) Option {
	return func(o *options) {
		o.certPath = certPath
	}
}

func WithKeyPath(keyPath string) Option {
	return func(o *options) {
		o.keyPath = keyPath
	}
}

func WithHTTPTimeouts(readTimeout, writeTimeout, idleTimeout time.Duration) Option {
	return func(o *options) {
		o.httpReadTimeout = readTimeout
		o.httpWriteTimeout = writeTimeout
		o.httpIdleTimeout = idleTimeout
	}
}

// WithMaxRequestBodyBytes sets the maximum number of bytes a handler reads
// from a request body. A request whose body exceeds the limit fails to
// decode and is rejected before it reaches the gRPC backend.
func WithMaxRequestBodyBytes(maxRequestBodyBytes int64) Option {
	return func(o *options) {
		o.maxRequestBodyBytes = maxRequestBodyBytes
	}
}

// WithIncomingHeaderMatcher sets the function that decides which HTTP request
// headers are forwarded to the gRPC services as metadata. When it is not set,
// grpc-gateway's runtime.DefaultHeaderMatcher is used.
func WithIncomingHeaderMatcher(matcher runtime.HeaderMatcherFunc) Option {
	return func(o *options) {
		o.incomingHeaderMatcher = matcher
	}
}
