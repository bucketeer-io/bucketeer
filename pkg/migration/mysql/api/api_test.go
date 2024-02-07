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

package api

import (
	"context"
	"errors"
	"testing"
	"time"

	libmigrate "github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/migration/mysql/migrate/mock"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	proto "github.com/bucketeer-io/bucketeer/proto/migration"
)

func TestNewMySQLService(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	cf := mock.NewMockClientFactory(mockController)
	logger := zap.NewNop()
	s := NewMySQLService(cf, WithLogger(logger))
	assert.IsType(t, &MySQLService{}, s)
}

func TestMigrateAllMasterSchema(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*MySQLService)
		req         *proto.MigrateAllMasterSchemaRequest
		expectedErr error
	}{
		{
			desc: "err: failed to new migrate client",
			setup: func(ms *MySQLService) {
				cf := mock.NewMockClientFactory(mockController)
				cf.EXPECT().New().Return(nil, errors.New("error"))
				ms.migrateClientFactory = cf
			},
			req:         &proto.MigrateAllMasterSchemaRequest{},
			expectedErr: errInternal,
		},
		{
			desc: "err: failed to run migration",
			setup: func(ms *MySQLService) {
				c := mock.NewMockClient(mockController)
				c.EXPECT().Up().Return(errors.New("error"))
				cf := mock.NewMockClientFactory(mockController)
				cf.EXPECT().New().Return(c, nil)
				ms.migrateClientFactory = cf
			},
			req:         &proto.MigrateAllMasterSchemaRequest{},
			expectedErr: errInternal,
		},
		{
			desc: "success: no change",
			setup: func(ms *MySQLService) {
				c := mock.NewMockClient(mockController)
				c.EXPECT().Up().Return(libmigrate.ErrNoChange)
				cf := mock.NewMockClientFactory(mockController)
				cf.EXPECT().New().Return(c, nil)
				ms.migrateClientFactory = cf
			},
			req:         &proto.MigrateAllMasterSchemaRequest{},
			expectedErr: nil,
		},
		{
			desc: "success",
			setup: func(ms *MySQLService) {
				c := mock.NewMockClient(mockController)
				c.EXPECT().Up().Return(nil)
				cf := mock.NewMockClientFactory(mockController)
				cf.EXPECT().New().Return(c, nil)
				ms.migrateClientFactory = cf
			},
			req:         &proto.MigrateAllMasterSchemaRequest{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newMySQLService(t)
			p.setup(service)
			_, err := service.MigrateAllMasterSchema(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestRollbackMasterSchema(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*MySQLService, int64)
		req         *proto.RollbackMasterSchemaRequest
		expectedErr error
	}{
		{
			desc: "err: failed to new migrate client",
			setup: func(ms *MySQLService, step int64) {
				cf := mock.NewMockClientFactory(mockController)
				cf.EXPECT().New().Return(nil, errors.New("error"))
				ms.migrateClientFactory = cf
			},
			req:         &proto.RollbackMasterSchemaRequest{},
			expectedErr: errInternal,
		},
		{
			desc: "err: failed to run migration",
			setup: func(ms *MySQLService, step int64) {
				c := mock.NewMockClient(mockController)
				c.EXPECT().Steps(-int(step)).Return(errors.New("error"))
				cf := mock.NewMockClientFactory(mockController)
				cf.EXPECT().New().Return(c, nil)
				ms.migrateClientFactory = cf
			},
			req:         &proto.RollbackMasterSchemaRequest{Step: 1},
			expectedErr: errInternal,
		},
		{
			desc: "success",
			setup: func(ms *MySQLService, step int64) {
				c := mock.NewMockClient(mockController)
				c.EXPECT().Steps(-int(step)).Return(nil)
				cf := mock.NewMockClientFactory(mockController)
				cf.EXPECT().New().Return(c, nil)
				ms.migrateClientFactory = cf
			},
			req:         &proto.RollbackMasterSchemaRequest{Step: 1},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := createContextWithToken(t)
			service := newMySQLService(t)
			p.setup(service, p.req.Step)
			_, err := service.RollbackMasterSchema(ctx, p.req)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestMySQLServicePermissionDenied(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		action   func(context.Context, *MySQLService) error
		expected error
	}{
		{
			desc: "MigrateAllMasterSchema",
			action: func(ctx context.Context, ms *MySQLService) error {
				_, err := ms.MigrateAllMasterSchema(ctx, &proto.MigrateAllMasterSchemaRequest{})
				return err
			},
			expected: errPermissionDenied,
		},
		{
			desc: "RollbackMasterSchema",
			action: func(ctx context.Context, ms *MySQLService) error {
				_, err := ms.RollbackMasterSchema(ctx, &proto.RollbackMasterSchemaRequest{})
				return err
			},
			expected: errPermissionDenied,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctx := createContextWithTokenRoleUnassigned(t)
			service := newMySQLService(t)
			actual := p.action(ctx, service)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newMySQLService(t *testing.T) *MySQLService {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &MySQLService{
		logger: logger.Named("api"),
	}
}

func createContextWithToken(t *testing.T) context.Context {
	t.Helper()
	token := &token.IDToken{
		Issuer:        "issuer",
		Subject:       "sub",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: true,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}

func createContextWithTokenRoleUnassigned(t *testing.T) context.Context {
	t.Helper()
	token := &token.IDToken{
		Issuer:        "issuer",
		Subject:       "sub",
		Audience:      "audience",
		Expiry:        time.Now().AddDate(100, 0, 0),
		IssuedAt:      time.Now(),
		Email:         "email",
		IsSystemAdmin: false,
	}
	ctx := context.TODO()
	return context.WithValue(ctx, rpc.Key, token)
}
