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
	// Create a client connection to the gRPC server
	conn, err := g.createClientConn()
	if err != nil {
		return fmt.Errorf("failed to create client connection: %v", err)
	}
	defer conn.Close()

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

	// Register the REST gateway handlers
	if err := registerFunc(ctx,
		mux,
		g.grpcAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return fmt.Errorf("failed to register gateway handlers: %v", err)
	}

	// Create and start the HTTP server
	g.httpServer = &http.Server{
		Addr:    g.restAddr,
		Handler: mux,
	}

	go func() {
		if err := g.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			g.logger.Error("failed to serve HTTP", zap.Error(err))
		}
	}()

	g.logger.Info("gateway started",
		zap.String("grpc_addr", g.grpcAddr),
		zap.String("rest_addr", g.restAddr),
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

func (g *Gateway) createClientConn() (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                g.opts.keepaliveTime,
			Timeout:             g.opts.keepaliveTimeout,
			PermitWithoutStream: g.opts.permitWithoutStream,
		}),
		grpc.WithInitialWindowSize(g.opts.initialWindowSize),
		grpc.WithInitialConnWindowSize(g.opts.initialConnWindowSize),
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.opts.timeout)
	defer cancel()

	return grpc.DialContext(ctx, g.grpcAddr, opts...)
}
