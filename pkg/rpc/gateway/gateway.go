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
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

type Gateway struct {
	httpServer *http.Server
	grpcAddr   string
	restAddr   string
	opts       *options
	logger     *zap.Logger
}

func NewGateway(grpcAddr, restAddr string, opts ...Option) (*Gateway, error) {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.logger == nil {
		logger, err := log.NewLogger()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %v", err)
		}
		options.logger = logger
	}

	return &Gateway{
		grpcAddr: grpcAddr,
		restAddr: restAddr,
		opts:     &options,
		logger:   options.logger.Named("gateway"),
	}, nil
}

func (g *Gateway) Start(ctx context.Context,
	registerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error,
) error {
	// Create a new ServeMux for the REST gateway
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// Create gRPC dial options with proper credentials and settings
	var dialOpts []grpc.DialOption
	if g.opts.certPath != "" {
		creds, err := credentials.NewClientTLSFromFile(g.opts.certPath, "")
		if err != nil {
			return fmt.Errorf("failed to create TLS credentials: %v", err)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		g.logger.Warn("starting gateway without TLS credentials")
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add perRPCCredentials if provided
	if g.opts.perRPCCredentials != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(g.opts.perRPCCredentials))
	}

	// Add keepalive parameters
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                g.opts.keepaliveTime,
		Timeout:             g.opts.keepaliveTimeout,
		PermitWithoutStream: g.opts.permitWithoutStream,
	}))

	// Add window size parameters
	dialOpts = append(dialOpts,
		grpc.WithInitialWindowSize(g.opts.initialWindowSize),
		grpc.WithInitialConnWindowSize(g.opts.initialConnWindowSize),
	)

	// Register the REST gateway handlers
	if err := registerFunc(ctx, mux, g.grpcAddr, dialOpts); err != nil {
		return fmt.Errorf("failed to register gateway handlers: %v", err)
	}

	// Create and start the HTTP server
	g.httpServer = &http.Server{
		Addr:    g.restAddr,
		Handler: mux,
	}

	go func() {
		var err error
		// Determine whether to use HTTPS or HTTP based on keyPath and certPath
		if g.opts.keyPath != "" && g.opts.certPath != "" {
			g.logger.Info("starting gateway with HTTPS",
				zap.String("cert_path", g.opts.certPath),
				zap.String("key_path", g.opts.keyPath))
			err = g.httpServer.ListenAndServeTLS(g.opts.certPath, g.opts.keyPath)
		} else {
			g.logger.Info("starting gateway with HTTP (no TLS)")
			err = g.httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			g.logger.Error("failed to serve HTTP/HTTPS", zap.Error(err))
		}
	}()

	g.logger.Info("gateway started",
		zap.String("grpc_addr", g.grpcAddr),
		zap.String("rest_addr", g.restAddr),
		zap.Bool("tls_enabled", g.opts.keyPath != "" && g.opts.certPath != ""),
	)

	return nil
}

func (g *Gateway) Stop(ctx context.Context) {
	if g.httpServer != nil {
		if err := g.httpServer.Shutdown(ctx); err != nil {
			g.logger.Error("failed to shutdown HTTP server", zap.Error(err))
		}
	}
}
