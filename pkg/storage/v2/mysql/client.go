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
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const dsnParams = "collation=utf8mb4_bin"

type options struct {
	connMaxIdleTime time.Duration
	maxOpenConns    int
	maxIdleConns    int
	logger          *zap.Logger
	metrics         metrics.Registerer
}

func defaultOptions() *options {
	return &options{
		connMaxIdleTime: 300 * time.Second,
		maxOpenConns:    10,
		maxIdleConns:    5,
		logger:          zap.NewNop(),
	}
}

type Option func(*options)

func WithConnMaxIdleTime(it time.Duration) Option {
	return func(opts *options) {
		opts.connMaxIdleTime = it
	}
}

func WithMaxOpenConns(moc int) Option {
	return func(opts *options) {
		opts.maxOpenConns = moc
	}
}

func WithMaxIdleConns(mic int) Option {
	return func(opts *options) {
		opts.maxIdleConns = mic
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
}

type QueryExecer interface {
	Queryer
	Execer
}

type Client interface {
	QueryExecer
	Close() error
	BeginTx(ctx context.Context) (Transaction, error)
	RunInTransaction(ctx context.Context, tx Transaction, f func() error) error
}

type client struct {
	db     *sql.DB
	opts   *options
	logger *zap.Logger
}

func NewClient(
	ctx context.Context,
	dbUser, dbPass, dbHost string,
	dbPort int,
	dbName string,
	opts ...Option,
) (Client, error) {
	dopts := defaultOptions()
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	logger := dopts.logger.Named("mysql")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dsnParams,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to open db", zap.Error(err))
		return nil, err
	}
	db.SetConnMaxIdleTime(dopts.connMaxIdleTime)
	db.SetMaxOpenConns(dopts.maxOpenConns)
	db.SetMaxIdleConns(dopts.maxIdleConns)
	if err := db.PingContext(ctx); err != nil {
		logger.Error("Failed to ping db", zap.Error(err))
		return nil, err
	}
	return &client{
		db:     db,
		opts:   dopts,
		logger: logger,
	}, nil
}

func (c *client) Close() error {
	return c.db.Close()
}

func (c *client) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	var err error
	defer record()(operationExec, &err)
	sret, err := c.db.ExecContext(ctx, query, args...)
	err = convertMySQLError(err)
	return &result{sret}, err
}

func (c *client) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	var err error
	defer record()(operationQuery, &err)
	srows, err := c.db.QueryContext(ctx, query, args...)
	return &rows{srows}, err
}

func (c *client) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	var err error
	defer record()(operationQueryRow, &err)
	r := &row{c.db.QueryRowContext(ctx, query, args...)}
	err = r.Err()
	return r
}

func (c *client) BeginTx(ctx context.Context) (Transaction, error) {
	var err error
	defer record()(operationBeginTx, &err)
	stx, err := c.db.BeginTx(ctx, nil)
	return &transaction{stx}, err
}

func (c *client) RunInTransaction(ctx context.Context, tx Transaction, f func() error) error {
	var err error
	defer record()(operationRunInTransaction, &err)
	defer func() {
		if err != nil {
			tx.Rollback() // nolint:errcheck
		}
	}()
	if err = f(); err == nil {
		err = tx.Commit()
	}
	return err
}
