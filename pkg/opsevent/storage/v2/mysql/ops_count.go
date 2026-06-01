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

package mysql

import (
	"context"
	_ "embed"
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/domain"
	v2os "github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

var (
	//go:embed sql/ops_count/insert_ops_count.sql
	insertOpsCountSQL string
	//go:embed sql/ops_count/select_ops_counts.sql
	selectOpsCountsSQL string
)

type opsCountStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewOpsCountStorage(qe mysqlstorage.QueryExecer) v2os.OpsCountStorage {
	return &opsCountStorage{qe: qe}
}

func (s *opsCountStorage) UpsertOpsCount(
	ctx context.Context,
	environmentId string,
	oc *domain.OpsCount,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertOpsCountSQL,
		oc.Id,
		oc.AutoOpsRuleId,
		oc.ClauseId,
		oc.UpdatedAt,
		oc.OpsEventCount,
		oc.EvaluationCount,
		oc.FeatureId,
		environmentId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *opsCountStorage) ListOpsCounts(
	ctx context.Context,
	params v2os.ListOpsCountsParams,
) ([]*proto.OpsCount, int, error) {
	options, err := listOpsCountsOptionsFromParams(params)
	if err != nil {
		return nil, 0, err
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectOpsCountsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	opsCounts := make([]*proto.OpsCount, 0, options.Limit)
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
		return nil, 0, rows.Err()
	}
	nextOffset := options.Offset + len(opsCounts)
	return opsCounts, nextOffset, nil
}

func listOpsCountsOptionsFromParams(p v2os.ListOpsCountsParams) (*mysqlstorage.ListOptions, error) {
	filters := []*mysqlstorage.FilterV2{
		{
			Column:   "environment_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		},
	}

	var inFilters []*mysqlstorage.InFilter
	if len(p.FeatureIDs) > 0 {
		values := make([]interface{}, len(p.FeatureIDs))
		for i, id := range p.FeatureIDs {
			values[i] = id
		}
		inFilters = append(inFilters, &mysqlstorage.InFilter{
			Column: "feature_id",
			Values: values,
		})
	}
	if len(p.AutoOpsRuleIDs) > 0 {
		values := make([]interface{}, len(p.AutoOpsRuleIDs))
		for i, id := range p.AutoOpsRuleIDs {
			values[i] = id
		}
		inFilters = append(inFilters, &mysqlstorage.InFilter{
			Column: "auto_ops_rule_id",
			Values: values,
		})
	}

	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, v2os.ErrInvalidCursor
	}

	return &mysqlstorage.ListOptions{
		Limit:     p.PageSize,
		Offset:    offset,
		Filters:   filters,
		InFilters: inFilters,
	}, nil
}
