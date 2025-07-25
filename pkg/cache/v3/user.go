// Copyright 2025 The Bucketeer Authors.
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
	"strings"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	userAttributeKind     = "user_attr"
	userAttributesMaxSize = int64(100)
)

type UserAttributesCache interface {
	GetUserAttributeKeyAll(environmentId string) ([]string, error)
	Put(userAttributes *userproto.UserAttributes, ttl time.Duration) error
}

type userAttributesCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewUserAttributesCache(c cache.MultiGetDeleteCountCache) UserAttributesCache {
	return &userAttributesCache{cache: c}
}

func (u *userAttributesCache) GetUserAttributeKeyAll(environmentId string) ([]string, error) {
	scanKey := u.key(environmentId) + ":*"
	var cursor uint64
	var allKeys []string

	for {
		var keys []string
		var err error
		cursor, keys, err = u.cache.Scan(cursor, scanKey, userAttributesMaxSize)
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			break
		}
	}

	// Extract UserAttributeKey from the full keys
	// Key format: environmentId:user_attr:attributeKey
	attributeKeys := u.extractAttributeKeys(allKeys)
	return attributeKeys, nil
}

func (u *userAttributesCache) extractAttributeKeys(fullKeys []string) []string {
	attributeKeys := []string{}
	for _, fullKey := range fullKeys {
		// Split by ":" and get the last part which is the attribute key
		parts := strings.Split(fullKey, ":")
		for i, part := range parts {
			// If userAttrKindIndex is found, use the next element as the key
			if part == userAttributeKind && i+1 < len(parts) {
				attributeKeys = append(attributeKeys, parts[i+1])
				break
			}
		}
	}
	return attributeKeys
}

func (u *userAttributesCache) Put(userAttributes *userproto.UserAttributes, ttl time.Duration) error {
	if userAttributes == nil {
		return errors.New("user attributes is nil")
	}
	pipe := u.cache.Pipeline(false)
	for _, attribute := range userAttributes.UserAttributes {
		key := u.key(userAttributes.EnvironmentId) + ":" + attribute.Key
		for _, value := range attribute.Values {
			pipe.SAdd(key, value)
		}
		pipe.Expire(key, ttl)
	}
	_, err := pipe.Exec()
	if err != nil {
		return err
	}
	return nil
}

// We use a Redis Cluster “hash tag” so all user_attr keys for an environment
// live in the same slot. On standalone mode the `{…}` is just a literal.
// The key format is: {environmentId}:user_attr:{attributeKey}
func (u *userAttributesCache) key(environmentId string) string {
	return fmt.Sprintf("{%s}:%s", environmentId, userAttributeKind)
}
