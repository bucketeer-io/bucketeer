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

	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	operationalstorage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/operationalstorage"
)

var (
	//go:embed sql/count_auto_ops_rules.sql
	countAutoOpsRulesSql string
)

type autoOpsRuleStorage struct {
	qe pgstorage.QueryExecer
}

func NewAutoOpsRuleStorage(qe pgstorage.QueryExecer) operationalstorage.AutoOpsRuleStorage {
	return &autoOpsRuleStorage{qe: qe}
}

func (s autoOpsRuleStorage) CountOpsEventRate(ctx context.Context) (int, error) {
	var count int
	err := s.qe.QueryRowContext(ctx, countAutoOpsRulesSql).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
