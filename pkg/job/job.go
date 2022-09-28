// Copyright 2022 The Bucketeer Authors.
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

package job

import (
	"context"
	"time"

	"github.com/robfig/cron"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

type options struct {
	timeout time.Duration
	metrics metrics.Registerer
	logger  *zap.Logger
}

type Option func(*options)

func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type Job interface {
	Run(context.Context) error
}

type Manager struct {
	cron    *cron.Cron
	metrics metrics.Registerer
	logger  *zap.Logger
	ctx     context.Context
	cancel  func()
	doneCh  chan struct{}
}

func NewManager(r metrics.Registerer, subsystem string, logger *zap.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	registerMetrics(r, subsystem)
	return &Manager{
		cron:    cron.New(),
		metrics: r,
		logger:  logger.Named("jobmanager"),
		ctx:     ctx,
		cancel:  cancel,
		doneCh:  make(chan struct{}),
	}
}

func (m *Manager) Run() error {
	m.logger.Info("Run started")
	defer close(m.doneCh)
	m.cron.Start()
	<-m.ctx.Done()
	m.logger.Info("Run finished")
	return nil
}

func (m *Manager) Stop() {
	m.logger.Info("Stop started")
	m.cancel()
	m.cron.Stop()
	<-m.doneCh
	m.logger.Info("Stop finished")
}

func (m *Manager) AddCronJob(name, cron string, job Job) error {
	return m.cron.AddFunc(cron, func() {
		m.logger.Info("Job started", zap.String("name", name))
		startTime := time.Now()
		startedJobCounter.WithLabelValues(name).Inc()
		err := job.Run(m.ctx)
		code := codeSuccess
		if err != nil {
			code = codeFail
			m.logger.Error("Job finished with an error", zap.String("name", name), zap.Error(err))
		} else {
			m.logger.Info("Job finished", zap.String("name", name))
		}
		finishedJobCounter.WithLabelValues(name, code).Inc()
		finishedJobHistogram.WithLabelValues(name, code).Observe(time.Since(startTime).Seconds())
	})
}
