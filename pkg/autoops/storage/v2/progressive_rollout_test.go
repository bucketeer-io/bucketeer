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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestNewProgressiveRolloutStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewProgressiveRolloutStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &progressiveRolloutStorage{}, db)
}

func TestCreateProgressiveRollout(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc                 string
		setup                func(*progressiveRolloutStorage)
		input                *domain.ProgressiveRollout
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "",
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.ProgressiveRollout{
				ProgressiveRollout: &proto.ProgressiveRollout{Id: "id-1"},
			},
			environmentNamespace: "ns0",
			expectedErr:          ErrProgressiveRolloutAlreadyExists,
		},
		{
			setup: func(s *progressiveRolloutStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.ProgressiveRollout{
				ProgressiveRollout: &proto.ProgressiveRollout{Id: "id-1"},
			},
			environmentNamespace: "ns0",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		storage := newProgressiveRolloutStorageWithMock(t, mockController)
		if p.setup != nil {
			p.setup(storage)
		}
		err := storage.CreateProgressiveRollout(context.Background(), p.input, p.environmentNamespace)
		assert.Equal(t, p.expectedErr, err)
	}
}

func newProgressiveRolloutStorageWithMock(t *testing.T, mockController *gomock.Controller) *progressiveRolloutStorage {
	t.Helper()
	return &progressiveRolloutStorage{mock.NewMockQueryExecer(mockController)}
}
