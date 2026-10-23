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
	aoproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

const (
	autoOpsRuleKind = "autoOpsRule"
)

type AutoOpsRulesCache interface {
	Get(environmentId string) (*aoproto.AutoOpsRules, error)
	Put(autoOps *aoproto.AutoOpsRules, environmentId string) error
	// PutIfNewer stores the snapshot unless a newer generation is cached.
	// generation is the UnixNano time captured before reading the snapshot.
	PutIfNewer(autoOps *aoproto.AutoOpsRules, environmentId string, generation int64) (bool, error)
	Evict(environmentId string) error
}

type autoOpsRulesCache struct {
	cache cache.MultiGetCache
	ttl   time.Duration
}

func NewAutoOpsRulesCache(c cache.MultiGetCache, ttl time.Duration) AutoOpsRulesCache {
	return &autoOpsRulesCache{cache: c, ttl: ttl}
}

func (c *autoOpsRulesCache) Get(environmentId string) (*aoproto.AutoOpsRules, error) {
	key := c.key(environmentId)
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	autoOps := &aoproto.AutoOpsRules{}
	err = proto.Unmarshal(b, autoOps)
	if err != nil {
		return nil, err
	}
	return autoOps, nil
}

func (c *autoOpsRulesCache) Put(autoOps *aoproto.AutoOpsRules, environmentId string) error {
	buffer, err := proto.Marshal(autoOps)
	if err != nil {
		return err
	}
	key := c.key(environmentId)
	return c.cache.Put(key, buffer, c.ttl)
}

func (c *autoOpsRulesCache) PutIfNewer(
	autoOps *aoproto.AutoOpsRules,
	environmentId string,
	generation int64,
) (bool, error) {
	buffer, err := proto.Marshal(autoOps)
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

func (c *autoOpsRulesCache) Evict(environmentId string) error {
	return evictKey(c.cache, c.key(environmentId))
}

// generationKey shares the hash slot of the snapshot key, as the conditional write requires
func (c *autoOpsRulesCache) generationKey(environmentId string) string {
	return fmt.Sprintf("{%s}:gen", c.key(environmentId))
}

func (c *autoOpsRulesCache) key(environmentId string) string {
	return fmt.Sprintf("%s:%s", environmentId, autoOpsRuleKind)
}
