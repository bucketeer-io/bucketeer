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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	//go:embed sql/organization/insert_organization.sql
	insertOrganizationSQL string
	//go:embed sql/organization/update_organization.sql
	updateOrganizationSQL string
	//go:embed sql/organization/select_organization.sql
	selectOrganizationSQL string
	//go:embed sql/organization/select_system_admin_organization.sql
	selectSystemAdminOrganizationSQL string
	//go:embed sql/organization/select_organizations.sql
	selectOrganizationsSQL string
	//go:embed sql/organization/count_organizations.sql
	countOrganizationsSQL string
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
	GetSystemAdminOrganization(ctx context.Context) (*domain.Organization, error)
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

func (s *organizationStorage) CreateOrganization(ctx context.Context, o *domain.Organization) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertOrganizationSQL,
		o.Id,
		o.Name,
		o.UrlCode,
		o.Description,
		o.Disabled,
		o.Archived,
		o.Trial,
		o.SystemAdmin,
		o.CreatedAt,
		o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrOrganizationAlreadyExists
		}
		return err
	}
	return nil
}

func (s *organizationStorage) UpdateOrganization(ctx context.Context, o *domain.Organization) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateOrganizationSQL,
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
	err := s.qe.QueryRowContext(
		ctx,
		selectOrganizationSQL,
		id,
	).Scan(
		&organization.Id,
		&organization.Name,
		&organization.UrlCode,
		&organization.Description,
		&organization.Disabled,
		&organization.Archived,
		&organization.Trial,
		&organization.SystemAdmin,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrOrganizationNotFound
		}
		return nil, err
	}
	return &domain.Organization{Organization: &organization}, nil
}

func (s *organizationStorage) GetSystemAdminOrganization(ctx context.Context) (*domain.Organization, error) {
	organization := proto.Organization{}
	err := s.qe.QueryRowContext(
		ctx,
		selectSystemAdminOrganizationSQL,
	).Scan(
		&organization.Id,
		&organization.Name,
		&organization.UrlCode,
		&organization.Description,
		&organization.Disabled,
		&organization.Archived,
		&organization.Trial,
		&organization.SystemAdmin,
		&organization.CreatedAt,
		&organization.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
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
	query := fmt.Sprintf(selectOrganizationsSQL, whereSQL, orderBySQL, limitOffsetSQL)
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
			&organization.SystemAdmin,
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
	countQuery := fmt.Sprintf(countOrganizationsSQL, whereSQL, orderBySQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return organizations, nextOffset, totalCount, nil
}
