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

	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

var (
	//go:embed sql/update_experiment_result.sql
	updateExperimentResultSQL string
)

type ExperimentResultStorage interface {
	UpdateExperimentResult(ctx context.Context, environmentNamespace string,
		experimentResult *domain.ExperimentResult) error
}

type experimentResultStorage struct {
	qe mysql.QueryExecer
}

func NewExperimentResultStorage(qe mysql.QueryExecer) ExperimentResultStorage {
	return &experimentResultStorage{qe: qe}
}

func (e experimentResultStorage) UpdateExperimentResult(
	ctx context.Context,
	environmentNamespace string,
	experimentResult *domain.ExperimentResult,
) error {
	if _, err := e.qe.ExecContext(
		ctx,
		updateExperimentResultSQL,
		experimentResult.Id,
		experimentResult.ExperimentId,
		experimentResult.UpdatedAt,
		mysql.JSONObject{Val: experimentResult},
		environmentNamespace,
	); err != nil {
		return err
	}
	return nil
}
