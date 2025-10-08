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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/health"
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
	// Global metrics pushing
	StartContinuousPushing(pushGatewayURL, serviceName string, interval time.Duration)
	StopContinuousPushing()
	PushFinalMetrics()
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
	mux          *http.ServeMux
	server       *http.Server
	defaultPath  string
	registries   map[string]*registry
	opts         *options
	logger       *zap.Logger
	healthClient *http.Client
	// Global metrics pusher
	continuousPusher *ContinuousMetricsPusher
}

type registry struct {
	*prometheus.Registry
}

// ContinuousMetricsPusher handles continuous pushing of all metrics to pushgateway
type ContinuousMetricsPusher struct {
	pushGatewayURL string
	serviceName    string
	logger         *zap.Logger
	registries     map[string]*registry
	ticker         *time.Ticker
	done           chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewContinuousMetricsPusher creates a new continuous metrics pusher
func NewContinuousMetricsPusher(
	pushGatewayURL string,
	serviceName string,
	logger *zap.Logger,
	registries map[string]*registry,
) *ContinuousMetricsPusher {
	ctx, cancel := context.WithCancel(context.Background())

	return &ContinuousMetricsPusher{
		pushGatewayURL: pushGatewayURL,
		serviceName:    serviceName,
		logger:         logger.Named("global-metrics-pusher"),
		registries:     registries,
		done:           make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start begins continuous metrics pushing for all registries
func (mp *ContinuousMetricsPusher) Start(interval time.Duration) {
	if mp.pushGatewayURL == "" {
		mp.logger.Info("Push gateway URL not configured, skipping continuous metrics pushing")
		close(mp.done)
		return
	}

	mp.logger.Info("Starting global continuous metrics pusher",
		zap.String("service", mp.serviceName),
		zap.String("pushgateway_url", mp.pushGatewayURL),
		zap.Duration("interval", interval),
		zap.Int("registries_count", len(mp.registries)))

	mp.ticker = time.NewTicker(interval)

	go func() {
		defer mp.ticker.Stop()
		defer close(mp.done)

		// Push immediately on start
		mp.pushAllMetrics()

		for {
			select {
			case <-mp.ticker.C:
				mp.pushAllMetrics()
			case <-mp.ctx.Done():
				// Final push before shutdown
				mp.logger.Debug("Performing final global metrics push before shutdown")
				mp.pushAllMetrics()
				return
			}
		}
	}()
}

// pushAllMetrics pushes all metrics from all registries to pushgateway
func (mp *ContinuousMetricsPusher) pushAllMetrics() {
	for path, registry := range mp.registries {
		// Create a unique job name for each registry path
		jobName := fmt.Sprintf("global_metrics_%s", mp.serviceName)
		if path != "/metrics" {
			// Add path suffix for non-default paths
			jobName = fmt.Sprintf("global_metrics_%s_%s", mp.serviceName, path)
		}

		pusher := push.New(mp.pushGatewayURL, jobName).
			Gatherer(registry.Registry). // Push all metrics from this registry
			Grouping("bucketeer_service", mp.serviceName).
			Grouping("metrics_path", path)

		if err := pusher.Push(); err != nil {
			mp.logger.Error("Failed to push global metrics to Push Gateway",
				zap.Error(err),
				zap.String("pushgateway_url", mp.pushGatewayURL),
				zap.String("service", mp.serviceName),
				zap.String("path", path))
		}
	}
}

// Stop gracefully stops the continuous metrics pusher
func (mp *ContinuousMetricsPusher) Stop() {
	mp.logger.Info("Stopping global continuous metrics pusher")
	mp.cancel()
	<-mp.done
	mp.logger.Info("Global continuous metrics pusher stopped")
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
		healthClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
	r := m.Registerer(path)
	r.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
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
	return nil
}

func (m *metrics) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Stop continuous pushing before shutting down metrics server
	m.StopContinuousPushing()

	if err := m.server.Shutdown(ctx); err != nil {
		m.logger.Error("Failed to shutdown metrics server", zap.Error(err))
	}
}

// StartContinuousPushing begins continuous pushing of all metrics to pushgateway
func (m *metrics) StartContinuousPushing(pushGatewayURL, serviceName string, interval time.Duration) {
	if pushGatewayURL != "" && m.continuousPusher == nil {
		m.continuousPusher = NewContinuousMetricsPusher(
			pushGatewayURL,
			serviceName,
			m.logger,
			m.registries,
		)
		m.continuousPusher.Start(interval)
		// Startup logging handled in ContinuousMetricsPusher.Start()
	}
}

// StopContinuousPushing stops the continuous metrics pusher
func (m *metrics) StopContinuousPushing() {
	if m.continuousPusher != nil {
		m.continuousPusher.Stop()
		m.continuousPusher = nil
	}
}

// PushFinalMetrics performs a final push of all metrics (including shutdown metrics)
func (m *metrics) PushFinalMetrics() {
	if m.continuousPusher != nil {
		// Final push - logging handled in pushAllMetrics errors if needed
		m.continuousPusher.pushAllMetrics()
	}
}

func (m *metrics) Check(ctx context.Context) health.Status {
	resultCh := make(chan health.Status, 1)
	go func() {
		req, err := http.NewRequestWithContext(ctx, "GET", m.opts.healthCheckURL, nil)
		if err != nil {
			m.logger.Error("Failed to create health check request", zap.Error(err))
			resultCh <- health.Unhealthy
			return
		}
		resp, err := m.healthClient.Do(req)
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
