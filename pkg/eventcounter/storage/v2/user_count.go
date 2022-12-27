// Copyright 2022 The Bucketeer Authors.
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
)

type UserCountStorage interface {
	GetMAUCount(
		ctx context.Context,
		environmentNamespace, yearMonth string,
	) (int64, int64, error)
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
			SUM(event_count) as event_count
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
