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
	"fmt"
	"strconv"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	//go:embed sql/segment_user/get_segment_user.sql
	getSegmentUserSQL string
	//go:embed sql/segment_user/select_segment_users.sql
	selectSegmentUsersSQL string
)

const batchSize = 1000

type segmentUserStorage struct {
	qe pgstorage.QueryExecer
}

func NewSegmentUserStorage(qe pgstorage.QueryExecer) v2fs.SegmentUserStorage {
	return &segmentUserStorage{qe: qe}
}

func (s *segmentUserStorage) UpsertSegmentUsers(
	ctx context.Context,
	users []*proto.SegmentUser,
	environmentId string,
) error {
	for i := 0; i < len(users); i += batchSize {
		j := i + batchSize
		if j > len(users) {
			j = len(users)
		}
		if err := s.upsertSegmentUsers(ctx, users[i:j], environmentId); err != nil {
			return err
		}
	}
	return nil
}

func (s *segmentUserStorage) upsertSegmentUsers(
	ctx context.Context,
	users []*proto.SegmentUser,
	environmentId string,
) error {
	var query strings.Builder
	query.WriteString(`
		INSERT INTO segment_user (
			id,
			segment_id,
			user_id,
			state,
			deleted,
			environment_id
		) VALUES
	`)
	args := []interface{}{}
	for i, u := range users {
		if i != 0 {
			query.WriteString(",")
		}
		base := i*6 + 1
		query.WriteString(fmt.Sprintf(" ($%d, $%d, $%d, $%d, $%d, $%d)",
			base, base+1, base+2, base+3, base+4, base+5))
		args = append(
			args,
			u.Id,
			u.SegmentId,
			u.UserId,
			int32(u.State),
			u.Deleted,
			environmentId,
		)
	}
	query.WriteString(`
		ON CONFLICT (id, environment_id) DO UPDATE SET
		segment_id = EXCLUDED.segment_id,
		user_id = EXCLUDED.user_id,
		state = EXCLUDED.state,
		deleted = EXCLUDED.deleted
	`)
	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	return err
}

func (s *segmentUserStorage) GetSegmentUser(
	ctx context.Context,
	id, environmentId string,
) (*domain.SegmentUser, error) {
	segmentUser := proto.SegmentUser{}
	var state int32
	err := s.qe.QueryRowContext(
		ctx,
		getSegmentUserSQL,
		id,
		environmentId,
	).Scan(
		&segmentUser.Id,
		&segmentUser.SegmentId,
		&segmentUser.UserId,
		&state,
		&segmentUser.Deleted,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2fs.ErrSegmentUserNotFound
		}
		return nil, err
	}
	segmentUser.State = proto.SegmentUser_State(state)
	return &domain.SegmentUser{SegmentUser: &segmentUser}, nil
}

func (s *segmentUserStorage) ListSegmentUsers(
	ctx context.Context,
	p v2fs.ListSegmentUsersParams,
) ([]*proto.SegmentUser, int, error) {
	filters := []*pgstorage.Filter{
		{Column: "segment_id", Operator: pgstorage.OperatorEqual, Value: p.SegmentID},
		{Column: "deleted", Operator: pgstorage.OperatorEqual, Value: false},
		{Column: "environment_id", Operator: pgstorage.OperatorEqual, Value: p.EnvironmentID},
	}
	if p.State != nil {
		filters = append(filters, &pgstorage.Filter{
			Column: "state", Operator: pgstorage.OperatorEqual, Value: *p.State,
		})
	}
	if p.UserID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column: "user_id", Operator: pgstorage.OperatorEqual, Value: p.UserID,
		})
	}
	limit := p.PageSize
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, 0, v2fs.ErrInvalidListSegmentUsersCursor
	}
	whereParts := make([]pgstorage.WherePart, 0, len(filters))
	for _, f := range filters {
		whereParts = append(whereParts, f)
	}
	whereSQL, whereArgs := pgstorage.ConstructWhereSQLString(whereParts)
	limitOffsetSQL := pgstorage.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectSegmentUsersSQL, whereSQL, limitOffsetSQL)
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
		return nil, 0, rows.Err()
	}
	nextOffset := offset + len(segmentUsers)
	return segmentUsers, nextOffset, nil
}
