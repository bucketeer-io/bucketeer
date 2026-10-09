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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package cache

import (
	"errors"
	"fmt"
	"time"

	redis "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage"
)

var (
	ErrNotFound    = errors.New("cache: not found")
	ErrInvalidType = errors.New("cache: not expected type")
)

type Cache interface {
	Getter
	Putter
	ConditionalPutter
}

type MultiGetCache interface {
	Cache
	MultiGetter
}

type MultiGetDeleteCountCache interface {
	MultiGetCache
	Deleter
	Counter
	PipeLiner
	Expirer
	PFGetter
	PFMerger
}

type Getter interface {
	Get(key interface{}) (interface{}, error)
}

type PFGetter interface {
	PFCount(keys ...string) (int64, error)
}

type MultiGetter interface {
	GetMulti(keys interface{}, ignoreNotFound bool) ([]interface{}, error)
	Scan(cursor, key, count interface{}) (uint64, []string, error)
	SMembers(key string) ([]string, error)
}

type Putter interface {
	Put(key interface{}, value interface{}, expiration time.Duration) error
}

// ConditionalPutter writes a value guarded by a monotonic generation so that
// concurrent writers cannot replace a newer snapshot with an older one.
type ConditionalPutter interface {
	// PutIfNewer atomically stores value under key and records generation under
	// genKey, unless genKey already describes the value currently stored at key
	// and holds a generation greater than the supplied one. In that case nothing
	// is written and accepted is false.
	//
	// The write is always accepted when key is missing or when the value at key
	// was written by a path that does not maintain genKey (see
	// IsStaleGeneration), so a stale generation marker can never block
	// repopulation. generation must not be negative.
	PutIfNewer(key, genKey string, value []byte, generation int64, expiration time.Duration) (accepted bool, err error)
	// DeleteWithGeneration removes key and genKey as a single atomic step so a
	// concurrent PutIfNewer can never observe (or be left with) a value whose
	// marker is gone, which would otherwise let a later stale write through.
	DeleteWithGeneration(key, genKey string) error
}

type Deleter interface {
	Delete(key string) error
}

type Counter interface {
	Increment(key string) (int64, error)
	IncrementBy(key string, value int64) (int64, error)
	PFAdd(key string, els ...string) (int64, error)
}

type PipeLiner interface {
	Pipeline(tx bool) redis.PipeClient
}

type Expirer interface {
	Expire(key string, expiration time.Duration) (bool, error)
}

type PFMerger interface {
	PFMerge(dest string, expiration time.Duration, keys ...string) error
}

// FIXME: remove after persistent-redis migration
type Lister interface {
	Keys(pattern string, maxSize int) ([]string, error)
}

func MakeKey(kind, id, environmentId string) string {
	if environmentId == storage.AdminEnvironmentID {
		return fmt.Sprintf("%s:%s", kind, id)
	}
	return fmt.Sprintf("%s:%s:%s", environmentId, kind, id)
}

func MakeKeyPrefix(kind, environmentId string) string {
	if environmentId == storage.AdminEnvironmentID {
		return fmt.Sprintf("%s:", kind)
	}
	return fmt.Sprintf("%s:%s:", environmentId, kind)
}

// MakeHashSlotKey creates a key to ensure that multiple keys are allocated in the same hash slot.
// https://redis.io/topics/cluster-spec#keys-hash-tags
func MakeHashSlotKey(hashTag, id, environmentId string) string {
	if environmentId == storage.AdminEnvironmentID {
		return fmt.Sprintf("{%s}%s", hashTag, id)
	}
	return fmt.Sprintf("{%s:%s}%s", environmentId, hashTag, id)
}

func Bytes(value interface{}) ([]byte, error) {
	b, ok := value.([]byte)
	if !ok {
		return nil, ErrInvalidType
	}
	return b, nil
}
