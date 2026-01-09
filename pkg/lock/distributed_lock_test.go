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

package lock

import (
	"context"
	"errors"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3/mock"
)

func TestNewDistributedLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	lock := NewDistributedLock(mockClient, time.Minute)

	assert.NotNil(t, lock)
	assert.Equal(t, time.Minute, lock.expiration)
}

func TestDistributedLock_Lock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	lock := NewDistributedLock(mockClient, time.Minute)

	ctx := context.Background()

	t.Run("successful lock acquisition", func(t *testing.T) {
		mockClient.EXPECT().
			SetNX(ctx, "test-key", gomock.Any(), time.Minute).
			Return(true, nil)

		acquired, value, err := lock.Lock(ctx, "test-key")

		assert.True(t, acquired)
		assert.NotEmpty(t, value)
		assert.NoError(t, err)
	})

	t.Run("failed lock acquisition", func(t *testing.T) {
		mockClient.EXPECT().
			SetNX(ctx, "test-key", gomock.Any(), time.Minute).
			Return(false, nil)

		acquired, value, err := lock.Lock(ctx, "test-key")

		assert.False(t, acquired)
		assert.NotEmpty(t, value)
		assert.NoError(t, err)
	})

	t.Run("error during lock acquisition", func(t *testing.T) {
		mockClient.EXPECT().
			SetNX(ctx, "test-key", gomock.Any(), time.Minute).
			Return(false, errors.New("redis error"))

		acquired, value, err := lock.Lock(ctx, "test-key")

		assert.False(t, acquired)
		assert.NotEmpty(t, value)
		assert.Error(t, err)
		assert.Equal(t, "redis error", err.Error())
	})
}

func TestDistributedLock_Unlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	lock := NewDistributedLock(mockClient, time.Minute)

	ctx := context.Background()

	t.Run("successful unlock", func(t *testing.T) {
		successCmd := goredis.NewCmdResult(int64(1), nil)
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"test-key"}, "test-value").
			Return(successCmd)

		unlocked, err := lock.Unlock(ctx, "test-key", "test-value")

		assert.True(t, unlocked)
		assert.NoError(t, err)
	})

	t.Run("unsuccessful unlock", func(t *testing.T) {
		failCmd := goredis.NewCmdResult(int64(0), nil)
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"test-key"}, "test-value").
			Return(failCmd)

		unlocked, err := lock.Unlock(ctx, "test-key", "test-value")

		assert.False(t, unlocked)
		assert.NoError(t, err)
	})

	t.Run("error during unlock", func(t *testing.T) {
		errorCmd := goredis.NewCmdResult(nil, errors.New("eval error"))
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"test-key"}, "test-value").
			Return(errorCmd)

		unlocked, err := lock.Unlock(ctx, "test-key", "test-value")

		assert.False(t, unlocked)
		assert.Error(t, err)
		assert.Equal(t, "eval error", err.Error())
	})
}

func TestDistributedLock_MultipleKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	lock := NewDistributedLock(mockClient, time.Minute)

	ctx := context.Background()

	t.Run("lock and unlock multiple keys", func(t *testing.T) {
		// Lock key1
		mockClient.EXPECT().
			SetNX(ctx, "key1", gomock.Any(), time.Minute).
			Return(true, nil)

		acquired1, value1, err := lock.Lock(ctx, "key1")
		assert.True(t, acquired1)
		assert.NotEmpty(t, value1)
		assert.NoError(t, err)

		// Lock key2
		mockClient.EXPECT().
			SetNX(ctx, "key2", gomock.Any(), time.Minute).
			Return(true, nil)

		acquired2, value2, err := lock.Lock(ctx, "key2")
		assert.True(t, acquired2)
		assert.NotEmpty(t, value2)
		assert.NoError(t, err)

		// Unlock key1
		successCmd1 := goredis.NewCmdResult(int64(1), nil)
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"key1"}, value1).
			Return(successCmd1)

		unlocked1, err := lock.Unlock(ctx, "key1", value1)
		assert.True(t, unlocked1)
		assert.NoError(t, err)

		// Unlock key2
		successCmd2 := goredis.NewCmdResult(int64(1), nil)
		mockClient.EXPECT().
			Eval(ctx, unlockScript, []string{"key2"}, value2).
			Return(successCmd2)

		unlocked2, err := lock.Unlock(ctx, "key2", value2)
		assert.True(t, unlocked2)
		assert.NoError(t, err)
	})
}
