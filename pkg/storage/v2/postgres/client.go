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

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	postgres = "postgres"
)

type options struct {
	connMaxLifetime time.Duration
	maxOpenConns    int
	maxIdleConns    int
	logger          *zap.Logger
}

func defaultOptions() *options {
	return &options{
		connMaxLifetime: 300 * time.Second,
		maxOpenConns:    10,
		maxIdleConns:    5,
		logger:          zap.NewNop(),
	}
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
}

type Client interface {
	Execer
	Close() error
}

type client struct {
	db     *sql.DB
	opts   *options
	logger *zap.Logger
}

type Option func(*options)

func WithConnMaxLifetime(cml time.Duration) Option {
	return func(opts *options) {
		opts.connMaxLifetime = cml
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
	logger := dopts.logger.Named(postgres)
	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		postgres, dbUser, dbPass, dbHost, dbPort, dbName,
	)
	db, err := sql.Open(postgres, dsn)
	if err != nil {
		logger.Error("Failed to open db", zap.Error(err))
		return nil, err
	}
	db.SetConnMaxLifetime(dopts.connMaxLifetime)
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

func (c *client) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	var err error
	sret, err := c.db.ExecContext(ctx, query, args...)
	err = convertPostgresError(err)
	return &result{sret}, err
}

func (c *client) Close() error {
	return c.db.Close()
}
