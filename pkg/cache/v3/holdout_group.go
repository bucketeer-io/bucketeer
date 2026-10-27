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
package v3

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	holdoutGroupKind = "holdout_group"
)

// HoldoutGroupCache stores one serving snapshot per environment.
// An empty Id is the sentinel for "nothing is serving".
type HoldoutGroupCache interface {
	Get(environmentID string) (*ftproto.HoldoutGroup, error)
	Put(group *ftproto.HoldoutGroup, environmentID string) error
	Evict(environmentID string) error
}

type holdoutGroupCache struct {
	cache cache.Cache
	ttl   time.Duration
}

func NewHoldoutGroupCache(c cache.Cache, ttl time.Duration) HoldoutGroupCache {
	return &holdoutGroupCache{cache: c, ttl: ttl}
}

func (c *holdoutGroupCache) Get(environmentID string) (*ftproto.HoldoutGroup, error) {
	value, err := c.cache.Get(c.key(environmentID))
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	group := &ftproto.HoldoutGroup{}
	if err := proto.Unmarshal(b, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (c *holdoutGroupCache) Put(group *ftproto.HoldoutGroup, environmentID string) error {
	buffer, err := proto.Marshal(group)
	if err != nil {
		return err
	}
	return c.cache.Put(c.key(environmentID), buffer, c.ttl)
}

func (c *holdoutGroupCache) Evict(environmentID string) error {
	return evictKey(c.cache, c.key(environmentID))
}

func (c *holdoutGroupCache) key(environmentID string) string {
	return fmt.Sprintf("%s:%s", environmentID, holdoutGroupKind)
}
