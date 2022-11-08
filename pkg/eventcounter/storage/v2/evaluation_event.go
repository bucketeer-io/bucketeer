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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v2

import (
	"context"
	"strconv"

	"github.com/lib/pq"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/postgres"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	esproto "github.com/bucketeer-io/bucketeer/proto/event/service"
)

type EventStorage interface {
	CreateEvaluationEvent(ctx context.Context, event *eventproto.EvaluationEvent, id, environmentNamespace string) error
	CreateGoalEvent(
		ctx context.Context,
		event *eventproto.GoalEvent,
		id, environmentNamespace string,
		evaluations []string,
	) error
	CreateUserEvent(
		ctx context.Context,
		event *esproto.UserEvent,
		id, environmentNamespace string,
	) error
}

type eventStorage struct {
	qe postgres.Execer
}

func NewEventStorage(qe postgres.Execer) EventStorage {
	return &eventStorage{qe: qe}
}

func (s *eventStorage) CreateEvaluationEvent(
	ctx context.Context,
	event *eventproto.EvaluationEvent,
	id, environmentNamespace string,
) error {
	query := `
		INSERT INTO evaluation_event (
			id,
			timestamp,
			feature_id,
			feature_version,
			variation_id,
			user_id,
			user_data,
			reason,
			tag,
			source_id,
			environment_namespace
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) ON CONFLICT DO NOTHING
	`
	userData := map[string]string{}
	if event.User != nil {
		userData = event.User.Data
	}
	reason := ""
	if event.Reason != nil {
		reason = event.Reason.Type.String()
	}
	_, err := s.qe.ExecContext(
		ctx,
		query,
		id,
		event.Timestamp,
		event.FeatureId,
		event.FeatureVersion,
		event.VariationId,
		event.UserId,
		postgres.JSONObject{Val: userData},
		reason,
		event.Tag,
		event.SourceId.String(),
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *eventStorage) CreateGoalEvent(
	ctx context.Context,
	event *eventproto.GoalEvent,
	id, environmentNamespace string,
	evaluations []string,
) error {
	query := `
		INSERT INTO goal_event (
			id,
			timestamp,
			goal_id,
			value,
			user_id,
			user_data,
			tag,
			source_id,
			environment_namespace,
			evaluations
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) ON CONFLICT DO NOTHING
	`
	userData := map[string]string{}
	if event.User != nil {
		userData = event.User.Data
	}
	_, err := s.qe.ExecContext(
		ctx,
		query,
		id,
		event.Timestamp,
		event.GoalId,
		strconv.FormatFloat(event.Value, 'f', -1, 64),
		event.UserId,
		postgres.JSONObject{Val: userData},
		event.Tag,
		event.SourceId.String(),
		environmentNamespace,
		pq.Array(evaluations),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *eventStorage) CreateUserEvent(
	ctx context.Context,
	event *esproto.UserEvent,
	id, environmentNamespace string,
) error {
	query := `
		INSERT INTO user_event (
			id,
			tag,
			user_id,
			timestamp,
			source_id,
			environment_namespace
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) ON CONFLICT DO NOTHING
	`
	_, err := s.qe.ExecContext(
		ctx,
		query,
		id,
		event.Tag,
		event.UserId,
		event.LastSeen,
		event.SourceId.String(),
		environmentNamespace,
	)
	if err != nil {
		return err
	}
	return nil
}
