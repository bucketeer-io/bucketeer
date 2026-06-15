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
	"strconv"

	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

var (
	//go:embed sql/experiment/select_experiment.sql
	selectExperimentSQL string
	//go:embed sql/experiment/select_experiments.sql
	selectExperimentsSQL string
	//go:embed sql/experiment/count_experiment.sql
	countExperimentSQL string
	//go:embed sql/experiment/summarize_experiment.sql
	summarizeExperimentSQL string
	//go:embed sql/experiment/update_experiment.sql
	updateExperimentSQL string
	//go:embed sql/experiment/insert_experiment.sql
	insertExperimentSQL string
)

type experimentStorage struct {
	qe pgstorage.QueryExecer
}

func NewExperimentStorage(qe pgstorage.QueryExecer) v2es.ExperimentStorage {
	return &experimentStorage{qe: qe}
}

func (s *experimentStorage) CreateExperiment(
	ctx context.Context,
	e *domain.Experiment,
	environmentId string,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertExperimentSQL,
		e.Id,
		e.GoalId,
		e.FeatureId,
		e.FeatureVersion,
		pgstorage.JSONObject{Val: e.Variations},
		e.StartAt,
		e.StopAt,
		e.Stopped,
		e.StoppedAt,
		e.CreatedAt,
		e.UpdatedAt,
		e.Archived,
		e.Deleted,
		pgstorage.JSONObject{Val: e.GoalIds},
		e.Name,
		e.Description,
		e.BaseVariationId,
		int32(e.Status),
		e.Maintainer,
		environmentId,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
			return v2es.ErrExperimentAlreadyExists
		}
		return err
	}
	return nil
}

func (s *experimentStorage) UpdateExperiment(
	ctx context.Context,
	e *domain.Experiment,
	environmentId string,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		updateExperimentSQL,
		e.GoalId,
		e.FeatureId,
		e.FeatureVersion,
		pgstorage.JSONObject{Val: e.Variations},
		e.StartAt,
		e.StopAt,
		e.Stopped,
		e.StoppedAt,
		e.CreatedAt,
		e.UpdatedAt,
		e.Archived,
		e.Deleted,
		pgstorage.JSONObject{Val: e.GoalIds},
		e.Name,
		e.Description,
		e.BaseVariationId,
		e.Maintainer,
		int32(e.Status),
		e.Id,
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
		return v2es.ErrExperimentUnexpectedAffectedRows
	}
	return nil
}

func (s *experimentStorage) GetExperiment(
	ctx context.Context,
	id, environmentId string,
) (*domain.Experiment, error) {
	experiment := proto.Experiment{}
	var status int32
	err := s.qe.QueryRowContext(
		ctx,
		selectExperimentSQL,
		id,
		environmentId,
	).Scan(
		&experiment.Id,
		&experiment.GoalId,
		&experiment.FeatureId,
		&experiment.FeatureVersion,
		&pgstorage.JSONObject{Val: &experiment.Variations},
		&experiment.StartAt,
		&experiment.StopAt,
		&experiment.Stopped,
		&experiment.StoppedAt,
		&experiment.CreatedAt,
		&experiment.UpdatedAt,
		&experiment.Archived,
		&experiment.Deleted,
		&pgstorage.JSONObject{Val: &experiment.GoalIds},
		&experiment.Name,
		&experiment.Description,
		&experiment.BaseVariationId,
		&experiment.Maintainer,
		&status,
		&pgstorage.JSONObject{Val: &experiment.Goals},
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
			return nil, v2es.ErrExperimentNotFound
		}
		return nil, err
	}
	experiment.Status = proto.Experiment_Status(status)
	return &domain.Experiment{Experiment: &experiment}, nil
}

func (s *experimentStorage) ListExperiments(
	ctx context.Context,
	params v2es.ListExperimentsParams,
) ([]*proto.Experiment, int, int64, error) {
	options, err := listExperimentsOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectExperimentsSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	limit := options.Limit
	offset := options.Offset
	experiments := make([]*proto.Experiment, 0, limit)
	for rows.Next() {
		experiment := proto.Experiment{}
		var status int32
		err := rows.Scan(
			&experiment.Id,
			&experiment.GoalId,
			&experiment.FeatureId,
			&experiment.FeatureVersion,
			&pgstorage.JSONObject{Val: &experiment.Variations},
			&experiment.StartAt,
			&experiment.StopAt,
			&experiment.Stopped,
			&experiment.StoppedAt,
			&experiment.CreatedAt,
			&experiment.UpdatedAt,
			&experiment.Archived,
			&experiment.Deleted,
			&pgstorage.JSONObject{Val: &experiment.GoalIds},
			&experiment.Name,
			&experiment.Description,
			&experiment.BaseVariationId,
			&experiment.Maintainer,
			&status,
			&pgstorage.JSONObject{Val: &experiment.Goals},
		)
		if err != nil {
			return nil, 0, 0, err
		}
		experiment.Status = proto.Experiment_Status(status)
		experiments = append(experiments, &experiment)
	}
	if rows.Err() != nil {
		return nil, 0, 0, rows.Err()
	}
	nextOffset := offset + len(experiments)
	var totalCount int64
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(countExperimentSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}

	return experiments, nextOffset, totalCount, nil
}

func (s *experimentStorage) GetExperimentSummary(
	ctx context.Context,
	environmentID string,
) (*v2es.ExperimentSummary, error) {
	summary := &v2es.ExperimentSummary{}
	err := s.qe.QueryRowContext(ctx, summarizeExperimentSQL, environmentID).Scan(
		&summary.TotalWaitingCount,
		&summary.TotalRunningCount,
		&summary.TotalStoppedCount,
	)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func listExperimentsOptionsFromParams(
	p v2es.ListExperimentsParams,
) (*pgstorage.ListOptions, error) {
	filters := []*pgstorage.Filter{
		{
			Column:   "deleted",
			Operator: pgstorage.OperatorEqual,
			Value:    false,
		},
		{
			Column:   "environment_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		},
	}
	if p.Archived != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "archived",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.Archived,
		})
	}
	if p.FeatureID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.FeatureID,
		})
	}
	if p.FeatureVersion != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "feature_version",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.FeatureVersion,
		})
	}
	if p.StartAt != 0 {
		// When a start timestamp is provided,
		// use it as the lower bound for filtering.
		filters = append(filters, &pgstorage.Filter{
			Column:   "start_at",
			Operator: pgstorage.OperatorGreaterThanOrEqual,
			Value:    p.StartAt,
		})
	}
	if p.StopAt != 0 {
		// When a stop timestamp is provided:
		// - If p.StartAt is also provided, treat p.StopAt as an absolute upper bound.
		// (This selects experiments with stop_at <= p.StopAt.)
		// - If p.StartAt is not provided, treat p.StopAt as a relative cutoff timestamp.
		// (This selects experiments with stop_at >= p.StopAt.)
		if p.StartAt != 0 {
			filters = append(filters, &pgstorage.Filter{
				Column:   "stop_at",
				Operator: pgstorage.OperatorLessThanOrEqual,
				Value:    p.StopAt,
			})
		} else {
			filters = append(filters, &pgstorage.Filter{
				Column:   "stop_at",
				Operator: pgstorage.OperatorGreaterThanOrEqual,
				Value:    p.StopAt,
			})
		}
	}
	if p.Maintainer != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "maintainer",
			Operator: pgstorage.OperatorEqual,
			Value:    p.Maintainer,
		})
	}
	var inFilters []*pgstorage.InFilter
	if len(p.Statuses) > 0 {
		statuses := make([]interface{}, 0, len(p.Statuses))
		for _, sts := range p.Statuses {
			statuses = append(statuses, sts)
		}
		inFilters = append(inFilters, &pgstorage.InFilter{
			Column: "status",
			Values: statuses,
		})
	}
	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{"name", "description"},
			Keyword: p.SearchKeyword,
		}
	}
	orders, err := experimentListOrders(p.OrderBy, p.OrderDirection)
	if err != nil {
		return nil, err
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil || offset < 0 {
		return nil, v2es.ErrInvalidCursor
	}
	limit := p.PageSize
	if limit < 0 {
		limit = 0
	}
	return &pgstorage.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		Orders:      orders,
		InFilters:   inFilters,
		SearchQuery: searchQuery,
	}, nil
}

func experimentListOrders(
	orderBy proto.ListExperimentsRequest_OrderBy,
	orderDirection proto.ListExperimentsRequest_OrderDirection,
) ([]*pgstorage.Order, error) {
	var column string
	switch orderBy {
	case proto.ListExperimentsRequest_DEFAULT,
		proto.ListExperimentsRequest_NAME:
		column = "ex.name"
	case proto.ListExperimentsRequest_CREATED_AT:
		column = "ex.created_at"
	case proto.ListExperimentsRequest_UPDATED_AT:
		column = "ex.updated_at"
	case proto.ListExperimentsRequest_START_AT:
		column = "ex.start_at"
	case proto.ListExperimentsRequest_STOP_AT:
		column = "ex.stop_at"
	case proto.ListExperimentsRequest_STATUS:
		column = "ex.status"
	case proto.ListExperimentsRequest_GOALS_COUNT:
		column = "jsonb_array_length(ex.goal_ids::jsonb)"
	default:
		return nil, v2es.ErrInvalidOrderBy
	}
	direction := pgstorage.OrderDirectionAsc
	if orderDirection == proto.ListExperimentsRequest_DESC {
		direction = pgstorage.OrderDirectionDesc
	}
	return []*pgstorage.Order{pgstorage.NewOrder(column, direction)}, nil
}
