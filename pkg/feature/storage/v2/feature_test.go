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

func TestNewFeatureStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewFeatureStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &featureStorage{}, storage)
}

func TestCreateFeatureMySQL(t *testing.T) {
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
			desc: "error",
			setup: func(s *featureStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *featureStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &featureStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateFeature(context.Background(), p.feature, "env")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateFeatureMySQL(t *testing.T) {
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
			desc: "error",
			setup: func(s *featureStorage) {
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *featureStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
			},
			feature:     &domain.Feature{Feature: &proto.Feature{}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &featureStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateFeature(context.Background(), p.feature, "env")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*featureStorage)
		featureID      string
		expected       *domain.Feature
		expectedErr    error
		expectedCalled bool
	}{
		{
			desc: "error",
			setup: func(s *featureStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("feature not found"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			featureID:      "feature",
			expected:       nil,
			expectedErr:    errors.New("feature not found"),
			expectedCalled: true,
		},
		{
			desc: "success",
			setup: func(s *featureStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			featureID: "feature",
			expected: &domain.Feature{Feature: &proto.Feature{
				AutoOpsSummary: &proto.AutoOpsSummary{},
			}},
			expectedErr:    nil,
			expectedCalled: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &featureStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			feature, err := storage.GetFeature(context.Background(), p.featureID, "env")
			assert.Equal(t, p.expected, feature)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListFeatureMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*featureStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.Feature
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "error",
			setup: func(s *featureStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
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
			desc: "success",
			setup: func(s *featureStorage) {
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
			whereParts: []mysql.WherePart{
				mysql.NewFilter("environment_id", "=", "env"),
			},
			expected:       []*proto.Feature{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &featureStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			features, cursor, _, err := storage.ListFeatures(
				context.Background(),
				p.whereParts,
				proto.FeatureLastUsedInfo_UNKNOWN,
				"env",
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, features)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListFeatureFilterByExperiment(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		setup          func(*featureStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.Feature
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "error",
			setup: func(s *featureStorage) {
				s.qe.(*mock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *featureStorage) {
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
			expected:       []*proto.Feature{},
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &featureStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			features, cursor, _, err := storage.ListFeaturesFilteredByExperiment(
				context.Background(),
				p.whereParts,
				proto.FeatureLastUsedInfo_UNKNOWN,
				"env",
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, features)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCountFeaturesByStatus(t *testing.T) {
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
			desc: "error",
			setup: func(s *featureStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			environmentID: "env",
			expected:      nil,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "success",
			setup: func(s *featureStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockClient).EXPECT().QueryRowContext(
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
			storage := &featureStorage{qe: mock.NewMockClient(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			count, err := storage.GetFeatureSummary(context.Background(), p.environmentID)
			assert.Equal(t, p.expected, count)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
