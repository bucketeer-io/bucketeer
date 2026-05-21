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
	"database/sql"
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
	//go:embed sql/segment/select_segments.sql
	selectSegmentsSQL string
	//go:embed sql/segment/count_segments.sql
	countSegmentsSQL string
	//go:embed sql/segment/get_segment.sql
	getSegmentSQLQuery string
	//go:embed sql/segment/update_segment.sql
	updateSegmentSQLQuery string
	//go:embed sql/segment/insert_segment.sql
	insertSegmentSQLQuery string
	//go:embed sql/segment/delete_segment.sql
	deleteSegmentSQLQuery string
	//go:embed sql/segment/select_all_in_use_segments.sql
	selectAllInUseSegmentsSQLQuery string
	//go:embed sql/segment/select_segment_users_by_segment.sql
	selectSegmentUsersBySegmentSQLQuery string
)

type segmentStorage struct {
	qe pgstorage.QueryExecer
}

func NewSegmentStorage(qe pgstorage.QueryExecer) v2fs.SegmentStorage {
	return &segmentStorage{qe: qe}
}

func (s *segmentStorage) CreateSegment(
	ctx context.Context,
	segment *domain.Segment,
	environmentId string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertSegmentSQLQuery,
		segment.Id,
		segment.Name,
		segment.Description,
		pgstorage.JSONObject{Val: segment.Rules},
		segment.CreatedAt,
		segment.UpdatedAt,
		segment.Version,
		segment.Deleted,
		segment.IncludedUserCount,
		segment.ExcludedUserCount,
		int32(segment.Status),
		environmentId,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2fs.ErrSegmentAlreadyExists
		}
		return err
	}
	return nil
}

func (s *segmentStorage) UpdateSegment(
	ctx context.Context,
	segment *domain.Segment,
	environmentId string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateSegmentSQLQuery,
		segment.Name,
		segment.Description,
		pgstorage.JSONObject{Val: segment.Rules},
		segment.CreatedAt,
		segment.UpdatedAt,
		segment.Version,
		segment.Deleted,
		segment.IncludedUserCount,
		segment.ExcludedUserCount,
		int32(segment.Status),
		segment.Id,
		environmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return v2fs.ErrSegmentUnexpectedAffectedRows
	}
	return nil
}

func (s *segmentStorage) GetSegment(
	ctx context.Context,
	id, environmentId string,
) (*domain.Segment, []string, error) {
	segment := proto.Segment{}
	var status int32
	featureIDs := new(sql.NullString)
	err := s.qe.QueryRowContext(
		ctx,
		getSegmentSQLQuery,
		environmentId,
		id,
		environmentId,
	).Scan(
		&segment.Id,
		&segment.Name,
		&segment.Description,
		&pgstorage.JSONObject{Val: &segment.Rules},
		&segment.CreatedAt,
		&segment.UpdatedAt,
		&segment.Version,
		&segment.Deleted,
		&segment.IncludedUserCount,
		&segment.ExcludedUserCount,
		&status,
		featureIDs,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, nil, v2fs.ErrSegmentNotFound
		}
		return nil, nil, err
	}
	array := []string{}
	if featureIDs.Valid {
		segment.IsInUseStatus = true
		array = strings.Split(featureIDs.String, ",")
	}
	segment.Status = proto.Segment_Status(status)
	return &domain.Segment{Segment: &segment}, array, nil
}

func listSegmentsOrders(
	orderBy proto.ListSegmentsRequest_OrderBy,
	orderDirection proto.ListSegmentsRequest_OrderDirection,
) ([]*pgstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListSegmentsRequest_DEFAULT,
		proto.ListSegmentsRequest_NAME:
		column = "name"
	case proto.ListSegmentsRequest_CREATED_AT:
		column = "created_at"
	case proto.ListSegmentsRequest_UPDATED_AT:
		column = "updated_at"
	case proto.ListSegmentsRequest_USERS:
		column = "included_user_count"
	case proto.ListSegmentsRequest_CONNECTIONS:
		column = "feature_ids"
	default:
		return nil, v2fs.ErrInvalidListSegmentsOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListSegmentsRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}

func (s *segmentStorage) ListSegments(
	ctx context.Context,
	p v2fs.ListSegmentsParams,
) ([]*proto.Segment, int, int64, map[string][]string, error) {
	orders, err := listSegmentsOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	filters := []*pgstorage.Filter{
		{
			Column:   "seg.deleted",
			Operator: pgstorage.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "seg.environment_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		},
	}
	if p.Status != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "seg.status",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Status,
		})
	}
	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{"seg.name", "seg.description"},
			Keyword: p.SearchKeyword,
		}
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, 0, 0, nil, v2fs.ErrInvalidListSegmentsCursor
	}
	options := &pgstorage.ListOptions{
		Limit:       int(p.PageSize),
		Offset:      offset,
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      orders,
	}
	var whereSQL, orderBySQL, limitOffsetSQL string
	var whereArgs []any
	whereParts := options.CreateWhereParts()
	whereSQL, whereArgs = pgstorage.ConstructWhereSQLString(whereParts)
	orderBySQL = pgstorage.ConstructOrderBySQLString(options.Orders)
	limitOffsetSQL = pgstorage.ConstructLimitOffsetSQLString(options.Limit, options.Offset)

	var isInUseStatusSQL string
	if p.IsInUseStatus != nil {
		if *p.IsInUseStatus {
			isInUseStatusSQL = "WHERE feature_ids IS NOT NULL"
		} else {
			isInUseStatusSQL = "WHERE feature_ids IS NULL"
		}
	}
	query := fmt.Sprintf(selectSegmentsSQL, whereSQL, isInUseStatusSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	defer rows.Close()
	segments := make([]*proto.Segment, 0, options.Limit)
	featureIDsMap := map[string][]string{}
	for rows.Next() {
		segment := proto.Segment{}
		var status int32
		featureIDs := new(sql.NullString)
		err := rows.Scan(
			&segment.Id,
			&segment.Name,
			&segment.Description,
			&pgstorage.JSONObject{Val: &segment.Rules},
			&segment.CreatedAt,
			&segment.UpdatedAt,
			&segment.Version,
			&segment.Deleted,
			&segment.IncludedUserCount,
			&segment.ExcludedUserCount,
			&status,
			featureIDs,
		)
		if err != nil {
			return nil, 0, 0, nil, err
		}
		array := []string{}
		if featureIDs.Valid {
			segment.IsInUseStatus = true
			array = strings.Split(featureIDs.String, ",")
		}
		featureIDsMap[segment.Id] = array
		segment.Status = proto.Segment_Status(status)
		segments = append(segments, &segment)
	}
	if rows.Err() != nil {
		return nil, 0, 0, nil, rows.Err()
	}
	nextOffset := offset + len(segments)
	var totalCount int64
	countConditionSQL := "> 0 THEN 1 ELSE 1"
	if p.IsInUseStatus != nil {
		if *p.IsInUseStatus {
			countConditionSQL = "> 0 THEN 1 ELSE NULL"
		} else {
			countConditionSQL = "> 0 THEN NULL ELSE 1"
		}
	}
	countQuery := fmt.Sprintf(countSegmentsSQL, countConditionSQL, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	return segments, nextOffset, totalCount, featureIDsMap, nil
}

func (s *segmentStorage) DeleteSegment(ctx context.Context, id string) error {
	result, err := s.qe.ExecContext(ctx, deleteSegmentSQLQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return v2fs.ErrSegmentUnexpectedAffectedRows
	}
	return nil
}

func (s *segmentStorage) ListAllInUseSegments(
	ctx context.Context,
) ([]*v2fs.InUseSegment, error) {
	rows, err := s.qe.QueryContext(ctx, selectAllInUseSegmentsSQLQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	segments := make([]*v2fs.InUseSegment, 0)
	for rows.Next() {
		var seg v2fs.InUseSegment
		err := rows.Scan(
			&seg.SegmentID,
			&seg.EnvironmentID,
			&seg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		segments = append(segments, &seg)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return segments, nil
}

func (s *segmentStorage) ListSegmentUsersBySegment(
	ctx context.Context,
	segmentID, environmentID string,
) ([]*proto.SegmentUser, error) {
	rows, err := s.qe.QueryContext(ctx, selectSegmentUsersBySegmentSQLQuery, segmentID, environmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*proto.SegmentUser, 0)
	for rows.Next() {
		var user proto.SegmentUser
		var state int32
		err := rows.Scan(
			&user.Id,
			&user.SegmentId,
			&user.UserId,
			&state,
			&user.Deleted,
		)
		if err != nil {
			return nil, err
		}
		user.State = proto.SegmentUser_State(state)
		users = append(users, &user)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return users, nil
}
