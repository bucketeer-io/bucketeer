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

package storage

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

var (
	ErrDomainPolicyAlreadyExists          = errors.New("domain_policy: already exists")
	ErrDomainPolicyNotFound               = errors.New("domain_policy: not found")
	ErrDomainPolicyUnexpectedAffectedRows = errors.New("domain_policy: unexpected affected rows")
)

var (
	//go:embed sql/domain_policy/insert_domain_policy.sql
	insertDomainPolicySQL string
	//go:embed sql/domain_policy/select_domain_policy.sql
	selectDomainPolicySQL string
	//go:embed sql/domain_policy/update_domain_policy.sql
	updateDomainPolicySQL string
	//go:embed sql/domain_policy/delete_domain_policy.sql
	deleteDomainPolicySQL string
	//go:embed sql/domain_policy/list_domain_policies.sql
	listDomainPoliciesSQL string
	//go:embed sql/domain_policy/count_domain_policies.sql
	countDomainPoliciesSQL string
)

type domainPolicyStorage struct {
	qe mysql.QueryExecer
}

// NewDomainPolicyStorage creates a new domain policy storage instance
func NewDomainPolicyStorage(qe mysql.QueryExecer) DomainPolicyStorage {
	return &domainPolicyStorage{qe: qe}
}

func (s *domainPolicyStorage) CreateDomainPolicy(ctx context.Context, policy *domain.DomainAuthPolicy) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertDomainPolicySQL,
		policy.Domain,
		mysql.JSONObject{Val: policy.AuthPolicy},
		policy.Enabled,
		policy.CreatedAt,
		policy.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrDomainPolicyAlreadyExists
		}
		return err
	}
	return nil
}

func (s *domainPolicyStorage) GetDomainPolicy(
	ctx context.Context,
	domainName string,
) (*domain.DomainAuthPolicy, error) {
	policy := authproto.DomainAuthPolicy{}
	err := s.qe.QueryRowContext(
		ctx,
		selectDomainPolicySQL,
		domainName,
	).Scan(
		&policy.Domain,
		&mysql.JSONObject{Val: &policy.AuthPolicy},
		&policy.Enabled,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrDomainPolicyNotFound
		}
		return nil, err
	}

	return domain.NewDomainAuthPolicyFromProto(&policy), nil
}

func (s *domainPolicyStorage) UpdateDomainPolicy(ctx context.Context, policy *domain.DomainAuthPolicy) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateDomainPolicySQL,
		mysql.JSONObject{Val: policy.AuthPolicy},
		policy.Enabled,
		policy.UpdatedAt,
		policy.Domain,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrDomainPolicyUnexpectedAffectedRows
	}
	return nil
}

func (s *domainPolicyStorage) DeleteDomainPolicy(ctx context.Context, domainName string) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteDomainPolicySQL,
		domainName,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrDomainPolicyUnexpectedAffectedRows
	}
	return nil
}

func (s *domainPolicyStorage) ListDomainPolicies(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*authproto.DomainAuthPolicy, int, int64, error) {
	var query string
	var whereArgs []any
	if options != nil {
		var whereSQL string
		whereParts := options.CreateWhereParts()
		whereSQL, whereArgs = mysql.ConstructWhereSQLString(whereParts)
		orderBySQL := mysql.ConstructOrderBySQLString(options.Orders)
		limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
		query = fmt.Sprintf(listDomainPoliciesSQL, whereSQL, orderBySQL, limitOffsetSQL)
	} else {
		query = listDomainPoliciesSQL
		whereArgs = []interface{}{}
	}

	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}

	policies := make([]*authproto.DomainAuthPolicy, 0, limit)
	for rows.Next() {
		policy := authproto.DomainAuthPolicy{}
		err := rows.Scan(
			&policy.Domain,
			&mysql.JSONObject{Val: &policy.AuthPolicy},
			&policy.Enabled,
			&policy.CreatedAt,
			&policy.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		policies = append(policies, &policy)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	nextOffset := offset + len(policies)

	var totalCount int64
	err = s.qe.QueryRowContext(ctx, countDomainPoliciesSQL).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}

	return policies, nextOffset, totalCount, nil
}
