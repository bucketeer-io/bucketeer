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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const measuredGroupsEndedSince = int64(1_700_000_000)

// expectMeasuredGroupsQuery expects the group query with the arguments in
// the order of its placeholders: a scheduled stop is compared on stop_at and
// a manual stop on stopped_at, both against endedSince.
func expectMeasuredGroupsQuery(storage *holdoutStorage) *gomock.Call {
	return storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
		gomock.Any(),
		selectMeasuredGroupsSQL,
		"env",
		false,
		int32(proto.HoldoutGroup_RUNNING),
		int32(proto.HoldoutGroup_ANALYSIS),
		int32(proto.HoldoutGroup_STOPPED),
		measuredGroupsEndedSince,
		int32(proto.HoldoutGroup_FORCE_STOPPED),
		measuredGroupsEndedSince,
	)
}

func measuredGroupRows(controller *gomock.Controller, ids ...string) *mock.MockRows {
	rows := mock.NewMockRows(controller)
	rows.EXPECT().Close().Return(nil)
	for _, id := range ids {
		rows.EXPECT().Next().Return(true)
		rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = id
			*(dest[1].(*string)) = "env"
			return nil
		})
	}
	rows.EXPECT().Next().Return(false)
	rows.EXPECT().Err().Return(nil)
	return rows
}

type measuredFeatureRow struct {
	id, groupID, featureID string
	addedAt, removedAt     int64
}

func measuredFeatureRows(controller *gomock.Controller, features ...measuredFeatureRow) *mock.MockRows {
	rows := mock.NewMockRows(controller)
	rows.EXPECT().Close().Return(nil)
	for _, f := range features {
		rows.EXPECT().Next().Return(true)
		rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...any) error {
			*(dest[0].(*string)) = f.id
			*(dest[1].(*string)) = f.groupID
			*(dest[2].(*string)) = f.featureID
			*(dest[3].(*string)) = "env"
			*(dest[4].(*string)) = "hold"
			*(dest[5].(*int64)) = f.addedAt
			*(dest[6].(*int64)) = f.removedAt
			return nil
		})
	}
	rows.EXPECT().Next().Return(false)
	rows.EXPECT().Err().Return(nil)
	return rows
}

func TestListMeasuredGroups(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	ctx := context.Background()

	t.Run("no measured group reads no target rows", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		expectMeasuredGroupsQuery(storage).Return(measuredGroupRows(controller), nil)
		// No second QueryContext is expected: gomock fails the test if the
		// target rows are read.
		groups, err := storage.ListMeasuredGroups(ctx, "env", measuredGroupsEndedSince)
		require.NoError(t, err)
		assert.Empty(t, groups)
	})

	t.Run("attaches every target row to its group, removed ones included", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		expectMeasuredGroupsQuery(storage).Return(measuredGroupRows(controller, "hg-1", "hg-2"), nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
			gomock.Any(),
			gomock.Cond(func(q any) bool {
				query := q.(string)
				return !strings.Contains(query, "removed_at = 0") &&
					strings.HasSuffix(query, " ORDER BY added_at ASC, id ASC")
			}),
			"env", "hg-1", "hg-2",
		).Return(measuredFeatureRows(
			controller,
			measuredFeatureRow{id: "r1", groupID: "hg-1", featureID: "flag-a", addedAt: 10, removedAt: 20},
			measuredFeatureRow{id: "r2", groupID: "hg-2", featureID: "flag-b", addedAt: 15},
			measuredFeatureRow{id: "r3", groupID: "hg-1", featureID: "flag-a", addedAt: 30},
		), nil)
		groups, err := storage.ListMeasuredGroups(ctx, "env", measuredGroupsEndedSince)
		require.NoError(t, err)
		require.Len(t, groups, 2)
		assert.Equal(t, "hg-1", groups[0].Id)
		require.Len(t, groups[0].Features, 2)
		assert.Equal(t, "r1", groups[0].Features[0].Id)
		assert.Equal(t, int64(20), groups[0].Features[0].RemovedAt)
		assert.Equal(t, "r3", groups[0].Features[1].Id)
		assert.Equal(t, int64(0), groups[0].Features[1].RemovedAt)
		assert.Equal(t, "hg-2", groups[1].Id)
		require.Len(t, groups[1].Features, 1)
		assert.Equal(t, "flag-b", groups[1].Features[0].FeatureId)
	})

	t.Run("group query error", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		expectMeasuredGroupsQuery(storage).Return(nil, errors.New("db"))
		groups, err := storage.ListMeasuredGroups(ctx, "env", measuredGroupsEndedSince)
		assert.Nil(t, groups)
		assert.EqualError(t, err, "db")
	})

	t.Run("target rows query error", func(t *testing.T) {
		storage := newHoldoutStorageWithMock(t, controller)
		expectMeasuredGroupsQuery(storage).Return(measuredGroupRows(controller, "hg-1"), nil)
		storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
			gomock.Any(), gomock.Any(), "env", "hg-1",
		).Return(nil, errors.New("db"))
		groups, err := storage.ListMeasuredGroups(ctx, "env", measuredGroupsEndedSince)
		assert.Nil(t, groups)
		assert.EqualError(t, err, "db")
	})
}

// ListCurrentFeatures shares its statement with ListMeasuredGroups and must
// keep leaving removed target rows out.
func TestListCurrentFeaturesLeavesRemovedRowsOut(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	storage := newHoldoutStorageWithMock(t, controller)
	storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
		gomock.Any(),
		gomock.Cond(func(q any) bool {
			return strings.HasSuffix(q.(string), " AND removed_at = 0")
		}),
		"env", "hg-1",
	).Return(measuredFeatureRows(
		controller,
		measuredFeatureRow{id: "r3", groupID: "hg-1", featureID: "flag-a", addedAt: 30},
	), nil)
	features, err := storage.ListCurrentFeatures(context.Background(), "env", []string{"hg-1"})
	require.NoError(t, err)
	require.Len(t, features, 1)
	assert.Equal(t, "r3", features[0].Id)
}

// The statement decides which groups are measured. A group without an
// internal experiment is never one, an archived group still is, and the two
// kinds of stop are dated by different columns.
func TestSelectMeasuredGroupsSQL(t *testing.T) {
	t.Parallel()
	assert.Contains(t, selectMeasuredGroupsSQL, "(experiment_id <> '' OR analysis_experiment_id <> '')")
	assert.Contains(t, selectMeasuredGroupsSQL, "stop_at >= ")
	assert.Contains(t, selectMeasuredGroupsSQL, "stopped_at >= ")
	assert.Contains(t, selectMeasuredGroupsSQL, "deleted = ")
	assert.NotContains(t, selectMeasuredGroupsSQL, "archived = ")
}
