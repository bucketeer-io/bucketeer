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
	"fmt"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

const (
	autoOpsRuleKind     = "autoOpsRule"
	autoOpsRuleCacheTTL = time.Minute
)

type AutoOpsRulesCache interface {
	Get(environmentNamespace string) (*aoproto.AutoOpsRules, error)
	Put(autoOps *aoproto.AutoOpsRules, environmentNamespace string) error
}

type autoOpsRulesCache struct {
	cache cache.MultiGetCache
}

func NewAutoOpsRulesCache(c cache.MultiGetCache) AutoOpsRulesCache {
	return &autoOpsRulesCache{cache: c}
}

func (c *autoOpsRulesCache) Get(environmentNamespace string) (*aoproto.AutoOpsRules, error) {
	key := c.key(environmentNamespace)
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

func (c *autoOpsRulesCache) Put(autoOps *aoproto.AutoOpsRules, environmentNamespace string) error {
	buffer, err := proto.Marshal(autoOps)
	if err != nil {
		return err
	}
	key := c.key(environmentNamespace)
	return c.cache.Put(key, buffer, autoOpsRuleCacheTTL)
}

func (c *autoOpsRulesCache) key(environmentNamespace string) string {
	return fmt.Sprintf("%s:%s", environmentNamespace, autoOpsRuleKind)
}
