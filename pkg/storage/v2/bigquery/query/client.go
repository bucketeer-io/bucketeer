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

package query

import (
	"context"

	"cloud.google.com/go/bigquery/storage/managedwriter"
	"go.uber.org/zap"
	"google.golang.org/protobuf/reflect/protodesc"
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
	Close() error
}

type query struct {
	client *managedwriter.ManagedStream
	opts   *queryOptions
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
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	c, err := managedwriter.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	managedStream, err := c.NewManagedStream(
		ctx,
		managedwriter.WithSchemaDescriptor(protodesc.ToDescriptorProto(desc)),
		managedwriter.WithDestinationTable(
			managedwriter.TableParentFromParts(project, dataset, table),
		),
		managedwriter.WithType(managedwriter.DefaultStream),
		managedwriter.EnableWriteRetries(true),
	)
	if err != nil {
		return nil, err
	}
	return &query{
		client: managedStream,
		opts:   dopts,
	}, nil
}

func (q *query) AppendRows(
	ctx context.Context,
	msgs [][]byte,
) error {
	var err error
	defer record()(operationQuery, &err)
	results := []*managedwriter.AppendResult{}
	for i := 0; i < len(msgs); i += 10 {
		end := i + 10
		if end > len(msgs) {
			end = len(msgs)
		}
		batch := msgs[i:end]
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

func (q *query) Close() error {
	return q.client.Close()
}
