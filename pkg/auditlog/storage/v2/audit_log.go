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
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

var (
	ErrAuditLogAlreadyExists = pkgErr.NewErrorAlreadyExists(pkgErr.AuditlogPackageName, "auditlog already exists")
	ErrAuditLogNotFound      = pkgErr.NewErrorNotFound(pkgErr.AuditlogPackageName, "auditlog not found", "auditlog")
	//go:embed sql/auditlog/select_audit_log_v2.sql
	selectAuditLogV2SQL string
	//go:embed sql/auditlog/insert_audit_logs_v2.sql
	insertAuditLogsV2SQL string
	//go:embed sql/auditlog/insert_audit_log_v2.sql
	insertAuditLogV2SQL string
	//go:embed sql/auditlog/select_audit_logs_v2.sql
	selectAuditLogsV2SQL string
	//go:embed sql/auditlog/select_audit_log_v2_count.sql
	selectAuditLogV2CountSQL string
)

type AuditLogStorage interface {
	GetAuditLog(ctx context.Context, id string, environmentID string) (*proto.AuditLog, error)
	CreateAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error
	CreateAuditLog(ctx context.Context, auditLog *domain.AuditLog) error
	ListAuditLogs(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*proto.AuditLog, int, int64, error)
}

type auditLogStorage struct {
	qe mysql.QueryExecer
}

func NewAuditLogStorage(qe mysql.QueryExecer) AuditLogStorage {
	return &auditLogStorage{qe}
}

func (s *auditLogStorage) GetAuditLog(ctx context.Context, id string, environmentID string) (*proto.AuditLog, error) {
	auditLog := &proto.AuditLog{}
	var et int32
	var t int32
	row := s.qe.QueryRowContext(ctx, selectAuditLogV2SQL, environmentID, id)
	err := row.Scan(
		&auditLog.Id,
		&auditLog.Timestamp,
		&et,
		&auditLog.EntityId,
		&t,
		&mysql.JSONObject{Val: &auditLog.Event},
		&mysql.JSONObject{Val: &auditLog.Editor},
		&mysql.JSONObject{Val: &auditLog.Options},
		&auditLog.EntityData,
		&auditLog.PreviousEntityData,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrAuditLogNotFound
		}
		return nil, err
	}
	auditLog.EntityType = eventproto.Event_EntityType(et)
	auditLog.Type = eventproto.Event_Type(t)
	return auditLog, nil
}

func (s *auditLogStorage) CreateAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error {
	if len(auditLogs) == 0 {
		return nil
	}
	var query strings.Builder
	args := []interface{}{}
	for i, al := range auditLogs {
		if i != 0 {
			query.WriteString(",")
		} else {
			query.WriteString(insertAuditLogsV2SQL)
		}
		query.WriteString(" (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(
			args,
			al.Id,
			al.Timestamp,
			int32(al.EntityType),
			al.EntityId,
			int32(al.Type),
			mysql.JSONObject{Val: al.Event},
			mysql.JSONObject{Val: al.Editor},
			mysql.JSONObject{Val: al.Options},
			al.EnvironmentId,
			al.EntityData,
			al.PreviousEntityData,
		)
	}
	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *auditLogStorage) CreateAuditLog(ctx context.Context, auditLog *domain.AuditLog) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertAuditLogV2SQL,
		auditLog.Id,
		auditLog.Timestamp,
		int32(auditLog.EntityType),
		auditLog.EntityId,
		int32(auditLog.Type),
		mysql.JSONObject{Val: auditLog.Event},
		mysql.JSONObject{Val: auditLog.Editor},
		mysql.JSONObject{Val: auditLog.Options},
		auditLog.EnvironmentId,
		auditLog.EntityData,
		auditLog.PreviousEntityData,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *auditLogStorage) ListAuditLogs(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.AuditLog, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectAuditLogsV2SQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	var limit, offset int
	if options != nil {
		limit = options.Limit
		offset = options.Offset
	}
	auditLogs := make([]*proto.AuditLog, 0, limit)
	for rows.Next() {
		auditLog := proto.AuditLog{}
		var et int32
		var t int32
		err := rows.Scan(
			&auditLog.Id,
			&auditLog.Timestamp,
			&et,
			&auditLog.EntityId,
			&t,
			&mysql.JSONObject{Val: &auditLog.Event},
			&mysql.JSONObject{Val: &auditLog.Editor},
			&mysql.JSONObject{Val: &auditLog.Options},
			&auditLog.EntityData,
			&auditLog.PreviousEntityData,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		auditLog.EntityType = eventproto.Event_EntityType(et)
		auditLog.Type = eventproto.Event_Type(t)
		auditLogs = append(auditLogs, &auditLog)
	}
	if rows.Err() != nil {
		return nil, 0, 0, err
	}
	nextOffset := offset + len(auditLogs)
	var totalCount int64
	countQuery, countWhereArgs := mysql.ConstructCountQuery(selectAuditLogV2CountSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return auditLogs, nextOffset, totalCount, nil
}
