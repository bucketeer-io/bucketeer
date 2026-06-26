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

	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	dwhstorage "github.com/bucketeer-io/bucketeer/v2/pkg/subscriber/storage/dwhstorage"
)

func TestNewEvaluationEventStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewEvaluationEventStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &evaluationEventStorage{}, storage)
}

func TestEvaluationEventStorage_CreateEvaluationEvents(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	validEvent := dwhstorage.EvaluationEventParams{
		ID:            "id-1",
		EnvironmentID: "env-1",
		FeatureID:     "feature-1",
		UserID:        "user-1",
		Timestamp:     1000000000,
	}

	patterns := []struct {
		desc        string
		setup       func(s *evaluationEventStorage)
		input       []dwhstorage.EvaluationEventParams
		expectedErr bool
	}{
		{
			desc:        "success: empty",
			setup:       func(s *evaluationEventStorage) {},
			input:       []dwhstorage.EvaluationEventParams{},
			expectedErr: false,
		},
		{
			desc: "success",
			setup: func(s *evaluationEventStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input:       []dwhstorage.EvaluationEventParams{validEvent},
			expectedErr: false,
		},
		{
			desc:        "error: missing required field",
			setup:       func(s *evaluationEventStorage) {},
			input:       []dwhstorage.EvaluationEventParams{{ID: "", EnvironmentID: "env"}},
			expectedErr: true,
		},
		{
			desc: "error: exec fails",
			setup: func(s *evaluationEventStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input:       []dwhstorage.EvaluationEventParams{validEvent},
			expectedErr: true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &evaluationEventStorage{qe: mock.NewMockQueryExecer(mockController)}
			p.setup(s)
			err := s.CreateEvaluationEvents(context.Background(), p.input)
			if p.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
