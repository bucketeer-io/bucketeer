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

package gateway

import (
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
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
	perRPCCredentials         credentials.PerRPCCredentials
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
	perRPCCredentials:         nil,
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

func WithPerRPCCredentials(creds credentials.PerRPCCredentials) Option {
	return func(o *options) {
		o.perRPCCredentials = creds
	}
}
