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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	esproto "github.com/bucketeer-io/bucketeer/v2/proto/event/service"
)

const (
	batchSize = 1000
)

type MAUStorage interface {
	UpsertMAU(ctx context.Context, event *esproto.UserEvent, environmentId string) error
	UpsertMAUs(ctx context.Context, events []*esproto.UserEvent, environmentId string) error
}

// mysqlMAUStorage is the temporal implementation.
// We plan to replace the mysql with the postgresql.
type mysqlMAUStorage struct {
	qe mysql.QueryExecer
}

func NewMysqlMAUStorage(qe mysql.QueryExecer) MAUStorage {
	return &mysqlMAUStorage{qe: qe}
}

func (s *mysqlMAUStorage) UpsertMAU(ctx context.Context, event *esproto.UserEvent, environmentId string) error {
	query := `
	INSERT INTO mau (
		user_id,
		yearmonth,
		source_id,
		event_count,
		created_at,
		updated_at,
		environment_id
	) VALUES (
		?, ?, ?, ?, ?, ?, ?
	) ON DUPLICATE KEY UPDATE
		event_count = event_count + 1,
		updated_at = VALUES(updated_at)
	`
	t := time.Unix(event.LastSeen, 0)
	yearMonth := fmt.Sprintf("%d%02d", t.Year(), t.Month())
	_, err := s.qe.ExecContext(
		ctx,
		query,
		event.UserId,
		yearMonth,
		event.SourceId,
		1,
		event.LastSeen,
		event.LastSeen,
		environmentId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *mysqlMAUStorage) UpsertMAUs(
	ctx context.Context,
	events []*esproto.UserEvent,
	environmentId string,
) error {
	// Upsert the events in batches
	for i := 0; i < len(events); i += batchSize {
		j := i + batchSize
		if j > len(events) {
			j = len(events)
		}
		if err := s.upsertMAUs(
			ctx,
			events[i:j],
			environmentId,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *mysqlMAUStorage) upsertMAUs(
	ctx context.Context,
	users []*esproto.UserEvent,
	environmentId string,
) error {
	var query strings.Builder
	query.WriteString(`
		INSERT INTO mau (
			user_id,
			yearmonth,
			source_id,
			event_count,
			created_at,
			updated_at,
			environment_id
		) VALUES
	`)
	args := []interface{}{}
	for i, event := range users {
		if i != 0 {
			query.WriteString(",")
		}
		t := time.Unix(event.LastSeen, 0)
		yearMonth := fmt.Sprintf("%d%02d", t.Year(), t.Month())
		query.WriteString(" (?, ?, ?, ?, ?, ?, ?)")
		args = append(
			args,
			event.UserId,
			yearMonth,
			event.SourceId,
			1,
			event.LastSeen,
			event.LastSeen,
			environmentId,
		)
	}
	query.WriteString(`
		ON DUPLICATE KEY UPDATE
		event_count = event_count + 1,
		updated_at = VALUES(updated_at)
	`)
	_, err := s.qe.ExecContext(
		ctx,
		query.String(),
		args...,
	)
	if err != nil {
		return err
	}
	return nil
}
