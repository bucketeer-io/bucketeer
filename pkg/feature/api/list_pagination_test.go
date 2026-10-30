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

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	tagstorage "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage"
	tagstoragemock "github.com/bucketeer-io/bucketeer/v2/pkg/tag/storage/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
	tagproto "github.com/bucketeer-io/bucketeer/v2/proto/tag"
)

// listPaginationCase describes a page_size / cursor pair of a list request
// and what the handler is expected to do with it: reject it without reading
// the storage, or pass the page size and cursor through to the storage.
type listPaginationCase struct {
	name     string
	pageSize int64
	cursor   string
	wantErr  error
}

func listPaginationCases() []listPaginationCase {
	return []listPaginationCase{
		{
			name:     "negative page size is rejected",
			pageSize: -1,
			cursor:   "0",
			wantErr:  statusInvalidPageSize.Err(),
		},
		{
			name:     "negative cursor is rejected",
			pageSize: 10,
			cursor:   "-5",
			wantErr:  statusInvalidCursor.Err(),
		},
		{
			name:     "non-numeric cursor is rejected",
			pageSize: 10,
			cursor:   "XXX",
			wantErr:  statusInvalidCursor.Err(),
		},
		{
			name:     "page size and cursor are passed to the storage",
			pageSize: 10,
			cursor:   "20",
		},
		{
			name:     "zero page size and empty cursor list everything from the start",
			pageSize: 0,
			cursor:   "",
		},
	}
}

// assertListPaginationResult checks the outcome of a list call: a rejected
// request returns the expected error and never reaches the storage; an
// accepted one reaches the storage with the page size of the request and
// with wantCursor.
func assertListPaginationResult(
	t *testing.T,
	tc listPaginationCase,
	err error,
	listed bool,
	listedPageSize int64,
	listedCursor string,
	wantCursor string,
) {
	t.Helper()
	if tc.wantErr != nil {
		require.Error(t, err)
		assert.Equal(t, tc.wantErr.Error(), err.Error())
		assert.False(t, listed)
		return
	}
	require.NoError(t, err)
	require.True(t, listed)
	assert.Equal(t, tc.pageSize, listedPageSize)
	assert.Equal(t, wantCursor, listedCursor)
}

// defaultedCursor is the cursor the handlers that parse the cursor themselves
// hand to the storage: an empty cursor becomes the first page.
func defaultedCursor(cursor string) string {
	if cursor == "" {
		return "0"
	}
	return cursor
}

func TestListFeaturesPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createFeatureService(mockController)
			var listed *v2fs.ListFeaturesParams
			if tc.wantErr == nil {
				service.featureStorage.(*mock.MockFeatureStorage).EXPECT().
					ListFeatures(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2fs.ListFeaturesParams) ([]*featureproto.Feature, int, int64, error) {
						listed = &p
						return []*featureproto.Feature{}, 0, 0, nil
					})
				service.featureStorage.(*mock.MockFeatureStorage).EXPECT().
					GetFeatureSummary(gomock.Any(), "ns0").
					Return(&featureproto.FeatureSummary{}, nil)
			}

			resp, err := service.ListFeatures(createContextWithToken(), &featureproto.ListFeaturesRequest{
				EnvironmentId: "ns0",
				PageSize:      tc.pageSize,
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor, tc.cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListEnabledFeaturesPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createFeatureService(mockController)
			var listed *v2fs.ListFeaturesParams
			if tc.wantErr == nil {
				service.featureStorage.(*mock.MockFeatureStorage).EXPECT().
					ListFeatures(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2fs.ListFeaturesParams) ([]*featureproto.Feature, int, int64, error) {
						listed = &p
						return []*featureproto.Feature{}, 0, 0, nil
					})
			}

			resp, err := service.ListEnabledFeatures(createContextWithToken(), &featureproto.ListEnabledFeaturesRequest{
				EnvironmentId: "ns0",
				PageSize:      tc.pageSize,
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor, tc.cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListSegmentsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createFeatureService(mockController)
			var listed *v2fs.ListSegmentsParams
			if tc.wantErr == nil {
				service.segmentStorage.(*mock.MockSegmentStorage).EXPECT().
					ListSegments(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2fs.ListSegmentsParams) ([]*featureproto.Segment, int, int64, map[string][]string, error) {
						listed = &p
						return []*featureproto.Segment{}, 0, 0, map[string][]string{}, nil
					})
				// The segments are completed with the features that use them.
				service.featureStorage.(*mock.MockFeatureStorage).EXPECT().
					ListFeatures(gomock.Any(), gomock.Any()).
					Return([]*featureproto.Feature{}, 0, int64(0), nil)
			}

			resp, err := service.ListSegments(createContextWithToken(), &featureproto.ListSegmentsRequest{
				EnvironmentId: "ns0",
				PageSize:      tc.pageSize,
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = listed.PageSize, listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor, tc.cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListSegmentUsersPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createFeatureService(mockController)
			var listed *v2fs.ListSegmentUsersParams
			if tc.wantErr == nil {
				service.segmentUserStorage.(*mock.MockSegmentUserStorage).EXPECT().
					ListSegmentUsers(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2fs.ListSegmentUsersParams) ([]*featureproto.SegmentUser, int, error) {
						listed = &p
						return []*featureproto.SegmentUser{}, 0, nil
					})
			}

			resp, err := service.ListSegmentUsers(createContextWithToken(), &featureproto.ListSegmentUsersRequest{
				EnvironmentId: "ns0",
				SegmentId:     "segment-id",
				PageSize:      tc.pageSize,
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = int64(listed.PageSize), listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor, tc.cursor)
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListFlagTriggersPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createFeatureService(mockController)
			var listed *v2fs.ListFlagTriggersParams
			if tc.wantErr == nil {
				service.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().
					ListFlagTriggers(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p v2fs.ListFlagTriggersParams) ([]*featureproto.FlagTrigger, int, int64, error) {
						listed = &p
						return []*featureproto.FlagTrigger{}, 0, 0, nil
					})
			}

			resp, err := service.ListFlagTriggers(createContextWithToken(), &featureproto.ListFlagTriggersRequest{
				EnvironmentId: "ns0",
				FeatureId:     "feature-id",
				PageSize:      int32(tc.pageSize),
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = int64(listed.PageSize), listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor, defaultedCursor(tc.cursor))
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}

func TestListTagsPaginationValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range listPaginationCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createFeatureService(mockController)
			var listed *tagstorage.ListTagsParams
			if tc.wantErr == nil {
				service.tagStorage.(*tagstoragemock.MockTagStorage).EXPECT().
					ListTags(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, p tagstorage.ListTagsParams) ([]*tagproto.Tag, int, int64, error) {
						listed = &p
						return []*tagproto.Tag{}, 0, 0, nil
					})
			}

			resp, err := service.ListTags(createContextWithToken(), &featureproto.ListTagsRequest{
				EnvironmentId: "ns0",
				PageSize:      tc.pageSize,
				Cursor:        tc.cursor,
			})
			var pageSize int64
			var cursor string
			if listed != nil {
				pageSize, cursor = int64(listed.PageSize), listed.Cursor
			}
			assertListPaginationResult(t, tc, err, listed != nil, pageSize, cursor, defaultedCursor(tc.cursor))
			assert.Equal(t, tc.wantErr == nil, resp != nil)
		})
	}
}
