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

package rpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/token"
)

type Server struct {
	certPath   string
	keyPath    string
	name       string
	logger     *zap.Logger
	port       int
	metrics    metrics.Registerer
	verifier   token.Verifier
	services   []Service
	handlers   []httpHandler
	rpcServer  *grpc.Server
	httpServer *http.Server
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

func WithHandler(path string, handler http.Handler) Option {
	return func(s *Server) {
		s.handlers = append(s.handlers, httpHandler{Handler: handler, path: path})
	}
}

func NewServer(service Service, certPath, keyPath, serverName string, opt ...Option) *Server {
	server := &Server{
		port:   9000,
		name:   serverName,
		logger: zap.NewNop(),
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
	s.logger.Info(fmt.Sprintf("Running on %d", s.port))
	s.runServer()
}

func (s *Server) Stop(timeout time.Duration) {
	s.logger.Info("Server is going to shut down")
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Server failed to shut down", zap.Error(err))
	}
	s.logger.Info("Server has shut down gracefully",
		zap.Duration("elapsedTime", time.Since(startTime)),
	)
}

func (s *Server) setupRPC() {
	creds, err := credentials.NewServerTLSFromFile(s.certPath, s.keyPath)
	if err != nil {
		s.logger.Fatal("Failed to read credentials: %v", zap.Error(err))
	}
	interceptor := chainUnaryServerInterceptors(
		LogUnaryServerInterceptor(s.logger),
		MetricsUnaryServerInterceptor(),
	)
	if s.verifier != nil {
		interceptor = chainUnaryServerInterceptors(
			interceptor,
			AuthUnaryServerInterceptor(s.verifier))
	}
	s.rpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(interceptor),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)
	for _, service := range s.services {
		service.Register(s.rpcServer)
	}
}

func (s *Server) setupHTTP() {
	mux := http.NewServeMux()
	for _, handler := range s.handlers {
		mux.Handle(handler.path, handler)
	}
	s.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", s.port),
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if isRPC(req) {
				s.rpcServer.ServeHTTP(resp, req)
			} else {
				mux.ServeHTTP(resp, req)
			}
		}),
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
