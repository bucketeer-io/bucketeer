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

	"github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/domain"
	operationaldatabase "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/operational_database"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
)

var (
	//go:embed sql/select_experiment_result.sql
	selectExperimentResultSQL string
)

type experimentResultStorage struct {
	qe mysql.QueryExecer
}

func NewExperimentResultStorage(qe mysql.QueryExecer) operationaldatabase.ExperimentResultStorage {
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
			return nil, operationaldatabase.ErrExperimentResultNotFound
		}
		return nil, err
	}
	er.GoalResults = erForGoalResults.GoalResults
	er.TotalEvaluationUserCount = erForGoalResults.TotalEvaluationUserCount
	er.SrmResult = erForGoalResults.SrmResult
	return &domain.ExperimentResult{ExperimentResult: &er}, nil
}
