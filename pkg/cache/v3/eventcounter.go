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
	"strconv"

	goredis "github.com/go-redis/redis"

	"github.com/bucketeer-io/bucketeer/pkg/cache"
)

type EventCounterCache interface {
	GetEventCounts(keys []string) ([]float64, error)
	GetUserCounts(keys []string) ([]float64, error)
}

type eventCounterCache struct {
	cache cache.MultiGetDeleteCountCache
}

func NewEventCountCache(c cache.MultiGetDeleteCountCache) EventCounterCache {
	return &eventCounterCache{cache: c}
}

func (c *eventCounterCache) GetEventCounts(keys []string) ([]float64, error) {
	result, err := c.cache.GetMulti(keys)
	if err != nil {
		return []float64{}, fmt.Errorf("err: %v, keys: %v", err, keys)
	}
	return getEventValues(result)
}

func getEventValues(vals []interface{}) ([]float64, error) {
	eventVals := make([]float64, 0, len(vals))
	for _, v := range vals {
		if v == nil {
			eventVals = append(eventVals, 0)
			continue
		}
		str, ok := v.(string)
		if !ok {
			return []float64{}, fmt.Errorf("failed to cast value: %v", v)
		}
		float, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return []float64{}, err
		}
		eventVals = append(eventVals, float)
	}
	return eventVals, nil
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
	return getUserValues(iCmds)
}

func getUserValues(cmds []*goredis.IntCmd) ([]float64, error) {
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
