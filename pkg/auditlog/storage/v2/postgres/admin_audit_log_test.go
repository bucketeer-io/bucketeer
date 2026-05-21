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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
)

func TestNewAdminAuditLogStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAdminAuditLogStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &adminAuditLogStorage{}, storage)
}

func TestCreateAdminAuditLogPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*adminAuditLogStorage)
		auditLog    *domain.AuditLog
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			auditLog:    &domain.AuditLog{AuditLog: &proto.AuditLog{}},
			expectedErr: errInternal,
		},
		{
			desc: "ErrAdminAuditLogAlreadyExists",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			auditLog:    &domain.AuditLog{AuditLog: &proto.AuditLog{}},
			expectedErr: v2als.ErrAdminAuditLogAlreadyExists,
		},
		{
			desc: "Success",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			auditLog:    &domain.AuditLog{AuditLog: &proto.AuditLog{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAdminAuditLog(context.Background(), p.auditLog)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateAdminAuditLogsPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*adminAuditLogStorage)
		auditLogs   []*domain.AuditLog
		expectedErr error
	}{
		{
			desc:        "Empty",
			setup:       nil,
			auditLogs:   []*domain.AuditLog{},
			expectedErr: nil,
		},
		{
			desc: "Error",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			auditLogs: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: "id-1"}},
			},
			expectedErr: errInternal,
		},
		{
			desc: "ErrAdminAuditLogAlreadyExists",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			auditLogs: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: "id-1"}},
			},
			expectedErr: v2als.ErrAdminAuditLogAlreadyExists,
		},
		{
			desc: "Success",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			auditLogs: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: "id-1"}},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAdminAuditLogs(context.Background(), p.auditLogs)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAdminAuditLogsPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*adminAuditLogStorage)
		params         v2als.ListAdminAuditLogsParams
		expected       []*proto.AuditLog
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "QueryError",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			params: v2als.ListAdminAuditLogsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "CountError",
			setup: func(s *adminAuditLogStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2als.ListAdminAuditLogsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "Success",
			setup: func(s *adminAuditLogStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2als.ListAdminAuditLogsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       []*proto.AuditLog{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			auditLogs, cursor, _, err := storage.ListAdminAuditLogs(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, auditLogs)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newAdminAuditLogStorageWithMock(t *testing.T, mockController *gomock.Controller) *adminAuditLogStorage {
	t.Helper()
	return &adminAuditLogStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}
