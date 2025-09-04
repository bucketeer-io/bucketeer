// Copyright 2025 The Bucketeer Authors.
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

package storage

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/coderef/domain"
	err "github.com/bucketeer-io/bucketeer/pkg/error"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

var (
	//go:embed sql/code_reference/insert_code_reference.sql
	insertCodeReferenceSQL string
	//go:embed sql/code_reference/update_code_reference.sql
	updateCodeReferenceSQL string
	//go:embed sql/code_reference/select_code_reference.sql
	selectCodeReferenceSQL string
	//go:embed sql/code_reference/select_code_references.sql
	selectCodeReferencesSQL string
	//go:embed sql/code_reference/count_code_references.sql
	countCodeReferencesSQL string
	//go:embed sql/code_reference/delete_code_reference.sql
	deleteCodeReferenceSQL string
)

var (
	ErrCodeReferenceNotFound = err.NewErrorNotFound(
		err.CoderefPackageName,
		"code reference not found", "code_reference",
	)
	ErrCodeReferenceUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.CoderefPackageName,
		"code reference unexpected affected rows",
	)
)

type codeReferenceStorage struct {
	client mysql.Client
}

func NewCodeReferenceStorage(client mysql.Client) CodeReferenceStorage {
	return &codeReferenceStorage{client}
}

func (s *codeReferenceStorage) RunInTransaction(ctx context.Context, f func() error) error {
	tx, err := s.client.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("coderef: begin tx: %w", err)
	}
	ctx = context.WithValue(ctx, transactionKey, tx)
	return s.client.RunInTransaction(ctx, tx, f)
}

func (s *codeReferenceStorage) qe(ctx context.Context) mysql.QueryExecer {
	tx, ok := ctx.Value(transactionKey).(mysql.Transaction)
	if ok {
		return tx
	}
	return s.client
}

func (s *codeReferenceStorage) CreateCodeReference(ctx context.Context, cr *domain.CodeReference) error {
	_, err := s.qe(ctx).ExecContext(
		ctx,
		insertCodeReferenceSQL,
		cr.Id,
		cr.FeatureId,
		cr.FilePath,
		cr.FileExtension,
		cr.LineNumber,
		cr.CodeSnippet,
		cr.ContentHash,
		mysql.JSONObject{Val: cr.Aliases},
		cr.RepositoryName,
		cr.RepositoryOwner,
		cr.RepositoryType,
		cr.RepositoryBranch,
		cr.CommitHash,
		cr.EnvironmentId,
		cr.CreatedAt,
		cr.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *codeReferenceStorage) UpdateCodeReference(ctx context.Context, cr *domain.CodeReference) error {
	result, err := s.qe(ctx).ExecContext(
		ctx,
		updateCodeReferenceSQL,
		cr.FilePath,
		cr.FileExtension,
		cr.LineNumber,
		cr.CodeSnippet,
		cr.ContentHash,
		mysql.JSONObject{Val: cr.Aliases},
		cr.RepositoryName,
		cr.RepositoryOwner,
		cr.RepositoryType,
		cr.RepositoryBranch,
		cr.CommitHash,
		cr.UpdatedAt,
		cr.Id,
		cr.EnvironmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrCodeReferenceUnexpectedAffectedRows
	}
	return nil
}

func (s *codeReferenceStorage) GetCodeReference(
	ctx context.Context,
	id string,
) (*domain.CodeReference, error) {
	codeRef := &domain.CodeReference{}
	row := s.qe(ctx).QueryRowContext(
		ctx,
		selectCodeReferenceSQL,
		id,
	)
	err := row.Scan(
		&codeRef.Id,
		&codeRef.FeatureId,
		&codeRef.FilePath,
		&codeRef.FileExtension,
		&codeRef.LineNumber,
		&codeRef.CodeSnippet,
		&codeRef.ContentHash,
		&mysql.JSONObject{Val: &codeRef.Aliases},
		&codeRef.RepositoryName,
		&codeRef.RepositoryOwner,
		&codeRef.RepositoryType,
		&codeRef.RepositoryBranch,
		&codeRef.CommitHash,
		&codeRef.EnvironmentId,
		&codeRef.CreatedAt,
		&codeRef.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrCodeReferenceNotFound
		}
		return nil, err
	}
	return codeRef, nil
}

func (s *codeReferenceStorage) ListCodeReferences(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*domain.CodeReference, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf("%s %s %s %s", selectCodeReferencesSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe(ctx).QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	codeRefs := make([]*domain.CodeReference, 0, limit)
	for rows.Next() {
		codeRef := &domain.CodeReference{}
		err := rows.Scan(
			&codeRef.Id,
			&codeRef.FeatureId,
			&codeRef.FilePath,
			&codeRef.FileExtension,
			&codeRef.LineNumber,
			&codeRef.CodeSnippet,
			&codeRef.ContentHash,
			&mysql.JSONObject{Val: &codeRef.Aliases},
			&codeRef.RepositoryName,
			&codeRef.RepositoryOwner,
			&codeRef.RepositoryType,
			&codeRef.RepositoryBranch,
			&codeRef.CommitHash,
			&codeRef.EnvironmentId,
			&codeRef.CreatedAt,
			&codeRef.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		codeRefs = append(codeRefs, codeRef)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(codeRefs)
	var total int64
	countQuery := fmt.Sprintf("%s %s", countCodeReferencesSQL, whereSQL)
	row := s.qe(ctx).QueryRowContext(ctx, countQuery, whereArgs...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, 0, err
	}
	return codeRefs, nextOffset, total, nil
}

func (s *codeReferenceStorage) DeleteCodeReference(ctx context.Context, id string) error {
	result, err := s.qe(ctx).ExecContext(ctx, deleteCodeReferenceSQL, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCodeReferenceNotFound
	}
	return nil
}
