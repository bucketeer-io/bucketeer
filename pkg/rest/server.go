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

package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

type Server struct {
	certPath   string
	keyPath    string
	port       int
	metrics    metrics.Registerer
	httpServer *http.Server
	mux        *http.ServeMux
	logger     *zap.Logger
	services   []Service
}

type Option func(*Server)

const httpName = "http"

func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithMetrics(registerer metrics.Registerer) Option {
	return func(s *Server) {
		s.metrics = registerer
	}
}

func WithService(service Service) Option {
	return func(s *Server) {
		s.services = append(s.services, service)
	}
}

func NewServer(certPath, keyPath string, opt ...Option) *Server {
	server := &Server{
		port:   8000,
		logger: zap.NewNop(),
		mux:    http.NewServeMux(),
	}
	for _, o := range opt {
		o(server)
	}
	server.logger = server.logger.Named(httpName)
	if len(certPath) == 0 {
		server.logger.Fatal("CertPath must not be empty")
	}
	server.certPath = certPath
	if len(keyPath) == 0 {
		server.logger.Fatal("KeyPath must not be empty")
	}
	server.keyPath = keyPath
	if len(server.services) == 0 {
		server.logger.Fatal("Service must not be nil")
	}
	return server
}

func (s *Server) Run() {
	if s.metrics != nil {
		registerMetrics(s.metrics)
	}
	s.setup()
	s.logger.Info(fmt.Sprintf("Rest server is running on %d", s.port))
	s.runServer()
}

func (s *Server) Stop(timeout time.Duration) {
	s.logger.Info("Rest server is going to shut down")
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Rest server failed to shut down", zap.Error(err))
	}
	s.logger.Info("Rest server has shut down gracefully",
		zap.Duration("elapsedTime", time.Since(startTime)),
	)
}

func (s *Server) setup() {
	mws := newMiddleWares()
	mws.Append(LogServerMiddleware(s.logger))
	mws.Append(MetricsServerMiddleware)
	for _, service := range s.services {
		service.Register(s.mux)
	}
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mws.Handle(s.mux),
	}
}

func (s *Server) runServer() {
	err := s.httpServer.ListenAndServeTLS(s.certPath, s.keyPath)
	if err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Failed to serve", zap.Error(err))
	}
}
