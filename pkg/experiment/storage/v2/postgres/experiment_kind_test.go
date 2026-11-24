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
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	v2es "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

// experimentColumnPattern matches a column of the experiment table in a select list.
var experimentColumnPattern = regexp.MustCompile(`(?m)^\s+ex\.(\w+),$`)

// The position of kind in the select list and in Scan must agree,
// because Scan reads the columns by position.
func TestExperimentSelectColumnOrder(t *testing.T) {
	t.Parallel()
	for name, query := range map[string]string{
		"select_experiment.sql":  selectExperimentSQL,
		"select_experiments.sql": selectExperimentsSQL,
	} {
		matches := experimentColumnPattern.FindAllStringSubmatch(query, -1)
		columns := make([]string, 0, len(matches))
		for _, m := range matches {
			columns = append(columns, m[1])
		}
		require.Len(t, columns, 20, name)
		assert.Equal(t, []string{"maintainer", "status", "kind"}, columns[17:], name)
	}
}

func TestCreateExperimentStoresKind(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	columns := strings.Fields(strings.NewReplacer(",", " ", "(", " ", ")", " ").Replace(
		insertExperimentSQL[:strings.Index(insertExperimentSQL, "VALUES")],
	))
	// The first three fields are "INSERT INTO experiment".
	columns = columns[3:]
	require.Equal(t, "kind", columns[len(columns)-1])

	for _, kind := range []proto.Experiment_Kind{
		proto.Experiment_KIND_EXPERIMENT,
		proto.Experiment_HOLDOUT,
	} {
		storage := newExperimentStorageWithMock(t, mockController)
		storage.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
			gomock.Any(), insertExperimentSQL, gomock.Any(),
		).DoAndReturn(func(_ context.Context, _ string, args ...interface{}) (postgres.Result, error) {
			// One argument for each column, and kind is the last one.
			require.Len(t, args, len(columns))
			assert.Equal(t, int32(kind), args[len(args)-1])
			return nil, nil
		})
		err := storage.CreateExperiment(
			context.Background(),
			&domain.Experiment{Experiment: &proto.Experiment{Id: "id-0", Kind: kind}},
			"ns0",
		)
		assert.NoError(t, err)
	}
}

func TestGetExperimentReadsKind(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	storage := newExperimentStorageWithMock(t, mockController)
	row := mock.NewMockRow(mockController)
	row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
		// The destinations follow the select list: status is the 19th and kind the 20th.
		require.Len(t, dest, 21)
		*dest[18].(*int32) = int32(proto.Experiment_STOPPED)
		*dest[19].(*int32) = int32(proto.Experiment_HOLDOUT)
		return nil
	})
	storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
		gomock.Any(), selectExperimentSQL, gomock.Any(),
	).Return(row)

	actual, err := storage.GetExperiment(context.Background(), "id-0", "ns0")
	require.NoError(t, err)
	assert.Equal(t, proto.Experiment_STOPPED, actual.Status)
	assert.Equal(t, proto.Experiment_HOLDOUT, actual.Kind)
}

func TestListExperimentsReadsKind(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	storage := newExperimentStorageWithMock(t, mockController)
	rows := mock.NewMockRows(mockController)
	rows.EXPECT().Close().Return(nil)
	gomock.InOrder(
		rows.EXPECT().Next().Return(true),
		rows.EXPECT().Next().Return(false),
	)
	rows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
		require.Len(t, dest, 21)
		*dest[18].(*int32) = int32(proto.Experiment_RUNNING)
		*dest[19].(*int32) = int32(proto.Experiment_HOLDOUT)
		return nil
	})
	rows.EXPECT().Err().Return(nil)
	storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(rows, nil)
	countRow := mock.NewMockRow(mockController)
	countRow.EXPECT().Scan(gomock.Any()).Return(nil)
	storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(countRow)

	actual, _, _, err := storage.ListExperiments(
		context.Background(),
		v2es.ListExperimentsParams{EnvironmentID: "ns0"},
	)
	require.NoError(t, err)
	require.Len(t, actual, 1)
	assert.Equal(t, proto.Experiment_RUNNING, actual[0].Status)
	assert.Equal(t, proto.Experiment_HOLDOUT, actual[0].Kind)
}

func TestListExperimentsKindsFilter(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		kinds         []proto.Experiment_Kind
		expectedWhere string
		expectedArgs  []interface{}
	}{
		{
			desc:          "no kinds: every kind is listed",
			kinds:         nil,
			expectedWhere: " WHERE deleted = $1 AND environment_id = $2 ",
			expectedArgs:  []interface{}{false, "ns0"},
		},
		{
			desc:          "one kind",
			kinds:         []proto.Experiment_Kind{proto.Experiment_KIND_EXPERIMENT},
			expectedWhere: " WHERE deleted = $1 AND environment_id = $2 AND  kind IN ($3) ",
			expectedArgs:  []interface{}{false, "ns0", int32(0)},
		},
		{
			desc: "two kinds",
			kinds: []proto.Experiment_Kind{
				proto.Experiment_KIND_EXPERIMENT,
				proto.Experiment_HOLDOUT,
			},
			expectedWhere: " WHERE deleted = $1 AND environment_id = $2 AND  kind IN ($3, $4) ",
			expectedArgs:  []interface{}{false, "ns0", int32(0), int32(1)},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newExperimentStorageWithMock(t, mockController)
			rows := mock.NewMockRows(mockController)
			rows.EXPECT().Close().Return(nil)
			rows.EXPECT().Next().Return(false)
			rows.EXPECT().Err().Return(nil)
			// The list and the total count use the same filter.
			storage.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
				gomock.Any(), gomock.Any(), gomock.Any(),
			).DoAndReturn(func(_ context.Context, query string, args ...interface{}) (postgres.Rows, error) {
				assert.True(t, strings.HasPrefix(query, selectExperimentsSQL+p.expectedWhere), query)
				assert.Equal(t, p.expectedArgs, args)
				return rows, nil
			})
			countRow := mock.NewMockRow(mockController)
			countRow.EXPECT().Scan(gomock.Any()).Return(nil)
			storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
				gomock.Any(), gomock.Any(), gomock.Any(),
			).DoAndReturn(func(_ context.Context, query string, args ...interface{}) postgres.Row {
				assert.Equal(t, countExperimentSQL+p.expectedWhere, query)
				assert.Equal(t, p.expectedArgs, args)
				return countRow
			})
			_, _, _, err := storage.ListExperiments(
				context.Background(),
				v2es.ListExperimentsParams{EnvironmentID: "ns0", Kinds: p.kinds},
			)
			assert.NoError(t, err)
		})
	}
}

func TestGetExperimentSummaryKindsFilter(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc          string
		kinds         []proto.Experiment_Kind
		expectedWhere string
		expectedArgs  []interface{}
	}{
		{
			desc:          "no kinds: every kind is counted",
			kinds:         nil,
			expectedWhere: " WHERE environment_id = $1 AND deleted = $2 AND archived = $3 ",
			expectedArgs:  []interface{}{"ns0", false, false},
		},
		{
			desc:          "one kind",
			kinds:         []proto.Experiment_Kind{proto.Experiment_KIND_EXPERIMENT},
			expectedWhere: " WHERE environment_id = $1 AND deleted = $2 AND archived = $3 AND  kind IN ($4) ",
			expectedArgs:  []interface{}{"ns0", false, false, int32(0)},
		},
		{
			desc: "two kinds",
			kinds: []proto.Experiment_Kind{
				proto.Experiment_KIND_EXPERIMENT,
				proto.Experiment_HOLDOUT,
			},
			expectedWhere: " WHERE environment_id = $1 AND deleted = $2 AND archived = $3 AND  kind IN ($4, $5) ",
			expectedArgs:  []interface{}{"ns0", false, false, int32(0), int32(1)},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newExperimentStorageWithMock(t, mockController)
			row := mock.NewMockRow(mockController)
			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
				*dest[0].(*int64) = 1
				*dest[1].(*int64) = 2
				*dest[2].(*int64) = 3
				return nil
			})
			storage.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
				gomock.Any(), gomock.Any(), gomock.Any(),
			).DoAndReturn(func(_ context.Context, query string, args ...interface{}) postgres.Row {
				assert.Equal(t, summarizeExperimentSQL+p.expectedWhere, query)
				assert.Equal(t, p.expectedArgs, args)
				return row
			})
			summary, err := storage.GetExperimentSummary(context.Background(), "ns0", p.kinds)
			require.NoError(t, err)
			assert.Equal(t, &v2es.ExperimentSummary{
				TotalWaitingCount: 1,
				TotalRunningCount: 2,
				TotalStoppedCount: 3,
			}, summary)
		})
	}
}

// A goal lists the experiments on feature flags that use it, and not an
// experiment whose kind is HOLDOUT. is_in_use_status counts every kind:
// a goal that only a holdout group measures is still in use.
func TestGoalQueriesAndExperimentKind(t *testing.T) {
	t.Parallel()
	for name, query := range map[string]string{
		"select_goal.sql":  selectGoalSQL,
		"select_goals.sql": selectGoalsSQL,
	} {
		inUseStart := strings.Index(query, "FROM experiment ex1")
		inUseEnd := strings.Index(query, "AS is_in_use_status")
		listStart := strings.Index(query, "FROM experiment ex2")
		listEnd := strings.Index(query, "AS experiments")
		require.True(t, 0 <= inUseStart && inUseStart < inUseEnd && inUseEnd < listStart && listStart < listEnd, name)
		assert.NotContains(t, query[inUseStart:inUseEnd], "kind", name)
		assert.Contains(t, query[listStart:listEnd], "ex2.kind = 0 AND", name)
	}
	assert.NotContains(t, countGoalSQL, "kind")
}
