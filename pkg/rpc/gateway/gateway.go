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
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

// HandlerRegistrar is a function that registers a gRPC-Gateway handler
type HandlerRegistrar func(context.Context, *runtime.ServeMux, []grpc.DialOption) error

type Gateway struct {
	httpServer *http.Server
	restAddr   string
	opts       *options
	logger     *zap.Logger
}

func NewGateway(restAddr string, opts ...Option) (*Gateway, error) {
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
		restAddr: restAddr,
		opts:     &options,
		logger:   options.logger.Named("gateway"),
	}, nil
}

func (g *Gateway) Start(ctx context.Context,
	registerFuncs ...HandlerRegistrar,
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

	// Register all the provided handler registrars
	for _, registerFunc := range registerFuncs {
		if err := registerFunc(ctx, mux, dialOpts); err != nil {
			return fmt.Errorf("failed to register gateway handler: %v", err)
		}
	}

	// Wrap the mux with the boolean conversion middleware
	handler := BooleanConversionMiddleware(mux)

	// Create and start the HTTP server
	g.httpServer = &http.Server{
		Addr:    g.restAddr,
		Handler: handler,
	}

	// Start the server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		var err error
		if g.opts.keyPath != "" && g.opts.certPath != "" {
			err = g.httpServer.ListenAndServeTLS(g.opts.certPath, g.opts.keyPath)
		} else {
			g.logger.Info("starting gateway with HTTP (no TLS)")
			err = g.httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Check if there was an immediate error
	select {
	case err := <-errChan:
		return fmt.Errorf("failed to start gateway: %v", err)
	default:
		// No immediate error, server is starting
		g.logger.Debug("gateway started",
			zap.String("rest_addr", g.restAddr),
			zap.Bool("tls_enabled", g.opts.keyPath != "" && g.opts.certPath != ""),
			zap.Int("handlers_registered", len(registerFuncs)))
		return nil
	}
}

// Stop gracefully shuts down the HTTP server
func (g *Gateway) Stop(timeout time.Duration) {
	if g.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := g.httpServer.Shutdown(shutdownCtx); err != nil {
			g.logger.Error("failed to shutdown HTTP server gracefully", zap.Error(err))
		}
	}
}
