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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	ErrOrganizationAlreadyExists          = errors.New("organization: already exists")
	ErrOrganizationNotFound               = errors.New("organization: not found")
	ErrOrganizationUnexpectedAffectedRows = errors.New("organization: unexpected affected rows")
)

type OrganizationStorage interface {
	CreateOrganization(ctx context.Context, p *domain.Organization) error
	UpdateOrganization(ctx context.Context, p *domain.Organization) error
	GetOrganization(ctx context.Context, id string) (*domain.Organization, error)
	ListOrganizations(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Organization, int, int64, error)
}

type organizationStorage struct {
	qe mysql.QueryExecer
}

func NewOrganizationStorage(qe mysql.QueryExecer) OrganizationStorage {
	return &organizationStorage{qe}
}

func (s *organizationStorage) CreateOrganization(ctx context.Context, p *domain.Organization) error {
	query := `
		INSERT INTO organization (
			id,
			name,
			url_code,
			description,
			disabled,
			archived,
			trial,
			created_at,
			updated_at
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		p.Id,
		p.Name,
		p.UrlCode,
		p.Description,
		p.Disabled,
		p.Archived,
		p.Trial,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrOrganizationAlreadyExists
		}
		return err
	}
	return nil
}

func (s *organizationStorage) UpdateOrganization(ctx context.Context, o *domain.Organization) error {
	query := `
		UPDATE 
			organization
		SET
			name = ?,
			description = ?,
			disabled = ?,
			archived = ?,
			trial = ?,
			created_at = ?,
			updated_at = ?
		WHERE
			id = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
		o.Name,
		o.Description,
		o.Disabled,
		o.Archived,
		o.Trial,
		o.CreatedAt,
		o.UpdatedAt,
		o.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrOrganizationUnexpectedAffectedRows
	}
	return nil
}

func (s *organizationStorage) GetOrganization(ctx context.Context, id string) (*domain.Organization, error) {
	organization := proto.Organization{}
	query := `
		SELECT
			id,
			name,
			url_code,
			description,
			disabled,
			archived,
			trial,
			created_at,
			updated_at
		FROM
			organization
		WHERE
			id = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&organization.Id,
		&organization.Name,
		&organization.UrlCode,
		&organization.Description,
		&organization.Disabled,
		&organization.Archived,
		&organization.Trial,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrOrganizationNotFound
		}
		return nil, err
	}
	return &domain.Organization{Organization: &organization}, nil
}

func (s *organizationStorage) ListOrganizations(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Organization, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			url_code,
			description,
			disabled,
			archived,
			trial,
			created_at,
			updated_at
		FROM
			organization
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	organizations := make([]*proto.Organization, 0, limit)
	for rows.Next() {
		organization := proto.Organization{}
		err := rows.Scan(
			&organization.Id,
			&organization.Name,
			&organization.UrlCode,
			&organization.Description,
			&organization.Disabled,
			&organization.Archived,
			&organization.Trial,
			&organization.CreatedAt,
			&organization.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		organizations = append(organizations, &organization)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(organizations)
	var totalCount int64
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(1)
		FROM
			organization
		%s %s
		`, whereSQL, orderBySQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return organizations, nextOffset, totalCount, nil
}
