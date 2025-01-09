// Copyright 2025 The Bucketeer Authors.
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
	"errors"
	"fmt"
	"strings"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrSegmentAlreadyExists          = errors.New("segment: already exists")
	ErrSegmentNotFound               = errors.New("segment: not found")
	ErrSegmentUnexpectedAffectedRows = errors.New("segment: unexpected affected rows")
)

type SegmentStorage interface {
	CreateSegment(ctx context.Context, segment *domain.Segment, environmentId string) error
	UpdateSegment(ctx context.Context, segment *domain.Segment, environmentId string) error
	GetSegment(ctx context.Context, id, environmentId string) (*domain.Segment, []string, error)
	ListSegments(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
		isInUseStatus *bool,
		environmentId string,
	) ([]*proto.Segment, int, int64, map[string][]string, error)
	DeleteSegment(ctx context.Context, id string) error
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
	query := `
		INSERT INTO segment (
			id,
			name,
			description,
			rules,
			created_at,
			updated_at,
			version,
			deleted,
			included_user_count,
			excluded_user_count,
			status,
			environment_id
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?
		)
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
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
	query := `
		UPDATE 
			segment
		SET
			name = ?,
			description = ?,
			rules = ?,
			created_at = ?,
			updated_at = ?,
			version = ?,
			deleted = ?,
			included_user_count = ?,
			excluded_user_count = ?,
			status = ?
		WHERE
			id = ? AND
			environment_id = ?
	`
	result, err := s.qe.ExecContext(
		ctx,
		query,
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
	query := `
		SELECT
			id,
			name,
			description,
			rules,
			created_at,
			updated_at,
			version,
			deleted,
			included_user_count,
			excluded_user_count,
			status,
			(
				SELECT 
					GROUP_CONCAT(id)
				FROM 
					feature
				WHERE
					environment_id = ? AND
					rules LIKE concat("%%", segment.id, "%%")
			) AS feature_ids
		FROM
			segment
		WHERE
			id = ? AND
			environment_id = ?
	`
	featureIDs := new(sql.NullString)
	err := s.qe.QueryRowContext(
		ctx,
		query,
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
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
	isInUseStatus *bool,
	environmentId string,
) ([]*proto.Segment, int, int64, map[string][]string, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	prepareArgs := make([]interface{}, 0, len(whereArgs)+1)
	prepareArgs = append(prepareArgs, environmentId)
	prepareArgs = append(prepareArgs, whereArgs...)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	var isInUseStatusSQL string
	if isInUseStatus != nil {
		if *isInUseStatus {
			isInUseStatusSQL = "HAVING feature_ids IS NOT NULL"
		} else {
			isInUseStatusSQL = "HAVING feature_ids IS NULL"
		}
	}
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			rules,
			created_at,
			updated_at,
			version,
			deleted,
			included_user_count,
			excluded_user_count,
			status,
			(
				SELECT 
					GROUP_CONCAT(id)
				FROM 
					feature
				WHERE
					environment_id = ? AND
					rules LIKE concat("%%", segment.id, "%%")
			) AS feature_ids
		FROM
			segment
		%s %s %s %s
		`, whereSQL, isInUseStatusSQL, orderBySQL, limitOffsetSQL,
	)
	rows, err := s.qe.QueryContext(ctx, query, prepareArgs...)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	defer rows.Close()
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
	countQuery := fmt.Sprintf(`
		SELECT
			COUNT(	
				CASE 
					WHEN (
						SELECT 
							COUNT(1)
						FROM 
							feature
						WHERE
							environment_id = ? AND
							rules LIKE concat("%%", segment.id, "%%")
					) %s
				END
			)
		FROM
			segment
		%s
		`, countConditionSQL, whereSQL,
	)
	err = s.qe.QueryRowContext(ctx, countQuery, prepareArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	return segments, nextOffset, totalCount, featureIDsMap, nil
}

func (s *segmentStorage) DeleteSegment(ctx context.Context, id string) error {
	query := `
		DELETE FROM segment
			WHERE id = ?
	`
	result, err := s.qe.ExecContext(ctx, query, id)
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
