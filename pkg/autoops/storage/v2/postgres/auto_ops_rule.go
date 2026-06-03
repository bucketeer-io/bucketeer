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
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

var (
	//go:embed sql/auto_ops_rule/insert_auto_ops_rule.sql
	insertAutoOpsRuleSQL string
	//go:embed sql/auto_ops_rule/update_auto_ops_rule.sql
	updateAutoOpsRuleSQL string
	//go:embed sql/auto_ops_rule/select_auto_ops_rule.sql
	selectAutoOpsRuleSQL string
	//go:embed sql/auto_ops_rule/select_auto_ops_rules.sql
	selectAutoOpsRulesSQL string
)

type autoOpsRuleStorage struct {
	qe pgstorage.QueryExecer
}

func NewAutoOpsRuleStorage(qe pgstorage.QueryExecer) v2as.AutoOpsRuleStorage {
	return &autoOpsRuleStorage{qe: qe}
}

func (s *autoOpsRuleStorage) CreateAutoOpsRule(
	ctx context.Context,
	e *domain.AutoOpsRule,
	environmentId string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertAutoOpsRuleSQL,
		e.Id,
		e.FeatureId,
		int32(e.OpsType),
		pgstorage.JSONObject{Val: e.Clauses},
		e.CreatedAt,
		e.UpdatedAt,
		e.Deleted,
		int32(e.AutoOpsStatus),
		environmentId,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2as.ErrAutoOpsRuleAlreadyExists
		}
		return err
	}
	return nil
}

func (s *autoOpsRuleStorage) UpdateAutoOpsRule(
	ctx context.Context,
	e *domain.AutoOpsRule,
	environmentId string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateAutoOpsRuleSQL,
		e.FeatureId,
		int32(e.OpsType),
		pgstorage.JSONObject{Val: e.Clauses},
		e.CreatedAt,
		e.UpdatedAt,
		e.Deleted,
		int32(e.AutoOpsStatus),
		e.Id,
		environmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return v2as.ErrAutoOpsRuleUnexpectedAffectedRows
	}
	return nil
}

func (s *autoOpsRuleStorage) GetAutoOpsRule(
	ctx context.Context,
	id, environmentId string,
) (*domain.AutoOpsRule, error) {
	autoOpsRule := proto.AutoOpsRule{}
	var opsType int32
	err := s.qe.QueryRowContext(
		ctx,
		selectAutoOpsRuleSQL,
		id,
		environmentId,
	).Scan(
		&autoOpsRule.Id,
		&autoOpsRule.FeatureId,
		&opsType,
		&pgstorage.JSONObject{Val: &autoOpsRule.Clauses},
		&autoOpsRule.CreatedAt,
		&autoOpsRule.UpdatedAt,
		&autoOpsRule.Deleted,
		&autoOpsRule.AutoOpsStatus,
		&autoOpsRule.FeatureName,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2as.ErrAutoOpsRuleNotFound
		}
		return nil, err
	}
	autoOpsRule.OpsType = proto.OpsType(opsType)
	return &domain.AutoOpsRule{AutoOpsRule: &autoOpsRule}, nil
}

func (s *autoOpsRuleStorage) ListAutoOpsRules(
	ctx context.Context,
	params v2as.ListAutoOpsRulesParams,
) ([]*proto.AutoOpsRule, int, error) {
	options, err := listAutoOpsRulesOptionsFromParams(params)
	if err != nil {
		return nil, 0, err
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectAutoOpsRulesSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	autoOpsRules := make([]*proto.AutoOpsRule, 0)
	for rows.Next() {
		autoOpsRule := proto.AutoOpsRule{}
		var opsType int32
		err := rows.Scan(
			&autoOpsRule.Id,
			&autoOpsRule.FeatureId,
			&opsType,
			&pgstorage.JSONObject{Val: &autoOpsRule.Clauses},
			&autoOpsRule.CreatedAt,
			&autoOpsRule.UpdatedAt,
			&autoOpsRule.Deleted,
			&autoOpsRule.AutoOpsStatus,
			&autoOpsRule.FeatureName,
		)
		if err != nil {
			return nil, 0, err
		}
		autoOpsRule.OpsType = proto.OpsType(opsType)
		autoOpsRules = append(autoOpsRules, &autoOpsRule)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	nextOffset := options.Offset + len(autoOpsRules)
	return autoOpsRules, nextOffset, nil
}

func listAutoOpsRulesOptionsFromParams(p v2as.ListAutoOpsRulesParams) (*pgstorage.ListOptions, error) {
	filters := []*pgstorage.Filter{
		{
			Column:   "aor.deleted",
			Operator: pgstorage.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "aor.environment_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		},
	}
	var inFilters []*pgstorage.InFilter
	if len(p.FeatureIDs) > 0 {
		values := make([]interface{}, len(p.FeatureIDs))
		for i, id := range p.FeatureIDs {
			values[i] = id
		}
		inFilters = append(inFilters, &pgstorage.InFilter{
			Column: "aor.feature_id",
			Values: values,
		})
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, v2as.ErrInvalidCursor
	}
	limit := p.PageSize
	if limit < 0 {
		limit = 0
	}
	return &pgstorage.ListOptions{
		Limit:     limit,
		Offset:    offset,
		Filters:   filters,
		InFilters: inFilters,
	}, nil
}
