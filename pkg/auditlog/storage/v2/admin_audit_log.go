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
	"strings"

	"github.com/bucketeer-io/bucketeer/pkg/auditlog/domain"
	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

var (
	ErrAdminAuditLogAlreadyExists = pkgErr.NewErrorAlreadyExists(
		pkgErr.AuditlogPackageName, "admin auditlog already exists")

	//go:embed sql/adminauditlog/insert_admin_audit_logs_v2.sql
	insertAdminAuditLogsV2SQL string
	//go:embed sql/adminauditlog/insert_admin_audit_log_v2.sql
	insertAdminAuditLogV2SQL string
	//go:embed sql/adminauditlog/select_admin_audit_log_v2.sql
	selectAdminAuditLogV2SQL string
	//go:embed sql/adminauditlog/select_admin_audit_log_v2_count.sql
	selectAdminAuditLogV2CountSQL string
)

type AdminAuditLogStorage interface {
	CreateAdminAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error
	CreateAdminAuditLog(ctx context.Context, auditLog *domain.AuditLog) error
	ListAdminAuditLogs(
		ctx context.Context,
		options *mysql.ListOptions,
	) ([]*proto.AuditLog, int, int64, error)
}

type adminAuditLogStorage struct {
	qe mysql.QueryExecer
}

func NewAdminAuditLogStorage(qe mysql.QueryExecer) AdminAuditLogStorage {
	return &adminAuditLogStorage{qe}
}

func (s *adminAuditLogStorage) CreateAdminAuditLogs(ctx context.Context, auditLogs []*domain.AuditLog) error {
	if len(auditLogs) == 0 {
		return nil
	}
	var query strings.Builder
	args := []interface{}{}
	for i, al := range auditLogs {
		if i != 0 {
			query.WriteString(",")
		} else {
			query.WriteString(insertAdminAuditLogsV2SQL)
		}
		query.WriteString(" (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
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
			al.EntityData,
			al.PreviousEntityData,
		)
	}
	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrAdminAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *adminAuditLogStorage) CreateAdminAuditLog(ctx context.Context, auditLog *domain.AuditLog) error {
	_, err := s.qe.ExecContext(
		ctx,
		insertAdminAuditLogV2SQL,
		auditLog.Id,
		auditLog.Timestamp,
		int32(auditLog.EntityType),
		auditLog.EntityId,
		int32(auditLog.Type),
		mysql.JSONObject{Val: auditLog.Event},
		mysql.JSONObject{Val: auditLog.Editor},
		mysql.JSONObject{Val: auditLog.Options},
		auditLog.EntityData,
		auditLog.PreviousEntityData,
	)
	if err != nil {
		if err == mysql.ErrDuplicateEntry {
			return ErrAdminAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *adminAuditLogStorage) ListAdminAuditLogs(
	ctx context.Context,
	options *mysql.ListOptions,
) ([]*proto.AuditLog, int, int64, error) {
	query, whereArgs := mysql.ConstructQueryAndWhereArgs(selectAdminAuditLogV2SQL, options)
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
	countQuery, countWhereArgs := mysql.ConstructCountQuery(selectAdminAuditLogV2CountSQL, options)
	err = s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, err
	}
	return auditLogs, nextOffset, totalCount, nil
}
