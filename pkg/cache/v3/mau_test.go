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

package v3

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
)

func TestMAUCache_MAUKey(t *testing.T) {
	t.Parallel()
	month := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	actual := mauKey("env-123", "ANDROID", month)
	assert.Equal(t, "env-123:mau:ANDROID:202601", actual)
}

func TestMAUCache_MergeIntoMAUBatch_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceIDs := []string{"ANDROID", "IOS"}
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	gomock.InOrder(
		mockCache.EXPECT().PFMerge("env-123:mau:ANDROID:202601", mauTTL, "env-123:mau:ANDROID:202601", "env-123:dau:ANDROID:20260115").Return(nil),
		mockCache.EXPECT().PFCount("env-123:mau:ANDROID:202601").Return(int64(100), nil),
		mockCache.EXPECT().Delete("env-123:dau:ANDROID:20260115").Return(nil),
		mockCache.EXPECT().PFMerge("env-123:mau:IOS:202601", mauTTL, "env-123:mau:IOS:202601", "env-123:dau:IOS:20260115").Return(nil),
		mockCache.EXPECT().PFCount("env-123:mau:IOS:202601").Return(int64(200), nil),
		mockCache.EXPECT().Delete("env-123:dau:IOS:20260115").Return(nil),
	)

	result, err := c.MergeIntoMAUBatch(envID, sourceIDs, date)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int64{"ANDROID": 100, "IOS": 200}, result)
}

func TestMAUCache_MergeIntoMAUBatch_EmptySourceIDs(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	c := NewMAUCache(mockCache)

	result, err := c.MergeIntoMAUBatch("env-123", []string{}, time.Now())
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestMAUCache_MergeIntoMAUBatch_PFMergeError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceIDs := []string{"ANDROID"}
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedErr := errors.New("merge error")

	mockCache.EXPECT().PFMerge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr)

	result, err := c.MergeIntoMAUBatch(envID, sourceIDs, date)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to merge MAU for source ANDROID")
	assert.Nil(t, result)
}

func TestMAUCache_MergeIntoMAUBatch_PFCountError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceIDs := []string{"ANDROID"}
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedErr := errors.New("count error")

	gomock.InOrder(
		mockCache.EXPECT().PFMerge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockCache.EXPECT().PFCount(gomock.Any()).Return(int64(0), expectedErr),
	)

	result, err := c.MergeIntoMAUBatch(envID, sourceIDs, date)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count MAU for source ANDROID")
	assert.Nil(t, result)
}

func TestMAUCache_MergeIntoMAUBatch_DeleteError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceIDs := []string{"ANDROID"}
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedErr := errors.New("delete error")

	gomock.InOrder(
		mockCache.EXPECT().PFMerge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockCache.EXPECT().PFCount(gomock.Any()).Return(int64(100), nil),
		mockCache.EXPECT().Delete(gomock.Any()).Return(expectedErr),
	)

	result, err := c.MergeIntoMAUBatch(envID, sourceIDs, date)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete DAU key for source ANDROID")
	assert.Nil(t, result)
}
