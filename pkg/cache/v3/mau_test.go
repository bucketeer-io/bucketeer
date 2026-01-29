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
	redismock "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
)

func TestMAUCache_DAUKey(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		envID    string
		sourceID string
		date     time.Time
		expected string
	}{
		{
			desc:     "success: builds correct key format",
			envID:    "env-123",
			sourceID: "ANDROID",
			date:     time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC),
			expected: "env-123:ANDROID:dau:20260128",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			c := &mauCache{}
			actual := c.dauKey(p.envID, p.sourceID, p.date)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestMAUCache_RecordDAU_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	mockPipe := redismock.NewMockPipeClient(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceID := "ANDROID"
	userID := "user-456"
	date := time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC)
	expectedKey := "env-123:ANDROID:dau:20260128"

	mockCache.EXPECT().Pipeline(false).Return(mockPipe)
	mockPipe.EXPECT().PFAdd(expectedKey, userID)
	mockPipe.EXPECT().Expire(expectedKey, dauTTL)
	mockPipe.EXPECT().Exec().Return(nil, nil)

	err := c.RecordDAU(envID, sourceID, userID, date)
	assert.NoError(t, err)
}

func TestMAUCache_RecordDAU_EmptyUserID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	mockPipe := redismock.NewMockPipeClient(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceID := "ANDROID"
	userID := "" // HERE
	date := time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC)

	mockPipe.EXPECT().PFAdd(gomock.Any(), gomock.Any()).Times(0)

	err := c.RecordDAU(envID, sourceID, userID, date)
	assert.NoError(t, err)
}

func TestMAUCache_RecordDAU_PipelineError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
	mockPipe := redismock.NewMockPipeClient(ctrl)
	c := NewMAUCache(mockCache)

	envID := "env-123"
	sourceID := "ANDROID"
	userID := "user-456"
	date := time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC)
	expectedKey := "env-123:ANDROID:dau:20260128"
	expectedErr := errors.New("redis connection error")

	mockCache.EXPECT().Pipeline(false).Return(mockPipe)
	mockPipe.EXPECT().PFAdd(expectedKey, userID)
	mockPipe.EXPECT().Expire(expectedKey, dauTTL)
	mockPipe.EXPECT().Exec().Return(nil, expectedErr)

	err := c.RecordDAU(envID, sourceID, userID, date)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to record DAU")
	assert.Contains(t, err.Error(), "redis connection error")
}
