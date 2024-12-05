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
	"strconv"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
)

type EventCounterCache interface {
	GetEventCounts(keys []string) ([]float64, error)
	GetEventCountsV2(keys [][]string) ([]float64, error)
	GetUserCount(key string) (int64, error)
	GetUserCounts(keys []string) ([]float64, error)
	GetUserCountsV2(keys []string) ([]float64, error)
	MergeMultiKeys(dest string, keys []string) error
	DeleteKey(key string) error
	UpdateUserCount(key, userID string) error
	ExpireKey(key string, expiration time.Duration) (bool, error)
}

type eventCounterCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewEventCountCache(c cache.MultiGetDeleteCountCache) EventCounterCache {
	return &eventCounterCache{cache: c}
}

func (c *eventCounterCache) GetEventCounts(keys []string) ([]float64, error) {
	values, err := c.cache.GetMulti(keys, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple keys: %w", err)
	}
	return c.getEventValues(values)
}

func (*eventCounterCache) getEventValues(values []interface{}) ([]float64, error) {
	eventVals := make([]float64, 0, len(values))
	for _, v := range values {
		var str string
		switch v := v.(type) {
		case []byte:
			str = string(v)
		case string:
			str = v
		default:
			return nil, fmt.Errorf("unexpected value type: %v", v)
		}
		if str == "" {
			str = "0"
		}
		float, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float: %w", err)
		}
		eventVals = append(eventVals, float)
	}
	return eventVals, nil
}

func (c *eventCounterCache) GetEventCountsV2(keys [][]string) ([]float64, error) {
	eventVals := make([]float64, 0, len(keys))
	for _, day := range keys {
		values, err := c.cache.GetMulti(day, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get multiple keys: %w", err)
		}
		dayTotal, err := c.getEventValues(values)
		if err != nil {
			return nil, err
		}
		var totalVal float64
		for _, v := range dayTotal {
			totalVal += v
		}
		eventVals = append(eventVals, totalVal)
	}
	return eventVals, nil
}

func (c *eventCounterCache) GetUserCount(key string) (int64, error) {
	return c.cache.PFCount(key)
}

func (c *eventCounterCache) GetUserCounts(keys []string) ([]float64, error) {
	userVals := make([]float64, 0, len(keys))
	for _, key := range keys {
		count, err := c.cache.PFCount(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get PFCount for key %s: %w", key, err)
		}
		userVals = append(userVals, float64(count))
	}
	return userVals, nil
}

func (c *eventCounterCache) GetUserCountsV2(keys []string) ([]float64, error) {
	return c.GetUserCounts(keys)
}

func (c *eventCounterCache) UpdateUserCount(key, userID string) error {
	_, err := c.cache.PFAdd(key, userID)
	if err != nil {
		return fmt.Errorf("failed to update user count: %w", err)
	}
	return nil
}

func (c *eventCounterCache) MergeMultiKeys(dest string, keys []string) error {
	if err := c.cache.PFMerge(dest, keys...); err != nil {
		return fmt.Errorf("failed to merge keys: %w", err)
	}
	return nil
}

func (c *eventCounterCache) DeleteKey(key string) error {
	if err := c.cache.Delete(key); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}
func (c *eventCounterCache) ExpireKey(key string, expiration time.Duration) (bool, error) {
	ok, err := c.cache.Expire(key, expiration)
	if err != nil {
		return false, fmt.Errorf("failed to set expiration for key: %w", err)
	}
	return ok, nil
}
