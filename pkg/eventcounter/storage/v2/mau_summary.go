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

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

type MAUSummaryStorage interface {
	UpsertMAUSummary(
		ctx context.Context,
		mauSummary *eventcounter.MAUSummary,
	) error
}

type mauSummaryStorage struct {
	qe mysql.QueryExecer
}

func NewMAUSummaryStorage(qe mysql.QueryExecer) MAUSummaryStorage {
	return &mauSummaryStorage{qe: qe}
}

func (s *mauSummaryStorage) UpsertMAUSummary(
	ctx context.Context,
	m *eventcounter.MAUSummary,
) error {
	query := `
		INSERT INTO mau_summary (
			yearmonth,
			source_id,
			user_count,
			request_count,
			evaluation_count,
			goal_count,
			is_all,
			is_finished,
			created_at,
			updated_at,
			environment_id
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		) ON DUPLICATE KEY UPDATE
			user_count = VALUES(user_count),
			request_count = VALUES(request_count),
			evaluation_count = VALUES(evaluation_count),
			goal_count = VALUES(goal_count),
			updated_at = VALUES(updated_at)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		m.Yearmonth,
		m.SourceId,
		m.UserCount,
		m.RequestCount,
		m.EvaluationCount,
		m.GoalCount,
		m.IsAll,
		m.IsFinished,
		m.CreatedAt,
		m.UpdatedAt,
		m.EnvironmentId,
	)
	if err != nil {
		return err
	}
	return nil
}
