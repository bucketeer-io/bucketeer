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
	"context"
	"errors"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache/mock"
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
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
	mockPipe := redismock.NewMockPipeClient(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceIDs := []string{"ANDROID", "IOS"}
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	mockCache.EXPECT().Pipeline(false).Return(mockPipe)

	// PFMERGE and Expire for each sourceID
	mockPipe.EXPECT().PFMerge("env-123:mau:ANDROID:202601", "env-123:mau:ANDROID:202601", "{env-123:ANDROID:au}:d:20260115")
	mockPipe.EXPECT().Expire("env-123:mau:ANDROID:202601", mauTTL)
	mockPipe.EXPECT().Del("{env-123:ANDROID:au}:d:20260115")
	mockPipe.EXPECT().PFMerge("env-123:mau:IOS:202601", "env-123:mau:IOS:202601", "{env-123:IOS:au}:d:20260115")
	mockPipe.EXPECT().Expire("env-123:mau:IOS:202601", mauTTL)
	mockPipe.EXPECT().Del("{env-123:IOS:au}:d:20260115")

	// PFCount for each sourceID - use real goredis.IntCmd with SetVal
	androidCmd := goredis.NewIntCmd(context.Background())
	androidCmd.SetVal(100)
	mockPipe.EXPECT().PFCount("env-123:mau:ANDROID:202601").Return(androidCmd)

	iosCmd := goredis.NewIntCmd(context.Background())
	iosCmd.SetVal(200)
	mockPipe.EXPECT().PFCount("env-123:mau:IOS:202601").Return(iosCmd)

	mockPipe.EXPECT().Exec().Return(nil, nil)

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

func TestMAUCache_MergeIntoMAUBatch_PipelineError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	mockPipe := redismock.NewMockPipeClient(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceIDs := []string{"ANDROID"}
	date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedErr := errors.New("pipeline error")

	mockCache.EXPECT().Pipeline(false).Return(mockPipe)
	mockPipe.EXPECT().PFMerge(gomock.Any(), gomock.Any(), gomock.Any())
	mockPipe.EXPECT().Expire(gomock.Any(), gomock.Any())
	mockPipe.EXPECT().Del(gomock.Any())
	dummyCmd := goredis.NewIntCmd(context.Background())
	mockPipe.EXPECT().PFCount(gomock.Any()).Return(dummyCmd)
	mockPipe.EXPECT().Exec().Return(nil, expectedErr)

	result, err := c.MergeIntoMAUBatch(envID, sourceIDs, date)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute merging MAU batch")
	assert.Nil(t, result)
}
