// Copyright 2022 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/storage"
)

var (
	ErrNotFound    = errors.New("cache: not found")
	ErrInvalidType = errors.New("cache: not expected type")
)

type Cache interface {
	Getter
	Putter
}

type MultiGetCache interface {
	Cache
	MultiGetter
}

type MultiGetDeleteCache interface {
	MultiGetCache
	Deleter
}

type Getter interface {
	Get(key interface{}) (interface{}, error)
}

type MultiGetter interface {
	GetMulti(keys interface{}) ([]interface{}, error)
	Scan(cursor, key, count interface{}) (uint64, []string, error)
}

type Putter interface {
	Put(key interface{}, value interface{}) error
}

type Deleter interface {
	Delete(key string) error
}

// FIXME: remove after persistent-redis migration
type Lister interface {
	Keys(pattern string, maxSize int) ([]string, error)
}

func MakeKey(kind, id, environmentNamespace string) string {
	if environmentNamespace == storage.AdminEnvironmentNamespace {
		return fmt.Sprintf("%s:%s", kind, id)
	}
	return fmt.Sprintf("%s:%s:%s", environmentNamespace, kind, id)
}

func MakeKeyPrefix(kind, environmentNamespace string) string {
	if environmentNamespace == storage.AdminEnvironmentNamespace {
		return fmt.Sprintf("%s:", kind)
	}
	return fmt.Sprintf("%s:%s:", environmentNamespace, kind)
}

// MakeHashSlotKey creates a key to ensure that multiple keys are allocated in the same hash slot.
// https://redis.io/topics/cluster-spec#keys-hash-tags
func MakeHashSlotKey(hashTag, id, environmentNamespace string) string {
	if environmentNamespace == storage.AdminEnvironmentNamespace {
		return fmt.Sprintf("{%s}%s", hashTag, id)
	}
	return fmt.Sprintf("{%s:%s}%s", environmentNamespace, hashTag, id)
}

func Bytes(value interface{}) ([]byte, error) {
	b, ok := value.([]byte)
	if !ok {
		return nil, ErrInvalidType
	}
	return b, nil
}
