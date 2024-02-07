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
	"time"

	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/stats"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

type options struct {
	dialTimeout       time.Duration
	perRPCCredentials credentials.PerRPCCredentials
	block             bool
	logger            *zap.Logger
	metrics           metrics.Registerer
	statsHandler      stats.Handler
}

var defaultOptions = options{
	block:        false,
	logger:       zap.NewNop(),
	statsHandler: &ocgrpc.ClientHandler{},
}

type Option func(*options)

func WithDialTimeout(d time.Duration) Option {
	return func(o *options) {
		o.dialTimeout = d
	}
}

func WithPerRPCCredentials(creds credentials.PerRPCCredentials) Option {
	return func(o *options) {
		o.perRPCCredentials = creds
	}
}

func WithBlock() Option {
	return func(o *options) {
		o.block = true
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithMetrics(registerer metrics.Registerer) Option {
	return func(o *options) {
		o.metrics = registerer
	}
}

func WithStatsHandler(handler stats.Handler) Option {
	return func(o *options) {
		o.statsHandler = handler
	}
}

func NewClientConn(addr string, certPath string, opts ...Option) (*grpc.ClientConn, error) {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}
	ctx := context.Background()
	if options.dialTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.dialTimeout)
		defer cancel()
	}
	cred, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, err
	}
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithUnaryInterceptor(options.unaryInterceptor()),
		grpc.WithStatsHandler(options.statsHandler),
	}
	if options.perRPCCredentials != nil {
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(options.perRPCCredentials))
	}
	if options.block {
		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	return grpc.DialContext(ctx, addr, dialOptions...)
}

func (o *options) unaryInterceptor() grpc.UnaryClientInterceptor {
	if o.metrics == nil {
		return ChainUnaryClientInterceptors(
			XRequestIDUnaryClientInterceptor(),
			LogUnaryClientInterceptor(o.logger),
		)
	}
	registerMetrics(o.metrics)
	return ChainUnaryClientInterceptors(
		XRequestIDUnaryClientInterceptor(),
		LogUnaryClientInterceptor(o.logger),
		MetricsUnaryClientInterceptor(),
	)
}
