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
	"fmt"

	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

type queryOptions struct {
	logger  *zap.Logger
	metrics metrics.Registerer
}

type QueryOption func(*queryOptions)

func WithLogger(l *zap.Logger) QueryOption {
	return func(opts *queryOptions) {
		opts.logger = l
	}
}

func WithMetrics(r metrics.Registerer) QueryOption {
	return func(opts *queryOptions) {
		opts.metrics = r
	}
}

type Query interface {
	AppendRows(ctx context.Context, msgs [][]byte) error
}

type query struct {
	client *managedwriter.ManagedStream
	opts   *queryOptions
	logger *zap.Logger
}

func NewQuery(
	ctx context.Context,
	project, dataset, table string,
	desc protoreflect.MessageDescriptor,
	opts ...QueryOption,
) (Query, error) {
	dopts := &queryOptions{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	// if dopts.metrics != nil {
	// 	registerMetrics(dopts.metrics)
	// }
	logger := dopts.logger.Named("bigtable_query")
	c, err := managedwriter.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	descriptorProto, err := adapt.NormalizeDescriptor(desc)
	if err != nil {
		return nil, err
	}
	tableName := fmt.Sprintf(
		"projects/%s/datasets/%s/tables/%s",
		project,
		dataset,
		table,
	)
	managedStream, err := c.NewManagedStream(
		ctx,
		managedwriter.WithSchemaDescriptor(descriptorProto),
		managedwriter.WithDestinationTable(tableName),
		managedwriter.WithType(managedwriter.DefaultStream),
		managedwriter.EnableWriteRetries(true),
	)
	if err != nil {
		return nil, err
	}
	return &query{
		client: managedStream,
		opts:   dopts,
		logger: logger,
	}, nil
}

func (q *query) AppendRows(
	ctx context.Context,
	msgs [][]byte,
) error {
	results := []*managedwriter.AppendResult{}
	for i := 0; i < len(msgs); i += 10 {
		batch := msgs[i : i+10]
		r, err := q.client.AppendRows(ctx, batch)
		if err != nil {
			return err
		}
		results = append(results, r)
	}
	for _, r := range results {
		_, err := r.GetResult(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
