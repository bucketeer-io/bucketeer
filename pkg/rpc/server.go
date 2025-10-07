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
	"fmt"
	"net/http"
	"strings"
	"sync"
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
	certPath        string
	keyPath         string
	name            string
	logger          *zap.Logger
	port            int
	metrics         metrics.Registerer
	verifier        token.Verifier
	services        []Service
	handlers        []httpHandler
	rpcServer       *grpc.Server
	httpServer      *http.Server
	grpcWebServer   *grpcweb.WrappedGrpcServer
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutdownTracker *ShutdownTracker
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

	// Initialize shutdown tracker
	server.shutdownTracker = NewShutdownTracker(serverName, server.logger)

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
	// Start shutdown tracking
	s.shutdownTracker.StartShutdown("sigterm")

	shutdownStart := time.Now()
	s.logger.Info("Starting server graceful shutdown",
		zap.String("server", s.name),
		zap.Duration("timeout", timeout))

	var wg sync.WaitGroup
	shutdownErrors := make(chan error, 2)

	// Shutdown gRPC server
	if s.rpcServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			grpcStart := time.Now()
			s.logger.Info("Starting gRPC server graceful shutdown")

			// Use a shorter timeout for gRPC to leave time for HTTP
			grpcTimeout := timeout - time.Second
			if grpcTimeout < time.Second {
				grpcTimeout = time.Second
			}

			done := make(chan struct{})
			go func() {
				s.rpcServer.GracefulStop()
				close(done)
			}()

			select {
			case <-done:
				s.shutdownTracker.TrackShutdownDuration("grpc", time.Since(grpcStart), true)
				s.logger.Info("gRPC server shutdown completed gracefully")
			case <-time.After(grpcTimeout):
				s.logger.Warn("gRPC server graceful shutdown timed out, forcing stop")
				s.rpcServer.Stop()
				s.shutdownTracker.TrackShutdownDuration("grpc", time.Since(grpcStart), false)
			}
		}()
	}

	// Shutdown HTTP server
	if s.httpServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			httpStart := time.Now()
			s.logger.Info("Starting HTTP server graceful shutdown")

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if err := s.httpServer.Shutdown(ctx); err != nil {
				s.logger.Error("HTTP server failed to shut down gracefully", zap.Error(err))
				s.shutdownTracker.TrackShutdownDuration("http", time.Since(httpStart), false)
				shutdownErrors <- err
			} else {
				s.shutdownTracker.TrackShutdownDuration("http", time.Since(httpStart), true)
				s.logger.Info("HTTP server shutdown completed gracefully")
			}
		}()
	}

	// Wait for all shutdowns to complete or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("Server shutdown completed",
			zap.String("server", s.name),
			zap.Duration("total_duration", time.Since(shutdownStart)))
	case <-time.After(timeout):
		s.logger.Warn("Server shutdown timed out",
			zap.String("server", s.name),
			zap.Duration("timeout", timeout))
	}

	// Log any shutdown errors
	close(shutdownErrors)
	for err := range shutdownErrors {
		s.logger.Error("Shutdown error", zap.Error(err))
	}
}

func (s *Server) setupRPC() {
	creds, err := credentials.NewServerTLSFromFile(s.certPath, s.keyPath)
	if err != nil {
		s.logger.Fatal("Failed to read credentials: %v", zap.Error(err))
	}

	interceptors := []grpc.UnaryServerInterceptor{
		MetricsUnaryServerInterceptor(),
		s.shutdownTracker.ShutdownAwareUnaryInterceptor(), // Add shutdown tracking
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

	// Create gRPC-Web wrapper
	s.grpcWebServer = grpcweb.WrapServer(s.rpcServer)
}

func (s *Server) setupHTTP() {
	mux := http.NewServeMux()
	for _, handler := range s.handlers {
		mux.Handle(handler.path, handler)
	}

	// Add simple internal shutdown coordination endpoint
	mux.HandleFunc("/internal/shutdown-ready", func(w http.ResponseWriter, r *http.Request) {
		if s.shutdownTracker.IsShuttingDown() {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("not ready"))
		}
	})

	// Wrap the main handler with shutdown tracking middleware
	mainHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
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
		Handler:      s.shutdownTracker.ShutdownAwareHTTPMiddleware(mainHandler),
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
