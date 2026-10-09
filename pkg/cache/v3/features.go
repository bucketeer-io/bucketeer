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
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	featuresKind = "features"
	// featuresGenerationSuffix is appended to the features key (wrapped in a
	// Redis hash tag) to form the generation marker key.
	featuresGenerationSuffix = "gen"
)

type FeaturesCache interface {
	Get(environmentId string) (*featureproto.Features, error)
	// Put unconditionally overwrites the cached snapshot. Snapshot writers that
	// may run concurrently (subscriber refresher, batch cacher) must use
	// PutIfNewer instead.
	Put(features *featureproto.Features, environmentId string) error
	// PutIfNewer stores the snapshot only if generation is not older than the
	// generation of the snapshot currently cached. generation should be the
	// time (UnixNano) at which the writer started reading the snapshot from the
	// source of truth, captured before the first read is issued. Returns false
	// when the write was rejected as stale; that is not an error.
	PutIfNewer(features *featureproto.Features, environmentId string, generation int64) (bool, error)
	// Evict removes the snapshot and its generation marker.
	Evict(environmentId string) error
}

type featuresCache struct {
	cache cache.Cache
	ttl   time.Duration
}

func NewFeaturesCache(c cache.Cache, ttl time.Duration) FeaturesCache {
	return &featuresCache{cache: c, ttl: ttl}
}

func (c *featuresCache) Get(environmentId string) (*featureproto.Features, error) {
	key := c.key(environmentId)
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	features := &featureproto.Features{}
	err = proto.Unmarshal(b, features)
	if err != nil {
		return nil, err
	}
	return features, nil
}

func (c *featuresCache) Put(features *featureproto.Features, environmentId string) error {
	if features == nil {
		return errors.New("features cannot be nil")
	}
	buffer, err := proto.Marshal(features)
	if err != nil {
		return err
	}
	key := c.key(environmentId)
	return c.cache.Put(key, buffer, c.ttl)
}

func (c *featuresCache) PutIfNewer(
	features *featureproto.Features,
	environmentId string,
	generation int64,
) (bool, error) {
	if features == nil {
		return false, errors.New("features cannot be nil")
	}
	buffer, err := proto.Marshal(features)
	if err != nil {
		return false, err
	}
	return c.cache.PutIfNewer(
		c.key(environmentId),
		c.generationKey(environmentId),
		buffer,
		generation,
		c.ttl,
	)
}

func (c *featuresCache) Evict(environmentId string) error {
	if err := evictKey(c.cache, c.key(environmentId)); err != nil {
		return err
	}
	return evictKey(c.cache, c.generationKey(environmentId))
}

func (c *featuresCache) key(environmentId string) string {
	return fmt.Sprintf("%s:%s", environmentId, featuresKind)
}

// generationKey returns the key holding the generation marker for the
// environment's snapshot. The features key is wrapped in a hash tag so both
// keys land in the same Redis Cluster slot, which the multi-key Lua script
// requires. The features key itself is unchanged so existing readers (and
// pods running older builds) keep working.
func (c *featuresCache) generationKey(environmentId string) string {
	return fmt.Sprintf("{%s}:%s", c.key(environmentId), featuresGenerationSuffix)
}
