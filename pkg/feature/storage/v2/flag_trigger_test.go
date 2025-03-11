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
//

package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewFlagTriggerStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewFlagTriggerStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &flagTriggerStorage{}, storage)
}

func TestFlagTriggerStorageCreateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *flagTriggerStorage)
		flagTrigger *domain.FlagTrigger
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *flagTriggerStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *flagTriggerStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &flagTriggerStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.CreateFlagTrigger(context.Background(), p.flagTrigger)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestFlagTriggerStorageUpdateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *flagTriggerStorage)
		flagTrigger *domain.FlagTrigger
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *flagTriggerStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *flagTriggerStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &flagTriggerStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.UpdateFlagTrigger(context.Background(), p.flagTrigger)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestFlagTriggerStorageDeleteFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *flagTriggerStorage)
		flagTrigger *domain.FlagTrigger
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *flagTriggerStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{
				Id:            "id",
				EnvironmentId: "env",
			}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *flagTriggerStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{
				Id:            "id",
				EnvironmentId: "env",
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &flagTriggerStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.DeleteFlagTrigger(context.Background(), p.flagTrigger.Id, p.flagTrigger.EnvironmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestFlagTriggerStorageGetFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *flagTriggerStorage)
		flagTrigger *domain.FlagTrigger
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *flagTriggerStorage) {
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{
				Id:            "id",
				EnvironmentId: "env",
			}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *flagTriggerStorage) {
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
			},
			flagTrigger: &domain.FlagTrigger{FlagTrigger: &proto.FlagTrigger{
				Id:            "id",
				EnvironmentId: "env",
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &flagTriggerStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			_, err := storage.GetFlagTrigger(context.Background(), p.flagTrigger.Id, p.flagTrigger.EnvironmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestFlagTriggerStorageGetFlagTriggerByToken(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *flagTriggerStorage)
		token       string
		expected    *proto.FlagTrigger
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *flagTriggerStorage) {
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
			},
			token:       "token",
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *flagTriggerStorage) {
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
			},
			token:       "token",
			expected:    &proto.FlagTrigger{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &flagTriggerStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			_, err := storage.GetFlagTriggerByToken(context.Background(), p.token)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestFlagTriggerStorageListFlagTriggers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(storage *flagTriggerStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.FlagTrigger
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "error",
			setup: func(s *flagTriggerStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:  []mysql.WherePart{},
			orders:      []*mysql.Order{},
			limit:       0,
			offset:      0,
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *flagTriggerStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			whereParts:     []mysql.WherePart{},
			orders:         []*mysql.Order{},
			limit:          0,
			offset:         0,
			expected:       []*proto.FlagTrigger{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &flagTriggerStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			expected, nextOffset, _, err := storage.ListFlagTriggers(context.Background(), p.whereParts, p.orders, p.limit, p.offset)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, expected)
			assert.Equal(t, p.expectedCursor, nextOffset)
		})
	}
}
