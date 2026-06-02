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

package postgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/coderef/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/coderef/storage"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	coderefproto "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
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
	//go:embed sql/code_reference/count_code_references_by_feature_ids.sql
	countCodeReferencesByFeatureIDsSQL string
)

type codeReferenceStorage struct {
	qe pgstorage.QueryExecer
}

func NewCodeReferenceStorage(qe pgstorage.QueryExecer) storage.CodeReferenceStorage {
	return &codeReferenceStorage{qe}
}

func (s *codeReferenceStorage) CreateCodeReference(ctx context.Context, cr *domain.CodeReference) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertCodeReferenceSQL,
		cr.Id,
		cr.FeatureId,
		cr.FilePath,
		cr.FileExtension,
		cr.LineNumber,
		cr.CodeSnippet,
		cr.ContentHash,
		pgstorage.JSONObject{Val: cr.Aliases},
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
	result, err := s.qe.ExecContext(
		ctx,
		updateCodeReferenceSQL,
		cr.FilePath,
		cr.FileExtension,
		cr.LineNumber,
		cr.CodeSnippet,
		cr.ContentHash,
		pgstorage.JSONObject{Val: cr.Aliases},
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
		return storage.ErrCodeReferenceUnexpectedAffectedRows
	}
	return nil
}

func (s *codeReferenceStorage) GetCodeReference(
	ctx context.Context,
	id string,
) (*domain.CodeReference, error) {
	codeRef := &domain.CodeReference{}
	row := s.qe.QueryRowContext(ctx, selectCodeReferenceSQL, id)
	err := row.Scan(
		&codeRef.Id,
		&codeRef.FeatureId,
		&codeRef.FilePath,
		&codeRef.FileExtension,
		&codeRef.LineNumber,
		&codeRef.CodeSnippet,
		&codeRef.ContentHash,
		&pgstorage.JSONObject{Val: &codeRef.Aliases},
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
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, storage.ErrCodeReferenceNotFound
		}
		return nil, err
	}
	return codeRef, nil
}

func (s *codeReferenceStorage) ListCodeReferences(
	ctx context.Context,
	params storage.ListCodeReferencesParams,
) ([]*domain.CodeReference, int, int64, error) {
	options, err := listCodeReferencesOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	whereParts := options.CreateWhereParts()
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(whereParts)
	orderBySQL := pgstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	query := fmt.Sprintf("%s %s %s %s", selectCodeReferencesSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	codeRefs := make([]*domain.CodeReference, 0, options.Limit)
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
			&pgstorage.JSONObject{Val: &codeRef.Aliases},
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
	nextOffset := options.Offset + len(codeRefs)
	var total int64
	countQuery := fmt.Sprintf("%s %s", countCodeReferencesSQL, whereSQL)
	row := s.qe.QueryRowContext(ctx, countQuery, whereArgs...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, 0, err
	}
	return codeRefs, nextOffset, total, nil
}

func (s *codeReferenceStorage) DeleteCodeReference(ctx context.Context, id string) error {
	result, err := s.qe.ExecContext(ctx, deleteCodeReferenceSQL, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrCodeReferenceNotFound
	}
	return nil
}

func (s *codeReferenceStorage) GetCodeReferenceCountsByFeatureIDs(
	ctx context.Context,
	environmentID string,
	featureIDs []string,
) (map[string]int64, error) {
	result := make(map[string]int64, len(featureIDs))
	if len(featureIDs) == 0 {
		return result, nil
	}

	// $1 is environment_id; feature_id placeholders start at $2.
	placeholders := make([]string, len(featureIDs))
	args := make([]interface{}, len(featureIDs)+1)
	args[0] = environmentID
	for i, id := range featureIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}

	query := fmt.Sprintf(countCodeReferencesByFeatureIDsSQL, strings.Join(placeholders, ", "))
	rows, err := s.qe.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var featureID string
		var count int64
		if err := rows.Scan(&featureID, &count); err != nil {
			return nil, err
		}
		result[featureID] = count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func listCodeReferencesOptionsFromParams(
	p storage.ListCodeReferencesParams,
) (*pgstorage.ListOptions, error) {
	filters := []*pgstorage.Filter{
		{Column: "environment_id", Operator: pgstorage.OperatorEqual, Value: p.EnvironmentID},
		{Column: "feature_id", Operator: pgstorage.OperatorEqual, Value: p.FeatureID},
	}
	if p.RepositoryName != "" {
		filters = append(filters, &pgstorage.Filter{
			Column: "repository_name", Operator: pgstorage.OperatorEqual, Value: p.RepositoryName,
		})
	}
	if p.RepositoryOwner != "" {
		filters = append(filters, &pgstorage.Filter{
			Column: "repository_owner", Operator: pgstorage.OperatorEqual, Value: p.RepositoryOwner,
		})
	}
	if p.RepositoryType != coderefproto.CodeReference_REPOSITORY_TYPE_UNSPECIFIED {
		filters = append(filters, &pgstorage.Filter{
			Column: "repository_type", Operator: pgstorage.OperatorEqual, Value: p.RepositoryType,
		})
	}
	if p.RepositoryBranch != "" {
		filters = append(filters, &pgstorage.Filter{
			Column: "repository_branch", Operator: pgstorage.OperatorEqual, Value: p.RepositoryBranch,
		})
	}
	if p.FileExtension != "" {
		filters = append(filters, &pgstorage.Filter{
			Column: "file_extension", Operator: pgstorage.OperatorEqual, Value: p.FileExtension,
		})
	}

	var column string
	switch p.OrderBy {
	case coderefproto.ListCodeReferencesRequest_DEFAULT:
		column = "id"
	case coderefproto.ListCodeReferencesRequest_CREATED_AT:
		column = "created_at"
	case coderefproto.ListCodeReferencesRequest_UPDATED_AT:
		column = "updated_at"
	default:
		return nil, storage.ErrInvalidOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if p.OrderDirection == coderefproto.ListCodeReferencesRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}

	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, storage.ErrInvalidCursor
	}
	limit := p.PageSize
	if limit < 0 {
		limit = 0
	}
	return &pgstorage.ListOptions{
		Limit:   limit,
		Offset:  offset,
		Filters: filters,
		Orders:  []*pgstorage.Order{pgstorage.NewOrder(column, direction)},
	}, nil
}
