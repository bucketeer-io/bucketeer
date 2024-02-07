// Copyright 2023 The Bucketeer Authors.
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
//

package migration

import (
	"context"
	"embed"
	"fmt"

	libmigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	proto "github.com/bucketeer-io/bucketeer/proto/batch"
)

const mysqlParams = "collation=utf8mb4_bin"

//go:embed mysql/*.sql
var migrationSchemaFS embed.FS

type MysqlSchemaMigration struct {
	mysqlMigrate *libmigrate.Migrate
	logger       *zap.Logger
}

func NewMySQLSchemaMigration(
	mysqlUser, mysqlPass, mysqlHost, mysqlDBName string,
	mysqlPort int,
	logger *zap.Logger,
) *MysqlSchemaMigration {
	d, err := iofs.New(migrationSchemaFS, "mysql")
	if err != nil {
		logger.Error("failed to create a new source instance", zap.Error(err))
	}
	databaseURL := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%d)/%s?%s",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDBName, mysqlParams,
	)
	m, err := libmigrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		logger.Error("failed to create a new migrate instance", zap.Error(err))
	}
	return &MysqlSchemaMigration{
		mysqlMigrate: m,
		logger:       logger,
	}
}

func (m MysqlSchemaMigration) Migrate(ctx context.Context, request *proto.MigrationRequest) error {
	var steps int
	if request.Direction == proto.MigrationRequest_UP {
		steps = int(request.Steps)
	} else if request.Direction == proto.MigrationRequest_DOWN {
		steps = -int(request.Steps)
	}
	err := m.mysqlMigrate.Steps(steps)
	if err != nil {
		m.logger.Error(
			"failed to migrate",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Int32("steps", request.Steps),
			)...,
		)
		return err
	}
	return nil
}

func (m MysqlSchemaMigration) CurrentVersion(
	ctx context.Context,
	request *proto.MigrationVersionRequest,
) (uint, bool, error) {
	version, dirty, err := m.mysqlMigrate.Version()
	if err != nil {
		m.logger.Error("failed to get current version",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return 0, false, err
	}
	if int(version) == database.NilVersion {
		m.logger.Error("failed to get current version",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(libmigrate.ErrNilVersion))...,
		)
		return 0, false, libmigrate.ErrNilVersion
	}
	return version, dirty, nil
}
