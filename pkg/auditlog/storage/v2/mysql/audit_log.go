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

type auditLogStorage struct {
	qe mysqlstorage.QueryExecer
}

func NewAuditLogStorage(qe mysqlstorage.QueryExecer) v2als.AuditLogStorage {
	return &auditLogStorage{qe}
}

func (s *auditLogStorage) GetAuditLog(
	ctx context.Context,
	id string,
	environmentID string,
) (*proto.AuditLog, error) {
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
		&mysqlstorage.JSONObject{Val: &auditLog.Event},
		&mysqlstorage.JSONObject{Val: &auditLog.Editor},
		&mysqlstorage.JSONObject{Val: &auditLog.Options},
		&auditLog.EntityData,
		&auditLog.PreviousEntityData,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrNoRows) {
			return nil, v2als.ErrAuditLogNotFound
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
			mysqlstorage.JSONObject{Val: al.Event},
			mysqlstorage.JSONObject{Val: al.Editor},
			mysqlstorage.JSONObject{Val: al.Options},
			al.EnvironmentId,
			al.EntityData,
			al.PreviousEntityData,
		)
	}
	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
			return v2als.ErrAuditLogAlreadyExists
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
		mysqlstorage.JSONObject{Val: auditLog.Event},
		mysqlstorage.JSONObject{Val: auditLog.Editor},
		mysqlstorage.JSONObject{Val: auditLog.Options},
		auditLog.EnvironmentId,
		auditLog.EntityData,
		auditLog.PreviousEntityData,
	)
	if err != nil {
		if errors.Is(err, mysqlstorage.ErrDuplicateEntry) {
			return v2als.ErrAuditLogAlreadyExists
		}
		return err
	}
	return nil
}

func (s *auditLogStorage) ListAuditLogs(
	ctx context.Context,
	params v2als.ListAuditLogsParams,
) ([]*proto.AuditLog, int, int64, error) {
	options, err := listAuditLogsOptionsFromParams(params)
	if err != nil {
		return nil, 0, 0, err
	}
	query, whereArgs := mysqlstorage.ConstructQueryAndWhereArgs(selectAuditLogsV2SQL, options)
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
	countQuery, countWhereArgs := mysqlstorage.ConstructCountQuery(selectAuditLogV2CountSQL, options)
	if err := s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return auditLogs, nextOffset, totalCount, nil
}

func listAuditLogsOptionsFromParams(p v2als.ListAuditLogsParams) (*mysqlstorage.ListOptions, error) {
	var filters []*mysqlstorage.FilterV2
	if p.EnvironmentID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "environment_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		})
	}
	if p.EntityType != nil {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "entity_type",
			Operator: mysqlstorage.OperatorEqual,
			Value:    *p.EntityType,
		})
	}
	if p.EntityID != "" {
		filters = append(filters, &mysqlstorage.FilterV2{
			Column:   "entity_id",
			Operator: mysqlstorage.OperatorEqual,
			Value:    p.EntityID,
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
	case proto.ListAuditLogsRequest_DEFAULT,
		proto.ListAuditLogsRequest_TIMESTAMP:
		column = "timestamp"
	default:
		return nil, v2als.ErrInvalidOrderBy
	}
	direction := mysqlstorage.OrderDirectionDesc
	if p.OrderDirection == proto.ListAuditLogsRequest_ASC {
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
