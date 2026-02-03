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
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	ErrSegmentAlreadyExists          = pkgErr.NewErrorAlreadyExists(pkgErr.FeaturePackageName, "segment already exists")
	ErrSegmentNotFound               = pkgErr.NewErrorNotFound(pkgErr.FeaturePackageName, "segment not found", "segment")
	ErrSegmentUnexpectedAffectedRows = pkgErr.NewErrorUnexpectedAffectedRows(
		pkgErr.FeaturePackageName,
		"segment unexpected affected rows",
	)

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

type SegmentStorage interface {
	CreateSegment(ctx context.Context, segment *domain.Segment, environmentId string) error
	UpdateSegment(ctx context.Context, segment *domain.Segment, environmentId string) error
	GetSegment(ctx context.Context, id, environmentId string) (*domain.Segment, []string, error)
	ListSegments(
		ctx context.Context,
		options *mysql.ListOptions,
		isInUseStatus *bool,
	) ([]*proto.Segment, int, int64, map[string][]string, error)
	DeleteSegment(ctx context.Context, id string) error
	// ListAllInUseSegments lists all segments that are in use (referenced by feature flags).
	// Returns lightweight segment info (id, environment_id, updated_at).
	ListAllInUseSegments(ctx context.Context) ([]*InUseSegment, error)
	// ListSegmentUsersBySegment lists all users for a specific segment.
	// This is called per-segment to avoid loading all users in a single query.
	ListSegmentUsersBySegment(ctx context.Context, segmentID, environmentID string) ([]*proto.SegmentUser, error)
}

// InUseSegment represents a segment that is in use by feature flags.
type InUseSegment struct {
	SegmentID     string
	EnvironmentID string
	UpdatedAt     int64
}

type segmentStorage struct {
	qe mysql.QueryExecer
}

func NewSegmentStorage(qe mysql.QueryExecer) SegmentStorage {
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
		mysql.JSONObject{Val: segment.Rules},
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
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrSegmentAlreadyExists
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
		mysql.JSONObject{Val: segment.Rules},
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
		return ErrSegmentUnexpectedAffectedRows
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
		&mysql.JSONObject{Val: &segment.Rules},
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
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, nil, ErrSegmentNotFound
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

func (s *segmentStorage) ListSegments(
	ctx context.Context,
	options *mysql.ListOptions,
	isInUseStatus *bool,
) ([]*proto.Segment, int, int64, map[string][]string, error) {
	// Because select_segments.sql defines the variable strings in a complex constructed way,
	//  we do not use ConstructQueryAndWhereArgs() here.
	var whereSQL, orderBySQL, limitOffsetSQL string
	var whereArgs []any
	if options != nil {
		whereParts := options.CreateWhereParts()
		whereSQL, whereArgs = mysql.ConstructWhereSQLString(whereParts)
		orderBySQL = mysql.ConstructOrderBySQLString(options.Orders)
		limitOffsetSQL = mysql.ConstructLimitOffsetSQLString(options.Limit, options.Offset)
	} else {
		whereArgs = []interface{}{}
	}
	var isInUseStatusSQL string
	if isInUseStatus != nil {
		if *isInUseStatus {
			isInUseStatusSQL = "HAVING feature_ids IS NOT NULL"
		} else {
			isInUseStatusSQL = "HAVING feature_ids IS NULL"
		}
	}
	query := fmt.Sprintf(selectSegmentsSQL, whereSQL, isInUseStatusSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	segments := make([]*proto.Segment, 0, limit)
	featureIDsMap := map[string][]string{}
	for rows.Next() {
		segment := proto.Segment{}
		var status int32
		featureIDs := new(sql.NullString)
		err := rows.Scan(
			&segment.Id,
			&segment.Name,
			&segment.Description,
			&mysql.JSONObject{Val: &segment.Rules},
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
		return nil, 0, 0, nil, err
	}
	nextOffset := offset + len(segments)
	var totalCount int64
	countConditionSQL := "> 0 THEN 1 ELSE 1"
	if isInUseStatus != nil {
		if *isInUseStatus {
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
		return ErrSegmentUnexpectedAffectedRows
	}
	return nil
}

// ListAllInUseSegments lists all segments that are in use (referenced by feature flags).
func (s *segmentStorage) ListAllInUseSegments(
	ctx context.Context,
) ([]*InUseSegment, error) {
	rows, err := s.qe.QueryContext(ctx, selectAllInUseSegmentsSQLQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	segments := make([]*InUseSegment, 0)
	for rows.Next() {
		var seg InUseSegment
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

// ListSegmentUsersBySegment lists all users for a specific segment.
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
