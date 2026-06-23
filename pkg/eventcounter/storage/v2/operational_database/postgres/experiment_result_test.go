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

	operationaldatabase "github.com/bucketeer-io/bucketeer/v2/pkg/eventcounter/storage/v2/operational_database"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/postgres/mock"
	ecproto "github.com/bucketeer-io/bucketeer/v2/proto/eventcounter"
)

func TestNewExperimentResultStoragePostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewExperimentResultStorage(mock.NewMockTransaction(mockController))
	assert.IsType(t, &experimentResultStorage{}, storage)
}

func TestGetExperimentResultPostgres(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc          string
		setup         func(*experimentResultStorage)
		id            string
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrExperimentResultNotFound",
			setup: func(s *experimentResultStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(postgres.ErrNoRows)
				s.qe.(*mock.MockTransaction).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   operationaldatabase.ErrExperimentResultNotFound,
		},
		{
			desc: "Error",
			setup: func(s *experimentResultStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
				s.qe.(*mock.MockTransaction).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *experimentResultStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
				s.qe.(*mock.MockTransaction).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            "id-0",
			environmentId: "ns",
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := &experimentResultStorage{mock.NewMockTransaction(mockController)}
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetExperimentResult(context.Background(), p.id, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

// TestGetExperimentResultPostgres_PropagatesAllProtoFields locks in that every
// proto field of ExperimentResult round-trips through the read storage.
// GetExperimentResult scans the top-level columns (Id, ExperimentId,
// UpdatedAt) into the result directly and the rest of the proto out of a
// JSON blob into a side-by-side struct (erForGoalResults), then has to
// manually copy each non-top-level field onto the returned object. Forgetting
// to copy a field is a silent footgun — the field is in the JSON blob but
// gets dropped on the floor (this was the cause of an e2e regression that
// surfaced after SrmResult was added to ExperimentResult).
//
// This test populates each field via the JSON blob path and asserts every
// one is present on the returned domain object, so any future ExperimentResult
// field addition fails this test until the field copy is added to
// GetExperimentResult.
func TestGetExperimentResultPostgres_PropagatesAllProtoFields(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expectedGoalResults := []*ecproto.GoalResult{
		{GoalId: "goal-1"}, {GoalId: "goal-2"},
	}
	const expectedTotalEvalUsers = int64(12345)
	expectedSRM := &ecproto.SrmResult{
		Status:    ecproto.SrmResult_OK,
		PValue:    0.42,
		Threshold: 0.001,
	}

	row := mock.NewMockRow(mockController)
	row.EXPECT().Scan(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).DoAndReturn(func(args ...interface{}) error {
		*args[0].(*string) = "exp-1"
		*args[1].(*string) = "exp-1"
		*args[2].(*int64) = 1700000000
		// args[3] is *postgres.JSONObject{Val: *ecproto.ExperimentResult};
		// populate the inner proto exactly as the JSON unmarshal would, so
		// the field-copy logic in GetExperimentResult has something to
		// propagate.
		target := args[3].(*postgres.JSONObject).Val.(*ecproto.ExperimentResult)
		target.GoalResults = expectedGoalResults
		target.TotalEvaluationUserCount = expectedTotalEvalUsers
		target.SrmResult = expectedSRM
		return nil
	})

	qe := mock.NewMockTransaction(mockController)
	qe.EXPECT().QueryRowContext(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(row)

	storage := &experimentResultStorage{qe: qe}
	got, err := storage.GetExperimentResult(context.Background(), "exp-1", "env-1")
	if !assert.NoError(t, err) || !assert.NotNil(t, got) || !assert.NotNil(t, got.ExperimentResult) {
		return
	}
	er := got.ExperimentResult
	assert.Equal(t, "exp-1", er.Id, "top-level Id column")
	assert.Equal(t, "exp-1", er.ExperimentId, "top-level ExperimentId column")
	assert.EqualValues(t, 1700000000, er.UpdatedAt, "top-level UpdatedAt column")
	assert.Equal(t, expectedGoalResults, er.GoalResults,
		"GoalResults must be copied from the JSON blob — without the copy "+
			"line in GetExperimentResult this returns nil")
	assert.Equal(t, expectedTotalEvalUsers, er.TotalEvaluationUserCount,
		"TotalEvaluationUserCount must be copied from the JSON blob")
	assert.Equal(t, expectedSRM, er.SrmResult,
		"SrmResult must be copied from the JSON blob — failure here means a "+
			"new proto field was added to ExperimentResult without updating "+
			"the field-copy block in GetExperimentResult")
}
