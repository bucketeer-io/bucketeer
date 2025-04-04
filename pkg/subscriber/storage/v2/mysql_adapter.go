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

package v2

import (
	"context"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigquery/writer"
)

// MySQLAdapter adapts MySQL writer to BigQuery writer interface
type MySQLAdapter struct {
	mysqlWriter MySQLWriter
	logger      *zap.Logger
}

// AdapterOption is a functional option for configuring the MySQL adapter
type AdapterOption func(*adapterOptions)

type adapterOptions struct {
	logger  *zap.Logger
	metrics metrics.Registerer
}

// WithAdapterLogger sets the logger for the MySQL adapter
func WithAdapterLogger(l *zap.Logger) AdapterOption {
	return func(opts *adapterOptions) {
		opts.logger = l
	}
}

// WithAdapterMetrics sets the metrics for the MySQL adapter
func WithAdapterMetrics(r metrics.Registerer) AdapterOption {
	return func(opts *adapterOptions) {
		opts.metrics = r
	}
}

// NewMySQLAdapter creates a new MySQL adapter that implements the BigQuery writer interface
func NewMySQLAdapter(mysqlWriter MySQLWriter, opts ...AdapterOption) writer.Writer {
	dopts := &adapterOptions{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}

	return &MySQLAdapter{
		mysqlWriter: mysqlWriter,
		logger:      dopts.logger.Named("mysql-adapter"),
	}
}

// AppendRows adapts the MySQL writer's AppendRows method to the BigQuery writer interface
func (a *MySQLAdapter) AppendRows(ctx context.Context, batches [][][]byte) ([]int, error) {
	return a.mysqlWriter.AppendRows(ctx, batches)
}

// Close adapts the MySQL writer's Close method to the BigQuery writer interface
func (a *MySQLAdapter) Close() error {
	return a.mysqlWriter.Close()
}
