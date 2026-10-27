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

package mysql

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

var (
	//go:embed sql/holdout/insert_holdout_group.sql
	insertHoldoutGroupSQL string
	//go:embed sql/holdout/update_holdout_group.sql
	updateHoldoutGroupSQL string
	//go:embed sql/holdout/select_holdout_group.sql
	selectHoldoutGroupSQL string
	//go:embed sql/holdout/select_holdout_groups.sql
	selectHoldoutGroupsSQL string
	//go:embed sql/holdout/count_holdout_group.sql
	countHoldoutGroupSQL string
	//go:embed sql/holdout/insert_holdout_group_feature.sql
	insertHoldoutGroupFeatureSQL string
	//go:embed sql/holdout/remove_holdout_group_feature.sql
	removeHoldoutGroupFeatureSQL string
	//go:embed sql/holdout/lock_environment.sql
	lockEnvironmentSQL string
	//go:embed sql/holdout/select_sampling_seed.sql
	selectSamplingSeedSQL string
	//go:embed sql/holdout/select_active_groups.sql
	selectActiveGroupsSQL string
	//go:embed sql/holdout/select_live_serving_group.sql
	selectLiveServingGroupSQL string
	//go:embed sql/holdout/select_live_feature_targets.sql
	selectLiveFeatureTargetsSQL string
)

type holdoutStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewHoldoutStorage(qe mysqlstorage.QueryExecer) v2hs.HoldoutStorage {
	return &holdoutStorage{qe: qe}
}

func (s *holdoutStorage) CreateHoldoutGroup(ctx context.Context, group *domain.HoldoutGroup) error {
	_, err := s.qe.ExecContext(ctx, insertHoldoutGroupSQL, insertArgs(group)...)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
			return v2hs.ErrHoldoutGroupAlreadyExists
		}
		return err
	}
	return nil
}

func (s *holdoutStorage) UpdateHoldoutGroup(ctx context.Context, group *domain.HoldoutGroup) error {
	result, err := s.qe.ExecContext(ctx, updateHoldoutGroupSQL, updateArgs(group)...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	return s.confirmHoldoutGroupUpdate(ctx, group, rows)
}

// confirmHoldoutGroupUpdate treats a zero-row update of an existing group as
// success. MySQL reports zero rows when every column is unchanged, which
// happens when the same admin adds or removes a target in the same second.
func (s *holdoutStorage) confirmHoldoutGroupUpdate(
	ctx context.Context,
	group *domain.HoldoutGroup,
	rows int64,
) error {
	if rows == 1 {
		return nil
	}
	if rows != 0 {
		return v2hs.ErrHoldoutGroupUnexpectedAffectedRows
	}
	_, err := s.GetHoldoutGroup(ctx, group.Id, group.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2hs.ErrHoldoutGroupNotFound) {
			return v2hs.ErrHoldoutGroupUnexpectedAffectedRows
		}
		return err
	}
	return nil
}

func (s *holdoutStorage) GetHoldoutGroup(
	ctx context.Context,
	id, environmentID string,
) (*domain.HoldoutGroup, error) {
	return s.get(ctx, selectHoldoutGroupSQL, id, environmentID)
}

func (s *holdoutStorage) GetHoldoutGroupForUpdate(
	ctx context.Context,
	id, environmentID string,
) (*domain.HoldoutGroup, error) {
	return s.get(ctx, selectHoldoutGroupSQL+" FOR UPDATE", id, environmentID)
}

func (s *holdoutStorage) GetLiveServingGroup(
	ctx context.Context,
	environmentID string,
) (*domain.HoldoutGroup, error) {
	group, err := scanHoldoutGroup(s.qe.QueryRowContext(
		ctx,
		selectLiveServingGroupSQL,
		environmentID,
		false,
		false,
		int32(proto.HoldoutGroup_RUNNING),
		int32(proto.HoldoutGroup_ANALYSIS),
	))
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, v2hs.ErrHoldoutGroupNotFound
		}
		return nil, err
	}
	features, err := s.ListCurrentFeatures(ctx, environmentID, []string{group.Id})
	if err != nil {
		return nil, err
	}
	for _, feature := range features {
		group.Features = append(group.Features, feature.HoldoutGroupFeature)
	}
	return group, nil
}

func (s *holdoutStorage) get(
	ctx context.Context,
	query, id, environmentID string,
) (*domain.HoldoutGroup, error) {
	group, err := scanHoldoutGroup(s.qe.QueryRowContext(ctx, query, id, environmentID, false))
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, v2hs.ErrHoldoutGroupNotFound
		}
		return nil, err
	}
	return group, nil
}

func (s *holdoutStorage) ListHoldoutGroups(
	ctx context.Context,
	params v2hs.ListHoldoutGroupsParams,
) ([]*domain.HoldoutGroup, int, int64, error) {
	options, err := listOptions(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, args := mysqlstorage.ConstructQueryAndWhereArgs(selectHoldoutGroupsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	groups := make([]*domain.HoldoutGroup, 0, options.Limit)
	for rows.Next() {
		group, err := scanHoldoutGroup(rows)
		if err != nil {
			return nil, 0, 0, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}
	var total int64
	countQuery, countArgs := mysqlstorage.ConstructCountQuery(countHoldoutGroupSQL, options)
	if err := s.qe.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, 0, err
	}
	return groups, options.Offset + len(groups), total, nil
}

func (s *holdoutStorage) CreateHoldoutGroupFeatures(
	ctx context.Context,
	features []*domain.HoldoutGroupFeature,
) error {
	for _, feature := range features {
		if _, err := s.qe.ExecContext(
			ctx,
			insertHoldoutGroupFeatureSQL,
			feature.Id,
			feature.HoldoutGroupId,
			feature.FeatureId,
			feature.EnvironmentId,
			feature.HoldoutVariationId,
			feature.AddedAt,
			feature.RemovedAt,
		); err != nil {
			if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
				return v2hs.ErrHoldoutGroupAlreadyExists
			}
			return err
		}
	}
	return nil
}

func (s *holdoutStorage) RemoveHoldoutGroupFeature(
	ctx context.Context,
	environmentID, groupID, featureID string,
	removedAt int64,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		removeHoldoutGroupFeatureSQL,
		removedAt,
		groupID,
		environmentID,
		featureID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return v2hs.ErrHoldoutFeatureNotFound
	}
	return nil
}

func (s *holdoutStorage) ListCurrentFeatures(
	ctx context.Context,
	environmentID string,
	groupIDs []string,
) ([]*domain.HoldoutGroupFeature, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(groupIDs))
	args := make([]interface{}, 0, len(groupIDs)+1)
	args = append(args, environmentID)
	for i, id := range groupIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}
	query := fmt.Sprintf(`SELECT
    id, holdout_group_id, feature_id, environment_id, holdout_variation_id, added_at, removed_at
FROM holdout_group_feature
WHERE environment_id = ? AND removed_at = 0 AND holdout_group_id IN (%s)`, strings.Join(placeholders, ", "))
	rows, err := s.qe.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFeatures(rows)
}

func (s *holdoutStorage) LockEnvironment(ctx context.Context, environmentID string) error {
	var id string
	err := s.qe.QueryRowContext(ctx, lockEnvironmentSQL, environmentID).Scan(&id)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return v2hs.ErrEnvironmentNotFound
		}
		return err
	}
	return nil
}

func (s *holdoutStorage) GetEnvironmentSamplingSeed(ctx context.Context, environmentID string) (string, error) {
	var seed string
	err := s.qe.QueryRowContext(ctx, selectSamplingSeedSQL, environmentID).Scan(&seed)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return seed, nil
}

func (s *holdoutStorage) ListActiveGroups(
	ctx context.Context,
	environmentID string,
) ([]domain.GroupSnapshot, error) {
	rows, err := s.qe.QueryContext(
		ctx,
		selectActiveGroupsSQL,
		environmentID,
		false,
		int32(proto.HoldoutGroup_WAITING),
		int32(proto.HoldoutGroup_RUNNING),
		int32(proto.HoldoutGroup_ANALYSIS),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []domain.GroupSnapshot
	for rows.Next() {
		var snapshot domain.GroupSnapshot
		var status int32
		if err := rows.Scan(
			&snapshot.ID,
			&snapshot.Name,
			&status,
			&snapshot.StartAt,
			&snapshot.StopAt,
			&snapshot.SegmentID,
		); err != nil {
			return nil, err
		}
		snapshot.Status = proto.HoldoutGroup_Status(status)
		groups = append(groups, snapshot)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *holdoutStorage) ListLiveFeatureTargets(
	ctx context.Context,
	environmentID, featureID string,
) ([]domain.FeatureTarget, error) {
	rows, err := s.qe.QueryContext(
		ctx,
		selectLiveFeatureTargetsSQL,
		environmentID,
		featureID,
		false,
		int32(proto.HoldoutGroup_WAITING),
		int32(proto.HoldoutGroup_RUNNING),
		int32(proto.HoldoutGroup_ANALYSIS),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var targets []domain.FeatureTarget
	for rows.Next() {
		var target domain.FeatureTarget
		var status int32
		if err := rows.Scan(
			&target.GroupID,
			&target.GroupName,
			&status,
			&target.FeatureID,
			&target.VariationID,
		); err != nil {
			return nil, err
		}
		target.Status = proto.HoldoutGroup_Status(status)
		targets = append(targets, target)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return targets, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanHoldoutGroup(sc rowScanner) (*domain.HoldoutGroup, error) {
	group := &proto.HoldoutGroup{}
	var (
		assignmentType int32
		status         int32
		bucketStart    v2hs.Decimal
		percentage     v2hs.Decimal
	)
	err := sc.Scan(
		&group.Id,
		&group.EnvironmentId,
		&group.Name,
		&group.Description,
		&assignmentType,
		&group.SamplingSeed,
		&bucketStart,
		&percentage,
		&group.SegmentId,
		&status,
		&group.StartAt,
		&group.StopAt,
		&group.StartAnalysisAt,
		&group.AnalysisStartedAt,
		&group.StoppedAt,
		&mysqlstorage.JSONObject{Val: &group.GoalIds},
		&group.ExperimentId,
		&group.AnalysisExperimentId,
		&group.CreatedBy,
		&group.UpdatedBy,
		&group.CreatedAt,
		&group.UpdatedAt,
		&group.Archived,
		&group.Deleted,
	)
	if err != nil {
		return nil, err
	}
	group.AssignmentType = proto.HoldoutGroup_AssignmentType(assignmentType)
	group.Status = proto.HoldoutGroup_Status(status)
	group.BucketStart = float64(bucketStart)
	group.Percentage = float64(percentage)
	group.GoalIds = v2hs.NormalizeGoalIDs(group.GoalIds)
	return &domain.HoldoutGroup{HoldoutGroup: group}, nil
}

func scanFeatures(rows mysqlstorage.Rows) ([]*domain.HoldoutGroupFeature, error) {
	var features []*domain.HoldoutGroupFeature
	for rows.Next() {
		feature := &proto.HoldoutGroupFeature{}
		if err := rows.Scan(
			&feature.Id,
			&feature.HoldoutGroupId,
			&feature.FeatureId,
			&feature.EnvironmentId,
			&feature.HoldoutVariationId,
			&feature.AddedAt,
			&feature.RemovedAt,
		); err != nil {
			return nil, err
		}
		features = append(features, &domain.HoldoutGroupFeature{HoldoutGroupFeature: feature})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return features, nil
}

func insertArgs(group *domain.HoldoutGroup) []interface{} {
	return []interface{}{
		group.Id,
		group.EnvironmentId,
		group.Name,
		group.Description,
		int32(group.AssignmentType),
		group.SamplingSeed,
		v2hs.Decimal(group.BucketStart),
		v2hs.Decimal(group.Percentage),
		group.SegmentId,
		int32(group.Status),
		group.StartAt,
		group.StopAt,
		group.StartAnalysisAt,
		group.AnalysisStartedAt,
		group.StoppedAt,
		mysqlstorage.JSONObject{Val: v2hs.NormalizeGoalIDs(group.GoalIds)},
		group.ExperimentId,
		group.AnalysisExperimentId,
		group.CreatedBy,
		group.UpdatedBy,
		group.CreatedAt,
		group.UpdatedAt,
		group.Archived,
		group.Deleted,
	}
}

func updateArgs(group *domain.HoldoutGroup) []interface{} {
	return []interface{}{
		group.Name,
		group.Description,
		int32(group.Status),
		group.StartAt,
		group.StopAt,
		group.StartAnalysisAt,
		group.AnalysisStartedAt,
		group.StoppedAt,
		mysqlstorage.JSONObject{Val: v2hs.NormalizeGoalIDs(group.GoalIds)},
		group.ExperimentId,
		group.AnalysisExperimentId,
		group.UpdatedBy,
		group.UpdatedAt,
		group.Archived,
		group.Deleted,
		group.Id,
		group.EnvironmentId,
		false,
	}
}

func listOptions(p v2hs.ListHoldoutGroupsParams) (*mysqlstorage.ListOptions, error) {
	filters := []*mysqlstorage.FilterV2{
		{Column: "deleted", Operator: mysqlstorage.OperatorEqual, Value: false},
		{Column: "environment_id", Operator: mysqlstorage.OperatorEqual, Value: p.EnvironmentID},
	}
	if p.Archived != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column: "archived", Operator: mysqlstorage.OperatorEqual, Value: *p.Archived,
		})
	}
	var inFilters []*mysqlstorage.InFilter
	if len(p.Statuses) > 0 {
		values := make([]interface{}, 0, len(p.Statuses))
		for _, status := range p.Statuses {
			values = append(values, int32(status))
		}
		inFilters = append(inFilters, &mysqlstorage.InFilter{Column: "status", Values: values})
	}
	var search *mysqlstorage.SearchQuery
	if p.SearchKeyword != "" {
		search = &mysqlstorage.SearchQuery{Columns: []string{"name", "description"}, Keyword: p.SearchKeyword}
	}
	orders, err := listOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, err
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, v2hs.ErrInvalidCursor
	}
	limit := p.PageSize
	if limit < 0 {
		limit = 0
	}
	return &mysqlstorage.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		InFilters:   inFilters,
		SearchQuery: search,
		Orders:      orders,
	}, nil
}

func listOrders(
	orderBy proto.ListHoldoutGroupsRequest_OrderBy,
	orderDirection proto.ListHoldoutGroupsRequest_OrderDirection,
) ([]*mysqlstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListHoldoutGroupsRequest_DEFAULT, proto.ListHoldoutGroupsRequest_NAME:
		column = "name"
	case proto.ListHoldoutGroupsRequest_CREATED_AT:
		column = "created_at"
	case proto.ListHoldoutGroupsRequest_UPDATED_AT:
		column = "updated_at"
	case proto.ListHoldoutGroupsRequest_STATUS:
		column = "status"
	default:
		return nil, v2hs.ErrInvalidOrderBy
	}
	direction := mysqlstorage.OrderDirectionAsc
	if orderDirection == proto.ListHoldoutGroupsRequest_DESC {
		direction = mysqlstorage.OrderDirectionDesc
	}
	return []*mysqlstorage.Order{
		mysqlstorage.NewOrder(column, direction),
		mysqlstorage.NewOrder("id", direction),
	}, nil
}
