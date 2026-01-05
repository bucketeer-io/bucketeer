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

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

const (
	userAttributeKind = "user_attr"
)

type UserAttributesCache interface {
	// Returns all attribute keys for the given environment.
	GetUserAttributeKeyAll(environmentID string) ([]string, error)
	// Stores values for each attribute key under the given TTL.
	Put(userAttributes *userproto.UserAttributes, ttl time.Duration) error
}

type userAttributesCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewUserAttributesCache(c cache.MultiGetDeleteCountCache) UserAttributesCache {
	return &userAttributesCache{cache: c}
}

// key returns the base prefix for all user_attr entries,
// e.g. "env123:user_attr"
func (u *userAttributesCache) key(environmentID string) string {
	return fmt.Sprintf("%s:%s", environmentID, userAttributeKind)
}

// Put writes each attribute's values into its own Set, and also
// adds the attribute.Key to an index Set so we can list them later.
func (u *userAttributesCache) Put(
	userAttributes *userproto.UserAttributes,
	ttl time.Duration,
) error {
	if userAttributes == nil {
		return errors.New("userAttributes cannot be nil")
	}

	pipe := u.cache.Pipeline(false)
	// indexKey holds the list of all attribute keys
	indexKey := u.key(userAttributes.EnvironmentId) + ":keys"

	for _, attribute := range userAttributes.UserAttributes {
		// 1) Store attribute values in their own set: env:user_attr:country -> ["US", "JP"]
		attrKey := fmt.Sprintf("%s:%s", u.key(userAttributes.EnvironmentId), attribute.Key)
		// convert []string → []interface{} for SAdd
		members := make([]interface{}, len(attribute.Values))
		for i, v := range attribute.Values {
			members[i] = v
		}
		// SAdd can add multiple values at once
		pipe.SAdd(attrKey, members...)
		pipe.Expire(attrKey, ttl)

		// 2) Store attribute key in index set: env:user_attr:keys -> ["country", "plan_type"]
		// This avoids scanning Redis keys and works efficiently with Redis clusters
		pipe.SAdd(indexKey, attribute.Key)
		pipe.Expire(indexKey, ttl)
	}

	_, err := pipe.Exec()
	return err
}

// GetUserAttributeKeyAll fetches the complete list of attribute keys
// via a single SMEMBERS call on the index Set.
func (u *userAttributesCache) GetUserAttributeKeyAll(
	environmentID string,
) ([]string, error) {
	indexKey := u.key(environmentID) + ":keys"
	keys, err := u.cache.SMembers(indexKey)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
