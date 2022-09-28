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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package bigtable

import (
	"context"
	"errors"

	"cloud.google.com/go/bigtable"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

var (
	ErrKeyNotFound              = errors.New("storage: key not found")
	ErrColumnFamilyNotFound     = errors.New("storage: column family not found")
	ErrColumnNotFound           = errors.New("storage: column not found")
	ErrInternal                 = errors.New("storage: internal")
	errFailedToWritePartialRows = errors.New("storage: failed to write partial rows")
)

type Reader interface {
	ReadRows(ctx context.Context, request *ReadRequest) (Rows, error)
}

type Writer interface {
	WriteRow(ctx context.Context, request *WriteRequest) error
	WriteRows(ctx context.Context, request *WriteRequest) error
}

type ReadWriter interface {
	Reader
	Writer
}

type Client interface {
	Reader
	Writer
	Close() error
}

type client struct {
	client *bigtable.Client
	opts   *options
	logger *zap.Logger
}

type options struct {
	metrics metrics.Registerer
	logger  *zap.Logger
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

func NewBigtableClient(
	ctx context.Context,
	projectID, instance string,
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
	logger := dopts.logger.Named("bigtable")
	c, err := bigtable.NewClient(ctx, projectID, instance)
	if err != nil {
		logger.Error("Failed to create bigtable client", zap.Error(err))
		return nil, err
	}
	return &client{
		client: c,
		opts:   dopts,
		logger: logger,
	}, nil
}

func (c *client) ReadRows(ctx context.Context, req *ReadRequest) (Rows, error) {
	var err error
	defer record()(operationReadRows, &err)
	// Row set
	var rowSet bigtable.RowSet
	if req.RowSet.get() != nil {
		rowSet = req.RowSet.get()
	} else {
		rowSet = bigtable.RowRange{} // Read all keys
	}
	var rs []bigtable.Row
	if len(req.RowFilters) == 0 {
		rs, err = c.readRows(ctx, req.TableName, rowSet)
	} else {
		rs, err = c.readRowsWithFilter(ctx, req.TableName, rowSet, req.RowFilters)
	}
	if err != nil {
		c.logger.Error("Failed to read rows", zap.Error(err))
		return nil, ErrInternal
	}
	if len(rs) == 0 {
		err = ErrKeyNotFound
		return nil, ErrKeyNotFound
	}
	return &rows{
		rows:         rs,
		columnFamily: req.ColumnFamily,
		logger:       c.logger,
	}, nil
}

func (c *client) readRows(
	ctx context.Context,
	tableName string,
	rowSet bigtable.RowSet,
) ([]bigtable.Row, error) {
	tbl := c.client.Open(tableName)
	var rs []bigtable.Row
	err := tbl.ReadRows(
		ctx,
		rowSet,
		func(row bigtable.Row) bool {
			rs = append(rs, row)
			return true
		},
	)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (c *client) readRowsWithFilter(
	ctx context.Context,
	tableName string,
	rowSet bigtable.RowSet,
	rowFilters []RowFilter,
) ([]bigtable.Row, error) {
	// Read filters
	rf := makeFilters(rowFilters)
	tbl := c.client.Open(tableName)
	var rs []bigtable.Row
	err := tbl.ReadRows(
		ctx,
		rowSet,
		func(row bigtable.Row) bool {
			rs = append(rs, row)
			return true
		},
		bigtable.RowFilter(rf),
	)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (c *client) WriteRow(ctx context.Context, req *WriteRequest) error {
	var err error
	defer record()(operationWriteRow, &err)
	mut := bigtable.NewMutation()
	mut.Set(req.ColumnFamily, req.ColumnName, bigtable.Now(), req.Items[0].Value)
	tbl := c.client.Open(req.TableName)
	if err = tbl.Apply(ctx, req.Items[0].Key, mut); err != nil {
		c.logger.Error("Failed to write row", zap.Error(err), zap.String("rowKey", req.Items[0].Key))
		return ErrInternal
	}
	return nil
}

func (c *client) WriteRows(ctx context.Context, req *WriteRequest) error {
	var err error
	var errs []error
	defer record()(operationWriteRows, &err)
	muts := make([]*bigtable.Mutation, 0, len(req.Items))
	rowKeys := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		mut := bigtable.NewMutation()
		mut.Set(req.ColumnFamily, req.ColumnName, bigtable.Now(), item.Value)
		muts = append(muts, mut)
		rowKeys = append(rowKeys, item.Key)
	}
	tbl := c.client.Open(req.TableName)
	errs, err = tbl.ApplyBulk(ctx, rowKeys, muts)
	if err != nil {
		c.logger.Error("Failed to write rows",
			zap.Error(err),
			zap.Strings("rowKeys", rowKeys))
		return ErrInternal
	}
	if errs != nil {
		err = errFailedToWritePartialRows
		c.logger.Error("Failed to write partial rows",
			zap.Int("errs size", len(errs)),
			zap.Errors("errs", errs),
			zap.Strings("rowKeys", rowKeys))
		return ErrInternal
	}
	return nil
}

func (c *client) Close() error {
	var err error
	defer record()(operationClose, &err)
	if err = c.client.Close(); err != nil {
		c.logger.Error("Failed to close bigtable client", zap.Error(err))
		return ErrInternal
	}
	return nil
}

func makeFilters(filters []RowFilter) bigtable.Filter {
	if len(filters) == 1 {
		return filters[0].get()
	}
	chainFilters := make([]bigtable.Filter, 0, len(filters))
	for _, filter := range filters {
		chainFilters = append(chainFilters, filter.get())
	}
	return bigtable.ChainFilters(chainFilters...)
}
