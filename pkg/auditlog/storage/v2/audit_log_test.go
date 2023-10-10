// Copyright 2023 The Bucketeer Authors.
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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/auditlog/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/auditlog"
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
				{AuditLog: &proto.AuditLog{Id: "id-0"}},
				{AuditLog: &proto.AuditLog{Id: "id-1"}},
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
				{AuditLog: &proto.AuditLog{Id: "id-0"}},
				{AuditLog: &proto.AuditLog{Id: "id-1"}},
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
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: []*domain.AuditLog{
				{AuditLog: &proto.AuditLog{Id: "id-0"}},
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
			err := storage.CreateAuditLogs(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAuditLogs(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*auditLogStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.AuditLog
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *auditLogStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:     nil,
			orders:         nil,
			limit:          0,
			offset:         0,
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *auditLogStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			whereParts: nil,
			orders: []*mysql.Order{
				mysql.NewOrder("timestamp", mysql.OrderDirectionDesc),
			},
			limit:          10,
			offset:         5,
			expected:       []*proto.AuditLog{},
			expectedCursor: 5,
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
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, auditLogs)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newAuditLogStorageWithMock(t *testing.T, mockController *gomock.Controller) *auditLogStorage {
	t.Helper()
	return &auditLogStorage{mock.NewMockQueryExecer(mockController)}
}
