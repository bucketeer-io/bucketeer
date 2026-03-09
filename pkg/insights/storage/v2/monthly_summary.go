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
	"fmt"
	"strings"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
)

var (
	//go:embed sql/upsert_monthly_summary.sql
	upsertMonthlySummarySQL string
	//go:embed sql/select_monthly_summary.sql
	selectMonthlySummarySQL string
)

type MonthlySummaryRecord struct {
	Yearmonth     string
	EnvironmentID string
	SourceID      string
	MAU           int64
	Requests      int64
}

type ListMonthlySummaryResult struct {
	Yearmonth       string
	EnvironmentID   string
	EnvironmentName string
	ProjectName     string
	SourceID        string
	MAU             int64
	Requests        int64
}

type MonthlySummaryStorage interface {
	UpsertMonthlySummaryBatch(ctx context.Context, records []MonthlySummaryRecord) error
	ListMonthlySummaries(ctx context.Context, environmentIDs, sourceIDs []string) ([]ListMonthlySummaryResult, error)
}

type monthlySummaryStorage struct {
	qe mysql.QueryExecer
}

func NewMonthlySummaryStorage(qe mysql.QueryExecer) MonthlySummaryStorage {
	return &monthlySummaryStorage{qe: qe}
}

func (s *monthlySummaryStorage) UpsertMonthlySummaryBatch(
	ctx context.Context,
	records []MonthlySummaryRecord,
) error {
	if len(records) == 0 {
		return nil
	}

	now := time.Now().Unix()

	placeholders := make([]string, len(records))
	args := make([]any, 0, len(records)*7)

	for i, rec := range records {
		placeholders[i] = "(?, ?, ?, ?, ?, ?, ?)"
		args = append(args, rec.EnvironmentID, rec.SourceID, rec.Yearmonth, rec.MAU, rec.Requests, now, now)
	}

	query := fmt.Sprintf(upsertMonthlySummarySQL, strings.Join(placeholders, ", "))

	_, err := s.qe.ExecContext(ctx, query, args...)
	return err
}

func (s *monthlySummaryStorage) ListMonthlySummaries(
	ctx context.Context,
	environmentIDs, sourceIDs []string,
) ([]ListMonthlySummaryResult, error) {
	if len(environmentIDs) == 0 || len(sourceIDs) == 0 {
		return nil, nil
	}

	envPlaceholders := make([]string, len(environmentIDs))
	srcPlaceholders := make([]string, len(sourceIDs))
	args := make([]any, 0, len(environmentIDs)+len(sourceIDs))

	for i, id := range environmentIDs {
		envPlaceholders[i] = "?"
		args = append(args, id)
	}
	for i, id := range sourceIDs {
		srcPlaceholders[i] = "?"
		args = append(args, id)
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

	var results []ListMonthlySummaryResult
	for rows.Next() {
		var row ListMonthlySummaryResult
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
