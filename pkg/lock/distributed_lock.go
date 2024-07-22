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
	"time"

	"github.com/google/uuid"

	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

const (
	unlockScript = `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
)

// DistributedLock represents a distributed lock
type DistributedLock struct {
	client     redisv3.Client
	key        string
	value      string
	expiration time.Duration
}

// NewDistributedLock creates a new DistributedLock
func NewDistributedLock(client redisv3.Client, key string, expiration time.Duration) *DistributedLock {
	return &DistributedLock{
		client:     client,
		key:        key,
		value:      uuid.New().String(),
		expiration: expiration,
	}
}

// Lock attempts to acquire the lock
func (dl *DistributedLock) Lock(ctx context.Context) (bool, error) {
	return dl.client.SetNX(ctx, dl.key, dl.value, dl.expiration)
}

// Unlock releases the lock
func (dl *DistributedLock) Unlock(ctx context.Context) (bool, error) {
	cmd := dl.client.Eval(ctx, unlockScript, []string{dl.key}, dl.value)
	res, err := cmd.Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
