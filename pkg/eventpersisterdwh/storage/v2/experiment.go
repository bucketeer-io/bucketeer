// Copyright 2023 The Bucketeer Authors.
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
//

package v2

import (
	"context"
	_ "embed"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

var (
	//go:embed sql/count_experiment.sql
	countExperimentSql string
)

type ExperimentStorage interface {
	ListRunningExperiments(ctx context.Context) ([]*proto.Experiment, error)
}

type experimentStorage struct {
	qe mysql.QueryExecer
}

func NewExperimentStorage(qe mysql.QueryExecer) ExperimentStorage {
	return &experimentStorage{qe: qe}
}

func (e experimentStorage) ListRunningExperiments(ctx context.Context) ([]*proto.Experiment, error) {
	rows, err := e.qe.QueryContext(ctx, countExperimentSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	experiments := make([]*proto.Experiment, 0)
	for rows.Next() {
		var experiment proto.Experiment
		if err := rows.Scan(
			&experiment.Id,
			&experiment.Status,
			&experiment.StopAt,
		); err != nil {
			return nil, err
		}
		experiments = append(experiments, &experiment)
	}
	return experiments, nil
}
