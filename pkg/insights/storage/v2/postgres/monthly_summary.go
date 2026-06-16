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
	"fmt"
	"strings"
	"time"

	v2is "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
)

var (
	//go:embed sql/upsert_monthly_summary.sql
	upsertMonthlySummarySQL string
	//go:embed sql/select_monthly_summary.sql
	selectMonthlySummarySQL string
)

const monthlySummaryColumns = 7

type monthlySummaryStorage struct {
	qe postgres.QueryExecer
}

func NewMonthlySummaryStorage(qe postgres.QueryExecer) v2is.MonthlySummaryStorage {
	return &monthlySummaryStorage{qe: qe}
}

func (s *monthlySummaryStorage) UpsertMonthlySummaryBatch(
	ctx context.Context,
	records []v2is.MonthlySummaryRecord,
) error {
	if len(records) == 0 {
		return nil
	}

	now := time.Now().Unix()

	placeholders := make([]string, len(records))
	args := make([]any, 0, len(records)*monthlySummaryColumns)

	for i, rec := range records {
		base := i * monthlySummaryColumns
		placeholders[i] = fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7,
		)
		args = append(args, rec.EnvironmentID, rec.SourceID, rec.Yearmonth, rec.MAU, rec.Requests, now, now)
	}

	query := fmt.Sprintf(upsertMonthlySummarySQL, strings.Join(placeholders, ", "))

	_, err := s.qe.ExecContext(ctx, query, args...)
	return err
}

func (s *monthlySummaryStorage) ListMonthlySummaries(
	ctx context.Context,
	environmentIDs, sourceIDs []string,
) ([]v2is.ListMonthlySummaryResult, error) {
	if len(environmentIDs) == 0 || len(sourceIDs) == 0 {
		return nil, nil
	}

	envPlaceholders := make([]string, len(environmentIDs))
	srcPlaceholders := make([]string, len(sourceIDs))
	args := make([]any, 0, len(environmentIDs)+len(sourceIDs))

	pos := 1
	for i, id := range environmentIDs {
		envPlaceholders[i] = fmt.Sprintf("$%d", pos)
		args = append(args, id)
		pos++
	}
	for i, id := range sourceIDs {
		srcPlaceholders[i] = fmt.Sprintf("$%d", pos)
		args = append(args, id)
		pos++
	}

	query := fmt.Sprintf(
		selectMonthlySummarySQL,
		strings.Join(envPlaceholders, ","),
		strings.Join(srcPlaceholders, ","),
	)

	rows, err := s.qe.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []v2is.ListMonthlySummaryResult
	for rows.Next() {
		var row v2is.ListMonthlySummaryResult
		if err := rows.Scan(
			&row.EnvironmentID,
			&row.EnvironmentName,
			&row.ProjectName,
			&row.SourceID,
			&row.Yearmonth,
			&row.MAU,
			&row.Requests,
		); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
