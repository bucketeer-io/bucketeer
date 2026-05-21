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

	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewFeatureLastUsedInfoStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewFeatureLastUsedInfoStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &featureLastUsedInfoStorage{}, storage)
}

func TestFeatureLastUsedInfoStorageGetFeatureLastUsedInfos(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *featureLastUsedInfoStorage)
		ids         []string
		envID       string
		expected    []*domain.FeatureLastUsedInfo
		expectedErr error
	}{
		{
			desc: "error: query fails",
			setup: func(s *featureLastUsedInfoStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("query error"))
			},
			ids:         []string{"id1"},
			envID:       "env",
			expected:    nil,
			expectedErr: errors.New("query error"),
		},
		{
			desc: "success: empty result",
			setup: func(s *featureLastUsedInfoStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			ids:         []string{"id1"},
			envID:       "env",
			expected:    []*domain.FeatureLastUsedInfo{},
			expectedErr: nil,
		},
		{
			desc: "error: scan fails",
			setup: func(s *featureLastUsedInfoStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any()).Return(errors.New("scan error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			ids:         []string{"id1"},
			envID:       "env",
			expected:    nil,
			expectedErr: errors.New("scan error"),
		},
		{
			desc: "success: one result",
			setup: func(s *featureLastUsedInfoStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				rows.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			ids:   []string{"id1"},
			envID: "env",
			expected: []*domain.FeatureLastUsedInfo{
				{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{}},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &featureLastUsedInfoStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			actual, err := storage.GetFeatureLastUsedInfos(context.Background(), p.ids, p.envID)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestFeatureLastUsedInfoStorageUpsertFeatureLastUsedInfo(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(storage *featureLastUsedInfoStorage)
		flui        *domain.FeatureLastUsedInfo
		envID       string
		expectedErr error
	}{
		{
			desc: "error",
			setup: func(s *featureLastUsedInfoStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("exec error"))
			},
			flui: &domain.FeatureLastUsedInfo{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
				FeatureId: "fid",
				Version:   1,
			}},
			envID:       "env",
			expectedErr: errors.New("exec error"),
		},
		{
			desc: "success",
			setup: func(s *featureLastUsedInfoStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			flui: &domain.FeatureLastUsedInfo{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
				FeatureId: "fid",
				Version:   1,
			}},
			envID:       "env",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &featureLastUsedInfoStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(storage)
			err := storage.UpsertFeatureLastUsedInfo(context.Background(), p.flui, p.envID)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
