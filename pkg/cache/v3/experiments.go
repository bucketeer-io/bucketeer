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
package v3

import (
	"fmt"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
)

type ExperimentsCache interface {
	Get(featureID string, featureVersion int32, environmentNamespace string) (*experimentproto.Experiments, error)
	Put(
		featureID string,
		featureVersion int32,
		experiments *experimentproto.Experiments,
		environmentNamespace string,
	) error
}

type experimentsCache struct {
	cache cache.Cache
}

func NewExperimentsCache(c cache.Cache) ExperimentsCache {
	return &experimentsCache{cache: c}
}

func (c *experimentsCache) Get(
	featureID string,
	featureVersion int32,
	environmentNamespace string,
) (*experimentproto.Experiments, error) {
	key := c.key(featureID, featureVersion, environmentNamespace)
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	experiments := &experimentproto.Experiments{}
	err = proto.Unmarshal(b, experiments)
	if err != nil {
		return nil, err
	}
	return experiments, nil
}

func (c *experimentsCache) Put(
	featureID string,
	featureVersion int32,
	experiments *experimentproto.Experiments,
	environmentNamespace string,
) error {
	buffer, err := proto.Marshal(experiments)
	if err != nil {
		return err
	}
	key := c.key(featureID, featureVersion, environmentNamespace)
	return c.cache.Put(key, buffer)
}

func (c *experimentsCache) key(featureID string, featureVersion int32, environmentNamespace string) string {
	return cache.MakeKey(
		"event_transformer:cache:experiments",
		fmt.Sprintf("%s:%d", featureID, featureVersion),
		environmentNamespace,
	)
}
