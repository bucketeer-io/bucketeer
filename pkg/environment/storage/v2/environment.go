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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
)

var (
	ErrEnvironmentAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.EnvironmentPackageName,
		"environment already exists")
	ErrEnvironmentNotFound = pkgErr.NewErrorNotFound(
		pkgErr.EnvironmentPackageName,
		"environment not found",
		"environment")
	ErrEnvironmentUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.EnvironmentPackageName,
		"environment unexpected affected rows")

	//go:embed sql/environment/insert_environment.sql
	insertEnvironmentSQL string
	//go:embed sql/environment/update_environment.sql
	updateEnvironmentSQL string
	//go:embed sql/environment/select_environment.sql
	selectEnvironmentSQL string
	//go:embed sql/environment/select_environments.sql
	selectEnvironmentsSQL string
	//go:embed sql/environment/count_environments.sql
	countEnvironmentsSQL string
	//go:embed sql/environment/list_auto_archive_enabled_environments.sql
	listAutoArchiveEnabledEnvironmentsSQL string
	//go:embed sql/environment/delete_target_from_environment.sql
	deleteTargetFromEnvironmentSQL string
	//go:embed sql/environment/delete_environment.sql
	deleteEnvironmentSQL string
	//go:embed sql/environment/count_target_entities_in_environment.sql
	countTargetEntitiesInEnvironmentSQL string

	allowedTables = map[string]bool{
		"subscription": true, "experiment_result": true, "push": true,
		"ops_count": true, "auto_ops_rule": true, "segment_user": true,
		"segment": true, "goal": true, "experiment": true, "tag": true,
		"ops_progressive_rollout": true, "flag_trigger": true,
		"code_reference": true, "feature": true, "api_key": true,
		"audit_log": true, "account_v2": true,
	}
)

type EnvironmentStorage interface {
	CreateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error
	UpdateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error
	GetEnvironmentV2(ctx context.Context, id string) (*domain.EnvironmentV2, error)
	ListEnvironmentsV2(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*proto.EnvironmentV2, int, int64, error)
	ListAutoArchiveEnabledEnvironments(ctx context.Context) ([]*domain.EnvironmentV2, error)
	DeleteTargetFromEnvironmentV2(
		ctx context.Context,
		environmentID string,
		targetID string,
	) error
	DeleteEnvironmentV2(ctx context.Context, whereParts []mysql.WherePart) error
	CountTargetEntitiesInEnvironmentV2(
		ctx context.Context,
		environmentID string,
		target string,
	) (int64, error)
}

type environmentStorage struct {
	qe mysql.QueryExecer
}

func NewEnvironmentStorage(qe mysql.QueryExecer) EnvironmentStorage {
	return &environmentStorage{qe}
}

func (s *environmentStorage) CreateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertEnvironmentSQL,
		e.Id,
		e.Name,
		e.UrlCode,
		e.Description,
		e.ProjectId,
		e.OrganizationId,
		e.Archived,
		e.RequireComment,
		e.CreatedAt,
		e.UpdatedAt,
		e.AutoArchiveEnabled,
		e.AutoArchiveUnusedDays,
		e.AutoArchiveCheckCodeRefs,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrEnvironmentAlreadyExists
		}
		return err
	}
	return nil
}

func (s *environmentStorage) UpdateEnvironmentV2(ctx context.Context, e *domain.EnvironmentV2) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateEnvironmentSQL,
		e.Name,
		e.Description,
		e.Archived,
		e.RequireComment,
		e.CreatedAt,
		e.UpdatedAt,
		e.AutoArchiveEnabled,
		e.AutoArchiveUnusedDays,
		e.AutoArchiveCheckCodeRefs,
		e.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrEnvironmentUnexpectedAffectedRows
	}
	return nil
}

func (s *environmentStorage) GetEnvironmentV2(ctx context.Context, id string) (*domain.EnvironmentV2, error) {
	environment := proto.EnvironmentV2{}
	err := s.qe.QueryRowContext(
		ctx,
		selectEnvironmentSQL,
		id,
	).Scan(
		&environment.Id,
		&environment.Name,
		&environment.UrlCode,
		&environment.Description,
		&environment.ProjectId,
		&environment.OrganizationId,
		&environment.Archived,
		&environment.RequireComment,
		&environment.CreatedAt,
		&environment.UpdatedAt,
		&environment.AutoArchiveEnabled,
		&environment.AutoArchiveUnusedDays,
		&environment.AutoArchiveCheckCodeRefs,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}
	return &domain.EnvironmentV2{EnvironmentV2: &environment}, nil
}

func (s *environmentStorage) ListEnvironmentsV2(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.EnvironmentV2, int, int64, error) {
	// Because select_environments.sql defines the variable strings in a complex constructed way,
	//  we do not use ConstructQueryAndWhereArgs() here.
	var query string
	var whereArgs []any
	if options != nil {
		var whereSQL string
		whereParts := options.CreateWhereParts()
		whereSQL, whereArgs = mysql.ConstructWhereSQLString(whereParts)
		orderBySQL := mysql.ConstructOrderBySQLString(options.Orders)
		limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
		query = fmt.Sprintf(selectEnvironmentsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	} else {
		query = selectEnvironmentsSQL
		whereArgs = []interface{}{}
	}

	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	defer rows.Close()
	environments := make([]*proto.EnvironmentV2, 0, limit)
	for rows.Next() {
		environment := proto.EnvironmentV2{}
		err := rows.Scan(
			&environment.Id,
			&environment.Name,
			&environment.UrlCode,
			&environment.Description,
			&environment.ProjectId,
			&environment.OrganizationId,
			&environment.Archived,
			&environment.RequireComment,
			&environment.CreatedAt,
			&environment.UpdatedAt,
			&environment.AutoArchiveEnabled,
			&environment.AutoArchiveUnusedDays,
			&environment.AutoArchiveCheckCodeRefs,
			&environment.FeatureFlagCount,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		environments = append(environments, &environment)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(environments)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countEnvironmentsSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return environments, nextOffset, totalCount, nil
}

func (s *environmentStorage) ListAutoArchiveEnabledEnvironments(ctx context.Context) ([]*domain.EnvironmentV2, error) {
	rows, err := s.qe.QueryContext(ctx, listAutoArchiveEnabledEnvironmentsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	environments := make([]*domain.EnvironmentV2, 0)
	for rows.Next() {
		environment := proto.EnvironmentV2{}
		err := rows.Scan(
			&environment.Id,
			&environment.Name,
			&environment.UrlCode,
			&environment.Description,
			&environment.ProjectId,
			&environment.OrganizationId,
			&environment.Archived,
			&environment.RequireComment,
			&environment.CreatedAt,
			&environment.UpdatedAt,
			&environment.AutoArchiveEnabled,
			&environment.AutoArchiveUnusedDays,
			&environment.AutoArchiveCheckCodeRefs,
		)
		if err != nil {
			return nil, err
		}
		environments = append(environments, &domain.EnvironmentV2{EnvironmentV2: &environment})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return environments, nil
}

func (s *environmentStorage) DeleteTargetFromEnvironmentV2(
	ctx context.Context,
	environmentID string,
	target string,
) error {
	if !allowedTables[target] {
		return fmt.Errorf("table %s is not allowed to delete from", target)
	}
	args := []interface{}{
		environmentID,
	}

	query := fmt.Sprintf(deleteTargetFromEnvironmentSQL, target)
	_, err := s.qe.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *environmentStorage) DeleteEnvironmentV2(ctx context.Context, whereParts []mysql.WherePart) error {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	query := fmt.Sprintf(deleteEnvironmentSQL, whereSQL)
	_, err := s.qe.ExecContext(
		ctx,
		query,
		whereArgs...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *environmentStorage) CountTargetEntitiesInEnvironmentV2(
	ctx context.Context,
	environmentID string,
	target string,
) (int64, error) {
	rows, err := s.qe.QueryContext(
		ctx,
		fmt.Sprintf(countTargetEntitiesInEnvironmentSQL, target),
		environmentID,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var count int64
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	if rows.Err() != nil {
		return 0, err
	}
	return count, nil
}
