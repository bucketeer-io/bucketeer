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

package mysql

import (
	"context"
	_ "embed"
	"errors"
	"strconv"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	mysqlstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

var (
	//go:embed sql/adminauditlog/insert_admin_audit_logs_v2.sql
	insertAdminAuditLogsV2SQL string
	//go:embed sql/adminauditlog/insert_admin_audit_log_v2.sql
	insertAdminAuditLogV2SQL string
	//go:embed sql/adminauditlog/select_admin_audit_log_v2.sql
	selectAdminAuditLogV2SQL string
	//go:embed sql/adminauditlog/select_admin_audit_log_v2_count.sql
	selectAdminAuditLogV2CountSQL string
)

type adminAuditLogStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewAdminAuditLogStorage(qe mysqlstorage.QueryExecer) v2als.AdminAuditLogStorage {
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
			mysqlstorage.JSONObject{Val: al.Event},
			mysqlstorage.JSONObject{Val: al.Editor},
			mysqlstorage.JSONObject{Val: al.Options},
			al.EntityData,
			al.PreviousEntityData,
		)
	}
	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
			return v2als.ErrAdminAuditLogAlreadyExists
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
		mysqlstorage.JSONObject{Val: auditLog.Event},
		mysqlstorage.JSONObject{Val: auditLog.Editor},
		mysqlstorage.JSONObject{Val: auditLog.Options},
		auditLog.EntityData,
		auditLog.PreviousEntityData,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
			return v2als.ErrAdminAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *adminAuditLogStorage) ListAdminAuditLogs(
	ctx context.Context,
	params v2als.ListAdminAuditLogsParams,
) ([]*proto.AuditLog, int, int64, error) {
	options, err := listAdminAuditLogsOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectAdminAuditLogV2SQL, options)
	rows, err := s.qe.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()
	auditLogs := make([]*proto.AuditLog, 0, options.Limit)
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
			&mysqlstorage.JSONObject{Val: &auditLog.Event},
			&mysqlstorage.JSONObject{Val: &auditLog.Editor},
			&mysqlstorage.JSONObject{Val: &auditLog.Options},
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
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}
	nextOffset := options.Offset + len(auditLogs)
	var totalCount int64
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(selectAdminAuditLogV2CountSQL, options)
	if err := s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return auditLogs, nextOffset, totalCount, nil
}

func listAdminAuditLogsOptionsFromParams(p v2als.ListAdminAuditLogsParams) (*mysqlstorage.ListOptions, error) {
	var filters []*mysqlstorage.FilterV2
	if p.EntityType != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "entity_type",
			Operator: mysqlstorage.OperatorEqual,
			Value:    *p.EntityType,
		})
	}
	if p.From != 0 {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "timestamp",
			Operator: mysqlstorage.OperatorGreaterThanOrEqual,
			Value:    p.From,
		})
	}
	if p.To != 0 {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "timestamp",
			Operator: mysqlstorage.OperatorLessThanOrEqual,
			Value:    p.To,
		})
	}
	var searchQuery *mysqlstorage.SearchQuery
	if p.SearchKeyword != "" {
		searchQuery = &mysqlstorage.SearchQuery{
			Columns: []string{"editor"},
			Keyword: p.SearchKeyword,
		}
	}
	var column string
	switch p.OrderBy {
	case proto.ListAdminAuditLogsRequest_DEFAULT,
		proto.ListAdminAuditLogsRequest_TIMESTAMP:
		column = "timestamp"
	default:
		return nil, v2als.ErrInvalidOrderBy
	}
	direction := mysqlstorage.OrderDirectionDesc
	if p.OrderDirection == proto.ListAdminAuditLogsRequest_ASC {
		direction = mysqlstorage.OrderDirectionAsc
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, err
	}
	return &mysqlstorage.ListOptions{
		Limit:       p.PageSize,
		Offset:      offset,
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      []*mysqlstorage.Order{mysqlstorage.NewOrder(column, direction)},
	}, nil
}
