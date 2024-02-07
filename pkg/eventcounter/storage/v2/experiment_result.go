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

	"github.com/bucketeer-io/bucketeer/pkg/eventcounter/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

var ErrExperimentResultNotFound = errors.New("experimentResult: experiment result not found")

type ExperimentResultStorage interface {
	GetExperimentResult(ctx context.Context, id, environmentNamespace string) (*domain.ExperimentResult, error)
}

type experimentResultStorage struct {
	qe mysql.QueryExecer
}

func NewExperimentResultStorage(qe mysql.QueryExecer) ExperimentResultStorage {
	return &experimentResultStorage{qe}
}

func (s *experimentResultStorage) GetExperimentResult(
	ctx context.Context,
	id, environmentNamespace string,
) (*domain.ExperimentResult, error) {
	er := proto.ExperimentResult{}
	er_for_goal_results := proto.ExperimentResult{}
	query := `
		SELECT
			id,
			experiment_id,
			updated_at,
			data
		FROM
			experiment_result
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
		&er.Id,
		&er.ExperimentId,
		&er.UpdatedAt,
		&mysql.JSONPBObject{Val: &er_for_goal_results},
	)
	if err != nil {
		if err == mysql.ErrNoRows {
			return nil, ErrExperimentResultNotFound
		}
		return nil, err
	}
	er.GoalResults = er_for_goal_results.GoalResults
	return &domain.ExperimentResult{ExperimentResult: &er}, nil
}
