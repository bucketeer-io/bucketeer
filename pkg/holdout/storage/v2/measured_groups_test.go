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

package v2_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2/mock"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestMeasuredGroups(t *testing.T) {
	t.Parallel()
	now := time.Unix(1_700_000_000, 0)
	// A stopped group stays in the list for as long as its internal
	// experiments stay in the list of experiments the pipeline reads.
	endedSince := now.Add(-experimentdomain.StoppedExperimentEventWindow).Unix()

	t.Run("keeps the order and the target rows of the storage", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		storage := mock.NewMockHoldoutStorage(controller)
		stopped := &proto.HoldoutGroup{
			Id:           "stopped",
			ExperimentId: "e1",
			Features: []*proto.HoldoutGroupFeature{
				{Id: "r1", FeatureId: "flag", AddedAt: 10, RemovedAt: 20},
			},
		}
		running := &proto.HoldoutGroup{Id: "running", ExperimentId: "e2"}
		storage.EXPECT().ListMeasuredGroups(gomock.Any(), "env", endedSince).Return(
			[]*domain.HoldoutGroup{{HoldoutGroup: stopped}, {HoldoutGroup: running}}, nil,
		)
		list, err := v2hs.MeasuredGroups(context.Background(), storage, "env", now)
		require.NoError(t, err)
		require.Len(t, list.HoldoutGroups, 2)
		assert.Same(t, stopped, list.HoldoutGroups[0])
		assert.Same(t, running, list.HoldoutGroups[1])
	})

	t.Run("no group is an empty list, not nil", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		storage := mock.NewMockHoldoutStorage(controller)
		storage.EXPECT().ListMeasuredGroups(gomock.Any(), "env", endedSince).Return(nil, nil)
		list, err := v2hs.MeasuredGroups(context.Background(), storage, "env", now)
		require.NoError(t, err)
		require.NotNil(t, list)
		assert.Empty(t, list.HoldoutGroups)
	})

	t.Run("storage error", func(t *testing.T) {
		t.Parallel()
		controller := gomock.NewController(t)
		storage := mock.NewMockHoldoutStorage(controller)
		storageErr := errors.New("db")
		storage.EXPECT().ListMeasuredGroups(gomock.Any(), "env", endedSince).Return(nil, storageErr)
		list, err := v2hs.MeasuredGroups(context.Background(), storage, "env", now)
		assert.Nil(t, list)
		assert.ErrorIs(t, err, storageErr)
	})
}
