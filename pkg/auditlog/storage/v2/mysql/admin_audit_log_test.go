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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/domain"
	v2als "github.com/bucketeer-io/bucketeer/v2/pkg/auditlog/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/auditlog"
)

func TestNewAdminAuditLogStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAdminAuditLogStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &adminAuditLogStorage{}, storage)
}

func TestCreateAdminAuditLogMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*adminAuditLogStorage)
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			expectedErr: v2als.ErrAdminAuditLogAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &adminAuditLogStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
			err := storage.CreateAdminAuditLog(
				context.Background(),
				&domain.AuditLog{AuditLog: &proto.AuditLog{}},
			)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreateAdminAuditLogsMySQL(t *testing.T) {
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
			desc:        "empty slice returns nil",
			setup:       nil,
			auditLogs:   []*domain.AuditLog{},
			expectedErr: nil,
		},
		{
			desc: "error",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			auditLogs: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{}},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "error: duplicate entry",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			auditLogs: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{}},
			},
			expectedErr: v2als.ErrAdminAuditLogAlreadyExists,
		},
		{
			desc: "success",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			auditLogs: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{}},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &adminAuditLogStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAdminAuditLogs(context.Background(), p.auditLogs)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAdminAuditLogsMySQL(t *testing.T) {
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
			desc: "error: query",
			setup: func(s *adminAuditLogStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params: v2als.ListAdminAuditLogsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "error: count",
			setup: func(s *adminAuditLogStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("count error"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			params: v2als.ListAdminAuditLogsParams{
				PageSize: 10,
				Cursor:   "0",
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("count error"),
		},
		{
			desc: "success",
			setup: func(s *adminAuditLogStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
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
			storage := &adminAuditLogStorage{qe: mock.NewMockClient(mockController)}
			p.setup(storage)
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
