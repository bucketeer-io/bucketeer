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
	_ "embed"
	"errors"
	"fmt"

	"github.com/bucketeer-io/bucketeer/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

var (
	ErrExperimentAlreadyExists          = errors.New("experiment: already exists")
	ErrExperimentNotFound               = errors.New("experiment: not found")
	ErrExperimentUnexpectedAffectedRows = errors.New("experiment: unexpected affected rows")

	//go:embed sql/experiment/select_experiment.sql
	selectExperimentSQL string
	//go:embed sql/experiment/select_experiments.sql
	selectExperimentsSQL string
	//go:embed sql/experiment/count_experiment.sql
	countExperimentSQL string
	//go:embed sql/experiment/update_experiment.sql
	updateExperimentSQL string
	//go:embed sql/experiment/insert_experiment.sql
	insertExperimentSQL string
)

type ExperimentStorage interface {
	CreateExperiment(ctx context.Context, e *domain.Experiment, environmentId string) error
	UpdateExperiment(ctx context.Context, e *domain.Experiment, environmentId string) error
	GetExperiment(ctx context.Context, id, environmentId string) (*domain.Experiment, error)
	ListExperiments(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.Experiment, int, int64, *proto.ListExperimentsResponse_Summary, error)
}

type experimentStorage struct {
	qe mysql.QueryExecer
}

func NewExperimentStorage(qe mysql.QueryExecer) ExperimentStorage {
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
		mysql.JSONObject{Val: e.Variations},
		e.StartAt,
		e.StopAt,
		e.Stopped,
		e.StoppedAt,
		e.CreatedAt,
		e.UpdatedAt,
		e.Archived,
		e.Deleted,
		mysql.JSONObject{Val: e.GoalIds},
		e.Name,
		e.Description,
		e.BaseVariationId,
		int32(e.Status),
		e.Maintainer,
		environmentId,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrExperimentAlreadyExists
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
		mysql.JSONObject{Val: e.Variations},
		e.StartAt,
		e.StopAt,
		e.Stopped,
		e.StoppedAt,
		e.CreatedAt,
		e.UpdatedAt,
		e.Archived,
		e.Deleted,
		mysql.JSONObject{Val: e.GoalIds},
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
		return ErrExperimentUnexpectedAffectedRows
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
		&mysql.JSONObject{Val: &experiment.Variations},
		&experiment.StartAt,
		&experiment.StopAt,
		&experiment.Stopped,
		&experiment.StoppedAt,
		&experiment.CreatedAt,
		&experiment.UpdatedAt,
		&experiment.Archived,
		&experiment.Deleted,
		&mysql.JSONObject{Val: &experiment.GoalIds},
		&experiment.Name,
		&experiment.Description,
		&experiment.BaseVariationId,
		&experiment.Maintainer,
		&status,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrExperimentNotFound
		}
		return nil, err
	}
	experiment.Status = proto.Experiment_Status(status)
	return &domain.Experiment{Experiment: &experiment}, nil
}

func (s *experimentStorage) ListExperiments(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.Experiment, int, int64, *proto.ListExperimentsResponse_Summary, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectExperimentsSQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	defer rows.Close()
	experiments := make([]*proto.Experiment, 0, limit)
	for rows.Next() {
		experiment := proto.Experiment{}
		var status int32
		err := rows.Scan(
			&experiment.Id,
			&experiment.GoalId,
			&experiment.FeatureId,
			&experiment.FeatureVersion,
			&mysql.JSONObject{Val: &experiment.Variations},
			&experiment.StartAt,
			&experiment.StopAt,
			&experiment.Stopped,
			&experiment.StoppedAt,
			&experiment.CreatedAt,
			&experiment.UpdatedAt,
			&experiment.Archived,
			&experiment.Deleted,
			&mysql.JSONObject{Val: &experiment.GoalIds},
			&experiment.Name,
			&experiment.Description,
			&experiment.BaseVariationId,
			&experiment.Maintainer,
			&status,
			&mysql.JSONObject{Val: &experiment.Goals},
		)
		if err != nil {
			return nil, 0, 0, nil, err
		}
		experiment.Status = proto.Experiment_Status(status)
		experiments = append(experiments, &experiment)
	}
	if rows.Err() != nil {
		return nil, 0, 0, nil, err
	}
	nextOffset := offset + len(experiments)
	var totalCount int64
	summary := &proto.ListExperimentsResponse_Summary{}
	countQuery := fmt.Sprintf(countExperimentSQL, whereSQL)
	err = s.qe.QueryRowContext(ctx, countQuery, whereArgs...).Scan(
		&totalCount,
		&summary.TotalWaitingCount,
		&summary.TotalRunningCount,
		&summary.TotalStoppedCount,
	)
	if err != nil {
		return nil, 0, 0, nil, err
	}
	return experiments, nextOffset, totalCount, summary, nil
}
