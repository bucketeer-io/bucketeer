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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
)

var errInternal = errors.New("internal error")

func TestNewAuditLogStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAuditLogStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &auditLogStorage{}, storage)
}

func TestGetAuditLogPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*auditLogStorage)
		id          string
		expectedErr error
	}{
		{
			desc: "ErrAuditLogNotFound",
			setup: func(s *auditLogStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "audit-log-id",
			expectedErr: v2als.ErrAuditLogNotFound,
		},
		{
			desc: "Error",
			setup: func(s *auditLogStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errInternal)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "audit-log-id",
			expectedErr: errInternal,
		},
		{
			desc: "Success",
			setup: func(s *auditLogStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "audit-log-id",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetAuditLog(context.Background(), p.id, "env-1")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateAuditLogPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*auditLogStorage)
		auditLog    *domain.AuditLog
		expectedErr error
	}{
		{
			desc: "Error",
			setup: func(s *auditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			auditLog:    &domain.AuditLog{AuditLog: &proto.AuditLog{}},
			expectedErr: errInternal,
		},
		{
			desc: "ErrAuditLogAlreadyExists",
			setup: func(s *auditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			auditLog:    &domain.AuditLog{AuditLog: &proto.AuditLog{}},
			expectedErr: v2als.ErrAuditLogAlreadyExists,
		},
		{
			desc: "Success",
			setup: func(s *auditLogStorage) {
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
			storage := newAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAuditLog(context.Background(), p.auditLog)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateAuditLogsPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*auditLogStorage)
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
			setup: func(s *auditLogStorage) {
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
			desc: "ErrAuditLogAlreadyExists",
			setup: func(s *auditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			auditLogs: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: "id-1"}},
			},
			expectedErr: v2als.ErrAuditLogAlreadyExists,
		},
		{
			desc: "Success",
			setup: func(s *auditLogStorage) {
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
			storage := newAuditLogStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAuditLogs(context.Background(), p.auditLogs)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAuditLogsPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*auditLogStorage)
		params         v2als.ListAuditLogsParams
		expected       []*proto.AuditLog
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "QueryError",
			setup: func(s *auditLogStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errInternal)
			},
			params: v2als.ListAuditLogsParams{
				PageSize:      10,
				Cursor:        "0",
				EnvironmentID: "env-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "CountError",
			setup: func(s *auditLogStorage) {
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
			params: v2als.ListAuditLogsParams{
				PageSize:      10,
				Cursor:        "0",
				EnvironmentID: "env-1",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errInternal,
		},
		{
			desc: "Success",
			setup: func(s *auditLogStorage) {
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
			params: v2als.ListAuditLogsParams{
				PageSize:      10,
				Cursor:        "0",
				EnvironmentID: "env-1",
			},
			expected:       []*proto.AuditLog{},
			expectedCursor: 0,
			expectedErr:    nil,
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
				p.params,
			)
			assert.Equal(t, p.expected, auditLogs)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newAuditLogStorageWithMock(t *testing.T, mockController *gomock.Controller) *auditLogStorage {
	t.Helper()
	return &auditLogStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}
