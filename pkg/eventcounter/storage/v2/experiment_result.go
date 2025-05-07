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

	"github.com/bucketeer-io/bucketeer/pkg/eventcounter/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

var ErrExperimentResultNotFound = errors.New("experimentResult: experiment result not found")

var (
	//go:embed sql/select_experiment_result.sql
	selectExperimentResultSQL string
)

type ExperimentResultStorage interface {
	GetExperimentResult(ctx context.Context, id, environmentId string) (*domain.ExperimentResult, error)
}

type experimentResultStorage struct {
	qe mysql.QueryExecer
}

func NewExperimentResultStorage(qe mysql.QueryExecer) ExperimentResultStorage {
	return &experimentResultStorage{qe}
}

func (s *experimentResultStorage) GetExperimentResult(
	ctx context.Context,
	id, environmentId string,
) (*domain.ExperimentResult, error) {
	er := proto.ExperimentResult{}
	erForGoalResults := proto.ExperimentResult{}
	err := s.qe.QueryRowContext(
		ctx,
		selectExperimentResultSQL,
		id,
		environmentId,
	).Scan(
		&er.Id,
		&er.ExperimentId,
		&er.UpdatedAt,
		&mysql.JSONPBObject{Val: &erForGoalResults},
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrExperimentResultNotFound
		}
		return nil, err
	}
	er.GoalResults = erForGoalResults.GoalResults
	er.TotalEvaluationUserCount = erForGoalResults.TotalEvaluationUserCount
	return &domain.ExperimentResult{ExperimentResult: &er}, nil
}
