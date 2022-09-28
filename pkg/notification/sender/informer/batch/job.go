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

package batch

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/job"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/notification/sender/informer"
)

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
}

var defaultOptions = options{
	logger: zap.NewNop(),
}

type Option func(*options)

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type jobInformer struct {
	manager *job.Manager
	opts    *options
	logger  *zap.Logger
	ctx     context.Context
	cancel  func()
	doneCh  chan struct{}
}

type Job struct {
	Name string
	Cron string
	Job  job.Job
}

func NewJobInformer(
	jobs []*Job,
	opts ...Option) (informer.Informer, error) {

	ctx, cancel := context.WithCancel(context.Background())
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	manager := job.NewManager(options.metrics, "ops_event_batch", options.logger)
	cji := &jobInformer{
		manager: manager,
		opts:    &options,
		logger:  options.logger.Named("sender"),
		ctx:     ctx,
		cancel:  cancel,
		doneCh:  make(chan struct{}),
	}
	if err := cji.registerJobs(jobs); err != nil {
		return nil, err
	}
	return cji, nil
}

func (i *jobInformer) Run() error {
	return i.manager.Run()
}

func (i *jobInformer) Stop() {
	i.manager.Stop()
}

// Check always returns healthy status.
// TODO: Implement Check() on job.Manager and do it here too.
func (i *jobInformer) Check(ctx context.Context) health.Status {
	return health.Healthy
}

func (i *jobInformer) registerJobs(jobs []*Job) error {
	for _, j := range jobs {
		if err := i.manager.AddCronJob(j.Name, j.Cron, j.Job); err != nil {
			i.logger.Error("Failed to add cron job",
				zap.String("name", j.Name),
				zap.String("cron", j.Cron),
				zap.Error(err))
			return err
		}
	}
	return nil
}
