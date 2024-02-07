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
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

const (
	environmentAPIKeyKind             = "environment_apikey"
	environmentAPIKeyTTL              = 1 * time.Minute
	EnvironmentAPIKeyEvictionInterval = 10 * time.Second
)

type EnvironmentAPIKeyCache interface {
	Get(string) (*accountproto.EnvironmentAPIKey, error)
	Put(*accountproto.EnvironmentAPIKey) error
}

type environmentAPIKeyCache struct {
	cache cache.Cache
}

func NewEnvironmentAPIKeyCache(c cache.Cache) EnvironmentAPIKeyCache {
	return &environmentAPIKeyCache{cache: c}
}

func (c *environmentAPIKeyCache) Get(id string) (*accountproto.EnvironmentAPIKey, error) {
	key := c.key(id)
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	environmentAPIKey := &accountproto.EnvironmentAPIKey{}
	if err := proto.Unmarshal(b, environmentAPIKey); err != nil {
		return nil, err
	}
	return environmentAPIKey, nil
}

func (c *environmentAPIKeyCache) Put(environmentAPIKey *accountproto.EnvironmentAPIKey) error {
	buffer, err := proto.Marshal(environmentAPIKey)
	if err != nil {
		return err
	}
	key := c.key(environmentAPIKey.ApiKey.Id)
	return c.cache.Put(key, buffer, environmentAPIKeyTTL)
}

func (c *environmentAPIKeyCache) key(id string) string {
	// always use AdminEnvironmentNamespace because we'd like to get APIKey and environment_namespace only by id
	return cache.MakeKey(environmentAPIKeyKind, id, storage.AdminEnvironmentNamespace)
}
