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

package cacher

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	mockcachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3/mock"
	ftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	mockftstorage "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestSegmentUserCacher_RefreshAllEnvironmentCaches(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	internalErr := errors.New("internal error")

	patterns := []struct {
		desc        string
		setup       func(*segmentUserCacher)
		expectedErr error
	}{
		{
			desc: "err: failed to list all in-use segments",
			setup: func(c *segmentUserCacher) {
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListAllInUseSegments(gomock.Any()).
					Return(nil, internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc: "success: empty segments",
			setup: func(c *segmentUserCacher) {
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListAllInUseSegments(gomock.Any()).
					Return([]*ftstorage.InUseSegment{}, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: refresh cache for all segments",
			setup: func(c *segmentUserCacher) {
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListAllInUseSegments(gomock.Any()).
					Return([]*ftstorage.InUseSegment{
						{SegmentID: "seg-id-1", EnvironmentID: "env-id-1", UpdatedAt: 1234567890},
						{SegmentID: "seg-id-2", EnvironmentID: "env-id-2", UpdatedAt: 1234567891},
					}, nil)
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListSegmentUsersBySegment(gomock.Any(), "seg-id-1", "env-id-1").
					Return([]*ftproto.SegmentUser{
						{Id: "user-id-1", SegmentId: "seg-id-1", UserId: "user-1"},
					}, nil)
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListSegmentUsersBySegment(gomock.Any(), "seg-id-2", "env-id-2").
					Return([]*ftproto.SegmentUser{
						{Id: "user-id-2", SegmentId: "seg-id-2", UserId: "user-2"},
						{Id: "user-id-3", SegmentId: "seg-id-2", UserId: "user-3"},
					}, nil)
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(gomock.Any(), "env-id-1").
					Return(nil)
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(gomock.Any(), "env-id-2").
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: segment with no users (still cached)",
			setup: func(c *segmentUserCacher) {
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListAllInUseSegments(gomock.Any()).
					Return([]*ftstorage.InUseSegment{
						{SegmentID: "seg-id-1", EnvironmentID: "env-id-1", UpdatedAt: 1234567890},
					}, nil)
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListSegmentUsersBySegment(gomock.Any(), "seg-id-1", "env-id-1").
					Return([]*ftproto.SegmentUser{}, nil)
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(gomock.Any(), "env-id-1").
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: continues on segment user fetch error",
			setup: func(c *segmentUserCacher) {
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListAllInUseSegments(gomock.Any()).
					Return([]*ftstorage.InUseSegment{
						{SegmentID: "seg-id-1", EnvironmentID: "env-id-1", UpdatedAt: 1234567890},
						{SegmentID: "seg-id-2", EnvironmentID: "env-id-2", UpdatedAt: 1234567891},
					}, nil)
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListSegmentUsersBySegment(gomock.Any(), "seg-id-1", "env-id-1").
					Return(nil, internalErr) // First segment fails
				c.segStorage.(*mockftstorage.MockSegmentStorage).EXPECT().
					ListSegmentUsersBySegment(gomock.Any(), "seg-id-2", "env-id-2").
					Return([]*ftproto.SegmentUser{
						{Id: "user-id-2", SegmentId: "seg-id-2", UserId: "user-2"},
					}, nil) // Second segment succeeds
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(gomock.Any(), "env-id-2").
					Return(nil)
			},
			expectedErr: nil, // Should not fail, continues with other segments
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newSegmentUserCacherWithMock(t, controller, 1)
			p.setup(cacher)
			err := cacher.RefreshAllEnvironmentCaches(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestSegmentUserCacher_PutCache(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id"
	segmentUsers := &ftproto.SegmentUsers{
		SegmentId: "seg-id-1",
		Users: []*ftproto.SegmentUser{
			{Id: "user-id-1"},
			{Id: "user-id-2"},
		},
	}

	patterns := []struct {
		desc  string
		setup func(*segmentUserCacher)
	}{
		{
			desc: "success: put to single cache",
			setup: func(c *segmentUserCacher) {
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(segmentUsers, envID).
					Return(nil)
			},
		},
		{
			desc: "err: cache put fails (logged but not returned)",
			setup: func(c *segmentUserCacher) {
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(segmentUsers, envID).
					Return(errors.New("cache error"))
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newSegmentUserCacherWithMock(t, controller, 1)
			p.setup(cacher)
			// putCache doesn't return error, it just logs and records metrics
			cacher.putCache(segmentUsers, envID, len(segmentUsers.Users))
		})
	}
}

func TestSegmentUserCacher_PutCacheMultipleInstances(t *testing.T) {
	t.Parallel()
	controller := gomock.NewController(t)
	defer controller.Finish()

	envID := "env-id"
	segmentUsers := &ftproto.SegmentUsers{
		SegmentId: "seg-id-1",
		Users:     []*ftproto.SegmentUser{{Id: "user-id-1"}},
	}

	patterns := []struct {
		desc  string
		setup func(*segmentUserCacher)
	}{
		{
			desc: "success: put to multiple caches",
			setup: func(c *segmentUserCacher) {
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(segmentUsers, envID).
					Return(nil)
				c.caches[1].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(segmentUsers, envID).
					Return(nil)
			},
		},
		{
			desc: "partial failure: one cache fails",
			setup: func(c *segmentUserCacher) {
				c.caches[0].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(segmentUsers, envID).
					Return(errors.New("cache error"))
				c.caches[1].(*mockcachev3.MockSegmentUsersCache).EXPECT().
					Put(segmentUsers, envID).
					Return(nil)
			},
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			cacher := newSegmentUserCacherWithMock(t, controller, 2)
			p.setup(cacher)
			cacher.putCache(segmentUsers, envID, len(segmentUsers.Users))
		})
	}
}

func newSegmentUserCacherWithMock(t *testing.T, controller *gomock.Controller, numCaches int) *segmentUserCacher {
	t.Helper()
	logger, err := log.NewLogger()
	require.NoError(t, err)

	caches := make([]cachev3.SegmentUsersCache, numCaches)
	for i := 0; i < numCaches; i++ {
		caches[i] = mockcachev3.NewMockSegmentUsersCache(controller)
	}

	return &segmentUserCacher{
		segStorage: mockftstorage.NewMockSegmentStorage(controller),
		caches:     caches,
		logger:     logger,
	}
}
