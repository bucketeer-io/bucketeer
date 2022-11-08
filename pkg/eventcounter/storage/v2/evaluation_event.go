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

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/postgres"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

type EvaluationEventStorage interface {
	CreateEvaluationEvent(ctx context.Context, event *eventproto.EvaluationEvent, id, environmentNamespace string) error
}

type eventStorage struct {
	qe postgres.Execer
}

func NewEvaluationEventStorage(qe postgres.Execer) EvaluationEventStorage {
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
