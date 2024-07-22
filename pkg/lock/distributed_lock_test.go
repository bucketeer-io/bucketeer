// Copyright 2024 The Bucketeer Authors.
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

package lock

import (
	"context"
	"errors"
	"testing"
	"time"

	goredis "github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/redis/v3/mock"
)

func TestNewDistributedLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	lock := NewDistributedLock(mockClient, "test-key", time.Minute)

	assert.NotNil(t, lock)
	assert.Equal(t, "test-key", lock.key)
	assert.Equal(t, time.Minute, lock.expiration)
	assert.NotEmpty(t, lock.value)
}

func TestDistributedLock_Lock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	lock := NewDistributedLock(mockClient, "test-key", time.Minute)

	ctx := context.Background()

	t.Run("successful lock acquisition", func(t *testing.T) {
		mockClient.EXPECT().
			SetNX(ctx, "test-key", gomock.Any(), time.Minute).
			Return(true, nil)

		acquired, err := lock.Lock(ctx)

		assert.True(t, acquired)
		assert.NoError(t, err)
	})

	t.Run("failed lock acquisition", func(t *testing.T) {
		mockClient.EXPECT().
			SetNX(ctx, "test-key", gomock.Any(), time.Minute).
			Return(false, nil)

		acquired, err := lock.Lock(ctx)

		assert.False(t, acquired)
		assert.NoError(t, err)
	})
}

func TestDistributedLock_Unlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	lock := NewDistributedLock(mockClient, "test-key", time.Minute)

	ctx := context.Background()

	t.Run("successful unlock", func(t *testing.T) {
		successCmd := goredis.NewCmdResult(int64(1), nil)
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"test-key"}, gomock.Any()).
			Return(successCmd)

		unlocked, err := lock.Unlock(ctx)

		assert.True(t, unlocked)
		assert.NoError(t, err)
	})

	t.Run("unsuccessful unlock", func(t *testing.T) {
		failCmd := goredis.NewCmdResult(int64(0), nil)
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"test-key"}, gomock.Any()).
			Return(failCmd)

		unlocked, err := lock.Unlock(ctx)

		assert.False(t, unlocked)
		assert.NoError(t, err)
	})

	t.Run("error during unlock", func(t *testing.T) {
		errorCmd := goredis.NewCmdResult(nil, errors.New("eval error"))
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"test-key"}, gomock.Any()).
			Return(errorCmd)

		unlocked, err := lock.Unlock(ctx)

		assert.False(t, unlocked)
		assert.Error(t, err)
		assert.Equal(t, "eval error", err.Error())
	})
}
