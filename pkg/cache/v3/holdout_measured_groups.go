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

	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/v2/pkg/cache"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	holdoutMeasuredGroupsKind             = "holdout_measured_groups"
	holdoutMeasuredGroupsGenerationSuffix = "gen"
)

// HoldoutMeasuredGroupsCache stores, per environment, the holdout groups the
// analytics pipeline measures: the groups that have an internal experiment
// and are serving or stopped recently enough for late events to count. Each
// group carries every target row, removed ones included. It is separate from
// HoldoutGroupCache, whose serving snapshot is cleared when the group stops
// and holds current targets only.
type HoldoutMeasuredGroupsCache interface {
	Get(environmentID string) (*ftproto.HoldoutGroups, error)
	// PutIfNewer stores the list unless a newer generation is cached.
	// generation is the time (UnixNano) captured before the writer read the
	// list from storage. Returns false when the write was rejected as stale;
	// that is not an error.
	PutIfNewer(groups *ftproto.HoldoutGroups, environmentID string, generation int64) (bool, error)
}

type holdoutMeasuredGroupsCache struct {
	cache cache.Cache
}

func NewHoldoutMeasuredGroupsCache(c cache.Cache) HoldoutMeasuredGroupsCache {
	return &holdoutMeasuredGroupsCache{cache: c}
}

func (c *holdoutMeasuredGroupsCache) Get(environmentID string) (*ftproto.HoldoutGroups, error) {
	value, err := c.cache.Get(c.key(environmentID))
	if err != nil {
		return nil, err
	}
	b, err := cache.Bytes(value)
	if err != nil {
		return nil, err
	}
	groups := &ftproto.HoldoutGroups{}
	if err := proto.Unmarshal(b, groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (c *holdoutMeasuredGroupsCache) PutIfNewer(
	groups *ftproto.HoldoutGroups,
	environmentID string,
	generation int64,
) (bool, error) {
	if groups == nil {
		return false, errors.New("holdout groups cannot be nil")
	}
	buffer, err := proto.Marshal(groups)
	if err != nil {
		return false, err
	}
	return c.cache.PutIfNewer(
		c.key(environmentID),
		c.generationKey(environmentID),
		buffer,
		generation,
		0,
	)
}

func (c *holdoutMeasuredGroupsCache) key(environmentID string) string {
	return fmt.Sprintf("%s:%s", environmentID, holdoutMeasuredGroupsKind)
}

// generationKey wraps the value key in a Redis hash tag so the value and the
// marker share a cluster slot.
func (c *holdoutMeasuredGroupsCache) generationKey(environmentID string) string {
	return fmt.Sprintf("{%s}:%s", c.key(environmentID), holdoutMeasuredGroupsGenerationSuffix)
}
