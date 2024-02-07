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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package v3

import (
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	segmentUsersKind    = "segment_users"
	segmentUsersMaxSize = int64(100)
	segmentUsersTTL     = time.Duration(0)
)

type SegmentUsersCache interface {
	Get(segmentID, environmentNamespace string) (*featureproto.SegmentUsers, error)
	GetAll(environmentNamespace string) ([]*featureproto.SegmentUsers, error)
	Put(segmentUsers *featureproto.SegmentUsers, environmentNamespace string) error
}

type segmentUsersCache struct {
	cache cache.MultiGetCache
}

func NewSegmentUsersCache(c cache.MultiGetCache) SegmentUsersCache {
	return &segmentUsersCache{cache: c}
}

func (c *segmentUsersCache) Get(segmentID, environmentNamespace string) (*featureproto.SegmentUsers, error) {
	key := c.key(segmentID, environmentNamespace)
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	segmentUsers := &featureproto.SegmentUsers{}
	err = proto.Unmarshal(b, segmentUsers)
	if err != nil {
		return nil, err
	}
	return segmentUsers, nil
}

func (c *segmentUsersCache) GetAll(environmentNamespace string) ([]*featureproto.SegmentUsers, error) {
	keys, err := c.scan(environmentNamespace)
	if err != nil {
		return nil, err
	}
	users, err := c.cache.GetMulti(keys)
	if err != nil {
		return nil, err
	}
	segmentUsers := []*featureproto.SegmentUsers{}
	for _, value := range users {
		b, err := cache.Bytes(value)
		if err != nil {
			return nil, err
		}
		su := &featureproto.SegmentUsers{}
		err = proto.Unmarshal(b, su)
		if err != nil {
			return nil, err
		}
		segmentUsers = append(segmentUsers, su)
	}
	return segmentUsers, nil
}

func (c *segmentUsersCache) Put(segmentUsers *featureproto.SegmentUsers, environmentNamespace string) error {
	buffer, err := proto.Marshal(segmentUsers)
	if err != nil {
		return err
	}
	key := c.key(segmentUsers.SegmentId, environmentNamespace)
	return c.cache.Put(key, buffer, segmentUsersTTL)
}

func (c *segmentUsersCache) scan(environmentNamespace string) ([]string, error) {
	keyPrefix := cache.MakeKeyPrefix(segmentUsersKind, environmentNamespace)
	key := keyPrefix + "*"
	var cursor uint64
	var k []string
	var err error
	keys := []string{}
	for {
		cursor, k, err = c.cache.Scan(cursor, key, segmentUsersMaxSize)
		if err != nil {
			break
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *segmentUsersCache) key(segmentID, environmentNamespace string) string {
	return cache.MakeKey(segmentUsersKind, segmentID, environmentNamespace)
}
