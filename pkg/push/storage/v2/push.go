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
	_ "embed"
	"errors"

	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"

	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

var (
	ErrPushAlreadyExists          = err.NewErrorAlreadyExists(err.PushPackageName, "push already exists")
	ErrPushNotFound               = err.NewErrorNotFound(err.PushPackageName, "push not found", "push")
	ErrPushUnexpectedAffectedRows = err.NewErrorUnexpectedAffectedRows(
		err.PushPackageName,
		"push unexpected affected rows",
	)

	//go:embed sql/push/insert_push.sql
	insertPushSQL string
	//go:embed sql/push/update_push.sql
	updatePushSQL string
	//go:embed sql/push/select_push.sql
	selectPushSQL string
	//go:embed sql/push/list_pushes.sql
	listPushesSQL string
	//go:embed sql/push/count_pushes.sql
	countPushesSQL string
	//go:embed sql/push/delete_push.sql
	deletePushSQL string
)

type PushStorage interface {
	CreatePush(ctx context.Context, e *domain.Push, environmentId string) error
	UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error
	GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error)
	ListPushes(
		ctx context.Context,
		option *mysql.ListOptions,
	) ([]*proto.Push, int, int64, error)
	DeletePush(ctx context.Context, id, environmentId string) error
}

type pushStorage struct {
	qe mysql.QueryExecer
}

func NewPushStorage(qe mysql.QueryExecer) PushStorage {
	return &pushStorage{qe: qe}
}

func (s *pushStorage) CreatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertPushSQL,
		e.Id,
		e.FcmServiceAccount,
		mysql.JSONObject{Val: e.Tags},
		e.Deleted,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
		environmentId,
		e.Disabled,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrPushAlreadyExists
		}
		return err
	}
	return nil
}

func (s *pushStorage) UpdatePush(ctx context.Context, e *domain.Push, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		updatePushSQL,
		e.FcmServiceAccount,
		mysql.JSONObject{Val: e.Tags},
		e.Deleted,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
		e.Disabled,
		e.Id,
		environmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrPushUnexpectedAffectedRows
	}
	return nil
}

func (s *pushStorage) GetPush(ctx context.Context, id, environmentId string) (*domain.Push, error) {
	push := proto.Push{}
	err := s.qe.QueryRowContext(
		ctx,
		selectPushSQL,
		id,
		environmentId,
	).Scan(
		&push.Id,
		&push.FcmServiceAccount,
		&mysql.JSONObject{Val: &push.Tags},
		&push.Deleted,
		&push.Name,
		&push.CreatedAt,
		&push.UpdatedAt,
		&push.Disabled,
		&push.EnvironmentId,
		&push.EnvironmentName,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrPushNotFound
		}
		return nil, err
	}
	return &domain.Push{Push: &push}, nil
}

func (s *pushStorage) ListPushes(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.Push, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(listPushesSQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	var offset int
	var limit int
	if options != nil {
		offset = options.Offset
		limit = options.Limit
	}

	defer rows.Close()
	pushes := make([]*proto.Push, 0, limit)
	for rows.Next() {
		push := proto.Push{}
		err := rows.Scan(
			&push.Id,
			&push.FcmServiceAccount,
			&mysql.JSONObject{Val: &push.Tags},
			&push.Deleted,
			&push.Name,
			&push.CreatedAt,
			&push.UpdatedAt,
			&push.Disabled,
			&push.EnvironmentId,
			&push.EnvironmentName,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		pushes = append(pushes, &push)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(pushes)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(countPushesSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return pushes, nextOffset, totalCount, nil
}

func (s *pushStorage) DeletePush(ctx context.Context, id, environmentId string) error {
	result, err := s.qe.ExecContext(
		ctx,
		deletePushSQL,
		id,
		environmentId,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrPushUnexpectedAffectedRows
	}
	return nil
}
