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
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrSegmentUserNotFound = errors.New("segmentUser: not found")
)

type SegmentUserStorage interface {
	UpsertSegmentUsers(ctx context.Context, users []*proto.SegmentUser, environmentNamespace string) error
	GetSegmentUser(ctx context.Context, id, environmentNamespace string) (*domain.SegmentUser, error)
	ListSegmentUsers(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.SegmentUser, int, error)
}

type segmentUserStorage struct {
	qe mysql.QueryExecer
}

func NewSegmentUserStorage(qe mysql.QueryExecer) SegmentUserStorage {
	return &segmentUserStorage{qe: qe}
}

func (s *segmentUserStorage) UpsertSegmentUsers(
	ctx context.Context,
	users []*proto.SegmentUser,
	environmentNamespace string,
) error {
	for _, u := range users {
		query := `
			INSERT INTO segment_user (
				id,
				segment_id,
				user_id,
				state,
				deleted,
				environment_namespace
			) VALUES (
				?, ?, ?, ?, ?, ?
			) ON DUPLICATE KEY UPDATE
				segment_id = VALUES(segment_id),
				user_id = VALUES(user_id),
				state = VALUES(state),
				deleted = VALUES(deleted)
		`
		_, err := s.qe.ExecContext(
			ctx,
			query,
			u.Id,
			u.SegmentId,
			u.UserId,
			int32(u.State),
			u.Deleted,
			environmentNamespace,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *segmentUserStorage) GetSegmentUser(
	ctx context.Context,
	id, environmentNamespace string,
) (*domain.SegmentUser, error) {
	segmentUser := proto.SegmentUser{}
	var state int32
	query := `
		SELECT
			id,
			segment_id,
			user_id,
			state,
			deleted
		FROM
			segment_user
		WHERE
			id = ? AND
			environment_namespace = ?
	`
	err := s.qe.QueryRowContext(
		ctx,
		query,
		id,
		environmentNamespace,
	).Scan(
		&segmentUser.Id,
		&segmentUser.SegmentId,
		&segmentUser.UserId,
		&state,
		&segmentUser.Deleted,
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrSegmentUserNotFound
		}
		return nil, err
	}
	segmentUser.State = proto.SegmentUser_State(state)
	return &domain.SegmentUser{SegmentUser: &segmentUser}, nil
}

func (s *segmentUserStorage) ListSegmentUsers(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.SegmentUser, int, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(`
		SELECT
			id,
			segment_id,
			user_id,
			state,
			deleted
		FROM
			segment_user
		%s %s %s
		`, whereSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	segmentUsers := make([]*proto.SegmentUser, 0, limit)
	for rows.Next() {
		segmentUser := proto.SegmentUser{}
		var state int32
		err := rows.Scan(
			&segmentUser.Id,
			&segmentUser.SegmentId,
			&segmentUser.UserId,
			&state,
			&segmentUser.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		segmentUser.State = proto.SegmentUser_State(state)
		segmentUsers = append(segmentUsers, &segmentUser)
	}
	if rows.Err() != nil {
		return nil, 0, err
	}
	nextOffset := offset + len(segmentUsers)
	return segmentUsers, nextOffset, nil
}
