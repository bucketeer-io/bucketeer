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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/health"
)

type Registerer interface {
	MustRegister(...prometheus.Collector)
	Unregister(prometheus.Collector) bool
}

type Metrics interface {
	DefaultRegisterer() Registerer
	Registerer(path string) Registerer
	Check(ctx context.Context) health.Status
	Run() error
	Stop()
}

type options struct {
	healthCheckURL string
	logger         *zap.Logger
}

type Option func(*options)

func WithHealthCheckURL(url string) Option {
	return func(opts *options) {
		opts.healthCheckURL = url
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type metrics struct {
	mux         *http.ServeMux
	server      *http.Server
	defaultPath string
	registries  map[string]*registry
	opts        *options
	logger      *zap.Logger
}

type registry struct {
	*prometheus.Registry
}

func NewMetrics(port int, path string, opts ...Option) Metrics {
	dopts := &options{
		healthCheckURL: fmt.Sprintf("http://localhost:%d/health", port),
		logger:         zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	mux := http.NewServeMux()
	m := &metrics{
		mux: mux,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		defaultPath: path,
		registries:  make(map[string]*registry),
		opts:        dopts,
		logger:      dopts.logger.Named("metrics"),
	}
	r := m.Registerer(path)
	r.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)
	return m
}

func (m *metrics) DefaultRegisterer() Registerer {
	r := m.registries[m.defaultPath]
	return r
}

func (m *metrics) Registerer(path string) Registerer {
	if r, ok := m.registries[path]; ok {
		return r
	}
	r := &registry{Registry: prometheus.NewRegistry()}
	m.registries[path] = r
	return r
}

func (m *metrics) Run() error {
	m.logger.Info("Run started")
	for p, r := range m.registries {
		m.mux.Handle(p, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	}
	m.mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("healthy")) // nolint:errcheck
	})
	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		m.logger.Error("Failed to listen and serve", zap.Error(err))
		return err
	}
	m.logger.Info("Run finished")
	return nil
}

func (m *metrics) Stop() {
	m.logger.Info("Stop started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	m.server.Shutdown(ctx) // nolint:errcheck
	m.logger.Info("Stop finished")
}

func (m *metrics) Check(ctx context.Context) health.Status {
	resultCh := make(chan health.Status, 1)
	go func() {
		resp, err := http.Get(m.opts.healthCheckURL)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil || resp.StatusCode != 200 {
			m.logger.Error("Unhealthy", zap.Any("response", resp), zap.Error(err))
			resultCh <- health.Unhealthy
			return
		}
		resultCh <- health.Healthy
	}()
	select {
	case <-ctx.Done():
		m.logger.Error("Unhealthy due to context Done is closed", zap.Error(ctx.Err()))
		return health.Unhealthy
	case status := <-resultCh:
		return status
	}
}
