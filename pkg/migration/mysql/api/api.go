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

	libmigrate "github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/migration/mysql/migrate"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	migrationproto "github.com/bucketeer-io/bucketeer/proto/migration"
)

var (
	errInternal         = status.Error(codes.Internal, "migration-mysql: internal")
	errUnauthenticated  = status.Error(codes.Unauthenticated, "migration-mysql: unauthenticated")
	errPermissionDenied = status.Error(codes.PermissionDenied, "migration-mysql: permission denied")
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type MySQLService struct {
	migrateClientFactory migrate.ClientFactory
	opts                 *options
	logger               *zap.Logger
}

func NewMySQLService(migrateClientFactory migrate.ClientFactory, opts ...Option) *MySQLService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &MySQLService{
		migrateClientFactory: migrateClientFactory,
		opts:                 dopts,
		logger:               dopts.logger.Named("api"),
	}
}

func (s *MySQLService) Register(server *grpc.Server) {
	migrationproto.RegisterMigrationMySQLServiceServer(server, s)
}

func (s *MySQLService) MigrateAllMasterSchema(
	ctx context.Context,
	req *migrationproto.MigrateAllMasterSchemaRequest,
) (*migrationproto.MigrateAllMasterSchemaResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	migrateClient, err := s.migrateClientFactory.New()
	if err != nil {
		s.logger.Error(
			"Failed to new migrate client",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, errInternal
	}
	if err := migrateClient.Up(); err != nil {
		if errors.Is(err, libmigrate.ErrNoChange) {
			s.logger.Info("No change")
			return &migrationproto.MigrateAllMasterSchemaResponse{}, nil
		}
		s.logger.Error(
			"Failed to run migration",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, errInternal
	}
	return &migrationproto.MigrateAllMasterSchemaResponse{}, nil
}

func (s *MySQLService) RollbackMasterSchema(
	ctx context.Context,
	req *migrationproto.RollbackMasterSchemaRequest,
) (*migrationproto.RollbackMasterSchemaResponse, error) {
	_, err := s.checkSystemAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	migrateClient, err := s.migrateClientFactory.New()
	if err != nil {
		s.logger.Error(
			"Failed to new migrate client",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, errInternal
	}
	if err := migrateClient.Steps(-int(req.Step)); err != nil {
		s.logger.Error(
			"Failed to run rollback",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, errInternal
	}
	return &migrationproto.RollbackMasterSchemaResponse{}, nil
}

func (s *MySQLService) checkSystemAdminRole(ctx context.Context) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Info(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, errUnauthenticated
		case codes.PermissionDenied:
			s.logger.Info(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, errPermissionDenied
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, errInternal
		}
	}
	return editor, nil
}
