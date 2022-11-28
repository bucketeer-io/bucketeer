// Copyright 2022 The Bucketeer Authors.
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

package v2

import (
	"context"
	"fmt"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	esproto "github.com/bucketeer-io/bucketeer/proto/event/service"
)

type MAUStorage interface {
	UpsertMAU(ctx context.Context, event *esproto.UserEvent, environmentNamespace string) error
}

// mysqlMAUStorage is the temporal implementation.
// We plan to replace the mysql with the postgresql.
type mysqlMAUStorage struct {
	qe mysql.QueryExecer
}

func NewMysqlMAUStorage(qe mysql.QueryExecer) MAUStorage {
	return &mysqlMAUStorage{qe: qe}
}

func (s *mysqlMAUStorage) UpsertMAU(ctx context.Context, event *esproto.UserEvent, environmentNamespace string) error {
	query := `
	INSERT INTO mau (
		user_id,
		yearmonth,
		source_id,
		event_count,
		created_at,
		updated_at,
		environment_namespace
	) VALUES (
		?, ?, ?, ?, ?, ?, ?
	) ON DUPLICATE KEY UPDATE
		event_count = event_count + 1,
		updated_at = VALUES(updated_at)
	`
	now := time.Now()
	yearMonth := fmt.Sprintf("%d%d", now.Year(), now.Month())
	_, err := s.qe.ExecContext(
		ctx,
		query,
		event.UserId,
		yearMonth,
		event.SourceId,
		1,
		now.Unix(),
		now.Unix(),
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	return nil
}
