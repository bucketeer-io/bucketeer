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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestNewHoldoutStorage(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	storage := NewHoldoutStorage(mock.NewMockClient(controller))
	assert.IsType(t, &holdoutStorage{}, storage)
}

func TestCreateHoldoutGroup(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	group := &domain.HoldoutGroup{HoldoutGroup: &proto.HoldoutGroup{Id: "hg"}}

	t.Run("duplicate", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(nil, mysql.ErrDuplicateEntry)
		err := storage.CreateHoldoutGroup(context.Background(), group)
		assert.ErrorIs(t, err, v2hs.ErrHoldoutGroupAlreadyExists)
	})
	t.Run("success", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(nil, nil)
		assert.NoError(t, storage.CreateHoldoutGroup(context.Background(), group))
	})
}

func TestUpdateHoldoutGroup(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	group := &domain.HoldoutGroup{HoldoutGroup: &proto.HoldoutGroup{Id: "hg"}}

	t.Run("unexpected rows", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		result := mock.NewMockResult(controller)
		result.EXPECT().RowsAffected().Return(int64(0), nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(result, nil)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		err := storage.UpdateHoldoutGroup(context.Background(), group)
		assert.ErrorIs(t, err, v2hs.ErrHoldoutGroupUnexpectedAffectedRows)
	})
	t.Run("unchanged existing row", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		result := mock.NewMockResult(controller)
		result.EXPECT().RowsAffected().Return(int64(0), nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(result, nil)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		assert.NoError(t, storage.UpdateHoldoutGroup(context.Background(), group))
	})
	t.Run("success", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		result := mock.NewMockResult(controller)
		result.EXPECT().RowsAffected().Return(int64(1), nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(result, nil)
		assert.NoError(t, storage.UpdateHoldoutGroup(context.Background(), group))
	})
}

func TestGetHoldoutGroup(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	t.Run("not found", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		group, err := storage.GetHoldoutGroup(context.Background(), "hg", "env")
		assert.Nil(t, group)
		assert.ErrorIs(t, err, v2hs.ErrHoldoutGroupNotFound)
	})
	t.Run("success", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		group, err := storage.GetHoldoutGroup(context.Background(), "hg", "env")
		require.NoError(t, err)
		assert.NotNil(t, group)
	})
	t.Run("for update", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		_, err := storage.GetHoldoutGroupForUpdate(context.Background(), "hg", "env")
		assert.ErrorIs(t, err, v2hs.ErrHoldoutGroupNotFound)
	})
}

func TestListHoldoutGroups(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	archived := true

	t.Run("invalid cursor", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		_, cursor, total, err := storage.ListHoldoutGroups(context.Background(), v2hs.ListHoldoutGroupsParams{
			Cursor: "nope",
		})
		assert.Zero(t, cursor)
		assert.Zero(t, total)
		assert.ErrorIs(t, err, v2hs.ErrInvalidCursor)
	})
	t.Run("invalid order", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		_, _, _, err := storage.ListHoldoutGroups(context.Background(), v2hs.ListHoldoutGroupsParams{
			OrderBy: proto.ListHoldoutGroupsRequest_OrderBy(-1),
		})
		assert.ErrorIs(t, err, v2hs.ErrInvalidOrderBy)
	})
	orders := []proto.ListHoldoutGroupsRequest_OrderBy{
		proto.ListHoldoutGroupsRequest_DEFAULT,
		proto.ListHoldoutGroupsRequest_NAME,
		proto.ListHoldoutGroupsRequest_CREATED_AT,
		proto.ListHoldoutGroupsRequest_UPDATED_AT,
		proto.ListHoldoutGroupsRequest_STATUS,
	}
	for _, order := range orders {
		t.Run(order.String(), func(t *testing.T) {
			storage := newHoldoutStorageWithMock(t, controller)
			rows := mock.NewMockRows(controller)
			rows.EXPECT().Close().Return(nil)
			rows.EXPECT().Next().Return(false)
			rows.EXPECT().Err().Return(nil)
			storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
				gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(rows, nil)
			count := mock.NewMockRow(controller)
			count.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
				*(dest[0].(*int64)) = 4
				return nil
			})
			storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
				gomock.Any(), gomock.Any(), gomock.Any(),
			).Return(count)
			groups, cursor, total, err := storage.ListHoldoutGroups(context.Background(), v2hs.ListHoldoutGroupsParams{
				EnvironmentID:  "env",
				Archived:       &archived,
				Statuses:       []proto.HoldoutGroup_Status{proto.HoldoutGroup_RUNNING},
				SearchKeyword:  "hold",
				OrderBy:        order,
				OrderDirection: proto.ListHoldoutGroupsRequest_DESC,
				PageSize:       10,
				Cursor:         "2",
			})
			require.NoError(t, err)
			assert.Empty(t, groups)
			assert.Equal(t, 2, cursor)
			assert.Equal(t, int64(4), total)
		})
	}
}

func TestHoldoutGroupFeatures(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	feature := &domain.HoldoutGroupFeature{HoldoutGroupFeature: &proto.HoldoutGroupFeature{Id: "row"}}

	t.Run("create empty", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		assert.NoError(t, storage.CreateHoldoutGroupFeatures(context.Background(), nil))
	})
	t.Run("create duplicate", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(nil, mysql.ErrDuplicateEntry)
		err := storage.CreateHoldoutGroupFeatures(context.Background(), []*domain.HoldoutGroupFeature{feature})
		assert.ErrorIs(t, err, v2hs.ErrHoldoutGroupAlreadyExists)
	})
	t.Run("create success", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(nil, nil)
		assert.NoError(t, storage.CreateHoldoutGroupFeatures(context.Background(), []*domain.HoldoutGroupFeature{feature}))
	})
	t.Run("remove missing", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		result := mock.NewMockResult(controller)
		result.EXPECT().RowsAffected().Return(int64(0), nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(result, nil)
		err := storage.RemoveHoldoutGroupFeature(context.Background(), "env", "hg", "flag", 10)
		assert.ErrorIs(t, err, v2hs.ErrHoldoutFeatureNotFound)
	})
	t.Run("remove success", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		result := mock.NewMockResult(controller)
		result.EXPECT().RowsAffected().Return(int64(1), nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(result, nil)
		assert.NoError(t, storage.RemoveHoldoutGroupFeature(context.Background(), "env", "hg", "flag", 10))
	})
	t.Run("list current none", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		features, err := storage.ListCurrentFeatures(context.Background(), "env", nil)
		require.NoError(t, err)
		assert.Nil(t, features)
	})
	t.Run("list current", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		rows := mock.NewMockRows(controller)
		rows.EXPECT().Close().Return(nil)
		rows.EXPECT().Next().Return(true)
		rows.EXPECT().Next().Return(false)
		rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = "row"
			*(dest[1].(*string)) = "hg"
			*(dest[2].(*string)) = "flag"
			*(dest[3].(*string)) = "env"
			*(dest[4].(*string)) = "var"
			*(dest[5].(*int64)) = 10
			*(dest[6].(*int64)) = 0
			return nil
		})
		rows.EXPECT().Err().Return(nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(rows, nil)
		features, err := storage.ListCurrentFeatures(context.Background(), "env", []string{"hg"})
		require.NoError(t, err)
		require.Len(t, features, 1)
		assert.Equal(t, "flag", features[0].FeatureId)
	})
}

func TestEnvironmentLockAndSeed(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	t.Run("environment missing", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		assert.ErrorIs(t, storage.LockEnvironment(context.Background(), "env"), v2hs.ErrEnvironmentNotFound)
	})
	t.Run("environment locked", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = "env"
			return nil
		})
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		assert.NoError(t, storage.LockEnvironment(context.Background(), "env"))
	})
	t.Run("no seed yet", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		seed, err := storage.GetEnvironmentSamplingSeed(context.Background(), "env")
		require.NoError(t, err)
		assert.Empty(t, seed)
	})
	t.Run("existing seed", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = "seed"
			return nil
		})
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		seed, err := storage.GetEnvironmentSamplingSeed(context.Background(), "env")
		require.NoError(t, err)
		assert.Equal(t, "seed", seed)
	})
}

func TestListActiveGroupsAndTargets(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	t.Run("active groups", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		rows := mock.NewMockRows(controller)
		rows.EXPECT().Close().Return(nil)
		rows.EXPECT().Next().Return(true)
		rows.EXPECT().Next().Return(false)
		rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = "hg"
			*(dest[1].(*string)) = "name"
			*(dest[2].(*int32)) = int32(proto.HoldoutGroup_RUNNING)
			*(dest[3].(*int64)) = 10
			*(dest[4].(*int64)) = 20
			*(dest[5].(*string)) = "seg"
			return nil
		})
		rows.EXPECT().Err().Return(nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(rows, nil)
		groups, err := storage.ListActiveGroups(context.Background(), "env")
		require.NoError(t, err)
		require.Len(t, groups, 1)
		assert.Equal(t, proto.HoldoutGroup_RUNNING, groups[0].Status)
		assert.Equal(t, "seg", groups[0].SegmentID)
	})
	t.Run("live targets", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		rows := mock.NewMockRows(controller)
		rows.EXPECT().Close().Return(nil)
		rows.EXPECT().Next().Return(true)
		rows.EXPECT().Next().Return(false)
		rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = "hg"
			*(dest[1].(*string)) = "name"
			*(dest[2].(*int32)) = int32(proto.HoldoutGroup_ANALYSIS)
			*(dest[3].(*string)) = "flag"
			*(dest[4].(*string)) = "var"
			return nil
		})
		rows.EXPECT().Err().Return(nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(rows, nil)
		targets, err := storage.ListLiveFeatureTargets(context.Background(), "env", "flag")
		require.NoError(t, err)
		require.Len(t, targets, 1)
		assert.Equal(t, "var", targets[0].VariationID)
		assert.Equal(t, proto.HoldoutGroup_ANALYSIS, targets[0].Status)
	})
}

func TestGetLiveServingGroup(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	t.Run("not found", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(),
			"env", false, false,
			int32(proto.HoldoutGroup_RUNNING), int32(proto.HoldoutGroup_ANALYSIS),
		).Return(row)
		group, err := storage.GetLiveServingGroup(context.Background(), "env")
		assert.Nil(t, group)
		assert.ErrorIs(t, err, v2hs.ErrHoldoutGroupNotFound)
	})
	t.Run("attaches current features", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = "hg-live"
			*(dest[1].(*string)) = "env"
			*(dest[9].(*int32)) = int32(proto.HoldoutGroup_ANALYSIS)
			return nil
		})
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(row)
		rows := mock.NewMockRows(controller)
		rows.EXPECT().Close().Return(nil)
		rows.EXPECT().Next().Return(true)
		rows.EXPECT().Next().Return(false)
		rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = "row"
			*(dest[1].(*string)) = "hg-live"
			*(dest[2].(*string)) = "flag"
			*(dest[3].(*string)) = "env"
			*(dest[4].(*string)) = "hold"
			*(dest[5].(*int64)) = int64(10)
			*(dest[6].(*int64)) = int64(0)
			return nil
		})
		rows.EXPECT().Err().Return(nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
			gomock.Any(), gomock.Any(), "env", "hg-live",
		).Return(rows, nil)
		group, err := storage.GetLiveServingGroup(context.Background(), "env")
		require.NoError(t, err)
		require.NotNil(t, group)
		assert.Equal(t, "hg-live", group.Id)
		assert.Equal(t, proto.HoldoutGroup_ANALYSIS, group.Status)
		require.Len(t, group.Features, 1)
		assert.Equal(t, "flag", group.Features[0].FeatureId)
		assert.Equal(t, "hold", group.Features[0].HoldoutVariationId)
	})
}

func TestGetClearedHoldoutUpdatedAt(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)

	t.Run("no stopped group", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(),
			"env",
			int32(proto.HoldoutGroup_STOPPED),
			int32(proto.HoldoutGroup_FORCE_STOPPED),
		).Return(row)
		updatedAt, err := storage.GetClearedHoldoutUpdatedAt(context.Background(), "env")
		require.NoError(t, err)
		assert.Equal(t, int64(0), updatedAt)
	})
	t.Run("newest stop time", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		row := mock.NewMockRow(controller)
		row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*int64)) = int64(1700000000)
			return nil
		})
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
			gomock.Any(), gomock.Any(),
			"env",
			int32(proto.HoldoutGroup_STOPPED),
			int32(proto.HoldoutGroup_FORCE_STOPPED),
		).Return(row)
		updatedAt, err := storage.GetClearedHoldoutUpdatedAt(context.Background(), "env")
		require.NoError(t, err)
		assert.Equal(t, int64(1700000000), updatedAt)
	})
}

func newHoldoutStorageWithMock(t *testing.T, controller *gomock.Controller) *holdoutStorage {
	t.Helper()
	return &holdoutStorage{qe: mock.NewMockQueryExecer(controller)}
}
