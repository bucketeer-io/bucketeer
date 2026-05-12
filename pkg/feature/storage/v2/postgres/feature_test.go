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
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	pgmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewFeatureStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewFeatureStorage(pgmock.NewMockQueryExecer(mockController))
	assert.IsType(t, &featureStorage{}, storage)
}

func TestCreateFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*featureStorage)
		feature     *domain.Feature
		expectedErr error
	}{
		{
			desc: "ErrFeatureAlreadyExists",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, postgres.ErrDuplicateEntry)
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: v2fs.ErrFeatureAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateFeature(context.Background(), p.feature, "env")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*featureStorage)
		feature     *domain.Feature
		expectedErr error
	}{
		{
			desc: "ErrFeatureUnexpectedAffectedRows",
			setup: func(s *featureStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: v2fs.ErrFeatureUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *featureStorage) {
				result := pgmock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateFeature(context.Background(), p.feature, "env")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetFeature(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*featureStorage)
		id          string
		expected    *domain.Feature
		expectedErr error
	}{
		{
			desc: "ErrFeatureNotFound",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "feature-id",
			expected:    nil,
			expectedErr: v2fs.ErrFeatureNotFound,
		},
		{
			desc: "Error",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "feature-id",
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id: "feature-id",
			expected: &domain.Feature{Feature: &proto.Feature{
				AutoOpsSummary: &proto.AutoOpsSummary{},
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			feature, err := storage.GetFeature(context.Background(), p.id, "env")
			assert.Equal(t, p.expected, feature)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetFeatureByVersion(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*featureStorage)
		id          string
		version     int32
		expected    *domain.Feature
		expectedErr error
	}{
		{
			desc: "ErrFeatureNotFound",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(postgres.ErrNoRows)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "feature-id",
			version:     1,
			expected:    nil,
			expectedErr: v2fs.ErrFeatureNotFound,
		},
		{
			desc: "Error",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "feature-id",
			version:     1,
			expected:    nil,
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:      "feature-id",
			version: 1,
			expected: &domain.Feature{Feature: &proto.Feature{
				AutoOpsSummary: &proto.AutoOpsSummary{},
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			feature, err := storage.GetFeatureByVersion(context.Background(), p.id, p.version, "env")
			assert.Equal(t, p.expected, feature)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*featureStorage)
		params         v2fs.ListFeaturesParams
		expected       []*proto.Feature
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params:         v2fs.ListFeaturesParams{},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *featureStorage) {
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
			params: v2fs.ListFeaturesParams{
				EnvironmentID: "env",
			},
			expected:       []*proto.Feature{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			features, cursor, _, err := storage.ListFeatures(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, features)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListFeaturesFilteredByExperiment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*featureStorage)
		params         v2fs.ListFeaturesFilteredByExperimentParams
		expected       []*proto.Feature
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			params:         v2fs.ListFeaturesFilteredByExperimentParams{},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *featureStorage) {
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
			params: v2fs.ListFeaturesFilteredByExperimentParams{
				ListFeaturesParams: v2fs.ListFeaturesParams{
					EnvironmentID: "env",
				},
			},
			expected:       []*proto.Feature{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			features, cursor, _, err := storage.ListFeaturesFilteredByExperiment(
				context.Background(),
				p.params,
			)
			assert.Equal(t, p.expected, features)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListFeaturesByEnvironment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*featureStorage)
		environmentID string
		expected      []*proto.Feature
		expectedErr   error
	}{
		{
			desc: "Error: query fails",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("query error"))
			},
			environmentID: "env-id",
			expected:      nil,
			expectedErr:   errors.New("query error"),
		},
		{
			desc: "Success: empty result",
			setup: func(s *featureStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), "env-id",
				).Return(rows, nil)
			},
			environmentID: "env-id",
			expected:      []*proto.Feature{},
			expectedErr:   nil,
		},
		{
			desc: "Error: scan fails",
			setup: func(s *featureStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any()).Return(errors.New("scan error"))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), "env-id",
				).Return(rows, nil)
			},
			environmentID: "env-id",
			expected:      nil,
			expectedErr:   errors.New("scan error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			features, err := storage.ListFeaturesByEnvironment(context.Background(), p.environmentID)
			assert.Equal(t, p.expected, features)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListAllEnvironmentFeatures(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*featureStorage)
		expected    []*proto.EnvironmentFeature
		expectedErr error
	}{
		{
			desc: "Error: query fails",
			setup: func(s *featureStorage) {
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("query error"))
			},
			expected:    nil,
			expectedErr: errors.New("query error"),
		},
		{
			desc: "Success: empty result",
			setup: func(s *featureStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    []*proto.EnvironmentFeature{},
			expectedErr: nil,
		},
		{
			desc: "Error: scan fails",
			setup: func(s *featureStorage) {
				rows := pgmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(true)
				rows.EXPECT().Scan(gomock.Any()).Return(errors.New("scan error"))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			expected:    nil,
			expectedErr: errors.New("scan error"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			envFeatures, err := storage.ListAllEnvironmentFeatures(context.Background())
			assert.Equal(t, p.expected, envFeatures)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetFeatureSummary(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*featureStorage)
		environmentID string
		expected      *proto.FeatureSummary
		expectedErr   error
	}{
		{
			desc: "Error",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			environmentID: "env",
			expected:      nil,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *featureStorage) {
				row := pgmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*pgmock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			environmentID: "env",
			expected:      &proto.FeatureSummary{},
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newFeatureStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			summary, err := storage.GetFeatureSummary(context.Background(), p.environmentID)
			assert.Equal(t, p.expected, summary)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newFeatureStorageWithMock(t *testing.T, mockController *gomock.Controller) *featureStorage {
	t.Helper()
	return &featureStorage{qe: pgmock.NewMockQueryExecer(mockController)}
}
