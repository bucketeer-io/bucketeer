// Copyright 2025 The Bucketeer Authors.
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

package calculator

import (
	"context"
	"fmt"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/lock"
	redisv3 "github.com/bucketeer-io/bucketeer/pkg/redis/v3"
)

const (
	experimentLockKind = "experiment_lock"
)

// ExperimentLock represents a distributed lock for experiments
type ExperimentLock struct {
	lock *lock.DistributedLock
}

// NewExperimentLock creates a new ExperimentLock
func NewExperimentLock(client redisv3.Client, lockTTL time.Duration) *ExperimentLock {
	return &ExperimentLock{
		lock: lock.NewDistributedLock(client, lockTTL),
	}
}

// Lock attempts to acquire the lock for a specific experiment
func (el *ExperimentLock) Lock(ctx context.Context, environmentID, experimentID string) (bool, string, error) {
	lockKey := el.newLockKey(environmentID, experimentID)
	return el.lock.Lock(ctx, lockKey)
}

// Unlock releases the lock for a specific experiment
func (el *ExperimentLock) Unlock(ctx context.Context, environmentID, experimentID, value string) (bool, error) {
	lockKey := el.newLockKey(environmentID, experimentID)
	return el.lock.Unlock(ctx, lockKey, value)
}

// newLockKey generates the lock key for a specific experiment
func (el *ExperimentLock) newLockKey(environmentID, experimentID string) string {
	return fmt.Sprintf("%s:%s:%s", environmentID, experimentLockKind, experimentID)
}
