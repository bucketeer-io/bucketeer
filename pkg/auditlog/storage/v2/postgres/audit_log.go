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

package postgres

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	pgstorage "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
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
	qe pgstorage.QueryExecer
}

func NewAuditLogStorage(qe pgstorage.QueryExecer) v2als.AuditLogStorage {
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
		&pgstorage.JSONObject{Val: &auditLog.Event},
		&pgstorage.JSONObject{Val: &auditLog.Editor},
		&pgstorage.JSONObject{Val: &auditLog.Options},
		&auditLog.EntityData,
		&auditLog.PreviousEntityData,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNoRows) {
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
	query.WriteString(insertAuditLogsV2SQL)
	colsPerRow := 11
	for i, al := range auditLogs {
		if i != 0 {
			query.WriteString(",")
		}
		base := i*colsPerRow + 1
		query.WriteString(fmt.Sprintf(
			" ($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base, base+1, base+2, base+3, base+4,
			base+5, base+6, base+7, base+8, base+9, base+10,
		))
		args = append(
			args,
			al.Id,
			al.Timestamp,
			int32(al.EntityType),
			al.EntityId,
			int32(al.Type),
			pgstorage.JSONObject{Val: al.Event},
			pgstorage.JSONObject{Val: al.Editor},
			pgstorage.JSONObject{Val: al.Options},
			al.EnvironmentId,
			al.EntityData,
			al.PreviousEntityData,
		)
	}
	_, err := s.qe.ExecContext(ctx, query.String(), args...)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
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
		pgstorage.JSONObject{Val: auditLog.Event},
		pgstorage.JSONObject{Val: auditLog.Editor},
		pgstorage.JSONObject{Val: auditLog.Options},
		auditLog.EnvironmentId,
		auditLog.EntityData,
		auditLog.PreviousEntityData,
	)
	if err != nil {
		if errors.Is(err, pgstorage.ErrDuplicateEntry) {
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
	query, whereArgs := pgstorage.ConstructQueryAndWhereArgs(selectAuditLogsV2SQL, options)
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
			&pgstorage.JSONObject{Val: &auditLog.Event},
			&pgstorage.JSONObject{Val: &auditLog.Editor},
			&pgstorage.JSONObject{Val: &auditLog.Options},
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
	countQuery, countWhereArgs := pgstorage.ConstructCountQuery(selectAuditLogV2CountSQL, options)
	if err := s.qe.QueryRowContext(ctx, countQuery, countWhereArgs...).Scan(&totalCount); err != nil {
		return nil, 0, 0, err
	}
	return auditLogs, nextOffset, totalCount, nil
}

func listAuditLogsOptionsFromParams(p v2als.ListAuditLogsParams) (*pgstorage.ListOptions, error) {
	var filters []*pgstorage.Filter
	if p.EnvironmentID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "environment_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EnvironmentID,
		})
	}
	if p.EntityType != nil {
		filters = append(filters, &pgstorage.Filter{
			Column:   "entity_type",
			Operator: pgstorage.OperatorEqual,
			Value:    *p.EntityType,
		})
	}
	if p.EntityID != "" {
		filters = append(filters, &pgstorage.Filter{
			Column:   "entity_id",
			Operator: pgstorage.OperatorEqual,
			Value:    p.EntityID,
		})
	}
	if p.From != 0 {
		filters = append(filters, &pgstorage.Filter{
			Column:   "timestamp",
			Operator: pgstorage.OperatorGreaterThanOrEqual,
			Value:    p.From,
		})
	}
	if p.To != 0 {
		filters = append(filters, &pgstorage.Filter{
			Column:   "timestamp",
			Operator: pgstorage.OperatorLessThanOrEqual,
			Value:    p.To,
		})
	}
	var searchQuery *pgstorage.SearchQuery
	if p.SearchKeyword != "" {
		// editor is a JSONB column, so we cast to text for LIKE matching
		searchQuery = &pgstorage.SearchQuery{
			Columns: []string{"editor::text"},
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
	direction := pgstorage.OrderDirectionDesc
	if p.OrderDirection == proto.ListAuditLogsRequest_ASC {
		direction = pgstorage.OrderDirectionAsc
	}
	cursor := p.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, err
	}
	return &pgstorage.ListOptions{
		Limit:       p.PageSize,
		Offset:      offset,
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      []*pgstorage.Order{pgstorage.NewOrder(column, direction)},
	}, nil
}
