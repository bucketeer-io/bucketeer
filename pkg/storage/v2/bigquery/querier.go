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

package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

type querierOptions struct {
	logger  *zap.Logger
	metrics metrics.Registerer
}

type QuerierOption func(*querierOptions)

func WithLogger(l *zap.Logger) QuerierOption {
	return func(opts *querierOptions) {
		opts.logger = l
	}
}

func WithMetrics(r metrics.Registerer) QuerierOption {
	return func(opts *querierOptions) {
		opts.metrics = r
	}
}

type Querier interface {
	ExecQuery(context.Context, string, []bigquery.QueryParameter) (*bigquery.RowIterator, error)
	Close() error
}

type querier struct {
	client *bigquery.Client
	logger *zap.Logger
}

func NewQuerier(
	ctx context.Context,
	project,
	location string,
	opts ...QuerierOption,
) (Querier, error) {
	dopts := &querierOptions{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	logger := dopts.logger.Named("bigquery")
	cli, err := bigquery.NewClient(ctx, project)
	if err != nil {
		logger.Error("Failed to create bigquery client", zap.Error(err))
		return nil, err
	}
	cli.Location = location
	return &querier{
		client: cli,
		logger: logger,
	}, nil
}

func (c *querier) ExecQuery(
	ctx context.Context,
	query string,
	params []bigquery.QueryParameter,
) (*bigquery.RowIterator, error) {
	var err error
	defer record()(operationQuery, &err)
	q := c.client.Query(query)
	q.Parameters = params
	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	c.logger.Debug(
		"Bigquery jobStatus",
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

func (c *querier) Close() error {
	return c.client.Close()
}
