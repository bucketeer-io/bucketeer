// Copyright 2026 The Bucketeer Authors.
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

// Package database provides a dialect-agnostic transactional client for API and
// composition layers. Query execution stays on mysql.Client / postgres.Client
// (QueryExecer); those clients attach the active transaction to context inside
// RunInTransactionV2, so callbacks here only need context.
package database

import (
	"context"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
)

// Client is the unified transactional surface (see RFC: polymorphic database storage).
type Client interface {
	RunInTransactionV2(ctx context.Context, f func(ctx context.Context) error) error
	Close() error
}

type mysqlStorageClientAdapter struct {
	inner mysql.Client
}

// NewMySQLStorageClient wraps mysql.Client for callers that must not import mysql.
func NewMySQLStorageClient(c mysql.Client) Client {
	return &mysqlStorageClientAdapter{inner: c}
}

func (a *mysqlStorageClientAdapter) Close() error {
	return a.inner.Close()
}

func (a *mysqlStorageClientAdapter) RunInTransactionV2(
	ctx context.Context,
	f func(ctx context.Context) error,
) error {
	return a.inner.RunInTransactionV2(ctx, func(ctx context.Context, _ mysql.Transaction) error {
		return f(ctx)
	})
}

type postgresStorageClientAdapter struct {
	inner postgres.Client
}

// NewPostgresStorageClient wraps postgres.Client for callers that must not import postgres.
func NewPostgresStorageClient(c postgres.Client) Client {
	return &postgresStorageClientAdapter{inner: c}
}

func (a *postgresStorageClientAdapter) Close() error {
	return a.inner.Close()
}

func (a *postgresStorageClientAdapter) RunInTransactionV2(
	ctx context.Context,
	f func(ctx context.Context) error,
) error {
	return a.inner.RunInTransactionV2(ctx, func(ctx context.Context, _ postgres.Transaction) error {
		return f(ctx)
	})
}
