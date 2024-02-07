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

	goredis "github.com/go-redis/redis"

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
}

type eventCounterCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewEventCountCache(c cache.MultiGetDeleteCountCache) EventCounterCache {
	return &eventCounterCache{cache: c}
}

func (c *eventCounterCache) GetEventCounts(keys []string) ([]float64, error) {
	pipe := c.cache.Pipeline()
	sCmds := make([]*goredis.StringCmd, 0, len(keys))
	for _, k := range keys {
		c := pipe.Get(k)
		sCmds = append(sCmds, c)
	}
	_, err := pipe.Exec()
	if err != nil {
		// Exec returns error of the first failed command.
		// https://pkg.go.dev/github.com/redis/go-redis/v9#Pipeline.Exec
		if err != goredis.Nil {
			return []float64{}, fmt.Errorf("err: %s, keys: %v", err.Error(), keys)
		}
	}
	return c.getEventValues(sCmds)
}

func (*eventCounterCache) getEventValues(cmds []*goredis.StringCmd) ([]float64, error) {
	eventVals := make([]float64, 0, len(cmds))
	for _, c := range cmds {
		str, err := c.Result()
		if err != nil {
			if err != goredis.Nil {
				return []float64{}, err
			}
			str = "0"
		}
		float, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return []float64{}, err
		}
		eventVals = append(eventVals, float)
	}
	return eventVals, nil
}

func (c *eventCounterCache) GetEventCountsV2(keys [][]string) ([]float64, error) {
	pipe := c.cache.Pipeline()
	stringCmds := make([][]*goredis.StringCmd, 0, len(keys))
	for _, day := range keys {
		hourlyCmds := []*goredis.StringCmd{}
		for _, hour := range day {
			c := pipe.Get(hour)
			hourlyCmds = append(hourlyCmds, c)
		}
		stringCmds = append(stringCmds, hourlyCmds)
	}
	_, err := pipe.Exec()
	if err != nil {
		// Exec returns error of the first failed command.
		// https://pkg.go.dev/github.com/redis/go-redis/v9#Pipeline.Exec
		if err != goredis.Nil {
			return []float64{}, fmt.Errorf("err: %s, keys: %v", err, keys)
		}
	}
	return c.getEventValuesV2(stringCmds)
}

func (*eventCounterCache) getEventValuesV2(cmds [][]*goredis.StringCmd) ([]float64, error) {
	eventVals := make([]float64, 0, len(cmds))
	for _, day := range cmds {
		var totalVal float64
		for _, hour := range day {
			str, err := hour.Result()
			if err != nil {
				if err != goredis.Nil {
					return []float64{}, err
				}
				str = "0"
			}
			float, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return []float64{}, err
			}
			totalVal += float
		}
		eventVals = append(eventVals, totalVal)
	}
	return eventVals, nil
}

func (c *eventCounterCache) GetUserCount(key string) (int64, error) {
	return c.cache.PFCount(key)
}

func (c *eventCounterCache) GetUserCounts(keys []string) ([]float64, error) {
	pipe := c.cache.Pipeline()
	iCmds := make([]*goredis.IntCmd, 0, len(keys))
	for _, k := range keys {
		c := pipe.PFCount(k)
		iCmds = append(iCmds, c)
	}
	_, err := pipe.Exec()
	if err != nil {
		return []float64{}, fmt.Errorf("err: %v, keys: %v", err, keys)
	}
	return c.getUserValues(iCmds)
}

func (*eventCounterCache) getUserValues(cmds []*goredis.IntCmd) ([]float64, error) {
	userVals := make([]float64, 0, len(cmds))
	for _, c := range cmds {
		val, err := c.Result()
		if err != nil {
			return []float64{}, err
		}
		userVals = append(userVals, float64(val))
	}
	return userVals, nil
}

func (c *eventCounterCache) GetUserCountsV2(
	keys []string,
) ([]float64, error) {
	count, err := c.getUserCountsV2(keys)
	if err != nil {
		return nil, fmt.Errorf(
			"err: %v, keys: %v",
			err, keys,
		)
	}
	return count, nil
}

func (c *eventCounterCache) getUserCountsV2(
	keys []string,
) ([]float64, error) {
	pipe := c.cache.Pipeline()
	iCmds := make([]*goredis.IntCmd, 0, len(keys))
	for _, k := range keys {
		c := pipe.PFCount(k)
		iCmds = append(iCmds, c)
	}
	_, err := pipe.Exec()
	if err != nil {
		return nil, err
	}
	return c.getUserValues(iCmds)
}

func (c *eventCounterCache) UpdateUserCount(key, userID string) error {
	_, err := c.cache.PFAdd(key, userID)
	if err != nil {
		return err
	}
	return nil
}

func (c *eventCounterCache) MergeMultiKeys(dest string, keys []string) error {
	if err := c.cache.PFMerge(dest, keys...); err != nil {
		return fmt.Errorf("err: %s, dest: %v, keys: %v", err, dest, keys)
	}
	return nil
}

func (c *eventCounterCache) DeleteKey(key string) error {
	if err := c.cache.Delete(key); err != nil {
		return err
	}
	return nil
}
