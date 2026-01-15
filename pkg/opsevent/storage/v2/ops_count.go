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

	"github.com/bucketeer-io/bucketeer/v2/pkg/opsevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

var (
	//go:embed sql/ops_count/insert_ops_count.sql
	insertOpsCountSQL string
	//go:embed sql/ops_count/select_ops_counts.sql
	selectOpsCountsSQL string
)

type OpsCountStorage interface {
	UpsertOpsCount(ctx context.Context, environmentId string, oc *domain.OpsCount) error
	ListOpsCounts(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*proto.OpsCount, int, error)
}

type opsCountStorage struct {
	qe mysql.QueryExecer
}

func NewOpsCountStorage(qe mysql.QueryExecer) OpsCountStorage {
	return &opsCountStorage{qe: qe}
}

func (s *opsCountStorage) UpsertOpsCount(ctx context.Context, environmentId string, oc *domain.OpsCount) error {
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
	options *mysql.ListOptions,
) ([]*proto.OpsCount, int, error) {
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectOpsCountsSQL, options)
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
