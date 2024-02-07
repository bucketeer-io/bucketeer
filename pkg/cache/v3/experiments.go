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
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

const (
	experimentsKind    = "experiment"
	experimentCacheTTL = time.Minute
)

type ExperimentsCache interface {
	Get(environmentNamespace string) (*exproto.Experiments, error)
	Put(experiments *exproto.Experiments, environmentNamespace string) error
}

type experimentsCache struct {
	cache cache.MultiGetCache
}

func NewExperimentsCache(c cache.MultiGetCache) ExperimentsCache {
	return &experimentsCache{cache: c}
}

func (c *experimentsCache) Get(environmentNamespace string) (*exproto.Experiments, error) {
	key := c.key(environmentNamespace)
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	experiments := &exproto.Experiments{}
	err = proto.Unmarshal(b, experiments)
	if err != nil {
		return nil, err
	}
	return experiments, nil
}

func (c *experimentsCache) Put(experiments *exproto.Experiments, environmentNamespace string) error {
	buffer, err := proto.Marshal(experiments)
	if err != nil {
		return err
	}
	key := c.key(environmentNamespace)
	return c.cache.Put(key, buffer, experimentCacheTTL)
}

func (c *experimentsCache) key(environmentNamespace string) string {
	return fmt.Sprintf("%s:%s", environmentNamespace, experimentsKind)
}
