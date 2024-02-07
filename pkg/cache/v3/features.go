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
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	featuresKind = "features"
	featuresTTL  = time.Duration(0)
)

type FeaturesCache interface {
	Get(environmentNamespace string) (*featureproto.Features, error)
	Put(features *featureproto.Features, environmentNamespace string) error
}

type featuresCache struct {
	cache cache.MultiGetCache
}

func NewFeaturesCache(c cache.MultiGetCache) FeaturesCache {
	return &featuresCache{cache: c}
}

func (c *featuresCache) Get(environmentNamespace string) (*featureproto.Features, error) {
	key := c.key(environmentNamespace)
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

func (c *featuresCache) Put(features *featureproto.Features, environmentNamespace string) error {
	buffer, err := proto.Marshal(features)
	if err != nil {
		return err
	}
	key := c.key(environmentNamespace)
	return c.cache.Put(key, buffer, featuresTTL)
}

func (c *featuresCache) key(environmentNamespace string) string {
	return fmt.Sprintf("%s:%s", environmentNamespace, featuresKind)
}
