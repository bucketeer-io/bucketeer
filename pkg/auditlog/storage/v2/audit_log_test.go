// Copyright 2024 The Bucketeer Authors.
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

package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/bucketeer-io/bucketeer/pkg/auditlog/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/auditlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewSubscriptionStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAuditLogStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &auditLogStorage{}, storage)
}

func TestCreateAuditLogs(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id0 := "id-0"
	id1 := "id-1"
	patterns := []struct {
		desc        string
		setup       func(*auditLogStorage)
		input       []*domain.AuditLog
		expectedErr error
	}{
		{
			desc: "ErrAuditLogAlreadyExists",
			setup: func(s *auditLogStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: id0}},
				{AuditLog: &proto.AuditLog{Id: id1}},
			},
			expectedErr: ErrAuditLogAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *auditLogStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: id0}},
				{AuditLog: &proto.AuditLog{Id: id1}},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc:        "Success: len == 0",
			setup:       nil,
			input:       nil,
			expectedErr: nil,
		},
		{
			desc: "Success",
			setup: func(s *auditLogStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Regex("^INSERT INTO audit_log\\s+\\(\\s*id,\\s*timestamp,\\s*entity_type,\\s*entity_id,\\s*type,\\s*event,\\s*editor,\\s*options,\\s*environment_namespace\\s*\\)\\s+VALUES\\s*\\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\),\\s*\\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)$"),
					[]interface{}{
						id0, int64(1), int32(2), "e0", int32(3), mysql.JSONObject{Val: nil}, mysql.JSONObject{Val: nil}, mysql.JSONObject{Val: nil}, "ns0",
						id1, int64(10), int32(3), "e2", int32(4), mysql.JSONObject{Val: nil}, mysql.JSONObject{Val: nil}, mysql.JSONObject{Val: nil}, "ns1",
					},
				).Return(nil, nil)
			},
			input: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: id0, Timestamp: 1, EntityType: 2, EntityId: "e0", Type: 3}, EnvironmentNamespace: "ns0"},
				{AuditLog: &proto.AuditLog{Id: id1, Timestamp: 10, EntityType: 3, EntityId: "e2", Type: 4}, EnvironmentNamespace: "ns1"},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAuditLogs(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAuditLogs(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	getSize := 2
	offset := 5
	limit := 10
	timestamp := 4
	entityType := 2

	patterns := []struct {
		desc                string
		setup               func(*auditLogStorage)
		whereParts          []mysql.WherePart
		orders              []*mysql.Order
		limit               int
		offset              int
		expectedResultCount int
		expectedCursor      int
		expectedErr         error
	}{
		{
			desc: "Error",
			setup: func(s *auditLogStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:          nil,
			orders:              nil,
			limit:               0,
			offset:              0,
			expectedResultCount: 0,
			expectedCursor:      0,
			expectedErr:         errors.New("error"),
		},
		{
			desc: "Success:No wereParts and no orderParts and no limit and no offset",
			setup: func(s *auditLogStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Regex("^SELECT\\s+id,\\s*timestamp,\\s*entity_type,\\s*entity_id,\\s*type,\\s*event,\\s*editor,\\s*options\\s+FROM\\s+audit_log\\s*$"),
					[]interface{}{},
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Regex("^SELECT\\s+COUNT\\(1\\)\\s+FROM\\s+audit_log\\s*$"),
					[]interface{}{},
				).Return(row)
			},
			whereParts:          nil,
			orders:              nil,
			limit:               0,
			offset:              0,
			expectedResultCount: 0,
			expectedCursor:      0,
			expectedErr:         nil,
		},
		{
			desc: "Success",
			setup: func(s *auditLogStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					return nextCallCount <= getSize
				}).Times(getSize + 1)
				rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(getSize)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Regex("^SELECT\\s+id,\\s*timestamp,\\s*entity_type,\\s*entity_id,\\s*type,\\s*event,\\s*editor,\\s*options\\s+FROM\\s+audit_log\\s+WHERE timestamp >= \\? AND entity_type = \\? ORDER BY id ASC, timestamp DESC LIMIT 10 OFFSET 5\\s*$"),
					timestamp, entityType,
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Regex("^SELECT\\s+COUNT\\(1\\)\\s+FROM\\s+audit_log\\s+WHERE timestamp >= \\? AND entity_type = \\? ORDER BY id ASC, timestamp DESC\\s*$"),
					timestamp, entityType,
				).Return(row)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("timestamp", ">=", timestamp),
				mysql.NewFilter("entity_type", "=", entityType),
			},

			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
				mysql.NewOrder("timestamp", mysql.OrderDirectionDesc),
			},
			limit:               limit,
			offset:              offset,
			expectedResultCount: getSize,
			expectedCursor:      offset + getSize,
			expectedErr:         nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			auditLogs, cursor, _, err := storage.ListAuditLogs(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expectedResultCount, len(auditLogs))
			if len(auditLogs) > 0 {
				assert.IsType(t, auditLogs, []*proto.AuditLog{})
			}
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newAuditLogStorageWithMock(t *testing.T, mockController *gomock.Controller) *auditLogStorage {
	t.Helper()
	return &auditLogStorage{mock.NewMockQueryExecer(mockController)}
}
