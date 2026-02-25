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
			desc:     "success: builds correct key format with hash tag",
			envID:    "env-123",
			sourceID: "ANDROID",
			date:     time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC),
			expected: "{env-123:ANDROID:au}:d:20260128",
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

func TestMAUCache_RecordDAU(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		envID       string
		sourceID    string
		userID      string
		date        time.Time
		setup       func(*mock.MockMultiGetDeleteCountCache, *redismock.MockPipeClient)
		expectedErr error
	}{
		{
			desc:     "success",
			envID:    "env-123",
			sourceID: "ANDROID",
			userID:   "user-456",
			date:     time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC),
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
				mc.EXPECT().Pipeline(false).Return(mp)
				mp.EXPECT().PFAdd("{env-123:ANDROID:au}:d:20260128", "user-456")
				mp.EXPECT().Expire("{env-123:ANDROID:au}:d:20260128", dauTTL)
				mp.EXPECT().Exec().Return(nil, nil)
			},
			expectedErr: nil,
		},
		{
			desc:     "empty userID: skip recording",
			envID:    "env-123",
			sourceID: "ANDROID",
			userID:   "",
			date:     time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC),
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
				// No pipeline operations expected
			},
			expectedErr: nil,
		},
		{
			desc:     "pipeline error",
			envID:    "env-123",
			sourceID: "ANDROID",
			userID:   "user-456",
			date:     time.Date(2026, 1, 28, 15, 30, 0, 0, time.UTC),
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
				mc.EXPECT().Pipeline(false).Return(mp)
				mp.EXPECT().PFAdd("{env-123:ANDROID:au}:d:20260128", "user-456")
				mp.EXPECT().Expire("{env-123:ANDROID:au}:d:20260128", dauTTL)
				mp.EXPECT().Exec().Return(nil, errors.New("redis connection error"))
			},
			expectedErr: errors.New("failed to record DAU"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
			mockPipe := redismock.NewMockPipeClient(ctrl)
			c := NewMAUCache(mockCache)

			p.setup(mockCache, mockPipe)

			err := c.RecordDAU(p.envID, p.sourceID, p.userID, p.date)
			if p.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), p.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
