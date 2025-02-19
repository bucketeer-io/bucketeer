// Copyright 2025 The Bucketeer Authors.
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
	"fmt"
	"strings"

	"github.com/bucketeer-io/bucketeer/pkg/auditlog/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

var (
	ErrAuditLogAlreadyExists = errors.New("auditlog: auditlog already exists")
	//go:embed sql/auditlog/insert_audit_logs_v2.sql
	insertAuditLogsV2SQL string
	//go:embed sql/auditlog/insert_audit_log_v2.sql
	insertAuditLogV2SQL string
	//go:embed sql/auditlog/select_audit_log_v2.sql
	selectAuditLogV2SQL string
	//go:embed sql/auditlog/select_audit_log_v2_count.sql
	selectAuditLogV2CountSQL string
)

type AuditLogStorage interface {
	CreateAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error
	CreateAuditLog(ctx context.Context, auditLog *domain.AuditLog) error
	ListAuditLogs(
		ctx context.Context,
		whereParts []mysql.WherePart,
		orders []*mysql.Order,
		limit, offset int,
	) ([]*proto.AuditLog, int, int64, error)
}

type auditLogStorage struct {
	client mysql.Client
}

func NewAuditLogStorage(client mysql.Client) AuditLogStorage {
	return &auditLogStorage{client}
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
	_, err := s.client.Qe(ctx).ExecContext(ctx, query.String(), args...)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *auditLogStorage) CreateAuditLog(ctx context.Context, auditLog *domain.AuditLog) error {
	_, err := s.client.Qe(ctx).ExecContext(
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
		if err == mysql.ErrDuplicateEntry {
			return ErrAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *auditLogStorage) ListAuditLogs(
	ctx context.Context,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	limit, offset int,
) ([]*proto.AuditLog, int, int64, error) {
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	orderBySQL := mysql.ConstructOrderBySQLString(orders)
	limitOffsetSQL := mysql.ConstructLimitOffsetSQLString(limit, offset)
	query := fmt.Sprintf(selectAuditLogV2SQL, whereSQL, orderBySQL, limitOffsetSQL)
	rows, err := s.client.Qe(ctx).QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
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
	countQuery := fmt.Sprintf(selectAuditLogV2CountSQL, whereSQL)
	err = s.client.Qe(ctx).QueryRowContext(ctx, countQuery, whereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return auditLogs, nextOffset, totalCount, nil
}
