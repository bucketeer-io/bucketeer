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

package database

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
)

func TestNewMySQLStorageClient_RunInTransactionV2(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("invokes callback with context", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mc := mysqlmock.NewMockClient(ctrl)
		mc.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).DoAndReturn(
			func(c context.Context, fn func(context.Context, mysql.Transaction) error) error {
				return fn(c, nil)
			},
		)
		db := NewMySQLStorageClient(mc)
		var ran bool
		err := db.RunInTransactionV2(ctx, func(c context.Context) error {
			assert.Equal(t, ctx, c)
			ran = true
			return nil
		})
		assert.NoError(t, err)
		assert.True(t, ran)
	})

	t.Run("propagates begin error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mc := mysqlmock.NewMockClient(ctrl)
		mc.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).Return(io.EOF)
		db := NewMySQLStorageClient(mc)
		err := db.RunInTransactionV2(ctx, func(context.Context) error { return nil })
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("propagates callback error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mc := mysqlmock.NewMockClient(ctrl)
		mc.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).DoAndReturn(
			func(c context.Context, fn func(context.Context, mysql.Transaction) error) error {
				return fn(c, nil)
			},
		)
		db := NewMySQLStorageClient(mc)
		err := db.RunInTransactionV2(ctx, func(context.Context) error {
			return errors.New("cb")
		})
		assert.EqualError(t, err, "cb")
	})
}

func TestNewMySQLStorageClient_Close(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	mc := mysqlmock.NewMockClient(ctrl)
	mc.EXPECT().Close().Return(io.EOF)
	db := NewMySQLStorageClient(mc)
	assert.ErrorIs(t, db.Close(), io.EOF)
}

func TestNewPostgresStorageClient_RunInTransactionV2(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("invokes callback with context", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		pc := pgmock.NewMockClient(ctrl)
		pc.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).DoAndReturn(
			func(c context.Context, fn func(context.Context, postgres.Transaction) error) error {
				return fn(c, nil)
			},
		)
		db := NewPostgresStorageClient(pc)
		var ran bool
		err := db.RunInTransactionV2(ctx, func(c context.Context) error {
			assert.Equal(t, ctx, c)
			ran = true
			return nil
		})
		assert.NoError(t, err)
		assert.True(t, ran)
	})

	t.Run("propagates begin error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		pc := pgmock.NewMockClient(ctrl)
		pc.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).Return(io.EOF)
		db := NewPostgresStorageClient(pc)
		err := db.RunInTransactionV2(ctx, func(context.Context) error { return nil })
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("propagates callback error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		pc := pgmock.NewMockClient(ctrl)
		pc.EXPECT().RunInTransactionV2(gomock.Any(), gomock.Any()).DoAndReturn(
			func(c context.Context, fn func(context.Context, postgres.Transaction) error) error {
				return fn(c, nil)
			},
		)
		db := NewPostgresStorageClient(pc)
		err := db.RunInTransactionV2(ctx, func(context.Context) error {
			return errors.New("cb")
		})
		assert.EqualError(t, err, "cb")
	})
}

func TestNewPostgresStorageClient_Close(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	pc := pgmock.NewMockClient(ctrl)
	pc.EXPECT().Close().Return(io.EOF)
	db := NewPostgresStorageClient(pc)
	assert.ErrorIs(t, db.Close(), io.EOF)
}
