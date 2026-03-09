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

func TestDAUKey(t *testing.T) {
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
			date:     time.Date(2026, 1, 28, 0, 0, 0, 0, time.UTC),
			expected: "env-123:dau:ANDROID:20260128",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := dauKey(p.envID, p.sourceID, p.date)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestDAUCache_RecordDAUBatch(t *testing.T) {
	testDate := time.Date(2026, 1, 28, 0, 0, 0, 0, time.UTC)
	t.Parallel()
	patterns := []struct {
		desc        string
		records     []DAURecord
		setup       func(*mock.MockMultiGetDeleteCountCache, *redismock.MockPipeClient)
		expectedErr error
	}{
		{
			desc: "success: single record with one user",
			records: []DAURecord{
				{Date: testDate, EnvID: "env-123", SourceID: "ANDROID", UserIDs: []string{"user-456"}},
			},
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
				mc.EXPECT().Pipeline(false).Return(mp)
				mp.EXPECT().PFAdd("env-123:dau:ANDROID:20260128", "user-456")
				mp.EXPECT().Expire("env-123:dau:ANDROID:20260128", dauTTL)
				mp.EXPECT().Exec().Return(nil, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "success: multiple records in one pipeline",
			records: []DAURecord{
				{Date: testDate, EnvID: "env-123", SourceID: "ANDROID", UserIDs: []string{"user-1"}},
				{Date: testDate, EnvID: "env-123", SourceID: "IOS", UserIDs: []string{"user-2"}},
				{Date: testDate, EnvID: "env-456", SourceID: "WEB", UserIDs: []string{"user-3"}},
			},
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
				mc.EXPECT().Pipeline(false).Return(mp)
				mp.EXPECT().PFAdd("env-123:dau:ANDROID:20260128", "user-1")
				mp.EXPECT().Expire("env-123:dau:ANDROID:20260128", dauTTL)
				mp.EXPECT().PFAdd("env-123:dau:IOS:20260128", "user-2")
				mp.EXPECT().Expire("env-123:dau:IOS:20260128", dauTTL)
				mp.EXPECT().PFAdd("env-456:dau:WEB:20260128", "user-3")
				mp.EXPECT().Expire("env-456:dau:WEB:20260128", dauTTL)
				mp.EXPECT().Exec().Return(nil, nil)
			},
			expectedErr: nil,
		},
		{
			desc:    "empty records: no-op",
			records: []DAURecord{},
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
			},
			expectedErr: nil,
		},
		{
			desc: "all empty userIDs: no-op",
			records: []DAURecord{
				{Date: testDate, EnvID: "env-123", SourceID: "ANDROID", UserIDs: []string{}},
				{Date: testDate, EnvID: "env-456", SourceID: "IOS", UserIDs: nil},
			},
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
			},
			expectedErr: nil,
		},
		{
			desc: "multiple users in single PFADD",
			records: []DAURecord{
				{Date: testDate, EnvID: "env-123", SourceID: "ANDROID", UserIDs: []string{"user-1", "user-2", "user-3"}},
			},
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
				mc.EXPECT().Pipeline(false).Return(mp)
				mp.EXPECT().PFAdd("env-123:dau:ANDROID:20260128", "user-1", "user-2", "user-3")
				mp.EXPECT().Expire("env-123:dau:ANDROID:20260128", dauTTL)
				mp.EXPECT().Exec().Return(nil, nil)
			},
			expectedErr: nil,
		},
		{
			desc: "pipeline error",
			records: []DAURecord{
				{Date: testDate, EnvID: "env-123", SourceID: "ANDROID", UserIDs: []string{"user-456"}},
			},
			setup: func(mc *mock.MockMultiGetDeleteCountCache, mp *redismock.MockPipeClient) {
				mc.EXPECT().Pipeline(false).Return(mp)
				mp.EXPECT().PFAdd("env-123:dau:ANDROID:20260128", "user-456")
				mp.EXPECT().Expire("env-123:dau:ANDROID:20260128", dauTTL)
				mp.EXPECT().Exec().Return(nil, errors.New("redis connection error"))
			},
			expectedErr: errors.New("failed to record DAU batch"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := mock.NewMockMultiGetDeleteCountCache(ctrl)
			mockPipe := redismock.NewMockPipeClient(ctrl)
			c := NewDAUCache(mockCache)

			p.setup(mockCache, mockPipe)

			err := c.RecordDAUBatch(p.records)
			if p.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), p.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
