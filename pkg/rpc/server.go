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

package rpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/bucketeer-io/bucketeer/v2/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
)

type Server struct {
	certPath      string
	keyPath       string
	name          string
	logger        *zap.Logger
	port          int
	metrics       metrics.Registerer
	verifier      token.Verifier
	services      []Service
	handlers      []httpHandler
	rpcServer     *grpc.Server
	httpServer    *http.Server
	grpcWebServer *grpcweb.WrappedGrpcServer // DEPRECATED: Remove once Node.js SDK migrates away from grpc-web
	readTimeout   time.Duration
	writeTimeout  time.Duration
	idleTimeout   time.Duration
}

type httpHandler struct {
	http.Handler
	path string
}

type Option func(*Server)

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithVerifier(verifier token.Verifier) Option {
	return func(s *Server) {
		s.verifier = verifier
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithService(service Service) Option {
	return func(s *Server) {
		s.services = append(s.services, service)
	}
}

func WithMetrics(registerer metrics.Registerer) Option {
	return func(s *Server) {
		s.metrics = registerer
	}
}

func WithTimeouts(readTimeout, writeTimeout, idleTimeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = readTimeout
		s.writeTimeout = writeTimeout
		s.idleTimeout = idleTimeout
	}
}

func WithHandler(path string, handler http.Handler) Option {
	return func(s *Server) {
		s.handlers = append(s.handlers, httpHandler{Handler: handler, path: path})
	}
}

func NewServer(service Service, certPath, keyPath, serverName string, opt ...Option) *Server {
	server := &Server{
		port:         9000,
		name:         serverName,
		logger:       zap.NewNop(),
		readTimeout:  30 * time.Second, // Default timeout
		writeTimeout: 30 * time.Second, // Default timeout
		idleTimeout:  60 * time.Second, // Default timeout
	}
	for _, o := range opt {
		o(server)
	}
	server.logger = server.logger.Named(fmt.Sprintf("rpc-server.%s", serverName))

	if len(certPath) == 0 {
		server.logger.Fatal("CertPath must not be empty")
	}
	server.certPath = certPath
	if len(keyPath) == 0 {
		server.logger.Fatal("KeyPath must not be empty")
	}
	server.keyPath = keyPath
	if service == nil {
		server.logger.Fatal("Service must not be nil")
	}
	server.services = append(server.services, service)
	return server
}

func (s *Server) Run() {
	if s.metrics != nil {
		registerMetrics(s.metrics)
	}
	s.setupRPC()
	s.setupHTTP()
	s.runServer()
}

func (s *Server) Stop(timeout time.Duration) {
	// Shutdown order is critical:
	// 1. HTTP server first (drains REST/gRPC-Gateway requests)
	// 2. gRPC server second (only pure gRPC connections remain)
	//
	// This ensures HTTP requests that call s.rpcServer.ServeHTTP() can complete
	// before we stop the underlying gRPC server.
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("HTTP server failed to shut down gracefully", zap.Error(err))
		}
	}

	// Note: We use Stop() instead of GracefulStop() because:
	// - GracefulStop() panics on HTTP-served connections (serverHandlerTransport.Drain not implemented)
	// - HTTP-served connections were already drained in step 1
	// - Pure gRPC clients have retry logic and Envoy connection draining to handle this
	if s.rpcServer != nil {
		s.rpcServer.Stop()
	}
}

func (s *Server) setupRPC() {
	creds, err := credentials.NewServerTLSFromFile(s.certPath, s.keyPath)
	if err != nil {
		s.logger.Fatal("Failed to read credentials: %v", zap.Error(err))
	}

	interceptors := []grpc.UnaryServerInterceptor{
		MetricsUnaryServerInterceptor(),
	}

	if s.verifier != nil {
		interceptors = append(interceptors, AuthUnaryServerInterceptor(s.verifier))
	}

	s.rpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(interceptors...),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)
	for _, service := range s.services {
		service.Register(s.rpcServer)
	}

	// DEPRECATED: grpc-web support for legacy Node.js SDK
	// TODO: Remove once Node.js SDK migrates to gRPC-Gateway (REST) or pure gRPC
	// This is an abandoned library (last updated 2021) with known issues.
	s.grpcWebServer = grpcweb.WrapServer(s.rpcServer)
}

func (s *Server) setupHTTP() {
	mux := http.NewServeMux()
	for _, handler := range s.handlers {
		mux.Handle(handler.path, handler)
	}

	// Wrap the main handler with gRPC-web support and routing
	mainHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// DEPRECATED: grpc-web support for legacy Node.js SDK
		// This check should be removed once Node.js SDK migrates away from grpc-web
		if s.grpcWebServer.IsGrpcWebRequest(req) {
			s.grpcWebServer.ServeHTTP(resp, req)
		} else if isRPC(req) {
			s.rpcServer.ServeHTTP(resp, req)
		} else {
			mux.ServeHTTP(resp, req)
		}
	})

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		Handler:      mainHandler,
	}
}

func (s *Server) runServer() {
	err := s.httpServer.ListenAndServeTLS(s.certPath, s.keyPath)
	if err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Failed to serve", zap.Error(err))
	}
}

func isRPC(req *http.Request) bool {
	if req.ProtoMajor == 2 &&
		strings.HasPrefix(req.Header.Get("Content-Type"), "application/grpc") {
		return true
	}
	return false
}
