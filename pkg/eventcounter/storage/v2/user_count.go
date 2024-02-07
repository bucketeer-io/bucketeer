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

	proto "github.com/bucketeer-io/bucketeer/proto/eventcounter"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

type UserCountStorage interface {
	GetMAUCount(
		ctx context.Context,
		environmentNamespace, yearMonth string,
	) (int64, int64, error)
	GetMAUCounts(
		ctx context.Context,
		yearMonth string,
	) ([]*proto.MAUSummary, error)
	GetMAUCountsGroupBySourceID(
		ctx context.Context,
		yearMonth string,
	) ([]*proto.MAUSummary, error)
}

type userCountStorage struct {
	qe mysql.QueryExecer
}

func NewUserCountStorage(qe mysql.QueryExecer) UserCountStorage {
	return &userCountStorage{qe: qe}
}

func (s *userCountStorage) GetMAUCount(
	ctx context.Context,
	environmentNamespace, yearMonth string,
) (int64, int64, error) {
	query := `
		SELECT
			count(*) as user_count,
			IFNULL(SUM(event_count), 0) as event_count
		FROM
			mau
		WHERE
			environment_namespace = ? AND
			yearmonth = ?
	`
	var userCount, eventCount int64
	err := s.qe.QueryRowContext(
		ctx,
		query,
		environmentNamespace,
		yearMonth,
	).Scan(
		&userCount,
		&eventCount,
	)
	if err != nil {
		return 0, 0, err
	}
	return userCount, eventCount, nil
}

func (s *userCountStorage) GetMAUCounts(
	ctx context.Context,
	yearMonth string,
) ([]*proto.MAUSummary, error) {
	summaries := make([]*proto.MAUSummary, 0)
	query := `
		SELECT
			environment_namespace as environment_id,
			count(*) as user_count,
			IFNULL(SUM(event_count), 0) as request_count
		FROM
			mau
		WHERE
			yearmonth = ?
		GROUP BY
			environment_namespace
	`
	rows, err := s.qe.QueryContext(
		ctx,
		query,
		yearMonth,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		summary := proto.MAUSummary{}
		err := rows.Scan(
			&summary.EnvironmentId,
			&summary.UserCount,
			&summary.RequestCount,
		)
		if err != nil {
			return nil, err
		}
		summary.Yearmonth = yearMonth
		summary.IsAll = true
		summaries = append(summaries, &summary)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return summaries, nil
}

func (s *userCountStorage) GetMAUCountsGroupBySourceID(
	ctx context.Context,
	yearMonth string,
) ([]*proto.MAUSummary, error) {
	summaries := make([]*proto.MAUSummary, 0)
	query := `
		SELECT
			environment_namespace as environment_id,
			source_id,
			count(*) as user_count,
			IFNULL(SUM(event_count), 0) as request_count
		FROM
			mau
		WHERE
			yearmonth = ?
		GROUP BY
			environment_namespace,
			source_id
	`
	rows, err := s.qe.QueryContext(
		ctx,
		query,
		yearMonth,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		summary := proto.MAUSummary{}
		err := rows.Scan(
			&summary.EnvironmentId,
			&summary.SourceId,
			&summary.UserCount,
			&summary.RequestCount,
		)
		if err != nil {
			return nil, err
		}
		summary.Yearmonth = yearMonth
		summary.IsAll = false
		summaries = append(summaries, &summary)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return summaries, nil
}
