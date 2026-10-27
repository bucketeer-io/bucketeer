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
//

package cacher

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	cachev3mock "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	envclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	holdoutmock "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2/mock"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestHoldoutGroupCacherRun(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	ctx := context.Background()

	t.Run("environment list fails", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("list"))
		err := cacher.Run(ctx)
		assert.EqualError(t, err, "list")
	})

	t.Run("missing group stores an empty snapshot", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(nil, v2hs.ErrHoldoutGroupNotFound)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetClearedHoldoutUpdatedAt(gomock.Any(), "env").
			Return(int64(0), nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				Put(&ftproto.HoldoutGroup{}, "env").
				Return(nil)
		}
		assert.NoError(t, cacher.Run(ctx))
	})

	t.Run("stopped group keeps its stop time on the empty snapshot", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(nil, v2hs.ErrHoldoutGroupNotFound)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetClearedHoldoutUpdatedAt(gomock.Any(), "env").
			Return(int64(1700000000), nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				Put(&ftproto.HoldoutGroup{UpdatedAt: 1700000000}, "env").
				Return(nil)
		}
		assert.NoError(t, cacher.Run(ctx))
	})

	t.Run("stores the live group", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		group := &ftproto.HoldoutGroup{Id: "hg", EnvironmentId: "env"}
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(&domain.HoldoutGroup{HoldoutGroup: group}, nil)
		for _, cache := range cacher.caches {
			cache.(*cachev3mock.MockHoldoutGroupCache).EXPECT().
				Put(group, "env").
				Return(nil)
		}
		assert.NoError(t, cacher.Run(ctx))
	})

	t.Run("storage error stops the job", func(t *testing.T) {
		cacher := newHoldoutGroupCacher(t, controller)
		cacher.environmentClient.(*envclientmock.MockClient).EXPECT().
			ListEnvironmentsV2(gomock.Any(), gomock.Any()).
			Return(&envproto.ListEnvironmentsV2Response{
				Environments: []*envproto.EnvironmentV2{{Id: "env"}},
			}, nil)
		cacher.storage.(*holdoutmock.MockHoldoutStorage).EXPECT().
			GetLiveServingGroup(gomock.Any(), "env").
			Return(nil, errors.New("db"))
		err := cacher.Run(ctx)
		assert.EqualError(t, err, "db")
	})
}

func newHoldoutGroupCacher(t *testing.T, controller *gomock.Controller) *holdoutGroupCacher {
	t.Helper()
	require.NotNil(t, controller)
	return &holdoutGroupCacher{
		environmentClient: envclientmock.NewMockClient(controller),
		storage:           holdoutmock.NewMockHoldoutStorage(controller),
		caches: []cachev3.HoldoutGroupCache{
			cachev3mock.NewMockHoldoutGroupCache(controller),
			cachev3mock.NewMockHoldoutGroupCache(controller),
		},
		logger: zap.NewNop(),
	}
}
