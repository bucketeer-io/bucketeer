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

package querier

import (
	"context"
	"os"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	bigqueryEmulatorHostEnv = "BIGQUERY_EMULATOR_HOST"
)

type options struct {
	logger  *zap.Logger
	metrics metrics.Registerer
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

type Client interface {
	ExecQuery(context.Context, string, []bigquery.QueryParameter) (*bigquery.RowIterator, error)
	Close() error
}

type client struct {
	cli    *bigquery.Client
	opts   *options
	logger *zap.Logger
}

func NewClient(
	ctx context.Context,
	project,
	location string,
	opts ...Option,
) (Client, error) {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	logger := dopts.logger.Named("bigquery-querier")
	var gcpOpts []option.ClientOption
	if bigqueryEmulatorEndpoint := os.Getenv(bigqueryEmulatorHostEnv); bigqueryEmulatorEndpoint != "" {
		gcpOpts = append(gcpOpts, option.WithoutAuthentication())
		gcpOpts = append(gcpOpts, option.WithEndpoint(bigqueryEmulatorEndpoint))
	}
	cli, err := bigquery.NewClient(ctx, project, gcpOpts...)
	if err != nil {
		logger.Error("Failed to create BigQuery client", zap.Error(err))
		return nil, err
	}
	cli.Location = location
	return &client{
		cli:    cli,
		opts:   dopts,
		logger: logger,
	}, nil
}

func (c *client) ExecQuery(
	ctx context.Context,
	query string,
	params []bigquery.QueryParameter,
) (*bigquery.RowIterator, error) {
	var err error
	defer record()(operationQuery, &err)
	q := c.cli.Query(query)
	q.Parameters = params
	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	c.logger.Debug(
		"BigQuery QueryJobStatus",
		zap.Any("status", status),
		zap.Any("query", query),
		zap.Any("params", params),
	)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}
	return job.Read(ctx)
}

func (c *client) Close() error {
	return c.cli.Close()
}
