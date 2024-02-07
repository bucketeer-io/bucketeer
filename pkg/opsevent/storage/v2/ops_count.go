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
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/opsevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

type OpsCountStorage interface {
	UpsertOpsCount(ctx context.Context, environmentNamespace string, oc *domain.OpsCount) error
	ListOpsCounts(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.OpsCount, int, error)
}

type opsCountStorage struct {
	qe mysql.QueryExecer
}

func NewOpsCountStorage(qe mysql.QueryExecer) OpsCountStorage {
	return &opsCountStorage{qe: qe}
}

func (s *opsCountStorage) UpsertOpsCount(ctx context.Context, environmentNamespace string, oc *domain.OpsCount) error {
	query := `
		INSERT INTO ops_count (
			id,
			auto_ops_rule_id,
			clause_id,
			updated_at,
			ops_event_count,
			evaluation_count,
			feature_id,
			environment_namespace
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		) ON DUPLICATE KEY UPDATE
			auto_ops_rule_id = VALUES(auto_ops_rule_id),
			clause_id = VALUES(clause_id),
			updated_at = VALUES(updated_at),
			ops_event_count = VALUES(ops_event_count),
			evaluation_count = VALUES(evaluation_count),
			feature_id = VALUES(feature_id)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		oc.Id,
		oc.AutoOpsRuleId,
		oc.ClauseId,
		oc.UpdatedAt,
		oc.OpsEventCount,
		oc.EvaluationCount,
		oc.FeatureId,
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *opsCountStorage) ListOpsCounts(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.OpsCount, int, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			auto_ops_rule_id,
			clause_id,
			updated_at,
			ops_event_count,
			evaluation_count,
			feature_id
		FROM
			ops_count
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	opsCounts := make([]*proto.OpsCount, 0, limit)
	for rows.Next() {
		opsCount := proto.OpsCount{}
		err := rows.Scan(
			&opsCount.Id,
			&opsCount.AutoOpsRuleId,
			&opsCount.ClauseId,
			&opsCount.UpdatedAt,
			&opsCount.OpsEventCount,
			&opsCount.EvaluationCount,
			&opsCount.FeatureId,
		)
		if err != nil {
			return nil, 0, err
		}
		opsCounts = append(opsCounts, &opsCount)
	}
	if rows.Err() != nil {
		return nil, 0, err
	}
	nextOffset := offset + len(opsCounts)
	return opsCounts, nextOffset, nil
}
